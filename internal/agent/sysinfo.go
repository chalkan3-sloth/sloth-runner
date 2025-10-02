package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// SystemInfo holds comprehensive system information
type SystemInfo struct {
	Hostname        string            `json:"hostname"`
	Platform        string            `json:"platform"`
	PlatformFamily  string            `json:"platform_family"`
	PlatformVersion string            `json:"platform_version"`
	Architecture    string            `json:"architecture"`
	CPUs            int               `json:"cpus"`
	Memory          *MemoryInfo       `json:"memory"`
	Disk            []*DiskInfo       `json:"disk"`
	Network         []*NetworkInfo    `json:"network"`
	Uptime          int64             `json:"uptime"`
	LoadAverage     []float64         `json:"load_average"`
	Kernel          string            `json:"kernel"`
	KernelVersion   string            `json:"kernel_version"`
	Virtualization  string            `json:"virtualization"`
	CollectedAt     time.Time         `json:"collected_at"`
	Packages        *PackageInfo      `json:"packages"`
	Services        []string          `json:"services"`
	Users           []string          `json:"users"`
	Environment     map[string]string `json:"environment"`
}

// MemoryInfo holds memory information
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Cached      uint64  `json:"cached"`
	Buffers     uint64  `json:"buffers"`
}

// DiskInfo holds disk information
type DiskInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// NetworkInfo holds network interface information
type NetworkInfo struct {
	Name        string   `json:"name"`
	Addresses   []string `json:"addresses"`
	MAC         string   `json:"mac"`
	MTU         int      `json:"mtu"`
	IsUp        bool     `json:"is_up"`
	Speed       int64    `json:"speed"`
}

// PackageInfo holds package manager information
type PackageInfo struct {
	Manager        string   `json:"manager"`
	InstalledCount int      `json:"installed_count"`
	UpdatesAvailable int    `json:"updates_available"`
}

// CollectSystemInfo collects comprehensive system information
func CollectSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{
		Architecture: runtime.GOARCH,
		CPUs:         runtime.NumCPU(),
		CollectedAt:  time.Now(),
		Environment:  make(map[string]string),
	}

	// Hostname
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	// Platform detection
	info.Platform = runtime.GOOS
	detectPlatformDetails(info)

	// Memory information
	info.Memory = collectMemoryInfo()

	// Disk information
	info.Disk = collectDiskInfo()

	// Network information
	info.Network = collectNetworkInfo()

	// Uptime
	info.Uptime = getUptime()

	// Load average (Unix-like systems)
	if runtime.GOOS != "windows" {
		info.LoadAverage = getLoadAverage()
	}

	// Kernel information
	info.Kernel, info.KernelVersion = getKernelInfo()

	// Virtualization detection
	info.Virtualization = detectVirtualization()

	// Package information
	info.Packages = collectPackageInfo()

	// Services (systemd-based)
	info.Services = collectServices()

	// Users
	info.Users = collectUsers()

	// Select environment variables
	for _, env := range []string{"PATH", "HOME", "USER", "SHELL", "LANG"} {
		if val := os.Getenv(env); val != "" {
			info.Environment[env] = val
		}
	}

	return info, nil
}

// detectPlatformDetails detects detailed platform information
func detectPlatformDetails(info *SystemInfo) {
	switch runtime.GOOS {
	case "linux":
		info.PlatformFamily = "linux"
		// Try to read /etc/os-release
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "ID=") {
					info.Platform = strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
				}
				if strings.HasPrefix(line, "VERSION_ID=") {
					info.PlatformVersion = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
				}
			}
		}
	case "darwin":
		info.PlatformFamily = "darwin"
		info.Platform = "macos"
		// Get macOS version
		if output, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
			info.PlatformVersion = strings.TrimSpace(string(output))
		}
	case "windows":
		info.PlatformFamily = "windows"
	}
}

