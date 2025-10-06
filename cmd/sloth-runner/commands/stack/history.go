package stack

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewHistoryCommand creates the stack history command
func NewHistoryCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history <stack-name>",
		Short: "Show execution history of a stack",
		Long:  `Show detailed execution history of a specific workflow stack.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			limit, _ := cmd.Flags().GetInt("limit")

			stackManager, err := stack.NewStackManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize stack manager: %w", err)
			}
			defer stackManager.Close()

			stackState, err := stackManager.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			executions, err := stackManager.GetStackExecutions(stackState.ID, limit)
			if err != nil {
				return fmt.Errorf("failed to get executions: %w", err)
			}

			if len(executions) == 0 {
				pterm.Info.Printf("No execution history for stack '%s'.\n", stackName)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth(false).
				WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
				WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
				Printf("Execution History: %s", stackName)

			pterm.Printf("\n")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "STARTED\tSTATUS\tDURATION\tTASKS\tSUCCESS\tFAILED\tERROR")
			fmt.Fprintln(w, "-------\t------\t--------\t-----\t-------\t------\t-----")

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

				errorMsg := ""
				if exec.ErrorMessage != "" {
					if len(exec.ErrorMessage) > 30 {
						errorMsg = exec.ErrorMessage[:30] + "..."
					} else {
						errorMsg = exec.ErrorMessage
					}
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\t%s\n",
					exec.StartedAt.Format("2006-01-02 15:04:05"),
					status,
					exec.Duration.String(),
					exec.TaskCount,
					exec.SuccessCount,
					exec.FailureCount,
					errorMsg)
			}
			return w.Flush()
		},
	}

	cmd.Flags().IntP("limit", "n", 10, "Number of executions to show")
	return cmd
}
