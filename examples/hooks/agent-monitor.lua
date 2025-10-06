-- Hook: Agent Monitor
-- Monitors agent registration and disconnection events

function on_event()
    local agent = event.data.agent or event.agent

    -- Create log entry
    local log_entry = string.format(
        "[%s] 🟢 Agent REGISTERED: %s (%s) - Version: %s\n",
        os.date("%Y-%m-%d %H:%M:%S"),
        agent.name or "unknown",
        agent.address or "unknown",
        agent.version or "unknown"
    )

    -- Append to agent events log
    file_ops.lineinfile({
        path = "/tmp/agent-events.log",
        line = log_entry,
        create = true
    })

    print("✅ Agent registered: " .. (agent.name or "unknown") .. " at " .. (agent.address or "unknown"))

    return true
end
