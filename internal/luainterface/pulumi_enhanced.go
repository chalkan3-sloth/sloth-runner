package luainterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	lua "github.com/yuin/gopher-lua"
)

// EnhancedPulumiModule provides advanced Pulumi functionality
type EnhancedPulumiModule struct{}

// PulumiStack represents an enhanced Pulumi stack
type PulumiStack struct {
	Name        string
	Project     string
	WorkDir     string
	Venv        *types.PythonVenv
	LoginURL    string
	ConfigFile  string
	Backend     string
	Environment map[string]string
	Secrets     map[string]string
	Tags        map[string]string
}

// PulumiResult represents Pulumi operation result
type PulumiResult struct {
	Success     bool
	ExitCode    int
	Stdout      string
	Stderr      string
	Duration    time.Duration
	Summary     map[string]interface{}
	Resources   []map[string]interface{}
	Outputs     map[string]interface{}
	Permalink   string
}

// NewEnhancedPulumiModule creates a new enhanced Pulumi module
func NewEnhancedPulumiModule() *EnhancedPulumiModule {
	return &EnhancedPulumiModule{}
}

// RegisterEnhancedPulumiModule registers the enhanced Pulumi module
func RegisterEnhancedPulumiModule(L *lua.LState) {
	module := NewEnhancedPulumiModule()
	
	// Create the pulumi table
	pulumiTable := L.NewTable()
	
	// Stack management
	L.SetField(pulumiTable, "stack", L.NewFunction(module.luaNewStack))
	L.SetField(pulumiTable, "new_stack", L.NewFunction(module.luaCreateStack))
	L.SetField(pulumiTable, "list_stacks", L.NewFunction(module.luaListStacks))
	L.SetField(pulumiTable, "select_stack", L.NewFunction(module.luaSelectStack))
	L.SetField(pulumiTable, "remove_stack", L.NewFunction(module.luaRemoveStack))
	
	// Operations
	L.SetField(pulumiTable, "up", L.NewFunction(module.luaUp))
	L.SetField(pulumiTable, "preview", L.NewFunction(module.luaPreview))
	L.SetField(pulumiTable, "destroy", L.NewFunction(module.luaDestroy))
	L.SetField(pulumiTable, "refresh", L.NewFunction(module.luaRefresh))
	L.SetField(pulumiTable, "cancel", L.NewFunction(module.luaCancel))
	
	// Configuration
	L.SetField(pulumiTable, "config_set", L.NewFunction(module.luaConfigSet))
	L.SetField(pulumiTable, "config_get", L.NewFunction(module.luaConfigGet))
	L.SetField(pulumiTable, "config_rm", L.NewFunction(module.luaConfigRm))
	L.SetField(pulumiTable, "config_list", L.NewFunction(module.luaConfigList))
	
	// Outputs and state
	L.SetField(pulumiTable, "outputs", L.NewFunction(module.luaOutputs))
	L.SetField(pulumiTable, "export", L.NewFunction(module.luaExport))
	L.SetField(pulumiTable, "import", L.NewFunction(module.luaImport))
	L.SetField(pulumiTable, "state", L.NewFunction(module.luaState))
	
	// Plugin management
	L.SetField(pulumiTable, "plugin_install", L.NewFunction(module.luaPluginInstall))
	L.SetField(pulumiTable, "plugin_ls", L.NewFunction(module.luaPluginList))
	L.SetField(pulumiTable, "plugin_rm", L.NewFunction(module.luaPluginRemove))
	
	// Advanced features
	L.SetField(pulumiTable, "watch", L.NewFunction(module.luaWatch))
	L.SetField(pulumiTable, "logs", L.NewFunction(module.luaLogs))
	L.SetField(pulumiTable, "history", L.NewFunction(module.luaHistory))
	L.SetField(pulumiTable, "console", L.NewFunction(module.luaConsole))
	
	// Utility functions
	L.SetField(pulumiTable, "version", L.NewFunction(module.luaVersion))
	L.SetField(pulumiTable, "whoami", L.NewFunction(module.luaWhoami))
	L.SetField(pulumiTable, "login", L.NewFunction(module.luaLogin))
	L.SetField(pulumiTable, "logout", L.NewFunction(module.luaLogout))
	
	// Register the pulumi table globally
	L.SetGlobal("pulumi", pulumiTable)
}

