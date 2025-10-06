package state

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewStatsCommand creates the state stats command
// TODO: Extract logic from main.go stateStatsCmd (lines ~2419-2460)
func NewStatsCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show state statistics",
		Long:  `Show statistics about tracked states.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			return nil
		},
	}
}
