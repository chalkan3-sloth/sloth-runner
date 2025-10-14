# Sugestões de Features para Goroutines e Channels

Este documento apresenta sugestões de novas features que podem ser implementadas no projeto relacionadas a goroutines e channels.

## 🔄 1. Advanced Channel Patterns

### 1.1. Buffered Channel with Overflow Strategy

Canal bufferizado com estratégias diferentes para lidar com overflow (drop oldest, drop newest, block, error).

```lua
-- Criar canal com estratégia de overflow
local ch = goroutine.channel_with_overflow(10, "drop_oldest")
-- Opções: "drop_oldest", "drop_newest", "block", "error"

-- Uso normal
ch:send(value)
```

**Implementação Sugerida**: Adicionar `luaChannelWithOverflow` em `goroutine.go`

### 1.2. Priority Channel

Canal que processa mensagens baseado em prioridade.

```lua
-- Criar canal de prioridade
local pq = goroutine.priority_channel()

-- Enviar com prioridade (maior = mais prioritário)
pq:send(value, 10)  -- alta prioridade
pq:send(value, 1)   -- baixa prioridade

-- Receber sempre pega o de maior prioridade primeiro
local value = pq:receive()
```

**Implementação Sugerida**: Usar `container/heap` internamente

### 1.3. Debounced/Throttled Channel

Canal que limita frequência de processamento de mensagens.

```lua
-- Debounce: só processa após período de silêncio
local debounced = goroutine.debounce_channel(input, 500)  -- 500ms

-- Throttle: limita taxa de processamento
local throttled = goroutine.throttle_channel(input, 100)  -- max 1 msg/100ms
```

**Implementação Sugerida**: Usar timers e goroutines para controlar fluxo

### 1.4. Broadcast Channel

Canal que envia mensagens para múltiplos subscribers.

```lua
-- Criar broadcaster
local bc = goroutine.broadcaster()

-- Subscribers
local sub1 = bc:subscribe()
local sub2 = bc:subscribe()
local sub3 = bc:subscribe()

-- Broadcast envia para todos
bc:broadcast("message")

-- Cada subscriber recebe independentemente
local msg1 = sub1:receive()
local msg2 = sub2:receive()
local msg3 = sub3:receive()
```

**Implementação Sugerida**: Manter lista de subscribers e enviar para cada um

### 1.5. Merge Multiple Channels

Combinar múltiplos canais em um único output.

```lua
-- Merge N canais em um
local ch1 = goroutine.channel(10)
local ch2 = goroutine.channel(10)
local ch3 = goroutine.channel(10)

local merged = goroutine.merge({ch1, ch2, ch3})

-- Recebe de qualquer canal que tenha dados
merged:range(function(value)
  log.info("Received: " .. value)
end)
```

**Implementação Sugerida**: Similar ao fan-in existente, mas generalizado

---

## 🔄 2. Pipeline Enhancements

### 2.1. Pipeline com Error Handling

Pipeline que propaga erros através das stages.

```lua
local output, errors = goroutine.pipeline_with_errors(input, {
  {
    workers = 2,
    fn = function(x)
      if x < 0 then
        return nil, "negative value not allowed"
      end
      return x * 2, nil
    end
  },
  {
    workers = 2,
    fn = function(x)
      return x + 10, nil
    end
  }
})

-- Processar valores e erros separadamente
goroutine.spawn(function()
  errors:range(function(err)
    log.error("Pipeline error: " .. err)
  end)
end)

output:range(function(value)
  log.info("Result: " .. value)
end)
```

**Implementação Sugerida**: Cada stage retorna tupla `(value, error)` e propaga erros para canal separado

### 2.2. Pipeline com Retry Logic

Pipeline com retry automático para stages que falharem.

```lua
local output = goroutine.pipeline_with_retry(input, {
  {
    workers = 2,
    max_retries = 3,
    retry_delay = 100,  -- ms
    fn = function(x)
      -- Pode falhar temporariamente
      return risky_operation(x)
    end
  }
})
```

