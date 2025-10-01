# ☁️ Multi-Cloud Deployment Example

Deploy applications across multiple cloud providers.

## Overview

Sloth Runner supports deployment to:
- ☁️ AWS
- 🔷 Azure
- 🌩️ GCP
- 🌊 DigitalOcean

## Example: Deploy to Multiple Clouds

```lua
local aws = require("aws")
local azure = require("azure")
local gcp = require("gcp")
local log = require("log")

-- Deploy to AWS
local aws_deploy = task("deploy_aws")
    :description("Deploy to AWS")
    :command(function()
        log.info("☁️ Deploying to AWS...")
        local result = aws.deploy({
            region = "us-east-1",
            service = "ecs",
            image = "myapp:latest"
        })
        return result.success
    end)
    :build()

-- Deploy to Azure
local azure_deploy = task("deploy_azure")
    :description("Deploy to Azure")
    :command(function()
        log.info("🔷 Deploying to Azure...")
        local result = azure.deploy({
            region = "eastus",
            service = "container-instances",
            image = "myapp:latest"
        })
        return result.success
    end)
    :build()

-- Deploy to GCP
local gcp_deploy = task("deploy_gcp")
    :description("Deploy to GCP")
    :command(function()
        log.info("🌩️ Deploying to GCP...")
        local result = gcp.deploy({
            region = "us-central1",
            service = "cloud-run",
            image = "myapp:latest"
        })
        return result.success
    end)
    :build()

-- Multi-cloud workflow
workflow.define("multi_cloud", {
    description = "Deploy to multiple clouds",
    tasks = {
        aws_deploy,
        azure_deploy,
        gcp_deploy
    },
    parallel = true  -- Deploy to all clouds simultaneously
})
```

## Features

- ✅ Parallel deployment
- ✅ Provider-specific configuration
- ✅ Unified interface
- ✅ Automatic failover

## Learn More

- [Multi-Cloud Excellence](../../multi-cloud-excellence.md)
- [AWS Module](../../modules/aws.md)
- [Azure Module](../../modules/azure.md)
- [GCP Module](../../modules/gcp.md)
