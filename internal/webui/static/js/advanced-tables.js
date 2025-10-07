/* ===================================
   Advanced Table Features
   - Sorting, filtering, pagination
   - Bulk operations
   - Row selection
   - Inline editing
   - Export capabilities
   =================================== */

// Prevent duplicate loading
if (typeof window.SLOTH_ADVANCED_TABLES_LOADED !== 'undefined') {
    console.warn('Advanced tables already loaded, skipping...');
} else {
    window.SLOTH_ADVANCED_TABLES_LOADED = true;

class AdvancedTable {
    constructor(tableElement, options = {}) {
        this.table = typeof tableElement === 'string'
            ? document.querySelector(tableElement)
            : tableElement;

        if (!this.table) {
            console.warn('Table element not found');
            return;
        }

        this.options = {
            selectable: options.selectable !== false,
            sortable: options.sortable !== false,
            filterable: options.filterable !== false,
            pagination: options.pagination !== false,
            pageSize: options.pageSize || 10,
            exportable: options.exportable !== false,
            bulkActions: options.bulkActions || [],
            onRowSelect: options.onRowSelect || null,
            onSort: options.onSort || null,
            onFilter: options.onFilter || null,
            ...options
        };

        this.data = [];
        this.filteredData = [];
        this.selectedRows = new Set();
        this.currentPage = 1;
        this.sortColumn = null;
        this.sortDirection = 'asc';
        this.filters = {};

        this.init();
    }

    init() {
        this.extractData();
        this.enhance();
        this.attachEventListeners();
    }

    extractData() {
        // Extract data from existing table
        const rows = this.table.querySelectorAll('tbody tr');
        this.data = Array.from(rows).map((row, index) => {
            const cells = row.querySelectorAll('td');
            const rowData = {
                _index: index,
                _element: row
            };

            cells.forEach((cell, cellIndex) => {
                const header = this.table.querySelector(`thead th:nth-child(${cellIndex + 1})`);
                const key = header?.dataset?.key || `col${cellIndex}`;
                rowData[key] = cell.textContent.trim();
                rowData[`_${key}_html`] = cell.innerHTML;
            });

            return rowData;
        });

        this.filteredData = [...this.data];
    }

    enhance() {
        // Add wrapper
        const wrapper = document.createElement('div');
        wrapper.className = 'advanced-table-wrapper';
        this.table.parentNode.insertBefore(wrapper, this.table);
        wrapper.appendChild(this.table);

        // Add controls
        this.addControls(wrapper);

        // Enhance table
        this.table.classList.add('advanced-table');

        // Add selection column if enabled
        if (this.options.selectable) {
            this.addSelectionColumn();
        }

        // Make headers sortable
        if (this.options.sortable) {
            this.makeSortable();
        }
    }

    addControls(wrapper) {
        const controls = document.createElement('div');
        controls.className = 'table-controls';
        controls.innerHTML = `
            <div class="table-controls-left">
                <div class="table-info">
                    <span class="selected-count" style="display: none;">
                        <strong>0</strong> selecionados
                    </span>
                    <span class="total-count">
                        Total: <strong>${this.data.length}</strong> itens
                    </span>
                </div>
                <div class="bulk-actions" style="display: none;">
                    ${this.options.bulkActions.map(action => `
                        <button class="btn btn-sm btn-${action.variant || 'secondary'} bulk-action-btn"
                                data-action="${action.id}">
                            <i class="bi bi-${action.icon}"></i> ${action.label}
                        </button>
                    `).join('')}
                </div>
            </div>
            <div class="table-controls-right">
                <div class="table-search">
                    <i class="bi bi-search"></i>
                    <input type="search"
                           class="form-control form-control-sm table-search-input"
                           placeholder="Buscar...">
                </div>
                <div class="table-actions">
                    <button class="btn btn-sm btn-outline-secondary" onclick="this.closest('.advanced-table-wrapper').__instance.refresh()">
                        <i class="bi bi-arrow-clockwise"></i>
                    </button>
                    ${this.options.exportable ? `
                        <div class="dropdown">
                            <button class="btn btn-sm btn-outline-secondary dropdown-toggle" data-bs-toggle="dropdown">
                                <i class="bi bi-download"></i> Exportar
                            </button>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item" href="#" onclick="this.closest('.advanced-table-wrapper').__instance.export('json')">
                                    <i class="bi bi-filetype-json"></i> JSON
                                </a></li>
                                <li><a class="dropdown-item" href="#" onclick="this.closest('.advanced-table-wrapper').__instance.export('csv')">
                                    <i class="bi bi-filetype-csv"></i> CSV
                                </a></li>
                            </ul>
                        </div>
                    ` : ''}
                </div>
            </div>
        `;

        wrapper.insertBefore(controls, this.table);

        // Store instance reference
        wrapper.__instance = this;

        // Add pagination if enabled
        if (this.options.pagination) {
            const pagination = document.createElement('div');
            pagination.className = 'table-pagination';
            wrapper.appendChild(pagination);
            this.updatePagination();
        }
    }

    addSelectionColumn() {
        // Add header checkbox
        const thead = this.table.querySelector('thead tr');
        const th = document.createElement('th');
        th.className = 'selection-cell';
        th.innerHTML = `
            <input type="checkbox" class="form-check-input select-all-checkbox">
        `;
        thead.insertBefore(th, thead.firstChild);

        // Add row checkboxes
        this.table.querySelectorAll('tbody tr').forEach(row => {
            const td = document.createElement('td');
            td.className = 'selection-cell';
            td.innerHTML = `
                <input type="checkbox" class="form-check-input row-checkbox">
            `;
            row.insertBefore(td, row.firstChild);
        });
    }

    makeSortable() {
        const headers = this.table.querySelectorAll('thead th:not(.selection-cell)');
        headers.forEach((th, index) => {
            th.classList.add('sortable');
            th.style.cursor = 'pointer';

            const sortIcon = document.createElement('i');
            sortIcon.className = 'bi bi-arrow-down-up sort-icon';
            th.appendChild(sortIcon);

            th.addEventListener('click', () => this.sort(index));
        });
    }

    attachEventListeners() {
        const wrapper = this.table.closest('.advanced-table-wrapper');

        // Search
        const searchInput = wrapper.querySelector('.table-search-input');
        if (searchInput) {
            searchInput.addEventListener('input', window.debounce((e) => {
                this.filter({ _search: e.target.value });
            }, 300));
        }

        // Select all
        const selectAll = wrapper.querySelector('.select-all-checkbox');
        if (selectAll) {
            selectAll.addEventListener('change', (e) => {
                this.toggleSelectAll(e.target.checked);
            });
        }

        // Row selection
        const rowCheckboxes = this.table.querySelectorAll('.row-checkbox');
        rowCheckboxes.forEach((checkbox, index) => {
            checkbox.addEventListener('change', (e) => {
                this.toggleRowSelection(index, e.target.checked);
            });
        });

        // Bulk actions
        const bulkActionButtons = wrapper.querySelectorAll('.bulk-action-btn');
        bulkActionButtons.forEach(btn => {
            btn.addEventListener('click', () => {
                const actionId = btn.dataset.action;
                const action = this.options.bulkActions.find(a => a.id === actionId);
                if (action && action.handler) {
                    const selectedData = Array.from(this.selectedRows).map(i => this.data[i]);
                    action.handler(selectedData, this.selectedRows);
                }
            });
        });
    }

    sort(columnIndex) {
        const th = this.table.querySelectorAll('thead th:not(.selection-cell)')[columnIndex];
        const key = th.dataset.key || `col${columnIndex}`;

        // Toggle sort direction
        if (this.sortColumn === key) {
            this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            this.sortColumn = key;
            this.sortDirection = 'asc';
        }

        // Update icons
        this.table.querySelectorAll('.sort-icon').forEach(icon => {
            icon.className = 'bi bi-arrow-down-up sort-icon';
        });

        const icon = th.querySelector('.sort-icon');
        icon.className = `bi bi-arrow-${this.sortDirection === 'asc' ? 'up' : 'down'} sort-icon`;

        // Sort data
        this.filteredData.sort((a, b) => {
            const aVal = a[key];
            const bVal = b[key];

            // Try numeric comparison
            const aNum = parseFloat(aVal);
            const bNum = parseFloat(bVal);

            if (!isNaN(aNum) && !isNaN(bNum)) {
                return this.sortDirection === 'asc' ? aNum - bNum : bNum - aNum;
            }

            // String comparison
            return this.sortDirection === 'asc'
                ? aVal.localeCompare(bVal)
                : bVal.localeCompare(aVal);
        });

        if (this.options.onSort) {
            this.options.onSort(key, this.sortDirection);
        }

        this.render();
    }

    filter(filters) {
        this.filters = { ...this.filters, ...filters };

        this.filteredData = this.data.filter(row => {
            // Search filter
            if (filters._search) {
                const searchLower = filters._search.toLowerCase();
                const matches = Object.values(row).some(val =>
                    String(val).toLowerCase().includes(searchLower)
                );
                if (!matches) return false;
            }

            // Custom filters
            for (const [key, value] of Object.entries(this.filters)) {
                if (key.startsWith('_')) continue;
                if (row[key] !== value) return false;
            }

            return true;
        });

        if (this.options.onFilter) {
            this.options.onFilter(this.filters, this.filteredData);
        }

        this.currentPage = 1;
        this.render();
        this.updatePagination();
    }

    toggleSelectAll(checked) {
        this.selectedRows.clear();

        if (checked) {
            this.filteredData.forEach(row => this.selectedRows.add(row._index));
        }

        this.table.querySelectorAll('.row-checkbox').forEach(checkbox => {
            checkbox.checked = checked;
        });

        this.updateSelectionUI();
    }

    toggleRowSelection(index, checked) {
        if (checked) {
            this.selectedRows.add(index);
        } else {
            this.selectedRows.delete(index);
        }

        this.updateSelectionUI();

        if (this.options.onRowSelect) {
            this.options.onRowSelect(Array.from(this.selectedRows), this.data[index]);
        }
    }

    updateSelectionUI() {
        const wrapper = this.table.closest('.advanced-table-wrapper');
        const selectedCount = wrapper.querySelector('.selected-count');
        const totalCount = wrapper.querySelector('.total-count');
        const bulkActions = wrapper.querySelector('.bulk-actions');

        if (this.selectedRows.size > 0) {
            selectedCount.style.display = 'inline';
            selectedCount.querySelector('strong').textContent = this.selectedRows.size;
            totalCount.style.display = 'none';
            if (bulkActions) bulkActions.style.display = 'flex';
        } else {
            selectedCount.style.display = 'none';
            totalCount.style.display = 'inline';
            if (bulkActions) bulkActions.style.display = 'none';
        }

        // Update select-all checkbox
        const selectAll = wrapper.querySelector('.select-all-checkbox');
        if (selectAll) {
            selectAll.checked = this.selectedRows.size === this.filteredData.length;
            selectAll.indeterminate = this.selectedRows.size > 0 && this.selectedRows.size < this.filteredData.length;
        }
    }

    render() {
        const tbody = this.table.querySelector('tbody');
        const startIndex = (this.currentPage - 1) * this.options.pageSize;
        const endIndex = startIndex + this.options.pageSize;
        const pageData = this.filteredData.slice(startIndex, endIndex);

        tbody.innerHTML = pageData.map(row => {
            const isSelected = this.selectedRows.has(row._index);
            const cells = Object.keys(row)
                .filter(key => !key.startsWith('_') || key.endsWith('_html'))
                .map(key => {
                    if (key.endsWith('_html')) {
                        return row[key];
                    }
                    return `<td>${row[key]}</td>`;
                })
                .join('');

            return `
                <tr>
                    ${this.options.selectable ? `
                        <td class="selection-cell">
                            <input type="checkbox" class="form-check-input row-checkbox" ${isSelected ? 'checked' : ''}>
                        </td>
                    ` : ''}
                    ${cells}
                </tr>
            `;
        }).join('');

        // Reattach event listeners
        this.attachEventListeners();
    }

    updatePagination() {
        const wrapper = this.table.closest('.advanced-table-wrapper');
        const pagination = wrapper.querySelector('.table-pagination');
        if (!pagination) return;

        const totalPages = Math.ceil(this.filteredData.length / this.options.pageSize);

        pagination.innerHTML = `
            <div class="pagination-info">
                Mostrando ${((this.currentPage - 1) * this.options.pageSize) + 1} -
                ${Math.min(this.currentPage * this.options.pageSize, this.filteredData.length)}
                de ${this.filteredData.length}
            </div>
            <div class="pagination-controls">
                <button class="btn btn-sm btn-outline-secondary"
                        ${this.currentPage === 1 ? 'disabled' : ''}
                        onclick="this.closest('.advanced-table-wrapper').__instance.prevPage()">
                    <i class="bi bi-chevron-left"></i>
                </button>
                <span class="page-number">Página ${this.currentPage} de ${totalPages}</span>
                <button class="btn btn-sm btn-outline-secondary"
                        ${this.currentPage === totalPages ? 'disabled' : ''}
                        onclick="this.closest('.advanced-table-wrapper').__instance.nextPage()">
                    <i class="bi bi-chevron-right"></i>
                </button>
            </div>
        `;
    }

    nextPage() {
        const totalPages = Math.ceil(this.filteredData.length / this.options.pageSize);
        if (this.currentPage < totalPages) {
            this.currentPage++;
            this.render();
            this.updatePagination();
        }
    }

    prevPage() {
        if (this.currentPage > 1) {
            this.currentPage--;
            this.render();
            this.updatePagination();
        }
    }

    refresh() {
        this.extractData();
        this.filter(this.filters);

        if (window.toastManager) {
            toastManager.success('Tabela atualizada!', null, { duration: 2000 });
        }
    }

    export(format) {
        const data = this.selectedRows.size > 0
            ? Array.from(this.selectedRows).map(i => this.data[i])
            : this.filteredData;

        if (format === 'json') {
            const json = JSON.stringify(data, null, 2);
            this.downloadFile(json, 'export.json', 'application/json');
        } else if (format === 'csv') {
            const csv = this.toCSV(data);
            this.downloadFile(csv, 'export.csv', 'text/csv');
        }

        if (window.toastManager) {
            toastManager.success(`Dados exportados em ${format.toUpperCase()}!`, null, { duration: 3000 });
        }
    }

    toCSV(data) {
        if (data.length === 0) return '';

        const keys = Object.keys(data[0]).filter(k => !k.startsWith('_'));
        const header = keys.join(',');
        const rows = data.map(row =>
            keys.map(key => `"${String(row[key]).replace(/"/g, '""')}"`).join(',')
        );

        return [header, ...rows].join('\n');
    }

    downloadFile(content, filename, type) {
        const blob = new Blob([content], { type });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    }
}

// Export
window.AdvancedTable = AdvancedTable;

// Helper function to initialize tables
window.initAdvancedTable = function(selector, options) {
    const elements = typeof selector === 'string'
        ? document.querySelectorAll(selector)
        : [selector];

    return Array.from(elements).map(el => new AdvancedTable(el, options));
};

} // End of SLOTH_ADVANCED_TABLES_LOADED check
