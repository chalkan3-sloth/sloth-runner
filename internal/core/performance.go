package core

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// PerformanceMetrics tracks runtime performance metrics
type PerformanceMetrics struct {
	mu                sync.RWMutex
	TaskExecutions    int64
	TotalDuration     time.Duration
	MemoryUsage       runtime.MemStats
	GoroutineCount    int
	LastGCTime        time.Time
	PeakMemoryUsage   uint64
	CacheHitRatio     float64
	ConnectionsActive int64
}

// NewPerformanceMetrics creates a new performance metrics tracker
func NewPerformanceMetrics() *PerformanceMetrics {
	pm := &PerformanceMetrics{}
	pm.startMonitoring()
	return pm
}

// startMonitoring begins background monitoring
func (pm *PerformanceMetrics) startMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			pm.updateMetrics()
		}
	}()
}

// updateMetrics collects current runtime metrics
func (pm *PerformanceMetrics) updateMetrics() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	runtime.ReadMemStats(&pm.MemoryUsage)
	pm.GoroutineCount = runtime.NumGoroutine()
	
	if pm.MemoryUsage.Alloc > pm.PeakMemoryUsage {
		pm.PeakMemoryUsage = pm.MemoryUsage.Alloc
	}
}

// RecordTaskExecution records a task execution
func (pm *PerformanceMetrics) RecordTaskExecution(duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.TaskExecutions++
	pm.TotalDuration += duration
}

// GetSnapshot returns current metrics snapshot
func (pm *PerformanceMetrics) GetSnapshot() PerformanceMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	snapshot := *pm
	return snapshot
}

// ResourcePool manages pooled resources for better performance
type ResourcePool struct {
	buffers      sync.Pool
	httpClients  sync.Pool
	luaStates    sync.Pool
	maxPoolSize  int
	currentSize  int64
	mu           sync.RWMutex
}

// NewResourcePool creates a new resource pool
func NewResourcePool(maxSize int) *ResourcePool {
	return &ResourcePool{
		maxPoolSize: maxSize,
		buffers: sync.Pool{
			New: func() interface{} {
				// Pre-allocate 4KB buffers
				return make([]byte, 0, 4096)
			},
		},
		httpClients: sync.Pool{
			New: func() interface{} {
				return &HttpClientConfig{
					Timeout:    30 * time.Second,
					MaxRetries: 3,
				}
			},
		},
	}
}

// GetBuffer gets a buffer from the pool
func (rp *ResourcePool) GetBuffer() []byte {
	return rp.buffers.Get().([]byte)
}

// PutBuffer returns a buffer to the pool
func (rp *ResourcePool) PutBuffer(buf []byte) {
	// Reset length but preserve capacity
	buf = buf[:0]
	rp.buffers.Put(buf)
}

// CircuitBreaker implements circuit breaker pattern for reliability
type CircuitBreaker struct {
	mu            sync.RWMutex
	name          string
	maxFailures   int64
	resetTimeout  time.Duration
	state         CircuitState
	failures      int64
	successes     int64
	lastFailure   time.Time
	lastSuccess   time.Time
	halfOpenMax   int64
}

type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, maxFailures int64, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:         name,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
		halfOpenMax:  5, // Allow 5 requests in half-open state
	}
}

// Execute wraps function execution with circuit breaker logic
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.allowRequest() {
		return fmt.Errorf("circuit breaker '%s' is open", cb.name)
	}
	
	err := fn()
	cb.recordResult(err)
	return err
}

// allowRequest determines if a request should be allowed
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	
	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		return time.Since(cb.lastFailure) > cb.resetTimeout
	case StateHalfOpen:
		return cb.successes < cb.halfOpenMax
	default:
		return false
	}
}

// recordResult records the result of a request
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		
		if cb.state == StateClosed && cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		} else if cb.state == StateHalfOpen {
			cb.state = StateOpen
			cb.successes = 0
		}
	} else {
		cb.successes++
		cb.lastSuccess = time.Now()
		
		if cb.state == StateHalfOpen && cb.successes >= cb.halfOpenMax {
			cb.state = StateClosed
			cb.failures = 0
		}
	}
	
	// Transition from open to half-open
	if cb.state == StateOpen && time.Since(cb.lastFailure) > cb.resetTimeout {
		cb.state = StateHalfOpen
		cb.successes = 0
	}
}

