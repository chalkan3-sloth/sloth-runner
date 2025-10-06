package hook

import (
	"encoding/json"
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/spf13/cobra"
)

// NewGetCommand creates the hook get command
func NewGetCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <hook-name>",
		Short: "Get hook details in JSON format",
		Long:  `Get detailed information about a specific hook in JSON format.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]

			// Create repository
			repo, err := hooks.NewRepository()
			if err != nil {
				return fmt.Errorf("failed to create repository: %w", err)
			}
			defer repo.Close()

			// Get hook
			hook, err := repo.GetByName(hookName)
			if err != nil {
				return fmt.Errorf("failed to get hook: %w", err)
			}

			// Marshal to JSON
			data, err := json.MarshalIndent(hook, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal hook: %w", err)
			}

			fmt.Println(string(data))

			return nil
		},
	}

	return cmd
}
