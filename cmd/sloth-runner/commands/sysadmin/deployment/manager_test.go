package deployment

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewDeploymentManager(t *testing.T) {
	manager := NewDeploymentManager()
	if manager == nil {
		t.Fatal("NewDeploymentManager() returned nil")
	}

	_, ok := manager.(*SystemDeployment)
	if !ok {
		t.Error("NewDeploymentManager() did not return *SystemDeployment")
	}
}

func TestDeploy(t *testing.T) {
	manager := NewDeploymentManager()

	options := DeployOptions{
		Version:      "v1.0.0",
		Agents:       []string{"agent1", "agent2"},
		Strategy:     StrategyDirect,
		DryRun:       false,
		HealthCheck:  true,
		BackupBefore: true,
	}

	result, err := manager.Deploy(options)
	if err != nil {
		t.Fatalf("Deploy() failed: %v", err)
	}

	if result == nil {
		t.Fatal("Deploy() returned nil result")
	}

	if !result.Success {
		t.Errorf("Deploy() failed: %s", result.Message)
	}

	if result.Version != "v1.0.0" {
		t.Errorf("Expected version v1.0.0, got %s", result.Version)
	}

	if len(result.AgentsUpdated) != 2 {
		t.Errorf("Expected 2 agents updated, got %d", len(result.AgentsUpdated))
	}

	if len(result.AgentsFailed) != 0 {
		t.Errorf("Expected 0 agents failed, got %d", len(result.AgentsFailed))
	}

	if result.Duration == 0 {
		t.Error("Duration is zero")
	}
}

func TestDeployDryRun(t *testing.T) {
	manager := NewDeploymentManager()

	options := DeployOptions{
		Version:  "v2.0.0",
		Agents:   []string{"agent1", "agent2", "agent3"},
		Strategy: StrategyRolling,
		DryRun:   true,
	}

	result, err := manager.Deploy(options)
	if err != nil {
		t.Fatalf("Deploy() dry run failed: %v", err)
	}

	if !result.Success {
		t.Error("Dry run should always succeed")
	}

	if len(result.AgentsUpdated) != 0 {
		t.Error("Dry run should not update any agents")
	}

	if result.Message == "" {
		t.Error("Dry run message is empty")
	}
}

func TestDeployStrategies(t *testing.T) {
	manager := NewDeploymentManager()

	strategies := []DeploymentStrategy{
		StrategyDirect,
		StrategyRolling,
		StrategyCanary,
		StrategyBlueGreen,
	}

	for _, strategy := range strategies {
		t.Run(string(strategy), func(t *testing.T) {
			options := DeployOptions{
				Version:  "v1.0.0",
				Agents:   []string{"agent1"},
				Strategy: strategy,
			}

			result, err := manager.Deploy(options)
			if err != nil {
				t.Fatalf("Deploy() with strategy %s failed: %v", strategy, err)
			}

			if !result.Success {
				t.Errorf("Deploy() with strategy %s failed", strategy)
			}
		})
	}
}

func TestRollback(t *testing.T) {
	// Cria arquivo temporário para histórico
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "deployment-history.json")

	manager := &SystemDeployment{
		historyFile: historyFile,
	}

	// Faz dois deployments para ter histórico
	options1 := DeployOptions{
		Version:  "v1.0.0",
		Agents:   []string{"agent1"},
		Strategy: StrategyDirect,
	}
	manager.Deploy(options1)

	options2 := DeployOptions{
		Version:  "v2.0.0",
		Agents:   []string{"agent1"},
		Strategy: StrategyDirect,
	}
	manager.Deploy(options2)

	// Testa rollback
	rollbackOptions := RollbackOptions{
		Agent:  "agent1",
		Verify: true,
	}

	result, err := manager.Rollback(rollbackOptions)
	if err != nil {
		t.Fatalf("Rollback() failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Rollback() failed: %s", result.Message)
	}

	if result.Version != "v1.0.0" {
		t.Errorf("Expected rollback to v1.0.0, got %s", result.Version)
	}
}

