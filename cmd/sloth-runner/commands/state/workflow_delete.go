//go:build cgo
// +build cgo

package state

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewWorkflowDeleteCommand creates the workflow delete command
func NewWorkflowDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <workflow-id>",
		Short: "Delete a workflow state",
		Long:  `Deletes a workflow state and all associated resources, outputs, and versions. This action is irreversible.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			force, _ := cmd.Flags().GetBool("force")

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			if err := sm.InitWorkflowSchema(); err != nil {
				return fmt.Errorf("failed to initialize workflow schema: %w", err)
			}

			// Get workflow to display info
			workflow, err := sm.GetWorkflowState(workflowID)
			if err != nil {
				return fmt.Errorf("workflow not found: %s", workflowID)
			}

			if !force {
				pterm.Warning.Printfln("About to delete workflow state:")
				fmt.Printf("  Name:      %s\n", workflow.Name)
				fmt.Printf("  Version:   %d\n", workflow.Version)
				fmt.Printf("  Status:    %s\n", workflow.Status)
				fmt.Printf("  Resources: %d\n", len(workflow.Resources))
				fmt.Println()
				pterm.Error.Println("This action is IRREVERSIBLE and will delete all state, resources, and versions!")
				fmt.Println()

				confirm, _ := pterm.DefaultInteractiveConfirm.Show("Are you absolutely sure?")
				if !confirm {
					pterm.Info.Println("Delete cancelled")
					return nil
				}
			}

			spinner, _ := pterm.DefaultSpinner.Start("Deleting workflow state...")

			if err := sm.DeleteWorkflowState(workflowID); err != nil {
				spinner.Fail("Delete failed")
				return fmt.Errorf("failed to delete workflow state: %w", err)
			}

			spinner.Success(fmt.Sprintf("Workflow state '%s' deleted successfully", workflow.Name))

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}
