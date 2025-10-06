# SLOTH-RUNNER-EVENTS(1) - Event Queue Management

## NAME

**sloth-runner events** - View and manage the event queue

## SYNOPSIS

```
sloth-runner events <command> [options]
```

## DESCRIPTION

The **events** command provides tools for monitoring and managing the event queue in Sloth Runner. Events are system-generated notifications about task execution, agent status changes, workflow state, and custom events dispatched from workflows.

The event system operates with:
- **Buffered queue** - 1000 event capacity for high throughput
- **Worker pool** - 100 concurrent workers processing events
- **Persistent storage** - Events stored in SQLite for history and auditing
- **Hook integration** - Events trigger registered hooks automatically

Events capture complete context about what happened in the system, when it happened, and any associated data. Each event includes:
- Unique event ID
- Event type (task.started, agent.disconnected, etc.)
- Timestamp
- Processing status
- Event-specific data (task details, agent info, etc.)
- Hook execution history

## AVAILABLE COMMANDS

- **list** - List events in the queue with filtering
- **show** - Show detailed event information including hook executions
- **get** - Get event details in JSON format
- **delete** - Delete a specific event
- **cleanup** - Remove old events to manage database size

## EVENTS LIST

Display events in the queue with optional filtering by status and type.

### Synopsis

```
sloth-runner events list [options]
```

### Options

```
-l, --limit <n>        Maximum number of events to display (default: 50)
-t, --type <type>      Filter by event type (task.started, agent.disconnected, etc.)
-s, --status <status>  Filter by processing status (pending, processing, completed, failed)
```

### Event Statuses

- **pending** - Event queued, waiting for processing
- **processing** - Event currently being processed by worker
- **completed** - Event successfully processed, all hooks executed
- **failed** - Event processing failed (hook execution error)

### Examples

List recent events (default 50):

```bash
sloth-runner events list
```

List only task failure events:

```bash
sloth-runner events list --type task.failed
```

List failed event processing:

```bash
sloth-runner events list --status failed
```

List more events:

```bash
sloth-runner events list --limit 200
```

Combine filters:

```bash
sloth-runner events list \
  --type task.completed \
  --status completed \
  --limit 100
```

Sample output:

```
Events:
┌──────────────────────┬────────────────────┬────────────────┬───────────┬──────────────────────┐
│ ID                   │ Type               │ Status         │ Timestamp │ Data Summary         │
├──────────────────────┼────────────────────┼───────────────┼───────────┼──────────────────────┤
│ evt_abc123def456     │ task.started       │ completed      │ 17:32:19  │ Task: deploy_app     │
│ evt_def456ghi789     │ task.completed     │ completed      │ 17:32:45  │ Task: deploy_app     │
│ evt_ghi789jkl012     │ task.failed        │ completed      │ 17:33:12  │ Task: health_check   │
│ evt_jkl012mno345     │ agent.disconnected │ completed      │ 17:35:20  │ Agent: web-server-03 │
│ evt_mno345pqr678     │ task.started       │ completed      │ 17:40:05  │ Task: backup_db      │
└──────────────────────┴────────────────────┴───────────────┴───────────┴──────────────────────┘

Total: 5 events
```

## EVENTS SHOW

Display detailed information about a specific event, including all hook executions triggered by the event.

### Synopsis

```
sloth-runner events show <event-id>
```

### Examples

Show event details:

```bash
sloth-runner events show evt_abc123def456
```

Sample output:

```
Event Details
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

ID:          evt_abc123def456
Type:        task.failed
Status:      completed
Created:     2025-10-06 17:32:22
Processed:   2025-10-06 17:32:22
Stack:       production

Event Data:
{
  "task": {
    "task_name": "deploy_application",
    "agent_name": "prod-server-01",
    "status": "failed",
    "exit_code": 1,
    "error": "Connection timeout to database server",
    "duration": "5m30s",
    "started_at": "2025-10-06T17:26:52Z",
    "finished_at": "2025-10-06T17:32:22Z"
  }
}

Hook Executions (2):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. task_failure_alert
   Status:    Success
   Duration:  15ms
   Timestamp: 2025-10-06 17:32:22

   Output:
   🚨 ALERT: Task deploy_application failed
   Alert sent to monitoring system
   Log written to /var/log/sloth/failures.log

2. task_logger
   Status:    Success
   Duration:  8ms
   Timestamp: 2025-10-06 17:32:22

   Output:
   Logged task failure to database
```

Using event details for debugging:

```bash
# Show event and grep for errors
sloth-runner events show evt_abc123 | grep -i error

# Get event ID from list and show details
EVENT_ID=$(sloth-runner events list --type task.failed --limit 1 | grep evt_ | awk '{print $1}')
sloth-runner events show $EVENT_ID
```

## EVENTS GET

Retrieve event details in JSON format for programmatic access and automation.

### Synopsis

```
sloth-runner events get <event-id>
```

### Examples

Get event as JSON:

```bash
sloth-runner events get evt_abc123def456
```

Sample output:

```json
{
  "id": "evt_abc123def456",
  "type": "task.failed",
  "status": "completed",
  "created_at": "2025-10-06T17:32:22Z",
  "processed_at": "2025-10-06T17:32:22Z",
  "stack": "production",
  "data": {
    "task": {
      "task_name": "deploy_application",
      "agent_name": "prod-server-01",
      "status": "failed",
      "exit_code": 1,
      "error": "Connection timeout to database server",
      "duration": "5m30s"
    }
  },
  "hook_executions": [
    {
      "hook_id": "hook_123",
      "hook_name": "task_failure_alert",
      "status": "success",
      "duration_ms": 15,
      "output": "🚨 ALERT: Task deploy_application failed\nAlert sent to monitoring system",
      "executed_at": "2025-10-06T17:32:22Z"
    }
  ]
}
```

Use in automation scripts:

```bash
# Extract task name from failed event
TASK_NAME=$(sloth-runner events get evt_abc123 | jq -r '.data.task.task_name')
echo "Failed task: $TASK_NAME"

# Check if event processing completed
STATUS=$(sloth-runner events get evt_abc123 | jq -r '.status')
if [ "$STATUS" = "completed" ]; then
    echo "Event processed successfully"
fi

# Get error message from failed task
ERROR=$(sloth-runner events get evt_abc123 | jq -r '.data.task.error')
echo "Error: $ERROR"
```

Parse and analyze multiple events:

```bash
# Get all failed task events from today
sloth-runner events list --type task.failed --limit 1000 | \
  grep evt_ | awk '{print $1}' | \
  while read EVENT_ID; do
    echo "Processing event: $EVENT_ID"
    sloth-runner events get $EVENT_ID | jq '.data.task'
  done
```

## EVENTS DELETE

Remove a specific event from the queue and database.

### Synopsis

```
sloth-runner events delete <event-id>
```

### Examples

Delete a specific event:

```bash
sloth-runner events delete evt_abc123def456
```

Delete multiple events:

```bash
# Delete all failed events older than 1 hour
sloth-runner events list --status failed --limit 1000 | \
  grep evt_ | awk '{print $1}' | \
  while read EVENT_ID; do
    sloth-runner events delete $EVENT_ID
  done
```

## EVENTS CLEANUP

Remove old events from the database to manage storage. This is useful for preventing database growth over time.

### Synopsis

```
sloth-runner events cleanup [options]
```

### Options

```
-H, --hours <n>    Remove events older than this many hours (default: 24)
```

### Examples

Clean up events older than 24 hours (default):

```bash
sloth-runner events cleanup
```

Clean up events older than 7 days:

```bash
sloth-runner events cleanup --hours 168
```

Clean up events older than 1 hour:

```bash
sloth-runner events cleanup --hours 1
```

Sample output:

```
Cleaning up events older than 24 hours...
Removed 1,247 events
Database size reduced by 3.2 MB
```

Automated cleanup with cron:

```bash
# Add to crontab to run daily at 2 AM
0 2 * * * /usr/local/bin/sloth-runner events cleanup --hours 168
```

## MONITORING WORKFLOWS

### Real-time Event Monitoring

Monitor events as they occur:

```bash
# Watch recent events (refresh every 2 seconds)
watch -n 2 'sloth-runner events list --limit 10'
```

Monitor specific event types:

```bash
# Monitor task failures in real-time
watch -n 1 'sloth-runner events list --type task.failed --limit 5'

# Monitor agent disconnections
watch -n 5 'sloth-runner events list --type agent.disconnected --limit 10'
```

### Event Analysis

Analyze event patterns:

```bash
# Count events by type
sloth-runner events list --limit 10000 | \
  grep -E '(task\.|agent\.|workflow\.)' | \
  awk '{print $3}' | sort | uniq -c | sort -rn

# Find most common errors
sloth-runner events list --type task.failed --limit 1000 | \
  grep evt_ | awk '{print $1}' | \
  while read id; do
    sloth-runner events get $id | jq -r '.data.task.error'
  done | sort | uniq -c | sort -rn
```

### Debugging with Events

Debug task execution issues:

```bash
# Find all events for a specific task
sloth-runner events list --limit 1000 | grep "deploy_app"

# Get timeline of task events
sloth-runner events list --limit 100 | \
  grep "deploy_app" | \
  while read line; do
    EVENT_ID=$(echo $line | awk '{print $1}')
    sloth-runner events show $EVENT_ID
  done
```

Debug hook execution:

```bash
# Find events where hooks failed
sloth-runner events list --status failed --limit 100

# Show hook execution details
sloth-runner events show evt_abc123 | grep -A 20 "Hook Executions"
```

## INTEGRATION EXAMPLES

### Alerting Integration

Send alerts based on events:

