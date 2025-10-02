# Agent Examples

This directory contains examples of sloth-runner workflows that demonstrate distributed task execution across multiple agents.

## Available Examples

### 1. `simple_ls.sloth`
A minimal example that runs `ls -la` on both agents:
- **ladyguica** (192.168.1.16)
- **keiteguica** (192.168.1.17)

**Usage:**
```bash
./sloth-runner run -f examples/agents/simple_ls.sloth simple_ls_workflow
```

### 2. `list_files_per_host.sloth`
A comprehensive example with multiple workflows:

#### Workflows:
- **list_files_workflow**: Lists files in `/home/chalkan3` on both hosts
- **system_info_workflow**: Collects system information (hostname, uptime, disk, memory)
- **complete_host_audit**: Runs all tasks in parallel for complete host audit

#### Tasks:
- `ls_ladyguica`: File listing on ladyguica
- `ls_keiteguica`: File listing on keiteguica  
- `sysinfo_ladyguica`: System info from ladyguica
- `sysinfo_keiteguica`: System info from keiteguica

**Usage:**
```bash
# List files on both hosts
./sloth-runner run -f examples/agents/list_files_per_host.sloth list_files_workflow

# Get system information
./sloth-runner run -f examples/agents/list_files_per_host.sloth system_info_workflow

# Complete audit (all tasks)
./sloth-runner run -f examples/agents/list_files_per_host.sloth complete_host_audit
```

## Current Status

✅ **Master server**: Running on 192.168.1.29:50053  
✅ **Agents connected**: ladyguica (192.168.1.16) and keiteguica (192.168.1.17)  
✅ **Workflow execution**: Working with automatic agent distribution  
⚠️ **Direct agent commands**: May have connectivity issues (under investigation)

## Usage Examples

### 1. Workflow-based execution (Recommended)
```bash
# Run distributed workflow
./sloth-runner run -f examples/agents/simple_ls.sloth distributed_ls_workflow
```

### 2. Direct agent commands (Advanced)
```bash
# List available agents
./sloth-runner agent list --master 192.168.1.29:50053

# Run command on specific agent (if connectivity allows)
./sloth-runner agent run ladyguica "hostname" --master 192.168.1.29:50053
```

### 3. Using the example script
```bash
# Run all examples
./examples/agents/agent_commands_example.sh
```

## Key Concepts

# Agent Examples

This directory contains examples of sloth-runner workflows that demonstrate distributed task execution across multiple agents.

## Current Implementation Status

✅ **Master server**: Running on 192.168.1.29:50053  
✅ **Agents connected**: ladyguica (192.168.1.16) and keiteguica (192.168.1.17)  
✅ **Workflow execution**: Working with automatic task distribution  
⚠️ **Agent targeting**: `:delegate_to()` not yet implemented in task API  
⚠️ **Direct agent commands**: Connectivity issues under investigation

## Available Examples

### 1. `simple_ls.sloth` - Basic Distributed Workflow
A working example that runs `ls` commands with automatic agent distribution:

**Usage:**
```bash
./sloth-runner run -f examples/agents/simple_ls.sloth distributed_ls_workflow
```

**Features:**
- 2 parallel tasks
- Hostname detection to show which machine executes each task
- Automatic task distribution among available agents

### 2. `list_files_per_host.sloth` - Advanced Multi-Task Workflow
A comprehensive example with multiple workflows:

#### Workflows:
- **list_files_workflow**: Lists files in `/home/chalkan3` 
- **system_info_workflow**: Collects system information (hostname, uptime, disk, memory)
- **complete_host_audit**: Runs all tasks for complete host audit

**Usage:**
```bash
# List files workflow
./sloth-runner run -f examples/agents/list_files_per_host.sloth list_files_workflow

# System information workflow  
./sloth-runner run -f examples/agents/list_files_per_host.sloth system_info_workflow

# Complete audit (all tasks)
./sloth-runner run -f examples/agents/list_files_per_host.sloth complete_host_audit
```

### 3. `agent_commands_example.sh` - Direct Agent Commands
A script demonstrating direct agent command execution:

**Usage:**
```bash
chmod +x examples/agents/agent_commands_example.sh
./examples/agents/agent_commands_example.sh
```

## Current Status

### ✅ What's Working
1. **Master-Agent Architecture**: Master server coordinates with connected agents
2. **Agent Registration**: Agents successfully register and maintain heartbeat
3. **Workflow Execution**: Tasks execute successfully in workflows
4. **Parallel Processing**: Multiple tasks run concurrently
5. **Agent Discovery**: Master can list and monitor agent status

### ⚠️ Current Limitations
1. **Task Targeting**: `:delegate_to("agent_name")` not implemented in current API
2. **Automatic Distribution**: Tasks currently execute locally instead of on remote agents
3. **Direct Agent Commands**: `agent run` has connectivity issues

## Agent Targeting (Future Feature)

**Planned Syntax:**
```lua
local my_task = task("my_task")
    :description("Task description")
    :command(function(this, params)
        -- Access agent info through params
        log.info("Executing on agent: " .. (params.agent_name or "local"))
        log.info("Hostname: " .. (os.getenv("HOSTNAME") or "unknown"))
        return true, "Task completed"
    end)
    :delegate_to("ladyguica")  -- Future: Target specific agent
    :timeout("30s")
    :build()
```

**Current Workaround:**
Use `agent run` command for direct execution on specific agents:
```bash
# Run command on specific agent
./sloth-runner agent run ladyguica "ls -la" --master 192.168.1.29:50053
./sloth-runner agent run keiteguica "hostname" --master 192.168.1.29:50053
```

### Parallel Execution
Control parallel execution with `max_parallel_tasks` in workflow config:
```lua
workflow.define("my_workflow")
    :tasks({ task1, task2, task3, task4 })
    :config({
        max_parallel_tasks = 2  -- Run up to 2 tasks simultaneously
    })
```

### Error Handling
Tasks should return `success, message, data`:
```lua
:command(function(this, params)
    local handle = io.popen("ls")
    local result = handle:read("*a")
    local success = handle:close()
    
    if success then
        return true, "Command succeeded", { output = result }
    else
        return false, "Command failed"
    end
end)
```

## Tips

1. **Always specify timeout** for tasks to prevent hanging
2. **Use descriptive agent names** for easier management
3. **Test connectivity** before running complex workflows
4. **Monitor logs** on both master and agents for debugging
5. **Use parallel execution** for independent tasks to improve performance