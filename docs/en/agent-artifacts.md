# Agent Artifacts Management

## Overview

**Artifacts** are files produced by tasks during workflow execution that need to be preserved, shared between tasks, or downloaded for inspection. The Agent Artifacts system provides a complete solution for managing these files across distributed agents.

Think of artifacts as the "build outputs" of your workflows - compiled binaries, test reports, logs, configuration files, or any other files that tasks generate and subsequent tasks or humans need to consume.

## Key Features

- **Distributed Storage**: Artifacts stored on remote agents
- **Metadata Tracking**: Associate artifacts with stacks and tasks
- **Streaming Transfer**: Efficient handling of large files
- **Lifecycle Management**: Automatic cleanup of old artifacts
- **Filtering & Search**: Find artifacts by stack, task, or age
- **Checksum Verification**: Ensure data integrity

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Artifacts System                      │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Upload     │  │   Download   │  │    List      │  │
│  │   (Stream)   │  │   (Stream)   │  │  (Metadata)  │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Delete     │  │   Cleanup    │  │    Show      │  │
│  │  (Remove)    │  │   (Policy)   │  │  (Details)   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
├─────────────────────────────────────────────────────────┤
│           Agent Storage: /var/lib/sloth-runner/         │
│              Metadata: SQLite + Checksums                │
└─────────────────────────────────────────────────────────┘
```

---

## CLI Commands

### 1. List Artifacts

List all artifacts stored on an agent with optional filtering.

```bash
sloth-runner agent artifacts list <agent-name> [flags]
```

**Flags**:
- `--stack, -s <name>` - Filter by stack name
- `--task, -t <name>` - Filter by task name
- `--limit, -l <number>` - Maximum artifacts to show (default: 50)

**Examples**:

```bash
# List all artifacts on an agent
sloth-runner agent artifacts list build-agent

# Filter by stack
sloth-runner agent artifacts list build-agent --stack production

# Filter by task and limit results
sloth-runner agent artifacts list build-agent --task build --limit 10
```

**Output**:
```
📦 Artifacts on agent: build-agent 📦

Name                                     Stack           Size         Created              Task
---------------------------------------- --------------- ------------ -------------------- --------------------
app-v1.2.3.bin                          production      2.5 MB       2025-10-10 14:30:00  build
test-results.xml                        production      450 KB       2025-10-10 14:31:00  test
deployment-logs.tar.gz                  production      5.3 MB       2025-10-10 14:32:00  deploy

ℹ Total: 3 artifacts
```

---

### 2. Download Artifacts

Download an artifact from a remote agent to your local machine.

```bash
sloth-runner agent artifacts download <agent-name> <artifact-name> [flags]
```

**Flags**:
- `--output, -o <path>` - Local output path (default: artifact name)

**Examples**:

```bash
# Download to current directory
sloth-runner agent artifacts download build-agent app.bin

# Specify output path
sloth-runner agent artifacts download build-agent app.bin --output ./dist/app

# Download to specific location
sloth-runner agent artifacts download build-agent report.pdf -o /tmp/reports/latest.pdf
```

**Output**:
```
ℹ Downloading artifact 'app.bin' from agent 'build-agent'...
✓ Downloaded app.bin to 'app.bin' (2.5 MB)
```

---

### 3. Upload Artifacts

Upload a local file as an artifact to a remote agent.

```bash
sloth-runner agent artifacts upload <agent-name> <file-path> [flags]
```

**Flags**:
- `--stack, -s <name>` - Associate with stack
- `--task, -t <name>` - Associate with task
- `--name, -n <name>` - Artifact name (default: filename)

**Examples**:

```bash
# Simple upload
sloth-runner agent artifacts upload build-agent ./app.bin

# Associate with stack and task
sloth-runner agent artifacts upload build-agent ./app.bin \
  --stack production \
  --task build

# Upload with custom name
sloth-runner agent artifacts upload build-agent ./binary \
  --name app-v2.0.bin \
  --stack production
```

**Output**:
```
ℹ Uploading './app.bin' to agent 'build-agent' as 'app.bin'...
✓ Uploaded app.bin (2.5 MB)
```

---

### 4. Show Artifact Details

Display detailed information about a specific artifact.

```bash
sloth-runner agent artifacts show <agent-name> <artifact-name>
```

**Example**:

```bash
sloth-runner agent artifacts show build-agent app.bin
```

**Output**:
```
📦 Artifact Details

Name:       app.bin
Path:       /var/lib/sloth-runner/artifacts/production/build/app.bin
Size:       2.5 MB
Checksum:   sha256:a1b2c3d4e5f6...
Stack:      production
Task:       build
Created:    2025-10-10 14:30:00
Modified:   2025-10-10 14:30:00
Downloads:  5
```

---

### 5. Delete Artifacts

Remove an artifact from an agent's storage.

```bash
sloth-runner agent artifacts delete <agent-name> <artifact-name> [flags]
```

**Flags**:
- `--force, -f` - Skip confirmation prompt

**Examples**:

```bash
# With confirmation
sloth-runner agent artifacts delete build-agent old-app.bin

