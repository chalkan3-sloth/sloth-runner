# Comando Sysadmin Health

O comando `sysadmin health` fornece ferramentas abrangentes para monitoramento proativo, verificação de conectividade e diagnóstico de problemas no sloth-runner.

## Visão Geral

```bash
sloth-runner sysadmin health [subcommand] [flags]
```

> **Nota:** Este comando faz parte da suite de ferramentas `sysadmin`. Para ver todos os comandos disponíveis para administradores de sistema, use `sloth-runner sysadmin --help`.

## Subcomandos

### 1. check - Executar Todos os Health Checks

Executa todos os health checks do sistema e exibe resultados consolidados.

**Uso:**
```bash
sloth-runner sysadmin health check [flags]
```

**Flags:**
- `-o, --output string` - Formato de saída: text, json (padrão: text)
- `-v, --verbose` - Mostrar saída detalhada com tempos de execução

**Health Checks Executados:**
1. **Database Connectivity** - Verifica conectividade com SQLite
2. **Data Directory** - Verifica existência e permissões de escrita
3. **Master Server** - Testa conexão com servidor master (se configurado)
4. **Log Directory** - Verifica diretório de logs
5. **Disk Space** - Verifica espaço em disco disponível
6. **Memory Usage** - Monitora uso de memória do processo

**Exemplos:**
```bash
# Executar todos os checks
sloth-runner sysadmin health check

# Output em JSON para parsing
sloth-runner sysadmin health check --output json

# Modo verbose com tempos de execução
sloth-runner sysadmin health check --verbose

# Salvar resultado em arquivo
sloth-runner sysadmin health check --output json > health-report.json
```

**Output Exemplo (Text):**
```
🏥 Health Check Report
═══════════════════════════════════════════════════

Timestamp: 2025-10-09 10:34:41
Status:    ✅ HEALTHY

📊 Summary:
   OK:      6

📋 Checks:
   ✅ Database Connectivity: Database is accessible
   ✅ Data Directory: Data directory is accessible and writable: /etc/sloth-runner
   ✅ Master Server: Master server is reachable: localhost:50053
   ✅ Log Directory: Log directory is accessible: /etc/sloth-runner/logs
   ✅ Disk Space: Disk space check passed
   ✅ Memory Usage: Memory usage normal: 2 MB allocated, 12 MB system
```

**Output Exemplo (JSON):**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-09T10:34:41-03:00",
  "checks": [
    {
      "name": "Database Connectivity",
      "status": "ok",
      "message": "Database is accessible",
      "duration": 455208,
      "timestamp": "2025-10-09T10:34:41-03:00"
    }
  ],
  "summary": {
    "ok": 6
  }
}
```

**Status Possíveis:**
- ✅ **healthy** - Todos os checks passaram
- ⚠️ **warning** - Alguns checks com avisos (não crítico)
- ❌ **error** - Um ou mais checks falharam
- 🔴 **critical** - Falhas críticas do sistema

**Casos de Uso:**
- Health check em CI/CD pipelines
- Monitoramento automatizado (cron/systemd timer)
- Validação pós-deployment
- Troubleshooting inicial de problemas

---

### 2. agent - Verificar Saúde de Agents

Verifica conectividade e status de agents específicos ou todos os agents registrados.

**Uso:**
```bash
sloth-runner sysadmin health agent [agent-name] [flags]
```

**Flags:**
- `--all` - Verificar todos os agents registrados

**Exemplos:**
```bash
# Verificar agent específico (detalhado)
sloth-runner sysadmin health agent do-sloth-runner-01

# Verificar todos os agents (resumo)
sloth-runner sysadmin health agent --all
```

**Output Exemplo (Agent Específico):**
```
🔍 Checking health of agent: do-sloth-runner-01

📋 Agent Information:
   Name:    do-sloth-runner-01
   Address: 68.183.52.244:50051
   Status:  Active
   Last Heartbeat: 2025-10-09 10:34:54 (2s ago)

🔌 Connectivity Test:
   ✅ Connection successful
   Duration: 278ms

✅ Overall Status: Agent is healthy
```

**Output Exemplo (Todos os Agents):**
```
🔍 Checking health of 5 agent(s)

Agent: do-sloth-runner-01
   Status: ✅ Healthy
   Address: 68.183.52.244:50051
   Last Heartbeat: 2025-10-09 10:34:38

Agent: do-sloth-runner-02
   Status: ✅ Healthy
   Address: 45.55.222.242:50051
   Last Heartbeat: 2025-10-09 10:34:41

Agent: keite-guica
   Status: ❌ Unreachable (context deadline exceeded)
   Address: 192.168.1.17:50051
   Last Heartbeat: 2025-10-09 07:17:57

📊 Summary:
   Total:   5
   Healthy: 2
   Error:   3
