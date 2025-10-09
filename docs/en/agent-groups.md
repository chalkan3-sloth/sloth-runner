# Agent Groups - Agent Group Management

The Agent Groups system allows you to organize and manage multiple agents efficiently, facilitating operations at scale.

## Table of Contents

- [Overview](#overview)
- [CLI Commands](#cli-commands)
  - [Basic Management](#basic-management)
  - [Bulk Operations](#bulk-operations)
  - [Templates](#templates)
  - [Auto-Discovery](#auto-discovery)
  - [Webhooks](#webhooks)
- [Web Interface](#web-interface)
- [REST API](#rest-api)
- [Use Cases](#use-cases)

## Overview

Agent Groups offers the following features:

- **Logical Grouping**: Organize agents by function, environment, region, etc.
- **Bulk Operations**: Execute commands on multiple agents simultaneously
- **Templates**: Create reusable groups with matching rules
- **Auto-Discovery**: Automatically discover and add agents based on rules
- **Webhooks**: Receive notifications for group events
- **Hierarchy**: Organize groups in hierarchical structures

## CLI Commands

### Basic Management

#### List Groups

```bash
# List all groups (table format)
sloth-runner group list

# List in JSON format
sloth-runner group list --output json
```

**Example output:**
```
NAME              AGENTS  DESCRIPTION                 TAGS
----              ------  -----------                 ----
production-web    5       Production web servers      env=production,role=web
staging-db        2       Staging database servers    env=staging,role=database
monitoring        3       Monitoring agents           role=monitoring
```

#### Create Group

```bash
# Create basic group
sloth-runner group create production-web

# Create with description
sloth-runner group create production-web \
  --description "Production web servers"

# Create with tags
sloth-runner group create production-web \
  --description "Production web servers" \
  --tag environment=production \
  --tag role=webserver \
  --tag region=us-east-1
```

#### Show Group Details

```bash
# Show in text format
sloth-runner group show production-web

# Show in JSON format
sloth-runner group show production-web --output json
```

**Example output:**
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

#### Delete Group

```bash
# Delete with confirmation
sloth-runner group delete production-web

# Force delete without confirmation
sloth-runner group delete production-web --force
```

#### Add Agents to Group

```bash
# Add single agent
sloth-runner group add-agent production-web server-01

# Add multiple agents
sloth-runner group add-agent production-web server-01 server-02 server-03
```

#### Remove Agents from Group

```bash
# Remove single agent
sloth-runner group remove-agent production-web server-01

# Remove multiple agents
sloth-runner group remove-agent production-web server-01 server-02
```

### Bulk Operations

Execute operations on all agents in a group simultaneously.

#### Available Operations

- `restart` - Restart all agents
- `update` - Update all agents
- `shutdown` - Shutdown all agents
- `execute` - Execute custom command

```bash
# Restart all agents in group
sloth-runner group bulk production-web restart

# Update all agents with custom timeout
sloth-runner group bulk production-web update --timeout 600

# Execute custom command
sloth-runner group bulk production-web execute \
  --command "systemctl restart nginx"

# Execute command with timeout
sloth-runner group bulk production-web execute \
  --command "apt-get update && apt-get upgrade -y" \
  --timeout 900
```

**Example output:**
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

Templates allow you to create reusable groups with automatic matching rules.

#### List Templates

```bash
# List all templates
sloth-runner group template list

# List in JSON format
sloth-runner group template list --output json
```

#### Create Template

```bash
# Create template with tag match rule
sloth-runner group template create web-servers \
  --description "Web server template" \
  --rule "tag_match:equals:web" \
  --tag "env:production"

# Create template with multiple rules
sloth-runner group template create prod-db \
  --description "Production database template" \
  --rule "tag_match:equals:database" \
  --rule "name_pattern:contains:prod" \
  --rule "status:equals:active"

# Create template with regex
sloth-runner group template create monitoring-agents \
  --rule "name_pattern:regex:^monitor-.*$" \
  --tag "role:monitoring"
```

**Rule Types:**
- `tag_match` - Match based on agent tags
- `name_pattern` - Match based on agent name
- `status` - Match based on agent status

**Operators:**
- `equals` - Exact equality
- `contains` - Contains substring
- `regex` - Regular expression

#### Apply Template

```bash
# Apply template to create/update group
sloth-runner group template apply web-servers production-web
```

**Example output:**
```
✅ Template applied successfully to group 'production-web'
   Matched agents: 5
```

#### Delete Template

```bash
# Delete with confirmation
sloth-runner group template delete web-servers

# Force delete without confirmation
sloth-runner group template delete web-servers --force
```

### Auto-Discovery

Configure rules to automatically discover and add agents to groups.

#### List Configurations

```bash
# List all auto-discovery configurations
sloth-runner group auto-discovery list

# JSON format
sloth-runner group auto-discovery list --output json
```

**Example output:**
```
ID              NAME            GROUP            SCHEDULE        ENABLED  RULES
--              ----            -----            --------        -------  -----
web-disc        web-discovery   production-web   */10 * * * *    Yes      2
db-disc         db-discovery    production-db    0 * * * *       Yes      1
```

#### Create Configuration

```bash
# Create auto-discovery for web servers (every 10 minutes)
sloth-runner group auto-discovery create web-discovery \
  --group production-web \
  --schedule "*/10 * * * *" \
  --rule "tag_match:equals:web" \
  --rule "status:equals:active" \
  --enabled

# Create for database servers (every hour)
sloth-runner group auto-discovery create db-discovery \
  --group production-db \
  --schedule "0 * * * *" \
  --rule "tag_match:equals:database" \
  --rule "name_pattern:contains:db" \
  --tag "auto_discovered:true"
```

**Schedule Format:** Cron expression (minute hour day month weekday)
- `*/5 * * * *` - Every 5 minutes
- `0 * * * *` - Every hour
- `0 0 * * *` - Daily at midnight
- `0 0 * * 0` - Weekly on Sundays

#### Run Manually

```bash
# Run auto-discovery manually
sloth-runner group auto-discovery run web-discovery
```

**Example output:**
```
✅ Auto-discovery run completed
   Matched agents: 3
```

#### Enable/Disable

```bash
# Enable configuration
sloth-runner group auto-discovery enable web-discovery

# Disable configuration
sloth-runner group auto-discovery disable web-discovery
```

#### Delete Configuration

```bash
# Delete with confirmation
sloth-runner group auto-discovery delete web-discovery

# Force delete without confirmation
sloth-runner group auto-discovery delete web-discovery --force
```

### Webhooks

Configure webhooks to receive notifications for group events.

#### List Webhooks

```bash
# List all webhooks
sloth-runner group webhook list

# JSON format
sloth-runner group webhook list --output json
```

**Example output:**
```
ID              NAME             URL                                      EVENTS  ENABLED
--              ----             ---                                      ------  -------
slack-1         slack-notify     https://hooks.slack.com/services/...     3       Yes
discord-1       discord-webhook  https://discord.com/api/webhooks/...     2       Yes
```

#### Create Webhook

```bash
# Slack webhook
sloth-runner group webhook create slack-notify \
  --url "https://hooks.slack.com/services/YOUR/WEBHOOK/URL" \
  --event "group.created" \
  --event "group.deleted" \
  --event "bulk.operation_end" \
  --enabled

# Webhook with secret and custom headers
sloth-runner group webhook create discord-webhook \
  --url "https://discord.com/api/webhooks/YOUR/WEBHOOK" \
  --event "group.agent_added" \
  --event "group.agent_removed" \
  --secret "my-secret-key" \
  --header "Content-Type:application/json" \
  --header "X-Custom-Header:value" \
  --enabled

# Webhook for all events
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

**Available Events:**
- `group.created` - New group created
- `group.updated` - Group modified
- `group.deleted` - Group deleted
- `group.agent_added` - Agent added to group
- `group.agent_removed` - Agent removed from group
- `bulk.operation_start` - Bulk operation started
- `bulk.operation_end` - Bulk operation completed

#### Enable/Disable Webhook

```bash
# Enable webhook
sloth-runner group webhook enable slack-notify

# Disable webhook
sloth-runner group webhook disable slack-notify
```

#### View Webhook Logs

```bash
# View recent logs for all webhooks
sloth-runner group webhook logs

# View logs for specific webhook
sloth-runner group webhook logs --webhook slack-notify

# View last 50 logs
sloth-runner group webhook logs --limit 50
```

**Example output:**
```
TIMESTAMP            WEBHOOK        EVENT              STATUS     ERROR
---------            -------        -----              ------     -----
2025-10-08 14:30:15  slack-notify   group.created      ✅ 200     -
2025-10-08 14:25:10  slack-notify   group.agent_added  ✅ 200     -
2025-10-08 14:20:05  discord-1      bulk.operation_end ✅ 200     -
2025-10-08 14:15:00  slack-notify   group.deleted      ❌ 500     Connection timeout
```

#### Delete Webhook

```bash
# Delete with confirmation
sloth-runner group webhook delete slack-notify

# Force delete without confirmation
sloth-runner group webhook delete slack-notify --force
```

## Web Interface

The web interface provides a visual way to manage agent groups.

### Access the Interface

```bash
# Start web server (default port 8080)
sloth-runner ui start

# Start on custom port
sloth-runner ui start --port 9090
```

Access: `http://localhost:8080/agent-groups`

### Interface Features

The web interface has 6 main tabs:

1. **Groups** - Basic group management
   - Create, edit, delete groups
   - View details and metrics
   - Add/remove agents

2. **Templates** - Template management
   - Create templates with rules
   - Apply templates to groups
   - View existing templates

3. **Hierarchy** - Hierarchical structure
   - View group tree
   - Create parent-child relationships
   - Navigate hierarchy

4. **Auto-Discovery** - Auto-discovery configuration
   - Create discovery configurations
   - Manage schedules
   - Run discovery manually

5. **Webhooks** - Webhook management
   - Configure webhooks
   - View execution logs
   - Test webhooks

6. **Bulk Operations** - Bulk operations
   - Execute commands on groups
   - View real-time results
   - Operation history

## REST API

All features are available via REST API.

### Configuration

```bash
# Configure API URL (default: http://localhost:8080)
export SLOTH_RUNNER_API_URL="http://localhost:8080"
```

### Main Endpoints

#### Groups

```bash
# List groups
GET /api/v1/agent-groups

# Create group
POST /api/v1/agent-groups
{
  "group_name": "production-web",
  "description": "Production web servers",
  "tags": {"env": "production"},
  "agent_names": []
}

# Get group
GET /api/v1/agent-groups/{group_id}

# Delete group
DELETE /api/v1/agent-groups/{group_id}

# Add agents
POST /api/v1/agent-groups/{group_id}/agents
{
  "agent_names": ["server-01", "server-02"]
}

# Remove agents
DELETE /api/v1/agent-groups/{group_id}/agents
{
  "agent_names": ["server-01"]
}
```

#### Bulk Operations

```bash
# Execute bulk operation
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
# List templates
GET /api/v1/agent-groups/templates

# Create template
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

# Apply template
POST /api/v1/agent-groups/templates/{template_id}/apply
{
  "group_id": "production-web"
}

# Delete template
DELETE /api/v1/agent-groups/templates/{template_id}
```

#### Auto-Discovery

```bash
# List configurations
GET /api/v1/agent-groups/auto-discovery

# Create configuration
POST /api/v1/agent-groups/auto-discovery
{
  "name": "web-discovery",
  "group_id": "production-web",
  "schedule": "*/10 * * * *",
  "rules": [...],
  "enabled": true
}

# Run discovery
POST /api/v1/agent-groups/auto-discovery/{config_id}/run

# Delete configuration
DELETE /api/v1/agent-groups/auto-discovery/{config_id}
```

#### Webhooks

```bash
# List webhooks
GET /api/v1/agent-groups/webhooks

# Create webhook
POST /api/v1/agent-groups/webhooks
{
  "name": "slack-notify",
  "url": "https://hooks.slack.com/...",
  "events": ["group.created", "group.deleted"],
  "secret": "optional-secret",
  "enabled": true
}

# View logs
GET /api/v1/agent-groups/webhooks/logs?limit=20&webhook_id=slack-1

# Delete webhook
DELETE /api/v1/agent-groups/webhooks/{webhook_id}
```

## Use Cases

### Case 1: Production Web Environment

Manage production web servers with auto-discovery and webhooks.

```bash
# 1. Create group
sloth-runner group create production-web \
  --description "Production web servers" \
  --tag environment=production \
  --tag role=webserver

# 2. Configure auto-discovery (every 10 minutes)
sloth-runner group auto-discovery create web-disc \
  --group production-web \
  --schedule "*/10 * * * *" \
  --rule "tag_match:equals:webserver" \
  --rule "tag_match:equals:production" \
  --enabled

# 3. Configure Slack webhook
sloth-runner group webhook create slack-prod-web \
  --url "https://hooks.slack.com/services/YOUR/WEBHOOK" \
  --event "group.agent_added" \
  --event "bulk.operation_end" \
  --enabled

# 4. Execute update on all servers
sloth-runner group bulk production-web execute \
  --command "apt-get update && apt-get upgrade -y" \
  --timeout 600
```

### Case 2: Restart Services on Multiple Servers

```bash
# 1. Create temporary group
sloth-runner group create nginx-restart \
  --description "Servers needing nginx restart"

# 2. Add servers
sloth-runner group add-agent nginx-restart \
  server-01 server-02 server-03 server-04 server-05

# 3. Execute restart
sloth-runner group bulk nginx-restart execute \
  --command "systemctl restart nginx && systemctl status nginx"

# 4. Delete group after use
sloth-runner group delete nginx-restart --force
```

### Case 3: Monitoring with Templates

```bash
# 1. Create template for monitoring agents
sloth-runner group template create monitoring \
  --description "Monitoring agents template" \
  --rule "tag_match:equals:monitoring" \
  --rule "status:equals:active"

# 2. Create group using template
sloth-runner group create monitoring-agents \
  --description "Active monitoring agents"

# 3. Apply template
sloth-runner group template apply monitoring monitoring-agents

# 4. Configure auto-discovery
sloth-runner group auto-discovery create monitoring-disc \
  --group monitoring-agents \
  --schedule "*/5 * * * *" \
  --rule "tag_match:equals:monitoring" \
  --enabled
```

### Case 4: Deploy to Multiple Environments

```bash
# Create groups per environment
for env in dev staging production; do
  sloth-runner group create ${env}-web \
    --description "${env} web servers" \
    --tag environment=${env} \
    --tag role=webserver

  # Auto-discovery per environment
  sloth-runner group auto-discovery create ${env}-disc \
    --group ${env}-web \
    --schedule "*/15 * * * *" \
    --rule "tag_match:equals:${env}" \
    --rule "tag_match:equals:webserver" \
    --enabled
done

# Deploy to staging first
sloth-runner group bulk staging-web execute \
  --command "git pull && npm install && pm2 restart app"

# Then deploy to production
sloth-runner group bulk production-web execute \
  --command "git pull && npm install && pm2 restart app"
```

## Environment Variables

```bash
# API URL (default: http://localhost:8080)
export SLOTH_RUNNER_API_URL="http://api.example.com:8080"

# Master server address (for agents)
export SLOTH_RUNNER_MASTER_ADDR="192.168.1.29:50053"
```

## Troubleshooting

### API Not Responding

```bash
# Check if UI server is running
ps aux | grep "sloth-runner ui"

# Start server if not running
sloth-runner ui start --port 8080
```

### Webhook Not Triggering

```bash
# View webhook logs
sloth-runner group webhook logs --webhook webhook-id --limit 50

# Check if webhook is enabled
sloth-runner group webhook list
```

### Auto-discovery Not Working

```bash
# Run manually to test
sloth-runner group auto-discovery run config-id

# Check if enabled
sloth-runner group auto-discovery list

# Enable if needed
sloth-runner group auto-discovery enable config-id
```

### Bulk Operation Failed on Some Agents

```bash
# Bulk command shows which agents failed
# Example output:
# server-03   ❌ FAILED   1200ms   Connection timeout

# Check agent status
sloth-runner agent get server-03

# Try individual operation
sloth-runner agent restart server-03
```

## Example Scripts

### Automated Backup Script

```bash
#!/bin/bash

# Create database servers group
sloth-runner group create db-backup \
  --description "Database servers for backup"

# Add servers
sloth-runner group add-agent db-backup db-01 db-02 db-03

# Execute backup
sloth-runner group bulk db-backup execute \
  --command "mysqldump -u root -p\$DB_PASSWORD --all-databases > /backup/db-\$(date +%Y%m%d).sql" \
  --timeout 1800

# Check result
if [ $? -eq 0 ]; then
  echo "✅ Backup completed successfully"
else
  echo "❌ Backup failed"
  exit 1
fi
```

### Security Update Script

```bash
#!/bin/bash

# Server groups by priority
GROUPS=("critical" "important" "normal")

for group in "${GROUPS[@]}"; do
  echo "Updating ${group} servers..."

  sloth-runner group bulk ${group}-servers execute \
    --command "apt-get update && apt-get upgrade -y && apt-get autoremove -y" \
    --timeout 900

  # Wait 5 minutes between groups
  if [ "$group" != "normal" ]; then
    echo "Waiting 5 minutes before next group..."
    sleep 300
  fi
done

echo "✅ All security updates completed"
```

## References

- [Module Documentation](modules/README.md)
- [Agent Management Documentation](agent-management.md)
- [Hooks Documentation](hooks.md)
- [API Reference](api-reference.md)
