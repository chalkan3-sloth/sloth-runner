package systemd

import (
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewSystemdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "systemd",
		Short:   "Systemd service management",
		Long:    `Manage systemd services: list, start, stop, restart, enable, disable, and view logs.`,
		Aliases: []string{"service", "svc"},
	}

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List systemd services",
		Run: func(cmd *cobra.Command, args []string) {
			status, _ := cmd.Flags().GetString("status")
			filter, _ := cmd.Flags().GetString("filter")
			typeFlag, _ := cmd.Flags().GetString("type")

			if err := runList(status, filter, typeFlag); err != nil {
				pterm.Error.Printf("Failed to list services: %v\n", err)
			}
		},
	}
	listCmd.Flags().StringP("status", "s", "all", "Filter by status (all, running, stopped, failed)")
	listCmd.Flags().StringP("filter", "f", "", "Filter by service name or description")
	listCmd.Flags().StringP("type", "t", "service", "Unit type (service, socket, timer)")
	cmd.AddCommand(listCmd)

	// status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show detailed service status",
		Run: func(cmd *cobra.Command, args []string) {
			service, _ := cmd.Flags().GetString("service")

			if service == "" {
				pterm.Error.Println("Service name is required (use --service)")
				return
			}

			if err := runStatus(service); err != nil {
				pterm.Error.Printf("Failed to get service status: %v\n", err)
			}
		},
	}
	statusCmd.Flags().StringP("service", "s", "", "Service name")
	cmd.AddCommand(statusCmd)

	// start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a service",
		Run: func(cmd *cobra.Command, args []string) {
			service, _ := cmd.Flags().GetString("service")

			if service == "" {
				pterm.Error.Println("Service name is required (use --service)")
				return
			}

			if err := runStart(service); err != nil {
				pterm.Error.Printf("Failed to start service: %v\n", err)
			}
		},
	}
	startCmd.Flags().StringP("service", "s", "", "Service name")
	cmd.AddCommand(startCmd)

	// stop command
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a service",
		Run: func(cmd *cobra.Command, args []string) {
			service, _ := cmd.Flags().GetString("service")

			if service == "" {
				pterm.Error.Println("Service name is required (use --service)")
				return
			}

			if err := runStop(service); err != nil {
				pterm.Error.Printf("Failed to stop service: %v\n", err)
			}
		},
	}
	stopCmd.Flags().StringP("service", "s", "", "Service name")
	cmd.AddCommand(stopCmd)

	// restart command
	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart a service",
		Run: func(cmd *cobra.Command, args []string) {
			service, _ := cmd.Flags().GetString("service")

			if service == "" {
				pterm.Error.Println("Service name is required (use --service)")
				return
			}

			if err := runRestart(service); err != nil {
				pterm.Error.Printf("Failed to restart service: %v\n", err)
			}
		},
	}
	restartCmd.Flags().StringP("service", "s", "", "Service name")
	cmd.AddCommand(restartCmd)

	// enable command
	enableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable a service at boot",
		Run: func(cmd *cobra.Command, args []string) {
			service, _ := cmd.Flags().GetString("service")

			if service == "" {
				pterm.Error.Println("Service name is required (use --service)")
				return
			}

			if err := runEnable(service); err != nil {
				pterm.Error.Printf("Failed to enable service: %v\n", err)
			}
		},
	}
	enableCmd.Flags().StringP("service", "s", "", "Service name")
	cmd.AddCommand(enableCmd)

	// disable command
	disableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable a service at boot",
		Run: func(cmd *cobra.Command, args []string) {
			service, _ := cmd.Flags().GetString("service")

			if service == "" {
				pterm.Error.Println("Service name is required (use --service)")
				return
			}

			if err := runDisable(service); err != nil {
				pterm.Error.Printf("Failed to disable service: %v\n", err)
			}
		},
	}
	disableCmd.Flags().StringP("service", "s", "", "Service name")
	cmd.AddCommand(disableCmd)

	// logs command
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "View service logs",
		Run: func(cmd *cobra.Command, args []string) {
			service, _ := cmd.Flags().GetString("service")
			lines, _ := cmd.Flags().GetInt("lines")
			follow, _ := cmd.Flags().GetBool("follow")

			if service == "" {
				pterm.Error.Println("Service name is required (use --service)")
				return
			}

			if err := runLogs(service, lines, follow); err != nil {
				pterm.Error.Printf("Failed to get logs: %v\n", err)
			}
		},
	}
	logsCmd.Flags().StringP("service", "s", "", "Service name")
	logsCmd.Flags().IntP("lines", "n", 50, "Number of log lines to show")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	cmd.AddCommand(logsCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the systemd command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showSystemdDocs()
		},
	})

	return cmd
}

