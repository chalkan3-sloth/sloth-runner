# Module API Examples - Modern DSL

This document provides comprehensive examples of all Sloth Runner modules using the modern DSL with table-based parameters.

## Core Principles

All modules in Sloth Runner follow consistent patterns:

1. **Table-based parameters**: All functions accept a single table parameter with named fields
2. **Delegate execution**: Use `:delegate_to()` to execute commands on remote agents
3. **Consistent error handling**: Functions return success/failure with descriptive messages

---

## Package Management (pkg)

### Installing Packages

```lua
task("install_packages", {
    description = "Install required packages",
    command = function()
        -- Install single package
        local result = pkg.install({
            packages = {"nginx"},
            delegate_to = "web-server"
        })
        
        if not result.success then
            return false, "Installation failed: " .. result.message
        end
        
        return true, "Packages installed successfully"
    end
})
```

### Multiple Operations

```lua
task("manage_packages", {
    description = "Complete package management",
    command = function()
        -- Update package database
        pkg.update({ delegate_to = "web-server" })
        
        -- Install multiple packages
        pkg.install({
            packages = {"nginx", "postgresql", "redis"},
            delegate_to = "web-server"
        })
        
        -- Upgrade all packages
        pkg.upgrade({ delegate_to = "web-server" })
        
        -- Remove package
        pkg.remove({
            packages = {"apache2"},
            delegate_to = "web-server"
        })
        
        return true, "Package management completed"
    end
})
```

---

## Systemd Management

### Service Operations

```lua
task("manage_services", {
    description = "Manage systemd services",
    command = function()
        -- Start service
        systemd.start({
            service = "nginx",
            delegate_to = "web-server"
        })
        
        -- Enable service
        systemd.enable({
            service = "nginx",
            delegate_to = "web-server"
        })
        
        -- Check status
        local status = systemd.status({
            service = "nginx",
            delegate_to = "web-server"
        })
        
        log.info("Service status: " .. status.state)
        
        return true, "Service management completed"
    end
})
```

### Restart with Verification

```lua
task("restart_service", {
    description = "Restart and verify service",
    command = function()
        -- Restart service
        systemd.restart({
            service = "nginx",
            delegate_to = "web-server"
        })
        
        -- Wait a moment for service to start
        goroutine.sleep(2000)
        
        -- Verify it's running
        local status = systemd.is_active({
            service = "nginx",
            delegate_to = "web-server"
        })
        
        if not status then
            return false, "Service failed to start"
        end
        
        return true, "Service restarted successfully"
    end
})
```

---

## Terraform

### Complete Infrastructure Lifecycle

```lua
task("terraform_deploy", {
    description = "Deploy infrastructure with Terraform",
    command = function()
        local workdir = "./terraform/production"
        
        -- Initialize
        local init_result = terraform.init({
            workdir = workdir,
            delegate_to = "deploy-agent"
        })
        if not init_result.success then
            return false, "Init failed: " .. init_result.stderr
        end
        
        -- Plan
        local plan_result = terraform.plan({
            workdir = workdir,
            out = "prod.tfplan",
            delegate_to = "deploy-agent"
        })
        if not plan_result.success then
            return false, "Plan failed: " .. plan_result.stderr
        end
        
        -- Apply
        local apply_result = terraform.apply({
            workdir = workdir,
            plan = "prod.tfplan",
            auto_approve = true,
            delegate_to = "deploy-agent"
        })
        if not apply_result.success then
            return false, "Apply failed: " .. apply_result.stderr
        end
        
        -- Get outputs
        local vpc_id, err = terraform.output({
            workdir = workdir,
            name = "vpc_id",
            delegate_to = "deploy-agent"
        })
        if not vpc_id then
            return false, "Failed to get output: " .. err
        end
        
        log.info("VPC ID: " .. vpc_id)
        return true, "Infrastructure deployed successfully"
    end
})
```

---

## Pulumi

### Stack Management

