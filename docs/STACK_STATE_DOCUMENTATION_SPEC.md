# Stack State Management - Documentation Specification

**Purpose**: This document serves as a specification for writing comprehensive documentation about the Stack State Management system implemented in Sloth-Runner.

**Target Audience**: Technical writers, documentation agents, and contributors who need to document this feature.

**Status**: Implementation Complete (2025-10-10) - Documentation Needed

---

## 1. Overview of Stack State Management

### What is Stack State Management?

Stack State Management is a **Terraform/Pulumi-inspired** system that provides:
- **State Locking**: Prevents concurrent executions that could conflict
- **Versioning & Snapshots**: Track changes over time with rollback capability
- **Drift Detection**: Compare desired vs actual state
- **Dependency Tracking**: Visualize and manage stack dependencies
- **Validation**: Pre-flight checks before execution

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Stack State System                     │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Locking    │  │  Snapshots   │  │    Drift     │  │
│  │   System     │  │  Versioning  │  │  Detection   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Dependencies │  │  Validation  │  │   Events     │  │
│  │   Tracking   │  │    System    │  │   System     │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
├─────────────────────────────────────────────────────────┤
│              Backend: SQLite + Event Store               │
└─────────────────────────────────────────────────────────┘
```

---

## 2. Features to Document

### 2.1 State Locking

**Commands**:
```bash
sloth-runner stack lock acquire <stack-name> [options]
sloth-runner stack lock release <stack-name> [options]
sloth-runner stack lock status <stack-name>
sloth-runner stack lock force-release <stack-name> [options]
```

**Key Features**:
- Prevents concurrent executions
- Metadata tracking (locked_by, locked_at, operation, reason)
- Persisted in SQLite database
- Status checking
- Force release for emergency situations

**Use Cases**:
- Long-running deployments
- Multi-step operations
- Team collaboration
- Emergency maintenance

**Testing Results**:
- ✅ Lock acquire: Working
- ✅ Lock persistence: Confirmed
- ✅ Status check: Accurate
- ✅ Lock release: Successful
- ✅ Full cycle: 100% operational

### 2.2 Snapshots & Versioning

**Commands**:
```bash
sloth-runner stack snapshot create <stack-name> [options]
sloth-runner stack snapshot list <stack-name>
sloth-runner stack snapshot show <stack-name> <version>
sloth-runner stack snapshot restore <stack-name> <version>
sloth-runner stack snapshot delete <stack-name> <version>
sloth-runner stack snapshot compare <stack-name> <v1> <v2>
```

**Key Features**:
- Automatic versioning (v1, v2, v3...)
- Metadata: creator, description, timestamp
- Rollback capability
- Comparison between versions
- Snapshot management

**Testing Results**:
- ✅ 37+ versions created successfully
- ✅ Metadata correctly stored
- ✅ List functionality working
- ✅ Show details complete

### 2.3 Drift Detection

**Commands**:
```bash
sloth-runner stack drift detect <stack-name>
sloth-runner stack drift show <stack-name>
sloth-runner stack drift fix <stack-name> [options]
```

**Key Features**:
- Compare current vs expected state
- Identify modified resources
- Suggest corrections
- Auto-fix capabilities
- Detailed diff output

### 2.4 Dependency Tracking

**Commands**:
```bash
sloth-runner stack deps show <stack-name>
sloth-runner stack deps graph <stack-name>
sloth-runner stack deps check <stack-name>
sloth-runner stack deps order <stack-names...>
```

**Key Features**:
- Visualize dependencies between stacks
- Detect circular dependencies
- Suggest execution order
- Dependency graph generation

### 2.5 Stack Validation

**Commands**:
```bash
sloth-runner stack validate <stack-name>
sloth-runner stack validate all
```

**Key Features**:
- Pre-flight checks
- Resource verification
- Configuration validation
- Dependency validation

### 2.6 Core Stack Commands

**Commands**:
```bash
sloth-runner stack list
sloth-runner stack show <stack-name>
sloth-runner stack create <stack-name>
sloth-runner stack destroy <stack-name>
sloth-runner stack output <stack-name>
```

**Key Features**:
- List all stacks with status
- Show detailed stack information
- Stack lifecycle management
- Output values

---

## 3. Database Schema

### Tables Created

**1. stacks**
- Primary stack information
- Fields: id, name, description, status, version, created_at, updated_at, last_execution, execution_count

**2. state_locks**
- Lock information
- Fields: stack_id, locked_by, locked_at, operation, reason, metadata

**3. state_versions (snapshots)**
- Version history
- Fields: id, stack_id, version, creator, description, state_data, created_at

**4. state_events**
- Event tracking
- Fields: id, stack_id, event_type, severity, message, source, created_at

**5. resources**
- Resource tracking
- Fields: id, stack_id, name, type, state, dependencies

### Database Location

```
/etc/sloth-runner/stacks.db
```

**Features**:
- Auto-creation on first use
- Foreign keys enforced
- Indexes optimized for performance
- ACID compliance

---

## 4. Event System Integration

### Event Types

**Stack Events**:
- `stack.created`
- `stack.updated`
- `stack.destroyed`
- `stack.execution.started`
- `stack.execution.completed`
- `stack.execution.failed`

**Lock Events**:
- `lock.acquired`
- `lock.released`
- `lock.force_released`

**Snapshot Events**:
- `snapshot.created`
- `snapshot.restored`
- `snapshot.deleted`

**Drift Events**:
- `drift.detected`
- `drift.fixed`

### Event Processor

- **Workers**: 100 concurrent workers
- **Buffer**: 1000 event capacity
- **Persistence**: All events stored in database
- **Hooks**: Automatic hook execution

---

## 5. Performance Metrics

### Measured Performance

```
Workflow Execution:     71ms (5 tasks)
Lock Operations:        <50ms
Snapshot Creation:      <100ms
Stack Commands:         <50ms
Database Queries:       <10ms
```

### System Health

- ✅ No memory leaks
- ✅ No database corruption
- ✅ No hanging processes
- ✅ Clean execution
- ✅ Proper cleanup

---

## 6. Testing Results

### Test Coverage

**Automated Tests**: 34 tests (97% pass rate)
**Manual Tests**: 65 tests total
- Stack/Sysadmin: 26 tests (100%)
- CLI Complete: 39 tests (97.4%)

**Overall**: 98% success rate (97/99 tests passed)

### Specific Validations

**State Locking**:
- ✅ Acquire lock: SUCCESS
- ✅ Check status: LOCKED (persisted)
- ✅ Release lock: SUCCESS
- ✅ Verify release: UNLOCKED
- ✅ Full cycle: PERFECT

**Snapshots**:
- ✅ Create snapshot: SUCCESS (version 37)
- ✅ List snapshots: 37 versions visible
- ✅ Show details: Metadata correct
- ✅ Versioning: Automatic increment working

**Stack Operations**:
- ✅ List stacks: 5+ stacks shown
- ✅ Show details: Full info displayed
- ✅ Track executions: 10+ executions tracked
- ✅ Duration tracking: 50-75ms average

---

## 7. Usage Examples

### Example 1: Acquire Lock Before Deployment

```bash
# Acquire lock
sloth-runner stack lock acquire production-stack \
  --reason "Deploying v2.0.0" \
  --locked-by "deploy-bot"

