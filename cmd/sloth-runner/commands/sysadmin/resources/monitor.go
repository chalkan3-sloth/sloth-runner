package resources

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ResourceMonitor interface para monitoramento de recursos do sistema
type ResourceMonitor interface {
	GetCPU() (*CPUStats, error)
	GetMemory() (*MemoryStats, error)
	GetDisk() ([]*DiskStats, error)
	GetNetwork() ([]*NetworkStats, error)
	GetProcesses(limit int) ([]*ProcessStats, error)
}

// CPUStats contém estatísticas de CPU
type CPUStats struct {
	Usage       float64   // Porcentagem total
	PerCore     []float64 // Por core
	LoadAverage [3]float64 // 1, 5, 15 min
	Cores       int
}

// MemoryStats contém estatísticas de memória
type MemoryStats struct {
	Total       uint64
	Used        uint64
	Free        uint64
	Available   uint64
	SwapTotal   uint64
	SwapUsed    uint64
	SwapFree    uint64
	UsagePercent float64
}

// DiskStats contém estatísticas de disco
type DiskStats struct {
	Filesystem  string
	MountPoint  string
	Total       uint64
	Used        uint64
	Available   uint64
	UsagePercent float64
}

// NetworkStats contém estatísticas de rede
type NetworkStats struct {
	Interface    string
	BytesRecv    uint64
	BytesSent    uint64
	PacketsRecv  uint64
	PacketsSent  uint64
	ErrorsRecv   uint64
	ErrorsSent   uint64
}

// ProcessStats contém estatísticas de processo
type ProcessStats struct {
	PID         int
	Name        string
	CPUPercent  float64
	MemoryBytes uint64
	MemoryPercent float64
}

// SystemMonitor implementação padrão do ResourceMonitor
type SystemMonitor struct{}

// NewMonitor cria um novo monitor de recursos
func NewMonitor() ResourceMonitor {
	return &SystemMonitor{}
}

// GetCPU retorna estatísticas de CPU
func (m *SystemMonitor) GetCPU() (*CPUStats, error) {
	stats := &CPUStats{
		Cores: runtime.NumCPU(),
	}

	// Get load average
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		loadavg, err := getLoadAverage()
		if err == nil {
			stats.LoadAverage = loadavg
		}
	}

	// Get CPU usage
	usage, err := getCPUUsage()
	if err != nil {
		return nil, err
	}
	stats.Usage = usage

	// Get per-core usage (simplified - retorna usage igual para todos os cores)
	stats.PerCore = make([]float64, stats.Cores)
	for i := range stats.PerCore {
		stats.PerCore[i] = usage
	}

	return stats, nil
}

// GetMemory retorna estatísticas de memória
func (m *SystemMonitor) GetMemory() (*MemoryStats, error) {
	if runtime.GOOS == "darwin" {
		return getMemoryMac()
	}
	return getMemoryLinux()
}

// GetDisk retorna estatísticas de disco
func (m *SystemMonitor) GetDisk() ([]*DiskStats, error) {
	return getDiskStats()
}

// GetNetwork retorna estatísticas de rede
func (m *SystemMonitor) GetNetwork() ([]*NetworkStats, error) {
	if runtime.GOOS == "darwin" {
		return getNetworkMac()
	}
	return getNetworkLinux()
}

// GetProcesses retorna os processos que mais consomem recursos
func (m *SystemMonitor) GetProcesses(limit int) ([]*ProcessStats, error) {
	return getTopProcesses(limit)
}

// Helper functions

func getLoadAverage() ([3]float64, error) {
	var result [3]float64

	if runtime.GOOS == "linux" {
		// Read from /proc/loadavg
		data, err := os.ReadFile("/proc/loadavg")
		if err != nil {
			return result, err
		}

		parts := strings.Fields(string(data))
		if len(parts) < 3 {
			return result, fmt.Errorf("invalid /proc/loadavg format")
		}

		for i := 0; i < 3; i++ {
			val, err := strconv.ParseFloat(parts[i], 64)
			if err != nil {
				return result, err
			}
			result[i] = val
		}
		return result, nil
	}

	if runtime.GOOS == "darwin" {
		// Use sysctl on macOS
		cmd := exec.Command("sysctl", "-n", "vm.loadavg")
		output, err := cmd.Output()
		if err != nil {
			return result, err
		}

		// Output format: { 1.23 2.34 3.45 }
		str := strings.TrimSpace(string(output))
		str = strings.Trim(str, "{}")
		parts := strings.Fields(str)

		if len(parts) < 3 {
			return result, fmt.Errorf("invalid loadavg format")
		}

		for i := 0; i < 3; i++ {
			val, err := strconv.ParseFloat(parts[i], 64)
			if err != nil {
				return result, err
			}
			result[i] = val
		}
		return result, nil
	}

	return result, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
}

