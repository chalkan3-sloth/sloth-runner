package config

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management for sloth-runner",
		Long: `Manage sloth-runner configuration files, validate settings, compare configurations
across agents, and dynamically update configuration without restarts.`,
		Example: `  # Validate configuration file
  sloth-runner sysadmin config validate

  # Compare configurations between agents
  sloth-runner sysadmin config diff --agents do-sloth-runner-01,do-sloth-runner-02

  # Export current configuration
  sloth-runner sysadmin config export --output config.yaml

  # Set configuration value
  sloth-runner sysadmin config set --key log.level --value debug

  # Get configuration value
  sloth-runner sysadmin config get --key log.level

  # Reset configuration to defaults
  sloth-runner sysadmin config reset --confirm`,
	}

	// validate subcommand
	validateCmd := &cobra.Command{
		Use:   "validate [config-file]",
		Short: "Validate configuration files",
		Long:  `Validate sloth-runner configuration files for syntax errors and invalid settings.`,
		Run: func(cmd *cobra.Command, args []string) {
			var configPath string

			// Se path fornecido, usa ele; senão, procura nos defaults
			if len(args) > 0 {
				configPath = args[0]
			} else {
				path, err := FindConfigFile()
				if err != nil {
					pterm.Error.Printf("No configuration file found. Specify path or create config in default location.\n")
					pterm.Info.Println("Default locations:")
					for _, p := range GetDefaultConfigPaths() {
						pterm.Info.Printf("  - %s\n", p)
					}
					return
				}
				configPath = path
			}

			if err := runValidateConfig(configPath); err != nil {
				pterm.Error.Printf("Validation failed: %v\n", err)
			}
		},
	}
	cmd.AddCommand(validateCmd)

	// diff subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "diff",
		Short: "Compare configurations between agents",
		Long:  `Compare configuration settings across multiple agents to identify differences.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Configuration diff not yet implemented")
			pterm.Info.Println("Future features: Side-by-side comparison, highlight differences, export diff report")
		},
	})

	// export subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Export current configuration",
		Long:  `Export the current configuration to a file (YAML, JSON, or TOML format).`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Configuration export not yet implemented")
			pterm.Info.Println("Future features: YAML/JSON/TOML formats, include secrets option, template export")
		},
	})

	// import subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "import",
		Short: "Import configuration from file",
		Long:  `Import configuration from a YAML, JSON, or TOML file and apply settings.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Configuration import not yet implemented")
			pterm.Info.Println("Future features: Merge or replace, dry-run mode, backup before import")
		},
	})

	// set subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "Set configuration value dynamically",
		Long:  `Set a configuration value at runtime without restarting the service.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Dynamic configuration set not yet implemented")
			pterm.Info.Println("Future features: Hot reload, validation before apply, rollback on error")
		},
	})

	// get subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get configuration value",
		Long:  `Retrieve and display current configuration values.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Configuration get not yet implemented")
			pterm.Info.Println("Future features: Key path search, JSON/YAML output, show source (file/env/default)")
		},
	})

	// reset subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Long:  `Reset configuration to default values. Use with caution!`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Configuration reset not yet implemented")
			pterm.Info.Println("Future features: Selective reset, backup before reset, confirmation prompt")
		},
	})

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the config command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showConfigDocs()
		},
	})

	return cmd
}

func showConfigDocs() {
	title := "SLOTH-RUNNER SYSADMIN CONFIG(1)"
	description := "sloth-runner sysadmin config - Configuration management for sloth-runner"
	synopsis := "sloth-runner sysadmin config [subcommand] [options]"

	options := [][]string{
		{"validate", "Validate configuration files for syntax errors and invalid settings."},
		{"diff", "Compare configuration settings across multiple agents to identify differences."},
		{"export", "Export the current configuration to a file (YAML, JSON, or TOML format)."},
		{"import", "Import configuration from a YAML, JSON, or TOML file and apply settings."},
		{"set", "Set a configuration value at runtime without restarting the service."},
		{"get", "Retrieve and display current configuration values."},
		{"reset", "Reset configuration to default values. Use with caution!"},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Validate configuration",
			"sloth-runner sysadmin config validate",
			"Checks YAML/JSON syntax and validates all settings",
		},
		{
			"Compare between agents",
			"sloth-runner sysadmin config diff --agents web-01,web-02",
			"Shows side-by-side comparison of configurations",
		},
		{
			"Set value dynamically",
			"sloth-runner sysadmin config set --key log.level --value debug",
			"Updates log level without restart (hot reload)",
		},
		{
			"Export to file",
			"sloth-runner sysadmin config export --output config.yaml",
			"Exports current configuration in YAML format",
		},
		{
			"Import from file",
			"sloth-runner sysadmin config import --input config.yaml --dry-run",
			"Shows what changes would be applied without making them",
		},
		{
			"Get configuration value",
			"sloth-runner sysadmin config get --key database.host",
			"Displays specific configuration value",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin health - Health checks",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin config --help")
}

// runValidateConfig valida arquivo de configuração
func runValidateConfig(configPath string) error {
	validator := NewValidator()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Validating %s...", configPath))

	result, err := validator.ValidateFile(configPath)
	if err != nil {
		spinner.Fail("Validation failed")
		return err
	}

	if result.Valid {
		spinner.Success("✅ Configuration is valid!")
	} else {
		spinner.Fail("❌ Configuration has errors")
	}

	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println("Configuration Validation Results")
	pterm.Println()

	// Summary
	errorCount := len(result.Errors)
	warningCount := len(result.Warnings)

	summaryData := pterm.TableData{
		{"Metric", "Value"},
		{"File", configPath},
		{"Status", func() string {
			if result.Valid {
				return pterm.FgGreen.Sprint("✅ Valid")
			}
			return pterm.FgRed.Sprint("❌ Invalid")
		}()},
		{"Errors", func() string {
			if errorCount == 0 {
				return pterm.FgGreen.Sprint("0")
			}
			return pterm.FgRed.Sprintf("%d", errorCount)
		}()},
		{"Warnings", func() string {
			if warningCount == 0 {
				return pterm.FgGreen.Sprint("0")
			}
			return pterm.FgYellow.Sprintf("%d", warningCount)
		}()},
	}

	pterm.DefaultTable.WithHasHeader().WithData(summaryData).Render()
	pterm.Println()

	// Show errors if any
	if len(result.Errors) > 0 {
		pterm.DefaultSection.Println("Errors")
		for _, err := range result.Errors {
			pterm.Error.Printf("  [%s] %s\n", err.Field, err.Message)
		}
		pterm.Println()
	}

	// Show warnings if any
	if len(result.Warnings) > 0 {
		pterm.DefaultSection.Println("Warnings")
		for _, warn := range result.Warnings {
			pterm.Warning.Printf("  [%s] %s\n", warn.Field, warn.Message)
		}
		pterm.Println()
	}

	// Final message
	if result.Valid {
		pterm.Success.Println("Configuration validation passed!")
	} else {
		pterm.Error.Println("Configuration validation failed. Please fix the errors above.")
		return fmt.Errorf("validation failed with %d error(s)", errorCount)
	}

	return nil
}
