package domain

import "context"

// CurrencyConverter defines the interface for currency conversion operations
type CurrencyConverter interface {
	// Convert performs currency conversion from one currency to another
	Convert(ctx context.Context, request *ConversionRequest) (*ConversionResult, error)
}

// PriceRepository defines the interface for fetching price data from external sources
type PriceRepository interface {
	// GetConversionPrice fetches the conversion price from the external API
	GetConversionPrice(ctx context.Context, amount float64, from, to string) (*ConversionResult, error)
}
