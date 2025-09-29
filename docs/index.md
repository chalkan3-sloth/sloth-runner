# 🦥 Sloth Runner - Advanced Task Orchestration Platform

<div class="hero">
  <h1>🚀 Intelligent Task Orchestration</h1>
  <p>A powerful, modern task runner with distributed execution, state management, and comprehensive monitoring capabilities.</p>
  <div class="hero-badges">
    <span class="badge distributed">🌐 Distributed</span>
    <span class="badge state">💾 Stateful</span>
    <span class="badge monitoring">📊 Observable</span>
    <span class="badge enterprise">🏢 Enterprise Ready</span>
  </div>
  <a href="TUTORIAL/" class="btn primary">🚀 Get Started</a>
  <a href="en/quick-start/" class="btn secondary">⚡ Quick Start</a>
</div>

## 🌟 **Core Features**

<div class="feature-grid">
  <div class="feature-card">
    <div class="icon">🗂️</div>
    <h3>Stack Management</h3>
    <p>Pulumi-style stack management with persistent state, exported outputs, and execution history tracking.</p>
    <ul>
      <li>🔒 Persistent stack state with SQLite</li>
      <li>📊 Exported outputs capture</li>
      <li>📈 Complete execution history</li>
      <li>🎯 Environment isolation</li>
    </ul>
  </div>

  <div class="feature-card">
    <div class="icon">🌐</div>
    <h3>Distributed by Design</h3>
    <p>Native master-agent architecture with real-time streaming, automatic failover, and intelligent load balancing.</p>
    <ul>
      <li>🔗 gRPC-based agent communication</li>
      <li>📡 Real-time command streaming</li>
      <li>🔄 Automatic failover and recovery</li>
      <li>⚖️ Intelligent load balancing</li>
    </ul>
  </div>
  
  <div class="feature-card">
    <div class="icon">💾</div>
    <h3>Advanced State Management</h3>
    <p>Built-in SQLite-based persistent state with atomic operations, distributed locks, and TTL support.</p>
    <ul>
      <li>🔒 Distributed locking mechanisms</li>
      <li>⚛️ Atomic operations support</li>
      <li>⏰ TTL-based data expiration</li>
      <li>🔍 Pattern-based queries</li>
    </ul>
  </div>
  
  <div class="feature-card">
    <div class="icon">📊</div>
    <h3>Enhanced Output System</h3>
    <p>Pulumi-style rich output formatting with configurable styles, progress indicators, and structured displays.</p>
    <ul>
      <li>🎨 Multiple output styles (basic, enhanced, rich, modern)</li>
      <li>📈 Real-time progress indicators</li>
      <li>🎯 Structured output sections</li>
      <li>🌈 Rich color formatting</li>
    </ul>
  </div>
      <li>🚨 Real-time alerting</li>
    </ul>
  </div>

  <div class="feature-card">
    <div class="icon">🔐</div>
    <h3>Production Security</h3>
    <p>mTLS authentication, RBAC authorization, audit logging, and compliance-ready security features.</p>
    <ul>
      <li>🔐 mTLS mutual authentication</li>
      <li>👤 Role-based access control</li>
      <li>📋 Comprehensive audit logging</li>
      <li>✅ Compliance-ready features</li>
    </ul>
  </div>
  
  <div class="feature-card">
    <div class="icon">☁️</div>
    <h3>Multi-Cloud Ready</h3>
    <p>Native integrations with AWS, GCP, Azure, and on-premises infrastructure for hybrid deployments.</p>
    <ul>
      <li>☁️ AWS, GCP, Azure support</li>
      <li>🏢 On-premises integration</li>
      <li>🔄 Hybrid cloud orchestration</li>
      <li>🌊 DigitalOcean support</li>
    </ul>
  </div>
  
  <div class="feature-card">
    <div class="icon">⚡</div>
    <h3>High Performance</h3>
    <p>Parallel execution, circuit breakers, and intelligent resource optimization for maximum throughput.</p>
    <ul>
      <li>🔄 Parallel task execution</li>
      <li>⚡ Circuit breaker patterns</li>
      <li>🎯 Resource optimization</li>
      <li>📊 Performance analytics</li>
    </ul>
  </div>
</div>

## 🚀 **Practical Examples**

### 🗂️ Stack Management with Pulumi-Style Output

```bash
# Create a new project from template
sloth-runner workflow init my-cicd --template cicd

# Deploy to development environment
sloth-runner run dev-app -f my-cicd.lua --output enhanced

# Deploy to production with stack persistence
sloth-runner run prod-app -f my-cicd.lua -o rich

# Check deployment status and outputs
sloth-runner stack show prod-app
```

### 📊 Stack with Exported Outputs

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
</div>

## 🚀 **Get Started in Minutes**

<div class="getting-started">
  <div class="step">
    <div class="step-number">1</div>
    <h4>Install</h4>
    <pre><code># Linux/macOS
curl -L https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m | sed 's/x86_64/amd64/').tar.gz | tar xz
chmod +x sloth-runner && sudo mv sloth-runner /usr/local/bin/</code></pre>
  </div>
  
  <div class="step">
    <div class="step-number">2</div>
    <h4>Create Your First Workflow</h4>
    <pre><code>echo 'local hello = task("hello")
  :command(function() 
    log.info("Hello from Sloth Runner! 🦥")
    return true 
  end)
  :build()

workflow.define("greeting", { tasks = { hello } })' > hello.lua</code></pre>
  </div>
  
  <div class="step">
    <div class="step-number">3</div>
    <h4>Run Your Workflow</h4>
    <pre><code>sloth-runner run -f hello.lua</code></pre>
  </div>
