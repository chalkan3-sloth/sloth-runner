# 📦 Stow Module

The `stow` module provides native GNU Stow integration for managing dotfiles and symlink farms in Sloth Runner. It's a **global module** (no `require()` needed) with full idempotency and task user support.

## Features

- ✅ **Automatic target directory creation** with proper ownership
- ✅ **Idempotent operations** - safe to run multiple times
- ✅ **Task user integration** - respects `:user()` directive
- ✅ **Multiple stow operations** - link, unlink, restow
- ✅ **Advanced options** - no-folding, verbose, and more

## Functions

### `stow.link()`

Creates symlinks for a package (stow operation).

**Parameters:**
```lua
{
    package = "package_name",      -- Required: package/directory to stow
    source_dir = "/path/to/stow",  -- Required: stow directory
    target_dir = "/path/to/target", -- Required: target directory
    create_target = true,          -- Optional: create target dir if missing (default: true)
    verbose = false,               -- Optional: verbose output
    no_folding = false            -- Optional: don't fold directories
}
```

**Returns:** `success (bool), message (string)`

**Example:**
```lua
local ok, msg = stow.link({
    package = "zsh",
    source_dir = "/home/user/dotfiles",
    target_dir = "/home/user",
    create_target = true,
    verbose = true
})

if not ok then
    return false, msg
end
```

**With automatic directory creation:**
```lua
-- This will create /home/user/.zsh if it doesn't exist
-- and set ownership to the task user
local ok, msg = stow.link({
    package = ".",
    source_dir = "/home/user/dotfiles/zsh",
    target_dir = "/home/user/.zsh",
    create_target = true  -- Creates dir with task user ownership
})
```

### `stow.unlink()`

Removes symlinks for a package (unstow operation).

**Parameters:**
```lua
{
    package = "package_name",      -- Required
    source_dir = "/path/to/stow",  -- Required
    target_dir = "/path/to/target", -- Required
    verbose = false                -- Optional
}
```

**Returns:** `success (bool), message (string)`

**Example:**
```lua
local ok, msg = stow.unlink({
    package = "vim",
    source_dir = "/home/user/dotfiles",
    target_dir = "/home/user"
})
```

### `stow.restow()`

Removes and re-creates symlinks for a package (useful for updates).

**Parameters:**
```lua
{
    package = "package_name",      -- Required
    source_dir = "/path/to/stow",  -- Required
    target_dir = "/path/to/target", -- Required
    verbose = false,               -- Optional
    no_folding = false            -- Optional
}
```

**Returns:** `success (bool), message (string)`

**Example:**
```lua
-- Refresh all links for the package
local ok, msg = stow.restow({
    package = "zshrc",
    source_dir = "/home/user/dotfiles",
    target_dir = "/home/user",
    verbose = true
})
```

### `stow.ensure_target()` 🆕

Ensures a target directory exists with proper ownership and permissions.

**Parameters:**
```lua
{
    path = "/path/to/directory",  -- Required: directory path
    owner = "username",           -- Optional: owner (uses task user if not specified)
    mode = "0755"                -- Optional: permissions in octal (default: "0755")
}
```

**Returns:** `success (bool), message (string)`

**Example:**
```lua
-- Create directory as task user
local ok, msg = stow.ensure_target({
    path = "/home/user/.config/nvim"
})

-- Create with specific owner and permissions
local ok, msg = stow.ensure_target({
    path = "/home/user/.local/bin",
    owner = "user",
    mode = "0700"
})
```

## Complete Examples

### Basic Dotfiles Setup

```lua
local stow_dotfiles = task("stow-dotfiles")
    :description("Stow all dotfiles")
    :user("myuser")
    :command(function(this, params)
        local packages = { "zsh", "vim", "tmux", "git" }

        for _, pkg in ipairs(packages) do
            local ok, msg = stow.link({
                package = pkg,
                source_dir = "/home/myuser/dotfiles",
                target_dir = "/home/myuser",
                create_target = true
            })

            if ok then
                log.info("✅ " .. pkg .. " stowed")
            else
                log.error("❌ " .. pkg .. ": " .. msg)
                return false, msg
            end
        end

        return true, "All dotfiles stowed"
    end)
    :build()
```

### Nested Directory Structure

```lua
local stow_zsh = task("stow-zsh-config")
    :description("Stow zsh configuration into .zsh directory")
    :user("igor")
    :command(function(this, params)
        -- Ensure target directory exists
        local ok_dir, msg_dir = stow.ensure_target({
            path = "/home/igor/.zsh",
            owner = "igor"
        })

        if not ok_dir then
            return false, "Failed to create .zsh: " .. msg_dir
        end

        -- Stow the configuration
        local ok_stow, msg_stow = stow.link({
            package = ".",
            source_dir = "/home/igor/dotfiles/zsh",
            target_dir = "/home/igor/.zsh",
            no_folding = false
        })

        if not ok_stow then
            return false, "Failed to stow: " .. msg_stow
        end

        return true, "Zsh config stowed to .zsh directory"
    end)
    :build()
```

### User Environment Setup

