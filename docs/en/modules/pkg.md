# 📦 Package Manager Module

The `pkg` module provides comprehensive cross-platform package management functionality. It automatically detects the system's package manager and provides a unified interface for managing packages.

## 🎯 Supported Package Managers

- **apt / apt-get** (Debian/Ubuntu)
- **yum / dnf** (RHEL/CentOS/Fedora)
- **pacman** (Arch Linux)
- **zypper** (openSUSE)
- **brew** (macOS - Homebrew)

## 📚 Functions Overview

| Function | Description |
|----------|-------------|
| `pkg.install({packages = ...})` | Install one or more packages |
| `pkg.remove({packages = ...})` | Remove one or more packages |
| `pkg.update({})` | Update package cache/list |
| `pkg.upgrade({})` | Upgrade all packages |
| `pkg.search({query = ...})` | Search for packages |
| `pkg.info({package = ...})` | Get package information |
| `pkg.list({})` | List installed packages |
| `pkg.is_installed({package = ...})` | Check if package is installed |
| `pkg.get_manager({})` | Get detected package manager |
| `pkg.clean({})` | Clean package cache |
| `pkg.autoremove({})` | Remove unused dependencies |
| `pkg.which({executable = ...})` | Find executable path |
| `pkg.version({package = ...})` | Get package version |
| `pkg.deps({package = ...})` | List package dependencies |
| `pkg.install_local({file = ...})` | Install from local file |

## 📖 Detailed Documentation

### Installation & Removal

#### `pkg.install({packages = ...})`

Installs one or more packages.

**Parameters:**
- `packages`: String (single package) or Table (multiple packages)

**Returns:**
- `success` (boolean): `true` on success, `false` on failure
- `output` (string): Command output

**Examples:**

=== "Modern DSL"
    ```lua
    
    local install_tools = task("install_tools")
        :description("Install development tools")
        :command(function(this, params)
            log.info("Installing tools...")
            
            -- Install multiple packages
            local tools = {"git", "curl", "wget", "vim"}
            local success, output = pkg.install({packages = tools})
            
            if success then
                log.info("✅ Tools installed successfully!")
                return true, "Installed"
            else
                log.error("❌ Failed: " .. output)
                return false, "Failed"
            end
        end)
        :timeout("300s")
        :build()
    
    workflow.define("setup")
        :tasks({ install_tools })
    ```

=== "With delegate_to"
    ```lua
    
    local install_on_agent = task("install_on_agent")
        :description("Install packages on remote agent")
        :command(function(this, params)
            log.info("Installing on remote agent...")
            
            local success, output = pkg.install({packages = {"htop", "ncdu"}})
            
            if success then
                log.info("✅ Installed on agent!")
                return true, "OK"
            else
                return false, "Failed"
            end
        end)
        :delegate_to("production-server")
        :timeout("300s")
        :build()
    
    workflow.define("remote_install")
        :tasks({ install_on_agent })
    ```

=== "Single Package"
    ```lua
    
    local install_nginx = task("install_nginx")
        :description("Install nginx web server")
        :command(function(this, params)
            -- Install single package
            local success, output = pkg.install({packages = "nginx"})
            
            if success then
                log.info("✅ nginx installed!")
                return true, "OK"
            else
                return false, "Failed: " .. output
            end
        end)
        :timeout("300s")
        :build()
    ```

#### `pkg.remove({packages = ...})`

Removes one or more packages.

**Parameters:**
- `packages`: String or Table

**Returns:**
- `success` (boolean), `output` (string)

**Example:**

```lua

local cleanup = task("cleanup")
    :description("Remove unnecessary packages")
    :command(function(this, params)
        local packages = {"package1", "package2"}
        local success, output = pkg.remove({packages = packages})
        
        if success then
            log.info("✅ Packages removed")
            return true, "Removed"
        end
        return false, "Failed"
    end)
    :timeout("180s")
    :build()
```

### Package Information

#### `pkg.search({query = ...})`

Searches for packages.

**Example:**

```lua

local search_python = task("search_python")
    :description("Search for Python packages")
    :command(function(this, params)
        local success, results = pkg.search({query = "python3"})
        
        if success then
            log.info("Search results:")
            local count = 0
            for line in results:gmatch("[^\r\n]+") do
                if count < 10 then
                    log.info("  • " .. line)
                end
                count = count + 1
            end
            return true, count .. " results"
        end
        return false, "Search failed"
    end)
    :timeout("60s")
    :build()
```

#### `pkg.info({package = ...})`

Gets package information.

**Example:**

```lua
local success, info = pkg.info({package = "curl"})
if success then
    log.info("Package info:\n" .. info)
end
```

#### `pkg.list({})`

Lists installed packages.

**Returns:** `success` (boolean), `packages` (table)

**Example:**

```lua
local success, packages = pkg.list({})
if success and type(packages) == "table" then
    local count = 0
    for _ in pairs(packages) do count = count + 1 end
    log.info("📦 Total: " .. count .. " packages")
end
```

### System Maintenance

#### `pkg.update({})`

Updates package cache.

**Example:**

