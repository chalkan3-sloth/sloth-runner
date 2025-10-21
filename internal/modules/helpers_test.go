package modules

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// Test LuaHelpers struct creation
func TestNewLuaHelpers(t *testing.T) {
	helpers := &LuaHelpers{}

	if helpers == nil {
		t.Error("Expected non-nil LuaHelpers")
	}
}

func TestGlobalHelpersInstance(t *testing.T) {
	if Helpers == nil {
		t.Error("Expected global Helpers instance to be initialized")
	}
}

func TestGlobalHelpersType(t *testing.T) {
	var _ *LuaHelpers = Helpers
}

// Test ReturnSuccess
func TestLuaHelpers_ReturnSuccess(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	result := L.NewTable()
	L.SetField(result, "status", lua.LString("ok"))

	count := helpers.ReturnSuccess(L, result)

	if count != 2 {
		t.Errorf("Expected ReturnSuccess to return 2 values, got %d", count)
	}

	// Check stack: should have result table and nil error
	if L.GetTop() != 2 {
		t.Errorf("Expected 2 values on stack, got %d", L.GetTop())
	}

	resultValue := L.Get(-2)
	if resultValue.Type() != lua.LTTable {
		t.Errorf("Expected table result, got %s", resultValue.Type())
	}

	errorValue := L.Get(-1)
	if errorValue.Type() != lua.LTNil {
		t.Errorf("Expected nil error, got %s", errorValue.Type())
	}
}

func TestLuaHelpers_ReturnSuccess_WithString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	result := lua.LString("success")

	count := helpers.ReturnSuccess(L, result)

	if count != 2 {
		t.Errorf("Expected 2 return values, got %d", count)
	}

	resultValue := L.Get(-2)
	if resultValue.String() != "success" {
		t.Errorf("Expected 'success', got %s", resultValue.String())
	}
}

func TestLuaHelpers_ReturnSuccess_WithNumber(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	result := lua.LNumber(42)

	helpers.ReturnSuccess(L, result)

	resultValue := L.Get(-2)
	if resultValue != lua.LNumber(42) {
		t.Errorf("Expected 42, got %v", resultValue)
	}
}

// Test ReturnError
func TestLuaHelpers_ReturnError(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	message := "something went wrong"

	count := helpers.ReturnError(L, message)

	if count != 2 {
		t.Errorf("Expected ReturnError to return 2 values, got %d", count)
	}

	// Check stack: should have nil result and error message
	if L.GetTop() != 2 {
		t.Errorf("Expected 2 values on stack, got %d", L.GetTop())
	}

	resultValue := L.Get(-2)
	if resultValue.Type() != lua.LTNil {
		t.Errorf("Expected nil result, got %s", resultValue.Type())
	}

	errorValue := L.Get(-1)
	if errorValue.String() != message {
		t.Errorf("Expected error message '%s', got '%s'", message, errorValue.String())
	}
}

func TestLuaHelpers_ReturnError_EmptyMessage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	helpers.ReturnError(L, "")

	errorValue := L.Get(-1)
	if errorValue.String() != "" {
		t.Errorf("Expected empty error message, got '%s'", errorValue.String())
	}
}

func TestLuaHelpers_ReturnError_LongMessage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	longMessage := "This is a very long error message that contains detailed information about what went wrong"

	helpers.ReturnError(L, longMessage)

	errorValue := L.Get(-1)
	if errorValue.String() != longMessage {
		t.Error("Expected long error message to be preserved")
	}
}

// Test ReturnFluentSuccess
func TestLuaHelpers_ReturnFluentSuccess(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	self := L.NewUserData()

	count := helpers.ReturnFluentSuccess(L, self)

	if count != 2 {
		t.Errorf("Expected ReturnFluentSuccess to return 2 values, got %d", count)
	}

	selfValue := L.Get(-2)
	if selfValue.Type() != lua.LTUserData {
		t.Errorf("Expected UserData self, got %s", selfValue.Type())
	}

	errorValue := L.Get(-1)
	if errorValue.Type() != lua.LTNil {
		t.Errorf("Expected nil error, got %s", errorValue.Type())
	}
}

