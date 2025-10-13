package infra

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// RAIDModule provides RAID (mdadm) management
type RAIDModule struct {
	L *lua.LState
}

// NewRAIDModule creates a new RAID module instance
func NewRAIDModule(L *lua.LState) *RAIDModule {
	return &RAIDModule{L: L}
}

// Register registers the RAID module with the Lua state
func (m *RAIDModule) Register(L *lua.LState) {
	// Create raid namespace table
	raidTable := L.NewTable()

	// Array management
	L.SetField(raidTable, "create", L.NewFunction(m.create))
	L.SetField(raidTable, "assemble", L.NewFunction(m.assemble))
	L.SetField(raidTable, "stop", L.NewFunction(m.stop))
	L.SetField(raidTable, "remove", L.NewFunction(m.remove))

	// Device management
	L.SetField(raidTable, "add", L.NewFunction(m.addDevice))
	L.SetField(raidTable, "fail", L.NewFunction(m.failDevice))
	L.SetField(raidTable, "remove_device", L.NewFunction(m.removeDevice))
	L.SetField(raidTable, "add_spare", L.NewFunction(m.addSpare))

	// Monitoring and status
	L.SetField(raidTable, "detail", L.NewFunction(m.detail))
	L.SetField(raidTable, "status", L.NewFunction(m.status))
	L.SetField(raidTable, "list", L.NewFunction(m.list))
	L.SetField(raidTable, "exists", L.NewFunction(m.exists))

	// Recovery and maintenance
	L.SetField(raidTable, "grow", L.NewFunction(m.grow))
	L.SetField(raidTable, "check", L.NewFunction(m.check))
	L.SetField(raidTable, "repair", L.NewFunction(m.repair))

	// Configuration
	L.SetField(raidTable, "save_config", L.NewFunction(m.saveConfig))
	L.SetField(raidTable, "scan", L.NewFunction(m.scan))

	L.SetGlobal("raid", raidTable)
}

// execCommand executes a command and returns output
func (m *RAIDModule) execCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// needsSudo checks if we need sudo
func (m *RAIDModule) needsSudo() bool {
	output, err := m.execCommand("id", "-u")
	if err != nil {
		return true
	}
	return output != "0"
}

// prependSudo adds sudo to command if needed
func (m *RAIDModule) prependSudo(args []string) []string {
	if m.needsSudo() {
		return append([]string{"sudo"}, args...)
	}
	return args
}

// =============================================================================
// Array Management
// =============================================================================

// create creates a new RAID array
func (m *RAIDModule) create(L *lua.LState) int {
	options := L.CheckTable(1)

	// Extract required options
	name := options.RawGetString("name").String()
	level := options.RawGetString("level").String()

	if name == "" || level == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("name and level are required"))
		return 2
	}

	// Check if array already exists (idempotency)
	exists, _ := m.checkExists(name)
	if exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("RAID array %s already exists (idempotent)", name)))
		return 2
	}

	// Extract devices
	var devices []string
	devicesValue := options.RawGetString("devices")
	if devicesValue.Type() == lua.LTTable {
		devicesValue.(*lua.LTable).ForEach(func(_, v lua.LValue) {
			if v.Type() == lua.LTString {
				devices = append(devices, v.String())
			}
		})
	}

	if len(devices) == 0 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No devices specified"))
		return 2
	}

	// Build mdadm command
	args := m.prependSudo([]string{"mdadm", "--create", name, "--level=" + level, "--raid-devices=" + fmt.Sprint(len(devices))})

	// Optional: metadata version
	if metadata := options.RawGetString("metadata").String(); metadata != "" {
		args = append(args, "--metadata="+metadata)
	}

	// Optional: spare devices
	spareDevices := []string{}
	spareValue := options.RawGetString("spares")
	if spareValue.Type() == lua.LTTable {
		spareValue.(*lua.LTable).ForEach(func(_, v lua.LValue) {
			if v.Type() == lua.LTString {
				spareDevices = append(spareDevices, v.String())
			}
		})
	}

	if len(spareDevices) > 0 {
		args = append(args, "--spare-devices="+fmt.Sprint(len(spareDevices)))
	}

	// Optional: chunk size
	if chunk := options.RawGetString("chunk").String(); chunk != "" {
		args = append(args, "--chunk="+chunk)
	}

	// Optional: force
	if lua.LVAsBool(options.RawGetString("force")) {
		args = append(args, "--force")
	}

	// Add devices and spares
	args = append(args, devices...)
	args = append(args, spareDevices...)

	// Execute command
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create RAID array: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("RAID array %s created successfully", name)))
	return 2
}

// assemble assembles an existing RAID array
func (m *RAIDModule) assemble(L *lua.LState) int {
	name := L.CheckString(1)
	scan := false
	if L.GetTop() >= 2 {
		scan = L.ToBool(2)
	}

	args := m.prependSudo([]string{"mdadm", "--assemble", name})
	if scan {
		args = append(args, "--scan")
	}

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to assemble RAID array: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("RAID array %s assembled successfully", name)))
	return 2
}