```

**Critérios de Health:**
- ✅ **Healthy** - Agent conectável e heartbeat recente (< 5min)
- ⚠️ **Stale** - Agent conectável mas heartbeat antigo (> 5min)
- ❌ **Unreachable** - Falha na conexão TCP/gRPC

**Casos de Uso:**
- Validar agents antes de executar workflows
- Troubleshooting de conectividade
- Monitoramento de fleet de agents
- Identificação de agents offline

---

### 3. master - Verificar Saúde do Servidor Master

Verifica conectividade TCP e gRPC com o servidor master.

**Uso:**
```bash
sloth-runner sysadmin health master [flags]
```

**Flags:**
- `--address string` - Endereço do master (default: da configuração)

**Exemplos:**
```bash
# Verificar master configurado
sloth-runner sysadmin health master

# Verificar master específico
sloth-runner sysadmin health master --address 192.168.1.100:50053
```

**Output Exemplo (Sucesso):**
```
🔍 Checking master server: localhost:50053

🔌 TCP Connectivity:
   ✅ TCP connection successful
   Duration: 1.234ms

🔌 gRPC Connectivity:
   ✅ gRPC connection successful
   Duration: 45.678ms

✅ Master server is healthy
```

**Output Exemplo (Falha):**
```
🔍 Checking master server: localhost:50053

🔌 TCP Connectivity:
   ❌ TCP connection failed: connection refused
```

**Testes Realizados:**
1. **TCP Connectivity** - Socket TCP básico (porta acessível)
2. **gRPC Connectivity** - Handshake gRPC completo

**Casos de Uso:**
- Validar configuração de master
- Troubleshooting de problemas de conexão
- Verificar firewall/network
- Health check de infraestrutura

---

### 4. watch - Monitoramento Contínuo

Monitora continuamente a saúde do sistema em intervalos especificados.

**Uso:**
```bash
sloth-runner sysadmin health watch [flags]
```

**Flags:**
- `-i, --interval string` - Intervalo entre checks (padrão: 30s)

**Intervalos Suportados:**
- `30s` - 30 segundos
- `1m` - 1 minuto
- `5m` - 5 minutos
- `1h` - 1 hora

**Exemplos:**
```bash
# Monitorar a cada 30 segundos
sloth-runner sysadmin health watch

# Monitorar a cada 1 minuto
sloth-runner sysadmin health watch --interval 1m

# Monitorar a cada 5 minutos
sloth-runner sysadmin health watch --interval 5m
```

**Output Exemplo:**
```
👀 Watching system health (interval: 30s)
Press Ctrl+C to stop

[10:34:41] ✅ HEALTHY | OK: 4
────────────────────────────────────────────────────────────
[10:35:11] ✅ HEALTHY | OK: 4
────────────────────────────────────────────────────────────
[10:35:41] ⚠️  WARNING | OK: 3 | WARN: 1
────────────────────────────────────────────────────────────
[10:36:11] ✅ HEALTHY | OK: 4
```

**Checks Monitorados:**
- Database Connectivity
- Data Directory
- Master Server Connection
- Log Directory

**Casos de Uso:**
- Monitoramento em terminal dedicado
- Acompanhamento durante manutenção
- Validação pós-deploy contínua
- Debug de problemas intermitentes

---

### 5. diagnostics - Relatório de Diagnóstico Completo

Gera relatório detalhado de diagnóstico incluindo informações do sistema, configuração e health checks.

**Uso:**
```bash
sloth-runner sysadmin health diagnostics [flags]
```

**Flags:**
- `-o, --output string` - Arquivo de saída (stdout se não especificado)

**Exemplos:**
```bash
# Exibir no terminal
sloth-runner sysadmin health diagnostics

# Salvar em arquivo
sloth-runner sysadmin health diagnostics --output diagnostics.json

# Pipe para análise
sloth-runner sysadmin health diagnostics | jq '.health_checks'
```

**Conteúdo do Relatório:**

```json
{
  "timestamp": "2025-10-09T10:34:56-03:00",
  "version": "dev",
  "system": {
    "os": "darwin",
    "arch": "arm64",
    "cpus": 8,
    "go_version": "go1.21.5"
  },
  "configuration": {
    "data_dir": "/etc/sloth-runner",
    "log_dir": "/etc/sloth-runner/logs",
    "master_address": "localhost:50053"
  },
  "health_checks": [
    {
      "name": "Database Connectivity",
      "status": "ok",
      "message": "Database is accessible",
      "duration": 382333,
      "timestamp": "2025-10-09T10:34:56-03:00"
    }
  ]
}
```

**Informações Incluídas:**
1. **Timestamp** - Data/hora da geração
2. **Version** - Versão do sloth-runner
3. **System** - SO, arquitetura, CPUs, versão Go
4. **Configuration** - Paths e configurações
5. **Health Checks** - Todos os checks com detalhes

**Casos de Uso:**
- Troubleshooting com suporte técnico
- Documentação de incidentes
- Análise de ambiente
- Auditoria de configuração

---

## Workflows Comuns

### Quick Health Check Diário
```bash
# Executar check rápido
sloth-runner sysadmin health check

