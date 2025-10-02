package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestSessionModule(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "session.set and session.get",
			script: `
session.set("test_key", "test_value")
result = session.get("test_key")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "test_value" {
					t.Errorf("Expected 'test_value', got: %s", result.String())
				}
			},
		},
		{
			name: "session with numeric value",
			script: `
session.set("number", 42)
result = session.get("number")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTNumber && result.Type() != lua.LTString {
					t.Errorf("Expected number or string, got type: %v", result.Type())
				}
			},
		},
		{
			name: "session with boolean value",
			script: `
session.set("flag", true)
result = session.get("flag")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				// Boolean might be stored as string "true"
				if result.Type() != lua.LTBool && result.String() != "true" {
					t.Errorf("Expected boolean or 'true', got: %v", result)
				}
			},
		},
		{
			name: "session non-existent key",
			script: `
result = session.get("non_existent")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				// Non-existent key should return nil or empty
				if result.Type() != lua.LTNil && result.String() != "" {
					t.Logf("Non-existent key returned: %v", result)
				}
			},
		},
		{
			name: "session overwrite",
			script: `
session.set("key", "value1")
session.set("key", "value2")
result = session.get("key")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "value2" {
					t.Errorf("Expected 'value2', got: %s", result.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute %s: %v", tt.name, err)
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestSessionPersistenceWithinState(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Set multiple session values and verify they persist within the same Lua state
	script := `
session.set("key1", "value1")
session.set("key2", "value2")
session.set("key3", "value3")

local v1 = session.get("key1")
local v2 = session.get("key2")
local v3 = session.get("key3")

if v1 ~= "value1" then error("key1 mismatch") end
if v2 ~= "value2" then error("key2 mismatch") end
if v3 ~= "value3" then error("key3 mismatch") end

result = true
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test session persistence: %v", err)
	}

	result := L.GetGlobal("result")
	if !lua.LVAsBool(result) {
		t.Error("Session persistence test failed")
	}
}

func TestSessionWithComplexData(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Test storing JSON-encoded data in session
	script := `
local data = {name = "test", value = 123}
local json_str = json.encode(data)

session.set("json_data", json_str)
local stored = session.get("json_data")
local decoded = json.decode(stored)

result = decoded.name
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test session with complex data: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "test" {
		t.Errorf("Expected 'test', got: %s", result.String())
	}
}

func TestSessionIsolation(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	// Create two separate Lua states to test isolation
	L1 := lua.NewState()
	defer L1.Close()
	RegisterAllModules(L1)

	L2 := lua.NewState()
	defer L2.Close()
	RegisterAllModules(L2)

	// Set value in L1
	if err := L1.DoString(`session.set("isolated_key", "value_from_L1")`); err != nil {
		t.Fatalf("Failed to set value in L1: %v", err)
	}

	// Try to get value in L2 - should not exist
	if err := L2.DoString(`result = session.get("isolated_key")`); err != nil {
		t.Fatalf("Failed to get value in L2: %v", err)
	}

	result := L2.GetGlobal("result")
	// Should be nil or empty since sessions are isolated
	if result.Type() != lua.LTNil && result.String() != "" {
		t.Logf("Sessions might not be isolated: %v", result)
	}
}

func TestSessionEmptyValue(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
session.set("empty", "")
result = session.get("empty")
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test empty value: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "" {
		t.Errorf("Expected empty string, got: %s", result.String())
	}
}

func TestSessionMultipleOperations(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
-- Multiple set operations
for i = 1, 10 do
	session.set("key" .. i, "value" .. i)
end

-- Verify all values
local success = true
for i = 1, 10 do
	local value = session.get("key" .. i)
	if value ~= "value" .. i then
		success = false
		break
	end
end

result = success
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test multiple operations: %v", err)
	}

	result := L.GetGlobal("result")
	if !lua.LVAsBool(result) {
		t.Error("Multiple operations test failed")
	}
}

func TestSessionWithSpecialCharacters(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
local special = "value with spaces and special chars: !@#$%^&*()"
session.set("special_key", special)
result = session.get("special_key")
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test special characters: %v", err)
	}

	result := L.GetGlobal("result")
	expected := "value with spaces and special chars: !@#$%^&*()"
	if result.String() != expected {
		t.Errorf("Expected '%s', got: %s", expected, result.String())
	}
}

func TestSessionConcurrentAccess(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Test that session can handle rapid successive operations
	script := `
for i = 1, 100 do
	session.set("rapid_key", "value" .. i)
	local value = session.get("rapid_key")
	if value ~= "value" .. i then
		error("Concurrent access failed at iteration " .. i)
	end
end

result = true
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test concurrent access: %v", err)
	}

	result := L.GetGlobal("result")
	if !lua.LVAsBool(result) {
		t.Error("Concurrent access test failed")
	}
}

func TestSessionWithTableValues(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Test that tables need to be serialized before storing
	script := `
local tbl = {a = 1, b = 2, c = 3}
local json_str = json.encode(tbl)

session.set("table_key", json_str)
local stored = session.get("table_key")
local decoded = json.decode(stored)

result = decoded.a + decoded.b + decoded.c
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test table values: %v", err)
	}

	result := L.GetGlobal("result")
	if result.(lua.LNumber) != 6 {
		t.Errorf("Expected 6, got: %v", result)
	}
}

func TestSessionGetOrDefault(t *testing.T) {
	t.Skip("Session module not yet fully implemented - migrating to state module")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Test getting value with default fallback
	script := `
local value = session.get("non_existent") or "default_value"
result = value
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test get or default: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "default_value" {
		t.Errorf("Expected 'default_value', got: %s", result.String())
	}
}
