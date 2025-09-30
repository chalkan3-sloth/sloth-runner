-- Terraform Advanced Examples
-- This file demonstrates advanced usage of the terraform module

-- Example 1: Basic Terraform Operations
local function basic_terraform_operations()
    local terraform = require("terraform")
    
    -- Initialize Terraform
    print("🔧 Initializing Terraform...")
    local init_success, err = terraform.init({
        workdir = "./terraform",
        upgrade = true
    })
    
    if not init_success then
        print("Error initializing Terraform: " .. err)
        return false
    end
    
    print("✅ Terraform initialized successfully")
    
    -- Validate configuration
    print("🔍 Validating Terraform configuration...")
    local valid, err = terraform.validate({
        workdir = "./terraform",
        json = true
    })
    
    if not valid then
        print("Terraform validation failed: " .. err)
        return false
    end
    
    print("✅ Terraform configuration is valid")
    
    -- Format code
    print("📝 Formatting Terraform code...")
    local fmt_success, err = terraform.fmt({
        workdir = "./terraform",
        recursive = true,
        diff = true
    })
    
    if not fmt_success then
        print("Error formatting code: " .. err)
        return false
    end
    
    print("✅ Code formatted successfully")
    return true
end

-- Example 2: Plan and Apply
local function plan_and_apply()
    local terraform = require("terraform")
    
    -- Create execution plan
    print("📋 Creating Terraform execution plan...")
    local plan_result, err = terraform.plan({
        workdir = "./terraform",
        out = "tfplan",
        var_file = "terraform.tfvars",
        detailed_exitcode = true
    })
    
    if err then
        print("Error creating plan: " .. err)
        return false
    end
    
    print("Plan result: " .. plan_result)
    
    -- Apply the plan
    print("🚀 Applying Terraform plan...")
    local apply_success, err = terraform.apply({
        workdir = "./terraform",
        plan = "tfplan"
    })
    
    if not apply_success then
        print("Error applying plan: " .. err)
        return false
    end
    
    print("✅ Terraform apply completed successfully")
    return true
end

-- Example 3: Variable Management
local function manage_variables()
    local terraform = require("terraform")
    
    -- Define variables
    local variables = {
        environment = "development",
        instance_count = "2",
        instance_type = "t3.micro",
        region = "us-west-2"
    }
    
    -- Create plan with variables
    local plan_result, err = terraform.plan({
        workdir = "./terraform",
        vars = variables,
        out = "dev-plan"
    })
    
    if err then
        print("Error creating plan with variables: " .. err)
        return false
    end
    
    print("✅ Plan created with variables")
    
    -- Apply with auto-approve for development
    local apply_success, err = terraform.apply({
        workdir = "./terraform",
        vars = variables,
        auto_approve = true,
        parallelism = 5
    })
    
    if not apply_success then
        print("Error applying with variables: " .. err)
        return false
    end
    
    print("✅ Applied successfully with variables")
    return true
end

-- Example 4: Workspace Management
local function manage_workspaces()
    local terraform = require("terraform")
    
    -- List existing workspaces
    local workspaces, err = terraform.workspace_list({workdir = "./terraform"})
    if err then
        print("Error listing workspaces: " .. err)
        return false
    end
    
    print("Available workspaces: " .. workspaces)
    
    -- Create development workspace
    local dev_success, err = terraform.workspace_new("development", {workdir = "./terraform"})
    if not dev_success and not string.find(err, "already exists") then
        print("Error creating development workspace: " .. err)
        return false
    end
    
    -- Create staging workspace
    local staging_success, err = terraform.workspace_new("staging", {workdir = "./terraform"})
    if not staging_success and not string.find(err, "already exists") then
        print("Error creating staging workspace: " .. err)
        return false
    end
    
    -- Select development workspace
    local select_success, err = terraform.workspace_select("development", {workdir = "./terraform"})
    if not select_success then
        print("Error selecting development workspace: " .. err)
        return false
    end
    
    print("✅ Workspaces managed successfully")
    return true
end

-- Example 5: State Management
local function manage_state()
    local terraform = require("terraform")
    
    -- List resources in state
    local resources, err = terraform.state_list({workdir = "./terraform"})
    if err then
        print("Error listing state resources: " .. err)
        return false
    end
    
    print("Resources in state: " .. resources)
    
    -- Show specific resource
    if string.find(resources, "aws_instance") then
        local resource_details, err = terraform.state_show("aws_instance.web[0]", {workdir = "./terraform"})
        if err then
            print("Error showing resource details: " .. err)
        else
            print("Resource details: " .. resource_details)
        end
    end
    
    -- Pull state
    local state_json, err = terraform.state_pull({workdir = "./terraform"})
    if err then
        print("Error pulling state: " .. err)
        return false
    end
    
    print("✅ State operations completed")
    return true
end

