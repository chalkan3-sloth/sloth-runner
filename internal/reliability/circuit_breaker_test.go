package reliability

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker(t *testing.T) {
	t.Run("Initial State Closed", func(t *testing.T) {
		cb := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures: 3,
			Timeout:     1 * time.Second,
		})

		assert.Equal(t, StateClosed, cb.GetState())
	})

	t.Run("Success Execution", func(t *testing.T) {
		cb := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures: 3,
			Timeout:     1 * time.Second,
		})

		result, err := cb.Execute(func() (interface{}, error) {
			return "success", nil
		})

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, StateClosed, cb.GetState())
	})

	t.Run("Failure Opens Circuit", func(t *testing.T) {
		cb := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures: 2,
			Timeout:     1 * time.Second,
		})

		// First failure
		_, err := cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failure 1")
		})
		assert.Error(t, err)
		assert.Equal(t, StateClosed, cb.GetState())

		// Second failure - should open circuit
		_, err = cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failure 2")
		})
		assert.Error(t, err)
		assert.Equal(t, StateOpen, cb.GetState())
	})

	t.Run("Open Circuit Rejects Calls", func(t *testing.T) {
		cb := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures: 1,
			Timeout:     100 * time.Millisecond,
		})

		// Fail to open circuit
		cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failure")
		})

		assert.Equal(t, StateOpen, cb.GetState())

		// Should reject calls
		_, err := cb.Execute(func() (interface{}, error) {
			return "should not execute", nil
		})

		assert.Error(t, err)
		assert.True(t, IsCircuitBreakerError(err))
	})

	t.Run("Half-Open Transition", func(t *testing.T) {
		cb := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures:      1,
			Timeout:          50 * time.Millisecond,
			SuccessThreshold: 2,
		})

		// Open circuit
		cb.Execute(func() (interface{}, error) {
			return nil, errors.New("failure")
		})
		assert.Equal(t, StateOpen, cb.GetState())

		// Wait for timeout
		time.Sleep(60 * time.Millisecond)

		// Should allow one call and transition to half-open
		_, err := cb.Execute(func() (interface{}, error) {
			return "success", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, StateHalfOpen, cb.GetState())
	})
}

func TestCircuitBreakerManager(t *testing.T) {
	manager := NewCircuitBreakerManager()

	t.Run("Create and Retrieve", func(t *testing.T) {
		config := CircuitBreakerConfig{MaxFailures: 5}
		cb1 := manager.GetOrCreate("test1", config)
		cb2 := manager.GetOrCreate("test1", config) // Should return same instance

		assert.Same(t, cb1, cb2)

		cb3, exists := manager.Get("test1")
		assert.True(t, exists)
		assert.Same(t, cb1, cb3)
	})

	t.Run("List All", func(t *testing.T) {
		config := CircuitBreakerConfig{MaxFailures: 5}
		manager.GetOrCreate("test2", config)
		manager.GetOrCreate("test3", config)

		names := manager.List()
		assert.Contains(t, names, "test1")
		assert.Contains(t, names, "test2")
		assert.Contains(t, names, "test3")
	})

	t.Run("Remove", func(t *testing.T) {
		manager.Remove("test2")
		_, exists := manager.Get("test2")
		assert.False(t, exists)
	})
}

func TestCircuitBreakerWithContext(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures: 2,
		Timeout:     1 * time.Second,
	})

	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		
		// Start execution in goroutine and cancel context
		var result interface{}
		var err error
		done := make(chan bool)
		
		go func() {
			result, err = cb.ExecuteWithContext(ctx, func(ctx context.Context) (interface{}, error) {
				// This should be interrupted by context cancellation
				select {
				case <-time.After(1 * time.Second):
					return "should not complete", nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			})
			done <- true
		}()
		
		// Cancel context after a short delay
		time.Sleep(10 * time.Millisecond)
		cancel()
		
		// Wait for execution to complete
		<-done

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		assert.Nil(t, result)
	})

	t.Run("Success With Context", func(t *testing.T) {
		ctx := context.Background()

		result, err := cb.ExecuteWithContext(ctx, func(ctx context.Context) (interface{}, error) {
			return "success with context", nil
		})

		assert.NoError(t, err)
		assert.Equal(t, "success with context", result)
	})
}