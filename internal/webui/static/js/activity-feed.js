/* ===================================
   Real-time Activity Feed
   Live updates of system events
   =================================== */

// Prevent duplicate loading
if (typeof window.SLOTH_ACTIVITY_FEED_LOADED !== 'undefined') {
    console.warn('Activity feed already loaded, skipping...');
} else {
    window.SLOTH_ACTIVITY_FEED_LOADED = true;

class ActivityFeed {
    constructor(container, options = {}) {
        this.container = typeof container === 'string'
            ? document.querySelector(container)
            : container;

        if (!this.container) {
            console.warn('Activity feed container not found');
            return;
        }

        this.options = {
            maxItems: options.maxItems || 50,
            autoScroll: options.autoScroll !== false,
            showTimestamps: options.showTimestamps !== false,
            filter: options.filter || null,
            onActivityClick: options.onActivityClick || null,
            types: options.types || ['info', 'success', 'warning', 'error'],
            ...options
        };

        this.activities = [];
        this.paused = false;

        this.init();
    }

    init() {
        this.setupContainer();
        this.loadActivities();

        // Start listening to events if WebSocket is available
        if (window.wsManager) {
            this.connectToWebSocket();
        }
    }

    setupContainer() {
        this.container.classList.add('activity-feed');
        this.container.innerHTML = `
            <div class="activity-feed-header">
                <h6><i class="bi bi-activity"></i> Atividade Recente</h6>
                <div class="activity-feed-controls">
                    <button class="btn btn-sm btn-link pause-btn" onclick="this.closest('.activity-feed').__instance.togglePause()" title="Pausar">
                        <i class="bi bi-pause-fill"></i>
                    </button>
                    <button class="btn btn-sm btn-link" onclick="this.closest('.activity-feed').__instance.clear()" title="Limpar">
                        <i class="bi bi-trash"></i>
                    </button>
                    <button class="btn btn-sm btn-link" onclick="this.closest('.activity-feed').__instance.refresh()" title="Atualizar">
                        <i class="bi bi-arrow-clockwise"></i>
                    </button>
                </div>
            </div>
            <div class="activity-feed-list" id="activity-feed-list"></div>
            <div class="activity-feed-footer">
                <small class="text-muted">Atualizando em tempo real...</small>
            </div>
        `;

        this.list = this.container.querySelector('.activity-feed-list');
        this.container.__instance = this;
    }

    async loadActivities() {
        try {
            const response = await fetch('/api/v1/events?limit=' + this.options.maxItems);
            if (response.ok) {
                const data = await response.json();
                this.activities = data.events || [];
                this.render();
            }
        } catch (error) {
            console.error('Failed to load activities:', error);
            this.addActivity({
                type: 'error',
                message: 'Falha ao carregar atividades',
                timestamp: Date.now() / 1000
            });
        }
    }

    connectToWebSocket() {
        // Subscribe to events
        wsManager.on('event', (event) => {
            if (!this.paused) {
                this.addActivity(event);
            }
        });

        wsManager.on('agent:status', (data) => {
            if (!this.paused) {
                this.addActivity({
                    type: data.status === 'connected' ? 'success' : 'warning',
                    icon: 'hdd-network',
                    message: `Agente ${data.agent} ${data.status === 'connected' ? 'conectado' : 'desconectado'}`,
                    timestamp: Date.now() / 1000,
                    link: `/agents#${data.agent}`
                });
            }
        });

        wsManager.on('workflow:status', (data) => {
            if (!this.paused) {
                const icons = {
                    running: 'play-circle',
                    completed: 'check-circle',
                    failed: 'x-circle'
                };

                const types = {
                    running: 'info',
                    completed: 'success',
                    failed: 'error'
                };

                this.addActivity({
                    type: types[data.status] || 'info',
                    icon: icons[data.status] || 'info-circle',
                    message: `Workflow "${data.workflow}" ${this.translateStatus(data.status)}`,
                    timestamp: Date.now() / 1000,
                    link: `/workflows#${data.workflow}`
                });
            }
        });
    }

    translateStatus(status) {
        const translations = {
            running: 'em execução',
            completed: 'concluído',
            failed: 'falhou',
            pending: 'pendente'
        };
        return translations[status] || status;
    }

    addActivity(activity) {
        // Check filter
        if (this.options.filter && !this.options.filter(activity)) {
            return;
        }

        // Add to beginning
        this.activities.unshift(activity);

        // Limit size
        if (this.activities.length > this.options.maxItems) {
            this.activities = this.activities.slice(0, this.options.maxItems);
        }

        // Render
        this.render();

        // Show notification for important events
        if (activity.type === 'error' && window.toastManager) {
            toastManager.error(activity.message, 'Erro', { duration: 5000 });
        }
    }

    render() {
        if (this.activities.length === 0) {
            this.list.innerHTML = `
                <div class="activity-feed-empty">
                    <i class="bi bi-inbox"></i>
                    <p>Nenhuma atividade recente</p>
                </div>
            `;
            return;
        }

        const html = this.activities.map((activity, index) => {
            const icon = activity.icon || this.getDefaultIcon(activity.type);
            const time = this.formatTime(activity.timestamp);

            return `
                <div class="activity-item activity-${activity.type}" data-index="${index}">
                    <div class="activity-icon activity-icon-${activity.type}">
                        <i class="bi bi-${icon}"></i>
                    </div>
                    <div class="activity-content">
                        <div class="activity-message">${activity.message}</div>
                        ${this.options.showTimestamps ? `
                            <div class="activity-time">${time}</div>
                        ` : ''}
                    </div>
                    ${activity.link ? `
                        <a href="${activity.link}" class="activity-link">
                            <i class="bi bi-arrow-right"></i>
                        </a>
                    ` : ''}
                </div>
            `;
        }).join('');

        this.list.innerHTML = html;

        // Attach click handlers
        this.list.querySelectorAll('.activity-item').forEach(item => {
            const index = parseInt(item.dataset.index);
            item.addEventListener('click', () => {
                if (this.options.onActivityClick) {
                    this.options.onActivityClick(this.activities[index]);
                }
            });
        });

        // Auto scroll
        if (this.options.autoScroll && !this.paused) {
            this.list.scrollTop = 0;
        }
    }

    getDefaultIcon(type) {
        const icons = {
            success: 'check-circle',
            error: 'x-circle',
            warning: 'exclamation-triangle',
            info: 'info-circle'
        };
        return icons[type] || 'circle';
    }

    formatTime(timestamp) {
        const now = Date.now() / 1000;
        const diff = now - timestamp;

        if (diff < 60) return 'agora';
        if (diff < 3600) return `${Math.floor(diff / 60)}m atrás`;
        if (diff < 86400) return `${Math.floor(diff / 3600)}h atrás`;
        if (diff < 604800) return `${Math.floor(diff / 86400)}d atrás`;

        return new Date(timestamp * 1000).toLocaleDateString();
    }

    togglePause() {
        this.paused = !this.paused;

        const pauseBtn = this.container.querySelector('.pause-btn');
        const icon = pauseBtn.querySelector('i');

        if (this.paused) {
            icon.className = 'bi bi-play-fill';
            pauseBtn.title = 'Retomar';

            const footer = this.container.querySelector('.activity-feed-footer small');
            footer.textContent = 'Pausado';
            footer.style.color = '#f59e0b';
        } else {
            icon.className = 'bi bi-pause-fill';
            pauseBtn.title = 'Pausar';

            const footer = this.container.querySelector('.activity-feed-footer small');
            footer.textContent = 'Atualizando em tempo real...';
            footer.style.color = '';
        }
    }

    clear() {
        if (confirm('Limpar todas as atividades?')) {
            this.activities = [];
            this.render();

            if (window.toastManager) {
                toastManager.success('Atividades limpas', null, { duration: 2000 });
            }
        }
    }

    refresh() {
        this.loadActivities();

        if (window.toastManager) {
            toastManager.success('Atividades atualizadas', null, { duration: 2000 });
        }
    }
}

// Export
window.ActivityFeed = ActivityFeed;

// Helper function to create activity feed
window.createActivityFeed = function(selector, options) {
    return new ActivityFeed(selector, options);
};

} // End of SLOTH_ACTIVITY_FEED_LOADED check
