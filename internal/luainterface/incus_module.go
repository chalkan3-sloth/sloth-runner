package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterIncusModule registers the Incus module into the Lua state
func RegisterIncusModule(L *lua.LState) {
	// Get the global agent client if available
	var agentClient interface{}
	if globalCore := L.GetGlobal("_agent_client"); globalCore != lua.LNil {
		agentClient = globalCore
	}
	
	// Create and register the Incus module
	incusModule := infra.NewIncusModule(nil) // Pass nil for now, agent will be resolved later
	incusModule.Register(L)
	
	// Register metatables
	infra.RegisterInstanceMetatable(L)
	infra.RegisterImageMetatable(L)
	infra.RegisterNetworkMetatable(L)
	infra.RegisterProfileMetatable(L)
	infra.RegisterStorageMetatable(L)
	infra.RegisterSnapshotMetatable(L)
	
	_ = agentClient // Mark as used
}
