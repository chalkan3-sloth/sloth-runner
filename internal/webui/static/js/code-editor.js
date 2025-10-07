/* ===================================
   Code Editor with Syntax Highlighting
   For Sloth DSL and YAML files
   =================================== */

// Prevent duplicate loading
if (typeof window.SLOTH_CODE_EDITOR_LOADED !== 'undefined') {
    console.warn('Code editor already loaded, skipping...');
} else {
    window.SLOTH_CODE_EDITOR_LOADED = true;

class CodeEditor {
    constructor(element, options = {}) {
        this.element = typeof element === 'string'
            ? document.querySelector(element)
            : element;

        if (!this.element) {
            console.warn('Code editor element not found');
            return;
        }

        this.options = {
            language: options.language || 'yaml',
            theme: options.theme || 'sloth',
            lineNumbers: options.lineNumbers !== false,
            readOnly: options.readOnly || false,
            onChange: options.onChange || null,
            onSave: options.onSave || null,
            value: options.value || '',
            ...options
        };

        this.init();
    }

    init() {
        this.createEditor();
        this.setupSyntaxHighlighting();
        this.attachEventListeners();
        if (this.options.value) {
            this.setValue(this.options.value);
        }
    }

    createEditor() {
        this.element.classList.add('code-editor-container');

        this.element.innerHTML = `
            <div class="code-editor-header">
                <div class="code-editor-info">
                    <span class="code-editor-language">${this.options.language.toUpperCase()}</span>
                    <span class="code-editor-lines">0 lines</span>
                </div>
                <div class="code-editor-actions">
                    <button class="btn btn-sm btn-link" onclick="this.closest('.code-editor-container').__instance.format()" title="Format (Shift+Alt+F)">
                        <i class="bi bi-code-square"></i>
                    </button>
                    <button class="btn btn-sm btn-link" onclick="this.closest('.code-editor-container').__instance.copyToClipboard()" title="Copy (Ctrl+C)">
                        <i class="bi bi-clipboard"></i>
                    </button>
                    <button class="btn btn-sm btn-link" onclick="this.closest('.code-editor-container').__instance.download()" title="Download">
                        <i class="bi bi-download"></i>
                    </button>
                    ${this.options.onSave ? `
                        <button class="btn btn-sm btn-primary" onclick="this.closest('.code-editor-container').__instance.save()" title="Save (Ctrl+S)">
                            <i class="bi bi-check-circle"></i> Save
                        </button>
                    ` : ''}
                </div>
            </div>
            <div class="code-editor-body">
                ${this.options.lineNumbers ? '<div class="code-editor-line-numbers"></div>' : ''}
                <div class="code-editor-content">
                    <pre class="code-editor-pre"><code class="language-${this.options.language}"></code></pre>
                    <textarea class="code-editor-textarea"
                              ${this.options.readOnly ? 'readonly' : ''}
                              spellcheck="false"></textarea>
                </div>
            </div>
            <div class="code-editor-footer">
                <span class="code-editor-cursor">Ln 1, Col 1</span>
                <span class="code-editor-status"></span>
            </div>
        `;

        this.textarea = this.element.querySelector('.code-editor-textarea');
        this.pre = this.element.querySelector('.code-editor-pre');
        this.code = this.element.querySelector('code');
        this.lineNumbers = this.element.querySelector('.code-editor-line-numbers');
        this.linesCount = this.element.querySelector('.code-editor-lines');
        this.cursorInfo = this.element.querySelector('.code-editor-cursor');

        // Store instance reference
        this.element.__instance = this;
    }

    setupSyntaxHighlighting() {
        // Simple syntax highlighting for YAML/Sloth DSL
        this.highlightPatterns = {
            yaml: [
                { pattern: /^(\s*)([a-zA-Z_-]+):/gm, className: 'yaml-key' },
                { pattern: /"([^"\\]*(\\.[^"\\]*)*)"/g, className: 'yaml-string' },
                { pattern: /'([^'\\]*(\\.[^'\\]*)*)'/g, className: 'yaml-string' },
                { pattern: /\b(true|false|null|yes|no)\b/g, className: 'yaml-boolean' },
                { pattern: /\b\d+\b/g, className: 'yaml-number' },
                { pattern: /#.*/g, className: 'yaml-comment' },
                { pattern: /\${([^}]+)}/g, className: 'yaml-variable' },
                { pattern: /\b(name|description|tasks|steps|commands|type|target|agent)\b/g, className: 'yaml-keyword' }
            ]
        };
    }

    attachEventListeners() {
        this.textarea.addEventListener('input', () => this.handleInput());
        this.textarea.addEventListener('keydown', (e) => this.handleKeydown(e));
        this.textarea.addEventListener('scroll', () => this.syncScroll());
        this.textarea.addEventListener('click', () => this.updateCursor());
        this.textarea.addEventListener('keyup', () => this.updateCursor());

        // Sync scrolling
        this.syncScroll();
    }

    handleInput() {
        this.highlight();
        this.updateLineNumbers();
        this.updateLinesCount();

        if (this.options.onChange) {
            this.options.onChange(this.getValue());
        }
    }

    handleKeydown(e) {
        // Tab key
        if (e.key === 'Tab') {
            e.preventDefault();
            this.insertTab();
        }

        // Ctrl+S to save
        if ((e.ctrlKey || e.metaKey) && e.key === 's') {
            e.preventDefault();
            this.save();
        }

        // Shift+Alt+F to format
        if (e.shiftKey && e.altKey && e.key === 'f') {
            e.preventDefault();
            this.format();
        }

        // Ctrl+/ to toggle comment
        if ((e.ctrlKey || e.metaKey) && e.key === '/') {
            e.preventDefault();
            this.toggleComment();
        }

        // Auto-indent on Enter
        if (e.key === 'Enter') {
            e.preventDefault();
            this.autoIndent();
        }
    }

    highlight() {
        const value = this.textarea.value;
        let highlighted = this.escapeHtml(value);

        // Apply syntax highlighting
        const patterns = this.highlightPatterns[this.options.language] || [];
        patterns.forEach(({ pattern, className }) => {
            highlighted = highlighted.replace(pattern, (match) => {
                return `<span class="${className}">${match}</span>`;
            });
        });

        this.code.innerHTML = highlighted;
    }

    updateLineNumbers() {
        if (!this.lineNumbers) return;

        const lines = this.textarea.value.split('\n').length;
        this.lineNumbers.innerHTML = Array.from({ length: lines }, (_, i) =>
            `<div class="line-number">${i + 1}</div>`
        ).join('');
    }

    updateLinesCount() {
        const lines = this.textarea.value.split('\n').length;
        this.linesCount.textContent = `${lines} lines`;
    }

    updateCursor() {
        const pos = this.textarea.selectionStart;
        const value = this.textarea.value.substring(0, pos);
        const lines = value.split('\n');
        const line = lines.length;
        const col = lines[lines.length - 1].length + 1;

        this.cursorInfo.textContent = `Ln ${line}, Col ${col}`;
    }

    syncScroll() {
        this.pre.scrollTop = this.textarea.scrollTop;
        this.pre.scrollLeft = this.textarea.scrollLeft;

        if (this.lineNumbers) {
            this.lineNumbers.scrollTop = this.textarea.scrollTop;
        }
    }

    insertTab() {
        const start = this.textarea.selectionStart;
        const end = this.textarea.selectionEnd;
        const value = this.textarea.value;

        this.textarea.value = value.substring(0, start) + '  ' + value.substring(end);
        this.textarea.selectionStart = this.textarea.selectionEnd = start + 2;

        this.handleInput();
    }

    autoIndent() {
        const start = this.textarea.selectionStart;
        const value = this.textarea.value;
        const lines = value.substring(0, start).split('\n');
        const currentLine = lines[lines.length - 1];

        // Get indentation of current line
        const indent = currentLine.match(/^\s*/)[0];

        // Insert newline with same indentation
        this.textarea.value = value.substring(0, start) + '\n' + indent + value.substring(start);
        this.textarea.selectionStart = this.textarea.selectionEnd = start + indent.length + 1;

        this.handleInput();
    }

    toggleComment() {
        const start = this.textarea.selectionStart;
        const value = this.textarea.value;
        const lines = value.split('\n');
        const currentLineIndex = value.substring(0, start).split('\n').length - 1;
        const currentLine = lines[currentLineIndex];

        if (currentLine.trim().startsWith('#')) {
            // Uncomment
            lines[currentLineIndex] = currentLine.replace(/^\s*#\s?/, '');
        } else {
            // Comment
            const indent = currentLine.match(/^\s*/)[0];
            lines[currentLineIndex] = indent + '# ' + currentLine.trim();
        }

        this.textarea.value = lines.join('\n');
        this.handleInput();
    }

    format() {
        // Simple YAML formatting
        const lines = this.textarea.value.split('\n');
        let formatted = [];
        let indent = 0;

        lines.forEach(line => {
            const trimmed = line.trim();
            if (!trimmed) {
                formatted.push('');
                return;
            }

            // Decrease indent for closing brackets
            if (trimmed.startsWith('-') && indent > 0) {
                // Keep same indent for list items
            } else if (line.match(/^\s*\w+:/)) {
                // Key: value
                formatted.push('  '.repeat(indent) + trimmed);
            } else {
                formatted.push('  '.repeat(indent) + trimmed);
            }
        });

        this.setValue(formatted.join('\n'));

        if (window.toastManager) {
            toastManager.success('Code formatted!', null, { duration: 2000 });
        }
    }

    async copyToClipboard() {
        if (window.ClipboardManager) {
            await ClipboardManager.copy(this.getValue(), 'Code copied to clipboard!');
        }
    }

    download() {
        const content = this.getValue();
        const filename = `workflow_${Date.now()}.${this.options.language === 'yaml' ? 'yaml' : 'sloth'}`;

        if (window.downloadFile) {
            downloadFile(content, filename, 'text/plain');
        }
    }

    save() {
        if (this.options.onSave) {
            const status = this.element.querySelector('.code-editor-status');
            status.innerHTML = '<i class="bi bi-hourglass-split"></i> Saving...';

            Promise.resolve(this.options.onSave(this.getValue()))
                .then(() => {
                    status.innerHTML = '<i class="bi bi-check-circle text-success"></i> Saved';
                    setTimeout(() => status.innerHTML = '', 2000);
                })
                .catch((error) => {
                    status.innerHTML = '<i class="bi bi-x-circle text-danger"></i> Error';
                    console.error('Save error:', error);
                });
        }
    }

    getValue() {
        return this.textarea.value;
    }

    setValue(value) {
        this.textarea.value = value;
        this.highlight();
        this.updateLineNumbers();
        this.updateLinesCount();
        this.updateCursor();
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Export
window.CodeEditor = CodeEditor;

// Helper function to create code editor
window.createCodeEditor = function(selector, options) {
    return new CodeEditor(selector, options);
};

} // End of SLOTH_CODE_EDITOR_LOADED check
