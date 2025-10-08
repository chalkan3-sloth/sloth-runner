package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"path/filepath"
	"sync"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/chalkan3-sloth/sloth-runner/internal/metrics"
	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// agentRegistryServer implements the AgentRegistry service.
type agentRegistryServer struct {
	pb.UnimplementedAgentRegistryServer
	mu               sync.RWMutex
	db               *AgentDB
	grpcServer       *grpc.Server
	dispatcher       *hooks.Dispatcher
	metricsDB        *metrics.MetricsDB
	metricsCollector *metrics.Collector
}

// newAgentRegistryServer creates a new agentRegistryServer.
func newAgentRegistryServer() *agentRegistryServer {
	// Create SQLite database for agents
	dbPath := filepath.Join(GetSlothRunnerDataDir(), "agents.db")
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

	// Initialize global hook dispatcher and wire up event system
	var dispatcher *hooks.Dispatcher
	if err := InitializeHookSystem(); err != nil {
		pterm.Error.Printf("Failed to initialize hook system: %v\n", err)
		pterm.Info.Println("Hooks will be disabled")
	} else {
		dispatcher = hooks.GetGlobalDispatcher()
		pterm.Success.Println("Hook system initialized")
	}

	// Initialize metrics database
	metricsDBPath := config.GetMetricsDBPath()
	if err := config.EnsureDataDir(); err != nil {
		pterm.Error.Printf("Failed to create metrics directory: %v\n", err)
	}

	slog.Info("Initializing metrics database", "path", metricsDBPath)
	metricsDB, err := metrics.NewMetricsDB(metricsDBPath)
	if err != nil {
		pterm.Error.Printf("Failed to initialize metrics database: %v\n", err)
		metricsDB = nil
	} else {
		pterm.Success.Printf("Metrics database initialized at: %s\n", metricsDBPath)
	}

	// Initialize metrics collector
	var metricsCollector *metrics.Collector
	if metricsDB != nil && db != nil {
		agentClient := services.NewAgentClient()
		metricsCollector = metrics.NewCollector(metrics.CollectorConfig{
			MetricsDB:     metricsDB,
			AgentClient:   agentClient,
			Interval:      30 * time.Second,
			RetentionDays: 7,
		})

		// Start metrics collector with function to get agent list
		ctx := context.Background()
		if err := metricsCollector.Start(ctx, func() []metrics.AgentInfo {
			agents, err := db.ListAgents()
			if err != nil {
				slog.Error("Failed to get agent list for metrics collection", "error", err)
				return nil
			}

			agentInfos := make([]metrics.AgentInfo, 0, len(agents))
			for _, agent := range agents {
				agentInfos = append(agentInfos, metrics.AgentInfo{
					Name:    agent.Name,
					Address: agent.Address,
				})
			}
			return agentInfos
		}); err != nil {
			pterm.Error.Printf("Failed to start metrics collector: %v\n", err)
		} else {
			pterm.Success.Println("Metrics collector started (interval: 30s, retention: 7 days)")
		}
	}

	return &agentRegistryServer{
		db:               db,
		dispatcher:       dispatcher,
		metricsDB:        metricsDB,
		metricsCollector: metricsCollector,
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

	// Dispatch agent registered event
	if s.dispatcher != nil {
		agent := &hooks.AgentEvent{
			Name:    req.AgentName,
			Address: req.AgentAddress,
		}
		if err := s.dispatcher.DispatchAgentRegistered(agent); err != nil {
			pterm.Debug.Printf("Failed to dispatch agent registered event: %v\n", err)
		}
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

	// Get agent info before removing (for event dispatch)
	agentInfo, _ := s.db.GetAgent(req.AgentName)

	// Remove agent from database
	if err := s.db.UnregisterAgent(req.AgentName); err != nil {
		pterm.Error.Printf("Failed to unregister agent %s: %v\n", req.AgentName, err)
		return &pb.UnregisterAgentResponse{Success: false, Message: fmt.Sprintf("Failed to unregister agent: %v", err)}, nil
	}

	// Dispatch agent disconnected event
	if s.dispatcher != nil && agentInfo != nil {
		agent := &hooks.AgentEvent{
			Name:    agentInfo.Name,
			Address: agentInfo.Address,
			Version: agentInfo.Version,
		}
		if err := s.dispatcher.DispatchAgentDisconnected(agent); err != nil {
			pterm.Debug.Printf("Failed to dispatch agent disconnected event: %v\n", err)
		}
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

// SendEvent receives an event from an agent and dispatches it through the global event system
func (s *agentRegistryServer) SendEvent(ctx context.Context, req *pb.SendEventRequest) (*pb.SendEventResponse, error) {
	if req.Event == nil {
		return &pb.SendEventResponse{
			Success: false,
			Message: "Event data is required",
		}, nil
	}

	// Check if dispatcher is available
	if s.dispatcher == nil {
		slog.Warn("Event dispatcher not available, event will be dropped",
			"event_type", req.Event.EventType,
			"agent", req.Event.AgentName)
		return &pb.SendEventResponse{
			Success: false,
			Message: "Event dispatcher not initialized on master",
			EventId: req.Event.EventId,
		}, nil
	}

	// Convert protobuf EventData to internal hooks.Event format
	event := &hooks.Event{
		ID:        req.Event.EventId,
		Type:      hooks.EventType(req.Event.EventType),
		Timestamp: time.Unix(req.Event.Timestamp, 0),
		Stack:     req.Event.Stack,
		Agent:     req.Event.AgentName,
		RunID:     req.Event.RunId,
		Status:    hooks.EventStatusPending,
		CreatedAt: time.Now(),
	}

	// Handle event data - prefer JSON format if available
	if req.Event.DataJson != "" {
		// Parse JSON string into Data map
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(req.Event.DataJson), &data); err == nil {
			event.Data = data
		}
	} else if len(req.Event.Data) > 0 {
		// Convert protobuf map to interface{} map
		event.Data = make(map[string]interface{})
		for k, v := range req.Event.Data {
			event.Data[k] = v
		}
	}

	// Dispatch the event through the global dispatcher
	if err := s.dispatcher.Dispatch(event); err != nil {
		slog.Error("Failed to dispatch agent event",
			"event_id", event.ID,
			"event_type", event.Type,
			"agent", event.Agent,
			"error", err)
		return &pb.SendEventResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to dispatch event: %v", err),
			EventId: event.ID,
		}, nil
	}

	slog.Debug("Event received from agent and dispatched",
		"event_id", event.ID,
		"event_type", event.Type,
		"agent", event.Agent,
		"stack", event.Stack,
		"run_id", event.RunID)

	return &pb.SendEventResponse{
		Success: true,
		Message: "Event received and dispatched successfully",
		EventId: event.ID,
	}, nil
}

// SendEventBatch receives multiple events from an agent and dispatches them through the global event system
func (s *agentRegistryServer) SendEventBatch(ctx context.Context, req *pb.SendEventBatchRequest) (*pb.SendEventBatchResponse, error) {
	if len(req.Events) == 0 {
		return &pb.SendEventBatchResponse{
			Success:          false,
			Message:          "No events provided",
			EventsReceived:   0,
			EventsProcessed:  0,
			FailedEventIds:   []string{},
		}, nil
	}

	// Check if dispatcher is available
	if s.dispatcher == nil {
		slog.Warn("Event dispatcher not available, batch events will be dropped",
			"event_count", len(req.Events))

		failedIds := make([]string, len(req.Events))
		for i, evt := range req.Events {
			failedIds[i] = evt.EventId
		}

		return &pb.SendEventBatchResponse{
			Success:          false,
			Message:          "Event dispatcher not initialized on master",
			EventsReceived:   int32(len(req.Events)),
			EventsProcessed:  0,
			FailedEventIds:   failedIds,
		}, nil
	}

	var processedCount int32
	var failedEventIds []string

	slog.Info("Processing event batch from agent",
		"batch_size", len(req.Events),
		"requested_batch_size", req.BatchSize)

	// Process each event in the batch
	for _, pbEvent := range req.Events {
		// Convert protobuf EventData to internal hooks.Event format
		event := &hooks.Event{
			ID:        pbEvent.EventId,
			Type:      hooks.EventType(pbEvent.EventType),
			Timestamp: time.Unix(pbEvent.Timestamp, 0),
			Stack:     pbEvent.Stack,
			Agent:     pbEvent.AgentName,
			RunID:     pbEvent.RunId,
			Status:    hooks.EventStatusPending,
			CreatedAt: time.Now(),
		}

		// Handle event data
		if pbEvent.DataJson != "" {
			// Parse JSON string into Data map
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(pbEvent.DataJson), &data); err == nil {
				event.Data = data
			}
		} else if len(pbEvent.Data) > 0 {
			event.Data = make(map[string]interface{})
			for k, v := range pbEvent.Data {
				event.Data[k] = v
			}
		}

		// Dispatch the event
		if err := s.dispatcher.Dispatch(event); err != nil {
			slog.Error("Failed to dispatch event in batch",
				"event_id", event.ID,
				"event_type", event.Type,
				"agent", event.Agent,
				"error", err)
			failedEventIds = append(failedEventIds, event.ID)
		} else {
			processedCount++
		}
	}

	successRate := float64(processedCount) / float64(len(req.Events)) * 100
	slog.Info("Event batch processed",
		"total", len(req.Events),
		"processed", processedCount,
		"failed", len(failedEventIds),
		"success_rate", fmt.Sprintf("%.1f%%", successRate))

	return &pb.SendEventBatchResponse{
		Success:          processedCount > 0,
		Message:          fmt.Sprintf("Processed %d/%d events successfully", processedCount, len(req.Events)),
		EventsReceived:   int32(len(req.Events)),
		EventsProcessed:  processedCount,
		FailedEventIds:   failedEventIds,
	}, nil
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