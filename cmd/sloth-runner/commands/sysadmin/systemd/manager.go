package systemd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// SystemdManager interface para gerenciamento de serviços systemd
type SystemdManager interface {
	List(options ListOptions) ([]*ServiceInfo, error)
	Status(service string) (*ServiceDetail, error)
	Start(service string) error
	Stop(service string) error
	Restart(service string) error
	Enable(service string) error
	Disable(service string) error
	Logs(service string, lines int, follow bool) (string, error)
}

// ListOptions opções para listar serviços
type ListOptions struct {
	Status string // running, stopped, failed, all
	Filter string // filtro por nome
	Type   string // service, socket, timer
}

// ServiceInfo informações básicas do serviço
type ServiceInfo struct {
	Name        string
	LoadState   string
	ActiveState string
	SubState    string
	Description string
}

// ServiceDetail informações detalhadas do serviço
type ServiceDetail struct {
	*ServiceInfo
	MainPID        int
	Memory         uint64
	CPUUsage       float64
	TasksCurrent   int
	TasksMax       int
	RestartCount   int
	ActiveSince    time.Time
	Fragment       string
	DropInPaths    []string
	Documentation  []string
	ExecStart      string
	ExecReload     string
	ExecStop       string
	User           string
	Group          string
	Restart        string
	TimeoutStartS  int
	TimeoutStopS   int
}

// SystemSystemdManager implementação padrão
type SystemSystemdManager struct{}

// NewSystemdManager cria um novo systemd manager
func NewSystemdManager() SystemdManager {
	return &SystemSystemdManager{}
}

// List lista serviços
func (m *SystemSystemdManager) List(options ListOptions) ([]*ServiceInfo, error) {
	args := []string{"list-units", "--type=service", "--all", "--no-pager", "--no-legend"}

	if options.Type != "" && options.Type != "service" {
		args[1] = fmt.Sprintf("--type=%s", options.Type)
	}

	cmd := exec.Command("systemctl", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	var services []*ServiceInfo
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		service := &ServiceInfo{
			Name:        fields[0],
			LoadState:   fields[1],
			ActiveState: fields[2],
			SubState:    fields[3],
		}

		// Description é tudo depois do 4º campo
		if len(fields) > 4 {
			service.Description = strings.Join(fields[4:], " ")
		}

		// Aplica filtros
		if options.Status != "" && options.Status != "all" {
			if options.Status == "running" && service.ActiveState != "active" {
				continue
			}
			if options.Status == "stopped" && service.ActiveState != "inactive" {
				continue
			}
			if options.Status == "failed" && service.ActiveState != "failed" {
				continue
			}
		}

		if options.Filter != "" {
			if !strings.Contains(strings.ToLower(service.Name), strings.ToLower(options.Filter)) &&
			   !strings.Contains(strings.ToLower(service.Description), strings.ToLower(options.Filter)) {
				continue
			}
		}

		services = append(services, service)
	}

	return services, nil
}

// Status obtém status detalhado de um serviço
func (m *SystemSystemdManager) Status(service string) (*ServiceDetail, error) {
	// Garante que termina com .service se não tiver extensão
	if !strings.Contains(service, ".") {
		service = service + ".service"
	}

	cmd := exec.Command("systemctl", "show", service, "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get service status: %v", err)
	}

	detail := &ServiceDetail{
		ServiceInfo: &ServiceInfo{
			Name: service,
		},
	}

	// Parse output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "LoadState":
			detail.LoadState = value
		case "ActiveState":
			detail.ActiveState = value
		case "SubState":
			detail.SubState = value
		case "Description":
			detail.Description = value
		case "MainPID":
			fmt.Sscanf(value, "%d", &detail.MainPID)
		case "MemoryCurrent":
			fmt.Sscanf(value, "%d", &detail.Memory)
		case "CPUUsageNSec":
			var nsec uint64
			fmt.Sscanf(value, "%d", &nsec)
			detail.CPUUsage = float64(nsec) / 1e9
		case "TasksCurrent":
			fmt.Sscanf(value, "%d", &detail.TasksCurrent)
		case "TasksMax":
			fmt.Sscanf(value, "%d", &detail.TasksMax)
		case "NRestarts":
			fmt.Sscanf(value, "%d", &detail.RestartCount)
		case "ActiveEnterTimestamp":
			if value != "" && value != "n/a" {
				// Parse timestamp format
				layout := "Mon 2006-01-02 15:04:05 MST"
				t, err := time.Parse(layout, value)
				if err == nil {
					detail.ActiveSince = t
				}
			}
		case "FragmentPath":
			detail.Fragment = value
		case "ExecStart":
			detail.ExecStart = value
		case "ExecReload":
			detail.ExecReload = value
		case "ExecStop":
			detail.ExecStop = value
		case "User":
			detail.User = value
		case "Group":
			detail.Group = value
		case "Restart":
			detail.Restart = value
		case "TimeoutStartUSec":
			var usec uint64
			fmt.Sscanf(value, "%d", &usec)
			detail.TimeoutStartS = int(usec / 1e6)
		case "TimeoutStopUSec":
			var usec uint64
			fmt.Sscanf(value, "%d", &usec)
			detail.TimeoutStopS = int(usec / 1e6)
		}
	}

	return detail, nil
}

// Start inicia um serviço
func (m *SystemSystemdManager) Start(service string) error {
	if !strings.Contains(service, ".") {
		service = service + ".service"
	}

	cmd := exec.Command("systemctl", "start", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start service: %v - %s", err, string(output))
	}

	return nil
}

// Stop para um serviço
func (m *SystemSystemdManager) Stop(service string) error {
	if !strings.Contains(service, ".") {
		service = service + ".service"
	}

	cmd := exec.Command("systemctl", "stop", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop service: %v - %s", err, string(output))
	}

	return nil
}

// Restart reinicia um serviço
func (m *SystemSystemdManager) Restart(service string) error {
	if !strings.Contains(service, ".") {
		service = service + ".service"
	}

	cmd := exec.Command("systemctl", "restart", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restart service: %v - %s", err, string(output))
	}

	return nil
}

// Enable habilita um serviço no boot
func (m *SystemSystemdManager) Enable(service string) error {
	if !strings.Contains(service, ".") {
		service = service + ".service"
	}

	cmd := exec.Command("systemctl", "enable", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable service: %v - %s", err, string(output))
	}

	return nil
}

// Disable desabilita um serviço no boot
func (m *SystemSystemdManager) Disable(service string) error {
	if !strings.Contains(service, ".") {
		service = service + ".service"
	}

	cmd := exec.Command("systemctl", "disable", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to disable service: %v - %s", err, string(output))
	}

	return nil
}

// Logs obtém logs de um serviço
func (m *SystemSystemdManager) Logs(service string, lines int, follow bool) (string, error) {
	if !strings.Contains(service, ".") {
		service = service + ".service"
	}

	args := []string{"-u", service, "--no-pager"}

	if lines > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", lines))
	}

	if follow {
		args = append(args, "-f")
	}

	cmd := exec.Command("journalctl", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get logs: %v - %s", err, stderr.String())
	}

	return stdout.String(), nil
}
