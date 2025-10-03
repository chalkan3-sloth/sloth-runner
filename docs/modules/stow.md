# Stow Module

The `stow` module provides comprehensive GNU Stow integration for managing dotfiles and symbolic link packages. It includes full idempotency support and state management, making it perfect for configuration management and dotfile deployment.

## Overview

GNU Stow is a symlink farm manager which takes distinct packages of software and/or data located in separate directories on the filesystem, and makes them appear to be installed in a single directory tree.

The Sloth Runner `stow` module wraps Stow functionality with:
- **Idempotent operations** - Only makes changes when necessary
- **State tracking** - Tracks what has been stowed
- **Remote execution** - Can delegate to agents via `:delegate_to()`
- **Verification** - Check integrity of stowed packages

## Functions

### stow.stow()

Stow a package by creating symlinks from the stow directory to the target directory.

**Parameters:**
- `package` (string, required): Name of the package to stow
- `dir` or `stow_dir` (string, optional): Stow directory (default: `$HOME/.dotfiles`)
- `target` (string, optional): Target directory (default: `$HOME`)
- `verbose` (boolean, optional): Enable verbose output
- `no_folding` (boolean, optional): Disable directory folding
- `ignore` (table, optional): List of patterns to ignore
- `override` (table, optional): List of patterns to override
- `defer` (table, optional): List of patterns to defer
- `delegate_to` (string, optional): Agent name for remote execution

**Returns:**
- `result` (table): Operation result with fields:
  - `changed` (boolean): Whether changes were made
  - `status` (string): Operation status ("stowed", "already_stowed")
  - `package` (string): Package name
  - `target` (string): Target directory
  - `links` (table): List of created symlinks
- `error` (string or nil): Error message if operation failed

**Example:**
```lua
task({
    name = "setup-nvim-config",
    run = function()
        local result, err = stow.stow({
            package = "nvim",
            dir = "/home/user/.dotfiles",
            target = "/home/user/.config"
        })
        
        if err then
            error("Failed to stow nvim config: " .. err)
        end
        
        if result.changed then
            print("✓ Neovim configuration stowed")
        else
            print("• Neovim configuration already in place")
        end
    end
})
```

### stow.unstow()

Remove stowed symlinks for a package.

**Parameters:** Same as `stow.stow()`

**Returns:** Same structure as `stow.stow()`

**Example:**
```lua
task({
    name = "remove-old-config",
    run = function()
        local result, err = stow.unstow({
            package = "old-vim",
            dir = "/home/user/.dotfiles"
        })
        
        if result.changed then
            print("✓ Old configuration removed")
        end
    end
})
```

### stow.restow()

Restow a package (unstow then stow) - useful for updating symlinks.

**Parameters:** Same as `stow.stow()`

**Returns:** Same structure as `stow.stow()`

**Example:**
```lua
task({
    name = "update-configs",
    run = function()
        stow.restow({
            package = "zsh",
            verbose = true
        })
    end
})
```

### stow.adopt()

Adopt existing files into the stow package. Files in the target directory will be moved to the stow directory and replaced with symlinks.

**Parameters:** Same as `stow.stow()`

**Example:**
```lua
task({
    name = "adopt-existing-configs",
    run = function()
        -- Move existing .bashrc into dotfiles and create symlink
        stow.adopt({
            package = "bash",
            dir = "/home/user/.dotfiles"
        })
    end
})
```

### stow.check()

Check what would happen without actually making changes (dry run).

**Parameters:** Same as `stow.stow()`

**Returns:**
- `result` (table):
  - `package` (string): Package name
  - `output` (string): Detailed output
  - `would_succeed` (boolean): Whether operation would succeed

**Example:**
```lua
task({
    name = "check-before-stow",
    run = function()
        local result, err = stow.check({
            package = "tmux"
        })
        
        if result.would_succeed then
            print("Safe to stow tmux config")
        else
            print("Conflicts detected: " .. result.output)
        end
    end
})
```

### stow.simulate()

Alias for `stow.check()`.

### stow.is_stowed()

Check if a package is currently stowed.

