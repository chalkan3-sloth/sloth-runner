-- Hook: File Changed Alert
-- Triggers when a file watcher detects changes
-- Event type: file.changed

return {
    name = "file_changed_alert",
    description = "Alert when monitored files change",
    event_types = {"file.changed"},
    enabled = true,

    execute = function(event)
        local path = event.data.path or "unknown"
        local old_size = event.data.old_size or 0
        local new_size = event.data.new_size or 0
        local old_hash = event.data.old_hash or ""
        local new_hash = event.data.new_hash or ""

        log.info("🔄 FILE CHANGED DETECTED!")
        log.info("  Path: " .. path)
        log.info("  Size: " .. tostring(old_size) .. " -> " .. tostring(new_size))

        if old_hash ~= "" and new_hash ~= "" then
            log.info("  Hash: " .. old_hash:sub(1,8) .. "... -> " .. new_hash:sub(1,8) .. "...")
        end

        -- Could send notification, trigger deployment, etc
        log.info("✅ File change processed")

        return true
    end
}
