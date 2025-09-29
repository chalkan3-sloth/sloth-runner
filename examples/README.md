# 🦥 Sloth Runner Examples - Modern DSL

Esta pasta contém uma coleção abrangente de exemplos que demonstram as capacidades do Sloth Runner, incluindo funcionalidades como distributed execution, state management, monitoring, e integração com cloud providers.

## 🚀 **Exemplos Práticos**

### 📊 **Production Ready Examples**

**🔄 [gitops_kubernetes_advanced.lua](./gitops_kubernetes_advanced.lua)** - Pipeline completo GitOps com Kubernetes:
- **🌐 Distributed Execution**: Execução distribuída através de agents
- **💾 State Management**: Gerenciamento de estado persistente
- **📊 Monitoring**: Métricas e alertas integrados
- **🔄 CI/CD Pipeline**: Pipeline completo de deploy
- **🛡️ Security**: mTLS e RBAC integrados

**🐍 [ai_powered_pipeline.lua](./ai_powered_pipeline.lua)** - Pipeline inteligente com analytics:
- **📈 Predictive Analytics**: Análise preditiva de performance
- **🎯 Adaptive Optimization**: Otimização automática de recursos
- **🔄 Self-Healing**: Auto-recovery de falhas
- **📊 Real-time Monitoring**: Monitoramento em tempo real

```bash
# Execute exemplos práticos
./sloth-runner run -f examples/gitops_kubernetes_advanced.lua
./sloth-runner run -f examples/ai_powered_pipeline.lua
```

---

## 🚀 **Modern DSL - Sintaxe Única e Moderna**

Todos os exemplos agora usam **EXCLUSIVAMENTE** a Modern DSL - uma linguagem específica de domínio que oferece:

- **🎯 Fluent API**: Sintaxe intuitiva e encadeável
- **📋 Workflow Definition**: Configuração declarativa de workflows  
- **🔄 Enhanced Features**: Retry strategies, circuit breakers, e padrões avançados
- **🛡️ Type Safety**: Melhor validação e detecção de erros
- **📊 Rich Metadata**: Informações detalhadas de tasks e workflows
- **🧹 Clean Syntax**: Sintaxe limpa sem formato legacy

### ✨ **Modern DSL - Sintaxe Única**

```lua
-- 🎯 Definição de Task com Fluent API
local build_task = task("build_application")
    :description("Build com recursos modernos")
    :command(function(params, deps)
        log.info("Construindo aplicação...")
        return exec.run("go build -o app ./cmd/main.go")
    end)
    :timeout("5m")
    :retries(3, "exponential")
    :artifacts({"app"})
    :tags({"build", "application"})
    :build()

-- 📋 Definição de Workflow Declarativo
workflow.define("ci_pipeline", {
    description = "Pipeline CI/CD completo",
    version = "2.0.0",
    metadata = {
        team = "devops",
        environment = "production"
    },
    
    tasks = {
        build_task,
        test_task,
        deploy_task
    },
    
    on_success = function(results)
        log.info("✅ Pipeline executado com sucesso!")
        notify.slack("Pipeline CI/CD concluído", results)
    end,
    
    on_failure = function(error, context)
        log.error("❌ Falha no pipeline: " .. error.message)
        notify.slack("Pipeline falhou", { error = error, context = context })
    end
})
```

---

## 📂 **Estrutura dos Exemplos**

### 🌟 **Exemplos Destacados**

| Exemplo | Descrição | Complexidade | Recursos |
|---------|-----------|--------------|----------|
| [**simple_ai_demo.lua**](./simple_ai_demo.lua) | Demo básico com analytics | ⭐ Básico | Analytics, Monitoring |
| [**gitops_native_demo.lua**](./gitops_native_demo.lua) | GitOps workflow | ⭐⭐ Intermediário | Git, State, Notifications |
| [**gitops_kubernetes_advanced.lua**](./gitops_kubernetes_advanced.lua) | K8s + GitOps avançado | ⭐⭐⭐ Avançado | K8s, GitOps, Distributed |
| [**ai_powered_pipeline.lua**](./ai_powered_pipeline.lua) | Pipeline com IA | ⭐⭐⭐ Avançado | Analytics, Prediction, Optimization |

### 🔧 **Exemplos por Categoria**

