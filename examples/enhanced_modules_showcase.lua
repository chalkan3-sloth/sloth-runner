-- MODERN DSL - Enhanced Modules Showcase
-- Demonstrating new HTTP, String, and Math modules

-- Task 1: API Testing with HTTP module
local api_test_task = task("api_test")
    :description("Test API endpoints using enhanced HTTP module")
    :command(function(params, deps)
        log.info("🌐 Testing HTTP module capabilities...")
        
        -- Test GET request
        local response = http.get("https://jsonplaceholder.typicode.com/posts/1")
        if not response.success then
            return false, "GET request failed", { error = response.status }
        end
        
        log.info("✅ GET request successful: " .. response.status)
        log.info("📄 Response body sample: " .. strings.truncate(response.body, 100))
        
        -- Test POST request with JSON data
        local post_response = http.post(
            "https://jsonplaceholder.typicode.com/posts",
            data.to_json({
                title = "Test Post",
                body = "This is a test post from Sloth Runner",
                userId = 1
            }),
            {
                ["Content-Type"] = "application/json"
            }
        )
        
        if not post_response.success then
            return false, "POST request failed", { error = post_response.status }
        end
        
        log.info("✅ POST request successful: " .. post_response.status)
        
        return true, "HTTP tests completed", {
            get_status = response.status_code,
            post_status = post_response.status_code,
            response_sample = strings.truncate(response.body, 50)
        }
    end)
    :timeout("30s")
    :retries(2, "exponential")
    :build()

-- Task 2: String Processing
local string_processing_task = task("string_processing")
    :description("Demonstrate string manipulation capabilities")
    :depends_on({"api_test"})
    :command(function(params, deps)
        log.info("🔤 Testing string processing capabilities...")
        
        local sample_text = "Hello, World! This is a Sample String for Processing."
        
        -- Basic string operations
        local upper_text = strings.upper(sample_text)
        local lower_text = strings.lower(sample_text)
        local trimmed = strings.trim("  " .. sample_text .. "  ")
        
        log.info("📝 Original: " .. sample_text)
        log.info("📝 Upper: " .. upper_text)
        log.info("📝 Lower: " .. lower_text)
        
        -- String validation
        local email_valid = strings.is_email("test@example.com")
        local url_valid = strings.is_url("https://example.com")
        local numeric_check = strings.is_numeric("12345")
        
        log.info("✉️  Email validation: " .. tostring(email_valid))
        log.info("🌐 URL validation: " .. tostring(url_valid))
        log.info("🔢 Numeric check: " .. tostring(numeric_check))
        
        -- Hashing
        local md5_hash = strings.md5(sample_text)
        local sha256_hash = strings.sha256(sample_text)
        
        log.info("🔒 MD5: " .. md5_hash)
        log.info("🔒 SHA256: " .. sha256_hash)
        
        -- Base64 encoding/decoding
        local encoded = strings.base64_encode(sample_text)
        local decoded = strings.base64_decode(encoded)
        
        log.info("🔐 Base64 encoded: " .. encoded)
        log.info("🔓 Base64 decoded: " .. decoded)
        
        return true, "String processing completed", {
            original = sample_text,
            md5 = md5_hash,
            sha256 = sha256_hash,
            base64 = encoded,
            email_valid = email_valid,
            url_valid = url_valid
        }
    end)
    :timeout("30s")
    :build()

