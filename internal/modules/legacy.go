package modules

import (
	"github.com/chalkan3/sloth-runner/internal/luainterface"
	"github.com/chalkan3/sloth-runner/internal/modules"
	"github.com/yuin/gopher-lua"
)

// LegacyModule wraps existing luainterface modules to work with the new system
type LegacyModule struct {
	*modules.BaseModule
	loaderFunc func(*lua.LState) int
}

// NewLegacyModule creates a module wrapper for existing code
func NewLegacyModule(info modules.ModuleInfo, loaderFunc func(*lua.LState) int) *LegacyModule {
	return &LegacyModule{
		BaseModule: modules.NewBaseModule(info),
		loaderFunc: loaderFunc,
	}
}

// Loader implements the ModuleLoader interface
func (m *LegacyModule) Loader(L *lua.LState) int {
	return m.loaderFunc(L)
}

// MigrateLegacyModules registers existing modules with the new system
func MigrateLegacyModules(registry *modules.ModuleRegistry) error {
	// Example: Migrate Docker module
	dockerInfo := modules.ModuleInfo{
		Name:        "docker",
		Version:     "1.0.0",
		Description: "Docker container management",
		Author:      "Sloth Runner Team",
		Category:    "devops",
		Dependencies: []string{},
	}
	
	dockerModule := NewLegacyModule(dockerInfo, func(L *lua.LState) int {
		dockerMod := luainterface.NewDockerModule()
		return dockerMod.Loader(L)
	})
	
	if err := registry.Register("docker", dockerModule); err != nil {
		return err
	}
	
	// Example: Migrate State module
	stateInfo := modules.ModuleInfo{
		Name:        "state",
		Version:     "1.0.0", 
		Description: "Persistent state management with SQLite backend",
		Author:      "Sloth Runner Team",
		Category:    "core",
		Dependencies: []string{},
	}
	
	stateModule := NewLegacyModule(stateInfo, func(L *lua.LState) int {
		stateMod := luainterface.GetGlobalStateModule()
		return stateMod.Loader(L)
	})
	
	if err := registry.Register("state", stateModule); err != nil {
		return err
	}
	
	// Example: Migrate AWS module
	awsInfo := modules.ModuleInfo{
		Name:        "aws",
		Version:     "1.0.0",
		Description: "Amazon Web Services integration",
		Author:      "Sloth Runner Team", 
		Category:    "cloud",
		Dependencies: []string{"aws-cli"},
	}
	
	awsModule := NewLegacyModule(awsInfo, func(L *lua.LState) int {
		// The old AWS module doesn't have a struct, so we need to recreate the loader
		awsTable := L.NewTable()
		// Add AWS functions here based on existing implementation
		L.Push(awsTable)
		return 1
	})
	
	if err := registry.Register("aws", awsModule); err != nil {
		return err
	}
	
	// Add more legacy modules as needed...
	
	return nil
}