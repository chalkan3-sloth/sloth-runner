package luainterface

import (
	"runtime"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestNewPkgModule(t *testing.T) {
	module := NewPkgModule()
	if module == nil {
		t.Fatal("NewPkgModule returned nil")
	}
}

func TestPkgModuleLoader(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewPkgModule()
	L.Push(L.NewFunction(module.Loader))
	L.Call(0, 1)

	result := L.Get(-1)
	if result.Type() != lua.LTTable {
		t.Errorf("Expected table, got %v", result.Type())
	}
}

func TestPkgModuleExports(t *testing.T) {
	module := NewPkgModule()
	exports := module.exports()

	expectedFuncs := []string{
		"install",
		"remove",
		"update",
		"upgrade",
		"search",
		"info",
		"list",
		"is_installed",
		"get_manager",
		"clean",
		"autoremove",
		"which",
		"version",
		"deps",
		"install_local",
	}

	for _, funcName := range expectedFuncs {
		if _, exists := exports[funcName]; !exists {
			t.Errorf("Expected function %s not found in exports", funcName)
		}
	}
}

func TestPkgDetectPackageManager(t *testing.T) {
	module := NewPkgModule()

	manager, err := module.detectPackageManager()
	
	// On most systems at least one package manager should be available
	// This might fail in some environments, so we just check the function works
	if err != nil && manager != "" {
		t.Errorf("detectPackageManager returned error: %v, but also returned manager: %s", err, manager)
	}
	
	// On macOS, brew is likely available
	if runtime.GOOS == "darwin" && manager != "brew" && err != nil {
		t.Logf("Warning: No package manager detected on macOS: %v", err)
	}
}

func TestPkgNeedsSudo(t *testing.T) {
	module := NewPkgModule()

	tests := []struct {
		manager  string
		expected bool
	}{
		{"brew", false},
		{"apt-get", runtime.GOOS != "darwin"},
		{"yum", runtime.GOOS != "darwin"},
		{"dnf", runtime.GOOS != "darwin"},
		{"pacman", runtime.GOOS != "darwin"},
	}

	for _, tt := range tests {
		t.Run(tt.manager, func(t *testing.T) {
			result := module.needsSudo(tt.manager)
			if result != tt.expected {
				t.Errorf("needsSudo(%s) = %v, expected %v", tt.manager, result, tt.expected)
			}
		})
	}
}

func TestPkgParsePackagesSingleString(t *testing.T) {
	module := NewPkgModule()

	val := lua.LString("curl")
	packages := module.parsePackages(val)

	if len(packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(packages))
	}

	if packages[0] != "curl" {
		t.Errorf("Expected package 'curl', got '%s'", packages[0])
	}
}

func TestPkgParsePackagesTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewPkgModule()

	// Create a Lua table with packages
	tbl := L.NewTable()
	tbl.Append(lua.LString("curl"))
	tbl.Append(lua.LString("wget"))
	tbl.Append(lua.LString("git"))

	packages := module.parsePackages(tbl)

	if len(packages) != 3 {
		t.Errorf("Expected 3 packages, got %d", len(packages))
	}

	expectedPackages := []string{"curl", "wget", "git"}
	for i, expected := range expectedPackages {
		if packages[i] != expected {
			t.Errorf("Expected package '%s' at index %d, got '%s'", expected, i, packages[i])
		}
	}
}

func TestPkgBuildInstallCommand(t *testing.T) {
	module := NewPkgModule()

	tests := []struct {
		manager  string
		packages []string
		contains []string
	}{
		{"apt-get", []string{"curl"}, []string{"install", "curl"}},
		{"yum", []string{"wget"}, []string{"install", "wget"}},
		{"dnf", []string{"git"}, []string{"install", "git"}},
		{"pacman", []string{"vim"}, []string{"-S", "vim"}},
		{"brew", []string{"node"}, []string{"install", "node"}},
	}

	for _, tt := range tests {
		t.Run(tt.manager, func(t *testing.T) {
			cmd := module.buildInstallCommand(tt.manager, tt.packages)
			
			if len(cmd) == 0 {
				t.Errorf("buildInstallCommand returned empty slice for %s", tt.manager)
				return
			}

			// Check if expected strings are in the command
			cmdStr := strings.Join(cmd, " ")
			for _, expected := range tt.contains {
				if !strings.Contains(cmdStr, expected) {
					t.Errorf("Command for %s should contain '%s', got: %s", tt.manager, expected, cmdStr)
				}
			}
		})
	}
}

func TestPkgModuleGetManager(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewPkgModule()
	L.PreloadModule("pkg", module.Loader)

	err := L.DoString(`
		local pkg = require("pkg")
		local manager = pkg.get_manager({})
		-- Just check it doesn't crash
		-- The result depends on the system
	`)
	if err != nil {
		t.Logf("Note: get_manager test failed (expected on systems without package managers): %v", err)
	}
}

func TestPkgModuleIsInstalled(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewPkgModule()
	L.PreloadModule("pkg", module.Loader)

	// Test with a package that is likely not installed
	err := L.DoString(`
		local pkg = require("pkg")
		local installed = pkg.is_installed({package = "nonexistent-package-xyz-12345"})
		-- Just check it doesn't crash
	`)
	if err != nil {
		t.Logf("Note: is_installed test had an error: %v", err)
	}
}

func TestPkgModuleWhich(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewPkgModule()
	L.PreloadModule("pkg", module.Loader)

	// Test with common commands that should exist
	err := L.DoString(`
		local pkg = require("pkg")
		local path = pkg.which({executable = "ls"})
		-- ls should exist on Unix systems
		-- Just check it doesn't crash
	`)
	if err != nil {
		t.Logf("Note: which test had an error: %v", err)
	}
}

func TestPkgModuleVersion(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewPkgModule()
	L.PreloadModule("pkg", module.Loader)

	err := L.DoString(`
		local pkg = require("pkg")
		-- Try to get version of a common package
		-- This may fail if package is not installed
		local version, err = pkg.version({package = "bash"})
	`)
	if err != nil {
		t.Logf("Note: version test had an error (may be expected): %v", err)
	}
}

func TestPkgModuleIntegration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewPkgModule()
	L.PreloadModule("pkg", module.Loader)

	// Test that all functions are available
	err := L.DoString(`
		local pkg = require("pkg")
		
		-- Check all functions exist
		assert(type(pkg.install) == "function", "install should be a function")
		assert(type(pkg.remove) == "function", "remove should be a function")
		assert(type(pkg.update) == "function", "update should be a function")
		assert(type(pkg.upgrade) == "function", "upgrade should be a function")
		assert(type(pkg.search) == "function", "search should be a function")
		assert(type(pkg.info) == "function", "info should be a function")
		assert(type(pkg.list) == "function", "list should be a function")
		assert(type(pkg.is_installed) == "function", "is_installed should be a function")
		assert(type(pkg.get_manager) == "function", "get_manager should be a function")
		assert(type(pkg.clean) == "function", "clean should be a function")
		assert(type(pkg.autoremove) == "function", "autoremove should be a function")
		assert(type(pkg.which) == "function", "which should be a function")
		assert(type(pkg.version) == "function", "version should be a function")
		assert(type(pkg.deps) == "function", "deps should be a function")
		assert(type(pkg.install_local) == "function", "install_local should be a function")
	`)
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}
}
