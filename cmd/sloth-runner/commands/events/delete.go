package events

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <event-id>",
		Short: "Delete an event from the queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := hooks.NewRepository()
			if err != nil {
				return err
			}
			defer repo.Close()

			if err := repo.EventQueue.DeleteEvent(args[0]); err != nil {
				return fmt.Errorf("failed to delete event: %w", err)
			}

			pterm.Success.Printf("Event %s deleted\n", args[0])
			return nil
		},
	}
}
