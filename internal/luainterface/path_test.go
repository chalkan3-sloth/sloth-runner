package luainterface

import (
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestPathModule(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{
			name: "path.join",
			script: `
				local p = path.join("usr", "local", "bin")
				assert(type(p) == "string", "should return a string")
				assert(string.find(p, "usr") ~= nil, "should contain usr")
				assert(string.find(p, "local") ~= nil, "should contain local")
				assert(string.find(p, "bin") ~= nil, "should contain bin")
			`,
		},
		{
			name: "path.dir",
			script: `
				local d = path.dir("/usr/local/bin/file.txt")
				assert(type(d) == "string", "should return a string")
			`,
		},
		{
			name: "path.base",
			script: `
				local b = path.base("/usr/local/bin/file.txt")
				assert(b == "file.txt", "should return the base name")
			`,
		},
		{
			name: "path.ext",
			script: `
				local e = path.ext("file.txt")
				assert(e == ".txt", "should return the extension")
			`,
		},
		{
			name: "path.clean",
			script: `
				local c = path.clean("/usr//local/../local/./bin")
				assert(type(c) == "string", "should return a string")
			`,
		},
		{
			name: "path.abs",
			script: `
				local a = path.abs(".")
				assert(type(a) == "string", "should return a string")
				assert(string.sub(a, 1, 1) == "/", "should be absolute path")
			`,
		},
		{
			name: "path.is_abs",
			script: `
				assert(path.is_abs("/usr/local") == true, "should be absolute")
				assert(path.is_abs("relative/path") == false, "should not be absolute")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			loadPathModule(L)

			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestPathExists(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	loadPathModule(L)

	L.SetGlobal("testFile", lua.LString(testFile))
	L.SetGlobal("tmpDir", lua.LString(tmpDir))

	script := `
		assert(path.exists(testFile) == true, "file should exist")
		assert(path.exists(tmpDir) == true, "directory should exist")
		assert(path.exists("/nonexistent/path/xyz") == false, "nonexistent path should return false")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
