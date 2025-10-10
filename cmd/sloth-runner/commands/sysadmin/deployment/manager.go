package deployment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DeploymentManager interface para gerenciamento de deployments
type DeploymentManager interface {
	Deploy(options DeployOptions) (*DeploymentResult, error)
	Rollback(options RollbackOptions) (*DeploymentResult, error)
	GetHistory(agent string, limit int) ([]*DeploymentRecord, error)
}

// DeployOptions opções para deployment
type DeployOptions struct {
	Version      string
	Agents       []string
	Strategy     DeploymentStrategy
	DryRun       bool
	HealthCheck  bool
	BackupBefore bool
}

// RollbackOptions opções para rollback
type RollbackOptions struct {
	Agent   string
	Version string // Se vazio, volta para versão anterior
	Verify  bool
}

// DeploymentStrategy estratégia de deployment
type DeploymentStrategy string

const (
	StrategyDirect     DeploymentStrategy = "direct"
	StrategyRolling    DeploymentStrategy = "rolling"
	StrategyCanary     DeploymentStrategy = "canary"
	StrategyBlueGreen  DeploymentStrategy = "blue-green"
)

// DeploymentResult resultado do deployment
type DeploymentResult struct {
	Success       bool
	AgentsUpdated []string
	AgentsFailed  []string
	Version       string
	Duration      time.Duration
	Message       string
}

// DeploymentRecord registro de deployment
type DeploymentRecord struct {
	ID        string
	Timestamp time.Time
	Agent     string
	Version   string
	Strategy  DeploymentStrategy
	Success   bool
	Duration  time.Duration
	Message   string
}

// SystemDeployment implementação padrão
type SystemDeployment struct {
	historyFile string
}

// NewDeploymentManager cria um novo deployment manager
func NewDeploymentManager() DeploymentManager {
	home, _ := os.UserHomeDir()
	historyFile := filepath.Join(home, ".sloth-runner", "deployment-history.json")

	// Cria diretório se não existir
	os.MkdirAll(filepath.Dir(historyFile), 0755)

	return &SystemDeployment{
		historyFile: historyFile,
	}
}

// Deploy realiza deployment
func (d *SystemDeployment) Deploy(options DeployOptions) (*DeploymentResult, error) {
	start := time.Now()

	result := &DeploymentResult{
		Version:       options.Version,
		AgentsUpdated: []string{},
		AgentsFailed:  []string{},
	}

	if options.DryRun {
		result.Success = true
		result.Message = fmt.Sprintf("Dry run: Would deploy version %s to %d agent(s)", options.Version, len(options.Agents))
		result.Duration = time.Since(start)
		return result, nil
	}

	// Simula deployment para cada agente
	for _, agent := range options.Agents {
		// Em uma implementação real, aqui seria feito:
		// 1. Backup do binário atual
		// 2. Download do novo binário
		// 3. Substituição do binário
		// 4. Restart do agent
		// 5. Health check

		// Por enquanto, apenas registra o deployment
		record := &DeploymentRecord{
			ID:        fmt.Sprintf("deploy-%d", time.Now().Unix()),
			Timestamp: time.Now(),
			Agent:     agent,
			Version:   options.Version,
			Strategy:  options.Strategy,
			Success:   true,
			Duration:  time.Second * 2, // Simulado
			Message:   fmt.Sprintf("Deployed version %s successfully", options.Version),
		}

		if err := d.saveRecord(record); err == nil {
			result.AgentsUpdated = append(result.AgentsUpdated, agent)
		} else {
			result.AgentsFailed = append(result.AgentsFailed, agent)
		}
	}

	result.Success = len(result.AgentsFailed) == 0
	result.Duration = time.Since(start)

	if result.Success {
		result.Message = fmt.Sprintf("Successfully deployed version %s to %d agent(s)", options.Version, len(result.AgentsUpdated))
	} else {
		result.Message = fmt.Sprintf("Deployment partially failed: %d succeeded, %d failed", len(result.AgentsUpdated), len(result.AgentsFailed))
	}

	return result, nil
}

// Rollback realiza rollback
func (d *SystemDeployment) Rollback(options RollbackOptions) (*DeploymentResult, error) {
	start := time.Now()

	// Busca histórico para encontrar versão anterior
	history, err := d.GetHistory(options.Agent, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment history: %v", err)
	}

	if len(history) < 2 {
		return nil, fmt.Errorf("no previous deployment found for rollback")
	}

	// Se versão não especificada, usa a penúltima do histórico
	targetVersion := options.Version
	if targetVersion == "" {
		// Busca a versão anterior que foi bem-sucedida
		for i := 1; i < len(history); i++ {
			if history[i].Success {
				targetVersion = history[i].Version
				break
			}
		}
	}

	if targetVersion == "" {
		return nil, fmt.Errorf("no successful previous deployment found")
	}

	// Realiza rollback (simula deployment da versão anterior)
	deployOptions := DeployOptions{
		Version:      targetVersion,
		Agents:       []string{options.Agent},
		Strategy:     StrategyDirect,
		HealthCheck:  options.Verify,
		BackupBefore: true,
	}

	result, err := d.Deploy(deployOptions)
	if err != nil {
		return nil, err
	}

	result.Message = fmt.Sprintf("Rolled back to version %s", targetVersion)
	result.Duration = time.Since(start)

	return result, nil
}

// GetHistory retorna histórico de deployments
func (d *SystemDeployment) GetHistory(agent string, limit int) ([]*DeploymentRecord, error) {
	records, err := d.loadHistory()
	if err != nil {
		return []*DeploymentRecord{}, nil // Retorna vazio se não há histórico
	}

	// Filtra por agent se especificado
	filtered := []*DeploymentRecord{}
	for _, record := range records {
		if agent == "" || record.Agent == agent {
			filtered = append(filtered, record)
		}
	}

	// Limita quantidade
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, nil
}

// saveRecord salva registro no histórico
func (d *SystemDeployment) saveRecord(record *DeploymentRecord) error {
	records, _ := d.loadHistory()
	records = append([]*DeploymentRecord{record}, records...) // Adiciona no início

	// Limita histórico a 100 registros
	if len(records) > 100 {
		records = records[:100]
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(d.historyFile, data, 0644)
}

// loadHistory carrega histórico
func (d *SystemDeployment) loadHistory() ([]*DeploymentRecord, error) {
	data, err := os.ReadFile(d.historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []*DeploymentRecord{}, nil
		}
		return nil, err
	}

	var records []*DeploymentRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}
