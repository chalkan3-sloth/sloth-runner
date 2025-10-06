package events

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewEventsCommand creates the events command
func NewEventsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Manage event queue",
		Long:  `The events command provides subcommands to list, view, and manage event queue.`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCommand(ctx))
	cmd.AddCommand(NewGetCommand(ctx))
	cmd.AddCommand(NewShowCommand(ctx))
	cmd.AddCommand(NewDeleteCommand(ctx))
	cmd.AddCommand(NewCleanupCommand(ctx))

	return cmd
}
