-- 🌟 Hello World - Seu primeiro script Sloth Runner
-- Este é o exemplo mais básico que demonstra a estrutura fundamental

TaskDefinitions = {
    hello_world = {
        description = "Um simples Hello World para começar",
        
        tasks = {
            {
                name = "say_hello",
                description = "Diz olá para o mundo",
                command = function()
                    -- 📝 Log é usado para saída formatada
                    log.info("🌟 Olá, mundo do Sloth Runner!")
                    log.info("🦥 Este é seu primeiro script funcionando!")
                    
                    -- ✅ Sempre retorne true para sucesso, false para falha
                    return true, "Hello World executado com sucesso!"
                end
            },
            
            {
                name = "show_system_info",
                description = "Mostra informações do sistema",
                depends_on = "say_hello", -- 🔗 Esta task depende da anterior
                command = function()
                    log.info("📊 Informações do Sistema:")
                    
                    -- 🖥️ Executar comandos do sistema
                    local os_info = exec.run("uname -a")
                    if os_info.success then
                        log.info("Sistema: " .. os_info.stdout)
                    end
                    
                    local date_info = exec.run("date")
                    if date_info.success then
                        log.info("Data atual: " .. date_info.stdout)
                    end
                    
                    return true, "Informações do sistema coletadas"
                end
            },
            
            {
                name = "create_welcome_file",
                description = "Cria um arquivo de boas-vindas",
                depends_on = "show_system_info",
                command = function(params)
                    local welcome_text = [[
🦥 Bem-vindo ao Sloth Runner!

Este arquivo foi criado automaticamente pelo seu primeiro script.

Algumas coisas que você pode fazer:
- Explorar os exemplos em examples/
- Usar help() para ver ajuda interativa
- Criar seus próprios scripts de automação

Happy coding! 🚀
]]
                    
                    -- 📁 Usar o módulo fs para operações de arquivo
                    local filename = "welcome.txt"
                    fs.write(filename, welcome_text)
                    
                    -- ✅ Verificar se o arquivo foi criado
                    if fs.exists(filename) then
                        log.info("✅ Arquivo '" .. filename .. "' criado com sucesso!")
                        return true, "Arquivo de boas-vindas criado"
                    else
                        log.error("❌ Falha ao criar arquivo")
                        return false, "Erro ao criar arquivo"
                    end
                end
            }
        }
    }
}