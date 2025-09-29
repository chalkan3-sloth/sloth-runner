# ⚡ Quick Start Guide

Get up and running with Sloth Runner in under 10 minutes! This guide will walk you through installation, basic usage, and your first distributed task execution.

## 🚀 **Installation**

### **Option 1: Download Binary**
```bash
# Download latest release
curl -L https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner-linux-amd64 -o sloth-runner
chmod +x sloth-runner
sudo mv sloth-runner /usr/local/bin/
```

### **Option 2: Build from Source**
```bash
# Clone repository
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Build binary
go build -o sloth-runner ./cmd/sloth-runner

# Add to PATH
export PATH=$PATH:$(pwd)
```

### **Option 3: Docker**
```bash
# Pull official image
docker pull slothrunner/sloth-runner:latest

# Create alias for easy usage
alias sloth-runner='docker run --rm -v $(pwd):/workspace slothrunner/sloth-runner'
```

## 📋 **Verify Installation**

```bash
# Check version
sloth-runner --version

# View available commands
sloth-runner --help
```

Expected output:
```
Sloth Runner v2.0.0
A powerful task orchestration platform with Lua scripting
```

## 🎯 **Your First Task**

Create your first Lua task file:

```bash
# Create a simple task file
cat > hello-world.lua << 'EOF'
Modern DSLs = {
    hello_world = {
        description = "My first Sloth Runner task",
        tasks = {
            greet = {
                name = "greet",
                description = "Say hello to the world",
                command = function()
                    log.info("🎉 Hello from Sloth Runner!")
                    
                    -- Get system information
                    local hostname, _ = exec.run("hostname")
                    local whoami, _ = exec.run("whoami")
                    
                    log.info("Running on: " .. hostname)
                    log.info("User: " .. whoami)
                    
                    -- Use state management
                    state.set("last_greeting", os.time())
                    local count = state.increment("greeting_count", 1)
                    
                    log.info("This is greeting #" .. count)
                    
                    return true, "Hello World task completed successfully!"
                end
            },
            
            system_info = {
                name = "system_info", 
                description = "Display system metrics",
                depends_on = "greet",
                command = function()
                    log.info("📊 System Information:")
                    
                    -- Get system metrics
                    local cpu = metrics.system_cpu()
                    local memory = metrics.system_memory()
                    local disk = metrics.system_disk()
                    
                    log.info("CPU Usage: " .. string.format("%.1f%%", cpu))
                    log.info("Memory: " .. string.format("%.1f%% (%.0f MB used)", 
                        memory.percent, memory.used_mb))
                    log.info("Disk: " .. string.format("%.1f%% (%.1f GB used)", 
                        disk.percent, disk.used_gb))
                    
                    -- Record metrics
                    metrics.gauge("quickstart_cpu", cpu)
                    metrics.gauge("quickstart_memory", memory.percent)
                    
                    return true, "System info collected"
                end
            }
        }
    }
}
EOF
```

## 🏃‍♂️ **Run Your First Task**

```bash
# Execute the task
sloth-runner run -f hello-world.lua

# Or run specific task
sloth-runner run -f hello-world.lua -t greet
```

Expected output:
```
2024-01-15 10:30:00 INFO 🎉 Hello from Sloth Runner!
2024-01-15 10:30:00 INFO Running on: my-computer
2024-01-15 10:30:00 INFO User: myuser
2024-01-15 10:30:00 INFO This is greeting #1
2024-01-15 10:30:01 INFO 📊 System Information:
2024-01-15 10:30:01 INFO CPU Usage: 15.2%
2024-01-15 10:30:01 INFO Memory: 45.8% (7520 MB used)
2024-01-15 10:30:01 INFO Disk: 67.3% (234.5 GB used)
✅ Task 'hello_world' completed successfully!
```

## 🌐 **Setting Up Distributed Execution**

### **Step 1: Start Master Server**
```bash
# Start master on your main machine (e.g., 192.168.1.100)
sloth-runner master --port 50053 --bind-address 192.168.1.100

# Or with enhanced features
sloth-runner master --port 50053 --metrics-port 8080 --dashboard-port 3000
```

### **Step 2: Deploy Remote Agents**

On remote machine 1 (192.168.1.101):
```bash
# Download sloth-runner binary to remote machine
scp sloth-runner user@192.168.1.101:/usr/local/bin/

# SSH and start agent
ssh user@192.168.1.101
sloth-runner agent start \
    --name agent-1 \
    --master 192.168.1.100:50053 \
    --port 50051 \
    --bind-address 192.168.1.101
```

