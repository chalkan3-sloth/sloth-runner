package security

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewScanner(t *testing.T) {
	scanner := NewScanner()
	if scanner == nil {
		t.Fatal("NewScanner() returned nil")
	}

	_, ok := scanner.(*SystemScanner)
	if !ok {
		t.Error("NewScanner() did not return *SystemScanner")
	}
}

func TestAudit(t *testing.T) {
	scanner := NewScanner()

	options := AuditOptions{
		Since:            24 * time.Hour,
		ShowFailedAuth:   true,
		DetectAnomalies:  true,
		OutputFormat:     "table",
	}

	report, err := scanner.Audit(options)
	if err != nil {
		t.Fatalf("Audit() failed: %v", err)
	}

	if report == nil {
		t.Fatal("Audit() returned nil report")
	}

	if report.Timestamp.IsZero() {
		t.Error("Report timestamp is zero")
	}

	if report.TotalEvents == 0 {
		t.Error("Total events is zero")
	}

	if len(report.Recommendations) == 0 {
		t.Error("Recommendations list is empty")
	}
}

func TestAuditWithFailedAuth(t *testing.T) {
	scanner := NewScanner()

	options := AuditOptions{
		Since:          24 * time.Hour,
		ShowFailedAuth: true,
	}

	report, err := scanner.Audit(options)
	if err != nil {
		t.Fatalf("Audit() with failed auth failed: %v", err)
	}

	if report.FailedAuthAttempts == 0 {
		t.Error("Expected failed auth attempts to be detected")
	}

	foundFailedAuth := false
	for _, event := range report.SuspiciousEvents {
		if event.Type == "Failed authentication" {
			foundFailedAuth = true
			break
		}
	}

	if !foundFailedAuth {
		t.Error("Failed authentication event not found in suspicious events")
	}
}

func TestAuditWithAnomalies(t *testing.T) {
	scanner := NewScanner()

	options := AuditOptions{
		Since:           24 * time.Hour,
		DetectAnomalies: true,
	}

	report, err := scanner.Audit(options)
	if err != nil {
		t.Fatalf("Audit() with anomalies failed: %v", err)
	}

	foundAnomaly := false
	for _, event := range report.SuspiciousEvents {
		if event.Type == "Multiple failed SSH attempts" {
			foundAnomaly = true
			break
		}
	}

	if !foundAnomaly {
		t.Error("Anomaly detection event not found")
	}
}

func TestScan(t *testing.T) {
	scanner := NewScanner()

	options := ScanOptions{
		Agent: "test-agent",
		Full:  false,
	}

	report, err := scanner.Scan(options)
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if report == nil {
		t.Fatal("Scan() returned nil report")
	}

	if report.Timestamp.IsZero() {
		t.Error("Report timestamp is zero")
	}

	if report.Agent != "test-agent" {
		t.Errorf("Expected agent test-agent, got %s", report.Agent)
	}

	if report.Score < 0 || report.Score > 100 {
		t.Errorf("Score out of range: %d", report.Score)
	}
}

func TestScanCVEOnly(t *testing.T) {
	scanner := NewScanner()

	options := ScanOptions{
		Agent:   "test-agent",
		CVEOnly: true,
	}

	report, err := scanner.Scan(options)
	if err != nil {
		t.Fatalf("Scan() with CVE-only failed: %v", err)
	}

	if len(report.Vulnerabilities) == 0 {
		t.Error("Expected CVE vulnerabilities to be found")
	}

	// Verifica que todas as vulnerabilidades têm CVE
	for _, vuln := range report.Vulnerabilities {
		if vuln.CVE == "" {
			t.Error("Vulnerability without CVE found in CVE-only scan")
		}
	}
}

func TestScanDependencyAudit(t *testing.T) {
	scanner := NewScanner()

	options := ScanOptions{
		Agent:           "test-agent",
		DependencyAudit: true,
	}

	report, err := scanner.Scan(options)
	if err != nil {
		t.Fatalf("Scan() with dependency audit failed: %v", err)
	}

	if len(report.Vulnerabilities) == 0 {
		t.Error("Expected dependency vulnerabilities to be found")
	}
}

func TestScanFull(t *testing.T) {
	scanner := NewScanner()

	options := ScanOptions{
		Agent: "test-agent",
		Full:  true,
	}

	report, err := scanner.Scan(options)
	if err != nil {
		t.Fatalf("Full scan failed: %v", err)
	}

	// Full scan deve incluir CVEs e dependências
	if len(report.Vulnerabilities) == 0 {
		t.Error("Expected vulnerabilities in full scan")
	}

	// Full scan pode incluir config issues e permission issues
	// mas não é garantido em todos os ambientes
}

func TestScanCVE(t *testing.T) {
	scanner := &SystemScanner{}

	vulns := scanner.scanCVE()

	if len(vulns) == 0 {
		t.Error("scanCVE() returned no vulnerabilities")
	}

	for _, vuln := range vulns {
		if vuln.CVE == "" {
			t.Error("Vulnerability missing CVE")
		}
		if vuln.Package == "" {
			t.Error("Vulnerability missing package name")
		}
		if vuln.Severity == "" {
			t.Error("Vulnerability missing severity")
		}
	}
}

