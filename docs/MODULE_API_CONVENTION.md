# Module API Convention - Table-Based Parameters

## Overview

All Sloth Runner modules now use **table-based named parameters** instead of positional parameters. This provides better readability, extensibility, and IDE support.

## Design Principles

### 1. All Parameters Are Named

**Bad (Old Style):**
```lua
pkg.install("vim")
pkg.search("python")
```

**Good (New Style):**
```lua
pkg.install({packages = "vim"})
pkg.search({query = "python"})
```

### 2. Empty Tables for Zero-Parameter Functions

Even functions without parameters require an empty table for consistency:

```lua
pkg.update({})
pkg.upgrade({})
pkg.get_manager({})
```

### 3. Flexible Value Types

Parameters should accept both single values and tables:

```lua
-- Single package
pkg.install({packages = "vim"})

-- Multiple packages
pkg.install({packages = {"vim", "git", "curl"}})
```

### 4. Optional Parameters Have Defaults

```lua
-- With default timeout
ssh.connect({
    host = "server.com",
    user = "admin"
})

-- With custom timeout
ssh.connect({
    host = "server.com",
    user = "admin",
    timeout = 60
})
```

## Implementation Pattern

### Go Implementation

```go
func (m *Module) functionName(L *lua.LState) int {
    // Always expect a table as first parameter
    opts := L.CheckTable(1)
    
    // Extract required parameters
    requiredParam := getTableString(opts, "param_name", "")
    if requiredParam == "" {
        L.Push(lua.LFalse)
        L.Push(lua.LString("param_name parameter is required"))
        return 2
    }
    
    // Extract optional parameters with defaults
    optionalParam := getTableString(opts, "optional", "default")
    
    // Implementation...
    
    return 2 // Return success/failure and result
}
```

### Helper Functions

```go
// Helper to extract string from table
func getTableString(tbl *lua.LTable, key string, def string) string {
    val := tbl.RawGetString(key)
    if val.Type() == lua.LTString {
        return val.String()
    }
    return def
}

// Helper to extract int from table
func getTableInt(tbl *lua.LTable, key string, def int) int {
    val := tbl.RawGetString(key)
    if lv, ok := val.(lua.LNumber); ok {
        return int(lv)
    }
    return def
}

// Helper to extract bool from table
func getTableBool(tbl *lua.LTable, key string, def bool) bool {
    val := tbl.RawGetString(key)
    return lua.LVAsBool(val)
}
```

## Module Examples

### pkg Module (Package Management)

```lua
local pkg = require("pkg")

-- Install packages
pkg.install({packages = "vim"})
pkg.install({packages = {"git", "curl", "wget"}})

-- Search for packages
pkg.search({query = "python3"})

-- Check if installed
pkg.is_installed({package = "nginx"})

-- Get package info
pkg.info({package = "curl"})

-- System maintenance
pkg.update({})
pkg.upgrade({})
pkg.clean({})
pkg.autoremove({})

-- Utilities
pkg.which({executable = "git"})
pkg.version({package = "bash"})
pkg.deps({package = "nginx"})

-- Remove packages
pkg.remove({packages = "old-package"})

-- Install from local file
pkg.install_local({file = "/path/to/package.deb"})

-- Get package manager name
local manager, _ = pkg.get_manager({})
```

### ssh Module (SSH Operations)

```lua
local ssh = require("ssh")

-- Connect to server
ssh.connect({
    host = "server.com",
    user = "admin",
    port = 22,
    key_path = "~/.ssh/id_rsa",
    timeout = 30
})

-- Execute command
ssh.exec({
    command = "ls -la",
    timeout = 60
})

-- Upload file
ssh.upload({
    local_path = "/local/file.txt",
    remote_path = "/remote/file.txt"
})

-- Download file
ssh.download({
    remote_path = "/remote/file.txt",
    local_path = "/local/file.txt"
})

-- Disconnect
ssh.disconnect({})
```

## Benefits

### 1. Self-Documenting Code

```lua
-- Clear what each parameter means
pkg.install({packages = "vim"})

-- vs ambiguous
pkg.install("vim")
```

### 2. Easy to Extend

```lua
-- Can add new parameters without breaking existing code
pkg.install({
    packages = "vim",
    options = {"-y"},
    update_cache = true
})
```

### 3. Parameter Order Independence

```lua
-- These are equivalent
ssh.connect({host = "server", user = "admin", port = 22})
ssh.connect({user = "admin", port = 22, host = "server"})
```

### 4. Better IDE Support

IDEs can provide autocompletion for parameter names:
- Type `pkg.install({` → IDE suggests `packages`
- Type `ssh.connect({` → IDE suggests `host`, `user`, `port`, etc.

### 5. Validation and Error Messages

```lua
-- Clear error message
pkg.install({})
-- Error: "packages parameter is required"

-- vs unclear
pkg.install()
-- Error: "expected string, got nil"
```

## Migration Guide

### For Module Developers

1. Update function signature to accept table:
   ```go
   func (m *Module) install(L *lua.LState) int {
       opts := L.CheckTable(1)
       // ...
   }
   ```

2. Extract parameters using helpers:
   ```go
   packages := getTableString(opts, "packages", "")
   ```

3. Validate required parameters:
   ```go
   if packages == "" {
       L.Push(lua.LFalse)
       L.Push(lua.LString("packages parameter is required"))
       return 2
   }
   ```

4. Update tests:
   ```go
   err := L.DoString(`
       local pkg = require("pkg")
       pkg.install({packages = "vim"})
   `)
   ```

5. Update documentation with examples.

### For Users

Update existing scripts:

**Before:**
```lua
pkg.install("vim")
pkg.search("python")
pkg.is_installed("git")
```

**After:**
```lua
pkg.install({packages = "vim"})
pkg.search({query = "python"})
pkg.is_installed({package = "git"})
```

## Checklist for New Modules

When creating a new module, ensure:

- [ ] All functions accept table as first parameter
- [ ] Required parameters are validated
- [ ] Optional parameters have sensible defaults
- [ ] Helper functions (getTableString, etc.) are used
- [ ] Error messages mention parameter names
- [ ] Documentation shows table syntax
- [ ] Tests use table syntax
- [ ] Functions work with delegate_to()

## See Also

- [pkg Module Documentation](../docs/en/modules/pkg.md)
- [ssh Module Documentation](../docs/en/modules/ssh.md)
- [Modern DSL Guide](../docs/en/modern-dsl/overview.md)
