package modules

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

// Test NewBaseModule
func TestNewBaseModule(t *testing.T) {
	info := ModuleInfo{
		Name:    "test",
		Version: "1.0.0",
	}

	module := NewBaseModule(info)

	if module == nil {
		t.Error("Expected non-nil module")
	}

	if module.info.Name != "test" {
		t.Error("Expected module info to be set")
	}
}

func TestNewBaseModule_EmptyInfo(t *testing.T) {
	module := NewBaseModule(ModuleInfo{})

	if module == nil {
		t.Error("Expected non-nil module")
	}

	if module.info.Name != "" {
		t.Error("Expected empty name")
	}
}

func TestNewBaseModule_CompleteInfo(t *testing.T) {
	info := ModuleInfo{
		Name:         "complete",
		Version:      "2.0.0",
		Description:  "Complete module",
		Author:       "Test Author",
		Category:     "core",
		Dependencies: []string{"dep1", "dep2"},
	}

	module := NewBaseModule(info)

	if len(module.info.Dependencies) != 2 {
		t.Error("Expected dependencies to be preserved")
	}
}

// Test BaseModule Info
func TestBaseModule_Info(t *testing.T) {
	info := ModuleInfo{
		Name:    "test",
		Version: "1.0.0",
	}

	module := NewBaseModule(info)
	retrieved := module.Info()

	if retrieved.Name != "test" {
		t.Error("Expected to retrieve correct info")
	}
}

func TestBaseModule_Info_Immutability(t *testing.T) {
	info := ModuleInfo{
		Name:    "original",
		Version: "1.0.0",
	}

	module := NewBaseModule(info)
	retrieved := module.Info()

	// Modifying retrieved should not affect module
	retrieved.Name = "modified"

	if module.Info().Name == "modified" {
		t.Error("Module info should not be affected by modifications to retrieved copy")
	}
}

// Test BaseModule Initialize
func TestBaseModule_Initialize(t *testing.T) {
	module := NewBaseModule(ModuleInfo{})

	err := module.Initialize()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBaseModule_Initialize_MultipleCalls(t *testing.T) {
	module := NewBaseModule(ModuleInfo{})

	err1 := module.Initialize()
	err2 := module.Initialize()
	err3 := module.Initialize()

	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Expected no errors on multiple initialize calls")
	}
}

// Test BaseModule Cleanup
func TestBaseModule_Cleanup(t *testing.T) {
	module := NewBaseModule(ModuleInfo{})

	err := module.Cleanup()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBaseModule_Cleanup_MultipleCalls(t *testing.T) {
	module := NewBaseModule(ModuleInfo{})

	err1 := module.Cleanup()
	err2 := module.Cleanup()

	if err1 != nil || err2 != nil {
		t.Error("Expected no errors on multiple cleanup calls")
	}
}

// Test ValidationResult
func TestValidationResult_Valid(t *testing.T) {
	result := ValidationResult{
		IsValid: true,
		Errors:  []string{},
	}

	if !result.IsValid {
		t.Error("Expected valid result")
	}

	if len(result.Errors) != 0 {
		t.Error("Expected no errors")
	}
}

func TestValidationResult_Invalid(t *testing.T) {
	result := ValidationResult{
		IsValid: false,
		Errors:  []string{"error1", "error2"},
	}

	if result.IsValid {
		t.Error("Expected invalid result")
	}

	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationResult_NoErrors(t *testing.T) {
	result := ValidationResult{
		IsValid: true,
	}

	if result.Errors != nil {
		t.Error("Expected nil errors slice")
	}
}

// Test ValidateRequired
func TestValidateRequired_AllPresent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("name", lua.LString("test"))
	table.RawSetString("version", lua.LString("1.0.0"))

	result := ValidateRequired(L, table, []string{"name", "version"})

	if !result.IsValid {
		t.Error("Expected valid result when all required fields present")
	}

	if len(result.Errors) != 0 {
		t.Error("Expected no errors")
	}
}

func TestValidateRequired_MissingField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("name", lua.LString("test"))

	result := ValidateRequired(L, table, []string{"name", "version"})

	if result.IsValid {
		t.Error("Expected invalid result when required field missing")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
}

func TestValidateRequired_EmptyString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("name", lua.LString(""))

	result := ValidateRequired(L, table, []string{"name"})

	if result.IsValid {
		t.Error("Expected invalid result for empty string")
	}
}

func TestValidateRequired_NoRequiredFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()

	result := ValidateRequired(L, table, []string{})

	if !result.IsValid {
		t.Error("Expected valid result when no fields required")
	}
}

func TestValidateRequired_MultipleErrors(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()

	result := ValidateRequired(L, table, []string{"field1", "field2", "field3"})

	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(result.Errors))
	}
}

func TestValidateRequired_NilValue(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("test", lua.LNil)

	result := ValidateRequired(L, table, []string{"test"})

	if result.IsValid {
		t.Error("Expected invalid result for nil value")
	}
}

