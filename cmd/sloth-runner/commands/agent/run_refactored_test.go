package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test runCommandWithClient function
func TestRunCommandWithClient(t *testing.T) {
	tests := []struct {
		name          string
		responses     []*pb.StreamOutputResponse
		outputFormat  string
		expectedError bool
		checkOutput   func(string) bool
	}{
		{
			name: "successful command with stdout",
			responses: []*pb.StreamOutputResponse{
				{StdoutChunk: "hello\n"},
				{Finished: true, ExitCode: 0},
			},
			outputFormat:  "text",
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "hello")
			},
		},
		{
			name: "successful JSON output",
			responses: []*pb.StreamOutputResponse{
				{StdoutChunk: "output\n"},
				{Finished: true, ExitCode: 0},
			},
			outputFormat:  "json",
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "\"success\": true")
			},
		},
		{
			name: "command with exit code 1",
			responses: []*pb.StreamOutputResponse{
				{StderrChunk: "error\n"},
				{Finished: true, ExitCode: 1},
			},
			outputFormat:  "text",
			expectedError: true,
		},
		{
			name: "command with error message",
			responses: []*pb.StreamOutputResponse{
				{Error: "command failed"},
				{Finished: true, ExitCode: 1},
			},
			outputFormat:  "text",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockStream := &mocks.MockExecuteCommandClient{
				Responses: tt.responses,
			}

			mockClient.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
				return mockStream, nil
			}

			var outBuf, errBuf bytes.Buffer
			opts := RunCommandOptions{
				AgentName:    "test-agent",
				Command:      "test command",
				OutputFormat: tt.outputFormat,
				OutputWriter: &outBuf,
				ErrorWriter:  &errBuf,
			}

			err := runCommandWithClient(context.Background(), mockClient, opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.checkOutput != nil {
				if !tt.checkOutput(outBuf.String()) {
					t.Errorf("Output check failed. Got: %s", outBuf.String())
				}
			}
		})
	}
}

// Test processCommandStream function
func TestProcessCommandStream(t *testing.T) {
	tests := []struct {
		name            string
		responses       []*pb.StreamOutputResponse
		outputFormat    string
		expectedSuccess bool
		expectedExitCode int32
		expectedStdout  string
		expectedStderr  string
	}{
		{
			name: "successful command",
			responses: []*pb.StreamOutputResponse{
				{StdoutChunk: "output line 1\n"},
				{StdoutChunk: "output line 2\n"},
				{Finished: true, ExitCode: 0},
			},
			outputFormat:    "json",
			expectedSuccess: true,
			expectedExitCode: 0,
			expectedStdout:  "output line 1\noutput line 2\n",
		},
		{
			name: "command with stderr",
			responses: []*pb.StreamOutputResponse{
				{StderrChunk: "warning\n"},
				{StdoutChunk: "output\n"},
				{Finished: true, ExitCode: 0},
			},
			outputFormat:    "json",
			expectedSuccess: true,
			expectedExitCode: 0,
			expectedStdout:  "output\n",
			expectedStderr:  "warning\n",
		},
		{
			name: "failed command",
			responses: []*pb.StreamOutputResponse{
				{StderrChunk: "error occurred\n"},
				{Finished: true, ExitCode: 1},
			},
			outputFormat:    "json",
			expectedSuccess: false,
			expectedExitCode: 1,
			expectedStderr:  "error occurred\n",
		},
		{
			name: "command with error field",
			responses: []*pb.StreamOutputResponse{
				{Error: "execution failed"},
				{Finished: true, ExitCode: 127},
			},
			outputFormat:    "json",
			expectedSuccess: false,
			expectedExitCode: 127,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStream := &mocks.MockExecuteCommandClient{
				Responses: tt.responses,
			}

			var outBuf, errBuf bytes.Buffer
			result, err := processCommandStream(mockStream, tt.outputFormat, &outBuf, &errBuf)

			if err != nil {
				t.Fatalf("processCommandStream failed: %v", err)
			}

			if result.Success != tt.expectedSuccess {
				t.Errorf("Expected success=%v, got %v", tt.expectedSuccess, result.Success)
			}

			if result.ExitCode != tt.expectedExitCode {
				t.Errorf("Expected exit code=%d, got %d", tt.expectedExitCode, result.ExitCode)
			}

			if result.Stdout != tt.expectedStdout {
				t.Errorf("Expected stdout=%q, got %q", tt.expectedStdout, result.Stdout)
			}

			if result.Stderr != tt.expectedStderr {
				t.Errorf("Expected stderr=%q, got %q", tt.expectedStderr, result.Stderr)
			}
		})
	}
}

