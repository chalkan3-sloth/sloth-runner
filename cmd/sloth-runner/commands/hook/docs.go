//go:build cgo
// +build cgo

package hook

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/docs"
	"github.com/spf13/cobra"
)

// NewDocsCommand creates the hook docs command
func NewDocsCommand(ctx *commands.AppContext) *cobra.Command {
	var viewMode string

	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Show hook command documentation",
		Long:  `Display comprehensive documentation for the hook command and all its subcommands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			viewer := docs.NewDocViewer(docs.ViewMode(viewMode))
			return viewer.ShowCommand("hook")
		},
	}

	cmd.Flags().StringVarP(&viewMode, "mode", "m", "terminal", "View mode: terminal, raw, browser")

	return cmd
}
