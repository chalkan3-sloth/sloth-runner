package services

import (
	"context"
	"fmt"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AgentService handles agent operations via gRPC
// This implements the Service Layer pattern for agent management
type AgentService struct {
	masterAddr string
	timeout    time.Duration
}

// NewAgentService creates a new agent service
func NewAgentService(masterAddr string) *AgentService {
	if masterAddr == "" {
		masterAddr = "192.168.1.29:50053" // Default master address
	}
	return &AgentService{
		masterAddr: masterAddr,
		timeout:    10 * time.Minute, // Increased timeout for long-running operations like image downloads
	}
}

// SetTimeout sets the timeout for operations
func (s *AgentService) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// ListAgents lists all registered agents
func (s *AgentService) ListAgents() ([]*pb.AgentInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	conn, err := s.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	registryClient := pb.NewAgentRegistryClient(conn)
	resp, err := registryClient.ListAgents(ctx, &pb.ListAgentsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %w", err)
	}

	return resp.GetAgents(), nil
}

// StopAgent stops a specific agent
func (s *AgentService) StopAgent(agentName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	conn, err := s.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	registryClient := pb.NewAgentRegistryClient(conn)
	resp, err := registryClient.StopAgent(ctx, &pb.StopAgentRequest{
		AgentName: agentName,
	})
	if err != nil {
		return fmt.Errorf("failed to stop agent: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("failed to stop agent: %s", resp.GetMessage())
	}

	return nil
}

// DeleteAgent deletes a specific agent
func (s *AgentService) DeleteAgent(agentName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	conn, err := s.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	registryClient := pb.NewAgentRegistryClient(conn)
	resp, err := registryClient.UnregisterAgent(ctx, &pb.UnregisterAgentRequest{
		AgentName: agentName,
	})
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("failed to delete agent: %s", resp.GetMessage())
	}

	return nil
}

// GetAgent gets information about a specific agent
func (s *AgentService) GetAgent(agentName string) (*pb.AgentInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	conn, err := s.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	registryClient := pb.NewAgentRegistryClient(conn)
	resp, err := registryClient.GetAgentInfo(ctx, &pb.GetAgentInfoRequest{
		AgentName: agentName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get agent info: %w", err)
	}

	if !resp.GetSuccess() {
		return nil, fmt.Errorf("failed to get agent info: %s", resp.GetMessage())
	}

	return resp.GetAgentInfo(), nil
}

// connect establishes a gRPC connection to the master
func (s *AgentService) connect() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(s.masterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master at %s: %w", s.masterAddr, err)
	}
	return conn, nil
}

// GetMasterAddr returns the master address
func (s *AgentService) GetMasterAddr() string {
	return s.masterAddr
}
