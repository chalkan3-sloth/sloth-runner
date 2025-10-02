package core

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yuin/gopher-lua"
)

func TestNewMonitoringModule(t *testing.T) {
	module := NewMonitoringModule()
	
	if module == nil {
		t.Fatal("Expected module to be created")
	}
	
	if module.info.Name != "monitor" {
		t.Errorf("Expected name 'monitor', got '%s'", module.info.Name)
	}
	
	if module.metrics == nil {
		t.Error("Expected metrics map to be initialized")
	}
}

func TestMonitoringModuleLoader(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	if err := L.DoString(`monitor = require("monitor")`); err != nil {
		t.Fatalf("Failed to load module: %v", err)
	}
	
	functions := []string{
		"counter_inc", "counter_add", "gauge_set", "gauge_inc", "gauge_dec",
		"timer_start", "timer_end", "histogram_observe", "get_metric",
		"list_metrics", "reset_metric", "clear_all", "export_prometheus",
		"export_json", "system_metrics", "memory_stats",
	}
	
	for _, fn := range functions {
		if err := L.DoString(`assert(monitor.` + fn + ` ~= nil, "` + fn + ` function not found")`); err != nil {
			t.Errorf("Function %s not found: %v", fn, err)
		}
	}
}

func TestCounterInc(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local val1 = monitor.counter_inc("test_counter")
		local val2 = monitor.counter_inc("test_counter")
		local val3 = monitor.counter_inc("test_counter")
		
		assert(val1 == 1, "First increment should be 1")
		assert(val2 == 2, "Second increment should be 2")
		assert(val3 == 3, "Third increment should be 3")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Verify metric exists
	key := module.generateKey("test_counter", nil)
	metric, exists := module.metrics[key]
	if !exists {
		t.Error("Expected metric to exist")
	}
	
	if metric.Value != 3 {
		t.Errorf("Expected value 3, got %f", metric.Value)
	}
	
	if metric.Type != "counter" {
		t.Errorf("Expected type 'counter', got '%s'", metric.Type)
	}
}

