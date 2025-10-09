// Agent Groups Management

let currentGroupId = null;

// Load all groups
async function loadGroups() {
    try {
        const response = await fetch('/api/v1/agent-groups');
        const data = await response.json();

        const container = document.getElementById('groups-container');
        if (!data.groups || data.groups.length === 0) {
            container.innerHTML = `
                <div class="col-12">
                    <div class="alert alert-info">
                        <i class="bi bi-info-circle"></i> No groups found. Create your first group!
                    </div>
                </div>
            `;
            return;
        }

        container.innerHTML = data.groups.map(group => `
            <div class="col-md-6 col-lg-4 mb-3">
                <div class="card hover-grow">
                    <div class="card-body">
                        <h5 class="card-title">
                            <i class="bi bi-collection"></i> ${escapeHtml(group.name)}
                        </h5>
                        <p class="card-text text-muted">${escapeHtml(group.description || 'No description')}</p>
                        <div class="mb-2">
                            <span class="badge bg-primary">${group.agent_count || 0} agents</span>
                            ${Object.keys(group.tags || {}).map(key =>
                                `<span class="badge bg-secondary">${escapeHtml(key)}: ${escapeHtml(group.tags[key])}</span>`
                            ).join(' ')}
                        </div>
                        <div class="btn-group" role="group">
                            <button class="btn btn-sm btn-outline-primary" onclick="viewGroupDetails('${escapeHtml(group.id)}')">
                                <i class="bi bi-eye"></i> View
                            </button>
                            <button class="btn btn-sm btn-outline-info" onclick="manageAgents('${escapeHtml(group.id)}')">
                                <i class="bi bi-people"></i> Manage
                            </button>
                            <button class="btn btn-sm btn-outline-success" onclick="viewMetrics('${escapeHtml(group.id)}')">
                                <i class="bi bi-graph-up"></i> Metrics
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');
    } catch (error) {
        showToast('Failed to load groups: ' + error.message, 'error');
    }
}

// Load templates
async function loadTemplates() {
    try {
        const response = await fetch('/api/v1/agent-groups/templates');
        const data = await response.json();

        const container = document.getElementById('templates-container');
        if (!data.templates || data.templates.length === 0) {
            container.innerHTML = `
                <div class="col-12">
                    <div class="alert alert-info">
                        <i class="bi bi-info-circle"></i> No templates found. Create a template to automate group creation!
                    </div>
                </div>
            `;
            return;
        }

        container.innerHTML = data.templates.map(template => `
            <div class="col-md-6 col-lg-4 mb-3">
                <div class="card hover-grow">
                    <div class="card-body">
                        <h5 class="card-title">
                            <i class="bi bi-file-earmark-text"></i> ${escapeHtml(template.name)}
                        </h5>
                        <p class="card-text text-muted">${escapeHtml(template.description || 'No description')}</p>
                        <div class="mb-2">
                            <span class="badge bg-info">${template.rules?.length || 0} rules</span>
                        </div>
                        <div class="mb-2">
                            <small class="text-muted">Rules:</small>
                            ${(template.rules || []).map(rule => `
                                <div class="small">
                                    <code>${escapeHtml(rule.type)} ${escapeHtml(rule.operator)} "${escapeHtml(rule.value)}"</code>
                                </div>
                            `).join('')}
                        </div>
                        <div class="btn-group" role="group">
                            <button class="btn btn-sm btn-success" onclick="applyTemplate('${escapeHtml(template.id)}')">
                                <i class="bi bi-play-fill"></i> Apply
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deleteTemplate('${escapeHtml(template.id)}')">
                                <i class="bi bi-trash"></i> Delete
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');
    } catch (error) {
        showToast('Failed to load templates: ' + error.message, 'error');
    }
}

// Load auto-discovery configs
async function loadAutoDiscovery() {
    try {
        const response = await fetch('/api/v1/agent-groups/auto-discovery');
        const data = await response.json();

        const container = document.getElementById('autodiscovery-container');
        if (!data.configs || data.configs.length === 0) {
            container.innerHTML = `
                <div class="col-12">
                    <div class="alert alert-info">
                        <i class="bi bi-info-circle"></i> No auto-discovery configs found.
                    </div>
                </div>
            `;
            return;
        }

        container.innerHTML = data.configs.map(config => `
            <div class="col-md-6 col-lg-4 mb-3">
                <div class="card hover-grow">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5 class="card-title">
                                <i class="bi bi-search"></i> ${escapeHtml(config.name)}
                            </h5>
                            <span class="badge ${config.enabled ? 'bg-success' : 'bg-secondary'}">
                                ${config.enabled ? 'Enabled' : 'Disabled'}
                            </span>
                        </div>
                        <p class="card-text text-muted">${escapeHtml(config.description || 'No description')}</p>
                        <div class="mb-2">
                            <small><strong>Target:</strong> ${escapeHtml(config.target_group)}</small><br>
                            ${config.schedule ? `<small><strong>Schedule:</strong> <code>${escapeHtml(config.schedule)}</code></small>` : ''}
                        </div>
                        <div class="btn-group" role="group">
                            <button class="btn btn-sm btn-primary" onclick="runDiscovery('${escapeHtml(config.id)}')">
                                <i class="bi bi-play-fill"></i> Run Now
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deleteAutoDiscovery('${escapeHtml(config.id)}')">
                                <i class="bi bi-trash"></i> Delete
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');
    } catch (error) {
        showToast('Failed to load auto-discovery configs: ' + error.message, 'error');
    }
}

// Load webhooks
async function loadWebhooks() {
    try {
        const response = await fetch('/api/v1/agent-groups/webhooks');
        const data = await response.json();

        const container = document.getElementById('webhooks-container');
        if (!data.webhooks || data.webhooks.length === 0) {
            container.innerHTML = `
                <div class="col-12">
                    <div class="alert alert-info">
                        <i class="bi bi-info-circle"></i> No webhooks configured.
                    </div>
                </div>
            `;
            return;
        }

        container.innerHTML = data.webhooks.map(webhook => `
            <div class="col-md-6 col-lg-4 mb-3">
                <div class="card hover-grow">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5 class="card-title">
                                <i class="bi bi-webhook"></i> ${escapeHtml(webhook.name)}
                            </h5>
                            <span class="badge ${webhook.enabled ? 'bg-success' : 'bg-secondary'}">
                                ${webhook.enabled ? 'Enabled' : 'Disabled'}
                            </span>
                        </div>
                        <div class="mb-2">
                            <small class="text-muted">${escapeHtml(webhook.url)}</small>
                        </div>
                        <div class="mb-2">
                            ${(webhook.events || []).map(event =>
                                `<span class="badge bg-info">${escapeHtml(event)}</span>`
                            ).join(' ')}
                        </div>
                        <div class="btn-group" role="group">
                            <button class="btn btn-sm btn-info" onclick="viewWebhookLogs('${escapeHtml(webhook.id)}')">
                                <i class="bi bi-clock-history"></i> Logs
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deleteWebhook('${escapeHtml(webhook.id)}')">
                                <i class="bi bi-trash"></i> Delete
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');
    } catch (error) {
        showToast('Failed to load webhooks: ' + error.message, 'error');
    }
}

// Load groups for bulk operations dropdown
async function loadGroupsForBulkOps() {
    try {
        const response = await fetch('/api/v1/agent-groups');
        const data = await response.json();

        const select = document.getElementById('bulkGroupSelect');
        const discoverySelect = document.getElementById('discoveryTargetGroup');

        if (data.groups && data.groups.length > 0) {
            const options = data.groups.map(group =>
                `<option value="${escapeHtml(group.id)}">${escapeHtml(group.name)}</option>`
            ).join('');

            if (select) select.innerHTML += options;
            if (discoverySelect) discoverySelect.innerHTML += options;
        }
    } catch (error) {
        console.error('Failed to load groups for dropdowns:', error);
    }
}

// Refresh all data
function refreshGroups() {
    loadGroups();
    loadTemplates();
    loadAutoDiscovery();
    loadWebhooks();
    loadGroupsForBulkOps();
    showToast('Data refreshed', 'success');
}

// Show create group modal
function showCreateGroupModal() {
    const modal = new bootstrap.Modal(document.getElementById('createGroupModal'));
    modal.show();
}

// Create group
async function createGroup(event) {
    event.preventDefault();

    const name = document.getElementById('groupName').value;
    const description = document.getElementById('groupDescription').value;
    const tagsText = document.getElementById('groupTags').value;

    let tags = {};
    if (tagsText) {
        try {
            tags = JSON.parse(tagsText);
        } catch (error) {
            showToast('Invalid JSON in tags field', 'error');
            return;
        }
    }

    try {
        const response = await fetch('/api/v1/agent-groups', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                group_name: name,
                description: description,
                tags: tags,
                agent_names: []
            })
        });

        if (response.ok) {
            showToast('Group created successfully', 'success');
            bootstrap.Modal.getInstance(document.getElementById('createGroupModal')).hide();
            loadGroups();
            document.getElementById('groupName').value = '';
            document.getElementById('groupDescription').value = '';
            document.getElementById('groupTags').value = '';
        } else {
            const error = await response.json();
            showToast('Failed to create group: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to create group: ' + error.message, 'error');
    }
}

