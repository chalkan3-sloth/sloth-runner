//go:build cgo
// +build cgo

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

// NewListCommand creates the stack list command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all workflow stacks",
		Long:  `List all workflow stacks with their current state and execution history.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			stackManager, err := stack.NewStackManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize stack manager: %w", err)
			}
			defer stackManager.Close()

			stacks, err := stackManager.ListStacks()
			if err != nil {
				return fmt.Errorf("failed to list stacks: %w", err)
			}

			if len(stacks) == 0 {
				pterm.Info.Println("No stacks found.")
				return nil
			}

			// Create table output
			pterm.DefaultHeader.WithFullWidth(false).
				WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
				WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
				Printf("Workflow Stacks")

			pterm.Printf("\n")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tSTATUS\tLAST RUN\tDURATION\tEXECUTIONS\tDESCRIPTION")
			fmt.Fprintln(w, "----\t------\t--------\t--------\t----------\t-----------")

			for _, s := range stacks {
				status := s.Status
				switch status {
				case "completed":
					status = pterm.Green(status)
				case "failed":
					status = pterm.Red(status)
				case "running":
					status = pterm.Yellow(status)
				default:
					status = pterm.Gray(status)
				}

				lastRun := "never"
				if s.CompletedAt != nil {
					lastRun = s.CompletedAt.Format("2006-01-02 15:04")
				} else if s.UpdatedAt.Year() > 1 {
					lastRun = s.UpdatedAt.Format("2006-01-02 15:04")
				}

				duration := "0s"
				if s.LastDuration > 0 {
					duration = s.LastDuration.String()
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n",
					s.Name, status, lastRun, duration, s.ExecutionCount, s.Description)
			}

			return w.Flush()
		},
	}
}
