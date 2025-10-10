package performance

import (
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sysadmin/resources"
)

// PerformanceCollector interface para coleta de métricas de performance
type PerformanceCollector interface {
	CollectMetrics() (*PerformanceMetrics, error)
	CollectSample(duration time.Duration) (*PerformanceSample, error)
}

// PerformanceMetrics métricas completas de performance
type PerformanceMetrics struct {
	Timestamp   time.Time
	CPU         *CPUPerformance
	Memory      *MemoryPerformance
	Disk        *DiskPerformance
	Network     *NetworkPerformance
	Overall     *OverallPerformance
}

// CPUPerformance métricas de CPU
type CPUPerformance struct {
	Usage       float64
	LoadAverage [3]float64
	Cores       int
	Status      PerformanceStatus
}

// MemoryPerformance métricas de memória
type MemoryPerformance struct {
	UsagePercent float64
	Total        uint64
	Used         uint64
	Available    uint64
	Status       PerformanceStatus
}

// DiskPerformance métricas de disco
type DiskPerformance struct {
	UsagePercent float64
	TotalSpace   uint64
	UsedSpace    uint64
	FreeSpace    uint64
	Status       PerformanceStatus
}

// NetworkPerformance métricas de rede
type NetworkPerformance struct {
	TotalBytesRecv uint64
	TotalBytesSent uint64
	ActiveInterfaces int
	Status       PerformanceStatus
}

// OverallPerformance status geral de performance
type OverallPerformance struct {
	Score  int // 0-100
	Status PerformanceStatus
	Issues []string
}

// PerformanceStatus status de performance
type PerformanceStatus string

const (
	StatusExcellent PerformanceStatus = "Excellent"
	StatusGood      PerformanceStatus = "Good"
	StatusWarning   PerformanceStatus = "Warning"
	StatusCritical  PerformanceStatus = "Critical"
)

// PerformanceSample amostra de performance ao longo do tempo
type PerformanceSample struct {
	Duration    time.Duration
	Samples     []*PerformanceMetrics
	AverageCPU  float64
	MaxCPU      float64
	MinCPU      float64
	AverageRAM  float64
	MaxRAM      float64
	MinRAM      float64
}

// SystemCollector implementação padrão
type SystemCollector struct {
	monitor resources.ResourceMonitor
}

// NewCollector cria um novo collector
func NewCollector() PerformanceCollector {
	return &SystemCollector{
		monitor: resources.NewMonitor(),
	}
}

// CollectMetrics coleta métricas de performance
func (c *SystemCollector) CollectMetrics() (*PerformanceMetrics, error) {
	metrics := &PerformanceMetrics{
		Timestamp: time.Now(),
	}

	// Coleta CPU
	cpuStats, err := c.monitor.GetCPU()
	if err == nil {
		metrics.CPU = &CPUPerformance{
			Usage:       cpuStats.Usage,
			LoadAverage: cpuStats.LoadAverage,
			Cores:       cpuStats.Cores,
			Status:      getCPUStatus(cpuStats.Usage),
		}
	}

	// Coleta Memory
	memStats, err := c.monitor.GetMemory()
	if err == nil {
		metrics.Memory = &MemoryPerformance{
			UsagePercent: memStats.UsagePercent,
			Total:        memStats.Total,
			Used:         memStats.Used,
			Available:    memStats.Available,
			Status:       getMemoryStatus(memStats.UsagePercent),
		}
	}

	// Coleta Disk
	diskStats, err := c.monitor.GetDisk()
	if err == nil && len(diskStats) > 0 {
		var totalSpace, usedSpace, freeSpace uint64
		for _, disk := range diskStats {
			totalSpace += disk.Total
			usedSpace += disk.Used
			freeSpace += disk.Available
		}

		usagePercent := 0.0
		if totalSpace > 0 {
			usagePercent = 100.0 * float64(usedSpace) / float64(totalSpace)
		}

		metrics.Disk = &DiskPerformance{
			UsagePercent: usagePercent,
			TotalSpace:   totalSpace,
			UsedSpace:    usedSpace,
			FreeSpace:    freeSpace,
			Status:       getDiskStatus(usagePercent),
		}
	}

	// Coleta Network
	netStats, err := c.monitor.GetNetwork()
	if err == nil {
		var totalRecv, totalSent uint64
		activeInterfaces := 0

		for _, net := range netStats {
			totalRecv += net.BytesRecv
			totalSent += net.BytesSent
			if net.BytesRecv > 0 || net.BytesSent > 0 {
				activeInterfaces++
			}
		}

		metrics.Network = &NetworkPerformance{
			TotalBytesRecv:   totalRecv,
			TotalBytesSent:   totalSent,
			ActiveInterfaces: activeInterfaces,
			Status:           StatusGood, // Network status is harder to determine
		}
	}

	// Calcula overall performance
	metrics.Overall = c.calculateOverallPerformance(metrics)

	return metrics, nil
}

