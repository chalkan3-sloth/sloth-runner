//go:build cgo
// +build cgo

package state

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewShowCommand creates the state show command
// TODO: Extract logic from main.go stateShowCmd (lines ~2303-2344)
func NewShowCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "show <key>",
		Short: "Show details of a specific state",
		Long:  `Show details of a specific tracked state by its key.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
