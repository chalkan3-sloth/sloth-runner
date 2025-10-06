-- Comprehensive Monitoring Hook
-- This hook demonstrates how to respond to various system events

function on_event()
    local event_type = event.type
    local data = event.data or {}

    log.info("Processing event: " .. event_type)

    -- Handle Agent Events
    if event_type == "agent.registered" then
        local agent_name = data.agent and data.agent.name or "unknown"
        log.info("New agent registered: " .. agent_name)
        -- Send notification to Slack/Discord
        -- http.post("https://hooks.slack.com/...", {...})
        return true

    elseif event_type == "agent.disconnected" then
        local agent_name = data.agent and data.agent.name or "unknown"
        log.error("Agent disconnected: " .. agent_name)
        -- Alert operations team
        return true

    elseif event_type == "agent.heartbeat_failed" then
        local agent_name = data.agent and data.agent.name or "unknown"
        log.warn("Agent heartbeat failed: " .. agent_name)
        -- Attempt to restart agent or escalate
        return true

    elseif event_type == "agent.resource_high" then
        local resource = data.resource or "unknown"
        local current = data.current or 0
        log.warn(string.format("High resource usage on %s: %.2f%%", resource, current))
        -- Scale up or alert
        return true

    -- Handle Task Events
    elseif event_type == "task.failed" then
        local task_name = data.task and data.task.task_name or "unknown"
        local error_msg = data.task and data.task.error or "no error message"
        log.error("Task failed: " .. task_name .. " - " .. error_msg)
        -- Auto-retry logic or escalation
        return true

    elseif event_type == "task.timeout" then
        local task_name = data.task and data.task.task_name or "unknown"
        log.warn("Task timeout: " .. task_name)
        -- Kill and restart or alert
        return true

    -- Handle Workflow Events
    elseif event_type == "workflow.started" then
        local workflow_name = data.workflow_name or "unknown"
        log.info("Workflow started: " .. workflow_name)
        -- Track workflow execution
        return true

    elseif event_type == "workflow.failed" then
        local workflow_name = data.workflow_name or "unknown"
        local failed_tasks = data.failed_tasks or {}
        log.error("Workflow failed: " .. workflow_name)
        for i, task in ipairs(failed_tasks) do
            log.error("  - Failed task: " .. task)
        end
        -- Send detailed failure report
        return true

    -- Handle System Events
    elseif event_type == "system.startup" then
        log.info("System startup detected")
        -- Initialize monitoring, send startup notification
        return true

    elseif event_type == "system.shutdown" then
        log.info("System shutdown detected")
        -- Graceful cleanup, save state
        return true

    elseif event_type == "system.disk_full" then
        log.error("Disk full warning!")
        -- Cleanup old files, expand disk, alert
        return true

    elseif event_type == "system.memory_low" then
        log.warn("Low memory warning")
        -- Kill non-essential processes, scale up
        return true

    -- Handle Scheduler Events
    elseif event_type == "schedule.missed" then
        local schedule_name = data.schedule_name or "unknown"
        log.warn("Scheduled job missed: " .. schedule_name)
        -- Alert and possibly reschedule
        return true

    elseif event_type == "schedule.triggered" then
        local schedule_name = data.schedule_name or "unknown"
        log.info("Schedule triggered: " .. schedule_name)
        return true

    -- Handle State Events
    elseif event_type == "state.corrupted" then
        local state_key = data.state_key or "unknown"
        log.error("State corruption detected: " .. state_key)
        -- Restore from backup, alert
        return true

    -- Handle Secret Events
    elseif event_type == "secret.rotation_needed" then
        local secret_name = data.secret_name or "unknown"
        log.warn("Secret rotation needed: " .. secret_name)
        -- Trigger rotation workflow
        return true

    elseif event_type == "secret.expired" then
        local secret_name = data.secret_name or "unknown"
        log.error("Secret expired: " .. secret_name)
        -- Emergency rotation or disable dependent services
        return true

    elseif event_type == "secret.accessed" then
        local secret_name = data.secret_name or "unknown"
        local accessed_by = data.accessed_by or "unknown"
        log.info("Secret accessed: " .. secret_name .. " by " .. accessed_by)
        -- Audit logging
        return true

    -- Handle Stack Events
    elseif event_type == "stack.drift_detected" then
        local stack_name = data.stack_name or "unknown"
        log.warn("Infrastructure drift detected: " .. stack_name)
        -- Auto-remediate or alert
        return true

    elseif event_type == "stack.failed" then
        local stack_name = data.stack_name or "unknown"
        log.error("Stack deployment failed: " .. stack_name)
        -- Rollback and alert
        return true

    -- Handle Backup Events
    elseif event_type == "backup.failed" then
        local backup_id = data.backup_id or "unknown"
        log.error("Backup failed: " .. backup_id)
        -- Retry backup, alert ops team
        return true

    elseif event_type == "backup.completed" then
        local backup_id = data.backup_id or "unknown"
        local size = data.size or 0
        log.info(string.format("Backup completed: %s (size: %d bytes)", backup_id, size))
        -- Verify backup integrity
        return true

    -- Handle Database Events
    elseif event_type == "db.query_slow" then
        local query = data.query or "unknown"
        local duration = data.duration or 0
        log.warn(string.format("Slow query detected (%.2fms): %s", duration, query))
        -- Add to query optimization list
        return true

    elseif event_type == "db.error" then
        local database = data.database or "unknown"
        local error_msg = data.error or "no error"
        log.error("Database error on " .. database .. ": " .. error_msg)
        -- Check connection pool, restart if needed
        return true

    -- Handle Network Events
    elseif event_type == "network.down" then
        log.error("Network connectivity lost")
        -- Switch to backup network, alert
        return true

    elseif event_type == "network.latency_high" then
        local latency = data.latency or 0
        log.warn(string.format("High network latency: %.2fms", latency))
        -- Check network health, reroute traffic
        return true

    -- Handle Security Events
    elseif event_type == "security.breach" then
        local action = data.action or "unknown"
        local user = data.user or "unknown"
        log.error("SECURITY BREACH: " .. action .. " by " .. user)
        -- Lock down system, alert security team
        return true

    elseif event_type == "security.login_failed" then
        local user = data.user or "unknown"
        local ip = data.ip_address or "unknown"
        log.warn("Failed login attempt: " .. user .. " from " .. ip)
        -- Track failed attempts, trigger lockout if needed
        return true

    elseif event_type == "security.unauthorized" then
        local user = data.user or "unknown"
        local resource = data.resource or "unknown"
        log.warn("Unauthorized access attempt: " .. user .. " -> " .. resource)
        -- Audit log, alert if pattern detected
        return true

    -- Handle Deploy Events
    elseif event_type == "deploy.completed" then
        local service = data.service or "unknown"
        local version = data.version or "unknown"
        local environment = data.environment or "unknown"
        log.info(string.format("Deployment completed: %s v%s to %s", service, version, environment))
        -- Run smoke tests, update monitoring
        return true

    elseif event_type == "deploy.failed" then
        local service = data.service or "unknown"
        log.error("Deployment failed: " .. service)
        -- Automatic rollback
        return true

    elseif event_type == "deploy.rollback" then
        local service = data.service or "unknown"
        local prev_version = data.prev_version or "unknown"
        log.warn("Rolling back " .. service .. " to " .. prev_version)
        return true

    -- Handle Health Check Events
    elseif event_type == "health.check_failed" then
        local service = data.service or "unknown"
        local check_name = data.check_name or "unknown"
        log.error("Health check failed: " .. service .. " (" .. check_name .. ")")
        -- Restart service, alert on-call
        return true

    elseif event_type == "health.degraded" then
        local service = data.service or "unknown"
        log.warn("Service degraded: " .. service)
        -- Scale up, investigate
        return true

    elseif event_type == "health.recovered" then
        local service = data.service or "unknown"
        log.info("Service recovered: " .. service)
        -- Clear alerts
        return true

    -- Handle File Events
    elseif event_type == "file.created" or event_type == "file.modified" or event_type == "file.deleted" then
        local path = data.path or "unknown"
        log.info("File event: " .. event_type .. " - " .. path)
        -- Trigger build, sync, backup, etc.
        return true

    -- Handle Custom Events
    elseif event_type == "custom" then
        local name = data.name or "unknown"
        local message = data.message or "no message"
        log.info("Custom event: " .. name .. " - " .. message)
        return true

    else
        log.warn("Unhandled event type: " .. event_type)
        return true
    end
end
