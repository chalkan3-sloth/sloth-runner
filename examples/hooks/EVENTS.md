# Sloth Runner - Event Reference

This document lists all available events that can trigger hooks in the Sloth Runner system.

## Event Categories

### 1. Agent Events
Events related to agent lifecycle and health.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `agent.registered` | Agent successfully registered | `agent.name`, `agent.address`, `agent.version` |
| `agent.connected` | Agent established connection | `agent.name`, `agent.address` |
| `agent.disconnected` | Agent lost connection | `agent.name`, `last_seen` |
| `agent.heartbeat_failed` | Agent heartbeat check failed | `agent.name`, `consecutive_failures` |
| `agent.updated` | Agent software updated | `agent.name`, `old_version`, `new_version` |
| `agent.version_mismatch` | Agent version incompatible | `agent.name`, `agent_version`, `server_version` |
| `agent.resource_high` | High resource usage on agent | `agent.name`, `resource`, `current`, `threshold` |

### 2. Task Events
Events related to individual task execution.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `task.started` | Task execution started | `task_name`, `agent_name`, `start_time` |
| `task.completed` | Task completed successfully | `task_name`, `agent_name`, `exit_code`, `duration` |
| `task.failed` | Task execution failed | `task_name`, `agent_name`, `error`, `exit_code` |
| `task.timeout` | Task exceeded time limit | `task_name`, `agent_name`, `timeout_duration` |
| `task.retrying` | Task being retried | `task_name`, `attempt`, `max_retries` |
| `task.cancelled` | Task was cancelled | `task_name`, `agent_name`, `cancelled_by` |

### 3. Workflow Events
Events related to workflow execution.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `workflow.started` | Workflow execution started | `workflow_name`, `workflow_id`, `task_count` |
| `workflow.completed` | Workflow completed successfully | `workflow_name`, `workflow_id`, `duration` |
| `workflow.failed` | Workflow execution failed | `workflow_name`, `failed_tasks[]`, `error` |
| `workflow.paused` | Workflow paused | `workflow_name`, `workflow_id`, `paused_by` |
| `workflow.resumed` | Workflow resumed | `workflow_name`, `workflow_id`, `resumed_by` |
| `workflow.cancelled` | Workflow was cancelled | `workflow_name`, `workflow_id`, `cancelled_by` |

### 4. System Events
Events related to system-level operations.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `system.startup` | System started | `version`, `start_time` |
| `system.shutdown` | System shutting down | `reason`, `uptime` |
| `system.error` | System-level error occurred | `component`, `error`, `severity` |
| `system.warning` | System warning | `component`, `message`, `severity` |
| `system.resource_high` | System resource usage high | `resource`, `current`, `threshold` |
| `system.disk_full` | Disk space critical | `mount_point`, `usage_percent`, `available_bytes` |
| `system.memory_low` | Memory running low | `total`, `available`, `usage_percent` |
| `system.cpu_high` | CPU usage high | `usage_percent`, `load_average` |

### 5. Scheduler Events
Events related to scheduled jobs.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `schedule.triggered` | Scheduled job triggered | `schedule_name`, `schedule_id`, `workflow_name` |
| `schedule.missed` | Scheduled job missed | `schedule_name`, `expected_time`, `actual_time` |
| `schedule.created` | New schedule created | `schedule_name`, `cron_expr`, `workflow_name` |
| `schedule.deleted` | Schedule deleted | `schedule_name`, `deleted_by` |
| `schedule.updated` | Schedule configuration updated | `schedule_name`, `changes` |
| `schedule.enabled` | Schedule enabled | `schedule_name`, `enabled_by` |
| `schedule.disabled` | Schedule disabled | `schedule_name`, `disabled_by` |

### 6. State Events
Events related to state management.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `state.created` | New state created | `state_key`, `namespace`, `value` |
| `state.updated` | State value updated | `state_key`, `old_value`, `new_value` |
| `state.deleted` | State deleted | `state_key`, `namespace` |
| `state.corrupted` | State corruption detected | `state_key`, `corruption_type` |
| `state.locked` | State locked | `state_key`, `locked_by`, `lock_id` |
| `state.unlocked` | State unlocked | `state_key`, `lock_id` |

### 7. Secret Events
Events related to secret management.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `secret.created` | New secret created | `secret_name`, `namespace`, `created_by` |
| `secret.accessed` | Secret accessed | `secret_name`, `accessed_by`, `access_time` |
| `secret.deleted` | Secret deleted | `secret_name`, `deleted_by` |
| `secret.updated` | Secret value updated | `secret_name`, `updated_by` |
| `secret.rotation_needed` | Secret needs rotation | `secret_name`, `last_rotated`, `rotation_policy` |
| `secret.expired` | Secret expired | `secret_name`, `expired_at` |

### 8. Stack Events
Events related to infrastructure stack management.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `stack.deployed` | Stack deployed successfully | `stack_name`, `provider`, `resources[]` |
| `stack.destroyed` | Stack destroyed | `stack_name`, `destroyed_by` |
| `stack.updated` | Stack configuration updated | `stack_name`, `changes` |
| `stack.drift_detected` | Configuration drift detected | `stack_name`, `drift_info`, `resources[]` |
| `stack.failed` | Stack operation failed | `stack_name`, `operation`, `error` |

### 9. Backup Events
Events related to backup and restore operations.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `backup.started` | Backup operation started | `backup_id`, `backup_type`, `source` |
| `backup.completed` | Backup completed successfully | `backup_id`, `size`, `duration`, `destination` |
| `backup.failed` | Backup operation failed | `backup_id`, `error`, `partial_data` |
| `restore.started` | Restore operation started | `backup_id`, `destination` |
| `restore.completed` | Restore completed successfully | `backup_id`, `duration` |
| `restore.failed` | Restore operation failed | `backup_id`, `error` |

