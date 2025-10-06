# 🦥 Sloth Management - Saved Workflows

## Overview

The **Sloth Management System** allows you to save `.sloth` workflow files in a SQLite database for easy reuse across your infrastructure. Instead of specifying file paths every time, you can reference workflows by name, making your automation more maintainable and consistent.

## Why Use Saved Sloths?

### Benefits

- **🎯 Centralized Management:** All workflows stored in one database
- **📝 Version Control:** Track usage statistics and last execution times
- **🔄 Reusability:** Reference workflows by name instead of file path
- **🎚️ Control:** Activate/deactivate workflows without deleting them
- **📊 Analytics:** Monitor which workflows are being used and how often
- **🔐 Consistency:** Ensure teams use the same approved workflows

### Use Cases

1. **Standard Operating Procedures:** Store approved deployment workflows
2. **Multi-Environment Deployments:** Reuse the same workflow across dev/staging/prod
3. **Team Collaboration:** Share workflows across team members
4. **Compliance:** Maintain auditable history of workflow executions
5. **CI/CD Integration:** Reference workflows by name in pipelines

## Quick Start

### Save Your First Workflow

```bash
# Create a workflow file
cat > my-deployment.sloth << 'EOF'
workflow({
    name = "deploy_app",
    description = "Deploy application to production",
    tasks = {
        {
            name = "build",
            run = function()
                print("Building application...")
                return {changed = true, message = "Build completed"}
            end
        },
        {
            name = "deploy",
            depends_on = {"build"},
            run = function()
                print("Deploying to production...")
                return {changed = true, message = "Deployment completed"}
            end
        }
    }
})
EOF

# Save it to the sloth database
sloth-runner sloth add prod-deploy \
    --file my-deployment.sloth \
    --description "Production deployment workflow"
```

### Use the Saved Workflow

```bash
# Run using the saved sloth
sloth-runner run deploy --sloth prod-deploy --yes

# The workflow will be loaded from the database automatically
```

## Commands Reference

### `sloth add` - Save a Workflow

Add a new `.sloth` file to the database:

```bash
sloth-runner sloth add <name> --file <path> [options]
```

**Arguments:**
- `<name>`: Unique identifier for the sloth (required)

**Flags:**
- `--file, -f <path>`: Path to the `.sloth` file (required)
- `--description, -d <text>`: Description of what the workflow does
- `--active`: Set as active (default: `true`)

**Examples:**

```bash
# Basic usage
sloth-runner sloth add backup --file /path/to/backup.sloth

# With description
sloth-runner sloth add db-migration \
    --file migration.sloth \
    --description "Database schema migration workflow"

# Save as inactive (won't be runnable until activated)
sloth-runner sloth add experimental-deploy \
    --file deploy-v2.sloth \
    --active=false \
    --description "Testing new deployment process"
```

**Output:**
```
ℹ INFO  Name: prod-deploy
ℹ INFO  File: /path/to/my-deployment.sloth
ℹ INFO  Description: Production deployment workflow
ℹ INFO  Active: true

Adding sloth 'prod-deploy'...
✓ SUCCESS  Sloth 'prod-deploy' added successfully
```

### `sloth list` - View All Workflows

List all saved sloths:

```bash
sloth-runner sloth list [--active]
```

**Flags:**
- `--active, -a`: Show only active sloths

**Examples:**

```bash
# List all sloths
sloth-runner sloth list

# List only active sloths
sloth-runner sloth list --active
```

**Output:**
```
     All Sloths

Name           | Description                   | Active | Usage | Last Used       | Created
prod-deploy    | Production deployment         | ✓      | 42    | 2025-10-06 14:30| 2025-10-01
db-backup      | Daily database backup         | ✓      | 120   | 2025-10-06 03:00| 2025-09-15
test-workflow  | Integration tests             | ✓      | 15    | 2025-10-05 18:20| 2025-10-03
old-deploy     | Legacy deployment (deprecated)|        | 5     | 2025-09-20 10:00| 2025-08-01
```

### `sloth get` - View Workflow Details

Get detailed information about a specific sloth:

```bash
sloth-runner sloth get <name>
```

**Examples:**

```bash
sloth-runner sloth get prod-deploy
```

**Output:**
```
     Sloth: prod-deploy

ℹ INFO  ID: a88f200f-274c-4b9b-8ccf-9eeff1984317
ℹ INFO  Name: prod-deploy
ℹ INFO  Description: Production deployment workflow
ℹ INFO  File Path: /path/to/my-deployment.sloth
ℹ INFO  Active: Yes
ℹ INFO  Created: 2025-10-01 09:00:00
ℹ INFO  Updated: 2025-10-01 09:00:00
ℹ INFO  Last Used: 2025-10-06 14:30:15
ℹ INFO  Usage Count: 42
ℹ INFO  File Hash: d966e74684f5b6a2883b902c041f528b848afc1aca440d35b5d91f59838c9f59
```

### `sloth remove` - Delete a Workflow

Remove a sloth from the database:

