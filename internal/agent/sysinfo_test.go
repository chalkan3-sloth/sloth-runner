package agent

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// Test SystemInfo struct creation
func TestSystemInfo_Creation(t *testing.T) {
	info := &SystemInfo{
		Hostname:     "test-host",
		Platform:     "linux",
		Architecture: "amd64",
		CPUs:         4,
	}

	if info.Hostname != "test-host" {
		t.Errorf("Expected hostname 'test-host', got '%s'", info.Hostname)
	}

	if info.CPUs != 4 {
		t.Errorf("Expected 4 CPUs, got %d", info.CPUs)
	}
}

func TestSystemInfo_ZeroValue(t *testing.T) {
	info := &SystemInfo{}

	if info.Hostname != "" {
		t.Error("Expected empty hostname")
	}

	if info.CPUs != 0 {
		t.Error("Expected 0 CPUs")
	}
}

// Test MemoryInfo struct
func TestMemoryInfo_Creation(t *testing.T) {
	mem := &MemoryInfo{
		Total:       8 * 1024 * 1024 * 1024, // 8GB
		Available:   4 * 1024 * 1024 * 1024, // 4GB
		Used:        4 * 1024 * 1024 * 1024, // 4GB
		UsedPercent: 50.0,
	}

	if mem.Total != 8*1024*1024*1024 {
		t.Error("Expected 8GB total memory")
	}

	if mem.UsedPercent != 50.0 {
		t.Errorf("Expected 50%% used, got %f%%", mem.UsedPercent)
	}
}

func TestMemoryInfo_HighUsage(t *testing.T) {
	mem := &MemoryInfo{
		Total:       16 * 1024 * 1024 * 1024,
		Used:        15 * 1024 * 1024 * 1024,
		Available:   1 * 1024 * 1024 * 1024,
		UsedPercent: 93.75,
	}

	if mem.UsedPercent < 90.0 {
		t.Error("Expected high memory usage")
	}
}

func TestMemoryInfo_ZeroValue(t *testing.T) {
	mem := &MemoryInfo{}

	if mem.Total != 0 || mem.Used != 0 {
		t.Error("Expected zero values")
	}
}

// Test DiskInfo struct
func TestDiskInfo_Creation(t *testing.T) {
	disk := &DiskInfo{
		Device:      "/dev/sda1",
		Mountpoint:  "/",
		Fstype:      "ext4",
		Total:       100 * 1024 * 1024 * 1024, // 100GB
		Used:        50 * 1024 * 1024 * 1024,  // 50GB
		Free:        50 * 1024 * 1024 * 1024,
		UsedPercent: 50.0,
	}

	if disk.Device != "/dev/sda1" {
		t.Errorf("Expected device '/dev/sda1', got '%s'", disk.Device)
	}

	if disk.UsedPercent != 50.0 {
		t.Errorf("Expected 50%% used, got %f%%", disk.UsedPercent)
	}
}

func TestDiskInfo_MultipleMountpoints(t *testing.T) {
	disks := []*DiskInfo{
		{Device: "/dev/sda1", Mountpoint: "/", Fstype: "ext4"},
		{Device: "/dev/sda2", Mountpoint: "/home", Fstype: "ext4"},
		{Device: "/dev/sdb1", Mountpoint: "/data", Fstype: "xfs"},
	}

	if len(disks) != 3 {
		t.Errorf("Expected 3 disks, got %d", len(disks))
	}

	for _, disk := range disks {
		if disk.Mountpoint == "" {
			t.Error("Expected non-empty mountpoint")
		}
	}
}

func TestDiskInfo_DifferentFilesystems(t *testing.T) {
	filesystems := []string{"ext4", "xfs", "btrfs", "zfs", "ntfs"}

	for _, fs := range filesystems {
		disk := &DiskInfo{
			Device:     "/dev/test",
			Mountpoint: "/test",
			Fstype:     fs,
		}

		if disk.Fstype != fs {
			t.Errorf("Expected filesystem '%s', got '%s'", fs, disk.Fstype)
		}
	}
}

