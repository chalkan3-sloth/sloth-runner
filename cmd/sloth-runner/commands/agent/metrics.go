package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewMetricsCommand creates the agent metrics command
// TODO: Extract logic from main.go agentMetricsCmd and subcommands
func NewMetricsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Manage agent metrics",
		Long:  `Manages metrics collection and export for agents.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// TODO: Add subcommands
	// - prometheus
	// - grafana
	// Extract from main.go lines ~2002-2223

	return cmd
}
