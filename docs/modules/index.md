# 📦 Modules Reference

Sloth Runner provides a comprehensive set of built-in modules for common operations.

## 🌟 Featured Examples

### 🚀 Incus: Deploy Web Cluster with Parallel Execution

Deploy a complete web cluster in seconds using Incus containers and parallel goroutines:

```lua
task({
    name = "deploy-web-cluster",
    delegate_to = "incus-host-01",
    run = function()
        -- Create isolated network
        incus.network({
            name = "web-dmz",
            type = "bridge"
        }):set_config({
            ["ipv4.address"] = "10.10.0.1/24"
        }):create()

        -- Parallel deploy
        goroutine.map({"web-01", "web-02", "web-03"}, function(name)
            incus.instance({
                name = name,
                image = "ubuntu:22.04"
            }):create()
              :start()
              :wait_running()
              :exec("apt install -y nginx")
        end)
    end
})
```

**[📖 Full Incus Documentation →](./incus.md)**

---

### 📊 Facts: Intelligent Deployment Based on System State

Use system facts to make smart deployment decisions:

```lua
task({
    name = "intelligent-deploy",
    run = function()
        -- Collect system info
        local info, err = facts.get_all({ agent = "prod-server-01" })
        
        log.info("Platform: " .. info.platform.os)
        log.info("Memory: " .. string.format("%.2f GB", 
            info.memory.total / 1024 / 1024 / 1024))
        
        -- Validate requirements
        local mem_gb = info.memory.total / 1024 / 1024 / 1024
        if mem_gb < 4 then
            error("Insufficient memory!")
        end
        
        -- Check if Docker is installed
        local docker, _ = facts.get_package({ 
            agent = "prod-server-01", 
            name = "docker" 
        })
        
        if not docker.installed then
            pkg.install({ packages = {"docker.io"} })
               :delegate_to("prod-server-01")
        end
        
        -- Deploy based on architecture
        if info.platform.architecture == "arm64" then
            -- Use ARM image
        else
            -- Use x86 image
        end
    end
})
```

**[📖 Full Facts Documentation →](./facts.md)**

---

## Overview

All modules are **globally available** - no `require()` needed! Just use them directly in your tasks.

```lua
-- Old way (still works)
local pkg = require("pkg")

-- New way (recommended) - modules are global!
pkg.install({ packages = {"nginx"} })
```

## Core Modules

### ⚡ Execution & System

- **[exec](./exec.md)** - Execute shell commands and processes
- **[fs](./fs.md)** - File system operations (read, write, copy, move)
- **[net](./net.md)** - Network operations (HTTP, TCP, DNS)
- **[log](./log.md)** - Structured logging with levels

### 🧪 Testing & Validation 🔥

- **[infra_test](./infra_test.md)** - Infrastructure testing and validation (NEW!)
  - Test files, permissions, services, ports, processes
  - Remote agent testing support
  - Fail-fast validation for deployments

### 💾 Data & State

- **[state](../state-module.md)** - Persistent state management
- **[data](./data.md)** - Data processing (JSON, YAML, CSV)
- **[metrics](./metrics.md)** - Metrics collection and reporting

## Cloud Providers

### ☁️ AWS
[AWS Module Documentation](./aws.md)

Amazon Web Services integration:
- EC2, ECS, Lambda
- S3, DynamoDB
- CloudFormation
- IAM, Secrets Manager

### 🔷 Azure
[Azure Module Documentation](./azure.md)

Microsoft Azure integration:
- Virtual Machines, Container Instances
- Blob Storage, Cosmos DB
- ARM Templates
- Key Vault

### 🌩️ GCP
[GCP Module Documentation](./gcp.md)

Google Cloud Platform integration:
- Compute Engine, Cloud Run
- Cloud Storage, Firestore
- Deployment Manager
- Secret Manager

### 🌊 DigitalOcean
[DigitalOcean Module Documentation](./digitalocean.md)

DigitalOcean integration:
- Droplets, Kubernetes
- Spaces (Object Storage)
- Load Balancers
- Databases

## Infrastructure Tools

### 🐳 Docker
[Docker Module Documentation](./docker.md)

Container management:
- Build images
- Run containers
- Manage networks
- Docker Compose

### ☸️ Kubernetes
Integration via kubectl and native API

### 🌍 Terraform
[Terraform Module Documentation](./terraform.md)

Infrastructure as Code:
- Plan and apply
- State management
- Output parsing
- Multi-workspace

