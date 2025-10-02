# Agent Examples - SQLite Integration & Name Resolution

This directory contains examples of sloth-runner workflows that demonstrate distributed task execution across multiple agents with **SQLite database integration** and **agent name resolution**.

## ✅ **IMPLEMENTED FEATURES**

### 🗄️ **SQLite Database for Agents**
- ✅ **Automatic database** creation in `.sloth-cache/agents.db`
- ✅ **Agent registration** when connecting to master
- ✅ **Heartbeat tracking** for active/inactive status
- ✅ **Automatic cleanup** of inactive agents
- ✅ **Name-to-IP resolution** for agent targeting

### 🎯 **Modern Delegate_to API**
- ✅ **Agent name resolution**: `:delegate_to("agent_name")`
- ✅ **Direct address support**: `:delegate_to("192.168.1.16:50051")`
- ✅ **Automatic fallback** for unresolved names
- ✅ **Database-backed lookup** with caching

## Current Implementation Status

✅ **Master server**: Running with SQLite backend  
✅ **Agents connected**: ladyguica and keiteguica registered  
✅ **Database active**: `.sloth-cache/agents.db` operational  
✅ **Name resolution**: Agent names resolve to addresses  
✅ **Modern API**: `:delegate_to()` fully implemented  

## Available Examples

### 1. `simple_ls.sloth` - Agent Name Resolution
Demonstrates using agent names instead of IP addresses:

```lua
local task_ladyguica = task("ls_task_ladyguica")
    :description("List files on ladyguica agent")
    :command(function(this, params)
        log.info("Executing on: " .. (params.agent_name or "local"))
        -- Task logic here
        return true, "Task completed"
    end)
    :delegate_to("ladyguica")  -- Agent name (not IP!)
    :timeout("30s")
    :build()
```

**Usage:**
```bash
./sloth-runner run -f examples/agents/simple_ls.sloth distributed_ls_workflow
```

### 2. Database Operations

**View registered agents:**
```bash
./sloth-runner agent list --master 192.168.1.29:50053
```

**View SQLite database directly:**
```bash
sqlite3 .sloth-cache/agents.db "SELECT name, address, status, datetime(last_heartbeat, 'unixepoch') as last_seen FROM agents;"
```

**Database schema:**
```sql
CREATE TABLE agents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    address TEXT NOT NULL,
    status TEXT DEFAULT 'Active',
    last_heartbeat INTEGER DEFAULT 0,
    registered_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
```

### 3. Agent Management

**Start master with SQLite:**
```bash
./sloth-runner master --port 50053 --daemon
```

**Connect agents:**
```bash
# On remote hosts
sloth-runner agent start --master 192.168.1.29:50053 --name ladyguica --port 50051
sloth-runner agent start --master 192.168.1.29:50053 --name keiteguica --port 50051
```

## Key Concepts

### Agent Name Resolution
The system automatically resolves agent names to IP addresses:

```lua
-- These are equivalent if ladyguica is registered at 192.168.1.16:50051
:delegate_to("ladyguica")           -- Resolved via SQLite database
:delegate_to("192.168.1.16:50051")  -- Direct address
```

### Database Integration
- **Automatic registration**: Agents auto-register when connecting
- **Heartbeat tracking**: Status updated every 30 seconds  
- **Cleanup**: Inactive agents removed after 24 hours
- **Persistent storage**: Survives master restarts

### Error Handling
```lua
-- System handles various scenarios:
:delegate_to("unknown_agent")    -- Error: agent not found
:delegate_to("inactive_agent")   -- Error: agent not active  
:delegate_to("192.168.1.99:50051") -- Error: connection refused
```

## Implementation Details

### Agent Resolver Interface
```go
type AgentResolver interface {
    GetAgentAddress(agentName string) (string, error)
}
```

### Resolution Logic
1. **Check format**: If contains `:`, treat as direct address
2. **Database lookup**: Query SQLite for agent name
3. **Activity check**: Ensure agent is active (heartbeat < 60s)
4. **Return address**: Provide `IP:port` for connection

## Troubleshooting

### Common Issues

**Agent not found:**
```bash
# Check registered agents
./sloth-runner agent list --master 192.168.1.29:50053

# Check database directly  
sqlite3 .sloth-cache/agents.db "SELECT * FROM agents WHERE name='agent_name';"
```

**Database locked:**
```bash
# Restart master if database issues
pkill -f "sloth-runner master"
./sloth-runner master --port 50053 --daemon
```

**Agent inactive:**
```bash
# Check agent status
./sloth-runner agent list --master 192.168.1.29:50053

# Restart agent if needed
sloth-runner agent start --master 192.168.1.29:50053 --name agent_name
```

## Advanced Usage

### Programmatic Access
```lua
-- Access agent resolver in tasks
local resolver = require("agent_resolver")
local address = resolver.get_agent_address("ladyguica")
```

### Database Cleanup
```lua
-- Manual cleanup of inactive agents
local db = require("agent_db")
local removed = db.cleanup_inactive_agents(24) -- Remove agents inactive for 24+ hours
```

## Performance Notes

- **Lookup caching**: Agent addresses cached for 60 seconds
- **Connection pooling**: gRPC connections reused when possible
- **Heartbeat batching**: Status updates batched every 30 seconds
- **Database indexing**: Indexed on name, status, and heartbeat

---

**The SQLite integration and agent name resolution system is fully implemented and operational!** 🚀