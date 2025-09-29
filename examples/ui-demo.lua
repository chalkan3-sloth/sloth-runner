-- Example Lua script for Sloth Runner UI Demo
-- This script demonstrates various features that can be managed via the web interface

-- Load required modules
local shell = require("shell")
local json = require("json")
local system = require("system")

-- Task configuration
task_group("ui-demo-tasks") {
    description = "Demonstration tasks for the Sloth Runner UI",
    
    -- Simple system info task
    task("system-info") {
        description = "Gather system information",
        script = function()
            print("🖥️  System Information:")
            print("OS: " .. system.os())
            print("Arch: " .. system.arch())
            print("User: " .. system.user())
            print("Hostname: " .. system.hostname())
            
            local uptime = shell.exec("uptime")
            print("Uptime: " .. uptime.stdout)
            
            return {
                status = "success",
                data = {
                    os = system.os(),
                    arch = system.arch(),
                    user = system.user(),
                    hostname = system.hostname()
                }
            }
        end
    },
    
    -- File operations task
    task("file-operations") {
        description = "Demonstrate file operations",
        script = function()
            print("📁 File Operations Demo:")
            
            -- Create a temporary directory
            local temp_dir = "/tmp/sloth-runner-demo"
            shell.exec("mkdir -p " .. temp_dir)
            print("Created directory: " .. temp_dir)
            
            -- Create some demo files
            local files = {"task1.txt", "task2.txt", "task3.txt"}
            for i, filename in ipairs(files) do
                local filepath = temp_dir .. "/" .. filename
                local content = "Demo file " .. i .. " created at " .. os.date()
                shell.exec("echo '" .. content .. "' > " .. filepath)
                print("Created file: " .. filepath)
            end
            
            -- List files
            local ls_result = shell.exec("ls -la " .. temp_dir)
            print("Directory contents:")
            print(ls_result.stdout)
            
            -- Cleanup
            shell.exec("rm -rf " .. temp_dir)
            print("Cleaned up temporary directory")
            
            return {
                status = "success",
                message = "File operations completed successfully"
            }
        end
    },
    
    -- Network task
    task("network-check") {
        description = "Check network connectivity",
        script = function()
            print("🌐 Network Connectivity Check:")
            
            local hosts = {"google.com", "github.com", "stackoverflow.com"}
            local results = {}
            
            for _, host in ipairs(hosts) do
                print("Testing connectivity to " .. host .. "...")
                local result = shell.exec("ping -c 1 " .. host)
                
                if result.code == 0 then
                    print("✅ " .. host .. " - OK")
                    results[host] = "OK"
                else
                    print("❌ " .. host .. " - FAILED")
                    results[host] = "FAILED"
                end
            end
            
            return {
                status = "success",
                data = results
            }
        end
    },
    
    -- Process monitoring task
    task("process-monitor") {
        description = "Monitor system processes",
        script = function()
            print("📊 Process Monitoring:")
            
            -- Get top processes
            local top_result = shell.exec("ps aux | head -10")
            print("Top processes:")
            print(top_result.stdout)
            
            -- Memory usage
            local mem_result = shell.exec("free -h")
            if mem_result.code == 0 then
                print("Memory usage:")
                print(mem_result.stdout)
            else
                -- macOS alternative
                local vm_result = shell.exec("vm_stat")
                print("Memory statistics:")
                print(vm_result.stdout)
            end
            
            -- Disk usage
            local disk_result = shell.exec("df -h")
            print("Disk usage:")
            print(disk_result.stdout)
            
            return {
                status = "success",
                message = "Process monitoring completed"
            }
        end
    },
    
    -- Long running task simulation
    task("long-running-demo") {
        description = "Simulate a long running task with progress updates",
        script = function()
            print("⏳ Starting long running task simulation...")
            
            local steps = {"Initializing", "Processing data", "Performing calculations", "Generating report", "Finalizing"}
            
            for i, step in ipairs(steps) do
                print(string.format("🔄 Step %d/%d: %s", i, #steps, step))
                
                -- Simulate work with sleep
                shell.exec("sleep 2")
                
                -- Update progress
                local progress = math.floor((i / #steps) * 100)
                print(string.format("Progress: %d%%", progress))
            end
            
            print("✅ Long running task completed successfully!")
            
            return {
                status = "success",
                message = "Long running task completed",
                duration = "10 seconds",
                steps_completed = #steps
            }
        end
    },
    
    -- Error simulation task
    task("error-simulation") {
        description = "Simulate task failure for testing",
        script = function()
            print("⚠️  Simulating task failure...")
            
            -- Simulate some work
            print("Performing initial checks...")
            shell.exec("sleep 1")
            
            print("Processing data...")
            shell.exec("sleep 1")
            
            -- Simulate an error
            print("❌ Error occurred: Simulated failure for testing purposes")
            
            error("Simulated task failure - this is intentional for demo purposes")
        end
    },
    
    -- JSON processing task
    task("json-processing") {
        description = "Process JSON data",
        script = function()
            print("📄 JSON Processing Demo:")
            
            -- Create sample JSON data
            local sample_data = {
                name = "Sloth Runner Demo",
                version = "1.0.0",
                features = {"Web UI", "Lua Scripting", "Agent Management"},
                metadata = {
                    created = os.date(),
                    author = "Sloth Runner",
                    type = "demo"
                }
            }
            
            -- Convert to JSON string
            local json_string = json.encode(sample_data)
            print("Generated JSON:")
            print(json_string)
            
            -- Parse back from JSON
            local parsed_data = json.decode(json_string)
            print("Parsed data:")
            print("Name: " .. parsed_data.name)
            print("Version: " .. parsed_data.version)
            print("Features: " .. table.concat(parsed_data.features, ", "))
            
            return {
                status = "success",
                data = parsed_data
            }
        end
    }
}

-- Pipeline example
pipeline("demo-pipeline") {
    description = "Demonstration pipeline for UI",
    
    stage("initialization") {
        task("system-info")
    },
    
    stage("operations") {
        parallel {
            task("file-operations"),
            task("network-check")
        }
    },
    
    stage("monitoring") {
        task("process-monitor")
    },
    
    stage("data-processing") {
        task("json-processing")
    },
    
    stage("completion") {
        task("long-running-demo")
    }
}