// View group details
async function viewGroupDetails(groupId) {
    try {
        const response = await fetch(`/api/v1/agent-groups/${groupId}`);
        const group = await response.json();

        currentGroupId = groupId;

        const body = document.getElementById('groupDetailsBody');
        body.innerHTML = `
            <div class="mb-3">
                <h6>Name</h6>
                <p>${escapeHtml(group.name)}</p>
            </div>
            <div class="mb-3">
                <h6>Description</h6>
                <p>${escapeHtml(group.description || 'No description')}</p>
            </div>
            <div class="mb-3">
                <h6>Tags</h6>
                ${Object.keys(group.tags || {}).map(key =>
                    `<span class="badge bg-secondary me-1">${escapeHtml(key)}: ${escapeHtml(group.tags[key])}</span>`
                ).join('') || 'No tags'}
            </div>
            <div class="mb-3">
                <h6>Agents (${group.agent_count || 0})</h6>
                <ul class="list-group">
                    ${(group.agent_names || []).map(name =>
                        `<li class="list-group-item">${escapeHtml(name)}</li>`
                    ).join('') || '<li class="list-group-item">No agents</li>'}
                </ul>
            </div>
        `;

        const modal = new bootstrap.Modal(document.getElementById('groupDetailsModal'));
        modal.show();
    } catch (error) {
        showToast('Failed to load group details: ' + error.message, 'error');
    }
}

