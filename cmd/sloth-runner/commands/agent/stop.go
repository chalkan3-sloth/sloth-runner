package agent

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewStopCommand creates the agent stop command
func NewStopCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "stop <agent-name>",
		Short: "Stops a running agent",
		Long:  `Stops a running agent by sending a shutdown request.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			masterAddr, _ := cmd.Flags().GetString("master")

			// Create agent service
			agentService := services.NewAgentService(masterAddr)

			// Stop agent
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Stopping agent '%s'...", agentName))
			if err := agentService.StopAgent(agentName); err != nil {
				spinner.Fail(fmt.Sprintf("Failed to stop agent: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("Agent '%s' stopped successfully", agentName))
			return nil
		},
	}
}

func init() {
	// This will be called from NewAgentCommand
}
