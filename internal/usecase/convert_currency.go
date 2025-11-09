package usecase

import (
	"context"

	"github.com/kerimovkk/currency-conversion-utility/internal/domain"
)

// ConvertCurrencyUseCase handles currency conversion business logic
type ConvertCurrencyUseCase struct {
	priceRepo domain.PriceRepository
}

// NewConvertCurrencyUseCase creates a new ConvertCurrencyUseCase instance
func NewConvertCurrencyUseCase(priceRepo domain.PriceRepository) *ConvertCurrencyUseCase {
	return &ConvertCurrencyUseCase{
		priceRepo: priceRepo,
	}
}

// Execute performs currency conversion
func (uc *ConvertCurrencyUseCase) Execute(
	ctx context.Context,
	amount float64,
	fromSymbol, toSymbol string,
) (*domain.ConversionResult, error) {
	// Validate and create currencies
	fromCurrency, err := domain.NewCurrency(fromSymbol)
	if err != nil {
		return nil, err
	}

	toCurrency, err := domain.NewCurrency(toSymbol)
	if err != nil {
		return nil, err
	}

	// Create conversion request
	request, err := domain.NewConversionRequest(amount, fromCurrency, toCurrency)
	if err != nil {
		return nil, err
	}

	// Fetch conversion from repository
	result, err := uc.priceRepo.GetConversionPrice(
		ctx,
		request.Amount,
		request.FromCurrency.String(),
		request.ToCurrency.String(),
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
