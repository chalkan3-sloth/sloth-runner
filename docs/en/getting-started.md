# 🚀 Getting Started with Sloth Runner

Welcome to **Sloth Runner** - the AI-powered GitOps task orchestration platform! This guide will get you up and running in minutes.

---

## 📦 Installation

### Quick Install (Recommended)

Install the latest release with our automated script:

```bash
curl -sSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/main/install.sh | bash
```

This script:
- ✅ Detects your OS and architecture automatically
- ✅ Downloads the latest release from GitHub
- ✅ Installs to `/usr/local/bin`
- ✅ Verifies installation

**Note:** Requires `sudo` privileges.

### Manual Installation

Download from [GitHub Releases](https://github.com/chalkan3-sloth/sloth-runner/releases):

```bash
# Linux AMD64
wget https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner-linux-amd64.tar.gz
tar xzf sloth-runner-linux-amd64.tar.gz
sudo mv sloth-runner /usr/local/bin/

# macOS ARM64 (Apple Silicon)
wget https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner-darwin-arm64.tar.gz
tar xzf sloth-runner-darwin-arm64.tar.gz
sudo mv sloth-runner /usr/local/bin/
```

### Verify Installation

```bash
sloth-runner version
```

---

## 🎯 Quick Start

### Your First Workflow

Create a simple workflow file `hello.sloth`:

```lua
-- hello.sloth
local hello_task = task("hello")
  :description("My first task")
  :command(function(this, params)
    print("🦥 Hello from Sloth Runner!")
    return true, "Task completed successfully"
  end)
  :build()

workflow
  .define("hello_workflow")
  :description("Hello world workflow")
  :version("1.0.0")
  :tasks({hello_task})
```

Run it:

```bash
sloth-runner run -f hello.sloth
```

### Modern Output Styles

Try different output formats:

```bash
# Basic output
sloth-runner run -f hello.sloth -o basic

# Enhanced with colors and icons
sloth-runner run -f hello.sloth -o enhanced

# Rich with progress bars
sloth-runner run -f hello.sloth -o rich

# Modern UI
sloth-runner run -f hello.sloth -o modern

# JSON for CI/CD
sloth-runner run -f hello.sloth -o json
```

---

## 📚 Core Concepts

### File Extension

> **📝 Important:** Sloth Runner uses `.sloth` extension for workflow files (not `.lua`). The syntax is still Lua - only the extension changed for better identification.

### Tasks

Tasks are the building blocks. Define with the builder pattern:

```lua
local build_task = task("build")
  :description("Build the application")
  :command(function(this, params)
    return exec.run("go build -o app")
  end)
  :build()

local test_task = task("test")
  :description("Run tests")
  :depends_on({"build"})
  :command(function(this, params)
    return exec.run("go test ./...")
  end)
  :build()

workflow
  .define("ci_pipeline")
  :description("CI pipeline")
  :version("1.0.0")
  :tasks({build_task, test_task})
```

### Task Groups

Organize related tasks:

```lua
task_group("ci")
  :description("CI pipeline")
  :tasks({"build", "test", "lint"})
```

---

## 🗂️ Stack Management

Stacks provide state persistence and environment isolation.

### Create a Stack

```bash
sloth-runner stack new prod-app \
  -f deploy.sloth \
  --description "Production deployment"
```

### Run with Stack

```bash
# Run workflow on stack
sloth-runner run prod-app --yes

# Check stack state
sloth-runner stack show prod-app

# List all stacks
sloth-runner stack list
```

### Stack Benefits

- ✅ State persistence between runs
- ✅ Environment isolation
- ✅ History tracking
- ✅ Resource management

---

## 🎨 Modern DSL Features

### Task Builder API

```lua
local deploy_task = task("deploy")
  :description("Deploy to production")
  :command(function(this, params)
    -- Check environment condition
    if os.getenv("ENV") ~= "prod" then
      return false, "Not in production environment"
    end

    log.info("Deploying...")
    local success, output = exec.run("kubectl apply -f k8s/")

    if success then
      log.info("✅ Deployed successfully!")
    else
      log.error("❌ Deployment failed: " .. output)
    end

    return success, output
  end)
  :timeout("5m")
  :retries(3)
  :build()

workflow
  .define("deployment")
  :description("Production deployment")
  :version("1.0.0")
  :tasks({deploy_task})
```

### Values Files

Parameterize workflows with values files:

**values.yaml:**
```yaml
environment: production
replicas: 3
image: myapp:v1.2.3
```

**deploy.sloth:**
```lua
local env = values.environment
local replicas = values.replicas

local deploy_task = task("deploy")
  :description("Deploy application")
  :command(function(this, params)
    log.info("Deploying to " .. env)
    log.info("Replicas: " .. replicas)
    return true, "Deployment configuration applied"
  end)
  :build()

workflow
  .define("deploy_workflow")
  :description("Deploy with values")
  :version("1.0.0")
  :tasks({deploy_task})
```

Run with values:

```bash
sloth-runner run -f deploy.sloth -v values.yaml
```

---

## 🤖 Built-in Modules

Sloth Runner includes powerful built-in modules:

### Example: Docker Deployment

```lua
local docker = require("docker")

local deploy_container = task("deploy_container")
  :description("Deploy nginx container")
  :command(function(this, params)
    -- Pull image
    docker.pull("nginx:latest")

    -- Run container
    docker.run({
      image = "nginx:latest",
      name = "web-server",
      ports = {"80:80"},
      detach = true
    })

    return true, "Container deployed successfully"
  end)
  :build()

workflow
  .define("docker_deploy")
  :description("Deploy Docker container")
  :version("1.0.0")
  :tasks({deploy_container})
```

### Available Modules

- 🐳 **docker** - Container management
- ☁️ **aws, azure, gcp** - Cloud providers
- 🔧 **systemd** - Service management
- 📦 **pkg** - Package management
- 📊 **metrics** - Monitoring
- 🔐 **vault** - Secrets management
- 🌍 **terraform** - Infrastructure as Code

[See all modules →](/modules/)

---

## 🎭 Common Workflows

### CI/CD Pipeline

```lua
local lint_task = task("lint")
  :description("Run linter")
  :command(function(this, params)
    return exec.run("golangci-lint run")
  end)
  :build()

local test_task = task("test")
  :description("Run tests")
  :depends_on({"lint"})
  :command(function(this, params)
    return exec.run("go test -v ./...")
  end)
  :build()

local build_task = task("build")
  :description("Build application")
  :depends_on({"test"})
  :command(function(this, params)
    return exec.run("go build -o app")
  end)
  :build()

local deploy_task = task("deploy")
  :description("Deploy application")
  :depends_on({"build"})
  :command(function(this, params)
    exec.run("docker build -t myapp .")
    exec.run("docker push myapp")
    exec.run("kubectl rollout restart deployment/myapp")
    return true, "Deployment completed"
  end)
  :build()

workflow
  .define("cicd_pipeline")
  :description("Complete CI/CD pipeline")
  :version("1.0.0")
  :tasks({lint_task, test_task, build_task, deploy_task})
```

Run the pipeline:

```bash
sloth-runner run -f pipeline.sloth -o rich
```

### Infrastructure Automation

```lua
local plan_task = task("plan")
  :description("Plan Terraform changes")
  :command(function(this, params)
    return terraform.plan({
      dir = "./terraform",
      var_file = "prod.tfvars"
    })
  end)
  :build()

local apply_task = task("apply")
  :description("Apply Terraform changes")
  :depends_on({"plan"})
  :command(function(this, params)
    return terraform.apply({
      dir = "./terraform",
      auto_approve = true
    })
  end)
  :build()

workflow
  .define("terraform_deploy")
  :description("Terraform deployment")
  :version("1.0.0")
  :tasks({plan_task, apply_task})
```

---

## 🌐 Distributed Execution

### Start Master Server

```bash
sloth-runner master --port 50053 --daemon
```

### Start Agents

On different servers:

```bash
# Web server agent
sloth-runner agent start \
  --master master.example.com:50053 \
  --name web-01 \
  --tags web,nginx

# Database agent  
sloth-runner agent start \
  --master master.example.com:50053 \
  --name db-01 \
  --tags database,postgres
```

### Distribute Tasks

```lua
local deploy_web = task("deploy_web")
  :description("Reload nginx")
  :delegate_to("web-01")
  :command(function(this, params)
    return exec.run("nginx -s reload")
  end)
  :build()

local backup_db = task("backup_db")
  :description("Backup database")
  :delegate_to("db-01")
  :command(function(this, params)
    return exec.run("pg_dump mydb > backup.sql")
  end)
  :build()

workflow
  .define("distributed_ops")
  :description("Distributed operations")
  :version("1.0.0")
  :tasks({deploy_web, backup_db})
```

---

## 📊 Web Dashboard

Start the UI for visual management:

```bash
sloth-runner ui --port 8080
```

Access at: `http://localhost:8080`

Features:
- 📈 Real-time task monitoring
- 🤖 Agent health dashboard
- 📅 Scheduler management
- 📦 Stack browser
- 📊 Metrics and analytics

---

## 🔄 Scheduler

Schedule recurring tasks:

```lua
-- In your workflow
schedule("nightly_backup")
  :cron("0 2 * * *")  -- 2 AM daily
  :task("backup")
  :build()
```

Manage from CLI:

```bash
# Enable scheduler
sloth-runner scheduler enable

# List scheduled tasks
sloth-runner scheduler list

# Disable scheduler
sloth-runner scheduler disable
```

---

## 💡 Tips & Best Practices

### 1. Use Stacks for State Management

```bash
# ✅ Good: Use stacks
sloth-runner stack new myapp
sloth-runner run myapp

# ❌ Avoid: Direct execution without state
sloth-runner run -f workflow.sloth
```

### 2. Choose the Right Output Format

```bash
# Interactive terminal
sloth-runner run -f deploy.sloth -o rich

# CI/CD pipelines
sloth-runner run -f ci.sloth -o json

# Simple scripts
sloth-runner run -f task.sloth -o basic
```

### 3. Use Values Files for Environments

```bash
# Development
sloth-runner run -f app.sloth -v dev-values.yaml

# Production
sloth-runner run -f app.sloth -v prod-values.yaml
```

### 4. Leverage Built-in Modules

```lua
-- ❌ Don't shell out unnecessarily
local install_bad = task("install")
  :description("Install nginx (bad)")
  :command(function(this, params)
    return exec.run("apt-get install nginx")
  end)
  :build()

-- ✅ Use built-in modules
local install_good = task("install")
  :description("Install nginx (good)")
  :command(function(this, params)
    return pkg.install({packages = {"nginx"}})
  end)
  :build()
```

---

## 📖 Next Steps

Now that you're started, explore more:

- 📘 [Core Concepts](/en/core-concepts/) - Deep dive into architecture
- 🎨 [Modern DSL](/modern-dsl/introduction/) - Advanced syntax
- 🔧 [CLI Reference](/en/CLI/) - All commands
- 📦 [Modules](/modules/) - Built-in capabilities
- 🎭 [Examples](/EXAMPLES/) - Real-world workflows
- 🤖 [AI Features](/en/ai-features/) - Intelligent optimization
- 🚀 [GitOps](/en/gitops-features/) - Automated deployments

---

## 🆘 Getting Help

- 📖 **Documentation**: [Full docs](https://chalkan3-sloth.github.io/sloth-runner/)
- 💬 **Community**: [GitHub Discussions](https://github.com/chalkan3-sloth/sloth-runner/discussions)
- 🐛 **Issues**: [Bug Reports](https://github.com/chalkan3-sloth/sloth-runner/issues)
- 📝 **Examples**: [Example Repository](https://github.com/chalkan3-sloth/sloth-runner/tree/main/examples)

---

**Ready to automate?** Create your first workflow and start orchestrating! 🚀

---

[English](./getting-started/) | [Português](../pt/getting-started/) | [中文](../zh/getting-started/)