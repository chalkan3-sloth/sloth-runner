# 🔧 Complete Modules Reference

## Overview

Sloth Runner has 40+ built-in modules that provide functionality ranging from basic system operations to complex cloud provider integrations. This documentation covers **all** available modules with practical examples.

---

## 📦 Package Management

### Module `pkg` - Package Management

Manages system packages using apt, yum, dnf, pacman, brew, etc.

**Functions:**

#### `pkg.install(name, options)`

Installs one or more packages.

```lua
-- Install a package
pkg.install("nginx")

-- Install multiple packages
pkg.install({"nginx", "postgresql", "redis"})

-- With options
pkg.install("nginx", {
    update_cache = true,  -- Update cache before installing
    state = "present"     -- present (default) or latest
})

-- Install specific version (apt)
pkg.install("nginx=1.18.0-0ubuntu1")
```

#### `pkg.remove(name, options)`

Removes one or more packages.

```lua
-- Remove a package
pkg.remove("nginx")

-- Remove multiple
pkg.remove({"nginx", "apache2"})

-- Remove with purge (apt)
pkg.remove("nginx", { purge = true })
```

#### `pkg.update()`

Updates package cache.

```lua
-- Update cache (apt update, yum update, etc)
pkg.update()
```

#### `pkg.upgrade(name)`

Upgrades installed packages.

```lua
-- Upgrade all packages
pkg.upgrade()

-- Upgrade specific package
pkg.upgrade("nginx")
```

**Complete example:**

```yaml
tasks:
  - name: Setup web server
    exec:
      script: |
        -- Update cache
        pkg.update()

        -- Install necessary packages
        pkg.install({
          "nginx",
          "certbot",
          "python3-certbot-nginx"
        }, { state = "latest" })

        -- Remove old web server
        pkg.remove("apache2", { purge = true })
```

---

### Module `user` - User Management

Manages system users and groups.

**Functions:**

#### `user.create(name, options)`

Creates a user.

```lua
-- Create simple user
user.create("deploy")

-- With full options
user.create("deploy", {
    uid = 1001,
    gid = 1001,
    groups = {"sudo", "docker"},
    shell = "/bin/bash",
    home = "/home/deploy",
    create_home = true,
    system = false,
    comment = "Deploy user"
})
```

#### `user.remove(name, options)`

Removes a user.

```lua
-- Remove user
user.remove("olduser")

-- Remove and delete home
user.remove("olduser", { remove_home = true })
```

#### `user.exists(name)`

Checks if user exists.

```lua
if user.exists("deploy") then
    log.info("User deploy exists")
else
    user.create("deploy")
end
```

#### `group.create(name, options)`

Creates a group.

```lua
group.create("developers")
group.create("developers", { gid = 2000 })
```

---

## 📁 File Operations

### Module `file` - File Operations

Manages files and directories.

**Functions:**

#### `file.copy(source, destination, options)`

Copies files or directories.

```lua
-- Copy file
file.copy("/src/app.conf", "/etc/app/app.conf")

-- With options
file.copy("/src/app.conf", "/etc/app/app.conf", {
    owner = "root",
    group = "root",
    mode = "0644",
    backup = true  -- Backup if destination exists
})

-- Copy directory recursively
file.copy("/src/configs/", "/etc/myapp/", {
    recursive = true
})
```

#### `file.create(path, options)`

Creates a file.

```lua
-- Create empty file
file.create("/var/log/myapp.log")

-- With content and permissions
file.create("/etc/myapp/config.yaml", {
    content = [[
        server:
          host: 0.0.0.0
          port: 8080
    ]],
    owner = "myapp",
    group = "myapp",
    mode = "0640"
})
```

#### `file.remove(path, options)`

Removes files or directories.

```lua
-- Remove file
file.remove("/tmp/cache.dat")

-- Remove directory recursively
file.remove("/var/cache/oldapp", { recursive = true })

-- Remove with force
file.remove("/var/log/*.log", { force = true })
```

