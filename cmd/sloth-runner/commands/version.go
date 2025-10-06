package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates the version command
func NewVersionCommand(ctx *AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("sloth-runner version %s\n", ctx.Version)
			fmt.Printf("Git commit: %s\n", ctx.Commit)
			fmt.Printf("Build date: %s\n", ctx.Date)
		},
	}
}
