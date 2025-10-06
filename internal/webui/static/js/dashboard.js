// Dashboard page logic

let activityFeedItems = [];
const MAX_ACTIVITY_ITEMS = 50;

// Load dashboard data
async function loadDashboardData() {
    try {
        // Load stats
        const statsResponse = await fetch('/api/v1/dashboard');
        const stats = await statsResponse.json();
        updateStats(stats);

        // Load recent events
        await loadRecentEvents();

        // Load active agents
        await loadActiveAgents();
    } catch (error) {
        console.error('Error loading dashboard data:', error);
    }
}

// Update stats cards
function updateStats(stats) {
    document.getElementById('stat-agents-active').textContent = stats.agents.active;
    document.getElementById('stat-agents-total').textContent = stats.agents.total;

    document.getElementById('stat-workflows-active').textContent = stats.workflows.active;
    document.getElementById('stat-workflows-total').textContent = stats.workflows.total;

    document.getElementById('stat-hooks-enabled').textContent = stats.hooks.enabled;
    document.getElementById('stat-hooks-total').textContent = stats.hooks.total;

    document.getElementById('stat-events-pending').textContent = stats.events.pending;
}

// Load recent events
async function loadRecentEvents() {
    try {
        const response = await fetch('/api/v1/events?limit=10');
        const data = await response.json();

        const tbody = document.getElementById('events-list');

        if (!data.events || data.events.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center text-muted">No recent events</td></tr>';
            return;
        }

        tbody.innerHTML = data.events.map(event => {
            const statusClass = getEventStatusClass(event.status);
            const createdDate = new Date(event.created_at * 1000).toLocaleString();
            const processedDate = event.processed_at ? new Date(event.processed_at * 1000).toLocaleString() : '-';

            return `
                <tr>
                    <td>${event.id}</td>
                    <td><code>${event.event_type}</code></td>
                    <td><span class="badge bg-${statusClass}">${event.status}</span></td>
                    <td>${event.hook_id || '-'}</td>
                    <td>${createdDate}</td>
                    <td>${processedDate}</td>
                </tr>
            `;
        }).join('');
    } catch (error) {
        console.error('Error loading events:', error);
    }
}

// Load active agents
async function loadActiveAgents() {
    try {
        const response = await fetch('/api/v1/agents');
        const data = await response.json();

        const container = document.getElementById('active-agents-list');

        if (!data.agents || data.agents.length === 0) {
            container.innerHTML = '<p class="text-muted text-center">No agents registered</p>';
            return;
        }

        const activeAgents = data.agents.filter(agent => agent.status === 'Active');

        if (activeAgents.length === 0) {
            container.innerHTML = '<p class="text-muted text-center">No active agents</p>';
            return;
        }

        container.innerHTML = activeAgents.map(agent => {
            const lastHeartbeat = new Date(agent.last_heartbeat * 1000).toLocaleString();

            return `
                <div class="d-flex align-items-center mb-3 p-2 border-start border-4 border-success">
                    <div class="flex-grow-1">
                        <h6 class="mb-1">${agent.name}</h6>
                        <small class="text-muted">${agent.address}</small>
                        ${agent.version ? `<br><small class="text-muted">Version: ${agent.version}</small>` : ''}
                    </div>
                    <div class="text-end">
                        <small class="text-muted d-block">Last heartbeat</small>
                        <small class="text-muted">${lastHeartbeat}</small>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Error loading agents:', error);
    }
}

// Add activity to feed
function addActivityItem(type, message, icon, color) {
    const timestamp = new Date().toLocaleTimeString();

    const item = {
        type,
        message,
        icon,
        color,
        timestamp: Date.now(),
        formattedTime: timestamp
    };

    activityFeedItems.unshift(item);

    // Keep only the latest items
    if (activityFeedItems.length > MAX_ACTIVITY_ITEMS) {
        activityFeedItems = activityFeedItems.slice(0, MAX_ACTIVITY_ITEMS);
    }

    updateActivityFeed();
}

// Update activity feed display
function updateActivityFeed() {
    const container = document.getElementById('activity-feed');

    if (activityFeedItems.length === 0) {
        container.innerHTML = '<p class="text-muted text-center">No recent activity</p>';
        return;
    }

    container.innerHTML = activityFeedItems.map(item => `
        <div class="activity-item d-flex align-items-start">
            <div class="activity-icon bg-${item.color} text-white">
                <i class="bi bi-${item.icon}"></i>
            </div>
            <div class="flex-grow-1">
                <p class="mb-1">${item.message}</p>
                <small class="activity-time">${item.formattedTime}</small>
            </div>
        </div>
    `).join('');
}

// Get event status CSS class
function getEventStatusClass(status) {
    switch (status) {
        case 'pending': return 'warning';
        case 'processing': return 'info';
        case 'completed': return 'success';
        case 'failed': return 'danger';
        default: return 'secondary';
    }
}

// WebSocket event handlers
wsManager.on('agent_update', (data) => {
    addActivityItem('agent', `Agent ${data.name} ${data.status}`, 'server', 'primary');
    loadDashboardData(); // Refresh stats
});

wsManager.on('event_update', (data) => {
    addActivityItem('event', `Event ${data.id}: ${data.status}`, 'bell', 'info');
    loadRecentEvents(); // Refresh events table
});

wsManager.on('hook_execution', (data) => {
    const status = data.success ? 'succeeded' : 'failed';
    const color = data.success ? 'success' : 'danger';
    addActivityItem('hook', `Hook ${data.hook_name} ${status}`, 'lightning', color);
});

wsManager.on('workflow_update', (data) => {
    addActivityItem('workflow', `Workflow ${data.name} ${data.status}`, 'diagram-3', 'success');
    loadDashboardData(); // Refresh stats
});

wsManager.on('system_alert', (data) => {
    addActivityItem('alert', data.message, 'exclamation-triangle', 'warning');
});

// Initial load
document.addEventListener('DOMContentLoaded', () => {
    loadDashboardData();

    // Refresh data every 30 seconds
    setInterval(loadDashboardData, 30000);
});
