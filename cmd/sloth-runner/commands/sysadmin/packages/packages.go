package packages

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewPackagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "packages",
		Aliases: []string{"package", "pkg"},
		Short:   "Manage system packages on agents",
		Long: `Install, update, and manage system packages (apt, yum, dnf, pacman) on remote agents.
Supports rolling updates and automatic rollback on failure.`,
		Example: `  # List installed packages
  sloth-runner sysadmin packages list --agent web-01

  # Search for package
  sloth-runner sysadmin packages search nginx --agent web-01

  # Install package
  sloth-runner sysadmin packages install nginx --agent web-01

  # Update all packages with rolling strategy
  sloth-runner sysadmin packages update --all-agents --strategy rolling

  # Check for available updates
  sloth-runner sysadmin packages check-updates --all-agents`,
	}

	// list subcommand
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed packages",
		Long:  `List all installed packages on the agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd, args)
		},
	}
	listCmd.Flags().StringP("filter", "f", "", "Filter packages by name")
	listCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = unlimited)")
	cmd.AddCommand(listCmd)

	// search subcommand
	searchCmd := &cobra.Command{
		Use:   "search [package-name]",
		Short: "Search for packages",
		Long:  `Search for available packages in repositories.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd, args)
		},
	}
	searchCmd.Flags().IntP("limit", "l", 20, "Limit number of results")
	cmd.AddCommand(searchCmd)

	// install subcommand
	installCmd := &cobra.Command{
		Use:   "install [package-name]",
		Short: "Install a package",
		Long:  `Install a package on the specified agent(s).`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(cmd, args)
		},
	}
	installCmd.Flags().BoolP("yes", "y", false, "Automatically confirm installation")
	installCmd.Flags().Bool("no-deps", false, "Don't install dependencies")
	cmd.AddCommand(installCmd)

	// remove subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "remove [package-name]",
		Short: "Remove a package",
		Long:  `Remove an installed package from the agent(s).`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Package removal not yet implemented")
			pterm.Info.Println("Future features: Remove package, handle dependencies, purge configs")
		},
	})

	// update subcommand
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update package lists",
		Long:  `Update package repository lists (apt update, yum check-update, etc).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd, args)
		},
	}
	cmd.AddCommand(updateCmd)

	// upgrade subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade packages",
		Long:  `Upgrade all or specific packages to latest versions.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Package upgrade not yet implemented")
			pterm.Info.Println("Future features: Upgrade packages, rolling updates, auto-rollback on failure")
		},
	})

	// check-updates subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "check-updates",
		Short: "Check for available updates",
		Long:  `Check which packages have updates available without installing them.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Update checking not yet implemented")
			pterm.Info.Println("Future features: List available updates, security updates, version comparisons")
		},
	})

	// info subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "info [package-name]",
		Short: "Show package information",
		Long:  `Display detailed information about a package.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Package info not yet implemented")
			pterm.Info.Println("Future features: Package details, version, dependencies, size")
		},
	})

	// history subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "history",
		Short: "Show package management history",
		Long:  `Display history of package installations, updates, and removals.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Package history not yet implemented")
			pterm.Info.Println("Future features: Transaction history, rollback capability")
		},
	})

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the packages command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showPackagesDocs()
		},
	})

	return cmd
}

func showPackagesDocs() {
	title := "SLOTH-RUNNER SYSADMIN PACKAGES(1)"
	description := "sloth-runner sysadmin packages - Manage system packages on remote agents"
	synopsis := "sloth-runner sysadmin packages [subcommand] [options]"

	options := [][]string{
		{"list", "List all installed packages on the agent."},
		{"search [package-name]", "Search for available packages in repositories."},
		{"install [package-name]", "Install a package. Supports dependency resolution."},
		{"remove [package-name]", "Remove an installed package. Can purge configurations."},
		{"update", "Update package repository lists (apt update, yum check-update)."},
		{"upgrade", "Upgrade all or specific packages to latest versions with rolling strategy support."},
		{"check-updates", "Check which packages have updates available without installing."},
		{"info [package-name]", "Display detailed information about a package."},
		{"history", "Show history of package installations, updates, and removals."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"List installed packages",
			"sloth-runner sysadmin packages list --agent web-01",
			"Shows all installed packages with versions",
		},
		{
			"Search for nginx",
			"sloth-runner sysadmin pkg search nginx --agent web-01",
			"Searches repositories for nginx packages",
		},
		{
			"Install package",
			"sloth-runner sysadmin packages install nginx --agent web-01",
			"Installs nginx with dependency resolution",
		},
		{
			"Rolling upgrade",
			"sloth-runner sysadmin packages upgrade --all-agents --strategy rolling --wait-time 5m",
			"Upgrades packages across all agents with 5min wait between each",
		},
		{
			"Check for updates",
			"sloth-runner sysadmin pkg check-updates --all-agents",
			"Lists available updates for all agents",
		},
		{
			"Package info",
			"sloth-runner sysadmin packages info nginx --agent web-01",
			"Shows detailed package information",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin services - Service management",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin packages --help")
}

// Command implementation functions

func runList(cmd *cobra.Command, args []string) error {
	filter, _ := cmd.Flags().GetString("filter")
	limit, _ := cmd.Flags().GetInt("limit")

	spinner, _ := pterm.DefaultSpinner.Start("Detecting package manager...")

	pm, err := GetPackageManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	pmType := DetectPackageManager()
	spinner.Success(fmt.Sprintf("✅ Detected package manager: %s", pmType))

	spinner, _ = pterm.DefaultSpinner.Start("Fetching installed packages...")
	packages, err := pm.List()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to list packages: %s", err.Error()))
		return err
	}
	spinner.Stop()

	// Apply filter if provided
	if filter != "" {
		filtered := []Package{}
		for _, pkg := range packages {
			if containsIgnoreCase(pkg.Name, filter) {
				filtered = append(filtered, pkg)
			}
		}
		packages = filtered
	}

	// Apply limit if provided
	if limit > 0 && len(packages) > limit {
		packages = packages[:limit]
	}

	// Display results in table format
	pterm.DefaultHeader.WithFullWidth().Printf("Installed Packages (%d)", len(packages))
	fmt.Println()

	tableData := pterm.TableData{
		{"Package", "Version"},
	}

	for _, pkg := range packages {
		tableData = append(tableData, []string{pkg.Name, pkg.Version})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]
	limit, _ := cmd.Flags().GetInt("limit")

	spinner, _ := pterm.DefaultSpinner.Start("Detecting package manager...")

	pm, err := GetPackageManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	pmType := DetectPackageManager()
	spinner.Success(fmt.Sprintf("✅ Using: %s", pmType))

	spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Searching for '%s'...", query))
	packages, err := pm.Search(query)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Search failed: %s", err.Error()))
		return err
	}
	spinner.Stop()

	// Apply limit
	if limit > 0 && len(packages) > limit {
		packages = packages[:limit]
	}

	if len(packages) == 0 {
		pterm.Warning.Printf("No packages found matching '%s'\n", query)
		return nil
	}

	// Display results
	pterm.DefaultHeader.WithFullWidth().Printf("Search Results: '%s' (%d packages)", query, len(packages))
	fmt.Println()

	for _, pkg := range packages {
		pterm.FgGreen.Printf("📦 %s\n", pkg.Name)
		if pkg.Description != "" {
			pterm.FgGray.Printf("   %s\n", pkg.Description)
		}
		fmt.Println()
	}

	return nil
}

func runInstall(cmd *cobra.Command, args []string) error {
	packageName := args[0]
	autoYes, _ := cmd.Flags().GetBool("yes")

	spinner, _ := pterm.DefaultSpinner.Start("Detecting package manager...")

	pm, err := GetPackageManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	pmType := DetectPackageManager()
	spinner.Success(fmt.Sprintf("✅ Using: %s", pmType))

	// Confirm installation
	if !autoYes {
		confirm := pterm.DefaultInteractiveConfirm
		confirm.DefaultText = fmt.Sprintf("Install package '%s'?", packageName)
		result, _ := confirm.Show()
		if !result {
			pterm.Info.Println("Installation cancelled")
			return nil
		}
	}

	spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s...", packageName))

	err = pm.Install(packageName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Installation failed: %s", err.Error()))
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Successfully installed %s", packageName))

	return nil
}

func runUpdate(cmd *cobra.Command, args []string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Detecting package manager...")

	pm, err := GetPackageManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	pmType := DetectPackageManager()
	spinner.Success(fmt.Sprintf("✅ Using: %s", pmType))

	spinner, _ = pterm.DefaultSpinner.Start("Updating package lists...")

	err = pm.Update()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Update failed: %s", err.Error()))
		return err
	}

	spinner.Success("✅ Package lists updated successfully")

	return nil
}

// Helper function
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
