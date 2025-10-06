package agent

import (
	"context"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewModulesCommand creates the agent modules command
func NewModulesCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modules <agent_name>",
		Short: "Check availability of external modules/tools on an agent",
		Long:  `Checks which external tools and modules are available on the specified agent for Lua tasks to use.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr, _ := cmd.Flags().GetString("master")

			return checkAgentModules(agentName, masterAddr)
		},
	}

	cmd.Flags().String("master", "localhost:50051", "Master server address")

	return cmd
}

func checkAgentModules(agentName, masterAddr string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create connection factory and get client
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	// Use refactored function with injected client
	opts := ModulesCheckOptions{
		AgentName: agentName,
		Writer:    os.Stdout,
	}

	return checkAgentModulesWithClient(ctx, client, opts)
}
