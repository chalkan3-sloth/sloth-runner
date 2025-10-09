package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/spf13/cobra"
)

// NewExecCommand creates the agent exec command
func NewExecCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec <agent_name> <command>",
		Short: "Executes a command on a remote agent",
		Long:  `Executes an arbitrary shell command on a specified remote agent. With --local, connects directly to agent using local database.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			command := args[1]
			local, _ := cmd.Flags().GetBool("local")
			outputFormat, _ := cmd.Flags().GetString("output")

			// If --local flag is set, connect directly to agent
			if local {
				return runCommandOnAgentDirect(agentName, command, outputFormat)
			}

			// Otherwise use master server
			masterAddr, _ := cmd.Flags().GetString("master")

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
	cmd.Flags().Bool("local", false, "Connect directly to agent using local database")

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

// runCommandOnAgentDirect connects directly to agent using local database
func runCommandOnAgentDirect(agentName, command, outputFormat string) error {
	// Get agent address from local database
	agentAddr, err := getAgentAddressFromLocalDB(agentName)
	if err != nil {
		return err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create direct gRPC connection to agent
	conn, err := createGRPCConnection(agentAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to agent at %s: %w", agentAddr, err)
	}
	defer conn.Close()

	// Create agent client
	client := pb.NewAgentClient(conn)

	// Use refactored function with agent client
	opts := RunCommandOptions{
		AgentName:    agentName,
		Command:      command,
		OutputFormat: outputFormat,
		OutputWriter: os.Stdout,
		ErrorWriter:  os.Stderr,
	}

	return runCommandDirectly(ctx, client, opts)
}

