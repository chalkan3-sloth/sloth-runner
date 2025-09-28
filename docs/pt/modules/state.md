# 💾 Módulo de Gerenciamento de Estado

O módulo **Gerenciamento de Estado** fornece capacidades poderosas de estado persistente com operações atômicas, locks distribuídos e funcionalidade TTL (Time To Live). Todos os dados são armazenados localmente usando SQLite com modo WAL para máxima performance e confiabilidade.

## 🚀 Recursos Principais

- **Persistência SQLite**: Armazenamento confiável com modo WAL
- **Operações Atômicas**: Increment, compare-and-swap, append thread-safe
- **Locks Distribuídos**: Seções críticas com timeout automático
- **TTL (Time To Live)**: Expiração automática de chaves
- **Tipos de Dados**: String, number, boolean, table, list
- **Pattern Matching**: Busca de chaves com wildcards
- **Limpeza Automática**: Cleanup em background de dados expirados
- **Estatísticas**: Métricas de uso e performance

## 📋 Uso Básico

### Definindo e Obtendo Valores

```lua
-- Definir valores
state.set("versao_app", "v1.2.3")
state.set("contador_usuarios", 1000)
state.set("configuracao", {
    debug = true,
    max_conexoes = 100
})

-- Obter valores
local versao = state.get("versao_app")
local contador = state.get("contador_usuarios")
local config = state.get("configuracao")

-- Obter com valor padrão
local tema = state.get("tema_ui", "escuro")

-- Verificar existência
if state.exists("versao_app") then
    log.info("Versão da aplicação está configurada")
end

-- Deletar chave
state.delete("chave_antiga")
```

### TTL (Time To Live)

```lua
-- Definir com TTL (60 segundos)
state.set("token_sessao", "abc123", 60)

-- Definir TTL para chave existente
state.set_ttl("sessao_usuario", 300) -- 5 minutos

-- Verificar TTL restante
local ttl = state.get_ttl("token_sessao")
log.info("Token expira em " .. ttl .. " segundos")
```

### Operações Atômicas

```lua
-- Incremento atômico
local contador = state.increment("visualizacoes_pagina", 1)
local contador_bulk = state.increment("downloads", 50)

-- Decremento atômico  
local restante = state.decrement("estoque", 5)

-- Append de string
state.set("mensagens_log", "Iniciando aplicação")
local novo_tamanho = state.append("mensagens_log", " -> Conectando ao banco")

-- Compare-and-swap atômico
local versao_antiga = state.get("versao_config")
local sucesso = state.compare_swap("versao_config", versao_antiga, versao_antiga + 1)
if sucesso then
    log.info("Configuração atualizada com segurança")
end
```

### Operações de Lista

```lua
-- Adicionar itens à lista
state.list_push("fila_deployment", {
    app = "frontend",
    versao = "v2.1.0",
    ambiente = "staging"
})

-- Verificar tamanho da lista
local tamanho_fila = state.list_length("fila_deployment")
log.info("Itens na fila: " .. tamanho_fila)

-- Processar lista (pop remove último item)
while state.list_length("fila_deployment") > 0 do
    local deployment = state.list_pop("fila_deployment")
    log.info("Processando deployment: " .. deployment.app)
    -- Processar deployment...
end
```

### Locks Distribuídos e Seções Críticas

```lua
-- Tentar adquirir lock (sem esperar)
local lock_adquirido = state.try_lock("lock_deployment", 30) -- 30 segundos TTL
if lock_adquirido then
    -- Trabalho crítico
    state.unlock("lock_deployment")
end

-- Lock com espera e timeout
local adquirido = state.lock("migracao_banco", 60) -- esperar até 60s
if adquirido then
    -- Executar migração
    state.unlock("migracao_banco")
end

-- Seção crítica com gerenciamento automático de lock
state.with_lock("secao_critica", function()
    log.info("Executando operação crítica...")
    
    -- Atualizar contador global
    local contador = state.increment("contador_global", 1)
    
    -- Atualizar timestamp
    state.set("ultima_operacao", os.time())
    
    log.info("Operação crítica concluída - contador: " .. contador)
    
    -- Lock é liberado automaticamente quando a função retorna
    return "operacao_sucesso"
end, 15) -- timeout de 15 segundos
```

