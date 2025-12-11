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
		r.Get("/valuations", getAllValuationsHandler)
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

// getTokenValuationHandler returns valuation metrics for a specific token
//
// @Summary Get valuation metrics for a token
// @Description Retrieve APR, stability, TVL, and valuation remarks for a specific LST token
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenSymbol path string true "Token symbol (e.g., wstETH, rETH)"
// @Success 200 {object} services.ValuationData "valuation metrics for the token"
// @Failure 400 {object} map[string]string "error: invalid token symbol"
// @Failure 500 {object} map[string]string "error: failed to calculate valuation"
// @Router /api/token/{tokenSymbol}/valuation [get]
func getTokenValuationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenSymbol := chi.URLParam(r, "id")

	// Validate that the token exists in our database
	token, err := db.GetTokenBySymbol(tokenSymbol)
	if err != nil {
		http.Error(w, `{"error": "Token not found or not supported"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Try to get from cache first
	if cachedValuation, err := services.GetCachedValuation(ctx, tokenSymbol); err == nil && cachedValuation != nil {
		if err := json.NewEncoder(w).Encode(cachedValuation); err != nil {
			log.Printf("Error encoding cached response: %v", err)
			http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		}
		return
	}

	// Cache miss - compute valuation
	priceHistory, err := coingeckoClient.GetPriceHistoryWithCache(ctx, tokenSymbol)
	if err != nil {
		log.Printf("Error fetching price history for %s: %v", tokenSymbol, err)
		http.Error(w, `{"error": "Failed to fetch price history"}`, http.StatusInternalServerError)
		return
	}

	// Fetch TVL data
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://ethereum-rpc.publicnode.com"
	}

	tvl, err := services.FetchTVL(ctx, tokenSymbol, token.ContractAddress, token.Decimals, rpcURL)
	if err != nil {
		log.Printf("Error fetching TVL for %s: %v", tokenSymbol, err)
		// Continue with TVL = 0 rather than failing completely
		tvl = 0
	}

	// Calculate valuation
	valuation, err := services.CalculateValuation(ctx, tokenSymbol, priceHistory, tvl)
	if err != nil {
		log.Printf("Error calculating valuation for %s: %v", tokenSymbol, err)
		http.Error(w, `{"error": "Failed to calculate valuation"}`, http.StatusInternalServerError)
		return
	}

	// Cache the result
	if cacheErr := services.SetCachedValuation(ctx, tokenSymbol, *valuation); cacheErr != nil {
		log.Printf("Warning: failed to cache valuation for %s: %v", tokenSymbol, cacheErr)
	}

	if err := json.NewEncoder(w).Encode(valuation); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

// getAllValuationsHandler returns valuation metrics for all tokens
//
// @Summary Get valuation metrics for all tokens
// @Description Retrieve APR, stability, TVL, and valuation remarks for all tracked LST tokens (sortable table data)
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "valuations: array of valuation objects, count: number of valuations"
// @Failure 500 {object} map[string]string "error: failed to fetch valuations"
// @Router /api/valuations [get]
func getAllValuationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get all tokens
	tokens, err := db.GetAllTokens()
	if err != nil {
		log.Printf("Error fetching tokens: %v", err)
		http.Error(w, `{"error": "Failed to fetch tokens"}`, http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://ethereum-rpc.publicnode.com"
	}

	var valuations []services.ValuationData

	// Calculate valuation for each token
	for _, token := range tokens {
		// Try cache first
		if cachedValuation, err := services.GetCachedValuation(ctx, token.Symbol); err == nil && cachedValuation != nil {
			valuations = append(valuations, *cachedValuation)
			continue
		}

		// Cache miss - compute valuation
		priceHistory, err := coingeckoClient.GetPriceHistoryWithCache(ctx, token.Symbol)
		if err != nil {
			log.Printf("Error fetching price history for %s: %v", token.Symbol, err)
			continue // Skip this token but continue with others
		}

		tvl, err := services.FetchTVL(ctx, token.Symbol, token.ContractAddress, token.Decimals, rpcURL)
		if err != nil {
			log.Printf("Error fetching TVL for %s: %v", token.Symbol, err)
			tvl = 0 // Continue with TVL = 0
		}

		valuation, err := services.CalculateValuation(ctx, token.Symbol, priceHistory, tvl)
		if err != nil {
			log.Printf("Error calculating valuation for %s: %v", token.Symbol, err)
			continue // Skip this token but continue with others
		}

		// Cache the result
		if cacheErr := services.SetCachedValuation(ctx, token.Symbol, *valuation); cacheErr != nil {
			log.Printf("Warning: failed to cache valuation for %s: %v", token.Symbol, cacheErr)
		}

		valuations = append(valuations, *valuation)
	}

	response := map[string]interface{}{
		"valuations": valuations,
		"count":      len(valuations),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func refreshCacheHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Cache refresh - TODO"}`))
}