On remote machine 2 (192.168.1.102):
```bash
# SSH and start agent  
ssh user@192.168.1.102
sloth-runner agent start \
    --name agent-2 \
    --master 192.168.1.100:50053 \
    --port 50051 \
    --bind-address 192.168.1.102
```

### **Step 3: Verify Agent Registration**
```bash
# List registered agents
sloth-runner agent list --master 192.168.1.100:50053
```

Expected output:
```
Registered Agents:
  agent-1    192.168.1.101:50051    Active    2s ago
  agent-2    192.168.1.102:50051    Active    1s ago
```

### **Step 4: Run Distributed Tasks**
```bash
# Execute command on specific agent
sloth-runner agent run agent-1 "echo 'Hello from Agent 1'" --master 192.168.1.100:50053

# Execute on all agents
sloth-runner agent run agent-1 "uptime" --master 192.168.1.100:50053 &
sloth-runner agent run agent-2 "uptime" --master 192.168.1.100:50053 &
wait
```

## 📊 **Exploring Advanced Features**

### **State Management Example**

```lua
-- Create state-demo.lua
Modern DSLs = {
    state_demo = {
        description = "Demonstrate state management capabilities",
        tasks = {
            setup_state = {
                name = "setup_state",
                description = "Initialize application state", 
                command = function()
                    -- Initialize configuration
                    state.set("app_config", {
                        version = "1.0.0",
                        environment = "development",
                        debug = true
                    })
                    
                    -- Set TTL for session data (5 minutes)
                    state.set("session_token", "abc123xyz", 300)
                    
                    -- Initialize counters
                    state.set("api_calls", 0)
                    state.set("errors", 0)
                    
                    log.info("✅ Application state initialized")
                    return true, "State setup completed"
                end
            },
            
            simulate_usage = {
                name = "simulate_usage",
                description = "Simulate application usage",
                depends_on = "setup_state",
                command = function()
                    -- Simulate API calls
                    for i = 1, 10 do
                        local calls = state.increment("api_calls", 1)
                        
                        -- Simulate occasional error
                        if math.random(1, 10) > 8 then
                            state.increment("errors", 1)
                            log.warn("Simulated error occurred")
                        end
                        
                        -- Add to processing queue
                        state.list_push("processing_queue", {
                            id = "req_" .. i,
                            timestamp = os.time(),
                            status = "pending"
                        })
                        
                        exec.run("sleep 0.1") -- Small delay
                    end
                    
                    local total_calls = state.get("api_calls")
                    local total_errors = state.get("errors")
                    local queue_size = state.list_length("processing_queue")
                    
                    log.info("📊 Usage Summary:")
                    log.info("  API Calls: " .. total_calls)
                    log.info("  Errors: " .. total_errors)
                    log.info("  Queue Size: " .. queue_size)
                    
                    return true, "Usage simulation completed"
                end
            },
            
            process_queue = {
                name = "process_queue",
                description = "Process items in queue with locking",
                depends_on = "simulate_usage",
                command = function()
                    -- Process queue with distributed lock
                    state.with_lock("queue_processing", function()
                        log.info("🔒 Processing queue with exclusive lock...")
                        
                        local processed = 0
                        while state.list_length("processing_queue") > 0 do
                            local item = state.list_pop("processing_queue")
                            log.info("Processing item: " .. item.id)
                            processed = processed + 1
                        end
                        
                        log.info("✅ Processed " .. processed .. " items")
                        state.set("last_processing_time", os.time())
                        
                    end, 30) -- 30 second timeout
                    
                    return true, "Queue processing completed"
                end
            }
        }
    }
}
```

Run the state demo:
```bash
sloth-runner run -f state-demo.lua
```

### **Metrics Monitoring Example**

