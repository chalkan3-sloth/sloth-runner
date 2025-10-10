# Módulo infra_test

O módulo `infra_test` fornece um framework completo de validação e teste de infraestrutura nativo ao Sloth Runner. Ele permite que você insira asserções de teste diretamente nas suas tasks para verificar o resultado de operações de deploy ou configuration management.

## Visão Geral

O `infra_test` é inspirado em ferramentas como Testinfra e InSpec, mas é nativo e integrado ao Sloth Runner, permitindo testes de infraestrutura diretamente nas tasks sem dependências externas.

## Características Principais

- ✅ **Execução Local e Remota**: Todos os testes podem ser executados localmente ou delegados para agentes remotos
- ✅ **Asserções Nativas**: Interrompe a execução da task em caso de falha
- ✅ **Zero Dependências**: Não requer instalação de ferramentas externas
- ✅ **Integração Total**: Funciona perfeitamente com o sistema de agents do Sloth Runner
- ✅ **Detecção Automática de Pacotes**: Suporta apt, yum, pacman, apk e brew automaticamente
- ✅ **Validação de Versões**: Verifica versões específicas de pacotes instalados

## Módulos de Teste Disponíveis

O `infra_test` oferece 6 categorias de testes:

1. **🗂️ Testes de Arquivo** - Verifica existência, permissões, conteúdo e proprietários
2. **🌐 Testes de Rede** - Valida portas, conectividade TCP/UDP e ping
3. **⚙️ Testes de Serviço** - Verifica status de serviços systemd/init
4. **🔄 Testes de Processo** - Valida processos em execução
5. **💻 Testes de Comando** - Executa comandos e valida saídas
6. **📦 Testes de Pacote** - Verifica instalação e versões de pacotes (NOVO!)

## Parâmetro Target

Todas as funções de teste aceitam um parâmetro opcional `target` para especificar onde o teste será executado:

| Parâmetro target | Comportamento |
|-----------------|---------------|
| Omitido ou `"local"` | Executa no agente local (onde a task está rodando) |
| String (nome do agente) | O teste é delegado ao agente remoto especificado |
| `"localhost"` | Força o teste no agente onde a task foi agendada |

## Referência Rápida de Funções

### 🗂️ Testes de Arquivo
- `file_exists(path, [target])` - Verifica existência
- `is_directory(path, [target])` - Verifica se é diretório
- `is_file(path, [target])` - Verifica se é arquivo
- `file_contains(path, pattern, [target])` - Verifica conteúdo
- `file_mode(path, mode, [target])` - Verifica permissões
- `file_owner(path, user, [target])` - Verifica proprietário
- `file_group(path, group, [target])` - Verifica grupo
- `file_size(path, bytes, [target])` - Verifica tamanho

### 🌐 Testes de Rede
- `port_is_listening(port, [target])` - Verifica porta aberta
- `port_is_tcp(port, [target])` - Verifica porta TCP
- `port_is_udp(port, [target])` - Verifica porta UDP
- `can_connect(host, port, [timeout])` - Testa conectividade TCP
- `ping(host, [count], [target])` - Testa conectividade ICMP

### ⚙️ Testes de Serviço
- `service_is_running(name, [target])` - Verifica se serviço está ativo
- `service_is_enabled(name, [target])` - Verifica se está habilitado

### 🔄 Testes de Processo
- `process_is_running(pattern, [target])` - Verifica processo
- `process_count(pattern, count, [target])` - Conta processos

### 💻 Testes de Comando
- `command_succeeds(cmd, [target])` - Verifica exit code 0
- `command_fails(cmd, [target])` - Verifica exit code != 0
- `command_stdout_contains(cmd, pattern, [target])` - Verifica saída
- `command_stderr_is_empty(cmd, [target])` - Verifica stderr vazio
- `command_output_equals(cmd, expected, [target])` - Verifica saída exata

### 📦 Testes de Pacote
- `package_is_installed(name, [target])` - Verifica instalação
- `package_version(name, version, [target])` - Verifica versão

## Modelo de Retorno

- **Sucesso**: A função não retorna nada (ou retorna `true`)
- **Falha**: A função lança um erro que interrompe a execução da task e marca a task como falha

---

## Testes de Arquivo (File Tests)

### file_exists(path, [target])

Verifica se um arquivo ou diretório existe.

**Parâmetros:**
- `path` (string): Caminho do arquivo ou diretório
- `target` (string, opcional): Agente onde executar o teste

