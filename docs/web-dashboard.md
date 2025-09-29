# 🎨 Web Dashboard & UI

Sloth Runner provides a **comprehensive web-based dashboard** for managing workflows, monitoring agents, and visualizing task execution in real-time.

## 🚀 Quick Start

### Starting the Dashboard

```bash
# Basic UI server
sloth-runner ui --port 8080

# Daemon mode with custom port
sloth-runner ui --port 3000 --daemon

# With debug logging
sloth-runner ui --port 8080 --debug
```

### Accessing the Dashboard
Once started, open your browser and navigate to:
```
http://localhost:8080
```

## 🎯 Core Features

### 📊 Real-Time Monitoring
- **Live task execution** tracking
- **Progress indicators** for running workflows
- **Resource utilization** metrics
- **Performance graphs** and charts

### 🌐 Agent Management
- **Visual agent topology** display
- **Agent health status** monitoring  
- **Real-time heartbeat** tracking
- **Agent capability** visualization

### 📈 Workflow Dashboard
- **Workflow execution** history
- **Task dependency** graphs
- **Success/failure** statistics
- **Duration analytics**

### 📋 Log Management
- **Centralized log** aggregation
- **Real-time log** streaming
- **Filterable** by agent, task, or workflow
- **Searchable** log history

## 🎨 Dashboard Sections

### 1. Overview Page
**Main dashboard** with key metrics and status:

```
┌─────────────────────┬─────────────────────┐
│   Active Workflows  │    Agent Status     │
│        12           │     🟢 8 Active     │
│                     │     🔴 2 Offline    │
├─────────────────────┼─────────────────────┤
│   Tasks Executed    │   Success Rate      │
│       1,247         │       98.5%         │
└─────────────────────┴─────────────────────┘
```

### 2. Workflow Manager
**Visual workflow** builder and manager:
- **Drag-and-drop** task creation
- **Visual dependency** mapping
- **Workflow templates** library
- **Version control** integration

### 3. Agent Console
**Comprehensive agent** management interface:

```
Agent Name      Status    Last Seen    Tasks    CPU   Memory
build-agent-1   🟢 Active  2s ago      3/5      45%   62%
test-agent-2    🟢 Active  1s ago      2/4      32%   48% 
deploy-agent-3  🔴 Offline 2m ago      0/3      --    --
```

### 4. Execution Monitor
**Real-time execution** tracking:
- **Live progress** bars for running tasks
- **Task output** streaming
- **Error highlighting** and alerts
- **Execution timeline** visualization

### 5. Log Viewer
**Advanced log** analysis interface:
- **Multi-level filtering** (ERROR, WARN, INFO, DEBUG)
- **Agent-specific** log views
- **Keyword search** across all logs
- **Export functionality** for debugging

## 🔧 Configuration

### Environment Variables
```bash
# UI Configuration
export SLOTH_UI_PORT=8080
export SLOTH_UI_DEBUG=true
export SLOTH_UI_THEME=dark

# Database Connection
export SLOTH_DB_PATH=/data/sloth.db

# Security Settings
export SLOTH_UI_AUTH_ENABLED=true
export SLOTH_UI_SESSION_SECRET=your-secret-key
```

### Configuration File
```yaml
# ui-config.yaml
ui:
  port: 8080
  debug: false
  theme: "light"
  
dashboard:
  refresh_interval: 5s
  max_log_lines: 1000
  chart_retention: "24h"

security:
  auth_enabled: false
  session_timeout: "30m"
  csrf_protection: true
```

## 📊 Monitoring Features

### Real-Time Metrics
The dashboard displays **live metrics** including:

#### System Metrics
- **CPU utilization** across all agents
- **Memory usage** patterns
- **Network I/O** statistics
- **Disk usage** monitoring

#### Workflow Metrics  
- **Tasks per minute** execution rate
- **Average task duration**
- **Success/failure ratios**
- **Queue depth** monitoring

#### Agent Metrics
- **Agent availability** percentages
- **Task distribution** across agents
- **Agent response times**
- **Heartbeat latency**

### Alerting System
**Built-in alerting** for critical events:

