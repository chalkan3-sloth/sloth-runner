package state

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewWorkflowCommand creates the workflow state management command
func NewWorkflowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflow execution states (Terraform/Pulumi-like)",
		Long:  `Workflow state management provides Terraform/Pulumi-like state tracking for workflow executions, including versioning, drift detection, and rollback capabilities.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all workflow state subcommands
	cmd.AddCommand(
		NewWorkflowListCommand(ctx),
		NewWorkflowShowCommand(ctx),
		NewWorkflowVersionsCommand(ctx),
		NewWorkflowRollbackCommand(ctx),
		NewWorkflowDriftCommand(ctx),
		NewWorkflowResourcesCommand(ctx),
		NewWorkflowOutputsCommand(ctx),
		NewWorkflowDeleteCommand(ctx),
	)

	return cmd
}
