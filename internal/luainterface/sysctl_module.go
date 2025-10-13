package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterSysctlModule registers the Sysctl module into the Lua state
func RegisterSysctlModule(L *lua.LState) {
	sysctlModule := infra.NewSysctlModule(L)
	sysctlModule.Register(L)
}
