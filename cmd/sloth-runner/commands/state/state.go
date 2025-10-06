package state

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewStateCommand creates the parent state command
// This command groups all state management subcommands
func NewStateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Manage state and idempotency tracking",
		Long:  `The state command provides subcommands to view, list, and manage resource state for idempotent operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all state subcommands
	cmd.AddCommand(
		NewListCommand(ctx),
		NewShowCommand(ctx),
		NewDeleteCommand(ctx),
		NewClearCommand(ctx),
		NewStatsCommand(ctx),
	)

	return cmd
}
