package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Strategy defines retry behavior configuration
type Strategy struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultStrategy returns a sensible default retry strategy
// Retries: 4 attempts total (1 original + 3 retries)
// Delays: 1s, 2s, 4s (exponential backoff with multiplier 2.0)
func DefaultStrategy() *Strategy {
	return &Strategy{
		MaxAttempts:  4,
		InitialDelay: 1 * time.Second,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// ShouldRetryFunc determines if an error should trigger a retry
type ShouldRetryFunc func(err error) bool

// Do executes the given function with retry logic
func Do(ctx context.Context, strategy *Strategy, shouldRetry ShouldRetryFunc, fn RetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt < strategy.MaxAttempts; attempt++ {
		// Execute the function
		err := fn(ctx)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if we should retry
		if !shouldRetry(err) {
			return err // Don't retry, return error immediately
		}

		// Don't sleep after the last attempt
		if attempt == strategy.MaxAttempts-1 {
			break
		}

		// Calculate delay with exponential backoff
		delay := strategy.calculateDelay(attempt)

		// Check if context is cancelled before sleeping
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("max retry attempts (%d) exceeded: %w", strategy.MaxAttempts, lastErr)
}

// calculateDelay calculates the delay for a given attempt using exponential backoff
func (s *Strategy) calculateDelay(attempt int) time.Duration {
	delay := float64(s.InitialDelay) * math.Pow(s.Multiplier, float64(attempt))
	delayDuration := time.Duration(delay)

	if delayDuration > s.MaxDelay {
		delayDuration = s.MaxDelay
	}

	return delayDuration
}