// CollectSample coleta amostras ao longo do tempo
func (c *SystemCollector) CollectSample(duration time.Duration) (*PerformanceSample, error) {
	sample := &PerformanceSample{
		Duration: duration,
		Samples:  []*PerformanceMetrics{},
	}

	// Coleta inicial
	initial, err := c.CollectMetrics()
	if err != nil {
		return nil, err
	}
	sample.Samples = append(sample.Samples, initial)

	// Coleta amostras durante a duração
	interval := time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeout := time.After(duration)

	for {
		select {
		case <-timeout:
			// Calcula estatísticas
			c.calculateSampleStats(sample)
			return sample, nil

		case <-ticker.C:
			metrics, err := c.CollectMetrics()
			if err == nil {
				sample.Samples = append(sample.Samples, metrics)
			}
		}
	}
}

// calculateSampleStats calcula estatísticas da amostra
func (c *SystemCollector) calculateSampleStats(sample *PerformanceSample) {
	if len(sample.Samples) == 0 {
		return
	}

	var sumCPU, sumRAM float64
	sample.MaxCPU = sample.Samples[0].CPU.Usage
	sample.MinCPU = sample.Samples[0].CPU.Usage
	sample.MaxRAM = sample.Samples[0].Memory.UsagePercent
	sample.MinRAM = sample.Samples[0].Memory.UsagePercent

	for _, s := range sample.Samples {
		if s.CPU != nil {
			sumCPU += s.CPU.Usage
			if s.CPU.Usage > sample.MaxCPU {
				sample.MaxCPU = s.CPU.Usage
			}
			if s.CPU.Usage < sample.MinCPU {
				sample.MinCPU = s.CPU.Usage
			}
		}

		if s.Memory != nil {
			sumRAM += s.Memory.UsagePercent
			if s.Memory.UsagePercent > sample.MaxRAM {
				sample.MaxRAM = s.Memory.UsagePercent
			}
			if s.Memory.UsagePercent < sample.MinRAM {
				sample.MinRAM = s.Memory.UsagePercent
			}
		}
	}

	count := float64(len(sample.Samples))
	sample.AverageCPU = sumCPU / count
	sample.AverageRAM = sumRAM / count
}

// calculateOverallPerformance calcula performance geral
func (c *SystemCollector) calculateOverallPerformance(metrics *PerformanceMetrics) *OverallPerformance {
	overall := &OverallPerformance{
		Score:  100,
		Status: StatusExcellent,
		Issues: []string{},
	}

	// Reduz score baseado em métricas
	if metrics.CPU != nil {
		if metrics.CPU.Usage > 90 {
			overall.Score -= 30
			overall.Issues = append(overall.Issues, "High CPU usage")
		} else if metrics.CPU.Usage > 70 {
			overall.Score -= 15
		}
	}

	if metrics.Memory != nil {
		if metrics.Memory.UsagePercent > 90 {
			overall.Score -= 30
			overall.Issues = append(overall.Issues, "High memory usage")
		} else if metrics.Memory.UsagePercent > 80 {
			overall.Score -= 15
		}
	}

	if metrics.Disk != nil {
		if metrics.Disk.UsagePercent > 90 {
			overall.Score -= 20
			overall.Issues = append(overall.Issues, "High disk usage")
		} else if metrics.Disk.UsagePercent > 80 {
			overall.Score -= 10
		}
	}

	// Determina status baseado no score
	if overall.Score >= 80 {
		overall.Status = StatusExcellent
	} else if overall.Score >= 60 {
		overall.Status = StatusGood
	} else if overall.Score >= 40 {
		overall.Status = StatusWarning
	} else {
		overall.Status = StatusCritical
	}

	if len(overall.Issues) == 0 {
		overall.Issues = append(overall.Issues, "No issues detected")
	}

	return overall
}

// getCPUStatus determina status do CPU
func getCPUStatus(usage float64) PerformanceStatus {
	if usage < 50 {
		return StatusExcellent
	} else if usage < 70 {
		return StatusGood
	} else if usage < 90 {
		return StatusWarning
	}
	return StatusCritical
}

// getMemoryStatus determina status da memória
func getMemoryStatus(usage float64) PerformanceStatus {
	if usage < 60 {
		return StatusExcellent
	} else if usage < 80 {
		return StatusGood
	} else if usage < 90 {
		return StatusWarning
	}
	return StatusCritical
}

// getDiskStatus determina status do disco
func getDiskStatus(usage float64) PerformanceStatus {
	if usage < 70 {
		return StatusExcellent
	} else if usage < 85 {
		return StatusGood
	} else if usage < 95 {
		return StatusWarning
	}
	return StatusCritical
}

// FormatBytes formata bytes
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
