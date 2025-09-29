# 🔄 GitOps Native Workflows

> **🌟 World's First Native GitOps Task Runner**  
> Sloth Runner revolutionizes deployment automation with built-in GitOps workflows, making infrastructure-as-code truly seamless.

## 🌊 Overview

GitOps Native brings declarative, Git-driven deployment workflows directly into your task automation. No external tools needed - everything is built-in and ready to use.

### ✨ Key GitOps Features

#### 🔄 **Declarative Workflows**
- **Zero Configuration**: Works out-of-the-box with any Git repository
- **Multi-Repository**: Manage multiple repos in a single workflow
- **Branch Strategies**: Support for GitFlow, GitHub Flow, and custom strategies

#### 🔍 **Intelligent Diff Preview**
- **Visual Changes**: See exactly what will change before deployment
- **Conflict Detection**: Automatic detection and resolution of conflicts
- **Impact Analysis**: Understand the impact of changes before applying

#### 🛡️ **Smart Rollback**
- **Automatic Rollback**: Roll back on failure detection
- **State Backup**: Automatic backup before every deployment
- **Multiple Strategies**: Rollback to previous commit, timestamp, or state

#### 🏢 **Multi-Environment Support**
- **Environment Isolation**: Separate workflows for dev/staging/production
- **Progressive Deployment**: Automatic promotion through environments
- **Approval Gates**: Manual approval for production deployments

## 🚀 Quick Start

### Simple GitOps Workflow

```lua
local gitops = require("gitops")

-- Create a GitOps workflow with minimal configuration
local workflow = gitops.workflow({
    repo = "https://github.com/company/infrastructure",
    branch = "main",
    auto_sync = true,
    diff_preview = true,
    rollback_on_failure = true
})

-- That's it! GitOps is now active
log.info("GitOps workflow created: " .. workflow.workflow_id)
```

### Advanced Multi-Environment Setup

```lua
local gitops = require("gitops")

-- Define environments with different configurations
local environments = {
    {
        name = "development",
        repo = "https://github.com/company/k8s-dev",
        branch = "develop",
        auto_sync = true,        -- Auto-deploy in dev
        sync_interval = "5m"
    },
    {
        name = "staging",
        repo = "https://github.com/company/k8s-staging", 
        branch = "staging",
        auto_sync = true,        -- Auto-deploy in staging
        sync_interval = "10m"
    },
    {
        name = "production",
        repo = "https://github.com/company/k8s-prod",
        branch = "main",
        auto_sync = false,       -- Manual deploys in production
        approval_required = true
    }
}

-- Create workflows for all environments
for _, env in ipairs(environments) do
    local workflow_id = gitops.create_workflow({
        name = env.name .. " Environment",
        repository = gitops.register_repository({
            url = env.repo,
            branch = env.branch
        }),
        auto_sync = env.auto_sync,
        diff_preview = true,
        rollback_on_failure = true
    })
    
    log.info("Created GitOps workflow for " .. env.name .. ": " .. workflow_id)
end
```

## 🔍 Diff Preview & Change Analysis

### Preview Changes Before Deployment

```lua
local deploy_task = task("preview_and_deploy")
    :description("Preview changes before deploying")
    :command(function(params, deps)
        local workflow_id = params.workflow_id
        
        -- Generate comprehensive diff
        local diff = gitops.generate_diff(workflow_id)
        
        if not diff then
            log.info("ℹ️ No changes detected")
            return {success = true, message = "No changes to deploy"}
        end
        
        -- Display change summary
        log.info("📊 Deployment Summary:")
        log.info("  📝 Total changes: " .. diff.summary.total_changes)
        log.info("  ✨ Created: " .. diff.summary.created_resources)
        log.info("  🔄 Updated: " .. diff.summary.updated_resources)
        log.info("  🗑️ Deleted: " .. diff.summary.deleted_resources)
        
        -- Check for conflicts
        if diff.summary.conflict_count > 0 then
            log.warn("⚠️ Conflicts detected:")
            for _, conflict in ipairs(diff.conflicts) do
                log.warn("  ❌ " .. conflict.resource .. ": " .. conflict.description)
            end
            
            return {success = false, message = "Conflicts must be resolved before deployment"}
        end
        
        -- Check for high-impact changes
        local high_impact_changes = 0
        for _, change in ipairs(diff.changes) do
            if change.impact == "high" or change.impact == "critical" then
                high_impact_changes = high_impact_changes + 1
                log.warn("⚠️ High-impact change: " .. change.resource .. " (" .. change.type .. ")")
            end
        end
        
        -- Show warnings
        if #diff.warnings > 0 then
            log.warn("⚠️ Warnings:")
            for _, warning in ipairs(diff.warnings) do
                log.warn("  • " .. warning)
            end
        end
        
        -- Require confirmation for high-impact changes
        if high_impact_changes > 0 then
            print("Proceed with " .. high_impact_changes .. " high-impact changes? (y/N)")
            local response = io.read()
            if response:lower() ~= "y" then
                return {success = false, message = "Deployment cancelled by user"}
            end
        end
        
        -- Execute deployment
        log.info("🚀 Executing deployment...")
        return gitops.sync_workflow(workflow_id)
    end)
    :build()
```

