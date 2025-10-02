package core

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestErrorSeverity_String(t *testing.T) {
	tests := []struct {
		severity ErrorSeverity
		expected string
	}{
		{SeverityLow, "low"},
		{SeverityMedium, "medium"},
		{SeverityHigh, "high"},
		{SeverityCritical, "critical"},
		{ErrorSeverity(999), "unknown"},
	}

	for _, tt := range tests {
		result := tt.severity.String()
		if result != tt.expected {
			t.Errorf("Severity %d: expected '%s', got '%s'", tt.severity, tt.expected, result)
		}
	}
}

func TestNewSlothError(t *testing.T) {
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh)

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if err.Code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", err.Code)
	}

	if err.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", err.Message)
	}

	if err.Severity != SeverityHigh {
		t.Errorf("Expected severity High, got %v", err.Severity)
	}

	if err.Details == nil {
		t.Error("Expected non-nil Details map")
	}

	if err.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}

	if len(err.Stack) == 0 {
		t.Error("Expected stack trace")
	}
}

func TestSlothError_Error(t *testing.T) {
	// Test without cause
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh)
	expected := "[TEST_CODE] Test message"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test with cause
	cause := errors.New("root cause")
	err = NewSlothError("TEST_CODE", "Test message", SeverityHigh).WithCause(cause)
	if !strings.Contains(err.Error(), "root cause") {
		t.Errorf("Expected error message to contain cause, got '%s'", err.Error())
	}
}

func TestSlothError_WithCause(t *testing.T) {
	cause := errors.New("root cause")
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh).WithCause(cause)

	if err.Cause != cause {
		t.Error("Expected cause to be set")
	}

	// Test fluent API
	if err2 := err.WithCause(cause); err2 != err {
		t.Error("Expected WithCause to return same instance")
	}
}

func TestSlothError_WithDetail(t *testing.T) {
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh).
		WithDetail("key1", "value1").
		WithDetail("key2", 123)

	if err.Details["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got '%v'", err.Details["key1"])
	}

	if err.Details["key2"] != 123 {
		t.Errorf("Expected key2=123, got '%v'", err.Details["key2"])
	}
}

func TestSlothError_WithContext(t *testing.T) {
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh).
		WithContext("test context")

	if err.Context != "test context" {
		t.Errorf("Expected context 'test context', got '%s'", err.Context)
	}
}

func TestSlothError_WithRetryable(t *testing.T) {
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh).
		WithRetryable(true)

	if !err.Retryable {
		t.Error("Expected error to be retryable")
	}

	err.WithRetryable(false)
	if err.Retryable {
		t.Error("Expected error to not be retryable")
	}
}

func TestSlothError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh).WithCause(cause)

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Error("Expected unwrapped error to be the cause")
	}

	// Test without cause
	err2 := NewSlothError("TEST_CODE", "Test message", SeverityHigh)
	if err2.Unwrap() != nil {
		t.Error("Expected nil for error without cause")
	}
}

func TestSlothError_FluentAPI(t *testing.T) {
	cause := errors.New("root cause")
	err := NewSlothError("TEST_CODE", "Test message", SeverityHigh).
		WithCause(cause).
		WithDetail("key1", "value1").
		WithDetail("key2", 123).
		WithContext("test context").
		WithRetryable(true)

	// Verify all fields are set correctly
	if err.Code != "TEST_CODE" {
		t.Error("Code not set correctly")
	}
	if err.Cause != cause {
		t.Error("Cause not set correctly")
	}
	if len(err.Details) != 2 {
		t.Error("Details not set correctly")
	}
	if err.Context != "test context" {
		t.Error("Context not set correctly")
	}
	if !err.Retryable {
		t.Error("Retryable not set correctly")
	}
}

func TestSlothError_MultipleSeverities(t *testing.T) {
	severities := []ErrorSeverity{
		SeverityLow,
		SeverityMedium,
		SeverityHigh,
		SeverityCritical,
	}

	for _, sev := range severities {
		err := NewSlothError("TEST", "message", sev)
		if err.Severity != sev {
			t.Errorf("Expected severity %v, got %v", sev, err.Severity)
		}
	}
}

