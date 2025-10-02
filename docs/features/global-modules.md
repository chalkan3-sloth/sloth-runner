# 🌍 Global Modules - No require() Needed!

Starting from the latest version of Sloth Runner, **all native infrastructure modules are available globally** without needing to call `require()`.

## ✨ What Changed?

### ❌ Before (Old Way)
```lua
-- Had to require every module
local pkg = require("pkg")
local user = require("user")
local systemd = require("systemd")
local file_ops = require("file_ops")

task("setup")
    :command(function()
        pkg.install({ name = "nginx" })
        user.create({ name = "webapp" })
        systemd.enable("nginx")
    end)
    :build()
```

### ✅ After (New Way - Recommended)
```lua
-- All native modules available globally!
task("setup")
    :command(function()
        pkg.install({ name = "nginx" })
        user.create({ name = "webapp" })
        systemd.enable("nginx")
        file_ops.copy({
            src = "nginx.conf",
            dest = "/etc/nginx/nginx.conf"
        })
    end)
    :build()
```

## 📦 Available Global Modules

All these modules are **automatically available** in your Sloth scripts:

| Module | Description |
|--------|-------------|
| `pkg` | Package management (apt, yum, pacman, etc.) |
| `user` | User and group management |
| `ssh` | SSH connections and file transfers |
| `file_ops` | File operations (Ansible-like) |
| `systemd` | Systemd service management |
| `state` | Persistent state management |
| `terraform` | Terraform infrastructure |
| `pulumi` | Pulumi infrastructure |
| `kubernetes` | Kubernetes operations |
| `helm` | Helm chart management |
| `salt` | Salt Stack integration |
| `infra_test` | Infrastructure testing |

## 🎯 Best Practices

### ✅ DO: Use Global Modules Directly

```lua
task("deploy_app")
    :command(function()
        -- Clean and simple!
        pkg.install({ name = "nginx" })
        systemd.enable("nginx")
        systemd.start("nginx")
        
        -- Test the deployment
        infra_test.service_is_running("nginx")
        infra_test.port_is_listening(80)
    end)
    :build()
```

### ✅ DO: Combine with :delegate_to()

```lua
task("deploy_remote")
    :delegate_to("prod-server-01")
    :command(function()
        -- All modules work remotely too!
        pkg.update()
        pkg.install({ name = "postgresql" })
        systemd.enable("postgresql")
    end)
    :build()
```

### ⚠️ STILL NEED require(): Non-Infrastructure Modules

Some modules still require `require()` because they're not infrastructure-focused:

```lua
task("complex_workflow")
    :command(function()
        -- These STILL need require()
        local git = require("git")
        local http = require("http")
        local data = require("data")
        local crypto = require("crypto")
        
        -- But native modules don't!
        pkg.install({ name = "git" })
        user.create({ name = "developer" })
        
        -- Mix and match as needed
        git.clone("https://github.com/user/repo", "/tmp/repo")
    end)
    :build()
```

**Modules that still need `require()`:**
- `git` - Git operations
- `http` - HTTP requests
- `data` - Data transformation
- `crypto` - Cryptography
- `math` - Advanced math
- `strings` - String manipulation
- `time` - Time operations
- `observability` - Monitoring
- `security` - Security operations

## 🔄 Migration Guide

If you have existing scripts with `require()` calls for native modules, they will **still work**! The old way is backwards compatible.

### Example Migration

**Old Script (Still Works):**
```lua
local pkg = require("pkg")
local systemd = require("systemd")

task("install_nginx")
    :command(function()
        pkg.install({ name = "nginx" })
        systemd.start("nginx")
    end)
    :build()
```

**New Script (Recommended):**
```lua
-- Just remove the require() lines!
task("install_nginx")
    :command(function()
        pkg.install({ name = "nginx" })
        systemd.start("nginx")
    end)
    :build()
```

## 🎨 Benefits

1. **Less Boilerplate**: No need to require common modules
2. **Cleaner Code**: Focus on what you're doing, not imports
3. **Better DX**: Faster to write infrastructure automation
4. **Backwards Compatible**: Old scripts keep working
5. **Consistent**: All native modules follow the same pattern

## 📚 Complete Example

### Complete Server Setup

```lua
task("setup_web_server")
    :description("Complete web server setup")
    :delegate_to("web-01")
    :command(function()
        -- Update system
        pkg.update()
        
        -- Install packages
        pkg.install({
            name = "nginx",
            version = "latest"
        })
        pkg.install({ name = "certbot" })
        
        -- Create user
        user.create({
            name = "webapp",
            home = "/var/www",
            shell = "/bin/bash",
            groups = {"www-data"}
        })
        
        -- Copy config
        file_ops.template({
            src = "templates/nginx.conf.j2",
            dest = "/etc/nginx/nginx.conf",
            vars = {
                server_name = "example.com",
                port = 80
            },
            mode = 0o644
        })
        
        -- Start service
        systemd.enable("nginx")
        systemd.start("nginx")
        
        -- Validate
        infra_test.service_is_running("nginx")
        infra_test.port_is_listening(80)
        infra_test.file_exists("/etc/nginx/nginx.conf")
        
        log.info("✅ Web server setup complete!")
        return true
    end)
    :build()
```

## 🔗 Related Documentation

- [Package Manager Module](../modules/pkg.md)
- [User Management Module](../modules/user.md)
- [Systemd Module](../modules/systemd.md)
- [File Operations Module](../modules/file_ops.md)
- [SSH Module](../modules/ssh.md)
- [Infrastructure Testing Module](../modules/infra_test.md)
- [Modern DSL Guide](../modern-dsl/README.md)
