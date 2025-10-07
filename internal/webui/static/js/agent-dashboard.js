// Individual Agent Dashboard JavaScript
let agentName = null;
let metricsChart = null;
let distributionChart = null;
let refreshInterval = null;

// Get agent name from URL parameter
function getAgentFromURL() {
    const params = new URLSearchParams(window.location.search);
    return params.get('agent');
}

// Initialize dashboard
document.addEventListener('DOMContentLoaded', async () => {
    agentName = getAgentFromURL();

    if (!agentName) {
        document.body.innerHTML = '<div class="container mt-5"><div class="alert alert-danger">No agent specified</div></div>';
        return;
    }

    // Set agent name in header
    document.getElementById('agent-name').textContent = agentName;

    // Load initial data
    await loadAgentData();
    await loadHistoricalMetrics();

    // Auto-refresh every 5 seconds
    refreshInterval = setInterval(async () => {
        await loadAgentData();
    }, 5000);

    // Setup tab switching
    document.querySelectorAll('[data-bs-toggle="tab"]').forEach(tab => {
        tab.addEventListener('shown.bs.tab', function (e) {
            const target = e.target.getAttribute('href');
            if (target === '#processes') {
                loadProcesses();
            } else if (target === '#network') {
                loadNetwork();
            } else if (target === '#disk') {
                loadDisk();
            } else if (target === '#system') {
                loadSystemInfo();
            }
        });
    });

    // Load processes tab initially
    await loadProcesses();
});

// Load agent data
async function loadAgentData() {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}`);
        const data = await response.json();

        // Update status
        const statusBadge = document.getElementById('agent-status');
        const addressSpan = document.getElementById('agent-address');

        if (data.status === 'active' || data.last_heartbeat) {
            statusBadge.className = 'badge bg-success';
            statusBadge.textContent = 'Active';
        } else {
            statusBadge.className = 'badge bg-secondary';
            statusBadge.textContent = 'Inactive';
        }

        addressSpan.textContent = data.address || '';

        // Parse system_info if available
        let systemInfo = {};
        if (data.system_info) {
            try {
                systemInfo = JSON.parse(data.system_info);
            } catch (e) {
                console.error('Failed to parse system_info:', e);
            }
        }

        // Update metric cards
        updateMetricCard('cpu', systemInfo.cpu);
        updateMetricCard('memory', systemInfo.memory);
        updateMetricCard('disk', systemInfo.disk);
        updateMetricCard('load', systemInfo.load_average);

        // Update distribution chart
        updateDistributionChart(systemInfo);

    } catch (error) {
        console.error('Failed to load agent data:', error);
    }
}

// Update metric card
function updateMetricCard(type, data) {
    if (!data) return;

    const valueEl = document.getElementById(`${type}-value`);
    const progressEl = document.getElementById(`${type}-progress`);

    switch (type) {
        case 'cpu':
            const cpuPercent = data.usage_percent || 0;
            valueEl.textContent = cpuPercent.toFixed(1) + '%';
            progressEl.style.width = cpuPercent + '%';
            break;

        case 'memory':
            const memPercent = data.used_percent || 0;
            valueEl.textContent = memPercent.toFixed(1) + '%';
            progressEl.style.width = memPercent + '%';
            break;

        case 'disk':
            if (data.partitions && data.partitions.length > 0) {
                const avgPercent = data.partitions.reduce((sum, p) => sum + (p.percent || 0), 0) / data.partitions.length;
                valueEl.textContent = avgPercent.toFixed(1) + '%';
                progressEl.style.width = avgPercent + '%';
            }
            break;

        case 'load':
            const load1 = data['1min'] || 0;
            valueEl.textContent = load1.toFixed(2);
            const loadDetails = document.getElementById('load-details');
            if (loadDetails) {
                loadDetails.textContent = `${(data['1min'] || 0).toFixed(2)} / ${(data['5min'] || 0).toFixed(2)} / ${(data['15min'] || 0).toFixed(2)}`;
            }
            break;
    }
}

// Load historical metrics
async function loadHistoricalMetrics() {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/metrics/history?limit=50`);
        const data = await response.json();

        const history = data.history || [];
        history.sort((a, b) => a.timestamp - b.timestamp);

        const labels = history.map(h => {
            const date = new Date(h.timestamp * 1000);
            return date.toLocaleTimeString();
        });

        const cpuData = history.map(h => h.cpu_percent || 0);
        const memoryData = history.map(h => h.memory_percent || 0);
        const diskData = history.map(h => h.disk_percent || 0);

        // Create or update historical metrics chart
        const ctx = document.getElementById('history-chart');
        if (metricsChart) {
            metricsChart.data.labels = labels;
            metricsChart.data.datasets[0].data = cpuData;
            metricsChart.data.datasets[1].data = memoryData;
            metricsChart.data.datasets[2].data = diskData;
            metricsChart.update();
        } else {
            metricsChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [
                        {
                            label: 'CPU %',
                            data: cpuData,
                            borderColor: 'rgb(13, 110, 253)',
                            backgroundColor: 'rgba(13, 110, 253, 0.1)',
                            tension: 0.4
                        },
                        {
                            label: 'Memory %',
                            data: memoryData,
                            borderColor: 'rgb(25, 135, 84)',
                            backgroundColor: 'rgba(25, 135, 84, 0.1)',
                            tension: 0.4
                        },
                        {
                            label: 'Disk %',
                            data: diskData,
                            borderColor: 'rgb(255, 193, 7)',
                            backgroundColor: 'rgba(255, 193, 7, 0.1)',
                            tension: 0.4
                        }
                    ]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            display: true,
                            position: 'top'
                        }
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                            max: 100
                        }
                    }
                }
            });
        }
    } catch (error) {
        console.error('Failed to load historical metrics:', error);
    }
}

