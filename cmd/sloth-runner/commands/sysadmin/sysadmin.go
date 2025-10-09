package sysadmin

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/debug"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/health"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/logs"
	"github.com/spf13/cobra"
)

// NewSysadminCmd creates the sysadmin command with all subcommands
func NewSysadminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sysadmin",
		Short: "System administration and operations tools",
		Long: `Comprehensive tools for system administrators to manage, monitor, and troubleshoot
the sloth-runner infrastructure including logs, health checks, diagnostics,
backups, and alerting.`,
		Example: `  # View logs in real-time
  sloth-runner sysadmin logs tail --follow

  # Check system health
  sloth-runner sysadmin health check

  # Monitor agents
  sloth-runner sysadmin health agent --all

  # Export logs for analysis
  sloth-runner sysadmin logs export --format json --output logs.json`,
	}

	// Add subcommands
	cmd.AddCommand(logs.NewLogsCmd())
	cmd.AddCommand(health.NewHealthCmd())
	cmd.AddCommand(debug.NewDebugCmd())

	return cmd
}
