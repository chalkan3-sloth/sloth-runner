package hook

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewListCommand creates the hook list command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all registered hooks",
		Long:  `List all event hooks that are currently registered in the system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			stack, _ := cmd.Flags().GetString("stack")

			// Create repository
			repo, err := hooks.NewRepository()
			if err != nil {
				return fmt.Errorf("failed to create repository: %w", err)
			}
			defer repo.Close()

			// Get hooks (filtered by stack if specified)
			var hookList []*hooks.Hook
			if stack != "" {
				hookList, err = repo.ListByStack(stack)
			} else {
				hookList, err = repo.List()
			}
			if err != nil {
				return fmt.Errorf("failed to list hooks: %w", err)
			}

			if len(hookList) == 0 {
				pterm.Info.Println("No hooks registered")
				return nil
			}

			// Display header
			pterm.DefaultSection.Println("Registered Event Hooks")

			// Prepare table data
			tableData := [][]string{
				{"Name", "Event Type", "Stack", "Status", "Run Count", "Last Run"},
			}

			for _, h := range hookList {
				status := "disabled"
				if h.Enabled {
					status = "enabled"
				}

				lastRun := "Never"
				if h.LastRun != nil {
					lastRun = h.LastRun.Format("2006-01-02 15:04:05")
				}

				stackName := "-"
				if h.Stack != "" {
					stackName = h.Stack
				}

				tableData = append(tableData, []string{
					h.Name,
					string(h.EventType),
					stackName,
					status,
					fmt.Sprintf("%d", h.RunCount),
					lastRun,
				})
			}

			// Display table
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			pterm.Info.Printf("\nTotal hooks: %d\n", len(hookList))

			return nil
		},
	}

	cmd.Flags().StringP("stack", "s", "", "Filter hooks by stack name")

	return cmd
}
