package agent

import (
	"context"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewStopCommand creates the agent stop command
func NewStopCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop <agent-name>",
		Short: "Stops a running agent",
		Long:  `Stops a running agent by sending a shutdown request.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr, _ := cmd.Flags().GetString("master")

			// Create connection factory and get client
			factory := NewDefaultConnectionFactory()
			client, cleanup, err := factory.CreateRegistryClient(masterAddr)
			if err != nil {
				return err
			}
			defer cleanup()

			// Use refactored function with injected client
			opts := StopAgentOptions{
				AgentName: agentName,
			}

			err = stopAgentWithClient(context.Background(), client, opts)

			// Track operation
			trackAgentStop(agentName, err == nil)

			return err
		},
	}

	cmd.Flags().String("master", "localhost:50051", "Master server address")

	return cmd
}

func init() {
	// This will be called from NewAgentCommand
}
