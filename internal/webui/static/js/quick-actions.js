/* ===================================
   Quick Actions Floating Menu
   Contextual actions based on page
   =================================== */

// Prevent duplicate loading
if (typeof window.SLOTH_QUICK_ACTIONS_LOADED !== 'undefined') {
    console.warn('Quick actions already loaded, skipping...');
} else {
    window.SLOTH_QUICK_ACTIONS_LOADED = true;

class QuickActionsMenu {
    constructor(options = {}) {
        this.options = {
            position: options.position || 'bottom-right',
            actions: options.actions || [],
            showLabels: options.showLabels !== false,
            ...options
        };

        this.isOpen = false;
        this.actions = new Map();

        this.init();
        this.registerDefaultActions();
    }

    init() {
        this.createMenu();
        this.attachEventListeners();
    }

    createMenu() {
        const menu = document.createElement('div');
        menu.id = 'quick-actions-menu';
        menu.className = `quick-actions-menu quick-actions-${this.options.position}`;
        menu.innerHTML = `
            <button class="quick-actions-trigger" title="Ações Rápidas">
                <i class="bi bi-lightning-fill"></i>
            </button>
            <div class="quick-actions-list"></div>
        `;

        document.body.appendChild(menu);

        this.menu = menu;
        this.trigger = menu.querySelector('.quick-actions-trigger');
        this.list = menu.querySelector('.quick-actions-list');
    }

    attachEventListeners() {
        this.trigger.addEventListener('click', () => this.toggle());

        // Close on click outside
        document.addEventListener('click', (e) => {
            if (!this.menu.contains(e.target) && this.isOpen) {
                this.close();
            }
        });

        // Close on escape
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isOpen) {
                this.close();
            }
        });
    }

    registerDefaultActions() {
        // Page-specific actions
        const page = window.location.pathname.substring(1) || 'dashboard';

        // Global actions
        this.register({
            id: 'theme-toggle',
            label: 'Alternar Tema',
            icon: 'moon-stars',
            color: '#7C3AED',
            action: () => {
                if (window.themeManager) {
                    themeManager.toggle();
                }
            }
        });

        this.register({
            id: 'command-palette',
            label: 'Command Palette',
            icon: 'terminal',
            color: '#4F46E5',
            shortcut: 'Ctrl+Shift+P',
            action: () => {
                if (window.commandPalette) {
                    commandPalette.open();
                }
            }
        });

        this.register({
            id: 'refresh',
            label: 'Atualizar Página',
            icon: 'arrow-clockwise',
            color: '#10b981',
            action: () => window.location.reload()
        });

        // Page-specific actions
        if (page === 'agents') {
            this.register({
                id: 'add-agent',
                label: 'Adicionar Agente',
                icon: 'plus-circle',
                color: '#EC4899',
                action: () => {
                    // Trigger add agent modal
                    if (window.toastManager) {
                        toastManager.info('Abrindo modal de adicionar agente...');
                    }
                }
            });
        }

        if (page === 'workflows') {
            this.register({
                id: 'new-workflow',
                label: 'Novo Workflow',
                icon: 'file-plus',
                color: '#EC4899',
                action: () => {
                    if (window.toastManager) {
                        toastManager.info('Criando novo workflow...');
                    }
                }
            });
        }

        if (page === 'terminal') {
            this.register({
                id: 'new-session',
                label: 'Nova Sessão',
                icon: 'terminal-plus',
                color: '#EC4899',
                action: () => {
                    if (window.toastManager) {
                        toastManager.info('Abrindo nova sessão de terminal...');
                    }
                }
            });
        }

        this.render();
    }

    register(action) {
        this.actions.set(action.id, action);
        this.render();
    }

    unregister(actionId) {
        this.actions.delete(actionId);
        this.render();
    }

    render() {
        const actionsArray = Array.from(this.actions.values());

        if (actionsArray.length === 0) {
            this.menu.style.display = 'none';
            return;
        }

        this.menu.style.display = 'flex';

        this.list.innerHTML = actionsArray.map(action => `
            <button class="quick-action-item"
                    data-action-id="${action.id}"
                    title="${action.label}${action.shortcut ? ` (${action.shortcut})` : ''}"
                    style="--action-color: ${action.color || '#4F46E5'}">
                <i class="bi bi-${action.icon}"></i>
                ${this.options.showLabels ? `<span>${action.label}</span>` : ''}
            </button>
        `).join('');

        // Attach click handlers
        this.list.querySelectorAll('.quick-action-item').forEach(btn => {
            btn.addEventListener('click', () => {
                const actionId = btn.dataset.actionId;
                const action = this.actions.get(actionId);
                if (action) {
                    action.action();
                    this.close();
                }
            });
        });
    }

    toggle() {
        if (this.isOpen) {
            this.close();
        } else {
            this.open();
        }
    }

    open() {
        this.isOpen = true;
        this.menu.classList.add('active');
        this.trigger.querySelector('i').style.transform = 'rotate(45deg)';
    }

    close() {
        this.isOpen = false;
        this.menu.classList.remove('active');
        this.trigger.querySelector('i').style.transform = 'rotate(0deg)';
    }
}

// Initialize
const quickActions = new QuickActionsMenu();
window.quickActions = quickActions;

} // End of SLOTH_QUICK_ACTIONS_LOADED check
