/* ===================================
   Drag & Drop File Upload
   =================================== */

// Prevent duplicate loading
if (typeof window.SLOTH_DRAGDROP_LOADED !== 'undefined') {
    console.warn('Drag-drop already loaded, skipping...');
} else {
    window.SLOTH_DRAGDROP_LOADED = true;

class DragDropUploader {
    constructor(element, options = {}) {
        this.element = typeof element === 'string' ? document.querySelector(element) : element;
        if (!this.element) {
            console.warn('Drag-drop element not found');
            return;
        }

        this.options = {
            accept: options.accept || '*/*',
            maxSize: options.maxSize || 10 * 1024 * 1024, // 10MB
            multiple: options.multiple !== false,
            onUpload: options.onUpload || null,
            onError: options.onError || null,
            uploadUrl: options.uploadUrl || '/api/v1/upload',
            ...options
        };

        this.init();
    }

    init() {
        this.createDropZone();
        this.attachEvents();
    }

    createDropZone() {
        this.element.classList.add('drag-drop-zone');

        if (!this.element.querySelector('.drag-drop-content')) {
            this.element.innerHTML = `
                <div class="drag-drop-content">
                    <div class="drag-drop-icon">
                        <i class="bi bi-cloud-arrow-up"></i>
                    </div>
                    <div class="drag-drop-text">
                        <strong>Drag & drop files here</strong>
                        <span>or click to browse</span>
                    </div>
                    <input type="file" class="drag-drop-input" ${this.options.multiple ? 'multiple' : ''} accept="${this.options.accept}">
                </div>
                <div class="drag-drop-preview"></div>
            `;
        }

        this.input = this.element.querySelector('.drag-drop-input');
        this.preview = this.element.querySelector('.drag-drop-preview');
    }

    attachEvents() {
        // Drag events
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            this.element.addEventListener(eventName, this.preventDefaults, false);
        });

        ['dragenter', 'dragover'].forEach(eventName => {
            this.element.addEventListener(eventName, () => this.highlight(), false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            this.element.addEventListener(eventName, () => this.unhighlight(), false);
        });

        this.element.addEventListener('drop', (e) => this.handleDrop(e), false);

        // Click to upload
        this.element.addEventListener('click', () => this.input.click());
        this.input.addEventListener('change', (e) => this.handleFiles(e.target.files));

        // Prevent input from being clicked when clicking on element
        this.input.addEventListener('click', (e) => e.stopPropagation());
    }

    preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    highlight() {
        this.element.classList.add('drag-over');
    }

    unhighlight() {
        this.element.classList.remove('drag-over');
    }

    handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;
        this.handleFiles(files);
    }

    handleFiles(files) {
        if (!files || files.length === 0) return;

        const fileArray = Array.from(files);

        // Validate files
        const validFiles = fileArray.filter(file => this.validateFile(file));

        if (validFiles.length === 0) return;

        // Show preview
        this.showPreview(validFiles);

        // Upload files
        if (this.options.onUpload) {
            this.options.onUpload(validFiles);
        } else {
            this.uploadFiles(validFiles);
        }
    }

    validateFile(file) {
        // Check file size
        if (file.size > this.options.maxSize) {
            const maxMB = (this.options.maxSize / 1024 / 1024).toFixed(1);
            this.showError(`File "${file.name}" is too large. Max size: ${maxMB}MB`);
            return false;
        }

        // Check file type if accept is specified
        if (this.options.accept !== '*/*') {
            const acceptTypes = this.options.accept.split(',').map(t => t.trim());
            const fileType = file.type || '';
            const fileExt = '.' + file.name.split('.').pop();

            const isAccepted = acceptTypes.some(type => {
                if (type.startsWith('.')) return fileExt === type;
                if (type.endsWith('/*')) return fileType.startsWith(type.replace('/*', ''));
                return fileType === type;
            });

            if (!isAccepted) {
                this.showError(`File type "${fileType || fileExt}" is not accepted`);
                return false;
            }
        }

        return true;
    }

    showPreview(files) {
        this.preview.innerHTML = files.map((file, index) => `
            <div class="file-preview-item" data-index="${index}">
                <div class="file-icon">
                    ${this.getFileIcon(file)}
                </div>
                <div class="file-info">
                    <div class="file-name">${file.name}</div>
                    <div class="file-size">${this.formatBytes(file.size)}</div>
                </div>
                <div class="file-progress">
                    <div class="progress">
                        <div class="progress-bar progress-bar-animated-gradient" role="progressbar" style="width: 0%"></div>
                    </div>
                </div>
                <button class="file-remove" onclick="window.dragDrop.removeFile(${index})">
                    <i class="bi bi-x-circle"></i>
                </button>
            </div>
        `).join('');

        this.preview.classList.add('show');
    }

    getFileIcon(file) {
        const type = file.type.split('/')[0];
        switch (type) {
            case 'image': return '<i class="bi bi-file-image text-primary"></i>';
            case 'video': return '<i class="bi bi-file-play text-danger"></i>';
            case 'audio': return '<i class="bi bi-file-music text-warning"></i>';
            case 'text': return '<i class="bi bi-file-text text-info"></i>';
            case 'application':
                if (file.type.includes('pdf')) return '<i class="bi bi-file-pdf text-danger"></i>';
                if (file.type.includes('zip')) return '<i class="bi bi-file-zip text-warning"></i>';
                return '<i class="bi bi-file-earmark text-secondary"></i>';
            default: return '<i class="bi bi-file-earmark text-secondary"></i>';
        }
    }

    formatBytes(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    }

    async uploadFiles(files) {
        for (let i = 0; i < files.length; i++) {
            await this.uploadFile(files[i], i);
        }
    }

    async uploadFile(file, index) {
        const formData = new FormData();
        formData.append('file', file);

        const progressBar = this.preview.querySelector(`[data-index="${index}"] .progress-bar`);

        try {
            const xhr = new XMLHttpRequest();

            // Progress tracking
            xhr.upload.addEventListener('progress', (e) => {
                if (e.lengthComputable) {
                    const percent = (e.loaded / e.total) * 100;
                    if (progressBar) {
                        progressBar.style.width = percent + '%';
                    }
                }
            });

            // Complete
            xhr.addEventListener('load', () => {
                if (xhr.status === 200) {
                    if (progressBar) {
                        progressBar.classList.remove('progress-bar-animated-gradient');
                        progressBar.classList.add('bg-success');
                        progressBar.style.width = '100%';
                    }

                    if (window.toastManager) {
                        toastManager.success(`File "${file.name}" uploaded successfully!`, null, { duration: 3000 });
                    }
                } else {
                    throw new Error(`Upload failed: ${xhr.statusText}`);
                }
            });

            // Error
            xhr.addEventListener('error', () => {
                this.showError(`Failed to upload "${file.name}"`);
                if (progressBar) {
                    progressBar.classList.remove('progress-bar-animated-gradient');
                    progressBar.classList.add('bg-danger');
                }
            });

            xhr.open('POST', this.options.uploadUrl);
            xhr.send(formData);

        } catch (error) {
            this.showError(`Failed to upload "${file.name}": ${error.message}`);
            if (progressBar) {
                progressBar.classList.remove('progress-bar-animated-gradient');
                progressBar.classList.add('bg-danger');
            }
        }
    }

    removeFile(index) {
        const item = this.preview.querySelector(`[data-index="${index}"]`);
        if (item) {
            item.style.animation = 'fadeOut 0.3s ease-out';
            setTimeout(() => item.remove(), 300);
        }
    }

    showError(message) {
        if (this.options.onError) {
            this.options.onError(message);
        } else if (window.toastManager) {
            toastManager.error(message);
        } else {
            console.error(message);
        }
    }

    reset() {
        this.preview.innerHTML = '';
        this.preview.classList.remove('show');
        this.input.value = '';
    }
}

// Export
window.DragDropUploader = DragDropUploader;

// Global instance for easy access
window.dragDrop = {
    instances: new Map(),

    init: function(selector, options) {
        const elements = document.querySelectorAll(selector);
        elements.forEach(el => {
            const instance = new DragDropUploader(el, options);
            this.instances.set(el, instance);
        });
    },

    getInstance: function(element) {
        return this.instances.get(element);
    },

    removeFile: function(index) {
        // Helper for onclick handlers
        const activeInstance = Array.from(this.instances.values())[0];
        if (activeInstance) {
            activeInstance.removeFile(index);
        }
    }
};

} // End of SLOTH_DRAGDROP_LOADED check
