# ⏰ Advanced Task Scheduler

Sloth Runner includes a **powerful task scheduling system** that enables automated execution of workflows with cron-style syntax, background daemon support, and comprehensive schedule management.

## 🚀 Quick Start

### Enable the Scheduler
```bash
# Start scheduler daemon
sloth-runner scheduler enable --config scheduler.yaml

# Start with custom configuration
sloth-runner scheduler enable --config /path/to/custom-scheduler.yaml
```

### Basic Schedule Configuration
```yaml
# scheduler.yaml
scheduler:
  enabled: true
  timezone: "UTC"
  max_concurrent_tasks: 5
  log_level: "info"
  
schedules:
  - name: "daily_backup"
    cron: "0 2 * * *"  # Every day at 2 AM
    workflow: "backup.lua"
    description: "Daily database backup"
    
  - name: "hourly_health_check"
    cron: "0 * * * *"  # Every hour
    workflow: "health-check.lua"
    timeout: "5m"
    
  - name: "weekly_cleanup"
    cron: "0 0 * * 0"  # Every Sunday at midnight
    workflow: "cleanup.lua"
    enabled: true
```

## 📅 Cron Syntax

### Standard Cron Format
```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
* * * * *
```

### Common Examples
```yaml
schedules:
  # Every minute
  - cron: "* * * * *"
    
  # Every 5 minutes
  - cron: "*/5 * * * *"
    
  # Every hour at minute 30
  - cron: "30 * * * *"
    
  # Every day at 3:30 AM
  - cron: "30 3 * * *"
    
  # Every Monday at 9 AM
  - cron: "0 9 * * 1"
    
  # First day of every month at midnight
  - cron: "0 0 1 * *"
    
  # Every 15 minutes during business hours (9-5, Mon-Fri)
  - cron: "*/15 9-17 * * 1-5"
```

### Extended Syntax
```yaml
schedules:
  # Using @yearly, @monthly, @weekly, @daily, @hourly
  - cron: "@daily"
    workflow: "daily-tasks.lua"
    
  - cron: "@weekly"
    workflow: "weekly-report.lua"
    
  # Using @every with duration
  - cron: "@every 30m"
    workflow: "monitoring.lua"
    
  - cron: "@every 2h30m"
    workflow: "periodic-sync.lua"
```

## 🔧 Advanced Configuration

### Complete Scheduler Config
```yaml
# advanced-scheduler.yaml
scheduler:
  enabled: true
  timezone: "America/New_York"
  max_concurrent_tasks: 10
  log_level: "debug"
  
  # Database settings
  database:
    path: "/data/scheduler.db"
    backup_interval: "24h"
    
  # Notification settings
  notifications:
    on_failure: true
    on_success: false
    channels:
      - type: "slack"
        webhook: "https://hooks.slack.com/..."
      - type: "email"
        smtp: "smtp.company.com:587"
        
  # Performance settings
  performance:
    worker_pool_size: 5
    queue_buffer_size: 100
    execution_timeout: "30m"

schedules:
  - name: "production_deployment"
    cron: "0 2 * * 1-5"  # Weekdays at 2 AM
    workflow: "deploy-prod.lua"
    description: "Production deployment pipeline"
    
    # Schedule-specific settings
    timeout: "45m"
    retry_attempts: 3
    retry_delay: "5m"
    
    # Environment variables for this schedule
    env:
      ENVIRONMENT: "production"
      DEPLOY_STRATEGY: "blue_green"
      
    # Only run if conditions are met
    conditions:
      - "production.health_check == 'healthy'"
      - "staging.last_deploy > '24h'"
      
    # Notification overrides
    notifications:
      on_success: true
      on_failure: true
      channels: ["slack", "email"]
```

## 🎯 Workflow Integration

