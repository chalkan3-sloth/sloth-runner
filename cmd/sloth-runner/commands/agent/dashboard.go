package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SystemInfo represents agent system information
type SystemInfo struct {
	Hostname      string  `json:"hostname"`
	OS            string  `json:"os"`
	Arch          string  `json:"arch"`
	CPUCores      int     `json:"cpu_cores"`
	TotalMemoryGB float64 `json:"total_memory_gb"`
	Uptime        int64   `json:"uptime"` // Uptime in seconds
}

// MetricsInfo represents current metrics
type MetricsInfo struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	DiskPercent   float64 `json:"disk_percent"`
	LoadAvg1      float64 `json:"load_avg_1min"`
	LoadAvg5      float64 `json:"load_avg_5min"`
	LoadAvg15     float64 `json:"load_avg_15min"`
}

// displayAgentDashboard shows a comprehensive dashboard for an agent
func displayAgentDashboard(ctx context.Context, client AgentRegistryClient, agentName string) error {
	// Print header with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	pterm.DefaultHeader.WithFullWidth().Println(fmt.Sprintf("📊 Agent Dashboard - %s | %s", agentName, timestamp))
	fmt.Println()

	// Get agent info from registry
	agentResp, err := client.GetAgentInfo(ctx, &pb.GetAgentInfoRequest{AgentName: agentName})
	if err != nil || !agentResp.Success {
		return fmt.Errorf("agent not found: %s", agentName)
	}

	agentInfo := agentResp.AgentInfo

	// Connect to agent directly using pb.AgentClient
	conn, err := grpc.Dial(agentInfo.AgentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	agentClient := pb.NewAgentClient(conn)

	// Get detailed metrics from agent
	detailedResp, err := agentClient.GetDetailedMetrics(ctx, &pb.DetailedMetricsRequest{})
	if err != nil {
		pterm.Warning.Printf("Failed to get detailed metrics: %v\n", err)
	}

	processResp, _ := agentClient.GetProcessList(ctx, &pb.ProcessListRequest{})
	networkResp, _ := agentClient.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
	diskResp, _ := agentClient.GetDiskInfo(ctx, &pb.DiskInfoRequest{})

	// Parse system info
	var sysInfo SystemInfo
	if agentInfo.SystemInfoJson != "" {
		if err := json.Unmarshal([]byte(agentInfo.SystemInfoJson), &sysInfo); err != nil {
			pterm.Warning.Printf("Failed to parse system info: %v\n", err)
		}
	}

	// Display agent status
	displayAgentStatus(agentInfo, sysInfo)
	fmt.Println()

	// Display system information and metrics
	if detailedResp != nil {
		displayDetailedMetricsInfo(detailedResp)
		fmt.Println()
	} else if sysInfo.OS != "" {
		displaySystemInfo(sysInfo)
		fmt.Println()
	}

	// Display network details
	if networkResp != nil {
		displayNetworkInfo(networkResp)
		fmt.Println()
	}

	// Display disk details
	if diskResp != nil {
		displayDiskInfo(diskResp)
		fmt.Println()
	}

	// Display top processes
	if processResp != nil && len(processResp.Processes) > 0 {
		displayProcessListInfo(processResp)
		fmt.Println()
	}

	// Display last heartbeat
	lastHeartbeat := time.Unix(agentInfo.LastHeartbeat, 0)
	timeSince := time.Since(lastHeartbeat)

	heartbeatColor := pterm.FgGreen
	if timeSince > 30*time.Second {
		heartbeatColor = pterm.FgYellow
	}
	if timeSince > 60*time.Second {
		heartbeatColor = pterm.FgRed
	}

	pterm.Info.Printf("Last Heartbeat: %s (%s ago)\n",
		lastHeartbeat.Format("2006-01-02 15:04:05"),
		heartbeatColor.Sprintf("%s", formatDuration(timeSince)))

	return nil
}

// getCurrentMetrics fetches current system metrics from agent via metrics database
func getCurrentMetrics(ctx context.Context, client AgentRegistryClient, agentName string) (*MetricsInfo, error) {
	// For now, return nil - metrics will be fetched from the metrics database
	// This is a placeholder for future implementation
	return nil, nil
}

// displayAgentStatus shows agent connection status
func displayAgentStatus(agentInfo *pb.AgentInfo, sysInfo SystemInfo) {
	statusColor := pterm.FgGreen
	statusIcon := "✅"
	statusText := "ONLINE"

	if agentInfo.Status != "Active" {
		statusColor = pterm.FgRed
		statusIcon = "❌"
		statusText = "OFFLINE"
	}

	tableData := pterm.TableData{
		{"Agent Name", agentInfo.AgentName},
		{"Address", agentInfo.AgentAddress},
		{"Status", statusColor.Sprintf("%s %s", statusIcon, statusText)},
		{"Version", agentInfo.Version},
	}

	if sysInfo.Hostname != "" {
		tableData = append(tableData, []string{"Hostname", sysInfo.Hostname})
	}

	pterm.DefaultTable.WithHasHeader(false).WithBoxed(true).WithData(tableData).Render()
}

// displaySystemInfo shows system information
func displaySystemInfo(sysInfo SystemInfo) {
	if sysInfo.OS == "" {
		return
	}

	pterm.DefaultSection.Println("System Information")

	tableData := pterm.TableData{
		{"Operating System", fmt.Sprintf("%s (%s)", sysInfo.OS, sysInfo.Arch)},
		{"CPU Cores", fmt.Sprintf("%d cores", sysInfo.CPUCores)},
		{"Total Memory", fmt.Sprintf("%.2f GB", sysInfo.TotalMemoryGB)},
	}

	if sysInfo.Uptime > 0 {
		uptimeDuration := time.Duration(sysInfo.Uptime) * time.Second
		tableData = append(tableData, []string{"Uptime", formatDuration(uptimeDuration)})
	}

	pterm.DefaultTable.WithHasHeader(false).WithData(tableData).Render()
}

// displayCurrentMetrics shows current resource usage
func displayCurrentMetrics(metrics *MetricsInfo) {
	pterm.DefaultSection.Println("Current Metrics")

	// CPU Usage
	cpuColor := getMetricColor(metrics.CPUPercent)
	cpuBar := createBar(metrics.CPUPercent, 40)
	fmt.Printf("  %s %s %.1f%%\n", pterm.Cyan("CPU Usage:      "), cpuColor.Sprint(cpuBar), metrics.CPUPercent)

	// Memory Usage
	memColor := getMetricColor(metrics.MemoryPercent)
	memBar := createBar(metrics.MemoryPercent, 40)
	fmt.Printf("  %s %s %.1f%%\n", pterm.Cyan("Memory Usage:   "), memColor.Sprint(memBar), metrics.MemoryPercent)

	// Disk Usage
	diskColor := getMetricColor(metrics.DiskPercent)
	diskBar := createBar(metrics.DiskPercent, 40)
	fmt.Printf("  %s %s %.1f%%\n", pterm.Cyan("Disk Usage:     "), diskColor.Sprint(diskBar), metrics.DiskPercent)

	fmt.Println()

	// Load Average
	tableData := pterm.TableData{
		{"Load Average (1min)", fmt.Sprintf("%.2f", metrics.LoadAvg1)},
		{"Load Average (5min)", fmt.Sprintf("%.2f", metrics.LoadAvg5)},
		{"Load Average (15min)", fmt.Sprintf("%.2f", metrics.LoadAvg15)},
	}

	pterm.DefaultTable.WithHasHeader(false).WithData(tableData).Render()
}

// getMetricColor returns color based on metric value
func getMetricColor(value float64) pterm.Color {
	if value < 50 {
		return pterm.FgGreen
	} else if value < 80 {
		return pterm.FgYellow
	}
	return pterm.FgRed
}

// getStatusIcon returns status icon based on metric value
func getStatusIcon(value float64) string {
	if value < 50 {
		return pterm.Green("✓")
	} else if value < 80 {
		return pterm.Yellow("⚠")
	}
	return pterm.Red("✗")
}

// createBar creates a visual bar for metrics
func createBar(percent float64, width int) string {
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}

// formatDuration formats a duration into human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
	}
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%dd %dh", days, hours)
}

