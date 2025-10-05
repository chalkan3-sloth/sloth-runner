package luainterface

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	lua "github.com/yuin/gopher-lua"
)

// SystemModule provides enhanced system information and monitoring
type SystemModule struct{}

// NewSystemModule creates a new system module
func NewSystemModule() *SystemModule {
	return &SystemModule{}
}

// RegisterSystemModule registers the system module with the Lua state
func RegisterSystemModule(L *lua.LState) {
	module := NewSystemModule()
	
	// Create the system table
	systemTable := L.NewTable()
	
	// CPU information
	L.SetField(systemTable, "cpu_info", L.NewFunction(module.luaCPUInfo))
	L.SetField(systemTable, "cpu_usage", L.NewFunction(module.luaCPUUsage))
	L.SetField(systemTable, "cpu_count", L.NewFunction(module.luaCPUCount))
	L.SetField(systemTable, "load_average", L.NewFunction(module.luaLoadAverage))
	
	// Memory information
	L.SetField(systemTable, "memory_info", L.NewFunction(module.luaMemoryInfo))
	L.SetField(systemTable, "memory_usage", L.NewFunction(module.luaMemoryUsage))
	L.SetField(systemTable, "swap_info", L.NewFunction(module.luaSwapInfo))
	
	// Disk information
	L.SetField(systemTable, "disk_usage", L.NewFunction(module.luaDiskUsage))
	L.SetField(systemTable, "disk_io", L.NewFunction(module.luaDiskIO))
	L.SetField(systemTable, "disk_partitions", L.NewFunction(module.luaDiskPartitions))
	
	// Network information
	L.SetField(systemTable, "network_interfaces", L.NewFunction(module.luaNetworkInterfaces))
	L.SetField(systemTable, "network_stats", L.NewFunction(module.luaNetworkStats))
	L.SetField(systemTable, "network_connections", L.NewFunction(module.luaNetworkConnections))
	
	// Process information
	L.SetField(systemTable, "processes", L.NewFunction(module.luaProcesses))
	L.SetField(systemTable, "process_info", L.NewFunction(module.luaProcessInfo))
	L.SetField(systemTable, "kill_process", L.NewFunction(module.luaKillProcess))
	
	// System information
	L.SetField(systemTable, "host_info", L.NewFunction(module.luaHostInfo))
	L.SetField(systemTable, "uptime", L.NewFunction(module.luaUptime))
	L.SetField(systemTable, "environment", L.NewFunction(module.luaEnvironment))
	L.SetField(systemTable, "users", L.NewFunction(module.luaUsers))
	
	// Performance monitoring
	L.SetField(systemTable, "performance_snapshot", L.NewFunction(module.luaPerformanceSnapshot))
	L.SetField(systemTable, "system_health", L.NewFunction(module.luaSystemHealth))
	
	// Register the system table globally
	L.SetGlobal("system", systemTable)
}