**Exemplo:**
```lua
local infra_test = require("infra_test")

workflow.define("test-deployment")
  :description("Test deployment configuration")
  :version("1.0.0")
  :tasks({
    task("verify-config")
      :description("Verify configuration files")
      :command(function(this, params)
        -- Verifica localmente
        infra_test.file_exists("/etc/nginx/nginx.conf")

        -- Verifica em agente remoto
        infra_test.file_exists("/etc/nginx/nginx.conf", "web-server-01")

        return true, "Configuration verified successfully"
      end)
      :build()
  })
  :delegate_to("prod-agent")
```

### is_directory(path, [target])

Verifica se o caminho é um diretório.

**Exemplo:**
```lua
infra_test.is_directory("/var/www/html")
infra_test.is_directory("/opt/app", "app-server")
```

### is_file(path, [target])

Verifica se o caminho é um arquivo regular.

**Exemplo:**
```lua
infra_test.is_file("/etc/hosts")
infra_test.is_file("/var/log/app.log", "log-server")
```

### file_contains(path, pattern, [target])

Verifica se o arquivo contém uma string ou padrão regex.

**Parâmetros:**
- `path` (string): Caminho do arquivo
- `pattern` (string): String ou expressão regular a buscar
- `target` (string, opcional): Agente onde executar

**Exemplo:**
```lua
-- Verifica string simples
infra_test.file_contains("/etc/nginx/nginx.conf", "worker_processes")

-- Verifica com regex
infra_test.file_contains("/var/log/app.log", "ERROR.*database", "app-server")
```

### file_mode(path, mode, [target])

Verifica as permissões do arquivo.

**Parâmetros:**
- `path` (string): Caminho do arquivo
- `mode` (string): Permissões esperadas (ex: "644", "0644", "0o644")
- `target` (string, opcional): Agente onde executar

**Exemplo:**
```lua
infra_test.file_mode("/etc/passwd", "644")
infra_test.file_mode("/root/.ssh/id_rsa", "0600", "bastion")
```

### file_owner(path, user, [target])

Verifica se o proprietário (usuário) do arquivo corresponde.

**Exemplo:**
```lua
infra_test.file_owner("/var/www/html", "www-data")
infra_test.file_owner("/opt/app/config.yaml", "appuser", "app-server")
```

### file_group(path, group, [target])

Verifica se o grupo do arquivo corresponde.

**Exemplo:**
```lua
infra_test.file_group("/var/www/html", "www-data")
infra_test.file_group("/etc/ssl/private", "ssl-cert", "web-server")
```

### file_size(path, size_in_bytes, [target])

Verifica o tamanho exato do arquivo em bytes.

**Exemplo:**
```lua
infra_test.file_size("/etc/machine-id", 33)
infra_test.file_size("/var/cache/app.db", 1048576, "cache-server")
```

---

## Testes de Rede e Porta (Network Tests)

### port_is_listening(port, [target])

Verifica se a porta está aberta/escutando no alvo.

**Exemplo:**
```lua
infra_test.port_is_listening(80)
infra_test.port_is_listening(443, "web-server")
infra_test.port_is_listening(5432, "db-server")
```

### port_is_tcp(port, [target])

Verifica se a porta está escutando usando o protocolo TCP.

**Exemplo:**
```lua
infra_test.port_is_tcp(22)
infra_test.port_is_tcp(3306, "mysql-server")
```

### port_is_udp(port, [target])

Verifica se a porta está escutando usando o protocolo UDP.

**Exemplo:**
```lua
infra_test.port_is_udp(53)
infra_test.port_is_udp(123, "ntp-server")
```

### can_connect(host, port, [timeout_ms])

Testa a conectividade TCP a partir do agente para um host externo/remoto.

**Parâmetros:**
- `host` (string): Host de destino
- `port` (number): Porta de destino
- `timeout_ms` (number, opcional): Timeout em milissegundos (padrão: 5000)

**Exemplo:**
```lua
infra_test.can_connect("google.com", 443)
infra_test.can_connect("database.internal", 5432, 3000)
```

### ping(host, [count], [target])

Testa a conectividade ICMP (ping) para um host.

**Parâmetros:**
- `host` (string): Host de destino
- `count` (number, opcional): Número de pacotes (padrão: 4)
- `target` (string, opcional): Agente onde executar

**Exemplo:**
```lua
infra_test.ping("8.8.8.8")
infra_test.ping("internal-router", 10)
infra_test.ping("remote-server", 5, "edge-agent")
```

---

