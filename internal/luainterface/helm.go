package luainterface

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// HelmModule provides advanced Helm integration
type HelmModule struct{}

// NewHelmModule creates a new HelmModule
func NewHelmModule() *HelmModule {
	return &HelmModule{}
}

// Loader returns the Lua loader for the helm module
func (mod *HelmModule) Loader(L *lua.LState) int {
	helmTable := L.NewTable()
	L.SetFuncs(helmTable, map[string]lua.LGFunction{
		"install":         mod.helmInstall,
		"upgrade":         mod.helmUpgrade,
		"uninstall":       mod.helmUninstall,
		"list":            mod.helmList,
		"status":          mod.helmStatus,
		"get":             mod.helmGet,
		"rollback":        mod.helmRollback,
		"history":         mod.helmHistory,
		"test":            mod.helmTest,
		"lint":            mod.helmLint,
		"template":        mod.helmTemplate,
		"dependency":      mod.helmDependency,
		"package":         mod.helmPackage,
		"repo_add":        mod.helmRepoAdd,
		"repo_update":     mod.helmRepoUpdate,
		"repo_list":       mod.helmRepoList,
		"repo_remove":     mod.helmRepoRemove,
		"search_repo":     mod.helmSearchRepo,
		"search_hub":      mod.helmSearchHub,
		"create":          mod.helmCreate,
		"plugin_install":  mod.helmPluginInstall,
		"plugin_list":     mod.helmPluginList,
		"plugin_uninstall": mod.helmPluginUninstall,
		"version":         mod.helmVersion,
	})
	L.Push(helmTable)
	return 1
}

