package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
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
	Alloc      uint64  `json:"alloc"`
	TotalAlloc uint64  `json:"total_alloc"`
	Sys        uint64  `json:"sys"`
	NumGC      uint32  `json:"num_gc"`
	UsedMB     float64 `json:"used_mb"`
	TotalMB    float64 `json:"total_mb"`
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

	return &SystemMetrics{
		Timestamp: time.Now().Unix(),
		CPU: CPUMetrics{
			UsagePercent: 0, // TODO: Implement CPU usage calculation
			Cores:        runtime.NumCPU(),
		},
		Memory: MemoryMetrics{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
			UsedMB:     float64(m.Alloc) / 1024 / 1024,
			TotalMB:    float64(m.Sys) / 1024 / 1024,
		},
		Goroutines: runtime.NumGoroutine(),
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
