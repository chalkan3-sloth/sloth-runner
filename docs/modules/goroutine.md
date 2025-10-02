# 🚀 Módulo Goroutine - Execução Paralela Poderosa

## 🌟 Visão Geral

O módulo `goroutine` traz o **poder das goroutines do Go para seus scripts Lua**, permitindo executar tarefas em paralelo com facilidade. Com este módulo, você pode:

- ⚡ **Executar múltiplas operações simultaneamente** - Reduzir tempo de execução de minutos para segundos
- 🏭 **Criar worker pools** - Controlar concorrência e processar grandes volumes de dados
- 🎯 **Async/Await pattern** - Escrever código assíncrono de forma limpa e legível
- 🔄 **WaitGroups** - Sincronizar múltiplas goroutines facilmente
- ⏱️ **Timeout e error handling** - Executar operações com limites de tempo

### 💼 Casos de Uso Reais

| Cenário | Tempo Sequencial | Com Goroutines | Ganho |
|---------|------------------|----------------|-------|
| 🚀 Deploy em 10 servidores | 5 minutos | **30 segundos** | **10x mais rápido** |
| 🏥 Health check de 20 serviços | 1 minuto | **5 segundos** | **12x mais rápido** |
| 📊 Processar 1000 registros | 10 segundos | **1 segundo** | **10x mais rápido** |

---

## 📦 Importação

```lua
local goroutine = require("goroutine")
```

## Funções Disponíveis

### 1. `goroutine.spawn(function)`

Executa uma função em uma nova goroutine.

**Parâmetros:**
- `function`: Função Lua a ser executada em paralelo

**Retorno:** Nenhum

**Exemplo:**
```lua
goroutine.spawn(function()
    log.info("Executando em paralelo!")
end)
```

---

### 2. `goroutine.spawn_many(count, function)`

Executa múltiplas instâncias de uma função em goroutines separadas.

**Parâmetros:**
- `count` (number): Número de goroutines a criar
- `function`: Função que recebe o ID da goroutine como parâmetro

**Retorno:** Nenhum

**Exemplo:**
```lua
goroutine.spawn_many(5, function(id)
    log.info("Goroutine #" .. tostring(id))
end)
```

---

### 3. `goroutine.wait_group()`

Cria um WaitGroup para sincronização de goroutines.

**Retorno:** Objeto WaitGroup com os métodos:
- `add(delta)`: Incrementa o contador
- `done()`: Decrementa o contador
- `wait()`: Aguarda até o contador chegar a zero

**Exemplo:**
```lua
local wg = goroutine.wait_group()

wg:add(3)

for i = 1, 3 do
    goroutine.spawn(function()
        -- Fazer trabalho
        log.info("Worker " .. i)
        wg:done()
    end)
end

wg:wait()  -- Aguarda todas as goroutines
```

---

### 4. `goroutine.pool_create(name, options)`

Cria um worker pool para gerenciar execução paralela de tarefas.

**Parâmetros:**
- `name` (string): Nome único do pool
- `options` (table): Configurações do pool
  - `workers` (number): Número de workers (padrão: 4)

**Retorno:** `true` em sucesso

**Exemplo:**
```lua
goroutine.pool_create("mypool", { workers = 10 })
```

---

### 5. `goroutine.pool_submit(name, function, ...)`

Submete uma tarefa para execução em um worker pool.

**Parâmetros:**
- `name` (string): Nome do pool
- `function`: Função a ser executada
- `...`: Argumentos opcionais para a função

**Retorno:** 
- `task_id` (string): ID da tarefa submetida
- `error` (string): Mensagem de erro se falhar

**Exemplo:**
```lua
local task_id = goroutine.pool_submit("mypool", function()
    return "Resultado"
end)

if task_id then
    log.info("Tarefa submetida: " .. task_id)
end
```

---

### 6. `goroutine.pool_wait(name)`

