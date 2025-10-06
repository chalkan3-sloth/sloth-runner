# SLOTH-RUNNER-RUN(1) - Execute Workflows and Tasks

## NAME

**sloth-runner run** - Execute workflows and tasks locally or on remote agents

## SYNOPSIS

```
sloth-runner run <stack-name> [options]
```

## DESCRIPTION

The **run** command is the primary way to execute workflows and tasks in Sloth Runner. It processes Lua workflow files (.sloth files) and executes the defined tasks either locally or on remote agents.

Key features:

- **Stack-based state management** - Every execution requires a stack name for state tracking
- **Local or remote execution** - Run tasks locally or delegate to remote agents
- **Saved workflows** - Use saved .sloth files from the database with `--sloth`
- **SSH profiles** - Execute on remote hosts via SSH (without agents)
- **Parameter passing** - Pass parameters to workflows via command line or values files
- **Interactive mode** - Step through tasks interactively
- **Multiple output formats** - Choose from basic, enhanced, rich, modern, or JSON output

## REQUIRED ARGUMENT

**<stack-name>**
The stack name is REQUIRED for all executions. Stacks provide isolated environments for:
- State management and persistence
- Execution history tracking
- Configuration and secrets management
- Multi-environment deployments (dev, staging, prod)

## OPTIONS

```
-f, --file <path>              Path to the .sloth workflow file
    --sloth <name>             Name of saved sloth file (overrides --file)
-d, --delegate-to <agent>      Execute on specified agent (can be used multiple times)
    --ssh <profile>            SSH profile for remote execution
-v, --values <path>            Path to values file for parameters
-o, --output <format>          Output style: basic, enhanced, rich, modern, json (default: basic)
    --interactive              Run in interactive mode (step-by-step)
    --yes                      Skip confirmation prompts
    --debug                    Enable debug logging
    --ssh-password-stdin       Read SSH password from stdin
```

## EXAMPLES

### Basic Execution

Execute a workflow with a stack:

```bash
sloth-runner run production --file deployments/deploy.sloth --yes
```

Execute saved workflow:

```bash
sloth-runner run staging --sloth deployment_pipeline --yes
```

Execute specific task from workflow:

```bash
sloth-runner run dev --file workflows/maintenance.sloth backup_database --yes
```

### Delegation to Remote Agents

Execute on a single agent:

```bash
sloth-runner run production \
  --file deployments/deploy.sloth \
  --delegate-to prod-web-01 \
  --yes
```

Execute on multiple agents (parallel execution):

```bash
sloth-runner run production \
  --file deployments/deploy.sloth \
  --delegate-to prod-web-01 \
  --delegate-to prod-web-02 \
  --delegate-to prod-web-03 \
  --yes
```

Mix local and delegated tasks:

```bash
# Tasks in the workflow decide if they run locally or on agents
sloth-runner run production \
  --file workflows/complex-deploy.sloth \
  --delegate-to prod-db \
  --delegate-to prod-web-01 \
  --yes
```

### Parameter Passing

Pass parameters via command line:

```bash
sloth-runner run production \
  --file deployments/deploy.sloth \
  --param version=v2.1.0 \
  --param environment=production \
  --yes
```

Use values file for parameters:

```bash
# values/production.yaml
version: v2.1.0
environment: production
replicas: 3
database_host: prod-db.internal
```

```bash
sloth-runner run production \
  --file deployments/deploy.sloth \
  --values values/production.yaml \
  --yes
```

### SSH Remote Execution

Execute on remote host via SSH (no agent needed):

```bash
sloth-runner run staging \
  --file maintenance/cleanup.sloth \
  --ssh staging-server \
  --yes
```

With password authentication:

```bash
echo "mypassword" | sloth-runner run staging \
  --file maintenance/cleanup.sloth \
  --ssh staging-server \
  --ssh-password-stdin \
  --yes
```

### Interactive Mode

Step through tasks one by one:

```bash
sloth-runner run dev \
  --file workflows/deployment.sloth \
  --interactive
```

Sample interactive session:

```
Task: build_application
Description: Build the application binary
Execute this task? [y/n/q]: y

✓ Task 'build_application' completed successfully

Task: run_tests
Description: Execute test suite
Execute this task? [y/n/q]: y

✓ Task 'run_tests' completed successfully

Task: deploy_to_staging
Description: Deploy to staging environment
Execute this task? [y/n/q]: n

⊘ Task 'deploy_to_staging' skipped

Workflow completed: 2 tasks executed, 1 skipped
```

### Output Formats

Basic output (default):

```bash
sloth-runner run dev --file workflows/build.sloth --yes
```

Enhanced output with colors and formatting:

```bash
sloth-runner run dev --file workflows/build.sloth --output enhanced --yes
```

Modern output with progress bars and fancy formatting:

```bash
sloth-runner run dev --file workflows/build.sloth --output modern --yes
```

JSON output for automation:

```bash
sloth-runner run dev --file workflows/build.sloth --output json --yes
```

Sample JSON output:

```json
{
  "stack": "dev",
  "workflow": "workflows/build.sloth",
  "status": "completed",
  "tasks_total": 5,
  "tasks_completed": 5,
  "tasks_failed": 0,
  "tasks_skipped": 0,
  "duration_seconds": 45.3,
  "started_at": "2025-10-06T17:30:00Z",
  "completed_at": "2025-10-06T17:30:45Z",
  "tasks": [
    {
      "name": "build_app",
      "status": "completed",
      "exit_code": 0,
      "duration": "12.5s"
    }
  ]
}
```

## WORKFLOW FILE FORMAT

Workflows are Lua scripts defining tasks and their execution logic:

```lua
-- Example: deployment workflow
local build = task("build_application")
    :description("Build the application")
    :command(function(this, params)
        log.info("Building version: " .. params.version)

        local result = exec.run("go build -o app ./cmd/app")
        if result.exit_code ~= 0 then
            error("Build failed: " .. result.stderr)
        end

        log.info("Build completed successfully")
        return true
    end)
    :build()

local test = task("run_tests")
    :description("Execute test suite")
    :depends_on("build_application")
    :command(function(this, params)
        log.info("Running tests...")

        local result = exec.run("go test ./...")
        if result.exit_code ~= 0 then
            error("Tests failed")
        end

        return true
    end)
    :build()

local deploy = task("deploy_app")
    :description("Deploy application")
    :depends_on("run_tests")
    :command(function(this, params)
        log.info("Deploying to " .. params.environment)

        -- Deploy logic here
        ssh.copy({
            src = "app",
            dest = params.target_host .. ":/opt/app/app",
            user = "deploy"
        })

        ssh.exec({
            host = params.target_host,
            user = "deploy",
            command = "systemctl restart app"
        })

        return true
    end)
    :build()

return {build, test, deploy}
```

## STACK MANAGEMENT

Stacks provide isolated execution environments:

```bash
# Create a stack for production
sloth-runner stack new production --description "Production environment"

# Run workflows in that stack
sloth-runner run production --file workflows/deploy.sloth --yes

# View stack history
sloth-runner stack history production

# List all stacks
sloth-runner stack list
```

Different stacks for different environments:

```bash
# Development
sloth-runner run dev --file workflows/deploy.sloth --param env=dev --yes

# Staging
sloth-runner run staging --file workflows/deploy.sloth --param env=staging --yes

# Production
sloth-runner run production --file workflows/deploy.sloth --param env=prod --yes
```

## DELEGATION PATTERNS

### Pattern 1: Specific Agent for Specific Task

```lua
local backup = task("backup_database")
    :command(function(this, params)
        -- This task will run on the delegated agent
        exec.run("pg_dump mydb > /backup/dump.sql")
        return true
    end)
    :build()
```

```bash
sloth-runner run prod --file backup.sloth --delegate-to prod-db --yes
```

### Pattern 2: Multiple Agents in Parallel

```lua
local health_check = task("check_health")
    :command(function(this, params)
        local status = http.get({url = "http://localhost:8080/health"})
        if status.code ~= 200 then
            error("Health check failed")
        end
        return true
    end)
    :build()
```

```bash
# Runs health_check on all three agents in parallel
sloth-runner run prod --file health.sloth \
  --delegate-to web-01 \
  --delegate-to web-02 \
  --delegate-to web-03 \
  --yes
```

### Pattern 3: Mixed Local and Remote

