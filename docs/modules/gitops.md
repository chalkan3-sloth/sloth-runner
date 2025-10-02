# 🔄 GitOps Module - Complete API Reference

The GitOps module provides native Git-driven deployment workflows with intelligent diff preview, automatic rollback, and multi-environment support.

## 📋 Module Overview

```lua
local gitops = require("gitops")
```

The GitOps module enables:

- **🌊 Declarative Workflows** - Git-driven deployment automation
- **🔍 Intelligent Diff Preview** - Visual change analysis before deployment
- **🛡️ Smart Rollback** - Automatic rollback on failure with state backup
- **🏢 Multi-Environment** - Separate workflows for dev/staging/production
- **☸️ Kubernetes Native** - First-class Kubernetes integration

## 🚀 Quick Setup

### `gitops.workflow(config)`

Create a simple GitOps workflow with minimal configuration.

```lua
local workflow = gitops.workflow({
    repo = "https://github.com/company/infrastructure",
    branch = "main",
    auto_sync = true,
    diff_preview = true,
    rollback_on_failure = true
})

-- Returns: {workflow_id: "workflow-123", repository_id: "repo-123", status: "created"}
```

**Parameters:**
- `repo` (string): Git repository URL
- `branch` (string): Git branch to track (default: "main")
- `auto_sync` (boolean): Enable automatic synchronization
- `diff_preview` (boolean): Enable diff preview before sync
- `rollback_on_failure` (boolean): Enable automatic rollback on failure

**Returns:**
```lua
{
    workflow_id = "workflow-1234567890",
    repository_id = "repo-1234567890", 
    status = "created",
    auto_sync = true,
    diff_preview = true,
    rollback_on_failure = true
}
```

## 🏗️ Repository Management

### `gitops.register_repository(config)`

Register a Git repository for GitOps workflows.

```lua
local repo_id = gitops.register_repository({
    id = "production-repo",           -- Optional custom ID
    url = "https://github.com/company/k8s-manifests",
    branch = "main",
    credentials = {                   -- Optional authentication
        type = "token",
        token = "ghp_xxxxxxxxxxxx"
    },
    poll_interval = "30s"            -- How often to check for changes
})
```

**Parameters:**
- `id` (string): Custom repository ID (auto-generated if not provided)
- `url` (string): Git repository URL
- `branch` (string): Git branch to track
- `credentials` (table): Authentication credentials
  - `type` (string): "token", "ssh", or "userpass"
  - `token` (string): Personal access token (for type="token")
  - `username` (string): Username (for type="userpass")
  - `password` (string): Password (for type="userpass")
  - `ssh_key` (string): SSH private key (for type="ssh")
- `poll_interval` (string): Polling interval for auto-sync

## 🔧 Workflow Management

### `gitops.create_workflow(config)`

Create a detailed GitOps workflow with advanced configuration.

```lua
local workflow_id = gitops.create_workflow({
    id = "production-workflow",       -- Optional custom ID
    name = "Production Infrastructure",
    repository = "production-repo",   -- Repository ID
    target_path = "k8s/production",   -- Path within repository
    auto_sync = false,                -- Manual sync for production
    diff_preview = true,
    rollback_on_failure = true,
    sync_policy = {                   -- Advanced sync configuration
        auto_prune = true,            -- Remove orphaned resources
        retry = {
            limit = 3,
            backoff = "exponential"
        },
        health_check = {
            enabled = true,
            timeout = "10m"
        }
    }
})
```

**Parameters:**
- `id` (string): Custom workflow ID
- `name` (string): Human-readable workflow name
- `repository` (string): Repository ID to use
- `target_path` (string): Path within repository to sync
- `auto_sync` (boolean): Enable automatic synchronization
- `diff_preview` (boolean): Enable diff preview
- `rollback_on_failure` (boolean): Enable automatic rollback
- `sync_policy` (table): Advanced synchronization policies

