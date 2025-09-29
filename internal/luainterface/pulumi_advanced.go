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

// PulumiModule provides advanced Pulumi integration
type PulumiModule struct{}

// NewPulumiModule creates a new PulumiModule
func NewPulumiModule() *PulumiModule {
	return &PulumiModule{}
}

// Loader returns the Lua loader for the pulumi module
func (mod *PulumiModule) Loader(L *lua.LState) int {
	pulumiTable := L.NewTable()
	L.SetFuncs(pulumiTable, map[string]lua.LGFunction{
		"stack":           mod.pulumiStack,
		"new_stack":       mod.pulumiNewStack,
		"select_stack":    mod.pulumiSelectStack,
		"list_stacks":     mod.pulumiListStacks,
		"destroy_stack":   mod.pulumiDestroyStack,
		"up":              mod.pulumiUp,
		"preview":         mod.pulumiPreview,
		"refresh":         mod.pulumiRefresh,
		"destroy":         mod.pulumiDestroy,
		"config_set":      mod.pulumiConfigSet,
		"config_get":      mod.pulumiConfigGet,
		"config_rm":       mod.pulumiConfigRm,
		"outputs":         mod.pulumiOutputs,
		"logs":            mod.pulumiLogs,
		"history":         mod.pulumiHistory,
		"export":          mod.pulumiExport,
		"import":          mod.pulumiImport,
		"plugin_install":  mod.pulumiPluginInstall,
		"plugin_ls":       mod.pulumiPluginLs,
		"policy_new":      mod.pulumiPolicyNew,
		"policy_enable":   mod.pulumiPolicyEnable,
		"policy_disable":  mod.pulumiPolicyDisable,
		"login":           mod.pulumiLogin,
		"logout":          mod.pulumiLogout,
		"whoami":          mod.pulumiWhoami,
	})
	L.Push(pulumiTable)
	return 1
}

// pulumiStack creates or selects a Pulumi stack (alias for new_stack)
func (mod *PulumiModule) pulumiStack(L *lua.LState) int {
	return mod.pulumiNewStack(L)
}

