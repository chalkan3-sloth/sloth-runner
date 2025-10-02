# 🤖 Sloth Runner Agent - Cheat Sheet

## 🚀 Quick Commands

### Master Server
```bash
# Start master
./start_master.sh
# OR
./sloth-runner master start --port 50053 --bind-address 192.168.1.29 --daemon

# Stop master
pkill -f "sloth-runner master"

# Check status
ps aux | grep "sloth-runner master"
```

### Local Agent
```bash
# Start agent
./start_local_agent.sh <name> [ip] [port]
# Example:
./start_local_agent.sh my-agent

# Stop agent
pkill -f "sloth-runner agent start.*--name my-agent"
```

### Remote Agent (via SSH)
```bash
# Start
./manage_remote_agent.sh start user@host agent-name agent-ip

# Stop
./manage_remote_agent.sh stop user@host agent-name

# Status
./manage_remote_agent.sh status user@host agent-name

# Install sloth-runner
./manage_remote_agent.sh install user@host
```

### List & Manage Agents
```bash
# List all agents
./sloth-runner agent list --master 192.168.1.29:50053

# Run command on agent
./sloth-runner agent run --master 192.168.1.29:50053 \
  --agent agent-name --command "hostname"

# Stop remote agent (from master)
./sloth-runner agent stop --master 192.168.1.29:50053 \
  --agent agent-name
```

## 📝 Lua Script Examples

### Basic Task with :delegate_to()
```lua
local task = task("my_task")
    :delegate_to("agent-name")
    :command(function(this, params)
        local exec = require("exec")
        return exec.run("hostname")
    end)
    :build()

task:run()
```

### Multiple Agents
```lua
local agents = {"agent1", "agent2", "agent3"}

for _, agent in ipairs(agents) do
    local t = task("check_" .. agent)
        :delegate_to(agent)
        :command(function()
            local exec = require("exec")
            return exec.run("uptime")
        end)
        :build()
    t:run()
end
```

### With Modules
```lua
local deploy = task("deploy_nginx")
    :delegate_to("web-server")
    :command(function()
        local pkg = require("pkg")
        local systemd = require("systemd")
        
        -- Install nginx
        pkg.install({"nginx"})
        
        -- Start and enable
        systemd.enable("nginx")
        systemd.start("nginx")
        
        return true, "Deployed successfully"
    end)
    :timeout("2m")
    :build()
```

## 🎯 Configuration

### Your Setup (Example)
- **Master**: 192.168.1.29:50053
- **Agent ladyguica**: 192.168.1.16:50051
- **Agent keiteguica**: 192.168.1.17:50051

### Ports
- **Master**: 50053 (default)
- **Agent**: 50051 (default)

### Files
- **Executável**: `~/.local/bin/sloth-runner`
- **Logs**: Check with `--daemon` off
- **Config**: Via command-line flags

## 🔧 Troubleshooting

### Agent não conecta?
1. ✅ Master running? `ps aux | grep master`
2. ✅ Network? `ping <master-ip>`
3. ✅ Port open? `telnet <master-ip> 50053`
4. ✅ Firewall? `sudo ufw allow 50053/tcp`

### Ver logs
```bash
# Run without daemon
./sloth-runner agent start --name test \
  --master 192.168.1.29:50053 \
  --bind-address 192.168.1.X
```

### Kill all
```bash
# Kill all sloth-runner processes
pkill -f sloth-runner
```

## 📚 Resources

- **Full Guide**: `docs/agent-setup.md`
- **Quick Start**: `AGENT_QUICK_START.md`
- **Examples**: `examples/distributed_execution.sloth`
- **README**: Distributed Task Execution section

## 🌟 Common Use Cases

### 1. Deploy to all web servers
```lua
local web_servers = {"web1", "web2", "web3"}
-- Deploy nginx to each
```

### 2. Database migration on specific server
```lua
local migrate = task("migrate")
    :delegate_to("db-primary")
    :command(function() ... end)
```

### 3. Health check all agents
```lua
-- Check service status on all agents
for _, agent in ipairs(agents) do
    check_health(agent)
end
```

### 4. Backup configurations
```lua
local backup = task("backup")
    :delegate_to("app-server")
    :command(function()
        local exec = require("exec")
        return exec.run("tar -czf /backup/config.tar.gz /etc")
    end)
```

---

💡 **Tip**: Always test with a local agent first before deploying to remote agents!

🔗 **Links**:
- GitHub: https://github.com/chalkan3-sloth/sloth-runner
- Docs: `docs/agent-setup.md`
