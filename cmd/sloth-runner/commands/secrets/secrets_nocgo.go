//go:build !cgo
// +build !cgo

package secrets

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewSecretsCommand creates a stub secrets command for non-CGO builds
func NewSecretsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage encrypted secrets (requires CGO)",
		Long: `The secrets command requires CGO support (SQLite) which is not available in this build.

Secrets management features are only available in CGO-enabled builds.
This typically means Linux amd64/arm64 binaries in our releases.

To use secrets management:
  1. Download the CGO-enabled version for Linux from GitHub releases
  2. Or compile from source with CGO_ENABLED=1

This binary can still be used for basic workflow execution and agent management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("secrets management requires CGO support (SQLite).\n" +
				"This binary was compiled without CGO for portability.\n\n" +
				"Please use a CGO-enabled build (available for Linux) or compile with CGO_ENABLED=1")
		},
	}

	return cmd
}
