package telemetry

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// MetricsData holds parsed Prometheus metrics
type MetricsData struct {
	TasksTotal          map[string]float64
	TasksRunning        float64
	TaskDurationP50     map[string]float64
	TaskDurationP99     map[string]float64
	GRPCRequestsTotal   map[string]float64
	GRPCDurationP50     map[string]float64
	AgentUptime         float64
	AgentVersion        string
	AgentOS             string
	AgentArch           string
	Goroutines          float64
	MemoryAllocated     float64
	ErrorsTotal         map[string]float64
	Timestamp           time.Time
}

// FetchMetrics fetches and parses Prometheus metrics from an endpoint
func FetchMetrics(endpoint string) (*MetricsData, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics endpoint returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return parseMetrics(string(body))
}

// parseMetrics parses Prometheus text format
func parseMetrics(content string) (*MetricsData, error) {
	data := &MetricsData{
		TasksTotal:        make(map[string]float64),
		TaskDurationP50:   make(map[string]float64),
		TaskDurationP99:   make(map[string]float64),
		GRPCRequestsTotal: make(map[string]float64),
		GRPCDurationP50:   make(map[string]float64),
		ErrorsTotal:       make(map[string]float64),
		Timestamp:         time.Now(),
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		metricLine := parts[0]
		valueStr := parts[1]
		value, _ := strconv.ParseFloat(valueStr, 64)

		// Parse metric name and labels
		var metricName string
		labels := make(map[string]string)

		if strings.Contains(metricLine, "{") {
			idx := strings.Index(metricLine, "{")
			metricName = metricLine[:idx]
			labelsPart := metricLine[idx+1 : len(metricLine)-1]

			for _, labelPair := range strings.Split(labelsPart, ",") {
				kv := strings.Split(labelPair, "=")
				if len(kv) == 2 {
					key := strings.TrimSpace(kv[0])
					val := strings.Trim(strings.TrimSpace(kv[1]), "\"")
					labels[key] = val
				}
			}
		} else {
			metricName = metricLine
		}

		// Store metrics
		switch metricName {
		case "sloth_tasks_total":
			status := labels["status"]
			group := labels["group"]
			key := fmt.Sprintf("%s:%s", group, status)
			data.TasksTotal[key] = value

		case "sloth_tasks_running":
			data.TasksRunning = value

		case "sloth_task_duration_seconds":
			if strings.Contains(metricLine, "quantile=\"0.5\"") {
				group := labels["group"]
				task := labels["task"]
				key := fmt.Sprintf("%s:%s", group, task)
				data.TaskDurationP50[key] = value
			} else if strings.Contains(metricLine, "quantile=\"0.99\"") {
				group := labels["group"]
				task := labels["task"]
				key := fmt.Sprintf("%s:%s", group, task)
				data.TaskDurationP99[key] = value
			}

		case "sloth_grpc_requests_total":
			method := labels["method"]
			status := labels["status"]
			key := fmt.Sprintf("%s:%s", method, status)
			data.GRPCRequestsTotal[key] = value

		case "sloth_grpc_request_duration_seconds":
			if strings.Contains(metricLine, "quantile=\"0.5\"") {
				method := labels["method"]
				data.GRPCDurationP50[method] = value
			}

		case "sloth_agent_uptime_seconds":
			data.AgentUptime = value

		case "sloth_agent_info":
			data.AgentVersion = labels["version"]
			data.AgentOS = labels["os"]
			data.AgentArch = labels["arch"]

		case "sloth_goroutines":
			data.Goroutines = value

		case "sloth_memory_allocated_bytes":
			data.MemoryAllocated = value

		case "sloth_errors_total":
			errorType := labels["type"]
			data.ErrorsTotal[errorType] = value
		}
	}

	return data, nil
}

