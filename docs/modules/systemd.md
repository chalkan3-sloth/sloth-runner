# ⚙️ Systemd Module

The `systemd` module provides comprehensive systemd service management functionality for Linux systems. It allows you to create, manage, and monitor systemd services programmatically.

## 🎯 Overview

The systemd module enables you to:
- Create and configure systemd service files
- Start, stop, restart, and reload services
- Enable and disable services
- Check service status and activity
- List all services
- Manage systemd daemon configuration

## 📚 Functions Overview

| Function | Description |
|----------|-------------|
| `systemd.create_service(name, config)` | Create a new systemd service |
| `systemd.start(service)` | Start a service |
| `systemd.stop(service)` | Stop a service |
| `systemd.restart(service)` | Restart a service |
| `systemd.reload(service)` | Reload a service |
| `systemd.enable(service)` | Enable service at boot |
| `systemd.disable(service)` | Disable service at boot |
| `systemd.status(service)` | Get service status |
| `systemd.is_active(service)` | Check if service is active |
| `systemd.is_enabled(service)` | Check if service is enabled |
| `systemd.daemon_reload()` | Reload systemd daemon |
| `systemd.remove_service(service)` | Remove a service |
| `systemd.list_services(opts)` | List all services |
| `systemd.show(service)` | Show detailed service info |

## 📖 Detailed Documentation

### Service Creation

#### `systemd.create_service(name, config)`

Creates a new systemd service file at `/etc/systemd/system/{name}.service`.

**Parameters:**
- `name` (string): Service name (without .service extension)
- `config` (table): Service configuration

**Configuration Options:**

```lua
{
    -- [Unit] section
    description = "Service description",
    after = "network.target",
    wants = "other.service",
    requires = "required.service",
    
    -- [Service] section (required)
    exec_start = "/path/to/executable",
    exec_stop = "/path/to/stop/script",
    exec_reload = "/path/to/reload/script",
    type = "simple",  -- simple, forking, oneshot, dbus, notify, idle
    user = "username",
    group = "groupname",
    working_directory = "/path/to/workdir",
    restart = "always",  -- no, on-success, on-failure, on-abnormal, on-abort, always
    restart_sec = "5s",
    environment = {
        VAR1 = "value1",
        VAR2 = "value2"
    },
    
    -- [Install] section
    wanted_by = "multi-user.target"
}
```

**Returns:**
- `success` (boolean): `true` if service was created
- `message` (string): Result message

**Examples:**

=== "Modern DSL"
    ```lua
    local systemd = require("systemd")
    
    local create_web_service = task("create_web_service")
        :description("Create web application service")
        :command(function(this, params)
            log.info("Creating web service...")
            
            local config = {
                description = "Web Application Server",
                after = "network.target",
                exec_start = "/usr/bin/node /app/server.js",
                type = "simple",
                user = "webapp",
                working_directory = "/app",
                restart = "always",
                restart_sec = "10s",
                environment = {
                    NODE_ENV = "production",
                    PORT = "3000"
                }
            }
            
            local success, msg = systemd.create_service("webapp", config)
            
            if success then
                log.info("✅ Service created!")
                -- Reload daemon and enable
                systemd.daemon_reload()
                systemd.enable("webapp")
                systemd.start("webapp")
                return true, "Service deployed"
            else
                log.error("❌ Failed: " .. msg)
                return false, msg
            end
        end)
        :timeout("60s")
        :build()
    
    workflow.define("deploy_service")
        :description("Deploy web application service")
        :version("1.0.0")
        :tasks({ create_web_service })
        :build()
    ```

