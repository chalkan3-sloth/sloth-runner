package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewGetCommand creates the agent get command
// TODO: Extract logic from main.go agentGetCmd
func NewGetCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "get <agent-name>",
		Short: "Get detailed information about an agent",
		Long:  `Retrieves detailed system information collected from a specific agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Extract from main.go lines ~1650-1857
			// This is a complex command with JSON and human-readable output
			return nil
		},
	}
}
