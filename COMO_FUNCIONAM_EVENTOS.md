# Como Funcionam os Eventos no Sloth Runner

## 📋 Resumo

Os eventos **NÃO são enviados automaticamente a cada hora**. Eles são **event-driven** (disparados por ações).

## 🎯 Quando os Eventos São Gerados

### 1. Agent Events

| Evento | Quando Ocorre |
|--------|---------------|
| `agent.registered` | Quando um agent se conecta ao master server |
| `agent.disconnected` | Quando um agent desconecta do master |
| `agent.heartbeat_failed` | Quando heartbeat de um agent falha |
| `agent.updated` | Quando um agent é atualizado |

**Como gerar**:
```bash
# Registrar um novo agent (gera agent.registered)
sloth-runner agent start --name my-agent --port 50051

# Parar o agent (gera agent.disconnected)
# Ctrl+C no terminal do agent
```

### 2. Task Events

| Evento | Quando Ocorre |
|--------|---------------|
| `task.started` | Quando uma task começa a executar |
| `task.completed` | Quando uma task completa com sucesso |
| `task.failed` | Quando uma task falha |

**Como gerar**:
```bash
# Executar um workflow (gera task.started + task.completed)
sloth-runner run my-workflow --file examples/basic.sloth

# Executar workflow com task que falha (gera task.failed)
sloth-runner run failing-workflow --file examples/fail.sloth
```

### 3. Custom Events (de workflows Lua)

```lua
-- Dentro de um workflow .sloth
event.dispatch("custom", {
    name = "deployment_complete",
    environment = "production",
    version = "1.2.3"
})
```

## 📊 Estado Atual do Sistema

Verificado em 2025-10-07:

```bash
$ sloth-runner events list --limit 50
```

**Eventos Coletados**:
- ✅ `agent.registered` - 20 eventos (agents se conectando)
- ✅ `agent.disconnected` - 1 evento
- ✅ `task.started` - 21 eventos
- ✅ `task.completed` - 14 eventos
- ✅ `task.failed` - Vários eventos (alguns com status failed)

**Datas**: Todos de 2025-10-06 (ontem) - **porque não houve atividade nova desde então**

## 🔄 Por Que Não Há Eventos Novos?

Os eventos são **disparados por ações**, não por tempo. Se você não está:

1. ❌ Executando workflows
2. ❌ Conectando/desconectando agents
3. ❌ Executando tasks

Então **não haverá eventos novos**.

## ✅ Como Testar a Coleta de Eventos

### Teste 1: Gerar Agent Event
```bash
# Terminal 1: Iniciar master
sloth-runner master --port 50053

# Terminal 2: Conectar um agent (gera agent.registered)
sloth-runner agent start --name test-agent --port 50051

# Verificar eventos
sloth-runner events list --type agent.registered --limit 5
```

### Teste 2: Gerar Task Events
```bash
# Criar um workflow simples
cat > /tmp/test.sloth <<'EOF'
task("hello")
    :description("Say hello")
    :command(function()
        print("Hello from Sloth Runner!")
        return true
    end)
    :build()

workflow.define("test", {
    tasks = { task("hello") }
})
EOF

# Executar (gera task.started + task.completed)
sloth-runner run test --file /tmp/test.sloth

# Verificar eventos
sloth-runner events list --type task.started --limit 5
sloth-runner events list --type task.completed --limit 5
```

### Teste 3: Gerar Task Failed Event
```bash
# Workflow que falha propositalmente
cat > /tmp/fail.sloth <<'EOF'
task("fail")
    :description("Task that fails")
    :command(function()
        error("Intentional failure for testing")
    end)
    :build()

workflow.define("test_fail", {
    tasks = { task("fail") }
})
EOF

# Executar (gera task.started + task.failed)
sloth-runner run test_fail --file /tmp/fail.sloth

# Verificar eventos de falha
sloth-runner events list --type task.failed --limit 5
```

## 🌐 Interface Web

Após gerar eventos, acesse a interface web:

```bash
# Iniciar Web UI
sloth-runner ui --port 8080

# Acessar no navegador
# Dashboard: http://localhost:8080/
# Events: http://localhost:8080/events
```

A interface mostrará:
- ✅ Contagem total de eventos
- ✅ Eventos pendentes/processing/completed/failed
- ✅ Tabela com eventos recentes
- ✅ Detalhes de cada evento (payload JSON)
- ✅ Auto-refresh a cada 5 segundos

## 🔍 Verificar Eventos em Tempo Real

```bash
# Terminal 1: Monitorar eventos
watch -n 2 'sloth-runner events list --limit 10'

# Terminal 2: Executar workflows
sloth-runner run my-workflow --file workflow.sloth

# Você verá novos eventos aparecerem no Terminal 1!
```

## 📈 Estatísticas do Sistema

```bash
# Ver estatísticas gerais
sloth-runner events list | head -n 20

# Contar eventos por tipo
sloth-runner events list --type agent.registered | grep "Total events"
sloth-runner events list --type task.started | grep "Total events"
sloth-runner events list --type task.completed | grep "Total events"
```

## 🎯 Eventos em Produção

Em um ambiente de produção real:

1. **Agents conectando/desconectando** → Gera eventos continuamente
2. **Workflows executando** (manual ou scheduled) → Gera task events
3. **Deployments** → Gera custom events via Lua
4. **Monitoramento** → Agents reportando health checks

**Resultado**: Fluxo constante de eventos para rastrear toda a atividade do sistema.

## 🔧 Hooks e Automação

Você pode configurar hooks para reagir aos eventos:

```bash
# Registrar um hook que reage a task.failed
sloth-runner hook add \
  --name alert-on-failure \
  --event task.failed \
  --command "curl -X POST https://slack.com/webhook -d '{\"text\":\"Task failed!\"}'"

# Agora quando uma task falhar, o webhook será chamado automaticamente!
```

## 📝 Resumo

| Aspecto | Descrição |
|---------|-----------|
| **Frequência** | Event-driven (não por tempo) |
| **Trigger** | Ações (tasks, agents, workflows) |
| **Persistência** | SQLite (`.sloth-cache/hooks.db`) |
| **Processamento** | Assíncrono (100 workers, buffer de 1000) |
| **Retenção** | Ilimitada (até fazer cleanup manual) |
| **Visualização** | CLI (`sloth-runner events list`) ou Web UI |

## 🚀 Conclusão

**Eventos aparecem quando há atividade**. Para ver eventos novos na UI:

1. Execute workflows
2. Conecte/desconecte agents
3. Execute tasks manualmente
4. Configure schedulers para execução periódica

A interface web mostrará todos os eventos em tempo real! 🎉
