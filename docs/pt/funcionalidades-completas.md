# 🚀 Funcionalidades Completas do Sloth Runner

## Visão Geral

Documentação completa de **todas** as funcionalidades do Sloth Runner - desde recursos básicos até funcionalidades enterprise avançadas. Este guia serve como índice mestre para explorar todas as capacidades da plataforma.

---

## 📋 Índice de Funcionalidades

### 🎯 Núcleo (Core)
- [Execução de Workflows](#execução-de-workflows)
- [Linguagem DSL Sloth](#linguagem-dsl-sloth)
- [Sistema de Módulos](#sistema-de-módulos)
- [State Management](#state-management)
- [Idempotência](#idempotência)

### 🌐 Distribuído
- [Arquitetura Master-Agent](#arquitetura-master-agent)
- [Delegação de Tarefas](#delegação-de-tarefas)
- [Comunicação gRPC](#comunicação-grpc)
- [Auto-Reconnection](#auto-reconnection)
- [Health Checks](#health-checks)

### 🎨 Interface
- [Web UI Moderna](#web-ui-moderna)
- [CLI Completo](#cli-completo)
- [REPL Interativo](#repl-interativo)
- [Terminal Remoto](#terminal-remoto)
- [API REST](#api-rest)

### 🔧 Automação
- [Scheduler (Cron)](#scheduler)
- [Hooks & Events](#hooks--events)
- [GitOps](#gitops)
- [CI/CD Integration](#cicd-integration)
- [Workflows Salvos (Sloths)](#sloths)

### 📊 Monitoramento
- [Telemetria](#telemetria)
- [Prometheus Metrics](#prometheus-metrics)
- [Grafana Dashboards](#grafana-dashboards)
- [Logs Centralizados](#logs-centralizados)
- [Agent Metrics](#agent-metrics)

### ☁️ Cloud & IaC
- [Multi-Cloud](#multi-cloud)
- [Terraform](#terraform)
- [Pulumi](#pulumi)
- [Kubernetes](#kubernetes)
- [Docker](#docker)

### 🔐 Segurança & Enterprise
- [Autenticação](#autenticação)
- [TLS/SSL](#tlsssl)
- [Audit Logs](#audit-logs)
- [Backups](#backups)
- [RBAC](#rbac)

### 🚀 Performance
- [Otimizações](#otimizações)
- [Parallel Execution](#parallel-execution)
- [Resource Limits](#resource-limits)
- [Caching](#caching)

---

## 🎯 Núcleo (Core)

### Execução de Workflows

**Descrição:** Motor central para execução de workflows definidos em arquivos Sloth.

**Características:**
- Execução sequencial e paralela de tasks
- Suporte a grupos de tasks
- Variáveis e templating
- Conditional execution
- Error handling e retry
- Dry-run mode
- Verbose output

**Comandos:**
```bash
sloth-runner run <workflow> --file <arquivo>
sloth-runner run <workflow> --file <arquivo> --yes
sloth-runner run <workflow> --file <arquivo> --group <grupo>
sloth-runner run <workflow> --file <arquivo> --values vars.yaml
```

**Exemplos:**
```yaml
# Workflow básico
tasks:
  - name: Install nginx
    exec:
      script: |
        pkg.update()
        pkg.install("nginx")

  - name: Configure nginx
    exec:
      script: |
        file.copy("/src/nginx.conf", "/etc/nginx/nginx.conf")
        systemd.service_restart("nginx")
```

**Documentação:** `/docs/en/quick-start.md`

---

### Linguagem DSL Sloth

**Descrição:** DSL declarativa baseada em YAML com scripting Lua embarcado.

**Features:**
- **YAML-based** - sintaxe familiar e legível
- **Lua scripting** - poder de uma linguagem completa
- **Type-safe** - validação de tipos
- **Templating** - Go templates e Jinja2
- **Módulos globais** - sem necessidade de require()
- **Modern syntax** - suporte a features modernas

**Estrutura:**
```yaml
# Metadata
version: "1.0"
description: "Meu workflow"

# Variáveis
vars:
  env: production
  version: "1.2.3"

# Groups
groups:
  deploy:
    - install_deps
    - build_app
    - deploy_app

# Tasks
tasks:
  - name: install_deps
    exec:
      script: |
        pkg.install({"nodejs", "npm"})

  - name: build_app
    exec:
      script: |
        exec.command("npm install")
        exec.command("npm run build")

  - name: deploy_app
    exec:
      script: |
        file.copy("./dist", "/var/www/app")
        systemd.service_restart("app")
    delegate_to: web-01
```

**Documentação:** `/docs/modern-dsl/introduction.md`

---

### Sistema de Módulos

**Descrição:** 40+ módulos integrados para todas as necessidades de automação.

**Categorias:**

#### 📦 Sistema
- `pkg` - Gerenciamento de pacotes (apt, yum, brew, etc.)
- `user` - Gerenciamento de usuários/grupos
- `file` - Operações com arquivos
- `systemd` - Gerenciamento de serviços
- `exec` - Execução de comandos

#### 🐳 Containers
- `docker` - Docker completo (containers, images, networks)
- `incus` - Incus/LXC containers e VMs
- `kubernetes` - Deploy e gerenciamento K8s

#### ☁️ Cloud
- `aws` - AWS (EC2, S3, RDS, Lambda, etc.)
- `azure` - Azure (VMs, Storage, etc.)
- `gcp` - GCP (Compute Engine, Cloud Storage, etc.)
- `digitalocean` - DigitalOcean (Droplets, Load Balancers)

#### 🏗️ IaC
- `terraform` - Terraform (init, plan, apply, destroy)
- `pulumi` - Pulumi
- `ansible` - Ansible playbooks

#### 🔧 Ferramentas
- `git` - Operações Git
- `ssh` - SSH remoto
- `net` - Networking (ping, http, download)
- `template` - Templates (Jinja2, Go)

#### 📊 Observabilidade
- `log` - Logging estruturado
- `metrics` - Métricas (Prometheus)
- `notifications` - Notificações (Slack, Email, Discord, Telegram)

#### 🚀 Avançado
- `goroutine` - Execução paralela
- `reliability` - Retry, circuit breaker, timeout
- `state` - State management
- `facts` - System information
- `infra_test` - Infrastructure testing

**Lista completa:** `sloth-runner modules list`

**Documentação:** `/docs/pt/modulos-completos.md`

---

### State Management

**Descrição:** Sistema de persistência de estado entre execuções.

**Features:**
- Key-value store persistente
- SQLite backend
- State scoping (global, workflow, task)
- Change detection
- State cleanup

**API:**
```lua
-- Salvar estado
state.set("last_deploy_version", "v1.2.3")
state.set("deploy_timestamp", os.time())

-- Ler estado
local last_version = state.get("last_deploy_version")

-- Detectar mudança
if state.changed("config_hash", new_hash) then
    log.info("Config changed, redeploying")
    deploy()
end

-- Limpar estado
state.clear("temporary_data")
```

**Documentação:** `/docs/state-management.md`

---

### Idempotência

**Descrição:** Garantia de que workflows podem ser executados múltiplas vezes com o mesmo resultado.

**Features:**
- **Check mode** - verifica antes de executar
- **State tracking** - rastreia o que foi mudado
- **Resource fingerprinting** - detecta mudanças
- **Rollback** - desfaz mudanças em caso de erro

**Exemplo:**
```lua
-- Idempotente - verifica antes de instalar
if not pkg.is_installed("nginx") then
    pkg.install("nginx")
end

-- Idempotente - verifica hash do arquivo
local current_hash = file.hash("/etc/nginx/nginx.conf")
if current_hash ~= expected_hash then
    file.copy("/src/nginx.conf", "/etc/nginx/nginx.conf")
    systemd.service_restart("nginx")
end
```

**Documentação:** `/docs/idempotency.md`

---

## 🌐 Distribuído

### Arquitetura Master-Agent

**Descrição:** Arquitetura distribuída com servidor master central e agentes remotos.

**Componentes:**
- **Master Server** - coordena agentes e workflows
- **Agent Nodes** - executam tarefas remotamente
- **gRPC Communication** - comunicação eficiente e type-safe
- **Auto-Discovery** - agentes se auto-registram
- **Health Monitoring** - heartbeats automáticos

**Topologia:**
```
                    ┌──────────────┐
                    │   Master     │
                    │  (gRPC:50053)│
                    └──────┬───────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐       ┌────▼────┐      ┌────▼────┐
    │ Agent 1 │       │ Agent 2 │      │ Agent 3 │
    │  web-01 │       │  web-02 │      │   db-01 │
    └─────────┘       └─────────┘      └─────────┘
```

**Setup:**
```bash
# Iniciar master
sloth-runner server --port 50053

# Instalar agente remoto
sloth-runner agent install web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user ubuntu \
  --master 192.168.1.1:50053

# Listar agentes
sloth-runner agent list
```

**Documentação:** `/docs/en/master-agent-architecture.md`

---

### Delegação de Tarefas

**Descrição:** Delega execução de tasks para agentes específicos.

**Features:**
- **Single delegation** - delega para um agente
- **Multi delegation** - delega para múltiplos agentes em paralelo
- **Round-robin** - distribui carga
- **Failover** - fallback se agente falhar
- **Conditional delegation** - delega baseado em condições

**Sintaxe:**
```yaml
# Delegar para um agente
tasks:
  - name: Deploy to web-01
    exec:
      script: |
        pkg.install("nginx")
    delegate_to: web-01

# Delegar para múltiplos agentes
tasks:
  - name: Deploy to all web servers
    exec:
      script: |
        pkg.install("nginx")
    delegate_to:
      - web-01
      - web-02
      - web-03

# CLI - delegar workflow inteiro
sloth-runner run deploy --file deploy.sloth --delegate-to web-01
```

**Uso com valores:**
```yaml
# Passar valores específicos por agente
tasks:
  - name: Configure
    exec:
      script: |
        local ip = values.ip_address
        file.write("/etc/config", "IP=" .. ip)
    delegate_to: "{{ item }}"
    loop:
      - web-01
      - web-02
    values:
      web-01:
        ip_address: "192.168.1.10"
      web-02:
        ip_address: "192.168.1.11"
```

**Documentação:** `/docs/guides/values-delegate-to.md`

---

### Comunicação gRPC

**Descrição:** Comunicação eficiente entre master e agentes usando gRPC.

**Features:**
- **Streaming** - bi-directional streaming
- **Type-safe** - Protocol Buffers
- **Efficient** - binary protocol
- **Multiplexing** - múltiplas chamadas em uma conexão
- **TLS** - suporte a TLS/SSL

**Serviços:**
```protobuf
service AgentService {
    rpc ExecuteTask(TaskRequest) returns (TaskResponse);
    rpc StreamLogs(LogRequest) returns (stream LogEntry);
    rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse);
    rpc GetMetrics(MetricsRequest) returns (MetricsResponse);
}
```

**Porta padrão:** 50053

---

### Auto-Reconnection

**Descrição:** Agentes se reconectam automaticamente ao master em caso de desconexão.

**Features:**
- **Exponential backoff** - aumenta intervalo entre tentativas
- **Max retries** - limite configurável
- **Circuit breaker** - para de tentar após muitas falhas
- **Connection pooling** - reusa conexões

**Configuração:**
```yaml
agent:
  reconnect:
    enabled: true
    initial_delay: 1s
    max_delay: 60s
    max_retries: -1  # infinito
```

**Documentação:** `/docs/en/agent-improvements.md`

---

### Health Checks

**Descrição:** Monitoramento contínuo de saúde dos agentes.

**Tipos de checks:**
- **Heartbeat** - ping periódico
- **Resource check** - CPU, memory, disk
- **Service check** - verifica serviços críticos
- **Custom checks** - checks definidos pelo usuário

**Endpoints:**
```bash
# Health endpoint
curl http://agent:9090/health

# Metrics endpoint
curl http://agent:9090/metrics
```

**Thresholds:**
```yaml
health:
  cpu_threshold: 90  # %
  memory_threshold: 85  # %
  disk_threshold: 90  # %
  heartbeat_interval: 30s
  heartbeat_timeout: 90s
```

---

## 🎨 Interface

### Web UI Moderna

**Descrição:** Interface web completa, responsiva e em tempo real.

**Features principais:**
- ✅ Dashboard com métricas e gráficos
- ✅ Gerenciamento de agentes com métricas em tempo real
- ✅ Editor de workflows com syntax highlighting
- ✅ Visualização de execuções e logs
- ✅ Terminal interativo (xterm.js)
- ✅ Dark mode / Light mode
- ✅ WebSocket para atualizações em tempo real
- ✅ Mobile responsive
- ✅ Command palette (Ctrl+Shift+P)
- ✅ Drag & drop
- ✅ Glassmorphism design
- ✅ Animações suaves

**Páginas:**
1. Dashboard (`/`)
2. Agents (`/agents`)
3. Agent Control (`/agent-control`)
4. Agent Dashboard (`/agent-dashboard`)
5. Workflows (`/workflows`)
6. Executions (`/executions`)
7. Hooks (`/hooks`)
8. Events (`/events`)
9. Scheduler (`/scheduler`)
10. Logs (`/logs`)
11. Terminal (`/terminal`)
12. Sloths (`/sloths`)
13. Settings (`/settings`)

**Tecnologias:**
- Bootstrap 5.3
- Chart.js 4.4
- xterm.js
- WebSockets
- Canvas API

**Iniciar:**
```bash
sloth-runner ui --port 8080
```

**Acesso:** http://localhost:8080

**Documentação:** `/docs/pt/web-ui-completo.md`

---

### CLI Completo

**Descrição:** Interface de linha de comando completa com 100+ comandos.

**Categorias de comandos:**

#### Execução
- `run` - Executar workflow
- `version` - Ver versão

#### Agentes
- `agent list` - Listar agentes
- `agent get` - Detalhes do agente
- `agent install` - Instalar agente remoto
- `agent update` - Atualizar agente
- `agent start/stop/restart` - Controlar agente
- `agent modules` - Listar módulos do agente
- `agent metrics` - Ver métricas

#### Sloths (Workflows Salvos)
- `sloth list` - Listar sloths
- `sloth add` - Adicionar sloth
- `sloth get` - Ver sloth
- `sloth update` - Atualizar sloth
- `sloth remove` - Remover sloth
- `sloth activate/deactivate` - Ativar/desativar

#### Hooks
- `hook list` - Listar hooks
- `hook add` - Adicionar hook
- `hook remove` - Remover hook
- `hook enable/disable` - Habilitar/desabilitar
- `hook test` - Testar hook

#### Eventos
- `events list` - Listar eventos
- `events watch` - Monitorar eventos em tempo real

#### Database
- `db backup` - Backup do database
- `db restore` - Restaurar database
- `db vacuum` - Otimizar database
- `db stats` - Estatísticas

#### SSH
- `ssh list` - Listar conexões SSH
- `ssh add` - Adicionar conexão
- `ssh remove` - Remover conexão
- `ssh test` - Testar conexão

#### Módulos
- `modules list` - Listar módulos
- `modules info` - Info do módulo

#### Servidor
- `server` - Iniciar servidor master
- `ui` - Iniciar Web UI
- `terminal` - Terminal interativo

#### Utilitários
- `completion` - Auto-completar shell
- `doctor` - Diagnóstico

**Documentação:** `/docs/pt/referencia-cli.md`

---

### REPL Interativo

**Descrição:** Read-Eval-Print Loop para testar código Lua interativamente.

**Features:**
- **Lua completo** - interpretador Lua completo
- **Módulos carregados** - todos os módulos disponíveis
- **History** - histórico de comandos
- **Auto-complete** - Tab completion
- **Multi-line** - suporte a código multi-linha
- **Pretty print** - output formatado

**Iniciar:**
```bash
sloth-runner repl
```

**Exemplo de sessão:**
```lua
> pkg.install("nginx")
[OK] nginx installed successfully

> file.exists("/etc/nginx/nginx.conf")
true

> local content = file.read("/etc/nginx/nginx.conf")
> print(#content .. " bytes")
2048 bytes

> for i=1,5 do
>>   print("Hello " .. i)
>> end
Hello 1
Hello 2
Hello 3
Hello 4
Hello 5
```

**Comandos especiais:**
- `.help` - ajuda
- `.exit` - sair
- `.clear` - limpar tela
- `.load <file>` - carregar arquivo
- `.save <file>` - salvar sessão

**Documentação:** `/docs/en/repl.md`

---

### Terminal Remoto

**Descrição:** Terminal interativo para agentes remotos via web UI.

**Features:**
- **xterm.js** - emulador de terminal completo
- **Multiple sessions** - múltiplas sessões simultâneas
- **Tabs** - gerenciamento em tabs
- **Command history** - histórico de comandos (↑↓)
- **Copy/paste** - Ctrl+Shift+C/V
- **Themes** - vários temas disponíveis
- **Upload/download** - transferência de arquivos

**Acesso:**
1. Web UI → Terminal
2. Selecionar agente
3. Conectar

**Comandos especiais:**
```bash
.clear       # Limpar terminal
.exit        # Fechar sessão
.upload <f>  # Upload arquivo
.download <f># Download arquivo
.theme <t>   # Mudar tema
```

**URL:** http://localhost:8080/terminal

---

### API REST

**Descrição:** API RESTful completa para integração externa.

**Endpoints principais:**

#### Agentes
```
GET    /api/v1/agents           # Lista agentes
GET    /api/v1/agents/:name     # Detalhes do agente
POST   /api/v1/agents/:name/restart  # Reinicia agente
DELETE /api/v1/agents/:name     # Remove agente
```

#### Workflows
```
POST   /api/v1/workflows/run    # Executa workflow
GET    /api/v1/workflows/:id    # Detalhes do workflow
```

#### Execuções
```
GET    /api/v1/executions       # Lista execuções
GET    /api/v1/executions/:id   # Detalhes da execução
```

#### Hooks
```
GET    /api/v1/hooks            # Lista hooks
POST   /api/v1/hooks            # Cria hook
DELETE /api/v1/hooks/:name      # Remove hook
```

#### Eventos
```
GET    /api/v1/events           # Lista eventos
```

#### Métricas
```
GET    /api/v1/metrics          # Métricas Prometheus
```

**Autenticação:**
```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/agents
```

**Exemplos:**
```bash
# Listar agentes
curl http://localhost:8080/api/v1/agents

# Executar workflow
curl -X POST http://localhost:8080/api/v1/workflows/run \
  -H "Content-Type: application/json" \
  -d '{
    "file": "/workflows/deploy.sloth",
    "workflow_name": "deploy",
    "delegate_to": ["web-01"]
  }'

# Ver métricas
curl http://localhost:8080/api/v1/metrics
```

**Documentação:** `/docs/web-ui/api-reference.md`

---

## 🔧 Automação

### Scheduler

**Descrição:** Agendador de workflows baseado em cron.

**Features:**
- **Cron expressions** - sintaxe cron completa
- **Visual builder** - construtor visual na Web UI
- **Timezone support** - suporte a timezones
- **Missed run policy** - política para runs perdidos
- **Overlap prevention** - previne execuções sobrepostas
- **Notifications** - notificações de sucesso/falha

**Criar job:**
```bash
# Via CLI (em breve)
sloth-runner scheduler add deploy-job \
  --workflow deploy.sloth \
  --schedule "0 3 * * *"  # Todo dia às 3h

# Via Web UI
http://localhost:8080/scheduler
```

**Sintaxe cron:**
```
┌───────────── minuto (0 - 59)
│ ┌───────────── hora (0 - 23)
│ │ ┌───────────── dia do mês (1 - 31)
│ │ │ ┌───────────── mês (1 - 12)
│ │ │ │ ┌───────────── dia da semana (0 - 6) (Domingo=0)
│ │ │ │ │
* * * * *

Exemplos:
0 * * * *     # A cada hora
0 3 * * *     # Todo dia às 3h
0 0 * * 0     # Todo domingo à meia-noite
*/15 * * * *  # A cada 15 minutos
```

**Documentação:** `/docs/pt/scheduler.md`

---

### Hooks & Events

**Descrição:** Sistema de hooks para reagir a eventos do sistema.

**Eventos disponíveis:**
- `workflow.started` - Workflow iniciado
- `workflow.completed` - Workflow completado
- `workflow.failed` - Workflow falhou
- `task.started` - Task iniciada
- `task.completed` - Task completada
- `task.failed` - Task falhou
- `agent.connected` - Agente conectado
- `agent.disconnected` - Agente desconectado

**Criar hook:**
```bash
sloth-runner hook add slack-notify \
  --event workflow.completed \
  --script /scripts/notify-slack.lua \
  --priority 10
```

**Script de hook (Lua):**
```lua
-- /scripts/notify-slack.lua
local event = hook.event
local payload = hook.payload

if event == "workflow.completed" then
    notifications.slack(
        "https://hooks.slack.com/services/XXX/YYY/ZZZ",
        string.format("✅ Workflow '%s' completed!", payload.workflow_name),
        { channel = "#deployments" }
    )
end
```

**Payload disponível:**
```lua
-- workflow.* events
{
    workflow_name = "deploy",
    status = "success",
    duration = 45.3,
    started_at = 1234567890,
    completed_at = 1234567935
}

-- agent.* events
{
    agent_name = "web-01",
    address = "192.168.1.100:50060",
    status = "connected"
}
```

**Documentação:** `/docs/architecture/hooks-events-system.md`

---

### GitOps

**Descrição:** Implementação completa de padrões GitOps.

**Features:**
- **Git-based** - Git como source of truth
- **Auto-sync** - sincronização automática
- **Drift detection** - detecta mudanças manuais
- **Rollback** - rollback automático
- **Multi-environment** - dev, staging, production
- **PR-based** - aprovação via Pull Requests

**Workflow GitOps:**
```yaml
# .sloth/gitops.yaml
repos:
  - name: k8s-manifests
    url: https://github.com/org/k8s-manifests.git
    branch: main
    path: production/
    sync_interval: 5m
    auto_sync: true
    prune: true

hooks:
  on_sync:
    - notify-slack
  on_drift:
    - alert-team
```

**CLI:**
```bash
# Sync manualmente
sloth-runner gitops sync k8s-manifests

# Ver status
sloth-runner gitops status

# Ver drift
sloth-runner gitops diff
```

**Documentação:** `/docs/en/gitops-features.md`

---

### CI/CD Integration

**Descrição:** Integração com pipelines CI/CD.

**Suporte:**
- GitHub Actions
- GitLab CI
- Jenkins
- CircleCI
- Travis CI
- Azure Pipelines

**GitHub Actions example:**
```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Install Sloth Runner
        run: |
          curl -L https://github.com/org/sloth-runner/releases/latest/download/sloth-runner-linux-amd64 -o sloth-runner
          chmod +x sloth-runner

      - name: Run deployment
        env:
          SLOTH_RUNNER_MASTER_ADDR: ${{ secrets.SLOTH_MASTER }}
        run: |
          ./sloth-runner run deploy \
            --file deployments/production.sloth \
            --delegate-to web-01 \
            --yes
```

---

### Sloths

**Descrição:** Repositório de workflows salvos e reutilizáveis.

**Features:**
- **Versionamento** - histórico de versões
- **Tags** - organização por tags
- **Search** - busca por nome/descrição/tags
- **Clone** - clonar sloth existente
- **Export/Import** - compartilhar sloths
- **Active/Inactive** - ativar/desativar sem deletar

**Comandos:**
```bash
# Adicionar sloth
sloth-runner sloth add deploy --file deploy.sloth

# Listar sloths
sloth-runner sloth list

# Ver sloth
sloth-runner sloth get deploy

# Executar sloth
sloth-runner run deploy --file $(sloth-runner sloth get deploy --show-path)

# Remover sloth
sloth-runner sloth remove deploy
```

**Documentação:** `/docs/features/sloth-management.md`

---

## 📊 Monitoramento

### Telemetria

**Descrição:** Sistema completo de observabilidade.

**Componentes:**
- Prometheus metrics
- Structured logging
- Distributed tracing
- Health checks
- Performance profiling

**Arquitetura:**
```
┌──────────┐    metrics    ┌────────────┐
│  Master  ├───────────────► Prometheus │
└──────────┘               └─────┬──────┘
                                 │
┌──────────┐    metrics          │
│ Agent 1  ├─────────────────────┤
└──────────┘                     │
                                 ▼
┌──────────┐    metrics    ┌──────────┐
│ Agent 2  ├───────────────►  Grafana │
└──────────┘               └──────────┘
```

**Endpoints:**
```
http://master:9090/metrics
http://agent:9091/metrics
```

**Documentação:** `/docs/en/telemetry/index.md`

---

### Prometheus Metrics

**Descrição:** Métricas exportadas em formato Prometheus.

**Métricas disponíveis:**

#### Workflows
```
sloth_workflow_executions_total{status="success|failed"}
sloth_workflow_duration_seconds{workflow="name"}
sloth_workflow_tasks_total{workflow="name"}
```

#### Agentes
```
sloth_agent_connected_total
sloth_agent_cpu_usage_percent{agent="name"}
sloth_agent_memory_usage_bytes{agent="name"}
sloth_agent_disk_usage_bytes{agent="name"}
```

#### Sistema
```
sloth_tasks_executed_total
sloth_hooks_triggered_total{event="type"}
sloth_db_size_bytes
```

**Scrape config:**
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'sloth-master'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'sloth-agents'
    static_configs:
      - targets:
        - 'agent1:9091'
        - 'agent2:9091'
```

**Documentação:** `/docs/en/telemetry/prometheus-metrics.md`

---

### Grafana Dashboards

**Descrição:** Dashboards pré-configurados para Grafana.

**Dashboards:**
1. **Overview** - visão geral do sistema
2. **Agents** - métricas de todos os agentes
3. **Workflows** - execuções e performance
4. **Resources** - CPU, memory, disk, network

**Importar dashboard:**
```bash
# Gerar dashboard JSON
sloth-runner agent metrics grafana web-01 --export dashboard.json

# Importar no Grafana
curl -X POST http://admin:admin@localhost:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -d @dashboard.json
```

**Features:**
- Auto-refresh (5s, 10s, 30s, 1m)
- Time range selector
- Variables (agent, workflow)
- Alertas configuráveis
- Export PNG/PDF

**Documentação:** `/docs/en/telemetry/grafana-dashboard.md`

---

### Logs Centralizados

**Descrição:** Sistema centralizado de logs estruturados.

**Features:**
- **Structured** - JSON structured logs
- **Levels** - debug, info, warn, error
- **Context** - metadata rica
- **Search** - busca por qualquer campo
- **Export** - JSON, CSV, texto
- **Retention** - política de retenção

**Formato:**
```json
{
  "timestamp": "2025-10-07T10:30:45Z",
  "level": "info",
  "message": "Workflow completed",
  "workflow": "deploy",
  "agent": "web-01",
  "duration": 45.3,
  "status": "success"
}
```

**Acesso:**
```bash
# CLI
sloth-runner logs --follow

# Web UI
http://localhost:8080/logs

# API
curl http://localhost:8080/api/v1/logs?level=error&since=1h
```

---

### Agent Metrics

**Descrição:** Métricas detalhadas de agentes em tempo real.

**Métricas coletadas:**
- CPU usage (%)
- Memory usage (bytes, %)
- Disk usage (bytes, %)
- Load average (1m, 5m, 15m)
- Network I/O (bytes/sec)
- Process count
- Uptime

**Visualização:**
```bash
# CLI
sloth-runner agent metrics web-01
sloth-runner agent metrics web-01 --watch

# Web UI - Agent Dashboard
http://localhost:8080/agent-dashboard?agent=web-01

# API
curl http://localhost:8080/api/v1/agents/web-01/metrics
```

**Formato:**
```json
{
  "cpu": {
    "cores": 4,
    "usage_percent": 45.2,
    "load_avg": [1.2, 0.8, 0.5]
  },
  "memory": {
    "total_bytes": 8589934592,
    "used_bytes": 4294967296,
    "usage_percent": 50.0
  },
  "disk": {
    "total_bytes": 107374182400,
    "used_bytes": 53687091200,
    "usage_percent": 50.0
  }
}
```

---

## ☁️ Cloud & IaC

### Multi-Cloud

**Descrição:** Suporte nativo para múltiplos provedores cloud.

**Provedores suportados:**
- ✅ AWS (EC2, S3, RDS, Lambda, ECS, EKS, etc.)
- ✅ Azure (VMs, Storage, AKS, Functions, etc.)
- ✅ GCP (Compute Engine, Cloud Storage, GKE, etc.)
- ✅ DigitalOcean (Droplets, Spaces, K8s, etc.)
- ✅ Linode
- ✅ Vultr
- ✅ Hetzner Cloud

**Exemplo multi-cloud:**
```yaml
# Deploy em AWS e GCP simultaneamente
tasks:
  - name: Deploy to AWS
    exec:
      script: |
        aws.ec2_instance_create({
          image_id = "ami-xxx",
          instance_type = "t3.medium"
        })
    delegate_to: aws-agent

  - name: Deploy to GCP
    exec:
      script: |
        gcp.compute_instance_create({
          machine_type = "e2-medium",
          image_family = "ubuntu-2204-lts"
        })
    delegate_to: gcp-agent
```

**Documentação:** `/docs/en/enterprise-features.md`

---

### Terraform

**Descrição:** Integração completa com Terraform.

**Features:**
- `terraform.init` - Inicializar
- `terraform.plan` - Planejar
- `terraform.apply` - Aplicar
- `terraform.destroy` - Destruir
- State management
- Backend config
- Variable files

**Exemplo:**
```lua
local tf_dir = "/infra/terraform"

-- Initialize
terraform.init(tf_dir, {
    backend_config = {
        bucket = "my-tf-state",
        key = "prod/terraform.tfstate"
    }
})

-- Plan
local plan = terraform.plan(tf_dir, {
    var_file = "production.tfvars",
    vars = {
        region = "us-east-1",
        environment = "production"
    }
})

-- Apply se tiver mudanças
if plan.changes > 0 then
    terraform.apply(tf_dir, {
        auto_approve = true
    })
end
```

**Documentação:** `/docs/modules/terraform.md`

---

### Pulumi

**Descrição:** Integração com Pulumi.

**Suporte:**
- Stack management
- Configuration
- Up/Deploy
- Destroy
- Preview

**Exemplo:**
```lua
-- Select stack
pulumi.stack_select("production")

-- Configure
pulumi.config_set("aws:region", "us-east-1")

-- Deploy
pulumi.up({
    yes = true,  -- auto-approve
    parallel = 10
})
```

**Documentação:** `/docs/modules/pulumi.md`

---

### Kubernetes

**Descrição:** Deploy e gerenciamento Kubernetes.

**Features:**
- Apply manifests
- Helm charts
- Namespaces
- ConfigMaps/Secrets
- Rollouts
- Health checks

**Exemplo:**
```lua
-- Apply manifests
kubernetes.apply("/k8s/deployment.yaml", {
    namespace = "production"
})

-- Helm install
helm.install("myapp", "charts/myapp", {
    namespace = "production",
    values = {
        image = {
            tag = "v1.2.3"
        }
    }
})

-- Wait for rollout
kubernetes.rollout_status("deployment/myapp", {
    namespace = "production",
    timeout = "5m"
})
```

**Documentação:** `/docs/en/gitops/kubernetes.md`

---

### Docker

**Descrição:** Automação Docker completa.

**Funcionalidades:**
- Container lifecycle (run, stop, remove)
- Image management (build, push, pull)
- Networks (create, connect)
- Volumes (create, mount)
- Docker Compose

**Exemplo deployment:**
```lua
-- Build image
docker.image_build(".", {
    tag = "myapp:v1.2.3",
    build_args = {
        VERSION = "1.2.3"
    }
})

-- Push to registry
docker.image_push("myapp:v1.2.3", {
    registry = "registry.example.com"
})

-- Deploy
docker.container_run("myapp:v1.2.3", {
    name = "app",
    ports = {"3000:3000"},
    env = {
        DATABASE_URL = "postgres://..."
    },
    restart = "unless-stopped"
})
```

**Documentação:** `/docs/modules/docker.md`

---

## 🔐 Segurança & Enterprise

### Autenticação

**Descrição:** Sistema de autenticação para Web UI e API.

**Métodos:**
- Username/Password
- JWT tokens
- OAuth2 (GitHub, Google, etc.)
- LDAP/AD
- SSO

**Setup:**
```yaml
# config.yaml
auth:
  enabled: true
  type: jwt
  jwt:
    secret: "your-secret-key"
    expiry: 24h
  oauth:
    providers:
      - github:
          client_id: "xxx"
          client_secret: "yyy"
```

---

### TLS/SSL

**Descrição:** Suporte a TLS/SSL para comunicação segura.

**Features:**
- gRPC TLS
- HTTPS Web UI
- Certificate management
- Auto-renewal (Let's Encrypt)

**Configuração:**
```bash
# Master com TLS
sloth-runner server \
  --tls-cert /etc/sloth/cert.pem \
  --tls-key /etc/sloth/key.pem

# Agent com TLS
sloth-runner agent start \
  --master-tls-cert /etc/sloth/master-cert.pem
```

---

### Audit Logs

**Descrição:** Logs de auditoria de todas as ações.

**Eventos auditados:**
- User login/logout
- Workflow execution
- Configuration changes
- API calls
- Admin actions

**Formato:**
```json
{
  "timestamp": "2025-10-07T10:30:45Z",
  "event": "workflow.executed",
  "user": "admin",
  "ip": "192.168.1.100",
  "resource": "deploy.sloth",
  "action": "execute",
  "result": "success"
}
```

---

### Backups

**Descrição:** Sistema automatizado de backups.

**Features:**
- Auto-backup configurável
- Compression (gzip)
- Retention policy
- Remote backup (S3, Azure Blob, etc.)
- Restore

**Comandos:**
```bash
# Backup manual
sloth-runner db backup --output /backup/sloth.db --compress

# Restore
sloth-runner db restore /backup/sloth.db.gz --decompress

# Automated backup (cron)
0 3 * * * sloth-runner db backup --output /backup/sloth-$(date +\%Y\%m\%d).db --compress
```

---

### RBAC

**Descrição:** Role-Based Access Control.

**Roles:**
- **Admin** - acesso total
- **Operator** - executar workflows, gerenciar agentes
- **Developer** - criar/editar workflows
- **Viewer** - apenas visualização

**Permissions:**
```yaml
roles:
  operator:
    permissions:
      - workflow:execute
      - agent:view
      - agent:restart
      - logs:view

  developer:
    permissions:
      - workflow:create
      - workflow:edit
      - workflow:execute
      - logs:view

  viewer:
    permissions:
      - workflow:view
      - agent:view
      - logs:view
```

---

## 🚀 Performance

### Otimizações

**Descrição:** Otimizações recentes de performance.

**Melhorias implementadas:**

#### Agent Optimizations
- ✅ **Ultra-low memory** - 32MB RAM footprint
- ✅ **Binary size reduction** - de 45MB → 12MB
- ✅ **Startup time** - <100ms
- ✅ **CPU efficiency** - 99% idle quando inativo

#### Database Optimizations
- ✅ **WAL mode** - Write-Ahead Logging
- ✅ **Connection pooling** - reuso de conexões
- ✅ **Prepared statements** - queries otimizadas
- ✅ **Indexes** - índices em campos críticos
- ✅ **Auto-vacuum** - limpeza automática

#### gRPC Optimizations
- ✅ **Connection reuse** - keepalive
- ✅ **Compression** - gzip compression
- ✅ **Multiplexing** - múltiplas streams
- ✅ **Buffer pooling** - reuso de buffers

**Benchmark:**
```
Antes:
- Agent memory: 128MB
- Binary size: 45MB
- Startup time: 2s

Depois:
- Agent memory: 32MB (75% redução)
- Binary size: 12MB (73% redução)
- Startup time: 95ms (95% mais rápido)
```

**Documentação:** `/docs/PERFORMANCE_OPTIMIZATIONS.md`

---

### Parallel Execution

**Descrição:** Execução paralela de tasks usando goroutines.

**Features:**
- **goroutine.parallel()** - executa funções em paralelo
- **Concurrency control** - limite de goroutines simultâneas
- **Error handling** - coleta erros de todas as goroutines
- **Wait groups** - sincronização automática

**Exemplo:**
```lua
-- Executar múltiplas tasks em paralelo
goroutine.parallel({
    function()
        pkg.install("nginx")
    end,
    function()
        pkg.install("postgresql")
    end,
    function()
        pkg.install("redis")
    end
})

-- Com limite de concorrência
goroutine.parallel({
    tasks = {
        function() exec.command("task1") end,
        function() exec.command("task2") end,
        function() exec.command("task3") end,
        function() exec.command("task4") end
    },
    max_concurrent = 2  -- Máximo 2 por vez
})
```

**Documentação:** `/docs/modules/goroutine.md`

---

### Resource Limits

**Descrição:** Limites de recursos configuráveis.

**Configuração:**
```yaml
# Agent config
resources:
  cpu:
    limit: 2  # cores
    reserve: 0.5
  memory:
    limit: 2GB
    reserve: 512MB
  disk:
    limit: 10GB
    min_free: 1GB
```

**Enforcement:**
- CPU throttling
- Memory limits (cgroup)
- Disk quota
- Task timeout

---

### Caching

**Descrição:** Sistema de cache para otimização.

**Tipos de cache:**

#### Module cache
- Módulos Lua compilados
- Reduce load time

#### State cache
- State em memória
- Reduce DB queries

#### Metrics cache
- Métricas agregadas
- Reduce computation

**Configuração:**
```yaml
cache:
  enabled: true
  ttl: 5m
  max_size: 100MB
  eviction: lru  # least recently used
```

---

## 📚 Recursos Adicionais

### Documentação
- [🚀 Quick Start](/docs/en/quick-start.md)
- [🏗️ Arquitetura](/docs/architecture/sloth-runner-architecture.md)
- [📖 Modern DSL](/docs/modern-dsl/introduction.md)
- [🎯 Exemplos Avançados](/docs/en/advanced-examples.md)

### Links Úteis
- [GitHub Repository](https://github.com/chalkan3/sloth-runner)
- [Issue Tracker](https://github.com/chalkan3/sloth-runner/issues)
- [Releases](https://github.com/chalkan3/sloth-runner/releases)

---

**Última atualização:** 2025-10-07

**Total de Funcionalidades Documentadas:** 100+
