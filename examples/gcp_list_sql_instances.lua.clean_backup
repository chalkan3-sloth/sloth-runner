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
  gcp_sql_lister = {
    description = "Lists all GCP Cloud SQL instances in a given project.",
    tasks = {
      {
        name = "list_gcp_sql_instances",
        command = function()
          log.info("Listing GCP Cloud SQL instances...")

          local result = gcp.client({ project = "chalkan3" })
            :sql()
            :instances()
            :list()

          if not result.success then
            log.error("Failed to list SQL instances: " .. result.stderr)
            return false, "Failed to list SQL instances."
          end

          log.info("Successfully listed SQL instances.")

          local instances, err = data.parse_json(result.stdout)
          if err then
            log.error("Failed to decode JSON response: " .. err)
            log.info("Raw output: " .. result.stdout)
            return false, "Failed to parse SQL instance list."
          end

          if #instances == 0 then
            log.info("No SQL instances found in project.")
          else
            log.info("Found " .. #instances .. " SQL instance(s):")
            for i, instance in ipairs(instances) do
              log.info("  - " .. instance.name .. " (DB Version: " .. instance.databaseVersion .. ", Region: " .. instance.region .. ")")
            end
          end

          return true, "SQL instances listed successfully."
        end
      }
    }
  }
}
