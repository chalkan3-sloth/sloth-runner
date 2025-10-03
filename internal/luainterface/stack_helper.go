package luainterface

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	lua "github.com/yuin/gopher-lua"
)

// currentStack holds the active stack for the workflow
var currentStack *stack.StackState
var currentStackManager *stack.StackManager

// SetCurrentStack sets the active stack for the workflow
func SetCurrentStack(s *stack.StackState, sm *stack.StackManager) {
	currentStack = s
	currentStackManager = sm
}

// GetCurrentStack returns the active stack
func GetCurrentStack() *stack.StackState {
	return currentStack
}

// RegisterStackFunctions registers stack management functions in Lua
func RegisterStackFunctions(L *lua.LState) {
	// Global stack module
	stackModule := L.NewTable()

	// Stack information functions
	L.SetField(stackModule, "get_name", L.NewFunction(stackGetName))
	L.SetField(stackModule, "get_id", L.NewFunction(stackGetID))
	L.SetField(stackModule, "get_status", L.NewFunction(stackGetStatus))
	L.SetField(stackModule, "set_output", L.NewFunction(stackSetOutput))
	L.SetField(stackModule, "get_output", L.NewFunction(stackGetOutput))

	// Resource management functions
	L.SetField(stackModule, "register_resource", L.NewFunction(stackRegisterResource))
	L.SetField(stackModule, "get_resource", L.NewFunction(stackGetResource))
	L.SetField(stackModule, "update_resource", L.NewFunction(stackUpdateResource))
	L.SetField(stackModule, "delete_resource", L.NewFunction(stackDeleteResource))
	L.SetField(stackModule, "resource_exists", L.NewFunction(stackResourceExists))

	L.SetGlobal("stack", stackModule)
}

// stackGetName returns the current stack name
func stackGetName(L *lua.LState) int {
	if currentStack == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(currentStack.Name))
	return 1
}

// stackGetID returns the current stack ID
func stackGetID(L *lua.LState) int {
	if currentStack == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(currentStack.ID))
	return 1
}

// stackGetStatus returns the current stack status
func stackGetStatus(L *lua.LState) int {
	if currentStack == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(currentStack.Status))
	return 1
}

// stackSetOutput sets an output value in the stack
func stackSetOutput(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.Get(2)

	if currentStack == nil {
		L.RaiseError("no active stack")
		return 0
	}

	if currentStack.Outputs == nil {
		currentStack.Outputs = make(map[string]interface{})
	}

	currentStack.Outputs[key] = stackLuaValueToGo(value)

	// Update stack in database
	if currentStackManager != nil {
		if err := currentStackManager.UpdateStack(currentStack); err != nil {
			L.RaiseError("failed to update stack: %v", err)
			return 0
		}
	}

	return 0
}

// stackGetOutput gets an output value from the stack
func stackGetOutput(L *lua.LState) int {
	key := L.CheckString(1)

	if currentStack == nil || currentStack.Outputs == nil {
		L.Push(lua.LNil)
		return 1
	}

	value, ok := currentStack.Outputs[key]
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(stackGoValueToLua(L, value))
	return 1
}

// stackRegisterResource registers a resource in the stack
func stackRegisterResource(L *lua.LState) int {
	params := L.CheckTable(1)

	if currentStack == nil {
		L.RaiseError("no active stack")
		return 0
	}

	if currentStackManager == nil {
		L.RaiseError("no stack manager available")
		return 0
	}

	// Extract parameters
	resourceType := getTableString(params, "type", "")
	resourceName := getTableString(params, "name", "")
	moduleName := getTableString(params, "module", "")

	if resourceType == "" || resourceName == "" || moduleName == "" {
		L.RaiseError("type, name, and module are required")
		return 0
	}

	// Build resource ID
	resourceID := fmt.Sprintf("%s/%s/%s", currentStack.ID, moduleName, resourceType+"/"+resourceName)

	// Get properties
	properties := make(map[string]interface{})
	if propsLua := params.RawGetString("properties"); propsLua != lua.LNil {
		if propsTable, ok := propsLua.(*lua.LTable); ok {
			properties = stackLuaTableToMap(propsTable)
		}
	}

	// Calculate checksum
	propsJSON, _ := json.Marshal(properties)
	checksum := fmt.Sprintf("%x", sha256.Sum256(propsJSON))

	// Check if resource already exists
	existing, _ := currentStackManager.GetResourceByStackAndName(currentStack.ID, resourceType, resourceName)

	if existing != nil {
		// Resource exists - check if it needs update
		if existing.Checksum == checksum {
			// No changes - resource is idempotent
			L.Push(lua.LString("unchanged"))
			L.Push(stackGoValueToLua(L, existing))
			return 2
		}

		// Resource changed - update it
		existing.Properties = properties
		existing.Checksum = checksum
		existing.State = "pending"

		if err := currentStackManager.UpdateResource(existing); err != nil {
			L.RaiseError("failed to update resource: %v", err)
			return 0
		}

		L.Push(lua.LString("changed"))
		L.Push(stackGoValueToLua(L, existing))
		return 2
	}

	// Create new resource
	resource := &stack.Resource{
		ID:         resourceID,
		StackID:    currentStack.ID,
		Type:       resourceType,
		Name:       resourceName,
		Module:     moduleName,
		Properties: properties,
		State:      "pending",
		Checksum:   checksum,
		Metadata:   make(map[string]interface{}),
	}

	// Add dependencies if provided
	if depsLua := params.RawGetString("dependencies"); depsLua != lua.LNil {
		if depsTable, ok := depsLua.(*lua.LTable); ok {
			deps := []string{}
			depsTable.ForEach(func(_, v lua.LValue) {
				if str, ok := v.(lua.LString); ok {
					deps = append(deps, string(str))
				}
			})
			resource.Dependencies = deps
		}
	}

	if err := currentStackManager.CreateResource(resource); err != nil {
		L.RaiseError("failed to create resource: %v", err)
		return 0
	}

	L.Push(lua.LString("created"))
	L.Push(stackGoValueToLua(L, resource))
	return 2
}

