-- Helm Advanced Examples
-- This file demonstrates advanced usage of the helm module

-- Example 1: Repository Management
local function manage_repositories()
    local helm = require("helm")
    
    -- Add popular Helm repositories
    local repos = {
        {"stable", "https://charts.helm.sh/stable"},
        {"bitnami", "https://charts.bitnami.com/bitnami"},
        {"ingress-nginx", "https://kubernetes.github.io/ingress-nginx"},
        {"prometheus-community", "https://prometheus-community.github.io/helm-charts"},
        {"grafana", "https://grafana.github.io/helm-charts"}
    }
    
    print("📦 Adding Helm repositories...")
    for _, repo in ipairs(repos) do
        local success, err = helm.repo_add(repo[1], repo[2], {
            force_update = true
        })
        
        if not success and not string.find(err, "already exists") then
            print("Warning: Failed to add repo " .. repo[1] .. ": " .. err)
        else
            print("✅ Added repository: " .. repo[1])
        end
    end
    
    -- Update repositories
    print("🔄 Updating repositories...")
    local update_success, err = helm.repo_update()
    if not update_success then
        print("Error updating repositories: " .. err)
        return false
    end
    
    print("✅ Repositories updated")
    
    -- List repositories
    print("📋 Listing repositories...")
    local repo_list, err = helm.repo_list({
        output = "table"
    })
    
    if err then
        print("Error listing repositories: " .. err)
    else
        print("Available repositories: " .. repo_list)
    end
    
    return true
end

-- Example 2: Search and Install Charts
local function search_and_install()
    local helm = require("helm")
    
    -- Search for charts
    print("🔍 Searching for nginx charts...")
    local search_results, err = helm.search_repo("nginx", {
        output = "table"
    })
    
    if err then
        print("Error searching repositories: " .. err)
    else
        print("Search results: " .. search_results)
    end
    
    -- Search Helm Hub
    print("🔍 Searching Helm Hub for prometheus...")
    local hub_results, err = helm.search_hub("prometheus", {
        max_results = "5"
    })
    
    if err then
        print("Error searching Helm Hub: " .. err)
    else
        print("Hub results: " .. hub_results)
    end
    
    -- Install a chart
    print("🚀 Installing nginx chart...")
    local install_success, err = helm.install("my-nginx", "bitnami/nginx", {
        namespace = "default",
        create_namespace = true,
        values = "service.type=ClusterIP,replicaCount=2",
        wait = true,
        timeout = "300s"
    })
    
    if not install_success then
        print("Error installing chart: " .. err)
        return false
    end
    
    print("✅ Chart installed successfully")
    return true
end

-- Example 3: Release Management
local function manage_releases()
    local helm = require("helm")
    
    -- List all releases
    print("📋 Listing all releases...")
    local releases, err = helm.list({
        all_namespaces = true,
        output = "table"
    })
    
    if err then
        print("Error listing releases: " .. err)
    else
        print("Current releases: " .. releases)
    end
    
    -- Get release status
    print("🔍 Getting release status...")
    local status, err = helm.status("my-nginx", {
        namespace = "default"
    })
    
    if err then
        print("Error getting status: " .. err)
    else
        print("Release status: " .. status)
    end
    
    -- Get release values
    print("📝 Getting release values...")
    local values, err = helm.get("values", "my-nginx", {
        namespace = "default"
    })
    
    if err then
        print("Error getting values: " .. err)
    else
        print("Release values: " .. values)
    end
    
    -- Get release manifest
    print("📄 Getting release manifest...")
    local manifest, err = helm.get("manifest", "my-nginx", {
        namespace = "default"
    })
    
    if err then
        print("Error getting manifest: " .. err)
    else
        print("Release manifest: " .. string.sub(manifest, 1, 500) .. "...")
    end
    
    return true
end

