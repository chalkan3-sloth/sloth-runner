-- MODERN DSL ONLY: Exec Module Testing
-- Legacy TaskDefinitions removed - Modern DSL syntax only

-- Task 1: Print template variables
local template_vars_task = task("print_template_vars")
    :description("Prints template variables with modern DSL")
    :command(function(params, input)
        local env = "{{.Env}}"
        local is_prod = {{.IsProduction}}
        local shards = {}
        {{- range .Shards }}
        table.insert(shards, {{.}})
        {{- end }}

        log.info("🌍 Environment: " .. env)
        log.warn("🏭 Is Production: " .. tostring(is_prod))
        log.debug("🔢 Shards: " .. table.concat(shards, ", "))
        log.error("⚠️  This is a test error message from Lua.")

        return true, "Template variables printed", {
            environment = env,
            production = is_prod,
            shard_count = #shards
        }
    end)
    :timeout("30s")
    :build()

-- Task 2: Echo command test
local echo_task = task("run_echo_command")
    :description("Runs echo command using exec.run with modern DSL")
    :depends_on({"print_template_vars"})
    :command(function(params, deps)
        log.info("🔊 Running echo command...")
        local result = exec.run("echo 'Hello from modern exec!'")
        
        if not result.success then
            return false, "Command failed: " .. result.stderr
        else
            log.info("✅ Echo command output: " .. result.stdout)
            return true, "Command executed successfully", {
                stdout = result.stdout, 
                stderr = result.stderr,
                execution_time = os.time()
            }
        end
    end)
    :build()

-- Task 3: List files test
local list_files_task = task("list_files")
    :description("Lists files using exec.run with modern DSL")
    :depends_on({"run_echo_command"})
    :command(function(params, deps)
        log.info("📁 Listing files...")
        local result = exec.run("ls -la")
        
        if not result.success then
            return false, "ls command failed: " .. result.stderr
        else
            log.info("✅ Files listed successfully")
            return true, "ls command executed successfully", {
                stdout = result.stdout, 
                stderr = result.stderr,
                file_count = string.match(result.stdout, "(%d+)") or "unknown"
            }
        end
    end)
    :build()

-- Task 4: Additional exec test
local additional_exec_task = task("additional_exec_test")
    :description("Additional exec module testing")
    :command(function(params, deps)
        log.info("🔧 Running additional exec test...")
        
        -- Test multiple commands
        local date_result = exec.run("date")
        local whoami_result = exec.run("whoami")
        
        if date_result.success and whoami_result.success then
            log.info("📅 Current date: " .. date_result.stdout)
            log.info("👤 Current user: " .. whoami_result.stdout)
            
            return true, "Additional tests completed", {
                date = date_result.stdout,
                user = whoami_result.stdout,
                test_count = 2
            }
        else
            return false, "Some commands failed"
        end
    end)
    :build()

-- Define exec testing workflow
workflow.define("exec_module_test", {
    description = "Comprehensive exec module testing - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        category = "testing",
        tags = {"exec", "commands", "modern-dsl"},
        author = "Sloth Runner Team"
    },
    
    tasks = {
        template_vars_task,
        echo_task,
        list_files_task,
        additional_exec_task
    },
    
    config = {
        max_parallel_tasks = 2,
        timeout = "10m",
        retry_policy = "linear"
    },
    
    on_start = function()
        log.info("🚀 Starting exec module test suite...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ All exec tests passed!")
            log.info("📊 Test results summary:")
            for task_name, result in pairs(results) do
                log.info("  " .. task_name .. ": " .. (result.test_count or "1") .. " tests")
            end
        else
            log.error("❌ Some exec tests failed!")
        end
        return true
    end
})
