# 🦥 Sloth Runner - Modern DSL & Lua API Reference ⚙️

This document provides a comprehensive reference for both the **Modern DSL** (Domain Specific Language) and the enhanced Lua modules in `sloth-runner`. The Modern DSL provides a fluent, intuitive way to define workflows while maintaining full access to powerful Lua modules.

---

## 🎯 Modern DSL API Reference

### Task Definition API

The Modern DSL introduces a fluent, chainable API for defining tasks with enhanced features.

#### `task(name)` - Task Builder

Creates a new task builder with the specified name.

```lua
local my_task = task("task_name")
    :description("Task description")
    :command(function(params, deps) ... end)
    :build()
```

#### Task Builder Methods

##### `:description(desc)`
Sets the task description.
- **Parameters:** `desc` (string) - Human-readable task description
- **Returns:** Task builder for chaining

##### `:command(cmd)`
Defines the task command - can be a string or function.
- **Parameters:** 
  - `cmd` (string|function) - Shell command or Lua function
  - For functions: `function(params, deps)` where:
    - `params` (table) - Task parameters and configuration
    - `deps` (table) - Outputs from dependent tasks
- **Returns:** Task builder for chaining

```lua
-- String command
:command("echo 'Hello World'")

-- Function command with enhanced error handling
:command(function(params, deps)
    log.info("Executing modern task...")
    
    -- Enhanced error handling
    local result, err = exec.run("complex-command")
    if not result.success then
        return false, "Command failed: " .. err
    end
    
    return true, "Task completed", {
        output = result.stdout,
        timestamp = os.time()
    }
end)
```

##### `:depends_on(dependencies)`
Specifies task dependencies.
- **Parameters:** `dependencies` (array) - List of task names this task depends on
- **Returns:** Task builder for chaining

```lua
:depends_on({"build", "test"})
```

##### `:timeout(duration)`
Sets task timeout.
- **Parameters:** `duration` (string) - Timeout duration (e.g., "30s", "5m", "1h")
- **Returns:** Task builder for chaining

##### `:retries(count, strategy)`
Configures retry behavior.
- **Parameters:** 
  - `count` (number) - Number of retry attempts
  - `strategy` (string) - Retry strategy: "linear", "exponential", "fixed"
- **Returns:** Task builder for chaining

```lua
:retries(3, "exponential")  -- 3 retries with exponential backoff
```

##### `:artifacts(list)`
Specifies artifacts produced by this task.
- **Parameters:** `list` (array) - List of artifact file/directory names
- **Returns:** Task builder for chaining

##### `:consumes(list)`
Specifies artifacts consumed by this task.
- **Parameters:** `list` (array) - List of artifact names to consume
- **Returns:** Task builder for chaining

##### `:condition(predicate)`
Sets conditional execution logic.
- **Parameters:** `predicate` (function|string) - Condition to evaluate
- **Returns:** Task builder for chaining

```lua
:condition(when("env.DEPLOY == 'production'"))
```

##### `:on_success(callback)`
Defines success callback.
- **Parameters:** `callback` (function) - `function(params, output)`
- **Returns:** Task builder for chaining

##### `:on_failure(callback)`
Defines failure callback.
- **Parameters:** `callback` (function) - `function(params, error)`
- **Returns:** Task builder for chaining

##### `:async(enabled)`
Enables/disables asynchronous execution.
- **Parameters:** `enabled` (boolean) - Whether to run asynchronously
- **Returns:** Task builder for chaining

##### `:build()`
Finalizes and returns the constructed task.
- **Returns:** Task object ready for use in workflows

### Workflow Definition API

#### `workflow.define(name, config)`

Defines a complete workflow with modern configuration.

```lua
workflow.define("workflow_name", {
    description = "Workflow description",
    version = "2.0.0",
    
    metadata = {
        author = "Developer Name",
        tags = {"tag1", "tag2"},
        created_at = os.date(),
        category = "ci/cd"
    },
    
    tasks = { task1, task2, task3 },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 4,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("Starting workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("Workflow completed successfully!")
        else
            log.error("Workflow failed!")
        end
        return true
    end,
    
    on_failure = function(task_name, error)
        log.error("Task " .. task_name .. " failed: " .. error)
        return true  -- Continue workflow
    end
})
```

#### Workflow Configuration Options

