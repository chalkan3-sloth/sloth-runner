package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	Services        []ServiceInfo     `json:"services"`
	Users           []UserInfo        `json:"users"`
	Environment     map[string]string `json:"environment"`
	Processes       *ProcessInfo      `json:"processes"`
	Mounts          []MountInfo       `json:"mounts"`
	Timezone        string            `json:"timezone"`
	BootTime        int64             `json:"boot_time"`
}

// ServiceInfo holds service information
type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	State  string `json:"state"`
}

// UserInfo holds user information
type UserInfo struct {
	Username string `json:"username"`
	UID      string `json:"uid"`
	GID      string `json:"gid"`
	Home     string `json:"home"`
	Shell    string `json:"shell"`
}

// ProcessInfo holds process statistics
type ProcessInfo struct {
	Total   int `json:"total"`
	Running int `json:"running"`
	Sleeping int `json:"sleeping"`
	Zombie  int `json:"zombie"`
}

// MountInfo holds filesystem mount information
type MountInfo struct {
	Device     string `json:"device"`
	Mountpoint string `json:"mountpoint"`
	FSType     string `json:"fstype"`
	Options    string `json:"options"`
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
	Manager          string              `json:"manager"`
	InstalledCount   int                 `json:"installed_count"`
	UpdatesAvailable int                 `json:"updates_available"`
	Packages         []PackageDetail     `json:"packages"`
	Updates          []PackageDetail     `json:"updates"`
}

// PackageDetail holds detailed package information
type PackageDetail struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Architecture string `json:"architecture,omitempty"`
	Description string `json:"description,omitempty"`
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

	// Processes
	info.Processes = collectProcessInfo()

	// Mounts
	info.Mounts = collectMounts()

	// Timezone
	info.Timezone = getTimezone()

	// Boot time
	info.BootTime = getBootTime()

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
	info := &PackageInfo{
		Packages: []PackageDetail{},
		Updates:  []PackageDetail{},
	}
	
	// Detect package manager and collect detailed package information
	if _, err := exec.LookPath("dpkg"); err == nil {
		info.Manager = "dpkg"
		// Get installed packages with details
		if output, err := exec.Command("dpkg-query", "-W", "-f=${Package}\t${Version}\t${Architecture}\t${binary:Summary}\n").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				fields := strings.Split(line, "\t")
				if len(fields) >= 2 {
					pkg := PackageDetail{
						Name:    fields[0],
						Version: fields[1],
					}
					if len(fields) >= 3 {
						pkg.Architecture = fields[2]
					}
					if len(fields) >= 4 {
						pkg.Description = fields[3]
					}
					info.Packages = append(info.Packages, pkg)
				}
			}
			info.InstalledCount = len(info.Packages)
		}
		// Check for updates
		if output, err := exec.Command("apt", "list", "--upgradable").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "Listing") {
					continue
				}
				// Parse: package/distro version arch [upgradable from: oldversion]
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					nameParts := strings.Split(fields[0], "/")
					pkg := PackageDetail{
						Name:    nameParts[0],
						Version: fields[1],
					}
					if len(fields) >= 3 {
						pkg.Architecture = fields[2]
					}
					info.Updates = append(info.Updates, pkg)
				}
			}
			info.UpdatesAvailable = len(info.Updates)
		}
	} else if _, err := exec.LookPath("rpm"); err == nil {
		info.Manager = "rpm"
		// Get installed packages with details
		if output, err := exec.Command("rpm", "-qa", "--queryformat", "%{NAME}\t%{VERSION}-%{RELEASE}\t%{ARCH}\t%{SUMMARY}\n").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				fields := strings.Split(line, "\t")
				if len(fields) >= 2 {
					pkg := PackageDetail{
						Name:    fields[0],
						Version: fields[1],
					}
					if len(fields) >= 3 {
						pkg.Architecture = fields[2]
					}
					if len(fields) >= 4 {
						pkg.Description = fields[3]
					}
					info.Packages = append(info.Packages, pkg)
				}
			}
			info.InstalledCount = len(info.Packages)
		}
		// Check for updates with yum or dnf
		var updateCmd *exec.Cmd
		if _, err := exec.LookPath("dnf"); err == nil {
			updateCmd = exec.Command("dnf", "check-update", "-q")
		} else {
			updateCmd = exec.Command("yum", "check-update", "-q")
		}
		if output, err := updateCmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.Contains(line, "Last metadata") || strings.Contains(line, "Security") {
					continue
				}
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pkg := PackageDetail{
						Name:    fields[0],
						Version: fields[1],
					}
					if len(fields) >= 3 {
						pkg.Architecture = fields[2]
					}
					info.Updates = append(info.Updates, pkg)
				}
			}
			info.UpdatesAvailable = len(info.Updates)
		}
	} else if _, err := exec.LookPath("pacman"); err == nil {
		info.Manager = "pacman"
		// Get installed packages with details
		if output, err := exec.Command("pacman", "-Q").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pkg := PackageDetail{
						Name:    fields[0],
						Version: fields[1],
					}
					info.Packages = append(info.Packages, pkg)
				}
			}
			info.InstalledCount = len(info.Packages)
		}
		// Check for updates
		if output, err := exec.Command("pacman", "-Qu").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					pkg := PackageDetail{
						Name:    fields[0],
						Version: fields[3], // New version
					}
					info.Updates = append(info.Updates, pkg)
				}
			}
			info.UpdatesAvailable = len(info.Updates)
		}
	} else if _, err := exec.LookPath("brew"); err == nil {
		info.Manager = "brew"
		// Get installed packages with details
		if output, err := exec.Command("brew", "list", "--versions").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pkg := PackageDetail{
						Name:    fields[0],
						Version: strings.Join(fields[1:], " "),
					}
					info.Packages = append(info.Packages, pkg)
				}
			}
			info.InstalledCount = len(info.Packages)
		}
		// Check for updates
		if output, err := exec.Command("brew", "outdated").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					pkg := PackageDetail{
						Name:    fields[0],
						Version: fields[2], // New version (format: name (old) < new)
					}
					info.Updates = append(info.Updates, pkg)
				}
			}
			info.UpdatesAvailable = len(info.Updates)
		}
	}
	
	return info
}

