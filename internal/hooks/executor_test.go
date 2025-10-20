//go:build cgo
// +build cgo

package hooks

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// TestNewExecutor tests creating a new executor
func TestNewExecutor(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	executor := NewExecutor(repo)

	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}

	if executor.repo != repo {
		t.Error("Executor repository not set correctly")
	}
}

// TestExecute_FileNotFound tests executing a hook with non-existent file
func TestExecute_FileNotFound(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	hook := &Hook{
		ID:       "test-hook-1",
		Name:     "nonexistent",
		FilePath: "/nonexistent/path/to/hook.lua",
		EventType: EventTaskStarted,
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"name": "test-task",
			},
		},
	}

	result, err := executor.Execute(hook, event)

	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	if result.Success {
		t.Error("Expected result.Success to be false")
	}

	if result.Error == "" {
		t.Error("Expected error message in result")
	}
}

// TestExecute_SimpleScript tests executing a simple Lua hook script
func TestExecute_SimpleScript(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	// Create a temporary hook file
	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "simple_hook.lua")

	script := `
-- Simple hook that just returns true
function on_event()
	log.info("Hook executed")
	return true
end
`

	if err := os.WriteFile(hookFile, []byte(script), 0644); err != nil {
		t.Fatalf("Failed to create hook file: %v", err)
	}

	hook := &Hook{
		ID:       "simple-hook",
		Name:     "simple",
		FilePath: hookFile,
		EventType: EventTaskStarted,
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"name": "test-task",
			},
		},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got failure: %s", result.Error)
	}

	if result.Output == "" {
		t.Error("Expected output from hook")
	}

	if result.Duration == 0 {
		t.Error("Expected non-zero duration")
	}
}

// TestExecute_HookReturnsFalse tests a hook that returns false
func TestExecute_HookReturnsFalse(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "failing_hook.lua")

	script := `
function on_event()
	log.error("Hook failed")
	return false
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "failing-hook",
		Name:     "failing",
		FilePath: hookFile,
		EventType: EventTaskFailed,
	}

	event := &Event{
		Type:      EventTaskFailed,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v (should complete even if hook returns false)", err)
	}

	if result.Success {
		t.Error("Expected result.Success to be false")
	}
}

// TestExecute_NoOnEventFunction tests a script without on_event function
func TestExecute_NoOnEventFunction(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "no_function.lua")

	script := `
-- Script without on_event function
local x = 42
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "no-func-hook",
		Name:     "no-func",
		FilePath: hookFile,
		EventType: EventTaskStarted,
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Should still succeed if no on_event function
	if !result.Success {
		t.Error("Expected success even without on_event function")
	}
}

// TestExecute_ScriptError tests a script with syntax error
func TestExecute_ScriptError(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "error_hook.lua")

	script := `
-- Script with syntax error
function on_event()
	invalid lua syntax here!!!
	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "error-hook",
		Name:     "error",
		FilePath: hookFile,
		EventType: EventTaskStarted,
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	result, err := executor.Execute(hook, event)

	if err == nil {
		t.Error("Expected error for invalid Lua script")
	}

	if result.Success {
		t.Error("Expected result.Success to be false")
	}

	if result.Error == "" {
		t.Error("Expected error message")
	}
}

// TestExecute_AccessEventData tests accessing event data from Lua
func TestExecute_AccessEventData(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "event_data.lua")

	script := `
function on_event()
	log.info("Event type: " .. event.type)
	log.info("Task name: " .. event.task.name)
	log.info("Task status: " .. event.task.status)
	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "event-data-hook",
		Name:     "event-data",
		FilePath: hookFile,
		EventType: EventTaskCompleted,
	}

	event := &Event{
		Type:      EventTaskCompleted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"name":   "my-task",
				"status": "success",
			},
		},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got failure: %s", result.Error)
	}

	// Check output contains expected strings
	if result.Output == "" {
		t.Error("Expected output from hook")
	}
}

