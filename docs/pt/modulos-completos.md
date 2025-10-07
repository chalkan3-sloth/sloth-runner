# 🔧 Referência Completa de Módulos

## Visão Geral

O Sloth Runner possui mais de 40 módulos integrados que fornecem funcionalidades desde operações básicas de sistema até integrações complexas com provedores cloud. Esta documentação cobre **todos** os módulos disponíveis com exemplos práticos.

---

## 📦 Gerenciamento de Pacotes

### Módulo `pkg` - Gerenciamento de Pacotes

Gerencia pacotes do sistema usando apt, yum, dnf, pacman, brew, etc.

**Funções:**

#### `pkg.install(nome, opções)`

Instala um ou mais pacotes.

```lua
-- Instalar um pacote
pkg.install("nginx")

-- Instalar múltiplos pacotes
pkg.install({"nginx", "postgresql", "redis"})

-- Com opções
pkg.install("nginx", {
    update_cache = true,  -- Atualiza cache antes de instalar
    state = "present"     -- present (padrão) ou latest
})

-- Instalar versão específica (apt)
pkg.install("nginx=1.18.0-0ubuntu1")
```

#### `pkg.remove(nome, opções)`

Remove um ou mais pacotes.

```lua
-- Remover um pacote
pkg.remove("nginx")

-- Remover múltiplos
pkg.remove({"nginx", "apache2"})

-- Remover com purge (apt)
pkg.remove("nginx", { purge = true })
```

#### `pkg.update()`

Atualiza o cache de pacotes.

```lua
-- Atualiza cache (apt update, yum update, etc)
pkg.update()
```

#### `pkg.upgrade(nome)`

Atualiza pacotes instalados.

```lua
-- Atualizar todos os pacotes
pkg.upgrade()

-- Atualizar pacote específico
pkg.upgrade("nginx")
```

**Exemplo completo:**

```yaml
tasks:
  - name: Setup web server
    exec:
      script: |
        -- Atualizar cache
        pkg.update()

        -- Instalar pacotes necessários
        pkg.install({
          "nginx",
          "certbot",
          "python3-certbot-nginx"
        }, { state = "latest" })

        -- Remover servidor web antigo
        pkg.remove("apache2", { purge = true })
```

---

### Módulo `user` - Gerenciamento de Usuários

Gerencia usuários e grupos do sistema.

**Funções:**

#### `user.create(nome, opções)`

Cria um usuário.

```lua
-- Criar usuário simples
user.create("deploy")

-- Com opções completas
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

#### `user.remove(nome, opções)`

Remove um usuário.

```lua
-- Remover usuário
user.remove("olduser")

-- Remover e deletar home
user.remove("olduser", { remove_home = true })
```

#### `user.exists(nome)`

Verifica se usuário existe.

```lua
if user.exists("deploy") then
    log.info("User deploy exists")
else
    user.create("deploy")
end
```

#### `group.create(nome, opções)`

Cria um grupo.

```lua
group.create("developers")
group.create("developers", { gid = 2000 })
```

---

## 📁 Operações de Arquivos

### Módulo `file` - Operações com Arquivos

Gerencia arquivos e diretórios.

**Funções:**

#### `file.copy(origem, destino, opções)`

Copia arquivos ou diretórios.

```lua
-- Copiar arquivo
file.copy("/src/app.conf", "/etc/app/app.conf")

-- Com opções
file.copy("/src/app.conf", "/etc/app/app.conf", {
    owner = "root",
    group = "root",
    mode = "0644",
    backup = true  -- Faz backup se destino existir
})

-- Copiar diretório recursivamente
file.copy("/src/configs/", "/etc/myapp/", {
    recursive = true
})
```

#### `file.create(caminho, opções)`

Cria um arquivo.

```lua
-- Criar arquivo vazio
file.create("/var/log/myapp.log")

-- Com conteúdo e permissões
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

#### `file.remove(caminho, opções)`

Remove arquivos ou diretórios.

