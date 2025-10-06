package agent

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test checkExistingAgent function
func TestCheckExistingAgent(t *testing.T) {
	tests := []struct {
		name          string
		agentName     string
		setupPIDFile  bool
		pidContent    string
		createProcess bool
		expectRunning bool
		expectError   bool
	}{
		{
			name:          "no existing agent",
			agentName:     "test-agent-nonexistent",
			setupPIDFile:  false,
			expectRunning: false,
			expectError:   false,
		},
		{
			name:          "agent with invalid PID file",
			agentName:     "test-agent-invalid",
			setupPIDFile:  true,
			pidContent:    "not-a-number",
			expectRunning: false,
			expectError:   true,
		},
		{
			name:          "agent with stale PID file",
			agentName:     "test-agent-stale",
			setupPIDFile:  true,
			pidContent:    "99999",
			expectRunning: false,
			expectError:   false,
		},
		{
			name:          "agent currently running",
			agentName:     "test-agent-running",
			setupPIDFile:  true,
			pidContent:    strconv.Itoa(os.Getpid()), // Use our own PID
			expectRunning: true,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pidFile := filepath.Join("/tmp", "sloth-runner-agent-"+tt.agentName+".pid")

			// Cleanup before test
			os.Remove(pidFile)

			// Setup PID file if needed
			if tt.setupPIDFile {
				if err := os.WriteFile(pidFile, []byte(tt.pidContent), 0644); err != nil {
					t.Fatalf("Failed to create test PID file: %v", err)
				}
				defer os.Remove(pidFile)
			}

			// Run the test
			info, err := checkExistingAgent(tt.agentName)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check running state
			if !tt.expectError && info != nil {
				if info.Running != tt.expectRunning {
					t.Errorf("Expected running=%v, got running=%v", tt.expectRunning, info.Running)
				}

				if info.PIDFile != pidFile {
					t.Errorf("Expected PIDFile=%q, got %q", pidFile, info.PIDFile)
				}
			}
		})
	}
}

// Test buildDaemonCommandArgs function
func TestBuildDaemonCommandArgs(t *testing.T) {
	tests := []struct {
		name          string
		opts          StartAgentOptions
		expectedArgs  []string
		checkContains []string
	}{
		{
			name: "basic options",
			opts: StartAgentOptions{
				Port:       50052,
				MasterAddr: "localhost:50051",
				AgentName:  "test-agent",
			},
			checkContains: []string{
				"agent", "start",
				"--port", "50052",
				"--name", "test-agent",
				"--master", "localhost:50051",
			},
		},
		{
			name: "with bind address",
			opts: StartAgentOptions{
				Port:        50052,
				MasterAddr:  "localhost:50051",
				AgentName:   "test-agent",
				BindAddress: "192.168.1.10",
			},
			checkContains: []string{
				"--bind-address", "192.168.1.10",
			},
		},
		{
			name: "with report address",
			opts: StartAgentOptions{
				Port:          50052,
				MasterAddr:    "localhost:50051",
				AgentName:     "test-agent",
				ReportAddress: "192.168.1.100",
			},
			checkContains: []string{
				"--report-address", "192.168.1.100",
			},
		},
		{
			name: "with telemetry enabled",
			opts: StartAgentOptions{
				Port:             50052,
				MasterAddr:       "localhost:50051",
				AgentName:        "test-agent",
				TelemetryEnabled: true,
			},
			checkContains: []string{
				"--telemetry",
			},
		},
		{
			name: "with custom metrics port",
			opts: StartAgentOptions{
				Port:        50052,
				MasterAddr:  "localhost:50051",
				AgentName:   "test-agent",
				MetricsPort: 8080,
			},
			checkContains: []string{
				"--metrics-port", "8080",
			},
		},
		{
			name: "all options",
			opts: StartAgentOptions{
				Port:             50052,
				MasterAddr:       "localhost:50051",
				AgentName:        "test-agent",
				BindAddress:      "192.168.1.10",
				ReportAddress:    "192.168.1.100",
				TelemetryEnabled: true,
				MetricsPort:      8080,
			},
			checkContains: []string{
				"agent", "start",
				"--port", "50052",
				"--name", "test-agent",
				"--master", "localhost:50051",
				"--bind-address", "192.168.1.10",
				"--report-address", "192.168.1.100",
				"--telemetry",
				"--metrics-port", "8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildDaemonCommandArgs(tt.opts)

			// Convert args to string for easy checking
			argsStr := strings.Join(args, " ")

			// Check that all expected strings are present
			for _, expected := range tt.checkContains {
				if !strings.Contains(argsStr, expected) {
					t.Errorf("Expected args to contain %q, got: %v", expected, args)
				}
			}
		})
	}
}

