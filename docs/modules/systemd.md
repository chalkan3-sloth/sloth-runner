# Systemd Module Documentation

The Systemd module provides comprehensive systemd service management functionality, allowing you to create, manage, and monitor systemd services directly from your Sloth Runner workflows.

## 📋 **Overview**

The Systemd module enables you to:
- Create systemd service files with full configuration
- Start, stop, restart, reload services
- Enable and disable services for boot
- Check service status and health
- List and monitor services
- Clean up and remove services

## 🚀 **Quick Start**

### Basic Service Creation and Management

```lua
local systemd = require("systemd")

-- Create a service
local service_config = {
    description = "My Application",
    exec_start = "/usr/bin/myapp",
    user = "myuser",
    restart = "always"
}

local success, msg = systemd.create_service("myapp", service_config)
if success then
    systemd.daemon_reload()
    systemd.enable("myapp")
    systemd.start("myapp")
end
```

## 🔧 **API Reference**

### `systemd.create_service(name, config)`

Creates a systemd service file with the specified configuration.

**Parameters:**
- `name` (string): Service name (without .service extension)
- `config` (table): Service configuration options

**Configuration Options:**
```lua
{
    -- [Unit] section
    description = "Service description",
    after = "network.target",
    wants = "network-online.target",
    requires = "postgresql.service",
    
    -- [Service] section
    exec_start = "/path/to/executable",      -- Required
    exec_stop = "/path/to/stop/command",
    exec_reload = "/bin/kill -USR1 $MAINPID",
    type = "simple",                         -- simple, forking, oneshot, notify, etc.
    user = "username",
    group = "groupname", 
    working_directory = "/path/to/workdir",
    restart = "always",                      -- always, on-failure, no, etc.
    restart_sec = "10",
    
    -- Environment variables
    environment = {
        NODE_ENV = "production",
        PORT = "3000"
    },
    
    -- [Install] section
    wanted_by = "multi-user.target"
}
```

**Returns:**
- `success` (boolean): True if service file was created
- `message` (string): Success message or error details

**Example:**
```lua
local config = {
    description = "Node.js Web Application",
    after = "network.target",
    exec_start = "/usr/bin/node /opt/webapp/server.js",
    type = "simple",
    user = "nodejs",
    working_directory = "/opt/webapp",
    restart = "always",
    restart_sec = "5",
    environment = {
        NODE_ENV = "production",
        PORT = "8080"
    }
}

local success, msg = systemd.create_service("webapp", config)
```

### Service Control Functions

#### `systemd.start(service_name)`
Starts a systemd service.

#### `systemd.stop(service_name)`
Stops a systemd service.

#### `systemd.restart(service_name)`
Restarts a systemd service.

#### `systemd.reload(service_name)`
Reloads a systemd service (if supported by the service).

#### `systemd.enable(service_name)`
Enables a service to start at boot.

#### `systemd.disable(service_name)`
Disables a service from starting at boot.

**All control functions return:**
- `success` (boolean): True if operation succeeded
- `output` (string): Command output or error message

### Service Status Functions

#### `systemd.status(service_name)`
Gets detailed status information for a service.

**Returns:**
- `output` (string): Full status output
- `error` (string or nil): Error message if any

#### `systemd.is_active(service_name)`
Checks if a service is currently active.

**Returns:**
- `active` (boolean): True if service is active
- `status` (string): Status string ("active", "inactive", etc.)

#### `systemd.is_enabled(service_name)`
Checks if a service is enabled for boot.

**Returns:**
- `enabled` (boolean): True if service is enabled
- `status` (string): Status string ("enabled", "disabled", etc.)

### System Management Functions

#### `systemd.daemon_reload()`
Reloads systemd daemon configuration.

**Returns:**
- `success` (boolean): True if reload succeeded
- `output` (string): Command output

#### `systemd.remove_service(service_name)`
Stops, disables, and removes a service file.

**Returns:**
- `success` (boolean): True if removal succeeded
- `message` (string): Success or error message

#### `systemd.list_services(options)`
Lists systemd services with optional filtering.

