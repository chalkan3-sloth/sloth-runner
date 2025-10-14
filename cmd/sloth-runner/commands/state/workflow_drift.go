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

// NewWorkflowDriftCommand creates the workflow drift detection command
func NewWorkflowDriftCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drift <workflow-id>",
		Short: "Detect drift in workflow resources",
		Long:  `Detects drift between the expected state and actual state of workflow resources. Similar to 'terraform plan' showing what has changed.`,
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

			detections, err := sm.GetDriftDetections(workflowID)
			if err != nil {
				return fmt.Errorf("failed to get drift detections: %w", err)
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(detections, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(detections) == 0 {
				pterm.Success.Printfln("No drift detected for workflow: %s", workflowID)
				return nil
			}

			// Count drifted resources
			driftedCount := 0
			for _, d := range detections {
				if d.Drifted {
					driftedCount++
				}
			}

			if driftedCount == 0 {
				pterm.Success.Printfln("No drift detected for workflow: %s", workflowID)
				pterm.Info.Printf("Checked %d resource(s)\n", len(detections))
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Drift Detection: %s", workflowID)
			pterm.Warning.Printfln("\n%d resource(s) have drifted from expected state", driftedCount)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "RESOURCE TYPE\tRESOURCE ID\tDRIFTED\tDETECTED AT")
			fmt.Fprintln(w, "-------------\t-----------\t-------\t-----------")

			for _, detection := range detections {
				if detection.Drifted {
					driftedStr := pterm.Red("YES")
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
						detection.ResourceType,
						pterm.Cyan(detection.ResourceID),
						driftedStr,
						detection.DetectedAt.Format("2006-01-02 15:04"),
					)
				}
			}

			w.Flush()
			fmt.Println()

			// Show differences for first drifted resource
			for _, detection := range detections {
				if detection.Drifted {
					pterm.DefaultSection.Printfln("Sample Drift Details: %s", detection.ResourceID)
					fmt.Println("Expected:")
					expectedJSON, _ := json.MarshalIndent(detection.Expected, "", "  ")
					fmt.Println(pterm.Green(string(expectedJSON)))
					fmt.Println("\nActual:")
					actualJSON, _ := json.MarshalIndent(detection.Actual, "", "  ")
					fmt.Println(pterm.Yellow(string(actualJSON)))
					break
				}
			}

			pterm.Info.Println("\nUse 'sloth-runner run <workflow>' to apply changes and fix drift")

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format: table or json")

	return cmd
}
