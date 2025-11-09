package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	version = "1.0.0"
)

// Args represents parsed command-line arguments
type Args struct {
	Amount      float64
	FromCurrency string
	ToCurrency   string
	Verbose     bool
	ShowHelp    bool
	ShowVersion bool
}

// ParseArgs parses command-line arguments
func ParseArgs(args []string) (*Args, error) {
	// Create a new flag set for custom parsing
	fs := flag.NewFlagSet("currency-converter", flag.ContinueOnError)

	// Define flags
	verbose := fs.Bool("verbose", false, "Enable verbose output")
	help := fs.Bool("help", false, "Show help message")
	version := fs.Bool("version", false, "Show version information")

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	result := &Args{
		Verbose:     *verbose,
		ShowHelp:    *help,
		ShowVersion: *version,
	}

	// If help or version requested, return early
	if result.ShowHelp || result.ShowVersion {
		return result, nil
	}

	// Get remaining arguments (non-flag arguments)
	remaining := fs.Args()

	// Validate argument count
	if len(remaining) != 3 {
		return nil, fmt.Errorf("invalid number of arguments: expected 3 (amount, from_currency, to_currency), got %d", len(remaining))
	}

	// Parse amount
	amount, err := strconv.ParseFloat(remaining[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount '%s': must be a number", remaining[0])
	}

	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount '%f': must be greater than zero", amount)
	}

	result.Amount = amount
	result.FromCurrency = remaining[1]
	result.ToCurrency = remaining[2]

	return result, nil
}

// ShowHelp displays help message
func ShowHelp() {
	fmt.Println("Currency Conversion Utility")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  app [options] <amount> <from_currency> <to_currency>")
	fmt.Println()
	fmt.Println("ARGUMENTS:")
	fmt.Println("  amount          Amount to convert (must be > 0)")
	fmt.Println("  from_currency   Source currency symbol (e.g., USD, BTC)")
	fmt.Println("  to_currency     Target currency symbol (e.g., EUR, ETH)")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --help          Show this help message")
	fmt.Println("  --version       Show version information")
	fmt.Println("  --verbose       Enable verbose output with detailed information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  app 123.45 USD BTC")
	fmt.Println("  app --verbose 100 BTC USD")
	fmt.Println("  app 50 EUR GBP")
	fmt.Println()
	fmt.Println("ENVIRONMENT VARIABLES:")
	fmt.Println("  CMC_API_KEY     CoinMarketCap API key (required)")
	fmt.Println("  CMC_API_URL     CoinMarketCap API base URL (optional)")
	fmt.Println()
}

// ShowVersion displays version information
func ShowVersion() {
	fmt.Printf("Currency Conversion Utility v%s\n", version)
	fmt.Println("Powered by CoinMarketCap API")
}

// ValidateEnvironment checks if required environment variables are set
func ValidateEnvironment() error {
	if os.Getenv("CMC_API_KEY") == "" {
		return fmt.Errorf("CMC_API_KEY environment variable is not set")
	}
	return nil
}
