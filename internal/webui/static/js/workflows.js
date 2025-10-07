// Workflows Management
let currentSloths = [];
let currentFilter = 'all';

// API Helper
const API = {
    async get(url) {
        const response = await fetch(url);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    },

    async post(url, data = {}) {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    },

    async put(url, data = {}) {
        const response = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    },

    async delete(url) {
        const response = await fetch(url, { method: 'DELETE' });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.ok;
    }
};

// Load workflows
async function loadWorkflows() {
    try {
        const data = await API.get('/api/v1/sloths');
        currentSloths = data.sloths || [];
        renderWorkflows();
    } catch (error) {
        console.error('Failed to load workflows:', error);
        showError('Failed to load workflows');
    }
}

// Filter workflows
function filterWorkflows(filter) {
    currentFilter = filter;

    // Update filter buttons
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');

    renderWorkflows();
}

// Render workflows
function renderWorkflows() {
    const container = document.getElementById('workflows-container');

    let filtered = currentSloths;
    if (currentFilter === 'active') {
        filtered = currentSloths.filter(s => s.is_active);
    } else if (currentFilter === 'inactive') {
        filtered = currentSloths.filter(s => !s.is_active);
    }

    if (filtered.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No workflows found.
                    <a href="#" onclick="showAddWorkflowModal()" class="alert-link">Add one now</a>
                </div>
            </div>
        `;
        return;
    }

    container.innerHTML = filtered.map(sloth => `
        <div class="col-md-6 col-lg-4 mb-3">
            <div class="card workflow-card h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-3">
                        <div>
                            <h5 class="card-title mb-1">
                                <i class="bi bi-diagram-3"></i> ${sloth.name}
                            </h5>
                            <span class="badge ${sloth.is_active ? 'bg-success' : 'bg-secondary'}">
                                ${sloth.is_active ? 'Active' : 'Inactive'}
                            </span>
                        </div>
                        <div class="dropdown">
                            <button class="btn btn-sm btn-outline-secondary" data-bs-toggle="dropdown">
                                <i class="bi bi-three-dots-vertical"></i>
                            </button>
                            <ul class="dropdown-menu dropdown-menu-end">
                                <li><a class="dropdown-item" href="#" onclick="viewWorkflow('${sloth.name}')">
                                    <i class="bi bi-eye"></i> View
                                </a></li>
                                <li><a class="dropdown-item" href="#" onclick="editWorkflow('${sloth.name}')">
                                    <i class="bi bi-pencil"></i> Edit
                                </a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item" href="#" onclick="toggleWorkflowStatus('${sloth.name}', ${sloth.is_active})">
                                    <i class="bi bi-${sloth.is_active ? 'pause' : 'play'}-circle"></i>
                                    ${sloth.is_active ? 'Deactivate' : 'Activate'}
                                </a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item text-danger" href="#" onclick="deleteWorkflow('${sloth.name}')">
                                    <i class="bi bi-trash"></i> Delete
                                </a></li>
                            </ul>
                        </div>
                    </div>

                    <p class="text-muted small mb-3">${sloth.path || 'No path specified'}</p>

                    <div class="d-flex gap-2">
                        <button class="btn btn-primary btn-sm flex-grow-1" onclick="runWorkflow('${sloth.name}')">
                            <i class="bi bi-play-fill"></i> Run
                        </button>
                        <button class="btn btn-outline-info btn-sm" onclick="viewWorkflow('${sloth.name}')">
                            <i class="bi bi-eye"></i>
                        </button>
                    </div>
                </div>
                <div class="card-footer bg-transparent border-top-0 pt-0">
                    <small class="text-muted">
                        <i class="bi bi-clock"></i>
                        ${sloth.created_at ? new Date(sloth.created_at).toLocaleDateString() : 'Unknown'}
                    </small>
                </div>
            </div>
        </div>
    `).join('');
}

// Run workflow
async function runWorkflow(name) {
    try {
        const sloth = currentSloths.find(s => s.name === name);
        if (!sloth) return;

        // Show run modal
        const modal = new bootstrap.Modal(document.getElementById('runWorkflowModal'));
        document.getElementById('runWorkflowName').value = name;
        document.getElementById('runWorkflowPath').textContent = sloth.path;

        // Load groups from sloth file
        // For now, we'll show a simple form
        modal.show();
    } catch (error) {
        console.error('Failed to prepare workflow run:', error);
        showError('Failed to prepare workflow run');
    }
}

// Execute workflow
async function executeWorkflow() {
    const name = document.getElementById('runWorkflowName').value;
    const group = document.getElementById('runWorkflowGroup').value;
    const agents = document.getElementById('runWorkflowAgents').value;

    try {
        const payload = {
            sloth_name: name,
            group: group || undefined,
            agents: agents ? agents.split(',').map(a => a.trim()) : undefined
        };

        const result = await API.post('/api/v1/executions', payload);

        // Close modal
        bootstrap.Modal.getInstance(document.getElementById('runWorkflowModal')).hide();

        // Show success and redirect to executions
        showSuccess(`Workflow "${name}" started successfully`);
        setTimeout(() => {
            window.location.href = '/executions';
        }, 1500);
    } catch (error) {
        console.error('Failed to execute workflow:', error);
        showError('Failed to execute workflow: ' + error.message);
    }
}

// View workflow
async function viewWorkflow(name) {
    try {
        const sloth = await API.get(`/api/v1/sloths/${name}`);

        const modal = new bootstrap.Modal(document.getElementById('viewWorkflowModal'));
        document.getElementById('viewWorkflowContent').innerHTML = `
            <h5>${sloth.name}</h5>
            <p><strong>Status:</strong>
                <span class="badge ${sloth.is_active ? 'bg-success' : 'bg-secondary'}">
                    ${sloth.is_active ? 'Active' : 'Inactive'}
                </span>
            </p>
            <p><strong>Path:</strong> <code>${sloth.path}</code></p>
            <p><strong>Created:</strong> ${new Date(sloth.created_at).toLocaleString()}</p>
            <hr>
            <pre class="bg-light p-3 rounded"><code>${sloth.content || 'No content available'}</code></pre>
        `;

        modal.show();
    } catch (error) {
        console.error('Failed to load workflow:', error);
        showError('Failed to load workflow details');
    }
}

// Toggle workflow status
async function toggleWorkflowStatus(name, isActive) {
    try {
        const endpoint = isActive ? 'deactivate' : 'activate';
        await API.post(`/api/v1/sloths/${name}/${endpoint}`);

        showSuccess(`Workflow ${isActive ? 'deactivated' : 'activated'} successfully`);
        loadWorkflows();
    } catch (error) {
        console.error('Failed to toggle workflow status:', error);
        showError('Failed to update workflow status');
    }
}

// Delete workflow
async function deleteWorkflow(name) {
    if (!confirm(`Are you sure you want to delete workflow "${name}"?`)) return;

    try {
        await API.delete(`/api/v1/sloths/${name}`);
        showSuccess('Workflow deleted successfully');
        loadWorkflows();
    } catch (error) {
        console.error('Failed to delete workflow:', error);
        showError('Failed to delete workflow');
    }
}

// Show add workflow modal
function showAddWorkflowModal() {
    const modal = new bootstrap.Modal(document.getElementById('addWorkflowModal'));
    document.getElementById('addWorkflowForm').reset();
    modal.show();
}

// Add workflow
async function addWorkflow() {
    const name = document.getElementById('addWorkflowName').value;
    const path = document.getElementById('addWorkflowPath').value;
    const content = document.getElementById('addWorkflowContent').value;

    if (!name || !path) {
        showError('Name and path are required');
        return;
    }

    try {
        await API.post('/api/v1/sloths', { name, path, content });

        bootstrap.Modal.getInstance(document.getElementById('addWorkflowModal')).hide();
        showSuccess('Workflow added successfully');
        loadWorkflows();
    } catch (error) {
        console.error('Failed to add workflow:', error);
        showError('Failed to add workflow');
    }
}

// Utility functions
function showSuccess(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-success border-0 position-fixed top-0 end-0 m-3';
    toast.setAttribute('role', 'alert');
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    document.body.appendChild(toast);
    const bsToast = new bootstrap.Toast(toast);
    bsToast.show();
    setTimeout(() => toast.remove(), 5000);
}

function showError(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-danger border-0 position-fixed top-0 end-0 m-3';
    toast.setAttribute('role', 'alert');
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    document.body.appendChild(toast);
    const bsToast = new bootstrap.Toast(toast);
    bsToast.show();
    setTimeout(() => toast.remove(), 5000);
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    loadWorkflows();
});