// Test NetworkInfo struct
func TestNetworkInfo_Creation(t *testing.T) {
	net := &NetworkInfo{
		Name:      "eth0",
		Addresses: []string{"192.168.1.100", "fe80::1"},
		MAC:       "00:11:22:33:44:55",
		MTU:       1500,
		IsUp:      true,
	}

	if net.Name != "eth0" {
		t.Errorf("Expected interface 'eth0', got '%s'", net.Name)
	}

	if len(net.Addresses) != 2 {
		t.Errorf("Expected 2 addresses, got %d", len(net.Addresses))
	}

	if !net.IsUp {
		t.Error("Expected interface to be up")
	}
}

func TestNetworkInfo_NoAddresses(t *testing.T) {
	net := &NetworkInfo{
		Name:      "lo",
		Addresses: []string{},
		IsUp:      false,
	}

	if len(net.Addresses) != 0 {
		t.Error("Expected no addresses")
	}

	if net.IsUp {
		t.Error("Expected interface to be down")
	}
}

func TestNetworkInfo_MultipleInterfaces(t *testing.T) {
	interfaces := []NetworkInfo{
		{Name: "lo", Addresses: []string{"127.0.0.1"}},
		{Name: "eth0", Addresses: []string{"192.168.1.100"}},
		{Name: "wlan0", Addresses: []string{"192.168.1.101"}},
	}

	if len(interfaces) != 3 {
		t.Errorf("Expected 3 interfaces, got %d", len(interfaces))
	}
}

// Test ServiceInfo struct
func TestServiceInfo_Creation(t *testing.T) {
	service := ServiceInfo{
		Name:   "nginx",
		Status: "loaded",
		State:  "active",
	}

	if service.Name != "nginx" {
		t.Errorf("Expected service 'nginx', got '%s'", service.Name)
	}

	if service.State != "active" {
		t.Errorf("Expected state 'active', got '%s'", service.State)
	}
}

func TestServiceInfo_MultipleServices(t *testing.T) {
	services := []ServiceInfo{
		{Name: "nginx", Status: "loaded", State: "active"},
		{Name: "mysql", Status: "loaded", State: "active"},
		{Name: "redis", Status: "loaded", State: "inactive"},
	}

	activeCount := 0
	for _, svc := range services {
		if svc.State == "active" {
			activeCount++
		}
	}

	if activeCount != 2 {
		t.Errorf("Expected 2 active services, got %d", activeCount)
	}
}

func TestServiceInfo_FailedState(t *testing.T) {
	service := ServiceInfo{
		Name:   "broken-service",
		Status: "loaded",
		State:  "failed",
	}

	if service.State != "failed" {
		t.Error("Expected failed state")
	}
}

