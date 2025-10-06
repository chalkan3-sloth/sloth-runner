package modules

import (
	lua "github.com/yuin/gopher-lua"
)

// LuaHelpers provides common helper functions for Lua module returns
type LuaHelpers struct{}

// ReturnSuccess returns a successful result with nil error
// This is the standard pattern: (result, nil)
// Usage: return helpers.ReturnSuccess(L, resultTable)
func (h *LuaHelpers) ReturnSuccess(L *lua.LState, result lua.LValue) int {
	L.Push(result)
	L.Push(lua.LNil)
	return 2
}

// ReturnError returns nil result with error message
// This is the standard pattern: (nil, error)
// Usage: return helpers.ReturnError(L, "error message")
func (h *LuaHelpers) ReturnError(L *lua.LState, message string) int {
	L.Push(lua.LNil)
	L.Push(lua.LString(message))
	return 2
}

// ReturnFluentSuccess returns self with nil error for fluent API
// This is the fluent API pattern: (self, nil)
// Usage: return helpers.ReturnFluentSuccess(L, selfUserData)
func (h *LuaHelpers) ReturnFluentSuccess(L *lua.LState, self lua.LValue) int {
	L.Push(self)
	L.Push(lua.LNil)
	return 2
}

// CreateResultTable creates a standardized result table
// with changed, message, and optional additional fields
func (h *LuaHelpers) CreateResultTable(L *lua.LState, changed bool, message string, fields map[string]lua.LValue) *lua.LTable {
	result := L.NewTable()
	result.RawSetString("changed", lua.LBool(changed))
	result.RawSetString("message", lua.LString(message))

	for key, value := range fields {
		result.RawSetString(key, value)
	}

	return result
}

// ReturnIdempotent returns an idempotent result (no changes needed)
// This is used when an operation doesn't need to be performed
// Usage: return helpers.ReturnIdempotent(L, "Resource already exists")
func (h *LuaHelpers) ReturnIdempotent(L *lua.LState, message string) int {
	result := h.CreateResultTable(L, false, message, nil)
	return h.ReturnSuccess(L, result)
}

// ReturnChanged returns a successful result with changed=true
// This is used when an operation was performed successfully
// Usage: return helpers.ReturnChanged(L, "Resource created", extraFields)
func (h *LuaHelpers) ReturnChanged(L *lua.LState, message string, fields map[string]lua.LValue) int {
	result := h.CreateResultTable(L, true, message, fields)
	return h.ReturnSuccess(L, result)
}

// GetStringField safely gets a string field from a Lua table with default value
func GetStringField(L *lua.LState, table *lua.LTable, key, defaultValue string) string {
	lv := table.RawGetString(key)
	if lv.Type() == lua.LTString {
		return lv.String()
	}
	return defaultValue
}

// GetBoolField safely gets a boolean field from a Lua table with default value
func GetBoolField(L *lua.LState, table *lua.LTable, key string, defaultValue bool) bool {
	lv := table.RawGetString(key)
	if lv.Type() == lua.LTBool {
		return bool(lv.(lua.LBool))
	}
	return defaultValue
}

// GetIntField safely gets an integer field from a Lua table with default value
func GetIntField(L *lua.LState, table *lua.LTable, key string, defaultValue int) int {
	lv := table.RawGetString(key)
	if lv.Type() == lua.LTNumber {
		return int(lv.(lua.LNumber))
	}
	return defaultValue
}

// GetTableField safely gets a table field from a Lua table
func GetTableField(L *lua.LState, table *lua.LTable, key string) *lua.LTable {
	lv := table.RawGetString(key)
	if lv.Type() == lua.LTTable {
		return lv.(*lua.LTable)
	}
	return nil
}

// Global helper instance
var Helpers = &LuaHelpers{}
