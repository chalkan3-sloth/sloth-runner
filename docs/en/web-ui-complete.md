# 🎨 Complete Web UI Guide

## Overview

Sloth Runner's Web UI is a modern, responsive, and intuitive interface for managing workflows, agents, hooks, events, and monitoring. Built with Bootstrap 5, Chart.js, and WebSockets for real-time updates.

**Access URL:** `http://localhost:8080` (default port)

---

## 🚀 Starting the Web UI

```bash
# Method 1: UI command
sloth-runner ui --port 8080

# Method 2: With specific bind
sloth-runner ui --port 8080 --bind 0.0.0.0

# Method 3: With environment variable
export SLOTH_RUNNER_UI_PORT=8080
sloth-runner ui
```

---

## 📱 Pages and Features

### 1. 🏠 Main Dashboard (`/`)

**Features:**

- **System overview** - Cards with main statistics
  - Total workflows
  - Total active/inactive agents
  - Recent executions
  - Success rate

- **Interactive charts** (Chart.js)
  - Executions per day (last 7 days)
  - Success vs failure rate
  - Agent resource usage
  - Workflow distribution

- **Real-time activity feed**
  - Workflows started/completed
  - Agents connected/disconnected
  - System events
  - WebSocket updates

- **Quick Actions** (floating button)
  - Execute workflow
  - Create new workflow
  - View agents
  - Settings

**Modern features:**
- 🎨 Dark/light mode (automatic toggle)
- 📊 Responsive charts
- 🔄 Auto-refresh every 30 seconds
- 🎯 Smooth animations (fade-in, hover effects)
- 📱 Mobile-first design

---

### 2. 🤖 Agent Management (`/agents`)

**Features:**

#### Agent Cards

Each agent is displayed in a modern card with:

- **Visual status**
  - 🟢 Online (green with pulse animation)
  - 🔴 Offline (gray)
  - 🟡 Connecting (yellow)

- **Real-time metrics**
  - CPU Usage (%) - animated progress chart
  - Memory Usage (%) - animated progress chart
  - Disk Usage (%) - animated progress chart
  - Load Average - converted to % based on CPU cores

- **Agent information**
  - Name and address
  - Agent version
  - Formatted uptime (d/h/m/s)
  - Registration date
  - Last heartbeat

- **Sparklines** (mini trend graphs)
  - CPU usage last 24h
  - Memory usage last 24h
  - Rendered with Canvas API

- **Action buttons**
  - 📊 Dashboard - go to agent dashboard
  - ℹ️ Details - modal with full details
  - 📄 Logs - view agent logs
  - 🔄 Restart - restart agent (only if online)

#### General Statistics

- Total agents
- Active agents
- Inactive agents
- Uptime rate (%)

#### Advanced Features

- **Auto-refresh** - updates list every 10 seconds
- **WebSocket** - real-time notifications when agents connect/disconnect
- **Filters** - filter by status (all/active/inactive)
- **Search** - search agents by name
- **Responsive grid** - cards automatically reorganize
- **Skeleton loaders** - professional loading states

**Layout:**

```
┌─────────────────────────────────────────┐
│  📊 Stats Cards                         │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐  │
│  │Total │ │Active│ │Inact.│ │Uptime│  │
│  └──────┘ └──────┘ └──────┘ └──────┘  │
├─────────────────────────────────────────┤
│  🤖 Agent Cards                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │ Agent 1 │ │ Agent 2 │ │ Agent 3 │  │
│  │ 🟢 80%  │ │ 🟢 45%  │ │ 🔴 N/A  │  │
│  │ [graph] │ │ [graph] │ │ [graph] │  │
│  │ [btns]  │ │ [btns]  │ │ [btns]  │  │
│  └─────────┘ └─────────┘ └─────────┘  │
└─────────────────────────────────────────┘
```

---

### 3. 🎛️ Agent Control (`/agent-control`)

**Features:**

Dedicated page for detailed control of each agent.

- **Agent list** with expanded cards
- **Detailed metrics**
  - CPU cores, load average
  - Memory total/used/free
  - Disk total/used/free
  - Network interfaces
  - Detailed uptime

- **Control actions**
  - Start/Stop/Restart agent
  - Update agent version
  - Check modules
  - Run command remotely
  - View logs

- **Gauge charts** (circular charts)
  - CPU usage
  - Memory usage
  - Disk usage
  - With dynamic colors based on thresholds

**Color thresholds:**
- 🟢 Green: 0-40%
- 🔵 Blue: 40-70%
- 🟡 Yellow: 70-90%
- 🔴 Red: 90-100%

---

### 4. 📊 Agent Dashboard (`/agent-dashboard`)

**Features:**

Individual dashboard for each agent with advanced metrics.

