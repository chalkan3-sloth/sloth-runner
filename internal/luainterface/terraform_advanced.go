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

// TerraformModule provides advanced Terraform integration
type TerraformModule struct{}

// NewTerraformModule creates a new TerraformModule
func NewTerraformModule() *TerraformModule {
	return &TerraformModule{}
}

// Loader returns the Lua loader for the terraform module
func (mod *TerraformModule) Loader(L *lua.LState) int {
	terraformTable := L.NewTable()
	L.SetFuncs(terraformTable, map[string]lua.LGFunction{
		"init":                mod.terraformInit,
		"plan":                mod.terraformPlan,
		"apply":               mod.terraformApply,
		"destroy":             mod.terraformDestroy,
		"validate":            mod.terraformValidate,
		"fmt":                 mod.terraformFmt,
		"output":              mod.terraformOutput,
		"show":                mod.terraformShow,
		"state_list":          mod.terraformStateList,
		"state_show":          mod.terraformStateShow,
		"state_mv":            mod.terraformStateMove,
		"state_rm":            mod.terraformStateRemove,
		"state_pull":          mod.terraformStatePull,
		"state_push":          mod.terraformStatePush,
		"import":              mod.terraformImport,
		"taint":               mod.terraformTaint,
		"untaint":             mod.terraformUntaint,
		"workspace_list":      mod.terraformWorkspaceList,
		"workspace_new":       mod.terraformWorkspaceNew,
		"workspace_select":    mod.terraformWorkspaceSelect,
		"workspace_delete":    mod.terraformWorkspaceDelete,
		"providers":           mod.terraformProviders,
		"providers_lock":      mod.terraformProvidersLock,
		"providers_mirror":    mod.terraformProvidersMirror,
		"refresh":             mod.terraformRefresh,
		"graph":               mod.terraformGraph,
		"version":             mod.terraformVersion,
		"force_unlock":        mod.terraformForceUnlock,
		"get":                 mod.terraformGet,
		"console":             mod.terraformConsole,
	})
	L.Push(terraformTable)
	return 1
}

