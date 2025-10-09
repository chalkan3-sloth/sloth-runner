# Comando Debug

O comando `debug` fornece ferramentas para diagnóstico e troubleshooting de problemas no sloth-runner, incluindo conectividade com agents, configuração de agents e análise de execuções de workflows.

## Visão Geral

```bash
sloth-runner sysadmin debug [subcommand] [flags]
```

## Subcomandos Disponíveis

### 🔌 connection - Debug de Conectividade

Testa e diagnostica conectividade com agents específicos, incluindo TCP, DNS, gRPC e latência.

```bash
sloth-runner sysadmin debug connection [agent-name] [flags]
```

**Flags:**
- `-t, --timeout int` - Timeout de conexão em segundos (padrão: 5)
- `-v, --verbose` - Saída detalhada com testes adicionais

**Testes Executados:**

1. **Test 1: TCP Connection** - Testa conexão TCP básica com o agent
2. **Test 2: DNS Resolution** - Resolve DNS do host (apenas com --verbose)
3. **Test 3: gRPC Connection** - Testa conexão gRPC completa
4. **Test 4: Agent RPC Call** - Executa chamada RPC real (GetResourceUsage)
5. **Test 5: Latency Test** - Testa latência com 3 pings (apenas com --verbose)

**Exemplos:**

```bash
# Teste básico de conexão
sloth-runner sysadmin debug connection do-sloth-runner-01

# Teste verbose com timeout customizado
sloth-runner sysadmin debug connection do-sloth-runner-02 --verbose --timeout 10

# Teste rápido de múltiplos agents
for agent in web-01 web-02 web-03; do
  sloth-runner sysadmin debug connection $agent
done
```

**Exemplo de Saída:**

```
🔍 Debugging connection to agent: do-sloth-runner-01

📋 Agent Info:
  Name:     do-sloth-runner-01
  Address:  68.183.52.244:50051
  Status:   Active
  Last HB:  2025-10-09 13:39:42

🔌 Test 1: TCP Connection to 68.183.52.244:50051
  ✅ SUCCESS
  Duration: 141ms

🔗 Test 3: gRPC Connection
  ✅ SUCCESS
  Duration: 318ms

💓 Test 4: Agent RPC Call
  ✅ SUCCESS
  Duration: 147ms

✅ All tests completed successfully!
```

**Exemplo de Saída Verbose:**

```
🔍 Debugging connection to agent: do-sloth-runner-02

📋 Agent Info:
  Name:     do-sloth-runner-02
  Address:  45.55.222.242:50051
  Status:   Active
  Last HB:  2025-10-09 13:39:47

🔌 Test 1: TCP Connection to 45.55.222.242:50051
  ✅ SUCCESS
  Duration: 139ms

🌐 Test 2: DNS Resolution
  ✅ SUCCESS
  IP: 45.55.222.242

🔗 Test 3: gRPC Connection
  ✅ SUCCESS
  Duration: 276ms

💓 Test 4: Agent RPC Call
  ✅ SUCCESS
  Duration: 135ms

⏱️  Test 5: Latency Test (3 pings)
  Ping 1: ✅ 134ms
  Ping 2: ✅ 137ms
  Ping 3: ✅ 134ms
  Average: 135ms

✅ All tests completed successfully!
```

---

### 🔧 agent - Diagnóstico de Agent

Obtém diagnósticos detalhados sobre um agent, incluindo configuração, status e informações de sistema.

```bash
sloth-runner sysadmin debug agent [agent-name] [flags]
```

**Flags:**
- `-f, --full` - Diagnóstico completo incluindo informações de sistema

**Informações Fornecidas:**

1. **Basic Information:**
   - Nome do agent
   - Endereço (IP:porta)
   - Status (Active/Inactive)
   - Versão do agent
   - Timestamp de registro
   - Última atualização
   - Último heartbeat

2. **Connection Status:**
   - ✅ HEALTHY - Heartbeat recente (< 60s)
   - ⚠️ WARNING - Heartbeat antigo (60s-300s)
   - ❌ CRITICAL - Sem heartbeat (> 300s)

3. **System Information** (com --full):
   - Arquitetura e hostname
   - Kernel e versão
   - CPUs e memória
   - Discos e montagens
   - Interfaces de rede
   - Pacotes instalados
   - Variáveis de ambiente
   - Load average

4. **Recommendations:**
   - Ações sugeridas baseadas no status
   - Verificações recomendadas
   - Passos de troubleshooting

**Exemplos:**

```bash
# Diagnóstico básico
sloth-runner sysadmin debug agent do-sloth-runner-01

# Diagnóstico completo com system info
sloth-runner sysadmin debug agent do-sloth-runner-02 --full

# Verificar todos os agents registrados
sloth-runner agent list | tail -n +2 | awk '{print $1}' | while read agent; do
  echo "=== $agent ==="
  sloth-runner sysadmin debug agent $agent
done
```

