# Início Rápido

Bem-vindo ao Sloth-Runner! Este guia o ajudará a começar a usar a ferramenta rapidamente com as novas funcionalidades de gerenciamento de stacks e output estilo Pulumi.

> **📝 Nota Importante:** A partir da versão atual, os arquivos de workflow do Sloth Runner usam a extensão `.sloth` em vez de `.lua`. A sintaxe Lua permanece a mesma, apenas a extensão do arquivo mudou para melhor identificação dos arquivos DSL do Sloth Runner.

## Instalação

Para instalar o `sloth-runner` em seu sistema, você pode usar o script `install.sh` fornecido. Este script detecta automaticamente seu sistema operacional e arquitetura, baixa a versão mais recente do GitHub e coloca o executável `sloth-runner` em `/usr/local/bin`.

```bash
bash <(curl -sL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh)
```

**Nota:** O script `install.sh` requer privilégios de `sudo` para mover o executável para `/usr/local/bin`.

## 🗂️ Novo: Stack Management (Recomendado)

### Execução com Stack

A nova funcionalidade principal do Sloth Runner é o **Stack Management**, similar ao Pulumi:

```bash
# Nova sintaxe - nome do stack como argumento posicional
sloth-runner run {nome-do-stack} --file workflow.sloth

# Exemplos práticos
sloth-runner run production-app -f deploy.sloth --output enhanced
sloth-runner run dev-environment -f test.sloth -o rich
sloth-runner run staging-api -f pipeline.sloth
```

### Gerenciamento de Stacks

```bash
# Criar um novo stack
sloth-runner stack new my-stack --description "Production deployment stack" --workflow-file deploy.sloth

# Criar stack com valores padrão
sloth-runner stack new dev-stack

# Listar todos os stacks
sloth-runner stack list

# Ver detalhes e outputs exportados
sloth-runner stack show production-app

# Remover stack antigo
sloth-runner stack delete old-environment
```

### 🆔 Listagem de Tasks e Grupos (Novo)

```bash
# Listar tasks e grupos com IDs únicos
sloth-runner list -f workflow.sloth

# Visualizar estrutura completa do workflow
sloth-runner list -f pipeline.sloth

# Exemplo de saída organizada:
# ## Task Group: deploy_group  
# ID: e8e77f72-5cf4-4e98-adce-fc839846c24a
# Tasks:
# NAME     ID           DESCRIPTION             DEPENDS ON
# build    a1c4fa46...  Build the application   -
# test     d8dc4623...  Run tests               build
# deploy   6253cb19...  Deploy to production    build, test
```

## 🚀 Scaffolding de Projetos

### Criar Novo Projeto

```bash
# Criar projeto a partir de template
sloth-runner workflow init meu-app --template cicd

# Listar templates disponíveis
sloth-runner workflow list-templates

# Executar o workflow gerado
cd meu-app
sloth-runner run dev-env -f meu-app.sloth --output enhanced
```

### Templates Disponíveis

- **basic** - Workflow básico com uma task
- **cicd** - Pipeline CI/CD completo
- **infrastructure** - Deploy de infraestrutura
- **microservices** - Deploy de microserviços
- **data-pipeline** - Pipeline de processamento de dados

## 🎨 Estilos de Output

### Output Configurável

```bash
# Output básico (padrão)
sloth-runner run meu-stack -f workflow.sloth

# Output melhorado estilo Pulumi
sloth-runner run meu-stack -f workflow.sloth --output enhanced
sloth-runner run meu-stack -f workflow.sloth -o rich
sloth-runner run meu-stack -f workflow.sloth --output modern
```

### Demonstração Visual

Com `--output enhanced`:

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

## 💡 Uso Básico (Modo Tradicional)

Para executar um arquivo de tarefa Lua sem stack:

```bash
# Execução simples
sloth-runner run -f examples/basic_pipeline.sloth

# Com output melhorado
sloth-runner run -f examples/basic_pipeline.sloth --output enhanced
```

Para listar as tarefas em um arquivo:

```bash
sloth-runner list -f examples/basic_pipeline.sloth
```

## 📊 Exemplos Práticos

### Deploy Multi-Ambiente

```bash
# Desenvolvimento
sloth-runner run dev-app -f deploy.sloth

# Staging
sloth-runner run staging-app -f deploy.sloth

# Produção com output rico
sloth-runner run prod-app -f deploy.sloth --output enhanced

# Verificar status de produção
sloth-runner stack show prod-app
```

### CI/CD Integration

```bash
# No pipeline CI/CD
sloth-runner run ${ENVIRONMENT}-${APP_NAME} -f pipeline.sloth

# Exemplo específico
sloth-runner run prod-frontend -f frontend-deploy.sloth
```

## 🗃️ Persistência de Estado

Os stacks são automaticamente persistidos em:

```
~/.sloth-runner/stacks.db
```

Cada stack mantém:
- Status atual da execução
- Outputs exportados da pipeline
- Histórico completo de execuções
- Metadados e configurações

## ⏰ Agendador de Tarefas

O Sloth-Runner inclui um poderoso agendador de tarefas que permite automatizar a execução de seus fluxos de trabalho em segundo plano usando sintaxe cron. Para mais detalhes sobre como configurar e usar o agendador, consulte a documentação completa em [Agendador de Tarefas](./scheduler.md).

## 📚 Próximos Passos

Agora que você tem o Sloth-Runner instalado e funcionando com as novas funcionalidades:

- Explore o [Stack Management](./stack-management.md) para gerenciamento avançado de estado
- Veja os [Conceitos Essenciais](./core-concepts.md) para entender como definir suas tarefas
- Experimente os [Módulos Built-in](./index.md#módulos-built-in) para automação avançada
- Consulte [Exemplos Avançados](./advanced-examples.md) para casos de uso complexos
