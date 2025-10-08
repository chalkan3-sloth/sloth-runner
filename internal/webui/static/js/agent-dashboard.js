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
    console.log('[Dashboard] DOM loaded, initializing...');
    agentName = getAgentFromURL();
    console.log('[Dashboard] Agent name from URL:', agentName);

    if (!agentName) {
        console.error('[Dashboard] No agent name provided in URL');
        document.body.innerHTML = '<div class="container mt-5"><div class="alert alert-danger">No agent specified. Use ?agent=AGENT_NAME</div></div>';
        return;
    }

    // Set agent name in header
    const agentNameEl = document.getElementById('agent-name');
    if (agentNameEl) {
        agentNameEl.textContent = agentName;
        console.log('[Dashboard] Set agent name in header');
    } else {
        console.error('[Dashboard] Element #agent-name not found');
    }

    // Load initial data
    await loadAgentData();
    await loadHistoricalMetrics();

    // Auto-refresh every 5 seconds
    refreshInterval = setInterval(async () => {
        await loadAgentData();
    }, 5000);

    // Setup tab switching with direct click handling
    setupTabSwitching();

    // Load processes tab initially
    await loadProcesses();
});

// Load agent data
async function loadAgentData() {
    console.log('[Dashboard] loadAgentData() called for agent:', agentName);
    try {
        // Fetch both agent info and stats
        console.log('[Dashboard] Fetching agent data...');
        const [agentResponse, statsResponse] = await Promise.all([
            fetch(`/api/v1/agents/${agentName}`),
            fetch(`/api/v1/agents/${agentName}/stats`)
        ]);

        console.log('[Dashboard] Agent response status:', agentResponse.status);
        console.log('[Dashboard] Stats response status:', statsResponse.status);

        const agentData = await agentResponse.json();
        const statsData = await statsResponse.json();

        console.log('[Dashboard] Agent data:', agentData);
        console.log('[Dashboard] Stats data:', statsData);

        // Update status
        const statusBadge = document.getElementById('agent-status');
        const addressSpan = document.getElementById('agent-address');

        if (agentData.status === 'active' || agentData.last_heartbeat) {
            statusBadge.className = 'badge bg-success';
            statusBadge.textContent = 'Active';
        } else {
            statusBadge.className = 'badge bg-secondary';
            statusBadge.textContent = 'Inactive';
        }

        addressSpan.textContent = agentData.address || '';

        // Build systemInfo object from stats data
        const systemInfo = {
            cpu: {
                usage_percent: statsData.cpu_percent || 0
            },
            memory: {
                used_percent: statsData.memory_percent || 0
            },
            disk: {
                partitions: [{
                    percent: statsData.disk_percent || 0
                }]
            },
            load_average: {
                '1min': statsData.load_avg?.[0] || 0,
                '5min': statsData.load_avg?.[1] || 0,
                '15min': statsData.load_avg?.[2] || 0
            }
        };

        // Hide all spinners first
        document.querySelectorAll('.spinner-border').forEach(spinner => {
            spinner.style.display = 'none';
        });

        // Update metric cards
        updateMetricCard('cpu', systemInfo.cpu);
        updateMetricCard('memory', systemInfo.memory);
        updateMetricCard('disk', systemInfo.disk);
        updateMetricCard('load', systemInfo.load_average);

        // Update distribution chart
        updateDistributionChart(systemInfo);

        // Populate Overview tab details
        populateCPUDetails(statsData);
        populateMemoryDetails(statsData);

    } catch (error) {
        console.error('Failed to load agent data:', error);
    }
}

// Update metric card
function updateMetricCard(type, data) {
    console.log(`[Dashboard] updateMetricCard called for ${type}`, data);
    if (!data) {
        console.warn(`[Dashboard] No data for ${type}`);
        return;
    }

    const valueEl = document.getElementById(`${type}-value`);
    const progressEl = document.getElementById(`${type}-progress`);

    console.log(`[Dashboard] Elements for ${type}:`, {valueEl, progressEl});

    if (!valueEl) {
        console.error(`[Dashboard] Element #${type}-value not found`);
        return;
    }

    switch (type) {
        case 'cpu':
            const cpuPercent = data.usage_percent || 0;
            valueEl.textContent = cpuPercent.toFixed(1) + '%';
            if (progressEl) progressEl.style.width = cpuPercent + '%';
            console.log(`[Dashboard] Updated CPU: ${cpuPercent.toFixed(1)}%`);
            break;

        case 'memory':
            const memPercent = data.used_percent || 0;
            valueEl.textContent = memPercent.toFixed(1) + '%';
            if (progressEl) progressEl.style.width = memPercent + '%';
            console.log(`[Dashboard] Updated Memory: ${memPercent.toFixed(1)}%`);
            break;

        case 'disk':
            if (data.partitions && data.partitions.length > 0) {
                const avgPercent = data.partitions.reduce((sum, p) => sum + (p.percent || 0), 0) / data.partitions.length;
                valueEl.textContent = avgPercent.toFixed(1) + '%';
                if (progressEl) progressEl.style.width = avgPercent + '%';
                console.log(`[Dashboard] Updated Disk: ${avgPercent.toFixed(1)}%`);
            }
            break;

        case 'load':
            const load1 = data['1min'] || 0;
            valueEl.textContent = load1.toFixed(2);
            const loadDetails = document.getElementById('load-details');
            if (loadDetails) {
                loadDetails.textContent = `${(data['1min'] || 0).toFixed(2)} / ${(data['5min'] || 0).toFixed(2)} / ${(data['15min'] || 0).toFixed(2)}`;
            }
            console.log(`[Dashboard] Updated Load: ${load1.toFixed(2)}`);
            break;
    }

    // Check computed styles for visibility debugging
    const computedStyle = window.getComputedStyle(valueEl);
    console.log(`[Dashboard] Computed style for ${type}-value:`, {
        display: computedStyle.display,
        visibility: computedStyle.visibility,
        opacity: computedStyle.opacity,
        color: computedStyle.color,
        fontSize: computedStyle.fontSize,
        position: computedStyle.position,
        zIndex: computedStyle.zIndex,
        offsetParent: valueEl.offsetParent !== null ? 'visible' : 'HIDDEN',
        textContent: valueEl.textContent,
        innerHTML: valueEl.innerHTML
    });
}