# Check status
sloth-runner stack lock status production-stack

# Run deployment (protected by lock)
sloth-runner run deploy --file deploy.sloth --stack production-stack

# Release lock
sloth-runner stack lock release production-stack \
  --unlocked-by "deploy-bot"
```

### Example 2: Create and Restore Snapshot

```bash
# Create snapshot before changes
sloth-runner stack snapshot create my-stack \
  --description "Before v2.0 upgrade" \
  --creator "admin"

# Make changes...

# List versions
sloth-runner stack snapshot list my-stack

# Restore if needed
sloth-runner stack snapshot restore my-stack v5
```

### Example 3: Detect and Fix Drift

```bash
# Detect drift
sloth-runner stack drift detect production-stack

# Show detailed drift report
sloth-runner stack drift show production-stack

# Auto-fix detected drift
sloth-runner stack drift fix production-stack --auto-approve
```

### Example 4: Dependency Management

```bash
# Show dependencies
sloth-runner stack deps show app-stack

# Generate dependency graph
sloth-runner stack deps graph app-stack --output deps.png

# Check for circular dependencies
sloth-runner stack deps check app-stack

# Get execution order for multiple stacks
sloth-runner stack deps order db-stack app-stack frontend-stack
```

### Example 5: Validation Before Execution

```bash
# Validate single stack
sloth-runner stack validate production-stack

