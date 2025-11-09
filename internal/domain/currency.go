package domain

import (
	"strings"
)

// Currency represents a currency symbol (e.g., BTC, USD, EUR)
type Currency struct {
	Symbol string
}

// NewCurrency creates a new Currency with validation
func NewCurrency(symbol string) (*Currency, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	if symbol == "" {
		return nil, ErrInvalidCurrency
	}

	// Currency symbols are typically 3-4 characters (ISO codes or crypto symbols)
	if len(symbol) < 2 || len(symbol) > 10 {
		return nil, ErrInvalidCurrency
	}

	return &Currency{Symbol: symbol}, nil
}

// String returns the currency symbol as a string
func (c *Currency) String() string {
	return c.Symbol
}

// Equals checks if two currencies are equal
func (c *Currency) Equals(other *Currency) bool {
	if other == nil {
		return false
	}
	return c.Symbol == other.Symbol
}