</div>

## 📚 **Learn More**

<div class="learn-more-grid">
  <a href="TUTORIAL/" class="learn-card">
    <div class="icon">🚀</div>
    <h4>Quick Tutorial</h4>
    <p>Get up and running with practical examples in 5 minutes</p>
  </a>
  
  <a href="en/advanced-examples/" class="learn-card">
    <div class="icon">📝</div>
    <h4>Advanced Examples</h4>
    <p>Production-ready workflows and real-world use cases</p>
  </a>
  
  <a href="en/core-concepts/" class="learn-card">
    <div class="icon">🧠</div>
    <h4>Core Concepts</h4>
    <p>Understanding tasks, workflows, and distributed execution</p>
  </a>
  
  <a href="en/enterprise-features/" class="learn-card">
    <div class="icon">🏢</div>
    <h4>Enterprise Features</h4>
    <p>Production-grade security, monitoring, and reliability</p>
  </a>
  
  <a href="en/distributed/" class="learn-card">
    <div class="icon">🌐</div>
    <h4>Distributed Execution</h4>
    <p>Master-agent architecture and multi-node coordination</p>
  </a>
  
  <a href="en/modules/" class="learn-card">
    <div class="icon">🔧</div>
    <h4>Module Reference</h4>
    <p>Complete API documentation for all built-in modules</p>
  </a>
</div>

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
- 📖 [Getting Started](en/getting-started.md)
- 🧠 [Core Concepts](en/core-concepts.md)
- ⚡ [Quick Start](en/quick-start.md)
- 💻 [CLI Reference](en/CLI.md)
- 🔄 [Interactive REPL](en/repl.md)
- 🎯 [Advanced Features](en/advanced-features.md)
- 🚀 [Agent Improvements](en/agent-improvements.md)

### 🇧🇷 **Documentação em Português**
- 📖 [Primeiros Passos](pt/getting-started.md)
- 🧠 [Conceitos Fundamentais](pt/core-concepts.md)
- ⚡ [Início Rápido](pt/quick-start.md)
- 💻 [Referência CLI](pt/CLI.md)
- 🔄 [REPL Interativo](pt/repl.md)
- 🎯 [Recursos Avançados](pt/advanced-features.md)
- 🚀 [Melhorias dos Agentes](pt/agent-improvements.md)

### 🇨🇳 **中文文档**
- 📖 [入门指南](zh/getting-started.md)
- 🧠 [核心概念](zh/core-concepts.md)
- ⚡ [快速开始](zh/quick-start.md)
- 💻 [CLI参考](zh/CLI.md)
- 🔄 [交互式REPL](zh/repl.md)
- 🎯 [高级功能](zh/advanced-features.md)
- 🚀 [代理改进](zh/agent-improvements.md)

## 🔧 **Module Reference**

### **📦 Built-in Modules**
| Module | Description | Language Support |
|--------|-------------|------------------|
| [💾 **State**](en/modules/state.md) | Persistent state management | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📊 **Metrics**](en/modules/metrics.md) | Monitoring and observability | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [⚡ **Exec**](en/modules/exec.md) | Command execution | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📁 **FS**](en/modules/fs.md) | File system operations | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📡 **Net**](en/modules/net.md) | Network operations | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📋 **Data**](en/modules/data.md) | Data processing utilities | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |
| [📝 **Log**](en/modules/log.md) | Structured logging | <span class="lang-badge en">EN</span> <span class="lang-badge pt">PT</span> <span class="lang-badge zh">ZH</span> |

### **☁️ Cloud Provider Modules**
| Module | Description | Status |
|--------|-------------|---------|
| [☁️ **AWS**](en/modules/aws.md) | Amazon Web Services | <span class="status-indicator implemented">Ready</span> |
| [🌩️ **GCP**](en/modules/gcp.md) | Google Cloud Platform | <span class="status-indicator implemented">Ready</span> |
| [🔷 **Azure**](en/modules/azure.md) | Microsoft Azure | <span class="status-indicator implemented">Ready</span> |
| [🌊 **DigitalOcean**](en/modules/digitalocean.md) | DigitalOcean | <span class="status-indicator beta">Beta</span> |

### **🛠️ Infrastructure Modules**
| Module | Description | Status |
|--------|-------------|---------|
| [🐳 **Docker**](en/modules/docker.md) | Container management | <span class="status-indicator implemented">Ready</span> |
| [🏗️ **Pulumi**](en/modules/pulumi.md) | Modern IaC | <span class="status-indicator implemented">Ready</span> |
| [🌍 **Terraform**](en/modules/terraform.md) | Infrastructure provisioning | <span class="status-indicator implemented">Ready</span> |
| [🧂 **Salt**](en/modules/salt.md) | Configuration management | <span class="status-indicator beta">Beta</span> |
| [🐍 **Python**](en/modules/python.md) | Python integration | <span class="status-indicator beta">Beta</span> |

## 🚀 **Get Started Today**

```bash
# Install Sloth Runner
curl -L https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner_linux_amd64.tar.gz | tar xz
chmod +x sloth-runner && sudo mv sloth-runner /usr/local/bin/

# Create your first workflow
echo 'local hello_task = task("greet"):command(function() log.info("Hello World! 🚀") return true end):build(); workflow.define("hello", { tasks = { hello_task } })' > hello.lua

# Run it!
sloth-runner run -f hello.lua
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