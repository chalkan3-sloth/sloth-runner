package luainterface

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	lua "github.com/yuin/gopher-lua"
)

// StateModule provides state management and persistence functionality
type StateModule struct {
	stateManager *state.StateManager
	locksMux     *sync.RWMutex
}

var (
	globalStateModule *StateModule
	stateModuleOnce   sync.Once
)

// GetGlobalStateModule returns the singleton state module instance
func GetGlobalStateModule() *StateModule {
	stateModuleOnce.Do(func() {
		globalStateModule = NewStateModule("")
	})
	return globalStateModule
}

// NewStateModule creates a new StateModule 
func NewStateModule(dbPath string) *StateModule {
	stateManager, err := state.NewStateManager(dbPath)
	if err != nil {
		slog.Error("Failed to initialize state manager", "error", err)
		return nil
	}

	module := &StateModule{
		stateManager: stateManager,
		locksMux:     &sync.RWMutex{},
	}

	slog.Info("State module initialized", "backend", "state_manager")
	return module
}

// Loader is the module loader function for Lua
func (s *StateModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set":           s.luaStateSet,
		"get":           s.luaStateGet,
		"delete":        s.luaStateDelete,
		"exists":        s.luaStateExists,
		"keys":          s.luaStateKeys,
		"clear":         s.luaStateClear,
		"increment":     s.luaStateIncrement,
		"stats":         s.luaStateStats,
		"set_with_ttl":  s.luaStateSetWithTTL,
	})
	L.Push(mod)
	return 1
}

// StateLoader is the global loader function
func StateLoader(L *lua.LState) int {
	return GetGlobalStateModule().Loader(L)
}

// luaStateSet sets a value in the state store
func (s *StateModule) luaStateSet(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckAny(2)
	
	// Convert Lua value to string (StateManager expects strings)
	valueStr := luaValueToString(value)
	
	err := s.stateManager.Set(key, valueStr)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// luaStateGet gets a value from the state store
func (s *StateModule) luaStateGet(L *lua.LState) int {
	key := L.CheckString(1)
	
	// Use interface{} to handle both implementations
	value, err := s.stateManager.Get(key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Convert to string safely
	strValue := fmt.Sprintf("%v", value)
	
	L.Push(stringToLua(strValue))
	return 1
}

// luaStateDelete deletes a key from the state store
func (s *StateModule) luaStateDelete(L *lua.LState) int {
	key := L.CheckString(1)
	
	err := s.stateManager.Delete(key)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// luaStateExists checks if a key exists
func (s *StateModule) luaStateExists(L *lua.LState) int {
	key := L.CheckString(1)
	
	exists, err := s.stateManager.Exists(key)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(exists))
	return 1
}

// luaStateKeys returns all keys with optional prefix
func (s *StateModule) luaStateKeys(L *lua.LState) int {
	prefix := ""
	if L.GetTop() >= 1 {
		prefix = L.CheckString(1)
	}
	
	keys, err := s.stateManager.List(prefix)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Convert map to Lua table
	table := L.NewTable()
	for key, value := range keys {
		table.RawSetString(key, lua.LString(value))
	}
	
	L.Push(table)
	return 1
}

// luaStateClear clears all state
func (s *StateModule) luaStateClear(L *lua.LState) int {
	err := s.stateManager.Clear()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// luaStateIncrement increments a numeric value
func (s *StateModule) luaStateIncrement(L *lua.LState) int {
	key := L.CheckString(1)
	delta := int64(1)
	if L.GetTop() >= 2 {
		delta = int64(L.CheckNumber(2))
	}
	
	newValue, err := s.stateManager.Increment(key, delta)
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LNumber(newValue))
	return 1
}

// luaStateStats returns state statistics
func (s *StateModule) luaStateStats(L *lua.LState) int {
	stats, err := s.stateManager.Stats()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	table := L.NewTable()
	table.RawSetString("total_keys", lua.LNumber(stats.TotalKeys))
	table.RawSetString("total_size", lua.LNumber(stats.TotalSize))
	table.RawSetString("last_modified", lua.LNumber(stats.LastModified))
	table.RawSetString("backend", lua.LString(stats.Backend))
	
	L.Push(table)
	return 1
}

// luaStateSetWithTTL sets a value with TTL
func (s *StateModule) luaStateSetWithTTL(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckAny(2)
	ttlSeconds := int(L.CheckNumber(3))
	
	// Convert Lua value to Go interface{}
	goValue := luaValueToGo(value)
	
	err := s.stateManager.SetWithTTL(key, goValue, ttlSeconds)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// Helper functions to convert between Lua and Go values

// luaValueToString converts a Lua value to string
func luaValueToString(value lua.LValue) string {
	switch v := value.(type) {
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return fmt.Sprintf("%g", float64(v))
	case lua.LBool:
		return fmt.Sprintf("%t", bool(v))
	case *lua.LNilType:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// stringToLua converts a string to Lua value
func stringToLua(str string) lua.LValue {
	return lua.LString(str)
}

func luaValueToGo(value lua.LValue) interface{} {
	switch v := value.(type) {
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LNilType:
		return nil
	case *lua.LTable:
		// Convert table to map[string]interface{}
		result := make(map[string]interface{})
		v.ForEach(func(k, val lua.LValue) {
			if key, ok := k.(lua.LString); ok {
				result[string(key)] = luaValueToGo(val)
			}
		})
		return result
	default:
		return fmt.Sprintf("%v", v)
	}
}

func goValueToLua(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case nil:
		return lua.LNil
	case map[string]interface{}:
		table := L.NewTable()
		for key, val := range v {
			table.RawSetString(key, goValueToLua(L, val))
		}
		return table
	default:
		// Try to convert to string
		if str, ok := v.(string); ok {
			return lua.LString(str)
		}
		// Try to parse as number
		if str := fmt.Sprintf("%v", v); str != "" {
			if num, err := strconv.ParseFloat(str, 64); err == nil {
				return lua.LNumber(num)
			}
			return lua.LString(str)
		}
		return lua.LNil
	}
}