Aguarda até que todas as tarefas do pool sejam concluídas.

**Parâmetros:**
- `name` (string): Nome do pool

**Retorno:** `true` em sucesso

**Exemplo:**
```lua
goroutine.pool_wait("mypool")
```

---

### 7. `goroutine.pool_close(name)`

Fecha um worker pool e libera recursos.

**Parâmetros:**
- `name` (string): Nome do pool

**Retorno:** `true` em sucesso

**Exemplo:**
```lua
goroutine.pool_close("mypool")
```

---

### 8. `goroutine.pool_stats(name)`

Retorna estatísticas de um worker pool.

**Parâmetros:**
- `name` (string): Nome do pool

**Retorno:** Table com estatísticas:
- `name` (string): Nome do pool
- `workers` (number): Número de workers
- `active` (number): Tarefas em execução
- `completed` (number): Tarefas concluídas
- `failed` (number): Tarefas que falharam
- `queued` (number): Tarefas na fila

**Exemplo:**
```lua
local stats = goroutine.pool_stats("mypool")
log.info("Concluídas: " .. stats.completed)
log.info("Ativas: " .. stats.active)
```

---

### 9. `goroutine.async(function)`

Executa uma função de forma assíncrona e retorna um handle.

**Parâmetros:**
- `function`: Função a ser executada

**Retorno:** Handle para await

**Exemplo:**
```lua
local handle = goroutine.async(function()
    -- Operação demorada
    return "resultado"
end)
```

---

### 10. `goroutine.await(handle)`

Aguarda a conclusão de uma operação async.

**Parâmetros:**
- `handle`: Handle retornado por `async()`

**Retorno:**
- `success` (boolean): Se a operação foi bem-sucedida
- `...`: Valores retornados pela função async

**Exemplo:**
```lua
local handle = goroutine.async(function()
    return "valor1", "valor2"
end)

local success, val1, val2 = goroutine.await(handle)
if success then
    log.info("Resultados: " .. val1 .. ", " .. val2)
end
```

---

### 11. `goroutine.await_all(handles)`

Aguarda a conclusão de múltiplas operações async.

**Parâmetros:**
- `handles` (table): Array de handles

**Retorno:** Table com resultados:
```lua
{
    { success = true, values = {...} },
    { success = false, error = "..." },
    ...
}
```

**Exemplo:**
```lua
local handles = {}
for i = 1, 5 do
    handles[i] = goroutine.async(function()
        return "Resultado " .. i
    end)
end

local results = goroutine.await_all(handles)
for i, result in ipairs(results) do
    if result.success then
        log.info("Task " .. i .. ": " .. result.values[1])
    end
end
```

---

### 12. `goroutine.sleep(milliseconds)`

Pausa a execução por um período especificado.

**Parâmetros:**
- `milliseconds` (number): Tempo em milissegundos

**Retorno:** Nenhum

**Exemplo:**
```lua
goroutine.sleep(1000)  -- Dorme por 1 segundo
```

---

### 13. `goroutine.timeout(milliseconds, function)`

Executa uma função com um timeout.

**Parâmetros:**
- `milliseconds` (number): Tempo máximo em milissegundos
- `function`: Função a ser executada

**Retorno:**
- `success` (boolean): `false` se timeout
- `...`: Valores retornados ou mensagem de erro

**Exemplo:**
```lua
local success, result = goroutine.timeout(5000, function()
    -- Operação que pode demorar
    return "resultado"
end)

if success then
    log.info("Concluído: " .. result)
else
    log.error("Timeout: " .. result)
end
```

---

## Exemplos Práticos

### Exemplo 1: Worker Pool para Processamento Paralelo

