# 🎨 Guia Completo da Web UI

## Visão Geral

A Web UI do Sloth Runner é uma interface moderna, responsiva e intuitiva para gerenciar workflows, agentes, hooks, eventos e monitoramento. Construída com Bootstrap 5, Chart.js e WebSockets para atualizações em tempo real.

**URL de Acesso:** `http://localhost:8080` (porta padrão)

---

## 🚀 Iniciar a Web UI

```bash
# Método 1: Comando ui
sloth-runner ui --port 8080

# Método 2: Com bind específico
sloth-runner ui --port 8080 --bind 0.0.0.0

# Método 3: Com variável de ambiente
export SLOTH_RUNNER_UI_PORT=8080
sloth-runner ui
```

---

## 📱 Páginas e Funcionalidades

### 1. 🏠 Dashboard Principal (`/`)

**Funcionalidades:**

- **Visão geral do sistema** - Cards com estatísticas principais
  - Total de workflows
  - Total de agentes ativos/inativos
  - Execuções recentes
  - Taxa de sucesso

- **Gráficos interativos** (Chart.js)
  - Execuções por dia (últimos 7 dias)
  - Taxa de sucesso vs falha
  - Uso de recursos dos agentes
  - Distribuição de workflows

- **Feed de atividades em tempo real**
  - Workflows iniciados/completados
  - Agentes conectados/desconectados
  - Eventos do sistema
  - Atualização via WebSocket

- **Quick Actions** (botão flutuante)
  - Executar workflow
  - Criar novo workflow
  - Ver agentes
  - Configurações

**Recursos modernos:**
- 🎨 Modo escuro/claro (toggle automático)
- 📊 Gráficos responsivos
- 🔄 Auto-refresh a cada 30 segundos
- 🎯 Animações suaves (fade-in, hover effects)
- 📱 Design mobile-first

---

### 2. 🤖 Gerenciamento de Agentes (`/agents`)

**Funcionalidades:**

#### Cards de Agentes

Cada agente é exibido em um card moderno com:

- **Status visual**
  - 🟢 Online (verde com pulse animation)
  - 🔴 Offline (cinza)
  - 🟡 Connecting (amarelo)

- **Métricas em tempo real**
  - CPU Usage (%) - gráfico de progresso animado
  - Memory Usage (%) - gráfico de progresso animado
  - Disk Usage (%) - gráfico de progresso animado
  - Load Average - convertido para % baseado em CPU cores

- **Informações do agente**
  - Nome e endereço
  - Versão do agente
  - Uptime formatado (d/h/m/s)
  - Data de registro
  - Último heartbeat

- **Sparklines** (mini gráficos de tendência)
  - CPU usage nas últimas 24h
  - Memory usage nas últimas 24h
  - Renderizados com Canvas API

- **Botões de ação**
  - 📊 Dashboard - vai para dashboard do agente
  - ℹ️ Details - modal com detalhes completos
  - 📄 Logs - visualizar logs do agente
  - 🔄 Restart - reiniciar agente (apenas se online)

#### Estatísticas Gerais

- Total de agentes
- Agentes ativos
- Agentes inativos
- Taxa de uptime (%)

#### Funcionalidades Avançadas

- **Auto-refresh** - atualiza lista a cada 10 segundos
- **WebSocket** - notificações em tempo real quando agentes conectam/desconectam
- **Filtros** - filtrar por status (todos/ativos/inativos)
- **Busca** - buscar agentes por nome
- **Grid responsivo** - cards se reorganizam automaticamente
- **Skeleton loaders** - loading states profissionais

**Layout:**

```
┌─────────────────────────────────────────┐
│  📊 Stats Cards                         │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐  │
│  │Total │ │Active│ │Inact.│ │Uptime│  │
│  └──────┘ └──────┘ └──────┘ └──────┘  │
├─────────────────────────────────────────┤
│  🤖 Agent Cards                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │ Agent 1 │ │ Agent 2 │ │ Agent 3 │  │
│  │ 🟢 80%  │ │ 🟢 45%  │ │ 🔴 N/A  │  │
│  │ [graph] │ │ [graph] │ │ [graph] │  │
│  │ [btns]  │ │ [btns]  │ │ [btns]  │  │
│  └─────────┘ └─────────┘ └─────────┘  │
└─────────────────────────────────────────┘
```

