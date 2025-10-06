package agent

import (
	"context"
	"fmt"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// StopAgentOptions contains options for stopping an agent
type StopAgentOptions struct {
	AgentName string
}

// stopAgentWithClient stops an agent using an injected client (testable)
func stopAgentWithClient(ctx context.Context, client AgentRegistryClient, opts StopAgentOptions) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Stopping agent '%s'...", opts.AgentName))

	// Call StopAgent on the registry
	resp, err := client.StopAgent(ctx, &pb.StopAgentRequest{
		AgentName: opts.AgentName,
	})

	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to stop agent: %v", err))
		return fmt.Errorf("failed to stop agent: %w", err)
	}

	if !resp.Success {
		spinner.Fail("Failed to stop agent")
		return fmt.Errorf("stop failed: %s", resp.Message)
	}

	spinner.Success(fmt.Sprintf("Agent '%s' stopped successfully", opts.AgentName))
	return nil
}