### Scheduled Workflow Example
```lua
-- scheduled-backup.lua
local backup = task("database_backup")
    :description("Automated database backup")
    :command(function(params, deps)
        local timestamp = os.date("%Y%m%d_%H%M%S")
        local backup_file = "/backups/db_" .. timestamp .. ".sql"
        
        log.info("Starting scheduled backup to: " .. backup_file)
        
        -- Perform backup
        local result = exec.run("pg_dump myapp_db > " .. backup_file)
        if not result.success then
            return false, "Backup failed: " .. result.stderr
        end
        
        -- Compress backup
        exec.run("gzip " .. backup_file)
        
        -- Upload to cloud storage
        local aws = require("aws")
        aws.s3.upload(backup_file .. ".gz", "backups-bucket")
        
        -- Clean old backups (keep last 7 days)
        exec.run("find /backups -name '*.gz' -mtime +7 -delete")
        
        return true, "Backup completed successfully", {
            backup_file = backup_file .. ".gz",
            size = fs.size(backup_file .. ".gz"),
            timestamp = timestamp
        }
    end)
    :timeout("30m")
    :on_failure(function(params, error)
        -- Send alert on failure
        local notifications = require("notifications")
        notifications.slack.send({
            channel = "#alerts",
            message = "⚠️ Scheduled backup failed: " .. error,
            color = "danger"
        })
    end)
    :build()

-- Export for scheduler
workflow.define("backup_workflow", {
    description = "Scheduled database backup workflow",
    tasks = { backup },
    
    -- Scheduler metadata
    schedule_metadata = {
        category = "maintenance",
        priority = "high",
        max_duration = "30m"
    }
})
```

### Conditional Scheduling
```lua
-- conditional-deploy.lua
local deployment = task("conditional_deployment")
    :description("Deploy only if tests pass and load is low")
    :command(function(params, deps)
        -- Check prerequisites
        local health = monitoring.get_health_status()
        local load = monitoring.get_system_load()
        local tests = state.get("last_test_results")
        
        -- Conditional logic
        if tests.status ~= "passed" then
            return false, "Tests not passing, skipping deployment"
        end
        
        if load.cpu > 0.8 then
            return false, "System load too high, deferring deployment"
        end
        
        if health.status ~= "healthy" then
            return false, "System health check failed"
        end
        
        -- Proceed with deployment
        log.info("Conditions met, proceeding with deployment")
        return exec.run("./deploy.sh")
    end)
    :build()
```

## 📊 Management Commands

### Scheduler Control
```bash
# Enable/disable scheduler
sloth-runner scheduler enable
sloth-runner scheduler disable

# List all scheduled tasks
sloth-runner scheduler list

# Show detailed schedule information
sloth-runner scheduler show backup_task

# Delete a scheduled task
sloth-runner scheduler delete old_cleanup

# Pause/resume specific schedule
sloth-runner scheduler pause weekly_report
sloth-runner scheduler resume weekly_report
```

### Schedule Status
```bash
# Get scheduler status
sloth-runner scheduler status

# Output:
# Scheduler Status: RUNNING
# Active Schedules: 12
# Next Execution: daily_backup in 2h 15m
# Last Execution: hourly_check (success) 45m ago
# Failed Tasks (24h): 1
```

### Execution History
```bash
# View execution history
sloth-runner scheduler history --limit 20

# Filter by schedule name
sloth-runner scheduler history --schedule backup_task

# Filter by status
sloth-runner scheduler history --status failed --last 7d
```

## 🔍 Monitoring & Logging

### Execution Logs
```bash
# View scheduler logs
tail -f /var/log/sloth-runner/scheduler.log

# Sample log output:
# 2024-01-15 02:00:00 INFO [scheduler] Starting execution: daily_backup
# 2024-01-15 02:00:01 INFO [backup_task] Database backup initiated
# 2024-01-15 02:05:23 INFO [backup_task] Backup completed: 2.4GB
# 2024-01-15 02:05:24 INFO [scheduler] Execution completed: daily_backup (success)
```

### Metrics Collection
```yaml
# Enable metrics in scheduler config
scheduler:
  metrics:
    enabled: true
    port: 9090
    path: "/metrics"
    
  # Prometheus metrics available:
  # - scheduler_executions_total
  # - scheduler_execution_duration_seconds
  # - scheduler_queue_size
  # - scheduler_worker_pool_utilization
```

