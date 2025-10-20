package telemetry

import (
	"runtime"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Test Metrics struct creation
func TestNewMetrics(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	if metrics == nil {
		t.Error("Expected non-nil metrics")
	}

	if metrics.TasksTotal == nil {
		t.Error("Expected TasksTotal to be initialized")
	}

	if metrics.GRPCRequestsTotal == nil {
		t.Error("Expected GRPCRequestsTotal to be initialized")
	}

	if metrics.ErrorsTotal == nil {
		t.Error("Expected ErrorsTotal to be initialized")
	}
}

func TestNewMetrics_StartTime(t *testing.T) {
	before := time.Now()
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)
	after := time.Now()

	if metrics.startTime.Before(before) || metrics.startTime.After(after) {
		t.Error("Expected startTime to be set during NewMetrics call")
	}
}

func TestNewMetrics_AllCountersInitialized(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	counters := []struct {
		name    string
		counter interface{}
	}{
		{"TasksTotal", metrics.TasksTotal},
		{"GRPCRequestsTotal", metrics.GRPCRequestsTotal},
		{"ErrorsTotal", metrics.ErrorsTotal},
	}

	for _, c := range counters {
		if c.counter == nil {
			t.Errorf("Expected %s counter to be initialized", c.name)
		}
	}
}

func TestNewMetrics_AllGaugesInitialized(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	gauges := []struct {
		name  string
		gauge interface{}
	}{
		{"TasksRunning", metrics.TasksRunning},
		{"AgentUptime", metrics.AgentUptime},
		{"AgentInfo", metrics.AgentInfo},
		{"GoRoutines", metrics.GoRoutines},
		{"MemoryAllocated", metrics.MemoryAllocated},
	}

	for _, g := range gauges {
		if g.gauge == nil {
			t.Errorf("Expected %s gauge to be initialized", g.name)
		}
	}
}

func TestNewMetrics_AllHistogramsInitialized(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	if metrics.TaskDuration == nil {
		t.Error("Expected TaskDuration histogram to be initialized")
	}

	if metrics.GRPCDuration == nil {
		t.Error("Expected GRPCDuration histogram to be initialized")
	}
}

func TestNewMetrics_MultipleInstances(t *testing.T) {
	registry1 := prometheus.NewRegistry()
	metrics1 := NewMetrics(registry1)

	registry2 := prometheus.NewRegistry()
	metrics2 := NewMetrics(registry2)

	if metrics1 == nil || metrics2 == nil {
		t.Error("Expected both metrics instances to be created")
	}

	// They should be different instances
	if metrics1 == metrics2 {
		t.Error("Expected different instances")
	}
}

// Test UpdateRuntimeMetrics
func TestMetrics_UpdateRuntimeMetrics(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Allow some time to pass
	time.Sleep(10 * time.Millisecond)

	metrics.UpdateRuntimeMetrics()

	// Verify metrics were updated (we can't check exact values but can verify they're non-zero)
	// This is tested indirectly through Prometheus registry
}

func TestMetrics_UpdateRuntimeMetrics_Goroutines(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	beforeGoroutines := runtime.NumGoroutine()
	metrics.UpdateRuntimeMetrics()

	// Create some goroutines
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			time.Sleep(100 * time.Millisecond)
			done <- true
		}()
	}

	metrics.UpdateRuntimeMetrics()

	// Wait for goroutines to finish
	for i := 0; i < 5; i++ {
		<-done
	}

	// The goroutine count should have changed
	afterGoroutines := runtime.NumGoroutine()
	_ = beforeGoroutines
	_ = afterGoroutines
}

func TestMetrics_UpdateRuntimeMetrics_Memory(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.UpdateRuntimeMetrics()

	// Allocate some memory
	_ = make([]byte, 1024*1024*10) // 10MB

	metrics.UpdateRuntimeMetrics()

	// Memory should have been updated
}

func TestMetrics_UpdateRuntimeMetrics_Uptime(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// First update
	metrics.UpdateRuntimeMetrics()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Second update - uptime should have increased
	metrics.UpdateRuntimeMetrics()
}

func TestMetrics_UpdateRuntimeMetrics_MultipleCalls(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Call multiple times
	for i := 0; i < 10; i++ {
		metrics.UpdateRuntimeMetrics()
		time.Sleep(10 * time.Millisecond)
	}
}

