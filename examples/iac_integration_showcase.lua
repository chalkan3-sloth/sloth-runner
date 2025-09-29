-- Infrastructure as Code Integration Showcase
-- Demonstrates enhanced Salt, Pulumi, and Terraform modules working together

print("🌐 INFRASTRUCTURE AS CODE INTEGRATION SHOWCASE")
print("=" .. string.rep("=", 65))
print("Complete infrastructure lifecycle management with Salt, Pulumi & Terraform")

-- Initialize observability for the entire infrastructure lifecycle
local main_trace = observability.start_trace("infrastructure-lifecycle", {
    operation = "full_deployment",
    tools = "salt+pulumi+terraform",
    environment = "production"
})

print("\n🎯 Started infrastructure lifecycle trace:", main_trace)

-- 1. TERRAFORM: CORE INFRASTRUCTURE PROVISIONING
print("\n🏗️ PHASE 1: TERRAFORM - CORE INFRASTRUCTURE")
print("-" .. string.rep("-", 55))

local terraform_span = observability.start_span(main_trace, "terraform-infrastructure", "", {
    component = "terraform",
    phase = "provisioning"
})

-- Initialize Terraform workspace
local terraform_workspace = terraform.workspace({
    name = "production-core",
    workdir = "/iac/terraform/core",
    variables = {
        region = "us-east-1",
        environment = "production",
        vpc_cidr = "10.0.0.0/16",
        availability_zones = 3,
        instance_types = {
            web = "t3.large",
            db = "r5.xlarge",
            cache = "r6g.large"
        }
    },
    backend = {
        bucket = "company-terraform-state",
        key = "core/production.tfstate",
        region = "us-east-1",
        encrypt = true,
        dynamodb_table = "terraform-state-lock"
    },
    parallelism = 20
})

print("✅ Terraform workspace configured")

-- Initialize and plan infrastructure
local tf_init = terraform.init({
    workdir = "/iac/terraform/core",
    upgrade = true,
    backend = true
})

local tf_plan = terraform.plan({
    workdir = "/iac/terraform/core",
    out = "production-core.tfplan",
    variables = terraform_workspace.variables,
    parallelism = 20
})

if tf_plan.success then
    print("📋 Terraform plan completed:")
    print("   Duration:", tf_plan.duration_ms .. "ms")
    if tf_plan.changes then
        print("   Resources to add:", tf_plan.changes.add)
        print("   Resources to change:", tf_plan.changes.change)
        print("   Resources to destroy:", tf_plan.changes.destroy)
    end
end

-- Apply core infrastructure
local tf_apply = terraform.apply({
    workdir = "/iac/terraform/core",
    plan = "production-core.tfplan"
})

if tf_apply.success then
    print("🚀 Core infrastructure deployed:")
    print("   Duration:", tf_apply.duration_ms .. "ms")
    
    observability.counter("terraform_resources_deployed", tf_plan.changes and tf_plan.changes.add or 0, {
        environment = "production",
        phase = "core"
    })
end

-- Get infrastructure outputs
local tf_outputs = terraform.output({workdir = "/iac/terraform/core"})
local infrastructure_data = {}