```bash
#!/bin/bash
# alert-on-failures.sh - Send alerts for task failures

sloth-runner events list --type task.failed --limit 1 | \
  grep evt_ | awk '{print $1}' | \
  while read EVENT_ID; do
    # Get event details
    EVENT_JSON=$(sloth-runner events get $EVENT_ID)

    TASK=$(echo $EVENT_JSON | jq -r '.data.task.task_name')
    ERROR=$(echo $EVENT_JSON | jq -r '.data.task.error')
    AGENT=$(echo $EVENT_JSON | jq -r '.data.task.agent_name')

    # Send to PagerDuty/Slack/Email
    curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
      -H 'Content-Type: application/json' \
      -d "{
        \"text\": \"Task Failed: $TASK on $AGENT\",
        \"blocks\": [{
          \"type\": \"section\",
          \"text\": {
            \"type\": \"mrkdwn\",
            \"text\": \"*Task:* $TASK\\n*Agent:* $AGENT\\n*Error:* $ERROR\"
          }
        }]
      }"
  done
```

### Metrics Collection

Export event metrics:

```bash
#!/bin/bash
# export-metrics.sh - Export event metrics for Prometheus/Grafana

# Count events by type
echo "# Event counts by type"
sloth-runner events list --limit 10000 | \
  awk '{print $3}' | grep -E '^(task|agent|workflow)' | \
  sort | uniq -c | \
  while read count type; do
    echo "sloth_events_total{type=\"$type\"} $count"
  done

# Count by status
echo "# Event counts by status"
sloth-runner events list --limit 10000 | \
  awk '{print $4}' | grep -E '^(pending|completed|failed)' | \
  sort | uniq -c | \
  while read count status; do
    echo "sloth_events_total{status=\"$status\"} $count"
  done
```

### Audit Log Generation

Generate audit logs from events:

```bash
#!/bin/bash
# generate-audit-log.sh - Create audit trail from events

OUTPUT="/var/log/sloth-runner/audit.log"

echo "Sloth Runner Audit Log - $(date)" > $OUTPUT
echo "========================================" >> $OUTPUT
echo "" >> $OUTPUT

sloth-runner events list --limit 10000 | \
  grep evt_ | awk '{print $1}' | \
  while read EVENT_ID; do
    EVENT=$(sloth-runner events get $EVENT_ID)

    TIMESTAMP=$(echo $EVENT | jq -r '.created_at')
    TYPE=$(echo $EVENT | jq -r '.type')
    DATA=$(echo $EVENT | jq -r '.data')

    echo "[$TIMESTAMP] $TYPE" >> $OUTPUT
    echo "  Data: $DATA" >> $OUTPUT
    echo "" >> $OUTPUT
  done

echo "Audit log generated: $OUTPUT"
```

## EVENT DATA STRUCTURE

Events contain type-specific data:

### Task Events (task.started, task.completed, task.failed)

```json
{
  "task": {
    "task_name": "deploy_application",
    "agent_name": "prod-server-01",
    "status": "failed",
    "exit_code": 1,
    "error": "Connection timeout",
    "duration": "5m30s",
    "started_at": "2025-10-06T17:26:52Z",
    "finished_at": "2025-10-06T17:32:22Z"
  }
}
```

### Agent Events (agent.connected, agent.disconnected, etc.)

```json
{
  "agent": {
    "name": "web-server-03",
    "address": "192.168.1.100:50060",
    "last_heartbeat": "2025-10-06T17:25:00Z",
    "status": "disconnected",
    "version": "v5.0.0"
  }
}
```

### Workflow Events (workflow.started, workflow.completed, etc.)

```json
{
  "workflow": {
    "name": "deployment_pipeline",
    "file": "workflows/deploy.sloth",
    "status": "completed",
    "tasks_total": 10,
    "tasks_completed": 10,
    "tasks_failed": 0,
    "duration": "15m22s"
  }
}
```

### Custom Events

```json
{
  "custom_field_1": "value1",
  "custom_field_2": "value2",
  "message": "User-defined event data"
}
```

## BEST PRACTICES

1. **Regular Cleanup** - Schedule cleanup to prevent database growth
   ```bash
   # Cron job: daily cleanup at 2 AM
   0 2 * * * sloth-runner events cleanup --hours 168
   ```

2. **Monitor Critical Events** - Set up alerts for important event types
   ```bash
   watch -n 10 'sloth-runner events list --type task.failed --limit 5'
   ```

3. **Use Filtering** - Filter events to find relevant information quickly
   ```bash
   sloth-runner events list --type task.failed --status completed
   ```

4. **Export for Analysis** - Export events to external systems for long-term storage
   ```bash
   sloth-runner events list --limit 10000 | tee events-backup.log
   ```

5. **Programmatic Access** - Use JSON output for automation
   ```bash
   sloth-runner events get evt_abc123 | jq '.data.task'
   ```

## FILES

```
.sloth-cache/hooks.db    SQLite database storing events and hook executions
```

## SEE ALSO

- **sloth-runner-hook(1)** - Manage event hooks
- **sloth-runner(1)** - Main sloth-runner command
- **sloth-runner-workflow(1)** - Workflow management

## AUTHOR

Written by the Sloth Runner development team.

## COPYRIGHT

Copyright © 2025. Released under MIT License.
