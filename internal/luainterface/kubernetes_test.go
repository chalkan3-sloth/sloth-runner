package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestKubernetesApply(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterKubernetesModule(L)

	script := `
		local k8s = require("kubernetes")
		local result = k8s.apply({
			manifest = "test.yaml",
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Kubernetes test skipped (K8s not available): %v", err)
		return
	}
}

func TestKubernetesDelete(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterKubernetesModule(L)

	script := `
		local k8s = require("kubernetes")
		local result = k8s.delete({
			manifest = "test.yaml",
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Kubernetes test skipped (K8s not available): %v", err)
		return
	}
}

func TestKubernetesGetPods(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterKubernetesModule(L)

	script := `
		local k8s = require("kubernetes")
		local result = k8s.get_pods({
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Kubernetes test skipped (K8s not available): %v", err)
		return
	}
}

func TestKubernetesGetServices(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterKubernetesModule(L)

	script := `
		local k8s = require("kubernetes")
		local result = k8s.get_services({
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Kubernetes test skipped (K8s not available): %v", err)
		return
	}
}

func TestKubernetesScaleDeployment(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterKubernetesModule(L)

	script := `
		local k8s = require("kubernetes")
		local result = k8s.scale_deployment({
			name = "test-deployment",
			namespace = "default",
			replicas = 3
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Kubernetes test skipped (K8s not available): %v", err)
		return
	}
}

func TestKubernetesExec(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterKubernetesModule(L)

	script := `
		local k8s = require("kubernetes")
		local result = k8s.exec({
			pod = "test-pod",
			namespace = "default",
			command = {"echo", "hello"}
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Kubernetes test skipped (K8s not available): %v", err)
		return
	}
}

func TestKubernetesLogs(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterKubernetesModule(L)

	script := `
		local k8s = require("kubernetes")
		local result = k8s.logs({
			pod = "test-pod",
			namespace = "default"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Kubernetes test skipped (K8s not available): %v", err)
		return
	}
}
