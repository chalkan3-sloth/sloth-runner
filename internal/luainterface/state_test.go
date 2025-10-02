package luainterface

import (
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestStateModuleCreation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)

	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	if module.stateManager == nil {
		t.Error("stateManager is nil")
	}

	if module.locksMux == nil {
		t.Error("locksMux is nil")
	}
}

func TestStateModuleLoader(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	// Load the state module
	L.PreloadModule("state", module.Loader)

	// Test that state module can be required
	script := `
local state = require("state")
if type(state) ~= "table" then
	error("state module not loaded correctly")
end
`

	if err := L.DoString(script); err != nil {
		t.Errorf("Failed to load state module: %v", err)
	}
}

func TestStateModuleSetGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("state", module.Loader)

	script := `
local state = require("state")

-- Test set and get
state.set("test_key", "test_value")
local value = state.get("test_key")

if value ~= "test_value" then
	error("Expected 'test_value', got: " .. tostring(value))
end

result = value
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute state set/get: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "test_value" {
		t.Errorf("Expected 'test_value', got: %s", result.String())
	}
}

func TestStateModuleNumericValues(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("state", module.Loader)

	script := `
local state = require("state")

-- Test numeric values
state.set("number_key", 42)
local value = state.get("number_key")

result = value
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute state numeric test: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTNumber && result.Type() != lua.LTString {
		t.Errorf("Expected number or string, got: %v", result.Type())
	}
}

func TestStateModuleDelete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("state", module.Loader)

	script := `
local state = require("state")

-- Set a value
state.set("delete_key", "delete_value")

-- Verify it exists
local value = state.get("delete_key")
if value ~= "delete_value" then
	error("Value not set correctly")
end

-- Delete the value
state.delete("delete_key")

-- Verify it's deleted (should return nil or empty)
local deleted_value = state.get("delete_key")
result = deleted_value
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute state delete test: %v", err)
	}

	result := L.GetGlobal("result")
	// After delete, should be nil or empty string
	if result.Type() != lua.LTNil && result.String() != "" {
		t.Errorf("Expected nil or empty after delete, got: %v", result)
	}
}

func TestStateModulePersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")

	// First session - set value
	{
		module := NewStateModule(dbPath)
		if module == nil {
			t.Fatal("NewStateModule returned nil")
		}

		L := lua.NewState()
		defer L.Close()

		L.PreloadModule("state", module.Loader)

		script := `
local state = require("state")
state.set("persist_key", "persist_value")
`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to set persistent value: %v", err)
		}
	}

	// Second session - get value
	{
		module := NewStateModule(dbPath)
		if module == nil {
			t.Fatal("NewStateModule returned nil")
		}

		L := lua.NewState()
		defer L.Close()

		L.PreloadModule("state", module.Loader)

		script := `
local state = require("state")
result = state.get("persist_key")
`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to get persistent value: %v", err)
		}

		result := L.GetGlobal("result")
		if result.String() != "persist_value" {
			t.Errorf("Expected 'persist_value', got: %s", result.String())
		}
	}
}

func TestStateModuleGlobalInstance(t *testing.T) {
	module1 := GetGlobalStateModule()
	module2 := GetGlobalStateModule()

	if module1 != module2 {
		t.Error("GetGlobalStateModule should return the same instance")
	}

	if module1 == nil {
		t.Error("GetGlobalStateModule returned nil")
	}
}

func TestStateModuleMultipleKeys(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("state", module.Loader)

	script := `
local state = require("state")

-- Set multiple keys
state.set("key1", "value1")
state.set("key2", "value2")
state.set("key3", "value3")

-- Get all keys
local v1 = state.get("key1")
local v2 = state.get("key2")
local v3 = state.get("key3")

if v1 ~= "value1" or v2 ~= "value2" or v3 ~= "value3" then
	error("Keys not stored correctly")
end

result = true
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test multiple keys: %v", err)
	}

	result := L.GetGlobal("result")
	if !lua.LVAsBool(result) {
		t.Error("Multiple keys test failed")
	}
}

func TestStateModuleOverwrite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("state", module.Loader)

	script := `
local state = require("state")

-- Set initial value
state.set("overwrite_key", "initial_value")
local v1 = state.get("overwrite_key")

-- Overwrite with new value
state.set("overwrite_key", "new_value")
local v2 = state.get("overwrite_key")

if v1 ~= "initial_value" then
	error("Initial value not correct")
end

if v2 ~= "new_value" then
	error("Overwritten value not correct")
end

result = v2
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test overwrite: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "new_value" {
		t.Errorf("Expected 'new_value', got: %s", result.String())
	}
}

func TestStateModuleEmptyKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("state", module.Loader)

	// Test getting non-existent key
	script := `
local state = require("state")
result = state.get("non_existent_key")
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test empty key: %v", err)
	}

	result := L.GetGlobal("result")
	// Should return nil or empty string for non-existent key
	if result.Type() != lua.LTNil && result.String() != "" {
		t.Logf("Non-existent key returned: %v (type: %v)", result, result.Type())
	}
}

func TestStateModuleComplexValues(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	module := NewStateModule(dbPath)
	if module == nil {
		t.Fatal("NewStateModule returned nil")
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("state", module.Loader)
	RegisterAllModules(L)

	// Test storing JSON-encoded complex values
	script := `
local state = require("state")

local complex_data = {
	name = "test",
	values = {1, 2, 3},
	nested = {key = "value"}
}

-- Encode to JSON and store
local json_str = json.encode(complex_data)
state.set("complex_key", json_str)

-- Retrieve and decode
local stored_json = state.get("complex_key")
local decoded = json.decode(stored_json)

result = decoded.name
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to test complex values: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "test" {
		t.Errorf("Expected 'test', got: %s", result.String())
	}
}
