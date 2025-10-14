package mocks

import (
	"context"
	"io"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// MockAgentRegistryClient is a mock implementation of pb.AgentRegistryClient
type MockAgentRegistryClient struct {
	ListAgentsFunc      func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error)
	GetAgentInfoFunc    func(ctx context.Context, in *pb.GetAgentInfoRequest, opts ...grpc.CallOption) (*pb.GetAgentInfoResponse, error)
	ExecuteCommandFunc  func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error)
	StopAgentFunc       func(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error)
	UnregisterAgentFunc func(ctx context.Context, in *pb.UnregisterAgentRequest, opts ...grpc.CallOption) (*pb.UnregisterAgentResponse, error)
	RegisterAgentFunc   func(ctx context.Context, in *pb.RegisterAgentRequest, opts ...grpc.CallOption) (*pb.RegisterAgentResponse, error)
	HeartbeatFunc       func(ctx context.Context, in *pb.HeartbeatRequest, opts ...grpc.CallOption) (*pb.HeartbeatResponse, error)
}

func (m *MockAgentRegistryClient) RegisterAgent(ctx context.Context, in *pb.RegisterAgentRequest, opts ...grpc.CallOption) (*pb.RegisterAgentResponse, error) {
	if m.RegisterAgentFunc != nil {
		return m.RegisterAgentFunc(ctx, in, opts...)
	}
	return &pb.RegisterAgentResponse{Success: true}, nil
}

func (m *MockAgentRegistryClient) ListAgents(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error) {
	if m.ListAgentsFunc != nil {
		return m.ListAgentsFunc(ctx, in, opts...)
	}
	return &pb.ListAgentsResponse{
		Agents: []*pb.AgentInfo{
			{
				AgentName:    "test-agent",
				AgentAddress: "localhost:50052",
				Status:       "Active",
			},
		},
	}, nil
}

func (m *MockAgentRegistryClient) StopAgent(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error) {
	if m.StopAgentFunc != nil {
		return m.StopAgentFunc(ctx, in, opts...)
	}
	return &pb.StopAgentResponse{Success: true}, nil
}

func (m *MockAgentRegistryClient) UnregisterAgent(ctx context.Context, in *pb.UnregisterAgentRequest, opts ...grpc.CallOption) (*pb.UnregisterAgentResponse, error) {
	if m.UnregisterAgentFunc != nil {
		return m.UnregisterAgentFunc(ctx, in, opts...)
	}
	return &pb.UnregisterAgentResponse{Success: true}, nil
}

func (m *MockAgentRegistryClient) ExecuteCommand(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error) {
	if m.ExecuteCommandFunc != nil {
		return m.ExecuteCommandFunc(ctx, in, opts...)
	}
	return &MockExecuteCommandClient{}, nil
}

func (m *MockAgentRegistryClient) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest, opts ...grpc.CallOption) (*pb.HeartbeatResponse, error) {
	if m.HeartbeatFunc != nil {
		return m.HeartbeatFunc(ctx, in, opts...)
	}
	return &pb.HeartbeatResponse{}, nil
}

func (m *MockAgentRegistryClient) GetAgentInfo(ctx context.Context, in *pb.GetAgentInfoRequest, opts ...grpc.CallOption) (*pb.GetAgentInfoResponse, error) {
	if m.GetAgentInfoFunc != nil {
		return m.GetAgentInfoFunc(ctx, in, opts...)
	}
	return &pb.GetAgentInfoResponse{
		Success: true,
		AgentInfo: &pb.AgentInfo{
			AgentName:         "test-agent",
			AgentAddress:      "localhost:50052",
			Status:            "Active",
			LastHeartbeat:     1234567890,
			LastInfoCollected: 1234567890,
			SystemInfoJson:    `{"hostname":"test-host","platform":"linux","architecture":"amd64"}`,
		},
	}, nil
}

// MockExecuteCommandClient is a mock for streaming command execution
type MockExecuteCommandClient struct {
	grpc.ClientStream
	Responses []*pb.StreamOutputResponse
	Index     int
	RecvError error // Error to return after responses are exhausted
}

func (m *MockExecuteCommandClient) Recv() (*pb.StreamOutputResponse, error) {
	if m.Index >= len(m.Responses) {
		if m.RecvError != nil {
			return nil, m.RecvError
		}
		return nil, io.EOF
	}
	resp := m.Responses[m.Index]
	m.Index++
	return resp, nil
}

func (m *MockExecuteCommandClient) CloseSend() error {
	return nil
}

