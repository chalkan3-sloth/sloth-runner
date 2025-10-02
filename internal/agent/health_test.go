package agent

import (
	"context"
	"testing"
	"time"
)

func TestNewHealthMonitor(t *testing.T) {
	agentID := "test-agent-123"
	hm := NewHealthMonitor(agentID)

	if hm == nil {
		t.Fatal("Expected non-nil HealthMonitor")
	}
	if hm.agentID != agentID {
		t.Errorf("Expected agentID '%s', got '%s'", agentID, hm.agentID)
	}
	if hm.thresholds == nil {
		t.Error("Expected default thresholds to be set")
	}
	if len(hm.collectors) != 0 {
		t.Error("Expected empty collectors list")
	}
}

func TestHealthMonitor_SetThresholds(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	
	customThresholds := &HealthThresholds{
		CPUWarning:     60.0,
		CPUCritical:    85.0,
		MemoryWarning:  75.0,
		MemoryCritical: 90.0,
		DiskWarning:    80.0,
		DiskCritical:   95.0,
		LoadWarning:    3.0,
		LoadCritical:   6.0,
	}

	hm.SetThresholds(customThresholds)

	if hm.thresholds.CPUWarning != 60.0 {
		t.Errorf("Expected CPUWarning 60.0, got %f", hm.thresholds.CPUWarning)
	}
	if hm.thresholds.MemoryCritical != 90.0 {
		t.Errorf("Expected MemoryCritical 90.0, got %f", hm.thresholds.MemoryCritical)
	}
}

func TestHealthMonitor_AddCollector(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	collector := NewTaskMetricsCollector()

	hm.AddCollector(collector)

	if len(hm.collectors) != 1 {
		t.Errorf("Expected 1 collector, got %d", len(hm.collectors))
	}
}

func TestHealthMonitor_CollectMetrics(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	ctx := context.Background()

	err := hm.CollectMetrics(ctx)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	metrics := hm.GetMetrics()
	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}
	if metrics.AgentID != "test-agent" {
		t.Errorf("Expected AgentID 'test-agent', got '%s'", metrics.AgentID)
	}
	if metrics.Status == HealthStatusUnknown {
		t.Error("Expected status to be determined")
	}
}

func TestHealthMonitor_GetMetrics_NoCollection(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	
	metrics := hm.GetMetrics()
	if metrics != nil {
		t.Error("Expected nil metrics before collection")
	}
}

func TestHealthMonitor_CollectSystemMetrics(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	ctx := context.Background()

	systemMetrics, err := hm.collectSystemMetrics(ctx)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if systemMetrics == nil {
		t.Fatal("Expected non-nil system metrics")
	}
	if systemMetrics.Hostname == "" {
		t.Error("Expected non-empty hostname")
	}
	if len(systemMetrics.LoadAverage) != 3 {
		t.Errorf("Expected 3 load average values, got %d", len(systemMetrics.LoadAverage))
	}
}

func TestHealthMonitor_CollectMemoryMetrics(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	ctx := context.Background()

	memMetrics, err := hm.collectMemoryMetrics(ctx)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if memMetrics == nil {
		t.Fatal("Expected non-nil memory metrics")
	}
	if memMetrics.Total == 0 {
		t.Error("Expected non-zero total memory")
	}
	if memMetrics.UsedPercent < 0 || memMetrics.UsedPercent > 100 {
		t.Errorf("Expected memory usage between 0-100%%, got %f", memMetrics.UsedPercent)
	}
}

func TestHealthMonitor_CollectCPUMetrics(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	ctx := context.Background()

	cpuMetrics, err := hm.collectCPUMetrics(ctx)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if cpuMetrics == nil {
		t.Fatal("Expected non-nil CPU metrics")
	}
	if cpuMetrics.LogicalCores <= 0 {
		t.Error("Expected positive number of logical cores")
	}
	if cpuMetrics.UsagePercent < 0 || cpuMetrics.UsagePercent > 100 {
		t.Errorf("Expected CPU usage between 0-100%%, got %f", cpuMetrics.UsagePercent)
	}
}

func TestHealthMonitor_CollectDiskMetrics(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	ctx := context.Background()

	diskMetrics, err := hm.collectDiskMetrics(ctx)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if diskMetrics == nil {
		t.Fatal("Expected non-nil disk metrics")
	}
	if diskMetrics.Total == 0 {
		t.Error("Expected non-zero total disk space")
	}
	if diskMetrics.UsedPercent < 0 || diskMetrics.UsedPercent > 100 {
		t.Errorf("Expected disk usage between 0-100%%, got %f", diskMetrics.UsedPercent)
	}
}

func TestHealthMonitor_DetermineHealthStatus_Healthy(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	
	metrics := &HealthMetrics{
		CPU:    &CPUMetrics{UsagePercent: 30.0},
		Memory: &MemoryMetrics{UsedPercent: 40.0},
		Disk:   &DiskMetrics{UsedPercent: 50.0},
		System: &HealthSystemMetrics{LoadAverage: []float64{1.0, 1.0, 1.0}},
	}

	status := hm.determineHealthStatus(metrics)
	if status != HealthStatusHealthy {
		t.Errorf("Expected status Healthy, got %s", status)
	}
}