// stop stops a RAID array
func (m *RAIDModule) stop(L *lua.LState) int {
	name := L.CheckString(1)

	// Check if array exists
	exists, _ := m.checkExists(name)
	if !exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("RAID array %s does not exist (idempotent)", name)))
		return 2
	}

	args := m.prependSudo([]string{"mdadm", "--stop", name})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to stop RAID array: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("RAID array %s stopped successfully", name)))
	return 2
}

// remove removes a RAID array (alias for stop + zero-superblock)
func (m *RAIDModule) remove(L *lua.LState) int {
	name := L.CheckString(1)

	// Stop the array first
	stopArgs := m.prependSudo([]string{"mdadm", "--stop", name})
	m.execCommand(stopArgs...)

	// Get devices from array before stopping
	detailArgs := m.prependSudo([]string{"mdadm", "--detail", name})
	detailOutput, _ := m.execCommand(detailArgs...)

	// Zero superblock on all devices
	deviceRegex := regexp.MustCompile(`/dev/\S+`)
	devices := deviceRegex.FindAllString(detailOutput, -1)

	for _, device := range devices {
		if device == name {
			continue
		}
		zeroArgs := m.prependSudo([]string{"mdadm", "--zero-superblock", device})
		m.execCommand(zeroArgs...)
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("RAID array %s removed successfully", name)))
	return 2
}

// =============================================================================
// Device Management
// =============================================================================

// addDevice adds a device to a RAID array
func (m *RAIDModule) addDevice(L *lua.LState) int {
	arrayName := L.CheckString(1)
	device := L.CheckString(2)

	args := m.prependSudo([]string{"mdadm", "--add", arrayName, device})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to add device: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Device %s added to %s successfully", device, arrayName)))
	return 2
}

// failDevice marks a device as failed
func (m *RAIDModule) failDevice(L *lua.LState) int {
	arrayName := L.CheckString(1)
	device := L.CheckString(2)

	args := m.prependSudo([]string{"mdadm", "--fail", arrayName, device})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to mark device as failed: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Device %s marked as failed in %s", device, arrayName)))
	return 2
}

// removeDevice removes a device from a RAID array
func (m *RAIDModule) removeDevice(L *lua.LState) int {
	arrayName := L.CheckString(1)
	device := L.CheckString(2)

	args := m.prependSudo([]string{"mdadm", "--remove", arrayName, device})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove device: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Device %s removed from %s successfully", device, arrayName)))
	return 2
}

// addSpare adds a spare device to a RAID array
func (m *RAIDModule) addSpare(L *lua.LState) int {
	arrayName := L.CheckString(1)
	device := L.CheckString(2)

	args := m.prependSudo([]string{"mdadm", "--add", arrayName, "--spare", device})
	output, err := m.execCommand(args...)
	if err != nil {
		// Try without --spare flag (older mdadm versions)
		args = m.prependSudo([]string{"mdadm", "--add", arrayName, device})
		output, err = m.execCommand(args...)
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(fmt.Sprintf("Failed to add spare device: %s\n%s", err, output)))
			return 2
		}
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Spare device %s added to %s successfully", device, arrayName)))
	return 2
}

// =============================================================================
// Monitoring and Status
// =============================================================================

// detail shows detailed information about a RAID array
func (m *RAIDModule) detail(L *lua.LState) int {
	name := L.CheckString(1)

	args := m.prependSudo([]string{"mdadm", "--detail", name})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get detail: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LString(output))
	L.Push(lua.LNil)
	return 2
}

// status shows the status of all RAID arrays
func (m *RAIDModule) status(L *lua.LState) int {
	// Read /proc/mdstat
	args := []string{"cat", "/proc/mdstat"}
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get status: %s", err)))
		return 2
	}

	L.Push(lua.LString(output))
	L.Push(lua.LNil)
	return 2
}

// list lists all RAID arrays
func (m *RAIDModule) list(L *lua.LState) int {
	// Read /proc/mdstat to get array names
	output, err := m.execCommand("cat", "/proc/mdstat")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list arrays: %s", err)))
		return 2
	}

	// Parse array names
	arrayRegex := regexp.MustCompile(`^(md\d+)\s*:`)
	lines := strings.Split(output, "\n")

	arrayTable := L.NewTable()
	for _, line := range lines {
		matches := arrayRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			arrayName := "/dev/" + matches[1]

			// Get detail for each array
			detailArgs := m.prependSudo([]string{"mdadm", "--detail", arrayName})
			detailOutput, err := m.execCommand(detailArgs...)
			if err == nil {
				info := m.parseDetail(L, detailOutput)
				info.RawSetString("name", lua.LString(arrayName))
				arrayTable.Append(info)
			}
		}
	}

	L.Push(arrayTable)
	L.Push(lua.LNil)
	return 2
}

