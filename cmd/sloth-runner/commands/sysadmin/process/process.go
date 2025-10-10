package process

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewProcessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "process",
		Short:   "Process management and monitoring",
		Long:    `List, monitor, and manage system processes with advanced filtering and sorting.`,
		Aliases: []string{"proc", "ps"},
	}

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List running processes",
		Run: func(cmd *cobra.Command, args []string) {
			sortBy, _ := cmd.Flags().GetString("sort")
			top, _ := cmd.Flags().GetInt("top")
			filter, _ := cmd.Flags().GetString("filter")
			user, _ := cmd.Flags().GetString("user")

			if err := runList(sortBy, top, filter, user); err != nil {
				pterm.Error.Printf("Failed to list processes: %v\n", err)
			}
		},
	}
	listCmd.Flags().StringP("sort", "s", "cpu", "Sort by (cpu, memory, name, pid)")
	listCmd.Flags().IntP("top", "t", 20, "Show top N processes")
	listCmd.Flags().StringP("filter", "f", "", "Filter by process name or command")
	listCmd.Flags().StringP("user", "u", "", "Filter by username")
	cmd.AddCommand(listCmd)

	// kill command
	killCmd := &cobra.Command{
		Use:   "kill",
		Short: "Kill a process",
		Run: func(cmd *cobra.Command, args []string) {
			pidStr, _ := cmd.Flags().GetString("pid")
			signal, _ := cmd.Flags().GetString("signal")
			force, _ := cmd.Flags().GetBool("force")

			if pidStr == "" {
				pterm.Error.Println("PID is required (use --pid)")
				return
			}

			pid, err := strconv.ParseInt(pidStr, 10, 32)
			if err != nil {
				pterm.Error.Printf("Invalid PID: %v\n", err)
				return
			}

			if force {
				signal = "SIGKILL"
			}

			if err := runKill(int32(pid), signal); err != nil {
				pterm.Error.Printf("Failed to kill process: %v\n", err)
			}
		},
	}
	killCmd.Flags().StringP("pid", "p", "", "Process ID to kill")
	killCmd.Flags().StringP("signal", "s", "SIGTERM", "Signal to send (SIGTERM, SIGKILL, SIGINT, SIGHUP)")
	killCmd.Flags().BoolP("force", "f", false, "Force kill (SIGKILL)")
	cmd.AddCommand(killCmd)

	// info command
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show detailed process information",
		Run: func(cmd *cobra.Command, args []string) {
			pidStr, _ := cmd.Flags().GetString("pid")

			if pidStr == "" {
				pterm.Error.Println("PID is required (use --pid)")
				return
			}

			pid, err := strconv.ParseInt(pidStr, 10, 32)
			if err != nil {
				pterm.Error.Printf("Invalid PID: %v\n", err)
				return
			}

			if err := runInfo(int32(pid)); err != nil {
				pterm.Error.Printf("Failed to get process info: %v\n", err)
			}
		},
	}
	infoCmd.Flags().StringP("pid", "p", "", "Process ID")
	cmd.AddCommand(infoCmd)

	// monitor command
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor a process over time",
		Run: func(cmd *cobra.Command, args []string) {
			pidStr, _ := cmd.Flags().GetString("pid")
			duration, _ := cmd.Flags().GetDuration("duration")

			if pidStr == "" {
				pterm.Error.Println("PID is required (use --pid)")
				return
			}

			pid, err := strconv.ParseInt(pidStr, 10, 32)
			if err != nil {
				pterm.Error.Printf("Invalid PID: %v\n", err)
				return
			}

			if err := runMonitor(int32(pid), duration); err != nil {
				pterm.Error.Printf("Failed to monitor process: %v\n", err)
			}
		},
	}
	monitorCmd.Flags().StringP("pid", "p", "", "Process ID")
	monitorCmd.Flags().DurationP("duration", "d", 10*time.Second, "Monitoring duration")
	cmd.AddCommand(monitorCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the process command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showProcessDocs()
		},
	})

	return cmd
}

