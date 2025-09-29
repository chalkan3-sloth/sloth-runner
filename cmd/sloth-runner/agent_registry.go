package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// agentRegistryServer implements the AgentRegistry service.
type agentRegistryServer struct {
	pb.UnimplementedAgentRegistryServer
	mu          sync.Mutex
	agents      map[string]*pb.AgentInfo
	grpcServer  *grpc.Server
	tlsCertFile string
	tlsKeyFile  string
	tlsCaFile   string
}

// newAgentRegistryServer creates a new agentRegistryServer.
func newAgentRegistryServer(tlsCertFile, tlsKeyFile, tlsCaFile string) *agentRegistryServer {
	return &agentRegistryServer{
		agents:      make(map[string]*pb.AgentInfo),
		tlsCertFile: tlsCertFile,
		tlsKeyFile:  tlsKeyFile,
		tlsCaFile:   tlsCaFile,
	}
}

// RegisterAgent registers a new agent.
func (s *agentRegistryServer) RegisterAgent(ctx context.Context, req *pb.RegisterAgentRequest) (*pb.RegisterAgentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pterm.Success.Printf("Agent registered: %s at %s\n", req.AgentName, req.AgentAddress)
	s.agents[req.AgentName] = &pb.AgentInfo{
		AgentName:    req.AgentName,
		AgentAddress: req.AgentAddress,
	}

	return &pb.RegisterAgentResponse{Success: true, Message: "Agent registered successfully"}, nil
}

// ListAgents lists all registered agents.
func (s *agentRegistryServer) ListAgents(ctx context.Context, req *pb.ListAgentsRequest) (*pb.ListAgentsResponse, error) {
	pterm.Info.Println("Listing registered agents")
	s.mu.Lock()
	defer s.mu.Unlock()

	var agents []*pb.AgentInfo
	for _, agent := range s.agents {
		status := "Inactive"
		if agent.LastHeartbeat > 0 && time.Now().Unix()-agent.LastHeartbeat < 60 { // Agent considered active if heartbeat within last 60 seconds
			status = "Active"
		}
		agents = append(agents, &pb.AgentInfo{
			AgentName:     agent.AgentName,
			AgentAddress:  agent.AgentAddress,
			LastHeartbeat: agent.LastHeartbeat,
			Status:        status,
		})
	}

	return &pb.ListAgentsResponse{Agents: agents}, nil
}

// Heartbeat updates the last heartbeat timestamp for an agent.
func (s *agentRegistryServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if agent, ok := s.agents[req.AgentName]; ok {
		agent.LastHeartbeat = time.Now().Unix()
		pterm.Debug.Printf("Heartbeat received from agent: %s\n", req.AgentName)
		return &pb.HeartbeatResponse{Success: true, Message: "Heartbeat received"}, nil
	}
	return &pb.HeartbeatResponse{Success: false, Message: "Agent not found"}, nil
}

// ExecuteCommand executes a command on a remote agent and streams the output back to the client.
func (s *agentRegistryServer) ExecuteCommand(req *pb.ExecuteCommandRequest, stream pb.AgentRegistry_ExecuteCommandServer) error {
	s.mu.Lock()
	agent, ok := s.agents[req.AgentName]
	s.mu.Unlock()

	if !ok {
		return fmt.Errorf("agent not found: %s", req.AgentName)
	}

	dialOpts, err := getDialOptions(s.tlsCertFile, s.tlsKeyFile, s.tlsCaFile) // No client certs for master to agent communication
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(agent.AgentAddress, dialOpts...)
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
	s.mu.Lock()
	agent, ok := s.agents[req.AgentName]
	s.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("agent not found: %s", req.AgentName)
	}

	dialOpts, err := getDialOptions(s.tlsCertFile, s.tlsKeyFile, s.tlsCaFile) // No client certs for master to agent communication
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(agent.AgentAddress, dialOpts...)
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

// Start starts the agent registry server.
func (s *agentRegistryServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if s.tlsCertFile != "" && s.tlsKeyFile != "" && s.tlsCaFile != "" {
		serverCert, err := tls.LoadX509KeyPair(s.tlsCertFile, s.tlsKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load server certificate: %v", err)
		}

		caCert, err := os.ReadFile(s.tlsCaFile)
		if err != nil {
			return fmt.Errorf("failed to read CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return fmt.Errorf("failed to append CA certificate")
		}

		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		})
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		pterm.Warning.Println("Starting master in insecure mode. TLS certificates not provided.")
	}

	s.grpcServer = grpc.NewServer(opts...)
	pb.RegisterAgentRegistryServer(s.grpcServer, s)
	pterm.Info.Printf("Agent registry listening at %v\n", lis.Addr())
	return s.grpcServer.Serve(lis)
}

// Stop stops the agent registry server.
func (s *agentRegistryServer) Stop() {
	s.grpcServer.GracefulStop()
}

func getDialOptions(tlsCertFile, tlsKeyFile, tlsCaFile string) ([]grpc.DialOption, error) {
	if tlsCaFile != "" {
		caCert, err := os.ReadFile(tlsCaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		var clientCerts []tls.Certificate
		if tlsCertFile != "" && tlsKeyFile != "" {
			clientCert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load client certificate: %v", err)
			}
			clientCerts = append(clientCerts, clientCert)
		}

		creds := credentials.NewTLS(&tls.Config{
			Certificates: clientCerts,
			RootCAs:      caCertPool,
			ServerName:   "localhost", // Explicitly set ServerName
		})
		return []grpc.DialOption{grpc.WithTransportCredentials(creds)}, nil
	}

	pterm.Warning.Println("Connecting in insecure mode. TLS certificates not provided.")
	return []grpc.DialOption{grpc.WithInsecure()}, nil
}