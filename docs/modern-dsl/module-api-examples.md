# 📚 Module API Examples - Modern DSL

This guide provides comprehensive examples of using Sloth Runner modules with the Modern DSL. All examples follow best practices and demonstrate real-world scenarios.

## 📦 Package Management Module (`pkg`)

The package management module provides cross-platform package installation and management.

### Basic Package Installation

```lua
name = "package-setup"
version = "1.0.0"

group "system_packages" {
    task "update_repositories" {
        module = "pkg",
        action = "update",
        -- Automatically detects the package manager (apt, yum, dnf, etc.)
    }

    task "install_essentials" {
        module = "pkg",
        action = "install",
        packages = {"curl", "git", "vim", "htop"},
        state = "present"  -- Ensure packages are installed
    }

    task "install_development_tools" {
        module = "pkg",
        action = "install",
        packages = {"gcc", "make", "python3-pip", "nodejs", "npm"},
        state = "latest"  -- Ensure latest versions
    }
}
```

### Advanced Package Management

```lua
-- Package management with version control
group "versioned_packages" {
    task "install_specific_versions" {
        module = "pkg",
        action = "install",
        packages = {
            {name = "postgresql", version = "14"},
            {name = "redis", version = "7.0*"},
            {name = "nginx", version = ">=1.20"}
        }
    }

    task "remove_unwanted" {
        module = "pkg",
        action = "remove",
        packages = {"apache2", "mysql-server"},
        state = "absent",
        purge = true  -- Remove configuration files too
    }
}
```

## ⚙️ Systemd Module

Managing system services with the systemd module.

### Service Management

```lua
group "service_configuration" {
    task "configure_nginx" {
        module = "systemd",
        action = "service",
        name = "nginx",
        state = "started",
        enabled = true,
        daemon_reload = true  -- Reload systemd if unit files changed
    }

    task "configure_multiple_services" {
        module = "systemd",
        action = "multi_service",
        services = {
            {name = "postgresql", state = "started", enabled = true},
            {name = "redis", state = "started", enabled = true},
            {name = "memcached", state = "stopped", enabled = false}
        }
    }
}
```

### Custom Service Creation

```lua
group "custom_service" {
    task "create_app_service" {
        module = "fs",
        action = "write",
        path = "/etc/systemd/system/myapp.service",
        content = [[
[Unit]
Description=My Application
After=network.target

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/bin/start.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
        ]],
        mode = "0644"
    }

    task "reload_and_start" {
        module = "systemd",
        action = "daemon_reload"
    }

    task "start_custom_service" {
        module = "systemd",
        action = "service",
        name = "myapp",
        state = "started",
        enabled = true
    }
}
```

## 🐳 Docker Module

Container management with Docker.

### Container Deployment

```lua
group "docker_deployment" {
    task "pull_images" {
        module = "docker",
        action = "pull",
        images = {
            "nginx:latest",
            "postgres:14-alpine",
            "redis:7-alpine"
        }
    }

    task "run_database" {
        module = "docker",
        action = "container",
        name = "app-postgres",
        image = "postgres:14-alpine",
        state = "started",
        restart_policy = "always",
        environment = {
            POSTGRES_DB = "myapp",
            POSTGRES_USER = "appuser",
            POSTGRES_PASSWORD = state.get("db_password")
        },
        volumes = {
            "/data/postgres:/var/lib/postgresql/data"
        },
        ports = {"5432:5432"}
    }

    task "run_application" {
        module = "docker",
        action = "container",
        name = "myapp",
        image = "myapp:latest",
        state = "started",
        restart_policy = "always",
        environment = {
            DATABASE_URL = "postgresql://appuser@app-postgres/myapp",
            REDIS_URL = "redis://app-redis:6379"
        },
        links = {"app-postgres", "app-redis"},
        ports = {"80:3000"},
        healthcheck = {
            test = ["CMD", "curl", "-f", "http://localhost:3000/health"],
            interval = "30s",
            timeout = "10s",
            retries = 3
        }
    }
}
```

### Docker Compose

```lua
group "compose_deployment" {
    task "deploy_stack" {
        module = "docker",
        action = "compose",
        project_name = "myapp",
        compose_file = "./docker-compose.yml",
        state = "present",
        pull = true,
        build = true
    }

    task "scale_services" {
        module = "docker",
        action = "compose_scale",
        project_name = "myapp",
        services = {
            web = 3,
            worker = 2
        }
    }
}
```

## 🏗️ Terraform Module

Infrastructure as Code with Terraform.

### Complete Terraform Workflow

