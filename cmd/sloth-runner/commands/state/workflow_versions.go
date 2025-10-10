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

// NewWorkflowVersionsCommand creates the workflow versions command
func NewWorkflowVersionsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions <workflow-id>",
		Short: "List workflow state versions",
		Long:  `Lists all versions of a workflow state, allowing you to see the state history and perform rollbacks.`,
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

			versions, err := sm.GetVersions(workflowID)
			if err != nil {
				return fmt.Errorf("failed to get versions: %w", err)
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(versions, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(versions) == 0 {
				pterm.Info.Printfln("No versions found for workflow: %s", workflowID)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("State Versions: %s", workflowID)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "VERSION\tCREATED AT\tCREATED BY\tDESCRIPTION")
			fmt.Fprintln(w, "-------\t----------\t----------\t-----------")

			for _, version := range versions {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					version.Version,
					version.CreatedAt.Format("2006-01-02 15:04:05"),
					pterm.Cyan(version.CreatedBy),
					version.Description,
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d version(s)\n", len(versions))
			pterm.Info.Println("\nUse 'sloth-runner state workflow rollback <workflow-id> <version>' to rollback to a previous version")

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format: table or json")

	return cmd
}
