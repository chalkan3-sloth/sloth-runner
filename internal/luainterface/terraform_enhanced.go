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

	lua "github.com/yuin/gopher-lua"
)

// EnhancedTerraformModule provides advanced Terraform functionality
type EnhancedTerraformModule struct{}

// TerraformWorkspace represents a Terraform workspace
type TerraformWorkspace struct {
	Name      string
	WorkDir   string
	VarFiles  []string
	Variables map[string]interface{}
	Backend   map[string]interface{}
	Providers map[string]interface{}
	Environment map[string]string
	StateFile string
	Parallelism int
	Timeout   time.Duration
}

// TerraformResult represents Terraform operation result
type TerraformResult struct {
	Success    bool
	ExitCode   int
	Stdout     string
	Stderr     string
	Duration   time.Duration
	Plan       *TerraformPlan
	State      *TerraformState
	Outputs    map[string]interface{}
	Resources  []map[string]interface{}
	Changes    *TerraformChanges
}

// TerraformPlan represents a Terraform plan
type TerraformPlan struct {
	FormatVersion string                 `json:"format_version"`
	Variables     map[string]interface{} `json:"variables"`
	PlannedValues map[string]interface{} `json:"planned_values"`
	Changes       []map[string]interface{} `json:"resource_changes"`
	Summary       map[string]int         `json:"summary"`
}

// TerraformState represents Terraform state
type TerraformState struct {
	Version   int                    `json:"version"`
	Resources []map[string]interface{} `json:"resources"`
	Outputs   map[string]interface{} `json:"outputs"`
}

// TerraformChanges represents planned changes
type TerraformChanges struct {
	Add    int
	Change int
	Destroy int
}

// NewEnhancedTerraformModule creates a new enhanced Terraform module
func NewEnhancedTerraformModule() *EnhancedTerraformModule {
	return &EnhancedTerraformModule{}
}

// RegisterEnhancedTerraformModule registers the enhanced Terraform module
func RegisterEnhancedTerraformModule(L *lua.LState) {
	module := NewEnhancedTerraformModule()
	
	// Create the terraform table
	terraformTable := L.NewTable()
	
	// Workspace management
	L.SetField(terraformTable, "workspace", L.NewFunction(module.luaNewWorkspace))
	L.SetField(terraformTable, "workspace_list", L.NewFunction(module.luaWorkspaceList))
	L.SetField(terraformTable, "workspace_select", L.NewFunction(module.luaWorkspaceSelect))
	L.SetField(terraformTable, "workspace_new", L.NewFunction(module.luaWorkspaceNew))
	L.SetField(terraformTable, "workspace_delete", L.NewFunction(module.luaWorkspaceDelete))
	
	// Basic operations
	L.SetField(terraformTable, "init", L.NewFunction(module.luaInit))
	L.SetField(terraformTable, "plan", L.NewFunction(module.luaPlan))
	L.SetField(terraformTable, "apply", L.NewFunction(module.luaApply))
	L.SetField(terraformTable, "destroy", L.NewFunction(module.luaDestroy))
	L.SetField(terraformTable, "refresh", L.NewFunction(module.luaRefresh))
	
	// State management
	L.SetField(terraformTable, "state_list", L.NewFunction(module.luaStateList))
	L.SetField(terraformTable, "state_show", L.NewFunction(module.luaStateShow))
	L.SetField(terraformTable, "state_pull", L.NewFunction(module.luaStatePull))
	L.SetField(terraformTable, "state_push", L.NewFunction(module.luaStatePush))
	L.SetField(terraformTable, "state_mv", L.NewFunction(module.luaStateMove))
	L.SetField(terraformTable, "state_rm", L.NewFunction(module.luaStateRemove))
	
	// Import and taint
	L.SetField(terraformTable, "import", L.NewFunction(module.luaImport))
	L.SetField(terraformTable, "taint", L.NewFunction(module.luaTaint))
	L.SetField(terraformTable, "untaint", L.NewFunction(module.luaUntaint))
	
	// Output and validation
	L.SetField(terraformTable, "output", L.NewFunction(module.luaOutput))
	L.SetField(terraformTable, "validate", L.NewFunction(module.luaValidate))
	L.SetField(terraformTable, "fmt", L.NewFunction(module.luaFmt))
	
	// Advanced operations
	L.SetField(terraformTable, "graph", L.NewFunction(module.luaGraph))
	L.SetField(terraformTable, "force_unlock", L.NewFunction(module.luaForceUnlock))
	L.SetField(terraformTable, "providers", L.NewFunction(module.luaProviders))
	
	// Utility functions
	L.SetField(terraformTable, "version", L.NewFunction(module.luaVersion))
	L.SetField(terraformTable, "env", L.NewFunction(module.luaEnv))
	L.SetField(terraformTable, "console", L.NewFunction(module.luaConsole))
	
	// Register the terraform table globally
	L.SetGlobal("terraform", terraformTable)
}

