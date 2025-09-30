-- GitOps Nativo - Exemplo Básico
-- Demonstra o uso da funcionalidade GitOps integrada no Sloth Runner

local gitops = require("gitops")
local log = require("log")

log.info("🔄 GitOps Nativo - Demonstração")
log.info("=" .. string.rep("=", 50))

-- Exemplo 1: GitOps Workflow Simples
log.info("📦 Criando GitOps workflow simples...")

local simple_workflow = gitops.workflow({
    repo = "https://github.com/company/infrastructure",
    branch = "main",
    auto_sync = true,
    diff_preview = true,
    rollback_on_failure = true
})

if simple_workflow then
    log.info("✅ GitOps workflow criado:")
    log.info("  ID: " .. simple_workflow.workflow_id)
    log.info("  Repository: " .. simple_workflow.repository_id)
    log.info("  Auto-sync: " .. tostring(simple_workflow.auto_sync))
    log.info("  Rollback on failure: " .. tostring(simple_workflow.rollback_on_failure))
else
    log.error("❌ Falha ao criar GitOps workflow")
end

-- Exemplo 2: Registrar repositório personalizado
log.info("")
log.info("🏗️ Registrando repositório personalizado...")

local repo_id = gitops.register_repository({
    url = "https://github.com/company/k8s-manifests",
    branch = "production",
    id = "production-repo"
})

if repo_id then
    log.info("✅ Repositório registrado: " .. repo_id)
else
    log.error("❌ Falha ao registrar repositório")
end

-- Exemplo 3: Criar workflow detalhado
log.info("")
log.info("⚙️ Criando workflow GitOps detalhado...")

local detailed_workflow_id = gitops.create_workflow({
    name = "Production Infrastructure",
    repository = repo_id or "production-repo",
    target_path = "k8s/production",
    auto_sync = false,  -- Sync manual para produção
    diff_preview = true,
    rollback_on_failure = true
})

if detailed_workflow_id then
    log.info("✅ Workflow detalhado criado: " .. detailed_workflow_id)
    
    -- Exemplo 4: Gerar diff preview
    log.info("")
    log.info("🔍 Gerando preview de mudanças...")
    
    local diff = gitops.generate_diff(detailed_workflow_id)
    if diff then
        log.info("📊 Diff Summary:")
        log.info("  Total changes: " .. diff.summary.total_changes)
        log.info("  Created resources: " .. diff.summary.created_resources)
        log.info("  Updated resources: " .. diff.summary.updated_resources)
        log.info("  Deleted resources: " .. diff.summary.deleted_resources)
        log.info("  Conflicts: " .. diff.summary.conflict_count)
        
        if #diff.changes > 0 then
            log.info("📝 Changes detected:")
            for i, change in ipairs(diff.changes) do
                log.info("  " .. i .. ". " .. change.type .. " " .. change.resource .. " (" .. change.impact .. " impact)")
            end
        end
        
        if #diff.warnings > 0 then
            log.warn("⚠️ Warnings:")
            for i, warning in ipairs(diff.warnings) do
                log.warn("  " .. i .. ". " .. warning)
            end
        end
    else
        log.warn("ℹ️ Nenhuma mudança detectada")
    end
    
    -- Exemplo 5: Sync manual (após revisar diff)
    log.info("")
    log.info("🚀 Executando sync manual...")
    
    local sync_success = gitops.sync_workflow(detailed_workflow_id)
    if sync_success then
        log.info("✅ Sync executado com sucesso!")
        
        -- Verificar status após sync
        local status = gitops.get_workflow_status(detailed_workflow_id)
        if status then
            log.info("📊 Status do workflow:")
            log.info("  Status: " .. status.status)
            if status.last_sync_result then
                log.info("  Last sync: " .. status.last_sync_result.status)
                log.info("  Commit: " .. status.last_sync_result.commit_hash)
                if status.last_sync_result.metrics then
                    log.info("  Duration: " .. status.last_sync_result.metrics.duration)
                    log.info("  Resources applied: " .. status.last_sync_result.metrics.resources_applied)
                end
            end
        end
    else
        log.error("❌ Sync falhou!")
        
        -- Exemplo 6: Rollback em caso de falha
        log.info("🔄 Iniciando rollback...")
        local rollback_success = gitops.rollback_workflow(detailed_workflow_id, "Sync failed, rolling back")
        if rollback_success then
            log.info("✅ Rollback executado com sucesso!")
        else
            log.error("❌ Rollback falhou!")
        end
    end
