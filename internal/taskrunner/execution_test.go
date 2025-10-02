package taskrunner

import (
	"errors"
	"testing"
)

func TestExecuteShellCondition_Success(t *testing.T) {
	// Command that should succeed
	result, err := executeShellCondition("true")
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !result {
		t.Error("Expected command to succeed")
	}
}

func TestExecuteShellCondition_Failure(t *testing.T) {
	// Command that should fail
	result, err := executeShellCondition("false")
	
	if err != nil {
		t.Errorf("Expected no error for failed command, got %v", err)
	}
	
	if result {
		t.Error("Expected command to fail")
	}
}

func TestExecuteShellCondition_EmptyCommand(t *testing.T) {
	result, err := executeShellCondition("")
	
	if err == nil {
		t.Error("Expected error for empty command")
	}
	
	if result {
		t.Error("Expected false result for empty command")
	}
	
	if err.Error() != "command cannot be empty" {
		t.Errorf("Expected 'command cannot be empty' error, got %v", err)
	}
}

func TestExecuteShellCondition_CommandNotFound(t *testing.T) {
	// bash -c with a non-existent command will fail with exit code 127
	// but it will return false, not an error (since bash itself executed successfully)
	result, _ := executeShellCondition("nonexistentcommand12345xyz")
	
	// The bash command will execute and return a non-zero exit code
	// This should result in false with no error
	if result {
		t.Error("Expected false result for non-existent command")
	}
	
	// It's ok to have an error or no error depending on bash behavior
	// The important part is that result is false
}

func TestExecuteShellCondition_ExitCode(t *testing.T) {
	tests := []struct {
		command  string
		expected bool
	}{
		{"exit 0", true},
		{"exit 1", false},
		{"exit 255", false},
		{"[ 1 -eq 1 ]", true},
		{"[ 1 -eq 2 ]", false},
		{"test -f /etc/hosts", true}, // This file typically exists
		{"test -f /nonexistent/file/xyz", false},
	}
	
	for _, tt := range tests {
		result, err := executeShellCondition(tt.command)
		
		if err != nil {
			t.Errorf("Command '%s' returned unexpected error: %v", tt.command, err)
			continue
		}
		
		if result != tt.expected {
			t.Errorf("Command '%s': expected %v, got %v", tt.command, tt.expected, result)
		}
	}
}

func TestExecuteShellCondition_ComplexCommand(t *testing.T) {
	// Test with pipes and redirects
	result, err := executeShellCondition("echo 'test' | grep test > /dev/null")
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !result {
		t.Error("Expected command with pipes to succeed")
	}
}

func TestTaskExecutionError(t *testing.T) {
	baseErr := errors.New("base error")
	taskErr := &TaskExecutionError{
		TaskName: "test-task",
		Err:      baseErr,
	}
	
	expectedMsg := "task 'test-task' failed: base error"
	if taskErr.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, taskErr.Error())
	}
}

func TestTaskExecutionError_WithNilError(t *testing.T) {
	taskErr := &TaskExecutionError{
		TaskName: "test-task",
		Err:      nil,
	}
	
	// Should not panic
	msg := taskErr.Error()
	if msg == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestResolveAgentAddress_WithPort(t *testing.T) {
	// Test with address that contains port
	address := "192.168.1.100:50053"
	
	result, err := resolveAgentAddress(address)
	
	if err != nil {
		t.Errorf("Expected no error for address with port, got %v", err)
	}
	
	if result != address {
		t.Errorf("Expected address to be returned as-is, got %s", result)
	}
}

func TestResolveAgentAddress_WithoutResolver(t *testing.T) {
	// Clear global resolver
	oldResolver := globalAgentResolver
	globalAgentResolver = nil
	defer func() { globalAgentResolver = oldResolver }()
	
	// Test with agent name (no port)
	agentName := "test-agent"
	
	_, err := resolveAgentAddress(agentName)
	
	if err == nil {
		t.Error("Expected error when resolver is not available")
	}
	
	expectedErrMsg := "no agent resolver available"
	if !contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// Mock resolver for testing
type mockAgentResolver struct {
	agents map[string]string
}

func (m *mockAgentResolver) GetAgentAddress(agentName string) (string, error) {
	if addr, ok := m.agents[agentName]; ok {
		return addr, nil
	}
	return "", errors.New("agent not found")
}

func TestResolveAgentAddress_WithResolver(t *testing.T) {
	// Setup mock resolver
	oldResolver := globalAgentResolver
	defer func() { globalAgentResolver = oldResolver }()
	
	mockResolver := &mockAgentResolver{
		agents: map[string]string{
			"test-agent": "192.168.1.100:50053",
			"prod-agent": "10.0.0.1:50053",
		},
	}
	globalAgentResolver = mockResolver
	
	tests := []struct {
		name      string
		input     string
		expected  string
		shouldErr bool
	}{
		{"Known agent", "test-agent", "192.168.1.100:50053", false},
		{"Another known agent", "prod-agent", "10.0.0.1:50053", false},
		{"Unknown agent", "unknown-agent", "", true},
		{"Direct address", "192.168.1.1:50053", "192.168.1.1:50053", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveAgentAddress(tt.input)
			
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			
			if !tt.shouldErr && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSetAgentResolver(t *testing.T) {
	oldResolver := globalAgentResolver
	defer func() { globalAgentResolver = oldResolver }()
	
	mockResolver := &mockAgentResolver{
		agents: make(map[string]string),
	}
	
	SetAgentResolver(mockResolver)
	
	if globalAgentResolver != mockResolver {
		t.Error("SetAgentResolver did not set the global resolver")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && 
		(s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		len(s) > len(substr)*2 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func BenchmarkExecuteShellCondition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		executeShellCondition("true")
	}
}

func BenchmarkExecuteShellCondition_WithPipe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		executeShellCondition("echo test | grep test > /dev/null")
	}
}
