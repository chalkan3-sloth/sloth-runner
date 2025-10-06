package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewRunCommand creates the agent run command
// TODO: Extract logic from main.go agentRunCmd
func NewRunCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Runs a command on an agent",
		Long:  `Executes a command on a remote agent and streams the output back.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Extract from main.go lines ~1375-1502
			return nil
		},
	}
}
