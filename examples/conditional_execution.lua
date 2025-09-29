-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:32 -03

local check_condition_task = task("check_condition_for_run")
local conditional_task = task("conditional_task")
local cleanup_task = task("cleanup_run_condition")
local abort_check_task = task("check_abort_condition")
local final_task = task("final_task")

local check_condition_task = task("check_condition_for_run")
local conditional_task = task("conditional_task")
local cleanup_task = task("cleanup_run_condition")
local abort_check_task = task("check_abort_condition")
local final_task = task("final_task")
local check_condition_task = task("check_condition_for_run")
    :description("Creates a file for conditional checking")
    :command(function()
        log.info("🔧 Creating condition file for testing...")
        local success, output = exec.run("touch /tmp/sloth_runner_run_condition")
        if success then
            log.info("✅ Condition file created at /tmp/sloth_runner_run_condition")
            return true, "Condition file created", { file_path = "/tmp/sloth_runner_run_condition" }
        else
            log.error("❌ Failed to create condition file")
            return false, "Failed to create condition file"
        end
    end)
    :timeout("10s")
    :build()
local conditional_task = task("conditional_task")
    :description("Runs conditionally based on file existence")
    :depends_on({"check_condition_for_run"})
    :run_if("test -f /tmp/sloth_runner_run_condition")
    :command(function(params, deps)
        log.info("🎯 Conditional task executing - condition was met!")
        log.info("📁 Condition file: " .. deps.check_condition_for_run.file_path)
        
        -- Perform conditional work
        local timestamp = os.time()
        log.info("⏰ Execution timestamp: " .. timestamp)
        
        return true, "Conditional task completed", {
            executed_at = timestamp,
            reason = "condition_file_exists"
        }
    end)
    :on_success(function()
        log.info("✨ Conditional task executed successfully!")
    end)
    :on_skip(function()
        log.warn("⏭️  Conditional task skipped - condition not met")
    end)
    :build()
local cleanup_task = task("cleanup_run_condition")
    :description("Cleans up the condition file")
    :depends_on({"conditional_task"})
    :command(function(params, deps)
        log.info("🧹 Cleaning up condition file...")
        
        -- Enhanced cleanup with verification
        local file_path = "/tmp/sloth_runner_run_condition"
        
        if fs.exists(file_path) then
            local success, output = exec.run("rm -f " .. file_path)
            if success then
                log.info("✅ Condition file cleaned up successfully")
                return true, "Cleanup completed", { cleaned_file = file_path }
            else
                log.error("❌ Failed to cleanup: " .. output)
                return false, "Cleanup failed"
            end
        else
            log.info("ℹ️  Condition file already cleaned up")
            return true, "Already clean", { status = "already_clean" }
        end
    end)
    :build()
local abort_check_task = task("check_abort_condition")
    :description("Task that aborts workflow if specific condition is met")
    :command(function()
        log.info("🛡️  Checking for abort conditions...")
        
        -- Multiple abort conditions
        local abort_conditions = {
            "/tmp/sloth_runner_abort_condition",
            "/tmp/emergency_stop",
            "/tmp/maintenance_mode"
        }
        
        for _, condition_file in ipairs(abort_conditions) do
            if fs.exists(condition_file) then
                log.error("🚨 ABORT CONDITION DETECTED: " .. condition_file)
                return false, "Abort condition: " .. condition_file .. " exists"
            end
        end
        
        -- Additional environment checks
        local env_check = os.getenv("FORCE_ABORT")
        if env_check == "true" then
            log.error("🚨 ABORT: Environment variable FORCE_ABORT is set")
            return false, "Environment abort condition"
        end
        
        log.info("✅ No abort conditions detected - safe to proceed")
        return true, "All checks passed", {
            check_time = os.time(),
            conditions_checked = #abort_conditions + 1,
            status = "safe_to_proceed"
        }
    end)
    :abort_if(function()
        -- Function-based abort condition
        return fs.exists("/tmp/sloth_runner_abort_condition")
    end)
    :on_success(function()
        log.info("🛡️  Security checks passed!")
    end)
    :build()
local final_task = task("final_task")
    :description("Final task - executes only if all conditions are met")
    :depends_on({"check_abort_condition"})
    :command(function(params, deps)
        log.info("🎉 Final task executing - all conditions passed!")
        
        -- Comprehensive validation summary
        local validation_summary = {
            abort_check_status = deps.check_abort_condition.status,
            conditions_verified = deps.check_abort_condition.conditions_checked,
            execution_time = os.time(),
            workflow_completion = "success"
        }
        
        log.info("📊 Validation Summary:")
        for key, value in pairs(validation_summary) do
            log.info("  " .. key .. ": " .. tostring(value))
        end
        
        return true, "Workflow completed successfully", validation_summary
    end)
    :on_success(function(params, output)
        log.info("🏆 Workflow completed successfully!")
        log.info("📈 Final status: " .. output.workflow_completion)
    end)
    :build()

workflow.define("conditional_execution_demo", {
    description = "Conditional execution demonstration - Modern DSL Only",
    version = "2.0.0",
    
    metadata = {
        category = "demonstration",
        tags = {"conditional", "abort", "run_if", "modern-dsl"},
        author = "Sloth Runner Team",
        complexity = "intermediate"
    },
    
    tasks = {
        check_condition_task,
        conditional_task,
        cleanup_task,
        abort_check_task,
        final_task
    },
    
    config = {
        timeout = "10m",
        fail_fast = true,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🎬 Starting conditional execution workflow...")
        log.info("💡 This demonstrates: run_if, abort_if, and conditional logic")
        return true
    end,
    
    on_abort = function(reason)
        log.warn("⚠️  Workflow aborted: " .. reason)
        log.info("🧹 Performing emergency cleanup...")
        -- Emergency cleanup
        exec.run("rm -f /tmp/sloth_runner_*")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Conditional execution demo completed successfully!")
            log.info("✅ All conditional logic worked as expected")
        else
            log.error("❌ Conditional execution demo failed!")
        end
        return true
    end
})