func showProcessDocs() {
	title := "SLOTH-RUNNER SYSADMIN PROCESS(1)"
	description := "sloth-runner sysadmin process - Process management and monitoring"
	synopsis := "sloth-runner sysadmin process [subcommand] [options]"

	options := [][]string{
		{"list", "List running processes with filtering and sorting options."},
		{"kill", "Terminate a process using various signals (SIGTERM, SIGKILL, etc.)."},
		{"info", "Display detailed information about a specific process."},
		{"monitor", "Monitor process resource usage over time."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"List top CPU consumers",
			"sloth-runner sysadmin process list --sort cpu --top 10",
			"Shows the top 10 processes by CPU usage",
		},
		{
			"List processes by memory",
			"sloth-runner sysadmin proc list --sort memory --top 20",
			"Shows the top 20 processes by memory usage",
		},
		{
			"Filter processes by name",
			"sloth-runner sysadmin ps list --filter nginx",
			"Shows all nginx processes",
		},
		{
			"Kill a process gracefully",
			"sloth-runner sysadmin process kill --pid 1234 --signal SIGTERM",
			"Sends SIGTERM to process 1234",
		},
		{
			"Force kill a process",
			"sloth-runner sysadmin process kill --pid 1234 --force",
			"Sends SIGKILL to process 1234",
		},
		{
			"Show process details",
			"sloth-runner sysadmin process info --pid 1234",
			"Shows detailed information about process 1234",
		},
		{
			"Monitor a process",
			"sloth-runner sysadmin process monitor --pid 1234 --duration 30s",
			"Monitors process 1234 for 30 seconds",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin performance - System performance monitoring",
		"sloth-runner sysadmin resources - Resource monitoring",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin process --help")
}

// runList lista processos
func runList(sortBy string, top int, filter string, user string) error {
	manager := NewProcessManager()

	pterm.DefaultHeader.WithFullWidth().Println("Process List")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading processes...")

	options := ListOptions{
		SortBy:     sortBy,
		Top:        top,
		Filter:     filter,
		UserFilter: user,
	}

	processes, err := manager.List(options)
	if err != nil {
		spinner.Fail("Failed to load processes")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d processes", len(processes)))
	pterm.Println()

	if len(processes) == 0 {
		pterm.Info.Println("No processes found")
		return nil
	}

	// Process table
	tableData := pterm.TableData{
		{"PID", "Name", "User", "CPU%", "Memory", "Threads", "Status", "Command"},
	}

	for _, p := range processes {
		tableData = append(tableData, []string{
			fmt.Sprintf("%d", p.PID),
			p.Name,
			p.Username,
			fmt.Sprintf("%.1f%%", p.CPUPercent),
			fmt.Sprintf("%.1f MB", p.MemoryMB),
			fmt.Sprintf("%d", p.NumThreads),
			p.Status,
			truncate(p.Cmdline, 40),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Info.Printf("Sorted by: %s | Showing: %d processes\n", sortBy, len(processes))
	if filter != "" {
		pterm.Info.Printf("Filter: %s\n", filter)
	}
	if user != "" {
		pterm.Info.Printf("User: %s\n", user)
	}

	return nil
}

// runKill mata um processo
func runKill(pid int32, signal string) error {
	manager := NewProcessManager()

	// Obtém info do processo antes de matar
	info, err := manager.Info(pid)
	if err != nil {
		return err
	}

	pterm.DefaultHeader.WithFullWidth().Println("Kill Process")
	pterm.Println()

	pterm.Warning.Printf("Killing process %d (%s) with signal %s\n", pid, info.Name, signal)
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Sending signal...")

	err = manager.Kill(pid, signal)
	if err != nil {
		spinner.Fail("Failed to kill process")
		return err
	}

	spinner.Success("✅ Process killed successfully")
	pterm.Println()

	pterm.Success.Printf("Process %d (%s) was terminated\n", pid, info.Name)

	return nil
}

// runInfo mostra informações detalhadas de um processo
func runInfo(pid int32) error {
	manager := NewProcessManager()

	pterm.DefaultHeader.WithFullWidth().Println("Process Information")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading process details...")

	detail, err := manager.Info(pid)
	if err != nil {
		spinner.Fail("Failed to load process")
		return err
	}

	spinner.Success("✅ Process information loaded")
	pterm.Println()

	// Basic Info
	pterm.DefaultSection.Println("Basic Information")
	basicData := pterm.TableData{
		{"Property", "Value"},
		{"PID", fmt.Sprintf("%d", detail.PID)},
		{"Name", detail.Name},
		{"User", detail.Username},
		{"Status", detail.Status},
		{"Parent PID", fmt.Sprintf("%d", detail.ParentPID)},
		{"Nice", fmt.Sprintf("%d", detail.Nice)},
		{"Threads", fmt.Sprintf("%d", detail.NumThreads)},
		{"File Descriptors", fmt.Sprintf("%d", detail.NumFDs)},
	}
	pterm.DefaultTable.WithHasHeader().WithData(basicData).Render()
	pterm.Println()

	// Resource Usage
	pterm.DefaultSection.Println("Resource Usage")
	resourceData := pterm.TableData{
		{"Metric", "Value"},
		{"CPU", fmt.Sprintf("%.1f%%", detail.CPUPercent)},
		{"Memory", fmt.Sprintf("%.1f MB (%.2f%%)", detail.MemoryMB, detail.MemoryPercent)},
	}
	if detail.IOCounters != nil {
		resourceData = append(resourceData, []string{
			"Read Bytes", fmt.Sprintf("%d", detail.IOCounters.ReadBytes),
		})
		resourceData = append(resourceData, []string{
			"Write Bytes", fmt.Sprintf("%d", detail.IOCounters.WriteBytes),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(resourceData).Render()
	pterm.Println()

	// Command Line
	if detail.Cmdline != "" {
		pterm.DefaultSection.Println("Command Line")
		fmt.Printf("    %s\n\n", detail.Cmdline)
	}

	// Connections
	if len(detail.Connections) > 0 {
		pterm.DefaultSection.Println(fmt.Sprintf("Network Connections (%d)", len(detail.Connections)))
		for i, conn := range detail.Connections {
			if i >= 5 {
				fmt.Printf("    ... and %d more\n", len(detail.Connections)-5)
				break
			}
			fmt.Printf("    • %s\n", conn)
		}
		pterm.Println()
	}

	// Open Files
	if len(detail.OpenFiles) > 0 {
		pterm.DefaultSection.Println(fmt.Sprintf("Open Files (%d)", len(detail.OpenFiles)))
		for i, file := range detail.OpenFiles {
			if i >= 5 {
				fmt.Printf("    ... and %d more\n", len(detail.OpenFiles)-5)
				break
			}
			fmt.Printf("    • %s\n", file)
		}
		pterm.Println()
	}

	return nil
}

// runMonitor monitora um processo ao longo do tempo
func runMonitor(pid int32, duration time.Duration) error {
	manager := NewProcessManager()

	pterm.DefaultHeader.WithFullWidth().Println("Process Monitor")
	pterm.Println()

	pterm.Info.Printf("Monitoring process %d for %v...\n", pid, duration)
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Collecting samples over %v...", duration))

	metrics, err := manager.Monitor(pid, duration)
	if err != nil {
		spinner.Fail("Failed to monitor process")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Collected %d samples", len(metrics.Samples)))
	pterm.Println()

	pterm.DefaultHeader.WithFullWidth().Println("Process Monitoring Report")
	pterm.Println()

	// Summary
	summaryData := pterm.TableData{
		{"Metric", "Average", "Maximum"},
		{"CPU Usage", fmt.Sprintf("%.1f%%", metrics.AvgCPU), fmt.Sprintf("%.1f%%", metrics.MaxCPU)},
		{"Memory", fmt.Sprintf("%.1f MB", metrics.AvgMemory), fmt.Sprintf("%.1f MB", metrics.MaxMemory)},
	}

	pterm.DefaultTable.WithHasHeader().WithData(summaryData).Render()
	pterm.Println()

	// Analysis
	pterm.DefaultSection.Println("Analysis")

	cpuVariance := metrics.MaxCPU - metrics.AvgCPU
	if cpuVariance > 20 {
		pterm.Warning.Printf("  • High CPU variance (%.1f%%), process may be unstable\n", cpuVariance)
	} else {
		pterm.Success.Println("  • CPU usage is stable")
	}

	memVariance := metrics.MaxMemory - metrics.AvgMemory
	if memVariance > 50 {
		pterm.Warning.Printf("  • High memory variance (%.1f MB), potential memory leak\n", memVariance)
	} else {
		pterm.Success.Println("  • Memory usage is stable")
	}

	pterm.Println()
	pterm.Info.Printf("Monitored for: %v\n", duration)
	pterm.Info.Printf("Total samples: %d\n", len(metrics.Samples))

	return nil
}

// truncate trunca uma string
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
