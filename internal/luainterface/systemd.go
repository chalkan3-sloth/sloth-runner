package luainterface

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// SystemdModule provides systemd service management functionality
type SystemdModule struct{}

// NewSystemdModule creates a new SystemdModule
func NewSystemdModule() *SystemdModule {
	return &SystemdModule{}
}

// Loader returns the Lua loader for the systemd module
func (mod *SystemdModule) Loader(L *lua.LState) int {
	// Create systemd module table
	systemdTable := L.NewTable()
	
	// Service management functions
	L.SetField(systemdTable, "create_service", L.NewFunction(mod.createService))
	L.SetField(systemdTable, "start", L.NewFunction(mod.startService))
	L.SetField(systemdTable, "stop", L.NewFunction(mod.stopService))
	L.SetField(systemdTable, "restart", L.NewFunction(mod.restartService))
	L.SetField(systemdTable, "reload", L.NewFunction(mod.reloadService))
	L.SetField(systemdTable, "enable", L.NewFunction(mod.enableService))
	L.SetField(systemdTable, "disable", L.NewFunction(mod.disableService))
	L.SetField(systemdTable, "status", L.NewFunction(mod.statusService))
	L.SetField(systemdTable, "is_active", L.NewFunction(mod.isActiveService))
	L.SetField(systemdTable, "is_enabled", L.NewFunction(mod.isEnabledService))
	L.SetField(systemdTable, "daemon_reload", L.NewFunction(mod.daemonReload))
	L.SetField(systemdTable, "remove_service", L.NewFunction(mod.removeService))
	L.SetField(systemdTable, "list_services", L.NewFunction(mod.listServices))
	L.SetField(systemdTable, "show", L.NewFunction(mod.showService))
	
	L.Push(systemdTable)
	return 1
}

// createService creates a systemd service file
// Usage: systemd.create_service({name="myapp", description="My App", exec_start="/usr/bin/myapp", ...})
func (mod *SystemdModule) createService(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	serviceName := opts.RawGetString("name").String()
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	// Build systemd service content
	var serviceContent strings.Builder
	serviceContent.WriteString("[Unit]\n")
	
	// Unit section
	if desc := opts.RawGetString("description"); desc != lua.LNil && desc.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("Description=%s\n", desc.String()))
	}
	if after := opts.RawGetString("after"); after != lua.LNil && after.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("After=%s\n", after.String()))
	}
	if wants := opts.RawGetString("wants"); wants != lua.LNil && wants.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("Wants=%s\n", wants.String()))
	}
	if requires := opts.RawGetString("requires"); requires != lua.LNil && requires.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("Requires=%s\n", requires.String()))
	}
	
	// Service section
	serviceContent.WriteString("\n[Service]\n")
	
	if execStart := opts.RawGetString("exec_start"); execStart != lua.LNil && execStart.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("ExecStart=%s\n", execStart.String()))
	} else {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("exec_start is required"))
		return 2
	}
	
	if execStop := opts.RawGetString("exec_stop"); execStop != lua.LNil && execStop.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("ExecStop=%s\n", execStop.String()))
	}
	if execReload := opts.RawGetString("exec_reload"); execReload != lua.LNil && execReload.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("ExecReload=%s\n", execReload.String()))
	}
	
	if serviceType := opts.RawGetString("type"); serviceType != lua.LNil && serviceType.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("Type=%s\n", serviceType.String()))
	} else {
		serviceContent.WriteString("Type=simple\n")
	}
	
	if user := opts.RawGetString("user"); user != lua.LNil && user.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("User=%s\n", user.String()))
	}
	if group := opts.RawGetString("group"); group != lua.LNil && group.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("Group=%s\n", group.String()))
	}
	if workingDir := opts.RawGetString("working_directory"); workingDir != lua.LNil && workingDir.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("WorkingDirectory=%s\n", workingDir.String()))
	}
	if restart := opts.RawGetString("restart"); restart != lua.LNil && restart.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("Restart=%s\n", restart.String()))
	} else {
		serviceContent.WriteString("Restart=always\n")
	}
	if restartSec := opts.RawGetString("restart_sec"); restartSec != lua.LNil && restartSec.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("RestartSec=%s\n", restartSec.String()))
	}
	
	// Environment variables
	if env := opts.RawGetString("environment"); env != lua.LNil {
		if envTable, ok := env.(*lua.LTable); ok {
			envTable.ForEach(func(key, value lua.LValue) {
				serviceContent.WriteString(fmt.Sprintf("Environment=%s=%s\n", key.String(), value.String()))
			})
		}
	}
	
	// Install section
	serviceContent.WriteString("\n[Install]\n")
	if wantedBy := opts.RawGetString("wanted_by"); wantedBy != lua.LNil && wantedBy.String() != "" {
		serviceContent.WriteString(fmt.Sprintf("WantedBy=%s\n", wantedBy.String()))
	} else {
		serviceContent.WriteString("WantedBy=multi-user.target\n")
	}
	
	// Write service file
	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	err := os.WriteFile(serviceFile, []byte(serviceContent.String()), 0644)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("Failed to create service file: %v", err)))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("Service file created: %s", serviceFile)))
	return 2
}

