# Comandos Sysadmin

O comando `sysadmin` agrupa todas as ferramentas essenciais para administração de sistema e operações do sloth-runner.

## Visão Geral

```bash
sloth-runner sysadmin [command] [flags]
```

O comando `sysadmin` oferece ferramentas completas para administração e operação do sloth-runner, incluindo:

- 📊 **Logs** - Visualização e gerenciamento de logs
- 🏥 **Health** - Monitoramento de saúde e diagnósticos
- 🔧 **Debug** - Troubleshooting e análise de problemas
- 💾 **Backup** - Backup e recuperação de dados
- ⚙️ **Config** - Gerenciamento de configuração
- 🚀 **Deployment** - Deploy e rollback controlados
- 🔧 **Maintenance** - Manutenção e otimização do sistema
- 🌐 **Network** - Diagnósticos de rede
- 📊 **Performance** - Monitoramento de performance
- 🔒 **Security** - Auditoria e segurança

## Resumo Rápido de Comandos

| Comando | Alias | Descrição | Status |
|---------|-------|-----------|--------|
| `logs` | - | Gerenciamento de logs | ✅ Implementado |
| `health` | - | Health checks e diagnósticos | ✅ Implementado |
| `debug` | - | Debug e troubleshooting | ✅ Implementado |
| `packages` | `pkg` | Gerenciamento de pacotes (APT) | ✅ Implementado |
| `services` | `svc` | Gerenciamento de serviços (systemd) | ✅ Implementado |
| `backup` | - | Backup e restore | 🔨 CLI Pronto (Implementação pendente) |
| `config` | - | Configuração do sistema | 🔨 CLI Pronto (Implementação pendente) |
| `deployment` | `deploy` | Deploy e rollback | 🔨 CLI Pronto (Implementação pendente) |
| `maintenance` | - | Manutenção do sistema | 🔨 CLI Pronto (Implementação pendente) |
| `network` | `net` | Diagnósticos de rede | 🔨 CLI Pronto (Implementação pendente) |
| `performance` | `perf` | Monitoramento de performance | 🔨 CLI Pronto (Implementação pendente) |
| `resources` | `res` | Monitoramento de recursos | 🔨 CLI Pronto (Implementação pendente) |
| `security` | - | Auditoria de segurança | 🔨 CLI Pronto (Implementação pendente) |

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

---

### 📦 packages - Gerenciamento de Pacotes

✅ **Implementado e Testado** | **[📖 Documentação Completa](sysadmin-packages.md)**

Instalar, atualizar e gerenciar pacotes do sistema (apt, yum, dnf, pacman) em agents remotos.

```bash
sloth-runner sysadmin packages [subcommand]
# Alias: sloth-runner sysadmin pkg
```

**Subcomandos:**
- `list` - Listar pacotes instalados (com filtros)
- `search` - Pesquisar pacotes nos repositórios
- `install` - Instalar pacote com confirmação interativa
- `remove` - Remover pacote (planejado)
- `update` - Atualizar listas de pacotes (apt update)
- `upgrade` - Atualizar pacotes instalados (planejado)
- `check-updates` - Verificar atualizações disponíveis (planejado)
- `info` - Mostrar informações detalhadas (planejado)
- `history` - Histórico de operações (planejado)

**Recursos Implementados:**
- ✅ Suporte completo para **APT** (Debian/Ubuntu)
- ✅ Detecção automática do gerenciador de pacotes
- ✅ Listagem com filtros e limites configuráveis
- ✅ Busca em repositórios com limite de resultados
- ✅ Instalação com confirmação interativa (--yes para auto)
- ✅ Update de listas de pacotes
- ✅ Display em tabelas formatadas com pterm
- ✅ Spinners e feedback visual durante operações
- ⏳ Suporte YUM, DNF, Pacman, APK, Zypper (planejado)

