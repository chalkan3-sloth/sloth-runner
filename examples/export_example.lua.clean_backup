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
  main = {
    description = "A task group to demonstrate the export function.",
    tasks = {
      {
        name = "export-data-task",
        description = "Exports a table and also returns a value.",
        command = function(params, inputs)
          log.info("Exporting some data...")

          -- Use the global export function to send a table to the runner
          export({
            exported_value = "this came from the export function",
            another_key = 12345,
            is_exported = true
          })

          log.info("Export complete. The task will now finish and return its own output.")

          -- The task's own return value will be merged with the exported data.
          -- If keys conflict, the exported value will win.
          return true, "Task finished successfully.", { task_return_value = "this came from the task's return" }
        end
      }
    }
  }
}
