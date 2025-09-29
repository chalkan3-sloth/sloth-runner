package luainterface

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// KubernetesModule provides advanced Kubernetes integration
type KubernetesModule struct{}

// NewKubernetesModule creates a new KubernetesModule
func NewKubernetesModule() *KubernetesModule {
	return &KubernetesModule{}
}

// Loader returns the Lua loader for the kubernetes module
func (mod *KubernetesModule) Loader(L *lua.LState) int {
	k8sTable := L.NewTable()
	L.SetFuncs(k8sTable, map[string]lua.LGFunction{
		"apply":             mod.k8sApply,
		"delete":            mod.k8sDelete,
		"get":               mod.k8sGet,
		"describe":          mod.k8sDescribe,
		"logs":              mod.k8sLogs,
		"exec":              mod.k8sExec,
		"port_forward":      mod.k8sPortForward,
		"scale":             mod.k8sScale,
		"rollout":           mod.k8sRollout,
		"patch":             mod.k8sPatch,
		"label":             mod.k8sLabel,
		"annotate":          mod.k8sAnnotate,
		"create_namespace":  mod.k8sCreateNamespace,
		"delete_namespace":  mod.k8sDeleteNamespace,
		"get_nodes":         mod.k8sGetNodes,
		"get_pods":          mod.k8sGetPods,
		"get_services":      mod.k8sGetServices,
		"get_deployments":   mod.k8sGetDeployments,
		"create_secret":     mod.k8sCreateSecret,
		"create_configmap":  mod.k8sCreateConfigMap,
		"wait_for_ready":    mod.k8sWaitForReady,
		"top":               mod.k8sTop,
		"events":            mod.k8sEvents,
	})
	L.Push(k8sTable)
	return 1
}

