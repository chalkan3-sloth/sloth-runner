package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterFirewallModule registers the firewall module into the Lua state
func RegisterFirewallModule(L *lua.LState) {
	// Get the global agent client if available
	var agentClient interface{}
	if globalCore := L.GetGlobal("_agent_client"); globalCore != lua.LNil {
		agentClient = globalCore
	}

	// Create and register the firewall module
	firewallModule := infra.NewFirewallModule(nil) // Pass nil for now, agent will be resolved later
	firewallModule.Register(L)

	// Register metatable for firewall rules
	infra.RegisterRuleMetatable(L)

	_ = agentClient // Mark as used
}
