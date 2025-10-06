package sloth

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewSlothCommand creates the main sloth command
func NewSlothCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sloth",
		Short: "Manage saved .sloth files",
		Long: `Manage .sloth files stored in the database.
Sloths are saved workflow files that can be reused across multiple runs
without specifying the file path each time.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(NewAddCommand(ctx))
	cmd.AddCommand(NewListCommand(ctx))
	cmd.AddCommand(NewGetCommand(ctx))
	cmd.AddCommand(NewRemoveCommand(ctx))
	cmd.AddCommand(NewDeleteCommand(ctx))
	cmd.AddCommand(NewActivateCommand(ctx))
	cmd.AddCommand(NewDeactivateCommand(ctx))

	return cmd
}
