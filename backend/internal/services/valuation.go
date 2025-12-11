package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/cache"
)

// ValuationData represents the valuation metrics for a token
type ValuationData struct {
	TokenSymbol string    `json:"token_symbol"`
	Price       float64   `json:"price"`
	APR         float64   `json:"apr"`
	Stability   float64   `json:"stability"`
	TVL         float64   `json:"tvl"`
	Remarks     string    `json:"remarks"`
	LastUpdated time.Time `json:"last_updated"`
}

// CachedValuationData represents cached valuation data
type CachedValuationData struct {
	Data       ValuationData `json:"data"`
	CachedAt   time.Time     `json:"cached_at"`
	ExpiresAt  time.Time     `json:"expires_at"`
}

// CalculateAPR calculates the 1-year monthly average APR from price history
func CalculateAPR(priceHistory []PricePoint, symbol string) (float64, error) {
	if len(priceHistory) < 360 { // Need at least ~1 year of data for 12 months
		return 0, fmt.Errorf("insufficient price data for APR calculation")
	}

	// Sort price points by timestamp (oldest first)
	sort.Slice(priceHistory, func(i, j int) bool {
		return priceHistory[i].Timestamp < priceHistory[j].Timestamp
	})

	// Step 1: Calculate monthly averages by grouping into exactly 12 chunks
	monthlyAverages := []float64{}
	totalDays := len(priceHistory)

	// Use exactly 12 months worth of data (360 days if available)
	daysToUse := 360
	if totalDays < daysToUse {
		daysToUse = totalDays
	}

	chunkSize := daysToUse / 12 // This will be 30 for 360 days

	for i := 0; i < daysToUse; i += chunkSize {
		end := i + chunkSize
		if end > daysToUse {
			end = daysToUse
		}

		chunk := priceHistory[i:end]
		if len(chunk) == 0 {
			continue
		}

		// Calculate average price for this chunk
		sum := 0.0
		for _, point := range chunk {
			sum += point.Price
		}
		avgPrice := sum / float64(len(chunk))
		monthlyAverages = append(monthlyAverages, avgPrice)

		// Log monthly average calculation
		fmt.Printf("Token %s: Month %d average = %.6f (from %d days)\n",
			symbol, len(monthlyAverages), avgPrice, len(chunk))
	}

	// We should have approximately 12 monthly averages
	if len(monthlyAverages) < 2 {
		return 0, fmt.Errorf("insufficient monthly data for APR calculation")
	}

	// Step 2: Calculate monthly returns (12 values total)
	monthlyReturns := []float64{}

	// First month return is 0 (compared to itself)
	monthlyReturns = append(monthlyReturns, 0.0)

	// Subsequent months: difference from previous month
	for i := 1; i < len(monthlyAverages); i++ {
		monthlyReturn := monthlyAverages[i] - monthlyAverages[i-1]
		monthlyReturns = append(monthlyReturns, monthlyReturn)
	}

	// We should have 12 monthly returns
	if len(monthlyReturns) != 12 {
		return 0, fmt.Errorf("expected 12 monthly returns, got %d", len(monthlyReturns))
	}

	// Step 3: Calculate final APR as sum of all monthly returns
	apr := 0.0
	for _, monthlyReturn := range monthlyReturns {
		apr += monthlyReturn
	}

	// Log for debugging
	fmt.Printf("Token %s: %d monthly averages, %d monthly returns, total price change=%.6f\n",
		symbol, len(monthlyAverages), len(monthlyReturns), apr)

	return apr, nil
}

// getQuarter returns the quarter (1-4) for a given time
func getQuarter(t time.Time) int {
	return (int(t.Month())-1)/3 + 1
}

// CalculateStability calculates the stability score based on daily return variance
func CalculateStability(dailyReturns []float64) float64 {
	if len(dailyReturns) < 2 {
		return 1.0 // Perfect stability if insufficient data
	}

	// Calculate mean
	sum := 0.0
	for _, dailyReturn := range dailyReturns {
		sum += dailyReturn
	}
	mean := sum / float64(len(dailyReturns))

	// Calculate standard deviation
	sumSquares := 0.0
	for _, dailyReturn := range dailyReturns {
		sumSquares += math.Pow(dailyReturn-mean, 2)
	}
	stdDev := math.Sqrt(sumSquares / float64(len(dailyReturns)))

	// Calculate coefficient of variation
	if mean == 0 {
		return 0.0
	}
	coefficientOfVariation := stdDev / math.Abs(mean)

	// Convert to stability score (higher is more stable)
	// Formula: stability = 1 / (1 + coefficientOfVariation)
	stability := 1.0 / (1.0 + coefficientOfVariation)

	return stability
}

