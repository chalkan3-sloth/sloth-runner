-- Hook: notify-agent-join.lua
-- Event: agent.registered
-- Description: Sends notification when a new agent joins the cluster

-- This function is called when the event is triggered
function on_event()
    -- Access event data
    local agent = event.agent

    log.info("🤖 New agent registered!")
    log.info("  Name: " .. agent.name)
    log.info("  Address: " .. agent.address)
    log.info("  Version: " .. (agent.version or "unknown"))

    -- Check if agent has tags
    if agent.tags then
        log.info("  Tags: " .. table.concat(agent.tags, ", "))

        -- Only notify for production agents
        if contains(agent.tags, "production") then
            log.info("🔔 This is a production agent - sending notification")

            -- Send notification (you can customize this)
            http.post("https://hooks.slack.com/YOUR_WEBHOOK_URL")
        end
    end

    -- Check system info if available
    if agent.system_info then
        log.info("  OS: " .. (agent.system_info.os or "unknown"))
        log.info("  CPUs: " .. (agent.system_info.cpus or "unknown"))
        log.info("  Memory: " .. (agent.system_info.memory or "unknown") .. " MB")
    end

    log.info("✅ Agent setup completed successfully!")

    -- Return true to indicate success
    return true
end
