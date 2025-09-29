# Sloth Runner - Web UI Dashboard

## 🎨 Interface Web Embarcada

O Sloth Runner agora possui uma interface web moderna e responsiva para gerenciar tarefas, agentes e monitorar o sistema de forma visual e intuitiva.

## 🚀 Como Usar

### Iniciando a Interface Web

```bash
# Executar em modo normal
sloth-runner ui

# Executar em porta específica
sloth-runner ui --port 8080

# Executar como daemon (em background)
sloth-runner ui --daemon

# Executar com debug habilitado
sloth-runner ui --debug
```

### Acessando a Interface

Após iniciar o servidor, acesse no seu navegador:
```
http://localhost:8080
```

## 🌟 Funcionalidades

### Dashboard Principal
- **Estatísticas em Tempo Real**: Visualize métricas como total de tarefas, agentes ativos, tarefas em execução
- **Cards Informativos**: Interface moderna com design responsivo
- **Tema Escuro/Claro**: Alternância de tema com um clique

### Gerenciamento de Tarefas
- **Criação de Tarefas**: Interface intuitiva para criar novas tarefas
- **Tipos Suportados**:
  - Shell Command
  - Lua Script  
  - Pipeline
- **Status em Tempo Real**: Acompanhe o progresso das tarefas
- **Console Integrado**: Visualize logs e outputs em tempo real

### Gerenciamento de Agentes
- **Lista de Agentes**: Visualize todos os agentes registrados
- **Status Online/Offline**: Indicadores visuais do status dos agentes
- **Adicionar Novos Agentes**: Interface para registrar novos agentes
- **Execução Remota**: Execute comandos em agentes específicos

### Recursos Avançados
- **WebSocket**: Atualizações em tempo real
- **Design Responsivo**: Funciona em desktop, tablet e mobile
- **Interface Moderna**: Design glassmorphism com gradientes
- **Notificações**: Sistema de notificações para feedback do usuário

## 🎯 Recursos da Interface

### Painel de Controle
```
┌─────────────────────────────────────────────────────┐
│ 🦥 Sloth Runner - Dashboard                          │
├─────────────────────────────────────────────────────┤
│ [📊 Stats] [🔄 Refresh] [➕ New Task] [🛑 Stop All]  │
├─────────────────────────────────────────────────────┤
│                                                     │
│  📈 Estatísticas:                                   │
│  ┌─────────┬─────────┬─────────┬─────────┐          │
│  │Tasks: 12│Agents: 3│Running:2│Done: 8  │          │
│  └─────────┴─────────┴─────────┴─────────┘          │
│                                                     │
│  📋 Tarefas Recentes    |  🖥️  Agentes             │
│  ┌─────────────────────┬─────────────────────┐      │
│  │ ✅ Build Project    │ 🟢 local-agent      │      │
│  │ 🔄 Run Tests        │ 🟢 prod-server-1    │      │
│  │ ⏳ Deploy Staging   │ 🔴 test-server      │      │
│  │ ❌ Backup DB        │                     │      │
│  └─────────────────────┴─────────────────────┘      │
│                                                     │
│  🖥️  Console:                                       │
│  ┌─────────────────────────────────────────────────┐│
│  │ $ [16:00:01] Task "Build Project" completed     ││
│  │ $ [16:00:15] WebSocket connection established   ││
│  │ $ [16:00:30] Agent "prod-server-1" online       ││
│  └─────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────┘
```

### Modais Interativos

#### Criar Nova Tarefa
- Nome da tarefa
- Tipo (Shell/Lua/Pipeline)
- Comando/Script
- Agente de destino
- Execução imediata

#### Adicionar Agente
- Nome do agente
- Endereço (host:porta)
- Configurações TLS

## 🛠️ Tecnologias

### Backend
- **Go**: Server HTTP embarcado
- **Gorilla Mux**: Roteamento HTTP
- **Gorilla WebSocket**: Comunicação em tempo real
- **Embed**: Arquivos estáticos embarcados

### Frontend
- **HTML5/CSS3/JavaScript**: Interface pura sem frameworks
- **Font Awesome**: Ícones
- **WebSocket API**: Comunicação bidirecional
- **Fetch API**: Requisições REST

## 🔧 Arquitetura

```
┌─────────────────┐    HTTP/WS     ┌─────────────────┐
│   Web Browser   │ ◄─────────────► │   UI Server     │
└─────────────────┘                └─────────────────┘
                                           │
                                           │ Internal API
                                           ▼
                                   ┌─────────────────┐
                                   │ Task Manager    │
                                   │ Agent Registry  │
                                   │ Lua Interface   │
                                   └─────────────────┘
```

## 📱 Screenshots

### Dashboard Principal
![Dashboard](docs/ui-dashboard.png)

### Criação de Tarefas
![New Task](docs/ui-new-task.png)

### Tema Escuro
![Dark Mode](docs/ui-dark-mode.png)

## 🎨 Personalização

### Tema Personalizado
O CSS pode ser facilmente customizado. Os arquivos estáticos estão em:
```
internal/ui/static/
├── index.html
└── enhanced.css
```

### Branding
Para personalizar o branding:
1. Edite o título em `index.html`
2. Modifique as cores em `enhanced.css`
3. Substitua ícones e logos

## 🔒 Segurança

- **CORS**: Configurado para desenvolvimento
- **WebSocket Origin**: Validação de origem
- **Static Files**: Servidos com segurança
- **No External Dependencies**: Tudo embarcado

## 🚧 Desenvolvimento

### Adicionando Novas Funcionalidades

1. **Novas APIs**: Adicione rotas em `internal/ui/server.go`
2. **Frontend**: Modifique `static/index.html`
3. **Estilos**: Atualize `static/enhanced.css`
4. **WebSocket**: Implemente novos tipos de mensagem

### Build e Deploy

```bash
# Compilar com UI embarcada
go build -o sloth-runner ./cmd/sloth-runner

# O binário contém todos os assets necessários
./sloth-runner ui
```

## 🎯 Próximas Funcionalidades

- [ ] Editor de Lua integrado
- [ ] Visualização de logs avançada
- [ ] Métricas e gráficos
- [ ] Autenticação/autorização
- [ ] API de configuração
- [ ] Templates de tarefas
- [ ] Scheduler visual
- [ ] Dashboard customizável

## 📋 Comandos Disponíveis

```bash
# Ajuda da UI
sloth-runner ui --help

# Executar UI
sloth-runner ui

# UI em porta específica
sloth-runner ui --port 9090

# UI como daemon
sloth-runner ui --daemon

# UI com debug
sloth-runner ui --debug

# Parar UI daemon
pkill -f "sloth-runner ui"
```

## 🎉 Conclusão

A interface web do Sloth Runner oferece uma experiência visual moderna e intuitiva para gerenciar todas as funcionalidades do sistema. Com design responsivo, atualizações em tempo real e interface amigável, facilita o uso tanto para desenvolvimento quanto para produção.

Acesse `http://localhost:8080` e explore todas as funcionalidades!