// TestExecute_AgentEvent tests executing hook with agent event data
func TestExecute_AgentEvent(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "agent_event.lua")

	script := `
function on_event()
	log.info("Agent: " .. event.agent.name)
	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "agent-hook",
		Name:     "agent",
		FilePath: hookFile,
		EventType: EventAgentRegistered,
	}

	event := &Event{
		Type:      EventAgentRegistered,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent": map[string]interface{}{
				"name": "test-agent",
				"host": "192.168.1.100",
			},
		},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}
}

// TestRegisterEvent tests event registration in Lua
func TestRegisterEvent(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"name":   "test-task",
				"status": "running",
			},
		},
	}

	err := executor.registerEvent(L, event)
	if err != nil {
		t.Errorf("registerEvent() error = %v", err)
	}

	// Verify event global exists
	eventTable := L.GetGlobal("event")
	if eventTable.Type() != lua.LTTable {
		t.Error("Expected event to be a table")
	}
}

// TestRegisterEvent_NilData tests registering event with nil data
func TestRegisterEvent_NilData(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      nil,
	}

	err := executor.registerEvent(L, event)
	if err != nil {
		t.Errorf("registerEvent() error = %v", err)
	}

	// Should still work with nil data
	eventTable := L.GetGlobal("event")
	if eventTable.Type() != lua.LTTable {
		t.Error("Expected event to be a table")
	}
}

// TestMapToLuaTable tests converting Go map to Lua table
func TestMapToLuaTable(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	goMap := map[string]interface{}{
		"string_field": "hello",
		"int_field":    42,
		"float_field":  3.14,
		"bool_field":   true,
	}

	luaTable, err := executor.mapToLuaTable(L, goMap)
	if err != nil {
		t.Errorf("mapToLuaTable() error = %v", err)
	}

	if luaTable == nil {
		t.Fatal("Expected non-nil table")
	}

	// Verify fields
	if luaTable.RawGetString("string_field").String() != "hello" {
		t.Error("String field not set correctly")
	}
}

// TestGoValueToLua_String tests converting string to Lua value
func TestGoValueToLua_String(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	val, err := executor.goValueToLua(L, "test string")
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTString {
		t.Error("Expected LTString type")
	}

	if val.String() != "test string" {
		t.Error("String value not correct")
	}
}

// TestGoValueToLua_Int tests converting int to Lua value
func TestGoValueToLua_Int(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	val, err := executor.goValueToLua(L, 42)
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTNumber {
		t.Error("Expected LTNumber type")
	}
}

// TestGoValueToLua_Int64 tests converting int64 to Lua value
func TestGoValueToLua_Int64(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	val, err := executor.goValueToLua(L, int64(9223372036854775807))
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTNumber {
		t.Error("Expected LTNumber type")
	}
}

// TestGoValueToLua_Float tests converting float64 to Lua value
func TestGoValueToLua_Float(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	val, err := executor.goValueToLua(L, 3.14159)
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTNumber {
		t.Error("Expected LTNumber type")
	}
}

// TestGoValueToLua_Bool tests converting bool to Lua value
func TestGoValueToLua_Bool(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	val, err := executor.goValueToLua(L, true)
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTBool {
		t.Error("Expected LTBool type")
	}

	if val == lua.LFalse {
		t.Error("Expected true value")
	}
}

// TestGoValueToLua_Slice tests converting slice to Lua array
func TestGoValueToLua_Slice(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	slice := []interface{}{"a", "b", "c"}

	val, err := executor.goValueToLua(L, slice)
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTTable {
		t.Error("Expected LTTable type for slice")
	}

	table := val.(*lua.LTable)
	if table.Len() != 3 {
		t.Errorf("Expected table length 3, got %d", table.Len())
	}
}

// TestGoValueToLua_StringSlice tests converting string slice to Lua array
func TestGoValueToLua_StringSlice(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	slice := []string{"x", "y", "z"}

	val, err := executor.goValueToLua(L, slice)
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTTable {
		t.Error("Expected LTTable type for string slice")
	}

	table := val.(*lua.LTable)
	if table.Len() != 3 {
		t.Errorf("Expected table length 3, got %d", table.Len())
	}
}

// TestGoValueToLua_Map tests converting map to Lua table
func TestGoValueToLua_Map(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{
		"key": "value",
	}

	val, err := executor.goValueToLua(L, m)
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTTable {
		t.Error("Expected LTTable type for map")
	}
}

// TestGoValueToLua_Nil tests converting nil to Lua value
func TestGoValueToLua_Nil(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	val, err := executor.goValueToLua(L, nil)
	if err != nil {
		t.Errorf("goValueToLua() error = %v", err)
	}

	if val.Type() != lua.LTNil {
		t.Error("Expected LTNil type")
	}
}

// TestGoValueToLua_UnsupportedType tests converting unsupported type
func TestGoValueToLua_UnsupportedType(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	// Struct that can be marshalled to JSON
	type TestStruct struct {
		Field string
	}

	val, err := executor.goValueToLua(L, TestStruct{Field: "test"})
	if err != nil {
		t.Errorf("goValueToLua() error = %v (should marshal to JSON)", err)
	}

	// Should return JSON string
	if val.Type() != lua.LTString {
		t.Error("Expected LTString for unsupported type")
	}
}

// TestRegisterCustomFunctions tests registering custom Lua functions
func TestRegisterCustomFunctions(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	L := lua.NewState()
	defer L.Close()

	var outBuf, errBuf []byte
	executor.registerCustomFunctions(L, os.Stdout, os.Stderr)

	// Verify log table exists
	logTable := L.GetGlobal("log")
	if logTable.Type() != lua.LTTable {
		t.Error("Expected log to be a table")
	}

	// Verify http table exists
	httpTable := L.GetGlobal("http")
	if httpTable.Type() != lua.LTTable {
		t.Error("Expected http to be a table")
	}

	// Verify contains function exists
	containsFunc := L.GetGlobal("contains")
	if containsFunc.Type() != lua.LTFunction {
		t.Error("Expected contains to be a function")
	}

	_ = outBuf
	_ = errBuf
}

// TestExecute_LogFunctions tests using log functions in hook
func TestExecute_LogFunctions(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "log_test.lua")

	script := `
