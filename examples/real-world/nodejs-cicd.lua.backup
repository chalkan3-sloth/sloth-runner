-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03

require('dotenv').config();

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
    nodejs_cicd_pipeline = {
        description = "Pipeline completo de CI/CD para aplicação Node.js",
        
        tasks = {
            {
                name = "setup_pipeline_environment",
                description = "Configura ambiente do pipeline e variáveis",
                command = function(params)
                    log.info("🚀 Configurando Pipeline CI/CD para Node.js...")
                    
                    -- Configurações do pipeline
                    local pipeline_config = {
                        app_name = params.app_name or "nodejs-demo-app",
                        version = params.version or "1.0.0",
                        git_repo = params.git_repo or "https://github.com/example/nodejs-app.git",
                        environments = {
                            development = {
                                domain = "dev.myapp.com",
                                replicas = 1,
                                resources = {cpu = "100m", memory = "128Mi"}
                            },
                            staging = {
                                domain = "staging.myapp.com", 
                                replicas = 2,
                                resources = {cpu = "200m", memory = "256Mi"}
                            },
                            production = {
                                domain = "myapp.com",
                                replicas = 3,
                                resources = {cpu = "500m", memory = "512Mi"}
                            }
                        },
                        docker_registry = "docker.io/mycompany",
                        notification_webhook = "https://hooks.slack.com/services/..."
                    }
                    
                    -- Configurar ferramentas necessárias
                    local tools = {
                        node_version = "18",
                        npm_version = "latest",
                        docker_required = true,
                        kubectl_required = true,
                        helm_required = true
                    }
                    
                    state.set("pipeline_config", pipeline_config)
                    state.set("required_tools", tools)
                    
                    -- Inicializar métricas do pipeline
                    local metrics = {
                        pipeline_start_time = os.time(),
                        stages_completed = 0,
                        total_stages = 8,
                        tests_run = 0,
                        tests_passed = 0,
                        build_duration = 0,
                        deploy_duration = 0
                    }
                    
                    state.set("pipeline_metrics", metrics)
                    
                    log.info("✅ Pipeline configurado:")
                    log.info("  📱 App: " .. pipeline_config.app_name .. " v" .. pipeline_config.version)
                    log.info("  🌍 Ambientes: " .. table.concat({"development", "staging", "production"}, ", "))
                    log.info("  📊 Total de etapas: " .. metrics.total_stages)
                    
                    return true, "Environment configured"
                end
            },
            
            {
                name = "validate_tools_and_dependencies",
                description = "Valida ferramentas e dependências necessárias",
                depends_on = "setup_pipeline_environment",
                command = function()
                    log.info("🔧 Validando ferramentas e dependências...")
                    
                    local tools = state.get("required_tools")
                    local validation_results = {}
                    
                    -- Verificar Node.js
                    log.info("📦 Verificando Node.js...")
                    local node_check = exec.run("node --version")
                    if node_check.success then
                        local node_version = string.gsub(node_check.stdout, "v", "")
                        log.info("  ✅ Node.js " .. string.gsub(node_version, "\n", "") .. " encontrado")
                        validation_results.node = true
                    else
                        log.error("  ❌ Node.js não encontrado")
                        validation_results.node = false
                    end
                    
                    -- Verificar npm
                    log.info("📦 Verificando npm...")
                    local npm_check = exec.run("npm --version")
                    if npm_check.success then
                        log.info("  ✅ npm " .. string.gsub(npm_check.stdout, "\n", "") .. " encontrado")
                        validation_results.npm = true
                    else
                        log.error("  ❌ npm não encontrado")
                        validation_results.npm = false
                    end
                    
                    -- Verificar Docker
                    log.info("🐳 Verificando Docker...")
                    local docker_check = exec.run("docker --version")
                    if docker_check.success then
                        log.info("  ✅ Docker encontrado")
                        validation_results.docker = true
                        
                        -- Verificar se Docker daemon está rodando
                        local docker_info = exec.run("docker info")
                        if docker_info.success then
                            log.info("  ✅ Docker daemon está rodando")
                        else
                            log.error("  ❌ Docker daemon não está rodando")
                            validation_results.docker = false
                        end
                    else
                        log.error("  ❌ Docker não encontrado")
                        validation_results.docker = false
                    end
                    
                    -- Verificar Git
                    log.info("🔄 Verificando Git...")
                    local git_check = exec.run("git --version")
                    if git_check.success then
                        log.info("  ✅ Git encontrado")
                        validation_results.git = true
                    else
                        log.error("  ❌ Git não encontrado")
                        validation_results.git = false
                    end
                    
                    -- Verificar conectividade de rede
                    log.info("🌐 Verificando conectividade de rede...")
                    local http = require("http")
                    local connectivity_check = http.get({
                        url = "https://registry.npmjs.org/-/ping",
                        timeout = 10
                    })
                    
                    if connectivity_check.success then
                        log.info("  ✅ Conectividade com npm registry OK")
                        validation_results.network = true
                    else
                        log.error("  ❌ Problemas de conectividade: " .. connectivity_check.error)
                        validation_results.network = false
                    end
                    
                    state.set("tool_validation", validation_results)
                    
                    -- Verificar se todas as ferramentas essenciais estão disponíveis
                    local essential_tools = {"node", "npm", "docker", "git", "network"}
                    local missing_tools = {}
                    
                    for _, tool in ipairs(essential_tools) do
                        if not validation_results[tool] then
                            table.insert(missing_tools, tool)
                        end
                    end
                    
                    if #missing_tools > 0 then
                        log.error("❌ Ferramentas essenciais faltando: " .. table.concat(missing_tools, ", "))
                        return false, "Missing essential tools"
                    end
                    
                    log.info("✅ Todas as ferramentas necessárias estão disponíveis!")
                    
                    return true, "Tools validated"
                end
            },
            
            {
                name = "checkout_and_prepare_code",
                description = "Faz checkout do código e prepara ambiente",
                depends_on = "validate_tools_and_dependencies",
                command = function()
                    log.info("📥 Fazendo checkout do código...")
                    
                    local config = state.get("pipeline_config")
                    local app_dir = "./" .. config.app_name
                    
                    -- Simular checkout (em produção real, usaria git clone)
                    if not fs.exists(app_dir) then
                        fs.mkdir(app_dir)
                    end
                    
                    -- Criar estrutura básica de uma app Node.js
                    log.info("🏗️  Criando estrutura da aplicação...")
                    
                    -- package.json
                    local package_json = {
                        name = config.app_name,
                        version = config.version,
                        description = "Demo Node.js application for CI/CD pipeline",
                        main = "index.js",
                        scripts = {
                            start = "node index.js",
                            test = "jest",
                            ["test:unit"] = "jest --testPathPattern=unit",
                            ["test:integration"] = "jest --testPathPattern=integration",
                            ["test:coverage"] = "jest --coverage",
                            lint = "eslint .",
                            ["lint:fix"] = "eslint . --fix",
                            build = "npm run lint && npm run test"
                        },
                        dependencies = {
                            express = "^4.18.0",
                            cors = "^2.8.5",
                            helmet = "^6.0.0",
                            dotenv = "^16.0.0"
                        },
                        devDependencies = {
                            jest = "^29.0.0",
                            supertest = "^6.2.0",
                            eslint = "^8.0.0",
                            nodemon = "^2.0.0"
                        },
                        engines = {
                            node = ">=18.0.0",
                            npm = ">=8.0.0"
                        }
                    }
                    
                    fs.write(app_dir .. "/package.json", data.to_json(package_json))
                    
                    -- index.js
                    local index_js = [[
const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
require('dotenv').config();

const app = express();
const port = process.env.PORT || 3000;

app.use(helmet());
app.use(cors());
app.use(express.json());

app.get('/', (req, res) => {
  res.json({
    message: 'Hello from Node.js CI/CD Demo!',
    version: process.env.npm_package_version || '1.0.0',
    environment: process.env.NODE_ENV || 'development',
    timestamp: new Date().toISOString()
  });
});
