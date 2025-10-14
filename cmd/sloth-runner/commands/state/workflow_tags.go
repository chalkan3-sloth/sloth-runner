//go:build cgo
// +build cgo

package state

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewWorkflowTagsCommand creates the workflow tags management command
func NewWorkflowTagsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage workflow tags",
		Long:  `Add, remove, or list tags for workflow states. Tags help organize and categorize workflows.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewWorkflowTagsAddCommand(ctx),
		NewWorkflowTagsRemoveCommand(ctx),
		NewWorkflowTagsListCommand(ctx),
	)

	return cmd
}

// NewWorkflowTagsAddCommand creates the tag add command
func NewWorkflowTagsAddCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "add <workflow-id> <tag> [tags...]",
		Short: "Add tags to a workflow",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			tags := args[1:]

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			if err := sm.InitWorkflowSchema(); err != nil {
				return err
			}
			if err := sm.ExtendWorkflowSchema(); err != nil {
				return err
			}

			// Verify workflow exists
			_, err = sm.GetWorkflowState(workflowID)
			if err != nil {
				return fmt.Errorf("workflow not found: %s", workflowID)
			}

			// Add tags
			for _, tag := range tags {
				if err := sm.AddTag(workflowID, tag); err != nil {
					pterm.Warning.Printfln("Failed to add tag '%s': %v", tag, err)
				} else {
					pterm.Success.Printfln("Added tag: %s", tag)
				}
			}

			return nil
		},
	}
}

// NewWorkflowTagsRemoveCommand creates the tag remove command
func NewWorkflowTagsRemoveCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <workflow-id> <tag> [tags...]",
		Short: "Remove tags from a workflow",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			tags := args[1:]

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			if err := sm.ExtendWorkflowSchema(); err != nil {
				return err
			}

			for _, tag := range tags {
				if err := sm.RemoveTag(workflowID, tag); err != nil {
					pterm.Warning.Printfln("Failed to remove tag '%s': %v", tag, err)
				} else {
					pterm.Success.Printfln("Removed tag: %s", tag)
				}
			}

			return nil
		},
	}
}

// NewWorkflowTagsListCommand creates the tag list command
func NewWorkflowTagsListCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list <workflow-id>",
		Short: "List tags for a workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]

			sm, err := state.NewStateManager("")
			if err != nil {
				return fmt.Errorf("failed to initialize state manager: %w", err)
			}
			defer sm.Close()

			if err := sm.ExtendWorkflowSchema(); err != nil {
				return err
			}

			tags, err := sm.GetTags(workflowID)
			if err != nil {
				return err
			}

			if len(tags) == 0 {
				pterm.Info.Println("No tags found")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Tags for workflow: %s", workflowID)
			fmt.Println()

			for _, tag := range tags {
				fmt.Printf("  • %s\n", pterm.Cyan(tag))
			}

			fmt.Println()
			pterm.Success.Printf("Total: %d tag(s)\n", len(tags))

			return nil
		},
	}
}
