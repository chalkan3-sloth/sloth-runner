package state

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewClearCommand creates the state clear command
// TODO: Extract logic from main.go stateClearCmd (lines ~2383-2418)
func NewClearCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clear all tracked states",
		Long:  `Clear all tracked states from the state database.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
