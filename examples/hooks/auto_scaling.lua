-- Auto Scaling Hook
-- Automatically scales infrastructure based on resource events

function on_event()
    local event_type = event.type

    -- Handle high CPU usage
    if event_type == "system.cpu_high" or event_type == "agent.resource_high" then
        local current = event.data.current or 0
        local resource = event.data.resource or "unknown"

        if resource == "cpu" and current > 80 then
            log.warn(string.format("High CPU detected: %.2f%% - triggering auto-scale", current))

            -- Dispatch event to trigger scaling workflow
            event.dispatch("deploy.started", {
                deploy_id = "auto-scale-" .. os.time(),
                service = "worker-pool",
                version = "current",
                environment = "production",
                reason = "cpu_high"
            })

            return true
        end
    end

    -- Handle low memory
    if event_type == "system.memory_low" then
        local available = event.data.available or 0
        log.warn(string.format("Low memory: %d bytes available", available))

        -- Scale up memory or add instances
        event.dispatch_custom("scale_memory", "Increasing memory allocation due to low availability")

        return true
    end

    -- Handle disk full
    if event_type == "system.disk_full" then
        local mount_point = event.data.mount_point or "/"
        local usage = event.data.usage_percent or 0

        log.error(string.format("Disk full on %s: %.2f%%", mount_point, usage))

        -- Trigger cleanup or volume expansion
        event.dispatch_custom("expand_disk", mount_point)

        return true
    end

    return true
end
