package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test stopAgentWithClient function
func TestStopAgentWithClient(t *testing.T) {
	tests := []struct {
		name          string
		agentName     string
		response      *pb.StopAgentResponse
		responseError error
		expectedError bool
	}{
		{
			name:      "successful stop",
			agentName: "test-agent",
			response: &pb.StopAgentResponse{
				Success: true,
				Message: "Agent stopped",
			},
			responseError: nil,
			expectedError: false,
		},
		{
			name:      "agent not found",
			agentName: "missing-agent",
			response: &pb.StopAgentResponse{
				Success: false,
				Message: "Agent not found",
			},
			responseError: nil,
			expectedError: true,
		},
		{
			name:          "grpc error",
			agentName:     "test-agent",
			response:      nil,
			responseError: grpc.ErrServerStopped,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockClient.StopAgentFunc = func(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error) {
				if in.AgentName != tt.agentName {
					t.Errorf("Expected agent name %q, got %q", tt.agentName, in.AgentName)
				}
				return tt.response, tt.responseError
			}

			opts := StopAgentOptions{
				AgentName: tt.agentName,
			}

			err := stopAgentWithClient(context.Background(), mockClient, opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test stopAgentWithClient with different error messages
func TestStopAgentWithClient_ErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		response       *pb.StopAgentResponse
		responseError  error
		expectedErrMsg string
	}{
		{
			name:           "grpc connection error",
			response:       nil,
			responseError:  grpc.ErrServerStopped,
			expectedErrMsg: "failed to stop agent",
		},
		{
			name: "unsuccessful stop",
			response: &pb.StopAgentResponse{
				Success: false,
				Message: "Agent is already stopped",
			},
			responseError:  nil,
			expectedErrMsg: "stop failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockClient.StopAgentFunc = func(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error) {
				return tt.response, tt.responseError
			}

			opts := StopAgentOptions{
				AgentName: "test-agent",
			}

			err := stopAgentWithClient(context.Background(), mockClient, opts)

			if err == nil {
				t.Error("Expected error but got none")
			}

			if !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("Expected error containing %q, got: %v", tt.expectedErrMsg, err)
			}
		})
	}
}

// Test stopAgentWithClient verifies correct request parameters
func TestStopAgentWithClient_RequestParameters(t *testing.T) {
	testAgentName := "my-test-agent"
	mockClient := mocks.NewMockAgentRegistryClient()

	var capturedRequest *pb.StopAgentRequest
	mockClient.StopAgentFunc = func(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error) {
		capturedRequest = in
		return &pb.StopAgentResponse{Success: true, Message: "Stopped"}, nil
	}

	opts := StopAgentOptions{
		AgentName: testAgentName,
	}

	err := stopAgentWithClient(context.Background(), mockClient, opts)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if capturedRequest == nil {
		t.Fatal("Request was not captured")
	}

	if capturedRequest.AgentName != testAgentName {
		t.Errorf("Expected agent name %q in request, got %q", testAgentName, capturedRequest.AgentName)
	}
}

// Test stopAgentWithClient with context cancellation
func TestStopAgentWithClient_ContextCancellation(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.StopAgentFunc = func(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error) {
		return nil, context.Canceled
	}

	opts := StopAgentOptions{
		AgentName: "test-agent",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := stopAgentWithClient(ctx, mockClient, opts)

	if err == nil {
		t.Error("Expected error due to context cancellation")
	}

	if !strings.Contains(err.Error(), "failed to stop agent") {
		t.Errorf("Expected error about failed stop, got: %v", err)
	}
}

// Test stopAgentWithClient with different failure messages
func TestStopAgentWithClient_FailureMessages(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "agent already stopped",
			message: "Agent is already in stopped state",
		},
		{
			name:    "agent not responding",
			message: "Agent is not responding to shutdown request",
		},
		{
			name:    "internal error",
			message: "Internal server error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockClient.StopAgentFunc = func(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error) {
				return &pb.StopAgentResponse{
					Success: false,
					Message: tt.message,
				}, nil
			}

			opts := StopAgentOptions{
				AgentName: "test-agent",
			}

			err := stopAgentWithClient(context.Background(), mockClient, opts)

			if err == nil {
				t.Error("Expected error but got none")
			}

			if !strings.Contains(err.Error(), tt.message) {
				t.Errorf("Expected error to contain %q, got: %v", tt.message, err)
			}
		})
	}
}