```lua
-- Remover arquivo
file.remove("/tmp/cache.dat")

-- Remover diretório recursivamente
file.remove("/var/cache/oldapp", { recursive = true })

-- Remover com force
file.remove("/var/log/*.log", { force = true })
```

#### `file.exists(caminho)`

Verifica se arquivo/diretório existe.

```lua
if file.exists("/etc/nginx/nginx.conf") then
    log.info("Nginx config found")
end
```

#### `file.chmod(caminho, modo)`

Altera permissões.

```lua
file.chmod("/usr/local/bin/myapp", "0755")
file.chmod("/etc/ssl/private/key.pem", "0600")
```

#### `file.chown(caminho, owner, group)`

Altera dono e grupo.

```lua
file.chown("/var/www/html", "www-data", "www-data")
```

#### `file.read(caminho)`

Lê conteúdo de arquivo.

```lua
local content = file.read("/etc/hostname")
log.info("Hostname: " .. content)
```

#### `file.write(caminho, conteúdo, opções)`

Escreve conteúdo em arquivo.

```lua
file.write("/etc/motd", "Welcome to Production Server\n")

-- Com append
file.write("/var/log/app.log", "Log entry\n", {
    append = true
})
```

---

### Módulo `template` - Templates

Processa templates com variáveis.

```lua
-- Template Jinja2/Go template
template.render("/templates/nginx.conf.j2", "/etc/nginx/nginx.conf", {
    server_name = "example.com",
    port = 80,
    root = "/var/www/html"
})
```

---

### Módulo `stow` - Gerenciamento de Dotfiles

Gerencia dotfiles usando GNU Stow.

```lua
-- Fazer stow de dotfiles
stow.link("~/.dotfiles/vim", "~")
stow.link("~/.dotfiles/zsh", "~")

-- Unstow
stow.unlink("~/.dotfiles/vim", "~")

-- Restow (unstow + stow)
stow.restow("~/.dotfiles/vim", "~")
```

---

## 🐚 Execução de Comandos

### Módulo `exec` - Execução de Comandos

Executa comandos do sistema.

**Funções:**

#### `exec.command(comando, opções)`

Executa um comando.

```lua
-- Comando simples
local result = exec.command("ls -la /tmp")

-- Com opções
local result = exec.command("systemctl restart nginx", {
    user = "root",
    cwd = "/etc/nginx",
    env = {
        PATH = "/usr/local/bin:/usr/bin:/bin"
    },
    timeout = 30  -- segundos
})

-- Verificar resultado
if result.exit_code == 0 then
    log.info("Success: " .. result.stdout)
else
    log.error("Failed: " .. result.stderr)
end
```

#### `exec.shell(script)`

Executa script shell.

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

#### `exec.script(caminho, opções)`

Executa script de arquivo.

```lua
exec.script("/scripts/deploy.sh")

exec.script("/scripts/deploy.sh", {
    interpreter = "/bin/bash",
    args = {"production", "v1.2.3"}
})
```

---

### Módulo `goroutine` - Execução Paralela

Executa tarefas em paralelo usando goroutines.

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

-- Com limite de concorrência
goroutine.parallel({
    tasks = {
        function() exec.command("task1") end,
        function() exec.command("task2") end,
        function() exec.command("task3") end
    },
    max_concurrent = 2  -- Máximo 2 por vez
})
```

---

## 🐳 Containers e Virtualização

### Módulo `docker` - Docker

Gerencia containers, imagens e redes Docker.

**Funções:**

#### `docker.container_run(imagem, opções)`

Executa um container.

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

#### `docker.container_stop(nome)`

Para um container.

```lua
docker.container_stop("web-server")
```

#### `docker.container_remove(nome, opções)`

Remove um container.

```lua
docker.container_remove("web-server")
docker.container_remove("web-server", { force = true, volumes = true })
```

#### `docker.image_pull(imagem, opções)`

Baixa uma imagem.

```lua
docker.image_pull("nginx:latest")
docker.image_pull("myregistry.com/myapp:v1.2.3", {
    auth = {
        username = "user",
        password = "pass"
    }
})
```

#### `docker.image_build(contexto, opções)`

Constrói uma imagem.

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

#### `docker.network_create(nome, opções)`

Cria uma rede.

```lua
docker.network_create("app-network", {
    driver = "bridge",
    subnet = "172.20.0.0/16"
})
```

#### `docker.compose_up(arquivo, opções)`

Executa docker-compose.

```lua
docker.compose_up("docker-compose.yml", {
    project_name = "myapp",
    detach = true,
    build = true
})
```

**Exemplo completo:**

```yaml
tasks:
  - name: Deploy application with Docker
    exec:
      script: |
        -- Criar rede
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

