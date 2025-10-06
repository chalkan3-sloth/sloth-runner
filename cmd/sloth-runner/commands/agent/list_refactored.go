package agent

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// ListAgentsOptions contains options for listing agents
type ListAgentsOptions struct {
	Writer io.Writer
}

// listAgentsWithClient lists agents using an injected client (testable)
func listAgentsWithClient(ctx context.Context, client AgentRegistryClient, opts ListAgentsOptions) error {
	// Request agents list
	resp, err := client.ListAgents(ctx, &pb.ListAgentsRequest{})
	if err != nil {
		return fmt.Errorf("failed to list agents: %w", err)
	}

	agents := resp.GetAgents()

	if len(agents) == 0 {
		fmt.Fprintln(opts.Writer, "No agents registered.")
		return nil
	}

	return formatAgentsTable(agents, opts.Writer)
}

// formatAgentsTable formats agents in table format (testable)
func formatAgentsTable(agents []*pb.AgentInfo, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	// Write header
	fmt.Fprintln(tw, "AGENT NAME\tADDRESS\tSTATUS\tVERSION\tUPDATE STATUS\tLAST HEARTBEAT\tLAST INFO COLLECTED")
	fmt.Fprintln(tw, "------------\t----------\t------\t-------\t-------------\t--------------\t-------------------")

	// Write agent rows
	for _, agent := range agents {
		status := agent.GetStatus()
		coloredStatus := formatStatus(status)

		lastHeartbeat := formatTimestamp(agent.GetLastHeartbeat(), "N/A")
		lastInfoCollected := formatTimestamp(agent.GetLastInfoCollected(), "Never")

		version := agent.GetVersion()
		if version == "" {
			version = "unknown"
		}

		// TODO: Implement version comparison logic
		updateStatus := "-"

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			agent.GetAgentName(),
			agent.GetAgentAddress(),
			coloredStatus,
			version,
			updateStatus,
			lastHeartbeat,
			lastInfoCollected)
	}

	return tw.Flush()
}

// formatStatus formats agent status with color (testable)
func formatStatus(status string) string {
	if status == "Active" {
		return pterm.Green(status)
	}
	return pterm.Red(status)
}

// formatTimestamp formats a unix timestamp or returns default value (testable)
func formatTimestamp(timestamp int64, defaultValue string) string {
	if timestamp > 0 {
		return time.Unix(timestamp, 0).Format(time.RFC3339)
	}
	return defaultValue
}
