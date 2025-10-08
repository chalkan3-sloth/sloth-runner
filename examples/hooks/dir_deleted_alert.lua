-- Hook: Directory Deleted Alert
-- Triggers when a directory watcher detects directory deletion
-- Event type: dir.deleted

return {
    name = "dir_deleted_alert",
    description = "Alert when monitored directories are deleted",
    event_types = {"dir.deleted"},
    enabled = true,

    execute = function(event)
        local path = event.data.path or "unknown"
        local watcher_id = event.data.watcher_id or "unknown"
        local agent_name = event.agent_name or "unknown"

        log.warn("🗑️  DIRECTORY DELETED DETECTED!")
        log.warn("  Path: " .. path)
        log.warn("  Watcher: " .. watcher_id)
        log.warn("  Agent: " .. agent_name)

        -- Could trigger backup restore, cleanup, alert ops team, etc
        log.info("✅ Directory deletion processed")

        return true
    end
}