```lua
local process_files_task = task("process_files")
    :description("Processa arquivos em paralelo")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        -- Criar pool com 5 workers
        goroutine.pool_create("fileprocessor", { workers = 5 })
        
        local files = {"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"}
        
        -- Submeter tarefas
        for _, file in ipairs(files) do
            goroutine.pool_submit("fileprocessor", function()
                log.info("Processando: " .. file)
                goroutine.sleep(1000)  -- Simula processamento
                return "Processado: " .. file
            end)
        end
        
        -- Aguardar conclusão
        goroutine.pool_wait("fileprocessor")
        
        -- Ver estatísticas
        local stats = goroutine.pool_stats("fileprocessor")
        log.info("Total processado: " .. stats.completed)
        
        -- Limpar
        goroutine.pool_close("fileprocessor")
        
        return true
    end)
    :delegate_to("mariguica")
    :build()
```

### Exemplo 2: Operações Assíncronas com Async/Await

```lua
local fetch_data_task = task("fetch_data")
    :description("Busca dados de múltiplas fontes em paralelo")
    :command(function(this, params)
        local goroutine = require("goroutine")
        local http = require("http")
        
        -- Iniciar buscas assíncronas
        local h1 = goroutine.async(function()
            return http.get("https://api1.example.com/data")
        end)
        
        local h2 = goroutine.async(function()
            return http.get("https://api2.example.com/data")
        end)
        
        local h3 = goroutine.async(function()
            return http.get("https://api3.example.com/data")
        end)
        
        -- Aguardar todos os resultados
        local results = goroutine.await_all({h1, h2, h3})
        
        -- Processar resultados
        local all_success = true
        for i, result in ipairs(results) do
            if not result.success then
                log.error("API " .. i .. " falhou: " .. result.error)
                all_success = false
            end
        end
        
        return all_success
    end)
    :delegate_to("mariguica")
    :timeout("30s")
    :build()
```

### Exemplo 3: Sincronização com WaitGroup

```lua
local parallel_tasks = task("parallel_tasks")
    :description("Executa múltiplas tarefas com sincronização")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        local wg = goroutine.wait_group()
        local results = {}
        
        -- Adicionar 3 tarefas
        wg:add(3)
        
        -- Task 1: Download
        goroutine.spawn(function()
            log.info("Baixando arquivo...")
            goroutine.sleep(2000)
            results.download = "OK"
            wg:done()
        end)
        
        -- Task 2: Processar
        goroutine.spawn(function()
            log.info("Processando dados...")
            goroutine.sleep(1500)
            results.process = "OK"
            wg:done()
        end)
        
        -- Task 3: Upload
        goroutine.spawn(function()
            log.info("Fazendo upload...")
            goroutine.sleep(1000)
            results.upload = "OK"
            wg:done()
        end)
        
        -- Aguardar todas
        log.info("Aguardando conclusão...")
        wg:wait()
        
        log.info("Todas as tarefas concluídas!")
        log.info("Download: " .. results.download)
        log.info("Process: " .. results.process)
        log.info("Upload: " .. results.upload)
        
        return true
    end)
    :delegate_to("mariguica")
    :build()
```

### Exemplo 4: Timeout para Operações Críticas

```lua
local critical_operation = task("critical_operation")
    :description("Operação com timeout de segurança")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        local success, result = goroutine.timeout(5000, function()
            -- Operação que pode travar
            log.info("Executando operação crítica...")
            goroutine.sleep(3000)  -- Simulação
            return "Operação concluída"
        end)
        
        if success then
            log.info("✅ " .. result)
            return true
        else
            log.error("❌ Timeout: " .. result)
            return false
        end
    end)
    :delegate_to("mariguica")
    :build()
```

## Melhores Práticas

### 1. Sempre Fechar Pools

```lua
-- ✅ BOM
goroutine.pool_create("mypool", { workers = 5 })
-- ... usar pool
goroutine.pool_wait("mypool")
goroutine.pool_close("mypool")

-- ❌ RUIM - vazamento de recursos
goroutine.pool_create("mypool", { workers = 5 })
-- ... esqueceu de fechar
```

