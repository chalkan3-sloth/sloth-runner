# Comandos Sysadmin

O comando `sysadmin` agrupa todas as ferramentas essenciais para administração de sistema e operações do sloth-runner.

## Visão Geral

```bash
sloth-runner sysadmin [command] [flags]
```

## Comandos Disponíveis

### 📊 logs - Gerenciamento e Visualização de Logs

Ferramentas abrangentes para gerenciar, visualizar e analisar logs do sloth-runner.

```bash
sloth-runner sysadmin logs [subcommand]
```

**Subcomandos:**
- `tail` - Visualizar logs em tempo real (com ou sem follow)
- `search` - Buscar logs com filtros avançados
- `export` - Exportar logs em múltiplos formatos (text, json, csv)
- `rotate` - Rotacionar logs manualmente
- `level` - Alterar nível de logging dinâmicamente
- `remote` - Buscar logs de agents remotos via gRPC (sem SSH)

**Exemplos Rápidos:**
```bash
# Acompanhar logs em tempo real
sloth-runner sysadmin logs tail --follow

# Buscar erros na última hora
sloth-runner sysadmin logs search --query "error" --since 1h

# Exportar para JSON
sloth-runner sysadmin logs export --format json --output logs.json

# Buscar logs remotos de agent
sloth-runner sysadmin logs remote --agent do-sloth-runner-01 --system syslog
```

📖 **[Ver documentação completa de logs](logs-command.md)**

---

### 🏥 health - Health Checks e Diagnósticos

Monitoramento proativo, verificação de conectividade e diagnóstico de problemas.

```bash
sloth-runner sysadmin health [subcommand]
```

**Subcomandos:**
- `check` - Executar todos os health checks do sistema
- `agent` - Verificar saúde de agents específicos ou todos
- `master` - Verificar conectividade com master server
- `watch` - Monitoramento contínuo com intervalos configuráveis
- `diagnostics` - Gerar relatório completo de diagnóstico

**Health Checks:**
- ✅ Database Connectivity
- ✅ Data Directory (permissions)
- ✅ Master Server (connection)
- ✅ Log Directory
- ✅ Disk Space
- ✅ Memory Usage

**Exemplos Rápidos:**
```bash
# Check geral do sistema
sloth-runner sysadmin health check

# Verificar todos os agents
sloth-runner sysadmin health agent --all

# Monitorar continuamente
sloth-runner sysadmin health watch --interval 30s

# Gerar diagnóstico completo
sloth-runner sysadmin health diagnostics --output diagnostics.json
```

📖 **[Ver documentação completa de health](health-command.md)**

---

### 🔧 debug - Debugging e Troubleshooting

Ferramentas abrangentes para diagnóstico e troubleshooting de agents, conexões e workflows.

```bash
sloth-runner sysadmin debug [subcommand]
```

**Subcomandos:**
- `connection` - Testa conectividade com agent (TCP, DNS, gRPC, latência)
- `agent` - Diagnóstico completo de agent (config, status, system info)
- `workflow` - Análise de execuções de workflows

**Exemplos Rápidos:**
```bash
# Testar conectividade com agent
sloth-runner sysadmin debug connection web-01 --verbose

# Diagnóstico completo de agent
sloth-runner sysadmin debug agent web-01 --full

# Analisar últimas execuções de workflow
sloth-runner sysadmin debug workflow deploy-prod --last 5
```

📖 **[Ver documentação completa de debug](debug-command.md)**

---

### 💾 backup - Backup e Restore

_(Em desenvolvimento)_

Ferramentas para backup e recuperação de dados do sloth-runner.

**Planejado:**
- `create` - Criar backup completo
- `restore` - Restaurar de backup
- `list` - Listar backups disponíveis
- `verify` - Verificar integridade de backups
- `schedule` - Agendar backups automáticos

---

### 🔔 alert - Alerting e Notificações

_(Em desenvolvimento)_

Sistema de alertas e notificações proativas.

**Planejado:**
- `configure` - Configurar canais de notificação
- `rule` - Gerenciar regras de alerta
- `test` - Testar configurações de alerta
- `history` - Histórico de alertas