func TestRollbackSpecificVersion(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "deployment-history.json")

	manager := &SystemDeployment{
		historyFile: historyFile,
	}

	// Faz deployments
	manager.Deploy(DeployOptions{
		Version:  "v1.0.0",
		Agents:   []string{"agent1"},
		Strategy: StrategyDirect,
	})

	manager.Deploy(DeployOptions{
		Version:  "v2.0.0",
		Agents:   []string{"agent1"},
		Strategy: StrategyDirect,
	})

	// Rollback para versão específica
	rollbackOptions := RollbackOptions{
		Agent:   "agent1",
		Version: "v1.0.0",
		Verify:  false,
	}

	result, err := manager.Rollback(rollbackOptions)
	if err != nil {
		t.Fatalf("Rollback() to specific version failed: %v", err)
	}

	if result.Version != "v1.0.0" {
		t.Errorf("Expected rollback to v1.0.0, got %s", result.Version)
	}
}

func TestRollbackNoHistory(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "deployment-history.json")

	manager := &SystemDeployment{
		historyFile: historyFile,
	}

	rollbackOptions := RollbackOptions{
		Agent: "agent1",
	}

	_, err := manager.Rollback(rollbackOptions)
	if err == nil {
		t.Error("Rollback() should fail when there's no history")
	}
}

func TestGetHistory(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "deployment-history.json")

	manager := &SystemDeployment{
		historyFile: historyFile,
	}

	// Faz alguns deployments
	for i := 1; i <= 5; i++ {
		options := DeployOptions{
			Version:  "v1.0." + string(rune('0'+i)),
			Agents:   []string{"agent1"},
			Strategy: StrategyDirect,
		}
		manager.Deploy(options)
	}

	// Testa obter histórico
	history, err := manager.GetHistory("agent1", 10)
	if err != nil {
		t.Fatalf("GetHistory() failed: %v", err)
	}

	if len(history) != 5 {
		t.Errorf("Expected 5 history records, got %d", len(history))
	}

	// Testa limite
	history, err = manager.GetHistory("agent1", 3)
	if err != nil {
		t.Fatalf("GetHistory() with limit failed: %v", err)
	}

	if len(history) != 3 {
		t.Errorf("Expected 3 history records with limit, got %d", len(history))
	}
}

func TestGetHistoryAllAgents(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "deployment-history.json")

	manager := &SystemDeployment{
		historyFile: historyFile,
	}

	// Deploy em múltiplos agents
	agents := []string{"agent1", "agent2", "agent3"}
	for _, agent := range agents {
		options := DeployOptions{
			Version:  "v1.0.0",
			Agents:   []string{agent},
			Strategy: StrategyDirect,
		}
		manager.Deploy(options)
	}

	// Testa obter histórico de todos os agents
	history, err := manager.GetHistory("", 10)
	if err != nil {
		t.Fatalf("GetHistory() for all agents failed: %v", err)
	}

	if len(history) != 3 {
		t.Errorf("Expected 3 history records for all agents, got %d", len(history))
	}
}

func TestSaveAndLoadHistory(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "deployment-history.json")

	manager := &SystemDeployment{
		historyFile: historyFile,
	}

	record := &DeploymentRecord{
		ID:        "test-123",
		Timestamp: time.Now(),
		Agent:     "agent1",
		Version:   "v1.0.0",
		Strategy:  StrategyDirect,
		Success:   true,
		Duration:  time.Second * 5,
		Message:   "Test deployment",
	}

	err := manager.saveRecord(record)
	if err != nil {
		t.Fatalf("saveRecord() failed: %v", err)
	}

	// Verifica se arquivo foi criado
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		t.Error("History file was not created")
	}

	// Carrega histórico
	records, err := manager.loadHistory()
	if err != nil {
		t.Fatalf("loadHistory() failed: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	if records[0].ID != "test-123" {
		t.Errorf("Expected record ID test-123, got %s", records[0].ID)
	}

	if records[0].Agent != "agent1" {
		t.Errorf("Expected agent agent1, got %s", records[0].Agent)
	}
}

func TestHistoryLimit100Records(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "deployment-history.json")

	manager := &SystemDeployment{
		historyFile: historyFile,
	}

	// Cria 150 records
	for i := 0; i < 150; i++ {
		record := &DeploymentRecord{
			ID:        "test-" + string(rune(i)),
			Timestamp: time.Now(),
			Agent:     "agent1",
			Version:   "v1.0.0",
			Strategy:  StrategyDirect,
			Success:   true,
			Duration:  time.Second,
		}
		manager.saveRecord(record)
	}

	// Verifica que só mantém 100
	records, err := manager.loadHistory()
	if err != nil {
		t.Fatalf("loadHistory() failed: %v", err)
	}

	if len(records) != 100 {
		t.Errorf("Expected history limited to 100 records, got %d", len(records))
	}
}
