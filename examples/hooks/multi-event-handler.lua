-- Hook: Multi-Event Handler
-- Handles task started events

function on_event()
    local task = event.task or event.data.task

    local log_entry = string.format(
        "[%s] 🏁 Task STARTED: %s on %s\n",
        os.date("%Y-%m-%d %H:%M:%S"),
        task.task_name or "unknown",
        task.agent_name or "unknown"
    )

    file_ops.lineinfile({
        path = "/tmp/task-starts.log",
        line = log_entry,
        create = true
    })

    print("🏁 Task started: " .. (task.task_name or "unknown"))

    return true
end