// DisplayDashboard displays a comprehensive dashboard of metrics including system monitoring
func DisplayDashboard(data *MetricsData, agentName string) {
	// Collect system metrics
	systemMetrics, err := CollectSystemMetrics()
	if err != nil {
		pterm.Warning.Printf("Failed to collect some system metrics: %v\n", err)
	}

	// Clear screen
	pterm.Println()

	// Header
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		Println(fmt.Sprintf("📊 Sloth Runner Complete System Monitor - Agent: %s", agentName))
	pterm.Println()

	// System Overview Section
	if systemMetrics != nil {
		displaySystemOverview(systemMetrics)
	}

	// CPU Metrics Section
	if systemMetrics != nil {
		displayCPUMetrics(systemMetrics)
	}

	// Memory Metrics Section
	if systemMetrics != nil {
		displayMemoryMetrics(systemMetrics)
	}

	// Disk Metrics Section
	if systemMetrics != nil && len(systemMetrics.DiskMetrics) > 0 {
		displayDiskMetrics(systemMetrics)
	}

	// Network Metrics Section
	if systemMetrics != nil && len(systemMetrics.NetworkInterfaces) > 0 {
		displayNetworkMetrics(systemMetrics)
	}

	// Process Metrics Section
	if systemMetrics != nil && len(systemMetrics.TopProcesses) > 0 {
		displayProcessMetrics(systemMetrics)
	}

	// Sloth Runner Agent Section
	pterm.DefaultSection.WithLevel(2).Println("🦥 Sloth Runner Agent")
	displaySlothRunnerMetrics(data)

	// Summary Box
	displaySummary(data, systemMetrics)
}

func displaySystemOverview(m *SystemMetrics) {
	pterm.DefaultSection.Println("🖥️  System Overview")

	systemData := [][]string{
		{"Hostname", m.Hostname},
		{"OS", fmt.Sprintf("%s %s", m.Platform, m.PlatformVersion)},
		{"Kernel", fmt.Sprintf("%s (%s)", m.KernelVersion, m.KernelArch)},
		{"Uptime", FormatUptime(m.Uptime)},
		{"Boot Time", time.Unix(int64(m.BootTime), 0).Format("2006-01-02 15:04:05")},
		{"Processes", fmt.Sprintf("%d (Zombies: %d)", m.ProcessCount, m.ZombieCount)},
		{"Network Connections", fmt.Sprintf("%d", m.NetworkConnections)},
	}

	pterm.DefaultTable.WithHasHeader(false).WithBoxed(true).WithData(systemData).Render()
	pterm.Println()
}

func displayCPUMetrics(m *SystemMetrics) {
	pterm.DefaultSection.Println("🔥 CPU Metrics")

	// CPU Info
	cpuInfoData := [][]string{
		{"Model", m.CPUModel},
		{"Cores", fmt.Sprintf("%d", m.CPUCores)},
		{"Threads", fmt.Sprintf("%d", m.CPUThreads)},
		{"Speed", fmt.Sprintf("%.2f GHz", m.CPUSpeed/1000)},
		{"Load Average", fmt.Sprintf("%.2f, %.2f, %.2f", m.LoadAverage1, m.LoadAverage5, m.LoadAverage15)},
	}

	pterm.DefaultTable.WithHasHeader(false).WithBoxed(true).WithData(cpuInfoData).Render()

	// Overall CPU Usage
	cpuBar := createColoredProgressBar("CPU Usage", m.CPUUsageTotal, 100, "percent")
	pterm.Println(cpuBar)

	// Per-core usage (if available)
	if len(m.CPUUsagePerCore) > 0 && len(m.CPUUsagePerCore) <= 16 { // Show only if reasonable number of cores
		pterm.Println()
		pterm.DefaultBox.WithTitle("CPU Cores Usage").WithTitleTopCenter().Println(
			createMultiCoreDisplay(m.CPUUsagePerCore),
		)
	}
	pterm.Println()
}

