package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"sync"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// agentRegistryServer implements the AgentRegistry service.
type agentRegistryServer struct {
	pb.UnimplementedAgentRegistryServer
	mu         sync.RWMutex
	db         *AgentDB
	grpcServer *grpc.Server
}

// newAgentRegistryServer creates a new agentRegistryServer.
func newAgentRegistryServer() *agentRegistryServer {
	// Create SQLite database for agents
	dbPath := filepath.Join(".", ".sloth-cache", "agents.db")
	db, err := NewAgentDB(dbPath)
	if err != nil {
		pterm.Error.Printf("Failed to initialize agent database: %v\n", err)
		pterm.Info.Println("Falling back to in-memory storage")
		db = nil
	} else {
		pterm.Success.Printf("Agent database initialized at: %s\n", dbPath)
		
		// Cleanup inactive agents older than 24 hours on startup
		if removed, err := db.CleanupInactiveAgents(24); err == nil && removed > 0 {
			pterm.Info.Printf("Cleaned up %d inactive agents\n", removed)
		}
	}
	
	return &agentRegistryServer{
		db: db,
	}
}

// RegisterAgent registers a new agent.
func (s *agentRegistryServer) RegisterAgent(ctx context.Context, req *pb.RegisterAgentRequest) (*pb.RegisterAgentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pterm.Success.Printf("Agent registered: %s at %s\n", req.AgentName, req.AgentAddress)
	
	// Save to SQLite database
	if s.db != nil {
		if err := s.db.RegisterAgent(req.AgentName, req.AgentAddress); err != nil {
			pterm.Error.Printf("Failed to save agent to database: %v\n", err)
			return &pb.RegisterAgentResponse{Success: false, Message: fmt.Sprintf("Failed to save agent: %v", err)}, nil
		}
		pterm.Debug.Printf("Agent %s saved to database\n", req.AgentName)
	}

	return &pb.RegisterAgentResponse{Success: true, Message: "Agent registered successfully"}, nil
}

// ListAgents lists all registered agents.
func (s *agentRegistryServer) ListAgents(ctx context.Context, req *pb.ListAgentsRequest) (*pb.ListAgentsResponse, error) {
	pterm.Info.Println("Listing registered agents")
	s.mu.RLock()
	defer s.mu.RUnlock()

	var agents []*pb.AgentInfo
	
	if s.db != nil {
		// Get agents from SQLite database
		dbAgents, err := s.db.ListAgents()
		if err != nil {
			pterm.Error.Printf("Failed to list agents from database: %v\n", err)
			return &pb.ListAgentsResponse{Agents: agents}, nil
		}
		
		for _, agent := range dbAgents {
			agents = append(agents, &pb.AgentInfo{
				AgentName:         agent.Name,
				AgentAddress:      agent.Address,
				LastHeartbeat:     agent.LastHeartbeat,
				Status:            agent.Status,
				LastInfoCollected: agent.LastInfoCollected,
				SystemInfoJson:    agent.SystemInfo,
				Version:           agent.Version,
			})
		}
	}

	return &pb.ListAgentsResponse{Agents: agents}, nil
}

// Heartbeat updates the last heartbeat timestamp for an agent.
func (s *agentRegistryServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		// Update heartbeat and optionally store system info
		if err := s.db.UpdateHeartbeat(req.AgentName); err != nil {
			pterm.Debug.Printf("Failed to update heartbeat for agent %s: %v\n", req.AgentName, err)
			return &pb.HeartbeatResponse{Success: false, Message: "Agent not found"}, nil
		}

		// If system info is provided, update it
		if req.SystemInfoJson != "" {
			if err := s.db.UpdateSystemInfo(req.AgentName, req.SystemInfoJson); err != nil {
				pterm.Debug.Printf("Failed to update system info for agent %s: %v\n", req.AgentName, err)
			}
		}

		// If version is provided, update it
		if req.Version != "" {
			if err := s.db.UpdateVersion(req.AgentName, req.Version); err != nil {
				pterm.Debug.Printf("Failed to update version for agent %s: %v\n", req.AgentName, err)
			}
		}

		pterm.Debug.Printf("Heartbeat received from agent: %s\n", req.AgentName)
		return &pb.HeartbeatResponse{Success: true, Message: "Heartbeat received"}, nil
	}
	
	return &pb.HeartbeatResponse{Success: false, Message: "Database not available"}, nil
}

