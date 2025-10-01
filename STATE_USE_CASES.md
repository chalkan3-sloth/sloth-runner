# 📊 **Para que serve o módulo State no Sloth Runner?**

O módulo `state` é fundamental para **persistir dados entre execuções** de workflows e **coordenar tarefas**. Aqui estão os principais usos:

## 🎯 **1. Controle de Versionamento e Deploy**

### Exemplo: Rastrear última versão deployada
```lua
local deploy_task = task("deploy_app")
    :command(function()
        local state = require("state")
        
        -- Verificar última versão deployada
        local last_version = state.get("last_deployed_version") or "v1.0.0"
        local current_version = "v1.2.3"
        
        log.info("Última versão: " .. last_version)
        log.info("Nova versão: " .. current_version)
        
        if last_version == current_version then
            log.info("⏭️ Versão já está deployada, pulando...")
            return true, "Version already deployed"
        end
        
        -- Fazer deploy...
        log.info("🚀 Deployando nova versão...")
        
        -- Salvar nova versão após sucesso
        state.set("last_deployed_version", current_version)
        state.set("last_deploy_time", os.date())
        
        return true, "Deploy successful"
    end)
```

## 🔄 **2. Coordenação entre Tarefas Distribuídas**

### Exemplo: Master-Worker coordination
```lua
-- Worker registra que terminou
local worker_task = task("worker_process")
    :command(function()
        local state = require("state")
        
        -- Processar dados...
        log.info("Worker processando dados...")
        
        -- Incrementar contador de workers concluídos
        local completed = state.increment("workers_completed") or 1
        log.info("Workers concluídos: " .. completed)
        
        -- Marcar este worker como concluído
        state.set("worker_" .. os.getenv("WORKER_ID") .. "_status", "completed")
        
        return true, "Worker completed"
    end)

-- Master verifica se todos terminaram
local master_task = task("master_coordinator")
    :command(function()
        local state = require("state")
        
        local total_workers = 5
        local completed = state.get("workers_completed") or 0
        
        if completed >= total_workers then
            log.info("🎉 Todos os workers terminaram!")
            
            -- Processar resultados finais
            state.set("final_processing_status", "in_progress")
            
            return true, "All workers completed"
        else
            log.info("⏳ Aguardando workers: " .. completed .. "/" .. total_workers)
            return false, "Workers still running"
        end
    end)
```

## 🛡️ **3. Circuit Breaker e Rate Limiting**

### Exemplo: Evitar sobrecarga de APIs
```lua
local api_task = task("call_external_api")
    :command(function()
        local state = require("state")
        
        -- Verificar se API está com problemas
        local failures = state.get("api_failures") or 0
        local last_failure = state.get("last_api_failure")
        
        -- Circuit breaker: se muitas falhas recentes, não tenta
        if failures >= 5 and last_failure then
            local time_diff = os.time() - tonumber(last_failure)
            if time_diff < 300 then -- 5 minutos
                log.warn("🔴 Circuit breaker ativo - API com problemas")
                return false, "Circuit breaker active"
            end
        end
        
        -- Rate limiting: máximo 10 calls por minuto
        local calls_this_minute = state.get("api_calls_" .. os.date("%Y%m%d%H%M")) or 0
        if calls_this_minute >= 10 then
            log.warn("⏱️ Rate limit atingido")
            return false, "Rate limit exceeded"
        end
        
        -- Fazer chamada da API...
        local success = call_api()
        
        if success then
            -- Resetar contador de falhas
            state.set("api_failures", 0)
            -- Incrementar calls deste minuto
            state.increment("api_calls_" .. os.date("%Y%m%d%H%M"))
        else
            -- Incrementar falhas
            state.increment("api_failures")
            state.set("last_api_failure", os.time())
        end
        
        return success, success and "API call successful" or "API call failed"
    end)
```

## 📈 **4. Métricas e Monitoramento**