### Módulo `incus` - Incus/LXC Containers

Gerencia containers e VMs Incus (LXC).

**Funções:**

#### `incus.launch(imagem, nome, opções)`

Cria e inicia um container/VM.

```lua
-- Container Ubuntu
incus.launch("ubuntu:22.04", "web-01", {
    type = "container",  -- ou "virtual-machine"
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

-- VM com cloud-init
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

#### `incus.exec(nome, comando)`

Executa comando em container.

```lua
incus.exec("web-01", "apt update && apt install -y nginx")
```

#### `incus.file_push(origem, nome, destino)`

Envia arquivo para container.

```lua
incus.file_push("/local/app.conf", "web-01", "/etc/app/app.conf")
```

#### `incus.file_pull(nome, origem, destino)`

Baixa arquivo de container.

```lua
incus.file_pull("web-01", "/var/log/app.log", "/backup/app.log")
```

#### `incus.stop(nome)`

Para um container.

```lua
incus.stop("web-01")
incus.stop("web-01", { force = true })
```

#### `incus.delete(nome)`

Remove um container.

```lua
incus.delete("web-01")
incus.delete("web-01", { force = true })
```

---

## ☁️ Provedores Cloud

### Módulo `aws` - Amazon Web Services

Gerencia recursos AWS (EC2, S3, RDS, etc).

**Funções:**

#### `aws.ec2_instance_create(opções)`

Cria instância EC2.

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

#### `aws.s3_bucket_create(nome, opções)`

Cria bucket S3.

```lua
aws.s3_bucket_create("my-app-backups", {
    region = "us-east-1",
    acl = "private",
    versioning = true,
    encryption = "AES256"
})
```

#### `aws.s3_upload(arquivo, bucket, key)`

Faz upload para S3.

```lua
aws.s3_upload("/backup/db.sql.gz", "my-backups", "db/2024/backup.sql.gz")
```

#### `aws.rds_instance_create(opções)`

Cria instância RDS.

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

### Módulo `azure` - Microsoft Azure

Gerencia recursos Azure.

```lua
-- Criar VM
azure.vm_create({
    name = "web-vm-01",
    resource_group = "production",
    location = "eastus",
    size = "Standard_D2s_v3",
    image = "Ubuntu2204",
    admin_username = "azureuser",
    ssh_key = file.read("~/.ssh/id_rsa.pub")
})

-- Criar Storage Account
azure.storage_account_create({
    name = "myappstorage",
    resource_group = "production",
    location = "eastus",
    sku = "Standard_LRS"
})
```

---

### Módulo `gcp` - Google Cloud Platform

Gerencia recursos GCP.

```lua
-- Criar instância Compute Engine
gcp.compute_instance_create({
    name = "web-instance-01",
    zone = "us-central1-a",
    machine_type = "e2-medium",
    image_project = "ubuntu-os-cloud",
    image_family = "ubuntu-2204-lts",
    tags = {"http-server", "https-server"}
})

-- Criar bucket Cloud Storage
gcp.storage_bucket_create("my-app-data", {
    location = "US",
    storage_class = "STANDARD"
})
```

---

### Módulo `digitalocean` - DigitalOcean

Gerencia recursos DigitalOcean.

```lua
-- Criar Droplet
digitalocean.droplet_create({
    name = "web-01",
    region = "nyc3",
    size = "s-2vcpu-4gb",
    image = "ubuntu-22-04-x64",
    ssh_keys = [123456],
    backups = true,
    monitoring = true
})

