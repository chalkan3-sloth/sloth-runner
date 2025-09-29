-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03

local json_test_task = task("test_json_parse_and_to_json")
local yaml_test_task = task("test_yaml_parse_and_to_yaml")

local json_test_task = task("test_json_parse_and_to_json")
local yaml_test_task = task("test_yaml_parse_and_to_yaml")
local json_test_task = task("test_json_parse_and_to_json")
    :description("Parses JSON and converts back with modern DSL")
    :command(function(params)
        log.info("Modern DSL: Testing JSON parse and to_json...")
        
        local json_str = '{"name": "TaskRunner", "version": 1.0, "active": true, "features": ["fs", "net", "data"]}'
        local parsed_data, err = data.parse_json(json_str)
        
        if err then
            return false, "Failed to parse JSON: " .. err
        end

        log.info("Parsed JSON name: " .. parsed_data.name)
        log.info("Parsed JSON version: " .. parsed_data.version)
        log.info("Parsed JSON active: " .. tostring(parsed_data.active))
        log.info("Parsed JSON features[1]: " .. parsed_data.features[1])

        -- Enhanced validation
        local validation_errors = {}
        if parsed_data.name ~= "TaskRunner" then
            table.insert(validation_errors, "name mismatch")
        end
        if parsed_data.version ~= 1.0 then
            table.insert(validation_errors, "version mismatch")
        end
        if parsed_data.active ~= true then
            table.insert(validation_errors, "active flag mismatch")
        end
        if parsed_data.features[1] ~= "fs" then
            table.insert(validation_errors, "features mismatch")
        end
        
        if #validation_errors > 0 then
            return false, "Parsed JSON data validation failed: " .. table.concat(validation_errors, ", ")
        end

        local new_json_str, err = data.to_json(parsed_data)
        if err then
            return false, "Failed to convert to JSON: " .. err
        end
        
        log.info("Converted back to JSON:\n" .. new_json_str)
        
        -- Enhanced validation
        if not string.find(new_json_str, "TaskRunner") or not string.find(new_json_str, "1") then
            return false, "Converted JSON data validation failed"
        end

        return true, "JSON operations successful", {
            original_size = #json_str,
            converted_size = #new_json_str,
            features_tested = {"parse_json", "to_json"}
        }
    end)
    :timeout("30s")
    :build()
local yaml_test_task = task("test_yaml_parse_and_to_yaml")
    :description("Parses YAML and converts back with modern DSL")
    :depends_on({"test_json_parse_and_to_json"})
    :command(function(params)
        log.info("Modern DSL: Testing YAML parse and to_yaml...")
        
        local yaml_str = [[ 
name: YAMLTest
version: 2.0
enabled: false
items:
  - item1
  - item2
config:
  key: value
]]
        local parsed_data, err = data.parse_yaml(yaml_str)
        if err then
            return false, "Failed to parse YAML: " .. err
        end

        log.info("Parsed YAML name: " .. parsed_data.name)
        log.info("Parsed YAML version: " .. parsed_data.version)
        log.info("Parsed YAML enabled: " .. tostring(parsed_data.enabled))
        log.info("Parsed YAML items[1]: " .. parsed_data.items[1])
        log.info("Parsed YAML config.key: " .. parsed_data.config.key)

        -- Enhanced validation
        local validation_errors = {}
        if parsed_data.name ~= "YAMLTest" then
            table.insert(validation_errors, "name mismatch")
        end
        if parsed_data.version ~= 2.0 then
            table.insert(validation_errors, "version mismatch")
        end
        if parsed_data.enabled ~= false then
            table.insert(validation_errors, "enabled flag mismatch")
        end
        if parsed_data.items[1] ~= "item1" then
            table.insert(validation_errors, "items mismatch")
        end
        if parsed_data.config.key ~= "value" then
            table.insert(validation_errors, "config mismatch")
        end
        
        if #validation_errors > 0 then
            return false, "Parsed YAML data validation failed: " .. table.concat(validation_errors, ", ")
        end

        local new_yaml_str, err = data.to_yaml(parsed_data)
        if err then
            return false, "Failed to convert to YAML: " .. err
        end
        
        log.info("Converted back to YAML:\n" .. new_yaml_str)
        
        -- Enhanced validation
        if not string.find(new_yaml_str, "YAMLTest") or not string.find(new_yaml_str, "item1") then
            return false, "Converted YAML data validation failed"
        end

        return true, "YAML operations successful", {
            original_size = #yaml_str,
            converted_size = #new_yaml_str,
            features_tested = {"parse_yaml", "to_yaml"}
        }
    end)
    :timeout("30s")
    :on_success(function(params, output)
        log.info("All data serialization tests completed successfully!")
    end)
    :build()

workflow.define("data_operations_test_modern", {
    description = "Data serialization/deserialization operations - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        category = "testing",
        tags = {"data", "json", "yaml", "serialization", "modern-dsl"}
    },
    
    tasks = {
        json_test_task,
        yaml_test_task
    },
    
    config = {
        max_parallel_tasks = 1, -- Sequential for testing
        timeout = "10m",
        retry_policy = "exponential"
    },
    
    on_complete = function(success, results)
        if success then
            log.info("All data operation tests passed!")
        else
            log.error("Some data operation tests failed!")
        end
        return true
    end
})
