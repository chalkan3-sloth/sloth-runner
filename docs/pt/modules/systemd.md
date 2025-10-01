# ⚙️ Módulo Systemd

O módulo `systemd` fornece funcionalidades abrangentes de gerenciamento de serviços systemd para sistemas Linux. Permite criar, gerenciar e monitorar serviços systemd programaticamente.

## 🎯 Visão Geral

O módulo systemd permite:
- Criar e configurar arquivos de serviço systemd
- Iniciar, parar, reiniciar e recarregar serviços
- Habilitar e desabilitar serviços
- Verificar status e atividade de serviços
- Listar todos os serviços
- Gerenciar configuração do daemon systemd

## 📚 Visão Geral das Funções

| Função | Descrição |
|--------|-----------|
| `systemd.create_service(nome, config)` | Criar novo serviço systemd |
| `systemd.start(servico)` | Iniciar um serviço |
| `systemd.stop(servico)` | Parar um serviço |
| `systemd.restart(servico)` | Reiniciar um serviço |
| `systemd.reload(servico)` | Recarregar um serviço |
| `systemd.enable(servico)` | Habilitar no boot |
| `systemd.disable(servico)` | Desabilitar do boot |
| `systemd.status(servico)` | Obter status do serviço |
| `systemd.is_active(servico)` | Verificar se está ativo |
| `systemd.is_enabled(servico)` | Verificar se está habilitado |
| `systemd.daemon_reload()` | Recarregar daemon |
| `systemd.remove_service(servico)` | Remover um serviço |
| `systemd.list_services(opts)` | Listar todos os serviços |
| `systemd.show(servico)` | Mostrar info detalhada |

## 📖 Documentação Detalhada

### Criação de Serviços

#### `systemd.create_service(nome, config)`

Cria um novo arquivo de serviço systemd em `/etc/systemd/system/{nome}.service`.

**Parâmetros:**
- `nome` (string): Nome do serviço (sem extensão .service)
- `config` (table): Configuração do serviço

**Opções de Configuração:**

```lua
{
    -- Seção [Unit]
    description = "Descrição do serviço",
    after = "network.target",
    wants = "outro.service",
    requires = "necessario.service",
    
    -- Seção [Service] (obrigatório)
    exec_start = "/caminho/para/executavel",
    exec_stop = "/caminho/para/script/parada",
    exec_reload = "/caminho/para/script/recarga",
    type = "simple",  -- simple, forking, oneshot, dbus, notify, idle
    user = "usuario",
    group = "grupo",
    working_directory = "/caminho/para/diretorio",
    restart = "always",  -- no, on-success, on-failure, on-abnormal, on-abort, always
    restart_sec = "5s",
    environment = {
        VAR1 = "valor1",
        VAR2 = "valor2"
    },
    
    -- Seção [Install]
    wanted_by = "multi-user.target"
}
```

**Retorna:**
- `sucesso` (boolean): `true` se o serviço foi criado
- `mensagem` (string): Mensagem de resultado

**Exemplos:**

=== "DSL Moderno"
    ```lua
    local systemd = require("systemd")
    
    local criar_servico_web = task("criar_servico_web")
        :description("Criar serviço de aplicação web")
        :command(function(this, params)
            log.info("Criando serviço web...")
            
            local config = {
                description = "Servidor de Aplicação Web",
                after = "network.target",
                exec_start = "/usr/bin/node /app/server.js",
                type = "simple",
                user = "webapp",
                working_directory = "/app",
                restart = "always",
                restart_sec = "10s",
                environment = {
                    NODE_ENV = "production",
                    PORT = "3000"
                }
            }
            
            local sucesso, msg = systemd.create_service("webapp", config)
            
            if sucesso then
                log.info("✅ Serviço criado!")
                systemd.daemon_reload()
                systemd.enable("webapp")
                systemd.start("webapp")
                return true, "Serviço implantado"
            else
                log.error("❌ Falha: " .. msg)
                return false, msg
            end
        end)
        :timeout("60s")
        :build()
    
    workflow.define("implantar_servico")
        :tasks({ criar_servico_web })
    ```

