package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

// Define a custom unexported type for the context key to avoid collisions
type contextKey string

// A variable for the userID key
const UserIDKey contextKey = "userID"

func authorizeRequest(r *http.Request) (bool, *http.Request) {
	// Get authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, r
	}

	// Verify authorization header has a valid format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return false, r
	}
	tokenStr := parts[1]

	// Verify token and extract user ID
	userID, err := VerifyToken(tokenStr)
	if err != nil {
		log.Println("Error verifying token:", err)
		return false, r
	}

	// Create a new context with userID attached to it
	newCtx := context.WithValue(r.Context(), UserIDKey, userID)
	// Create a new request with the updated context
	newReq := r.WithContext(newCtx)

	return true, newReq
}

// checkAuthorization determines whether the given HTTP request
// requires authorization check.
func checkAuthorization(r *http.Request) bool {
	// Add more rules here
	return r.Method == "GET" && r.URL.Path == "/protected"
}

// authMiddleware checks and authorizes requests.
// Unauthorized requests does not cross this middleweare
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If request requires authorization
		if checkAuthorization(r) {
			// Check authorization and get the potentially new request
			authorized, newReq := authorizeRequest(r)
			if !authorized {
				// Logging
				log.Printf("Unauthorized request for %s %s", r.Method, r.URL.Path)

				// Write unauthorised response
				data := ErrorResponse{"Unauthorized"}
				WriteJSON(w, http.StatusUnauthorized, data)
				return
			}
			// Update the request with the updated context
			r = newReq
		}

		// Forward request to the next handler
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs the time taken for request completion
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start recording time
		start := time.Now()

		// Forward request to the next handler
		next.ServeHTTP(w, r)

		// Log time taken
		log.Printf("Time taken for %s %s is %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// StackMiddlewares contains all the middlewares
// stacked in their relative execution order.
//
// Returns a new handler with middlewares stacked on top of it.
func StackMiddlewares(handler http.Handler) http.Handler {
	// Middleware execution order
	// Request -> Logging -> Auth -> Main Handler

	// Stack middlewares in reverse order
	// First middleware is stacked last
	h := authMiddleware(handler)
	h = loggingMiddleware(h)
	return h
}
