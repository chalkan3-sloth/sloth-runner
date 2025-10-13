package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterLVMModule registers the LVM module into the Lua state
func RegisterLVMModule(L *lua.LState) {
	lvmModule := infra.NewLVMModule(L)
	lvmModule.Register(L)
}
