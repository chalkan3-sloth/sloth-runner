-- 💾 State Basics - Aprendendo gerenciamento de estado
-- Este exemplo demonstra como usar o sistema de estado do Sloth Runner

TaskDefinitions = {
    state_basics = {
        description = "Exemplos básicos de gerenciamento de estado",
        
        tasks = {
            {
                name = "basic_state_operations",
                description = "Operações básicas de estado (set, get, delete)",
                command = function()
                    log.info("💾 Demonstrando operações básicas de estado...")
                    
                    -- 📝 Definir valores de diferentes tipos
                    state.set("app_name", "Sloth Runner")
                    state.set("version", "1.0.0")
                    state.set("user_count", 1337)
                    state.set("is_production", false)
                    
                    -- 📚 Definir estruturas complexas
                    local config = {
                        database = {
                            host = "localhost",
                            port = 5432,
                            name = "myapp"
                        },
                        features = {"auth", "api", "monitoring"},
                        limits = {
                            max_connections = 100,
                            timeout = 30
                        }
                    }
                    state.set("app_config", config)
                    
                    log.info("✅ Valores definidos no estado!")
                    
                    -- 📖 Ler valores do estado
                    local app_name = state.get("app_name")
                    local version = state.get("version")
                    local user_count = state.get("user_count")
                    local is_prod = state.get("is_production")
                    
                    log.info("📋 Valores lidos do estado:")
                    log.info("  App: " .. app_name .. " v" .. version)
                    log.info("  Usuários: " .. user_count)
                    log.info("  Produção: " .. tostring(is_prod))
                    
                    -- 🔍 Ler estrutura complexa
                    local saved_config = state.get("app_config")
                    log.info("  Database: " .. saved_config.database.host .. ":" .. saved_config.database.port)
                    log.info("  Features: " .. table.concat(saved_config.features, ", "))
                    
                    return true, "Operações básicas de estado executadas"
                end
            },
            
            {
                name = "default_values",
                description = "Usando valores padrão quando chaves não existem",
                depends_on = "basic_state_operations",
                command = function()
                    log.info("🎯 Demonstrando valores padrão...")
                    
                    -- 📖 Tentar ler uma chave que não existe (com padrão)
                    local theme = state.get("ui_theme", "dark")
                    local max_retries = state.get("max_retries", 3)
                    local debug_mode = state.get("debug_mode", true)
                    
                    log.info("🎨 Tema UI: " .. theme .. " (valor padrão)")
                    log.info("🔄 Max retries: " .. max_retries .. " (valor padrão)")
                    log.info("🐛 Debug mode: " .. tostring(debug_mode) .. " (valor padrão)")
                    
                    -- ✅ Verificar se uma chave existe
                    if state.exists("app_name") then
                        log.info("✅ 'app_name' existe no estado")
                    end
                    
                    if not state.exists("missing_key") then
                        log.info("❌ 'missing_key' não existe no estado")
                    end
                    
                    return true, "Valores padrão demonstrados"
                end
            },
            
            {
                name = "state_persistence",
                description = "Demonstra que o estado persiste entre execuções",
                depends_on = "default_values",
                command = function()
                    log.info("💿 Demonstrando persistência de estado...")
                    
                    -- 📊 Incrementar um contador de execuções
                    local execution_count = state.get("execution_count", 0)
                    execution_count = execution_count + 1
                    state.set("execution_count", execution_count)
                    
                    log.info("🔢 Esta é a execução número: " .. execution_count)
                    
                    -- ⏰ Salvar timestamp da última execução
                    local last_run = state.get("last_execution_time")
                    local current_time = os.time()
                    
                    if last_run then
                        local time_diff = current_time - last_run
                        log.info("⏱️  Tempo desde a última execução: " .. time_diff .. " segundos")
                    else
                        log.info("🆕 Esta é a primeira execução!")
                    end
                    
                    state.set("last_execution_time", current_time)
                    
                    -- 📝 Manter um log de execuções
                    local execution_log = state.get("execution_log", {})
                    table.insert(execution_log, {
                        run_number = execution_count,
                        timestamp = current_time,
                        date = os.date("%Y-%m-%d %H:%M:%S", current_time)
                    })
                    state.set("execution_log", execution_log)
                    
                    log.info("📚 Log de execuções atualizado")
                    
                    return true, "Persistência de estado demonstrada"
                end
            },
            
            {
                name = "cleanup_demo",
                description = "Demonstra como limpar dados desnecessários",
                depends_on = "state_persistence",
                command = function()
                    log.info("🧹 Demonstrando limpeza de estado...")
                    
                    -- 🗂️ Criar algumas chaves temporárias
                    state.set("temp_data_1", "temporary value 1")
                    state.set("temp_data_2", "temporary value 2")
                    state.set("temp_config", {key = "value"})
                    
                    -- 🔍 Verificar que as chaves existem
                    log.info("📋 Antes da limpeza:")
                    log.info("  temp_data_1 exists: " .. tostring(state.exists("temp_data_1")))
                    log.info("  temp_data_2 exists: " .. tostring(state.exists("temp_data_2")))
                    log.info("  temp_config exists: " .. tostring(state.exists("temp_config")))
                    
                    -- 🗑️ Deletar chaves individuais
                    state.delete("temp_data_1")
                    log.info("🗑️  Deletada: temp_data_1")
                    
                    state.delete("temp_data_2")
                    log.info("🗑️  Deletada: temp_data_2")
                    
                    -- ✅ Verificar depois da limpeza
                    log.info("📋 Depois da limpeza:")
                    log.info("  temp_data_1 exists: " .. tostring(state.exists("temp_data_1")))
                    log.info("  temp_data_2 exists: " .. tostring(state.exists("temp_data_2")))
                    log.info("  temp_config exists: " .. tostring(state.exists("temp_config")))
                    
                    -- 🧽 Limpar a chave restante
                    state.delete("temp_config")
                    log.info("🗑️  Deletada: temp_config")
                    
                    log.info("✅ Limpeza concluída!")
                    
                    return true, "Limpeza de estado demonstrada"
                end
            }
        }
    }
}