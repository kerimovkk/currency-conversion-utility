package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConversionRequest(t *testing.T) {
	btc, _ := NewCurrency("BTC")
	usd, _ := NewCurrency("USD")

	tests := []struct {
		name    string
		amount  float64
		from    *Currency
		to      *Currency
		wantErr bool
	}{
		{
			name:    "valid conversion request",
			amount:  100.0,
			from:    usd,
			to:      btc,
			wantErr: false,
		},
		{
			name:    "zero amount",
			amount:  0,
			from:    usd,
			to:      btc,
			wantErr: true,
		},
		{
			name:    "negative amount",
			amount:  -50.0,
			from:    usd,
			to:      btc,
			wantErr: true,
		},
		{
			name:    "nil from currency",
			amount:  100.0,
			from:    nil,
			to:      btc,
			wantErr: true,
		},
		{
			name:    "nil to currency",
			amount:  100.0,
			from:    usd,
			to:      nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConversionRequest(tt.amount, tt.from, tt.to)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.amount, got.Amount)
				assert.Equal(t, tt.from, got.FromCurrency)
				assert.Equal(t, tt.to, got.ToCurrency)
			}
		})
	}
}

func TestNewConversionResult(t *testing.T) {
	btc, _ := NewCurrency("BTC")
	usd, _ := NewCurrency("USD")
	now := time.Now()

	result := NewConversionResult(100.0, 0.0025, 0.000025, usd, btc, now, now)

	assert.NotNil(t, result)
	assert.Equal(t, 100.0, result.OriginalAmount)
	assert.Equal(t, 0.0025, result.ConvertedAmount)
	assert.Equal(t, 0.000025, result.ExchangeRate)
	assert.Equal(t, usd, result.FromCurrency)
	assert.Equal(t, btc, result.ToCurrency)
	assert.Equal(t, now, result.Timestamp)
	assert.Equal(t, now, result.LastUpdated)
}
