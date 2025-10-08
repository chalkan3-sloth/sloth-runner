package client

import (
	"context"
	"fmt"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RegisterWatcherOnAgent registers a watcher on a remote agent via gRPC
func RegisterWatcherOnAgent(ctx context.Context, agentAddr string, config *pb.WatcherConfig) (*pb.RegisterWatcherResponse, error) {
	conn, err := grpc.NewClient(agentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	client := pb.NewAgentClient(conn)

	req := &pb.RegisterWatcherRequest{
		Config: config,
	}

	return client.RegisterWatcher(ctx, req)
}

// ListWatchersOnAgent lists all watchers on a remote agent via gRPC
func ListWatchersOnAgent(ctx context.Context, agentAddr string) (*pb.ListWatchersResponse, error) {
	conn, err := grpc.NewClient(agentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	client := pb.NewAgentClient(conn)

	req := &pb.ListWatchersRequest{}
	return client.ListWatchers(ctx, req)
}

// RemoveWatcherFromAgent removes a watcher from a remote agent via gRPC
func RemoveWatcherFromAgent(ctx context.Context, agentAddr string, watcherID string) (*pb.RemoveWatcherResponse, error) {
	conn, err := grpc.NewClient(agentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	client := pb.NewAgentClient(conn)

	req := &pb.RemoveWatcherRequest{
		WatcherId: watcherID,
	}

	return client.RemoveWatcher(ctx, req)
}
