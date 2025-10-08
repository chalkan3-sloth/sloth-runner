-- Hook: File Created Deploy Trigger
-- Triggers when a new file is created
-- Event type: file.created

return {
    name = "file_created_deploy",
    description = "Trigger deployment when new files are created",
    event_types = {"file.created"},
    enabled = true,

    execute = function(event)
        local path = event.data.path or "unknown"
        local size = event.data.size or 0

        log.info("📄 NEW FILE CREATED!")
        log.info("  Path: " .. path)
        log.info("  Size: " .. tostring(size) .. " bytes")

        -- Example: trigger deployment if config file created
        if path:match("%.conf$") or path:match("%.yaml$") then
            log.info("🚀 Config file detected - would trigger deployment")
            -- event.dispatch("deployment.triggered", {
            --     trigger = "file_created",
            --     file = path
            -- })
        end

        return true
    end
}
