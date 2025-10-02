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
task("deploy_aws", {
    description = "Deploy to AWS",
    command = function()
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
    end
})

-- Deploy to Azure
task("deploy_azure", {
    description = "Deploy to Azure",
    command = function()
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
    end
})

-- Deploy to GCP
task("deploy_gcp", {
    description = "Deploy to GCP",
    command = function()
        log.info("🌩️ Deploying to GCP...")
        local result = gcp.exec({
            "storage", "rsync", "-r", "./build",
            "gs://my-app-bucket/"
        })
        if result.exit_code ~= 0 then
            return false, "GCP deployment failed: " .. result.stderr
        end
        return true, "GCP deployment completed"
    end
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