#### `file.exists(path)`

Checks if file/directory exists.

```lua
if file.exists("/etc/nginx/nginx.conf") then
    log.info("Nginx config found")
end
```

#### `file.chmod(path, mode)`

Changes permissions.

```lua
file.chmod("/usr/local/bin/myapp", "0755")
file.chmod("/etc/ssl/private/key.pem", "0600")
```

#### `file.chown(path, owner, group)`

Changes owner and group.

```lua
file.chown("/var/www/html", "www-data", "www-data")
```

#### `file.read(path)`

Reads file content.

```lua
local content = file.read("/etc/hostname")
log.info("Hostname: " .. content)
```

#### `file.write(path, content, options)`

Writes content to file.

```lua
file.write("/etc/motd", "Welcome to Production Server\n")

-- With append
file.write("/var/log/app.log", "Log entry\n", {
    append = true
})
```

---

### Module `template` - Templates

Processes templates with variables.

```lua
-- Jinja2/Go template
template.render("/templates/nginx.conf.j2", "/etc/nginx/nginx.conf", {
    server_name = "example.com",
    port = 80,
    root = "/var/www/html"
})
```

---

### Module `stow` - Dotfiles Management

Manages dotfiles using GNU Stow.

```lua
-- Stow dotfiles
stow.link("~/.dotfiles/vim", "~")
stow.link("~/.dotfiles/zsh", "~")

-- Unstow
stow.unlink("~/.dotfiles/vim", "~")

-- Restow (unstow + stow)
stow.restow("~/.dotfiles/vim", "~")
```

---

## 🐚 Command Execution

### Module `exec` - Command Execution

Executes system commands.

**Functions:**

#### `exec.command(command, options)`

Executes a command.

```lua
-- Simple command
local result = exec.command("ls -la /tmp")

-- With options
local result = exec.command("systemctl restart nginx", {
    user = "root",
    cwd = "/etc/nginx",
    env = {
        PATH = "/usr/local/bin:/usr/bin:/bin"
    },
    timeout = 30  -- seconds
})

-- Check result
if result.exit_code == 0 then
    log.info("Success: " .. result.stdout)
else
    log.error("Failed: " .. result.stderr)
end
```

#### `exec.shell(script)`

Executes shell script.

```lua
exec.shell([[
    #!/bin/bash
    set -e

    apt update
    apt install -y nginx
    systemctl enable nginx
    systemctl start nginx
]])
```

#### `exec.script(path, options)`

Executes script from file.

```lua
exec.script("/scripts/deploy.sh")

exec.script("/scripts/deploy.sh", {
    interpreter = "/bin/bash",
    args = {"production", "v1.2.3"}
})
```

---

### Module `goroutine` - Parallel Execution

Executes tasks in parallel using goroutines.

```lua
goroutine.parallel({
    function()
        pkg.install("nginx")
    end,
    function()
        pkg.install("postgresql")
    end,
    function()
        pkg.install("redis")
    end
})

-- With concurrency limit
goroutine.parallel({
    tasks = {
        function() exec.command("task1") end,
        function() exec.command("task2") end,
        function() exec.command("task3") end
    },
    max_concurrent = 2  -- Maximum 2 at a time
})
```

---

## 🐳 Containers and Virtualization

### Module `docker` - Docker

Manages Docker containers, images, and networks.

**Functions:**

#### `docker.container_run(image, options)`

Runs a container.

```lua
docker.container_run("nginx:latest", {
    name = "web-server",
    ports = {"80:80", "443:443"},
    volumes = {"/var/www:/usr/share/nginx/html:ro"},
    env = {
        NGINX_HOST = "example.com",
        NGINX_PORT = "80"
    },
    restart = "unless-stopped",
    detach = true
})
```

#### `docker.container_stop(name)`

Stops a container.

```lua
docker.container_stop("web-server")
```

