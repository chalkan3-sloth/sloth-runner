# 🏗️ Infrastructure as Code Example

Managing infrastructure with Sloth Runner and Terraform/Pulumi.

## Overview

Use Sloth Runner to orchestrate infrastructure deployments:
- 🌍 Terraform
- 🏗️ Pulumi
- ☁️ Multi-cloud

## Terraform Example

```lua
local tf_plan = task("tf_plan")
    :description("Plan infrastructure changes")
    :command(function(this, params)
        log.info("📋 Planning...")
        local client = terraform.init("./terraform")
        local success, output = client:plan()
        if not success then
            return false, "Plan failed: " .. output
        end
        return true, "Plan completed"
    end)
    :build()

local tf_apply = task("tf_apply")
    :description("Apply infrastructure changes")
    :depends_on({"tf_plan"})
    :command(function(this, params)
        log.info("🚀 Applying...")
        local client = terraform.init("./terraform")
        local success, output = client:apply({auto_approve = true})
        if not success then
            return false, "Apply failed: " .. output
        end
        return true, "Apply completed"
    end)
    :build()

workflow
    .define("terraform_deploy")
    :description("Terraform deployment workflow")
    :version("1.0.0")
    :tasks({tf_plan, tf_apply})
```

## Pulumi Example

```lua
local pulumi_deploy = task("pulumi_deploy")
    :description("Deploy with Pulumi")
    :command(function(this, params)
        local stack = pulumi.stack({
            name = "my-org/project/production",
            workdir = "./infra"
        })

        local success, output = stack:up({yes = true})
        if not success then
            return false, "Deploy failed: " .. output
        end
        return true, "Deploy completed"
    end)
    :build()

workflow
    .define("pulumi_deploy")
    :description("Pulumi deployment workflow")
    :version("1.0.0")
    :tasks({pulumi_deploy})
```

## Learn More

- [Terraform Module](../../modules/terraform.md)
- [Pulumi Module](../../modules/pulumi.md)
- [Multi-Cloud](../../multi-cloud-excellence.md)