---

### 3. 🎛️ Controle de Agentes (`/agent-control`)

**Funcionalidades:**

Página dedicada para controle detalhado de cada agente.

- **Lista de agentes** com cards expandidos
- **Métricas detalhadas**
  - CPU cores, load average
  - Memory total/used/free
  - Disk total/used/free
  - Network interfaces
  - Uptime detalhado

- **Ações de controle**
  - Start/Stop/Restart agent
  - Update agent version
  - Check modules
  - Run command remotely
  - View logs

- **Gauge charts** (gráficos circulares)
  - CPU usage
  - Memory usage
  - Disk usage
  - Com cores dinâmicas baseadas em thresholds

**Thresholds de cores:**
- 🟢 Verde: 0-40%
- 🔵 Azul: 40-70%
- 🟡 Amarelo: 70-90%
- 🔴 Vermelho: 90-100%

---

### 4. 📊 Dashboard do Agente (`/agent-dashboard`)

**Funcionalidades:**

Dashboard individual para cada agente com métricas avançadas.

- **Time-series charts** (Chart.js)
  - CPU usage over time (line chart)
  - Memory usage over time (area chart)
  - Disk I/O (bar chart)
  - Network traffic (line chart)

- **System information**
  - OS name, version, kernel
  - CPU model, cores, architecture
  - Total memory, swap
  - Filesystems montados

- **Process list**
  - Top processes por CPU
  - Top processes por Memory
  - Atualização em tempo real

- **Logs em tempo real**
  - Stream de logs do agente
  - Filtros por nível (debug/info/warn/error)
  - Auto-scroll
  - Download de logs

- **Time range selector**
  - Last 1 hour
  - Last 6 hours
  - Last 24 hours
  - Last 7 days
  - Custom range

---

### 5. 📝 Workflows (`/workflows`)

**Funcionalidades:**

#### Lista de Workflows

- **Cards de workflows** com:
  - Nome e descrição
  - Tags/labels
  - Última execução
  - Taxa de sucesso
  - Botões: Run, Edit, Delete

- **Filtros**
  - Por tags
  - Por status (ativo/inativo)
  - Por frequência de execução

- **Busca** - buscar por nome/descrição

#### Editor de Workflow

**Funcionalidades do editor:**

- **Code Editor profissional**
  - Syntax highlighting para YAML/Sloth DSL
  - Line numbers
  - Auto-indent
  - Bracket matching
  - Keywords: `tasks`, `name`, `exec`, `delegate_to`, etc.
  - Cores customizadas (Sloth theme)

- **Keyboard shortcuts**
  - `Tab` - indentação (2 espaços)
  - `Shift+Tab` - des-indentação
  - `Ctrl+S` / `Cmd+S` - salvar
  - `Shift+Alt+F` - formatar

- **Templates**
  - Basic workflow
  - Multi-task workflow
  - Distributed workflow (com delegate_to)
  - Docker deployment
  - Full example workflow

- **Validação em tempo real**
  - Syntax checking
  - Linting
  - Avisos e erros inline

- **Preview**
  - Visualizar estrutura do workflow
  - Dependências entre tasks
  - Variáveis utilizadas

**Exemplo de sintaxe highlighting:**

```yaml
# Keywords em roxo
tasks:
  - name: Deploy App          # Strings em verde
    exec:
      script: |                # Pipe em laranja
        pkg.install("nginx")   # Funções em azul
        # Comments em cinza
    delegate_to: web-01        # Keys em azul claro
```

---

### 6. ⚡ Execuções (`/executions`)

**Funcionalidades:**

Histórico completo de execuções de workflows.

- **Lista de execuções** com:
  - Workflow name
  - Status (success/failed/running)
  - Started/completed time
  - Duration
  - Triggered by (user/schedule/hook)
  - Agent name (se delegado)

- **Filtros avançados**
  - Por status
  - Por workflow
  - Por agente
  - Por data/hora
  - Por usuário

- **Detalhes da execução**
  - Output completo
  - Logs estruturados
  - Task-by-task breakdown
  - Variáveis utilizadas
  - Métricas de performance

