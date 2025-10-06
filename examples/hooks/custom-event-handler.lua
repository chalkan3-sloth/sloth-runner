-- Hook: Custom Event Handler
-- Demonstrates handling custom events dispatched from workflows

function on_event()
    -- Custom events have flexible data structure
    local event_data = tostring(event.data)

    local log_entry = string.format(
        "[%s] 🎯 Custom Event Received\n   Event Type: %s\n   Data: %s\n\n",
        os.date("%Y-%m-%d %H:%M:%S"),
        event.type or "custom",
        event_data
    )

    -- Log custom event
    file_ops.lineinfile({
        path = "/tmp/custom-events.log",
        line = log_entry,
        create = true
    })

    print("🎯 Custom event processed: " .. event.type)

    return true
end