```bash
sloth-runner sloth remove <name>
# or
sloth-runner sloth delete <name>
```

**Examples:**

```bash
sloth-runner sloth remove old-deploy
```

**Interactive Confirmation:**
```
Are you sure you want to remove sloth 'old-deploy'? [y/N]: y
✓ SUCCESS  Sloth 'old-deploy' removed successfully
```

### `sloth activate` / `sloth deactivate` - Control Availability

Control whether a sloth can be executed:

```bash
sloth-runner sloth activate <name>
sloth-runner sloth deactivate <name>
```

**Examples:**

```bash
# Deactivate a workflow (makes it non-runnable)
sloth-runner sloth deactivate experimental-deploy

# Activate it later when ready
sloth-runner sloth activate experimental-deploy
```

**Output:**
```
Deactivating sloth 'experimental-deploy'...
⚠ WARNING  Sloth 'experimental-deploy' is now inactive
```

## Integration with `run` Command

The power of saved sloths comes from seamless integration with the `run` command.

### Using `--sloth` Flag

```bash
sloth-runner run <task> --sloth <name> [options]
```

**Important:** The `--sloth` flag **takes precedence** over `--file`. If both are specified, `--file` is ignored.

**Examples:**

```bash
# Run a task using saved sloth
sloth-runner run deploy --sloth prod-deploy --yes

# With delegation to remote agents
sloth-runner run backup --sloth db-backup --delegate-to db-server --yes

# With custom values file
sloth-runner run deploy --sloth prod-deploy --values prod-values.yaml --yes
```

### How It Works

When you use `--sloth`:

1. ✅ Sloth content is retrieved from the database
2. ✅ System checks if the sloth is **active**
3. ✅ Temporary `.sloth` file is created with the content
4. ✅ Workflow is executed using the temporary file
5. ✅ Usage counter is incremented
6. ✅ Last used timestamp is updated
7. ✅ Temporary file is cleaned up

**If sloth is inactive:**
```bash
$ sloth-runner run deploy --sloth old-deploy
ERROR execution failed
└ err: failed to use sloth 'old-deploy': sloth is not active
```

## Workflow Examples

### Example 1: Environment-Specific Deployments

```bash
# Save different environment configurations
sloth-runner sloth add dev-deploy --file deploy.sloth --description "Development deployment"
sloth-runner sloth add staging-deploy --file deploy.sloth --description "Staging deployment"
sloth-runner sloth add prod-deploy --file deploy.sloth --description "Production deployment"

# Use with environment-specific values
sloth-runner run deploy --sloth dev-deploy --values dev-values.yaml --yes
sloth-runner run deploy --sloth staging-deploy --values staging-values.yaml --yes
sloth-runner run deploy --sloth prod-deploy --values prod-values.yaml --yes
```

### Example 2: Team Standard Procedures

```bash
# Save approved SOPs
sloth-runner sloth add incident-response \
    --file sops/incident.sloth \
    --description "Standard incident response procedure"

sloth-runner sloth add security-audit \
    --file sops/audit.sloth \
    --description "Monthly security audit checklist"

# Team members run approved procedures
sloth-runner run respond --sloth incident-response --yes
```

### Example 3: CI/CD Integration

```yaml
# .github/workflows/deploy.yml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Deploy using saved sloth
        run: |
          sloth-runner run deploy \
            --sloth prod-deploy \
            --delegate-to prod-server \
            --yes
```

### Example 4: Workflow Rotation

```bash
# Deactivate old version
sloth-runner sloth deactivate deploy-v1

# Activate new version
sloth-runner sloth activate deploy-v2

# Now all references to 'deploy-v2' will use the new version
sloth-runner run deploy --sloth deploy-v2 --yes
```

## Database and Storage

### Location

Sloths are stored in:
```
~/.sloth-runner/sloths.db
```

### Schema

```sql
CREATE TABLE sloths (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    file_path TEXT NOT NULL,
    content TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at DATETIME,
    usage_count INTEGER DEFAULT 0,
    tags TEXT,
    file_hash TEXT NOT NULL
);

CREATE INDEX idx_sloths_name ON sloths(name);
CREATE INDEX idx_sloths_active ON sloths(is_active);
```

### Backup

To backup your sloths:

```bash
# Backup database
cp ~/.sloth-runner/sloths.db ~/backups/sloths-$(date +%Y%m%d).db

# Restore from backup
cp ~/backups/sloths-20251006.db ~/.sloth-runner/sloths.db
```

## Best Practices

### 1. **Use Descriptive Names**

```bash
# Good
sloth-runner sloth add prod-k8s-deploy --file deploy.sloth

# Avoid
sloth-runner sloth add temp1 --file deploy.sloth
```

### 2. **Add Meaningful Descriptions**

```bash
sloth-runner sloth add db-migration \
    --file migration.sloth \
    --description "PostgreSQL schema migration for v2.0 - includes user table changes"
```

### 3. **Deactivate Instead of Delete**

Preserve history and usage statistics:

