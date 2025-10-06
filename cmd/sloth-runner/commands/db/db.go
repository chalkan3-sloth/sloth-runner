package db

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewDBCommand creates the parent db command
func NewDBCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Interact with sloth-runner databases",
		Long: `The db command provides tools to query and inspect sloth-runner databases.

You can query agents.db (agent information) or hooks.db (hooks and events).`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all db subcommands
	cmd.AddCommand(
		NewQueryCommand(ctx),
		NewTablesCommand(ctx),
		NewSchemaCommand(ctx),
	)

	return cmd
}
