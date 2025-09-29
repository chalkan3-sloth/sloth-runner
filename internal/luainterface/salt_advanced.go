package luainterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// SaltModule provides advanced SaltStack integration
type SaltModule struct{}

// NewSaltModule creates a new SaltModule
func NewSaltModule() *SaltModule {
	return &SaltModule{}
}

// Loader returns the Lua loader for the salt module
func (mod *SaltModule) Loader(L *lua.LState) int {
	saltTable := L.NewTable()
	L.SetFuncs(saltTable, map[string]lua.LGFunction{
		"client":       mod.saltClient,
		"execute":      mod.saltExecute,
		"state_apply":  mod.saltStateApply,
		"state_highstate": mod.saltStateHighstate,
		"pillar_get":   mod.saltPillarGet,
		"grains_get":   mod.saltGrainsGet,
		"test_ping":    mod.saltTestPing,
		"cp_get_file":  mod.saltCpGetFile,
		"service":      mod.saltService,
		"pkg":          mod.saltPkg,
		"file":         mod.saltFile,
		"user":         mod.saltUser,
		"group":        mod.saltGroup,
		"async_run":    mod.saltAsyncRun,
		"job_status":   mod.saltJobStatus,
	})
	L.Push(saltTable)
	return 1
}

// saltClient creates a new salt client object for method chaining
func (mod *SaltModule) saltClient(L *lua.LState) int {
	clientTable := L.NewTable()
	
	// Add target method that returns a table with cmd method
	targetFunc := L.NewFunction(func(L *lua.LState) int {
		target := L.CheckString(2)       // target pattern
		targetType := L.CheckString(3)   // "glob", "list", "grain", etc.
		
		targetTable := L.NewTable()
		targetTable.RawSetString("target", lua.LString(target))
		targetTable.RawSetString("type", lua.LString(targetType))
		
		// Add cmd method
		cmdFunc := L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckString(2) // command (not used in mock)
			
			// Mock response for testing
			L.Push(lua.LString(""))        // stdout
			L.Push(lua.LString("salt error")) // stderr  
			L.Push(lua.LString("error"))   // error
			return 3
		})
		targetTable.RawSetString("cmd", cmdFunc)
		
		L.Push(targetTable)
		return 1
	})
	clientTable.RawSetString("target", targetFunc)
	
	L.Push(clientTable)
	return 1
}

