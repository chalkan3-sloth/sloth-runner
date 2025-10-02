package luainterface

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// InfraTestModule provides infrastructure testing and validation capabilities
type InfraTestModule struct{}

// NewInfraTestModule creates a new InfraTestModule
func NewInfraTestModule() *InfraTestModule {
	return &InfraTestModule{}
}

// Loader returns the Lua loader for the infra_test module
func (m *InfraTestModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), m.exports())
	L.Push(mod)
	return 1
}

func (m *InfraTestModule) exports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		// File Tests
		"file_exists":   m.fileExists,
		"is_directory":  m.isDirectory,
		"is_file":       m.isFile,
		"file_contains": m.fileContains,
		"file_mode":     m.fileMode,
		"file_owner":    m.fileOwner,
		"file_group":    m.fileGroup,
		"file_size":     m.fileSize,

		// Network & Port Tests
		"port_is_listening": m.portIsListening,
		"port_is_tcp":       m.portIsTCP,
		"port_is_udp":       m.portIsUDP,
		"can_connect":       m.canConnect,
		"ping":              m.ping,

		// Service & Process Tests
		"service_is_running": m.serviceIsRunning,
		"service_is_enabled": m.serviceIsEnabled,
		"process_is_running": m.processIsRunning,
		"process_count":      m.processCount,

		// Command & Output Tests
		"command_succeeds":        m.commandSucceeds,
		"command_fails":           m.commandFails,
		"command_stdout_contains": m.commandStdoutContains,
		"command_stderr_is_empty": m.commandStderrIsEmpty,
		"command_output_equals":   m.commandOutputEquals,

		// Package Tests
		"package_is_installed": m.packageIsInstalled,
		"package_version":      m.packageVersion,
	}
}

// executeOnTarget executes a command on the specified target (agent)
func (m *InfraTestModule) executeOnTarget(L *lua.LState, target string, command string) (string, error) {
	if target == "" || target == "local" || target == "localhost" {
		// Execute locally
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()
		return string(output), err
	}

	// Execute on remote agent via delegate_to
	delegateFn := L.GetGlobal("delegate_to")
	if delegateFn.Type() == lua.LTNil {
		return "", fmt.Errorf("delegate_to function not available")
	}

	// Call delegate_to(target, function() ... end)
	err := L.CallByParam(lua.P{
		Fn:      delegateFn,
		NRet:    1,
		Protect: true,
	}, lua.LString(target), L.NewFunction(func(L *lua.LState) int {
		// Execute command
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LString(string(output)))
		return 1
	}))

	if err != nil {
		return "", err
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.Type() == lua.LTString {
		return result.String(), nil
	}

	return "", fmt.Errorf("unexpected return type from delegate_to")
}

// getOptionalTarget gets the optional target parameter
func (m *InfraTestModule) getOptionalTarget(L *lua.LState, argNum int) string {
	if L.GetTop() >= argNum {
		val := L.Get(argNum)
		if val.Type() == lua.LTString {
			return val.String()
		}
	}
	return "local"
}

// ============================================================================
// File Tests
// ============================================================================

