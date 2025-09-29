package core

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ErrorSeverity defines the severity level of errors
type ErrorSeverity int

const (
	SeverityLow ErrorSeverity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s ErrorSeverity) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// SlothError represents a structured error in the Sloth Runner system
type SlothError struct {
	Code      string
	Message   string
	Details   map[string]interface{}
	Severity  ErrorSeverity
	Retryable bool
	Cause     error
	Context   string
	Timestamp time.Time
	Stack     []byte
}

// Error implements the error interface
func (e *SlothError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *SlothError) Unwrap() error {
	return e.Cause
}

// NewSlothError creates a new structured error
func NewSlothError(code, message string, severity ErrorSeverity) *SlothError {
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)
	
	return &SlothError{
		Code:      code,
		Message:   message,
		Details:   make(map[string]interface{}),
		Severity:  severity,
		Retryable: false,
		Timestamp: time.Now(),
		Stack:     stack[:length],
	}
}

// WithCause adds a cause to the error
func (e *SlothError) WithCause(cause error) *SlothError {
	e.Cause = cause
	return e
}

// WithDetail adds details to the error
func (e *SlothError) WithDetail(key string, value interface{}) *SlothError {
	e.Details[key] = value
	return e
}

// WithContext adds context information
func (e *SlothError) WithContext(context string) *SlothError {
	e.Context = context
	return e
}

// WithRetryable marks the error as retryable
func (e *SlothError) WithRetryable(retryable bool) *SlothError {
	e.Retryable = retryable
	return e
}

// ErrorCollector aggregates errors from multiple operations
type ErrorCollector struct {
	mu     sync.RWMutex
	errors []*SlothError
	maxSize int
}

// NewErrorCollector creates a new error collector
func NewErrorCollector(maxSize int) *ErrorCollector {
	return &ErrorCollector{
		errors:  make([]*SlothError, 0),
		maxSize: maxSize,
	}
}

// Collect adds an error to the collection
func (ec *ErrorCollector) Collect(err error) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	
	var slothErr *SlothError
	if se, ok := err.(*SlothError); ok {
		slothErr = se
	} else {
		slothErr = NewSlothError("GENERIC_ERROR", err.Error(), SeverityMedium).WithCause(err)
	}
	
	ec.errors = append(ec.errors, slothErr)
	
	// Prevent unbounded growth
	if len(ec.errors) > ec.maxSize {
		ec.errors = ec.errors[1:]
	}
}

// GetErrors returns all collected errors
func (ec *ErrorCollector) GetErrors() []*SlothError {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	
	errors := make([]*SlothError, len(ec.errors))
	copy(errors, ec.errors)
	return errors
}

// HasErrors returns true if there are collected errors
func (ec *ErrorCollector) HasErrors() bool {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	
	return len(ec.errors) > 0
}

// GetBySeverity returns errors filtered by severity
func (ec *ErrorCollector) GetBySeverity(severity ErrorSeverity) []*SlothError {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	
	var filtered []*SlothError
	for _, err := range ec.errors {
		if err.Severity == severity {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// Clear removes all collected errors
func (ec *ErrorCollector) Clear() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	
	ec.errors = ec.errors[:0]
}

// ErrorRecovery handles panic recovery and error reporting
type ErrorRecovery struct {
	logger *slog.Logger
	onPanic func(interface{}, []byte)
}

// NewErrorRecovery creates a new error recovery handler
func NewErrorRecovery(logger *slog.Logger) *ErrorRecovery {
	return &ErrorRecovery{
		logger: logger,
		onPanic: func(v interface{}, stack []byte) {
			logger.Error("Panic recovered", 
				"panic", v, 
				"stack", string(stack))
		},
	}
}

// WithPanicHandler sets a custom panic handler
func (er *ErrorRecovery) WithPanicHandler(handler func(interface{}, []byte)) *ErrorRecovery {
	er.onPanic = handler
	return er
}

// SafeExecute executes a function with panic recovery
func (er *ErrorRecovery) SafeExecute(fn func() error) (recovered bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			recovered = true
			stack := make([]byte, 4096)
			length := runtime.Stack(stack, false)
			
			er.onPanic(r, stack[:length])
			
			err = NewSlothError("PANIC_RECOVERED", 
				fmt.Sprintf("Panic recovered: %v", r),
				SeverityCritical).WithDetail("panic_value", r)
		}
	}()
	
	err = fn()
	return false, err
}

// SafeGo runs a goroutine with panic recovery
func (er *ErrorRecovery) SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)
				er.onPanic(r, stack[:length])
			}
		}()
		fn()
	}()
}

// ErrorAggregator aggregates multiple errors into a single error
type ErrorAggregator struct {
	errors []error
	mu     sync.Mutex
}