**Implementação Sugerida**: Wrapper em cada worker que tenta novamente em caso de erro

### 2.3. Pipeline com Rate Limiting

Pipeline que limita taxa de processamento.

```lua
local output = goroutine.pipeline_with_rate_limit(input, {
  {
    workers = 2,
    rate_limit = 100,  -- max 100 ops/segundo
    fn = function(x)
      return expensive_api_call(x)
    end
  }
})
```

**Implementação Sugerida**: Usar token bucket ou leaky bucket algorithm

### 2.4. Pipeline com Metrics/Observability

Pipeline que coleta métricas de performance automaticamente.

```lua
local output, metrics = goroutine.pipeline_with_metrics(input, {
  {
    name = "multiply",
    workers = 2,
    fn = function(x) return x * 2 end
  },
  {
    name = "add",
    workers = 2,
    fn = function(x) return x + 10 end
  }
})

-- Métricas disponíveis
goroutine.spawn(function()
  goroutine.sleep(5000)
  local stats = metrics:get()
  log.info("Stage 'multiply': " .. stats.multiply.processed .. " items, avg time: " .. stats.multiply.avg_time .. "ms")
  log.info("Stage 'add': " .. stats.add.processed .. " items, avg time: " .. stats.add.avg_time .. "ms")
end)
```

**Implementação Sugerida**: Instrumentar cada stage com contadores e timers

### 2.5. Dynamic Pipeline

Pipeline que pode adicionar/remover stages dinamicamente.

```lua
local pipeline = goroutine.dynamic_pipeline(input)

-- Adicionar stages
pipeline:add_stage({
  name = "stage1",
  workers = 2,
  fn = function(x) return x * 2 end
})

pipeline:add_stage({
  name = "stage2",
  workers = 2,
  fn = function(x) return x + 10 end
})

-- Remover stage
pipeline:remove_stage("stage1")

-- Obter output
local output = pipeline:output()
```

**Implementação Sugerida**: Manter estrutura dinâmica de stages e reconectar canais

---

## 🔄 3. Advanced Concurrency Patterns

### 3.1. Worker Pool com Priority

Pool de workers que processa tarefas baseado em prioridade.

```lua
local pool = goroutine.priority_worker_pool({
  workers = 10,
  fn = function(task)
    log.info("Processing task: " .. task.id)
    return process(task)
  end
})

-- Submeter tarefas com prioridade
pool:submit({id = 1, data = "..."}, 10)  -- alta prioridade
pool:submit({id = 2, data = "..."}, 1)   -- baixa prioridade

-- Aguardar conclusão
pool:wait()
pool:close()
```

**Implementação Sugerida**: Usar priority queue interna

### 3.2. Rate Limiter

Limitar taxa de execução de operações.

```lua
local limiter = goroutine.rate_limiter({
  rate = 100,      -- 100 operações
  per = 1000,      -- por segundo (1000ms)
  burst = 10       -- burst máximo de 10
})

for i = 1, 1000 do
  limiter:wait()  -- bloqueia se limite atingido
  api_call()
end
```

**Implementação Sugerida**: Token bucket algorithm com goroutine de refill

### 3.3. Circuit Breaker

Padrão circuit breaker para proteger contra falhas.

```lua
local cb = goroutine.circuit_breaker({
  threshold = 5,        -- falhas consecutivas para abrir
  timeout = 30000,      -- tempo em open state (ms)
  half_open_max = 3     -- tentativas em half-open
})

local success, result = cb:call(function()
  return risky_operation()
end)

if not success then
  log.error("Circuit breaker open: " .. result)
end

-- Obter estado
log.info("State: " .. cb:state())  -- "closed", "open", "half-open"
```

**Implementação Sugerida**: State machine com contadores de falha

### 3.4. Bulkhead

Isolamento de recursos para prevenir falha em cascata.

```lua
local bulkhead = goroutine.bulkhead({
  partitions = {
    api_service = {max_concurrent = 10},
    db_service = {max_concurrent = 20},
    cache_service = {max_concurrent = 5}
  }
})

-- Executar com isolamento
bulkhead:execute("api_service", function()
  return api_call()
end)

-- Se uma partição falhar, outras não são afetadas
```

