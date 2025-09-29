# 🗂️ Stack Management

O Sloth Runner oferece um sistema completo de gerenciamento de stacks similar ao Pulumi, permitindo persistir o estado dos workflows e rastrear execuções ao longo do tempo.

## 🚀 Introdução

O **Stack Management** no Sloth Runner permite:

- **Persistir estado** entre execuções
- **Rastrear outputs** exportados da pipeline
- **Histórico completo** de execuções
- **Gestão via CLI** intuitiva
- **Isolamento** por ambiente/projeto

## 📝 Sintaxe Básica

### Executar com Stack

```bash
# Nova sintaxe - stack name como argumento posicional
sloth-runner run {stack-name} --file workflow.lua

# Exemplos práticos
sloth-runner run production-app -f deploy.lua --output enhanced
sloth-runner run dev-environment -f test.lua -o rich
sloth-runner run my-cicd -f pipeline.lua
```

### Gerenciar Stacks

```bash
# Listar todos os stacks
sloth-runner stack list

# Ver detalhes de um stack
sloth-runner stack show production-app

# Remover stack
sloth-runner stack delete old-environment
```

### 🆔 Listar Tasks e Grupos (Novo)

```bash
# Listar todos os grupos e tasks com IDs únicos
sloth-runner list -f workflow.lua

# Exemplo de saída:
# Workflow Tasks and Groups
# 
# ## Task Group: deploy_group
# ID: e8e77f72-5cf4-4e98-adce-fc839846c24a
# Description: Deployment tasks with IDs
#
# Tasks:
# NAME     ID           DESCRIPTION             DEPENDS ON
# ----     --           -----------             ----------
# build    a1c4fa46...  Build the application   -
# test     d8dc4623...  Run tests               build
# deploy   6253cb19...  Deploy to production    build, test
```

## 🎯 Conceitos Fundamentais

### Stack State

Cada stack mantém:

- **ID único** (UUID)
- **Nome** do stack
- **Status** atual (created, running, completed, failed)
- **Outputs exportados** da pipeline
- **Histórico** de execuções
- **Metadados** e configurações

### 🆔 IDs Únicos de Tasks e Grupos

**Novidade:** Cada task e task group agora possui **IDs únicos (UUID)** para rastreabilidade completa:

#### Task IDs
- **UUID único** gerado automaticamente para cada task
- **Persistente** entre execuções
- **Rastreável** durante debugging e logs
- **Visível** no comando `sloth-runner list`

#### Task Group IDs  
- **UUID único** para cada grupo de tasks
- **Identificação** clara de componentes do workflow
- **Organização** hierárquica com IDs
- **Debugging** facilitado com identificação precisa

#### Benefícios dos IDs
- 🔍 **Debugging melhorado** com identificação precisa
- 📊 **Observabilidade** enhanced para Enterprise
- 🎯 **Execução seletiva** por ID (futuro)
- 📈 **Rastreabilidade** completa de execuções

### Ciclo de Vida

1. **Criação**: Stack é criado automaticamente na primeira execução
2. **Execução**: Estado é atualizado durante a execução
3. **Persistência**: Outputs e resultados são salvos
4. **Reutilização**: Execuções subsequentes reutilizam o stack

## 💾 Persistência de Estado

### Banco de Dados

O Sloth Runner usa **SQLite** para persistir o estado:

```
~/.sloth-runner/stacks.db
```

### Tabelas

- **stacks**: Informações principais dos stacks
- **stack_executions**: Histórico detalhado de execuções

## 📊 Outputs Exportados

### Captura Automática

O sistema captura automaticamente:

- **Exports do TaskRunner** (`runner.Exports`)
- **Variável global `outputs`** do Lua
- **Metadados** de execução

### Exemplo de Export

```lua
local deploy_task = task("deploy")
    :command(function(params, deps)
        -- Lógica de deploy...
        
        -- Exportar outputs para o stack
        runner.Export({
            app_url = "https://myapp.example.com",
            version = "1.2.3",
            environment = "production",
            deployed_at = os.date()
        })
        
        return true, "Deploy successful", deploy_info
    end)
    :build()
```

## 🖥️ Interface CLI

### Lista de Stacks

```bash
$ sloth-runner stack list

Workflow Stacks     

NAME                  STATUS     LAST RUN           DURATION     EXECUTIONS
----                  ------     --------           --------     ----------
production-app        completed  2025-09-29 19:27   6.8s         5
dev-environment       running    2025-09-29 19:25   2.1s         12
staging-api           failed     2025-09-29 19:20   4.2s         3
```

