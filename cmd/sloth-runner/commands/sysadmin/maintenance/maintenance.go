package maintenance

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewMaintenanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "System maintenance and cleanup tasks",
		Long:  `Perform maintenance tasks like log rotation, cleanup, garbage collection, and optimization.`,
		Example: `  # Clean old logs
  sloth-runner sysadmin maintenance clean-logs --older-than 30d

  # Optimize databases
  sloth-runner sysadmin maintenance optimize-db

  # Cleanup temp files
  sloth-runner sysadmin maintenance cleanup --agent do-sloth-runner-01`,
	}

	// clean-logs command
	cleanLogsCmd := &cobra.Command{
		Use:   "clean-logs",
		Short: "Clean old log files",
		Run: func(cmd *cobra.Command, args []string) {
			olderThan, _ := cmd.Flags().GetDuration("older-than")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			if err := runCleanLogs(olderThan, dryRun); err != nil {
				pterm.Error.Printf("Failed to clean logs: %v\n", err)
			}
		},
	}
	cleanLogsCmd.Flags().Duration("older-than", 30*24*time.Hour, "Remove logs older than this duration (e.g., 30d, 7d)")
	cleanLogsCmd.Flags().Bool("dry-run", false, "Show what would be removed without actually removing")
	cmd.AddCommand(cleanLogsCmd)

	// optimize-db command
	optimizeDBCmd := &cobra.Command{
		Use:   "optimize-db [database-path]",
		Short: "Optimize and vacuum databases",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dbPath := args[0]
			full, _ := cmd.Flags().GetBool("full")

			if err := runOptimizeDB(dbPath, full); err != nil {
				pterm.Error.Printf("Failed to optimize database: %v\n", err)
			}
		},
	}
	optimizeDBCmd.Flags().Bool("full", false, "Run full optimization including REINDEX")
	cmd.AddCommand(optimizeDBCmd)

	// cleanup command
	cleanupCmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Clean temporary files and caches",
		Run: func(cmd *cobra.Command, args []string) {
			tempFiles, _ := cmd.Flags().GetBool("temp")
			cache, _ := cmd.Flags().GetBool("cache")
			oldLogs, _ := cmd.Flags().GetBool("logs")
			logAge, _ := cmd.Flags().GetDuration("log-age")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			paths, _ := cmd.Flags().GetStringSlice("paths")

			// Se nenhuma flag específica, faz limpeza geral
			if !tempFiles && !cache && !oldLogs {
				tempFiles = true
				cache = true
				oldLogs = true
			}

			options := CleanupOptions{
				TempFiles: tempFiles,
				Cache:     cache,
				OldLogs:   oldLogs,
				LogAge:    logAge,
				DryRun:    dryRun,
				Paths:     paths,
			}

			if err := runCleanup(options); err != nil {
				pterm.Error.Printf("Failed to cleanup: %v\n", err)
			}
		},
	}
	cleanupCmd.Flags().Bool("temp", false, "Clean temporary files")
	cleanupCmd.Flags().Bool("cache", false, "Clean cache files")
	cleanupCmd.Flags().Bool("logs", false, "Clean old log files")
	cleanupCmd.Flags().Duration("log-age", 30*24*time.Hour, "Age of logs to clean")
	cleanupCmd.Flags().Bool("dry-run", false, "Show what would be cleaned without actually cleaning")
	cleanupCmd.Flags().StringSlice("paths", []string{}, "Additional paths to clean")
	cmd.AddCommand(cleanupCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the maintenance command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showMaintenanceDocs()
		},
	})

	return cmd
}