func TestLuaHelpers_ReturnFluentSuccess_WithTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	self := L.NewTable()

	helpers.ReturnFluentSuccess(L, self)

	selfValue := L.Get(-2)
	if selfValue.Type() != lua.LTTable {
		t.Errorf("Expected table, got %s", selfValue.Type())
	}
}

// Test CreateResultTable
func TestLuaHelpers_CreateResultTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	result := helpers.CreateResultTable(L, true, "Operation successful", nil)

	if result == nil {
		t.Fatal("Expected non-nil result table")
	}

	changedValue := L.GetField(result, "changed")
	if changedValue != lua.LTrue {
		t.Error("Expected changed to be true")
	}

	messageValue := L.GetField(result, "message")
	if messageValue.String() != "Operation successful" {
		t.Errorf("Expected message 'Operation successful', got '%s'", messageValue.String())
	}
}

func TestLuaHelpers_CreateResultTable_NotChanged(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	result := helpers.CreateResultTable(L, false, "No changes needed", nil)

	changedValue := L.GetField(result, "changed")
	if changedValue != lua.LFalse {
		t.Error("Expected changed to be false")
	}
}

func TestLuaHelpers_CreateResultTable_WithFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	fields := map[string]lua.LValue{
		"path":  lua.LString("/tmp/file"),
		"count": lua.LNumber(42),
		"valid": lua.LTrue,
	}

	result := helpers.CreateResultTable(L, true, "Created", fields)

	pathValue := L.GetField(result, "path")
	if pathValue.String() != "/tmp/file" {
		t.Error("Expected path field to be set")
	}

	countValue := L.GetField(result, "count")
	if countValue != lua.LNumber(42) {
		t.Error("Expected count field to be set")
	}

	validValue := L.GetField(result, "valid")
	if validValue != lua.LTrue {
		t.Error("Expected valid field to be set")
	}
}

func TestLuaHelpers_CreateResultTable_EmptyFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	result := helpers.CreateResultTable(L, true, "OK", map[string]lua.LValue{})

	// Should still have changed and message
	changedValue := L.GetField(result, "changed")
	if changedValue == lua.LNil {
		t.Error("Expected changed field to exist")
	}

	messageValue := L.GetField(result, "message")
	if messageValue == lua.LNil {
		t.Error("Expected message field to exist")
	}
}

// Test ReturnIdempotent
func TestLuaHelpers_ReturnIdempotent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	count := helpers.ReturnIdempotent(L, "Resource already exists")

	if count != 2 {
		t.Errorf("Expected 2 return values, got %d", count)
	}

	result := L.Get(-2)
	if result.Type() != lua.LTTable {
		t.Errorf("Expected table result, got %s", result.Type())
	}

	resultTable := result.(*lua.LTable)
	changedValue := L.GetField(resultTable, "changed")
	if changedValue != lua.LFalse {
		t.Error("Expected idempotent result to have changed=false")
	}

	messageValue := L.GetField(resultTable, "message")
	if messageValue.String() != "Resource already exists" {
		t.Errorf("Expected message to be preserved")
	}

	errorValue := L.Get(-1)
	if errorValue.Type() != lua.LTNil {
		t.Error("Expected nil error")
	}
}

func TestLuaHelpers_ReturnIdempotent_EmptyMessage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	helpers.ReturnIdempotent(L, "")

	result := L.Get(-2).(*lua.LTable)
	messageValue := L.GetField(result, "message")
	if messageValue.String() != "" {
		t.Error("Expected empty message")
	}
}

// Test ReturnChanged
func TestLuaHelpers_ReturnChanged(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	fields := map[string]lua.LValue{
		"id": lua.LString("123"),
	}

	count := helpers.ReturnChanged(L, "Resource created", fields)

	if count != 2 {
		t.Errorf("Expected 2 return values, got %d", count)
	}

	result := L.Get(-2)
	if result.Type() != lua.LTTable {
		t.Errorf("Expected table result, got %s", result.Type())
	}

	resultTable := result.(*lua.LTable)
	changedValue := L.GetField(resultTable, "changed")
	if changedValue != lua.LTrue {
		t.Error("Expected changed=true")
	}

	messageValue := L.GetField(resultTable, "message")
	if messageValue.String() != "Resource created" {
		t.Error("Expected message to be preserved")
	}

	idValue := L.GetField(resultTable, "id")
	if idValue.String() != "123" {
		t.Error("Expected id field to be set")
	}

	errorValue := L.Get(-1)
	if errorValue.Type() != lua.LTNil {
		t.Error("Expected nil error")
	}
}