### Detalhes do Stack

```bash
$ sloth-runner stack show production-app

Stack: production-app     

ID: abc123-def456-789
Status: completed
Created: 2025-09-29 15:30:21
Updated: 2025-09-29 19:27:15
Executions: 5
Last Duration: 6.8s

     Outputs     

app_url: "https://myapp.example.com"
version: "1.2.3"
environment: "production"
deployed_at: "2025-09-29 19:27:15"

     Recent Executions     

STARTED            STATUS     DURATION   TASKS   SUCCESS   FAILED
-------            ------     --------   -----   -------   ------
2025-09-29 19:27   completed  6.8s       3       3         0
2025-09-29 18:45   completed  7.2s       3       3         0
2025-09-29 17:30   failed     4.1s       3       2         1
```

## 🎨 Output Styles

### Configurável por Execução

```bash
# Output básico (padrão)
sloth-runner run my-stack -f workflow.lua

# Output melhorado
sloth-runner run my-stack -f workflow.lua --output enhanced
sloth-runner run my-stack -f workflow.lua -o rich
sloth-runner run my-stack -f workflow.lua --output modern
```

### Estilo Pulumi

O output `enhanced` oferece formatação rica similar ao Pulumi:

```
🦥 Sloth Runner

     Workflow: production-app     

Started at: 2025-09-29 19:27:15

✓ build (2.1s) completed
✓ test (3.2s) completed  
✓ deploy (1.5s) completed

     Workflow Completed Successfully     

✓ production-app
Duration: 6.8s
Tasks executed: 3

     Outputs     

├─ exports:
  │ app_url: "https://myapp.example.com"
  │ version: "1.2.3"
  │ environment: "production"
```

## 🔧 Casos de Uso

### Ambientes Separados

```bash
# Desenvolvimento
sloth-runner run dev-app -f app.lua

# Staging  
sloth-runner run staging-app -f app.lua

# Produção
sloth-runner run prod-app -f app.lua --output enhanced
```

### CI/CD Integration

```bash
# No pipeline CI/CD
sloth-runner run ${ENVIRONMENT}-${APP_NAME} -f pipeline.lua

# Exemplos:
sloth-runner run prod-frontend -f frontend-deploy.lua
sloth-runner run staging-api -f api-deploy.lua
```

### Monitoramento

```bash
# Ver status de todos os ambientes
sloth-runner stack list

# Verificar último deploy de produção
sloth-runner stack show prod-app

# Limpar ambientes de teste
sloth-runner stack delete temp-test-env
```

## 🛠️ Melhores Práticas

### Nomeação de Stacks

```bash
# Usar padrão: {ambiente}-{aplicação}
sloth-runner run prod-frontend -f deploy.lua
sloth-runner run staging-api -f deploy.lua
sloth-runner run dev-database -f setup.lua
```

### Export de Outputs

```lua
-- Exportar informações relevantes
runner.Export({
    -- URLs importantes
    app_url = deploy_info.url,
    admin_url = deploy_info.admin_url,
    
    -- Informações de versão
    version = build_info.version,
    commit_hash = build_info.commit,
    
    -- Configurações de ambiente
    environment = config.environment,
    region = config.region,
    
    -- Timestamps
    deployed_at = os.date(),
    build_time = build_info.timestamp
})
```

### Gestão de Ciclo de Vida

```bash
# Desenvolvimento ativo
sloth-runner run dev-app -f app.lua

# Quando pronto para staging
sloth-runner run staging-app -f app.lua

# Deploy para produção
sloth-runner run prod-app -f app.lua --output enhanced

# Limpeza de ambientes antigos
sloth-runner stack delete old-test-branch
```

## 🔄 Migração

### Comandos Antigos vs Novos

```bash
# Antes
sloth-runner run -f workflow.lua --stack my-stack

# Agora
sloth-runner run my-stack -f workflow.lua
```

### Compatibilidade

- Workflows existentes continuam funcionando
- Stack é opcional - pode executar sem especificar
- Outputs são capturados automaticamente quando stack é usado

## 📚 Próximos Passos

- [Output Styles](output-styles.md) - Configuração de estilos de output
- [Workflow Scaffolding](workflow-scaffolding.md) - Criação de projetos
- [Examples](../examples/) - Exemplos práticos
- [CLI Reference](CLI.md) - Referência completa de comandos