```lua
task("pulumi_deploy", {
    description = "Deploy with Pulumi",
    command = function()
        -- Create stack object
        local stack = pulumi.stack({
            name = "my-org/production/us-east-1",
            workdir = "./pulumi/infra",
            delegate_to = "deploy-agent"
        })
        
        -- Preview changes
        local preview = stack:preview({
            delegate_to = "deploy-agent"
        })
        log.info("Preview:\n" .. preview.stdout)
        
        -- Deploy
        local result = stack:up({
            yes = true,
            config = {
                ["aws:region"] = "us-east-1",
                ["myapp:environment"] = "production"
            },
            delegate_to = "deploy-agent"
        })
        
        if not result.success then
            return false, "Deploy failed: " .. result.stderr
        end
        
        -- Get outputs
        local outputs, err = stack:outputs({
            delegate_to = "deploy-agent"
        })
        if err then
            return false, "Failed to get outputs: " .. err
        end
        
        log.info("Endpoint: " .. outputs.endpoint)
        return true, "Pulumi deployment completed"
    end
})
```

---

## AWS

### S3 Sync with aws-vault

```lua
task("deploy_to_s3", {
    description = "Deploy static site to S3",
    command = function()
        local result = aws.s3.sync({
            source = "./build",
            destination = "s3://my-app-bucket/static",
            profile = "production",
            delete = true,
            delegate_to = "deploy-agent"
        })
        
        if not result then
            return false, "S3 sync failed"
        end
        
        return true, "Site deployed to S3"
    end
})
```

### Secrets Manager

```lua
task("fetch_secrets", {
    description = "Fetch secrets from AWS",
    command = function()
        local db_password, err = aws.secretsmanager.get_secret({
            secret_id = "production/database/password",
            profile = "production",
            delegate_to = "app-server"
        })
        
        if not db_password then
            return false, "Failed to fetch secret: " .. err
        end
        
        -- Use the secret to configure application
        log.info("Secret retrieved successfully")
        return true, "Secrets configured"
    end
})
```

---

## Docker

### Build and Push Pipeline

```lua
task("docker_pipeline", {
    description = "Build and push Docker image",
    command = function()
        local image_tag = "myapp:v1.2.3"
        
        -- Build
        local build_result = docker.build({
            tag = image_tag,
            path = "./app",
            dockerfile = "./app/Dockerfile",
            build_args = {
                VERSION = "1.2.3",
                BUILD_DATE = os.date("%Y-%m-%d")
            },
            delegate_to = "build-agent"
        })
        
        if not build_result.success then
            return false, "Build failed: " .. build_result.stderr
        end
        
        -- Push
        local push_result = docker.push({
            tag = image_tag,
            delegate_to = "build-agent"
        })
        
        if not push_result.success then
            return false, "Push failed: " .. push_result.stderr
        end
        
        return true, "Docker image built and pushed"
    end
})
```

---

## File Operations

### Copy with Template Rendering

```lua
task("deploy_config", {
    description = "Deploy configuration files",
    command = function()
        -- Copy file
        fs.copy({
            src = "./configs/app.conf",
            dest = "/etc/myapp/app.conf",
            mode = "0644",
            owner = "myapp",
            group = "myapp",
            delegate_to = "app-server"
        })
        
        -- Render and copy template
        fs.template({
            src = "./templates/database.conf.tmpl",
            dest = "/etc/myapp/database.conf",
            vars = {
                db_host = "db.example.com",
                db_port = 5432,
                db_name = "production"
            },
            mode = "0600",
            owner = "myapp",
            delegate_to = "app-server"
        })
        
        return true, "Configuration deployed"
    end
})
```

### Archive Operations

```lua
task("backup_data", {
    description = "Create backup archive",
    command = function()
        -- Create archive
        fs.unarchive({
            src = "/tmp/backup.tar.gz",
            dest = "/var/backups/",
            remote_src = true,
            delegate_to = "backup-server"
        })
        
        return true, "Backup extracted"
    end
})
```

---

## User Management

### Complete User Setup

```lua
task("setup_user", {
    description = "Create and configure user",
    command = function()
        -- Create user
        user.create({
            name = "appuser",
            uid = 1001,
            shell = "/bin/bash",
            home = "/home/appuser",
            groups = {"sudo", "docker"},
            delegate_to = "app-server"
        })
        
        -- Set password
        user.password({
            name = "appuser",
            password = "hashed_password_here",
            delegate_to = "app-server"
        })
        
        -- Add SSH key
        user.authorized_key({
            user = "appuser",
            key = "ssh-rsa AAAAB3NzaC1yc2...",
            state = "present",
            delegate_to = "app-server"
        })
        
        return true, "User configured"
    end
})
```

---

## SSH Management

### SSH Configuration

