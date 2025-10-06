package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	agentInternal "github.com/chalkan3-sloth/sloth-runner/internal/agent"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
)

// StartAgentOptions contains configuration for starting an agent
type StartAgentOptions struct {
	Port             int
	MasterAddr       string
	AgentName        string
	Daemon           bool
	BindAddress      string
	ReportAddress    string
	TelemetryEnabled bool
	MetricsPort      int
}

// DaemonProcessInfo contains information about a running daemon process
type DaemonProcessInfo struct {
	PID     int
	PIDFile string
	Running bool
}

// AddressConfig contains addresses determined from flags
type AddressConfig struct {
	ListenAddr string
	ReportAddr string
}

// RegistrationResult contains the result of master registration
type RegistrationResult struct {
	Success       bool
	Error         error
	MasterAddr    string
	ReportAddress string
}

// HeartbeatResult contains the result of a heartbeat attempt
type HeartbeatResult struct {
	Success bool
	Error   error
}

// checkExistingAgent checks if an agent is already running based on PID file (testable)
func checkExistingAgent(agentName string) (*DaemonProcessInfo, error) {
	pidFile := filepath.Join("/tmp", fmt.Sprintf("sloth-runner-agent-%s.pid", agentName))

	info := &DaemonProcessInfo{
		PIDFile: pidFile,
		Running: false,
	}

	// Check if PID file exists
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return info, nil
	}

	// Read PID from file
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		return info, fmt.Errorf("failed to read PID file: %w", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil {
		return info, fmt.Errorf("invalid PID in file: %w", err)
	}

	info.PID = pid

	// Check if process is running
	process, err := os.FindProcess(pid)
	if err != nil {
		// Process not found
		return info, nil
	}

	// Send signal 0 to check if process exists (doesn't actually send signal)
	if err := process.Signal(syscall.Signal(0)); err == nil {
		// Process is running
		info.Running = true
	}

	return info, nil
}

// buildDaemonCommandArgs builds command arguments for daemon mode (testable)
func buildDaemonCommandArgs(opts StartAgentOptions) []string {
	args := []string{
		"agent", "start",
		"--port", strconv.Itoa(opts.Port),
		"--name", opts.AgentName,
		"--master", opts.MasterAddr,
	}

	if opts.BindAddress != "" {
		args = append(args, "--bind-address", opts.BindAddress)
	}

	if opts.ReportAddress != "" {
		args = append(args, "--report-address", opts.ReportAddress)
	}

	if opts.TelemetryEnabled {
		args = append(args, "--telemetry")
	}

	if opts.MetricsPort != 9090 {
		args = append(args, "--metrics-port", strconv.Itoa(opts.MetricsPort))
	}

	return args
}

// writePIDFile writes a PID to the PID file (testable)
func writePIDFile(pidFile string, pid int) error {
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

// removePIDFile removes a PID file (testable)
func removePIDFile(pidFile string) error {
	if _, err := os.Stat(pidFile); err == nil {
		return os.Remove(pidFile)
	}
	return nil
}

// determineListenAddress determines the listen address from options (testable)
func determineListenAddress(opts StartAgentOptions) string {
	if opts.BindAddress != "" {
		return fmt.Sprintf("%s:%d", opts.BindAddress, opts.Port)
	}
	return fmt.Sprintf(":%d", opts.Port)
}

// determineReportAddress determines the report address for master registration (testable)
func determineReportAddress(opts StartAgentOptions, actualListenAddr string) string {
	if opts.ReportAddress != "" {
		// If report address has no port, add the port
		if !strings.Contains(opts.ReportAddress, ":") {
			return fmt.Sprintf("%s:%d", opts.ReportAddress, opts.Port)
		}
		return opts.ReportAddress
	}

	if opts.BindAddress != "" {
		return fmt.Sprintf("%s:%d", opts.BindAddress, opts.Port)
	}

	return actualListenAddr
}

// registerAgentWithMaster registers the agent with the master server (testable)
func registerAgentWithMaster(ctx context.Context, client AgentRegistryClient, agentName, reportAddress string) error {
	regCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := client.RegisterAgent(regCtx, &pb.RegisterAgentRequest{
		AgentName:    agentName,
		AgentAddress: reportAddress,
	})

	return err
}

// sendHeartbeatWithSystemInfo sends a heartbeat with optional system info (testable)
func sendHeartbeatWithSystemInfo(ctx context.Context, client AgentRegistryClient, agentName string, includeSystemInfo bool) error {
	var sysInfoJSON string

	if includeSystemInfo {
		if sysInfo, err := agentInternal.CollectSystemInfo(); err == nil {
			if jsonStr, err := sysInfo.ToJSON(); err == nil {
				sysInfoJSON = jsonStr
			}
		}
	}

	hbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := client.Heartbeat(hbCtx, &pb.HeartbeatRequest{
		AgentName:      agentName,
		SystemInfoJson: sysInfoJSON,
	})

	return err
}

// startDaemonProcess starts the agent as a daemon process (uses exec, harder to test but extracted)
func startDaemonProcess(opts StartAgentOptions) (int, error) {
	cmdArgs := buildDaemonCommandArgs(opts)
	command := exec.Command(os.Args[0], cmdArgs...)

	// Open log file for stdout
	stdoutFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open agent.log for stdout: %w", err)
	}
	defer stdoutFile.Close()
	command.Stdout = stdoutFile

	// Open log file for stderr
	stderrFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open agent.log for stderr: %w", err)
	}
	defer stderrFile.Close()
	command.Stderr = stderrFile

	// Start the process
	if err := command.Start(); err != nil {
		return 0, fmt.Errorf("failed to start agent process: %w", err)
	}

	return command.Process.Pid, nil
}