=== "With delegate_to"
    ```lua
    local systemd = require("systemd")
    
    local deploy_remote_service = task("deploy_remote_service")
        :description("Deploy service on remote agent")
        :command(function(this, params)
            local config = {
                description = "Remote Monitoring Agent",
                after = "network.target",
                exec_start = "/opt/monitor/agent",
                type = "simple",
                user = "monitor",
                restart = "always"
            }
            
            local success, msg = systemd.create_service("monitor-agent", config)
            
            if success then
                systemd.daemon_reload()
                systemd.enable("monitor-agent")
                systemd.start("monitor-agent")
                log.info("✅ Deployed on " .. (this.agent or "local"))
                return true, "OK"
            end
            
            return false, "Failed"
        end)
        :delegate_to("production-server")
        :timeout("60s")
        :build()
    
    workflow.define("remote_deploy")
        :description("Deploy service on remote server")
        :version("1.0.0")
        :tasks({ deploy_remote_service })
        :build()
    ```

### Service Control

#### `systemd.start(service)`

Starts a systemd service.

**Parameters:**
- `service` (string): Service name

**Returns:**
- `success` (boolean), `output` (string)

**Example:**
```lua
local success, output = systemd.start("nginx")
if success then
    log.info("✅ Nginx started")
end
```

#### `systemd.stop(service)`

Stops a systemd service.

**Example:**
```lua
local success, output = systemd.stop("nginx")
```

#### `systemd.restart(service)`

Restarts a systemd service.

**Example:**
```lua
local success, output = systemd.restart("nginx")
```

#### `systemd.reload(service)`

Reloads a systemd service configuration without restarting.

**Example:**
```lua
local success, output = systemd.reload("nginx")
```

### Service Status

#### `systemd.status(service)`

Gets detailed status of a service.

**Returns:**
- `status` (string): Status output
- `error` (string): Error message if any

**Example:**
```lua
local status, err = systemd.status("nginx")
log.info("Status:\n" .. status)
```

#### `systemd.is_active(service)`

Checks if a service is currently active/running.

**Returns:**
- `active` (boolean): `true` if active
- `state` (string): Service state

**Example:**
```lua
local active, state = systemd.is_active("nginx")
if active then
    log.info("✅ Service is running")
else
    log.warn("❌ Service is " .. state)
end
```

#### `systemd.is_enabled(service)`

Checks if a service is enabled to start at boot.

**Returns:**
- `enabled` (boolean): `true` if enabled
- `state` (string): Enable state

**Example:**
```lua
local enabled, state = systemd.is_enabled("nginx")
```

### Service Management

#### `systemd.enable(service)`

Enables a service to start automatically at boot.

**Example:**
```lua
local success, output = systemd.enable("nginx")
```

#### `systemd.disable(service)`

Disables a service from starting at boot.

**Example:**
```lua
local success, output = systemd.disable("nginx")
```

#### `systemd.daemon_reload()`

Reloads systemd daemon configuration. Required after creating or modifying service files.

**Example:**
```lua
local success, output = systemd.daemon_reload()
```

#### `systemd.remove_service(service)`

Removes a systemd service completely (stops, disables, and deletes the service file).

**Example:**
```lua
local success, msg = systemd.remove_service("old-service")
```

### Service Information

#### `systemd.list_services(options)`

Lists systemd services with optional filters.

**Parameters:**
- `options` (table, optional): Filter options
  - `state`: Filter by state (e.g., "active", "failed", "inactive")
  - `no_header`: Boolean, exclude header in output

**Returns:**
- `list` (string): Service list
- `error` (string): Error if any

**Example:**
```lua
-- List all services
local list, err = systemd.list_services()
log.info(list)

-- List only active services
local active, err = systemd.list_services({ state = "active" })

-- List failed services without header
local failed, err = systemd.list_services({ 
    state = "failed", 
    no_header = true 
})
```

#### `systemd.show(service)`

Shows detailed properties of a service.

**Returns:**
- `info` (string): Detailed service information
- `error` (string): Error if any

**Example:**
```lua
local info, err = systemd.show("nginx")
log.info("Service details:\n" .. info)
```

## 🎯 Complete Examples

### Web Application Deployment

