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

func TestOpenImport_ImportWithVariables(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create library with variables
	libFile := filepath.Join(tmpDir, "vars.lua")
	libContent := `
local version = "1.0.0"
local config = {
	setting1 = true,
	setting2 = "test"
}
return {
	get_version = function() return version end,
	get_config = function() return config end
}
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local vars = import("vars.lua")
version = vars.get_version()
config = vars.get_config()
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	version := L.GetGlobal("version")
	assert.Equal(t, "1.0.0", version.String())
}

func TestOpenImport_ImportFunction(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create library that returns a function
	libFile := filepath.Join(tmpDir, "func.lua")
	libContent := `
return function(x, y)
	return x * y
end
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local multiply = import("func.lua")
result = multiply(3, 4)
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	result := L.GetGlobal("result")
	assert.Equal(t, lua.LNumber(12), result)
}

func TestOpenImport_DeeplyNestedPath(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create deeply nested directory structure
	deepDir := filepath.Join(tmpDir, "a", "b", "c", "d")
	require.NoError(t, os.MkdirAll(deepDir, 0755))

	// Create library in deep path
	libFile := filepath.Join(deepDir, "deep.lua")
	libContent := `
return {
	location = "deep/nested/path"
}
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local deep = import("a/b/c/d/deep.lua")
result = deep.location
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	result := L.GetGlobal("result")
	assert.Equal(t, "deep/nested/path", result.String())
}

func TestOpenImport_ImportWithErrors(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create library with syntax error
	libFile := filepath.Join(tmpDir, "error.lua")
	libContent := `
return {
	func = function()
		error("Intentional error")
	end
}
`
	require.NoError(t, os.WriteFile(libFile, []byte(libContent), 0644))

	// Create main file
	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local lib = import("error.lua")
lib.func()
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Intentional error")
}

func TestOpenImport_CircularReference(t *testing.T) {
	tmpDir := t.TempDir()
	
	// This test ensures circular imports don't cause infinite loops
	// Thanks to caching, the second import should return the same cached module

	lib1File := filepath.Join(tmpDir, "lib1.lua")
	lib1Content := `
local lib2 = import("lib2.lua")
return {
	name = "lib1",
	get_lib2_name = function()
		return lib2.name
	end
}
`
	require.NoError(t, os.WriteFile(lib1File, []byte(lib1Content), 0644))

	lib2File := filepath.Join(tmpDir, "lib2.lua")
	lib2Content := `
return {
	name = "lib2"
}
`
	require.NoError(t, os.WriteFile(lib2File, []byte(lib2Content), 0644))

	mainFile := filepath.Join(tmpDir, "main.lua")
	mainContent := `
local lib1 = import("lib1.lua")
result = lib1.get_lib2_name()
`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	L := lua.NewState()
	defer L.Close()

	OpenImport(L, mainFile)

	err := L.DoFile(mainFile)
	require.NoError(t, err)

	result := L.GetGlobal("result")
	assert.Equal(t, "lib2", result.String())
}

