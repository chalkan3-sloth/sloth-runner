package sloth

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewAddCommand creates the 'sloth add' command
func NewAddCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> --file <path>",
		Short: "Add a new .sloth file to the database",
		Long: `Add a new .sloth file to the database for future use.
The file content will be stored and can be referenced by name in run commands.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			filePath, _ := cmd.Flags().GetString("file")
			description, _ := cmd.Flags().GetString("description")
			active, _ := cmd.Flags().GetBool("active")

			if filePath == "" {
				return fmt.Errorf("--file flag is required")
			}

			// Create service
			service, err := services.NewSlothService()
			if err != nil {
				return err
			}
			defer service.Close()

			// Add sloth
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Adding sloth '%s'...", name))

			if err := service.AddSloth(cmd.Context(), name, filePath, description, active); err != nil {
				spinner.Fail(fmt.Sprintf("Failed to add sloth: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("✓ Sloth '%s' added successfully", name))
			pterm.Println()

			// Show details
			pterm.Info.Printf("Name: %s\n", pterm.Cyan(name))
			pterm.Info.Printf("File: %s\n", pterm.Gray(filePath))
			if description != "" {
				pterm.Info.Printf("Description: %s\n", description)
			}
			pterm.Info.Printf("Active: %v\n", active)

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "Path to the .sloth file (required)")
	cmd.Flags().StringP("description", "d", "", "Description of the sloth")
	cmd.Flags().Bool("active", true, "Set sloth as active (default: true)")

	return cmd
}