function on_event()
	log.info("Info message")
	log.warn("Warning message")
	log.debug("Debug message")
	log.error("Error message")
	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "log-hook",
		Name:     "log-test",
		FilePath: hookFile,
		EventType: EventTaskStarted,
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}

	// Output should contain log messages
	if result.Output == "" {
		t.Error("Expected output from logging")
	}
}

// TestExecute_HTTPPost tests using http.post in hook
func TestExecute_HTTPPost(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "http_test.lua")

	script := `
function on_event()
	local success = http.post("https://example.com/webhook")
	log.info("HTTP POST result: " .. tostring(success))
	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "http-hook",
		Name:     "http-test",
		FilePath: hookFile,
		EventType: EventTaskCompleted,
	}

	event := &Event{
		Type:      EventTaskCompleted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}
}

// TestExecute_ContainsFunction tests using contains helper function
func TestExecute_ContainsFunction(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "contains_test.lua")

	script := `
function on_event()
	local list = {"apple", "banana", "cherry"}
	local has_banana = contains(list, "banana")
	local has_grape = contains(list, "grape")

	if has_banana then
		log.info("Found banana")
	end

	if not has_grape then
		log.info("Grape not found")
	end

	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "contains-hook",
		Name:     "contains-test",
		FilePath: hookFile,
		EventType: EventTaskStarted,
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}
}

// TestExecute_ComplexEventData tests hook with complex nested event data
func TestExecute_ComplexEventData(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "complex_data.lua")

	script := `
function on_event()
	log.info("Workflow: " .. event.data.workflow.name)
	log.info("Total tasks: " .. tostring(event.data.workflow.task_count))

	-- Access nested data
	if event.data.metadata then
		for k, v in pairs(event.data.metadata) do
			log.info(k .. " = " .. tostring(v))
		end
	end

	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "complex-hook",
		Name:     "complex-data",
		FilePath: hookFile,
		EventType: EventWorkflowCompleted,
	}

	event := &Event{
		Type:      EventWorkflowCompleted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"workflow": map[string]interface{}{
				"name":       "deploy-prod",
				"task_count": 15,
			},
			"metadata": map[string]interface{}{
				"duration": 120,
				"status":   "success",
			},
		},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got: %s", result.Error)
	}
}

// TestExecute_MeasuresDuration tests that execution duration is measured
func TestExecute_MeasuresDuration(t *testing.T) {
	repo, _ := NewRepository()
	defer repo.Close()

	executor := NewExecutor(repo)

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "duration_test.lua")

	script := `
function on_event()
	-- Simulate some work
	local sum = 0
	for i = 1, 1000 do
		sum = sum + i
	end
	return true
end
`

	os.WriteFile(hookFile, []byte(script), 0644)

	hook := &Hook{
		ID:       "duration-hook",
		Name:     "duration-test",
		FilePath: hookFile,
		EventType: EventTaskStarted,
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	result, err := executor.Execute(hook, event)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if result.Duration == 0 {
		t.Error("Expected non-zero duration")
	}

	if result.Duration < 0 {
		t.Error("Duration should not be negative")
	}
}
