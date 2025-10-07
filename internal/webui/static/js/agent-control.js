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

    // Add tab event listeners
    setupTabListeners();
});

function setupTabListeners() {
    // Wait for modal to be in DOM
    setTimeout(() => {
        const processesTab = document.getElementById('processes-tab-btn');
        const networkTab = document.getElementById('network-tab-btn');
        const diskTab = document.getElementById('disk-tab-btn');
        const logsTab = document.getElementById('logs-tab-btn');

        if (processesTab) {
            processesTab.addEventListener('shown.bs.tab', () => {
                if (currentAgentName) {
                    loadAgentProcesses(currentAgentName);
                }
            });
        }

        if (networkTab) {
            networkTab.addEventListener('shown.bs.tab', () => {
                if (currentAgentName) {
                    loadAgentNetwork(currentAgentName);
                }
            });
        }

        if (diskTab) {
            diskTab.addEventListener('shown.bs.tab', () => {
                if (currentAgentName) {
                    loadAgentDisk(currentAgentName);
                }
            });
        }

        if (logsTab) {
            logsTab.addEventListener('shown.bs.tab', () => {
                if (currentAgentName) {
                    loadAgentLogs(currentAgentName);
                }
            });
        }
    }, 100);
}

function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/api/v1/ws`);

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
            if (currentFilter === 'online') return agent.status === 'Active' || agent.status === 'connected';
            if (currentFilter === 'offline') return agent.status !== 'Active' && agent.status !== 'connected';
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
    const isOnline = agent.status === 'Active' || agent.status === 'connected';
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
                    <button class="btn btn-sm btn-outline-primary agent-action-btn flex-fill"
                            onclick="window.location.href='/agent-dashboard?agent=${agent.name}'" ${!isOnline ? 'disabled' : ''}>
                        <i class="bi bi-speedometer2"></i> Dashboard
                    </button>
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
    console.log('[DEBUG] Loading agent overview for:', agentName);
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/resources`);
        console.log('[DEBUG] API response status:', response.status);
        const data = await response.json();
        console.log('[DEBUG] API data received:', data);

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
                console.log('[DEBUG] Destroying old chart:', id);
                resourceCharts[id].destroy();
            }
        });

        // Check if Chart.js is loaded
        if (typeof Chart === 'undefined') {
            console.error('[ERROR] Chart.js is NOT loaded!');
            content.innerHTML += '<div class="alert alert-danger mt-3">Chart.js library not loaded. Cannot display charts.</div>';
            return;
        }
        console.log('[DEBUG] Chart.js is loaded, version:', Chart.version);

        // Create doughnut charts for resources
        console.log('[DEBUG] Creating CPU chart...');
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

        // Load real historical data from API
        const historyData = await loadMetricsHistory(agentName);

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

        console.log('[DEBUG] All charts created successfully:', Object.keys(resourceCharts));

    } catch (error) {
        console.error('[ERROR] Error loading agent overview:', error);
        console.error('[ERROR] Stack trace:', error.stack);
        const content = document.getElementById('agent-overview-content');
        if (content) {
            content.innerHTML = `<div class="alert alert-danger">Error: ${error.message}</div>`;
        }
    }
}

