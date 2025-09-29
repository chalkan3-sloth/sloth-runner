-- 🛡️ Reliability Patterns - Padrões avançados de confiabilidade
-- Este exemplo demonstra circuit breakers, retries, timeouts e fallbacks

TaskDefinitions = {
    reliability_patterns = {
        description = "Demonstração de padrões avançados de confiabilidade e resiliência",
        
        tasks = {
            {
                name = "setup_reliability_demo",
                description = "Configura ambiente para demonstrar padrões de confiabilidade",
                command = function()
                    log.info("🛡️  Configurando demonstração de padrões de confiabilidade...")
                    
                    -- Configurar serviços simulados com diferentes níveis de confiabilidade
                    local services = {
                        stable_service = {
                            name = "Stable Service",
                            url = "https://httpbin.org/status/200",
                            reliability = 0.95, -- 95% uptime
                            avg_response_time = 100
                        },
                        unstable_service = {
                            name = "Unstable Service", 
                            url = "https://httpbin.org/status/500",
                            reliability = 0.30, -- 30% uptime (muito instável)
                            avg_response_time = 2000
                        },
                        slow_service = {
                            name = "Slow Service",
                            url = "https://httpbin.org/delay/10",
                            reliability = 0.80,
                            avg_response_time = 10000
                        },
                        intermittent_service = {
                            name = "Intermittent Service",
                            urls = {
                                "https://httpbin.org/status/200",  -- 50% success
                                "https://httpbin.org/status/503"   -- 50% failure
                            },
                            reliability = 0.50
                        }
                    }
                    
                    state.set("test_services", services)
                    
                    -- Inicializar métricas de confiabilidade
                    local metrics = {
                        total_requests = 0,
                        successful_requests = 0,
                        failed_requests = 0,
                        circuit_breaker_trips = 0,
                        fallback_executions = 0,
                        retry_attempts = 0
                    }
                    
                    state.set("reliability_metrics", metrics)
                    
                    log.info("✅ Ambiente configurado:")
                    for service_id, service in pairs(services) do
                        log.info("  🔧 " .. service.name .. " (confiabilidade: " .. (service.reliability * 100) .. "%)")
                    end
                    
                    return true, "Ambiente de confiabilidade configurado"
                end
            },
            
            {
                name = "demonstrate_retry_patterns",
                description = "Demonstra diferentes estratégias de retry",
                depends_on = "setup_reliability_demo",
                command = function()
                    log.info("🔄 Demonstrando padrões de retry...")
                    
                    local http = require("http")
                    local reliability = require("reliability")
                    local services = state.get("test_services")
                    local metrics = state.get("reliability_metrics")
                    
                    -- 1. Retry Linear (interval fixo)
                    log.info("📊 1. Retry Linear (intervalo fixo de 2s)...")
                    local linear_start = os.time()
                    local linear_result = reliability.retry(function()
                        metrics.total_requests = metrics.total_requests + 1
                        metrics.retry_attempts = metrics.retry_attempts + 1
                        
                        local result = http.get({
                            url = services.unstable_service.url,
                            timeout = 5
                        })
                        
                        if result.success and result.data.status_code == 200 then
                            metrics.successful_requests = metrics.successful_requests + 1
                            return result.data
                        else
                            metrics.failed_requests = metrics.failed_requests + 1
                            error("Service returned: " .. (result.data and result.data.status_code or "error"))
                        end
                    end, {
                        max_attempts = 3,
                        delay = 2000, -- 2 segundos
                        backoff = "linear"
                    })
                    local linear_end = os.time()
                    
                    if linear_result.success then
                        log.info("✅ Linear retry bem-sucedido após " .. linear_result.attempts .. " tentativas")
                    else
                        log.info("❌ Linear retry falhou: " .. linear_result.error)
                    end
                    log.info("  ⏱️  Tempo total: " .. (linear_end - linear_start) .. "s")
                    
                    -- 2. Retry Exponencial (backoff exponencial)
                    log.info("\n📊 2. Retry Exponencial (backoff com jitter)...")
                    local exp_start = os.time()
                    local exp_result = reliability.retry(function()
                        metrics.total_requests = metrics.total_requests + 1
                        metrics.retry_attempts = metrics.retry_attempts + 1
                        
                        local result = http.get({
                            url = services.unstable_service.url,
                            timeout = 5
                        })
                        
                        if result.success and result.data.status_code == 200 then
                            metrics.successful_requests = metrics.successful_requests + 1
                            return result.data
                        else
                            metrics.failed_requests = metrics.failed_requests + 1
                            error("Service returned: " .. (result.data and result.data.status_code or "error"))
                        end
                    end, {
                        max_attempts = 4,
                        initial_delay = 500, -- 500ms inicial
                        max_delay = 8000,    -- máximo 8s
                        backoff = "exponential",
                        jitter = true
                    })
                    local exp_end = os.time()
                    
                    if exp_result.success then
                        log.info("✅ Exponential retry bem-sucedido após " .. exp_result.attempts .. " tentativas")
                    else
                        log.info("❌ Exponential retry falhou: " .. exp_result.error)
                    end
                    log.info("  ⏱️  Tempo total: " .. (exp_end - exp_start) .. "s")
                    
                    -- 3. Retry com Condição Customizada
                    log.info("\n📊 3. Retry com Condição Customizada...")
                    local custom_result = reliability.retry(function()
                        metrics.total_requests = metrics.total_requests + 1
                        
                        local result = http.get({
                            url = services.slow_service.url,
                            timeout = 3 -- Timeout curto para forçar falha
                        })
                        
                        if result.success then
                            metrics.successful_requests = metrics.successful_requests + 1
                            return result.data
                        else
                            metrics.failed_requests = metrics.failed_requests + 1
                            error("Request timeout or error")
                        end
                    end, {
                        max_attempts = 2,
                        delay = 1000,
                        retry_on = function(error)
                            -- Só fazer retry em timeouts, não em outros erros
                            return string.find(string.lower(error), "timeout") ~= nil
                        end
                    })
                    
                    if custom_result.success then
                        log.info("✅ Custom retry bem-sucedido")
                    else
                        log.info("❌ Custom retry falhou (como esperado para serviço lento): " .. custom_result.error)
                    end
                    
                    state.set("reliability_metrics", metrics)
                    
                    return true, "Padrões de retry demonstrados"
                end
            },
            
            {
                name = "demonstrate_circuit_breaker",
                description = "Demonstra padrão Circuit Breaker",
                depends_on = "demonstrate_retry_patterns",
                command = function()
                    log.info("⚡ Demonstrando padrão Circuit Breaker...")
                    
                    local http = require("http")
                    local reliability = require("reliability")
                    local services = state.get("test_services")
                    local metrics = state.get("reliability_metrics")
                    
                    -- Configurar circuit breaker para serviço instável
                    local circuit_config = {
                        failure_threshold = 3,    -- Abrir após 3 falhas
                        recovery_timeout = 10,    -- Tentar recovery após 10s
                        success_threshold = 2     -- Fechar após 2 sucessos
                    }
                    
                    log.info("🔧 Configurando Circuit Breaker:")
                    log.info("  🚨 Limite de falhas: " .. circuit_config.failure_threshold)
                    log.info("  ⏰ Timeout de recovery: " .. circuit_config.recovery_timeout .. "s")
                    log.info("  ✅ Sucessos para fechar: " .. circuit_config.success_threshold)
                    
                    -- Criar circuit breaker
                    local circuit = reliability.circuit_breaker("unstable_service_cb", circuit_config)
                    
                    -- Função que será protegida pelo circuit breaker
                    local protected_call = function()
                        local result = http.get({
                            url = services.unstable_service.url,
                            timeout = 5
                        })
                        
                        if result.success and result.data.status_code == 200 then
                            return result.data
                        else
                            error("Service failed with status: " .. (result.data and result.data.status_code or "unknown"))
                        end
                    end
                    
                    -- Testar circuit breaker com múltiplas chamadas
                    log.info("\n🧪 Testando Circuit Breaker com " .. services.unstable_service.name .. ":")
                    
                    for i = 1, 10 do
                        log.info("📞 Chamada " .. i .. ":")
                        
                        local cb_result = circuit.call(protected_call)
                        metrics.total_requests = metrics.total_requests + 1
                        
                        if cb_result.success then
                            log.info("  ✅ Sucesso")
                            metrics.successful_requests = metrics.successful_requests + 1
                        elseif cb_result.circuit_open then
                            log.info("  ⚡ Circuit Breaker ABERTO - chamada não executada")
                            metrics.circuit_breaker_trips = metrics.circuit_breaker_trips + 1
                        else
                            log.info("  ❌ Falha: " .. cb_result.error)
                            metrics.failed_requests = metrics.failed_requests + 1
                        end
                        
                        -- Status do circuit breaker
                        local cb_stats = circuit.stats()
                        log.info("  📊 CB Status: " .. cb_stats.state .. " (falhas: " .. cb_stats.failures .. ")")
                        
                        -- Pequeno delay entre chamadas
                        exec.run("sleep 1")
                    end
                    
                    -- Mostrar estatísticas finais do circuit breaker
                    local final_stats = circuit.stats()
                    log.info("\n📈 Estatísticas Finais do Circuit Breaker:")
                    log.info("  Estado: " .. final_stats.state)
                    log.info("  Total de falhas: " .. final_stats.failures)
                    log.info("  Sucessos consecutivos: " .. final_stats.successes)
                    log.info("  Última falha: " .. (final_stats.last_failure_time or "N/A"))
                    
                    state.set("reliability_metrics", metrics)
                    
                    return true, "Circuit Breaker demonstrado"
                end
            },
            
            {
                name = "demonstrate_fallback_patterns",
                description = "Demonstra padrões de fallback e degradação graceful",
                depends_on = "demonstrate_circuit_breaker",
                command = function()
                    log.info("🔄 Demonstrando padrões de Fallback...")
                    
                    local http = require("http")
                    local reliability = require("reliability")
                    local services = state.get("test_services")
                    local metrics = state.get("reliability_metrics")
                    
                    -- 1. Fallback Simples - Cache local
                    log.info("📊 1. Fallback para Cache Local:")
                    
                    local cached_data = {
                        user_profile = {
                            id = 1,
                            name = "Cached User",
                            email = "cached@example.com",
                            source = "local_cache"
                        }
                    }
                    
                    local user_service_call = function()
                        local result = http.get({
                            url = services.unstable_service.url, -- Vai falhar
                            timeout = 3
                        })
                        
                        if result.success and result.data.status_code == 200 then
                            return {
                                id = 1,
                                name = "Live User", 
                                email = "live@example.com",
                                source = "live_service"
                            }
                        else
                            error("User service unavailable")
                        end
                    end
                    
                    local fallback_result = reliability.with_fallback(user_service_call, function(error)
                        log.info("  🔄 Serviço principal falhou, usando cache: " .. error)
                        metrics.fallback_executions = metrics.fallback_executions + 1
                        return cached_data.user_profile
                    end)
                    
                    if fallback_result.success then
                        local user = fallback_result.data
                        log.info("  ✅ Dados do usuário obtidos: " .. user.name .. " (fonte: " .. user.source .. ")")
                    end
                    
                    -- 2. Fallback com Múltiplas Estratégias
                    log.info("\n📊 2. Fallback com Múltiplas Estratégias:")
                    
                    local multi_fallback = function()
                        -- Estratégia 1: Serviço principal
                        local primary_result = pcall(function()
                            local result = http.get({
                                url = services.unstable_service.url,
                                timeout = 2
                            })
                            if not result.success or result.data.status_code ~= 200 then
                                error("Primary service failed")
                            end
                            return {source = "primary", data = "primary data"}
                        end)
                        
                        if primary_result then
                            return {source = "primary", data = "primary data"}
                        end
                        
                        log.info("  🔄 Primary falhou, tentando secondary...")
                        
                        -- Estratégia 2: Serviço secundário
                        local secondary_result = pcall(function()
                            local result = http.get({
                                url = services.stable_service.url, -- Este deve funcionar
                                timeout = 5
                            })
                            if not result.success or result.data.status_code ~= 200 then
                                error("Secondary service failed")
                            end
                            return {source = "secondary", data = "secondary data"}
                        end)
                        
                        if secondary_result then
                            log.info("  ✅ Secondary service funcionou")
                            return {source = "secondary", data = "secondary data"}
                        end
                        
                        log.info("  🔄 Secondary falhou, usando default...")
                        metrics.fallback_executions = metrics.fallback_executions + 1
                        
                        -- Estratégia 3: Dados padrão
                        return {source = "default", data = "default fallback data"}
                    end
                    
                    local multi_result = multi_fallback()
                    log.info("  📝 Resultado: " .. multi_result.data .. " (fonte: " .. multi_result.source .. ")")
                    
                    -- 3. Degradação Graceful
                    log.info("\n📊 3. Degradação Graceful de Funcionalidade:")
                    
                    local feature_availability = {
                        recommendations = false, -- Serviço de recomendações está down
                        user_profile = true,     -- Perfil básico funciona
                        basic_search = true      -- Busca básica funciona
                    }
                    
                    local get_user_experience = function()
                        local experience = {
                            basic_features = {},
                            advanced_features = {},
                            degraded_features = {}
                        }
                        
                        -- Features básicas sempre disponíveis
                        table.insert(experience.basic_features, "user_authentication")
                        table.insert(experience.basic_features, "basic_navigation")
                        
                        if feature_availability.user_profile then
                            table.insert(experience.basic_features, "user_profile")
                        else
                            table.insert(experience.degraded_features, "user_profile_limited")
                        end
                        
                        if feature_availability.basic_search then
                            table.insert(experience.basic_features, "search")
                        else
                            table.insert(experience.degraded_features, "offline_search")
                        end
                        
                        if feature_availability.recommendations then
                            table.insert(experience.advanced_features, "personalized_recommendations")
                        else
                            table.insert(experience.degraded_features, "default_recommendations")
                            log.info("  🔄 Recomendações personalizadas indisponíveis, usando padrões")
                        end
                        
                        return experience
                    end
                    
                    local user_exp = get_user_experience()
                    log.info("  ✅ Features básicas: " .. table.concat(user_exp.basic_features, ", "))
                    log.info("  ⚡ Features avançadas: " .. table.concat(user_exp.advanced_features, ", "))
                    log.info("  🔄 Features degradadas: " .. table.concat(user_exp.degraded_features, ", "))
                    
                    state.set("reliability_metrics", metrics)
                    
                    return true, "Padrões de fallback demonstrados"
                end
            },
            
            {
                name = "demonstrate_timeout_patterns",
                description = "Demonstra diferentes estratégias de timeout",
                depends_on = "demonstrate_fallback_patterns",
                command = function()
                    log.info("⏰ Demonstrando padrões de Timeout...")
                    
                    local http = require("http")
                    local reliability = require("reliability")
                    local services = state.get("test_services")
                    
                    -- 1. Timeout Simples
                    log.info("📊 1. Timeout Simples (3 segundos):")
                    local simple_start = os.time()
                    local simple_result = reliability.with_timeout(function()
                        return http.get({
                            url = services.slow_service.url, -- 10s delay, vai dar timeout
                            timeout = 15 -- timeout do HTTP maior que o do reliability
                        })
                    end, 3000) -- 3 segundos
                    local simple_end = os.time()
                    
                    if simple_result.success then
                        log.info("  ✅ Operação concluída")
                    else
                        log.info("  ⏰ Timeout após " .. (simple_end - simple_start) .. " segundos: " .. simple_result.error)
                    end
                    
                    -- 2. Timeout Adaptativo baseado na latência histórica
                    log.info("\n📊 2. Timeout Adaptativo:")
                    
                    -- Simular histórico de latências
                    local latency_history = {120, 150, 90, 200, 110, 180, 95} -- ms
                    local avg_latency = 0
                    for _, latency in ipairs(latency_history) do
                        avg_latency = avg_latency + latency
                    end
                    avg_latency = avg_latency / #latency_history
                    
                    local adaptive_timeout = avg_latency * 3 -- 3x a latência média + buffer
                    log.info("  📈 Latência média histórica: " .. string.format("%.1f", avg_latency) .. "ms")
                    log.info("  ⏰ Timeout adaptativo calculado: " .. adaptive_timeout .. "ms")
                    
                    local adaptive_result = reliability.with_timeout(function()
                        return http.get({
                            url = services.stable_service.url,
                            timeout = 5
                        })
                    end, adaptive_timeout)
                    
                    if adaptive_result.success then
                        log.info("  ✅ Operação bem-sucedida dentro do timeout adaptativo")
                    else
                        log.info("  ⏰ Timeout adaptativo: " .. adaptive_result.error)
                    end
                    
                    -- 3. Timeout Hierárquico (diferentes timeouts para diferentes operações)
                    log.info("\n📊 3. Timeout Hierárquico:")
                    
                    local operation_timeouts = {
                        critical = 1000,    -- 1s para operações críticas
                        standard = 5000,    -- 5s para operações padrão
                        background = 30000  -- 30s para operações em background
                    }
                    
                    local execute_with_priority = function(operation_type, operation_func)
                        local timeout = operation_timeouts[operation_type] or operation_timeouts.standard
                        log.info("  🎯 Executando operação '" .. operation_type .. "' (timeout: " .. timeout .. "ms)")
                        
                        return reliability.with_timeout(operation_func, timeout)
                    end
                    
                    -- Operação crítica (deve ser rápida)
                    local critical_result = execute_with_priority("critical", function()
                        return http.get({
                            url = services.stable_service.url,
                            timeout = 2
                        })
                    end)
                    
                    if critical_result.success then
                        log.info("  ✅ Operação crítica concluída rapidamente")
                    else
                        log.info("  ❌ Operação crítica falhou: " .. critical_result.error)
                    end
                    
                    -- Operação background (pode demorar mais)
                    local background_result = execute_with_priority("background", function()
                        -- Simular operação lenta mas que vai completar
                        exec.run("sleep 2")
                        return {status = "background operation completed"}
                    end)
                    
                    if background_result.success then
                        log.info("  ✅ Operação background concluída: " .. background_result.data.status)
                    else
                        log.info("  ❌ Operação background falhou: " .. background_result.error)
                    end
                    
                    return true, "Padrões de timeout demonstrados"
                end
            },
            
            {
                name = "generate_reliability_report",
                description = "Gera relatório completo de confiabilidade",
                depends_on = "demonstrate_timeout_patterns",
                command = function()
                    log.info("📊 Gerando relatório de confiabilidade...")
                    
                    local metrics = state.get("reliability_metrics")
                    local services = state.get("test_services")
                    
                    -- Calcular estatísticas
                    local success_rate = (metrics.successful_requests / math.max(metrics.total_requests, 1)) * 100
                    local failure_rate = (metrics.failed_requests / math.max(metrics.total_requests, 1)) * 100
                    
                    local report = {
                        timestamp = os.date("%Y-%m-%d %H:%M:%S"),
                        summary = {
                            total_requests = metrics.total_requests,
                            successful_requests = metrics.successful_requests,
                            failed_requests = metrics.failed_requests,
                            success_rate = success_rate,
                            failure_rate = failure_rate
                        },
                        patterns_demonstrated = {
                            retry_strategies = 3,
                            circuit_breaker_configs = 1,
                            fallback_mechanisms = 3,
                            timeout_patterns = 3
                        },
                        resilience_metrics = {
                            circuit_breaker_trips = metrics.circuit_breaker_trips,
                            fallback_executions = metrics.fallback_executions,
                            retry_attempts = metrics.retry_attempts
                        }
                    }
                    
                    -- Criar relatório detalhado
                    local report_content = string.format([[
🛡️  Relatório de Confiabilidade e Resiliência
==========================================

Data/Hora: %s

📊 Resumo das Operações:
- Total de Requests: %d
- Requests Bem-sucedidos: %d (%.1f%%)
- Requests Falharam: %d (%.1f%%)

🔄 Padrões Implementados:
- Estratégias de Retry: %d
  • Linear (intervalo fixo)
  • Exponencial (com jitter)
  • Condicional (customizado)

⚡ Circuit Breaker:
- Configurações Testadas: %d
- Trips Ativados: %d
- Proteção contra falhas em cascata ✅

🔄 Fallback Mechanisms:
- Tipos Implementados: %d
  • Cache local
  • Múltiplas estratégias
  • Degradação graceful
- Execuções de Fallback: %d

⏰ Timeout Patterns:
- Estratégias Testadas: %d
  • Timeout simples
  • Timeout adaptativo
  • Timeout hierárquico

🎯 Métricas de Resiliência:
- Circuit Breaker Trips: %d
- Fallback Executions: %d  
- Retry Attempts: %d

✨ Benefícios Demonstrados:
- 🛡️  Proteção contra falhas em cascata
- 🔄 Recovery automático de falhas temporárias
- ⚡ Resposta rápida mesmo com serviços instáveis
- 📉 Degradação graceful de funcionalidade
- 🎯 Timeouts otimizados por contexto

🏆 Conclusão:
O sistema demonstrou alta resiliência através da implementação
de múltiplos padrões de confiabilidade, garantindo operação
estável mesmo com serviços externos instáveis.
]], 
                        report.timestamp,
                        report.summary.total_requests,
                        report.summary.successful_requests,
                        report.summary.success_rate,
                        report.summary.failed_requests, 
                        report.summary.failure_rate,
                        report.patterns_demonstrated.retry_strategies,
                        report.patterns_demonstrated.circuit_breaker_configs,
                        report.resilience_metrics.circuit_breaker_trips,
                        report.patterns_demonstrated.fallback_mechanisms,
                        report.resilience_metrics.fallback_executions,
                        report.patterns_demonstrated.timeout_patterns,
                        report.resilience_metrics.circuit_breaker_trips,
                        report.resilience_metrics.fallback_executions,
                        report.resilience_metrics.retry_attempts
                    )
                    
                    fs.write("reliability_patterns_report.md", report_content)
                    state.set("reliability_report", report)
                    
                    log.info("📈 Estatísticas Finais:")
                    log.info("  📊 Taxa de sucesso: " .. string.format("%.1f", success_rate) .. "%")
                    log.info("  🔄 Execuções de fallback: " .. metrics.fallback_executions)
                    log.info("  ⚡ Trips do circuit breaker: " .. metrics.circuit_breaker_trips)
                    log.info("  🔄 Tentativas de retry: " .. metrics.retry_attempts)
                    
                    log.info("✅ Relatório salvo: reliability_patterns_report.md")
                    
                    return true, "Relatório de confiabilidade gerado"
                end
            }
        }
    }
}