// systemdCommand executes a systemctl command
func (mod *SystemdModule) systemdCommand(command, serviceName string) (string, error) {
	var cmd *exec.Cmd
	if serviceName != "" {
		cmd = exec.Command("systemctl", command, serviceName)
	} else {
		cmd = exec.Command("systemctl", command)
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("systemctl %s failed: %s", command, stderr.String())
	}
	
	return stdout.String(), nil
}

// startService starts a systemd service (with idempotency)
// Usage: systemd.start({name="nginx"})
func (mod *SystemdModule) startService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	// IDEMPOTENCY: Check if already active
	output, err := mod.systemdCommand("is-active", serviceName)
	if err == nil && strings.TrimSpace(output) == "active" {
		result := L.NewTable()
		result.RawSetString("changed", lua.LFalse)
		result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s is already active", serviceName)))
		L.Push(lua.LTrue)
		L.Push(result)
		return 2
	}
	
	output, err = mod.systemdCommand("start", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	result.RawSetString("changed", lua.LTrue)
	result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s started", serviceName)))
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// stopService stops a systemd service (with idempotency)
// Usage: systemd.stop({name="nginx"})
func (mod *SystemdModule) stopService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	// IDEMPOTENCY: Check if already inactive
	output, err := mod.systemdCommand("is-active", serviceName)
	if err != nil || strings.TrimSpace(output) != "active" {
		result := L.NewTable()
		result.RawSetString("changed", lua.LFalse)
		result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s is already inactive", serviceName)))
		L.Push(lua.LTrue)
		L.Push(result)
		return 2
	}
	
	output, err = mod.systemdCommand("stop", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	result.RawSetString("changed", lua.LTrue)
	result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s stopped", serviceName)))
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// restartService restarts a systemd service
// Usage: systemd.restart({name="nginx"})
func (mod *SystemdModule) restartService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	output, err := mod.systemdCommand("restart", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(output))
	return 2
}

// reloadService reloads a systemd service
// Usage: systemd.reload({name="nginx"})
func (mod *SystemdModule) reloadService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	output, err := mod.systemdCommand("reload", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(output))
	return 2
}

