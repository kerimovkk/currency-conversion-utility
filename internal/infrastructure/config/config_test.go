package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		wantErr     bool
		expectedURL string
	}{
		{
			name: "valid config with all env vars",
			setupEnv: func() {
				os.Setenv("CMC_API_KEY", "test-api-key")
				os.Setenv("CMC_API_URL", "https://test-api.example.com")
			},
			cleanupEnv: func() {
				os.Unsetenv("CMC_API_KEY")
				os.Unsetenv("CMC_API_URL")
			},
			wantErr:     false,
			expectedURL: "https://test-api.example.com",
		},
		{
			name: "valid config with default URL",
			setupEnv: func() {
				os.Setenv("CMC_API_KEY", "test-api-key")
			},
			cleanupEnv: func() {
				os.Unsetenv("CMC_API_KEY")
			},
			wantErr:     false,
			expectedURL: "https://sandbox-api.coinmarketcap.com",
		},
		{
			name: "missing API key",
			setupEnv: func() {
				os.Unsetenv("CMC_API_KEY")
			},
			cleanupEnv: func() {},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanupEnv()

			cfg, err := Load()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.NotEmpty(t, cfg.APIKey)
				assert.Equal(t, tt.expectedURL, cfg.APIURL)
			}
		})
	}
}
