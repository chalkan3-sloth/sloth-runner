package sloth

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewListCommand creates the 'sloth list' command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved sloths",
		Long:  `List all .sloth files saved in the database with their details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			activeOnly, _ := cmd.Flags().GetBool("active")

			// Create service
			service, err := services.NewSlothService()
			if err != nil {
				return err
			}
			defer service.Close()

			// List sloths
			sloths, err := service.ListSloths(cmd.Context(), activeOnly)
			if err != nil {
				return err
			}

			if len(sloths) == 0 {
				if activeOnly {
					pterm.Info.Println("No active sloths found")
				} else {
					pterm.Info.Println("No sloths found")
				}
				return nil
			}

			// Display header
			if activeOnly {
				pterm.DefaultHeader.Println("Active Sloths")
			} else {
				pterm.DefaultHeader.Println("All Sloths")
			}
			pterm.Println()

			// Build table data
			tableData := pterm.TableData{
				{"Name", "Description", "Active", "Usage", "Last Used", "Created"},
			}

			for _, s := range sloths {
				activeStatus := "✗"
				if s.IsActive {
					activeStatus = "✓"
				}

				lastUsed := "Never"
				if s.LastUsedAt != nil {
					lastUsed = s.LastUsedAt.Format("2006-01-02 15:04")
				}

				description := s.Description
				if description == "" {
					description = "-"
				}

				tableData = append(tableData, []string{
					s.Name,
					description,
					activeStatus,
					pterm.Sprintf("%d", s.UsageCount),
					lastUsed,
					s.CreatedAt.Format("2006-01-02"),
				})
			}

			// Render table
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			return nil
		},
	}

	cmd.Flags().BoolP("active", "a", false, "Show only active sloths")

	return cmd
}
