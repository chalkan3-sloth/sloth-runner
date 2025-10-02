package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMetricsCollector(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	
	if mc == nil {
		t.Fatal("Expected non-nil MetricsCollector")
	}
	if !mc.enabled {
		t.Error("Expected metrics collector to be enabled by default")
	}
	if mc.customMetrics == nil {
		t.Error("Expected custom metrics map to be initialized")
	}
}

func TestMetricsCollector_CollectSystemMetrics(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	
	mc.collectSystemMetrics()
	
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	if mc.systemMetrics.MemoryTotalMB == 0 {
		t.Error("Expected non-zero total memory")
	}
	if mc.systemMetrics.LastUpdated.IsZero() {
		t.Error("Expected last updated time to be set")
	}
}

func TestMetricsCollector_CollectRuntimeMetrics(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	
	mc.collectRuntimeMetrics()
	
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	if mc.runtimeMetrics.NumCPU == 0 {
		t.Error("Expected non-zero CPU count")
	}
	if mc.runtimeMetrics.NumGoroutines == 0 {
		t.Error("Expected non-zero goroutine count")
	}
	if mc.runtimeMetrics.LastUpdated.IsZero() {
		t.Error("Expected last updated time to be set")
	}
}

func TestMetricsCollector_GetSnapshot(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	mc.collectSystemMetrics()
	mc.collectRuntimeMetrics()
	
	snapshot := mc.GetSnapshot("test-agent", "1.0.0")
	
	if snapshot.AgentName != "test-agent" {
		t.Errorf("Expected agent name 'test-agent', got '%s'", snapshot.AgentName)
	}
	if snapshot.AgentVersion != "1.0.0" {
		t.Errorf("Expected agent version '1.0.0', got '%s'", snapshot.AgentVersion)
	}
	if snapshot.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}

func TestMetricsCollector_UpdateTaskMetrics(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	
	mc.UpdateTaskMetrics("task_started", 0)
	mc.UpdateTaskMetrics("task_completed", 100 * time.Millisecond)
	mc.UpdateTaskMetrics("task_started", 0)
	mc.UpdateTaskMetrics("task_failed", 50 * time.Millisecond)
	
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	if mc.taskMetrics.TotalExecuted != 2 {
		t.Errorf("Expected 2 total executions, got %d", mc.taskMetrics.TotalExecuted)
	}
	if mc.taskMetrics.TotalSucceeded != 1 {
		t.Errorf("Expected 1 success, got %d", mc.taskMetrics.TotalSucceeded)
	}
	if mc.taskMetrics.TotalFailed != 1 {
		t.Errorf("Expected 1 failure, got %d", mc.taskMetrics.TotalFailed)
	}
}

func TestMetricsCollector_CustomMetrics(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	
	mc.SetCustomMetric("test_metric", 42)
	mc.SetCustomMetric("test_string", "value")
	
	mc.mu.RLock()
	value, exists := mc.customMetrics["test_metric"]
	mc.mu.RUnlock()
	
	if !exists {
		t.Error("Expected custom metric to exist")
	}
	if value != 42 {
		t.Errorf("Expected metric value 42, got %v", value)
	}
	
	mc.mu.RLock()
	strValue, exists := mc.customMetrics["test_string"]
	mc.mu.RUnlock()
	
	if !exists {
		t.Error("Expected custom string metric to exist")
	}
	if strValue != "value" {
		t.Errorf("Expected metric value 'value', got %v", strValue)
	}
}

func TestMetricsCollector_HTTPHandler(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	mc.collectSystemMetrics()
	mc.collectRuntimeMetrics()
	
	req := httptest.NewRequest(http.MethodGet, "/metrics/json", nil)
	w := httptest.NewRecorder()
	
	mc.handleMetricsJSON(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	var snapshot MetricsSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&snapshot); err != nil {
		t.Errorf("Failed to decode metrics response: %v", err)
	}
	
	if snapshot.System.MemoryTotalMB == 0 {
		t.Error("Expected system metrics to be populated")
	}
}

func TestSystemMetrics_JSONMarshaling(t *testing.T) {
	metrics := SystemMetrics{
		CPUUsagePercent: 45.5,
		MemoryUsageMB:   1024.0,
		MemoryTotalMB:   4096.0,
		MemoryPercent:   25.0,
		ProcessCount:    150,
		LastUpdated:     time.Now(),
	}
	
	data, err := json.Marshal(metrics)
	if err != nil {
		t.Errorf("Failed to marshal metrics: %v", err)
	}
	
	var unmarshaled SystemMetrics
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal metrics: %v", err)
	}
	
	if unmarshaled.CPUUsagePercent != metrics.CPUUsagePercent {
		t.Error("CPU usage not preserved after marshaling")
	}
	if unmarshaled.MemoryUsageMB != metrics.MemoryUsageMB {
		t.Error("Memory usage not preserved after marshaling")
	}
}

