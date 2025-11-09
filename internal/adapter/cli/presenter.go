package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/kerimovkk/currency-conversion-utility/internal/domain"
)

// Presenter handles output formatting
type Presenter struct {
	verbose bool
}

// NewPresenter creates a new Presenter instance
func NewPresenter(verbose bool) *Presenter {
	return &Presenter{
		verbose: verbose,
	}
}

// PresentResult displays the conversion result
func (p *Presenter) PresentResult(result *domain.ConversionResult) {
	if p.verbose {
		p.presentVerbose(result)
	} else {
		p.presentSimple(result)
	}
}

// presentSimple displays a simple one-line result
func (p *Presenter) presentSimple(result *domain.ConversionResult) {
	fmt.Printf("%.8g %s = %.8g %s\n",
		result.OriginalAmount,
		result.FromCurrency.String(),
		result.ConvertedAmount,
		result.ToCurrency.String(),
	)
}

// presentVerbose displays detailed conversion information
func (p *Presenter) presentVerbose(result *domain.ConversionResult) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("CURRENCY CONVERSION RESULT")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Original Amount:    %.8g %s\n", result.OriginalAmount, result.FromCurrency.String())
	fmt.Printf("Converted Amount:   %.8g %s\n", result.ConvertedAmount, result.ToCurrency.String())
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Exchange Rate:      1 %s = %.8g %s\n",
		result.FromCurrency.String(),
		result.ExchangeRate,
		result.ToCurrency.String(),
	)
	fmt.Printf("Last Updated:       %s\n", result.LastUpdated.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Query Time:         %s\n", result.Timestamp.Format("2006-01-02 15:04:05 MST"))
	fmt.Println(strings.Repeat("=", 60))
}

// PresentError displays an error message in a user-friendly format
func (p *Presenter) PresentError(err error) {
	if p.verbose {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
