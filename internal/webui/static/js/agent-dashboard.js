// Agent Dashboard JavaScript
let currentAgent = null;
let agentCharts = {};
let ws = null;

// Initialize
document.addEventListener('DOMContentLoaded', async () => {
    await loadAgents();
    setupWebSocket();

    // Auto-refresh every 10 seconds
    setInterval(() => {
        if (currentAgent) {
            refreshAgentData(currentAgent);
        }
    }, 10000);
});

async function loadAgents() {
    try {
        const data = await API.get('/agents');
        const agentPills = document.getElementById('agentPills');

        if (!data.agents || data.agents.length === 0) {
            agentPills.innerHTML = `
                <li class="nav-item">
                    <span class="nav-link disabled">No agents available</span>
                </li>
            `;
            return;
        }

        agentPills.innerHTML = data.agents.map((agent, index) => {
            const statusClass = agent.status === 'connected' ? 'success' : 'secondary';
            const active = index === 0 ? 'active' : '';

            return `
                <li class="nav-item">
                    <button class="nav-link ${active}"
                            data-agent="${agent.name}"
                            onclick="selectAgent('${agent.name}')">
                        <i class="bi bi-circle-fill text-${statusClass}"></i>
                        ${agent.name}
                    </button>
                </li>
            `;
        }).join('');

        // Auto-select first agent
        if (data.agents.length > 0) {
            await selectAgent(data.agents[0].name);
        }
    } catch (error) {
        notify.error('Failed to load agents');
        console.error(error);
    }
}

async function selectAgent(agentName) {
    currentAgent = agentName;

    // Update active pill
    document.querySelectorAll('#agentPills .nav-link').forEach(link => {
        link.classList.remove('active');
        if (link.dataset.agent === agentName) {
            link.classList.add('active');
        }
    });

    // Check if dashboard already exists
    let dashboard = document.getElementById(`agent-${agentName}`);

    if (!dashboard) {
        // Create new dashboard from template
        dashboard = createAgentDashboard(agentName);
        document.getElementById('agentTabContent').appendChild(dashboard);
    }

    // Hide all dashboards
    document.querySelectorAll('#agentTabContent .tab-pane').forEach(pane => {
        pane.classList.remove('show', 'active');
    });

    // Show selected dashboard
    dashboard.classList.add('show', 'active');

    // Load agent data
    await refreshAgentData(agentName);
}

function createAgentDashboard(agentName) {
    const template = document.getElementById('agent-dashboard-template');
    const clone = template.content.cloneNode(true);
    const dashboard = clone.querySelector('.tab-pane');

    dashboard.id = `agent-${agentName}`;
    dashboard.dataset.agent = agentName;

    // Setup tab switching
    const tabButtons = dashboard.querySelectorAll('.nav-link[data-tab]');
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabName = button.dataset.tab;
            switchAgentTab(agentName, tabName);
        });
    });

    // Setup refresh buttons
    dashboard.querySelector('.refresh-logs')?.addEventListener('click', () => {
        refreshAgentLogs(agentName);
    });

    dashboard.querySelector('.download-logs')?.addEventListener('click', () => {
        downloadAgentLogs(agentName);
    });

    return dashboard;
}

function switchAgentTab(agentName, tabName) {
    const dashboard = document.getElementById(`agent-${agentName}`);

    // Update tab buttons
    dashboard.querySelectorAll('.nav-link[data-tab]').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.tab === tabName) {
            btn.classList.add('active');
        }
    });

    // Update tab content
    dashboard.querySelectorAll('.agent-sub-tab').forEach(tab => {
        tab.classList.remove('active');
        if (tab.dataset.tab === tabName) {
            tab.classList.add('active');
        }
    });

    // Load data for the tab
    switch(tabName) {
        case 'logs':
            refreshAgentLogs(agentName);
            break;
        case 'modules':
            refreshAgentModules(agentName);
            break;
        case 'info':
            refreshAgentInfo(agentName);
            break;
    }
}

async function refreshAgentData(agentName) {
    try {
        const data = await API.get(`/agents/${agentName}`);
        const dashboard = document.getElementById(`agent-${agentName}`);

        if (!dashboard) return;

        // Update overview stats
        const statusBadge = data.status === 'connected'
            ? '<span class="badge bg-success">Connected</span>'
            : '<span class="badge bg-secondary">Offline</span>';

        dashboard.querySelector('.agent-status').innerHTML = statusBadge;
        dashboard.querySelector('.agent-cpu').textContent = (data.cpu_percent || 0).toFixed(1) + '%';
        dashboard.querySelector('.agent-memory').textContent = (data.memory_percent || 0).toFixed(1) + '%';
        dashboard.querySelector('.agent-tasks').textContent = data.tasks_running || 0;
        dashboard.querySelector('.agent-tasks-detail').textContent =
            `${data.tasks_running || 0} / ${data.tasks_completed || 0}`;

        // Update charts
        updateAgentCharts(agentName, data);

    } catch (error) {
        console.error('Failed to refresh agent data:', error);
    }
}

