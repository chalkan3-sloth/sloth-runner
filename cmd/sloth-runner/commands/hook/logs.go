//go:build cgo
// +build cgo

package hook

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewLogsCommand creates the logs command
func NewLogsCommand(ctx *commands.AppContext) *cobra.Command {
	var (
		limit      int
		showOutput bool
		showError  bool
		onlyFailed bool
		format     string
	)

	cmd := &cobra.Command{
		Use:   "logs [hook-name]",
		Short: "Show execution logs for a hook",
		Long: `Display the execution history and logs for a specific hook.

This command shows:
- Execution timestamp
- Success/failure status
- Execution duration
- Output (if --output flag is used)
- Error messages (if --error flag is used)

Useful for debugging hooks and understanding their behavior.`,
		Example: `  # Show last 10 executions of a hook
  sloth-runner hook logs file_changed_alert

  # Show last 50 executions
  sloth-runner hook logs file_changed_alert --limit 50

  # Show executions with output
  sloth-runner hook logs file_changed_alert --output

  # Show only failed executions
  sloth-runner hook logs file_changed_alert --only-failed

  # Show executions with full details (output + errors)
  sloth-runner hook logs file_changed_alert --output --error

  # JSON format
  sloth-runner hook logs file_changed_alert --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]

			// Initialize repository
			repo, err := hooks.NewRepository()
			if err != nil {
				return fmt.Errorf("failed to open hooks database: %w", err)
			}
			defer repo.Close()

			// Get hook by name
			hook, err := repo.GetByName(hookName)
			if err != nil {
				return fmt.Errorf("hook not found: %s", hookName)
			}

			// Get execution history
			executions, err := repo.GetExecutionHistory(hook.ID, limit)
			if err != nil {
				return fmt.Errorf("failed to get execution history: %w", err)
			}

			// Filter only failed if requested
			if onlyFailed {
				filtered := make([]*hooks.HookResult, 0)
				for _, exec := range executions {
					if !exec.Success {
						filtered = append(filtered, exec)
					}
				}
				executions = filtered
			}

			if len(executions) == 0 {
				if onlyFailed {
					pterm.Success.Printf("No failed executions found for hook '%s'\n", hookName)
				} else {
					pterm.Info.Printf("No execution history found for hook '%s'\n", hookName)
				}
				return nil
			}

			// Display based on format
			switch format {
			case "json":
				return displayLogsJSON(hook, executions)
			case "table", "":
				return displayLogsTable(hook, executions, showOutput, showError)
			default:
				return fmt.Errorf("unsupported format: %s (use 'table' or 'json')", format)
			}
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "Maximum number of executions to show")
	cmd.Flags().BoolVarP(&showOutput, "output", "o", false, "Show hook output")
	cmd.Flags().BoolVarP(&showError, "error", "e", false, "Show error messages")
	cmd.Flags().BoolVar(&onlyFailed, "only-failed", false, "Show only failed executions")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json)")

	return cmd
}

func displayLogsTable(hook *hooks.Hook, executions []*hooks.HookResult, showOutput, showError bool) error {
	// Header
	pterm.DefaultHeader.WithFullWidth().Printf("Hook Execution Logs: %s", hook.Name)
	fmt.Println()

	pterm.Info.Printf("Event Type: %s\n", hook.EventType)
	pterm.Info.Printf("Total Executions Shown: %d\n", len(executions))

	successCount := 0
	failedCount := 0
	for _, exec := range executions {
		if exec.Success {
			successCount++
		} else {
			failedCount++
		}
	}
	pterm.Info.Printf("Success: %d, Failed: %d\n", successCount, failedCount)
	fmt.Println()

	// Table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Table header
	if showOutput || showError {
		fmt.Fprintln(w, "TIMESTAMP\tSTATUS\tDURATION\tDETAILS")
		fmt.Fprintln(w, "─────────────────────\t──────\t────────\t───────")
	} else {
		fmt.Fprintln(w, "TIMESTAMP\tSTATUS\tDURATION")
		fmt.Fprintln(w, "─────────────────────\t──────\t────────")
	}

	// Table rows
	for _, exec := range executions {
		timestamp := exec.ExecutedAt.Format("2006-01-02 15:04:05")
		status := "✅ SUCCESS"
		if !exec.Success {
			status = "❌ FAILED"
		}
		duration := exec.Duration.String()

		if showOutput || showError {
			details := ""
			if showOutput && exec.Output != "" {
				if len(exec.Output) > 60 {
					details = exec.Output[:57] + "..."
				} else {
					details = exec.Output
				}
			}
			if showError && exec.Error != "" {
				if details != "" {
					details += " | "
				}
				if len(exec.Error) > 60 {
					details += "ERROR: " + exec.Error[:50] + "..."
				} else {
					details += "ERROR: " + exec.Error
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", timestamp, status, duration, details)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\n", timestamp, status, duration)
		}
	}

	w.Flush()
	fmt.Println()

	// Show detailed output/errors if requested
	if (showOutput || showError) && len(executions) > 0 {
		fmt.Println()
		pterm.DefaultHeader.WithFullWidth().Println("Detailed Execution Logs")
		fmt.Println()

		for i, exec := range executions {
			fmt.Printf("─── Execution #%d ───\n", i+1)
			fmt.Printf("Timestamp: %s\n", exec.ExecutedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Status: ")
			if exec.Success {
				pterm.Success.Println("SUCCESS")
			} else {
				pterm.Error.Println("FAILED")
			}
			fmt.Printf("Duration: %s\n", exec.Duration)

			if showOutput && exec.Output != "" {
				fmt.Println()
				pterm.Info.Println("Output:")
				fmt.Println(exec.Output)
			}

			if showError && exec.Error != "" {
				fmt.Println()
				pterm.Error.Println("Error:")
				fmt.Println(exec.Error)
			}

			if i < len(executions)-1 {
				fmt.Println()
			}
		}
	}

	return nil
}

func displayLogsJSON(hook *hooks.Hook, executions []*hooks.HookResult) error {
	type LogEntry struct {
		HookName   string    `json:"hook_name"`
		HookID     string    `json:"hook_id"`
		EventType  string    `json:"event_type"`
		Timestamp  time.Time `json:"timestamp"`
		Success    bool      `json:"success"`
		Duration   string    `json:"duration"`
		DurationMS int64     `json:"duration_ms"`
		Output     string    `json:"output,omitempty"`
		Error      string    `json:"error,omitempty"`
	}

	type LogsOutput struct {
		Hook       string      `json:"hook"`
		TotalShown int         `json:"total_shown"`
		Success    int         `json:"success_count"`
		Failed     int         `json:"failed_count"`
		Executions []LogEntry  `json:"executions"`
	}

	successCount := 0
	failedCount := 0
	var entries []LogEntry

	for _, exec := range executions {
		if exec.Success {
			successCount++
		} else {
			failedCount++
		}

		entries = append(entries, LogEntry{
			HookName:   hook.Name,
			HookID:     hook.ID,
			EventType:  string(hook.EventType),
			Timestamp:  exec.ExecutedAt,
			Success:    exec.Success,
			Duration:   exec.Duration.String(),
			DurationMS: exec.Duration.Milliseconds(),
			Output:     exec.Output,
			Error:      exec.Error,
		})
	}

	output := LogsOutput{
		Hook:       hook.Name,
		TotalShown: len(executions),
		Success:    successCount,
		Failed:     failedCount,
		Executions: entries,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
