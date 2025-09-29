-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03


-- local example_task = task("task_name")
--     :description("Task description with modern DSL")
--     :command(function(params, deps)
--         -- Enhanced task logic
--         return true, "Task completed", { result = "success" }
--     end)
--     :timeout("30s")
--     :build()

-- workflow.define("workflow_name", {
--     description = "Workflow description - Modern DSL",
--     version = "2.0.0",
--     tasks = { example_task },
--     config = { timeout = "10m" }
-- })

-- Maintain backward compatibility with legacy format
TaskDefinitions = {
    conditional_functions_workflow = {
        description = "A workflow to demonstrate conditional execution with Lua functions.",
        tasks = {
            {
                name = "setup_task",
                description = "This task provides output for the conditional task.",
                command = function()
                    return true, "Setup complete", { should_run = true }
                end
            },
            {
                name = "conditional_task_with_function",
                description = "This task only runs if the run_if function returns true.",
                depends_on = "setup_task",
                run_if = function(params, deps)
                    log.info("Checking run_if condition for conditional_task_with_function...")
                    if deps.setup_task and deps.setup_task.should_run == true then
                        log.info("Condition met, task will run.")
                        return true
                    end
                    log.info("Condition not met, task will be skipped.")
                    return false
                end,
                command = "echo 'Conditional task is running because the function returned true.'"
            },
            {
                name = "abort_task_with_function",
                description = "This task will abort the execution if the abort_if function returns true.",
                params = {
                    abort_execution = "true"
                },
                abort_if = function(params, deps)
                    log.info("Checking abort_if condition for abort_task_with_function...")
                    if params.abort_execution == "true" then
                        log.info("Abort condition met, execution will stop.")
                        return true
                    end
                    log.info("Abort condition not met.")
                    return false
                end,
                command = "echo 'This should not be executed.'"
            },
            {
                name = "final_task_after_abort",
                description = "This task will not be reached if the abort condition is met.",
                depends_on = "abort_task_with_function",
                command = "echo 'This is the final task and should not be reached.'"
            }
        }
    }
}
