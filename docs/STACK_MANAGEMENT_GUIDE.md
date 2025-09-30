# 🗂️ Stack Management Guide

> **Pulumi-Style Stack Management** for Workflow State Persistence and Output Management

Sloth Runner provides sophisticated stack management capabilities similar to Pulumi's approach, allowing you to maintain persistent state, track execution history, and capture exported outputs from your workflows.

## 🌟 Overview

Stack management in Sloth Runner enables:

- ✅ **Persistent State**: SQLite-based state storage
- ✅ **Exported Outputs**: Capture and persist workflow outputs
- ✅ **Execution History**: Complete audit trail of runs
- ✅ **Environment Isolation**: Separate stacks for different environments
- ✅ **Unique Task IDs**: Enhanced traceability and debugging
- ✅ **JSON Output**: Machine-readable results for automation

## 🚀 Basic Stack Operations

### Create and Run a Stack

```bash
# Create a new stack by running workflow with stack name
sloth-runner run my-production-stack -f pipeline.lua

# Run with enhanced output
sloth-runner run my-production-stack -f pipeline.lua --output enhanced

# Run with JSON output for automation
sloth-runner run my-production-stack -f pipeline.lua --output json
```

### List All Stacks

```bash
sloth-runner stack list
```

Example output:
```
Workflow Stacks     

NAME                  STATUS      LAST RUN           DURATION       EXECUTIONS   DESCRIPTION
----                  ------      --------           --------       ----------   -----------
my-production-stack   completed   2025-09-30 08:22   5.192ms        1            Stack for workflow: my-production-stack
dev-environment       completed   2025-09-30 08:19   9.073s         3            Stack for workflow: dev-environment
staging-app           failed      2025-09-30 08:18   12.730ms       1            Stack for workflow: staging-app
```

### Show Stack Details

```bash
sloth-runner stack show my-production-stack
```

Example output:
```
Stack: my-production-stack

ID: 6d5ae01e-b3ae-4787-8038-ffd0fd307ac7
Description: Stack for workflow: my-production-stack
Version: 1.0.0
Status: completed
Created: 2025-09-30 08:22:45
Updated: 2025-09-30 08:22:46
Completed: 2025-09-30 08:22:46
Workflow File: pipeline.lua
Executions: 1
Last Duration: 5.192ms

Outputs
app_url: https://myapp.example.com
version: 1.2.3
environment: production
deployed_at: Mon Sep 30 08:22:45 2025

Recent Executions
STARTED         STATUS      DURATION    TASKS   SUCCESS   FAILED
-------         ------      --------    -----   -------   ------
2025-09-30 08:22   completed   5.192ms     1       1         0
```

### Delete a Stack

```bash
# Delete with confirmation
sloth-runner stack delete old-stack

# Force delete without confirmation
sloth-runner stack delete old-stack --force
```

## 📝 Exporting Outputs in Workflows

### Using runner.Export()

```lua
local deploy_task = task("deploy_application")
    :description("Deploy application and export outputs")
    :command(function(params, deps)
        log.info("🚀 Deploying application...")
        
        -- Perform deployment
        local result = exec.run("kubectl apply -f deployment.yaml")
        
        if result.success then
            -- Export outputs that will be persisted in stack
            runner.Export({
                app_url = "https://myapp.example.com",
                version = "1.2.3",
                environment = params.environment or "production",
                deployed_at = os.date(),
                health_endpoint = "https://myapp.example.com/health",
                database_url = "postgres://prod-db:5432/myapp",
                replicas = 3
            })
            
            log.info("✅ Application deployed successfully!")
            return true, result.stdout, { 
                status = "deployed",
                replicas = 3 
            }
        else
            log.error("❌ Deployment failed: " .. result.stderr)
            return false, result.stderr
        end
    end)
    :timeout("10m")
    :build()

workflow.define("production_deployment", {
    description = "Production Application Deployment",
    version = "1.0.0",
    tasks = { deploy_task }
})
```

### Using Global outputs Table

```lua
-- Alternative method using global outputs table
local setup_task = task("setup_infrastructure")
    :command(function()
        log.info("Setting up infrastructure...")
        
        -- Setup infrastructure
        local vpc_result = exec.run("terraform apply -auto-approve vpc.tf")
        local db_result = exec.run("terraform apply -auto-approve database.tf")
        
        -- Set global outputs that will be captured
        outputs = {
            vpc_id = "vpc-123456789",
            database_endpoint = "db.example.com:5432",
            security_group_id = "sg-987654321",
            created_at = os.time()
        }
        
        return true, "Infrastructure setup complete"
    end)
    :build()
```

