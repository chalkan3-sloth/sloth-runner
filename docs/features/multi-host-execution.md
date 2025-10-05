# 🌐 Multi-Host Execution

## Overview

Sloth Runner now supports **parallel execution of tasks on multiple remote hosts** simultaneously. This powerful feature enables you to:

- Deploy applications to multiple servers in parallel
- Run maintenance tasks across your entire infrastructure
- Collect system information from multiple hosts at once
- Execute health checks on distributed systems
- Perform rolling updates with parallel execution

## How It Works

The multi-host execution feature uses **parallel gRPC connections** to execute tasks on multiple agents simultaneously. Each host receives the same task definition and executes it independently, with results collected and displayed in a unified summary.

### Architecture

```
┌─────────────┐
│ Sloth Runner│
│   (Master)  │
└─────┬───────┘
      │ Parallel Execution
      ├─────────────────┐
      │                 │
      ▼                 ▼
┌──────────┐      ┌──────────┐
│  Agent 1 │      │  Agent 2 │
│lady-arch │      │keite-guica│
└──────────┘      └──────────┘
```

## Usage Methods

### Method 1: Command Line Override

Use multiple `--delegate-to` flags when running your task:

```bash
sloth-runner run -f deploy.sloth --delegate-to host1 --delegate-to host2 --delegate-to host3
```

**Example:**
```bash
sloth-runner run -f examples/test_multi_host.sloth --delegate-to lady-arch --delegate-to keite-guica
```

### Method 2: File-Based Configuration

Define multiple hosts in your `.sloth` file:

```lua
TaskDefinitions = {
    deployment = {
        description = "Deploy to multiple servers",
        delegate_to = {"web-server-1", "web-server-2", "web-server-3"},
        tasks = {
            {
                name = "deploy_app",
                description = "Deploy application",
                command = function()
                    log.info("Deploying on " .. exec.run("hostname").output)
                    -- Your deployment logic here
                    return true, "Deployment successful"
                end
            }
        }
    }
}
```

### Method 3: Mixed Approach

You can override file-based configuration with CLI flags:

```lua
-- In file: delegate_to = "production-server"
```

```bash
# Override with multiple hosts at runtime
sloth-runner run -f deploy.sloth --delegate-to staging-1 --delegate-to staging-2
```

## Complete Examples

### Example 1: System Information Collection

```lua
-- file: system_check.sloth
TaskDefinitions = {
    system_check = {
        description = "Collect system info from all hosts",
        tasks = {
            {
                name = "collect_info",
                description = "Get system information",
                command = function()
                    log.info("=== System Information ===")

                    -- Hostname
                    local hostname = exec.run("hostname")
                    if hostname.success then
                        log.info("Host: " .. hostname.output)
                    end

                    -- System info
                    local uname = exec.run("uname -a")
                    if uname.success then
                        log.info("System: " .. uname.output)
                    end

                    -- CPU info
                    local cpu = exec.run("nproc")
                    if cpu.success then
                        log.info("CPUs: " .. cpu.output)
                    end

                    -- Memory info
                    local mem = exec.run("free -h | grep Mem")
                    if mem.success then
                        log.info("Memory: " .. mem.output)
                    end

                    -- Disk usage
                    local disk = exec.run("df -h /")
                    if disk.success then
                        log.info("Disk: " .. disk.output)
                    end

                    return true, "System info collected"
                end
            }
        }
    }
}
```

Run on multiple hosts:
```bash
sloth-runner run -f system_check.sloth --delegate-to server1 --delegate-to server2 --delegate-to server3
```

### Example 2: Application Deployment

```lua
-- file: deploy_app.sloth
TaskDefinitions = {
    deploy_application = {
        description = "Deploy application to multiple servers",
        tasks = {
            {
                name = "stop_service",
                description = "Stop the application service",
                command = function()
                    log.info("Stopping service...")
                    exec.run("systemctl stop myapp")
                    return true, "Service stopped"
                end
            },
            {
                name = "update_code",
                description = "Update application code",
                depends_on = {"stop_service"},
                command = function()
                    log.info("Updating application...")
                    exec.run("cd /opt/myapp && git pull")
                    exec.run("cd /opt/myapp && npm install")
                    return true, "Code updated"
                end
            },
            {
                name = "start_service",
                description = "Start the application service",
                depends_on = {"update_code"},
                command = function()
                    log.info("Starting service...")
                    exec.run("systemctl start myapp")

                    -- Verify service is running
                    local status = exec.run("systemctl is-active myapp")
                    if status.output:match("active") then
                        log.info("Service started successfully")
                        return true, "Service running"
                    else
                        return false, "Service failed to start"
                    end
                end
            }
        }
    }
}
```