func displayMemoryMetrics(m *SystemMetrics) {
	pterm.DefaultSection.Println("💾 Memory Metrics")

	// Memory bars
	memBar := createColoredProgressBar("RAM Usage", m.MemoryUsedPct, 100, "percent")
	pterm.Println(memBar)

	swapBar := createColoredProgressBar("Swap Usage", m.SwapUsedPct, 100, "percent")
	pterm.Println(swapBar)

	// Memory details table
	memData := [][]string{
		{"Type", "Total", "Used", "Free", "Available"},
		{"RAM", FormatBytes(m.MemoryTotal), FormatBytes(m.MemoryUsed),
			FormatBytes(m.MemoryFree), FormatBytes(m.MemoryAvailable)},
		{"Swap", FormatBytes(m.SwapTotal), FormatBytes(m.SwapUsed),
			FormatBytes(m.SwapFree), "-"},
	}

	pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
		WithBoxed(true).WithData(memData).Render()
	pterm.Println()
}

func displayDiskMetrics(m *SystemMetrics) {
	pterm.DefaultSection.Println("💿 Disk Metrics")

	diskData := [][]string{
		{"Mount", "Device", "Type", "Size", "Used", "Free", "Use%"},
	}

	for _, disk := range m.DiskMetrics {
		// Skip very small partitions
		if disk.Total < 1024*1024*100 { // Less than 100MB
			continue
		}

		usageColor := pterm.FgGreen
		if disk.UsedPercent > 80 {
			usageColor = pterm.FgRed
		} else if disk.UsedPercent > 60 {
			usageColor = pterm.FgYellow
		}

		diskData = append(diskData, []string{
			disk.MountPoint,
			disk.Device,
			disk.Fstype,
			FormatBytes(disk.Total),
			FormatBytes(disk.Used),
			FormatBytes(disk.Free),
			pterm.NewStyle(usageColor).Sprintf("%.1f%%", disk.UsedPercent),
		})
	}

	if len(diskData) > 1 {
		pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
			WithBoxed(true).WithData(diskData).Render()
	}
	pterm.Println()
}

func displayNetworkMetrics(m *SystemMetrics) {
	pterm.DefaultSection.Println("🌐 Network Metrics")

	netData := [][]string{
		{"Interface", "Sent", "Received", "Packets TX/RX", "Errors", "Drops"},
	}

	for _, net := range m.NetworkInterfaces {
		netData = append(netData, []string{
			net.Name,
			FormatBytes(net.BytesSent),
			FormatBytes(net.BytesRecv),
			fmt.Sprintf("%d/%d", net.PacketsSent, net.PacketsRecv),
			fmt.Sprintf("%d/%d", net.Errout, net.Errin),
			fmt.Sprintf("%d/%d", net.Dropout, net.Dropin),
		})
	}

	if len(netData) > 1 {
		pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
			WithBoxed(true).WithData(netData).Render()
	}
	pterm.Println()
}

func displayProcessMetrics(m *SystemMetrics) {
	pterm.DefaultSection.Println("📋 Top Processes")

	processData := [][]string{
		{"PID", "Name", "User", "CPU%", "Memory", "Status"},
	}

	for i, proc := range m.TopProcesses {
		if i >= 10 { // Show top 10
			break
		}

		cpuColor := pterm.FgGreen
		if proc.CPUPercent > 50 {
			cpuColor = pterm.FgRed
		} else if proc.CPUPercent > 20 {
			cpuColor = pterm.FgYellow
		}

		statusColor := pterm.FgGreen
		if proc.Status == "Z" {
			statusColor = pterm.FgRed
		} else if proc.Status == "T" {
			statusColor = pterm.FgYellow
		}

		// Truncate long names
		name := proc.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		processData = append(processData, []string{
			fmt.Sprintf("%d", proc.PID),
			name,
			proc.Username,
			pterm.NewStyle(cpuColor).Sprintf("%.1f%%", proc.CPUPercent),
			fmt.Sprintf("%.1f MB", proc.MemoryMB),
			pterm.NewStyle(statusColor).Sprint(proc.Status),
		})
	}

	pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
		WithBoxed(true).WithData(processData).Render()
	pterm.Println()
}

