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

// DisplayDashboard displays a comprehensive dashboard of metrics
func DisplayDashboard(data *MetricsData, agentName string) {
	// Clear screen
	pterm.Println()

	// Header
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		Println(fmt.Sprintf("📊 Sloth Runner Metrics Dashboard - Agent: %s", agentName))
	pterm.Println()

	// Agent Info Section
	agentInfoData := [][]string{
		{"Version", data.AgentVersion},
		{"OS", data.AgentOS},
		{"Architecture", data.AgentArch},
		{"Uptime", formatDuration(data.AgentUptime)},
		{"Last Updated", data.Timestamp.Format("2006-01-02 15:04:05")},
	}

	pterm.DefaultSection.Println("🔧 Agent Information")
	pterm.DefaultTable.WithHasHeader(false).WithBoxed(true).WithData(agentInfoData).Render()
	pterm.Println()

	// System Resources Section
	pterm.DefaultSection.Println("💻 System Resources")

	memoryMB := data.MemoryAllocated / 1024 / 1024
	goroutinesBar := createProgressBar("Goroutines", int(data.Goroutines), 1000)
	memoryBar := createProgressBar("Memory (MB)", int(memoryMB), 512)

	pterm.Println(goroutinesBar)
	pterm.Println(memoryBar)
	pterm.Println()

	// Tasks Section
	pterm.DefaultSection.Println("📋 Task Metrics")

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

	if len(taskData) > 1 {
		pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
			WithBoxed(true).WithData(taskData).Render()
	} else {
		pterm.Info.Println("No tasks executed yet")
	}

	runningBar := createProgressBar("Running Tasks", int(data.TasksRunning), 10)
	pterm.Println(runningBar)
	pterm.Println()

	// Task Duration Section
	if len(data.TaskDurationP50) > 0 {
		pterm.DefaultSection.Println("⏱️  Task Performance")

		durationData := [][]string{
			{"Task", "P50 (ms)", "P99 (ms)", "Status"},
		}

		for key, p50 := range data.TaskDurationP50 {
			p99 := data.TaskDurationP99[key]
			parts := strings.Split(key, ":")
			taskName := key
			if len(parts) > 1 {
				taskName = parts[1]
			}

			status := "🟢 Fast"
			if p99 > 5.0 {
				status = "🔴 Slow"
			} else if p99 > 1.0 {
				status = "🟡 Normal"
			}

			durationData = append(durationData, []string{
				taskName,
				fmt.Sprintf("%.2f", p50*1000),
				fmt.Sprintf("%.2f", p99*1000),
				status,
			})
		}

		pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
			WithBoxed(true).WithData(durationData).Render()
		pterm.Println()
	}

	// gRPC Metrics Section
	if len(data.GRPCRequestsTotal) > 0 {
		pterm.DefaultSection.Println("🌐 gRPC Metrics")

		grpcData := [][]string{
			{"Method", "Requests", "Avg Latency (ms)"},
		}

		methodRequests := make(map[string]float64)
		for key, count := range data.GRPCRequestsTotal {
			parts := strings.Split(key, ":")
			if len(parts) > 0 {
				method := parts[0]
				methodRequests[method] += count
			}
		}

		for method, requests := range methodRequests {
			latency := data.GRPCDurationP50[method] * 1000
			grpcData = append(grpcData, []string{
				method,
				fmt.Sprintf("%.0f", requests),
				fmt.Sprintf("%.2f", latency),
			})
		}

		pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
			WithBoxed(true).WithData(grpcData).Render()
		pterm.Println()
	}

	// Errors Section
	if len(data.ErrorsTotal) > 0 {
		pterm.DefaultSection.Println("⚠️  Errors")

		errorData := [][]string{
			{"Error Type", "Count"},
		}

		for errorType, count := range data.ErrorsTotal {
			errorData = append(errorData, []string{
				errorType,
				pterm.Red(fmt.Sprintf("%.0f", count)),
			})
		}

		pterm.DefaultTable.WithHasHeader(true).WithHeaderStyle(pterm.NewStyle(pterm.FgLightCyan)).
			WithBoxed(true).WithData(errorData).Render()
		pterm.Println()
	}

	// Summary Box
	totalTasks := float64(0)
	for _, count := range taskSummary {
		totalTasks += count
	}

	summaryText := fmt.Sprintf(
		"Total Tasks: %s | Running: %s | Memory: %s | Goroutines: %s",
		pterm.Cyan(fmt.Sprintf("%.0f", totalTasks)),
		pterm.Yellow(fmt.Sprintf("%.0f", data.TasksRunning)),
		pterm.Green(fmt.Sprintf("%.0f MB", memoryMB)),
		pterm.Magenta(fmt.Sprintf("%.0f", data.Goroutines)),
	)

	pterm.DefaultBox.WithTitle("📈 Summary").WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).Println(summaryText)
}

// Helper functions

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

// DisplayHistoricalTrends displays metrics trends over time (for future enhancement)
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
