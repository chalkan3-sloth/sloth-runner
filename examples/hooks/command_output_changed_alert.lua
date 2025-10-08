-- Hook: Command Output Changed Alert
-- Triggers when a command watcher detects output changes
-- Event type: command.output_changed

return {
    name = "command_output_changed_alert",
    description = "Alert when monitored command output changes",
    event_types = {"command.output_changed"},
    enabled = true,

    execute = function(event)
        local command = event.data.command or "unknown"
        local old_output = event.data.old_output or ""
        local new_output = event.data.new_output or ""
        local watcher_id = event.data.watcher_id or "unknown"
        local agent_name = event.agent_name or "unknown"

        log.info("🔍 COMMAND OUTPUT CHANGED!")
        log.info("  Agent: " .. agent_name)
        log.info("  Command: " .. command)
        log.info("  Watcher: " .. watcher_id)

        -- Show truncated output for readability
        local max_len = 100
        if #new_output > max_len then
            log.info("  New output: " .. new_output:sub(1, max_len) .. "...")
        else
            log.info("  New output: " .. new_output)
        end

        -- Could trigger:
        -- - Config change detection
        -- - System state monitoring
        -- - Version change alerts
        log.info("✅ Command output change processed")

        return true
    end
}
