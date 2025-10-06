package scheduler

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewSchedulerCommand creates the parent scheduler command
// This command groups all scheduler-related subcommands
func NewSchedulerCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "Manage workflow scheduling",
		Long:  `The scheduler command provides subcommands to manage scheduled workflow executions.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all scheduler subcommands
	cmd.AddCommand(
		NewEnableCommand(ctx),
		NewDisableCommand(ctx),
		NewListCommand(ctx),
		NewDeleteCommand(ctx),
	)

	return cmd
}
