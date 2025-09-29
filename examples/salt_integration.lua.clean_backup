-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:32 -03


-- local example_task = task("task_name")
--     :description("Task description with modern DSL")
--     :command(function(params, deps)
--         -- Enhanced task logic
--         return true, "Task completed", { result = "success" }
--     end)
--     :timeout("30s")
--     :build()

-- workflow.define("workflow_name", {
--     description = "Workflow description - Modern DSL",
--     version = "2.0.0",
--     tasks = { example_task },
--     config = { timeout = "10m" }
-- })

-- Maintain backward compatibility with legacy format
TaskDefinitions = {
    salt_integration_group = {
        description = "Examples for integrating with SaltStack using the 'salt' module",
        tasks = {
            {
                name = "salt_ping_minion",
                description = "Pings a specific Salt minion using the fluent API",
                command = function(params, input)
                    log.info("Pinging Salt minion 'keiteguica'...")
                    local stdout, stderr, err = salt.target("keiteguica"):ping():result()

                    if err then
                        log.error("Salt ping failed: " .. err .. " Stderr: " .. stderr)
                        return false, "Salt ping failed"
                    else
                        log.info("Salt ping successful. Result: " .. tostring(stdout))
                        return true, "Salt minion pinged successfully", {result = stdout}
                    end
                end,
            },
            {
                name = "salt_run_command_on_all_minions",
                description = "Runs a shell command on all Salt minions using the fluent API",
                command = function(params, input)
                    log.info("Running 'ls -l /tmp' on all Salt minions...")
                    local stdout, stderr, err = salt.target("*"):cmd("cmd.run", "ls -l /tmp"):result()

                    if err then
                        log.error("Salt cmd.run failed: " .. err .. " Stderr: " .. stderr)
                        return false, "Salt cmd.run failed"
                    else
                        log.info("Salt cmd.run successful. Result: " .. tostring(stdout))
                        return true, "Salt cmd.run executed successfully", {result = stdout}
                    end
                end,
            },
        },
    },
}
