package events

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "show <event-id>",
		Short: "Show detailed event information including hook executions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := hooks.NewRepository()
			if err != nil {
				return err
			}
			defer repo.Close()

			event, err := repo.EventQueue.GetEvent(args[0])
			if err != nil {
				return err
			}

			pterm.DefaultSection.Printf("Event: %s\n", event.ID)
			pterm.Info.Printf("Type: %s\n", event.Type)
			pterm.Info.Printf("Status: %s\n", event.Status)
			pterm.Info.Printf("Created: %s\n", event.CreatedAt.Format("2006-01-02 15:04:05"))

			if event.ProcessedAt != nil {
				pterm.Info.Printf("Processed: %s\n", event.ProcessedAt.Format("2006-01-02 15:04:05"))
			}

			if event.Error != "" {
				pterm.Error.Printf("Error: %s\n", event.Error)
			}

			// Get hook executions for this event
			executions, err := repo.EventQueue.GetEventHookExecutions(args[0])
			if err != nil {
				return fmt.Errorf("failed to get hook executions: %w", err)
			}

			if len(executions) > 0 {
				pterm.Println()
				pterm.DefaultSection.Println("Hook Executions")

				tableData := [][]string{
					{"Hook Name", "Status", "Duration", "Executed At"},
				}

				for _, exec := range executions {
					status := "✅ Success"
					if !exec.Success {
						status = "❌ Failed"
					}

					tableData = append(tableData, []string{
						exec.HookName,
						status,
						exec.Duration.String(),
						exec.ExecutedAt.Format("2006-01-02 15:04:05"),
					})
				}

				pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

				// Show detailed output/errors for each execution
				for _, exec := range executions {
					if exec.Output != "" || exec.Error != "" {
						pterm.Println()
						pterm.DefaultSection.Printf("Hook: %s\n", exec.HookName)

						if exec.Output != "" {
							pterm.Info.Println("Output:")
							pterm.Println(exec.Output)
						}

						if exec.Error != "" {
							pterm.Error.Println("Error:")
							pterm.Println(exec.Error)
						}
					}
				}
			} else {
				pterm.Println()
				pterm.Info.Println("No hooks were executed for this event")
			}

			return nil
		},
	}
}