- **`description`** (string) - Workflow description
- **`version`** (string) - Workflow version
- **`metadata`** (table) - Rich metadata including author, tags, etc.
- **`tasks`** (array) - List of task objects
- **`config`** (table) - Workflow-level configuration
- **`on_start`** (function) - Executed before workflow starts
- **`on_complete`** (function) - Executed after workflow completes
- **`on_failure`** (function) - Executed when tasks fail

---

## ⚡ Enhanced Lua Modules

### `exec` Module - Enhanced Command Execution

The `exec` module provides advanced command execution with modern error handling.

#### `exec.run(command, options)`

Executes a shell command with enhanced options and error handling.

```lua
-- Simple execution
local result = exec.run("echo 'Hello World'")

-- Advanced execution with options
local result = exec.run("long-running-command", {
    timeout = "5m",
    env = {
        NODE_ENV = "production",
        API_KEY = "secret"
    },
    workdir = "/path/to/project",
    retry = {
        count = 3,
        strategy = "exponential"
    }
})

-- Check results
if result.success then
    log.info("Command output: " .. result.stdout)
else
    log.error("Command failed: " .. result.stderr)
end
```

**Parameters:**
- `command` (string) - Command to execute
- `options` (table, optional) - Execution options:
  - `timeout` (string) - Command timeout
  - `env` (table) - Environment variables
  - `workdir` (string) - Working directory
  - `retry` (table) - Retry configuration

**Returns:**
- `success` (boolean) - Whether command succeeded
- `stdout` (string) - Standard output
- `stderr` (string) - Standard error
- `exit_code` (number) - Exit code
- `duration` (number) - Execution time in seconds

### `async` Module - Modern Async Operations

New module for modern asynchronous patterns.

#### `async.parallel(tasks, options)`

Execute multiple tasks in parallel with modern configuration.

```lua
local results = async.parallel({
    frontend = function()
        return exec.run("npm run build:frontend")
    end,
    backend = function()
        return exec.run("go build ./cmd/server")
    end,
    tests = function()
        return exec.run("npm test")
    end
}, {
    max_workers = 3,
    timeout = "10m",
    fail_fast = true
})

-- Process results
for name, result in pairs(results) do
    if result.success then
        log.info(name .. " completed successfully")
    else
        log.error(name .. " failed: " .. result.error)
    end
end
```

#### `async.timeout(duration, task)`

Execute a task with timeout handling.

```lua
local result = async.timeout("5m", function()
    return exec.run("long-running-task")
end)
```

### `circuit` Module - Circuit Breaker Patterns

New module for resilience patterns.

#### `circuit.protect(name, task, options)`

Protect a task with circuit breaker pattern.

```lua
local result = circuit.protect("external_api", function()
    return net.http_get("https://api.example.com/data")
end, {
    failure_threshold = 5,
    timeout = "30s",
    reset_timeout = "1m"
})
```

### `perf` Module - Performance Monitoring

New module for performance tracking.

#### `perf.measure(name, task)`

Measure task performance and collect metrics.

```lua
local result, duration = perf.measure("database_query", function()
    return database.query("SELECT * FROM users")
end)

log.info("Query completed in " .. duration .. "ms")
```

### `utils` Module - Enhanced Utilities

New utilities module for common operations.

#### `utils.config(name, environment)`

Load configuration with environment support.

```lua
local config = utils.config("app_config", "production")
local db_host = config.database.host
```

#### `utils.secret(name)`

Secure secret retrieval.

```lua
local api_key = utils.secret("external_api_key")
```

#### `utils.template(template, variables)`

Template rendering with variable substitution.

```lua
local rendered = utils.template("Hello {% raw %}{{name}}{% endraw %}", {
    name = "World"
})
```

### `validate` Module - Input Validation

New validation module for type checking and input validation.

#### `validate.required(value, name)`

Validate required fields.

```lua
validate.required(params.api_key, "api_key")
```

#### `validate.type(value, expected_type, name)`

Validate value types.

```lua
validate.type(params.timeout, "string", "timeout")
validate.type(params.retries, "number", "retries")
```

---

## 🔄 Enhanced Existing Modules

### `fs` Module - Enhanced File Operations

Enhanced with better error handling and metadata support.

#### `fs.read(path, options)`

```lua
-- Simple read
local content, err = fs.read("config.yaml")

-- Enhanced read with options
local content, err = fs.read("large-file.txt", {
    encoding = "utf-8",
    max_size = "10MB"
})
```

#### `fs.write(path, content, options)`

```lua
-- Enhanced write with metadata
fs.write("output.json", data.to_json(result), {
    mode = 0644,
    backup = true,
    atomic = true
})
```