```lua
local systemd = require("systemd")

local deploy_webapp = task("deploy_webapp")
    :description("Deploy and configure web application")
    :command(function(this, params)
        log.info("🚀 Deploying web application...")
        
        -- Create service
        local config = {
            description = "Node.js Web Application",
            after = "network.target postgresql.service",
            requires = "postgresql.service",
            exec_start = "/usr/bin/node /var/www/app/server.js",
            exec_reload = "/bin/kill -HUP $MAINPID",
            type = "simple",
            user = "webapp",
            group = "webapp",
            working_directory = "/var/www/app",
            restart = "always",
            restart_sec = "10s",
            environment = {
                NODE_ENV = "production",
                PORT = "3000",
                DB_HOST = "localhost"
            },
            wanted_by = "multi-user.target"
        }
        
        local success, msg = systemd.create_service("webapp", config)
        if not success then
            return false, "Failed to create service: " .. msg
        end
        
        log.info("✅ Service file created")
        
        -- Reload daemon
        systemd.daemon_reload()
        log.info("✅ Daemon reloaded")
        
        -- Enable and start
        systemd.enable("webapp")
        log.info("✅ Service enabled")
        
        systemd.start("webapp")
        log.info("✅ Service started")
        
        -- Verify it's running
        local active, state = systemd.is_active("webapp")
        if active then
            log.info("✅ Service is running!")
            return true, "Deployment successful"
        else
            log.error("❌ Service failed to start: " .. state)
            return false, "Service not running"
        end
    end)
    :timeout("120s")
    :build()

workflow.define("deploy")
    :description("Deploy web application")
    :version("1.0.0")
    :tasks({ deploy_webapp })
    :build()
```

### Service Health Check

```lua
local systemd = require("systemd")

local health_check = task("health_check")
    :description("Check critical services health")
    :command(function(this, params)
        log.info("🔍 Health Check Starting...")
        log.info(string.rep("=", 60))
        
        local services = {
            "nginx",
            "postgresql",
            "redis",
            "webapp"
        }
        
        local all_healthy = true
        
        for _, service in ipairs(services) do
            local active, state = systemd.is_active(service)
            local enabled, enable_state = systemd.is_enabled(service)
            
            log.info("\n📦 " .. service .. ":")
            log.info("  Active: " .. (active and "✅ YES" or "❌ NO (" .. state .. ")"))
            log.info("  Enabled: " .. (enabled and "✅ YES" or "⚠️  NO"))
            
            if not active then
                all_healthy = false
                log.warn("  ⚠️  Service is not running!")
            end
        end
        
        log.info("\n" .. string.rep("=", 60))
        
        if all_healthy then
            log.info("✅ All services healthy")
            return true, "All OK"
        else
            log.error("❌ Some services are down")
            return false, "Services down"
        end
    end)
    :timeout("60s")
    :build()

workflow.define("health_check")
    :description("Check critical services health")
    :version("1.0.0")
    :tasks({ health_check })
    :build()
```

### Distributed Service Management

```lua
local systemd = require("systemd")

local restart_all_servers = task("restart_nginx")
    :description("Restart nginx on all servers")
    :command(function(this, params)
        log.info("🔄 Restarting nginx...")
        
        local success, output = systemd.restart("nginx")
        
        if success then
            -- Wait a bit for restart
            os.execute("sleep 2")
            
            -- Verify it's running
            local active, state = systemd.is_active("nginx")
            if active then
                log.info("✅ Nginx restarted on " .. (this.agent or "local"))
                return true, "OK"
            else
                log.error("❌ Nginx failed to start: " .. state)
                return false, "Failed"
            end
        end
        
        return false, "Restart failed"
    end)
    :delegate_to("web-server-1")
    :timeout("60s")
    :build()

workflow.define("rolling_restart")
    :description("Restart nginx across all servers")
    :version("1.0.0")
    :tasks({ restart_all_servers })
    :build()
```

### Service Monitoring