**Parameters:** Same as `stow.stow()`

**Returns:**
- `is_stowed` (boolean): True if package is stowed
- `error` (string or nil): Error message if check failed

**Example:**
```lua
task({
    name = "conditional-stow",
    run = function()
        local is_stowed = stow.is_stowed({
            package = "vim"
        })
        
        if not is_stowed then
            stow.stow({ package = "vim" })
        end
    end
})
```

### stow.get_links()

Get a list of symlinks created by a stowed package.

**Parameters:** Same as `stow.stow()`

**Returns:**
- `links` (table): List of symlink paths
- `error` (string or nil): Error message if operation failed

**Example:**
```lua
task({
    name = "list-config-links",
    run = function()
        local links = stow.get_links({
            package = "nvim"
        })
        
        for i, link in ipairs(links) do
            print("  " .. link)
        end
    end
})
```

### stow.list_packages()

List all available packages in the stow directory.

**Parameters:**
- `dir` or `stow_dir` (string, optional): Stow directory to scan

**Returns:**
- `packages` (table): List of package names
- `error` (string or nil): Error message if operation failed

**Example:**
```lua
task({
    name = "show-all-packages",
    run = function()
        local packages = stow.list_packages({
            dir = "/home/user/.dotfiles"
        })
        
        print("Available dotfile packages:")
        for i, pkg in ipairs(packages) do
            print("  • " .. pkg)
        end
    end
})
```

### stow.verify()

Verify the integrity of a stowed package by checking all symlinks.

**Parameters:** Same as `stow.stow()`

**Returns:**
- `result` (table):
  - `package` (string): Package name
  - `total_files` (number): Total files in package
  - `stowed_links` (number): Number of stowed symlinks
  - `is_complete` (boolean): Whether all files are stowed
  - `broken_links` (table): List of broken symlinks
  - `is_valid` (boolean): Whether all symlinks are valid
- `error` (string or nil): Error message if check failed

**Example:**
```lua
task({
    name = "verify-dotfiles",
    run = function()
        local result = stow.verify({
            package = "zsh"
        })
        
        if result.is_valid and result.is_complete then
            print("✓ All zsh configs properly linked")
        else
            print("⚠ Issues detected:")
            print("  Complete: " .. tostring(result.is_complete))
            print("  Valid: " .. tostring(result.is_valid))
        end
    end
})
```

## Complete Examples

### Basic Dotfiles Setup

```lua
task({
    name = "setup-dotfiles",
    description = "Deploy all dotfiles using Stow",
    run = function()
        local packages = {"vim", "zsh", "tmux", "git"}
        
        for _, pkg in ipairs(packages) do
            local result, err = stow.stow({
                package = pkg,
                dir = "/home/user/.dotfiles"
            })
            
            if err then
                error("Failed to stow " .. pkg .. ": " .. err)
            end
            
            if result.changed then
                print("✓ " .. pkg .. " configured")
            else
                print("• " .. pkg .. " already configured")
            end
        end
    end
})
```

### Remote Dotfiles Deployment

```lua
task({
    name = "deploy-team-configs",
    description = "Deploy standard configs to all dev machines",
    run = function()
        local servers = {"dev-01", "dev-02", "dev-03"}
        local configs = {"vim", "tmux", "git"}
        
        for _, server in ipairs(servers) do
            print("Deploying to " .. server .. "...")
            
            for _, config in ipairs(configs) do
                stow.stow({
                    package = config,
                    dir = "/opt/team-dotfiles",
                    target = "/home/developer",
                    delegate_to = server
                })
            end
        end
    end
})
```

### Advanced Configuration with Verification

