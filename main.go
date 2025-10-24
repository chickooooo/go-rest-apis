package main

import (
	"log"
	"net/http"
)

func main() {
	// Create server multiplexer
	mux := http.NewServeMux()

	// Get handler set
	handlerSet := NewHandlerSet()

	// Define routes
	mux.HandleFunc("POST /login", handlerSet.Login)
	mux.HandleFunc("GET /protected", handlerSet.Protected) // protected route

	mux.HandleFunc("GET /products", handlerSet.AllProducts)
	mux.HandleFunc("POST /products", handlerSet.AddProduct)

	mux.HandleFunc("GET /products/{id}", handlerSet.ProductByID)
	mux.HandleFunc("PUT /products/{id}", handlerSet.UpdateProduct)
	mux.HandleFunc("PATCH /products/{id}", handlerSet.UpdateProductPartial)
	mux.HandleFunc("DELETE /products/{id}", handlerSet.DeleteProduct)

	// Invalid route
	mux.HandleFunc("/", handlerSet.NotFound)

	// Add middlewares
	muxWithMiddlewares := StackMiddlewares(mux)

	// Start server
	log.Println("Server started at port 8000...")
	log.Fatal(http.ListenAndServe(":8000", muxWithMiddlewares))
}
