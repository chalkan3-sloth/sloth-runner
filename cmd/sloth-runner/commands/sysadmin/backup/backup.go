package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup and restore sloth-runner data",
		Long:  `Create backups of sloth-runner databases, configurations, and state files.`,
	}

	// create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new backup",
		Run: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.Flags().GetString("output")
			compress, _ := cmd.Flags().GetBool("compress")
			description, _ := cmd.Flags().GetString("description")
			include, _ := cmd.Flags().GetStringSlice("include")

			// Se nenhum path especificado, usa defaults
			if len(include) == 0 {
				include = GetDefaultBackupPaths()
			}

			// Gera nome de arquivo se não especificado
			if output == "" {
				timestamp := time.Now().Format("20060102-150405")
				output = filepath.Join(GetDefaultBackupDir(), fmt.Sprintf("sloth-runner-backup-%s.tar.gz", timestamp))
				compress = true
			}

			options := BackupOptions{
				OutputPath:  output,
				Include:     include,
				Compress:    compress,
				Description: description,
			}

			if err := runCreateBackup(options); err != nil {
				pterm.Error.Printf("Failed to create backup: %v\n", err)
			}
		},
	}
	createCmd.Flags().StringP("output", "o", "", "Output backup file path")
	createCmd.Flags().Bool("compress", true, "Compress backup with gzip")
	createCmd.Flags().StringP("description", "d", "", "Backup description")
	createCmd.Flags().StringSliceP("include", "i", []string{}, "Paths to include (defaults to sloth-runner data)")
	cmd.AddCommand(createCmd)

	// restore command
	restoreCmd := &cobra.Command{
		Use:   "restore [backup-file]",
		Short: "Restore from backup",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			targetDir, _ := cmd.Flags().GetString("target")
			databaseOnly, _ := cmd.Flags().GetBool("database-only")
			configOnly, _ := cmd.Flags().GetBool("config-only")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			// Target padrão é o home do usuário
			if targetDir == "" {
				targetDir = os.Getenv("HOME")
			}

			options := RestoreOptions{
				InputPath:    input,
				TargetDir:    targetDir,
				DatabaseOnly: databaseOnly,
				ConfigOnly:   configOnly,
				DryRun:       dryRun,
			}

			if err := runRestoreBackup(options); err != nil {
				pterm.Error.Printf("Failed to restore backup: %v\n", err)
			}
		},
	}
	restoreCmd.Flags().StringP("target", "t", "", "Target directory for restore (default: $HOME)")
	restoreCmd.Flags().Bool("database-only", false, "Restore only database files")
	restoreCmd.Flags().Bool("config-only", false, "Restore only configuration files")
	restoreCmd.Flags().Bool("dry-run", false, "Show what would be restored without actually restoring")
	cmd.AddCommand(restoreCmd)

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available backups",
		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := cmd.Flags().GetString("dir")

			// Se não especificado, usa diretório padrão
			if dir == "" {
				dir = GetDefaultBackupDir()
			}

			if err := runListBackups(dir); err != nil {
				pterm.Error.Printf("Failed to list backups: %v\n", err)
			}
		},
	}
	listCmd.Flags().StringP("dir", "d", "", "Backup directory to list (default: ~/.local/share/sloth-runner/backups)")
	cmd.AddCommand(listCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the backup command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showBackupDocs()
		},
	})

	return cmd
}

