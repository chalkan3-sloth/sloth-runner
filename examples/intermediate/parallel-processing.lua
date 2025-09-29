-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03

local http = require("http")
local start = os.time()
local result = http.get({url = "%s", timeout = 10})
local duration = os.time() - start
local start = os.time()
local content = "Arquivo de teste #%d\\nCriado em PARALELO em: " .. os.date()
local duration = os.time() - start
local start = os.time()
local result = 0
local duration = os.time() - start

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
    parallel_processing = {
        description = "Demonstração de processamento paralelo e otimização de performance",
        
        tasks = {
            {
                name = "setup_parallel_demo",
                description = "Configura ambiente para demonstração paralela",
                command = function()
                    log.info("⚡ Configurando demonstração de processamento paralelo...")
                    
                    -- Criar dados de teste
                    local test_data = {
                        urls = {
                            "https://httpbin.org/delay/1",
                            "https://httpbin.org/delay/2", 
                            "https://httpbin.org/delay/1",
                            "https://jsonplaceholder.typicode.com/posts/1",
                            "https://jsonplaceholder.typicode.com/users/1"
                        },
                        files_to_create = {
                            "test_file_1.txt",
                            "test_file_2.txt", 
                            "test_file_3.txt",
                            "test_file_4.txt",
                            "test_file_5.txt"
                        },
                        numbers = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
                    }
                    
                    state.set("parallel_test_data", test_data)
                    
                    log.info("✅ Dados de teste preparados:")
                    log.info("  📡 " .. #test_data.urls .. " URLs para requisições HTTP")
                    log.info("  📁 " .. #test_data.files_to_create .. " arquivos para criar")
                    log.info("  🔢 " .. #test_data.numbers .. " números para processar")
                    
                    return true, "Ambiente configurado"
                end
            },
            
            {
                name = "sequential_processing_benchmark",
                description = "Executa processamento sequencial para comparação",
                depends_on = "setup_parallel_demo",
                command = function()
                    log.info("🐌 Executando processamento SEQUENCIAL (benchmark)...")
                    
                    local http = require("http")
                    local test_data = state.get("parallel_test_data")
                    local start_time = os.time()
                    
                    local results = {
                        http_requests = {},
                        file_operations = {},
                        calculations = {}
                    }
                    
                    -- 1. Requisições HTTP sequenciais
                    log.info("📡 Fazendo " .. #test_data.urls .. " requisições HTTP sequenciais...")
                    for i, url in ipairs(test_data.urls) do
                        local request_start = os.time()
                        local result = http.get({
                            url = url,
                            timeout = 10
                        })
                        local request_end = os.time()
                        
                        table.insert(results.http_requests, {
                            url = url,
                            success = result.success,
                            duration = request_end - request_start,
                            status = result.success and result.data.status_code or "error"
                        })
                        
                        log.info("  " .. i .. "/" .. #test_data.urls .. " - " .. (result.success and "✅" or "❌") .. " (" .. (request_end - request_start) .. "s)")
                    end
                    
                    -- 2. Operações de arquivo sequenciais
                    log.info("📁 Criando " .. #test_data.files_to_create .. " arquivos sequenciais...")
                    for i, filename in ipairs(test_data.files_to_create) do
                        local file_start = os.time()
                        local content = "Arquivo de teste #" .. i .. "\nCriado sequencialmente em: " .. os.date()
                        fs.write(filename, content)
                        exec.run("sleep 0.5")  -- Simular operação que demora
                        local file_end = os.time()
                        
                        table.insert(results.file_operations, {
                            filename = filename,
                            duration = file_end - file_start
                        })
                    end
                    
                    -- 3. Cálculos sequenciais
                    log.info("🔢 Processando " .. #test_data.numbers .. " números sequenciais...")
                    for i, number in ipairs(test_data.numbers) do
                        local calc_start = os.time()
                        
                        -- Simular cálculo pesado
                        local result = 0
                        for j = 1, number * 1000000 do
                            result = result + math.sin(j)
                        end
                        
                        local calc_end = os.time()
                        
                        table.insert(results.calculations, {
                            number = number,
                            result = result,
                            duration = calc_end - calc_start
                        })
                    end
                    
                    local end_time = os.time()
                    local total_duration = end_time - start_time
                    
                    results.total_duration = total_duration
                    state.set("sequential_results", results)
                    
                    log.info("🐌 Processamento sequencial concluído!")
                    log.info("⏱️  Tempo total: " .. total_duration .. " segundos")
                    
                    return true, "Processamento sequencial concluído"
                end
            },
            
            {
                name = "parallel_processing_demo",
                description = "Executa o mesmo processamento em paralelo",
                depends_on = "sequential_processing_benchmark",
                command = function()
                    log.info("⚡ Executando processamento PARALELO...")
                    
                    local test_data = state.get("parallel_test_data")
                    local start_time = os.time()
                    
                    -- Criar tarefas paralelas para requisições HTTP
                    local http_tasks = {}
                    for i, url in ipairs(test_data.urls) do
                        table.insert(http_tasks, {
                            name = "http_request_" .. i,
                            command = string.format([[
local http = require("http")
local start = os.time()
local result = http.get({url = "%s", timeout = 10})
local duration = os.time() - start
return {
    url = "%s",
    success = result.success,
    duration = duration,
    status = result.success and result.data.status_code or "error"
}
