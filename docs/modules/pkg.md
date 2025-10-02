# 📦 Package Manager Module

The `pkg` module provides comprehensive cross-platform package management functionality. It automatically detects the system's package manager and provides a unified interface for managing packages.

## 🚀 Quick Start

### Modern Syntax (Recommended)

```lua
-- pkg is available globally, no require needed!
task("install_tools")
  :command(function()
    pkg.update()
    pkg.install({"git", "curl", "vim"})
    return true
  end)
  :build()
```

### Classic Syntax (Still Supported)

```lua
local pkg = require("pkg")

task("install_tools")
  :command(function()
    pkg.update()
    pkg.install({"git", "curl", "vim"})
    return true
  end)
  :build()
```

> 💡 **Tip**: Use the modern syntax! All built-in modules (`pkg`, `docker`, `systemd`, `git`, `terraform`, etc.) are available globally without `require()`.

---

## 🎯 Supported Package Managers

| Package Manager | Systems | Auto-Detected |
|----------------|---------|---------------|
| **apt / apt-get** | Debian, Ubuntu | ✅ |
| **yum / dnf** | RHEL, CentOS, Fedora, Amazon Linux | ✅ |
| **pacman** | Arch Linux, Manjaro | ✅ |
| **zypper** | openSUSE, SLES | ✅ |
| **brew** | macOS (Homebrew) | ✅ |

---

## 📚 API Reference

### Installation & Removal

#### `pkg.install(packages)`

Install one or more packages.

**Parameters:**
- `packages`: `string` (single) or `table` (multiple)

**Returns:**
- `success` (boolean)
- `output` (string)

**Examples:**

```lua
-- Single package
local ok, msg = pkg.install("nginx")

-- Multiple packages
local ok, msg = pkg.install({"git", "curl", "wget", "vim"})

-- In Modern DSL task
task("install_nginx")
  :command(function()
    log.info("Installing nginx...")
    local success, output = pkg.install("nginx")
    if success then
      log.success("✅ Nginx installed!")
    else
      log.error("❌ Failed: " .. output)
    end
    return success
  end)
  :build()
```

#### `pkg.remove(packages)`

Remove one or more packages.

**Parameters:**
- `packages`: `string` (single) or `table` (multiple)

**Returns:**
- `success` (boolean)
- `output` (string)

**Examples:**

```lua
-- Remove single package
pkg.remove("apache2")

-- Remove multiple packages
pkg.remove({"apache2", "php-fpm"})
```

---

### Package Information

#### `pkg.is_installed(package)`

Check if a package is installed.

**Parameters:**
- `package`: `string` - Package name

**Returns:**
- `installed` (boolean)

**Example:**

```lua
if pkg.is_installed("nginx") then
  log.info("✅ Nginx is installed")
else
  log.warn("⚠️  Nginx not found, installing...")
  pkg.install("nginx")
end
```

#### `pkg.info(package)`

Get detailed information about a package.

**Parameters:**
- `package`: `string` - Package name

**Returns:**
- `success` (boolean)
- `info` (string) - Package information

**Example:**

```lua
local ok, info = pkg.info("nginx")
if ok then
  print(info)  -- Shows version, description, dependencies, etc.
end
```

#### `pkg.version(package)`

Get the installed version of a package.

**Parameters:**
- `package`: `string` - Package name

**Returns:**
- `success` (boolean)
- `version` (string)

**Example:**

```lua
local ok, ver = pkg.version("nginx")
if ok then
  log.info("Nginx version: " .. ver)
end
```

#### `pkg.deps(package)`

List package dependencies.

**Parameters:**
- `package`: `string` - Package name

**Returns:**
- `success` (boolean)
- `dependencies` (table or string)

**Example:**

```lua
local ok, deps = pkg.deps("nginx")
if ok and type(deps) == "table" then
  for _, dep in ipairs(deps) do
    print("  - " .. dep)
  end
end
```

---

### Repository Management

#### `pkg.update()`

Update the package cache/repository list.

**Returns:**
- `success` (boolean)
- `output` (string)

**Example:**

```lua
task("update_cache")
  :command(function()
    log.info("Updating package cache...")
    local ok, msg = pkg.update()
    return ok
  end)
  :timeout("2m")
  :build()
```

#### `pkg.upgrade()`

Upgrade all installed packages to their latest versions.

