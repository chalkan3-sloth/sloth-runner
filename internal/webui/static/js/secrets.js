// Secrets Management
let currentStacks = [];
let currentStack = null;

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
    async delete(url) {
        const response = await fetch(url, { method: 'DELETE' });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.ok;
    }
};

async function loadStacks() {
    // Simulating stack discovery - in real implementation this would come from API
    currentStacks = ['production', 'staging', 'development', 'default'];
    renderStacks();
}

function renderStacks() {
    const container = document.getElementById('stacks-list');

    container.innerHTML = currentStacks.map(stack => `
        <button class="list-group-item list-group-item-action d-flex justify-content-between align-items-center ${stack === currentStack ? 'active' : ''}"
                onclick="selectStack('${stack}')">
            <span>
                <i class="bi bi-folder"></i> ${stack}
            </span>
            ${stack === currentStack ? '<i class="bi bi-chevron-right"></i>' : ''}
        </button>
    `).join('');
}

async function selectStack(stack) {
    currentStack = stack;
    renderStacks();
    await loadSecrets(stack);
}

async function loadSecrets(stack) {
    try {
        const data = await API.get(`/api/v1/secrets/${stack}`);
        const secrets = data.secrets || [];
        renderSecrets(secrets);
    } catch (error) {
        console.error('Failed to load secrets:', error);
        showError('Failed to load secrets');
    }
}

function renderSecrets(secrets) {
    const container = document.getElementById('secrets-container');

    if (!currentStack) {
        container.innerHTML = `
            <div class="text-center py-5">
                <i class="bi bi-folder2-open fs-1 text-muted"></i>
                <p class="text-muted mt-3">Select a stack to view secrets</p>
            </div>
        `;
        return;
    }

    if (secrets.length === 0) {
        container.innerHTML = `
            <div class="text-center py-5">
                <i class="bi bi-shield-lock fs-1 text-muted"></i>
                <p class="text-muted mt-3">No secrets in stack "${currentStack}"</p>
                <button class="btn btn-primary mt-2" onclick="showAddSecretModal()">
                    <i class="bi bi-plus-lg"></i> Add Secret
                </button>
            </div>
        `;
        return;
    }

    container.innerHTML = `
        <div class="d-flex justify-content-between align-items-center mb-3">
            <h5><i class="bi bi-shield-lock"></i> Secrets in "${currentStack}"</h5>
            <button class="btn btn-success btn-sm" onclick="showAddSecretModal()">
                <i class="bi bi-plus-lg"></i> Add Secret
            </button>
        </div>
        <div class="table-responsive">
            <table class="table table-hover">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Type</th>
                        <th>Created</th>
                        <th>Last Updated</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${secrets.map(secret => `
                        <tr>
                            <td>
                                <i class="bi bi-key-fill text-warning"></i>
                                <strong>${secret.name}</strong>
                            </td>
                            <td><span class="badge bg-secondary">${secret.type || 'generic'}</span></td>
                            <td>${secret.created_at ? new Date(secret.created_at).toLocaleDateString() : 'N/A'}</td>
                            <td>${secret.updated_at ? new Date(secret.updated_at).toLocaleDateString() : 'N/A'}</td>
                            <td>
                                <div class="btn-group btn-group-sm">
                                    <button class="btn btn-outline-primary" onclick="viewSecret('${currentStack}', '${secret.name}')" title="View">
                                        <i class="bi bi-eye"></i>
                                    </button>
                                    <button class="btn btn-outline-warning" onclick="editSecret('${currentStack}', '${secret.name}')" title="Edit">
                                        <i class="bi bi-pencil"></i>
                                    </button>
                                    <button class="btn btn-outline-danger" onclick="deleteSecret('${currentStack}', '${secret.name}')" title="Delete">
                                        <i class="bi bi-trash"></i>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

function showAddSecretModal() {
    if (!currentStack) {
        showError('Please select a stack first');
        return;
    }

    const modal = new bootstrap.Modal(document.getElementById('addSecretModal'));
    document.getElementById('addSecretStack').value = currentStack;
    document.getElementById('addSecretForm').reset();
    document.getElementById('addSecretStack').value = currentStack;
    modal.show();
}

async function addSecret() {
    const stack = document.getElementById('addSecretStack').value;
    const name = document.getElementById('addSecretName').value;
    const value = document.getElementById('addSecretValue').value;
    const type = document.getElementById('addSecretType').value;

    if (!name || !value) {
        showError('Name and value are required');
        return;
    }

    try {
        await API.post(`/api/v1/secrets/${stack}`, {
            name,
            value,
            type: type || 'generic'
        });

        bootstrap.Modal.getInstance(document.getElementById('addSecretModal')).hide();
        showSuccess('Secret added successfully');
        loadSecrets(stack);
    } catch (error) {
        console.error('Failed to add secret:', error);
        showError('Failed to add secret: ' + error.message);
    }
}

async function viewSecret(stack, name) {
    try {
        const data = await API.get(`/api/v1/secrets/${stack}/${name}`);

        const modal = new bootstrap.Modal(document.getElementById('viewSecretModal'));
        document.getElementById('viewSecretContent').innerHTML = `
            <div class="alert alert-warning">
                <i class="bi bi-exclamation-triangle"></i>
                <strong>Security Warning:</strong> Secret values are masked for security. Use CLI for full access.
            </div>
            <div class="mb-3">
                <strong>Stack:</strong> ${stack}
            </div>
            <div class="mb-3">
                <strong>Name:</strong> ${name}
            </div>
            <div class="mb-3">
                <strong>Type:</strong> <span class="badge bg-secondary">${data.type || 'generic'}</span>
            </div>
            <div class="mb-3">
                <strong>Value:</strong> <code>********</code>
                <small class="text-muted d-block">Value is hidden for security</small>
            </div>
        `;

        modal.show();
    } catch (error) {
        console.error('Failed to load secret:', error);
        showError('Failed to load secret');
    }
}

async function deleteSecret(stack, name) {
    if (!confirm(`Are you sure you want to delete secret "${name}" from stack "${stack}"?`)) return;

    try {
        await API.delete(`/api/v1/secrets/${stack}/${name}`);
        showSuccess('Secret deleted successfully');
        loadSecrets(stack);
    } catch (error) {
        console.error('Failed to delete secret:', error);
        showError('Failed to delete secret');
    }
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
});