-- Task 3: Mathematical Operations
local math_operations_task = task("math_operations")
    :description("Demonstrate mathematical capabilities")
    :depends_on({"string_processing"})
    :command(function(params, deps)
        log.info("🧮 Testing mathematical operations...")
        
        -- Generate some random numbers for testing
        math.seed(os.time())
        local numbers = {}
        for i = 1, 10 do
            numbers[i] = math.random_int(1, 100)
        end
        
        log.info("🎲 Generated numbers: " .. strings.join(numbers, ", "))
        
        -- Statistical operations
        local sum = math.sum(numbers)
        local mean = math.mean(numbers)
        local median = math.median(numbers)
        local std_dev = math.std_dev(numbers)
        local variance = math.variance(numbers)
        
        log.info("📊 Sum: " .. sum)
        log.info("📊 Mean: " .. string.format("%.2f", mean))
        log.info("📊 Median: " .. median)
        log.info("📊 Standard Deviation: " .. string.format("%.2f", std_dev))
        log.info("📊 Variance: " .. string.format("%.2f", variance))
        
        -- Trigonometric calculations
        local angle = math.pi / 4  -- 45 degrees in radians
        local sin_val = math.sin(angle)
        local cos_val = math.cos(angle)
        local tan_val = math.tan(angle)
        
        log.info("📐 sin(π/4): " .. string.format("%.4f", sin_val))
        log.info("📐 cos(π/4): " .. string.format("%.4f", cos_val))
        log.info("📐 tan(π/4): " .. string.format("%.4f", tan_val))
        
        -- Power operations
        local square_root = math.sqrt(16)
        local power_result = math.pow(2, 8)
        local cube_root = math.cbrt(27)
        
        log.info("⚡ √16 = " .. square_root)
        log.info("⚡ 2^8 = " .. power_result)
        log.info("⚡ ∛27 = " .. cube_root)
        
        return true, "Mathematical operations completed", {
            dataset_size = #numbers,
            sum = sum,
            mean = mean,
            median = median,
            std_dev = std_dev,
            sin_45_deg = sin_val,
            cos_45_deg = cos_val
        }
    end)
    :timeout("30s")
    :build()

-- Task 4: Advanced HTTP API Integration
local advanced_api_task = task("advanced_api")
    :description("Advanced HTTP operations with error handling")
    :depends_on({"math_operations"})
    :command(function(params, deps)
        log.info("🚀 Testing advanced HTTP features...")
        
        -- Custom HTTP request with timeout and headers
        local custom_response = http.request({
            method = "GET",
            url = "https://api.github.com/users/octocat",
            headers = {
                ["User-Agent"] = "Sloth-Runner/1.0",
                ["Accept"] = "application/vnd.github.v3+json"
            },
            timeout = "10s"
        })
        
        if not custom_response.success then
            return false, "GitHub API request failed", { 
                status = custom_response.status_code,
                error = custom_response.body 
            }
        end
        
        log.info("✅ GitHub API response: " .. custom_response.status)
        
        -- Parse JSON response if available
        local user_data = nil
        if custom_response.json then
            user_data = custom_response.json
            log.info("👤 GitHub user: " .. (user_data.login or "unknown"))
            log.info("📊 Public repos: " .. (user_data.public_repos or "0"))
        end
        
        -- String processing on API response
        local response_hash = strings.sha256(custom_response.body)
        local response_size = #custom_response.body
        
        log.info("📏 Response size: " .. response_size .. " bytes")
        log.info("🔒 Response hash: " .. strings.truncate(response_hash, 16) .. "...")
        
        return true, "Advanced API operations completed", {
            github_user = user_data and user_data.login or "unknown",
            public_repos = user_data and user_data.public_repos or 0,
            response_size = response_size,
            response_hash = strings.truncate(response_hash, 16)
        }
    end)
    :timeout("30s")
    :retries(3, "exponential")
    :on_success(function(params, output)
        log.info("🎉 All enhanced modules working perfectly!")
        log.info("📈 Processing summary:")
        log.info("  - HTTP requests: ✅")
        log.info("  - String operations: ✅")
        log.info("  - Mathematical calculations: ✅")
        log.info("  - API integration: ✅")
    end)
    :build()

-- Modern Workflow Definition
workflow.define("enhanced_modules_showcase", {
    description = "Showcase of Enhanced Modules - HTTP, Strings, and Math",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"showcase", "http", "strings", "math", "modern-dsl"},
        category = "demonstration",
        complexity = "intermediate",
        estimated_duration = "2m"
    },
    
    tasks = { 
        api_test_task,
        string_processing_task,
        math_operations_task,
        advanced_api_task
    },
    
    config = {
        timeout = "10m",
        retry_policy = "exponential",
        max_parallel_tasks = 1,  -- Sequential execution for demo
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting Enhanced Modules Showcase...")
        log.info("📦 Demonstrating: HTTP, Strings, and Math modules")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Enhanced Modules Showcase completed successfully!")
            log.info("🎯 All new modules are functional and ready for use")
            log.info("📚 Check the documentation for more advanced features")
        else
            log.error("❌ Enhanced Modules Showcase failed!")
            log.warn("🔍 Check the logs for specific module errors")
        end
        return true
    end
})