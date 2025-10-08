-- Hook: Port Scanner
-- Triggers when ports are opened or closed
-- Event types: port.opened, port.closed

return {
    name = "port_scanner",
    description = "Monitor port changes for security",
    event_types = {"port.opened", "port.closed"},
    enabled = true,

    execute = function(event)
        local port = event.data.port or 0
        local event_type = event.type

        if event_type == "port.opened" then
            log.info("🔓 PORT OPENED: " .. tostring(port))

            -- Security check for unexpected ports
            local expected_ports = {80, 443, 22, 3000, 8080, 9090}
            local is_expected = false
            for _, p in ipairs(expected_ports) do
                if p == port then
                    is_expected = true
                    break
                end
            end

            if not is_expected then
                log.warn("⚠️  UNEXPECTED PORT OPENED: " .. tostring(port))
                -- Could trigger security scan, alert team
            end

        elseif event_type == "port.closed" then
            log.info("🔒 PORT CLOSED: " .. tostring(port))
            -- Could check if service should be running
        end

        return true
    end
}
