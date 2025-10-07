// Events Management
let currentEvents = [];
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
    }
};

// Load events
async function loadEvents() {
    try {
        const endpoint = currentFilter === 'pending' ? '/api/v1/events/pending' : '/api/v1/events';
        const data = await API.get(endpoint);
        currentEvents = data.events || [];
        renderEvents();
        updateStats();
    } catch (error) {
        console.error('Failed to load events:', error);
        showError('Failed to load events');
    }
}

// Filter events
function filterEvents(filter) {
    currentFilter = filter;

    // Update filter buttons
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');

    loadEvents();
}

// Render events
function renderEvents() {
    const container = document.getElementById('events-table-body');

    if (currentEvents.length === 0) {
        container.innerHTML = `
            <tr>
                <td colspan="6" class="text-center py-4">
                    <div class="text-muted">
                        <i class="bi bi-inbox fs-1 d-block mb-2"></i>
                        No events found
                    </div>
                </td>
            </tr>
        `;
        return;
    }

    container.innerHTML = currentEvents.map(event => `
        <tr>
            <td><code class="small">${event.id || 'N/A'}</code></td>
            <td>
                <span class="badge bg-info">${event.type || 'unknown'}</span>
            </td>
            <td>
                <span class="badge ${getStatusBadgeClass(event.status)}">
                    ${event.status || 'unknown'}
                </span>
            </td>
            <td>${event.hook_id || '-'}</td>
            <td>
                <small class="text-muted">
                    ${event.created_at ? new Date(event.created_at).toLocaleString() : 'N/A'}
                </small>
            </td>
            <td>
                <small class="text-muted">
                    ${event.processed_at ? new Date(event.processed_at).toLocaleString() : '-'}
                </small>
            </td>
            <td>
                <div class="btn-group btn-group-sm">
                    <button class="btn btn-outline-primary" onclick="viewEvent('${event.id}')" title="View Details">
                        <i class="bi bi-eye"></i>
                    </button>
                    ${event.status === 'failed' ? `
                        <button class="btn btn-outline-warning" onclick="retryEvent('${event.id}')" title="Retry">
                            <i class="bi bi-arrow-clockwise"></i>
                        </button>
                    ` : ''}
                </div>
            </td>
        </tr>
    `).join('');
}

// Get status badge class
function getStatusBadgeClass(status) {
    const statusMap = {
        'pending': 'bg-warning',
        'processing': 'bg-info',
        'completed': 'bg-success',
        'failed': 'bg-danger',
        'cancelled': 'bg-secondary'
    };
    return statusMap[status] || 'bg-secondary';
}

// Update stats
function updateStats() {
    const pending = currentEvents.filter(e => e.status === 'pending').length;
    const processing = currentEvents.filter(e => e.status === 'processing').length;
    const completed = currentEvents.filter(e => e.status === 'completed').length;
    const failed = currentEvents.filter(e => e.status === 'failed').length;

    document.getElementById('stat-pending').textContent = pending;
    document.getElementById('stat-processing').textContent = processing;
    document.getElementById('stat-completed').textContent = completed;
    document.getElementById('stat-failed').textContent = failed;
}

// View event details
async function viewEvent(id) {
    try {
        const event = await API.get(`/api/v1/events/${id}`);

        const modal = new bootstrap.Modal(document.getElementById('viewEventModal'));
        document.getElementById('viewEventContent').innerHTML = `
            <div class="mb-3">
                <strong>ID:</strong> <code>${event.id}</code>
            </div>
            <div class="mb-3">
                <strong>Type:</strong> <span class="badge bg-info">${event.type}</span>
            </div>
            <div class="mb-3">
                <strong>Status:</strong> <span class="badge ${getStatusBadgeClass(event.status)}">${event.status}</span>
            </div>
            ${event.hook_id ? `
                <div class="mb-3">
                    <strong>Hook ID:</strong> <code>${event.hook_id}</code>
                </div>
            ` : ''}
            <div class="mb-3">
                <strong>Created:</strong> ${new Date(event.created_at).toLocaleString()}
            </div>
            ${event.processed_at ? `
                <div class="mb-3">
                    <strong>Processed:</strong> ${new Date(event.processed_at).toLocaleString()}
                </div>
            ` : ''}
            ${event.payload ? `
                <div class="mb-3">
                    <strong>Payload:</strong>
                    <pre class="bg-light p-3 rounded mt-2"><code>${JSON.stringify(event.payload, null, 2)}</code></pre>
                </div>
            ` : ''}
            ${event.error ? `
                <div class="mb-3">
                    <strong>Error:</strong>
                    <div class="alert alert-danger mt-2">${event.error}</div>
                </div>
            ` : ''}
        `;

        modal.show();
    } catch (error) {
        console.error('Failed to load event:', error);
        showError('Failed to load event details');
    }
}

// Retry event
async function retryEvent(id) {
    try {
        await API.post(`/api/v1/events/${id}/retry`);
        showSuccess('Event queued for retry');
        loadEvents();
    } catch (error) {
        console.error('Failed to retry event:', error);
        showError('Failed to retry event');
    }
}

// Utility functions
function showSuccess(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-success border-0 position-fixed top-0 end-0 m-3';
    toast.setAttribute('role', 'alert');
    toast.style.zIndex = '9999';
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
    toast.style.zIndex = '9999';
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

// Initialize and auto-refresh
document.addEventListener('DOMContentLoaded', () => {
    loadEvents();

    // Auto-refresh every 5 seconds
    setInterval(loadEvents, 5000);
});
