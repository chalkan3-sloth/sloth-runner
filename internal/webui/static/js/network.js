// Network Dashboard JavaScript

let trafficChart = null;
let distributionChart = null;
let bandwidthChart = null;
let refreshInterval = null;

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    loadNetworkData();

    // Auto-refresh every 30 seconds
    refreshInterval = setInterval(loadNetworkData, 30000);
});

// Cleanup on page unload
window.addEventListener('beforeunload', function() {
    if (refreshInterval) {
        clearInterval(refreshInterval);
    }
});

// Load all network data
async function loadNetworkData() {
    try {
        await Promise.all([
            loadNetworkSummary(),
            loadAllAgentsNetwork(),
            loadTopAgents()
        ]);
    } catch (error) {
        console.error('Error loading network data:', error);
        showToast('Failed to load network data', 'error');
    }
}

// Load network summary
async function loadNetworkSummary() {
    try {
        const response = await fetch('/api/v1/network/summary');
        if (!response.ok) throw new Error('Failed to fetch network summary');

        const data = await response.json();

        // Update summary cards
        document.getElementById('totalAgents').textContent = data.total_agents || 0;
        document.getElementById('activeInterfaces').textContent = data.active_interfaces || 0;
        document.getElementById('totalDownload').textContent = formatBytes(data.total_rx_bytes || 0);
        document.getElementById('totalUpload').textContent = formatBytes(data.total_tx_bytes || 0);

    } catch (error) {
        console.error('Error loading network summary:', error);
    }
}

// Load all agents network stats
async function loadAllAgentsNetwork() {
    try {
        const response = await fetch('/api/v1/network/all');
        if (!response.ok) throw new Error('Failed to fetch agents network stats');

        const data = await response.json();

        // Update table
        updateAgentsTable(data.agents || []);

        // Update charts
        updateTrafficChart(data.agents || []);
        updateDistributionChart(data.agents || []);
        updateBandwidthChart(data.agents || []);

    } catch (error) {
        console.error('Error loading agents network stats:', error);
        document.getElementById('agentsTableBody').innerHTML = `
            <tr><td colspan="6" class="text-center text-danger">
                Failed to load network data
            </td></tr>
        `;
    }
}

// Update agents table
function updateAgentsTable(agents) {
    const tbody = document.getElementById('agentsTableBody');

    if (agents.length === 0) {
        tbody.innerHTML = `
            <tr><td colspan="6" class="text-center">No agents found</td></tr>
        `;
        return;
    }

    tbody.innerHTML = agents.map(agent => `
        <tr>
            <td>
                <strong>${escapeHtml(agent.agent_name)}</strong>
            </td>
            <td>
                <span class="badge bg-primary">${agent.interfaces.length} interfaces</span>
            </td>
            <td>
                <i class="bi bi-arrow-down-circle text-success"></i>
                ${formatBytes(agent.total_rx_bytes)}
            </td>
            <td>
                <i class="bi bi-arrow-up-circle text-info"></i>
                ${formatBytes(agent.total_tx_bytes)}
            </td>
            <td>
                <strong>${formatBytes(agent.total_rx_bytes + agent.total_tx_bytes)}</strong>
            </td>
            <td>
                <button class="btn btn-sm btn-outline-primary"
                        onclick="showAgentDetails('${escapeHtml(agent.agent_name)}')">
                    <i class="bi bi-eye"></i> Details
                </button>
            </td>
        </tr>
    `).join('');
}