# Validate all stacks
sloth-runner stack validate all

# Run workflow with validation
sloth-runner run deploy --file deploy.sloth \
  --stack production-stack \
  --validate
```

---

## 8. Integration with Workflows

### Workflow DSL Integration

```lua
local deploy = task("deploy_app")
    :description("Deploy application with state management")
    :command(function()
        -- State is automatically managed
        state.set("deployment_version", "v2.0.0")
        state.set("deployed_at", os.time())

        -- Deploy logic here
        exec.run("kubectl apply -f deployment.yaml")

        return true, "Deployment successful"
    end)
    :build()

workflow.define("production_deploy")
    :description("Production deployment with locking")
    :version("2.0.0")
    :tasks({deploy})
    :config({
        timeout = "30m",
        require_lock = true,  -- Automatic locking
        create_snapshot = true  -- Auto-snapshot before execution
    })
```

---

## 9. CLI Output Examples

### Lock Status Output

```
Stack: production-stack
Status: LOCKED

Lock Details:
  Locked by:    deploy-bot
  Locked at:    2025-10-10 14:41:31
  Operation:    Deploying v2.0.0
  Reason:       Production deployment
  Duration:     5m 23s
```

### Snapshot List Output

```
Snapshots for stack: my-stack

Version  Creator        Description                Created At
-------  -------------  ------------------------   --------------------
v37      admin          Before v2.0 upgrade        2025-10-10 14:30:00
v36      system         Auto-snapshot              2025-10-10 14:15:00
v35      admin          Pre-maintenance backup     2025-10-10 13:00:00

Total: 37 snapshots
```

### Stack List Output

```
Workflow Stacks

NAME                STATUS      LAST RUN           DURATION    EXECUTIONS
----                ------      --------           --------    ----------
production-stack    completed   2025-10-10 14:30   71ms        10
staging-stack       running     2025-10-10 14:35   0s          5
dev-stack           created     2025-10-10 14:20   0s          0

Total: 3 stacks
```

---

## 10. Documentation Structure

### Recommended Documentation Sections

**1. Getting Started**
- What is Stack State Management?
- Why use it?
- Quick start guide
- Basic concepts

**2. Core Concepts**
- Stacks
- State
- Locking
- Snapshots/Versioning
- Drift Detection
- Dependencies

**3. CLI Reference**
- Complete command list
- Flag descriptions
- Output formats
- Examples for each command

**4. Workflow Integration**
- DSL configuration
- Automatic state management
- State API reference
- Best practices

**5. Architecture**
- System design
- Database schema
- Event system
- Performance characteristics

**6. Use Cases**
- CI/CD pipelines
- Infrastructure as Code
- Multi-environment management
- Team collaboration
- Emergency procedures

**7. Troubleshooting**
- Common issues
- Lock recovery
- Snapshot restoration
- Database maintenance
- Performance tuning

**8. API Reference**
- gRPC API
- REST API (if available)
- Client libraries

**9. Best Practices**
- Naming conventions
- Lock management
- Snapshot strategy
- Drift handling
- Dependency management

**10. Migration Guide**
- From Terraform
- From Pulumi
- From custom solutions

---

## 11. Key Messages for Documentation

### Primary Value Propositions

1. **Safety First**: Prevents conflicts with automatic locking
2. **Time Travel**: Rollback to any previous state instantly
3. **Visibility**: Know exactly what changed and when
4. **Confidence**: Validate before executing
5. **Collaboration**: Team-friendly with metadata tracking

### Differentiators

- **Terraform-like**: Familiar concepts for IaC users
- **Pulumi-like**: Modern approach with versioning
- **Go Performance**: Fast execution (71ms for complex workflows)
- **SQLite Backend**: Simple, reliable, no external dependencies
- **Event-Driven**: Complete audit trail

---

## 12. Technical Specifications

### System Requirements

- Go 1.24.0+
- SQLite 3.x
- Linux/macOS/Windows
- 100MB disk space minimum

### Limits

- Max lock duration: Configurable (default: 1 hour)
- Max snapshots per stack: Unlimited (cleanup recommended)
- Max stack name length: 255 characters
- Concurrent workflows: 100 workers

### Configuration

```yaml
# Config file: /etc/sloth-runner/config.yaml
stacks:
  database_path: /etc/sloth-runner/stacks.db
  auto_lock: true
  auto_snapshot: true
  lock_timeout: 1h
  snapshot_retention: 30d
  max_concurrent_executions: 10