// pulumiNewStack creates a new Pulumi stack
func (mod *PulumiModule) pulumiNewStack(L *lua.LState) int {
	stackName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	template := opts.RawGetString("template").String()
	backend := opts.RawGetString("backend").String()
	
	args := []string{"stack", "init", stackName}
	
	if template != "" {
		args = append(args, "--template", template)
	}
	
	env := make(map[string]string)
	if backend != "" {
		env["PULUMI_BACKEND_URL"] = backend
	}
	
	result, err := mod.executePulumiCommand(workdir, env, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiSelectStack selects an existing stack
func (mod *PulumiModule) pulumiSelectStack(L *lua.LState) int {
	stackName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "stack", "select", stackName)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiListStacks lists all stacks
func (mod *PulumiModule) pulumiListStacks(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "stack", "ls", "--json")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// pulumiDestroyStack destroys a stack
func (mod *PulumiModule) pulumiDestroyStack(L *lua.LState) int {
	stackName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	force := opts.RawGetString("force")
	
	args := []string{"stack", "rm", stackName}
	if lua.LVAsBool(force) {
		args = append(args, "--force")
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiUp performs a pulumi up operation
func (mod *PulumiModule) pulumiUp(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"up"}
	
	// Add options
	if skipPreview := opts.RawGetString("skip_preview"); lua.LVAsBool(skipPreview) {
		args = append(args, "--skip-preview")
	}
	
	if yes := opts.RawGetString("yes"); lua.LVAsBool(yes) {
		args = append(args, "--yes")
	}
	
	if parallel := opts.RawGetString("parallel"); parallel != lua.LNil {
		args = append(args, "--parallel", parallel.String())
	}
	
	if refresh := opts.RawGetString("refresh"); lua.LVAsBool(refresh) {
		args = append(args, "--refresh")
	}
	
	if diff := opts.RawGetString("diff"); lua.LVAsBool(diff) {
		args = append(args, "--diff")
	}
	
	// Add JSON output for programmatic processing
	args = append(args, "--json")
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiPreview performs a preview operation
func (mod *PulumiModule) pulumiPreview(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"preview", "--json"}
	
	if diff := opts.RawGetString("diff"); lua.LVAsBool(diff) {
		args = append(args, "--diff")
	}
	
	if refresh := opts.RawGetString("refresh"); lua.LVAsBool(refresh) {
		args = append(args, "--refresh")
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// pulumiRefresh refreshes the stack state
func (mod *PulumiModule) pulumiRefresh(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"refresh"}
	
	if yes := opts.RawGetString("yes"); lua.LVAsBool(yes) {
		args = append(args, "--yes")
	}
	
	if skipPreview := opts.RawGetString("skip_preview"); lua.LVAsBool(skipPreview) {
		args = append(args, "--skip-preview")
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiDestroy destroys stack resources
func (mod *PulumiModule) pulumiDestroy(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"destroy"}
	
	if yes := opts.RawGetString("yes"); lua.LVAsBool(yes) {
		args = append(args, "--yes")
	}
	
	if skipPreview := opts.RawGetString("skip_preview"); lua.LVAsBool(skipPreview) {
		args = append(args, "--skip-preview")
	}
	
	if parallel := opts.RawGetString("parallel"); parallel != lua.LNil {
		args = append(args, "--parallel", parallel.String())
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiConfigSet sets configuration values
func (mod *PulumiModule) pulumiConfigSet(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"config", "set", key, value}
	
	if secret := opts.RawGetString("secret"); lua.LVAsBool(secret) {
		args = append(args, "--secret")
	}
	
	if plaintext := opts.RawGetString("plaintext"); lua.LVAsBool(plaintext) {
		args = append(args, "--plaintext")
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiConfigGet gets configuration values
func (mod *PulumiModule) pulumiConfigGet(L *lua.LState) int {
	key := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"config", "get", key}
	
	if json := opts.RawGetString("json"); lua.LVAsBool(json) {
		args = append(args, "--json")
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(strings.TrimSpace(result)))
	L.Push(lua.LNil)
	return 2
}

// pulumiConfigRm removes configuration values
func (mod *PulumiModule) pulumiConfigRm(L *lua.LState) int {
	key := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "config", "rm", key)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiOutputs gets stack outputs
func (mod *PulumiModule) pulumiOutputs(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "stack", "output", "--json")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// pulumiLogs gets stack logs
func (mod *PulumiModule) pulumiLogs(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"logs"}
	
	if follow := opts.RawGetString("follow"); lua.LVAsBool(follow) {
		args = append(args, "--follow")
	}
	
	if since := opts.RawGetString("since"); since != lua.LNil {
		args = append(args, "--since", since.String())
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// pulumiHistory gets stack history
func (mod *PulumiModule) pulumiHistory(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"history", "--json"}
	
	if showSecrets := opts.RawGetString("show_secrets"); lua.LVAsBool(showSecrets) {
		args = append(args, "--show-secrets")
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// pulumiExport exports stack state
func (mod *PulumiModule) pulumiExport(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"stack", "export"}
	
	if file := opts.RawGetString("file"); file != lua.LNil {
		args = append(args, "--file", file.String())
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// pulumiImport imports stack state
func (mod *PulumiModule) pulumiImport(L *lua.LState) int {
	file := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "stack", "import", "--file", file)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiPluginInstall installs a plugin
func (mod *PulumiModule) pulumiPluginInstall(L *lua.LState) int {
	pluginKind := L.CheckString(1) // resource, language, analyzer
	pluginName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"plugin", "install", pluginKind, pluginName}
	
	if version := opts.RawGetString("version"); version != lua.LNil {
		args = append(args, "--version", version.String())
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiPluginLs lists installed plugins
func (mod *PulumiModule) pulumiPluginLs(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "plugin", "ls", "--json")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// pulumiPolicyNew creates a new policy pack
func (mod *PulumiModule) pulumiPolicyNew(L *lua.LState) int {
	policyName := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	template := opts.RawGetString("template").String()
	if template == "" {
		template = "typescript"
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "policy", "new", policyName, "--template", template)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiPolicyEnable enables a policy pack
func (mod *PulumiModule) pulumiPolicyEnable(L *lua.LState) int {
	policyPack := L.CheckString(1)
	orgName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "policy", "enable", policyPack, orgName)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiPolicyDisable disables a policy pack
func (mod *PulumiModule) pulumiPolicyDisable(L *lua.LState) int {
	policyPack := L.CheckString(1)
	orgName := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "policy", "disable", policyPack, orgName)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiLogin logs into Pulumi backend
func (mod *PulumiModule) pulumiLogin(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	backend := opts.RawGetString("backend").String()
	
	args := []string{"login"}
	if backend != "" {
		args = append(args, backend)
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiLogout logs out of Pulumi backend
func (mod *PulumiModule) pulumiLogout(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "logout")
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// pulumiWhoami shows current user info
func (mod *PulumiModule) pulumiWhoami(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executePulumiCommand(workdir, nil, "whoami", "--json")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// executePulumiCommand executes a pulumi command with environment variables
func (mod *PulumiModule) executePulumiCommand(workdir string, env map[string]string, cmdArgs ...string) (string, error) {
	// Check if pulumi command exists
	if _, err := exec.LookPath("pulumi"); err != nil {
		return "", fmt.Errorf("pulumi command not found in PATH: %w", err)
	}
	
	cmd := exec.Command("pulumi", cmdArgs...)
	cmd.Dir = workdir
	
	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Set timeout for long-running operations
	timeout := 600 * time.Second // 10 minutes
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
			return "", fmt.Errorf("pulumi command failed: %s", errorMsg)
		}
		return stdout.String(), nil
		
	case <-timer.C:
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("pulumi command timed out after %v", timeout)
	}
}

// PulumiLoader loads the pulumi module for Lua
func PulumiLoader(L *lua.LState) int {
	mod := NewPulumiModule()
	return mod.Loader(L)
}