- **Ações**
  - Re-run workflow
  - Download logs
  - Share execution (link)
  - Compare with previous

- **Status indicators**
  - 🟢 Success (verde)
  - 🔴 Failed (vermelho)
  - 🟡 Running (amarelo com spinner)
  - ⏸️ Paused (cinza)

---

### 7. 🎣 Hooks (`/hooks`)

**Funcionalidades:**

Gerenciamento completo de hooks (event handlers).

- **Lista de hooks**
  - Nome do hook
  - Event type
  - Script path
  - Priority
  - Enabled/disabled status
  - Last triggered

- **Tipos de eventos**
  - `workflow.started`
  - `workflow.completed`
  - `workflow.failed`
  - `task.started`
  - `task.completed`
  - `task.failed`
  - `agent.connected`
  - `agent.disconnected`

- **Criar/Editar hook**
  - Form com campos:
    - Name
    - Event type (dropdown)
    - Script (code editor)
    - Priority (0-100)
    - Enabled (toggle)
  - Syntax highlighting para Lua/Bash
  - Validação

- **Testar hook**
  - Dry-run com payload de teste
  - Ver output sem executar ações
  - Debug mode

- **Histórico de hooks**
  - Quando foi disparado
  - Payload recebido
  - Output do script
  - Success/failure

---

### 8. 📡 Eventos (`/events`)

**Funcionalidades:**

Monitoramento de eventos do sistema em tempo real.

- **Feed de eventos**
  - Timestamp
  - Event type
  - Source (workflow/agent/system)
  - Details/payload
  - Status

- **WebSocket stream**
  - Eventos aparecem em tempo real
  - Notificações sonoras (opcional)
  - Desktop notifications (opcional)

- **Filtros**
  - Por tipo de evento
  - Por source
  - Por status
  - Por time range

- **Exportar eventos**
  - JSON
  - CSV
  - Logs format

- **Estatísticas**
  - Eventos por hora
  - Eventos por tipo
  - Top sources

---

### 9. 📅 Scheduler (`/scheduler`)

**Funcionalidades:**

Agendamento de workflows (cron-like).

- **Jobs agendados**
  - Nome do job
  - Workflow associado
  - Cron expression
  - Next run time
  - Last run status
  - Enabled/disabled

- **Criar job**
  - Form com:
    - Name
    - Workflow (dropdown)
    - Schedule (cron ou visual builder)
    - Variables
    - Notifications (on success/failure)

- **Visual cron builder**
  - Seletor de minutos/horas/dias/meses
  - Preview: "Runs every day at 3:00 AM"
  - Templates comuns:
    - Every hour
    - Every day at midnight
    - Every Monday at 9 AM
    - Custom

- **Histórico de execuções**
  - Por job agendado
  - Success rate
  - Average duration

---

### 10. 📄 Logs (`/logs`)

**Funcionalidades:**

Visualização centralizada de logs.

- **Filtros avançados**
  - Por nível (debug/info/warn/error)
  - Por source (agent/workflow/system)
  - Por time range
  - Por texto (search)

- **Live tail**
  - Stream em tempo real
  - Auto-scroll
  - Pause/resume
  - Highlight patterns

- **Structured logs**
  - JSON format
  - Campos expandíveis
  - Syntax highlighting

- **Exportar**
  - Download como .log
  - Copy to clipboard
  - Share link

- **Log levels com cores**
  - 🔵 DEBUG (azul)
  - 🟢 INFO (verde)
  - 🟡 WARN (amarelo)
  - 🔴 ERROR (vermelho)

---

### 11. 🖥️ Terminal Interativo (`/terminal`)

**Funcionalidades:**

Terminal web interativo para agentes remotos.

- **xterm.js** - terminal completo
- **SSH-like experience**
- **Múltiplas sessões** (tabs)
- **Command history** (setas ↑↓)
- **Auto-complete** (Tab)
- **Copy/paste** (Ctrl+Shift+C/V)
- **Themes** - Solarized, Monokai, Dracula, etc.

**Comandos especiais:**
- `.clear` - limpar terminal
- `.exit` - fechar sessão
- `.upload <file>` - upload arquivo
- `.download <file>` - download arquivo

