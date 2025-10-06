package sloth

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewActivateCommand creates the 'sloth activate' command
func NewActivateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate <name>",
		Short: "Activate a sloth",
		Long:  `Set a sloth as active. Active sloths can be used in run commands.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Create service
			service, err := services.NewSlothService()
			if err != nil {
				return err
			}
			defer service.Close()

			// Activate sloth
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Activating sloth '%s'...", name))

			if err := service.ActivateSloth(cmd.Context(), name); err != nil {
				spinner.Fail(fmt.Sprintf("Failed to activate sloth: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("✓ Sloth '%s' is now active", name))

			return nil
		},
	}

	return cmd
}

// NewDeactivateCommand creates the 'sloth deactivate' command
func NewDeactivateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate <name>",
		Short: "Deactivate a sloth",
		Long:  `Set a sloth as inactive. Inactive sloths cannot be used in run commands.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Create service
			service, err := services.NewSlothService()
			if err != nil {
				return err
			}
			defer service.Close()

			// Deactivate sloth
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deactivating sloth '%s'...", name))

			if err := service.DeactivateSloth(cmd.Context(), name); err != nil {
				spinner.Fail(fmt.Sprintf("Failed to deactivate sloth: %v", err))
				return err
			}

			spinner.Warning(fmt.Sprintf("⚠ Sloth '%s' is now inactive", name))

			return nil
		},
	}

	return cmd
}
