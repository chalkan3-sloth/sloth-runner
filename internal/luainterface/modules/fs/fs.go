package fs

import (
	"io"
	"os"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

// Read reads a file and returns its content
func Read(L *lua.LState) int {
	path := L.CheckString(1)

	data, err := os.ReadFile(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(string(data)))
	return 1
}

// Write writes content to a file
func Write(L *lua.LState) int {
	path := L.CheckString(1)
	content := L.CheckString(2)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// Append appends content to a file
func Append(L *lua.LState) int {
	path := L.CheckString(1)
	content := L.CheckString(2)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// Exists checks if a file or directory exists
func Exists(L *lua.LState) int {
	path := L.CheckString(1)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		L.Push(lua.LBool(false))
	} else if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	} else {
		L.Push(lua.LBool(true))
	}
	return 1
}

// Mkdir creates a directory
func Mkdir(L *lua.LState) int {
	path := L.CheckString(1)

	if err := os.MkdirAll(path, 0755); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// Rm removes a file
func Rm(L *lua.LState) int {
	path := L.CheckString(1)

	if err := os.Remove(path); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// RmR removes a directory recursively
func RmR(L *lua.LState) int {
	path := L.CheckString(1)

	if err := os.RemoveAll(path); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// Ls lists directory contents
func Ls(L *lua.LState) int {
	path := L.CheckString(1)

	entries, err := os.ReadDir(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	tbl := L.NewTable()
	for _, entry := range entries {
		tbl.Append(lua.LString(entry.Name()))
	}

	L.Push(tbl)
	return 1
}

// TmpName returns a temporary file name
func TmpName(L *lua.LState) int {
	prefix := "sloth-"
	if L.GetTop() >= 1 {
		prefix = L.CheckString(1)
	}

	f, err := os.CreateTemp("", prefix+"*")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	name := f.Name()
	f.Close()
	os.Remove(name)

	L.Push(lua.LString(name))
	return 1
}

// Size returns the size of a file
func Size(L *lua.LState) int {
	path := L.CheckString(1)

	info, err := os.Stat(path)
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(info.Size()))
	return 1
}

// Copy copies a file
func Copy(L *lua.LState) int {
	src := L.CheckString(1)
	dst := L.CheckString(2)

	source, err := os.Open(src)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// Basename returns the base name of a path
func Basename(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Base(path)))
	return 1
}

// Dirname returns the directory name of a path
func Dirname(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Dir(path)))
	return 1
}

// Loader returns the fs module loader
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"read":     Read,
		"write":    Write,
		"append":   Append,
		"exists":   Exists,
		"mkdir":    Mkdir,
		"rm":       Rm,
		"rmr":      RmR,
		"ls":       Ls,
		"tmpname":  TmpName,
		"size":     Size,
		"copy":     Copy,
		"basename": Basename,
		"dirname":  Dirname,
	})
	L.Push(mod)
	return 1
}

// Open registers the fs module and loads it globally
func Open(L *lua.LState) {
	L.PreloadModule("fs", Loader)
	if err := L.DoString(`fs = require("fs")`); err != nil {
		panic(err)
	}
}
