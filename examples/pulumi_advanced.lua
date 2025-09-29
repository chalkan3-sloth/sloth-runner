-- Pulumi Advanced Examples
-- This file demonstrates advanced usage of the pulumi module

-- Example 1: Basic Stack Operations
local function create_and_manage_stack()
    local pulumi = require("pulumi")
    
    -- Create a new stack
    local success, err = pulumi.new_stack("dev-stack", {
        workdir = "./infrastructure",
        template = "typescript"
    })
    
    if not success then
        print("Error creating stack: " .. err)
        return false
    end
    
    print("✅ Stack created successfully")
    
    -- Select the stack
    local select_success, err = pulumi.select_stack("dev-stack", {
        workdir = "./infrastructure"
    })
    
    if not select_success then
        print("Error selecting stack: " .. err)
        return false
    end
    
    print("✅ Stack selected successfully")
    return true
end

-- Example 2: Configuration Management
local function manage_configuration()
    local pulumi = require("pulumi")
    
    -- Set configuration values
    local configs = {
        {"aws:region", "us-west-2"},
        {"app:environment", "development"},
        {"app:version", "1.0.0"}
    }
    
    for _, config in ipairs(configs) do
        local success, err = pulumi.config_set(config[1], config[2], {
            workdir = "./infrastructure"
        })
        
        if not success then
            print("Error setting config " .. config[1] .. ": " .. err)
            return false
        end
        
        print("Set config: " .. config[1] .. " = " .. config[2])
    end
    
    -- Set a secret configuration
    local secret_success, err = pulumi.config_set("database:password", "super-secret-password", {
        workdir = "./infrastructure",
        secret = true
    })
    
    if not secret_success then
        print("Error setting secret config: " .. err)
        return false
    end
    
    print("✅ Secret configuration set")
    
    -- Get configuration value
    local region, err = pulumi.config_get("aws:region", {
        workdir = "./infrastructure"
    })
    
    if err then
        print("Error getting config: " .. err)
        return false
    end
    
    print("Current AWS region: " .. region)
    return true
end

-- Example 3: Preview and Deploy
local function preview_and_deploy()
    local pulumi = require("pulumi")
    
    -- Preview changes
    print("🔍 Previewing infrastructure changes...")
    local preview_result, err = pulumi.preview({
        workdir = "./infrastructure",
        diff = true,
        refresh = true
    })
    
    if err then
        print("Error during preview: " .. err)
        return false
    end
    
    print("Preview result: " .. preview_result)
    
    -- Deploy infrastructure
    print("🚀 Deploying infrastructure...")
    local deploy_success, deploy_result = pulumi.up({
        workdir = "./infrastructure",
        yes = true,
        parallel = 10,
        refresh = true
    })
    
    if not deploy_success then
        print("Error during deployment: " .. deploy_result)
        return false
    end
    
    print("✅ Deployment successful!")
    print("Deployment result: " .. deploy_result)
    return true
end

-- Example 4: Multi-Stack Management
local function manage_multiple_stacks()
    local pulumi = require("pulumi")
    
    local environments = {"dev", "staging", "prod"}
    
    for _, env in ipairs(environments) do
        local stack_name = "app-" .. env
        
        -- Create stack
        local success, err = pulumi.new_stack(stack_name, {
            workdir = "./infrastructure"
        })
        
        if not success and not string.find(err, "already exists") then
            print("Error creating stack " .. stack_name .. ": " .. err)
            return false
        end
        
        -- Select stack
        pulumi.select_stack(stack_name, {workdir = "./infrastructure"})
        
        -- Set environment-specific configuration
        pulumi.config_set("app:environment", env, {workdir = "./infrastructure"})
        
        if env == "prod" then
            pulumi.config_set("app:instance_count", "3", {workdir = "./infrastructure"})
            pulumi.config_set("app:instance_type", "m5.large", {workdir = "./infrastructure"})
        else
            pulumi.config_set("app:instance_count", "1", {workdir = "./infrastructure"})
            pulumi.config_set("app:instance_type", "t3.micro", {workdir = "./infrastructure"})
        end
        
        print("✅ Configured stack: " .. stack_name)
    end
    
    -- List all stacks
    local stacks_json, err = pulumi.list_stacks({workdir = "./infrastructure"})
    if err then
        print("Error listing stacks: " .. err)
        return false
    end
    
    print("Available stacks: " .. stacks_json)
    return true
