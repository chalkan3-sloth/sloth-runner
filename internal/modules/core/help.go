package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chalkan3/sloth-runner/internal/modules"
	"github.com/yuin/gopher-lua"
)

// HelpModule provides interactive help and documentation
type HelpModule struct {
	*modules.BaseModule
}

// NewHelpModule creates a new help module
func NewHelpModule() *HelpModule {
	info := modules.ModuleInfo{
		Name:        "help",
		Version:     "1.0.0",
		Description: "Interactive help system for modules and functions",
		Author:      "Sloth Runner Team",
		Category:    "core",
		Dependencies: []string{},
	}
	
	return &HelpModule{
		BaseModule: modules.NewBaseModule(info),
	}
}

// Loader returns the Lua loader function
func (m *HelpModule) Loader(L *lua.LState) int {
	helpTable := L.NewTable()
	
	L.SetFuncs(helpTable, map[string]lua.LGFunction{
		"modules":   m.luaListModules,
		"module":    m.luaModuleHelp,
		"search":    m.luaSearchHelp,
		"examples":  m.luaExamples,
		"version":   m.luaVersion,
	})
	
	// Add global help function
	L.SetGlobal("help", L.NewFunction(m.luaGlobalHelp))
	
	L.Push(helpTable)
	return 1
}

// luaListModules lists all available modules
func (m *HelpModule) luaListModules(L *lua.LState) int {
	registry := modules.GetGlobalRegistry()
	moduleList := registry.List()
	sort.Strings(moduleList)
	
	fmt.Printf("📚 Available Modules (%d total):\n\n", len(moduleList))
	
	// Group by category
	categories := make(map[string][]string)
	allInfo := registry.GetInfo()
	
	for _, name := range moduleList {
		info := allInfo[name]
		category := info.Category
		if category == "" {
			category = "other"
		}
		categories[category] = append(categories[category], name)
	}
	
	// Display by category
	for category, names := range categories {
		fmt.Printf("🔹 %s:\n", strings.Title(category))
		for _, name := range names {
			info := allInfo[name]
			fmt.Printf("  • %-15s - %s\n", name, info.Description)
		}
		fmt.Println()
	}
	
	fmt.Println("💡 Use help.module('name') for detailed information about a specific module")
	fmt.Println("🔍 Use help.search('term') to search for functionality")
	
	return 0
}

// luaModuleHelp shows detailed help for a specific module
func (m *HelpModule) luaModuleHelp(L *lua.LState) int {
	moduleName := L.CheckString(1)
	
	registry := modules.GetGlobalRegistry()
	loader, exists := registry.Get(moduleName)
	if !exists {
		fmt.Printf("❌ Module '%s' not found\n", moduleName)
		fmt.Println("💡 Use help.modules() to see all available modules")
		return 0
	}
	
	info := loader.Info()
	
	fmt.Printf("📖 Module: %s\n", info.Name)
	fmt.Printf("🔖 Version: %s\n", info.Version)
	fmt.Printf("📝 Description: %s\n", info.Description)
	fmt.Printf("👤 Author: %s\n", info.Author)
	fmt.Printf("🏷️  Category: %s\n", info.Category)
	
	if len(info.Dependencies) > 0 {
		fmt.Printf("📦 Dependencies: %s\n", strings.Join(info.Dependencies, ", "))
	}
	
	// Show usage examples based on module type
	m.showModuleExamples(moduleName)
	
	return 0
}

// luaSearchHelp searches for modules or functions containing a term
func (m *HelpModule) luaSearchHelp(L *lua.LState) int {
	searchTerm := strings.ToLower(L.CheckString(1))
	
	registry := modules.GetGlobalRegistry()
	allInfo := registry.GetInfo()
	
	fmt.Printf("🔍 Search results for '%s':\n\n", searchTerm)
	
	var found bool
	for name, info := range allInfo {
		if strings.Contains(strings.ToLower(name), searchTerm) ||
		   strings.Contains(strings.ToLower(info.Description), searchTerm) ||
		   strings.Contains(strings.ToLower(info.Category), searchTerm) {
			fmt.Printf("✨ %s - %s\n", name, info.Description)
			found = true
		}
	}
	
	if !found {
		fmt.Println("❌ No modules found matching your search")
		fmt.Println("💡 Try a broader search term or use help.modules() to see all available modules")
	}
	
	return 0
}

