# Agent Groups - Gestão de Grupos de Agentes

O sistema de Agent Groups (Grupos de Agentes) permite organizar e gerenciar múltiplos agentes de forma eficiente, facilitando operações em escala.

## Índice

- [Visão Geral](#visão-geral)
- [Comandos CLI](#comandos-cli)
  - [Gerenciamento Básico](#gerenciamento-básico)
  - [Operações em Massa](#operações-em-massa)
  - [Templates](#templates)
  - [Auto-Discovery](#auto-discovery)
  - [Webhooks](#webhooks)
- [Interface Web](#interface-web)
- [API REST](#api-rest)
- [Casos de Uso](#casos-de-uso)

## Visão Geral

Agent Groups oferece as seguintes funcionalidades:

- **Agrupamento Lógico**: Organize agentes por função, ambiente, região, etc.
- **Operações em Massa**: Execute comandos em múltiplos agentes simultaneamente
- **Templates**: Crie grupos reutilizáveis com regras de matching
- **Auto-Discovery**: Descubra e adicione agentes automaticamente baseado em regras
- **Webhooks**: Receba notificações de eventos de grupos
- **Hierarquia**: Organize grupos em estruturas hierárquicas

## Comandos CLI

### Gerenciamento Básico

#### Listar Grupos

```bash
# Listar todos os grupos (formato tabela)
sloth-runner group list

# Listar em formato JSON
sloth-runner group list --output json
```

**Saída exemplo:**
```
NAME              AGENTS  DESCRIPTION                 TAGS
----              ------  -----------                 ----
production-web    5       Production web servers      env=production,role=web
staging-db        2       Staging database servers    env=staging,role=database
monitoring        3       Monitoring agents           role=monitoring
```

#### Criar Grupo

```bash
# Criar grupo básico
sloth-runner group create production-web

# Criar com descrição
sloth-runner group create production-web \
  --description "Production web servers"

# Criar com tags
sloth-runner group create production-web \
  --description "Production web servers" \
  --tag environment=production \
  --tag role=webserver \
  --tag region=us-east-1
```

#### Visualizar Detalhes de um Grupo

```bash
# Visualizar em formato texto
sloth-runner group show production-web

# Visualizar em formato JSON
sloth-runner group show production-web --output json
```

**Saída exemplo:**
```
Group: production-web
Description: Production web servers
Agent Count: 5

Tags:
  environment: production
  role: webserver
  region: us-east-1

Agents:
  • server-01
  • server-02
  • server-03
  • server-04
  • server-05
```

#### Deletar Grupo

```bash
# Deletar com confirmação
sloth-runner group delete production-web

# Deletar sem confirmação (force)
sloth-runner group delete production-web --force
```

#### Adicionar Agentes ao Grupo

```bash
# Adicionar um agente
sloth-runner group add-agent production-web server-01

# Adicionar múltiplos agentes
sloth-runner group add-agent production-web server-01 server-02 server-03
```

#### Remover Agentes do Grupo

```bash
# Remover um agente
sloth-runner group remove-agent production-web server-01

# Remover múltiplos agentes
sloth-runner group remove-agent production-web server-01 server-02
```

### Operações em Massa

Execute operações em todos os agentes de um grupo simultaneamente.

#### Operações Disponíveis

- `restart` - Reiniciar todos os agentes
- `update` - Atualizar todos os agentes
- `shutdown` - Desligar todos os agentes
- `execute` - Executar comando customizado

```bash
# Reiniciar todos os agentes do grupo
sloth-runner group bulk production-web restart

# Atualizar todos os agentes com timeout customizado
sloth-runner group bulk production-web update --timeout 600

# Executar comando customizado
sloth-runner group bulk production-web execute \
  --command "systemctl restart nginx"

# Executar comando com timeout
sloth-runner group bulk production-web execute \
  --command "apt-get update && apt-get upgrade -y" \
  --timeout 900
```

**Saída exemplo:**
```
⏳ Executing 'restart' operation on group 'production-web'...

✅ Bulk operation completed in 3420ms
Summary: 5/5 agents succeeded (100.0%)

AGENT       STATUS       DURATION  OUTPUT/ERROR
-----       ------       --------  ------------
server-01   ✅ SUCCESS   650ms     Agent restarted successfully
server-02   ✅ SUCCESS   720ms     Agent restarted successfully
server-03   ✅ SUCCESS   680ms     Agent restarted successfully
server-04   ✅ SUCCESS   710ms     Agent restarted successfully
server-05   ✅ SUCCESS   660ms     Agent restarted successfully
```

### Templates

Templates permitem criar grupos reutilizáveis com regras de matching automático.

#### Listar Templates

```bash
# Listar todos os templates
sloth-runner group template list

# Listar em formato JSON
sloth-runner group template list --output json
```

#### Criar Template

```bash
# Criar template com regra de tag match
sloth-runner group template create web-servers \
  --description "Web server template" \
  --rule "tag_match:equals:web" \
  --tag "env:production"

# Criar template com múltiplas regras
sloth-runner group template create prod-db \
  --description "Production database template" \
  --rule "tag_match:equals:database" \
  --rule "name_pattern:contains:prod" \
  --rule "status:equals:active"

# Criar template com regex
sloth-runner group template create monitoring-agents \
  --rule "name_pattern:regex:^monitor-.*$" \
  --tag "role:monitoring"
```

**Tipos de Regras:**
- `tag_match` - Match baseado em tags do agente
- `name_pattern` - Match baseado no nome do agente
- `status` - Match baseado no status do agente

**Operadores:**
- `equals` - Igualdade exata
- `contains` - Contém substring
- `regex` - Expressão regular

#### Aplicar Template

```bash
# Aplicar template para criar/atualizar grupo
sloth-runner group template apply web-servers production-web
```

**Saída exemplo:**
```
✅ Template applied successfully to group 'production-web'
   Matched agents: 5
```

#### Deletar Template

```bash
# Deletar com confirmação
sloth-runner group template delete web-servers

# Deletar sem confirmação
sloth-runner group template delete web-servers --force
```

### Auto-Discovery

Configure regras para descobrir e adicionar agentes automaticamente aos grupos.

#### Listar Configurações

```bash
# Listar todas as configurações de auto-discovery
sloth-runner group auto-discovery list

# Formato JSON
sloth-runner group auto-discovery list --output json
```

**Saída exemplo:**
```
ID              NAME            GROUP            SCHEDULE        ENABLED  RULES
--              ----            -----            --------        -------  -----
web-disc        web-discovery   production-web   */10 * * * *    Yes      2
db-disc         db-discovery    production-db    0 * * * *       Yes      1
```

#### Criar Configuração

```bash
# Criar auto-discovery para web servers (a cada 10 minutos)
sloth-runner group auto-discovery create web-discovery \
  --group production-web \
  --schedule "*/10 * * * *" \
  --rule "tag_match:equals:web" \
  --rule "status:equals:active" \
  --enabled

# Criar para database servers (a cada hora)
sloth-runner group auto-discovery create db-discovery \
  --group production-db \
  --schedule "0 * * * *" \
  --rule "tag_match:equals:database" \
  --rule "name_pattern:contains:db" \
  --tag "auto_discovered:true"
```

**Formato do Schedule:** Expressão cron (minuto hora dia mês dia-da-semana)
- `*/5 * * * *` - A cada 5 minutos
- `0 * * * *` - A cada hora
- `0 0 * * *` - Diariamente à meia-noite
- `0 0 * * 0` - Semanalmente aos domingos

#### Executar Manualmente

```bash
# Executar auto-discovery manualmente
sloth-runner group auto-discovery run web-discovery
```

**Saída exemplo:**
```
✅ Auto-discovery run completed
   Matched agents: 3
```

#### Habilitar/Desabilitar

```bash
# Habilitar configuração
sloth-runner group auto-discovery enable web-discovery

# Desabilitar configuração
sloth-runner group auto-discovery disable web-discovery
```

#### Deletar Configuração

```bash
# Deletar com confirmação
sloth-runner group auto-discovery delete web-discovery

# Deletar sem confirmação
sloth-runner group auto-discovery delete web-discovery --force
```

### Webhooks

Configure webhooks para receber notificações de eventos de grupos.

#### Listar Webhooks

```bash
# Listar todos os webhooks
sloth-runner group webhook list

# Formato JSON
sloth-runner group webhook list --output json
```

**Saída exemplo:**
```
ID              NAME             URL                                      EVENTS  ENABLED
--              ----             ---                                      ------  -------
slack-1         slack-notify     https://hooks.slack.com/services/...     3       Yes
discord-1       discord-webhook  https://discord.com/api/webhooks/...     2       Yes
```

#### Criar Webhook

```bash
# Webhook para Slack
sloth-runner group webhook create slack-notify \
  --url "https://hooks.slack.com/services/YOUR/WEBHOOK/URL" \
  --event "group.created" \
  --event "group.deleted" \
  --event "bulk.operation_end" \
  --enabled

# Webhook com secret e headers customizados
sloth-runner group webhook create discord-webhook \
  --url "https://discord.com/api/webhooks/YOUR/WEBHOOK" \
  --event "group.agent_added" \
  --event "group.agent_removed" \
  --secret "my-secret-key" \
  --header "Content-Type:application/json" \
  --header "X-Custom-Header:value" \
  --enabled

# Webhook para todos os eventos
sloth-runner group webhook create all-events \
  --url "https://example.com/webhook" \
  --event "group.created" \
  --event "group.updated" \
  --event "group.deleted" \
  --event "group.agent_added" \
  --event "group.agent_removed" \
  --event "bulk.operation_start" \
  --event "bulk.operation_end"
```

**Eventos Disponíveis:**
- `group.created` - Novo grupo criado
- `group.updated` - Grupo modificado
- `group.deleted` - Grupo deletado
- `group.agent_added` - Agente adicionado ao grupo
- `group.agent_removed` - Agente removido do grupo
- `bulk.operation_start` - Operação em massa iniciada
- `bulk.operation_end` - Operação em massa completada

#### Habilitar/Desabilitar Webhook

```bash
# Habilitar webhook
sloth-runner group webhook enable slack-notify

# Desabilitar webhook
sloth-runner group webhook disable slack-notify
```

#### Visualizar Logs de Webhooks

```bash
# Ver logs recentes de todos os webhooks
sloth-runner group webhook logs

# Ver logs de webhook específico
sloth-runner group webhook logs --webhook slack-notify

# Ver últimos 50 logs
sloth-runner group webhook logs --limit 50
```

**Saída exemplo:**
```
TIMESTAMP            WEBHOOK        EVENT              STATUS     ERROR
---------            -------        -----              ------     -----
2025-10-08 14:30:15  slack-notify   group.created      ✅ 200     -
2025-10-08 14:25:10  slack-notify   group.agent_added  ✅ 200     -
2025-10-08 14:20:05  discord-1      bulk.operation_end ✅ 200     -
2025-10-08 14:15:00  slack-notify   group.deleted      ❌ 500     Connection timeout
```

#### Deletar Webhook

```bash
# Deletar com confirmação
sloth-runner group webhook delete slack-notify

# Deletar sem confirmação
sloth-runner group webhook delete slack-notify --force
```

## Interface Web

A interface web oferece uma maneira visual de gerenciar grupos de agentes.

### Acessar a Interface

```bash
# Iniciar o servidor web (porta padrão 8080)
sloth-runner ui start

# Iniciar em porta customizada
sloth-runner ui start --port 9090
```

Acesse: `http://localhost:8080/agent-groups`

### Funcionalidades da Interface

A interface web possui 6 abas principais:

1. **Groups** - Gerenciamento básico de grupos
   - Criar, editar, deletar grupos
   - Visualizar detalhes e métricas
   - Adicionar/remover agentes

2. **Templates** - Gerenciamento de templates
   - Criar templates com regras
   - Aplicar templates a grupos
   - Visualizar templates existentes

3. **Hierarchy** - Estrutura hierárquica
   - Visualizar árvore de grupos
   - Criar relacionamentos pai-filho
   - Navegar pela hierarquia

4. **Auto-Discovery** - Configuração de auto-discovery
   - Criar configurações de descoberta
   - Gerenciar schedules
   - Executar discovery manualmente

5. **Webhooks** - Gerenciamento de webhooks
   - Configurar webhooks
   - Visualizar logs de execução
   - Testar webhooks

6. **Bulk Operations** - Operações em massa
   - Executar comandos em grupos
   - Visualizar resultados em tempo real
   - Histórico de operações

## API REST

Todas as funcionalidades estão disponíveis via API REST.

### Configuração

```bash
# Configurar URL da API (padrão: http://localhost:8080)
export SLOTH_RUNNER_API_URL="http://localhost:8080"
```

### Endpoints Principais

#### Grupos

```bash
# Listar grupos
GET /api/v1/agent-groups

# Criar grupo
POST /api/v1/agent-groups
{
  "group_name": "production-web",
  "description": "Production web servers",
  "tags": {"env": "production"},
  "agent_names": []
}

# Obter grupo
GET /api/v1/agent-groups/{group_id}

# Deletar grupo
DELETE /api/v1/agent-groups/{group_id}

# Adicionar agentes
POST /api/v1/agent-groups/{group_id}/agents
{
  "agent_names": ["server-01", "server-02"]
}

# Remover agentes
DELETE /api/v1/agent-groups/{group_id}/agents
{
  "agent_names": ["server-01"]
}
```

#### Operações em Massa

```bash
# Executar operação em massa
POST /api/v1/agent-groups/bulk-operation
{
  "group_id": "production-web",
  "operation": "restart",
  "params": {},
  "timeout": 300
}
```

#### Templates

```bash
# Listar templates
GET /api/v1/agent-groups/templates

# Criar template
POST /api/v1/agent-groups/templates
{
  "name": "web-servers",
  "description": "Web server template",
  "rules": [
    {
      "type": "tag_match",
      "operator": "equals",
      "value": "web"
    }
  ],
  "tags": {"env": "production"}
}

# Aplicar template
POST /api/v1/agent-groups/templates/{template_id}/apply
{
  "group_id": "production-web"
}

# Deletar template
DELETE /api/v1/agent-groups/templates/{template_id}
```

#### Auto-Discovery

```bash
# Listar configurações
GET /api/v1/agent-groups/auto-discovery

# Criar configuração
POST /api/v1/agent-groups/auto-discovery
{
  "name": "web-discovery",
  "group_id": "production-web",
  "schedule": "*/10 * * * *",
  "rules": [...],
  "enabled": true
}

# Executar discovery
POST /api/v1/agent-groups/auto-discovery/{config_id}/run

# Deletar configuração
DELETE /api/v1/agent-groups/auto-discovery/{config_id}
```

#### Webhooks

```bash
# Listar webhooks
GET /api/v1/agent-groups/webhooks

# Criar webhook
POST /api/v1/agent-groups/webhooks
{
  "name": "slack-notify",
  "url": "https://hooks.slack.com/...",
  "events": ["group.created", "group.deleted"],
  "secret": "optional-secret",
  "enabled": true
}

# Ver logs
GET /api/v1/agent-groups/webhooks/logs?limit=20&webhook_id=slack-1

# Deletar webhook
DELETE /api/v1/agent-groups/webhooks/{webhook_id}
```

## Casos de Uso

### Caso 1: Ambiente de Produção Web

Gerenciar servidores web de produção com auto-discovery e webhooks.

```bash
# 1. Criar grupo
sloth-runner group create production-web \
  --description "Production web servers" \
  --tag environment=production \
  --tag role=webserver

# 2. Configurar auto-discovery (a cada 10 minutos)
sloth-runner group auto-discovery create web-disc \
  --group production-web \
  --schedule "*/10 * * * *" \
  --rule "tag_match:equals:webserver" \
  --rule "tag_match:equals:production" \
  --enabled

# 3. Configurar webhook para Slack
sloth-runner group webhook create slack-prod-web \
  --url "https://hooks.slack.com/services/YOUR/WEBHOOK" \
  --event "group.agent_added" \
  --event "bulk.operation_end" \
  --enabled

# 4. Executar atualização em todos os servidores
sloth-runner group bulk production-web execute \
  --command "apt-get update && apt-get upgrade -y" \
  --timeout 600
```

### Caso 2: Reiniciar Serviços em Múltiplos Servidores

```bash
# 1. Criar grupo temporário
sloth-runner group create nginx-restart \
  --description "Servers needing nginx restart"

# 2. Adicionar servidores
sloth-runner group add-agent nginx-restart \
  server-01 server-02 server-03 server-04 server-05

# 3. Executar restart
sloth-runner group bulk nginx-restart execute \
  --command "systemctl restart nginx && systemctl status nginx"

# 4. Deletar grupo após uso
sloth-runner group delete nginx-restart --force
```

### Caso 3: Monitoramento com Templates

```bash
# 1. Criar template para agentes de monitoramento
sloth-runner group template create monitoring \
  --description "Monitoring agents template" \
  --rule "tag_match:equals:monitoring" \
  --rule "status:equals:active"

# 2. Criar grupo usando template
sloth-runner group create monitoring-agents \
  --description "Active monitoring agents"

# 3. Aplicar template
sloth-runner group template apply monitoring monitoring-agents

# 4. Configurar auto-discovery
sloth-runner group auto-discovery create monitoring-disc \
  --group monitoring-agents \
  --schedule "*/5 * * * *" \
  --rule "tag_match:equals:monitoring" \
  --enabled
```

### Caso 4: Deploy em Múltiplos Ambientes

```bash
# Criar grupos por ambiente
for env in dev staging production; do
  sloth-runner group create ${env}-web \
    --description "${env} web servers" \
    --tag environment=${env} \
    --tag role=webserver

  # Auto-discovery por ambiente
  sloth-runner group auto-discovery create ${env}-disc \
    --group ${env}-web \
    --schedule "*/15 * * * *" \
    --rule "tag_match:equals:${env}" \
    --rule "tag_match:equals:webserver" \
    --enabled
done

# Deploy em staging primeiro
sloth-runner group bulk staging-web execute \
  --command "git pull && npm install && pm2 restart app"

# Depois deploy em produção
sloth-runner group bulk production-web execute \
  --command "git pull && npm install && pm2 restart app"
```

## Variáveis de Ambiente

```bash
# URL da API (padrão: http://localhost:8080)
export SLOTH_RUNNER_API_URL="http://api.example.com:8080"

# Endereço do master server (para agentes)
export SLOTH_RUNNER_MASTER_ADDR="192.168.1.29:50053"
```

## Troubleshooting

### API não responde

```bash
# Verificar se o servidor UI está rodando
ps aux | grep "sloth-runner ui"

# Iniciar servidor se não estiver rodando
sloth-runner ui start --port 8080
```

### Webhook não dispara

```bash
# Ver logs de webhooks
sloth-runner group webhook logs --webhook webhook-id --limit 50

# Verificar se webhook está habilitado
sloth-runner group webhook list
```

### Auto-discovery não funciona

```bash
# Executar manualmente para testar
sloth-runner group auto-discovery run config-id

# Verificar se está habilitado
sloth-runner group auto-discovery list

# Habilitar se necessário
sloth-runner group auto-discovery enable config-id
```

### Operação em massa falhou em alguns agentes

```bash
# O comando bulk mostra quais agentes falharam
# Exemplo de saída:
# server-03   ❌ FAILED   1200ms   Connection timeout

# Verificar status do agente
sloth-runner agent get server-03

# Tentar operação individual
sloth-runner agent restart server-03
```

## Exemplos de Scripts

### Script de Backup Automatizado

```bash
#!/bin/bash

# Criar grupo de servidores de banco de dados
sloth-runner group create db-backup \
  --description "Database servers for backup"

# Adicionar servidores
sloth-runner group add-agent db-backup db-01 db-02 db-03

# Executar backup
sloth-runner group bulk db-backup execute \
  --command "mysqldump -u root -p\$DB_PASSWORD --all-databases > /backup/db-\$(date +%Y%m%d).sql" \
  --timeout 1800

# Verificar resultado
if [ $? -eq 0 ]; then
  echo "✅ Backup completed successfully"
else
  echo "❌ Backup failed"
  exit 1
fi
```

### Script de Atualização de Segurança

```bash
#!/bin/bash

# Grupos de servidores por prioridade
GROUPS=("critical" "important" "normal")

for group in "${GROUPS[@]}"; do
  echo "Updating ${group} servers..."

  sloth-runner group bulk ${group}-servers execute \
    --command "apt-get update && apt-get upgrade -y && apt-get autoremove -y" \
    --timeout 900

  # Esperar 5 minutos entre grupos
  if [ "$group" != "normal" ]; then
    echo "Waiting 5 minutes before next group..."
    sleep 300
  fi
done

echo "✅ All security updates completed"
```

## Referências

- [Documentação de Módulos](modules/README.md)
- [Documentação de Agentes](agent-management.md)
- [Documentação de Hooks](hooks.md)
- [API Reference](api-reference.md)