// Test CreateErrorResponse
func TestCreateErrorResponse_Basic(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	CreateErrorResponse(L, "test error")

	result := L.Get(-1).(*lua.LTable)
	success := result.RawGetString("success")
	errorMsg := result.RawGetString("error")

	if success != lua.LFalse {
		t.Error("Expected success to be false")
	}

	if errorMsg.String() != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", errorMsg.String())
	}
}

func TestCreateErrorResponse_WithDetails(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	CreateErrorResponse(L, "validation failed", "missing name", "missing version")

	result := L.Get(-1).(*lua.LTable)
	details := result.RawGetString("details").(*lua.LTable)

	detail1 := details.RawGetInt(1)
	detail2 := details.RawGetInt(2)

	if detail1.String() != "missing name" {
		t.Error("Expected first detail to be 'missing name'")
	}

	if detail2.String() != "missing version" {
		t.Error("Expected second detail to be 'missing version'")
	}
}

func TestCreateErrorResponse_NoDetails(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	CreateErrorResponse(L, "simple error")

	result := L.Get(-1).(*lua.LTable)
	details := result.RawGetString("details")

	if details != lua.LNil {
		t.Error("Expected no details field when no details provided")
	}
}

func TestCreateErrorResponse_EmptyMessage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	CreateErrorResponse(L, "")

	result := L.Get(-1).(*lua.LTable)
	errorMsg := result.RawGetString("error")

	if errorMsg.String() != "" {
		t.Error("Expected empty error message to be preserved")
	}
}

func TestCreateErrorResponse_StackPosition(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	returnValue := CreateErrorResponse(L, "test")

	if returnValue != 1 {
		t.Errorf("Expected return value 1, got %d", returnValue)
	}

	if L.GetTop() != 1 {
		t.Errorf("Expected 1 value on stack, got %d", L.GetTop())
	}
}

// Test CreateSuccessResponse
func TestCreateSuccessResponse_WithData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data := L.NewTable()
	data.RawSetString("result", lua.LString("success"))

	CreateSuccessResponse(L, data)

	result := L.Get(-1).(*lua.LTable)
	success := result.RawGetString("success")
	resultData := result.RawGetString("data")

	if success != lua.LTrue {
		t.Error("Expected success to be true")
	}

	if resultData == lua.LNil {
		t.Error("Expected data to be present")
	}
}

func TestCreateSuccessResponse_NoData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	CreateSuccessResponse(L, nil)

	result := L.Get(-1).(*lua.LTable)
	data := result.RawGetString("data")

	if data != lua.LNil {
		t.Error("Expected no data field when nil provided")
	}
}

func TestCreateSuccessResponse_StringData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data := lua.LString("test value")

	CreateSuccessResponse(L, data)

	result := L.Get(-1).(*lua.LTable)
	resultData := result.RawGetString("data")

	if resultData.String() != "test value" {
		t.Error("Expected data to be 'test value'")
	}
}

func TestCreateSuccessResponse_NumberData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data := lua.LNumber(42)

	CreateSuccessResponse(L, data)

	result := L.Get(-1).(*lua.LTable)
	resultData := result.RawGetString("data")

	if resultData.(lua.LNumber) != 42 {
		t.Error("Expected data to be 42")
	}
}

func TestCreateSuccessResponse_TableData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data := L.NewTable()
	data.RawSetString("key", lua.LString("value"))

	CreateSuccessResponse(L, data)

	result := L.Get(-1).(*lua.LTable)
	resultData := result.RawGetString("data").(*lua.LTable)
	value := resultData.RawGetString("key")

	if value.String() != "value" {
		t.Error("Expected nested data to be preserved")
	}
}

func TestCreateSuccessResponse_StackPosition(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	returnValue := CreateSuccessResponse(L, nil)

	if returnValue != 1 {
		t.Errorf("Expected return value 1, got %d", returnValue)
	}

	if L.GetTop() != 1 {
		t.Errorf("Expected 1 value on stack, got %d", L.GetTop())
	}
}

// Test WrapLuaFunction
func TestWrapLuaFunction_NoValidation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	called := false
	fn := func(L *lua.LState) int {
		called = true
		return 0
	}

	wrapped := WrapLuaFunction(fn, []string{})
	wrapped(L)

	if !called {
		t.Error("Expected wrapped function to be called")
	}
}

func TestWrapLuaFunction_WithValidation_Success(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	called := false
	fn := func(L *lua.LState) int {
		called = true
		return 0
	}

	table := L.NewTable()
	table.RawSetString("name", lua.LString("test"))
	L.Push(table)

	wrapped := WrapLuaFunction(fn, []string{"name"})
	wrapped(L)

	if !called {
		t.Error("Expected wrapped function to be called when validation passes")
	}
}

