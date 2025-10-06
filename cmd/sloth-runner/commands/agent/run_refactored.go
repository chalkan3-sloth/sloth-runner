package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// RunCommandOptions contains options for running a command
type RunCommandOptions struct {
	AgentName    string
	Command      string
	OutputFormat string
	OutputWriter io.Writer
	ErrorWriter  io.Writer
}

// CommandResult holds the result of command execution
type CommandResult struct {
	Success     bool
	ExitCode    int32
	Stdout      string
	Stderr      string
	Error       string
	HasFinished bool
}

// runCommandWithClient executes a command using an injected client (testable)
func runCommandWithClient(ctx context.Context, client AgentRegistryClient, opts RunCommandOptions) error {
	// Show header for text output
	if opts.OutputFormat != "json" {
		pterm.Info.WithWriter(opts.OutputWriter).Printf("🚀 Executing on agent: %s\n", opts.AgentName)
		pterm.Info.WithWriter(opts.OutputWriter).Printf("📝 Command: %s\n", opts.Command)
		fmt.Fprintln(opts.OutputWriter)
	}

	// Execute command
	stream, err := client.ExecuteCommand(ctx, &pb.ExecuteCommandRequest{
		AgentName: opts.AgentName,
		Command:   opts.Command,
	})
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	// Process stream
	result, err := processCommandStream(stream, opts.OutputFormat, opts.OutputWriter, opts.ErrorWriter)
	if err != nil {
		return err
	}

	// Format output
	if opts.OutputFormat == "json" {
		return formatCommandResultJSON(result, opts.AgentName, opts.Command, opts.OutputWriter)
	}

	return formatCommandResultText(result, opts.AgentName, opts.OutputWriter)
}

// processCommandStream processes the command execution stream (testable)
func processCommandStream(stream pb.AgentRegistry_ExecuteCommandClient, outputFormat string, outWriter, errWriter io.Writer) (*CommandResult, error) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	result := &CommandResult{
		ExitCode:    -1,
		HasFinished: false,
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("stream error: %w", err)
		}

		// Handle stdout
		if resp.GetStdoutChunk() != "" {
			if outputFormat == "json" {
				stdoutBuffer.WriteString(resp.GetStdoutChunk())
			} else {
				fmt.Fprint(outWriter, resp.GetStdoutChunk())
			}
		}

		// Handle stderr
		if resp.GetStderrChunk() != "" {
			if outputFormat == "json" {
				stderrBuffer.WriteString(resp.GetStderrChunk())
			} else {
				fmt.Fprint(errWriter, resp.GetStderrChunk())
			}
		}

		// Handle error
		if resp.GetError() != "" {
			result.Error = resp.GetError()
		}

		// Handle completion
		if resp.GetFinished() {
			result.ExitCode = resp.GetExitCode()
			result.HasFinished = true
			break
		}
	}

	// Store buffered output
	result.Stdout = stdoutBuffer.String()
	result.Stderr = stderrBuffer.String()

	// Determine success
	result.Success = (result.HasFinished && result.ExitCode == 0) || (!result.HasFinished && result.Error == "")

	return result, nil
}

// formatCommandResultJSON formats command result as JSON (testable)
func formatCommandResultJSON(result *CommandResult, agentName, command string, w io.Writer) error {
	output := map[string]interface{}{
		"agent":     agentName,
		"command":   command,
		"success":   result.Success,
		"exit_code": result.ExitCode,
		"stdout":    result.Stdout,
		"stderr":    result.Stderr,
		"error":     result.Error,
		"finished":  result.HasFinished,
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON output: %w", err)
	}

	fmt.Fprintln(w, string(jsonOutput))

	if !result.Success {
		return fmt.Errorf("command execution failed")
	}

	return nil
}

// formatCommandResultText formats command result as human-readable text (testable)
func formatCommandResultText(result *CommandResult, agentName string, w io.Writer) error {
	fmt.Fprintln(w)

	if result.Success {
		pterm.Success.WithWriter(w).Printf("✅ Command completed successfully on agent %s", agentName)
		if result.HasFinished {
			fmt.Fprintf(w, " (exit code: %d)\n", result.ExitCode)
		} else {
			fmt.Fprintln(w)
		}
		return nil
	}

	// Handle failure cases
	if result.HasFinished && result.ExitCode != 0 {
		pterm.Error.WithWriter(w).Printf("❌ Command failed on agent %s (exit code: %d)\n", agentName, result.ExitCode)
	} else if result.Error != "" {
		pterm.Error.WithWriter(w).Printf("❌ Command failed on agent %s: %s\n", agentName, result.Error)
	} else {
		pterm.Error.WithWriter(w).Printf("❌ Command failed on agent %s\n", agentName)
	}

	return fmt.Errorf("command execution failed")
}
