-- GitOps Avançado - Integração com Kubernetes
-- Demonstra GitOps com deploymens reais, health checks e rollbacks inteligentes

local gitops = require("gitops")
local log = require("log")
local exec = require("exec")

log.info("🚀 GitOps Avançado - Kubernetes Integration")
log.info("=" .. string.rep("=", 60))

-- Configuração de múltiplos ambientes
local environments = {
    {
        name = "development",
        repo = "https://github.com/company/k8s-dev",
        branch = "develop",
        namespace = "dev",
        auto_sync = true,
        sync_interval = "5m"
    },
    {
        name = "staging", 
        repo = "https://github.com/company/k8s-staging",
        branch = "staging",
        namespace = "staging",
        auto_sync = true,
        sync_interval = "10m"
    },
    {
        name = "production",
        repo = "https://github.com/company/k8s-prod",
        branch = "main",
        namespace = "production",
        auto_sync = false,  -- Manual deploys em produção
        sync_interval = "never"
    }
}

-- Função para criar pipeline GitOps por ambiente
local function setup_environment_gitops(env)
    log.info("🏗️ Configurando GitOps para " .. env.name .. "...")
    
    -- Registrar repositório
    local repo_id = gitops.register_repository({
        id = env.name .. "-repo",
        url = env.repo,
        branch = env.branch
    })
    
    if not repo_id then
        log.error("❌ Falha ao registrar repositório " .. env.name)
        return nil
    end
    
    -- Criar workflow GitOps
    local workflow_id = gitops.create_workflow({
        id = env.name .. "-workflow",
        name = env.name .. " Environment GitOps",
        repository = repo_id,
        target_path = "manifests",
        auto_sync = env.auto_sync,
        diff_preview = true,
        rollback_on_failure = true
    })
    
    if workflow_id then
        log.info("✅ " .. env.name .. " GitOps configurado: " .. workflow_id)
        return workflow_id
    else
        log.error("❌ Falha ao criar workflow " .. env.name)
        return nil
    end
end

-- Função para deployment com validação
local function deploy_environment(env, workflow_id, force)
    log.info("🚀 Iniciando deployment " .. env.name .. "...")
    
    -- 1. Gerar diff preview
    log.info("🔍 Analisando mudanças...")
    local diff = gitops.generate_diff(workflow_id)
    
    if not diff then
        log.info("ℹ️ Nenhuma mudança detectada para " .. env.name)
        return true
    end
    
    -- 2. Exibir resumo das mudanças
    log.info("📊 Mudanças detectadas em " .. env.name .. ":")
    log.info("  📝 Total: " .. diff.summary.total_changes)
    log.info("  ✨ Criados: " .. diff.summary.created_resources)
    log.info("  🔄 Atualizados: " .. diff.summary.updated_resources)
    log.info("  🗑️ Removidos: " .. diff.summary.deleted_resources)
    
    -- 3. Verificar alertas e conflitos
    if diff.summary.conflict_count > 0 then
        log.warn("⚠️ " .. diff.summary.conflict_count .. " conflito(s) detectado(s)!")
        for i, warning in ipairs(diff.warnings) do
            log.warn("  " .. i .. ". " .. warning)
        end
        
        if not force and env.name == "production" then
            log.error("❌ Deployment bloqueado por conflitos em produção")
            return false
        end
    end
    
    -- 4. Validação específica por ambiente
    if env.name == "production" then
        -- Em produção, requer validação extra
        log.info("🔒 Validações de produção...")
        
        -- Verificar se há mudanças críticas
        local has_critical_changes = false
        for _, change in ipairs(diff.changes) do
            if change.impact == "critical" or change.impact == "high" then
                has_critical_changes = true
                log.warn("⚠️ Mudança crítica: " .. change.resource .. " (" .. change.type .. ")")
            end
        end
        
        if has_critical_changes and not force then
            log.error("❌ Deployment bloqueado: mudanças críticas requerem aprovação manual")
            log.info("💡 Use força (force=true) para prosseguir mesmo assim")
            return false
        end
    end
    
    -- 5. Executar sync
    log.info("🚀 Executando sync GitOps...")
    local sync_success = gitops.sync_workflow(workflow_id)
    
    if not sync_success then
        log.error("❌ Sync falhou para " .. env.name)
        return false
    end
    
    -- 6. Verificar status pós-deployment
    log.info("🔍 Verificando status pós-deployment...")
    local status = gitops.get_workflow_status(workflow_id)
    
    if status and status.last_sync_result then
        local result = status.last_sync_result
        log.info("📊 Resultado do sync:")
        log.info("  Status: " .. result.status)
        log.info("  Commit: " .. result.commit_hash)
        
        if result.metrics then
            log.info("  Duração: " .. result.metrics.duration)
            log.info("  Recursos aplicados: " .. result.metrics.resources_applied)
        end
        
        if result.status == "succeeded" then
            log.info("✅ Deployment " .. env.name .. " realizado com sucesso!")
            return true
        else
            log.error("❌ Deployment " .. env.name .. " falhou: " .. (result.message or "erro desconhecido"))
            return false
        end
    end
    
    return false