// saltExecute executes a raw salt command
func (mod *SaltModule) saltExecute(L *lua.LState) int {
	target := L.CheckString(1)
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	// Optional arguments
	args := ""
	if L.GetTop() > 3 {
		args = L.CheckString(4)
	}
	
	// Optional configuration
	opts := L.OptTable(5, L.NewTable())
	timeoutVal := opts.RawGetString("timeout")
	timeout := 60 // default
	if timeoutVal.Type() == lua.LTNumber {
		timeout = int(lua.LVAsNumber(timeoutVal))
	}
	
	saltCmd := []string{"salt", target, fmt.Sprintf("%s.%s", module, function)}
	if args != "" {
		saltCmd = append(saltCmd, args)
	}
	saltCmd = append(saltCmd, "--out=json", fmt.Sprintf("--timeout=%d", timeout))
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltStateApply applies a specific state
func (mod *SaltModule) saltStateApply(L *lua.LState) int {
	target := L.CheckString(1)
	stateName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	pillar := ""
	if pillarValue := opts.RawGetString("pillar"); pillarValue != lua.LNil {
		pillar = pillarValue.String()
	}
	
	saltCmd := []string{"salt", target, "state.apply", stateName, "--out=json"}
	if pillar != "" {
		saltCmd = append(saltCmd, fmt.Sprintf("pillar='%s'", pillar))
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltStateHighstate runs highstate on target
func (mod *SaltModule) saltStateHighstate(L *lua.LState) int {
	target := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	test := opts.RawGetString("test")
	
	saltCmd := []string{"salt", target, "state.highstate", "--out=json"}
	if lua.LVAsBool(test) {
		saltCmd = append(saltCmd, "test=True")
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltPillarGet retrieves pillar data
func (mod *SaltModule) saltPillarGet(L *lua.LState) int {
	target := L.CheckString(1)
	pillarKey := L.CheckString(2)
	
	saltCmd := []string{"salt", target, "pillar.get", pillarKey, "--out=json"}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltGrainsGet retrieves grains data
func (mod *SaltModule) saltGrainsGet(L *lua.LState) int {
	target := L.CheckString(1)
	grainKey := ""
	if L.GetTop() > 1 {
		grainKey = L.CheckString(2)
	}
	
	var saltCmd []string
	if grainKey != "" {
		saltCmd = []string{"salt", target, "grains.get", grainKey, "--out=json"}
	} else {
		saltCmd = []string{"salt", target, "grains.items", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltTestPing tests connectivity to minions
func (mod *SaltModule) saltTestPing(L *lua.LState) int {
	target := L.CheckString(1)
	
	saltCmd := []string{"salt", target, "test.ping", "--out=json"}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltCpGetFile copies files from master to minions
func (mod *SaltModule) saltCpGetFile(L *lua.LState) int {
	target := L.CheckString(1)
	sourcePath := L.CheckString(2)
	destPath := L.CheckString(3)
	
	saltCmd := []string{"salt", target, "cp.get_file", sourcePath, destPath, "--out=json"}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltService manages services
func (mod *SaltModule) saltService(L *lua.LState) int {
	target := L.CheckString(1)
	action := L.CheckString(2) // start, stop, restart, status, enable, disable
	serviceName := L.CheckString(3)
	
	saltCmd := []string{"salt", target, fmt.Sprintf("service.%s", action), serviceName, "--out=json"}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltPkg manages packages
func (mod *SaltModule) saltPkg(L *lua.LState) int {
	target := L.CheckString(1)
	action := L.CheckString(2) // install, remove, upgrade, list_pkgs
	
	var saltCmd []string
	switch action {
	case "install", "remove":
		packageName := L.CheckString(3)
		saltCmd = []string{"salt", target, fmt.Sprintf("pkg.%s", action), packageName, "--out=json"}
	case "upgrade":
		saltCmd = []string{"salt", target, "pkg.upgrade", "--out=json"}
	case "list_pkgs":
		saltCmd = []string{"salt", target, "pkg.list_pkgs", "--out=json"}
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid action. Use: install, remove, upgrade, list_pkgs"))
		return 2
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltFile manages files
func (mod *SaltModule) saltFile(L *lua.LState) int {
	target := L.CheckString(1)
	action := L.CheckString(2) // exists, copy, remove, set_mode
	filePath := L.CheckString(3)
	
	var saltCmd []string
	switch action {
	case "exists":
		saltCmd = []string{"salt", target, "file.file_exists", filePath, "--out=json"}
	case "copy":
		destPath := L.CheckString(4)
		saltCmd = []string{"salt", target, "file.copy", filePath, destPath, "--out=json"}
	case "remove":
		saltCmd = []string{"salt", target, "file.remove", filePath, "--out=json"}
	case "set_mode":
		mode := L.CheckString(4)
		saltCmd = []string{"salt", target, "file.set_mode", filePath, mode, "--out=json"}
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid action. Use: exists, copy, remove, set_mode"))
		return 2
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltUser manages users
func (mod *SaltModule) saltUser(L *lua.LState) int {
	target := L.CheckString(1)
	action := L.CheckString(2) // add, delete, info, list_users
	
	var saltCmd []string
	switch action {
	case "add":
		username := L.CheckString(3)
		opts := L.OptTable(4, L.NewTable())
		
		saltCmd = []string{"salt", target, "user.add", username}
		
		if home := opts.RawGetString("home"); home != lua.LNil {
			saltCmd = append(saltCmd, fmt.Sprintf("home=%s", home.String()))
		}
		if shell := opts.RawGetString("shell"); shell != lua.LNil {
			saltCmd = append(saltCmd, fmt.Sprintf("shell=%s", shell.String()))
		}
		
		saltCmd = append(saltCmd, "--out=json")
		
	case "delete":
		username := L.CheckString(3)
		saltCmd = []string{"salt", target, "user.delete", username, "--out=json"}
		
	case "info":
		username := L.CheckString(3)
		saltCmd = []string{"salt", target, "user.info", username, "--out=json"}
		
	case "list_users":
		saltCmd = []string{"salt", target, "user.list_users", "--out=json"}
		
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid action. Use: add, delete, info, list_users"))
		return 2
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltGroup manages groups
func (mod *SaltModule) saltGroup(L *lua.LState) int {
	target := L.CheckString(1)
	action := L.CheckString(2) // add, delete, info, list_groups
	
	var saltCmd []string
	switch action {
	case "add":
		groupname := L.CheckString(3)
		saltCmd = []string{"salt", target, "group.add", groupname, "--out=json"}
		
	case "delete":
		groupname := L.CheckString(3)
		saltCmd = []string{"salt", target, "group.delete", groupname, "--out=json"}
		
	case "info":
		groupname := L.CheckString(3)
		saltCmd = []string{"salt", target, "group.info", groupname, "--out=json"}
		
	case "list_groups":
		saltCmd = []string{"salt", target, "group.getent", "--out=json"}
		
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid action. Use: add, delete, info, list_groups"))
		return 2
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltAsyncRun runs commands asynchronously
func (mod *SaltModule) saltAsyncRun(L *lua.LState) int {
	target := L.CheckString(1)
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	// Optional arguments
	args := ""
	if L.GetTop() > 3 {
		args = L.CheckString(4)
	}
	
	saltCmd := []string{"salt", target, fmt.Sprintf("%s.%s", module, function), "--async", "--out=json"}
	if args != "" {
		saltCmd = append(saltCmd, args)
	}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Parse the job ID from the result
	var jobData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jobData); err == nil {
		if jid, ok := jobData["jid"]; ok {
			L.Push(lua.LString(jid.(string)))
			L.Push(lua.LNil)
			return 2
		}
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// saltJobStatus checks the status of an asynchronous job
func (mod *SaltModule) saltJobStatus(L *lua.LState) int {
	jobID := L.CheckString(1)
	
	saltCmd := []string{"salt-run", "jobs.lookup_jid", jobID, "--out=json"}
	
	result, err := mod.executeSaltCommand(saltCmd...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// executeSaltCommand executes a salt command and returns the result
func (mod *SaltModule) executeSaltCommand(cmdArgs ...string) (string, error) {
	// Check if salt command exists
	if _, err := exec.LookPath("salt"); err != nil {
		return "", fmt.Errorf("salt command not found in PATH: %w", err)
	}
	
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	
	// Set environment variables
	cmd.Env = os.Environ()
	
	// Set working directory to home if config files are there
	if homeDir, err := os.UserHomeDir(); err == nil {
		if _, err := os.Stat(filepath.Join(homeDir, ".saltrc")); err == nil {
			cmd.Dir = homeDir
		}
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Set timeout
	timeout := 120 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			errorMsg := stderr.String()
			if errorMsg == "" {
				errorMsg = err.Error()
			}
			return "", fmt.Errorf("salt command failed: %s", errorMsg)
		}
		return stdout.String(), nil
		
	case <-timer.C:
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("salt command timed out after %v", timeout)
	}
}

// SaltLoader loads the salt module for Lua
func SaltLoader(L *lua.LState) int {
	mod := NewSaltModule()
	return mod.Loader(L)
}