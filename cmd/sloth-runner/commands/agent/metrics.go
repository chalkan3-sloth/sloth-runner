package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewMetricsCommand creates the agent metrics command
func NewMetricsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Manage agent metrics",
		Long:  `Manages metrics collection and export for agents.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(newPrometheusCommand())
	cmd.AddCommand(newGrafanaCommand())

	return cmd
}

func newPrometheusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prom <agent_name>",
		Short: "Get Prometheus metrics endpoint for an agent",
		Long:  `Retrieves the Prometheus metrics endpoint URL for a specific agent, or displays current metrics snapshot.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr, _ := cmd.Flags().GetString("master")
			showSnapshot, _ := cmd.Flags().GetBool("snapshot")

			return prometheusMetrics(agentName, masterAddr, showSnapshot)
		},
	}

	cmd.Flags().String("master", "localhost:50051", "Master server address")
	cmd.Flags().Bool("snapshot", false, "Display current metrics snapshot")

	return cmd
}

func newGrafanaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grafana <agent_name>",
		Short: "Display detailed metrics dashboard for an agent",
		Long:  `Shows a comprehensive terminal-based dashboard with detailed graphs and metrics visualization.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr, _ := cmd.Flags().GetString("master")
			watch, _ := cmd.Flags().GetBool("watch")
			interval, _ := cmd.Flags().GetInt("interval")

			return grafanaDashboard(agentName, masterAddr, watch, interval)
		},
	}

	cmd.Flags().String("master", "localhost:50051", "Master server address")
	cmd.Flags().Bool("watch", false, "Continuously update dashboard")
	cmd.Flags().Int("interval", 5, "Refresh interval in seconds (for watch mode)")

	return cmd
}

func prometheusMetrics(agentName, masterAddr string, showSnapshot bool) error {
	// Connect to master to get agent address
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create connection factory and get client
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	// Use refactored function with injected client
	opts := MetricsOptions{
		AgentName:    agentName,
		ShowSnapshot: showSnapshot,
		Writer:       os.Stdout,
	}

	return prometheusMetricsWithClient(ctx, client, opts)
}

func grafanaDashboard(agentName, masterAddr string, watch bool, interval int) error {
	// Connect to master to get agent address
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create connection factory and get client
	factory := NewDefaultConnectionFactory()
	client, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	// If watch mode, handle it here with loop
	if watch {
		pterm.Info.Printf("🔄 Watching metrics for %s (refresh every %ds, press Ctrl+C to stop)\n", agentName, interval)
		fmt.Println()

		for {
			// Clear screen for watch mode
			fmt.Print("\033[H\033[2J")

			opts := DashboardOptions{
				AgentName: agentName,
				Watch:     false, // Set to false to avoid recursion
				Interval:  interval,
			}

			if err := grafanaDashboardWithClient(ctx, client, opts); err != nil {
				return err
			}

			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	// One-time display
	opts := DashboardOptions{
		AgentName: agentName,
		Watch:     false,
		Interval:  interval,
	}

	return grafanaDashboardWithClient(ctx, client, opts)
}
