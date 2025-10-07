// Breadcrumbs Navigation Component
// Provides hierarchical navigation with auto-generation from path

class BreadcrumbManager {
    constructor(containerId = 'breadcrumbs') {
        this.containerId = containerId;
        this.separator = '/';
        this.homeLabel = 'Home';
        this.homeIcon = 'bi-house-door';
        this.routes = this.initializeRoutes();
    }

    /**
     * Initialize route mappings
     * Maps URL paths to readable names and icons
     */
    initializeRoutes() {
        return {
            '/': { label: 'Dashboard', icon: 'bi-speedometer2' },
            '/agents': { label: 'Agents', icon: 'bi-hdd-network' },
            '/workflows': { label: 'Workflows', icon: 'bi-diagram-3' },
            '/stacks': { label: 'Stacks', icon: 'bi-layers' },
            '/hooks': { label: 'Hooks', icon: 'bi-lightning' },
            '/events': { label: 'Events', icon: 'bi-bell' },
            '/scheduler': { label: 'Scheduler', icon: 'bi-calendar-event' },
            '/executions': { label: 'Executions', icon: 'bi-play-circle' },
            '/ssh': { label: 'SSH Profiles', icon: 'bi-key' },
            '/secrets': { label: 'Secrets', icon: 'bi-shield-lock' },
            '/backup': { label: 'Backup', icon: 'bi-cloud-upload' },
            '/terminal': { label: 'Terminal', icon: 'bi-terminal' },
            '/metrics': { label: 'Metrics', icon: 'bi-graph-up' },
            '/logs': { label: 'Logs', icon: 'bi-file-text' }
        };
    }

    /**
     * Generate breadcrumbs from current path
     */
    generate() {
        const path = window.location.pathname;
        const container = document.getElementById(this.containerId);

        if (!container) {
            console.warn(`Breadcrumb container #${this.containerId} not found`);
            return;
        }

        const breadcrumbs = this.buildBreadcrumbs(path);
        container.innerHTML = this.render(breadcrumbs);
        this.attachAnimations();
    }

    /**
     * Build breadcrumb items from path
     */
    buildBreadcrumbs(path) {
        const breadcrumbs = [];

        // Always start with home
        breadcrumbs.push({
            label: this.homeLabel,
            icon: this.homeIcon,
            url: '/',
            isActive: path === '/'
        });

        // Parse path segments
        if (path !== '/') {
            const segments = path.split('/').filter(s => s);
            let currentPath = '';

            segments.forEach((segment, index) => {
                currentPath += `/${segment}`;
                const isLast = index === segments.length - 1;

                // Get route info or use segment name
                const routeInfo = this.routes[currentPath] || {
                    label: this.formatSegment(segment),
                    icon: null
                };

                breadcrumbs.push({
                    label: routeInfo.label,
                    icon: routeInfo.icon,
                    url: currentPath,
                    isActive: isLast
                });
            });
        }

        return breadcrumbs;
    }

    /**
     * Format segment name (convert kebab-case to Title Case)
     */
    formatSegment(segment) {
        return segment
            .split('-')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ');
    }

    /**
     * Render breadcrumbs HTML
     */
    render(breadcrumbs) {
        const items = breadcrumbs.map((crumb, index) => {
            const isLast = index === breadcrumbs.length - 1;

            return `
                <li class="breadcrumb-item ${isLast ? 'active' : ''} fade-in-up stagger-${index + 1}">
                    ${!isLast ? `
                        <a href="${crumb.url}" class="breadcrumb-link">
                            ${crumb.icon ? `<i class="bi ${crumb.icon}"></i>` : ''}
                            <span>${crumb.label}</span>
                        </a>
                    ` : `
                        <span class="breadcrumb-current">
                            ${crumb.icon ? `<i class="bi ${crumb.icon}"></i>` : ''}
                            <span>${crumb.label}</span>
                        </span>
                    `}
                </li>
            `;
        }).join('');

        return `
            <nav aria-label="breadcrumb">
                <ol class="breadcrumb-list">
                    ${items}
                </ol>
            </nav>
        `;
    }

