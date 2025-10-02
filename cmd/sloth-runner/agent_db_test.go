package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*AgentDB, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_agents.db")
	
	db, err := NewAgentDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	return db, dbPath
}

func TestNewAgentDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	
	db, err := NewAgentDB(dbPath)
	if err != nil {
		t.Fatalf("NewAgentDB failed: %v", err)
	}
	defer db.Close()
	
	if db == nil {
		t.Fatal("Expected non-nil database")
	}
	
	// Verify database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("Database file was not created")
	}
}

func TestNewAgentDB_InvalidPath(t *testing.T) {
	// Try to create DB in a path that doesn't exist and can't be created
	dbPath := "/root/restricted/test.db"
	
	_, err := NewAgentDB(dbPath)
	if err == nil {
		t.Fatal("Expected error when creating DB in restricted path")
	}
}

func TestRegisterAgent(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	err := db.RegisterAgent("test-agent", "localhost:8080")
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}
	
	// Verify agent was registered
	agent, err := db.GetAgent("test-agent")
	if err != nil {
		t.Fatalf("GetAgent failed: %v", err)
	}
	
	if agent.Name != "test-agent" {
		t.Errorf("Expected name 'test-agent', got '%s'", agent.Name)
	}
	if agent.Address != "localhost:8080" {
		t.Errorf("Expected address 'localhost:8080', got '%s'", agent.Address)
	}
	if agent.Status != "Active" {
		t.Errorf("Expected status 'Active', got '%s'", agent.Status)
	}
}

func TestRegisterAgent_Update(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register agent
	err := db.RegisterAgent("test-agent", "localhost:8080")
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}
	
	// Get initial registration time
	agent1, _ := db.GetAgent("test-agent")
	
	time.Sleep(1100 * time.Millisecond) // Sleep more than 1 second to ensure different timestamps
	
	// Update agent with new address
	err = db.RegisterAgent("test-agent", "localhost:9090")
	if err != nil {
		t.Fatalf("RegisterAgent update failed: %v", err)
	}
	
	// Verify agent was updated
	agent2, err := db.GetAgent("test-agent")
	if err != nil {
		t.Fatalf("GetAgent failed: %v", err)
	}
	
	if agent2.Address != "localhost:9090" {
		t.Errorf("Expected address 'localhost:9090', got '%s'", agent2.Address)
	}
	
	// Verify registered_at remained the same
	if agent1.RegisteredAt != agent2.RegisteredAt {
		t.Error("RegisteredAt should not change on update")
	}
	
	// Verify updated_at changed
	if agent1.UpdatedAt >= agent2.UpdatedAt {
		t.Errorf("UpdatedAt should change on update. Old: %d, New: %d", agent1.UpdatedAt, agent2.UpdatedAt)
	}
}

func TestUpdateHeartbeat(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register agent
	db.RegisterAgent("test-agent", "localhost:8080")
	
	agent1, _ := db.GetAgent("test-agent")
	initialHeartbeat := agent1.LastHeartbeat
	
	time.Sleep(1100 * time.Millisecond) // Sleep more than 1 second to ensure different timestamps
	
	// Update heartbeat
	err := db.UpdateHeartbeat("test-agent")
	if err != nil {
		t.Fatalf("UpdateHeartbeat failed: %v", err)
	}
	
	agent2, _ := db.GetAgent("test-agent")
	
	if agent2.LastHeartbeat <= initialHeartbeat {
		t.Errorf("LastHeartbeat should have been updated. Old: %d, New: %d", initialHeartbeat, agent2.LastHeartbeat)
	}
}

func TestUpdateHeartbeat_NonExistentAgent(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	err := db.UpdateHeartbeat("non-existent-agent")
	if err == nil {
		t.Fatal("Expected error when updating heartbeat for non-existent agent")
	}
}

func TestGetAgent(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register multiple agents
	db.RegisterAgent("agent1", "host1:8080")
	db.RegisterAgent("agent2", "host2:8080")
	
	// Get specific agent
	agent, err := db.GetAgent("agent1")
	if err != nil {
		t.Fatalf("GetAgent failed: %v", err)
	}
	
	if agent.Name != "agent1" {
		t.Errorf("Expected name 'agent1', got '%s'", agent.Name)
	}
	if agent.Address != "host1:8080" {
		t.Errorf("Expected address 'host1:8080', got '%s'", agent.Address)
	}
}

