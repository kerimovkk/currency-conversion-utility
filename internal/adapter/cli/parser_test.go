package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Args
		wantErr bool
	}{
		{
			name: "valid args with amount and currencies",
			args: []string{"123.45", "USD", "BTC"},
			want: &Args{
				Amount:       123.45,
				FromCurrency: "USD",
				ToCurrency:   "BTC",
				Verbose:      false,
				ShowHelp:     false,
				ShowVersion:  false,
			},
			wantErr: false,
		},
		{
			name: "valid args with verbose flag",
			args: []string{"--verbose", "100", "BTC", "ETH"},
			want: &Args{
				Amount:       100,
				FromCurrency: "BTC",
				ToCurrency:   "ETH",
				Verbose:      true,
				ShowHelp:     false,
				ShowVersion:  false,
			},
			wantErr: false,
		},
		{
			name: "help flag",
			args: []string{"--help"},
			want: &Args{
				ShowHelp:    true,
				Verbose:     false,
				ShowVersion: false,
			},
			wantErr: false,
		},
		{
			name: "version flag",
			args: []string{"--version"},
			want: &Args{
				ShowVersion: true,
				ShowHelp:    false,
				Verbose:     false,
			},
			wantErr: false,
		},
		{
			name:    "too few arguments",
			args:    []string{"100", "USD"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    []string{"100", "USD", "BTC", "extra"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid amount (not a number)",
			args:    []string{"abc", "USD", "BTC"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid amount (zero)",
			args:    []string{"0", "USD", "BTC"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid amount (negative)",
			args:    []string{"-50", "USD", "BTC"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Amount, got.Amount)
				assert.Equal(t, tt.want.FromCurrency, got.FromCurrency)
				assert.Equal(t, tt.want.ToCurrency, got.ToCurrency)
				assert.Equal(t, tt.want.Verbose, got.Verbose)
				assert.Equal(t, tt.want.ShowHelp, got.ShowHelp)
				assert.Equal(t, tt.want.ShowVersion, got.ShowVersion)
			}
		})
	}
}
