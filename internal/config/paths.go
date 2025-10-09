package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetDataDir returns the directory for sloth-runner data files
// On Linux/macOS: /etc/sloth-runner (or $HOME/.sloth-runner if not root)
// On Windows: C:\ProgramData\sloth-runner (or %APPDATA%\sloth-runner if not admin)
func GetDataDir() string {
	// Check for environment variable override
	if dir := os.Getenv("SLOTH_RUNNER_DATA_DIR"); dir != "" {
		return dir
	}

	if runtime.GOOS == "windows" {
		// Windows: Use ProgramData or AppData
		if programData := os.Getenv("PROGRAMDATA"); programData != "" {
			return filepath.Join(programData, "sloth-runner")
		}
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "sloth-runner")
		}
		return filepath.Join("C:", "ProgramData", "sloth-runner")
	}

	// Unix-like systems (Linux, macOS, etc.)
	// If running as root or if /etc/sloth-runner is writable, use it
	if os.Geteuid() == 0 {
		return "/etc/sloth-runner"
	}

	// Check if /etc/sloth-runner exists and is writable
	etcDir := "/etc/sloth-runner"
	if info, err := os.Stat(etcDir); err == nil && info.IsDir() {
		// Check if writable
		testFile := filepath.Join(etcDir, ".write-test")
		if f, err := os.Create(testFile); err == nil {
			f.Close()
			os.Remove(testFile)
			return etcDir
		}
	}

	// Fallback to user's home directory
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".sloth-runner")
	}

	// Ultimate fallback to current directory
	return ".sloth-cache"
}

// GetAgentDBPath returns the full path to the agent database
func GetAgentDBPath() string {
	return filepath.Join(GetDataDir(), "agents.db")
}

// GetSlothDBPath returns the full path to the sloth repository database
func GetSlothDBPath() string {
	return filepath.Join(GetDataDir(), "sloth_repo.db")
}

// GetHookDBPath returns the full path to the hooks database
func GetHookDBPath() string {
	return filepath.Join(GetDataDir(), "hooks.db")
}

// GetSecretsDBPath returns the full path to the secrets database
func GetSecretsDBPath() string {
	return filepath.Join(GetDataDir(), "secrets.db")
}

// GetSSHDBPath returns the full path to the SSH profiles database
func GetSSHDBPath() string {
	return filepath.Join(GetDataDir(), "ssh_profiles.db")
}

// GetStackDBPath returns the full path to the stack database
func GetStackDBPath() string {
	return filepath.Join(GetDataDir(), "stacks.db")
}

// GetMetricsDBPath returns the full path to the metrics database
func GetMetricsDBPath() string {
	return filepath.Join(GetDataDir(), "metrics.db")
}

// GetMastersDBPath returns the full path to the masters database
func GetMastersDBPath() string {
	return filepath.Join(GetDataDir(), "masters.db")
}

// GetHistoryDBPath returns the full path to the execution history database
func GetHistoryDBPath() string {
	return filepath.Join(GetDataDir(), "history.db")
}

// GetLogDir returns the directory for log files
func GetLogDir() string {
	return filepath.Join(GetDataDir(), "logs")
}

// EnsureDataDir creates the data directory if it doesn't exist
func EnsureDataDir() error {
	dir := GetDataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Also ensure logs directory exists
	logDir := GetLogDir()
	return os.MkdirAll(logDir, 0755)
}

// GetMasterAddress returns the configured master server address
// Priority: 1. SLOTH_RUNNER_MASTER_ADDR env var, 2. Default master from DB, 3. master.conf file, 4. default localhost:50051
func GetMasterAddress() string {
	// First check environment variable
	if addr := os.Getenv("SLOTH_RUNNER_MASTER_ADDR"); addr != "" {
		return addr
	}

	// Then check database for default master (lazy import to avoid circular dependency)
	// This will be handled by the caller in most cases

	// Then check config file (legacy support)
	configPath := filepath.Join(GetDataDir(), "master.conf")
	if data, err := os.ReadFile(configPath); err == nil && len(data) > 0 {
		return strings.TrimSpace(string(data))
	}

	// No master configured - return empty string
	// Commands should fallback to local database when empty
	return ""
}

// GetMasterAddressOrName returns the master address, resolving names to addresses
// If the input is a name (no colon), it looks up the address from the database
// If the input is an address (has colon), it returns it directly
func GetMasterAddressOrName(nameOrAddr string) (string, error) {
	// If empty, get default
	if nameOrAddr == "" {
		return GetMasterAddress(), nil
	}

	// If it contains a colon, it's an address
	if filepath.Ext(nameOrAddr) != "" || len(nameOrAddr) > 0 && nameOrAddr[len(nameOrAddr)-1] >= '0' && nameOrAddr[len(nameOrAddr)-1] <= '9' {
		// Simple heuristic: if it looks like an address, return it
		return nameOrAddr, nil
	}

	// Otherwise, treat as a name and look it up
	// This will be implemented by the caller using masterdb
	return "", fmt.Errorf("master name resolution not implemented in config package - use masterdb directly")
}

// ResolveMasterAddress is a helper to get master address with proper error handling
// It should be used by commands that need to resolve master names
func ResolveMasterAddress(nameOrAddr string) (string, error) {
	if nameOrAddr == "" {
		return GetMasterAddress(), nil
	}
	return nameOrAddr, nil
}
