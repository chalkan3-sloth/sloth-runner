//go:build !cgo
// +build !cgo

package workflow

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewWorkflowCommand creates the workflow parent command for non-CGO builds
func NewWorkflowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows (limited functionality without CGO)",
		Long:  `Manage workflows including listing and previewing workflow files. Running workflows with state management requires CGO support.`,
	}

	// Add subcommands that don't require CGO
	cmd.AddCommand(NewListCommand(ctx))
	cmd.AddCommand(NewPreviewCommand(ctx))

	// Add a stub run command that returns an error
	cmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run a workflow (requires CGO)",
		Long: `Running workflows with state management requires CGO support (SQLite).

This binary was compiled without CGO for portability.
To run workflows with state management:
  1. Download the CGO-enabled version for Linux from GitHub releases
  2. Or compile from source with CGO_ENABLED=1

For basic task execution without state management, you can use the agent commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("workflow run requires CGO support (SQLite).\n" +
				"This binary was compiled without CGO for portability.\n\n" +
				"Please use a CGO-enabled build (available for Linux) or compile with CGO_ENABLED=1")
		},
	})

	return cmd
}