-- Example 4: Upgrade and Rollback
local function upgrade_and_rollback()
    local helm = require("helm")
    
    -- Upgrade the release with new values
    print("⬆️  Upgrading release...")
    local upgrade_success, err = helm.upgrade("my-nginx", "bitnami/nginx", {
        namespace = "default",
        values = "service.type=NodePort,replicaCount=3",
        wait = true,
        timeout = "300s"
    })
    
    if not upgrade_success then
        print("Error upgrading release: " .. err)
        return false
    end
    
    print("✅ Release upgraded successfully")
    
    -- Check release history
    print("📈 Checking release history...")
    local history, err = helm.history("my-nginx", {
        namespace = "default",
        max = "10"
    })
    
    if err then
        print("Error getting history: " .. err)
    else
        print("Release history: " .. history)
    end
    
    -- Rollback to previous version
    print("⏪ Rolling back to previous version...")
    local rollback_success, err = helm.rollback("my-nginx", "1", {
        namespace = "default",
        wait = true,
        timeout = "300s"
    })
    
    if not rollback_success then
        print("Error rolling back: " .. err)
        return false
    end
    
    print("✅ Rollback completed successfully")
    return true
end

-- Example 5: Custom Values and Configuration
local function custom_configuration()
    local helm = require("helm")
    
    -- Create a custom values file content
    local custom_values = [[
replicaCount: 3
image:
  repository: nginx
  tag: "1.21"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: "nginx"
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: my-app.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: my-app-tls
      hosts:
        - my-app.example.com

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80

nodeSelector: {}
tolerations: []
affinity: {}
]]
    
    -- Write values to a temporary file
    local fs = require("fs")
    local values_file = "/tmp/custom-values.yaml"
    local write_success, err = fs.write_file(values_file, custom_values)
    
    if not write_success then
        print("Error writing values file: " .. err)
        return false
    end
    
    -- Install with custom values file
    print("🚀 Installing with custom values...")
    local install_success, err = helm.install("custom-nginx", "bitnami/nginx", {
        namespace = "production",
        create_namespace = true,
        values_file = values_file,
        wait = true,
        timeout = "600s"
    })
    
    if not install_success then
        print("Error installing with custom values: " .. err)
        return false
    end
    
    print("✅ Chart installed with custom configuration")
    
    -- Clean up values file
    os.remove(values_file)
    
    return true
end

-- Example 6: Chart Development
local function chart_development()
    local helm = require("helm")
    
    -- Create a new chart
    print("📦 Creating new chart...")
    local create_success, err = helm.create("my-app")
    if not create_success then
        print("Error creating chart: " .. err)
        return false
    end
    
    print("✅ Chart created successfully")
    
    -- Lint the chart
    print("🔍 Linting chart...")
    local lint_success, err = helm.lint("./my-app", {
        strict = true
    })
    
    if not lint_success then
        print("Chart linting failed: " .. err)
        return false
    end
    
    print("✅ Chart linting passed")
    
    -- Template the chart
    print("📄 Templating chart...")
    local template_output, err = helm.template("my-app", "./my-app", {
        namespace = "test"
    })
    
    if err then
        print("Error templating chart: " .. err)
        return false
    end
    
    print("Template output: " .. string.sub(template_output, 1, 500) .. "...")
    
    -- Package the chart
    print("📦 Packaging chart...")
    local package_success, err = helm.package("./my-app", {
        destination = "./charts",
        version = "0.1.0",
        app_version = "1.0.0"
    })
    
    if not package_success then
        print("Error packaging chart: " .. err)
        return false
    end
    
    print("✅ Chart packaged successfully")
    
    -- Clean up
    os.execute("rm -rf ./my-app ./charts")
    
    return true
end

