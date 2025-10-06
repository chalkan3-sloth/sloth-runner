-- Deployment Automation Hook
-- Automates deployment lifecycle and rollback procedures

function on_event()
    local event_type = event.type

    -- Handle deployment started
    if event_type == "deploy.started" then
        local deploy_id = event.data.deploy_id or "unknown"
        local service = event.data.service or "unknown"
        local version = event.data.version or "unknown"
        local environment = event.data.environment or "unknown"

        log.info(string.format("🚀 Deployment started: %s v%s to %s (ID: %s)",
                               service, version, environment, deploy_id))

        -- Send notification to deployment channel
        event.dispatch_custom("notify_team",
            string.format("Deployment started: %s v%s → %s", service, version, environment))

        -- Start deployment monitoring
        event.dispatch("health.check_started", {
            service = service,
            deploy_id = deploy_id,
            check_interval = 30 -- seconds
        })

        return true
    end

    -- Handle successful deployment
    if event_type == "deploy.completed" then
        local deploy_id = event.data.deploy_id or "unknown"
        local service = event.data.service or "unknown"
        local version = event.data.version or "unknown"
        local duration = event.data.duration or "0s"

        log.info(string.format("✅ Deployment completed: %s v%s (duration: %s)",
                               service, version, duration))

        -- Run post-deployment smoke tests
        event.dispatch_custom("run_smoke_tests", service)

        -- Update service catalog
        event.dispatch("service.updated", {
            service = service,
            version = version,
            deploy_id = deploy_id,
            status = "deployed"
        })

        -- Send success notification
        event.dispatch_custom("notify_success",
            string.format("✅ %s v%s deployed successfully in %s", service, version, duration))

        return true
    end

    -- Handle deployment failure - CRITICAL
    if event_type == "deploy.failed" then
        local deploy_id = event.data.deploy_id or "unknown"
        local service = event.data.service or "unknown"
        local error_msg = event.data.error or "unknown error"
        local rollback_needed = event.data.rollback_needed

        log.error(string.format("❌ Deployment failed: %s (ID: %s)", service, deploy_id))
        log.error(string.format("Error: %s", error_msg))

        -- Automatic rollback if needed
        if rollback_needed then
            log.warn(string.format("Initiating automatic rollback for %s", service))

            -- Get previous version from state or metadata
            local prev_version = event.data.prev_version or "previous"

            event.dispatch("deploy.rollback", {
                deploy_id = deploy_id,
                service = service,
                prev_version = prev_version,
                reason = "automated_rollback_on_failure",
                original_error = error_msg
            })
        end

        -- Alert ops team
        event.dispatch_custom("alert_ops",
            string.format("🚨 Deployment failed: %s\\nError: %s", service, error_msg))

        return true
    end

    -- Handle rollback
    if event_type == "deploy.rollback" then
        local deploy_id = event.data.deploy_id or "unknown"
        local service = event.data.service or "unknown"
        local prev_version = event.data.prev_version or "previous"
        local reason = event.data.reason or "manual"

        log.warn(string.format("⏮️  Rolling back %s to %s (reason: %s)",
                               service, prev_version, reason))

        -- Execute rollback workflow
        event.dispatch_custom("execute_rollback", service)

        -- Monitor rollback health
        event.dispatch("health.check_started", {
            service = service,
            deploy_id = "rollback-" .. deploy_id,
            check_type = "rollback_verification"
        })

        -- Notify team
        event.dispatch_custom("notify_rollback",
            string.format("⏮️ Rolling back %s to %s", service, prev_version))

        return true
    end

    -- Handle health check failures during deployment
    if event_type == "health.check_failed" then
        local service = event.data.service or "unknown"
        local check_name = event.data.check_name or "unknown"
        local consecutive_failures = event.data.consecutive_failures or 1

        log.error(string.format("Health check failed: %s (%s) - failures: %d",
                                service, check_name, consecutive_failures))

        -- If multiple consecutive failures, trigger rollback
        if consecutive_failures >= 3 then
            log.error(string.format("Too many health check failures for %s - triggering rollback", service))

            event.dispatch("deploy.rollback", {
                deploy_id = "health-failure-" .. os.time(),
                service = service,
                reason = "health_check_failures",
                consecutive_failures = consecutive_failures
            })
        end

        return true
    end

    -- Handle service degradation
    if event_type == "health.degraded" then
        local service = event.data.service or "unknown"
        local metrics = event.data.metrics or {}

        log.warn(string.format("Service degraded: %s", service))

        -- Could trigger auto-scaling here
        event.dispatch_custom("scale_up", service)

        return true
    end

    -- Handle service recovery
    if event_type == "health.recovered" then
        local service = event.data.service or "unknown"
        local downtime = event.data.downtime_duration or "unknown"

        log.info(string.format("Service recovered: %s (downtime: %s)", service, downtime))

        -- Send recovery notification
        event.dispatch_custom("notify_recovery",
            string.format("✅ %s has recovered after %s", service, downtime))

        return true
    end

    return true
end