**Implementação Sugerida**: Semaphore por partição

### 3.5. Retry com Exponential Backoff

Retry inteligente com backoff exponencial.

```lua
local result, err = goroutine.retry({
  max_attempts = 5,
  initial_delay = 100,   -- ms
  max_delay = 10000,     -- ms
  multiplier = 2,        -- dobra a cada tentativa
  jitter = true          -- adiciona randomização
}, function()
  return risky_operation()
end)
```

**Implementação Sugerida**: Loop com sleep progressivo

---

## 🔄 4. Message Queue Patterns

### 4.1. Pub/Sub System

Sistema completo de publicação/subscrição.

```lua
local pubsub = goroutine.pubsub()

-- Subscribers
local sub1 = pubsub:subscribe("events")
local sub2 = pubsub:subscribe("events")
local sub3 = pubsub:subscribe("logs")

-- Publisher
pubsub:publish("events", {type = "user_login", user = "john"})
pubsub:publish("logs", {level = "info", msg = "system started"})

-- Cada subscriber recebe mensagens do seu tópico
sub1:range(function(msg)
  log.info("Sub1 received: " .. msg.type)
end)
```

**Implementação Sugerida**: Map de tópicos para lista de subscribers

### 4.2. Request/Reply Pattern

Padrão request-reply com timeout.

```lua
local rr = goroutine.request_reply()

-- Server
goroutine.spawn(function()
  rr:serve(function(request)
    -- Processar request
    return {result = process(request.data)}
  end)
end)

-- Client
local response, err = rr:request({data = "..."}, 5000)  -- 5s timeout
if err then
  log.error("Request failed: " .. err)
else
  log.info("Response: " .. response.result)
end
```

**Implementação Sugerida**: Correlation IDs e map de pending requests

### 4.3. Work Queue com Dead Letter Queue

Fila de trabalho com DLQ para mensagens que falharam.

```lua
local queue = goroutine.work_queue({
  workers = 10,
  max_retries = 3,
  dlq_enabled = true
})

-- Processar trabalho
queue:process(function(item)
  return process(item)
end)

-- Enfileirar itens
for i = 1, 100 do
  queue:enqueue({id = i, data = "..."})
end

-- Processar DLQ
local dlq = queue:dead_letter_queue()
dlq:range(function(failed_item)
  log.error("Failed permanently: " .. failed_item.id)
  -- Pode enviar alerta, salvar em DB, etc
end)
```

**Implementação Sugerida**: Worker pool + retry logic + canal separado para falhas

---

## 🔄 5. Synchronization Enhancements

### 5.1. Barrier

Sincronização tipo barrier para N goroutines.

```lua
local barrier = goroutine.barrier(3)  -- aguarda 3 goroutines

for i = 1, 3 do
  goroutine.spawn(function()
    log.info("Worker " .. i .. " doing work...")
    goroutine.sleep(math.random(1000, 3000))

    log.info("Worker " .. i .. " reached barrier")
    barrier:wait()  -- bloqueia até todos chegarem

    log.info("Worker " .. i .. " continuing after barrier")
  end)
end
```

**Implementação Sugerida**: WaitGroup + Mutex + Condition Variable

### 5.2. Latch

Countdown latch (como Java CountDownLatch).

```lua
local latch = goroutine.latch(5)  -- aguarda 5 eventos

-- Workers
for i = 1, 5 do
  goroutine.spawn(function()
    work()
    latch:count_down()
  end)
end

-- Aguardar todos completarem
latch:wait()
log.info("All workers completed")
```

**Implementação Sugerida**: Atomic counter + channel para sinalização

### 5.3. ReadersWriter Lock com Upgrade

RWMutex que permite upgrade de read para write lock.