## 📊 JSON Output for Automation

### Basic JSON Output

```bash
sloth-runner run my-stack -f workflow.lua --output json
```

Example JSON response:
```json
{
  "status": "success",
  "duration": "5.192125ms",
  "stack": {
    "id": "6d5ae01e-b3ae-4787-8038-ffd0fd307ac7",
    "name": "my-stack"
  },
  "tasks": {
    "deploy_application": {
      "status": "Success",
      "duration": "4.120ms",
      "error": ""
    },
    "run_tests": {
      "status": "Success", 
      "duration": "1.050ms",
      "error": ""
    }
  },
  "outputs": {
    "app_url": "https://myapp.example.com",
    "version": "1.2.3",
    "environment": "production",
    "deployed_at": "Mon Sep 30 08:22:45 2025"
  },
  "workflow": "my-stack",
  "execution_time": 1759237365
}
```

### Error Handling in JSON

```json
{
  "status": "failed",
  "duration": "2.341ms",
  "error": "task execution failed: deployment failed with exit code 1",
  "stack": {
    "id": "abc123...",
    "name": "failed-stack"
  },
  "tasks": {
    "failing_task": {
      "status": "Failed",
      "duration": "1.200ms",
      "error": "deployment failed with exit code 1"
    }
  },
  "outputs": {},
  "workflow": "failed-deployment",
  "execution_time": 1759237365
}
```

## 🆔 Task and Group IDs

### List Tasks with IDs

```bash
sloth-runner list -f workflow.lua
```

Example output:
```
Workflow Tasks and Groups

## Task Group: production_deployment
ID: 6e14b1ca-b8c8-4433-b813-ed4e9dabcb9c
Description: Production Application Deployment

Tasks:
NAME                    ID                     DESCRIPTION                      DEPENDS ON
----                    --                     -----------                      ----------
deploy_application      4600e1c7...           Deploy application and export    -
run_tests              8a92f3e1...           Run integration tests           deploy_application
notify_team            c4d8b2f9...           Send deployment notification     run_tests
```

### Using IDs for Debugging

```lua
local debug_task = task("debug_deployment")
    :description("Debug deployment with ID tracking")
    :command(function(params, deps)
        log.info("🔍 Debugging task ID: " .. params.task_id)
        log.info("🔍 Group ID: " .. params.group_id)
        
        -- Access previous task results by ID
        if deps.deploy_application then
            log.info("📋 Deploy task result: " .. deps.deploy_application.status)
        end
        
        return true, "Debug complete"
    end)
    :depends_on({"deploy_application"})
    :build()
```

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy with Sloth Runner
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Sloth Runner
        run: |
          curl -L https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner_linux_amd64.tar.gz | tar xz
          chmod +x sloth-runner
          sudo mv sloth-runner /usr/local/bin/
      
      - name: Deploy to Production
        id: deploy
        run: |
          # Run with JSON output for parsing
          OUTPUT=$(sloth-runner run prod-${{ github.sha }} -f deploy.lua --output json)
          echo "deployment_output=$OUTPUT" >> $GITHUB_OUTPUT
          
          # Check if deployment was successful
          STATUS=$(echo "$OUTPUT" | jq -r '.status')
          if [ "$STATUS" != "success" ]; then
            echo "Deployment failed!"
            exit 1
          fi
      
      - name: Extract Deployment Info
        run: |
          # Extract outputs from JSON
          APP_URL=$(echo '${{ steps.deploy.outputs.deployment_output }}' | jq -r '.outputs.app_url')
          VERSION=$(echo '${{ steps.deploy.outputs.deployment_output }}' | jq -r '.outputs.version')
          
          echo "🚀 Deployed successfully!"
          echo "📍 App URL: $APP_URL"
          echo "🏷️ Version: $VERSION"
```

### Jenkins Pipeline Example

```groovy
pipeline {
    agent any
    
    stages {
        stage('Deploy') {
            steps {
                script {
                    // Run deployment with JSON output
                    def result = sh(
                        script: "sloth-runner run jenkins-${env.BUILD_NUMBER} -f deploy.lua --output json",
                        returnStdout: true
                    ).trim()
                    
                    // Parse JSON result
                    def deployment = readJSON text: result
                    
                    if (deployment.status == 'success') {
                        echo "✅ Deployment successful!"
                        echo "📍 App URL: ${deployment.outputs.app_url}"
                        echo "🏷️ Version: ${deployment.outputs.version}"
                        
                        // Store outputs as build properties
                        env.APP_URL = deployment.outputs.app_url
                        env.DEPLOYED_VERSION = deployment.outputs.version
                    } else {
                        error("❌ Deployment failed: ${deployment.error}")
                    }
                }
            }
        }
        
        stage('Verify Deployment') {
            steps {
                script {
                    // Use exported outputs from previous stage
                    sh "curl -f ${env.APP_URL}/health"
                    echo "🎉 Health check passed for version ${env.DEPLOYED_VERSION}"
                }
            }
        }
    }
}
```

## 📈 Advanced Stack Management

### Stack Naming Strategies

```bash
# Environment-based naming
sloth-runner run dev-myapp -f deploy.lua
sloth-runner run staging-myapp -f deploy.lua  
sloth-runner run prod-myapp -f deploy.lua