func TestWrapLuaFunction_WithValidation_Failure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	called := false
	fn := func(L *lua.LState) int {
		called = true
		return 0
	}

	table := L.NewTable()
	L.Push(table)

	wrapped := WrapLuaFunction(fn, []string{"required_field"})
	wrapped(L)

	if called {
		t.Error("Expected wrapped function not to be called when validation fails")
	}

	// Check error response
	result := L.Get(-1).(*lua.LTable)
	success := result.RawGetString("success")

	if success != lua.LFalse {
		t.Error("Expected error response on validation failure")
	}
}

func TestWrapLuaFunction_NoTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	called := false
	fn := func(L *lua.LState) int {
		called = true
		return 0
	}

	L.Push(lua.LString("not a table"))

	wrapped := WrapLuaFunction(fn, []string{"field"})
	wrapped(L)

	if !called {
		t.Error("Expected wrapped function to be called when no table provided (validation skipped)")
	}
}

func TestWrapLuaFunction_ReturnValue(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	fn := func(L *lua.LState) int {
		L.Push(lua.LString("result"))
		return 1
	}

	wrapped := WrapLuaFunction(fn, []string{})
	returnCount := wrapped(L)

	if returnCount != 1 {
		t.Errorf("Expected return count 1, got %d", returnCount)
	}

	result := L.Get(-1)
	if result.String() != "result" {
		t.Error("Expected result to be preserved")
	}
}

// Test BaseModule lifecycle
func TestBaseModule_Lifecycle(t *testing.T) {
	module := NewBaseModule(ModuleInfo{Name: "test"})

	// Initialize
	err := module.Initialize()
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}

	// Use module (get info)
	info := module.Info()
	if info.Name != "test" {
		t.Error("Expected module to work after initialization")
	}

	// Cleanup
	err = module.Cleanup()
	if err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}
}

func TestBaseModule_MultipleLifecycles(t *testing.T) {
	module := NewBaseModule(ModuleInfo{})

	for i := 0; i < 3; i++ {
		module.Initialize()
		module.Info()
		module.Cleanup()
	}

	// Should still work after multiple cycles
	info := module.Info()
	if info.Name != "" {
		t.Error("Expected module to remain functional")
	}
}

// Test ValidationResult with various scenarios
func TestValidationResult_SingleError(t *testing.T) {
	result := ValidationResult{
		IsValid: false,
		Errors:  []string{"single error"},
	}

	if len(result.Errors) != 1 {
		t.Error("Expected single error")
	}
}

func TestValidationResult_MultipleErrors(t *testing.T) {
	result := ValidationResult{
		IsValid: false,
		Errors:  []string{"error1", "error2", "error3"},
	}

	if len(result.Errors) != 3 {
		t.Error("Expected 3 errors")
	}
}

func TestValidationResult_EmptyErrors(t *testing.T) {
	result := ValidationResult{
		IsValid: true,
		Errors:  []string{},
	}

	if len(result.Errors) != 0 {
		t.Error("Expected empty errors slice")
	}
}

// Test edge cases
func TestValidateRequired_ComplexTypes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()

	// Number field
	table.RawSetString("number", lua.LNumber(42))

	// Bool field
	table.RawSetString("bool", lua.LBool(true))

	// Table field
	nestedTable := L.NewTable()
	table.RawSetString("table", nestedTable)

	result := ValidateRequired(L, table, []string{"number", "bool", "table"})

	if !result.IsValid {
		t.Error("Expected validation to pass for complex types")
	}
}

func TestValidateRequired_NumberZero(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("count", lua.LNumber(0))

	result := ValidateRequired(L, table, []string{"count"})

	if !result.IsValid {
		t.Error("Expected zero number to be valid")
	}
}

func TestValidateRequired_BoolFalse(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("flag", lua.LBool(false))

	result := ValidateRequired(L, table, []string{"flag"})

	if !result.IsValid {
		t.Error("Expected false boolean to be valid")
	}
}

func TestCreateErrorResponse_ManyDetails(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	details := make([]string, 10)
	for i := 0; i < 10; i++ {
		details[i] = "error"
	}

	CreateErrorResponse(L, "many errors", details...)

	result := L.Get(-1).(*lua.LTable)
	detailsTable := result.RawGetString("details").(*lua.LTable)

	// Check we have all details
	count := 0
	detailsTable.ForEach(func(k, v lua.LValue) {
		count++
	})

	if count != 10 {
		t.Errorf("Expected 10 details, got %d", count)
	}
}

func TestCreateSuccessResponse_BoolData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data := lua.LBool(true)

	CreateSuccessResponse(L, data)

	result := L.Get(-1).(*lua.LTable)
	resultData := result.RawGetString("data")

	if resultData.(lua.LBool) != lua.LTrue {
		t.Error("Expected boolean data to be preserved")
	}
}

// TestValidateRequired_ConcurrentAccess tests concurrent validation calls
func TestValidateRequired_ConcurrentAccess(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("name", lua.LString("test"))
	table.RawSetString("version", lua.LString("1.0.0"))

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			result := ValidateRequired(L, table, []string{"name", "version"})
			if !result.IsValid {
				t.Error("Expected valid result")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
