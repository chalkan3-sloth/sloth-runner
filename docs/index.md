# 🦥 Sloth Runner - AI-Powered GitOps Task Orchestration Platform

**The world's first AI-powered task orchestration platform with native GitOps capabilities.** Sloth Runner combines intelligent optimization, predictive analytics, automated deployments, and enterprise-grade reliability into a single, powerful platform.

[![Go CI](https://github.com/chalkan3-sloth/sloth-runner/actions/workflows/ci.yml/badge.svg)](https://github.com/chalkan3-sloth/sloth-runner/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Lua Powered](https://img.shields.io/badge/Lua-Powered-purple.svg)](https://www.lua.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/chalkan3-sloth/sloth-runner/blob/main/LICENSE)

---

## 🚀 **Quick Start with GitOps**

Get started with a complete GitOps workflow in under 5 minutes:

### 1. Install Sloth Runner

```bash
curl -sSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/main/install.sh | bash
```

### 2. Run the GitOps Example

```bash
# Clone the repository
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Execute the complete GitOps workflow
sloth-runner run -f examples/deploy_git_terraform.sloth -v examples/values.yaml deploy_git_terraform
```

### 3. Watch the Magic Happen

```
✅ Repository cloned successfully
✅ Terraform initialized automatically  
✅ Infrastructure planned and validated
✅ Deployment completed successfully
```

---

## ✨ **Revolutionary Features**

### 🎯 **Modern DSL for GitOps**
*Clean, powerful Lua-based syntax designed for infrastructure workflows*

```lua
-- Complete GitOps workflow in clean, readable syntax
local clone_task = task("clone_infrastructure")
    :description("Clone Terraform infrastructure repository")
    :workdir("/tmp/infrastructure")
    :command(function(this, params)
        local git = require("git")
        
        log.info("📡 Cloning infrastructure repository...")
        local repository = git.clone(
            values.git.repository_url,
            this.workdir.get()
        )
        
        return true, "Repository cloned successfully", {
            repository_url = values.git.repository_url,
            clone_destination = this.workdir.get()
        }
    end)
    :timeout("5m")
    :retries(3, "exponential")
    :build()

local deploy_task = task("deploy_terraform")
    :description("Deploy infrastructure using Terraform")
    :command(function(this, params)
        local terraform = require("terraform")
        
        -- Terraform init runs automatically
        local client = terraform.init(this.workdir:get())
        
        -- Load configuration from values.yaml
        local tfvars = client:create_tfvars("terraform.tfvars", {
            environment = values.terraform.environment,
            region = values.terraform.region,
            instance_type = values.terraform.instance_type
        })
        
        -- Plan and apply infrastructure
        local plan_result = client:plan({ var_file = tfvars.filename })
        if plan_result.success then
            return client:apply({ 
                var_file = tfvars.filename,
                auto_approve = true 
            })
        end
        
        return false, "Terraform plan failed"
    end)
    :timeout("15m")
    :build()

-- Define the complete GitOps workflow
workflow.define("infrastructure_deployment")
    :description("Complete GitOps: Clone + Plan + Deploy")
    :version("1.0.0")
    :tasks({ clone_task, deploy_task })
    :config({
        timeout = "20m",
        max_parallel_tasks = 1
    })
    :on_complete(function(success, results)
        if success then
            log.info("🎉 Infrastructure deployed successfully!")
        end
    end)
```

### 🏗️ **Native GitOps Integration**
*Built-in support for Git and Terraform operations*

```lua
-- Git operations with automatic credential handling
local git = require("git")
local repo = git.clone("https://github.com/company/infrastructure", "/tmp/infra")
git.checkout(repo, "production")
git.pull(repo, "origin", "production")

-- Terraform lifecycle management
local terraform = require("terraform")
local client = terraform.init("/tmp/infra/terraform/")  -- Runs 'terraform init'
local plan = client:plan({ var_file = "production.tfvars" })
local apply = client:apply({ auto_approve = true })

-- Values-driven configuration
local config = {
    environment = values.terraform.environment or "production",
    region = values.terraform.region or "us-west-2",
    instance_count = values.terraform.instance_count or 3
}
```

### ⚙️ **External Configuration Management**
*Clean separation of code and configuration using values.yaml*

**values.yaml:**
```yaml
terraform:
  environment: "production"
  region: "us-west-2" 
  instance_type: "t3.medium"
  enable_monitoring: true

git:
  repository_url: "https://github.com/company/terraform-infrastructure"
  branch: "main"

workflow:
  timeout: "30m"
  max_parallel_tasks: 2
```

**Access in workflows:**
```lua
-- Load configuration from values.yaml
local terraform_config = {
    environment = values.terraform.environment,
    region = values.terraform.region,
    instance_type = values.terraform.instance_type
}
```

---

## ⚡ **Parallel Execution with Goroutines** 🚀

> **GAME CHANGER!** Execute múltiplas operações simultaneamente e reduza o tempo de deploy de **minutos para segundos**!

<div class="grid cards" markdown>

-   :material-rocket-launch:{ .lg .middle } **10x Mais Rápido**

    ---

    Deploy em 10 servidores em paralelo ao invés de sequencialmente.
    
    **Antes:** 5 minutos ⏱️  
    **Agora:** 30 segundos ⚡

-   :material-pool:{ .lg .middle } **Worker Pools**

    ---

    Controle a concorrência com worker pools para processar grandes volumes.
    
    Perfeito para APIs com rate limiting.

-   :material-clock-fast:{ .lg .middle } **Async/Await**

    ---

    Padrão moderno de programação assíncrona no Lua.
    
    Código limpo e fácil de entender.

-   :material-shield-check:{ .lg .middle } **Timeout Built-in**

    ---

    Proteção contra operações travadas com timeout automático.
    
    Seguro e confiável.

</div>

### 💡 **Exemplo Real: Deploy Paralelo**

```lua
local deploy_task = task("deploy_multi_server")
    :description("Deploy to 10 servers in parallel - 10x faster!")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        -- Lista de servidores para deploy
        local servers = {
            "web-01", "web-02", "web-03", "api-01", "api-02",
            "api-03", "db-01", "db-02", "cache-01", "cache-02"
        }
        
        log.info("🚀 Starting parallel deployment to " .. #servers .. " servers...")
        
        -- Criar handles assíncronos para cada servidor
        local handles = {}
        for _, server in ipairs(servers) do
            local handle = goroutine.async(function()
                log.info("📦 Deploying to " .. server)
                
                -- Simula deploy (upload, install, restart, health check)
                goroutine.sleep(500)
                
                return server, "deployed", os.date("%H:%M:%S")
            end)
            
            table.insert(handles, handle)
        end
        
        -- Aguardar TODOS os deploys completarem
        local results = goroutine.await_all(handles)
        
        -- Processar resultados
        log.info("📊 All " .. #results .. " servers deployed successfully!")
        
        return true, "Parallel deployment completed in ~3 seconds!"
    end)
    :timeout("2m")
    :build()

workflow.define("parallel_deployment")
    :description("Deploy to multiple servers in parallel")
    :tasks({ deploy_task })
```

**Performance Real:**

| Operação | Sequencial | Com Goroutines | Ganho |
|----------|------------|----------------|-------|
| 🚀 Deploy 10 servidores | 5 minutos | **30 segundos** | **10x** ⚡ |
| 🏥 Health check 20 serviços | 1 minuto | **5 segundos** | **12x** ⚡ |
| 📊 Processar 1000 itens | 10 segundos | **1 segundo** | **10x** ⚡ |

**[📖 Documentação Completa de Goroutines](modules/goroutine/)** | **[🧪 Mais Exemplos](https://github.com/chalkan3-sloth/sloth-runner/tree/main/examples)**

---

## 🌟 Core Features

### 🗂️ Stack Management
**Pulumi-style stack management** with persistent state, exported outputs, and execution history tracking.

- 🔒 **Persistent stack state** with SQLite in `/etc/sloth-runner/`
- 📊 **Exported outputs capture** from pipeline with JSON support
- 📈 **Complete execution history** tracking with duration metrics
- 🎯 **Environment isolation** by stack name
- 🆔 **Unique task and group IDs** for enhanced traceability
- 📋 **Task listing** with detailed relationship view
- 🗑️ **Stack deletion** with confirmation prompts
- 🎨 **Multiple output formats**: basic, enhanced, modern, json

```bash
# Create and run a stack with enhanced output
sloth-runner stack new my-production-stack -d "Production deployment" -f pipeline.sloth
sloth-runner run my-production-stack -f pipeline.sloth --output enhanced

# Run with JSON output for CI/CD integration
sloth-runner run my-stack -f workflow.sloth --output json

# List all stacks with status and metrics
sloth-runner stack list

# Show stack details with outputs and execution history
sloth-runner stack show my-production-stack

# List tasks with unique IDs and dependencies
sloth-runner list -f pipeline.sloth

# Delete stacks with confirmation
sloth-runner stack delete old-stack
sloth-runner stack delete old-stack --force  # skip confirmation
```

### 🌐 Distributed by Design
**Native master-agent architecture** with real-time streaming, automatic failover, and intelligent load balancing.

- 🔗 **gRPC-based** agent communication
- 📡 **Real-time** command streaming
- 🔄 **Automatic failover** and recovery
- ⚖️ **Intelligent load** balancing
- 🏗️ **Scalable architecture** for enterprise workloads
- 🔒 **TLS-secured** communication

```bash
# Start master server
sloth-runner master --port 50053 --daemon

# Start and manage agents
sloth-runner agent start --name worker-01 --master localhost:50053
sloth-runner agent list --master localhost:50053
sloth-runner agent run worker-01 "docker ps" --master localhost:50053
```

### 🎨 Web Dashboard & UI
**Modern web-based dashboard** for comprehensive workflow management and monitoring.

- 📊 **Real-time monitoring** dashboard
- 🎯 **Agent management** interface
- 📈 **Performance metrics** visualization
- 🔍 **Centralized logging** system
- 👥 **Team collaboration** features

```bash
# Start web dashboard
sloth-runner ui --port 8080
# Access at http://localhost:8080

# Run as daemon
sloth-runner ui --daemon --port 8080
```

### 🤖 AI/ML Integration
**Built-in artificial intelligence** capabilities for smart automation and decision making.

- 🧠 **OpenAI integration** for text processing
- 🤖 **Automated decision** making
- 📝 **Code generation** assistance
- 🔍 **Intelligent analysis** of workflows
- 🎯 **Smart recommendations**

```lua
-- AI-powered workflow optimization
local ai = require("ai")
local result = ai.openai.complete("Generate Docker build script")
local decision = ai.decide({
    cpu_usage = metrics.cpu,
    memory_usage = metrics.memory
})
```

### ⏰ Advanced Scheduling
**Enterprise-grade task scheduling** with cron-style syntax and background execution.

- ⏰ **Cron-style scheduling** syntax
- 🔄 **Background execution** daemon
- 📅 **Recurring tasks** management
- 🎯 **Event-driven** triggers
- 📊 **Schedule monitoring**

```bash
# Enable scheduler
sloth-runner scheduler enable --config scheduler.yaml

# List scheduled tasks
sloth-runner scheduler list

# Delete a scheduled task
sloth-runner scheduler delete backup-task
```

### 💾 Advanced State Management
**Built-in SQLite-based** persistent state with atomic operations, distributed locks, and TTL support.

- 🔒 **Distributed locking** mechanisms
- ⚛️ **Atomic operations** support
- ⏰ **TTL-based** data expiration
- 🔍 **Pattern-based** queries
- 🔄 **State replication** across agents

```lua
-- Advanced state operations
local state = require("state")
state.lock("deploy-resource", 30)  -- 30 second lock
state.set("config", data, 3600)    -- 1 hour TTL
state.atomic_increment("build-count")
```

### 🏗️ Project Scaffolding
**Template-based project initialization** similar to Pulumi new or Terraform init.

- 📋 **Multiple templates** (basic, cicd, infrastructure, microservices, data-pipeline)
- 🎯 **Interactive mode** with guided setup
- 📁 **Complete project** structure generation
- 🔧 **Configuration files** auto-generated

```bash
# List available templates
sloth-runner workflow list-templates

# Create new project from template
sloth-runner workflow init my-app --template cicd

# Interactive mode
sloth-runner workflow init my-app --interactive
```

### ☁️ Multi-Cloud Excellence
**Comprehensive cloud provider** support with advanced automation capabilities.

- ☁️ **AWS, GCP, Azure** native integration
- 🚀 **Terraform & Pulumi** advanced support
- 🔧 **Infrastructure as Code** automation
- 🔒 **Security & compliance** built-in
- 📊 **Cost optimization** tools

### 🔒 Enterprise Security
**Built-in security features** for enterprise compliance and data protection.

- 🔐 **Certificate management**
- 🔒 **Secret encryption** and storage
- 🛡️ **Vulnerability scanning**
- 📋 **Compliance checking**
- 📝 **Audit logging** system

### 📊 Enhanced Output System
**Pulumi-style rich output** formatting with configurable styles, progress indicators, and structured displays.

- 🎨 **Multiple output styles** (basic, enhanced, rich, modern, **json**)
- 📈 **Real-time progress** indicators
- 🎯 **Structured output** sections
- 🌈 **Rich color** formatting
- 📊 **Metrics visualization**
- 🔧 **JSON output** for automation and CI/CD integration

```bash
# Enhanced Pulumi-style output
sloth-runner run my-stack -f workflow.sloth --output enhanced

# JSON output for automation
sloth-runner run my-stack -f workflow.sloth --output json

# List tasks with unique IDs
sloth-runner list -f workflow.sloth
```

### 🔧 Rich Module Ecosystem
**Extensive collection** of pre-built modules for common automation tasks.

- 🌐 **Network & HTTP** operations
- 💽 **Database** integrations (MySQL, PostgreSQL, MongoDB, Redis)
- 📧 **Notification systems** (Email, Slack, Discord)
- 🐍 **Python/R integration** with virtual environments
- 🔗 **GitOps** advanced workflows
- 🧪 **Testing frameworks** and quality assurance

---

## 🚀 Quick Start Examples

### 🗂️ Stack Management with Pulumi-Style Output

```bash
# Create a new project from template
sloth-runner workflow init my-cicd --template cicd

# Deploy to development environment
sloth-runner run dev-app -f my-cicd.sloth --output enhanced

# Deploy to production with stack persistence
sloth-runner run prod-app -f my-cicd.sloth -o rich

# Check deployment status and outputs
sloth-runner stack show prod-app
```

### 📊 Stack with Exported Outputs & JSON Output

```lua
local deploy_task = task("deploy")
    :command(function(params, deps)
        -- Deploy application
        local result = exec.run("kubectl apply -f deployment.yaml")
        
        -- Export important outputs to stack
        runner.Export({
            app_url = "https://myapp.example.com",
            version = "1.2.3",
            environment = "production",
            deployed_at = os.date(),
            health_endpoint = "https://myapp.example.com/health"
        })
        
        return true, result.stdout, { status = "deployed" }
    end)
    :build()

workflow.define("production_deployment", {
    tasks = { deploy_task }
})
```

**Run with JSON output for automation:**
```bash
# Get structured JSON output for CI/CD integration
sloth-runner run prod-deployment -f deploy.sloth --output json

# Example JSON output:
{
  "status": "success",
  "duration": "5.192ms",
  "stack": {
    "id": "abc123...",
    "name": "prod-deployment"
  },
  "tasks": {
    "deploy": {
      "status": "Success",
      "duration": "4.120ms"
    }
  },
  "outputs": {
    "app_url": "https://myapp.example.com",
    "version": "1.2.3",
    "environment": "production"
  },
  "workflow": "production_deployment",
  "execution_time": 1759237365
}
```

---

## 📊 CLI Commands Overview

### Stack Management (NEW!)
```bash
# Execute with stack persistence (NEW SYNTAX)
sloth-runner run {stack-name} --file workflow.sloth

# Enhanced output styles
sloth-runner run {stack-name} --file workflow.sloth --output enhanced
sloth-runner run {stack-name} --file workflow.sloth --output json

# Manage stacks
sloth-runner stack list                    # List all stacks
sloth-runner stack show production-app     # Show stack details with outputs
sloth-runner stack delete old-env          # Delete stack

# List tasks with unique IDs
sloth-runner list --file workflow.sloth      # Show tasks and groups with IDs
```

### Project Scaffolding
```bash
# Create new projects
sloth-runner workflow init my-app --template cicd
sloth-runner workflow list-templates       # Available templates
```

### Distributed Agents & Web UI
```bash
# Start master server
sloth-runner master --port 50053 --daemon

# Start distributed agents
sloth-runner agent start --name web-builder --master localhost:50053
sloth-runner agent start --name db-manager --master localhost:50053

# Start web dashboard
sloth-runner ui --port 8080 --daemon
# Access dashboard at http://localhost:8080

# List connected agents
sloth-runner agent list --master localhost:50053

# Execute commands on specific agents
sloth-runner agent run web-builder "docker ps" --master localhost:50053
```

### Advanced Scheduling
```bash
# Enable background scheduler
sloth-runner scheduler enable --config scheduler.yaml

# List and manage scheduled tasks
sloth-runner scheduler list
sloth-runner scheduler delete backup-task
```

### 📊 Distributed Deployment with Monitoring

```lua
local monitoring = require("monitoring")
local state = require("state")

-- Production deployment with comprehensive monitoring
local deploy_task = task("production_deployment")
    :command(function(params, deps)
        -- Track deployment metrics
        monitoring.counter("deployments_started", 1)
        
        -- Use state for coordination
        local deploy_id = state.increment("deployment_counter", 1)
        state.set("current_deployment", deploy_id)
        
        -- Execute deployment
        local result = exec.run("kubectl apply -f production.yaml")
        
        if result.success then
            monitoring.gauge("deployment_status", 1)
            state.set("last_successful_deploy", os.time())
            log.info("✅ Deployment " .. deploy_id .. " completed successfully")
        else
            monitoring.gauge("deployment_status", 0)
            monitoring.counter("deployments_failed", 1)
            log.error("❌ Deployment " .. deploy_id .. " failed: " .. result.stderr)
        end
        
        return result
    end)
    :build()
```

### 🌐 Multi-Agent Distributed Execution

```lua
local distributed = require("distributed")

-- Execute tasks across multiple agents
workflow.define("distributed_pipeline", {
    tasks = {
        task("build_frontend")
            :agent("build-agent-1")
            :command("npm run build")
            :build(),
            
        task("build_backend")
            :agent("build-agent-2")
            :command("go build -o app ./cmd/server")
            :build(),
            
        task("run_tests")
            :agent("test-agent")
            :depends_on({"build_frontend", "build_backend"})
            :command("npm test && go test ./...")
            :build(),
            
        task("deploy")
            :agent("deploy-agent")
            :depends_on({"run_tests"})
            :command("./deploy.sh production")
            :build()
    }
})
```

### 💾 Advanced State Management

```lua
local state = require("state")

-- Complex state operations with locking
local update_config = task("update_configuration")
    :command(function(params, deps)
        -- Critical section with automatic locking
        return state.with_lock("config_update", function()
            local current_version = state.get("config_version") or 0
            local new_version = current_version + 1
            
            -- Atomic configuration update
            local success = state.compare_and_swap("config_version", current_version, new_version)
            
            if success then
                state.set("config_data", params.new_config)
                state.set("config_updated_at", os.time())
                log.info("Configuration updated to version " .. new_version)
                return { version = new_version, success = true }
            else
                log.error("Configuration update failed - version mismatch")
                return { success = false, error = "version_mismatch" }
            end
        end)
    end)
    :build()
```

### 🔄 CI/CD Pipeline with GitOps

```lua
local git = require("git")
local docker = require("docker")
local kubernetes = require("kubernetes")

-- Complete CI/CD pipeline
workflow.define("gitops_pipeline", {
    on_git_push = true,
    
    tasks = {
        task("checkout_code")
            :command(function()
                return git.clone(params.repository, "/tmp/build")
            end)
            :build(),
            
        task("run_tests")
            :depends_on({"checkout_code"})
            :command("cd /tmp/build && npm test")
            :retry_count(3)
            :build(),
            
        task("build_image")
            :depends_on({"run_tests"})
            :command(function()
                return docker.build({
                    path = "/tmp/build",
                    tag = "myapp:" .. params.git_sha,
                    push = true
                })
            end)
            :build(),
            
        task("deploy_staging")
            :depends_on({"build_image"})
            :command(function()
                return kubernetes.apply_manifest({
                    file = "/tmp/build/k8s/staging.yaml",
                    namespace = "staging",
                    image = "myapp:" .. params.git_sha
                })
            end)
            :build(),
            
        task("integration_tests")
            :depends_on({"deploy_staging"})
            :command("./run-integration-tests.sh staging")
            :build(),
            
        task("deploy_production")
            :depends_on({"integration_tests"})
            :condition(function() return params.branch == "main" end)
            :command(function()
                return kubernetes.apply_manifest({
                    file = "/tmp/build/k8s/production.yaml",
                    namespace = "production",
                    image = "myapp:" .. params.git_sha
                })
            end)
            :build()
    }
})
```

## 📊 **Module Reference**

<div class="modules-grid">
  <div class="module-category core">
    <h4>🔧 Core Modules</h4>
    <ul>
      <li><code>exec</code> - Command execution with streaming</li>
      <li><code>fs</code> - File system operations</li>
      <li><code>net</code> - Network utilities</li>
      <li><code>data</code> - Data processing utilities</li>
      <li><code>log</code> - Structured logging</li>
    </ul>
  </div>
  
  <div class="module-category state">
    <h4>💾 State & Monitoring</h4>
    <ul>
      <li><code>state</code> - Persistent state management</li>
      <li><code>metrics</code> - Monitoring and metrics</li>
      <li><code>monitoring</code> - System monitoring</li>
      <li><code>health</code> - Health check utilities</li>
    </ul>
  </div>
  
  <div class="module-category cloud">
    <h4>☁️ Cloud Providers</h4>
    <ul>
      <li><code>aws</code> - Amazon Web Services</li>
      <li><code>gcp</code> - Google Cloud Platform</li>
      <li><code>azure</code> - Microsoft Azure</li>
      <li><code>digitalocean</code> - DigitalOcean</li>
    </ul>
  </div>
  
  <div class="module-category infrastructure">
    <h4>🛠️ Infrastructure</h4>
    <ul>
      <li><code>kubernetes</code> - Kubernetes orchestration</li>
      <li><code>docker</code> - Container management</li>
      <li><code>terraform</code> - Infrastructure as Code</li>
      <li><code>pulumi</code> - Modern IaC</li>
      <li><code>salt</code> - Configuration management</li>
    </ul>
  </div>
  
  <div class="module-category integration">
    <h4>🔗 Integrations</h4>
    <ul>
      <li><code>git</code> - Git operations</li>
      <li><code>python</code> - Python script execution</li>
      <li><code>notification</code> - Alert notifications</li>
      <li><code>crypto</code> - Cryptographic operations</li>
    </ul>
  </div>
</div>

## 🎯 **Why Choose Sloth Runner?**

<div class="comparison">
  <div class="comparison-item">
    <h4>🏢 Enterprise Ready</h4>
    <ul>
      <li>🌍 Distributed execution across multiple agents</li>
      <li>🔒 Production-grade security with mTLS</li>
      <li>📊 Comprehensive monitoring and alerting</li>
      <li>💾 Reliable state management with persistence</li>
      <li>🔄 Circuit breakers and fault tolerance</li>
    </ul>
  </div>
  
  <div class="comparison-item">
    <h4>👩‍💻 Developer Experience</h4>
    <ul>
      <li>🧰 Rich Lua-based DSL for complex workflows</li>
      <li>📡 Real-time command output streaming</li>
      <li>🔄 Interactive REPL for debugging</li>
      <li>📚 Comprehensive documentation</li>
      <li>🎯 Intuitive task dependency management</li>
    </ul>
  </div>
  
  <div class="comparison-item">
    <h4>🚀 Performance & Reliability</h4>
    <ul>
      <li>⚡ High-performance parallel execution</li>
      <li>🔄 Automatic retry and error handling</li>
      <li>📈 Built-in performance monitoring</li>
      <li>🎛️ Resource optimization and throttling</li>
      <li>🛡️ Robust error recovery mechanisms</li>
    </ul>
  </div>

  <div class="comparison-item">
    <h4>🔧 Operational Excellence</h4>
    <ul>
      <li>📊 Prometheus-compatible metrics</li>
      <li>🔍 Distributed tracing support</li>
      <li>📋 Structured audit logging</li>
      <li>🚨 Flexible alerting mechanisms</li>
      <li>🔄 GitOps workflow integration</li>
    </ul>
  </div>

---

## 🚀 Get Started in Minutes

### Installation

=== "Quick Install (Recommended)"

    ```bash
    # One-line installer for Linux/macOS
    curl -sSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/main/install.sh | bash
    ```

=== "Manual Download"

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

=== "From Source"

    ```bash
    git clone https://github.com/chalkan3-sloth/sloth-runner.git
    cd sloth-runner
    go build -o sloth-runner ./cmd/sloth-runner
    sudo mv sloth-runner /usr/local/bin/
    ```

### Create Your First Workflow

!!! example "Simple Hello World"

    Create `hello.sloth`:

    ```lua
    task("hello")
      :description("My first Sloth Runner task")
      :command(function() 
        log.info("🦥 Hello from Sloth Runner!")
        return true 
      end)
      :build()

    workflow.define("greeting")
      :tasks({"hello"})
      :on_complete(function(success)
        if success then
          log.success("✅ Workflow completed!")
        end
      end)
    ```

### Run Your Workflow

=== "Basic Output"

    ```bash
    sloth-runner run -f hello.sloth
    ```

=== "Modern UI"

    ```bash
    sloth-runner run -f hello.sloth -o modern
    ```

=== "Rich Progress"

    ```bash
    sloth-runner run -f hello.sloth -o rich
    ```

---

## 📚 Learn More

<div class="grid cards" markdown>

-   :rocket:{ .lg .middle } **Quick Tutorial**

    ---

    Get up and running with practical examples in 5 minutes

    [:octicons-arrow-right-24: Start Tutorial](TUTORIAL/)

-   :material-file-document:{ .lg .middle } **Advanced Examples**

    ---

    Production-ready workflows and real-world use cases

    [:octicons-arrow-right-24: View Examples](en/advanced-examples/)

-   :brain:{ .lg .middle } **Core Concepts**

    ---

    Deep dive into Sloth Runner's architecture and features

    [:octicons-arrow-right-24: Learn Concepts](en/core-concepts/)

-   :books:{ .lg .middle } **API Reference**

    ---

    Complete documentation of all modules and functions

    [:octicons-arrow-right-24: Browse API](modules/)

-   :material-file-code:{ .lg .middle } **Modern DSL**

    ---

    Learn the modern task definition syntax

    [:octicons-arrow-right-24: DSL Guide](modern-dsl/introduction/)

-   :material-github:{ .lg .middle } **GitHub Repository**

    ---

    Source code, issues, and contributions

    [:octicons-arrow-right-24: View on GitHub](https://github.com/chalkan3-sloth/sloth-runner)

</div>

---

## 💾 **State Management & Persistence** <span class="status-indicator implemented">Implemented</span>
- **SQLite-based persistent state** with WAL mode for performance
- **Atomic operations**: increment, compare-and-swap, append
- **Distributed locks** with automatic timeout handling
- **TTL support** for automatic data expiration
- **Pattern matching** for bulk operations

```lua
-- Persistent state example
state.set("deployment_version", "v1.2.3")
local counter = state.increment("api_calls", 1)

-- Critical section with automatic locking
state.with_lock("deployment", function()
    -- Safe deployment logic
    local success = deploy_application()
    state.set("last_deploy", os.time())
    return success
end)
```

## 📊 **Metrics & Monitoring** <span class="status-indicator implemented">Implemented</span>
- **System metrics**: CPU, memory, disk, network monitoring
- **Custom metrics**: gauges, counters, histograms, timers
- **Health checks** with configurable thresholds
- **Prometheus endpoints** for external monitoring
- **Real-time alerting** based on conditions

```lua
-- Monitoring example
local cpu = metrics.system_cpu()
metrics.gauge("app_performance", response_time)
metrics.counter("requests_total", 1)

if cpu > 80 then
    metrics.alert("high_cpu", {
        level = "warning",
        message = "CPU usage critical: " .. cpu .. "%"
    })
end
```

## 🌐 **Distributed Agent System** <span class="status-indicator implemented">Implemented</span>
- **Master-agent architecture** with gRPC communication
- **Real-time streaming** of command output
- **Automatic agent registration** and health monitoring
- **Load balancing** across available agents
- **TLS encryption** for secure communication

```bash
# Start master server
sloth-runner master --port 50053

# Deploy agents on remote machines
sloth-runner agent start --name agent-1 --master master:50053

# Execute distributed commands
sloth-runner agent run agent-1 "deploy-script.sh"
```

## 📚 **Documentation by Language**

### 🇺🇸 **English Documentation**
- 📖 [Getting Started](en/getting-started/)
- 🧠 [Core Concepts](en/core-concepts/)
- ⚡ [Quick Start](en/quick-start/)
- 💻 [CLI Reference](en/CLI/)
- 🔄 [Interactive REPL](en/repl/)
- 🎯 [Advanced Features](en/advanced-features/)
- 🚀 [Agent Improvements](en/agent-improvements/)

### 🇧🇷 **Documentação em Português**
- 📖 [Primeiros Passos](pt/getting-started/)
- 🧠 [Conceitos Fundamentais](pt/core-concepts/)
- ⚡ [Início Rápido](pt/quick-start/)
- 💻 [Referência CLI](pt/CLI/)
- 🔄 [REPL Interativo](pt/repl/)
- 🎯 [Recursos Avançados](pt/advanced-features/)
- 🚀 [Melhorias dos Agentes](pt/agent-improvements/)

### 🇨🇳 **中文文档**
- 📖 [入门指南](zh/getting-started/)
- 🧠 [核心概念](zh/core-concepts/)
- ⚡ [快速开始](zh/quick-start/)
- 💻 [CLI参考](zh/CLI/)
- 🔄 [交互式REPL](zh/repl/)
- 🎯 [高级功能](zh/advanced-features/)
- 🚀 [代理改进](zh/agent-improvements/)

## 🔧 **Module Reference**

### **📦 Built-in Modules**
| Module | Description | Language Support |
|--------|-------------|------------------|
| [💾 **State**](en/modules/state/) | Persistent state management | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📊 **Metrics**](en/modules/metrics/) | Monitoring and observability | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [⚡ **Exec**](en/modules/exec/) | Command execution | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📁 **FS**](en/modules/fs/) | File system operations | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📡 **Net**](en/modules/net/) | Network operations | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📋 **Data**](en/modules/data/) | Data processing utilities | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📝 **Log**](en/modules/log/) | Structured logging | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |

### **☁️ Cloud Provider Modules**
| Module | Description | Status |
|--------|-------------|---------|
| [☁️ **AWS**](en/modules/aws/) | Amazon Web Services | <span class="status-indicator implemented">Ready</span> |
| [🌩️ **GCP**](en/modules/gcp/) | Google Cloud Platform | <span class="status-indicator implemented">Ready</span> |
| [🔷 **Azure**](en/modules/azure/) | Microsoft Azure | <span class="status-indicator implemented">Ready</span> |
| [🌊 **DigitalOcean**](en/modules/digitalocean/) | DigitalOcean | <span class="status-indicator beta">Beta</span> |

### **🛠️ Infrastructure Modules**
| Module | Description | Status |
|--------|-------------|---------|
| [🐳 **Docker**](en/modules/docker/) | Container management | <span class="status-indicator implemented">Ready</span> |
| [🏗️ **Pulumi**](en/modules/pulumi/) | Modern IaC | <span class="status-indicator implemented">Ready</span> |
| [🌍 **Terraform**](en/modules/terraform/) | Infrastructure provisioning | <span class="status-indicator implemented">Ready</span> |
| [🧂 **Salt**](en/modules/salt/) | Configuration management | <span class="status-indicator beta">Beta</span> |
| [🐍 **Python**](en/modules/python/) | Python integration | <span class="status-indicator beta">Beta</span> |