## Testes de Serviço e Processo (Service & Process Tests)

### service_is_running(name, [target])

Verifica se o serviço está ativo (via systemctl, service, etc.).

**Exemplo:**
```lua
infra_test.service_is_running("nginx")
infra_test.service_is_running("postgresql", "db-server")
```

### service_is_enabled(name, [target])

Verifica se o serviço está habilitado para iniciar no boot.

**Exemplo:**
```lua
infra_test.service_is_enabled("docker")
infra_test.service_is_enabled("nginx", "web-server")
```

### process_is_running(pattern, [target])

Verifica se um processo com um nome ou padrão de comando está em execução.

**Exemplo:**
```lua
infra_test.process_is_running("nginx")
infra_test.process_is_running("java.*spring-boot", "app-server")
```

### process_count(pattern, count, [target])

Verifica se o número de processos corresponde a um valor exato.

**Parâmetros:**
- `pattern` (string): Padrão para buscar processos
- `count` (number): Número esperado de processos
- `target` (string, opcional): Agente onde executar

**Exemplo:**
```lua
infra_test.process_count("nginx", 4)
infra_test.process_count("worker", 8, "worker-node")
```

---

## Testes de Comando e Saída (Command & Output Tests)

### command_succeeds(cmd, [target])

Verifica se o comando retorna o código de saída 0.

**Exemplo:**
```lua
infra_test.command_succeeds("which docker")
infra_test.command_succeeds("systemctl is-active nginx", "web-server")
```

### command_fails(cmd, [target])

Verifica se o comando retorna um código de saída diferente de zero.

**Exemplo:**
```lua
infra_test.command_fails("systemctl is-active fake-service")
infra_test.command_fails("test -f /nonexistent", "app-server")
```

### command_stdout_contains(cmd, pattern, [target])

Verifica se a saída padrão do comando contém uma string ou regex.

**Parâmetros:**
- `cmd` (string): Comando a executar
- `pattern` (string): String ou regex a buscar na saída
- `target` (string, opcional): Agente onde executar

**Exemplo:**
```lua
infra_test.command_stdout_contains("cat /etc/os-release", "Ubuntu")
infra_test.command_stdout_contains("docker --version", "version 20", "docker-host")
```

### command_stderr_is_empty(cmd, [target])

Verifica se a saída de erro do comando está vazia.

**Exemplo:**
```lua
infra_test.command_stderr_is_empty("ls /home")
infra_test.command_stderr_is_empty("cat /etc/hosts", "web-server")
```

### command_output_equals(cmd, expected_output, [target])

Verifica se a saída padrão é exatamente igual ao valor esperado.

**Parâmetros:**
- `cmd` (string): Comando a executar
- `expected_output` (string): Saída esperada
- `target` (string, opcional): Agente onde executar

**Exemplo:**
```lua
infra_test.command_output_equals("whoami", "root")
infra_test.command_output_equals("cat /etc/hostname", "web-01", "web-server")
```

---

## Testes de Pacote (Package Tests)

### package_is_installed(package_name, [target])

Verifica se um pacote está instalado no sistema. O módulo detecta automaticamente o gerenciador de pacotes disponível (apt/dpkg, yum/rpm, pacman, apk, brew).

**Parâmetros:**
- `package_name` (string): Nome do pacote
- `target` (string, opcional): Agente onde executar o teste

**Gerenciadores Suportados:**
- **Debian/Ubuntu**: dpkg
- **RedHat/CentOS/Fedora**: rpm
- **Arch Linux**: pacman
- **Alpine Linux**: apk
- **macOS**: brew

**Exemplo:**
```lua
local infra_test = require("infra_test")

-- Verifica se nginx está instalado localmente
infra_test.package_is_installed("nginx")

-- Verifica em agente remoto
infra_test.package_is_installed("postgresql", "db-server")

-- Verifica múltiplos pacotes
infra_test.package_is_installed("docker-ce")
infra_test.package_is_installed("docker-compose")
infra_test.package_is_installed("git")
```

### package_version(package_name, expected_version, [target])

Verifica a versão de um pacote instalado. Aceita versão exata ou prefixo.

**Parâmetros:**
- `package_name` (string): Nome do pacote
- `expected_version` (string): Versão esperada (ou prefixo da versão)
- `target` (string, opcional): Agente onde executar o teste

**Exemplo:**
```lua
-- Verifica versão exata
infra_test.package_version("nginx", "1.18.0")

-- Verifica prefixo de versão (ex: 1.18.x)
infra_test.package_version("nginx", "1.18", "web-server")

-- Verifica versão major
infra_test.package_version("postgresql", "14", "db-server")
```

