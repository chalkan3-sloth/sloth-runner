// Sloth Runner - Shared Navbar Component
// This script loads and manages the consistent navigation bar across all pages

const SlothNavbar = {
    // Get the navbar HTML
    getHTML: function() {
        return `
        <!-- Top Navigation Bar -->
        <nav class="navbar navbar-expand-lg">
            <div class="container-fluid">
                <a class="navbar-brand" href="/">
                    <img src="/static/img/sloth-logo.svg" alt="Sloth" class="sloth-logo" onerror="this.style.display='none'">
                    <span class="fw-bold">🦥 Sloth Runner</span>
                </a>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarNav">
                    <!-- Search Bar -->
                    <div class="navbar-search mx-3">
                        <div class="search-container">
                            <i class="bi bi-search search-icon"></i>
                            <input type="search" id="navbar-search" class="form-control" placeholder="Search... (Ctrl+K)" autocomplete="off">
                            <kbd class="search-shortcut">⌘K</kbd>
                        </div>
                        <div id="search-results" class="search-results"></div>
                    </div>

                    <ul class="navbar-nav me-auto">
                        <li class="nav-item">
                            <a class="nav-link" href="/" data-page="index">
                                <i class="bi bi-speedometer2"></i> Dashboard
                            </a>
                        </li>
                        <li class="nav-item dropdown">
                            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
                                <i class="bi bi-layers"></i> Management
                            </a>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item" href="/agents" data-page="agents"><i class="bi bi-hdd-network"></i> Agents</a></li>
                                <li><a class="dropdown-item" href="/agent-groups" data-page="agent-groups"><i class="bi bi-collection"></i> Agent Groups</a></li>
                                <li><a class="dropdown-item" href="/agent-control" data-page="agent-control"><i class="bi bi-gear"></i> Agent Control</a></li>
                                <li><a class="dropdown-item" href="/workflows" data-page="workflows"><i class="bi bi-diagram-3"></i> Workflows</a></li>
                                <li><a class="dropdown-item" href="/stacks" data-page="stacks"><i class="bi bi-layers"></i> Stacks</a></li>
                                <li><a class="dropdown-item" href="/hooks" data-page="hooks"><i class="bi bi-hook"></i> Hooks</a></li>
                                <li><a class="dropdown-item" href="/events" data-page="events"><i class="bi bi-bell"></i> Events</a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item" href="/secrets" data-page="secrets"><i class="bi bi-shield-lock"></i> Secrets</a></li>
                                <li><a class="dropdown-item" href="/ssh" data-page="ssh"><i class="bi bi-key"></i> SSH Profiles</a></li>
                            </ul>
                        </li>
                        <li class="nav-item dropdown">
                            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
                                <i class="bi bi-tools"></i> Operations
                            </a>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item" href="/executions" data-page="executions"><i class="bi bi-play-circle"></i> Executions</a></li>
                                <li><a class="dropdown-item" href="/scheduler" data-page="scheduler"><i class="bi bi-calendar-event"></i> Scheduler</a></li>
                                <li><a class="dropdown-item" href="/terminal" data-page="terminal"><i class="bi bi-terminal"></i> Terminal</a></li>
                            </ul>
                        </li>
                        <li class="nav-item dropdown">
                            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
                                <i class="bi bi-graph-up"></i> Monitoring
                            </a>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item" href="/agent-control" data-page="agent-control"><i class="bi bi-bar-chart"></i> Agent Dashboards</a></li>
                                <li><a class="dropdown-item" href="/metrics" data-page="metrics"><i class="bi bi-speedometer"></i> System Metrics</a></li>
                                <li><a class="dropdown-item" href="/logs" data-page="logs"><i class="bi bi-file-text"></i> Logs</a></li>
                            </ul>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link" href="/backup" data-page="backup">
                                <i class="bi bi-server"></i> Backup
                            </a>
                        </li>
                    </ul>
                    <ul class="navbar-nav">
                        <li class="nav-item">
                            <button class="nav-link btn btn-link theme-toggle" onclick="themeManager && themeManager.toggle()">
                                <i class="bi bi-moon-stars-fill"></i>
                            </button>
                        </li>
                        <li class="nav-item">
                            <span class="nav-link">
                                <i class="bi bi-circle-fill status-pulse" id="ws-status" style="color: #6c757d;"></i>
                                <span id="ws-status-text">Connecting...</span>
                            </span>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
        `;
    },

    // Initialize navbar
    init: function() {
        // Find navbar placeholder or create one
        let navbarContainer = document.getElementById('sloth-navbar');

        if (!navbarContainer) {
            // If no placeholder, insert at the beginning of body
            navbarContainer = document.createElement('div');
            navbarContainer.id = 'sloth-navbar';
            document.body.insertBefore(navbarContainer, document.body.firstChild);
        }

        // Insert navbar HTML
        navbarContainer.innerHTML = this.getHTML();

        // Set active page
        this.setActivePage();
    },

    // Set active page in navbar
    setActivePage: function() {
        const currentPath = window.location.pathname;
        const pageName = currentPath.substring(1) || 'index';

        // Remove all active classes
        document.querySelectorAll('.navbar .nav-link, .navbar .dropdown-item').forEach(link => {
            link.classList.remove('active');
        });

        // Add active class to current page
        const activeLink = document.querySelector(`[data-page="${pageName}"]`);
        if (activeLink) {
            activeLink.classList.add('active');

            // If it's a dropdown item, also mark the dropdown as active
            const dropdown = activeLink.closest('.dropdown');
            if (dropdown) {
                const dropdownToggle = dropdown.querySelector('.dropdown-toggle');
                if (dropdownToggle) {
                    dropdownToggle.classList.add('active');
                }
            }
        }
    },

    // Initialize search functionality
    initSearch: function() {
        const searchInput = document.getElementById('navbar-search');
        const searchResults = document.getElementById('search-results');

        if (!searchInput) return;

        // Search index
        const searchIndex = [
            { title: 'Dashboard', url: '/', icon: 'speedometer2', description: 'Overview and system status' },
            { title: 'Agents', url: '/agents', icon: 'hdd-network', description: 'Manage remote agents' },
            { title: 'Agent Groups', url: '/agent-groups', icon: 'collection', description: 'Organize and manage agent groups' },
            { title: 'Agent Control', url: '/agent-control', icon: 'gear', description: 'Control and monitor agents' },
            { title: 'Workflows', url: '/workflows', icon: 'diagram-3', description: 'Manage workflows and tasks' },
            { title: 'Stacks', url: '/stacks', icon: 'layers', description: 'Infrastructure stacks' },
            { title: 'Hooks', url: '/hooks', icon: 'hook', description: 'Event hooks and triggers' },
            { title: 'Events', url: '/events', icon: 'bell', description: 'System events log' },
            { title: 'Secrets', url: '/secrets', icon: 'shield-lock', description: 'Manage secrets and credentials' },
            { title: 'SSH Profiles', url: '/ssh', icon: 'key', description: 'SSH connection profiles' },
            { title: 'Executions', url: '/executions', icon: 'play-circle', description: 'Workflow executions history' },
            { title: 'Scheduler', url: '/scheduler', icon: 'calendar-event', description: 'Schedule tasks and workflows' },
            { title: 'Terminal', url: '/terminal', icon: 'terminal', description: 'Web-based terminal' },
            { title: 'Metrics', url: '/metrics', icon: 'speedometer', description: 'System metrics and monitoring' },
            { title: 'Logs', url: '/logs', icon: 'file-text', description: 'Application logs' },
            { title: 'Backup', url: '/backup', icon: 'server', description: 'Backup and restore' }
        ];

        // Search function
        const performSearch = window.debounce((query) => {
            if (!query || query.length < 2) {
                searchResults.classList.remove('show');
                return;
            }

            const results = searchIndex.filter(item =>
                item.title.toLowerCase().includes(query.toLowerCase()) ||
                item.description.toLowerCase().includes(query.toLowerCase())
            ).slice(0, 5);

            if (results.length === 0) {
                searchResults.innerHTML = '<div class="search-no-results">No results found</div>';
            } else {
                searchResults.innerHTML = results.map(item => `
                    <a href="${item.url}" class="search-result-item">
                        <i class="bi bi-${item.icon}"></i>
                        <div class="search-result-content">
                            <div class="search-result-title">${item.title}</div>
                            <div class="search-result-description">${item.description}</div>
                        </div>
                    </a>
                `).join('');
            }

            searchResults.classList.add('show');
        }, 300);

        // Event listeners
        searchInput.addEventListener('input', (e) => performSearch(e.target.value));

        searchInput.addEventListener('focus', () => {
            if (searchInput.value.length >= 2) {
                searchResults.classList.add('show');
            }
        });

        // Close on click outside
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.navbar-search')) {
                searchResults.classList.remove('show');
            }
        });

        // Keyboard navigation
        searchInput.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                searchResults.classList.remove('show');
                searchInput.blur();
            }
        });
    }
};

// Auto-initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        SlothNavbar.init();
        setTimeout(() => SlothNavbar.initSearch(), 100);
    });
} else {
    SlothNavbar.init();
    setTimeout(() => SlothNavbar.initSearch(), 100);
}