// DetermineValuationRemarks determines the valuation status based on price appreciation vs APR yield
func DetermineValuationRemarks(priceGain float64, apr float64) string {
	// Compare actual price gain over the year vs the APR (expected yield)
	// If price gain > APR, market has overpriced future yield (overvalued)
	// If price gain < APR, market has underpriced future yield (undervalued)

	if apr == 0 {
		return "Unknown"
	}

	ratio := priceGain / apr

	switch {
	case ratio < 0.5:
		return "Very Undervalued"  // Price gain much less than APR
	case ratio < 0.8:
		return "Undervalued"       // Price gain less than APR
	case ratio <= 1.2:
		return "Fair Value"        // Price gain roughly matches APR
	case ratio <= 1.5:
		return "Overvalued"        // Price gain exceeds APR
	default:
		return "Very Overvalued"   // Price gain significantly exceeds APR
	}
}

// GetExpectedPriceFromAPR calculates expected price based on APR and reference price
func GetExpectedPriceFromAPR(apr float64, referencePrice float64) float64 {
	// This is a simplified calculation - in practice, you'd use more sophisticated models
	// For now, we'll use the APR to estimate fair value
	// Expected price = reference price adjusted by APR differential
	// This is a placeholder - you might want to implement a more sophisticated model
	return referencePrice * (1 + apr) // Simplified assumption
}

// GetCachedValuation retrieves valuation data from cache
func GetCachedValuation(ctx context.Context, symbol string) (*ValuationData, error) {
	cacheKey := fmt.Sprintf("valuation:%s", symbol)

	cachedData, err := cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, nil // Cache miss
	}

	var cached CachedValuationData
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