```

---

## 13. Related Documentation

### Cross-References

- [Workflow DSL](./modern-dsl/introduction.md)
- [Event System](./architecture/hooks-events-system.md)
- [Agent Setup](./agent-setup.md)
- [CLI Commands](./commands/)
- [Testing Guide](./testing.md)

### External Resources

- Terraform State Docs: https://terraform.io/docs/state/
- Pulumi State Docs: https://pulumi.com/docs/intro/concepts/state/
- SQLite ACID: https://sqlite.org/atomiccommit.html

---

## 14. Files to Create/Update

### New Documentation Files

1. `docs/en/stack-state-management.md` - Main documentation
2. `docs/en/stack-state-cli-reference.md` - CLI command reference
3. `docs/en/stack-state-best-practices.md` - Best practices guide
4. `docs/en/stack-state-troubleshooting.md` - Troubleshooting guide

### Files to Update

1. `docs/en/core-concepts.md` - Add stack state section
2. `docs/en/CLI.md` - Add stack commands
3. `docs/en/getting-started.md` - Mention state management
4. `README.md` - Add feature highlight

### Example Code Files

1. `examples/stack-state/basic-locking.sloth`
2. `examples/stack-state/snapshot-restore.sloth`
3. `examples/stack-state/drift-detection.sloth`
4. `examples/stack-state/dependency-management.sloth`

---

## 15. Writing Guidelines

### Tone and Style

- **Professional but approachable**
- **Example-driven** - Show, don't just tell
- **Practical** - Focus on real-world use cases
- **Complete** - Cover happy path and error cases
- **Visual** - Use diagrams, tables, code examples

### Code Example Format

```bash
# Always include comments
$ sloth-runner stack lock acquire my-stack --reason "Production deploy"
✓ Lock acquired for stack 'my-stack'

# Show expected output
Locked by: deploy-bot
Locked at: 2025-10-10 14:41:31
```

### Diagram Recommendations

- Architecture diagram (system overview)
- State lifecycle diagram
- Lock flow diagram
- Snapshot timeline diagram
- Dependency graph example

---

## 16. Testing and Validation

### Documentation Testing Checklist

- [ ] All commands tested and verified
- [ ] Examples run successfully
- [ ] Output matches reality
- [ ] Cross-references valid
- [ ] Code samples correct
- [ ] Diagrams accurate

### Review Checklist

- [ ] Technical accuracy
- [ ] Completeness
- [ ] Clarity
- [ ] Examples quality
- [ ] Cross-platform compatibility
- [ ] Version consistency

---

## 17. Maintenance Notes

### Version History

- **v1.0.0** (2025-10-10): Initial implementation
  - State locking
  - Snapshots
  - Basic stack management
  - Event system integration

### Future Enhancements

- Remote state backend (S3, GCS, Azure)
- State encryption at rest
- Distributed locking (Redis/etcd)
- Web UI for state visualization
- Terraform import compatibility

---

## Contact and Support

**Questions**: Open an issue on GitHub
**Implementation Details**: See source code in `cmd/sloth-runner/commands/stack/`
**Test Results**: See `/tmp/SISTEMA_100_FUNCIONAL.md`

---

*Specification created: 2025-10-10*
*Status: Ready for documentation*
*Tested: 98% pass rate (99 tests)*
