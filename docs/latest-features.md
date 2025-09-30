# 🆕 Latest Features & Improvements

Esta página documenta as mais recentes funcionalidades implementadas no Sloth Runner.

## 🗂️ Stack Management Completo

### 📍 Database Location
O Sloth Runner agora armazena todos os dados em:
```
/etc/sloth-runner/stacks.db
```

Benefícios:
- ✅ **Persistência** entre sessões
- ✅ **Compartilhamento** entre usuários
- ✅ **Backup** centralizado
- ✅ **Performance** otimizada

### 🚀 Comando `run` Melhorado

**Nova sintaxe posicional para stack names:**

```bash
# Sintaxe nova (recomendada)
sloth-runner run {stack-name} --file workflow.lua --output enhanced

# Exemplos práticos
sloth-runner run production-api -f api-deploy.lua --output json
sloth-runner run staging-tests -f test-suite.lua --output modern
sloth-runner run dev-environment -f dev-setup.lua --output enhanced
```

### 🎨 Múltiplas Opções de Output

#### 1. **Basic Output** (padrão)
```bash
sloth-runner run my-stack -f workflow.lua --output basic
```

#### 2. **Enhanced Output** (estilo Pulumi)
```bash
sloth-runner run my-stack -f workflow.lua --output enhanced
```
- Progress bars animados
- Status colorido
- Resumo detalhado
- Outputs exportados

#### 3. **Modern Output** (visual aprimorado)
```bash
sloth-runner run my-stack -f workflow.lua --output modern
```

#### 4. **JSON Output** (para CI/CD)
```bash
sloth-runner run my-stack -f workflow.lua --output json
```

**Exemplo de output JSON:**
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
    "deploy": {
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

## 🆔 IDs Únicos para Tasks e Grupos

### ✨ Funcionalidade
Cada task e grupo agora possui um **UUID único** para rastreamento granular.

### 📋 Comando `list`
```bash
sloth-runner list -f workflow.lua
```

**Exemplo de output:**
```
Workflow Tasks and Groups

## Task Group: production-pipeline
ID: 7680eaa4-7b4d-4d71-8bdf-b659f719d75c
Description: Production deployment pipeline

Tasks:
NAME      ID              DESCRIPTION              DEPENDS ON
----      --              -----------              ----------
setup     a8cad56b...     Setup environment       -
build     628e29e7...     Build application       setup
test      3f07511c...     Run tests               build
deploy    9b12345a...     Deploy to production    test
```

## 🗂️ Comandos de Stack

### `sloth-runner stack list`
Lista todos os stacks com métricas:

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
Detalhes completos de um stack:

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
Remove stack com confirmação:

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

## 📊 Outputs Exportados

### 💡 Como Exportar
No seu workflow Lua, use a variável global `outputs`:

```lua
TaskDefinitions = {
    deploy_pipeline = {
        description = "Production deployment",
        tasks = {
            {
                name = "deploy",
                description = "Deploy application",
                command = function()
                    print("Deploying application...")
                    
                    -- Exportar outputs que serão persistidos
                    outputs = outputs or {}
                    outputs.app_version = "1.2.3"
                    outputs.deployment_url = "https://app.example.com"
                    outputs.deployment_time = os.time()
                    outputs.status = "deployed"
                    
                    return true
                end
            }
        }
    }
}
```

### 📈 Visualização
Os outputs são:
- **Persistidos** no stack
- **Versionados** por execução
- **Visíveis** no `stack show`
- **Exportáveis** em JSON

## 🛠️ Installation & Build

### 🔨 Build Manual
```bash
# Clone e build
git clone https://github.com/chalkan3/sloth-runner.git
cd sloth-runner
go build -o sloth-runner ./cmd/sloth-runner/

# Instalar no sistema
cp sloth-runner ~/.local/bin/
```

### 📦 Preparação do Sistema
```bash
# Criar diretório do database
sudo mkdir -p /etc/sloth-runner
sudo chmod 777 /etc/sloth-runner
```

## 🚀 Casos de Uso

### 1. **CI/CD Pipeline**
```bash
# Pipeline de produção
sloth-runner run cicd-main -f cicd-pipeline.lua --output json > deployment-report.json

# Verificar status
sloth-runner stack show cicd-main

# Deployment para múltiplos ambientes
sloth-runner run production -f deploy.lua --output enhanced
sloth-runner run staging -f deploy.lua --output enhanced
sloth-runner run development -f deploy.lua --output enhanced
```

### 2. **Desenvolvimento Local**
```bash
# Setup de ambiente
sloth-runner run dev-setup -f local-env.lua --output modern

# Testes automatizados
sloth-runner run test-suite -f tests.lua --output enhanced

# Monitoramento
sloth-runner stack list
sloth-runner stack show test-suite
```

### 3. **Análise e Debugging**
```bash
# Listar estrutura do workflow
sloth-runner list -f complex-workflow.lua

# Output detalhado em JSON para análise
sloth-runner run debug-session -f debug.lua --output json | jq .

# Histórico de execuções
sloth-runner stack show debug-session
```

## 🔧 Próximas Funcionalidades

As seguintes funcionalidades estão sendo desenvolvidas:

- 🌐 **Remote stacks** (shared state)
- 🔄 **Stack templates** para rápida criação
- 📊 **Métricas avançadas** de performance
- 🔒 **RBAC** para stacks empresariais
- 🧪 **Stack testing** e validação
- 📈 **Dashboard web** para stacks

## 💡 Dicas e Truques

### Automatização CI/CD
```bash
#!/bin/bash
# deploy.sh
set -e

STACK_NAME="production-$(date +%Y%m%d-%H%M%S)"
RESULT=$(sloth-runner run $STACK_NAME -f production.lua --output json)
STATUS=$(echo $RESULT | jq -r '.status')

if [ "$STATUS" = "success" ]; then
    echo "✅ Deployment successful!"
    echo $RESULT | jq '.outputs'
else
    echo "❌ Deployment failed!"
    echo $RESULT | jq '.error'
    exit 1
fi
```

### Monitoramento
```bash
# Script de monitoramento
#!/bin/bash
while true; do
    echo "=== Stack Status ==="
    sloth-runner stack list
    echo ""
    sleep 30
done
```

### Cleanup de Stacks Antigos
```bash
# Limpar stacks antigos
sloth-runner stack list | grep "failed\|old" | awk '{print $1}' | xargs -I {} sloth-runner stack delete {} --force
```

---

📚 **Para mais informações:** [Stack Management Guide](stack-management.md) | [Getting Started](getting-started.md) | [Advanced Features](advanced-features.md)