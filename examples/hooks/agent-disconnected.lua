-- Hook: Agent Disconnected Monitor
-- Logs when agents disconnect

function on_event()
    local agent = event.data.agent or event.agent

    -- Create log entry
    local log_entry = string.format(
        "[%s] 🔴 Agent DISCONNECTED: %s (%s)\n",
        os.date("%Y-%m-%d %H:%M:%S"),
        agent.name or "unknown",
        agent.address or "unknown"
    )

    -- Append to agent events log
    file_ops.lineinfile({
        path = "/tmp/agent-events.log",
        line = log_entry,
        create = true
    })

    print("⚠️  Agent disconnected: " .. (agent.name or "unknown"))

    return true
end
