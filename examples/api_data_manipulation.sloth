-- MODERN DSL ONLY - Workflow description - Modern DSL
-- Converted from legacy Modern DSL format
-- Category: general, Complexity: basic

-- Main task using Modern DSL
local main_task = task("api_data_manipulation_task")
    :description("Workflow description - Modern DSL")
    :command(function(params, deps)
        log.info("🚀 Executing api_data_manipulation with Modern DSL...")
        
        -- TODO: Replace with actual implementation from original file
        -- Original logic should be migrated here
        
        return true, "Task completed successfully", {
            task_name = "api_data_manipulation",
            execution_time = os.time(),
            status = "success"
        }
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :on_success(function(params, output)
        log.info("✅ api_data_manipulation task completed successfully")
    end)
    :on_failure(function(params, error)
        log.error("❌ api_data_manipulation task failed: " .. error)
    end)
    :build()

-- Additional tasks can be added here following the same pattern
-- local secondary_task = task("api_data_manipulation_secondary")
--     :description("Secondary task for api_data_manipulation")
--     :depends_on({"api_data_manipulation_task"})
--     :command(function(params, deps)
--         -- Secondary logic here
--         return true, "Secondary task completed", {}
--     end)
--     :build()

-- Modern Workflow Definition
workflow.define("api_data_manipulation_workflow", {
    description = "Workflow description - Modern DSL - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        category = "general",
        complexity = "basic",
        tags = {"api_data_manipulation", "modern-dsl", "general"},
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
        log.info("🚀 Starting api_data_manipulation workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ api_data_manipulation workflow completed successfully!")
        else
            log.error("❌ api_data_manipulation workflow failed!")
        end
        return true
    end
})

-- Migration Note:
-- This file has been converted from legacy Modern DSL format to Modern DSL
-- TODO: Review and implement the original logic in the Modern DSL structure above
-- Original backup saved as api_data_manipulation.lua.backup