-- Example 7: Testing Charts
local function test_charts()
    local helm = require("helm")
    
    -- Install a chart with tests
    print("🚀 Installing chart for testing...")
    local install_success, err = helm.install("test-app", "bitnami/nginx", {
        namespace = "testing",
        create_namespace = true,
        wait = true,
        timeout = "300s"
    })
    
    if not install_success then
        print("Error installing test chart: " .. err)
        return false
    end
    
    -- Run tests
    print("🧪 Running chart tests...")
    local test_success, err = helm.test("test-app", {
        namespace = "testing",
        logs = true,
        timeout = "300s"
    })
    
    if not test_success then
        print("Chart tests failed: " .. err)
        return false
    end
    
    print("✅ Chart tests passed")
    
    -- Cleanup test release
    helm.uninstall("test-app", {namespace = "testing"})
    
    return true
end

-- Example 8: Plugin Management
local function manage_plugins()
    local helm = require("helm")
    
    -- List installed plugins
    print("🔌 Listing installed plugins...")
    local plugins, err = helm.plugin_list()
    if err then
        print("Error listing plugins: " .. err)
    else
        print("Installed plugins: " .. plugins)
    end
    
    -- Install a useful plugin (diff)
    print("📦 Installing helm-diff plugin...")
    local install_success, err = helm.plugin_install("https://github.com/databus23/helm-diff")
    
    if not install_success and not string.find(err, "already exists") then
        print("Warning: Failed to install plugin: " .. err)
    else
        print("✅ Plugin installed successfully")
    end
    
    -- Install secrets plugin
    print("📦 Installing helm-secrets plugin...")
    local secrets_success, err = helm.plugin_install("https://github.com/jkroepke/helm-secrets")
    
    if not secrets_success and not string.find(err, "already exists") then
        print("Warning: Failed to install secrets plugin: " .. err)
    else
        print("✅ Secrets plugin installed")
    end
    
    return true
end

-- Example 9: Multi-Environment Deployment
local function multi_environment_deployment()
    local helm = require("helm")
    
    local environments = {
        {name = "development", replicas = "1", resources = "small"},
        {name = "staging", replicas = "2", resources = "medium"},
        {name = "production", replicas = "3", resources = "large"}
    }
    
    for _, env in ipairs(environments) do
        print("🚀 Deploying to " .. env.name .. " environment...")
        
        -- Environment-specific values
        local values = string.format(
            "replicaCount=%s,environment=%s,resources.preset=%s",
            env.replicas, env.name, env.resources
        )
        
        local install_success, err = helm.install("app-" .. env.name, "bitnami/nginx", {
            namespace = env.name,
            create_namespace = true,
            values = values,
            wait = true,
            timeout = "300s"
        })
        
        if not install_success then
            print("Error deploying to " .. env.name .. ": " .. err)
        else
            print("✅ Deployed to " .. env.name .. " successfully")
        end
    end
    
    return true
end

-- Example 10: Complete Helm Workflow
local function complete_helm_workflow()
    local helm = require("helm")
    
    print("🚀 Starting complete Helm workflow...")
    
    -- Step 1: Repository management
    print("Step 1: Managing repositories...")
    if not manage_repositories() then
        return false
    end
    
    -- Step 2: Search and install
    print("Step 2: Searching and installing...")
    if not search_and_install() then
        return false
    end
    
    -- Step 3: Release management
    print("Step 3: Managing releases...")
    if not manage_releases() then
        return false
    end
    
    -- Step 4: Upgrade and rollback
    print("Step 4: Upgrading and rolling back...")
    if not upgrade_and_rollback() then
        return false
    end
    
    -- Step 5: Custom configuration
    print("Step 5: Custom configuration...")
    if not custom_configuration() then
        return false
    end
    
    -- Step 6: Plugin management
    print("Step 6: Managing plugins...")
    manage_plugins()
    
    -- Step 7: Multi-environment deployment
    print("Step 7: Multi-environment deployment...")
    multi_environment_deployment()
    
    print("✅ Complete Helm workflow finished!")
    return true
end