### `gitops.get_workflow_status(workflow_id)`

Get the current status of a GitOps workflow.

```lua
local status = gitops.get_workflow_status("workflow-123")
```

**Returns:**
```lua
{
    id = "workflow-123",
    name = "Production Infrastructure",
    status = "synced",                -- active | syncing | synced | failed | degraded
    auto_sync = false,
    repository = "production-repo",
    last_sync_result = {              -- Last synchronization result
        id = "sync-1234567890",
        status = "succeeded",         -- running | succeeded | failed
        start_time = "2024-01-15T10:30:00Z",
        commit_hash = "abc123def456",
        message = "Sync completed successfully",
        metrics = {
            duration = "45.2s",
            resources_processed = 15,
            resources_applied = 8,
            resources_skipped = 7,
            conflicts_resolved = 0
        }
    }
}
```

### `gitops.list_workflows()`

List all registered GitOps workflows.

```lua
local workflows = gitops.list_workflows()
-- Returns array of workflow objects
```

## 🔄 Synchronization

### `gitops.sync_workflow(workflow_id)`

Manually trigger synchronization for a workflow.

```lua
local success = gitops.sync_workflow("workflow-123")
-- Returns: true on success, false on failure
```

### `gitops.start_auto_sync()`

Start the auto-sync controller for all workflows with `auto_sync = true`.

```lua
gitops.start_auto_sync()
-- Starts background polling for all auto-sync enabled workflows
```

### `gitops.stop_auto_sync()`

Stop the auto-sync controller.

```lua
gitops.stop_auto_sync()
-- Stops all background synchronization
```

## 🔍 Diff & Preview

### `gitops.generate_diff(workflow_id)`

Generate a comprehensive diff preview for pending changes.

```lua
local diff = gitops.generate_diff("workflow-123")
```

**Returns:**
```lua
{
    workflow_id = "workflow-123",
    generated_at = "2024-01-15T10:30:00Z",
    summary = {                       -- High-level summary
        total_changes = 5,
        created_resources = 2,
        updated_resources = 2,
        deleted_resources = 1,
        conflict_count = 0,
        warning_count = 1
    },
    changes = {                       -- Detailed changes
        {
            type = "create",          -- create | update | delete
            resource = "Deployment/web-app",
            desired_state = {...},    -- New resource definition
            diff = "+ Creating Deployment/web-app with 3 replicas",
            impact = "medium"         -- low | medium | high | critical
        },
        {
            type = "update", 
            resource = "Service/web-svc",
            current_state = {...},    -- Current resource state
            desired_state = {...},    -- Desired resource state
            diff = "~ Updating Service/web-svc:\n  port: 80 -> 8080",
            impact = "low"
        }
    },
    conflicts = {                     -- Detected conflicts
        {
            resource = "ConfigMap/app-config",
            type = "validation",      -- resource_exists | out_of_sync | validation
            description = "Resource modified outside of GitOps",
            current_state = {...},
            desired_state = {...},
            suggestions = [
                "Review manual changes before proceeding",
                "Consider updating the Git repository"
            ]
        }
    },
    warnings = [                      -- Warnings and recommendations
        "High-impact change detected: Deployment/critical-app"
    ]
}
```

### `gitops.preview_changes(workflow_id)`

Alias for `gitops.generate_diff()` for better readability.

```lua
local preview = gitops.preview_changes("workflow-123")
-- Same as gitops.generate_diff()
```

## 🛡️ Rollback

### `gitops.rollback_workflow(workflow_id, reason)`

Rollback a workflow to its previous state.

```lua
local success = gitops.rollback_workflow("workflow-123", "Health check failed")
-- Returns: true on success, false on failure
```

**Parameters:**
- `workflow_id` (string): Workflow to rollback
- `reason` (string): Reason for rollback (for audit logging)

## 🎯 Complete Examples

### Multi-Environment Setup