// Update distribution chart
function updateDistributionChart(systemInfo) {
    const cpuPercent = systemInfo.cpu?.usage_percent || 0;
    const memPercent = systemInfo.memory?.used_percent || 0;

    let diskPercent = 0;
    if (systemInfo.disk?.partitions && systemInfo.disk.partitions.length > 0) {
        diskPercent = systemInfo.disk.partitions.reduce((sum, p) => sum + (p.percent || 0), 0) / systemInfo.disk.partitions.length;
    }

    const ctx = document.getElementById('distribution-chart');

    if (distributionChart) {
        distributionChart.data.datasets[0].data = [cpuPercent, memPercent, diskPercent];
        distributionChart.update();
    } else {
        distributionChart = new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: ['CPU Usage', 'Memory Usage', 'Disk Usage'],
                datasets: [{
                    data: [cpuPercent, memPercent, diskPercent],
                    backgroundColor: [
                        'rgba(13, 110, 253, 0.8)',
                        'rgba(25, 135, 84, 0.8)',
                        'rgba(255, 193, 7, 0.8)'
                    ],
                    borderWidth: 2,
                    borderColor: '#fff'
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true,
                        position: 'bottom'
                    }
                }
            }
        });
    }
}

// Load processes
async function loadProcesses() {
    const container = document.getElementById('processes-content');
    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/processes`);
        const data = await response.json();

        if (data.processes && data.processes.length > 0) {
            const table = `
                <table class="table table-sm table-hover">
                    <thead>
                        <tr>
                            <th>PID</th>
                            <th>Name</th>
                            <th>CPU %</th>
                            <th>Memory %</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${data.processes.slice(0, 20).map(p => `
                            <tr>
                                <td>${p.pid}</td>
                                <td>${p.name || '-'}</td>
                                <td>${(p.cpu_percent || 0).toFixed(1)}%</td>
                                <td>${(p.memory_percent || 0).toFixed(1)}%</td>
                                <td><span class="badge bg-success">${p.status || 'running'}</span></td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
            container.innerHTML = table;
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No process data available</div>';
        }
    } catch (error) {
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load processes</div>';
    }
}

// Load network info
async function loadNetwork() {
    const container = document.getElementById('network-content');
    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/network`);
        const data = await response.json();

        if (data.interfaces && data.interfaces.length > 0) {
            const cards = data.interfaces.map(iface => `
                <div class="col-md-6 mb-3">
                    <div class="card">
                        <div class="card-body">
                            <h6 class="card-title">${iface.name}</h6>
                            <p class="mb-1"><strong>IP:</strong> ${iface.ip || '-'}</p>
                            <p class="mb-1"><strong>MAC:</strong> ${iface.mac || '-'}</p>
                            <p class="mb-0"><strong>Sent:</strong> ${formatBytes(iface.bytes_sent || 0)} | <strong>Received:</strong> ${formatBytes(iface.bytes_recv || 0)}</p>
                        </div>
                    </div>
                </div>
            `).join('');

            container.innerHTML = `<div class="row">${cards}</div>`;
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No network data available</div>';
        }
    } catch (error) {
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load network info</div>';
    }
}

