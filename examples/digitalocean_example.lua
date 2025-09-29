-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03

local log = require("log")
local droplet_to_delete_name = "my-test-droplet-to-delete"

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
  ["digitalocean-management"] = {
    description = "A pipeline to list and manage DigitalOcean resources.",

    tasks = {
      {
        name = "list_droplets",
        description = "Lists all Droplets in the account.",
        command = function()
          log.info("Listing all DigitalOcean Droplets...")
          local droplets, err = digitalocean.droplets.list()

          if not droplets then
            log.error("Failed to list Droplets: " .. err)
            return false, "doctl list failed."
          end

          log.info("Successfully retrieved Droplet list.")
          print("--- Droplets ---")
          for _, droplet in ipairs(droplets) do
            print(string.format("ID: %d, Name: %s, Status: %s, Region: %s", droplet.id, droplet.name, droplet.status, droplet.region.slug))
          end
          print("----------------")
          
          -- Pass the droplet list to the next task
          return true, "Droplets listed.", {droplets = droplets}
        end
      },
      {
        name = "delete_specific_droplet",
        description = "Finds a specific Droplet by name and deletes it.",
        depends_on = "list_droplets",
        command = function(params, deps)
          local droplets = deps.list_droplets.droplets
          local target_droplet_id = nil

          log.info("Searching for Droplet with name: " .. droplet_to_delete_name)
          for _, droplet in ipairs(droplets) do
            if droplet.name == droplet_to_delete_name then
              target_droplet_id = droplet.id
              break
            end
          end

          if not target_droplet_id then
            log.warn("Could not find a Droplet named '" .. droplet_to_delete_name .. "' to delete. Skipping.")
            -- We return true because not finding the droplet isn't a pipeline failure.
            return true, "Target Droplet not found."
          end

          log.info("Found Droplet with ID: " .. target_droplet_id .. ". Deleting now...")
          local ok, err = digitalocean.droplets.delete({
            id = tostring(target_droplet_id),
            force = true
          })

          if not ok then
            log.error("Failed to delete Droplet: " .. err)
            return false, "Droplet deletion failed."
          end

          log.info("Successfully initiated deletion of Droplet " .. droplet_to_delete_name)
          return true, "Droplet deleted."
        end
      }
    }
  }
}
