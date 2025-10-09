package agent

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/chalkan3-sloth/sloth-runner/internal/masterdb"
	"github.com/spf13/cobra"
	_ "github.com/mattn/go-sqlite3"
)

// addMasterFlag adds the --master flag to a command with the correct default value
func addMasterFlag(cmd *cobra.Command) {
	cmd.Flags().String("master", "", "Master server name or address (e.g., 'production' or '192.168.1.29:50053')")
}

// getMasterAddress gets the master address from flags or config
// Supports both master names (looked up from database) and direct addresses
func getMasterAddress(cmd *cobra.Command) string {
	masterFlag, _ := cmd.Flags().GetString("master")

	// If flag is empty, try to get default from database
	if masterFlag == "" {
		return getDefaultMasterAddress()
	}

	// If it contains a colon, it's an address - use it directly
	if strings.Contains(masterFlag, ":") {
		return masterFlag
	}

	// Otherwise, it's a name - look it up in the database
	db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
	if err != nil {
		// If can't open DB, fallback to treating it as address or use default
		return config.GetMasterAddress()
	}
	defer db.Close()

	master, err := db.Get(masterFlag)
	if err != nil {
		// If not found in DB, fallback to config default
		return config.GetMasterAddress()
	}

	return master.Address
}

// getDefaultMasterAddress returns the default master address from database or config
func getDefaultMasterAddress() string {
	// Try to get default from database first
	db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
	if err != nil {
		// Fallback to config file/env
		return config.GetMasterAddress()
	}
	defer db.Close()

	master, err := db.GetDefault()
	if err != nil {
		// No default in DB, fallback to config
		return config.GetMasterAddress()
	}

	return master.Address
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

// getAgentAddressFromLocalDB retrieves agent address from local database
// This is used when --local flag is set to connect directly to agent
func getAgentAddressFromLocalDB(agentName string) (string, error) {
	dbPath := config.GetAgentDBPath()

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("no local agent database found at: %s", dbPath)
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Query agent address
	var address string
	err = db.QueryRow("SELECT address FROM agents WHERE name = ? LIMIT 1", agentName).Scan(&address)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("agent '%s' not found in local database", agentName)
	}
	if err != nil {
		return "", fmt.Errorf("failed to query agent: %w", err)
	}

	return address, nil
}