// SetCachedValuation stores valuation data in cache for 10 minutes
func SetCachedValuation(ctx context.Context, symbol string, data ValuationData) error {
	cacheKey := fmt.Sprintf("valuation:%s", symbol)

	cached := CachedValuationData{
		Data:      data,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	cachedData, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	return cache.Set(ctx, cacheKey, string(cachedData), 10*time.Minute)
}

// CalculateValuation computes all valuation metrics for a token
func CalculateValuation(ctx context.Context, symbol string, priceHistory []PricePoint, tvl float64) (*ValuationData, error) {
	// Calculate APR
	apr, err := CalculateAPR(priceHistory, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate APR: %w", err)
	}

	// For stability, calculate daily returns from the same data used for APR
	dailyReturns := []float64{}
	for i := 1; i < len(priceHistory); i++ {
		if priceHistory[i-1].Price > 0 {
			dailyReturn := (priceHistory[i].Price / priceHistory[i-1].Price) - 1
			dailyReturns = append(dailyReturns, dailyReturn)
		}
	}

	// Calculate stability
	stability := CalculateStability(dailyReturns)

	// Get current price (latest price point)
	var currentPrice float64
	if len(priceHistory) > 0 {
		// Sort by timestamp and get latest
		sort.Slice(priceHistory, func(i, j int) bool {
			return priceHistory[i].Timestamp > priceHistory[j].Timestamp
		})
		currentPrice = priceHistory[0].Price
	}

	// Calculate last month average (most recent 30 days)
	var lastMonthAvg float64
	if len(priceHistory) >= 30 {
		// Take the most recent 30 days
		recentPrices := priceHistory[:30]
		sum := 0.0
		for _, point := range recentPrices {
			sum += point.Price
		}
		lastMonthAvg = sum / float64(len(recentPrices))
	} else {
		// Fallback to current price if insufficient data
		lastMonthAvg = currentPrice
	}

	// Calculate expected price using the new formula
	averageMonthlyReturn := apr / 12.0  // APR is sum of 12 monthly returns
	expectedPrice := (averageMonthlyReturn / 2.0) + lastMonthAvg

	// Determine valuation remarks: current price vs expected price
	remarks := determineValuationRemarks(currentPrice, expectedPrice)

	valuation := &ValuationData{
		TokenSymbol: symbol,
		Price:       currentPrice,
		APR:         apr,
		Stability:   stability,
		TVL:         tvl,
		Remarks:     remarks,
		LastUpdated: time.Now(),
	}

	return valuation, nil
}

// determineValuationRemarks implements the 5-level expected price valuation logic
func determineValuationRemarks(currentPrice, expectedPrice float64) string {
	// 5-level valuation: current price vs expected price
	// Expected price = (average_monthly_return / 2) + last_month_average

	if expectedPrice == 0 {
		return "Unknown"
	}

	// Calculate deviation percentage
	deviation := (currentPrice - expectedPrice) / expectedPrice

	// Define thresholds
	fairValueTolerance := 0.001  // ±0.1% for Fair Value
	significantThreshold := 0.01 // ±1% for "Very" levels

	switch {
	case deviation <= -significantThreshold:
		return "Very Undervalued"  // >1% below expected
	case deviation < -fairValueTolerance:
		return "Undervalued"  // 0.1%-1% below expected
	case deviation <= fairValueTolerance:
		return "Fair Value"  // ±0.1% of expected
	case deviation < significantThreshold:
		return "Overvalued"  // 0.1%-1% above expected
	default:
		return "Very Overvalued"  // >1% above expected
	}
}

// calculateMonthlyAPRs is a helper function to get monthly APRs for stability calculation
func calculateMonthlyAPRs(priceHistory []PricePoint) ([]float64, error) {
	if len(priceHistory) < 30 {
		return []float64{}, fmt.Errorf("insufficient data")
	}

	sort.Slice(priceHistory, func(i, j int) bool {
		return priceHistory[i].Timestamp < priceHistory[j].Timestamp
	})

	monthlyAPRs := []float64{}
	currentMonth := time.Unix(priceHistory[0].Timestamp/1000, 0).UTC()
	monthStartPrice := priceHistory[0].Price

	for i := 1; i < len(priceHistory); i++ {
		pointTime := time.Unix(priceHistory[i].Timestamp/1000, 0).UTC()

		if pointTime.Year() != currentMonth.Year() || pointTime.Month() != currentMonth.Month() {
			if monthStartPrice > 0 {
				monthlyReturn := (priceHistory[i-1].Price / monthStartPrice) - 1
				monthlyAPR := monthlyReturn * 12
				monthlyAPRs = append(monthlyAPRs, monthlyAPR)
			}

			currentMonth = pointTime
			monthStartPrice = priceHistory[i].Price
		}
	}

	// Final month
	if len(priceHistory) > 1 && monthStartPrice > 0 {
		finalReturn := (priceHistory[len(priceHistory)-1].Price / monthStartPrice) - 1
		finalAPR := finalReturn * 12
		monthlyAPRs = append(monthlyAPRs, finalAPR)
	}

	return monthlyAPRs, nil
}

// calculateQuarterlyAPRs calculates quarterly APRs for stability analysis
func calculateQuarterlyAPRs(priceHistory []PricePoint) ([]float64, error) {
	if len(priceHistory) < 90 { // Need at least 3 months of data
		return []float64{}, fmt.Errorf("insufficient data for quarterly APR calculation")
	}

	// Sort price points by timestamp (oldest first)
	sort.Slice(priceHistory, func(i, j int) bool {
		return priceHistory[i].Timestamp < priceHistory[j].Timestamp
	})

	// Group prices into quarterly periods and calculate quarterly averages
	quarterlyAvgPrices := []float64{}

	// Start from the first available data point
	currentQuarter := getQuarter(time.Unix(priceHistory[0].Timestamp/1000, 0).UTC())
	quarterPrices := []float64{priceHistory[0].Price}

	for i := 1; i < len(priceHistory); i++ {
		pointTime := time.Unix(priceHistory[i].Timestamp/1000, 0).UTC()
		pointQuarter := getQuarter(pointTime)

		// Check if we've moved to a new quarter
		if pointQuarter != currentQuarter {
			// Calculate quarterly average price
			if len(quarterPrices) > 0 {
				sum := 0.0
				for _, price := range quarterPrices {
					sum += price
				}
				avgPrice := sum / float64(len(quarterPrices))
				quarterlyAvgPrices = append(quarterlyAvgPrices, avgPrice)
			}

			// Start new quarter
			currentQuarter = pointQuarter
			quarterPrices = []float64{priceHistory[i].Price}
		} else {
			// Add price to current quarter
			quarterPrices = append(quarterPrices, priceHistory[i].Price)
		}
	}

	// Calculate final quarter if we have data
	if len(quarterPrices) > 0 {
		sum := 0.0
		for _, price := range quarterPrices {
			sum += price
		}
		avgPrice := sum / float64(len(quarterPrices))
		quarterlyAvgPrices = append(quarterlyAvgPrices, avgPrice)
	}

	if len(quarterlyAvgPrices) < 2 {
		return []float64{}, fmt.Errorf("insufficient quarterly data")
	}

	// Calculate quarterly APRs
	quarterlyAPRs := []float64{}
	for i := 1; i < len(quarterlyAvgPrices); i++ {
		if quarterlyAvgPrices[i-1] > 0 {
			quarterlyReturn := (quarterlyAvgPrices[i] / quarterlyAvgPrices[i-1]) - 1
			quarterlyAPR := quarterlyReturn * 4 // Annualize the quarterly return
			quarterlyAPRs = append(quarterlyAPRs, quarterlyAPR)
		}
	}

	return quarterlyAPRs, nil
}