// collectMemoryInfo collects memory information
func collectMemoryInfo() *MemoryInfo {
	mem := &MemoryInfo{}
	
	switch runtime.GOOS {
	case "linux":
		// Read /proc/meminfo
		if data, err := os.ReadFile("/proc/meminfo"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) < 2 {
					continue
				}
				var value uint64
				fmt.Sscanf(fields[1], "%d", &value)
				value *= 1024 // Convert from KB to bytes
				
				switch fields[0] {
				case "MemTotal:":
					mem.Total = value
				case "MemAvailable:":
					mem.Available = value
				case "MemFree:":
					mem.Free = value
				case "Cached:":
					mem.Cached = value
				case "Buffers:":
					mem.Buffers = value
				}
			}
			mem.Used = mem.Total - mem.Available
			if mem.Total > 0 {
				mem.UsedPercent = float64(mem.Used) / float64(mem.Total) * 100
			}
		}
	case "darwin":
		// Use vm_stat for macOS
		if output, err := exec.Command("vm_stat").Output(); err == nil {
			pageSize := uint64(4096)
			lines := strings.Split(string(output), "\n")
			var free, active, inactive, speculative, wired uint64
			
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) < 2 {
					continue
				}
				value := strings.TrimSuffix(fields[len(fields)-1], ".")
				var num uint64
				fmt.Sscanf(value, "%d", &num)
				
				if strings.Contains(line, "Pages free:") {
					free = num * pageSize
				} else if strings.Contains(line, "Pages active:") {
					active = num * pageSize
				} else if strings.Contains(line, "Pages inactive:") {
					inactive = num * pageSize
				} else if strings.Contains(line, "Pages speculative:") {
					speculative = num * pageSize
				} else if strings.Contains(line, "Pages wired down:") {
					wired = num * pageSize
				}
			}
			
			// Get total memory using sysctl
			if output, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
				fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &mem.Total)
			}
			
			mem.Free = free
			mem.Used = active + wired
			mem.Available = free + inactive + speculative
			mem.Cached = inactive
			if mem.Total > 0 {
				mem.UsedPercent = float64(mem.Used) / float64(mem.Total) * 100
			}
		}
	}
	
	return mem
}

// collectDiskInfo collects disk information
func collectDiskInfo() []*DiskInfo {
	var disks []*DiskInfo
	
	switch runtime.GOOS {
	case "linux", "darwin":
		// Use df command
		output, err := exec.Command("df", "-Pk").Output()
		if err != nil {
			return disks
		}
		
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i == 0 { // Skip header
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 6 {
				continue
			}
			
			// Skip tmpfs and other virtual filesystems on Linux
			if runtime.GOOS == "linux" {
				fstype := fields[0]
				if strings.HasPrefix(fstype, "tmpfs") || strings.HasPrefix(fstype, "devtmpfs") {
					continue
				}
			}
			
			disk := &DiskInfo{
				Device:     fields[0],
				Mountpoint: fields[5],
			}
			
			fmt.Sscanf(fields[1], "%d", &disk.Total)
			fmt.Sscanf(fields[2], "%d", &disk.Used)
			fmt.Sscanf(fields[3], "%d", &disk.Free)
			
			disk.Total *= 1024
			disk.Used *= 1024
			disk.Free *= 1024
			
			if disk.Total > 0 {
				disk.UsedPercent = float64(disk.Used) / float64(disk.Total) * 100
			}
			
			disks = append(disks, disk)
		}
	}
	
	return disks
}

