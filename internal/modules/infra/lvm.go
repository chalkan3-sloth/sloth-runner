package infra

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// LVMModule provides LVM (Logical Volume Manager) management with fluent API
type LVMModule struct {
	L *lua.LState
}

// NewLVMModule creates a new LVM module instance
func NewLVMModule(L *lua.LState) *LVMModule {
	return &LVMModule{L: L}
}

// Register registers the LVM module with the Lua state
func (m *LVMModule) Register(L *lua.LState) {
	// Create lvm namespace table
	lvmTable := L.NewTable()

	// Physical Volume (PV) operations
	L.SetField(lvmTable, "pv_create", L.NewFunction(m.pvCreate))
	L.SetField(lvmTable, "pv_remove", L.NewFunction(m.pvRemove))
	L.SetField(lvmTable, "pv_list", L.NewFunction(m.pvList))
	L.SetField(lvmTable, "pv_exists", L.NewFunction(m.pvExists))

	// Volume Group (VG) operations
	L.SetField(lvmTable, "vg_create", L.NewFunction(m.vgCreate))
	L.SetField(lvmTable, "vg_remove", L.NewFunction(m.vgRemove))
	L.SetField(lvmTable, "vg_extend", L.NewFunction(m.vgExtend))
	L.SetField(lvmTable, "vg_reduce", L.NewFunction(m.vgReduce))
	L.SetField(lvmTable, "vg_list", L.NewFunction(m.vgList))
	L.SetField(lvmTable, "vg_exists", L.NewFunction(m.vgExists))

	// Logical Volume (LV) operations
	L.SetField(lvmTable, "lv_create", L.NewFunction(m.lvCreate))
	L.SetField(lvmTable, "lv_remove", L.NewFunction(m.lvRemove))
	L.SetField(lvmTable, "lv_extend", L.NewFunction(m.lvExtend))
	L.SetField(lvmTable, "lv_reduce", L.NewFunction(m.lvReduce))
	L.SetField(lvmTable, "lv_list", L.NewFunction(m.lvList))
	L.SetField(lvmTable, "lv_exists", L.NewFunction(m.lvExists))
	L.SetField(lvmTable, "lv_resize", L.NewFunction(m.lvResize))

	// Snapshot operations
	L.SetField(lvmTable, "snapshot_create", L.NewFunction(m.snapshotCreate))
	L.SetField(lvmTable, "snapshot_merge", L.NewFunction(m.snapshotMerge))

	// Thin provisioning
	L.SetField(lvmTable, "thin_pool_create", L.NewFunction(m.thinPoolCreate))
	L.SetField(lvmTable, "thin_lv_create", L.NewFunction(m.thinLvCreate))

	// Utility functions
	L.SetField(lvmTable, "scan", L.NewFunction(m.scan))
	L.SetField(lvmTable, "info", L.NewFunction(m.info))

	L.SetGlobal("lvm", lvmTable)
}

// execCommand executes a command and returns output
func (m *LVMModule) execCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// needsSudo checks if we need sudo
func (m *LVMModule) needsSudo() bool {
	output, err := m.execCommand("id", "-u")
	if err != nil {
		return true
	}
	return output != "0"
}

// prependSudo adds sudo to command if needed
func (m *LVMModule) prependSudo(args []string) []string {
	if m.needsSudo() {
		return append([]string{"sudo"}, args...)
	}
	return args
}

// =============================================================================
// Physical Volume (PV) Operations
// =============================================================================

// pvCreate creates a physical volume
func (m *LVMModule) pvCreate(L *lua.LState) int {
	device := L.CheckString(1)
	force := false
	if L.GetTop() >= 2 {
		force = L.ToBool(2)
	}

	// Check if PV already exists (idempotency)
	exists, _ := m.checkPVExists(device)
	if exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Physical volume %s already exists (idempotent)", device)))
		return 2
	}

	args := m.prependSudo([]string{"pvcreate"})
	if force {
		args = append(args, "-f")
	}
	args = append(args, device)

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create PV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Physical volume %s created successfully", device)))
	return 2
}