```lua
local update_cache = task("update_cache")
    :description("Update package cache")
    :command(function(this, params)
        log.info("Updating...")
        return pkg.update({})
    end)
    :timeout("120s")
    :build()
```

#### `pkg.upgrade({})`

Upgrades all packages.

#### `pkg.clean({})`

Cleans package cache.

#### `pkg.autoremove({})`

Removes unused dependencies.

**Example:**

```lua
local maintenance = task("maintenance")
    :description("System maintenance")
    :command(function(this, params)
        -- Update
        pkg.update({})
        
        -- Upgrade
        pkg.upgrade({})
        
        -- Clean
        pkg.clean({})
        pkg.autoremove({})
        
        return true, "Maintenance complete"
    end)
    :timeout("600s")
    :build()
```

### Advanced Functions

#### `pkg.is_installed({package = ...})`

Checks if installed.

**Example:**

```lua

local check_requirements = task("check_requirements")
    :description("Check required packages")
    :command(function(this, params)
        local required = {"git", "curl", "wget"}
        local missing = {}
        
        for _, pkg_name in ipairs(required) do
            local installed, _ = pkg.is_installed({package = pkg_name})
            if not installed then
                table.insert(missing, pkg_name)
            end
        end
        
        if #missing > 0 then
            return false, "Missing: " .. table.concat(missing, ", ")
        end
        
        return true, "All OK"
    end)
    :build()
```

#### `pkg.get_manager({})`

Returns package manager name.

**Example:**

```lua
local manager, err = pkg.get_manager({})
log.info("Manager: " .. (manager or "unknown"))
```

#### `pkg.which({executable = ...})`

Finds executable path.

**Example:**

```lua
local path, err = pkg.which({executable = "git"})
if path then
    log.info("Git at: " .. path)
end
```

#### `pkg.version({package = ...})`

Gets package version.

#### `pkg.deps({package = ...})`

Lists dependencies.

#### `pkg.install_local({file = ...})`

Installs from local file (.deb, .rpm).

**Example:**

```lua
local success, output = pkg.install_local({file = "/path/to/package.deb"})
if success then
    log.info("✅ Package installed from local file")
end
```

## 🎯 Complete Examples

### Development Environment Setup

```lua

local update = task("update")
    :command(function() return pkg.update({}) end)
    :build()

local install_tools = task("install_tools")
    :command(function()
        local tools = {"git", "curl", "wget", "vim", "htop"}
        return pkg.install({packages = tools})
    end)
    :depends_on({"update"})
    :build()

local verify = task("verify")
    :command(function()
        for _, tool in ipairs({"git", "curl"}) do
            local installed, _ = pkg.is_installed({package = tool})
            if installed then
                local path, _ = pkg.which({executable = tool})
                log.info("✅ " .. tool .. " (" .. (path or "unknown") .. ")")
            end
        end
        return true, "OK"
    end)
    :depends_on({"install_tools"})
    :build()

workflow.define("setup_dev")
    :tasks({ update, install_tools, verify })
```

### Distributed Management

```lua

local update_servers = task("update_servers")
    :command(function() return pkg.update({}) end)
    :delegate_to("prod-server-1")
    :build()

local install_monitoring = task("install_monitoring")
    :command(function()
        return pkg.install({packages = {"htop", "iotop", "nethogs"}})
    end)
    :delegate_to("prod-server-1")
    :depends_on({"update_servers"})
    :build()

workflow.define("setup_monitoring")
    :tasks({ update_servers, install_monitoring })
```

### System Audit

```lua

local audit = task("audit")
    :command(function()
        log.info("📊 System Audit")
        log.info(string.rep("=", 60))
        
        local manager, _ = pkg.get_manager({})
        log.info("Manager: " .. (manager or "unknown"))
        
        local _, packages = pkg.list({})
        local count = 0
        if type(packages) == "table" then
            for _ in pairs(packages) do count = count + 1 end
        end
        log.info("Packages: " .. count)
        
        local critical = {"openssl", "curl"}
        for _, p in ipairs(critical) do
            local installed, _ = pkg.is_installed({package = p})
            log.info((installed and "✅" or "❌") .. " " .. p)
        end
        
        return true, "OK"
    end)
    :build()

workflow.define("audit")
    :tasks({ audit })
```

## 🚀 Best Practices

1. **Update before installing:**
   ```lua
   pkg.update({})
   pkg.install({packages = "package"})
   ```

2. **Check before installing:**
   ```lua
   local installed, _ = pkg.is_installed({package = "git"})
   if not installed then
       pkg.install({packages = "git"})
   end
   ```

3. **Cleanup after operations:**
   ```lua
   pkg.clean({})
   pkg.autoremove({})
   ```

4. **Use delegate_to for remote:**
   ```lua
   :delegate_to("server-name")
   ```

## ⚠️ Platform Notes

- **Linux**: Requires sudo
- **macOS**: Homebrew doesn't need sudo
- **Arch**: Uses pacman syntax
- **openSUSE**: Uses zypper

## 🔗 See Also

- [exec Module](exec.md)
- [Modern DSL Guide](../modern-dsl/overview.md)
- [Distributed Agents](../distributed.md)