end

-- Example 5: Plugin Management
local function manage_plugins()
    local pulumi = require("pulumi")
    
    -- Install required plugins
    local plugins = {
        {"resource", "aws", "6.0.0"},
        {"resource", "kubernetes", "4.0.0"},
        {"resource", "docker", "4.0.0"}
    }
    
    for _, plugin in ipairs(plugins) do
        local success, err = pulumi.plugin_install(plugin[1], plugin[2], {
            workdir = "./infrastructure",
            version = plugin[3]
        })
        
        if not success then
            print("Error installing plugin " .. plugin[2] .. ": " .. err)
        else
            print("✅ Installed plugin: " .. plugin[2] .. " v" .. plugin[3])
        end
    end
    
    -- List installed plugins
    local plugins_json, err = pulumi.plugin_ls({workdir = "./infrastructure"})
    if err then
        print("Error listing plugins: " .. err)
        return false
    end
    
    print("Installed plugins: " .. plugins_json)
    return true
end

-- Example 6: Output Management
local function manage_outputs()
    local pulumi = require("pulumi")
    
    -- Get all outputs
    local outputs_json, err = pulumi.outputs({workdir = "./infrastructure"})
    if err then
        print("Error getting outputs: " .. err)
        return false
    end
    
    print("Stack outputs: " .. outputs_json)
    
    -- Parse outputs and use them
    local data = require("data")
    local outputs = data.parse_json(outputs_json)
    
    if outputs and outputs.website_url then
        print("Website URL: " .. outputs.website_url.value)
    end
    
    if outputs and outputs.database_endpoint then
        print("Database endpoint: " .. outputs.database_endpoint.value)
    end
    
    return true
end

-- Example 7: Stack History and Rollback
local function manage_history()
    local pulumi = require("pulumi")
    
    -- Get stack history
    local history_json, err = pulumi.history({workdir = "./infrastructure"})
    if err then
        print("Error getting history: " .. err)
        return false
    end
    
    print("Stack history: " .. history_json)
    
    -- Parse history to find the last successful update
    local data = require("data")
    local history = data.parse_json(history_json)
    
    if history and #history > 1 then
        local last_update = history[1]
        print("Last update: " .. last_update.kind .. " at " .. last_update.startTime)
        
        if last_update.result == "failed" then
            print("⚠️  Last update failed, consider rolling back")
        end
    end
    
    return true
end

-- Example 8: Export and Import State
local function backup_and_restore()
    local pulumi = require("pulumi")
    
    -- Export stack state
    local state_json, err = pulumi.export({
        workdir = "./infrastructure",
        file = "./backups/stack-backup.json"
    })
    
    if err then
        print("Error exporting state: " .. err)
        return false
    end
    
    print("✅ State exported successfully")
    
    -- The import functionality would be used in disaster recovery scenarios
    -- pulumi.import("./backups/stack-backup.json", {workdir = "./infrastructure"})
    
    return true
end

-- Example 9: Policy Management (Pulumi Cloud)
local function manage_policies()
    local pulumi = require("pulumi")
    
    -- Create a new policy pack
    local success, err = pulumi.policy_new("security-policies", {
        workdir = "./policies",
        template = "typescript"
    })
    
    if not success and not string.find(err, "already exists") then
        print("Error creating policy pack: " .. err)
        return false
    end
    
    print("✅ Policy pack created")
    
    -- Enable policy pack for organization (requires Pulumi Cloud)
    -- pulumi.policy_enable("security-policies", "my-org", {workdir = "./policies"})
    
    return true
