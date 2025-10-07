// Modern Toast Notification System for Sloth Runner
// Supports multiple types, auto-dismiss, progress bars, actions, and stacking

class ToastManager {
    constructor() {
        this.container = null;
        this.toasts = [];
        this.maxToasts = 5;
        this.defaultDuration = 5000;
        this.init();
    }

    init() {
        // Create toast container if it doesn't exist
        if (!this.container) {
            this.container = document.createElement('div');
            this.container.id = 'toast-container';
            this.container.className = 'toast-container';
            document.body.appendChild(this.container);
        }
    }

    /**
     * Show a toast notification
     * @param {Object} options - Toast configuration
     * @param {string} options.message - Toast message
     * @param {string} options.type - Toast type: success, error, warning, info, loading
     * @param {string} options.title - Optional title
     * @param {number} options.duration - Duration in ms (0 for persistent)
     * @param {Array} options.actions - Optional action buttons
     * @param {boolean} options.dismissible - Can be dismissed manually
     * @param {Function} options.onDismiss - Callback when dismissed
     */
    show(options) {
        const {
            message,
            type = 'info',
            title = '',
            duration = this.defaultDuration,
            actions = [],
            dismissible = true,
            onDismiss = null,
            icon = null
        } = options;

        // Remove oldest toast if max limit reached
        if (this.toasts.length >= this.maxToasts) {
            this.dismiss(this.toasts[0].id);
        }

        const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
        const toast = this.createToast({ id, message, type, title, duration, actions, dismissible, icon, onDismiss });

        this.toasts.push({ id, element: toast, onDismiss });
        this.container.appendChild(toast);

        // Trigger animation
        requestAnimationFrame(() => {
            toast.classList.add('toast-show');
        });

        // Auto-dismiss after duration
        if (duration > 0) {
            setTimeout(() => {
                this.dismiss(id);
            }, duration);
        }

        return id;
    }

    createToast({ id, message, type, title, duration, actions, dismissible, icon, onDismiss }) {
        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.setAttribute('data-toast-id', id);
        toast.setAttribute('role', 'alert');

        // Get icon based on type
        const toastIcon = icon || this.getIcon(type);

        // Build toast content
        let html = `
            <div class="toast-content">
                <div class="toast-icon">
                    <i class="bi ${toastIcon}"></i>
                </div>
                <div class="toast-body">
                    ${title ? `<div class="toast-title">${title}</div>` : ''}
                    <div class="toast-message">${message}</div>
                    ${actions.length > 0 ? this.buildActions(actions, id) : ''}
                </div>
                ${dismissible ? `
                    <button class="toast-close" onclick="toastManager.dismiss('${id}')">
                        <i class="bi bi-x-lg"></i>
                    </button>
                ` : ''}
            </div>
        `;

        // Add progress bar for auto-dismiss toasts
        if (duration > 0) {
            html += `
                <div class="toast-progress">
                    <div class="toast-progress-bar" style="animation-duration: ${duration}ms"></div>
                </div>
            `;
        }

        toast.innerHTML = html;
        return toast;
    }

    buildActions(actions, toastId) {
        const actionsHtml = actions.map(action => `
            <button class="toast-action-btn" onclick="(function() {
                ${action.onClick ? `(${action.onClick.toString()})();` : ''}
                toastManager.dismiss('${toastId}');
            })()">
                ${action.icon ? `<i class="bi ${action.icon}"></i>` : ''}
                ${action.label}
            </button>
        `).join('');

        return `<div class="toast-actions">${actionsHtml}</div>`;
    }

    getIcon(type) {
        const icons = {
            success: 'bi-check-circle-fill',
            error: 'bi-x-circle-fill',
            warning: 'bi-exclamation-triangle-fill',
            info: 'bi-info-circle-fill',
            loading: 'bi-arrow-repeat'
        };
        return icons[type] || icons.info;
    }

    dismiss(id) {
        const index = this.toasts.findIndex(t => t.id === id);
        if (index === -1) return;

        const { element, onDismiss } = this.toasts[index];

        // Trigger exit animation
        element.classList.remove('toast-show');
        element.classList.add('toast-hide');

        // Remove from DOM after animation
        setTimeout(() => {
            if (element.parentNode) {
                element.parentNode.removeChild(element);
            }
            if (onDismiss) onDismiss();
        }, 300);

        // Remove from array
        this.toasts.splice(index, 1);
    }

    dismissAll() {
        [...this.toasts].forEach(({ id }) => this.dismiss(id));
    }

    // Helper methods for common toast types
    success(message, title = '', options = {}) {
        return this.show({ message, title, type: 'success', ...options });
    }

    error(message, title = '', options = {}) {
        return this.show({ message, title, type: 'error', duration: 7000, ...options });
    }

    warning(message, title = '', options = {}) {
        return this.show({ message, title, type: 'warning', ...options });
    }

    info(message, title = '', options = {}) {
        return this.show({ message, title, type: 'info', ...options });
    }

    loading(message, title = '', options = {}) {
        return this.show({
            message,
            title,
            type: 'loading',
            duration: 0,
            dismissible: false,
            ...options
        });
    }

    // Update an existing toast (useful for loading states)
    update(id, options) {
        const index = this.toasts.findIndex(t => t.id === id);
        if (index === -1) return;

        const oldToast = this.toasts[index];

        // Create new toast with updated options
        const newToast = this.createToast({
            id,
            ...options,
            duration: options.duration !== undefined ? options.duration : this.defaultDuration
        });

        // Replace in DOM
        oldToast.element.parentNode.replaceChild(newToast, oldToast.element);
        this.toasts[index].element = newToast;

        // Trigger animation
        requestAnimationFrame(() => {
            newToast.classList.add('toast-show');
        });

        // Auto-dismiss if duration is set
        if (options.duration && options.duration > 0) {
            setTimeout(() => {
                this.dismiss(id);
            }, options.duration);
        }
    }

    // Promise-based loading toast
    async promise(promise, messages = {}) {
        const {
            loading = 'Loading...',
            success = 'Success!',
            error = 'Error occurred'
        } = messages;

        const id = this.loading(loading);

        try {
            const result = await promise;
            this.update(id, {
                type: 'success',
                message: success,
                duration: 3000,
                dismissible: true
            });
            return result;
        } catch (err) {
            this.update(id, {
                type: 'error',
                message: typeof error === 'function' ? error(err) : error,
                duration: 7000,
                dismissible: true
            });
            throw err;
        }
    }
}

// Create global instance
const toastManager = new ToastManager();

// Export as global
window.toastManager = toastManager;
window.ToastManager = ToastManager;

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ToastManager;
}
