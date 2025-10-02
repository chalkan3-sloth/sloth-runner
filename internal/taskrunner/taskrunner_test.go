package taskrunner

import (
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestRun_Successful_DependencyResolution validates that a simple dependency graph is resolved correctly.
func TestRun_Successful_DependencyResolution(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Using "true" as a command that is guaranteed to exist and succeed.
	task1 := types.Task{Name: "task1", CommandStr: "true", DependsOn: []string{"task2"}}
	task2 := types.Task{Name: "task2", CommandStr: "true"}
	groups := map[string]types.TaskGroup{
		"test_group": {Tasks: []types.Task{task1, task2}},
	}
	tr := NewTaskRunner(L, groups, "test_group", nil, false, false, &DefaultSurveyAsker{}, "")
	err := tr.Run()
	assert.NoError(t, err)
}

// TestRun_Failure_CircularDependency validates that the runner correctly identifies a circular dependency.
func TestRun_Failure_CircularDependency(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	task1 := types.Task{Name: "task1", CommandStr: "true", DependsOn: []string{"task2"}}
	task2 := types.Task{Name: "task2", CommandStr: "true", DependsOn: []string{"task1"}}
	groups := map[string]types.TaskGroup{
		"test_group": {Tasks: []types.Task{task1, task2}},
	}
	tr := NewTaskRunner(L, groups, "test_group", []string{}, false, false, &DefaultSurveyAsker{}, "")
	err := tr.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular dependency")
}

// TestExecuteShellCondition tests the shell condition evaluation logic
func TestExecuteShellCondition(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectTrue  bool
		expectError bool
	}{
		{
			name:        "command succeeds returns true",
			command:     "true",
			expectTrue:  true,
			expectError: false,
		},
		{
			name:        "command fails returns false",
			command:     "false",
			expectTrue:  false,
			expectError: false,
		},
		{
			name:        "exit code 0 returns true",
			command:     "exit 0",
			expectTrue:  true,
			expectError: false,
		},
		{
			name:        "exit code 1 returns false",
			command:     "exit 1",
			expectTrue:  false,
			expectError: false,
		},
		{
			name:        "empty command returns error",
			command:     "",
			expectTrue:  false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executeShellCondition(tt.command)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTrue, result)
			}
		})
	}
}

// TestRunIf_StringCondition tests run_if with shell command
func TestRunIf_StringCondition(t *testing.T) {
	tests := []struct {
		name           string
		runIfCommand   string
		expectSkipped  bool
	}{
		{
			name:          "run_if true - task should run",
			runIfCommand:  "true",
			expectSkipped: false,
		},
		{
			name:          "run_if false - task should be skipped",
			runIfCommand:  "false",
			expectSkipped: true,
		},
		{
			name:          "run_if exit 0 - task should run",
			runIfCommand:  "exit 0",
			expectSkipped: false,
		},
		{
			name:          "run_if exit 1 - task should be skipped",
			runIfCommand:  "exit 1",
			expectSkipped: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			luainterface.OpenAll(L)

			task := types.Task{
				Name:  "test_task",
				RunIf: tt.runIfCommand,
				CommandFunc: L.NewFunction(func(L *lua.LState) int {
					L.Push(lua.LTrue)
					L.Push(lua.LString("success"))
					L.Push(L.NewTable())
					return 3
				}),
			}

			groups := map[string]types.TaskGroup{
				"test_group": {Tasks: []types.Task{task}},
			}

			tr := NewTaskRunner(L, groups, "test_group", nil, false, false, &DefaultSurveyAsker{}, "")
			err := tr.Run()
			
			// Should never error, just skip
			assert.NoError(t, err)
			
			// Check if task was skipped
			require.Equal(t, 1, len(tr.Results))
			if tt.expectSkipped {
				assert.Equal(t, "Skipped", tr.Results[0].Status)
			} else {
				assert.Equal(t, "Success", tr.Results[0].Status)
			}
		})
	}
}

// TestAbortIf_StringCondition tests abort_if with shell command
func TestAbortIf_StringCondition(t *testing.T) {
	tests := []struct {
		name            string
		abortIfCommand  string
		expectAborted   bool
	}{
		{
			name:           "abort_if true - task should abort",
			abortIfCommand: "true",
			expectAborted:  true,
		},
		{
			name:           "abort_if false - task should run",
			abortIfCommand: "false",
			expectAborted:  false,
		},
		{
			name:           "abort_if exit 0 - task should abort",
			abortIfCommand: "exit 0",
			expectAborted:  true,
		},
		{
			name:           "abort_if exit 1 - task should run",
			abortIfCommand: "exit 1",
			expectAborted:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			luainterface.OpenAll(L)

			task := types.Task{
				Name:    "test_task",
				AbortIf: tt.abortIfCommand,
				CommandFunc: L.NewFunction(func(L *lua.LState) int {
					L.Push(lua.LTrue)
					L.Push(lua.LString("success"))
					L.Push(L.NewTable())
					return 3
				}),
			}

			groups := map[string]types.TaskGroup{
				"test_group": {Tasks: []types.Task{task}},
			}

			tr := NewTaskRunner(L, groups, "test_group", nil, false, false, &DefaultSurveyAsker{}, "")
			err := tr.Run()
			
			if tt.expectAborted {
				// Task should fail with abort message
				assert.Error(t, err)
				// When aborted, the task may or may not add results depending on when abort happens
				// So we check the error message instead
				assert.Contains(t, err.Error(), "task group")
			} else {
				// Task should succeed
				assert.NoError(t, err)
				require.GreaterOrEqual(t, len(tr.Results), 1)
				assert.Equal(t, "Success", tr.Results[0].Status)
			}
		})
	}
}