**Exemplo de Saída Básica:**

```
🔍 Agent Diagnostics: do-sloth-runner-01

📋 Basic Information:
  Name:             do-sloth-runner-01
  Address:          68.183.52.244:50051
  Status:           Active
  Version:          6.12.1
  Registered At:    2025-10-07 19:36:52
  Last Updated:     2025-10-09 13:39:52
  Last Heartbeat:   2025-10-09 13:39:52 (5s ago)

🔌 Connection Status:
  ✅ HEALTHY - Recent heartbeat

💡 Recommendations:
```

**Exemplo de Saída com --full:**

```
🔍 Agent Diagnostics: do-sloth-runner-02

📋 Basic Information:
  Name:             do-sloth-runner-02
  Address:          45.55.222.242:50051
  Status:           Active
  Version:          6.12.1
  Registered At:    2025-10-07 19:51:27
  Last Updated:     2025-10-09 13:39:57
  Last Heartbeat:   2025-10-09 13:39:57 (5s ago)

🔌 Connection Status:
  ✅ HEALTHY - Recent heartbeat

💻 System Information:
{
  "architecture": "amd64",
  "hostname": "do-sloth-runner-02",
  "kernel": "Linux",
  "kernel_version": "5.15.0-113-generic",
  "cpus": 1,
  "memory": {
    "total": 476258304,
    "available": 289492992,
    "used": 186765312,
    "used_percent": 39.21
  },
  "disk": [
    {
      "device": "/dev/vda1",
      "mountpoint": "/",
      "total": 10213466112,
      "free": 7706378240,
      "used": 2490310656,
      "used_percent": 24.38
    }
  ],
  ...
}

💡 Recommendations:
```

---

### 📊 workflow - Debug de Workflow

Analisa execuções de workflows, mostra detalhes de tasks e identifica problemas.

```bash
sloth-runner sysadmin debug workflow [workflow-name|latest] [flags]
```

**Flags:**
- `-n, --last int` - Número de execuções para mostrar (padrão: 1)

**Informações Fornecidas:**

Para cada execução:
- ID da execução
- Nome do workflow
- Group (se aplicável)
- Agent (se aplicável)
- Status (running, completed, failed, etc.)
- Timestamp de início
- Timestamp de fim e duração
- Estatísticas de tasks (total, sucesso, falha)
- Mensagem de erro (se houver)

**Exemplos:**

```bash
# Debug da última execução
sloth-runner sysadmin debug workflow latest

# Debug de workflow específico
sloth-runner sysadmin debug workflow deploy-prod

# Debug das últimas 5 execuções
sloth-runner sysadmin debug workflow latest --last 5

# Debug das últimas 10 execuções de workflow específico
sloth-runner sysadmin debug workflow backup-daily --last 10
```

**Exemplo de Saída (sem execuções):**

```
🔍 Workflow Debug: latest

No workflow executions found
```

**Exemplo de Saída (com execuções):**

```
🔍 Workflow Debug: deploy-prod

📊 Execution #1:
  ID:           550e8400-e29b-41d4-a716-446655440000
  Workflow:     deploy-prod
  Group:        production
  Agent:        web-01
  Status:       completed
  Start Time:   2025-10-09 14:30:15
  End Time:     2025-10-09 14:35:42
  Duration:     5m27s
  Tasks:        12 total, 12 success, 0 failed

📊 Execution #2:
  ID:           6ba7b810-9dad-11d1-80b4-00c04fd430c8
  Workflow:     deploy-prod
  Agent:        web-02
  Status:       failed
  Start Time:   2025-10-09 13:15:30
  End Time:     2025-10-09 13:16:45
  Duration:     1m15s
  Tasks:        12 total, 8 success, 4 failed
  Error:        Task 'restart-service' failed: connection timeout
```

---

## Casos de Uso Comuns

### 1. Troubleshooting de Agent Não Responsivo

```bash
# Passo 1: Verificar diagnóstico do agent
sloth-runner sysadmin debug agent problematic-agent

# Passo 2: Testar conectividade
sloth-runner sysadmin debug connection problematic-agent --verbose

# Passo 3: Ver logs do agent
sloth-runner sysadmin logs remote --agent problematic-agent --system syslog --lines 100
```

### 2. Análise de Falha de Workflow

```bash
# Passo 1: Ver execuções recentes
sloth-runner sysadmin debug workflow failed-workflow --last 5

# Passo 2: Verificar agent usado
sloth-runner sysadmin debug agent web-01

# Passo 3: Ver logs do período
sloth-runner sysadmin logs search --query "error" --since 1h
```

### 3. Health Check Pré-Deploy

