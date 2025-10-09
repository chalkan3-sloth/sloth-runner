package history

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/chalkan3-sloth/sloth-runner/internal/execution"
	"github.com/spf13/cobra"
)

func NewHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "history",
		Aliases: []string{"hist", "executions"},
		Short:   "View execution history",
		Long:    `View, search, and analyze workflow execution history.`,
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newShowCmd())
	cmd.AddCommand(newStatsCmd())
	cmd.AddCommand(newCleanupCmd())

	return cmd
}

func newListCmd() *cobra.Command {
	var workflow string
	var status string
	var agent string
	var group string
	var since string
	var limit int
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List execution history",
		Long: `List recent workflow executions with filtering options.

Available status filters: running, completed, failed, cancelled`,
		Example: `  # List recent executions
  sloth-runner history list

  # List executions for a specific workflow
  sloth-runner history list --workflow deploy-app

  # List failed executions
  sloth-runner history list --status failed

  # List executions from last 24 hours
  sloth-runner history list --since 24h

  # List with JSON output
  sloth-runner history list -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listExecutions(workflow, status, agent, group, since, limit, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&workflow, "workflow", "w", "", "Filter by workflow name")
	cmd.Flags().StringVarP(&status, "status", "s", "", "Filter by status (running|completed|failed|cancelled)")
	cmd.Flags().StringVarP(&agent, "agent", "a", "", "Filter by agent name")
	cmd.Flags().StringVarP(&group, "group", "g", "", "Filter by group name")
	cmd.Flags().StringVar(&since, "since", "", "Show executions since (e.g., 24h, 7d, 30d)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 50, "Number of executions to show")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")

	return cmd
}

func newShowCmd() *cobra.Command {
	var verbose bool
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "show <execution-id>",
		Short: "Show detailed execution information",
		Long:  `Display detailed information about a specific execution including all task results.`,
		Example: `  # Show execution details
  sloth-runner history show abc123

  # Show with full output
  sloth-runner history show abc123 --verbose

  # Show as JSON
  sloth-runner history show abc123 -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showExecution(args[0], verbose, outputFormat)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show full output from execution")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")

	return cmd
}

func newStatsCmd() *cobra.Command {
	var workflow string
	var since string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show execution statistics",
		Long:  `Display statistics about workflow executions including success rates and performance metrics.`,
		Example: `  # Show overall statistics
  sloth-runner history stats

  # Show stats for specific workflow
  sloth-runner history stats --workflow deploy-app

  # Show stats for last 7 days
  sloth-runner history stats --since 7d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showStats(workflow, since, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&workflow, "workflow", "w", "", "Filter by workflow name")
	cmd.Flags().StringVar(&since, "since", "", "Statistics since (e.g., 24h, 7d, 30d)")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")

	return cmd
}

func newCleanupCmd() *cobra.Command {
	var days int
	var force bool

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Clean up old execution history",
		Long:  `Delete execution records older than the specified number of days.`,
		Example: `  # Delete executions older than 30 days
  sloth-runner history cleanup --days 30

  # Delete without confirmation
  sloth-runner history cleanup --days 90 --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cleanupHistory(days, force)
		},
	}

	cmd.Flags().IntVarP(&days, "days", "d", 90, "Delete executions older than this many days")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")

	return cmd
}

