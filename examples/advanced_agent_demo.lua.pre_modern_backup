-- Advanced Agent Capabilities Demo
-- Demonstra as novas funcionalidades de métricas, monitoramento e controle de qualidade

TaskDefinitions = {
    advanced_agent_demo = {
        description = "Demonstração das capacidades avançadas dos agentes com métricas e monitoramento",
        tasks = {
            -- Task 1: Monitoramento de Sistema
            system_monitoring = {
                name = "system_monitoring",
                description = "Coleta e análise de métricas do sistema",
                command = function()
                    log.info("🔍 Iniciando monitoramento avançado do sistema...")
                    
                    -- Coletar métricas do sistema
                    local cpu_usage = metrics.system_cpu()
                    local memory_info = metrics.system_memory()
                    local disk_info = metrics.system_disk()
                    local runtime_info = metrics.runtime_info()
                    
                    log.info("📊 Métricas do Sistema:")
                    log.info("   CPU: " .. string.format("%.1f%%", cpu_usage))
                    log.info("   Memória: " .. string.format("%.1f%% (%.0f/%.0f MB)", 
                        memory_info.percent, memory_info.used_mb, memory_info.total_mb))
                    log.info("   Disco: " .. string.format("%.1f%% (%.1f/%.1f GB)", 
                        disk_info.percent, disk_info.used_gb, disk_info.total_gb))
                    log.info("   Goroutines: " .. runtime_info.goroutines)
                    log.info("   Heap: " .. string.format("%.1f MB", runtime_info.heap_alloc_mb))
                    
                    -- Registrar métricas customizadas
                    metrics.gauge("agent_cpu_usage", cpu_usage, {environment = "production"})
                    metrics.gauge("agent_memory_percent", memory_info.percent)
                    metrics.gauge("agent_disk_percent", disk_info.percent)
                    
                    -- Verificar thresholds e alertas
                    if cpu_usage > 80 then
                        metrics.alert("high_cpu_usage", {
                            level = "warning",
                            message = "CPU usage is high: " .. string.format("%.1f%%", cpu_usage),
                            threshold = 80,
                            value = cpu_usage
                        })
                    end
                    
                    if memory_info.percent > 85 then
                        metrics.alert("high_memory_usage", {
                            level = "critical",
                            message = "Memory usage is critical: " .. string.format("%.1f%%", memory_info.percent),
                            threshold = 85,
                            value = memory_info.percent
                        })
                    end
                    
                    return true, "System monitoring completed successfully"
                end,
            },
            
            -- Task 2: Performance Benchmarking
            performance_benchmark = {
                name = "performance_benchmark",
                description = "Executa benchmarks de performance e registra métricas",
                depends_on = "system_monitoring",
                command = function()
                    log.info("🏃‍♂️ Iniciando benchmark de performance...")
                    
                    -- Benchmark de CPU
                    local cpu_benchmark_time = metrics.timer("cpu_benchmark", function()
                        -- Simular carga de CPU
                        local result = 0
                        for i = 1, 1000000 do
                            result = result + math.sqrt(i)
                        end
                        log.info("CPU benchmark result: " .. result)
                    end)
                    
                    log.info("⏱️  CPU Benchmark: " .. string.format("%.2f ms", cpu_benchmark_time))
                    
                    -- Benchmark de I/O
                    local io_benchmark_time = metrics.timer("io_benchmark", function()
                        -- Criar arquivo de teste
                        local test_data = string.rep("benchmark_data_", 1000)
                        fs.write("/tmp/benchmark_test.txt", test_data)
                        
                        -- Ler arquivo várias vezes
                        for i = 1, 100 do
                            local data = fs.read("/tmp/benchmark_test.txt")
                            if not data then
                                error("Failed to read benchmark file")
                            end
                        end
                        
                        -- Limpar arquivo de teste
                        exec.run("rm -f /tmp/benchmark_test.txt")
                    end)
                    
                    log.info("⏱️  I/O Benchmark: " .. string.format("%.2f ms", io_benchmark_time))
                    
                    -- Registrar métricas de performance
                    metrics.histogram("task_performance_cpu", cpu_benchmark_time, {type = "cpu", unit = "ms"})
                    metrics.histogram("task_performance_io", io_benchmark_time, {type = "io", unit = "ms"})
                    
                    -- Contador de benchmarks executados
                    local benchmark_count = metrics.counter("benchmarks_executed", 1, {agent = "myagent1"})
                    log.info("📈 Total benchmarks executados: " .. benchmark_count)
                    
                    return true, "Performance benchmark completed"
                end,
            },
            
            -- Task 3: Health Check Avançado
            advanced_health_check = {
                name = "advanced_health_check",
                description = "Realiza verificações de saúde abrangentes do sistema",
                depends_on = "performance_benchmark",
                command = function()
                    log.info("🏥 Executando health check avançado...")
                    
                    local health = metrics.health_status()
                    local overall_status = health.overall
                    
                    log.info("🔍 Status de Saúde Geral: " .. overall_status)
                    
                    -- Verificar cada componente
                    local components = {"cpu", "memory", "disk"}
                    for _, component in ipairs(components) do
                        local comp_info = health[component]
                        if comp_info then
                            local status_icon = "✅"
                            if comp_info.status == "warning" then
                                status_icon = "⚠️"
                            elseif comp_info.status == "critical" then
                                status_icon = "❌"
                            end
                            
                            log.info(string.format("   %s %s: %.1f%% (%s)", 
                                status_icon, component:upper(), comp_info.usage, comp_info.status))
                        end
                    end
                    
                    -- Verificações customizadas
                    local custom_checks = {
                        {
                            name = "disk_space_check",
                            command = "df -h / | tail -n 1 | awk '{print $5}' | sed 's/%//'"
                        },
                        {
                            name = "network_connectivity",
                            command = "ping -c 1 8.8.8.8 >/dev/null 2>&1 && echo 'OK' || echo 'FAILED'"
                        },
                        {
                            name = "process_count",
                            command = "ps aux | wc -l"
                        }
                    }
                    
                    for _, check in ipairs(custom_checks) do
                        local stdout, stderr, failed = exec.run(check.command)
                        if not failed then
                            local result = string.gsub(stdout, "\n", "")
                            log.info("✓ " .. check.name .. ": " .. result)
                            metrics.gauge("health_check_" .. check.name, tonumber(result) or 1, {status = "ok"})
                        else
                            log.error("✗ " .. check.name .. " failed: " .. stderr)
                            metrics.gauge("health_check_" .. check.name, 0, {status = "failed"})
                        end
                    end
                    
                    -- Registrar health score
                    local health_score = 100
                    if overall_status == "warning" then
                        health_score = 75
                    elseif overall_status == "critical" then
                        health_score = 25
                    end
                    
                    metrics.gauge("agent_health_score", health_score, {agent = "myagent1"})
                    
                    return true, "Health check completed - Status: " .. overall_status
                end,
            },
            
            -- Task 4: Load Testing Simulation
            load_testing = {
                name = "load_testing",
                description = "Simula carga de trabalho e mede capacidade do agente",
                depends_on = "advanced_health_check",
                command = function()
                    log.info("🔥 Iniciando simulação de load testing...")
                    
                    local concurrent_tasks = 10
                    local task_duration = 2 -- segundos
                    
                    -- Registrar início do load test
                    metrics.gauge("load_test_active", 1)
                    metrics.counter("load_tests_started", 1)
                    
                    local start_time = os.time()
                    
                    -- Simular múltiplas tarefas concorrentes
                    for i = 1, concurrent_tasks do
                        local task_time = metrics.timer("simulated_task_" .. i, function()
                            -- Simular trabalho (CPU + I/O)
                            local result = 0
                            for j = 1, 100000 do
                                result = result + math.sin(j)
                            end
                            
                            -- Simular I/O
                            exec.run("sleep " .. (task_duration / concurrent_tasks))
                            
                            return result
                        end)
                        
                        log.info(string.format("📋 Task %d completed in %.2f ms", i, task_time))
                        metrics.histogram("concurrent_task_duration", task_time, {task_id = tostring(i)})
                    end
                    
                    local total_time = os.time() - start_time
                    log.info("⏰ Load test completed in " .. total_time .. " seconds")
                    
                    -- Calcular métricas de throughput
                    local throughput = concurrent_tasks / total_time
                    metrics.gauge("load_test_throughput", throughput, {unit = "tasks_per_second"})
                    
                    -- Verificar estado do sistema após load test
                    local post_health = metrics.health_status()
                    log.info("📊 Sistema após load test: " .. post_health.overall)
                    
                    metrics.gauge("load_test_active", 0)
                    metrics.counter("load_tests_completed", 1)
                    
                    return true, string.format("Load test completed - Throughput: %.2f tasks/sec", throughput)
                end,
            },
            
            -- Task 5: Quality Assurance e Relatório
            quality_assurance_report = {
                name = "quality_assurance_report",
                description = "Gera relatório completo de qualidade e performance",
                depends_on = "load_testing",
                command = function()
                    log.info("📋 Gerando relatório de Quality Assurance...")
                    
                    -- Coletar todas as métricas customizadas
                    local custom_metrics = metrics.list_custom()
                    log.info("📊 Métricas coletadas: " .. #custom_metrics)
                    
                    -- Gerar relatório JSON
                    local report = {
                        timestamp = os.date("%Y-%m-%d %H:%M:%S"),
                        agent_info = {
                            name = "myagent1",
                            version = "1.0.0",
                            runtime = metrics.runtime_info()
                        },
                        system_health = metrics.health_status(),
                        performance_summary = {},
                        quality_score = 0
                    }
                    
                    -- Calcular score de qualidade baseado nas métricas
                    local quality_factors = {
                        cpu_efficiency = 100 - metrics.system_cpu(),
                        memory_efficiency = 100 - metrics.system_memory().percent,
                        disk_efficiency = 100 - metrics.system_disk().percent,
                        stability = 100 -- Baseado em não ter crashes
                    }
                    
                    local total_score = 0
                    for factor, score in pairs(quality_factors) do
                        total_score = total_score + score
                        report.performance_summary[factor] = score
                    end
                    
                    report.quality_score = total_score / 4
                    
                    -- Salvar relatório
                    local report_json = data.to_json(report)
                    local report_filename = "/tmp/agent_qa_report_" .. os.date("%Y%m%d_%H%M%S") .. ".json"
                    fs.write(report_filename, report_json)
                    
                    log.info("✅ Relatório salvo em: " .. report_filename)
                    log.info("🏆 Quality Score: " .. string.format("%.1f/100", report.quality_score))
                    
                    -- Registrar métricas finais
                    metrics.gauge("qa_report_generated", 1)
                    metrics.gauge("overall_quality_score", report.quality_score)
                    
                    -- Determinar se o agente passa no QA
                    local qa_passed = report.quality_score >= 70
                    local status = qa_passed and "PASSED" or "FAILED"
                    
                    metrics.gauge("qa_test_result", qa_passed and 1 or 0, {status = status:lower()})
                    
                    log.info("🎯 QA Result: " .. status)
                    
                    return true, "QA Report completed - Score: " .. string.format("%.1f", report.quality_score) .. " - " .. status
                end,
            },
            
            -- Task 6: Cleanup e Otimização
            cleanup_and_optimization = {
                name = "cleanup_and_optimization",
                description = "Limpeza do sistema e otimizações baseadas em métricas",
                depends_on = "quality_assurance_report",
                command = function()
                    log.info("🧹 Iniciando limpeza e otimização do sistema...")
                    
                    -- Limpeza de arquivos temporários
                    local cleanup_start = os.time()
                    exec.run("find /tmp -name 'benchmark_*' -delete 2>/dev/null")
                    exec.run("find /tmp -name 'agent_*' -type f -mtime +1 -delete 2>/dev/null")
                    
                    -- Force garbage collection
                    collectgarbage("collect")
                    
                    local cleanup_time = os.time() - cleanup_start
                    log.info("🗑️  Limpeza concluída em " .. cleanup_time .. " segundos")
                    
                    -- Verificar melhorias pós-limpeza
                    local post_cleanup_memory = metrics.system_memory()
                    log.info("💾 Memória após limpeza: " .. string.format("%.1f%%", post_cleanup_memory.percent))
                    
                    -- Registrar métricas de otimização
                    metrics.counter("system_cleanups_performed", 1)
                    metrics.gauge("post_cleanup_memory_percent", post_cleanup_memory.percent)
                    
                    -- Recomendações de otimização baseadas em métricas
                    local recommendations = {}
                    
                    if post_cleanup_memory.percent > 80 then
                        table.insert(recommendations, "Consider increasing memory allocation")
                    end
                    
                    local cpu_usage = metrics.system_cpu()
                    if cpu_usage > 70 then
                        table.insert(recommendations, "High CPU usage detected - consider load balancing")
                    end
                    
                    if #recommendations > 0 then
                        log.info("💡 Recomendações de otimização:")
                        for i, rec in ipairs(recommendations) do
                            log.info("   " .. i .. ". " .. rec)
                        end
                        metrics.gauge("optimization_recommendations_count", #recommendations)
                    else
                        log.info("✅ Sistema está otimizado")
                        metrics.gauge("optimization_recommendations_count", 0)
                    end
                    
                    return true, "Cleanup and optimization completed"
                end,
            }
        }
    }
}