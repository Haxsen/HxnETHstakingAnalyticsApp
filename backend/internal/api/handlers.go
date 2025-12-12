package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	tokenService     *services.TokenService
	valuationService *services.ValuationService
}

// NewHandler creates a new handler with dependencies
func NewHandler(tokenService *services.TokenService, valuationService *services.ValuationService) *Handler {
	return &Handler{
		tokenService:     tokenService,
		valuationService: valuationService,
	}
}

// GetTokensHandler returns all active tokens
//
// @Summary Get all tracked LST tokens
// @Description Retrieve a list of all active Liquid Staking Tokens being tracked
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "tokens: array of token objects, count: number of tokens"
// @Failure 500 {object} map[string]string "error: error message"
// @Router /api/tokens [get]
func (h *Handler) GetTokensHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokens, err := h.tokenService.GetAllTokens(r.Context())
	if err != nil {
		log.Printf("Error fetching tokens: %v", err)
		JSONError(w, "Failed to fetch tokens", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, map[string]interface{}{
		"tokens": tokens,
		"count":  len(tokens),
	})
}

// GetTokenHistoryHandler returns price history for a token
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
func (h *Handler) GetTokenHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenSymbol := chi.URLParam(r, "id") // Using 'id' param but it's actually the symbol

	// Validate that the token exists
	if err := h.tokenService.ValidateTokenExists(r.Context(), tokenSymbol); err != nil {
		JSONError(w, "Token not found or not supported", http.StatusBadRequest)
		return
	}

	// Fetch price history
	priceHistory, err := h.valuationService.GetTokenHistory(r.Context(), tokenSymbol)
	if err != nil {
		log.Printf("Error fetching price history for %s: %v", tokenSymbol, err)
		JSONError(w, "Failed to fetch price history", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, map[string]interface{}{
		"token_symbol": tokenSymbol,
		"price_history": priceHistory,
		"count": len(priceHistory),
	})
}

// GetTokenValuationHandler returns valuation metrics for a specific token
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
func (h *Handler) GetTokenValuationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenSymbol := chi.URLParam(r, "id")

	// Get token details
	token, err := h.tokenService.GetTokenBySymbol(r.Context(), tokenSymbol)
	if err != nil {
		JSONError(w, "Token not found or not supported", http.StatusBadRequest)
		return
	}

	// Get valuation using service
	valuation, err := h.valuationService.GetTokenValuation(r.Context(), tokenSymbol, token)
	if err != nil {
		log.Printf("Error getting valuation for %s: %v", tokenSymbol, err)
		JSONError(w, "Failed to calculate valuation", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(valuation); err != nil {
		log.Printf("Error encoding response: %v", err)
		JSONError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAllValuationsHandler returns valuation metrics for all tokens
//
// @Summary Get valuation metrics for all tokens
// @Description Retrieve APR, stability, TVL, and valuation remarks for all tracked LST tokens (sortable table data)
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "valuations: array of valuation objects, count: number of valuations"
// @Failure 500 {object} map[string]string "error: failed to fetch valuations"
// @Router /api/valuations [get]
func (h *Handler) GetAllValuationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get all tokens
	tokens, err := h.tokenService.GetAllTokens(r.Context())
	if err != nil {
		log.Printf("Error fetching tokens: %v", err)
		JSONError(w, "Failed to fetch tokens", http.StatusInternalServerError)
		return
	}

	// Get valuations for all tokens
	valuations, err := h.valuationService.GetAllTokenValuations(r.Context(), tokens)
	if err != nil {
		log.Printf("Error getting valuations: %v", err)
		JSONError(w, "Failed to fetch valuations", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, map[string]interface{}{
		"valuations": valuations,
		"count":      len(valuations),
	})
}

func (h *Handler) RefreshCacheHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Cache refresh - TODO"}`))
}
