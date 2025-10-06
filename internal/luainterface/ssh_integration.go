package luainterface

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/ssh"
)

var (
	// Global SSH executor for remote command execution
	globalSSHExecutor  *ssh.Executor
	globalSSHProfile   string
	globalSSHPassword  *string
)

// SetSSHExecutor sets the global SSH executor for remote command execution
func SetSSHExecutor(executor *ssh.Executor, profile string, password *string) {
	globalSSHExecutor = executor
	globalSSHProfile = profile
	globalSSHPassword = password

	slog.Info("SSH executor configured for remote execution", "profile", profile)
}

// ClearSSHExecutor clears the global SSH executor
func ClearSSHExecutor() {
	globalSSHExecutor = nil
	globalSSHProfile = ""
	globalSSHPassword = nil
}

// ExecuteCommandWithSSH executes a command either locally or via SSH
func ExecuteCommandWithSSH(command string) (string, error) {
	// If no SSH executor is set, execute locally
	if globalSSHExecutor == nil {
		return ExecuteLocalCommand(command)
	}

	// Execute remotely via SSH
	slog.Debug("Executing command via SSH", "profile", globalSSHProfile, "command", command)

	result, err := globalSSHExecutor.ExecuteCommand(globalSSHProfile, command, globalSSHPassword)
	if err != nil {
		return "", fmt.Errorf("SSH execution failed: %w", err)
	}

	// Combine stdout and stderr
	output := result.Output
	if result.Error != "" {
		if output != "" {
			output += "\n"
		}
		output += result.Error
	}

	// Check exit code
	if result.ExitCode != 0 {
		return output, fmt.Errorf("command exited with code %d", result.ExitCode)
	}

	return output, nil
}

// IsSSHExecutionEnabled checks if SSH execution is configured
func IsSSHExecutionEnabled() bool {
	return globalSSHExecutor != nil
}

// GetSSHProfile returns the current SSH profile name
func GetSSHProfile() string {
	return globalSSHProfile
}

// ExecuteLocalCommand executes a command locally (helper function)
func ExecuteLocalCommand(command string) (string, error) {
	cmd := ExecCommand("bash", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Return output even on error for debugging
		return string(output), err
	}

	return strings.TrimSpace(string(output)), nil
}