- **Time-series charts** (Chart.js)
  - CPU usage over time (line chart)
  - Memory usage over time (area chart)
  - Disk I/O (bar chart)
  - Network traffic (line chart)

- **System information**
  - OS name, version, kernel
  - CPU model, cores, architecture
  - Total memory, swap
  - Mounted filesystems

- **Process list**
  - Top processes by CPU
  - Top processes by Memory
  - Real-time updates

- **Real-time logs**
  - Agent log stream
  - Filters by level (debug/info/warn/error)
  - Auto-scroll
  - Log download

- **Time range selector**
  - Last 1 hour
  - Last 6 hours
  - Last 24 hours
  - Last 7 days
  - Custom range

---

### 5. 📝 Workflows (`/workflows`)

**Features:**

#### Workflow List

- **Workflow cards** with:
  - Name and description
  - Tags/labels
  - Last execution
  - Success rate
  - Buttons: Run, Edit, Delete

- **Filters**
  - By tags
  - By status (active/inactive)
  - By execution frequency

- **Search** - search by name/description

#### Workflow Editor

**Editor features:**

- **Professional code editor**
  - Syntax highlighting for YAML/Sloth DSL
  - Line numbers
  - Auto-indent
  - Bracket matching
  - Keywords: `tasks`, `name`, `exec`, `delegate_to`, etc.
  - Custom colors (Sloth theme)

- **Keyboard shortcuts**
  - `Tab` - indentation (2 spaces)
  - `Shift+Tab` - de-indentation
  - `Ctrl+S` / `Cmd+S` - save
  - `Shift+Alt+F` - format

- **Templates**
  - Basic workflow
  - Multi-task workflow
  - Distributed workflow (with delegate_to)
  - Docker deployment
  - Full example workflow

- **Real-time validation**
  - Syntax checking
  - Linting
  - Inline warnings and errors

- **Preview**
  - View workflow structure
  - Task dependencies
  - Variables used

**Syntax highlighting example:**

```yaml
# Keywords in purple
tasks:
  - name: Deploy App          # Strings in green
    exec:
      script: |                # Pipe in orange
        pkg.install("nginx")   # Functions in blue
        # Comments in gray
    delegate_to: web-01        # Keys in light blue
```

---

### 6. ⚡ Executions (`/executions`)

**Features:**

Complete history of workflow executions.

- **Execution list** with:
  - Workflow name
  - Status (success/failed/running)
  - Started/completed time
  - Duration
  - Triggered by (user/schedule/hook)
  - Agent name (if delegated)

- **Advanced filters**
  - By status
  - By workflow
  - By agent
  - By date/time
  - By user

- **Execution details**
  - Complete output
  - Structured logs
  - Task-by-task breakdown
  - Variables used
  - Performance metrics

- **Actions**
  - Re-run workflow
  - Download logs
  - Share execution (link)
  - Compare with previous

- **Status indicators**
  - 🟢 Success (green)
  - 🔴 Failed (red)
  - 🟡 Running (yellow with spinner)
  - ⏸️ Paused (gray)

---

### 7. 🎣 Hooks (`/hooks`)

**Features:**

Complete hook (event handler) management.

- **Hook list**
  - Hook name
  - Event type
  - Script path
  - Priority
  - Enabled/disabled status
  - Last triggered

- **Event types**
  - `workflow.started`
  - `workflow.completed`
  - `workflow.failed`
  - `task.started`
  - `task.completed`
  - `task.failed`
  - `agent.connected`
  - `agent.disconnected`

- **Create/Edit hook**
  - Form with fields:
    - Name
    - Event type (dropdown)
    - Script (code editor)
    - Priority (0-100)
    - Enabled (toggle)
  - Syntax highlighting for Lua/Bash
  - Validation

- **Test hook**
  - Dry-run with test payload
  - View output without executing actions
  - Debug mode

- **Hook history**
  - When triggered
  - Payload received
  - Script output
  - Success/failure

---

### 8. 📡 Events (`/events`)

**Features:**

Real-time system event monitoring.

- **Event feed**
  - Timestamp
  - Event type
  - Source (workflow/agent/system)
  - Details/payload
  - Status

- **WebSocket stream**
  - Events appear in real-time
  - Sound notifications (optional)
  - Desktop notifications (optional)

- **Filters**
  - By event type
  - By source
  - By status
  - By time range

- **Export events**
  - JSON
  - CSV
  - Log format

- **Statistics**
  - Events per hour
  - Events by type
  - Top sources

---

### 9. 📅 Scheduler (`/scheduler`)

**Features:**

Workflow scheduling (cron-like).

- **Scheduled jobs**
  - Job name
  - Associated workflow
  - Cron expression
  - Next run time
  - Last run status
  - Enabled/disabled

