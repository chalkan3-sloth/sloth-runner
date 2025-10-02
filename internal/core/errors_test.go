package core

import (
	"errors"
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
