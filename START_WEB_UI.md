# 🚀 Starting the Sloth Runner Web UI

## Quick Start

### 1. Build the Project (if not already built)
```bash
go build -o sloth-runner ./cmd/sloth-runner
```

### 2. Start the Web UI

#### Basic Usage (Port 8080)
```bash
./sloth-runner ui
```

#### Custom Port
```bash
./sloth-runner ui --port 9090
```

#### With Authentication
```bash
./sloth-runner ui --auth --username admin --password yourpassword
```

#### Debug Mode
```bash
./sloth-runner ui --debug
```

#### Full Options
```bash
./sloth-runner ui \
  --port 8080 \
  --auth \
  --username admin \
  --password secret123 \
  --debug
```

### 3. Access the Interface

Open your web browser and navigate to:
```
http://localhost:8080
```

Or with custom port:
```
http://localhost:9090
```

---

## 🎯 What You'll See

### Dashboard (Homepage)
When you first open the UI, you'll see:

1. **Navigation Bar**
   - Sloth Runner logo
   - Management dropdown (Agents, Workflows, Hooks, Events, Secrets, SSH)
   - Operations dropdown (Executions, Scheduler, Terminal)
   - Monitoring dropdown (Metrics, Logs)
   - Backup link
   - Theme toggle (Light/Dark mode)
   - WebSocket status indicator

2. **Quick Stats Cards** (Top Row)
   - 🖥️ Agents Online
   - 🔄 Active Workflows
   - ⚡ Active Hooks
   - 🔔 Pending Events

3. **System Health Chart**
   - Real-time CPU usage
   - Real-time Memory usage
   - Updates every 5 seconds

4. **Resource Usage Panel**
   - CPU Usage (%)
   - Memory Usage (%)
   - Active Goroutines

5. **Quick Actions**
   - Run Workflow button
   - Open Terminal button
   - View Metrics button

6. **Recent Activity**
   - Timeline of system events

7. **Connected Agents**
   - List of online agents

8. **Recent Events**
   - Table of latest hook events

---

## 🔍 Exploring Features

### Try These Actions First

#### 1. View Agents
```
Click: Management → Agents
```
- See all registered agents
- Check their status
- View agent metrics

#### 2. Create a Workflow
```
Click: Management → Workflows → Create Workflow
```
- Name your workflow
- Write workflow content
- Save and activate

#### 3. Execute a Workflow
```
Click: Operations → Executions
```
- Select a workflow
- Click "Execute"
- Watch live logs stream in real-time!

#### 4. Open Terminal
```
Click: Operations → Terminal → Create Session
```
- Start a web terminal
- Execute commands
- See output in real-time

#### 5. View Metrics
```
Click: Monitoring → Metrics
```
- Beautiful charts showing system health
- CPU, Memory, Goroutines
- Agent metrics table
- Event and Hook statistics

#### 6. Schedule a Workflow
```
Click: Operations → Scheduler → New Schedule
```
- Choose workflow
- Set cron expression
- Let it run automatically!

#### 7. Create a Backup
```
Click: Backup → Create & Download Backup
```
- Downloads complete system backup
- Includes all databases
- Ready for restore

---

## 🎨 Theme Switching

### Dark Mode
1. Look for the moon icon (🌙) in the top navigation bar
2. Click it to toggle between light and dark themes
3. Your preference is saved automatically

---

## 📊 Real-Time Features

The following features update automatically via WebSocket:

✅ **Dashboard Stats** - Updates every 5 seconds
✅ **System Health Chart** - Live CPU and Memory graphs
✅ **Execution Logs** - Real-time log streaming
✅ **Agent Status** - Instant connection status updates
✅ **Event Queue** - Live event processing updates
✅ **Metrics** - Auto-refreshing charts

---

## 🔧 Troubleshooting

### UI Not Loading?

#### Check if server is running:
```bash
# In another terminal
ps aux | grep sloth-runner
```

#### Check if port is already in use:
```bash
# macOS/Linux
lsof -i :8080

# Or try a different port
./sloth-runner ui --port 9090
```

#### Check for errors:
```bash
# Run in debug mode
./sloth-runner ui --debug
```

### WebSocket Not Connecting?

Look for the connection status in the top-right corner:
- 🔴 Red = Disconnected
- 🟢 Green = Connected

If disconnected:
1. Refresh the page (F5)
2. Check browser console (F12) for errors
3. Ensure server is running
4. Check firewall settings

### Pages Not Found (404)?

Make sure you built the project with the latest code:
```bash
# Clean build
rm sloth-runner
go build -o sloth-runner ./cmd/sloth-runner
./sloth-runner ui
```

### Charts Not Showing?

1. Open browser console (F12)
2. Check for JavaScript errors
3. Ensure Chart.js is loading (check Network tab)
4. Try clearing browser cache (Ctrl+Shift+R)

---

## 📱 Mobile/Tablet Access

The UI is fully responsive! Access from:
- 📱 Mobile phones
- 📲 Tablets
- 💻 Laptops
- 🖥️ Desktop computers

### Access from other devices on your network:

Find your local IP:
```bash
# macOS/Linux
ifconfig | grep "inet " | grep -v 127.0.0.1

# Example output: inet 192.168.1.100
```

Then on mobile/tablet, go to:
```
http://192.168.1.100:8080
```

---

## 🔐 Security Notes

### Running with Authentication

For production use, always enable authentication:
```bash
./sloth-runner ui --auth --username admin --password $(openssl rand -base64 32)
```

### HTTPS/TLS

For production, use a reverse proxy like nginx or traefik with TLS:
```nginx
# nginx example
server {
    listen 443 ssl;
    server_name sloth.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

---

## 🎯 Performance Tips

### For best performance:

1. **Use WebSocket** - Already enabled, gives you real-time updates
2. **Dark Mode** - Reduces screen strain and power consumption
3. **Close unused tabs** - Each page maintains its own WebSocket connection
4. **Use filters** - In logs and events pages to reduce data transfer
5. **Adjust refresh rates** - Some pages allow customizing refresh intervals

---

## 🆘 Getting Help

### Need Help?

1. Check the [Web UI Features Guide](./WEB_UI_FEATURES.md)
2. Look at browser console for errors (F12)
3. Check server logs
4. Visit GitHub repository for issues
5. Review documentation

### Found a Bug?

Please report it on GitHub with:
- Browser and version
- Steps to reproduce
- Expected vs actual behavior
- Screenshots if possible
- Console errors (F12 → Console)

---

## 🚀 Next Steps

Once you have the UI running:

1. **Explore all pages** - Click through each menu item
2. **Create test workflows** - Try running simple workflows
3. **Set up agents** - Connect remote agents
4. **Create hooks** - Automate event responses
5. **Schedule workflows** - Set up cron jobs
6. **Monitor metrics** - Watch your system health
7. **Create backups** - Regular backup schedule

---

## 📸 Screenshots

### Dashboard
![Dashboard with metrics and charts showing system health]

### Executions
![Workflow execution with live logs streaming]

### Metrics
![Beautiful charts showing CPU, memory, and system metrics]

### Terminal
![Web-based terminal with command execution]

### Dark Mode
![Dark theme showing all pages]

---

## 🎉 Enjoy!

You now have a fully functional, modern web interface for Sloth Runner!

Key highlights:
✅ Real-time updates via WebSocket
✅ Beautiful dark mode
✅ Live workflow execution with streaming logs
✅ Web-based terminal
✅ Comprehensive metrics and monitoring
✅ Backup and restore capabilities
✅ Cron-based workflow scheduling
✅ Event-driven hooks management

**Happy automating! 🚀**