```lua
local rwlock = goroutine.rwmutex_upgradeable()

-- Read lock
rwlock:rlock()
local value = read_data()

-- Tentar upgrade para write
if rwlock:try_upgrade() then
  write_data(value + 1)
  rwlock:unlock()
else
  -- Upgrade falhou, fazer unlock e tentar wlock
  rwlock:runlock()
  rwlock:lock()
  write_data(value + 1)
  rwlock:unlock()
end
```

**Implementação Sugerida**: Controle de lock state mais sofisticado

### 5.4. Weighted Semaphore

Semaphore que permite adquirir múltiplos recursos.

```lua
local sem = goroutine.weighted_semaphore(10)

-- Adquirir 3 recursos
sem:acquire(3)
work_with_resources(3)
sem:release(3)

-- Adquirir 1 recurso
sem:acquire(1)
work_with_resources(1)
sem:release(1)
```

**Implementação Sugerida**: Modificar semaphore existente para suportar pesos

---

## 🔄 6. Monitoring and Observability

### 6.1. Goroutine Profiler

Ferramenta para monitorar goroutines criadas pelo módulo.

```lua
local profiler = goroutine.profiler()
profiler:start()

-- ... executar código ...

goroutine.sleep(5000)

local stats = profiler:stats()
log.info("Active goroutines: " .. stats.active)
log.info("Total created: " .. stats.total_created)
log.info("Total completed: " .. stats.completed)
log.info("Average lifetime: " .. stats.avg_lifetime .. "ms")
```

**Implementação Sugerida**: Instrumentar `spawn()` para rastrear goroutines

### 6.2. Channel Inspector

Ferramenta para inspecionar estado de canais.

```lua
local ch = goroutine.channel(100)
local inspector = goroutine.channel_inspector(ch)

-- Métricas
log.info("Length: " .. inspector:len())
log.info("Capacity: " .. inspector:cap())
log.info("Total sent: " .. inspector:total_sent())
log.info("Total received: " .. inspector:total_received())
log.info("Is closed: " .. tostring(inspector:is_closed()))
log.info("Blocked senders: " .. inspector:blocked_senders())
log.info("Blocked receivers: " .. inspector:blocked_receivers())
```

**Implementação Sugerida**: Wrapper em volta de channel que coleta métricas

### 6.3. Deadlock Detector

Detector automático de deadlocks.

```lua
local detector = goroutine.deadlock_detector({
  check_interval = 1000,  -- verificar a cada 1s
  timeout = 5000          -- considerar deadlock após 5s sem progresso
})

detector:start()

-- Se deadlock detectado, callback é chamado
detector:on_deadlock(function(info)
  log.error("DEADLOCK DETECTED!")
  log.error("Stuck goroutines: " .. info.stuck_count)
  log.error("Blocked on channels: " .. info.blocked_channels)
  -- Pode forçar panic, enviar alerta, etc
end)
```

**Implementação Sugerida**: Monitorar goroutines bloqueadas e detectar falta de progresso

---

## 🔄 7. Optimization Patterns

### 7.1. Bounded Parallelism

Controlar nível de paralelismo dinamicamente.

```lua
local bounded = goroutine.bounded_parallel({
  max_parallel = 10,
  queue_size = 100
})

-- Submeter 1000 tarefas, mas apenas 10 executam simultaneamente
for i = 1, 1000 do
  bounded:submit(function()
    work(i)
  end)
end

bounded:wait_all()
```

**Implementação Sugerida**: Semaphore + work queue

### 7.2. Batch Processor

Processar itens em batches automaticamente.

```lua
local batch = goroutine.batch_processor({
  batch_size = 100,
  flush_interval = 1000,  -- ms
  fn = function(items)
    -- Processar batch de items
    bulk_insert_db(items)
  end
})

-- Adicionar items individualmente
for i = 1, 10000 do
  batch:add({id = i, data = "..."})
  -- Flush automático quando batch_size ou flush_interval atingido
end

batch:flush()  -- flush final
batch:close()
```

**Implementação Sugerida**: Buffer + timer para flush periódico

### 7.3. Adaptive Worker Pool

Pool que ajusta número de workers baseado na carga.