// helmInstall installs a Helm chart
func (mod *HelmModule) helmInstall(L *lua.LState) int {
	releaseName := L.CheckString(1)
	chart := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	createNamespace := lua.LVAsBool(opts.RawGetString("create_namespace"))
	version := opts.RawGetString("version").String()
	values := opts.RawGetString("values").String()
	valuesFile := opts.RawGetString("values_file").String()
	wait := lua.LVAsBool(opts.RawGetString("wait"))
	timeout := opts.RawGetString("timeout").String()
	dryRun := lua.LVAsBool(opts.RawGetString("dry_run"))
	atomic := lua.LVAsBool(opts.RawGetString("atomic"))
	
	args := []string{"install", releaseName, chart}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if createNamespace {
		args = append(args, "--create-namespace")
	}
	
	if version != "" {
		args = append(args, "--version", version)
	}
	
	if values != "" {
		args = append(args, "--set", values)
	}
	
	if valuesFile != "" {
		args = append(args, "--values", valuesFile)
	}
	
	if wait {
		args = append(args, "--wait")
	}
	
	if timeout != "" {
		args = append(args, "--timeout", timeout)
	}
	
	if dryRun {
		args = append(args, "--dry-run")
	}
	
	if atomic {
		args = append(args, "--atomic")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmUpgrade upgrades a Helm release
func (mod *HelmModule) helmUpgrade(L *lua.LState) int {
	releaseName := L.CheckString(1)
	chart := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	version := opts.RawGetString("version").String()
	values := opts.RawGetString("values").String()
	valuesFile := opts.RawGetString("values_file").String()
	wait := lua.LVAsBool(opts.RawGetString("wait"))
	timeout := opts.RawGetString("timeout").String()
	install := lua.LVAsBool(opts.RawGetString("install"))
	atomic := lua.LVAsBool(opts.RawGetString("atomic"))
	force := lua.LVAsBool(opts.RawGetString("force"))
	resetValues := lua.LVAsBool(opts.RawGetString("reset_values"))
	
	args := []string{"upgrade", releaseName, chart}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if version != "" {
		args = append(args, "--version", version)
	}
	
	if values != "" {
		args = append(args, "--set", values)
	}
	
	if valuesFile != "" {
		args = append(args, "--values", valuesFile)
	}
	
	if wait {
		args = append(args, "--wait")
	}
	
	if timeout != "" {
		args = append(args, "--timeout", timeout)
	}
	
	if install {
		args = append(args, "--install")
	}
	
	if atomic {
		args = append(args, "--atomic")
	}
	
	if force {
		args = append(args, "--force")
	}
	
	if resetValues {
		args = append(args, "--reset-values")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmUninstall uninstalls a Helm release
func (mod *HelmModule) helmUninstall(L *lua.LState) int {
	releaseName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	keepHistory := lua.LVAsBool(opts.RawGetString("keep_history"))
	dryRun := lua.LVAsBool(opts.RawGetString("dry_run"))
	
	args := []string{"uninstall", releaseName}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if keepHistory {
		args = append(args, "--keep-history")
	}
	
	if dryRun {
		args = append(args, "--dry-run")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmList lists Helm releases
func (mod *HelmModule) helmList(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	allNamespaces := lua.LVAsBool(opts.RawGetString("all_namespaces"))
	filter := opts.RawGetString("filter").String()
	output := opts.RawGetString("output").String()
	
	args := []string{"list"}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if allNamespaces {
		args = append(args, "--all-namespaces")
	}
	
	if filter != "" {
		args = append(args, "--filter", filter)
	}
	
	if output != "" {
		args = append(args, "--output", output)
	} else {
		args = append(args, "--output", "json")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmStatus gets the status of a release
func (mod *HelmModule) helmStatus(L *lua.LState) int {
	releaseName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	output := opts.RawGetString("output").String()
	
	args := []string{"status", releaseName}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if output != "" {
		args = append(args, "--output", output)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmGet gets information about a release
func (mod *HelmModule) helmGet(L *lua.LState) int {
	subcommand := L.CheckString(1) // all, hooks, manifest, notes, values
	releaseName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	revision := opts.RawGetString("revision").String()
	output := opts.RawGetString("output").String()
	
	args := []string{"get", subcommand, releaseName}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if revision != "" {
		args = append(args, "--revision", revision)
	}
	
	if output != "" {
		args = append(args, "--output", output)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmRollback rollbacks a release to a previous revision
func (mod *HelmModule) helmRollback(L *lua.LState) int {
	releaseName := L.CheckString(1)
	revision := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	wait := lua.LVAsBool(opts.RawGetString("wait"))
	timeout := opts.RawGetString("timeout").String()
	dryRun := lua.LVAsBool(opts.RawGetString("dry_run"))
	force := lua.LVAsBool(opts.RawGetString("force"))
	
	args := []string{"rollback", releaseName, revision}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if wait {
		args = append(args, "--wait")
	}
	
	if timeout != "" {
		args = append(args, "--timeout", timeout)
	}
	
	if dryRun {
		args = append(args, "--dry-run")
	}
	
	if force {
		args = append(args, "--force")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmHistory shows release history
func (mod *HelmModule) helmHistory(L *lua.LState) int {
	releaseName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	max := opts.RawGetString("max").String()
	output := opts.RawGetString("output").String()
	
	args := []string{"history", releaseName}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if max != "" {
		args = append(args, "--max", max)
	}
	
	if output != "" {
		args = append(args, "--output", output)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmTest runs tests for a release
func (mod *HelmModule) helmTest(L *lua.LState) int {
	releaseName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	timeout := opts.RawGetString("timeout").String()
	logs := lua.LVAsBool(opts.RawGetString("logs"))
	
	args := []string{"test", releaseName}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if timeout != "" {
		args = append(args, "--timeout", timeout)
	}
	
	if logs {
		args = append(args, "--logs")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmLint lints a chart
func (mod *HelmModule) helmLint(L *lua.LState) int {
	chartPath := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	strict := lua.LVAsBool(opts.RawGetString("strict"))
	values := opts.RawGetString("values").String()
	valuesFile := opts.RawGetString("values_file").String()
	
	args := []string{"lint", chartPath}
	
	if strict {
		args = append(args, "--strict")
	}
	
	if values != "" {
		args = append(args, "--set", values)
	}
	
	if valuesFile != "" {
		args = append(args, "--values", valuesFile)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmTemplate renders chart templates locally
func (mod *HelmModule) helmTemplate(L *lua.LState) int {
	releaseName := L.CheckString(1)
	chart := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	namespace := opts.RawGetString("namespace").String()
	values := opts.RawGetString("values").String()
	valuesFile := opts.RawGetString("values_file").String()
	outputDir := opts.RawGetString("output_dir").String()
	showOnly := opts.RawGetString("show_only").String()
	
	args := []string{"template", releaseName, chart}
	
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	
	if values != "" {
		args = append(args, "--set", values)
	}
	
	if valuesFile != "" {
		args = append(args, "--values", valuesFile)
	}
	
	if outputDir != "" {
		args = append(args, "--output-dir", outputDir)
	}
	
	if showOnly != "" {
		args = append(args, "--show-only", showOnly)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmDependency manages chart dependencies
func (mod *HelmModule) helmDependency(L *lua.LState) int {
	action := L.CheckString(1) // build, list, update
	chartPath := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	skipRefresh := lua.LVAsBool(opts.RawGetString("skip_refresh"))
	
	args := []string{"dependency", action, chartPath}
	
	if action == "update" && skipRefresh {
		args = append(args, "--skip-refresh")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmPackage packages a chart
func (mod *HelmModule) helmPackage(L *lua.LState) int {
	chartPath := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	destination := opts.RawGetString("destination").String()
	version := opts.RawGetString("version").String()
	appVersion := opts.RawGetString("app_version").String()
	sign := lua.LVAsBool(opts.RawGetString("sign"))
	
	args := []string{"package", chartPath}
	
	if destination != "" {
		args = append(args, "--destination", destination)
	}
	
	if version != "" {
		args = append(args, "--version", version)
	}
	
	if appVersion != "" {
		args = append(args, "--app-version", appVersion)
	}
	
	if sign {
		args = append(args, "--sign")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmRepoAdd adds a chart repository
func (mod *HelmModule) helmRepoAdd(L *lua.LState) int {
	repoName := L.CheckString(1)
	repoURL := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	username := opts.RawGetString("username").String()
	password := opts.RawGetString("password").String()
	forceUpdate := lua.LVAsBool(opts.RawGetString("force_update"))
	
	args := []string{"repo", "add", repoName, repoURL}
	
	if username != "" {
		args = append(args, "--username", username)
	}
	
	if password != "" {
		args = append(args, "--password", password)
	}
	
	if forceUpdate {
		args = append(args, "--force-update")
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmRepoUpdate updates chart repositories
func (mod *HelmModule) helmRepoUpdate(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	repoName := opts.RawGetString("repo").String()
	
	args := []string{"repo", "update"}
	
	if repoName != "" {
		args = append(args, repoName)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmRepoList lists chart repositories
func (mod *HelmModule) helmRepoList(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	output := opts.RawGetString("output").String()
	
	args := []string{"repo", "list"}
	
	if output != "" {
		args = append(args, "--output", output)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmRepoRemove removes a chart repository
func (mod *HelmModule) helmRepoRemove(L *lua.LState) int {
	repoName := L.CheckString(1)
	
	result, err := mod.executeHelmCommand("repo", "remove", repoName)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmSearchRepo searches chart repositories
func (mod *HelmModule) helmSearchRepo(L *lua.LState) int {
	keyword := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	version := opts.RawGetString("version").String()
	maxResults := opts.RawGetString("max_results").String()
	output := opts.RawGetString("output").String()
	
	args := []string{"search", "repo", keyword}
	
	if version != "" {
		args = append(args, "--version", version)
	}
	
	if maxResults != "" {
		args = append(args, "--max-col-width", maxResults)
	}
	
	if output != "" {
		args = append(args, "--output", output)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmSearchHub searches Helm Hub
func (mod *HelmModule) helmSearchHub(L *lua.LState) int {
	keyword := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	maxResults := opts.RawGetString("max_results").String()
	output := opts.RawGetString("output").String()
	
	args := []string{"search", "hub", keyword}
	
	if maxResults != "" {
		args = append(args, "--max-col-width", maxResults)
	}
	
	if output != "" {
		args = append(args, "--output", output)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmCreate creates a new chart
func (mod *HelmModule) helmCreate(L *lua.LState) int {
	chartName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	starter := opts.RawGetString("starter").String()
	
	args := []string{"create", chartName}
	
	if starter != "" {
		args = append(args, "--starter", starter)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmPluginInstall installs a Helm plugin
func (mod *HelmModule) helmPluginInstall(L *lua.LState) int {
	pluginPath := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	version := opts.RawGetString("version").String()
	
	args := []string{"plugin", "install", pluginPath}
	
	if version != "" {
		args = append(args, "--version", version)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmPluginList lists installed plugins
func (mod *HelmModule) helmPluginList(L *lua.LState) int {
	result, err := mod.executeHelmCommand("plugin", "list")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// helmPluginUninstall uninstalls a plugin
func (mod *HelmModule) helmPluginUninstall(L *lua.LState) int {
	pluginName := L.CheckString(1)
	
	result, err := mod.executeHelmCommand("plugin", "uninstall", pluginName)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// helmVersion shows Helm version
func (mod *HelmModule) helmVersion(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	short := lua.LVAsBool(opts.RawGetString("short"))
	template := opts.RawGetString("template").String()
	
	args := []string{"version"}
	
	if short {
		args = append(args, "--short")
	}
	
	if template != "" {
		args = append(args, "--template", template)
	}
	
	result, err := mod.executeHelmCommand(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// executeHelmCommand executes a helm command
func (mod *HelmModule) executeHelmCommand(cmdArgs ...string) (string, error) {
	// Check if helm command exists
	if _, err := exec.LookPath("helm"); err != nil {
		return "", fmt.Errorf("helm command not found in PATH: %w", err)
	}
	
	cmd := exec.Command("helm", cmdArgs...)
	
	// Set environment variables
	cmd.Env = os.Environ()
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Set timeout
	timeout := 600 * time.Second // 10 minutes for long operations
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
			return "", fmt.Errorf("helm command failed: %s", errorMsg)
		}
		return stdout.String(), nil
		
	case <-timer.C:
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("helm command timed out after %v", timeout)
	}
}