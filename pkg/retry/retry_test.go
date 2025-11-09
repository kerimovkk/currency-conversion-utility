package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo_SuccessFirstAttempt(t *testing.T) {
	strategy := DefaultStrategy()
	callCount := 0

	err := Do(context.Background(), strategy, func(err error) bool {
		return true
	}, func(ctx context.Context) error {
		callCount++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestDo_SuccessAfterRetries(t *testing.T) {
	strategy := &Strategy{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	callCount := 0
	testErr := errors.New("test error")

	err := Do(context.Background(), strategy, func(err error) bool {
		return err == testErr
	}, func(ctx context.Context) error {
		callCount++
		if callCount < 3 {
			return testErr
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, callCount)
}

func TestDo_MaxRetriesExceeded(t *testing.T) {
	strategy := &Strategy{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	callCount := 0
	testErr := errors.New("persistent error")

	err := Do(context.Background(), strategy, func(err error) bool {
		return err == testErr
	}, func(ctx context.Context) error {
		callCount++
		return testErr
	})

	assert.Error(t, err)
	assert.Equal(t, 3, callCount)
	assert.ErrorIs(t, err, testErr)
}

func TestDo_NoRetryOnNonRetryableError(t *testing.T) {
	strategy := DefaultStrategy()
	callCount := 0
	retryableErr := errors.New("retryable")
	nonRetryableErr := errors.New("non-retryable")

	err := Do(context.Background(), strategy, func(err error) bool {
		return err == retryableErr
	}, func(ctx context.Context) error {
		callCount++
		return nonRetryableErr
	})

	assert.Error(t, err)
	assert.Equal(t, 1, callCount)
	assert.Equal(t, nonRetryableErr, err)
}

func TestDo_ContextCancellation(t *testing.T) {
	strategy := &Strategy{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
	}

	ctx, cancel := context.WithCancel(context.Background())
	callCount := 0
	testErr := errors.New("test error")

	// Cancel context after first failure
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := Do(ctx, strategy, func(err error) bool {
		return err == testErr
	}, func(ctx context.Context) error {
		callCount++
		return testErr
	})

	assert.Error(t, err)
	assert.Less(t, callCount, 5) // Should not reach max attempts
}

func TestStrategy_calculateDelay(t *testing.T) {
	strategy := &Strategy{
		InitialDelay: 1 * time.Second,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
		{4, 10 * time.Second}, // Capped at MaxDelay
		{5, 10 * time.Second}, // Capped at MaxDelay
	}

	for _, tt := range tests {
		got := strategy.calculateDelay(tt.attempt)
		assert.Equal(t, tt.want, got, "attempt %d", tt.attempt)
	}
}
