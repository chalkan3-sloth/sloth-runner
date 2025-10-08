-- Hook: Directory Changed Alert
-- Triggers when a directory watcher detects changes
-- Event type: dir.changed

return {
    name = "dir_changed_alert",
    description = "Alert when monitored directories change",
    event_types = {"dir.changed"},
    enabled = true,

    execute = function(event)
        local path = event.data.path or "unknown"
        local files_added = event.data.files_added or 0
        local files_removed = event.data.files_removed or 0
        local files_modified = event.data.files_modified or 0
        local watcher_id = event.data.watcher_id or "unknown"

        log.info("📁 DIRECTORY CHANGED DETECTED!")
        log.info("  Path: " .. path)
        log.info("  Watcher: " .. watcher_id)
        log.info("  Files added: " .. tostring(files_added))
        log.info("  Files removed: " .. tostring(files_removed))
        log.info("  Files modified: " .. tostring(files_modified))

        local total_changes = files_added + files_removed + files_modified
        log.info("  Total changes: " .. tostring(total_changes))

        -- Could trigger rebuild, notification, etc
        log.info("✅ Directory change processed")

        return true
    end
}