## 🚀 **Get Started Today**

```bash
# Install Sloth Runner
curl -L https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner_linux_amd64.tar.gz | tar xz
chmod +x sloth-runner && sudo mv sloth-runner /usr/local/bin/

# Create your first workflow
echo 'local hello_task = task("greet"):command(function() log.info("Hello World! 🚀") return true end):build(); workflow.define("hello", { tasks = { hello_task } })' > hello.sloth

# Run it!
sloth-runner run -f hello.sloth
```

## 🤝 **Community & Support**

<div class="community-grid">
  <a href="https://github.com/chalkan3-sloth/sloth-runner" class="community-card">
    <div class="icon">🐙</div>
    <h4>GitHub</h4>
    <p>Source code, issues, and contributions</p>
  </a>
  
  <a href="https://github.com/chalkan3-sloth/sloth-runner/discussions" class="community-card">
    <div class="icon">💬</div>
    <h4>Discussions</h4>
    <p>Community Q&A and feature discussions</p>
  </a>
  
  <a href="https://github.com/chalkan3-sloth/sloth-runner/issues" class="community-card">
    <div class="icon">🐛</div>
    <h4>Issues</h4>
    <p>Bug reports and feature requests</p>
  </a>
  
  <a href="mailto:enterprise@sloth-runner.dev" class="community-card">
    <div class="icon">🏢</div>
    <h4>Enterprise</h4>
    <p>Commercial support and services</p>
  </a>
</div>

---

<div class="footer-cta">
  <h3>🦥 Ready to streamline your automation?</h3>
  <p>Join developers using Sloth Runner for reliable, scalable task orchestration.</p>
  <a href="TUTORIAL/" class="btn primary large">🚀 Start Building Today</a>
</div>