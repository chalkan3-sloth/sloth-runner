//go:build windows
// +build windows

package luainterface

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

// addUnixFileInfo is a no-op on Windows (UID/GID don't exist)
func addUnixFileInfo(L *lua.LState, result *lua.LTable, info os.FileInfo) {
	// Windows doesn't have UID/GID, so we don't add these fields
}
