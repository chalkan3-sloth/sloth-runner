package modules

import (
	"fmt"
	"sync"

	"github.com/yuin/gopher-lua"
)

// ModuleInfo contains metadata about a module
type ModuleInfo struct {
	Name        string
	Version     string
	Description string
	Author      string
	Dependencies []string
	Category    string
}

// ModuleLoader is the interface that all modules must implement
type ModuleLoader interface {
	// Loader returns the Lua loader function for the module
	Loader(L *lua.LState) int
	
	// Info returns metadata about the module
	Info() ModuleInfo
	
	// Initialize is called when the module is first loaded
	Initialize() error
	
	// Cleanup is called when the module is being unloaded
	Cleanup() error
}

// ModuleRegistry manages all registered modules
type ModuleRegistry struct {
	modules map[string]ModuleLoader
	mutex   sync.RWMutex
}

var (
	globalRegistry *ModuleRegistry
	registryOnce   sync.Once
)

// GetGlobalRegistry returns the singleton module registry
func GetGlobalRegistry() *ModuleRegistry {
	registryOnce.Do(func() {
		globalRegistry = NewModuleRegistry()
	})
	return globalRegistry
}

// NewModuleRegistry creates a new module registry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]ModuleLoader),
	}
}

// Register adds a module to the registry
func (r *ModuleRegistry) Register(name string, loader ModuleLoader) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module '%s' is already registered", name)
	}
	
	// Initialize the module
	if err := loader.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize module '%s': %w", name, err)
	}
	
	r.modules[name] = loader
	return nil
}

// Unregister removes a module from the registry
func (r *ModuleRegistry) Unregister(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	loader, exists := r.modules[name]
	if !exists {
		return fmt.Errorf("module '%s' is not registered", name)
	}
	
	// Cleanup the module
	if err := loader.Cleanup(); err != nil {
		return fmt.Errorf("failed to cleanup module '%s': %w", name, err)
	}
	
	delete(r.modules, name)
	return nil
}

// Get retrieves a module from the registry
func (r *ModuleRegistry) Get(name string) (ModuleLoader, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	loader, exists := r.modules[name]
	return loader, exists
}

// List returns all registered modules
func (r *ModuleRegistry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}

// ListByCategory returns modules filtered by category
func (r *ModuleRegistry) ListByCategory(category string) []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var names []string
	for name, loader := range r.modules {
		if loader.Info().Category == category {
			names = append(names, name)
		}
	}
	return names
}

// GetInfo returns information about all registered modules
func (r *ModuleRegistry) GetInfo() map[string]ModuleInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	info := make(map[string]ModuleInfo)
	for name, loader := range r.modules {
		info[name] = loader.Info()
	}
	return info
}

// LoadModule loads a module into the Lua state
func (r *ModuleRegistry) LoadModule(L *lua.LState, name string) error {
	loader, exists := r.Get(name)
	if !exists {
		return fmt.Errorf("module '%s' not found", name)
	}
	
	// Create a require function that loads the module
	requireFn := func(L *lua.LState) int {
		return loader.Loader(L)
	}
	
	// Register the module with Lua's require system
	L.PreloadModule(name, requireFn)
	return nil
}

// LoadAllModules loads all registered modules into the Lua state
func (r *ModuleRegistry) LoadAllModules(L *lua.LState) error {
	names := r.List()
	for _, name := range names {
		if err := r.LoadModule(L, name); err != nil {
			return fmt.Errorf("failed to load module '%s': %w", name, err)
		}
	}
	return nil
}