### 2. Usar WaitGroups para Sincronização

```lua
-- ✅ BOM
local wg = goroutine.wait_group()
wg:add(3)

for i = 1, 3 do
    goroutine.spawn(function()
        -- trabalho
        wg:done()
    end)
end

wg:wait()

-- ❌ RUIM - não garante ordem
for i = 1, 3 do
    goroutine.spawn(function()
        -- trabalho sem sincronização
    end)
end
```

### 3. Tratar Erros em Operações Async

```lua
-- ✅ BOM
local success, result = goroutine.await(handle)
if success then
    log.info("OK: " .. result)
else
    log.error("Erro: " .. result)
    -- Tratamento de erro
end

-- ❌ RUIM - assume sucesso
local _, result = goroutine.await(handle)
log.info(result)  -- pode ser erro!
```

### 4. Dimensionar Pools Adequadamente

```lua
-- ✅ BOM - baseado em cores disponíveis
local cpus = 4  -- ou detectar dinamicamente
goroutine.pool_create("cpu-bound", { workers = cpus })

-- ✅ BOM - I/O bound pode ter mais workers
goroutine.pool_create("io-bound", { workers = cpus * 2 })

-- ❌ RUIM - muito poucos workers
goroutine.pool_create("mypool", { workers = 1 })

-- ❌ RUIM - workers demais
goroutine.pool_create("mypool", { workers = 1000 })
```

### 5. Usar Timeouts para Operações Externas

```lua
-- ✅ BOM
local success, data = goroutine.timeout(10000, function()
    return fetch_external_api()
end)

-- ❌ RUIM - pode travar indefinidamente
fetch_external_api()
```

## Performance e Limitações

### Capacidades
- ✅ Execução verdadeiramente paralela usando goroutines do Go
- ✅ Overhead muito baixo para criar goroutines
- ✅ Suporta milhares de goroutines simultâneas
- ✅ Worker pools com gerenciamento eficiente de recursos
- ✅ Sincronização segura com WaitGroups

### Limitações
- ⚠️ Cada goroutine spawned cria um novo estado Lua (overhead de memória)
- ⚠️ Variáveis não são compartilhadas entre goroutines (use valores de retorno)
- ⚠️ Worker pools têm buffer limitado de tarefas (padrão: 2x workers)
- ⚠️ Async handles não podem ser reutilizados após await

## Troubleshooting

### Pool Queue Cheio
```lua
local task_id, err = goroutine.pool_submit("mypool", fn)
if not task_id then
    log.warn("Pool cheio: " .. err)
    -- Aguardar ou aumentar workers
end
```

### Detectar Goroutines Travadas
```lua
-- Usar timeout para detectar travamentos
local success, result = goroutine.timeout(5000, function()
    -- operação suspeita
end)

if not success then
    log.error("Possível deadlock detectado!")
end
```

### Monitorar Pool
```lua
-- Verificar periodicamente
local stats = goroutine.pool_stats("mypool")
if stats.failed > 0 then
    log.warn("Tarefas falharam: " .. stats.failed)
end

if stats.active == 0 and stats.queued == 0 then
    log.info("Pool está ocioso")
end
```

## Compatibilidade

- ✅ Funciona com `:delegate_to()` para execução remota
- ✅ Compatível com todos os outros módulos
- ✅ Suporta nested goroutines
- ✅ Thread-safe em todas as operações
- ✅ Funciona em Linux, macOS e Windows

---

## 🎯 Exemplos Completos e Prontos para Usar

### 🚀 Exemplo Real: Deploy Paralelo em Múltiplos Servidores

Este exemplo mostra como deployar uma aplicação em 6 servidores simultaneamente, reduzindo o tempo de 5 minutos para 30 segundos!