```lua
local gitops = require("gitops")
local log = require("log")

-- Define environments
local environments = {
    {
        name = "development",
        repo = "https://github.com/company/k8s-dev",
        branch = "develop",
        auto_sync = true,
        sync_interval = "5m"
    },
    {
        name = "staging",
        repo = "https://github.com/company/k8s-staging",
        branch = "staging", 
        auto_sync = true,
        sync_interval = "10m"
    },
    {
        name = "production",
        repo = "https://github.com/company/k8s-prod",
        branch = "main",
        auto_sync = false,      -- Manual deployments in production
        approval_required = true
    }
}

-- Create workflows for all environments
local workflows = {}
for _, env in ipairs(environments) do
    -- Register repository
    local repo_id = gitops.register_repository({
        id = env.name .. "-repo",
        url = env.repo,
        branch = env.branch
    })
    
    -- Create workflow
    local workflow_id = gitops.create_workflow({
        id = env.name .. "-workflow",
        name = env.name .. " Environment", 
        repository = repo_id,
        target_path = "manifests",
        auto_sync = env.auto_sync,
        diff_preview = true,
        rollback_on_failure = true
    })
    
    workflows[env.name] = workflow_id
    log.info("Created GitOps workflow for " .. env.name .. ": " .. workflow_id)
end

-- Start auto-sync controller
gitops.start_auto_sync()
```

### Production Deployment with Validation

```lua
local production_deploy = task("production_deploy")
    :description("Production deployment with full GitOps validation")
    :command(function(params, deps)
        local workflow_id = workflows.production
        
        -- Step 1: Generate diff and validate
        log.info("🔍 Analyzing changes for production deployment...")
        local diff = gitops.generate_diff(workflow_id)
        
        if not diff then
            log.info("ℹ️ No changes detected")
            return {success = true, message = "No changes to deploy"}
        end
        
        -- Step 2: Display change summary
        log.info("📊 Production Deployment Summary:")
        log.info("  📝 Total changes: " .. diff.summary.total_changes)
        log.info("  ✨ Created: " .. diff.summary.created_resources)
        log.info("  🔄 Updated: " .. diff.summary.updated_resources)
        log.info("  🗑️ Deleted: " .. diff.summary.deleted_resources)
        
        -- Step 3: Check for conflicts and high-impact changes
        if diff.summary.conflict_count > 0 then
            log.error("💥 Conflicts detected - manual resolution required")
            return {success = false, message = "Conflicts must be resolved"}
        end
        
        local high_impact_changes = 0
        for _, change in ipairs(diff.changes) do
            if change.impact == "high" or change.impact == "critical" then
                high_impact_changes = high_impact_changes + 1
                log.warn("⚠️ High-impact: " .. change.resource .. " (" .. change.type .. ")")
            end
        end
        
        -- Step 4: Show warnings
        if #diff.warnings > 0 then
            log.warn("⚠️ Warnings:")
            for _, warning in ipairs(diff.warnings) do
                log.warn("  • " .. warning)
            end
        end
        
        -- Step 5: Require approval for production
        if high_impact_changes > 0 then
            print("🔒 High-impact changes detected. Proceed? (y/N)")
            local response = io.read()
            if response:lower() ~= "y" then
                return {success = false, message = "Deployment cancelled"}
            end
        end
        
        -- Step 6: Execute deployment
        log.info("🚀 Executing production deployment...")
        local sync_success = gitops.sync_workflow(workflow_id)
        
        if not sync_success then
            log.error("💥 Production deployment failed!")
            return {success = false, message = "Deployment failed"}
        end
        
        -- Step 7: Verify deployment
        log.info("🔍 Verifying deployment...")
        local status = gitops.get_workflow_status(workflow_id)
        
        if status.status == "synced" and status.last_sync_result.status == "succeeded" then
            log.info("✅ Production deployment successful!")
            log.info("📊 Applied " .. status.last_sync_result.metrics.resources_applied .. " resources")
            log.info("⏱️ Completed in " .. status.last_sync_result.metrics.duration)
            return {success = true, message = "Production deployed successfully"}
        else
            log.error("💥 Deployment verification failed!")
            
            -- Automatic rollback
            log.warn("🔄 Initiating automatic rollback...")
            local rollback_success = gitops.rollback_workflow(workflow_id, "Deployment verification failed")
            
            if rollback_success then
                log.info("✅ Automatic rollback completed")
                return {success = false, message = "Deployment failed, rollback successful"}
            else
                log.error("💥 Rollback also failed!")
                return {success = false, message = "Deployment and rollback both failed"}
            end
        end
    end)
    :build()
```