### `net` Module - Enhanced HTTP Client

Enhanced with retry, circuit breakers, and better error handling.

#### `net.http_get(url, options)`

```lua
local response = net.http_get("https://api.example.com/data", {
    headers = {
        ["Authorization"] = "Bearer " .. token,
        ["Content-Type"] = "application/json"
    },
    timeout = "30s",
    retry = {
        count = 3,
        strategy = "exponential"
    },
    circuit_breaker = {
        name = "external_api",
        failure_threshold = 5
    }
})
```

### `state` Module - Advanced State Management

Enhanced with clustering, TTL, and atomic operations.

#### `state.set_with_ttl(key, value, ttl)`

Set value with time-to-live.

```lua
state.set_with_ttl("session_token", token, 3600)  -- 1 hour TTL
```

#### `state.atomic_update(key, update_function)`

Atomic state updates.

```lua
state.atomic_update("counter", function(current_value)
    return (current_value or 0) + 1
end)
```

---

## 📊 Legacy Format Support

The Modern DSL maintains **100% backward compatibility** with the legacy `Modern DSLs` format:

```lua
-- Legacy format (still fully supported)
Modern DSLs = {
    my_workflow = {
        description = "Legacy workflow",
        tasks = {
            {
                name = "my_task",
                command = "echo 'Hello'",
                depends_on = "other_task",
                timeout = "30s"
            }
        }
    }
}
```

This allows for gradual migration and ensures existing scripts continue to work without modification.

---

## 🎯 Migration Examples

### Converting Legacy to Modern DSL

**Before (Legacy):**
```lua
Modern DSLs = {
    build_pipeline = {
        description = "Build and test pipeline",
        tasks = {
            {
                name = "build",
                command = "go build ./...",
                timeout = "5m",
                retries = 2
            },
            {
                name = "test",
                command = "go test ./...",
                depends_on = "build",
                timeout = "10m"
            }
        }
    }
}
```

**After (Modern DSL):**
```lua
local build_task = task("build")
    :description("Build the application")
    :command("go build ./...")
    :timeout("5m")
    :retries(2, "exponential")
    :build()

local test_task = task("test")
    :description("Run tests")
    :command("go test ./...")
    :depends_on({"build"})
    :timeout("10m")
    :build()

workflow.define("build_pipeline", {
    description = "Build and test pipeline - Modern DSL",
    version = "2.0.0",
    tasks = { build_task, test_task },
    config = {
        timeout = "20m",
        retry_policy = "exponential"
    }
})
```

This enhanced API provides better error handling, more features, and improved maintainability while preserving all existing functionality.

### `fs.write(path, content)`

Writes a string to a file, overwriting it if it exists.

-   **Parameters:**
    -   `path` (string): The path to the file.
    -   `content` (string): The content to write.
-   **Returns:**
    -   `err` (string or nil): An error message on failure, otherwise `nil`.

### `fs.append(path, content)`

Appends a string to the end of a file, creating it if it doesn't exist.

-   **Parameters:**
    -   `path` (string): The path to the file.
    -   `content` (string): The content to append.
-   **Returns:**
    -   `err` (string or nil): An error message on failure, otherwise `nil`.

### `fs.exists(path)`

Checks if a file or directory exists at the given path.

-   **Parameters:**
    -   `path` (string): The path to check.
-   **Returns:**
    -   `exists` (boolean): `true` if the path exists, `false` otherwise.

### `fs.mkdir(path)`

Creates a directory, including any necessary parent directories.

-   **Parameters:**
    -   `path` (string): The directory path to create.
-   **Returns:**
    -   `err` (string or nil): An error message on failure, otherwise `nil`.

### `fs.rm(path)`

Removes a file or an empty directory.

-   **Parameters:**
    -   `path` (string): The path to remove.
-   **Returns:**
    -   `err` (string or nil): An error message on failure, otherwise `nil`.

### `fs.rm_r(path)`

Recursively removes a directory and all its contents.

-   **Parameters:**
    -   `path` (string): The path to the directory to remove.
-   **Returns:**
    -   `err` (string or nil): An error message on failure, otherwise `nil`.

### `fs.ls(path)`

Lists the names of files and directories inside a given path.

-   **Parameters:**
    -   `path` (string): The path to the directory.
-   **Returns:**
    -   `files` (table or nil): A Lua table (array) of file and directory names, or `nil` on error.
    -   `err` (string or nil): An error message on failure, otherwise `nil`.

**Example:**

