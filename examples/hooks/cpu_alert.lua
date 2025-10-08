-- Hook: CPU Alert
-- Triggers when CPU usage exceeds threshold
-- Event type: cpu.high_usage

return {
    name = "cpu_alert",
    description = "Alert on high CPU usage",
    event_types = {"cpu.high_usage"},
    enabled = true,

    execute = function(event)
        local cpu_percent = event.data.cpu_percent or 0
        local threshold = event.data.threshold or 0
        local load_1min = event.data.load_1min or 0

        log.warn("🔥 HIGH CPU USAGE DETECTED!")
        log.warn("  CPU: " .. string.format("%.2f%%", cpu_percent))
        log.warn("  Load: " .. string.format("%.2f", load_1min))
        log.warn("  Threshold: " .. string.format("%.2f%%", threshold))

        -- Example actions
        if cpu_percent > threshold * 1.5 then
            log.error("🚨 CRITICAL: CPU usage extremely high!")
            -- Could trigger autoscaling, alert team, kill runaway processes
        else
            log.warn("⚠️  WARNING: CPU usage elevated")
            -- Could log metrics, prepare for scaling
        end

        return true
    end
}