// pvRemove removes a physical volume
func (m *LVMModule) pvRemove(L *lua.LState) int {
	device := L.CheckString(1)
	force := false
	if L.GetTop() >= 2 {
		force = L.ToBool(2)
	}

	// Check if PV exists
	exists, _ := m.checkPVExists(device)
	if !exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Physical volume %s does not exist (idempotent)", device)))
		return 2
	}

	args := m.prependSudo([]string{"pvremove"})
	if force {
		args = append(args, "-f")
	}
	args = append(args, device)

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove PV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Physical volume %s removed successfully", device)))
	return 2
}

// pvList lists all physical volumes
func (m *LVMModule) pvList(L *lua.LState) int {
	args := m.prependSudo([]string{"pvs", "--noheadings", "-o", "pv_name,vg_name,pv_size,pv_free"})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list PVs: %s", err)))
		return 2
	}

	pvTable := L.NewTable()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			pvInfo := L.NewTable()
			pvInfo.RawSetString("name", lua.LString(fields[0]))
			pvInfo.RawSetString("vg", lua.LString(fields[1]))
			pvInfo.RawSetString("size", lua.LString(fields[2]))
			pvInfo.RawSetString("free", lua.LString(fields[3]))
			pvTable.Append(pvInfo)
		}
	}

	L.Push(pvTable)
	L.Push(lua.LNil)
	return 2
}

// pvExists checks if a physical volume exists
func (m *LVMModule) pvExists(L *lua.LState) int {
	device := L.CheckString(1)
	exists, _ := m.checkPVExists(device)
	L.Push(lua.LBool(exists))
	return 1
}

// checkPVExists helper to check PV existence
func (m *LVMModule) checkPVExists(device string) (bool, error) {
	args := m.prependSudo([]string{"pvs", device})
	_, err := m.execCommand(args...)
	return err == nil, err
}

// =============================================================================
// Volume Group (VG) Operations
// =============================================================================

// vgCreate creates a volume group
func (m *LVMModule) vgCreate(L *lua.LState) int {
	vgName := L.CheckString(1)
	devicesTable := L.CheckTable(2)

	// Check if VG already exists (idempotency)
	exists, _ := m.checkVGExists(vgName)
	if exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Volume group %s already exists (idempotent)", vgName)))
		return 2
	}

	// Extract devices from table
	var devices []string
	devicesTable.ForEach(func(_, v lua.LValue) {
		if v.Type() == lua.LTString {
			devices = append(devices, v.String())
		}
	})

	if len(devices) == 0 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No devices specified for volume group"))
		return 2
	}

	args := m.prependSudo(append([]string{"vgcreate", vgName}, devices...))
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create VG: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Volume group %s created successfully", vgName)))
	return 2
}

// vgRemove removes a volume group
func (m *LVMModule) vgRemove(L *lua.LState) int {
	vgName := L.CheckString(1)
	force := false
	if L.GetTop() >= 2 {
		force = L.ToBool(2)
	}

	// Check if VG exists
	exists, _ := m.checkVGExists(vgName)
	if !exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Volume group %s does not exist (idempotent)", vgName)))
		return 2
	}

	args := m.prependSudo([]string{"vgremove"})
	if force {
		args = append(args, "-f")
	}
	args = append(args, vgName)

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove VG: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Volume group %s removed successfully", vgName)))
	return 2
}

// vgExtend extends a volume group with new physical volumes
func (m *LVMModule) vgExtend(L *lua.LState) int {
	vgName := L.CheckString(1)
	devicesTable := L.CheckTable(2)

	// Extract devices from table
	var devices []string
	devicesTable.ForEach(func(_, v lua.LValue) {
		if v.Type() == lua.LTString {
			devices = append(devices, v.String())
		}
	})

	if len(devices) == 0 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No devices specified for volume group extension"))
		return 2
	}

	args := m.prependSudo(append([]string{"vgextend", vgName}, devices...))
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to extend VG: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Volume group %s extended successfully", vgName)))
	return 2
}

// vgReduce reduces a volume group by removing physical volumes
func (m *LVMModule) vgReduce(L *lua.LState) int {
	vgName := L.CheckString(1)
	devicesTable := L.CheckTable(2)

	// Extract devices from table
	var devices []string
	devicesTable.ForEach(func(_, v lua.LValue) {
		if v.Type() == lua.LTString {
			devices = append(devices, v.String())
		}
	})

	if len(devices) == 0 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No devices specified for volume group reduction"))
		return 2
	}

	args := m.prependSudo(append([]string{"vgreduce", vgName}, devices...))
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reduce VG: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Volume group %s reduced successfully", vgName)))
	return 2
}

