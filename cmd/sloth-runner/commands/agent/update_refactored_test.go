package agent

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test findAgentAddress function
func TestFindAgentAddress(t *testing.T) {
	tests := []struct {
		name           string
		agentName      string
		agents         []*pb.AgentInfo
		expectedAddr   string
		expectedError  bool
	}{
		{
			name:      "agent found",
			agentName: "test-agent",
			agents: []*pb.AgentInfo{
				{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
				{AgentName: "other-agent", AgentAddress: "192.168.1.11:50051"},
			},
			expectedAddr:  "192.168.1.10:50051",
			expectedError: false,
		},
		{
			name:      "agent not found",
			agentName: "missing-agent",
			agents: []*pb.AgentInfo{
				{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
			},
			expectedAddr:  "",
			expectedError: true,
		},
		{
			name:          "empty agent list",
			agentName:     "any-agent",
			agents:        []*pb.AgentInfo{},
			expectedAddr:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
				return &pb.ListAgentsResponse{Agents: tt.agents}, nil
			}

			addr, err := findAgentAddress(context.Background(), mockClient, tt.agentName)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if addr != tt.expectedAddr {
				t.Errorf("Expected address %q, got %q", tt.expectedAddr, addr)
			}
		})
	}
}

// Test findAgentAddress with ListAgents error
func TestFindAgentAddress_ListError(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
		return nil, grpc.ErrServerStopped
	}

	_, err := findAgentAddress(context.Background(), mockClient, "test-agent")

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to list agents") {
		t.Errorf("Expected 'failed to list agents' error, got: %v", err)
	}
}

// Test performAgentUpdate function
func TestPerformAgentUpdate(t *testing.T) {
	tests := []struct {
		name          string
		opts          UpdateAgentOptions
		response      *pb.UpdateAgentResponse
		expectedError bool
		checkResult   func(*UpdateAgentResult) bool
	}{
		{
			name: "successful update",
			opts: UpdateAgentOptions{
				AgentName:     "test-agent",
				TargetVersion: "1.2.0",
				Restart:       true,
			},
			response: &pb.UpdateAgentResponse{
				Success:    true,
				Message:    "Update successful",
				OldVersion: "1.0.0",
				NewVersion: "1.2.0",
			},
			expectedError: false,
			checkResult: func(r *UpdateAgentResult) bool {
				return r.Success && r.OldVersion == "1.0.0" && r.NewVersion == "1.2.0"
			},
		},
		{
			name: "update with skip restart",
			opts: UpdateAgentOptions{
				AgentName:     "test-agent",
				TargetVersion: "latest",
				Restart:       false,
			},
			response: &pb.UpdateAgentResponse{
				Success:    true,
				Message:    "Update successful, restart skipped",
				OldVersion: "1.1.0",
				NewVersion: "1.2.0",
			},
			expectedError: false,
			checkResult: func(r *UpdateAgentResult) bool {
				return r.Success && r.Message == "Update successful, restart skipped"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentClient()
			mockClient.UpdateAgentFunc = func(ctx context.Context, in *pb.UpdateAgentRequest, opts ...grpc.CallOption) (*pb.UpdateAgentResponse, error) {
				// Verify request parameters
				if in.TargetVersion != tt.opts.TargetVersion {
					t.Errorf("Expected version %q, got %q", tt.opts.TargetVersion, in.TargetVersion)
				}
				if in.SkipRestart != !tt.opts.Restart {
					t.Errorf("Expected SkipRestart=%v, got %v", !tt.opts.Restart, in.SkipRestart)
				}
				return tt.response, nil
			}

			result, err := performAgentUpdate(context.Background(), mockClient, tt.opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectedError && tt.checkResult != nil {
				if !tt.checkResult(result) {
					t.Errorf("Result check failed: %+v", result)
				}
			}
		})
	}
}

// Test performAgentUpdate with error
func TestPerformAgentUpdate_Error(t *testing.T) {
	mockClient := mocks.NewMockAgentClient()
	mockClient.UpdateAgentFunc = func(ctx context.Context, in *pb.UpdateAgentRequest, opts ...grpc.CallOption) (*pb.UpdateAgentResponse, error) {
		return nil, grpc.ErrServerStopped
	}

	opts := UpdateAgentOptions{
		AgentName:     "test-agent",
		TargetVersion: "1.2.0",
		Restart:       true,
	}

	_, err := performAgentUpdate(context.Background(), mockClient, opts)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to update agent") {
		t.Errorf("Expected 'failed to update agent' error, got: %v", err)
	}
}

