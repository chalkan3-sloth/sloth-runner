package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterRAIDModule registers the RAID module into the Lua state
func RegisterRAIDModule(L *lua.LState) {
	raidModule := infra.NewRAIDModule(L)
	raidModule.Register(L)
}
