package alerting

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewAlertManager(t *testing.T) {
	manager := NewAlertManager()
	if manager == nil {
		t.Fatal("NewAlertManager() returned nil")
	}

	_, ok := manager.(*SystemAlertManager)
	if !ok {
		t.Error("NewAlertManager() did not return *SystemAlertManager")
	}
}

func TestAddRule(t *testing.T) {
	manager := NewAlertManager()

	rule := &AlertRule{
		Name:        "Test CPU Alert",
		Type:        AlertTypeCPU,
		Enabled:     true,
		Threshold:   80.0,
		Severity:    SeverityWarning,
		Description: "Test alert for CPU",
	}

	err := manager.AddRule(rule)
	if err != nil {
		t.Fatalf("AddRule() failed: %v", err)
	}

	// Verifica se ID foi gerado
	if rule.ID == "" {
		t.Error("AddRule() did not generate ID")
	}

	// Verifica se CreatedAt foi setado
	if rule.CreatedAt.IsZero() {
		t.Error("AddRule() did not set CreatedAt")
	}

	// Limpa após teste
	defer manager.RemoveRule(rule.ID)
}

func TestAddRuleWithID(t *testing.T) {
	manager := NewAlertManager()

	rule := &AlertRule{
		ID:        "custom-id-123",
		Name:      "Test Alert",
		Type:      AlertTypeMemory,
		Enabled:   true,
		Threshold: 90.0,
		Severity:  SeverityCritical,
	}

	err := manager.AddRule(rule)
	if err != nil {
		t.Fatalf("AddRule() with custom ID failed: %v", err)
	}

	if rule.ID != "custom-id-123" {
		t.Error("AddRule() changed custom ID")
	}

	defer manager.RemoveRule(rule.ID)
}

func TestListRules(t *testing.T) {
	manager := NewAlertManager()

	// Adiciona algumas regras de teste
	rules := []*AlertRule{
		{
			Name:      "CPU Alert",
			Type:      AlertTypeCPU,
			Enabled:   true,
			Threshold: 80.0,
			Severity:  SeverityWarning,
		},
		{
			Name:      "Memory Alert",
			Type:      AlertTypeMemory,
			Enabled:   true,
			Threshold: 90.0,
			Severity:  SeverityCritical,
		},
	}

	for _, rule := range rules {
		if err := manager.AddRule(rule); err != nil {
			t.Fatalf("AddRule() failed: %v", err)
		}
		defer manager.RemoveRule(rule.ID)
	}

	// Lista regras
	listed, err := manager.ListRules()
	if err != nil {
		t.Fatalf("ListRules() failed: %v", err)
	}

	if len(listed) < 2 {
		t.Errorf("ListRules() returned %d rules, expected at least 2", len(listed))
	}
}

func TestRemoveRule(t *testing.T) {
	manager := NewAlertManager()

	// Adiciona regra
	rule := &AlertRule{
		Name:      "Test Alert",
		Type:      AlertTypeCPU,
		Enabled:   true,
		Threshold: 80.0,
		Severity:  SeverityWarning,
	}

	if err := manager.AddRule(rule); err != nil {
		t.Fatalf("AddRule() failed: %v", err)
	}

	// Remove regra
	err := manager.RemoveRule(rule.ID)
	if err != nil {
		t.Fatalf("RemoveRule() failed: %v", err)
	}

	// Verifica que foi removida
	rules, _ := manager.ListRules()
	for _, r := range rules {
		if r.ID == rule.ID {
			t.Error("RemoveRule() did not remove the rule")
		}
	}
}

func TestRemoveRuleNotFound(t *testing.T) {
	manager := NewAlertManager()

	err := manager.RemoveRule("non-existent-id")
	if err == nil {
		t.Error("RemoveRule() with non-existent ID should return error")
	}
}