// Setup tab switching with direct handling
function setupTabSwitching() {
    console.log('[Tabs] Setting up tab switching...');

    // Get all tab links
    const tabLinks = document.querySelectorAll('a[data-bs-toggle="tab"]');
    console.log('[Tabs] Found', tabLinks.length, 'tab links');

    // Store active tab
    let activeTab = '#overview';

    // Function to load tab content
    function loadTabContent(tabId) {
        console.log('[Tabs] Loading content for tab:', tabId);

        switch(tabId) {
            case '#processes':
                loadProcesses();
                break;
            case '#network':
                loadNetwork();
                break;
            case '#disk':
                loadDisk();
                break;
            case '#system':
                loadSystemInfo();
                break;
            case '#connections':
                loadConnections();
                break;
            case '#logs':
                loadLogs();
                break;
            case '#errors':
                loadErrors();
                break;
            case '#health':
                loadHealthDiagnostics();
                break;
            case '#performance':
                console.log('[Tabs] Loading performance history...');
                loadPerformanceHistory();
                break;
            case '#metrics':
                loadDetailedMetrics();
                break;
            case '#overview':
                console.log('[Tabs] Overview tab selected');
                break;
            default:
                console.log('[Tabs] Unknown tab:', tabId);
        }
    }

    // Attach click handlers to each tab
    tabLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            // Get the target tab ID
            const targetTab = this.getAttribute('href');
            console.log('[Tabs] Tab clicked:', targetTab);

            // Only load content if switching to a different tab
            if (targetTab !== activeTab) {
                activeTab = targetTab;

                // Wait a moment for Bootstrap to complete the tab switch
                setTimeout(() => {
                    loadTabContent(targetTab);
                }, 150);
            }
        });
    });

    // Also listen for Bootstrap tab shown event as backup
    const tabElements = document.querySelectorAll('button[data-bs-toggle="tab"], a[data-bs-toggle="tab"]');
    tabElements.forEach(tab => {
        tab.addEventListener('shown.bs.tab', function(event) {
            const targetTab = event.target.getAttribute('href') || event.target.getAttribute('data-bs-target');
            console.log('[Tabs] Bootstrap tab shown event:', targetTab);
            if (targetTab && targetTab !== activeTab) {
                activeTab = targetTab;
                loadTabContent(targetTab);
            }
        });
    });

    console.log('[Tabs] Tab switching setup complete');
}

