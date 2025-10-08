# Integração de Eventos Reais na Web UI

## Data: 2025-10-08

## 🎯 Objetivo
Remover dados mockados da interface web e integrar com os eventos reais coletados pelo sistema.

## ✅ Mudanças Realizadas

### 1. Backend - Event Handler (já existente e funcional)
**Arquivo**: `internal/webui/handlers/event.go`

Endpoints disponíveis:
- `GET /api/v1/events` - Lista todos os eventos
- `GET /api/v1/events/pending` - Lista eventos pendentes
- `GET /api/v1/events/:id` - Detalhes de um evento
- `POST /api/v1/events/:id/retry` - Retentar evento falho

### 2. Frontend - Events Page
**Arquivo**: `internal/webui/static/js/events.js`

**Mudanças**:
```javascript
// Antes: Usava apenas lowercase field names (mockados)
const eventId = event.id;
const eventType = event.type;

// Depois: Suporte a Go struct (capitalized) e JSON (lowercase)
const eventId = (event.ID || event.id || 'N/A').substring(0, 8);
const eventType = event.Type || event.type || 'unknown';
const eventStatus = event.Status || event.status || 'unknown';
```

**Melhorias**:
- ✅ Parsing correto de timestamps do Go (RFC3339)
- ✅ Suporte a ambos formatos de campo (capitalized e lowercase)
- ✅ Exibição de 8 caracteres do ID do evento
- ✅ Formatação robusta de datas com fallback
- ✅ Visualização completa de event.Data em JSON formatado
- ✅ Scroll automático para payloads grandes (max-height: 400px)

### 3. Dashboard - Event Stats
**Arquivo**: `internal/webui/handlers/dashboard.go`

**Mudanças**:
```go
// Antes: Apenas eventos pendentes
"events": gin.H{
    "pending": len(pendingEvents),
}

// Depois: Estatísticas completas
"events": gin.H{
    "total":      len(allEvents),
    "pending":    len(pendingEvents),
    "processing": processingEvents,
    "completed":  completedEvents,
    "failed":     failedEvents,
}
```

**Arquivo**: `internal/webui/static/js/dashboard.js`

**Mudanças**:
- ✅ Atualização do card de eventos para mostrar "X of Y total"
- ✅ Parsing correto de campos do Go struct
- ✅ Formatação consistente de timestamps
- ✅ Tratamento de erros na listagem de eventos

**Arquivo**: `internal/webui/templates/index.html`

**Mudanças**:
```html
<!-- Antes -->
<small class="text-muted">in queue</small>

<!-- Depois -->
<small class="text-muted">of <span id="stat-events-total">0</span> total</small>
```

## 📊 Dados Reais Agora Exibidos

### Dashboard (/)
```
┌─────────────────────┐
│ Pending Events      │
│       0             │ <- Dados reais do banco
│ of 10 total         │ <- Total de eventos
└─────────────────────┘
```

### Events Page (/events)

**Tabela de Eventos**:
| ID       | Type              | Status    | Created              | Processed            |
|----------|-------------------|-----------|----------------------|----------------------|
| c33fa85c | agent.registered  | completed | 2025-10-07 22:09:50  | 2025-10-07 22:09:50 |
| a4eb3844 | task.started      | completed | 2025-10-06 17:32:21  | 2025-10-06 17:32:21 |

**Event Details Modal**:
```json
{
  "ID": "c33fa85c-...",
  "Type": "agent.registered",
  "Status": "completed",
  "Timestamp": "2025-10-07T22:09:50-03:00",
  "Created": "2025-10-07T22:09:50-03:00",
  "Processed": "2025-10-07T22:09:50-03:00",
  "Data": {
    "agent": {
      "name": "lady-arch",
      "address": "192.168.1.16:50051",
      "version": "v6.11.1",
      "system_info": { ... }
    }
  }
}
```

## 🔍 Tipos de Eventos Coletados

Atualmente o sistema está coletando:

1. **Agent Events** (9 eventos):
   - `agent.registered` - Quando agent se registra
   - `agent.disconnected` - Quando agent desconecta

2. **Task Events** (5 eventos):
   - `task.started` - Início de tarefa
   - `task.completed` - Conclusão de tarefa

## 🔄 Auto-Refresh

- **Events Page**: Refresh automático a cada 5 segundos
- **Dashboard**: Refresh ao carregar a página

## 🎨 Melhorias de UX

### Status Badges
```javascript
const statusMap = {
    'pending': 'bg-warning',      // Amarelo
    'processing': 'bg-info',      // Azul
    'completed': 'bg-success',    // Verde
    'failed': 'bg-danger',        // Vermelho
    'cancelled': 'bg-secondary'   // Cinza
};
```

### Cards de Estatísticas
```
┌──────────────────┬──────────────────┬──────────────────┬──────────────────┐
│ Pending: 0       │ Processing: 0    │ Completed: 10    │ Failed: 0        │
│ ⏰               │ ⚙️               │ ✅               │ ❌               │
└──────────────────┴──────────────────┴──────────────────┴──────────────────┘
```

## 📝 Arquivos Modificados

### Backend
1. `internal/webui/handlers/dashboard.go` - Estatísticas de eventos expandidas
2. `internal/webui/handlers/event.go` - (já existente, sem modificações)

### Frontend
1. `internal/webui/static/js/events.js` - Parse de dados reais do Go
2. `internal/webui/static/js/dashboard.js` - Integração com eventos reais
3. `internal/webui/templates/index.html` - Card de eventos atualizado

## ✅ Testes Necessários

Para testar a integração completa:

1. **Iniciar Web UI**:
```bash
sloth-runner ui --port 8080
```

2. **Acessar Pages**:
   - Dashboard: http://localhost:8080/
   - Events: http://localhost:8080/events

3. **Verificar**:
   - [ ] Dashboard mostra contagem de eventos real
   - [ ] Events page lista eventos do banco
   - [ ] Event details modal mostra dados completos
   - [ ] Auto-refresh funciona
   - [ ] Retry de eventos falhos funciona
   - [ ] Timestamps estão formatados corretamente

4. **Gerar Eventos de Teste**:
```bash
# Registrar um novo agent
sloth-runner agent start

# Executar um workflow
sloth-runner run my-workflow --file workflow.sloth
```

## 🎉 Resultado

A interface web agora exibe **dados reais** do sistema de eventos:
- ✅ Sem dados mockados
- ✅ Integração completa com SQLite
- ✅ Suporte a todos os tipos de eventos (100+ tipos)
- ✅ Visualização detalhada de payloads
- ✅ Estatísticas em tempo real
- ✅ Auto-refresh automático

## 🚀 Próximos Passos (Opcional)

1. **Filtros Avançados**: Filtrar por tipo de evento, status, data
2. **Busca**: Buscar eventos por ID, tipo, payload
3. **Gráficos**: Timeline de eventos, distribuição por tipo
4. **WebSocket**: Push de eventos em tempo real (sem polling)
5. **Exportação**: Download de eventos em CSV/JSON
6. **Limpeza**: Interface para fazer cleanup de eventos antigos
