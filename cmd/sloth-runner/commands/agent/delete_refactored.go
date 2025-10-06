package agent

import (
	"context"
	"fmt"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// DeleteAgentOptions contains options for deleting an agent
type DeleteAgentOptions struct {
	AgentName string
}

// deleteAgentWithClient deletes an agent using an injected client (testable)
func deleteAgentWithClient(ctx context.Context, client AgentRegistryClient, opts DeleteAgentOptions) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deleting agent '%s'...", opts.AgentName))

	// Call UnregisterAgent on the registry
	resp, err := client.UnregisterAgent(ctx, &pb.UnregisterAgentRequest{
		AgentName: opts.AgentName,
	})

	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to delete agent: %v", err))
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	if !resp.Success {
		spinner.Fail("Failed to delete agent")
		return fmt.Errorf("delete failed: agent not found or already removed")
	}

	spinner.Success(fmt.Sprintf("Agent '%s' deleted successfully", opts.AgentName))
	return nil
}