**Parameters:**
- `options` (table, optional): Filtering options
  - `state` (string): Filter by state ("active", "failed", etc.)
  - `no_header` (boolean): Omit header from output

**Returns:**
- `output` (string): Service listing output
- `error` (string or nil): Error message if any

#### `systemd.show(service_name)`
Shows detailed service properties.

**Returns:**
- `output` (string): Detailed service information
- `error` (string or nil): Error message if any

## 💡 **Complete Examples**

### Example 1: Deploy Web Application

```lua
local deploy_webapp = task("deploy_webapp")
    :description("Deploy web application as systemd service")
    :command(function()
        local systemd = require("systemd")
        
        -- Service configuration
        local config = {
            description = "Production Web Application",
            after = "network.target postgresql.service",
            exec_start = "/opt/webapp/bin/server",
            exec_reload = "/bin/kill -USR1 $MAINPID",
            type = "simple",
            user = "webapp",
            group = "webapp",
            working_directory = "/opt/webapp",
            restart = "always",
            restart_sec = "10",
            environment = {
                NODE_ENV = "production",
                DATABASE_URL = "postgresql://localhost/webapp",
                PORT = "8080"
            }
        }
        
        -- Create and start service
        local success, msg = systemd.create_service("webapp", config)
        if not success then
            return false, "Failed to create service: " .. msg
        end
        
        systemd.daemon_reload()
        systemd.enable("webapp")
        systemd.start("webapp")
        
        -- Verify service is running
        local is_active = systemd.is_active("webapp")
        if not is_active then
            return false, "Service failed to start"
        end
        
        return true, "Web application deployed successfully"
    end)
    :build()
```

### Example 2: Database Service Management

```lua
local manage_database = task("manage_database")
    :description("Manage database service")
    :command(function()
        local systemd = require("systemd")
        
        local db_config = {
            description = "PostgreSQL Database Server",
            after = "network.target",
            exec_start = "/usr/bin/postgres -D /var/lib/postgresql/data",
            exec_reload = "/bin/kill -HUP $MAINPID",
            type = "forking",
            user = "postgres",
            group = "postgres",
            working_directory = "/var/lib/postgresql",
            restart = "always",
            restart_sec = "5"
        }
        
        -- Create database service
        systemd.create_service("custom-postgres", db_config)
        systemd.daemon_reload()
        
        -- Enable and start
        systemd.enable("custom-postgres")
        systemd.start("custom-postgres")
        
        -- Wait for startup
        for i = 1, 10 do
            local is_active = systemd.is_active("custom-postgres")
            if is_active then
                log.info("✅ Database service is active")
                break
            end
            log.info("⏳ Waiting for database startup...")
            os.execute("sleep 2")
        end
        
        return true, "Database service configured"
    end)
    :build()
```

### Example 3: Service Monitoring and Health Checks

```lua
local monitor_services = task("monitor_services")
    :description("Monitor critical services")
    :command(function()
        local systemd = require("systemd")
        
        local critical_services = {"webapp", "database", "nginx", "redis"}
        local failed_services = {}
        
        for _, service in ipairs(critical_services) do
            local is_active, status = systemd.is_active(service)
            local is_enabled = systemd.is_enabled(service)
            
            log.info("Service: " .. service)
            log.info("  Active: " .. tostring(is_active) .. " (" .. status .. ")")
            log.info("  Enabled: " .. tostring(is_enabled))
            
            if not is_active then
                table.insert(failed_services, service)
                
                -- Try to restart failed services
                log.warn("🔄 Attempting to restart " .. service)
                local restart_ok = systemd.restart(service)
                
                if restart_ok then
                    log.info("✅ Successfully restarted " .. service)
                else
                    log.error("❌ Failed to restart " .. service)
                end
            end
        end
        
        -- Check for any failed services system-wide
        local failed_list, error = systemd.list_services({state = "failed"})
        if not error and failed_list and failed_list ~= "" then
            log.warn("⚠️ System has failed services:")
            log.warn(failed_list)
        end
        
        return #failed_services == 0, 
               #failed_services == 0 and "All services healthy" or "Some services failed",
               { failed_services = failed_services }
    end)
    :build()
```

