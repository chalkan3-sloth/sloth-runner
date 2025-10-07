# 🔧 模块完整参考

## 概述

Sloth Runner 拥有 40 多个集成模块，提供从基本系统操作到复杂云提供商集成的功能。本文档涵盖了**所有**可用模块及实际示例。

---

## 📦 包管理

### 模块 `pkg` - 包管理

使用 apt、yum、dnf、pacman、brew 等管理系统包。

**函数：**

#### `pkg.install(名称, 选项)`

安装一个或多个包。

```lua
-- 安装单个包
pkg.install("nginx")

-- 安装多个包
pkg.install({"nginx", "postgresql", "redis"})

-- 使用选项
pkg.install("nginx", {
    update_cache = true,  -- 安装前更新缓存
    state = "present"     -- present（默认）或 latest
})

-- 安装特定版本（apt）
pkg.install("nginx=1.18.0-0ubuntu1")
```

#### `pkg.remove(名称, 选项)`

删除一个或多个包。

```lua
-- 删除单个包
pkg.remove("nginx")

-- 删除多个
pkg.remove({"nginx", "apache2"})

-- 使用 purge 删除（apt）
pkg.remove("nginx", { purge = true })
```

#### `pkg.update()`

更新包缓存。

```lua
-- 更新缓存（apt update、yum update 等）
pkg.update()
```

#### `pkg.upgrade(名称)`

升级已安装的包。

```lua
-- 升级所有包
pkg.upgrade()

-- 升级特定包
pkg.upgrade("nginx")
```

**完整示例：**

```yaml
tasks:
  - name: Setup web server
    exec:
      script: |
        -- 更新缓存
        pkg.update()

        -- 安装必需的包
        pkg.install({
          "nginx",
          "certbot",
          "python3-certbot-nginx"
        }, { state = "latest" })

        -- 删除旧的 web 服务器
        pkg.remove("apache2", { purge = true })
```

---

### 模块 `user` - 用户管理

管理系统用户和组。

**函数：**

#### `user.create(名称, 选项)`

创建用户。

```lua
-- 创建简单用户
user.create("deploy")

-- 使用完整选项
user.create("deploy", {
    uid = 1001,
    gid = 1001,
    groups = {"sudo", "docker"},
    shell = "/bin/bash",
    home = "/home/deploy",
    create_home = true,
    system = false,
    comment = "部署用户"
})
```

#### `user.remove(名称, 选项)`

删除用户。

```lua
-- 删除用户
user.remove("olduser")

-- 删除并删除家目录
user.remove("olduser", { remove_home = true })
```

#### `user.exists(名称)`

检查用户是否存在。

```lua
if user.exists("deploy") then
    log.info("用户 deploy 存在")
else
    user.create("deploy")
end
```

#### `group.create(名称, 选项)`

创建组。

```lua
group.create("developers")
group.create("developers", { gid = 2000 })
```

---

## 📁 文件操作

### 模块 `file` - 文件操作

管理文件和目录。

**函数：**

#### `file.copy(源, 目标, 选项)`

复制文件或目录。

```lua
-- 复制文件
file.copy("/src/app.conf", "/etc/app/app.conf")

-- 使用选项
file.copy("/src/app.conf", "/etc/app/app.conf", {
    owner = "root",
    group = "root",
    mode = "0644",
    backup = true  -- 如果目标存在则备份
})

-- 递归复制目录
file.copy("/src/configs/", "/etc/myapp/", {
    recursive = true
})
```

#### `file.create(路径, 选项)`

创建文件。

```lua
-- 创建空文件
file.create("/var/log/myapp.log")

-- 使用内容和权限
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

#### `file.remove(路径, 选项)`

删除文件或目录。

```lua
-- 删除文件
file.remove("/tmp/cache.dat")

-- 递归删除目录
file.remove("/var/cache/oldapp", { recursive = true })

-- 强制删除
file.remove("/var/log/*.log", { force = true })
```

#### `file.exists(路径)`

检查文件/目录是否存在。

```lua
if file.exists("/etc/nginx/nginx.conf") then
    log.info("找到 Nginx 配置")
