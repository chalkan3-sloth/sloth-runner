package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterNixOSModule registers the NixOS module for Lua
func RegisterNixOSModule(L *lua.LState) {
	infra.RegisterNixOSModule(L)
}
