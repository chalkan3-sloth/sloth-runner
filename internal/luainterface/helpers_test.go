package luainterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestLuaTableToGoMap_SimpleTypes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("string_key", lua.LString("string_value"))
	table.RawSetString("number_key", lua.LNumber(42))
	table.RawSetString("bool_key", lua.LBool(true))

	result := LuaTableToGoMap(L, table)

	assert.Equal(t, "string_value", result["string_key"])
	// LVAsNumber returns lua.LNumber which needs to be converted
	numVal := result["number_key"]
	assert.Equal(t, lua.LNumber(42), numVal)
	assert.Equal(t, true, result["bool_key"])
}

func TestLuaTableToGoMap_NestedTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	nested := L.NewTable()
	nested.RawSetString("inner_key", lua.LString("inner_value"))

	table := L.NewTable()
	table.RawSetString("nested", nested)

	result := LuaTableToGoMap(L, table)

	nestedResult, ok := result["nested"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "inner_value", nestedResult["inner_key"])
}

func TestLuaTableToGoMap_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.Append(lua.LString("item1"))
	table.Append(lua.LString("item2"))
	table.Append(lua.LString("item3"))

	result := LuaTableToGoMap(L, table)

	// Lua arrays start at 1
	assert.Equal(t, "item1", result["1"])
	assert.Equal(t, "item2", result["2"])
	assert.Equal(t, "item3", result["3"])
}

func TestLuaTableToGoMap_MixedKeys(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("key1", lua.LString("value1"))
	table.RawSet(lua.LNumber(1), lua.LString("indexed_value"))

	result := LuaTableToGoMap(L, table)

	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "indexed_value", result["1"])
}

func TestLuaTableToGoMap_EmptyTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	result := LuaTableToGoMap(L, table)

	assert.Empty(t, result)
}

func TestLuaTableToGoMap_NilValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("nil_key", lua.LNil)
	table.RawSetString("valid_key", lua.LString("value"))

	result := LuaTableToGoMap(L, table)

	// Nil values should be filtered out or represented as nil
	assert.Contains(t, result, "valid_key")
	assert.Equal(t, "value", result["valid_key"])
}

func TestLuaToGoValue_String(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := LuaToGoValue(L, lua.LString("test"))
	assert.Equal(t, "test", value)
}

func TestLuaToGoValue_Number(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := LuaToGoValue(L, lua.LNumber(42.5))
	assert.Equal(t, 42.5, value)
}

func TestLuaToGoValue_Bool(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	valueTrue := LuaToGoValue(L, lua.LBool(true))
	assert.Equal(t, true, valueTrue)

	valueFalse := LuaToGoValue(L, lua.LBool(false))
	assert.Equal(t, false, valueFalse)
}

func TestLuaToGoValue_Nil(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	value := LuaToGoValue(L, lua.LNil)
	assert.Nil(t, value)
}

func TestLuaToGoValue_Table(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("key", lua.LString("value"))

	value := LuaToGoValue(L, table)
	
	mapValue, ok := value.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", mapValue["key"])
}

func TestLuaToGoValue_ComplexTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("string", lua.LString("value"))
	table.RawSetString("number", lua.LNumber(42))
	table.RawSetString("bool", lua.LBool(true))

	nested := L.NewTable()
	nested.RawSetString("inner", lua.LString("nested value"))
	table.RawSetString("nested", nested)

	value := LuaToGoValue(L, table)
	
	mapValue, ok := value.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", mapValue["string"])
	assert.Equal(t, 42.0, mapValue["number"])
	assert.Equal(t, true, mapValue["bool"])
	
	nestedMap, ok := mapValue["nested"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "nested value", nestedMap["inner"])
}
