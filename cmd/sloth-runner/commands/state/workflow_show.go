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

// NewWorkflowShowCommand creates the workflow show command
func NewWorkflowShowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <workflow-name-or-id>",
		Short: "Show detailed workflow state",
		Long:  `Displays detailed information about a workflow execution state, including resources and outputs.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowNameOrID := args[0]
			outputFormat, _ := cmd.Flags().GetString("output")

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			if err := sm.InitWorkflowSchema(); err != nil {
				return fmt.Errorf("failed to initialize workflow schema: %w", err)
			}

			// Try to get by ID first, then by name
			workflow, err := sm.GetWorkflowState(workflowNameOrID)
			if err != nil {
				workflow, err = sm.GetWorkflowStateByName(workflowNameOrID)
				if err != nil {
					return fmt.Errorf("workflow not found: %s", workflowNameOrID)
				}
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(workflow, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			// Display workflow header
			pterm.DefaultHeader.WithFullWidth().Printfln("Workflow: %s (v%d)", workflow.Name, workflow.Version)
			fmt.Println()

			// Basic info
			pterm.DefaultSection.Println("Basic Information")
			fmt.Printf("ID:           %s\n", pterm.Cyan(workflow.ID))
			fmt.Printf("Name:         %s\n", pterm.Cyan(workflow.Name))
			fmt.Printf("Version:      %s\n", pterm.Cyan(fmt.Sprintf("%d", workflow.Version)))

			statusColor := pterm.Green
			switch workflow.Status {
			case state.WorkflowStatusFailed:
				statusColor = pterm.Red
			case state.WorkflowStatusRunning:
				statusColor = pterm.Yellow
			case state.WorkflowStatusRolledBack:
				statusColor = pterm.Magenta
			}
			fmt.Printf("Status:       %s\n", statusColor(string(workflow.Status)))

			fmt.Printf("Started At:   %s\n", workflow.StartedAt.Format("2006-01-02 15:04:05"))
			if workflow.CompletedAt != nil {
				fmt.Printf("Completed At: %s\n", workflow.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("Duration:     %s\n", formatDuration(workflow.Duration))

			if workflow.ErrorMsg != "" {
				fmt.Printf("Error:        %s\n", pterm.Red(workflow.ErrorMsg))
			}

			if workflow.LockedBy != "" {
				fmt.Printf("Locked By:    %s\n", pterm.Yellow(workflow.LockedBy))
			}
			fmt.Println()

			// Metadata
			if len(workflow.Metadata) > 0 {
				pterm.DefaultSection.Println("Metadata")
				for key, value := range workflow.Metadata {
					fmt.Printf("%s: %s\n", pterm.Cyan(key), value)
				}
				fmt.Println()
			}

			// Resources
			if len(workflow.Resources) > 0 {
				pterm.DefaultSection.Println("Resources")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(w, "TYPE\tNAME\tACTION\tSTATUS")
				fmt.Fprintln(w, "----\t----\t------\t------")

				for _, resource := range workflow.Resources {
					actionColor := pterm.Green
					switch resource.Action {
					case state.ResourceActionCreate:
						actionColor = pterm.Green
					case state.ResourceActionUpdate:
						actionColor = pterm.Yellow
					case state.ResourceActionDelete:
						actionColor = pterm.Red
					}

					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
						resource.Type,
						pterm.Cyan(resource.Name),
						actionColor(string(resource.Action)),
						resource.Status,
					)
				}
				w.Flush()
				fmt.Println()
			}

			// Outputs
			if len(workflow.Outputs) > 0 {
				pterm.DefaultSection.Println("Outputs")
				for key, value := range workflow.Outputs {
					fmt.Printf("%s = %s\n", pterm.Cyan(key), value)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format: table or json")

	return cmd
}