// displayDetailedMetricsInfo shows detailed metrics
func displayDetailedMetricsInfo(metrics *pb.DetailedMetricsResponse) {
	pterm.DefaultSection.Println("💻 System Information")

	uptime := time.Duration(metrics.UptimeSeconds) * time.Second
	tableData := pterm.TableData{
		{"OS Version", metrics.OsVersion},
		{"Kernel", metrics.KernelVersion},
		{"Uptime", formatDuration(uptime)},
		{"Process Count", fmt.Sprintf("%d", metrics.ProcessCount)},
	}
	pterm.DefaultTable.WithHasHeader(false).WithData(tableData).Render()

	fmt.Println()
	pterm.DefaultSection.Println("📊 Resource Usage")

	// CPU - calculate average from per-core usage
	var cpuPercent float64
	if len(metrics.Cpu.PerCoreUsage) > 0 {
		for _, usage := range metrics.Cpu.PerCoreUsage {
			cpuPercent += usage
		}
		cpuPercent /= float64(len(metrics.Cpu.PerCoreUsage))
	}
	cpuColor := getMetricColor(cpuPercent)
	cpuBar := createBar(cpuPercent, 40)
	cpuStatus := getStatusIcon(cpuPercent)
	fmt.Printf("  %s %s %s %.1f%%\n", pterm.Cyan("CPU Usage:      "), cpuColor.Sprint(cpuBar), cpuStatus, cpuPercent)

	// Memory Usage
	memPercent := metrics.Memory.Percent
	memColor := getMetricColor(memPercent)
	memBar := createBar(memPercent, 40)
	memStatus := getStatusIcon(memPercent)
	fmt.Printf("  %s %s %s %.1f%% (%s / %s)\n",
		pterm.Cyan("Memory Usage:   "),
		memColor.Sprint(memBar),
		memStatus,
		memPercent,
		formatBytes(metrics.Memory.UsedBytes),
		formatBytes(metrics.Memory.TotalBytes))

	fmt.Println()

	// Load Average
	loadData := pterm.TableData{
		{"Load Average (1min)", fmt.Sprintf("%.2f", metrics.LoadAvg_1Min)},
		{"Load Average (5min)", fmt.Sprintf("%.2f", metrics.LoadAvg_5Min)},
		{"Load Average (15min)", fmt.Sprintf("%.2f", metrics.LoadAvg_15Min)},
	}
	pterm.DefaultTable.WithHasHeader(false).WithData(loadData).Render()

	fmt.Println()
	pterm.DefaultSection.Println("🧠 Memory Details")
	memData := pterm.TableData{
		{"Total", formatBytes(metrics.Memory.TotalBytes)},
		{"Used", fmt.Sprintf("%s (%.1f%%)", formatBytes(metrics.Memory.UsedBytes), metrics.Memory.Percent)},
		{"Free", formatBytes(metrics.Memory.FreeBytes)},
		{"Available", formatBytes(metrics.Memory.AvailableBytes)},
		{"Cached", formatBytes(metrics.Memory.CachedBytes)},
		{"Buffers", formatBytes(metrics.Memory.BuffersBytes)},
	}
	pterm.DefaultTable.WithHasHeader(false).WithData(memData).Render()
}