func TestGetAgent_NotFound(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	_, err := db.GetAgent("non-existent-agent")
	if err == nil {
		t.Fatal("Expected error when getting non-existent agent")
	}
}

func TestListAgents(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register multiple agents
	db.RegisterAgent("agent1", "host1:8080")
	db.RegisterAgent("agent2", "host2:8080")
	db.RegisterAgent("agent3", "host3:8080")
	
	agents, err := db.ListAgents()
	if err != nil {
		t.Fatalf("ListAgents failed: %v", err)
	}
	
	if len(agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(agents))
	}
	
	// Verify agents are sorted by name
	if agents[0].Name != "agent1" || agents[1].Name != "agent2" || agents[2].Name != "agent3" {
		t.Error("Agents not sorted correctly")
	}
}

func TestListAgents_StatusDetermination(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register agent with recent heartbeat
	db.RegisterAgent("active-agent", "host1:8080")
	time.Sleep(100 * time.Millisecond)
	db.UpdateHeartbeat("active-agent")
	
	// Register agent and set old heartbeat to make it inactive
	db.RegisterAgent("inactive-agent", "host2:8080")
	// Set very old heartbeat (more than 60 seconds ago)
	oldTime := time.Now().Unix() - 120
	db.db.Exec("UPDATE agents SET last_heartbeat = ? WHERE name = ?", oldTime, "inactive-agent")
	
	agents, err := db.ListAgents()
	if err != nil {
		t.Fatalf("ListAgents failed: %v", err)
	}
	
	// Find agents and check status
	var activeAgent, inactiveAgent *AgentRecord
	for _, agent := range agents {
		if agent.Name == "active-agent" {
			activeAgent = agent
		} else if agent.Name == "inactive-agent" {
			inactiveAgent = agent
		}
	}
	
	if activeAgent == nil || inactiveAgent == nil {
		t.Fatal("Could not find test agents")
	}
	
	if activeAgent.Status != "Active" {
		t.Errorf("Expected 'Active', got '%s'", activeAgent.Status)
	}
	
	if inactiveAgent.Status != "Inactive" {
		t.Errorf("Expected 'Inactive', got '%s'", inactiveAgent.Status)
	}
}

func TestGetAgentAddress(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register agent and update heartbeat
	db.RegisterAgent("test-agent", "localhost:8080")
	db.UpdateHeartbeat("test-agent")
	
	address, err := db.GetAgentAddress("test-agent")
	if err != nil {
		t.Fatalf("GetAgentAddress failed: %v", err)
	}
	
	if address != "localhost:8080" {
		t.Errorf("Expected 'localhost:8080', got '%s'", address)
	}
}

func TestGetAgentAddress_InactiveAgent(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register agent and set old heartbeat (inactive)
	db.RegisterAgent("inactive-agent", "localhost:8080")
	// Set very old heartbeat (more than 60 seconds ago)
	oldTime := time.Now().Unix() - 120
	db.db.Exec("UPDATE agents SET last_heartbeat = ? WHERE name = ?", oldTime, "inactive-agent")
	
	_, err := db.GetAgentAddress("inactive-agent")
	if err == nil {
		t.Fatal("Expected error when getting address of inactive agent")
	}
}

func TestRemoveAgent(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register agent
	db.RegisterAgent("test-agent", "localhost:8080")
	
	// Verify agent exists
	_, err := db.GetAgent("test-agent")
	if err != nil {
		t.Fatal("Agent should exist before removal")
	}
	
	// Remove agent
	err = db.RemoveAgent("test-agent")
	if err != nil {
		t.Fatalf("RemoveAgent failed: %v", err)
	}
	
	// Verify agent was removed
	_, err = db.GetAgent("test-agent")
	if err == nil {
		t.Fatal("Agent should not exist after removal")
	}
}

func TestRemoveAgent_NonExistent(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	err := db.RemoveAgent("non-existent-agent")
	if err == nil {
		t.Fatal("Expected error when removing non-existent agent")
	}
}