func TestErrorCollector_Collect(t *testing.T) {
	collector := NewErrorCollector(10)
	
	err1 := errors.New("error 1")
	err2 := NewSlothError("CODE2", "error 2", SeverityHigh)
	
	collector.Collect(err1)
	collector.Collect(err2)
	
	if !collector.HasErrors() {
		t.Error("Expected collector to have errors")
	}
	
	collectedErrors := collector.GetErrors()
	if len(collectedErrors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(collectedErrors))
	}
}

func TestErrorCollector_MaxSize(t *testing.T) {
	collector := NewErrorCollector(3)
	
	// Add more errors than max size
	for i := 0; i < 5; i++ {
		collector.Collect(errors.New("error"))
	}
	
	errors := collector.GetErrors()
	if len(errors) > 3 {
		t.Errorf("Expected max 3 errors, got %d", len(errors))
	}
}

func TestErrorCollector_GetBySeverity(t *testing.T) {
	collector := NewErrorCollector(10)
	
	collector.Collect(NewSlothError("E1", "error 1", SeverityLow))
	collector.Collect(NewSlothError("E2", "error 2", SeverityHigh))
	collector.Collect(NewSlothError("E3", "error 3", SeverityHigh))
	collector.Collect(NewSlothError("E4", "error 4", SeverityLow))
	
	highErrors := collector.GetBySeverity(SeverityHigh)
	if len(highErrors) != 2 {
		t.Errorf("Expected 2 high severity errors, got %d", len(highErrors))
	}
	
	lowErrors := collector.GetBySeverity(SeverityLow)
	if len(lowErrors) != 2 {
		t.Errorf("Expected 2 low severity errors, got %d", len(lowErrors))
	}
}

func TestErrorCollector_Clear(t *testing.T) {
	collector := NewErrorCollector(10)
	
	collector.Collect(errors.New("error 1"))
	collector.Collect(errors.New("error 2"))
	
	if !collector.HasErrors() {
		t.Error("Expected collector to have errors")
	}
	
	collector.Clear()
	
	if collector.HasErrors() {
		t.Error("Expected collector to be empty after Clear")
	}
}

