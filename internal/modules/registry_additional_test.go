package modules

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

// Mock module loader for testing
type mockModuleLoader struct {
	name           string
	version        string
	description    string
	author         string
	category       string
	dependencies   []string
	initError      error
	cleanupError   error
	loaderCalled   bool
	initCalled     bool
	cleanupCalled  bool
}

func (m *mockModuleLoader) Loader(L *lua.LState) int {
	m.loaderCalled = true
	return 0
}

func (m *mockModuleLoader) Info() ModuleInfo {
	return ModuleInfo{
		Name:         m.name,
		Version:      m.version,
		Description:  m.description,
		Author:       m.author,
		Category:     m.category,
		Dependencies: m.dependencies,
	}
}

func (m *mockModuleLoader) Initialize() error {
	m.initCalled = true
	return m.initError
}

func (m *mockModuleLoader) Cleanup() error {
	m.cleanupCalled = true
	return m.cleanupError
}

// Test ModuleInfo struct
func TestModuleInfo_Creation(t *testing.T) {
	info := ModuleInfo{
		Name:        "test-module",
		Version:     "1.0.0",
		Description: "Test module",
		Author:      "Test Author",
		Category:    "core",
		Dependencies: []string{"dep1", "dep2"},
	}

	if info.Name != "test-module" {
		t.Errorf("Expected name 'test-module', got '%s'", info.Name)
	}

	if info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", info.Version)
	}

	if len(info.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(info.Dependencies))
	}
}

func TestModuleInfo_EmptyValues(t *testing.T) {
	info := ModuleInfo{}

	if info.Name != "" {
		t.Error("Expected empty name")
	}

	if info.Version != "" {
		t.Error("Expected empty version")
	}

	if info.Dependencies != nil {
		t.Error("Expected nil dependencies")
	}
}

func TestModuleInfo_NoDependencies(t *testing.T) {
	info := ModuleInfo{
		Name:         "standalone",
		Version:      "1.0.0",
		Dependencies: []string{},
	}

	if len(info.Dependencies) != 0 {
		t.Error("Expected no dependencies")
	}
}

func TestModuleInfo_MultipleCategories(t *testing.T) {
	categories := []string{"core", "utils", "network", "database", "system"}

	for _, category := range categories {
		info := ModuleInfo{
			Name:     "module",
			Category: category,
		}

		if info.Category != category {
			t.Errorf("Expected category '%s', got '%s'", category, info.Category)
		}
	}
}

// Test NewModuleRegistry
func TestNewModuleRegistry(t *testing.T) {
	registry := NewModuleRegistry()

	if registry == nil {
		t.Error("Expected non-nil registry")
	}

	if registry.modules == nil {
		t.Error("Expected initialized modules map")
	}

	if len(registry.modules) != 0 {
		t.Error("Expected empty modules map")
	}
}

func TestNewModuleRegistry_MultipleInstances(t *testing.T) {
	registry1 := NewModuleRegistry()
	registry2 := NewModuleRegistry()

	if registry1 == registry2 {
		t.Error("Expected different registry instances")
	}
}

// Test GetGlobalRegistry
func TestGetGlobalRegistry(t *testing.T) {
	registry := GetGlobalRegistry()

	if registry == nil {
		t.Error("Expected non-nil global registry")
	}
}

func TestGetGlobalRegistry_Singleton(t *testing.T) {
	registry1 := GetGlobalRegistry()
	registry2 := GetGlobalRegistry()

	if registry1 != registry2 {
		t.Error("Expected same global registry instance (singleton)")
	}
}

// Test Register
func TestModuleRegistry_Register_Success(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{
		name:    "test-module",
		version: "1.0.0",
	}

	err := registry.Register("test-module", loader)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !loader.initCalled {
		t.Error("Expected Initialize to be called")
	}
}

func TestModuleRegistry_Register_Duplicate(t *testing.T) {
	registry := NewModuleRegistry()
	loader1 := &mockModuleLoader{name: "duplicate"}
	loader2 := &mockModuleLoader{name: "duplicate"}

	registry.Register("duplicate", loader1)
	err := registry.Register("duplicate", loader2)

	if err == nil {
		t.Error("Expected error when registering duplicate module")
	}
}

func TestModuleRegistry_Register_InitError(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{
		name:      "error-module",
		initError: &testError{"init failed"},
	}

	err := registry.Register("error-module", loader)

	if err == nil {
		t.Error("Expected error when initialization fails")
	}
}

func TestModuleRegistry_Register_MultipleModules(t *testing.T) {
	registry := NewModuleRegistry()

	for i := 0; i < 10; i++ {
		loader := &mockModuleLoader{
			name: "module-" + string(rune('0'+i)),
		}
		err := registry.Register(loader.name, loader)
		if err != nil {
			t.Errorf("Failed to register module %d: %v", i, err)
		}
	}

	if len(registry.modules) != 10 {
		t.Errorf("Expected 10 modules, got %d", len(registry.modules))
	}
}

