# Facts Module

The `facts` module provides access to system information collected from agents by the Sloth Runner master. This allows you to query detailed information about remote systems, including hardware specs, installed packages, running services, and more.

## Overview

The facts module communicates with the Sloth Runner master to retrieve cached system information that agents periodically collect and report. This enables infrastructure discovery, validation, and conditional task execution based on real-time system state.

## Installation

The `facts` module is automatically available as a global module in all Sloth Runner tasks. No `require` statement is needed.

## Core Functions

### facts.get_all()

Retrieves all facts from an agent.

**Syntax:**
```lua
local info, err = facts.get_all({ agent = "agent-name" })
```

**Parameters:**
- `agent` (string): The name of the agent to query

**Returns:**
- Table containing all system information, or nil on error
- Error message if operation failed

**Example:**
```lua
local check_system = task("check_system")
    :description("Check system information")
    :command(function(this, params)
        local info, err = facts.get_all({ agent = "prod-server-01" })
        if err then
            return false, "Failed to get facts: " .. err
        end

        print("Hostname: " .. info.hostname)
        print("Platform: " .. info.platform)
        print("CPUs: " .. info.cpus)
        print("Memory Total: " .. info.memory.total)

        return true, "System information retrieved"
    end)
    :build()
```

### facts.get_hostname()

Gets the hostname of an agent.

**Syntax:**
```lua
local hostname, err = facts.get_hostname({ agent = "agent-name" })
```

**Example:**
```lua
local check_hostname = task("check_hostname")
    :description("Check agent hostname")
    :command(function(this, params)
        local hostname, err = facts.get_hostname({ agent = "web-01" })
        if err then
            return false, "Failed: " .. err
        end
        print("Hostname: " .. hostname)
        return true, "Hostname checked: " .. hostname
    end)
    :build()
```

### facts.get_platform()

Gets platform information (OS, kernel, architecture).

**Syntax:**
```lua
local platform, err = facts.get_platform({ agent = "agent-name" })
```

**Returns:**
Table with fields:
- `os`: Operating system name
- `family`: Platform family (e.g., "debian", "redhat")
- `version`: OS version
- `architecture`: System architecture (e.g., "amd64", "arm64")
- `kernel`: Kernel name
- `kernel_version`: Kernel version
- `virtualization`: Virtualization technology (if any)

**Example:**
```lua
local check_os = task("check_os")
    :description("Check OS and platform information")
    :command(function(this, params)
        local platform, err = facts.get_platform({ agent = "db-server" })
        if err then
            return false, "Failed: " .. err
        end

        print(string.format("OS: %s %s", platform.os, platform.version))
        print(string.format("Arch: %s", platform.architecture))
        print(string.format("Kernel: %s %s", platform.kernel, platform.kernel_version))

        if platform.virtualization ~= "" then
            print("Running on: " .. platform.virtualization)
        end

        return true, "Platform information retrieved"
    end)
    :build()
```

### facts.get_memory()

Gets memory information.

**Syntax:**
```lua
local memory, err = facts.get_memory({ agent = "agent-name" })
```

**Returns:**
Table with fields:
- `total`: Total memory in bytes
- `available`: Available memory in bytes
- `used`: Used memory in bytes
- `used_percent`: Memory usage percentage
- `free`: Free memory in bytes
- `cached`: Cached memory in bytes
- `buffers`: Buffer memory in bytes

**Example:**
```lua
local check_memory = task("check_memory")
    :description("Check memory usage")
    :command(function(this, params)
        local mem, err = facts.get_memory({ agent = "app-server" })
        if err then
            return false, "Failed: " .. err
        end

        local total_gb = mem.total / 1024 / 1024 / 1024
        local used_gb = mem.used / 1024 / 1024 / 1024

        print(string.format("Memory: %.2f GB / %.2f GB (%.1f%%)",
            used_gb, total_gb, mem.used_percent))

        if mem.used_percent > 90 then
            print("WARNING: Memory usage is critical!")
            return false, "Critical memory usage detected"
        end

        return true, "Memory check completed"
    end)
    :build()
```

