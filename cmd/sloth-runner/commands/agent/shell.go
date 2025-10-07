package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewShellCommand creates the agent shell command
func NewShellCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell <agent_name>",
		Short: "Open an interactive shell on a remote agent",
		Long:  `Opens an interactive bash shell on the specified agent via gRPC streaming.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr := getMasterAddress(cmd)

			return openAgentShell(agentName, masterAddr)
		},
	}

	addMasterFlag(cmd)

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

	pterm.Info.Printf("🔗 Connecting to agent %s at %s...\n", agentName, agentInfo.AgentAddress)

	// Connect to agent directly
	conn, err := grpc.Dial(agentInfo.AgentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	pterm.Success.Printf("✅ Connected! Type 'exit' to quit\n\n")

	// Save current terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel to signal when to stop
	done := make(chan struct{})
	errChan := make(chan error, 2)

	// Goroutine to read from shell and print to stdout
	go func() {
		for {
			output, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("stream receive error: %w", err)
				return
			}

			if len(output.Stdout) > 0 {
				os.Stdout.Write(output.Stdout)
			}

			if len(output.Stderr) > 0 {
				os.Stderr.Write(output.Stderr)
			}

			if output.Completed {
				if output.Error != "" {
					errChan <- fmt.Errorf("shell error: %s", output.Error)
				}
				close(done)
				return
			}
		}
	}()

	// Goroutine to read from stdin and send to shell (char by char)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err != io.EOF {
					errChan <- fmt.Errorf("stdin read error: %w", err)
				}
				return
			}

			if n > 0 {
				// Send raw input to shell
				if err := stream.Send(&pb.ShellInput{StdinData: buf[:n]}); err != nil {
					errChan <- fmt.Errorf("failed to send input: %w", err)
					return
				}
			}
		}
	}()

	// Wait for completion, error, or signal
	select {
	case <-sigChan:
		// Send Ctrl+C to remote shell
		stream.Send(&pb.ShellInput{StdinData: []byte{0x03}}) // Ctrl+C
		stream.CloseSend()
	case err := <-errChan:
		stream.CloseSend()
		term.Restore(int(os.Stdin.Fd()), oldState)
		return err
	case <-done:
		stream.CloseSend()
	}

	term.Restore(int(os.Stdin.Fd()), oldState)
	pterm.Info.Println("\n👋 Shell session closed")
	return nil
}