// Workspace management
func (t *EnhancedTerraformModule) luaNewWorkspace(L *lua.LState) int {
	options := L.CheckTable(1)
	
	workspace := &TerraformWorkspace{
		Name:        options.RawGetString("name").String(),
		WorkDir:     options.RawGetString("workdir").String(),
		Variables:   make(map[string]interface{}),
		Backend:     make(map[string]interface{}),
		Providers:   make(map[string]interface{}),
		Environment: make(map[string]string),
		Parallelism: 10,
		Timeout:     30 * time.Minute,
	}
	
	// Parse var files
	if varFilesTable := options.RawGetString("var_files"); varFilesTable.Type() == lua.LTTable {
		varFilesTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
			workspace.VarFiles = append(workspace.VarFiles, lua.LVAsString(value))
		})
	}
	
	// Parse variables
	if varsTable := options.RawGetString("variables"); varsTable.Type() == lua.LTTable {
		workspace.Variables = t.luaTableToMap(varsTable.(*lua.LTable))
	}
	
	// Parse backend config
	if backendTable := options.RawGetString("backend"); backendTable.Type() == lua.LTTable {
		workspace.Backend = t.luaTableToMap(backendTable.(*lua.LTable))
	}
	
	// Parse providers
	if providersTable := options.RawGetString("providers"); providersTable.Type() == lua.LTTable {
		workspace.Providers = t.luaTableToMap(providersTable.(*lua.LTable))
	}
	
	// Parse environment variables
	if envTable := options.RawGetString("env"); envTable.Type() == lua.LTTable {
		envTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			workspace.Environment[lua.LVAsString(key)] = lua.LVAsString(value)
		})
	}
	
	// Parse options
	if parallelism := int(options.RawGetString("parallelism").(lua.LNumber)); parallelism > 0 {
		workspace.Parallelism = parallelism
	}
	
	if timeout := options.RawGetString("timeout").String(); timeout != "" {
		if duration, err := time.ParseDuration(timeout); err == nil {
			workspace.Timeout = duration
		}
	}
	
	L.Push(t.workspaceToLua(L, workspace))
	return 1
}