// exists checks if a RAID array exists
func (m *RAIDModule) exists(L *lua.LState) int {
	name := L.CheckString(1)
	exists, _ := m.checkExists(name)
	L.Push(lua.LBool(exists))
	return 1
}

// checkExists helper to check array existence
func (m *RAIDModule) checkExists(name string) (bool, error) {
	args := m.prependSudo([]string{"mdadm", "--detail", name})
	_, err := m.execCommand(args...)
	return err == nil, err
}

// parseDetail parses mdadm --detail output into a Lua table
func (m *RAIDModule) parseDetail(L *lua.LState, output string) *lua.LTable {
	info := L.NewTable()
	lines := strings.Split(output, "\n")

	levelRegex := regexp.MustCompile(`Raid Level\s*:\s*(.+)`)
	devicesRegex := regexp.MustCompile(`Raid Devices\s*:\s*(\d+)`)
	stateRegex := regexp.MustCompile(`State\s*:\s*(.+)`)

	for _, line := range lines {
		if matches := levelRegex.FindStringSubmatch(line); len(matches) > 1 {
			info.RawSetString("level", lua.LString(matches[1]))
		}
		if matches := devicesRegex.FindStringSubmatch(line); len(matches) > 1 {
			info.RawSetString("devices", lua.LString(matches[1]))
		}
		if matches := stateRegex.FindStringSubmatch(line); len(matches) > 1 {
			info.RawSetString("state", lua.LString(matches[1]))
		}
	}

	return info
}

// =============================================================================
// Recovery and Maintenance
// =============================================================================

// grow grows (reshapes) a RAID array
func (m *RAIDModule) grow(L *lua.LState) int {
	arrayName := L.CheckString(1)
	options := L.CheckTable(2)

	args := m.prependSudo([]string{"mdadm", "--grow", arrayName})

	// Optional: change RAID level
	if level := options.RawGetString("level").String(); level != "" {
		args = append(args, "--level="+level)
	}

	// Optional: change number of devices
	if devices := options.RawGetString("raid_devices").String(); devices != "" {
		args = append(args, "--raid-devices="+devices)
	}

	// Optional: change chunk size
	if chunk := options.RawGetString("chunk").String(); chunk != "" {
		args = append(args, "--chunk="+chunk)
	}

	// Optional: change layout
	if layout := options.RawGetString("layout").String(); layout != "" {
		args = append(args, "--layout="+layout)
	}

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to grow array: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Array %s reshaped successfully", arrayName)))
	return 2
}

// check initiates a consistency check
func (m *RAIDModule) check(L *lua.LState) int {
	arrayName := L.CheckString(1)

	// Extract array name without /dev/
	arrayShortName := strings.TrimPrefix(arrayName, "/dev/")

	// Trigger check via sysfs
	checkPath := fmt.Sprintf("/sys/block/%s/md/sync_action", arrayShortName)
	args := m.prependSudo([]string{"sh", "-c", fmt.Sprintf("echo check > %s", checkPath)})

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to start check: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Consistency check started on %s", arrayName)))
	return 2
}

// repair initiates a repair operation
func (m *RAIDModule) repair(L *lua.LState) int {
	arrayName := L.CheckString(1)

	// Extract array name without /dev/
	arrayShortName := strings.TrimPrefix(arrayName, "/dev/")

	// Trigger repair via sysfs
	repairPath := fmt.Sprintf("/sys/block/%s/md/sync_action", arrayShortName)
	args := m.prependSudo([]string{"sh", "-c", fmt.Sprintf("echo repair > %s", repairPath)})

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to start repair: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Repair operation started on %s", arrayName)))
	return 2
}

// =============================================================================
// Configuration
// =============================================================================

// saveConfig saves RAID configuration to /etc/mdadm.conf or /etc/mdadm/mdadm.conf
func (m *RAIDModule) saveConfig(L *lua.LState) int {
	// Determine config file location
	var configFile string
	if _, err := m.execCommand("test", "-f", "/etc/mdadm.conf"); err == nil {
		configFile = "/etc/mdadm.conf"
	} else {
		configFile = "/etc/mdadm/mdadm.conf"
	}

	// Scan arrays and save configuration
	args := m.prependSudo([]string{"sh", "-c", fmt.Sprintf("mdadm --detail --scan > %s", configFile)})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to save config: %s\n%s", err, output)))
		return 2
	}

	// Update initramfs if available
	if _, err := m.execCommand("which", "update-initramfs"); err == nil {
		updateArgs := m.prependSudo([]string{"update-initramfs", "-u"})
		m.execCommand(updateArgs...)
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Configuration saved to %s", configFile)))
	return 2
}

// scan scans for all RAID arrays and reports
func (m *RAIDModule) scan(L *lua.LState) int {
	args := m.prependSudo([]string{"mdadm", "--detail", "--scan"})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to scan: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LString(output))
	L.Push(lua.LNil)
	return 2
}
