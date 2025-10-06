package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewListCommand creates the agent list command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all registered agents",
		Long:  `Lists all agents that are currently registered with the master.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				pterm.DefaultLogger.Level = pterm.LogLevelDebug
				slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			}

			masterAddr, _ := cmd.Flags().GetString("master")

			// Create context
			ctx := context.Background()

			// Create connection factory and get client
			factory := NewDefaultConnectionFactory()
			client, cleanup, err := factory.CreateRegistryClient(masterAddr)
			if err != nil {
				return fmt.Errorf("failed to connect to master: %w", err)
			}
			defer cleanup()

			// Use refactored function with injected client
			opts := ListAgentsOptions{
				Writer: os.Stdout,
			}

			return listAgentsWithClient(ctx, client, opts)
		},
	}

	cmd.Flags().String("master", "192.168.1.29:50053", "Master registry address")
	cmd.Flags().Bool("debug", false, "Enable debug logging")

	return cmd
}