function updateAgentCharts(agentName, data) {
    const dashboard = document.getElementById(`agent-${agentName}`);
    if (!dashboard) return;

    const chartKey = `agent-${agentName}`;

    // Line chart for CPU/Memory over time
    if (!agentCharts[chartKey + '-line']) {
        const ctx = dashboard.querySelector('.agent-chart');
        agentCharts[chartKey + '-line'] = new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: 'CPU %',
                    data: [],
                    borderColor: 'rgb(124, 179, 66)',
                    backgroundColor: 'rgba(124, 179, 66, 0.1)',
                    tension: 0.4
                }, {
                    label: 'Memory %',
                    data: [],
                    borderColor: 'rgb(139, 115, 85)',
                    backgroundColor: 'rgba(139, 115, 85, 0.1)',
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: { beginAtZero: true, max: 100 }
                }
            }
        });
    }

    // Update line chart
    const lineChart = agentCharts[chartKey + '-line'];
    const now = new Date().toLocaleTimeString();
    lineChart.data.labels.push(now);
    lineChart.data.datasets[0].data.push(data.cpu_percent || 0);
    lineChart.data.datasets[1].data.push(data.memory_percent || 0);

    if (lineChart.data.labels.length > 20) {
        lineChart.data.labels.shift();
        lineChart.data.datasets[0].data.shift();
        lineChart.data.datasets[1].data.shift();
    }

    lineChart.update('none');

    // Pie chart for current usage
    if (!agentCharts[chartKey + '-pie']) {
        const ctx = dashboard.querySelector('.agent-pie-chart');
        agentCharts[chartKey + '-pie'] = new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: ['Used', 'Free'],
                datasets: [{
                    data: [0, 100],
                    backgroundColor: [
                        'rgba(124, 179, 66, 0.8)',
                        'rgba(200, 178, 153, 0.3)'
                    ]
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false
            }
        });
    }

    // Update pie chart with average usage
    const avgUsage = ((data.cpu_percent || 0) + (data.memory_percent || 0)) / 2;
    const pieChart = agentCharts[chartKey + '-pie'];
    pieChart.data.datasets[0].data = [avgUsage, 100 - avgUsage];
    pieChart.update();

    // Task history bar chart
    if (!agentCharts[chartKey + '-tasks']) {
        const ctx = dashboard.querySelector('.agent-task-chart');
        agentCharts[chartKey + '-tasks'] = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: ['Running', 'Completed', 'Failed'],
                datasets: [{
                    label: 'Tasks',
                    data: [0, 0, 0],
                    backgroundColor: [
                        'rgba(41, 182, 246, 0.8)',
                        'rgba(124, 179, 66, 0.8)',
                        'rgba(239, 83, 80, 0.8)'
                    ]
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: { beginAtZero: true }
                }
            }
        });
    }

    // Update task chart
    const taskChart = agentCharts[chartKey + '-tasks'];
    taskChart.data.datasets[0].data = [
        data.tasks_running || 0,
        data.tasks_completed || 0,
        data.tasks_failed || 0
    ];
    taskChart.update();
}

async function refreshAgentLogs(agentName) {
    const dashboard = document.getElementById(`agent-${agentName}`);
    if (!dashboard) return;

    const logsContainer = dashboard.querySelector('.agent-logs');
    logsContainer.innerHTML = '<div class="text-center py-5"><div class="spinner-border text-success"></div><p class="mt-3">Loading logs...</p></div>';

    try {
        // Try to get agent-specific logs
        const logFile = `${agentName}.log`;
        const data = await API.get(`/logs/${logFile}?tail=100`);

        if (data.content && data.content.trim()) {
            const lines = data.content.split('\n').filter(l => l.trim());
            logsContainer.innerHTML = lines.map(line => {
                let color = '#7CB342'; // default green

                if (line.includes('ERROR') || line.includes('FATAL')) {
                    color = '#EF5350';
                } else if (line.includes('WARN')) {
                    color = '#FFA726';
                } else if (line.includes('INFO')) {
                    color = '#29B6F6';
                }

                return `<div style="color: ${color}; margin-bottom: 2px;">${escapeHtml(line)}</div>`;
            }).join('');
        } else {
            logsContainer.innerHTML = '<div class="text-muted text-center py-5">No logs available for this agent</div>';
        }

        // Auto-scroll to bottom
        logsContainer.scrollTop = logsContainer.scrollHeight;
    } catch (error) {
        logsContainer.innerHTML = `
            <div class="alert alert-warning m-3">
                <i class="bi bi-exclamation-triangle"></i>
                Unable to load logs for this agent. The log file may not exist yet.
            </div>
        `;
    }
}

