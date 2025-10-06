package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/modules"
	lua "github.com/yuin/gopher-lua"
)

// Executor executes hook scripts
type Executor struct {
	repo *Repository
}

// NewExecutor creates a new hook executor
func NewExecutor(repo *Repository) *Executor {
	return &Executor{
		repo: repo,
	}
}

// Execute executes a hook with the given event data
func (e *Executor) Execute(hook *Hook, event *Event) (*HookResult, error) {
	startTime := time.Now()

	result := &HookResult{
		HookID:     hook.ID,
		ExecutedAt: startTime,
	}

	// Check if file exists
	if _, err := os.Stat(hook.FilePath); os.IsNotExist(err) {
		result.Success = false
		result.Error = fmt.Sprintf("hook file not found: %s", hook.FilePath)
		result.Duration = time.Since(startTime)
		return result, fmt.Errorf("%s", result.Error)
	}

	// Create Lua state
	L := lua.NewState()
	defer L.Close()

	// Capture output
	var outputBuf bytes.Buffer
	var errorBuf bytes.Buffer

	// Register event data in Lua
	if err := e.registerEvent(L, event); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to register event: %v", err)
		result.Duration = time.Since(startTime)
		return result, err
	}

	// Register custom functions
	e.registerCustomFunctions(L, &outputBuf, &errorBuf)

	// Load all sloth-runner modules
	registry := modules.GetGlobalRegistry()
	if err := registry.LoadAllModules(L); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to load modules: %v", err)
		result.Output = outputBuf.String()
		result.Duration = time.Since(startTime)
		return result, err
	}

	// Execute the hook file
	if err := L.DoFile(hook.FilePath); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("execution error: %v", err)
		result.Output = outputBuf.String()
		result.Duration = time.Since(startTime)
		return result, err
	}

	// Check if hook function exists and call it
	hookFn := L.GetGlobal("on_event")
	if hookFn.Type() == lua.LTFunction {
		// Call the hook function
		if err := L.CallByParam(lua.P{
			Fn:      hookFn,
			NRet:    1,
			Protect: true,
		}); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("hook function error: %v", err)
			result.Output = outputBuf.String()
			result.Duration = time.Since(startTime)
			return result, err
		}

		// Get return value (success boolean)
		ret := L.Get(-1)
		L.Pop(1)

		if ret.Type() == lua.LTBool {
			result.Success = lua.LVAsBool(ret)
		} else {
			result.Success = true // Default to success if no boolean returned
		}
	} else {
		// No on_event function found, consider it successful if no errors
		result.Success = true
	}

	result.Output = outputBuf.String()
	if errorBuf.Len() > 0 && !result.Success {
		result.Error = errorBuf.String()
	}
	result.Duration = time.Since(startTime)

	return result, nil
}

// registerEvent registers the event data in Lua state
func (e *Executor) registerEvent(L *lua.LState, event *Event) error {
	// Create event table
	eventTable := L.NewTable()

	// Set event type
	eventTable.RawSetString("type", lua.LString(event.Type))
	eventTable.RawSetString("timestamp", lua.LNumber(event.Timestamp.Unix()))

	// Convert event data to Lua table
	if event.Data != nil {
		dataTable, err := e.mapToLuaTable(L, event.Data)
		if err != nil {
			return fmt.Errorf("failed to convert event data: %w", err)
		}
		eventTable.RawSetString("data", dataTable)

		// Also set common fields at event level for convenience
		if agent, ok := event.Data["agent"].(map[string]interface{}); ok {
			agentTable, _ := e.mapToLuaTable(L, agent)
			eventTable.RawSetString("agent", agentTable)
		}

		if task, ok := event.Data["task"].(map[string]interface{}); ok {
			taskTable, _ := e.mapToLuaTable(L, task)
			eventTable.RawSetString("task", taskTable)
		}
	}

	// Set as global
	L.SetGlobal("event", eventTable)

	return nil
}

// mapToLuaTable converts a Go map to Lua table
func (e *Executor) mapToLuaTable(L *lua.LState, m map[string]interface{}) (*lua.LTable, error) {
	table := L.NewTable()

	for k, v := range m {
		luaValue, err := e.goValueToLua(L, v)
		if err != nil {
			return nil, err
		}
		table.RawSetString(k, luaValue)
	}

	return table, nil
}

// goValueToLua converts Go value to Lua value
func (e *Executor) goValueToLua(L *lua.LState, val interface{}) (lua.LValue, error) {
	switch v := val.(type) {
	case string:
		return lua.LString(v), nil
	case int:
		return lua.LNumber(v), nil
	case int64:
		return lua.LNumber(v), nil
	case float64:
		return lua.LNumber(v), nil
	case bool:
		return lua.LBool(v), nil
	case []interface{}:
		arr := L.NewTable()
		for i, item := range v {
			luaVal, err := e.goValueToLua(L, item)
			if err != nil {
				return lua.LNil, err
			}
			arr.RawSetInt(i+1, luaVal)
		}
		return arr, nil
	case []string:
		arr := L.NewTable()
		for i, item := range v {
			arr.RawSetInt(i+1, lua.LString(item))
		}
		return arr, nil
	case map[string]interface{}:
		return e.mapToLuaTable(L, v)
	case nil:
		return lua.LNil, nil
	default:
		// Try to marshal as JSON and parse
		data, err := json.Marshal(v)
		if err != nil {
			return lua.LNil, fmt.Errorf("unsupported type: %T", v)
		}
		return lua.LString(string(data)), nil
	}
}

// registerCustomFunctions registers custom Lua functions
func (e *Executor) registerCustomFunctions(L *lua.LState, out, err io.Writer) {
	// log.info function
	L.SetGlobal("log", L.NewTable())
	logTable := L.GetGlobal("log").(*lua.LTable)

	logTable.RawSetString("info", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		fmt.Fprintf(out, "[INFO] %s\n", msg)
		return 0
	}))

	logTable.RawSetString("error", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		fmt.Fprintf(err, "[ERROR] %s\n", msg)
		return 0
	}))

	logTable.RawSetString("warn", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		fmt.Fprintf(out, "[WARN] %s\n", msg)
		return 0
	}))

	logTable.RawSetString("debug", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		fmt.Fprintf(out, "[DEBUG] %s\n", msg)
		return 0
	}))

	// http module for sending webhooks (basic implementation)
	L.SetGlobal("http", L.NewTable())
	httpTable := L.GetGlobal("http").(*lua.LTable)

	httpTable.RawSetString("post", L.NewFunction(func(L *lua.LState) int {
		url := L.CheckString(1)
		// Basic stub - in real implementation, would send HTTP POST
		fmt.Fprintf(out, "[HTTP] POST to %s\n", url)
		L.Push(lua.LTrue)
		return 1
	}))

	// Helper function to check if list contains value
	L.SetGlobal("contains", L.NewFunction(func(L *lua.LState) int {
		list := L.CheckTable(1)
		value := L.CheckString(2)

		found := false
		list.ForEach(func(k, v lua.LValue) {
			if v.String() == value {
				found = true
			}
		})

		L.Push(lua.LBool(found))
		return 1
	}))
}
