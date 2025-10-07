// Enhanced Workflows Management with Charts, Code Editor and Execution Interface
let currentSloths = [];
let currentFilter = 'all';
let executionTrendsChart = null;
let statusDistributionChart = null;
let currentExecutionId = null;
let workflowCodeEditor = null;

// API Helper
const API = {
    async get(url) {
        const response = await fetch(url);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    },

    async post(url, data = {}) {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    },

    async put(url, data = {}) {
        const response = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();
    },

    async delete(url) {
        const response = await fetch(url, { method: 'DELETE' });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.ok;
    }
};

// ============= INITIALIZATION =============

document.addEventListener('DOMContentLoaded', async () => {
    await initializeCharts();
    await loadWorkflows();
    await loadStatistics();

    // Initialize Code Editor
    if (typeof CodeEditor !== 'undefined') {
        workflowCodeEditor = new CodeEditor('#workflow-code-editor', {
            language: 'yaml',
            theme: 'sloth',
            lineNumbers: true,
            readOnly: false,
            value: getBasicTemplate(),
            onChange: (value) => {
                // Optional: handle changes
            }
        });
    }

    // Setup modal event listeners
    document.getElementById('addWorkflowModal').addEventListener('shown.bs.modal', () => {
        if (workflowCodeEditor && !workflowCodeEditor.getValue()) {
            workflowCodeEditor.setValue(getBasicTemplate());
        }
    });

    document.getElementById('runWorkflowModal').addEventListener('shown.bs.modal', () => {
        loadAgentSelectionList();
    });
});

// ============= CHARTS INITIALIZATION =============

