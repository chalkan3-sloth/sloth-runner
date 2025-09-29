package luainterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// EnhancedSaltModule provides advanced Salt functionality
type EnhancedSaltModule struct{}

// SaltClient represents an enhanced Salt client
type SaltClient struct {
	ConfigPath   string
	MasterHost   string
	MasterPort   int
	Username     string
	Password     string
	KeyFile      string
	Timeout      int
	RetryCount   int
	Environment  map[string]string
}

// SaltTarget represents an enhanced Salt target
type SaltTarget struct {
	Client     *SaltClient
	Target     string
	TargetType string
	Batch      string
	Async      bool
}

// SaltResult represents Salt command result
type SaltResult struct {
	Success   bool
	RetCode   int
	Stdout    string
	Stderr    string
	Duration  time.Duration
	Jid       string
	Returns   map[string]interface{}
}

// NewEnhancedSaltModule creates a new enhanced Salt module
func NewEnhancedSaltModule() *EnhancedSaltModule {
	return &EnhancedSaltModule{}
}

// RegisterEnhancedSaltModule registers the enhanced Salt module
func RegisterEnhancedSaltModule(L *lua.LState) {
	module := NewEnhancedSaltModule()
	
	// Create the salt table
	saltTable := L.NewTable()
	
	// Client management
	L.SetField(saltTable, "client", L.NewFunction(module.luaNewClient))
	L.SetField(saltTable, "master_client", L.NewFunction(module.luaNewMasterClient))
	L.SetField(saltTable, "local_client", L.NewFunction(module.luaNewLocalClient))
	
	// Quick access functions
	L.SetField(saltTable, "cmd", L.NewFunction(module.luaQuickCmd))
	L.SetField(saltTable, "state_apply", L.NewFunction(module.luaStateApply))
	L.SetField(saltTable, "grains", L.NewFunction(module.luaGrains))
	L.SetField(saltTable, "pillar", L.NewFunction(module.luaPillar))
	
	// Utility functions
	L.SetField(saltTable, "test_ping", L.NewFunction(module.luaTestPing))
	L.SetField(saltTable, "key_list", L.NewFunction(module.luaKeyList))
	L.SetField(saltTable, "key_accept", L.NewFunction(module.luaKeyAccept))
	L.SetField(saltTable, "highstate", L.NewFunction(module.luaHighstate))
	
	// Batch operations
	L.SetField(saltTable, "batch_cmd", L.NewFunction(module.luaBatchCmd))
	L.SetField(saltTable, "async_cmd", L.NewFunction(module.luaAsyncCmd))
	L.SetField(saltTable, "job_status", L.NewFunction(module.luaJobStatus))
	
	// Advanced operations
	L.SetField(saltTable, "orchestrate", L.NewFunction(module.luaOrchestrate))
	L.SetField(saltTable, "event_listen", L.NewFunction(module.luaEventListen))
	L.SetField(saltTable, "mine_get", L.NewFunction(module.luaMineGet))
	
	// Register the salt table globally
	L.SetGlobal("salt", saltTable)
}

// Client creation functions
func (s *EnhancedSaltModule) luaNewClient(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	client := &SaltClient{
		ConfigPath:  options.RawGetString("config").String(),
		MasterHost:  options.RawGetString("master").String(),
		MasterPort:  int(options.RawGetString("port").(lua.LNumber)),
		Username:    options.RawGetString("username").String(),
		Password:    options.RawGetString("password").String(),
		KeyFile:     options.RawGetString("key_file").String(),
		Timeout:     int(options.RawGetString("timeout").(lua.LNumber)),
		RetryCount:  int(options.RawGetString("retries").(lua.LNumber)),
		Environment: make(map[string]string),
	}
	
	// Set defaults
	if client.MasterPort == 0 {
		client.MasterPort = 4506
	}
	if client.Timeout == 0 {
		client.Timeout = 30
	}
	if client.RetryCount == 0 {
		client.RetryCount = 3
	}
	
	// Parse environment variables
	if envTable := options.RawGetString("env"); envTable.Type() == lua.LTTable {
		envTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			client.Environment[lua.LVAsString(key)] = lua.LVAsString(value)
		})
	}
	
	L.Push(s.clientToLua(L, client))
	return 1
}

