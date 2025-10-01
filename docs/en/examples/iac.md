# 🏗️ Infrastructure as Code Example

Managing infrastructure with Sloth Runner and Terraform/Pulumi.

## Overview

Use Sloth Runner to orchestrate infrastructure deployments:
- 🌍 Terraform
- 🏗️ Pulumi
- ☁️ Multi-cloud

## Terraform Example

```lua
local terraform = require("terraform")
local log = require("log")

local plan_task = task("tf_plan")
    :description("Plan infrastructure changes")
    :command(function()
        log.info("📋 Planning...")
        local result = terraform.plan({
            dir = "./terraform",
            var_file = "prod.tfvars"
        })
        return result.success, result.plan
    end)
    :build()

local apply_task = task("tf_apply")
    :description("Apply infrastructure changes")
    :depends_on({"tf_plan"})
    :command(function()
        log.info("🚀 Applying...")
        local result = terraform.apply({
            dir = "./terraform",
            auto_approve = true
        })
        return result.success, result.output
    end)
    :build()

workflow.define("infrastructure", {
    description = "Manage infrastructure",
    tasks = { plan_task, apply_task }
})
```

## Pulumi Example

```lua
local pulumi = require("pulumi")

local deploy_task = task("pulumi_deploy")
    :description("Deploy with Pulumi")
    :command(function()
        local result = pulumi.up({
            stack = "production",
            project = "./infra"
        })
        return result.success
    end)
    :build()
```

## Learn More

- [Terraform Module](../../modules/terraform.md)
- [Pulumi Module](../../modules/pulumi.md)
- [Multi-Cloud](../../multi-cloud-excellence.md)