### facts.get_disk()

Gets disk/filesystem information.

**Syntax:**
```lua
-- Get all disks
local disks, err = facts.get_disk({ agent = "agent-name" })

-- Get specific mountpoint
local disk, err = facts.get_disk({ 
    agent = "agent-name", 
    mountpoint = "/home" 
})
```

**Returns:**
- Array of disk information (if no mountpoint specified)
- Single disk table (if mountpoint specified)

Each disk table contains:
- `device`: Device name
- `mountpoint`: Mount path
- `fstype`: Filesystem type
- `total`: Total space in bytes
- `used`: Used space in bytes
- `free`: Free space in bytes
- `used_percent`: Usage percentage

**Example:**
```lua
local check_disk_space = task("check_disk_space")
    :description("Check disk space usage")
    :command(function(this, params)
        local disks, err = facts.get_disk({ agent = "file-server" })
        if err then
            return false, "Failed: " .. err
        end

        local has_warning = false
        for i, disk in ipairs(disks) do
            local used_gb = disk.used / 1024 / 1024 / 1024
            local total_gb = disk.total / 1024 / 1024 / 1024

            print(string.format("%s: %.2f GB / %.2f GB (%.1f%%)",
                disk.mountpoint, used_gb, total_gb, disk.used_percent))

            if disk.used_percent > 85 then
                print("  WARNING: Low disk space!")
                has_warning = true
            end
        end

        if has_warning then
            return false, "Low disk space detected on some volumes"
        end

        return true, "Disk space check completed"
    end)
    :build()
```

### facts.get_network()

Gets network interface information.

**Syntax:**
```lua
-- Get all interfaces
local interfaces, err = facts.get_network({ agent = "agent-name" })

-- Get specific interface
local iface, err = facts.get_network({ 
    agent = "agent-name", 
    interface = "eth0" 
})
```

**Returns:**
Table or array of tables with fields:
- `name`: Interface name
- `mac`: MAC address
- `mtu`: MTU size
- `is_up`: Interface status (boolean)
- `speed`: Link speed
- `addresses`: Array of IP addresses

**Example:**
```lua
local check_network = task("check_network")
    :description("Check network interfaces")
    :command(function(this, params)
        local ifaces, err = facts.get_network({ agent = "router" })
        if err then
            return false, "Failed: " .. err
        end

        for _, iface in ipairs(ifaces) do
            print(string.format("Interface: %s", iface.name))
            print(string.format("  MAC: %s", iface.mac))
            print(string.format("  Status: %s", iface.is_up and "UP" or "DOWN"))
            print("  IPs:")
            for _, addr in ipairs(iface.addresses) do
                print("    - " .. addr)
            end
        end

        return true, "Network interfaces checked"
    end)
    :build()
```

### facts.get_packages()

Gets information about installed packages.

**Syntax:**
```lua
local pkg_info, err = facts.get_packages({ agent = "agent-name" })
```

**Returns:**
Table with fields:
- `manager`: Package manager name
- `installed_count`: Number of installed packages
- `updates_available`: Number of available updates
- `packages`: Array of installed packages
- `updates`: Array of available updates

Each package contains:
- `name`: Package name
- `version`: Installed version
- `architecture`: Package architecture
- `description`: Package description

**Example:**
```lua
local check_packages = task("check_packages")
    :description("Check installed packages and updates")
    :command(function(this, params)
        local pkgs, err = facts.get_packages({ agent = "server-01" })
        if err then
            return false, "Failed: " .. err
        end

        print(string.format("Package Manager: %s", pkgs.manager))
        print(string.format("Installed Packages: %d", pkgs.installed_count))
        print(string.format("Updates Available: %d", pkgs.updates_available))

        if pkgs.updates_available > 0 then
            print("\nAvailable Updates:")
            for _, upd in ipairs(pkgs.updates) do
                print(string.format("  - %s: %s", upd.name, upd.version))
            end
        end

        return true, string.format("Package check completed: %d updates available", pkgs.updates_available)
    end)
    :build()
```