async function loadMetricsHistory(agentName) {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/metrics/history?limit=50`);
        if (!response.ok) {
            console.warn('Could not load metrics history, using mock data');
            return generateMockHistory({ cpu_percent: 0, memory_percent: 0, disk_percent: 0 });
        }

        const data = await response.json();
        const history = data.history || [];

        // If no history data, generate mock data
        if (history.length === 0) {
            console.warn('No metrics history available, using mock data');
            return generateMockHistory({ cpu_percent: 0, memory_percent: 0, disk_percent: 0 });
        }

        // Sort by timestamp ascending (oldest first)
        history.sort((a, b) => a.timestamp - b.timestamp);

        const labels = [];
        const cpu = [];
        const memory = [];
        const disk = [];

        history.forEach(point => {
            const date = new Date(point.timestamp * 1000);
            labels.push(date.toLocaleTimeString());
            cpu.push(point.cpu_percent || 0);
            memory.push(point.memory_percent || 0);
            disk.push(point.disk_percent || 0);
        });

        return { labels, cpu, memory, disk };
    } catch (error) {
        console.error('Error loading metrics history:', error);
        return generateMockHistory({ cpu_percent: 0, memory_percent: 0, disk_percent: 0 });
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

// Load agent processes when Processes tab is clicked
async function loadAgentProcesses(agentName) {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/processes`);
        const data = await response.json();

        const content = document.getElementById('agent-processes-content');

        if (!data.processes || data.processes.length === 0) {
            content.innerHTML = '<div class="alert alert-info">No processes found</div>';
            return;
        }

        content.innerHTML = `
            <table class="table table-sm table-hover">
                <thead>
                    <tr>
                        <th>PID</th>
                        <th>Name</th>
                        <th>CPU %</th>
                        <th>Memory %</th>
                        <th>Status</th>
                        <th>User</th>
                    </tr>
                </thead>
                <tbody>
                    ${data.processes.map(proc => `
                        <tr>
                            <td>${proc.pid}</td>
                            <td><code>${proc.name || 'N/A'}</code></td>
                            <td>${proc.cpu_percent?.toFixed(1) || '0.0'}%</td>
                            <td>${proc.memory_percent?.toFixed(1) || '0.0'}%</td>
                            <td><span class="badge bg-success">${proc.status || 'running'}</span></td>
                            <td>${proc.username || 'N/A'}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (error) {
        console.error('Error loading processes:', error);
        document.getElementById('agent-processes-content').innerHTML =
            `<div class="alert alert-danger">Error loading processes: ${error.message}</div>`;
    }
}