## 🚨 Error Handling & Alerts

### Failure Notifications
```yaml
schedules:
  - name: "critical_task"
    cron: "0 */6 * * *"
    workflow: "critical.lua"
    
    # Failure handling
    on_failure:
      retry:
        attempts: 3
        delay: "10m"
        backoff: "exponential"
        
      notifications:
        - type: "slack"
          webhook: "https://hooks.slack.com/..."
          message: "🚨 Critical task failed: {{.Error}}"
          
        - type: "email"
          to: ["ops@company.com"]
          subject: "URGENT: Critical Task Failure"
          
        - type: "webhook"
          url: "https://monitoring.company.com/alert"
          method: "POST"
```

### Circuit Breaker Pattern
```lua
-- circuit-breaker-task.lua
local task_with_circuit_breaker = task("resilient_task")
    :description("Task with circuit breaker protection")
    :command(function(params, deps)
        local circuit_breaker = require("reliability").circuit_breaker
        
        -- Configure circuit breaker
        local breaker = circuit_breaker.new({
            failure_threshold = 5,    -- Open after 5 failures
            recovery_timeout = "5m",  -- Try to close after 5 minutes
            success_threshold = 3     -- Close after 3 successes
        })
        
        return breaker.call(function()
            -- Your actual task logic here
            return external_api.call()
        end)
    end)
    :build()
```

## 🔒 Security Features

### Secure Configuration
```yaml
scheduler:
  security:
    # Run as specific user
    run_as_user: "scheduler"
    run_as_group: "schedulers"
    
    # Restrict file access
    working_directory: "/var/lib/sloth-runner"
    allowed_paths:
      - "/var/lib/sloth-runner"
      - "/tmp/scheduler"
      
    # Environment isolation
    clear_env: true
    allowed_env_vars:
      - "PATH"
      - "HOME"
      - "TZ"
```

### Audit Logging
```yaml
scheduler:
  audit:
    enabled: true
    log_file: "/var/log/sloth-runner/audit.log"
    include_payload: false  # Don't log sensitive data
    
    # Audit events
    events:
      - "schedule_created"
      - "schedule_modified" 
      - "schedule_deleted"
      - "execution_started"
      - "execution_completed"
      - "execution_failed"
```

## ⚡ Performance Optimization

### Resource Management
```yaml
scheduler:
  resources:
    # Memory limits per task
    max_memory_per_task: "512MB"
    
    # CPU limits
    max_cpu_per_task: "1.0"  # 1 CPU core
    
    # Execution time limits
    default_timeout: "15m"
    max_timeout: "2h"
    
    # Disk space monitoring
    min_free_disk: "1GB"
    cleanup_threshold: "5GB"
```

### Queue Management
```yaml
scheduler:
  queue:
    # Queue sizes
    pending_queue_size: 1000
    running_queue_size: 50
    
    # Priority levels
    priority_levels: 5
    default_priority: 3
    
    # Queue processing
    batch_size: 10
    processing_interval: "1s"
```

## 🎯 Best Practices

### Schedule Design
1. **Use appropriate timeouts** for long-running tasks
2. **Implement idempotent** workflows 
3. **Add proper error handling** and retries
4. **Monitor resource usage** regularly
5. **Use conditional execution** when appropriate

### Performance Tips
1. **Avoid overlapping schedules** for resource-intensive tasks
2. **Use priority levels** for critical vs. maintenance tasks
3. **Implement circuit breakers** for external dependencies
4. **Monitor queue depth** and adjust worker pool size
5. **Regular cleanup** of old execution logs

### Security Recommendations
1. **Run scheduler as dedicated user** with minimal privileges
2. **Validate all inputs** from scheduled workflows
3. **Use secure storage** for sensitive configuration
4. **Enable audit logging** for compliance
5. **Regular security reviews** of scheduled tasks

---

The Advanced Task Scheduler transforms Sloth Runner into a **comprehensive automation platform** capable of managing complex, time-based workflows with enterprise-grade reliability and monitoring! ⏰🚀