// Test UserInfo struct
func TestUserInfo_Creation(t *testing.T) {
	user := UserInfo{
		Username: "testuser",
		UID:      "1000",
		GID:      "1000",
		Home:     "/home/testuser",
		Shell:    "/bin/bash",
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	if user.Shell != "/bin/bash" {
		t.Errorf("Expected shell '/bin/bash', got '%s'", user.Shell)
	}
}

func TestUserInfo_SystemUsers(t *testing.T) {
	users := []UserInfo{
		{Username: "root", UID: "0", GID: "0", Home: "/root", Shell: "/bin/bash"},
		{Username: "nobody", UID: "65534", GID: "65534", Home: "/nonexistent", Shell: "/usr/sbin/nologin"},
	}

	for _, user := range users {
		if user.Username == "" {
			t.Error("Expected non-empty username")
		}
	}
}

func TestUserInfo_DifferentShells(t *testing.T) {
	shells := []string{"/bin/bash", "/bin/zsh", "/bin/sh", "/usr/bin/fish", "/usr/sbin/nologin"}

	for _, shell := range shells {
		user := UserInfo{
			Username: "test",
			Shell:    shell,
		}

		if user.Shell != shell {
			t.Errorf("Expected shell '%s', got '%s'", shell, user.Shell)
		}
	}
}

// Test ProcessInfo struct
func TestProcessInfo_Creation(t *testing.T) {
	proc := &ProcessInfo{
		Total:    150,
		Running:  3,
		Sleeping: 140,
		Zombie:   7,
	}

	if proc.Total != 150 {
		t.Errorf("Expected 150 total processes, got %d", proc.Total)
	}

	sum := proc.Running + proc.Sleeping + proc.Zombie
	if sum != proc.Total {
		t.Errorf("Process counts don't add up: %d != %d", sum, proc.Total)
	}
}

func TestProcessInfo_HighLoad(t *testing.T) {
	proc := &ProcessInfo{
		Total:    1000,
		Running:  150,
		Sleeping: 800,
		Zombie:   50,
	}

	if proc.Total != 1000 {
		t.Errorf("Expected 1000 total processes, got %d", proc.Total)
	}

	if proc.Zombie > 0 {
		// High zombie count may indicate issues
		if proc.Zombie > 100 {
			t.Error("Very high zombie count detected")
		}
	}
}

func TestProcessInfo_ZeroValue(t *testing.T) {
	proc := &ProcessInfo{}

	if proc.Total != 0 {
		t.Error("Expected zero total processes")
	}
}

// Test MountInfo struct
func TestMountInfo_Creation(t *testing.T) {
	mount := MountInfo{
		Device:     "/dev/sda1",
		Mountpoint: "/",
		FSType:     "ext4",
		Options:    "rw,relatime",
	}

	if mount.Device != "/dev/sda1" {
		t.Errorf("Expected device '/dev/sda1', got '%s'", mount.Device)
	}

	if mount.FSType != "ext4" {
		t.Errorf("Expected fstype 'ext4', got '%s'", mount.FSType)
	}
}

func TestMountInfo_MultipleMounts(t *testing.T) {
	mounts := []MountInfo{
		{Device: "/dev/sda1", Mountpoint: "/", FSType: "ext4"},
		{Device: "/dev/sda2", Mountpoint: "/home", FSType: "ext4"},
		{Device: "tmpfs", Mountpoint: "/tmp", FSType: "tmpfs"},
	}

	if len(mounts) != 3 {
		t.Errorf("Expected 3 mounts, got %d", len(mounts))
	}
}

func TestMountInfo_ReadOnlyMount(t *testing.T) {
	mount := MountInfo{
		Device:     "/dev/cdrom",
		Mountpoint: "/mnt/cdrom",
		FSType:     "iso9660",
		Options:    "ro",
	}

	if !strings.Contains(mount.Options, "ro") {
		t.Error("Expected read-only mount")
	}
}

// Test PackageInfo struct
func TestPackageInfo_Creation(t *testing.T) {
	pkg := &PackageInfo{
		Manager:          "apt",
		InstalledCount:   500,
		UpdatesAvailable: 10,
		Packages:         []PackageDetail{},
		Updates:          []PackageDetail{},
	}

	if pkg.Manager != "apt" {
		t.Errorf("Expected manager 'apt', got '%s'", pkg.Manager)
	}

	if pkg.InstalledCount != 500 {
		t.Errorf("Expected 500 packages, got %d", pkg.InstalledCount)
	}
}

func TestPackageInfo_MultipleManagers(t *testing.T) {
	managers := []string{"apt", "yum", "dnf", "pacman", "brew", "rpm"}

	for _, mgr := range managers {
		pkg := &PackageInfo{
			Manager: mgr,
		}

		if pkg.Manager != mgr {
			t.Errorf("Expected manager '%s', got '%s'", mgr, pkg.Manager)
		}
	}
}

func TestPackageInfo_NoUpdates(t *testing.T) {
	pkg := &PackageInfo{
		Manager:          "apt",
		InstalledCount:   100,
		UpdatesAvailable: 0,
	}

	if pkg.UpdatesAvailable != 0 {
		t.Error("Expected no updates available")
	}
}

// Test PackageDetail struct
func TestPackageDetail_Creation(t *testing.T) {
	pkg := PackageDetail{
		Name:         "nginx",
		Version:      "1.18.0-1",
		Architecture: "amd64",
		Description:  "High performance web server",
	}

	if pkg.Name != "nginx" {
		t.Errorf("Expected package 'nginx', got '%s'", pkg.Name)
	}

	if pkg.Version == "" {
		t.Error("Expected non-empty version")
	}
}

func TestPackageDetail_MultiplePackages(t *testing.T) {
	packages := []PackageDetail{
		{Name: "nginx", Version: "1.18.0"},
		{Name: "mysql", Version: "8.0.23"},
		{Name: "redis", Version: "6.0.10"},
	}

	if len(packages) != 3 {
		t.Errorf("Expected 3 packages, got %d", len(packages))
	}
}

func TestPackageDetail_DifferentArchitectures(t *testing.T) {
	architectures := []string{"amd64", "arm64", "i386", "armhf"}

	for _, arch := range architectures {
		pkg := PackageDetail{
			Name:         "test-package",
			Version:      "1.0.0",
			Architecture: arch,
		}

		if pkg.Architecture != arch {
			t.Errorf("Expected architecture '%s', got '%s'", arch, pkg.Architecture)
		}
	}
}

// Test SystemInfo JSON marshaling
func TestSystemInfo_ToJSON(t *testing.T) {
	info := &SystemInfo{
		Hostname:     "test-host",
		Platform:     "linux",
		Architecture: "amd64",
		CPUs:         4,
		CollectedAt:  time.Now(),
	}

	jsonStr, err := info.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	if !strings.Contains(jsonStr, "test-host") {
		t.Error("JSON should contain hostname")
	}
}

func TestSystemInfo_ToJSON_Empty(t *testing.T) {
	info := &SystemInfo{}

	jsonStr, err := info.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal empty struct: %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}
}

