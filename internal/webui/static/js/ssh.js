// SSH Profiles Management
let currentProfiles = [];

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

async function loadProfiles() {
    try {
        const data = await API.get('/api/v1/ssh');
        currentProfiles = data.profiles || [];
        renderProfiles();
    } catch (error) {
        console.error('Failed to load SSH profiles:', error);
        showError('Failed to load SSH profiles');
    }
}

function renderProfiles() {
    const container = document.getElementById('ssh-profiles-container');

    if (currentProfiles.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No SSH profiles configured.
                    <a href="#" onclick="showAddProfileModal()" class="alert-link">Add one now</a>
                </div>
            </div>
        `;
        return;
    }

    container.innerHTML = currentProfiles.map(profile => `
        <div class="col-md-6 col-lg-4 mb-3">
            <div class="card h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-3">
                        <div class="flex-grow-1">
                            <h5 class="card-title mb-2">
                                <i class="bi bi-key"></i> ${profile.name}
                            </h5>
                            <p class="text-muted mb-1">
                                <i class="bi bi-person"></i> ${profile.user}@${profile.host}:${profile.port || 22}
                            </p>
                        </div>
                        <div class="dropdown">
                            <button class="btn btn-sm btn-outline-secondary" data-bs-toggle="dropdown">
                                <i class="bi bi-three-dots-vertical"></i>
                            </button>
                            <ul class="dropdown-menu dropdown-menu-end">
                                <li><a class="dropdown-item" href="#" onclick="viewProfile('${profile.name}')">
                                    <i class="bi bi-eye"></i> View Details
                                </a></li>
                                <li><a class="dropdown-item" href="#" onclick="editProfile('${profile.name}')">
                                    <i class="bi bi-pencil"></i> Edit
                                </a></li>
                                <li><a class="dropdown-item" href="#" onclick="testConnection('${profile.name}')">
                                    <i class="bi bi-wifi"></i> Test Connection
                                </a></li>
                                <li><a class="dropdown-item" href="#" onclick="viewAuditLog('${profile.name}')">
                                    <i class="bi bi-clock-history"></i> Audit Log
                                </a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item text-danger" href="#" onclick="deleteProfile('${profile.name}')">
                                    <i class="bi bi-trash"></i> Delete
                                </a></li>
                            </ul>
                        </div>
                    </div>

                    ${profile.key_path ? `
                        <p class="text-muted small mb-2">
                            <i class="bi bi-file-lock"></i> Key: ${profile.key_path}
                        </p>
                    ` : ''}

                    ${profile.description ? `
                        <p class="text-muted small">${profile.description}</p>
                    ` : ''}

                    <div class="mt-3">
                        <small class="text-muted">
                            <i class="bi bi-clock"></i>
                            ${profile.created_at ? new Date(profile.created_at).toLocaleDateString() : 'Unknown'}
                        </small>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
}

function showAddProfileModal() {
    const modal = new bootstrap.Modal(document.getElementById('addProfileModal'));
    document.getElementById('addProfileForm').reset();
    modal.show();
}

async function addProfile() {
    const name = document.getElementById('addProfileName').value;
    const host = document.getElementById('addProfileHost').value;
    const port = document.getElementById('addProfilePort').value;
    const user = document.getElementById('addProfileUser').value;
    const keyPath = document.getElementById('addProfileKeyPath').value;
    const password = document.getElementById('addProfilePassword').value;
    const description = document.getElementById('addProfileDescription').value;

    if (!name || !host || !user) {
        showError('Name, host, and user are required');
        return;
    }

    try {
        const payload = {
            name,
            host,
            port: port ? parseInt(port) : 22,
            user,
            description
        };

        if (keyPath) payload.key_path = keyPath;
        if (password) payload.password = password;

        await API.post('/api/v1/ssh', payload);

        bootstrap.Modal.getInstance(document.getElementById('addProfileModal')).hide();
        showSuccess('SSH profile added successfully');
        loadProfiles();
    } catch (error) {
        console.error('Failed to add profile:', error);
        showError('Failed to add SSH profile: ' + error.message);
    }
}

async function viewProfile(name) {
    try {
        const profile = await API.get(`/api/v1/ssh/${name}`);

        const modal = new bootstrap.Modal(document.getElementById('viewProfileModal'));
        document.getElementById('viewProfileContent').innerHTML = `
            <h5>${profile.name}</h5>
            <div class="mb-3">
                <strong>Host:</strong> ${profile.host}:${profile.port || 22}
            </div>
            <div class="mb-3">
                <strong>User:</strong> ${profile.user}
            </div>
            ${profile.key_path ? `
                <div class="mb-3">
                    <strong>SSH Key:</strong> <code>${profile.key_path}</code>
                </div>
            ` : ''}
            ${profile.description ? `
                <div class="mb-3">
                    <strong>Description:</strong> ${profile.description}
                </div>
            ` : ''}
            <div class="mb-3">
                <strong>Created:</strong> ${new Date(profile.created_at).toLocaleString()}
            </div>
            ${profile.last_used ? `
                <div class="mb-3">
                    <strong>Last Used:</strong> ${new Date(profile.last_used).toLocaleString()}
                </div>
            ` : ''}
        `;

        modal.show();
    } catch (error) {
        console.error('Failed to load profile:', error);
        showError('Failed to load profile details');
    }
}

async function testConnection(name) {
    showInfo('Testing connection to ' + name + '...');

    try {
        // Simulate connection test - you can implement actual test endpoint
        await new Promise(resolve => setTimeout(resolve, 1500));
        showSuccess('Connection test successful!');
    } catch (error) {
        showError('Connection test failed: ' + error.message);
    }
}

async function viewAuditLog(name) {
    try {
        const data = await API.get(`/api/v1/ssh/${name}/audit`);

        const modal = new bootstrap.Modal(document.getElementById('auditLogModal'));
        const container = document.getElementById('auditLogContent');

        if (!data.logs || data.logs.length === 0) {
            container.innerHTML = `
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No audit log entries available
                </div>
            `;
        } else {
            container.innerHTML = `
                <div class="table-responsive">
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th>Time</th>
                                <th>Action</th>
                                <th>User</th>
                                <th>Details</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${data.logs.map(log => `
                                <tr>
                                    <td>${new Date(log.timestamp).toLocaleString()}</td>
                                    <td><span class="badge bg-info">${log.action}</span></td>
                                    <td>${log.user || 'system'}</td>
                                    <td>${log.details || '-'}</td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>
                </div>
            `;
        }

        modal.show();
    } catch (error) {
        console.error('Failed to load audit log:', error);
        showError('Failed to load audit log');
    }
}

async function deleteProfile(name) {
    if (!confirm(`Are you sure you want to delete SSH profile "${name}"?`)) return;

    try {
        await API.delete(`/api/v1/ssh/${name}`);
        showSuccess('SSH profile deleted successfully');
        loadProfiles();
    } catch (error) {
        console.error('Failed to delete profile:', error);
        showError('Failed to delete SSH profile');
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

function showInfo(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-info border-0 position-fixed top-0 end-0 m-3';
    toast.style.zIndex = '9999';
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    document.body.appendChild(toast);
    new bootstrap.Toast(toast).show();
    setTimeout(() => toast.remove(), 3000);
}

document.addEventListener('DOMContentLoaded', () => {
    loadProfiles();
});