// Update traffic chart
function updateTrafficChart(agents) {
    const ctx = document.getElementById('trafficChart');
    if (!ctx) return;

    const labels = agents.map(a => a.agent_name);
    const rxData = agents.map(a => a.total_rx_mb);
    const txData = agents.map(a => a.total_tx_mb);

    if (trafficChart) {
        trafficChart.destroy();
    }

    trafficChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Download (MB)',
                    data: rxData,
                    backgroundColor: 'rgba(75, 192, 192, 0.6)',
                    borderColor: 'rgba(75, 192, 192, 1)',
                    borderWidth: 1
                },
                {
                    label: 'Upload (MB)',
                    data: txData,
                    backgroundColor: 'rgba(255, 99, 132, 0.6)',
                    borderColor: 'rgba(255, 99, 132, 1)',
                    borderWidth: 1
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
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return context.dataset.label + ': ' + context.parsed.y.toFixed(2) + ' MB';
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value.toFixed(0) + ' MB';
                        }
                    }
                }
            }
        }
    });
}

// Update distribution chart
function updateDistributionChart(agents) {
    const ctx = document.getElementById('distributionChart');
    if (!ctx) return;

    const labels = agents.map(a => a.agent_name);
    const data = agents.map(a => a.total_rx_bytes + a.total_tx_bytes);

    const colors = [
        'rgba(255, 99, 132, 0.6)',
        'rgba(54, 162, 235, 0.6)',
        'rgba(255, 206, 86, 0.6)',
        'rgba(75, 192, 192, 0.6)',
        'rgba(153, 102, 255, 0.6)',
        'rgba(255, 159, 64, 0.6)',
        'rgba(199, 199, 199, 0.6)',
        'rgba(83, 102, 255, 0.6)',
        'rgba(255, 102, 196, 0.6)',
        'rgba(102, 255, 178, 0.6)'
    ];

    if (distributionChart) {
        distributionChart.destroy();
    }

    distributionChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: colors,
                borderColor: colors.map(c => c.replace('0.6', '1')),
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'right',
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = ((context.parsed / total) * 100).toFixed(1);
                            return context.label + ': ' + formatBytes(context.parsed) + ' (' + percentage + '%)';
                        }
                    }
                }
            }
        }
    });
}

// Update bandwidth chart
function updateBandwidthChart(agents) {
    const ctx = document.getElementById('bandwidthChart');
    if (!ctx) return;

    const labels = agents.map(a => a.agent_name);
    const totalData = agents.map(a => (a.total_rx_bytes + a.total_tx_bytes) / 1024 / 1024);

    if (bandwidthChart) {
        bandwidthChart.destroy();
    }

    bandwidthChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Total Traffic (MB)',
                data: totalData,
                borderColor: 'rgba(102, 126, 234, 1)',
                backgroundColor: 'rgba(102, 126, 234, 0.1)',
                fill: true,
                tension: 0.4,
                borderWidth: 2
            }]
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
                    callbacks: {
                        label: function(context) {
                            return 'Traffic: ' + context.parsed.y.toFixed(2) + ' MB';
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value.toFixed(0) + ' MB';
                        }
                    }
                }
            }
        }
    });
}

// Load top agents
async function loadTopAgents() {
    try {
        const response = await fetch('/api/v1/network/top');
        if (!response.ok) throw new Error('Failed to fetch top agents');

        const data = await response.json();
        const topAgents = data.top_agents || [];

        const container = document.getElementById('topAgentsList');

        if (topAgents.length === 0) {
            container.innerHTML = '<p class="text-muted text-center">No data available</p>';
            return;
        }

        container.innerHTML = topAgents.map((agent, index) => `
            <div class="agent-network-item">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <h6 class="mb-0">
                            <span class="badge bg-primary me-2">#${index + 1}</span>
                            ${escapeHtml(agent.agent_name)}
                        </h6>
                    </div>
                    <div class="text-end">
                        <div class="fw-bold">${formatBytes(agent.total_bytes)}</div>
                        <small class="text-muted">
                            <i class="bi bi-arrow-down"></i> ${formatBytes(agent.rx_bytes)} /
                            <i class="bi bi-arrow-up"></i> ${formatBytes(agent.tx_bytes)}
                        </small>
                    </div>
                </div>
                <div class="progress mt-2" style="height: 6px;">
                    <div class="progress-bar bg-gradient"
                         style="width: ${(agent.total_bytes / topAgents[0].total_bytes * 100).toFixed(1)}%"></div>
                </div>
            </div>
        `).join('');

    } catch (error) {
        console.error('Error loading top agents:', error);
        document.getElementById('topAgentsList').innerHTML =
            '<p class="text-danger text-center">Failed to load top agents</p>';
    }
}