func TestHealthMonitor_DetermineHealthStatus_Warning(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	
	metrics := &HealthMetrics{
		CPU:    &CPUMetrics{UsagePercent: 75.0}, // Above warning threshold (70%)
		Memory: &MemoryMetrics{UsedPercent: 50.0},
		Disk:   &DiskMetrics{UsedPercent: 50.0},
		System: &HealthSystemMetrics{LoadAverage: []float64{1.0, 1.0, 1.0}},
	}

	status := hm.determineHealthStatus(metrics)
	if status != HealthStatusWarning {
		t.Errorf("Expected status Warning, got %s", status)
	}
}

func TestHealthMonitor_DetermineHealthStatus_Critical(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	
	metrics := &HealthMetrics{
		CPU:    &CPUMetrics{UsagePercent: 95.0}, // Above critical threshold (90%)
		Memory: &MemoryMetrics{UsedPercent: 50.0},
		Disk:   &DiskMetrics{UsedPercent: 50.0},
		System: &HealthSystemMetrics{LoadAverage: []float64{1.0, 1.0, 1.0}},
	}

	status := hm.determineHealthStatus(metrics)
	if status != HealthStatusCritical {
		t.Errorf("Expected status Critical, got %s", status)
	}
}

func TestHealthMonitor_DetermineHealthStatus_Unknown(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	
	status := hm.determineHealthStatus(nil)
	if status != HealthStatusUnknown {
		t.Errorf("Expected status Unknown, got %s", status)
	}
}

func TestHealthMonitor_StartPeriodicCollection(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Start periodic collection in a goroutine
	done := make(chan bool)
	go func() {
		hm.StartPeriodicCollection(ctx, 50*time.Millisecond)
		done <- true
	}()

	// Wait for context to complete
	<-ctx.Done()
	<-done

	// Verify that metrics were collected
	metrics := hm.GetMetrics()
	if metrics == nil {
		t.Error("Expected metrics to be collected")
	}
}

func TestTaskMetricsCollector_New(t *testing.T) {
	collector := NewTaskMetricsCollector()
	
	if collector == nil {
		t.Fatal("Expected non-nil collector")
	}
	if collector.Name() != "tasks" {
		t.Errorf("Expected name 'tasks', got '%s'", collector.Name())
	}
}

func TestTaskMetricsCollector_Collect(t *testing.T) {
	collector := NewTaskMetricsCollector()
	
	metrics, err := collector.Collect()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if metrics == nil {
		t.Error("Expected non-nil metrics")
	}
}

func TestTaskMetricsCollector_UpdateTaskStats(t *testing.T) {
	collector := NewTaskMetricsCollector()
	
	stats := &TaskStatistics{
		TotalTasks:     100,
		RunningTasks:   5,
		CompletedTasks: 90,
		FailedTasks:    5,
		LastTaskTime:   time.Now(),
	}

	collector.UpdateTaskStats(stats)
	
	metrics, err := collector.Collect()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if totalTasks, ok := metrics["total_tasks"].(float64); !ok || int64(totalTasks) != 100 {
		t.Errorf("Expected total_tasks 100, got %v", metrics["total_tasks"])
	}
}

func TestHealthStatus_Constants(t *testing.T) {
	tests := []struct {
		status   HealthStatus
		expected string
	}{
		{HealthStatusHealthy, "healthy"},
		{HealthStatusWarning, "warning"},
		{HealthStatusCritical, "critical"},
		{HealthStatusUnknown, "unknown"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("Expected status '%s', got '%s'", tt.expected, string(tt.status))
		}
	}
}

func TestHealthThresholds_DefaultValues(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	
	if hm.thresholds.CPUWarning != 70.0 {
		t.Errorf("Expected default CPUWarning 70.0, got %f", hm.thresholds.CPUWarning)
	}
	if hm.thresholds.CPUCritical != 90.0 {
		t.Errorf("Expected default CPUCritical 90.0, got %f", hm.thresholds.CPUCritical)
	}
	if hm.thresholds.MemoryWarning != 80.0 {
		t.Errorf("Expected default MemoryWarning 80.0, got %f", hm.thresholds.MemoryWarning)
	}
	if hm.thresholds.MemoryCritical != 95.0 {
		t.Errorf("Expected default MemoryCritical 95.0, got %f", hm.thresholds.MemoryCritical)
	}
}

func TestHealthMonitor_ConcurrentAccess(t *testing.T) {
	hm := NewHealthMonitor("test-agent")
	ctx := context.Background()

	// Test concurrent metric collection and retrieval
	done := make(chan bool, 2)

	// Goroutine 1: Collect metrics
	go func() {
		for i := 0; i < 5; i++ {
			hm.CollectMetrics(ctx)
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 2: Get metrics
	go func() {
		for i := 0; i < 5; i++ {
			hm.GetMetrics()
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// No panic means concurrent access is safe
}