func (s *EnhancedSaltModule) luaNewMasterClient(L *lua.LState) int {
	host := L.CheckString(1)
	port := L.OptInt(2, 4506)
	
	client := &SaltClient{
		MasterHost: host,
		MasterPort: port,
		Timeout:    30,
		RetryCount: 3,
	}
	
	L.Push(s.clientToLua(L, client))
	return 1
}

func (s *EnhancedSaltModule) luaNewLocalClient(L *lua.LState) int {
	configPath := L.OptString(1, "/etc/salt")
	
	client := &SaltClient{
		ConfigPath: configPath,
		Timeout:    30,
		RetryCount: 3,
	}
	
	L.Push(s.clientToLua(L, client))
	return 1
}

// Quick access functions
func (s *EnhancedSaltModule) luaQuickCmd(L *lua.LState) int {
	target := L.CheckString(1)
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	var args []string
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, lua.LVAsString(L.Get(i)))
	}
	
	cmd := fmt.Sprintf("%s.%s", module, function)
	result := s.executeSaltCmd(target, "glob", cmd, args, nil)
	
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaStateApply(L *lua.LState) int {
	target := L.CheckString(1)
	state := L.CheckString(2)
	
	options := L.OptTable(3, L.NewTable())
	testMode := lua.LVAsBool(options.RawGetString("test"))
	pillar := options.RawGetString("pillar")
	
	var args []string
	args = append(args, state)
	
	if testMode {
		args = append(args, "test=True")
	}
	
	if pillar.Type() == lua.LTTable {
		pillarData, _ := json.Marshal(s.luaTableToMap(pillar.(*lua.LTable)))
		args = append(args, fmt.Sprintf("pillar='%s'", string(pillarData)))
	}
	
	result := s.executeSaltCmd(target, "glob", "state.apply", args, nil)
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaGrains(L *lua.LState) int {
	target := L.CheckString(1)
	item := L.OptString(2, "")
	
	var cmd string
	var args []string
	
	if item != "" {
		cmd = "grains.item"
		args = []string{item}
	} else {
		cmd = "grains.items"
	}
	
	result := s.executeSaltCmd(target, "glob", cmd, args, nil)
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaPillar(L *lua.LState) int {
	target := L.CheckString(1)
	item := L.OptString(2, "")
	
	var cmd string
	var args []string
	
	if item != "" {
		cmd = "pillar.item"
		args = []string{item}
	} else {
		cmd = "pillar.items"
	}
	
	result := s.executeSaltCmd(target, "glob", cmd, args, nil)
	L.Push(s.resultToLua(L, result))
	return 1
}

// Utility functions
func (s *EnhancedSaltModule) luaTestPing(L *lua.LState) int {
	target := L.CheckString(1)
	
	result := s.executeSaltCmd(target, "glob", "test.ping", []string{}, nil)
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaKeyList(L *lua.LState) int {
	status := L.OptString(1, "all")
	
	result := s.executeSaltKeyCmd("--list", status)
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaKeyAccept(L *lua.LState) int {
	key := L.CheckString(1)
	
	result := s.executeSaltKeyCmd("--accept", key, "--yes")
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaHighstate(L *lua.LState) int {
	target := L.CheckString(1)
	
	options := L.OptTable(2, L.NewTable())
	testMode := lua.LVAsBool(options.RawGetString("test"))
	
	var args []string
	if testMode {
		args = append(args, "test=True")
	}
	
	result := s.executeSaltCmd(target, "glob", "state.highstate", args, nil)
	L.Push(s.resultToLua(L, result))
	return 1
}

// Batch operations
func (s *EnhancedSaltModule) luaBatchCmd(L *lua.LState) int {
	target := L.CheckString(1)
	batchSize := L.CheckString(2)
	module := L.CheckString(3)
	function := L.CheckString(4)
	
	var args []string
	for i := 5; i <= L.GetTop(); i++ {
		args = append(args, lua.LVAsString(L.Get(i)))
	}
	
	cmd := fmt.Sprintf("%s.%s", module, function)
	result := s.executeSaltCmd(target, "glob", cmd, args, &SaltClient{})
	
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaAsyncCmd(L *lua.LState) int {
	target := L.CheckString(1)
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	var args []string
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, lua.LVAsString(L.Get(i)))
	}
	
	cmd := fmt.Sprintf("%s.%s", module, function)
	
	// Execute with --async flag
	saltArgs := []string{"--async", "-L", target, cmd}
	saltArgs = append(saltArgs, args...)
	
	result := s.executeSaltCommand(saltArgs)
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaJobStatus(L *lua.LState) int {
	jid := L.CheckString(1)
	
	result := s.executeSaltRunCmd("jobs.lookup_jid", jid)
	L.Push(s.resultToLua(L, result))
	return 1
}

// Advanced operations
func (s *EnhancedSaltModule) luaOrchestrate(L *lua.LState) int {
	orch := L.CheckString(1)
	
	options := L.OptTable(2, L.NewTable())
	pillar := options.RawGetString("pillar")
	
	var args []string
	args = append(args, orch)
	
	if pillar.Type() == lua.LTTable {
		pillarData, _ := json.Marshal(s.luaTableToMap(pillar.(*lua.LTable)))
		args = append(args, fmt.Sprintf("pillar='%s'", string(pillarData)))
	}
	
	result := s.executeSaltRunCmd("state.orchestrate", args...)
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaEventListen(L *lua.LState) int {
	tag := L.OptString(1, "*")
	timeout := L.OptInt(2, 10)
	
	args := []string{"--timeout", fmt.Sprintf("%d", timeout)}
	if tag != "*" {
		args = append(args, tag)
	}
	
	result := s.executeSaltCommand(append([]string{"--events"}, args...))
	L.Push(s.resultToLua(L, result))
	return 1
}

func (s *EnhancedSaltModule) luaMineGet(L *lua.LState) int {
	target := L.CheckString(1)
	function := L.CheckString(2)
	
	result := s.executeSaltCmd(target, "glob", "mine.get", []string{target, function}, nil)
	L.Push(s.resultToLua(L, result))
	return 1
}

// Helper functions
func (s *EnhancedSaltModule) executeSaltCmd(target, targetType, cmd string, args []string, client *SaltClient) *SaltResult {
	saltArgs := []string{"--out=json", "-t", "30"}
	
	if client != nil && client.ConfigPath != "" {
		saltArgs = append(saltArgs, "--config-dir="+client.ConfigPath)
	}
	
	saltArgs = append(saltArgs, fmt.Sprintf("-%s", targetType), target, cmd)
	saltArgs = append(saltArgs, args...)
	
	return s.executeSaltCommand(saltArgs)
}

func (s *EnhancedSaltModule) executeSaltKeyCmd(args ...string) *SaltResult {
	saltKeyArgs := []string{"--out=json"}
	saltKeyArgs = append(saltKeyArgs, args...)
	
	startTime := time.Now()
	cmd := exec.Command("salt-key", saltKeyArgs...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	duration := time.Since(startTime)
	
	result := &SaltResult{
		Success:  err == nil,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.RetCode = exitErr.ExitCode()
		} else {
			result.RetCode = -1
		}
	}
	
	// Try to parse JSON output
	if result.Success && result.Stdout != "" {
		var returns map[string]interface{}
		if json.Unmarshal([]byte(result.Stdout), &returns) == nil {
			result.Returns = returns
		}
	}
	
	return result
}

func (s *EnhancedSaltModule) executeSaltRunCmd(cmd string, args ...string) *SaltResult {
	saltRunArgs := []string{"--out=json", cmd}
	saltRunArgs = append(saltRunArgs, args...)
	
	startTime := time.Now()
	execCmd := exec.Command("salt-run", saltRunArgs...)
	
	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr
	
	err := execCmd.Run()
	duration := time.Since(startTime)
	
	result := &SaltResult{
		Success:  err == nil,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.RetCode = exitErr.ExitCode()
		} else {
			result.RetCode = -1
		}
	}
	
	return result
}

func (s *EnhancedSaltModule) executeSaltCommand(args []string) *SaltResult {
	startTime := time.Now()
	cmd := exec.Command("salt", args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	duration := time.Since(startTime)
	
	result := &SaltResult{
		Success:  err == nil,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.RetCode = exitErr.ExitCode()
		} else {
			result.RetCode = -1
		}
	}
	
	// Try to parse JSON output
	if result.Success && result.Stdout != "" {
		var returns map[string]interface{}
		if json.Unmarshal([]byte(result.Stdout), &returns) == nil {
			result.Returns = returns
		}
	}
	
	return result
}

func (s *EnhancedSaltModule) clientToLua(L *lua.LState, client *SaltClient) lua.LValue {
	clientTable := L.NewTable()
	L.SetField(clientTable, "config_path", lua.LString(client.ConfigPath))
	L.SetField(clientTable, "master_host", lua.LString(client.MasterHost))
	L.SetField(clientTable, "master_port", lua.LNumber(client.MasterPort))
	L.SetField(clientTable, "timeout", lua.LNumber(client.Timeout))
	L.SetField(clientTable, "retry_count", lua.LNumber(client.RetryCount))
	
	return clientTable
}

func (s *EnhancedSaltModule) resultToLua(L *lua.LState, result *SaltResult) lua.LValue {
	resultTable := L.NewTable()
	L.SetField(resultTable, "success", lua.LBool(result.Success))
	L.SetField(resultTable, "ret_code", lua.LNumber(result.RetCode))
	L.SetField(resultTable, "stdout", lua.LString(result.Stdout))
	L.SetField(resultTable, "stderr", lua.LString(result.Stderr))
	L.SetField(resultTable, "duration_ms", lua.LNumber(result.Duration.Milliseconds()))
	L.SetField(resultTable, "jid", lua.LString(result.Jid))
	
	if result.Returns != nil {
		L.SetField(resultTable, "returns", s.mapToLuaTable(L, result.Returns))
	}
	
	return resultTable
}

func (s *EnhancedSaltModule) luaTableToMap(table *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	table.ForEach(func(key, value lua.LValue) {
		switch v := value.(type) {
		case lua.LString:
			result[lua.LVAsString(key)] = string(v)
		case lua.LNumber:
			result[lua.LVAsString(key)] = float64(v)
		case lua.LBool:
			result[lua.LVAsString(key)] = bool(v)
		case *lua.LTable:
			result[lua.LVAsString(key)] = s.luaTableToMap(v)
		default:
			result[lua.LVAsString(key)] = lua.LVAsString(value)
		}
	})
	return result
}

func (s *EnhancedSaltModule) mapToLuaTable(L *lua.LState, data map[string]interface{}) *lua.LTable {
	table := L.NewTable()
	for key, value := range data {
		switch v := value.(type) {
		case string:
			table.RawSetString(key, lua.LString(v))
		case float64:
			table.RawSetString(key, lua.LNumber(v))
		case bool:
			table.RawSetString(key, lua.LBool(v))
		case map[string]interface{}:
			table.RawSetString(key, s.mapToLuaTable(L, v))
		default:
			table.RawSetString(key, lua.LString(fmt.Sprintf("%v", v)))
		}
	}
	return table
}