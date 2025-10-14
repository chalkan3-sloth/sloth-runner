//go:build cgo
// +build cgo

package stack

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewNewCommand creates the stack new command
func NewNewCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <stack-name>",
		Short: "Create a new workflow stack",
		Long:  `Create a new workflow stack with the specified name and optional configuration.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			description, _ := cmd.Flags().GetString("description")
			workflowFile, _ := cmd.Flags().GetString("workflow-file")
			stackVersion, _ := cmd.Flags().GetString("version")

			stackManager, err := stack.NewStackManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize stack manager: %w", err)
			}
			defer stackManager.Close()

			// Check if stack already exists
			if _, err := stackManager.GetStackByName(stackName); err == nil {
				return fmt.Errorf("stack '%s' already exists", stackName)
			}

			// Set defaults
			if description == "" {
				description = fmt.Sprintf("Workflow stack: %s", stackName)
			}
			if stackVersion == "" {
				stackVersion = "1.0.0"
			}

			// Create new stack
			stackID := uuid.New().String()
			newStack := &stack.StackState{
				ID:            stackID,
				Name:          stackName,
				Description:   description,
				Version:       stackVersion,
				Status:        "created",
				WorkflowFile:  workflowFile,
				TaskResults:   make(map[string]interface{}),
				Outputs:       make(map[string]interface{}),
				Configuration: make(map[string]interface{}),
				Metadata:      make(map[string]interface{}),
			}

			if err := stackManager.CreateStack(newStack); err != nil {
				return fmt.Errorf("failed to create stack: %w", err)
			}

			// Show success message
			pterm.Success.Printf("Stack '%s' created successfully.\n", stackName)
			pterm.Printf("\n")
			pterm.Printf("Stack Details:\n")
			pterm.Printf("  Name: %s\n", stackName)
			pterm.Printf("  ID: %s\n", stackID)
			pterm.Printf("  Description: %s\n", description)
			pterm.Printf("  Version: %s\n", stackVersion)
			if workflowFile != "" {
				pterm.Printf("  Workflow File: %s\n", workflowFile)
			}
			pterm.Printf("  Status: %s\n", "created")

			pterm.Printf("\n")
			pterm.Printf("Next steps:\n")
			if workflowFile != "" {
				pterm.Printf("  1. Run your workflow: sloth-runner run %s -f %s\n", stackName, workflowFile)
			} else {
				pterm.Printf("  1. Run your workflow: sloth-runner run %s -f <workflow-file>\n", stackName)
			}
			pterm.Printf("  2. View stack details: sloth-runner stack show %s\n", stackName)
			pterm.Printf("  3. List all stacks: sloth-runner stack list\n")

			return nil
		},
	}

	cmd.Flags().StringP("description", "d", "", "Description of the stack")
	cmd.Flags().StringP("workflow-file", "f", "", "Path to the workflow file")
	cmd.Flags().String("version", "1.0.0", "Version of the stack")

	return cmd
}
