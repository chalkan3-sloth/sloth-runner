//go:build !cgo
// +build !cgo

package sloth

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewSlothCommand creates a stub sloth command for non-CGO builds
func NewSlothCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sloth",
		Short: "Manage sloth files (requires CGO)",
		Long: `The sloth command requires CGO support (SQLite) which is not available in this build.

Sloth file management features are only available in CGO-enabled builds.
This typically means Linux amd64/arm64 binaries in our releases.

To use sloth file management:
  1. Download the CGO-enabled version for Linux from GitHub releases
  2. Or compile from source with CGO_ENABLED=1

This binary can still be used for basic workflow execution and agent management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("sloth file management requires CGO support (SQLite).\n" +
				"This binary was compiled without CGO for portability.\n\n" +
				"Please use a CGO-enabled build (available for Linux) or compile with CGO_ENABLED=1")
		},
	}

	return cmd
}
