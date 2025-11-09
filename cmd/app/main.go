package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kerimovkk/currency-conversion-utility/internal/adapter/cli"
	"github.com/kerimovkk/currency-conversion-utility/internal/adapter/repository"
	"github.com/kerimovkk/currency-conversion-utility/internal/infrastructure/config"
	infrahttp "github.com/kerimovkk/currency-conversion-utility/internal/infrastructure/http"
	"github.com/kerimovkk/currency-conversion-utility/internal/usecase"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Parse command-line arguments
	args, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "\nRun 'app --help' for usage information")
		return 1
	}

	// Handle --help flag
	if args.ShowHelp {
		cli.ShowHelp()
		return 0
	}

	// Handle --version flag
	if args.ShowVersion {
		cli.ShowVersion()
		return 0
	}

	// Validate environment variables
	if err := cli.ValidateEnvironment(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "\nPlease set the CMC_API_KEY environment variable")
		fmt.Fprintln(os.Stderr, "You can copy .env.example to .env and add your API key")
		return 1
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		return 1
	}

	// Initialize dependencies (Dependency Injection)
	httpClient := infrahttp.NewClient()
	priceRepo := repository.NewCoinMarketCapRepository(httpClient, cfg.APIKey, cfg.APIURL)
	convertUseCase := usecase.NewConvertCurrencyUseCase(priceRepo)

	// Create presenter
	presenter := cli.NewPresenter(args.Verbose)

	// Execute conversion with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := convertUseCase.Execute(ctx, args.Amount, args.FromCurrency, args.ToCurrency)
	if err != nil {
		presenter.PresentError(err)
		return 1
	}

	// Present result
	presenter.PresentResult(result)
	return 0
}
