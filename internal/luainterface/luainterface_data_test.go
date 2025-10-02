package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestGoValueToLua_String(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := GoValueToLua(L, "test string")
	if str, ok := value.(lua.LString); !ok || string(str) != "test string" {
		t.Errorf("Expected LString 'test string', got %v", value)
	}
}

func TestGoValueToLua_Int(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := GoValueToLua(L, 42)
	if num, ok := value.(lua.LNumber); !ok || int(num) != 42 {
		t.Errorf("Expected LNumber 42, got %v", value)
	}
}

func TestGoValueToLua_Float(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := GoValueToLua(L, 3.14)
	if num, ok := value.(lua.LNumber); !ok || num != 3.14 {
		t.Errorf("Expected LNumber 3.14, got %v", value)
	}
}

func TestGoValueToLua_Bool(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := GoValueToLua(L, true)
	if b, ok := value.(lua.LBool); !ok || bool(b) != true {
		t.Errorf("Expected LBool true, got %v", value)
	}
}

func TestGoValueToLua_Nil(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := GoValueToLua(L, nil)
	if value != lua.LNil {
		t.Errorf("Expected LNil, got %v", value)
	}
}

func TestGoValueToLua_Map(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	testMap := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	value := GoValueToLua(L, testMap)
	table, ok := value.(*lua.LTable)
	if !ok {
		t.Fatalf("Expected LTable, got %T", value)
	}

	if table.RawGetString("key1").String() != "value1" {
		t.Errorf("Expected key1='value1', got %v", table.RawGetString("key1"))
	}
	if int(table.RawGetString("key2").(lua.LNumber)) != 42 {
		t.Errorf("Expected key2=42, got %v", table.RawGetString("key2"))
	}
	if bool(table.RawGetString("key3").(lua.LBool)) != true {
		t.Errorf("Expected key3=true, got %v", table.RawGetString("key3"))
	}
}

func TestGoValueToLua_Slice(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	testSlice := []interface{}{"a", "b", "c"}

	value := GoValueToLua(L, testSlice)
	table, ok := value.(*lua.LTable)
	if !ok {
		t.Fatalf("Expected LTable, got %T", value)
	}

	if table.RawGetInt(1).String() != "a" {
		t.Errorf("Expected index 1='a', got %v", table.RawGetInt(1))
	}
	if table.RawGetInt(2).String() != "b" {
		t.Errorf("Expected index 2='b', got %v", table.RawGetInt(2))
	}
	if table.RawGetInt(3).String() != "c" {
		t.Errorf("Expected index 3='c', got %v", table.RawGetInt(3))
	}
}

func TestGoValueToLua_NestedStructure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	nested := map[string]interface{}{
		"outer": map[string]interface{}{
			"inner": []interface{}{1, 2, 3},
		},
	}

	value := GoValueToLua(L, nested)
	table, ok := value.(*lua.LTable)
	if !ok {
		t.Fatalf("Expected LTable, got %T", value)
	}

	outerTable := table.RawGetString("outer").(*lua.LTable)
	innerTable := outerTable.RawGetString("inner").(*lua.LTable)
	
	if int(innerTable.RawGetInt(1).(lua.LNumber)) != 1 {
		t.Errorf("Expected nested value 1, got %v", innerTable.RawGetInt(1))
	}
}

func TestLuaDataParseJson_Valid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	OpenData(L)

	script := `
		local data = require("data")
		local result = data.parse_json('{"name": "test", "value": 42}')
		return result.name, result.value
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	name := L.Get(-2).String()
	value := int(L.Get(-1).(lua.LNumber))

	if name != "test" {
		t.Errorf("Expected name='test', got '%s'", name)
	}
	if value != 42 {
		t.Errorf("Expected value=42, got %d", value)
	}
}

func TestLuaDataToJson_Valid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	OpenData(L)

	script := `
		local data = require("data")
		local tbl = {name = "test", value = 42}
		local json_str = data.to_json(tbl)
		return json_str
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	result := L.Get(-1).String()
	if result == "" {
		t.Error("Expected non-empty JSON string")
	}
	
	// Verify it contains expected data
	if !contains(result, "test") || !contains(result, "42") {
		t.Errorf("JSON doesn't contain expected values: %s", result)
	}
}

func TestLuaDataParseYaml_Valid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	OpenData(L)

	script := `
		local data = require("data")
		local result = data.parse_yaml("name: test\nvalue: 42")
		return result.name, result.value
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	name := L.Get(-2).String()
	value := int(L.Get(-1).(lua.LNumber))

	if name != "test" {
		t.Errorf("Expected name='test', got '%s'", name)
	}
	if value != 42 {
		t.Errorf("Expected value=42, got %d", value)
	}
}

func TestLuaDataToYaml_Valid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	OpenData(L)

	script := `
		local data = require("data")
		local tbl = {name = "test", value = 42}
		local yaml_str = data.to_yaml(tbl)
		return yaml_str
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	result := L.Get(-1).String()
	if result == "" {
		t.Error("Expected non-empty YAML string")
	}
	
	// Verify it contains expected data
	if !contains(result, "test") {
		t.Errorf("YAML doesn't contain expected values: %s", result)
	}
}

func TestOpenData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	OpenData(L)

	// Test that data module can be accessed and used
	script := `
		local data = require("data")
		return type(data), type(data.parse_json), type(data.to_json)
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	dataType := L.Get(-3).String()
	parseJsonType := L.Get(-2).String()
	toJsonType := L.Get(-1).String()

	if dataType != "table" {
		t.Errorf("Expected data to be 'table', got '%s'", dataType)
	}
	if parseJsonType != "function" {
		t.Errorf("Expected parse_json to be 'function', got '%s'", parseJsonType)
	}
	if toJsonType != "function" {
		t.Errorf("Expected to_json to be 'function', got '%s'", toJsonType)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || len(s) > len(substr)+1 && anyIndex(s, substr)))
}

func anyIndex(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