// Test displayUpdateSummary function
func TestDisplayUpdateSummary(t *testing.T) {
	tests := []struct {
		name        string
		result      *UpdateAgentResult
		opts        UpdateAgentOptions
		checkOutput func(string) bool
	}{
		{
			name: "update with restart",
			result: &UpdateAgentResult{
				Success:    true,
				OldVersion: "1.0.0",
				NewVersion: "1.2.0",
				Message:    "Update completed",
			},
			opts: UpdateAgentOptions{
				AgentName: "test-agent",
				Restart:   true,
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "test-agent") &&
					strings.Contains(out, "1.0.0") &&
					strings.Contains(out, "1.2.0") &&
					strings.Contains(out, "Update completed")
			},
		},
		{
			name: "update without restart",
			result: &UpdateAgentResult{
				Success:    true,
				OldVersion: "1.1.0",
				NewVersion: "1.3.0",
				Message:    "",
			},
			opts: UpdateAgentOptions{
				AgentName: "another-agent",
				Restart:   false,
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "another-agent") &&
					strings.Contains(out, "1.1.0") &&
					strings.Contains(out, "restart skipped")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			displayUpdateSummary(tt.result, tt.opts, &buf)

			if tt.checkOutput != nil {
				if !tt.checkOutput(buf.String()) {
					t.Errorf("Output check failed. Got: %s", buf.String())
				}
			}
		})
	}
}

// Test updateAgentWithClients function
func TestUpdateAgentWithClients(t *testing.T) {
	tests := []struct {
		name           string
		opts           UpdateAgentOptions
		setupMocks     func(*mocks.MockAgentRegistryClient, *mocks.MockAgentClient)
		expectedError  bool
		checkResult    func(*UpdateAgentResult) bool
	}{
		{
			name: "successful update",
			opts: UpdateAgentOptions{
				AgentName:     "test-agent",
				TargetVersion: "1.2.0",
				Restart:       true,
			},
			setupMocks: func(regClient *mocks.MockAgentRegistryClient, agentClient *mocks.MockAgentClient) {
				regClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{
							{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
						},
					}, nil
				}
				agentClient.UpdateAgentFunc = func(ctx context.Context, in *pb.UpdateAgentRequest, opts ...grpc.CallOption) (*pb.UpdateAgentResponse, error) {
					return &pb.UpdateAgentResponse{
						Success:    true,
						Message:    "Update successful",
						OldVersion: "1.0.0",
						NewVersion: "1.2.0",
					}, nil
				}
			},
			expectedError: false,
			checkResult: func(r *UpdateAgentResult) bool {
				return r.Success && r.NewVersion == "1.2.0"
			},
		},
		{
			name: "agent not found",
			opts: UpdateAgentOptions{
				AgentName:     "missing-agent",
				TargetVersion: "1.2.0",
				Restart:       true,
			},
			setupMocks: func(regClient *mocks.MockAgentRegistryClient, agentClient *mocks.MockAgentClient) {
				regClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{},
					}, nil
				}
			},
			expectedError: true,
		},
		{
			name: "update failed",
			opts: UpdateAgentOptions{
				AgentName:     "test-agent",
				TargetVersion: "1.2.0",
				Restart:       true,
			},
			setupMocks: func(regClient *mocks.MockAgentRegistryClient, agentClient *mocks.MockAgentClient) {
				regClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{
							{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
						},
					}, nil
				}
				agentClient.UpdateAgentFunc = func(ctx context.Context, in *pb.UpdateAgentRequest, opts ...grpc.CallOption) (*pb.UpdateAgentResponse, error) {
					return &pb.UpdateAgentResponse{
						Success: false,
						Message: "Update failed - version not found",
					}, nil
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegClient := mocks.NewMockAgentRegistryClient()
			mockAgentClient := mocks.NewMockAgentClient()

			tt.setupMocks(mockRegClient, mockAgentClient)

			agentClientFactory := func(addr string) (AgentClient, func(), error) {
				return mockAgentClient, func() {}, nil
			}

			var buf bytes.Buffer
			tt.opts.Writer = &buf

			result, err := updateAgentWithClients(context.Background(), mockRegClient, agentClientFactory, tt.opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectedError && tt.checkResult != nil {
				if !tt.checkResult(result) {
					t.Errorf("Result check failed: %+v", result)
				}
			}
		})
	}
}

// Test updateAgentWithClients with agent connection error
func TestUpdateAgentWithClients_ConnectionError(t *testing.T) {
	mockRegClient := mocks.NewMockAgentRegistryClient()
	mockRegClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
		return &pb.ListAgentsResponse{
			Agents: []*pb.AgentInfo{
				{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
			},
		}, nil
	}

	agentClientFactory := func(addr string) (AgentClient, func(), error) {
		return nil, nil, grpc.ErrServerStopped
	}

	var buf bytes.Buffer
	opts := UpdateAgentOptions{
		AgentName:     "test-agent",
		TargetVersion: "1.2.0",
		Restart:       true,
		Writer:        &buf,
	}

	_, err := updateAgentWithClients(context.Background(), mockRegClient, agentClientFactory, opts)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to connect to agent") {
		t.Errorf("Expected 'failed to connect to agent' error, got: %v", err)
	}
}
