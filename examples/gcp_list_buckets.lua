-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:32 -03


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
  gcp_bucket_lister = {
    description = "Lists all GCP Storage buckets in a given project.",
    tasks = {
      {
        name = "list_gcp_buckets",
        command = function()
          log.info("Listing GCP Storage buckets...")

          local result = gcp.client({ project = "chalkan3" })
            :storage()
            :buckets()
            :list()

          if not result.success then
            log.error("Failed to list buckets: " .. result.stderr)
            return false, "Failed to list buckets."
          end

          log.info("Successfully listed buckets.")

          local buckets, err = data.parse_json(result.stdout)
          if err then
            log.error("Failed to decode JSON response: " .. err)
            log.info("Raw output: " .. result.stdout)
            return false, "Failed to parse bucket list."
          end

          if #buckets == 0 then
            log.info("No buckets found in project.")
          else
            log.info("Found " .. #buckets .. " bucket(s):")
            for i, bucket in ipairs(buckets) do
              log.info("  - " .. bucket.name .. " (Location: " .. bucket.location .. ")")
            end
          end

          return true, "Buckets listed successfully."
        end
      }
    }
  }
}