function initializeCharts() {
    // Execution Trends Chart
    const trendsCtx = document.getElementById('executionTrendsChart');
    if (trendsCtx) {
        executionTrendsChart = new Chart(trendsCtx, {
            type: 'line',
            data: {
                labels: getLast7Days(),
                datasets: [
                    {
                        label: 'Successful',
                        data: [0, 0, 0, 0, 0, 0, 0],
                        borderColor: 'rgb(75, 192, 192)',
                        backgroundColor: 'rgba(75, 192, 192, 0.1)',
                        tension: 0.4
                    },
                    {
                        label: 'Failed',
                        data: [0, 0, 0, 0, 0, 0, 0],
                        borderColor: 'rgb(255, 99, 132)',
                        backgroundColor: 'rgba(255, 99, 132, 0.1)',
                        tension: 0.4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true,
                        position: 'top'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            precision: 0
                        }
                    }
                }
            }
        });
    }

    // Status Distribution Chart
    const statusCtx = document.getElementById('statusDistributionChart');
    if (statusCtx) {
        statusDistributionChart = new Chart(statusCtx, {
            type: 'doughnut',
            data: {
                labels: ['Success', 'Failed', 'Running', 'Pending'],
                datasets: [{
                    data: [0, 0, 0, 0],
                    backgroundColor: [
                        'rgba(75, 192, 192, 0.8)',
                        'rgba(255, 99, 132, 0.8)',
                        'rgba(54, 162, 235, 0.8)',
                        'rgba(255, 206, 86, 0.8)'
                    ],
                    borderWidth: 2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    }
}

function getLast7Days() {
    const days = [];
    for (let i = 6; i >= 0; i--) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        days.push(date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' }));
    }
    return days;
}

// ============= STATISTICS =============

async function loadStatistics() {
    try {
        // Load execution history for statistics
        const executions = await API.get('/api/v1/executions').catch(() => ({ executions: [] }));
        const executionsList = executions.executions || [];

        // Calculate statistics
        const totalWorkflows = currentSloths.length;
        const activeWorkflows = currentSloths.filter(s => s.is_active).length;
        const totalExecutions = executionsList.length;

        // Success rate (last 30 days)
        const thirtyDaysAgo = new Date();
        thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
        const recentExecutions = executionsList.filter(e => new Date(e.created_at) > thirtyDaysAgo);
        const successfulExecutions = recentExecutions.filter(e => e.status === 'completed' || e.status === 'success').length;
        const successRate = recentExecutions.length > 0 ? ((successfulExecutions / recentExecutions.length) * 100).toFixed(1) : 0;

        // Update statistics cards
        document.getElementById('stat-total-workflows').textContent = totalWorkflows;
        document.getElementById('stat-active-workflows').textContent = activeWorkflows;
        document.getElementById('stat-total-executions').textContent = totalExecutions;
        document.getElementById('stat-success-rate').textContent = `${successRate}%`;

        // Update charts with real data
        updateChartsWithExecutionData(executionsList);

    } catch (error) {
        console.error('Failed to load statistics:', error);
    }
}

function updateChartsWithExecutionData(executions) {
    // Get last 7 days execution data
    const last7Days = [];
    for (let i = 6; i >= 0; i--) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        date.setHours(0, 0, 0, 0);
        last7Days.push(date);
    }

    const successData = new Array(7).fill(0);
    const failedData = new Array(7).fill(0);

    executions.forEach(exec => {
        const execDate = new Date(exec.created_at);
        execDate.setHours(0, 0, 0, 0);

        const dayIndex = last7Days.findIndex(d => d.getTime() === execDate.getTime());
        if (dayIndex >= 0) {
            if (exec.status === 'completed' || exec.status === 'success') {
                successData[dayIndex]++;
            } else if (exec.status === 'failed' || exec.status === 'error') {
                failedData[dayIndex]++;
            }
        }
    });

    if (executionTrendsChart) {
        executionTrendsChart.data.datasets[0].data = successData;
        executionTrendsChart.data.datasets[1].data = failedData;
        executionTrendsChart.update();
    }

    // Update status distribution
    const statusCounts = {
        success: executions.filter(e => e.status === 'completed' || e.status === 'success').length,
        failed: executions.filter(e => e.status === 'failed' || e.status === 'error').length,
        running: executions.filter(e => e.status === 'running').length,
        pending: executions.filter(e => e.status === 'pending').length
    };

    if (statusDistributionChart) {
        statusDistributionChart.data.datasets[0].data = [
            statusCounts.success,
            statusCounts.failed,
            statusCounts.running,
            statusCounts.pending
        ];
        statusDistributionChart.update();
    }
}

// ============= WORKFLOWS MANAGEMENT =============

async function loadWorkflows() {
    try {
        const data = await API.get('/api/v1/sloths');
        currentSloths = data.sloths || [];
        renderWorkflows();
        await loadStatistics(); // Refresh statistics after loading workflows
    } catch (error) {
        console.error('Failed to load workflows:', error);
        showError('Failed to load workflows');
    }
}

function filterWorkflows(filter) {
    currentFilter = filter;

    // Update filter buttons
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');

    renderWorkflows();
}

function searchWorkflows() {
    renderWorkflows();
}

function renderWorkflows() {
    const container = document.getElementById('workflows-container');
    const searchTerm = document.getElementById('searchWorkflows')?.value.toLowerCase() || '';

    let filtered = currentSloths;

    // Apply status filter
    if (currentFilter === 'active') {
        filtered = filtered.filter(s => s.is_active);
    } else if (currentFilter === 'inactive') {
        filtered = filtered.filter(s => !s.is_active);
    }

    // Apply search filter
    if (searchTerm) {
        filtered = filtered.filter(s =>
            s.name.toLowerCase().includes(searchTerm) ||
            (s.path && s.path.toLowerCase().includes(searchTerm))
        );
    }

    if (filtered.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No workflows found.
                    <a href="#" onclick="showAddWorkflowModal()" class="alert-link">Create one now</a>
                </div>
            </div>
        `;
        return;
    }

    container.innerHTML = filtered.map(sloth => `
        <div class="col-md-6 col-lg-4 mb-3">
            <div class="card workflow-card h-100 shadow-sm">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-3">
                        <div class="flex-grow-1">
                            <h5 class="card-title mb-2">
                                <i class="bi bi-diagram-3 text-primary"></i> ${sloth.name}
                            </h5>
                            <span class="badge ${sloth.is_active ? 'bg-success' : 'bg-secondary'} mb-2">
                                ${sloth.is_active ? 'Active' : 'Inactive'}
                            </span>
                        </div>
                        <div class="dropdown">
                            <button class="btn btn-sm btn-outline-secondary" data-bs-toggle="dropdown">
                                <i class="bi bi-three-dots-vertical"></i>
                            </button>
                            <ul class="dropdown-menu dropdown-menu-end">
                                <li><a class="dropdown-item" href="#" onclick="viewWorkflow('${sloth.name}'); return false;">
                                    <i class="bi bi-eye"></i> View Details
                                </a></li>
                                <li><a class="dropdown-item" href="#" onclick="editWorkflow('${sloth.name}'); return false;">
                                    <i class="bi bi-pencil"></i> Edit
                                </a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item" href="#" onclick="toggleWorkflowStatus('${sloth.name}', ${sloth.is_active}); return false;">
                                    <i class="bi bi-${sloth.is_active ? 'pause' : 'play'}-circle"></i>
                                    ${sloth.is_active ? 'Deactivate' : 'Activate'}
                                </a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item text-danger" href="#" onclick="deleteWorkflow('${sloth.name}'); return false;">
                                    <i class="bi bi-trash"></i> Delete
                                </a></li>
                            </ul>
                        </div>
                    </div>

                    <p class="text-muted small mb-3">
                        <i class="bi bi-folder"></i> ${sloth.path || 'No path specified'}
                    </p>

                    <div class="d-flex gap-2 mb-2">
                        <button class="btn btn-success btn-sm flex-grow-1" onclick="runWorkflow('${sloth.name}')">
                            <i class="bi bi-play-fill"></i> Execute
                        </button>
                        <button class="btn btn-outline-info btn-sm" onclick="viewWorkflow('${sloth.name}')">
                            <i class="bi bi-eye"></i>
                        </button>
                    </div>
                </div>
                <div class="card-footer bg-transparent border-top-0 pt-0">
                    <small class="text-muted">
                        <i class="bi bi-clock"></i>
                        ${sloth.created_at ? new Date(sloth.created_at).toLocaleDateString() : 'Unknown'}
                    </small>
                </div>
            </div>
        </div>
    `).join('');
}

// ============= CODE EDITOR FUNCTIONS =============

function showAddWorkflowModal() {
    const modal = new bootstrap.Modal(document.getElementById('addWorkflowModal'));
    document.getElementById('addWorkflowForm').reset();
    document.getElementById('modal-title-text').textContent = 'Create New Workflow';

    // Set basic template in code editor
    if (workflowCodeEditor) {
        workflowCodeEditor.setValue(getBasicTemplate());
    }

    modal.show();
}

function getBasicTemplate() {
    return `-- Sloth Runner Workflow
-- Modern DSL with method chaining

local hello_task = task("hello")
    :description("Simple hello world task")
    :command(function(this, params)
        log.info("🦥 Hello from Sloth Runner!")
        shell.execute("echo 'Running task on agent'")
        return true, "Success"
    end)
    :timeout("60s")
    :build()

workflow
    .define("my_workflow")
    :description("My first workflow")
    :version("1.0.0")
    :tasks({ hello_task })
    :config({
        timeout = "5m",
        max_parallel_tasks = 1,
    })
    :on_complete(function(success, results)
        if success then
            log.info("✅ Workflow completed!")
        else
            log.error("❌ Workflow failed")
        end
    end)
`;
}

function getAdvancedTemplate() {
    return `-- Advanced Sloth Runner Workflow
-- Features: Multiple tasks, dependencies, package management

local update_packages = task("update_packages")
    :description("Atualiza cache de pacotes do sistema")
    :command(function(this, params)
        log.info("🔄 Atualizando cache de pacotes...")

        local success, output = pkg.update()

        if success then
            log.info("✅ Cache atualizado com sucesso!")
            return true, "Cache atualizado"
        else
            log.error("❌ Erro: " .. output)
            return false, "Falha"
        end
    end)
    :timeout("120s")
    :build()

local install_tools = task("install_tools")
    :description("Instala ferramentas de desenvolvimento")
    :command(function(this, params)
        log.info("🛠️  Instalando ferramentas...")

        local tools = { "git", "curl", "wget", "vim" }
        local success, output = pkg.install(tools)

        if success then
            log.info("✅ Ferramentas instaladas!")
            for _, tool in ipairs(tools) do
                log.info("  ✓ " .. tool)
            end
            return true, "Instalado"
        else
            log.warn("⚠️  " .. output)
            return true, "OK"
        end
    end)
    :depends_on({ "update_packages" })
    :timeout("300s")
    :build()

local verify_installation = task("verify_installation")
    :description("Verifica instalação das ferramentas")
    :command(function(this, params)
        log.info("🔍 Verificando instalação...")

        local tools = { "git", "curl" }
        for _, tool in ipairs(tools) do
            local success, output = pkg.info(tool)
            if success then
                log.info("✓ " .. tool .. " instalado")
            else
                log.warn("⚠️  " .. tool .. " não encontrado")
            end
        end

        return true, "Verificado"
    end)
    :depends_on({ "install_tools" })
    :timeout("60s")
    :build()

workflow
    .define("setup_environment")
    :description("Configura ambiente de desenvolvimento")
    :version("1.0.0")
    :tasks({
        update_packages,
        install_tools,
        verify_installation,
    })
    :config({
        timeout = "15m",
        max_parallel_tasks = 1,
    })
    :on_complete(function(success, results)
        if success then
            log.info("🎉 Ambiente configurado com sucesso!")
        else
            log.error("❌ Falha na configuração do ambiente")
        end
    end)
`;
}

function insertTemplate(type) {
    if (!workflowCodeEditor) {
        console.error('Code editor not initialized');
        return;
    }

    const template = type === 'basic' ? getBasicTemplate() : getAdvancedTemplate();
    workflowCodeEditor.setValue(template);

    if (window.toastManager) {
        toastManager.success(`${type === 'basic' ? 'Basic' : 'Advanced'} template loaded`);
    }
}

function formatCode() {
    if (!workflowCodeEditor) {
        console.error('Code editor not initialized');
        return;
    }

    // Use the code editor's built-in format function
    workflowCodeEditor.format();
}

async function validateWorkflow() {
    if (!workflowCodeEditor) {
        if (window.toastManager) {
            toastManager.error('Code editor not initialized');
        }
        return;
    }

    const code = workflowCodeEditor.getValue();

    // Basic YAML/Sloth DSL validation
    try {
        const issues = [];

        // Check for basic YAML structure
        const lines = code.split('\n');
        let indentStack = [0];

        for (let i = 0; i < lines.length; i++) {
            const line = lines[i];
            if (!line.trim() || line.trim().startsWith('#')) continue;

            const indent = line.match(/^\s*/)[0].length;

            // Check for valid indentation (must be multiple of 2)
            if (indent % 2 !== 0) {
                issues.push(`Line ${i + 1}: Invalid indentation (must be multiple of 2 spaces)`);
            }
        }

        if (issues.length > 0) {
            if (window.toastManager) {
                toastManager.warning('Validation issues found:\n' + issues.join('\n'));
            }
        } else {
            if (window.toastManager) {
                toastManager.success('Workflow syntax is valid ✓');
            }
        }
    } catch (error) {
        if (window.toastManager) {
            toastManager.error('Validation error: ' + error.message);
        }
    }
}

async function saveWorkflow() {
    const name = document.getElementById('addWorkflowName').value;
    const path = document.getElementById('addWorkflowPath').value;
    const description = document.getElementById('addWorkflowDescription').value;
    const tags = document.getElementById('addWorkflowTags').value;
    const isActive = document.getElementById('addWorkflowActive').checked;

    if (!name || !path) {
        if (window.toastManager) {
            toastManager.error('Name and path are required');
        }
        return;
    }

    if (!workflowCodeEditor) {
        if (window.toastManager) {
            toastManager.error('Code editor not initialized');
        }
        return;
    }

    const content = workflowCodeEditor.getValue();

    try {
        await API.post('/api/v1/sloths', {
            name,
            path,
            content,
            description,
            tags: tags ? tags.split(',').map(t => t.trim()) : [],
            is_active: isActive
        });

        bootstrap.Modal.getInstance(document.getElementById('addWorkflowModal')).hide();
        showSuccess('Workflow saved successfully');
        loadWorkflows();
    } catch (error) {
        console.error('Failed to save workflow:', error);
        showError('Failed to save workflow: ' + error.message);
    }
}

// ============= WORKFLOW EXECUTION =============

async function runWorkflow(name) {
    try {
        const sloth = currentSloths.find(s => s.name === name);
        if (!sloth) return;

        const modal = new bootstrap.Modal(document.getElementById('runWorkflowModal'));
        document.getElementById('runWorkflowName').value = name;
        document.getElementById('runWorkflowPath').textContent = sloth.path;

        // Reset form
        document.getElementById('runWorkflowGroup').value = '';
        document.getElementById('runWorkflowEnv').value = '';
        document.getElementById('runWorkflowMode').value = 'normal';
        document.getElementById('runWorkflowTimeout').value = '3600';
        document.getElementById('runWorkflowRetries').value = '0';
        document.getElementById('runWorkflowVerbose').checked = false;
        document.getElementById('runWorkflowNotify').checked = true;

        // Hide execution progress
        document.getElementById('execution-progress').classList.add('d-none');

        modal.show();
    } catch (error) {
        console.error('Failed to prepare workflow run:', error);
        showError('Failed to prepare workflow run');
    }
}

async function loadAgentSelectionList() {
    const container = document.getElementById('agent-selection-list');

    try {
        const data = await API.get('/api/v1/agents');
        const agents = data.agents || [];

        if (agents.length === 0) {
            container.innerHTML = `
                <div class="text-center text-muted py-3">
                    <i class="bi bi-info-circle"></i>
                    <small class="d-block mt-2">No agents available</small>
                </div>
            `;
            return;
        }

        container.innerHTML = agents.map(agent => `
            <div class="form-check mb-2">
                <input class="form-check-input agent-checkbox" type="checkbox" value="${agent.name}" id="agent-${agent.name}">
                <label class="form-check-label d-flex justify-content-between" for="agent-${agent.name}">
                    <span>
                        <i class="bi bi-hdd-network"></i> ${agent.name}
                    </span>
                    <span class="badge ${agent.status === 'active' ? 'bg-success' : 'bg-secondary'} ms-2">
                        ${agent.status || 'offline'}
                    </span>
                </label>
            </div>
        `).join('');

    } catch (error) {
        console.error('Failed to load agents:', error);
        container.innerHTML = `
            <div class="text-center text-danger py-3">
                <i class="bi bi-exclamation-triangle"></i>
                <small class="d-block mt-2">Failed to load agents</small>
            </div>
        `;
    }
}

async function executeWorkflow() {
    const name = document.getElementById('runWorkflowName').value;
    const group = document.getElementById('runWorkflowGroup').value;
    const mode = document.getElementById('runWorkflowMode').value;
    const envText = document.getElementById('runWorkflowEnv').value;
    const timeout = parseInt(document.getElementById('runWorkflowTimeout').value);
    const retries = parseInt(document.getElementById('runWorkflowRetries').value);
    const verbose = document.getElementById('runWorkflowVerbose').checked;
    const notify = document.getElementById('runWorkflowNotify').checked;

    // Get selected agents
    const selectedAgents = Array.from(document.querySelectorAll('.agent-checkbox:checked'))
        .map(cb => cb.value);

    // Parse environment variables
    const env = {};
    if (envText) {
        envText.split('\n').forEach(line => {
            const [key, ...valueParts] = line.split('=');
            if (key && valueParts.length > 0) {
                env[key.trim()] = valueParts.join('=').trim();
            }
        });
    }

    try {
        // Show progress
        const progressDiv = document.getElementById('execution-progress');
        const progressBar = document.getElementById('execution-progress-bar');
        const statusDiv = document.getElementById('execution-status');

        progressDiv.classList.remove('d-none');
        progressBar.style.width = '10%';
        statusDiv.textContent = 'Preparing execution...';

        const payload = {
            sloth_name: name,
            group: group || undefined,
            agents: selectedAgents.length > 0 ? selectedAgents : undefined,
            mode: mode,
            env: Object.keys(env).length > 0 ? env : undefined,
            timeout: timeout,
            max_retries: retries,
            verbose: verbose,
            notify_on_complete: notify
        };

        progressBar.style.width = '30%';
        statusDiv.textContent = 'Submitting workflow...';

        const result = await API.post('/api/v1/executions', payload);
        currentExecutionId = result.execution_id;

        progressBar.style.width = '100%';
        statusDiv.textContent = 'Workflow started successfully!';

        // Wait a moment then close modal and redirect
        setTimeout(() => {
            bootstrap.Modal.getInstance(document.getElementById('runWorkflowModal')).hide();
            showSuccess(`Workflow "${name}" started successfully`);

            setTimeout(() => {
                window.location.href = `/executions?id=${currentExecutionId}`;
            }, 1000);
        }, 1000);

    } catch (error) {
        console.error('Failed to execute workflow:', error);
        showError('Failed to execute workflow: ' + error.message);
        document.getElementById('execution-progress').classList.add('d-none');
    }
}

function cancelExecution() {
    if (currentExecutionId) {
        // TODO: Implement execution cancellation API
        showError('Execution cancellation not yet implemented');
    }
}

function scheduleWorkflow() {
    showError('Workflow scheduling not yet implemented. Use the Scheduler page.');
}

// ============= VIEW & EDIT =============

async function viewWorkflow(name) {
    try {
        const sloth = await API.get(`/api/v1/sloths/${name}`);

        const modal = new bootstrap.Modal(document.getElementById('viewWorkflowModal'));
        document.getElementById('viewWorkflowContent').innerHTML = `
            <div class="mb-3">
                <h5><i class="bi bi-diagram-3"></i> ${sloth.name}</h5>
                <div class="d-flex gap-2 mb-3">
                    <span class="badge ${sloth.is_active ? 'bg-success' : 'bg-secondary'}">
                        ${sloth.is_active ? 'Active' : 'Inactive'}
                    </span>
                </div>
            </div>

            <div class="row mb-3">
                <div class="col-md-6">
                    <strong>Path:</strong><br>
                    <code class="text-muted">${sloth.path}</code>
                </div>
                <div class="col-md-6">
                    <strong>Created:</strong><br>
                    <span class="text-muted">${new Date(sloth.created_at).toLocaleString()}</span>
                </div>
            </div>

            <hr>

            <h6 class="mb-2"><i class="bi bi-file-code"></i> Workflow Content:</h6>
            <pre class="bg-dark text-light p-3 rounded" style="max-height: 400px; overflow-y: auto;"><code>${escapeHtml(sloth.content || 'No content available')}</code></pre>
        `;

        modal.show();
    } catch (error) {
        console.error('Failed to load workflow:', error);
        showError('Failed to load workflow details');
    }
}

async function editWorkflow(name) {
    try {
        const sloth = await API.get(`/api/v1/sloths/${name}`);

        // Populate form
        document.getElementById('addWorkflowName').value = sloth.name;
        document.getElementById('addWorkflowPath').value = sloth.path;
        document.getElementById('addWorkflowActive').checked = sloth.is_active;
        document.getElementById('modal-title-text').textContent = 'Edit Workflow';

        // Set content in code editor
        if (workflowCodeEditor) {
            workflowCodeEditor.setValue(sloth.content || '');
        }

        const modal = new bootstrap.Modal(document.getElementById('addWorkflowModal'));
        modal.show();

    } catch (error) {
        console.error('Failed to load workflow for editing:', error);
        showError('Failed to load workflow');
    }
}

async function toggleWorkflowStatus(name, isActive) {
    try {
        const endpoint = isActive ? 'deactivate' : 'activate';
        await API.post(`/api/v1/sloths/${name}/${endpoint}`);

        showSuccess(`Workflow ${isActive ? 'deactivated' : 'activated'} successfully`);
        loadWorkflows();
    } catch (error) {
        console.error('Failed to toggle workflow status:', error);
        showError('Failed to update workflow status');
    }
}

async function deleteWorkflow(name) {
    if (!confirm(`Are you sure you want to delete workflow "${name}"?\n\nThis action cannot be undone.`)) return;

    try {
        await API.delete(`/api/v1/sloths/${name}`);
        showSuccess('Workflow deleted successfully');
        loadWorkflows();
    } catch (error) {
        console.error('Failed to delete workflow:', error);
        showError('Failed to delete workflow');
    }
}

// ============= UTILITY FUNCTIONS =============

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showSuccess(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-success border-0 position-fixed top-0 end-0 m-3';
    toast.style.zIndex = '9999';
    toast.setAttribute('role', 'alert');
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body"><i class="bi bi-check-circle me-2"></i>${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    document.body.appendChild(toast);
    const bsToast = new bootstrap.Toast(toast);
    bsToast.show();
    setTimeout(() => toast.remove(), 5000);
}

function showError(message) {
    const toast = document.createElement('div');
    toast.className = 'toast align-items-center text-white bg-danger border-0 position-fixed top-0 end-0 m-3';
    toast.style.zIndex = '9999';
    toast.setAttribute('role', 'alert');
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body"><i class="bi bi-exclamation-triangle me-2"></i>${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;
    document.body.appendChild(toast);
    const bsToast = new bootstrap.Toast(toast);
    bsToast.show();
    setTimeout(() => toast.remove(), 5000);
}
