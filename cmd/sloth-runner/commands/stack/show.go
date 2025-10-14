//go:build cgo
// +build cgo

package stack

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewShowCommand creates the stack show command
func NewShowCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "show <stack-name>",
		Short: "Show detailed information about a stack",
		Long:  `Show detailed information about a specific workflow stack including execution history.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			stackManager, err := stack.NewStackManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize stack manager: %w", err)
			}
			defer stackManager.Close()

			stackState, err := stackManager.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			// Show stack details
			pterm.DefaultHeader.WithFullWidth(false).
				WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
				WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
				Printf("Stack: %s", stackState.Name)

			pterm.Printf("\n")
			pterm.Printf("ID: %s\n", stackState.ID)
			pterm.Printf("Description: %s\n", stackState.Description)
			pterm.Printf("Version: %s\n", stackState.Version)
			pterm.Printf("Status: %s\n", stackState.Status)
			pterm.Printf("Created: %s\n", stackState.CreatedAt.Format("2006-01-02 15:04:05"))
			pterm.Printf("Updated: %s\n", stackState.UpdatedAt.Format("2006-01-02 15:04:05"))
			if stackState.CompletedAt != nil {
				pterm.Printf("Completed: %s\n", stackState.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			pterm.Printf("Workflow File: %s\n", stackState.WorkflowFile)
			pterm.Printf("Executions: %d\n", stackState.ExecutionCount)
			if stackState.LastDuration > 0 {
				pterm.Printf("Last Duration: %s\n", stackState.LastDuration.String())
			}
			if stackState.LastError != "" {
				pterm.Printf("Last Error: %s\n", pterm.Red(stackState.LastError))
			}

			// Show outputs if any
			if len(stackState.Outputs) > 0 {
				pterm.Printf("\n")
				pterm.DefaultHeader.WithFullWidth(false).
					WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
					WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
					Printf("Outputs")
				pterm.Printf("\n")
				for key, value := range stackState.Outputs {
					pterm.Printf("%s: %v\n", pterm.Cyan(key), value)
				}
			}

			// Show recent executions
			executions, err := stackManager.GetStackExecutions(stackState.ID, 5)
			if err != nil {
				slog.Warn("Failed to get executions", "error", err)
			} else if len(executions) > 0 {
				pterm.Printf("\n")
				pterm.DefaultHeader.WithFullWidth(false).
					WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
					WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
					Printf("Recent Executions")
				pterm.Printf("\n")

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(w, "STARTED\tSTATUS\tDURATION\tTASKS\tSUCCESS\tFAILED")
				fmt.Fprintln(w, "-------\t------\t--------\t-----\t-------\t------")

				for _, exec := range executions {
					status := exec.Status
					switch status {
					case "completed":
						status = pterm.Green(status)
					case "failed":
						status = pterm.Red(status)
					default:
						status = pterm.Gray(status)
					}

					fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\n",
						exec.StartedAt.Format("2006-01-02 15:04"),
						status,
						exec.Duration.String(),
						exec.TaskCount,
						exec.SuccessCount,
						exec.FailureCount)
				}
				w.Flush()
			}

			return nil
		},
	}
}
