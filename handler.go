package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var productStore []Product = make([]Product, 0)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	// Encode the response to JSON into a buffer first
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// Log error
		log.Printf("Error while encoding to JSON. Data: %+v. Error: %v.\n", data, err)

		// Set error headers and send generic error
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Set content type and write status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonBytes)
}

type HandlerSet struct{}

func NewHandlerSet() *HandlerSet {
	return &HandlerSet{}
}

func (h *HandlerSet) Login(w http.ResponseWriter, r *http.Request) {
	// Dummy user login
	loginSuccess := rand.Intn(2)
	if loginSuccess == 0 {
		data := ErrorResponse{"Invalid credentials"}
		WriteJSON(w, http.StatusUnauthorized, data)
		return
	}

	// Generate access and refresh token
	accessToken, err := NewAccessToken(999)
	if err != nil {
		log.Println("Error generating access token:", err)
		data := ErrorResponse{"Something went wrong"}
		WriteJSON(w, http.StatusInternalServerError, data)
		return
	}
	refreshToken, err := NewRefreshToken(999)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		data := ErrorResponse{"Something went wrong"}
		WriteJSON(w, http.StatusInternalServerError, data)
		return
	}

	// Send token response
	response := TokenResponse{
		Access:  accessToken,
		Refresh: refreshToken,
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *HandlerSet) Protected(w http.ResponseWriter, r *http.Request) {
	// Get authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		data := ErrorResponse{"Invalid credentials"}
		WriteJSON(w, http.StatusUnauthorized, data)
		return
	}

	// Verify authorization header has a valid format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		data := ErrorResponse{"Invalid credentials"}
		WriteJSON(w, http.StatusUnauthorized, data)
		return
	}
	tokenStr := parts[1]

	// Verify token and extract user ID
	userID, err := VerifyToken(tokenStr)
	if err != nil {
		log.Println("Error verifying token:", err)
		data := ErrorResponse{"Invalid credentials"}
		WriteJSON(w, http.StatusUnauthorized, data)
		return
	}

	WriteJSON(w, http.StatusOK, userID)
}

func (h *HandlerSet) AllProducts(w http.ResponseWriter, r *http.Request) {
	// Extract 'sort' query parameter
	queryParams := r.URL.Query()
	sortValue := queryParams.Get("sort")

	// Create copy of products
	products := make([]Product, len(productStore))
	copy(products, productStore)

	// Implement sort
	switch sortValue {
	case "id":
		sort.Sort(ByID(products))
	case "name":
		sort.Sort(ByName(products))
	}

	WriteJSON(w, http.StatusOK, products)
}

func (h *HandlerSet) NotFound(w http.ResponseWriter, r *http.Request) {
	data := ErrorResponse{"Not found"}
	WriteJSON(w, http.StatusNotFound, data)
}

func (h *HandlerSet) AddProduct(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		response := ErrorResponse{Message: "Invalid request"}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Add product
	if len(productStore) == 0 {
		product.ID = 1
	} else {
		product.ID = productStore[len(productStore)-1].ID + 1
	}
	productStore = append(productStore, product)

	WriteJSON(w, http.StatusCreated, product)
}

func (h *HandlerSet) ProductByID(w http.ResponseWriter, r *http.Request) {
	// Parse path parameter "id"
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		response := ErrorResponse{Message: "Invalid request"}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Get product matching the "id"
	var product *Product
	for _, p := range productStore {
		if p.ID == id {
			product = &p
			break
		}
	}
	if product == nil {
		response := ErrorResponse{Message: "Not found"}
		WriteJSON(w, http.StatusNotFound, response)
		return
	}

	WriteJSON(w, http.StatusOK, *product)
}

func (h *HandlerSet) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse path parameter "id"
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		response := ErrorResponse{Message: "Invalid request"}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Get product matching the "id"
	var productIdx int = -1
	for i, p := range productStore {
		if p.ID == id {
			productIdx = i
			break
		}
	}
	if productIdx == -1 {
		response := ErrorResponse{Message: "Not found"}
		WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Parse request body
	var newProduct Product
	err = json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		response := ErrorResponse{Message: "Invalid request"}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()
	newProduct.ID = id

	// Update product
	productStore[productIdx] = newProduct
	WriteJSON(w, http.StatusCreated, newProduct)
}

func (h *HandlerSet) UpdateProductPartial(w http.ResponseWriter, r *http.Request) {
	// Parse path parameter "id"
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		response := ErrorResponse{Message: "Invalid request"}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Get product matching the "id"
	var productIdx int = -1
	for i, p := range productStore {
		if p.ID == id {
			productIdx = i
			break
		}
	}
	if productIdx == -1 {
		response := ErrorResponse{Message: "Not found"}
		WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Parse request body
	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		response := ErrorResponse{Message: "Invalid request"}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Update product
	product := &productStore[productIdx]
	for field, value := range updates {
		switch field {
		case "name":
			if name, ok := value.(string); ok {
				product.Name = name
			}
		case "description":
			if desc, ok := value.(string); ok {
				product.Description = desc
			}
		case "price":
			if price, ok := value.(float64); ok {
				product.Price = price
			}
		}
	}

	WriteJSON(w, http.StatusOK, product)
}

func (h *HandlerSet) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Parse path parameter "id"
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		response := ErrorResponse{Message: "Invalid request"}
		WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Get product matching the "id"
	var productIdx int = -1
	for i, p := range productStore {
		if p.ID == id {
			productIdx = i
			break
		}
	}
	if productIdx == -1 {
		response := ErrorResponse{Message: "Not found"}
		WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Delete product
	product := productStore[productIdx]
	productStore = append(productStore[:productIdx], productStore[productIdx+1:]...)

	WriteJSON(w, http.StatusOK, product)
}