```bash
# Don't delete
# sloth-runner sloth remove old-workflow

# Instead, deactivate
sloth-runner sloth deactivate old-workflow
```

### 4. **Use Active-Only Listing**

For operational use, focus on active workflows:

```bash
sloth-runner sloth list --active
```

### 5. **Monitor Usage Statistics**

Regularly check which workflows are being used:

```bash
# Check details of frequently used workflows
sloth-runner sloth get prod-deploy
```

### 6. **Version Your Workflows**

Use naming conventions for versions:

```bash
sloth-runner sloth add deploy-v1 --file deploy-v1.sloth
sloth-runner sloth add deploy-v2 --file deploy-v2.sloth

# When ready to switch
sloth-runner sloth deactivate deploy-v1
sloth-runner sloth activate deploy-v2
```

### 7. **Combine with Delegation**

Use saved sloths with remote execution:

```bash
sloth-runner run setup \
    --sloth server-setup \
    --delegate-to prod-server-01 \
    --delegate-to prod-server-02 \
    --yes
```

## Troubleshooting

### Sloth is Inactive

**Error:**
```
ERROR: failed to use sloth 'my-workflow': sloth is not active
```

**Solution:**
```bash
sloth-runner sloth activate my-workflow
```

### Sloth Not Found

**Error:**
```
ERROR: sloth not found
```

**Solution:**
```bash
# List available sloths
sloth-runner sloth list

# Add the sloth if it doesn't exist
sloth-runner sloth add my-workflow --file workflow.sloth
```

### Duplicate Name

**Error:**
```
ERROR: sloth with this name already exists
```

**Solution:**
```bash
# Option 1: Use different name
sloth-runner sloth add my-workflow-v2 --file workflow.sloth

# Option 2: Remove existing and re-add
sloth-runner sloth remove my-workflow
sloth-runner sloth add my-workflow --file workflow.sloth
```

### File Not Found

**Error:**
```
ERROR: failed to read file: open workflow.sloth: no such file or directory
```

**Solution:**
```bash
# Verify file exists
ls -la workflow.sloth

# Use absolute path
sloth-runner sloth add my-workflow --file /absolute/path/to/workflow.sloth
```

## Architecture

The Sloth Management System follows clean architecture principles:

### Layers

1. **Commands:** CLI interface (`cmd/sloth-runner/commands/sloth/`)
2. **Services:** Business logic (`cmd/sloth-runner/services/sloth_service.go`)
3. **Repository:** Data access (`internal/sloth/sqlite_repository.go`)
4. **Domain:** Models and interfaces (`internal/sloth/sloth.go`)

### Design Patterns

- **Repository Pattern:** Abstracts database operations
- **Service Layer Pattern:** Encapsulates business logic
- **Factory Pattern:** Dependency injection for testability
- **Command Pattern:** CLI command structure

For detailed architecture information, see [Architecture Documentation](../architecture/sloth-runner-architecture.md).

## Performance Considerations

- ✅ **Fast Lookups:** Indexed by name for O(1) retrieval
- ✅ **Minimal Overhead:** Direct database access, no network calls
- ✅ **Efficient Storage:** SQLite with WAL mode
- ✅ **Quick Execution:** Temporary files created in `/tmp`

## Security

- 🔒 Database file permissions: `0600`
- 🔒 Temporary files cleaned up after execution
- 🔒 No sensitive data logged
- 🔒 SHA256 file hashing for integrity verification

## Comparison: File vs Sloth

| Aspect | Using `--file` | Using `--sloth` |
|--------|---------------|----------------|
| **Repeatability** | Specify path every time | Reference by name |
| **Version Control** | Manual file management | Automatic tracking |
| **Usage Analytics** | Not available | Tracked automatically |
| **Team Sharing** | Share file paths | Share sloth names |
| **Active Control** | Not available | Can activate/deactivate |
| **History** | No built-in history | Full execution history |
| **Portability** | Path-dependent | Database-portable |

## Future Enhancements

Planned features for the Sloth Management System:

1. **Update Command:** Update sloth content from file
2. **Tags:** Categorize sloths with tags
3. **Search:** Find sloths by description or tags
4. **Export/Import:** Share sloths between systems
5. **Version History:** Track changes to sloth content
6. **Clone:** Duplicate sloth with new name
7. **Validation:** Syntax check before saving
8. **Web UI:** Browser-based management interface

## Related Documentation

- [Architecture](../architecture/sloth-runner-architecture.md)
- [Run Command](../run-syntax.md)
- [DSL Reference](../modern-dsl/reference-guide.md)
- [Agent Management](../en/master-agent-architecture.md)
- [SSH Management](../ssh-management.md)

## Summary

The Sloth Management System transforms how you organize and execute workflows:

✅ **Centralized** - All workflows in one database
✅ **Trackable** - Usage statistics and history
✅ **Controllable** - Activate/deactivate as needed
✅ **Reusable** - Reference by name, not path
✅ **Scalable** - Works with distributed agents
✅ **Professional** - Enterprise-grade workflow management

Start using saved sloths today to streamline your infrastructure automation! 🚀