// Delete group
async function deleteGroup() {
    if (!currentGroupId) return;

    if (!confirm('Are you sure you want to delete this group?')) return;

    try {
        const response = await fetch(`/api/v1/agent-groups/${currentGroupId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            showToast('Group deleted successfully', 'success');
            bootstrap.Modal.getInstance(document.getElementById('groupDetailsModal')).hide();
            loadGroups();
            currentGroupId = null;
        } else {
            const error = await response.json();
            showToast('Failed to delete group: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to delete group: ' + error.message, 'error');
    }
}

// View metrics
async function viewMetrics(groupId) {
    try {
        const response = await fetch(`/api/v1/agent-groups/${groupId}/metrics`);
        const metrics = await response.json();

        showToast(`Group Metrics - Total: ${metrics.total_agents}, Healthy: ${metrics.healthy_agents}, CPU: ${metrics.avg_cpu_percent.toFixed(2)}%`, 'info');
    } catch (error) {
        showToast('Failed to load metrics: ' + error.message, 'error');
    }
}

// Manage agents
let currentManageGroupId = null;

async function manageAgents(groupId) {
    currentManageGroupId = groupId;

    // Open modal
    const modal = new bootstrap.Modal(document.getElementById('manageAgentsModal'));
    modal.show();

    // Load data
    await loadManageAgentsData();
}

async function loadManageAgentsData() {
    try {
        // Load all agents
        const agentsResponse = await fetch('/api/v1/agents');
        const agentsData = await agentsResponse.json();

        // Load group details
        const groupResponse = await fetch(`/api/v1/agent-groups/${currentManageGroupId}`);
        const groupData = await groupResponse.json();

        const allAgents = agentsData.agents || [];
        const groupMembers = groupData.agent_names || [];

        // Populate available agents (not in group)
        const availableAgents = allAgents.filter(agent => !groupMembers.includes(agent.name));
        const availableList = document.getElementById('availableAgentsList');

        if (availableAgents.length === 0) {
            availableList.innerHTML = '<p class="text-muted text-center py-3">No available agents</p>';
        } else {
            availableList.innerHTML = availableAgents.map(agent => `
                <div class="d-flex justify-content-between align-items-center p-2 border-bottom">
                    <div>
                        <i class="bi bi-hdd-network"></i> ${escapeHtml(agent.name)}
                        <br><small class="text-muted">${escapeHtml(agent.address)}</small>
                    </div>
                    <button class="btn btn-sm btn-success" onclick="addAgentToGroup('${escapeHtml(agent.name)}')">
                        <i class="bi bi-plus"></i> Add
                    </button>
                </div>
            `).join('');
        }

        // Populate group members
        const membersList = document.getElementById('groupMembersList');

        if (groupMembers.length === 0) {
            membersList.innerHTML = '<p class="text-muted text-center py-3">No members</p>';
        } else {
            membersList.innerHTML = groupMembers.map(memberName => `
                <div class="d-flex justify-content-between align-items-center p-2 border-bottom">
                    <div>
                        <i class="bi bi-hdd-network text-success"></i> ${escapeHtml(memberName)}
                    </div>
                    <button class="btn btn-sm btn-danger" onclick="removeAgentFromGroup('${escapeHtml(memberName)}')">
                        <i class="bi bi-dash"></i> Remove
                    </button>
                </div>
            `).join('');
        }
    } catch (error) {
        showToast('Failed to load agents: ' + error.message, 'error');
    }
}

async function addAgentToGroup(agentName) {
    try {
        const response = await fetch(`/api/v1/agent-groups/${currentManageGroupId}/agents`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                agent_names: [agentName]
            })
        });

        if (response.ok) {
            showToast(`Agent ${agentName} added to group`, 'success');
            await loadManageAgentsData();
            loadGroups(); // Refresh main view
        } else {
            const error = await response.json();
            showToast('Failed to add agent: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to add agent: ' + error.message, 'error');
    }
}

async function removeAgentFromGroup(agentName) {
    try {
        const response = await fetch(`/api/v1/agent-groups/${currentManageGroupId}/agents`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                agent_names: [agentName]
            })
        });

        if (response.ok) {
            showToast(`Agent ${agentName} removed from group`, 'success');
            await loadManageAgentsData();
            loadGroups(); // Refresh main view
        } else {
            const error = await response.json();
            showToast('Failed to remove agent: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to remove agent: ' + error.message, 'error');
    }
}

// Template functions
function showCreateTemplateModal() {
    const modal = new bootstrap.Modal(document.getElementById('createTemplateModal'));
    modal.show();
}

function addRule() {
    const container = document.getElementById('rulesContainer');
    const ruleDiv = document.createElement('div');
    ruleDiv.className = 'rule-item mb-2';
    ruleDiv.innerHTML = `
        <div class="row g-2">
            <div class="col-md-3">
                <select class="form-select rule-type">
                    <option value="name_pattern">Name Pattern</option>
                    <option value="status">Status</option>
                    <option value="tag_match">Tag Match</option>
                </select>
            </div>
            <div class="col-md-3">
                <select class="form-select rule-operator">
                    <option value="equals">Equals</option>
                    <option value="contains">Contains</option>
                    <option value="regex">Regex</option>
                </select>
            </div>
            <div class="col-md-5">
                <input type="text" class="form-control rule-value" placeholder="Value">
            </div>
            <div class="col-md-1">
                <button type="button" class="btn btn-sm btn-danger" onclick="removeRule(this)">
                    <i class="bi bi-trash"></i>
                </button>
            </div>
        </div>
    `;
    container.appendChild(ruleDiv);
}

function removeRule(button) {
    button.closest('.rule-item').remove();
}

async function createTemplate(event) {
    event.preventDefault();

    const name = document.getElementById('templateName').value;
    const description = document.getElementById('templateDescription').value;
    const tagsText = document.getElementById('templateTags').value;

    // Collect rules
    const rules = [];
    document.querySelectorAll('#rulesContainer .rule-item').forEach(item => {
        const type = item.querySelector('.rule-type').value;
        const operator = item.querySelector('.rule-operator').value;
        const value = item.querySelector('.rule-value').value;

        if (value) {
            rules.push({ type, operator, value });
        }
    });

    let tags = {};
    if (tagsText) {
        try {
            tags = JSON.parse(tagsText);
        } catch (error) {
            showToast('Invalid JSON in tags field', 'error');
            return;
        }
    }

    try {
        const response = await fetch('/api/v1/agent-groups/templates', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                id: name.toLowerCase().replace(/\s+/g, '-'),
                name: name,
                description: description,
                rules: rules,
                tags: tags
            })
        });

        if (response.ok) {
            showToast('Template created successfully', 'success');
            bootstrap.Modal.getInstance(document.getElementById('createTemplateModal')).hide();
            loadTemplates();
        } else {
            const error = await response.json();
            showToast('Failed to create template: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to create template: ' + error.message, 'error');
    }
}

async function applyTemplate(templateId) {
    const groupName = prompt('Enter name for the new group:');
    if (!groupName) return;

    try {
        const response = await fetch(`/api/v1/agent-groups/templates/${templateId}/apply`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                group_name: groupName,
                description: `Created from template ${templateId}`
            })
        });

        if (response.ok) {
            showToast('Template applied successfully', 'success');
            loadGroups();
        } else {
            const error = await response.json();
            showToast('Failed to apply template: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to apply template: ' + error.message, 'error');
    }
}

async function deleteTemplate(templateId) {
    if (!confirm('Are you sure you want to delete this template?')) return;

    try {
        const response = await fetch(`/api/v1/agent-groups/templates/${templateId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            showToast('Template deleted successfully', 'success');
            loadTemplates();
        } else {
            const error = await response.json();
            showToast('Failed to delete template: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to delete template: ' + error.message, 'error');
    }
}

// Auto-discovery functions
function showCreateAutoDiscoveryModal() {
    const modal = new bootstrap.Modal(document.getElementById('createAutoDiscoveryModal'));
    modal.show();
}

function addDiscoveryRule() {
    const container = document.getElementById('discoveryRulesContainer');
    const ruleDiv = document.createElement('div');
    ruleDiv.className = 'rule-item mb-2';
    ruleDiv.innerHTML = `
        <div class="row g-2">
            <div class="col-md-3">
                <select class="form-select discovery-rule-type">
                    <option value="name_pattern">Name Pattern</option>
                    <option value="status">Status</option>
                    <option value="tag_match">Tag Match</option>
                </select>
            </div>
            <div class="col-md-3">
                <select class="form-select discovery-rule-operator">
                    <option value="equals">Equals</option>
                    <option value="contains">Contains</option>
                    <option value="regex">Regex</option>
                </select>
            </div>
            <div class="col-md-5">
                <input type="text" class="form-control discovery-rule-value" placeholder="Value">
            </div>
            <div class="col-md-1">
                <button type="button" class="btn btn-sm btn-danger" onclick="removeDiscoveryRule(this)">
                    <i class="bi bi-trash"></i>
                </button>
            </div>
        </div>
    `;
    container.appendChild(ruleDiv);
}

function removeDiscoveryRule(button) {
    button.closest('.rule-item').remove();
}

async function createAutoDiscovery(event) {
    event.preventDefault();

    const name = document.getElementById('discoveryName').value;
    const targetGroup = document.getElementById('discoveryTargetGroup').value;
    const schedule = document.getElementById('discoverySchedule').value;
    const description = document.getElementById('discoveryDescription').value;
    const enabled = document.getElementById('discoveryEnabled').checked;

    // Collect rules
    const rules = [];
    document.querySelectorAll('#discoveryRulesContainer .rule-item').forEach(item => {
        const type = item.querySelector('.discovery-rule-type').value;
        const operator = item.querySelector('.discovery-rule-operator').value;
        const value = item.querySelector('.discovery-rule-value').value;

        if (value) {
            rules.push({ type, operator, value });
        }
    });

    try {
        const response = await fetch('/api/v1/agent-groups/auto-discovery', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                id: name.toLowerCase().replace(/\s+/g, '-'),
                name: name,
                target_group: targetGroup,
                schedule: schedule,
                description: description,
                enabled: enabled,
                rules: rules
            })
        });

        if (response.ok) {
            showToast('Auto-discovery config created successfully', 'success');
            bootstrap.Modal.getInstance(document.getElementById('createAutoDiscoveryModal')).hide();
            loadAutoDiscovery();
        } else {
            const error = await response.json();
            showToast('Failed to create config: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to create config: ' + error.message, 'error');
    }
}