func showMaintenanceDocs() {
	title := "SLOTH-RUNNER SYSADMIN MAINTENANCE(1)"
	description := "sloth-runner sysadmin maintenance - System maintenance and cleanup tasks"
	synopsis := "sloth-runner sysadmin maintenance [subcommand] [options]"

	options := [][]string{
		{"clean-logs", "Clean and rotate old log files. Compress archives and delete old entries."},
		{"optimize-db", "Optimize and vacuum databases with VACUUM, ANALYZE, and index rebuilding."},
		{"cleanup", "Clean temporary files, caches, and detect orphaned files."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Clean old logs",
			"sloth-runner sysadmin maintenance clean-logs --older-than 30d",
			"Removes log files older than 30 days and compresses recent ones",
		},
		{
			"Optimize database",
			"sloth-runner sysadmin maintenance optimize-db --full",
			"Runs full VACUUM, ANALYZE, and rebuilds indexes",
		},
		{
			"General cleanup",
			"sloth-runner sysadmin maintenance cleanup --dry-run",
			"Shows what would be cleaned without actually removing anything",
		},
		{
			"Automated maintenance",
			"sloth-runner sysadmin maintenance cleanup --all-agents",
			"Runs cleanup tasks across all registered agents",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin logs - Log management",
		"sloth-runner sysadmin backup - Backup and restore",
	}

	showDocs(title, description, synopsis, options, examples, seeAlso)
}

// showDocs displays formatted documentation similar to man pages
func showDocs(title, description, synopsis string, options [][]string, examples [][]string, seeAlso []string) {
	// Header
	pterm.DefaultHeader.WithFullWidth().Println(title)
	fmt.Println()

	// Name and Description
	pterm.DefaultSection.Println("NAME")
	fmt.Printf("    %s\n\n", description)

	// Synopsis
	if synopsis != "" {
		pterm.DefaultSection.Println("SYNOPSIS")
		fmt.Printf("    %s\n\n", synopsis)
	}

	// Options
	if len(options) > 0 {
		pterm.DefaultSection.Println("OPTIONS")
		for _, opt := range options {
			if len(opt) >= 2 {
				pterm.FgCyan.Printf("    %s\n", opt[0])
				fmt.Printf("        %s\n\n", opt[1])
			}
		}
	}

	// Examples
	if len(examples) > 0 {
		pterm.DefaultSection.Println("EXAMPLES")
		for i, ex := range examples {
			if len(ex) >= 2 {
				pterm.FgYellow.Printf("    Example %d: %s\n", i+1, ex[0])
				pterm.FgGreen.Printf("    $ %s\n", ex[1])
				if len(ex) >= 3 {
					fmt.Printf("        %s\n", ex[2])
				}
				fmt.Println()
			}
		}
	}

	// See Also
	if len(seeAlso) > 0 {
		pterm.DefaultSection.Println("SEE ALSO")
		for _, item := range seeAlso {
			fmt.Printf("    • %s\n", item)
		}
		fmt.Println()
	}

	// Footer
	pterm.FgGray.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	pterm.FgGray.Println("Documentation generated for sloth-runner sysadmin v2.0")
	pterm.FgGray.Println("For more information: sloth-runner sysadmin maintenance --help")
}

// runCleanLogs executa limpeza de logs
func runCleanLogs(olderThan time.Duration, dryRun bool) error {
	manager := NewMaintenanceManager()

	mode := "Cleaning"
	if dryRun {
		mode = "Analyzing"
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("%s log files older than %v...", mode, olderThan))

	result, err := manager.CleanLogs(olderThan, dryRun)
	if err != nil {
		spinner.Fail("Failed to clean logs")
		return err
	}

	if dryRun {
		spinner.Success(fmt.Sprintf("✅ Would remove %d log file(s), freeing %s", result.FilesRemoved, FormatBytes(result.SpaceFreed)))
	} else {
		spinner.Success(fmt.Sprintf("✅ Removed %d log file(s), freed %s", result.FilesRemoved, FormatBytes(result.SpaceFreed)))
	}

	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println("Log Cleanup Results")
	pterm.Println()

	// Summary table
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Files Removed", fmt.Sprintf("%d", result.FilesRemoved)},
		{"Space Freed", FormatBytes(result.SpaceFreed)},
		{"Duration", result.Duration.String()},
		{"Mode", func() string {
			if dryRun {
				return "Dry Run (no changes made)"
			}
			return "Live"
		}()},
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Show details if any
	if len(result.Details) > 0 && len(result.Details) <= 10 {
		pterm.DefaultSection.Println("Details")
		for _, detail := range result.Details {
			pterm.Info.Println("  " + detail)
		}
		pterm.Println()
	} else if len(result.Details) > 10 {
		pterm.Info.Printf("Processed %d files (showing first 10)\n", len(result.Details))
		for i := 0; i < 10; i++ {
			pterm.Info.Println("  " + result.Details[i])
		}
		pterm.Println()
	}

	return nil
}

// runOptimizeDB executa otimização de banco de dados
func runOptimizeDB(dbPath string, full bool) error {
	manager := NewMaintenanceManager()

	// Get database size before
	sizeBefore, err := GetDatabaseSize(dbPath)
	if err != nil {
		return fmt.Errorf("failed to get database size: %w", err)
	}

	mode := "standard"
	if full {
		mode = "full"
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Optimizing database with %s mode...", mode))

	start := time.Now()
	err = manager.OptimizeDatabase(dbPath, full)
	duration := time.Since(start)

	if err != nil {
		spinner.Fail("Failed to optimize database")
		return err
	}

	// Get database size after
	sizeAfter, err := GetDatabaseSize(dbPath)
	if err != nil {
		sizeAfter = sizeBefore // Fallback
	}

	saved := int64(sizeBefore) - int64(sizeAfter)
	savedPercent := 0.0
	if sizeBefore > 0 {
		savedPercent = float64(saved) / float64(sizeBefore) * 100.0
	}

	spinner.Success("✅ Database optimization completed")
	pterm.Println()

	pterm.DefaultHeader.WithFullWidth().Println("Database Optimization Results")
	pterm.Println()

	// Results table
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Database Path", dbPath},
		{"Mode", mode},
		{"Size Before", FormatBytes(sizeBefore)},
		{"Size After", FormatBytes(sizeAfter)},
	}

	if saved > 0 {
		tableData = append(tableData, []string{"Space Saved", fmt.Sprintf("%s (%.1f%%)", FormatBytes(uint64(saved)), savedPercent)})
	} else {
		tableData = append(tableData, []string{"Space Saved", "0 B (database already optimized)"})
	}

	tableData = append(tableData, []string{"Duration", duration.String()})

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	if full {
		pterm.Success.Println("Full optimization completed: VACUUM, ANALYZE, REINDEX")
	} else {
		pterm.Info.Println("Standard optimization completed: VACUUM, ANALYZE")
		pterm.Info.Println("Use --full for complete optimization including REINDEX")
	}

	return nil
}

// runCleanup executa limpeza geral
func runCleanup(options CleanupOptions) error {
	manager := NewMaintenanceManager()

	mode := "Cleaning"
	if options.DryRun {
		mode = "Analyzing"
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("%s temporary files and caches...", mode))

	result, err := manager.Cleanup(options)
	if err != nil {
		spinner.Fail("Failed to cleanup")
		return err
	}

	if options.DryRun {
		spinner.Success(fmt.Sprintf("✅ Would remove %d file(s), freeing %s", result.FilesRemoved, FormatBytes(result.SpaceFreed)))
	} else {
		spinner.Success(fmt.Sprintf("✅ Removed %d file(s), freed %s", result.FilesRemoved, FormatBytes(result.SpaceFreed)))
	}

	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println("Cleanup Results")
	pterm.Println()

	// Summary table
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Files Removed", fmt.Sprintf("%d", result.FilesRemoved)},
		{"Space Freed", FormatBytes(result.SpaceFreed)},
		{"Duration", result.Duration.String()},
	}

	// Show what was cleaned
	cleaned := []string{}
	if options.TempFiles {
		cleaned = append(cleaned, "Temp Files")
	}
	if options.Cache {
		cleaned = append(cleaned, "Cache")
	}
	if options.OldLogs {
		cleaned = append(cleaned, fmt.Sprintf("Logs (older than %v)", options.LogAge))
	}
	if len(options.Paths) > 0 {
		cleaned = append(cleaned, fmt.Sprintf("Custom Paths (%d)", len(options.Paths)))
	}

	if len(cleaned) > 0 {
		tableData = append(tableData, []string{"Cleaned", fmt.Sprintf("%v", cleaned)})
	}

	tableData = append(tableData, []string{"Mode", func() string {
		if options.DryRun {
			return "Dry Run (no changes made)"
		}
		return "Live"
	}()})

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Show summary
	if result.FilesRemoved == 0 {
		pterm.Info.Println("No files to clean - system is already clean!")
	} else if options.DryRun {
		pterm.Warning.Println("This was a dry run. Use without --dry-run to actually remove files.")
	} else {
		pterm.Success.Println("Cleanup completed successfully!")
	}

	return nil
}
