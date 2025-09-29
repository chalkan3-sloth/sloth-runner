# 📝 Modern DSL Best Practices

This guide provides best practices, patterns, and recommendations for writing effective workflows using the Modern DSL.

## 🎯 General Principles

### 1. **Clear and Descriptive Naming**

```lua
-- ❌ Poor naming
local t1 = task("t1"):command("npm run build"):build()
local t2 = task("t2"):command("npm test"):build()

-- ✅ Good naming
local build_frontend = task("build_frontend")
    :description("Build React frontend application")
    :command("npm run build:frontend")
    :build()

local run_unit_tests = task("run_unit_tests")
    :description("Execute Jest unit test suite")
    :command("npm run test:unit")
    :build()
```

### 2. **Comprehensive Documentation**

```lua
local deploy_to_production = task("deploy_to_production")
    :description("Deploy application to production Kubernetes cluster with health checks")
    :command(function(params, deps)
        -- Deploy using helm with production values
        local result = exec.run("helm upgrade --install myapp ./charts/myapp -f values.prod.yaml")
        
        if not result.success then
            return false, "Helm deployment failed: " .. result.stderr
        end
        
        -- Verify deployment health
        local health_check = k8s.wait_for_pods("app=myapp", "5m")
        if not health_check.ready then
            return false, "Pods not ready within timeout"
        end
        
        return true, "Production deployment successful", {
            release_name = "myapp",
            pods_ready = health_check.count,
            deployment_time = os.time()
        }
    end)
    :metadata({
        owner = "platform-team",
        runbook = "https://runbooks.company.com/production-deploy",
        escalation = "platform-oncall@company.com"
    })
    :build()
```

### 3. **Consistent Error Handling**

```lua
local process_data_task = task("process_data")
    :description("Process incoming data with comprehensive error handling")
    :command(function(params, deps)
        local input_data = deps.fetch_data.result
        
        -- Validate input
        if not input_data or #input_data == 0 then
            return false, "No input data received from fetch_data task"
        end
        
        -- Process with error handling
        local success, result, error = pcall(function()
            return data_processor.process(input_data, {
                format = params.output_format or "json",
                validation = true,
                sanitize = true
            })
        end)
        
        if not success then
            return false, "Data processing failed: " .. (error or "unknown error")
        end
        
        -- Validate output
        if not result or not result.processed_count then
            return false, "Processing completed but output validation failed"
        end
        
        log.info("Successfully processed " .. result.processed_count .. " records")
        
        return true, "Data processing completed", {
            processed_count = result.processed_count,
            output_file = result.output_file,
            processing_time = result.duration
        }
    end)
    :timeout("30m")
    :retries(2, "exponential")
    :on_failure(function(params, error)
        log.error("Data processing failed: " .. error)
        
        -- Send detailed failure notification
        notifications.send("slack", {
            channel = "#data-alerts",
            message = "🚨 Data processing pipeline failed\n" ..
                     "Error: " .. error .. "\n" ..
                     "Input size: " .. (params.input_size or "unknown") .. "\n" ..
                     "Runbook: https://runbooks.company.com/data-processing"
        })
    end)
    :build()
```

## 🏗️ Task Design Patterns

### Pattern 1: **Idempotent Tasks**

Design tasks that can be safely re-run:

```lua
local setup_database = task("setup_database")
    :description("Idempotent database setup with migration support")
    :command(function(params)
        -- Check if database exists
        local db_exists = database.exists(params.database_name)
        
        if not db_exists then
            log.info("Creating database: " .. params.database_name)
            local create_result = database.create(params.database_name)
            if not create_result.success then
                return false, "Failed to create database: " .. create_result.error
            end
        else
            log.info("Database already exists: " .. params.database_name)
        end
        
        -- Run migrations (idempotent)
        log.info("Running database migrations...")
        local migrate_result = database.migrate(params.database_name, {
            migrations_path = "./migrations",
            target_version = params.target_version
        })
        
        if not migrate_result.success then
            return false, "Migration failed: " .. migrate_result.error
        end
        
        return true, "Database setup completed", {
            database_name = params.database_name,
            current_version = migrate_result.current_version,
            migrations_applied = migrate_result.applied_count
        }
    end)
    :retries(3, "exponential")
    :build()
```

### Pattern 2: **Circuit Breaker for External Services**

