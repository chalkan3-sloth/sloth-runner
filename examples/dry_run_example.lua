-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:30 -03


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
    dry_run_demo = {
        description = "A workflow to demonstrate the dry-run functionality.",
        tasks = {
            {
                name = "task_one",
                description = "This task would normally do something.",
                command = "echo 'Executing Task One'"
            },
            {
                name = "task_two",
                description = "This task depends on the first one.",
                depends_on = "task_one",
                command = "echo 'Executing Task Two'"
            },
            {
                name = "task_three",
                description = "This task would also do something.",
                command = "echo 'Executing Task Three'"
            }
        }
    }
}
