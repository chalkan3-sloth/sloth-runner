// Network Graph Visualization for Agent Dependencies
// Uses D3.js-like force-directed graph to show agent relationships

class NetworkGraph {
    constructor(containerId, options = {}) {
        this.containerId = containerId;
        this.container = document.getElementById(containerId);
        this.width = options.width || this.container?.clientWidth || 800;
        this.height = options.height || 600;
        this.nodes = [];
        this.links = [];
        this.svg = null;
        this.simulation = null;

        this.options = {
            nodeRadius: 30,
            linkDistance: 150,
            charge: -300,
            colors: {
                active: '#10B981',
                inactive: '#6B7280',
                master: '#4F46E5',
                worker: '#3B82F6'
            },
            ...options
        };

        this.initSVG();
    }

    initSVG() {
        if (!this.container) return;

        this.container.innerHTML = `
            <svg width="${this.width}" height="${this.height}" class="network-graph">
                <defs>
                    <marker id="arrowhead" markerWidth="10" markerHeight="10"
                            refX="9" refY="3" orient="auto">
                        <polygon points="0 0, 10 3, 0 6" fill="#999" />
                    </marker>
                    <filter id="glow">
                        <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
                        <feMerge>
                            <feMergeNode in="coloredBlur"/>
                            <feMergeNode in="SourceGraphic"/>
                        </feMerge>
                    </filter>
                </defs>
                <g class="links"></g>
                <g class="nodes"></g>
            </svg>
        `;

        this.svg = this.container.querySelector('svg');
    }

    /**
     * Set graph data
     * @param {Array} nodes - [{id, label, type, status}, ...]
     * @param {Array} links - [{source, target, type}, ...]
     */
    setData(nodes, links) {
        this.nodes = nodes.map(n => ({
            ...n,
            x: this.width / 2 + (Math.random() - 0.5) * 100,
            y: this.height / 2 + (Math.random() - 0.5) * 100
        }));
        this.links = links;
        this.render();
    }

    render() {
        if (!this.svg) return;

        const linksGroup = this.svg.querySelector('.links');
        const nodesGroup = this.svg.querySelector('.nodes');

        // Clear existing
        linksGroup.innerHTML = '';
        nodesGroup.innerHTML = '';

        // Create links
        this.links.forEach(link => {
            const source = this.nodes.find(n => n.id === link.source);
            const target = this.nodes.find(n => n.id === link.target);

            if (source && target) {
                const line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
                line.setAttribute('x1', source.x);
                line.setAttribute('y1', source.y);
                line.setAttribute('x2', target.x);
                line.setAttribute('y2', target.y);
                line.setAttribute('stroke', '#999');
                line.setAttribute('stroke-width', '2');
                line.setAttribute('marker-end', 'url(#arrowhead)');
                line.setAttribute('class', 'network-link');
                linksGroup.appendChild(line);
            }
        });

        // Create nodes
        this.nodes.forEach(node => {
            const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
            g.setAttribute('class', 'network-node');
            g.setAttribute('transform', `translate(${node.x}, ${node.y})`);

            // Node circle
            const circle = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
            circle.setAttribute('r', this.options.nodeRadius);
            circle.setAttribute('fill', this.getNodeColor(node));
            circle.setAttribute('stroke', '#fff');
            circle.setAttribute('stroke-width', '3');
            circle.setAttribute('filter', 'url(#glow)');
            circle.setAttribute('class', 'node-circle');

            // Node label
            const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text.setAttribute('text-anchor', 'middle');
            text.setAttribute('dy', this.options.nodeRadius + 20);
            text.setAttribute('fill', 'var(--text-primary)');
            text.setAttribute('font-size', '12');
            text.setAttribute('font-weight', '600');
            text.textContent = node.label;

            // Status indicator
            if (node.status === 'active') {
                const statusCircle = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
                statusCircle.setAttribute('r', '6');
                statusCircle.setAttribute('cx', this.options.nodeRadius - 5);
                statusCircle.setAttribute('cy', -this.options.nodeRadius + 5);
                statusCircle.setAttribute('fill', this.options.colors.active);
                statusCircle.setAttribute('class', 'status-pulse');
                g.appendChild(statusCircle);
            }

            g.appendChild(circle);
            g.appendChild(text);
            nodesGroup.appendChild(g);

            // Add interaction
            this.addNodeInteraction(g, node);
        });

        this.startSimulation();
    }