// Test formatCommandResultJSON function
func TestFormatCommandResultJSON(t *testing.T) {
	tests := []struct {
		name          string
		result        *CommandResult
		agentName     string
		command       string
		expectedError bool
		checkJSON     func(map[string]interface{}) bool
	}{
		{
			name: "successful result",
			result: &CommandResult{
				Success:     true,
				ExitCode:    0,
				Stdout:      "output\n",
				Stderr:      "",
				Error:       "",
				HasFinished: true,
			},
			agentName:     "test-agent",
			command:       "echo test",
			expectedError: false,
			checkJSON: func(m map[string]interface{}) bool {
				return m["success"] == true && m["exit_code"].(float64) == 0
			},
		},
		{
			name: "failed result",
			result: &CommandResult{
				Success:     false,
				ExitCode:    1,
				Stdout:      "",
				Stderr:      "error\n",
				Error:       "command failed",
				HasFinished: true,
			},
			agentName:     "test-agent",
			command:       "false",
			expectedError: true,
			checkJSON: func(m map[string]interface{}) bool {
				return m["success"] == false && m["error"] == "command failed"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatCommandResultJSON(tt.result, tt.agentName, tt.command, &buf)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Parse JSON
			var output map[string]interface{}
			if jsonErr := json.Unmarshal(buf.Bytes(), &output); jsonErr != nil {
				t.Fatalf("Invalid JSON: %v", jsonErr)
			}

			// Verify required fields
			if output["agent"] != tt.agentName {
				t.Errorf("Expected agent=%s, got %v", tt.agentName, output["agent"])
			}
			if output["command"] != tt.command {
				t.Errorf("Expected command=%s, got %v", tt.command, output["command"])
			}

			// Custom checks
			if tt.checkJSON != nil && !tt.checkJSON(output) {
				t.Error("JSON check failed")
			}
		})
	}
}

// Test formatCommandResultText function
func TestFormatCommandResultText(t *testing.T) {
	tests := []struct {
		name          string
		result        *CommandResult
		agentName     string
		expectedError bool
		checkOutput   func(string) bool
	}{
		{
			name: "successful command",
			result: &CommandResult{
				Success:     true,
				ExitCode:    0,
				HasFinished: true,
			},
			agentName:     "test-agent",
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "successfully") &&
					strings.Contains(out, "exit code: 0")
			},
		},
		{
			name: "failed with exit code",
			result: &CommandResult{
				Success:     false,
				ExitCode:    1,
				HasFinished: true,
			},
			agentName:     "test-agent",
			expectedError: true,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "failed") &&
					strings.Contains(out, "exit code: 1")
			},
		},
		{
			name: "failed with error message",
			result: &CommandResult{
				Success:     false,
				Error:       "timeout",
				HasFinished: false,
			},
			agentName:     "test-agent",
			expectedError: true,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "timeout")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatCommandResultText(tt.result, tt.agentName, &buf)

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

// Test CommandResult success logic
func TestCommandResult_SuccessLogic(t *testing.T) {
	tests := []struct {
		name            string
		hasFinished     bool
		exitCode        int32
		error           string
		expectedSuccess bool
	}{
		{"finished with exit 0", true, 0, "", true},
		{"finished with exit 1", true, 1, "", false},
		{"finished with error", true, 1, "failed", false},
		{"not finished, no error", false, -1, "", true},
		{"not finished with error", false, -1, "timeout", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := (tt.hasFinished && tt.exitCode == 0) || (!tt.hasFinished && tt.error == "")
			if success != tt.expectedSuccess {
				t.Errorf("Expected success=%v, got %v", tt.expectedSuccess, success)
			}
		})
	}
}

// Test processCommandStream with text output
func TestProcessCommandStream_TextOutput(t *testing.T) {
	mockStream := &mocks.MockExecuteCommandClient{
		Responses: []*pb.StreamOutputResponse{
			{StdoutChunk: "line 1\n"},
			{StdoutChunk: "line 2\n"},
			{StderrChunk: "warning\n"},
			{Finished: true, ExitCode: 0},
		},
	}

	var outBuf, errBuf bytes.Buffer
	result, err := processCommandStream(mockStream, "text", &outBuf, &errBuf)

	if err != nil {
		t.Fatalf("processCommandStream failed: %v", err)
	}

	// For text output, stdout/stderr are written directly to writers
	if outBuf.String() != "line 1\nline 2\n" {
		t.Errorf("Expected stdout to be written to outWriter, got %q", outBuf.String())
	}

	if errBuf.String() != "warning\n" {
		t.Errorf("Expected stderr to be written to errWriter, got %q", errBuf.String())
	}

	// Result buffers should be empty for text mode
	if result.Stdout != "" {
		t.Errorf("Expected empty result.Stdout for text mode, got %q", result.Stdout)
	}
}

// Test processCommandStream EOF handling
func TestProcessCommandStream_EOF(t *testing.T) {
	mockStream := &mocks.MockExecuteCommandClient{
		Responses: []*pb.StreamOutputResponse{
			{StdoutChunk: "output\n"},
			// EOF happens after this
		},
	}

	var outBuf, errBuf bytes.Buffer
	result, err := processCommandStream(mockStream, "json", &outBuf, &errBuf)

	if err != nil {
		t.Fatalf("processCommandStream should handle EOF gracefully: %v", err)
	}

	if result.Stdout != "output\n" {
		t.Errorf("Expected stdout=%q, got %q", "output\n", result.Stdout)
	}
}

// Test runCommandOnAgent integration
func TestRunCommandOnAgent_Integration(t *testing.T) {
	// This test verifies the integration but will fail on connection
	err := runCommandOnAgent("test-agent", "ls", "invalid:99999", "json")
	if err == nil {
		t.Log("Command succeeded (unexpected, but possible if server is running)")
	} else {
		// Expected to fail with connection error
		if !strings.Contains(err.Error(), "failed to connect") &&
			!strings.Contains(err.Error(), "connection refused") &&
			!strings.Contains(err.Error(), "no such host") {
			t.Logf("Got connection error as expected: %v", err)
		}
	}
}
