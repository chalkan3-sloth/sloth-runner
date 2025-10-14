package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// runInteractiveShell handles the interactive shell session with proper terminal management
func runInteractiveShell(ctx context.Context, stream pb.Agent_InteractiveShellClient, agentName, agentAddress string) error {
	// Check if we have a TTY
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("interactive shell requires a TTY - please run from a real terminal, not via pipes or automation")
	}

	// Get initial terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width, height = 80, 24 // fallback to standard size
	}

	// Send initial window size to server
	if err := stream.Send(&pb.ShellInput{
		WindowRows: uint32(height),
		WindowCols: uint32(width),
	}); err != nil {
		return fmt.Errorf("failed to send initial window size: %w", err)
	}

	// Display welcome banner
	printImprovedShellBanner(agentName, agentAddress, width, height)

	// Save current terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}

	// Ensure terminal is restored on exit
	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		// Reset terminal: clear formatting, show cursor
		fmt.Print("\033[0m\033[?25h\r\n")
	}()

	// Handle signals: Ctrl+C, SIGTERM, and window resize
	sigChan := make(chan os.Signal, 1)
	setupShellSignals(sigChan)

	// Channels for coordination
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

	// Goroutine to read from stdin and send to shell
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

	// Main loop: handle signals and wait for completion
	for {
		select {
		case sig := <-sigChan:
			if isSigwinch(sig) {
				// Handle terminal resize
				newWidth, newHeight, _ := term.GetSize(int(os.Stdin.Fd()))
				if err := stream.Send(&pb.ShellInput{
					Resize:     true,
					WindowRows: uint32(newHeight),
					WindowCols: uint32(newWidth),
				}); err != nil {
					return fmt.Errorf("failed to send resize: %w", err)
				}
			} else if sig == os.Interrupt {
				// Send Ctrl+C to remote shell
				stream.Send(&pb.ShellInput{StdinData: []byte{0x03}}) // Ctrl+C
			} else if sig == syscall.SIGTERM {
				// Terminate the shell gracefully
				stream.Send(&pb.ShellInput{Terminate: true})
				stream.CloseSend()
				return nil
			}
		case err := <-errChan:
			stream.CloseSend()
			return err
		case <-done:
			stream.CloseSend()
			printShellGoodbye()
			return nil
		}
	}
}

// printImprovedShellBanner displays an enhanced welcome banner
func printImprovedShellBanner(agentName, address string, width, height int) {
	hostname, _ := os.Hostname()

	banner := pterm.DefaultBox.
		WithTitle("🚀 Sloth Runner Interactive Shell").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).
		Sprint(fmt.Sprintf(`
Connected to: %s
Address:      %s
Terminal:     %dx%d
From:         %s

Commands:
  • Type commands normally
  • Ctrl+D or 'exit' to quit
  • Ctrl+C to interrupt command
  • Ctrl+L to clear screen

Tip: Full-screen apps (vi, htop, nano) are fully supported!
`,
		pterm.FgGreen.Sprint(agentName),
		pterm.FgCyan.Sprint(address),
		width, height,
		pterm.FgYellow.Sprint(hostname),
	))

	fmt.Println(banner)
	fmt.Println()
}