#### `docker.container_remove(name, options)`

Removes a container.

```lua
docker.container_remove("web-server")
docker.container_remove("web-server", { force = true, volumes = true })
```

#### `docker.image_pull(image, options)`

Pulls an image.

```lua
docker.image_pull("nginx:latest")
docker.image_pull("myregistry.com/myapp:v1.2.3", {
    auth = {
        username = "user",
        password = "pass"
    }
})
```

#### `docker.image_build(context, options)`

Builds an image.

```lua
docker.image_build(".", {
    tag = "myapp:latest",
    dockerfile = "Dockerfile",
    build_args = {
        VERSION = "1.2.3",
        ENV = "production"
    }
})
```

#### `docker.network_create(name, options)`

Creates a network.

```lua
docker.network_create("app-network", {
    driver = "bridge",
    subnet = "172.20.0.0/16"
})
```

#### `docker.compose_up(file, options)`

Runs docker-compose.

```lua
docker.compose_up("docker-compose.yml", {
    project_name = "myapp",
    detach = true,
    build = true
})
```

**Complete example:**

```yaml
tasks:
  - name: Deploy application with Docker
    exec:
      script: |
        -- Create network
        docker.network_create("app-net")

        -- Database
        docker.container_run("postgres:14", {
            name = "app-db",
            network = "app-net",
            env = {
                POSTGRES_DB = "myapp",
                POSTGRES_USER = "myapp",
                POSTGRES_PASSWORD = "secret"
            },
            volumes = {"pgdata:/var/lib/postgresql/data"}
        })

        -- Application
        docker.container_run("myapp:latest", {
            name = "app",
            network = "app-net",
            ports = {"3000:3000"},
            env = {
                DATABASE_URL = "postgres://myapp:secret@app-db:5432/myapp"
            },
            depends_on = {"app-db"}
        })
```

---

### Module `incus` - Incus/LXC Containers

Manages Incus (LXC) containers and VMs.

**Functions:**

#### `incus.launch(image, name, options)`

Creates and starts a container/VM.

```lua
-- Ubuntu container
incus.launch("ubuntu:22.04", "web-01", {
    type = "container",  -- or "virtual-machine"
    config = {
        ["limits.cpu"] = "2",
        ["limits.memory"] = "2GB"
    },
    devices = {
        eth0 = {
            type = "nic",
            network = "lxdbr0"
        }
    }
})

-- VM with cloud-init
incus.launch("ubuntu:22.04", "vm-01", {
    type = "virtual-machine",
    config = {
        ["limits.cpu"] = "4",
        ["limits.memory"] = "4GB",
        ["cloud-init.user-data"] = [[
#cloud-init
packages:
  - nginx
  - postgresql
        ]]
    }
})
```

#### `incus.exec(name, command)`

Executes command in container.

```lua
incus.exec("web-01", "apt update && apt install -y nginx")
```

#### `incus.file_push(source, name, destination)`

Pushes file to container.

```lua
incus.file_push("/local/app.conf", "web-01", "/etc/app/app.conf")
```

#### `incus.file_pull(name, source, destination)`

Pulls file from container.

```lua
incus.file_pull("web-01", "/var/log/app.log", "/backup/app.log")
```

#### `incus.stop(name)`

Stops a container.

```lua
incus.stop("web-01")
incus.stop("web-01", { force = true })
```

#### `incus.delete(name)`

Removes a container.

```lua
incus.delete("web-01")
incus.delete("web-01", { force = true })
```

---

## ☁️ Cloud Providers

### Module `aws` - Amazon Web Services

Manages AWS resources (EC2, S3, RDS, etc).

**Functions:**

#### `aws.ec2_instance_create(options)`

Creates EC2 instance.

```lua
aws.ec2_instance_create({
    image_id = "ami-0c55b159cbfafe1f0",
    instance_type = "t3.medium",
    key_name = "my-key",
    security_groups = {"web-sg"},
    subnet_id = "subnet-12345",
    tags = {
        Name = "web-server-01",
        Environment = "production"
    },
    user_data = [[
#!/bin/bash
apt update
apt install -y nginx
    ]]
})
```