func showBackupDocs() {
	title := "SLOTH-RUNNER SYSADMIN BACKUP(1)"
	description := "sloth-runner sysadmin backup - Backup and restore sloth-runner data"
	synopsis := "sloth-runner sysadmin backup [subcommand] [options]"

	options := [][]string{
		{"create", "Create a new backup with full or incremental mode, compression, and encryption support."},
		{"restore", "Restore from backup with point-in-time recovery and selective restore capabilities."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Create full backup",
			"sloth-runner sysadmin backup create --output backup.tar.gz",
			"Creates a full backup compressed with gzip",
		},
		{
			"Create encrypted backup",
			"sloth-runner sysadmin backup create --output backup.tar.gz.enc --encrypt",
			"Creates an encrypted backup with password protection",
		},
		{
			"Incremental backup",
			"sloth-runner sysadmin backup create --output backup-inc.tar.gz --incremental",
			"Creates incremental backup with only changed files",
		},
		{
			"Restore from backup",
			"sloth-runner sysadmin backup restore --input backup.tar.gz",
			"Restores all data from the backup file",
		},
		{
			"Selective restore",
			"sloth-runner sysadmin backup restore --input backup.tar.gz --database-only",
			"Restores only the database, skipping configuration files",
		},
		{
			"Point-in-time recovery",
			"sloth-runner sysadmin backup restore --input backup.tar.gz --timestamp 2025-01-01T12:00:00Z",
			"Restores to a specific point in time",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin config - Configuration management",
		"sloth-runner sysadmin maintenance - System maintenance",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin backup --help")
}

// runCreateBackup cria um backup
func runCreateBackup(options BackupOptions) error {
	manager := NewBackupManager()

	// Cria diretório de backup se não existir
	backupDir := filepath.Dir(options.OutputPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	spinner, _ := pterm.DefaultSpinner.Start("Creating backup...")

	start := time.Now()
	info, err := manager.CreateBackup(options)
	duration := time.Since(start)

	if err != nil {
		spinner.Fail("Failed to create backup")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Backup created: %s", filepath.Base(options.OutputPath)))
	pterm.Println()

	pterm.DefaultHeader.WithFullWidth().Println("Backup Created Successfully")
	pterm.Println()

	// Results table
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Output File", options.OutputPath},
		{"Size", FormatBytes(info.Size)},
		{"Files Backed Up", fmt.Sprintf("%d", info.FileCount)},
		{"Compressed", fmt.Sprintf("%v", info.Compressed)},
		{"Duration", duration.String()},
	}

	if options.Description != "" {
		tableData = append(tableData, []string{"Description", options.Description})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Success.Println("Backup completed successfully!")
	pterm.Info.Printf("To restore: sloth-runner sysadmin backup restore %s\n", options.OutputPath)

	return nil
}

// runRestoreBackup restaura um backup
func runRestoreBackup(options RestoreOptions) error {
	manager := NewBackupManager()

	mode := "Restoring"
	if options.DryRun {
		mode = "Analyzing"
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("%s backup...", mode))

	info, err := manager.RestoreBackup(options)
	if err != nil {
		spinner.Fail("Failed to restore backup")
		return err
	}

	if options.DryRun {
		spinner.Success(fmt.Sprintf("✅ Would restore %d file(s)", info.FilesRestored))
	} else {
		spinner.Success(fmt.Sprintf("✅ Restored %d file(s)", info.FilesRestored))
	}

	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println("Backup Restore Results")
	pterm.Println()

	// Results table
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Files Restored", fmt.Sprintf("%d", info.FilesRestored)},
		{"Bytes Restored", FormatBytes(info.BytesRestored)},
		{"Duration", info.Duration.String()},
		{"Target Directory", options.TargetDir},
	}

	// Filter info
	filters := []string{}
	if options.DatabaseOnly {
		filters = append(filters, "Database Only")
	}
	if options.ConfigOnly {
		filters = append(filters, "Config Only")
	}
	if len(filters) > 0 {
		tableData = append(tableData, []string{"Filters", fmt.Sprintf("%v", filters)})
	}

	tableData = append(tableData, []string{"Mode", func() string {
		if options.DryRun {
			return "Dry Run (no changes made)"
		}
		return "Live"
	}()})

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Show errors if any
	if len(info.Errors) > 0 {
		pterm.Warning.Printf("Encountered %d error(s) during restore:\n", len(info.Errors))
		for i, err := range info.Errors {
			if i < 5 { // Show only first 5 errors
				pterm.Error.Println("  " + err)
			}
		}
		if len(info.Errors) > 5 {
			pterm.Info.Printf("  ... and %d more errors\n", len(info.Errors)-5)
		}
		pterm.Println()
	}

	if options.DryRun {
		pterm.Warning.Println("This was a dry run. Use without --dry-run to actually restore files.")
	} else {
		pterm.Success.Println("Restore completed successfully!")
	}

	return nil
}

// runListBackups lista backups disponíveis
func runListBackups(dir string) error {
	manager := NewBackupManager()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Scanning %s for backups...", dir))

	backups, err := manager.ListBackups(dir)
	if err != nil {
		spinner.Fail("Failed to list backups")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d backup(s)", len(backups)))
	pterm.Println()

	if len(backups) == 0 {
		pterm.Info.Println("No backups found in directory.")
		pterm.Info.Printf("Create a backup with: sloth-runner sysadmin backup create\n")
		return nil
	}

	pterm.DefaultHeader.WithFullWidth().Println("Available Backups")
	pterm.Println()

	// Backups table
	tableData := pterm.TableData{
		{"Filename", "Size", "Created", "Compressed"},
	}

	for _, backup := range backups {
		tableData = append(tableData, []string{
			filepath.Base(backup.Path),
			FormatBytes(backup.Size),
			backup.Created.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%v", backup.Compressed),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Info.Printf("Backup directory: %s\n", dir)

	return nil
}
