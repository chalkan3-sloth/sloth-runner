// Timeline Visualization Component for Workflow Executions
// Displays workflow execution history in a beautiful visual timeline

class TimelineVisualization {
    constructor(containerId) {
        this.containerId = containerId;
        this.container = document.getElementById(containerId);
        this.events = [];
    }

    /**
     * Add an event to the timeline
     */
    addEvent(event) {
        this.events.push(event);
        this.render();
    }

    /**
     * Set all events at once
     */
    setEvents(events) {
        this.events = events;
        this.render();
    }

    /**
     * Render the timeline
     */
    render() {
        if (!this.container) return;

        const html = `
            <div class="timeline">
                ${this.events.map((event, index) => this.renderEvent(event, index)).join('')}
            </div>
        `;

        this.container.innerHTML = html;
        this.attachAnimations();
    }

    /**
     * Render a single event
     */
    renderEvent(event, index) {
        const {
            title,
            description,
            timestamp,
            status = 'completed', // completed, running, failed, pending
            icon = 'bi-check-circle',
            metadata = {}
        } = event;

        const statusClass = `timeline-event-${status}`;
        const position = index % 2 === 0 ? 'left' : 'right';

        return `
            <div class="timeline-item ${statusClass} timeline-${position} fade-in-up" style="animation-delay: ${index * 0.1}s;">
                <div class="timeline-marker">
                    <i class="bi ${icon}"></i>
                </div>
                <div class="timeline-content">
                    <div class="timeline-header">
                        <h5 class="timeline-title">${title}</h5>
                        <span class="timeline-time">${this.formatTimestamp(timestamp)}</span>
                    </div>
                    ${description ? `<p class="timeline-description">${description}</p>` : ''}
                    ${Object.keys(metadata).length > 0 ? this.renderMetadata(metadata) : ''}
                    <span class="timeline-status-badge badge-${status}">${status}</span>
                </div>
            </div>
        `;
    }

    /**
     * Render metadata
     */
    renderMetadata(metadata) {
        const items = Object.entries(metadata).map(([key, value]) => `
            <div class="timeline-meta-item">
                <span class="timeline-meta-key">${key}:</span>
                <span class="timeline-meta-value">${value}</span>
            </div>
        `).join('');

        return `<div class="timeline-metadata">${items}</div>`;
    }

    /**
     * Format timestamp
     */
    formatTimestamp(timestamp) {
        const date = new Date(timestamp);
        const now = new Date();
        const diff = now - date;
        const minutes = Math.floor(diff / 60000);
        const hours = Math.floor(diff / 3600000);
        const days = Math.floor(diff / 86400000);

        if (minutes < 1) return 'Just now';
        if (minutes < 60) return `${minutes}m ago`;
        if (hours < 24) return `${hours}h ago`;
        if (days < 7) return `${days}d ago`;

        return date.toLocaleDateString();
    }

    /**
     * Attach animations
     */
    attachAnimations() {
        const items = this.container.querySelectorAll('.timeline-item');
        items.forEach((item, index) => {
            item.style.animationDelay = `${index * 0.1}s`;
        });
    }

    /**
     * Clear timeline
     */
    clear() {
        this.events = [];
        if (this.container) {
            this.container.innerHTML = '';
        }
    }

    /**
     * Filter events by status
     */
    filterByStatus(status) {
        const filtered = this.events.filter(e => e.status === status);
        this.setEvents(filtered);
    }
}

// CSS for timeline (embedded)
const timelineStyles = `
<style>
.timeline {
    position: relative;
    max-width: 1000px;
    margin: 0 auto;
    padding: 40px 0;
}

.timeline::before {
    content: '';
    position: absolute;
    left: 50%;
    top: 0;
    bottom: 0;
    width: 3px;
    background: linear-gradient(180deg,
        var(--primary-color) 0%,
        var(--secondary-color) 100%);
    transform: translateX(-50%);
}

.timeline-item {
    position: relative;
    margin-bottom: 40px;
    opacity: 0;
}

.timeline-item.timeline-left {
    padding-right: calc(50% + 40px);
    text-align: right;
}

.timeline-item.timeline-right {
    padding-left: calc(50% + 40px);
}

.timeline-marker {
    position: absolute;
    top: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 50px;
    height: 50px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    color: white;
    z-index: 2;
    box-shadow: 0 0 0 4px var(--bg-primary);
    transition: all 0.3s ease;
}

.timeline-event-completed .timeline-marker {
    background: linear-gradient(135deg, var(--success-color), var(--secondary-light));
}

.timeline-event-running .timeline-marker {
    background: linear-gradient(135deg, var(--info-color), #60A5FA);
    animation: pulse 2s infinite;
}

.timeline-event-failed .timeline-marker {
    background: linear-gradient(135deg, var(--danger-color), #F87171);
}

.timeline-event-pending .timeline-marker {
    background: linear-gradient(135deg, var(--warning-color), #FCD34D);
}

.timeline-content {
    background: var(--bg-card);
    border: 1px solid var(--border-color);
    border-radius: 12px;
    padding: 20px;
    box-shadow: var(--shadow-sm);
    transition: all 0.3s ease;
}

.timeline-item:hover .timeline-content {
    transform: translateY(-5px);
    box-shadow: var(--shadow-lg);
}

.timeline-item:hover .timeline-marker {
    transform: translateX(-50%) scale(1.2);
}

.timeline-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.timeline-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
}

.timeline-time {
    font-size: 12px;
    color: var(--text-muted);
}

.timeline-description {
    font-size: 14px;
    color: var(--text-secondary);
    margin-bottom: 12px;
}

.timeline-metadata {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 12px;
}

.timeline-meta-item {
    font-size: 13px;
}

.timeline-meta-key {
    color: var(--text-muted);
    font-weight: 500;
}

.timeline-meta-value {
    color: var(--text-primary);
    margin-left: 4px;
}

.timeline-status-badge {
    display: inline-block;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
}

.badge-completed {
    background: rgba(16, 185, 129, 0.1);
    color: var(--success-color);
}

.badge-running {
    background: rgba(59, 130, 246, 0.1);
    color: var(--info-color);
}

.badge-failed {
    background: rgba(239, 68, 68, 0.1);
    color: var(--danger-color);
}

.badge-pending {
    background: rgba(245, 158, 11, 0.1);
    color: var(--warning-color);
}

/* Mobile responsive */
@media (max-width: 768px) {
    .timeline::before {
        left: 30px;
    }

    .timeline-item {
        padding-left: 70px !important;
        padding-right: 0 !important;
        text-align: left !important;
    }

    .timeline-marker {
        left: 30px;
    }
}
</style>
`;

// Create global instance helper
function createTimeline(containerId) {
    return new TimelineVisualization(containerId);
}

// Export
if (typeof module !== 'undefined' && module.exports) {
    module.exports = TimelineVisualization;
}