// displayNetworkInfo shows network interface details
func displayNetworkInfo(networkResp *pb.NetworkInfoResponse) {
	pterm.DefaultSection.Println("🌐 Network Interfaces")

	if len(networkResp.Interfaces) == 0 {
		pterm.Info.Println("No network interfaces found")
		return
	}

	for _, iface := range networkResp.Interfaces {
		if !iface.IsUp {
			continue // Skip down interfaces
		}

		fmt.Printf("\n  %s %s\n", pterm.Cyan("Interface:"), pterm.Green(iface.Name))
		tableData := pterm.TableData{
			{"IP Addresses", strings.Join(iface.IpAddresses, ", ")},
			{"MAC Address", iface.MacAddress},
			{"Bytes Sent", formatBytes(iface.BytesSent)},
			{"Bytes Received", formatBytes(iface.BytesRecv)},
			{"Status", pterm.Green("UP")},
		}
		pterm.DefaultTable.WithHasHeader(false).WithData(tableData).Render()
	}
}

// displayDiskInfo shows disk partition details
func displayDiskInfo(diskResp *pb.DiskInfoResponse) {
	pterm.DefaultSection.Println("💾 Disk Partitions")

	if len(diskResp.Partitions) == 0 {
		pterm.Info.Println("No disk partitions found")
		return
	}

	for _, partition := range diskResp.Partitions {
		// Skip snap mounts (always show 100% and clutter output)
		if strings.HasPrefix(partition.Mountpoint, "/snap/") {
			continue
		}

		// Calculate actual percentage (fixing backend bug)
		var actualPercent float64
		if partition.TotalBytes > 0 {
			actualPercent = float64(partition.UsedBytes) / float64(partition.TotalBytes) * 100.0
		}

		color := getMetricColor(actualPercent)
		status := getStatusIcon(actualPercent)

		// Highlight critical partitions (< 10% free)
		mountDisplay := partition.Mountpoint
		if actualPercent > 90 {
			mountDisplay = pterm.Red(partition.Mountpoint) + " ⚠️  CRITICAL"
		} else {
			mountDisplay = pterm.Green(partition.Mountpoint)
		}

		fmt.Printf("\n  %s %s (%s)\n", pterm.Cyan("Partition:"), mountDisplay, partition.Device)

		tableData := pterm.TableData{
			{"Total", formatBytes(partition.TotalBytes)},
			{"Used", fmt.Sprintf("%s (%s%.1f%%%s) %s", formatBytes(partition.UsedBytes), color, actualPercent, pterm.Reset.Sprint(), status)},
			{"Free", formatBytes(partition.FreeBytes)},
			{"Filesystem", partition.Fstype},
		}
		pterm.DefaultTable.WithHasHeader(false).WithData(tableData).Render()
	}
}

// displayProcessListInfo shows top processes
func displayProcessListInfo(processResp *pb.ProcessListResponse) {
	pterm.DefaultSection.Println("⚙️  Top Processes (by CPU)")

	// Create table data
	tableData := pterm.TableData{
		{"PID", "Name", "CPU%", "Memory", "Status"},
	}

	for i, proc := range processResp.Processes {
		if i >= 10 { // Limit to top 10
			break
		}

		// Calculate memory in MB from percentage (assuming we have total memory)
		// For now, just show percentage with better formatting
		memDisplay := fmt.Sprintf("%.1f%%", proc.MemoryPercent)
		if proc.MemoryBytes > 0 {
			memDisplay = formatBytes(proc.MemoryBytes)
		}

		tableData = append(tableData, []string{
			fmt.Sprintf("%d", proc.Pid),
			truncateString(proc.Name, 30),
			fmt.Sprintf("%.1f%%", proc.CpuPercent),
			memDisplay,
			proc.Status,
		})
	}

	pterm.DefaultTable.WithHasHeader(true).WithData(tableData).Render()
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