### facts.get_package()

Checks if a specific package is installed.

**Syntax:**
```lua
local pkg, err = facts.get_package({ 
    agent = "agent-name", 
    name = "package-name" 
})
```

**Returns:**
Table with fields:
- `name`: Package name
- `installed`: Boolean indicating if package is installed
- `version`: Installed version (if installed)
- `architecture`: Package architecture (if installed)
- `description`: Package description (if installed)

**Example:**
```lua
local ensure_nginx = task("ensure_nginx")
    :description("Ensure nginx is installed")
    :command(function(this, params)
        local pkg, err = facts.get_package({
            agent = "web-server",
            name = "nginx"
        })

        if err then
            return false, "Failed: " .. err
        end

        if pkg.installed then
            print(string.format("nginx %s is installed", pkg.version))
            return true, "nginx is already installed"
        else
            print("nginx is not installed - installing...")
            pkg.install({ packages = {"nginx"} }):delegate_to("web-server")
            return true, "nginx installation initiated"
        end
    end)
    :build()
```

### facts.get_services()

Gets information about all services.

**Syntax:**
```lua
local services, err = facts.get_services({ agent = "agent-name" })
```

**Returns:**
Array of service tables with fields:
- `name`: Service name
- `status`: Service status
- `state`: Service state

**Example:**
```lua
local list_services = task("list_services")
    :description("List active services")
    :command(function(this, params)
        local services, err = facts.get_services({ agent = "app-server" })
        if err then
            return false, "Failed: " .. err
        end

        print("Active Services:")
        local count = 0
        for _, svc in ipairs(services) do
            if svc.status == "active" then
                print(string.format("  - %s: %s", svc.name, svc.state))
                count = count + 1
            end
        end

        return true, string.format("Found %d active services", count)
    end)
    :build()
```

### facts.get_service()

Gets status of a specific service.

**Syntax:**
```lua
local service, err = facts.get_service({ 
    agent = "agent-name", 
    name = "service-name" 
})
```

**Example:**
```lua
local check_nginx_status = task("check_nginx_status")
    :description("Check nginx service status")
    :command(function(this, params)
        local svc, err = facts.get_service({
            agent = "web-01",
            name = "nginx"
        })

        if err then
            return false, "Failed: " .. err
        end

        print(string.format("nginx: %s (%s)", svc.status, svc.state))

        if svc.status ~= "active" then
            print("WARNING: nginx is not active!")
            return false, "nginx is not active"
        end

        return true, "nginx is active and running"
    end)
    :build()
```

### facts.get_users()

Gets information about system users.

**Syntax:**
```lua
local users, err = facts.get_users({ agent = "agent-name" })
```

**Returns:**
Array of user tables with fields:
- `username`: User name
- `uid`: User ID
- `gid`: Group ID
- `home`: Home directory
- `shell`: Login shell

**Example:**
```lua
local list_users = task("list_users")
    :description("List system users")
    :command(function(this, params)
        local users, err = facts.get_users({ agent = "server" })
        if err then
            return false, "Failed: " .. err
        end

        print("System Users:")
        for _, user in ipairs(users) do
            print(string.format("  %s (UID: %s) - %s",
                user.username, user.uid, user.shell))
        end

        return true, string.format("Listed %d users", #users)
    end)
    :build()
```

### facts.get_user()

Gets information about a specific user.

**Syntax:**
```lua
local user, err = facts.get_user({ 
    agent = "agent-name", 
    username = "username" 
})
```

