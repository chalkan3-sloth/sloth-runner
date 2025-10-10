# Gerenciamento de Artifacts em Agents

## Visão Geral

**Artifacts** são arquivos produzidos por tasks durante a execução de workflows que precisam ser preservados, compartilhados entre tasks ou baixados para inspeção. O sistema de Agent Artifacts fornece uma solução completa para gerenciar esses arquivos através de agents distribuídos.

Pense em artifacts como os "outputs de build" dos seus workflows - binários compilados, relatórios de teste, logs, arquivos de configuração ou quaisquer outros arquivos que tasks geram e que tasks subsequentes ou humanos precisam consumir.

## Características Principais

- **Armazenamento Distribuído**: Artifacts armazenados em agents remotos
- **Rastreamento de Metadata**: Associar artifacts com stacks e tasks
- **Transferência via Streaming**: Manipulação eficiente de arquivos grandes
- **Gerenciamento de Ciclo de Vida**: Limpeza automática de artifacts antigos
- **Filtragem e Busca**: Encontrar artifacts por stack, task ou idade
- **Verificação de Checksum**: Garantir integridade dos dados

## Arquitetura

```
┌─────────────────────────────────────────────────────────┐
│                   Sistema de Artifacts                   │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Upload     │  │   Download   │  │    List      │  │
│  │   (Stream)   │  │   (Stream)   │  │  (Metadata)  │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Delete     │  │   Cleanup    │  │    Show      │  │
│  │  (Remover)   │  │  (Política)  │  │  (Detalhes)  │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
├─────────────────────────────────────────────────────────┤
│         Armazenamento: /var/lib/sloth-runner/           │
│              Metadata: SQLite + Checksums                │
└─────────────────────────────────────────────────────────┘
```

---

## Comandos CLI

### 1. Listar Artifacts

Lista todos os artifacts armazenados em um agent com filtragem opcional.

```bash
sloth-runner agent artifacts list <agent-name> [flags]
```

**Flags**:
- `--stack, -s <nome>` - Filtrar por nome do stack
- `--task, -t <nome>` - Filtrar por nome da task
- `--limit, -l <número>` - Máximo de artifacts a mostrar (padrão: 50)

**Exemplos**:

```bash
# Listar todos artifacts em um agent
sloth-runner agent artifacts list build-agent

# Filtrar por stack
sloth-runner agent artifacts list build-agent --stack production

# Filtrar por task e limitar resultados
sloth-runner agent artifacts list build-agent --task build --limit 10
```

**Saída**:
```
📦 Artifacts no agent: build-agent 📦

Nome                                     Stack           Tamanho      Criado               Task
---------------------------------------- --------------- ------------ -------------------- --------------------
app-v1.2.3.bin                          production      2.5 MB       2025-10-10 14:30:00  build
test-results.xml                        production      450 KB       2025-10-10 14:31:00  test
deployment-logs.tar.gz                  production      5.3 MB       2025-10-10 14:32:00  deploy

ℹ Total: 3 artifacts
```

---

### 2. Baixar Artifacts

Baixar um artifact de um agent remoto para sua máquina local.

```bash
sloth-runner agent artifacts download <agent-name> <artifact-name> [flags]
```

**Flags**:
- `--output, -o <caminho>` - Caminho de saída local (padrão: nome do artifact)

**Exemplos**:

```bash
# Baixar para diretório atual
sloth-runner agent artifacts download build-agent app.bin

# Especificar caminho de saída
sloth-runner agent artifacts download build-agent app.bin --output ./dist/app

# Baixar para localização específica
sloth-runner agent artifacts download build-agent report.pdf -o /tmp/reports/latest.pdf
```

**Saída**:
```
ℹ Baixando artifact 'app.bin' do agent 'build-agent'...
✓ Download de app.bin para 'app.bin' concluído (2.5 MB)
```

---

### 3. Enviar Artifacts

Enviar um arquivo local como artifact para um agent remoto.

```bash
sloth-runner agent artifacts upload <agent-name> <caminho-arquivo> [flags]
```

**Flags**:
- `--stack, -s <nome>` - Associar com stack
- `--task, -t <nome>` - Associar com task
- `--name, -n <nome>` - Nome do artifact (padrão: nome do arquivo)

**Exemplos**:

```bash
# Upload simples
sloth-runner agent artifacts upload build-agent ./app.bin

# Associar com stack e task
sloth-runner agent artifacts upload build-agent ./app.bin \
  --stack production \
  --task build

# Upload com nome customizado
sloth-runner agent artifacts upload build-agent ./binary \
  --name app-v2.0.bin \
  --stack production
```

**Saída**:
```
ℹ Enviando './app.bin' para agent 'build-agent' como 'app.bin'...
✓ Upload de app.bin concluído (2.5 MB)
```

---

### 4. Mostrar Detalhes do Artifact

