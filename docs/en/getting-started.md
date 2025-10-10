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
task("hello")
  :description("My first task")
  :command(function(this, params)
    print("🦥 Hello from Sloth Runner!")
    return true, "Task completed successfully"
  end)
  :build()
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
task("build")
  :description("Build the application")
  :command("go build -o app")
  :build()

task("test")
  :description("Run tests")
  :command("go test ./...")
  :depends_on("build")
  :build()
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
task("deploy")
  :description("Deploy to production")
  :condition(function(this, params) return os.getenv("ENV") == "prod" end)
  :command(function(this, params)
    log.info("Deploying...")
    return exec.run("kubectl apply -f k8s/")
  end)
  :on_success(function(this, params)
    log.success("✅ Deployed successfully!")
  end)
  :on_error(function(this, params, err)
    log.error("❌ Deployment failed: " .. err)
  end)
  :timeout(300)
  :retry(3)
  :build()
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

task("deploy")
  :command(function(this, params)
    log.info("Deploying to " .. env)
    log.info("Replicas: " .. replicas)
    return true, "Deployment configuration applied"
  end)
  :build()
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

task("deploy_container")
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
task("lint")
  :command("golangci-lint run")
  :build()

task("test")
  :command("go test -v ./...")
  :depends_on("lint")
  :build()

task("build")
  :command("go build -o app")
  :depends_on("test")
  :build()

task("deploy")
  :command(function(this, params)
    exec.run("docker build -t myapp .")
    exec.run("docker push myapp")
    exec.run("kubectl rollout restart deployment/myapp")
    return true, "Deployment completed"
  end)
  :depends_on("build")
  :build()
```

Run the pipeline:

```bash
sloth-runner run -f pipeline.sloth -o rich
```

### Infrastructure Automation

```lua
local terraform = require("terraform")

task("plan")
  :command(function(this, params)
    return terraform.plan({
      dir = "./terraform",
      var_file = "prod.tfvars"
    })
  end)
  :build()

task("apply")
  :command(function(this, params)
    return terraform.apply({
      dir = "./terraform",
      auto_approve = true
    })
  end)
  :depends_on("plan")
  :build()
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
task("deploy_web")
  :target_agent("web-01")
  :command("nginx -s reload")
  :build()

task("backup_db")
  :target_agent("db-01")
  :command("pg_dump mydb > backup.sql")
  :build()
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
task("install"):command("apt-get install nginx"):build()

-- ✅ Use built-in modules
local pkg = require("pkg")
task("install"):command(function(this, params)
  return pkg.install("nginx")
end):build()
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