// Stack management
func (p *EnhancedPulumiModule) luaNewStack(L *lua.LState) int {
	options := L.CheckTable(1)
	
	stack := &PulumiStack{
		Name:        options.RawGetString("name").String(),
		Project:     options.RawGetString("project").String(),
		WorkDir:     options.RawGetString("workdir").String(),
		LoginURL:    options.RawGetString("login_url").String(),
		ConfigFile:  options.RawGetString("config_file").String(),
		Backend:     options.RawGetString("backend").String(),
		Environment: make(map[string]string),
		Secrets:     make(map[string]string),
		Tags:        make(map[string]string),
	}
	
	// Handle venv
	if venvValue := options.RawGetString("venv"); venvValue.Type() == lua.LTUserData {
		if v, ok := venvValue.(*lua.LUserData).Value.(*types.PythonVenv); ok {
			stack.Venv = v
		}
	}
	
	// Parse environment variables
	if envTable := options.RawGetString("env"); envTable.Type() == lua.LTTable {
		envTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			stack.Environment[lua.LVAsString(key)] = lua.LVAsString(value)
		})
	}
	
	// Parse tags
	if tagsTable := options.RawGetString("tags"); tagsTable.Type() == lua.LTTable {
		tagsTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			stack.Tags[lua.LVAsString(key)] = lua.LVAsString(value)
		})
	}
	
	L.Push(p.stackToLua(L, stack))
	return 1
}

