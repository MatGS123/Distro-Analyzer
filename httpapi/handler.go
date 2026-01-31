// Package httpapi maneja las peticiones HTTP del sistema.
//
// Los handlers orquestan el flujo entre componentes pero no contienen lógica de negocio.
package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"distroanalyzer/analyze"
	"distroanalyzer/cache"
	"distroanalyzer/collect"
	"distroanalyzer/explain"
	"distroanalyzer/profile"
	"distroanalyzer/score"
	"distroanalyzer/store"
)

// Handler maneja las peticiones HTTP.
type Handler struct {
	collector  collect.Collector
	analyzer   analyze.Analyzer
	engine     *score.Engine
	explainer  explain.Explainer
	cache      cache.Cache
	store      store.Store
	templates  *template.Template
}

// NewHandler crea un nuevo handler HTTP.
func NewHandler(
	collector collect.Collector,
	analyzer analyze.Analyzer,
	engine *score.Engine,
	explainer explain.Explainer,
	cache cache.Cache,
	store store.Store,
	templatesDir string,
) (*Handler, error) {

	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 {
			return a * b
		},
	}

	tmpl, err := template.New("").
	Funcs(funcMap).
	ParseGlob(templatesDir + "/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Handler{
		collector: collector,
		analyzer:  analyzer,
		engine:    engine,
		explainer: explainer,
		cache:     cache,
		store:     store,
		templates: tmpl,
	}, nil
}

// Home muestra el formulario principal.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Analyze procesa un perfil de GitHub.
func (h *Handler) Analyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Verificar cache
	cacheKey := "profile:" + username
	cached, err := h.cache.Get(ctx, cacheKey)
	if err != nil {
		log.Printf("cache error: %v", err)
	}

	if cached != nil {
		log.Printf("cache hit for %s", username)
		h.renderResult(w, cached)
		return
	}

	// 2. Ejecutar pipeline completo
	prof, err := h.runPipeline(ctx, username)
	if err != nil {
		log.Printf("pipeline error for %s: %v", username, err)
		http.Error(w, fmt.Sprintf("Analysis failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 3. Guardar en cache (1 hora TTL)
	if err := h.cache.Set(ctx, cacheKey, prof, 1*time.Hour); err != nil {
		log.Printf("failed to cache profile: %v", err)
	}

	// 4. Persistir en DB
	if err := h.store.Save(ctx, prof); err != nil {
		log.Printf("failed to save profile: %v", err)
	}

	// 5. Renderizar resultado
	h.renderResult(w, prof)
}

// runPipeline ejecuta el flujo completo de análisis.
func (h *Handler) runPipeline(ctx context.Context, username string) (*profile.Profile, error) {
	// 1. Collect
	rawData, err := h.collector.Collect(username)
	if err != nil {
		return nil, fmt.Errorf("collection failed: %w", err)
	}

	// 2. Analyze
	signals, err := h.analyzer.Analyze(rawData)
		if err != nil {
			return nil, fmt.Errorf("analysis failed: %w", err)
		}

		// 3. Score
		scoreOut := h.engine.Score(signals)

		// 4. Explain
		explanation := h.explainer.Explain(scoreOut.Result, signals)
		scoreOut.Result.Explanation = explanation

		// 5. Construir Profile completo
		prof := &profile.Profile{
			Username:  username,
			Source:    "github",
			RawData:   *rawData,
			Signals:   *signals,
			Result:    *scoreOut.Result,
			CreatedAt: time.Now(),
		}

		prof.Recommendation = profile.Recommendation{
			DistroID:   scoreOut.BestDistroID,
			DistroName: scoreOut.BestDistroName,
		}

		// 6. Limpiar datos pesados
		prof.ClearLargeData()

		return prof, nil
}

// renderResult renderiza el template de resultado.
func (h *Handler) renderResult(w http.ResponseWriter, p *profile.Profile) {
	data := map[string]interface{}{
		"Profile": p,
	}

	if err := h.templates.ExecuteTemplate(w, "result.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// History muestra perfiles analizados previamente.
func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	profiles, err := h.store.List(ctx, 20, 0)
	if err != nil {
		log.Printf("failed to list profiles: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Profiles": profiles,
	}

	if err := h.templates.ExecuteTemplate(w, "history.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// API endpoint para respuestas JSON.
func (h *Handler) AnalyzeJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Verificar cache
	cacheKey := "profile:" + req.Username
	cached, err := h.cache.Get(ctx, cacheKey)
	if err != nil {
		log.Printf("cache error: %v", err)
	}

	var prof *profile.Profile
	if cached != nil {
		prof = cached
	} else {
		prof, err = h.runPipeline(ctx, req.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.cache.Set(ctx, cacheKey, prof, 1*time.Hour)
		h.store.Save(ctx, prof)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prof)
}
