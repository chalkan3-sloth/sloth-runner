package resources

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewResourcesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resources",
		Aliases: []string{"resource", "res"},
		Short:   "Monitor system resources on agents",
		Long: `Monitor CPU, memory, disk, and network usage on remote agents.
Set thresholds for alerts and track resource usage over time.`,
		Example: `  # Overview of all resources
  sloth-runner sysadmin resources overview

  # CPU usage details
  sloth-runner sysadmin resources cpu

  # Memory usage
  sloth-runner sysadmin resources memory

  # Disk usage
  sloth-runner sysadmin resources disk

  # Top consumers
  sloth-runner sysadmin resources top`,
	}

	// overview subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "overview",
		Short: "Show resource overview",
		Long:  `Display overview of CPU, memory, disk, and network resources.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOverview()
		},
	})

	// cpu subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "cpu",
		Short: "Show CPU usage",
		Long:  `Display detailed CPU usage information including per-core usage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCPU()
		},
	})

	// memory subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "memory",
		Short: "Show memory usage",
		Long:  `Display memory usage including RAM, swap, and buffers/cache.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMemory()
		},
	})

	// disk subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "disk",
		Short: "Show disk usage",
		Long:  `Display disk space usage for all mounted filesystems.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDisk()
		},
	})

	// io subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "io",
		Short: "Show disk I/O statistics",
		Long:  `Display disk I/O statistics including read/write operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Disk I/O monitoring not yet implemented")
			pterm.Info.Println("Future features: Read/write throughput, IOPS, latency")
		},
	})

	// network subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "network",
		Short: "Show network statistics",
		Long:  `Display network interface statistics including bandwidth usage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNetwork()
		},
	})

	// check subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check resources against thresholds",
		Long:  `Check resource usage against configured thresholds and alert if exceeded.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Resource threshold checking not yet implemented")
			pterm.Info.Println("Future features: Configurable thresholds, alerting, multi-agent checks")
		},
	})

	// history subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "history",
		Short: "Show resource usage history",
		Long:  `Display historical resource usage data.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Warning.Println("⚠️  Resource history not yet implemented")
			pterm.Info.Println("Future features: Time-series data, graphs, trend analysis")
		},
	})

	// top subcommand
	topCmd := &cobra.Command{
		Use:   "top",
		Short: "Show top resource consumers",
		Long:  `Display processes consuming the most resources (like htop).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, _ := cmd.Flags().GetInt("limit")
			return runTop(limit)
		},
	}
	topCmd.Flags().IntP("limit", "n", 10, "Number of processes to show")
	cmd.AddCommand(topCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the resources command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showResourcesDocs()
		},
	})

	return cmd
}

func runOverview() error {
	monitor := NewMonitor()

	spinner, _ := pterm.DefaultSpinner.Start("Collecting resource information...")

	// Get all stats
	cpu, err := monitor.GetCPU()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get CPU stats: %s", err.Error()))
		return err
	}

	memory, err := monitor.GetMemory()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get memory stats: %s", err.Error()))
		return err
	}

	disks, err := monitor.GetDisk()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get disk stats: %s", err.Error()))
		return err
	}

	networks, err := monitor.GetNetwork()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get network stats: %s", err.Error()))
		return err
	}

	spinner.Success("✅ Resource information collected")
	fmt.Println()

	// Display overview
	pterm.DefaultHeader.WithFullWidth().Println("System Resource Overview")
	fmt.Println()

	// CPU Section
	pterm.DefaultSection.Println("CPU")
	cpuColor := getUsageColor(cpu.Usage)
	pterm.Printf("  Usage: %s\n", cpuColor.Sprintf("%.1f%%", cpu.Usage))
	pterm.Printf("  Cores: %d\n", cpu.Cores)
	pterm.Printf("  Load Avg: %.2f, %.2f, %.2f (1m, 5m, 15m)\n", cpu.LoadAverage[0], cpu.LoadAverage[1], cpu.LoadAverage[2])
	fmt.Println()

	// Memory Section
	pterm.DefaultSection.Println("Memory")
	memColor := getUsageColor(memory.UsagePercent)
	pterm.Printf("  Total: %s\n", FormatBytes(memory.Total))
	pterm.Printf("  Used: %s (%s)\n", FormatBytes(memory.Used), memColor.Sprintf("%.1f%%", memory.UsagePercent))
	pterm.Printf("  Free: %s\n", FormatBytes(memory.Free))
	pterm.Printf("  Available: %s\n", FormatBytes(memory.Available))
	if memory.SwapTotal > 0 {
		swapUsedPercent := 0.0
		if memory.SwapTotal > 0 {
			swapUsedPercent = 100.0 * float64(memory.SwapUsed) / float64(memory.SwapTotal)
		}
		swapColor := getUsageColor(swapUsedPercent)
		pterm.Printf("  Swap: %s / %s (%s)\n", FormatBytes(memory.SwapUsed), FormatBytes(memory.SwapTotal), swapColor.Sprintf("%.1f%%", swapUsedPercent))
	}
	fmt.Println()

	// Disk Section
	pterm.DefaultSection.Println("Disk")
	for _, disk := range disks {
		diskColor := getUsageColor(disk.UsagePercent)
		pterm.Printf("  %s (%s)\n", disk.MountPoint, disk.Filesystem)
		pterm.Printf("    Usage: %s / %s (%s)\n", FormatBytes(disk.Used), FormatBytes(disk.Total), diskColor.Sprintf("%.1f%%", disk.UsagePercent))
	}
	fmt.Println()

	// Network Section
	pterm.DefaultSection.Println("Network")
	for _, net := range networks {
		pterm.Printf("  %s\n", net.Interface)
		pterm.Printf("    RX: %s", FormatBytes(net.BytesRecv))
		if net.PacketsRecv > 0 {
			pterm.Printf(" (%d packets)", net.PacketsRecv)
		}
		fmt.Println()
		pterm.Printf("    TX: %s", FormatBytes(net.BytesSent))
		if net.PacketsSent > 0 {
			pterm.Printf(" (%d packets)", net.PacketsSent)
		}
		fmt.Println()
	}

	return nil
}

func runCPU() error {
	monitor := NewMonitor()

	spinner, _ := pterm.DefaultSpinner.Start("Collecting CPU information...")

	cpu, err := monitor.GetCPU()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get CPU stats: %s", err.Error()))
		return err
	}

	spinner.Success("✅ CPU information collected")
	fmt.Println()

	// Display CPU info
	pterm.DefaultHeader.WithFullWidth().Println("CPU Information")
	fmt.Println()

	// Overall usage
	cpuColor := getUsageColor(cpu.Usage)
	pterm.DefaultBox.WithTitle("Overall CPU Usage").WithTitleTopCenter().Println(
		cpuColor.Sprintf("%.1f%%", cpu.Usage),
	)
	fmt.Println()

	// Details
	pterm.DefaultSection.Println("Details")
	pterm.Printf("  Cores: %d\n", cpu.Cores)
	pterm.Printf("  Load Average:\n")
	pterm.Printf("    1 min:  %.2f\n", cpu.LoadAverage[0])
	pterm.Printf("    5 min:  %.2f\n", cpu.LoadAverage[1])
	pterm.Printf("    15 min: %.2f\n", cpu.LoadAverage[2])
	fmt.Println()

	// Per-core usage (se disponível)
	if len(cpu.PerCore) > 0 {
		pterm.DefaultSection.Println("Per-Core Usage")
		for i, usage := range cpu.PerCore {
			coreColor := getUsageColor(usage)
			pterm.Printf("  Core %2d: %s\n", i, coreColor.Sprintf("%.1f%%", usage))
		}
	}

	return nil
}

func runMemory() error {
	monitor := NewMonitor()

	spinner, _ := pterm.DefaultSpinner.Start("Collecting memory information...")

	memory, err := monitor.GetMemory()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get memory stats: %s", err.Error()))
		return err
	}

	spinner.Success("✅ Memory information collected")
	fmt.Println()

	// Display memory info
	pterm.DefaultHeader.WithFullWidth().Println("Memory Information")
	fmt.Println()

	// Usage box
	memColor := getUsageColor(memory.UsagePercent)
	pterm.DefaultBox.WithTitle("Memory Usage").WithTitleTopCenter().Println(
		fmt.Sprintf("%s / %s (%s)", FormatBytes(memory.Used), FormatBytes(memory.Total), memColor.Sprintf("%.1f%%", memory.UsagePercent)),
	)
	fmt.Println()

	// RAM Details
	pterm.DefaultSection.Println("RAM")
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Total", FormatBytes(memory.Total)},
		{"Used", FormatBytes(memory.Used)},
		{"Free", FormatBytes(memory.Free)},
		{"Available", FormatBytes(memory.Available)},
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	fmt.Println()

	// Swap Details
	if memory.SwapTotal > 0 {
		pterm.DefaultSection.Println("Swap")
		swapUsedPercent := 0.0
		if memory.SwapTotal > 0 {
			swapUsedPercent = 100.0 * float64(memory.SwapUsed) / float64(memory.SwapTotal)
		}
		swapColor := getUsageColor(swapUsedPercent)
		swapTableData := pterm.TableData{
			{"Metric", "Value"},
			{"Total", FormatBytes(memory.SwapTotal)},
			{"Used", fmt.Sprintf("%s (%s)", FormatBytes(memory.SwapUsed), swapColor.Sprintf("%.1f%%", swapUsedPercent))},
			{"Free", FormatBytes(memory.SwapFree)},
		}
		pterm.DefaultTable.WithHasHeader().WithData(swapTableData).Render()
	} else {
		pterm.Info.Println("No swap configured")
	}

	return nil
}

func runDisk() error {
	monitor := NewMonitor()

	spinner, _ := pterm.DefaultSpinner.Start("Collecting disk information...")

	disks, err := monitor.GetDisk()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get disk stats: %s", err.Error()))
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d filesystems", len(disks)))
	fmt.Println()

	// Display disk info
	pterm.DefaultHeader.WithFullWidth().Println("Disk Usage")
	fmt.Println()

	// Table with all disks
	tableData := pterm.TableData{
		{"Filesystem", "Mount Point", "Total", "Used", "Available", "Usage %"},
	}

	for _, disk := range disks {
		usageColor := getUsageColor(disk.UsagePercent)
		tableData = append(tableData, []string{
			disk.Filesystem,
			disk.MountPoint,
			FormatBytes(disk.Total),
			FormatBytes(disk.Used),
			FormatBytes(disk.Available),
			usageColor.Sprintf("%.1f%%", disk.UsagePercent),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}

func runNetwork() error {
	monitor := NewMonitor()

	spinner, _ := pterm.DefaultSpinner.Start("Collecting network information...")

	networks, err := monitor.GetNetwork()
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get network stats: %s", err.Error()))
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d network interfaces", len(networks)))
	fmt.Println()

	// Display network info
	pterm.DefaultHeader.WithFullWidth().Println("Network Statistics")
	fmt.Println()

	for _, net := range networks {
		pterm.DefaultSection.Printf("Interface: %s", net.Interface)
		fmt.Println()

		tableData := pterm.TableData{
			{"Metric", "Received", "Sent"},
			{"Bytes", FormatBytes(net.BytesRecv), FormatBytes(net.BytesSent)},
		}

		if net.PacketsRecv > 0 || net.PacketsSent > 0 {
			tableData = append(tableData, []string{"Packets", fmt.Sprintf("%d", net.PacketsRecv), fmt.Sprintf("%d", net.PacketsSent)})
		}

		if net.ErrorsRecv > 0 || net.ErrorsSent > 0 {
			tableData = append(tableData, []string{"Errors", pterm.FgRed.Sprintf("%d", net.ErrorsRecv), pterm.FgRed.Sprintf("%d", net.ErrorsSent)})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		fmt.Println()
	}

	return nil
}

func runTop(limit int) error {
	monitor := NewMonitor()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Finding top %d processes...", limit))

	processes, err := monitor.GetProcesses(limit)
	if err != nil {
		spinner.Fail(fmt.Sprintf("❌ Failed to get process stats: %s", err.Error()))
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d processes", len(processes)))
	fmt.Println()

	// Display top processes
	pterm.DefaultHeader.WithFullWidth().Printf("Top %d Resource Consumers", limit)
	fmt.Println()

	tableData := pterm.TableData{
		{"PID", "Name", "CPU %", "Memory %"},
	}

	for _, proc := range processes {
		cpuColor := getUsageColor(proc.CPUPercent)
		memColor := getUsageColor(proc.MemoryPercent)

		tableData = append(tableData, []string{
			fmt.Sprintf("%d", proc.PID),
			proc.Name,
			cpuColor.Sprintf("%.1f%%", proc.CPUPercent),
			memColor.Sprintf("%.1f%%", proc.MemoryPercent),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}

// Helper function to get color based on usage percentage
func getUsageColor(usage float64) pterm.Color {
	if usage >= 90 {
		return pterm.FgRed
	} else if usage >= 70 {
		return pterm.FgYellow
	}
	return pterm.FgGreen
}

// showResourcesDocs and showDocs functions remain unchanged
func showResourcesDocs() {
	title := "SLOTH-RUNNER SYSADMIN RESOURCES(1)"
	description := "sloth-runner sysadmin resources - Monitor system resources on remote agents"
	synopsis := "sloth-runner sysadmin resources [subcommand] [options]"

	options := [][]string{
		{"overview", "Display overview of CPU, memory, disk, and network resources."},
		{"cpu", "Show detailed CPU usage including per-core usage and load average."},
		{"memory", "Display memory usage including RAM, swap, buffers, and cache."},
		{"disk", "Show disk space usage for all mounted filesystems."},
		{"io", "Display disk I/O statistics including read/write operations and IOPS."},
		{"network", "Show network interface statistics including bandwidth and packets."},
		{"check", "Check resource usage against configured thresholds and alert if exceeded."},
		{"history", "Display historical resource usage data with trend analysis."},
		{"top", "Show processes consuming the most resources (htop-like interface)."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Resource overview",
			"sloth-runner sysadmin resources overview",
			"Shows summary of all resource usage",
		},
		{
			"CPU details",
			"sloth-runner sysadmin res cpu",
			"Displays per-core CPU usage and load average",
		},
		{
			"Memory usage",
			"sloth-runner sysadmin resources memory",
			"Shows RAM, swap, buffers, and cache usage",
		},
		{
			"Disk usage",
			"sloth-runner sysadmin resources disk",
			"Shows disk usage for all filesystems",
		},
		{
			"Network statistics",
			"sloth-runner sysadmin res network",
			"Shows network interface statistics",
		},
		{
			"Top consumers",
			"sloth-runner sysadmin resources top --limit 20",
			"Shows top 20 CPU-consuming processes",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin performance - Performance monitoring",
		"sloth-runner sysadmin health - Health checks",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin resources --help")
}