func TestUnregisterAgent(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register agent
	db.RegisterAgent("test-agent", "localhost:8080")
	
	// Unregister agent
	err := db.UnregisterAgent("test-agent")
	if err != nil {
		t.Fatalf("UnregisterAgent failed: %v", err)
	}
	
	// Verify agent was removed
	_, err = db.GetAgent("test-agent")
	if err == nil {
		t.Fatal("Agent should not exist after unregistration")
	}
}

func TestCleanupInactiveAgents(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register multiple agents
	db.RegisterAgent("active-agent", "host1:8080")
	db.RegisterAgent("inactive-agent1", "host2:8080")
	db.RegisterAgent("inactive-agent2", "host3:8080")
	
	// Update heartbeat only for active agent
	db.UpdateHeartbeat("active-agent")
	
	// Set old heartbeat for inactive agents (simulate old activity)
	oldTime := time.Now().Unix() - (25 * 3600) // 25 hours ago
	db.db.Exec("UPDATE agents SET last_heartbeat = ? WHERE name IN (?, ?)", 
		oldTime, "inactive-agent1", "inactive-agent2")
	
	// Cleanup agents inactive for more than 24 hours
	count, err := db.CleanupInactiveAgents(24)
	if err != nil {
		t.Fatalf("CleanupInactiveAgents failed: %v", err)
	}
	
	if count != 2 {
		t.Errorf("Expected 2 agents to be cleaned up, got %d", count)
	}
	
	// Verify active agent still exists
	_, err = db.GetAgent("active-agent")
	if err != nil {
		t.Fatal("Active agent should still exist")
	}
	
	// Verify inactive agents were removed
	_, err = db.GetAgent("inactive-agent1")
	if err == nil {
		t.Fatal("inactive-agent1 should have been removed")
	}
	
	_, err = db.GetAgent("inactive-agent2")
	if err == nil {
		t.Fatal("inactive-agent2 should have been removed")
	}
}

func TestGetStats(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register multiple agents
	db.RegisterAgent("active1", "host1:8080")
	db.RegisterAgent("active2", "host2:8080")
	db.RegisterAgent("inactive1", "host3:8080")
	
	time.Sleep(100 * time.Millisecond)
	
	// Update heartbeat for active agents
	db.UpdateHeartbeat("active1")
	db.UpdateHeartbeat("active2")
	
	// Set old heartbeat for inactive agent
	oldTime := time.Now().Unix() - 120
	db.db.Exec("UPDATE agents SET last_heartbeat = ? WHERE name = ?", oldTime, "inactive1")
	
	// Get stats
	stats, err := db.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	
	totalAgents, ok := stats["total_agents"].(int)
	if !ok {
		t.Fatal("total_agents not found in stats")
	}
	if totalAgents != 3 {
		t.Errorf("Expected 3 total agents, got %d", totalAgents)
	}
	
	activeAgents, ok := stats["active_agents"].(int)
	if !ok {
		t.Fatal("active_agents not found in stats")
	}
	if activeAgents != 2 {
		t.Errorf("Expected 2 active agents, got %d", activeAgents)
	}
	
	inactiveAgents, ok := stats["inactive_agents"].(int)
	if !ok {
		t.Fatal("inactive_agents not found in stats")
	}
	if inactiveAgents != 1 {
		t.Errorf("Expected 1 inactive agent, got %d", inactiveAgents)
	}
}

func TestClose(t *testing.T) {
	db, _ := setupTestDB(t)
	
	err := db.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	
	// Try to close again
	err = db.Close()
	if err != nil {
		t.Fatalf("Second close should not error: %v", err)
	}
}

func TestConcurrentOperations(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Register initial agents
	for i := 0; i < 10; i++ {
		db.RegisterAgent(string(rune('a'+i)), "localhost:8080")
	}
	
	// Perform concurrent operations
	done := make(chan bool, 30)
	
	// Concurrent heartbeat updates
	for i := 0; i < 10; i++ {
		go func(name string) {
			db.UpdateHeartbeat(name)
			done <- true
		}(string(rune('a' + i)))
	}
	
	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func(name string) {
			db.GetAgent(name)
			done <- true
		}(string(rune('a' + i)))
	}
	
	// Concurrent list operations
	for i := 0; i < 10; i++ {
		go func() {
			db.ListAgents()
			done <- true
		}()
	}
	
	// Wait for all operations to complete
	for i := 0; i < 30; i++ {
		<-done
	}
}