func TestLuaHelpers_ReturnChanged_NoFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	helpers.ReturnChanged(L, "Updated", nil)

	result := L.Get(-2).(*lua.LTable)
	changedValue := L.GetField(result, "changed")
	if changedValue != lua.LTrue {
		t.Error("Expected changed=true")
	}
}

func TestLuaHelpers_ReturnChanged_MultipleFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}
	fields := map[string]lua.LValue{
		"id":       lua.LString("abc"),
		"name":     lua.LString("test"),
		"count":    lua.LNumber(10),
		"enabled":  lua.LTrue,
		"disabled": lua.LFalse,
	}

	helpers.ReturnChanged(L, "Created", fields)

	result := L.Get(-2).(*lua.LTable)

	// Check all fields are present
	if L.GetField(result, "id").String() != "abc" {
		t.Error("Expected id field")
	}
	if L.GetField(result, "name").String() != "test" {
		t.Error("Expected name field")
	}
	if L.GetField(result, "count") != lua.LNumber(10) {
		t.Error("Expected count field")
	}
	if L.GetField(result, "enabled") != lua.LTrue {
		t.Error("Expected enabled field")
	}
	if L.GetField(result, "disabled") != lua.LFalse {
		t.Error("Expected disabled field")
	}
}

// Test GetStringField
func TestGetStringField_ValidString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "name", lua.LString("test"))

	result := GetStringField(L, table, "name", "default")

	if result != "test" {
		t.Errorf("Expected 'test', got '%s'", result)
	}
}

func TestGetStringField_MissingField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()

	result := GetStringField(L, table, "missing", "default")

	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}

func TestGetStringField_WrongType(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "number", lua.LNumber(42))

	result := GetStringField(L, table, "number", "default")

	if result != "default" {
		t.Errorf("Expected default value for wrong type, got '%s'", result)
	}
}

func TestGetStringField_EmptyString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "empty", lua.LString(""))

	result := GetStringField(L, table, "empty", "default")

	if result != "" {
		t.Errorf("Expected empty string to be returned, got '%s'", result)
	}
}

func TestGetStringField_LongString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	longString := "This is a very long string with lots of text"
	L.SetField(table, "long", lua.LString(longString))

	result := GetStringField(L, table, "long", "default")

	if result != longString {
		t.Error("Expected long string to be preserved")
	}
}

// Test GetBoolField
func TestGetBoolField_ValidTrue(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "enabled", lua.LTrue)

	result := GetBoolField(L, table, "enabled", false)

	if result != true {
		t.Error("Expected true")
	}
}

func TestGetBoolField_ValidFalse(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "disabled", lua.LFalse)

	result := GetBoolField(L, table, "disabled", true)

	if result != false {
		t.Error("Expected false")
	}
}

func TestGetBoolField_MissingField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()

	result := GetBoolField(L, table, "missing", true)

	if result != true {
		t.Error("Expected default value true")
	}
}

func TestGetBoolField_WrongType(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "string", lua.LString("true"))

	result := GetBoolField(L, table, "string", false)

	if result != false {
		t.Error("Expected default value for wrong type")
	}
}

func TestGetBoolField_NumberType(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "number", lua.LNumber(1))

	result := GetBoolField(L, table, "number", true)

	if result != true {
		t.Error("Expected default value when field is not boolean")
	}
}

// Test GetIntField
func TestGetIntField_ValidInt(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "count", lua.LNumber(42))

	result := GetIntField(L, table, "count", 0)

	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestGetIntField_Zero(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "zero", lua.LNumber(0))

	result := GetIntField(L, table, "zero", 99)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestGetIntField_NegativeNumber(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "negative", lua.LNumber(-10))

	result := GetIntField(L, table, "negative", 0)

	if result != -10 {
		t.Errorf("Expected -10, got %d", result)
	}
}

