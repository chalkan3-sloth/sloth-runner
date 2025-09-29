# 🦥 Getting Started Tutorial - Modern DSL

Welcome to Sloth Runner! This guide will walk you through creating and running your first set of tasks using the **Modern DSL**.

## Prerequisites

Before you begin, make sure you have:
1.  Go (version 1.21+) installed on your system.
2.  The `sloth-runner` executable installed. If not, follow the installation instructions in the main [README.md](../README.md).

## Step 1: Create Your First Task File

Let's create a simple Lua file named `my_tasks.lua`. This file will define our tasks using the **Modern DSL**.

```lua
-- my_tasks.lua - Modern DSL

-- Define your first task using the fluent API
local hello_task = task("say_hello")
    :description("Prints a friendly greeting with Modern DSL")
    :command(function()
        log.info("🦥 Hello from Sloth Runner Modern DSL!")
        return true, "echo 'Hello from Sloth Runner Modern DSL! 🦥'", {
            greeting = "Hello Modern DSL",
            timestamp = os.time()
        }
    end)
    :timeout("30s")
    :build()

-- Define the workflow
workflow.define("hello_world_workflow", {
    description = "A simple workflow to say hello - Modern DSL",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"hello-world", "tutorial", "modern-dsl"}
    },
    
    tasks = { hello_task },
    
    config = {
        timeout = "5m"
    },
    
    on_start = function()
        log.info("🚀 Starting hello world workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Hello world workflow completed!")
        end
        return true
    end
})
```

This defines a workflow with a single task using the **Modern DSL** fluent API. The task executes a simple greeting command and returns structured output.

## Step 2: Run Your Task

Now, let's run the task using the `sloth-runner` CLI. Open your terminal in the same directory where you saved `my_tasks.lua` and run:

```bash
sloth-runner run -f my_tasks.lua
```

You should see the structured output and logging from the Modern DSL task execution!

## Step 3: Add a Dependent Task

Let's make it more interesting by adding a second task that depends on the first one. Modify `my_tasks.lua`:

```lua
-- my_tasks.lua - Modern DSL with Dependencies

-- First task with output
local hello_task = task("say_hello")
    :description("Prints a friendly greeting with Modern DSL")
    :command(function()
        log.info("🦥 Executing hello task...")
        return true, "echo 'Hello from Sloth Runner Modern DSL! 🦥'", { 
            message = "Hello World",
            source = "modern_dsl"
        }
    end)
    :timeout("30s")
    :build()

-- Second task that depends on the first
local show_message_task = task("show_message")
    :description("Shows the message from the first task")
    :depends_on({"say_hello"}) -- This creates the dependency
    :command(function(params, deps)
        -- The output from 'say_hello' is available in deps!
        local received_message = deps.say_hello.message
        local source = deps.say_hello.source
        
        log.info("📩 Received message: " .. received_message .. " from " .. source)
        
        local result = exec.run("echo 'The first task said: " .. received_message .. "'")
        return result.success, result.stdout, { 
            confirmation = "Message received!",
            processed_message = received_message
        }
    end)
    :build()

-- Updated workflow with both tasks
workflow.define("hello_workflow_with_dependencies", {
    description = "Workflow demonstrating task dependencies - Modern DSL",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"dependencies", "tutorial", "modern-dsl"}
    },
    
    tasks = { hello_task, show_message_task },
    
    config = {
        timeout = "10m",
        max_parallel_tasks = 2
    },
    
    on_start = function()
        log.info("🚀 Starting workflow with dependencies...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Workflow with dependencies completed!")
            log.info("📊 Results summary:")
            for task_name, result in pairs(results) do
                log.info("  " .. task_name .. ": " .. (result.success and "✅" or "❌"))
            end
        end
        return true
    end
})
```

**Changes:**
-   Both tasks now use the **Modern DSL** fluent API with `:build()`
-   The `hello_task` returns structured output with `message` and `source` fields
-   The `show_message_task` uses `:depends_on({"say_hello"})` to create the dependency
-   Dependencies are accessed through the `deps` parameter in Modern DSL
-   Enhanced logging and error handling with the `exec` module
-   Workflow definition includes both tasks with proper metadata

## Step 4: Run the Dependent Task

Now, let's run only the final task, `show_message`. Sloth Runner will automatically figure out that it needs to run `say_hello` first.

```bash
sloth-runner run -f my_tasks.lua -t show_message
```

You will see both tasks execute in the correct order, with enhanced Modern DSL logging and structured output.

## Step 5: Advanced Features

Let's add some advanced Modern DSL features to make our workflow more robust:

```lua
-- Advanced Modern DSL example
local advanced_task = task("advanced_processing")
    :description("Demonstrates advanced Modern DSL features")
    :depends_on({"show_message"})
    :command(function(params, deps)
        log.info("🔧 Processing with advanced features...")
        
        -- Use circuit breaker for resilience
        return circuit.protect("external_service", function()
            -- Simulate processing
            return true, "Advanced processing completed", {
                processed_count = 42,
                performance_metrics = {
                    duration = "2.5s",
                    memory_usage = "64MB"
                }
            }
        end)
    end)
    :retries(3, "exponential")
    :timeout("2m")
    :on_success(function(params, output)
        log.info("🎉 Advanced processing succeeded! Processed " .. output.processed_count .. " items")
    end)
    :on_failure(function(params, error)
        log.error("❌ Advanced processing failed: " .. error)
    end)
    :build()

-- Add to workflow
workflow.define("advanced_workflow", {
    description = "Advanced workflow with Modern DSL features",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"advanced", "circuit-breaker", "retry", "modern-dsl"},
        complexity = "intermediate"
    },
    
    tasks = { hello_task, show_message_task, advanced_task },
    
    config = {
        timeout = "15m",
        retry_policy = "exponential",
        max_parallel_tasks = 3,
        circuit_breaker = {
            failure_threshold = 5,
            recovery_timeout = "30s"
        }
    }
})
```

## What's Next?

Congratulations! You've successfully created and run a Modern DSL task pipeline with advanced features.

### Learn More:
-   📚 **[Modern DSL Introduction](modern-dsl/introduction.md)** - Complete Modern DSL guide
-   🎯 **[Task Definition API](modern-dsl/task-api.md)** - Full task builder reference  
-   📋 **[Workflow Definition](modern-dsl/workflow-api.md)** - Workflow configuration guide
-   🔧 **[Lua API Reference](LUA_API.md)** - Built-in modules (`exec`, `fs`, `net`, etc.)
-   📝 **[Examples](../examples/)** - Modern DSL examples from basic to advanced
-   🎨 **[Best Practices](modern-dsl/best-practices.md)** - Modern DSL patterns and guidelines

### Next Steps:
1. **Explore Examples**: Browse the `/examples` directory for real-world Modern DSL workflows
2. **Try Built-in Modules**: Experiment with `fs`, `net`, `data`, `state` modules in your tasks
3. **Add Error Handling**: Implement retry strategies and circuit breakers
4. **Build Complex Workflows**: Create multi-stage CI/CD pipelines with the Modern DSL