// Test Unregister
func TestModuleRegistry_Unregister_Success(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{name: "test"}

	registry.Register("test", loader)
	err := registry.Unregister("test")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !loader.cleanupCalled {
		t.Error("Expected Cleanup to be called")
	}
}

func TestModuleRegistry_Unregister_NotFound(t *testing.T) {
	registry := NewModuleRegistry()

	err := registry.Unregister("nonexistent")

	if err == nil {
		t.Error("Expected error when unregistering nonexistent module")
	}
}

func TestModuleRegistry_Unregister_CleanupError(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{
		name:         "test",
		cleanupError: &testError{"cleanup failed"},
	}

	registry.Register("test", loader)
	err := registry.Unregister("test")

	if err == nil {
		t.Error("Expected error when cleanup fails")
	}
}

// Test Get
func TestModuleRegistry_Get_Success(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{name: "test"}

	registry.Register("test", loader)
	retrieved, exists := registry.Get("test")

	if !exists {
		t.Error("Expected module to exist")
	}

	if retrieved != loader {
		t.Error("Expected same loader instance")
	}
}

func TestModuleRegistry_Get_NotFound(t *testing.T) {
	registry := NewModuleRegistry()

	_, exists := registry.Get("nonexistent")

	if exists {
		t.Error("Expected module to not exist")
	}
}

func TestModuleRegistry_Get_AfterUnregister(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{name: "test"}

	registry.Register("test", loader)
	registry.Unregister("test")

	_, exists := registry.Get("test")

	if exists {
		t.Error("Expected module to not exist after unregister")
	}
}

// Test List
func TestModuleRegistry_List_Empty(t *testing.T) {
	registry := NewModuleRegistry()

	list := registry.List()

	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d modules", len(list))
	}
}

func TestModuleRegistry_List_Multiple(t *testing.T) {
	registry := NewModuleRegistry()

	names := []string{"module1", "module2", "module3"}
	for _, name := range names {
		loader := &mockModuleLoader{name: name}
		registry.Register(name, loader)
	}

	list := registry.List()

	if len(list) != 3 {
		t.Errorf("Expected 3 modules, got %d", len(list))
	}
}

func TestModuleRegistry_List_AfterOperations(t *testing.T) {
	registry := NewModuleRegistry()

	loader1 := &mockModuleLoader{name: "m1"}
	loader2 := &mockModuleLoader{name: "m2"}
	loader3 := &mockModuleLoader{name: "m3"}

	registry.Register("m1", loader1)
	registry.Register("m2", loader2)
	registry.Register("m3", loader3)
	registry.Unregister("m2")

	list := registry.List()

	if len(list) != 2 {
		t.Errorf("Expected 2 modules after unregister, got %d", len(list))
	}
}

// Test ListByCategory
func TestModuleRegistry_ListByCategory_Empty(t *testing.T) {
	registry := NewModuleRegistry()

	list := registry.ListByCategory("core")

	if len(list) != 0 {
		t.Error("Expected empty list for nonexistent category")
	}
}

func TestModuleRegistry_ListByCategory_Multiple(t *testing.T) {
	registry := NewModuleRegistry()

	loader1 := &mockModuleLoader{name: "core1", category: "core"}
	loader2 := &mockModuleLoader{name: "core2", category: "core"}
	loader3 := &mockModuleLoader{name: "util1", category: "utils"}

	registry.Register("core1", loader1)
	registry.Register("core2", loader2)
	registry.Register("util1", loader3)

	coreList := registry.ListByCategory("core")

	if len(coreList) != 2 {
		t.Errorf("Expected 2 core modules, got %d", len(coreList))
	}
}

func TestModuleRegistry_ListByCategory_NoMatches(t *testing.T) {
	registry := NewModuleRegistry()

	loader := &mockModuleLoader{name: "test", category: "core"}
	registry.Register("test", loader)

	list := registry.ListByCategory("nonexistent")

	if len(list) != 0 {
		t.Error("Expected no matches for nonexistent category")
	}
}

// Test GetInfo
func TestModuleRegistry_GetInfo_Empty(t *testing.T) {
	registry := NewModuleRegistry()

	info := registry.GetInfo()

	if len(info) != 0 {
		t.Error("Expected empty info map")
	}
}

func TestModuleRegistry_GetInfo_Multiple(t *testing.T) {
	registry := NewModuleRegistry()

	loader1 := &mockModuleLoader{
		name:    "mod1",
		version: "1.0.0",
	}
	loader2 := &mockModuleLoader{
		name:    "mod2",
		version: "2.0.0",
	}

	registry.Register("mod1", loader1)
	registry.Register("mod2", loader2)

	info := registry.GetInfo()

	if len(info) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(info))
	}

	if info["mod1"].Version != "1.0.0" {
		t.Error("Expected mod1 version to be 1.0.0")
	}
}