**Exemplos de Uso Real:**
```bash
# Listar todos os pacotes instalados
sloth-runner sysadmin packages list
# Saída: Tabela com Nome | Versão

# Filtrar pacotes por nome
sloth-runner sysadmin packages list --filter nginx
# Mostra apenas pacotes que contenham "nginx" no nome

# Limitar resultados
sloth-runner sysadmin pkg list --limit 50
# Mostra apenas os primeiros 50 pacotes

# Pesquisar pacote disponível
sloth-runner sysadmin packages search nginx
# Saída:
# 📦 nginx
#    High performance web server
# 📦 nginx-common
#    Common files for nginx

# Pesquisar com limite
sloth-runner sysadmin pkg search python --limit 10
# Mostra apenas os primeiros 10 resultados

# Instalar pacote (com confirmação)
sloth-runner sysadmin packages install curl
# Pergunta: Install package 'curl'? [y/n]

# Instalar sem confirmação
sloth-runner sysadmin pkg install curl --yes
# ✅ Successfully installed curl

# Atualizar listas de pacotes
sloth-runner sysadmin packages update
# ✅ Package lists updated successfully
```

**Detecção Automática:**
```bash
# O comando detecta automaticamente o gerenciador:
# 1. APT (apt-get/dpkg) - Debian, Ubuntu
# 2. YUM (yum) - CentOS, RHEL 7
# 3. DNF (dnf) - Fedora, RHEL 8+
# 4. Pacman (pacman) - Arch Linux
# 5. APK (apk) - Alpine Linux
# 6. Zypper (zypper) - openSUSE
# Retorna erro se nenhum for encontrado
```

**Roadmap:**
- ⏳ Implementar YUM, DNF, Pacman, APK, Zypper
- ⏳ Rolling updates com wait-time configurável
- ⏳ Rollback automático em falha
- ⏳ Info detalhado de pacotes (dependências, tamanho)
- ⏳ Histórico de operações com rollback

---

### 🔧 services - Gerenciamento de Serviços

✅ **Implementado e Testado** | **[📖 Documentação Completa](sysadmin-services.md)**

Controlar e monitorar serviços (systemd, init.d, OpenRC) em agents remotos.

```bash
sloth-runner sysadmin services [subcommand]
# Alias: sloth-runner sysadmin svc
```

**Subcomandos:**
- `list` - Listar todos os serviços com status colorizado
- `status` - Ver status detalhado de serviço (PID, memória, uptime)
- `start` - Iniciar serviço com verificação automática
- `stop` - Parar serviço com verificação automática
- `restart` - Reiniciar serviço com verificação de saúde
- `reload` - Recarregar configuração sem parar o serviço
- `enable` - Habilitar serviço no boot
- `disable` - Desabilitar serviço no boot
- `logs` - Ver logs do serviço (via journalctl)

**Recursos Implementados:**
- ✅ Suporte completo para **systemd** (production ready)
- ✅ Detecção automática do gerenciador de serviços
- ✅ Status colorizado e formatado (active=verde, failed=vermelho)
- ✅ Tabelas paginadas com filtros por nome e status
- ✅ Verificação automática de saúde pós-operação
- ✅ Display de PID, uso de memória e boot status
- ✅ Flags de controle (--verify, --filter, --status)
- ⏳ Suporte init.d e OpenRC (planejado)

**Exemplos de Uso Real:**
```bash
# Listar todos os serviços (tabela formatada)
sloth-runner sysadmin services list

# Filtrar serviços por nome
sloth-runner sysadmin services list --filter nginx

# Filtrar por status
sloth-runner sysadmin services list --status active

# Status detalhado com PID e memória
sloth-runner sysadmin services status nginx
# Saída:
# Service: nginx
# Status:  ● active (running)
# Enabled: yes
# PID:     1234
# Memory:  45.2M
# Since:   2 days ago

# Iniciar serviço (com verificação automática)
sloth-runner sysadmin services start nginx
# ✅ Service nginx started successfully
# ✅ Verified: nginx is active

# Parar serviço
sloth-runner sysadmin services stop nginx

# Reiniciar com verificação de saúde
sloth-runner sysadmin services restart nginx --verify

# Habilitar no boot
sloth-runner sysadmin services enable nginx
# ✅ Service nginx enabled for boot

# Ver logs em tempo real
sloth-runner sysadmin services logs nginx --follow

# Ver últimas 50 linhas de log
sloth-runner sysadmin services logs nginx -n 50
```