**Returns:**
- `success` (boolean)
- `output` (string)

**Example:**

```lua
task("upgrade_system")
  :command(function()
    pkg.update()
    local ok, msg = pkg.upgrade()
    if ok then
      log.success("✅ System upgraded!")
    end
    return ok
  end)
  :timeout("30m")
  :build()
```

#### `pkg.search(query)`

Search for packages in repositories.

**Parameters:**
- `query`: `string` - Search term

**Returns:**
- `success` (boolean)
- `results` (string) - Search results

**Example:**

```lua
local ok, results = pkg.search("python3")
if ok then
  print(results)
end
```

---

### Maintenance

#### `pkg.clean()`

Clean the package manager cache.

**Returns:**
- `success` (boolean)
- `output` (string)

**Example:**

```lua
pkg.clean()  -- Free up disk space
```

#### `pkg.autoremove()`

Remove packages that were automatically installed as dependencies but are no longer needed.

**Returns:**
- `success` (boolean)
- `output` (string)

**Example:**

```lua
task("cleanup")
  :command(function()
    pkg.autoremove()
    pkg.clean()
    return true
  end)
  :build()
```

#### `pkg.list()`

List all installed packages.

**Returns:**
- `success` (boolean)
- `packages` (table or string) - List of installed packages

**Example:**

```lua
local ok, packages = pkg.list()
if ok and type(packages) == "table" then
  log.info("Installed packages: " .. #packages)
end
```

---

### Advanced Functions

#### `pkg.get_manager()`

Get the detected package manager name.

**Returns:**
- `manager` (string) - e.g., "apt", "yum", "pacman", "brew"

**Example:**

```lua
local pm = pkg.get_manager()
log.info("Using package manager: " .. pm)
```

#### `pkg.which(executable)`

Find the full path of an executable.

**Parameters:**
- `executable`: `string` - Command name

**Returns:**
- `success` (boolean)
- `path` (string) - Full path or error message

**Example:**

```lua
local ok, path = pkg.which("nginx")
if ok then
  log.info("Nginx binary: " .. path)
end
```

#### `pkg.install_local(file)`

Install a package from a local file.

**Parameters:**
- `file`: `string` - Path to package file (.deb, .rpm, etc.)

**Returns:**
- `success` (boolean)
- `output` (string)

**Example:**

```lua
pkg.install_local("/tmp/my-app_1.0.0_amd64.deb")
```

---

## 🎯 Complete Examples

### Development Environment Setup

```lua
task("setup_dev_env")
  :description("Install development tools")
  :command(function()
    log.info("🚀 Setting up development environment...")
    
    -- Update cache
    log.info("📦 Updating package cache...")
    pkg.update()
    
    -- Install dev tools
    local tools = {
      "git",
      "curl",
      "wget",
      "vim",
      "build-essential",  -- apt
      "htop",
      "jq"
    }
    
    log.info("🛠️  Installing tools...")
    local ok, msg = pkg.install(tools)
    
    if ok then
      log.success("✅ All tools installed!")
      
      -- Verify installations
      for _, tool in ipairs(tools) do
        if pkg.is_installed(tool) then
          local _, ver = pkg.version(tool)
          log.info("  ✓ " .. tool .. " " .. (ver or ""))
        end
      end
    else
      log.error("❌ Installation failed: " .. msg)
      return false
    end
    
    return true
  end)
  :timeout("10m")
  :build()
```

### Conditional Package Management

```lua
task("ensure_nginx")
  :description("Ensure Nginx is installed and running")
  :command(function()
    -- Check if already installed
    if pkg.is_installed("nginx") then
      log.info("✅ Nginx already installed")
      local _, ver = pkg.version("nginx")
      log.info("   Version: " .. ver)
    else
      log.info("Installing Nginx...")
      local ok, msg = pkg.install("nginx")
      if not ok then
        log.error("Failed: " .. msg)
        return false
      end
    end
    
    -- Start service (assuming systemd)
    local systemd = require("systemd")
    systemd.enable("nginx")
    systemd.start("nginx")
    
    return true
  end)
  :build()
```

### Multi-Package Workflow

```lua
task("update")
  :command(function()
    return pkg.update()
  end)
  :build()

task("install_web_stack")
  :depends_on("update")
  :command(function()
    return pkg.install({"nginx", "php-fpm", "mysql-server"})
  end)
  :build()

task("cleanup")
  :depends_on("install_web_stack")
  :command(function()
    pkg.autoremove()
    pkg.clean()
    return true
  end)
  :build()
```

