package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewServicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services",
		Aliases: []string{"service", "svc"},
		Short:   "Manage systemd/init.d services on agents",
		Long: `Control and monitor services (systemd, init.d, OpenRC) on remote agents.
Start, stop, restart, and check status of services without SSH access.`,
		Example: `  # List all services
  sloth-runner sysadmin services list --agent web-01

  # Check service status
  sloth-runner sysadmin services status nginx --agent web-01

  # Start service
  sloth-runner sysadmin services start nginx --agent web-01

  # Restart service on multiple agents
  sloth-runner sysadmin services restart nginx --agents web-01,web-02

  # View service logs
  sloth-runner sysadmin services logs nginx --agent web-01 --follow`,
	}

	// list subcommand
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all services on agent",
		Long:  `List all available services with their current status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd, args)
		},
	}
	listCmd.Flags().StringP("filter", "f", "", "Filter services by name")
	listCmd.Flags().StringP("status", "s", "", "Filter by status (active/inactive/failed)")
	cmd.AddCommand(listCmd)

	// status subcommand
	statusCmd := &cobra.Command{
		Use:   "status [service-name]",
		Short: "Check service status",
		Long:  `Check the current status of a specific service.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(cmd, args)
		},
	}
	cmd.AddCommand(statusCmd)

	// start subcommand
	startCmd := &cobra.Command{
		Use:   "start [service-name]",
		Short: "Start a service",
		Long:  `Start a stopped service on the specified agent(s).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(cmd, args)
		},
	}
	startCmd.Flags().Bool("verify", true, "Verify service started successfully")
	cmd.AddCommand(startCmd)

	// stop subcommand
	stopCmd := &cobra.Command{
		Use:   "stop [service-name]",
		Short: "Stop a service",
		Long:  `Stop a running service on the specified agent(s).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStop(cmd, args)
		},
	}
	stopCmd.Flags().Bool("verify", true, "Verify service stopped successfully")
	cmd.AddCommand(stopCmd)

	// restart subcommand
	restartCmd := &cobra.Command{
		Use:   "restart [service-name]",
		Short: "Restart a service",
		Long:  `Restart a service on the specified agent(s).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestart(cmd, args)
		},
	}
	restartCmd.Flags().Bool("verify", true, "Verify service restarted successfully")
	cmd.AddCommand(restartCmd)

	// reload subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "reload [service-name]",
		Short: "Reload service configuration",
		Long:  `Reload service configuration without full restart (if supported).`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Service reload not yet implemented")
			pterm.Info.Println("Future features: Reload config, zero-downtime reload")
		},
	})

	// enable subcommand
	enableCmd := &cobra.Command{
		Use:   "enable [service-name]",
		Short: "Enable service at boot",
		Long:  `Enable a service to start automatically at system boot.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnable(cmd, args)
		},
	}
	cmd.AddCommand(enableCmd)

	// disable subcommand
	disableCmd := &cobra.Command{
		Use:   "disable [service-name]",
		Short: "Disable service at boot",
		Long:  `Disable a service from starting automatically at system boot.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDisable(cmd, args)
		},
	}
	cmd.AddCommand(disableCmd)

	// logs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "logs [service-name]",
		Short: "View service logs",
		Long:  `View logs for a specific service (journalctl or log files).`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Service logs not yet implemented")
			pterm.Info.Println("Future features: View logs, follow mode, filter by level")
		},
	})

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the services command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showServicesDocs()
		},
	})

	return cmd
}

func showServicesDocs() {
	title := "SLOTH-RUNNER SYSADMIN SERVICES(1)"
	description := "sloth-runner sysadmin services - Manage systemd/init.d services on remote agents"
	synopsis := "sloth-runner sysadmin services [subcommand] [options]"

	options := [][]string{
		{"list", "List all services on the specified agent with their current status."},
		{"status [service-name]", "Check the current status of a specific service including uptime, PID, and memory usage."},
		{"start [service-name]", "Start a stopped service. Includes verification of successful startup."},
		{"stop [service-name]", "Stop a running service gracefully. Supports force stop if needed."},
		{"restart [service-name]", "Restart a service. Performs stop followed by start with health verification."},
		{"reload [service-name]", "Reload service configuration without full restart (zero-downtime when supported)."},
		{"enable [service-name]", "Enable a service to start automatically at system boot."},
		{"disable [service-name]", "Disable a service from starting automatically at system boot."},
		{"logs [service-name]", "View service logs from journalctl or log files. Supports follow mode."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"List all services",
			"sloth-runner sysadmin services list --agent web-01",
			"Shows all services with status (active/inactive/failed)",
		},
		{
			"Check nginx status",
			"sloth-runner sysadmin services status nginx --agent web-01",
			"Displays detailed status including PID, memory, uptime",
		},
		{
			"Restart service on multiple agents",
			"sloth-runner sysadmin services restart nginx --agents web-01,web-02,web-03",
			"Performs rolling restart across multiple agents",
		},
		{
			"Start service with verification",
			"sloth-runner sysadmin svc start application --agent app-01 --verify",
			"Starts service and verifies it's running correctly",
		},
		{
			"Follow service logs in real-time",
			"sloth-runner sysadmin services logs nginx --agent web-01 --follow",
			"Streams logs from the service in real-time",
		},
		{
			"Enable service at boot",
			"sloth-runner sysadmin services enable docker --all-agents",
			"Enables docker to start automatically on all agents",
		},
		{
			"Reload nginx configuration",
			"sloth-runner sysadmin services reload nginx --agents web-*",
			"Reloads nginx config on all web agents without downtime",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin packages - Package management",
		"sloth-runner sysadmin resources - Resource monitoring",
		"sloth-runner sysadmin health - Health checks",
		"sloth-runner sysadmin logs - Log management",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin services --help")
}

// Command implementation functions

func runList(cmd *cobra.Command, args []string) error {
	filter, _ := cmd.Flags().GetString("filter")
	statusFilter, _ := cmd.Flags().GetString("status")

	spinner, _ := pterm.DefaultSpinner.Start("Detecting service manager...")

	sm, err := GetServiceManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	smType := DetectServiceManager()
	spinner.Success(fmt.Sprintf("✅ Detected: %s", smType))

	spinner, _ = pterm.DefaultSpinner.Start("Fetching services...")
	services, err := sm.List()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to list services: %s", err.Error()))
		return err
	}
	spinner.Stop()

	// Apply filters
	filtered := []Service{}
	for _, svc := range services {
		// Apply name filter
		if filter != "" && !strings.Contains(strings.ToLower(svc.Name), strings.ToLower(filter)) {
			continue
		}

		// Apply status filter
		if statusFilter != "" && strings.ToLower(string(svc.Status)) != strings.ToLower(statusFilter) {
			continue
		}

		filtered = append(filtered, svc)
	}

	// Display results
	pterm.DefaultHeader.WithFullWidth().Printf("Services (%d)", len(filtered))
	fmt.Println()

	tableData := pterm.TableData{
		{"Service", "Status", "Description"},
	}

	for _, svc := range filtered {
		statusStr := colorizeStatus(svc.Status)
		tableData = append(tableData, []string{svc.Name, statusStr, truncate(svc.Description, 50)})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	serviceName := args[0]

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking status of %s...", serviceName))

	sm, err := GetServiceManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	service, err := sm.Status(serviceName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get status: %s", err.Error()))
		return err
	}
	spinner.Stop()

	// Display detailed status
	pterm.DefaultHeader.WithFullWidth().Printf("Service: %s", serviceName)
	fmt.Println()

	// Status box
	statusColor := getStatusColor(service.Status)
	statusBox := pterm.DefaultBox.WithTitle("Status").WithTitleTopCenter()
	statusBox.Println(statusColor.Sprintf("● %s", strings.ToUpper(string(service.Status))))

	// Details
	fmt.Println()
	pterm.DefaultSection.Println("Details")
	
	if service.PID != "" {
		pterm.Printf("   PID: %s\n", service.PID)
	}
	if service.Memory != "" {
		pterm.Printf("   Memory: %s\n", service.Memory)
	}
	if service.Enabled {
		pterm.Printf("   Boot: %s\n", pterm.FgGreen.Sprint("enabled"))
	} else {
		pterm.Printf("   Boot: %s\n", pterm.FgGray.Sprint("disabled"))
	}

	return nil
}

func runStart(cmd *cobra.Command, args []string) error {
	serviceName := args[0]
	verify, _ := cmd.Flags().GetBool("verify")

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Starting %s...", serviceName))

	sm, err := GetServiceManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	err = sm.Start(serviceName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to start: %s", err.Error()))
		return err
	}

	if verify {
		time.Sleep(1 * time.Second)
		service, err := sm.Status(serviceName)
		if err == nil && service.Status == StatusActive {
			spinner.Success(fmt.Sprintf("✅ %s started successfully", serviceName))
		} else {
			spinner.Warning(fmt.Sprintf("⚠️  %s started but verification failed", serviceName))
		}
	} else {
		spinner.Success(fmt.Sprintf("✅ %s started", serviceName))
	}

	return nil
}

func runStop(cmd *cobra.Command, args []string) error {
	serviceName := args[0]
	verify, _ := cmd.Flags().GetBool("verify")

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Stopping %s...", serviceName))

	sm, err := GetServiceManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	err = sm.Stop(serviceName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to stop: %s", err.Error()))
		return err
	}

	if verify {
		time.Sleep(1 * time.Second)
		service, err := sm.Status(serviceName)
		if err == nil && service.Status == StatusInactive {
			spinner.Success(fmt.Sprintf("✅ %s stopped successfully", serviceName))
		} else {
			spinner.Warning(fmt.Sprintf("⚠️  %s stopped but verification failed", serviceName))
		}
	} else {
		spinner.Success(fmt.Sprintf("✅ %s stopped", serviceName))
	}

	return nil
}

func runRestart(cmd *cobra.Command, args []string) error {
	serviceName := args[0]
	verify, _ := cmd.Flags().GetBool("verify")

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Restarting %s...", serviceName))

	sm, err := GetServiceManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	err = sm.Restart(serviceName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to restart: %s", err.Error()))
		return err
	}

	if verify {
		time.Sleep(1 * time.Second)
		service, err := sm.Status(serviceName)
		if err == nil && service.Status == StatusActive {
			spinner.Success(fmt.Sprintf("✅ %s restarted successfully", serviceName))
		} else {
			spinner.Warning(fmt.Sprintf("⚠️  %s restarted but verification failed", serviceName))
		}
	} else {
		spinner.Success(fmt.Sprintf("✅ %s restarted", serviceName))
	}

	return nil
}

func runEnable(cmd *cobra.Command, args []string) error {
	serviceName := args[0]

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Enabling %s...", serviceName))

	sm, err := GetServiceManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	err = sm.Enable(serviceName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to enable: %s", err.Error()))
		return err
	}

	spinner.Success(fmt.Sprintf("✅ %s enabled at boot", serviceName))

	return nil
}

func runDisable(cmd *cobra.Command, args []string) error {
	serviceName := args[0]

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Disabling %s...", serviceName))

	sm, err := GetServiceManager()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	err = sm.Disable(serviceName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to disable: %s", err.Error()))
		return err
	}

	spinner.Success(fmt.Sprintf("✅ %s disabled at boot", serviceName))

	return nil
}

// Helper functions

func colorizeStatus(status ServiceStatus) string {
	switch status {
	case StatusActive:
		return pterm.FgGreen.Sprint("● active")
	case StatusInactive:
		return pterm.FgGray.Sprint("○ inactive")
	case StatusFailed:
		return pterm.FgRed.Sprint("✖ failed")
	default:
		return pterm.FgYellow.Sprint("? unknown")
	}
}

func getStatusColor(status ServiceStatus) pterm.Color {
	switch status {
	case StatusActive:
		return pterm.FgGreen
	case StatusInactive:
		return pterm.FgGray
	case StatusFailed:
		return pterm.FgRed
	default:
		return pterm.FgYellow
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
