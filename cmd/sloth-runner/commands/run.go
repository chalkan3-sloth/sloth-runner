package commands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/handlers"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
)

// NewRunCommand creates the run command
// This demonstrates the Command Pattern with dependency injection
func NewRunCommand(ctx *AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <stack-name> [--file <workflow.sloth>] [--sloth <name>] [--ssh <profile>]",
		Short: "Run sloth-runner tasks with a stack (required)",
		Long: `Run sloth-runner tasks from Lua files with a stack for state management.
A stack name is REQUIRED for all executions to track state and history.
Tasks can be executed locally or remotely via SSH profiles.
When using --ssh, tasks will be executed on the remote host.

You can use a saved sloth file with --sloth <name> instead of --file.
If --sloth is specified, --file will be ignored.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract flags
			filePath, _ := cmd.Flags().GetString("file")
			slothName, _ := cmd.Flags().GetString("sloth")
			values, _ := cmd.Flags().GetString("values")
			yesFlag, _ := cmd.Flags().GetBool("yes")
			interactive, _ := cmd.Flags().GetBool("interactive")
			outputStyle, _ := cmd.Flags().GetString("output")
			debug, _ := cmd.Flags().GetBool("debug")
			delegateToHosts, _ := cmd.Flags().GetStringArray("delegate-to")
			sshProfile, _ := cmd.Flags().GetString("ssh")
			sshPasswordStdin, _ := cmd.Flags().GetBool("ssh-password-stdin")

			// Configure log level based on debug flag
			if debug {
				pterm.DefaultLogger.Level = pterm.LogLevelDebug
				slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			} else {
				pterm.DefaultLogger.Level = pterm.LogLevelInfo
				slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			}

			// Handle sloth reference - if --sloth is specified, use it instead of --file
			if slothName != "" {
				slothService, err := services.NewSlothService()
				if err != nil {
					return fmt.Errorf("failed to initialize sloth service: %w", err)
				}
				defer slothService.Close()

				// Get sloth content and write to temp file
				content, err := slothService.UseSloth(cmd.Context(), slothName)
				if err != nil {
					return fmt.Errorf("failed to use sloth '%s': %w", slothName, err)
				}

				// Write to temporary file
				tmpFile, err := slothService.WriteContentToFile(content)
				if err != nil {
					return fmt.Errorf("failed to create temp file from sloth: %w", err)
				}
				defer os.Remove(tmpFile) // Clean up temp file after execution

				// Use temp file as filePath
				filePath = tmpFile

				if debug {
					slog.Debug("using sloth file", "name", slothName, "temp_file", tmpFile)
				}
			}

			// Get stack name from first argument
			stackName := args[0]

			// Determine output writer
			writer := cmd.OutOrStdout()
			if ctx.TestMode && ctx.OutputWriter != nil {
				writer = ctx.OutputWriter
			}

			// Create stack service
			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			// Create handler configuration
			config := &handlers.RunConfig{
				StackName:        stackName,
				FilePath:         filePath,
				Values:           values,
				Interactive:      interactive,
				OutputStyle:      outputStyle,
				Debug:            debug,
				DelegateToHosts:  delegateToHosts,
				SSHProfile:       sshProfile,
				SSHPasswordStdin: sshPasswordStdin,
				YesFlag:          yesFlag,
				Context:          cmd.Context(),
				Writer:           writer,
				AgentRegistry:    ctx.AgentRegistry,
			}

			// Create and execute handler
			handler := handlers.NewRunHandler(stackService, config)
			return handler.Execute()
		},
	}

	// Add flags
	cmd.Flags().StringP("file", "f", "", "Path to the Lua task file")
	cmd.Flags().String("sloth", "", "Name of saved sloth file to use (takes precedence over --file)")
	cmd.Flags().StringP("values", "v", "", "Path to the values file")
	cmd.Flags().Bool("yes", false, "Skip confirmation prompts")
	cmd.Flags().Bool("interactive", false, "Run in interactive mode")
	cmd.Flags().StringP("output", "o", "basic", "Output style: basic, enhanced, rich, modern, json")
	cmd.Flags().Bool("debug", false, "Enable debug logging (shows technical details)")
	cmd.Flags().StringArrayP("delegate-to", "d", []string{}, "Execute tasks on specified agents (can be used multiple times)")
	cmd.Flags().String("ssh", "", "SSH profile name for remote execution")
	cmd.Flags().Bool("ssh-password-stdin", false, "Read SSH password from stdin (must be followed by -)")

	return cmd
}
