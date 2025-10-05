package core

import (
	"fmt"
	"os/exec"

	lua "github.com/yuin/gopher-lua"
)

// execAsTaskUserGit executes a git command, optionally as a specific user using sudo
func execAsTaskUserGit(L *lua.LState, args []string) *exec.Cmd {
	taskUser := L.GetGlobal("__TASK_USER__")
	if taskUser.Type() == lua.LTString && taskUser.String() != "" && taskUser.String() != "root" {
		// Run as specific user using sudo
		allArgs := append([]string{"-u", taskUser.String(), "git"}, args...)
		return exec.Command("sudo", allArgs...)
	}
	// Run as current user (root)
	return exec.Command("git", args...)
}

// RegisterGitModule registers the git module in the Lua state
func RegisterGitModule(L *lua.LState) {
	// Create git module table
	gitModule := L.NewTable()

	// Register functions
	L.SetField(gitModule, "clone", L.NewFunction(gitClone))
	L.SetField(gitModule, "pull", L.NewFunction(gitPull))
	L.SetField(gitModule, "status", L.NewFunction(gitStatus))
	L.SetField(gitModule, "checkout", L.NewFunction(gitCheckout))
	L.SetField(gitModule, "commit", L.NewFunction(gitCommit))
	L.SetField(gitModule, "push", L.NewFunction(gitPush))
	L.SetField(gitModule, "is_repo", L.NewFunction(gitIsRepo))
	L.SetField(gitModule, "ensure_clean", L.NewFunction(gitEnsureClean))

	// Set as global
	L.SetGlobal("git", gitModule)
}

// gitClone clones a git repository with idempotency
// Usage: local repo, err = git.clone({url = "...", local_path = "...", clean = false})
func gitClone(L *lua.LState) int {
	// Get parameters table
	params := L.CheckTable(1)

	url := getStringField(L, params, "url", "")
	localPath := getStringField(L, params, "local_path", "")
	branch := getStringField(L, params, "branch", "")
	depth := getIntField(L, params, "depth", 0)
	clean := getBoolField(L, params, "clean", false)

	if url == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("url is required"))
		return 2
	}

	if localPath == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("local_path is required"))
		return 2
	}

	// Check if directory already exists and is a git repo
	checkCmd := exec.Command("git", "-C", localPath, "rev-parse", "--git-dir")
	if err := checkCmd.Run(); err == nil {
		// Directory exists and is a git repo
		if !clean {
			// Return existing repo
			repoTable := L.NewTable()
			L.SetField(repoTable, "path", lua.LString(localPath))
			L.SetField(repoTable, "url", lua.LString(url))
			L.SetField(repoTable, "exists", lua.LBool(true))

			L.Push(repoTable)
			L.Push(lua.LNil)
			return 2
		}
		// Clean requested, remove directory
		exec.Command("rm", "-rf", localPath).Run()
	}

	// Build git clone command
	args := []string{"clone"}

	if branch != "" {
		args = append(args, "-b", branch)
	}

	if depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", depth))
	}

	args = append(args, url, localPath)

	// Execute git clone (as task user if specified)
	cmd := execAsTaskUserGit(L, args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("git clone failed: %s - %s", err.Error(), string(output))))
		return 2
	}

	// Return repo object (table with path)
	repoTable := L.NewTable()
	L.SetField(repoTable, "path", lua.LString(localPath))
	L.SetField(repoTable, "url", lua.LString(url))
	L.SetField(repoTable, "exists", lua.LBool(false))

	L.Push(repoTable)
	L.Push(lua.LNil)
	return 2
}

// gitPull pulls changes from remote repository
// Usage: local success, err = git.pull({path = "..."})
func gitPull(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	if path == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("path is required"))
		return 2
	}

	// Execute git pull
	cmd := exec.Command("git", "-C", path, "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("git pull failed: %s - %s", err.Error(), string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("pull successful"))
	return 2
}

// gitStatus gets the status of the repository
// Usage: local status, err = git.status({path = "..."})
func gitStatus(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	if path == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("path is required"))
		return 2
	}

	// Execute git status
	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("git status failed: %s", err.Error())))
		return 2
	}

	L.Push(lua.LString(string(output)))
	L.Push(lua.LNil)
	return 2
}

// gitCheckout checks out a branch or commit
// Usage: local success, err = git.checkout({path = "...", branch = "..."})
func gitCheckout(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	branch := getStringField(L, params, "branch", "")

	if path == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("path is required"))
		return 2
	}

	if branch == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("branch is required"))
		return 2
	}

	// Execute git checkout
	cmd := exec.Command("git", "-C", path, "checkout", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("git checkout failed: %s - %s", err.Error(), string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("checkout successful"))
	return 2
}

// gitCommit creates a commit
// Usage: local success, err = git.commit({path = "...", message = "..."})
func gitCommit(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	message := getStringField(L, params, "message", "")
	addAll := getBoolField(L, params, "add_all", false)

	if path == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("path is required"))
		return 2
	}

	if message == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("message is required"))
		return 2
	}

	// Add all files if requested
	if addAll {
		cmd := exec.Command("git", "-C", path, "add", ".")
		if output, err := cmd.CombinedOutput(); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("git add failed: %s - %s", err.Error(), string(output))))
			return 2
		}
	}

	// Execute git commit
	cmd := exec.Command("git", "-C", path, "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("git commit failed: %s - %s", err.Error(), string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("commit successful"))
	return 2
}

// gitPush pushes changes to remote repository
// Usage: local success, err = git.push({path = "..."})
func gitPush(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	remote := getStringField(L, params, "remote", "origin")
	branch := getStringField(L, params, "branch", "")

	if path == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("path is required"))
		return 2
	}

	// Build git push command
	args := []string{"-C", path, "push", remote}
	if branch != "" {
		args = append(args, branch)
	}

	// Execute git push
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("git push failed: %s - %s", err.Error(), string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("push successful"))
	return 2
}

// Helper function to get int field from table
func getIntField(L *lua.LState, tbl *lua.LTable, key string, defaultValue int) int {
	lv := tbl.RawGetString(key)
	if lv == lua.LNil {
		return defaultValue
	}
	if num, ok := lv.(lua.LNumber); ok {
		return int(num)
	}
	return defaultValue
}

// Helper function to get bool field from table
func getBoolField(L *lua.LState, tbl *lua.LTable, key string, defaultValue bool) bool {
	lv := tbl.RawGetString(key)
	if lv == lua.LNil {
		return defaultValue
	}
	if b, ok := lv.(lua.LBool); ok {
		return bool(b)
	}
	return defaultValue
}

// gitIsRepo checks if a directory is a git repository
// Usage: local is_repo = git.is_repo({path = "..."})
func gitIsRepo(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	if path == "" {
		L.Push(lua.LBool(false))
		return 1
	}

	// Check if directory is a git repo
	cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
	err := cmd.Run()

	L.Push(lua.LBool(err == nil))
	return 1
}

// gitEnsureClean ensures a directory is clean (removes it if exists)
// Usage: local success, err = git.ensure_clean({path = "..."})
func gitEnsureClean(L *lua.LState) int {
	params := L.CheckTable(1)

	path := getStringField(L, params, "path", "")
	if path == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("path is required"))
		return 2
	}

	// Remove directory if it exists
	cmd := exec.Command("rm", "-rf", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to clean directory: %s - %s", err.Error(), string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LNil)
	return 2
}