func showSystemdDocs() {
	title := "SLOTH-RUNNER SYSADMIN SYSTEMD(1)"
	description := "sloth-runner sysadmin systemd - Systemd service management"
	synopsis := "sloth-runner sysadmin systemd [subcommand] [options]"

	options := [][]string{
		{"list", "List systemd services with filtering by status, name, or type."},
		{"status", "Show detailed status information for a specific service."},
		{"start", "Start a systemd service."},
		{"stop", "Stop a systemd service."},
		{"restart", "Restart a systemd service."},
		{"enable", "Enable a service to start at boot."},
		{"disable", "Disable a service from starting at boot."},
		{"logs", "View service logs from journald."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"List all running services",
			"sloth-runner sysadmin systemd list --status running",
			"Shows only active services",
		},
		{
			"List failed services",
			"sloth-runner sysadmin systemd list --status failed",
			"Shows services in failed state",
		},
		{
			"Filter services by name",
			"sloth-runner sysadmin service list --filter nginx",
			"Shows all services containing 'nginx'",
		},
		{
			"Show service status",
			"sloth-runner sysadmin svc status --service nginx",
			"Displays detailed information about nginx service",
		},
		{
			"Start a service",
			"sloth-runner sysadmin systemd start --service nginx",
			"Starts the nginx service",
		},
		{
			"Restart a service",
			"sloth-runner sysadmin systemd restart --service nginx",
			"Restarts the nginx service",
		},
		{
			"Enable service at boot",
			"sloth-runner sysadmin systemd enable --service nginx",
			"Configures nginx to start on system boot",
		},
		{
			"View service logs",
			"sloth-runner sysadmin systemd logs --service nginx --lines 100",
			"Shows last 100 log lines for nginx",
		},
		{
			"Follow service logs",
			"sloth-runner sysadmin systemd logs --service nginx --follow",
			"Shows logs in real-time (like tail -f)",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin process - Process management",
		"sloth-runner sysadmin performance - System performance monitoring",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin systemd --help")
}

// runList lista serviços
func runList(status string, filter string, typeFlag string) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Systemd Services")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading services...")

	options := ListOptions{
		Status: status,
		Filter: filter,
		Type:   typeFlag,
	}

	services, err := manager.List(options)
	if err != nil {
		spinner.Fail("Failed to load services")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d services", len(services)))
	pterm.Println()

	if len(services) == 0 {
		pterm.Info.Println("No services found")
		return nil
	}

	// Service table
	tableData := pterm.TableData{
		{"Name", "Load", "Active", "Sub", "Description"},
	}

	for _, s := range services {
		// Color based on state
		activeState := s.ActiveState
		if s.ActiveState == "active" {
			activeState = pterm.FgGreen.Sprint(s.ActiveState)
		} else if s.ActiveState == "failed" {
			activeState = pterm.FgRed.Sprint(s.ActiveState)
		} else if s.ActiveState == "inactive" {
			activeState = pterm.FgYellow.Sprint(s.ActiveState)
		}

		tableData = append(tableData, []string{
			truncate(s.Name, 40),
			s.LoadState,
			activeState,
			s.SubState,
			truncate(s.Description, 40),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Info.Printf("Status filter: %s | Type: %s | Total: %d\n", status, typeFlag, len(services))
	if filter != "" {
		pterm.Info.Printf("Filter: %s\n", filter)
	}

	return nil
}

// runStatus mostra status detalhado de um serviço
func runStatus(service string) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Service Status")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading service details...")

	detail, err := manager.Status(service)
	if err != nil {
		spinner.Fail("Failed to load service")
		return err
	}

	spinner.Success("✅ Service information loaded")
	pterm.Println()

	// Basic Info
	pterm.DefaultSection.Println("Basic Information")
	basicData := pterm.TableData{
		{"Property", "Value"},
		{"Name", detail.Name},
		{"Description", detail.Description},
		{"Load State", detail.LoadState},
		{"Active State", getColoredState(detail.ActiveState)},
		{"Sub State", detail.SubState},
	}
	pterm.DefaultTable.WithHasHeader().WithData(basicData).Render()
	pterm.Println()

	// Process Info
	if detail.MainPID > 0 {
		pterm.DefaultSection.Println("Process Information")
		processData := pterm.TableData{
			{"Property", "Value"},
			{"Main PID", fmt.Sprintf("%d", detail.MainPID)},
			{"Memory", formatBytes(detail.Memory)},
			{"CPU Time", fmt.Sprintf("%.2fs", detail.CPUUsage)},
		}
		if detail.TasksCurrent > 0 {
			processData = append(processData, []string{"Tasks", fmt.Sprintf("%d (max: %d)", detail.TasksCurrent, detail.TasksMax)})
		}
		if detail.RestartCount > 0 {
			processData = append(processData, []string{"Restarts", fmt.Sprintf("%d", detail.RestartCount)})
		}
		pterm.DefaultTable.WithHasHeader().WithData(processData).Render()
		pterm.Println()
	}

	// Configuration
	pterm.DefaultSection.Println("Configuration")
	configData := pterm.TableData{
		{"Property", "Value"},
	}
	if detail.Fragment != "" {
		configData = append(configData, []string{"Unit File", detail.Fragment})
	}
	if detail.User != "" {
		configData = append(configData, []string{"User", detail.User})
	}
	if detail.Group != "" {
		configData = append(configData, []string{"Group", detail.Group})
	}
	if detail.Restart != "" {
		configData = append(configData, []string{"Restart", detail.Restart})
	}
	if detail.TimeoutStartS > 0 {
		configData = append(configData, []string{"Start Timeout", fmt.Sprintf("%ds", detail.TimeoutStartS)})
	}
	if detail.TimeoutStopS > 0 {
		configData = append(configData, []string{"Stop Timeout", fmt.Sprintf("%ds", detail.TimeoutStopS)})
	}
	if len(configData) > 1 {
		pterm.DefaultTable.WithHasHeader().WithData(configData).Render()
		pterm.Println()
	}

	// Exec Commands
	if detail.ExecStart != "" || detail.ExecStop != "" || detail.ExecReload != "" {
		pterm.DefaultSection.Println("Commands")
		if detail.ExecStart != "" {
			fmt.Printf("    Start: %s\n", truncate(detail.ExecStart, 80))
		}
		if detail.ExecReload != "" {
			fmt.Printf("    Reload: %s\n", truncate(detail.ExecReload, 80))
		}
		if detail.ExecStop != "" {
			fmt.Printf("    Stop: %s\n", truncate(detail.ExecStop, 80))
		}
		pterm.Println()
	}

	// Active Since
	if !detail.ActiveSince.IsZero() {
		pterm.Info.Printf("Active since: %s (%s ago)\n",
			detail.ActiveSince.Format("2006-01-02 15:04:05"),
			time.Since(detail.ActiveSince).Round(time.Second))
	}

	return nil
}

// runStart inicia um serviço
func runStart(service string) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Start Service")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Starting %s...", service))

	err := manager.Start(service)
	if err != nil {
		spinner.Fail("Failed to start service")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Service %s started successfully", service))
	pterm.Println()

	return nil
}

