package performance

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewPerformanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "performance",
		Short:   "Monitor and analyze system performance",
		Long:    `Track CPU, memory, disk I/O, and network performance metrics.`,
		Aliases: []string{"perf"},
	}

	// show command
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current performance metrics",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runShowPerformance(); err != nil {
				pterm.Error.Printf("Failed to show performance: %v\n", err)
			}
		},
	}
	cmd.AddCommand(showCmd)

	// monitor command
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor performance over time",
		Run: func(cmd *cobra.Command, args []string) {
			duration, _ := cmd.Flags().GetDuration("duration")

			if err := runMonitorPerformance(duration); err != nil {
				pterm.Error.Printf("Failed to monitor performance: %v\n", err)
			}
		},
	}
	monitorCmd.Flags().DurationP("duration", "d", 10*time.Second, "Monitoring duration")
	cmd.AddCommand(monitorCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the performance command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showPerformanceDocs()
		},
	})

	return cmd
}

func showPerformanceDocs() {
	title := "SLOTH-RUNNER SYSADMIN PERFORMANCE(1)"
	description := "sloth-runner sysadmin performance - Monitor and analyze system performance"
	synopsis := "sloth-runner sysadmin performance [subcommand] [options]"

	options := [][]string{
		{"show", "Show current performance metrics including CPU usage, memory stats, disk I/O, and network throughput."},
		{"monitor", "Monitor performance in real-time with live dashboards, alert thresholds, and historical trends."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"View current metrics",
			"sloth-runner sysadmin performance show --agent web-01",
			"Displays current CPU, memory, disk, and network metrics",
		},
		{
			"Continuous monitoring",
			"sloth-runner sysadmin perf monitor --interval 5s --all-agents",
			"Monitors all agents, updating every 5 seconds",
		},
		{
			"Performance snapshot",
			"sloth-runner sysadmin performance show --agent db-01 --json",
			"Exports performance metrics in JSON format",
		},
		{
			"Alert on thresholds",
			"sloth-runner sysadmin perf monitor --alert-if cpu>80 memory>90",
			"Alerts when CPU exceeds 80% or memory exceeds 90%",
		},
		{
			"Historical analysis",
			"sloth-runner sysadmin performance show --agent web-01 --since 24h",
			"Shows performance trends over the last 24 hours",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin resources - Resource monitoring",
		"sloth-runner sysadmin health - Health checks",
		"sloth-runner agent metrics - Agent metrics",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin performance --help")
}

// runShowPerformance mostra métricas de performance
func runShowPerformance() error {
	collector := NewCollector()

	spinner, _ := pterm.DefaultSpinner.Start("Collecting performance metrics...")

	metrics, err := collector.CollectMetrics()
	if err != nil {
		spinner.Fail("Failed to collect metrics")
		return err
	}

	spinner.Success("✅ Performance metrics collected")
	pterm.Println()

	// Overall Performance Score
	pterm.DefaultHeader.WithFullWidth().Println("System Performance Report")
	pterm.Println()

	// Overall score box
	scoreColor := getScoreColor(metrics.Overall.Status)
	pterm.DefaultBox.WithTitle("Overall Performance").WithTitleTopCenter().Println(
		scoreColor.Sprintf("%s (Score: %d/100)", metrics.Overall.Status, metrics.Overall.Score),
	)
	pterm.Println()

	// Detailed metrics table
	tableData := pterm.TableData{
		{"Component", "Usage", "Status", "Details"},
	}

	// CPU
	if metrics.CPU != nil {
		tableData = append(tableData, []string{
			"CPU",
			fmt.Sprintf("%.1f%%", metrics.CPU.Usage),
			getStatusColor(metrics.CPU.Status).Sprint(metrics.CPU.Status),
			fmt.Sprintf("%d cores, Load: %.2f", metrics.CPU.Cores, metrics.CPU.LoadAverage[0]),
		})
	}

	// Memory
	if metrics.Memory != nil {
		tableData = append(tableData, []string{
			"Memory",
			fmt.Sprintf("%.1f%%", metrics.Memory.UsagePercent),
			getStatusColor(metrics.Memory.Status).Sprint(metrics.Memory.Status),
			fmt.Sprintf("%s / %s", FormatBytes(metrics.Memory.Used), FormatBytes(metrics.Memory.Total)),
		})
	}

	// Disk
	if metrics.Disk != nil {
		tableData = append(tableData, []string{
			"Disk",
			fmt.Sprintf("%.1f%%", metrics.Disk.UsagePercent),
			getStatusColor(metrics.Disk.Status).Sprint(metrics.Disk.Status),
			fmt.Sprintf("%s / %s free", FormatBytes(metrics.Disk.FreeSpace), FormatBytes(metrics.Disk.TotalSpace)),
		})
	}

	// Network
	if metrics.Network != nil {
		tableData = append(tableData, []string{
			"Network",
			fmt.Sprintf("%d active", metrics.Network.ActiveInterfaces),
			getStatusColor(metrics.Network.Status).Sprint(metrics.Network.Status),
			fmt.Sprintf("RX: %s, TX: %s", FormatBytes(metrics.Network.TotalBytesRecv), FormatBytes(metrics.Network.TotalBytesSent)),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Issues
	if len(metrics.Overall.Issues) > 0 {
		pterm.DefaultSection.Println("Performance Analysis")
		for _, issue := range metrics.Overall.Issues {
			if issue == "No issues detected" {
				pterm.Success.Println("  • " + issue)
			} else {
				pterm.Warning.Println("  • " + issue)
			}
		}
		pterm.Println()
	}

	// Recommendations
	if metrics.Overall.Score < 80 {
		pterm.DefaultSection.Println("Recommendations")
		if metrics.CPU != nil && metrics.CPU.Usage > 70 {
			pterm.Info.Println("  • Consider scaling CPU resources or optimizing workload")
		}
		if metrics.Memory != nil && metrics.Memory.UsagePercent > 80 {
			pterm.Info.Println("  • Consider adding more memory or optimizing memory usage")
		}
		if metrics.Disk != nil && metrics.Disk.UsagePercent > 80 {
			pterm.Info.Println("  • Clean up disk space or expand storage capacity")
		}
		pterm.Println()
	}

	return nil
}

// runMonitorPerformance monitora performance ao longo do tempo
func runMonitorPerformance(duration time.Duration) error {
	collector := NewCollector()

	pterm.Info.Printf("Monitoring performance for %v...\n", duration)
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Collecting samples over %v...", duration))

	sample, err := collector.CollectSample(duration)
	if err != nil {
		spinner.Fail("Failed to collect samples")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Collected %d samples", len(sample.Samples)))
	pterm.Println()

	pterm.DefaultHeader.WithFullWidth().Println("Performance Monitoring Report")
	pterm.Println()

	// Summary statistics
	summaryData := pterm.TableData{
		{"Metric", "Average", "Min", "Max", "Status"},
		{
			"CPU Usage",
			fmt.Sprintf("%.1f%%", sample.AverageCPU),
			fmt.Sprintf("%.1f%%", sample.MinCPU),
			fmt.Sprintf("%.1f%%", sample.MaxCPU),
			getStatusColor(getCPUStatus(sample.AverageCPU)).Sprint(getCPUStatus(sample.AverageCPU)),
		},
		{
			"Memory Usage",
			fmt.Sprintf("%.1f%%", sample.AverageRAM),
			fmt.Sprintf("%.1f%%", sample.MinRAM),
			fmt.Sprintf("%.1f%%", sample.MaxRAM),
			getStatusColor(getMemoryStatus(sample.AverageRAM)).Sprint(getMemoryStatus(sample.AverageRAM)),
		},
	}

	pterm.DefaultTable.WithHasHeader().WithData(summaryData).Render()
	pterm.Println()

	// Analysis
	pterm.DefaultSection.Println("Analysis")

	cpuVariance := sample.MaxCPU - sample.MinCPU
	ramVariance := sample.MaxRAM - sample.MinRAM

	if cpuVariance > 30 {
		pterm.Warning.Printf("  • High CPU variance (%.1f%%), workload may be unstable\n", cpuVariance)
	} else {
		pterm.Success.Println("  • CPU usage is stable")
	}

	if ramVariance > 20 {
		pterm.Warning.Printf("  • High memory variance (%.1f%%), potential memory leaks\n", ramVariance)
	} else {
		pterm.Success.Println("  • Memory usage is stable")
	}

	if sample.AverageCPU < 50 && sample.AverageRAM < 70 {
		pterm.Success.Println("  • System resources are well-balanced")
	}

	pterm.Println()
	pterm.Info.Printf("Monitored for: %v\n", duration)
	pterm.Info.Printf("Total samples: %d\n", len(sample.Samples))

	return nil
}

// getStatusColor retorna cor para status
func getStatusColor(status PerformanceStatus) pterm.Color {
	switch status {
	case StatusExcellent:
		return pterm.FgGreen
	case StatusGood:
		return pterm.FgCyan
	case StatusWarning:
		return pterm.FgYellow
	case StatusCritical:
		return pterm.FgRed
	default:
		return pterm.FgWhite
	}
}

// getScoreColor retorna cor para score
func getScoreColor(status PerformanceStatus) pterm.Color {
	return getStatusColor(status)
}