Exibir informações detalhadas sobre um artifact específico.

```bash
sloth-runner agent artifacts show <agent-name> <artifact-name>
```

**Exemplo**:

```bash
sloth-runner agent artifacts show build-agent app.bin
```

**Saída**:
```
📦 Detalhes do Artifact

Nome:       app.bin
Caminho:    /var/lib/sloth-runner/artifacts/production/build/app.bin
Tamanho:    2.5 MB
Checksum:   sha256:a1b2c3d4e5f6...
Stack:      production
Task:       build
Criado:     2025-10-10 14:30:00
Modificado: 2025-10-10 14:30:00
Downloads:  5
```

---

### 5. Deletar Artifacts

Remover um artifact do armazenamento de um agent.

```bash
sloth-runner agent artifacts delete <agent-name> <artifact-name> [flags]
```

**Flags**:
- `--force, -f` - Pular prompt de confirmação

**Exemplos**:

```bash
# Com confirmação
sloth-runner agent artifacts delete build-agent old-app.bin

# Deletar forçado sem confirmação
sloth-runner agent artifacts delete build-agent old-app.bin --force
```

**Saída**:
```
Tem certeza que deseja deletar o artifact 'old-app.bin' do agent 'build-agent'? (y/N): y
✓ Artifact 'old-app.bin' deletado do agent 'build-agent'
```

---

### 6. Limpar Artifacts Antigos

Remover automaticamente artifacts mais antigos que uma duração especificada.

```bash
sloth-runner agent artifacts cleanup <agent-name> [flags]
```

**Flags**:
- `--older-than <duração>` - Remover artifacts mais antigos que isto (padrão: 30d)
- `--stack, -s <nome>` - Limitar limpeza a stack específico
- `--dry-run` - Visualizar o que seria deletado sem deletar

**Formato de Duração**:
- `7d` - 7 dias
- `30d` - 30 dias
- `6h` - 6 horas
- `1h` - 1 hora

**Exemplos**:

```bash
# Limpeza padrão (30 dias)
sloth-runner agent artifacts cleanup build-agent

# Limpar artifacts com mais de 7 dias
sloth-runner agent artifacts cleanup build-agent --older-than 7d

# Limpar apenas stack específico
sloth-runner agent artifacts cleanup build-agent --stack old-project

# Visualizar sem deletar (dry-run)
sloth-runner agent artifacts cleanup build-agent --older-than 7d --dry-run
```

**Saída**:
```
✓ Deletados 15 artifacts, liberados 45.2 MB
```

---

## Integração com Workflows

Artifacts são gerenciados automaticamente quando declarados em definições de tasks usando os métodos `:artifacts()` e `:consumes()`.

### Produzindo Artifacts

Use `:artifacts()` para declarar arquivos que devem ser salvos após a execução da task:

```lua
local build_task = task("build")
    :description("Compilar binário da aplicação")
    :command(function()
        -- Compilar a aplicação
        exec.run("go build -o app.bin ./cmd/app")

        log.info("Build concluído com sucesso")
        return true, "Binário criado: app.bin"
    end)
    :artifacts({"app.bin"})  -- Declarar artifact
    :build()
```

### Consumindo Artifacts

Use `:consumes()` para acessar artifacts de tasks anteriores:

```lua
local test_task = task("test")
    :description("Executar testes com binário compilado")
    :depends_on({"build"})
    :consumes({"app.bin"})  -- Consumir artifact da task build
    :command(function()
        -- O artifact é automaticamente copiado para o workdir desta task
        exec.run("chmod +x app.bin")
        exec.run("./app.bin --test")

        return true, "Testes passaram"
    end)
    :build()
```

### Exemplo Completo de CI/CD

```lua
-- Estágio de build
local build = task("build")
    :command(function()
        exec.run("go build -o app.bin")
        return true
    end)
    :artifacts({"app.bin"})
    :build()

-- Estágio de teste
local test = task("test")
    :depends_on({"build"})
    :consumes({"app.bin"})
    :command(function()
        exec.run("./app.bin --test")
        -- Gerar relatório de teste
        exec.run("./generate-report.sh > test-report.xml")
        return true
    end)
    :artifacts({"test-report.xml"})
    :build()

-- Estágio de deploy
local deploy = task("deploy")
    :depends_on({"test"})
    :consumes({"app.bin"})
    :command(function()
        exec.run("scp app.bin production:/opt/app/")
        return true
    end)
    :build()

-- Definir workflow
workflow.define("ci_pipeline")
    :description("Pipeline CI/CD completo com artifacts")
    :version("1.0.0")
    :tasks({build, test, deploy})
    :config({
        timeout = "30m",
        create_workdir_before_run = true
    })
```

---

## Casos de Uso

### 1. Pipeline CI/CD

**Cenário**: Compilar uma vez, deployar em qualquer lugar