# Force delete without confirmation
sloth-runner agent artifacts delete build-agent old-app.bin --force
```

**Output**:
```
Are you sure you want to delete artifact 'old-app.bin' from agent 'build-agent'? (y/N): y
✓ Deleted artifact 'old-app.bin' from agent 'build-agent'
```

---

### 6. Cleanup Old Artifacts

Automatically remove artifacts older than a specified duration.

```bash
sloth-runner agent artifacts cleanup <agent-name> [flags]
```

**Flags**:
- `--older-than <duration>` - Remove artifacts older than this (default: 30d)
- `--stack, -s <name>` - Limit cleanup to specific stack
- `--dry-run` - Preview what would be deleted without actually deleting

**Duration Format**:
- `7d` - 7 days
- `30d` - 30 days
- `6h` - 6 hours
- `1h` - 1 hour

**Examples**:

```bash
# Default cleanup (30 days)
sloth-runner agent artifacts cleanup build-agent

# Cleanup artifacts older than 7 days
sloth-runner agent artifacts cleanup build-agent --older-than 7d

# Cleanup specific stack only
sloth-runner agent artifacts cleanup build-agent --stack old-project

# Preview without deleting (dry-run)
sloth-runner agent artifacts cleanup build-agent --older-than 7d --dry-run
```

**Output**:
```
✓ Deleted 15 artifacts, freed 45.2 MB
```

---

## Workflow Integration

Artifacts are automatically managed when declared in task definitions using the `:artifacts()` and `:consumes()` methods.

### Producing Artifacts

Use `:artifacts()` to declare files that should be saved after task execution:

```lua
local build_task = task("build")
    :description("Build application binary")
    :command(function()
        -- Build the application
        exec.run("go build -o app.bin ./cmd/app")

        log.info("Build completed successfully")
        return true, "Binary created: app.bin"
    end)
    :artifacts({"app.bin"})  -- Declare artifact
    :build()
```

### Consuming Artifacts

Use `:consumes()` to access artifacts from previous tasks:

```lua
local test_task = task("test")
    :description("Run tests with built binary")
    :depends_on({"build"})
    :consumes({"app.bin"})  -- Consume artifact from build task
    :command(function()
        -- The artifact is automatically copied to this task's workdir
        exec.run("chmod +x app.bin")
        exec.run("./app.bin --test")

        return true, "Tests passed"
    end)
    :build()
```

### Complete CI/CD Example

```lua
-- Build stage
local build = task("build")
    :command(function()
        exec.run("go build -o app.bin")
        return true
    end)
    :artifacts({"app.bin"})
    :build()

-- Test stage
local test = task("test")
    :depends_on({"build"})
    :consumes({"app.bin"})
    :command(function()
        exec.run("./app.bin --test")
        -- Generate test report
        exec.run("./generate-report.sh > test-report.xml")
        return true
    end)
    :artifacts({"test-report.xml"})
    :build()

-- Deploy stage
local deploy = task("deploy")
    :depends_on({"test"})
    :consumes({"app.bin"})
    :command(function()
        exec.run("scp app.bin production:/opt/app/")
        return true
    end)
    :build()

-- Define workflow
workflow.define("ci_pipeline")
    :description("Complete CI/CD pipeline with artifacts")
    :version("1.0.0")
    :tasks({build, test, deploy})
    :config({
        timeout = "30m",
        create_workdir_before_run = true
    })
```

---

## Use Cases

### 1. CI/CD Pipeline

**Scenario**: Build once, deploy everywhere

```bash
# 1. Run build on build agent
sloth-runner run ci_build --file build.sloth --agent build-agent

# 2. Download artifact
sloth-runner agent artifacts download build-agent app-v1.2.3.bin

# 3. Upload to deployment agent
sloth-runner agent artifacts upload deploy-agent app-v1.2.3.bin \
  --stack production \
  --task deploy

# 4. Run deployment
sloth-runner run deploy --file deploy.sloth --agent deploy-agent
```

### 2. Debugging Failed Workflows

**Scenario**: Investigate what went wrong

```bash
# List artifacts from failed workflow
sloth-runner agent artifacts list prod-agent --task failed-task

# Download error logs
sloth-runner agent artifacts download prod-agent error.log

# Inspect locally
cat error.log

# Download core dump if available
sloth-runner agent artifacts download prod-agent core.dump
```

### 3. Artifact Retention Policy

**Scenario**: Keep storage under control

```bash
# Weekly cleanup of old artifacts
sloth-runner agent artifacts cleanup build-agent --older-than 30d

