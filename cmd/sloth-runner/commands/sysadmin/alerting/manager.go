package alerting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// AlertManager interface para gerenciamento de alertas
type AlertManager interface {
	AddRule(rule *AlertRule) error
	RemoveRule(id string) error
	ListRules() ([]*AlertRule, error)
	CheckRules() ([]*Alert, error)
	GetHistory(limit int) ([]*Alert, error)
}

// AlertType tipo de alerta
type AlertType string

const (
	AlertTypeCPU     AlertType = "cpu"
	AlertTypeMemory  AlertType = "memory"
	AlertTypeDisk    AlertType = "disk"
	AlertTypeService AlertType = "service"
	AlertTypeProcess AlertType = "process"
)

// Severity nível de severidade
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// AlertRule regra de alerta
type AlertRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        AlertType `json:"type"`
	Enabled     bool      `json:"enabled"`
	Threshold   float64   `json:"threshold"`
	Severity    Severity  `json:"severity"`
	Description string    `json:"description"`
	Target      string    `json:"target,omitempty"` // service name, process name, disk path
	CreatedAt   time.Time `json:"created_at"`
}

// Alert alerta disparado
type Alert struct {
	RuleID      string    `json:"rule_id"`
	RuleName    string    `json:"rule_name"`
	Type        AlertType `json:"type"`
	Severity    Severity  `json:"severity"`
	Message     string    `json:"message"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	TriggeredAt time.Time `json:"triggered_at"`
}

// SystemAlertManager implementação padrão
type SystemAlertManager struct {
	rulesFile   string
	historyFile string
}

// NewAlertManager cria um novo alert manager
func NewAlertManager() AlertManager {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".sloth-runner", "alerting")
	os.MkdirAll(dataDir, 0755)

	return &SystemAlertManager{
		rulesFile:   filepath.Join(dataDir, "rules.json"),
		historyFile: filepath.Join(dataDir, "history.json"),
	}
}

// AddRule adiciona uma regra de alerta
func (m *SystemAlertManager) AddRule(rule *AlertRule) error {
	rules, err := m.ListRules()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Gera ID se não tiver
	if rule.ID == "" {
		rule.ID = fmt.Sprintf("rule-%d", time.Now().Unix())
	}

	rule.CreatedAt = time.Now()
	rules = append(rules, rule)

	return m.saveRules(rules)
}

// RemoveRule remove uma regra de alerta
func (m *SystemAlertManager) RemoveRule(id string) error {
	rules, err := m.ListRules()
	if err != nil {
		return err
	}

	var newRules []*AlertRule
	found := false
	for _, rule := range rules {
		if rule.ID != id {
			newRules = append(newRules, rule)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("rule not found: %s", id)
	}

	return m.saveRules(newRules)
}

// ListRules lista todas as regras
func (m *SystemAlertManager) ListRules() ([]*AlertRule, error) {
	data, err := os.ReadFile(m.rulesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []*AlertRule{}, nil
		}
		return nil, err
	}

	var rules []*AlertRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// CheckRules verifica todas as regras e dispara alertas
func (m *SystemAlertManager) CheckRules() ([]*Alert, error) {
	rules, err := m.ListRules()
	if err != nil {
		return nil, err
	}

	var alerts []*Alert

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		alert, err := m.checkRule(rule)
		if err != nil {
			continue // Skip rules that fail to check
		}

		if alert != nil {
			alerts = append(alerts, alert)
		}
	}

	// Save alerts to history
	if len(alerts) > 0 {
		m.saveAlerts(alerts)
	}

	return alerts, nil
}

// GetHistory obtém histórico de alertas
func (m *SystemAlertManager) GetHistory(limit int) ([]*Alert, error) {
	data, err := os.ReadFile(m.historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Alert{}, nil
		}
		return nil, err
	}

	var history []*Alert
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	// Limita resultados
	if limit > 0 && len(history) > limit {
		history = history[:limit]
	}

	return history, nil
}

// checkRule verifica uma regra específica
func (m *SystemAlertManager) checkRule(rule *AlertRule) (*Alert, error) {
	var value float64
	var message string

	switch rule.Type {
	case AlertTypeCPU:
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			return nil, err
		}
		value = cpuPercent[0]
		message = fmt.Sprintf("CPU usage is %.1f%% (threshold: %.1f%%)", value, rule.Threshold)

	case AlertTypeMemory:
		vmem, err := mem.VirtualMemory()
		if err != nil {
			return nil, err
		}
		value = vmem.UsedPercent
		message = fmt.Sprintf("Memory usage is %.1f%% (threshold: %.1f%%)", value, rule.Threshold)

	case AlertTypeDisk:
		target := rule.Target
		if target == "" {
			target = "/"
		}
		usage, err := disk.Usage(target)
		if err != nil {
			return nil, err
		}
		value = usage.UsedPercent
		message = fmt.Sprintf("Disk %s usage is %.1f%% (threshold: %.1f%%)", target, value, rule.Threshold)

	case AlertTypeProcess:
		if rule.Target == "" {
			return nil, fmt.Errorf("process name required for process alerts")
		}
		running, err := m.isProcessRunning(rule.Target)
		if err != nil {
			return nil, err
		}
		if !running {
			value = 0
			message = fmt.Sprintf("Process %s is not running", rule.Target)
		} else {
			value = 1
			message = fmt.Sprintf("Process %s is running", rule.Target)
			// Para processos, threshold 0 significa que deve estar rodando
			// Se está rodando e threshold é 0, não dispara alerta
			if rule.Threshold == 0 {
				return nil, nil
			}
		}

	case AlertTypeService:
		if rule.Target == "" {
			return nil, fmt.Errorf("service name required for service alerts")
		}
		active, err := m.isServiceActive(rule.Target)
		if err != nil {
			return nil, err
		}
		if !active {
			value = 0
			message = fmt.Sprintf("Service %s is not active", rule.Target)
		} else {
			value = 1
			message = fmt.Sprintf("Service %s is active", rule.Target)
			// Para serviços, threshold 0 significa que deve estar ativo
			if rule.Threshold == 0 {
				return nil, nil
			}
		}

	default:
		return nil, fmt.Errorf("unknown alert type: %s", rule.Type)
	}

	// Check if threshold is exceeded
	if value >= rule.Threshold {
		return &Alert{
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Type:        rule.Type,
			Severity:    rule.Severity,
			Message:     message,
			Value:       value,
			Threshold:   rule.Threshold,
			TriggeredAt: time.Now(),
		}, nil
	}

	return nil, nil
}

// isProcessRunning verifica se um processo está rodando
func (m *SystemAlertManager) isProcessRunning(name string) (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range processes {
		procName, err := p.Name()
		if err != nil {
			continue
		}
		if procName == name {
			return true, nil
		}
	}

	return false, nil
}

// isServiceActive verifica se um serviço está ativo (mock para demonstração)
func (m *SystemAlertManager) isServiceActive(name string) (bool, error) {
	// Em produção, isso verificaria o systemctl/service status
	// Para demonstração, simula verificação
	return true, nil
}

// saveRules salva regras em arquivo
func (m *SystemAlertManager) saveRules(rules []*AlertRule) error {
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.rulesFile, data, 0644)
}

// saveAlerts salva alertas no histórico
func (m *SystemAlertManager) saveAlerts(alerts []*Alert) error {
	// Carrega histórico existente
	history, _ := m.GetHistory(0)

	// Adiciona novos alertas
	history = append(alerts, history...)

	// Limita a 1000 alertas
	if len(history) > 1000 {
		history = history[:1000]
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.historyFile, data, 0644)
}