## 🔍 Referência da API

### Operações Básicas
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.set(chave, valor, ttl?)` | chave: string, valor: any, ttl?: number | sucesso: boolean | Define um valor com TTL opcional |
| `state.get(chave, padrao?)` | chave: string, padrao?: any | valor: any | Obtém um valor ou retorna o padrão |
| `state.delete(chave)` | chave: string | sucesso: boolean | Remove uma chave |
| `state.exists(chave)` | chave: string | existe: boolean | Verifica se a chave existe |
| `state.clear(padrao?)` | padrao?: string | sucesso: boolean | Remove chaves por padrão |

### Operações TTL
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.set_ttl(chave, segundos)` | chave: string, segundos: number | sucesso: boolean | Define TTL para chave existente |
| `state.get_ttl(chave)` | chave: string | ttl: number | Obtém TTL restante (-1 = sem TTL, -2 = não existe) |

### Operações Atômicas
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.increment(chave, delta?)` | chave: string, delta?: number | novo_valor: number | Incrementa valor atomicamente |
| `state.decrement(chave, delta?)` | chave: string, delta?: number | novo_valor: number | Decrementa valor atomicamente |
| `state.append(chave, valor)` | chave: string, valor: string | novo_tamanho: number | Anexa string atomicamente |
| `state.compare_swap(chave, antigo, novo)` | chave: string, antigo: any, novo: any | sucesso: boolean | Compare-and-swap atômico |

### Operações de Lista
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.list_push(chave, item)` | chave: string, item: any | tamanho: number | Adiciona item ao final da lista |
| `state.list_pop(chave)` | chave: string | item: any \| nil | Remove e retorna último item |
| `state.list_length(chave)` | chave: string | tamanho: number | Obtém tamanho da lista |

### Locks Distribuídos
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.try_lock(nome, ttl)` | nome: string, ttl: number | sucesso: boolean | Tenta adquirir lock sem esperar |
| `state.lock(nome, timeout?)` | nome: string, timeout?: number | sucesso: boolean | Adquire lock com timeout |
| `state.unlock(nome)` | nome: string | sucesso: boolean | Libera lock |
| `state.with_lock(nome, funcao, timeout?)` | nome: string, funcao: function, timeout?: number | resultado: any | Executa função com lock automático |

### Utilitários
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `state.keys(padrao?)` | padrao?: string | chaves: table | Lista chaves por padrão |
| `state.stats()` | - | stats: table | Obtém estatísticas do sistema |

## 💡 Casos de Uso Práticos

### 1. Controle de Versão de Deploy

```lua
TaskDefinitions = {
    pipeline_deployment = {
        tasks = {
            preparar_deploy = {
                command = function()
                    -- Verificar última versão deployada
                    local ultima_versao = state.get("ultima_versao_deployada", "v0.0.0")
                    local nova_versao = "v1.2.3"
                    
                    -- Verificar se já foi deployada
                    if ultima_versao == nova_versao then
                        log.warn("Versão " .. nova_versao .. " já foi deployada")
                        return false, "Versão já deployada"
                    end
                    
                    -- Registrar início do deploy
                    state.set("status_deploy", "em_progresso")
                    state.set("inicio_deploy", os.time())
                    state.increment("total_deploys", 1)
                    
                    return true, "Preparação do deploy concluída"
                end
            },
            
            executar_deploy = {
                depends_on = "preparar_deploy",
                command = function()
                    -- Seção crítica para deployment
                    return state.with_lock("lock_deployment", function()
                        log.info("Executando deployment com lock...")
                        
                        -- Simular deployment
                        exec.run("sleep 5")
                        
                        -- Atualizar estado
                        state.set("ultima_versao_deployada", "v1.2.3")
                        state.set("status_deploy", "concluido")
                        state.set("fim_deploy", os.time())
                        
                        -- Registrar histórico
                        state.list_push("historico_deploy", {
                            versao = "v1.2.3",
                            timestamp = os.time(),
                            duracao = state.get("fim_deploy") - state.get("inicio_deploy")
                        })
                        
                        return true, "Deploy concluído com sucesso"
                    end, 300) -- timeout de 5 minutos
                end
            }
        }
    }
}
```

### 2. Cache Inteligente com TTL

```lua
-- Função helper para cache
function obter_dados_cached(chave_cache, funcao_busca, ttl)
    local cached = state.get(chave_cache)
    if cached then
        log.info("Cache hit: " .. chave_cache)
        return cached
    end
    
    log.info("Cache miss: " .. chave_cache .. " - buscando...")
    local dados = funcao_busca()
    state.set(chave_cache, dados, ttl or 300) -- 5 minutos padrão
    return dados