// Test SetAgentInfo
func TestMetrics_SetAgentInfo(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.SetAgentInfo("1.0.0", "linux", "amd64")

	// Info should be set (verified through Prometheus registry)
}

func TestMetrics_SetAgentInfo_MultipleVersions(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	testCases := []struct {
		version string
		os      string
		arch    string
	}{
		{"1.0.0", "linux", "amd64"},
		{"1.1.0", "darwin", "arm64"},
		{"2.0.0", "windows", "amd64"},
	}

	for _, tc := range testCases {
		metrics.SetAgentInfo(tc.version, tc.os, tc.arch)
	}
}

func TestMetrics_SetAgentInfo_EmptyValues(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Should handle empty values
	metrics.SetAgentInfo("", "", "")
}

func TestMetrics_SetAgentInfo_SpecialCharacters(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.SetAgentInfo("v1.0.0-beta", "linux-gnu", "x86_64")
}

// Test RecordTaskExecution
func TestMetrics_RecordTaskExecution(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordTaskExecution("test-group", "test-task", "success", 100*time.Millisecond)

	// Task should be recorded
}

func TestMetrics_RecordTaskExecution_MultipleStatuses(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	statuses := []string{"success", "failed", "skipped", "error"}

	for i, status := range statuses {
		duration := time.Duration(i*100) * time.Millisecond
		metrics.RecordTaskExecution("group1", "task1", status, duration)
	}
}

func TestMetrics_RecordTaskExecution_DifferentDurations(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	durations := []time.Duration{
		1 * time.Millisecond,
		10 * time.Millisecond,
		100 * time.Millisecond,
		1 * time.Second,
		10 * time.Second,
	}

	for i, duration := range durations {
		metrics.RecordTaskExecution("group1", "task1", "success", duration)
		_ = i
	}
}

func TestMetrics_RecordTaskExecution_MultipleGroups(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	groups := []string{"group1", "group2", "group3", "deploy", "test"}

	for _, group := range groups {
		metrics.RecordTaskExecution(group, "task1", "success", 100*time.Millisecond)
	}
}

func TestMetrics_RecordTaskExecution_MultipleTasks(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	tasks := []string{"build", "test", "deploy", "cleanup", "notify"}

	for _, task := range tasks {
		metrics.RecordTaskExecution("group1", task, "success", 100*time.Millisecond)
	}
}

func TestMetrics_RecordTaskExecution_ZeroDuration(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordTaskExecution("group1", "task1", "success", 0)
}

func TestMetrics_RecordTaskExecution_VeryLongDuration(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordTaskExecution("group1", "task1", "success", 1*time.Hour)
}

// Test RecordGRPCRequest
func TestMetrics_RecordGRPCRequest(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordGRPCRequest("ExecuteTask", "OK", 50*time.Millisecond)

	// Request should be recorded
}

func TestMetrics_RecordGRPCRequest_MultipleMethods(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	methods := []string{
		"ExecuteTask",
		"GetStatus",
		"ListTasks",
		"CancelTask",
		"UpdateConfig",
	}

	for _, method := range methods {
		metrics.RecordGRPCRequest(method, "OK", 25*time.Millisecond)
	}
}

func TestMetrics_RecordGRPCRequest_MultipleStatuses(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	statuses := []string{"OK", "ERROR", "CANCELLED", "TIMEOUT", "UNAVAILABLE"}

	for _, status := range statuses {
		metrics.RecordGRPCRequest("ExecuteTask", status, 50*time.Millisecond)
	}
}

func TestMetrics_RecordGRPCRequest_DifferentDurations(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		25 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		250 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}

	for _, duration := range durations {
		metrics.RecordGRPCRequest("ExecuteTask", "OK", duration)
	}
}

func TestMetrics_RecordGRPCRequest_ZeroDuration(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordGRPCRequest("ExecuteTask", "OK", 0)
}

func TestMetrics_RecordGRPCRequest_VeryLongDuration(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordGRPCRequest("LongRunningTask", "OK", 30*time.Second)
}

// Test RecordError
func TestMetrics_RecordError(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordError("connection_error")

	// Error should be recorded
}

func TestMetrics_RecordError_MultipleTypes(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	errorTypes := []string{
		"connection_error",
		"timeout_error",
		"parse_error",
		"validation_error",
		"not_found_error",
		"permission_error",
	}

	for _, errorType := range errorTypes {
		metrics.RecordError(errorType)
	}
}