```lua
local pool = goroutine.adaptive_pool({
  min_workers = 2,
  max_workers = 20,
  scale_up_threshold = 0.8,    -- 80% de utilização
  scale_down_threshold = 0.2,  -- 20% de utilização
  check_interval = 5000        -- verificar a cada 5s
})

pool:process(function(task)
  return work(task)
end)

-- Pool ajusta workers automaticamente
```

**Implementação Sugerida**: Monitorar queue length e CPU e ajustar dinamicamente

### 7.4. Coalescing Cache

Cache que coalece requisições simultâneas para a mesma key.

```lua
local cache = goroutine.coalescing_cache({
  ttl = 60000,  -- 60s
  loader = function(key)
    -- Função que carrega valor (chamada apenas uma vez por key, mesmo com múltiplos requests)
    return expensive_load(key)
  end
})

-- Múltiplas goroutines requisitam a mesma key
for i = 1, 100 do
  goroutine.spawn(function()
    local value = cache:get("expensive_key")
    -- Apenas 1 chamada a expensive_load() é feita
    process(value)
  end)
end
```

**Implementação Sugerida**: Map de pending requests + sync.Once por key

---

## 🎯 Priorização Sugerida

### Alta Prioridade (Maior Impacto / Menor Esforço)
1. **Rate Limiter** - Muito útil para APIs
2. **Retry com Exponential Backoff** - Pattern comum
3. **Worker Pool com Priority** - Extensão natural do pool existente
4. **Batch Processor** - Otimização importante para bulk operations
5. **Pipeline com Error Handling** - Feature crítica para pipelines de produção

### Média Prioridade
1. **Circuit Breaker** - Pattern importante para resiliência
2. **Broadcast Channel** - Pattern de comunicação útil
3. **Merge Multiple Channels** - Complementa fan-in existente
4. **Request/Reply Pattern** - Muito usado em microservices
5. **Channel Inspector** - Ótimo para debugging

### Baixa Prioridade (Mais Complexo / Casos Específicos)
1. **Deadlock Detector** - Complexo de implementar corretamente
2. **Dynamic Pipeline** - Caso de uso menos comum
3. **Adaptive Worker Pool** - Requer tuning cuidadoso
4. **Bulkhead** - Pattern avançado
5. **Coalescing Cache** - Otimização para casos específicos

---

## 📝 Considerações de Implementação

### Compatibilidade
- Manter compatibilidade com API existente
- Adicionar features como opt-in
- Documentar breaking changes claramente

### Performance
- Minimizar overhead de sincronização
- Usar estruturas de dados eficientes
- Profile antes e depois

### Testing
- Testes unitários para cada feature
- Testes de stress/concorrência
- Testes de edge cases (deadlocks, race conditions)

### Documentação
- Exemplos práticos para cada feature
- Patterns de uso recomendados
- Troubleshooting guides

---

## 🚀 Roadmap Sugerido

### Fase 1: Foundations (1-2 sprints)
- Rate Limiter
- Retry com Exponential Backoff
- Pipeline com Error Handling

### Fase 2: Advanced Patterns (2-3 sprints)
- Circuit Breaker
- Worker Pool com Priority
- Broadcast Channel
- Merge Multiple Channels

### Fase 3: Optimization (2 sprints)
- Batch Processor
- Bounded Parallelism
- Channel Inspector

### Fase 4: Advanced Features (3+ sprints)
- Request/Reply Pattern
- Pub/Sub System
- Adaptive Worker Pool
- Deadlock Detector

---

## 📚 Recursos Adicionais

### Inspiração
- Go standard library: `sync`, `context`, `sync/atomic`
- Go concurrency patterns: https://go.dev/blog/pipelines
- Hystrix (Circuit Breaker): https://github.com/Netflix/Hystrix
- Resilience4j: https://github.com/resilience4j/resilience4j

### Bibliotecas Go para Referência
- `golang.org/x/sync` - Extended sync primitives
- `github.com/sony/gobreaker` - Circuit breaker
- `github.com/cenkalti/backoff` - Backoff algorithms
- `golang.org/x/time/rate` - Rate limiting