async function refreshAgentModules(agentName) {
    const dashboard = document.getElementById(`agent-${agentName}`);
    if (!dashboard) return;

    const modulesContainer = dashboard.querySelector('.agent-modules');
    modulesContainer.innerHTML = '<div class="text-center py-5"><div class="spinner-border text-primary"></div><p class="mt-3">Checking modules...</p></div>';

    try {
        // This would call the agent modules check endpoint
        const modules = [
            { name: 'exec', available: true, version: '1.0' },
            { name: 'fs', available: true, version: '1.0' },
            { name: 'net', available: true, version: '1.0' },
            { name: 'pkg', available: true, version: '1.0' },
            { name: 'docker', available: false, version: '-' },
            { name: 'git', available: true, version: '1.0' },
        ];

        modulesContainer.innerHTML = `
            <div class="row g-3">
                ${modules.map(mod => {
                    const statusClass = mod.available ? 'success' : 'secondary';
                    const statusIcon = mod.available ? 'check-circle' : 'x-circle';

                    return `
                        <div class="col-md-4 col-lg-3">
                            <div class="card">
                                <div class="card-body text-center">
                                    <i class="bi bi-${statusIcon} fs-2 text-${statusClass}"></i>
                                    <h6 class="mt-2 mb-1">${mod.name}</h6>
                                    <small class="text-muted">v${mod.version}</small>
                                </div>
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
        `;
    } catch (error) {
        modulesContainer.innerHTML = `
            <div class="alert alert-danger">
                Failed to check modules: ${error.message}
            </div>
        `;
    }
}

async function refreshAgentInfo(agentName) {
    const dashboard = document.getElementById(`agent-${agentName}`);
    if (!dashboard) return;

    try {
        const data = await API.get(`/agents/${agentName}`);

        const infoContainer = dashboard.querySelector('.agent-info');
        infoContainer.innerHTML = `
            <table class="table table-borderless">
                <tbody>
                    <tr>
                        <th width="40%">Name:</th>
                        <td>${data.name || '-'}</td>
                    </tr>
                    <tr>
                        <th>Status:</th>
                        <td><span class="badge bg-${data.status === 'connected' ? 'success' : 'secondary'}">${data.status || '-'}</span></td>
                    </tr>
                    <tr>
                        <th>Address:</th>
                        <td><code>${data.address || '-'}</code></td>
                    </tr>
                    <tr>
                        <th>Version:</th>
                        <td>${data.version || '-'}</td>
                    </tr>
                    <tr>
                        <th>Registered:</th>
                        <td>${data.registered_at ? formatDate(data.registered_at) : '-'}</td>
                    </tr>
                    <tr>
                        <th>Last Heartbeat:</th>
                        <td>${data.last_heartbeat ? timeAgo(data.last_heartbeat) : '-'}</td>
                    </tr>
                    <tr>
                        <th>Platform:</th>
                        <td>${data.platform || '-'}</td>
                    </tr>
                    <tr>
                        <th>Architecture:</th>
                        <td>${data.arch || '-'}</td>
                    </tr>
                </tbody>
            </table>
        `;

        const historyContainer = dashboard.querySelector('.agent-history');
        historyContainer.innerHTML = `
            <div class="timeline">
                <div class="timeline-item">
                    <i class="bi bi-circle-fill text-success"></i>
                    <div>
                        <strong>Connected</strong>
                        <small class="text-muted d-block">${data.last_heartbeat ? timeAgo(data.last_heartbeat) : 'Recently'}</small>
                    </div>
                </div>
                <div class="timeline-item">
                    <i class="bi bi-circle-fill text-info"></i>
                    <div>
                        <strong>Registered</strong>
                        <small class="text-muted d-block">${data.registered_at ? formatDate(data.registered_at) : 'Unknown'}</small>
                    </div>
                </div>
            </div>

            <style>
                .timeline { position: relative; padding: 20px 0 20px 30px; }
                .timeline-item { position: relative; padding-bottom: 20px; }
                .timeline-item i { position: absolute; left: -30px; top: 5px; font-size: 10px; }
                .timeline-item:not(:last-child)::before {
                    content: '';
                    position: absolute;
                    left: -26px;
                    top: 15px;
                    width: 2px;
                    height: 100%;
                    background: var(--border-color);
                }
            </style>
        `;
    } catch (error) {
        console.error('Failed to load agent info:', error);
    }
}

function downloadAgentLogs(agentName) {
    const dashboard = document.getElementById(`agent-${agentName}`);
    if (!dashboard) return;

    const logsContainer = dashboard.querySelector('.agent-logs');
    const logs = Array.from(logsContainer.children).map(div => div.textContent).join('\n');

    downloadFile(logs, `${agentName}-logs.txt`, 'text/plain');
    notify.success('Logs downloaded!', 2000);
}

function setupWebSocket() {
    ws = new WebSocketManager();
    ws.onMessage((data) => {
        if (data.type === 'agent_update' && currentAgent === data.agent_name) {
            refreshAgentData(currentAgent);
        }
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