# Cleanup specific old projects
sloth-runner agent artifacts cleanup build-agent \
  --stack legacy-project \
  --older-than 7d

# Preview cleanup before executing
sloth-runner agent artifacts cleanup build-agent \
  --older-than 30d \
  --dry-run
```

### 4. Multi-Stage Builds

**Scenario**: Compile dependencies once, use in multiple builds

```lua
-- Stage 1: Build dependencies
local deps = task("build_deps")
    :command(function()
        exec.run("go mod download")
        exec.run("tar -czf deps.tar.gz $GOPATH/pkg/mod")
        return true
    end)
    :artifacts({"deps.tar.gz"})
    :build()

-- Stage 2: Build app (uses cached deps)
local build = task("build_app")
    :depends_on({"build_deps"})
    :consumes({"deps.tar.gz"})
    :command(function()
        exec.run("tar -xzf deps.tar.gz -C $GOPATH/pkg")
        exec.run("go build -o app.bin")
        return true
    end)
    :artifacts({"app.bin"})
    :build()
```

### 5. Artifact Transfer Between Agents

**Scenario**: Move artifacts across infrastructure

```bash
# Download from source agent
sloth-runner agent artifacts download source-agent app.bin \
  --output /tmp/app.bin

# Upload to target agent
sloth-runner agent artifacts upload target-agent /tmp/app.bin \
  --stack production \
  --task deploy

# Verify transfer
sloth-runner agent artifacts show target-agent app.bin
```

---

## Storage and Organization

### Directory Structure

Artifacts are organized by stack and task on the agent:

```
/var/lib/sloth-runner/artifacts/
├── production/
│   ├── build/
│   │   ├── app-v1.2.3.bin
│   │   └── app-v1.2.2.bin
│   ├── test/
│   │   └── test-results.xml
│   └── deploy/
│       └── deployment.log
├── staging/
│   └── build/
│       └── app-staging.bin
└── metadata.db (SQLite)
```

### Metadata Storage

Artifact metadata is stored in SQLite:

```sql
CREATE TABLE artifacts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    size INTEGER NOT NULL,
    checksum TEXT NOT NULL,
    stack TEXT,
    task TEXT,
    created_at INTEGER NOT NULL,
    modified_at INTEGER NOT NULL,
    download_count INTEGER DEFAULT 0
);

CREATE INDEX idx_artifacts_stack ON artifacts(stack);
CREATE INDEX idx_artifacts_task ON artifacts(task);
CREATE INDEX idx_artifacts_created ON artifacts(created_at);
```

---

## Best Practices

### 1. Naming Conventions

Use descriptive, versioned names:

```bash
# Good
app-v1.2.3.bin
report-2025-10-10.pdf
logs-production-20251010.tar.gz

# Avoid
output.txt
file.bin
temp.log
```

### 2. Always Associate Metadata

Link artifacts to stacks and tasks:

```bash
sloth-runner agent artifacts upload build-agent app.bin \
  --stack production \
  --task build-v1.2.3
```

### 3. Regular Cleanup

Schedule automatic cleanup to prevent storage issues:

```bash
# Cron job example (weekly cleanup)
0 2 * * 0 sloth-runner agent artifacts cleanup build-agent --older-than 30d
```

### 4. Size Awareness

Monitor artifact sizes:

```bash
# List artifacts sorted by size
sloth-runner agent artifacts list build-agent --limit 100 | \
  sort -k3 -hr | \
  head -10
```

### 5. Checksums for Critical Files

Always verify checksums for important artifacts:

```bash
# Show artifact details including checksum
sloth-runner agent artifacts show build-agent app.bin

# Verify after download
sha256sum app.bin
```

### 6. Artifact Lifecycle in Workflows

```lua
workflow.define("production_pipeline")
    :config({
        artifact_retention = "30d",  -- Auto-cleanup
        artifact_compression = true,  -- Compress large files
        verify_checksums = true       -- Always verify
    })
```

---

## Troubleshooting

### Artifact Not Found

```bash
# Check if artifact exists
sloth-runner agent artifacts list build-agent | grep artifact-name

# Show full details
sloth-runner agent artifacts show build-agent artifact-name
```

### Download Failures

```bash
# Check agent connectivity
sloth-runner agent get build-agent

# Retry with longer timeout
timeout 600 sloth-runner agent artifacts download build-agent large-file.bin
```

### Storage Full

```bash
# Aggressive cleanup
sloth-runner agent artifacts cleanup build-agent --older-than 1d

# List largest artifacts
sloth-runner agent artifacts list build-agent --limit 100 | \
  sort -k3 -hr | \
  head -20
