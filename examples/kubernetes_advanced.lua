-- Kubernetes Advanced Examples
-- This file demonstrates advanced usage of the kubernetes module

-- Example 1: Basic Kubernetes Operations
local function basic_k8s_operations()
    local k8s = require("kubernetes")
    
    -- Get cluster nodes
    print("🔍 Getting cluster nodes...")
    local nodes, err = k8s.get_nodes()
    if err then
        print("Error getting nodes: " .. err)
        return false
    end
    
    print("Cluster nodes: " .. nodes)
    
    -- Create a namespace
    print("📂 Creating namespace...")
    local ns_success, err = k8s.create_namespace("demo-app", {
        dry_run = false
    })
    
    if not ns_success and not string.find(err, "already exists") then
        print("Error creating namespace: " .. err)
        return false
    end
    
    print("✅ Namespace created successfully")
    
    -- Get pods in all namespaces
    print("🔍 Getting all pods...")
    local pods, err = k8s.get_pods({
        all_namespaces = true
    })
    
    if err then
        print("Error getting pods: " .. err)
    else
        print("All pods: " .. string.sub(pods, 1, 500) .. "...")
    end
    
    return true
end

-- Example 2: Deploy Application with YAML
local function deploy_application()
    local k8s = require("kubernetes")
    
    -- Deploy a simple nginx application
    local nginx_yaml = [[
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: demo-app
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.20
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
  namespace: demo-app
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
  type: ClusterIP
]]
    
    print("🚀 Deploying nginx application...")
    local deploy_success, err = k8s.apply(nginx_yaml, {
        namespace = "demo-app"
    })
    
    if not deploy_success then
        print("Error deploying application: " .. err)
        return false
    end
    
    print("✅ Application deployed successfully")
    
    -- Wait for deployment to be ready
    print("⏳ Waiting for deployment to be ready...")
    local ready, err = k8s.wait_for_ready("deployment", "nginx-deployment", {
        namespace = "demo-app",
        timeout = "300s"
    })
    
    if not ready then
        print("Warning: Deployment not ready: " .. (err or "timeout"))
    else
        print("✅ Deployment is ready")
    end
    
    return true
end

