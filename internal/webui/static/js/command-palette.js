/* ===================================
   Command Palette (Ctrl+Shift+P / Cmd+Shift+P)
   Inspirado em VS Code, Spotlight, Alfred
   =================================== */

// Prevent duplicate loading
if (typeof window.SLOTH_COMMAND_PALETTE_LOADED !== 'undefined') {
    console.warn('Command palette already loaded, skipping...');
} else {
    window.SLOTH_COMMAND_PALETTE_LOADED = true;

class CommandPalette {
    constructor() {
        this.isOpen = false;
        this.commands = new Map();
        this.recentCommands = [];
        this.maxRecent = 5;
        this.selectedIndex = 0;
        this.filteredCommands = [];

        this.init();
        this.registerDefaultCommands();
    }

    init() {
        this.createPalette();
        this.attachKeyboardShortcut();
        this.loadRecentCommands();
    }

    createPalette() {
        const palette = document.createElement('div');
        palette.id = 'command-palette';
        palette.className = 'command-palette';
        palette.innerHTML = `
            <div class="command-palette-backdrop"></div>
            <div class="command-palette-container">
                <div class="command-palette-header">
                    <i class="bi bi-terminal"></i>
                    <input type="text"
                           id="command-palette-input"
                           class="command-palette-input"
                           placeholder="Digite um comando ou pesquise... (? para ajuda)"
                           autocomplete="off"
                           spellcheck="false">
                    <button class="command-palette-close" onclick="commandPalette.close()">
                        <i class="bi bi-x-lg"></i>
                    </button>
                </div>
                <div class="command-palette-results" id="command-palette-results">
                    <div class="command-palette-empty">
                        Digite para buscar comandos...
                    </div>
                </div>
                <div class="command-palette-footer">
                    <div class="command-palette-hints">
                        <kbd>↑↓</kbd> Navegar
                        <kbd>Enter</kbd> Executar
                        <kbd>Esc</kbd> Fechar
                    </div>
                </div>
            </div>
        `;
        document.body.appendChild(palette);

        this.input = document.getElementById('command-palette-input');
        this.results = document.getElementById('command-palette-results');

        // Event listeners
        this.input.addEventListener('input', (e) => this.handleInput(e.target.value));
        this.input.addEventListener('keydown', (e) => this.handleKeyDown(e));

        palette.querySelector('.command-palette-backdrop').addEventListener('click', () => this.close());
    }

    attachKeyboardShortcut() {
        document.addEventListener('keydown', (e) => {
            // Ctrl+Shift+P or Cmd+Shift+P
            if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'P') {
                e.preventDefault();
                this.toggle();
            }
        });
    }

    registerDefaultCommands() {
        // Navigation commands
        this.register({
            id: 'nav:dashboard',
            title: 'Ir para Dashboard',
            description: 'Página principal com visão geral',
            icon: 'speedometer2',
            category: 'Navegação',
            keywords: ['home', 'inicio', 'painel'],
            action: () => window.location.href = '/'
        });

        this.register({
            id: 'nav:agents',
            title: 'Ir para Agentes',
            description: 'Gerenciar agentes remotos',
            icon: 'hdd-network',
            category: 'Navegação',
            keywords: ['agents', 'nodes', 'remote'],
            action: () => window.location.href = '/agents'
        });

        this.register({
            id: 'nav:workflows',
            title: 'Ir para Workflows',
            description: 'Gerenciar workflows e tarefas',
            icon: 'diagram-3',
            category: 'Navegação',
            keywords: ['tasks', 'jobs', 'automacao'],
            action: () => window.location.href = '/workflows'
        });

        this.register({
            id: 'nav:terminal',
            title: 'Abrir Terminal',
            description: 'Terminal web interativo',
            icon: 'terminal',
            category: 'Navegação',
            keywords: ['shell', 'console', 'ssh'],
            action: () => window.location.href = '/terminal'
        });

        this.register({
            id: 'nav:metrics',
            title: 'Ver Métricas',
            description: 'Monitoramento e métricas do sistema',
            icon: 'graph-up',
            category: 'Navegação',
            keywords: ['stats', 'monitoring', 'graficos'],
            action: () => window.location.href = '/metrics'
        });

        // Theme commands
        this.register({
            id: 'theme:toggle',
            title: 'Alternar Tema Claro/Escuro',
            description: 'Mudar entre tema claro e escuro',
            icon: 'moon-stars',
            category: 'Aparência',
            keywords: ['dark', 'light', 'mode', 'cores'],
            action: () => {
                if (window.themeManager) {
                    themeManager.toggle();
                    this.close();
                    if (window.toastManager) {
                        toastManager.success('Tema alterado!', null, { duration: 2000 });
                    }
                }
            }
        });

        // Data commands
        this.register({
            id: 'data:refresh',
            title: 'Atualizar Página',
            description: 'Recarregar dados da página atual',
            icon: 'arrow-clockwise',
            category: 'Dados',
            keywords: ['reload', 'refresh', 'atualizar'],
            action: () => {
                this.close();
                window.location.reload();
            }
        });

        this.register({
            id: 'data:export',
            title: 'Exportar Dados',
            description: 'Exportar dados visíveis em JSON',
            icon: 'download',
            category: 'Dados',
            keywords: ['export', 'download', 'save'],
            action: () => {
                this.close();
                if (window.toastManager) {
                    toastManager.info('Função de exportação será implementada em breve');
                }
            }
        });

        // Help commands
        this.register({
            id: 'help:shortcuts',
            title: 'Mostrar Atalhos de Teclado',
            description: 'Ver todos os atalhos disponíveis',
            icon: 'keyboard',
            category: 'Ajuda',
            keywords: ['shortcuts', 'hotkeys', 'atalhos'],
            action: () => {
                this.close();
                this.showShortcutsHelp();
            }
        });

        this.register({
            id: 'help:docs',
            title: 'Abrir Documentação',
            description: 'Documentação do Sloth Runner',
            icon: 'book',
            category: 'Ajuda',
            keywords: ['docs', 'help', 'manual', 'guia'],
            action: () => {
                this.close();
                window.open('/docs', '_blank');
            }
        });

        // System commands
        this.register({
            id: 'system:clear-cache',
            title: 'Limpar Cache do Navegador',
            description: 'Limpar cache e recarregar',
            icon: 'trash',
            category: 'Sistema',
            keywords: ['cache', 'clear', 'limpar'],
            action: () => {
                this.close();
                if (confirm('Isso irá recarregar a página. Continuar?')) {
                    localStorage.clear();
                    sessionStorage.clear();
                    window.location.reload(true);
                }
            }
        });

        this.register({
            id: 'system:copy-url',
            title: 'Copiar URL da Página',
            description: 'Copiar URL atual para área de transferência',
            icon: 'link-45deg',
            category: 'Sistema',
            keywords: ['copy', 'url', 'link'],
            action: async () => {
                this.close();
                if (window.ClipboardManager) {
                    await ClipboardManager.copy(window.location.href, 'URL copiada!');
                }
            }
        });

        // UI Demo commands
        this.register({
            id: 'demo:confetti',
            title: 'Lançar Confetti 🎉',
            description: 'Demonstração de animação confetti',
            icon: 'stars',
            category: 'Demo',
            keywords: ['demo', 'test', 'celebration'],
            action: () => {
                this.close();
                if (window.confetti) {
                    confetti.celebrate();
                }
            }
        });
    }

    register(command) {
        this.commands.set(command.id, command);
    }

    handleInput(query) {
        if (!query) {
            this.showRecent();
            return;
        }

        if (query === '?') {
            this.showHelp();
            return;
        }

        const lowerQuery = query.toLowerCase();
        this.filteredCommands = Array.from(this.commands.values()).filter(cmd => {
            const searchText = `${cmd.title} ${cmd.description} ${cmd.category} ${cmd.keywords?.join(' ')}`.toLowerCase();
            return searchText.includes(lowerQuery);
        });

        this.selectedIndex = 0;
        this.renderResults();
    }

    renderResults() {
        if (this.filteredCommands.length === 0) {
            this.results.innerHTML = `
                <div class="command-palette-empty">
                    <i class="bi bi-search"></i>
                    <p>Nenhum comando encontrado</p>
                </div>
            `;
            return;
        }

        // Group by category
        const grouped = {};
        this.filteredCommands.forEach(cmd => {
            const cat = cmd.category || 'Outros';
            if (!grouped[cat]) grouped[cat] = [];
            grouped[cat].push(cmd);
        });

        let html = '';
        let globalIndex = 0;

        Object.entries(grouped).forEach(([category, commands]) => {
            html += `<div class="command-category">${category}</div>`;
            commands.forEach(cmd => {
                const isSelected = globalIndex === this.selectedIndex;
                html += `
                    <div class="command-item ${isSelected ? 'selected' : ''}"
                         data-index="${globalIndex}"
                         onclick="commandPalette.executeCommand('${cmd.id}')">
                        <i class="bi bi-${cmd.icon}"></i>
                        <div class="command-info">
                            <div class="command-title">${cmd.title}</div>
                            <div class="command-description">${cmd.description}</div>
                        </div>
                        ${cmd.shortcut ? `<kbd>${cmd.shortcut}</kbd>` : ''}
                    </div>
                `;
                globalIndex++;
            });
        });

        this.results.innerHTML = html;
    }

    showRecent() {
        if (this.recentCommands.length === 0) {
            this.results.innerHTML = `
                <div class="command-palette-empty">
                    <i class="bi bi-clock-history"></i>
                    <p>Nenhum comando recente</p>
                    <small>Digite ? para ver a ajuda</small>
                </div>
            `;
            return;
        }

        let html = '<div class="command-category">Comandos Recentes</div>';
        this.recentCommands.forEach((cmdId, index) => {
            const cmd = this.commands.get(cmdId);
            if (!cmd) return;

            const isSelected = index === this.selectedIndex;
            html += `
                <div class="command-item ${isSelected ? 'selected' : ''}"
                     data-index="${index}"
                     onclick="commandPalette.executeCommand('${cmd.id}')">
                    <i class="bi bi-${cmd.icon}"></i>
                    <div class="command-info">
                        <div class="command-title">${cmd.title}</div>
                        <div class="command-description">${cmd.description}</div>
                    </div>
                </div>
            `;
        });

        this.results.innerHTML = html;
    }

    showHelp() {
        this.results.innerHTML = `
            <div class="command-palette-help">
                <h5><i class="bi bi-info-circle"></i> Ajuda do Command Palette</h5>
                <div class="help-section">
                    <h6>Navegação</h6>
                    <p>Use <kbd>↑</kbd> e <kbd>↓</kbd> para navegar pelos comandos</p>
                    <p>Pressione <kbd>Enter</kbd> para executar o comando selecionado</p>
                    <p>Pressione <kbd>Esc</kbd> para fechar</p>
                </div>
                <div class="help-section">
                    <h6>Busca</h6>
                    <p>Digite para filtrar comandos por título, descrição ou categoria</p>
                    <p>A busca é case-insensitive e busca em múltiplos campos</p>
                </div>
                <div class="help-section">
                    <h6>Categorias Disponíveis</h6>
                    <p><span class="badge bg-primary">Navegação</span> Ir para diferentes páginas</p>
                    <p><span class="badge bg-success">Aparência</span> Customizar interface</p>
                    <p><span class="badge bg-info">Dados</span> Manipular dados</p>
                    <p><span class="badge bg-warning">Sistema</span> Configurações do sistema</p>
                    <p><span class="badge bg-secondary">Ajuda</span> Documentação e ajuda</p>
                </div>
            </div>
        `;
    }

    handleKeyDown(e) {
        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                this.selectedIndex = Math.min(this.selectedIndex + 1, this.filteredCommands.length - 1);
                this.renderResults();
                this.scrollToSelected();
                break;

            case 'ArrowUp':
                e.preventDefault();
                this.selectedIndex = Math.max(this.selectedIndex - 1, 0);
                this.renderResults();
                this.scrollToSelected();
                break;

            case 'Enter':
                e.preventDefault();
                if (this.filteredCommands.length > 0) {
                    const cmd = this.filteredCommands[this.selectedIndex];
                    this.executeCommand(cmd.id);
                }
                break;

            case 'Escape':
                e.preventDefault();
                this.close();
                break;
        }
    }

    scrollToSelected() {
        const selected = this.results.querySelector('.command-item.selected');
        if (selected) {
            selected.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
        }
    }

    executeCommand(commandId) {
        const cmd = this.commands.get(commandId);
        if (!cmd) return;

        // Add to recent
        this.addToRecent(commandId);

        // Execute
        try {
            cmd.action();
        } catch (error) {
            console.error('Error executing command:', error);
            if (window.toastManager) {
                toastManager.error('Erro ao executar comando: ' + error.message);
            }
        }
    }

    addToRecent(commandId) {
        // Remove if already exists
        this.recentCommands = this.recentCommands.filter(id => id !== commandId);

        // Add to beginning
        this.recentCommands.unshift(commandId);

        // Limit size
        if (this.recentCommands.length > this.maxRecent) {
            this.recentCommands = this.recentCommands.slice(0, this.maxRecent);
        }

        // Save to localStorage
        localStorage.setItem('sloth-recent-commands', JSON.stringify(this.recentCommands));
    }

    loadRecentCommands() {
        try {
            const saved = localStorage.getItem('sloth-recent-commands');
            if (saved) {
                this.recentCommands = JSON.parse(saved);
            }
        } catch (error) {
            console.error('Error loading recent commands:', error);
        }
    }

    toggle() {
        if (this.isOpen) {
            this.close();
        } else {
            this.open();
        }
    }

    open() {
        this.isOpen = true;
        document.getElementById('command-palette').classList.add('active');
        this.input.value = '';
        this.filteredCommands = [];
        this.selectedIndex = 0;
        this.showRecent();

        setTimeout(() => this.input.focus(), 100);

        // Add body class to prevent scrolling
        document.body.style.overflow = 'hidden';
    }

    close() {
        this.isOpen = false;
        document.getElementById('command-palette').classList.remove('active');
        document.body.style.overflow = '';
    }

    showShortcutsHelp() {
        const shortcuts = [
            { key: 'Ctrl+Shift+P', description: 'Abrir Command Palette', icon: 'terminal' },
            { key: 'Ctrl+K', description: 'Busca global', icon: 'search' },
            { key: '/', description: 'Mostrar atalhos', icon: 'keyboard' },
            { key: 'Esc', description: 'Fechar modais', icon: 'x-circle' },
            { key: '?', description: 'Ajuda contextual', icon: 'question-circle' }
        ];

        const html = `
            <div class="shortcuts-help-modal">
                <h5><i class="bi bi-keyboard"></i> Atalhos de Teclado</h5>
                <div class="shortcuts-list">
                    ${shortcuts.map(s => `
                        <div class="shortcut-row">
                            <i class="bi bi-${s.icon}"></i>
                            <span class="shortcut-desc">${s.description}</span>
                            <kbd>${s.key}</kbd>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;

        if (window.toastManager) {
            toastManager.show({
                message: html,
                type: 'info',
                duration: 10000
            });
        }
    }
}

// Initialize
const commandPalette = new CommandPalette();
window.commandPalette = commandPalette;

} // End of SLOTH_COMMAND_PALETTE_LOADED check
