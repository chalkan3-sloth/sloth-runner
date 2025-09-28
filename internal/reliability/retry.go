package reliability

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryStrategy defines different retry strategies
type RetryStrategy int

const (
	FixedDelay RetryStrategy = iota
	ExponentialBackoff
	LinearBackoff
	CustomBackoff
)

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxAttempts   int                                    // Maximum number of attempts
	InitialDelay  time.Duration                          // Initial delay between retries
	MaxDelay      time.Duration                          // Maximum delay between retries
	Strategy      RetryStrategy                          // Retry strategy to use
	Multiplier    float64                                // Multiplier for exponential/linear backoff
	Jitter        bool                                   // Add random jitter to delays
	ShouldRetry   func(error) bool                      // Function to determine if error is retryable
	OnRetry       func(attempt int, delay time.Duration, err error) // Callback on retry
	CustomDelayFn func(attempt int) time.Duration       // Custom delay function
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Strategy:     ExponentialBackoff,
		Multiplier:   2.0,
		Jitter:       true,
		ShouldRetry: func(err error) bool {
			// By default, retry on all errors except circuit breaker errors
			return !IsCircuitBreakerError(err)
		},
	}
}

// RetryStats holds statistics about retry operations
type RetryStats struct {
	TotalAttempts    int           `json:"total_attempts"`
	SuccessfulRetries int          `json:"successful_retries"`
	FailedRetries    int           `json:"failed_retries"`
	TotalDelay       time.Duration `json:"total_delay"`
	AverageDelay     time.Duration `json:"average_delay"`
	LastAttemptTime  time.Time     `json:"last_attempt_time"`
}

// Retrier handles retry logic with various strategies
type Retrier struct {
	config RetryConfig
	stats  RetryStats
}

// NewRetrier creates a new retrier with the given configuration
func NewRetrier(config RetryConfig) *Retrier {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	if config.InitialDelay <= 0 {
		config.InitialDelay = 1 * time.Second
	}
	if config.MaxDelay <= 0 {
		config.MaxDelay = 30 * time.Second
	}
	if config.Multiplier <= 0 {
		config.Multiplier = 2.0
	}
	if config.ShouldRetry == nil {
		config.ShouldRetry = func(err error) bool { return true }
	}

	return &Retrier{
		config: config,
	}
}

// Execute executes a function with retry logic
func (r *Retrier) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return r.ExecuteWithContext(context.Background(), func(ctx context.Context) (interface{}, error) {
		return fn()
	})
}

// ExecuteWithContext executes a function with retry logic and context support
func (r *Retrier) ExecuteWithContext(ctx context.Context, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	var lastErr error
	startTime := time.Now()

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		r.stats.TotalAttempts++
		r.stats.LastAttemptTime = time.Now()

		// Execute the function
		result, err := fn(ctx)
		if err == nil {
			if attempt > 1 {
				r.stats.SuccessfulRetries++
			}
			return result, nil
		}

		lastErr = err

		// Check if we should retry
		if !r.config.ShouldRetry(err) {
			r.stats.FailedRetries++
			return result, err
		}

		// Don't retry if this is the last attempt
		if attempt >= r.config.MaxAttempts {
			r.stats.FailedRetries++
			break
		}

		// Calculate delay
		delay := r.calculateDelay(attempt)
		r.stats.TotalDelay += delay

		// Call retry callback if provided
		if r.config.OnRetry != nil {
			r.config.OnRetry(attempt, delay, err)
		}

		// Wait for the delay or context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	// Update average delay
	if r.stats.TotalAttempts > 0 {
		r.stats.AverageDelay = time.Duration(int64(r.stats.TotalDelay) / int64(r.stats.TotalAttempts))
	}

	return nil, &RetryExhaustedError{
		LastError:    lastErr,
		Attempts:     r.config.MaxAttempts,
		TotalTime:    time.Since(startTime),
	}
}

