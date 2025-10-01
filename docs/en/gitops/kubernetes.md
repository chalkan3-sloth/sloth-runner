# ☸️ Kubernetes Integration

Native Kubernetes integration for GitOps workflows.

## Overview

Sloth Runner integrates seamlessly with Kubernetes:

- 📦 Deploy applications
- 🔄 Manage resources
- 📊 Monitor status
- 🔁 Rolling updates

## Basic Usage

```lua
local k8s = require("kubernetes")

local deploy_task = task("k8s_deploy")
    :description("Deploy to Kubernetes")
    :command(function()
        -- Apply manifest
        local result = k8s.apply("deployment.yaml")
        
        -- Wait for rollout
        k8s.wait_for_rollout("deployment/myapp", {
            timeout = "5m"
        })
        
        return result.success
    end)
    :build()
```

## Features

### Manifest Management
- Apply/delete manifests
- Template rendering
- Diff preview

### Resource Monitoring
- Pod status
- Deployment health
- Service endpoints

### GitOps Workflow
- Git-based source of truth
- Automated sync
- Drift detection

## Learn More

- [GitOps Overview](../gitops-features.md)
- [Multi-Cloud Support](../../multi-cloud-excellence.md)
