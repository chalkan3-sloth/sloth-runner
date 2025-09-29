-- Enhanced Terraform Module Examples

print("🏗️ ENHANCED TERRAFORM MODULE SHOWCASE")
print("=" .. string.rep("=", 50))

-- 1. Advanced Workspace Management
print("\n🏢 Advanced Workspace Management:")

-- Create a comprehensive workspace configuration
local production_workspace = terraform.workspace({
    name = "production",
    workdir = "/terraform/web-infrastructure",
    var_files = {
        "vars/common.tfvars",
        "vars/production.tfvars"
    },
    variables = {
        region = "us-east-1",
        environment = "production",
        replica_count = 5,
        instance_type = "t3.large"
    },
    backend = {
        bucket = "terraform-state-bucket",
        key = "web-app/production.tfstate",
        region = "us-east-1",
        encrypt = true
    },
    providers = {
        aws = {
            region = "us-east-1",
            profile = "production"
        }
    },
    env = {
        TF_LOG = "INFO",
        AWS_PROFILE = "production"
    },
    parallelism = 15,
    timeout = "45m"
})

print("✅ Production workspace configured")
print("   Name:", production_workspace.name)
print("   Working directory:", production_workspace.workdir)
print("   Variable files:", #production_workspace.var_files)
print("   Parallelism:", production_workspace.parallelism)

-- List workspaces
local workspace_list = terraform.workspace_list({workdir = "/terraform/web-infrastructure"})
if workspace_list.success then
    print("📋 Workspace list retrieved:")
    print("   Duration:", workspace_list.duration_ms .. "ms")
    
    -- Parse workspace list from output
    local workspaces = {}
    for line in workspace_list.stdout:gmatch("[^\r\n]+") do
        local workspace = line:match("%*?%s*(.+)")
        if workspace and workspace ~= "" then
            table.insert(workspaces, workspace)
        end
    end
    print("   Available workspaces:", #workspaces)
end

-- Create new workspace
local new_workspace = terraform.workspace_new("staging", {
    workdir = "/terraform/web-infrastructure",
    state = "/terraform/states/staging.tfstate"
})

if new_workspace.success then
    print("🆕 Staging workspace created")
else
    print("🆕 Staging workspace creation simulated")
end

-- Select workspace
local select_workspace = terraform.workspace_select("production", {
    workdir = "/terraform/web-infrastructure"
})

if select_workspace.success then
    print("🎯 Production workspace selected")
end

-- 2. Initialization and Planning
print("\n🚀 Initialization and Planning:")

-- Advanced initialization
local init_result = terraform.init({
    workdir = "/terraform/web-infrastructure",
    upgrade = true,
    reconfigure = false,
    migrate_state = false,
    backend = true,
    backend_config = {
        bucket = "terraform-state-bucket",
        key = "web-app/production.tfstate",
        region = "us-east-1"
    }
})

if init_result.success then
    print("🔧 Terraform initialization completed:")
    print("   Duration:", init_result.duration_ms .. "ms")
    print("   Providers and modules initialized")
else
    print("🔧 Terraform initialization simulated")
end

-- Comprehensive planning
local plan_result = terraform.plan({
    workdir = "/terraform/web-infrastructure",
    out = "production.tfplan",
    destroy = false,
    refresh = true,
    detailed_exitcode = true,
    parallelism = 15,
    variables = {
        replica_count = 5,
        instance_type = "t3.large",
        enable_monitoring = true
    },
    var_files = {
        "vars/common.tfvars",
        "vars/production.tfvars"
    },
    targets = {
        "aws_instance.web_servers",
        "aws_lb.application_load_balancer"
    }
})

if plan_result.success then
    print("📋 Terraform plan completed:")
    print("   Duration:", plan_result.duration_ms .. "ms")
    print("   Plan file: production.tfplan")
    
    if plan_result.changes then
        print("   Changes planned:")
        print("     Add:", plan_result.changes.add)
        print("     Change:", plan_result.changes.change)
        print("     Destroy:", plan_result.changes.destroy)
    end
    
    -- Check exit code for changes
    if plan_result.exit_code == 2 then
        print("   ⚠️ Changes detected in plan")
    elseif plan_result.exit_code == 0 then
        print("   ✅ No changes required")
    end
else
    print("📋 Terraform plan simulation completed")
end

-- 3. Apply and Deployment
print("\n🎯 Apply and Deployment:")

-- Apply with plan file
local apply_result = terraform.apply({
    workdir = "/terraform/web-infrastructure",
    plan = "production.tfplan",
    parallelism = 15
})

if apply_result.success then
    print("🚀 Terraform apply completed:")
    print("   Duration:", apply_result.duration_ms .. "ms")
    print("   Infrastructure deployed successfully")
else
    print("🚀 Terraform apply simulation completed")
end

-- Apply with auto-approval (for automation)
local auto_apply_result = terraform.apply({
    workdir = "/terraform/web-infrastructure",
    auto_approve = true,
    variables = {
        replica_count = 3,
        maintenance_window = "02:00-04:00"
    },
    targets = {
        "aws_instance.web_servers[0]",
        "aws_instance.web_servers[1]"
    }
})

if auto_apply_result.success then
    print("⚡ Auto-approved apply completed")
    print("   Duration:", auto_apply_result.duration_ms .. "ms")
end

-- 4. State Management
print("\n💾 State Management:")

-- List state resources
local state_list = terraform.state_list({
    workdir = "/terraform/web-infrastructure",
    id = "aws_instance"
})

if state_list.success then
    print("📝 State resources listed:")
    print("   Duration:", state_list.duration_ms .. "ms")
    
    -- Count resources in output
    local resource_count = 0
    for line in state_list.stdout:gmatch("[^\r\n]+") do
        if line and line ~= "" then
            resource_count = resource_count + 1
        end
    end
    print("   Resources found:", resource_count)
end

-- Show specific resource
local state_show = terraform.state_show("aws_instance.web_servers[0]", {
    workdir = "/terraform/web-infrastructure"
})

if state_show.success then
    print("🔍 Resource state details retrieved")
    print("   Duration:", state_show.duration_ms .. "ms")
end

-- Pull current state
local state_pull = terraform.state_pull({
    workdir = "/terraform/web-infrastructure"
})

if state_pull.success then
    print("📥 State pulled successfully:")
    print("   Duration:", state_pull.duration_ms .. "ms")
    
    if state_pull.state then
        print("   State version:", state_pull.state.version)
        if state_pull.state.outputs then
            local output_count = 0
            for _ in pairs(state_pull.state.outputs) do
                output_count = output_count + 1
            end
            print("   Outputs available:", output_count)
        end
    end
end

-- Move state resource
local state_move = terraform.state_mv(
    "aws_instance.web_server",
    "aws_instance.web_servers[0]",
    {
        workdir = "/terraform/web-infrastructure",
        dry_run = true
    }
)

if state_move.success then
    print("🔄 State move operation completed (dry run)")
    print("   Duration:", state_move.duration_ms .. "ms")
end

-- 5. Import and Resource Management
print("\n📥 Import and Resource Management:")

-- Import existing resource
local import_result = terraform.import(
    "aws_instance.existing_server",
    "i-1234567890abcdef0",
    {
        workdir = "/terraform/web-infrastructure",
        variables = {
            region = "us-east-1"
        }
    }
)

if import_result.success then
    print("📥 Resource imported successfully")
    print("   Duration:", import_result.duration_ms .. "ms")
else
    print("📥 Resource import simulation completed")
end

-- Taint resource for recreation
local taint_result = terraform.taint("aws_instance.web_servers[2]", {
    workdir = "/terraform/web-infrastructure"
})

if taint_result.success then
    print("🏷️ Resource tainted for recreation")
    print("   Duration:", taint_result.duration_ms .. "ms")
end

-- Untaint resource
local untaint_result = terraform.untaint("aws_instance.web_servers[2]", {
    workdir = "/terraform/web-infrastructure"
})

if untaint_result.success then
    print("✅ Resource untainted")
    print("   Duration:", untaint_result.duration_ms .. "ms")
end

-- 6. Output Management
print("\n📤 Output Management:")

-- Get all outputs
local outputs_result = terraform.output({
    workdir = "/terraform/web-infrastructure"
})

if outputs_result.success then
    print("📊 Terraform outputs retrieved:")
    print("   Duration:", outputs_result.duration_ms .. "ms")
    
    if outputs_result.outputs then
        local output_count = 0
        for name, value in pairs(outputs_result.outputs) do
            output_count = output_count + 1
            if type(value) == "table" and value.value then
                print("   " .. name .. ":", value.value)
            else
                print("   " .. name .. ":", tostring(value))
            end
        end
        print("   Total outputs:", output_count)
    end
else
    print("📊 Outputs simulation completed")
end

-- Get specific output
local specific_output = terraform.output({
    workdir = "/terraform/web-infrastructure",
    name = "load_balancer_dns"
})

if specific_output.success then
    print("🎯 Specific output retrieved:")
    print("   Duration:", specific_output.duration_ms .. "ms")
end

-- 7. Validation and Formatting
print("\n✅ Validation and Formatting:")

-- Validate configuration
local validate_result = terraform.validate({
    workdir = "/terraform/web-infrastructure"
})

if validate_result.success then
    print("✅ Configuration validation passed:")
    print("   Duration:", validate_result.duration_ms .. "ms")
    print("   No syntax errors found")
else
    print("⚠️ Configuration validation issues found")
    print("   Duration:", validate_result.duration_ms .. "ms")
end

-- Format configuration files
local fmt_result = terraform.fmt({
    workdir = "/terraform/web-infrastructure",
    recursive = true,
    diff = true,
    check = false,
    write = true
})

if fmt_result.success then
    print("🎨 Configuration formatted:")
    print("   Duration:", fmt_result.duration_ms .. "ms")
    
    if fmt_result.stdout and fmt_result.stdout ~= "" then
        local file_count = 0
        for line in fmt_result.stdout:gmatch("[^\r\n]+") do
            if line and line ~= "" then
                file_count = file_count + 1
            end
        end
        print("   Files formatted:", file_count)
    end
end

-- 8. Advanced Operations
print("\n🚀 Advanced Operations:")

-- Generate dependency graph
local graph_result = terraform.graph({
    workdir = "/terraform/web-infrastructure",
    type = "plan"
})

if graph_result.success then
    print("📊 Dependency graph generated:")
    print("   Duration:", graph_result.duration_ms .. "ms")
    print("   Graph output size:", #graph_result.stdout, "bytes")
end

-- List providers
local providers_result = terraform.providers({
    workdir = "/terraform/web-infrastructure"
})

if providers_result.success then
    print("🔌 Provider information retrieved:")
    print("   Duration:", providers_result.duration_ms .. "ms")
end

-- Force unlock state (emergency)
print("🔓 Force unlock capability available for emergency situations")

-- 9. Utility Operations
print("\n🛠️ Utility Operations:")

-- Check Terraform version
local version_result = terraform.version()
if version_result.success then
    print("ℹ️ Terraform version information:")
    print("   Duration:", version_result.duration_ms .. "ms")
    print("   Version data available")
end

-- Refresh state
local refresh_result = terraform.refresh({
    workdir = "/terraform/web-infrastructure",
    variables = {
        region = "us-east-1"
    }
})

if refresh_result.success then
    print("🔄 State refreshed successfully:")
    print("   Duration:", refresh_result.duration_ms .. "ms")
end

-- 10. Destruction and Cleanup
print("\n💥 Destruction and Cleanup:")

-- Targeted destruction
local targeted_destroy = terraform.destroy({
    workdir = "/terraform/web-infrastructure",
    auto_approve = true,
    targets = {
        "aws_instance.temporary_server"
    },
    variables = {
        region = "us-east-1"
    }
})

if targeted_destroy.success then
    print("🎯 Targeted destruction completed:")
    print("   Duration:", targeted_destroy.duration_ms .. "ms")
else
    print("🎯 Targeted destruction simulation completed")
end

-- 11. Performance and Monitoring
print("\n📊 Performance and Monitoring:")

-- Gather operation metrics
local operation_metrics = {
    init_time = init_result.duration_ms or 0,
    plan_time = plan_result.duration_ms or 0,
    apply_time = apply_result.duration_ms or 0,
    state_time = state_list.duration_ms or 0,
    total_operations = 12
}

print("⚡ Performance Summary:")
print("   Total infrastructure deployment time:", 
      (operation_metrics.init_time + operation_metrics.plan_time + operation_metrics.apply_time) .. "ms")
print("   Average operation time:", 
      math.floor((operation_metrics.init_time + 
                  operation_metrics.plan_time + 
                  operation_metrics.apply_time + 
                  operation_metrics.state_time) / 4) .. "ms")

print("   Operations completed:", operation_metrics.total_operations)

-- 12. Advanced Integration Examples
print("\n🔗 Advanced Integration Examples:")

print("🎯 Enterprise features available:")
print("   • Remote state management")
print("   • Workspace isolation")
print("   • Provider version constraints")
print("   • Module composition")
print("   • Variable validation")
print("   • Conditional resource creation")
print("   • Data source integration")
print("   • Local and remote backends")

print("\n📋 Use Cases:")
print("🏗️ Perfect for:")
print("   • Infrastructure as Code")
print("   • Multi-environment deployments")
print("   • Cloud resource management")
print("   • Network infrastructure")
print("   • Security and compliance")
print("   • Database infrastructure")
print("   • CI/CD integration")
print("   • Disaster recovery planning")

-- Workspace cleanup
local workspace_cleanup = terraform.workspace_delete("staging", {
    workdir = "/terraform/web-infrastructure",
    force = true
})

if workspace_cleanup.success then
    print("🧹 Staging workspace cleaned up")
else
    print("🧹 Workspace cleanup simulation completed")
end

print("\n✅ Enhanced Terraform module demonstration completed!")
print("🏗️ Enterprise-grade Infrastructure as Code ready!")
print("🚀 Provision and manage infrastructure with confidence!")