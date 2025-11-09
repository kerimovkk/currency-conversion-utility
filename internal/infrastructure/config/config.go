package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	APIKey string
	APIURL string
}

// Load loads configuration from environment variables
// It attempts to load .env file if present, but doesn't fail if it's missing
func Load() (*Config, error) {
	// Try to load .env file
	_ = godotenv.Load()

	apiKey := os.Getenv("CMC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CMC_API_KEY environment variable is required")
	}

	apiURL := os.Getenv("CMC_API_URL")
	if apiURL == "" {
		// Default to sandbox URL if not specified
		apiURL = "https://sandbox-api.coinmarketcap.com"
	}

	return &Config{
		APIKey: apiKey,
		APIURL: apiURL,
	}, nil
}
