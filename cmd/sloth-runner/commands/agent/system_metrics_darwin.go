//go:build darwin
// +build darwin
package agent

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// MemoryInfo holds memory statistics
type MemoryInfo struct {
	Total uint64
	Used  uint64
	Free  uint64
}

// DiskInfo holds disk usage statistics
type DiskInfo struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

// ProcessInfo holds process information
type ProcessInfo struct {
	PID           int
	Name          string
	Status        string
	CPUPercent    float64
	MemoryPercent float64
	MemoryBytes   uint64
	User          string
	Command       string
	StartedAt     int64
}

// NetworkInterfaceInfo holds network interface information
type NetworkInterfaceInfo struct {
	Name        string
	IPAddresses []string
	MACAddress  string
	BytesSent   uint64
	BytesRecv   uint64
	IsUp        bool
}

// DiskPartitionInfo holds disk partition information
type DiskPartitionInfo struct {
	Device       string
	Mountpoint   string
	FSType       string
	TotalBytes   uint64
	UsedBytes    uint64
	FreeBytes    uint64
	Percent      float64
	IOReadBytes  uint64
	IOWriteBytes uint64
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp int64
	Level     string
	Message   string
}

// getCPUPercent returns current CPU usage percentage
func getCPUPercent() float64 {
	// Read /proc/stat for CPU times
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return 0.0
	}

	// First line is overall CPU
	fields := strings.Fields(lines[0])
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0.0
	}

	// Calculate total and idle
	var total, idle uint64
	for i := 1; i < len(fields); i++ {
		val, _ := strconv.ParseUint(fields[i], 10, 64)
		total += val
		if i == 4 { // idle is the 4th field
			idle = val
		}
	}

	// Simple approximation: (total - idle) / total * 100
	if total == 0 {
		return 0.0
	}

	used := total - idle
	return float64(used) / float64(total) * 100.0
}

// getMemoryInfo returns system memory information
func getMemoryInfo() (MemoryInfo, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, err
	}

	var memTotal, memAvailable uint64
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024 // Convert from KB to bytes

		if strings.HasPrefix(fields[0], "MemTotal:") {
			memTotal = value
		} else if strings.HasPrefix(fields[0], "MemAvailable:") {
			memAvailable = value
		}
	}

	memUsed := memTotal - memAvailable

	return MemoryInfo{
		Total: memTotal,
		Used:  memUsed,
		Free:  memAvailable,
	}, nil
}

// getDiskUsage returns disk usage for a given path
func getDiskUsage(path string) (DiskInfo, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return DiskInfo{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	var usedPercent float64
	if total > 0 {
		usedPercent = float64(used) / float64(total) * 100.0
	}

	return DiskInfo{
		Total:       total,
		Used:        used,
		Free:        free,
		UsedPercent: usedPercent,
	}, nil
}

// getLoadAverage returns system load average
func getLoadAverage() [3]float64 {
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return [3]float64{0, 0, 0}
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return [3]float64{0, 0, 0}
	}

	load1, _ := strconv.ParseFloat(fields[0], 64)
	load5, _ := strconv.ParseFloat(fields[1], 64)
	load15, _ := strconv.ParseFloat(fields[2], 64)

	return [3]float64{load1, load5, load15}
}

// getProcessCount returns the number of running processes
func getProcessCount() int {
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return 0
	}

	count := 0
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		// Check if directory name is numeric (PID)
		if _, err := strconv.Atoi(dir.Name()); err == nil {
			count++
		}
	}

	return count
}

// getSystemUptime returns system uptime in seconds
func getSystemUptime() uint64 {
	data, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}

	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0
	}

	uptime, _ := strconv.ParseFloat(fields[0], 64)
	return uint64(uptime)
}

// getProcesses returns list of running processes
func getProcesses() ([]ProcessInfo, error) {
	// Use ps command (optimized version is in process_linux.go for Linux)
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	processes := make([]ProcessInfo, 0)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 11 {
			continue
		}

		pid, _ := strconv.Atoi(fields[1])
		cpuPercent, _ := strconv.ParseFloat(fields[2], 64)
		memPercent, _ := strconv.ParseFloat(fields[3], 64)

		// Command is everything from field 10 onwards
		command := strings.Join(fields[10:], " ")

		processes = append(processes, ProcessInfo{
			PID:           pid,
			Name:          fields[10],
			User:          fields[0],
			CPUPercent:    cpuPercent,
			MemoryPercent: memPercent,
			Status:        fields[7],
			Command:       command,
			StartedAt:     time.Now().Unix(), // Approximate
		})
	}

	return processes, nil
}