### 🏗️ Pulumi
[Pulumi Module Documentation](./pulumi.md)

Modern Infrastructure as Code:
- Stack management
- State backends
- Output exports
- Preview changes

### 🧂 SaltStack
[SaltStack Module Documentation](./salt.md)

Configuration management:
- Execute states
- Run commands
- Manage minions
- Highstate application

## Version Control

### 🐙 Git
[Git Module Documentation](./git.md)

Git operations:
- Clone repositories
- Commit changes
- Push/pull
- Branch management
- Tag management

## Notifications

### 🔔 Notifications
[Notifications Module Documentation](./notifications.md)

Multi-channel notifications:
- Slack
- Email
- Webhook
- Discord
- Microsoft Teams

## System Management

### ⚙️ Systemd
[Systemd Module Documentation](./systemd.md)

Linux service management:
- Start/stop services
- Enable/disable
- Status checking
- Journal logs

### 📦 Package Management

- **[pkg](./pkg.md)** - Package manager integration
  - apt (Debian/Ubuntu)
  - yum/dnf (RedHat/CentOS)
  - pacman (Arch Linux)
  - brew (macOS)

## Module Usage Patterns

### Basic Usage

```lua
-- Load module
local exec = require("exec")

-- Use module
local result = exec.run("echo 'Hello World'")
if result.success then
    print(result.stdout)
end
```

### Error Handling

```lua
local fs = require("fs")

local success, content = pcall(function()
    return fs.read("/path/to/file")
end)

if not success then
    log.error("Failed to read file: " .. content)
end
```

### Combining Modules

```lua
local git = require("git")
local exec = require("exec")
local notification = require("notification")

-- Clone repo
git.clone("https://github.com/user/repo.git", "/tmp/repo")

-- Build
exec.run("cd /tmp/repo && make build")

-- Notify
notification.slack({
    webhook = os.getenv("SLACK_WEBHOOK"),
    message = "Build completed!"
})
```

## Module Configuration

Some modules require configuration:

```lua
-- AWS credentials
local aws = require("aws")
aws.config({
    region = "us-east-1",
    access_key = os.getenv("AWS_ACCESS_KEY"),
    secret_key = os.getenv("AWS_SECRET_KEY")
})

-- Use AWS
aws.s3.upload("bucket-name", "file.txt", "/local/file.txt")
```

## Custom Modules

You can also create custom modules:

```lua
-- mymodule.lua
local M = {}

function M.hello(name)
    return "Hello, " .. name .. "!"
end

return M
```

Use it in your workflow:
```lua
local mymodule = require("mymodule")
print(mymodule.hello("World"))
```

## Best Practices

### 1. Check Return Values
Always check if operations succeeded:
```lua
local result = exec.run("command")
if not result.success then
    return false, result.stderr
end
```

### 2. Handle Errors Gracefully
Use pcall for operations that might fail:
```lua
local ok, err = pcall(function()
    fs.remove("/important/file")
end)
```

### 3. Use Environment Variables
Never hardcode credentials:
```lua
local api_key = os.getenv("API_KEY")
if not api_key then
    error("API_KEY not set")
end
```

### 4. Log Important Operations
```lua
log.info("Starting deployment...")
local result = deploy()
log.info("Deployment " .. (result.success and "succeeded" or "failed"))
```

## Module Reference Quick Links

### Core
- [exec](./exec.md) | [fs](./fs.md) | [net](./net.md) | [log](./log.md)

### Testing 🔥
- [infra_test](./infra_test.md) - Infrastructure testing and validation

### Cloud
- [AWS](./aws.md) | [Azure](./azure.md) | [GCP](./gcp.md) | [DigitalOcean](./digitalocean.md)

### Infrastructure
- [Docker](./docker.md) | [Terraform](./terraform.md) | [Pulumi](./pulumi.md) | [Salt](./salt.md)

### Tools
- [Git](./git.md) | [Pkg](./pkg.md) | [Systemd](./systemd.md) | [Notifications](./notifications.md)

### Parallel & Testing 🔥
- [Goroutine](./goroutine.md) | [infra_test](./infra_test.md)

## Learn More

- [Modern DSL Guide](../modern-dsl/index.md)
- [Core Concepts](../core-concepts.md)
- [Advanced Examples](../EXAMPLES.md)
- [Lua API Reference](../LUA_API.md)

---

**Need help?** Check the [documentation home](../index.md) or [file an issue](https://github.com/chalkan3-sloth/sloth-runner/issues).