```lua
group "terraform_infrastructure" {
    description = "Deploy AWS infrastructure with Terraform",

    task "init_terraform" {
        module = "terraform",
        action = "init",
        working_dir = "./terraform/environments/production",
        backend_config = {
            bucket = "terraform-state-bucket",
            key = "production/terraform.tfstate",
            region = "us-east-1"
        }
    }

    task "plan_infrastructure" {
        module = "terraform",
        action = "plan",
        working_dir = "./terraform/environments/production",
        var_file = "./terraform/environments/production/terraform.tfvars",
        out = "/tmp/production.tfplan",
        variables = {
            environment = "production",
            instance_count = 3
        }
    }

    task "apply_infrastructure" {
        module = "terraform",
        action = "apply",
        working_dir = "./terraform/environments/production",
        plan_file = "/tmp/production.tfplan",
        auto_approve = true,

        on_success = function(result)
            -- Store outputs in state
            local outputs = terraform.get_outputs({
                working_dir = "./terraform/environments/production"
            })
            state.set("vpc_id", outputs.vpc_id)
            state.set("subnet_ids", outputs.subnet_ids)
            log.info("Infrastructure deployed: VPC " .. outputs.vpc_id)
        end
    }
}
```

## 🔀 Git Module

Version control operations.

### Git Operations

```lua
group "git_operations" {
    task "clone_repository" {
        module = "git",
        action = "clone",
        repository = "https://github.com/myorg/myapp.git",
        dest = "/opt/myapp",
        branch = "main",
        depth = 1  -- Shallow clone for faster operation
    }

    task "update_code" {
        module = "git",
        action = "pull",
        repo_path = "/opt/myapp",
        branch = "main",

        before = function()
            -- Save current commit hash
            local current = git.get_commit({path = "/opt/myapp"})
            state.set("previous_commit", current)
        end,

        after = function()
            local new_commit = git.get_commit({path = "/opt/myapp"})
            log.info("Updated from " .. state.get("previous_commit") .. " to " .. new_commit)
        end
    }

    task "tag_release" {
        module = "git",
        action = "tag",
        repo_path = "/opt/myapp",
        tag = "v" .. state.get("version"),
        message = "Release version " .. state.get("version"),
        push = true
    }
}
```

## 🔗 Stow Module (Dotfiles Management)

Managing dotfiles and configuration symlinks.

### Dotfiles Setup

```lua
group "dotfiles_management" {
    task "clone_dotfiles" {
        module = "git",
        action = "clone",
        repository = "https://github.com/user/dotfiles.git",
        dest = "$HOME/.dotfiles"
    }

    task "stow_configurations" {
        module = "stow",
        action = "stow",
        source = "$HOME/.dotfiles",
        packages = {"vim", "tmux", "zsh", "git"},
        target = "$HOME",

        on_conflict = function(conflict)
            log.warn("Conflict detected: " .. conflict.file)
            -- Backup existing file
            fs.move({
                src = conflict.file,
                dest = conflict.file .. ".backup"
            })
            return "retry"  -- Retry stow operation
        end
    }

    task "unstow_old_configs" {
        module = "stow",
        action = "unstow",
        source = "$HOME/.dotfiles",
        packages = {"old-config"},
        target = "$HOME"
    }
}
```

## 📦 Incus Module (LXC/VM Management)

Container and VM orchestration with Incus.

### Container Creation and Management

```lua
group "incus_containers" {
    task "create_web_container" {
        module = "incus",
        action = "container",
        name = "web-01",
        image = "ubuntu:22.04",
        state = "started",
        profiles = {"default", "web"},
        config = {
            "limits.cpu" = "2",
            "limits.memory" = "2GB",
            "security.nesting" = "true"
        },
        devices = {
            root = {
                type = "disk",
                pool = "default",
                size = "20GB"
            },
            eth0 = {
                type = "nic",
                network = "lxdbr0",
                ["ipv4.address"] = "10.0.0.100"
            }
        }
    }

    task "configure_container" {
        module = "incus",
        action = "exec",
        container = "web-01",
        command = {
            "apt-get", "update", "&&",
            "apt-get", "install", "-y", "nginx"
        }
    }

    task "snapshot_container" {
        module = "incus",
        action = "snapshot",
        container = "web-01",
        name = "after-nginx-install",
        stateful = false
    }
}
```

## 🧪 Infrastructure Testing Module (`infra_test`)

Validate infrastructure state and configuration.

### Comprehensive Testing Suite