// enableService enables a systemd service (with idempotency)
// Usage: systemd.enable({name="nginx"})
func (mod *SystemdModule) enableService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	// IDEMPOTENCY: Check if already enabled
	output, err := mod.systemdCommand("is-enabled", serviceName)
	if err == nil && strings.TrimSpace(output) == "enabled" {
		result := L.NewTable()
		result.RawSetString("changed", lua.LFalse)
		result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s is already enabled", serviceName)))
		L.Push(lua.LTrue)
		L.Push(result)
		return 2
	}
	
	output, err = mod.systemdCommand("enable", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	result.RawSetString("changed", lua.LTrue)
	result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s enabled", serviceName)))
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// disableService disables a systemd service (with idempotency)
// Usage: systemd.disable({name="nginx"})
func (mod *SystemdModule) disableService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	// IDEMPOTENCY: Check if already disabled
	output, err := mod.systemdCommand("is-enabled", serviceName)
	if err != nil || strings.TrimSpace(output) == "disabled" {
		result := L.NewTable()
		result.RawSetString("changed", lua.LFalse)
		result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s is already disabled", serviceName)))
		L.Push(lua.LTrue)
		L.Push(result)
		return 2
	}
	
	output, err = mod.systemdCommand("disable", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	result.RawSetString("changed", lua.LTrue)
	result.RawSetString("message", lua.LString(fmt.Sprintf("Service %s disabled", serviceName)))
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// statusService gets the status of a systemd service
// Usage: systemd.status({name="nginx"})
func (mod *SystemdModule) statusService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LString(""))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	output, err := mod.systemdCommand("status", serviceName)
	// Note: status can return non-zero exit codes for inactive services
	// We still want to return the output in this case
	
	L.Push(lua.LString(output))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 2
}

// isActiveService checks if a service is active
// Usage: systemd.is_active({name="nginx"})
func (mod *SystemdModule) isActiveService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	output, err := mod.systemdCommand("is-active", serviceName)
	isActive := err == nil && strings.TrimSpace(output) == "active"
	
	L.Push(lua.LBool(isActive))
	L.Push(lua.LString(strings.TrimSpace(output)))
	return 2
}

// isEnabledService checks if a service is enabled
// Usage: systemd.is_enabled({name="nginx"})
func (mod *SystemdModule) isEnabledService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	output, err := mod.systemdCommand("is-enabled", serviceName)
	isEnabled := err == nil && strings.TrimSpace(output) == "enabled"
	
	L.Push(lua.LBool(isEnabled))
	L.Push(lua.LString(strings.TrimSpace(output)))
	return 2
}

// daemonReload reloads systemd daemon
func (mod *SystemdModule) daemonReload(L *lua.LState) int {
	output, err := mod.systemdCommand("daemon-reload", "")
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(output))
	return 2
}

// removeService removes a systemd service file
// Usage: systemd.remove_service({name="myapp"})
func (mod *SystemdModule) removeService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	// Stop and disable service first
	mod.systemdCommand("stop", serviceName)
	mod.systemdCommand("disable", serviceName)
	
	// Remove service file
	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	err := os.Remove(serviceFile)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("Failed to remove service file: %v", err)))
		return 2
	}
	
	// Reload daemon
	mod.systemdCommand("daemon-reload", "")
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("Service removed: %s", serviceName)))
	return 2
}

// listServices lists systemd services
func (mod *SystemdModule) listServices(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	
	args := []string{"list-units", "--type=service"}
	
	if state := opts.RawGetString("state"); state != lua.LNil {
		args = append(args, "--state="+state.String())
	}
	
	if noHeader := opts.RawGetString("no_header"); lua.LVAsBool(noHeader) {
		args = append(args, "--no-legend")
	}
	
	cmd := exec.Command("systemctl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list services: %s", stderr.String())))
		return 2
	}
	
	L.Push(lua.LString(stdout.String()))
	L.Push(lua.LNil)
	return 2
}

// showService shows detailed information about a service
// Usage: systemd.show({name="nginx"})
func (mod *SystemdModule) showService(L *lua.LState) int {
	opts := L.CheckTable(1)
	serviceName := opts.RawGetString("name").String()
	
	if serviceName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("name parameter is required"))
		return 2
	}
	
	output, err := mod.systemdCommand("show", serviceName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(output))
	L.Push(lua.LNil)
	return 2
}

// SystemdLoader is the global loader function
func SystemdLoader(L *lua.LState) int {
	return NewSystemdModule().Loader(L)
}