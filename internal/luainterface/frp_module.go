package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterFrpModule registers the FRP module into the Lua state
func RegisterFrpModule(L *lua.LState) {
	// Get the global agent client if available
	var agentClient interface{}
	if globalCore := L.GetGlobal("_agent_client"); globalCore != lua.LNil {
		agentClient = globalCore
	}

	// Create and register the FRP module
	frpModule := infra.NewFrpModule(nil) // Pass nil for now, agent will be resolved later
	frpModule.Register(L)

	// Register metatables for FRP server and client
	infra.RegisterServerMetatable(L)
	infra.RegisterClientMetatable(L)

	_ = agentClient // Mark as used
}
