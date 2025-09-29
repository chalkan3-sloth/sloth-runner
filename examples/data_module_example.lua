-- Data Processing Module Examples

-- JSON operations
local user_data = {
    name = "John Doe",
    age = 30,
    email = "john@example.com",
    active = true
}

local json_str = data.json_encode(user_data)
print("JSON encoded:", json_str)

local pretty_json = data.json_pretty(user_data)
print("Pretty JSON:")
print(pretty_json)

local decoded = data.json_decode(json_str)
print("Decoded name:", decoded.name)

-- JSON validation
local valid, err = data.json_validate('{"valid": "json"}')
print("Valid JSON:", valid)

local invalid, err = data.json_validate('{"invalid": json}')
print("Invalid JSON:", invalid, err)

-- YAML operations
local yaml_str = data.yaml_encode(user_data)
print("YAML encoded:")
print(yaml_str)

local yaml_decoded = data.yaml_decode(yaml_str)
print("YAML decoded name:", yaml_decoded.name)

-- Data conversion
local yaml_to_json = data.yaml_to_json(yaml_str)
print("YAML to JSON:", yaml_to_json)

local json_to_yaml = data.json_to_yaml(json_str)
print("JSON to YAML:")
print(json_to_yaml)

-- CSV operations
local csv_data = "name,age,city\nJohn,30,New York\nJane,25,Los Angeles"
local parsed_csv = data.csv_parse(csv_data)
print("CSV parsed rows:", #parsed_csv)

local csv_table = {
    {"name", "age", "city"},
    {"Bob", "35", "Chicago"},
    {"Alice", "28", "Miami"}
}
local generated_csv = data.csv_generate(csv_table)
print("Generated CSV:")
print(generated_csv)

-- Data transformation
local table1 = {name = "John", age = 30}
local table2 = {city = "New York", age = 31}
local merged = data.deep_merge(table1, table2)
print("Merged data:", data.json_encode(merged))

-- Path operations
local nested_data = {
    user = {
        profile = {
            name = "John",
            settings = {
                theme = "dark"
            }
        }
    }
}

local theme = data.get_path(nested_data, "user.profile.settings.theme")
print("Theme from path:", theme)

data.set_path(nested_data, "user.profile.settings.language", "en")
print("After setting language:", data.json_encode(nested_data))

-- Flatten/unflatten
local flattened = data.flatten(nested_data)
print("Flattened:", data.json_encode(flattened))

local unflattened = data.unflatten(flattened)
print("Unflattened:", data.json_encode(unflattened))