func (p *EnhancedPulumiModule) luaCreateStack(L *lua.LState) int {
	name := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"stack", "init", name}
	
	if template := options.RawGetString("template").String(); template != "" {
		args = append(args, "--template", template)
	}
	
	if description := options.RawGetString("description").String(); description != "" {
		args = append(args, "--description", description)
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaListStacks(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"stack", "ls", "--json"}
	
	if organization := options.RawGetString("organization").String(); organization != "" {
		args = append(args, "--organization", organization)
	}
	
	if project := options.RawGetString("project").String(); project != "" {
		args = append(args, "--project", project)
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaSelectStack(L *lua.LState) int {
	name := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"stack", "select", name}
	
	if create := lua.LVAsBool(options.RawGetString("create")); create {
		args = append(args, "--create")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaRemoveStack(L *lua.LState) int {
	name := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"stack", "rm", name}
	
	if force := lua.LVAsBool(options.RawGetString("force")); force {
		args = append(args, "--force")
	}
	
	if yes := lua.LVAsBool(options.RawGetString("yes")); yes {
		args = append(args, "--yes")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

// Operations
func (p *EnhancedPulumiModule) luaUp(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"up"}
	
	if yes := lua.LVAsBool(options.RawGetString("yes")); yes {
		args = append(args, "--yes")
	}
	
	if skipPreview := lua.LVAsBool(options.RawGetString("skip_preview")); skipPreview {
		args = append(args, "--skip-preview")
	}
	
	if refresh := lua.LVAsBool(options.RawGetString("refresh")); refresh {
		args = append(args, "--refresh")
	}
	
	if diff := lua.LVAsBool(options.RawGetString("diff")); diff {
		args = append(args, "--diff")
	}
	
	if target := options.RawGetString("target").String(); target != "" {
		args = append(args, "--target", target)
	}
	
	if parallel := int(options.RawGetString("parallel").(lua.LNumber)); parallel > 0 {
		args = append(args, "--parallel", fmt.Sprintf("%d", parallel))
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaPreview(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"preview"}
	
	if refresh := lua.LVAsBool(options.RawGetString("refresh")); refresh {
		args = append(args, "--refresh")
	}
	
	if diff := lua.LVAsBool(options.RawGetString("diff")); diff {
		args = append(args, "--diff")
	}
	
	if target := options.RawGetString("target").String(); target != "" {
		args = append(args, "--target", target)
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaDestroy(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"destroy"}
	
	if yes := lua.LVAsBool(options.RawGetString("yes")); yes {
		args = append(args, "--yes")
	}
	
	if skipPreview := lua.LVAsBool(options.RawGetString("skip_preview")); skipPreview {
		args = append(args, "--skip-preview")
	}
	
	if target := options.RawGetString("target").String(); target != "" {
		args = append(args, "--target", target)
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaRefresh(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"refresh"}
	
	if yes := lua.LVAsBool(options.RawGetString("yes")); yes {
		args = append(args, "--yes")
	}
	
	if skipPreview := lua.LVAsBool(options.RawGetString("skip_preview")); skipPreview {
		args = append(args, "--skip-preview")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaCancel(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"cancel"}
	
	if yes := lua.LVAsBool(options.RawGetString("yes")); yes {
		args = append(args, "--yes")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

// Configuration
func (p *EnhancedPulumiModule) luaConfigSet(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)
	options := L.OptTable(3, L.NewTable())
	
	args := []string{"config", "set", key, value}
	
	if secret := lua.LVAsBool(options.RawGetString("secret")); secret {
		args = append(args, "--secret")
	}
	
	if path := lua.LVAsBool(options.RawGetString("path")); path {
		args = append(args, "--path")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaConfigGet(L *lua.LState) int {
	key := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"config", "get", key}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaConfigRm(L *lua.LState) int {
	key := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"config", "rm", key}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaConfigList(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"config", "ls", "--json"}
	
	if showSecrets := lua.LVAsBool(options.RawGetString("show_secrets")); showSecrets {
		args = append(args, "--show-secrets")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

// Outputs and state
func (p *EnhancedPulumiModule) luaOutputs(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"stack", "output", "--json"}
	
	if showSecrets := lua.LVAsBool(options.RawGetString("show_secrets")); showSecrets {
		args = append(args, "--show-secrets")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaExport(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"stack", "export"}
	
	if file := options.RawGetString("file").String(); file != "" {
		args = append(args, "--file", file)
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaImport(L *lua.LState) int {
	file := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"stack", "import", "--file", file}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaState(L *lua.LState) int {
	action := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"state", action}
	
	switch action {
	case "delete":
		if urn := options.RawGetString("urn").String(); urn != "" {
			args = append(args, urn)
		}
	case "unprotect":
		if urn := options.RawGetString("urn").String(); urn != "" {
			args = append(args, urn)
		}
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

// Plugin management
func (p *EnhancedPulumiModule) luaPluginInstall(L *lua.LState) int {
	kind := L.CheckString(1)
	name := L.CheckString(2)
	options := L.OptTable(3, L.NewTable())
	
	args := []string{"plugin", "install", kind, name}
	
	if version := options.RawGetString("version").String(); version != "" {
		args = append(args, version)
	}
	
	if reinstall := lua.LVAsBool(options.RawGetString("reinstall")); reinstall {
		args = append(args, "--reinstall")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaPluginList(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"plugin", "ls", "--json"}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaPluginRemove(L *lua.LState) int {
	kind := L.CheckString(1)
	name := L.CheckString(2)
	options := L.OptTable(3, L.NewTable())
	
	args := []string{"plugin", "rm", kind, name}
	
	if version := options.RawGetString("version").String(); version != "" {
		args = append(args, version)
	}
	
	if yes := lua.LVAsBool(options.RawGetString("yes")); yes {
		args = append(args, "--yes")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

// Advanced features
func (p *EnhancedPulumiModule) luaWatch(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"watch"}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaLogs(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"logs"}
	
	if follow := lua.LVAsBool(options.RawGetString("follow")); follow {
		args = append(args, "--follow")
	}
	
	if since := options.RawGetString("since").String(); since != "" {
		args = append(args, "--since", since)
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaHistory(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"stack", "history", "--json"}
	
	if showSecrets := lua.LVAsBool(options.RawGetString("show_secrets")); showSecrets {
		args = append(args, "--show-secrets")
	}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaConsole(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"console"}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

// Utility functions
func (p *EnhancedPulumiModule) luaVersion(L *lua.LState) int {
	result := p.executePulumiCmd([]string{"version"}, L.NewTable())
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaWhoami(L *lua.LState) int {
	result := p.executePulumiCmd([]string{"whoami"}, L.NewTable())
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaLogin(L *lua.LState) int {
	url := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"login", url}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

func (p *EnhancedPulumiModule) luaLogout(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"logout"}
	
	result := p.executePulumiCmd(args, options)
	L.Push(p.resultToLua(L, result))
	return 1
}

// Helper functions
func (p *EnhancedPulumiModule) executePulumiCmd(args []string, options *lua.LTable) *PulumiResult {
	startTime := time.Now()
	
	// Setup command with environment
	var commands []string
	
	// Handle virtual environment activation
	if venvPath := options.RawGetString("venv_path").String(); venvPath != "" {
		activateScript := filepath.Join(venvPath, "bin", "activate")
		commands = append(commands, fmt.Sprintf("source %s", activateScript))
	}
	
	// Handle login
	if loginURL := options.RawGetString("login_url").String(); loginURL != "" {
		commands = append(commands, fmt.Sprintf("pulumi login %s", loginURL))
	}
	
	// Add main command
	pulumiCmd := "pulumi " + strings.Join(args, " ")
	commands = append(commands, pulumiCmd)
	
	fullCommand := strings.Join(commands, " && ")
	cmd := exec.Command("bash", "-c", fullCommand)
	
	// Set working directory
	if workdir := options.RawGetString("workdir").String(); workdir != "" {
		cmd.Dir = workdir
	}
	
	// Setup environment
	env := os.Environ()
	
	// Add Pulumi to PATH
	if homeDir, err := os.UserHomeDir(); err == nil {
		pulumiPath := filepath.Join(homeDir, ".pulumi", "bin")
		newPath := fmt.Sprintf("PATH=%s:%s", pulumiPath, os.Getenv("PATH"))
		
		found := false
		for i, v := range env {
			if strings.HasPrefix(v, "PATH=") {
				env[i] = newPath
				found = true
				break
			}
		}
		if !found {
			env = append(env, newPath)
		}
	}
	
	// Add custom environment variables
	if envTable := options.RawGetString("env"); envTable.Type() == lua.LTTable {
		envTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			env = append(env, fmt.Sprintf("%s=%s", lua.LVAsString(key), lua.LVAsString(value)))
		})
	}
	
	cmd.Env = env
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	duration := time.Since(startTime)
	
	result := &PulumiResult{
		Success:  err == nil,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}
	
	// Try to parse JSON output for structured data
	if result.Success && result.Stdout != "" {
		var data map[string]interface{}
		if json.Unmarshal([]byte(result.Stdout), &data) == nil {
			if summary, ok := data["summary"]; ok {
				if summaryMap, ok := summary.(map[string]interface{}); ok {
					result.Summary = summaryMap
				}
			}
			if resources, ok := data["resources"]; ok {
				if resourcesList, ok := resources.([]interface{}); ok {
					for _, res := range resourcesList {
						if resMap, ok := res.(map[string]interface{}); ok {
							result.Resources = append(result.Resources, resMap)
						}
					}
				}
			}
			if outputs, ok := data["outputs"]; ok {
				if outputsMap, ok := outputs.(map[string]interface{}); ok {
					result.Outputs = outputsMap
				}
			}
			if permalink, ok := data["permalink"]; ok {
				if permalinkStr, ok := permalink.(string); ok {
					result.Permalink = permalinkStr
				}
			}
		}
	}
	
	return result
}

func (p *EnhancedPulumiModule) stackToLua(L *lua.LState, stack *PulumiStack) lua.LValue {
	stackTable := L.NewTable()
	L.SetField(stackTable, "name", lua.LString(stack.Name))
	L.SetField(stackTable, "project", lua.LString(stack.Project))
	L.SetField(stackTable, "workdir", lua.LString(stack.WorkDir))
	L.SetField(stackTable, "login_url", lua.LString(stack.LoginURL))
	L.SetField(stackTable, "config_file", lua.LString(stack.ConfigFile))
	L.SetField(stackTable, "backend", lua.LString(stack.Backend))
	
	// Environment variables
	envTable := L.NewTable()
	for key, value := range stack.Environment {
		envTable.RawSetString(key, lua.LString(value))
	}
	L.SetField(stackTable, "env", envTable)
	
	// Tags
	tagsTable := L.NewTable()
	for key, value := range stack.Tags {
		tagsTable.RawSetString(key, lua.LString(value))
	}
	L.SetField(stackTable, "tags", tagsTable)
	
	return stackTable
}

func (p *EnhancedPulumiModule) resultToLua(L *lua.LState, result *PulumiResult) lua.LValue {
	resultTable := L.NewTable()
	L.SetField(resultTable, "success", lua.LBool(result.Success))
	L.SetField(resultTable, "exit_code", lua.LNumber(result.ExitCode))
	L.SetField(resultTable, "stdout", lua.LString(result.Stdout))
	L.SetField(resultTable, "stderr", lua.LString(result.Stderr))
	L.SetField(resultTable, "duration_ms", lua.LNumber(result.Duration.Milliseconds()))
	L.SetField(resultTable, "permalink", lua.LString(result.Permalink))
	
	if result.Summary != nil {
		L.SetField(resultTable, "summary", p.mapToLuaTable(L, result.Summary))
	}
	
	if result.Outputs != nil {
		L.SetField(resultTable, "outputs", p.mapToLuaTable(L, result.Outputs))
	}
	
	if len(result.Resources) > 0 {
		resourcesTable := L.NewTable()
		for i, resource := range result.Resources {
			resourcesTable.RawSetInt(i+1, p.mapToLuaTable(L, resource))
		}
		L.SetField(resultTable, "resources", resourcesTable)
	}
	
	return resultTable
}

func (p *EnhancedPulumiModule) mapToLuaTable(L *lua.LState, data map[string]interface{}) *lua.LTable {
	table := L.NewTable()
	for key, value := range data {
		switch v := value.(type) {
		case string:
			table.RawSetString(key, lua.LString(v))
		case float64:
			table.RawSetString(key, lua.LNumber(v))
		case bool:
			table.RawSetString(key, lua.LBool(v))
		case map[string]interface{}:
			table.RawSetString(key, p.mapToLuaTable(L, v))
		case []interface{}:
			arrayTable := L.NewTable()
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					arrayTable.RawSetInt(i+1, p.mapToLuaTable(L, itemMap))
				} else {
					arrayTable.RawSetInt(i+1, lua.LString(fmt.Sprintf("%v", item)))
				}
			}
			table.RawSetString(key, arrayTable)
		default:
			table.RawSetString(key, lua.LString(fmt.Sprintf("%v", v)))
		}
	}
	return table
}