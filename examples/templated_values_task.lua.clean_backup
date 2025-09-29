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
  templated_values_group = {
    description = "Demonstrates templating in values.yaml.",
    tasks = {
      {
        name = "print_templated_value",
        description = "Prints a value from values.yaml that was templated with an environment variable.",
        command = function()
          log.info("Templated value: " .. values.my_value)
          return true, "Printed templated value."
        end
      }
    }
  }
}
