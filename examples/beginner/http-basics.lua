-- MODERN DSL ONLY - HTTP Basics Example
-- Demonstrates HTTP operations with Modern DSL

-- HTTP GET task with circuit breaker
local http_get_task = task("http_get")
    :description("Perform HTTP GET with circuit breaker protection")
    :command(function(params)
        log.info("🌐 Making HTTP GET request...")
        
        -- Use circuit breaker for external API
        local result = circuit.protect("http_api", function()
            return net.http_get("https://jsonplaceholder.typicode.com/posts/1")
        end)
        
        if result.success then
            log.info("✅ HTTP request successful")
            return true, "HTTP GET completed", {
                status_code = result.status_code,
                body = result.body,
                headers = result.headers
            }
        else
            return false, "HTTP request failed: " .. (result.error or "unknown error")
        end
    end)
    :timeout("30s")
    :retries(3, "exponential")
    :on_success(function(params, output)
        log.info("📊 Response received: " .. string.len(output.body or "") .. " bytes")
    end)
    :build()

-- HTTP POST task
local http_post_task = task("http_post")
    :description("Perform HTTP POST with data")
    :depends_on({"http_get"})
    :command(function(params, deps)
        log.info("📤 Making HTTP POST request...")
        
        local post_data = {
            title = "Modern DSL Post",
            body = "Posted from Sloth Runner Modern DSL",
            userId = 1
        }
        
        local result = net.http_post("https://jsonplaceholder.typicode.com/posts", {
            headers = { ["Content-Type"] = "application/json" },
            body = data.to_json(post_data)
        })
        
        if result.success then
            return true, "HTTP POST completed", {
                post_result = result.body,
                post_status = result.status_code
            }
        else
            return false, "HTTP POST failed"
        end
    end)
    :timeout("45s")
    :build()

-- Modern Workflow Definition
workflow.define("http_basics", {
    description = "HTTP Operations - Modern DSL",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"http", "api", "beginner", "modern-dsl"},
        created_at = os.date()
    },
    
    tasks = { http_get_task, http_post_task },
    
    config = {
        timeout = "10m",
        max_parallel_tasks = 1,
        circuit_breaker = {
            failure_threshold = 3,
            recovery_timeout = "1m"
        }
    },
    
    on_start = function()
        log.info("🚀 Starting HTTP basics workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ HTTP workflow completed successfully!")
        else
            log.error("❌ HTTP workflow failed!")
        end
        return true
    end
})
