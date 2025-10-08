-- Hook: Service Status Changed Alert
-- Triggers when a service changes status
-- Event type: service.status_changed

return {
    name = "service_status_changed_alert",
    description = "Alert when monitored services change status",
    event_types = {"service.status_changed"},
    enabled = true,

    execute = function(event)
        local service_name = event.data.service_name or "unknown"
        local old_status = event.data.old_status or "unknown"
        local new_status = event.data.new_status or "unknown"
        local watcher_id = event.data.watcher_id or "unknown"
        local agent_name = event.agent_name or "unknown"

        log.info("🔄 SERVICE STATUS CHANGED!")
        log.info("  Agent: " .. agent_name)
        log.info("  Service: " .. service_name)
        log.info("  Watcher: " .. watcher_id)
        log.info("  Status: " .. old_status .. " -> " .. new_status)

        -- Check if service went down
        if new_status == "inactive" or new_status == "failed" then
            log.warn("  ⚠️  Service is down!")
            -- Could trigger:
            -- - Restart service
            -- - Alert ops team
            -- - Failover to backup
        elseif new_status == "active" and old_status ~= "active" then
            log.info("  ✅ Service is now active")
        end

        log.info("✅ Service status change processed")

        return true
    end
}