---

## Casos de Uso Comuns

### 1. Monitoramento Diário

```bash
# Health check rápido
sloth-runner sysadmin health check

# Ver últimos logs
sloth-runner sysadmin logs tail -n 50

# Verificar agents
sloth-runner sysadmin health agent --all
```

### 2. Troubleshooting de Problema

```bash
# 1. Check de saúde geral
sloth-runner sysadmin health check --verbose

# 2. Buscar erros recentes
sloth-runner sysadmin logs search --query "error" --since 1h

# 3. Verificar agent específico
sloth-runner sysadmin health agent problematic-agent

# 4. Gerar diagnóstico para análise
sloth-runner sysadmin health diagnostics --output issue-$(date +%Y%m%d).json
```

### 3. Manutenção e Arquivamento

```bash
# 1. Verificar espaço e logs
sloth-runner sysadmin health check
ls -lh /etc/sloth-runner/logs/

# 2. Exportar logs para backup
sloth-runner sysadmin logs export --format json --since 30d --output backup.json

# 3. Rotacionar e comprimir
sloth-runner sysadmin logs rotate --force
gzip /etc/sloth-runner/logs/sloth-runner.log.*

# 4. Verificar saúde pós-manutenção
sloth-runner sysadmin health check
```

### 4. Monitoramento Contínuo

```bash
# Terminal 1: Health monitoring
sloth-runner sysadmin health watch --interval 30s

# Terminal 2: Log monitoring
sloth-runner sysadmin logs tail --follow --level warn

# Terminal 3: Operações
# ... suas operações ...
```

### 5. Análise de Performance

```bash
# 1. Exportar logs para análise
sloth-runner sysadmin logs export --format json --since 24h --output performance.json

# 2. Gerar diagnóstico
sloth-runner sysadmin health diagnostics --output diagnostics.json

# 3. Analisar com jq
cat performance.json | jq -r 'group_by(.agent) | .[] | {agent: .[0].agent, count: length}'
```

---

## Workflows de Automação

### Script de Health Check Diário

```bash
#!/bin/bash
# daily-health.sh

DATE=$(date +%Y-%m-%d)
REPORT_DIR="./health-reports/$DATE"
mkdir -p "$REPORT_DIR"

echo "Running daily health check for $DATE..."

# Health check
sloth-runner sysadmin health check --output json > "$REPORT_DIR/health.json"

# Agent status
sloth-runner sysadmin health agent --all > "$REPORT_DIR/agents.txt"

# Export recent logs
sloth-runner sysadmin logs export --format json --since 24h --output "$REPORT_DIR/logs.json"

# Diagnostics
sloth-runner sysadmin health diagnostics --output "$REPORT_DIR/diagnostics.json"

# Check for issues
HEALTH_STATUS=$(jq -r '.status' "$REPORT_DIR/health.json")
if [ "$HEALTH_STATUS" != "healthy" ]; then
  echo "⚠️  WARNING: System is $HEALTH_STATUS"
  # Enviar alerta (email, slack, etc)
fi

echo "Report saved to $REPORT_DIR"
```

### Systemd Timer para Monitoring

```ini
# /etc/systemd/system/sloth-sysadmin-check.service
[Unit]
Description=Sloth-runner Sysadmin Health Check
After=network.target

[Service]
Type=oneshot
User=sloth-runner
ExecStart=/usr/local/bin/sloth-runner sysadmin health check
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```ini
# /etc/systemd/system/sloth-sysadmin-check.timer
[Unit]
Description=Run sloth-runner health check every 5 minutes

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min

[Install]
WantedBy=timers.target
```

Ativar:
```bash
sudo systemctl enable sloth-sysadmin-check.timer
sudo systemctl start sloth-sysadmin-check.timer
sudo systemctl status sloth-sysadmin-check.timer
```

### Cron Job para Log Rotation

```cron
# Rotacionar logs semanalmente às 00:00 de domingo
0 0 * * 0 /usr/local/bin/sloth-runner sysadmin logs rotate --force