func (t *EnhancedTerraformModule) luaWorkspaceList(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	result := t.executeTerraformCmd([]string{"workspace", "list"}, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaWorkspaceSelect(L *lua.LState) int {
	name := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"workspace", "select", name}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaWorkspaceNew(L *lua.LState) int {
	name := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"workspace", "new", name}
	
	if state := options.RawGetString("state").String(); state != "" {
		args = append(args, "-state="+state)
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaWorkspaceDelete(L *lua.LState) int {
	name := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"workspace", "delete", name}
	
	if force := lua.LVAsBool(options.RawGetString("force")); force {
		args = append(args, "-force")
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

// Basic operations
func (t *EnhancedTerraformModule) luaInit(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"init", "-input=false", "-no-color"}
	
	if upgrade := lua.LVAsBool(options.RawGetString("upgrade")); upgrade {
		args = append(args, "-upgrade")
	}
	
	if reconfigure := lua.LVAsBool(options.RawGetString("reconfigure")); reconfigure {
		args = append(args, "-reconfigure")
	}
	
	if migrateState := lua.LVAsBool(options.RawGetString("migrate_state")); migrateState {
		args = append(args, "-migrate-state")
	}
	
	if backend := lua.LVAsBool(options.RawGetString("backend")); !backend {
		args = append(args, "-backend=false")
	}
	
	// Add backend config
	if backendConfigTable := options.RawGetString("backend_config"); backendConfigTable.Type() == lua.LTTable {
		backendConfigTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			args = append(args, fmt.Sprintf("-backend-config=%s=%s", lua.LVAsString(key), lua.LVAsString(value)))
		})
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaPlan(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"plan", "-input=false", "-no-color"}
	
	if out := options.RawGetString("out").String(); out != "" {
		args = append(args, "-out="+out)
	}
	
	if destroy := lua.LVAsBool(options.RawGetString("destroy")); destroy {
		args = append(args, "-destroy")
	}
	
	if refresh := lua.LVAsBool(options.RawGetString("refresh")); !refresh {
		args = append(args, "-refresh=false")
	}
	
	if detailed := lua.LVAsBool(options.RawGetString("detailed_exitcode")); detailed {
		args = append(args, "-detailed-exitcode")
	}
	
	if parallelism := int(options.RawGetString("parallelism").(lua.LNumber)); parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", parallelism))
	}
	
	// Add target resources
	if targetsTable := options.RawGetString("targets"); targetsTable.Type() == lua.LTTable {
		targetsTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
			args = append(args, "-target="+lua.LVAsString(value))
		})
	}
	
	// Add variables
	args = t.addVariables(args, options)
	
	result := t.executeTerraformCmd(args, options)
	
	// Parse plan output if JSON
	if jsonOutput := lua.LVAsBool(options.RawGetString("json")); jsonOutput {
		t.parsePlanOutput(result)
	}
	
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaApply(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"apply", "-input=false", "-no-color"}
	
	if autoApprove := lua.LVAsBool(options.RawGetString("auto_approve")); autoApprove {
		args = append(args, "-auto-approve")
	}
	
	if plan := options.RawGetString("plan").String(); plan != "" {
		args = append(args, plan)
	} else {
		// Add variables if no plan file
		args = t.addVariables(args, options)
	}
	
	if parallelism := int(options.RawGetString("parallelism").(lua.LNumber)); parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", parallelism))
	}
	
	// Add target resources
	if targetsTable := options.RawGetString("targets"); targetsTable.Type() == lua.LTTable {
		targetsTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
			args = append(args, "-target="+lua.LVAsString(value))
		})
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaDestroy(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"destroy", "-input=false", "-no-color"}
	
	if autoApprove := lua.LVAsBool(options.RawGetString("auto_approve")); autoApprove {
		args = append(args, "-auto-approve")
	}
	
	if parallelism := int(options.RawGetString("parallelism").(lua.LNumber)); parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", parallelism))
	}
	
	// Add target resources
	if targetsTable := options.RawGetString("targets"); targetsTable.Type() == lua.LTTable {
		targetsTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
			args = append(args, "-target="+lua.LVAsString(value))
		})
	}
	
	// Add variables
	args = t.addVariables(args, options)
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaRefresh(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"refresh", "-input=false", "-no-color"}
	
	// Add variables
	args = t.addVariables(args, options)
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