### Cross-Platform Package Management

```lua
task("install_docker")
  :command(function()
    local pm = pkg.get_manager()
    log.info("Package manager: " .. pm)
    
    pkg.update()
    
    if pm == "apt" then
      pkg.install({"docker.io", "docker-compose"})
    elseif pm == "yum" or pm == "dnf" then
      pkg.install({"docker", "docker-compose"})
    elseif pm == "pacman" then
      pkg.install({"docker", "docker-compose"})
    elseif pm == "brew" then
      pkg.install("docker")
    end
    
    return true
  end)
  :build()
```

---

## 🔍 Error Handling

```lua
task("safe_install")
  :command(function()
    local ok, msg = pkg.install("nginx")
    
    if not ok then
      log.error("Installation failed: " .. msg)
      
      -- Try alternative
      log.info("Trying alternative package...")
      ok, msg = pkg.install("nginx-full")
    end
    
    return ok, msg
  end)
  :on_error(function(err)
    log.error("Task failed: " .. err)
    -- Cleanup or rollback here
  end)
  :retry(3)  -- Retry up to 3 times
  :build()
```

---

## 💡 Best Practices

### 1. Always Update Before Installing

```lua
-- ✅ Good
pkg.update()
pkg.install("package")

-- ❌ Bad
pkg.install("package")  -- May get outdated version
```

### 2. Handle Installation Failures

```lua
-- ✅ Good
local ok, msg = pkg.install("nginx")
if not ok then
  log.error(msg)
  return false
end

-- ❌ Bad
pkg.install("nginx")  -- Ignores failures
```

### 3. Check Before Installing

```lua
-- ✅ Good
if not pkg.is_installed("nginx") then
  pkg.install("nginx")
end

-- ❌ Bad (slower, may fail if already installed)
pkg.install("nginx")
```

### 4. Use Timeouts for Long Operations

```lua
task("upgrade_all")
  :command(function()
    return pkg.upgrade()
  end)
  :timeout("30m")  -- ✅ Prevent hanging
  :build()
```

---

## 🐛 Troubleshooting

### Permission Denied

Most package operations require root:

```bash
# Run with sudo
sudo sloth-runner run -f workflow.sloth
```

### Package Not Found

```lua
-- Search first
local ok, results = pkg.search("package-name")
print(results)
```

### Lock File Errors

```lua
-- Wait and retry
task("install")
  :command(function()
    return pkg.install("package")
  end)
  :retry(3)
  :retry_delay("30s")
  :build()
```

---

## 🔗 Related Modules

- [systemd](/modules/systemd/) - Service management
- [docker](/modules/docker/) - Container management
- [terraform](/modules/terraform/) - Infrastructure as Code

---

## 📖 See Also