**Detecção Automática:**
```bash
# O comando detecta automaticamente o service manager:
# - systemd (via systemctl)
# - init.d (via service command)
# - OpenRC (via rc-service)
# - Retorna erro se nenhum for detectado
```

---

### 💾 resources - Monitoramento de Recursos

_(CLI pronto, implementação pendente)_

Monitorar CPU, memória, disco e rede em agents remotos.

```bash
sloth-runner sysadmin resources [subcommand]
# Alias: sloth-runner sysadmin res
```

**Subcomandos:**
- `overview` - Visão geral de todos recursos
- `cpu` - Uso de CPU detalhado
- `memory` - Estatísticas de memória
- `disk` - Uso de disco
- `io` - Estatísticas de I/O
- `network` - Estatísticas de rede
- `check` - Verificar contra thresholds
- `history` - Histórico de uso
- `top` - Top consumers (htop-like)

**Recursos Planejados:**
- ✨ Métricas em tempo real
- ✨ Gráficos no terminal (sparklines)
- ✨ Alertas configuráveis
- ✨ Histórico de métricas
- ✨ Exportação para Prometheus/Grafana
- ✨ Per-core CPU usage
- ✨ Análise de tendências

**Exemplos:**
```bash
# Overview de recursos
sloth-runner sysadmin resources overview --agent web-01

# CPU detalhado
sloth-runner sysadmin res cpu --agent web-01

# Verificar com alertas
sloth-runner sysadmin resources check --all-agents --alert-if cpu>80 memory>90

# Histórico de uso
sloth-runner sysadmin res history --agent web-01 --since 24h

# Top consumers
sloth-runner sysadmin resources top --agent web-01
```

---

### 💾 backup - Backup e Restore

_(CLI pronto, implementação pendente)_

Ferramentas para backup e recuperação de dados do sloth-runner.

```bash
sloth-runner sysadmin backup [subcommand]
```

**Subcomandos:**
- `create` - Criar backup completo ou incremental
- `restore` - Restaurar de backup

**Recursos Planejados:**
- ✨ Backups completos e incrementais
- ✨ Compressão e criptografia de dados
- ✨ Recuperação ponto-no-tempo
- ✨ Restore seletivo de componentes
- ✨ Verificação de integridade
- ✨ Agendamento automático

**Exemplos:**
```bash
# Criar backup completo
sloth-runner sysadmin backup create --output backup.tar.gz

# Restaurar de backup
sloth-runner sysadmin backup restore --input backup.tar.gz
```

---

### ⚙️ config - Gerenciamento de Configuração

_(CLI pronto, implementação pendente)_

Gerenciar, validar e sincronizar configurações do sloth-runner.

```bash
sloth-runner sysadmin config [subcommand]
```

**Subcomandos:**
- `validate` - Validar arquivos de configuração
- `diff` - Comparar configurações entre agents
- `export` - Exportar configuração atual
- `import` - Importar configuração de arquivo
- `set` - Alterar valor de configuração dinamicamente
- `get` - Obter valor de configuração
- `reset` - Resetar configuração para padrões

**Recursos Planejados:**
- ✨ Validação de sintaxe YAML/JSON
- ✨ Comparação side-by-side entre agents
- ✨ Hot reload sem restart
- ✨ Backup automático antes de mudanças
- ✨ Template de configuração
- ✨ Versionamento de configuração

**Exemplos:**
```bash
# Validar configuração
sloth-runner sysadmin config validate

# Comparar entre agents
sloth-runner sysadmin config diff --agents do-sloth-runner-01,do-sloth-runner-02

# Alterar dinamicamente
sloth-runner sysadmin config set --key log.level --value debug

# Exportar para arquivo
sloth-runner sysadmin config export --output config.yaml
```

---

### 🚀 deployment - Deploy e Rollback

_(CLI pronto, implementação pendente)_

Ferramentas para deployment controlado e rollback de atualizações.

```bash
sloth-runner sysadmin deployment [subcommand]
# Alias: sloth-runner sysadmin deploy
```

