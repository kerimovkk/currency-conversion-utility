package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/kerimovkk/currency-conversion-utility/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPriceRepository is a mock implementation of domain.PriceRepository
type MockPriceRepository struct {
	mock.Mock
}

func (m *MockPriceRepository) GetConversionPrice(ctx context.Context, amount float64, from, to string) (*domain.ConversionResult, error) {
	args := m.Called(ctx, amount, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ConversionResult), args.Error(1)
}

func TestConvertCurrencyUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		fromSymbol  string
		toSymbol    string
		mockSetup   func(*MockPriceRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name:       "successful conversion",
			amount:     100.0,
			fromSymbol: "USD",
			toSymbol:   "BTC",
			mockSetup: func(m *MockPriceRepository) {
				from, _ := domain.NewCurrency("USD")
				to, _ := domain.NewCurrency("BTC")
				result := domain.NewConversionResult(
					100.0,
					0.0025,
					0.000025,
					from,
					to,
					time.Now(),
					time.Now(),
				)
				m.On("GetConversionPrice", mock.Anything, 100.0, "USD", "BTC").Return(result, nil)
			},
			wantErr: false,
		},
		{
			name:       "invalid amount (zero)",
			amount:     0,
			fromSymbol: "USD",
			toSymbol:   "BTC",
			mockSetup: func(m *MockPriceRepository) {
				// Mock won't be called because validation happens first
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidAmount,
		},
		{
			name:       "invalid amount (negative)",
			amount:     -50.0,
			fromSymbol: "USD",
			toSymbol:   "BTC",
			mockSetup: func(m *MockPriceRepository) {
				// Mock won't be called because validation happens first
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidAmount,
		},
		{
			name:       "invalid from currency",
			amount:     100.0,
			fromSymbol: "",
			toSymbol:   "BTC",
			mockSetup: func(m *MockPriceRepository) {
				// Mock won't be called
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidCurrency,
		},
		{
			name:       "invalid to currency",
			amount:     100.0,
			fromSymbol: "USD",
			toSymbol:   "X",
			mockSetup: func(m *MockPriceRepository) {
				// Mock won't be called
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidCurrency,
		},
		{
			name:       "repository returns error",
			amount:     100.0,
			fromSymbol: "USD",
			toSymbol:   "BTC",
			mockSetup: func(m *MockPriceRepository) {
				m.On("GetConversionPrice", mock.Anything, 100.0, "USD", "BTC").
					Return(nil, domain.ErrRateLimitExceeded)
			},
			wantErr:     true,
			expectedErr: domain.ErrRateLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPriceRepository)
			tt.mockSetup(mockRepo)

			uc := NewConvertCurrencyUseCase(mockRepo)
			ctx := context.Background()

			result, err := uc.Execute(ctx, tt.amount, tt.fromSymbol, tt.toSymbol)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.amount, result.OriginalAmount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
