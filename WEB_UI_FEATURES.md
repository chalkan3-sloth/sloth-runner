# 🎨 Sloth Runner Web UI - Complete Feature Guide

## 📋 Table of Contents
- [Overview](#overview)
- [Getting Started](#getting-started)
- [Dashboard](#dashboard)
- [Management Features](#management-features)
- [Operations](#operations)
- [Monitoring](#monitoring)
- [System Features](#system-features)

---

## Overview

The Sloth Runner Web UI provides a comprehensive, modern interface for managing your infrastructure automation platform. Built with Bootstrap 5, Chart.js, and WebSocket for real-time updates.

### Key Features
✅ **Real-time Updates** - WebSocket-based live data streaming
✅ **Dark Mode** - Beautiful dark theme with smooth transitions
✅ **Responsive Design** - Works on desktop, tablet, and mobile
✅ **Interactive Charts** - Real-time metrics visualization
✅ **Terminal Access** - Web-based terminal for remote execution
✅ **Backup & Restore** - Complete system backup capabilities

---

## Getting Started

### Starting the Web UI

```bash
# Start the UI server (default port: 8080)
./sloth-runner ui

# Custom port
./sloth-runner ui --port 9090

# With authentication
./sloth-runner ui --auth --username admin --password secret
```

### Accessing the Interface

Open your browser and navigate to:
```
http://localhost:8080
```

---

## Dashboard

The main dashboard provides a comprehensive overview of your system:

### 📊 Quick Stats Cards
- **Agents Online** - Number of connected agents vs total agents
- **Active Workflows** - Currently active workflows
- **Active Hooks** - Number of enabled hooks
- **Pending Events** - Events waiting to be processed

### 📈 System Health Chart
Real-time line chart showing:
- CPU usage percentage
- Memory usage percentage
- Updates every 5 seconds

### 📦 Resource Usage
Live progress bars for:
- CPU Usage
- Memory Usage
- Active Goroutines

### ⚡ Quick Actions
One-click access to:
- **Run Workflow** - Execute workflows instantly
- **Open Terminal** - Web-based terminal
- **View Metrics** - Detailed system metrics

### 📜 Recent Activity
- Timeline of recent system events
- Agent connection/disconnection events
- Workflow executions

### 🖥️ Connected Agents
- List of all connected agents
- Agent status and health
- Quick access to agent details

### 🔔 Recent Events
- Table of recent hook events
- Event status and processing time
- Quick link to full events page

---

## Management Features

### 🖥️ Agents
**Path:** `/agents`

Manage distributed agents:
- View all registered agents
- Check agent status (connected/disconnected)
- View agent metrics (CPU, memory, tasks)
- Update agent configuration
- Remove agents
- Check available modules on each agent

### 🔄 Workflows
**Path:** `/workflows`

Manage automation workflows:
- List all workflows (sloths)
- Create new workflows
- Edit existing workflows
- Activate/deactivate workflows
- Run workflows immediately
- View workflow content and configuration

### 🪝 Hooks
**Path:** `/hooks`

Event-driven automation:
- Create hooks for various event types
- Enable/disable hooks
- Edit hook scripts (Lua)
- View hook execution history
- Monitor hook performance
- Configure event filters

### 🔔 Events
**Path:** `/events`

Event management:
- View all system events
- Filter by status (pending/processed/failed)
- Filter by event type
- Retry failed events
- View event details and metadata
- Monitor event processing queue

### 🔐 Secrets
**Path:** `/secrets`

Secure credential management:
- Store secrets per stack
- View secret keys (values hidden)
- Add new secrets
- Delete secrets
- Organize by stack/environment

### 🔑 SSH Profiles
**Path:** `/ssh`

SSH connection management:
- Create SSH connection profiles
- Store SSH credentials securely
- Manage multiple SSH targets
- View connection audit logs
- Edit profile configurations

---

## Operations

### ▶️ Executions
**Path:** `/executions`

Execute and monitor workflows in real-time:

#### Features:
- **Execute Workflow** - Start workflow execution
  - Select from active workflows
  - Choose delegation target (local or remote agent)
  - Pass variables as JSON

- **Live Logs** - Real-time log streaming
  - WebSocket-based live updates
  - Auto-scroll option
  - Line counter
  - Download logs as file

- **Execution Control**
  - Cancel running executions
  - View execution status
  - Monitor duration
  - Check exit codes

- **Execution History**
  - List of all executions
  - Status indicators (running/completed/failed)
  - Quick access to execution logs
  - Time tracking

### 📅 Scheduler
**Path:** `/scheduler`

Cron-based workflow scheduling:

#### Features:
- **Create Schedules**
  - Name your schedule
  - Select workflow to run
  - Define cron expression
  - Set variables
  - Enable/disable schedule

- **Cron Expression Help**
  - Format: `minute hour day month weekday`
  - Examples provided:
    - `0 0 * * *` - Daily at midnight
    - `*/15 * * * *` - Every 15 minutes
    - `0 9-17 * * 1-5` - Hourly, 9am-5pm, Mon-Fri

- **Schedule Management**
  - Enable/disable schedules
  - Trigger manual execution
  - Edit schedule configuration
  - View next run time
  - Track run count and history
  - Delete schedules

### 💻 Terminal
**Path:** `/terminal`

Web-based terminal access:

#### Features:
- **Create Sessions**
  - Local or remote agent execution
  - Custom shell commands
  - Multiple concurrent sessions

- **Terminal Interface**
  - WebSocket-based terminal
  - Command history (↑↓ arrows)
  - Ctrl+C support
  - Real-time output

- **Session Management**
  - List active sessions
  - Switch between sessions
  - Close sessions
  - Track session agent

---

## Monitoring

### 📊 Metrics
**Path:** `/metrics`

Comprehensive system monitoring with real-time charts:

#### System Overview Cards:
- **CPU Usage** - Percentage and core count
- **Memory Usage** - Percentage and actual usage
- **Goroutines** - Active goroutine count
- **Connected Agents** - Online agents

#### Charts:
- **CPU & Memory Usage** - Real-time line chart
- **Workflow Activity** - Bar chart of workflow stats
- **Event Statistics** - Doughnut chart (pending/processed/failed)
- **Hook Statistics** - Doughnut chart (active/executed/failed)

#### Agent Metrics Table:
For each agent:
- Status (connected/disconnected)
- CPU percentage
- Memory percentage
- Tasks running
- Tasks completed
- Last heartbeat time

#### Historical Metrics:
- Time range selector (1h/6h/24h/7d)
- Multi-metric comparison
- CPU, Memory, and Goroutine trends
- Exportable data

### 📄 Logs
**Path:** `/logs`

Centralized log management:

#### Features:
- **Log File Browser**
  - List all available log files
  - File size and modification time
  - Color-coded by log type

- **Log Viewer**
  - Syntax highlighting
  - Color-coded by level (ERROR/WARN/INFO/DEBUG)
  - Auto-scroll option
  - Line count

- **Filters**
  - Search/grep functionality
  - Filter by log level
  - Tail last N lines
  - Real-time filtering

- **Log Actions**
  - Download log files
  - Clear view
  - Refresh logs
  - Auto-refresh toggle

---

## System Features

### 💾 Backup & Restore
**Path:** `/backup`

Complete system backup solution:

#### Create Backup:
- One-click backup creation
- Includes all databases:
  - Agents database
  - Workflows (Sloths) database
  - Hooks database
  - Secrets database
  - SSH profiles database
- Downloads as `.tar.gz` archive
- Timestamp-based naming

#### Restore Backup:
- Upload backup file
- File validation
- Size display
- Confirmation dialog
- Automatic page reload after restore

#### Best Practices:
Built-in recommendations for:
- **Backup Frequency** - Daily or weekly backups
- **Storage** - Secure, separate location
- **Testing** - Periodic backup verification
- **Versioning** - Keep multiple versions
- **Documentation** - Backup procedures

- **Restore Precautions**
  - Always backup before restore
  - Verify file integrity
  - Stop running workflows
  - Notify team members
  - Verify after restoration

### 🎨 Theme System

#### Dark Mode:
- Toggle via navbar button
- Smooth transitions
- Persists in localStorage
- Affects all pages
- Optimized color schemes

#### Features:
- CSS custom properties
- Automatic chart theme updates
- High contrast for accessibility
- Professional color palette

### 🔔 Notification System

#### Toast Notifications:
- Success messages (green)
- Error messages (red)
- Warning messages (yellow)
- Info messages (blue)
- Auto-dismiss (configurable duration)
- Manual dismiss option
- Stacking support

#### Desktop Notifications:
- Browser notification API
- Permission request
- System tray notifications
- Works when tab is inactive

---

## Technical Details

### Technologies Used:
- **Frontend:** Bootstrap 5, HTML5, CSS3, JavaScript
- **Charts:** Chart.js 4.4
- **Icons:** Bootstrap Icons
- **WebSocket:** gorilla/websocket
- **Backend:** Go, Gin framework

### API Endpoints:
All endpoints use `/api/v1` prefix:
- `/dashboard` - Dashboard statistics
- `/agents/*` - Agent management
- `/sloths/*` - Workflow management
- `/hooks/*` - Hook management
- `/events/*` - Event management
- `/secrets/*` - Secret management
- `/ssh/*` - SSH profile management
- `/executions/*` - Workflow execution
- `/scheduler/*` - Schedule management
- `/terminal/*` - Terminal sessions
- `/metrics` - System metrics
- `/logs/*` - Log management
- `/backup/*` - Backup operations
- `/ws` - WebSocket connection

### Real-time Updates:
- WebSocket connection status indicator
- Automatic reconnection
- Live data streaming
- Bi-directional communication
- Event broadcasting

---

## Navigation Structure

### Main Menu:
1. **Dashboard** - System overview
2. **Management** (Dropdown)
   - Agents
   - Workflows
   - Hooks
   - Events
   - Secrets
   - SSH Profiles
3. **Operations** (Dropdown)
   - Executions
   - Scheduler
   - Terminal
4. **Monitoring** (Dropdown)
   - Metrics
   - Logs
5. **Backup** - System backup

### User Controls:
- **Theme Toggle** - Switch between light/dark mode
- **Connection Status** - WebSocket connection indicator

---

## Tips & Tricks

### Keyboard Shortcuts:
- **Terminal:** Use ↑↓ arrows for command history
- **Terminal:** Ctrl+C to interrupt processes

### Performance:
- Metrics update every 5 seconds
- Charts maintain last 20 data points
- Auto-refresh is configurable
- WebSocket for efficient updates

### Best Practices:
1. **Regular Backups** - Schedule weekly backups
2. **Monitor Metrics** - Check resource usage regularly
3. **Review Logs** - Investigate errors promptly
4. **Scheduler** - Use cron for recurring tasks
5. **Hooks** - Automate event responses
6. **Terminal** - Quick debugging access

---

## Troubleshooting

### WebSocket Not Connecting:
- Check if server is running
- Verify port is not blocked
- Check browser console for errors
- Try refreshing the page

### Metrics Not Loading:
- Ensure metrics endpoint is accessible
- Check backend logs
- Verify agent connectivity

### Terminal Not Working:
- Check shell is available on target
- Verify SSH configuration for remote agents
- Check terminal permissions

### Dark Mode Issues:
- Clear browser cache
- Check localStorage
- Try toggling theme manually

---

## Future Enhancements

Planned features:
- [ ] Code editor with syntax highlighting
- [ ] Workflow visual designer
- [ ] Advanced alerting system
- [ ] User management and RBAC
- [ ] Audit log viewer
- [ ] Performance profiling
- [ ] Multi-language support
- [ ] Export/import configurations
- [ ] API documentation viewer
- [ ] Health check dashboard

---

## Support

For issues or feature requests:
- Check GitHub repository
- Review documentation
- Submit issues on GitHub
- Community discussions

---

**Version:** 2.0
**Last Updated:** October 2025
**Status:** Production Ready 🚀
