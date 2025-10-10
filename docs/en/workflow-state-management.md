# Workflow State Management

**Terraform/Pulumi-like state management for sloth-runner workflows**

## Overview

Sloth-runner now includes advanced state management capabilities similar to Terraform and Pulumi. This system tracks workflow executions, resources, outputs, and provides versioning, drift detection, and rollback capabilities.

## Key Features

### 1. **State Tracking**
- Track every workflow execution with complete metadata
- Store resource states and outputs
- Maintain execution history and status

### 2. **Versioning**
- Automatic versioning of workflow states
- Complete state snapshots for each version
- Easy rollback to any previous version

### 3. **Drift Detection**
- Compare expected vs actual resource state
- Identify configuration drift
- Detailed diff visualization

### 4. **Resource Management**
- Track all resources created/modified by workflows
- Monitor resource lifecycle (create, update, delete)
- Resource dependency tracking

### 5. **State Locking**
- Prevent concurrent modifications
- Automatic lock cleanup
- Lock holder identification

## Architecture

### Database Schema

The state system uses SQLite with the following tables:

```sql
workflow_states        -- Main workflow execution states
workflow_resources     -- Resources managed by workflows
workflow_outputs       -- Workflow output values
state_versions         -- Historical state snapshots
drift_detections       -- Detected drift between states
```

### Data Structures

```go
type WorkflowState struct {
    ID          string
    Name        string
    Version     int
    Status      WorkflowStateStatus  // pending, running, success, failed, rolled_back
    StartedAt   time.Time
    CompletedAt *time.Time
    Duration    int64
    Metadata    map[string]string
    Resources   []Resource
    Outputs     map[string]string
    ErrorMsg    string
    LockedBy    string
}

type Resource struct {
    ID         string
    Type       string
    Name       string
    Action     ResourceAction  // create, update, delete, read, noop
    Status     string
    Attributes map[string]interface{}
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

## CLI Commands

All workflow state management commands are under `sloth-runner state workflow`:

### List Workflows

```bash
# List all workflow states
sloth-runner state workflow list

# Filter by name
sloth-runner state workflow list --name my-workflow

# Filter by status
sloth-runner state workflow list --status success

# JSON output
sloth-runner state workflow list -o json
```

**Output:**
```
════════════════════════════════════════════════
Workflow States
════════════════════════════════════════════════

NAME              VERSION  STATUS   STARTED           DURATION  RESOURCES
----              -------  ------   -------           --------  ---------
deploy-prod       5        success  2025-10-10 14:30  2m30s     15
update-infra      3        success  2025-10-10 12:15  1m45s     8
backup-dbs        1        failed   2025-10-10 10:00  30s       3

✔ Total: 3 workflow(s)
```

### Show Workflow Details

```bash
# Show by ID or name
sloth-runner state workflow show deploy-prod

# JSON output
sloth-runner state workflow show abc123 -o json
```

**Output:**
```
════════════════════════════════════════════════
Workflow: deploy-prod (v5)
════════════════════════════════════════════════

─── Basic Information ───────────────────────────

ID:           abc123def456
Name:         deploy-prod
Version:      5
Status:       success
Started At:   2025-10-10 14:30:15
Completed At: 2025-10-10 14:32:45
Duration:     2m30s

─── Metadata ────────────────────────────────────

environment: production
deployed_by: user@example.com
commit_sha: a1b2c3d4

─── Resources ───────────────────────────────────

TYPE            NAME              ACTION  STATUS
----            ----              ------  ------
docker_container  web-server-1    create  running
docker_container  web-server-2    create  running
nginx_config    prod.conf         update  applied
ssl_cert        example.com       create  active

─── Outputs ─────────────────────────────────────

load_balancer_ip = 192.168.1.100
deployment_url = https://example.com
```

### List Versions

```bash
# Show all versions of a workflow
sloth-runner state workflow versions abc123

# JSON output
sloth-runner state workflow versions abc123 -o json
```

**Output:**
```
════════════════════════════════════════════════
State Versions: abc123
════════════════════════════════════════════════

VERSION  CREATED AT           CREATED BY  DESCRIPTION
-------  ----------           ----------  -----------
5        2025-10-10 14:32:45  system      Version 5
4        2025-10-10 12:15:30  system      Version 4
3        2025-10-09 16:45:00  system      Version 3
2        2025-10-09 10:20:15  system      Version 2
1        2025-10-09 09:00:00  system      Version 1