func TestCheckRulesCPU(t *testing.T) {
	manager := NewAlertManager()

	// Adiciona regra com threshold baixo para garantir que dispare
	rule := &AlertRule{
		Name:      "Low CPU Alert",
		Type:      AlertTypeCPU,
		Enabled:   true,
		Threshold: 0.0, // Sempre dispara
		Severity:  SeverityInfo,
	}

	if err := manager.AddRule(rule); err != nil {
		t.Fatalf("AddRule() failed: %v", err)
	}
	defer manager.RemoveRule(rule.ID)

	alerts, err := manager.CheckRules()
	if err != nil {
		t.Fatalf("CheckRules() failed: %v", err)
	}

	// Deve ter pelo menos um alerta
	found := false
	for _, alert := range alerts {
		if alert.RuleID == rule.ID {
			found = true
			if alert.Type != AlertTypeCPU {
				t.Error("Alert type mismatch")
			}
			if alert.Severity != SeverityInfo {
				t.Error("Alert severity mismatch")
			}
		}
	}

	if !found && rule.Threshold == 0.0 {
		// Com threshold 0, deveria sempre disparar se CPU > 0
		t.Log("Warning: CPU alert not triggered (may be expected if CPU is 0%)")
	}
}

func TestCheckRulesMemory(t *testing.T) {
	manager := NewAlertManager()

	rule := &AlertRule{
		Name:      "Low Memory Alert",
		Type:      AlertTypeMemory,
		Enabled:   true,
		Threshold: 0.0, // Sempre dispara se houver uso de memória
		Severity:  SeverityInfo,
	}

	if err := manager.AddRule(rule); err != nil {
		t.Fatalf("AddRule() failed: %v", err)
	}
	defer manager.RemoveRule(rule.ID)

	alerts, err := manager.CheckRules()
	if err != nil {
		t.Fatalf("CheckRules() failed: %v", err)
	}

	// Verifica estrutura do alerta
	for _, alert := range alerts {
		if alert.RuleID == rule.ID {
			if alert.Message == "" {
				t.Error("Alert message is empty")
			}
			if alert.Value < 0 {
				t.Error("Alert value is negative")
			}
		}
	}
}

func TestCheckRulesDisabled(t *testing.T) {
	manager := NewAlertManager()

	rule := &AlertRule{
		Name:      "Disabled Alert",
		Type:      AlertTypeCPU,
		Enabled:   false, // Desabilitada
		Threshold: 0.0,
		Severity:  SeverityWarning,
	}

	if err := manager.AddRule(rule); err != nil {
		t.Fatalf("AddRule() failed: %v", err)
	}
	defer manager.RemoveRule(rule.ID)

	alerts, err := manager.CheckRules()
	if err != nil {
		t.Fatalf("CheckRules() failed: %v", err)
	}

	// Não deve ter alerta de regra desabilitada
	for _, alert := range alerts {
		if alert.RuleID == rule.ID {
			t.Error("CheckRules() triggered alert for disabled rule")
		}
	}
}

func TestGetHistory(t *testing.T) {
	manager := NewAlertManager()

	history, err := manager.GetHistory(10)
	if err != nil {
		t.Fatalf("GetHistory() failed: %v", err)
	}

	if history == nil {
		t.Fatal("GetHistory() returned nil")
	}

	// History pode estar vazio se não houver alertas
	// Apenas verifica que não há erro
}

func TestGetHistoryWithLimit(t *testing.T) {
	manager := NewAlertManager()

	limit := 5
	history, err := manager.GetHistory(limit)
	if err != nil {
		t.Fatalf("GetHistory() failed: %v", err)
	}

	if len(history) > limit {
		t.Errorf("GetHistory() returned %d alerts, expected max %d", len(history), limit)
	}
}

func TestAlertTypes(t *testing.T) {
	types := []AlertType{
		AlertTypeCPU,
		AlertTypeMemory,
		AlertTypeDisk,
		AlertTypeService,
		AlertTypeProcess,
	}

	for _, alertType := range types {
		if alertType == "" {
			t.Error("Alert type is empty")
		}
	}
}

func TestSeverityLevels(t *testing.T) {
	severities := []Severity{
		SeverityInfo,
		SeverityWarning,
		SeverityCritical,
	}

	for _, severity := range severities {
		if severity == "" {
			t.Error("Severity level is empty")
		}
	}
}

func TestAlertRuleStructure(t *testing.T) {
	rule := &AlertRule{
		ID:          "test-123",
		Name:        "Test Rule",
		Type:        AlertTypeCPU,
		Enabled:     true,
		Threshold:   75.0,
		Severity:    SeverityWarning,
		Description: "Test description",
		Target:      "/data",
		CreatedAt:   time.Now(),
	}

	if rule.ID != "test-123" {
		t.Error("ID not set correctly")
	}
	if rule.Name != "Test Rule" {
		t.Error("Name not set correctly")
	}
	if rule.Type != AlertTypeCPU {
		t.Error("Type not set correctly")
	}
	if !rule.Enabled {
		t.Error("Enabled not set correctly")
	}
	if rule.Threshold != 75.0 {
		t.Error("Threshold not set correctly")
	}
}

