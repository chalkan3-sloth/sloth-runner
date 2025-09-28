package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// MetricsCollector collects and exposes agent metrics
type MetricsCollector struct {
	mu             sync.RWMutex
	systemMetrics  SystemMetrics
	runtimeMetrics RuntimeMetrics
	taskMetrics    TaskMetrics
	customMetrics  map[string]interface{}
	
	httpServer     *http.Server
	collectTicker  *time.Ticker
	enabled        bool
}

// SystemMetrics represents system-level metrics
type SystemMetrics struct {
	CPUUsagePercent   float64   `json:"cpu_usage_percent"`
	MemoryUsageMB     float64   `json:"memory_usage_mb"`
	MemoryTotalMB     float64   `json:"memory_total_mb"`
	MemoryPercent     float64   `json:"memory_percent"`
	DiskUsageGB       float64   `json:"disk_usage_gb"`
	DiskTotalGB       float64   `json:"disk_total_gb"`
	DiskPercent       float64   `json:"disk_percent"`
	LoadAverage1m     float64   `json:"load_avg_1m"`
	LoadAverage5m     float64   `json:"load_avg_5m"`
	LoadAverage15m    float64   `json:"load_avg_15m"`
	NetworkRxBytes    uint64    `json:"network_rx_bytes"`
	NetworkTxBytes    uint64    `json:"network_tx_bytes"`
	ProcessCount      int       `json:"process_count"`
	Uptime            uint64    `json:"uptime_seconds"`
	LastUpdated       time.Time `json:"last_updated"`
}

// RuntimeMetrics represents Go runtime metrics
type RuntimeMetrics struct {
	NumGoroutines   int     `json:"num_goroutines"`
	NumCPU          int     `json:"num_cpu"`
	HeapAllocMB     float64 `json:"heap_alloc_mb"`
	HeapSysMB       float64 `json:"heap_sys_mb"`
	HeapInuseMB     float64 `json:"heap_inuse_mb"`
	StackInuseMB    float64 `json:"stack_inuse_mb"`
	NumGC           uint32  `json:"num_gc"`
	GCPauseMs       float64 `json:"gc_pause_ms"`
	NextGCMB        float64 `json:"next_gc_mb"`
	LastUpdated     time.Time `json:"last_updated"`
}

// TaskMetrics represents task execution metrics
type TaskMetrics struct {
	TotalExecuted      int64     `json:"total_executed"`
	TotalSucceeded     int64     `json:"total_succeeded"`
	TotalFailed        int64     `json:"total_failed"`
	CurrentRunning     int       `json:"current_running"`
	AverageExecTimeMs  float64   `json:"avg_exec_time_ms"`
	LastTaskStarted    time.Time `json:"last_task_started"`
	LastTaskCompleted  time.Time `json:"last_task_completed"`
	TaskQueue          int       `json:"task_queue_size"`
	LastUpdated        time.Time `json:"last_updated"`
}

// MetricsSnapshot represents a complete metrics snapshot
type MetricsSnapshot struct {
	AgentName      string                 `json:"agent_name"`
	AgentVersion   string                 `json:"agent_version"`
	Timestamp      time.Time              `json:"timestamp"`
	System         SystemMetrics          `json:"system"`
	Runtime        RuntimeMetrics         `json:"runtime"`
	Tasks          TaskMetrics            `json:"tasks"`
	Custom         map[string]interface{} `json:"custom"`
	HealthStatus   string                 `json:"health_status"`
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(agentName string, metricsPort int) *MetricsCollector {
	mc := &MetricsCollector{
		customMetrics: make(map[string]interface{}),
		enabled:      true,
	}

	// Setup HTTP server for metrics endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", mc.handleMetrics)
	mux.HandleFunc("/metrics/json", mc.handleMetricsJSON)
	mux.HandleFunc("/health", mc.handleHealth)
	
	mc.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: mux,
	}

	return mc
}

// Start begins metrics collection
func (mc *MetricsCollector) Start() error {
	mc.collectTicker = time.NewTicker(10 * time.Second)
	
	// Start collection goroutine
	go mc.collectLoop()
	
	// Start HTTP server
	go func() {
		if err := mc.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Metrics server error: %v\n", err)
		}
	}()

	return nil
}

// Stop stops metrics collection
func (mc *MetricsCollector) Stop() error {
	mc.enabled = false
	
	if mc.collectTicker != nil {
		mc.collectTicker.Stop()
	}
	
	if mc.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mc.httpServer.Shutdown(ctx)
	}
	
	return nil
}

// collectLoop runs the metrics collection loop
func (mc *MetricsCollector) collectLoop() {
	// Collect initial metrics
	mc.collectSystemMetrics()
	mc.collectRuntimeMetrics()
	
	for range mc.collectTicker.C {
		if !mc.enabled {
			break
		}
		
		mc.collectSystemMetrics()
		mc.collectRuntimeMetrics()
	}
}

// collectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) collectSystemMetrics() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// CPU usage
	if cpuPercents, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercents) > 0 {
		mc.systemMetrics.CPUUsagePercent = cpuPercents[0]
	}

	// Memory usage
	if memInfo, err := mem.VirtualMemory(); err == nil {
		mc.systemMetrics.MemoryUsageMB = float64(memInfo.Used) / 1024 / 1024
		mc.systemMetrics.MemoryTotalMB = float64(memInfo.Total) / 1024 / 1024
		mc.systemMetrics.MemoryPercent = memInfo.UsedPercent
	}

	// Disk usage
	if diskInfo, err := disk.Usage("/"); err == nil {
		mc.systemMetrics.DiskUsageGB = float64(diskInfo.Used) / 1024 / 1024 / 1024
		mc.systemMetrics.DiskTotalGB = float64(diskInfo.Total) / 1024 / 1024 / 1024
		mc.systemMetrics.DiskPercent = diskInfo.UsedPercent
	}

	// Load average
	if loadInfo, err := load.Avg(); err == nil {
		mc.systemMetrics.LoadAverage1m = loadInfo.Load1
		mc.systemMetrics.LoadAverage5m = loadInfo.Load5
		mc.systemMetrics.LoadAverage15m = loadInfo.Load15
	}

	// Network I/O
	if netStats, err := net.IOCounters(false); err == nil && len(netStats) > 0 {
		mc.systemMetrics.NetworkRxBytes = netStats[0].BytesRecv
		mc.systemMetrics.NetworkTxBytes = netStats[0].BytesSent
	}

	// Process count
	if processes, err := process.Processes(); err == nil {
		mc.systemMetrics.ProcessCount = len(processes)
	}

	// System uptime
	if hostInfo, err := host.Info(); err == nil {
		mc.systemMetrics.Uptime = hostInfo.Uptime
	}

	mc.systemMetrics.LastUpdated = time.Now()
}

// collectRuntimeMetrics collects Go runtime metrics
func (mc *MetricsCollector) collectRuntimeMetrics() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	mc.runtimeMetrics.NumGoroutines = runtime.NumGoroutine()
	mc.runtimeMetrics.NumCPU = runtime.NumCPU()
	mc.runtimeMetrics.HeapAllocMB = float64(m.HeapAlloc) / 1024 / 1024
	mc.runtimeMetrics.HeapSysMB = float64(m.HeapSys) / 1024 / 1024
	mc.runtimeMetrics.HeapInuseMB = float64(m.HeapInuse) / 1024 / 1024
	mc.runtimeMetrics.StackInuseMB = float64(m.StackInuse) / 1024 / 1024
	mc.runtimeMetrics.NumGC = m.NumGC
	mc.runtimeMetrics.NextGCMB = float64(m.NextGC) / 1024 / 1024

	// Calculate average GC pause
	if m.NumGC > 0 {
		totalPause := time.Duration(0)
		for _, pause := range m.PauseNs[:] {
			totalPause += time.Duration(pause)
		}
		mc.runtimeMetrics.GCPauseMs = float64(totalPause.Nanoseconds()) / float64(m.NumGC) / 1000000
	}

	mc.runtimeMetrics.LastUpdated = time.Now()
}

// UpdateTaskMetrics updates task-related metrics
func (mc *MetricsCollector) UpdateTaskMetrics(event string, duration time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	switch event {
	case "task_started":
		mc.taskMetrics.CurrentRunning++
		mc.taskMetrics.LastTaskStarted = time.Now()
	case "task_completed":
		mc.taskMetrics.TotalExecuted++
		mc.taskMetrics.TotalSucceeded++
		mc.taskMetrics.CurrentRunning--
		mc.taskMetrics.LastTaskCompleted = time.Now()
		
		// Update average execution time
		if mc.taskMetrics.TotalExecuted > 0 {
			currentAvg := mc.taskMetrics.AverageExecTimeMs
			newAvg := (currentAvg*float64(mc.taskMetrics.TotalExecuted-1) + float64(duration.Milliseconds())) / float64(mc.taskMetrics.TotalExecuted)
			mc.taskMetrics.AverageExecTimeMs = newAvg
		}
	case "task_failed":
		mc.taskMetrics.TotalExecuted++
		mc.taskMetrics.TotalFailed++
		mc.taskMetrics.CurrentRunning--
		mc.taskMetrics.LastTaskCompleted = time.Now()
	}

	mc.taskMetrics.LastUpdated = time.Now()
}

// SetCustomMetric sets a custom metric
func (mc *MetricsCollector) SetCustomMetric(key string, value interface{}) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.customMetrics[key] = value
}