else
    log.error("❌ Falha ao criar workflow detalhado")
end

-- Exemplo 7: Listar todos os workflows
log.info("")
log.info("📋 Listando todos os workflows GitOps...")

local workflows = gitops.list_workflows()
if #workflows > 0 then
    log.info("📊 " .. #workflows .. " workflow(s) encontrado(s):")
    for i, workflow in ipairs(workflows) do
        log.info("  " .. i .. ". " .. workflow.name .. " (" .. workflow.id .. ")")
        log.info("     Status: " .. workflow.status)
        log.info("     Repository: " .. workflow.repository)
    end
else
    log.info("ℹ️ Nenhum workflow encontrado")
end

-- Exemplo 8: Iniciar auto-sync para todos os workflows
log.info("")
log.info("🔄 Iniciando auto-sync controller...")

local auto_sync_started = gitops.start_auto_sync()
if auto_sync_started then
    log.info("✅ Auto-sync controller iniciado!")
    log.info("🔄 Workflows com auto_sync=true serão sincronizados automaticamente")
else
    log.error("❌ Falha ao iniciar auto-sync controller")
end

log.info("")
log.info("🎉 Demonstração GitOps Nativo concluída!")
log.info("📊 Funcionalidades demonstradas:")
log.info("  ✅ Criação de workflows GitOps")
log.info("  ✅ Registro de repositórios")
log.info("  ✅ Preview de mudanças (diff)")
log.info("  ✅ Sync manual e automático")
log.info("  ✅ Rollback inteligente")
log.info("  ✅ Monitoramento de status")
log.info("  ✅ Auto-sync controller")

-- Definir um workflow que usa GitOps
workflow.define("gitops_demo", {
    description = "Demonstração do GitOps Nativo",
    version = "1.0.0",
    
    tasks = {
        {
            name = "setup_gitops",
            description = "Configurar GitOps workflows",
            command = function(params, deps)
                log.info("🔧 Configurando GitOps workflows...")
                
                -- Criar workflow para diferentes ambientes
                local environments = {"development", "staging", "production"}
                
                for _, env in ipairs(environments) do
                    local workflow_config = {
                        repo = "https://github.com/company/" .. env .. "-config",
                        branch = env == "production" and "main" or "develop",
                        auto_sync = env ~= "production", -- Prod é manual
                        diff_preview = true,
                        rollback_on_failure = true
                    }
                    
                    local result = gitops.workflow(workflow_config)
                    if result then
                        log.info("✅ " .. env .. " workflow: " .. result.workflow_id)
                    else
                        log.error("❌ Falha ao criar workflow " .. env)
                    end
                end
                
                return {success = true, message = "GitOps workflows configurados"}
            end
        },
        
        {
            name = "monitor_gitops",
            description = "Monitorar status dos workflows GitOps",
            command = function(params, deps)
                log.info("📊 Monitorando workflows GitOps...")
                
                local workflows = gitops.list_workflows()
                local healthy_count = 0
                local total_count = #workflows
                
                for _, workflow in ipairs(workflows) do
                    local status = gitops.get_workflow_status(workflow.id)
                    if status and (status.status == "synced" or status.status == "active") then
                        healthy_count = healthy_count + 1
                        log.info("✅ " .. workflow.name .. ": " .. status.status)
                    else
                        log.warn("⚠️ " .. workflow.name .. ": " .. (status and status.status or "unknown"))
                    end
                end
                
                log.info("📈 Health Report: " .. healthy_count .. "/" .. total_count .. " workflows healthy")
                
                return {
                    success = true, 
                    message = "GitOps monitoring completed",
                    metrics = {
                        total_workflows = total_count,
                        healthy_workflows = healthy_count,
                        health_percentage = (healthy_count / total_count) * 100
                    }
                }
            end
        }
    },
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 GitOps Demo workflow completed successfully!")
            log.info("🔄 GitOps está pronto para uso em produção!")
        else
            log.error("💥 GitOps Demo workflow failed")
        end
    end
})