func TestGetIntField_MissingField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()

	result := GetIntField(L, table, "missing", 100)

	if result != 100 {
		t.Errorf("Expected default value 100, got %d", result)
	}
}

func TestGetIntField_WrongType(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "string", lua.LString("42"))

	result := GetIntField(L, table, "string", 0)

	if result != 0 {
		t.Error("Expected default value for wrong type")
	}
}

func TestGetIntField_LargeNumber(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "large", lua.LNumber(999999))

	result := GetIntField(L, table, "large", 0)

	if result != 999999 {
		t.Errorf("Expected 999999, got %d", result)
	}
}

// Test GetTableField
func TestGetTableField_ValidTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	nested := L.NewTable()
	L.SetField(nested, "key", lua.LString("value"))
	L.SetField(table, "nested", nested)

	result := GetTableField(L, table, "nested")

	if result == nil {
		t.Error("Expected non-nil table")
	}

	value := L.GetField(result, "key")
	if value.String() != "value" {
		t.Error("Expected nested table to be returned")
	}
}

func TestGetTableField_MissingField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()

	result := GetTableField(L, table, "missing")

	if result != nil {
		t.Error("Expected nil for missing field")
	}
}

func TestGetTableField_WrongType(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "string", lua.LString("not a table"))

	result := GetTableField(L, table, "string")

	if result != nil {
		t.Error("Expected nil for wrong type")
	}
}

func TestGetTableField_EmptyTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	empty := L.NewTable()
	L.SetField(table, "empty", empty)

	result := GetTableField(L, table, "empty")

	if result == nil {
		t.Error("Expected empty table to be returned")
	}
}

func TestGetTableField_NumberField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	L.SetField(table, "number", lua.LNumber(123))

	result := GetTableField(L, table, "number")

	if result != nil {
		t.Error("Expected nil when field is not a table")
	}
}

// Test edge cases and combinations
func TestLuaHelpers_ReturnSuccessAndErrorSequence(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}

	// Return success first
	helpers.ReturnSuccess(L, lua.LString("ok"))

	// Return error next
	helpers.ReturnError(L, "failed")

	// Stack should have 4 values total
	if L.GetTop() != 4 {
		t.Errorf("Expected 4 values on stack, got %d", L.GetTop())
	}

	// Last two should be error return (nil, "failed")
	resultValue := L.Get(-2)
	errorValue := L.Get(-1)

	if resultValue.Type() != lua.LTNil {
		t.Error("Expected nil result for error return")
	}

	if errorValue.String() != "failed" {
		t.Errorf("Expected error message 'failed', got '%s'", errorValue.String())
	}
}

func TestLuaHelpers_CreateResultTable_ComplexFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	helpers := &LuaHelpers{}

	// Create nested table
	nestedTable := L.NewTable()
	L.SetField(nestedTable, "nested_key", lua.LString("nested_value"))

	// Create array table
	arrayTable := L.NewTable()
	arrayTable.Append(lua.LString("item1"))
	arrayTable.Append(lua.LString("item2"))

	fields := map[string]lua.LValue{
		"nested": nestedTable,
		"array":  arrayTable,
		"nil":    lua.LNil,
	}

	result := helpers.CreateResultTable(L, true, "Complex result", fields)

	// Verify nested table
	nested := L.GetField(result, "nested")
	if nested.Type() != lua.LTTable {
		t.Error("Expected nested table field")
	}

	nestedValue := L.GetField(nested.(*lua.LTable), "nested_key")
	if nestedValue.String() != "nested_value" {
		t.Error("Expected nested value to be preserved")
	}

	// Verify array table
	array := L.GetField(result, "array")
	if array.Type() != lua.LTTable {
		t.Error("Expected array table field")
	}

	// Verify nil field
	nilField := L.GetField(result, "nil")
	if nilField.Type() != lua.LTNil {
		t.Error("Expected nil field to be stored")
	}
}
