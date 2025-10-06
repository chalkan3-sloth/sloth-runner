package scheduler

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewDisableCommand creates the scheduler disable command
// TODO: Extract logic from main.go schedulerDisableCmd (lines ~501-523)
func NewDisableCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "disable <stack-name>",
		Short: "Disable scheduled execution for a stack",
		Long:  `Disable scheduled execution for a workflow stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
