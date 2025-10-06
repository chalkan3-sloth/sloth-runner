package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewExecCommand creates the agent exec command
func NewExecCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec <agent_name> <command>",
		Short: "Executes a command on a remote agent",
		Long:  `Executes an arbitrary shell command on a specified remote agent.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			command := args[1]
			masterAddr, _ := cmd.Flags().GetString("master")
			outputFormat, _ := cmd.Flags().GetString("output")

			// Get master address from environment if not specified
			if masterAddr == "" {
				masterAddr = os.Getenv("SLOTH_RUNNER_MASTER_ADDR")
			}

			if masterAddr == "" {
				return fmt.Errorf("master address not specified. Use --master flag or set SLOTH_RUNNER_MASTER_ADDR environment variable")
			}

			return runCommandOnAgent(agentName, command, masterAddr, outputFormat)
		},
	}

	cmd.Flags().String("master", "", "Master server address (or use SLOTH_RUNNER_MASTER_ADDR env var)")
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
