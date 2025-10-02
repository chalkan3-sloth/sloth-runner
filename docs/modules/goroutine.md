# Módulo Goroutine

## Visão Geral

O módulo `goroutine` fornece capacidades de execução concorrente para scripts Sloth Lua, permitindo que você execute tarefas em paralelo usando goroutines do Go. Este módulo oferece diversas primitivas de concorrência incluindo spawn de goroutines, worker pools, async/await, e sincronização com WaitGroups.

## Importação

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
