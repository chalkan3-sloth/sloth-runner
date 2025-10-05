package telemetry

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemMetrics holds comprehensive system metrics
type SystemMetrics struct {
	// CPU Metrics
	CPUUsagePerCore []float64
	CPUUsageTotal   float64
	CPUCores        int
	CPUThreads      int
	CPUSpeed        float64
	CPUModel        string
	LoadAverage1    float64
	LoadAverage5    float64
	LoadAverage15   float64

	// Memory Metrics
	MemoryTotal     uint64
	MemoryUsed      uint64
	MemoryFree      uint64
	MemoryAvailable uint64
	MemoryUsedPct   float64
	SwapTotal       uint64
	SwapUsed        uint64
	SwapFree        uint64
	SwapUsedPct     float64

	// Disk Metrics
	DiskMetrics []DiskInfo

	// Network Metrics
	NetworkInterfaces []NetworkInfo
	NetworkConnections int

	// Process Metrics
	ProcessCount   int
	TopProcesses   []ProcessInfo
	ZombieCount    int

	// System Info
	Hostname       string
	OS             string
	Platform       string
	PlatformFamily string
	PlatformVersion string
	KernelVersion  string
	KernelArch     string
	Uptime         uint64
	BootTime       uint64
	Temperature    []TempInfo

	// Timestamp
	Timestamp      time.Time
}

// DiskInfo holds information about a disk/partition
type DiskInfo struct {
	Device      string
	MountPoint  string
	Fstype      string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
	InodesTotal uint64
	InodesUsed  uint64
	InodesPct   float64
}

// NetworkInfo holds network interface information
type NetworkInfo struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	Errin       uint64
	Errout      uint64
	Dropin      uint64
	Dropout     uint64
	Speed       uint64
}

// ProcessInfo holds process information
type ProcessInfo struct {
	PID         int32
	Name        string
	Username    string
	CPUPercent  float64
	MemPercent  float32
	MemoryMB    float64
	Status      string
	CreateTime  int64
}

// TempInfo holds temperature sensor information
type TempInfo struct {
	SensorKey   string
	Temperature float64
	High        float64
	Critical    float64
}

// CollectSystemMetrics gathers comprehensive system metrics
func CollectSystemMetrics() (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
	}

	// Collect CPU metrics
	if err := collectCPUMetrics(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect CPU metrics: %w", err)
	}

	// Collect Memory metrics
	if err := collectMemoryMetrics(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect memory metrics: %w", err)
	}

	// Collect Disk metrics
	if err := collectDiskMetrics(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect disk metrics: %w", err)
	}

	// Collect Network metrics
	if err := collectNetworkMetrics(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect network metrics: %w", err)
	}

	// Collect Process metrics
	if err := collectProcessMetrics(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect process metrics: %w", err)
	}

	// Collect System info
	if err := collectSystemInfo(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect system info: %w", err)
	}

	return metrics, nil
}

func collectCPUMetrics(m *SystemMetrics) error {
	// CPU usage per core
	perCoreUsage, err := cpu.Percent(time.Second, true)
	if err == nil {
		m.CPUUsagePerCore = perCoreUsage
	}

	// Total CPU usage
	totalUsage, err := cpu.Percent(time.Second, false)
	if err == nil && len(totalUsage) > 0 {
		m.CPUUsageTotal = totalUsage[0]
	}

	// CPU info
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		m.CPUModel = cpuInfo[0].ModelName
		m.CPUCores = int(cpuInfo[0].Cores)
		m.CPUSpeed = cpuInfo[0].Mhz
	}

	// Logical CPU count
	m.CPUThreads = runtime.NumCPU()

	// Load averages
	loadAvg, err := load.Avg()
	if err == nil {
		m.LoadAverage1 = loadAvg.Load1
		m.LoadAverage5 = loadAvg.Load5
		m.LoadAverage15 = loadAvg.Load15
	}

	return nil
}

func collectMemoryMetrics(m *SystemMetrics) error {
	// Virtual memory
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	m.MemoryTotal = vmem.Total
	m.MemoryUsed = vmem.Used
	m.MemoryFree = vmem.Free
	m.MemoryAvailable = vmem.Available
	m.MemoryUsedPct = vmem.UsedPercent

	// Swap memory
	swap, err := mem.SwapMemory()
	if err == nil {
		m.SwapTotal = swap.Total
		m.SwapUsed = swap.Used
		m.SwapFree = swap.Free
		m.SwapUsedPct = swap.UsedPercent
	}

	return nil
}

