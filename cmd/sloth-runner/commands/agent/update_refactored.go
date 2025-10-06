package agent

import (
	"context"
	"fmt"
	"io"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// UpdateAgentOptions contains options for updating an agent
type UpdateAgentOptions struct {
	AgentName     string
	TargetVersion string
	Restart       bool
	Writer        io.Writer
}

// UpdateAgentResult contains the result of an agent update
type UpdateAgentResult struct {
	Success    bool
	Message    string
	OldVersion string
	NewVersion string
}

// updateAgentWithClients updates an agent using injected clients (testable)
func updateAgentWithClients(
	ctx context.Context,
	registryClient AgentRegistryClient,
	agentClientFactory func(string) (AgentClient, func(), error),
	opts UpdateAgentOptions,
) (*UpdateAgentResult, error) {
	pterm.DefaultHeader.WithFullWidth().Printf("Agent Update - %s", opts.AgentName)
	fmt.Fprintln(opts.Writer)

	spinner, _ := pterm.DefaultSpinner.Start("Connecting to master server...")

	// Get agent info from master
	agentAddress, err := findAgentAddress(ctx, registryClient, opts.AgentName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to find agent: %v", err))
		return nil, err
	}

	spinner.UpdateText(fmt.Sprintf("Connecting to agent at %s...", agentAddress))

	// Connect to agent
	agentClient, cleanup, err := agentClientFactory(agentAddress)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to connect to agent: %v", err))
		return nil, fmt.Errorf("failed to connect to agent at %s: %w", agentAddress, err)
	}
	defer cleanup()

	spinner.UpdateText("Initiating agent update...")

	// Call UpdateAgent
	result, err := performAgentUpdate(ctx, agentClient, opts)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to initiate update: %v", err))
		return nil, err
	}

	if !result.Success {
		spinner.Fail("Update failed")
		pterm.Error.Println(result.Message)
		return result, fmt.Errorf("update failed: %s", result.Message)
	}

	spinner.Success("Update completed successfully")
	fmt.Fprintln(opts.Writer)

	// Display update summary
	displayUpdateSummary(result, opts, opts.Writer)

	pterm.Success.Printf("✅ Agent %s updated successfully\n", opts.AgentName)

	return result, nil
}

// findAgentAddress finds an agent's address from the registry (testable)
func findAgentAddress(ctx context.Context, client AgentRegistryClient, agentName string) (string, error) {
	listResp, err := client.ListAgents(ctx, &pb.ListAgentsRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to list agents: %w", err)
	}

	for _, agent := range listResp.GetAgents() {
		if agent.GetAgentName() == agentName {
			return agent.GetAgentAddress(), nil
		}
	}

	return "", fmt.Errorf("agent '%s' not found", agentName)
}

// performAgentUpdate executes the update on the agent (testable)
func performAgentUpdate(ctx context.Context, client AgentClient, opts UpdateAgentOptions) (*UpdateAgentResult, error) {
	resp, err := client.UpdateAgent(ctx, &pb.UpdateAgentRequest{
		TargetVersion: opts.TargetVersion,
		Force:         false,
		SkipRestart:   !opts.Restart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update agent: %w", err)
	}

	return &UpdateAgentResult{
		Success:    resp.Success,
		Message:    resp.Message,
		OldVersion: resp.OldVersion,
		NewVersion: resp.NewVersion,
	}, nil
}

// displayUpdateSummary displays the update summary (testable)
func displayUpdateSummary(result *UpdateAgentResult, opts UpdateAgentOptions, w io.Writer) {
	pterm.Info.Println("Update Summary:")
	fmt.Fprintf(w, "  Agent:            %s\n", pterm.Cyan(opts.AgentName))
	fmt.Fprintf(w, "  Previous Version: %s\n", pterm.Yellow(result.OldVersion))
	fmt.Fprintf(w, "  New Version:      %s\n", pterm.Green(result.NewVersion))

	if !opts.Restart {
		fmt.Fprintf(w, "  Status:           %s\n", pterm.Yellow("Update complete, restart skipped"))
	} else {
		fmt.Fprintf(w, "  Status:           %s\n", pterm.Green("Agent updated and restarted"))
	}

	fmt.Fprintln(w)

	// Show additional info if available
	if result.Message != "" {
		pterm.Info.Println("Details:")
		fmt.Fprintf(w, "  %s\n", result.Message)
		fmt.Fprintln(w)
	}
}
