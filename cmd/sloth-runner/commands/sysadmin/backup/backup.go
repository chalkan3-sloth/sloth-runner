package backup

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup and restore sloth-runner data",
		Long:  `Create backups of sloth-runner databases, configurations, and state files.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new backup",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Backup creation not yet implemented")
			pterm.Info.Println("Future features: Full/incremental backups, compression, encryption")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "restore",
		Short: "Restore from backup",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Backup restore not yet implemented")
			pterm.Info.Println("Future features: Point-in-time recovery, selective restore")
		},
	})

	return cmd
}