#### 📊 **Analytics & Intelligence**
- **[simple_ai_demo.lua](./simple_ai_demo.lua)**: Demonstração básica de analytics
- **[ai_powered_pipeline.lua](./ai_powered_pipeline.lua)**: Pipeline com análise preditiva
- **[test_ai_module.lua](./test_ai_module.lua)**: Teste dos módulos de IA

#### 🔄 **GitOps & CI/CD**
- **[gitops_native_demo.lua](./gitops_native_demo.lua)**: Workflow GitOps básico
- **[gitops_kubernetes_advanced.lua](./gitops_kubernetes_advanced.lua)**: GitOps avançado com K8s
- **[test_gitops_basic.lua](./test_gitops_basic.lua)**: Teste básico do GitOps

#### 🎯 **Recursos Específicos**
- **[ai_intelligence_showcase.lua](./ai_intelligence_showcase.lua)**: Showcase de inteligência
- **[iac_integration_showcase.lua](./iac_integration_showcase.lua)**: Integração IaC
- **[unified_fluent_workflow.lua](./unified_fluent_workflow.lua)**: Workflow fluente unificado

---

## 🎯 **Como Usar os Exemplos**

### 1. **📥 Preparação**
```bash
# Clone o repositório
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Compile o Sloth Runner
go build -o sloth-runner ./cmd/sloth-runner
```

### 2. **🚀 Execução Básica**
```bash
# Execute um exemplo simples
./sloth-runner run -f examples/simple_ai_demo.lua

# Execute com verbose para debug
./sloth-runner run -f examples/gitops_native_demo.lua --verbose

# Execute com parâmetros customizados
./sloth-runner run -f examples/ai_powered_pipeline.lua --param environment=staging
```

### 3. **🔧 Modificação dos Exemplos**
```bash
# Copie um exemplo como base
cp examples/simple_ai_demo.lua my_workflow.lua

# Edite conforme necessário
vim my_workflow.lua

# Execute seu workflow customizado
./sloth-runner run -f my_workflow.lua
```

---

## 📚 **Recursos de Aprendizagem**

### 🎓 **Para Iniciantes**
1. Comece com **[simple_ai_demo.lua](./simple_ai_demo.lua)** - exemplo mais simples
2. Entenda a sintaxe Modern DSL no código
3. Execute e observe os logs de saída
4. Modifique valores e re-execute para experimentar

### 🏗️ **Para Desenvolvedores**
1. Analise **[gitops_kubernetes_advanced.lua](./gitops_kubernetes_advanced.lua)** - exemplo complexo
2. Estude os padrões de error handling e retry
3. Observe como state management é usado
4. Implemente seus próprios workflows baseados nos exemplos

### 🚀 **Para DevOps**
1. Use **[ai_powered_pipeline.lua](./ai_powered_pipeline.lua)** como base para CI/CD
2. Adapte para sua infraestrutura específica
3. Configure alertas e monitoramento
4. Implemente estratégias de deployment

---

## 🤝 **Contribuindo com Exemplos**

Tem um exemplo interessante? Contribua!

1. **📝 Crie seu exemplo** seguindo a Modern DSL
2. **📋 Adicione documentação** inline no código
3. **✅ Teste completamente** antes de submeter
4. **📧 Abra um PR** com descrição detalhada

### 📐 **Padrões para Novos Exemplos**
- Use **exclusivamente Modern DSL**
- Inclua **comentários explicativos**
- Adicione **error handling apropriado**
- Demonstre **pelo menos 2-3 recursos** do Sloth Runner
- Seja **prático e realista** (evite ficção científica)

---

## 🔗 **Links Úteis**

- 📖 [Documentação Completa](../docs/)
- 🧠 [Core Concepts](../docs/en/core-concepts.md)
- ⚡ [Quick Start](../docs/en/quick-start.md)
- 🎯 [Advanced Features](../docs/en/advanced-features.md)
- 🤖 [AI Features](../docs/en/ai-features.md)

---

## 🆘 **Suporte & Ajuda**

- 🐛 [Issues](https://github.com/chalkan3-sloth/sloth-runner/issues)
- 💬 [Discussions](https://github.com/chalkan3-sloth/sloth-runner/discussions)
- 📧 Email: support@sloth-runner.dev

**Happy Automating! 🦥🚀**