func displaySlothRunnerMetrics(data *MetricsData) {
	// Agent Info
	agentInfoData := [][]string{
		{"Version", data.AgentVersion},
		{"Uptime", formatDuration(data.AgentUptime)},
		{"Goroutines", fmt.Sprintf("%.0f", data.Goroutines)},
		{"Memory", fmt.Sprintf("%.0f MB", data.MemoryAllocated/1024/1024)},
	}

	pterm.DefaultTable.WithHasHeader(false).WithBoxed(true).WithData(agentInfoData).Render()

	// Tasks Summary
	if len(data.TasksTotal) > 0 {
		taskSummary := summarizeTasks(data.TasksTotal)
		taskData := [][]string{
			{"Status", "Count"},
		}
		for status, count := range taskSummary {
			taskData = append(taskData, []string{
				formatStatus(status),
				fmt.Sprintf("%.0f", count),
			})
		}

		pterm.Println()
		pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
			WithBoxed(true).WithData(taskData).Render()
	}
	pterm.Println()
}

func displaySummary(data *MetricsData, systemMetrics *SystemMetrics) {
	summaryLines := []string{}

	if systemMetrics != nil {
		summaryLines = append(summaryLines,
			fmt.Sprintf("🖥️  CPU: %s | RAM: %s | Disk: %d mounts",
				pterm.Cyan(fmt.Sprintf("%.1f%%", systemMetrics.CPUUsageTotal)),
				pterm.Yellow(fmt.Sprintf("%.1f%%", systemMetrics.MemoryUsedPct)),
				len(systemMetrics.DiskMetrics)),
		)

		summaryLines = append(summaryLines,
			fmt.Sprintf("📊 Processes: %s | Network: %s | Uptime: %s",
				pterm.Green(fmt.Sprintf("%d", systemMetrics.ProcessCount)),
				pterm.Blue(fmt.Sprintf("%d conn", systemMetrics.NetworkConnections)),
				pterm.Magenta(FormatUptime(systemMetrics.Uptime))),
		)
	}

	// Sloth Runner summary
	totalTasks := float64(0)
	for _, count := range data.TasksTotal {
		totalTasks += count
	}

	summaryLines = append(summaryLines,
		fmt.Sprintf("🦥 Tasks: %s | Running: %s | Agent Memory: %s",
			pterm.Cyan(fmt.Sprintf("%.0f", totalTasks)),
			pterm.Yellow(fmt.Sprintf("%.0f", data.TasksRunning)),
			pterm.Green(fmt.Sprintf("%.0f MB", data.MemoryAllocated/1024/1024))),
	)

	pterm.DefaultBox.WithTitle("📈 System Summary").WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).
		Println(strings.Join(summaryLines, "\n"))
}

// Helper functions

func createMultiCoreDisplay(coreUsages []float64) string {
	lines := []string{}
	for i, usage := range coreUsages {
		bar := createMiniBar(usage)
		color := pterm.FgGreen
		if usage > 80 {
			color = pterm.FgRed
		} else if usage > 60 {
			color = pterm.FgYellow
		}

		lines = append(lines, fmt.Sprintf("Core %2d: %s %s",
			i,
			bar,
			pterm.NewStyle(color).Sprintf("%5.1f%%", usage),
		))

		// Two columns display
		if (i+1)%2 == 0 && i < len(coreUsages)-1 {
			lines[len(lines)-2] = lines[len(lines)-2] + "  |  " + lines[len(lines)-1]
			lines = lines[:len(lines)-1]
		}
	}
	return strings.Join(lines, "\n")
}

