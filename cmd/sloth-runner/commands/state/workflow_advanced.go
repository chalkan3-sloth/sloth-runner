//go:build cgo
// +build cgo

package state

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewWorkflowImportCommand creates the import command
func NewWorkflowImportCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import workflow state from JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			overwrite, _ := cmd.Flags().GetBool("overwrite")

			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			var export state.WorkflowStateExport
			if err := json.Unmarshal(data, &export); err != nil {
				return err
			}

			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()
			sm.ExtendWorkflowSchema()

			spinner, _ := pterm.DefaultSpinner.Start("Importing workflow state...")

			if err := sm.ImportWorkflowState(&export, overwrite); err != nil {
				spinner.Fail("Import failed")
				return err
			}

			spinner.Success(fmt.Sprintf("Imported workflow: %s", export.State.Name))
			return nil
		},
	}

	cmd.Flags().BoolP("overwrite", "f", false, "Overwrite if workflow exists")
	return cmd
}

// NewWorkflowExportCommand creates the export command
func NewWorkflowExportCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "export <workflow-id> <output-file>",
		Short: "Export workflow state to JSON",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			outputFile := args[1]

			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()
			sm.ExtendWorkflowSchema()

			export, err := sm.ExportWorkflowState(workflowID, os.Getenv("USER"))
			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(export, "", "  ")
			if err != nil {
				return err
			}

			if err := os.WriteFile(outputFile, data, 0644); err != nil {
				return err
			}

			pterm.Success.Printfln("Exported workflow to: %s", outputFile)
			return nil
		},
	}
}

// NewWorkflowBackupCommand creates the backup command
func NewWorkflowBackupCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup <workflow-id>",
		Short: "Create compressed backup of workflow state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			backupDir, _ := cmd.Flags().GetString("dir")

			if backupDir == "" {
				home, _ := os.UserHomeDir()
				backupDir = home + "/.sloth-runner/backups"
			}

			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()
			sm.ExtendWorkflowSchema()

			spinner, _ := pterm.DefaultSpinner.Start("Creating backup...")

			backupPath, err := sm.BackupWorkflowState(workflowID, backupDir, os.Getenv("USER"))
			if err != nil {
				spinner.Fail("Backup failed")
				return err
			}

			spinner.Success(fmt.Sprintf("Backup created: %s", backupPath))
			return nil
		},
	}

	cmd.Flags().StringP("dir", "d", "", "Backup directory (default: ~/.sloth-runner/backups)")
	return cmd
}

// NewWorkflowRestoreCommand creates the restore command
func NewWorkflowRestoreCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore <backup-file>",
		Short: "Restore workflow state from backup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			backupFile := args[0]
			overwrite, _ := cmd.Flags().GetBool("overwrite")

			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()
			sm.ExtendWorkflowSchema()

			spinner, _ := pterm.DefaultSpinner.Start("Restoring from backup...")

			if err := sm.RestoreWorkflowState(backupFile, overwrite); err != nil {
				spinner.Fail("Restore failed")
				return err
			}

			spinner.Success("Workflow restored successfully")
			return nil
		},
	}

	cmd.Flags().BoolP("overwrite", "f", false, "Overwrite if workflow exists")
	return cmd
}

// NewWorkflowDiffCommand creates the diff command
func NewWorkflowDiffCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "diff <workflow-id> <from-version> <to-version>",
		Short: "Show differences between two versions",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			var fromVer, toVer int
			fmt.Sscanf(args[1], "%d", &fromVer)
			fmt.Sscanf(args[2], "%d", &toVer)

			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()
			sm.ExtendWorkflowSchema()

			diff, err := sm.DiffVersions(workflowID, fromVer, toVer)
			if err != nil {
				return err
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("State Diff: v%d → v%d", fromVer, toVer)
			fmt.Println()

			if diff.StatusChange != "" {
				fmt.Printf("Status: %s\n", pterm.Yellow(diff.StatusChange))
			}

			if len(diff.AddedResources) > 0 {
				pterm.Success.Printfln("\n+ Added Resources (%d)", len(diff.AddedResources))
				for _, r := range diff.AddedResources {
					fmt.Printf("  + %s (%s)\n", pterm.Green(r.Name), r.Type)
				}
			}

			if len(diff.RemovedResources) > 0 {
				pterm.Error.Printfln("\n- Removed Resources (%d)", len(diff.RemovedResources))
				for _, r := range diff.RemovedResources {
					fmt.Printf("  - %s (%s)\n", pterm.Red(r.Name), r.Type)
				}
			}

			if len(diff.ModifiedResources) > 0 {
				pterm.Info.Printfln("\n~ Modified Resources (%d)", len(diff.ModifiedResources))
				for _, r := range diff.ModifiedResources {
					fmt.Printf("  ~ %s (%s)\n", pterm.Yellow(r.Name), r.Type)
				}
			}

			return nil
		},
	}
}