✔ Total: 5 version(s)

ℹ Use 'sloth-runner state workflow rollback <workflow-id> <version>' to rollback
```

### Rollback to Previous Version

```bash
# Rollback with confirmation
sloth-runner state workflow rollback abc123 3

# Force rollback (skip confirmation)
sloth-runner state workflow rollback abc123 3 --force
```

**Output:**
```
⚠ About to rollback workflow 'deploy-prod' from version 5 to version 3

Do you want to continue? [y/N]: y

⠸ Rolling back workflow...
✔ Rollback completed successfully

✔ Workflow 'deploy-prod' rolled back to version 3 (new version: 6)
```

### Detect Drift

```bash
# Check for drift in workflow resources
sloth-runner state workflow drift abc123

# JSON output
sloth-runner state workflow drift abc123 -o json
```

**Output:**
```
════════════════════════════════════════════════
Drift Detection: abc123
════════════════════════════════════════════════

⚠ 2 resource(s) have drifted from expected state

RESOURCE TYPE      RESOURCE ID        DRIFTED  DETECTED AT
-------------      -----------        -------  -----------
docker_container   web-server-1       YES      2025-10-10 15:00
nginx_config       prod.conf          YES      2025-10-10 15:01

─── Sample Drift Details: web-server-1 ──────────

Expected:
{
  "image": "nginx:1.25",
  "ports": ["80:80", "443:443"],
  "replicas": 2
}

Actual:
{
  "image": "nginx:1.24",
  "ports": ["80:80"],
  "replicas": 1
}

ℹ Use 'sloth-runner run <workflow>' to apply changes and fix drift
```

### List Resources

```bash
# List all resources in a workflow
sloth-runner state workflow resources abc123

# Filter by resource type
sloth-runner state workflow resources abc123 --type docker_container

# JSON output
sloth-runner state workflow resources abc123 -o json
```

**Output:**
```
════════════════════════════════════════════════
Workflow Resources: deploy-prod
════════════════════════════════════════════════

ID          TYPE              NAME            ACTION  STATUS   CREATED
--          ----              ----            ------  ------   -------
a1b2c3d4... docker_container  web-server-1    create  running  2025-10-10 14:30
e5f6g7h8... docker_container  web-server-2    create  running  2025-10-10 14:31
i9j0k1l2... nginx_config      prod.conf       update  applied  2025-10-10 14:32
m3n4o5p6... ssl_cert          example.com     create  active   2025-10-10 14:32

✔ Total: 4 resource(s)

Actions:
  create: 3
  update: 1
```

### Show Outputs

```bash
# Display workflow outputs
sloth-runner state workflow outputs abc123

# JSON output
sloth-runner state workflow outputs abc123 -o json
```

**Output:**
```
════════════════════════════════════════════════
Workflow Outputs: deploy-prod
════════════════════════════════════════════════

KEY                  VALUE
---                  -----
load_balancer_ip     192.168.1.100
deployment_url       https://example.com
database_endpoint    db.example.com:5432
redis_endpoint       redis.example.com:6379

✔ Total: 4 output(s)
```

### Delete Workflow State

```bash
# Delete with confirmation
sloth-runner state workflow delete abc123

# Force delete (skip confirmation)
sloth-runner state workflow delete abc123 --force
```

**Output:**
```
⚠ About to delete workflow state:
  Name:      deploy-prod
  Version:   5
  Status:    success
  Resources: 4

⚠ This action is IRREVERSIBLE and will delete all state, resources, and versions!

Are you absolutely sure? [y/N]: y

⠸ Deleting workflow state...
✔ Workflow state 'deploy-prod' deleted successfully
```

## Integration with Workflows

### Automatic State Tracking (Coming Soon)

When workflow state management is integrated with execution:

```go
// Workflow execution automatically creates state
workflow, err := runner.Execute("deploy.sloth")

// State is created with:
// - Unique workflow ID
// - Execution metadata
// - Resource tracking
// - Output collection

// Access state after execution
state, err := stateManager.GetWorkflowStateByName("deploy")
fmt.Printf("Workflow completed in %s\n", state.Duration)
```

### Manual State Management

For now, you can manually manage state in your workflows:

```lua
-- Example: Create workflow state
local state = require('state')
local workflow_id = state.create_workflow{
    name = "my-deployment",
    metadata = {
        environment = "production",
        user = os.getenv("USER")
    }
}

