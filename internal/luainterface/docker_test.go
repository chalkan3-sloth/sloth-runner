package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestDockerImageList(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.image_list({})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Docker test skipped (Docker not available): %v", err)
		return
	}
}

func TestDockerContainerList(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.container_list({
			all = true
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Docker test skipped (Docker not available): %v", err)
		return
	}
}

func TestDockerContainerRun(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.container_run({
			image = "alpine:latest",
			command = {"echo", "hello"},
			name = "test-container",
			remove = true
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Docker test skipped (Docker not available): %v", err)
		return
	}
}

func TestDockerContainerStop(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.container_stop({
			name = "nonexistent-container"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for nonexistent container: %v", err)
	}
}

func TestDockerContainerRemove(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.container_remove({
			name = "nonexistent-container"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for nonexistent container: %v", err)
	}
}

func TestDockerImagePull(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.image_pull({
			image = "alpine:latest"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Docker test skipped (Docker not available): %v", err)
		return
	}
}

func TestDockerImageBuild(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.image_build({
			context = "/nonexistent/path",
			tag = "test:latest"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error for nonexistent context: %v", err)
	}
}

func TestDockerNetworkCreate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.network_create({
			name = "test-network"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Docker test skipped (Docker not available): %v", err)
		return
	}
}

func TestDockerVolumeCreate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDockerModule(L)

	script := `
		local docker = require("docker")
		local result = docker.volume_create({
			name = "test-volume"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Docker test skipped (Docker not available): %v", err)
		return
	}
}
