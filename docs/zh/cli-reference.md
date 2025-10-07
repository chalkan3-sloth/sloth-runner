# 📚 CLI 命令完整参考

## 概述

Sloth Runner 提供了一个完整而强大的命令行界面（CLI），用于管理工作流、代理、模块、钩子、事件等。本文档涵盖了**所有**可用命令及实际示例。

---

## 🎯 主要命令

### `run` - 执行工作流

从文件执行 Sloth 工作流。

```bash
# 基本语法
sloth-runner run <workflow-name> --file <文件.sloth> [选项]

# 示例
sloth-runner run deploy --file deploy.sloth
sloth-runner run deploy --file deploy.sloth --yes                    # 非交互模式
sloth-runner run deploy --file deploy.sloth --group production       # 执行特定组
sloth-runner run deploy --file deploy.sloth --delegate-to agent1     # 委托给代理
sloth-runner run deploy --file deploy.sloth --delegate-to agent1 --delegate-to agent2  # 多个代理
sloth-runner run deploy --file deploy.sloth --values vars.yaml       # 传递变量
sloth-runner run deploy --file deploy.sloth --var "env=production"   # 内联变量
```

**选项：**
- `--file, -f` - Sloth 文件路径
- `--yes, -y` - 非交互模式（不询问确认）
- `--group, -g` - 仅执行特定组
- `--delegate-to` - 将执行委托给远程代理
- `--values` - 包含变量的 YAML 文件
- `--var` - 定义内联变量（可多次使用）
- `--verbose, -v` - 详细模式

---

## 🤖 代理管理

### `agent list` - 列出代理

列出在主服务器上注册的所有代理。

```bash
# 语法
sloth-runner agent list [选项]

# 示例
sloth-runner agent list                    # 列出所有代理
sloth-runner agent list --format json      # JSON 输出
sloth-runner agent list --format yaml      # YAML 输出
sloth-runner agent list --status active    # 仅活动代理
```

**选项：**
- `--format` - 输出格式：table（默认）、json、yaml
- `--status` - 按状态过滤：active、inactive、all

---

### `agent get` - 代理详情

获取特定代理的详细信息。

```bash
# 语法
sloth-runner agent get <agent-name> [选项]

# 示例
sloth-runner agent get web-server-01
sloth-runner agent get web-server-01 --format json
sloth-runner agent get web-server-01 --show-metrics       # 包含指标
```

**选项：**
- `--format` - 输出格式：table、json、yaml
- `--show-metrics` - 显示代理指标

---

### `agent install` - 安装远程代理

通过 SSH 在远程服务器上安装 Sloth Runner 代理。

```bash
# 语法
sloth-runner agent install <agent-name> --ssh-host <host> --ssh-user <user> [选项]

# 示例
sloth-runner agent install web-01 --ssh-host 192.168.1.100 --ssh-user root
sloth-runner agent install web-01 --ssh-host 192.168.1.100 --ssh-user root --ssh-port 2222
sloth-runner agent install web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user root \
  --master 192.168.1.1:50053 \
  --bind-address 0.0.0.0 \
  --port 50060 \
  --report-address 192.168.1.100:50060
```

**选项：**
- `--ssh-host` - 远程服务器的 SSH 主机（必需）
- `--ssh-user` - SSH 用户（必需）
- `--ssh-port` - SSH 端口（默认：22）
- `--ssh-key` - SSH 私钥路径
- `--master` - 主服务器地址（默认：localhost:50053）
- `--bind-address` - 代理绑定地址（默认：0.0.0.0）
- `--port` - 代理端口（默认：50060）
- `--report-address` - 代理向主服务器报告的地址

---

### `agent update` - 更新代理

将代理二进制文件更新到最新版本。

```bash
# 语法
sloth-runner agent update <agent-name> [选项]

# 示例
sloth-runner agent update web-01
sloth-runner agent update web-01 --version v1.2.3
sloth-runner agent update web-01 --restart           # 更新后重启
```

**选项：**
- `--version` - 特定版本（默认：latest）
- `--restart` - 更新后重启代理
- `--force` - 即使版本相同也强制更新

---

### `agent modules` - 代理模块

列出或检查代理上可用的模块。

