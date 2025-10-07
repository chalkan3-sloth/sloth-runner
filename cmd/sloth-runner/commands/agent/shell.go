package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
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

	// Connect to agent directly
	conn, err := grpc.Dial(agentInfo.AgentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	agentClient := pb.NewAgentClient(conn)

	// Check if we have a TTY
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("interactive shell requires a TTY - please run from a real terminal, not via pipes or automation")
	}

	// Start interactive shell stream
	stream, err := agentClient.InteractiveShell(ctx)
	if err != nil {
		return fmt.Errorf("failed to open shell: %w", err)
	}

	// Display welcome banner
	printShellBanner(agentName, agentInfo.AgentAddress)

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
				// Check if it's a normal shell exit (PTY closed)
				if strings.Contains(err.Error(), "PTY read error") ||
				   strings.Contains(err.Error(), "input/output error") {
					// Normal shell exit, not an error
					close(done)
					return
				}
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
	printShellGoodbye()
	return nil
}

// printShellBanner displays a welcome banner when connecting to the shell
func printShellBanner(agentName, address string) {
	banner := pterm.DefaultBox.WithTitle("Sloth Runner Interactive Shell").WithTitleTopCenter().Sprint(
		fmt.Sprintf("Connected to: %s\nAddress: %s\n\nCommands:\n  • Type commands normally\n  • Press Ctrl+D or type 'exit' to quit\n  • Press Ctrl+C to interrupt current command",
			pterm.FgGreen.Sprint(agentName),
			pterm.FgCyan.Sprint(address),
		),
	)
	fmt.Println(banner)
	fmt.Println()
}

// printShellGoodbye displays a goodbye message when exiting the shell
func printShellGoodbye() {
	fmt.Println()
	pterm.Success.Println("Shell session closed. Goodbye!")
}
