# 🔄 CI/CD Pipeline Example

Complete CI/CD pipeline example using Sloth Runner.

## Overview

This example demonstrates a full CI/CD pipeline with:
- Build
- Test
- Deploy
- Monitoring

## Complete Example

```lua
local exec = require("exec")
local git = require("git")
local log = require("log")

-- Build stage
local build_task = task("build")
    :description("Build application")
    :command(function()
        log.info("🔨 Building...")
        local result = exec.run("go build -o app ./cmd")
        return result.success, result.stdout
    end)
    :build()

-- Test stage
local test_task = task("test")
    :description("Run tests")
    :depends_on({"build"})
    :command(function()
        log.info("🧪 Testing...")
        local result = exec.run("go test -v ./...")
        return result.success, result.stdout
    end)
    :build()

-- Deploy stage
local deploy_task = task("deploy")
    :description("Deploy to production")
    :depends_on({"build", "test"})
    :command(function()
        log.info("🚀 Deploying...")
        local result = exec.run("kubectl apply -f k8s/")
        return result.success, result.stdout
    end)
    :build()

-- Complete CI/CD workflow
workflow.define("cicd_pipeline", {
    description = "Complete CI/CD pipeline",
    tasks = { build_task, test_task, deploy_task },
    
    on_success = function()
        log.info("✅ Pipeline completed successfully!")
    end,
    
    on_failure = function(error)
        log.error("❌ Pipeline failed: " .. error.message)
    end
})
```

## Features Demonstrated

- ✅ Multi-stage pipeline
- ✅ Task dependencies
- ✅ Error handling
- ✅ Logging
- ✅ Deployment automation

## Learn More

- [GitOps Features](../gitops-features.md)
- [Advanced Examples](../advanced-examples.md)
