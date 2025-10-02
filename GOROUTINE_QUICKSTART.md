# 🚀 Goroutine Module - Quick Start

## 🎯 O Que É?

O módulo **goroutine** permite que você execute código Lua em **paralelo** usando as goroutines do Go. Isso significa que você pode:

- ✅ Processar múltiplos arquivos simultaneamente
- ✅ Fazer várias requisições HTTP em paralelo
- ✅ Executar tarefas independentes ao mesmo tempo
- ✅ Acelerar significativamente seus workflows

## 🚀 Instalação

Já vem incluído! Apenas atualize seu sloth-runner:

```bash
cd /Users/chalkan3/.projects/task-runner
go build -o sloth-runner ./cmd/sloth-runner
```

## 💡 Exemplo em 30 Segundos

```lua
local goroutine = require("goroutine")

local my_task = task("parallel_demo")
    :command(function(this, params)
        -- Criar um pool de 5 workers
        goroutine.pool_create("demo", { workers = 5 })
        
        -- Submeter 10 tarefas que rodam em paralelo
        for i = 1, 10 do
            goroutine.pool_submit("demo", function()
                log.info("Executando tarefa " .. i .. " em paralelo!")
                goroutine.sleep(1000)  -- Simula trabalho
            end)
        end
        
        -- Aguardar todas terminarem
        goroutine.pool_wait("demo")
        goroutine.pool_close("demo")
        
        log.info("✅ Todas as 10 tarefas concluídas!")
        return true
    end)
    :build()

workflow.define("demo")
    :tasks({ my_task })
```

**Salve como `demo.sloth` e execute:**

```bash
./sloth-runner run demo.sloth
```

## 📚 Principais Funções

### 1. Worker Pools (Recomendado para Processamento em Lote)

```lua
-- Criar pool
goroutine.pool_create("mypool", { workers = 10 })

-- Submeter tarefas
for i = 1, 100 do
    goroutine.pool_submit("mypool", function()
        -- seu código aqui
    end)
end

-- Aguardar conclusão
goroutine.pool_wait("mypool")
goroutine.pool_close("mypool")
```

### 2. Async/Await (Recomendado para Operações I/O)

```lua
-- Iniciar operação async
local handle = goroutine.async(function()
    return fetch_data_from_api()
end)

-- Fazer outras coisas...

-- Aguardar resultado
local success, result = goroutine.await(handle)
if success then
    log.info("Resultado: " .. result)
end
```

### 3. Spawn Simples (Para Fire-and-Forget)

```lua
-- Executar em background
goroutine.spawn(function()
    log.info("Executando em paralelo!")
end)

-- Ou múltiplas
goroutine.spawn_many(5, function(id)
    log.info("Worker " .. id)
end)
```

### 4. Sincronização com WaitGroup

```lua
local wg = goroutine.wait_group()
wg:add(3)

goroutine.spawn(function() task1(); wg:done() end)
goroutine.spawn(function() task2(); wg:done() end)
goroutine.spawn(function() task3(); wg:done() end)

wg:wait()  -- Aguarda todas
```

## 🎯 Casos de Uso Reais

### 1. Processar Múltiplos Arquivos

```lua
local process_files = task("process_files")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        local files = {
            "/data/file1.csv",
            "/data/file2.csv",
            "/data/file3.csv",
            -- ... mais arquivos
        }
        
        goroutine.pool_create("processor", { workers = 5 })
        
        for _, file in ipairs(files) do
            goroutine.pool_submit("processor", function()
                log.info("Processando " .. file)
                -- processar arquivo
                log.info("✅ " .. file .. " concluído")
            end)
        end
        
        goroutine.pool_wait("processor")
        goroutine.pool_close("processor")
        
        return true
    end)
    :build()
```

### 2. Health Check de Múltiplos Serviços

```lua
local health_check = task("health_check")
    :command(function(this, params)
        local goroutine = require("goroutine")
        local http = require("http")
        
        local services = {
            "https://api1.com/health",
            "https://api2.com/health",
            "https://api3.com/health"
        }
        
        local handles = {}
        
        -- Checar todos em paralelo
        for i, url in ipairs(services) do
            handles[i] = goroutine.async(function()
                local resp = http.get(url)
                return resp.status == 200
            end)
        end
        
        -- Aguardar todos
        local results = goroutine.await_all(handles)
        
        -- Verificar resultados
        local all_healthy = true
        for i, result in ipairs(results) do
            if not result.success or not result.values[1] then
                log.error("❌ Serviço " .. i .. " não está saudável")
                all_healthy = false
            end
        end
        
        return all_healthy
    end)
    :build()
```

### 3. Pipeline CI/CD Paralelo

```lua
local ci_pipeline = task("ci_pipeline")
    :command(function(this, params)
        local goroutine = require("goroutine")
        
        log.info("🏗️ Stage 1: Build")
        local wg = goroutine.wait_group()
        wg:add(3)
        
        -- Build em paralelo
        goroutine.spawn(function()
            log.info("Building frontend...")
            -- build frontend
            wg:done()
        end)
        
        goroutine.spawn(function()
            log.info("Building backend...")
            -- build backend
            wg:done()
        end)
        
        goroutine.spawn(function()
            log.info("Building mobile...")
            -- build mobile
            wg:done()
        end)
        
        wg:wait()
        log.info("✅ Builds concluídos")
        
        log.info("🧪 Stage 2: Tests")
        goroutine.pool_create("tests", { workers = 4 })
        
        local test_suites = {"unit", "integration", "e2e", "security"}
        for _, suite in ipairs(test_suites) do
            goroutine.pool_submit("tests", function()
                log.info("Running " .. suite .. " tests...")
                -- run tests
            end)
        end
        
        goroutine.pool_wait("tests")
        goroutine.pool_close("tests")
        log.info("✅ Tests concluídos")
        
        return true
    end)
    :build()
```

## ⚡ Performance

### Sem Goroutines (Sequencial)
```
Task 1: 2s
Task 2: 2s
Task 3: 2s
Total: 6s ❌
```

### Com Goroutines (Paralelo)
```
Task 1: 2s ┐
Task 2: 2s ├─ Simultâneo
Task 3: 2s ┘
Total: 2s ✅ 3x mais rápido!
```

## 🎓 Próximos Passos

- 📖 Leia a [documentação completa](docs/modules/goroutine.md)
- 🎬 Veja [exemplos avançados](../../sandbox/test_goroutine.sloth)
- 🧪 Teste com seus próprios workflows

## 🐛 Troubleshooting

### Problema: "pool not found"
**Solução:** Certifique-se de criar o pool antes de usar:
```lua
goroutine.pool_create("mypool", { workers = 5 })
```

### Problema: Tarefas não terminam
**Solução:** Use timeout para detectar travamentos:
```lua
local success = goroutine.timeout(5000, function()
    -- operação
end)
```

### Problema: Resultados incorretos
**Solução:** Use WaitGroup para garantir sincronização:
```lua
local wg = goroutine.wait_group()
wg:add(n)
-- ... spawn tasks com wg:done()
wg:wait()
```

## 💬 Suporte

- 📝 Issues: https://github.com/chalkan3-sloth/sloth-runner/issues
- 📖 Docs: `docs/modules/goroutine.md`
- 🎯 Exemplos: `sandbox/test_goroutine*.sloth`

---

**Comece a acelerar seus workflows agora! 🚀**
