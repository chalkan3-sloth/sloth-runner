// Stack Management
let currentStacks = [];

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

async function loadStacks() {
    try {
        const data = await API.get('/api/v1/stacks');
        currentStacks = data.stacks || [];
        renderStacks();
        updateStats();
    } catch (error) {
        console.error('Failed to load stacks:', error);
        showError('Failed to load stacks');
    }
}

function renderStacks() {
    const container = document.getElementById('stacks-container');

    if (currentStacks.length === 0) {
        container.innerHTML = `
            <div class="col-12 text-center py-5">
                <i class="bi bi-layers fs-1 text-muted"></i>
                <p class="text-muted mt-3">No stacks found</p>
                <button class="btn btn-primary mt-2" onclick="showCreateStackModal()">
                    <i class="bi bi-plus-lg"></i> Create First Stack
                </button>
            </div>
        `;
        return;
    }

    container.innerHTML = currentStacks.map(stack => {
        const typeIcon = getStackTypeIcon(stack.type);
        const typeColor = getStackTypeColor(stack.type);
        const statusBadge = stack.active ?
            '<span class="badge bg-success">Active</span>' :
            '<span class="badge bg-secondary">Inactive</span>';

        return `
            <div class="col-md-6 col-lg-4 mb-4">
                <div class="card h-100">
                    <div class="card-header d-flex justify-content-between align-items-center bg-${typeColor}">
                        <h5 class="mb-0 text-white">
                            <i class="bi bi-${typeIcon}"></i> ${stack.name}
                        </h5>
                        ${statusBadge}
                    </div>
                    <div class="card-body">
                        <p class="text-muted small">${stack.description || 'No description'}</p>

                        <div class="row text-center mb-3">
                            <div class="col-6">
                                <div class="border rounded p-2">
                                    <h6 class="mb-0">${stack.variable_count || 0}</h6>
                                    <small class="text-muted">Variables</small>
                                </div>
                            </div>
                            <div class="col-6">
                                <div class="border rounded p-2">
                                    <h6 class="mb-0">${stack.secret_count || 0}</h6>
                                    <small class="text-muted">Secrets</small>
                                </div>
                            </div>
                        </div>

                        <div class="d-grid gap-2">
                            <button class="btn btn-outline-primary btn-sm" onclick="viewStack('${stack.name}')">
                                <i class="bi bi-eye"></i> View Details
                            </button>
                            <div class="btn-group btn-group-sm">
                                <button class="btn btn-outline-secondary" onclick="showAddVariableModal('${stack.name}')" title="Add Variable">
                                    <i class="bi bi-plus-circle"></i> Variable
                                </button>
                                <button class="btn btn-outline-warning" onclick="showEditStackModal('${stack.name}')" title="Edit Stack">
                                    <i class="bi bi-pencil"></i> Edit
                                </button>
                                <button class="btn btn-outline-danger" onclick="deleteStack('${stack.name}')" title="Delete Stack">
                                    <i class="bi bi-trash"></i> Delete
                                </button>
                            </div>
                        </div>
                    </div>
                    <div class="card-footer text-muted small">
                        <i class="bi bi-clock"></i> ${stack.created_at ? new Date(stack.created_at).toLocaleDateString() : 'N/A'}
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

function updateStats() {
    const totalStacks = currentStacks.length;
    const activeStacks = currentStacks.filter(s => s.active).length;
    const totalVariables = currentStacks.reduce((sum, s) => sum + (s.variable_count || 0), 0);
    const totalSecrets = currentStacks.reduce((sum, s) => sum + (s.secret_count || 0), 0);

    document.getElementById('total-stacks').textContent = totalStacks;
    document.getElementById('active-stacks').textContent = activeStacks;
    document.getElementById('total-variables').textContent = totalVariables;
    document.getElementById('total-secrets').textContent = totalSecrets;
}

function showCreateStackModal() {
    const modal = new bootstrap.Modal(document.getElementById('createStackModal'));
    document.getElementById('createStackForm').reset();
    modal.show();
}

async function createStack() {
    const name = document.getElementById('createStackName').value.trim();
    const description = document.getElementById('createStackDescription').value.trim();
    const type = document.getElementById('createStackType').value;
    const active = document.getElementById('createStackActive').checked;

    if (!name) {
        showError('Stack name is required');
        return;
    }

    try {
        await API.post('/api/v1/stacks', {
            name,
            description,
            type,
            active
        });

        bootstrap.Modal.getInstance(document.getElementById('createStackModal')).hide();
        showSuccess('Stack created successfully');
        loadStacks();
    } catch (error) {
        console.error('Failed to create stack:', error);
        showError('Failed to create stack: ' + error.message);
    }
}

async function viewStack(name) {
    const modal = new bootstrap.Modal(document.getElementById('viewStackModal'));
    const content = document.getElementById('viewStackContent');

    content.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-primary"></div></div>';
    modal.show();

    try {
        const data = await API.get(`/api/v1/stacks/${name}`);

        content.innerHTML = `
            <div class="mb-4">
                <h5><i class="bi bi-layers"></i> ${data.name}</h5>
                <p class="text-muted">${data.description || 'No description'}</p>
                <div class="d-flex gap-2">
                    <span class="badge bg-${getStackTypeColor(data.type)}">${data.type}</span>
                    ${data.active ? '<span class="badge bg-success">Active</span>' : '<span class="badge bg-secondary">Inactive</span>'}
                </div>
            </div>

            <ul class="nav nav-tabs mb-3" role="tablist">
                <li class="nav-item">
                    <button class="nav-link active" data-bs-toggle="tab" data-bs-target="#variables-tab">
                        <i class="bi bi-code-square"></i> Variables (${data.variables?.length || 0})
                    </button>
                </li>
                <li class="nav-item">
                    <button class="nav-link" data-bs-toggle="tab" data-bs-target="#secrets-tab">
                        <i class="bi bi-shield-lock"></i> Secrets (${data.secrets?.length || 0})
                    </button>
                </li>
            </ul>

            <div class="tab-content">
                <div class="tab-pane fade show active" id="variables-tab">
                    ${renderVariables(data.variables || [])}
                </div>
                <div class="tab-pane fade" id="secrets-tab">
                    ${renderSecrets(data.secrets || [])}
                </div>
            </div>
        `;
    } catch (error) {
        console.error('Failed to load stack:', error);
        content.innerHTML = '<div class="alert alert-danger">Failed to load stack details</div>';
    }
}

function renderVariables(variables) {
    if (variables.length === 0) {
        return '<div class="text-muted text-center py-4">No variables in this stack</div>';
    }

    return `
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Key</th>
                        <th>Value</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${variables.map(v => `
                        <tr>
                            <td><code>${v.key}</code></td>
                            <td><code>${escapeHtml(v.value)}</code></td>
                            <td>
                                <button class="btn btn-sm btn-outline-danger" onclick="deleteVariable('${v.stack}', '${v.key}')">
                                    <i class="bi bi-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

function renderSecrets(secrets) {
    if (secrets.length === 0) {
        return '<div class="text-muted text-center py-4">No secrets in this stack</div>';
    }

    return `
        <div class="alert alert-warning">
            <i class="bi bi-exclamation-triangle"></i> Secret values are masked for security
        </div>
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Type</th>
                        <th>Created</th>
                    </tr>
                </thead>
                <tbody>
                    ${secrets.map(s => `
                        <tr>
                            <td><i class="bi bi-key-fill text-warning"></i> ${s.name}</td>
                            <td><span class="badge bg-secondary">${s.type || 'generic'}</span></td>
                            <td>${s.created_at ? new Date(s.created_at).toLocaleDateString() : 'N/A'}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

function showEditStackModal(name) {
    const stack = currentStacks.find(s => s.name === name);
    if (!stack) return;

    document.getElementById('editStackName').value = stack.name;
    document.getElementById('editStackDescription').value = stack.description || '';
    document.getElementById('editStackType').value = stack.type || 'custom';
    document.getElementById('editStackActive').checked = stack.active;

    const modal = new bootstrap.Modal(document.getElementById('editStackModal'));
    modal.show();
}

async function updateStack() {
    const name = document.getElementById('editStackName').value;
    const description = document.getElementById('editStackDescription').value.trim();
    const type = document.getElementById('editStackType').value;
    const active = document.getElementById('editStackActive').checked;

    try {
        await API.put(`/api/v1/stacks/${name}`, {
            description,
            type,
            active
        });

        bootstrap.Modal.getInstance(document.getElementById('editStackModal')).hide();
        showSuccess('Stack updated successfully');
        loadStacks();
    } catch (error) {
        console.error('Failed to update stack:', error);
        showError('Failed to update stack: ' + error.message);
    }
}

async function deleteStack(name) {
    if (!confirm(`Are you sure you want to delete stack "${name}"? This will also delete all associated variables and secrets.`)) {
        return;
    }

    try {
        await API.delete(`/api/v1/stacks/${name}`);
        showSuccess('Stack deleted successfully');
        loadStacks();
    } catch (error) {
        console.error('Failed to delete stack:', error);
        showError('Failed to delete stack: ' + error.message);
    }
}

function showAddVariableModal(stackName) {
    document.getElementById('addVariableStack').value = stackName;
    document.getElementById('addVariableForm').reset();
    document.getElementById('addVariableStack').value = stackName;

    const modal = new bootstrap.Modal(document.getElementById('addVariableModal'));
    modal.show();
}

async function addVariable() {
    const stack = document.getElementById('addVariableStack').value;
    const key = document.getElementById('addVariableKey').value.trim();
    const value = document.getElementById('addVariableValue').value;

    if (!key || !value) {
        showError('Key and value are required');
        return;
    }

    try {
        await API.post(`/api/v1/stacks/${stack}/variables`, {
            key,
            value
        });

        bootstrap.Modal.getInstance(document.getElementById('addVariableModal')).hide();
        showSuccess('Variable added successfully');
        loadStacks();
    } catch (error) {
        console.error('Failed to add variable:', error);
        showError('Failed to add variable: ' + error.message);
    }
}

async function deleteVariable(stack, key) {
    if (!confirm(`Delete variable "${key}"?`)) return;

    try {
        await API.delete(`/api/v1/stacks/${stack}/variables/${key}`);
        showSuccess('Variable deleted successfully');
        loadStacks();

        // Refresh view if modal is open
        const modal = bootstrap.Modal.getInstance(document.getElementById('viewStackModal'));
        if (modal) {
            viewStack(stack);
        }
    } catch (error) {
        console.error('Failed to delete variable:', error);
        showError('Failed to delete variable');
    }
}

function getStackTypeIcon(type) {
    switch (type) {
        case 'production': return 'shield-check';
        case 'staging': return 'gear';
        case 'development': return 'code-slash';
        case 'testing': return 'bug';
        default: return 'layers';
    }
}

function getStackTypeColor(type) {
    switch (type) {
        case 'production': return 'danger';
        case 'staging': return 'warning';
        case 'development': return 'info';
        case 'testing': return 'secondary';
        default: return 'primary';
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showSuccess(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-success border-0 position-fixed top-0 end-0 m-3';
    toast.style.zIndex = '9999';
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    document.body.appendChild(toast);
    new bootstrap.Toast(toast).show();
    setTimeout(() => toast.remove(), 5000);
}

function showError(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-danger border-0 position-fixed top-0 end-0 m-3';
    toast.style.zIndex = '9999';
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    document.body.appendChild(toast);
    new bootstrap.Toast(toast).show();
    setTimeout(() => toast.remove(), 5000);
}

document.addEventListener('DOMContentLoaded', () => {
    loadStacks();

    // Refresh every 30 seconds
    setInterval(loadStacks, 30000);
});
