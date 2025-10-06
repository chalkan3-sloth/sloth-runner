# Sloth Runner Web UI

Uma interface web completa para gerenciar e monitorar o Sloth Runner.

## Funcionalidades

### 🎯 Dashboard
- Visão geral em tempo real do sistema
- Estatísticas de agentes, workflows, hooks e eventos
- Feed de atividades em tempo real
- Monitoramento de agentes ativos
- Tabela de eventos recentes

### 🖥️ Gerenciamento de Agentes
- Listar todos os agentes distribuídos
- Visualizar detalhes completos de cada agente
- Monitorar status e heartbeat em tempo real
- Remover agentes inativos
- Visualizar informações do sistema de cada agente

### 📋 Workflows (Sloths)
- Listar todos os workflows registrados
- Criar novos workflows via interface
- Editar workflows existentes
- Ativar/desativar workflows
- Visualizar estatísticas de uso
- Executar workflows (integração futura)

### ⚡ Hooks
- Gerenciar hooks de eventos
- Habilitar/desabilitar hooks
- Visualizar histórico de execuções
- Monitorar taxa de sucesso/falha
- Ver logs de execução detalhados

### 📬 Eventos
- Monitorar fila de eventos em tempo real
- Visualizar eventos pendentes e processados
- Retentar eventos que falharam
- Filtrar eventos por tipo e status
- Visualizar histórico completo

### 🔐 Secrets
- Visualizar quais secrets existem (nomes apenas)
- Integração segura com CLI
- Documentação de comandos CLI

### 🔌 SSH Profiles
- Gerenciar perfis de conexão SSH
- Visualizar logs de auditoria
- Monitorar uso de perfis
- Estatísticas de conexão

## Como Usar

### Iniciar a UI

```bash
# Porta padrão (8080)
sloth-runner ui

# Porta personalizada
sloth-runner ui --port 3000

# Com autenticação
sloth-runner ui --auth --username admin --password mysecret

# Modo debug
sloth-runner ui --debug
```

### Acessar

Abra o navegador em: `http://localhost:8080`

## Arquitetura

### Backend
- **Framework**: Gin (Go)
- **WebSocket**: gorilla/websocket para atualizações em tempo real
- **Banco de dados**: SQLite (integra com os bancos existentes)
- **API**: RESTful JSON API

### Frontend
- **HTML5/CSS3/JavaScript puro** (sem frameworks pesados)
- **Bootstrap 5** para UI responsiva
- **WebSocket client** para atualizações em tempo real
- **Chart.js** para gráficos (futuro)

### Dados Persistentes
Todos os dados vêm diretamente dos bancos SQLite existentes:
- `.sloth-cache/agents.db` - Agentes
- `/etc/sloth-runner/sloths.db` - Workflows
- `.sloth-cache/hooks.db` - Hooks e eventos
- `~/.sloth-runner/secrets.db` - Secrets (somente nomes)
- `~/.sloth-runner/ssh_profiles.db` - Perfis SSH

## Endpoints da API

### Dashboard
- `GET /api/v1/dashboard` - Estatísticas gerais

### Agents
- `GET /api/v1/agents` - Listar agentes
- `GET /api/v1/agents/:name` - Detalhes do agente
- `DELETE /api/v1/agents/:name` - Remover agente

### Workflows
- `GET /api/v1/sloths` - Listar workflows
- `GET /api/v1/sloths/:name` - Detalhes do workflow
- `POST /api/v1/sloths` - Criar workflow
- `PUT /api/v1/sloths/:name` - Atualizar workflow
- `DELETE /api/v1/sloths/:name` - Deletar workflow
- `POST /api/v1/sloths/:name/activate` - Ativar
- `POST /api/v1/sloths/:name/deactivate` - Desativar

### Hooks
- `GET /api/v1/hooks` - Listar hooks
- `GET /api/v1/hooks/:id` - Detalhes do hook
- `POST /api/v1/hooks` - Criar hook
- `PUT /api/v1/hooks/:id` - Atualizar hook
- `DELETE /api/v1/hooks/:id` - Deletar hook
- `POST /api/v1/hooks/:id/enable` - Habilitar
- `POST /api/v1/hooks/:id/disable` - Desabilitar
- `GET /api/v1/hooks/:id/history` - Histórico de execuções

