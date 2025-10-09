package sysadmin

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/debug"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/health"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/logs"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sysadmin/backup"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sysadmin/deployment"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sysadmin/maintenance"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sysadmin/network"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sysadmin/performance"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sysadmin/security"
	"github.com/spf13/cobra"
)

// NewSysadminCmd creates the sysadmin command with all subcommands
func NewSysadminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sysadmin",
		Short: "System administration and operations tools",
		Long: `Comprehensive tools for system administrators to manage, monitor, and troubleshoot
the sloth-runner infrastructure including logs, health checks, diagnostics,
backups, performance monitoring, network diagnostics, and security auditing.`,
		Example: `  # View logs in real-time
  sloth-runner sysadmin logs tail --follow

  # Check system health
  sloth-runner sysadmin health check

  # Create backup
  sloth-runner sysadmin backup create --output backup.tar.gz

  # Monitor performance
  sloth-runner sysadmin performance show --agent do-sloth-runner-01

  # Network diagnostics
  sloth-runner sysadmin network ping --agent do-sloth-runner-01`,
	}

	// Add core subcommands
	cmd.AddCommand(logs.NewLogsCmd())
	cmd.AddCommand(health.NewHealthCmd())
	cmd.AddCommand(debug.NewDebugCmd())

	// Add new sysadmin tools
	cmd.AddCommand(backup.NewBackupCmd())
	cmd.AddCommand(performance.NewPerformanceCmd())
	cmd.AddCommand(network.NewNetworkCmd())
	cmd.AddCommand(security.NewSecurityCmd())
	cmd.AddCommand(maintenance.NewMaintenanceCmd())
	cmd.AddCommand(deployment.NewDeploymentCmd())

	return cmd
}