### 10. Database Events
Events related to database operations.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `db.connected` | Database connection established | `database`, `connection_pool_size` |
| `db.disconnected` | Database connection lost | `database`, `reason` |
| `db.query_slow` | Slow query detected | `database`, `query`, `duration`, `threshold` |
| `db.error` | Database error occurred | `database`, `operation`, `error` |
| `db.migration` | Database migration executed | `database`, `version`, `migration_name` |

### 11. Network Events
Events related to network connectivity and performance.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `network.down` | Network connectivity lost | `interface`, `last_successful_ping` |
| `network.up` | Network connectivity restored | `interface`, `downtime_duration` |
| `network.slow` | Network performance degraded | `interface`, `bandwidth`, `expected_bandwidth` |
| `network.latency_high` | High network latency detected | `remote_host`, `latency`, `threshold` |

### 12. Security Events
Events related to security and access control.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `security.breach` | Security breach detected | `user`, `action`, `resource`, `ip_address`, `severity` |
| `security.unauthorized` | Unauthorized access attempt | `user`, `resource`, `action`, `ip_address` |
| `security.login_failed` | Login attempt failed | `user`, `ip_address`, `reason`, `attempt_count` |
| `security.login_success` | Successful login | `user`, `ip_address`, `session_id` |
| `security.permission_denied` | Permission denied | `user`, `resource`, `required_permission` |

### 13. File System Events
Events related to file system changes.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `file.created` | File created | `path`, `size`, `mode`, `watcher` |
| `file.modified` | File modified | `path`, `size`, `modification_time` |
| `file.deleted` | File deleted | `path`, `watcher` |
| `file.renamed` | File renamed | `old_path`, `new_path` |
| `dir.created` | Directory created | `path`, `watcher` |
| `dir.deleted` | Directory deleted | `path`, `watcher` |

### 14. Deploy Events
Events related to application deployment.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `deploy.started` | Deployment started | `deploy_id`, `service`, `version`, `environment` |
| `deploy.completed` | Deployment completed | `deploy_id`, `service`, `version`, `duration` |
| `deploy.failed` | Deployment failed | `deploy_id`, `service`, `error`, `rollback_needed` |
| `deploy.rollback` | Deployment rolled back | `deploy_id`, `service`, `prev_version`, `reason` |

### 15. Health Check Events
Events related to service health monitoring.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `health.check_passed` | Health check passed | `service`, `check_name`, `response_time` |
| `health.check_failed` | Health check failed | `service`, `check_name`, `error`, `consecutive_failures` |
| `health.degraded` | Service performance degraded | `service`, `metrics`, `threshold` |
| `health.recovered` | Service recovered | `service`, `downtime_duration` |

### 16. Custom Events
User-defined custom events.

| Event Type | Description | Data Fields |
|------------|-------------|-------------|
| `custom` | Custom user-defined event | `name`, `message`, `payload{}` |

## Usage Examples

### Dispatching Events from Workflows

```lua
-- Dispatch a deployment event
event.dispatch("deploy.completed", {
    deploy_id = "deploy-12345",
    service = "api-server",
    version = "v2.3.1",
    environment = "production",
    duration = "5m30s"
})

-- Dispatch a security event
event.dispatch("security.login_failed", {
    user = "admin",
    ip_address = "192.168.1.100",
    reason = "invalid_password",
    attempt_count = 3
})

-- Dispatch a custom event
event.dispatch_custom("deployment_complete", "App deployed successfully to production")
```

### Creating Hooks

```bash
# Create a hook for failed deployments
sloth-runner hook add deploy-alert \
  --event deploy.failed \
  --file /path/to/alert_on_deploy_failure.lua \
  --description "Alert ops team on deployment failures"

# Create a hook for security breaches
sloth-runner hook add security-breach \
  --event security.breach \
  --file /path/to/security_incident_response.lua \
  --description "Automated security incident response"

# Create a hook for agent disconnections
sloth-runner hook add agent-reconnect \
  --event agent.disconnected \
  --file /path/to/auto_reconnect_agent.lua \
  --description "Attempt to reconnect disconnected agents"
```

### Hook Script Example

```lua
-- Hook script for handling failed deployments
function on_event()
    local deploy_id = event.data.deploy_id
    local service = event.data.service
    local error_msg = event.data.error

    -- Log the failure
    log.error("Deployment failed: " .. service .. " (" .. deploy_id .. ")")
    log.error("Error: " .. error_msg)

    -- Send alert to Slack
    local slack_webhook = "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
    local payload = {
        text = string.format("🚨 Deployment Failed\\nService: %s\\nDeploy ID: %s\\nError: %s",
                            service, deploy_id, error_msg)
    }

    http.post(slack_webhook, payload)

    -- Trigger rollback workflow if needed
    if event.data.rollback_needed then
        event.dispatch("deploy.rollback", {
            deploy_id = deploy_id,
            service = service,
            reason = "automated_rollback_on_failure"
        })
    end

    return true
end
```

## Best Practices

1. **Event Naming**: Follow the `category.action` convention
2. **Data Consistency**: Always include relevant context in event data
3. **Error Handling**: Hooks should handle errors gracefully
4. **Idempotency**: Hooks should be safe to run multiple times
5. **Performance**: Keep hook execution fast; use async operations for slow tasks
6. **Security**: Validate event data before acting on it
7. **Logging**: Log all hook executions for audit trail
8. **Testing**: Test hooks with sample events before production use
