-- MODERN DSL ONLY - Workflow description - Modern DSL
-- Converted from legacy Modern DSL format
-- Category: general, Complexity: basic

-- Main task using Modern DSL
local main_task = task("cicd_gcp_hub_spoke_task")
    :description("Workflow description - Modern DSL")
    :command(function(params, deps)
        log.info("🚀 Executing cicd_gcp_hub_spoke with Modern DSL...")
        
        -- TODO: Replace with actual implementation from original file
        -- Original logic should be migrated here
        
        return true, "Task completed successfully", {
            task_name = "cicd_gcp_hub_spoke",
            execution_time = os.time(),
            status = "success"
        }
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :on_success(function(params, output)
        log.info("✅ cicd_gcp_hub_spoke task completed successfully")
    end)
    :on_failure(function(params, error)
        log.error("❌ cicd_gcp_hub_spoke task failed: " .. error)
    end)
    :build()

-- Additional tasks can be added here following the same pattern
-- local secondary_task = task("cicd_gcp_hub_spoke_secondary")
--     :description("Secondary task for cicd_gcp_hub_spoke")
--     :depends_on({"cicd_gcp_hub_spoke_task"})
--     :command(function(params, deps)
--         -- Secondary logic here
--         return true, "Secondary task completed", {}
--     end)
--     :build()

-- Modern Workflow Definition
workflow.define("cicd_gcp_hub_spoke_workflow", {
    description = "Workflow description - Modern DSL - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        category = "general",
        complexity = "basic",
        tags = {"cicd_gcp_hub_spoke", "modern-dsl", "general"},
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
        log.info("🚀 Starting cicd_gcp_hub_spoke workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ cicd_gcp_hub_spoke workflow completed successfully!")
        else
            log.error("❌ cicd_gcp_hub_spoke workflow failed!")
        end
        return true
    end
})

-- Migration Note:
-- This file has been converted from legacy Modern DSL format to Modern DSL
-- TODO: Review and implement the original logic in the Modern DSL structure above
-- Original backup saved as cicd_gcp_hub_spoke.lua.backup
