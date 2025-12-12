package services

import (
	"context"
	"fmt"

	"github.com/Haxsen/HxnETHstakingAnalyticsApp/backend/internal/db"
)

// TokenService handles token-related business logic
type TokenService struct {
}

// NewTokenService creates a new token service
func NewTokenService() *TokenService {
	return &TokenService{}
}

// Token represents a token entity
type Token struct {
	ID             int    `json:"id"`
	Symbol         string `json:"symbol"`
	Name           string `json:"name"`
	ContractAddress string `json:"contract_address"`
	Decimals       int    `json:"decimals"`
	Blockchain     string `json:"blockchain"`
	IsActive       bool   `json:"is_active"`
}

// GetAllTokens retrieves all active tokens
func (s *TokenService) GetAllTokens(ctx context.Context) ([]Token, error) {
	dbTokens, err := db.GetAllTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens from database: %w", err)
	}

	// Convert db models to service models
	tokens := make([]Token, len(dbTokens))
	for i, dbToken := range dbTokens {
		tokens[i] = Token{
			ID:             dbToken.ID,
			Symbol:         dbToken.Symbol,
			Name:           dbToken.Name,
			ContractAddress: dbToken.ContractAddress,
			Decimals:       dbToken.Decimals,
			Blockchain:     dbToken.Blockchain,
			IsActive:       dbToken.IsActive,
		}
	}

	return tokens, nil
}

// GetTokenBySymbol retrieves a token by its symbol
func (s *TokenService) GetTokenBySymbol(ctx context.Context, symbol string) (*Token, error) {
	dbToken, err := db.GetTokenBySymbol(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get token by symbol: %w", err)
	}

	token := &Token{
		ID:             dbToken.ID,
		Symbol:         dbToken.Symbol,
		Name:           dbToken.Name,
		ContractAddress: dbToken.ContractAddress,
		Decimals:       dbToken.Decimals,
		Blockchain:     dbToken.Blockchain,
		IsActive:       dbToken.IsActive,
	}

	return token, nil
}

// ValidateTokenExists checks if a token exists and is active
func (s *TokenService) ValidateTokenExists(ctx context.Context, symbol string) error {
	token, err := s.GetTokenBySymbol(ctx, symbol)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	if !token.IsActive {
		return fmt.Errorf("token %s is not active", symbol)
	}

	return nil
}
