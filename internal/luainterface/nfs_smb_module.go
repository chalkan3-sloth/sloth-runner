package luainterface

import (
	"github.com/chalkan3-sloth/sloth-runner/internal/modules/infra"
	lua "github.com/yuin/gopher-lua"
)

// RegisterNFSSMBModule registers the NFS/SMB module into the Lua state
func RegisterNFSSMBModule(L *lua.LState) {
	nfsSmbModule := infra.NewNFSSMBModule(L)
	nfsSmbModule.Register(L)
}
