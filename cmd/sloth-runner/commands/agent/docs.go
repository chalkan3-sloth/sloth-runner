package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/docs"
	"github.com/spf13/cobra"
)

// NewDocsCommand creates the agent docs command
func NewDocsCommand(ctx *commands.AppContext) *cobra.Command {
	var viewMode string

	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Show agent command documentation",
		Long:  `Display comprehensive documentation for the agent command and all its subcommands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			viewer := docs.NewDocViewer(docs.ViewMode(viewMode))
			return viewer.ShowCommand("agent")
		},
	}

	cmd.Flags().StringVarP(&viewMode, "mode", "m", "terminal", "View mode: terminal, raw, browser")

	return cmd
}
