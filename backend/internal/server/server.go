package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/api"
	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/cache"
	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/db"
	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/services"
	_ "github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/docs" // Generated swagger docs
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server holds all dependencies and provides HTTP server functionality
type Server struct {
	router          *chi.Mux
	handler         *api.Handler
	coingeckoClient *services.CoinGeckoClient
	port            string
}

// Config holds server configuration
type Config struct {
	Port                string
	CORSAllowedOrigins  string
	CoinGeckoAPIKey     string
	EthereumRPCURL      string
}

// NewServer creates a new server with all dependencies injected
func NewServer(cfg *Config) (*Server, error) {
	// Initialize database
	if err := db.InitDB(); err != nil {
		return nil, err
	}

	// Initialize Redis cache
	if err := cache.InitRedis(); err != nil {
		log.Printf("Failed to initialize Redis (continuing without cache): %v", err)
	}

	// Initialize CoinGecko client
	coingeckoClient := services.NewCoinGeckoClient(cfg.CoinGeckoAPIKey)
	log.Println("CoinGecko client initialized")

	// Initialize services
	tokenService := services.NewTokenService()
	valuationService := services.NewValuationService(coingeckoClient)

	// Initialize API handlers
	handler := api.NewHandler(tokenService, valuationService)

	// Set default port
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	// Initialize router
	r := chi.NewRouter()

	server := &Server{
		router:          r,
		handler:         handler,
		coingeckoClient: coingeckoClient,
		port:            port,
	}

	// Setup middleware and routes
	server.setupMiddleware(cfg)
	server.setupRoutes()

	return server, nil
}

// setupMiddleware configures middleware for the server
func (s *Server) setupMiddleware(cfg *Config) {
	// CORS middleware
	corsOrigins := cfg.CORSAllowedOrigins
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000,http://127.0.0.1:3000"
	}
	allowedOrigins := strings.Split(corsOrigins, ",")

	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Standard middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Swagger UI
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+s.port+"/swagger/doc.json"),
	))

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/tokens", s.handler.GetTokensHandler)
		r.Get("/token/{id}/history", s.handler.GetTokenHistoryHandler)
		r.Get("/token/{id}/valuation", s.handler.GetTokenValuationHandler)
		r.Get("/valuations", s.handler.GetAllValuationsHandler)
		r.Post("/cache/refresh", s.handler.RefreshCacheHandler)
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Server starting on port %s", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}

// Close gracefully shuts down server dependencies
func (s *Server) Close() {
	db.CloseDB()
	cache.CloseRedis()
}