- **Create job**
  - Form with:
    - Name
    - Workflow (dropdown)
    - Schedule (cron or visual builder)
    - Variables
    - Notifications (on success/failure)

- **Visual cron builder**
  - Minute/hour/day/month selector
  - Preview: "Runs every day at 3:00 AM"
  - Common templates:
    - Every hour
    - Every day at midnight
    - Every Monday at 9 AM
    - Custom

- **Execution history**
  - Per scheduled job
  - Success rate
  - Average duration

---

### 10. 📄 Logs (`/logs`)

**Features:**

Centralized log viewing.

- **Advanced filters**
  - By level (debug/info/warn/error)
  - By source (agent/workflow/system)
  - By time range
  - By text (search)

- **Live tail**
  - Real-time stream
  - Auto-scroll
  - Pause/resume
  - Highlight patterns

- **Structured logs**
  - JSON format
  - Expandable fields
  - Syntax highlighting

- **Export**
  - Download as .log
  - Copy to clipboard
  - Share link

- **Log levels with colors**
  - 🔵 DEBUG (blue)
  - 🟢 INFO (green)
  - 🟡 WARN (yellow)
  - 🔴 ERROR (red)

---

### 11. 🖥️ Interactive Terminal (`/terminal`)

**Features:**

Interactive web terminal for remote agents.

- **xterm.js** - complete terminal
- **SSH-like experience**
- **Multiple sessions** (tabs)
- **Command history** (arrow keys ↑↓)
- **Auto-complete** (Tab)
- **Copy/paste** (Ctrl+Shift+C/V)
- **Themes** - Solarized, Monokai, Dracula, etc.

**Special commands:**
- `.clear` - clear terminal
- `.exit` - close session
- `.upload <file>` - upload file
- `.download <file>` - download file

---

### 12. 📦 Saved Sloths (`/sloths`)

**Features:**

Repository of saved workflows.

- **Sloth list**
  - Name
  - Description
  - Tags
  - Created/updated date
  - Active/inactive status
  - Use count

- **Actions**
  - Run sloth
  - Edit content
  - Clone sloth
  - Export (YAML)
  - Delete
  - Activate/Deactivate

- **Tag management**
  - Create tags
  - Color tags
  - Filter by tags

- **Versioning**
  - View version history
  - Compare versions (diff)
  - Restore previous version

---

### 13. ⚙️ Settings (`/settings`)

**Features:**

#### General Settings

- Server info
  - Master address
  - gRPC port
  - Web UI port
  - Database path

- Preferences
  - Theme (light/dark/auto)
  - Language (en/pt/zh)
  - Timezone
  - Date format

#### Notifications

- Email settings
  - SMTP host, port
  - Username/password
  - From address

- Slack integration
  - Webhook URL
  - Default channel
  - Mentions

- Telegram/Discord
  - Bot token
  - Chat ID / Webhook

#### Security

- Authentication
  - Enable/disable auth
  - User management
  - API tokens