#### `aws.s3_bucket_create(name, options)`

Creates S3 bucket.

```lua
aws.s3_bucket_create("my-app-backups", {
    region = "us-east-1",
    acl = "private",
    versioning = true,
    encryption = "AES256"
})
```

#### `aws.s3_upload(file, bucket, key)`

Uploads to S3.

```lua
aws.s3_upload("/backup/db.sql.gz", "my-backups", "db/2024/backup.sql.gz")
```

#### `aws.rds_instance_create(options)`

Creates RDS instance.

```lua
aws.rds_instance_create({
    identifier = "myapp-db",
    engine = "postgres",
    engine_version = "14.7",
    instance_class = "db.t3.medium",
    allocated_storage = 100,
    master_username = "admin",
    master_password = "SecurePass123!",
    vpc_security_groups = {"sg-12345"}
})
```

---

### Module `azure` - Microsoft Azure

Manages Azure resources.

```lua
-- Create VM
azure.vm_create({
    name = "web-vm-01",
    resource_group = "production",
    location = "eastus",
    size = "Standard_D2s_v3",
    image = "Ubuntu2204",
    admin_username = "azureuser",
    ssh_key = file.read("~/.ssh/id_rsa.pub")
})

-- Create Storage Account
azure.storage_account_create({
    name = "myappstorage",
    resource_group = "production",
    location = "eastus",
    sku = "Standard_LRS"
})
```

---

### Module `gcp` - Google Cloud Platform

Manages GCP resources.

```lua
-- Create Compute Engine instance
gcp.compute_instance_create({
    name = "web-instance-01",
    zone = "us-central1-a",
    machine_type = "e2-medium",
    image_project = "ubuntu-os-cloud",
    image_family = "ubuntu-2204-lts",
    tags = {"http-server", "https-server"}
})

-- Create Cloud Storage bucket
gcp.storage_bucket_create("my-app-data", {
    location = "US",
    storage_class = "STANDARD"
})
```

---

### Module `digitalocean` - DigitalOcean

Manages DigitalOcean resources.

```lua
-- Create Droplet
digitalocean.droplet_create({
    name = "web-01",
    region = "nyc3",
    size = "s-2vcpu-4gb",
    image = "ubuntu-22-04-x64",
    ssh_keys = [123456],
    backups = true,
    monitoring = true
})

-- Create Load Balancer
digitalocean.loadbalancer_create({
    name = "web-lb",
    region = "nyc3",
    forwarding_rules = {
        {
            entry_protocol = "https",
            entry_port = 443,
            target_protocol = "http",
            target_port = 80,
            tls_passthrough = false
        }
    },
    droplet_ids = {123456, 789012}
})
```

---

## 🏗️ Infrastructure as Code

### Module `terraform` - Terraform

Manages infrastructure with Terraform.

**Functions:**

#### `terraform.init(dir, options)`

Initializes Terraform.

```lua
terraform.init("/infra/terraform", {
    backend_config = {
        bucket = "my-tf-state",
        key = "prod/terraform.tfstate",
        region = "us-east-1"
    }
})
```

#### `terraform.plan(dir, options)`

Runs plan.

```lua
local plan = terraform.plan("/infra/terraform", {
    var_file = "prod.tfvars",
    vars = {
        environment = "production",
        region = "us-east-1"
    },
    out = "tfplan"
})
```

#### `terraform.apply(dir, options)`

Applies changes.

```lua
terraform.apply("/infra/terraform", {
    plan_file = "tfplan",
    auto_approve = true
})
```

#### `terraform.destroy(dir, options)`

Destroys resources.

```lua
terraform.destroy("/infra/terraform", {
    var_file = "prod.tfvars",
    auto_approve = false  -- Ask for confirmation
})
```

