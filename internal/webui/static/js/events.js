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
                <td colspan="7" class="text-center py-4">
                    <div class="text-muted">
                        <i class="bi bi-inbox fs-1 d-block mb-2"></i>
                        No events found
                    </div>
                </td>
            </tr>
        `;
        return;
    }

    container.innerHTML = currentEvents.map(event => {
        // Parse timestamps from Go time format
        const createdAt = event.CreatedAt || event.created_at;
        const processedAt = event.ProcessedAt || event.processed_at;

        // Format timestamp
        const formatTimestamp = (ts) => {
            if (!ts) return '-';
            try {
                return new Date(ts).toLocaleString();
            } catch {
                return ts;
            }
        };

        return `
        <tr>
            <td><code class="small">${(event.ID || event.id || 'N/A').substring(0, 8)}</code></td>
            <td>
                <span class="badge bg-info">${event.Type || event.type || 'unknown'}</span>
            </td>
            <td>
                <span class="badge ${getStatusBadgeClass(event.Status || event.status)}">
                    ${event.Status || event.status || 'unknown'}
                </span>
            </td>
            <td>-</td>
            <td>
                <small class="text-muted">
                    ${formatTimestamp(createdAt)}
                </small>
            </td>
            <td>
                <small class="text-muted">
                    ${formatTimestamp(processedAt)}
                </small>
            </td>
            <td>
                <div class="btn-group btn-group-sm">
                    <button class="btn btn-outline-primary" onclick="viewEvent('${event.ID || event.id}')" title="View Details">
                        <i class="bi bi-eye"></i>
                    </button>
                    ${(event.Status || event.status) === 'failed' ? `
                        <button class="btn btn-outline-warning" onclick="retryEvent('${event.ID || event.id}')" title="Retry">
                            <i class="bi bi-arrow-clockwise"></i>
                        </button>
                    ` : ''}
                </div>
            </td>
        </tr>
        `;
    }).join('');
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
    const pending = currentEvents.filter(e => (e.Status || e.status) === 'pending').length;
    const processing = currentEvents.filter(e => (e.Status || e.status) === 'processing').length;
    const completed = currentEvents.filter(e => (e.Status || e.status) === 'completed').length;
    const failed = currentEvents.filter(e => (e.Status || e.status) === 'failed').length;

    document.getElementById('stat-pending').textContent = pending;
    document.getElementById('stat-processing').textContent = processing;
    document.getElementById('stat-completed').textContent = completed;
    document.getElementById('stat-failed').textContent = failed;
}

// View event details
async function viewEvent(id) {
    try {
        const event = await API.get(`/api/v1/events/${id}`);

        // Support both Go struct (capitalized) and JSON (lowercase) field names
        const eventId = event.ID || event.id;
        const eventType = event.Type || event.type;
        const eventStatus = event.Status || event.status;
        const eventData = event.Data || event.data || event.payload;
        const eventError = event.Error || event.error;
        const createdAt = event.CreatedAt || event.created_at;
        const processedAt = event.ProcessedAt || event.processed_at;
        const timestamp = event.Timestamp || event.timestamp;

        const modal = new bootstrap.Modal(document.getElementById('viewEventModal'));
        document.getElementById('viewEventContent').innerHTML = `
            <div class="mb-3">
                <strong>ID:</strong> <code>${eventId}</code>
            </div>
            <div class="mb-3">
                <strong>Type:</strong> <span class="badge bg-info">${eventType}</span>
            </div>
            <div class="mb-3">
                <strong>Status:</strong> <span class="badge ${getStatusBadgeClass(eventStatus)}">${eventStatus}</span>
            </div>
            ${timestamp ? `
                <div class="mb-3">
                    <strong>Timestamp:</strong> ${new Date(timestamp).toLocaleString()}
                </div>
            ` : ''}
            <div class="mb-3">
                <strong>Created:</strong> ${new Date(createdAt).toLocaleString()}
            </div>
            ${processedAt ? `
                <div class="mb-3">
                    <strong>Processed:</strong> ${new Date(processedAt).toLocaleString()}
                </div>
            ` : ''}
            ${eventData ? `
                <div class="mb-3">
                    <strong>Data:</strong>
                    <pre class="bg-light p-3 rounded mt-2" style="max-height: 400px; overflow-y: auto;"><code>${JSON.stringify(eventData, null, 2)}</code></pre>
                </div>
            ` : ''}
            ${eventError ? `
                <div class="mb-3">
                    <strong>Error:</strong>
                    <div class="alert alert-danger mt-2">${eventError}</div>
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
