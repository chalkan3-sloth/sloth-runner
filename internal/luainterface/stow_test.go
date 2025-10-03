package luainterface

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	lua "github.com/yuin/gopher-lua"
)

func TestStowModule(t *testing.T) {
	// Create temporary directories for testing
	tmpDir := t.TempDir()
	stowDir := filepath.Join(tmpDir, "dotfiles")
	targetDir := filepath.Join(tmpDir, "target")
	
	os.MkdirAll(stowDir, 0755)
	os.MkdirAll(targetDir, 0755)

	// Create a test package
	pkgDir := filepath.Join(stowDir, "testpkg")
	os.MkdirAll(pkgDir, 0755)
	
	testFile := filepath.Join(pkgDir, "testfile.txt")
	os.WriteFile(testFile, []byte("test content"), 0644)

	L := lua.NewState()
	defer L.Close()

	sm, err := state.NewStateManager(filepath.Join(tmpDir, "state.db"))
	if err != nil {
		t.Fatalf("Failed to create state manager: %v", err)
	}
	defer sm.Close()

	stow := NewStowModule(sm)
	L.PreloadModule("stow", stow.Loader)

	t.Run("ParseConfig", func(t *testing.T) {
		script := `
		local stow = require("stow")
		
		-- Test with package only
		result, err = stow.stow({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		-- Should return something (result or error)
		return (result ~= nil) or (err ~= nil)
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected stow to return a result or error")
		}
	})

	t.Run("StowPackage", func(t *testing.T) {
		// Check if stow command is available
		if _, err := os.Stat("/usr/bin/stow"); os.IsNotExist(err) {
			if _, err := os.Stat("/usr/local/bin/stow"); os.IsNotExist(err) {
				t.Skip("stow command not found, skipping integration test")
			}
		}

		script := `
		local stow = require("stow")
		
		result, err = stow.stow({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err then
			return false, err
		end
		
		return result.status == "stowed" or result.status == "already_stowed"
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			errMsg := ""
			if L.GetTop() >= 2 {
				errMsg = L.Get(-2).String()
			}
			t.Errorf("Expected stow to succeed, got error: %s", errMsg)
		}
	})

	t.Run("IsStowed", func(t *testing.T) {
		script := `
		local stow = require("stow")
		
		is_stowed, err = stow.is_stowed({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		return is_stowed
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		// This may be false if stow command is not available
		// So we just check that it returns a boolean
		result := L.Get(-1)
		if result.Type() != lua.LTBool {
			t.Error("Expected is_stowed to return a boolean")
		}
	})

	t.Run("GetLinks", func(t *testing.T) {
		script := `
		local stow = require("stow")
		
		links, err = stow.get_links({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err then
			return false
		end
		
		return type(links) == "table"
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected get_links to return a table")
		}
	})

	t.Run("ListPackages", func(t *testing.T) {
		// Create another package
		os.MkdirAll(filepath.Join(stowDir, "anotherpkg"), 0755)

		script := `
		local stow = require("stow")
		
		packages, err = stow.list_packages({
			dir = "` + stowDir + `"
		})
		
		if err then
			return false
		end
		
		local count = 0
		for _ in pairs(packages) do
			count = count + 1
		end
		
		return count >= 2
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected list_packages to return at least 2 packages")
		}
	})

	t.Run("Verify", func(t *testing.T) {
		script := `
		local stow = require("stow")
		
		result, err = stow.verify({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err then
			return false
		end
		
		return result.package == "testpkg" and result.total_files ~= nil
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected verify to return package info")
		}
	})

	t.Run("CheckSimulate", func(t *testing.T) {
		if _, err := os.Stat("/usr/bin/stow"); os.IsNotExist(err) {
			if _, err := os.Stat("/usr/local/bin/stow"); os.IsNotExist(err) {
				t.Skip("stow command not found, skipping integration test")
			}
		}

		script := `
		local stow = require("stow")
		
		result, err = stow.check({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err then
			return false
		end
		
		return result.package == "testpkg"
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected check to return result")
		}
	})

	t.Run("Idempotency", func(t *testing.T) {
		if _, err := os.Stat("/usr/bin/stow"); os.IsNotExist(err) {
			if _, err := os.Stat("/usr/local/bin/stow"); os.IsNotExist(err) {
				t.Skip("stow command not found, skipping integration test")
			}
		}

		script := `
		local stow = require("stow")
		
		-- First stow
		result1, err1 = stow.stow({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err1 then
			return false, "first stow failed: " .. err1
		end
		
		-- Second stow (should be idempotent)
		result2, err2 = stow.stow({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err2 then
			return false, "second stow failed: " .. err2
		end
		
		-- Second call should report no changes
		return result2.status == "already_stowed" and not result2.changed
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if L.GetTop() >= 2 && L.Get(-2).Type() == lua.LTString {
			t.Errorf("Idempotency test failed: %s", L.Get(-2).String())
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected second stow to be idempotent (no changes)")
		}
	})

	t.Run("Unstow", func(t *testing.T) {
		if _, err := os.Stat("/usr/bin/stow"); os.IsNotExist(err) {
			if _, err := os.Stat("/usr/local/bin/stow"); os.IsNotExist(err) {
				t.Skip("stow command not found, skipping integration test")
			}
		}

		script := `
		local stow = require("stow")
		
		-- Ensure it's stowed first
		stow.stow({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		-- Now unstow
		result, err = stow.unstow({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err then
			return false
		end
		
		return result.status == "unstowed" or result.status == "not_stowed"
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected unstow to succeed")
		}
	})

	t.Run("ErrorHandling_MissingPackage", func(t *testing.T) {
		script := `
		local stow = require("stow")
		
		result, err = stow.stow({})
		
		return err ~= nil and string.match(err, "required")
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected error when package is missing")
		}
	})

	t.Run("ErrorHandling_InvalidPackage", func(t *testing.T) {
		script := `
		local stow = require("stow")
		
		result, err = stow.is_stowed({
			package = "nonexistent",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		-- Should return false, not error
		return not result and err == nil
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected is_stowed to return false for nonexistent package")
		}
	})
}

func TestStowModuleStateManagement(t *testing.T) {
	tmpDir := t.TempDir()
	stowDir := filepath.Join(tmpDir, "dotfiles")
	targetDir := filepath.Join(tmpDir, "target")
	
	os.MkdirAll(stowDir, 0755)
	os.MkdirAll(targetDir, 0755)

	pkgDir := filepath.Join(stowDir, "testpkg")
	os.MkdirAll(pkgDir, 0755)
	os.WriteFile(filepath.Join(pkgDir, "test.txt"), []byte("test"), 0644)

	L := lua.NewState()
	defer L.Close()

	sm, err := state.NewStateManager(filepath.Join(tmpDir, "state.db"))
	if err != nil {
		t.Fatalf("Failed to create state manager: %v", err)
	}
	defer sm.Close()

	stow := NewStowModule(sm)
	L.PreloadModule("stow", stow.Loader)

	t.Run("StateRecording", func(t *testing.T) {
		if _, err := os.Stat("/usr/bin/stow"); os.IsNotExist(err) {
			if _, err := os.Stat("/usr/local/bin/stow"); os.IsNotExist(err) {
				t.Skip("stow command not found, skipping integration test")
			}
		}

		script := `
		local stow = require("stow")
		
		result, err = stow.stow({
			package = "testpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		return err == nil
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		// Check if state was recorded
		resourceID := "stow:testpkg:" + targetDir
		stateData, err := sm.Get(resourceID)
		if err != nil && lua.LVAsBool(L.Get(-1)) {
			// Only fail if stow actually succeeded
			t.Errorf("Expected state to be recorded for stowed package: %v", err)
		}

		if stateData != nil {
			t.Logf("State recorded: %v", stateData)
		}
	})
}

func TestStowModuleAdvanced(t *testing.T) {
	tmpDir := t.TempDir()
	stowDir := filepath.Join(tmpDir, "dotfiles")
	targetDir := filepath.Join(tmpDir, "target")
	
	os.MkdirAll(stowDir, 0755)
	os.MkdirAll(targetDir, 0755)

	// Create package with nested structure
	pkgDir := filepath.Join(stowDir, "nestedpkg")
	os.MkdirAll(filepath.Join(pkgDir, ".config", "app"), 0755)
	os.WriteFile(filepath.Join(pkgDir, ".config", "app", "config.txt"), []byte("config"), 0644)

	L := lua.NewState()
	defer L.Close()

	sm, err := state.NewStateManager(filepath.Join(tmpDir, "state.db"))
	if err != nil {
		t.Fatalf("Failed to create state manager: %v", err)
	}
	defer sm.Close()

	stow := NewStowModule(sm)
	L.PreloadModule("stow", stow.Loader)

	t.Run("NestedPackage", func(t *testing.T) {
		if _, err := os.Stat("/usr/bin/stow"); os.IsNotExist(err) {
			if _, err := os.Stat("/usr/local/bin/stow"); os.IsNotExist(err) {
				t.Skip("stow command not found, skipping integration test")
			}
		}

		script := `
		local stow = require("stow")
		
		result, err = stow.stow({
			package = "nestedpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		if err then
			return false
		end
		
		-- Verify
		verify_result, verify_err = stow.verify({
			package = "nestedpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `"
		})
		
		return verify_result.total_files > 0
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected nested package to be processed correctly")
		}
	})

	t.Run("WithOptions", func(t *testing.T) {
		if _, err := os.Stat("/usr/bin/stow"); os.IsNotExist(err) {
			if _, err := os.Stat("/usr/local/bin/stow"); os.IsNotExist(err) {
				t.Skip("stow command not found, skipping integration test")
			}
		}

		script := `
		local stow = require("stow")
		
		result, err = stow.stow({
			package = "nestedpkg",
			dir = "` + stowDir + `",
			target = "` + targetDir + `",
			verbose = true,
			no_folding = true,
			ignore = {"*.bak", "*.swp"}
		})
		
		return err == nil or string.match(err, "already")
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if !lua.LVAsBool(L.Get(-1)) {
			t.Error("Expected stow with options to succeed")
		}
	})
}