-- Criar Load Balancer
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

### Módulo `terraform` - Terraform

Gerencia infraestrutura com Terraform.

**Funções:**

#### `terraform.init(dir, opções)`

Inicializa Terraform.

```lua
terraform.init("/infra/terraform", {
    backend_config = {
        bucket = "my-tf-state",
        key = "prod/terraform.tfstate",
        region = "us-east-1"
    }
})
```

#### `terraform.plan(dir, opções)`

Executa plan.

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

#### `terraform.apply(dir, opções)`

Aplica mudanças.

```lua
terraform.apply("/infra/terraform", {
    plan_file = "tfplan",
    auto_approve = true
})
```

#### `terraform.destroy(dir, opções)`

Destrói recursos.

```lua
terraform.destroy("/infra/terraform", {
    var_file = "prod.tfvars",
    auto_approve = false  -- Pede confirmação
})
```

**Exemplo completo:**

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

        -- Apply se plan estiver ok
        if plan.changes > 0 then
            terraform.apply(tf_dir, {
                auto_approve = true
            })
        end
```

---

### Módulo `pulumi` - Pulumi

Gerencia infraestrutura com Pulumi.

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

## 🔐 Git e Controle de Versão

### Módulo `git` - Git

Operações com repositórios Git.

**Funções:**

#### `git.clone(url, destino, opções)`

Clona um repositório.

```lua
git.clone("https://github.com/user/repo.git", "/opt/app")

-- Com opções
git.clone("https://github.com/user/repo.git", "/opt/app", {
    branch = "main",
    depth = 1,  -- Shallow clone
    auth = {
        username = "user",
        password = "token"
    }
})
```

#### `git.pull(dir, opções)`

Atualiza repositório.

```lua
git.pull("/opt/app")
git.pull("/opt/app", { branch = "develop" })
```

#### `git.checkout(dir, ref)`

Faz checkout de branch/tag.

```lua
git.checkout("/opt/app", "v1.2.3")
git.checkout("/opt/app", "develop")
```

#### `git.commit(dir, mensagem, opções)`

Cria commit.

```lua
git.commit("/opt/app", "Update config files", {
    author = "Deploy Bot <bot@example.com>",
    add_all = true
})
```

#### `git.push(dir, opções)`

Push para remote.

```lua
git.push("/opt/app")
git.push("/opt/app", {
    remote = "origin",
    branch = "main",
    force = false
})
```

---

### Módulo `gitops` - GitOps

Implementa padrões GitOps.

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

## 🌐 Rede e SSH

### Módulo `net` - Networking

Operações de rede.

```lua
-- Check se host está online
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

-- Download arquivo
net.download("https://example.com/file.tar.gz", "/tmp/file.tar.gz")
```

---

### Módulo `ssh` - SSH

Executa comandos via SSH.

```lua
-- Conectar e executar
ssh.exec("user@192.168.1.100", "ls -la /opt", {
    key = "~/.ssh/id_rsa",
    port = 22
})

-- Upload arquivo
ssh.upload("user@192.168.1.100", "/local/file.txt", "/remote/file.txt")

-- Download arquivo
ssh.download("user@192.168.1.100", "/remote/log.txt", "/local/log.txt")
```

---

## ⚙️ Serviços e Systemd

### Módulo `systemd` - Systemd

Gerencia serviços systemd.

**Funções:**

#### `systemd.service_start(nome)`

Inicia um serviço.

```lua
systemd.service_start("nginx")
```

#### `systemd.service_stop(nome)`

Para um serviço.

```lua
systemd.service_stop("nginx")
```

#### `systemd.service_restart(nome)`

Reinicia um serviço.

```lua
systemd.service_restart("nginx")
```

#### `systemd.service_enable(nome)`

Habilita serviço no boot.

```lua
systemd.service_enable("nginx")
```

#### `systemd.service_disable(nome)`

Desabilita serviço no boot.

```lua
systemd.service_disable("apache2")
```

#### `systemd.service_status(nome)`

Verifica status.

```lua
local status = systemd.service_status("nginx")
if status.active then
    log.info("Nginx is running")