// Test writePIDFile and removePIDFile functions
func TestPIDFileOperations(t *testing.T) {
	pidFile := filepath.Join("/tmp", "test-sloth-runner.pid")
	testPID := 12345

	// Cleanup before test
	os.Remove(pidFile)
	defer os.Remove(pidFile)

	// Test writing PID file
	if err := writePIDFile(pidFile, testPID); err != nil {
		t.Fatalf("Failed to write PID file: %v", err)
	}

	// Verify file exists and contains correct PID
	content, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatalf("Failed to read PID file: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		t.Fatalf("Invalid PID in file: %v", err)
	}

	if pid != testPID {
		t.Errorf("Expected PID %d, got %d", testPID, pid)
	}

	// Test removing PID file
	if err := removePIDFile(pidFile); err != nil {
		t.Fatalf("Failed to remove PID file: %v", err)
	}

	// Verify file no longer exists
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Error("PID file still exists after removal")
	}

	// Test removing non-existent file (should not error)
	if err := removePIDFile(pidFile); err != nil {
		t.Errorf("Unexpected error removing non-existent file: %v", err)
	}
}

// Test determineListenAddress function
func TestDetermineListenAddress(t *testing.T) {
	tests := []struct {
		name     string
		opts     StartAgentOptions
		expected string
	}{
		{
			name: "default listen on all interfaces",
			opts: StartAgentOptions{
				Port: 50052,
			},
			expected: ":50052",
		},
		{
			name: "bind to specific address",
			opts: StartAgentOptions{
				Port:        50052,
				BindAddress: "192.168.1.10",
			},
			expected: "192.168.1.10:50052",
		},
		{
			name: "bind to localhost",
			opts: StartAgentOptions{
				Port:        8080,
				BindAddress: "127.0.0.1",
			},
			expected: "127.0.0.1:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineListenAddress(tt.opts)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test determineReportAddress function
func TestDetermineReportAddress(t *testing.T) {
	tests := []struct {
		name             string
		opts             StartAgentOptions
		actualListenAddr string
		expected         string
	}{
		{
			name: "use report address with port",
			opts: StartAgentOptions{
				Port:          50052,
				ReportAddress: "192.168.1.100:50052",
			},
			actualListenAddr: "0.0.0.0:50052",
			expected:         "192.168.1.100:50052",
		},
		{
			name: "use report address without port",
			opts: StartAgentOptions{
				Port:          50052,
				ReportAddress: "192.168.1.100",
			},
			actualListenAddr: "0.0.0.0:50052",
			expected:         "192.168.1.100:50052",
		},
		{
			name: "use bind address when no report address",
			opts: StartAgentOptions{
				Port:        50052,
				BindAddress: "192.168.1.10",
			},
			actualListenAddr: "192.168.1.10:50052",
			expected:         "192.168.1.10:50052",
		},
		{
			name: "use actual listen address as fallback",
			opts: StartAgentOptions{
				Port: 50052,
			},
			actualListenAddr: "0.0.0.0:50052",
			expected:         "0.0.0.0:50052",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineReportAddress(tt.opts, tt.actualListenAddr)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test registerAgentWithMaster function
func TestRegisterAgentWithMaster(t *testing.T) {
	tests := []struct {
		name          string
		agentName     string
		reportAddress string
		setupMock     func(*mocks.MockAgentRegistryClient)
		expectError   bool
	}{
		{
			name:          "successful registration",
			agentName:     "test-agent",
			reportAddress: "192.168.1.10:50052",
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.RegisterAgentFunc = func(ctx context.Context, in *pb.RegisterAgentRequest, opts ...grpc.CallOption) (*pb.RegisterAgentResponse, error) {
					if in.AgentName != "test-agent" {
						t.Errorf("Expected agent name 'test-agent', got %q", in.AgentName)
					}
					if in.AgentAddress != "192.168.1.10:50052" {
						t.Errorf("Expected address '192.168.1.10:50052', got %q", in.AgentAddress)
					}
					return &pb.RegisterAgentResponse{Success: true}, nil
				}
			},
			expectError: false,
		},
		{
			name:          "registration failure",
			agentName:     "test-agent",
			reportAddress: "192.168.1.10:50052",
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.RegisterAgentFunc = func(ctx context.Context, in *pb.RegisterAgentRequest, opts ...grpc.CallOption) (*pb.RegisterAgentResponse, error) {
					return nil, grpc.ErrServerStopped
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			tt.setupMock(mockClient)

			err := registerAgentWithMaster(context.Background(), mockClient, tt.agentName, tt.reportAddress)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test sendHeartbeatWithSystemInfo function
func TestSendHeartbeatWithSystemInfo(t *testing.T) {
	tests := []struct {
		name              string
		agentName         string
		includeSystemInfo bool
		setupMock         func(*mocks.MockAgentRegistryClient)
		expectError       bool
		checkSystemInfo   bool
	}{
		{
			name:              "heartbeat without system info",
			agentName:         "test-agent",
			includeSystemInfo: false,
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.HeartbeatFunc = func(ctx context.Context, in *pb.HeartbeatRequest, opts ...grpc.CallOption) (*pb.HeartbeatResponse, error) {
					if in.AgentName != "test-agent" {
						t.Errorf("Expected agent name 'test-agent', got %q", in.AgentName)
					}
					if in.SystemInfoJson != "" {
						t.Errorf("Expected empty system info, got: %q", in.SystemInfoJson)
					}
					return &pb.HeartbeatResponse{Success: true}, nil
				}
			},
			expectError: false,
		},
		{
			name:              "heartbeat with system info",
			agentName:         "test-agent",
			includeSystemInfo: true,
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.HeartbeatFunc = func(ctx context.Context, in *pb.HeartbeatRequest, opts ...grpc.CallOption) (*pb.HeartbeatResponse, error) {
					if in.AgentName != "test-agent" {
						t.Errorf("Expected agent name 'test-agent', got %q", in.AgentName)
					}
					// System info may or may not be collected successfully, so just check it's not causing errors
					return &pb.HeartbeatResponse{Success: true}, nil
				}
			},
			expectError:     false,
			checkSystemInfo: true,
		},
		{
			name:              "heartbeat failure",
			agentName:         "test-agent",
			includeSystemInfo: false,
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.HeartbeatFunc = func(ctx context.Context, in *pb.HeartbeatRequest, opts ...grpc.CallOption) (*pb.HeartbeatResponse, error) {
					return nil, grpc.ErrServerStopped
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			tt.setupMock(mockClient)

			err := sendHeartbeatWithSystemInfo(context.Background(), mockClient, tt.agentName, tt.includeSystemInfo)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test context cancellation in registerAgentWithMaster
func TestRegisterAgentWithMaster_ContextCancellation(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.RegisterAgentFunc = func(ctx context.Context, in *pb.RegisterAgentRequest, opts ...grpc.CallOption) (*pb.RegisterAgentResponse, error) {
		return nil, context.Canceled
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := registerAgentWithMaster(ctx, mockClient, "test-agent", "192.168.1.10:50052")

	if err == nil {
		t.Error("Expected error due to context cancellation")
	}
}

// Test context cancellation in sendHeartbeatWithSystemInfo
func TestSendHeartbeatWithSystemInfo_ContextCancellation(t *testing.T) {
	mockClient := mocks.NewMockAgentRegistryClient()
	mockClient.HeartbeatFunc = func(ctx context.Context, in *pb.HeartbeatRequest, opts ...grpc.CallOption) (*pb.HeartbeatResponse, error) {
		return nil, context.Canceled
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := sendHeartbeatWithSystemInfo(ctx, mockClient, "test-agent", false)

	if err == nil {
		t.Error("Expected error due to context cancellation")
	}
}
