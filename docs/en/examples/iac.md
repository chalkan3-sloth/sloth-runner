# 🏗️ Infrastructure as Code Example

Managing infrastructure with Sloth Runner and Terraform/Pulumi.

## Overview

Use Sloth Runner to orchestrate infrastructure deployments:
- 🌍 Terraform
- 🏗️ Pulumi
- ☁️ Multi-cloud

## Terraform Example

```lua
task("tf_plan", {
    description = "Plan infrastructure changes",
    command = function()
        log.info("📋 Planning...")
        local result = terraform.plan({
            workdir = "./terraform"
        })
        if not result.success then
            return false, "Plan failed: " .. result.stderr
        end
        return true, "Plan completed"
    end
})

task("tf_apply", {
    description = "Apply infrastructure changes",
    depends_on = {"tf_plan"},
    command = function()
        log.info("🚀 Applying...")
        local result = terraform.apply({
            workdir = "./terraform",
            auto_approve = true
        })
        if not result.success then
            return false, "Apply failed: " .. result.stderr
        end
        return true, "Apply completed"
    end
})
```

## Pulumi Example

```lua
task("pulumi_deploy", {
    description = "Deploy with Pulumi",
    command = function()
        local stack = pulumi.stack({
            name = "my-org/project/production",
            workdir = "./infra"
        })
        
        local result = stack:up({ yes = true })
        if not result.success then
            return false, "Deploy failed: " .. result.stderr
        end
        return true, "Deploy completed"
    end
})
```

## Learn More

- [Terraform Module](../../modules/terraform.md)
- [Pulumi Module](../../modules/pulumi.md)
- [Multi-Cloud](../../multi-cloud-excellence.md)
