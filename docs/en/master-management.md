# 🎯 Master Server Management

Sloth Runner provides a powerful master server management system that allows you to configure and manage multiple master servers, making it easy to work across different environments (production, staging, development, etc.) without constantly specifying connection details.

## 📋 Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Usage Examples](#usage-examples)
- [Best Practices](#best-practices)

## Overview

### What is a Master Server?

A master server in Sloth Runner is a central coordination point that:
- Maintains the agent registry
- Routes tasks to appropriate agents
- Provides a unified API for agent management
- Handles agent health monitoring and metrics

### Why Manage Multiple Masters?

In real-world scenarios, you typically have multiple environments:
- **Production**: Your live infrastructure
- **Staging**: Pre-production testing environment
- **Development**: Local or development infrastructure
- **Testing**: Isolated testing environments

The master management system allows you to:
- Register multiple master servers with friendly names
- Quickly switch between environments
- Avoid repeatedly typing IP addresses and ports
- Prevent mistakes from using the wrong environment

## Quick Start

### 1. Add Your First Master

```bash
# Add production master
sloth-runner master add production 192.168.1.29:50053 \
  --description "Production master server"

# Add staging master
sloth-runner master add staging 10.0.0.5:50053 \
  --description "Staging environment"

# Add local development master
sloth-runner master add local localhost:50053 \
  --description "Local development"
```

### 2. List All Masters

```bash
sloth-runner master list
```

Output:
```
# Registered Master Servers

⭐ production (default)
   Address: 192.168.1.29:50053
   Description: Production master server
   Created: 2025-10-07 16:28:57

staging
   Address: 10.0.0.5:50053
   Description: Staging environment
   Created: 2025-10-07 16:30:15

local
   Address: localhost:50053
   Description: Local development
   Created: 2025-10-07 16:30:42

Total: 3 master(s)
```

### 3. Select Default Master

```bash
# Switch to staging
sloth-runner master select staging

# All subsequent commands will use staging master
sloth-runner agent list
sloth-runner agent shell my-agent
```

### 4. Override Default Master

```bash
# Use a different master for a specific command
sloth-runner agent list --master production

# Or use a direct address
sloth-runner agent list --master 192.168.1.29:50053
```

## Command Reference

### `master add`

Add a new master server configuration.

**Syntax:**
```bash
sloth-runner master add <name> <address> [flags]
```

**Arguments:**
- `<name>`: Unique name for the master (cannot contain `:`)
- `<address>`: Master server address in format `HOST:PORT`

**Flags:**
- `-d, --description`: Description of the master server

**Examples:**
```bash
# Basic usage
sloth-runner master add production 192.168.1.29:50053

# With description
sloth-runner master add staging 10.0.0.5:50053 \
  --description "Staging environment"
```

**Notes:**
- The first master added automatically becomes the default
- Master names must be unique
- Attempting to add a duplicate name will result in an error

---

### `master list`

List all registered master servers.

**Syntax:**
```bash
sloth-runner master list
```

**Output includes:**
- Master name (with ⭐ indicating default)
- Address
- Description (if provided)
- Creation timestamp

---

### `master select`

Set a master server as the default.

**Syntax:**
```bash
sloth-runner master select <name>
```

**Arguments:**
- `<name>`: Name of the master to set as default

**Examples:**
```bash
# Switch to production
sloth-runner master select production

# Switch to local development
sloth-runner master select local
```

**Notes:**
- The default master is used by all commands unless overridden with `--master` flag
- Only one master can be default at a time

---

### `master show`

Show details of a master server.

**Syntax:**
```bash
sloth-runner master show [name]
```

**Arguments:**
- `[name]`: (Optional) Name of the master to show. If omitted, shows the default master.

**Examples:**
```bash
# Show default master
sloth-runner master show

# Show specific master
sloth-runner master show production
```

---

### `master update`

Update a master server's address or description.

**Syntax:**
```bash
sloth-runner master update <name> <new_address> [flags]
```

**Arguments:**
- `<name>`: Name of the master to update
- `<new_address>`: New address in format `HOST:PORT`

**Flags:**
- `-d, --description`: Update description

**Examples:**
```bash
# Update address
sloth-runner master update production 192.168.1.30:50053

# Update address and description
sloth-runner master update staging 10.0.0.6:50053 \
  --description "New staging server"
```

---

### `master remove`

Remove a master server configuration.

**Syntax:**
```bash
sloth-runner master remove <name>
```

**Aliases:** `rm`, `delete`

**Arguments:**
- `<name>`: Name of the master to remove

**Examples:**
```bash
sloth-runner master remove old-production
sloth-runner master rm staging-old
```

**Notes:**
- Cannot remove the default master if other masters exist
- You must select a different master as default first
- Requires confirmation before deletion

---

### `master start`

Start the master gRPC server (unchanged from previous behavior).

**Syntax:**
```bash
sloth-runner master start [flags]
```

**Flags:**
- `-p, --port`: Port for the master gRPC server (default: 50053)
- `--bind`: Address to bind the master server (default: 0.0.0.0)
- `--daemon`: Run master server as daemon

---

## Usage Examples

### Multi-Environment Workflow

```bash
# Setup
sloth-runner master add production 192.168.1.29:50053
sloth-runner master add staging 10.0.0.5:50053
sloth-runner master add local localhost:50053

# Work with production
sloth-runner master select production
sloth-runner agent list
sloth-runner run deploy --delegate-to prod-server

# Switch to staging for testing
sloth-runner master select staging
sloth-runner agent list
sloth-runner run test --delegate-to staging-server

# Quick check on production without switching default
sloth-runner agent list --master production
```

### Using Names vs. Addresses

```bash
# By name (recommended)
sloth-runner agent list --master production

# By direct address (also works)
sloth-runner agent list --master 192.168.1.29:50053

# Default (uses selected master)
sloth-runner agent list
```

### Managing Master Lifecycle

```bash
# Add new master
sloth-runner master add production 192.168.1.29:50053

# View details
sloth-runner master show production

# Update if IP changes
sloth-runner master update production 192.168.1.30:50053

# Remove when no longer needed
sloth-runner master remove production
```

## Best Practices

### 1. Use Descriptive Names

```bash
# Good
sloth-runner master add prod-us-east 192.168.1.29:50053
sloth-runner master add prod-eu-west 10.0.0.5:50053

# Avoid
sloth-runner master add m1 192.168.1.29:50053
sloth-runner master add server2 10.0.0.5:50053
```

### 2. Add Descriptions

```bash
sloth-runner master add production 192.168.1.29:50053 \
  --description "Production US-East datacenter - handles web services"
```

### 3. Set Logical Defaults

- Set your most frequently used environment as default
- For development, use `local` as default
- For operations teams, use `production` as default

### 4. Document Your Masters

Maintain a team wiki or document listing:
- Master names and their purposes
- Network access requirements
- Contact information for each environment

### 5. Security Considerations

- Master server addresses are stored in SQLite database at `/etc/sloth-runner/masters.db` (or `~/.sloth-runner/masters.db`)
- Database file permissions: 0644 (readable by all, writable by owner)
- No sensitive credentials are stored (only addresses and metadata)
- Use firewall rules to restrict master server access

## Resolution Priority

When determining which master to use, Sloth Runner follows this priority order:

1. **`--master` flag** (if provided)
   - If value contains `:`, treated as direct address
   - Otherwise, looked up as a name in the database

2. **`SLOTH_RUNNER_MASTER_ADDR` environment variable** (if set)

3. **Default master from database** (if configured via `master select`)

4. **Legacy `master.conf` file** (for backward compatibility)

5. **Fallback to `localhost:50051`**

## Database Location

Master configurations are stored in:
- **Root/System**: `/etc/sloth-runner/masters.db`
- **User**: `~/.sloth-runner/masters.db`
- **Custom**: Set via `SLOTH_RUNNER_DATA_DIR` environment variable

## Migration from Old System

If you were using the old `master.conf` file or `SLOTH_RUNNER_MASTER_ADDR`:

```bash
# Old way
export SLOTH_RUNNER_MASTER_ADDR=192.168.1.29:50053
sloth-runner agent list

# New way
sloth-runner master add production 192.168.1.29:50053
sloth-runner master select production
sloth-runner agent list  # No export needed!
```

The old methods still work for backward compatibility, but the new master management system is recommended for better organization.

## Troubleshooting

### "Master not found" Error

```bash
$ sloth-runner agent list --master prod
Error: master 'prod' not found

# Solution: List available masters
$ sloth-runner master list

# Or add the master
$ sloth-runner master add prod 192.168.1.29:50053
```

### "Cannot delete default master" Error

```bash
$ sloth-runner master remove production
Error: cannot delete default master 'production'. Select a different default first

# Solution: Select another master first
$ sloth-runner master select staging
$ sloth-runner master remove production
```

### Connection Issues

```bash
# Check master details
$ sloth-runner master show production

# Verify address is correct
$ ping 192.168.1.29

# Test connection directly
$ sloth-runner agent list --master 192.168.1.29:50053
```

## See Also

- [Agent Management](./agent-management.md)
- [Distributed Execution](./distributed.md)
- [Getting Started](./getting-started.md)