// luaExamples shows examples for common tasks
func (m *HelpModule) luaExamples(L *lua.LState) int {
	category := L.OptString(1, "")
	
	examples := m.getExamples()
	
	if category == "" {
		fmt.Println("📚 Available Example Categories:")
		for cat := range examples {
			fmt.Printf("  • %s\n", cat)
		}
		fmt.Println("\n💡 Use help.examples('category') to see examples for a specific category")
		return 0
	}
	
	if categoryExamples, exists := examples[category]; exists {
		fmt.Printf("📖 Examples for '%s':\n\n", category)
		for i, example := range categoryExamples {
			fmt.Printf("%d. %s:\n", i+1, example.Title)
			fmt.Printf("```lua\n%s\n```\n\n", example.Code)
		}
	} else {
		fmt.Printf("❌ No examples found for category '%s'\n", category)
	}
	
	return 0
}

// luaVersion shows version information
func (m *HelpModule) luaVersion(L *lua.LState) int {
	fmt.Println("🦥 Sloth Runner Help System")
	fmt.Println("Version: 1.0.0")
	fmt.Println("Built with ❤️  by the Sloth Runner Team")
	
	return 0
}

// luaGlobalHelp provides context-sensitive help
func (m *HelpModule) luaGlobalHelp(L *lua.LState) int {
	if L.GetTop() == 0 {
		// Show general help
		fmt.Println("🆘 Sloth Runner Help")
		fmt.Println("==================")
		fmt.Println()
		fmt.Println("Available commands:")
		fmt.Println("  help()               - Show this help")
		fmt.Println("  help.modules()       - List all available modules")
		fmt.Println("  help.module('name')  - Show help for specific module")
		fmt.Println("  help.search('term')  - Search for functionality")
		fmt.Println("  help.examples()      - Show example categories")
		fmt.Println("  help.version()       - Show version information")
		fmt.Println()
		fmt.Println("💡 Quick Start:")
		fmt.Println("  1. Run 'help.modules()' to see available modules")
		fmt.Println("  2. Use 'local mod = require(\"module_name\")' to load a module")
		fmt.Println("  3. Check 'help.module(\"module_name\")' for detailed usage")
		return 0
	}
	
	// Show help for specific topic
	topic := L.CheckString(1)
	return m.luaModuleHelp(L)
}

// Example represents a code example
type Example struct {
	Title string
	Code  string
}

// getExamples returns categorized examples
func (m *HelpModule) getExamples() map[string][]Example {
	return map[string][]Example{
		"http": {
			{
				Title: "Simple GET request",
				Code: `local http = require("http")
local result = http.get({
    url = "https://api.github.com/user",
    headers = {
        ["Authorization"] = "token YOUR_TOKEN"
    }
})
if result.success then
    print("Status:", result.data.status_code)
    print("Body:", result.data.body)
end`,
			},
			{
				Title: "POST with JSON data",
				Code: `local http = require("http")
local result = http.post({
    url = "https://api.example.com/data",
    json = {
        name = "John Doe",
        email = "john@example.com"
    },
    headers = {
        ["Authorization"] = "Bearer token"
    }
})`,
			},
		},
		"docker": {
			{
				Title: "List running containers",
				Code: `local docker = require("docker")
local result = docker.exec({"ps"})
if result.success then
    print(result.stdout)
end`,
			},
		},
		"state": {
			{
				Title: "Store and retrieve data",
				Code: `local state = require("state")
-- Store data
state.set("user_config", {
    theme = "dark",
    language = "en"
})

-- Retrieve data
local config = state.get("user_config")
print("Theme:", config.theme)`,
			},
		},
	}
}

// showModuleExamples shows specific examples for a module
func (m *HelpModule) showModuleExamples(moduleName string) {
	examples := m.getExamples()
	if moduleExamples, exists := examples[moduleName]; exists {
		fmt.Printf("\n📋 Examples:\n")
		for i, example := range moduleExamples {
			fmt.Printf("\n%d. %s:\n", i+1, example.Title)
			fmt.Printf("```lua\n%s\n```\n", example.Code)
		}
	} else {
		fmt.Printf("\n💡 No specific examples available for this module yet.\n")
		fmt.Printf("Check the documentation or contribute examples!\n")
	}
}