```lua
workflow
    .define("user_dotfiles_setup")
    :description("Complete user dotfiles setup")
    :tasks({
        task("install-deps")
            :delegate_to("server1")
            :command(function()
                pkg.install({ packages = { "stow", "git", "zsh" } })
                return true
            end)
            :build(),

        task("create-user")
            :delegate_to("server1")
            :command(function()
                user.create({
                    username = "myuser",
                    shell = "/bin/zsh",
                    create_home = true
                })
                return true
            end)
            :build(),

        task("clone-dotfiles")
            :delegate_to("server1")
            :user("myuser")
            :command(function()
                exec.run("git clone https://github.com/user/dotfiles.git ~/dotfiles")
                return true
            end)
            :build(),

        task("stow-all")
            :delegate_to("server1")
            :user("myuser")
            :command(function()
                -- Stow zsh to .zsh directory
                stow.link({
                    package = ".",
                    source_dir = "/home/myuser/dotfiles/zsh",
                    target_dir = "/home/myuser/.zsh",
                    create_target = true
                })

                -- Stow zshrc to home
                stow.link({
                    package = "zshrc",
                    source_dir = "/home/myuser/dotfiles",
                    target_dir = "/home/myuser"
                })

                return true, "All dotfiles stowed"
            end)
            :build()
    })
```

### Multiple Packages with Error Handling

```lua
local stow_multiple = task("stow-multiple")
    :user("myuser")
    :command(function(this, params)
        local packages = {
            { name = "zsh", target = "/home/myuser" },
            { name = "vim", target = "/home/myuser" },
            { name = "scripts", target = "/home/myuser/.local/bin" },
        }

        local results = {}
        local failed = {}

        for _, pkg_info in ipairs(packages) do
            local ok, msg = stow.link({
                package = pkg_info.name,
                source_dir = "/home/myuser/dotfiles",
                target_dir = pkg_info.target,
                create_target = true,
                verbose = true
            })

            if ok then
                table.insert(results, pkg_info.name)
                log.info("✅ " .. pkg_info.name .. ": " .. msg)
            else
                table.insert(failed, pkg_info.name)
                log.error("❌ " .. pkg_info.name .. ": " .. msg)
            end
        end

        if #failed > 0 then
            return false, "Failed to stow: " .. table.concat(failed, ", ")
        end

        return true, "Successfully stowed: " .. table.concat(results, ", ")
    end)
    :build()
```

## Best Practices

### 1. **Always use `create_target = true`** for new setups
```lua
-- Good: Automatically creates missing directories
stow.link({
    package = "zsh",
    source_dir = "~/dotfiles",
    target_dir = "~/.config/zsh",
    create_target = true
})
```

### 2. **Use `:user()` directive** for proper ownership
```lua
task("stow-config")
    :user("myuser")  -- All stow operations will run as myuser
    :command(function()
        stow.link({ ... })
    end)
    :build()
```

### 3. **Explicitly create complex directory structures**
```lua
-- For complex structures, ensure directories first
stow.ensure_target({ path = "/home/user/.config/nvim" })
stow.ensure_target({ path = "/home/user/.local/share" })

-- Then stow
stow.link({ package = "nvim", ... })
```

### 4. **Use `restow` for updates**
```lua
-- When dotfiles change, use restow
stow.restow({
    package = "vim",
    source_dir = "~/dotfiles",
    target_dir = "~"
})
```

### 5. **Check results and log appropriately**
```lua
local ok, msg = stow.link({ ... })

if ok then
    log.info("✅ " .. msg)
else
    log.error("❌ " .. msg)
    return false, msg
end
```

## Idempotency

All stow operations are **fully idempotent**:

- `stow.link()` - Checks if links already exist before creating
- `stow.unlink()` - Only removes links if they exist
- `stow.restow()` - Safe to run multiple times
- `stow.ensure_target()` - Only creates directory if missing

**Example:**
```lua
-- Safe to run multiple times
stow.link({
    package = "zsh",
    source_dir = "/home/user/dotfiles",
    target_dir = "/home/user"
})
-- First run: Creates symlinks
-- Second run: Detects existing links, returns success
```

## Task User Integration

The stow module respects the task `:user()` directive:

```lua
task("stow-as-user")
    :user("igor")  -- Run as igor
    :command(function()
        -- This will:
        -- 1. Create /home/igor/.zsh owned by igor
        -- 2. Run stow as igor
        stow.link({
            package = ".",
            source_dir = "/home/igor/dotfiles/zsh",
            target_dir = "/home/igor/.zsh"
        })
    end)
    :build()
```

## Troubleshooting

### Links not created
```bash
# Check stow is installed
pkg.install({ packages = { "stow" } })

# Check source directory exists
log.info("Source: " .. exec.run("ls -la /home/user/dotfiles"))

# Use verbose mode
stow.link({ ..., verbose = true })
```

### Permission denied
```bash
# Ensure proper task user
task("fix-perms")
    :user("targetuser")  # Must match target directory owner
    :command(function()
        stow.link({ ... })
    end)
    :build()
```

### Directory already exists
```bash
# Use ensure_target to handle existing directories
stow.ensure_target({ path = "/home/user/.config" })
stow.link({
    package = "config",
    target_dir = "/home/user/.config",
    create_target = false  # Already ensured above
})
```

## Migration from Manual exec.run()

**Before (manual stow):**
```lua
exec.run("sudo -u igor -- /bin/sh -c 'mkdir -p /home/igor/.zsh'")
exec.run("sudo -u igor -- /bin/sh -c 'stow -d /home/igor/dotfiles/zsh -t /home/igor/.zsh .'")
```

**After (using stow module):**
```lua
stow.link({
    package = ".",
    source_dir = "/home/igor/dotfiles/zsh",
    target_dir = "/home/igor/.zsh",
    create_target = true  -- Handles mkdir and ownership
})
```

**Benefits:**
- ✅ Automatic directory creation
- ✅ Proper ownership handling
- ✅ Idempotent by default
- ✅ Better error messages
- ✅ Cleaner code

## See Also

- [file_ops module](./file_ops.md) - For file operations
- [user module](./user.md) - For user management
- [exec module](./exec.md) - For command execution