```bash
# 语法
sloth-runner agent modules <agent-name> [选项]

# 示例
sloth-runner agent modules web-01                      # 列出所有模块
sloth-runner agent modules web-01 --check pkg          # 检查 'pkg' 模块是否可用
sloth-runner agent modules web-01 --check docker      # 检查是否安装了 Docker
sloth-runner agent modules web-01 --format json       # JSON 输出
```

**选项：**
- `--check` - 检查特定模块
- `--format` - 输出格式：table、json、yaml

---

### `agent start` - 启动代理

在本地启动代理服务。

```bash
# 语法
sloth-runner agent start [选项]

# 示例
sloth-runner agent start                                    # 使用默认配置启动
sloth-runner agent start --master 192.168.1.1:50053         # 连接到特定主服务器
sloth-runner agent start --port 50060                       # 使用特定端口
sloth-runner agent start --name my-agent                    # 定义代理名称
sloth-runner agent start --bind 0.0.0.0                     # 绑定所有接口
sloth-runner agent start --foreground                       # 前台运行
```

**选项：**
- `--master` - 主服务器地址（默认：localhost:50053）
- `--port` - 代理端口（默认：50060）
- `--name` - 代理名称（默认：主机名）
- `--bind` - 绑定地址（默认：0.0.0.0）
- `--report-address` - 代理报告的地址
- `--foreground` - 前台运行（非守护进程）

---

### `agent stop` - 停止代理

停止代理服务。

```bash
# 语法
sloth-runner agent stop [选项]

# 示例
sloth-runner agent stop                # 停止本地代理
sloth-runner agent stop --name web-01  # 停止特定代理
```

---

### `agent restart` - 重启代理

重启代理服务。

```bash
# 语法
sloth-runner agent restart [agent-name]

# 示例
sloth-runner agent restart               # 重启本地代理
sloth-runner agent restart web-01        # 重启远程代理
```

---

### `agent metrics` - 代理指标

查看代理的性能和资源指标。

```bash
# 语法
sloth-runner agent metrics <agent-name> [选项]

# 示例
sloth-runner agent metrics web-01
sloth-runner agent metrics web-01 --format json
sloth-runner agent metrics web-01 --watch              # 持续更新
sloth-runner agent metrics web-01 --interval 5         # 5 秒间隔
```

**选项：**
- `--format` - 格式：table、json、yaml、prometheus
- `--watch` - 持续更新
- `--interval` - 更新间隔（秒）（默认：2）

---

### `agent metrics grafana` - Grafana 仪表板

为代理生成并显示 Grafana 仪表板。

```bash
# 语法
sloth-runner agent metrics grafana <agent-name> [选项]

# 示例
sloth-runner agent metrics grafana web-01
sloth-runner agent metrics grafana web-01 --export dashboard.json
```

**选项：**
- `--export` - 将仪表板导出到 JSON 文件

---

## 📦 Sloth 管理（已保存的工作流）

### `sloth list` - 列出 Sloth

列出本地仓库中保存的所有工作流。

```bash
# 语法
sloth-runner sloth list [选项]

# 示例
sloth-runner sloth list                   # 列出所有
sloth-runner sloth list --active          # 仅活动的 sloth
sloth-runner sloth list --inactive        # 仅非活动的 sloth
sloth-runner sloth list --format json     # JSON 输出
```

**选项：**
- `--active` - 仅活动的 sloth
- `--inactive` - 仅非活动的 sloth
- `--format` - 格式：table、json、yaml

---

### `sloth add` - 添加 Sloth

将新工作流添加到仓库。

```bash
# 语法
sloth-runner sloth add <name> --file <路径> [选项]

# 示例
sloth-runner sloth add deploy --file deploy.sloth
sloth-runner sloth add deploy --file deploy.sloth --description "生产部署"
sloth-runner sloth add deploy --file deploy.sloth --tags "prod,deploy"
```

**选项：**
- `--file` - Sloth 文件路径（必需）
- `--description` - Sloth 描述
- `--tags` - 逗号分隔的标签

---

### `sloth get` - 获取 Sloth

显示特定 sloth 的详细信息。

```bash
# 语法
sloth-runner sloth get <name> [选项]

# 示例
sloth-runner sloth get deploy
sloth-runner sloth get deploy --format json
sloth-runner sloth get deploy --show-content    # 显示工作流内容
```

