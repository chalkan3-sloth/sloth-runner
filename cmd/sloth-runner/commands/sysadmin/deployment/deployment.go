package deployment

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewDeploymentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployment",
		Short: "Deployment and rollback management",
		Long:  `Manage deployments, rollbacks, and version control across agents.`,
		Aliases: []string{"deploy"},
		Example: `  # Deploy new version
  sloth-runner sysadmin deployment deploy --version v6.25.0 --agents all

  # Rollback to previous version
  sloth-runner sysadmin deployment rollback --agent do-sloth-runner-01

  # Show deployment history
  sloth-runner sysadmin deployment history --agent do-sloth-runner-01`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "deploy",
		Short: "Deploy new version to agents",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Deployment not yet implemented")
			pterm.Info.Println("Future features: Rolling updates, canary deployments, blue-green")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "rollback",
		Short: "Rollback to previous version",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Rollback not yet implemented")
			pterm.Info.Println("Future features: One-click rollback, version history, safety checks")
		},
	})

	return cmd
}
