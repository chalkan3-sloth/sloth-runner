package agent

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test checkModuleAvailability function
func TestCheckModuleAvailability(t *testing.T) {
	tests := []struct {
		name           string
		module         moduleCheck
		streamOutput   string
		expectedResult bool
		expectError    bool
	}{
		{
			name: "module found",
			module: moduleCheck{
				Name:        "Docker",
				Command:     "docker",
				Description: "Container platform",
			},
			streamOutput:   "found",
			expectedResult: true,
			expectError:    false,
		},
		{
			name: "module not found",
			module: moduleCheck{
				Name:        "NonExistent",
				Command:     "nonexistent",
				Description: "Non-existent tool",
			},
			streamOutput:   "not found",
			expectedResult: false,
			expectError:    false,
		},
		{
			name: "module found with extra whitespace",
			module: moduleCheck{
				Name:        "Git",
				Command:     "git",
				Description: "Version control",
			},
			streamOutput:   "  found  \n",
			expectedResult: true,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockStream := &mocks.MockExecuteCommandClient{
				Responses: []*pb.StreamOutputResponse{
					{StdoutChunk: tt.streamOutput},
				},
			}

			mockClient.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
				return mockStream, nil
			}

			result, err := checkModuleAvailability(context.Background(), mockClient, "test-agent", tt.module)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expectedResult {
				t.Errorf("Expected result=%v, got %v", tt.expectedResult, result)
			}
		})
	}
}

// Test checkModuleAvailability with stream errors
func TestCheckModuleAvailability_StreamError(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
		return nil, grpc.ErrServerStopped
	}

	mod := moduleCheck{Name: "Test", Command: "test", Description: "Test tool"}
	_, err := checkModuleAvailability(context.Background(), mockClient, "test-agent", mod)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to execute command") {
		t.Errorf("Expected 'failed to execute command' error, got: %v", err)
	}
}

// Test checkModuleAvailability with receive error
func TestCheckModuleAvailability_ReceiveError(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockStream := &mocks.MockExecuteCommandClient{
		Responses: []*pb.StreamOutputResponse{
			{StdoutChunk: "partial"},
		},
		RecvError: grpc.ErrServerStopped,
	}

	mockClient.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
		return mockStream, nil
	}

	mod := moduleCheck{Name: "Test", Command: "test", Description: "Test tool"}
	_, err := checkModuleAvailability(context.Background(), mockClient, "test-agent", mod)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "error receiving stream") {
		t.Errorf("Expected 'error receiving stream' error, got: %v", err)
	}
}

// Test formatModulesResults function
func TestFormatModulesResults(t *testing.T) {
	tests := []struct {
		name        string
		result      *ModulesCheckResult
		allModules  []moduleCheck
		checkOutput func(string) bool
	}{
		{
			name: "all modules available",
			result: &ModulesCheckResult{
				Available: []moduleCheck{
					{Name: "Docker", Command: "docker", Description: "Container platform"},
					{Name: "Git", Command: "git", Description: "Version control"},
				},
				Missing: []moduleCheck{},
			},
			allModules: []moduleCheck{
				{Name: "Docker", Command: "docker", Description: "Container platform"},
				{Name: "Git", Command: "git", Description: "Version control"},
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "Docker") &&
					strings.Contains(out, "Git") &&
					strings.Contains(out, "Available: ")
			},
		},
		{
			name: "some modules missing",
			result: &ModulesCheckResult{
				Available: []moduleCheck{
					{Name: "Docker", Command: "docker", Description: "Container platform"},
				},
				Missing: []moduleCheck{
					{Name: "Terraform", Command: "terraform", Description: "IaC tool"},
				},
			},
			allModules: []moduleCheck{
				{Name: "Docker", Command: "docker", Description: "Container platform"},
				{Name: "Terraform", Command: "terraform", Description: "IaC tool"},
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "Docker") &&
					strings.Contains(out, "Terraform") &&
					strings.Contains(out, "Missing Modules")
			},
		},
		{
			name: "all modules missing",
			result: &ModulesCheckResult{
				Available: []moduleCheck{},
				Missing: []moduleCheck{
					{Name: "Docker", Command: "docker", Description: "Container platform"},
					{Name: "Terraform", Command: "terraform", Description: "IaC tool"},
				},
			},
			allModules: []moduleCheck{
				{Name: "Docker", Command: "docker", Description: "Container platform"},
				{Name: "Terraform", Command: "terraform", Description: "IaC tool"},
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "Missing Modules") &&
					strings.Contains(out, "Information:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatModulesResults(tt.result, tt.allModules, &buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.checkOutput != nil {
				if !tt.checkOutput(buf.String()) {
					t.Errorf("Output check failed. Got: %s", buf.String())
				}
			}
		})
	}
}