end
```

#### `file.chmod(路径, 模式)`

更改权限。

```lua
file.chmod("/usr/local/bin/myapp", "0755")
file.chmod("/etc/ssl/private/key.pem", "0600")
```

#### `file.chown(路径, 所有者, 组)`

更改所有者和组。

```lua
file.chown("/var/www/html", "www-data", "www-data")
```

#### `file.read(路径)`

读取文件内容。

```lua
local content = file.read("/etc/hostname")
log.info("主机名：" .. content)
```

#### `file.write(路径, 内容, 选项)`

写入内容到文件。

```lua
file.write("/etc/motd", "欢迎来到生产服务器\n")

-- 使用追加
file.write("/var/log/app.log", "日志条目\n", {
    append = true
})
```

---

### 模块 `template` - 模板

使用变量处理模板。

```lua
-- Jinja2/Go 模板
template.render("/templates/nginx.conf.j2", "/etc/nginx/nginx.conf", {
    server_name = "example.com",
    port = 80,
    root = "/var/www/html"
})
```

---

### 模块 `stow` - Dotfiles 管理

使用 GNU Stow 管理 dotfiles。

```lua
-- Stow dotfiles
stow.link("~/.dotfiles/vim", "~")
stow.link("~/.dotfiles/zsh", "~")

-- Unstow
stow.unlink("~/.dotfiles/vim", "~")

-- Restow（unstow + stow）
stow.restow("~/.dotfiles/vim", "~")
```

---

## 🐚 命令执行

### 模块 `exec` - 命令执行

执行系统命令。

**函数：**

#### `exec.command(命令, 选项)`

执行命令。

```lua
-- 简单命令
local result = exec.command("ls -la /tmp")

-- 使用选项
local result = exec.command("systemctl restart nginx", {
    user = "root",
    cwd = "/etc/nginx",
    env = {
        PATH = "/usr/local/bin:/usr/bin:/bin"
    },
    timeout = 30  -- 秒
})

-- 检查结果
if result.exit_code == 0 then
    log.info("成功：" .. result.stdout)
else
    log.error("失败：" .. result.stderr)
end
```

#### `exec.shell(脚本)`

执行 shell 脚本。

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

#### `exec.script(路径, 选项)`

从文件执行脚本。

```lua
exec.script("/scripts/deploy.sh")

exec.script("/scripts/deploy.sh", {
    interpreter = "/bin/bash",
    args = {"production", "v1.2.3"}
})
```

---

### 模块 `goroutine` - 并行执行

使用 goroutine 并行执行任务。

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

-- 限制并发
goroutine.parallel({
    tasks = {
        function() exec.command("task1") end,
        function() exec.command("task2") end,
        function() exec.command("task3") end
    },
    max_concurrent = 2  -- 最多同时 2 个
})
```

---

## 🐳 容器和虚拟化

### 模块 `docker` - Docker

管理 Docker 容器、镜像和网络。

**函数：**

#### `docker.container_run(镜像, 选项)`

运行容器。

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

#### `docker.container_stop(名称)`

停止容器。

```lua
docker.container_stop("web-server")
```

#### `docker.container_remove(名称, 选项)`

删除容器。

```lua
docker.container_remove("web-server")
docker.container_remove("web-server", { force = true, volumes = true })
```

#### `docker.image_pull(镜像, 选项)`

拉取镜像。

```lua
docker.image_pull("nginx:latest")
docker.image_pull("myregistry.com/myapp:v1.2.3", {
    auth = {
        username = "user",
        password = "pass"
    }
})
```

#### `docker.image_build(上下文, 选项)`

构建镜像。

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

#### `docker.network_create(名称, 选项)`

创建网络。

```lua
docker.network_create("app-network", {
    driver = "bridge",
    subnet = "172.20.0.0/16"
})
```

#### `docker.compose_up(文件, 选项)`

运行 docker-compose。

```lua
docker.compose_up("docker-compose.yml", {
    project_name = "myapp",
    detach = true,
    build = true
})
```

**完整示例：**

