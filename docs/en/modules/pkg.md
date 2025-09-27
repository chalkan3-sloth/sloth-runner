# `pkg` Module

The `pkg` module provides functions for managing system packages. It automatically detects the package manager (`apt`, `yum`, `brew`) and uses `sudo` when necessary.

## `pkg.install(packages)`

Installs one or more packages.

*   **`packages`**: A string representing a single package, or a Lua table of strings representing multiple packages to install.

**Returns:**

*   `true` on success, `false` on failure.
*   The command's output (stdout and stderr).

**Examples:**

```lua
-- Install a single package
local success, output = pkg.install("htop")
if not success then
  log.error("Failed to install htop: " .. output)
end

-- Install multiple packages
local success, output = pkg.install({"htop", "git"})
if not success then
  log.error("Failed to install htop and git: " .. output)
end
```

## `pkg.remove(packages)`

Removes one or more packages.

*   **`packages`**: A string representing a single package, or a Lua table of strings representing multiple packages to remove.

**Returns:**

*   `true` on success, `false` on failure.
*   The command's output (stdout and stderr).

**Examples:**

```lua
-- Remove a single package
local success, output = pkg.remove("htop")
if not success then
  log.error("Failed to remove htop: " .. output)
end

-- Remove multiple packages
local success, output = pkg.remove({"htop", "git"})
if not success then
  log.error("Failed to remove htop and git: " .. output)
end
```

## `pkg.update()`

Updates the package list.

**Returns:**

*   `true` on success, `false` on failure.
*   The command's output (stdout and stderr).

**Example:**

```lua
local success, output = pkg.update()
if not success then
  log.error("Failed to update package list: " .. output)
end
```