```

### Checksum Mismatch

If you encounter checksum errors:

1. Download again (may be network corruption)
2. Check agent disk health
3. Verify source file integrity

### Slow Transfers

For large files:
- Use compression
- Check network bandwidth
- Consider direct agent-to-agent transfer
- Split large artifacts into chunks

---

## Performance Considerations

### Transfer Optimization

- **Streaming**: Files are transferred in 64KB chunks
- **Compression**: Gzip compression for large files
- **Parallel**: Multiple artifacts can be downloaded concurrently

### Storage Efficiency

- **Deduplication**: Same checksum = same file (future feature)
- **Cleanup**: Regular cleanup prevents storage bloat
- **Compression**: Store compressed when appropriate

### Network Usage

```bash
# Estimate bandwidth
# 100MB artifact = ~10s on 100Mbps connection
# 1GB artifact = ~100s on 100Mbps connection
```

---

## Security Considerations

### Access Control

- Only authenticated users can access artifacts
- Agent-level permissions apply
- Consider network segmentation for sensitive artifacts

### Encryption

- In-transit: TLS for all transfers
- At-rest: Use encrypted storage volumes (future feature)

### Audit Trail

All artifact operations are logged:

```bash
# View artifact access logs
sloth-runner events list --type artifact.*
```

---

## Advanced Topics

### Artifact Retention Policies

Implement custom retention policies:

```bash
#!/bin/bash
# cleanup-policy.sh

# Keep last 5 versions of production artifacts
sloth-runner agent artifacts list production-agent --stack production | \
  tail -n +6 | \
  awk '{print $1}' | \
  xargs -I {} sloth-runner agent artifacts delete production-agent {} --force

# Remove all artifacts from old stacks
for stack in $(sloth-runner stack list --status=destroyed --format=json | jq -r '.[].name'); do
    sloth-runner agent artifacts cleanup production-agent --stack $stack --older-than 0d
done
```

### Monitoring

Track artifact usage:

```bash
# Monitor storage
sloth-runner agent artifacts list build-agent | \
  awk 'NR>2 {sum+=$3} END {print "Total: " sum " MB"}'

# Alert on excessive artifacts
TOTAL=$(sloth-runner agent artifacts list build-agent | grep "Total:" | awk '{print $2}')
if [ $TOTAL -gt 1000 ]; then
    echo "WARNING: Too many artifacts ($TOTAL)"
    # Send alert
fi
```

### Bulk Operations

```bash
# Batch download
cat artifacts.txt | while read artifact; do
    sloth-runner agent artifacts download build-agent $artifact
done

# Batch upload
find ./dist -name "*.bin" -exec \
    sloth-runner agent artifacts upload build-agent {} --stack production \;
```

---

## Integration with Other Systems

### Jenkins Pipeline

```groovy
pipeline {
    stages {
        stage('Build') {
            steps {
                sh 'sloth-runner run build --file build.sloth --agent build-agent'
            }
        }
        stage('Retrieve Artifact') {
            steps {
                sh 'sloth-runner agent artifacts download build-agent app.bin'
                archiveArtifacts 'app.bin'
            }
        }
    }
}
```

### GitLab CI

```yaml
build:
  script:
    - sloth-runner run build --file build.sloth --agent build-agent
    - sloth-runner agent artifacts download build-agent app.bin
  artifacts:
    paths:
      - app.bin
```

### GitHub Actions

```yaml
- name: Download Artifact
  run: |
    sloth-runner agent artifacts download build-agent app.bin

- name: Upload to GitHub
  uses: actions/upload-artifact@v2
  with:
    name: app-binary
    path: app.bin
```

---

## Future Enhancements

Planned features for the artifact system:

- **Deduplication**: Store identical files only once
- **Compression**: Automatic compression for large files
- **Encryption**: At-rest encryption
- **S3 Backend**: Store artifacts in S3/GCS/Azure
- **Versioning**: Automatic version tracking
- **Signing**: Cryptographic signing of artifacts
- **Mirroring**: Multi-region artifact replication

---

## Related Documentation

- [Workflow DSL - Artifacts](./core-concepts.md#artifact-management)
- [Agent Commands](./CLI.md#agent-commands)
- [Event System](./advanced-features.md#events)
- [Storage Configuration](./getting-started.md#configuration)

---

## FAQ

**Q: What's the maximum artifact size?**
A: No hard limit, but streaming is used for efficient transfer. 10GB+ files are supported.

**Q: Are artifacts versioned automatically?**
A: Not yet - use naming conventions (app-v1.0.0.bin) until auto-versioning is implemented.

**Q: Can I share artifacts between stacks?**
A: Yes, download from one stack and upload to another.

**Q: How long are artifacts kept?**
A: Forever unless you use `cleanup`. Implement retention policies for automatic cleanup.

**Q: Can I use artifacts with local workflows?**
A: Yes, artifacts work with both local and remote agent execution.

---

*Last updated: 2025-10-10*
*Version: 1.0.0*