- [Getting Started](/en/getting-started/)
- [Modern DSL](/modern-dsl/introduction/)
- [Examples Repository](https://github.com/chalkan3-sloth/sloth-runner/tree/main/examples)

---

**Package management made simple across all platforms!** 📦✨
    ```lua
    local pkg = require("pkg")
    
    local install_tools = task("install_tools")
        :description("Install development tools")
        :command(function(this, params)
            log.info("Installing tools...")
            
            -- Install multiple packages
            local tools = {"git", "curl", "wget", "vim"}
            local success, output = pkg.install(tools)
            
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
    local pkg = require("pkg")
    
    local install_on_agent = task("install_on_agent")
        :description("Install packages on remote agent")
        :command(function(this, params)
            log.info("Installing on remote agent...")
            
            local success, output = pkg.install({"htop", "ncdu"})
            
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

#### `pkg.remove(packages)`

Removes one or more packages.

**Parameters:**
- `packages`: String or Table

**Returns:**
- `success` (boolean), `output` (string)

**Example:**

```lua
local pkg = require("pkg")

local cleanup = task("cleanup")
    :description("Remove unnecessary packages")
    :command(function(this, params)
        local packages = {"package1", "package2"}
        local success, output = pkg.remove(packages)
        
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

#### `pkg.search(query)`

Searches for packages.

**Example:**

```lua
local pkg = require("pkg")

local search_python = task("search_python")
    :description("Search for Python packages")
    :command(function(this, params)
        local success, results = pkg.search("python3")
        
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

#### `pkg.info(package)`

Gets package information.

**Example:**

```lua
local success, info = pkg.info("curl")
if success then
    log.info("Package info:\n" .. info)
end
```

#### `pkg.list()`

Lists installed packages.

**Returns:** `success` (boolean), `packages` (table)

**Example:**

```lua
local success, packages = pkg.list()
if success and type(packages) == "table" then
    local count = 0
    for _ in pairs(packages) do count = count + 1 end
    log.info("📦 Total: " .. count .. " packages")
end
```

### System Maintenance

#### `pkg.update()`

Updates package cache.

**Example:**

```lua
local update_cache = task("update_cache")
    :description("Update package cache")
    :command(function(this, params)
        log.info("Updating...")
        return pkg.update()
    end)
    :timeout("120s")
    :build()
```

#### `pkg.upgrade()`

Upgrades all packages.

#### `pkg.clean()`

Cleans package cache.

#### `pkg.autoremove()`

Removes unused dependencies.

**Example:**

```lua
local maintenance = task("maintenance")
    :description("System maintenance")
    :command(function(this, params)
        -- Update
        pkg.update()
        
        -- Upgrade
        pkg.upgrade()
        
        -- Clean
        pkg.clean()
        pkg.autoremove()
        
        return true, "Maintenance complete"
    end)
    :timeout("600s")
    :build()
```

### Advanced Functions

#### `pkg.is_installed(package)`

Checks if installed.

**Example:**

```lua
local pkg = require("pkg")

local check_requirements = task("check_requirements")
    :description("Check required packages")
    :command(function(this, params)
        local required = {"git", "curl", "wget"}
        local missing = {}
        
        for _, pkg_name in ipairs(required) do
            local installed, _ = pkg.is_installed(pkg_name)
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

#### `pkg.get_manager()`

Returns package manager name.

**Example:**

```lua
local manager, err = pkg.get_manager()
log.info("Manager: " .. (manager or "unknown"))
```

#### `pkg.which(executable)`

Finds executable path.

**Example:**

```lua
local path, err = pkg.which("git")
if path then
    log.info("Git at: " .. path)
end
```

#### `pkg.version(package)`

Gets package version.

#### `pkg.deps(package)`

Lists dependencies.

#### `pkg.install_local(filepath)`

Installs from local file (.deb, .rpm).

## 🎯 Complete Examples

### Development Environment Setup

```lua
local pkg = require("pkg")

local update = task("update")
    :command(function() return pkg.update() end)
    :build()

local install_tools = task("install_tools")
    :command(function()
        local tools = {"git", "curl", "wget", "vim", "htop"}
        return pkg.install(tools)
    end)
    :depends_on({"update"})
    :build()

local verify = task("verify")
    :command(function()
        for _, tool in ipairs({"git", "curl"}) do
            if pkg.is_installed(tool) then
                local path = pkg.which(tool)
                log.info("✅ " .. tool .. " (" .. path .. ")")
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
local pkg = require("pkg")

local update_servers = task("update_servers")
    :command(function() return pkg.update() end)
    :delegate_to("prod-server-1")
    :build()

local install_monitoring = task("install_monitoring")
    :command(function()
        return pkg.install({"htop", "iotop", "nethogs"})
    end)
    :delegate_to("prod-server-1")
    :depends_on({"update_servers"})
    :build()

workflow.define("setup_monitoring")
    :tasks({ update_servers, install_monitoring })
```

### System Audit

```lua
local pkg = require("pkg")

local audit = task("audit")
    :command(function()
        log.info("📊 System Audit")
        log.info("=".rep(60))
        
        local manager = pkg.get_manager()
        log.info("Manager: " .. manager)
        
        local _, packages = pkg.list()
        local count = 0
        for _ in pairs(packages) do count = count + 1 end
        log.info("Packages: " .. count)
        
        local critical = {"openssl", "curl"}
        for _, p in ipairs(critical) do
            local installed = pkg.is_installed(p)
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
   pkg.update()
   pkg.install("package")
   ```

2. **Check before installing:**
   ```lua
   if not pkg.is_installed("git") then
       pkg.install("git")
   end
   ```

3. **Cleanup after operations:**
   ```lua
   pkg.clean()
   pkg.autoremove()
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
