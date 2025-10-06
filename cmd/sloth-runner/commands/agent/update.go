package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewUpdateCommand creates the agent update command
// NOTE: This command already exists in cmd/sloth-runner/agent_update.go
// TODO: Integrate existing agentUpdateCmd into this modular structure
func NewUpdateCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "update <agent-name>",
		Short: "Update the sloth-runner agent to the latest version",
		Long:  `Updates the sloth-runner agent binary to the latest version from GitHub releases via gRPC.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Move logic from agent_update.go to this modular structure
			return nil
		},
	}
}
