package httpapi

import (
	"net/http"
)

// NewRouter crea el router HTTP con todas las rutas.
func NewRouter(h *Handler, staticDir string) http.Handler {
	mux := http.NewServeMux()

	// Rutas HTML
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/analyze", h.Analyze)
	mux.HandleFunc("/history", h.History)

	// API JSON
	mux.HandleFunc("/api/analyze", h.AnalyzeJSON)

	// Archivos estáticos
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Middleware de logging
	return loggingMiddleware(mux)
}

// loggingMiddleware registra todas las peticiones.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request
		// log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
