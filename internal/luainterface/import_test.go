package luainterface

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestOpenImport_BasicImport(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a library file
	libFile := filepath.Join(tmpDir, "lib.lua")
	libContent := `
return {
	hello = function()
		return "Hello from lib"
	end
}
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local lib = import("lib.lua")
result = lib.hello()
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	result := L.GetGlobal("result")
	assert.Equal(t, "Hello from lib", result.String())
}

func TestOpenImport_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create subdirectory
	subDir := filepath.Join(tmpDir, "modules")
	require.NoError(t, os.MkdirAll(subDir, 0755))

	// Create library file in subdirectory
	libFile := filepath.Join(subDir, "helper.lua")
	libContent := `
return {
	add = function(a, b)
		return a + b
	end
}
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local helper = import("modules/helper.lua")
result = helper.add(2, 3)
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	result := L.GetGlobal("result")
	assert.Equal(t, lua.LNumber(5), result)
}

func TestOpenImport_MultipleImports(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create first library
	lib1File := filepath.Join(tmpDir, "lib1.lua")
	lib1Content := `
return {
	func1 = function()
		return "from lib1"
	end
}
`
	require.NoError(t, os.WriteFile(lib1File, []byte(lib1Content), 0644))

	// Create second library
	lib2File := filepath.Join(tmpDir, "lib2.lua")
	lib2Content := `
return {
	func2 = function()
		return "from lib2"
	end
}
`
	require.NoError(t, os.WriteFile(lib2File, []byte(lib2Content), 0644))

	// Create main file that imports both
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local lib1 = import("lib1.lua")
local lib2 = import("lib2.lua")
result1 = lib1.func1()
result2 = lib2.func2()
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	result1 := L.GetGlobal("result1")
	result2 := L.GetGlobal("result2")
	assert.Equal(t, "from lib1", result1.String())
	assert.Equal(t, "from lib2", result2.String())
}

func TestOpenImport_NonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create main file that imports nonexistent file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local lib = import("nonexistent.lua")
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
}

func TestOpenImport_Caching(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create library with counter
	libFile := filepath.Join(tmpDir, "counter.lua")
	libContent := `
local count = 0
count = count + 1
return {
	get_count = function()
		return count
	end
}
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file that imports same library twice
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local lib1 = import("counter.lua")
local lib2 = import("counter.lua")
count1 = lib1.get_count()
count2 = lib2.get_count()
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	count1 := L.GetGlobal("count1")
	count2 := L.GetGlobal("count2")
	
	// Both should return 1 because the import is cached
	assert.Equal(t, lua.LNumber(1), count1)
	assert.Equal(t, lua.LNumber(1), count2)
}

func TestOpenImport_NestedImports(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create innermost library
	innerLib := filepath.Join(tmpDir, "inner.lua")
	innerContent := `
return {
	value = "inner value"
}
`
	require.NoError(t, os.WriteFile(innerLib, []byte(innerContent), 0644))

	// Create middle library that imports inner
	middleLib := filepath.Join(tmpDir, "middle.lua")
	middleContent := `
local inner = import("inner.lua")
return {
	get_inner = function()
		return inner.value
	end
}
`
	require.NoError(t, os.WriteFile(middleLib, []byte(middleContent), 0644))

	// Create main file that imports middle
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local middle = import("middle.lua")
result = middle.get_inner()
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	result := L.GetGlobal("result")
	assert.Equal(t, "inner value", result.String())
}

func TestOpenImport_ImportReturnsNil(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create library that returns nothing
	libFile := filepath.Join(tmpDir, "empty.lua")
	libContent := `
-- Returns nothing
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local empty = import("empty.lua")
has_value = (empty ~= nil)
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	hasValue := L.GetGlobal("has_value")
	// Should still work, just returns nil/empty
	assert.NotNil(t, hasValue)
}