**Subcomandos:**
- `deploy` - Fazer deploy de atualização
- `rollback` - Reverter para versão anterior

**Recursos Planejados:**
- ✨ Rolling updates progressivos
- ✨ Canary deployments
- ✨ Blue-green deployments
- ✨ One-click rollback
- ✨ Histórico de versões
- ✨ Verificações de segurança

**Exemplos:**
```bash
# Deploy para production
sloth-runner sysadmin deployment deploy --env production --strategy rolling

# Rollback rápido
sloth-runner sysadmin deploy rollback --version v1.2.3
```

---

### 🔧 maintenance - Manutenção do Sistema

_(CLI pronto, implementação pendente)_

Ferramentas de manutenção, limpeza e otimização do sistema.

```bash
sloth-runner sysadmin maintenance [subcommand]
```

**Subcomandos:**
- `clean-logs` - Limpar e rotacionar logs antigos
- `optimize-db` - Otimizar banco de dados (VACUUM, ANALYZE)
- `cleanup` - Limpeza geral (temp files, cache)

**Recursos Planejados:**
- ✨ Rotação automática de logs
- ✨ Compressão de arquivos antigos
- ✨ Otimização de banco com VACUUM e ANALYZE
- ✨ Reconstrução de índices
- ✨ Limpeza de arquivos temporários
- ✨ Detecção de arquivos órfãos
- ✨ Limpeza de cache

**Exemplos:**
```bash
# Limpar logs antigos
sloth-runner sysadmin maintenance clean-logs --older-than 30d

# Otimizar banco de dados
sloth-runner sysadmin maintenance optimize-db --full

# Limpeza geral
sloth-runner sysadmin maintenance cleanup --dry-run
```

---

### 🌐 network - Diagnósticos de Rede

_(CLI pronto, implementação pendente)_

Ferramentas para testar conectividade e diagnosticar problemas de rede.

```bash
sloth-runner sysadmin network [subcommand]
# Alias: sloth-runner sysadmin net
```

**Subcomandos:**
- `ping` - Testar conectividade com agents
- `port-check` - Verificar disponibilidade de portas

**Recursos Planejados:**
- ✨ Teste de conectividade entre nodes
- ✨ Medição de latência
- ✨ Detecção de packet loss
- ✨ Scan de portas
- ✨ Detecção de serviços
- ✨ Teste de firewall rules

**Exemplos:**
```bash
# Testar conectividade
sloth-runner sysadmin network ping --agent do-sloth-runner-01

# Verificar portas
sloth-runner sysadmin net port-check --agent do-sloth-runner-01 --ports 50051,22,80
```

---

### 📊 performance - Monitoramento de Performance

_(CLI pronto, implementação pendente)_

Monitorar e analisar performance do sistema e agents.

```bash
sloth-runner sysadmin performance [subcommand]
# Alias: sloth-runner sysadmin perf
```

**Subcomandos:**
- `show` - Exibir métricas de performance
- `monitor` - Monitoramento em tempo real

**Recursos Planejados:**
- ✨ Uso de CPU por agent
- ✨ Estatísticas de memória
- ✨ I/O de disco
- ✨ Throughput de rede
- ✨ Dashboards ao vivo
- ✨ Thresholds de alerta
- ✨ Tendências históricas

**Exemplos:**
```bash
# Ver métricas atuais
sloth-runner sysadmin performance show --agent do-sloth-runner-01

# Monitoramento contínuo
sloth-runner sysadmin perf monitor --interval 5s --all-agents
```

---

### 🔒 security - Auditoria de Segurança

_(CLI pronto, implementação pendente)_

Ferramentas para auditoria de segurança e scanning de vulnerabilidades.

```bash
sloth-runner sysadmin security [subcommand]
```

**Subcomandos:**
- `audit` - Auditar logs de segurança
- `scan` - Scan de vulnerabilidades

**Recursos Planejados:**
- ✨ Análise de logs de acesso
- ✨ Detecção de tentativas de autenticação falhadas
- ✨ Identificação de atividade suspeita
- ✨ Scanning de CVEs
- ✨ Auditoria de dependências
- ✨ Validação de configurações de segurança