-- Example 3: Scale and Manage Deployments
local function manage_deployments()
    local k8s = require("kubernetes")
    
    -- Get current deployments
    print("📊 Getting current deployments...")
    local deployments, err = k8s.get_deployments({
        namespace = "demo-app"
    })
    
    if err then
        print("Error getting deployments: " .. err)
        return false
    end
    
    print("Current deployments: " .. deployments)
    
    -- Scale the deployment
    print("📈 Scaling deployment to 5 replicas...")
    local scale_success, err = k8s.scale("deployment", "nginx-deployment", "5", {
        namespace = "demo-app"
    })
    
    if not scale_success then
        print("Error scaling deployment: " .. err)
        return false
    end
    
    print("✅ Deployment scaled successfully")
    
    -- Check rollout status
    print("🔍 Checking rollout status...")
    local status, err = k8s.rollout("status", "deployment", "nginx-deployment", {
        namespace = "demo-app"
    })
    
    if err then
        print("Error checking rollout status: " .. err)
    else
        print("Rollout status: " .. status)
    end
    
    -- Get pod logs
    print("📋 Getting pod logs...")
    local pods, err = k8s.get_pods({
        namespace = "demo-app",
        selector = "app=nginx"
    })
    
    if not err then
        -- Parse the first pod name (in a real scenario, you'd parse JSON)
        local logs, err = k8s.logs("nginx-deployment-", {
            namespace = "demo-app",
            tail = "10"
        })
        if logs then
            print("Recent logs: " .. logs)
        end
    end
    
    return true
end

-- Example 4: ConfigMaps and Secrets
local function manage_config_and_secrets()
    local k8s = require("kubernetes")
    
    -- Create a ConfigMap
    print("📝 Creating ConfigMap...")
    local cm_success, err = k8s.create_configmap("app-config", {
        namespace = "demo-app",
        from_literal = "database_host=mysql.example.com"
    })
    
    if not cm_success and not string.find(err, "already exists") then
        print("Error creating ConfigMap: " .. err)
        return false
    end
    
    print("✅ ConfigMap created")
    
    -- Create a Secret
    print("🔐 Creating Secret...")
    local secret_success, err = k8s.create_secret("generic", "app-secrets", {
        namespace = "demo-app",
        from_literal = "database_password=supersecret"
    })
    
    if not secret_success and not string.find(err, "already exists") then
        print("Error creating Secret: " .. err)
        return false
    end
    
    print("✅ Secret created")
    
    -- List ConfigMaps and Secrets
    local configs, err = k8s.get("configmaps", {
        namespace = "demo-app"
    })
    
    if err then
        print("Error getting ConfigMaps: " .. err)
    else
        print("ConfigMaps: " .. configs)
    end
    
    local secrets, err = k8s.get("secrets", {
        namespace = "demo-app"
    })
    
    if err then
        print("Error getting Secrets: " .. err)
    else
        print("Secrets: " .. secrets)
    end
    
    return true
end

-- Example 5: Labels and Annotations
local function manage_labels_and_annotations()
    local k8s = require("kubernetes")
    
    -- Add labels to deployment
    print("🏷️  Adding labels to deployment...")
    local label_success, err = k8s.label("deployment", "nginx-deployment", "environment=demo version=1.0", {
        namespace = "demo-app",
        overwrite = true
    })
    
    if not label_success then
        print("Error adding labels: " .. err)
        return false
    end
    
    print("✅ Labels added successfully")
    
    -- Add annotations
    print("📋 Adding annotations...")
    local annotate_success, err = k8s.annotate("deployment", "nginx-deployment", "deployment.kubernetes.io/revision=1", {
        namespace = "demo-app",
        overwrite = true
    })
    
    if not annotate_success then
        print("Error adding annotations: " .. err)
        return false
    end
    
    print("✅ Annotations added successfully")
    
    -- Get detailed information about the deployment
    print("🔍 Getting deployment details...")
    local details, err = k8s.describe("deployment", "nginx-deployment", {
        namespace = "demo-app"
    })
    
    if err then
        print("Error describing deployment: " .. err)
    else
        print("Deployment details: " .. string.sub(details, 1, 500) .. "...")
    end
    
    return true
end

-- Example 6: Resource Monitoring and Events
local function monitor_resources()
    local k8s = require("kubernetes")
    
    -- Get resource usage
    print("📊 Getting node resource usage...")
    local node_usage, err = k8s.top("node")
    if err then
        print("Error getting node usage: " .. err)
    else
        print("Node usage: " .. node_usage)
    end
    
    -- Get pod resource usage
    print("📊 Getting pod resource usage...")
    local pod_usage, err = k8s.top("pod", {
        namespace = "demo-app"
    })
    
    if err then
        print("Error getting pod usage: " .. err)
    else
        print("Pod usage: " .. pod_usage)
    end
    
    -- Get recent events
    print("📰 Getting recent events...")
    local events, err = k8s.events({
        namespace = "demo-app"
    })
    
    if err then
        print("Error getting events: " .. err)
    else
        print("Recent events: " .. string.sub(events, 1, 1000) .. "...")
    end
    
    return true
end

-- Example 7: Execute Commands in Pods
local function execute_in_pods()
    local k8s = require("kubernetes")
    
    -- Get running pods
    print("🔍 Getting running pods...")
    local pods, err = k8s.get_pods({
        namespace = "demo-app",
        selector = "app=nginx"
    })
    
    if err then
        print("Error getting pods: " .. err)
        return false
    end
    
    -- Execute command in pod (assuming we have the pod name)
    print("💻 Executing command in pod...")
    local exec_result, err = k8s.exec("nginx-deployment-", "nginx -v", {
        namespace = "demo-app"
    })
    
    if err then
        print("Error executing command: " .. err)
    else
        print("Command output: " .. exec_result)
    end
    
    -- Get logs from specific container
    print("📋 Getting container logs...")
    local logs, err = k8s.logs("nginx-deployment-", {
        namespace = "demo-app",
        container = "nginx",
        tail = "20",
        since = "1h"
    })
    
    if err then
        print("Error getting logs: " .. err)
    else
        print("Container logs: " .. logs)
    end
    
    return true
end

-- Example 8: Patch Resources
local function patch_resources()
    local k8s = require("kubernetes")
    
    -- Patch deployment to change image
    print("🔧 Patching deployment image...")
    local patch = '{"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"nginx:1.21"}]}}}}'
    
    local patch_success, err = k8s.patch("deployment", "nginx-deployment", patch, {
        namespace = "demo-app",
        type = "merge"
    })
    
    if not patch_success then
        print("Error patching deployment: " .. err)
        return false
    end
    
    print("✅ Deployment patched successfully")
    
    -- Check rollout status
    print("🔍 Checking rollout progress...")
    local status, err = k8s.rollout("status", "deployment", "nginx-deployment", {
        namespace = "demo-app"
    })
    
    if err then
        print("Error checking rollout: " .. err)
    else
        print("Rollout status: " .. status)
    end
    
    return true
end

-- Example 9: Port Forwarding (for development)
local function setup_port_forwarding()
    local k8s = require("kubernetes")
    
    -- Set up port forwarding to service
    print("🔗 Setting up port forwarding...")
    local pf_success, err = k8s.port_forward("service", "nginx-service", "8080:80", {
        namespace = "demo-app"
    })
    
    if not pf_success then
        print("Error setting up port forwarding: " .. err)
        return false
    end
    
    print("✅ Port forwarding established on localhost:8080")
    print("ℹ️  You can now access the service at http://localhost:8080")
    
    return true
end

-- Example 10: Complete Application Lifecycle
local function complete_app_lifecycle()
    local k8s = require("kubernetes")
    
    print("🚀 Starting complete application lifecycle...")
    
    -- Step 1: Basic setup
    print("Step 1: Basic setup...")
    if not basic_k8s_operations() then
        return false
    end
    
    -- Step 2: Deploy application
    print("Step 2: Deploying application...")
    if not deploy_application() then
        return false
    end
    
    -- Step 3: Configure resources
    print("Step 3: Configuring resources...")
    if not manage_config_and_secrets() then
        return false
    end
    
    -- Step 4: Add metadata
    print("Step 4: Adding metadata...")
    if not manage_labels_and_annotations() then
        return false
    end
    
    -- Step 5: Scale and manage
    print("Step 5: Scaling and managing...")
    if not manage_deployments() then
        return false
    end
    
    -- Step 6: Monitor resources
    print("Step 6: Monitoring resources...")
    monitor_resources()
    
    -- Step 7: Update application
    print("Step 7: Updating application...")
    patch_resources()
    
    print("✅ Application lifecycle completed successfully!")
    return true
end

-- Example 11: Cleanup Resources
local function cleanup_resources()
    local k8s = require("kubernetes")
    
    print("🧹 Cleaning up resources...")
    
    -- Delete deployment
    local del_deploy, err = k8s.delete("deployment", "nginx-deployment", {
        namespace = "demo-app"
    })
    
    if not del_deploy then
        print("Error deleting deployment: " .. err)
    else
        print("✅ Deployment deleted")
    end
    
    -- Delete service
    local del_svc, err = k8s.delete("service", "nginx-service", {
        namespace = "demo-app"
    })
    
    if not del_svc then
        print("Error deleting service: " .. err)
    else
        print("✅ Service deleted")
    end
    
    -- Delete ConfigMap and Secret
    k8s.delete("configmap", "app-config", {namespace = "demo-app"})
    k8s.delete("secret", "app-secrets", {namespace = "demo-app"})
    
    -- Delete namespace
    local del_ns, err = k8s.delete_namespace("demo-app")
    if not del_ns then
        print("Error deleting namespace: " .. err)
    else
        print("✅ Namespace deleted")
    end
    
    print("✅ Cleanup completed")
    return true
end

-- Export task definitions
TaskDefinitions = {
    kubernetes_examples = {
        description = "Advanced Kubernetes management examples",
        workdir = ".",
        tasks = {
            basic_ops = {
                description = "Basic Kubernetes operations",
                command = basic_k8s_operations
            },
            deploy_app = {
                description = "Deploy application with YAML",
                command = deploy_application
            },
            manage_deployments = {
                description = "Scale and manage deployments",
                command = manage_deployments
            },
            config_secrets = {
                description = "Manage ConfigMaps and Secrets",
                command = manage_config_and_secrets
            },
            labels_annotations = {
                description = "Manage labels and annotations",
                command = manage_labels_and_annotations
            },
            monitor_resources = {
                description = "Monitor resource usage and events",
                command = monitor_resources
            },
            execute_in_pods = {
                description = "Execute commands in pods",
                command = execute_in_pods
            },
            patch_resources = {
                description = "Patch and update resources",
                command = patch_resources
            },
            port_forward = {
                description = "Setup port forwarding",
                command = setup_port_forwarding
            },
            complete_lifecycle = {
                description = "Complete application lifecycle",
                command = complete_app_lifecycle
            },
            cleanup = {
                description = "Cleanup all resources",
                command = cleanup_resources
            }
        }
    }
}