---

### 12. 📦 Sloths Salvos (`/sloths`)

**Funcionalidades:**

Repositório de workflows salvos.

- **Lista de sloths**
  - Nome
  - Description
  - Tags
  - Created/updated date
  - Active/inactive status
  - Use count

- **Ações**
  - Run sloth
  - Edit content
  - Clone sloth
  - Export (YAML)
  - Delete
  - Activate/Deactivate

- **Tags management**
  - Criar tags
  - Colorir tags
  - Filtrar por tags

- **Versionamento**
  - Ver histórico de versões
  - Comparar versões (diff)
  - Restaurar versão anterior

---

### 13. ⚙️ Configurações (`/settings`)

**Funcionalidades:**

#### General Settings

- Server info
  - Master address
  - gRPC port
  - Web UI port
  - Database path

- Preferences
  - Theme (light/dark/auto)
  - Language (en/pt/zh)
  - Timezone
  - Date format

#### Notifications

- Email settings
  - SMTP host, port
  - Username/password
  - From address

- Slack integration
  - Webhook URL
  - Default channel
  - Mentions

- Telegram/Discord
  - Bot token
  - Chat ID / Webhook

#### Security

- Authentication
  - Enable/disable auth
  - Users management
  - API tokens

- TLS/SSL
  - Enable HTTPS
  - Certificate upload
  - Auto-renewal (Let's Encrypt)

#### Database

- Backup settings
  - Auto-backup enabled
  - Backup schedule
  - Retention policy

- Maintenance
  - Vacuum database
  - View stats
  - Clear old data

---

## 🎨 Recursos Visuais Modernos

### Dark Mode / Light Mode

**Auto-detection** baseado em preferência do sistema + toggle manual

**Temas:**

```css
/* Light Mode */
--bg-primary: #ffffff
--text-primary: #212529
--accent: #4F46E5

/* Dark Mode */
--bg-primary: #1a1d23
--text-primary: #e9ecef
--accent: #818CF8
```

**Toggle:** Botão no navbar com ícones ☀️ (light) / 🌙 (dark)

---

### Animações e Transições

- **Fade-in** ao carregar páginas
- **Hover effects** em cards e botões
- **Pulse animation** em status indicators
- **Skeleton loaders** durante loading
- **Smooth scrolling**
- **Ripple effect** em botões (Material Design)
- **Page transitions** suaves

---

### Glassmorphism

Cards com efeito de vidro fosco:

```css
.glass-card {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}
```

---

### Toasts / Notifications

Sistema de notificações moderno:

- **Tipos:**
  - ℹ️ Info (azul)
  - ✅ Success (verde)
  - ⚠️ Warning (amarelo)
  - ❌ Error (vermelho)
  - ⏳ Loading (com spinner)

- **Posições:**
  - Top-right (padrão)
  - Top-left
  - Bottom-right
  - Bottom-left
  - Center

- **Features:**
  - Auto-dismiss (configurável)
  - Close button
  - Progress bar
  - Action buttons
  - Stacking múltiplos toasts

---

### Confetti Effects

Efeitos de confetti em eventos especiais:

- ✅ Workflow completado com sucesso
- 🤖 Novo agente conectado
- 🎯 Milestone alcançado
- 🎉 Deploy completado

```javascript
confetti.burst({
    particleCount: 100,
    spread: 70,
    origin: { y: 0.6 }
});
```

---

### Drag & Drop

- **Reordenar tasks** em workflows
- **Upload de arquivos** (drop zone)
- **Reorganizar dashboard** widgets

---

## ⌨️ Command Palette (Ctrl+Shift+P)

Quick access a todas as ações:

```
🔍 Search commands...

> Run Workflow
> View Agents
> Create Workflow
> Open Terminal
> View Logs
> Settings
> Toggle Dark Mode
> Export Data
...
```

**Features:**
- Fuzzy search
- Keyboard navigation (↑↓ Enter)
- Recent commands
- Shortcuts visíveis
- Context-aware (mostra ações baseadas na página atual)

---

## 📊 Gráficos e Visualizações

### Chart.js Components

**Tipos de gráficos:**

1. **Line Charts** - métricas temporais
2. **Bar Charts** - comparações
3. **Doughnut Charts** - distribuições
4. **Area Charts** - tendências
5. **Sparklines** - mini gráficos inline

**Features:**
- Responsivos
- Tooltips interativos
- Legendas
- Zoom/pan
- Exportar como PNG
- Temas dark/light

---

## 🔄 WebSocket Real-Time Updates

Conexão WebSocket para atualizações em tempo real:

**Eventos em tempo real:**
- Agent connected/disconnected
- Workflow started/completed
- New logs
- System alerts
- Metrics updates

**URL:** `ws://localhost:8080/ws`

**Reconexão automática** se conexão cair

---

## 📱 Mobile Responsive

Design mobile-first:

- **Breakpoints:**
  - Mobile: < 768px
  - Tablet: 768px - 1024px
  - Desktop: > 1024px

- **Mobile features:**
  - Hamburger menu
  - Touch-friendly buttons
  - Swipe gestures
  - Simplified charts
  - Bottom navigation

---

## 🔐 Autenticação (Opcional)

**Login page** se auth estiver habilitado:

- Username/password
- Remember me
- Forgot password
- OAuth (GitHub, Google, etc.)

**JWT tokens** para API

**Roles:**
- Admin - acesso total
- Operator - executar workflows
- Viewer - apenas visualizar

---

## 🛠️ Developer Tools

### API Explorer

Explorar e testar API REST:

```
GET  /api/v1/agents
GET  /api/v1/agents/:name
POST /api/v1/workflows/run
GET  /api/v1/executions
...
```

**Features:**
- Try it out (executar no browser)
- Request/response examples
- Authentication headers
- cURL examples

---

### Logs Browser

Navegar logs do sistema:

- Server logs
- Agent logs
- Application logs
- Audit logs

---

## 🎓 Guias Rápidos

### Quick Start Tour

Tour interativo para novos usuários:

1. Welcome → Agents page
2. Create your first workflow
3. Run a workflow
4. View execution results
5. Set up notifications

**Features:**
- Tooltips com dicas
- Highlight elements
- Skip/Next buttons
- Don't show again (cookie)

---

## 💡 Dicas de Uso

### Atalhos de Teclado

```
Global:
Ctrl+Shift+P  - Command palette
Ctrl+K        - Search
Ctrl+/        - Help
Esc           - Close modals

Editor:
Ctrl+S        - Save
Ctrl+F        - Find
Ctrl+Z        - Undo
Ctrl+Y        - Redo
```

---

### Bookmarklets

Salvar páginas importantes:

```
Dashboard:          /
My Workflows:       /workflows
Active Executions:  /executions?status=running
Agent Metrics:      /agent-dashboard
```

---

### Browser Extensions

**Disponíveis:**
- Chrome Extension - quick access
- Firefox Add-on

---

## 🔧 Customização

### Custom CSS

Adicionar CSS customizado em Settings:

```css
/* Custom theme */
:root {
    --primary-color: #FF6B6B;
    --accent-color: #4ECDC4;
}
```

---

### Widgets Personalizados

Criar widgets customizados para dashboard:

- Custom charts
- External integrations
- Iframe embeds

---

## 📚 Próximos Passos

- [📋 Referência CLI](referencia-cli.md)
- [🔧 Módulos](modulos-completos.md)
- [🎯 Exemplos](../en/advanced-examples.md)
- [🏗️ Arquitetura](../architecture/sloth-runner-architecture.md)

---

## 🐛 Troubleshooting

### Web UI não carrega

```bash
# Verificar se servidor está rodando
lsof -i :8080

# Ver logs
sloth-runner ui --port 8080 --verbose

# Limpar cache do browser
Ctrl+Shift+Del
```

### WebSocket não conecta

```bash
# Verificar firewall
sudo ufw allow 8080

# Testar conexão
wscat -c ws://localhost:8080/ws
```

### Métricas não aparecem

```bash
# Verificar se agente está enviando métricas
sloth-runner agent metrics <agent-name>

# Reiniciar agente
sloth-runner agent restart <agent-name>
```

---

**Última atualização:** 2025-10-07

**Desenvolvido com:** Bootstrap 5, Chart.js, xterm.js, WebSockets, Canvas API
