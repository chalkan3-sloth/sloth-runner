-- Hook: Task Completion Logger
-- Logs all task completions to a file

function on_event()
    local task = event.task or event.data.task

    -- Create log entry
    local log_entry = string.format(
        "[%s] Task '%s' completed on agent '%s' with exit code %d (duration: %s)\n",
        os.date("%Y-%m-%d %H:%M:%S"),
        task.task_name or "unknown",
        task.agent_name or "unknown",
        task.exit_code or 0,
        task.duration or "unknown"
    )

    -- Append to log file
    file_ops.lineinfile({
        path = "/tmp/task-completions.log",
        line = log_entry,
        create = true
    })

    print("✅ Logged task completion: " .. (task.task_name or "unknown"))

    return true
end
