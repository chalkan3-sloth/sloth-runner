package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewListCommand creates the state list command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [resource-type]",
		Short: "List all tracked states or filter by resource type",
		Long:  `Lists all resources being tracked for idempotency. Optionally filter by resource type.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agent, _ := cmd.Flags().GetString("agent")
			outputFormat, _ := cmd.Flags().GetString("output")

			var resourceType string
			if len(args) > 0 {
				resourceType = args[0]
			}

			sm, err := state.NewStateManager(filepath.Join(os.TempDir(), "sloth-state", agent+".db"))
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			states, err := sm.List(resourceType)
			if err != nil {
				return fmt.Errorf("failed to list states: %w", err)
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(states, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(states) == 0 {
				if resourceType != "" {
					pterm.Info.Printf("No states found for resource type: %s\n", resourceType)
				} else {
					pterm.Info.Println("No states tracked yet")
				}
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Println("State Tracking")
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "KEY\tVALUE\tUPDATED")
			fmt.Fprintln(w, "---\t-----\t-------")

			for key, value := range states {
				truncValue := value
				if len(truncValue) > 50 {
					truncValue = truncValue[:47] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n",
					pterm.Cyan(key),
					truncValue,
					"")
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d state(s)\n", len(states))

			return nil
		},
	}

	cmd.Flags().String("agent", "local", "Agent name for state storage")
	cmd.Flags().StringP("output", "o", "table", "Output format: table or json")

	return cmd
}
