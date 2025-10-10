# FRP Module - Fast Reverse Proxy

The FRP module provides integration with [FRP (Fast Reverse Proxy)](https://github.com/fatedier/frp) for exposing local servers behind NAT and firewalls. It offers a fluent API for managing both FRP servers (frps) and clients (frpc) with support for installation, configuration, and lifecycle management.

## Features

- **Installation Management**: Download and install FRP binaries automatically
- **Server Configuration**: Configure FRP server with TOML-based settings
- **Client Configuration**: Set up FRP client with multiple proxy configurations
- **Systemd Integration**: Manage FRP as systemd services
- **TOML to Lua Tables**: Automatic conversion between TOML configs and Lua tables
- **Fluent API**: Chain configuration methods for clean, readable code
- **Agent Delegation**: Execute operations on remote agents via `delegate_to()`

## Installation

The FRP module is automatically registered globally. No `require()` needed.

```lua
-- FRP is available immediately
local server = frp.server("my-server")
local client = frp.client("my-client")
```

## Module Functions

### frp.install(version?, target?)

Install FRP binaries (both frps and frpc) on local or remote machine.

**Parameters:**
- `version` (string, optional): FRP version to install (default: "latest")
- `target` (string, optional): Remote agent name for installation

**Returns:** `(result, error)`

**Example:**

```lua
-- Install latest version locally
local result, err = frp.install()

-- Install specific version on remote agent
local result, err = frp.install("0.52.0", "my-agent")
```

### frp.server(name?)

Create a new FRP server builder instance.

**Parameters:**
- `name` (string, optional): Server name (default: "frps")

**Returns:** `FrpServer` builder instance

**Example:**

```lua
local server = frp.server("production-server")
```

### frp.client(name?)

Create a new FRP client builder instance.

**Parameters:**
- `name` (string, optional): Client name (default: "frpc")

**Returns:** `FrpClient` builder instance

**Example:**

```lua
local client = frp.client("my-client")
```

## FRP Server API

### server:config(config_table)

Set server configuration from Lua table (fluent method).

**Parameters:**
- `config_table` (table): Configuration options

**Common Configuration Options:**
- `bindAddr` - Server bind address (default: "0.0.0.0")
- `bindPort` - Server bind port (default: 7000)
- `vhostHTTPPort` - HTTP vhost port (default: 80)
- `vhostHTTPSPort` - HTTPS vhost port (default: 443)
- `auth.method` - Authentication method ("token")
- `auth.token` - Authentication token
- `webServer` - Web dashboard configuration
- `log` - Logging configuration

**Returns:** `(self, nil)` for method chaining

**Example:**

```lua
local server = frp.server("my-server")
    :config({
        bindAddr = "0.0.0.0",
        bindPort = 7000,
        vhostHTTPPort = 8080,
        auth = {
            method = "token",
            token = "my_secret_token"
        },
        webServer = {
            addr = "0.0.0.0",
            port = 7500,
            user = "admin",
            password = "admin123"
        },
        log = {
            to = "/var/log/frp/frps.log",
            level = "info",
            maxDays = 7
        }
    })
```

### server:config_path(path)

Set the configuration file path (fluent method).

**Parameters:**
- `path` (string): Path to TOML configuration file

**Returns:** `(self, nil)`

**Example:**

```lua
local server = frp.server()
    :config_path("/etc/frp/frps.toml")
```

### server:version(version)

Set FRP version for installation (fluent method).

**Parameters:**
- `version` (string): Version string (e.g., "0.52.0" or "latest")

**Returns:** `(self, nil)`

**Example:**

```lua
local server = frp.server()
    :version("0.52.0")
```

### server:delegate_to(agent)

Execute operations on a remote agent (fluent method).

**Parameters:**
- `agent` (string): Agent name

**Returns:** `(self, nil)`

**Example:**

```lua
local server = frp.server()
    :delegate_to("production-agent")
```

### server:save_config()

Save configuration to TOML file (action method).

**Returns:** `(result, error)`

**Example:**

```lua
local result, err = server:save_config()
if err then
    log.error("Failed to save config: " .. err)
    return false
end
```

### server:load_config()

Load configuration from TOML file (action method).

**Returns:** `(config_table, error)`

**Example:**

```lua
local config, err = server:load_config()
if err then
    log.error("Failed to load config: " .. err)
else
    log.info("Loaded config with " .. #config .. " keys")
end
```

### server:install()

Install FRP binaries and systemd service (action method).

**Returns:** `(result, error)`

**Example:**

```lua
local result, err = server:install()
if err then
    log.error("Installation failed: " .. err)
    return false
end
log.info(result)
```

### server:start()

Start FRP server service (action method).

**Returns:** `(result, error)`

**Example:**

```lua
local result, err = server:start()
```

### server:stop()

Stop FRP server service (action method).

**Returns:** `(result, error)`

**Example:**

```lua
local result, err = server:stop()
```

### server:restart()

Restart FRP server service (action method).

**Returns:** `(result, error)`

**Example:**

```lua
local result, err = server:restart()
```

### server:status()

Get FRP server service status (action method).

**Returns:** `(status_output, error)`

**Example:**

```lua
local status, err = server:status()
log.info(status)
```

### server:enable()

Enable FRP server to start on boot (action method).

**Returns:** `(result, error)`

**Example:**

```lua
local result, err = server:enable()
```

### server:disable()

Disable FRP server from starting on boot (action method).

**Returns:** `(result, error)`

**Example:**

```lua
local result, err = server:disable()
```

## FRP Client API

### client:config(config_table)

Set client configuration from Lua table (fluent method).

**Parameters:**
- `config_table` (table): Configuration options

**Returns:** `(self, nil)`

**Example:**

```lua
local client = frp.client()
    :config({
        auth = {
            method = "token",
            token = "my_token"
        },
        log = {
            to = "/var/log/frp/frpc.log",
            level = "info"
        }
    })
```

### client:server(address, port)

Set FRP server connection details (fluent method).

**Parameters:**
- `address` (string): Server address
- `port` (number): Server port

**Returns:** `(self, nil)`

**Example:**

```lua
local client = frp.client()
    :server("frp.example.com", 7000)
```

### client:proxy(proxy_config)

Add a proxy configuration (fluent method). Can be called multiple times.

**Parameters:**
- `proxy_config` (table): Proxy configuration

**Common Proxy Options:**
- `name` - Unique proxy name
- `type` - Proxy type ("tcp", "http", "https", "udp")
- `localIP` - Local service IP (default: "127.0.0.1")
- `localPort` - Local service port
- `remotePort` - Remote port (for TCP/UDP)
- `customDomains` - Custom domains (for HTTP/HTTPS)

**Returns:** `(self, nil)`

**Example:**

```lua
local client = frp.client()
    :proxy({
        name = "web",
        type = "http",
        localIP = "127.0.0.1",
        localPort = 3000,
        customDomains = {"myapp.example.com"}
    })
    :proxy({
        name = "ssh",
        type = "tcp",
        localPort = 22,
        remotePort = 6000
    })
```

### client:config_path(path)

Set configuration file path (fluent method).

**Parameters:**
- `path` (string): Path to configuration file

**Returns:** `(self, nil)`

### client:version(version)

Set FRP version (fluent method).

**Parameters:**
- `version` (string): Version string

**Returns:** `(self, nil)`

### client:delegate_to(agent)

Execute on remote agent (fluent method).

**Parameters:**
- `agent` (string): Agent name

**Returns:** `(self, nil)`

### client:save_config()

Save configuration to TOML file (action method).

**Returns:** `(result, error)`

### client:load_config()

Load configuration from TOML file (action method).

**Returns:** `(config_table, error)`

### client:install()

Install FRP binaries and systemd service (action method).

**Returns:** `(result, error)`

### client:start()

Start FRP client service (action method).

**Returns:** `(result, error)`

### client:stop()

Stop FRP client service (action method).

**Returns:** `(result, error)`

### client:restart()

Restart FRP client service (action method).

**Returns:** `(result, error)`

### client:status()

Get service status (action method).

**Returns:** `(status_output, error)`

### client:enable()

Enable service on boot (action method).

**Returns:** `(result, error)`

### client:disable()

Disable service on boot (action method).

**Returns:** `(result, error)`

## Complete Examples

### Example 1: FRP Server Setup

```lua
local install_task = task("install_frp")
    :command(function()
        local result, err = frp.install("latest")
        if err then
            return false, err
        end
        return true, "FRP installed"
    end)
    :build()

local config_task = task("configure_server")
    :depends_on({"install_frp"})
    :command(function()
        local server = frp.server("production")
            :config({
                bindAddr = "0.0.0.0",
                bindPort = 7000,
                auth = {
                    method = "token",
                    token = "secure_token_123"
                }
            })

        local result, err = server:save_config()
        if err then
            return false, err
        end
        return true, "Configuration saved"
    end)
    :build()

local start_task = task("start_server")
    :depends_on({"configure_server"})
    :command(function()
        local server = frp.server()

        local _, err = server:enable()
        if err then return false, err end

        local _, err = server:start()
        if err then return false, err end

        return true, "Server started"
    end)
    :build()

workflow
    .define("frp_server_setup")
    :tasks({install_task, config_task, start_task})
```

### Example 2: FRP Client with Multiple Services

```lua
local setup_client = task("setup_frp_client")
    :command(function()
        local client = frp.client("services")
            :server("frp.example.com", 7000)
            :config({
                auth = {
                    method = "token",
                    token = "secure_token_123"
                }
            })
            -- Web application
            :proxy({
                name = "web",
                type = "http",
                localPort = 3000,
                customDomains = {"myapp.example.com"}
            })
            -- SSH access
            :proxy({
                name = "ssh",
                type = "tcp",
                localPort = 22,
                remotePort = 6000
            })
            -- PostgreSQL database
            :proxy({
                name = "postgres",
                type = "tcp",
                localPort = 5432,
                remotePort = 6432
            })

        local result, err = client:save_config()
        if err then
            return false, err
        end

        local _, err = client:install()
        if err then return false, err end

        local _, err = client:enable()
        if err then return false, err end

        local _, err = client:start()
        if err then return false, err end

        return true, "Client configured and started"
    end)
    :build()
```

### Example 3: Remote Agent Deployment

```lua
local deploy_server = task("deploy_frp_server")
    :command(function()
        local server = frp.server("remote-server")
            :delegate_to("production-agent")
            :version("0.52.0")
            :config({
                bindAddr = "0.0.0.0",
                bindPort = 7000,
                vhostHTTPPort = 80,
                vhostHTTPSPort = 443
            })

        local _, err = server:install()
        if err then return false, "Install failed: " .. err end

        local _, err = server:save_config()
        if err then return false, "Config failed: " .. err end

        local _, err = server:enable()
        if err then return false, "Enable failed: " .. err end

        local _, err = server:start()
        if err then return false, "Start failed: " .. err end

        return true, "Server deployed on remote agent"
    end)
    :build()
```

### Example 4: Configuration Management

```lua
local backup_config = task("backup_frp_config")
    :command(function()
        local server = frp.server()

        -- Load current configuration
        local config, err = server:load_config()
        if err then
            return false, "Failed to load config: " .. err
        end

        -- Backup to file
        local backup_file = "/tmp/frps_backup_" .. os.date("%Y%m%d_%H%M%S") .. ".lua"
        fs.write(backup_file, "return " .. inspect(config))

        log.info("Configuration backed up to: " .. backup_file)
        return true, "Backup created"
    end)
    :build()

local update_config = task("update_frp_config")
    :depends_on({"backup_frp_config"})
    :command(function()
        local server = frp.server()
            :config({
                bindPort = 7001,  -- Change port
                auth = {
                    method = "token",
                    token = "new_token_456"
                }
            })

        local _, err = server:save_config()
        if err then
            return false, err
        end

        local _, err = server:restart()
        if err then
            return false, err
        end

        return true, "Configuration updated and service restarted"
    end)
    :build()
```

## TOML Configuration

The FRP module automatically handles TOML configuration files. Configuration provided via Lua tables is converted to TOML format when saved.

### Server Configuration Example (frps.toml)

```toml
bindAddr = "0.0.0.0"
bindPort = 7000
vhostHTTPPort = 80
vhostHTTPSPort = 443

[auth]
method = "token"
token = "secure_token_123"

[webServer]
addr = "0.0.0.0"
port = 7500
user = "admin"
password = "admin123"

[log]
to = "/var/log/frp/frps.log"
level = "info"
maxDays = 7
```

### Client Configuration Example (frpc.toml)

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[auth]
method = "token"
token = "secure_token_123"

[[proxies]]
name = "web"
type = "http"
localIP = "127.0.0.1"
localPort = 3000
customDomains = ["myapp.example.com"]

[[proxies]]
name = "ssh"
type = "tcp"
localPort = 22
remotePort = 6000
```

## Best Practices

1. **Use Authentication**: Always configure token-based authentication for production
2. **Version Pinning**: Specify exact FRP version instead of "latest" for reproducibility
3. **Remote Deployment**: Use `delegate_to()` for managing FRP on remote machines
4. **Configuration Backup**: Backup configurations before updates
5. **Enable on Boot**: Use `:enable()` to ensure FRP starts after system reboots
6. **Monitoring**: Check service status regularly with `:status()`
7. **Logging**: Configure proper log files for troubleshooting
8. **Service Management**: Use systemd integration for reliable service management

## Troubleshooting

### Installation Issues

```lua
-- Check if binaries are installed
local result = exec.run("which frps frpc")
log.info(result)

-- Verify version
local version = exec.run("frps --version")
log.info(version)
```

### Configuration Issues

```lua
-- Validate TOML syntax
local config, err = server:load_config()
if err then
    log.error("Config validation failed: " .. err)
end

-- Check file permissions
local perms = exec.run("ls -l /etc/frp/frps.toml")
log.info(perms)
```

### Service Issues

```lua
-- Check service status
local status, _ = server:status()
log.info(status)

-- View logs
local logs = exec.run("journalctl -u frps -n 50")
log.info(logs)

-- Check if port is in use
local port_check = exec.run("ss -tlnp | grep 7000")
log.info(port_check)
```

## See Also

- [Incus Module](./incus.md) - Container/VM management
- [Systemd Module](./systemd.md) - Service management
- [Exec Module](./exec.md) - Command execution
- [FRP Official Documentation](https://github.com/fatedier/frp)
