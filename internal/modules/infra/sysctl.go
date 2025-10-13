package infra

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// SysctlModule provides kernel parameter management
type SysctlModule struct {
	L *lua.LState
}

// NewSysctlModule creates a new Sysctl module instance
func NewSysctlModule(L *lua.LState) *SysctlModule {
	return &SysctlModule{L: L}
}

// Register registers the Sysctl module with the Lua state
func (m *SysctlModule) Register(L *lua.LState) {
	sysctlTable := L.NewTable()

	L.SetField(sysctlTable, "get", L.NewFunction(m.get))
	L.SetField(sysctlTable, "set", L.NewFunction(m.set))
	L.SetField(sysctlTable, "set_persistent", L.NewFunction(m.setPersistent))
	L.SetField(sysctlTable, "list", L.NewFunction(m.list))
	L.SetField(sysctlTable, "reload", L.NewFunction(m.reload))
	L.SetField(sysctlTable, "exists", L.NewFunction(m.exists))
	L.SetField(sysctlTable, "apply", L.NewFunction(m.apply))

	L.SetGlobal("sysctl", sysctlTable)
}

func (m *SysctlModule) execCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func (m *SysctlModule) needsSudo() bool {
	output, err := m.execCommand("id", "-u")
	if err != nil {
		return true
	}
	return output != "0"
}

func (m *SysctlModule) prependSudo(args []string) []string {
	if m.needsSudo() {
		return append([]string{"sudo"}, args...)
	}
	return args
}

// get retrieves a kernel parameter value
func (m *SysctlModule) get(L *lua.LState) int {
	param := L.CheckString(1)

	args := []string{"sysctl", "-n", param}
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get parameter: %s", err)))
		return 2
	}

	L.Push(lua.LString(output))
	L.Push(lua.LNil)
	return 2
}

// set sets a kernel parameter (runtime only)
func (m *SysctlModule) set(L *lua.LState) int {
	param := L.CheckString(1)
	value := L.CheckAny(2)

	// Convert value to string
	valueStr := ""
	switch value.Type() {
	case lua.LTString:
		valueStr = value.String()
	case lua.LTNumber:
		valueStr = fmt.Sprintf("%v", float64(value.(lua.LNumber)))
	case lua.LTBool:
		if lua.LVAsBool(value) {
			valueStr = "1"
		} else {
			valueStr = "0"
		}
	default:
		L.Push(lua.LFalse)
		L.Push(lua.LString("Unsupported value type"))
		return 2
	}

	// Check current value (idempotency)
	currentArgs := []string{"sysctl", "-n", param}
	currentValue, _ := m.execCommand(currentArgs...)
	if currentValue == valueStr {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Parameter %s already set to %s (idempotent)", param, valueStr)))
		return 2
	}

	args := m.prependSudo([]string{"sysctl", "-w", fmt.Sprintf("%s=%s", param, valueStr)})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to set parameter: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Parameter %s set to %s", param, valueStr)))
	return 2
}

// setPersistent sets a kernel parameter persistently in /etc/sysctl.conf or /etc/sysctl.d/
func (m *SysctlModule) setPersistent(L *lua.LState) int {
	param := L.CheckString(1)
	value := L.CheckAny(2)
	configFile := "/etc/sysctl.d/99-sloth-runner.conf"
	if L.GetTop() >= 3 {
		configFile = L.CheckString(3)
	}

	// Convert value to string
	valueStr := ""
	switch value.Type() {
	case lua.LTString:
		valueStr = value.String()
	case lua.LTNumber:
		valueStr = fmt.Sprintf("%v", float64(value.(lua.LNumber)))
	case lua.LTBool:
		if lua.LVAsBool(value) {
			valueStr = "1"
		} else {
			valueStr = "0"
		}
	default:
		L.Push(lua.LFalse)
		L.Push(lua.LString("Unsupported value type"))
		return 2
	}

	// Read existing config
	var existingLines []string
	content, err := os.ReadFile(configFile)
	if err == nil {
		existingLines = strings.Split(string(content), "\n")
	}

	// Check if parameter already exists with correct value (idempotency)
	paramLine := fmt.Sprintf("%s = %s", param, valueStr)
	found := false
	updated := false

	for i, line := range existingLines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, param+"=") || strings.HasPrefix(trimmedLine, param+" =") {
			found = true
			if trimmedLine == paramLine || strings.TrimSpace(strings.Replace(trimmedLine, " ", "", -1)) == strings.TrimSpace(strings.Replace(paramLine, " ", "", -1)) {
				L.Push(lua.LTrue)
				L.Push(lua.LString(fmt.Sprintf("Parameter %s already persistent with value %s (idempotent)", param, valueStr)))
				return 2
			}
			existingLines[i] = paramLine
			updated = true
			break
		}
	}

	if !found {
		existingLines = append(existingLines, paramLine)
	}

	// Ensure directory exists
	dir := filepath.Dir(configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create config directory: %s", err)))
		return 2
	}

	// Write config file
	newContent := strings.Join(existingLines, "\n")
	writeArgs := m.prependSudo([]string{"sh", "-c", fmt.Sprintf("echo '%s' > %s", newContent, configFile)})
	output, err := m.execCommand(writeArgs...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to write config: %s\n%s", err, output)))
		return 2
	}

	// Apply the setting immediately
	applyArgs := m.prependSudo([]string{"sysctl", "-w", fmt.Sprintf("%s=%s", param, valueStr)})
	m.execCommand(applyArgs...)

	action := "added"
	if updated {
		action = "updated"
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Parameter %s %s in %s and applied", param, action, configFile)))
	return 2
}

// list lists all kernel parameters
func (m *SysctlModule) list(L *lua.LState) int {
	pattern := ""
	if L.GetTop() >= 1 {
		pattern = L.CheckString(1)
	}

	args := []string{"sysctl", "-a"}
	if pattern != "" {
		args = append(args, pattern)
	}

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list parameters: %s", err)))
		return 2
	}

	// Parse output into table
	paramTable := L.NewTable()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			paramInfo := L.NewTable()
			paramInfo.RawSetString("name", lua.LString(strings.TrimSpace(parts[0])))
			paramInfo.RawSetString("value", lua.LString(strings.TrimSpace(parts[1])))
			paramTable.Append(paramInfo)
		}
	}

	L.Push(paramTable)
	L.Push(lua.LNil)
	return 2
}

// reload reloads sysctl configuration from files
func (m *SysctlModule) reload(L *lua.LState) int {
	configFile := ""
	if L.GetTop() >= 1 {
		configFile = L.CheckString(1)
	}

	args := m.prependSudo([]string{"sysctl", "-p"})
	if configFile != "" {
		args = append(args, configFile)
	}

	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reload config: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString("Sysctl configuration reloaded successfully"))
	return 2
}

// exists checks if a kernel parameter exists
func (m *SysctlModule) exists(L *lua.LState) int {
	param := L.CheckString(1)

	args := []string{"sysctl", "-n", param}
	_, err := m.execCommand(args...)
	L.Push(lua.LBool(err == nil))
	return 1
}

// apply applies all sysctl settings from /etc/sysctl.conf and /etc/sysctl.d/
func (m *SysctlModule) apply(L *lua.LState) int {
	args := m.prependSudo([]string{"sysctl", "--system"})
	output, err := m.execCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to apply settings: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString("All sysctl settings applied successfully"))
	return 2
}
