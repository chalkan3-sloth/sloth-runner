package agent

import (
	"context"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// AgentRegistryClient interface for dependency injection
type AgentRegistryClient interface {
	ListAgents(ctx context.Context, in *pb.ListAgentsRequest, opts ...grpc.CallOption) (*pb.ListAgentsResponse, error)
	GetAgentInfo(ctx context.Context, in *pb.GetAgentInfoRequest, opts ...grpc.CallOption) (*pb.GetAgentInfoResponse, error)
	ExecuteCommand(ctx context.Context, in *pb.ExecuteCommandRequest, opts ...grpc.CallOption) (pb.AgentRegistry_ExecuteCommandClient, error)
	StopAgent(ctx context.Context, in *pb.StopAgentRequest, opts ...grpc.CallOption) (*pb.StopAgentResponse, error)
	UnregisterAgent(ctx context.Context, in *pb.UnregisterAgentRequest, opts ...grpc.CallOption) (*pb.UnregisterAgentResponse, error)
	RegisterAgent(ctx context.Context, in *pb.RegisterAgentRequest, opts ...grpc.CallOption) (*pb.RegisterAgentResponse, error)
	Heartbeat(ctx context.Context, in *pb.HeartbeatRequest, opts ...grpc.CallOption) (*pb.HeartbeatResponse, error)
}

// AgentClient interface for dependency injection
type AgentClient interface {
	UpdateAgent(ctx context.Context, in *pb.UpdateAgentRequest, opts ...grpc.CallOption) (*pb.UpdateAgentResponse, error)
	RegisterWatcher(ctx context.Context, in *pb.RegisterWatcherRequest, opts ...grpc.CallOption) (*pb.RegisterWatcherResponse, error)
	ListWatchers(ctx context.Context, in *pb.ListWatchersRequest, opts ...grpc.CallOption) (*pb.ListWatchersResponse, error)
	RemoveWatcher(ctx context.Context, in *pb.RemoveWatcherRequest, opts ...grpc.CallOption) (*pb.RemoveWatcherResponse, error)
}

// AgentService provides agent operations with injected dependencies
type AgentService struct {
	registryClient AgentRegistryClient
	agentClient    AgentClient
}

// NewAgentService creates a new agent service with injected clients
func NewAgentService(registryClient AgentRegistryClient, agentClient AgentClient) *AgentService {
	return &AgentService{
		registryClient: registryClient,
		agentClient:    agentClient,
	}
}

// ConnectionFactory creates gRPC connections and clients
type ConnectionFactory interface {
	CreateRegistryClient(masterAddr string) (AgentRegistryClient, func(), error)
	CreateAgentClient(agentAddr string) (AgentClient, func(), error)
}

// DefaultConnectionFactory implements ConnectionFactory with real gRPC connections
type DefaultConnectionFactory struct{}

// NewDefaultConnectionFactory creates a new default connection factory
func NewDefaultConnectionFactory() *DefaultConnectionFactory {
	return &DefaultConnectionFactory{}
}

// CreateRegistryClient creates a real registry client connection
func (f *DefaultConnectionFactory) CreateRegistryClient(masterAddr string) (AgentRegistryClient, func(), error) {
	conn, err := createGRPCConnection(masterAddr)
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewAgentRegistryClient(conn)
	cleanup := func() { conn.Close() }

	return client, cleanup, nil
}

// CreateAgentClient creates a real agent client connection
func (f *DefaultConnectionFactory) CreateAgentClient(agentAddr string) (AgentClient, func(), error) {
	conn, err := createGRPCConnection(agentAddr)
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewAgentClient(conn)
	cleanup := func() { conn.Close() }

	return client, cleanup, nil
}

// RegisterWatcherOnAgent registers a watcher on a remote agent via gRPC
func RegisterWatcherOnAgent(ctx context.Context, agentAddr string, config *pb.WatcherConfig) (*pb.RegisterWatcherResponse, error) {
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateAgentClient(agentAddr)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	req := &pb.RegisterWatcherRequest{
		Config: config,
	}

	return client.RegisterWatcher(ctx, req)
}

// ListWatchersOnAgent lists all watchers on a remote agent via gRPC
func ListWatchersOnAgent(ctx context.Context, agentAddr string) (*pb.ListWatchersResponse, error) {
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateAgentClient(agentAddr)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	req := &pb.ListWatchersRequest{}
	return client.ListWatchers(ctx, req)
}

// RemoveWatcherFromAgent removes a watcher from a remote agent via gRPC
func RemoveWatcherFromAgent(ctx context.Context, agentAddr string, watcherID string) (*pb.RemoveWatcherResponse, error) {
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateAgentClient(agentAddr)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	req := &pb.RemoveWatcherRequest{
		WatcherId: watcherID,
	}

	return client.RemoveWatcher(ctx, req)
}