func listExecutions(workflow, status, agent, group, since string, limit int, outputFormat string) error {
	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer db.Close()

	filters := make(map[string]interface{})
	if workflow != "" {
		filters["workflow"] = workflow
	}
	if status != "" {
		filters["status"] = status
	}
	if agent != "" {
		filters["agent"] = agent
	}
	if group != "" {
		filters["group"] = group
	}
	if since != "" {
		sinceTime, err := parseDuration(since)
		if err != nil {
			return fmt.Errorf("invalid since format: %w", err)
		}
		filters["since"] = sinceTime
	}

	executions, err := db.ListExecutions(filters, limit, 0)
	if err != nil {
		return fmt.Errorf("failed to list executions: %w", err)
	}

	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(executions)
	}

	if len(executions) == 0 {
		fmt.Println("No executions found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tWORKFLOW\tSTATUS\tAGENT/GROUP\tDURATION\tSTART TIME\tTASKS")
	fmt.Fprintln(w, "--\t--------\t------\t-----------\t--------\t----------\t-----")

	for _, exec := range executions {
		statusIcon := getStatusIcon(exec.Status)

		target := exec.AgentName
		if exec.GroupName != "" {
			target = fmt.Sprintf("group:%s", exec.GroupName)
		}
		if target == "" {
			target = "-"
		}

		duration := "-"
		if exec.Duration > 0 {
			duration = formatDuration(exec.Duration)
		}

		startTime := time.Unix(exec.StartTime, 0).Format("2006-01-02 15:04")

		tasks := fmt.Sprintf("%d/%d", exec.TasksSuccess, exec.TasksTotal)
		if exec.TasksFailed > 0 {
			tasks += fmt.Sprintf(" (%d failed)", exec.TasksFailed)
		}

		fmt.Fprintf(w, "%s\t%s\t%s %s\t%s\t%s\t%s\t%s\n",
			exec.ID[:8],
			exec.WorkflowName,
			statusIcon, exec.Status,
			target,
			duration,
			startTime,
			tasks,
		)
	}

	w.Flush()
	return nil
}