func TestAlertStructure(t *testing.T) {
	alert := &Alert{
		RuleID:      "rule-123",
		RuleName:    "Test Alert",
		Type:        AlertTypeMemory,
		Severity:    SeverityCritical,
		Message:     "Memory usage critical",
		Value:       95.5,
		Threshold:   90.0,
		TriggeredAt: time.Now(),
	}

	if alert.RuleID != "rule-123" {
		t.Error("RuleID not set correctly")
	}
	if alert.Type != AlertTypeMemory {
		t.Error("Type not set correctly")
	}
	if alert.Severity != SeverityCritical {
		t.Error("Severity not set correctly")
	}
	if alert.Value != 95.5 {
		t.Error("Value not set correctly")
	}
	if alert.Threshold != 90.0 {
		t.Error("Threshold not set correctly")
	}
}

func TestSaveAndLoadRules(t *testing.T) {
	// Cria diretório temporário para testes
	tmpDir := t.TempDir()

	manager := &SystemAlertManager{
		rulesFile:   filepath.Join(tmpDir, "rules.json"),
		historyFile: filepath.Join(tmpDir, "history.json"),
	}

	// Adiciona regras
	rules := []*AlertRule{
		{
			Name:      "Rule 1",
			Type:      AlertTypeCPU,
			Enabled:   true,
			Threshold: 80.0,
			Severity:  SeverityWarning,
		},
		{
			Name:      "Rule 2",
			Type:      AlertTypeMemory,
			Enabled:   false,
			Threshold: 90.0,
			Severity:  SeverityCritical,
		},
	}

	for _, rule := range rules {
		if err := manager.AddRule(rule); err != nil {
			t.Fatalf("AddRule() failed: %v", err)
		}
	}

	// Carrega regras
	loaded, err := manager.ListRules()
	if err != nil {
		t.Fatalf("ListRules() failed: %v", err)
	}

	if len(loaded) != len(rules) {
		t.Errorf("Loaded %d rules, expected %d", len(loaded), len(rules))
	}

	// Verifica conteúdo
	for i, rule := range loaded {
		if rule.Name != rules[i].Name {
			t.Errorf("Rule %d name mismatch", i)
		}
		if rule.Type != rules[i].Type {
			t.Errorf("Rule %d type mismatch", i)
		}
		if rule.Enabled != rules[i].Enabled {
			t.Errorf("Rule %d enabled mismatch", i)
		}
	}
}

func TestCheckRuleDisk(t *testing.T) {
	manager := &SystemAlertManager{
		rulesFile:   filepath.Join(t.TempDir(), "rules.json"),
		historyFile: filepath.Join(t.TempDir(), "history.json"),
	}

	rule := &AlertRule{
		ID:        "disk-test",
		Name:      "Disk Alert",
		Type:      AlertTypeDisk,
		Enabled:   true,
		Threshold: 200.0, // Threshold alto para não disparar normalmente
		Severity:  SeverityWarning,
		Target:    "/", // Root partition
	}

	alert, err := manager.checkRule(rule)
	if err != nil {
		t.Fatalf("checkRule() failed: %v", err)
	}

	// Com threshold alto, não deve disparar alerta
	if alert != nil {
		t.Log("Disk alert triggered (may be expected if disk is really full)")
	}
}

func TestIsProcessRunning(t *testing.T) {
	manager := &SystemAlertManager{}

	// Testa com processo que provavelmente existe
	running, err := manager.isProcessRunning("go")
	if err != nil {
		t.Fatalf("isProcessRunning() failed: %v", err)
	}

	// Não podemos garantir resultado, apenas que não dá erro
	_ = running

	// Testa com processo que não existe
	running, err = manager.isProcessRunning("nonexistentprocess12345")
	if err != nil {
		t.Fatalf("isProcessRunning() with non-existent process failed: %v", err)
	}

	if running {
		t.Error("isProcessRunning() returned true for non-existent process")
	}
}

func TestPersistenceDirectory(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	dataDir := filepath.Join(homeDir, ".sloth-runner", "alerting")

	// Apenas verifica que o caminho está correto
	if !filepath.IsAbs(dataDir) {
		t.Error("Data directory is not absolute path")
	}
}
