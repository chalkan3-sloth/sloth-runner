-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03


-- local example_task = task("task_name")
--     :description("Task description with modern DSL")
--     :command(function(params, deps)
--         -- Enhanced task logic
--         return true, "Task completed", { result = "success" }
--     end)
--     :timeout("30s")
--     :build()

-- workflow.define("workflow_name", {
--     description = "Workflow description - Modern DSL",
--     version = "2.0.0",
--     tasks = { example_task },
--     config = { timeout = "10m" }
-- })

-- Maintain backward compatibility with legacy format
TaskDefinitions = {
    http_basics = {
        description = "Exemplos básicos de uso do módulo HTTP",
        
        tasks = {
            {
                name = "simple_get_request",
                description = "Faz uma requisição GET simples",
                command = function()
                    log.info("🌐 Fazendo requisição GET simples...")
                    
                    -- Carregar o módulo HTTP
                    local http = require("http")
                    
                    -- Fazer uma requisição GET simples
                    local result = http.get({
                        url = "https://api.github.com/zen"
                    })
                    
                    if result.success then
                        log.info("✅ Requisição bem-sucedida!")
                        log.info("📝 GitHub Zen: " .. result.data.body)
                        log.info("⏱️  Tempo de resposta: " .. result.data.elapsed_ms .. "ms")
                        log.info("📊 Status: " .. result.data.status_code)
                    else
                        log.error("❌ Falha na requisição: " .. result.error)
                        return false, "GET request falhou"
                    end
                    
                    return true, "GET request executado com sucesso"
                end
            },
            
            {
                name = "request_with_headers",
                description = "Requisição com headers customizados",
                depends_on = "simple_get_request",
                command = function()
                    log.info("🔧 Fazendo requisição com headers customizados...")
                    
                    local http = require("http")
                    
                    local result = http.get({
                        url = "https://httpbin.org/headers",
                        headers = {
                            ["User-Agent"] = "Sloth-Runner/1.0",
                            ["Accept"] = "application/json",
                            ["X-Custom-Header"] = "Hello from Sloth Runner!"
                        }
                    })
                    
                    if result.success then
                        log.info("✅ Requisição com headers bem-sucedida!")
                        
                        -- Tentar parsear a resposta JSON
                        if result.data.json then
                            log.info("📋 Headers recebidos:")
                            log.info(data.to_json(result.data.json.headers))
                        end
                    else
                        log.error("❌ Falha na requisição: " .. result.error)
                        return false, "Request com headers falhou"
                    end
                    
                    return true, "Request com headers executado"
                end
            },
            
            {
                name = "post_with_json",
                description = "Requisição POST com dados JSON",
                depends_on = "request_with_headers",
                command = function()
                    log.info("📤 Enviando dados JSON via POST...")
                    
                    local http = require("http")
                    
                    -- Dados para enviar
                    local post_data = {
                        name = "Sloth Runner User",
                        message = "Hello from automation!",
                        timestamp = os.time(),
                        features = {"easy", "powerful", "flexible"}
                    }
                    
                    local result = http.post({
                        url = "https://httpbin.org/post",
                        json = post_data,
                        headers = {
                            ["User-Agent"] = "Sloth-Runner-Example/1.0"
                        }
                    })
                    
                    if result.success then
                        log.info("✅ POST com JSON bem-sucedido!")
                        
                        if result.data.json and result.data.json.json then
                            local received_data = result.data.json.json
                            log.info("📨 Dados recebidos de volta:")
                            log.info("  Nome: " .. received_data.name)
                            log.info("  Mensagem: " .. received_data.message)
                            log.info("  Features: " .. table.concat(received_data.features, ", "))
                        end
                    else
                        log.error("❌ Falha no POST: " .. result.error)
                        return false, "POST request falhou"
                    end
                    
                    return true, "POST com JSON executado"
                end
            },
            
            {
                name = "error_handling_demo",
                description = "Demonstração de tratamento de erros HTTP",
                depends_on = "post_with_json",
                command = function()
                    log.info("🚨 Demonstrando tratamento de erros HTTP...")
                    
                    local http = require("http")
                    
                    -- Tentar acessar uma URL que retorna 404
                    local result = http.get({
                        url = "https://httpbin.org/status/404"
                    })
                    
                    if result.success then
                        if result.data.status_code == 404 then
                            log.info("✅ Erro 404 tratado corretamente!")
                            log.info("📊 Status Code: " .. result.data.status_code)
                        else
                            log.info("🤔 Status inesperado: " .. result.data.status_code)
                        end
                    else
                        log.error("❌ Erro na requisição: " .. result.error)
                    end
                    
                    -- Tentar acessar uma URL inválida
                    local bad_result = http.get({
                        url = "https://this-domain-does-not-exist-12345.com",
                        timeout = 5  -- 5 segundos de timeout
                    })
                    
                    if not bad_result.success then
                        log.info("✅ Erro de conexão tratado corretamente!")
                        log.info("🔍 Tipo de erro: " .. bad_result.error)
                    end
                    
                    return true, "Tratamento de erros demonstrado"
                end
            }
        }
    }
}