**Example:**
```lua
local check_user = task("check_user")
    :description("Check if deploy user exists")
    :command(function(this, params)
        local user, err = facts.get_user({
            agent = "server",
            username = "deploy"
        })

        if err then
            print("User 'deploy' not found")
            -- Create user
            user.create({
                name = "deploy",
                home = "/home/deploy",
                shell = "/bin/bash"
            }):delegate_to("server")
            return true, "User 'deploy' created"
        else
            print(string.format("User 'deploy' exists: %s", user.home))
            return true, "User 'deploy' already exists"
        end
    end)
    :build()
```

### facts.get_processes()

Gets process statistics.

**Syntax:**
```lua
local procs, err = facts.get_processes({ agent = "agent-name" })
```

**Returns:**
Table with fields:
- `total`: Total number of processes
- `running`: Running processes
- `sleeping`: Sleeping processes
- `zombie`: Zombie processes

**Example:**
```lua
local check_processes = task("check_processes")
    :description("Check process statistics")
    :command(function(this, params)
        local procs, err = facts.get_processes({ agent = "server" })
        if err then
            return false, "Failed: " .. err
        end

        print(string.format("Processes: %d total, %d running, %d sleeping",
            procs.total, procs.running, procs.sleeping))

        if procs.zombie > 0 then
            print(string.format("WARNING: %d zombie processes!", procs.zombie))
            return false, string.format("Found %d zombie processes", procs.zombie)
        end

        return true, "Process check completed"
    end)
    :build()
```

### facts.get_mounts()

Gets filesystem mount information.

**Syntax:**
```lua
local mounts, err = facts.get_mounts({ agent = "agent-name" })
```

**Returns:**
Array of mount tables with fields:
- `device`: Device name
- `mountpoint`: Mount path
- `fstype`: Filesystem type
- `options`: Mount options

**Example:**
```lua
local list_mounts = task("list_mounts")
    :description("List mounted filesystems")
    :command(function(this, params)
        local mounts, err = facts.get_mounts({ agent = "server" })
        if err then
            return false, "Failed: " .. err
        end

        print("Mounted Filesystems:")
        for _, mount in ipairs(mounts) do
            print(string.format("  %s on %s type %s (%s)",
                mount.device, mount.mountpoint, mount.fstype, mount.options))
        end

        return true, string.format("Listed %d mount points", #mounts)
    end)
    :build()
```

### facts.get_uptime()

Gets system uptime information.

**Syntax:**
```lua
local uptime, err = facts.get_uptime({ agent = "agent-name" })
```

**Returns:**
Table with fields:
- `seconds`: Uptime in seconds
- `boot_time`: Boot time (Unix timestamp)
- `timezone`: System timezone

**Example:**
```lua
local check_uptime = task("check_uptime")
    :description("Check system uptime")
    :command(function(this, params)
        local uptime, err = facts.get_uptime({ agent = "server" })
        if err then
            return false, "Failed: " .. err
        end

        local days = math.floor(uptime.seconds / 86400)
        local hours = math.floor((uptime.seconds % 86400) / 3600)
        local mins = math.floor((uptime.seconds % 3600) / 60)

        print(string.format("Uptime: %d days, %d hours, %d minutes",
            days, hours, mins))
        print(string.format("Timezone: %s", uptime.timezone))

        return true, string.format("System uptime: %d days", days)
    end)
    :build()
```

### facts.get_load()

Gets system load average.

**Syntax:**
```lua
local load, err = facts.get_load({ agent = "agent-name" })
```

**Returns:**
Array with load averages [1min, 5min, 15min]

**Example:**
```lua
local check_load = task("check_load")
    :description("Check system load average")
    :command(function(this, params)
        local load, err = facts.get_load({ agent = "server" })
        if err then
            return false, "Failed: " .. err
        end

        print(string.format("Load Average: %.2f, %.2f, %.2f",
            load[1], load[2], load[3]))

        if load[1] > 4.0 then
            print("WARNING: High load!")
            return false, "System load is too high"
        end

        return true, "Load average checked"
    end)
    :build()
```

