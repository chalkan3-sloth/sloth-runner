package core

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestNewHelpModule(t *testing.T) {
	hm := NewHelpModule()
	
	if hm == nil {
		t.Fatal("Expected non-nil help module")
	}

	info := hm.Info()
	if info.Name != "help" {
		t.Errorf("Expected name 'help', got '%s'", info.Name)
	}

	if info.Version == "" {
		t.Error("Expected non-empty version")
	}

	if info.Description == "" {
		t.Error("Expected non-empty description")
	}
}

func TestHelpModule_Info(t *testing.T) {
	hm := NewHelpModule()
	info := hm.Info()

	// Verify all fields
	if info.Name == "" {
		t.Error("Expected non-empty name")
	}

	if info.Version == "" {
		t.Error("Expected non-empty version")
	}

	if info.Description == "" {
		t.Error("Expected non-empty description")
	}

	if info.Author == "" {
		t.Error("Expected non-empty author")
	}

	if info.Category == "" {
		t.Error("Expected non-empty category")
	}

	if info.Dependencies == nil {
		t.Error("Expected non-nil dependencies slice")
	}
}

func TestHelpModule_Load(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	returned := hm.Load(L)

	if returned != 1 {
		t.Errorf("Expected Load to return 1, got %d", returned)
	}

	// Check that help global exists
	helpGlobal := L.GetGlobal("help")
	if helpGlobal.Type() != lua.LTTable {
		t.Errorf("Expected help global to be a table, got %s", helpGlobal.Type())
	}

	// Check that functions exist in the table
	helpTable := helpGlobal.(*lua.LTable)
	
	functions := []string{"help", "modules", "search", "examples"}
	for _, fn := range functions {
		val := L.GetField(helpTable, fn)
		if val.Type() != lua.LTFunction {
			t.Errorf("Expected %s to be a function, got %s", fn, val.Type())
		}
	}
}

func TestHelpModule_ShowHelp_NoArgs(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	hm.Load(L)

	// Call help.help()
	if err := L.DoString(`result = help.help()`); err != nil {
		t.Fatalf("Failed to call help.help(): %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTString {
		t.Errorf("Expected string result, got %s", result.Type())
	}

	helpText := lua.LVAsString(result)
	if !strings.Contains(helpText, "Sloth Runner") {
		t.Error("Expected help text to contain 'Sloth Runner'")
	}

	if !strings.Contains(helpText, "help()") {
		t.Error("Expected help text to contain function documentation")
	}
}

func TestHelpModule_ShowHelp_WithTopic(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	hm.Load(L)

	// Call help.help("http")
	if err := L.DoString(`result = help.help("http")`); err != nil {
		t.Fatalf("Failed to call help.help('http'): %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTString {
		t.Errorf("Expected string result, got %s", result.Type())
	}

	helpText := lua.LVAsString(result)
	if !strings.Contains(helpText, "http") {
		t.Error("Expected help text to contain 'http'")
	}
}

func TestHelpModule_ListModules(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	hm.Load(L)

	// Call help.modules()
	if err := L.DoString(`result = help.modules()`); err != nil {
		t.Fatalf("Failed to call help.modules(): %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTString {
		t.Errorf("Expected string result, got %s", result.Type())
	}

	modulesText := lua.LVAsString(result)
	if !strings.Contains(modulesText, "Available Modules") {
		t.Error("Expected modules text to contain 'Available Modules'")
	}

	// Should contain some core modules
	if !strings.Contains(modulesText, "help") {
		t.Error("Expected modules list to contain 'help'")
	}
}

func TestHelpModule_SearchHelp(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	hm.Load(L)

	// Call help.search("http")
	if err := L.DoString(`result = help.search("http")`); err != nil {
		t.Fatalf("Failed to call help.search('http'): %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTString {
		t.Errorf("Expected string result, got %s", result.Type())
	}

	searchText := lua.LVAsString(result)
	if !strings.Contains(searchText, "Search results") {
		t.Error("Expected search text to contain 'Search results'")
	}

	if !strings.Contains(searchText, "http") {
		t.Error("Expected search text to contain query term")
	}
}

func TestHelpModule_ShowExamples_NoArgs(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	hm.Load(L)

	// Call help.examples()
	if err := L.DoString(`result = help.examples()`); err != nil {
		t.Fatalf("Failed to call help.examples(): %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTString {
		t.Errorf("Expected string result, got %s", result.Type())
	}

	examplesText := lua.LVAsString(result)
	if !strings.Contains(examplesText, "Examples") {
		t.Error("Expected examples text to contain 'Examples'")
	}

	if !strings.Contains(examplesText, "http.get") {
		t.Error("Expected examples text to contain code examples")
	}
}

func TestHelpModule_ShowExamples_WithCategory(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	hm.Load(L)

	// Call help.examples("http")
	if err := L.DoString(`result = help.examples("http")`); err != nil {
		t.Fatalf("Failed to call help.examples('http'): %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTString {
		t.Errorf("Expected string result, got %s", result.Type())
	}

	examplesText := lua.LVAsString(result)
	if !strings.Contains(examplesText, "http") {
		t.Error("Expected examples text to contain category name")
	}
}

func TestHelpModule_Loader(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	hm := NewHelpModule()
	returned := hm.Loader(L)

	if returned != 1 {
		t.Errorf("Expected Loader to return 1, got %d", returned)
	}
}