-- Example 6: Output Management
local function manage_outputs()
    local terraform = require("terraform")
    
    -- Get all outputs
    local outputs_json, err = terraform.output({
        workdir = "./terraform",
        json = true
    })
    
    if err then
        print("Error getting outputs: " .. err)
        return false
    end
    
    print("All outputs: " .. outputs_json)
    
    -- Get specific output
    local website_url, err = terraform.output({
        workdir = "./terraform",
        name = "website_url",
        raw = true
    })
    
    if err then
        print("Error getting website URL: " .. err)
    else
        print("Website URL: " .. website_url)
    end
    
    -- Parse JSON outputs
    local data = require("data")
    local outputs = data.parse_json(outputs_json)
    
    if outputs then
        for key, value in pairs(outputs) do
            if value.value then
                print("Output " .. key .. ": " .. tostring(value.value))
            end
        end
    end
    
    return true
end

-- Example 7: Import Existing Resources
local function import_resources()
    local terraform = require("terraform")
    
    -- Import existing AWS instance
    local import_success, err = terraform.import(
        "aws_instance.existing_server",
        "i-1234567890abcdef0",
        {
            workdir = "./terraform",
            var_file = "terraform.tfvars"
        }
    )
    
    if not import_success then
        print("Error importing resource: " .. err)
        return false
    end
    
    print("✅ Resource imported successfully")
    
    -- Refresh state to sync with reality
    local refresh_success, err = terraform.refresh({
        workdir = "./terraform",
        var_file = "terraform.tfvars"
    })
    
    if not refresh_success then
        print("Error refreshing state: " .. err)
        return false
    end
    
    print("✅ State refreshed successfully")
    return true
end

-- Example 8: Provider Management
local function manage_providers()
    local terraform = require("terraform")
    
    -- List providers
    local providers, err = terraform.providers({workdir = "./terraform"})
    if err then
        print("Error listing providers: " .. err)
        return false
    end
    
    print("Providers: " .. providers)
    
    -- Lock provider versions
    local lock_success, err = terraform.providers_lock({
        workdir = "./terraform",
        platform = "linux_amd64"
    })
    
    if not lock_success then
        print("Error locking providers: " .. err)
        return false
    end
    
    print("✅ Provider versions locked")
    
    -- Mirror providers for offline use
    local mirror_success, err = terraform.providers_mirror("./provider-cache", {
        workdir = "./terraform",
        platform = "linux_amd64"
    })
    
    if not mirror_success then
        print("Error mirroring providers: " .. err)
        return false
    end
    
    print("✅ Providers mirrored successfully")
    return true
end

