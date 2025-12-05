package db

import (
	"time"
)

// Token represents a token in the database
type Token struct {
	ID             int       `json:"id"`
	Symbol         string    `json:"symbol"`
	Name           string    `json:"name"`
	ContractAddress string   `json:"contract_address"`
	Decimals       int       `json:"decimals"`
	Blockchain     string    `json:"blockchain"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetAllTokens retrieves all active tokens from the database
func GetAllTokens() ([]Token, error) {
	query := `
		SELECT id, symbol, name, contract_address, decimals, blockchain, is_active, created_at, updated_at
		FROM tokens
		WHERE is_active = true
		ORDER BY symbol
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []Token
	for rows.Next() {
		var token Token
		err := rows.Scan(
			&token.ID,
			&token.Symbol,
			&token.Name,
			&token.ContractAddress,
			&token.Decimals,
			&token.Blockchain,
			&token.IsActive,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// GetTokenByID retrieves a token by its ID
func GetTokenByID(id int) (*Token, error) {
	query := `
		SELECT id, symbol, name, contract_address, decimals, blockchain, is_active, created_at, updated_at
		FROM tokens
		WHERE id = $1 AND is_active = true
	`

	var token Token
	err := DB.QueryRow(query, id).Scan(
		&token.ID,
		&token.Symbol,
		&token.Name,
		&token.ContractAddress,
		&token.Decimals,
		&token.Blockchain,
		&token.IsActive,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetTokenBySymbol retrieves a token by its symbol
func GetTokenBySymbol(symbol string) (*Token, error) {
	query := `
		SELECT id, symbol, name, contract_address, decimals, blockchain, is_active, created_at, updated_at
		FROM tokens
		WHERE symbol = $1 AND is_active = true
	`

	var token Token
	err := DB.QueryRow(query, symbol).Scan(
		&token.ID,
		&token.Symbol,
		&token.Name,
		&token.ContractAddress,
		&token.Decimals,
		&token.Blockchain,
		&token.IsActive,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}