```lua
-- Create metrics-demo.lua  
Modern DSLs = {
    metrics_demo = {
        description = "Demonstrate metrics and monitoring",
        tasks = {
            collect_metrics = {
                name = "collect_metrics",
                description = "Collect system and custom metrics",
                command = function()
                    log.info("📊 Collecting system metrics...")
                    
                    -- System metrics
                    local cpu = metrics.system_cpu()
                    local memory = metrics.system_memory() 
                    local disk = metrics.system_disk()
                    
                    log.info("System Status:")
                    log.info("  CPU: " .. string.format("%.1f%%", cpu))
                    log.info("  Memory: " .. string.format("%.1f%%", memory.percent))
                    log.info("  Disk: " .. string.format("%.1f%%", disk.percent))
                    
                    -- Custom metrics
                    metrics.gauge("demo_cpu_usage", cpu)
                    metrics.counter("demo_executions", 1)
                    
                    -- Performance timer
                    local processing_time = metrics.timer("data_processing", function()
                        -- Simulate data processing
                        local sum = 0
                        for i = 1, 1000000 do
                            sum = sum + math.sqrt(i)
                        end
                        return sum
                    end)
                    
                    log.info("⏱️ Processing took: " .. string.format("%.2f ms", processing_time))
                    
                    -- Health check
                    local health = metrics.health_status()
                    log.info("🏥 Overall health: " .. health.overall)
                    
                    -- Alert if CPU is high
                    if cpu > 50 then
                        metrics.alert("high_cpu_demo", {
                            level = "warning",
                            message = "CPU usage is elevated: " .. string.format("%.1f%%", cpu),
                            value = cpu
                        })
                    end
                    
                    return true, "Metrics collection completed"
                end
            }
        }
    }
}
```

Run the metrics demo:
```bash
sloth-runner run -f metrics-demo.lua
```

## 🌐 **Access Web Dashboard**

If you started the master with dashboard support:

```bash
# Open web dashboard
open http://192.168.1.100:3000

# View metrics endpoint
curl http://192.168.1.100:8080/metrics

# Check health status
curl http://192.168.1.100:8080/health
```

## 📚 **What's Next?**

### **Explore Core Concepts**
- 📖 [Core Concepts](core-concepts.md) - Understand tasks, workflows, and state
- 🔧 [CLI Commands](CLI.md) - Master all available commands
- 🌙 [Lua API](../LUA_API.md) - Deep dive into scripting capabilities

### **Advanced Features**
- 💾 [State Management](modules/state.md) - Persistent state and locks
- 📊 [Metrics & Monitoring](modules/metrics.md) - Observability and alerting
- 🚀 [Agent Improvements](agent-improvements.md) - Enterprise features

### **Cloud Integrations**
- ☁️ [AWS Integration](modules/aws.md) - Deploy and manage AWS resources
- 🌩️ [GCP Integration](modules/gcp.md) - Google Cloud Platform tasks
- 🔷 [Azure Integration](modules/azure.md) - Microsoft Azure automation

### **Infrastructure as Code**
- 🐳 [Docker](modules/docker.md) - Container management
- 🏗️ [Pulumi](modules/pulumi.md) - Modern infrastructure as code
- 🌍 [Terraform](modules/terraform.md) - Infrastructure provisioning

## 🆘 **Getting Help**

### **Documentation**
- 📚 [Full Documentation](../index.md)
- 🔍 [API Reference](../LUA_API.md)
- 💡 [Examples](../EXAMPLES.md)

### **Community**
- 💬 [GitHub Discussions](https://github.com/chalkan3-sloth/sloth-runner/discussions)
- 🐛 [Issue Tracker](https://github.com/chalkan3-sloth/sloth-runner/issues)
- 📧 [Email Support](mailto:support@sloth-runner.dev)

### **Quick Troubleshooting**

**Agent won't connect to master?**
```bash
# Check network connectivity
telnet 192.168.1.100 50053

# Verify master is running
sloth-runner agent list --master 192.168.1.100:50053

# Check firewall settings
sudo ufw status
```

**Tasks failing with permission errors?**
```bash
# Check user permissions
ls -la /usr/local/bin/sloth-runner

# Run with appropriate user
sudo -u myuser sloth-runner run -f task.lua
```

**State database issues?**
```bash
# Check state database location
ls -la ~/.sloth-runner/

# View state statistics
sloth-runner state stats

# Clear corrupted state (careful!)
rm ~/.sloth-runner/state.db*
```

## 🎉 **Congratulations!**

You've successfully:
- ✅ Installed Sloth Runner
- ✅ Executed your first task
- ✅ Set up distributed agents
- ✅ Explored state management
- ✅ Monitored system metrics

You're now ready to build powerful, distributed task orchestration workflows with Sloth Runner! 🚀

<div class="hero">
  <h2>🚀 Ready for More?</h2>
  <p>Explore advanced features and build production-ready workflows</p>
  <a href="advanced-features.md" class="btn">Advanced Features →</a>
  <a href="advanced-examples.md" class="btn">More Examples →</a>
</div>