func TestScanDependencies(t *testing.T) {
	scanner := &SystemScanner{}

	vulns := scanner.scanDependencies()

	if len(vulns) == 0 {
		t.Error("scanDependencies() returned no vulnerabilities")
	}

	for _, vuln := range vulns {
		if vuln.Package == "" {
			t.Error("Dependency vulnerability missing package name")
		}
	}
}

func TestScanConfiguration(t *testing.T) {
	scanner := &SystemScanner{}

	// Cria arquivo de config temporário com permissões inseguras
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configFile, []byte("test: config"), 0644)

	// Altera permissões para inseguras
	os.Chmod(configFile, 0666)

	// Nota: scanConfiguration verifica o arquivo real em ~/.sloth-runner/config.yaml
	// então esse teste pode não detectar o arquivo temporário
	issues := scanner.scanConfiguration()

	// O teste não deve falhar se não houver issues, pois depende do ambiente
	if len(issues) > 0 {
		for _, issue := range issues {
			if issue.File == "" {
				t.Error("Config issue missing file name")
			}
			if issue.Issue == "" {
				t.Error("Config issue missing issue description")
			}
		}
	}
}

func TestScanPermissions(t *testing.T) {
	scanner := &SystemScanner{}

	// Este teste depende do ambiente ter um diretório .ssh
	issues := scanner.scanPermissions()

	// O teste não deve falhar se não houver issues
	if len(issues) > 0 {
		for _, issue := range issues {
			if issue.Path == "" {
				t.Error("Permission issue missing path")
			}
			if issue.CurrentPerm == "" {
				t.Error("Permission issue missing current permission")
			}
			if issue.ExpectedPerm == "" {
				t.Error("Permission issue missing expected permission")
			}
		}
	}
}

func TestCalculateScore(t *testing.T) {
	scanner := &SystemScanner{}

	tests := []struct {
		name     string
		report   *ScanReport
		minScore int
		maxScore int
		severity SeverityLevel
	}{
		{
			name: "No issues - perfect score",
			report: &ScanReport{
				Vulnerabilities:  []Vulnerability{},
				ConfigIssues:     []ConfigIssue{},
				PermissionIssues: []PermissionIssue{},
			},
			minScore: 100,
			maxScore: 100,
			severity: SeverityInfo,
		},
		{
			name: "Critical vulnerability",
			report: &ScanReport{
				Vulnerabilities: []Vulnerability{
					{Severity: SeverityCritical},
				},
			},
			minScore: 60,
			maxScore: 70,
			severity: SeverityCritical,
		},
		{
			name: "High severity issues",
			report: &ScanReport{
				Vulnerabilities: []Vulnerability{
					{Severity: SeverityHigh},
				},
				ConfigIssues: []ConfigIssue{
					{Severity: SeverityHigh},
				},
			},
			minScore: 45,
			maxScore: 65,
			severity: SeverityHigh,
		},
		{
			name: "Multiple medium issues",
			report: &ScanReport{
				Vulnerabilities: []Vulnerability{
					{Severity: SeverityMedium},
					{Severity: SeverityMedium},
				},
			},
			minScore: 75,
			maxScore: 85,
			severity: SeverityMedium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, severity := scanner.calculateScore(tt.report)

			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Score %d not in expected range [%d, %d]", score, tt.minScore, tt.maxScore)
			}

			if severity != tt.severity {
				t.Errorf("Expected severity %s, got %s", tt.severity, severity)
			}

			// Verifica que score está no range válido
			if score < 0 || score > 100 {
				t.Errorf("Score out of valid range: %d", score)
			}
		})
	}
}

func TestSeverityLevels(t *testing.T) {
	levels := []SeverityLevel{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
		SeverityInfo,
	}

	for _, level := range levels {
		if level == "" {
			t.Errorf("Severity level is empty")
		}
	}
}

func TestVulnerabilityStructure(t *testing.T) {
	vuln := Vulnerability{
		CVE:         "CVE-2024-1234",
		Package:     "test-package",
		Version:     "1.0.0",
		Description: "Test vulnerability",
		Severity:    SeverityHigh,
		FixVersion:  "1.0.1",
	}

	if vuln.CVE == "" {
		t.Error("CVE is empty")
	}
	if vuln.Package == "" {
		t.Error("Package is empty")
	}
	if vuln.Severity == "" {
		t.Error("Severity is empty")
	}
}

func TestSecurityEventStructure(t *testing.T) {
	event := SecurityEvent{
		Timestamp:   time.Now(),
		Type:        "test-event",
		Source:      "test-source",
		Description: "Test description",
		Severity:    SeverityMedium,
	}

	if event.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}
	if event.Type == "" {
		t.Error("Type is empty")
	}
	if event.Severity == "" {
		t.Error("Severity is empty")
	}
}
