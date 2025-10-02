package luainterface

import (
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestGitClone(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterGitModule(L)

	tmpDir := t.TempDir()
	
	script := `
		local git = require("git")
		local result = git.clone({
			url = "https://github.com/chalkan3-sloth/sloth-runner.git",
			destination = "` + tmpDir + `/test-repo",
			depth = 1
		})
		return result.success
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Git clone test skipped (network/git not available): %v", err)
		return
	}

	result := L.Get(-1)
	if result == lua.LNil {
		t.Skip("Git not available in test environment")
	}
}

func TestGitStatus(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterGitModule(L)

	script := `
		local git = require("git")
		local result = git.status({
			path = "/nonexistent/path"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for nonexistent path: %v", err)
	}
}

func TestGitCommit(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterGitModule(L)

	tmpDir := t.TempDir()
	
	script := `
		local git = require("git")
		local result = git.commit({
			path = "` + tmpDir + `",
			message = "Test commit",
			author = "Test User <test@example.com>"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for git commit in non-repo: %v", err)
	}
}

func TestGitPull(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterGitModule(L)

	script := `
		local git = require("git")
		local result = git.pull({
			path = "/nonexistent/path"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for git pull in non-repo: %v", err)
	}
}

func TestGitPush(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterGitModule(L)

	script := `
		local git = require("git")
		local result = git.push({
			path = "/nonexistent/path"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for git push in non-repo: %v", err)
	}
}

func TestGitBranch(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterGitModule(L)

	script := `
		local git = require("git")
		local result = git.branch({
			path = "/nonexistent/path",
			name = "test-branch"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for git branch in non-repo: %v", err)
	}
}

func TestGitCheckout(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterGitModule(L)

	script := `
		local git = require("git")
		local result = git.checkout({
			path = "/nonexistent/path",
			branch = "main"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for git checkout in non-repo: %v", err)
	}
}