```lua
local dir = "/tmp/sloth-test"
fs.mkdir(dir)
fs.write(dir .. "/hello.txt", "Hello from Sloth! 🦥")
local files, err = fs.ls(dir)
if err then
    log.error("Could not list files: " .. err)
else
    log.info("Files in " .. dir .. ": " .. data.to_json(files))
end
fs.rm_r(dir)
```

---

## `net` Module

The `net` module provides networking utilities.

### `net.http_get(url)`

Performs an HTTP GET request.

-   **Parameters:**
    -   `url` (string): The URL to request.
-   **Returns:**
    -   `body` (string or nil): The response body.
    -   `status_code` (number): The HTTP status code (e.g., `200`).
    -   `headers` (table or nil): A Lua table of response headers.
    -   `err` (string or nil): An error message on failure.

### `net.http_post(url, body, [headers])`

Performs an HTTP POST request.

-   **Parameters:**
    -   `url` (string): The URL to post to.
    -   `body` (string): The request body.
    -   `headers` (table, optional): A Lua table of request headers.
-   **Returns:**
    -   `body` (string or nil): The response body.
    -   `status_code` (number): The HTTP status code.
    -   `headers` (table or nil): A Lua table of response headers.
    -   `err` (string or nil): An error message on failure.

### `net.download(url, destination_path)`

Downloads a file from a URL to a local path.

-   **Parameters:**
    -   `url` (string): The URL of the file to download.
    -   `destination_path` (string): The local path to save the file.
-   **Returns:**
    -   `err` (string or nil): An error message on failure.

**Example:**

```lua
log.info("Fetching a random cat fact...")
local body, status, _, err = net.http_get("https://catfact.ninja/fact")
if err or status ~= 200 then
    log.error("Failed to fetch cat fact: " .. (err or "status " .. status))
else
    local fact_data, json_err = data.parse_json(body)
    if json_err then
        log.error("Could not parse cat fact JSON: " .. json_err)
    else
        log.info("🐱 Cat Fact: " .. fact_data.fact)
    end
end
```

---

## `data` Module

The `data` module provides functions for data serialization and deserialization.

### `data.to_json(table)`

Converts a Lua table to a JSON string.

-   **Parameters:**
    -   `table` (table): The Lua table to convert.
-   **Returns:**
    -   `json_string` (string or nil): The resulting JSON string.
    -   `err` (string or nil): An error message on failure.

### `data.parse_json(json_string)`

Parses a JSON string into a Lua table.

-   **Parameters:**
    -   `json_string` (string): The JSON string to parse.
-   **Returns:**
    -   `table` (table or nil): The resulting Lua table.
    -   `err` (string or nil): An error message on failure.

### `data.to_yaml(table)`

Converts a Lua table to a YAML string.

-   **Parameters:**
    -   `table` (table): The Lua table to convert.
-   **Returns:**
    -   `yaml_string` (string or nil): The resulting YAML string.
    -   `err` (string or nil): An error message on failure.

### `data.parse_yaml(yaml_string)`

Parses a YAML string into a Lua table.

-   **Parameters:**
    -   `yaml_string` (string): The YAML string to parse.
-   **Returns:**
    -   `table` (table or nil): The resulting Lua table.
    -   `err` (string or nil): An error message on failure.

---

## `log` Module

The `log` module provides simple logging functions.

### `log.info(message)`
### `log.warn(message)`
### `log.error(message)`
### `log.debug(message)`

-   **Parameters:**
    -   `message` (string): The message to log.

**Example:**

```lua
log.info("Starting the task.")
log.warn("This is a warning.")
log.error("Something went wrong!")
log.debug("Here is some debug info.")
```

---

## `salt` Module

The `salt` module allows for direct execution of SaltStack commands.

### `salt.cmd(command_type, [arg1, arg2, ...])`

Executes a SaltStack command.

-   **Parameters:**
    -   `command_type` (string): The type of command, either `"salt"` or `"salt-call"`.
    -   `arg...` (string, optional): A variable number of string arguments for the command.
-   **Returns:**
    -   `stdout` (string): The standard output of the command.
    -   `stderr` (string): The standard error output of the command.
    -   `err` (string or nil): An error message if the command fails, otherwise `nil`.

**Example:**

```lua
-- Ping all minions
local stdout, stderr, err = salt.cmd("salt", "*", "test.ping")
if err then
    log.error("Salt command failed: " .. stderr)
else
    log.info("Salt ping result:\n" .. stdout)
end
```

