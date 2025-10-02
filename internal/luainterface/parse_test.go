package luainterface

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestParseLuaScript_BasicTaskGroup(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				description = "First task",
				command = "echo hello"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Len(t, taskGroups, 1)

	group, exists := taskGroups["test_group"]
	require.True(t, exists)
	assert.Equal(t, "Test group", group.Description)
	assert.Len(t, group.Tasks, 1)
	assert.Equal(t, "task1", group.Tasks[0].Name)
	assert.Equal(t, "echo hello", group.Tasks[0].CommandStr)
}

func TestParseLuaScript_WithWorkdir(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		workdir = "/tmp/test",
		create_workdir_before_run = true,
		tasks = {
			task1 = {
				name = "task1",
				command = "ls",
				workdir = "/tmp/task1"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	group := taskGroups["test_group"]
	assert.Equal(t, "/tmp/test", group.Workdir)
	assert.True(t, group.CreateWorkdirBeforeRun)
	assert.Equal(t, "/tmp/task1", group.Tasks[0].Workdir)
}

func TestParseLuaScript_WithDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo task1"
			},
			task2 = {
				name = "task2",
				command = "echo task2",
				depends_on = "task1"
			},
			task3 = {
				name = "task3",
				command = "echo task3",
				depends_on = {"task1", "task2"}
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	tasks := taskGroups["test_group"].Tasks
	
	// Find tasks by name
	var task2, task3 *struct {
		DependsOn []string
	}
	for i := range tasks {
		if tasks[i].Name == "task2" {
			task2 = &struct{ DependsOn []string }{tasks[i].DependsOn}
		}
		if tasks[i].Name == "task3" {
			task3 = &struct{ DependsOn []string }{tasks[i].DependsOn}
		}
	}

	require.NotNil(t, task2)
	require.NotNil(t, task3)
	assert.Equal(t, []string{"task1"}, task2.DependsOn)
	assert.ElementsMatch(t, []string{"task1", "task2"}, task3.DependsOn)
}

func TestParseLuaScript_WithArtifacts(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo test",
				artifacts = "output.txt"
			},
			task2 = {
				name = "task2",
				command = "echo test",
				artifacts = {"file1.txt", "file2.txt"}
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	tasks := taskGroups["test_group"].Tasks
	
	for _, task := range tasks {
		if task.Name == "task1" {
			assert.Equal(t, []string{"output.txt"}, task.Artifacts)
		}
		if task.Name == "task2" {
			assert.ElementsMatch(t, []string{"file1.txt", "file2.txt"}, task.Artifacts)
		}
	}
}

func TestParseLuaScript_WithRetries(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo test",
				retries = 3,
				timeout = "30s"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	task := taskGroups["test_group"].Tasks[0]
	assert.Equal(t, 3, task.Retries)
	assert.Equal(t, "30s", task.Timeout)
}

func TestParseLuaScript_WithAsync(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo test",
				async = true
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	task := taskGroups["test_group"].Tasks[0]
	assert.True(t, task.Async)
}

func TestParseLuaScript_WithParams(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo test",
				params = {
					key1 = "value1",
					key2 = "value2"
				}
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	task := taskGroups["test_group"].Tasks[0]
	assert.Equal(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, task.Params)
}

func TestParseLuaScript_WithDelegateTo_String(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		delegate_to = "remote-agent",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo test"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	group := taskGroups["test_group"]
	assert.Equal(t, "remote-agent", group.DelegateTo)
}

func TestParseLuaScript_WithDelegateTo_Table(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		delegate_to = {
			agent = "remote-agent",
			parallel = true
		},
		tasks = {
			task1 = {
				name = "task1",
				command = "echo test"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	group := taskGroups["test_group"]
	delegateMap, ok := group.DelegateTo.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "remote-agent", delegateMap["agent"])
	assert.Equal(t, true, delegateMap["parallel"])
}

func TestParseLuaScript_WithValues(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = Values.description,
		tasks = {
			task1 = {
				name = "task1",
				command = "echo " .. Values.message
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	L := lua.NewState()
	defer L.Close()
	
	valuesTable := L.NewTable()
	valuesTable.RawSetString("description", lua.LString("Test from values"))
	valuesTable.RawSetString("message", lua.LString("hello world"))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, valuesTable)
	require.NoError(t, err)

	group := taskGroups["test_group"]
	assert.Equal(t, "Test from values", group.Description)
	assert.Contains(t, group.Tasks[0].CommandStr, "hello world")
}

func TestParseLuaScript_InvalidScript(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		-- Invalid syntax
		description = 
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	_, err := ParseLuaScript(context.Background(), scriptPath, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute Lua script")
}

func TestParseLuaScript_NoTaskDefinitions(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
-- No TaskDefinitions
local x = 42
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	_, err := ParseLuaScript(context.Background(), scriptPath, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid task definitions found")
}

func TestParseLuaScript_WithConsumes(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo test",
				artifacts = "output.txt"
			},
			task2 = {
				name = "task2",
				command = "cat input.txt",
				consumes = "output.txt",
				depends_on = "task1"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	tasks := taskGroups["test_group"].Tasks
	for _, task := range tasks {
		if task.Name == "task2" {
			assert.Equal(t, []string{"output.txt"}, task.Consumes)
		}
	}
}

func TestParseLuaScript_WithNextIfFail(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	test_group = {
		description = "Test group",
		tasks = {
			task1 = {
				name = "task1",
				command = "test_command",
				next_if_fail = "fallback_task"
			},
			fallback_task = {
				name = "fallback_task",
				command = "echo fallback"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)

	tasks := taskGroups["test_group"].Tasks
	for _, task := range tasks {
		if task.Name == "task1" {
			assert.Equal(t, []string{"fallback_task"}, task.NextIfFail)
		}
	}
}

func TestParseLuaScript_MultipleTaskGroups(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sloth")

	script := `
TaskDefinitions = {
	group1 = {
		description = "First group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo group1"
			}
		}
	},
	group2 = {
		description = "Second group",
		tasks = {
			task1 = {
				name = "task1",
				command = "echo group2"
			}
		}
	}
}
`
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0644))

	taskGroups, err := ParseLuaScript(context.Background(), scriptPath, nil)
	require.NoError(t, err)
	require.Len(t, taskGroups, 2)

	assert.Contains(t, taskGroups, "group1")
	assert.Contains(t, taskGroups, "group2")
	assert.Equal(t, "First group", taskGroups["group1"].Description)
	assert.Equal(t, "Second group", taskGroups["group2"].Description)
}