// vgList lists all volume groups
func (m *LVMModule) vgList(L *lua.LState) int {
	args := m.prependSudo([]string{"vgs", "--noheadings", "-o", "vg_name,pv_count,lv_count,vg_size,vg_free"})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list VGs: %s", err)))
		return 2
	}

	vgTable := L.NewTable()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			vgInfo := L.NewTable()
			vgInfo.RawSetString("name", lua.LString(fields[0]))
			vgInfo.RawSetString("pv_count", lua.LString(fields[1]))
			vgInfo.RawSetString("lv_count", lua.LString(fields[2]))
			vgInfo.RawSetString("size", lua.LString(fields[3]))
			vgInfo.RawSetString("free", lua.LString(fields[4]))
			vgTable.Append(vgInfo)
		}
	}

	L.Push(vgTable)
	L.Push(lua.LNil)
	return 2
}

// vgExists checks if a volume group exists
func (m *LVMModule) vgExists(L *lua.LState) int {
	vgName := L.CheckString(1)
	exists, _ := m.checkVGExists(vgName)
	L.Push(lua.LBool(exists))
	return 1
}

// checkVGExists helper to check VG existence
func (m *LVMModule) checkVGExists(vgName string) (bool, error) {
	args := m.prependSudo([]string{"vgs", vgName})
	_, err := m.execCommand(args...)
	return err == nil, err
}

// =============================================================================
// Logical Volume (LV) Operations
// =============================================================================

// lvCreate creates a logical volume
func (m *LVMModule) lvCreate(L *lua.LState) int {
	vgName := L.CheckString(1)
	lvName := L.CheckString(2)
	size := L.CheckString(3) // e.g., "10G", "50%FREE"

	// Check if LV already exists (idempotency)
	fullName := fmt.Sprintf("%s/%s", vgName, lvName)
	exists, _ := m.checkLVExists(fullName)
	if exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Logical volume %s already exists (idempotent)", fullName)))
		return 2
	}

	args := m.prependSudo([]string{"lvcreate", "-n", lvName, "-L", size, vgName})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create LV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Logical volume %s created successfully", fullName)))
	return 2
}

// lvRemove removes a logical volume
func (m *LVMModule) lvRemove(L *lua.LState) int {
	lvPath := L.CheckString(1) // e.g., "vg_name/lv_name" or "/dev/vg_name/lv_name"
	force := false
	if L.GetTop() >= 2 {
		force = L.ToBool(2)
	}

	// Check if LV exists
	exists, _ := m.checkLVExists(lvPath)
	if !exists {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Logical volume %s does not exist (idempotent)", lvPath)))
		return 2
	}

	args := m.prependSudo([]string{"lvremove"})
	if force {
		args = append(args, "-f")
	}
	args = append(args, lvPath)

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove LV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Logical volume %s removed successfully", lvPath)))
	return 2
}

// lvExtend extends a logical volume
func (m *LVMModule) lvExtend(L *lua.LState) int {
	lvPath := L.CheckString(1)
	size := L.CheckString(2) // e.g., "+10G", "20G"

	args := m.prependSudo([]string{"lvextend", "-L", size, lvPath})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to extend LV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Logical volume %s extended successfully", lvPath)))
	return 2
}

// lvReduce reduces a logical volume
func (m *LVMModule) lvReduce(L *lua.LState) int {
	lvPath := L.CheckString(1)
	size := L.CheckString(2) // e.g., "-10G", "20G"
	force := false
	if L.GetTop() >= 3 {
		force = L.ToBool(3)
	}

	args := m.prependSudo([]string{"lvreduce", "-L", size})
	if force {
		args = append(args, "-f")
	}
	args = append(args, lvPath)

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reduce LV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Logical volume %s reduced successfully", lvPath)))
	return 2
}

// lvResize resizes a logical volume and filesystem together
func (m *LVMModule) lvResize(L *lua.LState) int {
	lvPath := L.CheckString(1)
	size := L.CheckString(2)
	resizeFs := true
	if L.GetTop() >= 3 {
		resizeFs = L.ToBool(3)
	}

	args := m.prependSudo([]string{"lvresize", "-L", size})
	if resizeFs {
		args = append(args, "-r") // Resize filesystem as well
	}
	args = append(args, lvPath)

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to resize LV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Logical volume %s resized successfully", lvPath)))
	return 2
}

