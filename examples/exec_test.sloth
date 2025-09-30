-- MODERN DSL ONLY - Exec Module Testing
-- Comprehensive testing of the exec module capabilities

-- Task 1: Print template variables and environment info
local template_vars_task = task("print_template_vars")
    :description("Prints template variables and environment information")
    :command(function(params, input)
        log.info("🌍 Environment Information:")
        
        -- Get template variables (these would be populated by the runner)
        local env = "{{.Env}}"
        local is_prod = {{.IsProduction}}
        local shards = {}
        {{- range .Shards }}
        table.insert(shards, {{.}})
        {{- end }}

        -- Display environment info
        log.info("  Environment: " .. (env or "development"))
        log.warn("  Is Production: " .. tostring(is_prod or false))
        log.debug("  Shards: " .. table.concat(shards, ", "))
        log.error("  ⚠️  This is a test error message from Lua.")

        -- Also get system environment variables
        local user = os.getenv("USER") or os.getenv("USERNAME") or "unknown"
        local home = os.getenv("HOME") or os.getenv("USERPROFILE") or "unknown"
        local path = os.getenv("PATH") or "unknown"
        
        log.info("💻 System Environment:")
        log.info("  User: " .. user)
        log.info("  Home: " .. home)
        log.info("  PATH length: " .. string.len(path) .. " characters")

        return true, "Template variables printed", {
            environment = env or "development",
            production = is_prod or false,
            shard_count = #shards,
            system_user = user,
            home_dir = home
        }
    end)
    :timeout("30s")
    :build()

-- Task 2: Basic echo command test
local echo_task = task("run_echo_command")
    :description("Test basic echo command execution")
    :depends_on({"print_template_vars"})
    :command(function(params, deps)
        log.info("📢 Testing basic echo command...")
        
        local message = "Hello from Modern DSL exec test!"
        local success, output = exec.run("echo '" .. message .. "'")
        
        if success then
            log.info("✅ Echo command successful")
            log.info("📤 Output: " .. output)
            return true, "Echo test passed", {
                command = "echo",
                output = output,
                message = message
            }
        else
            log.error("❌ Echo command failed: " .. output)
            return false, "Echo test failed"
        end
    end)
    :build()

-- Task 3: Command with timeout test
local timeout_task = task("test_command_timeout")
    :description("Test command execution with timeout handling")
    :depends_on({"run_echo_command"})
    :command(function(params, deps)
        log.info("⏱️  Testing command timeout handling...")
        
        -- Test quick command (should succeed)
        local success1, output1 = exec.run("echo 'Quick command'", {timeout = 5})
        if not success1 then
            log.error("❌ Quick command failed: " .. output1)
            return false, "Quick command test failed"
        end
        
        log.info("✅ Quick command test passed")
        
        -- Test command with artificial delay (but within timeout)
        local success2, output2 = exec.run("sleep 1 && echo 'Delayed command'", {timeout = 3})
        if not success2 then
            log.error("❌ Delayed command failed: " .. output2)
            return false, "Delayed command test failed"
        end
        
        log.info("✅ Delayed command test passed")
        
        return true, "Timeout tests completed", {
            quick_command_output = output1,
            delayed_command_output = output2,
            tests_passed = 2
        }
    end)
    :timeout("10s")
    :build()

-- Task 4: Environment variable test
local env_test_task = task("test_environment_variables")
    :description("Test environment variable handling in commands")
    :depends_on({"test_command_timeout"})
    :command(function(params, deps)
        log.info("🔧 Testing environment variable handling...")
        
        -- Set custom environment variables for command
        local env_vars = {
            CUSTOM_VAR = "sloth_runner_test",
            TEST_NUMBER = "42",
            TEST_BOOL = "true"
        }
        
        local success, output = exec.run("echo \"Custom: $CUSTOM_VAR, Number: $TEST_NUMBER, Bool: $TEST_BOOL\"", {
            env = env_vars,
            timeout = 5
        })
        
        if success then
            log.info("✅ Environment variable test passed")
            log.info("📤 Output: " .. output)
            
            -- Verify output contains expected values
            local contains_custom = string.find(output, "sloth_runner_test") ~= nil
            local contains_number = string.find(output, "42") ~= nil
            local contains_bool = string.find(output, "true") ~= nil
            
            return true, "Environment test completed", {
                output = output,
                contains_custom_var = contains_custom,
                contains_number = contains_number,
                contains_bool = contains_bool,
                all_vars_present = contains_custom and contains_number and contains_bool
            }
        else
            log.error("❌ Environment variable test failed: " .. output)
            return false, "Environment test failed"
        end
    end)
    :build()

