# Comando Sysadmin Logs

O comando `sysadmin logs` fornece ferramentas abrangentes para gerenciamento e visualização de logs do sloth-runner master e agents.

## Visão Geral

```bash
sloth-runner sysadmin logs [subcommand] [flags]
```

> **Nota:** Este comando faz parte da suite de ferramentas `sysadmin`. Para ver todos os comandos disponíveis para administradores de sistema, use `sloth-runner sysadmin --help`.

## Subcomandos

### 1. tail - Visualizar Logs em Tempo Real

Exibe as últimas N linhas de logs e opcionalmente acompanha novas entradas em tempo real.

**Uso:**
```bash
sloth-runner sysadmin logs tail [flags]
```

**Flags:**
- `-f, --follow` - Acompanhar logs em tempo real (como tail -f)
- `-n, --lines int` - Número de linhas a exibir (padrão: 10)
- `-a, --agent string` - Filtrar por nome do agent
- `-l, --level string` - Filtrar por nível de log (debug, info, warn, error)

**Exemplos:**
```bash
# Mostrar últimas 10 linhas
sloth-runner sysadmin logs tail

# Acompanhar logs em tempo real
sloth-runner sysadmin logs tail --follow

# Filtrar por agent específico
sloth-runner sysadmin logs tail --agent do-sloth-runner-01 --follow

# Filtrar por nível de erro
sloth-runner sysadmin logs tail --level error --follow

# Combinar filtros
sloth-runner sysadmin logs tail --agent web-01 --level warn -n 20
```

**Casos de Uso:**
- Monitoramento em tempo real de agents
- Debugging de problemas específicos
- Acompanhamento de deployments
- Identificação rápida de erros

---

### 2. search - Buscar nos Logs

Busca através dos logs usando queries de texto e filtros avançados.

**Uso:**
```bash
sloth-runner sysadmin logs search [flags]
```

**Flags:**
- `-q, --query string` - Query de busca (case-insensitive)
- `--since string` - Logs desde (ex: 1h, 30m, 24h, 7d)
- `--until string` - Logs até (ex: 1h, 30m)
- `-a, --agent string` - Filtrar por nome do agent
- `-l, --level string` - Filtrar por nível de log
- `--limit int` - Número máximo de resultados (padrão: 100)

**Exemplos:**
```bash
# Buscar por erros na última hora
sloth-runner sysadmin logs search --query "error" --since 1h

# Buscar em agent específico
sloth-runner sysadmin logs search --query "failed" --agent web-01

# Buscar com intervalo de tempo
sloth-runner sysadmin logs search --query "timeout" --since 2h --until 1h

# Buscar erros de conexão nos últimos 7 dias
sloth-runner sysadmin logs search --query "connection" --level error --since 7d

# Limitar resultados
sloth-runner sysadmin logs search --query "deploy" --limit 50
```

**Formatos de Tempo:**
- `1h` - 1 hora atrás
- `30m` - 30 minutos atrás
- `24h` - 24 horas atrás (1 dia)
- `7d` - 7 dias atrás

**Casos de Uso:**
- Investigação de incidentes históricos
- Análise de padrões de erro
- Auditoria de operações específicas
- Troubleshooting pós-mortem

---

### 3. export - Exportar Logs

Exporta logs em vários formatos para análise externa ou arquivamento.

**Uso:**
```bash
sloth-runner sysadmin logs export [flags]
```

**Flags:**
- `-f, --format string` - Formato de saída: text, json, csv (padrão: text)
- `-o, --output string` - Arquivo de saída (stdout se não especificado)
- `--since string` - Exportar logs desde (ex: 1h, 24h, 7d)
- `-a, --agent string` - Filtrar por nome do agent

**Exemplos:**
```bash
# Exportar para JSON
sloth-runner sysadmin logs export --format json --output logs.json

# Exportar para CSV com filtro de tempo
sloth-runner sysadmin logs export --format csv --output logs.csv --since 24h

# Exportar logs de agent específico
sloth-runner sysadmin logs export --agent web-01 --output web-01.log

# Exportar últimos 7 dias em JSON
sloth-runner sysadmin logs export --format json --since 7d --output weekly-logs.json

# Exportar para stdout (pipe para outras ferramentas)
sloth-runner sysadmin logs export --format csv | grep ERROR
```

