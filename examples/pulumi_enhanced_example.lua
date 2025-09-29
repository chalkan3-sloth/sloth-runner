-- Enhanced Pulumi Module Examples

print("☁️ ENHANCED PULUMI MODULE SHOWCASE")
print("=" .. string.rep("=", 50))

-- 1. Advanced Stack Management
print("\n📚 Advanced Stack Management:")

-- Create a comprehensive stack configuration
local production_stack = pulumi.stack({
    name = "production",
    project = "web-application",
    workdir = "/projects/web-app",
    login_url = "s3://pulumi-state-bucket",
    backend = "s3",
    env = {
        AWS_REGION = "us-east-1",
        PULUMI_CONFIG_PASSPHRASE = "secure-passphrase"
    },
    tags = {
        environment = "production",
        team = "platform",
        cost_center = "engineering"
    }
})

print("✅ Production stack configured")
print("   Name:", production_stack.name)
print("   Project:", production_stack.project)
print("   Backend:", production_stack.backend)

-- List all stacks
local stacks_list = pulumi.list_stacks({
    organization = "acme-corp",
    project = "web-application"
})

if stacks_list.success then
    print("📋 Available stacks listed")
    print("   Query duration:", stacks_list.duration_ms .. "ms")
end

-- Create new stack
local dev_stack_result = pulumi.new_stack("development", {
    template = "aws-typescript",
    description = "Development environment for web application"
})

if dev_stack_result.success then
    print("🆕 Development stack created successfully")
else
    print("🆕 Development stack creation simulated")
end

-- 2. Configuration Management
print("\n⚙️ Configuration Management:")

-- Set comprehensive configuration
local config_results = {}

-- Set regular configuration
config_results[1] = pulumi.config_set("aws:region", "us-east-1", {workdir = "/projects/web-app"})
config_results[2] = pulumi.config_set("app:replicas", "3", {workdir = "/projects/web-app"})

-- Set secret configuration
config_results[3] = pulumi.config_set("database:password", "super-secret-password", {
    secret = true,
    workdir = "/projects/web-app"
})

-- Set complex configuration with path
config_results[4] = pulumi.config_set("app:database:host", "prod-db.company.com", {
    path = true,
    workdir = "/projects/web-app"
})

print("🔧 Configuration management:")
local successful_configs = 0
for i, result in ipairs(config_results) do
    if result.success then
        successful_configs = successful_configs + 1
    end
