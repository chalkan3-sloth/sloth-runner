package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SecurityScanner interface para security scanning
type SecurityScanner interface {
	Audit(options AuditOptions) (*AuditReport, error)
	Scan(options ScanOptions) (*ScanReport, error)
}

// AuditOptions opções para audit
type AuditOptions struct {
	Since            time.Duration
	ShowFailedAuth   bool
	DetectAnomalies  bool
	OutputFormat     string
}

// ScanOptions opções para scan
type ScanOptions struct {
	Agent          string
	Full           bool
	CVEOnly        bool
	DependencyAudit bool
}

// AuditReport relatório de auditoria
type AuditReport struct {
	Timestamp         time.Time
	TotalEvents       int
	FailedAuthAttempts int
	SuspiciousEvents  []SecurityEvent
	Recommendations   []string
}

// ScanReport relatório de scan
type ScanReport struct {
	Timestamp       time.Time
	Agent           string
	Vulnerabilities []Vulnerability
	ConfigIssues    []ConfigIssue
	PermissionIssues []PermissionIssue
	Severity        SeverityLevel
	Score           int // 0-100, onde 100 = muito seguro
}

// SecurityEvent evento de segurança
type SecurityEvent struct {
	Timestamp   time.Time
	Type        string
	Source      string
	Description string
	Severity    SeverityLevel
}

// Vulnerability vulnerabilidade encontrada
type Vulnerability struct {
	CVE         string
	Package     string
	Version     string
	Description string
	Severity    SeverityLevel
	FixVersion  string
}

// ConfigIssue problema de configuração
type ConfigIssue struct {
	File        string
	Issue       string
	Description string
	Severity    SeverityLevel
}

// PermissionIssue problema de permissões
type PermissionIssue struct {
	Path         string
	CurrentPerm  string
	ExpectedPerm string
	Description  string
	Severity     SeverityLevel
}

// SeverityLevel nível de severidade
type SeverityLevel string

const (
	SeverityCritical SeverityLevel = "Critical"
	SeverityHigh     SeverityLevel = "High"
	SeverityMedium   SeverityLevel = "Medium"
	SeverityLow      SeverityLevel = "Low"
	SeverityInfo     SeverityLevel = "Info"
)

// SystemScanner implementação padrão
type SystemScanner struct{}

// NewScanner cria um novo scanner
func NewScanner() SecurityScanner {
	return &SystemScanner{}
}

// Audit realiza auditoria de segurança
func (s *SystemScanner) Audit(options AuditOptions) (*AuditReport, error) {
	report := &AuditReport{
		Timestamp:         time.Now(),
		FailedAuthAttempts: 0,
		SuspiciousEvents:  []SecurityEvent{},
		Recommendations:   []string{},
	}

	// Simula análise de logs
	// Em uma implementação real, leria logs de /var/log/auth.log, etc
	report.TotalEvents = 1247

	// Simula eventos suspeitos
	if options.DetectAnomalies {
		report.SuspiciousEvents = append(report.SuspiciousEvents, SecurityEvent{
			Timestamp:   time.Now().Add(-2 * time.Hour),
			Type:        "Multiple failed SSH attempts",
			Source:      "192.168.1.100",
			Description: "5 failed SSH login attempts in 2 minutes",
			Severity:    SeverityMedium,
		})
	}

	// Simula falhas de autenticação
	if options.ShowFailedAuth {
		report.FailedAuthAttempts = 12
		report.SuspiciousEvents = append(report.SuspiciousEvents, SecurityEvent{
			Timestamp:   time.Now().Add(-30 * time.Minute),
			Type:        "Failed authentication",
			Source:      "unknown-user@192.168.1.50",
			Description: "Multiple failed login attempts from unknown user",
			Severity:    SeverityHigh,
		})
	}

	// Recomendações
	if report.FailedAuthAttempts > 10 {
		report.Recommendations = append(report.Recommendations, "Consider implementing fail2ban or rate limiting")
	}
	if len(report.SuspiciousEvents) > 0 {
		report.Recommendations = append(report.Recommendations, "Review and investigate suspicious events")
		report.Recommendations = append(report.Recommendations, "Enable two-factor authentication for sensitive accounts")
	}
	if len(report.Recommendations) == 0 {
		report.Recommendations = append(report.Recommendations, "No security issues detected - system appears secure")
	}

	return report, nil
}

// Scan realiza scan de vulnerabilidades
func (s *SystemScanner) Scan(options ScanOptions) (*ScanReport, error) {
	report := &ScanReport{
		Timestamp:       time.Now(),
		Agent:           options.Agent,
		Vulnerabilities: []Vulnerability{},
		ConfigIssues:    []ConfigIssue{},
		PermissionIssues: []PermissionIssue{},
		Score:           100,
		Severity:        SeverityInfo,
	}

	// Scan de vulnerabilidades CVE
	if options.CVEOnly || options.Full {
		vulns := s.scanCVE()
		report.Vulnerabilities = append(report.Vulnerabilities, vulns...)
	}

	// Scan de dependências
	if options.DependencyAudit || options.Full {
		depVulns := s.scanDependencies()
		report.Vulnerabilities = append(report.Vulnerabilities, depVulns...)
	}

	// Scan de configurações
	if options.Full {
		configIssues := s.scanConfiguration()
		report.ConfigIssues = append(report.ConfigIssues, configIssues...)
	}

	// Scan de permissões
	if options.Full {
		permIssues := s.scanPermissions()
		report.PermissionIssues = append(report.PermissionIssues, permIssues...)
	}

	// Calcula score e severity
	report.Score, report.Severity = s.calculateScore(report)

	return report, nil
}