### Kubernetes-Specific GitOps

```lua
local k8s_deploy = task("kubernetes_gitops_deploy")
    :description("Kubernetes-native GitOps deployment")
    :command(function(params, deps)
        local workflow_id = params.workflow_id
        
        -- Generate diff with Kubernetes-specific analysis
        local diff = gitops.generate_diff(workflow_id)
        
        -- Kubernetes-specific validations
        local k8s_issues = {}
        for _, change in ipairs(diff.changes) do
            -- Check for dangerous Kubernetes operations
            if change.type == "delete" then
                if change.resource:match("Namespace") then
                    table.insert(k8s_issues, "🚨 CRITICAL: Deleting namespace " .. change.resource)
                elseif change.resource:match("PersistentVolume") then
                    table.insert(k8s_issues, "⚠️ WARNING: Deleting PersistentVolume " .. change.resource)
                end
            end
            
            if change.type == "update" and change.resource:match("Deployment") then
                log.info("📦 Deployment update: " .. change.resource)
                -- Could add image change detection here
            end
        end
        
        if #k8s_issues > 0 then
            log.warn("🚨 Kubernetes-specific issues detected:")
            for _, issue in ipairs(k8s_issues) do
                log.warn("  " .. issue)
            end
            
            print("Proceed despite Kubernetes warnings? (y/N)")
            local response = io.read()
            if response:lower() ~= "y" then
                return {success = false, message = "Deployment cancelled due to K8s issues"}
            end
        end
        
        -- Execute Kubernetes deployment
        local sync_success = gitops.sync_workflow(workflow_id)
        
        if sync_success then
            -- Kubernetes-specific post-deployment checks
            log.info("🔍 Running Kubernetes health checks...")
            
            -- Could add kubectl-based health checks here
            -- kubectl get pods --all-namespaces
            -- kubectl get services
            -- kubectl get ingress
            
            return {success = true, message = "Kubernetes deployment successful"}
        else
            return {success = false, message = "Kubernetes deployment failed"}
        end
    end)
    :build()
```

## 🎯 Best Practices

### 1. **Environment Separation**
```lua
-- Use different repositories for different environments
local env_repos = {
    dev = "company/k8s-dev",
    staging = "company/k8s-staging", 
    prod = "company/k8s-prod"
}
```

### 2. **Always Preview in Production**
```lua
-- Never deploy to production without reviewing changes
if environment == "production" then
    local diff = gitops.generate_diff(workflow_id)
    if diff.summary.conflict_count > 0 or has_high_impact_changes(diff) then
        -- Require manual approval
    end
end
```

### 3. **Descriptive Rollback Reasons**
```lua
-- Provide clear audit trail
gitops.rollback_workflow(workflow_id, "Health check failed after 5 minutes - CPU usage > 90%")
```

### 4. **Monitor Sync Results**
```lua
-- Always verify deployment success
local status = gitops.get_workflow_status(workflow_id)
if status.last_sync_result.status ~= "succeeded" then
    -- Handle failure appropriately
end
```

### 5. **Use Auto-Sync Judiciously**
```lua
-- Auto-sync for dev/staging, manual for production
local auto_sync = environment ~= "production"
```

## 🔧 Advanced Features