```lua
task("configure_ssh", {
    description = "Configure SSH access",
    command = function()
        -- Add known host
        ssh.known_host({
            name = "github.com",
            key = "github.com ssh-rsa AAAAB3NzaC1yc2...",
            delegate_to = "app-server"
        })
        
        -- Configure SSH client
        ssh.config({
            host = "production",
            hostname = "prod.example.com",
            user = "deploy",
            identity_file = "~/.ssh/deploy_rsa",
            port = 2222,
            delegate_to = "jump-box"
        })
        
        return true, "SSH configured"
    end
})
```

---

## Infrastructure Testing (infra_test)

### Complete Validation Suite

```lua
task("validate_infrastructure", {
    description = "Validate infrastructure state",
    command = function()
        -- File tests
        infra_test.file_exists({
            path = "/etc/nginx/nginx.conf",
            delegate_to = "web-server"
        })
        
        infra_test.file_mode({
            path = "/etc/nginx/nginx.conf",
            mode = "0644",
            delegate_to = "web-server"
        })
        
        infra_test.file_contains({
            path = "/etc/nginx/nginx.conf",
            pattern = "server_name example.com",
            delegate_to = "web-server"
        })
        
        -- Service tests
        infra_test.service_is_running({
            name = "nginx",
            delegate_to = "web-server"
        })
        
        infra_test.service_is_enabled({
            name = "nginx",
            delegate_to = "web-server"
        })
        
        -- Port tests
        infra_test.port_is_listening({
            port = 80,
            delegate_to = "web-server"
        })
        
        infra_test.port_is_tcp({
            port = 80,
            delegate_to = "web-server"
        })
        
        -- Package tests
        infra_test.package_is_installed({
            name = "nginx",
            delegate_to = "web-server"
        })
        
        -- Process tests
        infra_test.process_is_running({
            pattern = "nginx",
            delegate_to = "web-server"
        })
        
        -- Command tests
        infra_test.command_succeeds({
            cmd = "nginx -t",
            delegate_to = "web-server"
        })
        
        -- Network tests
        infra_test.can_connect({
            host = "database.example.com",
            port = 5432,
            timeout_ms = 5000
        })
        
        return true, "All validation tests passed"
    end
})
```

---

## Goroutines (Parallel Execution)

### Parallel Package Installation

```lua
task("parallel_setup", {
    description = "Setup multiple servers in parallel",
    command = function()
        local servers = {"web-1", "web-2", "web-3"}
        local wg = goroutine.wait_group()
        local results = {}
        
        for _, server in ipairs(servers) do
            wg:add(1)
            goroutine.go(function()
                log.info("Installing packages on " .. server)
                
                pkg.update({ delegate_to = server })
                pkg.install({
                    packages = {"nginx", "certbot"},
                    delegate_to = server
                })
                
                systemd.enable({
                    service = "nginx",
                    delegate_to = server
                })
                
                systemd.start({
                    service = "nginx",
                    delegate_to = server
                })
                
                results[server] = true
                wg:done()
            end)
        end
        
        -- Wait for all goroutines to complete
        wg:wait()
        
        -- Check results
        for server, success in pairs(results) do
            if not success then
                return false, "Setup failed on " .. server
            end
            log.info("✓ " .. server .. " configured successfully")
        end
        
        return true, "All servers configured in parallel"
    end
})
```

### Parallel Cloud Deployment

```lua
task("multi_cloud_deploy", {
    description = "Deploy to multiple clouds simultaneously",
    command = function()
        local wg = goroutine.wait_group()
        local errors = {}
        
        -- Deploy to AWS
        wg:add(1)
        goroutine.go(function()
            local ok, err = aws.s3.sync({
                source = "./build",
                destination = "s3://my-app-aws/",
                delete = true
            })
            if not ok then
                errors["aws"] = err
            end
            wg:done()
        end)
        
        -- Deploy to Azure
        wg:add(1)
        goroutine.go(function()
            local result = azure.exec({
                "storage", "blob", "upload-batch",
                "--destination", "mycontainer",
                "--source", "./build"
            })
            if result.exit_code ~= 0 then
                errors["azure"] = result.stderr
            end
            wg:done()
        end)
        
        -- Deploy to GCP
        wg:add(1)
        goroutine.go(function()
            local result = gcp.exec({
                "storage", "rsync", "-r", "./build",
                "gs://my-app-gcp/"
            })
            if result.exit_code ~= 0 then
                errors["gcp"] = result.stderr
            end
            wg:done()
        end)
        
        -- Wait for all deployments
        wg:wait()
        
        -- Check for errors
        if next(errors) then
            local msg = "Deployment failures:\n"
            for cloud, err in pairs(errors) do
                msg = msg .. "  " .. cloud .. ": " .. err .. "\n"
            end
            return false, msg
        end
        
        return true, "Successfully deployed to all clouds"
    end
})
```

