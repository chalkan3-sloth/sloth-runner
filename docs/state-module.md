# 🔄 State Management Module

O módulo **State** do sloth-runner fornece funcionalidades avançadas de gerenciamento de estado persistente com operações atômicas, locks distribuídos e TTL (Time To Live). Todos os dados são armazenados localmente usando SQLite com WAL mode para máxima performance e confiabilidade.

## 📦 Recursos Principais

- **Persistência SQLite**: Armazenamento confiável com WAL mode
- **Operações Atômicas**: Increment, compare-and-swap, append
- **Locks Distribuídos**: Seções críticas com timeout automático
- **TTL (Time To Live)**: Expiração automática de chaves
- **Tipos de Dados**: String, number, boolean, table, list
- **Pattern Matching**: Busca por chaves com wildcards
- **Cleanup Automático**: Limpeza de dados expirados em background
- **Estatísticas**: Métricas de uso e performance

## 🚀 Como Usar

### Operações Básicas

```lua
-- Definir valores
state.set("app_version", "v1.2.3")
state.set("user_count", 1000)
state.set("config", {
    debug = true,
    max_connections = 100
})

-- Recuperar valores
local version = state.get("app_version")
local count = state.get("user_count")
local config = state.get("config")

-- Valor padrão se chave não existir
local theme = state.get("ui_theme", "dark")

-- Verificar existência
if state.exists("app_version") then
    log.info("Versão da app está configurada")
end

-- Deletar chave
state.delete("old_key")
```

### TTL (Time To Live)

```lua
-- Definir com TTL (60 segundos)
state.set("session_token", "abc123", 60)

-- Definir TTL para chave existente
state.set_ttl("user_session", 300) -- 5 minutos

-- Verificar TTL restante
local ttl = state.get_ttl("session_token")
log.info("Token expira em " .. ttl .. " segundos")
```

### Operações Atômicas

```lua
-- Incremento atômico
local counter = state.increment("page_views", 1)
local bulk_counter = state.increment("downloads", 50)

-- Decremento atômico  
local remaining = state.decrement("inventory", 5)

-- Append a string
state.set("log_messages", "Iniciando aplicação")
local new_length = state.append("log_messages", " -> Conectando ao banco")

-- Compare-and-swap atômico
local old_version = state.get("config_version")
local success = state.compare_swap("config_version", old_version, old_version + 1)
if success then
    log.info("Configuração atualizada com segurança")
end
```

### Operações de Lista

```lua
-- Adicionar itens à lista
state.list_push("deployment_queue", {
    app = "frontend",
    version = "v2.1.0",
    environment = "staging"
})

state.list_push("deployment_queue", {
    app = "backend", 
    version = "v1.8.2",
    environment = "production"
})

-- Verificar tamanho da lista
local queue_size = state.list_length("deployment_queue")
log.info("Itens na fila: " .. queue_size)

-- Processar lista (pop remove último item)
while state.list_length("deployment_queue") > 0 do
    local deployment = state.list_pop("deployment_queue")
    log.info("Processando deployment: " .. deployment.app)
    -- Processar deployment...
end
```

### Locks Distribuídos e Seções Críticas

```lua
-- Tentar adquirir lock (sem esperar)
local lock_acquired = state.try_lock("deployment_lock", 30) -- 30 segundos TTL
if lock_acquired then
    -- Trabalho crítico
    state.unlock("deployment_lock")
end

-- Lock com espera e timeout
local acquired = state.lock("database_migration", 60) -- espera até 60s
if acquired then
    -- Executar migração
    state.unlock("database_migration")
end

-- Seção crítica com gerenciamento automático de lock
state.with_lock("critical_section", function()
    log.info("Executando operação crítica...")
    
    -- Atualizar contador global
    local counter = state.increment("global_counter", 1)
    
    -- Atualizar timestamp
    state.set("last_operation", os.time())
    
    log.info("Operação crítica concluída - contador: " .. counter)
    
    -- Lock é liberado automaticamente quando função retorna
    return "operacao_sucesso"
end, 15) -- timeout de 15 segundos
```

### Busca e Limpeza por Padrões

