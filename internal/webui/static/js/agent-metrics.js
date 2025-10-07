/* ===================================
   Agent Metrics & Visualizations
   Professional UX Components
   =================================== */

// Prevent duplicate loading
if (typeof window.AGENT_METRICS_LOADED !== 'undefined') {
    console.warn('Agent metrics already loaded, skipping...');
} else {
    window.AGENT_METRICS_LOADED = true;

// Sparkline Chart Generator
class Sparkline {
    constructor(canvas, data, options = {}) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.data = data;
        this.options = {
            color: options.color || '#4F46E5',
            lineWidth: options.lineWidth || 2,
            fillOpacity: options.fillOpacity || 0.2,
            smooth: options.smooth !== false,
            ...options
        };
        this.draw();
    }

    draw() {
        const { ctx, canvas, data, options } = this;
        const width = canvas.width;
        const height = canvas.height;

        ctx.clearRect(0, 0, width, height);

        if (data.length === 0) return;

        const max = Math.max(...data);
        const min = Math.min(...data);
        const range = max - min || 1;

        const xStep = width / (data.length - 1 || 1);
        const yScale = height / range;

        // Create path
        ctx.beginPath();
        ctx.strokeStyle = options.color;
        ctx.lineWidth = options.lineWidth;
        ctx.lineCap = 'round';
        ctx.lineJoin = 'round';

        data.forEach((value, i) => {
            const x = i * xStep;
            const y = height - ((value - min) * yScale);

            if (i === 0) {
                ctx.moveTo(x, y);
            } else {
                if (options.smooth) {
                    const prevX = (i - 1) * xStep;
                    const prevY = height - ((data[i - 1] - min) * yScale);
                    const cpX = (prevX + x) / 2;
                    ctx.quadraticCurveTo(prevX, prevY, cpX, (prevY + y) / 2);
                } else {
                    ctx.lineTo(x, y);
                }
            }
        });

        ctx.stroke();

        // Fill area under line
        if (options.fillOpacity > 0) {
            ctx.lineTo(width, height);
            ctx.lineTo(0, height);
            ctx.closePath();
            ctx.fillStyle = options.color + Math.floor(options.fillOpacity * 255).toString(16).padStart(2, '0');
            ctx.fill();
        }
    }

    update(newData) {
        this.data = newData;
        this.draw();
    }
}

// Resource Ring Progress
class ResourceRing {
    constructor(container, value, options = {}) {
        this.container = container;
        this.value = value;
        this.options = {
            max: options.max || 100,
            size: options.size || 60,
            strokeWidth: options.strokeWidth || 6,
            color: this.getColorForValue(value, options.max),
            ...options
        };
        this.render();
    }

    getColorForValue(value, max = 100) {
        const percent = (value / max) * 100;
        if (percent >= 90) return '#ef4444'; // Critical
        if (percent >= 70) return '#f59e0b'; // Warning
        if (percent >= 40) return '#3b82f6'; // Good
        return '#10b981'; // Excellent
    }

    render() {
        const { size, strokeWidth, color } = this.options;
        const radius = (size - strokeWidth) / 2;
        const circumference = 2 * Math.PI * radius;
        const percent = (this.value / this.options.max) * 100;
        const offset = circumference - (percent / 100) * circumference;

        this.container.innerHTML = `
            <div class="resource-ring">
                <svg width="${size}" height="${size}">
                    <circle class="resource-ring-bg"
                            cx="${size/2}" cy="${size/2}" r="${radius}"/>
                    <circle class="resource-ring-progress"
                            cx="${size/2}" cy="${size/2}" r="${radius}"
                            stroke="${color}"
                            stroke-dasharray="${circumference}"
                            stroke-dashoffset="${offset}"/>
                </svg>
                <div class="resource-ring-text">${Math.round(percent)}%</div>
            </div>
        `;
    }

    update(newValue) {
        this.value = newValue;
        this.options.color = this.getColorForValue(newValue, this.options.max);
        this.render();
    }
}

// Health Indicator
class HealthIndicator {
    static getStatus(metrics) {
        const { cpu = 0, memory = 0, disk = 0 } = metrics;

        // Critical if any resource is above 90%
        if (cpu > 90 || memory > 90 || disk > 90) {
            return { status: 'critical', icon: 'exclamation-triangle-fill', class: 'error' };
        }

        // Warning if any resource is above 70%
        if (cpu > 70 || memory > 70 || disk > 70) {
            return { status: 'warning', icon: 'exclamation-circle-fill', class: 'warning' };
        }

        // Good if all resources are below 70%
        return { status: 'healthy', icon: 'check-circle-fill', class: 'online' };
    }

    static getBadgeHTML(metrics) {
        const { status, icon, class: statusClass } = this.getStatus(metrics);
        return `
            <span class="status-indicator ${statusClass}">
                <span class="status-pulse"></span>
                <i class="bi bi-${icon}"></i>
                ${status}
            </span>
        `;
    }
}

// Progress Bar with Animation
class AnimatedProgress {
    static create(value, max = 100, options = {}) {
        const percent = (value / max) * 100;
        const level = this.getLevel(percent);

        return `
            <div class="metric-progress">
                <div class="metric-progress-bar ${level}"
                     style="width: ${percent}%"
                     data-value="${value}"
                     data-max="${max}">
                </div>
            </div>
        `;
    }

