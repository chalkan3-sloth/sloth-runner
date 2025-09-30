-- MODERN DSL ONLY - Azure Integration Example
-- Demonstrates Azure operations using Modern DSL

-- Azure Resource Groups task
local azure_list_rg = task("azure_list_resource_groups")
    :description("List Azure resource groups")
    :command(function(params)
        log.info("🏢 Listing Azure resource groups...")
        
        local result = exec.run("az group list --output table", {
            timeout = "45s",
            capture_output = true
        })
        
        if result.success then
            log.info("✅ Successfully listed resource groups")
            return true, result.output, {
                resource_groups = result.output,
                command_type = "resource_group_list"
            }
        else
            return false, "Failed to list resource groups: " .. (result.error or "unknown error")
        end
    end)
    :timeout("90s")
    :retries(2, "exponential")
    :on_success(function(params, output)
        log.info("📋 Resource groups listed successfully")
    end)
    :build()

-- Azure Storage Accounts task
local azure_list_storage = task("azure_list_storage")
    :description("List Azure storage accounts")
    :command(function(params)
        log.info("💾 Listing Azure storage accounts...")
        
        local result = exec.run("az storage account list --output table", {
            timeout = "60s",
            capture_output = true
        })
        
        if result.success then
            return true, result.output, {
                storage_accounts = result.output,
                account_type = "storage"
            }
        else
            return false, "Failed to list storage accounts"
        end
    end)
    :timeout("120s")
    :build()

-- Azure Virtual Machines task
local azure_list_vms = task("azure_list_vms")
    :description("List Azure virtual machines")
    :depends_on({"azure_list_resource_groups"})
    :command(function(params, deps)
        log.info("🖥️  Listing Azure virtual machines...")
        
        local result = exec.run("az vm list --output table", {
            timeout = "90s",
            capture_output = true
        })
        
        if result.success then
            return true, result.output, {
                virtual_machines = result.output,
                dependency_completed = "resource_groups_listed"
            }
        else
            return false, "Failed to list virtual machines"
        end
    end)
    :timeout("150s")
    :build()

-- Modern Workflow Definition
workflow.define("azure_operations", {
    description = "Azure Operations Workflow - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"azure", "cloud", "resource-groups", "storage", "vms", "modern-dsl"},
        created_at = os.date(),
        prerequisites = "Azure CLI authenticated and configured"
    },
    
    tasks = {
        azure_list_rg,
        azure_list_storage,
        azure_list_vms
    },
    
    config = {
        timeout = "25m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting Azure operations workflow...")
        log.info("🔑 Ensure Azure CLI is authenticated")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Azure operations workflow completed successfully!")
            log.info("☁️  Azure resources have been inventoried")
        else
            log.error("❌ Azure operations workflow failed!")
            log.warn("🔍 Check Azure CLI authentication and permissions")
        end
        return true
    end
})