---

## Complete Real-World Example

### Full Application Deployment

```lua
-- Complete deployment workflow with testing and verification

task("prepare_deployment", {
    description = "Prepare for deployment",
    command = function()
        log.info("Preparing deployment environment")
        
        -- Create deployment user
        user.create({
            name = "deploy",
            shell = "/bin/bash",
            delegate_to = "app-server"
        })
        
        -- Setup SSH access
        user.authorized_key({
            user = "deploy",
            key = fs.read_file("~/.ssh/deploy.pub"),
            delegate_to = "app-server"
        })
        
        return true, "Deployment prepared"
    end
})

task("install_dependencies", {
    description = "Install required packages",
    depends_on = {"prepare_deployment"},
    command = function()
        pkg.update({ delegate_to = "app-server" })
        pkg.install({
            packages = {"nginx", "postgresql", "redis", "nodejs"},
            delegate_to = "app-server"
        })
        
        return true, "Dependencies installed"
    end
})

task("deploy_application", {
    description = "Deploy application code",
    depends_on = {"install_dependencies"},
    command = function()
        -- Deploy configuration
        fs.template({
            src = "./templates/nginx.conf.tmpl",
            dest = "/etc/nginx/sites-available/myapp",
            vars = {
                domain = "example.com",
                port = 3000
            },
            delegate_to = "app-server"
        })
        
        -- Deploy application files
        fs.copy({
            src = "./build/",
            dest = "/var/www/myapp/",
            mode = "0755",
            owner = "deploy",
            delegate_to = "app-server"
        })
        
        return true, "Application deployed"
    end
})

task("configure_services", {
    description = "Configure and start services",
    depends_on = {"deploy_application"},
    command = function()
        -- Enable services
        systemd.enable({
            service = "nginx",
            delegate_to = "app-server"
        })
        
        systemd.enable({
            service = "postgresql",
            delegate_to = "app-server"
        })
        
        -- Restart services
        systemd.restart({
            service = "nginx",
            delegate_to = "app-server"
        })
        
        systemd.restart({
            service = "postgresql",
            delegate_to = "app-server"
        })
        
        return true, "Services configured"
    end
})

task("verify_deployment", {
    description = "Verify deployment success",
    depends_on = {"configure_services"},
    command = function()
        log.info("Running verification tests...")
        
        -- File tests
        infra_test.file_exists({
            path = "/var/www/myapp/index.html",
            delegate_to = "app-server"
        })
        
        -- Service tests
        infra_test.service_is_running({
            name = "nginx",
            delegate_to = "app-server"
        })
        
        infra_test.service_is_running({
            name = "postgresql",
            delegate_to = "app-server"
        })
        
        -- Port tests
        infra_test.port_is_listening({
            port = 80,
            delegate_to = "app-server"
        })
        
        -- HTTP test
        infra_test.command_succeeds({
            cmd = "curl -f http://localhost",
            delegate_to = "app-server"
        })
        
        log.info("✓ All verification tests passed")
        return true, "Deployment verified successfully"
    end
})
```

---

## Migration Guide

### Old DSL vs New DSL

**Old (Positional Arguments):**
```lua
local result = pkg.install("nginx", "app-server")
```

**New (Table Parameters):**
```lua
local result = pkg.install({
    packages = {"nginx"},
    delegate_to = "app-server"
})
```

### Benefits of New Syntax

1. **Clarity**: Named parameters make code self-documenting
2. **Flexibility**: Easy to add optional parameters
3. **Consistency**: All modules follow the same pattern
4. **Maintainability**: Easier to understand and modify

---

## Best Practices

1. **Always use table parameters** for better code readability
2. **Use delegate_to** for remote execution instead of manual SSH
3. **Check return values** and handle errors appropriately
4. **Use goroutines** for parallel operations to improve performance
5. **Leverage infra_test** to verify infrastructure state
6. **Keep tasks focused** - one task should do one thing well

---

## See Also

- [Module Index](../modules/index.md)
- [Infrastructure Testing](../modules/infra_test.md)
- [Goroutines](../modules/goroutine.md)
- [Distributed Agents](../distributed-agents.md)
