package agent

import (
	"context"
	"encoding/json"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// HealthStatus represents the overall health status of an agent
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusWarning   HealthStatus = "warning"
	HealthStatusCritical  HealthStatus = "critical"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// HealthMetrics contains comprehensive health information about an agent
type HealthMetrics struct {
	AgentID     string                 `json:"agent_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Status      HealthStatus           `json:"status"`
	Uptime      time.Duration          `json:"uptime"`
	TasksTotal  int64                  `json:"tasks_total"`
	TasksActive int64                  `json:"tasks_active"`
	TasksFailed int64                  `json:"tasks_failed"`
	
	// System metrics
	System   *HealthSystemMetrics   `json:"system"`
	Memory   *MemoryMetrics   `json:"memory"`
	CPU      *CPUMetrics      `json:"cpu"`
	Disk     *DiskMetrics     `json:"disk"`
	Network  *NetworkMetrics  `json:"network"`
	
	// Custom metrics
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// HealthSystemMetrics contains system-level information
type HealthSystemMetrics struct {
	Hostname        string  `json:"hostname"`
	Platform        string  `json:"platform"`
	PlatformFamily  string  `json:"platform_family"`
	PlatformVersion string  `json:"platform_version"`
	KernelVersion   string  `json:"kernel_version"`
	LoadAverage     []float64 `json:"load_average"`
	Processes       uint64  `json:"processes"`
}

// MemoryMetrics contains memory usage information
type MemoryMetrics struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Buffers     uint64  `json:"buffers,omitempty"`
	Cached      uint64  `json:"cached,omitempty"`
}

// CPUMetrics contains CPU usage information
type CPUMetrics struct {
	LogicalCores  int     `json:"logical_cores"`
	PhysicalCores int     `json:"physical_cores"`
	UsagePercent  float64 `json:"usage_percent"`
	ModelName     string  `json:"model_name"`
	Frequency     float64 `json:"frequency"`
}

// DiskMetrics contains disk usage information
type DiskMetrics struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Path        string  `json:"path"`
}

// NetworkMetrics contains network statistics
type NetworkMetrics struct {
	BytesReceived uint64 `json:"bytes_received"`
	BytesSent     uint64 `json:"bytes_sent"`
	PacketsReceived uint64 `json:"packets_received"`
	PacketsSent   uint64 `json:"packets_sent"`
	Errors        uint64 `json:"errors"`
	Drops         uint64 `json:"drops"`
}

// HealthMonitor monitors agent health and system metrics
type HealthMonitor struct {
	agentID     string
	startTime   time.Time
	mu          sync.RWMutex
	metrics     *HealthMetrics
	thresholds  *HealthThresholds
	collectors  []MetricCollector
}

// HealthThresholds defines warning and critical thresholds
type HealthThresholds struct {
	CPUWarning     float64 `json:"cpu_warning"`
	CPUCritical    float64 `json:"cpu_critical"`
	MemoryWarning  float64 `json:"memory_warning"`
	MemoryCritical float64 `json:"memory_critical"`
	DiskWarning    float64 `json:"disk_warning"`
	DiskCritical   float64 `json:"disk_critical"`
	LoadWarning    float64 `json:"load_warning"`
	LoadCritical   float64 `json:"load_critical"`
}

// MetricCollector interface for custom metric collection
type MetricCollector interface {
	Name() string
	Collect() (map[string]interface{}, error)
}

// NewHealthMonitor creates a new health monitor instance
func NewHealthMonitor(agentID string) *HealthMonitor {
	return &HealthMonitor{
		agentID:   agentID,
		startTime: time.Now(),
		thresholds: &HealthThresholds{
			CPUWarning:     70.0,
			CPUCritical:    90.0,
			MemoryWarning:  80.0,
			MemoryCritical: 95.0,
			DiskWarning:    85.0,
			DiskCritical:   95.0,
			LoadWarning:    2.0,
			LoadCritical:   5.0,
		},
		collectors: make([]MetricCollector, 0),
	}
}

// SetThresholds updates the health thresholds
func (hm *HealthMonitor) SetThresholds(thresholds *HealthThresholds) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.thresholds = thresholds
}

// AddCollector adds a custom metric collector
func (hm *HealthMonitor) AddCollector(collector MetricCollector) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.collectors = append(hm.collectors, collector)
}

// CollectMetrics gathers all system and custom metrics
func (hm *HealthMonitor) CollectMetrics(ctx context.Context) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	metrics := &HealthMetrics{
		AgentID:   hm.agentID,
		Timestamp: time.Now(),
		Uptime:    time.Since(hm.startTime),
		Custom:    make(map[string]interface{}),
	}

	// Collect system metrics
	if systemMetrics, err := hm.collectSystemMetrics(ctx); err == nil {
		metrics.System = systemMetrics
	}

	// Collect memory metrics
	if memMetrics, err := hm.collectMemoryMetrics(ctx); err == nil {
		metrics.Memory = memMetrics
	}

	// Collect CPU metrics
	if cpuMetrics, err := hm.collectCPUMetrics(ctx); err == nil {
		metrics.CPU = cpuMetrics
	}

	// Collect disk metrics
	if diskMetrics, err := hm.collectDiskMetrics(ctx); err == nil {
		metrics.Disk = diskMetrics
	}

	// Collect network metrics
	if netMetrics, err := hm.collectNetworkMetrics(ctx); err == nil {
		metrics.Network = netMetrics
	}

	// Collect custom metrics
	for _, collector := range hm.collectors {
		if customMetrics, err := collector.Collect(); err == nil {
			metrics.Custom[collector.Name()] = customMetrics
		}
	}

	// Determine overall health status
	metrics.Status = hm.determineHealthStatus(metrics)

	hm.metrics = metrics
	return nil
}

// GetMetrics returns the current health metrics
func (hm *HealthMonitor) GetMetrics() *HealthMetrics {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	
	if hm.metrics == nil {
		return nil
	}

	// Return a copy to prevent concurrent access issues
	data, _ := json.Marshal(hm.metrics)
	var copy HealthMetrics
	json.Unmarshal(data, &copy)
	return &copy
}

// collectSystemMetrics gathers system-level metrics
func (hm *HealthMonitor) collectSystemMetrics(ctx context.Context) (*HealthSystemMetrics, error) {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	loadInfo, err := load.AvgWithContext(ctx)
	if err != nil {
		loadInfo = &load.AvgStat{}
	}

	return &HealthSystemMetrics{
		Hostname:        hostInfo.Hostname,
		Platform:        hostInfo.Platform,
		PlatformFamily:  hostInfo.PlatformFamily,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		LoadAverage:     []float64{loadInfo.Load1, loadInfo.Load5, loadInfo.Load15},
		Processes:       hostInfo.Procs,
	}, nil
}

// collectMemoryMetrics gathers memory usage metrics
func (hm *HealthMonitor) collectMemoryMetrics(ctx context.Context) (*MemoryMetrics, error) {
	memInfo, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return &MemoryMetrics{
		Total:       memInfo.Total,
		Available:   memInfo.Available,
		Used:        memInfo.Used,
		UsedPercent: memInfo.UsedPercent,
		Free:        memInfo.Free,
		Buffers:     memInfo.Buffers,
		Cached:      memInfo.Cached,
	}, nil
}

// collectCPUMetrics gathers CPU usage metrics
func (hm *HealthMonitor) collectCPUMetrics(ctx context.Context) (*CPUMetrics, error) {
	// Get CPU usage percentage
	cpuPercent, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return nil, err
	}

	// Get CPU info
	cpuInfo, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	var usage float64
	if len(cpuPercent) > 0 {
		usage = cpuPercent[0]
	}

	metrics := &CPUMetrics{
		LogicalCores:  runtime.NumCPU(),
		PhysicalCores: runtime.NumCPU(),
		UsagePercent:  usage,
	}

	if len(cpuInfo) > 0 {
		metrics.ModelName = cpuInfo[0].ModelName
		metrics.Frequency = cpuInfo[0].Mhz
	}

	return metrics, nil
}

// collectDiskMetrics gathers disk usage metrics
func (hm *HealthMonitor) collectDiskMetrics(ctx context.Context) (*DiskMetrics, error) {
	diskUsage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return nil, err
	}

	return &DiskMetrics{
		Total:       diskUsage.Total,
		Free:        diskUsage.Free,
		Used:        diskUsage.Used,
		UsedPercent: diskUsage.UsedPercent,
		Path:        diskUsage.Path,
	}, nil
}

// collectNetworkMetrics gathers network statistics
func (hm *HealthMonitor) collectNetworkMetrics(ctx context.Context) (*NetworkMetrics, error) {
	// This is a simplified implementation
	// In a real implementation, you'd want to track interface-specific metrics
	return &NetworkMetrics{
		BytesReceived:   0,
		BytesSent:       0,
		PacketsReceived: 0,
		PacketsSent:     0,
		Errors:         0,
		Drops:          0,
	}, nil
}

// determineHealthStatus analyzes metrics and determines overall health
func (hm *HealthMonitor) determineHealthStatus(metrics *HealthMetrics) HealthStatus {
	if metrics == nil {
		return HealthStatusUnknown
	}

	// Check CPU
	if metrics.CPU != nil && metrics.CPU.UsagePercent > hm.thresholds.CPUCritical {
		return HealthStatusCritical
	}

	// Check Memory
	if metrics.Memory != nil && metrics.Memory.UsedPercent > hm.thresholds.MemoryCritical {
		return HealthStatusCritical
	}

	// Check Disk
	if metrics.Disk != nil && metrics.Disk.UsedPercent > hm.thresholds.DiskCritical {
		return HealthStatusCritical
	}

	// Check Load Average
	if metrics.System != nil && len(metrics.System.LoadAverage) > 0 {
		if metrics.System.LoadAverage[0] > hm.thresholds.LoadCritical {
			return HealthStatusCritical
		}
	}

	// Check for warnings
	if (metrics.CPU != nil && metrics.CPU.UsagePercent > hm.thresholds.CPUWarning) ||
	   (metrics.Memory != nil && metrics.Memory.UsedPercent > hm.thresholds.MemoryWarning) ||
	   (metrics.Disk != nil && metrics.Disk.UsedPercent > hm.thresholds.DiskWarning) ||
	   (metrics.System != nil && len(metrics.System.LoadAverage) > 0 && metrics.System.LoadAverage[0] > hm.thresholds.LoadWarning) {
		return HealthStatusWarning
	}

	return HealthStatusHealthy
}

// StartPeriodicCollection starts periodic metric collection
func (hm *HealthMonitor) StartPeriodicCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := hm.CollectMetrics(ctx); err != nil {
				// Log error but continue monitoring
				continue
			}
		}
	}
}

// TaskMetricsCollector implements MetricCollector for task-related metrics
type TaskMetricsCollector struct {
	taskStats *TaskStatistics
}

// TaskStatistics holds task execution statistics
type TaskStatistics struct {
	TotalTasks     int64         `json:"total_tasks"`
	RunningTasks   int64         `json:"running_tasks"`
	CompletedTasks int64         `json:"completed_tasks"`
	FailedTasks    int64         `json:"failed_tasks"`
	AverageRuntime time.Duration `json:"average_runtime"`
	LastTaskTime   time.Time     `json:"last_task_time"`
}

// NewTaskMetricsCollector creates a task metrics collector
func NewTaskMetricsCollector() *TaskMetricsCollector {
	return &TaskMetricsCollector{
		taskStats: &TaskStatistics{},
	}
}

// Name returns the collector name
func (tmc *TaskMetricsCollector) Name() string {
	return "tasks"
}

// Collect returns task metrics
func (tmc *TaskMetricsCollector) Collect() (map[string]interface{}, error) {
	data, err := json.Marshal(tmc.taskStats)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateTaskStats updates task statistics
func (tmc *TaskMetricsCollector) UpdateTaskStats(stats *TaskStatistics) {
	tmc.taskStats = stats
}