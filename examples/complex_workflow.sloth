-- Complex Workflow Example - Modern DSL
-- Demonstrates advanced Modern DSL features and patterns

-- Task 1: Environment setup with validation
local setup_environment_task = task("setup_environment")
    :description("Set up and validate the execution environment")
    :command(function(params, deps)
        log.info("🔧 Setting up environment...")
        
        -- Validate required tools
        local required_tools = {"git", "docker", "curl"}
        local missing_tools = {}
        
        for _, tool in ipairs(required_tools) do
            local result = exec.run("which " .. tool)
            if not result.success then
                table.insert(missing_tools, tool)
            end
        end
        
        if #missing_tools > 0 then
            return false, "Missing required tools: " .. table.concat(missing_tools, ", ")
        end
        
        -- Create workspace directory
        local workspace = "./complex-workflow-workspace"
        fs.mkdir(workspace)
        
        -- Set environment variables
        os.setenv("WORKSPACE_DIR", workspace)
        os.setenv("BUILD_ID", "build-" .. os.time())
        
        return true, "Environment setup completed", {
            workspace_dir = workspace,
            build_id = os.getenv("BUILD_ID"),
            validated_tools = required_tools,
            timestamp = os.date()
        }
    end)
    :timeout("2m")
    :retries(2, "linear")
    :on_success(function(params, output)
        log.info("✅ Environment setup completed")
        log.info("📁 Workspace: " .. output.workspace_dir)
        log.info("🔖 Build ID: " .. output.build_id)
    end)
    :build()

-- Task 2: Parallel data processing
local process_data_task = task("process_data")
    :description("Process data in parallel with Modern DSL async patterns")
    :depends_on({"setup_environment"})
    :command(function(params, deps)
        log.info("⚡ Processing data in parallel...")
        
        local workspace = deps.setup_environment.workspace_dir
        
        -- Modern DSL async parallel processing
        local results = async.parallel({
            process_logs = function()
                log.info("Processing logs...")
                -- Simulate log processing
                exec.run("sleep 2")
                return {
                    processed_files = 150,
                    size_mb = 45.2,
                    status = "completed"
                }
            end,
            
            process_metrics = function()
                log.info("Processing metrics...")
                -- Simulate metrics processing
                exec.run("sleep 1.5")
                return {
                    data_points = 2500,
                    anomalies_detected = 3,
                    status = "completed"
                }
            end,
            
            process_events = function()
                log.info("Processing events...")
                -- Simulate event processing
                exec.run("sleep 1")
                return {
                    events_processed = 890,
                    errors_found = 0,
                    status = "completed"
                }
            end
        }, {
            max_workers = 3,
            timeout = "5m",
            fail_fast = false
        })
        
        -- Validate all processing completed
        local all_success = true
        for name, result in pairs(results) do
            if not result or result.status ~= "completed" then
                all_success = false
                log.error("Failed to process: " .. name)
            end
        end
        
        if all_success then
            return true, "Parallel data processing completed", {
                processing_results = results,
                total_files = results.process_logs.processed_files,
                total_metrics = results.process_metrics.data_points,
                total_events = results.process_events.events_processed
            }
        else
            return false, "Some data processing tasks failed"
        end
    end)
    :timeout("10m")
    :async(true)
    :build()

-- Task 3: API integration with circuit breaker
local api_integration_task = task("api_integration")
    :description("Integrate with external APIs using circuit breaker pattern")
    :depends_on({"process_data"})
    :command(function(params, deps)
        log.info("🔌 Integrating with external APIs...")
        
        local api_results = {}
        
        -- Use circuit breaker for external API calls
        local weather_result = circuit.protect("weather_api", function()
            log.info("Calling weather API...")
            return net.http_get("https://api.openweathermap.org/data/2.5/weather?q=London&appid=demo", {
                timeout = "10s",
                retries = 2
            })
        end)
        
        if weather_result.success then
            api_results.weather = {
                status = "success",
                data = weather_result.data
            }
        else
            api_results.weather = {
                status = "failed",
                error = weather_result.error
            }
        end
        
        -- Another API call with different circuit breaker
        local status_result = circuit.protect("status_api", function()
            log.info("Calling status API...")
            return net.http_get("https://httpstat.us/200", {
                timeout = "5s"
            })
        end)
        
        api_results.status_check = status_result.success and "healthy" or "unhealthy"
        
        return true, "API integration completed", {
            api_results = api_results,
            circuit_breaker_stats = circuit.get_stats(),
            integration_timestamp = os.time()
        }
    end)
    :retries(3, "exponential")
    :on_failure(function(params, error)
        log.warn("API integration failed, continuing with degraded functionality")
    end)
    :build()

