package agent

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/spf13/cobra"
)

// addMasterFlag adds the --master flag to a command with the correct default value
func addMasterFlag(cmd *cobra.Command) {
	cmd.Flags().String("master", config.GetMasterAddress(), "Master server address")
}

// getMasterAddress gets the master address from flags or config
func getMasterAddress(cmd *cobra.Command) string {
	if masterAddr, _ := cmd.Flags().GetString("master"); masterAddr != "" {
		return masterAddr
	}
	return config.GetMasterAddress()
}

// formatBytes formats bytes to human-readable format
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