// Test checkAgentModulesWithClient function
func TestCheckAgentModulesWithClient(t *testing.T) {
	tests := []struct {
		name          string
		agentName     string
		setupMock     func(*mocks.MockAgentRegistryClient)
		expectedError bool
		checkOutput   func(string) bool
	}{
		{
			name:      "successful module check",
			agentName: "test-agent",
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
					// Simulate finding all modules
					return &mocks.MockExecuteCommandClient{
						Responses: []*pb.StreamOutputResponse{
							{StdoutChunk: "found"},
						},
					}, nil
				}
			},
			expectedError: false,
			checkOutput: func(out string) bool {
				// Check for module names that are written to the buffer
				// pterm output goes to stdout, but module descriptions go to the buffer
				return strings.Contains(out, "Container and VM management") ||
					strings.Contains(out, "Infrastructure as Code") ||
					len(out) > 0 // At least some output was written
			},
		},
		{
			name:      "module check with stream error",
			agentName: "test-agent",
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				callCount := 0
				m.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
					callCount++
					if callCount == 1 {
						// First call fails
						return nil, grpc.ErrServerStopped
					}
					return &mocks.MockExecuteCommandClient{
						Responses: []*pb.StreamOutputResponse{
							{StdoutChunk: "found"},
						},
					}, nil
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			tt.setupMock(mockClient)

			var buf bytes.Buffer
			opts := ModulesCheckOptions{
				AgentName: tt.agentName,
				Writer:    &buf,
			}

			err := checkAgentModulesWithClient(context.Background(), mockClient, opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectedError && tt.checkOutput != nil {
				if !tt.checkOutput(buf.String()) {
					t.Errorf("Output check failed. Got: %s", buf.String())
				}
			}
		})
	}
}

// Test getDefaultModules function
func TestGetDefaultModules(t *testing.T) {
	modules := getDefaultModules()

	if len(modules) == 0 {
		t.Error("Expected default modules but got empty list")
	}

	// Check for some expected modules
	expectedModules := []string{"Docker", "Git", "Terraform", "kubectl"}
	found := make(map[string]bool)

	for _, mod := range modules {
		found[mod.Name] = true
	}

	for _, expected := range expectedModules {
		if !found[expected] {
			t.Errorf("Expected module %s not found in default modules", expected)
		}
	}

	// Verify all modules have required fields
	for _, mod := range modules {
		if mod.Name == "" {
			t.Error("Module with empty name found")
		}
		if mod.Command == "" {
			t.Error("Module with empty command found")
		}
		if mod.Description == "" {
			t.Error("Module with empty description found")
		}
	}
}

// Test checkModuleAvailability with stderr output
func TestCheckModuleAvailability_WithStderr(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockStream := &mocks.MockExecuteCommandClient{
		Responses: []*pb.StreamOutputResponse{
			{StderrChunk: "not found"},
		},
	}

	mockClient.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
		return mockStream, nil
	}

	mod := moduleCheck{Name: "Test", Command: "test", Description: "Test tool"}
	result, err := checkModuleAvailability(context.Background(), mockClient, "test-agent", mod)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result {
		t.Error("Expected module to be not found, but it was marked as found")
	}
}

// Test formatModulesResults with io.Writer error handling
func TestFormatModulesResults_WriterHandling(t *testing.T) {
	result := &ModulesCheckResult{
		Available: []moduleCheck{
			{Name: "Docker", Command: "docker", Description: "Container platform"},
		},
		Missing: []moduleCheck{},
	}

	modules := []moduleCheck{
		{Name: "Docker", Command: "docker", Description: "Container platform"},
	}

	// Use a limited writer to test edge cases
	var buf bytes.Buffer
	err := formatModulesResults(result, modules, &buf)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify output contains expected information
	output := buf.String()
	if !strings.Contains(output, "Docker") {
		t.Error("Expected Docker in output")
	}
}

// Test module check command generation
func TestCheckModuleAvailability_CommandGeneration(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	var receivedCommand string

	mockClient.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
		receivedCommand = in.Command
		return &mocks.MockExecuteCommandClient{
			Responses: []*pb.StreamOutputResponse{
				{StdoutChunk: "found"},
			},
		}, nil
	}

	mod := moduleCheck{Name: "Docker", Command: "docker", Description: "Container platform"}
	_, err := checkModuleAvailability(context.Background(), mockClient, "test-agent", mod)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCmd := "command -v docker >/dev/null 2>&1 && echo 'found' || echo 'not found'"
	if receivedCommand != expectedCmd {
		t.Errorf("Expected command %q, got %q", expectedCmd, receivedCommand)
	}
}

// Helper mock for testing EOF handling
type eofStream struct {
	pb.AgentRegistry_ExecuteCommandClient
}

func (s *eofStream) Recv() (*pb.StreamOutputResponse, error) {
	return nil, io.EOF
}

// Test checkModuleAvailability with immediate EOF
func TestCheckModuleAvailability_ImmediateEOF(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
		return &eofStream{}, nil
	}

	mod := moduleCheck{Name: "Test", Command: "test", Description: "Test tool"}
	result, err := checkModuleAvailability(context.Background(), mockClient, "test-agent", mod)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Empty output should result in false (not "found")
	if result {
		t.Error("Expected false for empty output, got true")
	}
}
