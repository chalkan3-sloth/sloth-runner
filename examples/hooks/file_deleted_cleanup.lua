-- Hook: File Deleted Cleanup
-- Triggers when a monitored file is deleted
-- Event type: file.deleted

return {
    name = "file_deleted_cleanup",
    description = "Cleanup actions when files are deleted",
    event_types = {"file.deleted"},
    enabled = true,

    execute = function(event)
        local path = event.data.path or "unknown"

        log.warn("🗑️  FILE DELETED!")
        log.warn("  Path: " .. path)

        -- Example: cleanup related resources
        log.info("🧹 Running cleanup tasks")
        -- Could remove cache, update index, notify team, etc

        return true
    end
}