func TestCounterAdd(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local val = monitor.counter_add("test_counter", 5)
		assert(val == 5)
		val = monitor.counter_add("test_counter", 10)
		assert(val == 15)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGaugeSet(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local val = monitor.gauge_set("cpu_usage", 75.5)
		assert(val == 75.5)
		val = monitor.gauge_set("cpu_usage", 80.0)
		assert(val == 80.0)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	key := module.generateKey("cpu_usage", nil)
	metric := module.metrics[key]
	
	if metric.Type != "gauge" {
		t.Errorf("Expected type 'gauge', got '%s'", metric.Type)
	}
	
	if metric.Value != 80.0 {
		t.Errorf("Expected value 80.0, got %f", metric.Value)
	}
}

func TestGaugeInc(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.gauge_set("temp", 20)
		local val = monitor.gauge_inc("temp")
		assert(val == 21)
		val = monitor.gauge_inc("temp", 5)
		assert(val == 26)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGaugeDec(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.gauge_set("temp", 30)
		local val = monitor.gauge_dec("temp")
		assert(val == 29)
		val = monitor.gauge_dec("temp", 5)
		assert(val == 24)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestTimer(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local timer = monitor.timer_start("operation_time")
		assert(timer ~= nil)
		
		-- Simulate some work
		local sum = 0
		for i = 1, 1000000 do
			sum = sum + i
		end
		
		local duration = monitor.timer_end(timer)
		assert(duration ~= nil and duration >= 0)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestTimerNotStarted(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local duration, err = monitor.timer_end("invalid_timer")
		assert(duration == nil)
		assert(err ~= nil)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestHistogramObserve(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.histogram_observe("response_time", 100)
		monitor.histogram_observe("response_time", 200)
		monitor.histogram_observe("response_time", 300)
		
		local metric = monitor.get_metric("response_time")
		assert(metric.value == 200) -- Average of 100, 200, 300
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	key := module.generateKey("response_time", nil)
	metric := module.metrics[key]
	
	if metric.Type != "histogram" {
		t.Errorf("Expected type 'histogram', got '%s'", metric.Type)
	}
	
	count := metric.Metadata["count"].(float64)
	if count != 3 {
		t.Errorf("Expected count 3, got %f", count)
	}
	
	sum := metric.Metadata["sum"].(float64)
	if sum != 600 {
		t.Errorf("Expected sum 600, got %f", sum)
	}
}

func TestMetricsWithLabels(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.counter_inc("http_requests", {method = "GET", status = "200"})
		monitor.counter_inc("http_requests", {method = "GET", status = "200"})
		monitor.counter_inc("http_requests", {method = "POST", status = "201"})
		
		-- Verify metrics were created (check via list)
		local metrics = monitor.list_metrics()
		assert(#metrics >= 2, "Should have at least 2 metrics")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Verify directly in Go
	module.mutex.RLock()
	defer module.mutex.RUnlock()
	
	if len(module.metrics) < 2 {
		t.Errorf("Expected at least 2 metrics, got %d", len(module.metrics))
	}
}

func TestGetMetric(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.counter_inc("test_metric")
		
		local metric = monitor.get_metric("test_metric")
		assert(metric ~= nil)
		assert(metric.name == "test_metric")
		assert(metric.type == "counter")
		assert(metric.value == 1)
		assert(metric.last_updated ~= nil)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGetMetricNonExistent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local metric = monitor.get_metric("nonexistent")
		assert(metric == nil)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestListMetrics(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.counter_inc("metric1")
		monitor.gauge_set("metric2", 42)
		monitor.counter_inc("metric3")
		
		local metrics = monitor.list_metrics()
		assert(#metrics == 3)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestResetMetric(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.counter_inc("test")
		monitor.counter_inc("test")
		monitor.counter_inc("test")
		
		local metric = monitor.get_metric("test")
		assert(metric.value == 3)
		
		local result = monitor.reset_metric("test")
		assert(result == true)
		
		metric = monitor.get_metric("test")
		assert(metric.value == 0)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestClearAll(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.counter_inc("metric1")
		monitor.gauge_set("metric2", 42)
		monitor.counter_inc("metric3")
		
		local metrics = monitor.list_metrics()
		assert(#metrics == 3)
		
		monitor.clear_all()
		
		metrics = monitor.list_metrics()
		assert(#metrics == 0)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestExportPrometheus(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.counter_inc("http_requests_total", {method = "GET"})
		monitor.gauge_set("cpu_usage", 75.5)
		
		local output = monitor.export_prometheus()
		assert(output ~= nil and output ~= "")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Get the output
	if err := L.DoString(`return monitor.export_prometheus()`); err != nil {
		t.Fatalf("Failed to get output: %v", err)
	}
	
	output := L.Get(-1).String()
	L.Pop(1)
	
	// Verify Prometheus format
	if !strings.Contains(output, "# HELP") {
		t.Error("Expected output to contain HELP line")
	}
	
	if !strings.Contains(output, "# TYPE") {
		t.Error("Expected output to contain TYPE line")
	}
	
	if !strings.Contains(output, "http_requests_total") {
		t.Error("Expected output to contain http_requests_total metric")
	}
	
	if !strings.Contains(output, "cpu_usage") {
		t.Error("Expected output to contain cpu_usage metric")
	}
}

func TestExportJSON(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		monitor.counter_inc("requests")
		monitor.gauge_set("temperature", 25.5)
		
		local output = monitor.export_json()
		assert(output ~= nil and output ~= "")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Get the output and verify it's valid JSON
	if err := L.DoString(`return monitor.export_json()`); err != nil {
		t.Fatalf("Failed to get output: %v", err)
	}
	
	output := L.Get(-1).String()
	L.Pop(1)
	
	var metrics map[string]interface{}
	if err := json.Unmarshal([]byte(output), &metrics); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}
}

func TestSystemMetrics(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local metrics = monitor.system_metrics()
		
		assert(metrics.goroutines ~= nil)
		assert(metrics.cpu_count ~= nil)
		assert(metrics.memory_alloc ~= nil)
		assert(metrics.memory_total_alloc ~= nil)
		assert(metrics.memory_sys ~= nil)
		assert(metrics.gc_count ~= nil)
		
		assert(metrics.cpu_count > 0)
		assert(metrics.memory_alloc >= 0)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestMemoryStats(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewMonitoringModule()
	L.PreloadModule("monitor", module.Loader)
	
	script := `
		monitor = require("monitor")
		local stats = monitor.memory_stats()
		
		assert(stats.alloc ~= nil)
		assert(stats.total_alloc ~= nil)
		assert(stats.sys ~= nil)
		assert(stats.heap_alloc ~= nil)
		assert(stats.heap_sys ~= nil)
		assert(stats.gc_count ~= nil)
		
		assert(stats.alloc >= 0)
		assert(stats.total_alloc >= 0)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestMetricHistory(t *testing.T) {
	module := NewMonitoringModule()
	
	metric := &Metric{
		Name:    "test",
		Type:    "counter",
		Value:   0,
		History: []MetricPoint{},
	}
	
	// Add some points
	for i := 0; i < 10; i++ {
		metric.Value = float64(i)
		module.addToHistory(metric)
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}
	
	if len(metric.History) != 10 {
		t.Errorf("Expected 10 history points, got %d", len(metric.History))
	}
	
	// Add more to test limit
	for i := 0; i < 100; i++ {
		metric.Value = float64(i + 10)
		module.addToHistory(metric)
	}
	
	if len(metric.History) != 100 {
		t.Errorf("Expected history to be limited to 100, got %d", len(metric.History))
	}
}

func TestGenerateKey(t *testing.T) {
	module := NewMonitoringModule()
	
	tests := []struct {
		name     string
		metricName string
		labels   map[string]string
		wantKey  bool // just check if key is generated
	}{
		{"no labels", "metric1", nil, true},
		{"with labels", "metric2", map[string]string{"method": "GET"}, true},
		{"multiple labels", "metric3", map[string]string{"method": "POST", "status": "200"}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := module.generateKey(tt.metricName, tt.labels)
			
			if key == "" {
				t.Error("Expected non-empty key")
			}
			
			if len(tt.labels) == 0 && key != tt.metricName {
				t.Errorf("Expected key to be metric name for no labels, got '%s'", key)
			}
			
			if len(tt.labels) > 0 && !strings.Contains(key, "{") {
				t.Error("Expected key to contain labels")
			}
		})
	}
}

func TestMonitoringModuleInfo(t *testing.T) {
	module := NewMonitoringModule()
	info := module.Info()
	
	if info.Name != "monitor" {
		t.Errorf("Expected name 'monitor', got '%s'", info.Name)
	}
	
	if info.Category != "core" {
		t.Errorf("Expected category 'core', got '%s'", info.Category)
	}
	
	if info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", info.Version)
	}
}
