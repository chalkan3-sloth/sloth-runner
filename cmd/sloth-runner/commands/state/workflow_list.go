//go:build cgo
// +build cgo

package state

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewWorkflowListCommand creates the workflow list command
func NewWorkflowListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workflow execution states",
		Long:  `Lists all workflow execution states with their current status, version, and duration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			status, _ := cmd.Flags().GetString("status")
			outputFormat, _ := cmd.Flags().GetString("output")

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			// Initialize workflow schema
			if err := sm.InitWorkflowSchema(); err != nil {
				return fmt.Errorf("failed to initialize workflow schema: %w", err)
			}

			workflows, err := sm.ListWorkflowStates(name, state.WorkflowStateStatus(status))
			if err != nil {
				return fmt.Errorf("failed to list workflow states: %w", err)
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(workflows, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(workflows) == 0 {
				pterm.Info.Println("No workflow states found")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Println("Workflow States")
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tVERSION\tSTATUS\tSTARTED\tDURATION\tRESOURCES")
			fmt.Fprintln(w, "----\t-------\t------\t-------\t--------\t---------")

			for _, workflow := range workflows {
				statusColor := pterm.Green
				switch workflow.Status {
				case state.WorkflowStatusFailed:
					statusColor = pterm.Red
				case state.WorkflowStatusRunning:
					statusColor = pterm.Yellow
				case state.WorkflowStatusRolledBack:
					statusColor = pterm.Magenta
				}

				duration := formatDuration(time.Duration(workflow.Duration) * time.Second)
				started := workflow.StartedAt.Format("2006-01-02 15:04")

				fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%d\n",
					pterm.Cyan(workflow.Name),
					workflow.Version,
					statusColor(string(workflow.Status)),
					started,
					duration,
					len(workflow.Resources),
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d workflow(s)\n", len(workflows))

			return nil
		},
	}

	cmd.Flags().String("name", "", "Filter by workflow name")
	cmd.Flags().String("status", "", "Filter by status (pending, running, success, failed, rolled_back)")
	cmd.Flags().StringP("output", "o", "table", "Output format: table or json")

	return cmd
}