```lua
local systemd = require("systemd")

local monitor_services = task("monitor_services")
    :description("Monitor and report service status")
    :command(function(this, params)
        log.info("📊 Service Monitoring Report")
        log.info(string.rep("=", 60))
        
        -- List all failed services
        local failed, _ = systemd.list_services({ 
            state = "failed",
            no_header = true 
        })
        
        if failed and failed ~= "" then
            log.warn("\n⚠️  Failed Services:")
            log.warn(failed)
        else
            log.info("\n✅ No failed services")
        end
        
        -- List active services count
        local active, _ = systemd.list_services({ 
            state = "active",
            no_header = true 
        })
        
        if active then
            local count = 0
            for _ in active:gmatch("[^\r\n]+") do
                count = count + 1
            end
            log.info("\n📊 Active services: " .. count)
        end
        
        log.info("\n" .. string.rep("=", 60))
        return true, "Report complete"
    end)
    :timeout("60s")
    :build()

workflow.define("monitor")
    :description("Monitor and report service status")
    :version("1.0.0")
    :tasks({ monitor_services })
    :build()
```

### Service Update Workflow

```lua
local systemd = require("systemd")

local update_service = task("update_service")
    :description("Update service configuration")
    :command(function(this, params)
        local service_name = "webapp"
        
        log.info("🔄 Updating " .. service_name .. "...")
        
        -- Check if running
        local was_active, _ = systemd.is_active(service_name)
        
        -- Stop if running
        if was_active then
            log.info("Stopping service...")
            systemd.stop(service_name)
        end
        
        -- Update service configuration
        local new_config = {
            description = "Updated Web Application",
            after = "network.target",
            exec_start = "/usr/bin/node /app/server.js",
            type = "simple",
            user = "webapp",
            working_directory = "/app",
            restart = "always",
            environment = {
                NODE_ENV = "production",
                PORT = "3000",
                VERSION = "2.0"  -- New version
            }
        }
        
        systemd.create_service(service_name, new_config)
        systemd.daemon_reload()
        
        -- Start if it was running before
        if was_active then
            log.info("Starting service...")
            systemd.start(service_name)
            
            -- Verify
            local active, _ = systemd.is_active(service_name)
            if active then
                log.info("✅ Service updated and running")
                return true, "Updated"
            end
        end
        
        return true, "Configuration updated"
    end)
    :timeout("120s")
    :build()

workflow.define("update")
    :description("Update service configuration")
    :version("1.0.0")
    :tasks({ update_service })
    :build()
```

## 🚀 Best Practices

1. **Always reload daemon after creating/modifying services:**
   ```lua
   systemd.create_service("myservice", config)
   systemd.daemon_reload()
   ```

2. **Verify service started successfully:**
   ```lua
   systemd.start("myservice")
   local active, state = systemd.is_active("myservice")
   if not active then
       log.error("Service failed: " .. state)
   end
   ```

3. **Enable services for persistence:**
   ```lua
   systemd.enable("myservice")  -- Start at boot
   ```

4. **Use proper service types:**
   - `simple`: Default, process doesn't fork
   - `forking`: Process forks and parent exits
   - `oneshot`: Process exits before systemd continues
   - `notify`: Process sends notification when ready

5. **Set restart policies:**
   ```lua
   restart = "always"  -- Always restart
   restart_sec = "10s"  -- Wait 10s between restarts
   ```

6. **Use delegate_to for distributed management:**
   ```lua
   :delegate_to("server-name")
   ```

## ⚠️ Security Considerations

- Service files are created in `/etc/systemd/system/` (requires root/sudo)
- Always specify `user` and `group` to avoid running as root
- Use `WorkingDirectory` to isolate service environment
- Validate environment variables before setting them
- Use proper file permissions (0644 for service files)

## 🐧 Platform Support

- **Linux**: Full support (systemd-based distributions)
- **Ubuntu/Debian**: ✅ Supported
- **CentOS/RHEL**: ✅ Supported
- **Fedora**: ✅ Supported
- **Arch Linux**: ✅ Supported
- **macOS**: ❌ Not supported (use launchd instead)
- **Windows**: ❌ Not supported (use sc.exe or nssm)

## 🔗 See Also

- [exec Module](exec.md) - For running custom systemctl commands
- [Modern DSL Guide](../modern-dsl/overview.md) - DSL syntax reference
- [Distributed Agents](../distributed.md) - Remote execution with delegate_to
- [Official systemd documentation](https://www.freedesktop.org/wiki/Software/systemd/)
