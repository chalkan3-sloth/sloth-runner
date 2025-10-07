package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
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
	// Print header
	pterm.DefaultHeader.WithFullWidth().Println(fmt.Sprintf("📊 Agent Dashboard - %s", agentName))
	fmt.Println()

	// Get agent info from registry
	agentResp, err := client.GetAgentInfo(ctx, &pb.GetAgentInfoRequest{AgentName: agentName})
	if err != nil || !agentResp.Success {
		return fmt.Errorf("agent not found: %s", agentName)
	}

	agentInfo := agentResp.AgentInfo

	// Parse system info
	var sysInfo SystemInfo
	if agentInfo.SystemInfoJson != "" {
		if err := json.Unmarshal([]byte(agentInfo.SystemInfoJson), &sysInfo); err != nil {
			pterm.Warning.Printf("Failed to parse system info: %v\n", err)
		}
	}

	// Get current metrics
	metrics, err := getCurrentMetrics(ctx, client, agentName)
	if err != nil {
		pterm.Warning.Printf("Failed to get current metrics: %v\n", err)
	}

	// Display agent status
	displayAgentStatus(agentInfo, sysInfo)
	fmt.Println()

	// Display system information
	displaySystemInfo(sysInfo)
	fmt.Println()

	// Display current metrics
	if metrics != nil {
		displayCurrentMetrics(metrics)
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

