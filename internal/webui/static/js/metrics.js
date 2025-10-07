// System Metrics Page
let metricsChart = null;
let cpuChart = null;
let memoryChart = null;

const API = {
    async get(url) {
        const response = await fetch(url);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    }
};

// Initialize charts
function initCharts() {
    // Main metrics chart
    const ctx1 = document.getElementById('metricsChart');
    if (ctx1) {
        metricsChart = new Chart(ctx1, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: 'CPU %',
                    data: [],
                    borderColor: 'rgb(13, 110, 253)',
                    backgroundColor: 'rgba(13, 110, 253, 0.1)',
                    tension: 0.4
                }, {
                    label: 'Memory %',
                    data: [],
                    borderColor: 'rgb(25, 135, 84)',
                    backgroundColor: 'rgba(25, 135, 84, 0.1)',
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

    // CPU usage pie
    const ctx2 = document.getElementById('cpuChart');
    if (ctx2) {
        cpuChart = new Chart(ctx2, {
            type: 'doughnut',
            data: {
                labels: ['Used', 'Available'],
                datasets: [{
                    data: [0, 100],
                    backgroundColor: ['rgb(13, 110, 253)', 'rgb(233, 236, 239)']
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false
            }
        });
    }

    // Memory usage pie
    const ctx3 = document.getElementById('memoryChart');
    if (ctx3) {
        memoryChart = new Chart(ctx3, {
            type: 'doughnut',
            data: {
                labels: ['Used', 'Free'],
                datasets: [{
                    data: [0, 100],
                    backgroundColor: ['rgb(25, 135, 84)', 'rgb(233, 236, 239)']
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false
            }
        });
    }
}

// Load metrics
async function loadMetrics() {
    try {
        const data = await API.get('/api/v1/metrics');

        // Update stats cards
        if (data.cpu) {
            document.getElementById('cpu-usage').textContent = data.cpu.usage_percent.toFixed(1) + '%';
            if (cpuChart) {
                cpuChart.data.datasets[0].data = [data.cpu.usage_percent, 100 - data.cpu.usage_percent];
                cpuChart.update('none');
            }
        }

        if (data.memory) {
            const memUsedGB = (data.memory.used / 1024 / 1024 / 1024).toFixed(2);
            const memTotalGB = (data.memory.total / 1024 / 1024 / 1024).toFixed(2);
            document.getElementById('memory-usage').textContent = `${memUsedGB} / ${memTotalGB} GB`;

            if (memoryChart) {
                memoryChart.data.datasets[0].data = [data.memory.used_percent, 100 - data.memory.used_percent];
                memoryChart.update('none');
            }
        }

        if (data.disk) {
            const diskUsedGB = (data.disk.used / 1024 / 1024 / 1024).toFixed(2);
            const diskTotalGB = (data.disk.total / 1024 / 1024 / 1024).toFixed(2);
            document.getElementById('disk-usage').textContent = `${diskUsedGB} / ${diskTotalGB} GB`;
        }

        if (data.network) {
            const rxMB = (data.network.bytes_recv / 1024 / 1024).toFixed(2);
            const txMB = (data.network.bytes_sent / 1024 / 1024).toFixed(2);
            document.getElementById('network-usage').textContent = `↓${rxMB} MB  ↑${txMB} MB`;
        }

        // Update line chart
        if (metricsChart && data.cpu && data.memory) {
            const now = new Date().toLocaleTimeString();
            metricsChart.data.labels.push(now);
            metricsChart.data.datasets[0].data.push(data.cpu.usage_percent);
            metricsChart.data.datasets[1].data.push(data.memory.used_percent);

            if (metricsChart.data.labels.length > 20) {
                metricsChart.data.labels.shift();
                metricsChart.data.datasets[0].data.shift();
                metricsChart.data.datasets[1].data.shift();
            }

            metricsChart.update('none');
        }

        // Update system info
        if (data.host_info) {
            document.getElementById('hostname').textContent = data.host_info.hostname || 'N/A';
            document.getElementById('platform').textContent = `${data.host_info.os || 'Unknown'} ${data.host_info.platform || ''}`;
            document.getElementById('uptime').textContent = formatUptime(data.host_info.uptime || 0);
        }

        if (data.goroutines !== undefined) {
            document.getElementById('goroutines').textContent = data.goroutines;
        }

    } catch (error) {
        console.error('Failed to load metrics:', error);
    }
}

// Format uptime
function formatUptime(seconds) {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    if (days > 0) return `${days}d ${hours}h ${minutes}m`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    return `${minutes}m`;
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    initCharts();
    loadMetrics();

    // Refresh every 3 seconds
    setInterval(loadMetrics, 3000);
});
