package main

type ErrorResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// ByID implements sort.Interface for []Product based on the ID field.
type ByID []Product

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID } // Ascending ID sort

// ByName implements sort.Interface for []Product based on the Name field.
type ByName []Product

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name } // Ascending Name sort