end

-- Example 10: Complete Infrastructure Lifecycle
local function complete_infrastructure_lifecycle()
    local pulumi = require("pulumi")
    
    print("🚀 Starting complete infrastructure lifecycle...")
    
    -- Step 1: Setup
    print("Step 1: Setting up stack...")
    if not create_and_manage_stack() then
        return false
    end
    
    -- Step 2: Configuration
    print("Step 2: Managing configuration...")
    if not manage_configuration() then
        return false
    end
    
    -- Step 3: Plugin management
    print("Step 3: Installing plugins...")
    if not manage_plugins() then
        return false
    end
    
    -- Step 4: Preview changes
    print("Step 4: Previewing changes...")
    local preview_result, err = pulumi.preview({
        workdir = "./infrastructure",
        diff = true
    })
    
    if err then
        print("Preview failed: " .. err)
        return false
    end
    
    -- Step 5: Deploy (in non-production scenarios)
    print("Step 5: Deployment (preview only in this example)...")
    print("To deploy, run: pulumi.up({workdir = './infrastructure', yes = true})")
    
    -- Step 6: Monitor outputs
    print("Step 6: Managing outputs...")
    manage_outputs()
    
    -- Step 7: Backup state
    print("Step 7: Backing up state...")
    backup_and_restore()
    
    print("✅ Infrastructure lifecycle completed!")
    return true
end

-- Example 11: Refresh and Cleanup
local function refresh_and_cleanup()
    local pulumi = require("pulumi")
    
    -- Refresh stack state
    print("🔄 Refreshing stack state...")
    local refresh_success, err = pulumi.refresh({
        workdir = "./infrastructure",
        yes = true
    })
    
    if not refresh_success then
        print("Error refreshing stack: " .. err)
        return false
    end
    
    print("✅ Stack refreshed")
    
    -- Destroy resources (only in dev/test environments)
    local env = os.getenv("ENVIRONMENT") or "dev"
    if env == "dev" or env == "test" then
        print("🗑️  Destroying development resources...")
        local destroy_success, err = pulumi.destroy({
            workdir = "./infrastructure",
            yes = true,
            skip_preview = true
        })
        
        if not destroy_success then
            print("Error destroying resources: " .. err)
            return false
        end
        
        print("✅ Resources destroyed")
    else
        print("ℹ️  Skipping destruction in " .. env .. " environment")
    end
    
    return true
end

-- Export task definitions
TaskDefinitions = {
    pulumi_examples = {
        description = "Advanced Pulumi infrastructure management examples",
        workdir = ".",
        tasks = {
            create_stack = {
                description = "Create and manage Pulumi stack",
                command = create_and_manage_stack
            },
            manage_config = {
                description = "Manage Pulumi configuration",
                command = manage_configuration
            },
            preview_deploy = {
                description = "Preview and deploy infrastructure",
                command = preview_and_deploy
            },
            multi_stack = {
                description = "Manage multiple stacks",
                command = manage_multiple_stacks
            },
            manage_plugins = {
                description = "Install and manage plugins",
                command = manage_plugins
            },
            manage_outputs = {
                description = "Manage stack outputs",
                command = manage_outputs
            },
            manage_history = {
                description = "View stack history",
                command = manage_history
            },
            backup_restore = {
                description = "Export and import stack state",
                command = backup_and_restore
            },
            manage_policies = {
                description = "Create and manage policy packs",
                command = manage_policies
            },
            complete_lifecycle = {
                description = "Complete infrastructure lifecycle",
                command = complete_infrastructure_lifecycle
            },
            refresh_cleanup = {
                description = "Refresh state and cleanup resources",
                command = refresh_and_cleanup
            }
        }
    }
}