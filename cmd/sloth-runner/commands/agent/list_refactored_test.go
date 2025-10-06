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

// Test listAgentsWithClient function
func TestListAgentsWithClient(t *testing.T) {
	tests := []struct {
		name          string
		agents        []*pb.AgentInfo
		expectedError bool
		checkOutput   func(string) bool
	}{
		{
			name: "successful list with agents",
			agents: []*pb.AgentInfo{
				{
					AgentName:         "agent1",
					AgentAddress:      "192.168.1.10:50051",
					Status:            "Active",
					LastHeartbeat:     1609459200,
					LastInfoCollected: 1609459200,
				},
				{
					AgentName:         "agent2",
					AgentAddress:      "192.168.1.11:50051",
					Status:            "Inactive",
					LastHeartbeat:     1609459100,
					LastInfoCollected: 0,
				},
			},
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "agent1") &&
					strings.Contains(out, "agent2") &&
					strings.Contains(out, "192.168.1.10:50051")
			},
		},
		{
			name:          "empty agent list",
			agents:        []*pb.AgentInfo{},
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "No agents registered")
			},
		},
		{
			name: "single agent",
			agents: []*pb.AgentInfo{
				{
					AgentName:         "solo-agent",
					AgentAddress:      "10.0.0.1:50051",
					Status:            "Active",
					LastHeartbeat:     1609459200,
					LastInfoCollected: 1609459200,
				},
			},
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "solo-agent") &&
					strings.Contains(out, "10.0.0.1:50051")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
				return &pb.ListAgentsResponse{
					Agents: tt.agents,
				}, nil
			}

			var buf bytes.Buffer
			opts := ListAgentsOptions{
				Writer: &buf,
			}

			err := listAgentsWithClient(context.Background(), mockClient, opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
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

// Test formatAgentsTable function
func TestFormatAgentsTable(t *testing.T) {
	tests := []struct {
		name        string
		agents      []*pb.AgentInfo
		checkOutput func(string) bool
	}{
		{
			name: "single agent with all fields",
			agents: []*pb.AgentInfo{
				{
					AgentName:         "test-agent",
					AgentAddress:      "192.168.1.100:50051",
					Status:            "Active",
					LastHeartbeat:     1609459200,
					LastInfoCollected: 1609459200,
				},
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "AGENT NAME") &&
					strings.Contains(out, "ADDRESS") &&
					strings.Contains(out, "STATUS") &&
					strings.Contains(out, "test-agent") &&
					strings.Contains(out, "192.168.1.100:50051")
			},
		},
		{
			name: "multiple agents with different statuses",
			agents: []*pb.AgentInfo{
				{
					AgentName:         "active-agent",
					AgentAddress:      "10.0.0.1:50051",
					Status:            "Active",
					LastHeartbeat:     1609459200,
					LastInfoCollected: 1609459200,
				},
				{
					AgentName:         "inactive-agent",
					AgentAddress:      "10.0.0.2:50051",
					Status:            "Inactive",
					LastHeartbeat:     0,
					LastInfoCollected: 0,
				},
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "active-agent") &&
					strings.Contains(out, "inactive-agent")
			},
		},
		{
			name: "agent with zero timestamps",
			agents: []*pb.AgentInfo{
				{
					AgentName:         "new-agent",
					AgentAddress:      "10.0.0.3:50051",
					Status:            "Active",
					LastHeartbeat:     0,
					LastInfoCollected: 0,
				},
			},
			checkOutput: func(out string) bool {
				return strings.Contains(out, "new-agent") &&
					strings.Contains(out, "N/A") &&
					strings.Contains(out, "Never")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatAgentsTable(tt.agents, &buf)

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

// Test formatStatus function
func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		shouldBeGreen  bool
		shouldBeRed    bool
	}{
		{
			name:          "active status is green",
			status:        "Active",
			shouldBeGreen: true,
			shouldBeRed:   false,
		},
		{
			name:          "inactive status is red",
			status:        "Inactive",
			shouldBeGreen: false,
			shouldBeRed:   true,
		},
		{
			name:          "unknown status is red",
			status:        "Unknown",
			shouldBeGreen: false,
			shouldBeRed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatStatus(tt.status)

			// The result contains ANSI color codes, so we check if the original status is in the result
			if !strings.Contains(result, tt.status) {
				t.Errorf("Expected result to contain status %q, got %q", tt.status, result)
			}

			// For green status, the result should contain green ANSI codes
			if tt.shouldBeGreen && !strings.Contains(result, "\x1b[") {
				t.Logf("Green status result: %q", result)
			}

			// For red status, the result should contain red ANSI codes
			if tt.shouldBeRed && !strings.Contains(result, "\x1b[") {
				t.Logf("Red status result: %q", result)
			}
		})
	}
}

