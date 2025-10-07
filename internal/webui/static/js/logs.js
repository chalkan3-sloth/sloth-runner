// Logs Viewer
let currentLogs = [];
let currentLogFile = null;
let logStream = null;
let autoScroll = true;

const API = {
    async get(url) {
        const response = await fetch(url);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    }
};

async function loadLogFiles() {
    try {
        const data = await API.get('/api/v1/logs');
        currentLogs = data.logs || [];
        renderLogFiles();
    } catch (error) {
        console.error('Failed to load log files:', error);
        showError('Failed to load log files');
    }
}

function renderLogFiles() {
    const container = document.getElementById('log-files-list');

    if (currentLogs.length === 0) {
        container.innerHTML = `
            <div class="text-center py-3 text-muted">
                <i class="bi bi-file-earmark-text"></i>
                <p class="mb-0 mt-2">No log files found</p>
            </div>
        `;
        return;
    }

    container.innerHTML = currentLogs.map(log => `
        <button class="list-group-item list-group-item-action d-flex justify-content-between align-items-center ${log.name === currentLogFile ? 'active' : ''}"
                onclick="selectLogFile('${log.name}')">
            <div>
                <i class="bi bi-file-earmark-text"></i>
                <strong>${log.name}</strong>
                <br>
                <small class="text-muted">${formatBytes(log.size || 0)}</small>
            </div>
            ${log.name === currentLogFile ? '<i class="bi bi-chevron-right"></i>' : ''}
        </button>
    `).join('');
}

async function selectLogFile(filename) {
    currentLogFile = filename;
    renderLogFiles();
    await loadLogContent(filename);
}

async function loadLogContent(filename) {
    try {
        const data = await API.get(`/api/v1/logs/${filename}`);
        renderLogContent(data.content || '');
    } catch (error) {
        console.error('Failed to load log content:', error);
        showError('Failed to load log content');
    }
}

function renderLogContent(content) {
    const container = document.getElementById('log-content');

    if (!currentLogFile) {
        container.innerHTML = `
            <div class="text-center py-5">
                <i class="bi bi-file-earmark-text fs-1 text-muted"></i>
                <p class="text-muted mt-3">Select a log file to view</p>
            </div>
        `;
        return;
    }

    container.innerHTML = `
        <div class="d-flex justify-content-between align-items-center mb-3">
            <h5><i class="bi bi-file-earmark-text"></i> ${currentLogFile}</h5>
            <div class="btn-group btn-group-sm">
                <button class="btn btn-outline-primary" onclick="downloadLog('${currentLogFile}')" title="Download">
                    <i class="bi bi-download"></i> Download
                </button>
                <button class="btn btn-outline-secondary" onclick="toggleAutoScroll()" title="Auto-scroll" id="autoscroll-btn">
                    <i class="bi bi-arrow-down-circle${autoScroll ? '-fill' : ''}"></i>
                </button>
                <button class="btn btn-outline-danger" onclick="clearLogViewer()" title="Clear">
                    <i class="bi bi-x-circle"></i> Clear
                </button>
            </div>
        </div>
        <div class="log-viewer" id="log-viewer-content">
            <pre class="mb-0"><code>${escapeHtml(content)}</code></pre>
        </div>
    `;

    if (autoScroll) {
        scrollToBottom();
    }
}

function toggleAutoScroll() {
    autoScroll = !autoScroll;
    const btn = document.getElementById('autoscroll-btn');
    if (btn) {
        btn.innerHTML = `<i class="bi bi-arrow-down-circle${autoScroll ? '-fill' : ''}"></i>`;
    }
}

function clearLogViewer() {
    const viewer = document.getElementById('log-viewer-content');
    if (viewer) {
        viewer.innerHTML = '<pre class="mb-0"><code></code></pre>';
    }
}

function scrollToBottom() {
    const viewer = document.getElementById('log-viewer-content');
    if (viewer) {
        viewer.scrollTop = viewer.scrollHeight;
    }
}

function downloadLog(filename) {
    window.open(`/api/v1/logs/${filename}/download`, '_blank');
}

function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
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
    loadLogFiles();

    // Refresh log list every 10 seconds
    setInterval(loadLogFiles, 10000);
});
