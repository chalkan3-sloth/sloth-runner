# 🗂️ Stack Management

O Sloth Runner oferece um sistema completo de gerenciamento de stacks similar ao Pulumi, permitindo persistir o estado dos workflows e rastrear execuções ao longo do tempo.

## 🚀 Introdução

O **Stack Management** no Sloth Runner permite:

- **Persistir estado** entre execuções
- **Rastrear outputs** exportados da pipeline
- **Histórico completo** de execuções
- **Gestão via CLI** intuitiva
- **Isolamento** por ambiente/projeto
- **Database persistente** em `/etc/sloth-runner/`

## 📝 Sintaxe Básica

### Executar com Stack

```bash
# Nova sintaxe - stack name como argumento posicional
sloth-runner run {stack-name} --file workflow.lua

# Exemplos práticos
sloth-runner run production-app -f deploy.lua --output enhanced
sloth-runner run dev-environment -f test.lua -o json
sloth-runner run my-cicd -f pipeline.lua --output modern
```

### Gerenciar Stacks

```bash
# Listar todos os stacks
sloth-runner stack list

# Ver detalhes de um stack
sloth-runner stack show production-app

# Remover stack (com confirmação)
sloth-runner stack delete old-environment

# Remover stack (sem confirmação)
sloth-runner stack delete old-environment --force
```

### 🆔 Listar Tasks e Grupos

```bash
# Listar todos os grupos e tasks com IDs únicos
sloth-runner list -f workflow.lua

# Exemplo de saída:
# Workflow Tasks and Groups
# 
# ## Task Group: production-pipeline
# ID: 7680eaa4-7b4d-4d71-8bdf-b659f719d75c
# Description: Production deployment pipeline
# 
# Tasks:
# NAME      ID              DESCRIPTION              DEPENDS ON
# ----      --              -----------              ----------
# setup     a8cad56b...     Setup environment       -
# build     628e29e7...     Build application       setup
# test      3f07511c...     Run tests               build
# deploy    9b12345a...     Deploy to production    test
```

## 🎨 Opções de Output

### Output Styles Disponíveis

```bash
# Output básico (padrão)
sloth-runner run my-stack -f workflow.lua --output basic

# Output melhorado (estilo Pulumi)
sloth-runner run my-stack -f workflow.lua --output enhanced

# Output moderno com cores
sloth-runner run my-stack -f workflow.lua --output modern

# Output em JSON estruturado
sloth-runner run my-stack -f workflow.lua --output json
```

### Exemplo de Output JSON

```json
{
  "status": "success",
  "duration": "2.19075ms",
  "tasks": {
    "setup": {
      "status": "Success",
      "duration": "193.667µs",
      "error": ""
    },
    "build": {
      "status": "Success", 
      "duration": "420.666µs",
      "error": ""
    }
  },
  "outputs": {
    "app_version": "1.2.3",
    "deployment_url": "https://app.example.com",
    "status": "deployed"
  },
  "stack": {
    "name": "my-stack",
    "id": "42d60749-7092-4398-ae2d-2702e4a16e0a"
  },
  "workflow": "my-stack",
  "execution_time": 1759238671
}
```

## 💾 Database Storage

### Local de Armazenamento

O Sloth Runner agora armazena todos os dados de stack em:

```
/etc/sloth-runner/stacks.db
```

Esta localização permite:
- **Persistência** entre sessões
- **Compartilhamento** entre usuários
- **Backup** centralizado
- **Performance** otimizada

### Estrutura do Database

```sql
-- Tabela principal de stacks
CREATE TABLE stacks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    version TEXT,
    status TEXT NOT NULL DEFAULT 'created',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    workflow_file TEXT,
    task_results TEXT,
    outputs TEXT,
    configuration TEXT,
    metadata TEXT,
    execution_count INTEGER DEFAULT 0,
    last_duration INTEGER DEFAULT 0,
    last_error TEXT,
    resource_version TEXT NOT NULL DEFAULT '1'
);

-- Tabela de execuções
CREATE TABLE stack_executions (
    id TEXT PRIMARY KEY,
    stack_id TEXT NOT NULL,
    started_at DATETIME NOT NULL,
    completed_at DATETIME,
    duration INTEGER,
    status TEXT NOT NULL,
    task_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    failure_count INTEGER DEFAULT 0,
    outputs TEXT,
    error_message TEXT,
    FOREIGN KEY (stack_id) REFERENCES stacks(id)
);
```

## 🔍 Comandos de Stack

### `sloth-runner stack list`

Lista todos os stacks disponíveis:

```bash
$ sloth-runner stack list

Workflow Stacks
NAME              STATUS      LAST RUN           DURATION     EXECUTIONS   DESCRIPTION
----              ------      --------           --------     ----------   -----------
production-app    completed   2024-01-15 14:30   2.5s         15           Production deployment
dev-environment   completed   2024-01-15 10:15   1.2s         8            Development environment
staging-tests     failed      2024-01-14 16:45   45ms         3            Staging test suite
```

