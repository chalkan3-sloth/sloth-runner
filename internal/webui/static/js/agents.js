// Agents page logic

let currentAgentName = null;

// Load all agents
async function loadAgents() {
    const container = document.getElementById('agents-container');
    const statsContainer = document.getElementById('agent-stats');

    // Show skeleton loader
    await SkeletonLoader.loadWithSkeleton(
        'agents-container',
        SkeletonLoader.agentCards(6),
        async () => {
            try {
                const response = await fetch('/api/v1/agents');
                const data = await response.json();

                if (!data.agents || data.agents.length === 0) {
                    container.innerHTML = `
                        <div class="col-12">
                            <div class="empty-state fade-in-up">
                                <i class="bi bi-server"></i>
                                <h5>No agents registered</h5>
                                <p>Agents will appear here when they connect to the master server.</p>
                            </div>
                        </div>
                    `;
                    return;
                }

                // Update stats
                updateAgentStats(data.agents, statsContainer);

                // Render agent cards with staggered animation
                container.innerHTML = data.agents.map((agent, index) =>
                    createAgentCard(agent, index + 1)
                ).join('');

                // Initialize sparklines if available
                if (typeof initializeSparklines !== 'undefined') {
                    setTimeout(() => initializeSparklines(data.agents), 100);
                }

                toastManager.success(`Loaded ${data.agents.length} agents`, 'Success');
            } catch (error) {
                console.error('Error loading agents:', error);
                container.innerHTML = `
                    <div class="col-12">
                        <div class="alert alert-danger fade-in-up" role="alert">
                            <i class="bi bi-exclamation-triangle"></i> Error loading agents: ${error.message}
                        </div>
                    </div>
                `;
                toastManager.error('Failed to load agents', 'Error');
            }
        }
    );
}

