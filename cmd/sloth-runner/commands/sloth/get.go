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

// NewGetCommand creates the 'sloth get' command
func NewGetCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get details of a specific sloth",
		Long:  `Get detailed information about a specific .sloth file saved in the database.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			showContent, _ := cmd.Flags().GetBool("content")

			// Create service
			service, err := services.NewSlothService()
			if err != nil {
				return err
			}
			defer service.Close()

			// Get sloth
			sloth, err := service.GetSloth(cmd.Context(), name)
			if err != nil {
				return fmt.Errorf("failed to get sloth: %w", err)
			}

			// Display details
			pterm.DefaultHeader.Printf("Sloth: %s", sloth.Name)
			pterm.Println()

			pterm.Info.Printf("ID: %s\n", pterm.Gray(sloth.ID))
			pterm.Info.Printf("Name: %s\n", pterm.Cyan(sloth.Name))
			pterm.Info.Printf("Description: %s\n", sloth.Description)
			pterm.Info.Printf("File Path: %s\n", pterm.Gray(sloth.FilePath))

			activeStatus := "No"
			if sloth.IsActive {
				activeStatus = pterm.Green("Yes")
			} else {
				activeStatus = pterm.Red("No")
			}
			pterm.Info.Printf("Active: %s\n", activeStatus)

			pterm.Info.Printf("Created: %s\n", sloth.CreatedAt.Format("2006-01-02 15:04:05"))
			pterm.Info.Printf("Updated: %s\n", sloth.UpdatedAt.Format("2006-01-02 15:04:05"))

			if sloth.LastUsedAt != nil {
				pterm.Info.Printf("Last Used: %s\n", sloth.LastUsedAt.Format("2006-01-02 15:04:05"))
			} else {
				pterm.Info.Println("Last Used: Never")
			}

			pterm.Info.Printf("Usage Count: %d\n", sloth.UsageCount)
			pterm.Info.Printf("File Hash: %s\n", pterm.Gray(sloth.FileHash))

			if sloth.Tags != "" {
				pterm.Info.Printf("Tags: %s\n", sloth.Tags)
			}

			// Show content if requested
			if showContent {
				pterm.Println()
				pterm.DefaultBox.WithTitle("Content").Println(sloth.Content)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("content", "c", false, "Show file content")

	return cmd
}