```javascript
// Alert configuration
{
  "alerts": [
    {
      "name": "Agent Offline",
      "condition": "agent.status == 'offline'",
      "duration": "30s",
      "action": "email"
    },
    {
      "name": "High Task Failure Rate", 
      "condition": "task.failure_rate > 0.1",
      "duration": "5m",
      "action": "slack"
    }
  ]
}
```

## 🎨 Themes & Customization

### Built-in Themes
- **Light theme** - Clean, professional appearance
- **Dark theme** - Reduced eye strain for long sessions  
- **High contrast** - Accessibility focused
- **Custom theme** - Company branding support

### Dashboard Customization
```css
/* Custom CSS styling */
.dashboard-widget {
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.status-active {
  color: #28a745;
}

.status-offline {
  color: #dc3545;
}
```

## 🔒 Security Features

### Authentication
```bash
# Enable authentication
sloth-runner ui --auth-enabled

# With custom auth provider
sloth-runner ui --auth-provider ldap --auth-config auth.yaml
```

### Authorization
**Role-based access control** (RBAC):
- **Admin** - Full system access
- **Operator** - Workflow management
- **Viewer** - Read-only access

### Session Management
- **Secure session** cookies
- **Automatic logout** after inactivity
- **CSRF protection** enabled
- **XSS prevention** measures

## 📱 Mobile Responsiveness

The dashboard is **fully responsive** and works on:
- **Desktop computers** (optimal experience)
- **Tablets** (touch-friendly interface)
- **Mobile phones** (essential functions)

## 🔌 API Integration

### REST API Endpoints
```bash
# Get dashboard data
GET /api/v1/dashboard/overview

# List workflows
GET /api/v1/workflows

# Get agent status
GET /api/v1/agents

# Stream logs
GET /api/v1/logs/stream?agent=build-agent-1
```

### WebSocket API
```javascript
// Real-time updates via WebSocket
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');

ws.onmessage = function(event) {
  const update = JSON.parse(event.data);
  if (update.type === 'agent_status') {
    updateAgentDisplay(update.data);
  }
};
```

## 🎯 Use Cases

### DevOps Monitoring
- **CI/CD pipeline** visualization
- **Deployment status** tracking
- **Infrastructure monitoring**
- **Performance analytics**

### Team Collaboration
- **Shared workflow** visibility
- **Task assignment** tracking
- **Progress reporting**
- **Issue identification**

### Troubleshooting
- **Error diagnosis** tools
- **Log correlation** across agents
- **Performance bottleneck** identification
- **System health** monitoring

## 🚀 Advanced Features

### Workflow Builder
**Visual workflow designer** with:
- **Drag-and-drop** task creation
- **Automatic dependency** detection
- **Template library** access
- **Real-time validation**

### Performance Analytics
**Advanced analytics** dashboard:
- **Historical trends** analysis
- **Capacity planning** insights
- **Optimization recommendations**
- **Custom report** generation

### Integration Hub
**External tool** integrations:
- **Slack notifications**
- **Email alerts**
- **Webhook support**
- **Custom plugins**

## 🐛 Troubleshooting

### Common Issues

#### Dashboard Won't Load
```bash
# Check UI server status
ps aux | grep sloth-runner

# Verify port availability
netstat -tlnp | grep 8080

# Check logs
tail -f ui.log
```

#### Slow Performance
```bash
# Enable debug mode
sloth-runner ui --debug

# Check database size
ls -lh sloth.db

# Monitor memory usage
top -p $(pgrep sloth-runner)
```

### Debug Tools
```bash
# Browser developer tools
# Network tab - Check API response times
# Console tab - Look for JavaScript errors
# Performance tab - Analyze rendering bottlenecks
```

## 📊 Metrics Collection

### Prometheus Integration
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'sloth-runner-ui'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

### Grafana Dashboards
Pre-built **Grafana dashboards** available for:
- **System overview** metrics
- **Workflow performance** tracking
- **Agent health** monitoring
- **Error rate** analysis

---

The Web Dashboard transforms Sloth Runner into a **visual, user-friendly platform** that makes complex workflow management accessible to teams of all technical levels! 🎨✨