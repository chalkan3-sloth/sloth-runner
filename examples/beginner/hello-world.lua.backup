-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03

local hello_task = task("say_hello")
local verify_task = task("verify_state")

local hello_task = task("say_hello")
local verify_task = task("verify_state")
local hello_task = task("say_hello")
    :description("Says hello to the world with modern DSL")
    :command(function()
        log.info("🌟 Modern DSL: Olá, mundo do Sloth Runner!")
        log.info("🦥 Este é seu primeiro script funcionando com Modern DSL!")
        
        -- Enhanced state management
        local greeting_data = {
            greeting = "Olá do Sloth Runner Modern DSL!",
            start_time = os.date("%Y-%m-%d %H:%M:%S"),
            dsl_version = "2.0.0",
            features = {"modern-syntax", "fluent-api", "enhanced-logging"}
        }
        
        state.set("greeting", greeting_data.greeting)
        state.set("start_time", greeting_data.start_time)
        state.set("dsl_version", greeting_data.dsl_version)
        
        log.info("📊 Dados salvos no state: " .. data.to_json(greeting_data))
        
        return true, "Modern DSL Hello World executado com sucesso!", greeting_data
    end)
    :timeout("30s")
    :on_success(function(params, output)
        log.info("✨ Modern DSL task completed at: " .. output.start_time)
    end)
    :build()
local verify_task = task("verify_state")
    :description("Verifies state data with modern DSL")
    :depends_on({"say_hello"})
    :command(function()
        log.info("🔍 Modern DSL: Verificando dados do state...")
        
        local greeting = state.get("greeting")
        local start_time = state.get("start_time")
        local dsl_version = state.get("dsl_version")
        
        log.info("Greeting: " .. (greeting or "N/A"))
        log.info("Start Time: " .. (start_time or "N/A"))
        log.info("DSL Version: " .. (dsl_version or "N/A"))
        
        -- Enhanced verification
        local verification_results = {
            greeting_exists = state.exists("greeting"),
            start_time_exists = state.exists("start_time"),
            dsl_version_exists = state.exists("dsl_version"),
            all_data_present = true
        }
        
        verification_results.all_data_present = verification_results.greeting_exists and 
                                              verification_results.start_time_exists and 
                                              verification_results.dsl_version_exists
        
        if verification_results.all_data_present then
            log.info("✅ Modern DSL: Todos os dados verificados com sucesso!")
        else
            log.warn("⚠️ Modern DSL: Alguns dados estão faltando")
        end
        
        return verification_results.all_data_present, 
               "State verification " .. (verification_results.all_data_present and "passed" or "failed"),
               verification_results
    end)
    :build()

workflow.define("hello_world_modern", {
    description = "Hello World demonstration - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        category = "beginner",
        tags = {"hello-world", "beginner", "state", "modern-dsl"},
        author = "Sloth Runner Team"
    },
    
    tasks = {
        hello_task,
        verify_task
    },
    
    config = {
        timeout = "5m",
        clean_workdir_after_run = true
    },
    
    on_start = function()
        log.info("🚀 Starting hello world workflow with Modern DSL...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✨ Hello World workflow completed successfully!")
            log.info("🎯 Modern DSL demonstration completed")
        end
        return true
    end
})
