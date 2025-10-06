-- Hook: alert-agent-down.lua
-- Event: agent.disconnected
-- Description: Sends alert when an agent disconnects

function on_event()
    local agent = event.agent

    log.error("🔴 Agent disconnected!")
    log.error("  Name: " .. agent.name)
    log.error("  Address: " .. agent.address)
    log.error("  Timestamp: " .. os.date("%Y-%m-%d %H:%M:%S", event.timestamp))

    -- Send alert notification
    log.info("📨 Sending alert notification...")
    http.post("https://hooks.slack.com/YOUR_WEBHOOK_URL")

    return true
end