// GetStats returns current circuit breaker statistics
func (cb *CircuitBreaker) GetStats() CircuitStats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	
	return CircuitStats{
		Name:        cb.name,
		State:       cb.state.String(),
		Failures:    cb.failures,
		Successes:   cb.successes,
		LastFailure: cb.lastFailure,
		LastSuccess: cb.lastSuccess,
	}
}

type CircuitStats struct {
	Name        string
	State       string
	Failures    int64
	Successes   int64
	LastFailure time.Time
	LastSuccess time.Time
}

// RetryConfig defines retry behavior configuration
type RetryConfig struct {
	MaxAttempts   int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	Multiplier    float64
	Jitter        bool
	RetryableFunc func(error) bool
}

// DefaultRetryConfig returns sensible defaults
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableFunc: func(err error) bool {
			// Default: retry on most errors except context cancellation
			return err != nil && err != context.Canceled
		},
	}
}

// HttpClientConfig optimized HTTP client configuration
type HttpClientConfig struct {
	Timeout         time.Duration
	MaxRetries      int
	RetryDelay      time.Duration
	MaxConnections  int
	IdleConnTimeout time.Duration
	KeepAlive       time.Duration
}

// ContextManager manages context lifecycle and cancellation
type ContextManager struct {
	contexts map[string]context.CancelFunc
	mu       sync.RWMutex
}

// NewContextManager creates a new context manager
func NewContextManager() *ContextManager {
	return &ContextManager{
		contexts: make(map[string]context.CancelFunc),
	}
}

// CreateContext creates a new managed context
func (cm *ContextManager) CreateContext(id string, timeout time.Duration) context.Context {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// Cancel existing context if any
	if cancel, exists := cm.contexts[id]; exists {
		cancel()
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	cm.contexts[id] = cancel
	
	return ctx
}

// CancelContext cancels a specific context
func (cm *ContextManager) CancelContext(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if cancel, exists := cm.contexts[id]; exists {
		cancel()
		delete(cm.contexts, id)
	}
}

// CancelAll cancels all managed contexts
func (cm *ContextManager) CancelAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	for id, cancel := range cm.contexts {
		cancel()
		delete(cm.contexts, id)
	}
}

// MemoryManager helps prevent memory leaks and manage allocations
type MemoryManager struct {
	allocations map[string]interface{}
	mu          sync.RWMutex
	maxSize     int64
	currentSize int64
}

// NewMemoryManager creates a new memory manager
func NewMemoryManager(maxSizeMB int64) *MemoryManager {
	return &MemoryManager{
		allocations: make(map[string]interface{}),
		maxSize:     maxSizeMB * 1024 * 1024, // Convert MB to bytes
	}
}

// Track tracks a memory allocation
func (mm *MemoryManager) Track(id string, obj interface{}, estimatedSize int64) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	if mm.currentSize+estimatedSize > mm.maxSize {
		return fmt.Errorf("memory limit exceeded: %d + %d > %d", 
			mm.currentSize, estimatedSize, mm.maxSize)
	}
	
	mm.allocations[id] = obj
	mm.currentSize += estimatedSize
	
	return nil
}

// Untrack removes tracking for an allocation
func (mm *MemoryManager) Untrack(id string, estimatedSize int64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	delete(mm.allocations, id)
	mm.currentSize -= estimatedSize
	
	if mm.currentSize < 0 {
		mm.currentSize = 0
	}
}

// GetStats returns memory management statistics
func (mm *MemoryManager) GetStats() MemoryStats {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	return MemoryStats{
		TrackedAllocations: len(mm.allocations),
		CurrentSize:        mm.currentSize,
		MaxSize:           mm.maxSize,
		UsagePercent:      float64(mm.currentSize) / float64(mm.maxSize) * 100,
	}
}

type MemoryStats struct {
	TrackedAllocations int
	CurrentSize        int64
	MaxSize           int64
	UsagePercent      float64
}