```yaml
tasks:
  - name: 使用 Docker 部署应用
    exec:
      script: |
        -- 创建网络
        docker.network_create("app-net")

        -- 数据库
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

        -- 应用程序
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

### 模块 `incus` - Incus/LXC 容器

管理 Incus（LXC）容器和虚拟机。

**函数：**

#### `incus.launch(镜像, 名称, 选项)`

创建并启动容器/虚拟机。

```lua
-- Ubuntu 容器
incus.launch("ubuntu:22.04", "web-01", {
    type = "container",  -- 或 "virtual-machine"
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

-- 带 cloud-init 的虚拟机
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

#### `incus.exec(名称, 命令)`

在容器中执行命令。

```lua
incus.exec("web-01", "apt update && apt install -y nginx")
```

#### `incus.file_push(源, 名称, 目标)`

发送文件到容器。

```lua
incus.file_push("/local/app.conf", "web-01", "/etc/app/app.conf")
```

#### `incus.file_pull(名称, 源, 目标)`

从容器下载文件。

```lua
incus.file_pull("web-01", "/var/log/app.log", "/backup/app.log")
```

#### `incus.stop(名称)`

停止容器。

```lua
incus.stop("web-01")
incus.stop("web-01", { force = true })
```

#### `incus.delete(名称)`

删除容器。

```lua
incus.delete("web-01")
incus.delete("web-01", { force = true })
```

---

## ☁️ 云提供商

### 模块 `aws` - Amazon Web Services

管理 AWS 资源（EC2、S3、RDS 等）。

**函数：**

#### `aws.ec2_instance_create(选项)`

创建 EC2 实例。

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

#### `aws.s3_bucket_create(名称, 选项)`

创建 S3 存储桶。

```lua
aws.s3_bucket_create("my-app-backups", {
    region = "us-east-1",
    acl = "private",
    versioning = true,
    encryption = "AES256"
})
```

#### `aws.s3_upload(文件, 存储桶, 键)`

上传到 S3。

```lua
aws.s3_upload("/backup/db.sql.gz", "my-backups", "db/2024/backup.sql.gz")
```

#### `aws.rds_instance_create(选项)`

创建 RDS 实例。

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

### 模块 `azure` - Microsoft Azure

管理 Azure 资源。

```lua
-- 创建虚拟机
azure.vm_create({
    name = "web-vm-01",
    resource_group = "production",
    location = "eastus",
    size = "Standard_D2s_v3",
    image = "Ubuntu2204",
    admin_username = "azureuser",
    ssh_key = file.read("~/.ssh/id_rsa.pub")
})

-- 创建存储账户
azure.storage_account_create({
    name = "myappstorage",
    resource_group = "production",
    location = "eastus",
    sku = "Standard_LRS"
})
```

---

### 模块 `gcp` - Google Cloud Platform

管理 GCP 资源。

```lua
-- 创建 Compute Engine 实例
gcp.compute_instance_create({
    name = "web-instance-01",
    zone = "us-central1-a",
    machine_type = "e2-medium",
    image_project = "ubuntu-os-cloud",
    image_family = "ubuntu-2204-lts",
    tags = {"http-server", "https-server"}
})

-- 创建 Cloud Storage 存储桶
gcp.storage_bucket_create("my-app-data", {
    location = "US",
    storage_class = "STANDARD"
})
```

---

### 模块 `digitalocean` - DigitalOcean

管理 DigitalOcean 资源。

```lua
-- 创建 Droplet
digitalocean.droplet_create({
    name = "web-01",
    region = "nyc3",
    size = "s-2vcpu-4gb",
    image = "ubuntu-22-04-x64",
    ssh_keys = [123456],
    backups = true,
    monitoring = true
})

-- 创建负载均衡器
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

## 🏗️ 基础设施即代码

### 模块 `terraform` - Terraform

使用 Terraform 管理基础设施。

**函数：**

#### `terraform.init(目录, 选项)`

初始化 Terraform。

```lua
terraform.init("/infra/terraform", {
    backend_config = {
        bucket = "my-tf-state",
        key = "prod/terraform.tfstate",
        region = "us-east-1"
    }
})
```

#### `terraform.plan(目录, 选项)`

执行计划。

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

#### `terraform.apply(目录, 选项)`

应用更改。

```lua
terraform.apply("/infra/terraform", {
    plan_file = "tfplan",
    auto_approve = true
})
```

#### `terraform.destroy(目录, 选项)`

销毁资源。

```lua
terraform.destroy("/infra/terraform", {
    var_file = "prod.tfvars",
    auto_approve = false  -- 请求确认
})
```

**完整示例：**

```yaml
tasks:
  - name: 部署基础设施
    exec:
      script: |
        local tf_dir = "/infra/terraform"

        -- 初始化
        terraform.init(tf_dir)

        -- 计划
        local plan = terraform.plan(tf_dir, {
            var_file = "production.tfvars"
        })

        -- 如果计划有变更则应用
        if plan.changes > 0 then
            terraform.apply(tf_dir, {
                auto_approve = true
            })
        end
```

---

### 模块 `pulumi` - Pulumi

使用 Pulumi 管理基础设施。

```lua
-- 初始化堆栈
pulumi.stack_init("production")

-- 配置
pulumi.config_set("aws:region", "us-east-1")

-- 部署
pulumi.up({
    stack = "production",
    yes = true  -- 自动批准
})

-- 销毁
pulumi.destroy({
    stack = "production",
    yes = false
})
```

---

## 🔐 Git 和版本控制

### 模块 `git` - Git

Git 仓库操作。

**函数：**

#### `git.clone(url, 目标, 选项)`

克隆仓库。

```lua
git.clone("https://github.com/user/repo.git", "/opt/app")

-- 使用选项
git.clone("https://github.com/user/repo.git", "/opt/app", {
    branch = "main",
    depth = 1,  -- 浅克隆
    auth = {
        username = "user",
        password = "token"
    }
})
```

#### `git.pull(目录, 选项)`

更新仓库。

```lua
git.pull("/opt/app")
git.pull("/opt/app", { branch = "develop" })
```

#### `git.checkout(目录, 引用)`

切换分支/标签。

```lua
git.checkout("/opt/app", "v1.2.3")
git.checkout("/opt/app", "develop")
```

#### `git.commit(目录, 消息, 选项)`

创建提交。

```lua
git.commit("/opt/app", "更新配置文件", {
    author = "部署机器人 <bot@example.com>",
    add_all = true
})
```

#### `git.push(目录, 选项)`

推送到远程。

```lua
git.push("/opt/app")
git.push("/opt/app", {
    remote = "origin",
    branch = "main",
    force = false
})
```

---

### 模块 `gitops` - GitOps

实现 GitOps 模式。

```lua
-- 从 Git 同步
gitops.sync({
    repo = "https://github.com/org/k8s-manifests.git",
    branch = "main",
    path = "production/",
    destination = "/opt/k8s-manifests"
})

-- 应用清单
gitops.apply({
    source = "/opt/k8s-manifests",
    namespace = "production"
})
```

---

## 🌐 网络和 SSH

### 模块 `net` - 网络

网络操作。

```lua
-- 检查主机是否在线
if net.ping("example.com") then
    log.info("主机在线")
end

-- 端口扫描
local open = net.port_open("example.com", 443)

-- HTTP 请求
local response = net.http_get("https://api.example.com/status")
if response.status == 200 then
    log.info(response.body)
end

-- 下载文件
net.download("https://example.com/file.tar.gz", "/tmp/file.tar.gz")
```

---

### 模块 `ssh` - SSH

通过 SSH 执行命令。

```lua
-- 连接并执行
ssh.exec("user@192.168.1.100", "ls -la /opt", {
    key = "~/.ssh/id_rsa",
    port = 22
})

-- 上传文件
ssh.upload("user@192.168.1.100", "/local/file.txt", "/remote/file.txt")

-- 下载文件
ssh.download("user@192.168.1.100", "/remote/log.txt", "/local/log.txt")
```

---

## ⚙️ 服务和 Systemd

### 模块 `systemd` - Systemd

管理 systemd 服务。

**函数：**

#### `systemd.service_start(名称)`

启动服务。

```lua
systemd.service_start("nginx")
```

#### `systemd.service_stop(名称)`

停止服务。

```lua
systemd.service_stop("nginx")
```

#### `systemd.service_restart(名称)`

重启服务。

```lua
systemd.service_restart("nginx")
```

#### `systemd.service_enable(名称)`

在启动时启用服务。

```lua
systemd.service_enable("nginx")
```

#### `systemd.service_disable(名称)`

在启动时禁用服务。

```lua
systemd.service_disable("apache2")
```

#### `systemd.service_status(名称)`

检查状态。

```lua
local status = systemd.service_status("nginx")
if status.active then
    log.info("Nginx 正在运行")
end
```

#### `systemd.unit_reload()`

重新加载 systemd 单元。

```lua
systemd.unit_reload()
```

**完整示例：**

```yaml
tasks:
  - name: 部署和配置 nginx
    exec:
      script: |
        -- 安装
        pkg.install("nginx")

        -- 配置
        file.copy("/deploy/nginx.conf", "/etc/nginx/nginx.conf")

        -- 启用并启动
        systemd.service_enable("nginx")
        systemd.service_start("nginx")

        -- 验证
        local status = systemd.service_status("nginx")
        if not status.active then
            error("Nginx 启动失败")
        end
```

---

## 📊 指标和监控

### 模块 `metrics` - 指标

收集和发送指标。

```lua
-- 计数器
metrics.counter("requests_total", 1, {
    method = "GET",
    status = "200"
})

-- 仪表
metrics.gauge("memory_usage_bytes", 1024*1024*512)

-- 直方图
metrics.histogram("request_duration_seconds", 0.234)

-- 自定义指标
metrics.custom("app_users_active", 42, {
    type = "gauge",
    labels = {
        region = "us-east-1"
    }
})
```

---

### 模块 `log` - 日志

高级日志系统。

```lua
-- 日志级别
log.debug("调试消息")
log.info("信息消息")
log.warn("警告消息")
log.error("错误消息")

-- 带结构化字段
log.info("用户登录", {
    user_id = 123,
    ip = "192.168.1.100",
    timestamp = os.time()
})

-- 带堆栈跟踪的错误
log.error("连接失败", {
    error = err,
    component = "数据库"
})
```

---

## 🔔 通知

### 模块 `notifications` - 通知

向各种服务发送通知。

**函数：**

#### `notifications.slack(webhook, 消息, 选项)`

发送到 Slack。

```lua
notifications.slack(
    "https://hooks.slack.com/services/XXX/YYY/ZZZ",
    "部署成功完成！:rocket:",
    {
        channel = "#deployments",
        username = "Sloth Runner",
        icon_emoji = ":sloth:"
    }
)
```

#### `notifications.email(选项)`

发送电子邮件。

```lua
notifications.email({
    from = "noreply@example.com",
    to = "admin@example.com",
    subject = "部署状态",
    body = "生产部署已完成",
    smtp_host = "smtp.gmail.com",
    smtp_port = 587,
    smtp_user = "user@gmail.com",
    smtp_pass = "password"
})
```

#### `notifications.discord(webhook, 消息)`

发送到 Discord。

```lua
notifications.discord(
    "https://discord.com/api/webhooks/XXX/YYY",
    "部署完成！:tada:"
)
```

#### `notifications.telegram(令牌, 聊天ID, 消息)`

发送到 Telegram。

```lua
notifications.telegram(
    "bot123456:ABC-DEF",
    "123456789",
    "部署成功完成"
)
```

---

## 🧪 测试和验证

### 模块 `infra_test` - 基础设施测试

测试和验证基础设施。

```lua
-- 测试端口
infra_test.port("example.com", 443, {
    timeout = 5,
    should_be_open = true
})

-- 测试 HTTP
infra_test.http("https://example.com", {
    status_code = 200,
    contains = "欢迎",
    timeout = 10
})

-- 测试命令
infra_test.command("systemctl is-active nginx", {
    exit_code = 0,
    stdout_contains = "active"
})

-- 测试文件
infra_test.file("/etc/nginx/nginx.conf", {
    exists = true,
    mode = "0644",
    owner = "root"
})
```

---

## 📡 Facts - 系统信息

### 模块 `facts` - 系统信息

收集系统信息。

```lua
-- 获取所有信息
local facts = facts.gather()

-- 访问信息
log.info("操作系统：" .. facts.os.name)
log.info("内核：" .. facts.kernel.version)
log.info("CPU 核心：" .. facts.cpu.cores)
log.info("内存：" .. facts.memory.total)
log.info("主机名：" .. facts.hostname)

-- 特定信息
local cpu = facts.cpu()
local mem = facts.memory()
local disk = facts.disk()
local network = facts.network()
```

---

## 🔄 状态和持久化

### 模块 `state` - 状态管理

管理执行之间的状态。

```lua
-- 保存状态
state.set("last_deploy_version", "v1.2.3")
state.set("last_deploy_time", os.time())

-- 获取状态
local last_version = state.get("last_deploy_version")
if last_version == nil then
    log.info("首次部署")
end

-- 检查是否更改
if state.changed("app_config_hash", new_hash) then
    log.info("配置已更改，重启应用")
    systemd.service_restart("myapp")
end

-- 清除状态
state.clear("temporary_data")
```

---

## 🐍 Python 集成

### 模块 `python` - Python

执行 Python 代码。

```lua
-- 运行 Python 脚本
python.run([[
import requests
import json

response = requests.get('https://api.github.com/repos/user/repo')
data = response.json()
print(f"Stars: {data['stargazers_count']}")
]])

-- 运行 Python 文件
python.run_file("/scripts/deploy.py", {
    args = {"production", "v1.2.3"},
    venv = "/opt/venv"
})

-- 安装包
python.pip_install({"requests", "boto3"})
```

---

## 🧂 配置管理

### 模块 `salt` - SaltStack

与 SaltStack 集成。

```lua
-- 应用 Salt 状态
salt.state_apply("webserver", {
    pillar = {
        nginx_port = 80,
        domain = "example.com"
    }
})

-- 运行 Salt 命令
salt.cmd_run("service.restart", {"nginx"})
```

---

## 📊 数据处理

### 模块 `data` - 数据处理

处理和转换数据。

```lua
-- 解析 JSON
local json_data = data.json_parse('{"name": "value"}')

-- 生成 JSON
local json_str = data.json_encode({
    name = "app",
    version = "1.0"
})

-- 解析 YAML
local yaml_data = data.yaml_parse([[
name: myapp
version: 1.0
]])

-- 解析 TOML
local toml_data = data.toml_parse([[
[server]
host = "0.0.0.0"
port = 8080
]])

-- 模板处理
local result = data.template([[
你好 {{ name }}，版本 {{ version }}
]], {
    name = "用户",
    version = "1.0"
})
```

---

## 🔐 可靠性和重试

### 模块 `reliability` - 可靠性

添加重试、断路器等。

```lua
-- 带退避的重试
reliability.retry(function()
    -- 可能失败的操作
    exec.command("curl https://api.example.com")
end, {
    max_attempts = 3,
    initial_delay = 1,  -- 秒
    max_delay = 30,
    backoff_factor = 2  -- 指数退避
})

-- 断路器
reliability.circuit_breaker(function()
    -- 受保护的操作
    http.get("https://external-api.com/data")
end, {
    failure_threshold = 5,
    timeout = 60,  -- 重试前等待的秒数
    success_threshold = 2  -- 关闭前的成功次数
})

-- 超时
reliability.timeout(function()
    -- 带超时的操作
    exec.command("long-running-command")
end, 30)  -- 30 秒
```

---

## 🎯 全局模块（无需 require！）

以下模块无需 `require()` 即可全局使用：

- `log` - 日志
- `exec` - 命令执行
- `file` - 文件操作
- `pkg` - 包管理
- `systemd` - Systemd
- `docker` - Docker
- `git` - Git
- `state` - 状态管理
- `facts` - 系统信息
- `metrics` - 指标

---

## 📚 下一步

- [📋 CLI 参考](cli-reference.md) - 所有 CLI 命令
- [🎨 Web UI](web-ui-complete.md) - Web 界面指南
- [🎯 示例](../en/advanced-examples.md) - 实际示例

---

**最后更新：** 2025-10-07
