package httpapi

import (
	"net/http"
	"time"
	"context"
)

// TimeoutMiddleware agrega timeout a las peticiones.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Crear context con timeout
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			// Reemplazar request context
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware agrega headers CORS.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