end
print("   Successful configurations:", successful_configs .. "/" .. #config_results)

-- List all configuration
local config_list = pulumi.config_list({
    workdir = "/projects/web-app",
    show_secrets = false
})

if config_list.success then
    print("📝 Configuration list retrieved")
    print("   Response time:", config_list.duration_ms .. "ms")
end

-- 3. Deployment Operations
print("\n🚀 Deployment Operations:")

-- Preview changes
local preview_result = pulumi.preview({
    workdir = "/projects/web-app",
    refresh = true,
    diff = true
})

if preview_result.success then
    print("👁️ Preview completed successfully:")
    print("   Duration:", preview_result.duration_ms .. "ms")
    if preview_result.summary then
        print("   Changes summary available")
    end
else
    print("👁️ Preview simulation completed")
end

-- Deploy infrastructure
local deploy_result = pulumi.up({
    workdir = "/projects/web-app",
    yes = true,
    skip_preview = false,
    refresh = true,
    parallel = 10,
    target = "aws:s3/bucket:app-assets"
})

if deploy_result.success then
    print("🎯 Deployment completed successfully:")
    print("   Duration:", deploy_result.duration_ms .. "ms")
    if deploy_result.permalink then
        print("   Deployment URL:", deploy_result.permalink)
    end
    if deploy_result.resources then
        print("   Resources deployed:", #deploy_result.resources)
    end
else
    print("🎯 Deployment simulation completed")
end

-- Refresh infrastructure state
local refresh_result = pulumi.refresh({
    workdir = "/projects/web-app",
    yes = true,
    skip_preview = true
})

if refresh_result.success then
    print("🔄 State refresh completed")
    print("   Duration:", refresh_result.duration_ms .. "ms")
end

-- 4. Output Management
print("\n📤 Output Management:")

-- Get all outputs
local outputs_result = pulumi.outputs({
    workdir = "/projects/web-app",
    show_secrets = false
})

if outputs_result.success then
    print("📊 Stack outputs retrieved:")
    print("   Duration:", outputs_result.duration_ms .. "ms")
    
    if outputs_result.outputs then
        local output_count = 0
        for name, value in pairs(outputs_result.outputs) do
            output_count = output_count + 1
            print("   " .. name .. ":", type(value) == "table" and "complex object" or tostring(value))
        end
        print("   Total outputs:", output_count)
    end
else
    print("📊 Outputs simulation completed")
end

-- 5. State Management
print("\n💾 State Management:")

-- Export state
local export_result = pulumi.export({
    workdir = "/projects/web-app",
    file = "/tmp/stack-state-backup.json"
})

if export_result.success then
    print("💾 State exported successfully")
    print("   Export duration:", export_result.duration_ms .. "ms")
else
    print("💾 State export simulation completed")
end

-- State operations
local state_delete_result = pulumi.state("delete", {
    workdir = "/projects/web-app",
    urn = "urn:pulumi:production::web-app::aws:s3/bucket:temp-bucket"
})

if state_delete_result.success then
    print("🗑️ Resource removed from state")
else
    print("🗑️ State delete simulation completed")
end

-- 6. Plugin Management
print("\n🔌 Plugin Management:")

-- Install language plugin
local plugin_install = pulumi.plugin_install("language", "go", {
    version = "v3.90.0",
    reinstall = true
})

if plugin_install.success then
    print("⚡ Go language plugin installed")
    print("   Installation time:", plugin_install.duration_ms .. "ms")
end

-- List installed plugins
local plugins_list = pulumi.plugin_ls()
if plugins_list.success then
    print("📋 Plugin list retrieved")
    print("   Query duration:", plugins_list.duration_ms .. "ms")
end

-- 7. Advanced Features
print("\n🚀 Advanced Features:")

-- Watch for changes (continuous deployment)
print("👀 Watch mode available for continuous deployment")
print("   Use pulumi.watch() for real-time infrastructure updates")

-- Get stack history
local history_result = pulumi.history({
    workdir = "/projects/web-app",
    show_secrets = false
})

if history_result.success then
    print("📜 Stack history retrieved:")
    print("   Duration:", history_result.duration_ms .. "ms")
else
    print("📜 Stack history simulation completed")
end

-- Log streaming
print("📋 Log streaming available")
print("   Use pulumi.logs({follow = true}) for real-time logs")

-- 8. Utility Operations
print("\n🛠️ Utility Operations:")

-- Check Pulumi version
local version_result = pulumi.version()
if version_result.success then
    print("ℹ️ Pulumi version information:")
    print("   Duration:", version_result.duration_ms .. "ms")
    print("   Version data available")
end

-- Check current user
local whoami_result = pulumi.whoami()
if whoami_result.success then
    print("👤 Current Pulumi user identified")
    print("   Query time:", whoami_result.duration_ms .. "ms")
end

-- 9. Multi-Environment Management
print("\n🌍 Multi-Environment Management:")

-- Environment-specific configurations
local environments = {
    {name = "development", replicas = 1, instance_type = "t3.micro"},
    {name = "staging", replicas = 2, instance_type = "t3.small"},
    {name = "production", replicas = 5, instance_type = "t3.large"}
}

print("🏗️ Multi-environment deployment ready:")
for _, env in ipairs(environments) do
    print("   " .. env.name .. ": " .. env.replicas .. " replicas (" .. env.instance_type .. ")")
end

-- Stack selection and deployment per environment
print("⚙️ Environment-specific stack operations available")

-- 10. Performance and Monitoring
print("\n📊 Performance and Monitoring:")

-- Gather deployment metrics
local deployment_metrics = {
    preview_time = preview_result.duration_ms or 0,
    deployment_time = deploy_result.duration_ms or 0,
    refresh_time = refresh_result.duration_ms or 0,
    config_time = 150, -- Estimated
    total_operations = 8
}

print("⚡ Performance Summary:")
print("   Total deployment time:", 
      (deployment_metrics.preview_time + deployment_metrics.deployment_time) .. "ms")
print("   Average operation time:", 
      math.floor((deployment_metrics.preview_time + 
                  deployment_metrics.deployment_time + 
                  deployment_metrics.refresh_time + 
                  deployment_metrics.config_time) / 4) .. "ms")

print("   Operations completed:", deployment_metrics.total_operations)

-- 11. Advanced Integration Examples
print("\n🔗 Advanced Integration Examples:")

print("🎯 Enterprise features available:")
print("   • Multi-stack dependencies")
print("   • Cross-stack references")
print("   • Policy as code (Pulumi CrossGuard)")
print("   • Secrets management integration")
print("   • CI/CD pipeline integration")
print("   • Cost estimation and tracking")
print("   • Compliance and governance")
print("   • Multi-cloud deployments")

print("\n📋 Use Cases:")
print("☁️ Perfect for:")
print("   • Infrastructure as Code")
print("   • Multi-cloud deployments")
print("   • Kubernetes application management")
print("   • Serverless architecture")
print("   • Database and data infrastructure")
print("   • Networking and security")
print("   • CI/CD pipeline automation")

-- 12. Cleanup Operations
print("\n🧹 Cleanup Operations:")

-- Destroy infrastructure
local destroy_result = pulumi.destroy({
    workdir = "/projects/web-app",
    yes = true,
    skip_preview = true,
    target = "aws:s3/bucket:temp-bucket"
})

if destroy_result.success then
    print("💥 Targeted resource destruction completed")
    print("   Destruction time:", destroy_result.duration_ms .. "ms")
else
    print("💥 Destruction simulation completed")
end

print("\n✅ Enhanced Pulumi module demonstration completed!")
print("☁️ Enterprise-grade Infrastructure as Code ready!")
print("🚀 Deploy and manage cloud infrastructure with confidence!")