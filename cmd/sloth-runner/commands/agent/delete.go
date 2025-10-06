package agent

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates the agent delete command
func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <agent-name>",
		Short: "Deletes an agent from the registry",
		Long:  `Removes an agent from the registry. This does not stop the agent if it's running.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			force, _ := cmd.Flags().GetBool("force")
			masterAddr, _ := cmd.Flags().GetString("master")

			// Ask for confirmation unless --force is used
			if !force {
				confirm := false
				prompt := &survey.Confirm{
					Message: fmt.Sprintf("Are you sure you want to delete agent '%s'?", agentName),
					Default: false,
				}
				if err := survey.AskOne(prompt, &confirm); err != nil {
					return err
				}
				if !confirm {
					pterm.Warning.Println("Deletion cancelled")
					return nil
				}
			}

			// Create agent service
			agentService := services.NewAgentService(masterAddr)

			// Delete agent
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deleting agent '%s'...", agentName))
			if err := agentService.DeleteAgent(agentName); err != nil {
				spinner.Fail(fmt.Sprintf("Failed to delete agent: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("Agent '%s' deleted successfully", agentName))
			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")
	cmd.Flags().String("master", "192.168.1.29:50053", "Master registry address")

	return cmd
}