**Complete example:**

```yaml
tasks:
  - name: Deploy infrastructure
    exec:
      script: |
        local tf_dir = "/infra/terraform"

        -- Initialize
        terraform.init(tf_dir)

        -- Plan
        local plan = terraform.plan(tf_dir, {
            var_file = "production.tfvars"
        })

        -- Apply if plan is ok
        if plan.changes > 0 then
            terraform.apply(tf_dir, {
                auto_approve = true
            })
        end
```

---

### Module `pulumi` - Pulumi

Manages infrastructure with Pulumi.

```lua
-- Initialize stack
pulumi.stack_init("production")

-- Configure
pulumi.config_set("aws:region", "us-east-1")

-- Up
pulumi.up({
    stack = "production",
    yes = true  -- Auto-approve
})

-- Destroy
pulumi.destroy({
    stack = "production",
    yes = false
})
```

---

## 🔐 Git and Version Control

### Module `git` - Git

Git repository operations.

**Functions:**

#### `git.clone(url, destination, options)`

Clones a repository.

```lua
git.clone("https://github.com/user/repo.git", "/opt/app")

-- With options
git.clone("https://github.com/user/repo.git", "/opt/app", {
    branch = "main",
    depth = 1,  -- Shallow clone
    auth = {
        username = "user",
        password = "token"
    }
})
```

#### `git.pull(dir, options)`

Updates repository.

```lua
git.pull("/opt/app")
git.pull("/opt/app", { branch = "develop" })
```

#### `git.checkout(dir, ref)`

Checks out branch/tag.

```lua
git.checkout("/opt/app", "v1.2.3")
git.checkout("/opt/app", "develop")
```

#### `git.commit(dir, message, options)`

Creates commit.

```lua
git.commit("/opt/app", "Update config files", {
    author = "Deploy Bot <bot@example.com>",
    add_all = true
})
```

#### `git.push(dir, options)`

Pushes to remote.

```lua
git.push("/opt/app")
git.push("/opt/app", {
    remote = "origin",
    branch = "main",
    force = false
})
```

---

### Module `gitops` - GitOps

Implements GitOps patterns.

```lua
-- Sync from Git
gitops.sync({
    repo = "https://github.com/org/k8s-manifests.git",
    branch = "main",
    path = "production/",
    destination = "/opt/k8s-manifests"
})

-- Apply manifests
gitops.apply({
    source = "/opt/k8s-manifests",
    namespace = "production"
})
```

---

## 🌐 Network and SSH

### Module `net` - Networking

Network operations.

```lua
-- Check if host is online
if net.ping("example.com") then
    log.info("Host is up")
end

-- Port scan
local open = net.port_open("example.com", 443)

-- HTTP request
local response = net.http_get("https://api.example.com/status")
if response.status == 200 then
    log.info(response.body)
end

-- Download file
net.download("https://example.com/file.tar.gz", "/tmp/file.tar.gz")
```

---

### Module `ssh` - SSH

Executes commands via SSH.

```lua
-- Connect and execute
ssh.exec("user@192.168.1.100", "ls -la /opt", {
    key = "~/.ssh/id_rsa",
    port = 22
})

-- Upload file
ssh.upload("user@192.168.1.100", "/local/file.txt", "/remote/file.txt")

-- Download file
ssh.download("user@192.168.1.100", "/remote/log.txt", "/local/log.txt")
```

---

## ⚙️ Services and Systemd

### Module `systemd` - Systemd

Manages systemd services.

**Functions:**

#### `systemd.service_start(name)`

Starts a service.

```lua
systemd.service_start("nginx")
```

#### `systemd.service_stop(name)`

Stops a service.

```lua
systemd.service_stop("nginx")
```

#### `systemd.service_restart(name)`

Restarts a service.

```lua
systemd.service_restart("nginx")
```

#### `systemd.service_enable(name)`

Enables service on boot.