---

## Exemplos Completos

### Exemplo 1: Teste de Deploy de Aplicação

```lua
local infra_test = require("infra_test")
local pkg = require("pkg")

workflow.define("deploy-and-test-app")
  :description("Deploy and test nginx application")
  :version("1.0.0")
  :tasks({
    task("install-nginx")
      :description("Install nginx package")
      :command(function(this, params)
        pkg.install("nginx")
        return true, "Nginx installed successfully"
      end)
      :build(),

    task("verify-installation")
      :description("Verify nginx installation")
      :command(function(this, params)
        -- Verifica se o pacote foi instalado
        infra_test.package_is_installed("nginx")

        -- Verifica se os arquivos existem
        infra_test.file_exists("/usr/sbin/nginx")
        infra_test.file_exists("/etc/nginx/nginx.conf")

        -- Verifica se o serviço está rodando e habilitado
        infra_test.service_is_running("nginx")
        infra_test.service_is_enabled("nginx")

        -- Verifica se a porta está aberta
        infra_test.port_is_tcp(80)

        -- Verifica se o processo está ativo
        infra_test.process_is_running("nginx")

        return true, "Installation verified successfully"
      end)
      :build(),

    task("verify-config")
      :description("Verify nginx configuration")
      :command(function(this, params)
        -- Verifica permissões e proprietário
        infra_test.file_mode("/etc/nginx/nginx.conf", "644")
        infra_test.file_owner("/var/www/html", "www-data")

        -- Verifica conteúdo da configuração
        infra_test.file_contains("/etc/nginx/nginx.conf", "worker_processes")

        return true, "Configuration verified successfully"
      end)
      :build()
  })
  :delegate_to("web-server-01")
```

### Exemplo 2: Validação Multi-Agent

```lua
local infra_test = require("infra_test")

workflow.define("test-infrastructure")
  :description("Test infrastructure across multiple agents")
  :version("1.0.0")
  :tasks({
    task("test-web-servers")
      :description("Test multiple web servers")
      :command(function(this, params)
        -- Testa múltiplos servidores web
        local servers = {"web-01", "web-02", "web-03"}

        for _, server in ipairs(servers) do
          print("Testing " .. server)

          infra_test.service_is_running("nginx", server)
          infra_test.port_is_listening(80, server)
          infra_test.port_is_listening(443, server)
          infra_test.file_exists("/var/www/html/index.html", server)
        end

        return true, "All web servers tested successfully"
      end)
      :build(),

    task("test-connectivity")
      :description("Test connectivity between servers")
      :command(function(this, params)
        -- Testa conectividade entre servidores
        infra_test.can_connect("db-server.internal", 5432)
        infra_test.can_connect("cache-server.internal", 6379)
        infra_test.ping("load-balancer", 5)

        return true, "Connectivity tests passed"
      end)
      :build()
  })
```

### Exemplo 3: Teste de Configuração Completa

```lua
local infra_test = require("infra_test")
local systemd = require("systemd")

workflow.define("deploy-microservice")
  :description("Deploy and validate microservice")
  :version("1.0.0")
  :tasks({
    task("create-service")
      :description("Create systemd service for myapp")
      :command(function(this, params)
        systemd.create_service("myapp", {
          description = "My Application",
          exec_start = "/opt/myapp/bin/start.sh",
          user = "appuser",
          working_directory = "/opt/myapp"
        })

        systemd.enable("myapp")
        systemd.start("myapp")

        return true, "Service created and started successfully"
      end)
      :build(),

    task("validate-deployment")
      :description("Validate microservice deployment")
      :command(function(this, params)
        -- Verifica estrutura de diretórios
        infra_test.is_directory("/opt/myapp")
        infra_test.is_directory("/opt/myapp/bin")
        infra_test.is_directory("/opt/myapp/logs")

        -- Verifica arquivos
        infra_test.is_file("/opt/myapp/bin/start.sh")
        infra_test.file_mode("/opt/myapp/bin/start.sh", "755")
        infra_test.file_owner("/opt/myapp", "appuser")

        -- Verifica serviço
        infra_test.service_is_running("myapp")
        infra_test.service_is_enabled("myapp")

        -- Verifica processo
        infra_test.process_is_running("myapp")

        -- Verifica porta da aplicação
        infra_test.port_is_listening(8080)

        -- Testa endpoint da aplicação
        infra_test.command_succeeds("curl -s http://localhost:8080/health")
        infra_test.command_stdout_contains(
          "curl -s http://localhost:8080/health",
          "\"status\":\"up\""
        )

        return true, "Deployment validated successfully"
      end)
      :build()
  })
  :delegate_to("app-server-prod")
```

