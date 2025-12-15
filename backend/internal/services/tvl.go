package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/cache"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TVLData represents TVL information for a token
type TVLData struct {
	TokenSymbol string    `json:"token_symbol"`
	TVL         float64   `json:"tvl"`
	LastUpdated time.Time `json:"last_updated"`
}

// CachedTVLData represents cached TVL data
type CachedTVLData struct {
	Data       TVLData    `json:"data"`
	CachedAt   time.Time  `json:"cached_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
}

// TVLFetcher handles TVL data fetching from blockchain
type TVLFetcher struct {
	ethClient *ethclient.Client
}

// NewTVLFetcher creates a new TVL fetcher
func NewTVLFetcher(rpcURL string) (*TVLFetcher, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	return &TVLFetcher{
		ethClient: client,
	}, nil
}

// ERC20 ABI for totalSupply function
const erc20ABI = `[{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

// FetchTVLFromContract fetches TVL by calling totalSupply on the token contract
func (t *TVLFetcher) FetchTVLFromContract(ctx context.Context, contractAddress string, decimals int) (float64, error) {
	// Parse the contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Call totalSupply function
	address := common.HexToAddress(contractAddress)
	callData, err := parsedABI.Pack("totalSupply")
	if err != nil {
		return 0, fmt.Errorf("failed to pack call data: %w", err)
	}

	result, err := t.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: callData,
	}, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to call contract: %w", err)
	}

	// Unpack the result
	outputs, err := parsedABI.Unpack("totalSupply", result)
	if err != nil {
		return 0, fmt.Errorf("failed to unpack result: %w", err)
	}

	if len(outputs) == 0 {
		return 0, fmt.Errorf("no outputs from totalSupply call")
	}

	totalSupply, ok := outputs[0].(*big.Int)
	if !ok {
		return 0, fmt.Errorf("unexpected output type from totalSupply")
	}

	// Convert to float with proper decimal adjustment
	tvlFloat := new(big.Float).SetInt(totalSupply)
	decimalsMultiplier := new(big.Float).SetFloat64(math.Pow(10, float64(decimals)))
	tvlFloat.Quo(tvlFloat, decimalsMultiplier)

	tvl, _ := tvlFloat.Float64()
	return tvl, nil
}

// GetCachedTVL retrieves TVL data from cache
func GetCachedTVL(ctx context.Context, symbol string) (*TVLData, error) {
	cacheKey := fmt.Sprintf("tvl:%s", symbol)

	cachedData, err := cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, nil // Cache miss
	}

	var cached CachedTVLData
	if err := json.Unmarshal([]byte(cachedData), &cached); err != nil {
		return nil, nil // Invalid cache data
	}

	// Check if cache is expired
	if time.Now().After(cached.ExpiresAt) {
		cache.Delete(ctx, cacheKey)
		return nil, nil
	}

	return &cached.Data, nil
}

// SetCachedTVL stores TVL data in cache
func SetCachedTVL(ctx context.Context, symbol string, data TVLData) error {
	cacheDurationStr := os.Getenv("TVL_CACHE_DURATION")
	cacheDuration := 5 * time.Minute
	if cacheDurationStr != "" {
		if parsed, err := time.ParseDuration(cacheDurationStr); err == nil {
			cacheDuration = parsed
		}
	}

	cacheKey := fmt.Sprintf("tvl:%s", symbol)

	cached := CachedTVLData{
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

// FetchTVL fetches TVL data with caching
func FetchTVL(ctx context.Context, symbol, contractAddress string, decimals int, rpcURL string) (float64, error) {
	// Try to get from cache first
	if cachedData, err := GetCachedTVL(ctx, symbol); err == nil && cachedData != nil {
		return cachedData.TVL, nil
	}

	// Cache miss - fetch from blockchain
	fetcher, err := NewTVLFetcher(rpcURL)
	if err != nil {
		return 0, fmt.Errorf("failed to create TVL fetcher: %w", err)
	}

	tvl, err := fetcher.FetchTVLFromContract(ctx, contractAddress, decimals)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch TVL from contract: %w", err)
	}

	// Cache the result
	tvlData := TVLData{
		TokenSymbol: symbol,
		TVL:         tvl,
		LastUpdated: time.Now(),
	}

	if cacheErr := SetCachedTVL(ctx, symbol, tvlData); cacheErr != nil {
		// Log cache error but don't fail the request
		fmt.Printf("Warning: failed to cache TVL for %s: %v\n", symbol, cacheErr)
	}

	return tvl, nil
}