// Load agent network info when Network tab is clicked
async function loadAgentNetwork(agentName) {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/network`);
        const data = await response.json();

        const content = document.getElementById('agent-network-content');

        if (!data.interfaces || data.interfaces.length === 0) {
            content.innerHTML = '<div class="alert alert-info">No network interfaces found</div>';
            return;
        }

        content.innerHTML = `
            <div class="row">
                <div class="col-12">
                    <h6>Network Interfaces</h6>
                    ${data.interfaces.map(iface => `
                        <div class="card mb-3">
                            <div class="card-body">
                                <h6>${iface.name}</h6>
                                <div class="row">
                                    <div class="col-md-6">
                                        <table class="table table-sm">
                                            <tr>
                                                <td>IP Address:</td>
                                                <td><strong>${iface.addrs?.join(', ') || 'N/A'}</strong></td>
                                            </tr>
                                            <tr>
                                                <td>MAC Address:</td>
                                                <td><code>${iface.hardware_addr || 'N/A'}</code></td>
                                            </tr>
                                            <tr>
                                                <td>MTU:</td>
                                                <td>${iface.mtu || 'N/A'}</td>
                                            </tr>
                                            <tr>
                                                <td>Status:</td>
                                                <td><span class="badge bg-${iface.flags?.includes('up') ? 'success' : 'secondary'}">${iface.flags?.join(', ') || 'N/A'}</span></td>
                                            </tr>
                                        </table>
                                    </div>
                                    <div class="col-md-6">
                                        <table class="table table-sm">
                                            <tr>
                                                <td>Bytes Sent:</td>
                                                <td><strong>${formatBytes(iface.bytes_sent)}</strong></td>
                                            </tr>
                                            <tr>
                                                <td>Bytes Received:</td>
                                                <td><strong>${formatBytes(iface.bytes_recv)}</strong></td>
                                            </tr>
                                            <tr>
                                                <td>Packets Sent:</td>
                                                <td>${iface.packets_sent || 0}</td>
                                            </tr>
                                            <tr>
                                                <td>Packets Received:</td>
                                                <td>${iface.packets_recv || 0}</td>
                                            </tr>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    } catch (error) {
        console.error('Error loading network info:', error);
        document.getElementById('agent-network-content').innerHTML =
            `<div class="alert alert-danger">Error loading network info: ${error.message}</div>`;
    }
}

// Load agent disk info when Disk tab is clicked
async function loadAgentDisk(agentName) {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/disk`);
        const data = await response.json();

        const content = document.getElementById('agent-disk-content');

        if (!data.partitions || data.partitions.length === 0) {
            content.innerHTML = '<div class="alert alert-info">No disk partitions found</div>';
            return;
        }

        content.innerHTML = `
            <div class="row">
                ${data.partitions.map(partition => {
                    const usedPercent = partition.usage?.percent || 0;
                    const usedClass = usedPercent > 80 ? 'danger' : usedPercent > 60 ? 'warning' : 'success';

                    return `
                        <div class="col-md-6 mb-3">
                            <div class="card">
                                <div class="card-body">
                                    <h6><i class="bi bi-hdd"></i> ${partition.mountpoint}</h6>
                                    <table class="table table-sm">
                                        <tr>
                                            <td>Device:</td>
                                            <td><code>${partition.device || 'N/A'}</code></td>
                                        </tr>
                                        <tr>
                                            <td>Filesystem:</td>
                                            <td>${partition.fstype || 'N/A'}</td>
                                        </tr>
                                        <tr>
                                            <td>Total:</td>
                                            <td><strong>${formatBytes(partition.usage?.total)}</strong></td>
                                        </tr>
                                        <tr>
                                            <td>Used:</td>
                                            <td><strong>${formatBytes(partition.usage?.used)}</strong></td>
                                        </tr>
                                        <tr>
                                            <td>Free:</td>
                                            <td><strong>${formatBytes(partition.usage?.free)}</strong></td>
                                        </tr>
                                    </table>
                                    <div class="progress" style="height: 25px;">
                                        <div class="progress-bar bg-${usedClass}" style="width: ${usedPercent}%">
                                            ${usedPercent.toFixed(1)}%
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
        `;
    } catch (error) {
        console.error('Error loading disk info:', error);
        document.getElementById('agent-disk-content').innerHTML =
            `<div class="alert alert-danger">Error loading disk info: ${error.message}</div>`;
    }
}

// Load agent logs when Logs tab is clicked
async function loadAgentLogs(agentName) {
    const logsStream = document.getElementById('logs-stream');
    logsStream.textContent = '📡 Connecting to log stream...\n\n';

    try {
        // Close previous EventSource if exists
        if (window.agentLogStreams && window.agentLogStreams[agentName]) {
            window.agentLogStreams[agentName].close();
        }

        const url = `/api/v1/agents/${agentName}/logs/stream`;
        const eventSource = new EventSource(url);

        eventSource.onopen = () => {
            logsStream.innerHTML += '✅ Connected to log stream. Waiting for logs...\n\n';
        };

        eventSource.onmessage = (event) => {
            try {
                const log = JSON.parse(event.data);
                const timestamp = new Date(log.timestamp * 1000).toLocaleTimeString();
                const levelColors = {
                    'ERROR': '#EF4444',
                    'WARN': '#F59E0B',
                    'INFO': '#10B981',
                    'DEBUG': '#3B82F6'
                };
                const color = levelColors[log.level] || '#9CA3AF';
                logsStream.innerHTML += `<span style="color: #6B7280;">[${timestamp}]</span> <span style="color: ${color}; font-weight: bold;">[${log.level}]</span> ${escapeHtml(log.message)}\n`;
                logsStream.scrollTop = logsStream.scrollHeight;
            } catch (e) {
                // If not JSON, just append as plain text
                logsStream.innerHTML += escapeHtml(event.data) + '\n';
                logsStream.scrollTop = logsStream.scrollHeight;
            }
        };

        eventSource.onerror = (error) => {
            console.error('EventSource error:', error);
            logsStream.innerHTML += '\n❌ Log stream error or disconnected\n';
            eventSource.close();
        };

        // Store EventSource reference
        if (!window.agentLogStreams) {
            window.agentLogStreams = {};
        }
        window.agentLogStreams[agentName] = eventSource;

    } catch (error) {
        console.error('Error loading logs:', error);
        logsStream.textContent = `Error loading logs: ${error.message}`;
    }
}

// Helper function to escape HTML
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}
