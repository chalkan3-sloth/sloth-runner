/* ===================================
   Micro-interactions & UX Enhancements
   =================================== */

// Prevent duplicate loading
if (typeof window.SLOTH_INTERACTIONS_LOADED !== 'undefined') {
    console.warn('Micro-interactions already loaded, skipping...');
} else {
    window.SLOTH_INTERACTIONS_LOADED = true;

// Ripple effect for buttons
class RippleEffect {
    constructor() {
        this.init();
    }

    init() {
        // Add ripple to all buttons and clickable elements
        document.addEventListener('click', (e) => {
            const element = e.target.closest('.btn, .card-clickable, .list-group-item-action');
            if (element && !element.classList.contains('no-ripple')) {
                this.createRipple(e, element);
            }
        });
    }

    createRipple(event, element) {
        const ripple = document.createElement('span');
        const rect = element.getBoundingClientRect();
        const size = Math.max(rect.width, rect.height);
        const x = event.clientX - rect.left - size / 2;
        const y = event.clientY - rect.top - size / 2;

        ripple.style.width = ripple.style.height = `${size}px`;
        ripple.style.left = `${x}px`;
        ripple.style.top = `${y}px`;
        ripple.classList.add('ripple-effect');

        // Remove old ripples
        const oldRipples = element.querySelectorAll('.ripple-effect');
        oldRipples.forEach(old => old.remove());

        // Add position relative if needed
        if (getComputedStyle(element).position === 'static') {
            element.style.position = 'relative';
            element.style.overflow = 'hidden';
        }

        element.appendChild(ripple);

        setTimeout(() => ripple.remove(), 600);
    }
}

// Smooth scroll to top
class ScrollToTop {
    constructor() {
        this.button = this.createButton();
        this.init();
    }

    createButton() {
        const btn = document.createElement('button');
        btn.className = 'scroll-to-top';
        btn.innerHTML = '<i class="bi bi-arrow-up"></i>';
        btn.setAttribute('aria-label', 'Scroll to top');
        document.body.appendChild(btn);
        return btn;
    }

    init() {
        window.addEventListener('scroll', () => {
            if (window.pageYOffset > 300) {
                this.button.classList.add('visible');
            } else {
                this.button.classList.remove('visible');
            }
        });

        this.button.addEventListener('click', () => {
            window.scrollTo({
                top: 0,
                behavior: 'smooth'
            });
        });
    }
}

// Loading state manager
class LoadingManager {
    static show(message = 'Loading...') {
        const existing = document.querySelector('.loading-overlay');
        if (existing) return;

        const overlay = document.createElement('div');
        overlay.className = 'loading-overlay';
        overlay.innerHTML = `
            <div class="spinner-wrapper">
                <div class="spinner-border text-primary" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
                <div class="loading-text">${message}</div>
            </div>
        `;
        document.body.appendChild(overlay);

        // Add progress bar
        const progressBar = document.createElement('div');
        progressBar.className = 'progress-bar-top loading';
        document.body.appendChild(progressBar);

        return {
            overlay,
            progressBar,
            remove: () => {
                overlay.remove();
                progressBar.remove();
            }
        };
    }

    static hide() {
        const overlay = document.querySelector('.loading-overlay');
        const progressBar = document.querySelector('.progress-bar-top');

        if (overlay) {
            overlay.style.animation = 'fadeOut 0.3s ease-out';
            setTimeout(() => overlay.remove(), 300);
        }

        if (progressBar) {
            setTimeout(() => progressBar.remove(), 300);
        }
    }

    static showProgress(element) {
        if (!element) return;
        element.classList.add('table-loading');
        return () => element.classList.remove('table-loading');
    }
}

// Keyboard shortcuts
class KeyboardShortcuts {
    constructor() {
        this.shortcuts = new Map();
        this.init();
    }

    init() {
        document.addEventListener('keydown', (e) => {
            // Don't trigger shortcuts when typing in inputs
            if (e.target.matches('input, textarea, select')) return;

            const key = this.getKeyCombo(e);
            const action = this.shortcuts.get(key);

            if (action) {
                e.preventDefault();
                action();
            }
        });

        // Register default shortcuts
        this.register('/', () => this.showShortcutsHelp());
        this.register('Escape', () => this.closeModals());
        this.register('Control+k', () => this.focusSearch());
        this.register('Meta+k', () => this.focusSearch()); // Mac
    }

    getKeyCombo(e) {
        const parts = [];
        if (e.ctrlKey || e.metaKey) parts.push(e.metaKey ? 'Meta' : 'Control');
        if (e.shiftKey) parts.push('Shift');
        if (e.altKey) parts.push('Alt');
        parts.push(e.key);
        return parts.join('+');
    }

    register(combo, action) {
        this.shortcuts.set(combo, action);
    }

    showShortcutsHelp() {
        const shortcuts = [
            { key: '/', description: 'Show keyboard shortcuts' },
            { key: 'Esc', description: 'Close modals/dialogs' },
            { key: 'Ctrl+K', description: 'Focus search' },
            { key: 'Ctrl+/', description: 'Toggle sidebar' }
        ];

        const html = `
            <div class="keyboard-shortcuts-help">
                <h5><i class="bi bi-keyboard"></i> Keyboard Shortcuts</h5>
                <div class="shortcuts-list">
                    ${shortcuts.map(s => `
                        <div class="shortcut-item">
                            <kbd>${s.key}</kbd>
                            <span>${s.description}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;

        if (window.toastManager) {
            toastManager.show({
                message: html,
                type: 'info',
                duration: 8000
            });
        }
    }

    closeModals() {
        // Close Bootstrap modals
        const modals = document.querySelectorAll('.modal.show');
        modals.forEach(modal => {
            const bsModal = bootstrap.Modal.getInstance(modal);
            if (bsModal) bsModal.hide();
        });

        // Close dropdowns
        const dropdowns = document.querySelectorAll('.dropdown-menu.show');
        dropdowns.forEach(dropdown => {
            dropdown.classList.remove('show');
        });
    }

    focusSearch() {
        const searchInput = document.querySelector('#navbar-search, input[type="search"]');
        if (searchInput) {
            searchInput.focus();
            searchInput.select();
        }
    }
}

// Smart tooltips
class SmartTooltips {
    constructor() {
        this.init();
    }

    init() {
        // Initialize Bootstrap tooltips with custom options
        const tooltipTriggerList = [].slice.call(
            document.querySelectorAll('[data-bs-toggle="tooltip"]')
        );
        tooltipTriggerList.map(el => new bootstrap.Tooltip(el, {
            delay: { show: 500, hide: 100 },
            animation: true
        }));

        // Auto-add tooltips to truncated text
        this.addTruncationTooltips();
    }

    addTruncationTooltips() {
        document.querySelectorAll('.text-truncate').forEach(el => {
            if (el.scrollWidth > el.clientWidth) {
                el.setAttribute('data-bs-toggle', 'tooltip');
                el.setAttribute('title', el.textContent);
                new bootstrap.Tooltip(el);
            }
        });
    }
}

// Copy to clipboard with feedback
class ClipboardManager {
    static async copy(text, successMessage = 'Copied!') {
        try {
            await navigator.clipboard.writeText(text);
            if (window.toastManager) {
                toastManager.success(successMessage, null, { duration: 2000 });
            }
            return true;
        } catch (error) {
            if (window.notify) {
                notify.error('Failed to copy to clipboard');
            }
            return false;
        }
    }

    static initCopyButtons() {
        document.querySelectorAll('[data-clipboard-copy]').forEach(btn => {
            btn.addEventListener('click', async () => {
                const text = btn.getAttribute('data-clipboard-copy') ||
                             btn.getAttribute('data-clipboard-target') &&
                             document.querySelector(btn.getAttribute('data-clipboard-target'))?.textContent;

                if (text) {
                    await this.copy(text);

                    // Visual feedback
                    const originalIcon = btn.innerHTML;
                    btn.innerHTML = '<i class="bi bi-check"></i>';
                    setTimeout(() => btn.innerHTML = originalIcon, 2000);
                }
            });
        });
    }
}

// Auto-refresh manager
class AutoRefreshManager {
    constructor() {
        this.intervals = new Map();
    }

    start(key, callback, intervalMs = 5000) {
        this.stop(key);
        const id = setInterval(callback, intervalMs);
        this.intervals.set(key, id);
    }

    stop(key) {
        const id = this.intervals.get(key);
        if (id) {
            clearInterval(id);
            this.intervals.delete(key);
        }
    }

    stopAll() {
        this.intervals.forEach(id => clearInterval(id));
        this.intervals.clear();
    }
}

// Empty state creator
class EmptyStateManager {
    static create(options = {}) {
        const {
            icon = 'inbox',
            title = 'No data found',
            subtitle = 'Get started by adding your first item',
            actionText = null,
            actionCallback = null
        } = options;

        const actionHtml = actionText ? `
            <button class="btn btn-primary mt-3 hover-grow" onclick="(${actionCallback})()">
                <i class="bi bi-plus-circle"></i> ${actionText}
            </button>
        ` : '';

        return `
            <div class="empty-state">
                <div class="empty-state-icon">
                    <i class="bi bi-${icon}"></i>
                </div>
                <h5 class="mt-3">${title}</h5>
                <p class="text-muted">${subtitle}</p>
                ${actionHtml}
            </div>
        `;
    }
}

// Initialize all micro-interactions
document.addEventListener('DOMContentLoaded', () => {
    new RippleEffect();
    new ScrollToTop();
    new KeyboardShortcuts();
    new SmartTooltips();
    ClipboardManager.initCopyButtons();
});

// Export for use in other scripts
window.RippleEffect = RippleEffect;
window.ScrollToTop = ScrollToTop;
window.LoadingManager = LoadingManager;
window.KeyboardShortcuts = KeyboardShortcuts;
window.SmartTooltips = SmartTooltips;
window.ClipboardManager = ClipboardManager;
window.AutoRefreshManager = AutoRefreshManager;
window.EmptyStateManager = EmptyStateManager;

// Create global instances
window.loadingManager = LoadingManager;
window.clipboardManager = ClipboardManager;
window.autoRefresh = new AutoRefreshManager();
window.emptyState = EmptyStateManager;

} // End of SLOTH_INTERACTIONS_LOADED check
