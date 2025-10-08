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
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	_ "github.com/mattn/go-sqlite3"
)

// NewListCommand creates the agent list command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all registered agents",
		Long:  `Lists all agents that are currently registered. By default, tries to read from local database first, then falls back to master server if specified.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				pterm.DefaultLogger.Level = pterm.LogLevelDebug
				slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			}

			local, _ := cmd.Flags().GetBool("local")

			// If --local flag is set, use local database
			if local {
				return listAgentsFromLocalDB(debug)
			}

			// Get master address (supports both names and addresses)
			masterAddr := getMasterAddress(cmd)

			// If no master address, use local database
			if masterAddr == "" {
				if debug {
					slog.Debug("No master address configured, using local database")
				}
				return listAgentsFromLocalDB(debug)
			}

			// Try master server first, fallback to local DB if it fails
			ctx := context.Background()

			// Create connection factory and get client with timeout
			factory := NewDefaultConnectionFactory()
			client, cleanup, err := factory.CreateRegistryClient(masterAddr)
			if err != nil {
				if debug {
					slog.Debug("Failed to connect to master, falling back to local database", "error", err)
				}
				pterm.Warning.Printf("Could not connect to master at %s, using local database\n", masterAddr)
				return listAgentsFromLocalDB(debug)
			}
			defer cleanup()

			// Use refactored function with injected client
			opts := ListAgentsOptions{
				Writer: os.Stdout,
			}

			return listAgentsWithClient(ctx, client, opts)
		},
	}

	cmd.Flags().String("master", "", "Master registry address (if empty, uses local database)")
	cmd.Flags().Bool("local", false, "Force reading from local database")
	cmd.Flags().Bool("debug", false, "Enable debug logging")

	return cmd
}

// listAgentsFromLocalDB reads agents directly from the local SQLite database
func listAgentsFromLocalDB(debug bool) error {
	// Database path
	dbPath := config.GetAgentDBPath()

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		pterm.Warning.Println("No local agent database found")
		pterm.Info.Printf("Expected database at: %s\n", dbPath)
		return nil
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Query agents
	query := `SELECT id, name, address, status, last_heartbeat, registered_at, updated_at,
			  last_info_collected, COALESCE(system_info, ''), COALESCE(version, '')
			  FROM agents ORDER BY name`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query agents: %w", err)
	}
	defer rows.Close()

	// Display header
	pterm.DefaultSection.Println("Registered Agents (from local database)")

	var agents [][]string
	now := time.Now().Unix()
	count := 0

	for rows.Next() {
		var id int
		var name, address, status string
		var lastHeartbeat, registeredAt, updatedAt, lastInfoCollected int64
		var systemInfo, version string

		err := rows.Scan(&id, &name, &address, &status, &lastHeartbeat, &registeredAt, &updatedAt, &lastInfoCollected, &systemInfo, &version)
		if err != nil {
			return fmt.Errorf("failed to scan agent row: %w", err)
		}

		// Update status based on heartbeat (agent is active if heartbeat within last 60 seconds)
		if lastHeartbeat > 0 && now-lastHeartbeat < 60 {
			status = "Active"
		} else {
			status = "Inactive"
		}

		// Format last heartbeat
		lastHB := "Never"
		if lastHeartbeat > 0 {
			hbTime := time.Unix(lastHeartbeat, 0)
			lastHB = fmt.Sprintf("%s (%s ago)", hbTime.Format("15:04:05"), time.Since(hbTime).Round(time.Second))
		}

		// Format version
		if version == "" {
			version = "Unknown"
		}

		agents = append(agents, []string{
			name,
			address,
			status,
			lastHB,
			version,
		})
		count++
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating agent rows: %w", err)
	}

	if count == 0 {
		pterm.Info.Println("No agents registered")
		return nil
	}

	// Display agents in table format
	pterm.DefaultTable.WithHasHeader().WithData(pterm.TableData{
		{"Agent Name", "Address", "Status", "Last Heartbeat", "Version"},
	}).WithData(agents).Render()

	pterm.Info.Printf("\nTotal agents: %d\n", count)

	return nil
}