// getNetworkInterfaces returns list of network interfaces
func getNetworkInterfaces() ([]NetworkInterfaceInfo, error) {
	// Use ip command to get interface info
	cmd := exec.Command("ip", "-j", "addr", "show")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to reading /proc/net/dev
		return getNetworkInterfacesFallback()
	}

	// Parse JSON output properly
	var ipOutput []struct {
		Ifname   string `json:"ifname"`
		Operstate string `json:"operstate"`
		Address  string `json:"address"`
		AddrInfo []struct {
			Family string `json:"family"`
			Local  string `json:"local"`
		} `json:"addr_info"`
	}

	if err := json.Unmarshal(output, &ipOutput); err != nil {
		// Fallback if JSON parsing fails
		return getNetworkInterfacesFallback()
	}

	// Get bytes sent/recv from /proc/net/dev
	netStats := make(map[string]struct{ sent, recv uint64 })
	if data, err := ioutil.ReadFile("/proc/net/dev"); err == nil {
		lines := strings.Split(string(data), "\n")
		for i := 2; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			name := strings.TrimSpace(parts[0])
			fields := strings.Fields(parts[1])
			if len(fields) >= 9 {
				recv, _ := strconv.ParseUint(fields[0], 10, 64)
				sent, _ := strconv.ParseUint(fields[8], 10, 64)
				netStats[name] = struct{ sent, recv uint64 }{sent, recv}
			}
		}
	}

	interfaces := make([]NetworkInterfaceInfo, 0)
	for _, iface := range ipOutput {
		// Get IPs
		ips := make([]string, 0)
		for _, addr := range iface.AddrInfo {
			if addr.Family == "inet" || addr.Family == "inet6" {
				ips = append(ips, addr.Local)
			}
		}

		// Get stats
		stats := netStats[iface.Ifname]

		interfaces = append(interfaces, NetworkInterfaceInfo{
			Name:        iface.Ifname,
			IPAddresses: ips,
			MACAddress:  iface.Address,
			BytesSent:   stats.sent,
			BytesRecv:   stats.recv,
			IsUp:        iface.Operstate == "UP" || iface.Operstate == "UNKNOWN",
		})
	}

	return interfaces, nil
}

// getNetworkInterfacesFallback reads network info from /proc
func getNetworkInterfacesFallback() ([]NetworkInterfaceInfo, error) {
	data, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	interfaces := make([]NetworkInterfaceInfo, 0)
	lines := strings.Split(string(data), "\n")

	// Skip first two header lines
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		fields := strings.Fields(parts[1])

		if len(fields) < 8 {
			continue
		}

		bytesRecv, _ := strconv.ParseUint(fields[0], 10, 64)
		bytesSent, _ := strconv.ParseUint(fields[8], 10, 64)

		interfaces = append(interfaces, NetworkInterfaceInfo{
			Name:      name,
			BytesRecv: bytesRecv,
			BytesSent: bytesSent,
			IsUp:      true,
		})
	}

	return interfaces, nil
}

// getDiskPartitions returns list of disk partitions
func getDiskPartitions() ([]DiskPartitionInfo, error) {
	// Read /proc/mounts
	data, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return nil, err
	}

	partitions := make([]DiskPartitionInfo, 0)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		device := fields[0]
		mountpoint := fields[1]
		fstype := fields[2]

		// Skip virtual filesystems
		if strings.HasPrefix(device, "/dev/") {
			var stat syscall.Statfs_t
			if err := syscall.Statfs(mountpoint, &stat); err != nil {
				continue
			}

			total := stat.Blocks * uint64(stat.Bsize)
			free := stat.Bavail * uint64(stat.Bsize)
			used := total - free

			var percent float64
			if total > 0 {
				percent = float64(used) / float64(total) * 100.0
			}

			partitions = append(partitions, DiskPartitionInfo{
				Device:     device,
				Mountpoint: mountpoint,
				FSType:     fstype,
				TotalBytes: total,
				UsedBytes:  used,
				FreeBytes:  free,
				Percent:    percent,
			})
		}
	}

	return partitions, nil
}

// collectLogs collects system logs and sends to channel
func collectLogs(logChan chan LogEntry) {
	// Read from journalctl or syslog
	cmd := exec.Command("journalctl", "-f", "-n", "10", "-o", "short")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Failed to open journalctl pipe", "error", err)
		return
	}

	if err := cmd.Start(); err != nil {
		slog.Error("Failed to start journalctl", "error", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		logChan <- LogEntry{
			Timestamp: time.Now().Unix(),
			Level:     "INFO",
			Message:   line,
		}
	}
}

// getNetworkBytes returns total network RX and TX bytes across all interfaces
func getNetworkBytes() (uint64, uint64) {
	interfaces, err := getNetworkInterfaces()
	if err != nil {
		return 0, 0
	}

	var totalRx, totalTx uint64
	for _, iface := range interfaces {
		totalRx += iface.BytesRecv
		totalTx += iface.BytesSent
	}

	return totalRx, totalTx
}
