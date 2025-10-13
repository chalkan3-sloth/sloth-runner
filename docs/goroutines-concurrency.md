# Goroutines & Concurrency - Sloth Runner

## 📖 Visão Geral

O módulo `goroutine` do Sloth Runner fornece um sistema completo de concorrência inspirado em Go, permitindo execução paralela, sincronização avançada e padrões de concorrência de alto nível. Este sistema traz a simplicidade e poder do modelo de concorrência de Go para o ambiente Lua.

## 🎯 Índice

- [Primitivas Básicas](#primitivas-básicas)
  - [Spawn & Spawn Many](#spawn--spawn-many)
  - [Wait Group](#wait-group)
  - [Sleep](#sleep)
- [Channels](#channels)
  - [Criação de Channels](#criação-de-channels)
  - [Operações com Channels](#operações-com-channels)
  - [Channel Range](#channel-range)
- [Sincronização](#sincronização)
  - [Mutex](#mutex)
  - [RWMutex](#rwmutex)
  - [Semaphore](#semaphore)
  - [Condition Variables](#condition-variables)
  - [Sync.Once](#synconce)
- [Operações Atômicas](#operações-atômicas)
- [Context](#context)
- [Worker Pools](#worker-pools)
- [Async/Await](#asyncawait)
- [Select Statement](#select-statement)
- [Padrões Avançados](#padrões-avançados)
  - [Pipeline](#pipeline)
  - [Fan-Out](#fan-out)
  - [Fan-In](#fan-in)
- [Exemplos Práticos](#exemplos-práticos)

---

## Primitivas Básicas

### Spawn & Spawn Many

Execute funções em goroutines separadas (threads leves).

#### `goroutine.spawn(function)`

Executa uma função em uma goroutine separada.

```lua
goroutine.spawn(function()
    log.info("Running in separate goroutine")
    goroutine.sleep(1000)
    log.info("Done!")
end)
```

#### `goroutine.spawn_many(count, function)`

Executa múltiplas goroutines com identificadores únicos.

```lua
goroutine.spawn_many(5, function(id)
    log.info("Worker " .. id .. " started")
    goroutine.sleep(math.random(100, 500))
    log.info("Worker " .. id .. " finished")
end)
```

### Wait Group

Sincronize múltiplas goroutines esperando que todas completem.

```lua
local wg = goroutine.wait_group()

wg:add(3)

for i = 1, 3 do
    goroutine.spawn(function()
        log.info("Task " .. i .. " running")
        goroutine.sleep(100)
        wg:done()
    end)
end

wg:wait() -- Espera todas as 3 tasks completarem
log.info("All tasks completed!")
```

**Métodos:**
- `wg:add(delta)` - Adiciona delta ao contador
- `wg:done()` - Decrementa o contador em 1
- `wg:wait()` - Bloqueia até o contador chegar a 0

### Sleep

Pausa a execução por um período especificado.

```lua
goroutine.sleep(1000) -- Dorme por 1000ms (1 segundo)
```

---

## Channels

Channels são o mecanismo principal de comunicação entre goroutines, permitindo troca segura de dados.

### Criação de Channels

```lua
-- Channel unbuffered (blocking)
local ch = goroutine.channel()

-- Channel buffered com capacidade 10
local ch = goroutine.channel(10)

-- Channel send-only
local sendCh = goroutine.channel(5, "send")

-- Channel receive-only
local recvCh = goroutine.channel(5, "receive")
```

### Operações com Channels

#### Send e Receive (Bloqueantes)

```lua
local ch = goroutine.channel(5)

-- Producer
goroutine.spawn(function()
    for i = 1, 10 do
        ch:send(i)
        log.info("Sent: " .. i)
    end
    ch:close()
end)

-- Consumer
goroutine.spawn(function()
    while true do
        local value, ok = ch:receive()
        if not ok then
            log.info("Channel closed")
            break
        end
        log.info("Received: " .. value)
    end
end)
```

#### Try Send e Try Receive (Não-Bloqueantes)

```lua
local ch = goroutine.channel(1)

-- Try send (retorna false se channel estiver cheio)
local ok = ch:try_send("value")
if ok then
    log.info("Sent successfully")
else
    log.info("Channel full, send failed")
end

-- Try receive (retorna nil, false se channel estiver vazio)
local value, ok = ch:try_receive()
if ok then
    log.info("Received: " .. value)
else
    log.info("Channel empty")
end
```

#### Métodos Utilitários

```lua
local ch = goroutine.channel(10)

-- Verificar capacidade e tamanho
local capacity = ch:cap()      -- 10
local length = ch:len()        -- Número atual de elementos

-- Verificar se está fechado
local closed = ch:is_closed()  -- true/false

-- Fechar channel
ch:close()
```

### Channel Range

Itera sobre um channel até que ele seja fechado. Similar ao `for range` do Go.

```lua
local ch = goroutine.channel(10)

-- Producer
goroutine.spawn(function()
    for i = 1, 100 do
        ch:send(i)
    end
    ch:close()
end)

-- Consumer usando range
ch:range(function(value)
    log.info("Processing: " .. value)
    goroutine.sleep(10)
end)

log.info("All values processed!")
```

**Características:**
- Para automaticamente quando o channel é fechado
- Tratamento de erros integrado
- Sintaxe limpa e idiomática

---

## Sincronização

### Mutex

Exclusão mútua para proteger recursos compartilhados.

```lua
local mu = goroutine.mutex()
local counter = 0

for i = 1, 10 do
    goroutine.spawn(function()
        for j = 1, 100 do
            mu:lock()
            counter = counter + 1
            mu:unlock()
        end
    end)
end

goroutine.sleep(1000)
log.info("Counter: " .. counter) -- 1000 (thread-safe)
```

**Métodos:**
- `mu:lock()` - Adquire o lock (bloqueante)
- `mu:unlock()` - Libera o lock
- `mu:try_lock()` - Tenta adquirir o lock (não-bloqueante)

### RWMutex

Read-Write Mutex permite múltiplos leitores ou um único escritor.

```lua
local rwmu = goroutine.rwmutex()
local data = { value = 0 }

-- Múltiplos leitores
for i = 1, 5 do
    goroutine.spawn(function()
        rwmu:rlock()
        log.info("Reader " .. i .. ": " .. data.value)
        goroutine.sleep(100)
        rwmu:runlock()
    end)
end

-- Um escritor
goroutine.spawn(function()
    rwmu:lock()
    data.value = 42
    log.info("Writer updated value")
    rwmu:unlock()
end)
```

**Métodos:**
- `rwmu:rlock()` - Read lock (múltiplos permitidos)
- `rwmu:runlock()` - Read unlock
- `rwmu:lock()` - Write lock (exclusivo)
- `rwmu:unlock()` - Write unlock
- `rwmu:try_rlock()` - Try read lock (não-bloqueante)
- `rwmu:try_lock()` - Try write lock (não-bloqueante)

### Semaphore

Controla o acesso a um número limitado de recursos.

```lua
local sem = goroutine.semaphore(3) -- Máximo 3 recursos simultâneos

for i = 1, 10 do
    goroutine.spawn(function()
        sem:acquire() -- Adquire um token
        log.info("Worker " .. i .. " acquired resource")

        -- Usa o recurso
        goroutine.sleep(500)

        log.info("Worker " .. i .. " releasing resource")
        sem:release() -- Libera o token
    end)
end
```

**Métodos:**
- `sem:acquire()` - Adquire um token (bloqueante)
- `sem:release()` - Libera um token
- `sem:try_acquire()` - Tenta adquirir (não-bloqueante)
- `sem:available()` - Número de tokens disponíveis
- `sem:capacity()` - Capacidade total

### Condition Variables

Sincronização complexa com condições.

```lua
local cond = goroutine.cond()
local mu = cond:get_mutex()
local ready = false

-- Waiter
goroutine.spawn(function()
    mu:lock()
    while not ready do
        log.info("Waiting for signal...")
        cond:wait() -- Libera lock, espera sinal, re-adquire lock
    end
    log.info("Received signal!")
    mu:unlock()
end)

-- Signaler
goroutine.sleep(1000)
mu:lock()
ready = true
cond:signal() -- Acorda um waiter
-- ou cond:broadcast() -- Acorda todos os waiters
mu:unlock()
```

**Métodos:**
- `cond:wait()` - Espera por sinal
- `cond:signal()` - Acorda um waiter
- `cond:broadcast()` - Acorda todos os waiters
- `cond:get_mutex()` - Obtém o mutex associado

### Sync.Once

Garante que uma função execute apenas uma vez.

```lua
local once = goroutine.once()
local config = nil

local initConfig = function()
    log.info("Initializing config (runs only once)")
    config = { host = "localhost", port = 8080 }
end

-- Múltiplas goroutines tentam inicializar
for i = 1, 5 do
    goroutine.spawn(function()
        once:call(initConfig) -- Apenas uma executa
        log.info("Config ready: " .. config.host)
    end)
end
```

---

## Operações Atômicas

Operações lock-free para contadores e valores compartilhados.

```lua
local counter = goroutine.atomic_int(0)

-- Spawn 10 goroutines incrementando
for i = 1, 10 do
    goroutine.spawn(function()
        for j = 1, 100 do
            counter:add(1) -- Atomic increment
        end
    end)
end

goroutine.sleep(1000)
log.info("Final count: " .. counter:load()) -- 1000
```

**Métodos:**
- `counter:add(delta)` - Adiciona atomicamente, retorna novo valor
- `counter:load()` - Lê o valor atomicamente
- `counter:store(value)` - Armazena valor atomicamente
- `counter:swap(new)` - Troca valor, retorna o antigo
- `counter:compare_and_swap(old, new)` - CAS operation

**Exemplo CAS Loop (Lock-Free):**

```lua
local value = goroutine.atomic_int(0)

local function safe_increment()
    while true do
        local old = value:load()
        local new = old + 1
        if value:compare_and_swap(old, new) then
            return new
        end
        -- CAS falhou, retry
    end
end
```

---

## Context

Gerenciamento de cancelamento, timeouts e deadlines.

### Context Básico

```lua
local ctx = goroutine.context()

-- Verificar cancelamento
if ctx:is_cancelled() then
    log.info("Context cancelled")
end

-- Cancelar manualmente
ctx:cancel()

-- Obter erro
local err = ctx:err() -- "context canceled" ou nil
```

### Context com Timeout

```lua
local ctx = goroutine.context()
local timeoutCtx = ctx:with_timeout(5000) -- 5 segundos

goroutine.spawn(function()
    for i = 1, 100 do
        if timeoutCtx:is_cancelled() then
            log.info("Timeout reached at iteration " .. i)
            return
        end

        -- Trabalho
        goroutine.sleep(100)
    end
end)
```

### Context com Deadline

```lua
local ctx = goroutine.context()
local deadline_ms = os.time() * 1000 + 10000 -- 10 segundos no futuro
local deadlineCtx = ctx:with_deadline(deadline_ms)

-- Verificar deadline
local dl, has_deadline = deadlineCtx:deadline()
if has_deadline then
    log.info("Deadline em: " .. dl .. "ms")
end
```

### Context Hierárquico (Cascading)

```lua
local parentCtx = goroutine.context()

-- Criar child contexts
local child1, cancel1 = parentCtx:with_cancel()
local child2, cancel2 = parentCtx:with_cancel()

-- Cancelar parent cancela todos os children
parentCtx:cancel()
```

**Métodos:**
- `ctx:with_cancel()` - Cria child context com cancelamento
- `ctx:with_timeout(ms)` - Cria context com timeout
- `ctx:with_deadline(deadline_ms)` - Cria context com deadline
- `ctx:is_cancelled()` - Verifica se foi cancelado
- `ctx:err()` - Retorna erro se cancelado
- `ctx:cancel()` - Cancela o context
- `ctx:deadline()` - Retorna deadline (se existir)

---

## Worker Pools

Pools de workers para processamento paralelo gerenciado.

```lua
-- Criar pool com 10 workers
goroutine.pool_create("mypool", { workers = 10 })

-- Submeter tasks
for i = 1, 100 do
    goroutine.pool_submit("mypool", function()
        log.info("Processing task " .. i)
        goroutine.sleep(100)
    end)
end

-- Obter estatísticas
local stats = goroutine.pool_stats("mypool")
log.info("Active: " .. stats.active)
log.info("Completed: " .. stats.completed)
log.info("Failed: " .. stats.failed)

-- Fechar pool
goroutine.pool_close("mypool")
```

**Funções:**
- `pool_create(name, options)` - Cria um pool
- `pool_submit(name, fn, ...)` - Submete task ao pool
- `pool_wait(name)` - Espera todas as tasks completarem
- `pool_stats(name)` - Retorna estatísticas do pool
- `pool_close(name)` - Fecha o pool

---

## Async/Await

Programação assíncrona estilo Promise.

```lua
-- Executar função assíncrona
local handle1 = goroutine.async(function()
    goroutine.sleep(1000)
    return "Result 1"
end)

local handle2 = goroutine.async(function()
    goroutine.sleep(500)
    return "Result 2"
end)

-- Await individual
local success, result = goroutine.await(handle1)
if success then
    log.info("Got: " .. result)
end

-- Await all
local results = goroutine.await_all({handle1, handle2})
for i, result in ipairs(results) do
    if result.success then
        log.info("Result " .. i .. ": " .. result.values[1])
    else
        log.error("Error " .. i .. ": " .. result.error)
    end
end
```

### Timeout com Async

```lua
local success, result = goroutine.timeout(2000, function()
    -- Operação que pode levar muito tempo
    goroutine.sleep(5000)
    return "Done"
end)

if success then
    log.info("Completed: " .. result)
else
    log.info("Timeout: " .. result) -- "timeout exceeded"
end
```

---

## Select Statement

Multiplexação de operações em channels.

### Select Básico

```lua
local ch1 = goroutine.channel(1)
local ch2 = goroutine.channel(1)

goroutine.spawn(function()
    goroutine.sleep(100)
    ch1:send("from ch1")
end)

goroutine.spawn(function()
    goroutine.sleep(200)
    ch2:send("from ch2")
end)

-- Select espera o primeiro canal disponível
goroutine.select({
    {
        channel = ch1,
        receive = true,
        handler = function(value)
            log.info("Ch1: " .. value)
        end
    },
    {
        channel = ch2,
        receive = true,
        handler = function(value)
            log.info("Ch2: " .. value)
        end
    },
    {
        default = true,
        handler = function()
            log.info("No channel ready")
        end
    }
})
```

### Select com Timeout

```lua
local ch1 = goroutine.channel(1)
local ch2 = goroutine.channel(1)

local timedout, result = goroutine.select_timeout(1000, {
    {
        channel = ch1,
        receive = true,
        handler = function(value)
            log.info("Received from ch1: " .. value)
        end
    },
    {
        channel = ch2,
        receive = true,
        handler = function(value)
            log.info("Received from ch2: " .. value)
        end
    }
})

if timedout then
    log.info("Select timed out after 1000ms")
else
    log.info("Case " .. result .. " executed")
end
```

### Select com Send

```lua
local ch = goroutine.channel(1)

goroutine.select({
    {
        channel = ch,
        send = "my value",
        handler = function()
            log.info("Successfully sent")
        end
    },
    {
        default = true,
        handler = function()
            log.info("Channel full, can't send")
        end
    }
})
```

---

## Padrões Avançados

### Pipeline

Crie pipelines de processamento em múltiplas etapas com paralelização configurável.

```lua
local input = goroutine.channel(100)

-- Pipeline: Stage1 → Stage2 → Stage3
local output = goroutine.pipeline(input, {
    -- Stage 1: Parse (2 workers)
    {
        workers = 2,
        fn = function(raw)
            return parse(raw)
        end
    },
    -- Stage 2: Transform (5 workers)
    {
        workers = 5,
        fn = function(data)
            return transform(data)
        end
    },
    -- Stage 3: Enrich (3 workers)
    {
        workers = 3,
        fn = function(data)
            return enrich(data)
        end
    }
})

-- Producer
goroutine.spawn(function()
    for i = 1, 1000 do
        input:send(raw_data[i])
    end
    input:close()
end)

-- Consumer
output:range(function(processed)
    save_to_db(processed)
end)
```

**Caso de Uso:** ETL (Extract, Transform, Load), processamento de logs, data streaming

### Fan-Out

Distribui trabalho de um canal para múltiplos workers paralelos.

```lua
local jobs = goroutine.channel(100)
local outputs = goroutine.fan_out(jobs, 10) -- 10 workers paralelos

-- Spawn workers para processar
for i, outCh in ipairs(outputs) do
    goroutine.spawn(function()
        local worker_id = i
        outCh:range(function(job)
            log.info("Worker " .. worker_id .. " processing: " .. job.id)
            process_job(job)
        end)
    end)
end

-- Feed jobs
for i = 1, 1000 do
    jobs:send({ id = i, data = "..." })
end
jobs:close()
```

**Caso de Uso:** Web scraping paralelo, processamento de imagens, API calls paralelas

### Fan-In

Merge múltiplos canais em um único output.

```lua
-- Múltiplas fontes de dados
local source1 = goroutine.channel(10)
local source2 = goroutine.channel(10)
local source3 = goroutine.channel(10)

-- Merge em um canal
local merged = goroutine.fan_in({source1, source2, source3})

-- Producers
goroutine.spawn(function()
    for i = 1, 100 do
        source1:send("S1-" .. i)
        goroutine.sleep(10)
    end
    source1:close()
end)

goroutine.spawn(function()
    for i = 1, 100 do
        source2:send("S2-" .. i)
        goroutine.sleep(15)
    end
    source2:close()
end)

goroutine.spawn(function()
    for i = 1, 100 do
        source3:send("S3-" .. i)
        goroutine.sleep(20)
    end
    source3:close()
end)

-- Consumer único processando de todas as fontes
merged:range(function(value)
    log.info("Processing: " .. value)
end)
```

**Caso de Uso:** Agregação de logs de múltiplos servidores, merge de resultados distribuídos

---

## Exemplos Práticos

### 1. Web Scraper Paralelo

```lua
workflow.define("parallel_scraper", {
    description = "Scrape múltiplas URLs em paralelo",
    tasks = {
        scrape = {
            description = "Parallel web scraping",
            command = function()
                local urls = {
                    "https://example.com/page1",
                    "https://example.com/page2",
                    -- ... mais URLs
                }

                local urlCh = goroutine.channel(#urls)
                local sem = goroutine.semaphore(10) -- Máx 10 concurrent
                local wg = goroutine.wait_group()

                -- Feed URLs
                for _, url in ipairs(urls) do
                    urlCh:send(url)
                end
                urlCh:close()

                -- Workers
                urlCh:range(function(url)
                    wg:add(1)
                    goroutine.spawn(function()
                        sem:acquire()

                        log.info("Scraping: " .. url)
                        local content = http.get(url)
                        save_content(url, content)

                        sem:release()
                        wg:done()
                    end)
                end)

                wg:wait()
                return true, "Scraping completed"
            end
        }
    }
})
```

### 2. Pipeline de Processamento de Dados

```lua
workflow.define("data_pipeline", {
    description = "ETL pipeline com goroutines",
    tasks = {
        process = {
            command = function()
                local rawData = goroutine.channel(1000)

                -- Pipeline: Extract → Validate → Transform → Load
                local loaded = goroutine.pipeline(rawData, {
                    {
                        workers = 3,
                        fn = function(raw)
                            return extract(raw)
                        end
                    },
                    {
                        workers = 5,
                        fn = function(data)
                            if validate(data) then
                                return data
                            end
                            return nil
                        end
                    },
                    {
                        workers = 4,
                        fn = function(data)
                            return transform(data)
                        end
                    },
                    {
                        workers = 2,
                        fn = function(data)
                            load_to_db(data)
                            return data
                        end
                    }
                })

                -- Feed data
                goroutine.spawn(function()
                    for line in io.lines("data.csv") do
                        rawData:send(line)
                    end
                    rawData:close()
                end)

                -- Wait completion
                local count = 0
                loaded:range(function(data)
                    count = count + 1
                end)

                log.info("Processed " .. count .. " records")
                return true
            end
        }
    }
})
```

### 3. Rate-Limited API Calls

```lua
workflow.define("api_calls", {
    description = "Rate-limited parallel API calls",
    tasks = {
        call_api = {
            command = function()
                local requests = goroutine.channel(100)
                local sem = goroutine.semaphore(5) -- 5 concurrent
                local results = {}
                local mu = goroutine.mutex()
                local wg = goroutine.wait_group()

                -- Feed requests
                for i = 1, 100 do
                    requests:send({ id = i, endpoint = "/api/data/" .. i })
                end
                requests:close()

                -- Workers
                requests:range(function(req)
                    wg:add(1)
                    goroutine.spawn(function()
                        sem:acquire()

                        local response = http.get(req.endpoint)

                        mu:lock()
                        table.insert(results, response)
                        mu:unlock()

                        goroutine.sleep(200) -- Rate limiting
                        sem:release()
                        wg:done()
                    end)
                end)

                wg:wait()
                log.info("Completed " .. #results .. " API calls")
                return true
            end
        }
    }
})
```

### 4. Real-Time Log Aggregation

```lua
workflow.define("log_aggregator", {
    description = "Aggregate logs from multiple sources",
    tasks = {
        aggregate = {
            command = function()
                local servers = {"server1", "server2", "server3"}
                local logChannels = {}

                -- Create channel per server
                for _, server in ipairs(servers) do
                    local ch = goroutine.channel(100)
                    table.insert(logChannels, ch)

                    -- Tail logs from each server
                    goroutine.spawn(function()
                        while true do
                            local log = fetch_log(server)
                            if log then
                                ch:send({
                                    server = server,
                                    message = log,
                                    timestamp = os.time()
                                })
                            end
                            goroutine.sleep(100)
                        end
                    end)
                end

                -- Fan-in all logs
                local allLogs = goroutine.fan_in(logChannels)

                -- Process merged logs
                allLogs:range(function(logEntry)
                    log.info("[" .. logEntry.server .. "] " .. logEntry.message)
                    store_log(logEntry)
                end)

                return true
            end
        }
    }
})
```

### 5. Producer-Consumer com Bounded Queue

```lua
local function bounded_queue_example()
    local queue = goroutine.channel(10) -- Buffer de 10
    local wg = goroutine.wait_group()

    -- Producer
    wg:add(1)
    goroutine.spawn(function()
        for i = 1, 100 do
            queue:send(i)
            log.info("Produced: " .. i)
        end
        queue:close()
        wg:done()
    end)

    -- Consumers (3 paralelos)
    for consumer_id = 1, 3 do
        wg:add(1)
        goroutine.spawn(function()
            local id = consumer_id
            queue:range(function(item)
                log.info("Consumer " .. id .. " processing: " .. item)
                goroutine.sleep(50) -- Simulate work
            end)
            wg:done()
        end)
    end

    wg:wait()
    log.info("Queue processing completed")
end
```

---

## 🎓 Melhores Práticas

### 1. Sempre Feche Channels

```lua
local ch = goroutine.channel(10)

-- ✅ BOM
goroutine.spawn(function()
    for i = 1, 10 do
        ch:send(i)
    end
    ch:close() -- Importante!
end)

-- ❌ RUIM - Channel nunca fecha, range fica travado
goroutine.spawn(function()
    for i = 1, 10 do
        ch:send(i)
    end
    -- Esqueceu de fechar!
end)
```

### 2. Use WaitGroups para Sincronização

```lua
local wg = goroutine.wait_group()

-- ✅ BOM
wg:add(5)
for i = 1, 5 do
    goroutine.spawn(function()
        -- trabalho
        wg:done()
    end)
end
wg:wait() -- Espera todas completarem

-- ❌ RUIM - Pode sair antes das goroutines terminarem
for i = 1, 5 do
    goroutine.spawn(function()
        -- trabalho
    end)
end
goroutine.sleep(1000) -- Não confiável!
```

### 3. Proteja Recursos Compartilhados

```lua
local data = { count = 0 }
local mu = goroutine.mutex()

-- ✅ BOM
mu:lock()
data.count = data.count + 1
mu:unlock()

-- ❌ RUIM - Race condition!
data.count = data.count + 1
```

### 4. Use Semaphores para Limitar Concorrência

```lua
local sem = goroutine.semaphore(10) -- Máx 10 simultâneos

-- ✅ BOM
for i = 1, 1000 do
    goroutine.spawn(function()
        sem:acquire()
        -- trabalho pesado
        sem:release()
    end)
end

-- ❌ RUIM - Pode criar 1000 goroutines e sobrecarregar sistema
for i = 1, 1000 do
    goroutine.spawn(function()
        -- trabalho pesado
    end)
end
```

### 5. Use Context para Cancelamento

```lua
local ctx = goroutine.context()
local timeoutCtx = ctx:with_timeout(5000)

-- ✅ BOM
goroutine.spawn(function()
    while not timeoutCtx:is_cancelled() do
        -- trabalho
        goroutine.sleep(100)
    end
    log.info("Gracefully stopped")
end)

-- ❌ RUIM - Sem forma de parar gracefully
goroutine.spawn(function()
    while true do
        -- trabalho sem controle de parada
        goroutine.sleep(100)
    end
end)
```

### 6. Prefira Atomic para Contadores Simples

```lua
-- ✅ BOM - Lock-free e mais rápido
local counter = goroutine.atomic_int(0)
counter:add(1)

-- ❌ RUIM - Overhead de mutex para contador simples
local mu = goroutine.mutex()
local counter = 0
mu:lock()
counter = counter + 1
mu:unlock()
```

---

## 📊 Tabela de Referência Rápida

| Primitiva | Use Quando | Complexidade | Performance |
|-----------|-----------|--------------|-------------|
| `spawn` | Executar função assíncrona | Baixa | Alta |
| `wait_group` | Esperar grupo de goroutines | Baixa | Alta |
| `channel` | Comunicar entre goroutines | Média | Alta |
| `mutex` | Proteger recurso compartilhado | Baixa | Média |
| `rwmutex` | Muitos leitores, poucos escritores | Média | Alta |
| `semaphore` | Limitar concorrência | Baixa | Alta |
| `atomic_int` | Contador lock-free | Baixa | Muito Alta |
| `once` | Inicialização única | Baixa | Alta |
| `cond` | Sincronização complexa | Alta | Média |
| `context` | Cancelamento/Timeout | Média | Alta |
| `select` | Multiplexar channels | Média | Alta |
| `pipeline` | Processamento em etapas | Alta | Alta |
| `fan_out` | Distribuir trabalho | Média | Alta |
| `fan_in` | Agregar resultados | Média | Alta |

---

## 🔗 Recursos Adicionais

- **Exemplos Completos**: `/examples/goroutine_*.sloth`
- **Código Fonte**: `/internal/modules/core/goroutine.go`
- **Padrões de Concorrência**: Este documento

---

## 📝 Conclusão

O sistema de goroutines do Sloth Runner fornece todas as ferramentas necessárias para criar aplicações altamente concorrentes e eficientes. Com primitivas inspiradas em Go e padrões de alto nível como Pipeline, Fan-Out e Fan-In, você pode construir sistemas complexos de processamento distribuído com código limpo e manutenível.

**Principais Vantagens:**
- ✅ Modelo de concorrência simples e poderoso
- ✅ Primitivas de sincronização completas
- ✅ Padrões de alto nível prontos para uso
- ✅ Performance próxima a Go nativo
- ✅ Integração perfeita com Lua

---

**Versão**: 1.0.0
**Última Atualização**: 2025
**Autor**: Sloth Runner Team
