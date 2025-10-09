package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// runInteractiveShellRobust is a rock-solid implementation of the interactive shell
func runInteractiveShellRobust(ctx context.Context, stream pb.Agent_InteractiveShellClient, agentName, agentAddress string) error {
	// Check TTY
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("interactive shell requires a TTY")
	}

	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width, height = 80, 24
	}

	// Send initial window size
	if err := stream.Send(&pb.ShellInput{
		WindowRows: uint32(height),
		WindowCols: uint32(width),
	}); err != nil {
		return fmt.Errorf("failed to send initial window size: %w", err)
	}

	// Print banner
	printRobustBanner(agentName, agentAddress, width, height)

	// Set raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}

	// Critical: Always restore terminal
	var terminalRestored bool
	var restoreMu sync.Mutex
	restoreTerminal := func() {
		restoreMu.Lock()
		defer restoreMu.Unlock()
		if !terminalRestored {
			term.Restore(int(os.Stdin.Fd()), oldState)
			fmt.Print("\033[0m\033[?25h\r\n") // Reset + show cursor + CRLF
			terminalRestored = true
		}
	}
	defer restoreTerminal()

	// Signal handling
	sigChan := make(chan os.Signal, 10) // Buffered to avoid missing signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGWINCH)
	defer signal.Stop(sigChan)

	// Coordination channels
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stdinDone := make(chan struct{})
	stdoutDone := make(chan struct{})
	errorChan := make(chan error, 3)

	// Goroutine 1: Read from remote shell, write to local stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(stdoutDone)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			output, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				// Only report unexpected errors
				if ctx.Err() == nil {
					select {
					case errorChan <- fmt.Errorf("receive error: %w", err):
					default:
					}
				}
				return
			}

			// Write stdout data
			if len(output.Stdout) > 0 {
				written := 0
				for written < len(output.Stdout) {
					n, err := os.Stdout.Write(output.Stdout[written:])
					if err != nil {
						return
					}
					written += n
				}
			}

			// Write stderr data
			if len(output.Stderr) > 0 {
				written := 0
				for written < len(output.Stderr) {
					n, err := os.Stderr.Write(output.Stderr[written:])
					if err != nil {
						return
					}
					written += n
				}
			}

			// Check for completion
			if output.Completed {
				if output.Error != "" {
					select {
					case errorChan <- fmt.Errorf("remote error: %s", output.Error):
					default:
					}
				}
				return
			}
		}
	}()

	// Goroutine 2: Read from local stdin, send to remote shell
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(stdinDone)

		buf := make([]byte, 256) // Smaller buffer for more responsive input
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err != io.EOF && ctx.Err() == nil {
					select {
					case errorChan <- fmt.Errorf("stdin error: %w", err):
					default:
					}
				}
				return
			}

			if n > 0 {
				// Make a copy to avoid data races
				data := make([]byte, n)
				copy(data, buf[:n])

				if err := stream.Send(&pb.ShellInput{StdinData: data}); err != nil {
					if ctx.Err() == nil {
						select {
						case errorChan <- fmt.Errorf("send error: %w", err):
						default:
						}
					}
					return
				}
			}
		}
	}()

	// Main event loop
mainLoop:
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGWINCH:
				// Terminal resized
				newWidth, newHeight, err := term.GetSize(int(os.Stdin.Fd()))
				if err == nil {
					stream.Send(&pb.ShellInput{
						Resize:     true,
						WindowRows: uint32(newHeight),
						WindowCols: uint32(newWidth),
					})
				}

			case os.Interrupt:
				// Ctrl+C - send to remote
				stream.Send(&pb.ShellInput{StdinData: []byte{0x03}})

			case syscall.SIGTERM:
				// Terminate gracefully
				stream.Send(&pb.ShellInput{Terminate: true})
				cancel()
				break mainLoop
			}

		case err := <-errorChan:
			// Got an error from one of the goroutines
			cancel()
			stream.CloseSend()
			wg.Wait()
			return err

		case <-stdoutDone:
			// Remote shell closed
			cancel()
			stream.CloseSend()
			wg.Wait()
			printShellGoodbye()
			return nil

		case <-stdinDone:
			// Local stdin closed (shouldn't happen normally)
			cancel()
			stream.CloseSend()
			wg.Wait()
			return nil
		}
	}

	// Cleanup
	cancel()
	stream.CloseSend()
	wg.Wait()
	return nil
}

// printRobustBanner shows a clean, informative banner
func printRobustBanner(agentName, address string, width, height int) {
	hostname, _ := os.Hostname()

	fmt.Printf("\033[2J\033[H") // Clear screen and move to top

	banner := pterm.DefaultBox.
		WithTitle("🚀 Sloth Runner Shell").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).
		Sprint(fmt.Sprintf(`
Agent:    %s
Address:  %s
Terminal: %dx%d
Client:   %s

[Ctrl+D or 'exit' to quit] [Ctrl+C to interrupt]
`,
			pterm.FgGreen.Sprintf(agentName),
			pterm.FgCyan.Sprintf(address),
			width, height,
			pterm.FgYellow.Sprintf(hostname),
		))

	fmt.Println(banner)
	fmt.Println()
}