### Exemplo 4: Teste de Segurança

```lua
local infra_test = require("infra_test")

workflow.define("security-audit")
  :description("Perform security audit on production server")
  :version("1.0.0")
  :tasks({
    task("check-file-permissions")
      :description("Check critical file permissions")
      :command(function(this, params)
        -- Verifica permissões críticas
        infra_test.file_mode("/etc/passwd", "644")
        infra_test.file_mode("/etc/shadow", "640")
        infra_test.file_mode("/root/.ssh/id_rsa", "600")

        -- Verifica proprietários
        infra_test.file_owner("/etc/shadow", "root")
        infra_test.file_group("/etc/shadow", "shadow")

        return true, "File permissions verified successfully"
      end)
      :build(),

    task("check-services")
      :description("Check service security status")
      :command(function(this, params)
        -- Verifica que serviços desnecessários não estão rodando
        infra_test.command_fails("systemctl is-active telnet")
        infra_test.command_fails("systemctl is-active ftp")

        -- Verifica que serviços críticos estão rodando
        infra_test.service_is_running("sshd")
        infra_test.service_is_running("fail2ban")

        return true, "Service security checks passed"
      end)
      :build(),

    task("check-firewall")
      :description("Check firewall rules")
      :command(function(this, params)
        -- Verifica regras de firewall
        infra_test.command_succeeds("iptables -L | grep -q 'Chain INPUT'")
        infra_test.command_stdout_contains(
          "iptables -L INPUT",
          "ACCEPT.*tcp.*dpt:ssh"
        )

        return true, "Firewall rules verified successfully"
      end)
      :build()
  })
  :delegate_to("prod-server")
```

### Exemplo 5: Teste de Pacotes e Dependências

```lua
local infra_test = require("infra_test")
local pkg = require("pkg")

workflow.define("setup-development-environment")
  :description("Setup and verify development environment")
  :version("1.0.0")
  :tasks({
    task("install-packages")
      :description("Install required development packages")
      :command(function(this, params)
        pkg.install("git")
        pkg.install("docker-ce")
        pkg.install("nodejs")
        pkg.install("python3")

        return true, "Development packages installed successfully"
      end)
      :build(),

    task("verify-packages")
      :description("Verify package installations and versions")
      :command(function(this, params)
        -- Verifica se todos os pacotes foram instalados
        infra_test.package_is_installed("git")
        infra_test.package_is_installed("docker-ce")
        infra_test.package_is_installed("nodejs")
        infra_test.package_is_installed("python3")

        -- Verifica versões específicas
        infra_test.package_version("nodejs", "18")
        infra_test.package_version("python3", "3.10")

        -- Verifica binários disponíveis
        infra_test.command_succeeds("which git")
        infra_test.command_succeeds("which docker")
        infra_test.command_succeeds("which node")
        infra_test.command_succeeds("which python3")

        -- Verifica versões via comando
        infra_test.command_stdout_contains("node --version", "v18")
        infra_test.command_stdout_contains("python3 --version", "Python 3.10")

        return true, "Package verification completed successfully"
      end)
      :build(),

    task("verify-docker-service")
      :description("Verify Docker service is running")
      :command(function(this, params)
        infra_test.service_is_running("docker")
        infra_test.service_is_enabled("docker")
        infra_test.port_is_listening(2375)

        return true, "Docker service verified successfully"
      end)
      :build()
  })
  :delegate_to("dev-machine")
```

### Exemplo 6: Auditoria de Pacotes Multi-Agent

