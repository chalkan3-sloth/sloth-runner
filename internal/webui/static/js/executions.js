// Workflow Executions Management
let currentExecutions = [];
let currentExecutionId = null;
let autoScroll = true;
let refreshInterval = null;

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

async function loadExecutions() {
    try {
        const data = await API.get('/api/v1/executions');
        currentExecutions = data.executions || [];
        renderExecutions();
    } catch (error) {
        console.error('Failed to load executions:', error);
        showError('Failed to load executions');
    }
}

function renderExecutions() {
    const container = document.getElementById('executionsList');

    if (currentExecutions.length === 0) {
        container.innerHTML = `
            <div class="text-center py-4 text-muted">
                <i class="bi bi-play-circle"></i>
                <p class="mb-0 mt-2">No executions yet</p>
            </div>
        `;
        return;
    }

    container.innerHTML = currentExecutions.map(exec => {
        const statusClass = getStatusClass(exec.status);
        const statusIcon = getStatusIcon(exec.status);
        const duration = exec.end_time ?
            formatDuration(new Date(exec.end_time) - new Date(exec.start_time)) :
            'Running...';

        return `
            <button class="list-group-item list-group-item-action ${exec.id === currentExecutionId ? 'active' : ''}"
                    onclick="viewExecution('${exec.id}')">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <i class="bi bi-${statusIcon} text-${statusClass}"></i>
                        <strong>${exec.workflow_id || 'Unknown'}</strong>
                        <br>
                        <small class="text-muted">ID: ${exec.id.substring(0, 8)}</small>
                    </div>
                    <div class="text-end">
                        <span class="badge bg-${statusClass}">${exec.status}</span>
                        <br>
                        <small class="text-muted">${duration}</small>
                    </div>
                </div>
            </button>
        `;
    }).join('');
}

async function viewExecution(id) {
    currentExecutionId = id;
    renderExecutions();

    try {
        const data = await API.get(`/api/v1/executions/${id}`);

        document.getElementById('executionInfo').style.display = 'block';
        document.getElementById('infoWorkflow').textContent = data.workflow_id || 'Unknown';
        document.getElementById('infoStatus').innerHTML = getStatusBadge(data.status);
        document.getElementById('infoStartTime').textContent = formatTimestamp(data.start_time);

        const duration = data.end_time ?
            formatDuration(new Date(data.end_time) - new Date(data.start_time)) :
            'Running...';
        document.getElementById('infoDuration').textContent = duration;

        // Render logs
        const logsContainer = document.getElementById('logsContainer');
        if (data.logs && data.logs.length > 0) {
            logsContainer.innerHTML = data.logs.map(log =>
                `<div class="log-line">${escapeHtml(log)}</div>`
            ).join('');
            document.getElementById('logCount').textContent = `${data.logs.length} lines`;
        } else {
            logsContainer.innerHTML = '<div class="text-muted text-center py-5">No logs available</div>';
            document.getElementById('logCount').textContent = '0 lines';
        }

        // Show cancel button if running
        document.getElementById('executionControls').style.display =
            data.status === 'running' ? 'block' : 'none';

        if (autoScroll) {
            logsContainer.scrollTop = logsContainer.scrollHeight;
        }

    } catch (error) {
        console.error('Failed to load execution:', error);
        showError('Failed to load execution details');
    }
}

async function cancelExecution() {
    if (!currentExecutionId) return;

    if (!confirm('Are you sure you want to cancel this execution?')) return;

    try {
        await API.post(`/api/v1/executions/${currentExecutionId}/cancel`);
        showSuccess('Execution cancelled');
        await viewExecution(currentExecutionId);
        await loadExecutions();
    } catch (error) {
        console.error('Failed to cancel execution:', error);
        showError('Failed to cancel execution');
    }
}

function clearLogs() {
    document.getElementById('logsContainer').innerHTML = `
        <div class="text-muted text-center py-5">
            <i class="bi bi-info-circle fs-1"></i>
            <p class="mt-3">Logs cleared. Select an execution to view logs.</p>
        </div>
    `;
    document.getElementById('logCount').textContent = '0 lines';
}

function downloadLogs() {
    if (!currentExecutionId) {
        showError('No execution selected');
        return;
    }

    const logsContainer = document.getElementById('logsContainer');
    const logs = Array.from(logsContainer.querySelectorAll('.log-line'))
        .map(line => line.textContent)
        .join('\n');

    if (!logs) {
        showError('No logs to download');
        return;
    }

    const blob = new Blob([logs], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `execution-${currentExecutionId}.log`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

function toggleAutoScroll() {
    autoScroll = !autoScroll;
    const btn = document.getElementById('autoScrollBtn');
    btn.innerHTML = `<i class="bi bi-arrow-down-circle${autoScroll ? '-fill' : ''}"></i> Auto-scroll: ${autoScroll ? 'ON' : 'OFF'}`;
}

function getStatusClass(status) {
    switch (status) {
        case 'running': return 'primary';
        case 'completed': case 'success': return 'success';
        case 'failed': case 'error': return 'danger';
        case 'cancelled': return 'warning';
        default: return 'secondary';
    }
}

function getStatusIcon(status) {
    switch (status) {
        case 'running': return 'arrow-repeat';
        case 'completed': case 'success': return 'check-circle';
        case 'failed': case 'error': return 'x-circle';
        case 'cancelled': return 'dash-circle';
        default: return 'circle';
    }
}

function getStatusBadge(status) {
    const statusClass = getStatusClass(status);
    const statusIcon = getStatusIcon(status);
    return `<span class="badge bg-${statusClass}"><i class="bi bi-${statusIcon}"></i> ${status}</span>`;
}

function formatTimestamp(timestamp) {
    if (!timestamp) return 'N/A';
    const date = new Date(timestamp);
    return date.toLocaleString();
}

function formatDuration(ms) {
    if (!ms || ms < 0) return '0s';

    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);

    if (hours > 0) {
        return `${hours}h ${minutes % 60}m ${seconds % 60}s`;
    } else if (minutes > 0) {
        return `${minutes}m ${seconds % 60}s`;
    }
    return `${seconds}s`;
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

// Add CSS for log lines
const style = document.createElement('style');
style.textContent = `
    .log-viewer {
        background: #1a1a1a;
        color: #e9ecef;
        padding: 1rem;
        height: 600px;
        overflow-y: auto;
        font-family: 'Courier New', monospace;
        font-size: 13px;
    }
    .log-line {
        margin-bottom: 2px;
        line-height: 1.4;
    }
`;
document.head.appendChild(style);

document.addEventListener('DOMContentLoaded', () => {
    loadExecutions();

    // Setup cancel button
    const cancelBtn = document.getElementById('cancelBtn');
    if (cancelBtn) {
        cancelBtn.addEventListener('click', cancelExecution);
    }

    // Refresh every 5 seconds
    refreshInterval = setInterval(() => {
        loadExecutions();
        if (currentExecutionId) {
            viewExecution(currentExecutionId);
        }
    }, 5000);
});