// Load disk info
async function loadDisk() {
    const container = document.getElementById('disk-content');
    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/disk`);
        const data = await response.json();

        if (data.partitions && data.partitions.length > 0) {
            const cards = data.partitions.map(part => `
                <div class="col-md-6 mb-3">
                    <div class="card">
                        <div class="card-body">
                            <h6 class="card-title">${part.mountpoint || part.device}</h6>
                            <div class="progress mb-2" style="height: 20px;">
                                <div class="progress-bar" role="progressbar" style="width: ${part.percent}%"
                                     aria-valuenow="${part.percent}" aria-valuemin="0" aria-valuemax="100">
                                    ${part.percent.toFixed(1)}%
                                </div>
                            </div>
                            <p class="mb-0">
                                <strong>Used:</strong> ${formatBytes(part.used || 0)} / ${formatBytes(part.total || 0)}
                                <strong class="ms-3">Free:</strong> ${formatBytes(part.free || 0)}
                            </p>
                        </div>
                    </div>
                </div>
            `).join('');

            container.innerHTML = `<div class="row">${cards}</div>`;
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No disk data available</div>';
        }
    } catch (error) {
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load disk info</div>';
    }
}

// Load system info
async function loadSystemInfo() {
    const container = document.getElementById('system-content');
    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}`);
        const data = await response.json();

        let systemInfo = {};
        if (data.system_info) {
            try {
                systemInfo = JSON.parse(data.system_info);
            } catch (e) {
                console.error('Failed to parse system_info:', e);
            }
        }

        const table = `
            <table class="table table-bordered">
                <tbody>
                    <tr>
                        <th width="30%">Agent Name</th>
                        <td>${data.name || '-'}</td>
                    </tr>
                    <tr>
                        <th>Version</th>
                        <td>${data.version || '-'}</td>
                    </tr>
                    <tr>
                        <th>Address</th>
                        <td>${data.address || '-'}</td>
                    </tr>
                    <tr>
                        <th>Platform</th>
                        <td>${systemInfo.platform || '-'}</td>
                    </tr>
                    <tr>
                        <th>Architecture</th>
                        <td>${systemInfo.arch || '-'}</td>
                    </tr>
                    <tr>
                        <th>CPU Cores</th>
                        <td>${systemInfo.cpu?.cores || '-'}</td>
                    </tr>
                    <tr>
                        <th>Total Memory</th>
                        <td>${formatBytes(systemInfo.memory?.total || 0)}</td>
                    </tr>
                    <tr>
                        <th>Hostname</th>
                        <td>${systemInfo.hostname || '-'}</td>
                    </tr>
                    <tr>
                        <th>Last Heartbeat</th>
                        <td>${data.last_heartbeat || '-'}</td>
                    </tr>
                </tbody>
            </table>
        `;

        container.innerHTML = table;
    } catch (error) {
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load system info</div>';
    }
}

// Format bytes
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    if (refreshInterval) {
        clearInterval(refreshInterval);
    }
    if (metricsChart) {
        metricsChart.destroy();
    }
    if (distributionChart) {
        distributionChart.destroy();
    }
});
