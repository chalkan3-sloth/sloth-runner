package events

import (
	"encoding/json"
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/spf13/cobra"
)

func NewGetCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "get <event-id>",
		Short: "Get event details in JSON format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := hooks.NewRepository()
			if err != nil {
				return err
			}
			defer repo.Close()

			event, err := repo.EventQueue.GetEvent(args[0])
			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}
}
