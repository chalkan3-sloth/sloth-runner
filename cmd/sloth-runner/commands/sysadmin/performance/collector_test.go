package performance

import (
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	collector := NewCollector()
	if collector == nil {
		t.Fatal("NewCollector() returned nil")
	}

	_, ok := collector.(*SystemCollector)
	if !ok {
		t.Error("NewCollector() did not return *SystemCollector")
	}
}

func TestCollectMetrics(t *testing.T) {
	collector := NewCollector()

	metrics, err := collector.CollectMetrics()
	if err != nil {
		t.Fatalf("CollectMetrics() failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("CollectMetrics() returned nil metrics")
	}

	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	// Verifica métricas de CPU
	if metrics.CPU != nil {
		if metrics.CPU.Usage < 0 || metrics.CPU.Usage > 100 {
			t.Errorf("CPU usage out of range: %f", metrics.CPU.Usage)
		}
		if metrics.CPU.Cores <= 0 {
			t.Errorf("Invalid CPU cores: %d", metrics.CPU.Cores)
		}
	}

	// Verifica métricas de memória
	if metrics.Memory != nil {
		if metrics.Memory.UsagePercent < 0 || metrics.Memory.UsagePercent > 100 {
			t.Errorf("Memory usage out of range: %f", metrics.Memory.UsagePercent)
		}
		if metrics.Memory.Total == 0 {
			t.Error("Total memory is zero")
		}
	}

	// Verifica overall performance
	if metrics.Overall == nil {
		t.Fatal("Overall performance is nil")
	}

	if metrics.Overall.Score < 0 || metrics.Overall.Score > 100 {
		t.Errorf("Overall score out of range: %d", metrics.Overall.Score)
	}

	if len(metrics.Overall.Issues) == 0 {
		t.Error("Issues list is empty")
	}
}

func TestCollectSample(t *testing.T) {
	collector := NewCollector()

	duration := 3 * time.Second
	sample, err := collector.CollectSample(duration)
	if err != nil {
		t.Fatalf("CollectSample() failed: %v", err)
	}

	if sample == nil {
		t.Fatal("CollectSample() returned nil")
	}

	if sample.Duration != duration {
		t.Errorf("Expected duration %v, got %v", duration, sample.Duration)
	}

	if len(sample.Samples) < 2 {
		t.Errorf("Expected at least 2 samples, got %d", len(sample.Samples))
	}

	// Verifica estatísticas
	if sample.AverageCPU < 0 || sample.AverageCPU > 100 {
		t.Errorf("Average CPU out of range: %f", sample.AverageCPU)
	}

	if sample.MaxCPU < sample.MinCPU {
		t.Errorf("Max CPU (%f) less than Min CPU (%f)", sample.MaxCPU, sample.MinCPU)
	}

	if sample.AverageRAM < 0 || sample.AverageRAM > 100 {
		t.Errorf("Average RAM out of range: %f", sample.AverageRAM)
	}

	if sample.MaxRAM < sample.MinRAM {
		t.Errorf("Max RAM (%f) less than Min RAM (%f)", sample.MaxRAM, sample.MinRAM)
	}
}

func TestGetCPUStatus(t *testing.T) {
	tests := []struct {
		usage    float64
		expected PerformanceStatus
	}{
		{30.0, StatusExcellent},
		{60.0, StatusGood},
		{80.0, StatusWarning},
		{95.0, StatusCritical},
	}

	for _, tt := range tests {
		result := getCPUStatus(tt.usage)
		if result != tt.expected {
			t.Errorf("getCPUStatus(%f) = %v, want %v", tt.usage, result, tt.expected)
		}
	}
}

func TestGetMemoryStatus(t *testing.T) {
	tests := []struct {
		usage    float64
		expected PerformanceStatus
	}{
		{50.0, StatusExcellent},
		{70.0, StatusGood},
		{85.0, StatusWarning},
		{95.0, StatusCritical},
	}

	for _, tt := range tests {
		result := getMemoryStatus(tt.usage)
		if result != tt.expected {
			t.Errorf("getMemoryStatus(%f) = %v, want %v", tt.usage, result, tt.expected)
		}
	}
}

func TestGetDiskStatus(t *testing.T) {
	tests := []struct {
		usage    float64
		expected PerformanceStatus
	}{
		{50.0, StatusExcellent},
		{75.0, StatusGood},
		{90.0, StatusWarning},
		{98.0, StatusCritical},
	}

	for _, tt := range tests {
		result := getDiskStatus(tt.usage)
		if result != tt.expected {
			t.Errorf("getDiskStatus(%f) = %v, want %v", tt.usage, result, tt.expected)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected string
	}{
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		result := FormatBytes(tt.bytes)
		if result != tt.expected {
			t.Errorf("FormatBytes(%d) = %s, want %s", tt.bytes, result, tt.expected)
		}
	}
}

func TestCalculateOverallPerformance(t *testing.T) {
	collector := &SystemCollector{}

	tests := []struct {
		name     string
		metrics  *PerformanceMetrics
		minScore int
		maxScore int
	}{
		{
			name: "Low usage - excellent",
			metrics: &PerformanceMetrics{
				CPU: &CPUPerformance{
					Usage: 30.0,
				},
				Memory: &MemoryPerformance{
					UsagePercent: 50.0,
				},
				Disk: &DiskPerformance{
					UsagePercent: 40.0,
				},
			},
			minScore: 90,
			maxScore: 100,
		},
		{
			name: "High usage - critical",
			metrics: &PerformanceMetrics{
				CPU: &CPUPerformance{
					Usage: 95.0,
				},
				Memory: &MemoryPerformance{
					UsagePercent: 95.0,
				},
				Disk: &DiskPerformance{
					UsagePercent: 95.0,
				},
			},
			minScore: 0,
			maxScore: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overall := collector.calculateOverallPerformance(tt.metrics)

			if overall.Score < tt.minScore || overall.Score > tt.maxScore {
				t.Errorf("Score %d not in expected range [%d, %d]", overall.Score, tt.minScore, tt.maxScore)
			}

			if len(overall.Issues) == 0 {
				t.Error("Issues list is empty")
			}
		})
	}
}
