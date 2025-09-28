package reliability

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateHalfOpen
	StateOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig holds configuration for a circuit breaker
type CircuitBreakerConfig struct {
	MaxFailures     int           // Maximum failures before opening
	Timeout         time.Duration // How long to wait before transitioning to half-open
	SuccessThreshold int          // Successes needed in half-open to close
	OnStateChange   func(from, to CircuitBreakerState) // Optional state change callback
}

// CircuitBreakerStats holds statistics about the circuit breaker
type CircuitBreakerStats struct {
	Requests      uint64
	TotalSuccess  uint64
	TotalFailures uint64
	ConsecutiveSuccess uint64
	ConsecutiveFailures uint64
	LastSuccessTime time.Time
	LastFailureTime time.Time
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config   CircuitBreakerConfig
	state    CircuitBreakerState
	stats    CircuitBreakerStats
	mu       sync.RWMutex
	lastStateChange time.Time
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.MaxFailures <= 0 {
		config.MaxFailures = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 60 * time.Second
	}
	if config.SuccessThreshold <= 0 {
		config.SuccessThreshold = 1
	}

	return &CircuitBreaker{
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Execute executes the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	if !cb.canExecute() {
		cb.mu.RLock()
		state := cb.state
		cb.mu.RUnlock()
		return nil, &CircuitBreakerError{
			State:   state,
			Message: fmt.Sprintf("circuit breaker is %s", state),
		}
	}

	cb.mu.Lock()
	cb.stats.Requests++
	cb.mu.Unlock()

	result, err := fn()

	if err != nil {
		cb.onFailure()
		return result, err
	}

	cb.onSuccess()
	return result, nil
}

// ExecuteWithContext executes the given function with context and circuit breaker protection
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	if !cb.canExecute() {
		cb.mu.RLock()
		state := cb.state
		cb.mu.RUnlock()
		return nil, &CircuitBreakerError{
			State:   state,
			Message: fmt.Sprintf("circuit breaker is %s", state),
		}
	}

	cb.mu.Lock()
	cb.stats.Requests++
	cb.mu.Unlock()

	result, err := fn(ctx)

	if err != nil {
		cb.onFailure()
		return result, err
	}

	cb.onSuccess()
	return result, nil
}

// canExecute determines if the circuit breaker allows execution
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if timeout has elapsed
		return time.Since(cb.lastStateChange) >= cb.config.Timeout
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// onSuccess handles successful execution
func (cb *CircuitBreaker) onSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.stats.TotalSuccess++
	cb.stats.ConsecutiveSuccess++
	cb.stats.ConsecutiveFailures = 0
	cb.stats.LastSuccessTime = time.Now()

	switch cb.state {
	case StateHalfOpen:
		if cb.stats.ConsecutiveSuccess >= uint64(cb.config.SuccessThreshold) {
			cb.setState(StateClosed)
		}
	case StateOpen:
		// Transition to half-open after timeout
		if time.Since(cb.lastStateChange) >= cb.config.Timeout {
			cb.setState(StateHalfOpen)
		}
	}
}

// onFailure handles failed execution
func (cb *CircuitBreaker) onFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.stats.TotalFailures++
	cb.stats.ConsecutiveFailures++
	cb.stats.ConsecutiveSuccess = 0
	cb.stats.LastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.stats.ConsecutiveFailures >= uint64(cb.config.MaxFailures) {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		cb.setState(StateOpen)
	}
}

// setState changes the circuit breaker state
func (cb *CircuitBreaker) setState(newState CircuitBreakerState) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()

	// Reset consecutive counters on state change
	if newState == StateClosed {
		cb.stats.ConsecutiveFailures = 0
	} else if newState == StateHalfOpen {
		cb.stats.ConsecutiveSuccess = 0
	}

	// Call state change callback if provided
	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(oldState, newState)
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetStats returns a copy of the current statistics
func (cb *CircuitBreaker) GetStats() CircuitBreakerStats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.stats
}

// Reset resets the circuit breaker to its initial state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	oldState := cb.state
	cb.state = StateClosed
	cb.stats = CircuitBreakerStats{}
	cb.lastStateChange = time.Now()

	if cb.config.OnStateChange != nil && oldState != StateClosed {
		go cb.config.OnStateChange(oldState, StateClosed)
	}
}

// CircuitBreakerError represents an error from the circuit breaker
type CircuitBreakerError struct {
	State   CircuitBreakerState
	Message string
}

func (e *CircuitBreakerError) Error() string {
	return e.Message
}

// IsCircuitBreakerError checks if an error is a circuit breaker error
func IsCircuitBreakerError(err error) bool {
	var cbErr *CircuitBreakerError
	return errors.As(err, &cbErr)
}

// CircuitBreakerManager manages multiple named circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate gets an existing circuit breaker or creates a new one
func (cbm *CircuitBreakerManager) GetOrCreate(name string, config CircuitBreakerConfig) *CircuitBreaker {
	cbm.mu.RLock()
	if cb, exists := cbm.breakers[name]; exists {
		cbm.mu.RUnlock()
		return cb
	}
	cbm.mu.RUnlock()

	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	// Double-check after acquiring write lock
	if cb, exists := cbm.breakers[name]; exists {
		return cb
	}

	cb := NewCircuitBreaker(config)
	cbm.breakers[name] = cb
	return cb
}

// Get gets an existing circuit breaker by name
func (cbm *CircuitBreakerManager) Get(name string) (*CircuitBreaker, bool) {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()
	cb, exists := cbm.breakers[name]
	return cb, exists
}

// Remove removes a circuit breaker by name
func (cbm *CircuitBreakerManager) Remove(name string) {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()
	delete(cbm.breakers, name)
}

// List returns all circuit breaker names
func (cbm *CircuitBreakerManager) List() []string {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	names := make([]string, 0, len(cbm.breakers))
	for name := range cbm.breakers {
		names = append(names, name)
	}
	return names
}

// GetAllStats returns statistics for all circuit breakers
func (cbm *CircuitBreakerManager) GetAllStats() map[string]CircuitBreakerStats {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	stats := make(map[string]CircuitBreakerStats)
	for name, cb := range cbm.breakers {
		stats[name] = cb.GetStats()
	}
	return stats
}