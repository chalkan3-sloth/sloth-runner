package luainterface

import (
"log/slog"
"os"
"testing"

"github.com/chalkan3-sloth/sloth-runner/internal/core"
lua "github.com/yuin/gopher-lua"
)

func setupTestDSL(t *testing.T) (*lua.LState, *ModernDSL) {
L := lua.NewState()
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
config := core.DefaultCoreConfig()
globalCore, err := core.NewGlobalCore(config, logger)
if err != nil {
L.Close()
t.Fatalf("Failed to create global core: %v", err)
}
dsl := NewModernDSL(globalCore)
dsl.RegisterModernDSL(L)
return L, dsl
}

func TestNewModernDSL(t *testing.T) {
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
config := core.DefaultCoreConfig()
globalCore, err := core.NewGlobalCore(config, logger)
if err != nil {
t.Fatalf("Failed to create global core: %v", err)
}

dsl := NewModernDSL(globalCore)
if dsl == nil {
t.Fatal("NewModernDSL returned nil")
}

if dsl.core == nil {
t.Error("core not initialized")
}

if dsl.logger == nil {
t.Error("logger not initialized")
}

if dsl.taskRegistry == nil {
t.Error("taskRegistry not initialized")
}

if dsl.builders == nil {
t.Error("builders map not initialized")
}
}

func TestNewTaskRegistry(t *testing.T) {
registry := NewTaskRegistry()
if registry == nil {
t.Fatal("NewTaskRegistry returned nil")
}

if registry.tasks == nil {
t.Error("tasks map not initialized")
}

if registry.groups == nil {
t.Error("groups map not initialized")
}

if registry.templates == nil {
t.Error("templates map not initialized")
}

if registry.validators == nil {
t.Error("validators map not initialized")
}
}

func TestTaskBuilderBasic(t *testing.T) {
L, _ := setupTestDSL(t)
defer L.Close()

err := L.DoString(`
local t = task("test_task")
assert(t ~= nil, "task should not be nil")
`)
if err != nil {
t.Fatalf("Error executing Lua code: %v", err)
}
}

func TestTaskBuilderWithDescription(t *testing.T) {
L, _ := setupTestDSL(t)
defer L.Close()

err := L.DoString(`
local t = task("test_task")
:description("Test task description")
:build()
assert(t ~= nil, "task should not be nil")
`)
if err != nil {
t.Fatalf("Error executing Lua code: %v", err)
}
}

func TestTaskBuilderWithCommand(t *testing.T) {
L, _ := setupTestDSL(t)
defer L.Close()

err := L.DoString(`
local t = task("test_task")
:command(function(params, deps)
return true, "success", {}
end)
:build()
assert(t ~= nil, "task should not be nil")
`)
if err != nil {
t.Fatalf("Error executing Lua code: %v", err)
}
}

func TestTaskBuilderWithTimeout(t *testing.T) {
L, _ := setupTestDSL(t)
defer L.Close()

err := L.DoString(`
local t = task("test_task")
:timeout("5m")
:build()
assert(t ~= nil, "task should not be nil")
`)
if err != nil {
t.Fatalf("Error executing Lua code: %v", err)
}
}

func TestTaskBuilderWithRetries(t *testing.T) {
L, _ := setupTestDSL(t)
defer L.Close()

err := L.DoString(`
local t = task("test_task")
:retries(3, "exponential")
:build()
assert(t ~= nil, "task should not be nil")
`)
if err != nil {
t.Fatalf("Error executing Lua code: %v", err)
}
}

func TestTaskDefinitionFields(t *testing.T) {
td := &TaskDefinition{
Name:        "test_task",
Description: "Test description",
Version:     "1.0.0",
Tags:        []string{"test", "example"},
Category:    "testing",
Workdir:     "/tmp",
}

if td.Name != "test_task" {
t.Errorf("Expected name 'test_task', got %s", td.Name)
}

if td.Description != "Test description" {
t.Errorf("Expected description 'Test description', got %s", td.Description)
}

if td.Workdir != "/tmp" {
t.Errorf("Expected workdir '/tmp', got %s", td.Workdir)
}

if len(td.Tags) != 2 {
t.Errorf("Expected 2 tags, got %d", len(td.Tags))
}
}

