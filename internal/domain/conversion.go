package domain

import "time"

// ConversionRequest represents a request to convert one currency to another
type ConversionRequest struct {
	Amount       float64
	FromCurrency *Currency
	ToCurrency   *Currency
}

// NewConversionRequest creates a new ConversionRequest with validation
func NewConversionRequest(amount float64, from, to *Currency) (*ConversionRequest, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if from == nil || to == nil {
		return nil, ErrInvalidCurrency
	}

	return &ConversionRequest{
		Amount:       amount,
		FromCurrency: from,
		ToCurrency:   to,
	}, nil
}

// ConversionResult represents the result of a currency conversion
type ConversionResult struct {
	OriginalAmount   float64
	ConvertedAmount  float64
	FromCurrency     *Currency
	ToCurrency       *Currency
	ExchangeRate     float64
	Timestamp        time.Time
	LastUpdated      time.Time
}

// NewConversionResult creates a new ConversionResult
func NewConversionResult(
	originalAmount, convertedAmount, exchangeRate float64,
	from, to *Currency,
	timestamp, lastUpdated time.Time,
) *ConversionResult {
	return &ConversionResult{
		OriginalAmount:  originalAmount,
		ConvertedAmount: convertedAmount,
		FromCurrency:    from,
		ToCurrency:      to,
		ExchangeRate:    exchangeRate,
		Timestamp:       timestamp,
		LastUpdated:     lastUpdated,
	}
}
