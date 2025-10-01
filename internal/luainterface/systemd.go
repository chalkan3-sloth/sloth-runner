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
func (mod *SystemdModule) createService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	serviceConfig := L.CheckTable(2)
	
	// Build systemd service content
	var serviceContent strings.Builder
	serviceContent.WriteString("[Unit]\n")
	
	// Unit section
	if desc := serviceConfig.RawGetString("description"); desc != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("Description=%s\n", desc.String()))
	}
	if after := serviceConfig.RawGetString("after"); after != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("After=%s\n", after.String()))
	}
	if wants := serviceConfig.RawGetString("wants"); wants != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("Wants=%s\n", wants.String()))
	}
	if requires := serviceConfig.RawGetString("requires"); requires != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("Requires=%s\n", requires.String()))
	}
	
	// Service section
	serviceContent.WriteString("\n[Service]\n")
	
	if execStart := serviceConfig.RawGetString("exec_start"); execStart != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("ExecStart=%s\n", execStart.String()))
	} else {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("exec_start is required"))
		return 2
	}
	
	if execStop := serviceConfig.RawGetString("exec_stop"); execStop != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("ExecStop=%s\n", execStop.String()))
	}
	if execReload := serviceConfig.RawGetString("exec_reload"); execReload != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("ExecReload=%s\n", execReload.String()))
	}
	
	if serviceType := serviceConfig.RawGetString("type"); serviceType != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("Type=%s\n", serviceType.String()))
	} else {
		serviceContent.WriteString("Type=simple\n")
	}
	
	if user := serviceConfig.RawGetString("user"); user != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("User=%s\n", user.String()))
	}
	if group := serviceConfig.RawGetString("group"); group != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("Group=%s\n", group.String()))
	}
	if workingDir := serviceConfig.RawGetString("working_directory"); workingDir != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("WorkingDirectory=%s\n", workingDir.String()))
	}
	if restart := serviceConfig.RawGetString("restart"); restart != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("Restart=%s\n", restart.String()))
	} else {
		serviceContent.WriteString("Restart=always\n")
	}
	if restartSec := serviceConfig.RawGetString("restart_sec"); restartSec != lua.LNil {
		serviceContent.WriteString(fmt.Sprintf("RestartSec=%s\n", restartSec.String()))
	}
	
	// Environment variables
	if env := serviceConfig.RawGetString("environment"); env != lua.LNil {
		if envTable, ok := env.(*lua.LTable); ok {
			envTable.ForEach(func(key, value lua.LValue) {
				serviceContent.WriteString(fmt.Sprintf("Environment=%s=%s\n", key.String(), value.String()))
			})
		}
	}
	
	// Install section
	serviceContent.WriteString("\n[Install]\n")
	if wantedBy := serviceConfig.RawGetString("wanted_by"); wantedBy != lua.LNil {
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

// startService starts a systemd service
func (mod *SystemdModule) startService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
	output, err := mod.systemdCommand("start", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(output))
	return 2
}

// stopService stops a systemd service
func (mod *SystemdModule) stopService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
	output, err := mod.systemdCommand("stop", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(output))
	return 2
}

// restartService restarts a systemd service
func (mod *SystemdModule) restartService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
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
func (mod *SystemdModule) reloadService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
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

// enableService enables a systemd service
func (mod *SystemdModule) enableService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
	output, err := mod.systemdCommand("enable", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(output))
	return 2
}

// disableService disables a systemd service
func (mod *SystemdModule) disableService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
	output, err := mod.systemdCommand("disable", serviceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(output))
	return 2
}

// statusService gets the status of a systemd service
func (mod *SystemdModule) statusService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
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
func (mod *SystemdModule) isActiveService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
	output, err := mod.systemdCommand("is-active", serviceName)
	isActive := err == nil && strings.TrimSpace(output) == "active"
	
	L.Push(lua.LBool(isActive))
	L.Push(lua.LString(strings.TrimSpace(output)))
	return 2
}

// isEnabledService checks if a service is enabled
func (mod *SystemdModule) isEnabledService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
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
func (mod *SystemdModule) removeService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
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
func (mod *SystemdModule) showService(L *lua.LState) int {
	serviceName := L.CheckString(1)
	
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