// NewErrorAggregator creates a new error aggregator
func NewErrorAggregator() *ErrorAggregator {
	return &ErrorAggregator{}
}

// Add adds an error to the aggregator
func (ea *ErrorAggregator) Add(err error) {
	if err == nil {
		return
	}
	
	ea.mu.Lock()
	defer ea.mu.Unlock()
	
	ea.errors = append(ea.errors, err)
}

// ToError converts aggregated errors to a single error
func (ea *ErrorAggregator) ToError() error {
	ea.mu.Lock()
	defer ea.mu.Unlock()
	
	if len(ea.errors) == 0 {
		return nil
	}
	
	if len(ea.errors) == 1 {
		return ea.errors[0]
	}
	
	var messages []string
	severity := SeverityLow
	
	for _, err := range ea.errors {
		messages = append(messages, err.Error())
		
		if se, ok := err.(*SlothError); ok {
			if se.Severity > severity {
				severity = se.Severity
			}
		}
	}
	
	return NewSlothError("MULTIPLE_ERRORS", 
		fmt.Sprintf("Multiple errors occurred: %s", strings.Join(messages, "; ")),
		severity).WithDetail("error_count", len(ea.errors))
}

// HasErrors returns true if there are aggregated errors
func (ea *ErrorAggregator) HasErrors() bool {
	ea.mu.Lock()
	defer ea.mu.Unlock()
	
	return len(ea.errors) > 0
}

// Count returns the number of aggregated errors
func (ea *ErrorAggregator) Count() int {
	ea.mu.Lock()
	defer ea.mu.Unlock()
	
	return len(ea.errors)
}

// Clear removes all aggregated errors
func (ea *ErrorAggregator) Clear() {
	ea.mu.Lock()
	defer ea.mu.Unlock()
	
	ea.errors = ea.errors[:0]
}

// Timeout management for better error handling
type TimeoutManager struct {
	timeouts map[string]context.CancelFunc
	mu       sync.RWMutex
	logger   *slog.Logger
}

// NewTimeoutManager creates a new timeout manager
func NewTimeoutManager(logger *slog.Logger) *TimeoutManager {
	return &TimeoutManager{
		timeouts: make(map[string]context.CancelFunc),
		logger:   logger,
	}
}

// StartTimeout starts a timeout for an operation
func (tm *TimeoutManager) StartTimeout(id string, duration time.Duration, onTimeout func()) context.Context {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	// Cancel existing timeout if any
	if cancel, exists := tm.timeouts[id]; exists {
		cancel()
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	tm.timeouts[id] = cancel
	
	// Monitor for timeout
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			tm.logger.Warn("Operation timed out", "id", id, "duration", duration)
			if onTimeout != nil {
				onTimeout()
			}
		}
		
		tm.mu.Lock()
		delete(tm.timeouts, id)
		tm.mu.Unlock()
	}()
	
	return ctx
}

// CancelTimeout cancels a specific timeout
func (tm *TimeoutManager) CancelTimeout(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if cancel, exists := tm.timeouts[id]; exists {
		cancel()
		delete(tm.timeouts, id)
	}
}

// CancelAll cancels all active timeouts
func (tm *TimeoutManager) CancelAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	for id, cancel := range tm.timeouts {
		cancel()
		delete(tm.timeouts, id)
	}
}

// Error classification utilities
func ClassifyError(err error) ErrorSeverity {
	if err == nil {
		return SeverityLow
	}
	
	errStr := strings.ToLower(err.Error())
	
	// Critical errors
	if strings.Contains(errStr, "panic") ||
		strings.Contains(errStr, "fatal") ||
		strings.Contains(errStr, "deadlock") {
		return SeverityCritical
	}
	
	// High severity errors
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "permission denied") ||
		strings.Contains(errStr, "not found") {
		return SeverityHigh
	}
	
	// Medium severity errors
	if strings.Contains(errStr, "invalid") ||
		strings.Contains(errStr, "failed to") ||
		strings.Contains(errStr, "cannot") {
		return SeverityMedium
	}
	
	// Default to low
	return SeverityLow
}

// IsRetryableError determines if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's a SlothError with explicit retryable flag
	if se, ok := err.(*SlothError); ok {
		return se.Retryable
	}
	
	errStr := strings.ToLower(err.Error())
	
	// Non-retryable errors
	nonRetryable := []string{
		"permission denied",
		"unauthorized",
		"forbidden",
		"not found",
		"invalid argument",
		"bad request",
		"conflict",
	}
	
	for _, pattern := range nonRetryable {
		if strings.Contains(errStr, pattern) {
			return false
		}
	}
	
	// Retryable errors
	retryable := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"service unavailable",
		"too many requests",
		"internal server error",
		"bad gateway",
		"gateway timeout",
	}
	
	for _, pattern := range retryable {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	
	// Default to non-retryable for unknown errors
	return false
}