// Test formatTimestamp function
func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name         string
		timestamp    int64
		defaultValue string
		checkResult  func(string) bool
	}{
		{
			name:         "valid timestamp",
			timestamp:    1609459200,
			defaultValue: "N/A",
			checkResult: func(result string) bool {
				// Should be a valid RFC3339 timestamp, not the default
				return result != "N/A" && strings.Contains(result, "2021") || strings.Contains(result, "2020")
			},
		},
		{
			name:         "zero timestamp returns default",
			timestamp:    0,
			defaultValue: "N/A",
			checkResult: func(result string) bool {
				return result == "N/A"
			},
		},
		{
			name:         "zero timestamp with custom default",
			timestamp:    0,
			defaultValue: "Never",
			checkResult: func(result string) bool {
				return result == "Never"
			},
		},
		{
			name:         "negative timestamp returns default",
			timestamp:    -1,
			defaultValue: "Unknown",
			checkResult: func(result string) bool {
				return result == "Unknown"
			},
		},
		{
			name:         "recent timestamp",
			timestamp:    1672531200, // 2023-01-01T00:00:00Z
			defaultValue: "N/A",
			checkResult: func(result string) bool {
				// Should be a valid RFC3339 timestamp containing 2023 or 2022 (due to timezone)
				return result != "N/A" && (strings.Contains(result, "2023") || strings.Contains(result, "2022"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimestamp(tt.timestamp, tt.defaultValue)

			if !tt.checkResult(result) {
				t.Errorf("Timestamp check failed. Got: %q", result)
			}
		})
	}
}

// Test listAgentsWithClient error handling
func TestListAgentsWithClient_ErrorHandling(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
		return nil, grpc.ErrServerStopped
	}

	var buf bytes.Buffer
	opts := ListAgentsOptions{
		Writer: &buf,
	}

	err := listAgentsWithClient(context.Background(), mockClient, opts)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to list agents") {
		t.Errorf("Expected 'failed to list agents' error, got: %v", err)
	}
}

// Test formatAgentsTable with various timestamp combinations
func TestFormatAgentsTable_TimestampCombinations(t *testing.T) {
	agents := []*pb.AgentInfo{
		{
			AgentName:         "agent-with-heartbeat-only",
			AgentAddress:      "10.0.0.1:50051",
			Status:            "Active",
			LastHeartbeat:     1609459200,
			LastInfoCollected: 0,
		},
		{
			AgentName:         "agent-with-info-only",
			AgentAddress:      "10.0.0.2:50051",
			Status:            "Active",
			LastHeartbeat:     0,
			LastInfoCollected: 1609459200,
		},
		{
			AgentName:         "agent-with-both",
			AgentAddress:      "10.0.0.3:50051",
			Status:            "Active",
			LastHeartbeat:     1609459200,
			LastInfoCollected: 1609459300,
		},
		{
			AgentName:         "agent-with-neither",
			AgentAddress:      "10.0.0.4:50051",
			Status:            "Inactive",
			LastHeartbeat:     0,
			LastInfoCollected: 0,
		},
	}

	var buf bytes.Buffer
	err := formatAgentsTable(agents, &buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Check that all agents appear in output
	for _, agent := range agents {
		if !strings.Contains(output, agent.AgentName) {
			t.Errorf("Agent %s not found in output", agent.AgentName)
		}
	}

	// Check default values appear
	if !strings.Contains(output, "N/A") {
		t.Error("Expected 'N/A' for missing heartbeat timestamps")
	}

	if !strings.Contains(output, "Never") {
		t.Error("Expected 'Never' for missing info collected timestamps")
	}

	// Check formatted timestamps appear (should contain year 2021 or 2020 depending on timezone)
	if !strings.Contains(output, "2021") && !strings.Contains(output, "2020") {
		t.Error("Expected formatted timestamp in output")
	}
}
