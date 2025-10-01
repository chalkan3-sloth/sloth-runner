# Systemd Examples

This directory contains examples demonstrating Systemd integration with Sloth Runner.

## 📁 Available Examples

### `systemd_demo.sloth`
Comprehensive demonstration of Systemd service management capabilities.

**Features:**
- Service creation and configuration
- Service lifecycle management (start, stop, restart, enable, disable)
- Status monitoring and health checks
- Service removal and cleanup
- Blue-green deployment patterns
- Service discovery and listing

**Usage:**
```bash
sloth-runner run -f examples/systemd/systemd_demo.sloth systemd_demo_workflow
```

**What it demonstrates:**
1. 🏗️ **Service Creation**: Creates webapp and database services
2. ⚙️ **Configuration**: Sets up service files with proper settings
3. 🔄 **Lifecycle Management**: Starts, stops, restarts services
4. 📊 **Monitoring**: Checks service status and health
5. 🔄 **Blue-Green Deploy**: Demonstrates zero-downtime deployments
6. 🧹 **Cleanup**: Removes services and cleans up

## 🎯 Use Cases

### **Service Deployment**
```lua
local systemd = require("systemd")

local config = {
    description = "My Web Application",
    exec_start = "/usr/bin/node /opt/webapp/server.js",
    user = "webapp",
    restart = "always",
    environment = { NODE_ENV = "production" }
}

systemd.create_service("webapp", config)
systemd.daemon_reload()
systemd.enable("webapp")
systemd.start("webapp")
```

### **Blue-Green Deployment**
```lua
-- Deploy new version
systemd.start("webapp-green")

-- Health check new version
if systemd.is_active("webapp-green") then
    -- Switch traffic
    systemd.stop("webapp-blue")
    -- Update load balancer config
end
```

### **Service Monitoring**
```lua
local services = {"webapp", "database", "nginx"}
for _, service in ipairs(services) do
    if not systemd.is_active(service) then
        log.warn("Service " .. service .. " is down")
        systemd.restart(service)
    end
end
```

### **Bulk Operations**
```lua
local services = systemd.list_services({ state = "failed" })
for _, service in ipairs(services) do
    log.info("Restarting failed service: " .. service.name)
    systemd.restart(service.name)
end
```

## 🔧 Service Configuration Options

### **Basic Service Configuration**
```lua
local config = {
    description = "Service description",
    exec_start = "/path/to/executable",
    user = "service-user",
    group = "service-group",
    working_directory = "/opt/service",
    restart = "always",  -- no, on-success, on-failure, on-abnormal, on-watchdog, on-abort, always
    restart_sec = "5s"
}
```

### **Advanced Configuration**
```lua
local config = {
    description = "Advanced service configuration",
    exec_start = "/usr/bin/python3 /opt/app/main.py",
    exec_reload = "/bin/kill -HUP $MAINPID",
    exec_stop = "/bin/kill -TERM $MAINPID",
    user = "appuser",
    group = "appgroup",
    working_directory = "/opt/app",
    environment = {
        NODE_ENV = "production",
        DATABASE_URL = "postgresql://localhost/myapp",
        LOG_LEVEL = "info"
    },
    restart = "always",
    restart_sec = "10s",
    start_limit_interval = "60s",
    start_limit_burst = "3",
    wanted_by = "multi-user.target"
}
```

## 📋 Available Functions

| Function | Description | Example |
|----------|-------------|---------|
| `create_service(name, config)` | Create systemd service | `systemd.create_service("webapp", config)` |
| `start(service)` | Start service | `systemd.start("webapp")` |
| `stop(service)` | Stop service | `systemd.stop("webapp")` |
| `restart(service)` | Restart service | `systemd.restart("webapp")` |
| `reload(service)` | Reload service | `systemd.reload("webapp")` |
| `enable(service)` | Enable on boot | `systemd.enable("webapp")` |
| `disable(service)` | Disable on boot | `systemd.disable("webapp")` |
| `status(service)` | Get status | `systemd.status("webapp")` |
| `is_active(service)` | Check if active | `systemd.is_active("webapp")` |
| `is_enabled(service)` | Check if enabled | `systemd.is_enabled("webapp")` |
| `daemon_reload()` | Reload systemd | `systemd.daemon_reload()` |
| `remove_service(service)` | Remove service | `systemd.remove_service("webapp")` |
| `list_services(options)` | List services | `systemd.list_services({state="active"})` |
| `show(service)` | Show properties | `systemd.show("webapp")` |

## ⚙️ Prerequisites

- Linux system with systemd
- Appropriate permissions (usually requires sudo for service management)
- Sloth Runner compiled and available

## 🛡️ Security Considerations

### **User Isolation**
```lua
local config = {
    user = "webapp",           -- Run as non-root user
    group = "webapp",          -- Specific group
    no_new_privileges = true,  -- Security hardening
    private_tmp = true,        -- Isolated /tmp
    protect_system = "strict"  -- Read-only filesystem
}
```

### **Resource Limits**
```lua
local config = {
    memory_limit = "512M",
    cpu_quota = "50%",
    tasks_max = 100
}
```

## 🚀 Getting Started

1. **Run the demo:**
   ```bash
   sudo sloth-runner run -f examples/systemd/systemd_demo.sloth systemd_demo_workflow
   ```

2. **Check created services:**
   ```bash
   systemctl status sloth-webapp
   systemctl status sloth-database
   ```

3. **Monitor logs:**
   ```bash
   journalctl -u sloth-webapp -f
   ```

## 🔧 Customization

You can customize the examples by:
- Modifying service configurations
- Adding custom health checks
- Implementing rolling deployments
- Adding service dependencies
- Creating service templates for different application types

## 📚 Related Documentation

- [Systemd Module Documentation](../../docs/modules/systemd.md)
- [Service Management Guide](../../docs/guides/service-management.md)
- [Deployment Patterns](../../docs/patterns/deployment.md)