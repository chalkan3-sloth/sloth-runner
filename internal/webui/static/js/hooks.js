// Hooks Management
let currentHooks = [];

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

// Load hooks
async function loadHooks() {
    try {
        const data = await API.get('/api/v1/hooks');
        currentHooks = data.hooks || [];
        renderHooks();
    } catch (error) {
        console.error('Failed to load hooks:', error);
        showError('Failed to load hooks');
    }
}

// Render hooks
function renderHooks() {
    const container = document.getElementById('hooks-container');

    if (currentHooks.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No hooks configured.
                    <a href="#" onclick="showAddHookModal()" class="alert-link">Add one now</a>
                </div>
            </div>
        `;
        return;
    }

    container.innerHTML = currentHooks.map(hook => `
        <div class="col-md-6 col-lg-4 mb-3">
            <div class="card h-100 ${hook.enabled ? '' : 'opacity-75'}">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-3">
                        <div class="flex-grow-1">
                            <h5 class="card-title mb-1">
                                <i class="bi bi-lightning-charge"></i> ${hook.name || hook.id}
                            </h5>
                            <div class="mb-2">
                                <span class="badge ${hook.enabled ? 'bg-success' : 'bg-secondary'}">
                                    ${hook.enabled ? 'Enabled' : 'Disabled'}
                                </span>
                                <span class="badge bg-info">${hook.event_type || 'unknown'}</span>
                            </div>
                        </div>
                        <div class="dropdown">
                            <button class="btn btn-sm btn-outline-secondary" data-bs-toggle="dropdown">
                                <i class="bi bi-three-dots-vertical"></i>
                            </button>
                            <ul class="dropdown-menu dropdown-menu-end">
                                <li><a class="dropdown-item" href="#" onclick="viewHook('${hook.id}')">
                                    <i class="bi bi-eye"></i> View Details
                                </a></li>
                                <li><a class="dropdown-item" href="#" onclick="viewHookHistory('${hook.id}')">
                                    <i class="bi bi-clock-history"></i> View History
                                </a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item" href="#" onclick="toggleHook('${hook.id}', ${hook.enabled})">
                                    <i class="bi bi-${hook.enabled ? 'pause' : 'play'}-circle"></i>
                                    ${hook.enabled ? 'Disable' : 'Enable'}
                                </a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item text-danger" href="#" onclick="deleteHook('${hook.id}')">
                                    <i class="bi bi-trash"></i> Delete
                                </a></li>
                            </ul>
                        </div>
                    </div>

                    <p class="text-muted small mb-2">
                        <strong>Action:</strong> ${hook.action || 'N/A'}
                    </p>

                    ${hook.filter ? `
                        <p class="text-muted small mb-2">
                            <strong>Filter:</strong> <code class="small">${JSON.stringify(hook.filter)}</code>
                        </p>
                    ` : ''}

                    <div class="mt-3">
                        <small class="text-muted">
                            <i class="bi bi-clock"></i>
                            ${hook.created_at ? new Date(hook.created_at).toLocaleString() : 'Unknown'}
                        </small>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
}

// View hook details
async function viewHook(id) {
    try {
        const hook = await API.get(`/api/v1/hooks/${id}`);

        const modal = new bootstrap.Modal(document.getElementById('viewHookModal'));
        document.getElementById('viewHookContent').innerHTML = `
            <h5>${hook.name || hook.id}</h5>
            <p><strong>Status:</strong>
                <span class="badge ${hook.enabled ? 'bg-success' : 'bg-secondary'}">
                    ${hook.enabled ? 'Enabled' : 'Disabled'}
                </span>
            </p>
            <p><strong>Event Type:</strong> <span class="badge bg-info">${hook.event_type}</span></p>
            <p><strong>Action:</strong> <code>${hook.action}</code></p>
            ${hook.config ? `<p><strong>Config:</strong> <pre class="bg-light p-2 rounded"><code>${JSON.stringify(hook.config, null, 2)}</code></pre></p>` : ''}
            ${hook.filter ? `<p><strong>Filter:</strong> <pre class="bg-light p-2 rounded"><code>${JSON.stringify(hook.filter, null, 2)}</code></pre></p>` : ''}
            <p><strong>Created:</strong> ${new Date(hook.created_at).toLocaleString()}</p>
            ${hook.updated_at ? `<p><strong>Updated:</strong> ${new Date(hook.updated_at).toLocaleString()}</p>` : ''}
        `;

        modal.show();
    } catch (error) {
        console.error('Failed to load hook:', error);
        showError('Failed to load hook details');
    }
}

// View hook history
async function viewHookHistory(id) {
    try {
        const history = await API.get(`/api/v1/hooks/${id}/history`);

        const modal = new bootstrap.Modal(document.getElementById('hookHistoryModal'));
        const container = document.getElementById('hookHistoryContent');

        if (!history || history.length === 0) {
            container.innerHTML = `
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No execution history available
                </div>
            `;
        } else {
            container.innerHTML = `
                <div class="table-responsive">
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th>Time</th>
                                <th>Status</th>
                                <th>Event</th>
                                <th>Duration</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${history.map(h => `
                                <tr>
                                    <td>${new Date(h.timestamp).toLocaleString()}</td>
                                    <td><span class="badge bg-${h.success ? 'success' : 'danger'}">${h.success ? 'Success' : 'Failed'}</span></td>
                                    <td>${h.event_type || 'N/A'}</td>
                                    <td>${h.duration || 'N/A'}</td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>
                </div>
            `;
        }

        modal.show();
    } catch (error) {
        console.error('Failed to load hook history:', error);
        showError('Failed to load hook history');
    }
}

// Toggle hook status
async function toggleHook(id, isEnabled) {
    try {
        const endpoint = isEnabled ? 'disable' : 'enable';
        await API.post(`/api/v1/hooks/${id}/${endpoint}`);

        showSuccess(`Hook ${isEnabled ? 'disabled' : 'enabled'} successfully`);
        loadHooks();
    } catch (error) {
        console.error('Failed to toggle hook:', error);
        showError('Failed to update hook status');
    }
}

// Delete hook
async function deleteHook(id) {
    if (!confirm('Are you sure you want to delete this hook?')) return;

    try {
        await API.delete(`/api/v1/hooks/${id}`);
        showSuccess('Hook deleted successfully');
        loadHooks();
    } catch (error) {
        console.error('Failed to delete hook:', error);
        showError('Failed to delete hook');
    }
}

// Show add hook modal
function showAddHookModal() {
    const modal = new bootstrap.Modal(document.getElementById('addHookModal'));
    document.getElementById('addHookForm').reset();
    modal.show();
}

// Add hook
async function addHook() {
    const name = document.getElementById('addHookName').value;
    const eventType = document.getElementById('addHookEventType').value;
    const action = document.getElementById('addHookAction').value;
    const filterText = document.getElementById('addHookFilter').value;

    if (!name || !eventType || !action) {
        showError('Name, event type, and action are required');
        return;
    }

    try {
        const payload = {
            name,
            event_type: eventType,
            action,
            enabled: true
        };

        if (filterText) {
            try {
                payload.filter = JSON.parse(filterText);
            } catch (e) {
                showError('Invalid JSON in filter field');
                return;
            }
        }

        await API.post('/api/v1/hooks', payload);

        bootstrap.Modal.getInstance(document.getElementById('addHookModal')).hide();
        showSuccess('Hook added successfully');
        loadHooks();
    } catch (error) {
        console.error('Failed to add hook:', error);
        showError('Failed to add hook: ' + error.message);
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
    loadHooks();
});
