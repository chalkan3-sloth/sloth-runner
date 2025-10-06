package state

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates the state delete command
// TODO: Extract logic from main.go stateDeleteCmd (lines ~2345-2382)
func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <key>",
		Short: "Delete a specific state",
		Long:  `Delete a specific tracked state by its key.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