async function runDiscovery(configId) {
    try {
        const response = await fetch(`/api/v1/agent-groups/auto-discovery/${configId}/run`, {
            method: 'POST'
        });

        if (response.ok) {
            const result = await response.json();
            showToast(`Discovery completed: ${result.agents_added} agents added`, 'success');
            loadGroups();
        } else {
            const error = await response.json();
            showToast('Failed to run discovery: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to run discovery: ' + error.message, 'error');
    }
}

async function deleteAutoDiscovery(configId) {
    if (!confirm('Are you sure you want to delete this auto-discovery config?')) return;

    try {
        const response = await fetch(`/api/v1/agent-groups/auto-discovery/${configId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            showToast('Config deleted successfully', 'success');
            loadAutoDiscovery();
        } else {
            const error = await response.json();
            showToast('Failed to delete config: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to delete config: ' + error.message, 'error');
    }
}

// Webhook functions
function showCreateWebhookModal() {
    const modal = new bootstrap.Modal(document.getElementById('createWebhookModal'));
    modal.show();
}

async function createWebhook(event) {
    event.preventDefault();

    const name = document.getElementById('webhookName').value;
    const url = document.getElementById('webhookURL').value;
    const secret = document.getElementById('webhookSecret').value;
    const retryCount = parseInt(document.getElementById('webhookRetry').value);
    const enabled = document.getElementById('webhookEnabled').checked;

    // Collect selected events
    const events = [];
    document.querySelectorAll('.webhook-event:checked').forEach(checkbox => {
        events.push(checkbox.value);
    });

    if (events.length === 0) {
        showToast('Please select at least one event', 'error');
        return;
    }

    try {
        const response = await fetch('/api/v1/agent-groups/webhooks', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                id: name.toLowerCase().replace(/\s+/g, '-'),
                name: name,
                url: url,
                events: events,
                secret: secret,
                retry_count: retryCount,
                timeout: 30,
                enabled: enabled
            })
        });

        if (response.ok) {
            showToast('Webhook created successfully', 'success');
            bootstrap.Modal.getInstance(document.getElementById('createWebhookModal')).hide();
            loadWebhooks();
        } else {
            const error = await response.json();
            showToast('Failed to create webhook: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to create webhook: ' + error.message, 'error');
    }
}

async function deleteWebhook(webhookId) {
    if (!confirm('Are you sure you want to delete this webhook?')) return;

    try {
        const response = await fetch(`/api/v1/agent-groups/webhooks/${webhookId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            showToast('Webhook deleted successfully', 'success');
            loadWebhooks();
        } else {
            const error = await response.json();
            showToast('Failed to delete webhook: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to delete webhook: ' + error.message, 'error');
    }
}

async function viewWebhookLogs(webhookId) {
    try {
        const response = await fetch(`/api/v1/agent-groups/webhooks/${webhookId}/logs`);
        const data = await response.json();

        if (data.logs && data.logs.length > 0) {
            const logsHtml = data.logs.map(log => `
                <div class="alert ${log.success ? 'alert-success' : 'alert-danger'}">
                    <strong>${escapeHtml(log.event_type)}</strong> -
                    ${new Date(log.timestamp * 1000).toLocaleString()}
                    ${log.status_code ? ` (${log.status_code})` : ''}
                    ${log.error ? `<br><small>${escapeHtml(log.error)}</small>` : ''}
                </div>
            `).join('');

            showToast(`<h6>Recent Webhook Deliveries</h6>${logsHtml}`, 'info');
        } else {
            showToast('No logs found for this webhook', 'info');
        }
    } catch (error) {
        showToast('Failed to load webhook logs: ' + error.message, 'error');
    }
}

// Bulk operations
async function executeBulkOperation(event) {
    event.preventDefault();

    const groupId = document.getElementById('bulkGroupSelect').value;
    const operation = document.getElementById('bulkOperation').value;
    const timeout = parseInt(document.getElementById('bulkTimeout').value);

    let params = {};
    if (operation === 'execute_command') {
        params.command = document.getElementById('bulkCommand').value;
    }

    try {
        showToast('Executing bulk operation...', 'info');

        const response = await fetch('/api/v1/agent-groups/bulk-operation', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                group_id: groupId,
                operation: operation,
                params: params,
                timeout: timeout
            })
        });

        if (response.ok) {
            const result = await response.json();
            displayBulkResults(result);
        } else {
            const error = await response.json();
            showToast('Failed to execute operation: ' + error.error, 'error');
        }
    } catch (error) {
        showToast('Failed to execute operation: ' + error.message, 'error');
    }
}