// NewWorkflowSearchCommand creates the search command
func NewWorkflowSearchCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Advanced search for workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			query := state.StateSearchQuery{}
			query.Name, _ = cmd.Flags().GetString("name")
			query.ResourceType, _ = cmd.Flags().GetString("resource-type")
			query.HasErrors, _ = cmd.Flags().GetBool("has-errors")
			query.Limit, _ = cmd.Flags().GetInt("limit")

			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()
			sm.ExtendWorkflowSchema()

			workflows, err := sm.SearchWorkflows(query)
			if err != nil {
				return err
			}

			if len(workflows) == 0 {
				pterm.Info.Println("No workflows found matching criteria")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Println("Search Results")
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tSTATUS\tRESOURCES\tSTARTED")
			fmt.Fprintln(w, "----\t------\t---------\t-------")

			for _, workflow := range workflows {
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
					pterm.Cyan(workflow.Name),
					string(workflow.Status),
					len(workflow.Resources),
					workflow.StartedAt.Format("2006-01-02 15:04"),
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Found: %d workflow(s)\n", len(workflows))

			return nil
		},
	}

	cmd.Flags().String("name", "", "Filter by workflow name (partial match)")
	cmd.Flags().String("resource-type", "", "Filter by resource type")
	cmd.Flags().Bool("has-errors", false, "Show only workflows with errors")
	cmd.Flags().Int("limit", 50, "Maximum results")

	return cmd
}

// NewWorkflowPruneCommand creates the prune command
func NewWorkflowPruneCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove old workflow states",
		RunE: func(cmd *cobra.Command, args []string) error {
			days, _ := cmd.Flags().GetInt("older-than")
			keepSuccessful, _ := cmd.Flags().GetBool("keep-successful")
			force, _ := cmd.Flags().GetBool("force")

			olderThan := time.Duration(days) * 24 * time.Hour

			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()

			if !force {
				pterm.Warning.Printfln("About to delete workflow states older than %d days", days)
				if keepSuccessful {
					pterm.Info.Println("Successful workflows will be kept")
				}
				confirm, _ := pterm.DefaultInteractiveConfirm.Show("Continue?")
				if !confirm {
					return nil
				}
			}

			count, err := sm.PruneOldStates(olderThan, keepSuccessful)
			if err != nil {
				return err
			}

			pterm.Success.Printfln("Pruned %d workflow state(s)", count)
			return nil
		},
	}

	cmd.Flags().Int("older-than", 30, "Delete states older than N days")
	cmd.Flags().Bool("keep-successful", false, "Keep successful workflows")
	cmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	return cmd
}

// NewWorkflowAnalyticsCommand creates the analytics command
func NewWorkflowAnalyticsCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "analytics",
		Short: "Show workflow analytics and statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			sm, err := state.NewStateManager("")
			if err != nil {
				return err
			}
			defer sm.Close()

			sm.InitWorkflowSchema()
			sm.ExtendWorkflowSchema()

			analytics, err := sm.GetAnalytics()
			if err != nil {
				return err
			}

			pterm.DefaultHeader.WithFullWidth().Println("Workflow Analytics")
			fmt.Println()

			// Overview
			pterm.DefaultSection.Println("Overview")
			fmt.Printf("Total Workflows:  %d\n", analytics.TotalWorkflows)
			fmt.Printf("Total Resources:  %d\n", analytics.TotalResources)
			fmt.Printf("Success Rate:     %.1f%%\n", analytics.SuccessRate)
			fmt.Printf("Avg Duration:     %.0fs\n", analytics.AverageDuration)
			fmt.Println()

			// Status distribution
			if len(analytics.StatusDistribution) > 0 {
				pterm.DefaultSection.Println("Status Distribution")
				for status, count := range analytics.StatusDistribution {
					fmt.Printf("%s: %d\n", status, count)
				}
				fmt.Println()
			}

			// Resource types
			if len(analytics.ResourceTypes) > 0 {
				pterm.DefaultSection.Println("Resource Types")
				for rtype, count := range analytics.ResourceTypes {
					fmt.Printf("%s: %d\n", rtype, count)
				}
				fmt.Println()
			}

			// Top workflows
			if len(analytics.TopWorkflows) > 0 {
				pterm.DefaultSection.Println("Top Workflows")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(w, "WORKFLOW\tEXECS\tSUCCESS\tFAILURE\tAVG DURATION")
				fmt.Fprintln(w, "--------\t-----\t-------\t-------\t------------")

				for _, stats := range analytics.TopWorkflows {
					fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%.0fs\n",
						pterm.Cyan(stats.WorkflowName),
						stats.ExecutionCount,
						stats.SuccessCount,
						stats.FailureCount,
						stats.AverageDuration,
					)
				}

				w.Flush()
			}

			return nil
		},
	}
}
