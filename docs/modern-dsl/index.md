# 🎨 Modern DSL - Sloth Runner

Welcome to the Sloth Runner Modern DSL documentation!

## Overview

Sloth Runner uses a modern, expressive DSL (Domain Specific Language) based on Lua that makes it easy to define workflows, tasks, and orchestrate complex operations.

## Key Features

- 🔄 **Chainable API** - Fluent, readable syntax
- 🎯 **Type Safety** - Clear error messages
- 🧩 **Modular Design** - Reusable components
- 📦 **Rich Standard Library** - Built-in modules for common tasks
- 🌐 **Distributed Execution** - Native support for distributed workflows
- 💾 **State Management** - Persistent state across runs
- 📊 **Stack Management** - Pulumi-style stack support

## Quick Example

```lua
-- Define a simple task
local build_task = task("build")
    :description("Build the application")
    :command(function(params, deps)
        local exec = require("exec")
        local result = exec.run("go build -o app ./cmd")
        return result.success, result.stdout
    end)
    :timeout("5m")
    :retries(3)
    :build()

-- Create a workflow
workflow.define("my_workflow", {
    description = "My first workflow",
    version = "1.0.0",
    tasks = { build_task }
})
```

## Documentation Sections

### 📖 Getting Started

- [Introduction](./introduction.md) - Start here to learn the basics
- [Best Practices](./best-practices.md) - Learn how to write effective workflows
- [Reference Guide](./reference-guide.md) - Complete API reference

### 🔧 Core Concepts

#### Tasks
Tasks are the building blocks of workflows. They define individual units of work.

```lua
local my_task = task("task_name")
    :description("What this task does")
    :command(function(params, deps)
        -- Your code here
        return true  -- success
    end)
    :build()
```

#### Workflows
Workflows orchestrate multiple tasks with dependencies.

```lua
workflow.define("workflow_name", {
    description = "Workflow description",
    tasks = { task1, task2, task3 },
    on_success = function(results)
        print("Success!")
    end
})
```

#### Dependencies
Tasks can depend on other tasks:

```lua
local test_task = task("test")
    :depends_on({"build"})  -- Runs after build
    :command(function(params, deps)
        local build_result = deps.build
        -- Use build result
    end)
    :build()
```

### 📦 Built-in Modules

Sloth Runner provides a rich set of built-in modules:

- **exec** - Execute commands
- **fs** - File system operations
- **net** - Network operations
- **log** - Logging
- **state** - State management
- **metrics** - Metrics collection

[See all modules](../modules/index.md)

### 🎯 Common Patterns

#### Error Handling
```lua
:command(function()
    local success, error = pcall(function()
        -- Your code
    end)
    return success, error
end)
```

#### Conditional Execution
```lua
:condition(function(params)
    return params.environment == "production"
end)
```

#### Callbacks
```lua
:on_success(function(params, output)
    log.info("Task succeeded!")
end)
:on_failure(function(params, error)
    log.error("Task failed: " .. error)
end)
```

## Advanced Features

### 🗂️ Stack Management

```lua
-- Use stacks for environment isolation
sloth-runner stack run -f workflow.sloth --stack production
```

### 🌐 Distributed Execution

```lua
local remote_task = task("remote_work")
    :agent("worker-01")
    :run_on("remote_cluster")
    :command(function()
        -- Runs on remote agent
    end)
    :build()
```

### 📊 Output Formats

```bash
# Enhanced output with emojis
sloth-runner run -f workflow.sloth --output enhanced

# JSON output for automation
sloth-runner run -f workflow.sloth --output json

# Modern styled output
sloth-runner run -f workflow.sloth --output modern
```

## Examples

### CI/CD Pipeline
See [CI/CD Example](../en/examples/cicd.md)

### Infrastructure as Code
See [IaC Example](../en/examples/iac.md)

### Multi-Cloud Deployment
See [Multi-Cloud Example](../en/examples/multi-cloud.md)

## Learn More

- [Core Concepts](../core-concepts.md)
- [Advanced Features](../advanced-features.md)
- [Examples](../EXAMPLES.md)
- [Lua API Reference](../LUA_API.md)

## Community & Support

- 📚 [Documentation Home](../index.md)
- 🐛 [Report Issues](https://github.com/chalkan3-sloth/sloth-runner/issues)
- 💬 [Discussions](https://github.com/chalkan3-sloth/sloth-runner/discussions)

---

**Ready to get started?** Check out the [Introduction](./introduction.md)!