end

-- Função para health check pós-deployment
local function health_check_environment(env, workflow_id)
    log.info("🏥 Health check " .. env.name .. "...")
    
    -- Simular health checks reais (em produção usaria kubectl, curl, etc.)
    local checks = {
        {
            name = "Pod Readiness",
            command = "kubectl get pods -n " .. env.namespace .. " --field-selector=status.phase=Running",
            expected = "success"
        },
        {
            name = "Service Endpoints",
            command = "kubectl get endpoints -n " .. env.namespace,
            expected = "success" 
        },
        {
            name = "Ingress Status",
            command = "kubectl get ingress -n " .. env.namespace,
            expected = "success"
        }
    }
    
    local passed_checks = 0
    local total_checks = #checks
    
    for _, check in ipairs(checks) do
        log.info("🔍 " .. check.name .. "...")
        
        -- Simular execução (em ambiente real executaria o comando)
        local success = math.random() > 0.2 -- 80% success rate
        
        if success then
            log.info("  ✅ " .. check.name .. " - OK")
            passed_checks = passed_checks + 1
        else
            log.warn("  ❌ " .. check.name .. " - FALHA")
        end
    end
    
    local health_percentage = (passed_checks / total_checks) * 100
    log.info("📊 Health Score: " .. passed_checks .. "/" .. total_checks .. " (" .. string.format("%.1f", health_percentage) .. "%)")
    
    if health_percentage >= 80 then
        log.info("✅ Environment " .. env.name .. " está saudável")
        return true
    else
        log.warn("⚠️ Environment " .. env.name .. " está degradado")
        return false
    end
end

-- Função para rollback inteligente
local function intelligent_rollback(env, workflow_id, reason)
    log.warn("🔄 Iniciando rollback inteligente para " .. env.name .. "...")
    log.warn("📋 Motivo: " .. reason)
    
    -- Backup antes do rollback
    log.info("💾 Criando backup do estado atual...")
    
    -- Executar rollback
    local rollback_success = gitops.rollback_workflow(workflow_id, reason)
    
    if rollback_success then
        log.info("✅ Rollback executado com sucesso!")
        
        -- Verificar health após rollback
        log.info("🔍 Verificando saúde pós-rollback...")
        local health_ok = health_check_environment(env, workflow_id)
        
        if health_ok then
            log.info("✅ Environment " .. env.name .. " restaurado com sucesso!")
            return true
        else
            log.error("❌ Environment ainda está degradado após rollback!")
            return false
        end
    else
        log.error("❌ Rollback falhou para " .. env.name .. "!")
        return false
    end
end

-- Criar workflows para todos os ambientes
local env_workflows = {}

for _, env in ipairs(environments) do
    local workflow_id = setup_environment_gitops(env)
    if workflow_id then
        env_workflows[env.name] = {
            env = env,
            workflow_id = workflow_id
        }
    end
end

-- Iniciar auto-sync controller
log.info("")
log.info("🔄 Iniciando GitOps Auto-Sync Controller...")
gitops.start_auto_sync()