```lua
systemd.service_enable("nginx")
```

#### `systemd.service_disable(name)`

Disables service on boot.

```lua
systemd.service_disable("apache2")
```

#### `systemd.service_status(name)`

Checks status.

```lua
local status = systemd.service_status("nginx")
if status.active then
    log.info("Nginx is running")
end
```

#### `systemd.unit_reload()`

Reloads systemd units.

```lua
systemd.unit_reload()
```

**Complete example:**

```yaml
tasks:
  - name: Deploy and configure nginx
    exec:
      script: |
        -- Install
        pkg.install("nginx")

        -- Configure
        file.copy("/deploy/nginx.conf", "/etc/nginx/nginx.conf")

        -- Enable and start
        systemd.service_enable("nginx")
        systemd.service_start("nginx")

        -- Verify
        local status = systemd.service_status("nginx")
        if not status.active then
            error("Nginx failed to start")
        end
```

---

## 📊 Metrics and Monitoring

### Module `metrics` - Metrics

Collects and sends metrics.

```lua
-- Counter
metrics.counter("requests_total", 1, {
    method = "GET",
    status = "200"
})

-- Gauge
metrics.gauge("memory_usage_bytes", 1024*1024*512)

-- Histogram
metrics.histogram("request_duration_seconds", 0.234)

-- Custom metric
metrics.custom("app_users_active", 42, {
    type = "gauge",
    labels = {
        region = "us-east-1"
    }
})
```

---

### Module `log` - Logging

Advanced logging system.

```lua
-- Log levels
log.debug("Debug message")
log.info("Info message")
log.warn("Warning message")
log.error("Error message")

-- With structured fields
log.info("User login", {
    user_id = 123,
    ip = "192.168.1.100",
    timestamp = os.time()
})

-- Error with stack trace
log.error("Failed to connect", {
    error = err,
    component = "database"
})
```

---

## 🔔 Notifications

### Module `notifications` - Notifications

Sends notifications to various services.

**Functions:**

#### `notifications.slack(webhook, message, options)`

Sends to Slack.

```lua
notifications.slack(
    "https://hooks.slack.com/services/XXX/YYY/ZZZ",
    "Deploy completed successfully! :rocket:",
    {
        channel = "#deployments",
        username = "Sloth Runner",
        icon_emoji = ":sloth:"
    }
)
```

#### `notifications.email(options)`

Sends email.

```lua
notifications.email({
    from = "noreply@example.com",
    to = "admin@example.com",
    subject = "Deploy Status",
    body = "Production deploy completed",
    smtp_host = "smtp.gmail.com",
    smtp_port = 587,
    smtp_user = "user@gmail.com",
    smtp_pass = "password"
})
```

#### `notifications.discord(webhook, message)`

Sends to Discord.

```lua
notifications.discord(
    "https://discord.com/api/webhooks/XXX/YYY",
    "Deploy completed! :tada:"
)
```

#### `notifications.telegram(token, chat_id, message)`

Sends to Telegram.

```lua
notifications.telegram(
    "bot123456:ABC-DEF",
    "123456789",
    "Deploy finished successfully"
)
```

---

## 🧪 Testing and Validation

### Module `infra_test` - Infrastructure Testing

Tests and validates infrastructure.

```lua
-- Test port
infra_test.port("example.com", 443, {
    timeout = 5,
    should_be_open = true
})

-- Test HTTP
infra_test.http("https://example.com", {
    status_code = 200,
    contains = "Welcome",
    timeout = 10
})

-- Test command
infra_test.command("systemctl is-active nginx", {
    exit_code = 0,
    stdout_contains = "active"
})

-- Test file
infra_test.file("/etc/nginx/nginx.conf", {
    exists = true,
    mode = "0644",
    owner = "root"
})
```

---

## 📡 Facts - System Information

### Module `facts` - System Facts

Collects system information.

