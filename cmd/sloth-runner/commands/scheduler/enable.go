package scheduler

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewEnableCommand creates the scheduler enable command
// TODO: Extract logic from main.go schedulerEnableCmd (lines ~476-500)
func NewEnableCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "enable <stack-name> <cron-expression>",
		Short: "Enable scheduled execution for a stack",
		Long:  `Enable scheduled execution for a workflow stack using cron expression.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
