package agent

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	_ "github.com/mattn/go-sqlite3"
)

// NewGetCommand creates the agent get command
func NewGetCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <agent_name>",
		Short: "Get detailed information about an agent",
		Long:  `Retrieves detailed system information collected from a specific agent. By default, tries to read from local database first, then falls back to master server if specified.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			outputFormat, _ := cmd.Flags().GetString("output")
			local, _ := cmd.Flags().GetBool("local")
			debug, _ := cmd.Flags().GetBool("debug")

			if debug {
				pterm.DefaultLogger.Level = pterm.LogLevelDebug
				slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			}

			// If --local flag is set, use local database
			if local {
				return getAgentInfoFromLocalDB(agentName, outputFormat, debug)
			}

			// Get master address (supports both names and addresses)
			masterAddr := getMasterAddress(cmd)

			// If no master address, use local database
			if masterAddr == "" {
				if debug {
					slog.Debug("No master address configured, using local database")
				}
				return getAgentInfoFromLocalDB(agentName, outputFormat, debug)
			}

			return getAgentInfo(agentName, masterAddr, outputFormat)
		},
	}

	addMasterFlag(cmd)
	cmd.Flags().StringP("output", "o", "text", "Output format: text or json")
	cmd.Flags().Bool("local", false, "Force reading from local database")
	cmd.Flags().Bool("debug", false, "Enable debug logging")

	return cmd
}

func getAgentInfo(agentName, masterAddr, outputFormat string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create connection factory and get client
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	// Use refactored function with injected client
	opts := GetAgentInfoOptions{
		AgentName:    agentName,
		OutputFormat: outputFormat,
		Writer:       os.Stdout,
	}

	return getAgentInfoWithClient(ctx, client, opts)
}

// getAgentInfoFromLocalDB reads agent info directly from the local SQLite database
func getAgentInfoFromLocalDB(agentName, outputFormat string, debug bool) error {
	// Database path
	dbPath := config.GetAgentDBPath()

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("no local agent database found at: %s", dbPath)
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Query agent
	query := `SELECT id, name, address, status, last_heartbeat, registered_at, updated_at,
			  last_info_collected, COALESCE(system_info, ''), COALESCE(version, '')
			  FROM agents WHERE name = ? LIMIT 1`

	var id int
	var name, address, status string
	var lastHeartbeat, registeredAt, updatedAt, lastInfoCollected int64
	var systemInfoJSON, version string

	err = db.QueryRow(query, agentName).Scan(&id, &name, &address, &status, &lastHeartbeat,
		&registeredAt, &updatedAt, &lastInfoCollected, &systemInfoJSON, &version)
	if err == sql.ErrNoRows {
		return fmt.Errorf("agent '%s' not found in local database", agentName)
	}
	if err != nil {
		return fmt.Errorf("failed to query agent: %w", err)
	}

	if debug {
		slog.Debug("Agent found in local database", "name", name, "address", address)
	}

	// Update status based on heartbeat (agent is active if heartbeat within last 60 seconds)
	now := time.Now().Unix()
	if lastHeartbeat > 0 && now-lastHeartbeat < 60 {
		status = "Active"
	} else {
		status = "Inactive"
	}

	// Create AgentInfo object for formatting
	agentInfo := &pb.AgentInfo{
		AgentName:         name,
		AgentAddress:      address,
		Status:            status,
		LastHeartbeat:     lastHeartbeat,
		LastInfoCollected: lastInfoCollected,
		SystemInfoJson:    systemInfoJSON,
		Version:           version,
	}

	// Use existing formatting functions
	opts := GetAgentInfoOptions{
		AgentName:    agentName,
		OutputFormat: outputFormat,
		Writer:       os.Stdout,
	}

	if outputFormat == "json" {
		return formatAgentInfoJSON(agentInfo, opts.Writer)
	}

	return formatAgentInfoText(agentInfo, opts.Writer)
}