-- Example 9: Taint and Untaint Resources
local function taint_operations()
    local terraform = require("terraform")
    
    -- Taint a resource for recreation
    local taint_success, err = terraform.taint("aws_instance.web[0]", {workdir = "./terraform"})
    if not taint_success then
        print("Error tainting resource: " .. err)
        return false
    end
    
    print("✅ Resource tainted successfully")
    
    -- Create plan to see tainted resource
    local plan_result, err = terraform.plan({
        workdir = "./terraform",
        detailed_exitcode = true
    })
    
    if err then
        print("Error creating plan: " .. err)
    else
        print("Plan with tainted resource: " .. plan_result)
    end
    
    -- Untaint the resource (if we don't want to recreate it)
    local untaint_success, err = terraform.untaint("aws_instance.web[0]", {workdir = "./terraform"})
    if not untaint_success then
        print("Error untainting resource: " .. err)
        return false
    end
    
    print("✅ Resource untainted successfully")
    return true
end

-- Example 10: Graph Generation
local function generate_graph()
    local terraform = require("terraform")
    
    -- Generate dependency graph
    local graph, err = terraform.graph({
        workdir = "./terraform",
        type = "plan"
    })
    
    if err then
        print("Error generating graph: " .. err)
        return false
    end
    
    print("Dependency graph: " .. graph)
    
    -- Save graph to file
    local fs = require("fs")
    local success, write_err = fs.write_file("./terraform-graph.dot", graph)
    if not success then
        print("Error writing graph file: " .. write_err)
        return false
    end
    
    print("✅ Graph saved to terraform-graph.dot")
    print("To visualize: dot -Tpng terraform-graph.dot -o terraform-graph.png")
    return true
end

-- Example 11: Console Operations
local function console_operations()
    local terraform = require("terraform")
    
    -- Evaluate expressions using Terraform console
    local expressions = {
        "var.environment",
        "local.common_tags",
        "data.aws_availability_zones.available.names",
        "length(var.subnet_cidrs)"
    }
    
    for _, expr in ipairs(expressions) do
        local result, err = terraform.console(expr, {workdir = "./terraform"})
        if err then
            print("Error evaluating '" .. expr .. "': " .. err)
        else
            print("Expression '" .. expr .. "' = " .. result)
        end
    end
    
    return true
end

-- Example 12: Complete Infrastructure Lifecycle
local function complete_infrastructure_lifecycle()
    local terraform = require("terraform")
    
    print("🚀 Starting complete Terraform infrastructure lifecycle...")
    
    -- Step 1: Initialize and validate
    print("Step 1: Initialize and validate...")
    if not basic_terraform_operations() then
        return false
    end
    
    -- Step 2: Workspace management
    print("Step 2: Managing workspaces...")
    if not manage_workspaces() then
        return false
    end
    
    -- Step 3: Provider management
    print("Step 3: Managing providers...")
    manage_providers()
    
    -- Step 4: Plan with variables
    print("Step 4: Planning with variables...")
    local plan_result, err = terraform.plan({
        workdir = "./terraform",
        vars = {
            environment = "development",
            instance_count = "1"
        },
        out = "deployment.tfplan"
    })
    
    if err then
        print("Planning failed: " .. err)
        return false
    end
    
    -- Step 5: Show plan details
    print("Step 5: Showing plan details...")
    local plan_details, err = terraform.show({
        workdir = "./terraform",
        file = "deployment.tfplan",
        json = true
    })
    
    if err then
        print("Error showing plan: " .. err)
    else
        print("Plan details: " .. string.sub(plan_details, 1, 500) .. "...")
    end
    
    -- Step 6: Apply (commented out for safety)
    print("Step 6: Apply (would be executed here)...")
    print("To apply: terraform.apply({workdir = './terraform', plan = 'deployment.tfplan'})")
    
    -- Step 7: Manage outputs
    print("Step 7: Managing outputs...")
    manage_outputs()
    
    -- Step 8: State management
    print("Step 8: Managing state...")
    manage_state()
    
    print("✅ Infrastructure lifecycle completed!")
    return true
end

-- Example 13: Disaster Recovery Operations
local function disaster_recovery()
    local terraform = require("terraform")
    
    print("🆘 Performing disaster recovery operations...")
    
    -- Force unlock state if needed
    local lock_id = os.getenv("TF_LOCK_ID")
    if lock_id then
        local unlock_success, err = terraform.force_unlock(lock_id, {
            workdir = "./terraform",
            force = true
        })
        
        if not unlock_success then
            print("Error force unlocking: " .. err)
        else
            print("✅ State unlocked successfully")
        end
    end
    
    -- Refresh state to sync with reality
    local refresh_success, err = terraform.refresh({
        workdir = "./terraform"
    })
    
    if not refresh_success then
        print("Error refreshing state: " .. err)
        return false
    end
    
    -- Get current state
    local state, err = terraform.state_pull({workdir = "./terraform"})
    if err then
        print("Error pulling state: " .. err)
        return false
    end
    
    -- Save state backup
    local fs = require("fs")
    local timestamp = os.date("%Y%m%d_%H%M%S")
    local backup_file = "./terraform-state-backup-" .. timestamp .. ".json"
    
    local backup_success, err = fs.write_file(backup_file, state)
    if not backup_success then
        print("Error saving state backup: " .. err)
        return false
    end
    
    print("✅ State backup saved to " .. backup_file)
    return true
end

-- Example 14: Destroy Infrastructure
local function destroy_infrastructure()
    local terraform = require("terraform")
    
    -- Only allow destruction in development environment
    local env = os.getenv("ENVIRONMENT") or "development"
    if env ~= "development" and env ~= "test" then
        print("⚠️  Destruction not allowed in " .. env .. " environment")
        return false
    end
    
    print("🗑️  Destroying infrastructure in " .. env .. " environment...")
    
    -- Plan destruction
    local plan_result, err = terraform.plan({
        workdir = "./terraform",
        destroy = true,
        out = "destroy.tfplan"
    })
    
    if err then
        print("Error planning destruction: " .. err)
        return false
    end
    
    print("Destruction plan: " .. plan_result)
    
    -- Destroy infrastructure
    local destroy_success, err = terraform.destroy({
        workdir = "./terraform",
        auto_approve = true
    })
    
    if not destroy_success then
        print("Error destroying infrastructure: " .. err)
        return false
    end
    
    print("✅ Infrastructure destroyed successfully")
    return true
end

-- Export task definitions
TaskDefinitions = {
    terraform_examples = {
        description = "Advanced Terraform infrastructure management examples",
        workdir = ".",
        tasks = {
            basic_ops = {
                description = "Basic Terraform operations (init, validate, fmt)",
                command = basic_terraform_operations
            },
            plan_apply = {
                description = "Plan and apply Terraform configuration",
                command = plan_and_apply
            },
            manage_vars = {
                description = "Manage Terraform variables",
                command = manage_variables
            },
            manage_workspaces = {
                description = "Create and manage Terraform workspaces",
                command = manage_workspaces
            },
            manage_state = {
                description = "Manage Terraform state",
                command = manage_state
            },
            manage_outputs = {
                description = "Manage Terraform outputs",
                command = manage_outputs
            },
            import_resources = {
                description = "Import existing resources",
                command = import_resources
            },
            manage_providers = {
                description = "Manage Terraform providers",
                command = manage_providers
            },
            taint_ops = {
                description = "Taint and untaint resources",
                command = taint_operations
            },
            generate_graph = {
                description = "Generate dependency graph",
                command = generate_graph
            },
            console_ops = {
                description = "Terraform console operations",
                command = console_operations
            },
            complete_lifecycle = {
                description = "Complete infrastructure lifecycle",
                command = complete_infrastructure_lifecycle
            },
            disaster_recovery = {
                description = "Disaster recovery operations",
                command = disaster_recovery
            },
            destroy = {
                description = "Destroy infrastructure (dev/test only)",
                command = destroy_infrastructure
            }
        }
    }
}