```lua
group "infrastructure_validation" {
    description = "Validate deployed infrastructure",

    task "test_web_server" {
        module = "infra_test",
        action = "suite",
        tests = {
            -- File tests
            {type = "file_exists", path = "/etc/nginx/nginx.conf"},
            {type = "file_contains", path = "/etc/nginx/nginx.conf", pattern = "worker_processes auto"},
            {type = "file_mode", path = "/etc/nginx/nginx.conf", mode = "0644"},
            {type = "file_owner", path = "/etc/nginx/nginx.conf", owner = "root", group = "root"},

            -- Service tests
            {type = "service_running", name = "nginx"},
            {type = "service_enabled", name = "nginx"},

            -- Port tests
            {type = "port_listening", port = 80, protocol = "tcp"},
            {type = "port_listening", port = 443, protocol = "tcp"},

            -- Process tests
            {type = "process_running", pattern = "nginx: master process"},

            -- Command tests
            {type = "command_succeeds", command = "nginx -t"},
            {type = "command_output", command = "nginx -v", pattern = "nginx/1.2"},

            -- HTTP tests
            {type = "http_response", url = "http://localhost", status = 200},
            {type = "http_contains", url = "http://localhost", content = "Welcome"},

            -- SSL tests
            {type = "ssl_certificate_valid", host = "example.com", port = 443},
            {type = "ssl_days_remaining", host = "example.com", min_days = 30}
        },

        on_failure = function(test, error)
            notification.send("slack", {
                channel = "#ops-alerts",
                message = "Infrastructure test failed: " .. test.type .. " - " .. error
            })
        end
    }
}
```

## ⚡ Parallel Execution with Goroutines

Execute tasks concurrently for improved performance.

### Parallel Multi-Server Deployment

```lua
group "parallel_deployment" {
    task "deploy_to_all_servers" {
        module = "goroutine",
        action = "parallel",

        execute = function()
            local servers = {"web-01", "web-02", "web-03", "app-01", "app-02"}
            local results = {}
            local wg = goroutine.WaitGroup()

            for _, server in ipairs(servers) do
                wg:Add(1)
                goroutine.Go(function()
                    log.info("Deploying to " .. server)

                    -- Copy application files
                    fs.sync({
                        source = "./dist/",
                        dest = server .. ":/opt/app/",
                        delete = true
                    })

                    -- Restart service
                    systemd.restart({
                        service = "myapp",
                        host = server
                    })

                    -- Verify deployment
                    local ok = infra_test.http_response({
                        url = "http://" .. server .. ":3000/health",
                        status = 200
                    })

                    results[server] = ok
                    wg:Done()
                end)
            end

            wg:Wait()

            -- Check results
            local failed = {}
            for server, ok in pairs(results) do
                if not ok then
                    table.insert(failed, server)
                end
            end

            if #failed > 0 then
                return false, "Deployment failed on: " .. table.concat(failed, ", ")
            end

            return true, "Successfully deployed to all servers"
        end
    }
}
```

### Parallel Cloud Provisioning

```lua
group "multi_cloud_provisioning" {
    task "provision_all_clouds" {
        module = "goroutine",
        action = "parallel",

        execute = function()
            local tasks = {
                aws = function()
                    terraform.apply({
                        working_dir = "./terraform/aws",
                        auto_approve = true
                    })
                end,

                azure = function()
                    terraform.apply({
                        working_dir = "./terraform/azure",
                        auto_approve = true
                    })
                end,

                gcp = function()
                    terraform.apply({
                        working_dir = "./terraform/gcp",
                        auto_approve = true
                    })
                end
            }

            local results = goroutine.RunParallel(tasks)

            for cloud, result in pairs(results) do
                if not result.success then
                    log.error(cloud .. " provisioning failed: " .. result.error)
                    return false, "Multi-cloud provisioning failed"
                end
            end

            return true, "All clouds provisioned successfully"
        end
    }
}
```

## 💾 State Management Module

Persistent state management across task executions.

### State Operations

```lua
group "state_management" {
    task "save_deployment_info" {
        module = "state",
        action = "batch",

        operations = function()
            -- Set multiple values
            state.set("deployment.version", "2.1.0")
            state.set("deployment.timestamp", os.time())
            state.set("deployment.commit", git.get_commit({path = "."}))

            -- Set complex objects
            state.set("deployment.servers", {
                web = {"web-01", "web-02"},
                app = {"app-01", "app-02"},
                db = {"db-01"}
            })

            -- Increment counter
            local deploy_count = state.get("deployment.count", 0)
            state.set("deployment.count", deploy_count + 1)

            -- Conditional state
            if state.get("environment") == "production" then
                state.set("deployment.approval_required", true)
            end
        end
    }

    task "load_previous_state" {
        module = "state",
        action = "load",

        execute = function()
            local last_version = state.get("deployment.version", "unknown")
            local last_timestamp = state.get("deployment.timestamp", 0)

            if last_timestamp > 0 then
                local last_date = os.date("%Y-%m-%d %H:%M:%S", last_timestamp)
                log.info("Last deployment: " .. last_version .. " at " .. last_date)
            end

            -- Get nested values
            local servers = state.get("deployment.servers")
            if servers then
                for role, list in pairs(servers) do
                    log.info(role .. " servers: " .. table.concat(list, ", "))
                end
            end
        end
    }
}
```

