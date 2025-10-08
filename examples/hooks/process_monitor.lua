-- Hook: Process Monitor
-- Triggers when processes start or stop
-- Event types: process.started, process.stopped

return {
    name = "process_monitor",
    description = "Monitor critical process lifecycle",
    event_types = {"process.started", "process.stopped"},
    enabled = true,

    execute = function(event)
        local process = event.data.process or "unknown"
        local event_type = event.type

        if event_type == "process.started" then
            log.info("✅ PROCESS STARTED: " .. process)
            -- Could update monitoring dashboard, send notification
        elseif event_type == "process.stopped" then
            log.warn("⚠️  PROCESS STOPPED: " .. process)
            -- Could trigger restart, alert on-call, etc
            log.warn("  Action: Checking if restart needed")
        end

        return true
    end
}