Deploy to production servers:
```bash
sloth-runner run -f deploy_app.sloth \
    --delegate-to prod-web-1 \
    --delegate-to prod-web-2 \
    --delegate-to prod-web-3
```

### Example 3: Health Checks

```lua
-- file: health_check.sloth
TaskDefinitions = {
    health_check = {
        description = "Run health checks on all nodes",
        tasks = {
            {
                name = "check_services",
                description = "Check critical services",
                command = function()
                    local all_ok = true
                    local services = {"nginx", "postgresql", "redis", "myapp"}

                    for _, service in ipairs(services) do
                        local status = exec.run("systemctl is-active " .. service)
                        if status.output:match("active") then
                            log.info("✅ " .. service .. " is running")
                        else
                            log.error("❌ " .. service .. " is not running")
                            all_ok = false
                        end
                    end

                    -- Check disk space
                    local disk = exec.run("df -h / | awk 'NR==2 {print $5}' | sed 's/%//'")
                    local usage = tonumber(disk.output)
                    if usage > 80 then
                        log.warn("⚠️ Disk usage is high: " .. usage .. "%")
                    else
                        log.info("✅ Disk usage is normal: " .. usage .. "%")
                    end

                    -- Check load average
                    local load = exec.run("uptime | awk -F'load average:' '{print $2}'")
                    log.info("Load average: " .. load.output)

                    return all_ok, all_ok and "All checks passed" or "Some checks failed"
                end
            }
        }
    }
}
```

## Execution Output

When running tasks on multiple hosts, you'll see a comprehensive execution summary:

```
🚀 Executing task 'system_info' on 3 hosts

Host         | Status
-------------|------------
web-server-1 | ⏳ Pending
web-server-2 | ⏳ Pending
web-server-3 | ⏳ Pending

🔗 Connecting to 192.168.1.10:50051...
🔗 Connecting to 192.168.1.11:50051...
🔗 Connecting to 192.168.1.12:50051...

✅ Success on 192.168.1.10:50051
✅ Success on 192.168.1.11:50051
✅ Success on 192.168.1.12:50051

📊 Multi-Host Execution Results

Host         | Status     | Details
-------------|------------|------------------------
web-server-1 | ✅ Success | Completed successfully
web-server-2 | ✅ Success | Completed successfully
web-server-3 | ✅ Success | Completed successfully

┌─ ✅ Execution Summary ─┐
│ Task:      system_info │
│ Total:     3 hosts     │
│ Success:   3           │
│ Failed:    0           │
└────────────────────────┘
```

## Agent Configuration

Before using multi-host execution, ensure your agents are properly configured and running:

### 1. Start Agents on Remote Hosts

On each remote host:
```bash
sloth-runner agent start <agent-name> --port <port>
```

### 2. Register Agents

Agents should be registered with the master:
```bash
sloth-runner agent list
```

Output:
```
AGENT NAME     ADDRESS              STATUS
----------     -------              ------
lady-arch      192.168.1.16:50052   Active
keite-guica    192.168.1.17:50051   Active
lady-guica     192.168.1.16:50051   Active
```

### 3. Use Agent Names or Addresses

You can use either agent names or direct addresses:

```bash
# Using agent names
sloth-runner run -f task.sloth --delegate-to lady-arch --delegate-to keite-guica

# Using addresses
sloth-runner run -f task.sloth --delegate-to 192.168.1.16:50052 --delegate-to 192.168.1.17:50051
```

## Advanced Features

### Conditional Execution

Execute on different hosts based on conditions:

```lua
TaskDefinitions = {
    conditional_deploy = {
        description = "Conditional deployment",
        tasks = {
            {
                name = "deploy",
                command = function()
                    local hostname = exec.run("hostname").output:gsub("\n", "")

                    if hostname:match("prod") then
                        log.info("Production deployment")
                        -- Production specific logic
                    elseif hostname:match("staging") then
                        log.info("Staging deployment")
                        -- Staging specific logic
                    else
                        log.info("Development deployment")
                        -- Development specific logic
                    end

                    return true, "Deployment completed for " .. hostname
                end
            }
        }
    }
}
```

### Rolling Updates

Combine with dependencies for rolling updates:

