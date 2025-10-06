// Agent Control Center JavaScript
let agents = [];
let selectedAgents = new Set();
let currentFilter = 'all';
let ws = null;
let currentAgentName = null;

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    initWebSocket();
    loadAgents();
    setInterval(loadAgents, 5000); // Refresh every 5 seconds
});

function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

    ws.onopen = () => {
        console.log('WebSocket connected');
        document.getElementById('ws-status').classList.add('online');
        document.getElementById('ws-status').classList.remove('offline');
        document.getElementById('ws-status-text').textContent = 'Connected';
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleWebSocketMessage(data);
    };

    ws.onclose = () => {
        document.getElementById('ws-status').classList.remove('online');
        document.getElementById('ws-status').classList.add('offline');
        document.getElementById('ws-status-text').textContent = 'Disconnected';
        setTimeout(initWebSocket, 5000);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

function handleWebSocketMessage(data) {
    if (data.type === 'agent_status') {
        updateAgentStatus(data.agent_name, data.status);
    } else if (data.type === 'metrics_update') {
        updateAgentMetrics(data.agent_name, data.metrics);
    } else if (data.type === 'log_entry' && currentAgentName === data.agent_name) {
        appendLogEntry(data.log);
    }
}

async function loadAgents() {
    try {
        const response = await fetch('/api/v1/agents');
        const data = await response.json();
        agents = data.agents || [];
        renderAgents();
        updateOverviewStats();
    } catch (error) {
        console.error('Error loading agents:', error);
    }
}

function renderAgents() {
    const grid = document.getElementById('agents-grid');

    let filteredAgents = agents;
    if (currentFilter !== 'all') {
        filteredAgents = agents.filter(agent => {
            if (currentFilter === 'online') return agent.status === 'connected';
            if (currentFilter === 'offline') return agent.status !== 'connected';
            if (currentFilter === 'warning') return agent.cpu_percent > 70 || agent.memory_percent > 70;
            if (currentFilter === 'error') return agent.cpu_percent > 90 || agent.memory_percent > 90;
            return true;
        });
    }

    if (filteredAgents.length === 0) {
        grid.innerHTML = `
            <div class="col-12 text-center py-5">
                <i class="bi bi-inbox" style="font-size: 4rem; color: var(--text-muted);"></i>
                <p class="mt-3 text-muted">No agents match the current filter</p>
            </div>
        `;
        return;
    }

    grid.innerHTML = filteredAgents.map(agent => renderAgentCard(agent)).join('');
}

function renderAgentCard(agent) {
    const isOnline = agent.status === 'connected';
    const cpuPercent = agent.cpu_percent || 0;
    const memoryPercent = agent.memory_percent || 0;
    const isSelected = selectedAgents.has(agent.name);

    let cardClass = 'agent-control-card';
    if (!isOnline) {
        cardClass += ' offline';
    } else if (cpuPercent > 90 || memoryPercent > 90) {
        cardClass += ' error';
    } else if (cpuPercent > 70 || memoryPercent > 70) {
        cardClass += ' warning';
    }

    const cpuGaugeClass = cpuPercent > 80 ? 'high' : cpuPercent > 50 ? 'medium' : 'low';
    const memGaugeClass = memoryPercent > 80 ? 'high' : memoryPercent > 50 ? 'medium' : 'low';

    return `
        <div class="card ${cardClass}">
            <div class="card-body">
                <div class="d-flex justify-content-between align-items-start mb-3">
                    <div class="form-check">
                        <input class="form-check-input" type="checkbox"
                               id="select-${agent.name}"
                               ${isSelected ? 'checked' : ''}
                               onchange="toggleAgentSelection('${agent.name}')">
                        <label class="form-check-label" for="select-${agent.name}">
                            <h5 class="mb-0">
                                <i class="bi bi-server"></i>
                                ${agent.name}
                            </h5>
                        </label>
                    </div>
                    <span class="live-indicator ${isOnline ? 'online' : 'offline'}"></span>
                </div>

                <div class="mb-3">
                    <div class="d-flex justify-content-between mb-1">
                        <small class="text-muted">CPU</small>
                        <small class="fw-bold">${cpuPercent.toFixed(1)}%</small>
                    </div>
                    <div class="resource-gauge">
                        <div class="resource-gauge-fill ${cpuGaugeClass}" style="width: ${cpuPercent}%"></div>
                    </div>
                </div>

                <div class="mb-3">
                    <div class="d-flex justify-content-between mb-1">
                        <small class="text-muted">Memory</small>
                        <small class="fw-bold">${memoryPercent.toFixed(1)}%</small>
                    </div>
                    <div class="resource-gauge">
                        <div class="resource-gauge-fill ${memGaugeClass}" style="width: ${memoryPercent}%"></div>
                    </div>
                </div>

                <div class="mb-3">
                    <div class="d-flex justify-content-between">
                        <small class="text-muted">
                            <i class="bi bi-clock"></i>
                            ${formatUptime(agent.uptime_seconds || 0)}
                        </small>
                        <small class="text-muted">
                            <i class="bi bi-hdd"></i>
                            ${agent.disk_percent?.toFixed(0) || 0}%
                        </small>
                    </div>
                </div>

                ${agent.groups && agent.groups.length > 0 ? `
                    <div class="mb-3">
                        ${agent.groups.map(g => `<span class="group-badge">${g}</span>`).join('')}
                    </div>
                ` : ''}

                <div class="d-flex gap-2">
                    <button class="btn btn-sm btn-primary agent-action-btn flex-fill"
                            onclick="openAgentDetail('${agent.name}')" ${!isOnline ? 'disabled' : ''}>
                        <i class="bi bi-eye"></i> Details
                    </button>
                    <button class="btn btn-sm btn-outline-success agent-action-btn"
                            onclick="quickCommand('${agent.name}')" ${!isOnline ? 'disabled' : ''}>
                        <i class="bi bi-terminal"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-warning agent-action-btn"
                            onclick="restartAgent('${agent.name}')" ${!isOnline ? 'disabled' : ''}>
                        <i class="bi bi-arrow-clockwise"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger agent-action-btn"
                            onclick="shutdownAgent('${agent.name}')" ${!isOnline ? 'disabled' : ''}>
                        <i class="bi bi-power"></i>
                    </button>
                </div>
            </div>
        </div>
    `;
}

function updateOverviewStats() {
    const total = agents.length;
    const online = agents.filter(a => a.status === 'connected').length;
    const avgCpu = agents.reduce((sum, a) => sum + (a.cpu_percent || 0), 0) / (agents.length || 1);
    const avgMem = agents.reduce((sum, a) => sum + (a.memory_percent || 0), 0) / (agents.length || 1);

    document.getElementById('total-agents').textContent = total;
    document.getElementById('online-agents').textContent = online;
    document.getElementById('avg-cpu').textContent = avgCpu.toFixed(1) + '%';
    document.getElementById('avg-memory').textContent = avgMem.toFixed(1) + '%';
}

function filterAgents(filter) {
    currentFilter = filter;
    // Update active chip
    document.querySelectorAll('.filter-chip').forEach(chip => {
        chip.classList.remove('active');
    });
    event.target.classList.add('active');
    renderAgents();
}

function toggleAgentSelection(agentName) {
    if (selectedAgents.has(agentName)) {
        selectedAgents.delete(agentName);
    } else {
        selectedAgents.add(agentName);
    }
    updateBulkActionBar();
}

function updateBulkActionBar() {
    const bar = document.getElementById('bulk-action-bar');
    const count = document.getElementById('selected-count');

    count.textContent = selectedAgents.size;

    if (selectedAgents.size > 0) {
        bar.classList.add('active');
    } else {
        bar.classList.remove('active');
    }
}

function clearSelection() {
    selectedAgents.clear();
    document.querySelectorAll('input[type="checkbox"]').forEach(cb => {
        cb.checked = false;
    });
    updateBulkActionBar();
}

async function openAgentDetail(agentName) {
    currentAgentName = agentName;
    document.getElementById('detail-agent-name').textContent = agentName;

    const modal = new bootstrap.Modal(document.getElementById('agentDetailModal'));
    modal.show();

    // Load agent details
    await loadAgentOverview(agentName);
}

let resourceCharts = {};

async function loadAgentOverview(agentName) {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/resources`);
        const data = await response.json();

        const content = document.getElementById('agent-overview-content');
        content.innerHTML = `
            <div class="row">
                <div class="col-md-6">
                    <h6>System Resources</h6>
                    <table class="table table-sm">
                        <tr>
                            <td>CPU Usage:</td>
                            <td><strong>${data.cpu_percent?.toFixed(1)}%</strong></td>
                        </tr>
                        <tr>
                            <td>Memory:</td>
                            <td><strong>${formatBytes(data.memory_used_bytes)} / ${formatBytes(data.memory_total_bytes)}</strong></td>
                        </tr>
                        <tr>
                            <td>Disk:</td>
                            <td><strong>${formatBytes(data.disk_used_bytes)} / ${formatBytes(data.disk_total_bytes)}</strong></td>
                        </tr>
                        <tr>
                            <td>Processes:</td>
                            <td><strong>${data.process_count}</strong></td>
                        </tr>
                        <tr>
                            <td>Uptime:</td>
                            <td><strong>${formatUptime(data.uptime_seconds)}</strong></td>
                        </tr>
                    </table>
                </div>
                <div class="col-md-6">
                    <h6>Load Average</h6>
                    <table class="table table-sm">
                        <tr>
                            <td>1 min:</td>
                            <td><strong>${data.load_avg_1min?.toFixed(2)}</strong></td>
                        </tr>
                        <tr>
                            <td>5 min:</td>
                            <td><strong>${data.load_avg_5min?.toFixed(2)}</strong></td>
                        </tr>
                        <tr>
                            <td>15 min:</td>
                            <td><strong>${data.load_avg_15min?.toFixed(2)}</strong></td>
                        </tr>
                    </table>
                </div>
            </div>

            <div class="row mt-4">
                <div class="col-md-4">
                    <h6 class="text-center">CPU Usage</h6>
                    <canvas id="cpu-chart" style="max-height: 200px;"></canvas>
                </div>
                <div class="col-md-4">
                    <h6 class="text-center">Memory Usage</h6>
                    <canvas id="memory-chart" style="max-height: 200px;"></canvas>
                </div>
                <div class="col-md-4">
                    <h6 class="text-center">Disk Usage</h6>
                    <canvas id="disk-chart" style="max-height: 200px;"></canvas>
                </div>
            </div>

            <div class="row mt-4">
                <div class="col-12">
                    <h6>Resource History (Last 5 minutes)</h6>
                    <canvas id="history-chart" style="height: 250px;"></canvas>
                </div>
            </div>
        `;

        // Destroy old charts if they exist
        ['cpu-chart', 'memory-chart', 'disk-chart', 'history-chart'].forEach(id => {
            if (resourceCharts[id]) {
                resourceCharts[id].destroy();
            }
        });

        // Create doughnut charts for resources
        const cpuChart = new Chart(document.getElementById('cpu-chart'), {
            type: 'doughnut',
            data: {
                labels: ['Used', 'Available'],
                datasets: [{
                    data: [data.cpu_percent, 100 - data.cpu_percent],
                    backgroundColor: ['#4F46E5', '#E5E7EB'],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: {
                            label: (context) => `${context.label}: ${context.parsed.toFixed(1)}%`
                        }
                    }
                }
            }
        });

        const memoryChart = new Chart(document.getElementById('memory-chart'), {
            type: 'doughnut',
            data: {
                labels: ['Used', 'Available'],
                datasets: [{
                    data: [data.memory_percent, 100 - data.memory_percent],
                    backgroundColor: ['#10B981', '#E5E7EB'],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: {
                            label: (context) => `${context.label}: ${context.parsed.toFixed(1)}%`
                        }
                    }
                }
            }
        });

        const diskChart = new Chart(document.getElementById('disk-chart'), {
            type: 'doughnut',
            data: {
                labels: ['Used', 'Available'],
                datasets: [{
                    data: [data.disk_percent, 100 - data.disk_percent],
                    backgroundColor: ['#F59E0B', '#E5E7EB'],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: {
                            label: (context) => `${context.label}: ${context.parsed.toFixed(1)}%`
                        }
                    }
                }
            }
        });

        // Generate mock history data for demonstration
        const historyData = generateMockHistory(data);

        const historyChart = new Chart(document.getElementById('history-chart'), {
            type: 'line',
            data: {
                labels: historyData.labels,
                datasets: [
                    {
                        label: 'CPU %',
                        data: historyData.cpu,
                        borderColor: '#4F46E5',
                        backgroundColor: 'rgba(79, 70, 229, 0.1)',
                        tension: 0.4,
                        fill: true
                    },
                    {
                        label: 'Memory %',
                        data: historyData.memory,
                        borderColor: '#10B981',
                        backgroundColor: 'rgba(16, 185, 129, 0.1)',
                        tension: 0.4,
                        fill: true
                    },
                    {
                        label: 'Disk %',
                        data: historyData.disk,
                        borderColor: '#F59E0B',
                        backgroundColor: 'rgba(245, 158, 11, 0.1)',
                        tension: 0.4,
                        fill: true
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: {
                    mode: 'index',
                    intersect: false
                },
                plugins: {
                    legend: {
                        position: 'top'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            callback: (value) => value + '%'
                        }
                    }
                }
            }
        });

        resourceCharts = {
            'cpu-chart': cpuChart,
            'memory-chart': memoryChart,
            'disk-chart': diskChart,
            'history-chart': historyChart
        };

    } catch (error) {
        console.error('Error loading agent overview:', error);
    }
}

function generateMockHistory(currentData) {
    const labels = [];
    const cpu = [];
    const memory = [];
    const disk = [];

    const now = new Date();
    for (let i = 30; i >= 0; i--) {
        const time = new Date(now - i * 10000);
        labels.push(time.toLocaleTimeString());

        // Generate realistic fluctuating data around current values
        cpu.push(Math.max(0, Math.min(100, currentData.cpu_percent + (Math.random() - 0.5) * 20)));
        memory.push(Math.max(0, Math.min(100, currentData.memory_percent + (Math.random() - 0.5) * 10)));
        disk.push(Math.max(0, Math.min(100, currentData.disk_percent + (Math.random() - 0.5) * 5)));
    }

    return { labels, cpu, memory, disk };
}

async function executeCommand() {
    const command = document.getElementById('execute-command').value;
    if (!command || !currentAgentName) return;

    const output = document.getElementById('command-output');
    output.textContent = `$ ${command}\n\nExecuting...\n`;

    try {
        const response = await fetch(`/api/v1/agents/${currentAgentName}/command`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ command })
        });

        const reader = response.body.getReader();
        const decoder = new TextDecoder();

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const chunk = decoder.decode(value);
            output.textContent += chunk;
            output.scrollTop = output.scrollHeight;
        }
    } catch (error) {
        output.textContent += `\n\nError: ${error.message}`;
    }
}

function showBulkCommandModal() {
    const modal = new bootstrap.Modal(document.getElementById('bulkCommandModal'));
    modal.show();
}

async function runBulkCommand() {
    const command = document.getElementById('bulk-command').value;
    const parallel = document.getElementById('parallel-execution').checked;

    if (!command) {
        alert('Please enter a command');
        return;
    }

    const agentNames = Array.from(selectedAgents);
    if (agentNames.length === 0) {
        alert('Please select at least one agent');
        return;
    }

    try {
        const response = await fetch('/api/v1/agents/bulk/execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                agent_names: agentNames,
                command: command,
                parallel: parallel
            })
        });

        // Handle streaming response
        const reader = response.body.getReader();
        const decoder = new TextDecoder();

        console.log('Bulk execution started...');
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const chunk = decoder.decode(value);
            const results = chunk.split('\n').filter(line => line.trim());

            results.forEach(line => {
                try {
                    const result = JSON.parse(line);
                    console.log(`[${result.agent_name}] ${result.success ? '✓' : '✗'} ${result.output || result.error}`);
                } catch (e) {
                    // Ignore parsing errors
                }
            });
        }

        alert('Bulk command execution completed');
        bootstrap.Modal.getInstance(document.getElementById('bulkCommandModal')).hide();
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

function showCreateGroupModal() {
    const modal = new bootstrap.Modal(document.getElementById('createGroupModal'));
    modal.show();
}

async function createGroup() {
    const name = document.getElementById('group-name').value;
    const description = document.getElementById('group-description').value;

    if (!name) {
        alert('Please enter a group name');
        return;
    }

    const agentNames = Array.from(selectedAgents);

    try {
        const response = await fetch('/api/v1/agent-groups', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                group_name: name,
                description: description,
                agent_names: agentNames
            })
        });

        if (response.ok) {
            alert('Group created successfully');
            bootstrap.Modal.getInstance(document.getElementById('createGroupModal')).hide();
            clearSelection();
            loadAgents();
        } else {
            const error = await response.json();
            alert(`Error: ${error.message}`);
        }
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

async function restartAgent(agentName) {
    if (!confirm(`Restart agent ${agentName}?`)) return;

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/restart`, {
            method: 'POST'
        });

        if (response.ok) {
            alert(`Agent ${agentName} restart initiated`);
        } else {
            alert('Failed to restart agent');
        }
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

async function shutdownAgent(agentName) {
    if (!confirm(`Shutdown agent ${agentName}? This will stop the agent service.`)) return;

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/shutdown`, {
            method: 'POST'
        });

        if (response.ok) {
            alert(`Agent ${agentName} shutdown initiated`);
            loadAgents();
        } else {
            alert('Failed to shutdown agent');
        }
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

function quickCommand(agentName) {
    const command = prompt('Enter command to execute:');
    if (!command) return;

    currentAgentName = agentName;
    openAgentDetail(agentName);

    // Wait for modal to open then execute
    setTimeout(() => {
        document.querySelector('[data-bs-target="#command-tab"]').click();
        setTimeout(() => {
            document.getElementById('execute-command').value = command;
            executeCommand();
        }, 300);
    }, 500);
}

function refreshAgents() {
    loadAgents();
}

// Utility functions
function formatBytes(bytes) {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

function formatUptime(seconds) {
    if (!seconds) return '0s';
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const mins = Math.floor((seconds % 3600) / 60);

    if (days > 0) return `${days}d ${hours}h`;
    if (hours > 0) return `${hours}h ${mins}m`;
    return `${mins}m`;
}

function updateAgentStatus(agentName, status) {
    const agent = agents.find(a => a.name === agentName);
    if (agent) {
        agent.status = status;
        renderAgents();
        updateOverviewStats();
    }
}

function updateAgentMetrics(agentName, metrics) {
    const agent = agents.find(a => a.name === agentName);
    if (agent) {
        Object.assign(agent, metrics);
        renderAgents();
        updateOverviewStats();
    }
}

function appendLogEntry(log) {
    const logsStream = document.getElementById('logs-stream');
    if (logsStream) {
        const timestamp = new Date(log.timestamp * 1000).toLocaleTimeString();
        logsStream.innerHTML += `[${timestamp}] [${log.level}] ${log.message}\n`;
        logsStream.scrollTop = logsStream.scrollHeight;
    }
}
