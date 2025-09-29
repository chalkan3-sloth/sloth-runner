# 🚀 New Modules Enhancement

We've expanded Sloth Runner with powerful new modules that enhance automation capabilities and provide enterprise-grade functionality.

## 🌐 HTTP Module

The new HTTP module provides comprehensive HTTP client functionality with advanced features:

### Features
- **Full HTTP Methods**: GET, POST, PUT, DELETE, PATCH
- **Custom Headers**: Full header customization
- **JSON Support**: Automatic JSON parsing and serialization
- **Timeout Control**: Configurable request timeouts
- **Response Processing**: Structured response handling

### Usage Examples

```lua
-- Simple GET request
local response = http.get("https://api.example.com/data")
if response.success then
    log.info("Response: " .. response.body)
end

-- POST with JSON data
local post_response = http.post(
    "https://api.example.com/users",
    data.to_json({
        name = "John Doe",
        email = "john@example.com"
    }),
    {
        ["Content-Type"] = "application/json",
        ["Authorization"] = "Bearer " .. token
    }
)

-- Advanced request with custom options
local custom_response = http.request({
    method = "PUT",
    url = "https://api.example.com/users/123",
    body = {
        name = "Updated Name",
        status = "active"
    },
    headers = {
        ["Authorization"] = "Bearer " .. token
    },
    timeout = "30s"
})
```

## 🔤 String Module

Advanced string processing capabilities with validation, encoding, and manipulation:

### Features
- **String Manipulation**: trim, upper, lower, split, join, replace
- **Regular Expressions**: match, match_all, replace_regex
- **Encoding/Decoding**: Base64, URL encoding
- **Hashing**: MD5, SHA1, SHA256
- **Validation**: Email, URL, numeric, alphabetic checks
- **Formatting**: Padding, truncation

### Usage Examples

```lua
-- Basic string operations
local text = "  Hello, World!  "
local clean_text = strings.trim(text)
local upper_text = strings.upper(clean_text)
local words = strings.split(clean_text, " ")

-- Validation
local is_valid_email = strings.is_email("user@example.com")
local is_valid_url = strings.is_url("https://example.com")
local is_numeric = strings.is_numeric("12345")

-- Hashing and encoding
local md5_hash = strings.md5("sensitive data")
local sha256_hash = strings.sha256("important content")
local encoded = strings.base64_encode("secret message")

-- Regular expressions
local matches = strings.match("Version 1.2.3", "Version (%d+)%.(%d+)%.(%d+)")
if matches then
    log.info("Major: " .. matches[2] .. ", Minor: " .. matches[3])
end

-- Text formatting
local padded = strings.pad_left("42", 5, "0")  -- "00042"
local truncated = strings.truncate("Very long text here", 10, "...")  -- "Very lo..."
```

## 🧮 Math Module

Comprehensive mathematical operations and statistical functions:

### Features
- **Basic Math**: abs, ceil, floor, round, min, max, clamp
- **Power & Roots**: pow, sqrt, cbrt
- **Trigonometry**: sin, cos, tan, asin, acos, atan, atan2
- **Logarithms**: log, log10, log2, exp
- **Random Numbers**: random, random_int, random_float with seeding
- **Statistics**: sum, mean, median, mode, variance, standard deviation
- **Constants**: pi, e, phi (golden ratio)

### Usage Examples

```lua
-- Basic operations
local absolute = math.abs(-42)  -- 42
local maximum = math.max(10, 20, 5, 30)  -- 30
local clamped = math.clamp(15, 0, 10)  -- 10

-- Statistical analysis
local dataset = {10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
local mean_value = math.mean(dataset)  -- 55
local median_value = math.median(dataset)  -- 55
local std_deviation = math.std_dev(dataset)  -- ~30.28

-- Random number generation
math.seed(os.time())  -- Seed with current time
local random_int = math.random_int(1, 100)  -- Random integer 1-100
local random_float = math.random_float(0.0, 1.0)  -- Random float 0.0-1.0

-- Trigonometry
local angle = math.pi / 4  -- 45 degrees in radians
local sine = math.sin(angle)  -- ~0.7071
local cosine = math.cos(angle)  -- ~0.7071

-- Power and logarithms
local square_root = math.sqrt(16)  -- 4
local power = math.pow(2, 8)  -- 256
local natural_log = math.log(math.e)  -- 1
```

## 🔧 Integration with Modern DSL

All new modules integrate seamlessly with the Modern DSL:

```lua
local api_processing_task = task("api_data_processing")
    :description("Process API data with enhanced modules")
    :command(function(params, deps)
        -- Fetch data with HTTP module
        local response = http.get("https://api.example.com/metrics")
        if not response.success then
            return false, "API request failed"
        end
        
        -- Process response with string module
        local data_hash = strings.sha256(response.body)
        local clean_data = strings.trim(response.body)
        
        -- Parse and analyze with math module
        local data = data.from_json(clean_data)
        if data.values then
            local mean_val = math.mean(data.values)
            local std_dev = math.std_dev(data.values)
            
            log.info("Data analysis complete:")
            log.info("  Mean: " .. mean_val)
            log.info("  Std Dev: " .. std_dev)
            log.info("  Hash: " .. strings.truncate(data_hash, 16))
        end
        
        return true, "Data processed successfully", {
            hash = data_hash,
            mean = mean_val,
            std_dev = std_dev
        }
    end)
    :timeout("60s")
    :retries(3, "exponential")
    :build()
```

## 🎯 Use Cases

### 1. API Testing and Integration
- Automated API endpoint testing
- Data synchronization between services
- Webhook processing and validation
- External service health checks

### 2. Data Processing Pipelines
- Text processing and normalization
- Data validation and cleanup
- Statistical analysis of metrics
- Content transformation workflows

### 3. Security and Monitoring
- Hash verification for integrity checks
- Encoded data processing
- Performance metrics calculation
- Alert threshold analysis

### 4. DevOps Automation
- Configuration file processing
- Log analysis with pattern matching
- Metric aggregation and reporting
- Automated testing with mathematical validation

## 📚 Documentation

Each module includes comprehensive error handling and follows Modern DSL patterns:

- **Consistent Return Values**: All functions return appropriate types or error information
- **Integration Ready**: Works seamlessly with existing modules (state, exec, fs, etc.)
- **Performance Optimized**: Efficient implementations suitable for production use
- **Type Safe**: Proper type checking and validation

## 🚀 Getting Started

To use the new modules, simply update to the latest version of Sloth Runner. The modules are automatically available in all Lua scripts:

```lua
-- No imports needed - modules are globally available
local response = http.get("https://example.com")
local processed = strings.upper(response.body)
local random_num = math.random_int(1, 100)
```

Check out the [enhanced modules showcase example](../examples/enhanced_modules_showcase.lua) for a complete demonstration of all features.

---

**These new modules significantly expand Sloth Runner's capabilities, making it a more powerful platform for automation, data processing, and integration workflows!**