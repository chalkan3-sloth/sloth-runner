//go:build !cgo
// +build !cgo

package hook

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewHookCommand creates a stub hook command for non-CGO builds
func NewHookCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Manage workflow hooks (requires CGO)",
		Long: `The hook command requires CGO support (SQLite) which is not available in this build.

Hook management features are only available in CGO-enabled builds.
This typically means Linux amd64/arm64 binaries in our releases.

To use hook management:
  1. Download the CGO-enabled version for Linux from GitHub releases
  2. Or compile from source with CGO_ENABLED=1

This binary can still be used for basic workflow execution and agent management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("hook management requires CGO support (SQLite).\n" +
				"This binary was compiled without CGO for portability.\n\n" +
				"Please use a CGO-enabled build (available for Linux) or compile with CGO_ENABLED=1")
		},
	}

	return cmd
}
