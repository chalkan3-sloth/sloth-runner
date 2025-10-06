-- Security Incident Response Hook
-- Automated response to security events

local MAX_LOGIN_ATTEMPTS = 5
local LOCKOUT_DURATION = 3600 -- 1 hour in seconds

function on_event()
    local event_type = event.type

    -- Handle security breaches - CRITICAL
    if event_type == "security.breach" then
        local user = event.data.user or "unknown"
        local action = event.data.action or "unknown"
        local ip = event.data.ip_address or "unknown"
        local severity = event.data.severity or "high"

        log.error(string.format("🚨 SECURITY BREACH DETECTED 🚨"))
        log.error(string.format("User: %s | IP: %s | Action: %s | Severity: %s", user, ip, action, severity))

        -- Immediate actions:
        -- 1. Lock the affected user account
        -- 2. Block the IP address
        -- 3. Alert security team
        -- 4. Create incident ticket

        -- Send critical alert
        local alert_data = {
            severity = "critical",
            title = "Security Breach Detected",
            message = string.format("Breach attempt by %s from %s", user, ip),
            user = user,
            ip_address = ip,
            action = action,
            timestamp = os.time()
        }

        -- Dispatch to alerting system
        event.dispatch("security.alert", alert_data)

        return true
    end

    -- Handle failed login attempts
    if event_type == "security.login_failed" then
        local user = event.data.user or "unknown"
        local ip = event.data.ip_address or "unknown"
        local attempt_count = event.data.attempt_count or 1

        log.warn(string.format("Failed login: %s from %s (attempt %d)", user, ip, attempt_count))

        -- Lock account after too many attempts
        if attempt_count >= MAX_LOGIN_ATTEMPTS then
            log.error(string.format("Locking account %s due to %d failed attempts", user, attempt_count))

            -- Dispatch account lock event
            event.dispatch("security.account_locked", {
                user = user,
                ip_address = ip,
                attempts = attempt_count,
                lockout_duration = LOCKOUT_DURATION,
                reason = "excessive_failed_logins"
            })

            -- Block IP temporarily
            event.dispatch_custom("block_ip", string.format("Blocking %s for %d seconds", ip, LOCKOUT_DURATION))
        end

        return true
    end

    -- Handle unauthorized access attempts
    if event_type == "security.unauthorized" then
        local user = event.data.user or "unknown"
        local resource = event.data.resource or "unknown"
        local ip = event.data.ip_address or "unknown"

        log.warn(string.format("Unauthorized access: %s trying to access %s from %s", user, resource, ip))

        -- Check if this is a pattern (multiple unauthorized attempts)
        -- In production, you'd query a state store for historical data
        -- For now, just log and alert

        event.dispatch("security.alert", {
            severity = "medium",
            title = "Unauthorized Access Attempt",
            message = string.format("%s attempted to access %s", user, resource),
            user = user,
            resource = resource,
            ip_address = ip
        })

        return true
    end

    -- Handle permission denied events
    if event_type == "security.permission_denied" then
        local user = event.data.user or "unknown"
        local resource = event.data.resource or "unknown"
        local required_permission = event.data.required_permission or "unknown"

        log.info(string.format("Permission denied: %s lacks %s for %s", user, required_permission, resource))

        -- This is informational, but log for audit trail
        return true
    end

    -- Handle successful logins for audit
    if event_type == "security.login_success" then
        local user = event.data.user or "unknown"
        local ip = event.data.ip_address or "unknown"
        local session_id = event.data.session_id or "unknown"

        log.info(string.format("Successful login: %s from %s (session: %s)", user, ip, session_id))

        -- Audit log successful logins
        -- In production, store this in a secure audit database
        return true
    end

    return true
end
