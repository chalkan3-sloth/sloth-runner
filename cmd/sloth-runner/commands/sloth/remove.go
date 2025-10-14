//go:build cgo
// +build cgo

package sloth

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewRemoveCommand creates the 'sloth remove' command
func NewRemoveCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a sloth from the database",
		Long:  `Remove a .sloth file from the database. This action cannot be undone.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			force, _ := cmd.Flags().GetBool("force")

			// Create service
			service, err := services.NewSlothService()
			if err != nil {
				return err
			}
			defer service.Close()

			// Confirm if not forced
			if !force {
				result, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText(fmt.Sprintf("Are you sure you want to remove sloth '%s'?", name)).
					Show()

				if !result {
					pterm.Info.Println("Operation cancelled")
					return nil
				}
			}

			// Remove sloth
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing sloth '%s'...", name))

			if err := service.RemoveSloth(cmd.Context(), name); err != nil {
				spinner.Fail(fmt.Sprintf("Failed to remove sloth: %v", err))
				trackSlothDelete(name, false)
				return err
			}

			spinner.Success(fmt.Sprintf("✓ Sloth '%s' removed successfully", name))

			// Track operation
			trackSlothDelete(name, true)

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}

// NewDeleteCommand creates the 'sloth delete' command (alias for remove)
func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a sloth from the database (alias for remove)",
		Long:  `Delete a .sloth file from the database. This action cannot be undone.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			force, _ := cmd.Flags().GetBool("force")

			// Create service
			service, err := services.NewSlothService()
			if err != nil {
				return err
			}
			defer service.Close()

			// Confirm if not forced
			if !force {
				result, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText(fmt.Sprintf("Are you sure you want to delete sloth '%s'?", name)).
					Show()

				if !result {
					pterm.Info.Println("Operation cancelled")
					return nil
				}
			}

			// Delete sloth
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deleting sloth '%s'...", name))

			if err := service.DeleteSloth(cmd.Context(), name); err != nil {
				spinner.Fail(fmt.Sprintf("Failed to delete sloth: %v", err))
				trackSlothDelete(name, false)
				return err
			}

			spinner.Success(fmt.Sprintf("✓ Sloth '%s' deleted successfully", name))

			// Track operation
			trackSlothDelete(name, true)

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}