### Custom Sync Policies

```lua
local workflow_id = gitops.create_workflow({
    name = "Advanced Sync Policy",
    repository = repo_id,
    sync_policy = {
        auto_prune = true,            -- Remove resources not in Git
        retry = {
            limit = 5,
            backoff = "exponential",  -- exponential | linear | fixed
            max_duration = "10m"
        },
        health_check = {
            enabled = true,
            timeout = "10m",
            failure_mode = "rollback"  -- ignore | fail | rollback
        },
        pre_sync_hooks = [            -- Commands to run before sync
            "kubectl cluster-info",
            "helm repo update"
        ],
        post_sync_hooks = [           -- Commands to run after sync
            "kubectl rollout status deployment/app",
            "curl -f http://app/health"
        ]
    }
})
```

### Multi-Repository Coordination

```lua
-- Coordinate deployments across multiple repositories
local repos = {
    frontend = gitops.workflow({repo = "company/frontend-config"}),
    backend = gitops.workflow({repo = "company/backend-config"}),
    database = gitops.workflow({repo = "company/database-config"})
}

-- Deploy in dependency order
gitops.sync_workflow(repos.database.workflow_id)
gitops.sync_workflow(repos.backend.workflow_id) 
gitops.sync_workflow(repos.frontend.workflow_id)
```

## 🚀 Integration Examples

### With AI Module

```lua
local ai = require("ai")
local gitops = require("gitops")

local intelligent_deploy = task("ai_gitops_deploy")
    :command(function(params, deps)
        local deploy_cmd = "kubectl apply -f manifests/"
        
        -- AI failure prediction before GitOps deployment
        local prediction = ai.predict_failure("ai_gitops_deploy", deploy_cmd)
        
        if prediction.failure_probability > 0.25 then
            log.warn("🤖 AI detected high deployment risk: " .. 
                    string.format("%.1f%%", prediction.failure_probability * 100))
            
            for _, rec in ipairs(prediction.recommendations) do
                log.info("💡 AI Recommendation: " .. rec)
            end
        end
        
        -- GitOps deployment with AI insights
        local workflow_id = params.gitops_workflow_id
        local success = gitops.sync_workflow(workflow_id)
        
        -- Record execution for AI learning
        ai.record_execution({
            task_name = "ai_gitops_deploy",
            command = deploy_cmd,
            success = success,
            execution_time = "30s",
            ai_prediction_used = true,
            predicted_failure_probability = prediction.failure_probability
        })
        
        return {success = success}
    end)
    :build()
```

### With Modern DSL Workflows

```lua
workflow.define("gitops_pipeline", {
    description = "Complete GitOps deployment pipeline",
    version = "2.0.0",
    
    metadata = {
        author = "DevOps Team",
        tags = {"gitops", "kubernetes", "production"}
    },
    
    tasks = {
        production_deploy,
        k8s_deploy
    },
    
    on_task_start = function(task_name)
        log.info("🚀 Starting GitOps task: " .. task_name)
    end,
    
    on_task_complete = function(task_name, success, output)
        if success then
            log.info("✅ GitOps task completed: " .. task_name)
        else
            log.error("❌ GitOps task failed: " .. task_name)
            
            -- Could trigger emergency rollback here
            if task_name == "production_deploy" then
                log.warn("🔄 Triggering emergency rollback...")
                gitops.rollback_workflow(production_workflow_id, "Emergency rollback due to task failure")
            end
        end
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 GitOps pipeline completed successfully!")
        else
            log.error("💥 GitOps pipeline failed - check logs for details")
        end
    end
})
```

## 📚 See Also

- [GitOps Features Overview](../gitops-features.md)
- [GitOps Quick Setup](../gitops/quick-setup.md)
- [Multi-Environment GitOps](../gitops/multi-env.md)
- [Kubernetes Integration](../gitops/kubernetes.md)
- [Rollback Strategies](../gitops/rollback.md)