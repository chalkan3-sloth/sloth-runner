// Utility functions and notification system
// Prevent duplicate loading
if (typeof window.SLOTH_UTILS_LOADED !== 'undefined') {
    console.warn('Utils.js already loaded, skipping...');
} else {
    window.SLOTH_UTILS_LOADED = true;

// Theme Management
class ThemeManager {
    constructor() {
        this.theme = localStorage.getItem('theme') || 'light';
        this.apply();
    }

    toggle() {
        this.theme = this.theme === 'light' ? 'dark' : 'light';
        this.apply();
        localStorage.setItem('theme', this.theme);
    }

    apply() {
        document.documentElement.setAttribute('data-theme', this.theme);
    }

    isDark() {
        return this.theme === 'dark';
    }
}

const themeManager = new ThemeManager();

// Notification System
class NotificationManager {
    constructor() {
        this.container = null;
        this.init();
    }

    init() {
        // Create container if it doesn't exist
        if (!document.querySelector('.notification-container')) {
            this.container = document.createElement('div');
            this.container.className = 'notification-container';
            document.body.appendChild(this.container);
        } else {
            this.container = document.querySelector('.notification-container');
        }
    }

    show(message, type = 'info', duration = 5000) {
        const notification = document.createElement('div');
        notification.className = `notification alert alert-${type} alert-dismissible fade show`;
        notification.setAttribute('role', 'alert');

        const iconMap = {
            success: 'check-circle-fill',
            danger: 'exclamation-triangle-fill',
            warning: 'exclamation-circle-fill',
            info: 'info-circle-fill'
        };

        notification.innerHTML = `
            <i class="bi bi-${iconMap[type]} me-2"></i>
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;

        this.container.appendChild(notification);

        // Auto remove after duration
        if (duration > 0) {
            setTimeout(() => {
                notification.classList.add('removing');
                setTimeout(() => notification.remove(), 300);
            }, duration);
        }

        // Desktop notification if permitted
        if (Notification.permission === 'granted') {
            new Notification('Sloth Runner', {
                body: message,
                icon: '/static/img/logo.png'
            });
        }

        return notification;
    }

    success(message, duration) {
        return this.show(message, 'success', duration);
    }

    error(message, duration) {
        return this.show(message, 'danger', duration);
    }

    warning(message, duration) {
        return this.show(message, 'warning', duration);
    }

    info(message, duration) {
        return this.show(message, 'info', duration);
    }
}

const notify = new NotificationManager();

// Request Desktop Notification Permission
if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission();
}

// API Helper Functions
class API {
    static async get(endpoint) {
        try {
            const response = await fetch(`/api/v1${endpoint}`);
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            return await response.json();
        } catch (error) {
            notify.error(`API Error: ${error.message}`);
            throw error;
        }
    }

    static async post(endpoint, data) {
        try {
            const response = await fetch(`/api/v1${endpoint}`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(data)
            });
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            return await response.json();
        } catch (error) {
            notify.error(`API Error: ${error.message}`);
            throw error;
        }
    }

    static async put(endpoint, data) {
        try {
            const response = await fetch(`/api/v1${endpoint}`, {
                method: 'PUT',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(data)
            });
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            return await response.json();
        } catch (error) {
            notify.error(`API Error: ${error.message}`);
            throw error;
        }
    }

    static async delete(endpoint) {
        try {
            const response = await fetch(`/api/v1${endpoint}`, {
                method: 'DELETE'
            });
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            return await response.json();
        } catch (error) {
            notify.error(`API Error: ${error.message}`);
            throw error;
        }
    }
}

// Date/Time formatting
function formatDate(timestamp) {
    return new Date(timestamp * 1000).toLocaleString();
}

function formatDuration(ms) {
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);

    if (hours > 0) return `${hours}h ${minutes % 60}m`;
    if (minutes > 0) return `${minutes}m ${seconds % 60}s`;
    return `${seconds}s`;
}

function timeAgo(timestamp) {
    const seconds = Math.floor((Date.now() - timestamp * 1000) / 1000);

    if (seconds < 60) return `${seconds}s ago`;
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
    return `${Math.floor(seconds / 86400)}d ago`;
}

// File size formatting
function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Loading overlay
function showLoading(element, message = 'Loading...') {
    const overlay = document.createElement('div');
    overlay.className = 'text-center py-5';
    overlay.innerHTML = `
        <div class="spinner-border text-primary mb-3" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
        <p class="text-muted">${message}</p>
    `;
    element.innerHTML = '';
    element.appendChild(overlay);
}

// Error display
function showError(element, message) {
    element.innerHTML = `
        <div class="alert alert-danger" role="alert">
            <i class="bi bi-exclamation-triangle-fill me-2"></i>
            ${message}
        </div>
    `;
}

// Empty state
function showEmptyState(element, icon, title, subtitle) {
    element.innerHTML = `
        <div class="empty-state">
            <i class="bi bi-${icon}"></i>
            <h5>${title}</h5>
            <p>${subtitle}</p>
        </div>
    `;
}

// Copy to clipboard
async function copyToClipboard(text) {
    try {
        await navigator.clipboard.writeText(text);
        notify.success('Copied to clipboard!', 2000);
    } catch (error) {
        notify.error('Failed to copy');
    }
}

// Download file
function downloadFile(data, filename, type = 'text/plain') {
    const blob = new Blob([data], { type });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
}

// Confirm dialog
function confirmAction(message, onConfirm) {
    if (confirm(message)) {
        onConfirm();
    }
}

// Debounce function
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Export for use in other scripts
window.ThemeManager = ThemeManager;
window.themeManager = themeManager;
window.notify = notify;
window.API = API;
window.formatDate = formatDate;
window.formatDuration = formatDuration;
window.timeAgo = timeAgo;
window.formatBytes = formatBytes;
window.showLoading = showLoading;
window.showError = showError;
window.showEmptyState = showEmptyState;
window.copyToClipboard = copyToClipboard;
window.downloadFile = downloadFile;
window.confirmAction = confirmAction;
window.debounce = debounce;

} // End of SLOTH_UTILS_LOADED check