// collectNetworkInfo collects network interface information
func collectNetworkInfo() []*NetworkInfo {
	var networks []*NetworkInfo
	
	switch runtime.GOOS {
	case "linux":
		// Use ip command
		if output, err := exec.Command("ip", "-o", "addr", "show").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			interfaceMap := make(map[string]*NetworkInfo)
			
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) < 4 {
					continue
				}
				
				name := fields[1]
				if _, exists := interfaceMap[name]; !exists {
					interfaceMap[name] = &NetworkInfo{
						Name:      name,
						Addresses: []string{},
					}
				}
				
				// Check if interface is up
				if strings.Contains(line, "UP") {
					interfaceMap[name].IsUp = true
				}
				
				// Get IP address
				if len(fields) >= 4 && (fields[2] == "inet" || fields[2] == "inet6") {
					addr := strings.Split(fields[3], "/")[0]
					interfaceMap[name].Addresses = append(interfaceMap[name].Addresses, addr)
				}
			}
			
			// Get MAC addresses
			if output, err := exec.Command("ip", "link", "show").Output(); err == nil {
				lines := strings.Split(string(output), "\n")
				var currentInterface string
				for _, line := range lines {
					if strings.HasPrefix(line, " ") {
						// This is a continuation line
						if currentInterface != "" && strings.Contains(line, "link/ether") {
							fields := strings.Fields(line)
							for i, field := range fields {
								if field == "link/ether" && i+1 < len(fields) {
									if iface, exists := interfaceMap[currentInterface]; exists {
										iface.MAC = fields[i+1]
									}
								}
							}
						}
					} else {
						fields := strings.Fields(line)
						if len(fields) >= 2 {
							currentInterface = strings.TrimSuffix(fields[1], ":")
						}
					}
				}
			}
			
			for _, iface := range interfaceMap {
				networks = append(networks, iface)
			}
		}
	case "darwin":
		// Use ifconfig for macOS
		if output, err := exec.Command("ifconfig").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			var current *NetworkInfo
			
			for _, line := range lines {
				if !strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, " ") && line != "" {
					// New interface
					fields := strings.Fields(line)
					if len(fields) > 0 {
						if current != nil {
							networks = append(networks, current)
						}
						current = &NetworkInfo{
							Name:      strings.TrimSuffix(fields[0], ":"),
							Addresses: []string{},
						}
						if strings.Contains(line, "UP") {
							current.IsUp = true
						}
					}
				} else if current != nil {
					// Parse interface details
					line = strings.TrimSpace(line)
					if strings.HasPrefix(line, "inet ") {
						fields := strings.Fields(line)
						if len(fields) >= 2 {
							current.Addresses = append(current.Addresses, fields[1])
						}
					} else if strings.HasPrefix(line, "ether ") {
						fields := strings.Fields(line)
						if len(fields) >= 2 {
							current.MAC = fields[1]
						}
					}
				}
			}
			if current != nil {
				networks = append(networks, current)
			}
		}
	}
	
	return networks
}

// getUptime returns system uptime in seconds
func getUptime() int64 {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/uptime"); err == nil {
			fields := strings.Fields(string(data))
			if len(fields) > 0 {
				var uptime float64
				fmt.Sscanf(fields[0], "%f", &uptime)
				return int64(uptime)
			}
		}
	case "darwin":
		if output, err := exec.Command("sysctl", "-n", "kern.boottime").Output(); err == nil {
			// Parse: { sec = 1234567890, usec = 0 } Mon Jan  1 00:00:00 2024
			line := string(output)
			if idx := strings.Index(line, "sec = "); idx != -1 {
				line = line[idx+6:]
				if idx := strings.Index(line, ","); idx != -1 {
					var bootTime int64
					fmt.Sscanf(line[:idx], "%d", &bootTime)
					return time.Now().Unix() - bootTime
				}
			}
		}
	}
	return 0
}

// getLoadAverage returns system load average
func getLoadAverage() []float64 {
	loads := make([]float64, 3)
	
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/loadavg"); err == nil {
			fields := strings.Fields(string(data))
			if len(fields) >= 3 {
				fmt.Sscanf(fields[0], "%f", &loads[0])
				fmt.Sscanf(fields[1], "%f", &loads[1])
				fmt.Sscanf(fields[2], "%f", &loads[2])
			}
		}
	case "darwin":
		if output, err := exec.Command("sysctl", "-n", "vm.loadavg").Output(); err == nil {
			// Parse: { 1.23 2.34 3.45 }
			line := strings.Trim(string(output), "{}\n ")
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				fmt.Sscanf(fields[0], "%f", &loads[0])
				fmt.Sscanf(fields[1], "%f", &loads[1])
				fmt.Sscanf(fields[2], "%f", &loads[2])
			}
		}
	}
	
	return loads
}

// getKernelInfo returns kernel name and version
func getKernelInfo() (string, string) {
	kernel := runtime.GOOS
	version := ""
	
	if output, err := exec.Command("uname", "-r").Output(); err == nil {
		version = strings.TrimSpace(string(output))
	}
	
	if output, err := exec.Command("uname", "-s").Output(); err == nil {
		kernel = strings.TrimSpace(string(output))
	}
	
	return kernel, version
}

