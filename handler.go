package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

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

func (h *HandlerSet) AllProducts(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, productStore)
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
