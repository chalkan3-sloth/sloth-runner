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
-- Build stage
local build_task = task("build")
    :description("Build application")
    :command(function(this, params)
        log.info("🔨 Building...")
        local success, output = exec.run("go build -o app ./cmd")
        if not success then
            return false, "Build failed: " .. output
        end
        return true, output
    end)
    :build()

-- Test stage
local test_task = task("test")
    :description("Run tests")
    :depends_on({"build"})
    :command(function(this, params)
        log.info("🧪 Testing...")
        local success, output = exec.run("go test -v ./...")
        if not success then
            return false, "Tests failed: " .. output
        end
        return true, output
    end)
    :build()

-- Deploy stage
local deploy_task = task("deploy")
    :description("Deploy to production")
    :depends_on({"build", "test"})
    :command(function(this, params)
        log.info("🚀 Deploying...")
        local success, output = exec.run("kubectl apply -f k8s/")
        if not success then
            return false, "Deploy failed: " .. output
        end
        return true, "Deployment completed successfully"
    end)
    :build()

-- Define workflow
workflow
    .define("cicd_pipeline")
    :description("Complete CI/CD pipeline")
    :version("1.0.0")
    :tasks({build_task, test_task, deploy_task})
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