// lvList lists all logical volumes
func (m *LVMModule) lvList(L *lua.LState) int {
	args := m.prependSudo([]string{"lvs", "--noheadings", "-o", "lv_name,vg_name,lv_size,lv_attr,pool_lv"})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list LVs: %s", err)))
		return 2
	}

	lvTable := L.NewTable()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			lvInfo := L.NewTable()
			lvInfo.RawSetString("name", lua.LString(fields[0]))
			lvInfo.RawSetString("vg", lua.LString(fields[1]))
			lvInfo.RawSetString("size", lua.LString(fields[2]))
			lvInfo.RawSetString("attr", lua.LString(fields[3]))
			if len(fields) >= 5 {
				lvInfo.RawSetString("pool", lua.LString(fields[4]))
			}
			lvTable.Append(lvInfo)
		}
	}

	L.Push(lvTable)
	L.Push(lua.LNil)
	return 2
}

// lvExists checks if a logical volume exists
func (m *LVMModule) lvExists(L *lua.LState) int {
	lvPath := L.CheckString(1)
	exists, _ := m.checkLVExists(lvPath)
	L.Push(lua.LBool(exists))
	return 1
}

// checkLVExists helper to check LV existence
func (m *LVMModule) checkLVExists(lvPath string) (bool, error) {
	args := m.prependSudo([]string{"lvs", lvPath})
	_, err := m.execCommand(args...)
	return err == nil, err
}

// =============================================================================
// Snapshot Operations
// =============================================================================

// snapshotCreate creates a snapshot of a logical volume
func (m *LVMModule) snapshotCreate(L *lua.LState) int {
	lvPath := L.CheckString(1)
	snapshotName := L.CheckString(2)
	size := L.CheckString(3) // Snapshot size, e.g., "1G"

	args := m.prependSudo([]string{"lvcreate", "-s", "-n", snapshotName, "-L", size, lvPath})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create snapshot: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Snapshot %s created successfully", snapshotName)))
	return 2
}

// snapshotMerge merges a snapshot back to its origin
func (m *LVMModule) snapshotMerge(L *lua.LState) int {
	snapshotPath := L.CheckString(1)

	args := m.prependSudo([]string{"lvconvert", "--merge", snapshotPath})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to merge snapshot: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Snapshot %s merged successfully", snapshotPath)))
	return 2
}

// =============================================================================
// Thin Provisioning
// =============================================================================

// thinPoolCreate creates a thin pool
func (m *LVMModule) thinPoolCreate(L *lua.LState) int {
	vgName := L.CheckString(1)
	poolName := L.CheckString(2)
	size := L.CheckString(3)

	args := m.prependSudo([]string{"lvcreate", "-T", "-L", size, "-n", poolName, vgName})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create thin pool: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Thin pool %s created successfully", poolName)))
	return 2
}

// thinLvCreate creates a thin logical volume
func (m *LVMModule) thinLvCreate(L *lua.LState) int {
	poolPath := L.CheckString(1) // e.g., "vg_name/pool_name"
	lvName := L.CheckString(2)
	virtualSize := L.CheckString(3)

	args := m.prependSudo([]string{"lvcreate", "-V", virtualSize, "-T", poolPath, "-n", lvName})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create thin LV: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Thin logical volume %s created successfully", lvName)))
	return 2
}

// =============================================================================
// Utility Functions
// =============================================================================

// scan scans for all LVM volumes
func (m *LVMModule) scan(L *lua.LState) int {
	args := m.prependSudo([]string{"pvscan"})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to scan: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(output))
	return 2
}

// info provides detailed information about a volume
func (m *LVMModule) info(L *lua.LState) int {
	volumePath := L.CheckString(1)

	// Try pvdisplay, vgdisplay, or lvdisplay based on path
	var cmdName string
	if matched, _ := regexp.MatchString(`^/dev/[^/]+$`, volumePath); matched {
		cmdName = "pvdisplay"
	} else if !strings.Contains(volumePath, "/") {
		cmdName = "vgdisplay"
	} else {
		cmdName = "lvdisplay"
	}

	args := m.prependSudo([]string{cmdName, volumePath})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get info: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LString(output))
	L.Push(lua.LNil)
	return 2
}
