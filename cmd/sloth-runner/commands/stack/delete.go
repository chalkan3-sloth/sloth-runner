//go:build cgo
// +build cgo

package stack

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates the stack delete command
func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <stack-name>",
		Short: "Delete a workflow stack",
		Long:  `Delete a workflow stack and all its execution history.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			force, _ := cmd.Flags().GetBool("force")

			stackManager, err := stack.NewStackManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize stack manager: %w", err)
			}
			defer stackManager.Close()

			stackState, err := stackManager.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			if !force {
				pterm.Warning.Printf("This will permanently delete stack '%s' and all its execution history.\n", stackName)
				confirm := pterm.DefaultInteractiveConfirm.WithDefaultValue(false)
				result, _ := confirm.Show("Are you sure?")
				if !result {
					pterm.Info.Println("Operation cancelled.")
					return nil
				}
			}

			if err := stackManager.DeleteStack(stackState.ID); err != nil {
				return fmt.Errorf("failed to delete stack: %w", err)
			}

			pterm.Success.Printf("Stack '%s' deleted successfully.\n", stackName)
			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")
	return cmd
}