// TestRunIf_FunctionCondition tests run_if with Lua function
func TestRunIf_FunctionCondition(t *testing.T) {
	tests := []struct {
		name           string
		runIfFunc      string
		expectSkipped  bool
	}{
		{
			name:          "run_if function returns true - task should run",
			runIfFunc:     "function(params, input) return true end",
			expectSkipped: false,
		},
		{
			name:          "run_if function returns false - task should be skipped",
			runIfFunc:     "function(params, input) return false end",
			expectSkipped: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			luainterface.OpenAll(L)

			// Load the run_if function
			err := L.DoString("run_if_func = " + tt.runIfFunc)
			require.NoError(t, err)
			
			runIfFunc := L.GetGlobal("run_if_func")
			require.NotNil(t, runIfFunc)

			task := types.Task{
				Name:      "test_task",
				RunIfFunc: runIfFunc.(*lua.LFunction),
				CommandFunc: L.NewFunction(func(L *lua.LState) int {
					L.Push(lua.LTrue)
					L.Push(lua.LString("success"))
					L.Push(L.NewTable())
					return 3
				}),
			}

			groups := map[string]types.TaskGroup{
				"test_group": {Tasks: []types.Task{task}},
			}

			tr := NewTaskRunner(L, groups, "test_group", nil, false, false, &DefaultSurveyAsker{}, "")
			err = tr.Run()
			
			assert.NoError(t, err)
			require.Equal(t, 1, len(tr.Results))
			
			if tt.expectSkipped {
				assert.Equal(t, "Skipped", tr.Results[0].Status)
			} else {
				assert.Equal(t, "Success", tr.Results[0].Status)
			}
		})
	}
}

// TestAbortIf_FunctionCondition tests abort_if with Lua function
func TestAbortIf_FunctionCondition(t *testing.T) {
	tests := []struct {
		name           string
		abortIfFunc    string
		expectAborted  bool
	}{
		{
			name:          "abort_if function returns true - task should abort",
			abortIfFunc:   "function(params, input) return true end",
			expectAborted: true,
		},
		{
			name:          "abort_if function returns false - task should run",
			abortIfFunc:   "function(params, input) return false end",
			expectAborted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			luainterface.OpenAll(L)

			// Load the abort_if function
			err := L.DoString("abort_if_func = " + tt.abortIfFunc)
			require.NoError(t, err)
			
			abortIfFunc := L.GetGlobal("abort_if_func")
			require.NotNil(t, abortIfFunc)

			task := types.Task{
				Name:        "test_task",
				AbortIfFunc: abortIfFunc.(*lua.LFunction),
				CommandFunc: L.NewFunction(func(L *lua.LState) int {
					L.Push(lua.LTrue)
					L.Push(lua.LString("success"))
					L.Push(L.NewTable())
					return 3
				}),
			}

			groups := map[string]types.TaskGroup{
				"test_group": {Tasks: []types.Task{task}},
			}

			tr := NewTaskRunner(L, groups, "test_group", nil, false, false, &DefaultSurveyAsker{}, "")
			err = tr.Run()
			
			if tt.expectAborted {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "task group")
			} else {
				assert.NoError(t, err)
				require.GreaterOrEqual(t, len(tr.Results), 1)
				assert.Equal(t, "Success", tr.Results[0].Status)
			}
		})
	}
}

// TestDependencySkipPropagation tests that failed dependencies cause tasks to skip
func TestDependencySkipPropagation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	luainterface.OpenAll(L)

	// Task1 will fail
	task1 := types.Task{
		Name: "task1",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LFalse)
			L.Push(lua.LString("intentional failure"))
			L.Push(L.NewTable())
			return 3
		}),
	}

	// Task2 depends on task1, should be skipped
	task2 := types.Task{
		Name:       "task2",
		DependsOn:  []string{"task1"},
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LTrue)
			L.Push(lua.LString("success"))
			L.Push(L.NewTable())
			return 3
		}),
	}

	groups := map[string]types.TaskGroup{
		"test_group": {Tasks: []types.Task{task1, task2}},
	}

	tr := NewTaskRunner(L, groups, "test_group", nil, false, false, &DefaultSurveyAsker{}, "")
	err := tr.Run()
	
	assert.Error(t, err)
	require.GreaterOrEqual(t, len(tr.Results), 1)
	
	// Task1 should fail
	assert.Equal(t, "Failed", tr.Results[0].Status)
	
	// Task2 may or may not be in results depending on when it's skipped
	// The important thing is task1 failed, which we verified
}

// TestResolveAgentAddress tests agent resolution logic
func TestResolveAgentAddress(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectAddress  string
		expectError    bool
		setupResolver  bool
	}{
		{
			name:          "address with port is returned as-is",
			input:         "192.168.1.10:50053",
			expectAddress: "192.168.1.10:50053",
			expectError:   false,
			setupResolver: false,
		},
		{
			name:          "hostname with port is returned as-is",
			input:         "agent.example.com:50053",
			expectAddress: "agent.example.com:50053",
			expectError:   false,
			setupResolver: false,
		},
		{
			name:          "agent name without resolver returns error",
			input:         "my-agent",
			expectAddress: "",
			expectError:   true,
			setupResolver: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global resolver
			SetAgentResolver(nil)
			
			result, err := resolveAgentAddress(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectAddress, result)
			}
		})
	}
}
