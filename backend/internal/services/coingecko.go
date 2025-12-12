package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CoinGeckoClient handles API calls to CoinGecko
type CoinGeckoClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// PricePoint represents a single price data point
type PricePoint struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

// PriceHistory represents the response from CoinGecko market chart API
type PriceHistory struct {
	Prices [][]interface{} `json:"prices"` // [[timestamp, price], ...]
}

// NewCoinGeckoClient creates a new CoinGecko API client
func NewCoinGeckoClient(apiKey string) *CoinGeckoClient {
	return &CoinGeckoClient{
		apiKey:  apiKey,
		baseURL: "https://api.coingecko.com/api/v3",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// symbolToCoinGeckoID maps our token symbols to CoinGecko coin IDs
var symbolToCoinGeckoID = map[string]string{
	"wstETH":  "wrapped-steth",
	"ankrETH": "ankreth",
	"rETH":    "rocket-pool-eth",
	"wBETH":   "wrapped-beacon-eth",
	"pufETH":  "pufeth",
	"LSETH":   "liquid-staked-ethereum",
	"RSETH":   "kelp-dao-restaked-eth",
	"METH":    "mantle-staked-ether",
	"CBETH":   "coinbase-wrapped-staked-eth",
	"TETH":    "treehouse-eth",
	"SFRXETH": "staked-frax-ether",
	"CDCETH":  "crypto-com-staked-eth",
	"UNIETH":  "universal-eth",
}

// GetCoinGeckoID returns the CoinGecko ID for a given symbol
func (c *CoinGeckoClient) GetCoinGeckoID(symbol string) (string, error) {
	id, exists := symbolToCoinGeckoID[symbol]
	if !exists {
		return "", fmt.Errorf("unsupported token symbol: %s", symbol)
	}
	return id, nil
}

// GetPriceHistory fetches 1-year price history for a token
func (c *CoinGeckoClient) GetPriceHistory(symbol string) ([]PricePoint, error) {
	coinID, err := c.GetCoinGeckoID(symbol)
	if err != nil {
		return nil, err
	}

	// CoinGecko market chart endpoint for 1 year of daily data
	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=eth&days=365&interval=daily", c.baseURL, coinID)

	// Add API key if available
	if c.apiKey != "" {
		url += "&x_cg_demo_api_key=" + c.apiKey
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CoinGecko API error (status %d): %s", resp.StatusCode, string(body))
	}

	var history PriceHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert the prices array to our PricePoint format
	var pricePoints []PricePoint
	for _, price := range history.Prices {
		if len(price) >= 2 {
			timestamp, ok1 := price[0].(float64)
			priceValue, ok2 := price[1].(float64)
			if ok1 && ok2 {
				pricePoints = append(pricePoints, PricePoint{
					Timestamp: int64(timestamp),
					Price:     priceValue,
				})
			}
		}
	}

	return pricePoints, nil
}
