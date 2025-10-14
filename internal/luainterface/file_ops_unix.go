//go:build unix || darwin || linux
// +build unix darwin linux

package luainterface

import (
	"os"
	"syscall"

	lua "github.com/yuin/gopher-lua"
)

// addUnixFileInfo adds Unix-specific file information (UID/GID) to the result table
func addUnixFileInfo(L *lua.LState, result *lua.LTable, info os.FileInfo) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		L.SetField(result, "uid", lua.LNumber(stat.Uid))
		L.SetField(result, "gid", lua.LNumber(stat.Gid))
	}
}