    /**
     * Set custom breadcrumbs manually
     */
    set(breadcrumbs) {
        const container = document.getElementById(this.containerId);
        if (!container) return;

        container.innerHTML = this.render(breadcrumbs);
        this.attachAnimations();
    }

    /**
     * Add a breadcrumb item dynamically
     */
    add(label, url, icon = null) {
        const container = document.getElementById(this.containerId);
        if (!container) return;

        const breadcrumbList = container.querySelector('.breadcrumb-list');
        if (!breadcrumbList) return;

        // Remove active class from previous last item
        const previousLast = breadcrumbList.querySelector('.breadcrumb-item.active');
        if (previousLast) {
            previousLast.classList.remove('active');
            const link = previousLast.querySelector('.breadcrumb-current');
            if (link) {
                const href = link.previousElementSibling?.href || '#';
                link.outerHTML = `
                    <a href="${href}" class="breadcrumb-link">
                        ${icon ? `<i class="bi ${icon}"></i>` : ''}
                        <span>${link.textContent}</span>
                    </a>
                `;
            }
        }

        // Add new item
        const newItem = document.createElement('li');
        newItem.className = 'breadcrumb-item active';
        newItem.innerHTML = `
            <span class="breadcrumb-current">
                ${icon ? `<i class="bi ${icon}"></i>` : ''}
                <span>${label}</span>
            </span>
        `;

        breadcrumbList.appendChild(newItem);
    }

    /**
     * Attach entrance animations
     */
    attachAnimations() {
        const items = document.querySelectorAll('.breadcrumb-item');
        items.forEach((item, index) => {
            item.style.animationDelay = `${index * 0.1}s`;
        });
    }

    /**
     * Register a custom route
     */
    registerRoute(path, label, icon = null) {
        this.routes[path] = { label, icon };
    }

    /**
     * Register multiple routes
     */
    registerRoutes(routes) {
        Object.assign(this.routes, routes);
    }
}

// CSS for breadcrumbs (should be included in main CSS or as separate file)
const breadcrumbStyles = `
<style>
.breadcrumb-container {
    margin-bottom: 24px;
}

.breadcrumb-list {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    list-style: none;
    padding: 0;
    margin: 0;
}

.breadcrumb-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: var(--text-secondary);
}

.breadcrumb-item:not(:last-child)::after {
    content: '/';
    color: var(--text-muted);
    margin-left: 8px;
    font-weight: 300;
}

.breadcrumb-link {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--text-secondary);
    text-decoration: none;
    padding: 6px 12px;
    border-radius: 8px;
    transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    font-weight: 500;
}

.breadcrumb-link:hover {
    color: var(--primary-color);
    background: rgba(79, 70, 229, 0.1);
    transform: translateY(-2px);
}

.breadcrumb-link i {
    font-size: 16px;
}

.breadcrumb-current {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--text-primary);
    font-weight: 600;
    padding: 6px 12px;
    background: linear-gradient(135deg,
        rgba(79, 70, 229, 0.1) 0%,
        rgba(124, 58, 237, 0.1) 100%);
    border-radius: 8px;
    border-left: 3px solid var(--primary-color);
}

.breadcrumb-current i {
    font-size: 16px;
    color: var(--primary-color);
}

.breadcrumb-item.active {
    color: var(--text-primary);
}

/* Mobile responsive */
@media (max-width: 768px) {
    .breadcrumb-list {
        font-size: 13px;
    }

    .breadcrumb-link,
    .breadcrumb-current {
        padding: 4px 8px;
    }

    .breadcrumb-link i,
    .breadcrumb-current i {
        font-size: 14px;
    }

    /* Hide icons on very small screens */
    @media (max-width: 480px) {
        .breadcrumb-link i,
        .breadcrumb-current i {
            display: none;
        }
    }
}

/* Animation support */
.breadcrumb-item.fade-in-up {
    opacity: 0;
    animation: fadeInUp 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
`;

// Create global instance
const breadcrumbs = new BreadcrumbManager();

// Auto-generate on page load
document.addEventListener('DOMContentLoaded', () => {
    breadcrumbs.generate();
});

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = BreadcrumbManager;
}
