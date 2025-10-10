package luainterface

import (
	"os"
	"path/filepath"
	"testing"

	luafs "github.com/chalkan3-sloth/sloth-runner/internal/luainterface/modules/fs"
	lua "github.com/yuin/gopher-lua"
)

func TestLuaFsRead(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	// Create temp file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := "test content"
	tmpFile.WriteString(content)
	tmpFile.Close()

	script := `
		local fs = require("fs")
		local content = fs.read("` + tmpFile.Name() + `")
		return content
	`

	err = L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	result := L.Get(-1).String()
	if result != content {
		t.Errorf("Expected '%s', got '%s'", content, result)
	}
}

func TestLuaFsWrite(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "test write content"

	script := `
		local fs = require("fs")
		fs.write("` + filePath + `", "` + content + `")
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	// Verify file was written
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(data) != content {
		t.Errorf("Expected '%s', got '%s'", content, string(data))
	}
}

func TestLuaFsAppend(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	
	// Write initial content
	os.WriteFile(filePath, []byte("initial\n"), 0644)

	script := `
		local fs = require("fs")
		fs.append("` + filePath + `", "appended")
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	// Verify content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := "initial\nappended"
	if string(data) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(data))
	}
}

func TestLuaFsExists_True(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	script := `
		local fs = require("fs")
		local exists = fs.exists("` + tmpFile.Name() + `")
		return exists
	`

	err = L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	result := L.Get(-1)
	if bool(result.(lua.LBool)) != true {
		t.Error("Expected file to exist")
	}
}

func TestLuaFsExists_False(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	script := `
		local fs = require("fs")
		local exists = fs.exists("/nonexistent/file.txt")
		return exists
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	result := L.Get(-1)
	if bool(result.(lua.LBool)) != false {
		t.Error("Expected file to not exist")
	}
}

func TestLuaFsMkdir(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newdir")

	script := `
		local fs = require("fs")
		fs.mkdir("` + newDir + `")
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("Expected directory to be created")
	}
}

func TestLuaFsRm(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	script := `
		local fs = require("fs")
		fs.rm("` + tmpFile.Name() + `")
	`

	err = L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(tmpFile.Name()); !os.IsNotExist(err) {
		t.Error("Expected file to be deleted")
	}
}

func TestLuaFsRmR(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("test"), 0644)

	script := `
		local fs = require("fs")
		fs.rmr("` + subDir + `")
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	// Verify directory was deleted
	if _, err := os.Stat(subDir); !os.IsNotExist(err) {
		t.Error("Expected directory to be deleted recursively")
	}
}

func TestLuaFsLs(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("test"), 0644)

	script := `
		local fs = require("fs")
		local files = fs.ls("` + tmpDir + `")
		local count = 0
		for _, f in ipairs(files) do
			count = count + 1
		end
		return count
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	count := int(L.Get(-1).(lua.LNumber))
	if count < 2 {
		t.Errorf("Expected at least 2 files, got %d", count)
	}
}

func TestLuaFsTmpName(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	script := `
		local fs = require("fs")
		local tmpname = fs.tmpname()
		return tmpname
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	result := L.Get(-1).String()
	if result == "" {
		t.Error("Expected non-empty temp file name")
	}
}

func TestLuaFsSize(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	luafs.Open(L)

	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := "test content with specific size"
	tmpFile.WriteString(content)
	tmpFile.Close()

	script := `
		local fs = require("fs")
		local size = fs.size("` + tmpFile.Name() + `")
		return size
	`

	err = L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	result := L.Get(-1)
	size := int(result.(lua.LNumber))
	expectedSize := len(content)

	if size != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, size)
	}
}

func TestOpenFs(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	luafs.Open(L)

	// Test that fs module can be accessed and used
	script := `
		local fs = require("fs")
		return type(fs), type(fs.read), type(fs.write)
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	fsType := L.Get(-3).String()
	readType := L.Get(-2).String()
	writeType := L.Get(-1).String()

	if fsType != "table" {
		t.Errorf("Expected fs to be 'table', got '%s'", fsType)
	}
	if readType != "function" {
		t.Errorf("Expected read to be 'function', got '%s'", readType)
	}
	if writeType != "function" {
		t.Errorf("Expected write to be 'function', got '%s'", writeType)
	}
}

func TestFsIntegration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	luafs.Open(L)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "integration_test.txt")

	// Test write, read, exists workflow
	script := `
		local fs = require("fs")
		fs.write("` + testFile + `", "test content")
		local exists = fs.exists("` + testFile + `")
		local content = fs.read("` + testFile + `")
		return exists, content
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}

	exists := L.Get(-2)
	content := L.Get(-1)

	if bool(exists.(lua.LBool)) != true {
		t.Error("Expected file to exist")
	}

	if content.String() != "test content" {
		t.Errorf("Expected 'test content', got '%s'", content.String())
	}
}