// k8sApply applies Kubernetes resources
func (mod *KubernetesModule) k8sApply(L *lua.LState) int {
	manifest := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	dryRun := lua.LVAsBool(opts.RawGetString("dry_run"))
	
	args := []string{"apply"}
	
	if strings.HasPrefix(manifest, "http") || strings.Contains(manifest, ".yaml") || strings.Contains(manifest, ".yml") {
		args = append(args, "-f", manifest)
	} else {
		// Treat as inline YAML
		args = append(args, "-f", "-")
	}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if dryRun {
		args = append(args, "--dry-run=client")
	}
	
	var cmd *exec.Cmd
	if strings.HasPrefix(manifest, "http") || strings.Contains(manifest, ".yaml") || strings.Contains(manifest, ".yml") {
		cmd = exec.Command("kubectl", args...)
	} else {
		cmd = exec.Command("kubectl", args...)
		cmd.Stdin = strings.NewReader(manifest)
	}
	
	result, err := mod.executeKubectlCommand(cmd)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sDelete deletes Kubernetes resources
func (mod *KubernetesModule) k8sDelete(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	force := lua.LVAsBool(opts.RawGetString("force"))
	gracePeriod := opts.RawGetString("grace_period").String()
	
	args := []string{"delete", resourceType, resourceName}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if force {
		args = append(args, "--force")
	}
	
	if gracePeriod != "" {
		args = append(args, "--grace-period", gracePeriod)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sGet gets Kubernetes resources
func (mod *KubernetesModule) k8sGet(L *lua.LState) int {
	resourceType := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	resourceName := opts.RawGetString("name").String()
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	output := opts.RawGetString("output").String()
	selector := opts.RawGetString("selector").String()
	allNamespaces := lua.LVAsBool(opts.RawGetString("all_namespaces"))
	
	args := []string{"get", resourceType}
	
	if resourceName != "" {
		args = append(args, resourceName)
	}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if output != "" {
		args = append(args, "-o", output)
	} else {
		args = append(args, "-o", "json")
	}
	
	if selector != "" {
		args = append(args, "-l", selector)
	}
	
	if allNamespaces {
		args = append(args, "--all-namespaces")
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sDescribe describes Kubernetes resources
func (mod *KubernetesModule) k8sDescribe(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	
	args := []string{"describe", resourceType, resourceName}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sLogs gets logs from pods
func (mod *KubernetesModule) k8sLogs(L *lua.LState) int {
	podName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	container := opts.RawGetString("container").String()
	follow := lua.LVAsBool(opts.RawGetString("follow"))
	tail := opts.RawGetString("tail").String()
	since := opts.RawGetString("since").String()
	
	args := []string{"logs", podName}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if container != "" {
		args = append(args, "-c", container)
	}
	
	if follow {
		args = append(args, "-f")
	}
	
	if tail != "" {
		args = append(args, "--tail", tail)
	}
	
	if since != "" {
		args = append(args, "--since", since)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sExec executes commands in pods
func (mod *KubernetesModule) k8sExec(L *lua.LState) int {
	podName := L.CheckString(1)
	command := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	container := opts.RawGetString("container").String()
	stdin := lua.LVAsBool(opts.RawGetString("stdin"))
	tty := lua.LVAsBool(opts.RawGetString("tty"))
	
	args := []string{"exec", podName}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if container != "" {
		args = append(args, "-c", container)
	}
	
	if stdin {
		args = append(args, "-i")
	}
	
	if tty {
		args = append(args, "-t")
	}
	
	args = append(args, "--", "sh", "-c", command)
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sPortForward forwards ports from pods or services
func (mod *KubernetesModule) k8sPortForward(L *lua.LState) int {
	resourceType := L.CheckString(1) // pod or service
	resourceName := L.CheckString(2)
	ports := L.CheckString(3) // "8080:80" or "8080"
	
	opts := L.OptTable(4, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	
	args := []string{"port-forward", resourceType + "/" + resourceName, ports}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	// Port forwarding is typically a long-running operation
	cmd := exec.Command("kubectl", args...)
	err := cmd.Start()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Port forwarding started for %s/%s on %s", resourceType, resourceName, ports)))
	return 2
}

// k8sScale scales deployments, replicasets, or statefulsets
func (mod *KubernetesModule) k8sScale(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	replicas := L.CheckString(3)
	
	opts := L.OptTable(4, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	
	args := []string{"scale", resourceType, resourceName, "--replicas=" + replicas}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sRollout manages rollouts
func (mod *KubernetesModule) k8sRollout(L *lua.LState) int {
	action := L.CheckString(1) // status, history, undo, restart
	resourceType := L.CheckString(2)
	resourceName := L.CheckString(3)
	
	opts := L.OptTable(4, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	revision := opts.RawGetString("revision").String()
	
	args := []string{"rollout", action, resourceType + "/" + resourceName}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if action == "undo" && revision != "" {
		args = append(args, "--to-revision", revision)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sPatch patches Kubernetes resources
func (mod *KubernetesModule) k8sPatch(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	patch := L.CheckString(3)
	
	opts := L.OptTable(4, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	patchType := opts.RawGetString("type").String()
	
	if patchType == "" {
		patchType = "strategic"
	}
	
	args := []string{"patch", resourceType, resourceName, "--patch", patch, "--type", patchType}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sLabel adds or updates labels
func (mod *KubernetesModule) k8sLabel(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	labels := L.CheckString(3) // "key1=value1 key2=value2"
	
	opts := L.OptTable(4, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	overwrite := lua.LVAsBool(opts.RawGetString("overwrite"))
	
	args := []string{"label", resourceType, resourceName}
	labelParts := strings.Fields(labels)
	args = append(args, labelParts...)
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if overwrite {
		args = append(args, "--overwrite")
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sAnnotate adds or updates annotations
func (mod *KubernetesModule) k8sAnnotate(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	annotations := L.CheckString(3) // "key1=value1 key2=value2"
	
	opts := L.OptTable(4, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	overwrite := lua.LVAsBool(opts.RawGetString("overwrite"))
	
	args := []string{"annotate", resourceType, resourceName}
	annotationParts := strings.Fields(annotations)
	args = append(args, annotationParts...)
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if overwrite {
		args = append(args, "--overwrite")
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sCreateNamespace creates a new namespace
func (mod *KubernetesModule) k8sCreateNamespace(L *lua.LState) int {
	namespaceName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	context := opts.RawGetString("context").String()
	dryRun := lua.LVAsBool(opts.RawGetString("dry_run"))
	
	args := []string{"create", "namespace", namespaceName}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if dryRun {
		args = append(args, "--dry-run=client")
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sDeleteNamespace deletes a namespace
func (mod *KubernetesModule) k8sDeleteNamespace(L *lua.LState) int {
	namespaceName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	context := opts.RawGetString("context").String()
	force := lua.LVAsBool(opts.RawGetString("force"))
	
	args := []string{"delete", "namespace", namespaceName}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if force {
		args = append(args, "--force")
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// Helper functions for common resource queries

// k8sGetNodes gets cluster nodes
func (mod *KubernetesModule) k8sGetNodes(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	context := opts.RawGetString("context").String()
	
	args := []string{"get", "nodes", "-o", "json"}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sGetPods gets pods
func (mod *KubernetesModule) k8sGetPods(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	selector := opts.RawGetString("selector").String()
	allNamespaces := lua.LVAsBool(opts.RawGetString("all_namespaces"))
	
	args := []string{"get", "pods", "-o", "json"}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if selector != "" {
		args = append(args, "-l", selector)
	}
	
	if allNamespaces {
		args = append(args, "--all-namespaces")
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sGetServices gets services
func (mod *KubernetesModule) k8sGetServices(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	
	args := []string{"get", "services", "-o", "json"}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sGetDeployments gets deployments
func (mod *KubernetesModule) k8sGetDeployments(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	
	args := []string{"get", "deployments", "-o", "json"}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sCreateSecret creates a secret
func (mod *KubernetesModule) k8sCreateSecret(L *lua.LState) int {
	secretType := L.CheckString(1) // generic, docker-registry, tls
	secretName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	fromLiteral := opts.RawGetString("from_literal").String()
	fromFile := opts.RawGetString("from_file").String()
	
	args := []string{"create", "secret", secretType, secretName}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if fromLiteral != "" {
		args = append(args, "--from-literal", fromLiteral)
	}
	
	if fromFile != "" {
		args = append(args, "--from-file", fromFile)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sCreateConfigMap creates a configmap
func (mod *KubernetesModule) k8sCreateConfigMap(L *lua.LState) int {
	configMapName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	fromLiteral := opts.RawGetString("from_literal").String()
	fromFile := opts.RawGetString("from_file").String()
	
	args := []string{"create", "configmap", configMapName}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if fromLiteral != "" {
		args = append(args, "--from-literal", fromLiteral)
	}
	
	if fromFile != "" {
		args = append(args, "--from-file", fromFile)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sWaitForReady waits for resources to be ready
func (mod *KubernetesModule) k8sWaitForReady(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	timeout := opts.RawGetString("timeout").String()
	
	if timeout == "" {
		timeout = "300s"
	}
	
	args := []string{"wait", "--for=condition=ready", resourceType + "/" + resourceName, "--timeout=" + timeout}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// k8sTop gets resource usage statistics
func (mod *KubernetesModule) k8sTop(L *lua.LState) int {
	resourceType := L.CheckString(1) // node or pod
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	
	args := []string{"top", resourceType}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// k8sEvents gets cluster events
func (mod *KubernetesModule) k8sEvents(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	context := opts.RawGetString("context").String()
	watch := lua.LVAsBool(opts.RawGetString("watch"))
	
	args := []string{"get", "events"}
	
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	
	if context != "" {
		args = append(args, "--context", context)
	}
	
	if watch {
		args = append(args, "--watch")
	}
	
	args = append(args, "--sort-by=.firstTimestamp")
	
	result, err := mod.executeKubectlCommand(exec.Command("kubectl", args...))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// executeKubectlCommand executes a kubectl command
func (mod *KubernetesModule) executeKubectlCommand(cmd *exec.Cmd) (string, error) {
	// Check if kubectl command exists
	if _, err := exec.LookPath("kubectl"); err != nil {
		return "", fmt.Errorf("kubectl command not found in PATH: %w", err)
	}
	
	// Set environment variables
	cmd.Env = os.Environ()
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Set timeout
	timeout := 300 * time.Second // 5 minutes default
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			errorMsg := stderr.String()
			if errorMsg == "" {
				errorMsg = err.Error()
			}
			return "", fmt.Errorf("kubectl command failed: %s", errorMsg)
		}
		return stdout.String(), nil
		
	case <-timer.C:
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("kubectl command timed out after %v", timeout)
	}
}