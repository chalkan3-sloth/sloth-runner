# 🦥 Sloth Runner - Enterprise Task Orchestration Platform

<div class="hero">
  <h1>🚀 The Modern Task Orchestration Platform</h1>
  <p>Powerful, flexible, and scalable task automation with Lua scripting, distributed agents, and enterprise-grade reliability.</p>
  <a href="TUTORIAL/" class="btn">🚀 Get Started</a>
  <a href="en/quick-start/" class="btn">⚡ Quick Start</a>
</div>

## 🌟 **Why Choose Sloth Runner?**

<div class="feature-grid">
  <div class="feature-card">
    <div class="icon">🌙</div>
    <h3>Lua Scripting Power</h3>
    <p>Write flexible, readable task definitions in Lua with full access to system resources and cloud APIs.</p>
  </div>
  
  <div class="feature-card">
    <div class="icon">🌐</div>
    <h3>Distributed by Design</h3>
    <p>Native master-agent architecture with real-time streaming, automatic failover, and intelligent load balancing.</p>
  </div>
  
  <div class="feature-card">
    <div class="icon">💾</div>
    <h3>Persistent State Management</h3>
    <p>Built-in SQLite-based state with atomic operations, distributed locks, and TTL support for complex workflows.</p>
  </div>
  
  <div class="feature-card">
    <div class="icon">📊</div>
    <h3>Advanced Monitoring</h3>
    <p>Real-time metrics, health checks, and Prometheus-compatible endpoints for complete observability.</p>
  </div>
  
  <div class="feature-card">
    <div class="icon">🔐</div>
    <h3>Enterprise Security</h3>
    <p>mTLS authentication, RBAC authorization, audit logging, and compliance-ready security features.</p>
  </div>
  
  <div class="feature-card">
    <div class="icon">☁️</div>
    <h3>Multi-Cloud Ready</h3>
    <p>Native integrations with AWS, GCP, Azure, and on-premises infrastructure for hybrid deployments.</p>
  </div>
</div>

## 🚀 **Key Features**

### **💾 State Management & Persistence** <span class="status-indicator implemented">Implemented</span>
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

### **📊 Metrics & Monitoring** <span class="status-indicator implemented">Implemented</span>
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

### **🌐 Distributed Agent System** <span class="status-indicator implemented">Implemented</span>
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
curl -L https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner-linux-amd64 -o sloth-runner
chmod +x sloth-runner && sudo mv sloth-runner /usr/local/bin/

# Create your first Modern DSL task
echo 'local hello_task = task("greet"):command(function() log.info("Hello Modern DSL! 🚀") return true end):build(); workflow.define("hello", { tasks = { hello_task } })' > hello.lua

# Run it!
sloth-runner run -f hello.lua
```

## 🤝 **Community & Support**

- 💬 [GitHub Discussions](https://github.com/chalkan3-sloth/sloth-runner/discussions)
- 🐛 [Issue Tracker](https://github.com/chalkan3-sloth/sloth-runner/issues)
- 📧 [Email Support](mailto:support@sloth-runner.dev)
- 💼 [Enterprise Support](mailto:enterprise@sloth-runner.dev)

<div class="hero">
  <h2>Ready to Transform Your Task Automation? 🚀</h2>
  <p>Join thousands of developers using Sloth Runner for reliable, scalable task orchestration</p>
  <a href="en/quick-start/" class="btn">Start Your Journey →</a>
  <a href="https://github.com/chalkan3-sloth/sloth-runner" class="btn">View on GitHub →</a>
</div>