**选项：**
- `--format` - 格式：table、json、yaml
- `--show-content` - 显示完整的工作流内容

---

### `sloth update` - 更新 Sloth

更新现有的 sloth。

```bash
# 语法
sloth-runner sloth update <name> [选项]

# 示例
sloth-runner sloth update deploy --file deploy-v2.sloth
sloth-runner sloth update deploy --description "新描述"
sloth-runner sloth update deploy --tags "prod,deploy,updated"
```

**选项：**
- `--file` - 新的 Sloth 文件
- `--description` - 新描述
- `--tags` - 新标签

---

### `sloth remove` - 删除 Sloth

从仓库中删除 sloth。

```bash
# 语法
sloth-runner sloth remove <name>

# 示例
sloth-runner sloth remove deploy
sloth-runner sloth remove deploy --force    # 不确认直接删除
```

**选项：**
- `--force` - 不询问确认直接删除

---

### `sloth activate` - 激活 Sloth

激活已停用的 sloth。

```bash
# 语法
sloth-runner sloth activate <name>

# 示例
sloth-runner sloth activate deploy
```

---

### `sloth deactivate` - 停用 Sloth

停用 sloth（不删除，仅标记为非活动）。

```bash
# 语法
sloth-runner sloth deactivate <name>

# 示例
sloth-runner sloth deactivate deploy
```

---

## 🎣 钩子管理

### `hook list` - 列出钩子

列出所有已注册的钩子。

```bash
# 语法
sloth-runner hook list [选项]

# 示例
sloth-runner hook list
sloth-runner hook list --format json
sloth-runner hook list --event workflow.started    # 按事件过滤
```

**选项：**
- `--format` - 格式：table、json、yaml
- `--event` - 按事件类型过滤

---

### `hook add` - 添加钩子

添加新钩子。

```bash
# 语法
sloth-runner hook add <name> --event <事件> --script <路径> [选项]

# 示例
sloth-runner hook add notify-slack --event workflow.completed --script notify.sh
sloth-runner hook add backup --event task.completed --script backup.lua
sloth-runner hook add validate --event workflow.started --script validate.lua --priority 10
```

**选项：**
- `--event` - 事件类型（必需）
- `--script` - 脚本路径（必需）
- `--priority` - 执行优先级（默认：0）
- `--enabled` - 启用钩子（默认：true）

**可用事件：**
- `workflow.started`
- `workflow.completed`
- `workflow.failed`
- `task.started`
- `task.completed`
- `task.failed`
- `agent.connected`
- `agent.disconnected`

---

### `hook remove` - 删除钩子

删除钩子。

```bash
# 语法
sloth-runner hook remove <name>

# 示例
sloth-runner hook remove notify-slack
sloth-runner hook remove notify-slack --force
```

---

### `hook enable` - 启用钩子

启用已禁用的钩子。

```bash
# 语法
sloth-runner hook enable <name>

# 示例
sloth-runner hook enable notify-slack
```

---

### `hook disable` - 禁用钩子

禁用钩子。

```bash
# 语法
sloth-runner hook disable <name>

# 示例
sloth-runner hook disable notify-slack
```

---

### `hook test` - 测试钩子

测试钩子的执行。

```bash
# 语法
sloth-runner hook test <name> [选项]

# 示例
sloth-runner hook test notify-slack
sloth-runner hook test notify-slack --payload '{"message": "test"}'
```

**选项：**
- `--payload` - 测试数据 JSON

---

## 📡 事件管理

### `events list` - 列出事件

列出系统最近的事件。

```bash
# 语法
sloth-runner events list [选项]

# 示例
sloth-runner events list
sloth-runner events list --limit 50               # 最近 50 个事件
sloth-runner events list --type workflow.started  # 按类型过滤
sloth-runner events list --since 1h               # 最近一小时的事件
sloth-runner events list --format json
```

**选项：**
- `--limit` - 最大事件数（默认：100）
- `--type` - 按事件类型过滤
- `--since` - 按时间过滤（例如：1h、30m、24h）
- `--format` - 格式：table、json、yaml

---

### `events watch` - 监视事件

实时监视事件。

```bash
# 语法
sloth-runner events watch [选项]

# 示例
sloth-runner events watch
sloth-runner events watch --type workflow.completed    # 仅工作流完成事件
sloth-runner events watch --filter "status=success"    # 带过滤器
```