// MockAgentClient is a mock implementation of pb.AgentClient
type MockAgentClient struct {
	UpdateAgentFunc     func(ctx context.Context, in *pb.UpdateAgentRequest, opts ...grpc.CallOption) (*pb.UpdateAgentResponse, error)
	RegisterWatcherFunc func(ctx context.Context, in *pb.RegisterWatcherRequest, opts ...grpc.CallOption) (*pb.RegisterWatcherResponse, error)
	ListWatchersFunc    func(ctx context.Context, in *pb.ListWatchersRequest, opts ...grpc.CallOption) (*pb.ListWatchersResponse, error)
	GetWatcherFunc      func(ctx context.Context, in *pb.GetWatcherRequest, opts ...grpc.CallOption) (*pb.GetWatcherResponse, error)
	RemoveWatcherFunc   func(ctx context.Context, in *pb.RemoveWatcherRequest, opts ...grpc.CallOption) (*pb.RemoveWatcherResponse, error)
}

func (m *MockAgentClient) ExecuteTask(ctx context.Context, in *pb.ExecuteTaskRequest, opts ...grpc.CallOption) (*pb.ExecuteTaskResponse, error) {
	return &pb.ExecuteTaskResponse{Success: true}, nil
}

func (m *MockAgentClient) RunCommand(ctx context.Context, in *pb.RunCommandRequest, opts ...grpc.CallOption) (pb.Agent_RunCommandClient, error) {
	return nil, nil
}

func (m *MockAgentClient) Shutdown(ctx context.Context, in *pb.ShutdownRequest, opts ...grpc.CallOption) (*pb.ShutdownResponse, error) {
	return &pb.ShutdownResponse{}, nil
}

func (m *MockAgentClient) UpdateAgent(ctx context.Context, in *pb.UpdateAgentRequest, opts ...grpc.CallOption) (*pb.UpdateAgentResponse, error) {
	if m.UpdateAgentFunc != nil {
		return m.UpdateAgentFunc(ctx, in, opts...)
	}
	return &pb.UpdateAgentResponse{
		Success:    true,
		Message:    "Update successful",
		OldVersion: "1.0.0",
		NewVersion: "1.1.0",
	}, nil
}

func (m *MockAgentClient) RegisterWatcher(ctx context.Context, in *pb.RegisterWatcherRequest, opts ...grpc.CallOption) (*pb.RegisterWatcherResponse, error) {
	if m.RegisterWatcherFunc != nil {
		return m.RegisterWatcherFunc(ctx, in, opts...)
	}
	return &pb.RegisterWatcherResponse{
		Success:   true,
		Message:   "Watcher registered successfully",
		WatcherId: "mock-watcher-id",
	}, nil
}

func (m *MockAgentClient) ListWatchers(ctx context.Context, in *pb.ListWatchersRequest, opts ...grpc.CallOption) (*pb.ListWatchersResponse, error) {
	if m.ListWatchersFunc != nil {
		return m.ListWatchersFunc(ctx, in, opts...)
	}
	return &pb.ListWatchersResponse{
		Watchers: []*pb.WatcherConfig{},
	}, nil
}

func (m *MockAgentClient) GetWatcher(ctx context.Context, in *pb.GetWatcherRequest, opts ...grpc.CallOption) (*pb.GetWatcherResponse, error) {
	if m.GetWatcherFunc != nil {
		return m.GetWatcherFunc(ctx, in, opts...)
	}
	return &pb.GetWatcherResponse{
		Found:   true,
		Watcher: &pb.WatcherConfig{},
	}, nil
}

func (m *MockAgentClient) RemoveWatcher(ctx context.Context, in *pb.RemoveWatcherRequest, opts ...grpc.CallOption) (*pb.RemoveWatcherResponse, error) {
	if m.RemoveWatcherFunc != nil {
		return m.RemoveWatcherFunc(ctx, in, opts...)
	}
	return &pb.RemoveWatcherResponse{
		Success: true,
		Message: "Watcher removed successfully",
	}, nil
}

// NewMockAgentRegistryClient creates a new mock with default implementations
func NewMockAgentRegistryClient() *MockAgentRegistryClient {
	return &MockAgentRegistryClient{}
}

// NewMockAgentClient creates a new mock agent client
func NewMockAgentClient() *MockAgentClient {
	return &MockAgentClient{}
}

// WithListAgents sets a custom ListAgents function
func (m *MockAgentRegistryClient) WithListAgents(f func(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error)) *MockAgentRegistryClient {
	m.ListAgentsFunc = f
	return m
}

// WithGetAgentInfo sets a custom GetAgentInfo function
func (m *MockAgentRegistryClient) WithGetAgentInfo(f func(ctx context.Context, in *pb.GetAgentInfoRequest, opts ...grpc.CallOption) (*pb.GetAgentInfoResponse, error)) *MockAgentRegistryClient {
	m.GetAgentInfoFunc = f
	return m
}

// WithExecuteCommand sets a custom ExecuteCommand function
func (m *MockAgentRegistryClient) WithExecuteCommand(f func(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error)) *MockAgentRegistryClient {
	m.ExecuteCommandFunc = f
	return m
}
