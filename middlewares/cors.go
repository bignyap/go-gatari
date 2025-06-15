package middlewares

import (
	"log"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the Authorization header
		token := r.Header.Get("Authorization")
		if token != "valid-token" { // Replace "valid-token" with your actual token logic
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Call the next handler if authenticated
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request method and URL
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
