package agent

import (
	"context"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewRunCommand creates the agent run command
func NewRunCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <agent_name> <command>",
		Short: "Executes a command on a remote agent",
		Long:  `Executes an arbitrary shell command on a specified remote agent.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			command := args[1]
			masterAddr, _ := cmd.Flags().GetString("master")
			outputFormat, _ := cmd.Flags().GetString("output")

			return runCommandOnAgent(agentName, command, masterAddr, outputFormat)
		},
	}

	cmd.Flags().String("master", "localhost:50051", "Master server address")
	cmd.Flags().StringP("output", "o", "text", "Output format: text or json")

	return cmd
}

func runCommandOnAgent(agentName, command, masterAddr, outputFormat string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create connection factory and get client
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	// Use refactored function with injected client
	opts := RunCommandOptions{
		AgentName:    agentName,
		Command:      command,
		OutputFormat: outputFormat,
		OutputWriter: os.Stdout,
		ErrorWriter:  os.Stderr,
	}

	return runCommandWithClient(ctx, client, opts)
}