# Feature branch deployment
sloth-runner run feature-auth-system -f deploy.lua

# Version-based deployment
sloth-runner run myapp-v1.2.3 -f deploy.lua

# Time-based deployment
sloth-runner run deploy-$(date +%Y%m%d-%H%M%S) -f deploy.lua
```

### Stack Lifecycle Management

```bash
# Create development stack
sloth-runner run dev-app -f app.lua --output enhanced

# Promote to staging
sloth-runner run staging-app -f app.lua -v staging-values.yaml

# Production deployment with audit trail
sloth-runner run prod-app -f app.lua -v prod-values.yaml --output json | tee deployment-audit.json

# Cleanup old stacks
sloth-runner stack list | grep "failed" | awk '{print $1}' | xargs -I {} sloth-runner stack delete {} --force
```

## 🛠️ Best Practices

### 1. **Meaningful Stack Names**
```bash
# ✅ Good: Environment and purpose clear
sloth-runner run prod-ecommerce-api -f deploy.lua
sloth-runner run staging-user-service -f deploy.lua

# ❌ Bad: Generic names
sloth-runner run test1 -f deploy.lua
sloth-runner run my-stack -f deploy.lua
```

### 2. **Comprehensive Output Exports**
```lua
-- ✅ Good: Export all relevant information
runner.Export({
    app_url = "https://api.example.com",
    version = "1.2.3",
    environment = "production",
    deployed_at = os.date(),
    health_endpoint = "https://api.example.com/health",
    database_url = "postgres://prod-db:5432/app",
    replicas = 3,
    resource_limits = {
        cpu = "500m",
        memory = "1Gi"
    }
})

-- ❌ Bad: Minimal or no exports
runner.Export({
    status = "deployed"
})
```

### 3. **Error Handling and Rollback**
```lua
local deploy_task = task("deploy_with_rollback")
    :command(function(params, deps)
        local current_version = state.get("current_version") or "unknown"
        
        -- Attempt deployment
        local result = exec.run("kubectl apply -f deployment.yaml")
        
        if result.success then
            -- Verify deployment health
            local health_check = net.get("https://app.example.com/health", {timeout = 30})
            
            if health_check.status_code == 200 then
                runner.Export({
                    status = "deployed",
                    previous_version = current_version,
                    current_version = params.version,
                    deployed_at = os.date()
                })
                state.set("current_version", params.version)
                return true, "Deployment successful"
            else
                -- Rollback on health check failure
                log.error("Health check failed, rolling back...")
                exec.run("kubectl rollout undo deployment/myapp")
                return false, "Deployment failed health check, rolled back"
            end
        else
            return false, "Deployment command failed: " .. result.stderr
        end
    end)
    :on_failure(function(params, error)
        log.error("💥 Deployment failed: " .. error)
        -- Additional cleanup or notification logic
    end)
    :build()
```

### 4. **Stack Monitoring and Alerting**
```lua
local monitor_task = task("post_deployment_monitoring")
    :command(function(params, deps)
        -- Set up monitoring after deployment
        monitoring.gauge("deployment_status", 1)
        monitoring.counter("deployments_total", 1)
        
        -- Export monitoring information
        runner.Export({
            monitoring_enabled = true,
            metrics_endpoint = "https://metrics.example.com/myapp",
            alerts_configured = {
                "high_error_rate",
                "slow_response_time",
                "pod_restart_rate"
            }
        })
        
        return true, "Monitoring configured"
    end)
    :depends_on({"deploy_with_rollback"})
    :build()
```

## 🔗 Related Documentation

- [🚀 Getting Started](getting-started.md)
- [📊 CLI Reference](CLI.md)
- [💾 State Management](state-module.md)
- [🎯 Advanced Features](advanced-features.md)
- [🌐 Distributed Execution](distributed-agents.md)

---

*Stack Management in Sloth Runner provides enterprise-grade workflow state management with the simplicity and power you need for modern automation workflows.*