```lua
-- Criar chaves com padrão
state.set("user:1:name", "Alice")
state.set("user:1:email", "alice@example.com") 
state.set("user:2:name", "Bob")
state.set("session:abc123", "user_1_session")
state.set("cache:products:123", product_data)

-- Buscar chaves por padrão
local user_keys = state.keys("user:*")        -- Todas as chaves de usuário
local user1_keys = state.keys("user:1:*")     -- Apenas dados do user 1
local cache_keys = state.keys("cache:*")      -- Todas as chaves de cache

-- Listar todas as chaves
local all_keys = state.keys() -- ou state.keys("*")

-- Limpeza seletiva
state.clear("cache:*")        -- Limpa todo o cache
state.clear("session:*")      -- Limpa todas as sessões
state.clear("temp_*")         -- Limpa chaves temporárias
```

### Monitoramento e Estatísticas

```lua
-- Obter estatísticas do sistema
local stats = state.stats()
log.info("Estatísticas do State:")
log.info("  Total de chaves: " .. stats.total_keys)
log.info("  Chaves expiradas: " .. stats.expired_keys)  
log.info("  Locks ativos: " .. stats.active_locks)
log.info("  Tamanho do DB: " .. stats.db_size_bytes .. " bytes")
log.info("  Caminho do DB: " .. stats.db_path)
```

## 📋 Referência Completa da API

### Operações Básicas
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.set(key, value, ttl?)` | key: string, value: any, ttl?: number | success: boolean | Define um valor com TTL opcional |
| `state.get(key, default?)` | key: string, default?: any | value: any | Recupera um valor ou retorna default |
| `state.delete(key)` | key: string | success: boolean | Remove uma chave |
| `state.exists(key)` | key: string | exists: boolean | Verifica se chave existe |
| `state.clear(pattern?)` | pattern?: string | success: boolean | Remove chaves por padrão |

### TTL Operations
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.set_ttl(key, seconds)` | key: string, seconds: number | success: boolean | Define TTL para chave existente |
| `state.get_ttl(key)` | key: string | ttl: number | Retorna TTL restante (-1 = sem TTL, -2 = não existe) |

### Operações Atômicas
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.increment(key, delta?)` | key: string, delta?: number | new_value: number | Incrementa valor atomicamente |
| `state.decrement(key, delta?)` | key: string, delta?: number | new_value: number | Decrementa valor atomicamente |
| `state.append(key, value)` | key: string, value: string | new_length: number | Anexa string atomicamente |
| `state.compare_swap(key, old, new)` | key: string, old: any, new: any | success: boolean | Compare-and-swap atômico |

### Operações de Lista
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.list_push(key, item)` | key: string, item: any | length: number | Adiciona item ao final da lista |
| `state.list_pop(key)` | key: string | item: any \| nil | Remove e retorna último item |
| `state.list_length(key)` | key: string | length: number | Retorna tamanho da lista |

### Locks Distribuídos
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.try_lock(name, ttl)` | name: string, ttl: number | success: boolean | Tenta adquirir lock sem esperar |
| `state.lock(name, timeout?)` | name: string, timeout?: number | success: boolean | Adquire lock com timeout |
| `state.unlock(name)` | name: string | success: boolean | Libera lock |
| `state.with_lock(name, fn, timeout?)` | name: string, fn: function, timeout?: number | result: any | Executa função com lock automático |

### Utilitários
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.keys(pattern?)` | pattern?: string | keys: table | Lista chaves por padrão |
| `state.stats()` | - | stats: table | Estatísticas do sistema |

## 🏗️ Casos de Uso Práticos

### 1. Controle de Versão de Deploy

```lua
Modern DSLs = {
    deployment_pipeline = {
        tasks = {
            prepare_deploy = {
                command = function()
                    -- Verificar última versão deployada
                    local last_version = state.get("last_deployed_version", "v0.0.0")
                    local new_version = "v1.2.3"
                    
                    -- Verificar se já está deployado
                    if last_version == new_version then
                        log.warn("Versão " .. new_version .. " já deployada")
                        return false, "Version already deployed"
                    end
                    
                    -- Registrar início do deploy
                    state.set("deploy_status", "in_progress")
                    state.set("deploy_start_time", os.time())
                    state.increment("total_deploys", 1)
                    
                    return true, "Deploy preparation completed"
                end
            },
            
            execute_deploy = {
                depends_on = "prepare_deploy",
                command = function()
                    -- Seção crítica para deploy
                    return state.with_lock("deployment_lock", function()
                        log.info("Executando deploy com lock...")
                        
                        -- Simular deploy
                        exec.run("sleep 5")
                        
                        -- Atualizar estado
                        state.set("last_deployed_version", "v1.2.3")
                        state.set("deploy_status", "completed")
                        state.set("deploy_end_time", os.time())
                        
                        -- Registrar histórico
                        state.list_push("deploy_history", {
                            version = "v1.2.3",
                            timestamp = os.time(),
                            duration = state.get("deploy_end_time") - state.get("deploy_start_time")
                        })
                        
                        return true, "Deploy completed successfully"
                    end, 300) -- 5 minutos timeout
                end
            }
        }
    }
}
```

