//go:build !cgo
// +build !cgo

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd is not available in non-CGO builds
var runCmd *cobra.Command

func init() {
	// runCmd is disabled in non-CGO builds
	// The command would require SQLite which needs CGO
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run workflow (requires CGO build)",
		Long: `The 'run' command requires CGO support for SQLite-based state management.

This binary was compiled without CGO for maximum portability.

To use the 'run' command:
  • Download a CGO-enabled build for Linux (amd64 or arm64)
  • Or compile from source with: CGO_ENABLED=1 go build

For simple command execution without state management, use other commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("'run' command requires CGO support (not available in this build)")
		},
	}

	// Note: runCmd is not automatically added to rootCmd in non-CGO builds
	// because it's not functional. Users need CGO-enabled builds for this feature.
}