end
```

#### `systemd.unit_reload()`

Recarrega unidades systemd.

```lua
systemd.unit_reload()
```

**Exemplo completo:**

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

## 📊 Métricas e Monitoramento

### Módulo `metrics` - Métricas

Coleta e envia métricas.

```lua
-- Contador
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

### Módulo `log` - Logging

Sistema de logging avançado.

```lua
-- Níveis de log
log.debug("Debug message")
log.info("Info message")
log.warn("Warning message")
log.error("Error message")

-- Com campos estruturados
log.info("User login", {
    user_id = 123,
    ip = "192.168.1.100",
    timestamp = os.time()
})

-- Error com stack trace
log.error("Failed to connect", {
    error = err,
    component = "database"
})
```

---

## 🔔 Notificações

### Módulo `notifications` - Notificações

Envia notificações para vários serviços.

**Funções:**

#### `notifications.slack(webhook, mensagem, opções)`

Envia para Slack.

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

#### `notifications.email(opções)`

Envia email.

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

#### `notifications.discord(webhook, mensagem)`

Envia para Discord.

```lua
notifications.discord(
    "https://discord.com/api/webhooks/XXX/YYY",
    "Deploy completed! :tada:"
)
```

#### `notifications.telegram(token, chat_id, mensagem)`

Envia para Telegram.

```lua
notifications.telegram(
    "bot123456:ABC-DEF",
    "123456789",
    "Deploy finished successfully"
)
```

---

## 🧪 Testes e Validação

### Módulo `infra_test` - Testes de Infraestrutura

Testa e valida infraestrutura.

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

## 📡 Facts - Informações do Sistema

### Módulo `facts` - System Facts

Coleta informações do sistema.

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

## 🔄 Estado e Persistência

### Módulo `state` - State Management

Gerencia estado entre execuções.

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

### Módulo `python` - Python

Executa código Python.

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

### Módulo `salt` - SaltStack

Integração com SaltStack.

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

### Módulo `data` - Data Processing

Processa e transforma dados.

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

### Módulo `reliability` - Confiabilidade

Adiciona retry, circuit breaker, etc.

```lua
-- Retry com backoff
reliability.retry(function()
    -- Operação que pode falhar
    exec.command("curl https://api.example.com")
end, {
    max_attempts = 3,
    initial_delay = 1,  -- segundos
    max_delay = 30,
    backoff_factor = 2  -- exponential backoff
})

-- Circuit breaker
reliability.circuit_breaker(function()
    -- Operação protegida
    http.get("https://external-api.com/data")
end, {
    failure_threshold = 5,
    timeout = 60,  -- segundos antes de tentar novamente
    success_threshold = 2  -- sucessos antes de fechar
})

-- Timeout
reliability.timeout(function()
    -- Operação com timeout
    exec.command("long-running-command")
end, 30)  -- 30 segundos
```

---

## 🎯 Módulos Globais (Sem require!)

Os seguintes módulos estão disponíveis globalmente sem necessidade de `require()`:

- `log` - Logging
- `exec` - Execução de comandos
- `file` - Operações com arquivos
- `pkg` - Gerenciamento de pacotes
- `systemd` - Systemd
- `docker` - Docker
- `git` - Git
- `state` - State management
- `facts` - System facts
- `metrics` - Métricas

---

## 📚 Próximos Passos

- [📋 Referência CLI](referencia-cli.md) - Todos os comandos CLI
- [🎨 Web UI](web-ui-completo.md) - Guia da interface web
- [🎯 Exemplos](../en/advanced-examples.md) - Exemplos práticos

---

**Última atualização:** 2025-10-07