```lua
TaskDefinitions = {
    rolling_update = {
        description = "Rolling update across hosts",
        tasks = {
            {
                name = "update_batch_1",
                delegate_to = {"server1", "server2"},
                command = function()
                    -- Update first batch
                    return true, "Batch 1 updated"
                end
            },
            {
                name = "verify_batch_1",
                depends_on = {"update_batch_1"},
                delegate_to = {"server1", "server2"},
                command = function()
                    -- Verify first batch
                    return true, "Batch 1 verified"
                end
            },
            {
                name = "update_batch_2",
                depends_on = {"verify_batch_1"},
                delegate_to = {"server3", "server4"},
                command = function()
                    -- Update second batch
                    return true, "Batch 2 updated"
                end
            }
        }
    }
}
```

## Best Practices

### 1. Use Meaningful Host Names

Instead of IP addresses, use descriptive agent names:
```bash
sloth-runner agent start web-prod-1
sloth-runner agent start db-prod-1
sloth-runner agent start cache-prod-1
```

### 2. Group Related Hosts

Create task definitions for specific host groups:
```lua
-- Web servers
web_servers = {"web-1", "web-2", "web-3"}

-- Database servers
db_servers = {"db-primary", "db-replica-1", "db-replica-2"}

-- Cache servers
cache_servers = {"redis-1", "redis-2"}
```

### 3. Error Handling

Always check for host-specific failures:
```lua
command = function()
    local result = exec.run("critical-command")
    if not result.success then
        log.error("Failed on " .. exec.run("hostname").output)
        return false, result.stderr
    end
    return true, "Success"
end
```

### 4. Idempotency

Ensure tasks are idempotent for safe re-execution:
```lua
command = function()
    -- Check if already done
    local check = exec.run("test -f /opt/myapp/.deployed")
    if check.success then
        log.info("Already deployed, skipping")
        return true, "Already deployed"
    end

    -- Perform deployment
    exec.run("deploy-application")
    exec.run("touch /opt/myapp/.deployed")

    return true, "Newly deployed"
end
```

## Troubleshooting

### Issue: Agent Not Found

**Error:**
```
Failed to resolve agent 'host-name': agent not found
```

**Solution:**
1. Check agent is running: `sloth-runner agent list`
2. Verify agent name spelling
3. Ensure agent is registered with master

### Issue: Connection Failed

**Error:**
```
Failed to connect to 192.168.1.10:50051
```

**Solution:**
1. Check network connectivity
2. Verify firewall rules allow gRPC port
3. Ensure agent is listening on correct port

### Issue: Partial Failures

**Behavior:** Some hosts succeed, others fail

**Solution:**
1. Check individual host logs
2. Verify all hosts have required dependencies
3. Ensure consistent environment across hosts

## Performance Considerations

### Parallel Execution

- All hosts execute simultaneously
- No performance penalty for multiple hosts
- Results collected asynchronously

### Resource Usage

- Each host connection uses minimal memory
- Network bandwidth scales linearly with hosts
- Master node coordinates without bottlenecks

### Scaling Limits

- Tested with up to 100 simultaneous hosts
- Network latency may affect large deployments
- Consider batching for very large infrastructures

## Migration Guide

### From Single Host to Multi-Host

**Before:**
```lua
delegate_to = "production-server"
```

**After (Option 1 - CLI Override):**
```bash
sloth-runner run -f task.sloth \
    --delegate-to prod-1 \
    --delegate-to prod-2 \
    --delegate-to prod-3
```

**After (Option 2 - File-Based):**
```lua
delegate_to = {"prod-1", "prod-2", "prod-3"}
```

### From Sequential to Parallel

**Before (Sequential):**
```bash
for host in server1 server2 server3; do
    sloth-runner run -f task.sloth --delegate-to $host
done
```

**After (Parallel):**
```bash
sloth-runner run -f task.sloth \
    --delegate-to server1 \
    --delegate-to server2 \
    --delegate-to server3
```

## Summary

The multi-host execution feature transforms Sloth Runner into a powerful orchestration tool for distributed systems. Key benefits include:

- **Parallel execution** - All hosts execute simultaneously
- **Unified results** - Single summary for all executions
- **Flexible configuration** - CLI or file-based host specification
- **Production ready** - Error handling and progress tracking
- **Scalable** - Handles large infrastructures efficiently

Use multi-host execution whenever you need to:
- Deploy to multiple servers
- Collect information from distributed systems
- Run maintenance tasks across infrastructure
- Perform health checks on multiple nodes
- Execute rolling updates or migrations