func TestSystemInfo_ToJSON_WithMemory(t *testing.T) {
	info := &SystemInfo{
		Hostname: "test",
		Memory: &MemoryInfo{
			Total:       8 * 1024 * 1024 * 1024,
			Used:        4 * 1024 * 1024 * 1024,
			UsedPercent: 50.0,
		},
	}

	jsonStr, err := info.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	if !strings.Contains(jsonStr, "memory") {
		t.Error("JSON should contain memory info")
	}
}

// Test SystemInfo JSON unmarshaling
func TestFromJSON_Valid(t *testing.T) {
	jsonStr := `{"hostname":"test-host","platform":"linux","architecture":"amd64","cpus":4}`

	info, err := FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if info.Hostname != "test-host" {
		t.Errorf("Expected hostname 'test-host', got '%s'", info.Hostname)
	}

	if info.CPUs != 4 {
		t.Errorf("Expected 4 CPUs, got %d", info.CPUs)
	}
}

func TestFromJSON_Invalid(t *testing.T) {
	jsonStr := `{invalid json}`

	_, err := FromJSON(jsonStr)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestFromJSON_Empty(t *testing.T) {
	jsonStr := `{}`

	info, err := FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty JSON: %v", err)
	}

	if info == nil {
		t.Error("Expected non-nil info")
	}
}

// Test SystemInfo roundtrip (JSON -> struct -> JSON)
func TestSystemInfo_JSONRoundtrip(t *testing.T) {
	original := &SystemInfo{
		Hostname:     "roundtrip-test",
		Platform:     "linux",
		Architecture: "amd64",
		CPUs:         8,
		Memory: &MemoryInfo{
			Total: 16 * 1024 * 1024 * 1024,
			Used:  8 * 1024 * 1024 * 1024,
		},
	}

	// Marshal to JSON
	jsonStr, err := original.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal from JSON
	restored, err := FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare
	if restored.Hostname != original.Hostname {
		t.Errorf("Hostname mismatch: expected '%s', got '%s'", original.Hostname, restored.Hostname)
	}

	if restored.CPUs != original.CPUs {
		t.Errorf("CPUs mismatch: expected %d, got %d", original.CPUs, restored.CPUs)
	}
}