```lua
local infra_test = require("infra_test")

workflow.define("audit-packages")
  :description("Audit packages across multiple servers")
  :version("1.0.0")
  :tasks({
    task("audit-web-servers")
      :description("Audit web server packages")
      :command(function(this, params)
        local servers = {"web-01", "web-02", "web-03"}
        local required_packages = {
          "nginx",
          "certbot",
          "ufw",
          "fail2ban"
        }

        for _, server in ipairs(servers) do
          print("Auditing " .. server)

          for _, pkg_name in ipairs(required_packages) do
            infra_test.package_is_installed(pkg_name, server)
          end

          -- Verifica versão do nginx
          infra_test.package_version("nginx", "1.18", server)

          -- Verifica que pacotes inseguros não estão instalados
          infra_test.command_fails("dpkg -l telnetd", server)
          infra_test.command_fails("dpkg -l rsh-server", server)
        end

        return true, "Web server audit completed successfully"
      end)
      :build(),

    task("audit-database-servers")
      :description("Audit database server packages")
      :command(function(this, params)
        local db_servers = {"db-01", "db-02"}

        for _, server in ipairs(db_servers) do
          print("Auditing database: " .. server)

          -- Verifica pacotes do PostgreSQL
          infra_test.package_is_installed("postgresql-14", server)
          infra_test.package_is_installed("postgresql-contrib", server)

          -- Verifica serviço
          infra_test.service_is_running("postgresql", server)
          infra_test.port_is_listening(5432, server)

          -- Verifica versão
          infra_test.command_stdout_contains(
            "psql --version",
            "14.",
            server
          )
        end

        return true, "Database server audit completed successfully"
      end)
      :build()
  })
```

---

## Melhores Práticas

1. **Organize Testes por Contexto**: Agrupe testes relacionados em tasks separadas
2. **Use Nomes Descritivos**: Nomeie suas tasks de forma clara (ex: "verify-nginx-config")
3. **Teste Progressivamente**: Comece com testes básicos (existência) e avance para testes complexos (conteúdo, permissões)
4. **Teste em Múltiplos Agentes**: Use o parâmetro `target` para validar configurações em vários servidores
5. **Combine com Módulos**: Integre `infra_test` com `pkg`, `systemd`, e outros módulos para ciclos completos de deploy+teste
6. **Valide Pacotes**: Sempre verifique se pacotes foram instalados corretamente após operações de instalação
7. **Use Versões Específicas**: Para ambientes de produção, valide versões específicas de pacotes críticos

## Casos de Uso Recomendados

### 1. Deploy com Validação
Combine instalação de pacotes com validação imediata:
```lua
workflow.define("deploy-with-validation")
  :description("Deploy with immediate validation")
  :version("1.0.0")
  :tasks({
    task("install")
      :description("Install nginx")
      :command(function(this, params)
        pkg.install("nginx")
        return true, "Nginx installed"
      end)
      :build(),
    task("validate")
      :description("Validate installation")
      :command(function(this, params)
        infra_test.package_is_installed("nginx")
        infra_test.service_is_running("nginx")
        infra_test.port_is_listening(80)
        return true, "Validation complete"
      end)
      :build()
  })
```

### 2. Auditoria de Conformidade
Valide que todos os servidores estão em conformidade:
```lua
workflow.define("compliance-check")
  :description("Check security compliance")
  :version("1.0.0")
  :tasks({
    task("check-security-packages")
      :description("Verify security packages")
      :command(function(this, params)
        infra_test.package_is_installed("fail2ban")
        infra_test.package_is_installed("ufw")
        infra_test.service_is_running("fail2ban")
        return true, "Security compliance verified"
      end)
      :build()
  })
```

### 3. Validação de Dependências
Verifique que todas as dependências necessárias estão presentes:
```lua
workflow.define("check-dependencies")
  :description("Check all required dependencies")
  :version("1.0.0")
  :tasks({
    task("verify")
      :description("Verify all dependencies are installed")
      :command(function(this, params)
        local deps = {"python3", "python3-pip", "python3-venv"}
        for _, dep in ipairs(deps) do
          infra_test.package_is_installed(dep)
        end
        return true, "All dependencies verified"
      end)
      :build()
  })
```

## Notas Importantes

- ⚠️ Todos os testes são síncronos e bloqueiam a execução até completarem
- ⚠️ Uma falha em qualquer teste interrompe a task imediatamente
- ⚠️ Testes em agentes remotos requerem que o agente esteja conectado e ativo
- ⚠️ Comandos shell são executados com `sh -c`, portanto use sintaxe POSIX-compatível

## Diferenças com Outras Ferramentas

### vs Testinfra
- ✅ Integrado nativamente ao Sloth Runner (sem Python/pip)
- ✅ Usa o sistema de agents nativo
- ✅ Sintaxe Lua consistente com o resto do workflow

### vs InSpec
- ✅ Mais leve e sem dependências Ruby
- ✅ Integração total com tasks e workflows
- ✅ Execução em tempo real durante o deploy

### vs Serverspec
- ✅ Não requer instalação de gems
- ✅ Melhor performance para testes rápidos
- ✅ Suporte nativo a execução paralela (via goroutines)