// stackGetResource retrieves a resource from the stack
func stackGetResource(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)

	if currentStack == nil || currentStackManager == nil {
		L.Push(lua.LNil)
		return 1
	}

	resource, err := currentStackManager.GetResourceByStackAndName(currentStack.ID, resourceType, resourceName)
	if err != nil || resource == nil {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(stackGoValueToLua(L, resource))
	return 1
}

// stackUpdateResource updates a resource state
func stackUpdateResource(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)
	params := L.CheckTable(3)

	if currentStack == nil || currentStackManager == nil {
		L.RaiseError("no active stack or stack manager")
		return 0
	}

	resource, err := currentStackManager.GetResourceByStackAndName(currentStack.ID, resourceType, resourceName)
	if err != nil || resource == nil {
		L.RaiseError("resource not found: %s/%s", resourceType, resourceName)
		return 0
	}

	// Update state if provided
	if stateLua := params.RawGetString("state"); stateLua != lua.LNil {
		if stateStr, ok := stateLua.(lua.LString); ok {
			resource.State = string(stateStr)
		}
	}

	// Update error message if provided
	if errorLua := params.RawGetString("error"); errorLua != lua.LNil {
		if errorStr, ok := errorLua.(lua.LString); ok {
			resource.ErrorMessage = string(errorStr)
		}
	}

	// Update last applied if state is "applied"
	if resource.State == "applied" {
		now := time.Now()
		resource.LastApplied = &now
	}

	if err := currentStackManager.UpdateResource(resource); err != nil {
		L.RaiseError("failed to update resource: %v", err)
		return 0
	}

	L.Push(lua.LTrue)
	return 1
}

// stackDeleteResource deletes a resource from the stack
func stackDeleteResource(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)

	if currentStack == nil || currentStackManager == nil {
		L.RaiseError("no active stack or stack manager")
		return 0
	}

	resource, err := currentStackManager.GetResourceByStackAndName(currentStack.ID, resourceType, resourceName)
	if err != nil || resource == nil {
		L.RaiseError("resource not found: %s/%s", resourceType, resourceName)
		return 0
	}

	if err := currentStackManager.DeleteResource(resource.ID); err != nil {
		L.RaiseError("failed to delete resource: %v", err)
		return 0
	}

	L.Push(lua.LTrue)
	return 1
}

// stackResourceExists checks if a resource exists
func stackResourceExists(L *lua.LState) int {
	resourceType := L.CheckString(1)
	resourceName := L.CheckString(2)

	if currentStack == nil || currentStackManager == nil {
		L.Push(lua.LFalse)
		return 1
	}

	resource, _ := currentStackManager.GetResourceByStackAndName(currentStack.ID, resourceType, resourceName)
	L.Push(lua.LBool(resource != nil))
	return 1
}

// Helper functions specific to stack resources

func stackLuaTableToMap(t *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	t.ForEach(func(k, v lua.LValue) {
		key := k.String()
		result[key] = stackLuaValueToGo(v)
	})
	return result
}

func stackLuaValueToGo(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LTable:
		return stackLuaTableToMap(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LNilType:
		return nil
	default:
		return lv.String()
	}
}

func stackMapToLuaTable(L *lua.LState, m map[string]interface{}) *lua.LTable {
	table := L.NewTable()
	for k, v := range m {
		table.RawSetString(k, stackGoValueToLua(L, v))
	}
	return table
}

func stackResourceToLuaTable(L *lua.LState, r *stack.Resource) *lua.LTable {
	table := L.NewTable()
	table.RawSetString("id", lua.LString(r.ID))
	table.RawSetString("stack_id", lua.LString(r.StackID))
	table.RawSetString("type", lua.LString(r.Type))
	table.RawSetString("name", lua.LString(r.Name))
	table.RawSetString("module", lua.LString(r.Module))
	table.RawSetString("state", lua.LString(r.State))
	table.RawSetString("checksum", lua.LString(r.Checksum))
	table.RawSetString("properties", stackMapToLuaTable(L, r.Properties))
	return table
}

func stackGoValueToLua(L *lua.LState, val interface{}) lua.LValue {
	switch v := val.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case map[string]interface{}:
		return stackMapToLuaTable(L, v)
	case *stack.Resource:
		return stackResourceToLuaTable(L, v)
	case nil:
		return lua.LNil
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}
