package telemetry

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Test Metrics structure
func TestMetrics_Structure(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewMetrics(registry)

	if m == nil {
		t.Error("Expected non-nil Metrics")
	}

	if m.TasksTotal == nil {
		t.Error("Expected TasksTotal to be initialized")
	}

	if m.GRPCRequestsTotal == nil {
		t.Error("Expected GRPCRequestsTotal to be initialized")
	}

	if m.ErrorsTotal == nil {
		t.Error("Expected ErrorsTotal to be initialized")
	}
}

func TestMetrics_Gauges(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewMetrics(registry)

	if m.TasksRunning == nil {
		t.Error("Expected TasksRunning gauge")
	}

	if m.AgentUptime == nil {
		t.Error("Expected AgentUptime gauge")
	}

	if m.AgentInfo == nil {
		t.Error("Expected AgentInfo gauge")
	}

	if m.GoRoutines == nil {
		t.Error("Expected GoRoutines gauge")
	}

	if m.MemoryAllocated == nil {
		t.Error("Expected MemoryAllocated gauge")
	}
}

func TestMetrics_Histograms(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewMetrics(registry)

	if m.TaskDuration == nil {
		t.Error("Expected TaskDuration histogram")
	}

	if m.GRPCDuration == nil {
		t.Error("Expected GRPCDuration histogram")
	}
}

func TestMetrics_StartTime(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewMetrics(registry)

	if m.startTime.IsZero() {
		t.Error("Expected startTime to be set")
	}

	if m.startTime.After(time.Now()) {
		t.Error("Expected startTime to be in the past")
	}
}

// Test Server structure
func TestNewServer(t *testing.T) {
	server := NewServer(9090, true)

	if server == nil {
		t.Error("Expected non-nil Server")
	}

	if server.port != 9090 {
		t.Errorf("Expected port 9090, got %d", server.port)
	}

	if !server.enabled {
		t.Error("Expected server to be enabled")
	}
}

func TestNewServer_DefaultPort(t *testing.T) {
	server := NewServer(0, true)

	if server.port != 9090 {
		t.Errorf("Expected default port 9090, got %d", server.port)
	}
}

func TestNewServer_CustomPort(t *testing.T) {
	server := NewServer(8080, true)

	if server.port != 8080 {
		t.Errorf("Expected port 8080, got %d", server.port)
	}
}

func TestNewServer_Disabled(t *testing.T) {
	server := NewServer(9090, false)

	if server.enabled {
		t.Error("Expected server to be disabled")
	}
}

func TestNewServer_HasMetrics(t *testing.T) {
	server := NewServer(9090, true)

	if server.metrics == nil {
		t.Error("Expected metrics to be initialized")
	}
}

func TestNewServer_HasRegistry(t *testing.T) {
	server := NewServer(9090, true)

	if server.registry == nil {
		t.Error("Expected registry to be initialized")
	}
}

// Test global collector
func TestInitGlobal(t *testing.T) {
	// Reset global state
	globalCollector = nil
	once = sync.Once{}

	collector := InitGlobal(9090, true)

	if collector == nil {
		t.Error("Expected non-nil global collector")
	}

	if collector.port != 9090 {
		t.Errorf("Expected port 9090, got %d", collector.port)
	}
}

func TestInitGlobal_Singleton(t *testing.T) {
	// Reset global state
	globalCollector = nil
	once = sync.Once{}

	collector1 := InitGlobal(9090, true)
	collector2 := InitGlobal(8080, true)

	// Should return the same instance
	if collector1 != collector2 {
		t.Error("Expected singleton behavior")
	}

	// Port should be from first initialization
	if collector2.port != 9090 {
		t.Error("Expected first initialization port to persist")
	}
}

func TestGetGlobal(t *testing.T) {
	// Reset and initialize
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, true)

	collector := GetGlobal()

	if collector == nil {
		t.Error("Expected non-nil collector")
	}

	if collector.port != 9090 {
		t.Errorf("Expected port 9090, got %d", collector.port)
	}
}

func TestGetGlobal_BeforeInit(t *testing.T) {
	// Reset global state
	globalCollector = nil
	once = sync.Once{}

	collector := GetGlobal()

	if collector != nil {
		t.Error("Expected nil collector before initialization")
	}
}

// Test convenience functions
func TestRecordTaskExecution_Enabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, true)

	// Should not panic
	RecordTaskExecution("group1", "task1", "success", 1.5)
}

func TestRecordTaskExecution_Disabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, false)

	// Should not panic even when disabled
	RecordTaskExecution("group1", "task1", "success", 1.5)
}

func TestRecordTaskExecution_NilCollector(t *testing.T) {
	globalCollector = nil

	// Should not panic with nil collector
	RecordTaskExecution("group1", "task1", "success", 1.5)
}

func TestRecordGRPCRequest_Enabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, true)

	// Should not panic
	RecordGRPCRequest("Execute", "success", 0.5)
}

func TestRecordGRPCRequest_Disabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, false)

	// Should not panic
	RecordGRPCRequest("Execute", "success", 0.5)
}

func TestRecordError_Enabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, true)

	// Should not panic
	RecordError("validation")
}