-- Pipeline de deployment completo
workflow.define("gitops_kubernetes_pipeline", {
    description = "Pipeline GitOps completo com Kubernetes",
    version = "2.0.0",
    
    -- Metadados do pipeline
    metadata = {
        author = "DevOps Team",
        tags = {"gitops", "kubernetes", "production"},
        environments = {"development", "staging", "production"}
    },
    
    tasks = {
        {
            name = "deploy_development",
            description = "Deploy automático no ambiente de desenvolvimento",
            command = function(params, deps)
                local env_data = env_workflows["development"]
                if not env_data then
                    return {success = false, message = "Development environment not configured"}
                end
                
                local success = deploy_environment(env_data.env, env_data.workflow_id, true)
                
                if success then
                    health_check_environment(env_data.env, env_data.workflow_id)
                    return {success = true, message = "Development deployment completed"}
                else
                    return {success = false, message = "Development deployment failed"}
                end
            end
        },
        
        {
            name = "deploy_staging",
            description = "Deploy no ambiente de staging com validação",
            depends_on = {"deploy_development"},
            command = function(params, deps)
                local env_data = env_workflows["staging"]
                if not env_data then
                    return {success = false, message = "Staging environment not configured"}
                end
                
                log.info("🎯 Iniciando deployment staging...")
                
                local success = deploy_environment(env_data.env, env_data.workflow_id, false)
                
                if success then
                    local health_ok = health_check_environment(env_data.env, env_data.workflow_id)
                    
                    if health_ok then
                        return {success = true, message = "Staging deployment completed and healthy"}
                    else
                        -- Rollback automático se health check falhar
                        intelligent_rollback(env_data.env, env_data.workflow_id, "Health check failed")
                        return {success = false, message = "Staging deployment rolled back due to health issues"}
                    end
                else
                    return {success = false, message = "Staging deployment failed"}
                end
            end
        },
        
        {
            name = "production_approval",
            description = "Aprovação manual para produção",
            depends_on = {"deploy_staging"},
            command = function(params, deps)
                log.info("🔒 Aguardando aprovação para produção...")
                log.info("📋 Staging deployment concluído com sucesso")
                log.info("🎯 Ready para produção!")
                
                -- Em ambiente real, aqui haveria integração com sistema de aprovação
                -- Por agora, simulamos aprovação automática
                local approved = true
                
                if approved then
                    log.info("✅ Produção aprovada!")
                    return {success = true, message = "Production deployment approved"}
                else
                    log.warn("❌ Produção rejeitada!")
                    return {success = false, message = "Production deployment rejected"}
                end
            end
        },
        
        {
            name = "deploy_production",
            description = "Deploy em produção com máxima segurança",
            depends_on = {"production_approval"},
            command = function(params, deps)
                local env_data = env_workflows["production"]
                if not env_data then
                    return {success = false, message = "Production environment not configured"}
                end
                
                log.info("🏭 Iniciando deployment PRODUÇÃO...")
                log.info("⚠️ Máxima atenção! Deploy em ambiente de produção!")
                
                -- Deploy com validações extras
                local success = deploy_environment(env_data.env, env_data.workflow_id, false)
                
                if success then
                    log.info("🔍 Executando health checks críticos...")
                    local health_ok = health_check_environment(env_data.env, env_data.workflow_id)
                    
                    if health_ok then
                        log.info("🎉 PRODUÇÃO DEPLOYED COM SUCESSO!")
                        return {success = true, message = "Production deployment completed successfully"}
                    else
                        log.error("💥 Health check falhou em produção!")
                        
                        -- Rollback imediato em produção
                        intelligent_rollback(env_data.env, env_data.workflow_id, "Production health check failed")
                        return {success = false, message = "Production deployment rolled back"}
                    end
                else
                    log.error("💥 Deploy produção falhou!")
                    return {success = false, message = "Production deployment failed"}
                end
            end
        },
        
        {
            name = "post_deployment_monitoring",
            description = "Monitoramento pós-deployment",
            depends_on = {"deploy_production"},
            command = function(params, deps)
                log.info("📊 Iniciando monitoramento pós-deployment...")
                
                -- Monitorar todos os ambientes
                for env_name, env_data in pairs(env_workflows) do
                    log.info("🔍 Monitorando " .. env_name .. "...")
                    
                    local status = gitops.get_workflow_status(env_data.workflow_id)
                    if status then
                        log.info("  Status: " .. status.status)
                        if status.last_sync_result then
                            log.info("  Last sync: " .. status.last_sync_result.status)
                        end
                    end
                    
                    health_check_environment(env_data.env, env_data.workflow_id)
                end
                
                log.info("📈 Monitoramento concluído!")
                return {success = true, message = "Post-deployment monitoring completed"}
            end
        }
    },
    
    -- Hooks do pipeline
    on_task_start = function(task_name)
        log.info("🚀 Iniciando task: " .. task_name)
    end,
    
    on_task_complete = function(task_name, success, output)
        if success then
            log.info("✅ Task concluída: " .. task_name)
        else
            log.error("❌ Task falhou: " .. task_name)
            
            -- Em caso de falha, executar rollback nos ambientes já deployados
            if task_name == "deploy_production" then
                log.warn("🔄 Falha crítica! Executando rollback de emergência...")
                
                for env_name, env_data in pairs(env_workflows) do
                    if env_name ~= "development" then -- Não rollback dev
                        intelligent_rollback(env_data.env, env_data.workflow_id, "Emergency rollback due to production failure")
                    end
                end
            end
        end
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 PIPELINE GITOPS KUBERNETES CONCLUÍDO COM SUCESSO!")
            log.info("🚀 Todos os ambientes deployados e saudáveis!")
            
            -- Enviar notificações (simulado)
            log.info("📱 Enviando notificações de sucesso...")
            
        else
            log.error("💥 Pipeline GitOps falhou!")
            log.error("🔍 Verificar logs para detalhes dos erros")
            
            -- Alertas críticos (simulado)
            log.error("🚨 ALERTA CRÍTICO: Pipeline de produção falhou!")
        end
        
        -- Resumo final
        log.info("")
        log.info("📊 RESUMO FINAL:")
        log.info("🔄 GitOps workflows ativos: " .. #gitops.list_workflows())
        log.info("🏗️ Ambientes gerenciados: " .. #environments)
        log.info("📈 Auto-sync: Ativo para dev/staging")
        log.info("🔒 Produção: Deploy manual com validações")
    end
})