// Agents page logic

let currentAgentName = null;

// Load all agents
async function loadAgents() {
    try {
        const response = await fetch('/api/v1/agents');
        const data = await response.json();

        const container = document.getElementById('agents-container');

        if (!data.agents || data.agents.length === 0) {
            container.innerHTML = `
                <div class="col-12">
                    <div class="empty-state">
                        <i class="bi bi-server"></i>
                        <h5>No agents registered</h5>
                        <p>Agents will appear here when they connect to the master server.</p>
                    </div>
                </div>
            `;
            return;
        }

        container.innerHTML = data.agents.map(agent => createAgentCard(agent)).join('');
    } catch (error) {
        console.error('Error loading agents:', error);
        const container = document.getElementById('agents-container');
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-danger" role="alert">
                    <i class="bi bi-exclamation-triangle"></i> Error loading agents: ${error.message}
                </div>
            </div>
        `;
    }
}

// Create agent card HTML
function createAgentCard(agent) {
    const isActive = agent.status === 'Active';
    const statusClass = isActive ? 'success' : 'secondary';
    const statusIcon = isActive ? 'check-circle-fill' : 'x-circle-fill';

    const registeredDate = new Date(agent.registered_at * 1000).toLocaleString();
    const lastHeartbeat = agent.last_heartbeat > 0
        ? new Date(agent.last_heartbeat * 1000).toLocaleString()
        : 'Never';

    const cardClass = isActive ? 'agent-card' : 'agent-card inactive';

    return `
        <div class="col-md-6 col-lg-4 mb-3">
            <div class="card ${cardClass}">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <h5 class="card-title mb-0">
                            <i class="bi bi-server me-2"></i>${agent.name}
                        </h5>
                        <span class="badge bg-${statusClass}">
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

                    <div class="mt-3">
                        <button class="btn btn-sm btn-primary w-100" onclick="showAgentDetails('${agent.name}')">
                            <i class="bi bi-info-circle"></i> View Details
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;
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

    try {
        const response = await fetch(`/api/v1/agents/${currentAgentName}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('agentDetailsModal'));
            modal.hide();
            loadAgents();
        } else {
            const error = await response.json();
            alert(`Error deleting agent: ${error.error}`);
        }
    } catch (error) {
        console.error('Error deleting agent:', error);
        alert(`Error deleting agent: ${error.message}`);
    }
}

// Refresh agents
function refreshAgents() {
    loadAgents();
}

// WebSocket handlers
wsManager.on('agent_update', (data) => {
    loadAgents();
});

// Setup event listeners
document.addEventListener('DOMContentLoaded', () => {
    loadAgents();

    document.getElementById('deleteAgentBtn').addEventListener('click', deleteAgent);

    // Auto-refresh every 10 seconds
    setInterval(loadAgents, 10000);
});
