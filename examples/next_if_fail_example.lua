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
    next_if_fail_demo = {
        description = "A workflow to demonstrate the next_if_fail functionality.",
        tasks = {
            {
                name = "task_that_fails",
                description = "This task is designed to fail.",
                command = function()
                    log.error("This task is intentionally failing.")
                    return false, "Intentional failure"
                end
            },
            {
                name = "task_after_failure",
                description = "This task runs only if task_that_fails fails.",
                next_if_fail = "task_that_fails",
                command = "echo 'This task ran because the previous one failed.'"
            },
            {
                name = "task_that_should_be_skipped",
                description = "This task depends on the failing task and should be skipped.",
                depends_on = "task_that_fails",
                command = "echo 'This should not be printed.'"
            },
            {
                name = "final_task",
                description = "This task depends on the task that runs after failure.",
                depends_on = "task_after_failure",
                command = "echo 'This is the final task.'"
            }
        }
    }
}
