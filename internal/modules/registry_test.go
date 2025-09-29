package modules

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

// TestModule implements ModuleLoader for testing
type TestModule struct {
	*BaseModule
	initialized bool
	cleaned     bool
}

func NewTestModule() *TestModule {
	info := ModuleInfo{
		Name:        "test",
		Version:     "1.0.0",
		Description: "Test module",
		Author:      "Test",
		Category:    "test",
		Dependencies: []string{},
	}
	
	return &TestModule{
		BaseModule: NewBaseModule(info),
	}
}

func (m *TestModule) Loader(L *lua.LState) int {
	testTable := L.NewTable()
	L.SetFuncs(testTable, map[string]lua.LGFunction{
		"hello": func(L *lua.LState) int {
			L.Push(lua.LString("Hello from test module!"))
			return 1
		},
	})
	L.Push(testTable)
	return 1
}

func (m *TestModule) Initialize() error {
	m.initialized = true
	return nil
}

func (m *TestModule) Cleanup() error {
	m.cleaned = true
	return nil
}

func TestModuleRegistry(t *testing.T) {
	registry := NewModuleRegistry()
	testModule := NewTestModule()
	
	// Test registration
	err := registry.Register("test", testModule)
	if err != nil {
		t.Fatalf("Failed to register module: %v", err)
	}
	
	// Test module was initialized
	if !testModule.initialized {
		t.Error("Module was not initialized during registration")
	}
	
	// Test retrieval
	loader, exists := registry.Get("test")
	if !exists {
		t.Error("Module not found after registration")
	}
	
	if loader != testModule {
		t.Error("Retrieved module is not the same as registered")
	}
	
	// Test listing
	modules := registry.List()
	if len(modules) != 1 || modules[0] != "test" {
		t.Errorf("Expected ['test'], got %v", modules)
	}
	
	// Test info
	info := registry.GetInfo()
	if len(info) != 1 {
		t.Errorf("Expected 1 module info, got %d", len(info))
	}
	
	testInfo := info["test"]
	if testInfo.Name != "test" || testInfo.Version != "1.0.0" {
		t.Errorf("Module info mismatch: %+v", testInfo)
	}
	
	// Test unregistration
	err = registry.Unregister("test")
	if err != nil {
		t.Fatalf("Failed to unregister module: %v", err)
	}
	
	// Test module was cleaned up
	if !testModule.cleaned {
		t.Error("Module was not cleaned up during unregistration")
	}
	
	// Test module no longer exists
	_, exists = registry.Get("test")
	if exists {
		t.Error("Module still exists after unregistration")
	}
}

func TestModuleLoading(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	registry := NewModuleRegistry()
	testModule := NewTestModule()
	
	err := registry.Register("test", testModule)
	if err != nil {
		t.Fatalf("Failed to register module: %v", err)
	}
	
	// Load module into Lua state
	err = registry.LoadModule(L, "test")
	if err != nil {
		t.Fatalf("Failed to load module: %v", err)
	}
	
	// Test the module in Lua
	script := `
		local test = require("test")
		result = test.hello()
	`
	
	err = L.DoString(script)
	if err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}
	
	// Check result
	result := L.GetGlobal("result")
	if result.Type() != lua.LTString || result.String() != "Hello from test module!" {
		t.Errorf("Unexpected result: %v", result)
	}
}

func TestValidationHelpers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	// Create test table
	table := L.NewTable()
	table.RawSetString("name", lua.LString("John"))
	table.RawSetString("age", lua.LString("25"))
	table.RawSetString("email", lua.LString(""))
	
	// Test required field validation
	required := []string{"name", "age", "email"}
	validation := ValidateRequired(L, table, required)
	
	if validation.IsValid {
		t.Error("Validation should have failed for empty email")
	}
	
	if len(validation.Errors) != 1 || validation.Errors[0] != "missing required field: email" {
		t.Errorf("Unexpected validation errors: %v", validation.Errors)
	}
	
	// Fix the table and test again
	table.RawSetString("email", lua.LString("john@example.com"))
	validation = ValidateRequired(L, table, required)
	
	if !validation.IsValid {
		t.Errorf("Validation should have passed: %v", validation.Errors)
	}
}

func TestResponseHelpers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	// Test error response
	CreateErrorResponse(L, "test error", "detail1", "detail2")
	
	result := L.Get(-1).(*lua.LTable)
	success := result.RawGetString("success")
	if success.Type() != lua.LTBool || bool(success.(lua.LBool)) != false {
		t.Error("Error response should have success=false")
	}
	
	errorMsg := result.RawGetString("error")
	if errorMsg.String() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", errorMsg.String())
	}
	
	L.Pop(1)
	
	// Test success response
	data := lua.LString("test data")
	CreateSuccessResponse(L, data)
	
	result = L.Get(-1).(*lua.LTable)
	success = result.RawGetString("success")
	if success.Type() != lua.LTBool || bool(success.(lua.LBool)) != true {
		t.Error("Success response should have success=true")
	}
	
	resultData := result.RawGetString("data")
	if resultData.String() != "test data" {
		t.Errorf("Expected 'test data', got '%s'", resultData.String())
	}
}

func TestCategoryFiltering(t *testing.T) {
	registry := NewModuleRegistry()
	
	// Register modules with different categories
	testModule1 := NewTestModule()
	testModule1.info.Category = "core"
	registry.Register("test1", testModule1)
	
	testModule2 := NewTestModule()
	testModule2.info.Category = "cloud"
	registry.Register("test2", testModule2)
	
	testModule3 := NewTestModule()
	testModule3.info.Category = "core"
	registry.Register("test3", testModule3)
	
	// Test category filtering
	coreModules := registry.ListByCategory("core")
	if len(coreModules) != 2 {
		t.Errorf("Expected 2 core modules, got %d", len(coreModules))
	}
	
	cloudModules := registry.ListByCategory("cloud")
	if len(cloudModules) != 1 || cloudModules[0] != "test2" {
		t.Errorf("Expected ['test2'], got %v", cloudModules)
	}
	
	nonExistentModules := registry.ListByCategory("nonexistent")
	if len(nonExistentModules) != 0 {
		t.Errorf("Expected no modules for nonexistent category, got %v", nonExistentModules)
	}
}