// Load historical metrics
async function loadHistoricalMetrics() {
    try {
        const response = await fetch(`/api/v1/agents/${agentName}/metrics/history?duration=1h&maxPoints=50`);
        const data = await response.json();

        const history = data.datapoints || [];
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
let allProcesses = [];
async function loadProcesses() {
    const container = document.getElementById('processes-content');
    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/processes`);
        const data = await response.json();

        if (data.processes && data.processes.length > 0) {
            allProcesses = data.processes; // Store all processes for filtering

            // Add filter input and table
            const html = `
                <div class="mb-3">
                    <input type="text" class="form-control" id="process-filter" placeholder="Filter processes by name or PID...">
                </div>
                <div id="process-table-container"></div>
            `;
            container.innerHTML = html;

            // Render initial table
            renderProcessTable(allProcesses.slice(0, 20));

            // Setup filter
            const filterInput = document.getElementById('process-filter');
            if (filterInput) {
                filterInput.addEventListener('input', function(e) {
                    const filterValue = e.target.value.toLowerCase();
                    const filtered = allProcesses.filter(p => {
                        const name = (p.name || '').toLowerCase();
                        const pid = String(p.pid);
                        return name.includes(filterValue) || pid.includes(filterValue);
                    });
                    renderProcessTable(filtered.slice(0, 20));
                });
            }
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No process data available</div>';
        }
    } catch (error) {
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load processes</div>';
    }
}

// Helper function to render process table
function renderProcessTable(processes) {
    const container = document.getElementById('process-table-container');
    if (!container) return;

    if (processes.length === 0) {
        container.innerHTML = '<div class="text-muted text-center py-3">No processes match your filter</div>';
        return;
    }

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
                ${processes.map(p => `
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
}

// Load network info
async function loadNetwork() {
    const container = document.getElementById('network-content');
    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/network`);
        const data = await response.json();

        if (data.interfaces && data.interfaces.length > 0) {
            const cards = data.interfaces.map(iface => {
                // Format IP addresses (array)
                const ipAddresses = iface.ip_addresses && iface.ip_addresses.length > 0
                    ? iface.ip_addresses.join(', ')
                    : '-';

                // Status badge
                const statusBadge = iface.is_up
                    ? '<span class="badge bg-success ms-2">UP</span>'
                    : '<span class="badge bg-secondary ms-2">DOWN</span>';

                return `
                    <div class="col-md-6 mb-3">
                        <div class="card">
                            <div class="card-body">
                                <h6 class="card-title">${iface.name || '-'}${statusBadge}</h6>
                                <p class="mb-1"><strong>IP:</strong> ${ipAddresses}</p>
                                <p class="mb-1"><strong>MAC:</strong> ${iface.mac_address || '-'}</p>
                                <p class="mb-0"><strong>Sent:</strong> ${formatBytes(iface.bytes_sent || 0)} | <strong>Received:</strong> ${formatBytes(iface.bytes_recv || 0)}</p>
                            </div>
                        </div>
                    </div>
                `;
            }).join('');

            container.innerHTML = `
                <div class="row">
                    ${data.hostname ? `<div class="col-12 mb-3"><div class="alert alert-info"><strong>Hostname:</strong> ${data.hostname}</div></div>` : ''}
                    ${cards}
                </div>
            `;
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No network data available</div>';
        }
    } catch (error) {
        console.error('Failed to load network info:', error);
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
        // Fetch both agent info and detailed metrics for system information
        const [agentResponse, detailedResponse, networkResponse] = await Promise.all([
            fetch(`/api/v1/agents/${agentName}`),
            fetch(`/api/v1/agents/${agentName}/metrics/detailed`),
            fetch(`/api/v1/agents/${agentName}/network`)
        ]);

        const data = await agentResponse.json();
        const detailed = await detailedResponse.json();
        const network = await networkResponse.json();

        // Try to parse system_info if it exists (for backward compatibility)
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
                        <td>${systemInfo.platform || data.platform || '-'}</td>
                    </tr>
                    <tr>
                        <th>Architecture</th>
                        <td>${systemInfo.arch || data.arch || '-'}</td>
                    </tr>
                    <tr>
                        <th>CPU Cores</th>
                        <td>${detailed.cpu?.core_count || systemInfo.cpu?.cores || data.cpu_count || '-'}</td>
                    </tr>
                    <tr>
                        <th>CPU Model</th>
                        <td>${detailed.cpu?.model_name || '-'}</td>
                    </tr>
                    <tr>
                        <th>Total Memory</th>
                        <td>${formatBytes(detailed.memory?.total || systemInfo.memory?.total || 0)}</td>
                    </tr>
                    <tr>
                        <th>Hostname</th>
                        <td>${network.hostname || systemInfo.hostname || '-'}</td>
                    </tr>
                    <tr>
                        <th>Last Heartbeat</th>
                        <td>${data.last_heartbeat || '-'}</td>
                    </tr>
                    <tr>
                        <th>Status</th>
                        <td>${data.status === 'active' || data.last_heartbeat ? '<span class="badge bg-success">Active</span>' : '<span class="badge bg-secondary">Inactive</span>'}</td>
                    </tr>
                </tbody>
            </table>
        `;

        container.innerHTML = table;
    } catch (error) {
        console.error('Failed to load system info:', error);
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

// Escape HTML
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return String(text).replace(/[&<>"']/g, m => map[m]);
}

// Load connections
async function loadConnections() {
    const container = document.getElementById('connections-content');
    if (!container) return;

    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/connections`);
        const data = await response.json();

        if (data.connections && data.connections.length > 0) {
            // Update counters
            const totalConnections = data.connections.length;
            const established = data.connections.filter(c => c.state === 'ESTABLISHED').length;
            const listening = data.connections.filter(c => c.state === 'LISTEN').length;
            const timeWait = data.connections.filter(c => c.state === 'TIME_WAIT').length;

            const connTotalEl = document.getElementById('conn-total');
            const connEstablishedEl = document.getElementById('conn-established');
            const connListeningEl = document.getElementById('conn-listening');
            const connTimeWaitEl = document.getElementById('conn-time-wait');

            if (connTotalEl) connTotalEl.textContent = totalConnections;
            if (connEstablishedEl) connEstablishedEl.textContent = established;
            if (connListeningEl) connListeningEl.textContent = listening;
            if (connTimeWaitEl) connTimeWaitEl.textContent = timeWait;

            // Populate table
            const table = `
                <table class="table table-sm table-hover">
                    <thead>
                        <tr>
                            <th>Local Address</th>
                            <th>Remote Address</th>
                            <th>State</th>
                            <th>PID</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${data.connections.map(conn => `
                            <tr>
                                <td>${conn.local_addr || '-'}:${conn.local_port || '-'}</td>
                                <td>${conn.remote_addr || '-'}:${conn.remote_port || '-'}</td>
                                <td><span class="badge bg-${conn.state === 'ESTABLISHED' ? 'success' : conn.state === 'LISTEN' ? 'info' : 'secondary'}">${conn.state || '-'}</span></td>
                                <td>${conn.pid || '-'}</td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
            container.innerHTML = table;
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No connection data available</div>';
        }
    } catch (error) {
        console.error('Failed to load connections:', error);
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load connections</div>';
    }
}

// Load logs
async function loadLogs() {
    const container = document.getElementById('logs-content');
    if (!container) return;

    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const maxLines = document.getElementById('log-max-lines')?.value || 100;
        const response = await fetch(`/api/v1/agents/${agentName}/logs?max_lines=${maxLines}`);
        const data = await response.json();

        if (data.logs && data.logs.length > 0) {
            const logEntries = data.logs.map(log => {
                let badgeClass = 'bg-primary';
                if (log.level === 'warning') badgeClass = 'bg-warning';
                else if (log.level === 'error') badgeClass = 'bg-danger';
                else if (log.level === 'debug') badgeClass = 'bg-secondary';

                return `
                    <div class="log-entry mb-2 p-2 border-bottom">
                        <div class="d-flex align-items-center">
                            <small class="text-muted me-2">${log.timestamp || ''}</small>
                            <span class="badge ${badgeClass} me-2">${log.level || 'info'}</span>
                            <small class="text-muted">${log.source || ''}</small>
                        </div>
                        <div class="mt-1">${log.message || ''}</div>
                    </div>
                `;
            }).join('');

            container.innerHTML = logEntries;
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No logs available</div>';
        }
    } catch (error) {
        console.error('Failed to load logs:', error);
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load logs</div>';
    }
}

// Load errors
async function loadErrors() {
    const container = document.getElementById('errors-content');
    if (!container) return;

    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const includeWarnings = document.getElementById('include-warnings')?.checked !== false;
        const response = await fetch(`/api/v1/agents/${agentName}/errors?include_warnings=${includeWarnings}`);
        const data = await response.json();

        if (data.errors && data.errors.length > 0) {
            // Find most common error
            const errorMessages = data.errors.filter(e => e.severity === 'error').map(e => e.message);
            const mostCommon = errorMessages.reduce((acc, msg) => {
                acc[msg] = (acc[msg] || 0) + 1;
                return acc;
            }, {});
            const mostCommonError = Object.keys(mostCommon).length > 0
                ? Object.keys(mostCommon).reduce((a, b) => mostCommon[a] > mostCommon[b] ? a : b)
                : 'None';

            const mostCommonEl = document.getElementById('most-common-error');
            const errorSummaryEl = document.getElementById('error-summary');
            if (mostCommonEl) {
                mostCommonEl.textContent = mostCommonError;
            }
            if (errorSummaryEl && mostCommonError !== 'None') {
                errorSummaryEl.style.display = 'block';
            }

            // Populate table
            const table = `
                <table class="table table-sm table-hover">
                    <thead>
                        <tr>
                            <th>Timestamp</th>
                            <th>Severity</th>
                            <th>Source</th>
                            <th>Message</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${data.errors.map(err => `
                            <tr>
                                <td><small>${err.timestamp || '-'}</small></td>
                                <td><span class="badge bg-${err.severity === 'error' ? 'danger' : 'warning'}">${err.severity || '-'}</span></td>
                                <td><small>${err.source || '-'}</small></td>
                                <td>${err.message || '-'}</td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
            container.innerHTML = table;
        } else {
            container.innerHTML = '<div class="text-muted text-center py-4">No errors available</div>';
        }
    } catch (error) {
        console.error('Failed to load errors:', error);
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load errors</div>';
    }
}

// Load health diagnostics
async function loadHealthDiagnostics() {
    const healthScoreCard = document.getElementById('health-score');
    const scoreValue = document.getElementById('score-value');
    const healthStatus = document.getElementById('health-status');
    const healthIssues = document.getElementById('health-issues');

    if (!healthIssues) return;

    healthIssues.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const deepCheck = document.getElementById('deep-check')?.checked || false;
        const response = await fetch(`/api/v1/agents/${agentName}/health/diagnose?deep_check=${deepCheck}`);
        const data = await response.json();

        // Show health score card
        if (healthScoreCard) {
            healthScoreCard.style.display = 'block';
        }

        // Update health score
        if (scoreValue) {
            scoreValue.textContent = data.health_score || 0;
        }

        // Update health status
        const status = data.health_status || 'unknown';
        let statusBadgeClass = 'bg-secondary';
        if (status === 'healthy') statusBadgeClass = 'bg-success';
        else if (status === 'degraded') statusBadgeClass = 'bg-warning';
        else if (status === 'unhealthy') statusBadgeClass = 'bg-danger';

        if (healthStatus) {
            healthStatus.innerHTML = `<span class="badge ${statusBadgeClass}">${status}</span>`;
        }

        // Populate health issues
        if (data.issues && data.issues.length > 0) {
            const issueCards = data.issues.map(issue => {
                let severityBadge = 'bg-info';
                if (issue.severity === 'critical') severityBadge = 'bg-danger';
                else if (issue.severity === 'warning') severityBadge = 'bg-warning';

                const suggestions = issue.suggestions && issue.suggestions.length > 0
                    ? `<ul class="mb-0 mt-2">${issue.suggestions.map(s => `<li>${s}</li>`).join('')}</ul>`
                    : '';

                return `
                    <div class="card mb-3">
                        <div class="card-body">
                            <div class="d-flex justify-content-between align-items-start">
                                <h6 class="card-title">${issue.category || 'Issue'}</h6>
                                <span class="badge ${severityBadge}">${issue.severity || 'info'}</span>
                            </div>
                            <p class="mb-1">${issue.description || ''}</p>
                            ${issue.current_value ? `<small class="text-muted">Current: ${issue.current_value}</small>` : ''}
                            ${issue.threshold ? `<small class="text-muted"> | Threshold: ${issue.threshold}</small>` : ''}
                            ${suggestions}
                        </div>
                    </div>
                `;
            }).join('');

            healthIssues.innerHTML = issueCards;
        } else {
            healthIssues.innerHTML = '<div class="alert alert-success">No health issues detected</div>';
        }
    } catch (error) {
        console.error('Failed to load health diagnostics:', error);
        healthIssues.innerHTML = '<div class="alert alert-warning">Failed to load health diagnostics</div>';
    }
}

// Alias for HTML onclick
function runHealthCheck() {
    loadHealthDiagnostics();
}

// Load performance history
let performanceChart = null;

async function loadPerformanceHistory() {
    console.log('[Performance] loadPerformanceHistory() called for agent:', agentName);
    const container = document.getElementById('performance-chart-container');
    if (!container) {
        console.error('[Performance] Container #performance-chart-container not found');
        return;
    }

    // Show loading state
    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div><p class="mt-2">Loading performance history...</p></div>';

    try {
        // Check if Chart.js is loaded
        if (typeof Chart === 'undefined') {
            console.error('[Performance] Chart.js library not loaded!');
            container.innerHTML = '<div class="alert alert-danger m-3">Chart library not available</div>';
            return;
        }
        console.log('[Performance] Chart.js is available:', Chart);

        // Fetch 1 hour of metrics history
        console.log('[Performance] Fetching metrics from API...');
        const response = await fetch(`/api/v1/agents/${agentName}/metrics/history?duration=1h&maxPoints=60`);
        const data = await response.json();
        console.log('[Performance] Received data:', data);

        if (!data.datapoints || data.datapoints.length === 0) {
            container.innerHTML = `
                <div class="alert alert-info m-3">
                    <i class="bi bi-info-circle"></i> No performance history available yet.
                    <br><small>Metrics are collected every 30 seconds. Please wait a few minutes.</small>
                </div>
            `;
            return;
        }

        // Prepare chart data
        console.log('[Performance] Preparing chart data...');
        const labels = data.datapoints.map(dp => {
            const date = new Date(dp.timestamp * 1000);
            return date.toLocaleTimeString();
        });

        const cpuData = data.datapoints.map(dp => dp.cpu_percent);
        const memoryData = data.datapoints.map(dp => dp.memory_percent);
        const diskData = data.datapoints.map(dp => dp.disk_percent);
        console.log('[Performance] Labels:', labels.length, 'CPU points:', cpuData.length);

        // Create canvas
        console.log('[Performance] Creating canvas element...');
        container.innerHTML = '<canvas id="performanceChart"></canvas>';
        const canvas = document.getElementById('performanceChart');
        if (!canvas) {
            console.error('[Performance] Failed to create canvas element!');
            return;
        }
        const ctx = canvas.getContext('2d');
        console.log('[Performance] Canvas context:', ctx);

        // Create chart
        console.log('[Performance] Creating Chart.js chart...');
        const chart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'CPU %',
                        data: cpuData,
                        borderColor: 'rgb(255, 99, 132)',
                        backgroundColor: 'rgba(255, 99, 132, 0.1)',
                        tension: 0.4
                    },
                    {
                        label: 'Memory %',
                        data: memoryData,
                        borderColor: 'rgb(54, 162, 235)',
                        backgroundColor: 'rgba(54, 162, 235, 0.1)',
                        tension: 0.4
                    },
                    {
                        label: 'Disk %',
                        data: diskData,
                        borderColor: 'rgb(255, 206, 86)',
                        backgroundColor: 'rgba(255, 206, 86, 0.1)',
                        tension: 0.4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'top',
                    },
                    title: {
                        display: true,
                        text: 'Resource Usage (Last Hour)'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            callback: function(value) {
                                return value + '%';
                            }
                        }
                    }
                }
            }
        });
        console.log('[Performance] Chart created successfully:', chart);

    } catch (error) {
        console.error('[Performance] Error in loadPerformanceHistory:', error);
        console.error('[Performance] Error stack:', error.stack);
        container.innerHTML = `<div class="alert alert-danger m-3">Failed to load performance history: ${error.message}</div>`;
    }
}

// Load detailed metrics
async function loadDetailedMetrics() {
    const container = document.getElementById('detailed-metrics-container');
    if (!container) return;

    container.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';

    try {
        const response = await fetch(`/api/v1/agents/${agentName}/metrics/detailed`);
        const data = await response.json();

        let html = '';

        // CPU Details
        if (data.cpu) {
            const perCoreUsage = data.cpu.per_core_usage && data.cpu.per_core_usage.length > 0
                ? data.cpu.per_core_usage.map((usage, idx) => `
                    <div class="col-md-3 mb-2">
                        <small>Core ${idx}: ${usage.toFixed(1)}%</small>
                        <div class="progress" style="height: 10px;">
                            <div class="progress-bar" style="width: ${usage}%"></div>
                        </div>
                    </div>
                `).join('')
                : '<div class="text-muted">No per-core data</div>';

            html += `
                <div class="card mb-3">
                    <div class="card-header"><strong>CPU Details</strong></div>
                    <div class="card-body">
                        <p><strong>Cores:</strong> ${data.cpu.core_count || data.cpu.cores || '-'}</p>
                        <p><strong>Model:</strong> ${data.cpu.model_name || data.cpu.model || '-'}</p>
                        <p><strong>MHz:</strong> ${data.cpu.mhz || '-'}</p>
                        <h6 class="mt-3">Per-Core Usage:</h6>
                        <div class="row">${perCoreUsage}</div>
                    </div>
                </div>
            `;
        }

        // Memory Details
        if (data.memory) {
            html += `
                <div class="card mb-3">
                    <div class="card-header"><strong>Memory Details</strong></div>
                    <div class="card-body">
                        <p><strong>Total:</strong> ${formatBytes(data.memory.total || 0)}</p>
                        <p><strong>Used:</strong> ${formatBytes(data.memory.used || 0)}</p>
                        <p><strong>Available:</strong> ${formatBytes(data.memory.available || 0)}</p>
                        <p><strong>Swap Total:</strong> ${formatBytes(data.memory.swap_total || 0)}</p>
                        <p><strong>Swap Used:</strong> ${formatBytes(data.memory.swap_used || 0)}</p>
                        <p><strong>Cache:</strong> ${formatBytes(data.memory.cached || 0)}</p>
                        <p><strong>Buffers:</strong> ${formatBytes(data.memory.buffers || 0)}</p>
                    </div>
                </div>
            `;
        }

        // Disk Details
        if (data.disk && data.disk.partitions) {
            const partitionsHtml = data.disk.partitions.map(part => `
                <div class="col-md-6 mb-3">
                    <div class="card">
                        <div class="card-body">
                            <h6>${part.mountpoint || part.device}</h6>
                            <div class="progress mb-2" style="height: 20px;">
                                <div class="progress-bar" style="width: ${part.percent || 0}%">
                                    ${(part.percent || 0).toFixed(1)}%
                                </div>
                            </div>
                            <small><strong>Total:</strong> ${formatBytes(part.total || 0)}</small><br>
                            <small><strong>Used:</strong> ${formatBytes(part.used || 0)}</small><br>
                            <small><strong>Free:</strong> ${formatBytes(part.free || 0)}</small>
                        </div>
                    </div>
                </div>
            `).join('');

            html += `
                <div class="card mb-3">
                    <div class="card-header"><strong>Disk Partitions</strong></div>
                    <div class="card-body">
                        <div class="row">${partitionsHtml}</div>
                    </div>
                </div>
            `;
        }

        // Network Details
        if (data.network && data.network.interfaces) {
            const interfacesHtml = data.network.interfaces.map(iface => {
                const ipAddresses = iface.ip_addresses && iface.ip_addresses.length > 0
                    ? iface.ip_addresses.join(', ')
                    : '-';

                const statusBadge = iface.is_up
                    ? '<span class="badge bg-success ms-2">UP</span>'
                    : '<span class="badge bg-secondary ms-2">DOWN</span>';

                return `
                    <div class="col-md-6 mb-3">
                        <div class="card">
                            <div class="card-body">
                                <h6>${iface.name || '-'}${statusBadge}</h6>
                                <p class="mb-1"><strong>IP:</strong> ${ipAddresses}</p>
                                <p class="mb-1"><strong>MAC:</strong> ${iface.mac_address || '-'}</p>
                                <p class="mb-1"><strong>Sent:</strong> ${formatBytes(iface.bytes_sent || 0)}</p>
                                <p class="mb-0"><strong>Received:</strong> ${formatBytes(iface.bytes_recv || 0)}</p>
                            </div>
                        </div>
                    </div>
                `;
            }).join('');

            html += `
                <div class="card mb-3">
                    <div class="card-header"><strong>Network Interfaces</strong></div>
                    <div class="card-body">
                        <div class="row">${interfacesHtml}</div>
                    </div>
                </div>
            `;
        }

        container.innerHTML = html || '<div class="text-muted text-center py-4">No detailed metrics available</div>';
    } catch (error) {
        console.error('Failed to load detailed metrics:', error);
        container.innerHTML = '<div class="alert alert-warning m-3">Failed to load detailed metrics</div>';
    }
}

// Refresh helper functions for HTML onclick handlers
function refreshProcesses() {
    loadProcesses();
}

function refreshConnections() {
    loadConnections();
}

function refreshLogs() {
    loadLogs();
}

function refreshErrors() {
    loadErrors();
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
    if (performanceChart) {
        performanceChart.destroy();
    }
});

// Populate CPU Details section in Overview tab
function populateCPUDetails(statsData) {
    const cpuDetails = document.getElementById('cpu-details');
    if (!cpuDetails) return;

    const cpuPercent = statsData.cpu_percent || 0;
    const loadAvg = statsData.load_avg || [0, 0, 0];
    const processCount = statsData.process_count || 0;

    cpuDetails.innerHTML = `
        <div class="row">
            <div class="col-6">
                <small class="text-muted">Current Usage</small>
                <h4>${cpuPercent.toFixed(1)}%</h4>
            </div>
            <div class="col-6">
                <small class="text-muted">Processes</small>
                <h4>${processCount}</h4>
            </div>
        </div>
        <div class="row mt-2">
            <div class="col-12">
                <small class="text-muted">Load Average</small>
                <p class="mb-0">${loadAvg[0].toFixed(2)} (1m) | ${loadAvg[1].toFixed(2)} (5m) | ${loadAvg[2].toFixed(2)} (15m)</p>
            </div>
        </div>
    `;
}

// Populate Memory Details section in Overview tab
function populateMemoryDetails(statsData) {
    const memoryDetails = document.getElementById('memory-details');
    if (!memoryDetails) return;

    const memUsedBytes = statsData.memory_used_bytes || 0;
    const memTotalBytes = statsData.memory_total_bytes || 0;
    const memPercent = statsData.memory_percent || 0;

    const usedGB = (memUsedBytes / (1024 * 1024 * 1024)).toFixed(2);
    const totalGB = (memTotalBytes / (1024 * 1024 * 1024)).toFixed(2);
    const freeGB = ((memTotalBytes - memUsedBytes) / (1024 * 1024 * 1024)).toFixed(2);

    memoryDetails.innerHTML = `
        <div class="row">
            <div class="col-6">
                <small class="text-muted">Used</small>
                <h4>${usedGB} GB</h4>
            </div>
            <div class="col-6">
                <small class="text-muted">Free</small>
                <h4>${freeGB} GB</h4>
            </div>
        </div>
        <div class="row mt-2">
            <div class="col-12">
                <small class="text-muted">Total Memory</small>
                <p class="mb-0">${totalGB} GB (${memPercent.toFixed(1)}% used)</p>
            </div>
        </div>
    `;
}

// ========== EVENTS TAB ==========
async function loadAgentEvents() {
    const currentAgentName = agentName || getAgentFromURL();
    if (!currentAgentName) return;

    const typeFilter = document.getElementById('event-type-filter')?.value || '';
    const statusFilter = document.getElementById('event-status-filter')?.value || '';
    const limit = document.getElementById('event-limit')?.value || '100';

    try {
        const params = new URLSearchParams({
            agent: currentAgentName,
            limit: limit
        });

        if (typeFilter) params.append('type', typeFilter);
        if (statusFilter) params.append('status', statusFilter);

        const response = await fetch(`/api/v1/events/by-agent?${params}`);
        const data = await response.json();

        // Update statistics
        updateEventStats(data.stats || {});

        // Render events
        renderEvents(data.events || []);
    } catch (error) {
        console.error('Failed to load events:', error);
        document.getElementById('events-content').innerHTML = `
            <div class="alert alert-danger">
                <i class="bi bi-exclamation-triangle"></i>
                Failed to load events: ${error.message}
            </div>
        `;
    }
}

function updateEventStats(stats) {
    document.getElementById('events-total').textContent = stats.total || 0;
    document.getElementById('events-pending').textContent = stats.pending || 0;
    document.getElementById('events-completed').textContent = stats.completed || 0;
    document.getElementById('events-failed').textContent = stats.failed || 0;

    // Update badge
    const badge = document.getElementById('events-badge');
    if (stats.total > 0) {
        badge.textContent = stats.total;
        badge.style.display = 'inline';
    } else {
        badge.style.display = 'none';
    }
}

function renderEvents(events) {
    const container = document.getElementById('events-content');

    if (!events || events.length === 0) {
        container.innerHTML = `
            <div class="text-center py-4 text-muted">
                <i class="bi bi-bell-slash fs-1"></i>
                <p class="mt-2">No events found for this agent</p>
            </div>
        `;
        return;
    }

    const html = `
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Timestamp</th>
                        <th>Type</th>
                        <th>Status</th>
                        <th>Stack</th>
                        <th>Run ID</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${events.map(event => renderEventRow(event)).join('')}
                </tbody>
            </table>
        </div>
    `;

    container.innerHTML = html;
}

function renderEventRow(event) {
    const status = event.status || 'unknown';
    const statusBadge = getStatusBadge(status);
    const typeBadge = getEventTypeBadge(event.type);
    const timestamp = event.created_at ? new Date(event.created_at).toLocaleString() : '-';

    return `
        <tr>
            <td><small>${timestamp}</small></td>
            <td>${typeBadge}</td>
            <td>${statusBadge}</td>
            <td><code class="small">${event.stack || '-'}</code></td>
            <td><small class="text-muted">${event.run_id ? event.run_id.substring(0, 8) : '-'}</small></td>
            <td>
                <button class="btn btn-sm btn-outline-primary" onclick="viewEventDetails('${event.id}')" title="View Details">
                    <i class="bi bi-eye"></i>
                </button>
                ${status === 'failed' ? `
                    <button class="btn btn-sm btn-outline-warning" onclick="retryEvent('${event.id}')" title="Retry">
                        <i class="bi bi-arrow-clockwise"></i>
                    </button>
                ` : ''}
            </td>
        </tr>
    `;
}

function getStatusBadge(status) {
    const badges = {
        'pending': '<span class="badge bg-secondary">Pending</span>',
        'processing': '<span class="badge bg-info">Processing</span>',
        'completed': '<span class="badge bg-success">Completed</span>',
        'failed': '<span class="badge bg-danger">Failed</span>'
    };
    return badges[status] || `<span class="badge bg-secondary">${status}</span>`;
}

function getEventTypeBadge(type) {
    const icons = {
        'task.started': '<i class="bi bi-play-circle text-primary"></i>',
        'task.completed': '<i class="bi bi-check-circle text-success"></i>',
        'task.failed': '<i class="bi bi-x-circle text-danger"></i>',
        'file.created': '<i class="bi bi-file-plus text-success"></i>',
        'file.changed': '<i class="bi bi-file-text text-warning"></i>',
        'file.deleted': '<i class="bi bi-file-minus text-danger"></i>',
        'process.started': '<i class="bi bi-play-fill text-info"></i>',
        'process.stopped': '<i class="bi bi-stop-fill text-warning"></i>'
    };

    const icon = icons[type] || '<i class="bi bi-bell"></i>';
    return `${icon} <span class="ms-1">${type}</span>`;
}

async function viewEventDetails(eventId) {
    try {
        const response = await fetch(`/api/v1/events/${eventId}`);
        const event = await response.json();

        // Show modal with event details
        const modalHtml = `
            <div class="modal fade" id="eventDetailsModal" tabindex="-1">
                <div class="modal-dialog modal-lg">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title">Event Details</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                        </div>
                        <div class="modal-body">
                            <dl class="row">
                                <dt class="col-sm-3">ID</dt>
                                <dd class="col-sm-9"><code>${event.id}</code></dd>

                                <dt class="col-sm-3">Type</dt>
                                <dd class="col-sm-9">${getEventTypeBadge(event.type)}</dd>

                                <dt class="col-sm-3">Status</dt>
                                <dd class="col-sm-9">${getStatusBadge(event.status)}</dd>

                                <dt class="col-sm-3">Agent</dt>
                                <dd class="col-sm-9">${event.agent || '-'}</dd>

                                <dt class="col-sm-3">Stack</dt>
                                <dd class="col-sm-9"><code>${event.stack || '-'}</code></dd>

                                <dt class="col-sm-3">Run ID</dt>
                                <dd class="col-sm-9"><code>${event.run_id || '-'}</code></dd>

                                <dt class="col-sm-3">Created At</dt>
                                <dd class="col-sm-9">${new Date(event.created_at).toLocaleString()}</dd>

                                ${event.processed_at ? `
                                    <dt class="col-sm-3">Processed At</dt>
                                    <dd class="col-sm-9">${new Date(event.processed_at).toLocaleString()}</dd>
                                ` : ''}

                                ${event.error ? `
                                    <dt class="col-sm-3">Error</dt>
                                    <dd class="col-sm-9"><span class="text-danger">${event.error}</span></dd>
                                ` : ''}

                                <dt class="col-sm-3">Data</dt>
                                <dd class="col-sm-9"><pre class="bg-light p-2 rounded">${JSON.stringify(event.data, null, 2)}</pre></dd>
                            </dl>
                        </div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                        </div>
                    </div>
                </div>
            </div>
        `;

        // Remove existing modal if any
        const existingModal = document.getElementById('eventDetailsModal');
        if (existingModal) existingModal.remove();

        // Add modal to DOM
        document.body.insertAdjacentHTML('beforeend', modalHtml);

        // Show modal
        const modal = new bootstrap.Modal(document.getElementById('eventDetailsModal'));
        modal.show();
    } catch (error) {
        console.error('Failed to load event details:', error);
        alert('Failed to load event details');
    }
}

async function retryEvent(eventId) {
    if (!confirm('Are you sure you want to retry this event?')) return;

    try {
        const response = await fetch(`/api/v1/events/${eventId}/retry`, {
            method: 'POST'
        });

        if (response.ok) {
            showToast('Event queued for retry', 'success');
            await loadAgentEvents();
        } else {
            throw new Error('Failed to retry event');
        }
    } catch (error) {
        console.error('Failed to retry event:', error);
        showToast('Failed to retry event', 'error');
    }
}

function refreshEvents() {
    loadAgentEvents();
}

// =============================================================================
// HOOKS TAB
// =============================================================================

async function loadAgentHooks() {
    const currentAgentName = agentName || getAgentFromURL();
    if (!currentAgentName) return;

    try {
        const limit = document.getElementById('hook-limit')?.value || 100;
        const params = new URLSearchParams({
            agent: currentAgentName,
            limit: limit
        });

        const response = await fetch(`/api/v1/events/hook-executions/by-agent?${params}`);
        const data = await response.json();

        // Update statistics
        updateHookStats(data.stats || {});

        // Render hooks
        renderHooks(data.executions || []);

        // Update badge
        const badge = document.getElementById('hooks-badge');
        if (badge && data.stats && data.stats.total > 0) {
            badge.textContent = data.stats.total;
            badge.style.display = 'inline';
        }
    } catch (error) {
        console.error('Failed to load hooks:', error);
        document.getElementById('hooks-content').innerHTML = `
            <div class="alert alert-danger">
                <i class="bi bi-exclamation-triangle"></i> Failed to load hooks: ${error.message}
            </div>
        `;
    }
}

function updateHookStats(stats) {
    document.getElementById('hooks-total').textContent = stats.total || 0;
    document.getElementById('hooks-success').textContent = stats.success || 0;
    document.getElementById('hooks-failed').textContent = stats.failed || 0;
}

function renderHooks(hooks) {
    const container = document.getElementById('hooks-content');

    if (hooks.length === 0) {
        container.innerHTML = `
            <div class="text-center py-4 text-muted">
                <i class="bi bi-lightning-charge fs-1"></i>
                <p class="mt-2">No hook executions found for this agent</p>
            </div>
        `;
        return;
    }

    let html = '<div class="table-responsive"><table class="table table-sm table-hover">';
    html += `
        <thead>
            <tr>
                <th>Time</th>
                <th>Hook Name</th>
                <th>Event Type</th>
                <th>Status</th>
                <th>Duration</th>
                <th>Stack</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
    `;

    hooks.forEach(hook => {
        const status = hook.success ?
            '<span class="badge bg-success">Success</span>' :
            '<span class="badge bg-danger">Failed</span>';

        const durationMs = hook.duration / 1000000; // Convert nanoseconds to milliseconds
        const duration = durationMs < 1000 ?
            `${durationMs.toFixed(2)}ms` :
            `${(durationMs / 1000).toFixed(2)}s`;

        const time = new Date(hook.executed_at).toLocaleString();

        html += `
            <tr>
                <td><small class="text-muted">${time}</small></td>
                <td>
                    <strong>${hook.hook_name}</strong><br>
                    <small class="text-muted">${hook.hook_id.substring(0, 8)}</small>
                </td>
                <td>
                    <span class="badge bg-info">${hook.event_type || 'N/A'}</span><br>
                    <small class="text-muted">${hook.event_id.substring(0, 8)}</small>
                </td>
                <td>${status}</td>
                <td><code>${duration}</code></td>
                <td><small class="text-muted">${hook.event_stack || '-'}</small></td>
                <td>
                    <button class="btn btn-sm btn-outline-primary" onclick="viewHookDetails('${hook.id}', '${hook.event_id}')">
                        <i class="bi bi-eye"></i>
                    </button>
                </td>
            </tr>
        `;
    });

    html += '</tbody></table></div>';
    container.innerHTML = html;
}

async function viewHookDetails(hookExecId, eventId) {
    try {
        // Fetch the event details to get more context
        const response = await fetch(`/api/v1/events/${eventId}`);
        const event = await response.json();

        // Show modal with details
        const modalHtml = `
            <div class="modal fade" id="hookDetailsModal" tabindex="-1">
                <div class="modal-dialog modal-lg">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title">
                                <i class="bi bi-lightning-charge"></i> Hook Execution Details
                            </h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                        </div>
                        <div class="modal-body">
                            <h6>Event Information</h6>
                            <dl class="row">
                                <dt class="col-sm-3">Event ID</dt>
                                <dd class="col-sm-9"><code>${eventId}</code></dd>

                                <dt class="col-sm-3">Event Type</dt>
                                <dd class="col-sm-9"><span class="badge bg-info">${event.type}</span></dd>

                                <dt class="col-sm-3">Timestamp</dt>
                                <dd class="col-sm-9">${new Date(event.timestamp).toLocaleString()}</dd>

                                <dt class="col-sm-3">Stack</dt>
                                <dd class="col-sm-9">${event.stack || '-'}</dd>

                                <dt class="col-sm-3">Run ID</dt>
                                <dd class="col-sm-9"><code>${event.run_id || '-'}</code></dd>
                            </dl>

                            <hr>

                            <h6>Event Data</h6>
                            <pre class="bg-light p-3 rounded"><code>${JSON.stringify(event.data, null, 2)}</code></pre>
                        </div>
                    </div>
                </div>
            </div>
        `;

        // Remove any existing modal
        document.getElementById('hookDetailsModal')?.remove();

        // Add new modal
        document.body.insertAdjacentHTML('beforeend', modalHtml);

        // Show modal
        const modal = new bootstrap.Modal(document.getElementById('hookDetailsModal'));
        modal.show();
    } catch (error) {
        console.error('Failed to load hook details:', error);
        showToast('Failed to load hook details', 'error');
    }
}

function refreshHooks() {
    loadAgentHooks();
}

// Add event listeners when hooks tab is shown
document.getElementById('hooks-tab')?.addEventListener('shown.bs.tab', function() {
    loadAgentHooks();
});

// Add event listeners for filters
document.getElementById('hook-limit')?.addEventListener('change', loadAgentHooks);

// Add event listeners when events tab is shown
document.getElementById('events-tab')?.addEventListener('shown.bs.tab', function() {
    loadAgentEvents();
});

// Add event listeners for filters
document.getElementById('event-type-filter')?.addEventListener('change', loadAgentEvents);
document.getElementById('event-status-filter')?.addEventListener('change', loadAgentEvents);
document.getElementById('event-limit')?.addEventListener('change', loadAgentEvents);

// ==================== NETWORK METRICS ====================

let networkTrafficChart = null;
let networkHistory = [];
const MAX_NETWORK_HISTORY = 60; // 5 minutes at 5s intervals
let lastNetworkBytes = null;

// Load network metrics
async function loadNetworkMetrics() {
    if (!agentName) return;

    try {
        const response = await fetch(`/api/v1/network/agent/${agentName}`);
        if (!response.ok) throw new Error('Failed to fetch network metrics');

        const data = await response.json();

        // Update summary cards
        document.getElementById('net-total-rx').textContent = formatBytes(data.total_rx_bytes || 0);
        document.getElementById('net-total-tx').textContent = formatBytes(data.total_tx_bytes || 0);

        // Calculate bandwidth (MB/s) if we have previous data
        if (lastNetworkBytes) {
            const timeDiff = 5; // 5 seconds between updates
            const rxDiff = data.total_rx_bytes - lastNetworkBytes.rx;
            const txDiff = data.total_tx_bytes - lastNetworkBytes.tx;

            const rxBandwidth = (rxDiff / timeDiff) / 1024 / 1024; // MB/s
            const txBandwidth = (txDiff / timeDiff) / 1024 / 1024; // MB/s

            document.getElementById('net-bandwidth-rx').textContent = rxBandwidth.toFixed(2) + ' MB/s';
            document.getElementById('net-bandwidth-tx').textContent = txBandwidth.toFixed(2) + ' MB/s';

            // Add to history for chart
            networkHistory.push({
                timestamp: new Date(),
                rx: rxBandwidth,
                tx: txBandwidth
            });

            // Keep only recent history
            if (networkHistory.length > MAX_NETWORK_HISTORY) {
                networkHistory.shift();
            }

            // Update chart
            updateNetworkTrafficChart();
        } else {
            document.getElementById('net-bandwidth-rx').textContent = '0.00 MB/s';
            document.getElementById('net-bandwidth-tx').textContent = '0.00 MB/s';
        }

        // Store current values for next calculation
        lastNetworkBytes = {
            rx: data.total_rx_bytes,
            tx: data.total_tx_bytes
        };

        // Update interfaces table
        updateNetworkInterfacesTable(data.interfaces || []);

    } catch (error) {
        console.error('Error loading network metrics:', error);
    }
}

// Update network traffic chart
function updateNetworkTrafficChart() {
    const canvas = document.getElementById('network-traffic-chart');
    if (!canvas) return;

    const ctx = canvas.getContext('2d');

    const labels = networkHistory.map(h => h.timestamp.toLocaleTimeString());
    const rxData = networkHistory.map(h => h.rx);
    const txData = networkHistory.map(h => h.tx);

    if (networkTrafficChart) {
        networkTrafficChart.destroy();
    }

    networkTrafficChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Download (MB/s)',
                    data: rxData,
                    borderColor: 'rgba(75, 192, 192, 1)',
                    backgroundColor: 'rgba(75, 192, 192, 0.1)',
                    fill: true,
                    tension: 0.4,
                    borderWidth: 2
                },
                {
                    label: 'Upload (MB/s)',
                    data: txData,
                    borderColor: 'rgba(255, 99, 132, 1)',
                    backgroundColor: 'rgba(255, 99, 132, 0.1)',
                    fill: true,
                    tension: 0.4,
                    borderWidth: 2
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
                },
                tooltip: {
                    mode: 'index',
                    intersect: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value.toFixed(2) + ' MB/s';
                        }
                    }
                },
                x: {
                    ticks: {
                        maxTicksLimit: 10
                    }
                }
            },
            interaction: {
                mode: 'nearest',
                axis: 'x',
                intersect: false
            }
        }
    });
}

// Update network interfaces table
function updateNetworkInterfacesTable(interfaces) {
    const tbody = document.getElementById('network-interfaces-table');
    if (!tbody) return;

    if (interfaces.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" class="text-center text-muted">No network interfaces found</td></tr>';
        return;
    }

    tbody.innerHTML = interfaces.map(iface => `
        <tr>
            <td><strong>${escapeHtml(iface.name)}</strong></td>
            <td>
                ${iface.is_up ?
                    '<span class="badge bg-success">UP</span>' :
                    '<span class="badge bg-secondary">DOWN</span>'}
            </td>
            <td>
                ${iface.ip_addresses && iface.ip_addresses.length > 0 ?
                    iface.ip_addresses.map(ip => `<span class="badge bg-info me-1">${escapeHtml(ip)}</span>`).join('') :
                    '<span class="text-muted">-</span>'}
            </td>
            <td><small class="font-monospace">${escapeHtml(iface.mac_address || '-')}</small></td>
            <td><span class="text-success">${formatBytes(iface.bytes_recv || 0)}</span></td>
            <td><span class="text-primary">${formatBytes(iface.bytes_sent || 0)}</span></td>
            <td>${(iface.packets_recv || 0).toLocaleString()}</td>
            <td>${(iface.packets_sent || 0).toLocaleString()}</td>
        </tr>
    `).join('');
}

// Add network tab listener
document.getElementById('network-tab')?.addEventListener('shown.bs.tab', function() {
    console.log('Network tab shown, loading metrics...');
    loadNetworkMetrics();

    // Start auto-refresh for network tab
    if (window.networkRefreshInterval) {
        clearInterval(window.networkRefreshInterval);
    }
    window.networkRefreshInterval = setInterval(loadNetworkMetrics, 5000);
});

// Stop network refresh when tab is hidden
document.getElementById('network-tab')?.addEventListener('hidden.bs.tab', function() {
    if (window.networkRefreshInterval) {
        clearInterval(window.networkRefreshInterval);
        window.networkRefreshInterval = null;
    }
});

