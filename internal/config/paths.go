package config

import (
	"os"
	"path/filepath"
	"runtime"
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

// EnsureDataDir creates the data directory if it doesn't exist
func EnsureDataDir() error {
	dir := GetDataDir()
	return os.MkdirAll(dir, 0755)
}

// GetMasterAddress returns the configured master server address
// Priority: 1. SLOTH_RUNNER_MASTER_ADDR env var, 2. master.conf file, 3. default localhost:50051
func GetMasterAddress() string {
	// First check environment variable
	if addr := os.Getenv("SLOTH_RUNNER_MASTER_ADDR"); addr != "" {
		return addr
	}

	// Then check config file
	configPath := filepath.Join(GetDataDir(), "master.conf")
	if data, err := os.ReadFile(configPath); err == nil && len(data) > 0 {
		return string(data)
	}

	// Default to localhost
	return "localhost:50051"
}