// ExecuteCommand executes a command on a remote agent and streams the output back to the client.
func (s *agentRegistryServer) ExecuteCommand(req *pb.ExecuteCommandRequest, stream pb.AgentRegistry_ExecuteCommandServer) error {
	s.mu.RLock()
	var agentAddress string
	var err error
	
	if s.db != nil {
		agentAddress, err = s.db.GetAgentAddress(req.AgentName)
		if err != nil {
			s.mu.RUnlock()
			return fmt.Errorf("agent not found or inactive: %s", req.AgentName)
		}
	} else {
		s.mu.RUnlock()
		return fmt.Errorf("database not available")
	}
	s.mu.RUnlock()

	conn, err := grpc.Dial(agentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %v", err)
	}
	defer conn.Close()

	client := pb.NewAgentClient(conn)
	agentStream, err := client.RunCommand(stream.Context(), &pb.RunCommandRequest{Command: req.Command})
	if err != nil {
		return fmt.Errorf("failed to call RunCommand on agent: %v", err)
	}

	for {
		resp, err := agentStream.Recv()
		if err == io.EOF {
			break // Stream has ended
		}
		if err != nil {
			return fmt.Errorf("error receiving stream from agent: %v", err)
		}

		// Stream the response directly to the client
		if err := stream.Send(resp); err != nil {
			return fmt.Errorf("error sending stream to client: %v", err)
		}

		if resp.GetFinished() {
			break
		}
	}

	return nil
}

// StopAgent stops a remote agent.
func (s *agentRegistryServer) StopAgent(ctx context.Context, req *pb.StopAgentRequest) (*pb.StopAgentResponse, error) {
	s.mu.RLock()
	var agentAddress string
	var err error
	
	if s.db != nil {
		agentAddress, err = s.db.GetAgentAddress(req.AgentName)
		if err != nil {
			s.mu.RUnlock()
			return nil, fmt.Errorf("agent not found or inactive: %s", req.AgentName)
		}
	} else {
		s.mu.RUnlock()
		return nil, fmt.Errorf("database not available")
	}
	s.mu.RUnlock()

	conn, err := grpc.Dial(agentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to agent: %v", err)
	}
	defer conn.Close()

	client := pb.NewAgentClient(conn)
	_, err = client.Shutdown(ctx, &pb.ShutdownRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to stop agent: %v", err)
	}

	return &pb.StopAgentResponse{Success: true, Message: "Agent stopped successfully"}, nil
}

// UnregisterAgent unregisters and removes an agent from the database.
func (s *agentRegistryServer) UnregisterAgent(ctx context.Context, req *pb.UnregisterAgentRequest) (*pb.UnregisterAgentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return &pb.UnregisterAgentResponse{Success: false, Message: "Database not available"}, nil
	}

	// Check if agent exists (using GetAgent instead of GetAgentAddress to allow inactive agents)
	_, err := s.db.GetAgent(req.AgentName)
	if err != nil {
		return &pb.UnregisterAgentResponse{Success: false, Message: fmt.Sprintf("Agent not found: %s", req.AgentName)}, nil
	}

	// Remove agent from database
	if err := s.db.UnregisterAgent(req.AgentName); err != nil {
		pterm.Error.Printf("Failed to unregister agent %s: %v\n", req.AgentName, err)
		return &pb.UnregisterAgentResponse{Success: false, Message: fmt.Sprintf("Failed to unregister agent: %v", err)}, nil
	}

	pterm.Success.Printf("Agent unregistered: %s\n", req.AgentName)
	return &pb.UnregisterAgentResponse{Success: true, Message: "Agent unregistered successfully"}, nil
}

// GetAgentInfo retrieves detailed information about a specific agent.
func (s *agentRegistryServer) GetAgentInfo(ctx context.Context, req *pb.GetAgentInfoRequest) (*pb.GetAgentInfoResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db == nil {
		return &pb.GetAgentInfoResponse{Success: false, Message: "Database not available"}, nil
	}

	// Get agent from database
	agent, err := s.db.GetAgent(req.AgentName)
	if err != nil {
		return &pb.GetAgentInfoResponse{Success: false, Message: fmt.Sprintf("Agent not found: %s", req.AgentName)}, nil
	}

	return &pb.GetAgentInfoResponse{
		Success: true,
		Message: "Agent info retrieved successfully",
		AgentInfo: &pb.AgentInfo{
			AgentName:         agent.Name,
			AgentAddress:      agent.Address,
			LastHeartbeat:     agent.LastHeartbeat,
			Status:            agent.Status,
			LastInfoCollected: agent.LastInfoCollected,
			SystemInfoJson:    agent.SystemInfo,
			Version:           agent.Version,
		},
	}, nil
}

// GetAgentAddress retrieves the address of an agent by name (new method for taskrunner)
func (s *agentRegistryServer) GetAgentAddress(agentName string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.db != nil {
		return s.db.GetAgentAddress(agentName)
	}
	
	return "", fmt.Errorf("database not available")
}

// Start starts the agent registry server.
func (s *agentRegistryServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	pterm.Warning.Println("Starting master in insecure mode.")
	
	s.grpcServer = grpc.NewServer()
	pb.RegisterAgentRegistryServer(s.grpcServer, s)
	pterm.Info.Printf("Agent registry listening at %v\n", lis.Addr())
	return s.grpcServer.Serve(lis)
}

// Stop stops the agent registry server.
func (s *agentRegistryServer) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	
	// Close database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			pterm.Error.Printf("Failed to close agent database: %v\n", err)
		} else {
			pterm.Info.Println("Agent database closed successfully")
		}
	}
}