package agent

import (
	"context"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewUpdateCommand creates the agent update command
func NewUpdateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <agent_name>",
		Short: "Update an agent to the latest version",
		Long:  `Updates the specified agent to the latest available version from GitHub releases.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr, _ := cmd.Flags().GetString("master")
			version, _ := cmd.Flags().GetString("version")
			restart, _ := cmd.Flags().GetBool("restart")

			return updateAgent(agentName, masterAddr, version, restart)
		},
	}

	cmd.Flags().String("master", "localhost:50051", "Master server address")
	cmd.Flags().String("version", "latest", "Version to update to (default: latest)")
	cmd.Flags().Bool("restart", true, "Restart agent service after update")

	return cmd
}

func updateAgent(agentName, masterAddr, version string, restart bool) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create connection factory
	factory := NewDefaultConnectionFactory()

	// Get registry client
	registryClient, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	// Create agent client factory function
	agentClientFactory := func(addr string) (AgentClient, func(), error) {
		return factory.CreateAgentClient(addr)
	}

	// Use refactored function with injected clients
	opts := UpdateAgentOptions{
		AgentName:     agentName,
		TargetVersion: version,
		Restart:       restart,
		Writer:        os.Stdout,
	}

	_, err = updateAgentWithClients(ctx, registryClient, agentClientFactory, opts)
	return err
}