function displayBulkResults(result) {
    const container = document.getElementById('bulkResultsContainer');
    const resultsDiv = document.getElementById('bulkResults');

    const successRate = (result.success_count / result.total_agents * 100).toFixed(2);

    resultsDiv.innerHTML = `
        <div class="alert alert-info">
            <strong>Summary:</strong> ${result.success_count}/${result.total_agents} agents succeeded (${successRate}%)
            <br>
            <strong>Duration:</strong> ${result.duration_ms}ms
        </div>
        <div class="table-responsive">
            <table class="table table-sm">
                <thead>
                    <tr>
                        <th>Agent</th>
                        <th>Status</th>
                        <th>Output/Error</th>
                        <th>Duration</th>
                    </tr>
                </thead>
                <tbody>
                    ${Object.entries(result.results).map(([agent, res]) => `
                        <tr class="${res.success ? 'table-success' : 'table-danger'}">
                            <td>${escapeHtml(agent)}</td>
                            <td>
                                <i class="bi ${res.success ? 'bi-check-circle text-success' : 'bi-x-circle text-danger'}"></i>
                                ${res.success ? 'Success' : 'Failed'}
                            </td>
                            <td><small>${escapeHtml(res.output || res.error || '')}</small></td>
                            <td>${res.duration_ms}ms</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;

    container.style.display = 'block';
    showToast('Bulk operation completed', 'success');
}

// Utility function
function escapeHtml(text) {
    if (text === null || text === undefined) return '';
    const div = document.createElement('div');
    div.textContent = text.toString();
    return div.innerHTML;
}
