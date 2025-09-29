-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:30 -03

local hello_task = task("hello_world")
local build_task = task("build_application")
local deploy_task = task("deploy_staging")
local config_task = task("setup_config")
local resource_task = task("process_large_dataset")

local hello_task = task("hello_world")
local build_task = task("build_application")
local deploy_task = task("deploy_staging")
local config_task = task("setup_config")
local resource_task = task("process_large_dataset")
local hello_task = task("hello_world")
    :description("A simple hello world task")
    :command("echo 'Hello from modern DSL!'")
    :timeout("30s")
    :build()
local build_task = task("build_application")
    :description("Build the application with modern features")
    :command(function(params, deps)
        log.info("Starting modern build process...")
        
        -- Use performance monitoring
        local result, duration = perf.measure(function()
            -- Parallel compilation
            local compile_results = async.parallel({
                frontend = function()
                    return exec.run("npm run build:frontend")
                end,
                backend = function()
                    return exec.run("go build -o app ./cmd/server")
                end
            }, 2) -- max 2 workers
            
            return compile_results
        end)
        
        return true, "Build completed", {
            artifacts = {"app", "dist/"},
            duration = duration
        }
    end)
    :depends_on({"setup_environment"})
    :async(true)
    :timeout("10m")
    :retries(2, "exponential")
    :artifacts({"app", "dist/"})
    :build()
local deploy_task = task("deploy_staging")
    :description("Deploy to staging environment")
    :command(function(params, deps)
        local build_info = deps.build_application
        
        -- Use circuit breaker for external calls
        local result = circuit.protect("staging_api", function()
            return exec.run("kubectl apply -f staging-deployment.yaml")
        end)
        
        return result.success, result.message
    end)
    :depends_on({"build_application"})
    :condition(when("env.STAGE == 'staging'"))
    :timeout("5m")
    :on_failure(function()
        -- Saga pattern - compensate on failure
        saga.compensate("deploy_rollback")
    end)
    :build()
local config_task = task("setup_config")
    :description("Setup configuration with modern utilities")
    :command(function(params)
        -- Use modern utilities
        local config = utils.config("app_config", "production")
        local secret = utils.secret("database_password")
        
        -- Template rendering
        local rendered = template.render("config.yaml.tmpl", {
            database_host = config.database.host,
            database_password = secret
        })
        
        fs.write("config.yaml", rendered)
        
        return true, "Configuration setup complete"
    end)
    :security_policy({
        sandbox = true,
        allowed_ops = {"fs.write", "template.render"}
    })
    :build()
local resource_task = task("process_large_dataset")
    :description("Process large dataset with resource management")
    :command(function(params)
        -- Allocate resources
        local cpu_resource = resource.allocate("cpu", "4 cores")
        local memory_resource = resource.allocate("memory", "8GB")
        
        local result = perf.measure(function()
            return exec.run("python process_dataset.py --input large_dataset.csv")
        end)
        
        -- Release resources
        resource.release(cpu_resource)
        resource.release(memory_resource)
        
        return result.success, result.message
    end)
    :resources({
        cpu = "4 cores",
        memory = "8GB",
        disk = "10GB"
    })
    :timeout("1h")
    :build()

workflow.define("ci_cd_pipeline", {
    description = "Complete CI/CD pipeline with modern features",
    version = "2.0.0",
    
    stages = {
        {
            name = "preparation",
            tasks = chain({
                "setup_workspace",
                "validate_environment",
                "load_secrets"
            })
        },
        
        {
            name = "build_and_test",
            tasks = parallel({
                "build_application",
                "run_tests",
                "security_scan"
            }, {
                max_workers = 3,
                fail_fast = true
            })
        },
        
        {
            name = "deployment",
            condition = when("test.success && build.success")
                :then("deploy_staging")
                :else("notify_failure"),
            
            tasks = {
                "deploy_staging",
                "smoke_test",
                "promote_production"
            }
        }
    }
})