# Se houver warnings/errors, investigar
sloth-runner sysadmin health agent --all
sloth-runner sysadmin health master
```

### Validação Pré-Deployment
```bash
# 1. Verificar sistema
sloth-runner sysadmin health check

# 2. Verificar todos os agents
sloth-runner sysadmin health agent --all

# 3. Verificar master
sloth-runner sysadmin health master

# 4. Gerar diagnóstico para registro
sloth-runner sysadmin health diagnostics --output pre-deploy-$(date +%Y%m%d).json
```

### Troubleshooting de Agent Problemático
```bash
# 1. Verificar health específico
sloth-runner sysadmin health agent problematic-agent

# 2. Ver logs do agent
sloth-runner logs tail --agent problematic-agent --level error

# 3. Verificar conectividade de rede
ping agent-host
telnet agent-host 50051

# 4. Ver informações completas
sloth-runner agent get problematic-agent
```

### Monitoramento Contínuo Durante Manutenção
```bash
# Terminal 1: Monitorar health
sloth-runner sysadmin health watch --interval 30s

# Terminal 2: Monitorar logs
sloth-runner logs tail --follow

# Terminal 3: Executar manutenção
# ... operações de manutenção ...
```

### Geração de Relatório de Saúde
```bash
#!/bin/bash
# daily-health-report.sh

REPORT_DATE=$(date +%Y-%m-%d)
REPORT_DIR="./health-reports"
mkdir -p $REPORT_DIR

echo "Generating health report for $REPORT_DATE..."

# Health check
sloth-runner sysadmin health check --output json > "$REPORT_DIR/check-$REPORT_DATE.json"

# Agent status
sloth-runner sysadmin health agent --all > "$REPORT_DIR/agents-$REPORT_DATE.txt"

# Full diagnostics
sloth-runner sysadmin health diagnostics --output "$REPORT_DIR/diagnostics-$REPORT_DATE.json"

echo "Report generated in $REPORT_DIR"
```

---

## Integração com Monitoramento

### Prometheus Exporter (Script Exemplo)
```bash
#!/bin/bash
# health-exporter.sh - Export health metrics for Prometheus

# Run health check
RESULT=$(sloth-runner sysadmin health check --output json)

# Parse and export metrics
echo "# HELP sloth_health_status Overall health status (0=error, 1=warning, 2=healthy)"
echo "# TYPE sloth_health_status gauge"
STATUS=$(echo $RESULT | jq -r '.status')
case $STATUS in
  "healthy") echo "sloth_health_status 2" ;;
  "warning") echo "sloth_health_status 1" ;;
  *) echo "sloth_health_status 0" ;;
esac

# Check-specific metrics
echo "# HELP sloth_health_checks Health checks summary"
echo "# TYPE sloth_health_checks gauge"
echo $RESULT | jq -r '.summary | to_entries | .[] | "sloth_health_checks{status=\"\(.key)\"} \(.value)"'
```

### Nagios/Icinga Check
```bash
#!/bin/bash
# check_sloth_health.sh - Nagios plugin

OUTPUT=$(sloth-runner sysadmin health check --output json)
STATUS=$(echo $OUTPUT | jq -r '.status')

case $STATUS in
  "healthy")
    echo "OK - Sloth-runner is healthy"
    exit 0
    ;;
  "warning")
    WARNINGS=$(echo $OUTPUT | jq -r '.summary.warning')
    echo "WARNING - $WARNINGS warnings detected"
    exit 1
    ;;
  *)
    ERRORS=$(echo $OUTPUT | jq -r '.summary.error // 0')
    echo "CRITICAL - $ERRORS errors detected"
    exit 2
    ;;
esac
```

### Systemd Timer (Health Check Periódico)
```ini
# /etc/systemd/system/sloth-health-check.service
[Unit]
Description=Sloth-runner Health Check
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/sloth-runner sysadmin health check
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```ini
# /etc/systemd/system/sloth-health-check.timer
[Unit]
Description=Run Sloth-runner Health Check every 5 minutes
Requires=sloth-health-check.service

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min
Unit=sloth-health-check.service

[Install]
WantedBy=timers.target
```

Ativar:
```bash
sudo systemctl enable sloth-health-check.timer
sudo systemctl start sloth-health-check.timer
```

### Alerting com Email
```bash
#!/bin/bash
# health-alert.sh - Send email on health issues

RESULT=$(sloth-runner sysadmin health check --output json)
STATUS=$(echo $RESULT | jq -r '.status')

if [ "$STATUS" != "healthy" ]; then
  SUBJECT="⚠️ Sloth-runner Health Alert: $STATUS"
  BODY=$(sloth-runner sysadmin health check)

  echo "$BODY" | mail -s "$SUBJECT" ops@example.com
fi
```