=== "Com delegate_to"
    ```lua
    local systemd = require("systemd")
    
    local implantar_servico_remoto = task("implantar_servico_remoto")
        :description("Implantar serviço no agente remoto")
        :command(function(this, params)
            local config = {
                description = "Agente de Monitoramento Remoto",
                after = "network.target",
                exec_start = "/opt/monitor/agent",
                type = "simple",
                user = "monitor",
                restart = "always"
            }
            
            local sucesso, msg = systemd.create_service("monitor-agent", config)
            
            if sucesso then
                systemd.daemon_reload()
                systemd.enable("monitor-agent")
                systemd.start("monitor-agent")
                log.info("✅ Implantado em " .. (this.agent or "local"))
                return true, "OK"
            end
            
            return false, "Falhou"
        end)
        :delegate_to("servidor-producao")
        :timeout("60s")
        :build()
    
    workflow.define("implantacao_remota")
        :tasks({ implantar_servico_remoto })
    ```

### Controle de Serviços

#### `systemd.start(servico)`

Inicia um serviço systemd.

**Exemplo:**
```lua
local sucesso, saida = systemd.start("nginx")
if sucesso then
    log.info("✅ Nginx iniciado")
end
```

#### `systemd.stop(servico)`

Para um serviço systemd.

#### `systemd.restart(servico)`

Reinicia um serviço systemd.

#### `systemd.reload(servico)`

Recarrega configuração sem reiniciar.

### Status do Serviço

#### `systemd.status(servico)`

Obtém status detalhado de um serviço.

**Exemplo:**
```lua
local status, err = systemd.status("nginx")
log.info("Status:\n" .. status)
```

#### `systemd.is_active(servico)`

Verifica se um serviço está ativo/rodando.

**Retorna:**
- `ativo` (boolean): `true` se ativo
- `estado` (string): Estado do serviço

**Exemplo:**
```lua
local ativo, estado = systemd.is_active("nginx")
if ativo then
    log.info("✅ Serviço está rodando")
else
    log.warn("❌ Serviço está " .. estado)
end
```

#### `systemd.is_enabled(servico)`

Verifica se um serviço está habilitado para iniciar no boot.

**Exemplo:**
```lua
local habilitado, estado = systemd.is_enabled("nginx")
```

### Gerenciamento de Serviços

#### `systemd.enable(servico)`

Habilita um serviço para iniciar automaticamente no boot.

#### `systemd.disable(servico)`

Desabilita um serviço do boot.

#### `systemd.daemon_reload()`

Recarrega configuração do daemon. Necessário após criar ou modificar arquivos de serviço.

#### `systemd.remove_service(servico)`

Remove completamente um serviço (para, desabilita e deleta o arquivo).

### Informações de Serviços

#### `systemd.list_services(opcoes)`

Lista serviços systemd com filtros opcionais.

**Parâmetros:**
- `opcoes` (table, opcional): Opções de filtro
  - `state`: Filtrar por estado (ex: "active", "failed", "inactive")
  - `no_header`: Boolean, excluir cabeçalho

**Exemplo:**
```lua
-- Listar todos os serviços
local lista, err = systemd.list_services()

-- Listar apenas ativos
local ativos, err = systemd.list_services({ state = "active" })

-- Listar falhados sem cabeçalho
local falhados, err = systemd.list_services({ 
    state = "failed", 
    no_header = true 
})
```

#### `systemd.show(servico)`

Mostra propriedades detalhadas de um serviço.

## 🎯 Exemplos Completos

### Implantação de Aplicação Web

