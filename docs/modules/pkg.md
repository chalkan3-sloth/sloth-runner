# 📦 Package Manager Module

The `pkg` module provides comprehensive cross-platform package management functionality. It automatically detects the system's package manager and provides a unified interface for managing packages.

## 🚀 Quick Start

```lua
-- pkg is available globally, no require needed!
task("install_tools")
  :description("Install essential development tools")
  :command(function(this, params)
    pkg.update()
    pkg.install({"git", "curl", "vim"})
    return true, "Tools installed"
  end)
  :build()
```

> 💡 **Tip**: The `pkg` module is available globally! All infrastructure modules are automatically available without `require()`.

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
  :description("Install nginx web server")
  :command(function(this, params)
    log.info("Installing nginx...")
    local success, output = pkg.install("nginx")
    if success then
      log.success("✅ Nginx installed!")
      return true, "Nginx installed successfully"
    else
      log.error("❌ Failed: " .. output)
      return false, "Failed to install nginx: " .. output
    end
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
  :description("Update package cache")
  :command(function(this, params)
    log.info("Updating package cache...")
    local ok, msg = pkg.update()
    if ok then
      return true, "Cache updated successfully"
    else
      return false, "Failed to update cache: " .. msg
    end
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
  :description("Upgrade all system packages")
  :command(function(this, params)
    pkg.update()
    local ok, msg = pkg.upgrade()
    if ok then
      log.success("✅ System upgraded!")
      return true, "System upgraded successfully"
    else
      return false, "Failed to upgrade system: " .. msg
    end
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
  :description("Clean up package cache and unused dependencies")
  :command(function(this, params)
    pkg.autoremove()
    pkg.clean()
    return true, "Cleanup completed"
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
  :command(function(this, params)
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
      return true, "All development tools installed successfully"
    else
      log.error("❌ Installation failed: " .. msg)
      return false, "Installation failed: " .. msg
    end
  end)
  :timeout("10m")
  :build()
```

### Conditional Package Management

```lua
task("ensure_nginx")
  :description("Ensure Nginx is installed and running")
  :command(function(this, params)
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
        return false, "Failed to install nginx: " .. msg
      end
    end

    -- Start service (assuming systemd)
    systemd.enable("nginx")
    systemd.start("nginx")

    return true, "Nginx is installed and running"
  end)
  :build()
```

### Multi-Package Workflow

```lua
task("update")
  :description("Update package cache")
  :command(function(this, params)
    local ok, msg = pkg.update()
    if ok then
      return true, "Package cache updated"
    else
      return false, "Failed to update cache: " .. msg
    end
  end)
  :build()

task("install_web_stack")
  :description("Install web server stack")
  :depends_on("update")
  :command(function(this, params)
    local ok, msg = pkg.install({"nginx", "php-fpm", "mysql-server"})
    if ok then
      return true, "Web stack installed successfully"
    else
      return false, "Failed to install web stack: " .. msg
    end
  end)
  :build()

task("cleanup")
  :description("Clean up after installation")
  :depends_on("install_web_stack")
  :command(function(this, params)
    pkg.autoremove()
    pkg.clean()
    return true, "Cleanup completed"
  end)
  :build()
```

### Cross-Platform Package Management

```lua
task("install_docker")
  :description("Install Docker on any platform")
  :command(function(this, params)
    local pm = pkg.get_manager()
    log.info("Package manager: " .. pm)

    pkg.update()

    local ok, msg
    if pm == "apt" then
      ok, msg = pkg.install({"docker.io", "docker-compose"})
    elseif pm == "yum" or pm == "dnf" then
      ok, msg = pkg.install({"docker", "docker-compose"})
    elseif pm == "pacman" then
      ok, msg = pkg.install({"docker", "docker-compose"})
    elseif pm == "brew" then
      ok, msg = pkg.install("docker")
    else
      return false, "Unsupported package manager: " .. pm
    end

    if ok then
      return true, "Docker installed successfully"
    else
      return false, "Failed to install Docker: " .. msg
    end
  end)
  :build()
```

---

## 🔍 Error Handling

```lua
task("safe_install")
  :description("Install nginx with fallback")
  :command(function(this, params)
    local ok, msg = pkg.install("nginx")

    if not ok then
      log.error("Installation failed: " .. msg)

      -- Try alternative
      log.info("Trying alternative package...")
      ok, msg = pkg.install("nginx-full")
    end

    if ok then
      return true, "Nginx installed successfully"
    else
      return false, "Failed to install nginx: " .. msg
    end
  end)
  :on_error(function(this, params, err)
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
  :description("Upgrade all packages")
  :command(function(this, params)
    local ok, msg = pkg.upgrade()
    if ok then
      return true, "All packages upgraded"
    else
      return false, "Failed to upgrade: " .. msg
    end
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
  :description("Install package with retry")
  :command(function(this, params)
    local ok, msg = pkg.install("package")
    if ok then
      return true, "Package installed"
    else
      return false, "Failed to install: " .. msg
    end
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
