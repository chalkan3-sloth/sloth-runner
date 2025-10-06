package events

import (
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List events in the queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := hooks.NewRepository()
			if err != nil {
				return err
			}
			defer repo.Close()

			status, _ := cmd.Flags().GetString("status")
			eventType, _ := cmd.Flags().GetString("type")
			limit, _ := cmd.Flags().GetInt("limit")

			var statusFilter hooks.EventStatus
			if status != "" {
				statusFilter = hooks.EventStatus(status)
			}

			var typeFilter hooks.EventType
			if eventType != "" {
				typeFilter = hooks.EventType(eventType)
			}

			events, err := repo.EventQueue.ListEvents(typeFilter, statusFilter, limit)
			if err != nil {
				return err
			}

			pterm.DefaultSection.Println("Event Queue")
			data := pterm.TableData{
				{"ID", "Type", "Status", "Created", "Processed"},
			}

			for _, event := range events {
				processedAt := "N/A"
				if event.ProcessedAt != nil {
					processedAt = event.ProcessedAt.Format(time.RFC3339)
				}
				data = append(data, []string{
					event.ID[:8],
					string(event.Type),
					string(event.Status),
					event.CreatedAt.Format(time.RFC3339),
					processedAt,
				})
			}

			pterm.DefaultTable.WithHasHeader().WithData(data).Render()
			pterm.Info.Printf("\nTotal events: %d\n", len(events))

			return nil
		},
	}

	cmd.Flags().StringP("status", "s", "", "Filter by status (pending, processing, completed, failed)")
	cmd.Flags().StringP("type", "t", "", "Filter by event type")
	cmd.Flags().IntP("limit", "l", 50, "Limit number of results")

	return cmd
}
