[English](./README.md) | [Português](./README.pt.md) | [中文](./README.zh.md)

# 🦥 Sloth Runner

A **modern task orchestration platform** built with Go and powered by **Lua scripting**. Sloth Runner provides a fluent Modern DSL for defining complex workflows, distributed execution capabilities, and comprehensive automation tools for DevOps teams.

**Sloth Runner** simplifies task automation with its intuitive Lua-based DSL, distributed master-agent architecture, and extensive built-in modules for common DevOps operations.

[![Go CI](https://github.com/chalkan3-sloth/sloth-runner/actions/workflows/ci.yml/badge.svg)](https://github.com/chalkan3-sloth/sloth-runner/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Lua Powered](https://img.shields.io/badge/Lua-Powered-purple.svg)](https://www.lua.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](./LICENSE)

---

## ✨ **Key Features**

### 🎯 **Modern DSL (Domain Specific Language)**
*Clean, powerful Lua-based syntax for complex workflows*

```lua
-- Define GitOps tasks with fluent API
local clone_task = task("clone_infrastructure")
    :description("Clone Terraform infrastructure repository")
    :workdir("/tmp/infrastructure")
    :command(function(this, params)
        local git = require("git")
        log.info("📡 Cloning infrastructure repository...")
        
        local repository = git.clone(
            "https://github.com/company/terraform-infrastructure",
            this.workdir.get()
        )
        
        return true, "Repository cloned successfully"
    end)
    :timeout("5m")
    :retries(3, "exponential")
    :on_success(function(this, params, output)
        log.info("✅ Infrastructure code ready for deployment!")
    end)
    :build()

-- Define workflows with advanced configuration
workflow.define("infrastructure_pipeline")
    :description("Complete Infrastructure Deployment Pipeline")
    :version("2.0.0")
    :tasks({ clone_task, terraform_plan_task, terraform_apply_task })
    :config({
        timeout = "30m",
        max_parallel_tasks = 2,
        environment = "production"
    })
    :on_complete(function(success, results)
        if success then
            log.info("🎉 Infrastructure successfully deployed!")
        end
    end)
```

### 🌐 **Distributed Master-Agent Architecture**
*Scale task execution across multiple machines*

- **Master Server:** Central orchestration and control
- **Agents:** Lightweight workers on remote machines  
- **gRPC Communication:** Reliable, high-performance communication
- **Load Balancing:** Intelligent task distribution
- **Health Monitoring:** Real-time agent status tracking

### 💾 **Advanced State Management**
*Persistent state with SQLite backend and advanced features*

```lua
-- Persistent state operations
state.set("deployment_count", 1, 3600) -- TTL of 1 hour
local count = state.increment("deployment_count")

-- Atomic operations
state.compare_swap("status", "deploying", "deployed")

-- Distributed locks
state.with_lock("deployment_lock", function()
    -- Critical section - only one task can execute this
    return deploy_application()
end, 30) -- 30 second timeout
```

### 🔧 **Rich Lua Module Ecosystem**
*Comprehensive built-in modules for common operations*

- **`git`**: Git operations with automatic credential handling
- **`terraform`**: Terraform lifecycle management (init, plan, apply)
- **`exec`**: Execute shell commands with enhanced error handling
- **`fs`**: File system operations with validation
- **`net`**: HTTP client with retry and timeout support
- **`data`**: JSON/YAML processing and validation
- **`log`**: Structured logging with context
- **`state`**: Persistent state management
- **`async`**: Parallel execution and modern async patterns
- **`utils`**: Configuration management and utilities

#### GitOps Module Examples

```lua
-- Git operations
local git = require("git")
local repo = git.clone("https://github.com/company/infra", "/tmp/infra")
git.checkout(repo, "production")

-- Terraform operations  
local terraform = require("terraform")
local client = terraform.init("/tmp/infra/terraform/")  -- Runs terraform init
local plan = client:plan({ var_file = "prod.tfvars" })
local apply = client:apply({ auto_approve = true })

-- State management
local state = require("state")
state.set("deployment_version", "v2.1.0", 3600)
local version = state.get("deployment_version")
```

## 🚀 **Quick Start**

### Installation

```bash
# Install latest version
curl -sSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/main/install.sh | bash

# Or download from releases
wget https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner-linux-amd64.tar.gz

# Or build from source
go install github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner@latest
```

### Your First GitOps Workflow

Create a complete GitOps workflow that clones a repository and deploys infrastructure with Terraform:

**1. Create your workflow file (`deploy.sloth`):**

```lua
-- Complete GitOps Workflow with Git + Terraform
local clone_repo_task = task("clone_repo")
    :description("Clone Git repository with infrastructure code")
    :workdir("/tmp/infrastructure")
    :command(function(this, params)
        local git = require("git")
        
        log.info("📡 Cloning repository...")
        local repository = git.clone(
            "https://github.com/your-org/terraform-infrastructure",
            this.workdir.get()
        )
        
        return true, "Repository cloned successfully", {
            repository_url = "https://github.com/your-org/terraform-infrastructure",
            clone_destination = this.workdir.get()
        }
    end)
    :timeout("5m")
    :build()

local deploy_infrastructure = task("deploy_terraform")
    :description("Deploy infrastructure using Terraform")
    :workdir("/tmp/infrastructure/environments/prod/")
    :command(function(this, params)
        local terraform = require("terraform")
        
        -- Terraform init is called automatically
        log.info("🔄 Initializing Terraform...")
        local client = terraform.init(this.workdir:get())
        
        -- Load configuration from values.yaml
        local terraform_config = {
            environment = Values.terraform.environment or "prod",
            instance_type = Values.terraform.instance_type or "t3.micro",
            region = Values.terraform.region or "us-east-1"
        }
        
        -- Create terraform.tfvars from configuration
        local tfvars = client:create_tfvars("terraform.tfvars", terraform_config)
        
        -- Plan and apply
        local plan_result = client:plan({ var_file = tfvars.filename })
        if plan_result.success then
            log.info("🚀 Applying Terraform changes...")
            local apply_result = client:apply({
                var_file = tfvars.filename,
                auto_approve = true
            })
            
            return apply_result.success, "Infrastructure deployment", {
                terraform_used = true,
                plan_success = true,
                apply_success = apply_result.success
            }
        end
        
        return false, "Terraform plan failed"
    end)
    :timeout("15m")
    :build()

-- Define the complete GitOps workflow
workflow.define("gitops_deploy")
    :description("Complete GitOps workflow: Git clone + Terraform deploy")
    :version("1.0.0")
    :tasks({ clone_repo_task, deploy_infrastructure })
    :config({
        timeout = "20m",
        max_parallel_tasks = 1
    })
```

**2. Create your configuration file (`values.yaml`):**

```yaml
terraform:
  environment: "production"
  instance_type: "t3.small"
  region: "us-west-2"
  
workflow:
  timeout: "30m"
  environment: "prod"
```

**3. Run your GitOps workflow:**

```bash
# Execute the complete workflow
sloth-runner run -f deploy.sloth -v values.yaml gitops_deploy

# Watch the magic happen:
# ✅ Repository cloned
# ✅ Terraform initialized  
# ✅ Infrastructure planned
# ✅ Infrastructure deployed
```

### Try the Complete Example

We provide a working example that you can run immediately:

```bash
# Clone the repository
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Run the GitOps example
sloth-runner run -f examples/deploy_git_terraform.sloth -v examples/values.yaml deploy_git_terraform
```

This example demonstrates:
- 🔄 **Git repository cloning**
- 🏗️ **Terraform infrastructure deployment**  
- ⚙️ **External configuration with values.yaml**
- 📊 **Comprehensive error handling and logging**
- 🎯 **Modern DSL syntax and best practices**

### Hello World Example

Create your first workflow with the Modern DSL:

```lua
-- hello-world.sloth
local hello_task = task("say_hello")
    :description("Simple hello world demonstration")
    :command(function(params)
        log.info("🌟 Hello World from Sloth Runner!")
        log.info("📅 Current time: " .. os.date())
        
        return true, "echo 'Hello, Modern Sloth Runner!'", {
            message = "Hello World",
            timestamp = os.time(),
            status = "success"
        }
    end)
    :timeout("30s")
    :on_success(function(params, output)
        log.info("✅ Hello World task completed successfully!")
        log.info("💬 Message: " .. output.message)
    end)
    :build()

-- Define the workflow
workflow.define("hello_world_workflow", {
    description = "Simple Hello World demonstration",
    version = "1.0.0",
    tasks = { hello_task },
    
    config = {
        timeout = "5m",
        max_parallel_tasks = 1
    }
})
```

```bash
# Run the workflow
./sloth-runner run -f hello-world.sloth
```

### Basic Pipeline Example

```lua
-- pipeline.sloth
-- Task 1: Fetch data
local fetch_data = task("fetch_data")
    :description("Fetches raw data from an API")
    :command(function(params)
        log.info("🔄 Fetching data...")
        return true, "echo 'Fetched raw data'", { 
            raw_data = "api_data",
            source = "external_api" 
        }
    end)
    :timeout("30s")
    :build()

-- Task 2: Process data (depends on fetch_data)
local process_data = task("process_data")
    :description("Processes the fetched data")
    :depends_on({"fetch_data"})
    :command(function(params, deps)
        local raw_data = deps.fetch_data.raw_data
        log.info("🔧 Processing data: " .. raw_data)
        
        return true, "echo 'Processed data'", { 
            processed_data = "processed_" .. raw_data 
        }
    end)
    :build()

-- Task 3: Store result
local store_result = task("store_result")
    :description("Stores the processed data")
    :depends_on({"process_data"})
    :command(function(params, deps)
        local final_data = deps.process_data.processed_data
        log.info("💾 Storing result: " .. final_data)
        
        -- Store in state for persistence
        state.set("last_result", final_data, 3600) -- 1 hour TTL
        
        return true, "echo 'Data stored successfully'"
    end)
    :build()

-- Define the complete pipeline
workflow.define("data_pipeline", {
    description = "Data Processing Pipeline",
    version = "1.0.0",
    tasks = { fetch_data, process_data, store_result },
    
    config = {
        timeout = "10m",
        max_parallel_tasks = 2
    }
})
```

### 🏢 **Enterprise & Production Features**
*Production-ready capabilities for enterprise environments*

- **🔒 Security:** RBAC, secrets management, and audit logging
- **📊 Monitoring:** Metrics, health checks, and observability  
- **⏰ Scheduler:** Cron-based with dependency-aware scheduling
- **📦 Artifacts:** Versioned artifacts with metadata and retention
- **🛡️ Reliability:** Circuit breakers and failure handling patterns
- **🌐 Clustering:** Master-agent architecture with load balancing

### 💻 **Modern CLI Interface**
*Comprehensive command-line interface for all operations*

```bash
# Core commands
sloth-runner run -f workflow.sloth        # Execute workflows
sloth-runner run --interactive          # Interactive task selection
sloth-runner ui                         # Start web dashboard

# 🆔 NEW: Stack Management (Pulumi-style)
sloth-runner stack new my-app --description "My application stack"  # Create stack
sloth-runner run my-stack -f workflow.sloth --output enhanced  # Run with stack
sloth-runner run my-stack -f workflow.sloth --output json     # JSON output
sloth-runner stack list                                      # List all stacks  
sloth-runner stack show my-stack                            # Show stack details
sloth-runner list -f workflow.sloth                          # List tasks with IDs

# Distributed execution
sloth-runner master                     # Start master server
sloth-runner agent start --name agent1  # Start agent
sloth-runner agent list                 # List all agents
sloth-runner agent run agent1 "command" # Execute on specific agent

# Utilities
sloth-runner scheduler enable           # Enable task scheduler
sloth-runner scheduler list             # List scheduled tasks
sloth-runner version                    # Show version information
```

## 🌟 **Advanced Examples**

### Complete CI/CD Pipeline

```lua
-- ci-cd-pipeline.sloth
local test_task = task("run_tests")
    :description("Run application tests")
    :command(function(params, deps)
        log.info("🧪 Running tests...")
        local result = exec.run("go test ./...")
        if not result.success then
            return false, "Tests failed: " .. result.stderr
        end
        return true, result.stdout, { test_results = "passed" }
    end)
    :timeout("10m")
    :build()

local build_task = task("build_app")
    :description("Build application")
    :depends_on({"run_tests"})
    :command(function(params, deps)
        log.info("🔨 Building application...")
        return exec.run("go build -o app ./cmd/main.go")
    end)
    :artifacts({"app"})
    :build()

local docker_task = task("build_docker")
    :description("Build Docker image")
    :depends_on({"build_app"})
    :command(function(params, deps)
        local tag = params.image_tag or "latest"
        log.info("🐳 Building Docker image with tag: " .. tag)
        
        local result = exec.run("docker build -t myapp:" .. tag .. " .")
        if result.success then
            state.set("docker_image", "myapp:" .. tag, 86400) -- 24 hours
        end
        return result.success, result.stdout
    end)
    :build()

local deploy_task = task("deploy_app")
    :description("Deploy to production")
    :depends_on({"build_docker"})
    :run_if(function(params, deps)
        -- Only deploy if all tests passed and image is built
        return deps.run_tests.test_results == "passed" and 
               state.exists("docker_image")
    end)
    :command(function(params, deps)
        local image = state.get("docker_image")
        log.info("🚀 Deploying image: " .. image)
        
        -- Deploy with rollback capability
        local result = exec.run("kubectl set image deployment/myapp app=" .. image)
        if result.success then
            -- Wait for rollout
            exec.run("kubectl rollout status deployment/myapp --timeout=300s")
        end
        return result.success, result.stdout
    end)
    :on_failure(function(params, error)
        log.error("❌ Deployment failed, rolling back...")
        exec.run("kubectl rollout undo deployment/myapp")
    end)
    :build()

workflow.define("ci_cd_pipeline", {
    description = "Complete CI/CD Pipeline",
    version = "1.0.0",
    
    metadata = {
        author = "DevOps Team",
        tags = {"ci", "cd", "production"}
    },
    
    tasks = { test_task, build_task, docker_task, deploy_task },
    
    config = {
        timeout = "45m",
        max_parallel_tasks = 2,
        fail_fast = true
    },
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 CI/CD Pipeline completed successfully!")
            -- Send notification
            net.post("https://hooks.slack.com/webhook", {
                text = "✅ Deployment successful for commit " .. (os.getenv("GIT_COMMIT") or "unknown")
            })
        else
            log.error("❌ CI/CD Pipeline failed!")
        end
        return true
    end
})
```

### Distributed Task Execution

```lua
-- distributed-workflow.sloth
local web_deploy = task("deploy_web_servers")
    :description("Deploy to web servers")
    :agent_selector("web-server-*")  -- Select agents matching pattern
    :command(function(params, deps)
        log.info("🌐 Deploying to web server: " .. params.agent_name)
        return exec.run("docker pull myapp:latest && docker-compose up -d")
    end)
    :parallel(true)  -- Run on all matching agents in parallel
    :build()

local db_migrate = task("run_db_migrations")
    :description("Run database migrations")
    :agent_selector("db-server-1")  -- Run only on specific agent
    :command(function(params, deps)
        log.info("🗄️ Running database migrations...")
        return exec.run("./migrate.sh --env=production")
    end)
    :build()

local health_check = task("verify_deployment")
    :description("Verify deployment health")
    :depends_on({"deploy_web_servers", "run_db_migrations"})
    :command(function(params, deps)
        log.info("🏥 Checking application health...")
        
        local health_url = "http://loadbalancer:8080/health"
        local response = net.get(health_url, {timeout = 30})
        
        if response.status_code == 200 then
            log.info("✅ Health check passed!")
            return true, "Health check successful"
        else
            log.error("❌ Health check failed!")
            return false, "Health check failed: " .. response.status_code
        end
    end)
    :retries(3, "exponential")
    :build()

workflow.define("distributed_deployment", {
    description = "Distributed Application Deployment",
    version = "1.0.0",
    tasks = { web_deploy, db_migrate, health_check },
    
    config = {
        timeout = "30m",
        require_all_agents = false  -- Continue even if some agents are offline
    }
})
```

## 📚 **Documentation**

- **🚀 [Getting Started](docs/getting-started.md)** - Complete setup and first steps
- **📖 [Modern DSL Reference](docs/LUA_API.md)** - Complete language and API reference  
- **🏗️ [Architecture Guide](docs/distributed.md)** - Master-agent architecture details
- **🧪 [Examples](docs/EXAMPLES.md)** - Real-world usage examples and patterns
- **🔧 [Advanced Features](docs/advanced-features.md)** - Enterprise capabilities
- **📊 [State Management](docs/state.md)** - Persistent state and data handling
- **🛡️ [Security Guide](docs/security.md)** - RBAC, secrets, and audit logging
- **📈 [Monitoring](docs/monitoring.md)** - Metrics, health checks, and observability

### 🛠️ **Development Tools**

- **🦥 [Neovim Plugin](docs/nvim-plugin.md)** - IDE-grade support with syntax highlighting, code completion, and integrated runner
  - [English](docs/en/nvim-plugin.md) | [Português](docs/pt/nvim-plugin.md) | [中文](docs/zh/nvim-plugin.md)

## 🎯 **Why Choose Sloth Runner?**

### 💡 **Developer Experience**
- **📝 Clean, intuitive syntax** with Modern DSL fluent API
- **🧪 Interactive development** with REPL and comprehensive testing
- **📚 Extensive documentation** with real-world examples
- **🔧 Rich ecosystem** of 15+ built-in Lua modules

### 🏢 **Enterprise Ready**
- **🔒 Production-grade security** with RBAC and secrets management
- **📊 Comprehensive monitoring** with metrics and health checks
- **🌐 Distributed architecture** with reliable master-agent topology  
- **⚡ High performance** with parallel execution and state persistence

### 🚀 **Modern Architecture**
- **🎯 Modern DSL only** - no legacy syntax or backwards compatibility issues
- **💾 Advanced state management** with SQLite persistence and TTL
- **🔄 Intelligent retry logic** with exponential backoff and circuit breakers
- **🪝 Rich lifecycle hooks** for comprehensive workflow control

## 🤝 **Community & Support**

- **📖 [Documentation](https://github.com/chalkan3-sloth/sloth-runner/tree/main/docs)** - Comprehensive guides and references
- **🐛 [Issue Tracker](https://github.com/chalkan3-sloth/sloth-runner/issues)** - Report bugs and request features
- **💡 [Discussions](https://github.com/chalkan3-sloth/sloth-runner/discussions)** - Ideas and general discussions
- **🏢 [Enterprise Support](mailto:enterprise@sloth-runner.dev)** - Commercial support and services

## 📈 **Project Status**

### ✅ **Current Features (Stable)**
- ✅ Modern DSL with fluent API
- ✅ Distributed master-agent architecture
- ✅ Advanced state management with SQLite
- ✅ Rich Lua module ecosystem (exec, fs, net, data, log, etc.)
- ✅ Enterprise features (RBAC, monitoring, scheduling)
- ✅ Comprehensive CLI interface
- ✅ Template system and scaffolding tools

### 🚧 **In Development**
- 🔄 Enhanced web UI with real-time monitoring
- 🔄 Additional cloud provider integrations
- 🔄 Advanced workflow visualization
- 🔄 Performance optimizations

### 🔮 **Planned Features**
- 📋 Workflow versioning and rollback
- 🔗 Integration with popular CI/CD platforms
- 📊 Advanced analytics and reporting
- 🎯 Custom plugin system

## 📄 **License**

MIT License - see [LICENSE](LICENSE) file for details.

---

**🦥 Sloth Runner** - *Modern task orchestration made simple*

*Built with ❤️ by the Sloth Runner Team*

### 🌟 **Modern DSL Examples**

#### Fluent Task Definition
```lua
local build_task = task("build_application")
    :description("Build application with modern features")
    :command(function(params, deps)
        log.info("Building application...")
        return exec.run("go build -o app ./cmd/main.go")
    end)
    :timeout("5m")
    :retries(3, "exponential")
    :artifacts({"app"})
    :on_success(function(params, output)
        log.info("Build completed successfully!")
    end)
    :build()
```

#### Workflow Definition
```lua
workflow.define("ci_pipeline", {
    description = "Continuous Integration Pipeline - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "DevOps Team",
        tags = {"ci", "build", "deploy"}
    },
    
    tasks = { build_task, test_task, deploy_task },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 4
    },
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 CI Pipeline completed successfully!")
        end
        return true
    end
})
```

### 🔧 **Enhanced Lua API Modules**
*   **`exec` module:** Enhanced shell command execution with modern error handling
*   **`fs` module:** Advanced file system operations with metadata and validation
*   **`net` module:** HTTP client with retry, timeout, and circuit breaker support
*   **`data` module:** JSON/YAML processing with schema validation
*   **`log` module:** Structured logging with contextual information
*   **`state` module:** Advanced persistent state with TTL, atomic operations, and clustering
*   **`async` module:** Modern async patterns - parallel execution, timeouts, promises
*   **`perf` module:** Performance monitoring and metrics collection
*   **`circuit` module:** Circuit breaker patterns for resilience
*   **`utils` module:** Configuration management, secrets, and utilities
*   **`validate` module:** Input validation and type checking

### 📊 **JSON Output Support**
*   **Structured Output:** `--output json` for machine-readable results
*   **Complete Information:** Status, duration, tasks, stack info, and outputs
*   **Error Handling:** Detailed error information in JSON format
*   **Integration Ready:** Perfect for CI/CD pipelines and monitoring tools

```json
{
  "status": "success",
  "duration": "9.073s",
  "stack": {
    "id": "3ec19a86...",
    "name": "production-deploy"
  },
  "tasks": {
    "build_app": {
      "status": "Success",
      "duration": "2.020s"
    }
  },
  "outputs": {
    "deployment_url": "https://app.example.com"
  }
}
```

### 🏢 **Enterprise Features**
*   **🌐 Distributed Architecture:** Master-agent with enhanced load balancing
*   **💾 Advanced State Management:** SQLite-based with clustering and replication
*   **🛡️ Enterprise Reliability:** Circuit breakers, saga patterns, and failure handling
*   **📊 Comprehensive Monitoring:** Metrics, health checks, and observability
*   **⏰ Smart Scheduler:** Cron-based with dependency-aware scheduling
*   **📦 Artifact Management:** Versioned artifacts with metadata and retention policies
*   **🔐 Security:** RBAC, secrets management, and audit logging

### 💻 **Modern CLI Interface**
*   `run`: Execute workflows with Modern DSL support and multiple output formats
*   `run --output json`: JSON structured output for automation and integration
*   `stack`: Complete stack management (list, show, delete)
*   `list`: Display workflows with enhanced metadata and unique IDs
*   `workflow init`: Project scaffolding with templates
*   `validate`: Enhanced validation for both DSL formats
*   `agent`: Enhanced distributed agent management
*   `ui`: Web-based dashboard for monitoring and management

---

## 🚀 Quick Start with Modern DSL

### 1. **Installation**
```bash
# Download latest release
curl -sSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/main/install.sh | bash

# Or build from source
go install github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner@latest
```

### 2. **Create Your First Modern Workflow**
```lua
-- hello-world-modern.sloth
local hello_task = task("say_hello")
    :description("Modern DSL hello world")
    :command(function()
        log.info("🚀 Hello from Modern DSL!")
        return true, "Hello completed", {
            message = "Welcome to Sloth Runner Modern DSL!",
            timestamp = os.time()
        }
    end)
    :timeout("30s")
    :build()

workflow.define("hello_world", {
    description = "Hello World - Modern DSL",
    version = "1.0.0",
    tasks = { hello_task }
})
```

### 3. **Run Your Workflow**
```bash
./sloth-runner run -f hello-world-modern.sloth
```

---

## 📚 Complete Documentation

### 🎯 **Modern DSL Guide**
- [Modern DSL Introduction](./docs/modern-dsl/introduction.md)
- [Task Definition API](./docs/modern-dsl/task-api.md)  
- [Workflow Definition](./docs/modern-dsl/workflow-api.md)
- [Migration Guide](./docs/modern-dsl/migration-guide.md)
- [Best Practices](./docs/modern-dsl/best-practices.md)

### 📖 **Core Documentation**
- [Getting Started](./docs/getting-started.md)
- [Modern DSL Examples](./examples/README.md)
- [Lua API Reference](./docs/LUA_API.md)
- [Enterprise Features](./docs/enterprise.md)
- [Distributed Architecture](./docs/distributed.md)

### 🔧 **Advanced Topics**
- [Circuit Breakers & Resilience](./docs/resilience.md)
- [Performance Monitoring](./docs/monitoring.md)
- [State Management](./docs/state.md)
- [Security & RBAC](./docs/security.md)

---

## 🎯 Modern DSL Benefits

| Feature | Legacy Format | Modern DSL |
|---------|---------------|------------|
| **Syntax** | Procedural | Fluent, chainable |
| **Type Safety** | Runtime errors | Compile-time validation |
| **Error Handling** | Basic | Enhanced with context |
| **Metadata** | Limited | Rich, structured |
| **Retry Logic** | Manual | Built-in strategies |
| **Dependencies** | Simple | Advanced with conditions |
| **Monitoring** | Basic logging | Comprehensive metrics |
| **Testing** | Manual | Integrated test framework |
| **Learning Curve** | ~~Two syntaxes~~ | Single, intuitive syntax |

---

## 🌟 **Modern DSL Examples - Complete Collection**

## 📁 **Example Workflows**

### 🟢 **Beginner Examples**
```bash
# Hello World with Modern DSL
./sloth-runner run -f examples/beginner/hello-world.sloth

# Simple state management
./sloth-runner run -f examples/simple_state_test.sloth

# Basic exec module testing
./sloth-runner run -f examples/exec_test.sloth

# Simple pipeline processing
./sloth-runner run -f examples/basic_pipeline.sloth
```

### 🟡 **Intermediate Examples**  
```bash
# Parallel execution with modern async
./sloth-runner run -f examples/parallel_execution.sloth

# Conditional execution and logic
./sloth-runner run -f examples/conditional_execution.sloth

# Enhanced pipeline with modern features
./sloth-runner run -f examples/basic_pipeline_modern.sloth

# Terraform infrastructure management
./sloth-runner run -f examples/terraform_example.sloth
```

### 🔴 **Advanced Examples**
```bash
# Advanced state management
./sloth-runner run -f examples/state_management_demo.sloth

# Enterprise reliability patterns
./sloth-runner run -f examples/reliability_demo.sloth
```

### 🌍 **Real-World Examples**
```bash
# Complete CI/CD pipeline
./sloth-runner run -f examples/real-world/nodejs-cicd.sloth

# Microservices deployment
./sloth-runner run -f examples/real-world/microservices-deploy.sloth
```

---

## 🤝 **Contributing**

We welcome contributions! Please see our [Contributing Guide](./CONTRIBUTING.md) for:
- Modern DSL development guidelines
- Code standards and testing
- Documentation improvements
- Example contributions

---