## 🔄 Sync Strategies

### Automatic Synchronization

```lua
-- Enable auto-sync for non-production environments
local dev_workflow = gitops.workflow({
    repo = "https://github.com/company/dev-config",
    auto_sync = true,
    sync_interval = "5m",     -- Check for changes every 5 minutes
    diff_preview = true,
    rollback_on_failure = true
})

-- Start the auto-sync controller
gitops.start_auto_sync()
log.info("🔄 Auto-sync controller started")
```

### Manual Synchronization with Validation

```lua
local production_deploy = task("production_deploy")
    :description("Manual production deployment with full validation")
    :command(function(params, deps)
        local workflow_id = params.workflow_id
        
        -- Step 1: Generate and review diff
        local diff = gitops.generate_diff(workflow_id)
        
        -- Step 2: Run pre-deployment validations
        log.info("🔍 Running pre-deployment validations...")
        
        -- Check for breaking changes
        local breaking_changes = false
        for _, change in ipairs(diff.changes) do
            if change.type == "delete" and change.resource:match("PersistentVolume") then
                breaking_changes = true
                log.error("💥 Breaking change detected: Deleting PersistentVolume")
            end
        end
        
        if breaking_changes then
            return {success = false, message = "Breaking changes detected - manual review required"}
        end
        
        -- Step 3: Execute deployment
        log.info("🚀 Executing production deployment...")
        local sync_result = gitops.sync_workflow(workflow_id)
        
        if not sync_result then
            log.error("💥 Deployment failed!")
            return {success = false, message = "Deployment failed"}
        end
        
        -- Step 4: Verify deployment
        log.info("🔍 Verifying deployment...")
        local status = gitops.get_workflow_status(workflow_id)
        
        if status.status == "synced" and status.last_sync_result.status == "succeeded" then
            log.info("✅ Production deployment successful!")
            return {success = true, message = "Production deployed successfully"}
        else
            log.error("💥 Deployment verification failed!")
            return {success = false, message = "Deployment verification failed"}
        end
    end)
    :build()
```

## 🛡️ Rollback Strategies

### Automatic Rollback on Failure

```lua
local resilient_deploy = task("resilient_deploy")
    :description("Deploy with automatic rollback on failure")
    :command(function(params, deps)
        local workflow_id = params.workflow_id
        
        -- Deploy with automatic rollback enabled
        local sync_result = gitops.sync_workflow(workflow_id)
        
        if not sync_result then
            log.warn("🔄 Deployment failed, automatic rollback initiated")
            
            -- GitOps will automatically rollback due to rollback_on_failure = true
            -- But we can also trigger manual rollback
            local rollback_result = gitops.rollback_workflow(workflow_id, "Deployment failed")
            
            if rollback_result then
                log.info("✅ Rollback completed successfully")
                return {success = false, message = "Deployment failed but rollback successful"}
            else
                log.error("💥 Rollback failed!")
                return {success = false, message = "Deployment and rollback both failed"}
            end
        end
        
        return {success = true, message = "Deployment successful"}
    end)
    :build()
```

### Manual Rollback

```lua
local manual_rollback = task("manual_rollback")
    :description("Manual rollback to previous state")
    :command(function(params, deps)
        local workflow_id = params.workflow_id
        local reason = params.reason or "Manual rollback requested"
        
        log.info("🔄 Initiating manual rollback...")
        log.info("📋 Reason: " .. reason)
        
        local rollback_result = gitops.rollback_workflow(workflow_id, reason)
        
        if rollback_result then
            log.info("✅ Manual rollback completed successfully")
            
            -- Verify rollback
            local status = gitops.get_workflow_status(workflow_id)
            log.info("📊 Current status: " .. status.status)
            
            return {success = true, message = "Manual rollback completed"}
        else
            log.error("💥 Manual rollback failed!")
            return {success = false, message = "Manual rollback failed"}
        end
    end)
    :build()
```

## ☸️ Kubernetes Integration

### Native Kubernetes Workflows

```lua
local k8s_gitops = task("kubernetes_gitops")
    :description("GitOps for Kubernetes manifests")
    :command(function(params, deps)
        -- Create GitOps workflow for Kubernetes
        local k8s_workflow = gitops.workflow({
            repo = "https://github.com/company/k8s-manifests",
            branch = "main",
            target_path = "manifests/production",  -- Focus on specific directory
            auto_sync = false,
            diff_preview = true,
            rollback_on_failure = true
        })
        
        -- Preview Kubernetes changes
        local diff = gitops.generate_diff(k8s_workflow.workflow_id)
        
        -- Kubernetes-specific validations
        local k8s_issues = {}
        for _, change in ipairs(diff.changes) do
            -- Check for dangerous operations
            if change.type == "delete" and change.resource:match("Namespace") then
                table.insert(k8s_issues, "Deleting namespace: " .. change.resource)
            end
            
            if change.type == "update" and change.resource:match("Deployment") then
                -- Check for image changes
                log.info("📦 Deployment update detected: " .. change.resource)
            end
        end
        
        if #k8s_issues > 0 then
            log.warn("⚠️ Kubernetes issues detected:")
            for _, issue in ipairs(k8s_issues) do
                log.warn("  • " .. issue)
            end
        end
        
        -- Deploy to Kubernetes
        return gitops.sync_workflow(k8s_workflow.workflow_id)
    end)
    :build()
```

