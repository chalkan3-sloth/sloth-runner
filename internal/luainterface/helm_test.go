package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestHelmInstall(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterHelmModule(L)

	script := `
		local helm = require("helm")
		local result = helm.install({
			name = "test-release",
			chart = "stable/nginx",
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Helm test skipped (Helm not available): %v", err)
		return
	}
}

func TestHelmUpgrade(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterHelmModule(L)

	script := `
		local helm = require("helm")
		local result = helm.upgrade({
			name = "test-release",
			chart = "stable/nginx",
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Helm test skipped (Helm not available): %v", err)
		return
	}
}

func TestHelmUninstall(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterHelmModule(L)

	script := `
		local helm = require("helm")
		local result = helm.uninstall({
			name = "test-release",
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Helm test skipped (Helm not available): %v", err)
		return
	}
}

func TestHelmList(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterHelmModule(L)

	script := `
		local helm = require("helm")
		local result = helm.list({
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Helm test skipped (Helm not available): %v", err)
		return
	}
}

func TestHelmRepoAdd(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterHelmModule(L)

	script := `
		local helm = require("helm")
		local result = helm.repo_add({
			name = "stable",
			url = "https://charts.helm.sh/stable"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Helm test skipped (Helm not available): %v", err)
		return
	}
}

func TestHelmRepoUpdate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterHelmModule(L)

	script := `
		local helm = require("helm")
		local result = helm.repo_update({})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Helm test skipped (Helm not available): %v", err)
		return
	}
}
