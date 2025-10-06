package agent

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	"github.com/chalkan3-sloth/sloth-runner/internal/telemetry"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test extractHost function
func TestExtractHost(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected string
	}{
		{
			name:     "address with port",
			address:  "192.168.1.10:50051",
			expected: "192.168.1.10",
		},
		{
			name:     "address without port",
			address:  "192.168.1.10",
			expected: "192.168.1.10",
		},
		{
			name:     "hostname with port",
			address:  "agent-server:8080",
			expected: "agent-server",
		},
		{
			name:     "localhost with port",
			address:  "localhost:50051",
			expected: "localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHost(tt.address)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test displayMetricsEndpoint function
func TestDisplayMetricsEndpoint(t *testing.T) {
	tests := []struct {
		name            string
		metricsEndpoint string
		host            string
		checkOutput     func(string) bool
	}{
		{
			name:            "basic endpoint",
			metricsEndpoint: "http://192.168.1.10:9090/metrics",
			host:            "192.168.1.10",
			checkOutput: func(out string) bool {
				return strings.Contains(out, "192.168.1.10:9090") &&
					(strings.Contains(out, "URL") || strings.Contains(out, "url"))
			},
		},
		{
			name:            "localhost endpoint",
			metricsEndpoint: "http://localhost:9090/metrics",
			host:            "localhost",
			checkOutput: func(out string) bool {
				return strings.Contains(out, "localhost:9090") &&
					strings.Contains(out, "curl") &&
					strings.Contains(out, "Prometheus")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := displayMetricsEndpoint(tt.metricsEndpoint, tt.host, &buf)

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

// Test prometheusMetricsWithClient function
func TestPrometheusMetricsWithClient(t *testing.T) {
	tests := []struct {
		name          string
		opts          MetricsOptions
		setupMock     func(*mocks.MockAgentRegistryClient)
		expectedError bool
		checkOutput   func(string) bool
	}{
		{
			name: "display endpoint without snapshot",
			opts: MetricsOptions{
				AgentName:    "test-agent",
				ShowSnapshot: false,
			},
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{
							{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
						},
					}, nil
				}
			},
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "192.168.1.10:9090")
			},
		},
		{
			name: "agent not found",
			opts: MetricsOptions{
				AgentName:    "missing-agent",
				ShowSnapshot: false,
			},
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{},
					}, nil
				}
			},
			expectedError: true,
		},
		{
			name: "display snapshot with telemetry unavailable",
			opts: MetricsOptions{
				AgentName:    "test-agent",
				ShowSnapshot: true,
			},
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{
							{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
						},
					}, nil
				}
				m.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
					return &mocks.MockExecuteCommandClient{
						Responses: []*pb.StreamOutputResponse{
							{StdoutChunk: "curl: (7) Failed to connect to localhost port 9090: Connection refused"},
						},
					}, nil
				}
			},
			expectedError: false,
			checkOutput: func(out string) bool {
				// Just check that we got some output (error message or help text)
				return len(out) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			tt.setupMock(mockClient)

			var buf bytes.Buffer
			tt.opts.Writer = &buf

			err := prometheusMetricsWithClient(context.Background(), mockClient, tt.opts)

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

// Test displayMetricsSnapshot function
func TestDisplayMetricsSnapshot(t *testing.T) {
	tests := []struct {
		name          string
		agentName     string
		setupMock     func(*mocks.MockAgentRegistryClient)
		expectedError bool
		checkOutput   func(string) bool
	}{
		{
			name:      "successful metrics fetch",
			agentName: "test-agent",
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
					return &mocks.MockExecuteCommandClient{
						Responses: []*pb.StreamOutputResponse{
							{StdoutChunk: "# TYPE go_goroutines gauge\ngo_goroutines 42\n"},
						},
					}, nil
				}
			},
			expectedError: false,
			checkOutput: func(out string) bool {
				return strings.Contains(out, "go_goroutines") &&
					strings.Contains(out, "Metrics Snapshot")
			},
		},
		{
			name:      "connection refused",
			agentName: "test-agent",
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
					return &mocks.MockExecuteCommandClient{
						Responses: []*pb.StreamOutputResponse{
							{StdoutChunk: "Connection refused"},
						},
					}, nil
				}
			},
			expectedError: false,
			checkOutput: func(out string) bool {
				// Just check that we got some output
				return len(out) > 0
			},
		},
		{
			name:      "execute command error",
			agentName: "test-agent",
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ExecuteCommandFunc = func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
					return nil, grpc.ErrServerStopped
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
			err := displayMetricsSnapshot(context.Background(), mockClient, tt.agentName, &buf)

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

// Test grafanaDashboardWithClient function
func TestGrafanaDashboardWithClient(t *testing.T) {
	tests := []struct {
		name          string
		opts          DashboardOptions
		setupMock     func(*mocks.MockAgentRegistryClient)
		expectedError bool
	}{
		{
			name: "successful dashboard display",
			opts: DashboardOptions{
				AgentName: "test-agent",
				Watch:     false,
				MetricsFetcher: func(endpoint string) (*telemetry.MetricsData, error) {
					return &telemetry.MetricsData{}, nil
				},
				DashboardDisplay: func(data *telemetry.MetricsData, agentName string) {
					// Mock display function
				},
			},
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{
							{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
						},
					}, nil
				}
			},
			expectedError: false,
		},
		{
			name: "metrics fetch failure",
			opts: DashboardOptions{
				AgentName: "test-agent",
				Watch:     false,
				MetricsFetcher: func(endpoint string) (*telemetry.MetricsData, error) {
					return nil, io.EOF
				},
			},
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{
							{AgentName: "test-agent", AgentAddress: "192.168.1.10:50051"},
						},
					}, nil
				}
			},
			expectedError: true,
		},
		{
			name: "agent not found",
			opts: DashboardOptions{
				AgentName: "missing-agent",
				Watch:     false,
			},
			setupMock: func(m *mocks.MockAgentRegistryClient) {
				m.ListAgentsFunc = func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
					return &pb.ListAgentsResponse{
						Agents: []*pb.AgentInfo{},
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

			err := grafanaDashboardWithClient(context.Background(), mockClient, tt.opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
