-- Teste simples das funcionalidades GitOps
local gitops = require("gitops")
local log = require("log")

log.info("🔄 Testando GitOps Nativo...")

-- Teste básico de criação de workflow
local workflow_result = gitops.workflow({
    repo = "https://github.com/test/repo",
    branch = "main",
    auto_sync = true,
    diff_preview = true,
    rollback_on_failure = true
})

if workflow_result then
    log.info("✅ GitOps workflow criado com sucesso!")
    log.info("📋 Workflow ID: " .. workflow_result.workflow_id)
    log.info("📦 Repository ID: " .. workflow_result.repository_id)
    log.info("🔄 Auto-sync: " .. tostring(workflow_result.auto_sync))
else
    log.error("❌ Falha ao criar GitOps workflow")
end

-- Teste de status
if workflow_result then
    local status = gitops.get_workflow_status(workflow_result.workflow_id)
    if status then
        log.info("📊 Status do workflow: " .. status.status)
    end
end

-- Teste de listagem
local workflows = gitops.list_workflows()
log.info("📋 Total de workflows: " .. #workflows)

log.info("🎉 Teste GitOps concluído!")

-- Workflow simples para testar
workflow.define("test_gitops", {
    description = "Teste das funcionalidades GitOps",
    tasks = {
        {
            name = "test_task",
            description = "Tarefa de teste",
            command = function()
                log.info("✅ GitOps está funcionando!")
                return {success = true, message = "GitOps test completed"}
            end
        }
    }
})