**Formatos de Saída:**

**TEXT (padrão):**
```
[2025-10-09 10:00:01] INFO : Starting sloth-runner master server on port 50053
[2025-10-09 10:00:05] INFO do-sloth-runner-01: Heartbeat received
```

**JSON:**
```json
[
  {
    "timestamp": "2025-10-09T10:00:01Z",
    "level": "INFO",
    "agent": "",
    "message": "Starting sloth-runner master server on port 50053"
  }
]
```

**CSV:**
```csv
Timestamp,Level,Agent,Message
2025-10-09T10:00:01Z,INFO,,"Starting sloth-runner master server on port 50053"
2025-10-09T10:00:05Z,INFO,do-sloth-runner-01,"Heartbeat received"
```

**Casos de Uso:**
- Importação para sistemas de análise (ELK, Splunk)
- Geração de relatórios
- Backup de logs
- Análise em Excel/Python/R

---

### 4. rotate - Rotacionar Logs

Rotaciona manualmente os arquivos de log para arquivamento.

**Uso:**
```bash
sloth-runner sysadmin logs rotate [flags]
```

**Flags:**
- `-f, --force` - Forçar rotação mesmo se arquivo for pequeno

**Exemplos:**
```bash
# Rotacionar logs (apenas se > 10MB)
sloth-runner sysadmin logs rotate

# Forçar rotação independente do tamanho
sloth-runner sysadmin logs rotate --force
```

**Comportamento:**
- Por padrão, só rotaciona se o arquivo for maior que 10MB
- Arquivo original é renomeado com timestamp: `sloth-runner.log.20251009-150405`
- Novo arquivo vazio é criado automaticamente
- Logs rotacionados devem ser arquivados/comprimidos manualmente

**Casos de Uso:**
- Manutenção preventiva de espaço em disco
- Preparação para backup
- Limpeza antes de debugging intensivo
- Separação de logs por período

---

### 5. level - Alterar Nível de Log

Altera dinamicamente o nível de logging do servidor master.

**Uso:**
```bash
sloth-runner sysadmin logs level [debug|info|warn|error]
```

**Níveis de Log:**
- `debug` - Logs detalhados para debugging
- `info` - Informações operacionais normais (padrão)
- `warn` - Avisos e problemas não críticos
- `error` - Apenas erros críticos

**Exemplos:**
```bash
# Ativar modo debug
sloth-runner sysadmin logs level debug

# Voltar para info
sloth-runner sysadmin logs level info

# Apenas erros
sloth-runner sysadmin logs level error
```

**Nota:** Esta funcionalidade requer API do master (em desenvolvimento). Atualmente mostra uma mensagem informativa.

**Casos de Uso:**
- Debugging temporário sem reiniciar serviço
- Redução de ruído em produção
- Troubleshooting de problemas específicos

---

### 6. remote - Buscar Logs de Agents Remotos

Busca logs do sistema operacional de agents remotos via gRPC, sem necessidade de SSH interativo.

**Uso:**
```bash
sloth-runner sysadmin logs remote [flags]
```

**Flags:**
- `-a, --agent string` - Nome do agent (obrigatório)
- `-p, --path string` - Caminho customizado do arquivo de log
- `-s, --system string` - Tipo de log do sistema (syslog, messages, journalctl, kern, auth)
- `-n, --lines int` - Número de linhas a exibir (padrão: 50)
- `-f, --follow` - Acompanhar saída de log em tempo real

**Tipos de Logs do Sistema:**
- `syslog` - Log geral do sistema (/var/log/syslog ou /var/log/messages)
- `messages` - Log de mensagens do sistema (/var/log/messages)
- `journalctl` - Logs do systemd journal
- `kern` - Logs do kernel (/var/log/kern.log)
- `auth` - Logs de autenticação (/var/log/auth.log ou /var/log/secure)