    static getLevel(percent) {
        if (percent >= 90) return 'critical';
        if (percent >= 70) return 'warning';
        if (percent >= 40) return 'good';
        return 'excellent';
    }
}

// Agent Card Builder
class AgentCardBuilder {
    constructor(agent) {
        this.agent = agent;
    }

    build() {
        const { agent } = this;
        const isOnline = agent.status === 'Active' || agent.status === 'online';

        // Mock metrics (replace with real data from API)
        const metrics = {
            cpu: agent.cpu || Math.random() * 100,
            memory: agent.memory || Math.random() * 100,
            disk: agent.disk || Math.random() * 100,
        };

        const health = HealthIndicator.getStatus(metrics);
        const uptime = this.formatUptime(agent.uptime || 0);

        return `
            <div class="col-lg-4 col-md-6 mb-4">
                <div class="card agent-card ${health.class} border-0 shadow-sm hover-lift fade-in-up">
                    <div class="card-body">
                        <!-- Header -->
                        <div class="agent-card-header">
                            <div class="agent-name">
                                <i class="bi bi-server"></i>
                                ${agent.name}
                            </div>
                            ${HealthIndicator.getBadgeHTML(metrics)}
                        </div>

                        <!-- Info Pills -->
                        <div class="agent-info-pills">
                            <div class="info-pill agent-tooltip" data-tooltip="Agent Address">
                                <i class="bi bi-geo-alt"></i>
                                ${agent.address || 'N/A'}
                            </div>
                            ${agent.version ? `
                                <div class="info-pill agent-tooltip" data-tooltip="Agent Version">
                                    <i class="bi bi-tag"></i>
                                    ${agent.version}
                                </div>
                            ` : ''}
                            <div class="info-pill agent-tooltip" data-tooltip="Uptime">
                                <i class="bi bi-clock"></i>
                                ${uptime}
                            </div>
                        </div>

                        <!-- Metrics -->
                        <div class="agent-metrics">
                            <!-- CPU -->
                            <div class="metric-row">
                                <div class="metric-label">
                                    <i class="bi bi-cpu"></i>
                                    CPU
                                </div>
                                <div class="metric-value">${Math.round(metrics.cpu)}%</div>
                            </div>
                            ${AnimatedProgress.create(metrics.cpu, 100)}

                            <!-- Memory -->
                            <div class="metric-row mt-3">
                                <div class="metric-label">
                                    <i class="bi bi-memory"></i>
                                    Memory
                                </div>
                                <div class="metric-value">${Math.round(metrics.memory)}%</div>
                            </div>
                            ${AnimatedProgress.create(metrics.memory, 100)}

                            <!-- Disk (optional) -->
                            ${metrics.disk !== undefined ? `
                                <div class="metric-row mt-3">
                                    <div class="metric-label">
                                        <i class="bi bi-hdd"></i>
                                        Disk
                                    </div>
                                    <div class="metric-value">${Math.round(metrics.disk)}%</div>
                                </div>
                                ${AnimatedProgress.create(metrics.disk, 100)}
                            ` : ''}
                        </div>

                        <!-- Sparkline (last 24h activity) -->
                        <div class="sparkline-container mt-3" id="sparkline-${agent.name}"></div>

                        <!-- Quick Actions -->
                        <div class="agent-actions mt-3">
                            <button class="agent-action-btn primary" onclick="viewAgentDashboard('${agent.name}')">
                                <i class="bi bi-speedometer2"></i>
                                Dashboard
                            </button>
                            <button class="agent-action-btn" onclick="viewAgentDetails('${agent.name}')">
                                <i class="bi bi-info-circle"></i>
                                Details
                            </button>
                            <button class="agent-action-btn" onclick="viewAgentLogs('${agent.name}')">
                                <i class="bi bi-file-text"></i>
                                Logs
                            </button>
                            ${isOnline ? `
                                <button class="agent-action-btn" onclick="restartAgent('${agent.name}')">
                                    <i class="bi bi-arrow-clockwise"></i>
                                    Restart
                                </button>
                            ` : ''}
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    formatUptime(seconds) {
        if (seconds < 60) return `${seconds}s`;
        if (seconds < 3600) return `${Math.floor(seconds / 60)}m`;
        if (seconds < 86400) return `${Math.floor(seconds / 3600)}h`;
        return `${Math.floor(seconds / 86400)}d`;
    }
}

// Initialize sparklines after cards are rendered
function initializeSparklines(agents) {
    agents.forEach(agent => {
        const container = document.getElementById(`sparkline-${agent.name}`);
        if (container) {
            const canvas = document.createElement('canvas');
            canvas.className = 'sparkline-canvas';
            canvas.width = container.offsetWidth * 2; // Retina
            canvas.height = 80; // 40px display height
            container.appendChild(canvas);

            // Generate sample data (replace with real metrics)
            const data = Array.from({ length: 24 }, () => Math.random() * 100);
            new Sparkline(canvas, data, {
                color: '#4F46E5',
                lineWidth: 3,
                fillOpacity: 0.1,
                smooth: true
            });
        }
    });
}

// Export to global scope
window.AgentCardBuilder = AgentCardBuilder;
window.Sparkline = Sparkline;
window.ResourceRing = ResourceRing;
window.HealthIndicator = HealthIndicator;
window.AnimatedProgress = AnimatedProgress;
window.initializeSparklines = initializeSparklines;

} // End of duplicate loading check
