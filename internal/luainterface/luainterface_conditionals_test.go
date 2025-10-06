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
workflow({
	name = "test_group",
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