**选项：**
- `--type` - 按事件类型过滤
- `--filter` - 过滤表达式

---

## 🗄️ 数据库管理

### `db backup` - 备份数据库

创建 SQLite 数据库备份。

```bash
# 语法
sloth-runner db backup [选项]

# 示例
sloth-runner db backup
sloth-runner db backup --output /backup/sloth-backup.db
sloth-runner db backup --compress                     # 使用 gzip 压缩
```

**选项：**
- `--output` - 备份文件路径
- `--compress` - 压缩备份

---

### `db restore` - 恢复数据库

从备份恢复数据库。

```bash
# 语法
sloth-runner db restore <备份文件> [选项]

# 示例
sloth-runner db restore /backup/sloth-backup.db
sloth-runner db restore /backup/sloth-backup.db.gz --decompress
```

**选项：**
- `--decompress` - 解压 gzip 备份

---

### `db vacuum` - 优化数据库

优化和压缩 SQLite 数据库。

```bash
# 语法
sloth-runner db vacuum

# 示例
sloth-runner db vacuum
```

---

### `db stats` - 数据库统计

显示数据库统计信息。

```bash
# 语法
sloth-runner db stats [选项]

# 示例
sloth-runner db stats
sloth-runner db stats --format json
```

**选项：**
- `--format` - 格式：table、json、yaml

---

## 🌐 SSH 管理

### `ssh list` - 列出 SSH 连接

列出已保存的 SSH 连接。

```bash
# 语法
sloth-runner ssh list [选项]

# 示例
sloth-runner ssh list
sloth-runner ssh list --format json
```

**选项：**
- `--format` - 格式：table、json、yaml

---

### `ssh add` - 添加 SSH 连接

添加新的 SSH 连接。

```bash
# 语法
sloth-runner ssh add <name> --host <host> --user <user> [选项]

# 示例
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu --port 2222
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu --key ~/.ssh/id_rsa
```

**选项：**
- `--host` - SSH 主机（必需）
- `--user` - SSH 用户（必需）
- `--port` - SSH 端口（默认：22）
- `--key` - SSH 私钥路径

---

### `ssh remove` - 删除 SSH 连接

删除已保存的 SSH 连接。

```bash
# 语法
sloth-runner ssh remove <name>

# 示例
sloth-runner ssh remove web-server
```

---

### `ssh test` - 测试 SSH 连接

测试 SSH 连接。

```bash
# 语法
sloth-runner ssh test <name>

# 示例
sloth-runner ssh test web-server
```

---

## 📋 模块

### `modules list` - 列出模块

列出所有可用模块。

```bash
# 语法
sloth-runner modules list [选项]

# 示例
sloth-runner modules list
sloth-runner modules list --format json
sloth-runner modules list --category cloud         # 按类别过滤
```

**选项：**
- `--format` - 格式：table、json、yaml
- `--category` - 按类别过滤

---

### `modules info` - 模块信息

显示模块的详细信息。

```bash
# 语法
sloth-runner modules info <module-name>

# 示例
sloth-runner modules info pkg
sloth-runner modules info docker
sloth-runner modules info terraform
```

---

## 🖥️ 服务器和 UI

### `server` - 启动主服务器

启动主服务器（gRPC）。

```bash
# 语法
sloth-runner server [选项]

# 示例
sloth-runner server                          # 在默认端口启动（50053）
sloth-runner server --port 50053             # 指定端口
sloth-runner server --bind 0.0.0.0           # 绑定所有接口
sloth-runner server --tls-cert cert.pem --tls-key key.pem  # 使用 TLS
```

**选项：**
- `--port` - 服务器端口（默认：50053）
- `--bind` - 绑定地址（默认：0.0.0.0）
- `--tls-cert` - TLS 证书
- `--tls-key` - TLS 私钥

---

### `ui` - 启动 Web UI

启动 Web 界面。

```bash
# 语法
sloth-runner ui [选项]

# 示例
sloth-runner ui                              # 在默认端口启动（8080）
sloth-runner ui --port 8080                  # 指定端口
sloth-runner ui --bind 0.0.0.0               # 绑定所有接口
```

**选项：**
- `--port` - Web UI 端口（默认：8080）
- `--bind` - 绑定地址（默认：0.0.0.0）

---

### `terminal` - 交互式终端

