package core

import (
	"fmt"
	"sort"

	"github.com/yuin/gopher-lua"
)

// ModuleInfo represents module metadata
type ModuleInfo struct {
	Name         string
	Version      string
	Description  string
	Author       string
	Functions    []string
	Examples     []string
	Categories   []string
	Category     string
	Dependencies []string
}

// BaseModule provides common functionality for all modules
type BaseModule struct {
	info ModuleInfo
}

// NewBaseModule creates a new base module
func NewBaseModule(info ModuleInfo) *BaseModule {
	return &BaseModule{info: info}
}

// Info returns module information
func (bm *BaseModule) Info() ModuleInfo {
	return bm.info
}

// HelpModule provides interactive help and documentation
type HelpModule struct {
	*BaseModule
}

// NewHelpModule creates a new help module
func NewHelpModule() *HelpModule {
	info := ModuleInfo{
		Name:        "help",
		Version:     "1.0.0", 
		Description: "Interactive help system for Sloth Runner modules",
		Author:      "Sloth Runner Team",
		Category:    "core",
		Dependencies: []string{},
		Functions: []string{
			"help([topic]) - Show general help or help for specific topic",
			"modules() - List all available modules",
			"search(term) - Search modules and functions",
			"examples([category]) - Show usage examples",
		},
		Examples: []string{
			"help() -- Show general help",
			"help('http') -- Show HTTP module help", 
			"help.modules() -- List all modules",
			"help.search('request') -- Search for request-related functions",
		},
	}
	
	return &HelpModule{
		BaseModule: NewBaseModule(info),
	}
}

// Load registers the help module functions
func (hm *HelpModule) Load(L *lua.LState) int {
	helpTable := L.NewTable()
	
	// Register functions
	L.SetFuncs(helpTable, map[string]lua.LGFunction{
		"help":     hm.showHelp,
		"modules":  hm.listModules, 
		"search":   hm.searchHelp,
		"examples": hm.showExamples,
	})
	
	// Set as global
	L.SetGlobal("help", helpTable)
	L.Push(helpTable)
	return 1
}

func (hm *HelpModule) showHelp(L *lua.LState) int {
	topicArg := L.OptString(1, "")
	
	if topicArg == "" {
		help := `
🦥 Sloth Runner Help System

Available Commands:
  help()           - Show this help
  help.modules()   - List all available modules  
  help.search()    - Search for functions
  help.examples()  - Show usage examples

Core Modules:
  • help     - This help system
  • http     - HTTP client functionality
  • validate - Data validation utilities

For module-specific help: help('module_name')
`
		L.Push(lua.LString(help))
	} else {
		help := fmt.Sprintf("Help for '%s' module - Implementation coming soon", topicArg)
		L.Push(lua.LString(help))
	}
	return 1
}

func (hm *HelpModule) listModules(L *lua.LState) int {
	modules := []string{"help", "http", "validate"}
	sort.Strings(modules)
	
	result := fmt.Sprintf("📚 Available Modules (%d total):\n\n", len(modules))
	for _, name := range modules {
		result += fmt.Sprintf("  • %s\n", name)
	}
	
	L.Push(lua.LString(result))
	return 1
}

func (hm *HelpModule) searchHelp(L *lua.LState) int {
	query := L.CheckString(1)
	
	result := fmt.Sprintf("🔍 Search results for '%s':\n\nFound matches in: help, http, validate modules", query)
	L.Push(lua.LString(result))
	return 1
}

func (hm *HelpModule) showExamples(L *lua.LState) int {
	categoryArg := L.OptString(1, "")
	
	examples := `
📝 Sloth Runner Examples

HTTP Module:
  response = http.get("https://api.example.com")
  data = http.post("https://api.example.com", {key = "value"})

Validation Module:
  valid = validate.email("test@example.com")
  clean = validate.sanitize(user_input)

Help System:
  help()                    -- Show general help
  help.modules()           -- List modules
  help.search("http")      -- Search functions
`
	
	if categoryArg != "" {
		examples = fmt.Sprintf("Examples for category '%s':\n%s", categoryArg, examples)
	}
	
	L.Push(lua.LString(examples))
	return 1
}