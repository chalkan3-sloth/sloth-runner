# 📦 Modules Reference

Sloth Runner provides a comprehensive set of built-in modules for common operations.

## Overview

Modules are Lua libraries that provide additional functionality to your workflows. They are loaded using the `require()` function.

```lua
local exec = require("exec")
local fs = require("fs")
local log = require("log")
```

## Core Modules

### ⚡ Execution & System

- **[exec](./exec.md)** - Execute shell commands and processes
- **[fs](./fs.md)** - File system operations (read, write, copy, move)
- **[net](./net.md)** - Network operations (HTTP, TCP, DNS)
- **[log](./log.md)** - Structured logging with levels

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
Package manager integration:
- apt (Debian/Ubuntu)
- yum/dnf (RedHat/CentOS)
- brew (macOS)
- chocolatey (Windows)

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

### Cloud
- [AWS](./aws.md) | [Azure](./azure.md) | [GCP](./gcp.md) | [DigitalOcean](./digitalocean.md)

### Infrastructure
- [Docker](./docker.md) | [Terraform](./terraform.md) | [Pulumi](./pulumi.md) | [Salt](./salt.md)

### Tools
- [Git](./git.md) | [Systemd](./systemd.md) | [Notifications](./notifications.md)

## Learn More

- [Modern DSL Guide](../modern-dsl/index.md)
- [Core Concepts](../core-concepts.md)
- [Advanced Examples](../EXAMPLES.md)
- [Lua API Reference](../LUA_API.md)

---

**Need help?** Check the [documentation home](../index.md) or [file an issue](https://github.com/chalkan3-sloth/sloth-runner/issues).
