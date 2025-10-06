package stack

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewStackCommand creates the parent stack command
// This command groups all stack-related subcommands
func NewStackCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack",
		Short: "Manage workflow stacks",
		Long:  `The stack command provides subcommands to manage workflow stacks and their state.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all stack subcommands
	cmd.AddCommand(
		NewListCommand(ctx),
		NewShowCommand(ctx),
		NewNewCommand(ctx),
		NewDeleteCommand(ctx),
		NewHistoryCommand(ctx),
	)

	return cmd
}