func TestRecordError_Disabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, false)

	// Should not panic
	RecordError("validation")
}

func TestIncrementRunningTasks_Enabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, true)

	// Should not panic
	IncrementRunningTasks()
}

func TestIncrementRunningTasks_Disabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, false)

	// Should not panic
	IncrementRunningTasks()
}

func TestDecrementRunningTasks_Enabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, true)

	// Should not panic
	DecrementRunningTasks()
}

func TestDecrementRunningTasks_Disabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, false)

	// Should not panic
	DecrementRunningTasks()
}

func TestSetAgentInfo_Enabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, true)

	// Should not panic
	SetAgentInfo("1.0.0", "linux", "amd64")
}

func TestSetAgentInfo_Disabled(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}
	InitGlobal(9090, false)

	// Should not panic
	SetAgentInfo("1.0.0", "linux", "amd64")
}

// Test Server fields
func TestServer_Structure(t *testing.T) {
	server := &Server{
		port:    9090,
		enabled: true,
	}

	if server.port != 9090 {
		t.Error("Expected port to be set")
	}

	if !server.enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestServer_DisabledServer(t *testing.T) {
	server := NewServer(9090, false)

	err := server.Start()
	if err != nil {
		t.Errorf("Disabled server should not return error: %v", err)
	}
}

// Test different port configurations
func TestNewServer_VariousPorts(t *testing.T) {
	ports := []int{8080, 8090, 9090, 9091, 3000}

	for _, port := range ports {
		server := NewServer(port, true)
		if server.port != port {
			t.Errorf("Expected port %d, got %d", port, server.port)
		}
	}
}

// Test metrics initialization
func TestMetrics_CounterVectors(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewMetrics(registry)

	// TasksTotal should have labels: status, group
	if m.TasksTotal == nil {
		t.Error("TasksTotal should be initialized")
	}

	// GRPCRequestsTotal should have labels: method, status
	if m.GRPCRequestsTotal == nil {
		t.Error("GRPCRequestsTotal should be initialized")
	}

	// ErrorsTotal should have label: type
	if m.ErrorsTotal == nil {
		t.Error("ErrorsTotal should be initialized")
	}
}

func TestMetrics_GaugeVectors(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewMetrics(registry)

	// AgentInfo should have labels: version, os, arch
	if m.AgentInfo == nil {
		t.Error("AgentInfo should be initialized")
	}
}

func TestMetrics_HistogramVectors(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewMetrics(registry)

	// TaskDuration should have labels: group, task
	if m.TaskDuration == nil {
		t.Error("TaskDuration should be initialized")
	}

	// GRPCDuration should have label: method
	if m.GRPCDuration == nil {
		t.Error("GRPCDuration should be initialized")
	}
}

// Test that multiple metrics can be created
func TestNewMetrics_Multiple(t *testing.T) {
	registry1 := prometheus.NewRegistry()
	registry2 := prometheus.NewRegistry()

	m1 := NewMetrics(registry1)
	m2 := NewMetrics(registry2)

	if m1 == nil || m2 == nil {
		t.Error("Expected both metrics to be created")
	}

	// They should be different instances
	if m1 == m2 {
		t.Error("Expected different instances")
	}
}

// Test edge cases
func TestNewServer_ZeroPort(t *testing.T) {
	server := NewServer(0, true)

	// Should default to 9090
	if server.port != 9090 {
		t.Errorf("Expected default port 9090 for zero, got %d", server.port)
	}
}

func TestNewServer_NegativePort(t *testing.T) {
	server := NewServer(-1, true)

	// Should still set the port (validation happens elsewhere)
	if server.port != -1 {
		t.Errorf("Expected port -1, got %d", server.port)
	}
}

func TestNewServer_HighPort(t *testing.T) {
	server := NewServer(65535, true)

	if server.port != 65535 {
		t.Errorf("Expected port 65535, got %d", server.port)
	}
}

// Test concurrent access to global collector
func TestInitGlobal_Concurrent(t *testing.T) {
	globalCollector = nil
	once = sync.Once{}

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(port int) {
			InitGlobal(9000+port, true)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// All should return the same instance
	collector := GetGlobal()
	if collector == nil {
		t.Error("Expected non-nil collector")
	}

	// Port should be from one of the initializations (race condition determines which)
	if collector.port < 9000 || collector.port > 9009 {
		t.Errorf("Expected port between 9000-9009, got %d", collector.port)
	}
}

func TestRecordTaskExecution_WithNilGlobal(t *testing.T) {
	globalCollector = nil

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Should not panic with nil global collector: %v", r)
		}
	}()

	RecordTaskExecution("group", "task", "success", 1.0)
}

func TestRecordGRPCRequest_WithNilGlobal(t *testing.T) {
	globalCollector = nil

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Should not panic with nil global collector: %v", r)
		}
	}()

	RecordGRPCRequest("method", "success", 1.0)
}

func TestRecordError_WithNilGlobal(t *testing.T) {
	globalCollector = nil

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Should not panic with nil global collector: %v", r)
		}
	}()

	RecordError("error_type")
}