```lua
local call_external_api = task("call_external_api")
    :description("Call external API with circuit breaker protection")
    :command(function(params)
        -- Use circuit breaker for external API calls
        local api_result = circuit.protect("payment_api", function()
            return net.http_post("https://api.payment.com/process", {
                headers = {
                    ["Authorization"] = "Bearer " .. utils.secret("payment_api_token"),
                    ["Content-Type"] = "application/json"
                },
                body = data.to_json(params.payment_data),
                timeout = "10s"
            })
        end, {
            failure_threshold = 5,
            recovery_timeout = "30s",
            half_open_max_calls = 3
        })
        
        if not api_result.success then
            -- Circuit might be open
            if api_result.circuit_open then
                return false, "Payment API circuit breaker is open - service unavailable"
            else
                return false, "Payment API call failed: " .. api_result.error
            end
        end
        
        -- Validate API response
        if api_result.status_code ~= 200 then
            return false, "Payment API returned error: " .. api_result.status_code
        end
        
        local response_data = data.parse_json(api_result.body)
        if not response_data.transaction_id then
            return false, "Invalid response from payment API - missing transaction_id"
        end
        
        return true, "Payment processed successfully", {
            transaction_id = response_data.transaction_id,
            amount = response_data.amount,
            status = response_data.status
        }
    end)
    :timeout("30s")
    :retries(3, "exponential")
    :build()
```

### Pattern 3: **Parallel Processing with Aggregation**

```lua
local parallel_data_processing = task("parallel_data_processing")
    :description("Process data in parallel and aggregate results")
    :command(function(params, deps)
        local input_files = deps.prepare_data.file_list
        
        log.info("Processing " .. #input_files .. " files in parallel...")
        
        -- Process files in parallel
        local results = async.parallel(
            table.map(input_files, function(file)
                return function()
                    return data_processor.process_file(file, {
                        format = params.output_format,
                        validation = true
                    })
                end
            end),
            {
                max_workers = params.max_workers or 4,
                timeout = "20m",
                fail_fast = false  -- Process all files even if some fail
            }
        )
        
        -- Aggregate results
        local successful_files = {}
        local failed_files = {}
        local total_records = 0
        
        for i, result in ipairs(results) do
            if result.success then
                table.insert(successful_files, {
                    file = input_files[i],
                    records = result.record_count
                })
                total_records = total_records + result.record_count
            else
                table.insert(failed_files, {
                    file = input_files[i],
                    error = result.error
                })
            end
        end
        
        -- Report results
        log.info("Processing completed:")
        log.info("  Successful files: " .. #successful_files)
        log.info("  Failed files: " .. #failed_files)
        log.info("  Total records processed: " .. total_records)
        
        if #failed_files > 0 then
            log.warn("Some files failed to process:")
            for _, failed in ipairs(failed_files) do
                log.warn("  " .. failed.file .. ": " .. failed.error)
            end
        end
        
        -- Decide if task should succeed or fail
        local success_rate = #successful_files / #input_files
        if success_rate < (params.min_success_rate or 0.8) then
            return false, "Processing failed - success rate " .. 
                   string.format("%.1f%%", success_rate * 100) .. 
                   " below threshold"
        end
        
        return true, "Parallel processing completed", {
            total_files = #input_files,
            successful_files = #successful_files,
            failed_files = #failed_files,
            total_records = total_records,
            success_rate = success_rate
        }
    end)
    :timeout("30m")
    :depends_on({"prepare_data"})
    :build()
```

## 🔄 Workflow Design Patterns

### Pattern 1: **Multi-Environment Deployment**

