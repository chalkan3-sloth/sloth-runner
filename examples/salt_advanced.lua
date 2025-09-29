-- SaltStack Advanced Examples
-- This file demonstrates advanced usage of the salt module

-- Example 1: Basic Salt Command Execution
local function test_salt_connectivity()
    local salt = require("salt")
    
    -- Test connectivity to all minions
    local result, err = salt.test_ping("*")
    if err then
        print("Error testing connectivity: " .. err)
        return false
    end
    
    print("Connectivity test result: " .. result)
    return true
end

-- Example 2: State Management
local function deploy_nginx_state()
    local salt = require("salt")
    
    -- Apply nginx state to web servers
    local result, err = salt.state_apply("web*", "nginx", {
        pillar = "nginx_version: '1.20.2'"
    })
    
    if err then
        print("Error applying nginx state: " .. err)
        return false
    end
    
    print("Nginx deployment result: " .. result)
    return true
end

-- Example 3: Package Management
local function install_packages()
    local salt = require("salt")
    
    -- Install vim on all Ubuntu servers
    local result, err = salt.pkg("ubuntu*", "install", "vim")
    if err then
        print("Error installing vim: " .. err)
        return false
    end
    
    print("Package installation result: " .. result)
    
    -- List installed packages
    local packages, err = salt.pkg("ubuntu*", "list_pkgs")
    if err then
        print("Error listing packages: " .. err)
        return false
    end
    
    print("Installed packages: " .. packages)
    return true
end

-- Example 4: Service Management
local function manage_services()
    local salt = require("salt")
    
    -- Start nginx service
    local result, err = salt.service("web*", "start", "nginx")
    if err then
        print("Error starting nginx: " .. err)
        return false
    end
    
    -- Enable nginx service to start at boot
    local enable_result, err = salt.service("web*", "enable", "nginx")
    if err then
        print("Error enabling nginx: " .. err)
        return false
    end
    
    -- Check service status
    local status, err = salt.service("web*", "status", "nginx")
    if err then
        print("Error checking nginx status: " .. err)
        return false
    end
    
    print("Service status: " .. status)
    return true
end

-- Example 5: User Management
local function create_users()
    local salt = require("salt")
    
    -- Create a new user
    local result, err = salt.user("*", "add", "deploy", {
        home = "/home/deploy",
        shell = "/bin/bash"
    })
    
    if err then
        print("Error creating user: " .. err)
        return false
    end
    
    -- Get user information
    local user_info, err = salt.user("*", "info", "deploy")
    if err then
        print("Error getting user info: " .. err)
        return false
    end
    
    print("User info: " .. user_info)
    return true
end

-- Example 6: File Management
local function manage_files()
    local salt = require("salt")
    
    -- Check if file exists
    local exists, err = salt.file("*", "exists", "/etc/nginx/nginx.conf")
    if err then
        print("Error checking file: " .. err)
        return false
    end
    
    -- Copy configuration file
    if exists then
        local copy_result, err = salt.file("*", "copy", "/etc/nginx/nginx.conf", "/etc/nginx/nginx.conf.backup")
        if err then
            print("Error copying file: " .. err)
            return false
        end
        
        -- Set file permissions
        local perm_result, err = salt.file("*", "set_mode", "/etc/nginx/nginx.conf.backup", "644")
        if err then
            print("Error setting permissions: " .. err)
            return false
        end
    end
    
    return true
end

-- Example 7: Asynchronous Operations
local function async_operations()
    local salt = require("salt")
    
    -- Start a long-running operation asynchronously
    local job_id, err = salt.async_run("*", "pkg", "upgrade")
    if err then
        print("Error starting async operation: " .. err)
        return false
    end
    
    print("Started async job: " .. job_id)
    
    -- Check job status periodically
    for i = 1, 10 do
        os.execute("sleep 5")  -- Wait 5 seconds
        
        local status, err = salt.job_status(job_id)
        if err then
            print("Error checking job status: " .. err)
            break
        end
        
        print("Job status (attempt " .. i .. "): " .. status)
        
        -- Parse the status to see if job is complete
        local status_data = require("data").parse_json(status)
        if status_data and status_data ~= {} then
            print("Job completed!")
            break
        end
    end
    
    return true
end

-- Example 8: Grains and Pillar Data
local function get_system_info()
    local salt = require("salt")
    
    -- Get all grains (system information)
    local grains, err = salt.grains_get("*")
    if err then
        print("Error getting grains: " .. err)
        return false
    end
    
    print("System grains: " .. grains)
    
    -- Get specific grain
    local os_info, err = salt.grains_get("*", "os")
    if err then
        print("Error getting OS info: " .. err)
        return false
    end
    
    print("Operating system: " .. os_info)
    
    -- Get pillar data
    local pillar_data, err = salt.pillar_get("*", "users")
    if err then
        print("Error getting pillar data: " .. err)
        return false
    end
    
    print("Pillar users: " .. pillar_data)
    return true
end

-- Example 9: Complete Infrastructure Setup
local function complete_setup()
    local salt = require("salt")
    
    print("🚀 Starting complete infrastructure setup...")
    
    -- Step 1: Test connectivity
    print("Step 1: Testing connectivity...")
    if not test_salt_connectivity() then
        return false
    end
    
    -- Step 2: Update packages
    print("Step 2: Updating packages...")
    local update_result, err = salt.pkg("*", "upgrade")
    if err then
        print("Warning: Package update failed: " .. err)
    end
    
    -- Step 3: Install required packages
    print("Step 3: Installing required packages...")
    local packages = {"nginx", "git", "htop", "curl"}
    for _, pkg in ipairs(packages) do
        local result, err = salt.pkg("*", "install", pkg)
        if err then
            print("Warning: Failed to install " .. pkg .. ": " .. err)
        else
            print("Installed " .. pkg .. " successfully")
        end
    end
    
    -- Step 4: Configure services
    print("Step 4: Configuring services...")
    local services = {"nginx"}
    for _, service in ipairs(services) do
        salt.service("*", "start", service)
        salt.service("*", "enable", service)
        print("Configured " .. service .. " service")
    end
    
    -- Step 5: Apply high state
    print("Step 5: Applying high state...")
    local highstate_result, err = salt.state_highstate("*", {test = false})
    if err then
        print("Warning: High state application failed: " .. err)
    else
        print("High state applied successfully")
    end
    
    print("✅ Infrastructure setup completed!")
    return true
end

-- Export task definitions
TaskDefinitions = {
    salt_examples = {
        description = "Advanced SaltStack management examples",
        workdir = ".",
        tasks = {
            test_connectivity = {
                description = "Test Salt connectivity to all minions",
                command = test_salt_connectivity
            },
            deploy_nginx = {
                description = "Deploy nginx using Salt states",
                command = deploy_nginx_state
            },
            manage_packages = {
                description = "Install and manage packages",
                command = install_packages
            },
            manage_services = {
                description = "Start and manage services",
                command = manage_services
            },
            create_users = {
                description = "Create and manage users",
                command = create_users
            },
            manage_files = {
                description = "File operations with Salt",
                command = manage_files
            },
            async_operations = {
                description = "Demonstrate asynchronous operations",
                command = async_operations
            },
            system_info = {
                description = "Get system information via grains and pillar",
                command = get_system_info
            },
            complete_setup = {
                description = "Complete infrastructure setup",
                command = complete_setup
            }
        }
    }
}