打开远程代理的交互式终端。

```bash
# 语法
sloth-runner terminal <agent-name>

# 示例
sloth-runner terminal web-01
```

---

## 🔧 实用工具

### `version` - 版本

显示 Sloth Runner 版本。

```bash
# 语法
sloth-runner version

# 示例
sloth-runner version
sloth-runner version --format json
```

---

### `completion` - 自动补全

为 shell 生成自动补全脚本。

```bash
# 语法
sloth-runner completion <shell>

# 示例
sloth-runner completion bash > /etc/bash_completion.d/sloth-runner
sloth-runner completion zsh > ~/.zsh/completion/_sloth-runner
sloth-runner completion fish > ~/.config/fish/completions/sloth-runner.fish
```

**支持的 shell：** bash、zsh、fish、powershell

---

### `doctor` - 诊断

执行系统和配置诊断。

```bash
# 语法
sloth-runner doctor [选项]

# 示例
sloth-runner doctor
sloth-runner doctor --format json
sloth-runner doctor --verbose             # 详细输出
```

**选项：**
- `--format` - 格式：text、json
- `--verbose` - 详细输出

---

## 🔐 环境变量

Sloth Runner 使用以下环境变量：

```bash
# 主服务器地址
export SLOTH_RUNNER_MASTER_ADDR="192.168.1.1:50053"

# 代理端口
export SLOTH_RUNNER_AGENT_PORT="50060"

# Web UI 端口
export SLOTH_RUNNER_UI_PORT="8080"

# 数据库路径
export SLOTH_RUNNER_DB_PATH="~/.sloth-runner/sloth.db"

# 日志级别
export SLOTH_RUNNER_LOG_LEVEL="info"  # debug, info, warn, error

# 启用调试模式
export SLOTH_RUNNER_DEBUG="true"
```

---

## 📊 常见使用示例

### 1. 生产部署与委托

```bash
sloth-runner run production-deploy \
  --file deployments/prod.sloth \
  --delegate-to web-01 \
  --delegate-to web-02 \
  --values prod-vars.yaml \
  --yes
```

### 2. 监视所有代理的指标

```bash
# 在一个终端中
sloth-runner agent metrics web-01 --watch

# 在另一个终端中
sloth-runner agent metrics web-02 --watch
```

### 3. 自动化备份

```bash
# 创建带时间戳的压缩备份
sloth-runner db backup \
  --output /backup/sloth-$(date +%Y%m%d-%H%M%S).db \
  --compress
```

### 4. 带通知钩子的工作流

```bash
# 添加通知钩子
sloth-runner hook add slack-notify \
  --event workflow.completed \
  --script /scripts/notify-slack.lua

# 执行工作流（钩子将自动触发）
sloth-runner run deploy --file deploy.sloth --yes
```

### 5. 在多个服务器上安装代理

```bash
# 循环在多个主机上安装
for host in 192.168.1.{10..20}; do
  sloth-runner agent install "agent-$host" \
    --ssh-host "$host" \
    --ssh-user ubuntu \
    --master 192.168.1.1:50053
done
```

---

## 🎓 下一步

- [📖 模块指南](modules-complete.md) - 所有模块的完整文档
- [🎨 Web UI](web-ui-complete.md) - Web 界面完整指南
- [🎯 高级示例](../en/advanced-examples.md) - 工作流实际示例
- [🏗️ 架构](../architecture/sloth-runner-architecture.md) - 系统架构

---

## 💡 提示和技巧

### 有用的别名

添加到您的 `.bashrc` 或 `.zshrc`：

```bash
alias sr='sloth-runner'
alias sra='sloth-runner agent'
alias srr='sloth-runner run'
alias srl='sloth-runner sloth list'
alias srui='sloth-runner ui --port 8080'
```

### 自动补全

```bash
# Bash
sloth-runner completion bash > /etc/bash_completion.d/sloth-runner
source /etc/bash_completion.d/sloth-runner

# Zsh
sloth-runner completion zsh > ~/.zsh/completion/_sloth-runner
```

### 调试模式

```bash
export SLOTH_RUNNER_DEBUG=true
export SLOTH_RUNNER_LOG_LEVEL=debug
sloth-runner run deploy --file deploy.sloth --verbose
```

---

**最后更新：** 2025-10-07