```lua
local deploy_to_environment = function(environment)
    return task("deploy_to_" .. environment)
        :description("Deploy application to " .. environment .. " environment")
        :command(function(params, deps)
            local build_info = deps.build_application
            
            log.info("Deploying to " .. environment .. " environment...")
            
            -- Environment-specific configuration
            local env_config = {
                staging = {
                    replicas = 2,
                    resources = {cpu = "100m", memory = "256Mi"},
                    ingress = "staging.example.com"
                },
                production = {
                    replicas = 5,
                    resources = {cpu = "500m", memory = "1Gi"},
                    ingress = "api.example.com"
                }
            }
            
            local config = env_config[environment]
            if not config then
                return false, "Unknown environment: " .. environment
            end
            
            -- Deploy with environment-specific settings
            local deploy_result = k8s.deploy({
                image = build_info.image_tag,
                namespace = environment,
                replicas = config.replicas,
                resources = config.resources,
                ingress_host = config.ingress
            })
            
            if not deploy_result.success then
                return false, "Deployment to " .. environment .. " failed: " .. deploy_result.error
            end
            
            -- Environment-specific health checks
            local health_timeout = environment == "production" and "10m" or "5m"
            local health_check = k8s.wait_for_rollout({
                deployment = "myapp",
                namespace = environment,
                timeout = health_timeout
            })
            
            if not health_check.ready then
                return false, "Health check failed for " .. environment .. " deployment"
            end
            
            return true, "Successfully deployed to " .. environment, {
                environment = environment,
                replicas_ready = health_check.ready_replicas,
                deployment_time = os.time(),
                endpoint = "https://" .. config.ingress
            }
        end)
        :condition(when("params.deploy_" .. environment .. " == true"))
        :timeout(environment == "production" and "20m" or "10m")
        :retries(environment == "production" and 1 or 2)
        :build()
end

-- Define environment-specific workflows
workflow.define("deploy_pipeline", {
    description = "Multi-environment deployment pipeline",
    version = "2.0.0",
    
    tasks = {
        build_application,
        run_tests,
        deploy_to_environment("staging"),
        deploy_to_environment("production")
    },
    
    config = {
        timeout = "1h",
        max_parallel_tasks = 2
    }
})
```

### Pattern 2: **Blue-Green Deployment**

```lua
local blue_green_deployment = task("blue_green_deploy")
    :description("Blue-green deployment with automatic rollback")
    :command(function(params, deps)
        local build_info = deps.build_application
        local current_env = k8s.get_active_environment("myapp")
        local target_env = current_env == "blue" and "green" or "blue"
        
        log.info("Current active environment: " .. current_env)
        log.info("Deploying to target environment: " .. target_env)
        
        -- Deploy to target environment
        local deploy_result = k8s.deploy({
            image = build_info.image_tag,
            environment = target_env,
            namespace = "production",
            replicas = 3
        })
        
        if not deploy_result.success then
            return false, "Deployment to " .. target_env .. " failed: " .. deploy_result.error
        end
        
        -- Wait for deployment to be ready
        local health_check = k8s.wait_for_rollout({
            deployment = "myapp-" .. target_env,
            namespace = "production",
            timeout = "10m"
        })
        
        if not health_check.ready then
            return false, "Target environment " .. target_env .. " not ready"
        end
        
        -- Run smoke tests against target environment
        local smoke_tests = testing.run_smoke_tests({
            endpoint = "http://myapp-" .. target_env .. ":8080",
            timeout = "5m"
        })
        
        if not smoke_tests.passed then
            log.error("Smoke tests failed, keeping current environment active")
            return false, "Smoke tests failed: " .. smoke_tests.error
        end
        
        -- Switch traffic to new environment
        log.info("Switching traffic from " .. current_env .. " to " .. target_env)
        local switch_result = k8s.switch_traffic({
            service = "myapp",
            from_env = current_env,
            to_env = target_env
        })
        
        if not switch_result.success then
            return false, "Traffic switch failed: " .. switch_result.error
        end
        
        -- Wait and verify new environment is stable
        sleep(30)  -- Allow some traffic to flow
        
        local stability_check = monitoring.check_stability({
            service = "myapp",
            environment = target_env,
            duration = "2m",
            error_rate_threshold = 0.01
        })
        
        if not stability_check.stable then
            log.error("New environment unstable, rolling back...")
            
            -- Rollback traffic
            k8s.switch_traffic({
                service = "myapp",
                from_env = target_env,
                to_env = current_env
            })
            
            return false, "Deployment unstable, rolled back: " .. stability_check.reason
        end
        
        -- Success - clean up old environment
        log.info("Deployment successful, cleaning up old environment")
        k8s.scale_down({
            deployment = "myapp-" .. current_env,
            replicas = 0
        })
        
        return true, "Blue-green deployment completed successfully", {
            previous_env = current_env,
            current_env = target_env,
            deployment_time = os.time(),
            image_deployed = build_info.image_tag
        }
    end)
    :depends_on({"build_application", "run_tests"})
    :timeout("30m")
    :on_failure(function(params, error)
        log.error("Blue-green deployment failed: " .. error)
        
        -- Send critical alert
        alerts.send("pagerduty", {
            severity = "critical",
            summary = "Blue-green deployment failed",
            details = error,
            runbook = "https://runbooks.company.com/blue-green-rollback"
        })
    end)
    :build()
```