### Example 4: Blue-Green Deployment

```lua
local blue_green_deploy = task("blue_green_deploy")
    :description("Blue-green deployment using systemd")
    :command(function()
        local systemd = require("systemd")
        local state = require("state")
        
        -- Determine current and new colors
        local current_color = state.get("active_color") or "blue"
        local new_color = current_color == "blue" and "green" or "blue"
        
        log.info("Current active: " .. current_color)
        log.info("Deploying to: " .. new_color)
        
        -- Configure new service
        local new_config = {
            description = "Web App " .. string.upper(new_color),
            exec_start = "/opt/webapp-" .. new_color .. "/bin/server",
            user = "webapp",
            working_directory = "/opt/webapp-" .. new_color,
            restart = "always",
            environment = {
                NODE_ENV = "production",
                PORT = new_color == "blue" and "8080" or "8081"
            }
        }
        
        -- Deploy new version
        systemd.create_service("webapp-" .. new_color, new_config)
        systemd.daemon_reload()
        systemd.enable("webapp-" .. new_color)
        systemd.start("webapp-" .. new_color)
        
        -- Verify new service is healthy
        local is_active = systemd.is_active("webapp-" .. new_color)
        if not is_active then
            return false, "New service failed to start"
        end
        
        -- Switch traffic (update load balancer config here)
        log.info("🔄 Switching traffic to " .. new_color)
        
        -- Stop old service
        systemd.stop("webapp-" .. current_color)
        systemd.disable("webapp-" .. current_color)
        
        -- Update state
        state.set("active_color", new_color)
        
        return true, "Blue-green deployment completed", {
            previous_color = current_color,
            active_color = new_color
        }
    end)
    :build()
```

## 🛡️ **Best Practices**

### 1. Always Use daemon-reload After Creating Services

```lua
systemd.create_service("myapp", config)
systemd.daemon_reload()  -- Required!
systemd.enable("myapp")
```

### 2. Check Service Status Before Operations

```lua
local is_active = systemd.is_active("myapp")
if is_active then
    log.info("Service is already running")
else
    systemd.start("myapp")
end
```

### 3. Handle Errors Gracefully

```lua
local success, msg = systemd.start("myapp")
if not success then
    log.error("Failed to start service: " .. msg)
    
    -- Check what went wrong
    local status, error = systemd.status("myapp")
    if error then
        log.error("Status check failed: " .. error)
    else
        log.info("Service status: " .. status)
    end
    
    return false, "Service startup failed"
end
```

### 4. Use Proper Service Types

```lua
-- For simple long-running processes
type = "simple"

-- For services that fork
type = "forking"

-- For one-time execution
type = "oneshot"

-- For services that notify systemd when ready
type = "notify"
```

### 5. Set Appropriate Restart Policies

```lua
-- Always restart on exit
restart = "always"

-- Only restart on failure
restart = "on-failure"

-- Never restart
restart = "no"

-- Restart delay
restart_sec = "10"
```

## 🔗 **Integration with Other Modules**

### With State Module (Service Coordination)

```lua
local state = require("state")
local systemd = require("systemd")

-- Wait for dependency service
while not systemd.is_active("database") do
    log.info("Waiting for database...")
    os.execute("sleep 2")
end

-- Record deployment
state.set("last_webapp_deploy", os.time())
systemd.start("webapp")
```

### With Git Module (GitOps Deployment)

```lua
local git = require("git")
local systemd = require("systemd")

-- Clone latest code
git.clone("https://github.com/company/webapp", "/opt/webapp")

-- Update service and restart
systemd.restart("webapp")
```

## 🚨 **Security Considerations**

1. **Run services as non-root users**
2. **Use specific working directories**
3. **Limit service capabilities**
4. **Set proper file permissions**
5. **Use environment variables for secrets**

## 📚 **Related Documentation**

- [State Module](state.md) - For service coordination
- [Git Module](git.md) - For GitOps deployments
- [Modern DSL](../modern-dsl/syntax.md) - Task definition syntax