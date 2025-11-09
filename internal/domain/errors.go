package domain

import "errors"

var (
	// ErrInvalidCurrency indicates that the provided currency symbol is invalid
	ErrInvalidCurrency = errors.New("invalid currency symbol")

	// ErrInvalidAmount indicates that the provided amount is invalid (e.g., negative or zero)
	ErrInvalidAmount = errors.New("invalid amount: must be greater than zero")

	// ErrAPIKeyMissing indicates that the API key is not configured
	ErrAPIKeyMissing = errors.New("API key is missing")

	// ErrAPIFailure indicates a general API failure
	ErrAPIFailure = errors.New("API request failed")

	// ErrUnauthorized indicates authentication failure (401)
	ErrUnauthorized = errors.New("unauthorized: invalid API key")

	// ErrForbidden indicates insufficient permissions (403)
	ErrForbidden = errors.New("forbidden: API key does not have access to this endpoint")

	// ErrRateLimitExceeded indicates rate limit was exceeded (429)
	ErrRateLimitExceeded = errors.New("rate limit exceeded: too many requests")

	// ErrServerError indicates server-side error (5xx)
	ErrServerError = errors.New("server error: please try again later")

	// ErrNetworkFailure indicates network connectivity issues
	ErrNetworkFailure = errors.New("network failure: could not connect to API")

	// ErrInvalidResponse indicates the API response could not be parsed
	ErrInvalidResponse = errors.New("invalid API response")
)
