-- MODERN DSL ONLY - Hello World Example
-- Demonstrates basic Modern DSL task creation

-- Hello World task using Modern DSL
local hello_task = task("hello_world")
    :description("Simple hello world demonstration")
    :command(function(params)
        log.info("🌟 Hello World from Modern DSL!")
        log.info("📅 Current time: " .. os.date())
        
        return true, "echo 'Hello, Modern Sloth Runner!'", {
            message = "Hello World",
            timestamp = os.time(),
            status = "success"
        }
    end)
    :timeout("30s")
    :on_success(function(params, output)
        log.info("✅ Hello World task completed successfully!")
        log.info("💬 Message: " .. output.message)
    end)
    :build()

-- Modern Workflow Definition
workflow.define("hello_world_workflow", {
    description = "Simple Hello World - Modern DSL",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"hello-world", "beginner", "modern-dsl"},
        created_at = os.date()
    },
    
    tasks = { hello_task },
    
    config = {
        timeout = "5m",
        max_parallel_tasks = 1
    },
    
    on_start = function()
        log.info("🚀 Starting Hello World workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Hello World workflow completed!")
        else
            log.error("❌ Hello World workflow failed!")
        end
        return true
    end
})