```bash
# 1. Executar build no agent de build
sloth-runner run ci_build --file build.sloth --agent build-agent

# 2. Baixar artifact
sloth-runner agent artifacts download build-agent app-v1.2.3.bin

# 3. Enviar para agent de deployment
sloth-runner agent artifacts upload deploy-agent app-v1.2.3.bin \
  --stack production \
  --task deploy

# 4. Executar deployment
sloth-runner run deploy --file deploy.sloth --agent deploy-agent
```

### 2. Debug de Workflows Falhados

**Cenário**: Investigar o que deu errado

```bash
# Listar artifacts do workflow que falhou
sloth-runner agent artifacts list prod-agent --task failed-task

# Baixar logs de erro
sloth-runner agent artifacts download prod-agent error.log

# Inspecionar localmente
cat error.log

# Baixar core dump se disponível
sloth-runner agent artifacts download prod-agent core.dump
```

### 3. Política de Retenção de Artifacts

**Cenário**: Manter armazenamento sob controle

```bash
# Limpeza semanal de artifacts antigos
sloth-runner agent artifacts cleanup build-agent --older-than 30d

# Limpar projetos antigos específicos
sloth-runner agent artifacts cleanup build-agent \
  --stack legacy-project \
  --older-than 7d

# Visualizar limpeza antes de executar
sloth-runner agent artifacts cleanup build-agent \
  --older-than 30d \
  --dry-run
```

---

## Melhores Práticas

### 1. Convenções de Nomenclatura

Use nomes descritivos e versionados:

```bash
# Bom
app-v1.2.3.bin
report-2025-10-10.pdf
logs-production-20251010.tar.gz

# Evite
output.txt
file.bin
temp.log
```

### 2. Sempre Associe Metadata

Vincule artifacts a stacks e tasks:

```bash
sloth-runner agent artifacts upload build-agent app.bin \
  --stack production \
  --task build-v1.2.3
```

### 3. Limpeza Regular

Agende limpeza automática para evitar problemas de armazenamento:

```bash
# Exemplo de cron job (limpeza semanal)
0 2 * * 0 sloth-runner agent artifacts cleanup build-agent --older-than 30d
```

### 4. Monitoramento de Tamanho

Monitore o tamanho dos artifacts:

```bash
# Listar artifacts ordenados por tamanho
sloth-runner agent artifacts list build-agent --limit 100 | \
  sort -k3 -hr | \
  head -10
```

---

## Troubleshooting

### Artifact Não Encontrado

```bash
# Verificar se artifact existe
sloth-runner agent artifacts list build-agent | grep artifact-name

# Mostrar detalhes completos
sloth-runner agent artifacts show build-agent artifact-name
```

### Falhas no Download

```bash
# Verificar conectividade do agent
sloth-runner agent get build-agent

# Tentar novamente com timeout maior
timeout 600 sloth-runner agent artifacts download build-agent large-file.bin
```

### Armazenamento Cheio

```bash
# Limpeza agressiva
sloth-runner agent artifacts cleanup build-agent --older-than 1d

# Listar maiores artifacts
sloth-runner agent artifacts list build-agent --limit 100 | \
  sort -k3 -hr | \
  head -20
```

---

## Considerações de Performance

### Otimização de Transferência

- **Streaming**: Arquivos transferidos em chunks de 64KB
- **Compressão**: Compressão gzip para arquivos grandes
- **Paralelo**: Múltiplos artifacts podem ser baixados concorrentemente

### Eficiência de Armazenamento

- **Deduplicação**: Mesmo checksum = mesmo arquivo (feature futura)
- **Limpeza**: Limpeza regular previne inchaço de armazenamento
- **Compressão**: Armazenar comprimido quando apropriado

---

## Documentação Relacionada

- [DSL de Workflow - Artifacts](./core-concepts.md#artifact-management)
- [Comandos de Agent](./CLI.md#agent-commands)
- [Sistema de Eventos](./advanced-features.md#events)
- [Configuração de Armazenamento](./getting-started.md#configuration)

---

## FAQ

**P: Qual o tamanho máximo de um artifact?**
R: Sem limite rígido, mas streaming é usado para transferência eficiente. Arquivos de 10GB+ são suportados.

**P: Artifacts são versionados automaticamente?**
R: Ainda não - use convenções de nomenclatura (app-v1.0.0.bin) até que versionamento automático seja implementado.

**P: Posso compartilhar artifacts entre stacks?**
R: Sim, baixe de um stack e envie para outro.

**P: Por quanto tempo artifacts são mantidos?**
R: Para sempre, a menos que você use `cleanup`. Implemente políticas de retenção para limpeza automática.

**P: Posso usar artifacts com workflows locais?**
R: Sim, artifacts funcionam com execução local e remota em agents.

---

*Última atualização: 2025-10-10*
*Versão: 1.0.0*
