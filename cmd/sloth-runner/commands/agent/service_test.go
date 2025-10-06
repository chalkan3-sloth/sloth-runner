package agent

import (
	"testing"
)

// Test NewDefaultConnectionFactory
func TestNewDefaultConnectionFactory(t *testing.T) {
	factory := NewDefaultConnectionFactory()
	if factory == nil {
		t.Fatal("Expected non-nil factory")
	}
}

// Test CreateAgentClient - just test that it attempts connection
func TestCreateAgentClient(t *testing.T) {
	factory := NewDefaultConnectionFactory()

	// Test with invalid address (will fail to connect)
	_, cleanup, err := factory.CreateAgentClient("invalid:99999")

	// We expect an error since the address is invalid
	if err == nil {
		t.Error("Expected error with invalid address")
		if cleanup != nil {
			cleanup()
		}
	}
}

// Test CreateRegistryClient - just test that it attempts connection
func TestCreateRegistryClient(t *testing.T) {
	factory := NewDefaultConnectionFactory()

	// Test with invalid address (will fail to connect)
	_, cleanup, err := factory.CreateRegistryClient("invalid:99999")

	// We expect an error since the address is invalid
	if err == nil {
		t.Error("Expected error with invalid address")
		if cleanup != nil {
			cleanup()
		}
	}
}

// Test NewAgentService
func TestNewAgentService(t *testing.T) {
	// This function is very simple - just constructs a struct
	// We can test it with nil clients since we're only testing construction
	service := NewAgentService(nil, nil)
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
}
