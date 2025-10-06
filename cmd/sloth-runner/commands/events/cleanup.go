package events

import (
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewCleanupCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup old events",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := hooks.NewRepository()
			if err != nil {
				return err
			}
			defer repo.Close()

			hours, _ := cmd.Flags().GetInt("hours")
			removed, err := repo.EventQueue.CleanupOldEvents(time.Duration(hours) * time.Hour)
			if err != nil {
				return fmt.Errorf("failed to cleanup events: %w", err)
			}

			pterm.Success.Printf("Cleaned up %d old events\n", removed)
			return nil
		},
	}

	cmd.Flags().IntP("hours", "H", 24, "Remove events older than this many hours")
	return cmd
}
