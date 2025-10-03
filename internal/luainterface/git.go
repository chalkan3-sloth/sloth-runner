package luainterface

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

const (
	luaGitRepoTypeName = "git_repo"
)

// GitRepo represents a git repository in Go.
type GitRepo struct {
	Path string
}

// --- Module Functions ---

// git.clone(url, path) -> repo
// Idempotent: If path exists and is a git repo, skip clone and return repo object
func gitClone(L *lua.LState) int {
	url := L.CheckString(1)
	path := L.CheckString(2)

	// Check if directory already exists
	if _, err := os.Stat(path); err == nil {
		// Directory exists, check if it's a git repository
		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			// It's already a git repository, return repo object without cloning
			repo := &GitRepo{Path: path}
			ud := L.NewUserData()
			ud.Value = repo
			L.SetMetatable(ud, L.GetTypeMetatable(luaGitRepoTypeName))
			L.Push(ud)
			return 1
		}
		// Directory exists but is not a git repo
		L.RaiseError("git clone failed: directory '%s' exists but is not a git repository", path)
		return 0
	}

	// Directory doesn't exist, proceed with clone
	cmd := ExecCommand("git", "clone", url, path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		L.RaiseError("git clone failed: %s", stderr.String())
	}

	repo := &GitRepo{Path: path}
	ud := L.NewUserData()
	ud.Value = repo
	L.SetMetatable(ud, L.GetTypeMetatable(luaGitRepoTypeName))
	L.Push(ud)
	return 1
}

// --- Helper ---

func checkGitRepo(L *lua.LState) *GitRepo {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*GitRepo); ok {
		return v
	}
	L.ArgError(1, "git repo expected")
	return nil
}

// --- Object Methods ---

// repo:checkout(branch, { create = false })
func repoCheckout(L *lua.LState) int {
	repo := checkGitRepo(L)
	branch := L.CheckString(2)
	opts := L.OptTable(3, L.NewTable())
	create := L.GetField(opts, "create").(lua.LBool)

	var cmd *exec.Cmd
	if create {
		cmd = ExecCommand("git", "checkout", "-b", branch)
	} else {
		cmd = ExecCommand("git", "checkout", branch)
	}
	cmd.Dir = repo.Path
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		L.RaiseError("git checkout to branch '%s' failed: %s", branch, stderr.String())
	}

	L.Push(L.Get(1)) // return self
	return 1
}

// repo:pull()
func repoPull(L *lua.LState) int {
	repo := checkGitRepo(L)

	cmd := ExecCommand("git", "pull")
	cmd.Dir = repo.Path
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		L.RaiseError("git pull failed: %s", stderr.String())
	}

	L.Push(L.Get(1)) // return self
	return 1
}

// repo:push()
func repoPush(L *lua.LState) int {
	repo := checkGitRepo(L)

	cmd := ExecCommand("git", "push")
	cmd.Dir = repo.Path
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		L.RaiseError("git push failed: %s", stderr.String())
	}

	L.Push(L.Get(1)) // return self
	return 1
}

// __index metamethod
func repoIndex(L *lua.LState) int {
	repo := checkGitRepo(L)
	key := L.CheckString(2)

	switch key {
	case "path":
		L.Push(lua.LString(repo.Path))
	default:
		// Fallback to methods in the metatable
		mt := L.GetTypeMetatable(luaGitRepoTypeName)
		L.Push(L.GetField(mt, key))
	}
	return 1
}

var gitRepoMethods = map[string]lua.LGFunction{
	"checkout": repoCheckout,
	"pull":     repoPull,
	"push":     repoPush,
}

// GitLoader loads the git module.
func GitLoader(L *lua.LState) int {
	// Register the repo type
	mt := L.NewTypeMetatable(luaGitRepoTypeName)
	L.SetField(mt, "__index", L.NewFunction(repoIndex))
	L.SetFuncs(mt, gitRepoMethods)

	// Register the module functions
	mod := L.NewTable()
	L.SetField(mod, "clone", L.NewFunction(gitClone))

	L.Push(mod)
	return 1
}

func OpenGit(L *lua.LState) {
	// Register as preload module (for backward compatibility with require("git"))
	L.PreloadModule("git", GitLoader)
	
	// Also register as global module (for use without require)
	// This makes git.clone() available directly
	mt := L.NewTypeMetatable(luaGitRepoTypeName)
	L.SetField(mt, "__index", L.NewFunction(repoIndex))
	L.SetFuncs(mt, gitRepoMethods)

	// Create git table with functions
	gitTable := L.NewTable()
	L.SetField(gitTable, "clone", L.NewFunction(gitClone))
	
	// Set git as global
	L.SetGlobal("git", gitTable)
}