func getCPUUsage() (float64, error) {
	if runtime.GOOS == "linux" {
		return getCPUUsageLinux()
	}
	if runtime.GOOS == "darwin" {
		return getCPUUsageMac()
	}
	return 0, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
}

func getCPUUsageLinux() (float64, error) {
	// Read /proc/stat twice with a small delay
	stat1, err := readProcStat()
	if err != nil {
		return 0, err
	}

	time.Sleep(100 * time.Millisecond)

	stat2, err := readProcStat()
	if err != nil {
		return 0, err
	}

	// Calculate CPU usage
	totalDelta := stat2.total - stat1.total
	idleDelta := stat2.idle - stat1.idle

	if totalDelta == 0 {
		return 0, nil
	}

	usage := 100.0 * float64(totalDelta-idleDelta) / float64(totalDelta)
	return usage, nil
}

type procStat struct {
	total uint64
	idle  uint64
}

func readProcStat() (*procStat, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read /proc/stat")
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return nil, fmt.Errorf("invalid /proc/stat format")
	}

	// Parse CPU times
	var times []uint64
	for i := 1; i < len(fields); i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return nil, err
		}
		times = append(times, val)
	}

	// Calculate total and idle
	var total uint64
	for _, t := range times {
		total += t
	}

	idle := times[3] // idle is the 4th field

	return &procStat{total: total, idle: idle}, nil
}

func getCPUUsageMac() (float64, error) {
	// Use top command on macOS
	cmd := exec.Command("top", "-l", "1", "-n", "0")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Parse top output for CPU usage
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "CPU usage") {
			// Format: "CPU usage: 12.34% user, 5.67% sys, 82.00% idle"
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				continue
			}

			values := strings.Split(parts[1], ",")
			if len(values) < 3 {
				continue
			}

			// Get idle percentage
			idleStr := strings.TrimSpace(values[2])
			idleStr = strings.TrimSuffix(idleStr, "% idle")
			idle, err := strconv.ParseFloat(strings.TrimSpace(idleStr), 64)
			if err != nil {
				continue
			}

			usage := 100.0 - idle
			return usage, nil
		}
	}

	return 0, fmt.Errorf("failed to parse CPU usage from top")
}

func getMemoryLinux() (*MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &MemoryStats{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		val *= 1024 // Convert from KB to bytes

		switch key {
		case "MemTotal":
			stats.Total = val
		case "MemFree":
			stats.Free = val
		case "MemAvailable":
			stats.Available = val
		case "SwapTotal":
			stats.SwapTotal = val
		case "SwapFree":
			stats.SwapFree = val
		}
	}

	stats.Used = stats.Total - stats.Available
	stats.SwapUsed = stats.SwapTotal - stats.SwapFree

	if stats.Total > 0 {
		stats.UsagePercent = 100.0 * float64(stats.Used) / float64(stats.Total)
	}

	return stats, nil
}

func getMemoryMac() (*MemoryStats, error) {
	stats := &MemoryStats{}

	// Get total memory
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	total, err := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return nil, err
	}
	stats.Total = total

	// Get vm_stat for memory usage
	cmd = exec.Command("vm_stat")
	output, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse vm_stat output
	lines := strings.Split(string(output), "\n")
	var pageSize uint64 = 4096 // Default page size
	var freePages, activePages, inactivePages, wiredPages uint64

	for _, line := range lines {
		if strings.Contains(line, "page size of") {
			// Extract page size from format: "Mach Virtual Memory Statistics: (page size of 16384 bytes)"
			// Look for the number between "of" and "bytes"
			startIdx := strings.Index(line, "page size of ")
			if startIdx != -1 {
				startIdx += len("page size of ")
				endIdx := strings.Index(line[startIdx:], " bytes")
				if endIdx != -1 {
					sizeStr := line[startIdx : startIdx+endIdx]
					val, err := strconv.ParseUint(strings.TrimSpace(sizeStr), 10, 64)
					if err == nil {
						pageSize = val
					}
				}
			}
			continue
		}

		// Split on colon to separate key from value
		// Format: "Pages wired down:                        146831."
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIdx])
		valStr := strings.TrimSpace(line[colonIdx+1:])
		valStr = strings.TrimSuffix(valStr, ".")
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			continue
		}

		switch key {
		case "Pages free":
			freePages = val
		case "Pages active":
			activePages = val
		case "Pages inactive":
			inactivePages = val
		case "Pages wired down":
			wiredPages = val
		}
	}

	stats.Free = freePages * pageSize
	stats.Used = (activePages + wiredPages) * pageSize
	stats.Available = (freePages + inactivePages) * pageSize

	if stats.Total > 0 {
		stats.UsagePercent = 100.0 * float64(stats.Used) / float64(stats.Total)
	}

	return stats, nil
}

