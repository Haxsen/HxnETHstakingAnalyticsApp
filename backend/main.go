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
	"log"
	"os"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Create server configuration
	cfg := &server.Config{
		Port:               os.Getenv("PORT"),
		CORSAllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
		CoinGeckoAPIKey:    os.Getenv("COINGECKO_API_KEY"),
		EthereumRPCURL:     os.Getenv("ETHEREUM_RPC_URL"),
	}

	// Create and start server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	// Start the server (this blocks)
	log.Fatal(srv.Start())
}
