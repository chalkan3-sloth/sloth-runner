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

// NewWorkflowResourcesCommand creates the workflow resources command
func NewWorkflowResourcesCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resources <workflow-id>",
		Short: "List workflow resources",
		Long:  `Lists all resources managed by a workflow execution, including their type, action, and status.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			outputFormat, _ := cmd.Flags().GetString("output")
			resourceType, _ := cmd.Flags().GetString("type")

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

			resources := workflow.Resources

			// Filter by type if specified
			if resourceType != "" {
				filtered := []state.Resource{}
				for _, r := range resources {
					if r.Type == resourceType {
						filtered = append(filtered, r)
					}
				}
				resources = filtered
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(resources, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(resources) == 0 {
				pterm.Info.Println("No resources found")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Workflow Resources: %s", workflow.Name)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tTYPE\tNAME\tACTION\tSTATUS\tCREATED")
			fmt.Fprintln(w, "--\t----\t----\t------\t------\t-------")

			for _, resource := range resources {
				actionColor := pterm.Green
				switch resource.Action {
				case state.ResourceActionCreate:
					actionColor = pterm.Green
				case state.ResourceActionUpdate:
					actionColor = pterm.Yellow
				case state.ResourceActionDelete:
					actionColor = pterm.Red
				case state.ResourceActionNoop:
					actionColor = pterm.Gray
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					resource.ID[:8]+"...",
					resource.Type,
					pterm.Cyan(resource.Name),
					actionColor(string(resource.Action)),
					resource.Status,
					resource.CreatedAt.Format("2006-01-02 15:04"),
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d resource(s)\n", len(resources))

			// Summary by action
			actionCounts := make(map[state.ResourceAction]int)
			for _, r := range resources {
				actionCounts[r.Action]++
			}

			if len(actionCounts) > 0 {
				fmt.Println("\nActions:")
				for action, count := range actionCounts {
					fmt.Printf("  %s: %d\n", action, count)
				}
			}

			return nil
		},
	}

	cmd.Flags().String("type", "", "Filter by resource type")
	cmd.Flags().StringP("output", "o", "table", "Output format: table or json")

	return cmd
}
