package state

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewWorkflowRollbackCommand creates the workflow rollback command
func NewWorkflowRollbackCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback <workflow-id> <version>",
		Short: "Rollback workflow to a previous version",
		Long:  `Rolls back a workflow state to a specific version. This creates a new version with the state from the specified version.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			version := 0
			if _, err := fmt.Sscanf(args[1], "%d", &version); err != nil {
				return fmt.Errorf("invalid version number: %s", args[1])
			}

			force, _ := cmd.Flags().GetBool("force")

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			if err := sm.InitWorkflowSchema(); err != nil {
				return fmt.Errorf("failed to initialize workflow schema: %w", err)
			}

			// Check current state
			current, err := sm.GetWorkflowState(workflowID)
			if err != nil {
				return fmt.Errorf("failed to get current workflow state: %w", err)
			}

			if !force {
				pterm.Warning.Printfln("About to rollback workflow '%s' from version %d to version %d",
					current.Name, current.Version, version)

				confirm, _ := pterm.DefaultInteractiveConfirm.Show("Do you want to continue?")
				if !confirm {
					pterm.Info.Println("Rollback cancelled")
					return nil
				}
			}

			spinner, _ := pterm.DefaultSpinner.Start("Rolling back workflow...")

			if err := sm.RollbackToVersion(workflowID, version); err != nil {
				spinner.Fail("Rollback failed")
				return fmt.Errorf("failed to rollback: %w", err)
			}

			spinner.Success("Rollback completed successfully")

			// Get updated state
			updated, err := sm.GetWorkflowState(workflowID)
			if err == nil {
				pterm.Success.Printfln("\nWorkflow '%s' rolled back to version %d (new version: %d)",
					updated.Name, version, updated.Version)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}
