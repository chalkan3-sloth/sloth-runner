package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewStartCommand creates the agent start command
// TODO: Extract logic from main.go agentStartCmd
func NewStartCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Starts the sloth-runner in agent mode",
		Long:  `The agent start command starts the sloth-runner as a background agent that can execute tasks remotely.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Extract from main.go lines ~1152-1373
			return nil
		},
	}
}
