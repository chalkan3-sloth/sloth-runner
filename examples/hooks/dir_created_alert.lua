-- Hook: Directory Created Alert
-- Triggers when a directory watcher detects directory creation
-- Event type: dir.created

return {
    name = "dir_created_alert",
    description = "Alert when monitored directories are created",
    event_types = {"dir.created"},
    enabled = true,

    execute = function(event)
        local path = event.data.path or "unknown"
        local watcher_id = event.data.watcher_id or "unknown"
        local agent_name = event.agent_name or "unknown"

        log.info("📂 DIRECTORY CREATED DETECTED!")
        log.info("  Path: " .. path)
        log.info("  Watcher: " .. watcher_id)
        log.info("  Agent: " .. agent_name)

        -- Could trigger permissions setup, initialization scripts, etc
        log.info("✅ Directory creation processed")

        return true
    end
}
