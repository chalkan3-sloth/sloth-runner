package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
)

func TestNewAgentRegistryServer(t *testing.T) {
	// Create temp directory for test database
	tmpDir, err := os.MkdirTemp("", "agent-registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	server := newAgentRegistryServer()
	if server == nil {
		t.Fatal("newAgentRegistryServer returned nil")
	}
}

func TestAgentRegistryRegisterAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "agents.db")
	db, err := NewAgentDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	server := &agentRegistryServer{
		db: db,
	}

	ctx := context.Background()
	req := &pb.RegisterAgentRequest{
		AgentName:    "test-agent",
		AgentAddress: "localhost:50051",
	}

	resp, err := server.RegisterAgent(ctx, req)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success=true, got success=%v", resp.Success)
	}

	if resp.Message != "Agent registered successfully" {
		t.Errorf("Expected message='Agent registered successfully', got message='%s'", resp.Message)
	}
}

func TestAgentRegistryListAgents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "agents.db")
	db, err := NewAgentDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Register an agent first
	err = db.RegisterAgent("test-agent", "localhost:50051")
	if err != nil {
		t.Fatalf("Failed to register agent: %v", err)
	}

	server := &agentRegistryServer{
		db: db,
	}

	ctx := context.Background()
	req := &pb.ListAgentsRequest{}

	resp, err := server.ListAgents(ctx, req)
	if err != nil {
		t.Fatalf("ListAgents failed: %v", err)
	}

	if len(resp.Agents) == 0 {
		t.Error("Expected at least one agent, got 0")
	}

	if resp.Agents[0].AgentName != "test-agent" {
		t.Errorf("Expected agent name='test-agent', got name='%s'", resp.Agents[0].AgentName)
	}
}

func TestHeartbeat(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "agents.db")
	db, err := NewAgentDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Register an agent first
	err = db.RegisterAgent("test-agent", "localhost:50051")
	if err != nil {
		t.Fatalf("Failed to register agent: %v", err)
	}

	server := &agentRegistryServer{
		db: db,
	}

	ctx := context.Background()
	req := &pb.HeartbeatRequest{
		AgentName: "test-agent",
	}

	resp, err := server.Heartbeat(ctx, req)
	if err != nil {
		t.Fatalf("Heartbeat failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success=true, got success=%v", resp.Success)
	}

	if resp.Message != "Heartbeat received" {
		t.Errorf("Expected message='Heartbeat received', got message='%s'", resp.Message)
	}
}

func TestHeartbeatNonExistentAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "agents.db")
	db, err := NewAgentDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	server := &agentRegistryServer{
		db: db,
	}

	ctx := context.Background()
	req := &pb.HeartbeatRequest{
		AgentName: "non-existent-agent",
	}

	resp, err := server.Heartbeat(ctx, req)
	if err != nil {
		t.Fatalf("Heartbeat failed: %v", err)
	}

	if resp.Success {
		t.Errorf("Expected success=false for non-existent agent, got success=%v", resp.Success)
	}
}

func TestRegisterAgentWithoutDB(t *testing.T) {
	server := &agentRegistryServer{
		db: nil, // No database
	}

	ctx := context.Background()
	req := &pb.RegisterAgentRequest{
		AgentName:    "test-agent",
		AgentAddress: "localhost:50051",
	}

	resp, err := server.RegisterAgent(ctx, req)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	// Should succeed even without database (fallback mode)
	if !resp.Success {
		t.Error("Expected success when falling back to in-memory storage")
	}
}

func TestListAgentsWithoutDB(t *testing.T) {
	server := &agentRegistryServer{
		db: nil, // No database
	}

	ctx := context.Background()
	req := &pb.ListAgentsRequest{}

	resp, err := server.ListAgents(ctx, req)
	if err != nil {
		t.Fatalf("ListAgents failed: %v", err)
	}

	// Should return empty list
	if len(resp.Agents) != 0 {
		t.Errorf("Expected 0 agents when database is not available, got %d", len(resp.Agents))
	}
}

func TestHeartbeatWithoutDB(t *testing.T) {
	server := &agentRegistryServer{
		db: nil, // No database
	}

	ctx := context.Background()
	req := &pb.HeartbeatRequest{
		AgentName: "test-agent",
	}

	resp, err := server.Heartbeat(ctx, req)
	if err != nil {
		t.Fatalf("Heartbeat failed: %v", err)
	}

	if resp.Success {
		t.Error("Expected failure when database is not available")
	}
}

func TestConcurrentRegistrations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "agents.db")
	db, err := NewAgentDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	server := &agentRegistryServer{
		db: db,
	}

	// Register multiple agents concurrently
	numAgents := 10
	done := make(chan bool, numAgents)

	for i := 0; i < numAgents; i++ {
		go func(id int) {
			ctx := context.Background()
			req := &pb.RegisterAgentRequest{
				AgentName:    fmt.Sprintf("agent-%d", id),
				AgentAddress: fmt.Sprintf("localhost:5005%d", id),
			}
			
			_, err := server.RegisterAgent(ctx, req)
			if err != nil {
				t.Errorf("RegisterAgent failed for agent-%d: %v", id, err)
			}
			done <- true
		}(i)
	}

	// Wait for all registrations
	timeout := time.After(5 * time.Second)
	for i := 0; i < numAgents; i++ {
		select {
		case <-done:
			// Success
		case <-timeout:
			t.Fatal("Timeout waiting for concurrent registrations")
		}
	}

	// Verify all agents were registered
	ctx := context.Background()
	listResp, err := server.ListAgents(ctx, &pb.ListAgentsRequest{})
	if err != nil {
		t.Fatalf("ListAgents failed: %v", err)
	}

	if len(listResp.Agents) != numAgents {
		t.Errorf("Expected %d agents, got %d", numAgents, len(listResp.Agents))
	}
}