// GetSnapshot returns a complete metrics snapshot
func (mc *MetricsCollector) GetSnapshot(agentName, version string) MetricsSnapshot {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return MetricsSnapshot{
		AgentName:      agentName,
		AgentVersion:   version,
		Timestamp:      time.Now(),
		System:         mc.systemMetrics,
		Runtime:        mc.runtimeMetrics,
		Tasks:          mc.taskMetrics,
		Custom:         mc.customMetrics,
		HealthStatus:   mc.getHealthStatus(),
	}
}

// getHealthStatus determines overall health status
func (mc *MetricsCollector) getHealthStatus() string {
	// Simple health check logic
	if mc.systemMetrics.CPUUsagePercent > 90 {
		return "critical"
	} else if mc.systemMetrics.CPUUsagePercent > 70 || mc.systemMetrics.MemoryPercent > 85 {
		return "warning"
	} else if mc.systemMetrics.DiskPercent > 90 {
		return "critical"
	} else if mc.systemMetrics.DiskPercent > 80 {
		return "warning"
	}
	return "healthy"
}

// HTTP handlers
func (mc *MetricsCollector) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snapshot := mc.GetSnapshot("", "")
	
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# HELP sloth_agent_cpu_usage_percent CPU usage percentage\n")
	fmt.Fprintf(w, "# TYPE sloth_agent_cpu_usage_percent gauge\n")
	fmt.Fprintf(w, "sloth_agent_cpu_usage_percent %.2f\n", snapshot.System.CPUUsagePercent)
	
	fmt.Fprintf(w, "# HELP sloth_agent_memory_usage_mb Memory usage in MB\n")
	fmt.Fprintf(w, "# TYPE sloth_agent_memory_usage_mb gauge\n")
	fmt.Fprintf(w, "sloth_agent_memory_usage_mb %.2f\n", snapshot.System.MemoryUsageMB)
	
	fmt.Fprintf(w, "# HELP sloth_agent_disk_usage_percent Disk usage percentage\n")
	fmt.Fprintf(w, "# TYPE sloth_agent_disk_usage_percent gauge\n")
	fmt.Fprintf(w, "sloth_agent_disk_usage_percent %.2f\n", snapshot.System.DiskPercent)
	
	fmt.Fprintf(w, "# HELP sloth_agent_tasks_total Total tasks executed\n")
	fmt.Fprintf(w, "# TYPE sloth_agent_tasks_total counter\n")
	fmt.Fprintf(w, "sloth_agent_tasks_total %d\n", snapshot.Tasks.TotalExecuted)
	
	fmt.Fprintf(w, "# HELP sloth_agent_tasks_running Currently running tasks\n")
	fmt.Fprintf(w, "# TYPE sloth_agent_tasks_running gauge\n")
	fmt.Fprintf(w, "sloth_agent_tasks_running %d\n", snapshot.Tasks.CurrentRunning)
	
	fmt.Fprintf(w, "# HELP sloth_agent_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE sloth_agent_goroutines gauge\n")
	fmt.Fprintf(w, "sloth_agent_goroutines %d\n", snapshot.Runtime.NumGoroutines)
}

func (mc *MetricsCollector) handleMetricsJSON(w http.ResponseWriter, r *http.Request) {
	snapshot := mc.GetSnapshot("", "")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

func (mc *MetricsCollector) handleHealth(w http.ResponseWriter, r *http.Request) {
	snapshot := mc.GetSnapshot("", "")
	
	status := map[string]interface{}{
		"status": snapshot.HealthStatus,
		"timestamp": time.Now(),
		"checks": map[string]interface{}{
			"cpu": map[string]interface{}{
				"usage": snapshot.System.CPUUsagePercent,
				"status": func() string {
					if snapshot.System.CPUUsagePercent > 90 {
						return "critical"
					} else if snapshot.System.CPUUsagePercent > 70 {
						return "warning"
					}
					return "healthy"
				}(),
			},
			"memory": map[string]interface{}{
				"usage": snapshot.System.MemoryPercent,
				"status": func() string {
					if snapshot.System.MemoryPercent > 90 {
						return "critical"
					} else if snapshot.System.MemoryPercent > 85 {
						return "warning"
					}
					return "healthy"
				}(),
			},
			"disk": map[string]interface{}{
				"usage": snapshot.System.DiskPercent,
				"status": func() string {
					if snapshot.System.DiskPercent > 95 {
						return "critical"
					} else if snapshot.System.DiskPercent > 80 {
						return "warning"
					}
					return "healthy"
				}(),
			},
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	if snapshot.HealthStatus == "critical" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else if snapshot.HealthStatus == "warning" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	
	json.NewEncoder(w).Encode(status)
}