```lua
local systemd = require("systemd")

local implantar_webapp = task("implantar_webapp")
    :description("Implantar e configurar aplicação web")
    :command(function(this, params)
        log.info("🚀 Implantando aplicação web...")
        
        local config = {
            description = "Aplicação Web Node.js",
            after = "network.target postgresql.service",
            requires = "postgresql.service",
            exec_start = "/usr/bin/node /var/www/app/server.js",
            type = "simple",
            user = "webapp",
            working_directory = "/var/www/app",
            restart = "always",
            environment = {
                NODE_ENV = "production",
                PORT = "3000"
            }
        }
        
        local sucesso, msg = systemd.create_service("webapp", config)
        if not sucesso then
            return false, "Falha ao criar serviço: " .. msg
        end
        
        systemd.daemon_reload()
        systemd.enable("webapp")
        systemd.start("webapp")
        
        local ativo, estado = systemd.is_active("webapp")
        if ativo then
            log.info("✅ Serviço rodando!")
            return true, "Implantação bem-sucedida"
        else
            return false, "Serviço não iniciou"
        end
    end)
    :timeout("120s")
    :build()

workflow.define("implantar")
    :tasks({ implantar_webapp })
```

### Verificação de Saúde

```lua
local systemd = require("systemd")

local verificacao_saude = task("verificacao_saude")
    :description("Verificar saúde dos serviços críticos")
    :command(function(this, params)
        log.info("🔍 Verificação de Saúde...")
        
        local servicos = {"nginx", "postgresql", "redis"}
        local todos_saudaveis = true
        
        for _, servico in ipairs(servicos) do
            local ativo, estado = systemd.is_active(servico)
            
            log.info("\n📦 " .. servico .. ":")
            log.info("  Ativo: " .. (ativo and "✅ SIM" or "❌ NÃO"))
            
            if not ativo then
                todos_saudaveis = false
            end
        end
        
        if todos_saudaveis then
            return true, "Todos OK"
        else
            return false, "Serviços inoperantes"
        end
    end)
    :timeout("60s")
    :build()

workflow.define("verificar_saude")
    :tasks({ verificacao_saude })
```

### Gerenciamento Distribuído

```lua
local systemd = require("systemd")

local reiniciar_todos_servidores = task("reiniciar_nginx")
    :description("Reiniciar nginx em todos os servidores")
    :command(function(this, params)
        log.info("🔄 Reiniciando nginx...")
        
        local sucesso, saida = systemd.restart("nginx")
        
        if sucesso then
            local ativo, estado = systemd.is_active("nginx")
            if ativo then
                log.info("✅ Nginx reiniciado")
                return true, "OK"
            end
        end
        
        return false, "Falha"
    end)
    :delegate_to("servidor-web-1")
    :timeout("60s")
    :build()

workflow.define("reinicio_escalonado")
    :tasks({ reiniciar_todos_servidores })
```

## 🚀 Melhores Práticas

1. **Sempre recarregue daemon após criar/modificar:**
   ```lua
   systemd.create_service("meuservico", config)
   systemd.daemon_reload()
   ```

2. **Verifique se iniciou com sucesso:**
   ```lua
   systemd.start("meuservico")
   local ativo, estado = systemd.is_active("meuservico")
   ```

3. **Habilite serviços para persistência:**
   ```lua
   systemd.enable("meuservico")
   ```

4. **Use tipos de serviço apropriados:**
   - `simple`: Padrão, processo não faz fork
   - `forking`: Processo faz fork e pai sai
   - `oneshot`: Processo sai antes do systemd continuar

5. **Configure políticas de reinício:**
   ```lua
   restart = "always"
   restart_sec = "10s"
   ```

## ⚠️ Considerações de Segurança

- Arquivos criados em `/etc/systemd/system/` (requer root/sudo)
- Sempre especifique `user` e `group`
- Use `WorkingDirectory` para isolar ambiente
- Valide variáveis de ambiente
- Use permissões apropriadas (0644)

## 🐧 Suporte de Plataforma

- **Linux**: ✅ Suporte completo
- **Ubuntu/Debian**: ✅ Suportado
- **CentOS/RHEL**: ✅ Suportado
- **Fedora**: ✅ Suportado
- **Arch Linux**: ✅ Suportado
- **macOS**: ❌ Não suportado
- **Windows**: ❌ Não suportado

## 🔗 Veja Também

- [Módulo exec](exec.md)
- [Guia DSL Moderno](../modern-dsl/overview.md)
- [Agentes Distribuídos](../distributed.md)