// Show agent details modal
async function showAgentDetails(agentName) {
    try {
        const response = await fetch(`/api/v1/network/agent/${agentName}`);
        if (!response.ok) throw new Error('Failed to fetch agent details');

        const data = await response.json();

        document.getElementById('modalAgentName').textContent = agentName;

        const interfacesList = document.getElementById('agentInterfacesList');
        interfacesList.innerHTML = data.interfaces.map(iface => `
            <div class="card mb-3">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <div>
                            <h5 class="card-title mb-1">
                                <i class="bi bi-ethernet"></i> ${escapeHtml(iface.name)}
                                ${iface.is_up ?
                                    '<span class="interface-badge interface-up">UP</span>' :
                                    '<span class="interface-badge interface-down">DOWN</span>'}
                            </h5>
                            <p class="text-muted mb-2">
                                <i class="bi bi-hdd-network"></i> ${escapeHtml(iface.mac_address)}
                            </p>
                        </div>
                    </div>

                    ${iface.ip_addresses && iface.ip_addresses.length > 0 ? `
                        <div class="mb-2">
                            <strong>IP Addresses:</strong>
                            ${iface.ip_addresses.map(ip => `
                                <span class="badge bg-secondary me-1">${escapeHtml(ip)}</span>
                            `).join('')}
                        </div>
                    ` : ''}

                    <div class="row">
                        <div class="col-md-6">
                            <div class="mb-2">
                                <i class="bi bi-arrow-down-circle text-success"></i>
                                <strong>Download:</strong> ${formatBytes(iface.bytes_recv)}
                            </div>
                            <div>
                                <i class="bi bi-box-arrow-in-down"></i>
                                <strong>Packets RX:</strong> ${iface.packets_recv.toLocaleString()}
                            </div>
                        </div>
                        <div class="col-md-6">
                            <div class="mb-2">
                                <i class="bi bi-arrow-up-circle text-info"></i>
                                <strong>Upload:</strong> ${formatBytes(iface.bytes_sent)}
                            </div>
                            <div>
                                <i class="bi bi-box-arrow-up"></i>
                                <strong>Packets TX:</strong> ${iface.packets_sent.toLocaleString()}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');

        const modal = new bootstrap.Modal(document.getElementById('agentDetailsModal'));
        modal.show();

    } catch (error) {
        console.error('Error loading agent details:', error);
        showToast('Failed to load agent details', 'error');
    }
}

// Utility functions
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
}

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

function showToast(message, type = 'info') {
    // Simple toast implementation (you can use a library like Bootstrap Toast)
    console.log(`[${type.toUpperCase()}] ${message}`);

    // If Bootstrap toast is available
    if (typeof bootstrap !== 'undefined') {
        const toastHtml = `
            <div class="toast align-items-center text-white bg-${type === 'error' ? 'danger' : type === 'success' ? 'success' : 'primary'} border-0" role="alert">
                <div class="d-flex">
                    <div class="toast-body">${escapeHtml(message)}</div>
                    <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
                </div>
            </div>
        `;

        let toastContainer = document.querySelector('.toast-container');
        if (!toastContainer) {
            toastContainer = document.createElement('div');
            toastContainer.className = 'toast-container position-fixed top-0 end-0 p-3';
            document.body.appendChild(toastContainer);
        }

        toastContainer.insertAdjacentHTML('beforeend', toastHtml);
        const toastElement = toastContainer.lastElementChild;
        const toast = new bootstrap.Toast(toastElement);
        toast.show();

        toastElement.addEventListener('hidden.bs.toast', () => {
            toastElement.remove();
        });
    }
}
