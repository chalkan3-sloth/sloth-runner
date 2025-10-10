# Core Concepts - Modern DSL

This document explains the fundamental concepts of `sloth-runner` using the **Modern DSL**, helping you understand how to define and orchestrate complex workflows with the new fluent API.

---

## Modern DSL Overview

The Modern DSL provides an intuitive, fluent API for defining workflows. You can use chainable methods to build tasks and define workflows declaratively.

```lua
-- my_pipeline.sloth - Modern DSL
local my_task = task("task_name")
    :description("Task description")
    :command(function(this, params)
        -- Task logic here
        return true, "Task completed"
    end)
    :build()

workflow
    .define("workflow_name")
    :description("Workflow description - Modern DSL")
    :version("1.0.0")
    :tasks({my_task})
```

---

## Task Definition with Modern DSL

Tasks are now defined using the `task()` function and fluent API methods:

### Basic Task Structure

```lua
local my_task = task("task_name")
    :description("What this task does")
    :command(function(this, params)
        -- Task logic here
        return true, "Success message", {output_data = "value"}
    end)
    :timeout("5m")
    :retries(3)
    :build()
```

### Task Builder Methods

**Core Properties:**
*   `:description(string)` - Human-readable task description
*   `:command(function|string)` - Task execution logic
*   `:timeout(string)` - Maximum execution time (e.g., "10s", "5m", "1h")
*   `:retries(number, strategy)` - Retry configuration with strategy ("exponential", "linear", "fixed")
*   `:depends_on(array)` - Array of task names this task depends on

**Advanced Features:**
*   `:async(boolean)` - Enable asynchronous execution
*   `:artifacts(array)` - Files to save after successful execution
*   `:consumes(array)` - Artifacts from other tasks to use
*   `:run_if(function|string)` - Conditional execution logic
*   `:abort_if(function|string)` - Condition to abort entire workflow

**Lifecycle Hooks:**
*   `:on_success(function)` - Execute when task succeeds
*   `:on_failure(function)` - Execute when task fails  
*   `:on_timeout(function)` - Execute when task times out
*   `:pre_hook(function)` - Execute before main command
*   `:post_hook(function)` - Execute after main command

**Example:**
```lua
-- Define a workflow with workdir management
workflow
  .define("my_workflow")
  :description("A workflow that manages its own temporary directory")
  :config({
    workdir = "/tmp/my_workflow",
    cleanup = "on_success"  -- or "always", "never"
  })
  :tasks({
    task("setup")
      :description("Setup task")
      :command(function()
        print("Setting up...")
        return true
      end)
      :build()
  })

---

## Individual Tasks

A task is a single unit of work. It's defined as a table with several available properties to control its behavior.

### Basic Properties

*   `name` (string): The unique name of the task within its group.
*   `description` (string): A brief description of what the task does.
*   `command` (string or function): The core action of the task.
    *   **As a string:** It's executed as a shell command.
    *   **As a function:** The Lua function is executed. It receives two arguments: `params` (a table of its parameters) and `deps` (a table containing the outputs of its dependencies). The function must return:
        1.  `boolean`: `true` for success, `false` for failure.
        2.  `string`: A message describing the result.
        3.  `table` (optional): A table of outputs that other tasks can depend on.

### Dependency and Execution Flow

*   `depends_on` (string or table): A list of task names that must complete successfully before this task can run.
*   `next_if_fail` (string or table): A list of task names to run *only if* this task fails. This is useful for cleanup or notification tasks.
*   `async` (boolean): If `true`, the task runs in the background, and the runner does not wait for it to complete before starting the next task in the execution order.

### Error Handling and Robustness

*   `retries` (number): The number of times to retry a task if it fails. Default is `0`.
*   `timeout` (string): A duration (e.g., `"10s"`, `"1m"`) after which the task will be terminated if it's still running.

### Conditional Execution

*   `run_if` (string or function): The task will be skipped unless this condition is met.
    *   **As a string:** A shell command. An exit code of `0` means the condition is met.
    *   **As a function:** A Lua function that returns `true` if the task should run.
*   `abort_if` (string or function): The entire workflow will be aborted if this condition is met.
    *   **As a string:** A shell command. An exit code of `0` means abort.
    *   **As a function:** A Lua function that returns `true` to abort.

### Lifecycle Hooks

*   `pre_exec` (function): A Lua function that runs *before* the main `command`.
*   `post_exec` (function): A Lua function that runs *after* the main `command` has completed successfully.

### Reusability

*   `uses` (table): Specifies a pre-defined task from another file (loaded via `import`) to use as a base. The current task definition can then override properties like `params` or `description`.
*   `params` (table): A dictionary of key-value pairs that can be passed to the task's `command` function.
*   `artifacts` (string or table): A file pattern (glob) or a list of patterns specifying which files from the task's `workdir` should be saved as artifacts after a successful run.
*   `consumes` (string or table): The name of an artifact (or a list of names) from a previous task that should be copied into this task's `workdir` before it runs.

---

## Artifact Management

Sloth-Runner allows tasks to share files with each other through an artifact mechanism. One task can "produce" one or more files as artifacts, and subsequent tasks can "consume" those artifacts.

This is useful for CI/CD pipelines where a build step might generate a binary (the artifact), which is then used by a testing or deployment step.

### How It Works

1.  **Producing Artifacts:** Add the `artifacts` key to your task definition. The value can be a single file pattern (e.g., `"report.txt"`) or a list (e.g., `{"*.log", "app.bin"}`). After the task runs successfully, the runner will find files in the task's `workdir` matching these patterns and copy them to a shared artifact storage for the pipeline.

2.  **Consuming Artifacts:** Add the `consumes` key to another task's definition (which typically `depends_on` the producer task). The value should be the filename of the artifact you want to use (e.g., `"report.txt"`). Before this task runs, the runner will copy the named artifact from the shared storage into this task's `workdir`, making it available to the `command`.

### Artifacts Example

```lua
local build_task = task("build")
    :description("Creates a binary and declares it as an artifact")
    :command(function(this, params)
        exec.run("echo 'binary_content' > app.bin")
        return true, "Binary created"
    end)
    :artifacts({"app.bin"})
    :build()