## 📊 GitOps API Reference

### Workflow Management

```lua
-- Create simple workflow
local workflow = gitops.workflow({
    repo = "https://github.com/org/repo",
    branch = "main",
    auto_sync = true,
    diff_preview = true,
    rollback_on_failure = true
})

-- Create detailed workflow
local workflow_id = gitops.create_workflow({
    name = "Production Infrastructure",
    repository = repo_id,
    target_path = "k8s/production",
    auto_sync = false,
    diff_preview = true,
    rollback_on_failure = true
})
```

### Repository Management

```lua
-- Register repository
local repo_id = gitops.register_repository({
    url = "https://github.com/company/infrastructure",
    branch = "main",
    credentials = {
        type = "token",
        token = "ghp_xxxxx"
    }
})
```

### Sync Operations

```lua
-- Manual sync
local success = gitops.sync_workflow(workflow_id)

-- Get workflow status
local status = gitops.get_workflow_status(workflow_id)

-- List all workflows
local workflows = gitops.list_workflows()
```

### Diff and Preview

```lua
-- Generate diff
local diff = gitops.generate_diff(workflow_id)

-- Alias for diff
local preview = gitops.preview_changes(workflow_id)
```

### Rollback Operations

```lua
-- Rollback workflow
local success = gitops.rollback_workflow(workflow_id, "Reason for rollback")
```

### Auto-Sync Control

```lua
-- Start auto-sync for all auto_sync=true workflows
gitops.start_auto_sync()

-- Stop auto-sync
gitops.stop_auto_sync()
```

## 🎯 Best Practices

### 1. **Environment Strategy**
```lua
-- Use different repositories for different environments
local environments = {
    dev = {repo = "company/k8s-dev", auto_sync = true},
    staging = {repo = "company/k8s-staging", auto_sync = true},
    prod = {repo = "company/k8s-prod", auto_sync = false}
}
```

### 2. **Always Preview in Production**
```lua
-- Never deploy to production without diff preview
if environment == "production" then
    local diff = gitops.generate_diff(workflow_id)
    if diff.summary.conflict_count > 0 then
        error("Conflicts detected in production deployment!")
    end
end
```

### 3. **Use Descriptive Rollback Reasons**
```lua
-- Provide clear reasons for rollbacks
gitops.rollback_workflow(workflow_id, "Health check failed after 5 minutes")
```

### 4. **Monitor Sync Results**
```lua
-- Always check sync results
local status = gitops.get_workflow_status(workflow_id)
if status.last_sync_result.status ~= "succeeded" then
    -- Handle failure
end
```

## 🔧 Advanced Configuration

### Multi-Repository Workflows

```lua
-- Coordinate multiple repositories
local frontend_workflow = gitops.workflow({
    repo = "https://github.com/company/frontend-config"
})

local backend_workflow = gitops.workflow({
    repo = "https://github.com/company/backend-config"
})

local database_workflow = gitops.workflow({
    repo = "https://github.com/company/database-config"
})

-- Deploy in sequence
gitops.sync_workflow(database_workflow.workflow_id)
gitops.sync_workflow(backend_workflow.workflow_id)
gitops.sync_workflow(frontend_workflow.workflow_id)
```

### Custom Sync Policies

```lua
local workflow_id = gitops.create_workflow({
    name = "Custom Sync Policy",
    repository = repo_id,
    sync_policy = {
        auto_prune = true,
        retry = {
            limit = 5,
            backoff = "exponential"
        },
        health_check = {
            enabled = true,
            timeout = "10m"
        }
    }
})
```

## 🧪 Examples

Explore our comprehensive [GitOps Examples](../examples/gitops/) directory:

- **Multi-Environment Deployments**: Dev/Staging/Prod workflows
- **Kubernetes GitOps**: Native K8s integration
- **Blue-Green Deployments**: Zero-downtime deployment strategies
- **Canary Releases**: Gradual rollout strategies
- **Disaster Recovery**: Backup and restore workflows

## 🚀 What's Next?

GitOps Native is continuously evolving. Upcoming features include:

- **🎯 ArgoCD Integration**: Seamless integration with ArgoCD
- **🔄 Flux Compatibility**: Work with Flux workflows  
- **📊 Advanced Metrics**: Deployment success rates and performance metrics
- **🌐 Multi-Cluster**: Deploy across multiple Kubernetes clusters
- **🛡️ Policy Enforcement**: OPA/Gatekeeper integration for policy validation

---

**🔄 Ready to revolutionize your deployments?** Start with our [GitOps Quick Setup Guide](gitops/quick-setup.md) or explore the [complete API reference](../modules/gitops.md).