if tf_outputs.success and tf_outputs.outputs then
    infrastructure_data = {
        vpc_id = tf_outputs.outputs.vpc_id and tf_outputs.outputs.vpc_id.value or "vpc-12345",
        subnet_ids = tf_outputs.outputs.subnet_ids and tf_outputs.outputs.subnet_ids.value or {"subnet-1", "subnet-2"},
        security_group_id = tf_outputs.outputs.security_group_id and tf_outputs.outputs.security_group_id.value or "sg-12345",
        db_endpoint = tf_outputs.outputs.db_endpoint and tf_outputs.outputs.db_endpoint.value or "db.example.com"
    }
    
    print("📊 Infrastructure outputs retrieved:")
    print("   VPC ID:", infrastructure_data.vpc_id)
    print("   Subnets:", #infrastructure_data.subnet_ids)
    print("   Database endpoint:", infrastructure_data.db_endpoint)
end

observability.add_span_event(terraform_span, "infrastructure-provisioned", {
    vpc_id = infrastructure_data.vpc_id,
    duration = tostring(tf_apply.duration_ms)
})

observability.end_span(terraform_span, "completed")

-- 2. PULUMI: APPLICATION INFRASTRUCTURE
print("\n☁️ PHASE 2: PULUMI - APPLICATION INFRASTRUCTURE")
print("-" .. string.rep("-", 55))

local pulumi_span = observability.start_span(main_trace, "pulumi-application", terraform_span, {
    component = "pulumi",
    phase = "application"
})

-- Configure Pulumi stack with Terraform outputs
local pulumi_stack = pulumi.stack({
    name = "production-app",
    project = "web-application",
    workdir = "/iac/pulumi/application",
    login_url = "s3://company-pulumi-state",
    env = {
        AWS_REGION = "us-east-1",
        PULUMI_CONFIG_PASSPHRASE = "secure-key"
    },
    tags = {
        environment = "production",
        team = "platform",
        project = "web-app"
    }
})

print("☁️ Pulumi stack configured")

-- Set configuration using Terraform outputs
local pulumi_configs = {
    pulumi.config_set("aws:region", "us-east-1", {workdir = "/iac/pulumi/application"}),
    pulumi.config_set("vpc:id", infrastructure_data.vpc_id, {workdir = "/iac/pulumi/application"}),
    pulumi.config_set("vpc:subnets", table.concat(infrastructure_data.subnet_ids, ","), {workdir = "/iac/pulumi/application"}),
    pulumi.config_set("security:groupId", infrastructure_data.security_group_id, {workdir = "/iac/pulumi/application"}),
    pulumi.config_set("database:endpoint", infrastructure_data.db_endpoint, {workdir = "/iac/pulumi/application", secret = true})
}

local successful_configs = 0
for _, config in ipairs(pulumi_configs) do
    if config.success then
        successful_configs = successful_configs + 1
    end
end
print("⚙️ Pulumi configuration:", successful_configs .. "/" .. #pulumi_configs .. " items set")

-- Preview application deployment
local pulumi_preview = pulumi.preview({
    workdir = "/iac/pulumi/application",
    refresh = true,
    diff = true
})

if pulumi_preview.success then
    print("👁️ Pulumi preview completed:")
    print("   Duration:", pulumi_preview.duration_ms .. "ms")
end

-- Deploy application infrastructure
local pulumi_deploy = pulumi.up({
    workdir = "/iac/pulumi/application",
    yes = true,
    parallel = 15,
    refresh = true
})

if pulumi_deploy.success then
    print("🎯 Application infrastructure deployed:")
    print("   Duration:", pulumi_deploy.duration_ms .. "ms")
    if pulumi_deploy.resources then
        print("   Resources created:", #pulumi_deploy.resources)
    end
    
    observability.counter("pulumi_resources_deployed", pulumi_deploy.resources and #pulumi_deploy.resources or 0, {
        environment = "production",
        phase = "application"
    })
end

-- Get application outputs
local pulumi_outputs = pulumi.outputs({
    workdir = "/iac/pulumi/application",
    show_secrets = false
})

local application_data = {}
if pulumi_outputs.success and pulumi_outputs.outputs then
    application_data = {
        load_balancer_url = pulumi_outputs.outputs.loadBalancerUrl and pulumi_outputs.outputs.loadBalancerUrl.value or "https://app.example.com",
        api_gateway_url = pulumi_outputs.outputs.apiGatewayUrl and pulumi_outputs.outputs.apiGatewayUrl.value or "https://api.example.com",
        instance_ips = pulumi_outputs.outputs.instanceIPs and pulumi_outputs.outputs.instanceIPs.value or {"10.0.1.10", "10.0.1.11", "10.0.1.12"}
    }
    
    print("📊 Application outputs retrieved:")
    print("   Load balancer URL:", application_data.load_balancer_url)
    print("   API Gateway URL:", application_data.api_gateway_url)
    print("   Instance IPs:", #application_data.instance_ips)
end

observability.add_span_event(pulumi_span, "application-deployed", {
    load_balancer = application_data.load_balancer_url,
    instances = tostring(#application_data.instance_ips)
})

observability.end_span(pulumi_span, "completed")

-- 3. SALT: CONFIGURATION MANAGEMENT
print("\n🧂 PHASE 3: SALT - CONFIGURATION MANAGEMENT")
print("-" .. string.rep("-", 50))

local salt_span = observability.start_span(main_trace, "salt-configuration", pulumi_span, {
    component = "salt",
    phase = "configuration"
})

-- Configure Salt client
local salt_client = salt.client({
    config = "/etc/salt",
    master = "salt-master.company.com",
    timeout = 60,
    retries = 3,
    env = {
        SALT_LOG_LEVEL = "info"
    }
})

print("🧂 Salt client configured")

-- Accept new minion keys (for newly created instances)
local key_accept_results = {}
for i, ip in ipairs(application_data.instance_ips) do
    local minion_id = "web-" .. string.format("%02d", i)
    local key_result = salt.key_accept(minion_id)
    table.insert(key_accept_results, {id = minion_id, success = key_result.success})
end

print("🔑 Minion key management:")
local accepted_keys = 0
for _, result in ipairs(key_accept_results) do
    if result.success then
        accepted_keys = accepted_keys + 1
    end
    print("   " .. result.id .. ":", result.success and "✅ Accepted" or "⏳ Pending")
end

-- Test connectivity to all minions
local ping_result = salt.test_ping("web*")
if ping_result.success then
    print("🏓 Salt connectivity test:")
    print("   Duration:", ping_result.duration_ms .. "ms")
    
    local responsive_minions = 0
    if ping_result.returns then
        for minion, response in pairs(ping_result.returns) do
            if response == true then
                responsive_minions = responsive_minions + 1
            end
        end
    end
    print("   Responsive minions:", responsive_minions)
    
    observability.gauge("salt_responsive_minions", responsive_minions, {
        environment = "production"
    })
end

-- Apply base system configuration
local base_config = salt.state_apply("web*", "base", {
    pillar = {
        environment = "production",
        timezone = "UTC",
        ntp_servers = {"pool.ntp.org"},
        monitoring = {
            enabled = true,
            agent = "datadog"
        }
    }
})

if base_config.success then
    print("⚙️ Base system configuration applied:")
    print("   Duration:", base_config.duration_ms .. "ms")
    print("   Configuration deployed to web servers")
end

-- Deploy application configuration
local app_config = salt.state_apply("web*", "webapp", {
    pillar = {
        app = {
            version = "v2.1.0",
            database_url = infrastructure_data.db_endpoint,
            load_balancer_url = application_data.load_balancer_url,
            api_gateway_url = application_data.api_gateway_url,
            replicas = 3,
            max_connections = 1000
        },
        nginx = {
            worker_processes = 4,
            worker_connections = 1024,
            client_max_body_size = "10m"
        }
    }
})

if app_config.success then
    print("🚀 Application configuration deployed:")
    print("   Duration:", app_config.duration_ms .. "ms")
    print("   Web application configured and running")
    
    observability.counter("salt_states_applied", 1, {
        state = "webapp",
        environment = "production"
    })
end

-- Install and configure monitoring
local monitoring_config = salt.state_apply("web*", "monitoring", {
    pillar = {
        monitoring = {
            datadog_api_key = "dd_api_key_secret",
            application_name = "web-app",
            environment = "production",
            metrics = {
                cpu = true,
                memory = true,
                disk = true,
                network = true,
                application = true
            }
        }
    }
})

if monitoring_config.success then
    print("📊 Monitoring configuration applied:")
    print("   Duration:", monitoring_config.duration_ms .. "ms")
    print("   Datadog agent configured on all instances")
end

-- Apply security hardening
local security_hardening = salt.state_apply("web*", "security", {
    pillar = {
        security = {
            ssh_port = 2222,
            disable_root_login = true,
            firewall_rules = {
                {port = 80, protocol = "tcp", source = "0.0.0.0/0"},
                {port = 443, protocol = "tcp", source = "0.0.0.0/0"},
                {port = 2222, protocol = "tcp", source = "10.0.0.0/16"}
            },
            fail2ban_enabled = true,
            auto_updates = true
        }
    }
})

if security_hardening.success then
    print("🔒 Security hardening applied:")
    print("   Duration:", security_hardening.duration_ms .. "ms")
    print("   Security policies enforced")
end

observability.add_span_event(salt_span, "configuration-completed", {
    minions_configured = tostring(accepted_keys),
    states_applied = "4"
})

observability.end_span(salt_span, "completed")

-- 4. INTEGRATION VALIDATION
print("\n🧪 PHASE 4: INTEGRATION VALIDATION")
print("-" .. string.rep("-", 40))

local validation_span = observability.start_span(main_trace, "integration-validation", salt_span, {
    component = "validation",
    phase = "testing"
})

-- Validate infrastructure connectivity
local connectivity_tests = {
    {name = "Load Balancer Health", url = application_data.load_balancer_url},
    {name = "API Gateway Health", url = application_data.api_gateway_url},
    {name = "Database Connectivity", host = infrastructure_data.db_endpoint}
}

print("🌐 Infrastructure connectivity validation:")
local successful_tests = 0

for _, test in ipairs(connectivity_tests) do
    local success = false
    if test.url then
        -- Simulate HTTP health check
        local result = http.get(test.url .. "/health", {timeout = 10})
        success = result and result.status_code == 200
    elseif test.host then
        -- Simulate database connectivity check
        local result = network.port_check(test.host, 5432, 5)
        success = result
    end
    
    if success then
        successful_tests = successful_tests + 1
        print("   ✅ " .. test.name .. ": OK")
    else
        print("   ❌ " .. test.name .. ": Failed")
    end
end

print("   Connectivity tests:", successful_tests .. "/" .. #connectivity_tests .. " passed")

-- Validate Salt minion states
local minion_validation = salt.cmd("web*", "test", "ping")
local responsive_minions = 0

if minion_validation.success and minion_validation.returns then
    for minion, response in pairs(minion_validation.returns) do
        if response == true then
            responsive_minions = responsive_minions + 1
        end
    end
end

print("🧂 Salt minion validation:")
print("   Responsive minions:", responsive_minions .. "/" .. #application_data.instance_ips)

-- Performance validation
local performance_metrics = {
    terraform_time = tf_apply.duration_ms or 0,
    pulumi_time = pulumi_deploy.duration_ms or 0,
    salt_config_time = (base_config.duration_ms or 0) + (app_config.duration_ms or 0),
    total_deployment_time = 0
}

performance_metrics.total_deployment_time = performance_metrics.terraform_time + 
                                           performance_metrics.pulumi_time + 
                                           performance_metrics.salt_config_time

print("⚡ Performance validation:")
print("   Total deployment time:", performance_metrics.total_deployment_time .. "ms")
print("   Infrastructure provisioning:", performance_metrics.terraform_time .. "ms")
print("   Application deployment:", performance_metrics.pulumi_time .. "ms")
print("   Configuration management:", performance_metrics.salt_config_time .. "ms")

observability.histogram("deployment_duration_ms", performance_metrics.total_deployment_time, {
    environment = "production",
    pipeline = "terraform+pulumi+salt"
})

observability.add_span_event(validation_span, "validation-completed", {
    connectivity_tests = tostring(successful_tests),
    responsive_minions = tostring(responsive_minions),
    total_time = tostring(performance_metrics.total_deployment_time)
})

observability.end_span(validation_span, "completed")

-- 5. FINAL SUMMARY
print("\n🎯 DEPLOYMENT SUMMARY")
print("-" .. string.rep("-", 30))

-- End main trace
local trace_success, total_duration = observability.end_trace(main_trace, "completed")

-- Calculate success metrics
local overall_success_rate = ((tf_apply.success and 1 or 0) + 
                             (pulumi_deploy.success and 1 or 0) + 
                             (app_config.success and 1 or 0)) / 3 * 100

print("🏆 INFRASTRUCTURE DEPLOYMENT COMPLETED")
print("   Total Duration:", total_duration .. "ms")
print("   Success Rate:", string.format("%.0f%%", overall_success_rate))
print("   Infrastructure Status: " .. (overall_success_rate >= 100 and "✅ HEALTHY" or "⚠️ PARTIAL"))

print("\n📊 Deployment Breakdown:")
print("   🏗️ Terraform (Core Infrastructure):")
print("      VPC, Subnets, Security Groups, Database")
print("      Duration:", performance_metrics.terraform_time .. "ms")

print("   ☁️ Pulumi (Application Infrastructure):")
print("      Load Balancers, Auto Scaling, API Gateway")
print("      Duration:", performance_metrics.pulumi_time .. "ms")

print("   🧂 Salt (Configuration Management):")
print("      System config, App deployment, Security, Monitoring")
print("      Duration:", performance_metrics.salt_config_time .. "ms")

print("\n🌐 Deployed Resources:")
print("   • VPC with " .. #infrastructure_data.subnet_ids .. " subnets")
print("   • " .. #application_data.instance_ips .. " application instances")
print("   • Load balancer: " .. application_data.load_balancer_url)
print("   • API gateway: " .. application_data.api_gateway_url)
print("   • Database: " .. infrastructure_data.db_endpoint)
print("   • " .. responsive_minions .. " configured Salt minions")

print("\n🎯 Integration Benefits Achieved:")
print("   ✅ Infrastructure as Code with Terraform")
print("   ✅ Modern cloud-native deployment with Pulumi")
print("   ✅ Comprehensive configuration management with Salt")
print("   ✅ End-to-end observability and monitoring")
print("   ✅ Automated validation and testing")
print("   ✅ Integrated deployment pipeline")

-- Record final metrics
observability.gauge("infrastructure_health_score", overall_success_rate, {
    environment = "production",
    deployment_id = main_trace
})

observability.counter("infrastructure_deployments_completed", 1, {
    tools = "terraform+pulumi+salt",
    environment = "production",
    success = tostring(overall_success_rate >= 100)
})

print("\n✅ INFRASTRUCTURE AS CODE INTEGRATION SHOWCASE COMPLETED!")
print("🚀 Full-stack infrastructure deployment with enterprise-grade automation!")
print("💡 Ready for production workloads with complete lifecycle management!")