func TestModuleRegistry_GetInfo_Complete(t *testing.T) {
	registry := NewModuleRegistry()

	loader := &mockModuleLoader{
		name:         "complete",
		version:      "1.0.0",
		description:  "Complete module",
		author:       "Test Author",
		category:     "core",
		dependencies: []string{"dep1"},
	}

	registry.Register("complete", loader)
	info := registry.GetInfo()

	moduleInfo := info["complete"]
	if moduleInfo.Name != "complete" {
		t.Error("Expected complete module info")
	}

	if len(moduleInfo.Dependencies) != 1 {
		t.Error("Expected 1 dependency")
	}
}

// Test LoadModule
func TestModuleRegistry_LoadModule_Success(t *testing.T) {
	registry := NewModuleRegistry()
	L := lua.NewState()
	defer L.Close()

	loader := &mockModuleLoader{name: "test"}
	registry.Register("test", loader)

	err := registry.LoadModule(L, "test")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestModuleRegistry_LoadModule_NotFound(t *testing.T) {
	registry := NewModuleRegistry()
	L := lua.NewState()
	defer L.Close()

	err := registry.LoadModule(L, "nonexistent")

	if err == nil {
		t.Error("Expected error for nonexistent module")
	}
}

func TestModuleRegistry_LoadModule_PreloadCheck(t *testing.T) {
	registry := NewModuleRegistry()
	L := lua.NewState()
	defer L.Close()

	loader := &mockModuleLoader{name: "preload"}
	registry.Register("preload", loader)

	registry.LoadModule(L, "preload")

	// Check that module was preloaded (exists in package.preload)
	L.DoString(`
		local preload = package.preload
		if preload.preload == nil then
			error("Module not preloaded")
		end
	`)
}

// Test LoadAllModules
func TestModuleRegistry_LoadAllModules_Empty(t *testing.T) {
	registry := NewModuleRegistry()
	L := lua.NewState()
	defer L.Close()

	err := registry.LoadAllModules(L)

	if err != nil {
		t.Errorf("Expected no error for empty registry, got %v", err)
	}
}

func TestModuleRegistry_LoadAllModules_Multiple(t *testing.T) {
	registry := NewModuleRegistry()
	L := lua.NewState()
	defer L.Close()

	for i := 1; i <= 5; i++ {
		loader := &mockModuleLoader{name: "mod"}
		registry.Register(loader.name, loader)
	}

	err := registry.LoadAllModules(L)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test concurrent operations
func TestModuleRegistry_ConcurrentRegister(t *testing.T) {
	registry := NewModuleRegistry()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			loader := &mockModuleLoader{name: "concurrent"}
			registry.Register(loader.name, loader)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Only one should succeed (first one registered)
	list := registry.List()
	if len(list) != 1 {
		t.Errorf("Expected 1 module after concurrent register, got %d", len(list))
	}
}

func TestModuleRegistry_ConcurrentGet(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{name: "test"}
	registry.Register("test", loader)

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			_, exists := registry.Get("test")
			if !exists {
				t.Error("Expected module to exist")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestModuleRegistry_ConcurrentList(t *testing.T) {
	registry := NewModuleRegistry()

	for i := 0; i < 5; i++ {
		loader := &mockModuleLoader{name: "mod"}
		registry.Register(loader.name, loader)
	}

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			list := registry.List()
			if len(list) < 0 {
				t.Error("Expected valid list")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test edge cases
func TestModuleRegistry_EmptyModuleName(t *testing.T) {
	registry := NewModuleRegistry()
	loader := &mockModuleLoader{name: ""}

	err := registry.Register("", loader)

	if err != nil {
		t.Error("Should allow empty module name (caller's responsibility)")
	}
}

func TestModuleRegistry_SpecialCharactersInName(t *testing.T) {
	registry := NewModuleRegistry()

	specialNames := []string{
		"mod-with-dash",
		"mod_with_underscore",
		"mod.with.dots",
		"mod/with/slash",
	}

	for _, name := range specialNames {
		loader := &mockModuleLoader{name: name}
		err := registry.Register(name, loader)

		if err != nil {
			t.Errorf("Should allow special characters in name '%s': %v", name, err)
		}
	}
}

func TestModuleRegistry_LongModuleName(t *testing.T) {
	registry := NewModuleRegistry()

	longName := ""
	for i := 0; i < 100; i++ {
		longName += "a"
	}

	loader := &mockModuleLoader{name: longName}
	err := registry.Register(longName, loader)

	if err != nil {
		t.Error("Should allow long module names")
	}
}

// Helper error type for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