// collectServices collects running services with detailed information
func collectServices() []ServiceInfo {
	var services []ServiceInfo
	
	if _, err := exec.LookPath("systemctl"); err == nil {
		output, err := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend").Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					serviceName := strings.TrimSuffix(fields[0], ".service")
					service := ServiceInfo{
						Name:   serviceName,
						Status: fields[1], // loaded/not-found
						State:  fields[2], // active/inactive/failed
					}
					services = append(services, service)
				}
			}
		}
	}
	
	return services
}

// collectUsers collects system users with detailed information
func collectUsers() []UserInfo {
	var users []UserInfo
	
	switch runtime.GOOS {
	case "linux", "darwin":
		if data, err := os.ReadFile("/etc/passwd"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				fields := strings.Split(line, ":")
				if len(fields) >= 7 && fields[0] != "" {
					user := UserInfo{
						Username: fields[0],
						UID:      fields[2],
						GID:      fields[3],
						Home:     fields[5],
						Shell:    fields[6],
					}
					users = append(users, user)
				}
			}
		}
	}
	
	return users
}

// collectProcessInfo collects process statistics
func collectProcessInfo() *ProcessInfo {
	info := &ProcessInfo{}
	
	switch runtime.GOOS {
	case "linux":
		// Read /proc for process count
		if entries, err := os.ReadDir("/proc"); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				// Check if directory name is a number (PID)
				if _, err := fmt.Sscanf(entry.Name(), "%d", new(int)); err == nil {
					info.Total++
					
					// Read status file to get process state
					statusPath := filepath.Join("/proc", entry.Name(), "status")
					if data, err := os.ReadFile(statusPath); err == nil {
						lines := strings.Split(string(data), "\n")
						for _, line := range lines {
							if strings.HasPrefix(line, "State:") {
								fields := strings.Fields(line)
								if len(fields) >= 2 {
									switch fields[1] {
									case "R":
										info.Running++
									case "S", "D":
										info.Sleeping++
									case "Z":
										info.Zombie++
									}
								}
								break
							}
						}
					}
				}
			}
		}
	case "darwin":
		// Use ps command
		if output, err := exec.Command("ps", "axo", "state").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for i, line := range lines {
				if i == 0 { // Skip header
					continue
				}
				state := strings.TrimSpace(line)
				if state == "" {
					continue
				}
				info.Total++
				// macOS process states: R=running, S=sleeping, Z=zombie, etc.
				switch state[0] {
				case 'R':
					info.Running++
				case 'S', 'I':
					info.Sleeping++
				case 'Z':
					info.Zombie++
				}
			}
		}
	}
	
	return info
}

// collectMounts collects filesystem mount information
func collectMounts() []MountInfo {
	var mounts []MountInfo
	
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/mounts"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					mount := MountInfo{
						Device:     fields[0],
						Mountpoint: fields[1],
						FSType:     fields[2],
						Options:    fields[3],
					}
					mounts = append(mounts, mount)
				}
			}
		}
	case "darwin":
		if output, err := exec.Command("mount").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				// Format: device on mountpoint (fstype, options)
				if strings.Contains(line, " on ") {
					parts := strings.Split(line, " on ")
					if len(parts) == 2 {
						device := parts[0]
						rest := parts[1]
						
						// Extract mountpoint and details
						if idx := strings.Index(rest, " ("); idx != -1 {
							mountpoint := rest[:idx]
							details := strings.Trim(rest[idx+2:], ")")
							
							mount := MountInfo{
								Device:     device,
								Mountpoint: mountpoint,
							}
							
							// Parse details (fstype, options)
							detailParts := strings.SplitN(details, ", ", 2)
							if len(detailParts) >= 1 {
								mount.FSType = detailParts[0]
							}
							if len(detailParts) >= 2 {
								mount.Options = detailParts[1]
							}
							
							mounts = append(mounts, mount)
						}
					}
				}
			}
		}
	}
	
	return mounts
}

// getTimezone returns the system timezone
func getTimezone() string {
	now := time.Now()
	zone, _ := now.Zone()
	return zone
}

// getBootTime returns the system boot time as Unix timestamp
func getBootTime() int64 {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/stat"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "btime ") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						var btime int64
						fmt.Sscanf(fields[1], "%d", &btime)
						return btime
					}
				}
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
					return bootTime
				}
			}
		}
	}
	return 0
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