### facts.get_kernel()

Gets kernel information.

**Syntax:**
```lua
local kernel, err = facts.get_kernel({ agent = "agent-name" })
```

**Returns:**
Table with fields:
- `name`: Kernel name
- `version`: Kernel version

**Example:**
```lua
local check_kernel = task("check_kernel")
    :description("Check kernel information")
    :command(function(this, params)
        local kernel, err = facts.get_kernel({ agent = "server" })
        if err then
            return false, "Failed: " .. err
        end

        print(string.format("Kernel: %s %s", kernel.name, kernel.version))

        return true, string.format("Kernel: %s %s", kernel.name, kernel.version)
    end)
    :build()
```

### facts.query()

Performs a query on facts (experimental).

**Syntax:**
```lua
local result, err = facts.query({ 
    agent = "agent-name", 
    path = "$.memory.total" 
})
```

**Note:** This function is experimental and may change in future versions.

## 🔥 Exemplo Destacado: Validação e Deploy Inteligente

Este exemplo demonstra como usar facts para tomar decisões inteligentes durante o deploy, validando o sistema alvo e adaptando o comportamento baseado nas condições reais.

```lua
local intelligent_deploy = task("intelligent-deploy")
    :description("Intelligent deployment with system validation")
    :command(function(this, params)
        local target_agent = values.target or "prod-server-01"

        -- 📊 Coletar informações do sistema alvo
        local info, err = facts.get_all({ agent = target_agent })
        if err then
            return false, "Cannot reach agent: " .. err
        end

        log.info("🔍 Analyzing " .. info.hostname)
        log.info("   Platform: " .. info.platform.os .. " " .. info.platform.version)
        log.info("   Memory: " .. string.format("%.2f GB", info.memory.total / 1024 / 1024 / 1024))
        log.info("   Arch: " .. info.platform.architecture)

        -- ✅ Validação de requisitos mínimos
        local mem_gb = info.memory.total / 1024 / 1024 / 1024
        if mem_gb < 4 then
            return false, "Insufficient memory: need 4GB, have " .. string.format("%.2f GB", mem_gb)
        end

        -- ✅ Verificar espaço em disco
        local root_disk, _ = facts.get_disk({
            agent = target_agent,
            mountpoint = "/"
        })
        if root_disk.used_percent > 85 then
            log.warn("⚠️  Disk usage high: " .. root_disk.used_percent .. "%")
        end

        -- ✅ Verificar se Docker já está instalado
        local docker_pkg, _ = facts.get_package({
            agent = target_agent,
            name = "docker"
        })

        if not docker_pkg.installed then
            log.info("📦 Installing Docker...")
            pkg.install({ packages = {"docker.io"} }):delegate_to(target_agent)
        else
            log.info("✅ Docker already installed: " .. docker_pkg.version)
        end

        -- 🚀 Deploy baseado na arquitetura
        local image_tag = "latest"
        if info.platform.architecture == "arm64" then
            image_tag = "latest-arm64"
        end

        log.info("🚀 Deploying with image: myapp:" .. image_tag)

        -- Continue with deployment...
        log.info("✅ Deploy completed successfully!")

        return true, "Deploy completed successfully"
    end)
    :build()
```

**Recursos demonstrados:**

- 📊 Coleta completa de informações do sistema remoto
- ✅ Validação de requisitos (memória, disco, pacotes)
- 🧠 Deploy condicional baseado em arquitetura (x86/ARM)
- 🔄 Instalação automática de dependências
- 🎯 Uso de `values` para parametrização

---

## Complete Examples

### Example 1: Pre-deployment System Validation