```lua
-- examples/parallel_deployment.sloth
local deploy_to_servers = task("deploy_multi_server")
    :description("Deploy application to multiple servers in parallel")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        local servers = {
            {name = "web-01", host = "192.168.1.10"},
            {name = "web-02", host = "192.168.1.11"},
            {name = "web-03", host = "192.168.1.12"},
            {name = "api-01", host = "192.168.1.20"},
            {name = "api-02", host = "192.168.1.21"},
            {name = "db-01", host = "192.168.1.30"},
        }
        
        log.info("🚀 Starting parallel deployment to " .. #servers .. " servers...")
        
        -- Create async handles for parallel deployment
        local handles = {}
        for _, server in ipairs(servers) do
            local handle = goroutine.async(function()
                log.info("📦 Deploying to " .. server.name .. " (" .. server.host .. ")")
                
                -- Simulate deployment steps
                local steps = {
                    "Uploading application files...",
                    "Installing dependencies...",
                    "Restarting services...",
                    "Running health checks..."
                }
                
                for _, step in ipairs(steps) do
                    log.info("  → " .. server.name .. ": " .. step)
                    goroutine.sleep(500)  -- Sleep 500ms to simulate work
                end
                
                return server.name, server.host, "success", os.date("%Y-%m-%d %H:%M:%S")
            end)
            
            table.insert(handles, handle)
        end
        
        log.info("⏳ Waiting for all deployments to complete...")
        
        -- Wait for all async operations to complete
        local results = goroutine.await_all(handles)
        
        -- Process results
        local success_count = 0
        local failed_count = 0
        
        log.info("\n📊 Deployment Results:")
        log.info("═══════════════════════════════════════")
        
        for i, result in ipairs(results) do
            if result.success then
                success_count = success_count + 1
                local server_name = result.values[1]
                local deployed_at = result.values[4]
                log.info("✅ " .. server_name .. " → Deployed successfully at " .. deployed_at)
            else
                failed_count = failed_count + 1
                log.error("❌ " .. (result.error or "Unknown deployment failure"))
            end
        end
        
        log.info("═══════════════════════════════════════")
        log.info("📈 Summary: " .. success_count .. " successful, " .. failed_count .. " failed")
        
        return success_count == #servers, "Deployment completed", {
            total = #servers,
            success = success_count,
            failed = failed_count
        }
    end)
    :timeout("2m")
    :build()

workflow.define("parallel_deployment")
    :description("Deploy to multiple servers in parallel")
    :version("1.0.0")
    :tasks({ deploy_to_servers })
    :config({ timeout = "5m" })
```

**Como executar:**
```bash
sloth-runner run -f examples/parallel_deployment.sloth
```

### 🏥 Exemplo Real: Health Check Paralelo

Verifique a saúde de múltiplos serviços simultaneamente:

```lua
-- examples/parallel_health_check.sloth
local parallel_health_check = task("check_services_health")
    :description("Check health of multiple services in parallel")
    :command(function(this, params)
        local goroutine = require("goroutine")
        local http = require("http")
        
        local services = {
            {name = "API Gateway", url = "http://localhost:8080/health"},
            {name = "Auth Service", url = "http://localhost:8081/health"},
            {name = "Database Service", url = "http://localhost:8082/health"},
            {name = "Cache Service", url = "http://localhost:8083/health"},
            {name = "Queue Service", url = "http://localhost:8084/health"},
        }
        
        log.info("🏥 Starting parallel health checks for " .. #services .. " services...")
        
        local handles = {}
        for _, service in ipairs(services) do
            local handle = goroutine.async(function()
                local start_time = os.clock()
                local success, response = pcall(function()
                    return http.get(service.url, {
                        timeout = 5,
                        headers = { ["User-Agent"] = "Sloth-Runner-HealthCheck/1.0" }
                    })
                end)
                
                local elapsed = (os.clock() - start_time) * 1000
                
                if success and response and response.status_code == 200 then
                    return service.name, "healthy", elapsed, response.body or ""
                else
                    local error_msg = response and response.error or "Connection failed"
                    return service.name, "unhealthy", elapsed, error_msg
                end
            end)
            
            table.insert(handles, handle)
        end
        
        log.info("⏳ Waiting for all health checks to complete...")
        
        local results = goroutine.await_all(handles)
        
        local healthy_count = 0
        local unhealthy_count = 0
        
        log.info("\n🏥 Health Check Results:")
        log.info("═══════════════════════════════════════════════")
        
        for _, result in ipairs(results) do
            if result.success then
                local name = result.values[1]
                local status = result.values[2]
                local time_ms = string.format("%.2f", result.values[3])
                
                if status == "healthy" then
                    healthy_count = healthy_count + 1
                    log.info("✅ " .. name .. ": " .. status .. " (" .. time_ms .. "ms)")
                else
                    unhealthy_count = unhealthy_count + 1
                    local error = result.values[4]
                    log.error("❌ " .. name .. ": " .. status .. " - " .. error)
                end
            else
                unhealthy_count = unhealthy_count + 1
                log.error("❌ Error: " .. (result.error or "Unknown error"))
            end
        end
        
        log.info("═══════════════════════════════════════════════")
        log.info("📊 Summary: " .. healthy_count .. " healthy, " .. unhealthy_count .. " unhealthy")
        
        return unhealthy_count == 0, "Health check completed", {
            total = #services,
            healthy = healthy_count,
            unhealthy = unhealthy_count
        }
    end)
    :timeout("30s")
    :build()

workflow.define("health_check_workflow")
    :description("Parallel health check for multiple services")
    :version("1.0.0")
    :tasks({ parallel_health_check })
```

### 🏭 Exemplo Real: Worker Pool para Processar Grande Volume

Processe milhares de itens com controle de concorrência:

```lua
-- examples/worker_pool_example.sloth
local process_with_pool = task("worker_pool_processing")
    :description("Process tasks using a worker pool")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        log.info("🏭 Creating worker pool with 5 workers...")
        goroutine.pool_create("data_processing", { workers = 5 })
        
        local tasks = {}
        for i = 1, 50 do
            tasks[i] = {
                id = i,
                data = "Task #" .. i,
                priority = math.random(1, 3)
            }
        end
        
        log.info("📋 Submitting " .. #tasks .. " tasks to worker pool...")
        
        for _, task_data in ipairs(tasks) do
            goroutine.pool_submit("data_processing", function()
                log.info("⚙️ Processing " .. task_data.data)
                goroutine.sleep(100 * task_data.priority)
                return {
                    id = task_data.id,
                    status = "completed",
                    processed_at = os.date("%H:%M:%S")
                }
            end)
        end
        
        log.info("⏳ Waiting for all tasks to complete...")
        goroutine.pool_wait("data_processing")
        
        local stats = goroutine.pool_stats("data_processing")
        
        log.info("\n📊 Worker Pool Statistics:")
        log.info("═══════════════════════════════════════")
        log.info("👷 Workers: " .. stats.workers)
        log.info("✅ Completed: " .. stats.completed)
        log.info("❌ Failed: " .. stats.failed)
        log.info("═══════════════════════════════════════")
        
        goroutine.pool_close("data_processing")
        
        return true, "All tasks processed successfully", {
            total_tasks = #tasks,
            completed = stats.completed,
            failed = stats.failed
        }
    end)
    :timeout("5m")
    :build()

workflow.define("worker_pool_workflow")
    :description("Process multiple tasks with a worker pool")
    :version("1.0.0")
    :tasks({ process_with_pool })
```

---

## 📚 Mais Recursos

- 📖 [Documentação Completa](../../README.md#-parallel-execution-with-goroutines-)
- 🧪 [Mais Exemplos](../../examples/)
- 🎯 [Benchmarks de Performance](../performance.md)
- 💬 [Discussões e Suporte](https://github.com/chalkan3-sloth/sloth-runner/discussions)
