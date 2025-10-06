package commands

import (
	"github.com/spf13/cobra"
)

// NewListCommand creates the list command
// TODO: Extract logic from main.go listCmd (lines ~574-641)
func NewListCommand(ctx *AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available workflows and tasks",
		Long:  `List all available workflows and tasks from sloth files in the current directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			// - Scan directory for .sloth files
			// - Parse and display available workflows
			return nil
		},
	}
}
