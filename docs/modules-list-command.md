# Module Documentation Command

The `sloth-runner modules list` command provides comprehensive documentation for all built-in modules available in Sloth Runner.

## Usage

```bash
sloth-runner modules list [flags]
```

## Flags

- `--module, -m`: Show details for a specific module
- `--format, -f`: Output format: `pretty` (default) or `json`

## Examples

### List All Modules

Display all available modules with their functions:

```bash
sloth-runner modules list
```

Output:
```
Sloth Runner - Available Modules

# pkg - Package management for multiple Linux distributions

┌──────────────────┬─────────────────────────────┐
│ Function         │ Description                 │
├──────────────────┼─────────────────────────────┤
│ pkg.install      │ Install packages            │
│ pkg.remove       │ Remove packages             │
│ pkg.update       │ Update package cache        │
│ pkg.upgrade      │ Upgrade all packages        │
│ pkg.is_installed │ Check if package installed  │
└──────────────────┴─────────────────────────────┘

# systemd - Systemd service management
...
```

### View Specific Module Documentation

Get detailed documentation including examples for a specific module:

```bash
sloth-runner modules list --module pkg
```

Output:
```
Sloth Runner - Available Modules

# pkg - Package management for multiple Linux distributions

┌──────────────────┬─────────────────────────────┐
│ Function         │ Description                 │
├──────────────────┼─────────────────────────────┤
│ pkg.install      │ Install packages            │
│ pkg.remove       │ Remove packages             │
│ pkg.update       │ Update package cache        │
│ pkg.upgrade      │ Upgrade all packages        │
│ pkg.is_installed │ Check if package installed  │
└──────────────────┴─────────────────────────────┘

# Examples

INFO pkg.install
  Parameters: {packages = {...}, target = 'agent_name'}
  pkg.install({
      packages = {"nginx", "curl"},
      target = "web-server"
  })

INFO pkg.remove
  Parameters: {packages = {...}, target = 'agent_name'}
  pkg.remove({
      packages = {"apache2"},
      target = "web-server"
  })

INFO pkg.update
  Parameters: {target = 'agent_name'}
  pkg.update({
      target = "web-server"
  })

INFO pkg.upgrade
  Parameters: {target = 'agent_name'}
  pkg.upgrade({
      target = "web-server"
  })

INFO pkg.is_installed
  Parameters: {package = 'name', target = 'agent_name'}
  local installed = pkg.is_installed({
      package = "nginx",
      target = "web-server"
  })
```

### JSON Output

Get machine-readable JSON output for integration with other tools:

```bash
sloth-runner modules list --format json
```

Output:
```json
[
  {
    "Name": "pkg",
    "Description": "Package management for multiple Linux distributions",
    "Functions": [
      {
        "Name": "pkg.install",
        "Description": "Install one or more packages",
        "Example": "pkg.install({\n    packages = {\"nginx\", \"curl\"},\n    target = \"web-server\"\n})",
        "Parameters": "{packages = {...}, target = 'agent_name'}"
      }
    ]
  }
]
```

### Query JSON Output with jq

Count total modules:
```bash
sloth-runner modules list --format json | jq 'length'
# Output: 21
```

List all module names:
```bash
sloth-runner modules list --format json | jq -r '.[].Name'
# Output:
# pkg
# systemd
# user
# ssh
# file
# ...
```

Count functions per module:
```bash
sloth-runner modules list --format json | \
  jq -r '.[] | .Name + ": " + (.Functions | length | tostring) + " functions"'
# Output:
# pkg: 5 functions
# systemd: 8 functions
# user: 8 functions
# ...
```

Get all function names from a specific module:
```bash
sloth-runner modules list --module systemd --format json | \
  jq -r '.[].Functions[].Name'
# Output:
# systemd.enable
# systemd.disable
# systemd.start
# systemd.stop
# ...
```

## Available Modules

Sloth Runner includes **21 built-in modules** with **69+ functions**:

### Core Infrastructure
- **pkg** (5 functions) - Package management for multiple Linux distributions
- **systemd** (8 functions) - Systemd service management
- **user** (8 functions) - Linux user and group management
- **ssh** (4 functions) - SSH key and configuration management
- **file** (9 functions) - File and directory operations

### Development & Utilities
- **http** (2 functions) - HTTP client operations
- **cmd** (1 function) - Execute shell commands
- **json** (2 functions) - JSON encoding/decoding
- **yaml** (2 functions) - YAML encoding/decoding
- **log** (4 functions) - Logging functions
- **crypto** (3 functions) - Cryptographic operations
- **database** (3 functions) - Database operations (PostgreSQL, MySQL, SQLite)

### DevOps & Cloud
- **terraform** (4 functions) - Terraform operations
- **pulumi** (3 functions) - Pulumi operations
- **docker** (2 functions) - Docker operations
- **kubernetes** (2 functions) - Kubernetes operations

### Cloud Providers
- **aws** (2 functions) - AWS operations
- **azure** (1 function) - Azure operations
- **gcp** (1 function) - Google Cloud Platform operations

### Integrations
- **slack** (1 function) - Slack notifications
- **goroutine** (2 functions) - Concurrent execution with goroutines

## Use Cases

### 1. Learning Sloth Runner

When starting with Sloth Runner, use `modules list` to discover available modules:

```bash
# See all available modules
sloth-runner modules list

# Learn about a specific module
sloth-runner modules list --module pkg
```

### 2. Quick Reference

Keep the command handy for quick reference while writing workflows:

```bash
# Check function signature
sloth-runner modules list --module file | grep "file.template"

# See example usage
sloth-runner modules list --module systemd
```

### 3. Integration with Editors

Use JSON output to integrate with text editors or IDEs:

```bash
# Generate autocomplete data
sloth-runner modules list --format json > modules.json
```

### 4. Documentation Generation

Generate custom documentation for your team:

```bash
# Extract all examples
sloth-runner modules list --format json | \
  jq -r '.[] | .Functions[] | "## " + .Name + "\n\n```lua\n" + .Example + "\n```\n"'
```

### 5. CI/CD Integration

Validate module availability in CI pipelines:

```bash
# Check if a module exists
if sloth-runner modules list --module pkg &>/dev/null; then
  echo "pkg module is available"
fi
```

## Tips

1. **Use module filter** for focused documentation:
   ```bash
   sloth-runner modules list --module <name>
   ```

2. **JSON output** is perfect for scripting:
   ```bash
   sloth-runner modules list --format json
   ```

3. **Combine with jq** for powerful queries:
   ```bash
   sloth-runner modules list --format json | jq '.[] | select(.Name == "pkg")'
   ```

4. **Quick search** with grep:
   ```bash
   sloth-runner modules list | grep -A 5 "pkg.install"
   ```

## See Also

- [Module Reference](modules.md) - Detailed module documentation
- [Package Management](pkg.md) - pkg module details
- [Systemd Management](systemd.md) - systemd module details
- [User Management](user.md) - user module details
