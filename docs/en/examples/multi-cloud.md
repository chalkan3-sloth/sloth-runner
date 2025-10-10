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
-- Deploy to AWS
local deploy_aws = task("deploy_aws")
    :description("Deploy to AWS")
    :command(function(this, params)
        log.info("☁️ Deploying to AWS...")
        local result = aws.s3.sync({
            source = "./build",
            destination = "s3://my-app-bucket/static",
            delete = true
        })
        if not result then
            return false, "AWS deployment failed"
        end
        return true, "AWS deployment completed"
    end)
    :build()

-- Deploy to Azure
local deploy_azure = task("deploy_azure")
    :description("Deploy to Azure")
    :command(function(this, params)
        log.info("🔷 Deploying to Azure...")
        local result = azure.exec({
            "storage", "blob", "upload-batch",
            "--destination", "mycontainer",
            "--source", "./build"
        })
        if result.exit_code ~= 0 then
            return false, "Azure deployment failed: " .. result.stderr
        end
        return true, "Azure deployment completed"
    end)
    :build()

-- Deploy to GCP
local deploy_gcp = task("deploy_gcp")
    :description("Deploy to GCP")
    :command(function(this, params)
        log.info("🌩️ Deploying to GCP...")
        local result = gcp.exec({
            "storage", "rsync", "-r", "./build",
            "gs://my-app-bucket/"
        })
        if result.exit_code ~= 0 then
            return false, "GCP deployment failed: " .. result.stderr
        end
        return true, "GCP deployment completed"
    end)
    :build()

-- Multi-cloud deployment workflow
workflow
    .define("multi_cloud_deploy")
    :description("Deploy to multiple cloud providers")
    :version("1.0.0")
    :tasks({deploy_aws, deploy_azure, deploy_gcp})
    :config({
        max_parallel_tasks = 3  -- Deploy to all clouds in parallel
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