func createMiniBar(percentage float64) string {
	barWidth := 20
	filled := int(percentage / 100 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	color := pterm.FgGreen
	if percentage > 80 {
		color = pterm.FgRed
	} else if percentage > 60 {
		color = pterm.FgYellow
	}

	bar := strings.Repeat("█", filled)
	empty := strings.Repeat("░", barWidth-filled)

	return pterm.NewStyle(color).Sprint(bar) + empty
}

func createColoredProgressBar(label string, current, max float64, format string) string {
	if max == 0 {
		max = 1
	}

	percentage := current / max * 100
	if percentage > 100 {
		percentage = 100
	}

	barWidth := 40
	filled := int(percentage / 100 * float64(barWidth))

	bar := strings.Repeat("█", filled)
	empty := strings.Repeat("░", barWidth-filled)

	color := pterm.FgGreen
	if percentage > 80 {
		color = pterm.FgRed
	} else if percentage > 60 {
		color = pterm.FgYellow
	}

	coloredBar := pterm.NewStyle(color).Sprint(bar)

	var valueStr string
	if format == "percent" {
		valueStr = fmt.Sprintf("%.1f%%", percentage)
	} else {
		valueStr = fmt.Sprintf("%.0f/%.0f", current, max)
	}

	return fmt.Sprintf("%s: [%s%s] %s",
		pterm.Bold.Sprint(label),
		coloredBar,
		empty,
		valueStr,
	)
}

func formatDuration(seconds float64) string {
	duration := time.Duration(seconds) * time.Second

	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

func formatStatus(status string) string {
	switch status {
	case "success":
		return pterm.Green("✓ Success")
	case "failed":
		return pterm.Red("✗ Failed")
	case "skipped":
		return pterm.Yellow("⊘ Skipped")
	default:
		return status
	}
}

func summarizeTasks(tasksTotal map[string]float64) map[string]float64 {
	summary := make(map[string]float64)

	for key, count := range tasksTotal {
		parts := strings.Split(key, ":")
		if len(parts) > 1 {
			status := parts[1]
			summary[status] += count
		}
	}

	return summary
}

func createProgressBar(label string, current, max int) string {
	if max == 0 {
		max = 1
	}

	percentage := float64(current) / float64(max) * 100
	if percentage > 100 {
		percentage = 100
	}

	barWidth := 40
	filled := int(percentage / 100 * float64(barWidth))

	bar := strings.Repeat("█", filled)
	empty := strings.Repeat("░", barWidth-filled)

	color := pterm.FgGreen
	if percentage > 80 {
		color = pterm.FgRed
	} else if percentage > 60 {
		color = pterm.FgYellow
	}

	coloredBar := pterm.NewStyle(color).Sprint(bar)

	return fmt.Sprintf("%s: [%s%s] %d/%d (%.1f%%)",
		pterm.Bold.Sprint(label),
		coloredBar,
		empty,
		current,
		max,
		percentage,
	)
}

// DisplayHistoricalTrends displays metrics trends over time
func DisplayHistoricalTrends(history []*MetricsData, agentName string) {
	if len(history) == 0 {
		pterm.Warning.Println("No historical data available")
		return
	}

	pterm.DefaultHeader.WithFullWidth().Println(
		fmt.Sprintf("📈 Metrics Trends - Agent: %s", agentName),
	)
	pterm.Println()

	// Sort by timestamp
	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp.Before(history[j].Timestamp)
	})

	// Memory trend
	pterm.DefaultSection.Println("Memory Usage Trend")
	for i, data := range history {
		memoryMB := data.MemoryAllocated / 1024 / 1024
		timestamp := data.Timestamp.Format("15:04:05")

		barLength := int(memoryMB / 10)
		if barLength > 50 {
			barLength = 50
		}

		bar := strings.Repeat("▓", barLength)
		fmt.Printf("%s │ %s %.0f MB\n", timestamp, bar, memoryMB)

		if i < len(history)-1 && i%5 == 4 {
			pterm.Println()
		}
	}
}

// ContinuousMonitor runs continuous monitoring with refresh
func ContinuousMonitor(endpoint string, agentName string, refreshInterval time.Duration) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Fetch and display metrics
		data, err := FetchMetrics(endpoint)
		if err != nil {
			pterm.Error.Printf("Failed to fetch metrics: %v\n", err)
		} else {
			DisplayDashboard(data, agentName)
		}

		// Wait for next refresh
		<-ticker.C
	}
}