// Update agent statistics
function updateAgentStats(agents, container) {
    const activeAgents = agents.filter(a => a.status === 'Active').length;
    const totalAgents = agents.length;
    const inactiveAgents = totalAgents - activeAgents;

    const activePercent = totalAgents > 0 ? (activeAgents / totalAgents * 100).toFixed(0) : 0;

    container.innerHTML = `
        <div class="col-lg-3 col-md-6 mb-3">
            <div class="card stat-card border-0 shadow-sm glass-card hover-lift fade-in-up stagger-1">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <h6 class="text-muted mb-2">Total Agents</h6>
                            <h2 class="mb-0 text-primary">${totalAgents}</h2>
                        </div>
                        <div class="stat-icon bg-primary bg-opacity-10 text-primary">
                            <i class="bi bi-hdd-network fs-1"></i>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-lg-3 col-md-6 mb-3">
            <div class="card stat-card border-0 shadow-sm glass-card hover-lift fade-in-up stagger-2">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <h6 class="text-muted mb-2">Active Agents</h6>
                            <h2 class="mb-0 text-success">${activeAgents}</h2>
                        </div>
                        <div class="stat-icon bg-success bg-opacity-10 text-success">
                            <i class="bi bi-check-circle fs-1"></i>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-lg-3 col-md-6 mb-3">
            <div class="card stat-card border-0 shadow-sm glass-card hover-lift fade-in-up stagger-3">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <h6 class="text-muted mb-2">Inactive Agents</h6>
                            <h2 class="mb-0 text-secondary">${inactiveAgents}</h2>
                        </div>
                        <div class="stat-icon bg-secondary bg-opacity-10 text-secondary">
                            <i class="bi bi-x-circle fs-1"></i>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-lg-3 col-md-6 mb-3">
            <div class="card stat-card border-0 shadow-sm glass-card hover-lift fade-in-up stagger-4">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <h6 class="text-muted mb-2">Uptime Rate</h6>
                            <h2 class="mb-0 text-info">${activePercent}%</h2>
                        </div>
                        <div class="stat-icon bg-info bg-opacity-10 text-info">
                            <i class="bi bi-graph-up fs-1"></i>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// Create agent card HTML
function createAgentCard(agent, index = 1) {
    // Use new AgentCardBuilder if available
    if (typeof AgentCardBuilder !== 'undefined') {
        const builder = new AgentCardBuilder(agent);
        return builder.build();
    }

    // Fallback to original card design
    const isActive = agent.status === 'Active';
    const statusClass = isActive ? 'success' : 'secondary';
    const statusIcon = isActive ? 'check-circle-fill' : 'x-circle-fill';

    const registeredDate = new Date(agent.registered_at * 1000).toLocaleString();
    const lastHeartbeat = agent.last_heartbeat > 0
        ? new Date(agent.last_heartbeat * 1000).toLocaleString()
        : 'Never';

    const cardClass = isActive ? 'agent-card' : 'agent-card inactive';
    const stagger = Math.min(index, 5); // Max stagger-5

    return `
        <div class="col-md-6 col-lg-4 mb-3">
            <div class="card ${cardClass} border-0 shadow-sm glass-card hover-lift fade-in-up stagger-${stagger}">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <h5 class="card-title mb-0">
                            <i class="bi bi-server me-2"></i>${agent.name}
                        </h5>
                        <span class="badge bg-${statusClass}">
                            ${isActive ? '<span class="pulse"></span>' : ''}
                            <i class="bi bi-${statusIcon}"></i> ${agent.status}
                        </span>
                    </div>

                    <p class="card-text text-muted small mb-2">
                        <i class="bi bi-geo-alt"></i> ${agent.address}
                    </p>

                    ${agent.version ? `
                        <p class="card-text text-muted small mb-2">
                            <i class="bi bi-tag"></i> Version: ${agent.version}
                        </p>
                    ` : ''}

                    <hr>

                    <div class="row g-2 small text-muted">
                        <div class="col-6">
                            <strong>Registered:</strong><br>
                            ${registeredDate}
                        </div>
                        <div class="col-6">
                            <strong>Last Heartbeat:</strong><br>
                            ${lastHeartbeat}
                        </div>
                    </div>

                    <div class="mt-3 d-grid gap-2">
                        <button class="btn btn-sm btn-primary hover-grow btn-ripple" onclick="showAgentDetails('${agent.name}')">
                            <i class="bi bi-info-circle"></i> View Details
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// View agent dashboard
function viewAgentDashboard(name) {
    window.location.href = `/agent-dashboard?agent=${encodeURIComponent(name)}`;
}

// View agent details (alias for compatibility)
function viewAgentDetails(name) {
    showAgentDetails(name);
}

// View agent logs
function viewAgentLogs(name) {
    window.location.href = `/logs?agent=${encodeURIComponent(name)}`;
}

// Restart agent
async function restartAgent(name) {
    if (!confirm(`Are you sure you want to restart agent "${name}"?`)) {
        return;
    }

    try {
        const response = await fetch(`/api/v1/agents/${name}/restart`, {
            method: 'POST'
        });

        if (response.ok) {
            if (window.toastManager) {
                toastManager.success(`Agent "${name}" restart initiated`, 'Success');
            }
            if (window.confetti) {
                confetti.burst({ particleCount: 30, spread: 40 });
            }
            setTimeout(() => refreshAgents(), 2000);
        } else {
            throw new Error('Restart failed');
        }
    } catch (error) {
        console.error('Error restarting agent:', error);
        if (window.toastManager) {
            toastManager.error(`Failed to restart agent "${name}"`, 'Error');
        }
    }
}

// Show agent details
async function showAgentDetails(name) {
    currentAgentName = name;

    const modal = new bootstrap.Modal(document.getElementById('agentDetailsModal'));
    const body = document.getElementById('agentDetailsBody');

    body.innerHTML = `
        <div class="text-center py-3">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    `;

    modal.show();

    try {
        const response = await fetch(`/api/v1/agents/${name}`);
        const agent = await response.json();

        const registeredDate = new Date(agent.registered_at * 1000).toLocaleString();
        const updatedDate = new Date(agent.updated_at * 1000).toLocaleString();
        const lastHeartbeat = agent.last_heartbeat > 0
            ? new Date(agent.last_heartbeat * 1000).toLocaleString()
            : 'Never';
        const lastInfoCollected = agent.last_info_collected > 0
            ? new Date(agent.last_info_collected * 1000).toLocaleString()
            : 'Never';

        let systemInfo = '';
        if (agent.system_info) {
            try {
                const info = JSON.parse(agent.system_info);
                systemInfo = `
                    <h6 class="mt-3">System Information</h6>
                    <pre class="small">${JSON.stringify(info, null, 2)}</pre>
                `;
            } catch (e) {
                systemInfo = `
                    <h6 class="mt-3">System Information</h6>
                    <pre class="small">${agent.system_info}</pre>
                `;
            }
        }

        body.innerHTML = `
            <div class="row">
                <div class="col-md-6">
                    <h6>Basic Information</h6>
                    <table class="table table-sm">
                        <tr>
                            <td><strong>Name:</strong></td>
                            <td>${agent.name}</td>
                        </tr>
                        <tr>
                            <td><strong>Address:</strong></td>
                            <td>${agent.address}</td>
                        </tr>
                        <tr>
                            <td><strong>Status:</strong></td>
                            <td><span class="badge bg-${agent.status === 'Active' ? 'success' : 'secondary'}">${agent.status}</span></td>
                        </tr>
                        <tr>
                            <td><strong>Version:</strong></td>
                            <td>${agent.version || 'Unknown'}</td>
                        </tr>
                    </table>
                </div>
                <div class="col-md-6">
                    <h6>Timestamps</h6>
                    <table class="table table-sm">
                        <tr>
                            <td><strong>Registered:</strong></td>
                            <td>${registeredDate}</td>
                        </tr>
                        <tr>
                            <td><strong>Updated:</strong></td>
                            <td>${updatedDate}</td>
                        </tr>
                        <tr>
                            <td><strong>Last Heartbeat:</strong></td>
                            <td>${lastHeartbeat}</td>
                        </tr>
                        <tr>
                            <td><strong>Last Info Collected:</strong></td>
                            <td>${lastInfoCollected}</td>
                        </tr>
                    </table>
                </div>
            </div>
            ${systemInfo}
        `;
    } catch (error) {
        console.error('Error loading agent details:', error);
        body.innerHTML = `
            <div class="alert alert-danger" role="alert">
                <i class="bi bi-exclamation-triangle"></i> Error loading agent details: ${error.message}
            </div>
        `;
    }
}

// Delete agent
async function deleteAgent() {
    if (!currentAgentName) return;

    if (!confirm(`Are you sure you want to delete agent "${currentAgentName}"?`)) {
        return;
    }

    const loadingToast = toastManager.loading(`Deleting agent "${currentAgentName}"...`);

    try {
        const response = await fetch(`/api/v1/agents/${currentAgentName}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            toastManager.update(loadingToast, {
                type: 'success',
                message: `Agent "${currentAgentName}" deleted successfully`,
                duration: 3000
            });

            const modal = bootstrap.Modal.getInstance(document.getElementById('agentDetailsModal'));
            modal.hide();
            loadAgents();
        } else {
            const error = await response.json();
            toastManager.update(loadingToast, {
                type: 'error',
                message: `Error deleting agent: ${error.error}`,
                duration: 5000
            });
        }
    } catch (error) {
        console.error('Error deleting agent:', error);
        toastManager.update(loadingToast, {
            type: 'error',
            message: `Error deleting agent: ${error.message}`,
            duration: 5000
        });
    }
}

// Refresh agents
function refreshAgents() {
    toastManager.info('Refreshing agents...', 'Refresh');
    loadAgents();
}

// Show add agent modal (placeholder)
function showAddAgentModal() {
    toastManager.info('Add agent feature coming soon!', 'Coming Soon');
}

// WebSocket handlers
wsManager.on('agent_update', (data) => {
    loadAgents();
});

wsManager.on('agent_connected', (data) => {
    confetti.agentConnected();
    toastManager.success(`Agent "${data.name}" connected!`, 'Agent Connected');
    loadAgents();
});

// Setup event listeners
document.addEventListener('DOMContentLoaded', () => {
    loadAgents();

    document.getElementById('deleteAgentBtn').addEventListener('click', deleteAgent);

    // Auto-refresh every 10 seconds
    setInterval(loadAgents, 10000);
});
