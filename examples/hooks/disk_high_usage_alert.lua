-- Hook: Disk High Usage Alert
-- Triggers when disk usage exceeds threshold
-- Event type: disk.high_usage

return {
    name = "disk_high_usage_alert",
    description = "Alert when disk usage is high",
    event_types = {"disk.high_usage"},
    enabled = true,

    execute = function(event)
        local percent = event.data.disk_percent or 0
        local used_gb = event.data.used_gb or 0
        local total_gb = event.data.total_gb or 0
        local threshold = event.data.threshold or 0
        local path = event.data.path or "unknown"
        local watcher_id = event.data.watcher_id or "unknown"
        local agent_name = event.agent_name or "unknown"

        log.warn("⚠️  HIGH DISK USAGE DETECTED!")
        log.warn("  Agent: " .. agent_name)
        log.warn("  Path: " .. path)
        log.warn("  Watcher: " .. watcher_id)
        log.warn("  Disk usage: " .. string.format("%.1f%%", percent))
        log.warn("  Used: " .. string.format("%.1f GB", used_gb) .. " / " .. string.format("%.1f GB", total_gb))
        log.warn("  Threshold: " .. string.format("%.1f%%", threshold))

        -- Could trigger:
        -- - Clean up old files
        -- - Rotate logs
        -- - Expand disk
        -- - Alert ops team
        log.info("✅ Disk usage alert processed")

        return true
    end
}
