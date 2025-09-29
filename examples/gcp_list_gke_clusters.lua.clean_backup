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
  gcp_gke_lister = {
    description = "Lists all GKE clusters in a given project.",
    tasks = {
      {
        name = "list_gke_clusters",
        command = function()
          log.info("Listing GKE clusters...")

          local result = gcp.client({ project = "chalkan3" })
            :gke()
            :clusters()
            :list()

          if not result.success then
            log.error("Failed to list GKE clusters: " .. result.stderr)
            return false, "Failed to list GKE clusters."
          end

          log.info("Successfully listed GKE clusters.")

          local clusters, err = data.parse_json(result.stdout)
          if err then
            log.error("Failed to decode JSON response: " .. err)
            log.info("Raw output: " .. result.stdout)
            return false, "Failed to parse GKE cluster list."
          end

          if #clusters == 0 then
            log.info("No GKE clusters found in project.")
          else
            log.info("Found " .. #clusters .. " GKE cluster(s):")
            for i, cluster in ipairs(clusters) do
              log.info("  - " .. cluster.name .. " (Status: " .. cluster.status .. ", Location: " .. cluster.location .. ")")
            end
          end

          return true, "GKE clusters listed successfully."
        end
      }
    }
  }
}