// scanCVE simula scan de CVE
func (s *SystemScanner) scanCVE() []Vulnerability {
	// Em uma implementação real, consultaria banco de CVEs
	return []Vulnerability{
		{
			CVE:         "CVE-2024-1234",
			Package:     "openssl",
			Version:     "1.1.1",
			Description: "Remote code execution vulnerability in SSL/TLS implementation",
			Severity:    SeverityHigh,
			FixVersion:  "1.1.1w",
		},
	}
}

// scanDependencies simula scan de dependências
func (s *SystemScanner) scanDependencies() []Vulnerability {
	// Em uma implementação real, usaria tools como `go mod verify`, `npm audit`, etc
	return []Vulnerability{
		{
			CVE:         "CVE-2024-5678",
			Package:     "golang.org/x/crypto",
			Version:     "0.10.0",
			Description: "Cryptographic vulnerability in SSH implementation",
			Severity:    SeverityMedium,
			FixVersion:  "0.14.0",
		},
	}
}

// scanConfiguration simula scan de configuração
func (s *SystemScanner) scanConfiguration() []ConfigIssue {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".sloth-runner", "config.yaml")

	issues := []ConfigIssue{}

	// Verifica se arquivo de config existe e tem permissões corretas
	if info, err := os.Stat(configPath); err == nil {
		perm := info.Mode().Perm()
		if perm&0077 != 0 {
			issues = append(issues, ConfigIssue{
				File:        configPath,
				Issue:       "Insecure file permissions",
				Description: "Configuration file is readable by other users",
				Severity:    SeverityMedium,
			})
		}
	}

	return issues
}

// scanPermissions simula scan de permissões
func (s *SystemScanner) scanPermissions() []PermissionIssue {
	home, _ := os.UserHomeDir()
	sshDir := filepath.Join(home, ".ssh")

	issues := []PermissionIssue{}

	// Verifica permissões do diretório .ssh
	if info, err := os.Stat(sshDir); err == nil {
		perm := info.Mode().Perm()
		expectedPerm := os.FileMode(0700)
		if perm != expectedPerm {
			issues = append(issues, PermissionIssue{
				Path:         sshDir,
				CurrentPerm:  fmt.Sprintf("%o", perm),
				ExpectedPerm: fmt.Sprintf("%o", expectedPerm),
				Description:  "SSH directory should have 700 permissions",
				Severity:     SeverityHigh,
			})
		}

		// Verifica chaves privadas
		filepath.Walk(sshDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			// Verifica se é chave privada
			if strings.HasPrefix(info.Name(), "id_") && !strings.HasSuffix(info.Name(), ".pub") {
				perm := info.Mode().Perm()
				expectedPerm := os.FileMode(0600)
				if perm != expectedPerm {
					issues = append(issues, PermissionIssue{
						Path:         path,
						CurrentPerm:  fmt.Sprintf("%o", perm),
						ExpectedPerm: fmt.Sprintf("%o", expectedPerm),
						Description:  "Private key should have 600 permissions",
						Severity:     SeverityCritical,
					})
				}
			}

			return nil
		})
	}

	return issues
}

// calculateScore calcula score de segurança
func (s *SystemScanner) calculateScore(report *ScanReport) (int, SeverityLevel) {
	score := 100
	maxSeverity := SeverityInfo

	// Reduz score baseado em vulnerabilidades
	for _, vuln := range report.Vulnerabilities {
		switch vuln.Severity {
		case SeverityCritical:
			score -= 30
			maxSeverity = SeverityCritical
		case SeverityHigh:
			score -= 20
			if maxSeverity != SeverityCritical {
				maxSeverity = SeverityHigh
			}
		case SeverityMedium:
			score -= 10
			if maxSeverity == SeverityInfo || maxSeverity == SeverityLow {
				maxSeverity = SeverityMedium
			}
		case SeverityLow:
			score -= 5
			if maxSeverity == SeverityInfo {
				maxSeverity = SeverityLow
			}
		}
	}

	// Reduz score baseado em config issues
	for _, issue := range report.ConfigIssues {
		switch issue.Severity {
		case SeverityCritical:
			score -= 25
		case SeverityHigh:
			score -= 15
		case SeverityMedium:
			score -= 8
		case SeverityLow:
			score -= 3
		}
	}

	// Reduz score baseado em permission issues
	for _, issue := range report.PermissionIssues {
		switch issue.Severity {
		case SeverityCritical:
			score -= 25
		case SeverityHigh:
			score -= 15
		case SeverityMedium:
			score -= 8
		case SeverityLow:
			score -= 3
		}
	}

	if score < 0 {
		score = 0
	}

	return score, maxSeverity
}
