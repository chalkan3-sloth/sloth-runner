package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// execAsTaskUser executes a command, optionally as a specific user using sudo
func execAsTaskUser(L *lua.LState, command string, args []string) *exec.Cmd {
	taskUser := L.GetGlobal("__TASK_USER__")
	if taskUser.Type() == lua.LTString && taskUser.String() != "" && taskUser.String() != "root" {
		// Run as specific user using sudo
		allArgs := append([]string{"-u", taskUser.String(), command}, args...)
		return exec.Command("sudo", allArgs...)
	}
	// Run as current user (root)
	return exec.Command(command, args...)
}

// RegisterStowModule registers the stow module in the Lua state
func RegisterStowModule(L *lua.LState) {
	// Create stow module table
	stowModule := L.NewTable()

	// Register functions
	L.SetField(stowModule, "link", L.NewFunction(stowLink))
	L.SetField(stowModule, "unlink", L.NewFunction(stowUnlink))
	L.SetField(stowModule, "restow", L.NewFunction(stowRestow))
	L.SetField(stowModule, "ensure_target", L.NewFunction(stowEnsureTarget))

	// Set as global
	L.SetGlobal("stow", stowModule)
}

// stowLink creates symlinks for a package
// Usage: local success, msg = stow.link({package = "...", source_dir = "...", target_dir = "...", create_target = true})
func stowLink(L *lua.LState) int {
	// Get parameters table
	params := L.CheckTable(1)

	pkg := getStringField(L, params, "package", "")
	sourceDir := getStringField(L, params, "source_dir", "")
	targetDir := getStringField(L, params, "target_dir", "")
	verbose := getBoolField(L, params, "verbose", false)
	noFolding := getBoolField(L, params, "no_folding", false)
	createTarget := getBoolField(L, params, "create_target", true) // Default true

	if pkg == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("package name is required"))
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

	// Create target directory if requested and it doesn't exist
	if createTarget {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			// Get task user to determine ownership
			taskUser := L.GetGlobal("__TASK_USER__")

			// Create directory with proper permissions
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				L.Push(lua.LBool(false))
				L.Push(lua.LString(fmt.Sprintf("failed to create target directory: %v", err)))
				return 2
			}

			// If task user is set and not root, change ownership
			if taskUser.Type() == lua.LTString && taskUser.String() != "" && taskUser.String() != "root" {
				chownCmd := exec.Command("chown", "-R", taskUser.String()+":"+taskUser.String(), targetDir)
				if err := chownCmd.Run(); err != nil {
					L.Push(lua.LBool(false))
					L.Push(lua.LString(fmt.Sprintf("failed to set ownership of target directory: %v", err)))
					return 2
				}
			}
		}
	}

	// IDEMPOTENCY CHECK: Verify if package is already stowed
	// Check if symlinks already exist by running stow with --no-folding in simulation mode
	checkArgs := []string{"-d", sourceDir, "-t", targetDir, "-n", "-v"}
	if noFolding {
		checkArgs = append(checkArgs, "--no-folding")
	}
	checkArgs = append(checkArgs, pkg)

	checkCmd := execAsTaskUser(L, "stow", checkArgs)
	checkOutput, checkErr := checkCmd.CombinedOutput()

	// If simulation shows no changes needed (empty output), it's already stowed
	// Note: "LINK:" in output means stow WILL create links, so NOT already stowed
	outputStr := strings.TrimSpace(string(checkOutput))
	if checkErr == nil && outputStr == "" {
		// No output means no changes needed - already stowed
		L.Push(lua.LBool(true))
		L.Push(lua.LString("package already stowed"))
		return 2
	}

	// Build stow command for actual execution
	args := []string{"-d", sourceDir, "-t", targetDir}

	if verbose {
		args = append(args, "-v")
	}

	if noFolding {
		args = append(args, "--no-folding")
	}

	args = append(args, pkg)

	// Execute stow (as task user if specified)
	cmd := execAsTaskUser(L, "stow", args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if error is about existing links (idempotent case)
		errStr := strings.ToLower(string(output))
		if strings.Contains(errStr, "already") || strings.Contains(errStr, "existing target") {
			L.Push(lua.LBool(true))
			L.Push(lua.LString("package already stowed"))
			return 2
		}
		
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

	// Execute stow (as task user if specified)
	cmd := execAsTaskUser(L, "stow", args)
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

	// Execute stow (as task user if specified)
	cmd := execAsTaskUser(L, "stow", args)
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

// stowEnsureTarget ensures a target directory exists with proper ownership
// Usage: local success, msg = stow.ensure_target({path = "/home/user/.zsh", owner = "user"})
func stowEnsureTarget(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	owner := getStringField(L, params, "owner", "")
	mode := getStringField(L, params, "mode", "0755")

	if path == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("path is required"))
		return 2
	}

	// If owner is not specified, try to get from __TASK_USER__
	if owner == "" {
		taskUser := L.GetGlobal("__TASK_USER__")
		if taskUser.Type() == lua.LTString && taskUser.String() != "" {
			owner = taskUser.String()
		}
	}

	// Check if directory already exists
	if info, err := os.Stat(path); err == nil {
		// Directory exists
		if !info.IsDir() {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("%s exists but is not a directory", path)))
			return 2
		}

		// Directory already exists
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("directory %s already exists", path)))
		return 2
	}

	// Parse mode (simplified - assumes octal)
	var fileMode os.FileMode = 0755
	if mode != "" {
		fmt.Sscanf(mode, "%o", &fileMode)
	}

	// Create directory
	if err := os.MkdirAll(path, fileMode); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
		return 2
	}

	// Set ownership if specified
	if owner != "" && owner != "root" {
		chownCmd := exec.Command("chown", "-R", owner+":"+owner, path)
		if err := chownCmd.Run(); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to set ownership: %v", err)))
			return 2
		}
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("directory %s created successfully", path)))
	return 2
}
