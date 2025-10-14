//go:build cgo
// +build cgo

package scheduler

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates the scheduler delete command
// TODO: Extract logic from main.go schedulerDeleteCmd (lines ~548-573)
func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <stack-name>",
		Short: "Delete scheduled execution for a stack",
		Long:  `Delete scheduled execution configuration for a workflow stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
