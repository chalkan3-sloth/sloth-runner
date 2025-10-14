//go:build cgo
// +build cgo

package stack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewMigrateCommand creates the migration command
func NewMigrateCommand(ctx *commands.AppContext) *cobra.Command {
	var sourceDB, targetDB string
	var dryRun bool
	var outputScript string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate from old workflow_state to unified state backend",
		Long: `Migrates data from the old workflow_state system to the new unified state backend.
This consolidates workflow states, resources, and outputs into a single Pulumi/Terraform-like state management system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate SQL script if requested
			if outputScript != "" {
				script := stack.GenerateMigrationScript()
				if err := os.WriteFile(outputScript, []byte(script), 0644); err != nil {
					return fmt.Errorf("failed to write migration script: %w", err)
				}
				pterm.Success.Printfln("Migration script written to: %s", outputScript)
				return nil
			}

			// Validate source database
			if sourceDB == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				sourceDB = filepath.Join(homeDir, ".sloth-runner", "state.db")
			}

			// Validate target database
			if targetDB == "" {
				targetDB = "/etc/sloth-runner/stacks.db"
			}

			// Check if source exists
			if _, err := os.Stat(sourceDB); os.IsNotExist(err) {
				return fmt.Errorf("source database not found: %s", sourceDB)
			}

			if dryRun {
				pterm.Info.Println("Dry run mode - no data will be modified")
				pterm.Info.Printfln("Source: %s", sourceDB)
				pterm.Info.Printfln("Target: %s", targetDB)
				return nil
			}

			// Confirm migration
			pterm.Warning.Println("This will migrate data from the old system to the new unified state backend")
			pterm.Info.Printfln("Source: %s", sourceDB)
			pterm.Info.Printfln("Target: %s", targetDB)
			fmt.Println()

			result, _ := pterm.DefaultInteractiveConfirm.Show("Continue with migration?")
			if !result {
				pterm.Info.Println("Migration cancelled")
				return nil
			}

			// Perform migration
			spinner, _ := pterm.DefaultSpinner.Start("Migrating data...")

			migrator, err := stack.NewMigrator(sourceDB, targetDB)
			if err != nil {
				spinner.Fail("Failed to initialize migrator")
				return fmt.Errorf("failed to initialize migrator: %w", err)
			}
			defer migrator.Close()

			report, err := migrator.PerformMigration()
			if err != nil {
				spinner.Fail("Migration failed")
				return fmt.Errorf("migration failed: %w", err)
			}

			spinner.Success("Migration completed")
			fmt.Println()

			// Display report
			pterm.DefaultHeader.WithFullWidth().Println("Migration Report")
			fmt.Println()

			data := [][]string{
				{"Stacks Migrated", fmt.Sprintf("%d", report.StacksMigrated)},
				{"Resources Migrated", fmt.Sprintf("%d", report.ResourcesMigrated)},
				{"Duration", report.Duration},
			}

			pterm.DefaultTable.WithHasHeader(false).WithData(data).Render()

			if len(report.Errors) > 0 {
				fmt.Println()
				pterm.Warning.Println("Errors encountered:")
				for _, err := range report.Errors {
					fmt.Printf("  - %s\n", err)
				}
			}

			// Export report
			reportJSON, _ := json.MarshalIndent(report, "", "  ")
			reportFile := filepath.Join(filepath.Dir(targetDB), "migration_report.json")
			os.WriteFile(reportFile, reportJSON, 0644)

			fmt.Println()
			pterm.Success.Printfln("Migration report saved to: %s", reportFile)

			return nil
		},
	}

	cmd.Flags().StringVar(&sourceDB, "source", "", "Source database path (default: ~/.sloth-runner/state.db)")
	cmd.Flags().StringVar(&targetDB, "target", "", "Target database path (default: /etc/sloth-runner/stacks.db)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without making changes")
	cmd.Flags().StringVar(&outputScript, "generate-script", "", "Generate SQL migration script to file")

	return cmd
}
