package telemetry

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds all Prometheus metrics for the sloth-runner agent
type Metrics struct {
	// Counters
	TasksTotal        *prometheus.CounterVec
	GRPCRequestsTotal *prometheus.CounterVec
	ErrorsTotal       *prometheus.CounterVec

	// Gauges
	TasksRunning    prometheus.Gauge
	AgentUptime     prometheus.Gauge
	AgentInfo       *prometheus.GaugeVec
	GoRoutines      prometheus.Gauge
	MemoryAllocated prometheus.Gauge

	// Histograms
	TaskDuration *prometheus.HistogramVec
	GRPCDuration *prometheus.HistogramVec

	// Internal
	startTime time.Time
}

// NewMetrics creates and registers all metrics with Prometheus
func NewMetrics(registry *prometheus.Registry) *Metrics {
	m := &Metrics{
		startTime: time.Now(),

		// Counters
		TasksTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "sloth_tasks_total",
				Help: "Total number of tasks executed by status",
			},
			[]string{"status", "group"},
		),

		GRPCRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "sloth_grpc_requests_total",
				Help: "Total number of gRPC requests received",
			},
			[]string{"method", "status"},
		),

		ErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "sloth_errors_total",
				Help: "Total number of errors by type",
			},
			[]string{"type"},
		),

		// Gauges
		TasksRunning: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "sloth_tasks_running",
				Help: "Number of tasks currently running",
			},
		),

		AgentUptime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "sloth_agent_uptime_seconds",
				Help: "Agent uptime in seconds",
			},
		),

		AgentInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "sloth_agent_info",
				Help: "Agent version and build information",
			},
			[]string{"version", "os", "arch"},
		),

		GoRoutines: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "sloth_goroutines",
				Help: "Number of goroutines currently running",
			},
		),

		MemoryAllocated: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "sloth_memory_allocated_bytes",
				Help: "Memory allocated in bytes",
			},
		),

		// Histograms
		TaskDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "sloth_task_duration_seconds",
				Help:    "Task execution duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"group", "task"},
		),

		GRPCDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "sloth_grpc_request_duration_seconds",
				Help: "gRPC request duration in seconds",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method"},
		),
	}

	// Register all metrics
	registry.MustRegister(
		m.TasksTotal,
		m.GRPCRequestsTotal,
		m.ErrorsTotal,
		m.TasksRunning,
		m.AgentUptime,
		m.AgentInfo,
		m.GoRoutines,
		m.MemoryAllocated,
		m.TaskDuration,
		m.GRPCDuration,
	)

	return m
}

// UpdateRuntimeMetrics updates runtime metrics (goroutines, memory, uptime)
func (m *Metrics) UpdateRuntimeMetrics() {
	// Update goroutines
	m.GoRoutines.Set(float64(runtime.NumGoroutine()))

	// Update memory stats
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m.MemoryAllocated.Set(float64(mem.Alloc))

	// Update uptime
	m.AgentUptime.Set(time.Since(m.startTime).Seconds())
}

// SetAgentInfo sets the agent version information
func (m *Metrics) SetAgentInfo(version, os, arch string) {
	m.AgentInfo.WithLabelValues(version, os, arch).Set(1)
}

// RecordTaskExecution records metrics for a task execution
func (m *Metrics) RecordTaskExecution(group, task, status string, duration time.Duration) {
	m.TasksTotal.WithLabelValues(status, group).Inc()
	m.TaskDuration.WithLabelValues(group, task).Observe(duration.Seconds())
}

// RecordGRPCRequest records metrics for a gRPC request
func (m *Metrics) RecordGRPCRequest(method, status string, duration time.Duration) {
	m.GRPCRequestsTotal.WithLabelValues(method, status).Inc()
	m.GRPCDuration.WithLabelValues(method).Observe(duration.Seconds())
}

// RecordError records an error occurrence
func (m *Metrics) RecordError(errorType string) {
	m.ErrorsTotal.WithLabelValues(errorType).Inc()
}

// IncrementRunningTasks increments the running tasks gauge
func (m *Metrics) IncrementRunningTasks() {
	m.TasksRunning.Inc()
}

// DecrementRunningTasks decrements the running tasks gauge
func (m *Metrics) DecrementRunningTasks() {
	m.TasksRunning.Dec()
}
