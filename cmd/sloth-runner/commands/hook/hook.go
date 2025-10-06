package hook

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewHookCommand creates the parent hook command
func NewHookCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Manage event hooks",
		Long:  `The hook command provides subcommands to add, list, get, show, delete, enable and disable event hooks.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all hook subcommands
	cmd.AddCommand(
		NewAddCommand(ctx),
		NewListCommand(ctx),
		NewGetCommand(ctx),
		NewShowCommand(ctx),
		NewDeleteCommand(ctx),
		NewEnableCommand(ctx),
		NewDisableCommand(ctx),
		NewTestCommand(ctx),
		NewDocsCommand(ctx),
	)

	return cmd
}
