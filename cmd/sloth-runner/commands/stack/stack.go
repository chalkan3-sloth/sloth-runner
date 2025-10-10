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
		// Basic stack management
		NewListCommand(ctx),
		NewShowCommand(ctx),
		NewNewCommand(ctx),
		NewDeleteCommand(ctx),
		NewHistoryCommand(ctx),

		// State management (Pulumi/Terraform-like)
		NewStateCommand(ctx),
		NewMigrateCommand(ctx),

		// Operations tracking
		NewOperationsCommand(ctx),

		// Advanced state management
		NewSnapshotCommand(ctx),   // Snapshot management (create, list, restore, compare)
		NewDriftCommand(ctx),      // Drift detection and auto-fix
		NewLockCommand(ctx),       // State locking (prevent concurrent modifications)
		NewValidateCommand(ctx),   // State validation and repair
		NewEventsCommand(ctx),     // Event viewing and statistics
		NewDepsCommand(ctx),       // Dependency graph visualization and analysis
	)

	return cmd
}