---

## Interpretação de Resultados

### Database Connectivity
- ✅ **OK**: SQLite acessível e respondendo
- ❌ **Error**: Arquivo não encontrado, permissões, ou corrupção

**Ações:**
```bash
# Verificar arquivo existe
ls -la /etc/sloth-runner/agents.db

# Verificar permissões
sudo chown $USER /etc/sloth-runner/agents.db

# Verificar integridade (se suspeita de corrupção)
sqlite3 /etc/sloth-runner/agents.db "PRAGMA integrity_check;"
```

### Data Directory
- ✅ **OK**: Diretório existe e tem permissão de escrita
- ⚠️ **Warning**: Diretório existe mas sem permissão de escrita
- ❌ **Error**: Diretório não existe

**Ações:**
```bash
# Criar diretório
sudo mkdir -p /etc/sloth-runner

# Ajustar permissões
sudo chown -R $USER /etc/sloth-runner
sudo chmod -R 755 /etc/sloth-runner
```

### Master Server
- ✅ **OK**: Master acessível via gRPC
- ⚠️ **Warning**: Master não configurado (normal para standalone)
- ❌ **Error**: Master configurado mas inacessível

**Ações:**
```bash
# Verificar configuração
echo $SLOTH_RUNNER_MASTER_ADDR

# Verificar conectividade TCP
telnet master-host 50053

# Verificar firewall
sudo ufw status

# Ver logs do master
sloth-runner logs tail --level error
```

### Log Directory
- ✅ **OK**: Diretório de logs existe e acessível
- ⚠️ **Warning**: Diretório não encontrado (será criado)
- ❌ **Error**: Path existe mas não é diretório

**Ações:**
```bash
# Criar diretório de logs
mkdir -p /etc/sloth-runner/logs

# Verificar espaço em disco
df -h /etc/sloth-runner/logs
```

### Disk Space
- ✅ **OK**: Espaço suficiente disponível
- ⚠️ **Warning**: Espaço limitado (< 10%)
- ❌ **Error**: Espaço crítico (< 5%)

**Ações:**
```bash
# Ver uso de disco
df -h /etc/sloth-runner

# Limpar logs antigos
sloth-runner logs rotate --force
gzip /etc/sloth-runner/logs/sloth-runner.log.*

# Limpar databases antigas (cuidado!)
# Fazer backup antes!
```

### Memory Usage
- ✅ **OK**: Uso de memória normal (< 1GB)
- ⚠️ **Warning**: Uso elevado (> 1GB)
- ❌ **Error**: Uso crítico (> 2GB)

**Ações:**
```bash
# Ver uso detalhado
ps aux | grep sloth-runner

# Verificar memory leaks
# Reiniciar se necessário
sudo systemctl restart sloth-runner
```

---

## Boas Práticas

1. **Health Checks Regulares:** Execute `health check` diariamente
2. **Monitoramento de Agents:** Use `health agent --all` antes de workflows críticos
3. **Alerting:** Configure alertas para falhas de health checks
4. **Documentação:** Salve diagnostics durante incidentes
5. **Automação:** Use systemd timers ou cron para checks periódicos
6. **Baseline:** Estabeleça baseline de saúde em ambiente saudável
7. **Trending:** Monitore tendências de performance ao longo do tempo

---

## Troubleshooting

### "No master server configured" Warning
```bash
# Isso é normal se você não está usando master/agent architecture
# Para configurar master:
export SLOTH_RUNNER_MASTER_ADDR=localhost:50053

# Ou adicionar ao ~/.bashrc
echo 'export SLOTH_RUNNER_MASTER_ADDR=localhost:50053' >> ~/.bashrc
```

### Todos os Agents "Unreachable"
```bash
# 1. Verificar se master está rodando
ps aux | grep sloth-runner

# 2. Verificar network
ping agent-host

# 3. Verificar portas
netstat -tlnp | grep 50051

# 4. Ver logs
sloth-runner logs tail --follow
```

### Health Check Trava/Timeout
```bash
# Possíveis causas:
# 1. Database lock
# 2. Network lento
# 3. Master não respondendo

# Solução:
# Kill processo travado
pkill -9 sloth-runner

# Verificar database
sqlite3 /etc/sloth-runner/agents.db "PRAGMA integrity_check;"

# Reiniciar services
sudo systemctl restart sloth-runner
```

---

## Ver Também

- [logs](logs-command.md) - Gerenciar e visualizar logs
- [agent](../agent.md) - Comandos de gerenciamento de agents
- [master](../master.md) - Servidor master
- [Monitoring Guide](../monitoring.md) - Guia de monitoramento
