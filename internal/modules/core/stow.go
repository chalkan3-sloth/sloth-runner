package core

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// RegisterStowModule registers the stow module in the Lua state
func RegisterStowModule(L *lua.LState) {
	// Create stow module table
	stowModule := L.NewTable()

	// Register functions
	L.SetField(stowModule, "link", L.NewFunction(stowLink))
	L.SetField(stowModule, "unlink", L.NewFunction(stowUnlink))
	L.SetField(stowModule, "restow", L.NewFunction(stowRestow))

	// Set as global
	L.SetGlobal("stow", stowModule)
}

// stowLink creates symlinks for a package
// Usage: local success, msg = stow.link({package = "...", source_dir = "...", target_dir = "..."})
func stowLink(L *lua.LState) int {
	// Get parameters table
	params := L.CheckTable(1)

	pkg := getStringField(L, params, "package", "")
	sourceDir := getStringField(L, params, "source_dir", "")
	targetDir := getStringField(L, params, "target_dir", "")
	verbose := getBoolField(L, params, "verbose", false)
	noFolding := getBoolField(L, params, "no_folding", false)

	if pkg == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("package is required"))
		return 2
	}

	if sourceDir == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("source_dir is required"))
		return 2
	}

	if targetDir == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("target_dir is required"))
		return 2
	}

	// Build stow command
	args := []string{"-d", sourceDir, "-t", targetDir}

	if verbose {
		args = append(args, "-v")
	}

	if noFolding {
		args = append(args, "--no-folding")
	}

	args = append(args, pkg)

	// Execute stow
	cmd := exec.Command("stow", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("stow link failed: %s - %s", err.Error(), strings.TrimSpace(string(output)))))
		return 2
	}

	msg := "stow link successful"
	if verbose && len(output) > 0 {
		msg = strings.TrimSpace(string(output))
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(msg))
	return 2
}

// stowUnlink removes symlinks for a package
// Usage: local success, msg = stow.unlink({package = "...", source_dir = "...", target_dir = "..."})
func stowUnlink(L *lua.LState) int {
	// Get parameters table
	params := L.CheckTable(1)

	pkg := getStringField(L, params, "package", "")
	sourceDir := getStringField(L, params, "source_dir", "")
	targetDir := getStringField(L, params, "target_dir", "")
	verbose := getBoolField(L, params, "verbose", false)

	if pkg == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("package is required"))
		return 2
	}

	if sourceDir == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("source_dir is required"))
		return 2
	}

	if targetDir == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("target_dir is required"))
		return 2
	}

	// Build stow command with -D (delete/unlink) flag
	args := []string{"-D", "-d", sourceDir, "-t", targetDir}

	if verbose {
		args = append(args, "-v")
	}

	args = append(args, pkg)

	// Execute stow
	cmd := exec.Command("stow", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("stow unlink failed: %s - %s", err.Error(), strings.TrimSpace(string(output)))))
		return 2
	}

	msg := "stow unlink successful"
	if verbose && len(output) > 0 {
		msg = strings.TrimSpace(string(output))
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(msg))
	return 2
}

// stowRestow removes and re-creates symlinks for a package
// Usage: local success, msg = stow.restow({package = "...", source_dir = "...", target_dir = "..."})
func stowRestow(L *lua.LState) int {
	// Get parameters table
	params := L.CheckTable(1)

	pkg := getStringField(L, params, "package", "")
	sourceDir := getStringField(L, params, "source_dir", "")
	targetDir := getStringField(L, params, "target_dir", "")
	verbose := getBoolField(L, params, "verbose", false)
	noFolding := getBoolField(L, params, "no_folding", false)

	if pkg == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("package is required"))
		return 2
	}

	if sourceDir == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("source_dir is required"))
		return 2
	}

	if targetDir == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("target_dir is required"))
		return 2
	}

	// Expand paths
	sourcePath, _ := filepath.Abs(sourceDir)
	targetPath, _ := filepath.Abs(targetDir)

	// Build stow command with -R (restow) flag
	args := []string{"-R", "-d", sourcePath, "-t", targetPath}

	if verbose {
		args = append(args, "-v")
	}

	if noFolding {
		args = append(args, "--no-folding")
	}

	args = append(args, pkg)

	// Execute stow
	cmd := exec.Command("stow", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("stow restow failed: %s - %s", err.Error(), strings.TrimSpace(string(output)))))
		return 2
	}

	msg := "stow restow successful"
	if verbose && len(output) > 0 {
		msg = strings.TrimSpace(string(output))
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(msg))
	return 2
}