**Exemplos:**
```bash
# Buscar syslog de um agent
sloth-runner sysadmin logs remote --agent do-sloth-runner-01 --system syslog

# Buscar journalctl
sloth-runner sysadmin logs remote --agent web-01 --system journalctl --lines 100

# Buscar arquivo de log customizado
sloth-runner sysadmin logs remote --agent app-01 --path /var/log/nginx/error.log

# Buscar logs de autenticação
sloth-runner sysadmin logs remote --agent db-01 --system auth --lines 20

# Buscar logs do kernel
sloth-runner sysadmin logs remote --agent web-02 --system kern --lines 30

# Acompanhar logs em tempo real (follow mode)
sloth-runner sysadmin logs remote --agent app-01 --system syslog --follow
```

**Comportamento:**
- Conecta ao agent via gRPC (porta configurada no registro)
- Executa comandos `tail`, `journalctl` no agent remoto
- Retorna logs do sistema operacional do agent
- Suporta modo follow para monitoramento contínuo
- Não requer configuração SSH adicional
- Funciona automaticamente se o agent estiver registrado

**Casos de Uso:**
- Debugging de problemas no sistema operacional do agent
- Monitoramento de logs de serviços no agent (nginx, apache, etc.)
- Análise de logs de autenticação e segurança
- Troubleshooting de kernel e hardware
- Coleta centralizada de logs sem ferramentas externas
- Investigação de incidentes em agents remotos

**Comparação com Outros Comandos:**
- `tail`: Logs locais do sloth-runner master
- `remote`: Logs do sistema operacional dos agents remotos
- `search`: Busca em logs locais já coletados
- `export`: Exporta logs locais para análise

**Vantagens:**
- ✅ Sem necessidade de SSH interativo
- ✅ Sem configuração adicional de chaves SSH
- ✅ Funciona através do firewall (usa porta gRPC já aberta)
- ✅ Controle de acesso via registro de agents
- ✅ Suporta múltiplos tipos de logs
- ✅ Modo follow para monitoramento contínuo

**Limitações:**
- Requer agent ativo e registrado
- Limitado aos logs acessíveis pelo usuário do agent
- Não suporta rotação remota (use comandos específicos do SO)

---

## Localização dos Logs

Os logs são armazenados em:
- **Linux/macOS (root):** `/etc/sloth-runner/logs/sloth-runner.log`
- **Linux/macOS (user):** `~/.sloth-runner/logs/sloth-runner.log`
- **Windows:** `C:\ProgramData\sloth-runner\logs\sloth-runner.log`

Pode ser customizado via variável de ambiente:
```bash
export SLOTH_RUNNER_DATA_DIR=/custom/path
```

---

## Formato dos Logs

Formato padrão das entradas de log:

```
YYYY-MM-DD HH:MM:SS LEVEL [agent=AGENT_NAME] Message
```

**Exemplo:**
```
2025-10-09 10:00:05 INFO agent=do-sloth-runner-01 Heartbeat received from agent
2025-10-09 10:00:10 WARN agent=do-sloth-runner-01 High CPU usage detected: 85%
2025-10-09 10:00:15 ERROR agent=do-sloth-runner-02 Connection timeout after 5 seconds
```

---

## Workflows Comuns

### Debugging de Problema Recente
```bash
# 1. Ver últimos erros
sloth-runner sysadmin logs tail --level error -n 50

# 2. Buscar padrão específico
sloth-runner sysadmin logs search --query "timeout" --since 1h

# 3. Acompanhar em tempo real
sloth-runner sysadmin logs tail --follow --level error
```

### Análise de Agent Específico
```bash
# 1. Exportar logs do agent
sloth-runner sysadmin logs export --agent web-01 --since 24h --output web-01.log

# 2. Contar erros
sloth-runner sysadmin logs search --agent web-01 --level error --since 24h

# 3. Monitorar em tempo real
sloth-runner sysadmin logs tail --agent web-01 --follow
```

### Manutenção e Arquivamento
```bash
# 1. Verificar tamanho atual
ls -lh /etc/sloth-runner/logs/

# 2. Exportar para backup
sloth-runner sysadmin logs export --format json --since 30d --output backup.json

# 3. Rotacionar logs
sloth-runner sysadmin logs rotate --force

# 4. Comprimir arquivo antigo
gzip /etc/sloth-runner/logs/sloth-runner.log.20251009-*
```