## 📊 Monitoring and Observability Best Practices

### 1. **Comprehensive Metrics Collection**

```lua
workflow.define("monitored_pipeline", {
    description = "Pipeline with comprehensive monitoring",
    version = "2.0.0",
    
    tasks = { build_task, test_task, deploy_task },
    
    config = {
        monitoring = {
            metrics = {
                enabled = true,
                custom_metrics = {
                    "pipeline_duration_seconds",
                    "build_size_bytes", 
                    "test_coverage_percentage",
                    "deployment_success_rate"
                }
            },
            
            alerts = {
                enabled = true,
                rules = {
                    {
                        name = "pipeline_duration_high",
                        condition = "pipeline_duration_seconds > 1800",  -- 30 minutes
                        severity = "warning",
                        message = "Pipeline taking longer than expected"
                    },
                    {
                        name = "deployment_failure_rate_high",
                        condition = "deployment_success_rate < 0.95",
                        severity = "critical",
                        message = "Deployment success rate below 95%"
                    }
                }
            }
        }
    },
    
    on_start = function()
        metrics.start_timer("pipeline_duration")
        metrics.increment("pipeline_starts_total")
        return true
    end,
    
    on_complete = function(success, results)
        local duration = metrics.stop_timer("pipeline_duration")
        
        metrics.record_gauge("pipeline_duration_seconds", duration)
        
        if success then
            metrics.increment("pipeline_success_total")
        else
            metrics.increment("pipeline_failure_total")
        end
        
        return true
    end
})
```

### 2. **Structured Logging**

```lua
local structured_logging_task = task("process_with_logging")
    :description("Task with comprehensive structured logging")
    :command(function(params)
        local correlation_id = utils.uuid()
        
        log.info("Starting data processing", {
            correlation_id = correlation_id,
            input_size = params.input_size,
            processing_mode = params.mode,
            timestamp = os.time()
        })
        
        -- Processing with progress logging
        local total_items = params.input_size
        local processed_items = 0
        
        for i = 1, total_items do
            -- Process item
            local item_result = process_item(i)
            processed_items = processed_items + 1
            
            -- Log progress every 1000 items
            if i % 1000 == 0 then
                log.info("Processing progress", {
                    correlation_id = correlation_id,
                    processed_items = processed_items,
                    total_items = total_items,
                    progress_percentage = math.floor((processed_items / total_items) * 100),
                    items_per_second = calculate_rate(processed_items)
                })
            end
            
            if not item_result.success then
                log.error("Item processing failed", {
                    correlation_id = correlation_id,
                    item_id = i,
                    error = item_result.error,
                    retry_count = item_result.retry_count
                })
            end
        end
        
        log.info("Data processing completed", {
            correlation_id = correlation_id,
            total_items = total_items,
            processed_items = processed_items,
            success_rate = processed_items / total_items,
            duration_seconds = calculate_duration()
        })
        
        return true, "Processing completed", {
            correlation_id = correlation_id,
            processed_items = processed_items,
            success_rate = processed_items / total_items
        }
    end)
    :build()
```

## 🔐 Security Best Practices

### 1. **Secret Management**

```lua
local secure_deployment = task("secure_deploy")
    :description("Deployment with proper secret management")
    :command(function(params)
        -- Retrieve secrets securely
        local db_password = utils.secret("database_password")
        local api_key = utils.secret("external_api_key")
        local ssl_cert = utils.secret("ssl_certificate")
        
        if not db_password or not api_key or not ssl_cert then
            return false, "Required secrets not available"
        end
        
        -- Use secrets in deployment without logging them
        local deploy_result = k8s.deploy({
            image = params.image_tag,
            secrets = {
                DATABASE_PASSWORD = db_password,
                API_KEY = api_key,
                SSL_CERT = ssl_cert
            },
            security_context = {
                run_as_non_root = true,
                read_only_root_filesystem = true,
                capabilities = {
                    drop = {"ALL"}
                }
            }
        })
        
        -- Clear secrets from memory
        db_password = nil
        api_key = nil
        ssl_cert = nil
        
        if not deploy_result.success then
            return false, "Secure deployment failed: " .. deploy_result.error
        end
        
        return true, "Secure deployment completed", {
            deployment_id = deploy_result.deployment_id,
            security_scan_passed = true
        }
    end)
    :security({
        secrets_required = {"database_password", "external_api_key", "ssl_certificate"},
        rbac_role = "secure-deployer",
        audit_logging = true
    })
    :build()
```