func TestErrorRecovery_SafeExecute(t *testing.T) {
	// Create a simple logger for testing
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	recovery := NewErrorRecovery(logger)
	
	// Test normal execution
	recovered, err := recovery.SafeExecute(func() error {
		return nil
	})
	
	if recovered {
		t.Error("Expected no recovery for normal execution")
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Test with error
	testErr := errors.New("test error")
	recovered, err = recovery.SafeExecute(func() error {
		return testErr
	})
	
	if recovered {
		t.Error("Expected no recovery for returned error")
	}
	if err != testErr {
		t.Errorf("Expected test error, got %v", err)
	}
}

func TestErrorRecovery_SafeExecute_Panic(t *testing.T) {
	// Create a simple logger for testing
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	recovery := NewErrorRecovery(logger)
	
	recovered, err := recovery.SafeExecute(func() error {
		panic("test panic")
	})
	
	if !recovered {
		t.Error("Expected panic to be recovered")
	}
	
	if err == nil {
		t.Error("Expected error after panic recovery")
	}
	
	slothErr, ok := err.(*SlothError)
	if !ok {
		t.Error("Expected SlothError after panic")
	}
	
	if slothErr.Code != "PANIC_RECOVERED" {
		t.Errorf("Expected PANIC_RECOVERED code, got %s", slothErr.Code)
	}
}

func TestErrorAggregator_Add(t *testing.T) {
	aggregator := NewErrorAggregator()
	
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	
	aggregator.Add(err1)
	aggregator.Add(err2)
	aggregator.Add(nil) // Should be ignored
	
	if !aggregator.HasErrors() {
		t.Error("Expected aggregator to have errors")
	}
	
	if aggregator.Count() != 2 {
		t.Errorf("Expected 2 errors, got %d", aggregator.Count())
	}
}

func TestErrorAggregator_ToError(t *testing.T) {
	aggregator := NewErrorAggregator()
	
	// No errors
	err := aggregator.ToError()
	if err != nil {
		t.Error("Expected nil error when aggregator is empty")
	}
	
	// Single error
	singleErr := errors.New("single error")
	aggregator.Add(singleErr)
	err = aggregator.ToError()
	if err != singleErr {
		t.Error("Expected same error for single error")
	}
	
	// Multiple errors
	aggregator.Clear()
	aggregator.Add(errors.New("error 1"))
	aggregator.Add(errors.New("error 2"))
	err = aggregator.ToError()
	
	if err == nil {
		t.Error("Expected error for multiple errors")
	}
	
	if !strings.Contains(err.Error(), "Multiple errors occurred") {
		t.Errorf("Expected 'Multiple errors occurred', got %s", err.Error())
	}
}

func TestErrorAggregator_Clear(t *testing.T) {
	aggregator := NewErrorAggregator()
	
	aggregator.Add(errors.New("error 1"))
	aggregator.Add(errors.New("error 2"))
	
	if !aggregator.HasErrors() {
		t.Error("Expected aggregator to have errors")
	}
	
	aggregator.Clear()
	
	if aggregator.HasErrors() {
		t.Error("Expected aggregator to be empty after Clear")
	}
	
	if aggregator.Count() != 0 {
		t.Errorf("Expected 0 errors after Clear, got %d", aggregator.Count())
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		err      error
		expected ErrorSeverity
	}{
		{nil, SeverityLow},
		{errors.New("panic occurred"), SeverityCritical},
		{errors.New("fatal error"), SeverityCritical},
		{errors.New("timeout exceeded"), SeverityHigh},
		{errors.New("connection refused"), SeverityHigh},
		{errors.New("permission denied"), SeverityHigh},
		{errors.New("file not found"), SeverityHigh},
		{errors.New("invalid input"), SeverityMedium},
		{errors.New("failed to connect"), SeverityMedium},
		{errors.New("cannot parse"), SeverityMedium},
		{errors.New("some other error"), SeverityLow},
	}
	
	for _, tt := range tests {
		result := ClassifyError(tt.err)
		if result != tt.expected {
			errMsg := "nil"
			if tt.err != nil {
				errMsg = tt.err.Error()
			}
			t.Errorf("ClassifyError(%s): expected %s, got %s", 
				errMsg, tt.expected.String(), result.String())
		}
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("permission denied"), false},
		{errors.New("unauthorized"), false},
		{errors.New("not found"), false},
		{errors.New("invalid argument"), false},
		{errors.New("timeout exceeded"), true},
		{errors.New("connection refused"), true},
		{errors.New("connection reset"), true},
		{errors.New("temporary failure"), true},
		{errors.New("service unavailable"), true},
		{errors.New("too many requests"), true},
		{errors.New("internal server error"), true},
		{errors.New("bad gateway"), true},
		{errors.New("gateway timeout"), true},
		{errors.New("some random error"), false},
	}
	
	for _, tt := range tests {
		result := IsRetryableError(tt.err)
		if result != tt.expected {
			errMsg := "nil"
			if tt.err != nil {
				errMsg = tt.err.Error()
			}
			t.Errorf("IsRetryableError(%s): expected %v, got %v", 
				errMsg, tt.expected, result)
		}
	}
}

func TestIsRetryableError_SlothError(t *testing.T) {
	// Test with explicit retryable flag
	retryableErr := NewSlothError("TEST", "test error", SeverityHigh).
		WithRetryable(true)
	
	if !IsRetryableError(retryableErr) {
		t.Error("Expected retryable error to return true")
	}
	
	nonRetryableErr := NewSlothError("TEST", "test error", SeverityHigh).
		WithRetryable(false)
	
	if IsRetryableError(nonRetryableErr) {
		t.Error("Expected non-retryable error to return false")
	}
}