```lua
-- Build locally
local build = task("build")
    :local(true)  -- Force local execution
    :command(function()
        exec.run("docker build -t myapp .")
        return true
    end)
    :build()

-- Deploy to agent
local deploy = task("deploy")
    :depends_on("build")
    :command(function()
        -- Runs on delegated agent
        exec.run("docker pull myapp && docker restart app")
        return true
    end)
    :build()
```

```bash
sloth-runner run prod --file deploy.sloth --delegate-to prod-web --yes
```

## COMMON USE CASES

### Application Deployment

```bash
sloth-runner run production \
  --file workflows/deploy-app.sloth \
  --values values/production.yaml \
  --delegate-to prod-web-01 \
  --delegate-to prod-web-02 \
  --yes
```

### Database Migration

```bash
sloth-runner run production \
  --file workflows/db-migration.sloth \
  --param version=v2.0.0 \
  --delegate-to prod-db \
  --yes
```

### Infrastructure Provisioning

```bash
sloth-runner run staging \
  --file workflows/provision-vm.sloth \
  --param vm_name=staging-web-04 \
  --param vm_memory=4096 \
  --yes
```

### Backup and Maintenance

```bash
# Daily backup
sloth-runner run production \
  --file workflows/backup.sloth \
  --delegate-to prod-db \
  --yes

# Weekly cleanup
sloth-runner run production \
  --file workflows/cleanup.sloth \
  --delegate-to prod-web-01 \
  --yes
```

### Multi-Environment Testing

```bash
# Test in dev
sloth-runner run dev \
  --file workflows/integration-tests.sloth \
  --param test_suite=smoke \
  --yes

# Test in staging
sloth-runner run staging \
  --file workflows/integration-tests.sloth \
  --param test_suite=full \
  --yes
```

## ERROR HANDLING

### Task Failure

When a task fails, the workflow stops:

```
Error: Task 'deploy_app' failed
  Exit Code: 1
  Error: Connection timeout to database server
  Duration: 5m30s

Workflow stopped at task: deploy_app
Tasks completed: 2/5
```

### Retry Failed Tasks

```bash
# View failed execution in stack history
sloth-runner stack history production

# Re-run the workflow
sloth-runner run production --file workflows/deploy.sloth --yes
```

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
sloth-runner run dev \
  --file workflows/deploy.sloth \
  --debug \
  --yes
```

## AUTOMATION AND CI/CD

### GitLab CI

```yaml
deploy:production:
  stage: deploy
  script:
    - sloth-runner run production \
        --file deployments/deploy.sloth \
        --param version=$CI_COMMIT_TAG \
        --delegate-to prod-web-01 \
        --yes \
        --output json > deployment-result.json
  artifacts:
    reports:
      deployment-result.json
  only:
    - tags
```

### GitHub Actions

```yaml
- name: Deploy to Production
  run: |
    sloth-runner run production \
      --file workflows/deploy.sloth \
      --param version=${{ github.ref_name }} \
      --delegate-to prod-web-01 \
      --yes
```

### Cron Jobs

```bash
# Daily backup at 2 AM
0 2 * * * /usr/local/bin/sloth-runner run production --file /opt/workflows/backup.sloth --delegate-to prod-db --yes

# Weekly cleanup on Sundays at 3 AM
0 3 * * 0 /usr/local/bin/sloth-runner run production --file /opt/workflows/cleanup.sloth --yes
```

## FILES

```
.sloth-cache/state.db         Task execution state and history
.sloth-cache/sloth.db        Saved workflow files
```

## ENVIRONMENT VARIABLES

```
SLOTH_RUNNER_MASTER_ADDR     Master server address for agent delegation
```

## EXIT CODES

```
0    All tasks completed successfully
1    One or more tasks failed
2    Invalid arguments or workflow file not found
3    Stack not found or invalid
```

## SEE ALSO

- **sloth-runner-workflow(1)** - Workflow management commands
- **sloth-runner-stack(1)** - Stack management
- **sloth-runner-sloth(1)** - Saved workflow management
- **sloth-runner-agent(1)** - Agent management for delegation

## AUTHOR

Written by the Sloth Runner development team.

## COPYRIGHT

Copyright © 2025. Released under MIT License.
