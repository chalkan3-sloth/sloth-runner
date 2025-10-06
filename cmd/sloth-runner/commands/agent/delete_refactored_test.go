package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test deleteAgentWithClient function
func TestDeleteAgentWithClient(t *testing.T) {
	tests := []struct {
		name          string
		agentName     string
		response      *pb.UnregisterAgentResponse
		responseError error
		expectedError bool
	}{
		{
			name:      "successful deletion",
			agentName: "test-agent",
			response: &pb.UnregisterAgentResponse{
				Success: true,
			},
			responseError: nil,
			expectedError: false,
		},
		{
			name:      "agent not found",
			agentName: "missing-agent",
			response: &pb.UnregisterAgentResponse{
				Success: false,
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
			mockClient.UnregisterAgentFunc = func(ctx context.Context, in *pb.UnregisterAgentRequest, opts ...grpc.CallOption) (*pb.UnregisterAgentResponse, error) {
				if in.AgentName != tt.agentName {
					t.Errorf("Expected agent name %q, got %q", tt.agentName, in.AgentName)
				}
				return tt.response, tt.responseError
			}

			opts := DeleteAgentOptions{
				AgentName: tt.agentName,
			}

			err := deleteAgentWithClient(context.Background(), mockClient, opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test deleteAgentWithClient with different error messages
func TestDeleteAgentWithClient_ErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		response       *pb.UnregisterAgentResponse
		responseError  error
		expectedErrMsg string
	}{
		{
			name:           "grpc connection error",
			response:       nil,
			responseError:  grpc.ErrServerStopped,
			expectedErrMsg: "failed to delete agent",
		},
		{
			name: "unsuccessful deletion",
			response: &pb.UnregisterAgentResponse{
				Success: false,
			},
			responseError:  nil,
			expectedErrMsg: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockClient.UnregisterAgentFunc = func(ctx context.Context, in *pb.UnregisterAgentRequest, opts ...grpc.CallOption) (*pb.UnregisterAgentResponse, error) {
				return tt.response, tt.responseError
			}

			opts := DeleteAgentOptions{
				AgentName: "test-agent",
			}

			err := deleteAgentWithClient(context.Background(), mockClient, opts)

			if err == nil {
				t.Error("Expected error but got none")
			}

			if !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("Expected error containing %q, got: %v", tt.expectedErrMsg, err)
			}
		})
	}
}

// Test deleteAgentWithClient verifies correct request parameters
func TestDeleteAgentWithClient_RequestParameters(t *testing.T) {
	testAgentName := "my-test-agent"
	mockClient := mocks.NewMockAgentRegistryClient()

	var capturedRequest *pb.UnregisterAgentRequest
	mockClient.UnregisterAgentFunc = func(ctx context.Context, in *pb.UnregisterAgentRequest, opts ...grpc.CallOption) (*pb.UnregisterAgentResponse, error) {
		capturedRequest = in
		return &pb.UnregisterAgentResponse{Success: true}, nil
	}

	opts := DeleteAgentOptions{
		AgentName: testAgentName,
	}

	err := deleteAgentWithClient(context.Background(), mockClient, opts)

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

// Test deleteAgentWithClient with context cancellation
func TestDeleteAgentWithClient_ContextCancellation(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.UnregisterAgentFunc = func(ctx context.Context, in *pb.UnregisterAgentRequest, opts ...grpc.CallOption) (*pb.UnregisterAgentResponse, error) {
		return nil, context.Canceled
	}

	opts := DeleteAgentOptions{
		AgentName: "test-agent",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := deleteAgentWithClient(ctx, mockClient, opts)

	if err == nil {
		t.Error("Expected error due to context cancellation")
	}

	if !strings.Contains(err.Error(), "failed to delete agent") {
		t.Errorf("Expected error about failed deletion, got: %v", err)
	}
}
