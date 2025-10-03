package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/core"
	lua "github.com/yuin/gopher-lua"
)

// RegisterStowModule registers the Stow module into the Lua state
func RegisterStowModule(L *lua.LState) {
	core.RegisterStowModule(L)
}