local test_task = task("test")
    :description("Consumes the binary to run tests")
    :depends_on({"build"})
    :consumes({"app.bin"})
    :command(function(this, params)
        -- At this point, 'app.bin' exists in this task's workdir
        local success, content = exec.run("cat app.bin")
        if content:find("binary_content") then
            log.info("Successfully consumed artifact!")
            return true, "Artifact validated"
        else
            return false, "Artifact content was incorrect!"
        end
    end)
    :build()

workflow
    .define("ci_pipeline")
    :description("Demonstrates the use of artifacts")
    :version("1.0.0")
    :tasks({build_task, test_task})
    :config({
        timeout = "10m"
    })
```

---

## Global Functions

`sloth-runner` provides global functions in the Lua environment to help orchestrate workflows.

### `import(path)`

Loads another sloth file and returns the value it returns. This is the primary mechanism for creating reusable task modules. The path is relative to the file calling `import`.

**Example (`reusable_tasks.sloth`):**
```lua
-- Import a module that returns task definitions
local docker_tasks = import("shared/docker.sloth")

-- Use the imported task with custom parameters
local build_app = docker_tasks.build_image("my-app")
    :description("Build my-app Docker image")
    :timeout("10m")
    :build()

workflow
    .define("main")
    :description("Main workflow using reusable tasks")
    :version("1.0.0")
    :tasks({build_app})
```

### `parallel(tasks)`

Executes a list of tasks concurrently and waits for all of them to complete.

*   `tasks` (table): A list of task tables to run in parallel.

**Example:**
```lua
command = function()
  log.info("Starting 3 tasks in parallel...")
  local results, err = parallel({
    { name = "short_task", command = "sleep 1" },
    { name = "medium_task", command = "sleep 2" },
    { name = "long_task", command = "sleep 3" }
  })
  if err then
    return false, "Parallel execution failed"
  end
  return true, "All parallel tasks finished."
end
```

### `export(table)`

Exports data from any point in a script to the CLI. When the `--return` flag is used, all exported tables are merged with the final task's output into a single JSON object.

*   `table`: A Lua table to be exported.

**Example:**
```lua
command = function()
  export({ important_value = "data from the middle of a task" })
  return true, "Task done", { final_output = "some result" }
end
```
Running with `--return` would produce:
```json
{
  "important_value": "data from the middle of a task",
  "final_output": "some result"
}
```