// detectVirtualization detects if running in a virtualized environment
func detectVirtualization() string {
	switch runtime.GOOS {
	case "linux":
		// Check systemd-detect-virt
		if output, err := exec.Command("systemd-detect-virt").Output(); err == nil {
			virt := strings.TrimSpace(string(output))
			if virt != "none" {
				return virt
			}
		}
		
		// Check /sys/class/dmi/id/product_name
		if data, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
			product := strings.ToLower(strings.TrimSpace(string(data)))
			if strings.Contains(product, "virtualbox") {
				return "virtualbox"
			} else if strings.Contains(product, "vmware") {
				return "vmware"
			} else if strings.Contains(product, "kvm") {
				return "kvm"
			} else if strings.Contains(product, "qemu") {
				return "qemu"
			}
		}
	case "darwin":
		// Check for common VM indicators
		if output, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
			brand := strings.ToLower(string(output))
			if strings.Contains(brand, "qemu") {
				return "qemu"
			}
		}
	}
	
	return "none"
}

// collectPackageInfo collects package manager information
func collectPackageInfo() *PackageInfo {
	info := &PackageInfo{}
	
	// Detect package manager and count packages
	if _, err := exec.LookPath("dpkg"); err == nil {
		info.Manager = "dpkg"
		if output, err := exec.Command("dpkg", "-l").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			count := 0
			for _, line := range lines {
				if strings.HasPrefix(line, "ii ") {
					count++
				}
			}
			info.InstalledCount = count
		}
		// Check for updates
		if output, err := exec.Command("apt", "list", "--upgradable").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			info.UpdatesAvailable = len(lines) - 1 // Subtract header
			if info.UpdatesAvailable < 0 {
				info.UpdatesAvailable = 0
			}
		}
	} else if _, err := exec.LookPath("rpm"); err == nil {
		info.Manager = "rpm"
		if output, err := exec.Command("rpm", "-qa").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			info.InstalledCount = len(lines) - 1
		}
		// Check for updates
		if output, err := exec.Command("yum", "check-update").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			count := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" && !strings.Contains(line, "Last metadata") {
					count++
				}
			}
			info.UpdatesAvailable = count
		}
	} else if _, err := exec.LookPath("pacman"); err == nil {
		info.Manager = "pacman"
		if output, err := exec.Command("pacman", "-Q").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			info.InstalledCount = len(lines) - 1
		}
		// Check for updates
		if output, err := exec.Command("pacman", "-Qu").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			info.UpdatesAvailable = len(lines) - 1
			if info.UpdatesAvailable < 0 {
				info.UpdatesAvailable = 0
			}
		}
	} else if _, err := exec.LookPath("brew"); err == nil {
		info.Manager = "brew"
		if output, err := exec.Command("brew", "list").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			info.InstalledCount = len(lines) - 1
		}
		// Check for updates
		if output, err := exec.Command("brew", "outdated").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			info.UpdatesAvailable = len(lines) - 1
			if info.UpdatesAvailable < 0 {
				info.UpdatesAvailable = 0
			}
		}
	}
	
	return info
}

// collectServices collects running services (systemd-based systems)
func collectServices() []string {
	var services []string
	
	if _, err := exec.LookPath("systemctl"); err == nil {
		output, err := exec.Command("systemctl", "list-units", "--type=service", "--state=running", "--no-pager", "--no-legend").Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					serviceName := strings.TrimSuffix(fields[0], ".service")
					services = append(services, serviceName)
				}
			}
		}
	}
	
	return services
}

// collectUsers collects system users
func collectUsers() []string {
	var users []string
	
	switch runtime.GOOS {
	case "linux", "darwin":
		if data, err := os.ReadFile("/etc/passwd"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				fields := strings.Split(line, ":")
				if len(fields) > 0 && fields[0] != "" {
					users = append(users, fields[0])
				}
			}
		}
	}
	
	return users
}

// ToJSON converts SystemInfo to JSON string
func (s *SystemInfo) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON parses JSON string to SystemInfo
func FromJSON(jsonStr string) (*SystemInfo, error) {
	var info SystemInfo
	err := json.Unmarshal([]byte(jsonStr), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