- TLS/SSL
  - Enable HTTPS
  - Certificate upload
  - Auto-renewal (Let's Encrypt)

#### Database

- Backup settings
  - Auto-backup enabled
  - Backup schedule
  - Retention policy

- Maintenance
  - Vacuum database
  - View stats
  - Clear old data

---

## 🎨 Modern Visual Features

### Dark Mode / Light Mode

**Auto-detection** based on system preference + manual toggle

**Themes:**

```css
/* Light Mode */
--bg-primary: #ffffff
--text-primary: #212529
--accent: #4F46E5

/* Dark Mode */
--bg-primary: #1a1d23
--text-primary: #e9ecef
--accent: #818CF8
```

**Toggle:** Button in navbar with icons ☀️ (light) / 🌙 (dark)

---

### Animations and Transitions

- **Fade-in** when loading pages
- **Hover effects** on cards and buttons
- **Pulse animation** on status indicators
- **Skeleton loaders** during loading
- **Smooth scrolling**
- **Ripple effect** on buttons (Material Design)
- **Smooth page transitions**

---

### Glassmorphism

Cards with frosted glass effect:

```css
.glass-card {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}
```

---

### Toasts / Notifications

Modern notification system:

- **Types:**
  - ℹ️ Info (blue)
  - ✅ Success (green)
  - ⚠️ Warning (yellow)
  - ❌ Error (red)
  - ⏳ Loading (with spinner)

- **Positions:**
  - Top-right (default)
  - Top-left
  - Bottom-right
  - Bottom-left
  - Center

- **Features:**
  - Auto-dismiss (configurable)
  - Close button
  - Progress bar
  - Action buttons
  - Multiple toast stacking

---

### Confetti Effects

Confetti effects on special events:

- ✅ Workflow completed successfully
- 🤖 New agent connected
- 🎯 Milestone reached
- 🎉 Deploy completed

```javascript
confetti.burst({
    particleCount: 100,
    spread: 70,
    origin: { y: 0.6 }
});
```

---

### Drag & Drop

- **Reorder tasks** in workflows
- **File upload** (drop zone)
- **Reorganize dashboard** widgets

---

## ⌨️ Command Palette (Ctrl+Shift+P)

Quick access to all actions:

```
🔍 Search commands...

> Run Workflow
> View Agents
> Create Workflow
> Open Terminal
> View Logs
> Settings
> Toggle Dark Mode
> Export Data
...
```

**Features:**
- Fuzzy search
- Keyboard navigation (↑↓ Enter)
- Recent commands
- Visible shortcuts
- Context-aware (shows actions based on current page)

---

## 📊 Charts and Visualizations

### Chart.js Components

**Chart types:**

1. **Line Charts** - temporal metrics
2. **Bar Charts** - comparisons
3. **Doughnut Charts** - distributions
4. **Area Charts** - trends
5. **Sparklines** - inline mini charts

**Features:**
- Responsive
- Interactive tooltips
- Legends
- Zoom/pan
- Export as PNG
- Dark/light themes

---

## 🔄 WebSocket Real-Time Updates

WebSocket connection for real-time updates:

**Real-time events:**
- Agent connected/disconnected
- Workflow started/completed
- New logs
- System alerts
- Metrics updates

**URL:** `ws://localhost:8080/ws`

**Auto-reconnect** if connection drops

---

## 📱 Mobile Responsive

Mobile-first design:

- **Breakpoints:**
  - Mobile: < 768px
  - Tablet: 768px - 1024px
  - Desktop: > 1024px

- **Mobile features:**
  - Hamburger menu
  - Touch-friendly buttons
  - Swipe gestures
  - Simplified charts
  - Bottom navigation

---

## 🔐 Authentication (Optional)

**Login page** if auth is enabled:

- Username/password
- Remember me
- Forgot password
- OAuth (GitHub, Google, etc.)

**JWT tokens** for API

**Roles:**
- Admin - full access
- Operator - execute workflows
- Viewer - view only

---

## 🛠️ Developer Tools

### API Explorer

Explore and test REST API:

```
GET  /api/v1/agents
GET  /api/v1/agents/:name
POST /api/v1/workflows/run
GET  /api/v1/executions
...
```

**Features:**
- Try it out (execute in browser)
- Request/response examples
- Authentication headers
- cURL examples

---

### Logs Browser

Browse system logs:

- Server logs
- Agent logs
- Application logs
- Audit logs

---

## 🎓 Quick Guides

### Quick Start Tour

Interactive tour for new users:

1. Welcome → Agents page
2. Create your first workflow
3. Run a workflow
4. View execution results
5. Set up notifications

**Features:**
- Tooltips with tips
- Highlight elements
- Skip/Next buttons
- Don't show again (cookie)

---

## 💡 Usage Tips

### Keyboard Shortcuts

```
Global:
Ctrl+Shift+P  - Command palette
Ctrl+K        - Search
Ctrl+/        - Help
Esc           - Close modals

Editor:
Ctrl+S        - Save
Ctrl+F        - Find
Ctrl+Z        - Undo
Ctrl+Y        - Redo
```

---

### Bookmarklets

Save important pages:

```
Dashboard:          /
My Workflows:       /workflows
Active Executions:  /executions?status=running
Agent Metrics:      /agent-dashboard
```

---

### Browser Extensions

**Available:**
- Chrome Extension - quick access
- Firefox Add-on

---

## 🔧 Customization

### Custom CSS

Add custom CSS in Settings:

```css
/* Custom theme */
:root {
    --primary-color: #FF6B6B;
    --accent-color: #4ECDC4;
}
```

---

### Custom Widgets

Create custom widgets for dashboard:

- Custom charts
- External integrations
- Iframe embeds

---

## 📚 Next Steps

- [📋 CLI Reference](cli-reference.md)
- [🔧 Modules](modules-complete.md)
- [🎯 Examples](../en/advanced-examples.md)
- [🏗️ Architecture](../architecture/sloth-runner-architecture.md)

---

## 🐛 Troubleshooting

### Web UI won't load

```bash
# Check if server is running
lsof -i :8080

# View logs
sloth-runner ui --port 8080 --verbose

# Clear browser cache
Ctrl+Shift+Del
```

### WebSocket won't connect

```bash
# Check firewall
sudo ufw allow 8080

# Test connection
wscat -c ws://localhost:8080/ws
```

### Metrics not appearing

```bash
# Check if agent is sending metrics
sloth-runner agent metrics <agent-name>

# Restart agent
sloth-runner agent restart <agent-name>
```

---

**Last updated:** 2025-10-07

**Built with:** Bootstrap 5, Chart.js, xterm.js, WebSockets, Canvas API