```bash
#!/bin/bash
# pre-deploy-check.sh

echo "Checking all agents before deploy..."

AGENTS=$(sloth-runner agent list | tail -n +2 | awk '{print $1}')

for agent in $AGENTS; do
  echo "=== Checking $agent ==="

  # Test connection
  if ! sloth-runner sysadmin debug connection $agent; then
    echo "❌ Connection failed for $agent"
    exit 1
  fi

  # Check diagnostics
  sloth-runner sysadmin debug agent $agent | grep -q "HEALTHY"
  if [ $? -ne 0 ]; then
    echo "⚠️  Agent $agent is not healthy"
    exit 1
  fi
done

echo "✅ All agents are healthy. Proceeding with deploy..."
```

### 4. Performance Analysis

```bash
# Analisar performance de workflows
sloth-runner sysadmin debug workflow data-processing --last 20 | \
  grep "Duration" | \
  awk '{print $2}' | \
  sort -n

# Testar latência de todos os agents
for agent in $(sloth-runner agent list | tail -n +2 | awk '{print $1}'); do
  echo "=== $agent ==="
  sloth-runner sysadmin debug connection $agent --verbose 2>&1 | grep "Average"
done
```

### 5. Automated Monitoring

```bash
#!/bin/bash
# monitor-agents.sh

while true; do
  DATE=$(date +"%Y-%m-%d %H:%M:%S")

  for agent in $(sloth-runner agent list | tail -n +2 | awk '{print $1}'); do
    # Test connection
    if ! sloth-runner sysadmin debug connection $agent &>/dev/null; then
      echo "[$DATE] ❌ Connection failed: $agent" | tee -a /var/log/agent-monitor.log
      # Enviar alerta
      curl -X POST "https://alerts.example.com/webhook" \
        -d "{\"agent\": \"$agent\", \"status\": \"connection_failed\"}"
    fi

    # Check health
    STATUS=$(sloth-runner sysadmin debug agent $agent | grep "Connection Status" -A 1 | tail -n 1)
    if echo "$STATUS" | grep -q "CRITICAL"; then
      echo "[$DATE] ⚠️  Critical status: $agent" | tee -a /var/log/agent-monitor.log
    fi
  done

  sleep 300  # Check every 5 minutes
done
```

---

## Workflows de Troubleshooting

### Workflow 1: Agent Não Aparece

**Sintoma:** Agent não aparece em `sloth-runner agent list`

**Debug:**
```bash
# 1. Verificar se o agent está rodando
ssh user@agent-host "ps aux | grep sloth-runner"

# 2. Ver logs do agent
ssh user@agent-host "tail -100 /var/log/sloth-runner-agent.log"

# 3. Verificar conectividade de rede
sloth-runner sysadmin debug connection agent-name

# 4. Verificar master server
sloth-runner sysadmin health master
```

### Workflow 2: Workflow Falha Intermitentemente

**Sintoma:** Workflow funciona às vezes, falha outras vezes

**Debug:**
```bash
# 1. Ver histórico de execuções
sloth-runner sysadmin debug workflow problematic-workflow --last 10

# 2. Identificar padrão (agent, horário, etc.)

# 3. Testar latência do agent
sloth-runner sysadmin debug connection target-agent --verbose

# 4. Ver system info do agent
sloth-runner sysadmin debug agent target-agent --full

# 5. Verificar recursos (memória, disco)
# Analisar output do --full acima
```

### Workflow 3: Conectividade Lenta

**Sintoma:** Workflows executam lentamente

**Debug:**
```bash
# 1. Testar latência de todos os agents
for agent in $(sloth-runner agent list | tail -n +2 | awk '{print $1}'); do
  echo "=== $agent ==="
  sloth-runner sysadmin debug connection $agent --verbose
done

# 2. Comparar latências

# 3. Verificar recursos dos agents lentos
sloth-runner sysadmin debug agent slow-agent --full

# 4. Ver workflows recentes no agent
sloth-runner sysadmin debug workflow latest --last 5
```

### Workflow 4: Agent Desconecta Frequentemente

**Sintoma:** Agent mostra "CRITICAL - No heartbeat"

**Debug:**
```bash
# 1. Verificar status atual
sloth-runner sysadmin debug agent unstable-agent

# 2. Testar conectividade
sloth-runner sysadmin debug connection unstable-agent --verbose

# 3. Ver logs do agent remotamente
sloth-runner sysadmin logs remote --agent unstable-agent --system syslog --lines 200

# 4. Verificar logs do master
sloth-runner sysadmin logs tail --follow | grep unstable-agent

# 5. Verificar recursos do sistema
sloth-runner sysadmin debug agent unstable-agent --full | jq '.memory, .disk'
```

---

## Interpretação de Resultados

### Connection Tests

**TCP Connection FAILED:**
- Firewall bloqueando porta
- Agent não está rodando
- Endereço IP incorreto
- Rede inacessível

