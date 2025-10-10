package hook

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates the hook delete command
func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <hook-name>",
		Short: "Delete a hook",
		Long:  `Delete an event hook from the system.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]
			force, _ := cmd.Flags().GetBool("force")

			// Create repository
			repo, err := hooks.NewRepository()
			if err != nil {
				return fmt.Errorf("failed to create repository: %w", err)
			}
			defer repo.Close()

			// Get hook to verify it exists
			hook, err := repo.GetByName(hookName)
			if err != nil {
				return fmt.Errorf("failed to get hook: %w", err)
			}

			// Confirm deletion if not forced
			if !force {
				pterm.Warning.Printf("Are you sure you want to delete hook '%s'?\n", hookName)
				pterm.Info.Printf("Event type: %s\n", hook.EventType)
				pterm.Info.Printf("File: %s\n", hook.FilePath)

				result, _ := pterm.DefaultInteractiveConfirm.Show()
				if !result {
					pterm.Info.Println("Deletion cancelled")
					return nil
				}
			}

			// Delete hook
			if err := repo.Delete(hook.ID); err != nil {
				trackHookDelete(hookName, false)
				return fmt.Errorf("failed to delete hook: %w", err)
			}

			pterm.Success.Printf("Hook '%s' deleted successfully\n", hookName)

			// Track operation
			trackHookDelete(hookName, true)

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")

	return cmd
}