-- Track resources
state.add_resource(workflow_id, {
    type = "docker_container",
    name = "web-server",
    action = "create",
    status = "running",
    attributes = {
        image = "nginx:latest",
        ports = {"80:80"}
    }
})

-- Set outputs
state.set_output(workflow_id, "container_id", container.id)
state.set_output(workflow_id, "url", "http://localhost")

-- Mark workflow complete
state.complete_workflow(workflow_id)
```

## Use Cases

### 1. Infrastructure as Code

Track infrastructure changes like Terraform:

```bash
# Deploy infrastructure
sloth-runner run infrastructure.sloth

# Check current state
sloth-runner state workflow show infrastructure

# Detect drift from code
sloth-runner state workflow drift infrastructure

# Reapply to fix drift
sloth-runner run infrastructure.sloth
```

### 2. Deployment Management

Track deployments with full history:

```bash
# List all deployments
sloth-runner state workflow list --name deploy-prod

# Show specific deployment
sloth-runner state workflow show deploy-prod-v5

# Rollback deployment
sloth-runner state workflow rollback deploy-prod-v5 4
```

### 3. Audit and Compliance

Complete audit trail of all changes:

```bash
# List all workflow executions
sloth-runner state workflow list

# Get detailed history
sloth-runner state workflow versions <id>

# Export state for compliance
sloth-runner state workflow show <id> -o json > audit-trail.json
```

## Best Practices

### 1. Naming Conventions

Use descriptive names for workflows:
- `deploy-<environment>` - e.g., `deploy-prod`, `deploy-staging`
- `backup-<resource>` - e.g., `backup-database`, `backup-files`
- `update-<component>` - e.g., `update-ssl-certs`, `update-configs`

### 2. Metadata

Always include useful metadata:
```go
metadata := map[string]string{
    "environment": "production",
    "deployed_by": user,
    "commit_sha": gitCommit,
    "ticket": "JIRA-1234",
}
```

### 3. Resource Tracking

Track all meaningful resources:
- Infrastructure components
- Configuration files
- Deployed services
- External resources (S3, databases, etc.)

### 4. Drift Detection

Run drift detection regularly:
```bash
# Add to cron
0 */6 * * * sloth-runner state workflow drift prod-infra
```

### 5. State Cleanup

Clean up old workflow states periodically:
```bash
# Delete failed/old workflows
sloth-runner state workflow delete <old-workflow-id>
```

## Comparison with Terraform/Pulumi

| Feature | Terraform | Pulumi | Sloth-Runner |
|---------|-----------|--------|--------------|
| State Tracking | ✅ | ✅ | ✅ |
| Versioning | ✅ | ✅ | ✅ |
| Drift Detection | ✅ | ✅ | ✅ |
| Rollback | ⚠️ Limited | ✅ | ✅ |
| State Locking | ✅ | ✅ | ✅ |
| Resource Tracking | ✅ | ✅ | ✅ |
| Outputs | ✅ | ✅ | ✅ |
| Backend | Multiple | Service | SQLite |

## Storage Location

State is stored in SQLite database at:
```
~/.sloth-runner/state.db
```

Database includes:
- Workflow states
- Resources
- Outputs
- Versions
- Drift detections
- State locks

## FAQ

### How is this different from regular state management?

Regular state is key-value based for idempotency. Workflow state is execution-based with versioning, tracking complete workflow runs like Terraform tracks infrastructure changes.

### Can I use both systems together?

Yes! Regular state for idempotency, workflow state for execution tracking. They complement each other.

### How long is state kept?

Forever, unless manually deleted. Implement your own retention policy if needed.

### Can I migrate state between environments?

Yes, export as JSON and import in another environment:
```bash
sloth-runner state workflow show prod -o json > prod-state.json
# Import in another system (feature coming soon)
```

### Is state locked during workflows?

Yes, workflows can acquire locks to prevent concurrent modifications:
```go
sm.WithLock("deploy-prod", "user-123", 5*time.Minute, func() error {
    // Safe to modify state here
})
```

## Next Steps

- [Examples](./workflow-state-examples.md) - Practical examples
- [API Reference](./workflow-state-api.md) - Programming API
- [Integration Guide](./workflow-state-integration.md) - Integrate with your workflows

## Related Commands

- `sloth-runner state list` - List regular state (key-value)
- `sloth-runner state show` - Show key-value state
- `sloth-runner state workflow` - Workflow state management