func getDiskStats() ([]*DiskStats, error) {
	cmd := exec.Command("df", "-k")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var stats []*DiskStats
	lines := strings.Split(string(output), "\n")

	for i, line := range lines {
		if i == 0 {
			// Skip header
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// Skip special filesystems
		if strings.HasPrefix(fields[0], "map") || strings.HasPrefix(fields[0], "devfs") {
			continue
		}

		total, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		used, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			continue
		}
		avail, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			continue
		}

		// Convert from 1K blocks to bytes
		total *= 1024
		used *= 1024
		avail *= 1024

		usagePercent := 0.0
		if total > 0 {
			usagePercent = 100.0 * float64(used) / float64(total)
		}

		stat := &DiskStats{
			Filesystem:   fields[0],
			MountPoint:   fields[len(fields)-1],
			Total:        total,
			Used:         used,
			Available:    avail,
			UsagePercent: usagePercent,
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func getNetworkLinux() ([]*NetworkStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stats []*NetworkStats
	scanner := bufio.NewScanner(file)

	// Skip first two lines (headers)
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		iface := strings.TrimSpace(parts[0])
		fields := strings.Fields(parts[1])
		if len(fields) < 16 {
			continue
		}

		// Skip loopback
		if iface == "lo" {
			continue
		}

		bytesRecv, _ := strconv.ParseUint(fields[0], 10, 64)
		packetsRecv, _ := strconv.ParseUint(fields[1], 10, 64)
		errorsRecv, _ := strconv.ParseUint(fields[2], 10, 64)
		bytesSent, _ := strconv.ParseUint(fields[8], 10, 64)
		packetsSent, _ := strconv.ParseUint(fields[9], 10, 64)
		errorsSent, _ := strconv.ParseUint(fields[10], 10, 64)

		stat := &NetworkStats{
			Interface:   iface,
			BytesRecv:   bytesRecv,
			BytesSent:   bytesSent,
			PacketsRecv: packetsRecv,
			PacketsSent: packetsSent,
			ErrorsRecv:  errorsRecv,
			ErrorsSent:  errorsSent,
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func getNetworkMac() ([]*NetworkStats, error) {
	cmd := exec.Command("netstat", "-ib")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var stats []*NetworkStats
	lines := strings.Split(string(output), "\n")

	for i, line := range lines {
		if i == 0 {
			// Skip header
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		iface := fields[0]

		// Skip loopback and inactive interfaces
		if iface == "lo0" || strings.HasSuffix(iface, "*") {
			continue
		}

		// Only process lines with Link# (which is in fields[2])
		if !strings.HasPrefix(fields[2], "<Link#") {
			continue
		}

		bytesRecv, _ := strconv.ParseUint(fields[6], 10, 64)
		bytesSent, _ := strconv.ParseUint(fields[9], 10, 64)

		stat := &NetworkStats{
			Interface: iface,
			BytesRecv: bytesRecv,
			BytesSent: bytesSent,
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func getTopProcesses(limit int) ([]*ProcessStats, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "darwin" {
		cmd = exec.Command("ps", "-arcwwwxo", "pid,comm,%cpu,%mem", "-m")
	} else {
		cmd = exec.Command("ps", "aux", "--sort=-%cpu")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var stats []*ProcessStats
	lines := strings.Split(string(output), "\n")

	for i, line := range lines {
		if i == 0 || len(stats) >= limit {
			// Skip header or reached limit
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		var pid int
		var name string
		var cpuPercent, memPercent float64

		if runtime.GOOS == "darwin" {
			pid, _ = strconv.Atoi(fields[0])
			name = fields[1]
			cpuPercent, _ = strconv.ParseFloat(fields[2], 64)
			memPercent, _ = strconv.ParseFloat(fields[3], 64)
		} else {
			pid, _ = strconv.Atoi(fields[1])
			cpuPercent, _ = strconv.ParseFloat(fields[2], 64)
			memPercent, _ = strconv.ParseFloat(fields[3], 64)
			name = fields[10]
		}

		stat := &ProcessStats{
			PID:           pid,
			Name:          name,
			CPUPercent:    cpuPercent,
			MemoryPercent: memPercent,
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

// FormatBytes converte bytes para formato legível
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