### Exemplo: Coletar estatísticas
```lua
local metrics_task = task("collect_metrics")
    :command(function()
        local state = require("state")
        
        -- Incrementar contadores
        state.increment("total_deployments")
        state.increment("deployments_today_" .. os.date("%Y%m%d"))
        
        -- Rastrear tempo médio
        local start_time = state.get("deploy_start_time")
        if start_time then
            local duration = os.time() - tonumber(start_time)
            
            -- Calcular média móvel
            local avg_duration = state.get("avg_deploy_duration") or duration
            local new_avg = (avg_duration + duration) / 2
            state.set("avg_deploy_duration", new_avg)
            
            log.info("⏱️ Deploy duration: " .. duration .. "s")
            log.info("📊 Average duration: " .. new_avg .. "s")
        end
        
        return true, "Metrics collected"
    end)
```

## 🔐 **5. Controle de Acesso e Locks Distribuídos**

### Exemplo: Prevenir execuções simultâneas
```lua
local critical_task = task("critical_operation")
    :command(function()
        local state = require("state")
        
        -- Verificar se já há uma operação crítica rodando
        local lock_owner = state.get("critical_operation_lock")
        local my_id = "instance_" .. os.getenv("HOSTNAME")
        
        if lock_owner and lock_owner ~= my_id then
            log.warn("🔒 Operação crítica já está rodando em: " .. lock_owner)
            return false, "Critical operation already running"
        end
        
        -- Adquirir lock
        state.set("critical_operation_lock", my_id)
        state.set("critical_operation_started", os.time())
        
        -- Operação crítica...
        log.info("🔧 Executando operação crítica...")
        
        -- Liberar lock
        state.delete("critical_operation_lock")
        
        return true, "Critical operation completed"
    end)
```

## 🗄️ **6. Cache de Resultados Caros**

### Exemplo: Cache de downloads ou processamento
```lua
local expensive_task = task("process_large_dataset")
    :command(function()
        local state = require("state")
        
        local data_hash = "dataset_v2_hash_abc123"
        local cache_key = "processed_" .. data_hash
        
        -- Verificar se já processamos este dataset
        local cached_result = state.get(cache_key)
        if cached_result then
            log.info("📋 Usando resultado em cache")
            return true, "Used cached result", { result = cached_result }
        end
        
        -- Processar dataset (operação cara)
        log.info("⚙️ Processando dataset...")
        local result = process_dataset()
        
        -- Salvar no cache por 24 horas
        state.set_with_ttl(cache_key, result, 24 * 3600)
        
        return true, "Dataset processed", { result = result }
    end)
```

## 🎛️ **7. Configuração Dinâmica**

### Exemplo: Feature flags e configuração em runtime
```lua
local feature_task = task("deploy_with_features")
    :command(function()
        local state = require("state")
        
        -- Verificar feature flags
        local new_algorithm_enabled = state.get("feature_new_algorithm") == "true"
        local beta_features = state.get("beta_features_enabled") == "true"
        
        log.info("🚀 New algorithm: " .. tostring(new_algorithm_enabled))
        log.info("🧪 Beta features: " .. tostring(beta_features))
        
        if new_algorithm_enabled then
            -- Deploy com novo algoritmo
            deploy_with_new_algorithm()
        else
            -- Deploy tradicional
            deploy_traditional()
        end
        
        return true, "Deploy completed with current features"
    end)
```

## 🎯 **Resumo dos Casos de Uso:**

| **Caso de Uso** | **Exemplo** |
|----------------|-------------|
| 🔄 **Controle de Estado** | Última versão deployada, status de serviços |
| 🤝 **Coordenação** | Workers distribuídos, master-slave |
| 🛡️ **Proteção** | Circuit breakers, rate limiting |
| 📊 **Métricas** | Contadores, durações, estatísticas |
| 🔒 **Locks** | Operações críticas, prevenir concorrência |
| 💾 **Cache** | Resultados caros, processamento |
| ⚙️ **Configuração** | Feature flags, parâmetros dinâmicos |

**O módulo state é essencial para workflows profissionais que precisam de coordenação, persistência e controle de estado entre execuções!** 🚀