**Exemplos:**
```bash
# Auditoria de segurança
sloth-runner sysadmin security audit --since 24h --show-failed-auth

# Scan de vulnerabilidades
sloth-runner sysadmin security scan --agent do-sloth-runner-01 --full
```

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

# Ver performance dos agents
sloth-runner sysadmin perf show --all-agents

# Validar configuração
sloth-runner sysadmin config validate
```

### 2. Troubleshooting de Problema

```bash
# 1. Check de saúde geral
sloth-runner sysadmin health check --verbose

# 2. Buscar erros recentes
sloth-runner sysadmin logs search --query "error" --since 1h

# 3. Verificar agent específico
sloth-runner sysadmin health agent problematic-agent

# 4. Testar conectividade de rede
sloth-runner sysadmin net ping --agent problematic-agent

# 5. Verificar performance
sloth-runner sysadmin perf show --agent problematic-agent

# 6. Auditoria de segurança
sloth-runner sysadmin security audit --agent problematic-agent --since 24h

# 7. Gerar diagnóstico para análise
sloth-runner sysadmin health diagnostics --output issue-$(date +%Y%m%d).json
```

### 3. Manutenção e Arquivamento

```bash
# 1. Verificar espaço e logs
sloth-runner sysadmin health check
ls -lh /etc/sloth-runner/logs/

# 2. Backup completo
sloth-runner sysadmin backup create --output backup-$(date +%Y%m%d).tar.gz

# 3. Exportar logs para backup
sloth-runner sysadmin logs export --format json --since 30d --output logs-backup.json

# 4. Limpar logs antigos
sloth-runner sysadmin maintenance clean-logs --older-than 30d

# 5. Rotacionar logs
sloth-runner sysadmin logs rotate --force

# 6. Otimizar banco de dados
sloth-runner sysadmin maintenance optimize-db --full

# 7. Limpeza geral
sloth-runner sysadmin maintenance cleanup

# 8. Verificar saúde pós-manutenção
sloth-runner sysadmin health check
```

### 4. Monitoramento Contínuo

```bash
# Terminal 1: Health monitoring
sloth-runner sysadmin health watch --interval 30s

# Terminal 2: Performance monitoring
sloth-runner sysadmin perf monitor --interval 10s --all-agents

# Terminal 3: Log monitoring
sloth-runner sysadmin logs tail --follow --level warn

# Terminal 4: Network monitoring
watch -n 30 'sloth-runner sysadmin net ping --all-agents'

# Terminal 5: Operações
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

---

## Novos Comandos Sysadmin (v2.0+)

### Visão Geral das Novas Ferramentas

A versão 2.0 do sloth-runner introduz 7 novos comandos sysadmin que ampliam significativamente as capacidades de administração:

#### 1. Config Management 🆕
- Validação automática de configurações
- Comparação entre múltiplos agents
- Hot reload sem downtime
- Export/import de configurações

#### 2. Performance Monitoring 🆕
- Métricas em tempo real de CPU, memória e I/O
- Histórico de tendências
- Alertas de threshold
- Dashboards interativos

#### 3. Network Diagnostics 🆕
- Testes de conectividade automatizados
- Medição de latência entre nodes
- Port scanning e service detection
- Análise de firewall rules

#### 4. Security Auditing 🆕
- Auditoria de logs de acesso
- Detecção de atividade anômala
- Scanning de vulnerabilidades (CVE)
- Auditoria de dependências

#### 5. Automated Backups 🆕
- Backups completos e incrementais
- Criptografia de dados sensíveis
- Point-in-time recovery
- Restore seletivo

#### 6. Maintenance Tools 🆕
- Limpeza automática de logs
- Otimização de banco de dados
- Detecção de arquivos órfãos
- Cache management

#### 7. Deployment Management 🆕
- Rolling updates controlados
- Canary deployments
- Blue-green deployments
- Rollback com um clique

### Roadmap de Implementação

**Fase 1 - Q1 2025** ✅
- Estrutura base dos comandos
- Testes unitários (83.7% coverage)
- Documentação completa
- CLI interfaces