```lua
-- Get all facts
local facts = facts.gather()

-- Access facts
log.info("OS: " .. facts.os.name)
log.info("Kernel: " .. facts.kernel.version)
log.info("CPU Cores: " .. facts.cpu.cores)
log.info("Memory: " .. facts.memory.total)
log.info("Hostname: " .. facts.hostname)

-- Specific facts
local cpu = facts.cpu()
local mem = facts.memory()
local disk = facts.disk()
local network = facts.network()
```

---

## 🔄 State and Persistence

### Module `state` - State Management

Manages state between executions.

```lua
-- Save state
state.set("last_deploy_version", "v1.2.3")
state.set("last_deploy_time", os.time())

-- Get state
local last_version = state.get("last_deploy_version")
if last_version == nil then
    log.info("First deploy")
end

-- Check if changed
if state.changed("app_config_hash", new_hash) then
    log.info("Config changed, restarting app")
    systemd.service_restart("myapp")
end

-- Clear state
state.clear("temporary_data")
```

---

## 🐍 Python Integration

### Module `python` - Python

Executes Python code.

```lua
-- Run Python script
python.run([[
import requests
import json

response = requests.get('https://api.github.com/repos/user/repo')
data = response.json()
print(f"Stars: {data['stargazers_count']}")
]])

-- Run Python file
python.run_file("/scripts/deploy.py", {
    args = {"production", "v1.2.3"},
    venv = "/opt/venv"
})

-- Install packages
python.pip_install({"requests", "boto3"})
```

---

## 🧂 Configuration Management

### Module `salt` - SaltStack

SaltStack integration.

```lua
-- Apply Salt state
salt.state_apply("webserver", {
    pillar = {
        nginx_port = 80,
        domain = "example.com"
    }
})

-- Run Salt command
salt.cmd_run("service.restart", {"nginx"})
```

---

## 📊 Data Processing

### Module `data` - Data Processing

Processes and transforms data.

```lua
-- Parse JSON
local json_data = data.json_parse('{"name": "value"}')

-- Generate JSON
local json_str = data.json_encode({
    name = "app",
    version = "1.0"
})

-- Parse YAML
local yaml_data = data.yaml_parse([[
name: myapp
version: 1.0
]])

-- Parse TOML
local toml_data = data.toml_parse([[
[server]
host = "0.0.0.0"
port = 8080
]])

-- Template processing
local result = data.template([[
Hello {{ name }}, version {{ version }}
]], {
    name = "User",
    version = "1.0"
})
```

---

## 🔐 Reliability & Retry

### Module `reliability` - Reliability

Adds retry, circuit breaker, etc.

```lua
-- Retry with backoff
reliability.retry(function()
    -- Operation that may fail
    exec.command("curl https://api.example.com")
end, {
    max_attempts = 3,
    initial_delay = 1,  -- seconds
    max_delay = 30,
    backoff_factor = 2  -- exponential backoff
})

-- Circuit breaker
reliability.circuit_breaker(function()
    -- Protected operation
    http.get("https://external-api.com/data")
end, {
    failure_threshold = 5,
    timeout = 60,  -- seconds before retry
    success_threshold = 2  -- successes before closing
})

-- Timeout
reliability.timeout(function()
    -- Operation with timeout
    exec.command("long-running-command")
end, 30)  -- 30 seconds
```

---

## 🎯 Global Modules (No require!)

The following modules are available globally without needing `require()`:

- `log` - Logging
- `exec` - Command execution
- `file` - File operations
- `pkg` - Package management
- `systemd` - Systemd
- `docker` - Docker
- `git` - Git
- `state` - State management
- `facts` - System facts
- `metrics` - Metrics

---

## 📚 Next Steps

- [📋 CLI Reference](cli-reference.md) - All CLI commands
- [🎨 Web UI](web-ui-complete.md) - Web interface guide
- [🎯 Examples](../en/advanced-examples.md) - Practical examples

---

**Last updated:** 2025-10-07
