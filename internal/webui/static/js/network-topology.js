// Network Topology Visualization with D3.js

let simulation = null;
let svg = null;
let g = null;
let link = null;
let node = null;
let label = null;
let topologyData = null;

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    initializeTopology();
    loadTopology();
});

// Initialize D3 topology visualization
function initializeTopology() {
    const container = d3.select('#topology-container');
    const width = container.node().clientWidth;
    const height = container.node().clientHeight;

    // Create SVG
    svg = container.append('svg')
        .attr('width', width)
        .attr('height', height);

    // Add zoom behavior
    const zoom = d3.zoom()
        .scaleExtent([0.1, 4])
        .on('zoom', (event) => {
            g.attr('transform', event.transform);
        });

    svg.call(zoom);

    // Create main group for zoom/pan
    g = svg.append('g');

    // Create force simulation
    simulation = d3.forceSimulation()
        .force('link', d3.forceLink().id(d => d.id).distance(150))
        .force('charge', d3.forceManyBody().strength(-150))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('collision', d3.forceCollide().radius(50));
}

// Load topology data from API
async function loadTopology() {
    try {
        const response = await fetch('/api/v1/network/topology');
        if (!response.ok) throw new Error('Failed to fetch topology');

        topologyData = await response.json();

        // Update stats
        document.getElementById('totalNodes').textContent = topologyData.stats.total_nodes || 0;
        document.getElementById('activeAgents').textContent = topologyData.stats.active_agents || 0;
        document.getElementById('totalLinks').textContent = topologyData.stats.total_links || 0;

        // Render topology
        renderTopology(topologyData);

    } catch (error) {
        console.error('Error loading topology:', error);
        alert('Failed to load network topology: ' + error.message);
    }
}

// Render topology visualization
function renderTopology(data) {
    // Clear existing elements
    g.selectAll('*').remove();

    const nodes = data.nodes || [];
    const links = data.links || [];

    // Create links
    link = g.append('g')
        .selectAll('line')
        .data(links)
        .join('line')
        .attr('class', d => 'link ' + d.type)
        .attr('stroke-width', 2);

    // Create link labels
    const linkLabels = g.append('g')
        .selectAll('text')
        .data(links)
        .join('text')
        .attr('class', 'link-label')
        .text(d => d.type);

    // Create nodes
    const nodeGroup = g.append('g')
        .selectAll('g')
        .data(nodes)
        .join('g')
        .attr('class', d => {
            let classes = 'node';
            if (d.type === 'master') {
                classes += ' node-master';
            } else if (d.type === 'agent') {
                classes += ' node-agent';
                if (d.status === 'Active') {
                    classes += ' active';
                } else {
                    classes += ' offline';
                }
            }
            return classes;
        })
        .call(drag(simulation))
        .on('mouseenter', showNodeTooltip)
        .on('mouseleave', hideNodeTooltip)
        .on('click', nodeClicked);

    // Add circles to nodes
    nodeGroup.append('circle')
        .attr('r', d => d.type === 'master' ? 30 : 20);

    // Add icons to nodes
    nodeGroup.append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', 5)
        .style('font-size', d => d.type === 'master' ? '20px' : '16px')
        .style('fill', 'white')
        .text(d => d.type === 'master' ? '\u2605' : '\u25CF'); // ★ for master, ● for agents

    // Add labels
    label = nodeGroup.append('text')
        .attr('class', 'node-label')
        .attr('text-anchor', 'middle')
        .attr('dy', d => d.type === 'master' ? 45 : 35)
        .style('font-size', '11px')
        .style('font-weight', 'bold')
        .style('fill', '#333')
        .text(d => d.name);

    node = nodeGroup;

    // Update force simulation
    simulation
        .nodes(nodes)
        .on('tick', ticked);

    simulation.force('link')
        .links(links);

    simulation.alpha(1).restart();

    // Tick function for animation
    function ticked() {
        link
            .attr('x1', d => d.source.x)
            .attr('y1', d => d.source.y)
            .attr('x2', d => d.target.x)
            .attr('y2', d => d.target.y);

        linkLabels
            .attr('x', d => (d.source.x + d.target.x) / 2)
            .attr('y', d => (d.source.y + d.target.y) / 2);

        node
            .attr('transform', d => `translate(${d.x},${d.y})`);
    }
}

// Drag behavior
function drag(simulation) {
    function dragstarted(event, d) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
    }

    function dragged(event, d) {
        d.fx = event.x;
        d.fy = event.y;
    }

    function dragended(event, d) {
        if (!event.active) simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
    }

    return d3.drag()
        .on('start', dragstarted)
        .on('drag', dragged)
        .on('end', dragended);
}

// Show node tooltip
function showNodeTooltip(event, d) {
    const tooltip = document.getElementById('nodeTooltip');

    let content = `<strong>${d.name}</strong><br>`;
    content += `Type: ${d.type}<br>`;

    if (d.status) {
        content += `Status: <span style="color: ${d.status === 'Active' ? '#56ab2f' : '#999'}">${d.status}</span><br>`;
    }

    if (d.ip_addresses && d.ip_addresses.length > 0) {
        content += `IPs: ${d.ip_addresses.join(', ')}`;
    }

    tooltip.innerHTML = content;
    tooltip.style.display = 'block';
    tooltip.style.left = (event.pageX + 10) + 'px';
    tooltip.style.top = (event.pageY - 30) + 'px';
}

// Hide node tooltip
function hideNodeTooltip() {
    const tooltip = document.getElementById('nodeTooltip');
    tooltip.style.display = 'none';
}

// Node click handler
function nodeClicked(event, d) {
    if (d.type === 'agent') {
        // Navigate to agent details or network stats for this agent
        if (confirm(`View network details for ${d.name}?`)) {
            window.location.href = `/agent-dashboard?agent=${d.name}`;
        }
    }
}

// Update force strength
function updateForce() {
    const strength = parseInt(document.getElementById('forceStrength').value);
    simulation.force('charge').strength(strength);
    simulation.alpha(0.3).restart();
}

// Toggle labels
function toggleLabels() {
    const showLabels = document.getElementById('showLabels').checked;

    if (label) {
        label.style('display', showLabels ? 'block' : 'none');
    }
}

// Utility function for escaping HTML
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return String(text).replace(/[&<>"']/g, m => map[m]);
}