```lua
task({
    name = "managed-nvim-setup",
    description = "Complete Neovim configuration with verification",
    run = function()
        -- Check if already configured
        local is_stowed = stow.is_stowed({
            package = "nvim",
            target = os.getenv("HOME") .. "/.config"
        })
        
        if not is_stowed then
            -- Dry run first
            local check_result = stow.check({
                package = "nvim",
                target = os.getenv("HOME") .. "/.config"
            })
            
            if not check_result.would_succeed then
                print("Conflicts detected, adopting existing files...")
                stow.adopt({
                    package = "nvim",
                    target = os.getenv("HOME") .. "/.config"
                })
            end
            
            -- Now stow
            stow.stow({
                package = "nvim",
                target = os.getenv("HOME") .. "/.config",
                ignore = {"*.swp", "*.bak"}
            })
        end
        
        -- Verify integrity
        local verify = stow.verify({
            package = "nvim",
            target = os.getenv("HOME") .. "/.config"
        })
        
        if verify.is_valid and verify.is_complete then
            print("✓ Neovim fully configured and verified")
        else
            error("Configuration verification failed!")
        end
    end
})
```

### Parallel Dotfiles Deployment

```lua
task({
    name = "parallel-dotfile-deployment",
    description = "Deploy dotfiles to multiple servers in parallel",
    run = function()
        local servers = {"web-01", "web-02", "web-03", "web-04"}
        local packages = {"bash", "vim", "git", "tmux"}
        
        -- Deploy to all servers in parallel
        goroutine.map(servers, function(server)
            print("Deploying to " .. server)
            
            for _, pkg in ipairs(packages) do
                local result = stow.stow({
                    package = pkg,
                    dir = "/opt/dotfiles",
                    delegate_to = server
                })
                
                if result.changed then
                    print("[" .. server .. "] ✓ " .. pkg)
                end
            end
            
            -- Verify all packages
            for _, pkg in ipairs(packages) do
                local verify = stow.verify({
                    package = pkg,
                    delegate_to = server
                })
                
                if not verify.is_valid then
                    error("[" .. server .. "] Verification failed for " .. pkg)
                end
            end
            
            print("[" .. server .. "] All dotfiles deployed and verified")
        end)
    end
})
```

### Dotfiles Rotation Strategy

```lua
task({
    name = "rotate-configs",
    description = "Safely rotate configuration versions",
    run = function()
        local old_version = "vim-v1"
        local new_version = "vim-v2"
        
        -- Check current state
        local is_old_stowed = stow.is_stowed({
            package = old_version
        })
        
        if is_old_stowed then
            print("Removing old version...")
            stow.unstow({ package = old_version })
        end
        
        -- Stow new version
        print("Deploying new version...")
        local result = stow.stow({
            package = new_version,
            verbose = true
        })
        
        if not result.changed then
            print("No changes needed - already on new version")
        else
            -- Verify new deployment
            local verify = stow.verify({ package = new_version })
            
            if verify.is_valid and verify.is_complete then
                print("✓ Successfully rotated to " .. new_version)
            else
                -- Rollback
                print("⚠ Verification failed, rolling back...")
                stow.unstow({ package = new_version })
                stow.stow({ package = old_version })
                error("Rollback completed due to verification failure")
            end
        end
    end
})
```

## Idempotency

All `stow` operations are idempotent:

- **stow()**: Will skip if package is already fully stowed
- **unstow()**: Will skip if package is not stowed
- **restow()**: Only restows if necessary

The module tracks state and only makes changes when needed, making it safe to run repeatedly.

## State Management

When a state stack is active, the `stow` module tracks:

- Which packages are stowed
- Where they are stowed
- The hash of the configuration
- Symlink details

This enables:
- Change detection
- Rollback capabilities
- Audit trails

## Requirements

- GNU Stow must be installed on the target system
- Stow package directory must exist and contain valid packages
- Appropriate permissions for creating symlinks in target directory

## Error Handling

All functions return `(result, error)`. Always check for errors:

```lua
local result, err = stow.stow({ package = "vim" })
if err then
    error("Operation failed: " .. err)
end
```

## Best Practices

1. **Always verify** after stowing critical configurations
2. **Use check()** before stowing to detect conflicts
3. **Leverage idempotency** - safe to run multiple times
4. **Version your packages** for easy rollback
5. **Use delegate_to** for remote deployments
6. **Combine with goroutines** for parallel deployment

## See Also

- [File Operations Module](file-ops.md)
- [User Module](user.md)
- [SSH Module](ssh.md)
