package services

import (
	"context"
	"fmt"
	"os"
)

// ValuationService handles valuation-related business logic
type ValuationService struct {
	coingeckoClient *CoinGeckoClient
}

// NewValuationService creates a new valuation service
func NewValuationService(coingeckoClient *CoinGeckoClient) *ValuationService {
	return &ValuationService{
		coingeckoClient: coingeckoClient,
	}
}

// GetTokenHistory retrieves price history for a token
func (s *ValuationService) GetTokenHistory(ctx context.Context, symbol string) ([]PricePoint, error) {
	priceHistory, err := s.coingeckoClient.GetPriceHistoryWithCache(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history for %s: %w", symbol, err)
	}
	return priceHistory, nil
}

// GetTokenValuation retrieves valuation metrics for a specific token
func (s *ValuationService) GetTokenValuation(ctx context.Context, symbol string, token *Token) (*ValuationData, error) {
	// Try to get from cache first
	if cachedValuation, err := GetCachedValuation(ctx, symbol); err == nil && cachedValuation != nil {
		return cachedValuation, nil
	}

	// Cache miss - compute valuation
	priceHistory, err := s.GetTokenHistory(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}

	// Fetch TVL data
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://ethereum-rpc.publicnode.com"
	}

	tvl, err := FetchTVL(ctx, symbol, token.ContractAddress, token.Decimals, rpcURL)
	if err != nil {
		// Continue with TVL = 0 rather than failing completely
		tvl = 0
	}

	// Calculate valuation
	valuation, err := CalculateValuation(ctx, symbol, priceHistory, tvl)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate valuation: %w", err)
	}

	// Cache the result
	if cacheErr := SetCachedValuation(ctx, symbol, *valuation); cacheErr != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to cache valuation for %s: %v\n", symbol, cacheErr)
	}

	return valuation, nil
}

// GetAllTokenValuations retrieves valuation metrics for all tokens
func (s *ValuationService) GetAllTokenValuations(ctx context.Context, tokens []Token) ([]ValuationData, error) {
	var valuations []ValuationData

	for _, token := range tokens {
		valuation, err := s.GetTokenValuation(ctx, token.Symbol, &token)
		if err != nil {
			// Log error but continue with other tokens
			fmt.Printf("Error getting valuation for %s: %v\n", token.Symbol, err)
			continue
		}

		valuations = append(valuations, *valuation)
	}

	return valuations, nil
}
