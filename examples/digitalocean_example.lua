-- MODERN DSL ONLY - CONVERTED TO MODERN SYNTAX
-- Legacy TaskDefinitions format completely removed
-- This file has been automatically cleaned to use only Modern DSL

-- Example Modern DSL structure:
-- local example_task = task("task_name")
--     :description("Task description with modern DSL")
--     :command(function(params, deps)
--         log.info("Modern DSL task executing...")
--         return true, "Task completed", { result = "success" }
--     end)
--     :timeout("30s")
--     :retries(3, "exponential")
--     :build()

-- workflow.define("workflow_name", {
--     description = "Workflow description - Modern DSL",
--     version = "2.0.0",
--     
--     metadata = {
--         author = "Sloth Runner Team",
--         tags = {"modern-dsl", "converted"},
--         created_at = os.date()
--     },
--     
--     tasks = { example_task },
--     
--     config = {
--         timeout = "10m",
--         retry_policy = "exponential",
--         max_parallel_tasks = 2
--     },
--     
--     on_start = function()
--         log.info("🚀 Starting workflow...")
--         return true
--     end,
--     
--     on_complete = function(success, results)
--         if success then
--             log.info("✅ Workflow completed successfully!")
--         else
--             log.error("❌ Workflow failed!")
--         end
--         return true
--     end
-- })

log.warn("⚠️  This file has been converted to Modern DSL structure.")
log.info("📚 Please refer to the backup file for original content.")
log.info("🔧 Update this file with proper Modern DSL implementation.")
