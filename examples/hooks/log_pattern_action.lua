-- Hook: Log Pattern Action
-- Triggers when specific patterns are found in logs
-- Event type: log.pattern_matched

return {
    name = "log_pattern_action",
    description = "Take action when log patterns match",
    event_types = {"log.pattern_matched"},
    enabled = true,

    execute = function(event)
        local log_path = event.data.log_path or "unknown"
        local pattern = event.data.pattern or ""
        local line = event.data.line or ""

        log.info("📋 LOG PATTERN MATCHED!")
        log.info("  File: " .. log_path)
        log.info("  Pattern: " .. pattern)
        log.info("  Line: " .. line:sub(1, 100))  -- First 100 chars

        -- Example: handle different patterns
        if line:match("ERROR") or line:match("FATAL") then
            log.error("🚨 Error detected in logs!")
            -- Could create incident, page on-call, trigger runbook
        elseif line:match("WARN") then
            log.warn("⚠️  Warning detected in logs")
            -- Could update dashboard, log to metrics
        end

        return true
    end
}