### `sloth-runner stack show <name>`

Mostra detalhes completos de um stack:

```bash
$ sloth-runner stack show production-app

Stack: production-app

ID: 6fcf2daf-ac69-40d0-95ad-ce39d5bd24b8
Description: Production deployment pipeline
Version: 1.0.0
Status: completed
Created: 2024-01-15 08:00:00
Updated: 2024-01-15 14:30:15
Completed: 2024-01-15 14:30:15
Workflow File: production-deploy.lua
Executions: 15
Last Duration: 2.5s

Outputs
app_version: 1.2.3
deployment_url: https://app.example.com
database_version: 5.7.2
status: deployed

Recent Executions
STARTED            STATUS      DURATION   TASKS   SUCCESS   FAILED
-------            ------      --------   -----   -------   ------
2024-01-15 14:30   completed   2.5s       5       5         0
2024-01-15 12:15   completed   2.1s       5       5         0
2024-01-15 10:30   failed      1.8s       5       3         2
```

### `sloth-runner stack delete <name>`

Remove um stack:

```bash
# Com confirmação interativa
$ sloth-runner stack delete old-environment
⚠ This will permanently delete stack 'old-environment' and all its execution history.
Are you sure? (y/N): y
✓ Stack 'old-environment' deleted successfully.

# Sem confirmação (modo força)
$ sloth-runner stack delete old-environment --force
✓ Stack 'old-environment' deleted successfully.
```

## 🎯 Casos de Uso Avançados

### Pipeline CI/CD

```lua
-- cicd-pipeline.lua
TaskDefinitions = {
    cicd_pipeline = {
        description = "Complete CI/CD pipeline",
        tasks = {
            {
                name = "checkout",
                description = "Checkout source code",
                command = "git clone https://github.com/user/repo.git ."
            },
            {
                name = "build",
                description = "Build application",
                depends_on = {"checkout"},
                command = function()
                    print("Building application...")
                    -- Build logic here
                    outputs = outputs or {}
                    outputs.build_version = "1.2.3"
                    outputs.build_hash = "abc123def"
                    return true
                end
            },
            {
                name = "test",
                description = "Run tests",
                depends_on = {"build"},
                command = "npm test"
            },
            {
                name = "deploy",
                description = "Deploy to production",
                depends_on = {"test"},
                command = function()
                    print("Deploying to production...")
                    outputs = outputs or {}
                    outputs.deployment_url = "https://app.example.com"
                    outputs.deployment_time = os.time()
                    return true
                end
            }
        }
    }
}
```

```bash
# Executar pipeline
sloth-runner run cicd-main -f cicd-pipeline.lua --output enhanced

# Verificar status
sloth-runner stack show cicd-main

# Ver outputs em JSON
sloth-runner run cicd-main -f cicd-pipeline.lua --output json
```

### Múltiplos Ambientes

```bash
# Diferentes stacks por ambiente
sloth-runner run production -f deploy.lua --output enhanced
sloth-runner run staging -f deploy.lua --output enhanced
sloth-runner run development -f deploy.lua --output enhanced

# Listar todos os ambientes
sloth-runner stack list
```

## 📊 Monitoramento e Analytics

### Histórico de Execuções

Cada stack mantém um histórico completo:

- **Timestamps** de início e fim
- **Duração** de cada execução
- **Status** (success/failed)
- **Contadores** de tasks (total/success/failed)
- **Outputs** exportados
- **Mensagens de erro**

### Outputs Exportados

Os workflows podem exportar dados que são:

- **Persistidos** no stack
- **Versionados** por execução
- **Acessíveis** via CLI
- **Exportáveis** em JSON

```lua
-- Exportar outputs no workflow
TaskDefinitions = {
    export_example = {
        description = "Example with exports",
        tasks = {
            {
                name = "process",
                command = function()
                    -- Criar/atualizar outputs globais
                    outputs = outputs or {}
                    outputs.processed_items = 42
                    outputs.timestamp = os.time()
                    outputs.status = "success"
                    outputs.version = "1.0.0"
                    return true
                end
            }
        }
    }
}
```

## 🔧 Configuração Avançada

### Configuração de Database

Por padrão, o database é criado em `/etc/sloth-runner/stacks.db`, mas pode ser personalizado:

```bash
# Usando variável de ambiente
export SLOTH_RUNNER_DB_PATH="/custom/path/stacks.db"
sloth-runner run my-stack -f workflow.lua
```

### Backup e Restore

```bash
# Backup do database
cp /etc/sloth-runner/stacks.db /backup/location/

# Restore do database
cp /backup/location/stacks.db /etc/sloth-runner/
```

## 🚀 Próximos Passos

1. **Experimente** diferentes outputs: `basic`, `enhanced`, `modern`, `json`
2. **Crie** múltiplos stacks para diferentes projetos
3. **Monitore** execuções com `stack list` e `stack show`
4. **Use** IDs únicos para rastreamento granular
5. **Implemente** pipelines de CI/CD com persistência de estado 
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