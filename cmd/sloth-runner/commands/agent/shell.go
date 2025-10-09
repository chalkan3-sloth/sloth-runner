package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// NewShellCommand creates the agent shell command
func NewShellCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell <agent_name>",
		Short: "Open an interactive shell on a remote agent",
		Long:  `Opens an interactive bash shell on the specified agent via gRPC streaming. With --local, connects directly using local database.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			local, _ := cmd.Flags().GetBool("local")

			// If --local flag is set, connect directly
			if local {
				return openAgentShellDirect(agentName)
			}

			masterAddr := getMasterAddress(cmd)
			return openAgentShell(agentName, masterAddr)
		},
	}

	addMasterFlag(cmd)
	cmd.Flags().Bool("local", false, "Connect directly to agent using local database")

	return cmd
}

func openAgentShell(agentName, masterAddr string) error {
	ctx := context.Background()

	// Connect to master to get agent address
	factory := NewDefaultConnectionFactory()
	registryClient, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to master: %w", err)
	}
	defer cleanup()

	// Get agent info
	agentResp, err := registryClient.GetAgentInfo(ctx, &pb.GetAgentInfoRequest{AgentName: agentName})
	if err != nil || !agentResp.Success {
		return fmt.Errorf("agent not found: %s", agentName)
	}

	agentInfo := agentResp.AgentInfo

	// Connect to agent directly with keep-alive
	conn, err := grpc.Dial(agentInfo.AgentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}))
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	agentClient := pb.NewAgentClient(conn)

	// Start interactive shell stream
	stream, err := agentClient.InteractiveShell(ctx)
	if err != nil {
		return fmt.Errorf("failed to open shell: %w", err)
	}

	// Use the robust interactive shell handler
	return runInteractiveShellRobust(ctx, stream, agentName, agentInfo.AgentAddress)
}

// printShellGoodbye displays a goodbye message when exiting the shell
func printShellGoodbye() {
	fmt.Println()
	pterm.Success.Println("✨ Shell session closed. Goodbye!")
}

// openAgentShellDirect opens shell directly to agent using local database
func openAgentShellDirect(agentName string) error {
	ctx := context.Background()

	// Get agent address from local database
	agentAddr, err := getAgentAddressFromLocalDB(agentName)
	if err != nil {
		return err
	}

	// Connect to agent directly with keep-alive
	conn, err := grpc.Dial(agentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}))
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	agentClient := pb.NewAgentClient(conn)

	// Start interactive shell stream
	stream, err := agentClient.InteractiveShell(ctx)
	if err != nil {
		return fmt.Errorf("failed to open shell: %w", err)
	}

	// Use the robust interactive shell handler
	return runInteractiveShellRobust(ctx, stream, agentName, agentAddr)
}
