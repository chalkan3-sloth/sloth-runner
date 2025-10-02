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
task("build", {
    description = "Build application",
    command = function()
        log.info("🔨 Building...")
        local result = exec.run({ cmd = "go build -o app ./cmd" })
        if not result.success then
            return false, "Build failed: " .. result.stderr
        end
        return true, result.stdout
    end
})

-- Test stage
task("test", {
    description = "Run tests",
    depends_on = {"build"},
    command = function()
        log.info("🧪 Testing...")
        local result = exec.run({ cmd = "go test -v ./..." })
        if not result.success then
            return false, "Tests failed: " .. result.stderr
        end
        return true, result.stdout
    end
})

-- Deploy stage
task("deploy", {
    description = "Deploy to production",
    depends_on = {"build", "test"},
    command = function()
        log.info("🚀 Deploying...")
        local result = exec.run({ cmd = "kubectl apply -f k8s/" })
        if not result.success then
            return false, "Deploy failed: " .. result.stderr
        end
        return true, "Deployment completed successfully"
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
