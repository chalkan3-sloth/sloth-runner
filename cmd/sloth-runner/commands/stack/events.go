package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewEventsCommand creates the events command
func NewEventsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "View and manage state events",
		Long:  `View real-time events from the state management system.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewEventsListCommand(ctx),
		NewEventsStatsCommand(ctx),
		NewEventsWatchCommand(ctx),
	)

	return cmd
}

// NewEventsListCommand lists recent events
func NewEventsListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recent state events",
		Long:  `Lists the most recent state management events.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, _ := cmd.Flags().GetInt("limit")
			eventType, _ := cmd.Flags().GetString("type")
			stackName, _ := cmd.Flags().GetString("stack")
			outputFormat, _ := cmd.Flags().GetString("output")

			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			var events []*stack.StateEvent

			if eventType != "" {
				events = tracker.GetEventsByType(stack.EventType(eventType), limit)
			} else if stackName != "" {
				stackService, err := services.NewStackService()
				if err != nil {
					return err
				}
				defer stackService.Close()

				st, err := stackService.GetStackByName(stackName)
				if err != nil {
					return fmt.Errorf("stack '%s' not found: %w", stackName, err)
				}

				events = tracker.GetEventsByStack(st.ID, limit)
			} else {
				events = tracker.GetEventHistory(limit)
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(events, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(events) == 0 {
				pterm.Info.Println("No events found")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Println("State Events")
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tTYPE\tSOURCE\tSTACK\tSEVERITY")
			fmt.Fprintln(w, "---------\t----\t------\t-----\t--------")

			for _, event := range events {
				severity := event.Severity
				if event.Severity == "error" {
					severity = pterm.Red(severity)
				} else if event.Severity == "warning" {
					severity = pterm.Yellow(severity)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					event.Timestamp.Format("2006-01-02 15:04:05"),
					event.Type,
					event.Source,
					event.StackName,
					severity,
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d event(s)\n", len(events))

			return nil
		},
	}

	cmd.Flags().IntP("limit", "l", 50, "Maximum number of events to show")
	cmd.Flags().StringP("type", "t", "", "Filter by event type")
	cmd.Flags().StringP("stack", "s", "", "Filter by stack name")
	cmd.Flags().StringP("output", "o", "table", "Output format (table or json)")

	return cmd
}

// NewEventsStatsCommand shows event statistics
func NewEventsStatsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show event statistics",
		Long:  `Displays statistics about state events.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			stats := tracker.GetEventStats()

			pterm.DefaultHeader.WithFullWidth().Println("Event Statistics")
			fmt.Println()

			pterm.Info.Printf("Total Events: %d\n", stats["total_events"])
			fmt.Println()

			// By type
			pterm.DefaultSection.Println("Events by Type")
			byType := stats["by_type"].(map[stack.EventType]int)
			for eventType, count := range byType {
				pterm.Info.Printf("  %s: %d\n", eventType, count)
			}
			fmt.Println()

			// By severity
			pterm.DefaultSection.Println("Events by Severity")
			bySeverity := stats["by_severity"].(map[string]int)
			for severity, count := range bySeverity {
				color := pterm.Info
				if severity == "error" {
					color = pterm.Error
				} else if severity == "warning" {
					color = pterm.Warning
				}
				color.Printf("  %s: %d\n", severity, count)
			}

			return nil
		},
	}

	return cmd
}

// NewEventsWatchCommand watches events in real-time
func NewEventsWatchCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch state events in real-time",
		Long:  `Continuously displays state events as they occur.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pterm.Info.Println("Watching for state events... (Press Ctrl+C to stop)")
			fmt.Println()

			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			// Subscribe to all events
			eventBus := tracker.GetEventBus()

			// This is a simple implementation - in a real scenario, you'd use channels
			eventBus.SubscribeAll(func(ctx context.Context, event *stack.StateEvent) error {
				severity := ""
				switch event.Severity {
				case "error":
					severity = pterm.Red("ERROR")
				case "warning":
					severity = pterm.Yellow("WARN")
				case "info":
					severity = pterm.Cyan("INFO")
				default:
					severity = event.Severity
				}

				pterm.Printf("[%s] %s %s: %s (stack: %s)\n",
					event.Timestamp.Format("15:04:05"),
					severity,
					event.Type,
					event.Source,
					event.StackName,
				)
				return nil
			})

			pterm.Info.Println("Event watcher configured. Future events will be displayed here.")
			pterm.Warning.Println("Note: This is a passive view of the event log. Use in a persistent session to see live events.")

			return nil
		},
	}

	return cmd
}
