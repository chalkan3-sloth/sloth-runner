package telemetry

import (
	"sync"
)

var (
	// globalCollector is the singleton instance of the telemetry server
	globalCollector *Server
	once            sync.Once
)

// InitGlobal initializes the global telemetry collector
func InitGlobal(port int, enabled bool) *Server {
	once.Do(func() {
		globalCollector = NewServer(port, enabled)
	})
	return globalCollector
}

// GetGlobal returns the global telemetry collector
// Returns nil if not initialized
func GetGlobal() *Server {
	return globalCollector
}

// RecordTaskExecution is a convenience function to record task execution on the global collector
func RecordTaskExecution(group, task, status string, duration float64) {
	if globalCollector != nil && globalCollector.enabled {
		globalCollector.metrics.TasksTotal.WithLabelValues(status, group).Inc()
		globalCollector.metrics.TaskDuration.WithLabelValues(group, task).Observe(duration)
	}
}

// RecordGRPCRequest is a convenience function to record gRPC requests on the global collector
func RecordGRPCRequest(method, status string, duration float64) {
	if globalCollector != nil && globalCollector.enabled {
		globalCollector.metrics.GRPCRequestsTotal.WithLabelValues(method, status).Inc()
		globalCollector.metrics.GRPCDuration.WithLabelValues(method).Observe(duration)
	}
}

// RecordError is a convenience function to record errors on the global collector
func RecordError(errorType string) {
	if globalCollector != nil && globalCollector.enabled {
		globalCollector.metrics.ErrorsTotal.WithLabelValues(errorType).Inc()
	}
}

// IncrementRunningTasks is a convenience function to increment running tasks
func IncrementRunningTasks() {
	if globalCollector != nil && globalCollector.enabled {
		globalCollector.metrics.TasksRunning.Inc()
	}
}

// DecrementRunningTasks is a convenience function to decrement running tasks
func DecrementRunningTasks() {
	if globalCollector != nil && globalCollector.enabled {
		globalCollector.metrics.TasksRunning.Dec()
	}
}

// SetAgentInfo sets the agent version information on the global collector
func SetAgentInfo(version, os, arch string) {
	if globalCollector != nil && globalCollector.enabled {
		globalCollector.metrics.SetAgentInfo(version, os, arch)
	}
}
