//go:build cgo
// +build cgo

package scheduler

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewListCommand creates the scheduler list command
// TODO: Extract logic from main.go schedulerListCmd (lines ~524-547)
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all scheduled workflows",
		Long:  `List all workflows that have scheduled execution enabled.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