**DNS Resolution FAILED:**
- Problema de DNS
- Hostname incorreto
- /etc/hosts não configurado

**gRPC Connection FAILED:**
- Agent não está respondendo gRPC
- Certificados/TLS incorretos (se usando SSL)
- Versão incompatível do protocolo

**Agent RPC Call FAILED:**
- Agent rodando mas não funcional
- Recursos do agent esgotados
- Bug no agent

### Agent Health Status

**✅ HEALTHY (< 60s desde último heartbeat):**
- Agent funcionando normalmente
- Conectividade boa
- Nenhuma ação necessária

**⚠️ WARNING (60s-300s):**
- Possível problema de rede
- Agent com carga alta
- Monitorar de perto

**❌ CRITICAL (> 300s):**
- Agent provavelmente offline
- Problema sério de rede/firewall
- Investigação imediata necessária

### Workflow Status

**completed:**
- Workflow executou com sucesso
- Todas as tasks completaram

**failed:**
- Uma ou mais tasks falharam
- Ver error_message para detalhes
- Verificar agent que executou

**running:**
- Workflow ainda executando
- Verificar se está travado
- Pode precisar de timeout

---

## Integração com Outras Ferramentas

### Com JQ para análise

```bash
# Exportar diagnóstico para JSON (future feature)
# sloth-runner sysadmin debug agent web-01 --output json | jq .

# Analisar memória de múltiplos agents
for agent in $(sloth-runner agent list | tail -n +2 | awk '{print $1}'); do
  echo -n "$agent: "
  sloth-runner sysadmin debug agent $agent --full 2>/dev/null | \
    jq -r '.memory.used_percent'
done
```

### Com Prometheus

```bash
# Exportar métricas de latência para Prometheus
# (implementar como custom exporter)

#!/bin/bash
# prometheus-agent-exporter.sh

while true; do
  for agent in $(sloth-runner agent list | tail -n +2 | awk '{print $1}'); do
    LATENCY=$(sloth-runner sysadmin debug connection $agent --verbose 2>&1 | \
              grep "Average:" | awk '{print $2}' | sed 's/ms//')

    echo "sloth_agent_latency{agent=\"$agent\"} $LATENCY" >> /var/lib/prometheus/node-exporter/agent-metrics.prom
  done
  sleep 60
done
```

### Com Grafana

```bash
# Coletar dados para dashboard Grafana

# Latência de agents
sloth-runner sysadmin debug connection $AGENT --verbose | grep Average

# Status de agents
sloth-runner sysadmin debug agent $AGENT | grep "Connection Status"

# Taxa de sucesso de workflows
sloth-runner sysadmin debug workflow latest --last 100 | \
  grep Status | awk '{print $2}' | sort | uniq -c
```

---

## Tips e Best Practices

### 1. Troubleshooting Proativo

✅ Execute debug periodicamente, não apenas quando há problemas
✅ Monitore latência de agents em produção
✅ Mantenha histórico de diagnósticos para comparação
✅ Configure alertas baseados em thresholds

### 2. Debugging Eficiente

✅ Sempre comece com `debug agent` para overview
✅ Use `debug connection --verbose` para problemas de rede
✅ Analise últimas 5-10 execuções de workflow, não apenas a última
✅ Capture diagnóstico ANTES de reiniciar agent com problema

### 3. Documentação

✅ Documente problemas recorrentes e soluções
✅ Crie playbooks para cenários comuns
✅ Mantenha inventário atualizado de agents e configurações
✅ Use comments em scripts de troubleshooting

### 4. Automação

✅ Automatize health checks regulares
✅ Crie scripts para diagnósticos comuns
✅ Integre com sistema de alertas
✅ Log de todas as operações de debug para auditoria

---

## Troubleshooting do Comando Debug

### "agent not found"

```bash
# Verificar agents registrados
sloth-runner agent list

# Verificar database
sqlite3 /etc/sloth-runner/agents.db "SELECT name, address FROM agents"
```

### "failed to open database"

```bash
# Verificar permissões
ls -la /etc/sloth-runner/*.db

# Ajustar se necessário
sudo chown $USER:$USER /etc/sloth-runner/*.db
```

### "connection timeout"

```bash
# Aumentar timeout
sloth-runner sysadmin debug connection agent-name --timeout 30

# Verificar firewall
sudo iptables -L | grep 50051

# Testar conectividade básica
telnet agent-host 50051
```

### "No workflow executions found"

Isso é normal se:
- Nenhum workflow foi executado ainda
- Workflows foram executados antes do histórico ser implementado
- Database foi limpo/resetado

---

## Referências

- [Comando Agent](agent.md) - Gerenciamento de agents
- [Comando Health](health-command.md) - Health checks
- [Comando Logs](logs-command.md) - Gerenciamento de logs
- [Sysadmin Overview](sysadmin.md) - Visão geral de comandos sysadmin
