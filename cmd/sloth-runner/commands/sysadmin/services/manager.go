package services

import (
	"fmt"
	"os/exec"
	"strings"
)

// ServiceManagerType representa o tipo de gerenciador de serviços
type ServiceManagerType string

const (
	SYSTEMD ServiceManagerType = "systemd"
	INITD   ServiceManagerType = "init.d"
	OPENRC  ServiceManagerType = "openrc"
	NONE    ServiceManagerType = "none"
)

// ServiceStatus representa o status de um serviço
type ServiceStatus string

const (
	StatusActive   ServiceStatus = "active"
	StatusInactive ServiceStatus = "inactive"
	StatusFailed   ServiceStatus = "failed"
	StatusUnknown  ServiceStatus = "unknown"
)

// Service representa um serviço do sistema
type Service struct {
	Name        string
	Status      ServiceStatus
	Enabled     bool
	Description string
	PID         string
	Memory      string
	Uptime      string
}

// ServiceManager interface para operações de gerenciamento de serviços
type ServiceManager interface {
	List() ([]Service, error)
	Status(serviceName string) (*Service, error)
	Start(serviceName string) error
	Stop(serviceName string) error
	Restart(serviceName string) error
	Reload(serviceName string) error
	Enable(serviceName string) error
	Disable(serviceName string) error
	Logs(serviceName string, follow bool, lines int) error
}

// DetectServiceManager detecta qual gerenciador de serviços está disponível
func DetectServiceManager() ServiceManagerType {
	// Verificar systemd
	if commandExists("systemctl") {
		// Verificar se systemd está realmente rodando
		cmd := exec.Command("systemctl", "is-system-running")
		if err := cmd.Run(); err == nil || cmd.ProcessState.ExitCode() == 1 {
			return SYSTEMD
		}
	}

	// Verificar OpenRC
	if commandExists("rc-service") {
		return OPENRC
	}

	// Verificar init.d
	if commandExists("service") || dirExists("/etc/init.d") {
		return INITD
	}

	return NONE
}

// GetServiceManager retorna uma instância do gerenciador apropriado
func GetServiceManager() (ServiceManager, error) {
	smType := DetectServiceManager()

	switch smType {
	case SYSTEMD:
		return &SystemdManager{}, nil
	case INITD:
		return &InitdManager{}, nil
	case OPENRC:
		return &OpenRCManager{}, nil
	default:
		return nil, fmt.Errorf("no supported service manager found")
	}
}

// SystemdManager implementa ServiceManager para systemd
type SystemdManager struct{}

func (s *SystemdManager) List() ([]Service, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []Service
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "UNIT") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		serviceName := strings.TrimSuffix(fields[0], ".service")
		status := parseServiceStatus(fields[2])

		services = append(services, Service{
			Name:        serviceName,
			Status:      status,
			Description: strings.Join(fields[4:], " "),
		})
	}

	return services, nil
}

func (s *SystemdManager) Status(serviceName string) (*Service, error) {
	// Get basic status
	cmd := exec.Command("systemctl", "status", serviceName)
	output, _ := cmd.Output() // Ignorar erro pois pode retornar não-zero para serviços inativos

	service := &Service{
		Name:   serviceName,
		Status: StatusUnknown,
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Active:") {
			if strings.Contains(line, "active (running)") {
				service.Status = StatusActive
			} else if strings.Contains(line, "inactive") {
				service.Status = StatusInactive
			} else if strings.Contains(line, "failed") {
				service.Status = StatusFailed
			}
		} else if strings.Contains(line, "Main PID:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				service.PID = fields[2]
			}
		} else if strings.Contains(line, "Memory:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				service.Memory = fields[1]
			}
		}
	}

	// Check if enabled
	cmd = exec.Command("systemctl", "is-enabled", serviceName)
	output, _ = cmd.Output()
	service.Enabled = strings.TrimSpace(string(output)) == "enabled"

	return service, nil
}

func (s *SystemdManager) Start(serviceName string) error {
	cmd := exec.Command("systemctl", "start", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start %s: %w\nOutput: %s", serviceName, err, string(output))
	}
	return nil
}

func (s *SystemdManager) Stop(serviceName string) error {
	cmd := exec.Command("systemctl", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop %s: %w\nOutput: %s", serviceName, err, string(output))
	}
	return nil
}

func (s *SystemdManager) Restart(serviceName string) error {
	cmd := exec.Command("systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart %s: %w\nOutput: %s", serviceName, err, string(output))
	}
	return nil
}

func (s *SystemdManager) Reload(serviceName string) error {
	cmd := exec.Command("systemctl", "reload", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to reload %s: %w\nOutput: %s", serviceName, err, string(output))
	}
	return nil
}

func (s *SystemdManager) Enable(serviceName string) error {
	cmd := exec.Command("systemctl", "enable", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable %s: %w\nOutput: %s", serviceName, err, string(output))
	}
	return nil
}

func (s *SystemdManager) Disable(serviceName string) error {
	cmd := exec.Command("systemctl", "disable", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable %s: %w\nOutput: %s", serviceName, err, string(output))
	}
	return nil
}

func (s *SystemdManager) Logs(serviceName string, follow bool, lines int) error {
	args := []string{"journalctl", "-u", serviceName, "-n", fmt.Sprintf("%d", lines)}
	if follow {
		args = append(args, "-f")
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = nil // Will be handled by caller
	cmd.Stderr = nil

	return cmd.Run()
}

// Stubs for init.d and OpenRC

type InitdManager struct{}

func (i *InitdManager) List() ([]Service, error) {
	return nil, fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Status(serviceName string) (*Service, error) {
	return nil, fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Start(serviceName string) error {
	return fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Stop(serviceName string) error {
	return fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Restart(serviceName string) error {
	return fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Reload(serviceName string) error {
	return fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Enable(serviceName string) error {
	return fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Disable(serviceName string) error {
	return fmt.Errorf("init.d support not yet implemented")
}

func (i *InitdManager) Logs(serviceName string, follow bool, lines int) error {
	return fmt.Errorf("init.d support not yet implemented")
}

type OpenRCManager struct{}

func (o *OpenRCManager) List() ([]Service, error) {
	return nil, fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Status(serviceName string) (*Service, error) {
	return nil, fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Start(serviceName string) error {
	return fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Stop(serviceName string) error {
	return fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Restart(serviceName string) error {
	return fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Reload(serviceName string) error {
	return fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Enable(serviceName string) error {
	return fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Disable(serviceName string) error {
	return fmt.Errorf("OpenRC support not yet implemented")
}

func (o *OpenRCManager) Logs(serviceName string, follow bool, lines int) error {
	return fmt.Errorf("OpenRC support not yet implemented")
}

// Helper functions

func parseServiceStatus(statusStr string) ServiceStatus {
	statusStr = strings.ToLower(statusStr)
	switch {
	case strings.Contains(statusStr, "active"):
		return StatusActive
	case strings.Contains(statusStr, "inactive"):
		return StatusInactive
	case strings.Contains(statusStr, "failed"):
		return StatusFailed
	default:
		return StatusUnknown
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func dirExists(path string) bool {
	cmd := exec.Command("test", "-d", path)
	return cmd.Run() == nil
}
