package luainterface

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestParseLuaScript_RunIfFunction tests parsing tasks with run_if as a Lua function
func TestParseLuaScript_RunIfFunction(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_runif_func.sloth")

	script := `
workflow.define("test_group", {
	description = "Test group with run_if function",
	tasks = {
		{
			name = "conditional_task",
			description = "Task with run_if function",
			command = "echo 'running'",
			run_if = function(params, input)
				return params.should_run == "true"
			end
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.Equal(t, "conditional_task", task.Name)
	assert.NotNil(t, task.RunIfFunc, "RunIfFunc should be parsed")
	assert.Empty(t, task.RunIf, "RunIf string should be empty when function is used")
}

// TestParseLuaScript_RunIfString tests parsing tasks with run_if as a shell command string
func TestParseLuaScript_RunIfString(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_runif_str.sloth")

	script := `
workflow({
	name = "test_group",
	description = "Test group with run_if string",
	tasks = {
		{
			name = "conditional_task",
			description = "Task with run_if string",
			command = "echo 'running'",
			run_if = "test -f /tmp/trigger_file"
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.Equal(t, "conditional_task", task.Name)
	assert.Equal(t, "test -f /tmp/trigger_file", task.RunIf)
	assert.Nil(t, task.RunIfFunc, "RunIfFunc should be nil when string is used")
}

// TestParseLuaScript_AbortIfFunction tests parsing tasks with abort_if as a Lua function
func TestParseLuaScript_AbortIfFunction(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_abortif_func.sloth")

	script := `
workflow({
	name = "test_group",
	description = "Test group with abort_if function",
	tasks = {
		{
			name = "abortable_task",
			description = "Task with abort_if function",
			command = "echo 'running'",
			abort_if = function(params, input)
				return params.should_abort == "true"
			end
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.Equal(t, "abortable_task", task.Name)
	assert.NotNil(t, task.AbortIfFunc, "AbortIfFunc should be parsed")
	assert.Empty(t, task.AbortIf, "AbortIf string should be empty when function is used")
}

// TestParseLuaScript_AbortIfString tests parsing tasks with abort_if as a shell command string
func TestParseLuaScript_AbortIfString(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_abortif_str.sloth")

	script := `
workflow({
	name = "test_group",
	description = "Test group with abort_if string",
	tasks = {
		{
			name = "abortable_task",
			description = "Task with abort_if string",
			command = "echo 'running'",
			abort_if = "test -f /tmp/abort_trigger"
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.Equal(t, "abortable_task", task.Name)
	assert.Equal(t, "test -f /tmp/abort_trigger", task.AbortIf)
	assert.Nil(t, task.AbortIfFunc, "AbortIfFunc should be nil when string is used")
}

// TestParseLuaScript_BothRunIfAndAbortIf tests parsing tasks with both run_if and abort_if
func TestParseLuaScript_BothRunIfAndAbortIf(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_both_conditions.sloth")

	script := `
workflow({
	name = "test_group",
	description = "Test group with both conditions",
	tasks = {
		{
			name = "complex_task",
			description = "Task with both run_if and abort_if",
			command = "echo 'running'",
			run_if = function(params, input)
				return params.enabled == "true"
			end,
			abort_if = "test -f /tmp/emergency_stop"
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.Equal(t, "complex_task", task.Name)
	assert.NotNil(t, task.RunIfFunc, "RunIfFunc should be parsed")
	assert.Equal(t, "test -f /tmp/emergency_stop", task.AbortIf)
}

// TestParseLuaScript_NoConditions tests parsing tasks without any conditions
func TestParseLuaScript_NoConditions(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_no_conditions.sloth")

	script := `
workflow({
	name = "test_group",
	description = "Test group without conditions",
	tasks = {
		{
			name = "simple_task",
			description = "Simple task",
			command = "echo 'running'"
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.Equal(t, "simple_task", task.Name)
	assert.Nil(t, task.RunIfFunc)
	assert.Empty(t, task.RunIf)
	assert.Nil(t, task.AbortIfFunc)
	assert.Empty(t, task.AbortIf)
}

// TestExecuteLuaFunction_RunIfTrue tests executing a run_if function that returns true
func TestExecuteLuaFunction_RunIfTrue(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a simple run_if function that returns true
	script := `
function should_run(params, input)
	return true
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("should_run").(*lua.LFunction)
	params := map[string]string{}
	input := L.NewTable()

	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Function should return true")
}

// TestExecuteLuaFunction_RunIfFalse tests executing a run_if function that returns false
func TestExecuteLuaFunction_RunIfFalse(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function should_not_run(params, input)
	return false
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("should_not_run").(*lua.LFunction)
	params := map[string]string{}
	input := L.NewTable()

	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result, "Function should return false")
}

// TestExecuteLuaFunction_RunIfWithParams tests run_if function with parameters
func TestExecuteLuaFunction_RunIfWithParams(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_params(params, input)
	return params.env == "production"
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_params").(*lua.LFunction)
	input := L.NewTable()

	// Test with production
	params := map[string]string{"env": "production"}
	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Should return true for production env")

	// Test with dev
	params = map[string]string{"env": "dev"}
	result, _, _, err = ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result, "Should return false for dev env")
}

// TestExecuteLuaFunction_RunIfWithInput tests run_if function with input from dependencies
func TestExecuteLuaFunction_RunIfWithInput(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_input(params, input)
	return input.previous_task_status == "success"
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_input").(*lua.LFunction)
	params := map[string]string{}

	// Test with success status
	input := L.NewTable()
	input.RawSetString("previous_task_status", lua.LString("success"))
	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Should return true when previous task succeeded")

	// Test with failure status
	input = L.NewTable()
	input.RawSetString("previous_task_status", lua.LString("failed"))
	result, _, _, err = ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result, "Should return false when previous task failed")
}

// TestExecuteLuaFunction_AbortIfError tests abort_if function that returns an error
func TestExecuteLuaFunction_AbortIfError(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_abort_with_error(params, input)
	error("Intentional abort error")
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	abortIfFunc := L.GetGlobal("check_abort_with_error").(*lua.LFunction)
	params := map[string]string{}
	input := L.NewTable()

	_, _, _, err = ExecuteLuaFunction(L, abortIfFunc, params, input, 1, nil)
	assert.Error(t, err, "Should return error from Lua function")
	assert.Contains(t, err.Error(), "Intentional abort error")
}

// TestExecuteLuaFunction_ComplexCondition tests complex conditional logic
func TestExecuteLuaFunction_ComplexCondition(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function complex_condition(params, input)
	local env = params.env or "dev"
	local hour = tonumber(params.hour) or 0
	local is_prod = env == "production"
	local is_business_hours = hour >= 9 and hour <= 17
	
	-- Only run in production during business hours
	return is_prod and is_business_hours
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("complex_condition").(*lua.LFunction)
	input := L.NewTable()

	testCases := []struct {
		name     string
		params   map[string]string
		expected bool
	}{
		{
			name:     "Production during business hours",
			params:   map[string]string{"env": "production", "hour": "10"},
			expected: true,
		},
		{
			name:     "Production outside business hours",
			params:   map[string]string{"env": "production", "hour": "22"},
			expected: false,
		},
		{
			name:     "Dev during business hours",
			params:   map[string]string{"env": "dev", "hour": "10"},
			expected: false,
		},
		{
			name:     "Dev outside business hours",
			params:   map[string]string{"env": "dev", "hour": "22"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _, _, err := ExecuteLuaFunction(L, runIfFunc, tc.params, input, 1, nil)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result, fmt.Sprintf("Expected %v for %s", tc.expected, tc.name))
		})
	}
}

// TestParseLuaTask_AllConditionalFields tests that parseLuaTask correctly handles all conditional fields
func TestParseLuaTask_AllConditionalFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Register a Lua function for testing
	script := `
function my_run_if(params, input)
	return true
end

function my_abort_if(params, input)
	return false
end

test_task = {
	name = "test_task",
	description = "Test task with all conditional fields",
	command = "echo 'test'",
	run_if = my_run_if,
	abort_if = my_abort_if,
}
`
	err := L.DoString(script)
	require.NoError(t, err)

	taskTable := L.GetGlobal("test_task").(*lua.LTable)
	task := parseLuaTask(L, taskTable)

	assert.Equal(t, "test_task", task.Name)
	assert.NotNil(t, task.RunIfFunc)
	assert.NotNil(t, task.AbortIfFunc)
}

// TestConditionalEdgeCases tests edge cases for conditional logic
func TestConditionalEdgeCases(t *testing.T) {
	t.Run("RunIf returns non-boolean", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()

		script := `
function returns_string(params, input)
	return "not a boolean"
end
`
		err := L.DoString(script)
		require.NoError(t, err)

		runIfFunc := L.GetGlobal("returns_string").(*lua.LFunction)
		params := map[string]string{}
		input := L.NewTable()

		// ExecuteLuaFunction expects boolean as first return
		// If non-boolean is returned, it should return false and set message
		result, msg, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
		require.NoError(t, err)
		assert.False(t, result, "Should return false for non-boolean return type")
		assert.Contains(t, msg, "unexpected first return type")
	})

	t.Run("Nil function pointers", func(t *testing.T) {
		task := types.Task{
			Name:        "test",
			RunIfFunc:   nil,
			AbortIfFunc: nil,
			RunIf:       "",
			AbortIf:     "",
		}

		assert.Nil(t, task.RunIfFunc)
		assert.Nil(t, task.AbortIfFunc)
		assert.Empty(t, task.RunIf)
		assert.Empty(t, task.AbortIf)
	})
}

// TestExecuteLuaFunction_RunIfWithNilParams tests run_if with nil parameters
func TestExecuteLuaFunction_RunIfWithNilParams(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_nil_params(params, input)
	if params == nil then
		return false
	end
	return true
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_nil_params").(*lua.LFunction)
	input := L.NewTable()

	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, nil, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Should handle nil params gracefully")
}

// TestExecuteLuaFunction_RunIfWithNilInput tests run_if with nil input
func TestExecuteLuaFunction_RunIfWithNilInput(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_nil_input(params, input)
	if input == nil then
		return true
	end
	return false
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_nil_input").(*lua.LFunction)
	params := map[string]string{}

	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, nil, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Should handle nil input gracefully")
}

// TestExecuteLuaFunction_MultipleReturnValues tests functions returning multiple values
func TestExecuteLuaFunction_MultipleReturnValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function multi_return(params, input)
	return true, "success message", "extra data"
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("multi_return").(*lua.LFunction)
	params := map[string]string{}
	input := L.NewTable()

	result, msg, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result)
	assert.Equal(t, "success message", msg)
}

// TestParseLuaScript_MixedConditionTypes tests mixing function and string conditions
func TestParseLuaScript_MixedConditionTypes(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_mixed.sloth")

	script := `
workflow({
	name = "mixed_group",
	description = "Mixed condition types",
	tasks = {
		{
			name = "task1",
			command = "echo 'task1'",
			run_if = function(params, input) return true end
		},
		{
			name = "task2",
			command = "echo 'task2'",
			run_if = "test -f /tmp/file"
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "mixed_group")

	group := taskGroups["mixed_group"]
	require.Len(t, group.Tasks, 2)
	assert.NotNil(t, group.Tasks[0].RunIfFunc)
	assert.Empty(t, group.Tasks[1].RunIfFunc)
	assert.NotEmpty(t, group.Tasks[1].RunIf)
}

// TestExecuteLuaFunction_EmptyParamsMap tests with empty params map
func TestExecuteLuaFunction_EmptyParamsMap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_empty_params(params, input)
	for k, v in pairs(params) do
		return false
	end
	return true
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_empty_params").(*lua.LFunction)
	params := map[string]string{}
	input := L.NewTable()

	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Should handle empty params map")
}

// TestExecuteLuaFunction_ComplexInputTable tests with complex input table
func TestExecuteLuaFunction_ComplexInputTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_complex_input(params, input)
	if input.nested and input.nested.value == "test" then
		return true
	end
	return false
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_complex_input").(*lua.LFunction)
	params := map[string]string{}

	input := L.NewTable()
	nested := L.NewTable()
	nested.RawSetString("value", lua.LString("test"))
	input.RawSetString("nested", nested)

	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result)
}

// TestParseLuaScript_EmptyTaskList tests workflow with no tasks
func TestParseLuaScript_EmptyTaskList(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_empty.sloth")

	script := `
workflow({
	name = "empty_group",
	description = "Empty task list",
	tasks = {}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "empty_group")

	group := taskGroups["empty_group"]
	assert.Len(t, group.Tasks, 0)
}

// TestParseLuaScript_MultipleWorkflows tests multiple workflows in one script
func TestParseLuaScript_MultipleWorkflows(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_multiple.sloth")

	script := `
workflow({
	name = "group1",
	tasks = {
		{name = "task1", command = "echo 'task1'"}
	}
})

workflow({
	name = "group2",
	tasks = {
		{name = "task2", command = "echo 'task2'"}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	assert.Len(t, taskGroups, 2)
	assert.Contains(t, taskGroups, "group1")
	assert.Contains(t, taskGroups, "group2")
}

// TestExecuteLuaFunction_NumericParams tests with numeric parameter values
func TestExecuteLuaFunction_NumericParams(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_numeric(params, input)
	local threshold = tonumber(params.threshold) or 0
	return threshold > 50
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_numeric").(*lua.LFunction)
	input := L.NewTable()

	params := map[string]string{"threshold": "100"}
	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result)

	params = map[string]string{"threshold": "25"}
	result, _, _, err = ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestExecuteLuaFunction_BooleanParamStrings tests boolean params as strings
func TestExecuteLuaFunction_BooleanParamStrings(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_bool_string(params, input)
	return params.enabled == "true"
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_bool_string").(*lua.LFunction)
	input := L.NewTable()

	params := map[string]string{"enabled": "true"}
	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result)

	params = map[string]string{"enabled": "false"}
	result, _, _, err = ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result)
}

// TestExecuteLuaFunction_AbortIfTrue tests abort condition returning true
func TestExecuteLuaFunction_AbortIfTrue(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function should_abort(params, input)
	return params.abort == "true"
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	abortIfFunc := L.GetGlobal("should_abort").(*lua.LFunction)
	params := map[string]string{"abort": "true"}
	input := L.NewTable()

	result, _, _, err := ExecuteLuaFunction(L, abortIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Abort condition should return true")
}

// TestExecuteLuaFunction_AbortIfFalse tests abort condition returning false
func TestExecuteLuaFunction_AbortIfFalse(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function should_not_abort(params, input)
	return params.abort == "true"
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	abortIfFunc := L.GetGlobal("should_not_abort").(*lua.LFunction)
	params := map[string]string{"abort": "false"}
	input := L.NewTable()

	result, _, _, err := ExecuteLuaFunction(L, abortIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result, "Abort condition should return false")
}

// TestParseLuaScript_BothFunctionConditions tests task with both conditions as functions
func TestParseLuaScript_BothFunctionConditions(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_both_func.sloth")

	script := `
workflow({
	name = "test_group",
	tasks = {
		{
			name = "complex_task",
			command = "echo 'test'",
			run_if = function(params, input) return true end,
			abort_if = function(params, input) return false end
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.NotNil(t, task.RunIfFunc)
	assert.NotNil(t, task.AbortIfFunc)
	assert.Empty(t, task.RunIf)
	assert.Empty(t, task.AbortIf)
}

// TestParseLuaScript_BothStringConditions tests task with both conditions as strings
func TestParseLuaScript_BothStringConditions(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test_both_str.sloth")

	script := `
workflow({
	name = "test_group",
	tasks = {
		{
			name = "complex_task",
			command = "echo 'test'",
			run_if = "test -f /tmp/run",
			abort_if = "test -f /tmp/abort"
		}
	}
})
`
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	require.NoError(t, err)

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Contains(t, taskGroups, "test_group")

	group := taskGroups["test_group"]
	require.Len(t, group.Tasks, 1)

	task := group.Tasks[0]
	assert.Nil(t, task.RunIfFunc)
	assert.Nil(t, task.AbortIfFunc)
	assert.Equal(t, "test -f /tmp/run", task.RunIf)
	assert.Equal(t, "test -f /tmp/abort", task.AbortIf)
}

// TestExecuteLuaFunction_StringReturnAsMessage tests function returning only a string
func TestExecuteLuaFunction_StringReturnAsMessage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function returns_only_string(params, input)
	return "just a message"
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("returns_only_string").(*lua.LFunction)
	params := map[string]string{}
	input := L.NewTable()

	result, msg, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result, "Non-boolean should result in false")
	assert.NotEmpty(t, msg, "Message should be set")
}

// TestExecuteLuaFunction_NilReturn tests function returning nil
func TestExecuteLuaFunction_NilReturn(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function returns_nil(params, input)
	return nil
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("returns_nil").(*lua.LFunction)
	params := map[string]string{}
	input := L.NewTable()

	result, msg, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.False(t, result, "Nil should result in false")
	assert.NotEmpty(t, msg, "Message should indicate unexpected type")
}

// TestExecuteLuaFunction_LargeParamSet tests with many parameters
func TestExecuteLuaFunction_LargeParamSet(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	script := `
function check_many_params(params, input)
	local count = 0
	for k, v in pairs(params) do
		count = count + 1
	end
	return count >= 10
end
`
	err := L.DoString(script)
	require.NoError(t, err)

	runIfFunc := L.GetGlobal("check_many_params").(*lua.LFunction)
	input := L.NewTable()

	params := make(map[string]string)
	for i := 0; i < 15; i++ {
		params[fmt.Sprintf("param%d", i)] = fmt.Sprintf("value%d", i)
	}

	result, _, _, err := ExecuteLuaFunction(L, runIfFunc, params, input, 1, nil)
	require.NoError(t, err)
	assert.True(t, result, "Should handle many parameters")
}
