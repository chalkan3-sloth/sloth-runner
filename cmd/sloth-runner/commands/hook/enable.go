package hook

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewEnableCommand creates the hook enable command
func NewEnableCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable <hook-name>",
		Short: "Enable a hook",
		Long:  `Enable an event hook so it will be triggered when events occur.`,
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

			if hook.Enabled {
				pterm.Info.Printf("Hook '%s' is already enabled\n", hookName)
				return nil
			}

			// Enable hook
			if err := repo.Enable(hook.ID); err != nil {
				return fmt.Errorf("failed to enable hook: %w", err)
			}

			pterm.Success.Printf("Hook '%s' enabled successfully\n", hookName)

			return nil
		},
	}

	return cmd
}