// calculateDelay calculates the delay for the given attempt based on the strategy
func (r *Retrier) calculateDelay(attempt int) time.Duration {
	var delay time.Duration

	switch r.config.Strategy {
	case FixedDelay:
		delay = r.config.InitialDelay

	case ExponentialBackoff:
		delay = time.Duration(float64(r.config.InitialDelay) * math.Pow(r.config.Multiplier, float64(attempt-1)))

	case LinearBackoff:
		delay = time.Duration(float64(r.config.InitialDelay) * float64(attempt) * r.config.Multiplier)

	case CustomBackoff:
		if r.config.CustomDelayFn != nil {
			delay = r.config.CustomDelayFn(attempt)
		} else {
			delay = r.config.InitialDelay
		}

	default:
		delay = r.config.InitialDelay
	}

	// Apply maximum delay constraint
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}

	// Add jitter if enabled
	if r.config.Jitter {
		jitter := time.Duration(rand.Float64() * float64(delay) * 0.1) // Up to 10% jitter
		if rand.Float64() < 0.5 {
			delay -= jitter
		} else {
			delay += jitter
		}
	}

	// Ensure minimum delay
	if delay < 0 {
		delay = r.config.InitialDelay
	}

	return delay
}

// GetStats returns the current retry statistics
func (r *Retrier) GetStats() RetryStats {
	return r.stats
}

// Reset resets the retry statistics
func (r *Retrier) Reset() {
	r.stats = RetryStats{}
}

// RetryExhaustedError represents an error when all retry attempts are exhausted
type RetryExhaustedError struct {
	LastError error
	Attempts  int
	TotalTime time.Duration
}

func (e *RetryExhaustedError) Error() string {
	return fmt.Sprintf("retry exhausted after %d attempts (%v): %v", e.Attempts, e.TotalTime, e.LastError)
}

func (e *RetryExhaustedError) Unwrap() error {
	return e.LastError
}

// RetryableFunc is a convenience type for functions that can be retried
type RetryableFunc func() (interface{}, error)

// RetryableContextFunc is a convenience type for context-aware functions that can be retried
type RetryableContextFunc func(context.Context) (interface{}, error)

// Retry is a convenience function for simple retry operations
func Retry(maxAttempts int, initialDelay time.Duration, fn RetryableFunc) (interface{}, error) {
	retrier := NewRetrier(RetryConfig{
		MaxAttempts:  maxAttempts,
		InitialDelay: initialDelay,
		Strategy:     ExponentialBackoff,
		Multiplier:   2.0,
		Jitter:       true,
	})
	return retrier.Execute(fn)
}

// RetryWithContext is a convenience function for simple retry operations with context
func RetryWithContext(ctx context.Context, maxAttempts int, initialDelay time.Duration, fn RetryableContextFunc) (interface{}, error) {
	retrier := NewRetrier(RetryConfig{
		MaxAttempts:  maxAttempts,
		InitialDelay: initialDelay,
		Strategy:     ExponentialBackoff,
		Multiplier:   2.0,
		Jitter:       true,
	})
	return retrier.ExecuteWithContext(ctx, fn)
}

// Common retry predicates

// RetryOnAnyError retries on any error
func RetryOnAnyError(err error) bool {
	return err != nil
}

// RetryOnSpecificErrors retries only on specific error types
func RetryOnSpecificErrors(errorTypes ...error) func(error) bool {
	return func(err error) bool {
		for _, errType := range errorTypes {
			if err.Error() == errType.Error() {
				return true
			}
		}
		return false
	}
}

// RetryOnTemporaryErrors retries on temporary/network errors
func RetryOnTemporaryErrors(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for temporary network errors
	if temp, ok := err.(interface{ Temporary() bool }); ok {
		return temp.Temporary()
	}
	
	// Check for timeout errors
	if timeout, ok := err.(interface{ Timeout() bool }); ok {
		return timeout.Timeout()
	}
	
	// Don't retry circuit breaker errors
	if IsCircuitBreakerError(err) {
		return false
	}
	
	return true
}

// RetryableError wraps an error to indicate it should be retried
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable: %v", e.Err)
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// NonRetryableError wraps an error to indicate it should not be retried
type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string {
	return fmt.Sprintf("non-retryable: %v", e.Err)
}

func (e *NonRetryableError) Unwrap() error {
	return e.Err
}

// IsRetryableError checks if an error is marked as retryable
func IsRetryableError(err error) bool {
	var retryableErr *RetryableError
	return err != nil && (err.Error() != "" || errors.As(err, &retryableErr))
}

// IsNonRetryableError checks if an error is marked as non-retryable
func IsNonRetryableError(err error) bool {
	var nonRetryableErr *NonRetryableError
	return err != nil && errors.As(err, &nonRetryableErr)
}