func TestRuntimeMetrics_JSONMarshaling(t *testing.T) {
	metrics := RuntimeMetrics{
		NumGoroutines: 25,
		NumCPU:        8,
		HeapAllocMB:   128.5,
		NumGC:         42,
		LastUpdated:   time.Now(),
	}
	
	data, err := json.Marshal(metrics)
	if err != nil {
		t.Errorf("Failed to marshal metrics: %v", err)
	}
	
	var unmarshaled RuntimeMetrics
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal metrics: %v", err)
	}
	
	if unmarshaled.NumGoroutines != metrics.NumGoroutines {
		t.Error("Goroutine count not preserved after marshaling")
	}
	if unmarshaled.NumCPU != metrics.NumCPU {
		t.Error("CPU count not preserved after marshaling")
	}
}

func TestTaskMetrics_JSONMarshaling(t *testing.T) {
	metrics := TaskMetrics{
		TotalExecuted:     100,
		TotalSucceeded:    95,
		TotalFailed:       5,
		CurrentRunning:    3,
		AverageExecTimeMs: 123.45,
		LastUpdated:       time.Now(),
	}
	
	data, err := json.Marshal(metrics)
	if err != nil {
		t.Errorf("Failed to marshal metrics: %v", err)
	}
	
	var unmarshaled TaskMetrics
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal metrics: %v", err)
	}
	
	if unmarshaled.TotalExecuted != metrics.TotalExecuted {
		t.Error("Total executed not preserved after marshaling")
	}
	if unmarshaled.TotalSucceeded != metrics.TotalSucceeded {
		t.Error("Total succeeded not preserved after marshaling")
	}
}

func TestMetricsCollector_Enable_Disable(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	
	if !mc.enabled {
		t.Error("Expected collector to be enabled by default")
	}
	
	mc.Stop()
	if mc.enabled {
		t.Error("Expected collector to be disabled after Stop")
	}
	
	mc.enabled = true
	if !mc.enabled {
		t.Error("Expected collector to be enabled again")
	}
}

func TestMetricsCollector_ConcurrentAccess(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	
	done := make(chan bool, 3)
	
	// Concurrent reads
	go func() {
		for i := 0; i < 10; i++ {
			mc.GetSnapshot("test", "1.0")
			time.Sleep(5 * time.Millisecond)
		}
		done <- true
	}()
	
	// Concurrent writes (task metrics)
	go func() {
		for i := 0; i < 10; i++ {
			mc.UpdateTaskMetrics("task_started", 0)
			mc.UpdateTaskMetrics("task_completed", 10 * time.Millisecond)
			time.Sleep(5 * time.Millisecond)
		}
		done <- true
	}()
	
	// Concurrent custom metrics
	go func() {
		for i := 0; i < 10; i++ {
			mc.SetCustomMetric("test", i)
			mc.mu.RLock()
			_ = mc.customMetrics["test"]
			mc.mu.RUnlock()
			time.Sleep(5 * time.Millisecond)
		}
		done <- true
	}()
	
	// Wait for all goroutines
	<-done
	<-done
	<-done
	
	// No panic means thread-safety is working
}

func TestMetricsSnapshot_CompleteData(t *testing.T) {
	mc := NewMetricsCollector("test-agent", 9090)
	mc.collectSystemMetrics()
	mc.collectRuntimeMetrics()
	mc.UpdateTaskMetrics("task_started", 0)
	mc.UpdateTaskMetrics("task_completed", 50 * time.Millisecond)
	mc.SetCustomMetric("custom_field", "test_value")
	
	snapshot := mc.GetSnapshot("test-agent", "v1.2.3")
	
	if snapshot.AgentName == "" {
		t.Error("Expected agent name to be set")
	}
	if snapshot.AgentVersion == "" {
		t.Error("Expected agent version to be set")
	}
	if snapshot.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
	if snapshot.System.MemoryTotalMB == 0 {
		t.Error("Expected system metrics to be populated")
	}
	if snapshot.Runtime.NumGoroutines == 0 {
		t.Error("Expected runtime metrics to be populated")
	}
	if snapshot.Tasks.TotalExecuted == 0 {
		t.Error("Expected task metrics to be populated")
	}
	if len(snapshot.Custom) == 0 {
		t.Error("Expected custom metrics to be populated")
	}
}

func BenchmarkMetricsCollector_CollectSystemMetrics(b *testing.B) {
	mc := NewMetricsCollector("bench-agent", 9090)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.collectSystemMetrics()
	}
}

func BenchmarkMetricsCollector_CollectRuntimeMetrics(b *testing.B) {
	mc := NewMetricsCollector("bench-agent", 9090)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.collectRuntimeMetrics()
	}
}

func BenchmarkMetricsCollector_GetSnapshot(b *testing.B) {
	mc := NewMetricsCollector("bench-agent", 9090)
	mc.collectSystemMetrics()
	mc.collectRuntimeMetrics()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.GetSnapshot("bench-agent", "1.0.0")
	}
}