# Health check a cada hora
0 * * * * /usr/local/bin/sloth-runner sysadmin health check --output json > /var/log/sloth-health-$(date +\%Y\%m\%d\%H).json
```

---

## Integração com Ferramentas Externas

### Prometheus

```bash
# Export métricas para Prometheus
sloth-runner sysadmin health check --output json | \
  jq -r '.summary | to_entries | .[] | "sloth_health_\(.key) \(.value)"'
```

### ELK Stack

```bash
# Export logs para Logstash
sloth-runner sysadmin logs export --format json --since 1h | \
  curl -X POST "http://logstash:5044" -H "Content-Type: application/json" -d @-
```

### Grafana

```bash
# Gerar métricas para Grafana
sloth-runner sysadmin health diagnostics | \
  jq '.health_checks[] | {time: .timestamp, metric: .name, value: (.status == "ok" | if . then 1 else 0 end)}'
```

---

## Boas Práticas

### 1. Monitoramento Regular
- ✅ Execute `health check` diariamente
- ✅ Configure alertas para status não-healthy
- ✅ Monitore crescimento de logs
- ✅ Verifique agents antes de workflows críticos

### 2. Gerenciamento de Logs
- ✅ Rotacione logs regularmente (semanal/mensal)
- ✅ Comprima logs antigos
- ✅ Archive logs com mais de 30 dias
- ✅ Mantenha apenas últimos 3 meses online

### 3. Troubleshooting
- ✅ Sempre comece com `health check`
- ✅ Use `logs search` para investigar erros
- ✅ Gere `diagnostics` ao reportar issues
- ✅ Documente problemas e soluções

### 4. Automação
- ✅ Use systemd timers ou cron
- ✅ Automatize backups
- ✅ Configure alerting
- ✅ Mantenha histórico de health checks

### 5. Segurança
- ✅ Restrinja acesso aos comandos sysadmin
- ✅ Proteja logs exportados (podem conter info sensível)
- ✅ Use permissões adequadas em diretórios
- ✅ Monitore tentativas de acesso não autorizado

---

## Troubleshooting Geral

### "Permission denied" em comandos sysadmin

```bash
# Verificar propriedade dos diretórios
ls -la /etc/sloth-runner/

# Ajustar permissões
sudo chown -R $USER:$USER /etc/sloth-runner
sudo chmod -R 755 /etc/sloth-runner
```

### "Database is locked"

```bash
# Verificar processos usando o database
lsof /etc/sloth-runner/*.db

# Se necessário, parar processos conflitantes
sudo systemctl stop sloth-runner
```

### Logs não aparecem ou estão incompletos

```bash
# Verificar se logs estão sendo escritos
ls -lh /etc/sloth-runner/logs/

# Ver logs do sistema
journalctl -u sloth-runner -f

# Verificar nível de log atual
sloth-runner sysadmin logs level info
```

### Health checks falhando intermitentemente

```bash
# Monitorar continuamente
sloth-runner sysadmin health watch --interval 10s

# Verificar recursos do sistema
top
df -h
free -h

# Ver logs de erro
sloth-runner sysadmin logs tail --level error --follow
```

---

## Performance Tips

### Otimização de Logs

```bash
# 1. Usar filtros sempre que possível
sloth-runner sysadmin logs search --query "error" --since 1h  # ✅ Rápido
sloth-runner sysadmin logs search --query "error"              # ❌ Lento

# 2. Limitar resultados
sloth-runner sysadmin logs search --query "error" --limit 100

# 3. Rotacionar logs grandes
sloth-runner sysadmin logs rotate --force
```

### Otimização de Health Checks

```bash
# Para checks rápidos, use apenas checks essenciais
# (futura feature: --quick flag)

# Em produção, use intervalos adequados
sloth-runner sysadmin health watch --interval 1m  # ✅ Bom para produção
sloth-runner sysadmin health watch --interval 10s # ❌ Muito frequente
```

---

## Ver Também

- [Agent Management](../agent.md) - Gerenciar agents
- [Workflow Execution](../workflow.md) - Executar workflows
- [Master Server](../master.md) - Servidor master
- [CLI Reference](../cli-reference.md) - Referência completa de comandos