end

-- Uso em tasks
TaskDefinitions = {
    processamento_dados = {
        tasks = {
            buscar_dados_usuario = {
                command = function()
                    local dados_usuario = obter_dados_cached("usuario:123:perfil", function()
                        -- Simular busca custosa
                        return {
                            nome = "Alice",
                            email = "alice@exemplo.com",
                            preferencias = {"modo_escuro", "notificacoes"}
                        }
                    end, 600) -- Cache por 10 minutos
                    
                    log.info("Dados do usuário: " .. data.to_json(dados_usuario))
                    return true, "Dados do usuário obtidos"
                end
            }
        }
    }
}
```

### 3. Rate Limiting

```lua
function verificar_rate_limit(identificador, max_requisicoes, janela_segundos)
    local chave = "rate_limit:" .. identificador
    local contador_atual = state.get(chave, 0)
    
    if contador_atual >= max_requisicoes then
        return false, "Rate limit excedido"
    end
    
    -- Incrementar contador
    if contador_atual == 0 then
        -- Primeira requisição na janela
        state.set(chave, 1, janela_segundos)
    else
        -- Incrementar contador existente
        state.increment(chave, 1)
    end
    
    return true, "Requisição permitida"
end

-- Uso em tasks
TaskDefinitions = {
    tarefas_api = {
        tasks = {
            fazer_chamada_api = {
                command = function()
                    local permitido, msg = verificar_rate_limit("chamadas_api", 100, 3600) -- 100 chamadas/hora
                    
                    if not permitido then
                        log.error(msg)
                        return false, msg
                    end
                    
                    -- Fazer chamada da API
                    log.info("Fazendo chamada da API...")
                    return true, "Chamada da API concluída"
                end
            }
        }
    }
}
```

## ⚙️ Configuração e Armazenamento

### Localização do Banco de Dados

Por padrão, o banco de dados SQLite é criado em:
- **Linux/macOS**: `~/.sloth-runner/state.db`
- **Windows**: `%USERPROFILE%\.sloth-runner\state.db`

### Características Técnicas

- **Engine**: SQLite 3 com modo WAL
- **Acesso Concorrente**: Suporte a múltiplas conexões simultâneas
- **Limpeza Automática**: Limpeza automática de dados expirados a cada 5 minutos
- **Timeout de Lock**: Locks expirados são limpos automaticamente
- **Serialização**: JSON para objetos complexos, formato nativo para tipos simples

### Limitações

- **Escopo Local**: Estado é persistido apenas na máquina local
- **Concorrência**: Locks são efetivos apenas no processo local
- **Tamanho**: Adequado para datasets pequenos a médios (< 1GB)

## 🔄 Melhores Práticas

1. **Use TTL para dados temporários** para evitar crescimento desnecessário
2. **Use locks para seções críticas** para evitar condições de corrida
3. **Use padrões para operações em lote** para gerenciar chaves relacionadas
4. **Monitore o tamanho do armazenamento** usando `state.stats()`
5. **Use operações atômicas** em vez de padrões read-modify-write
6. **Limpe chaves expiradas** regularmente com `state.clear(padrao)`

O módulo **Gerenciamento de Estado** transforma o sloth-runner em uma plataforma stateful e confiável para orquestração complexa de tarefas! 🚀