package agent

import (
	"context"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewGetCommand creates the agent get command
func NewGetCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <agent_name>",
		Short: "Get detailed information about an agent",
		Long:  `Retrieves detailed system information collected from a specific agent.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr, _ := cmd.Flags().GetString("master")
			outputFormat, _ := cmd.Flags().GetString("output")

			return getAgentInfo(agentName, masterAddr, outputFormat)
		},
	}

	cmd.Flags().String("master", "localhost:50051", "Master server address")
	cmd.Flags().StringP("output", "o", "text", "Output format: text or json")

	return cmd
}

func getAgentInfo(agentName, masterAddr, outputFormat string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create connection factory and get client
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	// Use refactored function with injected client
	opts := GetAgentInfoOptions{
		AgentName:    agentName,
		OutputFormat: outputFormat,
		Writer:       os.Stdout,
	}

	return getAgentInfoWithClient(ctx, client, opts)
}
