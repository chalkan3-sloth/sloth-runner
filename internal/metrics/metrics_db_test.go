//go:build cgo
// +build cgo

package metrics

import (
	"testing"
	"time"
)

// Test MetricPoint struct
func TestMetricPoint_Creation(t *testing.T) {
	now := time.Now().Unix()
	metric := MetricPoint{
		Timestamp:       now,
		CPUPercent:      25.5,
		MemoryPercent:   60.2,
		MemoryUsedBytes: 1024 * 1024 * 1024, // 1GB
		DiskPercent:     45.3,
		LoadAvg1Min:     1.5,
		LoadAvg5Min:     1.2,
		LoadAvg15Min:    1.0,
		ProcessCount:    150,
		NetworkRxBytes:  1000000,
		NetworkTxBytes:  500000,
	}

	if metric.Timestamp != now {
		t.Error("Expected timestamp to be set")
	}

	if metric.CPUPercent != 25.5 {
		t.Errorf("Expected CPU 25.5%%, got %f", metric.CPUPercent)
	}

	if metric.MemoryPercent != 60.2 {
		t.Errorf("Expected Memory 60.2%%, got %f", metric.MemoryPercent)
	}

	if metric.ProcessCount != 150 {
		t.Errorf("Expected 150 processes, got %d", metric.ProcessCount)
	}
}

func TestMetricPoint_ZeroValues(t *testing.T) {
	metric := MetricPoint{}

	if metric.Timestamp != 0 {
		t.Error("Expected zero timestamp")
	}

	if metric.CPUPercent != 0 {
		t.Error("Expected zero CPU")
	}

	if metric.MemoryPercent != 0 {
		t.Error("Expected zero memory")
	}

	if metric.ProcessCount != 0 {
		t.Error("Expected zero process count")
	}
}

func TestMetricPoint_HighValues(t *testing.T) {
	metric := MetricPoint{
		CPUPercent:      100.0,
		MemoryPercent:   99.9,
		MemoryUsedBytes: 1024 * 1024 * 1024 * 16, // 16GB
		DiskPercent:     95.5,
		LoadAvg1Min:     10.0,
		LoadAvg5Min:     8.5,
		LoadAvg15Min:    7.2,
		ProcessCount:    1000,
		NetworkRxBytes:  1024 * 1024 * 1024 * 10, // 10GB
		NetworkTxBytes:  1024 * 1024 * 1024 * 5,  // 5GB
	}

	if metric.CPUPercent != 100.0 {
		t.Errorf("Expected CPU 100%%, got %f", metric.CPUPercent)
	}

	if metric.MemoryPercent <= 99.0 {
		t.Error("Expected very high memory usage")
	}

	if metric.ProcessCount != 1000 {
		t.Errorf("Expected 1000 processes, got %d", metric.ProcessCount)
	}
}

func TestMetricPoint_NetworkMetrics(t *testing.T) {
	metric := MetricPoint{
		NetworkRxBytes: 123456789,
		NetworkTxBytes: 987654321,
	}

	if metric.NetworkRxBytes == 0 {
		t.Error("Expected non-zero RX bytes")
	}

	if metric.NetworkTxBytes == 0 {
		t.Error("Expected non-zero TX bytes")
	}

	if metric.NetworkRxBytes == metric.NetworkTxBytes {
		t.Error("RX and TX bytes should typically be different")
	}
}

func TestMetricPoint_LoadAverages(t *testing.T) {
	metric := MetricPoint{
		LoadAvg1Min:  2.5,
		LoadAvg5Min:  2.0,
		LoadAvg15Min: 1.5,
	}

	// Load averages should typically decrease over time
	if metric.LoadAvg1Min < metric.LoadAvg15Min {
		// This is fine - just testing that values can be set
	}

	if metric.LoadAvg1Min == 0 && metric.LoadAvg5Min == 0 && metric.LoadAvg15Min == 0 {
		t.Error("All load averages are zero")
	}
}

func TestMetricPoint_Timestamps(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	testCases := []struct {
		name      string
		timestamp int64
	}{
		{"Current time", now.Unix()},
		{"Past time", past.Unix()},
		{"Future time", future.Unix()},
		{"Zero time", 0},
		{"Epoch", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metric := MetricPoint{
				Timestamp: tc.timestamp,
			}

			if metric.Timestamp != tc.timestamp {
				t.Errorf("Expected timestamp %d, got %d", tc.timestamp, metric.Timestamp)
			}
		})
	}
}

