package workflow

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewWorkflowCommand creates the workflow parent command
func NewWorkflowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows",
		Long:  `Manage workflows including running, listing, and previewing workflow files.`,
	}

	// Add subcommands
	cmd.AddCommand(NewRunCommand(ctx))
	cmd.AddCommand(NewListCommand(ctx))
	cmd.AddCommand(NewPreviewCommand(ctx))

	return cmd
}
