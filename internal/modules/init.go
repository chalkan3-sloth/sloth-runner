package modules

import (
	"log/slog"

	"github.com/chalkan3/sloth-runner/internal/modules/core"
	"github.com/yuin/gopher-lua"
)

// InitializeCoreModules registers all core modules
func InitializeCoreModules(registry *ModuleRegistry) error {
	coreModules := []ModuleLoader{
		core.NewHTTPModule(),
		core.NewHelpModule(), 
		core.NewValidateModule(),
	}
	
	for _, module := range coreModules {
		info := module.Info()
		if err := registry.Register(info.Name, module); err != nil {
			slog.Error("Failed to register core module", "module", info.Name, "error", err)
			return err
		}
		slog.Info("Registered core module", "module", info.Name, "version", info.Version)
	}
	
	return nil
}

// SetupLuaEnvironment sets up a Lua state with all registered modules
func SetupLuaEnvironment() (*lua.LState, error) {
	L := lua.NewState()
	
	// Initialize core modules if not already done
	registry := GetGlobalRegistry()
	if len(registry.List()) == 0 {
		if err := InitializeCoreModules(registry); err != nil {
			L.Close()
			return nil, err
		}
	}
	
	// Load all modules into Lua state
	if err := registry.LoadAllModules(L); err != nil {
		L.Close()
		return nil, err
	}
	
	// Add global functions
	setupGlobalFunctions(L, registry)
	
	slog.Info("Lua environment setup complete", "modules_loaded", len(registry.List()))
	return L, nil
}

// setupGlobalFunctions adds convenience global functions
func setupGlobalFunctions(L *lua.LState, registry *ModuleRegistry) {
	// Add a global function to list available modules
	L.SetGlobal("modules", L.NewFunction(func(L *lua.LState) int {
		moduleList := registry.List()
		table := L.NewTable()
		for i, name := range moduleList {
			table.RawSetInt(i+1, lua.LString(name))
		}
		L.Push(table)
		return 1
	}))
	
	// Add a global function to get module info
	L.SetGlobal("module_info", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		loader, exists := registry.Get(name)
		if !exists {
			L.Push(lua.LNil)
			return 1
		}
		
		info := loader.Info()
		table := L.NewTable()
		table.RawSetString("name", lua.LString(info.Name))
		table.RawSetString("version", lua.LString(info.Version))
		table.RawSetString("description", lua.LString(info.Description))
		table.RawSetString("author", lua.LString(info.Author))
		table.RawSetString("category", lua.LString(info.Category))
		
		// Add dependencies
		if len(info.Dependencies) > 0 {
			deps := L.NewTable()
			for i, dep := range info.Dependencies {
				deps.RawSetInt(i+1, lua.LString(dep))
			}
			table.RawSetString("dependencies", deps)
		}
		
		L.Push(table)
		return 1
	}))
}