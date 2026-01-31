package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distroanalyzer/analyze"
	"distroanalyzer/cache"
	"distroanalyzer/collect"
	"distroanalyzer/explain"
	"distroanalyzer/httpapi"
	"distroanalyzer/score"
	"distroanalyzer/store"
)

func main() {
	// 1. Cargar configuración
	cfg := loadConfig()

	log.Printf("Starting DistroAnalyzer on %s", cfg.ServerAddr)

	// 2. Inicializar componentes
	components, err := initComponents(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize components: %v", err)
	}
	defer components.cleanup()

	// 3. Crear handler HTTP
	handler, err := httpapi.NewHandler(
		components.collector,
		components.analyzer,
		components.engine,
		components.explainer,
		components.cache,
		components.store,
		cfg.TemplatesDir,
	)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	// 4. Crear router
	router := httpapi.NewRouter(handler, cfg.StaticDir)

	// 5. Configurar servidor HTTP
	server := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Iniciar servidor en goroutine
	go func() {
		log.Printf("Server listening on %s (Using Cerebras API)", cfg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 7. Esperar señal de shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 8. Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// Config contiene la configuración de la aplicación.
type Config struct {
	ServerAddr      string
	TemplatesDir    string
	StaticDir       string
	DBPath          string
	RedisAddr       string
	RedisPass       string
	RedisDB         int
	GithubToken     string
	CerebrasAPIKey  string // Cambiado de Gemini
	CerebrasModel   string // Cambiado de Gemini
	UseRedis        bool
}

// loadConfig carga la configuración desde variables de entorno.
func loadConfig() *Config {
	return &Config{
		ServerAddr:      getEnv("SERVER_ADDR", ":8080"),
		TemplatesDir:    getEnv("TEMPLATES_DIR", "./web/templates"),
		StaticDir:       getEnv("STATIC_DIR", "./web/static"),
		DBPath:          getEnv("DB_PATH", "./data/distroanalyzer.db"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:       getEnv("REDIS_PASSWORD", ""),
		RedisDB:         0,
		GithubToken:     getEnv("GITHUB_TOKEN", ""),
		CerebrasAPIKey:  getEnv("CEREBRAS_API_KEY", ""), // Busca la nueva variable
		CerebrasModel:   getEnv("CEREBRAS_MODEL", "llama3.1-8b"), // Modelo por defecto de Cerebras
		UseRedis:        getEnv("USE_REDIS", "false") == "true",
	}
}

// Components agrupa todos los componentes inicializados.
type Components struct {
	collector collect.Collector
	analyzer  analyze.Analyzer
	engine    *score.Engine
	explainer explain.Explainer
	cache     cache.Cache
	store     store.Store
}

func (c *Components) cleanup() {
	if c.store != nil {
		if closer, ok := c.store.(interface{ Close() error }); ok {
			closer.Close()
		}
	}
	if c.cache != nil {
		if closer, ok := c.cache.(interface{ Close() error }); ok {
			closer.Close()
		}
	}
}

// initComponents inicializa todos los componentes del sistema.
func initComponents(cfg *Config) (*Components, error) {
	// 1. Collector (GitHub)
	collector := collect.NewGitHubCollector(cfg.GithubToken)

	// 2. Analyzer (Cerebras)
	if cfg.CerebrasAPIKey == "" {
		log.Println("WARNING: CEREBRAS_API_KEY not set, analysis will fail")
	}

	// Usamos el constructor que actualizamos en ai.go
	analyzer, err := analyze.NewAIAnalyzer(cfg.CerebrasAPIKey, cfg.CerebrasModel)
	if err != nil {
		return nil, err
	}

	// 3. Scoring engine
	distros := score.Top50Distros()
	engine := score.NewEngine(distros)

	// 4. Explainer
	explainer := explain.NewSimpleExplainer()

	// 5. Cache
	var cacheImpl cache.Cache
	if cfg.UseRedis {
		log.Printf("Using Redis cache at %s", cfg.RedisAddr)
		redisCache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
		if err != nil {
			log.Printf("Redis connection failed, falling back to memory cache: %v", err)
			cacheImpl = cache.NewMemoryCache()
		} else {
			cacheImpl = redisCache
		}
	} else {
		log.Println("Using in-memory cache")
		cacheImpl = cache.NewMemoryCache()
	}

	// 6. Store
	log.Printf("Using SQLite database at %s", cfg.DBPath)
	storeImpl, err := store.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	return &Components{
		collector: collector,
		analyzer:  analyzer,
		engine:    engine,
		explainer: explainer,
		cache:     cacheImpl,
		store:     storeImpl,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