    getNodeColor(node) {
        if (node.type === 'master') return this.options.colors.master;
        if (node.type === 'worker') return this.options.colors.worker;
        if (node.status === 'active') return this.options.colors.active;
        return this.options.colors.inactive;
    }

    addNodeInteraction(element, node) {
        let isDragging = false;
        let startX, startY;

        element.addEventListener('mouseenter', () => {
            element.style.cursor = 'pointer';
            const circle = element.querySelector('.node-circle');
            circle.setAttribute('stroke-width', '5');
        });

        element.addEventListener('mouseleave', () => {
            if (!isDragging) {
                const circle = element.querySelector('.node-circle');
                circle.setAttribute('stroke-width', '3');
            }
        });

        element.addEventListener('mousedown', (e) => {
            isDragging = true;
            startX = e.clientX - node.x;
            startY = e.clientY - node.y;
            element.style.cursor = 'grabbing';
        });

        document.addEventListener('mousemove', (e) => {
            if (isDragging) {
                node.x = e.clientX - startX;
                node.y = e.clientY - startY;
                this.updatePositions();
            }
        });

        document.addEventListener('mouseup', () => {
            if (isDragging) {
                isDragging = false;
                element.style.cursor = 'pointer';
            }
        });

        element.addEventListener('click', () => {
            this.onNodeClick(node);
        });
    }

    updatePositions() {
        const linksGroup = this.svg.querySelector('.links');
        const nodesGroup = this.svg.querySelector('.nodes');

        // Update links
        const lines = linksGroup.querySelectorAll('.network-link');
        this.links.forEach((link, i) => {
            const source = this.nodes.find(n => n.id === link.source);
            const target = this.nodes.find(n => n.id === link.target);

            if (source && target && lines[i]) {
                lines[i].setAttribute('x1', source.x);
                lines[i].setAttribute('y1', source.y);
                lines[i].setAttribute('x2', target.x);
                lines[i].setAttribute('y2', target.y);
            }
        });

        // Update nodes
        const nodeGroups = nodesGroup.querySelectorAll('.network-node');
        this.nodes.forEach((node, i) => {
            if (nodeGroups[i]) {
                nodeGroups[i].setAttribute('transform', `translate(${node.x}, ${node.y})`);
            }
        });
    }

    startSimulation() {
        // Simple force simulation
        let iterations = 0;
        const maxIterations = 100;

        const simulate = () => {
            if (iterations++ >= maxIterations) return;

            // Apply forces
            this.nodes.forEach(node => {
                let fx = 0, fy = 0;

                // Repulsion from other nodes
                this.nodes.forEach(other => {
                    if (node.id !== other.id) {
                        const dx = node.x - other.x;
                        const dy = node.y - other.y;
                        const dist = Math.sqrt(dx * dx + dy * dy) || 1;
                        const force = this.options.charge / (dist * dist);
                        fx += (dx / dist) * force;
                        fy += (dy / dist) * force;
                    }
                });

                // Attraction to center
                const centerX = this.width / 2;
                const centerY = this.height / 2;
                fx += (centerX - node.x) * 0.01;
                fy += (centerY - node.y) * 0.01;

                node.x += fx;
                node.y += fy;

                // Keep in bounds
                node.x = Math.max(50, Math.min(this.width - 50, node.x));
                node.y = Math.max(50, Math.min(this.height - 50, node.y));
            });

            this.updatePositions();
            requestAnimationFrame(simulate);
        };

        simulate();
    }

    onNodeClick(node) {
        console.log('Node clicked:', node);
        // Emit event or call callback
        if (this.options.onNodeClick) {
            this.options.onNodeClick(node);
        }
    }

    destroy() {
        if (this.container) {
            this.container.innerHTML = '';
        }
    }
}

// Export
if (typeof module !== 'undefined' && module.exports) {
    module.exports = NetworkGraph;
}