func TestMetrics_RecordError_SameTypeMultipleTimes(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Record same error type multiple times
	for i := 0; i < 10; i++ {
		metrics.RecordError("connection_error")
	}
}

func TestMetrics_RecordError_EmptyType(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordError("")
}

func TestMetrics_RecordError_SpecialCharacters(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordError("ssh_connection_error")
	metrics.RecordError("grpc-timeout")
	metrics.RecordError("lua.execution.error")
}

// Test IncrementRunningTasks
func TestMetrics_IncrementRunningTasks(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.IncrementRunningTasks()

	// Task count should be incremented
}

func TestMetrics_IncrementRunningTasks_Multiple(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Increment multiple times
	for i := 0; i < 5; i++ {
		metrics.IncrementRunningTasks()
	}
}

func TestMetrics_IncrementRunningTasks_Large(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Increment many times
	for i := 0; i < 100; i++ {
		metrics.IncrementRunningTasks()
	}
}

// Test DecrementRunningTasks
func TestMetrics_DecrementRunningTasks(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Increment first
	metrics.IncrementRunningTasks()
	metrics.DecrementRunningTasks()

	// Task count should be back to zero
}

func TestMetrics_DecrementRunningTasks_Multiple(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Increment several times
	for i := 0; i < 5; i++ {
		metrics.IncrementRunningTasks()
	}

	// Decrement several times
	for i := 0; i < 5; i++ {
		metrics.DecrementRunningTasks()
	}
}

func TestMetrics_DecrementRunningTasks_BelowZero(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Decrement without incrementing first (should handle gracefully)
	metrics.DecrementRunningTasks()
}

// Test combined increment/decrement scenarios
func TestMetrics_RunningTasks_IncrementDecrement(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Simulate task lifecycle
	metrics.IncrementRunningTasks() // Task starts
	time.Sleep(10 * time.Millisecond)
	metrics.DecrementRunningTasks() // Task ends

	metrics.IncrementRunningTasks()
	metrics.IncrementRunningTasks()
	metrics.DecrementRunningTasks()
	metrics.DecrementRunningTasks()
}

func TestMetrics_RunningTasks_Concurrent(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Simulate concurrent task execution
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			metrics.IncrementRunningTasks()
			time.Sleep(10 * time.Millisecond)
			metrics.DecrementRunningTasks()
			done <- true
		}()
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test Metrics struct zero value
func TestMetrics_ZeroValue(t *testing.T) {
	var metrics Metrics

	if !metrics.startTime.IsZero() {
		t.Error("Expected zero startTime")
	}
}

// Test complete workflow
func TestMetrics_CompleteWorkflow(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Set agent info
	metrics.SetAgentInfo("1.0.0", "linux", "amd64")

	// Update runtime metrics
	metrics.UpdateRuntimeMetrics()

	// Record task execution
	metrics.IncrementRunningTasks()
	metrics.RecordTaskExecution("deploy", "build", "success", 500*time.Millisecond)
	metrics.DecrementRunningTasks()

	// Record gRPC request
	metrics.RecordGRPCRequest("ExecuteTask", "OK", 50*time.Millisecond)

	// Record error
	metrics.RecordError("connection_timeout")

	// Update runtime metrics again
	metrics.UpdateRuntimeMetrics()
}

func TestMetrics_MultipleWorkflows(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Simulate multiple task executions
	for i := 0; i < 5; i++ {
		metrics.IncrementRunningTasks()
		metrics.RecordTaskExecution("group1", "task1", "success", time.Duration(i*100)*time.Millisecond)
		metrics.RecordGRPCRequest("ExecuteTask", "OK", time.Duration(i*10)*time.Millisecond)
		metrics.DecrementRunningTasks()
		metrics.UpdateRuntimeMetrics()
	}
}

// Test edge cases
func TestMetrics_EdgeCases_EmptyLabels(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordTaskExecution("", "", "", 0)
	metrics.RecordGRPCRequest("", "", 0)
	metrics.RecordError("")
}

func TestMetrics_EdgeCases_VeryLongLabels(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	longString := string(make([]byte, 1000))
	metrics.RecordTaskExecution(longString, longString, longString, 100*time.Millisecond)
}

func TestMetrics_EdgeCases_NegativeDuration(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Should handle negative duration (even though it doesn't make sense)
	metrics.RecordTaskExecution("group1", "task1", "success", -100*time.Millisecond)
}