-- Example 11: Cleanup All Resources
local function cleanup_all()
    local helm = require("helm")
    
    print("🧹 Cleaning up all Helm resources...")
    
    -- Get all releases
    local releases, err = helm.list({
        all_namespaces = true,
        output = "json"
    })
    
    if err then
        print("Error listing releases for cleanup: " .. err)
        return false
    end
    
    -- Uninstall specific releases
    local release_names = {
        "my-nginx",
        "custom-nginx", 
        "app-development",
        "app-staging",
        "app-production"
    }
    
    for _, release in ipairs(release_names) do
        print("🗑️  Uninstalling " .. release .. "...")
        local uninstall_success, err = helm.uninstall(release, {
            namespace = "default"
        })
        
        if not uninstall_success then
            print("Warning: Failed to uninstall " .. release .. ": " .. (err or "unknown error"))
        else
            print("✅ Uninstalled " .. release)
        end
    end
    
    -- Clean up namespaces (using kubectl)
    local k8s = require("kubernetes")
    local namespaces = {"production", "testing", "development", "staging"}
    
    for _, ns in ipairs(namespaces) do
        k8s.delete_namespace(ns)
        print("✅ Cleaned up namespace: " .. ns)
    end
    
    print("✅ Cleanup completed")
    return true
end

-- Example 12: Dependency Management
local function manage_dependencies()
    local helm = require("helm")
    
    -- This example shows how to work with chart dependencies
    print("📦 Managing chart dependencies...")
    
    -- Create a chart with dependencies
    local chart_success, err = helm.create("web-app")
    if not chart_success then
        print("Error creating chart: " .. err)
        return false
    end
    
    -- In a real scenario, you would modify Chart.yaml to include dependencies
    -- For demo purposes, we'll work with the created chart
    
    -- Build dependencies
    print("🔨 Building dependencies...")
    local build_success, err = helm.dependency("build", "./web-app")
    if not build_success then
        print("Warning: Dependency build failed: " .. err)
    else
        print("✅ Dependencies built successfully")
    end
    
    -- List dependencies
    print("📋 Listing dependencies...")
    local list_success, err = helm.dependency("list", "./web-app")
    if not list_success then
        print("Warning: Failed to list dependencies: " .. err)
    else
        print("Dependencies: " .. err) -- err contains the output in this case
    end
    
    -- Update dependencies
    print("🔄 Updating dependencies...")
    local update_success, err = helm.dependency("update", "./web-app")
    if not update_success then
        print("Warning: Dependency update failed: " .. err)
    else
        print("✅ Dependencies updated successfully")
    end
    
    -- Clean up
    os.execute("rm -rf ./web-app")
    
    return true
end

-- Export task definitions
TaskDefinitions = {
    helm_examples = {
        description = "Advanced Helm package manager examples",
        workdir = ".",
        tasks = {
            manage_repos = {
                description = "Add and manage Helm repositories",
                command = manage_repositories
            },
            search_install = {
                description = "Search and install Helm charts",
                command = search_and_install
            },
            manage_releases = {
                description = "Manage Helm releases",
                command = manage_releases
            },
            upgrade_rollback = {
                description = "Upgrade and rollback releases",
                command = upgrade_and_rollback
            },
            custom_config = {
                description = "Deploy with custom configuration",
                command = custom_configuration
            },
            chart_development = {
                description = "Create and develop charts",
                command = chart_development
            },
            test_charts = {
                description = "Test Helm charts",
                command = test_charts
            },
            manage_plugins = {
                description = "Install and manage Helm plugins",
                command = manage_plugins
            },
            multi_env_deploy = {
                description = "Deploy to multiple environments",
                command = multi_environment_deployment
            },
            complete_workflow = {
                description = "Complete Helm workflow",
                command = complete_helm_workflow
            },
            manage_dependencies = {
                description = "Manage chart dependencies",
                command = manage_dependencies
            },
            cleanup_all = {
                description = "Cleanup all Helm resources",
                command = cleanup_all
            }
        }
    }
}