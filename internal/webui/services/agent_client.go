package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AgentClient manages connections to agents
type AgentClient struct {
	connections map[string]*grpc.ClientConn
	clients     map[string]pb.AgentClient
	mu          sync.RWMutex
	timeout     time.Duration
}

// NewAgentClient creates a new agent client manager
func NewAgentClient() *AgentClient {
	return &AgentClient{
		connections: make(map[string]*grpc.ClientConn),
		clients:     make(map[string]pb.AgentClient),
		timeout:     30 * time.Second,
	}
}

// GetClient returns a gRPC client for the given agent address
func (ac *AgentClient) GetClient(ctx context.Context, agentAddress string) (pb.AgentClient, error) {
	ac.mu.RLock()
	client, exists := ac.clients[agentAddress]
	ac.mu.RUnlock()

	if exists {
		return client, nil
	}

	// Create new connection
	ac.mu.Lock()
	defer ac.mu.Unlock()

	// Double check after acquiring write lock
	if client, exists := ac.clients[agentAddress]; exists {
		return client, nil
	}

	conn, err := grpc.Dial(agentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to agent %s: %w", agentAddress, err)
	}

	client = pb.NewAgentClient(conn)
	ac.connections[agentAddress] = conn
	ac.clients[agentAddress] = client

	return client, nil
}

// RunCommand executes a command on an agent and returns a streaming response
func (ac *AgentClient) RunCommand(ctx context.Context, agentAddress, command string) (pb.Agent_RunCommandClient, error) {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return nil, err
	}

	stream, err := client.RunCommand(ctx, &pb.RunCommandRequest{
		Command: command,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	return stream, nil
}

// GetResourceUsage retrieves resource usage from an agent
func (ac *AgentClient) GetResourceUsage(ctx context.Context, agentAddress string) (*pb.ResourceUsageResponse, error) {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.timeout)
	defer cancel()

	resp, err := client.GetResourceUsage(ctx, &pb.ResourceUsageRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get resource usage: %w", err)
	}

	return resp, nil
}

// GetProcessList retrieves the process list from an agent
func (ac *AgentClient) GetProcessList(ctx context.Context, agentAddress string) (*pb.ProcessListResponse, error) {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.timeout)
	defer cancel()

	resp, err := client.GetProcessList(ctx, &pb.ProcessListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get process list: %w", err)
	}

	return resp, nil
}

// GetNetworkInfo retrieves network information from an agent
func (ac *AgentClient) GetNetworkInfo(ctx context.Context, agentAddress string) (*pb.NetworkInfoResponse, error) {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.timeout)
	defer cancel()

	resp, err := client.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get network info: %w", err)
	}

	return resp, nil
}

// GetDiskInfo retrieves disk information from an agent
func (ac *AgentClient) GetDiskInfo(ctx context.Context, agentAddress string) (*pb.DiskInfoResponse, error) {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.timeout)
	defer cancel()

	resp, err := client.GetDiskInfo(ctx, &pb.DiskInfoRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get disk info: %w", err)
	}

	return resp, nil
}

// StreamLogs streams logs from an agent
func (ac *AgentClient) StreamLogs(ctx context.Context, agentAddress string) (pb.Agent_StreamLogsClient, error) {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return nil, err
	}

	stream, err := client.StreamLogs(ctx, &pb.StreamLogsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to stream logs: %w", err)
	}

	return stream, nil
}

// StreamMetrics streams metrics from an agent
func (ac *AgentClient) StreamMetrics(ctx context.Context, agentAddress string) (pb.Agent_StreamMetricsClient, error) {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return nil, err
	}

	stream, err := client.StreamMetrics(ctx, &pb.StreamMetricsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to stream metrics: %w", err)
	}

	return stream, nil
}

// RestartService restarts the agent service
func (ac *AgentClient) RestartService(ctx context.Context, agentAddress string) error {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.timeout)
	defer cancel()

	_, err = client.RestartService(ctx, &pb.RestartServiceRequest{})
	if err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	return nil
}

// Shutdown shuts down an agent
func (ac *AgentClient) Shutdown(ctx context.Context, agentAddress string) error {
	client, err := ac.GetClient(ctx, agentAddress)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, ac.timeout)
	defer cancel()

	_, err = client.Shutdown(ctx, &pb.ShutdownRequest{})
	if err != nil {
		return fmt.Errorf("failed to shutdown agent: %w", err)
	}

	return nil
}

// CloseConnection closes the connection to a specific agent
func (ac *AgentClient) CloseConnection(agentAddress string) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	conn, exists := ac.connections[agentAddress]
	if !exists {
		return nil
	}

	err := conn.Close()
	delete(ac.connections, agentAddress)
	delete(ac.clients, agentAddress)

	return err
}

// CloseAll closes all agent connections
func (ac *AgentClient) CloseAll() error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	var errs []error
	for addr, conn := range ac.connections {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection to %s: %w", addr, err))
		}
	}

	ac.connections = make(map[string]*grpc.ClientConn)
	ac.clients = make(map[string]pb.AgentClient)

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}
