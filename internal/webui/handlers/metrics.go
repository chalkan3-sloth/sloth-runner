package handlers

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// MetricsHandler handles system metrics
type MetricsHandler struct {
	wsHub *WebSocketHub
}

// SystemMetrics represents system metrics
type SystemMetrics struct {
	Timestamp       int64              `json:"timestamp"`
	CPU             CPUMetrics         `json:"cpu"`
	Memory          MemoryMetrics      `json:"memory"`
	Disk            DiskMetrics        `json:"disk"`
	Network         NetworkMetrics     `json:"network"`
	HostInfo        HostInfo           `json:"host_info"`
	Goroutines      int                `json:"goroutines"`
	AgentMetrics    []AgentMetric      `json:"agent_metrics"`
	WorkflowMetrics WorkflowMetrics    `json:"workflow_metrics"`
	HookMetrics     HookMetrics        `json:"hook_metrics"`
	EventMetrics    EventMetrics       `json:"event_metrics"`
}

// CPUMetrics represents CPU metrics
type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"`
	Cores        int     `json:"cores"`
}

// MemoryMetrics represents memory metrics
type MemoryMetrics struct {
	Alloc        uint64  `json:"alloc"`
	TotalAlloc   uint64  `json:"total_alloc"`
	Sys          uint64  `json:"sys"`
	NumGC        uint32  `json:"num_gc"`
	Used         uint64  `json:"used"`
	Total        uint64  `json:"total"`
	UsedPercent  float64 `json:"used_percent"`
}

// DiskMetrics represents disk metrics
type DiskMetrics struct {
	Used        uint64  `json:"used"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// NetworkMetrics represents network metrics
type NetworkMetrics struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

// HostInfo represents host information
type HostInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	Uptime          uint64 `json:"uptime"`
}

// AgentMetric represents metrics for a single agent
type AgentMetric struct {
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	ResponseTimeMs int64   `json:"response_time_ms"`
	CPUUsage       float64 `json:"cpu_usage,omitempty"`
	MemoryUsage    float64 `json:"memory_usage,omitempty"`
}

// WorkflowMetrics represents workflow execution metrics
type WorkflowMetrics struct {
	TotalExecutions     int     `json:"total_executions"`
	SuccessfulExecutions int    `json:"successful_executions"`
	FailedExecutions    int     `json:"failed_executions"`
	AverageDurationMs   float64 `json:"average_duration_ms"`
}

// HookMetrics represents hook execution metrics
type HookMetrics struct {
	TotalExecutions   int     `json:"total_executions"`
	SuccessRate       float64 `json:"success_rate"`
	AverageDurationMs float64 `json:"average_duration_ms"`
}

// EventMetrics represents event queue metrics
type EventMetrics struct {
	PendingCount    int     `json:"pending_count"`
	ProcessedCount  int     `json:"processed_count"`
	FailedCount     int     `json:"failed_count"`
	ProcessingRate  float64 `json:"processing_rate"` // events per minute
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(wsHub *WebSocketHub) *MetricsHandler {
	handler := &MetricsHandler{
		wsHub: wsHub,
	}

	// Start metrics collection goroutine
	go handler.collectMetricsPeriodically()

	return handler
}

// GetMetrics returns current system metrics
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	metrics := h.collectMetrics()
	c.JSON(http.StatusOK, metrics)
}

// GetHistoricalMetrics returns historical metrics
func (h *MetricsHandler) GetHistoricalMetrics(c *gin.Context) {
	// TODO: Implement metrics storage and retrieval
	c.JSON(http.StatusOK, gin.H{
		"metrics": []SystemMetrics{},
		"message": "Historical metrics not yet implemented",
	})
}

// collectMetrics gathers current system metrics
func (h *MetricsHandler) collectMetrics() *SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Collect CPU metrics
	cpuPercent, _ := cpu.Percent(time.Second, false)
	var cpuUsage float64
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// Collect Memory metrics
	vmStat, _ := mem.VirtualMemory()

	// Collect Disk metrics
	diskStat, _ := disk.Usage("/")

	// Collect Network metrics
	netStats, _ := net.IOCounters(false)
	var networkMetrics NetworkMetrics
	if len(netStats) > 0 {
		networkMetrics = NetworkMetrics{
			BytesSent:   netStats[0].BytesSent,
			BytesRecv:   netStats[0].BytesRecv,
			PacketsSent: netStats[0].PacketsSent,
			PacketsRecv: netStats[0].PacketsRecv,
		}
	}

	// Collect Host info
	hostInfo, _ := host.Info()
	hostname, _ := os.Hostname()

	var hostData HostInfo
	if hostInfo != nil {
		hostData = HostInfo{
			Hostname:        hostname,
			OS:              hostInfo.OS,
			Platform:        hostInfo.Platform,
			PlatformVersion: hostInfo.PlatformVersion,
			Uptime:          hostInfo.Uptime,
		}
	} else {
		hostData = HostInfo{
			Hostname: hostname,
			OS:       runtime.GOOS,
			Platform: runtime.GOARCH,
		}
	}

	var memMetrics MemoryMetrics
	if vmStat != nil {
		memMetrics = MemoryMetrics{
			Alloc:        m.Alloc,
			TotalAlloc:   m.TotalAlloc,
			Sys:          m.Sys,
			NumGC:        m.NumGC,
			Used:         vmStat.Used,
			Total:        vmStat.Total,
			UsedPercent:  vmStat.UsedPercent,
		}
	} else {
		memMetrics = MemoryMetrics{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
			Used:       m.Alloc,
			Total:      m.Sys,
		}
	}

	var diskMetrics DiskMetrics
	if diskStat != nil {
		diskMetrics = DiskMetrics{
			Used:        diskStat.Used,
			Total:       diskStat.Total,
			Free:        diskStat.Free,
			UsedPercent: diskStat.UsedPercent,
		}
	}

	return &SystemMetrics{
		Timestamp: time.Now().Unix(),
		CPU: CPUMetrics{
			UsagePercent: cpuUsage,
			Cores:        runtime.NumCPU(),
		},
		Memory:      memMetrics,
		Disk:        diskMetrics,
		Network:     networkMetrics,
		HostInfo:    hostData,
		Goroutines:  runtime.NumGoroutine(),
		AgentMetrics: []AgentMetric{},
		WorkflowMetrics: WorkflowMetrics{
			TotalExecutions:      0,
			SuccessfulExecutions: 0,
			FailedExecutions:     0,
			AverageDurationMs:    0,
		},
		HookMetrics: HookMetrics{
			TotalExecutions:   0,
			SuccessRate:       0,
			AverageDurationMs: 0,
		},
		EventMetrics: EventMetrics{
			PendingCount:   0,
			ProcessedCount: 0,
			FailedCount:    0,
			ProcessingRate: 0,
		},
	}
}

// collectMetricsPeriodically collects and broadcasts metrics every 5 seconds
func (h *MetricsHandler) collectMetricsPeriodically() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := h.collectMetrics()
		if h.wsHub != nil {
			h.wsHub.Broadcast("system_metrics", metrics)
		}
	}
}