### 2. Cache com TTL Automático

```lua
-- Função helper para cache
function get_cached_data(cache_key, fetch_function, ttl)
    local cached = state.get(cache_key)
    if cached then
        log.info("Cache hit: " .. cache_key)
        return cached
    end
    
    log.info("Cache miss: " .. cache_key .. " - fetching...")
    local data = fetch_function()
    state.set(cache_key, data, ttl or 300) -- 5 minutos default
    return data
end

-- Uso em tasks
Modern DSLs = {
    data_processing = {
        tasks = {
            fetch_user_data = {
                command = function()
                    local user_data = get_cached_data("user:123:profile", function()
                        -- Simulação de busca custosa
                        return {
                            name = "Alice",
                            email = "alice@example.com",
                            preferences = {"dark_mode", "notifications"}
                        }
                    end, 600) -- Cache por 10 minutos
                    
                    log.info("User data: " .. data.to_json(user_data))
                    return true, "User data retrieved"
                end
            }
        }
    }
}
```

### 3. Controle de Rate Limiting

```lua
function check_rate_limit(identifier, max_requests, window_seconds)
    local key = "rate_limit:" .. identifier
    local current_count = state.get(key, 0)
    
    if current_count >= max_requests then
        return false, "Rate limit exceeded"
    end
    
    -- Incrementar contador
    if current_count == 0 then
        -- Primeira requisição na janela
        state.set(key, 1, window_seconds)
    else
        -- Incrementar contador existente
        state.increment(key, 1)
    end
    
    return true, "Request allowed"
end

-- Uso em tasks
Modern DSLs = {
    api_tasks = {
        tasks = {
            make_api_call = {
                command = function()
                    local allowed, msg = check_rate_limit("api_calls", 100, 3600) -- 100 calls/hora
                    
                    if not allowed then
                        log.error(msg)
                        return false, msg
                    end
                    
                    -- Fazer chamada da API
                    log.info("Making API call...")
                    return true, "API call completed"
                end
            }
        }
    }
}
```

## 🔧 Configuração e Personalização

### Localização do Banco de Dados

Por padrão, o banco de dados SQLite é criado em:
- **Linux/macOS**: `~/.sloth-runner/state.db`
- **Windows**: `%USERPROFILE%\.sloth-runner\state.db`

### Características Técnicas

- **Engine**: SQLite 3 com WAL mode
- **Concurrent Access**: Suporte a múltiplas conexões simultâneas
- **Auto-cleanup**: Limpeza automática de dados expirados a cada 5 minutos
- **Lock Timeout**: Locks expirados são limpos automaticamente
- **Serialização**: JSON para objetos complexos, formato nativo para tipos simples

### Limitações

- **Escopo Local**: Estado é persistido apenas na máquina local
- **Concorrência**: Locks são efetivos apenas no processo local
- **Tamanho**: Adequado para datasets pequenos a médios (< 1GB)

## 🎯 Próximos Passos

Para evoluir o módulo state, considere implementar:

1. **State Distribuído**: Sincronização entre múltiplos agentes
2. **Backup/Restore**: Funcionalidades de backup automático
3. **Compressão**: Compressão de dados grandes
4. **Indices Customizados**: Índices personalizados para queries complexas
5. **Webhooks**: Notificações em mudanças de estado
6. **Métricas Avançadas**: Histogramas de performance e uso

O módulo **State** transforma o sloth-runner em uma plataforma ainda mais poderosa para orquestração de tarefas complexas com gerenciamento de estado robusto e confiável! 🚀