### Events
- `GET /api/v1/events` - Listar eventos
- `GET /api/v1/events/pending` - Eventos pendentes
- `GET /api/v1/events/:id` - Detalhes do evento
- `POST /api/v1/events/:id/retry` - Retentar evento

### Secrets
- `GET /api/v1/secrets/:stack` - Listar secrets (nomes apenas)

### SSH
- `GET /api/v1/ssh` - Listar perfis SSH
- `GET /api/v1/ssh/:name` - Detalhes do perfil
- `POST /api/v1/ssh` - Criar perfil
- `PUT /api/v1/ssh/:name` - Atualizar perfil
- `DELETE /api/v1/ssh/:name` - Deletar perfil
- `GET /api/v1/ssh/:name/audit` - Logs de auditoria

### WebSocket
- `WS /api/v1/ws` - Conexão WebSocket para atualizações em tempo real

## Mensagens WebSocket

A UI recebe atualizações em tempo real através de WebSocket:

```json
{
  "type": "agent_update",
  "timestamp": 1234567890,
  "data": {
    "name": "agent-1",
    "status": "active"
  }
}
```

Tipos de mensagens:
- `agent_update` - Atualização de status de agente
- `event_update` - Novo evento ou mudança de status
- `hook_execution` - Execução de hook
- `workflow_update` - Mudança em workflow
- `system_alert` - Alerta do sistema

## Segurança

### Autenticação
Suporta HTTP Basic Authentication:
```bash
sloth-runner ui --auth --username admin --password mysecret
```

### CORS
CORS habilitado por padrão para desenvolvimento.

### Secrets
Secrets são **somente leitura** na UI por segurança. Para gerenciar secrets, use a CLI:
```bash
sloth-runner secret add <stack> <name> <value>
```

## Desenvolvimento

### Estrutura de Arquivos
```
internal/webui/
├── server.go                    # Servidor principal
├── handlers/                    # Handlers da API
│   ├── agent.go
│   ├── dashboard.go
│   ├── event.go
│   ├── hook.go
│   ├── secret.go
│   ├── sloth.go
│   ├── ssh.go
│   ├── websocket.go
│   └── wrappers.go             # Wrappers de BD
├── middleware/                  # Middlewares
│   ├── auth.go
│   └── cors.go
├── static/                      # Arquivos estáticos
│   ├── css/
│   │   └── main.css
│   └── js/
│       ├── websocket.js
│       ├── dashboard.js
│       └── agents.js
└── templates/                   # Templates HTML
    ├── index.html
    ├── agents.html
    ├── workflows.html
    ├── hooks.html
    ├── events.html
    ├── secrets.html
    └── ssh.html
```

### Adicionar Nova Funcionalidade

1. **Backend**: Criar handler em `internal/webui/handlers/`
2. **Rota**: Adicionar rota em `server.go`
3. **Frontend**: Criar template em `templates/` e JS em `static/js/`
4. **WebSocket**: Adicionar tipo de mensagem se necessário

## Roadmap

- [ ] Gráficos e métricas visuais
- [ ] Execução de workflows pela UI
- [ ] Editor de código com syntax highlighting
- [ ] Logs em tempo real
- [ ] Notificações push
- [ ] Dashboard customizável
- [ ] Exportação de dados
- [ ] Tema dark mode
- [ ] Multi-idioma
- [ ] Suporte a múltiplos usuários

## Screenshots

### Dashboard
![Dashboard](../images/ui-dashboard.png)

### Agents
![Agents](../images/ui-agents.png)

### Workflows
![Workflows](../images/ui-workflows.png)

## Troubleshooting

### Porta já em uso
```bash
# Use outra porta
sloth-runner ui --port 8081
```

### WebSocket não conecta
Verifique se não há proxy ou firewall bloqueando a conexão WebSocket.

### Dados não aparecem
Verifique se os bancos de dados SQLite existem nos caminhos corretos:
```bash
ls -la .sloth-cache/
ls -la ~/.sloth-runner/
```

## Contribuindo

Para contribuir com a UI:

1. Fork o repositório
2. Crie uma branch para sua feature
3. Implemente a feature (backend + frontend)
4. Teste localmente
5. Submeta um Pull Request

## Licença

Mesma licença do projeto principal Sloth Runner.
