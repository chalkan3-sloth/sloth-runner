-- Hook: Task Failure Alert
-- Writes task failures to a file for monitoring

function on_event()
    local task = event.task or event.data.task

    -- Create alert entry with full details
    local alert_entry = string.format([[
========================================
TASK FAILURE ALERT
========================================
Time: %s
Task: %s
Agent: %s
Status: %s
Exit Code: %d
Error: %s
Duration: %s
========================================

]],
        os.date("%Y-%m-%d %H:%M:%S"),
        task.task_name or "unknown",
        task.agent_name or "unknown",
        task.status or "failed",
        task.exit_code or -1,
        task.error or "unknown error",
        task.duration or "unknown"
    )

    -- Write to failures log
    file_ops.lineinfile({
        path = "/tmp/task-failures.log",
        line = alert_entry,
        create = true
    })

    print("🚨 ALERT: Task " .. (task.task_name or "unknown") .. " failed on " .. (task.agent_name or "unknown"))
    print("   Error: " .. (task.error or "unknown"))

    return true
end
