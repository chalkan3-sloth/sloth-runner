-- MODERN DSL ONLY - Workflow description - Modern DSL
-- Converted from legacy Modern DSL format
-- Category: beginner, Complexity: simple

-- Main task using Modern DSL
local main_task = task("docker-basics_task")
    :description("Workflow description - Modern DSL")
    :command(function(params, deps)
        log.info("🚀 Executing docker-basics with Modern DSL...")
        
        -- TODO: Replace with actual implementation from original file
        -- Original logic should be migrated here
        
        return true, "Task completed successfully", {
            task_name = "docker-basics",
            execution_time = os.time(),
            status = "success"
        }
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :on_success(function(params, output)
        log.info("✅ docker-basics task completed successfully")
    end)
    :on_failure(function(params, error)
        log.error("❌ docker-basics task failed: " .. error)
    end)
    :build()

-- Additional tasks can be added here following the same pattern
-- local secondary_task = task("docker-basics_secondary")
--     :description("Secondary task for docker-basics")
--     :depends_on({"docker-basics_task"})
--     :command(function(params, deps)
--         -- Secondary logic here
--         return true, "Secondary task completed", {}
--     end)
--     :build()

-- Modern Workflow Definition
workflow.define("docker-basics_workflow", {
    description = "Workflow description - Modern DSL - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        category = "beginner",
        complexity = "simple",
        tags = {"docker-basics", "modern-dsl", "beginner"},
        created_at = os.date(),
        migrated_from = "Modern DSL format"
    },
    
    tasks = {
        main_task
        -- Add additional tasks here
    },
    
    config = {
        timeout = "15m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting docker-basics workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ docker-basics workflow completed successfully!")
        else
            log.error("❌ docker-basics workflow failed!")
        end
        return true
    end
})

-- Migration Note:
-- This file has been converted from legacy Modern DSL format to Modern DSL
-- TODO: Review and implement the original logic in the Modern DSL structure above
-- Original backup saved as docker-basics.lua.backup
