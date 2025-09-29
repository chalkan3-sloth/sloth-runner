# 🗂️ Stack Management

Sloth Runner provides a complete stack management system similar to Pulumi, allowing you to persist workflow state and track executions over time.

## 🚀 Introduction

**Stack Management** in Sloth Runner enables:

- **Persist state** between executions
- **Track outputs** exported from pipeline
- **Complete history** of executions
- **Intuitive CLI** management
- **Isolation** by environment/project

## 📝 Basic Syntax

### Running with Stack

```bash
# New syntax - stack name as positional argument
sloth-runner run {stack-name} --file workflow.lua

# Practical examples
sloth-runner run production-app -f deploy.lua --output enhanced
sloth-runner run dev-environment -f test.lua -o rich
sloth-runner run my-cicd -f pipeline.lua
```

### Managing Stacks

```bash
# List all stacks
sloth-runner stack list

# Show stack details
sloth-runner stack show production-app

# Delete stack
sloth-runner stack delete old-environment
```

## 🎯 Core Concepts

### Stack State

Each stack maintains:

- **Unique ID** (UUID)
- **Stack name**
- **Current status** (created, running, completed, failed)
- **Exported outputs** from pipeline
- **Execution history**
- **Metadata** and configurations

### Lifecycle

1. **Creation**: Stack is automatically created on first execution
2. **Execution**: State is updated during execution
3. **Persistence**: Outputs and results are saved
4. **Reuse**: Subsequent executions reuse the stack

## 💾 State Persistence

### Database

Sloth Runner uses **SQLite** to persist state:

```
~/.sloth-runner/stacks.db
```

### Tables

- **stacks**: Main stack information
- **stack_executions**: Detailed execution history

## 📊 Exported Outputs

### Automatic Capture

The system automatically captures:

- **TaskRunner exports** (`runner.Exports`)
- **Global `outputs` variable** from Lua
- **Execution metadata**

### Export Example

```lua
local deploy_task = task("deploy")
    :command(function(params, deps)
        -- Deploy logic...
        
        -- Export outputs to stack
        runner.Export({
            app_url = "https://myapp.example.com",
            version = "1.2.3",
            environment = "production",
            deployed_at = os.date()
        })
        
        return true, "Deploy successful", deploy_info
    end)
    :build()
```

## 🖥️ CLI Interface

### Stack List

```bash
$ sloth-runner stack list

Workflow Stacks     

NAME                  STATUS     LAST RUN           DURATION     EXECUTIONS
----                  ------     --------           --------     ----------
production-app        completed  2025-09-29 19:27   6.8s         5
dev-environment       running    2025-09-29 19:25   2.1s         12
staging-api           failed     2025-09-29 19:20   4.2s         3
```

### Stack Details

```bash
$ sloth-runner stack show production-app

Stack: production-app     

ID: abc123-def456-789
Status: completed
Created: 2025-09-29 15:30:21
Updated: 2025-09-29 19:27:15
Executions: 5
Last Duration: 6.8s

     Outputs     

app_url: "https://myapp.example.com"
version: "1.2.3"
environment: "production"
deployed_at: "2025-09-29 19:27:15"

     Recent Executions     

STARTED            STATUS     DURATION   TASKS   SUCCESS   FAILED
-------            ------     --------   -----   -------   ------
2025-09-29 19:27   completed  6.8s       3       3         0
2025-09-29 18:45   completed  7.2s       3       3         0
2025-09-29 17:30   failed     4.1s       3       2         1
```

## 🎨 Output Styles

### Configurable per Execution

```bash
# Basic output (default)
sloth-runner run my-stack -f workflow.lua

# Enhanced output
sloth-runner run my-stack -f workflow.lua --output enhanced
sloth-runner run my-stack -f workflow.lua -o rich
sloth-runner run my-stack -f workflow.lua --output modern
```

### Pulumi Style

The `enhanced` output provides rich formatting similar to Pulumi:

```
🦥 Sloth Runner

     Workflow: production-app     

Started at: 2025-09-29 19:27:15

✓ build (2.1s) completed
✓ test (3.2s) completed  
✓ deploy (1.5s) completed

     Workflow Completed Successfully     

✓ production-app
Duration: 6.8s
Tasks executed: 3

     Outputs     

├─ exports:
  │ app_url: "https://myapp.example.com"
  │ version: "1.2.3"
  │ environment: "production"
```

## 🔧 Use Cases

### Separate Environments

```bash
# Development
sloth-runner run dev-app -f app.lua

# Staging  
sloth-runner run staging-app -f app.lua

# Production
sloth-runner run prod-app -f app.lua --output enhanced
```

### CI/CD Integration

```bash
# In CI/CD pipeline
sloth-runner run ${ENVIRONMENT}-${APP_NAME} -f pipeline.lua

# Examples:
sloth-runner run prod-frontend -f frontend-deploy.lua
sloth-runner run staging-api -f api-deploy.lua
```

### Monitoring

```bash
# View status of all environments
sloth-runner stack list

# Check last production deployment
sloth-runner stack show prod-app

# Clean up test environments
sloth-runner stack delete temp-test-env
```

## 🛠️ Best Practices

### Stack Naming

```bash
# Use pattern: {environment}-{application}
sloth-runner run prod-frontend -f deploy.lua
sloth-runner run staging-api -f deploy.lua
sloth-runner run dev-database -f setup.lua
```

### Output Exports

```lua
-- Export relevant information
runner.Export({
    -- Important URLs
    app_url = deploy_info.url,
    admin_url = deploy_info.admin_url,
    
    -- Version information
    version = build_info.version,
    commit_hash = build_info.commit,
    
    -- Environment settings
    environment = config.environment,
    region = config.region,
    
    -- Timestamps
    deployed_at = os.date(),
    build_time = build_info.timestamp
})
```

### Lifecycle Management

```bash
# Active development
sloth-runner run dev-app -f app.lua

# When ready for staging
sloth-runner run staging-app -f app.lua

# Deploy to production
sloth-runner run prod-app -f app.lua --output enhanced

# Clean up old environments
sloth-runner stack delete old-test-branch
```

## 🔄 Migration

### Old vs New Commands

```bash
# Before
sloth-runner run -f workflow.lua --stack my-stack

# Now
sloth-runner run my-stack -f workflow.lua
```

### Compatibility

- Existing workflows continue to work
- Stack is optional - can run without specifying
- Outputs are captured automatically when stack is used

## 📚 Next Steps

- [Output Styles](output-styles.md) - Output style configuration
- [Workflow Scaffolding](workflow-scaffolding.md) - Project creation
- [Examples](../examples/) - Practical examples
- [CLI Reference](CLI.md) - Complete command reference