// fileExists verifies if a file or directory exists
func (m *InfraTestModule) fileExists(L *lua.LState) int {
	path := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("test -e '%s' && echo 'true' || echo 'false'", path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("file_exists: execution failed: %v", err)
		return 0
	}

	exists := strings.TrimSpace(output) == "true"
	if !exists {
		L.RaiseError("file_exists: path '%s' does not exist on target '%s'", path, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// isDirectory verifies if the path is a directory
func (m *InfraTestModule) isDirectory(L *lua.LState) int {
	path := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("test -d '%s' && echo 'true' || echo 'false'", path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("is_directory: execution failed: %v", err)
		return 0
	}

	isDir := strings.TrimSpace(output) == "true"
	if !isDir {
		L.RaiseError("is_directory: path '%s' is not a directory on target '%s'", path, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// isFile verifies if the path is a regular file
func (m *InfraTestModule) isFile(L *lua.LState) int {
	path := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("test -f '%s' && echo 'true' || echo 'false'", path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("is_file: execution failed: %v", err)
		return 0
	}

	isFile := strings.TrimSpace(output) == "true"
	if !isFile {
		L.RaiseError("is_file: path '%s' is not a regular file on target '%s'", path, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// fileContains verifies if the file contains a string or regex pattern
func (m *InfraTestModule) fileContains(L *lua.LState) int {
	path := L.CheckString(1)
	pattern := L.CheckString(2)
	target := m.getOptionalTarget(L, 3)

	command := fmt.Sprintf("cat '%s'", path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("file_contains: failed to read file '%s': %v", path, err)
		return 0
	}

	// Try regex first, fallback to simple string contains
	matched, err := regexp.MatchString(pattern, output)
	if err != nil {
		// Fallback to simple string contains
		matched = strings.Contains(output, pattern)
	}

	if !matched {
		L.RaiseError("file_contains: file '%s' does not contain pattern '%s' on target '%s'", path, pattern, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// fileMode verifies the permissions of the file
func (m *InfraTestModule) fileMode(L *lua.LState) int {
	path := L.CheckString(1)
	expectedMode := L.CheckString(2)
	target := m.getOptionalTarget(L, 3)

	command := fmt.Sprintf("stat -c '%%a' '%s' 2>/dev/null || stat -f '%%A' '%s' 2>/dev/null", path, path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("file_mode: failed to get mode for '%s': %v", path, err)
		return 0
	}

	actualMode := strings.TrimSpace(output)
	// Normalize expected mode (remove 0o prefix if present)
	expectedMode = strings.TrimPrefix(expectedMode, "0o")
	expectedMode = strings.TrimPrefix(expectedMode, "0")

	if actualMode != expectedMode {
		L.RaiseError("file_mode: file '%s' has mode '%s' but expected '%s' on target '%s'",
			path, actualMode, expectedMode, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// fileOwner verifies the owner (user) of the file
func (m *InfraTestModule) fileOwner(L *lua.LState) int {
	path := L.CheckString(1)
	expectedUser := L.CheckString(2)
	target := m.getOptionalTarget(L, 3)

	command := fmt.Sprintf("stat -c '%%U' '%s' 2>/dev/null || stat -f '%%Su' '%s' 2>/dev/null", path, path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("file_owner: failed to get owner for '%s': %v", path, err)
		return 0
	}

	actualUser := strings.TrimSpace(output)
	if actualUser != expectedUser {
		L.RaiseError("file_owner: file '%s' is owned by '%s' but expected '%s' on target '%s'",
			path, actualUser, expectedUser, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// fileGroup verifies the group of the file
func (m *InfraTestModule) fileGroup(L *lua.LState) int {
	path := L.CheckString(1)
	expectedGroup := L.CheckString(2)
	target := m.getOptionalTarget(L, 3)

	command := fmt.Sprintf("stat -c '%%G' '%s' 2>/dev/null || stat -f '%%Sg' '%s' 2>/dev/null", path, path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("file_group: failed to get group for '%s': %v", path, err)
		return 0
	}

	actualGroup := strings.TrimSpace(output)
	if actualGroup != expectedGroup {
		L.RaiseError("file_group: file '%s' has group '%s' but expected '%s' on target '%s'",
			path, actualGroup, expectedGroup, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// fileSize verifies the size of the file in bytes
func (m *InfraTestModule) fileSize(L *lua.LState) int {
	path := L.CheckString(1)
	expectedSize := L.CheckInt(2)
	target := m.getOptionalTarget(L, 3)

	command := fmt.Sprintf("stat -c '%%s' '%s' 2>/dev/null || stat -f '%%z' '%s' 2>/dev/null", path, path)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("file_size: failed to get size for '%s': %v", path, err)
		return 0
	}

	actualSize, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		L.RaiseError("file_size: failed to parse size: %v", err)
		return 0
	}

	if actualSize != expectedSize {
		L.RaiseError("file_size: file '%s' has size %d but expected %d on target '%s'",
			path, actualSize, expectedSize, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// ============================================================================
// Network & Port Tests
// ============================================================================

// portIsListening verifies if a port is listening
func (m *InfraTestModule) portIsListening(L *lua.LState) int {
	port := L.CheckInt(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("netstat -ln 2>/dev/null | grep ':%d ' || ss -ln 2>/dev/null | grep ':%d '", port, port)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil || strings.TrimSpace(output) == "" {
		L.RaiseError("port_is_listening: port %d is not listening on target '%s'", port, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// portIsTCP verifies if a port is listening on TCP
func (m *InfraTestModule) portIsTCP(L *lua.LState) int {
	port := L.CheckInt(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("netstat -lnt 2>/dev/null | grep ':%d ' || ss -lnt 2>/dev/null | grep ':%d '", port, port)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil || strings.TrimSpace(output) == "" {
		L.RaiseError("port_is_tcp: port %d is not listening on TCP on target '%s'", port, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// portIsUDP verifies if a port is listening on UDP
func (m *InfraTestModule) portIsUDP(L *lua.LState) int {
	port := L.CheckInt(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("netstat -lnu 2>/dev/null | grep ':%d ' || ss -lnu 2>/dev/null | grep ':%d '", port, port)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil || strings.TrimSpace(output) == "" {
		L.RaiseError("port_is_udp: port %d is not listening on UDP on target '%s'", port, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// canConnect tests TCP connectivity to a host:port
func (m *InfraTestModule) canConnect(L *lua.LState) int {
	host := L.CheckString(1)
	port := L.CheckInt(2)
	timeoutMs := 5000
	if L.GetTop() >= 3 {
		timeoutMs = L.CheckInt(3)
	}

	timeoutSec := timeoutMs / 1000
	command := fmt.Sprintf("timeout %d bash -c 'cat < /dev/null > /dev/tcp/%s/%d' 2>/dev/null && echo 'success' || echo 'failed'",
		timeoutSec, host, port)

	cmd := exec.Command("sh", "-c", command)
	output, _ := cmd.CombinedOutput()

	if !strings.Contains(string(output), "success") {
		L.RaiseError("can_connect: failed to connect to %s:%d within %dms", host, port, timeoutMs)
	}

	L.Push(lua.LBool(true))
	return 1
}

// ping tests ICMP connectivity to a host
func (m *InfraTestModule) ping(L *lua.LState) int {
	host := L.CheckString(1)
	count := 4
	if L.GetTop() >= 2 {
		count = L.CheckInt(2)
	}
	target := m.getOptionalTarget(L, 3)

	command := fmt.Sprintf("ping -c %d %s > /dev/null 2>&1 && echo 'success' || echo 'failed'", count, host)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil || !strings.Contains(output, "success") {
		L.RaiseError("ping: failed to ping %s with %d packets on target '%s'", host, count, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// ============================================================================
// Service & Process Tests
// ============================================================================

// serviceIsRunning verifies if a service is active/running
func (m *InfraTestModule) serviceIsRunning(L *lua.LState) int {
	serviceName := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("systemctl is-active %s 2>/dev/null || service %s status 2>/dev/null | grep -q running",
		serviceName, serviceName)
	output, err := m.executeOnTarget(L, target, command)

	isActive := strings.Contains(strings.TrimSpace(output), "active") || err == nil
	if !isActive {
		L.RaiseError("service_is_running: service '%s' is not running on target '%s'", serviceName, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// serviceIsEnabled verifies if a service is enabled to start at boot
func (m *InfraTestModule) serviceIsEnabled(L *lua.LState) int {
	serviceName := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("systemctl is-enabled %s 2>/dev/null", serviceName)
	output, err := m.executeOnTarget(L, target, command)

	isEnabled := strings.Contains(strings.TrimSpace(output), "enabled") || err == nil
	if !isEnabled {
		L.RaiseError("service_is_enabled: service '%s' is not enabled on target '%s'", serviceName, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// processIsRunning verifies if a process matching the pattern is running
func (m *InfraTestModule) processIsRunning(L *lua.LState) int {
	pattern := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	command := fmt.Sprintf("pgrep -f '%s' > /dev/null && echo 'running' || echo 'not running'", pattern)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil || !strings.Contains(output, "running") {
		L.RaiseError("process_is_running: no process matching '%s' found on target '%s'", pattern, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// processCount verifies the number of processes matching a pattern
func (m *InfraTestModule) processCount(L *lua.LState) int {
	pattern := L.CheckString(1)
	expectedCount := L.CheckInt(2)
	target := m.getOptionalTarget(L, 3)

	command := fmt.Sprintf("pgrep -f '%s' | wc -l", pattern)
	output, err := m.executeOnTarget(L, target, command)

	if err != nil {
		L.RaiseError("process_count: failed to count processes for '%s': %v", pattern, err)
		return 0
	}

	actualCount, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		L.RaiseError("process_count: failed to parse count: %v", err)
		return 0
	}

	if actualCount != expectedCount {
		L.RaiseError("process_count: found %d processes matching '%s' but expected %d on target '%s'",
			actualCount, pattern, expectedCount, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// ============================================================================
// Command & Output Tests
// ============================================================================

// commandSucceeds verifies that a command returns exit code 0
func (m *InfraTestModule) commandSucceeds(L *lua.LState) int {
	command := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	_, err := m.executeOnTarget(L, target, command)
	if err != nil {
		L.RaiseError("command_succeeds: command '%s' failed on target '%s': %v", command, target, err)
	}

	L.Push(lua.LBool(true))
	return 1
}

// commandFails verifies that a command returns non-zero exit code
func (m *InfraTestModule) commandFails(L *lua.LState) int {
	command := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	_, err := m.executeOnTarget(L, target, command)
	if err == nil {
		L.RaiseError("command_fails: command '%s' succeeded but was expected to fail on target '%s'", command, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// commandStdoutContains verifies that command stdout contains a pattern
func (m *InfraTestModule) commandStdoutContains(L *lua.LState) int {
	command := L.CheckString(1)
	pattern := L.CheckString(2)
	target := m.getOptionalTarget(L, 3)

	output, err := m.executeOnTarget(L, target, command)
	if err != nil {
		L.RaiseError("command_stdout_contains: command '%s' failed: %v", command, err)
		return 0
	}

	matched, err := regexp.MatchString(pattern, output)
	if err != nil {
		matched = strings.Contains(output, pattern)
	}

	if !matched {
		L.RaiseError("command_stdout_contains: output of '%s' does not contain pattern '%s' on target '%s'",
			command, pattern, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// commandStderrIsEmpty verifies that command stderr is empty
func (m *InfraTestModule) commandStderrIsEmpty(L *lua.LState) int {
	command := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	fullCommand := fmt.Sprintf("%s 2>&1 1>/dev/null", command)
	output, _ := m.executeOnTarget(L, target, fullCommand)

	if strings.TrimSpace(output) != "" {
		L.RaiseError("command_stderr_is_empty: command '%s' produced stderr output on target '%s': %s",
			command, target, output)
	}

	L.Push(lua.LBool(true))
	return 1
}

// commandOutputEquals verifies that command stdout equals expected output
func (m *InfraTestModule) commandOutputEquals(L *lua.LState) int {
	command := L.CheckString(1)
	expectedOutput := L.CheckString(2)
	target := m.getOptionalTarget(L, 3)

	output, err := m.executeOnTarget(L, target, command)
	if err != nil {
		L.RaiseError("command_output_equals: command '%s' failed: %v", command, err)
		return 0
	}

	actualOutput := strings.TrimSpace(output)
	expectedOutput = strings.TrimSpace(expectedOutput)

	if actualOutput != expectedOutput {
		L.RaiseError("command_output_equals: output of '%s' is '%s' but expected '%s' on target '%s'",
			command, actualOutput, expectedOutput, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// ============================================================================
// Package Tests
// ============================================================================

// detectPackageManager detects the package manager available on the target
func (m *InfraTestModule) detectPackageManager(L *lua.LState, target string) (string, error) {
	managers := []struct {
		name    string
		command string
	}{
		{"dpkg", "which dpkg"},
		{"rpm", "which rpm"},
		{"pacman", "which pacman"},
		{"apk", "which apk"},
		{"brew", "which brew"},
	}

	for _, mgr := range managers {
		output, err := m.executeOnTarget(L, target, mgr.command)
		if err == nil && strings.TrimSpace(output) != "" {
			return mgr.name, nil
		}
	}

	return "", fmt.Errorf("no supported package manager found")
}

// packageIsInstalled verifies if a package is installed
func (m *InfraTestModule) packageIsInstalled(L *lua.LState) int {
	packageName := L.CheckString(1)
	target := m.getOptionalTarget(L, 2)

	pkgManager, err := m.detectPackageManager(L, target)
	if err != nil {
		L.RaiseError("package_is_installed: %v on target '%s'", err, target)
		return 0
	}

	var command string
	switch pkgManager {
	case "dpkg":
		command = fmt.Sprintf("dpkg -l '%s' 2>/dev/null | grep -q '^ii' && echo 'installed' || echo 'not installed'", packageName)
	case "rpm":
		command = fmt.Sprintf("rpm -q '%s' > /dev/null 2>&1 && echo 'installed' || echo 'not installed'", packageName)
	case "pacman":
		command = fmt.Sprintf("pacman -Q '%s' > /dev/null 2>&1 && echo 'installed' || echo 'not installed'", packageName)
	case "apk":
		command = fmt.Sprintf("apk info -e '%s' > /dev/null 2>&1 && echo 'installed' || echo 'not installed'", packageName)
	case "brew":
		command = fmt.Sprintf("brew list '%s' > /dev/null 2>&1 && echo 'installed' || echo 'not installed'", packageName)
	default:
		L.RaiseError("package_is_installed: unsupported package manager '%s'", pkgManager)
		return 0
	}

	output, err := m.executeOnTarget(L, target, command)
	if err != nil || !strings.Contains(output, "installed") {
		L.RaiseError("package_is_installed: package '%s' is not installed on target '%s'", packageName, target)
	}

	L.Push(lua.LBool(true))
	return 1
}

// packageVersion verifies the version of an installed package
func (m *InfraTestModule) packageVersion(L *lua.LState) int {
	packageName := L.CheckString(1)
	expectedVersion := L.CheckString(2)
	target := m.getOptionalTarget(L, 3)

	pkgManager, err := m.detectPackageManager(L, target)
	if err != nil {
		L.RaiseError("package_version: %v on target '%s'", err, target)
		return 0
	}

	var command string
	switch pkgManager {
	case "dpkg":
		command = fmt.Sprintf("dpkg-query -W -f='${Version}' '%s' 2>/dev/null", packageName)
	case "rpm":
		command = fmt.Sprintf("rpm -q --queryformat '%%{VERSION}-%%{RELEASE}' '%s' 2>/dev/null", packageName)
	case "pacman":
		command = fmt.Sprintf("pacman -Q '%s' 2>/dev/null | awk '{print $2}'", packageName)
	case "apk":
		command = fmt.Sprintf("apk info '%s' 2>/dev/null | grep '%s-' | sed 's/%s-//'", packageName, packageName, packageName)
	case "brew":
		command = fmt.Sprintf("brew list --versions '%s' 2>/dev/null | awk '{print $2}'", packageName)
	default:
		L.RaiseError("package_version: unsupported package manager '%s'", pkgManager)
		return 0
	}

	output, err := m.executeOnTarget(L, target, command)
	if err != nil {
		L.RaiseError("package_version: failed to get version for package '%s': %v", packageName, err)
		return 0
	}

	actualVersion := strings.TrimSpace(output)
	if actualVersion == "" {
		L.RaiseError("package_version: package '%s' is not installed on target '%s'", packageName, target)
		return 0
	}

	// Check if versions match (can be exact or prefix match)
	if !strings.HasPrefix(actualVersion, expectedVersion) && actualVersion != expectedVersion {
		L.RaiseError("package_version: package '%s' has version '%s' but expected '%s' on target '%s'",
			packageName, actualVersion, expectedVersion, target)
	}

	L.Push(lua.LBool(true))
	return 1
}
