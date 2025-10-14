//go:build cgo
// +build cgo

package state

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewWorkflowOutputsCommand creates the workflow outputs command
func NewWorkflowOutputsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outputs <workflow-id>",
		Short: "Show workflow outputs",
		Long:  `Displays all output values from a workflow execution. Similar to 'terraform output'.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			outputFormat, _ := cmd.Flags().GetString("output")

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			if err := sm.InitWorkflowSchema(); err != nil {
				return fmt.Errorf("failed to initialize workflow schema: %w", err)
			}

			workflow, err := sm.GetWorkflowState(workflowID)
			if err != nil {
				return fmt.Errorf("failed to get workflow state: %w", err)
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(workflow.Outputs, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(workflow.Outputs) == 0 {
				pterm.Info.Printfln("No outputs defined for workflow: %s", workflow.Name)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Workflow Outputs: %s", workflow.Name)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "KEY\tVALUE")
			fmt.Fprintln(w, "---\t-----")

			for key, value := range workflow.Outputs {
				fmt.Fprintf(w, "%s\t%s\n",
					pterm.Cyan(key),
					value,
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d output(s)\n", len(workflow.Outputs))

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format: table or json")

	return cmd
}