func collectDiskMetrics(m *SystemMetrics) error {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	var diskInfos []DiskInfo
	for _, partition := range partitions {
		// Skip special filesystems
		if strings.HasPrefix(partition.Fstype, "fuse") ||
			strings.HasPrefix(partition.Device, "shm") ||
			strings.HasPrefix(partition.Mountpoint, "/snap") ||
			strings.HasPrefix(partition.Mountpoint, "/run") ||
			strings.HasPrefix(partition.Mountpoint, "/sys") ||
			strings.HasPrefix(partition.Mountpoint, "/proc") {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		diskInfo := DiskInfo{
			Device:      partition.Device,
			MountPoint:  partition.Mountpoint,
			Fstype:      partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
			InodesTotal: usage.InodesTotal,
			InodesUsed:  usage.InodesUsed,
		}

		if usage.InodesTotal > 0 {
			diskInfo.InodesPct = float64(usage.InodesUsed) / float64(usage.InodesTotal) * 100
		}

		diskInfos = append(diskInfos, diskInfo)
	}

	m.DiskMetrics = diskInfos
	return nil
}

func collectNetworkMetrics(m *SystemMetrics) error {
	// Network interfaces
	netIO, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	var networkInfos []NetworkInfo
	for _, io := range netIO {
		// Skip loopback
		if io.Name == "lo" || io.Name == "lo0" {
			continue
		}

		info := NetworkInfo{
			Name:        io.Name,
			BytesSent:   io.BytesSent,
			BytesRecv:   io.BytesRecv,
			PacketsSent: io.PacketsSent,
			PacketsRecv: io.PacketsRecv,
			Errin:       io.Errin,
			Errout:      io.Errout,
			Dropin:      io.Dropin,
			Dropout:     io.Dropout,
		}

		networkInfos = append(networkInfos, info)
	}

	m.NetworkInterfaces = networkInfos

	// Connection count
	connections, err := net.Connections("all")
	if err == nil {
		m.NetworkConnections = len(connections)
	}

	return nil
}

func collectProcessMetrics(m *SystemMetrics) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	m.ProcessCount = len(processes)

	var processInfos []ProcessInfo
	var zombieCount int

	for _, p := range processes {
		// Get process info
		name, _ := p.Name()
		username, _ := p.Username()
		cpuPct, _ := p.CPUPercent()
		memPct, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		statusList, _ := p.Status()
		createTime, _ := p.CreateTime()

		// Get first status from list
		var status string
		if len(statusList) > 0 {
			status = statusList[0]
		}

		if status == "Z" {
			zombieCount++
		}

		var memMB float64
		if memInfo != nil {
			memMB = float64(memInfo.RSS) / 1024 / 1024
		}

		info := ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			Username:   username,
			CPUPercent: cpuPct,
			MemPercent: memPct,
			MemoryMB:   memMB,
			Status:     status,
			CreateTime: createTime,
		}

		processInfos = append(processInfos, info)
	}

	// Sort by CPU usage and get top 10
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPUPercent > processInfos[j].CPUPercent
	})

	topCount := 10
	if len(processInfos) < topCount {
		topCount = len(processInfos)
	}

	m.TopProcesses = processInfos[:topCount]
	m.ZombieCount = zombieCount

	return nil
}

func collectSystemInfo(m *SystemMetrics) error {
	// Host info
	hostInfo, err := host.Info()
	if err != nil {
		return err
	}

	m.Hostname = hostInfo.Hostname
	m.OS = hostInfo.OS
	m.Platform = hostInfo.Platform
	m.PlatformFamily = hostInfo.PlatformFamily
	m.PlatformVersion = hostInfo.PlatformVersion
	m.KernelVersion = hostInfo.KernelVersion
	m.KernelArch = hostInfo.KernelArch
	m.Uptime = hostInfo.Uptime
	m.BootTime = hostInfo.BootTime

	// Temperature sensors (if available)
	temps, err := host.SensorsTemperatures()
	if err == nil {
		var tempInfos []TempInfo
		for _, temp := range temps {
			info := TempInfo{
				SensorKey:   temp.SensorKey,
				Temperature: temp.Temperature,
				High:        temp.High,
				Critical:    temp.Critical,
			}
			tempInfos = append(tempInfos, info)
		}
		m.Temperature = tempInfos
	}

	return nil
}

// FormatBytes converts bytes to human-readable format
func FormatBytes(bytes uint64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// FormatUptime converts uptime seconds to human-readable format
func FormatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}