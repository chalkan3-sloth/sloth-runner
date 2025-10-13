package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterCronModule registers the Cron module into the Lua state
func RegisterCronModule(L *lua.LState) {
	cronModule := infra.NewCronModule(L)
	cronModule.Register(L)
}