### 2. **Input Validation**

```lua
local validated_task = task("process_user_input")
    :description("Process user input with comprehensive validation")
    :command(function(params)
        -- Validate required parameters
        validate.required(params.user_id, "user_id")
        validate.required(params.action, "action")
        
        -- Validate parameter types and formats
        validate.type(params.user_id, "number", "user_id")
        validate.pattern(params.action, "^[a-zA-Z0-9_]+$", "action")
        
        -- Validate parameter ranges
        validate.range(params.user_id, 1, 1000000, "user_id")
        validate.enum(params.action, {"create", "update", "delete"}, "action")
        
        -- Sanitize input
        local sanitized_input = utils.sanitize({
            user_id = params.user_id,
            action = params.action,
            data = params.data and utils.escape_html(params.data) or nil
        })
        
        -- Process with validated and sanitized input
        local result = user_processor.process(sanitized_input)
        
        if not result.success then
            return false, "Processing failed: " .. result.error
        end
        
        return true, "User input processed successfully", {
            user_id = sanitized_input.user_id,
            action = sanitized_input.action,
            result_id = result.id
        }
    end)
    :validation(function(params)
        -- Pre-execution validation
        if not params.user_id or not params.action then
            return false, "Missing required parameters"
        end
        return true
    end)
    :build()
```

## 🎯 Performance Optimization Best Practices

### 1. **Efficient Resource Usage**

```lua
local optimized_task = task("resource_optimized_processing")
    :description("Processing task optimized for resource usage")
    :command(function(params)
        -- Set resource limits
        process.set_memory_limit("2GB")
        process.set_cpu_limit("2 cores")
        
        -- Use streaming for large datasets
        local input_stream = data.open_stream(params.input_file)
        local output_stream = data.create_stream(params.output_file)
        
        local processed_count = 0
        local batch_size = 1000
        
        while true do
            local batch = input_stream:read_batch(batch_size)
            if not batch or #batch == 0 then
                break
            end
            
            -- Process batch efficiently
            local processed_batch = data_processor.process_batch(batch, {
                parallel_workers = 4,
                memory_efficient = true
            })
            
            -- Write results
            output_stream:write_batch(processed_batch)
            processed_count = processed_count + #processed_batch
            
            -- Memory cleanup
            batch = nil
            processed_batch = nil
            
            -- Yield control periodically
            if processed_count % 10000 == 0 then
                log.info("Processed " .. processed_count .. " records...")
                coroutine.yield()
            end
        end
        
        input_stream:close()
        output_stream:close()
        
        return true, "Processing completed efficiently", {
            processed_count = processed_count,
            memory_usage = process.get_memory_usage(),
            cpu_usage = process.get_cpu_usage()
        }
    end)
    :resources({
        cpu = "2 cores",
        memory = "2GB",
        disk = "10GB"
    })
    :build()
```

### 2. **Caching Strategies**

```lua
local cached_computation = task("cached_expensive_computation")
    :description("Expensive computation with intelligent caching")
    :command(function(params)
        local cache_key = "computation_" .. params.dataset_id .. "_" .. params.algorithm_version
        
        -- Check cache first
        local cached_result = cache.get(cache_key)
        if cached_result then
            log.info("Using cached result for " .. cache_key)
            return true, "Computation completed (cached)", cached_result
        end
        
        log.info("Cache miss, performing computation...")
        
        -- Perform expensive computation
        local start_time = os.time()
        local computation_result = expensive_algorithm.compute({
            dataset_id = params.dataset_id,
            algorithm_version = params.algorithm_version,
            parameters = params.computation_params
        })
        local computation_time = os.time() - start_time
        
        if not computation_result.success then
            return false, "Computation failed: " .. computation_result.error
        end
        
        -- Cache result with TTL
        local cache_ttl = computation_time > 300 and "1h" or "30m"  -- Longer cache for expensive computations
        cache.set(cache_key, computation_result.data, cache_ttl)
        
        log.info("Computation completed in " .. computation_time .. "s, cached with TTL " .. cache_ttl)
        
        return true, "Computation completed", {
            result = computation_result.data,
            computation_time = computation_time,
            cache_key = cache_key,
            cached = false
        }
    end)
    :build()
```

---

**Following these best practices will help you build robust, maintainable, and efficient workflows using the Modern DSL. Remember to adapt these patterns to your specific use cases and requirements!**