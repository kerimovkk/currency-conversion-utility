package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCurrency(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		want    string
		wantErr bool
	}{
		{
			name:    "valid uppercase symbol",
			symbol:  "BTC",
			want:    "BTC",
			wantErr: false,
		},
		{
			name:    "valid lowercase symbol (should be converted to uppercase)",
			symbol:  "eth",
			want:    "ETH",
			wantErr: false,
		},
		{
			name:    "valid symbol with spaces (should be trimmed)",
			symbol:  "  USD  ",
			want:    "USD",
			wantErr: false,
		},
		{
			name:    "empty symbol",
			symbol:  "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "too short symbol",
			symbol:  "B",
			want:    "",
			wantErr: true,
		},
		{
			name:    "too long symbol",
			symbol:  "VERYLONGCURRENCYSYMBOL",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCurrency(tt.symbol)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got.Symbol)
			}
		})
	}
}

func TestCurrency_String(t *testing.T) {
	currency, _ := NewCurrency("BTC")
	assert.Equal(t, "BTC", currency.String())
}

func TestCurrency_Equals(t *testing.T) {
	btc1, _ := NewCurrency("BTC")
	btc2, _ := NewCurrency("BTC")
	eth, _ := NewCurrency("ETH")

	assert.True(t, btc1.Equals(btc2))
	assert.False(t, btc1.Equals(eth))
	assert.False(t, btc1.Equals(nil))
}
