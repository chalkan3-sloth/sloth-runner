# Módulo infra_test

O módulo `infra_test` fornece um framework completo de validação e teste de infraestrutura nativo ao Sloth Runner. Ele permite que você insira asserções de teste diretamente nas suas tasks para verificar o resultado de operações de deploy ou configuration management.

## Visão Geral

O `infra_test` é inspirado em ferramentas como Testinfra e InSpec, mas é nativo e integrado ao Sloth Runner, permitindo testes de infraestrutura diretamente nas tasks sem dependências externas.

## Características Principais

- ✅ **Execução Local e Remota**: Todos os testes podem ser executados localmente ou delegados para agentes remotos
- ✅ **Asserções Nativas**: Interrompe a execução da task em caso de falha
- ✅ **Zero Dependências**: Não requer instalação de ferramentas externas
- ✅ **Integração Total**: Funciona perfeitamente com o sistema de agents do Sloth Runner

## Parâmetro Target

Todas as funções de teste aceitam um parâmetro opcional `target` para especificar onde o teste será executado:

| Parâmetro target | Comportamento |
|-----------------|---------------|
| Omitido ou `"local"` | Executa no agente local (onde a task está rodando) |
| String (nome do agente) | O teste é delegado ao agente remoto especificado |
| `"localhost"` | Força o teste no agente onde a task foi agendada |

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

workflow("test-deployment")
  :task("verify-config", function()
    -- Verifica localmente
    infra_test.file_exists("/etc/nginx/nginx.conf")
    
    -- Verifica em agente remoto
    infra_test.file_exists("/etc/nginx/nginx.conf", "web-server-01")
  end)
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

## Exemplos Completos

### Exemplo 1: Teste de Deploy de Aplicação

```lua
local infra_test = require("infra_test")
local pkg = require("pkg")

workflow("deploy-and-test-app")
  :task("install-nginx", function()
    pkg.install("nginx")
  end)
  
  :task("verify-installation", function()
    -- Verifica se o nginx foi instalado
    infra_test.file_exists("/usr/sbin/nginx")
    infra_test.file_exists("/etc/nginx/nginx.conf")
    
    -- Verifica se o serviço está rodando e habilitado
    infra_test.service_is_running("nginx")
    infra_test.service_is_enabled("nginx")
    
    -- Verifica se a porta está aberta
    infra_test.port_is_tcp(80)
    
    -- Verifica se o processo está ativo
    infra_test.process_is_running("nginx")
  end)
  
  :task("verify-config", function()
    -- Verifica permissões e proprietário
    infra_test.file_mode("/etc/nginx/nginx.conf", "644")
    infra_test.file_owner("/var/www/html", "www-data")
    
    -- Verifica conteúdo da configuração
    infra_test.file_contains("/etc/nginx/nginx.conf", "worker_processes")
  end)
  
  :delegate_to("web-server-01")
```

### Exemplo 2: Validação Multi-Agent

```lua
local infra_test = require("infra_test")

workflow("test-infrastructure")
  :task("test-web-servers", function()
    -- Testa múltiplos servidores web
    local servers = {"web-01", "web-02", "web-03"}
    
    for _, server in ipairs(servers) do
      print("Testing " .. server)
      
      infra_test.service_is_running("nginx", server)
      infra_test.port_is_listening(80, server)
      infra_test.port_is_listening(443, server)
      infra_test.file_exists("/var/www/html/index.html", server)
    end
  end)
  
  :task("test-connectivity", function()
    -- Testa conectividade entre servidores
    infra_test.can_connect("db-server.internal", 5432)
    infra_test.can_connect("cache-server.internal", 6379)
    infra_test.ping("load-balancer", 5)
  end)
```

### Exemplo 3: Teste de Configuração Completa

```lua
local infra_test = require("infra_test")
local systemd = require("systemd")

workflow("deploy-microservice")
  :task("create-service", function()
    systemd.create_service("myapp", {
      description = "My Application",
      exec_start = "/opt/myapp/bin/start.sh",
      user = "appuser",
      working_directory = "/opt/myapp"
    })
    
    systemd.enable("myapp")
    systemd.start("myapp")
  end)
  
  :task("validate-deployment", function()
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
  end)
  
  :delegate_to("app-server-prod")
```

### Exemplo 4: Teste de Segurança

```lua
local infra_test = require("infra_test")

workflow("security-audit")
  :task("check-file-permissions", function()
    -- Verifica permissões críticas
    infra_test.file_mode("/etc/passwd", "644")
    infra_test.file_mode("/etc/shadow", "640")
    infra_test.file_mode("/root/.ssh/id_rsa", "600")
    
    -- Verifica proprietários
    infra_test.file_owner("/etc/shadow", "root")
    infra_test.file_group("/etc/shadow", "shadow")
  end)
  
  :task("check-services", function()
    -- Verifica que serviços desnecessários não estão rodando
    infra_test.command_fails("systemctl is-active telnet")
    infra_test.command_fails("systemctl is-active ftp")
    
    -- Verifica que serviços críticos estão rodando
    infra_test.service_is_running("sshd")
    infra_test.service_is_running("fail2ban")
  end)
  
  :task("check-firewall", function()
    -- Verifica regras de firewall
    infra_test.command_succeeds("iptables -L | grep -q 'Chain INPUT'")
    infra_test.command_stdout_contains(
      "iptables -L INPUT",
      "ACCEPT.*tcp.*dpt:ssh"
    )
  end)
  
  :delegate_to("prod-server")
```

---

## Melhores Práticas

1. **Organize Testes por Contexto**: Agrupe testes relacionados em tasks separadas
2. **Use Nomes Descritivos**: Nomeie suas tasks de forma clara (ex: "verify-nginx-config")
3. **Teste Progressivamente**: Comece com testes básicos (existência) e avance para testes complexos (conteúdo, permissões)
4. **Teste em Múltiplos Agentes**: Use o parâmetro `target` para validar configurações em vários servidores
5. **Combine com Módulos**: Integre `infra_test` com `pkg`, `systemd`, e outros módulos para ciclos completos de deploy+teste

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