// terraformInit initializes a Terraform working directory
func (mod *TerraformModule) terraformInit(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"init"}
	
	// Add options
	if upgrade := opts.RawGetString("upgrade"); lua.LVAsBool(upgrade) {
		args = append(args, "-upgrade")
	}
	
	if reconfigure := opts.RawGetString("reconfigure"); lua.LVAsBool(reconfigure) {
		args = append(args, "-reconfigure")
	}
	
	if migrate := opts.RawGetString("migrate_state"); lua.LVAsBool(migrate) {
		args = append(args, "-migrate-state")
	}
	
	if backend := opts.RawGetString("backend"); !lua.LVAsBool(backend) {
		args = append(args, "-backend=false")
	}
	
	if backendConfig := opts.RawGetString("backend_config"); backendConfig != lua.LNil {
		args = append(args, "-backend-config="+backendConfig.String())
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformPlan creates an execution plan
func (mod *TerraformModule) terraformPlan(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"plan"}
	
	// Add options
	if out := opts.RawGetString("out"); out != lua.LNil {
		args = append(args, "-out="+out.String())
	}
	
	if destroy := opts.RawGetString("destroy"); lua.LVAsBool(destroy) {
		args = append(args, "-destroy")
	}
	
	if refresh := opts.RawGetString("refresh"); !lua.LVAsBool(refresh) {
		args = append(args, "-refresh=false")
	}
	
	if varFile := opts.RawGetString("var_file"); varFile != lua.LNil {
		args = append(args, "-var-file="+varFile.String())
	}
	
	if vars := opts.RawGetString("vars"); vars != lua.LNil {
		// Parse variables from table or string
		if varsTable, ok := vars.(*lua.LTable); ok {
			varsTable.ForEach(func(key, value lua.LValue) {
				args = append(args, fmt.Sprintf("-var=%s=%s", key.String(), value.String()))
			})
		}
	}
	
	if detailed := opts.RawGetString("detailed_exitcode"); lua.LVAsBool(detailed) {
		args = append(args, "-detailed-exitcode")
	}
	
	if jsonOutput := opts.RawGetString("json"); lua.LVAsBool(jsonOutput) {
		args = append(args, "-json")
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformApply applies the changes required to reach the desired state
func (mod *TerraformModule) terraformApply(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"apply"}
	
	// Add options
	if planFile := opts.RawGetString("plan"); planFile != lua.LNil {
		args = append(args, planFile.String())
	} else {
		if autoApprove := opts.RawGetString("auto_approve"); lua.LVAsBool(autoApprove) {
			args = append(args, "-auto-approve")
		}
	}
	
	if varFile := opts.RawGetString("var_file"); varFile != lua.LNil {
		args = append(args, "-var-file="+varFile.String())
	}
	
	if vars := opts.RawGetString("vars"); vars != lua.LNil {
		if varsTable, ok := vars.(*lua.LTable); ok {
			varsTable.ForEach(func(key, value lua.LValue) {
				args = append(args, fmt.Sprintf("-var=%s=%s", key.String(), value.String()))
			})
		}
	}
	
	if parallelism := opts.RawGetString("parallelism"); parallelism != lua.LNil {
		args = append(args, "-parallelism="+parallelism.String())
	}
	
	if jsonOutput := opts.RawGetString("json"); lua.LVAsBool(jsonOutput) {
		args = append(args, "-json")
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformDestroy destroys Terraform-managed infrastructure
func (mod *TerraformModule) terraformDestroy(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"destroy"}
	
	if autoApprove := opts.RawGetString("auto_approve"); lua.LVAsBool(autoApprove) {
		args = append(args, "-auto-approve")
	}
	
	if varFile := opts.RawGetString("var_file"); varFile != lua.LNil {
		args = append(args, "-var-file="+varFile.String())
	}
	
	if vars := opts.RawGetString("vars"); vars != lua.LNil {
		if varsTable, ok := vars.(*lua.LTable); ok {
			varsTable.ForEach(func(key, value lua.LValue) {
				args = append(args, fmt.Sprintf("-var=%s=%s", key.String(), value.String()))
			})
		}
	}
	
	if target := opts.RawGetString("target"); target != lua.LNil {
		args = append(args, "-target="+target.String())
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformValidate validates the Terraform files
func (mod *TerraformModule) terraformValidate(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"validate"}
	
	if jsonOutput := opts.RawGetString("json"); lua.LVAsBool(jsonOutput) {
		args = append(args, "-json")
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformFmt formats Terraform configuration files
func (mod *TerraformModule) terraformFmt(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"fmt"}
	
	if recursive := opts.RawGetString("recursive"); lua.LVAsBool(recursive) {
		args = append(args, "-recursive")
	}
	
	if check := opts.RawGetString("check"); lua.LVAsBool(check) {
		args = append(args, "-check")
	}
	
	if diff := opts.RawGetString("diff"); lua.LVAsBool(diff) {
		args = append(args, "-diff")
	}
	
	if write := opts.RawGetString("write"); !lua.LVAsBool(write) {
		args = append(args, "-write=false")
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformOutput reads an output variable
func (mod *TerraformModule) terraformOutput(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"output"}
	
	if jsonOutput := opts.RawGetString("json"); lua.LVAsBool(jsonOutput) {
		args = append(args, "-json")
	}
	
	if raw := opts.RawGetString("raw"); lua.LVAsBool(raw) {
		args = append(args, "-raw")
	}
	
	if outputName := opts.RawGetString("name"); outputName != lua.LNil {
		args = append(args, outputName.String())
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(strings.TrimSpace(result)))
	L.Push(lua.LNil)
	return 2
}

// terraformShow provides human-readable output from a state or plan file
func (mod *TerraformModule) terraformShow(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"show"}
	
	if jsonOutput := opts.RawGetString("json"); lua.LVAsBool(jsonOutput) {
		args = append(args, "-json")
	}
	
	if file := opts.RawGetString("file"); file != lua.LNil {
		args = append(args, file.String())
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformStateList lists resources in the state
func (mod *TerraformModule) terraformStateList(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"state", "list"}
	
	if address := opts.RawGetString("address"); address != lua.LNil {
		args = append(args, address.String())
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformStateShow shows the attributes of a resource in the state
func (mod *TerraformModule) terraformStateShow(L *lua.LState) int {
	address := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "state", "show", address)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformStateMove moves resources in the state
func (mod *TerraformModule) terraformStateMove(L *lua.LState) int {
	source := L.CheckString(1)
	destination := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "state", "mv", source, destination)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformStateRemove removes instances from the state
func (mod *TerraformModule) terraformStateRemove(L *lua.LState) int {
	address := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "state", "rm", address)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformStatePull downloads and outputs the state from remote state
func (mod *TerraformModule) terraformStatePull(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "state", "pull")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformStatePush uploads a local state file to remote state
func (mod *TerraformModule) terraformStatePush(L *lua.LState) int {
	statePath := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"state", "push"}
	
	if force := opts.RawGetString("force"); lua.LVAsBool(force) {
		args = append(args, "-force")
	}
	
	args = append(args, statePath)
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformImport imports existing infrastructure into Terraform state
func (mod *TerraformModule) terraformImport(L *lua.LState) int {
	address := L.CheckString(1)
	id := L.CheckString(2)
	
	opts := L.OptTable(3, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"import"}
	
	if varFile := opts.RawGetString("var_file"); varFile != lua.LNil {
		args = append(args, "-var-file="+varFile.String())
	}
	
	if vars := opts.RawGetString("vars"); vars != lua.LNil {
		if varsTable, ok := vars.(*lua.LTable); ok {
			varsTable.ForEach(func(key, value lua.LValue) {
				args = append(args, fmt.Sprintf("-var=%s=%s", key.String(), value.String()))
			})
		}
	}
	
	args = append(args, address, id)
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformTaint marks a resource for recreation
func (mod *TerraformModule) terraformTaint(L *lua.LState) int {
	address := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "taint", address)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformUntaint removes the taint from a resource
func (mod *TerraformModule) terraformUntaint(L *lua.LState) int {
	address := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "untaint", address)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformWorkspaceList lists workspaces
func (mod *TerraformModule) terraformWorkspaceList(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "workspace", "list")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformWorkspaceNew creates a new workspace
func (mod *TerraformModule) terraformWorkspaceNew(L *lua.LState) int {
	name := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "workspace", "new", name)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformWorkspaceSelect selects a workspace
func (mod *TerraformModule) terraformWorkspaceSelect(L *lua.LState) int {
	name := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "workspace", "select", name)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformWorkspaceDelete deletes a workspace
func (mod *TerraformModule) terraformWorkspaceDelete(L *lua.LState) int {
	name := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"workspace", "delete"}
	
	if force := opts.RawGetString("force"); lua.LVAsBool(force) {
		args = append(args, "-force")
	}
	
	args = append(args, name)
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformProviders prints information about providers
func (mod *TerraformModule) terraformProviders(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, "providers")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformProvidersLock writes provider dependency locks
func (mod *TerraformModule) terraformProvidersLock(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"providers", "lock"}
	
	if platform := opts.RawGetString("platform"); platform != lua.LNil {
		args = append(args, "-platform="+platform.String())
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformProvidersMirror saves local copies of providers
func (mod *TerraformModule) terraformProvidersMirror(L *lua.LState) int {
	targetDir := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"providers", "mirror"}
	
	if platform := opts.RawGetString("platform"); platform != lua.LNil {
		args = append(args, "-platform="+platform.String())
	}
	
	args = append(args, targetDir)
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformRefresh updates state to match remote systems
func (mod *TerraformModule) terraformRefresh(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"refresh"}
	
	if varFile := opts.RawGetString("var_file"); varFile != lua.LNil {
		args = append(args, "-var-file="+varFile.String())
	}
	
	if vars := opts.RawGetString("vars"); vars != lua.LNil {
		if varsTable, ok := vars.(*lua.LTable); ok {
			varsTable.ForEach(func(key, value lua.LValue) {
				args = append(args, fmt.Sprintf("-var=%s=%s", key.String(), value.String()))
			})
		}
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformGraph generates a dependency graph
func (mod *TerraformModule) terraformGraph(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"graph"}
	
	if graphType := opts.RawGetString("type"); graphType != lua.LNil {
		args = append(args, "-type="+graphType.String())
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformVersion displays version information
func (mod *TerraformModule) terraformVersion(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"version"}
	
	if jsonOutput := opts.RawGetString("json"); lua.LVAsBool(jsonOutput) {
		args = append(args, "-json")
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// terraformForceUnlock manually unlocks the state
func (mod *TerraformModule) terraformForceUnlock(L *lua.LState) int {
	lockID := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"force-unlock"}
	
	if force := opts.RawGetString("force"); lua.LVAsBool(force) {
		args = append(args, "-force")
	}
	
	args = append(args, lockID)
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformGet downloads and installs modules
func (mod *TerraformModule) terraformGet(L *lua.LState) int {
	opts := L.OptTable(1, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	args := []string{"get"}
	
	if update := opts.RawGetString("update"); lua.LVAsBool(update) {
		args = append(args, "-update")
	}
	
	result, err := mod.executeTerraformCommand(workdir, nil, args...)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(result))
	return 2
}

// terraformConsole starts an interactive console
func (mod *TerraformModule) terraformConsole(L *lua.LState) int {
	expression := L.CheckString(1)
	
	opts := L.OptTable(2, L.NewTable())
	workdir := opts.RawGetString("workdir").String()
	if workdir == "" {
		workdir = "."
	}
	
	// For console, we need to pass the expression via stdin
	cmd := exec.Command("terraform", "console")
	cmd.Dir = workdir
	cmd.Stdin = strings.NewReader(expression + "\n")
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("terraform console failed: %s", stderr.String())))
		return 2
	}
	
	L.Push(lua.LString(strings.TrimSpace(stdout.String())))
	L.Push(lua.LNil)
	return 2
}

// executeTerraformCommand executes a terraform command with environment variables
func (mod *TerraformModule) executeTerraformCommand(workdir string, env map[string]string, cmdArgs ...string) (string, error) {
	// Check if terraform command exists
	if _, err := exec.LookPath("terraform"); err != nil {
		return "", fmt.Errorf("terraform command not found in PATH: %w", err)
	}
	
	cmd := exec.Command("terraform", cmdArgs...)
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
			return "", fmt.Errorf("terraform command failed: %s", errorMsg)
		}
		return stdout.String(), nil
		
	case <-timer.C:
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("terraform command timed out after %v", timeout)
	}
}