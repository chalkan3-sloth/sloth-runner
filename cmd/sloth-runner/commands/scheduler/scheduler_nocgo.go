//go:build !cgo
// +build !cgo

package scheduler

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewSchedulerCommand creates a stub scheduler command for non-CGO builds
func NewSchedulerCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "Manage workflow scheduler (requires CGO)",
		Long: `The scheduler command requires CGO support (SQLite) which is not available in this build.

Scheduler features are only available in CGO-enabled builds.
This typically means Linux amd64/arm64 binaries in our releases.

To use the scheduler:
  1. Download the CGO-enabled version for Linux from GitHub releases
  2. Or compile from source with CGO_ENABLED=1

This binary can still be used for basic workflow execution and agent management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("scheduler requires CGO support (SQLite).\n" +
				"This binary was compiled without CGO for portability.\n\n" +
				"Please use a CGO-enabled build (available for Linux) or compile with CGO_ENABLED=1")
		},
	}

	return cmd
}
