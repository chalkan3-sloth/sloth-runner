package agent

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/telemetry"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// MetricsOptions contains options for metrics commands
type MetricsOptions struct {
	AgentName    string
	ShowSnapshot bool
	Writer       io.Writer
}

// DashboardOptions contains options for dashboard commands
type DashboardOptions struct {
	AgentName        string
	MetricsEndpoint  string
	Watch            bool
	Interval         int
	MetricsFetcher   func(string) (*telemetry.MetricsData, error)
	DashboardDisplay func(*telemetry.MetricsData, string)
}

// prometheusMetricsWithClient fetches prometheus metrics using injected client (testable)
func prometheusMetricsWithClient(ctx context.Context, client AgentRegistryClient, opts MetricsOptions) error {
	pterm.DefaultHeader.WithFullWidth().Printf("Prometheus Metrics - Agent: %s", opts.AgentName)
	fmt.Fprintln(opts.Writer)

	// Find agent address
	agentAddress, err := findAgentAddress(ctx, client, opts.AgentName)
	if err != nil {
		return err
	}

	// Extract host and build metrics endpoint
	host := extractHost(agentAddress)
	metricsEndpoint := fmt.Sprintf("http://%s:9090/metrics", host)

	if opts.ShowSnapshot {
		return displayMetricsSnapshot(ctx, client, opts.AgentName, opts.Writer)
	}

	return displayMetricsEndpoint(metricsEndpoint, host, opts.Writer)
}

// extractHost extracts host from address (removes port if present) (testable)
func extractHost(address string) string {
	if strings.Contains(address, ":") {
		return strings.Split(address, ":")[0]
	}
	return address
}

// displayMetricsSnapshot fetches and displays current metrics snapshot (testable)
func displayMetricsSnapshot(ctx context.Context, client AgentRegistryClient, agentName string, w io.Writer) error {
	pterm.Info.Println("📊 Fetching metrics snapshot...")
	fmt.Fprintln(w)

	// Execute curl command on agent to fetch metrics
	curlCmd := "curl -s http://localhost:9090/metrics"
	stream, err := client.ExecuteCommand(ctx, &pb.ExecuteCommandRequest{
		AgentName: agentName,
		Command:   curlCmd,
	})

	if err != nil {
		return fmt.Errorf("failed to fetch metrics: %v", err)
	}

	var output strings.Builder
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving metrics: %v", err)
		}
		output.WriteString(resp.GetStdoutChunk())
	}

	metrics := output.String()
	if strings.Contains(metrics, "Connection refused") || strings.Contains(metrics, "Failed to connect") {
		pterm.Error.Println("❌ Telemetry server is not running on this agent")
		fmt.Fprintln(w)
		pterm.Info.Println("💡 Start the agent with telemetry enabled:")
		fmt.Fprintf(w, "  sloth-runner agent start --name %s --telemetry\n", agentName)
		return nil
	}

	fmt.Fprintln(w, pterm.DefaultBox.WithTitle("Metrics Snapshot").Sprint(metrics))
	return nil
}

// displayMetricsEndpoint displays the metrics endpoint information (testable)
func displayMetricsEndpoint(metricsEndpoint, host string, w io.Writer) error {
	pterm.Success.Println("✅ Metrics Endpoint:")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s %s\n", pterm.Cyan("URL:"), pterm.Green(metricsEndpoint))
	fmt.Fprintln(w)

	pterm.Info.Println("📝 Usage:")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s\n", pterm.Gray("# View metrics in browser:"))
	fmt.Fprintf(w, "  open %s\n", metricsEndpoint)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s\n", pterm.Gray("# Fetch metrics via curl:"))
	fmt.Fprintf(w, "  curl %s\n", metricsEndpoint)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s\n", pterm.Gray("# Configure Prometheus scraper:"))
	fmt.Fprintln(w, "  - job_name: 'sloth-runner-agents'")
	fmt.Fprintln(w, "    static_configs:")
	fmt.Fprintf(w, "      - targets: ['%s:9090']\n", host)
	fmt.Fprintln(w)

	pterm.Info.Println("💡 Tip: Use --snapshot flag to display current metrics")
	return nil
}

// grafanaDashboardWithClient displays grafana dashboard using injected client and functions (testable)
func grafanaDashboardWithClient(ctx context.Context, client AgentRegistryClient, opts DashboardOptions) error {
	// Find agent address
	agentAddress, err := findAgentAddress(ctx, client, opts.AgentName)
	if err != nil {
		return err
	}

	// Extract host and build metrics endpoint
	host := extractHost(agentAddress)
	metricsEndpoint := fmt.Sprintf("http://%s:9090/metrics", host)

	// Use injected fetcher or default
	fetchMetrics := opts.MetricsFetcher
	if fetchMetrics == nil {
		fetchMetrics = telemetry.FetchMetrics
	}

	// Use injected display or default
	displayDashboard := opts.DashboardDisplay
	if displayDashboard == nil {
		displayDashboard = telemetry.DisplayDashboard
	}

	// Function to fetch and display dashboard
	showDashboard := func() error {
		data, err := fetchMetrics(metricsEndpoint)
		if err != nil {
			pterm.Error.Printf("❌ Failed to fetch metrics: %v\n", err)
			fmt.Println()
			pterm.Info.Println("💡 Ensure the agent is running with telemetry enabled:")
			fmt.Printf("  sloth-runner agent start --name %s --telemetry\n", opts.AgentName)
			return err
		}

		displayDashboard(data, opts.AgentName)
		return nil
	}

	// Display dashboard in watch mode or one-time
	if opts.Watch {
		return fmt.Errorf("watch mode not supported in test environment")
	}

	return showDashboard()
}
