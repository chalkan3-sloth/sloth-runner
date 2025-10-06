package luainterface

import (
	"context"
	"os"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestTaskExecutionFlow(t *testing.T) {
	script := `
workflow({
	name = "test_group",
	description = "Test execution flow",
	tasks = {
		{
			name = "setup",
			description = "Setup task",
			command = "echo setup"
		},
		{
			name = "main",
			description = "Main task",
			command = "echo main"
		},
		{
			name = "cleanup",
			description = "Cleanup task",
			command = "echo cleanup"
		}
	}
})
`

	tmpfile, err := os.CreateTemp("", "test-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	groups, err := ParseLuaScript(context.Background(), tmpfile.Name(), nil)
	if err != nil {
		t.Fatalf("ParseLuaScript() error = %v", err)
	}

	group := groups["test_group"]
	if len(group.Tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(group.Tasks))
	}

	// Verify task order
	expectedNames := []string{"setup", "main", "cleanup"}
	for i, task := range group.Tasks {
		if task.Name != expectedNames[i] {
			t.Errorf("Task %d: expected name '%s', got '%s'", i, expectedNames[i], task.Name)
		}
	}
}

func TestConditionalTaskExecution(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
local should_run = true

local tasks = {}

if should_run then
	table.insert(tasks, {
		name = "conditional_task",
		command = "echo running"
	})
end

workflow({
	name = "test_group",
	description = "Test conditional execution",
	tasks = tasks
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute conditional script: %v", err)
	}

	// Verify task was added
	workflows := L.GetGlobal("__workflows__")
	if workflows.Type() != lua.LTTable {
		t.Fatal("__workflows__ not a table")
	}

	testGroup := workflows.(*lua.LTable).RawGetString("test_group")
	if testGroup.Type() != lua.LTTable {
		t.Fatal("test_group not a table")
	}

	tasks := testGroup.(*lua.LTable).RawGetString("tasks")
	if tasks.Type() != lua.LTTable {
		t.Fatal("tasks not a table")
	}

	// Check task count
	count := 0
	tasks.(*lua.LTable).ForEach(func(_, _ lua.LValue) {
		count++
	})

	if count != 1 {
		t.Errorf("Expected 1 task, got %d", count)
	}
}

func TestTaskWithVariablesExecution(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Set up variables
	valuesTable := L.NewTable()
	valuesTable.RawSetString("app_name", lua.LString("myapp"))
	valuesTable.RawSetString("version", lua.LString("1.0.0"))
	L.SetGlobal("values", valuesTable)

	script := `
local app = values.app_name
local ver = values.version

workflow({
	name = "test_group",
	description = "Deploy " .. app .. " version " .. ver,
	tasks = {
		{
			name = "deploy_" .. app,
			command = "echo deploying " .. app .. " " .. ver
		}
	}
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script with variables: %v", err)
	}

	workflows := L.GetGlobal("__workflows__")
	testGroup := workflows.(*lua.LTable).RawGetString("test_group")
	description := testGroup.(*lua.LTable).RawGetString("description").String()

	expected := "Deploy myapp version 1.0.0"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}

func TestLoopTaskGeneration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
local services = {"web", "api", "worker"}

local tasks = {}

for _, service in ipairs(services) do
	table.insert(tasks, {
		name = "deploy_" .. service,
		command = "echo deploying " .. service
	})
end

workflow({
	name = "deploy",
	description = "Deploy all services",
	tasks = tasks
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute loop generation: %v", err)
	}

	workflows := L.GetGlobal("__workflows__")
	deployGroup := workflows.(*lua.LTable).RawGetString("deploy")
	tasks := deployGroup.(*lua.LTable).RawGetString("tasks")

	count := 0
	tasks.(*lua.LTable).ForEach(func(_, _ lua.LValue) {
		count++
	})

	if count != 3 {
		t.Errorf("Expected 3 tasks from loop, got %d", count)
	}
}

func TestTaskWithDelegation(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	script := `
workflow({
	name = "distributed",
	description = "Distributed tasks",
	tasks = {
		{
			name = "local_task",
			command = "echo local"
		},
		{
			name = "remote_task",
			command = "echo remote",
			delegate_to = "remote-agent"
		}
	}
})
`

	tmpfile, err := os.CreateTemp("", "test-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	groups, err := ParseLuaScript(context.Background(), tmpfile.Name(), nil)
	if err != nil {
		t.Fatalf("ParseLuaScript() error = %v", err)
	}

	group := groups["distributed"]
	
	if group.Tasks[0].DelegateTo != "" {
		t.Errorf("First task should not be delegated, got: %s", group.Tasks[0].DelegateTo)
	}

	if group.Tasks[1].DelegateTo != "remote-agent" {
		t.Errorf("Second task should be delegated to 'remote-agent', got: %s", group.Tasks[1].DelegateTo)
	}
}

func TestTaskWithHelperFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
-- Define helper function
local function make_task(name, cmd)
	return {
		name = name,
		command = cmd,
		description = "Generated task: " .. name
	}
end

workflow({
	name = "helpers",
	description = "Using helper functions",
	tasks = {
		make_task("task1", "echo 1"),
		make_task("task2", "echo 2"),
		make_task("task3", "echo 3")
	}
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute helper functions: %v", err)
	}

	workflows := L.GetGlobal("__workflows__")
	helpersGroup := workflows.(*lua.LTable).RawGetString("helpers")
	tasks := helpersGroup.(*lua.LTable).RawGetString("tasks")

	count := 0
	tasks.(*lua.LTable).ForEach(func(_, taskValue lua.LValue) {
		count++
		if taskValue.Type() == lua.LTTable {
			taskTable := taskValue.(*lua.LTable)
			desc := taskTable.RawGetString("description").String()
			if desc == "" {
				t.Error("Task description is empty")
			}
		}
	})

	if count != 3 {
		t.Errorf("Expected 3 tasks, got %d", count)
	}
}

func TestNestedModuleUsage(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
local config = {
	environment = "production",
	replicas = 3
}

local config_json = json.encode(config)
local config_yaml = yaml.encode(config)

workflow({
	name = "config_tasks",
	description = "Tasks using modules",
	tasks = {
		{
			name = "show_json",
			command = "echo " .. config_json
		},
		{
			name = "show_yaml",
			command = "echo " .. config_yaml
		}
	}
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute nested module usage: %v", err)
	}

	// Verify __workflows__ was created
	workflows := L.GetGlobal("__workflows__")
	if workflows.Type() != lua.LTTable {
		t.Fatal("__workflows__ not created")
	}
}

func TestTaskErrorHandlingExecution(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Test script with potential error
	script := `
local tasks = {
	{
		name = "safe_task",
		command = "echo safe"
	}
}

-- Try to add invalid task (should not crash)
pcall(function()
	table.insert(tasks, nil)
end)

workflow({
	name = "error_handling",
	description = "Test error handling",
	tasks = tasks
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute error handling script: %v", err)
	}

	workflows := L.GetGlobal("__workflows__")
	errorGroup := workflows.(*lua.LTable).RawGetString("error_handling")
	tasks := errorGroup.(*lua.LTable).RawGetString("tasks")

	// Should still have at least the safe task
	count := 0
	tasks.(*lua.LTable).ForEach(func(_, value lua.LValue) {
		if value.Type() != lua.LTNil {
			count++
		}
	})

	if count < 1 {
		t.Error("Expected at least 1 valid task")
	}
}

func TestDynamicTaskConfiguration(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Set environment variable
	os.Setenv("TEST_MODE", "integration")
	defer os.Unsetenv("TEST_MODE")

	script := `
local mode = env.get("TEST_MODE")

workflow({
	name = "dynamic",
	description = "Dynamic configuration based on " .. mode,
	tasks = {
		{
			name = "test_task",
			command = "echo " .. mode
		}
	}
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute dynamic configuration: %v", err)
	}

	workflows := L.GetGlobal("__workflows__")
	dynamicGroup := workflows.(*lua.LTable).RawGetString("dynamic")
	description := dynamicGroup.(*lua.LTable).RawGetString("description").String()

	if description != "Dynamic configuration based on integration" {
		t.Errorf("Expected 'integration' in description, got: %s", description)
	}
}

func TestTaskWithWorkdirExecution(t *testing.T) {
	script := `
workflow({
	name = "workdir_test",
	description = "Test workdir",
	workdir = "/tmp/test",
	create_workdir_before_run = true,
	tasks = {
		{
			name = "task_in_workdir",
			command = "pwd"
		}
	}
})
`

	tmpfile, err := os.CreateTemp("", "test-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	groups, err := ParseLuaScript(context.Background(), tmpfile.Name(), nil)
	if err != nil {
		t.Fatalf("ParseLuaScript() error = %v", err)
	}

	group := groups["workdir_test"]
	if group.Workdir != "/tmp/test" {
		t.Errorf("Expected workdir '/tmp/test', got: %s", group.Workdir)
	}

	if !group.CreateWorkdirBeforeRun {
		t.Error("Expected CreateWorkdirBeforeRun to be true")
	}
}

func TestComplexTaskScenario(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
-- Configuration
local config = {
	app_name = "myapp",
	environments = {"dev", "staging", "prod"},
	version = "2.0.0"
}

-- Helper to create deployment task
local function create_deploy_task(env)
	return {
		name = "deploy_to_" .. env,
		description = "Deploy " .. config.app_name .. " to " .. env,
		command = strings.join({
			"deploy",
			"--app", config.app_name,
			"--env", env,
			"--version", config.version
		}, " ")
	}
end

local tasks = {}

-- Generate tasks for each environment
for _, env in ipairs(config.environments) do
	table.insert(tasks, create_deploy_task(env))
end

-- Add post-deployment task
table.insert(tasks, {
	name = "notify",
	description = "Send deployment notification",
	command = "echo Deployment complete"
})

-- Build workflow
workflow({
	name = "multi_env_deploy",
	description = "Multi-environment deployment",
	workdir = "/deployments",
	tasks = tasks
})
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute complex scenario: %v", err)
	}

	workflows := L.GetGlobal("__workflows__")
	deployGroup := workflows.(*lua.LTable).RawGetString("multi_env_deploy")
	tasks := deployGroup.(*lua.LTable).RawGetString("tasks")

	// Should have 3 deploy tasks + 1 notify task = 4 total
	count := 0
	tasks.(*lua.LTable).ForEach(func(_, _ lua.LValue) {
		count++
	})

	if count != 4 {
		t.Errorf("Expected 4 tasks in complex scenario, got %d", count)
	}
}
