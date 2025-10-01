# 🏢 Multi-Environment GitOps

Deploy and manage applications across multiple environments with GitOps.

## Overview

Sloth Runner provides native multi-environment support:

- 🔄 Development
- 🧪 Staging
- 🚀 Production
- 🌍 Multi-region

## Environment Configuration

```lua
-- Define environment-specific settings
local environments = {
    dev = {
        replicas = 1,
        resources = { cpu = "100m", memory = "128Mi" }
    },
    staging = {
        replicas = 2,
        resources = { cpu = "200m", memory = "256Mi" }
    },
    production = {
        replicas = 5,
        resources = { cpu = "1", memory = "1Gi" }
    }
}

-- Deploy to specific environment
workflow.define("multi_env_deploy", {
    environment = params.env or "dev",
    tasks = { deploy_task }
})
```

## Features

- ✅ Environment isolation
- ✅ Progressive rollout
- ✅ Environment-specific secrets
- ✅ Cross-environment promotion

## Learn More

- [GitOps Features](../gitops-features.md)
- [Stack Management](../../stack-management.md)