// State management
func (t *EnhancedTerraformModule) luaStateList(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"state", "list"}
	
	if id := options.RawGetString("id").String(); id != "" {
		args = append(args, "-id="+id)
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaStateShow(L *lua.LState) int {
	address := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"state", "show", address}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaStatePull(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"state", "pull"}
	
	result := t.executeTerraformCmd(args, options)
	
	// Try to parse state as JSON
	if result.Success && result.Stdout != "" {
		var state TerraformState
		if json.Unmarshal([]byte(result.Stdout), &state) == nil {
			result.State = &state
		}
	}
	
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaStatePush(L *lua.LState) int {
	file := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"state", "push", file}
	
	if force := lua.LVAsBool(options.RawGetString("force")); force {
		args = append(args, "-force")
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaStateMove(L *lua.LState) int {
	source := L.CheckString(1)
	destination := L.CheckString(2)
	options := L.OptTable(3, L.NewTable())
	
	args := []string{"state", "mv", source, destination}
	
	if dryRun := lua.LVAsBool(options.RawGetString("dry_run")); dryRun {
		args = append(args, "-dry-run")
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaStateRemove(L *lua.LState) int {
	address := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"state", "rm", address}
	
	if dryRun := lua.LVAsBool(options.RawGetString("dry_run")); dryRun {
		args = append(args, "-dry-run")
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

// Import and taint
func (t *EnhancedTerraformModule) luaImport(L *lua.LState) int {
	address := L.CheckString(1)
	id := L.CheckString(2)
	options := L.OptTable(3, L.NewTable())
	
	args := []string{"import", address, id}
	
	// Add variables
	args = t.addVariables(args, options)
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaTaint(L *lua.LState) int {
	address := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"taint", address}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaUntaint(L *lua.LState) int {
	address := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"untaint", address}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

// Output and validation
func (t *EnhancedTerraformModule) luaOutput(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"output", "-json"}
	
	if name := options.RawGetString("name").String(); name != "" {
		args = append(args, name)
	}
	
	if raw := lua.LVAsBool(options.RawGetString("raw")); raw {
		args = append(args, "-raw")
	}
	
	result := t.executeTerraformCmd(args, options)
	
	// Parse outputs
	if result.Success && result.Stdout != "" {
		var outputs map[string]interface{}
		if json.Unmarshal([]byte(result.Stdout), &outputs) == nil {
			result.Outputs = outputs
		}
	}
	
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaValidate(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"validate", "-json"}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaFmt(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"fmt"}
	
	if recursive := lua.LVAsBool(options.RawGetString("recursive")); recursive {
		args = append(args, "-recursive")
	}
	
	if diff := lua.LVAsBool(options.RawGetString("diff")); diff {
		args = append(args, "-diff")
	}
	
	if check := lua.LVAsBool(options.RawGetString("check")); check {
		args = append(args, "-check")
	}
	
	if write := lua.LVAsBool(options.RawGetString("write")); !write {
		args = append(args, "-write=false")
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

// Advanced operations
func (t *EnhancedTerraformModule) luaGraph(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"graph"}
	
	if graphType := options.RawGetString("type").String(); graphType != "" {
		args = append(args, "-type="+graphType)
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaForceUnlock(L *lua.LState) int {
	lockId := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	args := []string{"force-unlock", lockId}
	
	if force := lua.LVAsBool(options.RawGetString("force")); force {
		args = append(args, "-force")
	}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaProviders(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"providers"}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

// Utility functions
func (t *EnhancedTerraformModule) luaVersion(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"version", "-json"}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaEnv(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"env", "list"}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

func (t *EnhancedTerraformModule) luaConsole(L *lua.LState) int {
	options := L.OptTable(1, L.NewTable())
	
	args := []string{"console"}
	
	result := t.executeTerraformCmd(args, options)
	L.Push(t.resultToLua(L, result))
	return 1
}

// Helper functions
func (t *EnhancedTerraformModule) executeTerraformCmd(args []string, options *lua.LTable) *TerraformResult {
	startTime := time.Now()
	
	cmd := exec.Command("terraform", args...)
	
	// Set working directory
	if workdir := options.RawGetString("workdir").String(); workdir != "" {
		cmd.Dir = workdir
	}
	
	// Setup environment
	env := os.Environ()
	
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
	
	result := &TerraformResult{
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
	
	return result
}

func (t *EnhancedTerraformModule) addVariables(args []string, options *lua.LTable) []string {
	// Add variable files
	if varFilesTable := options.RawGetString("var_files"); varFilesTable.Type() == lua.LTTable {
		varFilesTable.(*lua.LTable).ForEach(func(_, value lua.LValue) {
			args = append(args, "-var-file="+lua.LVAsString(value))
		})
	}
	
	// Add variables
	if varsTable := options.RawGetString("variables"); varsTable.Type() == lua.LTTable {
		varsTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			args = append(args, fmt.Sprintf("-var=%s=%s", lua.LVAsString(key), lua.LVAsString(value)))
		})
	}
	
	return args
}

func (t *EnhancedTerraformModule) parsePlanOutput(result *TerraformResult) {
	if result.Stdout != "" {
		var plan TerraformPlan
		if json.Unmarshal([]byte(result.Stdout), &plan) == nil {
			result.Plan = &plan
			
			// Calculate changes summary
			changes := &TerraformChanges{}
			for _, change := range plan.Changes {
				if actions, ok := change["change"].(map[string]interface{}); ok {
					if actionsList, ok := actions["actions"].([]interface{}); ok {
						for _, action := range actionsList {
							switch action.(string) {
							case "create":
								changes.Add++
							case "update":
								changes.Change++
							case "delete":
								changes.Destroy++
							}
						}
					}
				}
			}
			result.Changes = changes
		}
	}
}

func (t *EnhancedTerraformModule) workspaceToLua(L *lua.LState, workspace *TerraformWorkspace) lua.LValue {
	workspaceTable := L.NewTable()
	L.SetField(workspaceTable, "name", lua.LString(workspace.Name))
	L.SetField(workspaceTable, "workdir", lua.LString(workspace.WorkDir))
	L.SetField(workspaceTable, "state_file", lua.LString(workspace.StateFile))
	L.SetField(workspaceTable, "parallelism", lua.LNumber(workspace.Parallelism))
	L.SetField(workspaceTable, "timeout", lua.LString(workspace.Timeout.String()))
	
	// Var files
	varFilesTable := L.NewTable()
	for i, file := range workspace.VarFiles {
		varFilesTable.RawSetInt(i+1, lua.LString(file))
	}
	L.SetField(workspaceTable, "var_files", varFilesTable)
	
	// Variables
	varsTable := L.NewTable()
	for key, value := range workspace.Variables {
		varsTable.RawSetString(key, lua.LString(fmt.Sprintf("%v", value)))
	}
	L.SetField(workspaceTable, "variables", varsTable)
	
	// Environment
	envTable := L.NewTable()
	for key, value := range workspace.Environment {
		envTable.RawSetString(key, lua.LString(value))
	}
	L.SetField(workspaceTable, "env", envTable)
	
	return workspaceTable
}

func (t *EnhancedTerraformModule) resultToLua(L *lua.LState, result *TerraformResult) lua.LValue {
	resultTable := L.NewTable()
	L.SetField(resultTable, "success", lua.LBool(result.Success))
	L.SetField(resultTable, "exit_code", lua.LNumber(result.ExitCode))
	L.SetField(resultTable, "stdout", lua.LString(result.Stdout))
	L.SetField(resultTable, "stderr", lua.LString(result.Stderr))
	L.SetField(resultTable, "duration_ms", lua.LNumber(result.Duration.Milliseconds()))
	
	if result.Plan != nil {
		planTable := L.NewTable()
		L.SetField(planTable, "format_version", lua.LString(result.Plan.FormatVersion))
		
		if result.Plan.Summary != nil {
			summaryTable := L.NewTable()
			for key, value := range result.Plan.Summary {
				summaryTable.RawSetString(key, lua.LNumber(value))
			}
			L.SetField(planTable, "summary", summaryTable)
		}
		
		L.SetField(resultTable, "plan", planTable)
	}
	
	if result.Changes != nil {
		changesTable := L.NewTable()
		L.SetField(changesTable, "add", lua.LNumber(result.Changes.Add))
		L.SetField(changesTable, "change", lua.LNumber(result.Changes.Change))
		L.SetField(changesTable, "destroy", lua.LNumber(result.Changes.Destroy))
		L.SetField(resultTable, "changes", changesTable)
	}
	
	if result.Outputs != nil {
		L.SetField(resultTable, "outputs", t.mapToLuaTable(L, result.Outputs))
	}
	
	if result.State != nil {
		stateTable := L.NewTable()
		L.SetField(stateTable, "version", lua.LNumber(result.State.Version))
		
		if result.State.Outputs != nil {
			L.SetField(stateTable, "outputs", t.mapToLuaTable(L, result.State.Outputs))
		}
		
		L.SetField(resultTable, "state", stateTable)
	}
	
	return resultTable
}

func (t *EnhancedTerraformModule) luaTableToMap(table *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	table.ForEach(func(key, value lua.LValue) {
		switch v := value.(type) {
		case lua.LString:
			result[lua.LVAsString(key)] = string(v)
		case lua.LNumber:
			result[lua.LVAsString(key)] = float64(v)
		case lua.LBool:
			result[lua.LVAsString(key)] = bool(v)
		case *lua.LTable:
			result[lua.LVAsString(key)] = t.luaTableToMap(v)
		default:
			result[lua.LVAsString(key)] = lua.LVAsString(value)
		}
	})
	return result
}

func (t *EnhancedTerraformModule) mapToLuaTable(L *lua.LState, data map[string]interface{}) *lua.LTable {
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
			table.RawSetString(key, t.mapToLuaTable(L, v))
		case []interface{}:
			arrayTable := L.NewTable()
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					arrayTable.RawSetInt(i+1, t.mapToLuaTable(L, itemMap))
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