# 🚀 Sloth Runner 完整功能

## 概述

Sloth Runner 的**所有**功能的完整文档 - 从基本功能到高级企业功能。本指南是探索平台所有能力的主索引。

---

## 📋 功能索引

### 🎯 核心（Core）
- [工作流执行](#工作流执行)
- [Sloth DSL 语言](#sloth-dsl-语言)
- [模块系统](#模块系统)
- [状态管理](#状态管理)
- [幂等性](#幂等性)

### 🌐 分布式
- [Master-Agent 架构](#master-agent-架构)
- [任务委托](#任务委托)
- [gRPC 通信](#grpc-通信)
- [自动重连](#自动重连)
- [健康检查](#健康检查)

### 🎨 界面
- [现代 Web UI](#现代-web-ui)
- [完整 CLI](#完整-cli)
- [交互式 REPL](#交互式-repl)
- [远程终端](#远程终端)
- [REST API](#rest-api)

### 🔧 自动化
- [调度器 (Cron)](#调度器)
- [钩子和事件](#钩子和事件)
- [GitOps](#gitops)
- [CI/CD 集成](#cicd-集成)
- [保存的工作流 (Sloths)](#sloths)

### 📊 监控
- [遥测](#遥测)
- [Prometheus 指标](#prometheus-指标)
- [Grafana 仪表板](#grafana-仪表板)
- [集中式日志](#集中式日志)
- [代理指标](#代理指标)

### ☁️ 云和 IaC
- [多云](#多云)
- [Terraform](#terraform)
- [Pulumi](#pulumi)
- [Kubernetes](#kubernetes)
- [Docker](#docker)

### 🔐 安全和企业
- [身份验证](#身份验证)
- [TLS/SSL](#tlsssl)
- [审计日志](#审计日志)
- [备份](#备份)
- [RBAC](#rbac)

### 🚀 性能
- [优化](#优化)
- [并行执行](#并行执行)
- [资源限制](#资源限制)
- [缓存](#缓存)

---

## 🎯 核心（Core）

### 工作流执行

**描述：** 执行在 Sloth 文件中定义的工作流的核心引擎。

**特性：**
- 任务的顺序和并行执行
- 支持任务组
- 变量和模板
- 条件执行
- 错误处理和重试
- 试运行模式
- 详细输出

**命令：**
```bash
sloth-runner run <workflow> --file <文件>
sloth-runner run <workflow> --file <文件> --yes
sloth-runner run <workflow> --file <文件> --group <组>
sloth-runner run <workflow> --file <文件> --values vars.yaml
```

**示例：**
```yaml
# 基本工作流
tasks:
  - name: 安装 nginx
    exec:
      script: |
        pkg.update()
        pkg.install("nginx")

  - name: 配置 nginx
    exec:
      script: |
        file.copy("/src/nginx.conf", "/etc/nginx/nginx.conf")
        systemd.service_restart("nginx")
```

**文档：** `/docs/en/quick-start.md`

---

### Sloth DSL 语言

**描述：** 基于 YAML 的声明式 DSL，内嵌 Lua 脚本。

**特性：**
- **基于 YAML** - 熟悉且可读的语法
- **Lua 脚本** - 完整语言的强大功能
- **类型安全** - 类型验证
- **模板** - Go 模板和 Jinja2
- **全局模块** - 无需 require()
- **现代语法** - 支持现代特性

**结构：**
```yaml
# 元数据
version: "1.0"
description: "我的工作流"

# 变量
vars:
  env: production
  version: "1.2.3"

# 组
groups:
  deploy:
    - install_deps
    - build_app
    - deploy_app

# 任务
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

**文档：** `/docs/modern-dsl/introduction.md`

---

### 模块系统

**描述：** 40 多个集成模块，满足所有自动化需求。

**类别：**

#### 📦 系统
- `pkg` - 包管理（apt、yum、brew 等）
- `user` - 用户/组管理
- `file` - 文件操作
- `systemd` - 服务管理
- `exec` - 命令执行

#### 🐳 容器
- `docker` - 完整 Docker（容器、镜像、网络）
- `incus` - Incus/LXC 容器和虚拟机
- `kubernetes` - K8s 部署和管理

#### ☁️ 云
- `aws` - AWS（EC2、S3、RDS、Lambda 等）
- `azure` - Azure（虚拟机、存储等）
- `gcp` - GCP（Compute Engine、Cloud Storage 等）
- `digitalocean` - DigitalOcean（Droplet、负载均衡器）

#### 🏗️ IaC
- `terraform` - Terraform（init、plan、apply、destroy）
- `pulumi` - Pulumi
- `ansible` - Ansible playbook

#### 🔧 工具
- `git` - Git 操作
- `ssh` - 远程 SSH
- `net` - 网络（ping、http、download）
- `template` - 模板（Jinja2、Go）

#### 📊 可观测性
- `log` - 结构化日志
- `metrics` - 指标（Prometheus）
- `notifications` - 通知（Slack、Email、Discord、Telegram）

#### 🚀 高级
- `goroutine` - 并行执行
- `reliability` - 重试、断路器、超时
- `state` - 状态管理
- `facts` - 系统信息
- `infra_test` - 基础设施测试

**完整列表：** `sloth-runner modules list`

**文档：** `/docs/zh/modules-complete.md`

---

### 状态管理

**描述：** 执行之间的状态持久化系统。

**特性：**
- 持久键值存储
- SQLite 后端
- 状态作用域（全局、工作流、任务）
- 变更检测
- 状态清理

**API：**
```lua
-- 保存状态
state.set("last_deploy_version", "v1.2.3")
state.set("deploy_timestamp", os.time())

-- 读取状态
local last_version = state.get("last_deploy_version")

-- 检测变更
if state.changed("config_hash", new_hash) then
    log.info("配置已变更，重新部署")
    deploy()
end

-- 清除状态
state.clear("temporary_data")
```

**文档：** `/docs/state-management.md`

---

### 幂等性

**描述：** 保证工作流可以多次执行并获得相同结果。

**特性：**
- **检查模式** - 执行前检查
- **状态跟踪** - 跟踪已变更的内容
- **资源指纹** - 检测变更
- **回滚** - 出错时撤消更改

**示例：**
```lua
-- 幂等 - 安装前检查
if not pkg.is_installed("nginx") then
    pkg.install("nginx")
end

-- 幂等 - 检查文件哈希
local current_hash = file.hash("/etc/nginx/nginx.conf")
if current_hash ~= expected_hash then
    file.copy("/src/nginx.conf", "/etc/nginx/nginx.conf")
    systemd.service_restart("nginx")
end
```

**文档：** `/docs/idempotency.md`

---

## 🌐 分布式

### Master-Agent 架构

**描述：** 带有中央主服务器和远程代理的分布式架构。

**组件：**
- **主服务器** - 协调代理和工作流
- **代理节点** - 远程执行任务
- **gRPC 通信** - 高效且类型安全的通信
- **自动发现** - 代理自注册
- **健康监控** - 自动心跳

**拓扑：**
```
                    ┌──────────────┐
                    │   主服务器   │
                    │  (gRPC:50053)│
                    └──────┬───────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐       ┌────▼────┐      ┌────▼────┐
    │代理 1   │       │代理 2   │      │代理 3   │
    │  web-01 │       │  web-02 │      │   db-01 │
    └─────────┘       └─────────┘      └─────────┘
```

**设置：**
```bash
# 启动主服务器
sloth-runner server --port 50053

# 安装远程代理
sloth-runner agent install web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user ubuntu \
  --master 192.168.1.1:50053

# 列出代理
sloth-runner agent list
```

**文档：** `/docs/en/master-agent-architecture.md`

---

### 任务委托

**描述：** 将任务执行委托给特定代理。

**特性：**
- **单一委托** - 委托给一个代理
- **多重委托** - 并行委托给多个代理
- **轮询** - 分配负载
- **故障转移** - 如果代理失败则回退
- **条件委托** - 基于条件委托

**语法：**
```yaml
# 委托给一个代理
tasks:
  - name: 部署到 web-01
    exec:
      script: |
        pkg.install("nginx")
    delegate_to: web-01

# 委托给多个代理
tasks:
  - name: 部署到所有 web 服务器
    exec:
      script: |
        pkg.install("nginx")
    delegate_to:
      - web-01
      - web-02
      - web-03

# CLI - 委托整个工作流
sloth-runner run deploy --file deploy.sloth --delegate-to web-01
```

**带值的使用：**
```yaml
# 为每个代理传递特定值
tasks:
  - name: 配置
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

**文档：** `/docs/guides/values-delegate-to.md`

---

### gRPC 通信

**描述：** 使用 gRPC 在主服务器和代理之间进行高效通信。

**特性：**
- **流式传输** - 双向流式传输
- **类型安全** - Protocol Buffers
- **高效** - 二进制协议
- **多路复用** - 单连接上的多个调用
- **TLS** - TLS/SSL 支持

**服务：**
```protobuf
service AgentService {
    rpc ExecuteTask(TaskRequest) returns (TaskResponse);
    rpc StreamLogs(LogRequest) returns (stream LogEntry);
    rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse);
    rpc GetMetrics(MetricsRequest) returns (MetricsResponse);
}
```

**默认端口：** 50053

---

### 自动重连

**描述：** 断开连接时代理自动重新连接到主服务器。

**特性：**
- **指数退避** - 尝试之间增加间隔
- **最大重试** - 可配置限制
- **断路器** - 多次失败后停止尝试
- **连接池** - 重用连接

**配置：**
```yaml
agent:
  reconnect:
    enabled: true
    initial_delay: 1s
    max_delay: 60s
    max_retries: -1  # 无限
```

**文档：** `/docs/en/agent-improvements.md`

---

### 健康检查

**描述：** 持续监控代理健康。

**检查类型：**
- **心跳** - 定期 ping
- **资源检查** - CPU、内存、磁盘
- **服务检查** - 检查关键服务
- **自定义检查** - 用户定义的检查

**端点：**
```bash
# 健康端点
curl http://agent:9090/health

# 指标端点
curl http://agent:9090/metrics
```

**阈值：**
```yaml
health:
  cpu_threshold: 90  # %
  memory_threshold: 85  # %
  disk_threshold: 90  # %
  heartbeat_interval: 30s
  heartbeat_timeout: 90s
```

---

## 🎨 界面

### 现代 Web UI

**描述：** 完整、响应式和实时的 Web 界面。

**主要功能：**
- ✅ 带指标和图表的仪表板
- ✅ 带实时指标的代理管理
- ✅ 带语法高亮的工作流编辑器
- ✅ 执行和日志可视化
- ✅ 交互式终端（xterm.js）
- ✅ 深色/浅色模式
- ✅ WebSocket 实时更新
- ✅ 移动响应式
- ✅ 命令面板（Ctrl+Shift+P）
- ✅ 拖放
- ✅ 毛玻璃设计
- ✅ 平滑动画

**页面：**
1. 仪表板 (`/`)
2. 代理 (`/agents`)
3. 代理控制 (`/agent-control`)
4. 代理仪表板 (`/agent-dashboard`)
5. 工作流 (`/workflows`)
6. 执行 (`/executions`)
7. 钩子 (`/hooks`)
8. 事件 (`/events`)
9. 调度器 (`/scheduler`)
10. 日志 (`/logs`)
11. 终端 (`/terminal`)
12. Sloth (`/sloths`)
13. 设置 (`/settings`)

**技术：**
- Bootstrap 5.3
- Chart.js 4.4
- xterm.js
- WebSockets
- Canvas API

**启动：**
```bash
sloth-runner ui --port 8080
```

**访问：** http://localhost:8080

**文档：** `/docs/zh/web-ui-complete.md`

---

### 完整 CLI

**描述：** 具有 100 多个命令的完整命令行界面。

**命令类别：**

#### 执行
- `run` - 执行工作流
- `version` - 查看版本

#### 代理
- `agent list` - 列出代理
- `agent get` - 代理详情
- `agent install` - 安装远程代理
- `agent update` - 更新代理
- `agent start/stop/restart` - 控制代理
- `agent modules` - 列出代理模块
- `agent metrics` - 查看指标

#### Sloth（保存的工作流）
- `sloth list` - 列出 sloth
- `sloth add` - 添加 sloth
- `sloth get` - 查看 sloth
- `sloth update` - 更新 sloth
- `sloth remove` - 删除 sloth
- `sloth activate/deactivate` - 激活/停用

#### 钩子
- `hook list` - 列出钩子
- `hook add` - 添加钩子
- `hook remove` - 删除钩子
- `hook enable/disable` - 启用/禁用
- `hook test` - 测试钩子

#### 事件
- `events list` - 列出事件
- `events watch` - 实时监控事件

#### 数据库
- `db backup` - 备份数据库
- `db restore` - 恢复数据库
- `db vacuum` - 优化数据库
- `db stats` - 统计信息

#### SSH
- `ssh list` - 列出 SSH 连接
- `ssh add` - 添加连接
- `ssh remove` - 删除连接
- `ssh test` - 测试连接

#### 模块
- `modules list` - 列出模块
- `modules info` - 模块信息

#### 服务器
- `server` - 启动主服务器
- `ui` - 启动 Web UI
- `terminal` - 交互式终端

#### 实用工具
- `completion` - Shell 自动补全
- `doctor` - 诊断

**文档：** `/docs/zh/cli-reference.md`

---

### 交互式 REPL

**描述：** 交互式测试 Lua 代码的读取-求值-打印循环。

**特性：**
- **完整 Lua** - 完整的 Lua 解释器
- **模块已加载** - 所有模块可用
- **历史记录** - 命令历史
- **自动补全** - Tab 补全
- **多行** - 支持多行代码
- **美化打印** - 格式化输出

**启动：**
```bash
sloth-runner repl
```

**会话示例：**
```lua
> pkg.install("nginx")
[OK] nginx 安装成功

> file.exists("/etc/nginx/nginx.conf")
true

> local content = file.read("/etc/nginx/nginx.conf")
> print(#content .. " 字节")
2048 字节

> for i=1,5 do
>>   print("你好 " .. i)
>> end
你好 1
你好 2
你好 3
你好 4
你好 5
```

**特殊命令：**
- `.help` - 帮助
- `.exit` - 退出
- `.clear` - 清屏
- `.load <file>` - 加载文件
- `.save <file>` - 保存会话

**文档：** `/docs/en/repl.md`

---

### 远程终端

**描述：** 通过 Web UI 连接远程代理的交互式终端。

**特性：**
- **xterm.js** - 完整的终端模拟器
- **多会话** - 同时多个会话
- **标签页** - 标签管理
- **命令历史** - 命令历史（↑↓）
- **复制/粘贴** - Ctrl+Shift+C/V
- **主题** - 多种主题可用
- **上传/下载** - 文件传输

**访问：**
1. Web UI → 终端
2. 选择代理
3. 连接

**特殊命令：**
```bash
.clear       # 清除终端
.exit        # 关闭会话
.upload <f>  # 上传文件
.download <f># 下载文件
.theme <t>   # 更换主题
```

**URL：** http://localhost:8080/terminal

---

### REST API

**描述：** 用于外部集成的完整 RESTful API。

**主要端点：**

#### 代理
```
GET    /api/v1/agents           # 列出代理
GET    /api/v1/agents/:name     # 代理详情
POST   /api/v1/agents/:name/restart  # 重启代理
DELETE /api/v1/agents/:name     # 删除代理
```

#### 工作流
```
POST   /api/v1/workflows/run    # 执行工作流
GET    /api/v1/workflows/:id    # 工作流详情
```

#### 执行
```
GET    /api/v1/executions       # 列出执行
GET    /api/v1/executions/:id   # 执行详情
```

#### 钩子
```
GET    /api/v1/hooks            # 列出钩子
POST   /api/v1/hooks            # 创建钩子
DELETE /api/v1/hooks/:name      # 删除钩子
```

#### 事件
```
GET    /api/v1/events           # 列出事件
```

#### 指标
```
GET    /api/v1/metrics          # Prometheus 指标
```

**身份验证：**
```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/agents
```

**示例：**
```bash
# 列出代理
curl http://localhost:8080/api/v1/agents

# 执行工作流
curl -X POST http://localhost:8080/api/v1/workflows/run \
  -H "Content-Type: application/json" \
  -d '{
    "file": "/workflows/deploy.sloth",
    "workflow_name": "deploy",
    "delegate_to": ["web-01"]
  }'

# 查看指标
curl http://localhost:8080/api/v1/metrics
```

**文档：** `/docs/web-ui/api-reference.md`

---

## 🔧 自动化

### 调度器

**描述：** 基于 cron 的工作流调度器。

**特性：**
- **Cron 表达式** - 完整的 cron 语法
- **可视化构建器** - Web UI 中的可视化构建器
- **时区支持** - 支持时区
- **未执行策略** - 未执行运行的策略
- **重叠预防** - 防止执行重叠
- **通知** - 成功/失败通知

**创建作业：**
```bash
# 通过 CLI（即将推出）
sloth-runner scheduler add deploy-job \
  --workflow deploy.sloth \
  --schedule "0 3 * * *"  # 每天凌晨 3 点

# 通过 Web UI
http://localhost:8080/scheduler
```

**Cron 语法：**
```
┌───────────── 分钟 (0 - 59)
│ ┌───────────── 小时 (0 - 23)
│ │ ┌───────────── 日 (1 - 31)
│ │ │ ┌───────────── 月 (1 - 12)
│ │ │ │ ┌───────────── 星期 (0 - 6)（周日=0）
│ │ │ │ │
* * * * *

示例：
0 * * * *     # 每小时
0 3 * * *     # 每天凌晨 3 点
0 0 * * 0     # 每周日午夜
*/15 * * * *  # 每 15 分钟
```

**文档：** `/docs/zh/scheduler.md`

---

### 钩子和事件

**描述：** 响应系统事件的钩子系统。

**可用事件：**
- `workflow.started` - 工作流已开始
- `workflow.completed` - 工作流已完成
- `workflow.failed` - 工作流失败
- `task.started` - 任务已开始
- `task.completed` - 任务已完成
- `task.failed` - 任务失败
- `agent.connected` - 代理已连接
- `agent.disconnected` - 代理已断开

**创建钩子：**
```bash
sloth-runner hook add slack-notify \
  --event workflow.completed \
  --script /scripts/notify-slack.lua \
  --priority 10
```

**钩子脚本（Lua）：**
```lua
-- /scripts/notify-slack.lua
local event = hook.event
local payload = hook.payload

if event == "workflow.completed" then
    notifications.slack(
        "https://hooks.slack.com/services/XXX/YYY/ZZZ",
        string.format("✅ 工作流 '%s' 已完成！", payload.workflow_name),
        { channel = "#deployments" }
    )
end
```

**可用负载：**
```lua
-- workflow.* 事件
{
    workflow_name = "deploy",
    status = "success",
    duration = 45.3,
    started_at = 1234567890,
    completed_at = 1234567935
}

-- agent.* 事件
{
    agent_name = "web-01",
    address = "192.168.1.100:50060",
    status = "connected"
}
```

**文档：** `/docs/architecture/hooks-events-system.md`

---

### GitOps

**描述：** 完整的 GitOps 模式实现。

**特性：**
- **基于 Git** - Git 作为真相来源
- **自动同步** - 自动同步
- **漂移检测** - 检测手动更改
- **回滚** - 自动回滚
- **多环境** - dev、staging、production
- **基于 PR** - 通过 Pull Request 批准

**GitOps 工作流：**
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

**CLI：**
```bash
# 手动同步
sloth-runner gitops sync k8s-manifests

# 查看状态
sloth-runner gitops status

# 查看漂移
sloth-runner gitops diff
```

**文档：** `/docs/en/gitops-features.md`

---

### CI/CD 集成

**描述：** 与 CI/CD 管道集成。

**支持：**
- GitHub Actions
- GitLab CI
- Jenkins
- CircleCI
- Travis CI
- Azure Pipelines

**GitHub Actions 示例：**
```yaml
# .github/workflows/deploy.yml
name: 部署

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: 安装 Sloth Runner
        run: |
          curl -L https://github.com/org/sloth-runner/releases/latest/download/sloth-runner-linux-amd64 -o sloth-runner
          chmod +x sloth-runner

      - name: 运行部署
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

**描述：** 保存和可重用的工作流仓库。

**特性：**
- **版本控制** - 版本历史
- **标签** - 按标签组织
- **搜索** - 按名称/描述/标签搜索
- **克隆** - 克隆现有 sloth
- **导出/导入** - 共享 sloth
- **活动/非活动** - 激活/停用而不删除

**命令：**
```bash
# 添加 sloth
sloth-runner sloth add deploy --file deploy.sloth

# 列出 sloth
sloth-runner sloth list

# 查看 sloth
sloth-runner sloth get deploy

# 执行 sloth
sloth-runner run deploy --file $(sloth-runner sloth get deploy --show-path)

# 删除 sloth
sloth-runner sloth remove deploy
```

**文档：** `/docs/features/sloth-management.md`

---

## 📊 监控

### 遥测

**描述：** 完整的可观测性系统。

**组件：**
- Prometheus 指标
- 结构化日志
- 分布式跟踪
- 健康检查
- 性能分析

**架构：**
```
┌──────────┐    指标    ┌────────────┐
│  主服务器│───────────► Prometheus │
└──────────┘            └─────┬──────┘
                              │
┌──────────┐    指标          │
│ 代理 1   ├──────────────────┤
└──────────┘                  │
                              ▼
┌──────────┐    指标    ┌──────────┐
│ 代理 2   ├───────────►  Grafana │
└──────────┘            └──────────┘
```

**端点：**
```
http://master:9090/metrics
http://agent:9091/metrics
```

**文档：** `/docs/en/telemetry/index.md`

---

### Prometheus 指标

**描述：** 以 Prometheus 格式导出的指标。

**可用指标：**

#### 工作流
```
sloth_workflow_executions_total{status="success|failed"}
sloth_workflow_duration_seconds{workflow="name"}
sloth_workflow_tasks_total{workflow="name"}
```

#### 代理
```
sloth_agent_connected_total
sloth_agent_cpu_usage_percent{agent="name"}
sloth_agent_memory_usage_bytes{agent="name"}
sloth_agent_disk_usage_bytes{agent="name"}
```

#### 系统
```
sloth_tasks_executed_total
sloth_hooks_triggered_total{event="type"}
sloth_db_size_bytes
```

**抓取配置：**
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

**文档：** `/docs/en/telemetry/prometheus-metrics.md`

---

### Grafana 仪表板

**描述：** Grafana 的预配置仪表板。

**仪表板：**
1. **概览** - 系统概览
2. **代理** - 所有代理的指标
3. **工作流** - 执行和性能
4. **资源** - CPU、内存、磁盘、网络

**导入仪表板：**
```bash
# 生成仪表板 JSON
sloth-runner agent metrics grafana web-01 --export dashboard.json

# 导入到 Grafana
curl -X POST http://admin:admin@localhost:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -d @dashboard.json
```

**功能：**
- 自动刷新（5s、10s、30s、1m）
- 时间范围选择器
- 变量（代理、工作流）
- 可配置的警报
- 导出 PNG/PDF

**文档：** `/docs/en/telemetry/grafana-dashboard.md`

---

### 集中式日志

**描述：** 结构化日志的集中式系统。

**特性：**
- **结构化** - JSON 结构化日志
- **级别** - debug、info、warn、error
- **上下文** - 丰富的元数据
- **搜索** - 按任何字段搜索
- **导出** - JSON、CSV、文本
- **保留** - 保留策略

**格式：**
```json
{
  "timestamp": "2025-10-07T10:30:45Z",
  "level": "info",
  "message": "工作流已完成",
  "workflow": "deploy",
  "agent": "web-01",
  "duration": 45.3,
  "status": "success"
}
```

**访问：**
```bash
# CLI
sloth-runner logs --follow

# Web UI
http://localhost:8080/logs

# API
curl http://localhost:8080/api/v1/logs?level=error&since=1h
```

---

### 代理指标

**描述：** 实时详细的代理指标。

**收集的指标：**
- CPU 使用率（%）
- 内存使用率（字节、%）
- 磁盘使用率（字节、%）
- 负载平均值（1m、5m、15m）
- 网络 I/O（字节/秒）
- 进程数
- 运行时间

**可视化：**
```bash
# CLI
sloth-runner agent metrics web-01
sloth-runner agent metrics web-01 --watch

# Web UI - 代理仪表板
http://localhost:8080/agent-dashboard?agent=web-01

# API
curl http://localhost:8080/api/v1/agents/web-01/metrics
```

**格式：**
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

## ☁️ 云和 IaC

### 多云

**描述：** 对多个云提供商的原生支持。

**支持的提供商：**
- ✅ AWS（EC2、S3、RDS、Lambda、ECS、EKS 等）
- ✅ Azure（虚拟机、存储、AKS、Functions 等）
- ✅ GCP（Compute Engine、Cloud Storage、GKE 等）
- ✅ DigitalOcean（Droplet、Spaces、K8s 等）
- ✅ Linode
- ✅ Vultr
- ✅ Hetzner Cloud

**多云示例：**
```yaml
# 同时部署到 AWS 和 GCP
tasks:
  - name: 部署到 AWS
    exec:
      script: |
        aws.ec2_instance_create({
          image_id = "ami-xxx",
          instance_type = "t3.medium"
        })
    delegate_to: aws-agent

  - name: 部署到 GCP
    exec:
      script: |
        gcp.compute_instance_create({
          machine_type = "e2-medium",
          image_family = "ubuntu-2204-lts"
        })
    delegate_to: gcp-agent
```

**文档：** `/docs/en/enterprise-features.md`

---

### Terraform

**描述：** 与 Terraform 的完整集成。

**特性：**
- `terraform.init` - 初始化
- `terraform.plan` - 计划
- `terraform.apply` - 应用
- `terraform.destroy` - 销毁
- 状态管理
- 后端配置
- 变量文件

**示例：**
```lua
local tf_dir = "/infra/terraform"

-- 初始化
terraform.init(tf_dir, {
    backend_config = {
        bucket = "my-tf-state",
        key = "prod/terraform.tfstate"
    }
})

-- 计划
local plan = terraform.plan(tf_dir, {
    var_file = "production.tfvars",
    vars = {
        region = "us-east-1",
        environment = "production"
    }
})

-- 如果有变更则应用
if plan.changes > 0 then
    terraform.apply(tf_dir, {
        auto_approve = true
    })
end
```

**文档：** `/docs/modules/terraform.md`

---

### Pulumi

**描述：** 与 Pulumi 集成。

**支持：**
- 堆栈管理
- 配置
- 部署
- 销毁
- 预览

**示例：**
```lua
-- 选择堆栈
pulumi.stack_select("production")

-- 配置
pulumi.config_set("aws:region", "us-east-1")

-- 部署
pulumi.up({
    yes = true,  -- 自动批准
    parallel = 10
})
```

**文档：** `/docs/modules/pulumi.md`

---

### Kubernetes

**描述：** Kubernetes 部署和管理。

**特性：**
- 应用清单
- Helm 图表
- 命名空间
- ConfigMaps/Secrets
- 部署
- 健康检查

**示例：**
```lua
-- 应用清单
kubernetes.apply("/k8s/deployment.yaml", {
    namespace = "production"
})

-- Helm 安装
helm.install("myapp", "charts/myapp", {
    namespace = "production",
    values = {
        image = {
            tag = "v1.2.3"
        }
    }
})

-- 等待部署
kubernetes.rollout_status("deployment/myapp", {
    namespace = "production",
    timeout = "5m"
})
```

**文档：** `/docs/en/gitops/kubernetes.md`

---

### Docker

**描述：** 完整的 Docker 自动化。

**功能：**
- 容器生命周期（运行、停止、删除）
- 镜像管理（构建、推送、拉取）
- 网络（创建、连接）
- 卷（创建、挂载）
- Docker Compose

**部署示例：**
```lua
-- 构建镜像
docker.image_build(".", {
    tag = "myapp:v1.2.3",
    build_args = {
        VERSION = "1.2.3"
    }
})

-- 推送到注册表
docker.image_push("myapp:v1.2.3", {
    registry = "registry.example.com"
})

-- 部署
docker.container_run("myapp:v1.2.3", {
    name = "app",
    ports = {"3000:3000"},
    env = {
        DATABASE_URL = "postgres://..."
    },
    restart = "unless-stopped"
})
```

**文档：** `/docs/modules/docker.md`

---

## 🔐 安全和企业

### 身份验证

**描述：** Web UI 和 API 的身份验证系统。

**方法：**
- 用户名/密码
- JWT 令牌
- OAuth2（GitHub、Google 等）
- LDAP/AD
- SSO

**设置：**
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

**描述：** TLS/SSL 支持安全通信。

**特性：**
- gRPC TLS
- HTTPS Web UI
- 证书管理
- 自动续订（Let's Encrypt）

**配置：**
```bash
# 带 TLS 的主服务器
sloth-runner server \
  --tls-cert /etc/sloth/cert.pem \
  --tls-key /etc/sloth/key.pem

# 带 TLS 的代理
sloth-runner agent start \
  --master-tls-cert /etc/sloth/master-cert.pem
```

---

### 审计日志

**描述：** 所有操作的审计日志。

**审计的事件：**
- 用户登录/登出
- 工作流执行
- 配置更改
- API 调用
- 管理员操作

**格式：**
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

### 备份

**描述：** 自动备份系统。

**特性：**
- 可配置的自动备份
- 压缩（gzip）
- 保留策略
- 远程备份（S3、Azure Blob 等）
- 恢复

**命令：**
```bash
# 手动备份
sloth-runner db backup --output /backup/sloth.db --compress

# 恢复
sloth-runner db restore /backup/sloth.db.gz --decompress

# 自动备份（cron）
0 3 * * * sloth-runner db backup --output /backup/sloth-$(date +\%Y\%m\%d).db --compress
```

---

### RBAC

**描述：** 基于角色的访问控制。

**角色：**
- **管理员** - 完全访问
- **操作员** - 执行工作流、管理代理
- **开发者** - 创建/编辑工作流
- **查看者** - 仅查看

**权限：**
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

## 🚀 性能

### 优化

**描述：** 最近的性能优化。

**实施的改进：**

#### 代理优化
- ✅ **超低内存** - 32MB RAM 占用
- ✅ **二进制大小减少** - 从 45MB → 12MB
- ✅ **启动时间** - <100ms
- ✅ **CPU 效率** - 空闲时 99% 闲置

#### 数据库优化
- ✅ **WAL 模式** - 预写日志
- ✅ **连接池** - 连接重用
- ✅ **预编译语句** - 优化的查询
- ✅ **索引** - 关键字段上的索引
- ✅ **自动清理** - 自动清理

#### gRPC 优化
- ✅ **连接重用** - keepalive
- ✅ **压缩** - gzip 压缩
- ✅ **多路复用** - 多个流
- ✅ **缓冲池** - 缓冲区重用

**基准测试：**
```
之前：
- 代理内存：128MB
- 二进制大小：45MB
- 启动时间：2s

之后：
- 代理内存：32MB（减少 75%）
- 二进制大小：12MB（减少 73%）
- 启动时间：95ms（快 95%）
```

**文档：** `/docs/PERFORMANCE_OPTIMIZATIONS.md`

---

### 并行执行

**描述：** 使用 goroutine 并行执行任务。

**特性：**
- **goroutine.parallel()** - 并行执行函数
- **并发控制** - 限制同时 goroutine
- **错误处理** - 收集所有 goroutine 的错误
- **等待组** - 自动同步

**示例：**
```lua
-- 并行执行多个任务
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

-- 限制并发
goroutine.parallel({
    tasks = {
        function() exec.command("task1") end,
        function() exec.command("task2") end,
        function() exec.command("task3") end,
        function() exec.command("task4") end
    },
    max_concurrent = 2  -- 最多同时 2 个
})
```

**文档：** `/docs/modules/goroutine.md`

---

### 资源限制

**描述：** 可配置的资源限制。

**配置：**
```yaml
# 代理配置
resources:
  cpu:
    limit: 2  # 核心
    reserve: 0.5
  memory:
    limit: 2GB
    reserve: 512MB
  disk:
    limit: 10GB
    min_free: 1GB
```

**强制：**
- CPU 限制
- 内存限制（cgroup）
- 磁盘配额
- 任务超时

---

### 缓存

**描述：** 优化的缓存系统。

**缓存类型：**

#### 模块缓存
- 编译的 Lua 模块
- 减少加载时间

#### 状态缓存
- 内存中的状态
- 减少数据库查询

#### 指标缓存
- 聚合指标
- 减少计算

**配置：**
```yaml
cache:
  enabled: true
  ttl: 5m
  max_size: 100MB
  eviction: lru  # 最近最少使用
```

---

## 📚 其他资源

### 文档
- [🚀 快速入门](/docs/en/quick-start.md)
- [🏗️ 架构](/docs/architecture/sloth-runner-architecture.md)
- [📖 现代 DSL](/docs/modern-dsl/introduction.md)
- [🎯 高级示例](/docs/en/advanced-examples.md)

### 有用的链接
- [GitHub 仓库](https://github.com/chalkan3/sloth-runner)
- [问题跟踪](https://github.com/chalkan3/sloth-runner/issues)
- [发布](https://github.com/chalkan3/sloth-runner/releases)

---

**最后更新：** 2025-10-07

**已记录的功能总数：** 100+