-- Task 4: Conditional deployment
local conditional_deploy_task = task("conditional_deployment")
    :description("Deploy only if all conditions are met")
    :depends_on({"api_integration"})
    :run_if(function(params, deps)
        -- Complex conditional logic
        local data_results = deps.process_data
        local api_results = deps.api_integration
        
        -- Deploy only if:
        -- 1. Data processing was successful
        -- 2. No critical errors found
        -- 3. At least one API is healthy
        local should_deploy = data_results and 
                             data_results.total_events > 0 and
                             api_results.api_results.status_check == "healthy"
        
        if should_deploy then
            log.info("✅ All deployment conditions met")
        else
            log.warn("⚠️ Deployment conditions not met, skipping deployment")
        end
        
        return should_deploy
    end)
    :command(function(params, deps)
        log.info("🚀 Executing conditional deployment...")
        
        local workspace = deps.setup_environment.workspace_dir
        local build_id = deps.setup_environment.build_id
        
        -- Simulate deployment process
        local deployment_steps = {
            "Preparing deployment package",
            "Validating deployment configuration", 
            "Deploying to staging environment",
            "Running smoke tests",
            "Promoting to production"
        }
        
        for i, step in ipairs(deployment_steps) do
            log.info(string.format("[%d/%d] %s", i, #deployment_steps, step))
            exec.run("sleep 0.5") -- Simulate work
        end
        
        return true, "Deployment completed successfully", {
            deployment_id = "deploy-" .. build_id,
            deployed_at = os.date(),
            deployment_steps = deployment_steps,
            environment = "production"
        }
    end)
    :timeout("15m")
    :artifacts({"deployment-logs"})
    :on_success(function(params, output)
        log.info("🎉 Deployment successful!")
        log.info("🆔 Deployment ID: " .. output.deployment_id)
    end)
    :build()

-- Task 5: Cleanup and reporting
local cleanup_task = task("cleanup_and_report")
    :description("Clean up resources and generate final report")
    :depends_on({"conditional_deployment"})
    :command(function(params, deps)
        log.info("🧹 Cleaning up and generating report...")
        
        local workspace = deps.setup_environment.workspace_dir
        
        -- Generate comprehensive report
        local report = {
            workflow_id = "complex-workflow-" .. os.time(),
            execution_summary = {
                total_tasks = 5,
                successful_tasks = 0,
                failed_tasks = 0,
                skipped_tasks = 0
            },
            performance_metrics = {
                data_processed = deps.process_data and deps.process_data.total_files or 0,
                api_calls_made = 2,
                deployment_completed = deps.conditional_deployment ~= nil
            },
            resource_usage = {
                workspace_created = workspace,
                temporary_files_cleaned = true,
                circuit_breaker_stats = deps.api_integration.circuit_breaker_stats
            },
            generated_at = os.date()
        }
        
        -- Count successful tasks
        for task_name, result in pairs(deps) do
            if result and result ~= false then
                report.execution_summary.successful_tasks = report.execution_summary.successful_tasks + 1
            else
                report.execution_summary.failed_tasks = report.execution_summary.failed_tasks + 1
            end
        end
        
        -- Save report
        local report_json = data.to_json(report)
        fs.write_file(workspace .. "/final-report.json", report_json)
        
        -- Cleanup temporary files (but keep report)
        log.info("🗑️ Cleaning up temporary resources...")
        
        return true, "Cleanup and reporting completed", {
            report = report,
            report_file = workspace .. "/final-report.json",
            cleanup_completed = true
        }
    end)
    :artifacts({"final-report.json"})
    :on_success(function(params, output)
        log.info("📊 Final report generated: " .. output.report_file)
        log.info("✅ Cleanup completed successfully")
    end)
    :build()

-- Define the complex workflow with comprehensive configuration
workflow.define("complex_workflow_modern", {
    description = "Complex workflow demonstrating advanced Modern DSL patterns",
    version = "3.0.0",
    
    metadata = {
        author = "Sloth Runner Advanced Team",
        tags = {"complex", "async", "circuit-breaker", "conditional", "modern-dsl"},
        complexity = "advanced",
        estimated_duration = "20m",
        requirements = ["git", "docker", "curl"],
        features_demonstrated = [
            "Environment validation",
            "Parallel async processing", 
            "Circuit breaker patterns",
            "Conditional execution",
            "Artifact management",
            "Comprehensive reporting"
        ]
    },
    
    tasks = {
        setup_environment_task,
        process_data_task,
        api_integration_task,
        conditional_deploy_task,
        cleanup_task
    },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 3,
        create_workdir_before_run = true,
        clean_workdir_after_run = false,
        
        circuit_breaker = {
            failure_threshold = 5,
            recovery_timeout = "30s",
            success_threshold = 3
        },
        
        performance_monitoring = true,
        metrics_collection = true
    },
    
    on_start = function()
        log.info("🚀 Starting complex workflow with Modern DSL...")
        log.info("🎯 This workflow demonstrates advanced patterns and enterprise features")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Complex workflow completed successfully!")
            log.info("📈 Advanced Modern DSL patterns executed flawlessly")
            
            -- Log final statistics
            local total_tasks = 0
            local successful_tasks = 0
            for task_name, result in pairs(results) do
                total_tasks = total_tasks + 1
                if result.success then
                    successful_tasks = successful_tasks + 1
                end
            end
            
            log.info("📊 Final Statistics:")
            log.info("  - Total tasks: " .. total_tasks)
            log.info("  - Successful: " .. successful_tasks)
            log.info("  - Success rate: " .. math.floor((successful_tasks/total_tasks)*100) .. "%")
            
        else
            log.error("❌ Complex workflow failed!")
            log.error("🔧 Check individual task results and logs for debugging")
        end
        return true
    end
})