// Test SystemInfo with complex nested structures
func TestSystemInfo_ComplexStructure(t *testing.T) {
	info := &SystemInfo{
		Hostname: "complex-test",
		Memory: &MemoryInfo{
			Total: 16 * 1024 * 1024 * 1024,
		},
		Disk: []*DiskInfo{
			{Device: "/dev/sda1", Mountpoint: "/"},
			{Device: "/dev/sda2", Mountpoint: "/home"},
		},
		Network: []*NetworkInfo{
			{Name: "eth0", Addresses: []string{"192.168.1.100"}},
		},
		Services: []ServiceInfo{
			{Name: "nginx", State: "active"},
		},
		Users: []UserInfo{
			{Username: "root", UID: "0"},
		},
		Environment: map[string]string{
			"PATH": "/usr/bin:/bin",
			"HOME": "/root",
		},
	}

	if len(info.Disk) != 2 {
		t.Errorf("Expected 2 disks, got %d", len(info.Disk))
	}

	if len(info.Network) != 1 {
		t.Errorf("Expected 1 network interface, got %d", len(info.Network))
	}

	if info.Environment["PATH"] == "" {
		t.Error("Expected PATH in environment")
	}
}

// Test SystemInfo marshaling preserves all fields
func TestSystemInfo_MarshalPreservesFields(t *testing.T) {
	info := &SystemInfo{
		Hostname:        "marshal-test",
		Platform:        "linux",
		PlatformFamily:  "debian",
		PlatformVersion: "11",
		Architecture:    "amd64",
		CPUs:            4,
		Uptime:          86400,
		LoadAverage:     []float64{1.5, 1.2, 1.0},
		Kernel:          "Linux",
		KernelVersion:   "5.10.0",
		Virtualization:  "kvm",
		Timezone:        "UTC",
		BootTime:        1234567890,
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var restored SystemInfo
	err = json.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if restored.Hostname != info.Hostname {
		t.Error("Hostname not preserved")
	}

	if len(restored.LoadAverage) != len(info.LoadAverage) {
		t.Error("LoadAverage not preserved")
	}
}

// Test edge cases
func TestSystemInfo_EdgeCases_NilPointers(t *testing.T) {
	info := &SystemInfo{
		Hostname: "test",
		Memory:   nil,
		Packages: nil,
		Processes: nil,
	}

	jsonStr, err := info.ToJSON()
	if err != nil {
		t.Fatalf("Should handle nil pointers: %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON")
	}
}

func TestSystemInfo_EdgeCases_EmptySlices(t *testing.T) {
	info := &SystemInfo{
		Hostname:    "test",
		Disk:        []*DiskInfo{},
		Network:     []*NetworkInfo{},
		Services:    []ServiceInfo{},
		Users:       []UserInfo{},
		Mounts:      []MountInfo{},
		LoadAverage: []float64{},
	}

	if len(info.Disk) != 0 {
		t.Error("Expected empty disk slice")
	}

	jsonStr, err := info.ToJSON()
	if err != nil {
		t.Fatalf("Should handle empty slices: %v", err)
	}

	if !strings.Contains(jsonStr, "disk") {
		t.Error("JSON should contain disk field")
	}
}

func TestSystemInfo_EdgeCases_LargeValues(t *testing.T) {
	info := &SystemInfo{
		Hostname: "large-test",
		Memory: &MemoryInfo{
			Total: 1024 * 1024 * 1024 * 1024, // 1TB
		},
		CPUs:   128,
		Uptime: 365 * 24 * 3600, // 1 year in seconds
	}

	if info.Memory.Total < 1024*1024*1024*1024 {
		t.Error("Large memory value not preserved")
	}

	if info.CPUs != 128 {
		t.Error("Large CPU count not preserved")
	}
}

func TestGetTimezone_Format(t *testing.T) {
	tz := getTimezone()

	if tz == "" {
		t.Error("Expected non-empty timezone")
	}
}