```lua
local validate_system = task("validate_system")
    :description("Validate system requirements before deployment")
    :command(function(this, params)
        local hostname, _ = facts.get_hostname({ agent = "prod-app-01" })
        print("Validating: " .. hostname)

        -- Check OS version
        local platform, err = facts.get_platform({ agent = "prod-app-01" })
        if err then
            return false, "Cannot get platform: " .. err
        end

        if platform.os ~= "linux" then
            return false, "Expected Linux, got: " .. platform.os
        end

        -- Check memory
        local mem, err = facts.get_memory({ agent = "prod-app-01" })
        if err then
            return false, "Cannot get memory: " .. err
        end

        local mem_gb = mem.total / 1024 / 1024 / 1024
        if mem_gb < 8 then
            return false, string.format("Insufficient memory: %.2f GB (need 8 GB)", mem_gb)
        end

        -- Check disk space
        local disk, err = facts.get_disk({
            agent = "prod-app-01",
            mountpoint = "/"
        })
        if err then
            return false, "Cannot get disk info: " .. err
        end

        if disk.used_percent > 80 then
            return false, string.format("Disk usage too high: %.1f%%", disk.used_percent)
        end

        -- Check required package
        local pkg, err = facts.get_package({
            agent = "prod-app-01",
            name = "docker"
        })
        if err then
            return false, "Cannot check package: " .. err
        end

        if not pkg.installed then
            return false, "Docker is not installed"
        end

        print("✓ All validations passed!")
        return true, "All system validations passed"
    end)
    :build()
```

### Example 2: Dynamic Inventory Based on Facts

```lua
local discover_web_servers = task("discover_web_servers")
    :description("Discover web servers running nginx")
    :command(function(this, params)
        local agents = {"server-01", "server-02", "server-03"}
        local web_servers = {}

        for _, agent in ipairs(agents) do
            -- Check if nginx is running
            local svc, err = facts.get_service({
                agent = agent,
                name = "nginx"
            })

            if not err and svc.status == "active" then
                -- Get IP address
                local iface, _ = facts.get_network({
                    agent = agent,
                    interface = "eth0"
                })

                if iface and #iface.addresses > 0 then
                    table.insert(web_servers, {
                        name = agent,
                        ip = iface.addresses[1]
                    })
                end
            end
        end

        print("Discovered Web Servers:")
        for _, server in ipairs(web_servers) do
            print(string.format("  - %s: %s", server.name, server.ip))
        end

        return true, string.format("Discovered %d web servers", #web_servers)
    end)
    :build()
```

### Example 3: Conditional Deployment Based on System State

```lua
local deploy_app = task("deploy_app")
    :description("Deploy application with conditional strategy based on system state")
    :command(function(this, params)
        local agent = "app-server-01"

        -- Get current system state
        local platform, _ = facts.get_platform({ agent = agent })
        local mem, _ = facts.get_memory({ agent = agent })

        -- Decide deployment strategy based on available resources
        if mem.available < 2 * 1024 * 1024 * 1024 then  -- Less than 2GB
            print("Low memory - using minimal deployment")
            -- Deploy with minimal resources
        else
            print("Sufficient memory - using full deployment")
            -- Deploy with full resources
        end

        -- Check if old version is installed
        local old_app, _ = facts.get_package({
            agent = agent,
            name = "myapp"
        })

        if old_app.installed then
            print("Stopping old version: " .. old_app.version)
            systemd.stop({ unit = "myapp" }):delegate_to(agent)
        end

        -- Continue with deployment...
        print("Deploying new version...")

        return true, "Application deployed successfully"
    end)
    :build()
```

## Best Practices

1. **Error Handling**: Always check for errors when calling facts functions
2. **Agent Availability**: Ensure the agent is online before querying facts
3. **Caching**: Facts are cached by the master; they may not reflect real-time state
4. **Performance**: Avoid excessive fact queries in loops; cache results when possible
5. **Validation**: Use facts for pre-deployment validation to catch issues early

## See Also

- [Agent Module](./agent.md) - Managing agents
- [InfraTest Module](./infra_test.md) - Infrastructure testing
- [Package Module](./pkg.md) - Package management
- [Systemd Module](./systemd.md) - Service management