func TestDependencyStruct(t *testing.T) {
dep := Dependency{
Name:     "dep_task",
Optional: false,
}

if dep.Name != "dep_task" {
t.Errorf("Expected name 'dep_task', got %s", dep.Name)
}

if dep.Optional {
t.Error("Expected optional to be false")
}
}

func TestRetryConfig(t *testing.T) {
retry := RetryConfig{
MaxAttempts: 3,
Delay:       1000,
}

if retry.MaxAttempts != 3 {
t.Errorf("Expected 3 max attempts, got %d", retry.MaxAttempts)
}

	if retry.Delay != 1000 {
		t.Errorf("Expected delay 1000, got %d", retry.Delay)
	}
}

func TestTaskDefinitionWithMultipleOptions(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("test-task", {
			desc = "A test task with multiple options",
			run_if = function() return true end,
			delegate_to = "agent1",
			retry = {max_attempts = 5, delay = 2000},
			timeout = 3600,
			async = true,
			env = {KEY1 = "value1", KEY2 = "value2"}
		}, function()
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse task: %v", err)
	}
}

func TestTaskDependencies(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("task1", function()
			return {status = "ok"}
		end)

		task("task2", {
			depends_on = {"task1"}
		}, function()
			return {status = "ok"}
		end)

		task("task3", {
			depends_on = {"task1", "task2"}
		}, function()
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse tasks with dependencies: %v", err)
	}
}

func TestTaskWithArtifacts(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("build-task", {
			artifacts = {
				{name = "binary", path = "/tmp/output"},
				{name = "logs", path = "/var/log/build.log"}
			}
		}, function()
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse task with artifacts: %v", err)
	}
}

func TestTaskWithVariables(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("task-with-vars", {
			values = {
				version = "1.0.0",
				environment = "production",
				replicas = 3
			}
		}, function(ctx)
			assert(ctx.values.version == "1.0.0", "version should be 1.0.0")
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse task with variables: %v", err)
	}
}

func TestTaskWithTags(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("tagged-task", {
			tags = {"build", "deploy", "production"}
		}, function()
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse task with tags: %v", err)
	}
}

func TestTaskErrorHandling(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("error-task", {
			on_error = function(err)
				-- Handle error
				print("Error occurred: " .. tostring(err))
			end
		}, function()
			error("something went wrong")
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse task with error handler: %v", err)
	}
}

func TestNestedTaskDefinitions(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("parent-task", function()
			task("child-task-1", function()
				return {status = "ok"}
			end)
			
			task("child-task-2", function()
				return {status = "ok"}
			end)
			
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse nested tasks: %v", err)
	}
}

func TestTaskWithConsumesProduces(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("producer", {
			produces = {"artifact1", "artifact2"}
		}, function()
			return {
				artifacts = {artifact1 = "data1", artifact2 = "data2"}
			}
		end)

		task("consumer", {
			consumes = {"artifact1"}
		}, function(ctx)
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse tasks with consumes/produces: %v", err)
	}
}

func TestTaskWithNextIfFail(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("risky-task", {
			next_if_fail = "fallback-task"
		}, function()
			error("task failed")
		end)

		task("fallback-task", function()
			return {status = "recovered"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse task with next_if_fail: %v", err)
	}
}

func TestTaskWithWorkdir(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		task("workdir-task", {
			workdir = "/tmp/workspace"
		}, function()
			return {status = "ok"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse task with workdir: %v", err)
	}
}

func TestMultipleTaskGroups(t *testing.T) {
	L, _ := setupTestDSL(t)
	defer L.Close()

	script := `
		-- Group 1: Build tasks
		task("compile", function()
			return {status = "compiled"}
		end)

		task("test", {depends_on = {"compile"}}, function()
			return {status = "tested"}
		end)

		-- Group 2: Deploy tasks
		task("deploy-staging", {depends_on = {"test"}}, function()
			return {status = "deployed to staging"}
		end)

		task("deploy-production", {depends_on = {"deploy-staging"}}, function()
			return {status = "deployed to production"}
		end)

		-- Group 3: Cleanup tasks
		task("cleanup", function()
			return {status = "cleaned"}
		end)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to parse multiple task groups: %v", err)
	}
}
