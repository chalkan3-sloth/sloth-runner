[English](./README.md) | [Português](./README.pt.md) | [中文](./README.zh.md)

# 🦥 Sloth Runner 🚀

A next-generation flexible and extensible task runner application written in Go, powered by **Modern DSL** Lua scripting. `sloth-runner` allows you to define complex workflows, manage task dependencies, and integrate with external systems through an intuitive **Modern DSL**.

[![Go CI](https://github.com/chalkan3/sloth-runner/actions/workflows/ci.yml/badge.svg)](https://github.com/chalkan3/sloth-runner/actions/workflows/ci.yml)

---

## ✨ Key Features

### 🎯 **Modern DSL (Domain Specific Language) - Only**
*   **🔮 Fluent API:** Define tasks using intuitive, chainable methods
*   **📋 Workflow Definition:** Declarative workflow configuration with metadata
*   **🔄 Enhanced Features:** Built-in retry strategies, circuit breakers, and advanced patterns
*   **🛡️ Type Safety:** Better validation and error detection
*   **📊 Rich Metadata:** Comprehensive task and workflow information
*   **🧹 Clean Syntax:** Single, consistent syntax (legacy format removed)

### 🏗️ **Core Capabilities**
*   **📜 Modern DSL Only:** Clean, powerful syntax without legacy baggage
*   **🔗 Advanced Dependency Management:** Smart dependency resolution with conditional execution
*   **⚡ Enhanced Async Execution:** Parallel task execution with modern async patterns
*   **🪝 Lifecycle Hooks:** Rich pre/post-execution hooks with enhanced error handling
*   **🎯 Circuit Breakers:** Built-in resilience patterns for external dependencies

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

### 🏢 **Enterprise Features**
*   **🌐 Distributed Architecture:** Master-agent with enhanced load balancing
*   **💾 Advanced State Management:** SQLite-based with clustering and replication
*   **🛡️ Enterprise Reliability:** Circuit breakers, saga patterns, and failure handling
*   **📊 Comprehensive Monitoring:** Metrics, health checks, and observability
*   **⏰ Smart Scheduler:** Cron-based with dependency-aware scheduling
*   **📦 Artifact Management:** Versioned artifacts with metadata and retention policies
*   **🔐 Security:** RBAC, secrets management, and audit logging

### 💻 **Modern CLI Interface**
*   `run`: Execute workflows with Modern DSL support
*   `validate`: Enhanced validation for both DSL formats
*   `migrate`: Convert legacy scripts to Modern DSL
*   `list`: Display workflows with enhanced metadata
*   `test`: Advanced testing framework for workflows
*   `repl`: Interactive REPL with Modern DSL support
*   `template`: Modern DSL templates and scaffolding
*   `agent`: Enhanced distributed agent management

---

## 🚀 Quick Start with Modern DSL

### 1. **Installation**
```bash
# Download latest release
curl -sSL https://raw.githubusercontent.com/chalkan3/sloth-runner/main/install.sh | bash

# Or build from source
go install github.com/chalkan3/sloth-runner/cmd/sloth-runner@latest
```

### 2. **Create Your First Modern Workflow**
```lua
-- hello-world-modern.lua
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
./sloth-runner run -f hello-world-modern.lua
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

### 🟢 **Beginner Examples** - 100% Modern DSL
```bash
# Hello World with Modern DSL
./sloth-runner run -f examples/beginner/hello-world.lua

# Simple state management
./sloth-runner run -f examples/simple_state_test.lua

# Basic exec module testing
./sloth-runner run -f examples/exec_test.lua

# Simple pipeline processing
./sloth-runner run -f examples/basic_pipeline.lua
```

### 🟡 **Intermediate Examples**  
```bash
# Parallel execution with modern async
./sloth-runner run -f examples/parallel_execution.lua

# Conditional execution and logic
./sloth-runner run -f examples/conditional_execution.lua

# Enhanced pipeline with modern features
./sloth-runner run -f examples/basic_pipeline_modern.lua

# Terraform infrastructure management
./sloth-runner run -f examples/terraform_example.lua
```

### 🔴 **Advanced Examples**
```bash
# Advanced state management
./sloth-runner run -f examples/state_management_demo.lua

# Enterprise reliability patterns
./sloth-runner run -f examples/reliability_demo.lua
```

### 🌍 **Real-World Examples**
```bash
# Complete CI/CD pipeline
./sloth-runner run -f examples/real-world/nodejs-cicd.lua

# Microservices deployment
./sloth-runner run -f examples/real-world/microservices-deploy.lua
```

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](./CONTRIBUTING.md) for:
- Modern DSL development guidelines
- Code standards and testing
- Documentation improvements
- Example contributions

---

## 🎉 What's New in Modern DSL

- **🎯 Fluent API**: Intuitive, chainable task definitions
- **📋 Workflow Metadata**: Rich workflow information and versioning
- **🔄 Enhanced Retry**: Built-in exponential backoff and circuit breakers
- **⚡ Async Patterns**: Modern parallel execution and timeouts
- **📊 Monitoring**: Integrated performance metrics and observability
- **🛡️ Type Safety**: Better validation and error prevention
- **🎨 Better UX**: More readable and maintainable code
- **🧹 Clean Syntax**: Single, powerful syntax without legacy overhead

Join us in the next generation of workflow automation! 🚀

## 🚀 Getting Started

### Installation

To install `sloth-runner` on your system, you can use the provided `install.sh` script. This script automatically detects your operating system and architecture, downloads the latest release from GitHub, and places the `sloth-runner` executable in `/usr/local/bin`.

```bash
bash <(curl -sL https://raw.githubusercontent.com/chalkan3/sloth-runner/master/install.sh)
```

**Note:** The `install.sh` script requires `sudo` privileges to move the executable to `/usr/local/bin`.

### Basic Usage

To run a Lua task file:

```bash
sloth-runner run -f examples/basic_pipeline.lua
```

To list the tasks in a file:

```bash
sloth-runner list -f examples/basic_pipeline.lua
```

---

## 📜 Defining Tasks in Modern DSL

Tasks are defined using the Modern DSL builder pattern within workflows. Each task uses the fluent API with methods like `:description()`, `:command()`, `:timeout()`, `:depends_on()`, etc.

Example (`examples/basic_pipeline.lua`):

```lua
-- Import reusable tasks from another file. The path is relative.
local docker_tasks = import("examples/shared/docker.lua")

-- Task 1: Fetch data with modern async execution
local fetch_data = task("fetch_data")
    :description("Fetches raw data from an API")
    :command(function(params)
        log.info("🔄 Fetching data...")
        -- Simulates an API call with circuit breaker protection
        return circuit.protect("external_api", function()
            return true, "echo 'Fetched raw data'", { raw_data = "api_data" }
        end)
    end)
    :async(true)
    :timeout("30s")
    :retries(3, "exponential")
    :build()

-- Task 2: Flaky task with enhanced retry logic
local flaky_task = task("flaky_task")
    :description("This task fails intermittently and will retry")
    :command(function()
        if math.random() > 0.5 then
            log.info("✅ Flaky task succeeded.")
            return true, "echo 'Success!'"
        else
            log.error("❌ Flaky task failed, will retry...")
            return false, "Random failure"
        end
    end)
    :retries(3, "exponential")
    :on_failure(function(params, error)
        log.warn("Retry attempt failed: " .. error)
    end)
    :build()

-- Task 3: Process data with dependency injection
local process_data = task("process_data")
    :description("Processes the fetched data")
    :depends_on({"fetch_data", "flaky_task"})
    :command(function(params, deps)
        local raw_data = deps.fetch_data.raw_data
        log.info("🔧 Processing data: " .. raw_data)
        return true, "echo 'Processed data'", { 
            processed_data = "processed_" .. raw_data 
        }
    end)
    :build()

-- Task 4: Long-running task with timeout
local long_running_task = task("long_running_task")
    :description("A task that will be terminated if it runs too long")
    :command("echo 'Starting long task...'; sleep 10; echo 'This will not be printed.';")
    :timeout("5s")
    :on_timeout(function()
        log.warn("⏰ Task timed out as expected")
    end)
    :build()

-- Task 5: Cleanup task with conditional execution
local cleanup_on_fail = task("cleanup_on_fail")
    :description("Runs only if the long-running task fails")
    :depends_on({"long_running_task"})
    :run_if(function(params, deps)
        return not deps.long_running_task.success
    end)
    :command("echo 'Cleanup task executed due to previous failure.'")
    :build()

-- Task 6: Reusable task usage
local build_image = task("build_image")
    :description("Builds the application's Docker image")
    :command(function()
        return docker_tasks.build({
            image_name = "my-awesome-app",
            tag = "v1.2.3",
            context = "./app_context"
        })
    end)
    :artifacts({"docker-image"})
    :build()

-- Task 7: Conditional deployment
local conditional_deploy = task("conditional_deploy")
    :description("Deploys the application only if the build artifact exists")
    :depends_on({"build_image"})
    :run_if("test -f ./app_context/artifact.txt")
    :command("echo 'Deploying application...'")
    :build()

-- Task 8: Gatekeeper with abort condition
local gatekeeper_check = task("gatekeeper_check")
    :description("Aborts the workflow if a critical condition is not met")
    :command("echo 'This command will not be executed if aborted.'")
    :abort_if(function(params, deps)
        -- Lua function condition
        log.warn("🛡️  Checking gatekeeper condition...")
        if params.force_proceed ~= "true" then
            log.error("❌ Gatekeeper check failed. Aborting workflow.")
            return true -- Abort
        end
        return false -- Do not abort
    end)
    :build()

-- Define the complete workflow
workflow.define("full_pipeline_demo", {
    description = "A comprehensive pipeline demonstrating various features - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"demo", "comprehensive", "modern-dsl"},
        complexity = "advanced"
    },
    
    tasks = {
        fetch_data,
        flaky_task,
        process_data,
        long_running_task,
        cleanup_on_fail,
        build_image,
        conditional_deploy,
        gatekeeper_check
    },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 3,
        fail_fast = false
    },
    
    on_start = function()
        log.info("🚀 Starting comprehensive pipeline demo...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Comprehensive pipeline completed successfully!")
        else
            log.error("❌ Pipeline failed - check individual task results")
        end
        return true
    end
})
```
```

---

## 🌐 Master-Agent Architecture

`sloth-runner` is designed with a master-agent architecture to facilitate distributed task execution. This allows you to orchestrate and run tasks across multiple remote machines from a central control point.

### Core Concepts

*   **Master Server:** The central component that manages agents and orchestrates tasks.
*   **Agent:** A lightweight process that runs on a remote machine, executes tasks, and reports status.
*   **Communication:** Master and agents communicate using gRPC.

### Installation and Startup

#### 1. Start the Master Server

On your local machine or a designated control server, start the `sloth-runner` master:

```bash
go run ./cmd/sloth-runner master -p 50053 --daemon
```

#### 2. Start the Agents

On each remote machine, start the `sloth-runner` agent, ensuring you provide a unique name and the correct master address:

```bash
sloth-runner agent start --name <agent_name> --master <master_ip>:<master_port> --port <agent_port> --bind-address <agent_ip> --daemon
```

*   `--name`: A unique name for the agent (e.g., `web-server-1`).
*   `--master`: The address of the master server.
*   `--port`: The port for the agent to listen on.
*   `--bind-address`: The agent's accessible IPv4 address.
*   `--daemon`: Runs the agent as a background process.

### Agent Management

#### Listing Agents

Use `sloth-runner agent list` to view all registered agents and their status:

```bash
go run ./cmd/sloth-runner agent list
```

**Example Output:**

```
AGENT NAME     ADDRESS              STATUS            LAST HEARTBEAT
------------   ----------           ------            --------------
agent2         192.168.1.17:50052   Active   1758984185
agent1         192.168.1.16:50051   Active   1758984183
```

#### Running Commands on Agents

Execute commands on a specific agent using `sloth-runner agent run`:

```bash
go run ./cmd/sloth-runner agent run agent1 'echo "Hello from agent1 on $(hostname)"'
```

**Example Output:**

```
┌─  SUCCESS  Command Execution Result on agent1 ──────────┐
|  SUCCESS  Command executed successfully!                |
|  INFO  Command: echo "Hello from agent1 on $(hostname)" |
| # Stdout:                                               |
| Hello from agent1 on ladyguica                          |
|                                                         |
|                                                         |
└─────────────────────────────────────────────────────────┘
```

#### Stopping Agents

Gracefully stop a remote agent using `sloth-runner agent stop`:

```bash
go run ./cmd/sloth-runner agent stop <agent_name>
```

---

## 📄 Templates

`sloth-runner` provides several templates to quickly scaffold new task definition files.

| Template Name      | Description                                                                    |
| :----------------- | :----------------------------------------------------------------------------- |
| `simple`           | Generates a single group with a 'hello world' task. Ideal for getting started. |
| `python`           | Creates a pipeline to set up a Python environment, install dependencies, and run a script. |
| `parallel`         | Demonstrates how to run multiple tasks concurrently.                           |
| `python-pulumi`    | Pipeline to deploy Pulumi infrastructure managed with Python.                  |
| `python-pulumi-salt` | Provisions infrastructure with Pulumi and configures it using SaltStack.       |
| `git-python-pulumi` | CI/CD Pipeline: Clones a repo, sets up the environment, and deploys with Pulumi. |
| `dummy`            | Generates a dummy task that does nothing.                                      |

---

## 💻 CLI Commands

`sloth-runner` provides a simple and powerful command-line interface.

### `sloth-runner run`

Executes tasks defined in a Lua template file.

**Usage:** `sloth-runner run [flags]`

**Description:**
The run command executes tasks defined in a Lua template file.
You can specify the file, environment variables, and target specific tasks or groups.

**Flags:**

*   `-f, --file string`: Path to the Lua task configuration template file (default: "examples/basic_pipeline.lua")
*   `-e, --env string`: Environment for the tasks (e.g., Development, Production) (default: "Development")
*   `-p, --prod`: Set to true for production environment (default: false)
*   `--shards string`: Comma-separated list of shard numbers (e.g., 1,2,3) (default: "1,2,3")
*   `-t, --tasks string`: Comma-separated list of specific tasks to run (e.g., task1,task2)
*   `-g, --group string`: Run tasks only from a specific task group
*   `-v, --values string`: Path to a YAML file with values to be passed to Lua tasks
*   `-d, --dry-run`: Simulate the execution of tasks without actually running them (default: false)
*   `--return`: Return the output of the target tasks as JSON (default: false)
*   `-y, --yes`: Bypass interactive task selection and run all tasks (default: false)

### `sloth-runner list`

Lists all available task groups and tasks.

**Usage:** `sloth-runner list [flags]`

**Description:**
The list command displays all task groups and their respective tasks, along with their descriptions and dependencies.

**Flags:**

*   `-f, --file string`: Path to the Lua task configuration template file (default: "examples/basic_pipeline.lua")
*   `-e, --env string`: Environment for the tasks (e.g., Development, Production) (default: "Development")
*   `-p, --prod`: Set to true for production environment (default: false)
*   `--shards string`: Comma-separated list of shard numbers (e.g., 1,2,3) (default: "1,2,3")
*   `-v, --values string`: Path to a YAML file with values to be passed to Lua tasks

### `sloth-runner validate`

Validates the syntax and structure of a Lua task file.

**Usage:** `sloth-runner validate [flags]`

**Description:**
The validate command checks a Lua task file for syntax errors and ensures that the Modern DSL workflow structure is correctly formatted.

**Flags:**

*   `-f, --file string`: Path to the Lua task configuration template file (default: "examples/basic_pipeline.lua")
*   `-e, --env string`: Environment for the tasks (e.g., Development, Production) (default: "Development")
*   `-p, --prod`: Set to true for production environment (default: false)
*   `--shards string`: Comma-separated list of shard numbers (e.g., 1,2,3) (default: "1,2,3")
*   `-v, --values string`: Path to a YAML file with values to be passed to Lua tasks

### `sloth-runner test`

Executes a Lua test file for a task workflow.

**Usage:** `sloth-runner test -w <workflow-file> -f <test-file>`

**Description:**
The test command runs a specified Lua test file against a workflow.
Inside the test file, you can use the 'test' and 'assert' modules to validate task behaviors.

**Flags:**

*   `-f, --file string`: Path to the Lua test file (required)
*   `-w, --workflow string`: Path to the Lua workflow file to be tested (required)

### `sloth-runner repl`

Starts an interactive REPL session.

**Usage:** `sloth-runner repl [flags]`

**Description:**
The repl command starts an interactive Read-Eval-Print Loop that allows you
to execute Lua code and interact with all the built-in sloth-runner modules.
You can optionally load a workflow file to have its context available.

**Flags:**

*   `-f, --file string`: Path to a Lua workflow file to load into the REPL session

### `sloth-runner scheduler`

Manages the background task scheduler.

**Usage:** `sloth-runner scheduler [command]`

**Description:**
The scheduler command provides subcommands to enable, disable, list, and delete the sloth-runner background task scheduler.

#### `sloth-runner scheduler enable`

Starts the sloth-runner scheduler in the background.

**Usage:** `sloth-runner scheduler enable [flags]`

**Description:**
The enable command starts the sloth-runner scheduler as a persistent background process.
It will monitor and execute tasks defined in the scheduler configuration file.

**Flags:**

*   `-c, --scheduler-config string`: Path to the scheduler configuration file (default: "scheduler.yaml")

#### `sloth-runner scheduler disable`

Stops the running sloth-runner scheduler.

**Usage:** `sloth-runner scheduler disable`

**Description:**
The disable command stops the background sloth-runner scheduler process.

#### `sloth-runner scheduler list`

Lists all scheduled tasks.

**Usage:** `sloth-runner scheduler list [flags]`

**Description:**
The list command displays all scheduled tasks defined in the scheduler configuration file.

**Flags:**

*   `-c, --scheduler-config string`: Path to the scheduler configuration file (default: "scheduler.yaml")

#### `sloth-runner scheduler delete <task_name>`

Deletes a specific scheduled task.

**Usage:** `sloth-runner scheduler delete <task_name> [flags]`

**Description:**
The delete command removes a specific scheduled task from the scheduler configuration file.

**Arguments:**

*   `<task_name>`: The name of the scheduled task to delete.

**Flags:**

*   `-c, --scheduler-config string`: Path to the scheduler configuration file (default: "scheduler.yaml")

### Agent Commands

*   `sloth-runner agent start [-p <port>]`: Starts the sloth-runner in agent mode.
*   `sloth-runner agent list`: Lists all registered agents.
*   `sloth-runner agent run <agent_name> <command>`: Executes a command on a remote agent.
*   `sloth-runner agent stop <agent_name>`: Stops a remote agent.

### `sloth-runner version`

Print the version number of sloth-runner.

**Usage:** `sloth-runner version`

**Description:**
All software has versions. This is sloth-runner's

### `sloth-runner template list`

Lists all available templates.

**Usage:** `sloth-runner template list`

**Description:**
Displays a table of all available templates that can be used with the 'new' command.

### `sloth-runner new <group-name>`

Generates a new task definition file from a template.

**Usage:** `sloth-runner new <group-name> [flags]`

**Description:**
The new command creates a boilerplate Lua task definition file.
You can choose from different templates and specify an output file.
Run 'sloth-runner template list' to see all available templates.

**Arguments:**

*   `<group-name>`: The name of the task group to generate.

**Flags:**

*   `-o, --output string`: Output file path (default: stdout)
*   `-t, --template string`: Template to use. See `template list` for options. (default: "simple")
*   `--set key=value`: Pass key-value pairs to the template for dynamic content generation.

### `sloth-runner check dependencies`

Checks for required external CLI tools.

**Usage:** `sloth-runner check dependencies`

**Description:**
Verifies that all external command-line tools used by the various modules (e.g., docker, aws, doctl) are installed and available in the system's PATH.

---

## 🚀 Advanced Features

`sloth-runner` also includes several advanced features for more complex workflows and development scenarios.

*   **Interactive Task Runner:** Step through tasks one by one for debugging and development.
*   **Enhanced `values.yaml` Templating:** Use Go template syntax to inject environment variables into your `values.yaml` files for dynamic configurations.

For detailed information on these and other advanced features, refer to the [Advanced Features documentation](./docs/advanced-features.md).

---

## ⚙️ Lua API

`sloth-runner` exposes several Go functionalities as Lua modules, allowing your tasks to interact with the system and external services.

*   **`exec` module:** Execute shell commands.
*   **`fs` module:** Perform file system operations.
*   **`net` module:** Make HTTP requests and download files.
*   **`data` module:** Parse and serialize data in JSON and YAML format.
*   **`log` module:** Log messages with different severity levels.
*   **`salt` module:** Execute SaltStack commands.

For detailed API usage, please refer to the examples in the `/examples` directory.
# Test pipeline fix
# Test pipeline fix