## 🔔 Notifications Module

Send notifications to various channels.

### Multi-Channel Notifications

```lua
group "notifications" {
    task "send_deployment_notifications" {
        module = "notification",
        action = "multi_send",

        execute = function()
            local version = state.get("deployment.version")
            local environment = state.get("environment")

            -- Slack notification
            notification.send("slack", {
                webhook_url = os.getenv("SLACK_WEBHOOK"),
                channel = "#deployments",
                username = "Deployment Bot",
                icon_emoji = ":rocket:",
                message = "Deployment started",
                attachments = {
                    {
                        color = "good",
                        title = "Deployment Details",
                        fields = {
                            {title = "Version", value = version, short = true},
                            {title = "Environment", value = environment, short = true},
                            {title = "Triggered by", value = os.getenv("USER"), short = true}
                        }
                    }
                }
            })

            -- Email notification
            notification.send("email", {
                to = {"ops-team@example.com"},
                subject = "Deployment: " .. version,
                body = "Deployment of version " .. version .. " to " .. environment .. " has started.",
                smtp_server = "smtp.example.com",
                smtp_port = 587,
                smtp_user = "notifications@example.com",
                smtp_password = os.getenv("SMTP_PASSWORD")
            })

            -- Discord notification
            notification.send("discord", {
                webhook_url = os.getenv("DISCORD_WEBHOOK"),
                content = "**Deployment Alert**",
                embeds = {
                    {
                        title = "New Deployment",
                        description = "Version " .. version .. " is being deployed",
                        color = 5763719,  -- Green
                        fields = {
                            {name = "Environment", value = environment},
                            {name = "Status", value = "In Progress"}
                        }
                    }
                }
            })
        end
    }
}
```

## 🌐 Network Module

Network operations and testing.

### Network Configuration and Testing

```lua
group "network_setup" {
    task "configure_firewall" {
        module = "net",
        action = "firewall",
        rules = {
            {action = "allow", port = 22, protocol = "tcp", source = "10.0.0.0/8"},
            {action = "allow", port = 80, protocol = "tcp"},
            {action = "allow", port = 443, protocol = "tcp"},
            {action = "deny", port = 3306, protocol = "tcp", source = "0.0.0.0/0"}
        }
    }

    task "test_connectivity" {
        module = "net",
        action = "test",
        tests = {
            {type = "ping", host = "8.8.8.8", count = 3},
            {type = "tcp", host = "github.com", port = 443},
            {type = "http", url = "https://api.github.com", status = 200},
            {type = "dns", hostname = "example.com", expected = "93.184.216.34"}
        }
    }

    task "download_artifacts" {
        module = "net",
        action = "download",
        downloads = {
            {
                url = "https://releases.example.com/app-v2.tar.gz",
                dest = "/tmp/app.tar.gz",
                checksum = "sha256:abcdef123456...",
                headers = {
                    ["Authorization"] = "Bearer " .. os.getenv("GITHUB_TOKEN")
                }
            }
        }
    }
}
```

## 📊 Complete Real-World Example

A complete deployment pipeline using multiple modules:

```lua
name = "production-deployment"
version = "1.0.0"
description = "Complete production deployment pipeline"

-- Configuration
local config = {
    app_name = "myapp",
    environment = "production",
    version = os.getenv("VERSION") or "latest",
    servers = {
        web = {"web-01", "web-02", "web-03"},
        app = {"app-01", "app-02"},
        db = {"db-01", "db-02"}
    }
}

-- Pre-deployment checks
group "pre_deployment" {
    task "validate_environment" {
        module = "infra_test",
        action = "suite",
        tests = {
            -- Check all servers are accessible
            {type = "ping", hosts = config.servers.web},
            {type = "ping", hosts = config.servers.app},
            {type = "ping", hosts = config.servers.db},

            -- Check required services
            {type = "service_running", name = "docker", hosts = config.servers.app},
            {type = "port_listening", port = 5432, hosts = config.servers.db}
        }
    }

    task "backup_database" {
        module = "exec",
        action = "command",
        command = "pg_dump -h db-01 -U postgres myapp > /backups/myapp-$(date +%Y%m%d-%H%M%S).sql",
        host = "backup-server"
    }
}

-- Build and prepare
group "build" {
    task "build_application" {
        module = "docker",
        action = "build",
        context = "./",
        dockerfile = "./Dockerfile",
        tag = config.app_name .. ":" .. config.version,
        build_args = {
            VERSION = config.version,
            BUILD_DATE = os.date("%Y-%m-%d")
        }
    }

    task "push_to_registry" {
        module = "docker",
        action = "push",
        image = config.app_name .. ":" .. config.version,
        registry = "registry.example.com"
    }
}

-- Deploy to servers
group "deployment" {
    task "deploy_web_servers" {
        module = "goroutine",
        action = "parallel",

        execute = function()
            local tasks = {}

            for _, server in ipairs(config.servers.web) do
                tasks[server] = function()
                    -- Update nginx config
                    fs.template({
                        source = "./templates/nginx.conf.j2",
                        dest = "/etc/nginx/sites-available/" .. config.app_name,
                        variables = {
                            app_name = config.app_name,
                            upstream_servers = config.servers.app
                        },
                        host = server
                    })

                    -- Reload nginx
                    systemd.reload({
                        service = "nginx",
                        host = server
                    })
                end
            end

            return goroutine.RunParallel(tasks)
        end
    }

    task "deploy_app_servers" {
        module = "goroutine",
        action = "parallel",

        execute = function()
            local tasks = {}

            for _, server in ipairs(config.servers.app) do
                tasks[server] = function()
                    -- Pull new image
                    docker.pull({
                        image = config.app_name .. ":" .. config.version,
                        host = server
                    })

                    -- Stop old container
                    docker.stop({
                        container = config.app_name,
                        host = server
                    })

                    -- Start new container
                    docker.run({
                        name = config.app_name,
                        image = config.app_name .. ":" .. config.version,
                        ports = {"3000:3000"},
                        environment = {
                            NODE_ENV = config.environment,
                            DATABASE_URL = "postgresql://db-01:5432/myapp"
                        },
                        restart = "always",
                        host = server
                    })
                end
            end

            return goroutine.RunParallel(tasks)
        end
    }
}

-- Post-deployment
group "post_deployment" {
    task "verify_deployment" {
        module = "infra_test",
        action = "suite",
        tests = {
            -- Check application health
            {type = "http_response", url = "https://example.com/health", status = 200},
            {type = "http_contains", url = "https://example.com/version", content = config.version},

            -- Check all app servers
            {type = "docker_container_running", name = config.app_name, hosts = config.servers.app}
        }
    }

    task "send_notifications" {
        module = "notification",
        action = "send",
        channel = "slack",
        message = "Deployment completed successfully!",
        details = {
            Version = config.version,
            Environment = config.environment,
            Servers = table.concat(config.servers.app, ", ")
        }
    }

    task "update_monitoring" {
        module = "exec",
        action = "command",
        command = "curl -X POST https://monitoring.example.com/api/deployments",
        data = {
            application = config.app_name,
            version = config.version,
            timestamp = os.time()
        }
    }
}

-- Rollback group (if needed)
group "rollback" {
    enabled = false,  -- Enable manually if rollback needed

    task "rollback_to_previous" {
        module = "state",
        action = "execute",

        execute = function()
            local previous_version = state.get("previous_version")
            if not previous_version then
                return false, "No previous version found"
            end

            log.warning("Rolling back to version: " .. previous_version)

            -- Re-deploy previous version
            for _, server in ipairs(config.servers.app) do
                docker.stop({
                    container = config.app_name,
                    host = server
                })

                docker.run({
                    name = config.app_name,
                    image = config.app_name .. ":" .. previous_version,
                    ports = {"3000:3000"},
                    restart = "always",
                    host = server
                })
            end

            return true, "Rolled back to " .. previous_version
        end
    }
}
```

## 🎯 Best Practices

1. **Use Functions for Reusability**: Create helper functions for common patterns
2. **Leverage Parallel Execution**: Use goroutines for operations that can run concurrently
3. **Implement Proper Error Handling**: Always check return values and handle failures
4. **Use State Management**: Track deployment state for rollbacks and auditing
5. **Test Infrastructure**: Use `infra_test` module to verify deployments
6. **Send Notifications**: Keep team informed about deployment status
7. **Version Everything**: Tag and version your deployments
8. **Document Complex Logic**: Add comments to explain non-obvious code

## 📖 Next Steps

- [Best Practices Guide](best-practices.md)
- [Reference Guide](reference-guide.md)
- [Migration from YAML](migration-guide.md)
- [Custom Module Development](custom-modules.md)