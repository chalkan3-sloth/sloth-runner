package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/core"
	lua "github.com/yuin/gopher-lua"
)

// RegisterGitModule registers the Git module into the Lua state
func RegisterGitModule(L *lua.LState) {
	core.RegisterGitModule(L)
}