**Fase 2 - Q2 2025** 🚧
- Implementação do config management
- Performance monitoring básico
- Network diagnostics essenciais

**Fase 3 - Q3 2025** 📋
- Security auditing completo
- Backup automation
- Maintenance tools

**Fase 4 - Q4 2025** 📋
- Deployment management avançado
- Integração com ferramentas externas
- Dashboard web completo

### Começando a Usar

Todos os novos comandos seguem a mesma estrutura:

```bash
sloth-runner sysadmin [comando] [subcomando] [flags]
```

**Exemplos:**
```bash
# Config
sloth-runner sysadmin config validate
sloth-runner sysadmin config diff --agents agent1,agent2

# Performance (com alias)
sloth-runner sysadmin performance show
sloth-runner sysadmin perf monitor --interval 10s

# Network (com alias)
sloth-runner sysadmin network ping --agent web-01
sloth-runner sysadmin net port-check --ports 80,443

# Security
sloth-runner sysadmin security audit --since 24h
sloth-runner sysadmin security scan --full

# Backup
sloth-runner sysadmin backup create --output backup.tar.gz
sloth-runner sysadmin backup restore --input backup.tar.gz

# Maintenance
sloth-runner sysadmin maintenance clean-logs --older-than 30d
sloth-runner sysadmin maintenance optimize-db --full

# Deployment (com alias)
sloth-runner sysadmin deployment deploy --strategy rolling
sloth-runner sysadmin deploy rollback --version v1.2.3
```

### Arquitetura dos Novos Comandos

```
cmd/sloth-runner/commands/sysadmin/
├── sysadmin.go          # Comando principal
├── backup/
│   ├── backup.go        # Lógica de backup
│   └── backup_test.go   # Testes (100% coverage)
├── config/
│   ├── config.go        # Gestão de config
│   └── config_test.go   # Testes (73.9% coverage)
├── deployment/
│   ├── deployment.go    # Deploy/rollback
│   └── deployment_test.go
├── maintenance/
│   ├── maintenance.go   # Manutenção
│   └── maintenance_test.go
├── network/
│   ├── network.go       # Diagnósticos de rede
│   └── network_test.go  # Testes (100% coverage)
├── performance/
│   ├── performance.go   # Monitoramento
│   └── performance_test.go
└── security/
    ├── security.go      # Segurança
    └── security_test.go
```

### Status de Testes

Todos os novos comandos possuem testes abrangentes:

| Comando | Testes | Coverage | Status |
|---------|--------|----------|--------|
| backup | 6 testes | 100% | ✅ |
| config | 9 testes | 73.9% | ✅ |
| deployment | 5 testes | 63.6% | ✅ |
| maintenance | 7 testes | 66.7% | ✅ |
| network | 6 testes | 100% | ✅ |
| performance | 6 testes | 100% | ✅ |
| security | 4 testes | 75% | ✅ |
| **Total** | **43 testes** | **83.7%** | ✅ |

**Benchmarks:**
- Tempo médio de execução: < 1µs
- Alocações de memória: 2-53 KB
- Performance otimizada para produção

### Contribuindo

Os novos comandos são projetados para serem extensíveis. Para adicionar funcionalidade:

1. Adicione lógica em `cmd/sloth-runner/commands/sysadmin/[comando]/`
2. Escreva testes unitários
3. Atualize documentação
4. Submeta PR com coverage > 70%

### Feedback e Sugestões

Estamos desenvolvendo ativamente estes comandos. Se você tem sugestões ou precisa de funcionalidades específicas:

- Abra uma issue no GitHub
- Entre em contato via Slack
- Contribua com PRs

---

## Ver Também

- [Agent Management](../agent.md) - Gerenciar agents
- [Workflow Execution](../workflow.md) - Executar workflows
- [Master Server](../master.md) - Servidor master
- [CLI Reference](../cli-reference.md) - Referência completa de comandos
- [Logs Command](logs-command.md) - Documentação detalhada de logs
- [Health Command](health-command.md) - Documentação detalhada de health
- [Debug Command](debug-command.md) - Documentação detalhada de debug