// CPU information
func (s *SystemModule) luaCPUInfo(L *lua.LState) int {
	cpuInfo, err := cpu.Info()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for i, info := range cpuInfo {
		cpuTable := L.NewTable()
		L.SetField(cpuTable, "model", lua.LString(info.ModelName))
		L.SetField(cpuTable, "family", lua.LString(info.Family))
		L.SetField(cpuTable, "speed_mhz", lua.LNumber(info.Mhz))
		L.SetField(cpuTable, "cores", lua.LNumber(info.Cores))
		L.SetField(cpuTable, "cache_size", lua.LNumber(info.CacheSize))
		L.SetField(cpuTable, "vendor_id", lua.LString(info.VendorID))
		
		result.RawSetInt(i+1, cpuTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaCPUUsage(L *lua.LState) int {
	interval := time.Duration(L.OptNumber(1, 1)) * time.Second
	
	usage, err := cpu.Percent(interval, false)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	if len(usage) > 0 {
		L.Push(lua.LNumber(usage[0]))
	} else {
		L.Push(lua.LNumber(0))
	}
	return 1
}

func (s *SystemModule) luaCPUCount(L *lua.LState) int {
	logical := L.OptBool(1, true)
	
	var count int
	var err error
	
	if logical {
		count, err = cpu.Counts(true)
	} else {
		count, err = cpu.Counts(false)
	}
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LNumber(count))
	return 1
}

func (s *SystemModule) luaLoadAverage(L *lua.LState) int {
	loadAvg, err := load.Avg()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	L.SetField(result, "load1", lua.LNumber(loadAvg.Load1))
	L.SetField(result, "load5", lua.LNumber(loadAvg.Load5))
	L.SetField(result, "load15", lua.LNumber(loadAvg.Load15))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Memory information
func (s *SystemModule) luaMemoryInfo(L *lua.LState) int {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	result := L.NewTable()
	L.SetField(result, "total", lua.LNumber(memInfo.Total))
	L.SetField(result, "available", lua.LNumber(memInfo.Available))
	L.SetField(result, "used", lua.LNumber(memInfo.Used))
	L.SetField(result, "free", lua.LNumber(memInfo.Free))
	L.SetField(result, "percent", lua.LNumber(memInfo.UsedPercent))
	L.SetField(result, "buffers", lua.LNumber(memInfo.Buffers))
	L.SetField(result, "cached", lua.LNumber(memInfo.Cached))

	L.Push(result)
	return 1
}

func (s *SystemModule) luaMemoryUsage(L *lua.LState) int {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LNumber(memInfo.UsedPercent))
	return 1
}

func (s *SystemModule) luaSwapInfo(L *lua.LState) int {
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	L.SetField(result, "total", lua.LNumber(swapInfo.Total))
	L.SetField(result, "used", lua.LNumber(swapInfo.Used))
	L.SetField(result, "free", lua.LNumber(swapInfo.Free))
	L.SetField(result, "percent", lua.LNumber(swapInfo.UsedPercent))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Disk information
func (s *SystemModule) luaDiskUsage(L *lua.LState) int {
	path := L.CheckString(1)
	
	diskInfo, err := disk.Usage(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	L.SetField(result, "total", lua.LNumber(diskInfo.Total))
	L.SetField(result, "used", lua.LNumber(diskInfo.Used))
	L.SetField(result, "free", lua.LNumber(diskInfo.Free))
	L.SetField(result, "percent", lua.LNumber(diskInfo.UsedPercent))
	L.SetField(result, "inodes_total", lua.LNumber(diskInfo.InodesTotal))
	L.SetField(result, "inodes_used", lua.LNumber(diskInfo.InodesUsed))
	L.SetField(result, "inodes_free", lua.LNumber(diskInfo.InodesFree))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaDiskIO(L *lua.LState) int {
	diskStats, err := disk.IOCounters()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for device, stats := range diskStats {
		deviceTable := L.NewTable()
		L.SetField(deviceTable, "read_count", lua.LNumber(stats.ReadCount))
		L.SetField(deviceTable, "write_count", lua.LNumber(stats.WriteCount))
		L.SetField(deviceTable, "read_bytes", lua.LNumber(stats.ReadBytes))
		L.SetField(deviceTable, "write_bytes", lua.LNumber(stats.WriteBytes))
		L.SetField(deviceTable, "read_time", lua.LNumber(stats.ReadTime))
		L.SetField(deviceTable, "write_time", lua.LNumber(stats.WriteTime))
		
		result.RawSetString(device, deviceTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaDiskPartitions(L *lua.LState) int {
	partitions, err := disk.Partitions(false)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for i, partition := range partitions {
		partTable := L.NewTable()
		L.SetField(partTable, "device", lua.LString(partition.Device))
		L.SetField(partTable, "mountpoint", lua.LString(partition.Mountpoint))
		L.SetField(partTable, "fstype", lua.LString(partition.Fstype))
		L.SetField(partTable, "opts", lua.LString(strings.Join(partition.Opts, ",")))
		
		result.RawSetInt(i+1, partTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Network information
func (s *SystemModule) luaNetworkInterfaces(L *lua.LState) int {
	interfaces, err := net.Interfaces()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for i, iface := range interfaces {
		ifaceTable := L.NewTable()
		L.SetField(ifaceTable, "name", lua.LString(iface.Name))
		L.SetField(ifaceTable, "mtu", lua.LNumber(iface.MTU))
		L.SetField(ifaceTable, "hardware_addr", lua.LString(iface.HardwareAddr))
		L.SetField(ifaceTable, "flags", lua.LString(strings.Join(iface.Flags, ",")))
		
		// Add addresses
		addrs := L.NewTable()
		for j, addr := range iface.Addrs {
			addrTable := L.NewTable()
			L.SetField(addrTable, "addr", lua.LString(addr.Addr))
			addrs.RawSetInt(j+1, addrTable)
		}
		L.SetField(ifaceTable, "addresses", addrs)
		
		result.RawSetInt(i+1, ifaceTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaNetworkStats(L *lua.LState) int {
	stats, err := net.IOCounters(false)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for i, stat := range stats {
		statTable := L.NewTable()
		L.SetField(statTable, "name", lua.LString(stat.Name))
		L.SetField(statTable, "bytes_sent", lua.LNumber(stat.BytesSent))
		L.SetField(statTable, "bytes_recv", lua.LNumber(stat.BytesRecv))
		L.SetField(statTable, "packets_sent", lua.LNumber(stat.PacketsSent))
		L.SetField(statTable, "packets_recv", lua.LNumber(stat.PacketsRecv))
		L.SetField(statTable, "err_in", lua.LNumber(stat.Errin))
		L.SetField(statTable, "err_out", lua.LNumber(stat.Errout))
		L.SetField(statTable, "drop_in", lua.LNumber(stat.Dropin))
		L.SetField(statTable, "drop_out", lua.LNumber(stat.Dropout))
		
		result.RawSetInt(i+1, statTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaNetworkConnections(L *lua.LState) int {
	kind := L.OptString(1, "inet")
	
	connections, err := net.Connections(kind)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for i, conn := range connections {
		connTable := L.NewTable()
		L.SetField(connTable, "fd", lua.LNumber(conn.Fd))
		L.SetField(connTable, "family", lua.LNumber(conn.Family))
		L.SetField(connTable, "type", lua.LNumber(conn.Type))
		L.SetField(connTable, "local_addr", lua.LString(conn.Laddr.IP))
		L.SetField(connTable, "local_port", lua.LNumber(conn.Laddr.Port))
		L.SetField(connTable, "remote_addr", lua.LString(conn.Raddr.IP))
		L.SetField(connTable, "remote_port", lua.LNumber(conn.Raddr.Port))
		L.SetField(connTable, "status", lua.LString(conn.Status))
		L.SetField(connTable, "pid", lua.LNumber(conn.Pid))
		
		result.RawSetInt(i+1, connTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Process information
func (s *SystemModule) luaProcesses(L *lua.LState) int {
	pids, err := process.Pids()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for i, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue
		}
		
		procTable := L.NewTable()
		L.SetField(procTable, "pid", lua.LNumber(pid))
		
		if name, err := proc.Name(); err == nil {
			L.SetField(procTable, "name", lua.LString(name))
		}
		
		if status, err := proc.Status(); err == nil {
			L.SetField(procTable, "status", lua.LString(strings.Join(status, ",")))
		}
		
		if cpuPercent, err := proc.CPUPercent(); err == nil {
			L.SetField(procTable, "cpu_percent", lua.LNumber(cpuPercent))
		}
		
		if memInfo, err := proc.MemoryInfo(); err == nil {
			memTable := L.NewTable()
			L.SetField(memTable, "rss", lua.LNumber(memInfo.RSS))
			L.SetField(memTable, "vms", lua.LNumber(memInfo.VMS))
			L.SetField(procTable, "memory", memTable)
		}
		
		result.RawSetInt(i+1, procTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaProcessInfo(L *lua.LState) int {
	pid := L.CheckInt(1)
	
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	L.SetField(result, "pid", lua.LNumber(pid))
	
	if name, err := proc.Name(); err == nil {
		L.SetField(result, "name", lua.LString(name))
	}
	
	if exe, err := proc.Exe(); err == nil {
		L.SetField(result, "exe", lua.LString(exe))
	}
	
	if cmdline, err := proc.Cmdline(); err == nil {
		L.SetField(result, "cmdline", lua.LString(cmdline))
	}
	
	if cwd, err := proc.Cwd(); err == nil {
		L.SetField(result, "cwd", lua.LString(cwd))
	}
	
	if status, err := proc.Status(); err == nil {
		L.SetField(result, "status", lua.LString(strings.Join(status, ",")))
	}
	
	if createTime, err := proc.CreateTime(); err == nil {
		L.SetField(result, "create_time", lua.LNumber(createTime))
	}
	
	if cpuPercent, err := proc.CPUPercent(); err == nil {
		L.SetField(result, "cpu_percent", lua.LNumber(cpuPercent))
	}
	
	if memInfo, err := proc.MemoryInfo(); err == nil {
		memTable := L.NewTable()
		L.SetField(memTable, "rss", lua.LNumber(memInfo.RSS))
		L.SetField(memTable, "vms", lua.LNumber(memInfo.VMS))
		L.SetField(result, "memory", memTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaKillProcess(L *lua.LState) int {
	pid := L.CheckInt(1)
	signal := L.OptInt(2, int(syscall.SIGTERM))
	
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	err = proc.SendSignal(syscall.Signal(signal))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// System information
func (s *SystemModule) luaHostInfo(L *lua.LState) int {
	hostInfo, err := host.Info()
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	result := L.NewTable()
	L.SetField(result, "hostname", lua.LString(hostInfo.Hostname))
	L.SetField(result, "uptime", lua.LNumber(hostInfo.Uptime))
	L.SetField(result, "boot_time", lua.LNumber(hostInfo.BootTime))
	L.SetField(result, "procs", lua.LNumber(hostInfo.Procs))
	L.SetField(result, "os", lua.LString(hostInfo.OS))
	L.SetField(result, "platform", lua.LString(hostInfo.Platform))
	L.SetField(result, "platform_family", lua.LString(hostInfo.PlatformFamily))
	L.SetField(result, "platform_version", lua.LString(hostInfo.PlatformVersion))
	L.SetField(result, "kernel_version", lua.LString(hostInfo.KernelVersion))
	L.SetField(result, "kernel_arch", lua.LString(hostInfo.KernelArch))
	L.SetField(result, "virtualization_system", lua.LString(hostInfo.VirtualizationSystem))
	L.SetField(result, "virtualization_role", lua.LString(hostInfo.VirtualizationRole))
	L.SetField(result, "host_id", lua.LString(hostInfo.HostID))

	L.Push(result)
	return 1
}

func (s *SystemModule) luaUptime(L *lua.LState) int {
	hostInfo, err := host.Info()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	uptime := time.Duration(hostInfo.Uptime) * time.Second
	
	result := L.NewTable()
	L.SetField(result, "seconds", lua.LNumber(hostInfo.Uptime))
	L.SetField(result, "human", lua.LString(uptime.String()))
	L.SetField(result, "days", lua.LNumber(hostInfo.Uptime/86400))
	L.SetField(result, "hours", lua.LNumber((hostInfo.Uptime%86400)/3600))
	L.SetField(result, "minutes", lua.LNumber((hostInfo.Uptime%3600)/60))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaEnvironment(L *lua.LState) int {
	envVars := os.Environ()
	result := L.NewTable()

	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			result.RawSetString(parts[0], lua.LString(parts[1]))
		}
	}

	L.Push(result)
	return 1
}

func (s *SystemModule) luaUsers(L *lua.LState) int {
	// Try to read /etc/passwd on Unix systems
	if runtime.GOOS != "windows" {
		file, err := os.Open("/etc/passwd")
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		defer file.Close()
		
		result := L.NewTable()
		scanner := bufio.NewScanner(file)
		userIndex := 1
		
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") && line != "" {
				parts := strings.Split(line, ":")
				if len(parts) >= 7 {
					userTable := L.NewTable()
					L.SetField(userTable, "username", lua.LString(parts[0]))
					L.SetField(userTable, "uid", lua.LString(parts[2]))
					L.SetField(userTable, "gid", lua.LString(parts[3]))
					L.SetField(userTable, "home", lua.LString(parts[5]))
					L.SetField(userTable, "shell", lua.LString(parts[6]))
					
					result.RawSetInt(userIndex, userTable)
					userIndex++
				}
			}
		}
		
		L.Push(lua.LTrue)
		L.Push(result)
		return 2
	}
	
	L.Push(lua.LNil)
	L.Push(lua.LString("user enumeration not supported on this platform"))
	return 2
}

// Performance monitoring
func (s *SystemModule) luaPerformanceSnapshot(L *lua.LState) int {
	result := L.NewTable()
	
	// CPU usage
	if cpuUsage, err := cpu.Percent(time.Second, false); err == nil && len(cpuUsage) > 0 {
		L.SetField(result, "cpu_percent", lua.LNumber(cpuUsage[0]))
	}
	
	// Memory usage
	if memInfo, err := mem.VirtualMemory(); err == nil {
		memTable := L.NewTable()
		L.SetField(memTable, "percent", lua.LNumber(memInfo.UsedPercent))
		L.SetField(memTable, "used_gb", lua.LNumber(float64(memInfo.Used)/1024/1024/1024))
		L.SetField(memTable, "total_gb", lua.LNumber(float64(memInfo.Total)/1024/1024/1024))
		L.SetField(result, "memory", memTable)
	}
	
	// Load average
	if loadAvg, err := load.Avg(); err == nil {
		loadTable := L.NewTable()
		L.SetField(loadTable, "load1", lua.LNumber(loadAvg.Load1))
		L.SetField(loadTable, "load5", lua.LNumber(loadAvg.Load5))
		L.SetField(loadTable, "load15", lua.LNumber(loadAvg.Load15))
		L.SetField(result, "load", loadTable)
	}
	
	// Disk usage for root
	if diskInfo, err := disk.Usage("/"); err == nil {
		diskTable := L.NewTable()
		L.SetField(diskTable, "percent", lua.LNumber(diskInfo.UsedPercent))
		L.SetField(diskTable, "used_gb", lua.LNumber(float64(diskInfo.Used)/1024/1024/1024))
		L.SetField(diskTable, "total_gb", lua.LNumber(float64(diskInfo.Total)/1024/1024/1024))
		L.SetField(result, "disk", diskTable)
	}
	
	// Timestamp
	L.SetField(result, "timestamp", lua.LNumber(time.Now().Unix()))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SystemModule) luaSystemHealth(L *lua.LState) int {
	result := L.NewTable()
	score := 100.0
	issues := L.NewTable()
	issueCount := 1
	
	// Check CPU usage
	if cpuUsage, err := cpu.Percent(time.Second, false); err == nil && len(cpuUsage) > 0 {
		if cpuUsage[0] > 90 {
			score -= 30
			issues.RawSetInt(issueCount, lua.LString("High CPU usage: "+fmt.Sprintf("%.1f%%", cpuUsage[0])))
			issueCount++
		} else if cpuUsage[0] > 80 {
			score -= 15
			issues.RawSetInt(issueCount, lua.LString("Elevated CPU usage: "+fmt.Sprintf("%.1f%%", cpuUsage[0])))
			issueCount++
		}
	}
	
	// Check memory usage
	if memInfo, err := mem.VirtualMemory(); err == nil {
		if memInfo.UsedPercent > 90 {
			score -= 30
			issues.RawSetInt(issueCount, lua.LString("High memory usage: "+fmt.Sprintf("%.1f%%", memInfo.UsedPercent)))
			issueCount++
		} else if memInfo.UsedPercent > 80 {
			score -= 15
			issues.RawSetInt(issueCount, lua.LString("Elevated memory usage: "+fmt.Sprintf("%.1f%%", memInfo.UsedPercent)))
			issueCount++
		}
	}
	
	// Check disk usage
	if diskInfo, err := disk.Usage("/"); err == nil {
		if diskInfo.UsedPercent > 95 {
			score -= 25
			issues.RawSetInt(issueCount, lua.LString("Critical disk usage: "+fmt.Sprintf("%.1f%%", diskInfo.UsedPercent)))
			issueCount++
		} else if diskInfo.UsedPercent > 85 {
			score -= 10
			issues.RawSetInt(issueCount, lua.LString("High disk usage: "+fmt.Sprintf("%.1f%%", diskInfo.UsedPercent)))
			issueCount++
		}
	}
	
	// Check load average
	if loadAvg, err := load.Avg(); err == nil {
		if cpuCount, err := cpu.Counts(true); err == nil {
			loadRatio := loadAvg.Load1 / float64(cpuCount)
			if loadRatio > 2.0 {
				score -= 20
				issues.RawSetInt(issueCount, lua.LString("High load average: "+fmt.Sprintf("%.2f", loadAvg.Load1)))
				issueCount++
			} else if loadRatio > 1.5 {
				score -= 10
				issues.RawSetInt(issueCount, lua.LString("Elevated load average: "+fmt.Sprintf("%.2f", loadAvg.Load1)))
				issueCount++
			}
		}
	}
	
	// Determine health status
	var status string
	if score >= 90 {
		status = "excellent"
	} else if score >= 75 {
		status = "good"
	} else if score >= 60 {
		status = "fair"
	} else if score >= 40 {
		status = "poor"
	} else {
		status = "critical"
	}
	
	L.SetField(result, "score", lua.LNumber(score))
	L.SetField(result, "status", lua.LString(status))
	L.SetField(result, "issues", issues)
	L.SetField(result, "timestamp", lua.LNumber(time.Now().Unix()))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}