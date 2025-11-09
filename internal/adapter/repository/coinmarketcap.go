package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kerimovkk/currency-conversion-utility/internal/domain"
	"github.com/kerimovkk/currency-conversion-utility/pkg/retry"
)

// CoinMarketCapRepository implements domain.PriceRepository for CoinMarketCap API
type CoinMarketCapRepository struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	retry      *retry.Strategy
}

// NewCoinMarketCapRepository creates a new CoinMarketCap API client
func NewCoinMarketCapRepository(httpClient *http.Client, apiKey, baseURL string) *CoinMarketCapRepository {
	return &CoinMarketCapRepository{
		httpClient: httpClient,
		apiKey:     apiKey,
		baseURL:    baseURL,
		retry:      retry.DefaultStrategy(),
	}
}

// APIResponse represents the structure of the CoinMarketCap API response
type APIResponse struct {
	Status StatusObject `json:"status"`
	Data   interface{}  `json:"data"`
}

// StatusObject represents the status object in API responses
type StatusObject struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage *string   `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
}

// PriceConversionData represents the data returned by /v1/tools/price-conversion
type PriceConversionData struct {
	Symbol      string                 `json:"symbol"`
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Amount      float64                `json:"amount"`
	LastUpdated time.Time              `json:"last_updated"`
	Quote       map[string]QuoteDetail `json:"quote"`
}

// QuoteDetail contains the conversion details for a specific currency
type QuoteDetail struct {
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}

// GetConversionPrice fetches the conversion price from CoinMarketCap API
func (c *CoinMarketCapRepository) GetConversionPrice(
	ctx context.Context,
	amount float64,
	from, to string,
) (*domain.ConversionResult, error) {
	var result *domain.ConversionResult
	var lastErr error

	// Execute with retry logic
	err := retry.Do(ctx, c.retry, c.shouldRetry, func(ctx context.Context) error {
		var err error
		result, err = c.fetchConversionPrice(ctx, amount, from, to)
		lastErr = err
		return err
	})

	if err != nil {
		return nil, lastErr
	}

	return result, nil
}

// fetchConversionPrice performs the actual API call
func (c *CoinMarketCapRepository) fetchConversionPrice(
	ctx context.Context,
	amount float64,
	from, to string,
) (*domain.ConversionResult, error) {
	// Build request URL
	endpoint := fmt.Sprintf("%s/v1/tools/price-conversion", c.baseURL)

	params := url.Values{}
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	params.Add("symbol", from)
	params.Add("convert", to)

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrAPIFailure, err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNetworkFailure, err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response body", domain.ErrAPIFailure)
	}

	// Handle HTTP error status codes
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleHTTPError(resp.StatusCode, body)
	}

	// Parse response
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidResponse, err)
	}

	// Check for API-level errors
	if apiResp.Status.ErrorCode != 0 {
		return nil, c.handleAPIError(apiResp.Status.ErrorCode, apiResp.Status.ErrorMessage)
	}

	return c.parseConversionData(apiResp.Data, amount, from, to)
}

// handleHTTPError converts HTTP error codes to domain errors
func (c *CoinMarketCapRepository) handleHTTPError(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", domain.ErrInvalidCurrency, string(body))
	case http.StatusUnauthorized:
		return domain.ErrUnauthorized
	case http.StatusForbidden:
		return domain.ErrForbidden
	case http.StatusTooManyRequests:
		return domain.ErrRateLimitExceeded
	case http.StatusInternalServerError, http.StatusBadGateway,
		 http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return domain.ErrServerError
	default:
		return fmt.Errorf("%w: HTTP %d", domain.ErrAPIFailure, statusCode)
	}
}

// handleAPIError converts API error codes to domain errors
func (c *CoinMarketCapRepository) handleAPIError(errorCode int, errorMsg *string) error {
	msg := ""
	if errorMsg != nil {
		msg = *errorMsg
	}

	switch errorCode {
	case 1001, 1002:
		return domain.ErrUnauthorized
	case 1005, 1006, 1007:
		return domain.ErrForbidden
	case 1008, 1009, 1010, 1011:
		return domain.ErrRateLimitExceeded
	default:
		return fmt.Errorf("%w: code %d - %s", domain.ErrAPIFailure, errorCode, msg)
	}
}

// shouldRetry determines if an error should trigger a retry
func (c *CoinMarketCapRepository) shouldRetry(err error) bool {
	// Retry on rate limit and server errors
	return err == domain.ErrRateLimitExceeded ||
		   err == domain.ErrServerError ||
		   err == domain.ErrNetworkFailure
}

// parseConversionData extracts conversion result from API response data
func (c *CoinMarketCapRepository) parseConversionData(
	data interface{},
	amount float64,
	from, to string,
) (*domain.ConversionResult, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal data", domain.ErrInvalidResponse)
	}

	var convData PriceConversionData
	if err := json.Unmarshal(dataBytes, &convData); err != nil {
		return nil, fmt.Errorf("%w: failed to parse conversion data", domain.ErrInvalidResponse)
	}

	// Extract quote for target currency
	quoteDetail, ok := convData.Quote[to]
	if !ok {
		return nil, fmt.Errorf("%w: no quote found for %s", domain.ErrInvalidResponse, to)
	}

	// Create currency objects
	fromCurrency, err := domain.NewCurrency(from)
	if err != nil {
		return nil, err
	}

	toCurrency, err := domain.NewCurrency(to)
	if err != nil {
		return nil, err
	}

	// Calculate converted amount and exchange rate
	convertedAmount := quoteDetail.Price
	exchangeRate := convertedAmount / amount

	// Create and return result
	return domain.NewConversionResult(
		amount,
		convertedAmount,
		exchangeRate,
		fromCurrency,
		toCurrency,
		time.Now(),
		quoteDetail.LastUpdated,
	), nil
}
