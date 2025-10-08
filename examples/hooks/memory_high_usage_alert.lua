-- Hook: Memory High Usage Alert
-- Triggers when memory usage exceeds threshold
-- Event type: memory.high_usage

return {
    name = "memory_high_usage_alert",
    description = "Alert when memory usage is high",
    event_types = {"memory.high_usage"},
    enabled = true,

    execute = function(event)
        local percent = event.data.memory_percent or 0
        local used_mb = event.data.used_mb or 0
        local total_mb = event.data.total_mb or 0
        local threshold = event.data.threshold or 0
        local watcher_id = event.data.watcher_id or "unknown"
        local agent_name = event.agent_name or "unknown"

        log.warn("⚠️  HIGH MEMORY USAGE DETECTED!")
        log.warn("  Agent: " .. agent_name)
        log.warn("  Watcher: " .. watcher_id)
        log.warn("  Memory usage: " .. string.format("%.1f%%", percent))
        log.warn("  Used: " .. tostring(used_mb) .. " MB / " .. tostring(total_mb) .. " MB")
        log.warn("  Threshold: " .. string.format("%.1f%%", threshold))

        -- Could trigger:
        -- - Kill memory-intensive processes
        -- - Restart services
        -- - Scale up resources
        -- - Send alerts to ops team
        log.info("✅ Memory alert processed")

        return true
    end
}