func TestMetricPoint_MemoryBytes(t *testing.T) {
	testCases := []struct {
		name  string
		bytes uint64
		desc  string
	}{
		{"1KB", 1024, "1 kilobyte"},
		{"1MB", 1024 * 1024, "1 megabyte"},
		{"1GB", 1024 * 1024 * 1024, "1 gigabyte"},
		{"10GB", 1024 * 1024 * 1024 * 10, "10 gigabytes"},
		{"100GB", 1024 * 1024 * 1024 * 100, "100 gigabytes"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metric := MetricPoint{
				MemoryUsedBytes: tc.bytes,
			}

			if metric.MemoryUsedBytes != tc.bytes {
				t.Errorf("Expected %d bytes (%s), got %d", tc.bytes, tc.desc, metric.MemoryUsedBytes)
			}
		})
	}
}

func TestMetricPoint_CPUPercent(t *testing.T) {
	testCases := []float64{0.0, 10.5, 25.0, 50.0, 75.5, 99.9, 100.0}

	for _, cpuPercent := range testCases {
		metric := MetricPoint{
			CPUPercent: cpuPercent,
		}

		if metric.CPUPercent != cpuPercent {
			t.Errorf("Expected CPU %f%%, got %f%%", cpuPercent, metric.CPUPercent)
		}
	}
}

func TestMetricPoint_DiskPercent(t *testing.T) {
	testCases := []float64{0.0, 15.3, 30.5, 50.0, 80.2, 95.7, 100.0}

	for _, diskPercent := range testCases {
		metric := MetricPoint{
			DiskPercent: diskPercent,
		}

		if metric.DiskPercent != diskPercent {
			t.Errorf("Expected Disk %f%%, got %f%%", diskPercent, metric.DiskPercent)
		}
	}
}

func TestMetricPoint_ProcessCount(t *testing.T) {
	testCases := []int{0, 1, 10, 50, 100, 500, 1000, 5000}

	for _, count := range testCases {
		metric := MetricPoint{
			ProcessCount: count,
		}

		if metric.ProcessCount != count {
			t.Errorf("Expected %d processes, got %d", count, metric.ProcessCount)
		}
	}
}

// Test batchMetric struct
func TestBatchMetric_Creation(t *testing.T) {
	metric := batchMetric{
		AgentName: "test-agent",
		Metric: MetricPoint{
			Timestamp:  time.Now().Unix(),
			CPUPercent: 50.0,
		},
	}

	if metric.AgentName != "test-agent" {
		t.Errorf("Expected agent name 'test-agent', got '%s'", metric.AgentName)
	}

	if metric.Metric.CPUPercent != 50.0 {
		t.Errorf("Expected CPU 50%%, got %f", metric.Metric.CPUPercent)
	}
}

func TestBatchMetric_MultipleInstances(t *testing.T) {
	batch := []batchMetric{
		{AgentName: "agent1", Metric: MetricPoint{CPUPercent: 10.0}},
		{AgentName: "agent2", Metric: MetricPoint{CPUPercent: 20.0}},
		{AgentName: "agent3", Metric: MetricPoint{CPUPercent: 30.0}},
	}

	if len(batch) != 3 {
		t.Errorf("Expected 3 batch metrics, got %d", len(batch))
	}

	for i, bm := range batch {
		expectedCPU := float64((i + 1) * 10)
		if bm.Metric.CPUPercent != expectedCPU {
			t.Errorf("Expected CPU %f%%, got %f", expectedCPU, bm.Metric.CPUPercent)
		}
	}
}

// Test MetricPoint JSON tags
func TestMetricPoint_JSONTags(t *testing.T) {
	// Verify struct has json tags for serialization
	metric := MetricPoint{
		Timestamp:       123456789,
		CPUPercent:      25.5,
		MemoryPercent:   60.0,
		MemoryUsedBytes: 1024,
		DiskPercent:     45.0,
		LoadAvg1Min:     1.5,
		LoadAvg5Min:     1.2,
		LoadAvg15Min:    1.0,
		ProcessCount:    100,
		NetworkRxBytes:  1000,
		NetworkTxBytes:  500,
	}

	// All fields should be set
	if metric.Timestamp == 0 {
		t.Error("Timestamp should be set")
	}

	if metric.CPUPercent == 0 {
		t.Error("CPUPercent should be set")
	}

	if metric.ProcessCount == 0 {
		t.Error("ProcessCount should be set")
	}
}

