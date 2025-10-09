package group

import (
	"github.com/spf13/cobra"
)

// NewGroupCmd creates the group command
func NewGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Manage agent groups",
		Long: `Manage agent groups for organizing and controlling multiple agents.

Agent groups allow you to:
  • Group agents logically for easier management
  • Execute bulk operations on multiple agents
  • Create reusable templates for group creation
  • Set up auto-discovery rules
  • Configure webhooks for group events
  • Organize groups hierarchically`,
		Example: `  # List all groups
  sloth-runner group list

  # Create a new group
  sloth-runner group create production-web --description "Production web servers"

  # Add agents to a group
  sloth-runner group add-agent production-web server-01 server-02

  # Execute bulk operation
  sloth-runner group bulk production-web restart`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewShowCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewAddAgentCmd())
	cmd.AddCommand(NewRemoveAgentCmd())
	cmd.AddCommand(NewBulkCmd())
	cmd.AddCommand(NewTemplateCmd())
	cmd.AddCommand(NewAutoDiscoveryCmd())
	cmd.AddCommand(NewWebhookCmd())

	return cmd
}
