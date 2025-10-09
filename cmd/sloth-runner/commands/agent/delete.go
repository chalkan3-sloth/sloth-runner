package agent

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	_ "github.com/mattn/go-sqlite3"
)

// NewDeleteCommand creates the agent delete command
func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <agent-name>",
		Short: "Deletes an agent from the registry",
		Long:  `Removes an agent from the registry. This does not stop the agent if it's running. By default, tries to delete from local database first, then falls back to master server if specified.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			force, _ := cmd.Flags().GetBool("force")
			local, _ := cmd.Flags().GetBool("local")
			debug, _ := cmd.Flags().GetBool("debug")

			if debug {
				pterm.DefaultLogger.Level = pterm.LogLevelDebug
				slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			}

			// Ask for confirmation unless --force is used
			if !force {
				confirm := false
				prompt := &survey.Confirm{
					Message: fmt.Sprintf("Are you sure you want to delete agent '%s'?", agentName),
					Default: false,
				}
				if err := survey.AskOne(prompt, &confirm); err != nil {
					return err
				}
				if !confirm {
					pterm.Warning.Println("Deletion cancelled")
					return nil
				}
			}

			// If --local flag is set, use local database
			if local {
				return deleteAgentFromLocalDB(agentName, debug)
			}

			// Get master address (supports both names and addresses)
			masterAddr := getMasterAddress(cmd)

			// If no master address, use local database
			if masterAddr == "" {
				if debug {
					slog.Debug("No master address configured, using local database")
				}
				return deleteAgentFromLocalDB(agentName, debug)
			}

			// Create connection factory and get client
			factory := NewDefaultConnectionFactory()
			client, cleanup, err := factory.CreateRegistryClient(masterAddr)
			if err != nil {
				return err
			}
			defer cleanup()

			// Use refactored function with injected client
			opts := DeleteAgentOptions{
				AgentName: agentName,
			}

			return deleteAgentWithClient(context.Background(), client, opts)
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")
	addMasterFlag(cmd)
	cmd.Flags().Bool("local", false, "Force deleting from local database")
	cmd.Flags().Bool("debug", false, "Enable debug logging")

	return cmd
}

// deleteAgentFromLocalDB deletes an agent directly from the local SQLite database
func deleteAgentFromLocalDB(agentName string, debug bool) error {
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

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deleting agent '%s' from local database...", agentName))

	// First check if agent exists
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM agents WHERE name = ?", agentName).Scan(&exists)
	if err != nil {
		spinner.Fail("Failed to check agent existence")
		return fmt.Errorf("failed to check agent: %w", err)
	}

	if exists == 0 {
		spinner.Fail(fmt.Sprintf("Agent '%s' not found", agentName))
		return fmt.Errorf("agent '%s' not found in local database", agentName)
	}

	// Delete the agent
	result, err := db.Exec("DELETE FROM agents WHERE name = ?", agentName)
	if err != nil {
		spinner.Fail("Failed to delete agent")
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		spinner.Fail("Failed to verify deletion")
		return fmt.Errorf("failed to verify deletion: %w", err)
	}

	if rowsAffected == 0 {
		spinner.Fail("No agent was deleted")
		return fmt.Errorf("delete failed: no rows affected")
	}

	if debug {
		slog.Debug("Agent deleted from local database", "name", agentName, "rows_affected", rowsAffected)
	}

	spinner.Success(fmt.Sprintf("Agent '%s' deleted successfully from local database", agentName))
	return nil
}
