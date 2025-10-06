-- Hook: Workflow Tracker
-- Tracks workflow execution statistics

function on_event()
    -- This would normally come from event data
    -- For now we'll create a mock entry
    local log_entry = string.format(
        "[%s] ✅ Workflow completed - Event ID: %s\n",
        os.date("%Y-%m-%d %H:%M:%S"),
        event.ID or "unknown"
    )

    -- Append to workflow log
    file_ops.lineinfile({
        path = "/tmp/workflow-completions.log",
        line = log_entry,
        create = true
    })

    print("✅ Workflow completion tracked")

    return true
end
