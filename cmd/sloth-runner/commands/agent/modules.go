package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewModulesCommand creates the agent modules command
// TODO: Extract logic from main.go agentModulesCheckCmd
func NewModulesCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "modules",
		Short: "Check available modules on an agent",
		Long:  `Checks which modules are available on the specified agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Extract from main.go lines ~1859-2000
			return nil
		},
	}
}
