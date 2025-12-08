// Package main ETH Staking Analytics Backend API
//
// This is the backend API for the ETH Staking Analytics application.
// It provides endpoints for retrieving LST token data, price history, and valuation metrics.
//
// Terms Of Service: https://github.com/Haxsen/HxnETHstakingAnalyticsApp
//
// Schemes: http, https
// Host: localhost:8080
// BasePath: /api
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/cache"
	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/db"
	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/services"
	_ "github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/docs" // Generated swagger docs
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var coingeckoClient *services.CoinGeckoClient

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Initialize Redis cache
	if err := cache.InitRedis(); err != nil {
		log.Printf("Failed to initialize Redis (continuing without cache): %v", err)
	} else {
		defer cache.CloseRedis()
	}

	// Initialize CoinGecko client
	apiKey := os.Getenv("COINGECKO_API_KEY")
	coingeckoClient = services.NewCoinGeckoClient(apiKey)
	log.Println("CoinGecko client initialized")

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize router
	r := chi.NewRouter()

	// CORS middleware
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000,http://127.0.0.1:3000"
	}
	allowedOrigins := strings.Split(corsOrigins, ",")

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/tokens", getTokensHandler)
		r.Get("/token/{id}/history", getTokenHistoryHandler)
		r.Get("/token/{id}/valuation", getTokenValuationHandler)
		r.Post("/cache/refresh", refreshCacheHandler)
	})

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// getTokensHandler returns all active tokens
//
// @Summary Get all tracked LST tokens
// @Description Retrieve a list of all active Liquid Staking Tokens being tracked
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "tokens: array of token objects, count: number of tokens"
// @Failure 500 {object} map[string]string "error: error message"
// @Router /api/tokens [get]
func getTokensHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokens, err := db.GetAllTokens()
	if err != nil {
		log.Printf("Error fetching tokens: %v", err)
		http.Error(w, `{"error": "Failed to fetch tokens"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tokens": tokens,
		"count":  len(tokens),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

// getTokenHistoryHandler returns price history for a token
//
// @Summary Get price history for a token
// @Description Retrieve 1-year price history for a specific LST token
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenSymbol path string true "Token symbol (e.g., wstETH, rETH)"
// @Success 200 {object} map[string]interface{} "price_history: array of price points, count: number of data points"
// @Failure 400 {object} map[string]string "error: invalid token symbol"
// @Failure 500 {object} map[string]string "error: failed to fetch price data"
// @Router /api/token/{tokenSymbol}/history [get]
func getTokenHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenSymbol := chi.URLParam(r, "id") // Using 'id' param but it's actually the symbol

	// Validate that the token exists in our database
	_, err := db.GetTokenBySymbol(tokenSymbol)
	if err != nil {
		http.Error(w, `{"error": "Token not found or not supported"}`, http.StatusBadRequest)
		return
	}

	// Fetch price history with caching
	ctx := r.Context()
	priceHistory, err := coingeckoClient.GetPriceHistoryWithCache(ctx, tokenSymbol)
	if err != nil {
		log.Printf("Error fetching price history for %s: %v", tokenSymbol, err)
		http.Error(w, `{"error": "Failed to fetch price history"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token_symbol": tokenSymbol,
		"price_history": priceHistory,
		"count": len(priceHistory),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func getTokenValuationHandler(w http.ResponseWriter, r *http.Request) {
	tokenID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Token valuation for ` + tokenID + ` - TODO"}`))
}

func refreshCacheHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Cache refresh - TODO"}`))
}