-- Task 5: Working directory test
local workdir_task = task("test_working_directory")
    :description("Test working directory handling")
    :depends_on({"test_environment_variables"})
    :command(function(params, deps)
        log.info("📁 Testing working directory handling...")
        
        -- Get current directory first
        local success1, current_dir = exec.run("pwd")
        if not success1 then
            log.error("❌ Failed to get current directory: " .. current_dir)
            return false, "Working directory test failed"
        end
        
        log.info("📍 Current directory: " .. current_dir)
        
        -- Test command with specific working directory
        local success2, output2 = exec.run("pwd && ls -la", {
            workdir = "/tmp",
            timeout = 10
        })
        
        if success2 then
            log.info("✅ Working directory test passed")
            log.info("📁 Temp directory listing completed")
            
            -- Check if output contains /tmp
            local in_tmp_dir = string.find(output2, "/tmp") ~= nil
            
            return true, "Working directory test completed", {
                original_dir = current_dir,
                temp_dir_output = output2,
                executed_in_tmp = in_tmp_dir
            }
        else
            log.error("❌ Working directory test failed: " .. output2)
            return false, "Working directory test failed"
        end
    end)
    :build()

-- Task 6: Error handling test
local error_handling_task = task("test_error_handling")
    :description("Test error handling for failed commands")
    :depends_on({"test_working_directory"})
    :command(function(params, deps)
        log.info("🚨 Testing error handling...")
        
        -- Intentionally run a command that should fail
        local success, output = exec.run("ls /nonexistent/directory/that/does/not/exist", {
            timeout = 5
        })
        
        if not success then
            log.info("✅ Error handling test passed (command failed as expected)")
            log.info("📤 Error output: " .. output)
            
            -- This is expected behavior - command should fail
            return true, "Error handling test completed", {
                command_failed_as_expected = true,
                error_output = output,
                error_captured = string.len(output) > 0
            }
        else
            log.warn("⚠️  Command unexpectedly succeeded: " .. output)
            return true, "Error handling test completed", {
                command_failed_as_expected = false,
                unexpected_success = true,
                output = output
            }
        end
    end)
    :build()

-- Exec Test Workflow
workflow.define("exec_module_test", {
    description = "Comprehensive exec module testing - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"exec", "testing", "commands", "modern-dsl"},
        complexity = "beginner",
        estimated_duration = "2m"
    },
    
    tasks = {
        template_vars_task,
        echo_task,
        timeout_task,
        env_test_task,
        workdir_task,
        error_handling_task
    },
    
    config = {
        timeout = "10m",
        retry_policy = "linear",
        max_parallel_tasks = 1, -- Sequential execution for proper testing
        fail_fast = false -- Continue even if some tests fail
    },
    
    on_start = function()
        log.info("🚀 Starting exec module comprehensive test...")
        log.info("🧪 Testing: echo, timeout, env vars, workdir, error handling")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Exec module test completed successfully!")
            log.info("✅ All exec module features tested")
            
            -- Summary of test results
            local test_summary = {
                template_vars = results.print_template_vars ~= nil,
                echo_command = results.run_echo_command ~= nil,
                timeout_handling = results.test_command_timeout ~= nil,
                env_variables = results.test_environment_variables ~= nil,
                working_directory = results.test_working_directory ~= nil,
                error_handling = results.test_error_handling ~= nil
            }
            
            local passed_tests = 0
            for test, passed in pairs(test_summary) do
                if passed then passed_tests = passed_tests + 1 end
            end
            
            log.info("📊 Test Summary: " .. passed_tests .. "/6 tests completed")
        else
            log.error("❌ Exec module test failed!")
            log.warn("🔍 Check individual test results for details")
        end
        return true
    end
})
