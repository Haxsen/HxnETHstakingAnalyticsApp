package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/cache"
)

// CachedPriceHistory represents cached price history data
type CachedPriceHistory struct {
	Symbol     string      `json:"symbol"`
	Data       []PricePoint `json:"data"`
	CachedAt   time.Time    `json:"cached_at"`
	ExpiresAt  time.Time    `json:"expires_at"`
}

// GetCachedPriceHistory retrieves price history from cache
func GetCachedPriceHistory(ctx context.Context, symbol string) ([]PricePoint, error) {
	cacheKey := fmt.Sprintf("price_history:%s", symbol)

	cachedData, err := cache.Get(ctx, cacheKey)
	if err != nil {
		// Cache miss or error - not a failure
		return nil, nil
	}

	var cached CachedPriceHistory
	if err := json.Unmarshal([]byte(cachedData), &cached); err != nil {
		// Invalid cache data - treat as cache miss
		return nil, nil
	}

	// Check if cache is expired
	if time.Now().After(cached.ExpiresAt) {
		// Cache expired - remove it
		cache.Delete(ctx, cacheKey)
		return nil, nil
	}

	return cached.Data, nil
}

// SetCachedPriceHistory stores price history in cache
func SetCachedPriceHistory(ctx context.Context, symbol string, data []PricePoint) error {
	cacheDurationStr := os.Getenv("PRICE_HISTORY_CACHE_DURATION")
	cacheDuration := 1 * time.Hour
	if cacheDurationStr != "" {
		if parsed, err := time.ParseDuration(cacheDurationStr); err == nil {
			cacheDuration = parsed
		}
	}

	cacheKey := fmt.Sprintf("price_history:%s", symbol)

	cached := CachedPriceHistory{
		Symbol:    symbol,
		Data:      data,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(cacheDuration),
	}

	cachedData, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	return cache.Set(ctx, cacheKey, string(cachedData), cacheDuration)
}

// GetPriceHistoryWithCache fetches price history with caching
func (c *CoinGeckoClient) GetPriceHistoryWithCache(ctx context.Context, symbol string) ([]PricePoint, error) {
	// Try to get from cache first
	if cachedData, err := GetCachedPriceHistory(ctx, symbol); err == nil && cachedData != nil {
		return cachedData, nil
	}

	// Cache miss - fetch from API
	data, err := c.GetPriceHistory(symbol)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if cacheErr := SetCachedPriceHistory(ctx, symbol, data); cacheErr != nil {
		// Log cache error but don't fail the request
		// (we successfully got data from API)
		fmt.Printf("Warning: failed to cache price history for %s: %v\n", symbol, cacheErr)
	}

	return data, nil
}