// runStop para um serviço
func runStop(service string) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Stop Service")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Stopping %s...", service))

	err := manager.Stop(service)
	if err != nil {
		spinner.Fail("Failed to stop service")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Service %s stopped successfully", service))
	pterm.Println()

	return nil
}

// runRestart reinicia um serviço
func runRestart(service string) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Restart Service")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Restarting %s...", service))

	err := manager.Restart(service)
	if err != nil {
		spinner.Fail("Failed to restart service")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Service %s restarted successfully", service))
	pterm.Println()

	return nil
}

// runEnable habilita um serviço no boot
func runEnable(service string) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Enable Service")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Enabling %s...", service))

	err := manager.Enable(service)
	if err != nil {
		spinner.Fail("Failed to enable service")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Service %s enabled at boot", service))
	pterm.Println()

	return nil
}

// runDisable desabilita um serviço no boot
func runDisable(service string) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Disable Service")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Disabling %s...", service))

	err := manager.Disable(service)
	if err != nil {
		spinner.Fail("Failed to disable service")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Service %s disabled from boot", service))
	pterm.Println()

	return nil
}

// runLogs exibe logs de um serviço
func runLogs(service string, lines int, follow bool) error {
	manager := NewSystemdManager()

	pterm.DefaultHeader.WithFullWidth().Println("Service Logs")
	pterm.Println()

	if follow {
		pterm.Info.Printf("Following logs for %s (Ctrl+C to stop)...\n\n", service)
	} else {
		pterm.Info.Printf("Showing last %d lines for %s...\n\n", lines, service)
	}

	logs, err := manager.Logs(service, lines, follow)
	if err != nil {
		return err
	}

	fmt.Print(logs)

	return nil
}

// getColoredState retorna estado colorido
func getColoredState(state string) string {
	switch state {
	case "active":
		return pterm.FgGreen.Sprint(state)
	case "failed":
		return pterm.FgRed.Sprint(state)
	case "inactive":
		return pterm.FgYellow.Sprint(state)
	default:
		return state
	}
}

// formatBytes formata bytes para formato legível
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// truncate trunca uma string
func truncate(s string, maxLen int) string {
	// Remove caracteres de controle e espaços extras
	s = strings.TrimSpace(s)
	s = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, s)

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