func showExecution(id string, verbose bool, outputFormat string) error {
	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer db.Close()

	exec, err := db.GetExecution(id)
	if err != nil {
		return fmt.Errorf("failed to get execution: %w", err)
	}

	tasks, err := db.GetTaskExecutions(id)
	if err != nil {
		return fmt.Errorf("failed to get task executions: %w", err)
	}

	if outputFormat == "json" {
		result := map[string]interface{}{
			"execution": exec,
			"tasks":     tasks,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Text format
	fmt.Printf("\n")
	fmt.Printf("Execution ID:    %s\n", exec.ID)
	fmt.Printf("Workflow:        %s\n", exec.WorkflowName)
	fmt.Printf("File:            %s\n", exec.WorkflowFile)
	fmt.Printf("Status:          %s %s\n", getStatusIcon(exec.Status), exec.Status)

	if exec.AgentName != "" {
		fmt.Printf("Agent:           %s\n", exec.AgentName)
	}
	if exec.GroupName != "" {
		fmt.Printf("Group:           %s\n", exec.GroupName)
	}
	if exec.User != "" {
		fmt.Printf("User:            %s\n", exec.User)
	}

	fmt.Printf("Start Time:      %s\n", time.Unix(exec.StartTime, 0).Format("2006-01-02 15:04:05"))
	if exec.EndTime > 0 {
		fmt.Printf("End Time:        %s\n", time.Unix(exec.EndTime, 0).Format("2006-01-02 15:04:05"))
	}
	if exec.Duration > 0 {
		fmt.Printf("Duration:        %s\n", formatDuration(exec.Duration))
	}

	fmt.Printf("Exit Code:       %d\n", exec.ExitCode)
	fmt.Printf("Tasks:           %d total, %d success, %d failed\n",
		exec.TasksTotal, exec.TasksSuccess, exec.TasksFailed)

	if exec.ErrorMessage != "" {
		fmt.Printf("\nError:\n%s\n", exec.ErrorMessage)
	}

	if verbose && exec.Output != "" {
		fmt.Printf("\nOutput:\n%s\n", exec.Output)
	}

	// Show tasks
	if len(tasks) > 0 {
		fmt.Printf("\nTasks:\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "TASK\tSTATUS\tDURATION\tCHANGED")
		fmt.Fprintln(w, "----\t------\t--------\t-------")

		for _, task := range tasks {
			statusIcon := getStatusIcon(task.Status)
			duration := "-"
			if task.Duration > 0 {
				duration = formatDuration(task.Duration)
			}
			changed := "no"
			if task.Changed {
				changed = "yes"
			}

			fmt.Fprintf(w, "%s\t%s %s\t%s\t%s\n",
				task.TaskName,
				statusIcon, task.Status,
				duration,
				changed,
			)

			if verbose && task.Error != "" {
				fmt.Fprintf(w, "  Error: %s\n", task.Error)
			}
		}

		w.Flush()
	}

	fmt.Printf("\n")
	return nil
}

func showStats(workflow, since string, outputFormat string) error {
	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer db.Close()

	filters := make(map[string]interface{})
	if workflow != "" {
		filters["workflow"] = workflow
	}
	if since != "" {
		sinceTime, err := parseDuration(since)
		if err != nil {
			return fmt.Errorf("invalid since format: %w", err)
		}
		filters["since"] = sinceTime
	}

	stats, err := db.GetStatistics(filters)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}

	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(stats)
	}

	// Text format
	fmt.Printf("\nExecution Statistics")
	if workflow != "" {
		fmt.Printf(" (workflow: %s)", workflow)
	}
	if since != "" {
		fmt.Printf(" (since: %s)", since)
	}
	fmt.Printf("\n\n")

	total := int(stats["total"].(int))
	completed := int(stats["completed"].(int))
	failed := int(stats["failed"].(int))
	running := int(stats["running"].(int))
	successRate := stats["success_rate"].(float64)

	avgDuration := int64(0)
	if avgVal, ok := stats["avg_duration"].(int64); ok {
		avgDuration = avgVal
	} else if avgVal, ok := stats["avg_duration"].(int); ok {
		avgDuration = int64(avgVal)
	}

	fmt.Printf("Total Executions:    %d\n", total)
	fmt.Printf("  Completed:         %d\n", completed)
	fmt.Printf("  Failed:            %d\n", failed)
	fmt.Printf("  Running:           %d\n", running)
	fmt.Printf("\n")
	fmt.Printf("Success Rate:        %.1f%%\n", successRate)
	if avgDuration > 0 {
		fmt.Printf("Average Duration:    %s\n", formatDuration(avgDuration))
	}
	fmt.Printf("\n")

	return nil
}

func cleanupHistory(days int, force bool) error {
	if !force {
		fmt.Printf("This will delete all executions older than %d days.\n", days)
		fmt.Print("Are you sure? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" && confirm != "y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer db.Close()

	deleted, err := db.DeleteOldExecutions(days)
	if err != nil {
		return fmt.Errorf("failed to delete old executions: %w", err)
	}

	fmt.Printf("Deleted %d old execution(s)\n", deleted)
	return nil
}

func getStatusIcon(status execution.ExecutionStatus) string {
	switch status {
	case execution.StatusCompleted:
		return "✅"
	case execution.StatusFailed:
		return "❌"
	case execution.StatusRunning:
		return "⏳"
	case execution.StatusCancelled:
		return "🚫"
	default:
		return "❓"
	}
}

func formatDuration(ms int64) string {
	duration := time.Duration(ms) * time.Millisecond

	if duration < time.Second {
		return fmt.Sprintf("%dms", ms)
	}
	if duration < time.Minute {
		return fmt.Sprintf("%.1fs", duration.Seconds())
	}
	if duration < time.Hour {
		return fmt.Sprintf("%.1fm", duration.Minutes())
	}
	return fmt.Sprintf("%.1fh", duration.Hours())
}

func parseDuration(s string) (int64, error) {
	var value int
	var unit string
	_, err := fmt.Sscanf(s, "%d%s", &value, &unit)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	var since time.Time

	switch unit {
	case "h":
		since = now.Add(-time.Duration(value) * time.Hour)
	case "d":
		since = now.AddDate(0, 0, -value)
	case "w":
		since = now.AddDate(0, 0, -value*7)
	case "m":
		since = now.AddDate(0, -value, 0)
	default:
		return 0, fmt.Errorf("unknown unit: %s (use h, d, w, or m)", unit)
	}

	return since.Unix(), nil
}
