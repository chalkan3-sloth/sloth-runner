package hook

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewShowCommand creates the hook show command
func NewShowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <hook-name>",
		Short: "Show detailed information about a hook",
		Long:  `Show detailed information about a specific hook in human-readable format.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]
			showHistory, _ := cmd.Flags().GetBool("history")
			historyLimit, _ := cmd.Flags().GetInt("limit")

			// Create repository
			repo, err := hooks.NewRepository()
			if err != nil {
				return fmt.Errorf("failed to create repository: %w", err)
			}
			defer repo.Close()

			// Get hook
			hook, err := repo.GetByName(hookName)
			if err != nil {
				return fmt.Errorf("failed to get hook: %w", err)
			}

			// Display hook information
			pterm.DefaultSection.Printf("Hook: %s\n", hook.Name)

			pterm.Info.Printf("ID: %s\n", hook.ID)
			pterm.Info.Printf("Description: %s\n", hook.Description)
			pterm.Info.Printf("Event Type: %s\n", hook.EventType)
			pterm.Info.Printf("File Path: %s\n", hook.FilePath)

			if hook.Enabled {
				pterm.Success.Println("Status: Enabled")
			} else {
				pterm.Warning.Println("Status: Disabled")
			}

			pterm.Info.Printf("Run Count: %d\n", hook.RunCount)

			if hook.LastRun != nil {
				pterm.Info.Printf("Last Run: %s\n", hook.LastRun.Format("2006-01-02 15:04:05"))
			} else {
				pterm.Info.Println("Last Run: Never")
			}

			pterm.Info.Printf("Created: %s\n", hook.CreatedAt.Format("2006-01-02 15:04:05"))
			pterm.Info.Printf("Updated: %s\n", hook.UpdatedAt.Format("2006-01-02 15:04:05"))

			// Show execution history if requested
			if showHistory {
				fmt.Println()
				pterm.DefaultSection.Println("Execution History")

				history, err := repo.GetExecutionHistory(hook.ID, historyLimit)
				if err != nil {
					return fmt.Errorf("failed to get execution history: %w", err)
				}

				if len(history) == 0 {
					pterm.Info.Println("No execution history")
					return nil
				}

				tableData := [][]string{
					{"Executed At", "Success", "Duration", "Output/Error"},
				}

				for _, exec := range history {
					success := "✓"
					output := exec.Output
					if !exec.Success {
						success = "✗"
						if exec.Error != "" {
							output = exec.Error
						}
					}

					// Truncate output if too long
					if len(output) > 100 {
						output = output[:97] + "..."
					}

					tableData = append(tableData, []string{
						exec.ExecutedAt.Format("2006-01-02 15:04:05"),
						success,
						exec.Duration.String(),
						output,
					})
				}

				pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
			}

			return nil
		},
	}

	cmd.Flags().Bool("history", false, "Show execution history")
	cmd.Flags().Int("limit", 10, "Limit number of history entries")

	return cmd
}
