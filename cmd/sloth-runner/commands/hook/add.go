package hook

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewAddCommand creates the hook add command
func NewAddCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <hook-name> --file <path-to-hook-file>",
		Short: "Add a new event hook",
		Long: `Add a new event hook that will be triggered when specific events occur.

The hook file should be a Lua script that defines event handlers.

Example:
  sloth-runner hook add notify-agent-join --file hooks/notify.lua`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]
			filePath, _ := cmd.Flags().GetString("file")
			eventType, _ := cmd.Flags().GetString("event")
			description, _ := cmd.Flags().GetString("description")
			stack, _ := cmd.Flags().GetString("stack")
			enabled, _ := cmd.Flags().GetBool("enabled")

			if filePath == "" {
				return fmt.Errorf("--file flag is required")
			}

			if eventType == "" {
				return fmt.Errorf("--event flag is required")
			}

			// Validate file exists
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return fmt.Errorf("invalid file path: %w", err)
			}

			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				return fmt.Errorf("hook file not found: %s", absPath)
			}

			// Validate event type
			validEvents := []string{
				"agent.registered",
				"agent.disconnected",
				"agent.heartbeat_failed",
				"agent.updated",
				"task.started",
				"task.completed",
				"task.failed",
				"workflow.started",
				"workflow.completed",
				"workflow.failed",
			}

			validEvent := false
			for _, e := range validEvents {
				if eventType == e {
					validEvent = true
					break
				}
			}

			if !validEvent {
				pterm.Error.Printf("Invalid event type: %s\n", eventType)
				pterm.Info.Println("Valid event types:")
				for _, e := range validEvents {
					pterm.Info.Printf("  - %s\n", e)
				}
				return fmt.Errorf("invalid event type")
			}

			// Create repository
			repo, err := hooks.NewRepository()
			if err != nil {
				return fmt.Errorf("failed to create repository: %w", err)
			}
			defer repo.Close()

			// Create hook
			hook := &hooks.Hook{
				Name:        hookName,
				Description: description,
				EventType:   hooks.EventType(eventType),
				FilePath:    absPath,
				Stack:       stack,
				Enabled:     enabled,
			}

			if err := repo.Add(hook); err != nil {
				return fmt.Errorf("failed to add hook: %w", err)
			}

			pterm.Success.Printf("Hook '%s' added successfully!\n", hookName)
			pterm.Info.Printf("Event type: %s\n", eventType)
			pterm.Info.Printf("File: %s\n", absPath)
			if stack != "" {
				pterm.Info.Printf("Stack: %s\n", stack)
			}
			pterm.Info.Printf("Enabled: %v\n", enabled)

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "Path to the hook file (required)")
	cmd.Flags().StringP("event", "e", "", "Event type to trigger the hook (required)")
	cmd.Flags().StringP("description", "d", "", "Hook description")
	cmd.Flags().StringP("stack", "s", "", "Stack name for hook isolation")
	cmd.Flags().Bool("enabled", true, "Enable the hook immediately")

	return cmd
}