### Investigação de Incidente
```bash
# 1. Identificar janela de tempo do problema
sloth-runner sysadmin logs search --query "error" --since 4h --until 2h

# 2. Focar em agent com problema
sloth-runner sysadmin logs search --agent problematic-agent --since 4h --until 2h

# 3. Exportar para análise detalhada
sloth-runner sysadmin logs export --agent problematic-agent --since 4h --format json --output incident.json

# 4. Analisar com jq
cat incident.json | jq '.[] | select(.level=="ERROR")'
```

### Debugging de Agent Remoto
```bash
# 1. Ver logs do sistema do agent
sloth-runner sysadmin logs remote --agent web-01 --system syslog --lines 100

# 2. Verificar autenticação
sloth-runner sysadmin logs remote --agent web-01 --system auth --lines 50

# 3. Verificar logs de serviço específico
sloth-runner sysadmin logs remote --agent web-01 --path /var/log/nginx/error.log --lines 50

# 4. Monitorar em tempo real
sloth-runner sysadmin logs remote --agent web-01 --system syslog --follow
```

### Análise de Segurança em Múltiplos Agents
```bash
# Verificar tentativas de login em todos os agents
for agent in web-01 web-02 db-01; do
  echo "=== $agent ==="
  sloth-runner sysadmin logs remote --agent $agent --system auth --lines 20 | grep -i "failed\|invalid"
done

# Verificar logs do kernel em busca de problemas
sloth-runner sysadmin logs remote --agent problem-server --system kern --lines 100
```

---

## Integração com Outras Ferramentas

### Pipeline com grep/awk/sed
```bash
# Contar erros por agent
sloth-runner sysadmin logs export --format csv | grep ERROR | awk -F',' '{print $3}' | sort | uniq -c

# Extrair timestamps de erros
sloth-runner sysadmin logs export | grep ERROR | awk '{print $1, $2}'
```

### Análise com jq (JSON)
```bash
# Top 5 agents com mais logs
sloth-runner sysadmin logs export --format json | jq -r '.[].agent' | sort | uniq -c | sort -rn | head -5

# Erros por hora
sloth-runner sysadmin logs export --format json | jq -r 'select(.level=="ERROR") | .timestamp' | cut -d'T' -f2 | cut -d':' -f1 | sort | uniq -c
```

### Watch para Monitoramento Contínuo
```bash
# Atualizar contagem de erros a cada 10 segundos
watch -n 10 'sloth-runner sysadmin logs search --query "error" --since 1h | tail -1'
```

---

## Troubleshooting

### Logs Não Aparecem
```bash
# Verificar se diretório existe
ls -la /etc/sloth-runner/logs/

# Verificar permissões
ls -la /etc/sloth-runner/logs/sloth-runner.log

# Criar diretório se necessário
sudo mkdir -p /etc/sloth-runner/logs
sudo chown $USER /etc/sloth-runner/logs
```

### Performance com Logs Grandes
```bash
# Usar limit para reduzir resultados
sloth-runner sysadmin logs search --query "error" --limit 100

# Filtrar por tempo recente
sloth-runner sysadmin logs search --query "error" --since 1h

# Rotacionar logs grandes
sloth-runner sysadmin logs rotate --force
```

### Erro de Formato de Tempo
```bash
# Formato correto: número + unidade (h/m/d)
sloth-runner sysadmin logs search --since 2h    # ✓ Correto
sloth-runner sysadmin logs search --since 2     # ✗ Errado
sloth-runner sysadmin logs search --since 2hr   # ✗ Errado
```

---

## Boas Práticas

1. **Rotação Regular:** Configure rotação automática ou execute semanalmente
2. **Arquivamento:** Comprima e arquive logs antigos (>30 dias)
3. **Monitoramento:** Use `logs tail --follow` em terminais dedicados
4. **Filtros Específicos:** Sempre use filtros (agent, level, since) para melhor performance
5. **Backup:** Exporte logs críticos antes de rotacionar
6. **Análise:** Use format JSON para análises complexas com jq/Python

---

## Ver Também

- [health](health-command.md) - Verificar saúde do sistema
- [agent](../agent.md) - Gerenciar agents
- [Troubleshooting Guide](../troubleshooting.md) - Guia de solução de problemas