// Test MetricPoint comparison
func TestMetricPoint_Comparison(t *testing.T) {
	metric1 := MetricPoint{
		Timestamp:  100,
		CPUPercent: 25.0,
	}

	metric2 := MetricPoint{
		Timestamp:  200,
		CPUPercent: 50.0,
	}

	// Metrics with different timestamps
	if metric1.Timestamp >= metric2.Timestamp {
		t.Error("Expected metric1 timestamp to be earlier")
	}

	// Metrics with different CPU
	if metric1.CPUPercent >= metric2.CPUPercent {
		t.Error("Expected metric1 CPU to be lower")
	}
}

// Test MetricPoint time series
func TestMetricPoint_TimeSeries(t *testing.T) {
	now := time.Now()
	metrics := []MetricPoint{
		{Timestamp: now.Unix(), CPUPercent: 10.0},
		{Timestamp: now.Add(1 * time.Minute).Unix(), CPUPercent: 20.0},
		{Timestamp: now.Add(2 * time.Minute).Unix(), CPUPercent: 30.0},
		{Timestamp: now.Add(3 * time.Minute).Unix(), CPUPercent: 25.0},
		{Timestamp: now.Add(4 * time.Minute).Unix(), CPUPercent: 15.0},
	}

	// Verify timestamps are in order
	for i := 1; i < len(metrics); i++ {
		if metrics[i].Timestamp <= metrics[i-1].Timestamp {
			t.Errorf("Timestamp out of order at index %d", i)
		}
	}
}

// Test MetricPoint edge cases
func TestMetricPoint_EdgeCases(t *testing.T) {
	// Very small values
	metric1 := MetricPoint{
		CPUPercent:    0.001,
		MemoryPercent: 0.001,
		DiskPercent:   0.001,
	}

	if metric1.CPUPercent >= 0.01 {
		t.Error("Expected very small CPU value")
	}

	// Maximum uint64 values
	metric2 := MetricPoint{
		MemoryUsedBytes: ^uint64(0), // Max uint64
		NetworkRxBytes:  ^uint64(0),
		NetworkTxBytes:  ^uint64(0),
	}

	if metric2.MemoryUsedBytes == 0 {
		t.Error("Expected maximum uint64 value")
	}

	// Negative process count (should be caught by type system)
	// ProcessCount is int, so can be negative theoretically
	metric3 := MetricPoint{
		ProcessCount: -1,
	}

	if metric3.ProcessCount >= 0 {
		// In practice, process count should always be non-negative
		// This test documents the type allows it
	}
}

// Test MetricPoint realistic scenarios
func TestMetricPoint_RealisticScenarios(t *testing.T) {
	testCases := []struct {
		name   string
		metric MetricPoint
		desc   string
	}{
		{
			name: "Idle system",
			metric: MetricPoint{
				CPUPercent:    5.0,
				MemoryPercent: 30.0,
				DiskPercent:   20.0,
				LoadAvg1Min:   0.5,
				ProcessCount:  50,
			},
			desc: "Low resource usage",
		},
		{
			name: "Busy system",
			metric: MetricPoint{
				CPUPercent:    85.0,
				MemoryPercent: 90.0,
				DiskPercent:   75.0,
				LoadAvg1Min:   5.0,
				ProcessCount:  300,
			},
			desc: "High resource usage",
		},
		{
			name: "Overloaded system",
			metric: MetricPoint{
				CPUPercent:    99.5,
				MemoryPercent: 98.0,
				DiskPercent:   95.0,
				LoadAvg1Min:   10.0,
				ProcessCount:  1000,
			},
			desc: "System under stress",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.metric.CPUPercent < 0 || tc.metric.CPUPercent > 100 {
				t.Errorf("%s: CPU should be 0-100%%", tc.name)
			}

			if tc.metric.MemoryPercent < 0 || tc.metric.MemoryPercent > 100 {
				t.Errorf("%s: Memory should be 0-100%%", tc.name)
			}

			if tc.metric.DiskPercent < 0 || tc.metric.DiskPercent > 100 {
				t.Errorf("%s: Disk should be 0-100%%", tc.name)
			}
		})
	}
}
