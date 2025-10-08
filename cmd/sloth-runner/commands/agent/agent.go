package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewAgentCommand creates the parent agent command
// This command groups all agent-related subcommands
func NewAgentCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manages sloth-runner agents",
		Long:  `The agent command provides subcommands to start, stop, list, and manage sloth-runner agents.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all agent subcommands
	cmd.AddCommand(
		NewStartCommand(ctx),
		NewStopCommand(ctx),
		NewListCommand(ctx),
		NewDeleteCommand(ctx),
		NewGetCommand(ctx),
		NewExecCommand(ctx),
		NewModulesCommand(ctx),
		NewMetricsCommand(ctx),
		NewUpdateCommand(ctx),
		NewInstallCommand(ctx),
		NewDocsCommand(ctx),
		NewShellCommand(ctx),
		NewWatcherCommand(ctx),
	)

	return cmd
}
