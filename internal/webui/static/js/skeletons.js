// Skeleton Loading Components Generator
// Provides ready-to-use skeleton templates for consistent loading states

const SkeletonLoader = {
    // Dashboard skeletons
    dashboardStats() {
        return `
            <div class="row mb-4">
                ${Array(4).fill(0).map(() => `
                    <div class="col-lg-3 col-md-6 mb-3">
                        <div class="skeleton-stat-card">
                            <div class="skeleton-stat-card-content">
                                <div class="skeleton-stat-card-text">
                                    <div class="skeleton skeleton-stat-label"></div>
                                    <div class="skeleton skeleton-stat-value"></div>
                                    <div class="skeleton skeleton-text-sm"></div>
                                </div>
                                <div class="skeleton skeleton-stat-icon"></div>
                            </div>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;
    },

    dashboardCharts() {
        return `
            <div class="row mb-4">
                <div class="col-lg-8">
                    <div class="card border-0 shadow-sm">
                        <div class="card-header bg-transparent border-0">
                            <div class="skeleton skeleton-title"></div>
                        </div>
                        <div class="card-body">
                            <div class="skeleton-chart">
                                <div class="skeleton-chart-bars">
                                    ${Array(6).fill(0).map(() => `
                                        <div class="skeleton skeleton-chart-bar"></div>
                                    `).join('')}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="col-lg-4">
                    <div class="card border-0 shadow-sm">
                        <div class="card-header bg-transparent border-0">
                            <div class="skeleton skeleton-title"></div>
                        </div>
                        <div class="card-body">
                            ${Array(3).fill(0).map(() => `
                                <div class="mb-4">
                                    <div class="skeleton skeleton-text mb-2"></div>
                                    <div class="progress" style="height: 10px;">
                                        <div class="skeleton" style="width: 100%; height: 100%;"></div>
                                    </div>
                                </div>
                            `).join('')}
                        </div>
                    </div>
                </div>
            </div>
        `;
    },

    activityFeed(count = 5) {
        return `
            ${Array(count).fill(0).map(() => `
                <div class="skeleton-activity-item">
                    <div class="skeleton skeleton-activity-icon"></div>
                    <div class="skeleton-activity-content">
                        <div class="skeleton skeleton-activity-message"></div>
                        <div class="skeleton skeleton-activity-time"></div>
                    </div>
                </div>
            `).join('')}
        `;
    },

    // Agent skeletons
    agentCards(count = 3) {
        return `
            <div class="row">
                ${Array(count).fill(0).map(() => `
                    <div class="col-lg-4 col-md-6 mb-4">
                        <div class="skeleton-agent-card">
                            <div class="skeleton-agent-header">
                                <div class="skeleton-agent-info">
                                    <div class="skeleton skeleton-agent-name"></div>
                                    <div class="skeleton skeleton-agent-address"></div>
                                </div>
                                <div class="skeleton skeleton-agent-status"></div>
                            </div>
                            <div class="skeleton-agent-metrics">
                                ${Array(3).fill(0).map(() => `
                                    <div class="skeleton-agent-metric">
                                        <div class="skeleton skeleton-agent-metric-label"></div>
                                        <div class="skeleton skeleton-agent-metric-value"></div>
                                    </div>
                                `).join('')}
                            </div>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;
    },

    agentList(count = 5) {
        return `
            <div class="skeleton-list">
                ${Array(count).fill(0).map(() => `
                    <div class="skeleton-list-item">
                        <div class="skeleton skeleton-avatar"></div>
                        <div style="flex: 1;">
                            <div class="skeleton skeleton-text" style="width: 60%; margin-bottom: 8px;"></div>
                            <div class="skeleton skeleton-text-sm" style="width: 40%;"></div>
                        </div>
                        <div class="skeleton skeleton-badge"></div>
                    </div>
                `).join('')}
            </div>
        `;
    },

    // Workflow skeletons
    workflowCards(count = 4) {
        return `
            <div class="row">
                ${Array(count).fill(0).map(() => `
                    <div class="col-lg-6 mb-4">
                        <div class="skeleton-workflow-card">
                            <div class="skeleton-workflow-header">
                                <div class="skeleton skeleton-workflow-title"></div>
                                <div class="skeleton-workflow-actions">
                                    ${Array(3).fill(0).map(() => `
                                        <div class="skeleton skeleton-workflow-action"></div>
                                    `).join('')}
                                </div>
                            </div>
                            <div class="skeleton-workflow-meta">
                                ${Array(3).fill(0).map(() => `
                                    <div class="skeleton-workflow-meta-item">
                                        <div class="skeleton skeleton-text-sm" style="width: 60px;"></div>
                                        <div class="skeleton skeleton-text" style="width: 80px;"></div>
                                    </div>
                                `).join('')}
                            </div>
                            <div class="skeleton skeleton-workflow-description"></div>
                            <div class="skeleton skeleton-workflow-description" style="width: 70%;"></div>
                            <div class="skeleton-workflow-tags mt-3">
                                ${Array(3).fill(0).map(() => `
                                    <div class="skeleton skeleton-workflow-tag"></div>
                                `).join('')}
                            </div>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;
    },

    // Table skeletons
    table(rows = 5, columns = 6) {
        return `
            <div class="table-responsive">
                <table class="table">
                    <thead>
                        <tr>
                            ${Array(columns).fill(0).map(() => `
                                <th><div class="skeleton skeleton-text" style="width: 80%;"></div></th>
                            `).join('')}
                        </tr>
                    </thead>
                    <tbody>
                        ${Array(rows).fill(0).map(() => `
                            <tr>
                                ${Array(columns).fill(0).map(() => `
                                    <td><div class="skeleton skeleton-text"></div></td>
                                `).join('')}
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            </div>
        `;
    },

    eventTable(rows = 10) {
        return `
            <div class="table-responsive">
                <table class="table table-hover">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Type</th>
                            <th>Status</th>
                            <th>Hook</th>
                            <th>Created</th>
                            <th>Processed</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${Array(rows).fill(0).map(() => `
                            <tr>
                                <td><div class="skeleton skeleton-text" style="width: 50px;"></div></td>
                                <td><div class="skeleton skeleton-text" style="width: 100px;"></div></td>
                                <td><div class="skeleton skeleton-badge"></div></td>
                                <td><div class="skeleton skeleton-text" style="width: 80px;"></div></td>
                                <td><div class="skeleton skeleton-text" style="width: 140px;"></div></td>
                                <td><div class="skeleton skeleton-text" style="width: 140px;"></div></td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            </div>
        `;
    },

    // Card skeletons
    card() {
        return `
            <div class="skeleton-card">
                <div class="skeleton-card-header">
                    <div class="skeleton skeleton-avatar"></div>
                    <div style="flex: 1;">
                        <div class="skeleton skeleton-text" style="width: 60%;"></div>
                        <div class="skeleton skeleton-text-sm" style="width: 40%;"></div>
                    </div>
                </div>
                <div class="skeleton-card-body">
                    <div class="skeleton skeleton-text"></div>
                    <div class="skeleton skeleton-text"></div>
                    <div class="skeleton skeleton-text" style="width: 80%;"></div>
                </div>
            </div>
        `;
    },

    // Generic content
    content() {
        return `
            <div class="skeleton skeleton-title"></div>
            ${Array(5).fill(0).map(() => `
                <div class="skeleton skeleton-text"></div>
            `).join('')}
            <div class="skeleton skeleton-text" style="width: 60%;"></div>
        `;
    },

    // Form skeleton
    form() {
        return `
            ${Array(3).fill(0).map(() => `
                <div class="mb-3">
                    <div class="skeleton skeleton-text" style="width: 120px; height: 14px; margin-bottom: 8px;"></div>
                    <div class="skeleton" style="height: 40px; border-radius: 8px;"></div>
                </div>
            `).join('')}
            <div class="skeleton skeleton-button"></div>
        `;
    },

    // Metrics/Stats
    metrics(count = 4) {
        return `
            <div class="row">
                ${Array(count).fill(0).map(() => `
                    <div class="col-md-${12/count} mb-3">
                        <div class="skeleton-card">
                            <div class="d-flex justify-content-between align-items-center">
                                <div style="flex: 1;">
                                    <div class="skeleton skeleton-text-sm" style="width: 100px; margin-bottom: 12px;"></div>
                                    <div class="skeleton" style="height: 32px; width: 60px; border-radius: 6px;"></div>
                                </div>
                                <div class="skeleton skeleton-icon" style="width: 48px; height: 48px; border-radius: 12px;"></div>
                            </div>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;
    },

    // Terminal skeleton
    terminal() {
        return `
            <div class="terminal-container" style="min-height: 400px;">
                ${Array(10).fill(0).map((_, i) => `
                    <div class="skeleton skeleton-text" style="width: ${Math.random() * 40 + 60}%; margin: 8px 0;"></div>
                `).join('')}
            </div>
        `;
    },

    // Modal content skeleton
    modal() {
        return `
            <div class="skeleton skeleton-title mb-4"></div>
            ${Array(4).fill(0).map(() => `
                <div class="mb-3">
                    <div class="skeleton skeleton-text"></div>
                    <div class="skeleton skeleton-text" style="width: 80%;"></div>
                </div>
            `).join('')}
        `;
    },

    // Helper function to show skeleton and replace with content
    async loadWithSkeleton(containerId, skeletonHtml, loadFunction) {
        const container = document.getElementById(containerId);
        if (!container) return;

        // Show skeleton
        container.innerHTML = skeletonHtml;

        try {
            // Load actual content
            await loadFunction();
        } catch (error) {
            // Show error state
            container.innerHTML = `
                <div class="empty-state">
                    <i class="bi bi-exclamation-triangle"></i>
                    <h5>Failed to load content</h5>
                    <p>${error.message}</p>
                    <button class="btn btn-primary" onclick="location.reload()">
                        <i class="bi bi-arrow-clockwise"></i> Retry
                    </button>
                </div>
            `;
        }
    },

    // Show skeleton in a specific container
    show(containerId, type = 'content', options = {}) {
        const container = document.getElementById(containerId);
        if (!container) return;

        const skeletonMap = {
            'dashboard-stats': () => this.dashboardStats(),
            'dashboard-charts': () => this.dashboardCharts(),
            'activity-feed': () => this.activityFeed(options.count || 5),
            'agent-cards': () => this.agentCards(options.count || 3),
            'agent-list': () => this.agentList(options.count || 5),
            'workflow-cards': () => this.workflowCards(options.count || 4),
            'table': () => this.table(options.rows || 5, options.columns || 6),
            'event-table': () => this.eventTable(options.rows || 10),
            'card': () => this.card(),
            'content': () => this.content(),
            'form': () => this.form(),
            'metrics': () => this.metrics(options.count || 4),
            'terminal': () => this.terminal(),
            'modal': () => this.modal()
        };

        const skeletonHtml = skeletonMap[type] ? skeletonMap[type]() : this.content();
        container.innerHTML = skeletonHtml;
    },

    // Hide skeleton and show content
    hide(containerId, content = '') {
        const container = document.getElementById(containerId);
        if (!container) return;

        if (content) {
            container.innerHTML = content;
        }
    }
};

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = SkeletonLoader;
}
