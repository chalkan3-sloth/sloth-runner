# 🚀 Complete Sloth Runner Features

## Overview

Complete documentation of **all** Sloth Runner features - from basic capabilities to advanced enterprise functionality. This guide serves as the master index for exploring all platform capabilities.

---

## 📋 Feature Index

### 🎯 Core
- [Workflow Execution](#workflow-execution)
- [Sloth DSL Language](#sloth-dsl-language)
- [Module System](#module-system)
- [State Management](#state-management)
- [Idempotency](#idempotency)

### 🌐 Distributed
- [Master-Agent Architecture](#master-agent-architecture)
- [Task Delegation](#task-delegation)
- [gRPC Communication](#grpc-communication)
- [Auto-Reconnection](#auto-reconnection)
- [Health Checks](#health-checks)

### 🎨 Interface
- [Modern Web UI](#modern-web-ui)
- [Complete CLI](#complete-cli)
- [Interactive REPL](#interactive-repl)
- [Remote Terminal](#remote-terminal)
- [REST API](#rest-api)

### 🔧 Automation
- [Scheduler (Cron)](#scheduler)
- [Hooks & Events](#hooks--events)
- [GitOps](#gitops)
- [CI/CD Integration](#cicd-integration)
- [Saved Workflows (Sloths)](#sloths)

### 📊 Monitoring
- [Telemetry](#telemetry)
- [Prometheus Metrics](#prometheus-metrics)
- [Grafana Dashboards](#grafana-dashboards)
- [Centralized Logs](#centralized-logs)
- [Agent Metrics](#agent-metrics)

### ☁️ Cloud & IaC
- [Multi-Cloud](#multi-cloud)
- [Terraform](#terraform)
- [Pulumi](#pulumi)
- [Kubernetes](#kubernetes)
- [Docker](#docker)

### 🔐 Security & Enterprise
- [Authentication](#authentication)
- [TLS/SSL](#tlsssl)
- [Audit Logs](#audit-logs)
- [Backups](#backups)
- [RBAC](#rbac)

### 🚀 Performance
- [Optimizations](#optimizations)
- [Parallel Execution](#parallel-execution)
- [Resource Limits](#resource-limits)
- [Caching](#caching)

---

## 🎯 Core

### Workflow Execution

**Description:** Central engine for executing workflows defined in Sloth files.

**Characteristics:**
- Sequential and parallel task execution
- Task group support
- Variables and templating
- Conditional execution
- Error handling and retry
- Dry-run mode
- Verbose output

**Commands:**
```bash
sloth-runner run <workflow> --file <file>
sloth-runner run <workflow> --file <file> --yes
sloth-runner run <workflow> --file <file> --group <group>
sloth-runner run <workflow> --file <file> --values vars.yaml
```

**Examples:**
```yaml
# Basic workflow
tasks:
  - name: Install nginx
    exec:
      script: |
        pkg.update()
        pkg.install("nginx")

  - name: Configure nginx
    exec:
      script: |
        file.copy("/src/nginx.conf", "/etc/nginx/nginx.conf")
        systemd.service_restart("nginx")
```

**Documentation:** `/docs/en/quick-start.md`

---

### Sloth DSL Language

**Description:** Declarative DSL based on YAML with embedded Lua scripting.

**Features:**
- **YAML-based** - familiar and readable syntax
- **Lua scripting** - power of a complete language
- **Type-safe** - type validation
- **Templating** - Go templates and Jinja2
- **Global modules** - no require() needed
- **Modern syntax** - supports modern features

**Structure:**
```yaml
# Metadata
version: "1.0"
description: "My workflow"

# Variables
vars:
  env: production
  version: "1.2.3"

# Groups
groups:
  deploy:
    - install_deps
    - build_app
    - deploy_app

# Tasks
tasks:
  - name: install_deps
    exec:
      script: |
        pkg.install({"nodejs", "npm"})

  - name: build_app
    exec:
      script: |
        exec.command("npm install")
        exec.command("npm run build")

  - name: deploy_app
    exec:
      script: |
        file.copy("./dist", "/var/www/app")
        systemd.service_restart("app")
    delegate_to: web-01
```

**Documentation:** `/docs/modern-dsl/introduction.md`

---

### Module System

**Description:** 40+ integrated modules providing all automation needs.

**Categories:**

#### 📦 System
- `pkg` - Package management (apt, yum, brew, etc.)
- `user` - User/group management
- `file` - File operations
- `systemd` - Service management
- `exec` - Command execution

#### 🐳 Containers
- `docker` - Complete Docker (containers, images, networks)
- `incus` - Incus/LXC containers and VMs
- `kubernetes` - K8s deploy and management

#### ☁️ Cloud
- `aws` - AWS (EC2, S3, RDS, Lambda, etc.)
- `azure` - Azure (VMs, Storage, etc.)
- `gcp` - GCP (Compute Engine, Cloud Storage, etc.)
- `digitalocean` - DigitalOcean (Droplets, Load Balancers)

#### 🏗️ IaC
- `terraform` - Terraform (init, plan, apply, destroy)
- `pulumi` - Pulumi
- `ansible` - Ansible playbooks

#### 🔧 Tools
- `git` - Git operations
- `ssh` - Remote SSH
- `net` - Networking (ping, http, download)
- `template` - Templates (Jinja2, Go)

#### 📊 Observability
- `log` - Structured logging
- `metrics` - Metrics (Prometheus)
- `notifications` - Notifications (Slack, Email, Discord, Telegram)

#### 🚀 Advanced
- `goroutine` - Parallel execution
- `reliability` - Retry, circuit breaker, timeout
- `state` - State management
- `facts` - System information
- `infra_test` - Infrastructure testing

**Complete list:** `sloth-runner modules list`

**Documentation:** `/docs/en/modules-complete.md`

---

### State Management

**Description:** State persistence system between executions.

**Features:**
- Persistent key-value store
- SQLite backend
- State scoping (global, workflow, task)
- Change detection
- State cleanup

**API:**
```lua
-- Save state
state.set("last_deploy_version", "v1.2.3")
state.set("deploy_timestamp", os.time())

-- Read state
local last_version = state.get("last_deploy_version")

-- Detect change
if state.changed("config_hash", new_hash) then
    log.info("Config changed, redeploying")
    deploy()
end

-- Clear state
state.clear("temporary_data")
```

**Documentation:** `/docs/state-management.md`

---

### Idempotency

**Description:** Ensures workflows can be executed multiple times with the same result.

**Features:**
- **Check mode** - checks before executing
- **State tracking** - tracks what was changed
- **Resource fingerprinting** - detects changes
- **Rollback** - undoes changes on error

**Example:**
```lua
-- Idempotent - checks before installing
if not pkg.is_installed("nginx") then
    pkg.install("nginx")
end

-- Idempotent - checks file hash
local current_hash = file.hash("/etc/nginx/nginx.conf")
if current_hash ~= expected_hash then
    file.copy("/src/nginx.conf", "/etc/nginx/nginx.conf")
    systemd.service_restart("nginx")
end
```

**Documentation:** `/docs/idempotency.md`

---

## 🌐 Distributed

### Master-Agent Architecture

**Description:** Distributed architecture with central master server and remote agents.

**Components:**
- **Master Server** - coordinates agents and workflows
- **Agent Nodes** - execute tasks remotely
- **gRPC Communication** - efficient and type-safe communication
- **Auto-Discovery** - agents self-register
- **Health Monitoring** - automatic heartbeats

**Topology:**
```
                    ┌──────────────┐
                    │   Master     │
                    │  (gRPC:50053)│
                    └──────┬───────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐       ┌────▼────┐      ┌────▼────┐
    │ Agent 1 │       │ Agent 2 │      │ Agent 3 │
    │  web-01 │       │  web-02 │      │   db-01 │
    └─────────┘       └─────────┘      └─────────┘
```

**Setup:**
```bash
# Start master
sloth-runner server --port 50053

# Install remote agent
sloth-runner agent install web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user ubuntu \
  --master 192.168.1.1:50053

# List agents
sloth-runner agent list
```

**Documentation:** `/docs/en/master-agent-architecture.md`

---

### Task Delegation

**Description:** Delegates task execution to specific agents.

**Features:**
- **Single delegation** - delegate to one agent
- **Multi delegation** - delegate to multiple agents in parallel
- **Round-robin** - distribute load
- **Failover** - fallback if agent fails
- **Conditional delegation** - delegate based on conditions

**Syntax:**
```yaml
# Delegate to one agent
tasks:
  - name: Deploy to web-01
    exec:
      script: |
        pkg.install("nginx")
    delegate_to: web-01

# Delegate to multiple agents
tasks:
  - name: Deploy to all web servers
    exec:
      script: |
        pkg.install("nginx")
    delegate_to:
      - web-01
      - web-02
      - web-03

# CLI - delegate entire workflow
sloth-runner run deploy --file deploy.sloth --delegate-to web-01
```

**Use with values:**
```yaml
# Pass agent-specific values
tasks:
  - name: Configure
    exec:
      script: |
        local ip = values.ip_address
        file.write("/etc/config", "IP=" .. ip)
    delegate_to: "{{ item }}"
    loop:
      - web-01
      - web-02
    values:
      web-01:
        ip_address: "192.168.1.10"
      web-02:
        ip_address: "192.168.1.11"
```

**Documentation:** `/docs/guides/values-delegate-to.md`

---

### gRPC Communication

**Description:** Efficient communication between master and agents using gRPC.

**Features:**
- **Streaming** - bi-directional streaming
- **Type-safe** - Protocol Buffers
- **Efficient** - binary protocol
- **Multiplexing** - multiple calls in one connection
- **TLS** - TLS/SSL support

**Services:**
```protobuf
service AgentService {
    rpc ExecuteTask(TaskRequest) returns (TaskResponse);
    rpc StreamLogs(LogRequest) returns (stream LogEntry);
    rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse);
    rpc GetMetrics(MetricsRequest) returns (MetricsResponse);
}
```

**Default port:** 50053

---

### Auto-Reconnection

**Description:** Agents automatically reconnect to master on disconnection.

**Features:**
- **Exponential backoff** - increases interval between attempts
- **Max retries** - configurable limit
- **Circuit breaker** - stops trying after many failures
- **Connection pooling** - reuses connections

**Configuration:**
```yaml
agent:
  reconnect:
    enabled: true
    initial_delay: 1s
    max_delay: 60s
    max_retries: -1  # infinite
```

**Documentation:** `/docs/en/agent-improvements.md`

---

### Health Checks

**Description:** Continuous agent health monitoring.

**Check types:**
- **Heartbeat** - periodic ping
- **Resource check** - CPU, memory, disk
- **Service check** - checks critical services
- **Custom checks** - user-defined checks

**Endpoints:**
```bash
# Health endpoint
curl http://agent:9090/health

# Metrics endpoint
curl http://agent:9090/metrics
```

**Thresholds:**
```yaml
health:
  cpu_threshold: 90  # %
  memory_threshold: 85  # %
  disk_threshold: 90  # %
  heartbeat_interval: 30s
  heartbeat_timeout: 90s
```

---

## 🎨 Interface

### Modern Web UI

**Description:** Complete, responsive, real-time web interface.

**Main features:**
- ✅ Dashboard with metrics and charts
- ✅ Agent management with real-time metrics
- ✅ Workflow editor with syntax highlighting
- ✅ Execution and log viewing
- ✅ Interactive terminal (xterm.js)
- ✅ Dark mode / Light mode
- ✅ WebSocket for real-time updates
- ✅ Mobile responsive
- ✅ Command palette (Ctrl+Shift+P)
- ✅ Drag & drop
- ✅ Glassmorphism design
- ✅ Smooth animations

**Pages:**
1. Dashboard (`/`)
2. Agents (`/agents`)
3. Agent Control (`/agent-control`)
4. Agent Dashboard (`/agent-dashboard`)
5. Workflows (`/workflows`)
6. Executions (`/executions`)
7. Hooks (`/hooks`)
8. Events (`/events`)
9. Scheduler (`/scheduler`)
10. Logs (`/logs`)
11. Terminal (`/terminal`)
12. Sloths (`/sloths`)
13. Settings (`/settings`)

**Technologies:**
- Bootstrap 5.3
- Chart.js 4.4
- xterm.js
- WebSockets
- Canvas API

**Start:**
```bash
sloth-runner ui --port 8080
```

**Access:** http://localhost:8080

**Documentation:** `/docs/en/web-ui-complete.md`

---

### Complete CLI

**Description:** Complete command-line interface with 100+ commands.

**Command categories:**

#### Execution
- `run` - Execute workflow
- `version` - View version

#### Agents
- `agent list` - List agents
- `agent get` - Agent details
- `agent install` - Install remote agent
- `agent update` - Update agent
- `agent start/stop/restart` - Control agent
- `agent modules` - List agent modules
- `agent metrics` - View metrics

#### Sloths (Saved Workflows)
- `sloth list` - List sloths
- `sloth add` - Add sloth
- `sloth get` - View sloth
- `sloth update` - Update sloth
- `sloth remove` - Remove sloth
- `sloth activate/deactivate` - Activate/deactivate

#### Hooks
- `hook list` - List hooks
- `hook add` - Add hook
- `hook remove` - Remove hook
- `hook enable/disable` - Enable/disable
- `hook test` - Test hook

#### Events
- `events list` - List events
- `events watch` - Monitor events in real-time

#### Database
- `db backup` - Database backup
- `db restore` - Restore database
- `db vacuum` - Optimize database
- `db stats` - Statistics

#### SSH
- `ssh list` - List SSH connections
- `ssh add` - Add connection
- `ssh remove` - Remove connection
- `ssh test` - Test connection

#### Modules
- `modules list` - List modules
- `modules info` - Module info

#### Server
- `server` - Start master server
- `ui` - Start Web UI
- `terminal` - Interactive terminal

#### Utilities
- `completion` - Shell auto-completion
- `doctor` - Diagnostics

**Documentation:** `/docs/en/cli-reference.md`

---

### Interactive REPL

**Description:** Read-Eval-Print Loop for interactive Lua code testing.

**Features:**
- **Complete Lua** - full Lua interpreter
- **Loaded modules** - all modules available
- **History** - command history
- **Auto-complete** - Tab completion
- **Multi-line** - multi-line code support
- **Pretty print** - formatted output

**Start:**
```bash
sloth-runner repl
```

**Example session:**
```lua
> pkg.install("nginx")
[OK] nginx installed successfully

> file.exists("/etc/nginx/nginx.conf")
true

> local content = file.read("/etc/nginx/nginx.conf")
> print(#content .. " bytes")
2048 bytes

> for i=1,5 do
>>   print("Hello " .. i)
>> end
Hello 1
Hello 2
Hello 3
Hello 4
Hello 5
```

**Special commands:**
- `.help` - help
- `.exit` - exit
- `.clear` - clear screen
- `.load <file>` - load file
- `.save <file>` - save session

**Documentation:** `/docs/en/repl.md`

---

### Remote Terminal

**Description:** Interactive terminal for remote agents via web UI.

**Features:**
- **xterm.js** - complete terminal emulator
- **Multiple sessions** - multiple simultaneous sessions
- **Tabs** - tab management
- **Command history** - command history (↑↓)
- **Copy/paste** - Ctrl+Shift+C/V
- **Themes** - various themes available
- **Upload/download** - file transfer

**Access:**
1. Web UI → Terminal
2. Select agent
3. Connect

**Special commands:**
```bash
.clear       # Clear terminal
.exit        # Close session
.upload <f>  # Upload file
.download <f># Download file
.theme <t>   # Change theme
```

**URL:** http://localhost:8080/terminal

---

### REST API

**Description:** Complete RESTful API for external integration.

**Main endpoints:**

#### Agents
```
GET    /api/v1/agents           # List agents
GET    /api/v1/agents/:name     # Agent details
POST   /api/v1/agents/:name/restart  # Restart agent
DELETE /api/v1/agents/:name     # Remove agent
```

#### Workflows
```
POST   /api/v1/workflows/run    # Execute workflow
GET    /api/v1/workflows/:id    # Workflow details
```

#### Executions
```
GET    /api/v1/executions       # List executions
GET    /api/v1/executions/:id   # Execution details
```

#### Hooks
```
GET    /api/v1/hooks            # List hooks
POST   /api/v1/hooks            # Create hook
DELETE /api/v1/hooks/:name      # Remove hook
```

#### Events
```
GET    /api/v1/events           # List events
```

#### Metrics
```
GET    /api/v1/metrics          # Prometheus metrics
```

**Authentication:**
```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/agents
```

**Examples:**
```bash
# List agents
curl http://localhost:8080/api/v1/agents

# Execute workflow
curl -X POST http://localhost:8080/api/v1/workflows/run \
  -H "Content-Type: application/json" \
  -d '{
    "file": "/workflows/deploy.sloth",
    "workflow_name": "deploy",
    "delegate_to": ["web-01"]
  }'

# View metrics
curl http://localhost:8080/api/v1/metrics
```

**Documentation:** `/docs/web-ui/api-reference.md`

---

## 🔧 Automation

### Scheduler

**Description:** Cron-based workflow scheduler.

**Features:**
- **Cron expressions** - complete cron syntax
- **Visual builder** - visual builder in Web UI
- **Timezone support** - timezone support
- **Missed run policy** - policy for missed runs
- **Overlap prevention** - prevents overlapping executions
- **Notifications** - success/failure notifications

**Create job:**
```bash
# Via CLI (coming soon)
sloth-runner scheduler add deploy-job \
  --workflow deploy.sloth \
  --schedule "0 3 * * *"  # Every day at 3am

# Via Web UI
http://localhost:8080/scheduler
```

**Cron syntax:**
```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday=0)
│ │ │ │ │
* * * * *

Examples:
0 * * * *     # Every hour
0 3 * * *     # Every day at 3am
0 0 * * 0     # Every Sunday at midnight
*/15 * * * *  # Every 15 minutes
```

**Documentation:** `/docs/scheduler.md`

---

### Hooks & Events

**Description:** Hook system for reacting to system events.

**Available events:**
- `workflow.started` - Workflow started
- `workflow.completed` - Workflow completed
- `workflow.failed` - Workflow failed
- `task.started` - Task started
- `task.completed` - Task completed
- `task.failed` - Task failed
- `agent.connected` - Agent connected
- `agent.disconnected` - Agent disconnected

**Create hook:**
```bash
sloth-runner hook add slack-notify \
  --event workflow.completed \
  --script /scripts/notify-slack.lua \
  --priority 10
```

**Hook script (Lua):**
```lua
-- /scripts/notify-slack.lua
local event = hook.event
local payload = hook.payload

if event == "workflow.completed" then
    notifications.slack(
        "https://hooks.slack.com/services/XXX/YYY/ZZZ",
        string.format("✅ Workflow '%s' completed!", payload.workflow_name),
        { channel = "#deployments" }
    )
end
```

**Available payload:**
```lua
-- workflow.* events
{
    workflow_name = "deploy",
    status = "success",
    duration = 45.3,
    started_at = 1234567890,
    completed_at = 1234567935
}

-- agent.* events
{
    agent_name = "web-01",
    address = "192.168.1.100:50060",
    status = "connected"
}
```

**Documentation:** `/docs/architecture/hooks-events-system.md`

---

### GitOps

**Description:** Complete GitOps pattern implementation.

**Features:**
- **Git-based** - Git as source of truth
- **Auto-sync** - automatic synchronization
- **Drift detection** - detects manual changes
- **Rollback** - automatic rollback
- **Multi-environment** - dev, staging, production
- **PR-based** - approval via Pull Requests

**GitOps workflow:**
```yaml
# .sloth/gitops.yaml
repos:
  - name: k8s-manifests
    url: https://github.com/org/k8s-manifests.git
    branch: main
    path: production/
    sync_interval: 5m
    auto_sync: true
    prune: true

hooks:
  on_sync:
    - notify-slack
  on_drift:
    - alert-team
```

**CLI:**
```bash
# Manual sync
sloth-runner gitops sync k8s-manifests

# View status
sloth-runner gitops status

# View drift
sloth-runner gitops diff
```

**Documentation:** `/docs/en/gitops-features.md`

---

### CI/CD Integration

**Description:** Integration with CI/CD pipelines.

**Support:**
- GitHub Actions
- GitLab CI
- Jenkins
- CircleCI
- Travis CI
- Azure Pipelines

**GitHub Actions example:**
```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Install Sloth Runner
        run: |
          curl -L https://github.com/org/sloth-runner/releases/latest/download/sloth-runner-linux-amd64 -o sloth-runner
          chmod +x sloth-runner

      - name: Run deployment
        env:
          SLOTH_RUNNER_MASTER_ADDR: ${{ secrets.SLOTH_MASTER }}
        run: |
          ./sloth-runner run deploy \
            --file deployments/production.sloth \
            --delegate-to web-01 \
            --yes
```

---

### Sloths

**Description:** Repository of saved and reusable workflows.

**Features:**
- **Versioning** - version history
- **Tags** - organization by tags
- **Search** - search by name/description/tags
- **Clone** - clone existing sloth
- **Export/Import** - share sloths
- **Active/Inactive** - activate/deactivate without deleting

**Commands:**
```bash
# Add sloth
sloth-runner sloth add deploy --file deploy.sloth

# List sloths
sloth-runner sloth list

# View sloth
sloth-runner sloth get deploy

# Execute sloth
sloth-runner run deploy --file $(sloth-runner sloth get deploy --show-path)

# Remove sloth
sloth-runner sloth remove deploy
```

**Documentation:** `/docs/features/sloth-management.md`

---

## 📊 Monitoring

### Telemetry

**Description:** Complete observability system.

**Components:**
- Prometheus metrics
- Structured logging
- Distributed tracing
- Health checks
- Performance profiling

**Architecture:**
```
┌──────────┐    metrics    ┌────────────┐
│  Master  ├───────────────► Prometheus │
└──────────┘               └─────┬──────┘
                                 │
┌──────────┐    metrics          │
│ Agent 1  ├─────────────────────┤
└──────────┘                     │
                                 ▼
┌──────────┐    metrics    ┌──────────┐
│ Agent 2  ├───────────────►  Grafana │
└──────────┘               └──────────┘
```

**Endpoints:**
```
http://master:9090/metrics
http://agent:9091/metrics
```

**Documentation:** `/docs/en/telemetry/index.md`

---

### Prometheus Metrics

**Description:** Metrics exported in Prometheus format.

**Available metrics:**

#### Workflows
```
sloth_workflow_executions_total{status="success|failed"}
sloth_workflow_duration_seconds{workflow="name"}
sloth_workflow_tasks_total{workflow="name"}
```

#### Agents
```
sloth_agent_connected_total
sloth_agent_cpu_usage_percent{agent="name"}
sloth_agent_memory_usage_bytes{agent="name"}
sloth_agent_disk_usage_bytes{agent="name"}
```

#### System
```
sloth_tasks_executed_total
sloth_hooks_triggered_total{event="type"}
sloth_db_size_bytes
```

**Scrape config:**
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'sloth-master'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'sloth-agents'
    static_configs:
      - targets:
        - 'agent1:9091'
        - 'agent2:9091'
```

**Documentation:** `/docs/en/telemetry/prometheus-metrics.md`

---

### Grafana Dashboards

**Description:** Pre-configured dashboards for Grafana.

**Dashboards:**
1. **Overview** - system overview
2. **Agents** - metrics for all agents
3. **Workflows** - executions and performance
4. **Resources** - CPU, memory, disk, network

**Import dashboard:**
```bash
# Generate dashboard JSON
sloth-runner agent metrics grafana web-01 --export dashboard.json

# Import to Grafana
curl -X POST http://admin:admin@localhost:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -d @dashboard.json
```

**Features:**
- Auto-refresh (5s, 10s, 30s, 1m)
- Time range selector
- Variables (agent, workflow)
- Configurable alerts
- Export PNG/PDF

**Documentation:** `/docs/en/telemetry/grafana-dashboard.md`

---

### Centralized Logs

**Description:** Centralized structured logging system.

**Features:**
- **Structured** - JSON structured logs
- **Levels** - debug, info, warn, error
- **Context** - rich metadata
- **Search** - search by any field
- **Export** - JSON, CSV, text
- **Retention** - retention policy

**Format:**
```json
{
  "timestamp": "2025-10-07T10:30:45Z",
  "level": "info",
  "message": "Workflow completed",
  "workflow": "deploy",
  "agent": "web-01",
  "duration": 45.3,
  "status": "success"
}
```

**Access:**
```bash
# CLI
sloth-runner logs --follow

# Web UI
http://localhost:8080/logs

# API
curl http://localhost:8080/api/v1/logs?level=error&since=1h
```

---

### Agent Metrics

**Description:** Detailed real-time agent metrics.

**Collected metrics:**
- CPU usage (%)
- Memory usage (bytes, %)
- Disk usage (bytes, %)
- Load average (1m, 5m, 15m)
- Network I/O (bytes/sec)
- Process count
- Uptime

**Visualization:**
```bash
# CLI
sloth-runner agent metrics web-01
sloth-runner agent metrics web-01 --watch

# Web UI - Agent Dashboard
http://localhost:8080/agent-dashboard?agent=web-01

# API
curl http://localhost:8080/api/v1/agents/web-01/metrics
```

**Format:**
```json
{
  "cpu": {
    "cores": 4,
    "usage_percent": 45.2,
    "load_avg": [1.2, 0.8, 0.5]
  },
  "memory": {
    "total_bytes": 8589934592,
    "used_bytes": 4294967296,
    "usage_percent": 50.0
  },
  "disk": {
    "total_bytes": 107374182400,
    "used_bytes": 53687091200,
    "usage_percent": 50.0
  }
}
```

---

## ☁️ Cloud & IaC

### Multi-Cloud

**Description:** Native support for multiple cloud providers.

**Supported providers:**
- ✅ AWS (EC2, S3, RDS, Lambda, ECS, EKS, etc.)
- ✅ Azure (VMs, Storage, AKS, Functions, etc.)
- ✅ GCP (Compute Engine, Cloud Storage, GKE, etc.)
- ✅ DigitalOcean (Droplets, Spaces, K8s, etc.)
- ✅ Linode
- ✅ Vultr
- ✅ Hetzner Cloud

**Multi-cloud example:**
```yaml
# Deploy to AWS and GCP simultaneously
tasks:
  - name: Deploy to AWS
    exec:
      script: |
        aws.ec2_instance_create({
          image_id = "ami-xxx",
          instance_type = "t3.medium"
        })
    delegate_to: aws-agent

  - name: Deploy to GCP
    exec:
      script: |
        gcp.compute_instance_create({
          machine_type = "e2-medium",
          image_family = "ubuntu-2204-lts"
        })
    delegate_to: gcp-agent
```

**Documentation:** `/docs/en/enterprise-features.md`

---

### Terraform

**Description:** Complete Terraform integration.

**Features:**
- `terraform.init` - Initialize
- `terraform.plan` - Plan
- `terraform.apply` - Apply
- `terraform.destroy` - Destroy
- State management
- Backend config
- Variable files

**Example:**
```lua
local tf_dir = "/infra/terraform"

-- Initialize
terraform.init(tf_dir, {
    backend_config = {
        bucket = "my-tf-state",
        key = "prod/terraform.tfstate"
    }
})

-- Plan
local plan = terraform.plan(tf_dir, {
    var_file = "production.tfvars",
    vars = {
        region = "us-east-1",
        environment = "production"
    }
})

-- Apply if there are changes
if plan.changes > 0 then
    terraform.apply(tf_dir, {
        auto_approve = true
    })
end
```

**Documentation:** `/docs/modules/terraform.md`

---

### Pulumi

**Description:** Pulumi integration.

**Support:**
- Stack management
- Configuration
- Up/Deploy
- Destroy
- Preview

**Example:**
```lua
-- Select stack
pulumi.stack_select("production")

-- Configure
pulumi.config_set("aws:region", "us-east-1")

-- Deploy
pulumi.up({
    yes = true,  -- auto-approve
    parallel = 10
})
```

**Documentation:** `/docs/modules/pulumi.md`

---

### Kubernetes

**Description:** Kubernetes deploy and management.

**Features:**
- Apply manifests
- Helm charts
- Namespaces
- ConfigMaps/Secrets
- Rollouts
- Health checks

**Example:**
```lua
-- Apply manifests
kubernetes.apply("/k8s/deployment.yaml", {
    namespace = "production"
})

-- Helm install
helm.install("myapp", "charts/myapp", {
    namespace = "production",
    values = {
        image = {
            tag = "v1.2.3"
        }
    }
})

-- Wait for rollout
kubernetes.rollout_status("deployment/myapp", {
    namespace = "production",
    timeout = "5m"
})
```

**Documentation:** `/docs/en/gitops/kubernetes.md`

---

### Docker

**Description:** Complete Docker automation.

**Functionality:**
- Container lifecycle (run, stop, remove)
- Image management (build, push, pull)
- Networks (create, connect)
- Volumes (create, mount)
- Docker Compose

**Deployment example:**
```lua
-- Build image
docker.image_build(".", {
    tag = "myapp:v1.2.3",
    build_args = {
        VERSION = "1.2.3"
    }
})

-- Push to registry
docker.image_push("myapp:v1.2.3", {
    registry = "registry.example.com"
})

-- Deploy
docker.container_run("myapp:v1.2.3", {
    name = "app",
    ports = {"3000:3000"},
    env = {
        DATABASE_URL = "postgres://..."
    },
    restart = "unless-stopped"
})
```

**Documentation:** `/docs/modules/docker.md`

---

## 🔐 Security & Enterprise

### Authentication

**Description:** Authentication system for Web UI and API.

**Methods:**
- Username/Password
- JWT tokens
- OAuth2 (GitHub, Google, etc.)
- LDAP/AD
- SSO

**Setup:**
```yaml
# config.yaml
auth:
  enabled: true
  type: jwt
  jwt:
    secret: "your-secret-key"
    expiry: 24h
  oauth:
    providers:
      - github:
          client_id: "xxx"
          client_secret: "yyy"
```

---

### TLS/SSL

**Description:** TLS/SSL support for secure communication.

**Features:**
- gRPC TLS
- HTTPS Web UI
- Certificate management
- Auto-renewal (Let's Encrypt)

**Configuration:**
```bash
# Master with TLS
sloth-runner server \
  --tls-cert /etc/sloth/cert.pem \
  --tls-key /etc/sloth/key.pem

# Agent with TLS
sloth-runner agent start \
  --master-tls-cert /etc/sloth/master-cert.pem
```

---

### Audit Logs

**Description:** Audit logs for all actions.

**Audited events:**
- User login/logout
- Workflow execution
- Configuration changes
- API calls
- Admin actions

**Format:**
```json
{
  "timestamp": "2025-10-07T10:30:45Z",
  "event": "workflow.executed",
  "user": "admin",
  "ip": "192.168.1.100",
  "resource": "deploy.sloth",
  "action": "execute",
  "result": "success"
}
```

---

### Backups

**Description:** Automated backup system.

**Features:**
- Configurable auto-backup
- Compression (gzip)
- Retention policy
- Remote backup (S3, Azure Blob, etc.)
- Restore

**Commands:**
```bash
# Manual backup
sloth-runner db backup --output /backup/sloth.db --compress

# Restore
sloth-runner db restore /backup/sloth.db.gz --decompress

# Automated backup (cron)
0 3 * * * sloth-runner db backup --output /backup/sloth-$(date +\%Y\%m\%d).db --compress
```

---

### RBAC

**Description:** Role-Based Access Control.

**Roles:**
- **Admin** - full access
- **Operator** - execute workflows, manage agents
- **Developer** - create/edit workflows
- **Viewer** - view only

**Permissions:**
```yaml
roles:
  operator:
    permissions:
      - workflow:execute
      - agent:view
      - agent:restart
      - logs:view

  developer:
    permissions:
      - workflow:create
      - workflow:edit
      - workflow:execute
      - logs:view

  viewer:
    permissions:
      - workflow:view
      - agent:view
      - logs:view
```

---

## 🚀 Performance

### Optimizations

**Description:** Recent performance optimizations.

**Implemented improvements:**

#### Agent Optimizations
- ✅ **Ultra-low memory** - 32MB RAM footprint
- ✅ **Binary size reduction** - from 45MB → 12MB
- ✅ **Startup time** - <100ms
- ✅ **CPU efficiency** - 99% idle when inactive

#### Database Optimizations
- ✅ **WAL mode** - Write-Ahead Logging
- ✅ **Connection pooling** - connection reuse
- ✅ **Prepared statements** - optimized queries
- ✅ **Indexes** - indexes on critical fields
- ✅ **Auto-vacuum** - automatic cleanup

#### gRPC Optimizations
- ✅ **Connection reuse** - keepalive
- ✅ **Compression** - gzip compression
- ✅ **Multiplexing** - multiple streams
- ✅ **Buffer pooling** - buffer reuse

**Benchmark:**
```
Before:
- Agent memory: 128MB
- Binary size: 45MB
- Startup time: 2s

After:
- Agent memory: 32MB (75% reduction)
- Binary size: 12MB (73% reduction)
- Startup time: 95ms (95% faster)
```

**Documentation:** `/PERFORMANCE_OPTIMIZATIONS.md`

---

### Parallel Execution

**Description:** Parallel task execution using goroutines.

**Features:**
- **goroutine.parallel()** - execute functions in parallel
- **Concurrency control** - limit of simultaneous goroutines
- **Error handling** - collects errors from all goroutines
- **Wait groups** - automatic synchronization

**Example:**
```lua
-- Execute multiple tasks in parallel
goroutine.parallel({
    function()
        pkg.install("nginx")
    end,
    function()
        pkg.install("postgresql")
    end,
    function()
        pkg.install("redis")
    end
})

-- With concurrency limit
goroutine.parallel({
    tasks = {
        function() exec.command("task1") end,
        function() exec.command("task2") end,
        function() exec.command("task3") end,
        function() exec.command("task4") end
    },
    max_concurrent = 2  -- Maximum 2 at a time
})
```

**Documentation:** `/docs/modules/goroutine.md`

---

### Resource Limits

**Description:** Configurable resource limits.

**Configuration:**
```yaml
# Agent config
resources:
  cpu:
    limit: 2  # cores
    reserve: 0.5
  memory:
    limit: 2GB
    reserve: 512MB
  disk:
    limit: 10GB
    min_free: 1GB
```

**Enforcement:**
- CPU throttling
- Memory limits (cgroup)
- Disk quota
- Task timeout

---

### Caching

**Description:** Caching system for optimization.

**Cache types:**

#### Module cache
- Compiled Lua modules
- Reduce load time

#### State cache
- State in memory
- Reduce DB queries

#### Metrics cache
- Aggregated metrics
- Reduce computation

**Configuration:**
```yaml
cache:
  enabled: true
  ttl: 5m
  max_size: 100MB
  eviction: lru  # least recently used
```

---

## 📚 Additional Resources

### Documentation
- [🚀 Quick Start](/docs/en/quick-start.md)
- [🏗️ Architecture](/docs/architecture/sloth-runner-architecture.md)
- [📖 Modern DSL](/docs/modern-dsl/introduction.md)
- [🎯 Advanced Examples](/docs/en/advanced-examples.md)

### Useful Links
- [GitHub Repository](https://github.com/chalkan3/sloth-runner)
- [Issue Tracker](https://github.com/chalkan3/sloth-runner/issues)
- [Releases](https://github.com/chalkan3/sloth-runner/releases)

---

**Last updated:** 2025-10-07

**Total Documented Features:** 100+
