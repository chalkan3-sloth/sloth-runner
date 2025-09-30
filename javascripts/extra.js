// Custom JavaScript for Sloth Runner Documentation

document.addEventListener('DOMContentLoaded', function() {
    // Add fade-in animation to main content
    const mainContent = document.querySelector('.md-content');
    if (mainContent) {
        mainContent.classList.add('fade-in');
    }

    // Add slide-in animation to navigation items
    const navItems = document.querySelectorAll('.md-nav__item');
    navItems.forEach((item, index) => {
        setTimeout(() => {
            item.classList.add('slide-in');
        }, index * 50);
    });

    // Enhanced code block interactions
    enhanceCodeBlocks();
    
    // Add copy buttons to code blocks
    addCopyButtons();
    
    // Initialize tooltips
    initializeTooltips();
    
    // Add language indicators
    addLanguageIndicators();
    
    // Initialize progress bars
    initializeProgressBars();
    
    // Add smooth scrolling
    enableSmoothScrolling();
    
    // Initialize search enhancements
    enhanceSearch();
});

// Enhance code blocks with syntax highlighting and features
function enhanceCodeBlocks() {
    const codeBlocks = document.querySelectorAll('pre code');
    codeBlocks.forEach(block => {
        // Add line numbers for longer code blocks
        if (block.textContent.split('\n').length > 5) {
            addLineNumbers(block);
        }
        
        // Add language label
        const language = getCodeLanguage(block);
        if (language) {
            addLanguageLabel(block, language);
        }
        
        // Make code blocks focusable for better accessibility
        block.setAttribute('tabindex', '0');
        
        // Add hover effects
        block.addEventListener('mouseenter', function() {
            this.style.transform = 'scale(1.02)';
            this.style.transition = 'transform 0.3s ease';
        });
        
        block.addEventListener('mouseleave', function() {
            this.style.transform = 'scale(1)';
        });
    });
}

// Add copy buttons to code blocks
function addCopyButtons() {
    const codeBlocks = document.querySelectorAll('pre code');
    codeBlocks.forEach(block => {
        const pre = block.parentElement;
        const button = document.createElement('button');
        button.className = 'copy-button';
        button.innerHTML = '📋 Copy';
        button.style.cssText = `
            position: absolute;
            top: 8px;
            right: 8px;
            background: rgba(0, 0, 0, 0.7);
            color: white;
            border: none;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            cursor: pointer;
            opacity: 0;
            transition: opacity 0.3s ease;
        `;
        
        pre.style.position = 'relative';
        pre.appendChild(button);
        
        // Show button on hover
        pre.addEventListener('mouseenter', () => {
            button.style.opacity = '1';
        });
        
        pre.addEventListener('mouseleave', () => {
            button.style.opacity = '0';
        });
        
        // Copy functionality
        button.addEventListener('click', async () => {
            try {
                await navigator.clipboard.writeText(block.textContent);
                button.innerHTML = '✅ Copied!';
                setTimeout(() => {
                    button.innerHTML = '📋 Copy';
                }, 2000);
            } catch (err) {
                // Fallback for older browsers
                const textarea = document.createElement('textarea');
                textarea.value = block.textContent;
                document.body.appendChild(textarea);
                textarea.select();
                document.execCommand('copy');
                document.body.removeChild(textarea);
                
                button.innerHTML = '✅ Copied!';
                setTimeout(() => {
                    button.innerHTML = '📋 Copy';
                }, 2000);
            }
        });
    });
}

// Add line numbers to code blocks
function addLineNumbers(codeBlock) {
    const lines = codeBlock.textContent.split('\n');
    const lineNumbersDiv = document.createElement('div');
    lineNumbersDiv.className = 'line-numbers';
    lineNumbersDiv.style.cssText = `
        position: absolute;
        left: 0;
        top: 0;
        bottom: 0;
        width: 40px;
        background: rgba(0, 0, 0, 0.1);
        padding: 16px 8px;
        font-family: monospace;
        font-size: 12px;
        color: rgba(0, 0, 0, 0.5);
        user-select: none;
        border-right: 1px solid rgba(0, 0, 0, 0.1);
    `;
    
    lines.forEach((_, index) => {
        if (index < lines.length - 1) { // Don't count empty last line
            const lineDiv = document.createElement('div');
            lineDiv.textContent = (index + 1).toString();
            lineDiv.style.lineHeight = '1.5';
            lineNumbersDiv.appendChild(lineDiv);
        }
    });
    
    const pre = codeBlock.parentElement;
    pre.style.position = 'relative';
    pre.style.paddingLeft = '50px';
    pre.insertBefore(lineNumbersDiv, codeBlock);
}

// Get programming language from code block
function getCodeLanguage(codeBlock) {
    const classes = codeBlock.className.split(' ');
    for (let cls of classes) {
        if (cls.startsWith('language-')) {
            return cls.replace('language-', '');
        }
    }
    return null;
}

// Add language label to code block
function addLanguageLabel(codeBlock, language) {
    const label = document.createElement('div');
    label.className = 'language-label';
    label.textContent = language.toUpperCase();
    label.style.cssText = `
        position: absolute;
        top: 8px;
        left: 8px;
        background: var(--sr-primary, #2196f3);
        color: white;
        padding: 2px 8px;
        border-radius: 4px;
        font-size: 10px;
        font-weight: bold;
        text-transform: uppercase;
        z-index: 10;
    `;
    
    const pre = codeBlock.parentElement;
    pre.style.position = 'relative';
    pre.appendChild(label);
}

// Initialize tooltips
function initializeTooltips() {
    // Create tooltips for API functions
    const apiElements = document.querySelectorAll('code');
    apiElements.forEach(element => {
        if (element.textContent.includes('state.') || element.textContent.includes('metrics.')) {
            element.classList.add('tooltip');
            const tooltipText = document.createElement('span');
            tooltipText.className = 'tooltip-text';
            tooltipText.textContent = getAPITooltip(element.textContent);
            element.appendChild(tooltipText);
        }
    });
}

// Get tooltip text for API functions
function getAPITooltip(functionName) {
    const tooltips = {
        'state.set': 'Set a key-value pair with optional TTL',
        'state.get': 'Get value by key with optional default',
        'state.delete': 'Delete a key from storage',
        'state.increment': 'Atomically increment a numeric value',
        'state.lock': 'Acquire a distributed lock',
        'metrics.gauge': 'Set a gauge metric (current value)',
        'metrics.counter': 'Increment a counter metric',
        'metrics.timer': 'Time function execution',
        'metrics.system_cpu': 'Get current CPU usage',
        'metrics.system_memory': 'Get memory usage information'
    };
    
    for (let [key, description] of Object.entries(tooltips)) {
        if (functionName.includes(key)) {
            return description;
        }
    }
    
    return 'Sloth Runner API function';
}

// Add language indicators to navigation
function addLanguageIndicators() {
    const navItems = document.querySelectorAll('.md-nav__link');
    navItems.forEach(link => {
        const text = link.textContent;
        if (text.includes('🇺🇸') || text.includes('English')) {
            addLanguageBadge(link, 'en', 'English');
        } else if (text.includes('🇧🇷') || text.includes('Português')) {
            addLanguageBadge(link, 'pt', 'Português');
        } else if (text.includes('🇨🇳') || text.includes('中文')) {
            addLanguageBadge(link, 'zh', '中文');
        }
    });
}

// Add language badge to navigation item
function addLanguageBadge(element, lang, name) {
    const badge = document.createElement('span');
    badge.className = `lang-badge ${lang}`;
    badge.textContent = name;
    badge.style.cssText = `
        margin-left: 8px;
        font-size: 0.7rem;
        padding: 2px 6px;
        border-radius: 8px;
        background: var(--sr-primary, #2196f3);
        color: white;
    `;
    element.appendChild(badge);
}

// Initialize progress bars for feature completion
function initializeProgressBars() {
    const progressBars = document.querySelectorAll('.progress-bar');
    progressBars.forEach(bar => {
        const progress = bar.querySelector('.progress');
        const targetWidth = progress.dataset.width || '0%';
        
        // Animate progress bar
        setTimeout(() => {
            progress.style.width = targetWidth;
        }, 500);
    });
}

// Enable smooth scrolling for anchor links
function enableSmoothScrolling() {
    const anchorLinks = document.querySelectorAll('a[href^="#"]');
    anchorLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            const targetId = this.getAttribute('href').substring(1);
            const targetElement = document.getElementById(targetId);
            
            if (targetElement) {
                e.preventDefault();
                targetElement.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
                
                // Add highlight effect to target element
                targetElement.style.backgroundColor = 'rgba(33, 150, 243, 0.1)';
                targetElement.style.transition = 'background-color 0.3s ease';
                
                setTimeout(() => {
                    targetElement.style.backgroundColor = '';
                }, 2000);
            }
        });
    });
}

// Enhance search functionality
function enhanceSearch() {
    const searchInput = document.querySelector('.md-search__input');
    if (searchInput) {
        // Add search suggestions
        searchInput.addEventListener('input', function() {
            const query = this.value.toLowerCase();
            if (query.length > 2) {
                showSearchSuggestions(query);
            }
        });
    }
}

// Show search suggestions
function showSearchSuggestions(query) {
    const suggestions = [
        'state management',
        'metrics monitoring', 
        'distributed locks',
        'TTL expiration',
        'atomic operations',
        'system metrics',
        'custom metrics',
        'health checks',
        'API reference',
        'lua examples'
    ];
    
    const filteredSuggestions = suggestions.filter(s => 
        s.toLowerCase().includes(query)
    );
    
    // Display suggestions (implementation depends on search framework)
    console.log('Search suggestions:', filteredSuggestions);
}

// Add keyboard shortcuts
document.addEventListener('keydown', function(e) {
    // Ctrl/Cmd + K for search
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        const searchInput = document.querySelector('.md-search__input');
        if (searchInput) {
            searchInput.focus();
        }
    }
    
    // Escape to close search
    if (e.key === 'Escape') {
        const searchInput = document.querySelector('.md-search__input');
        if (searchInput && document.activeElement === searchInput) {
            searchInput.blur();
        }
    }
});

// Add reading time estimator
function addReadingTime() {
    const content = document.querySelector('.md-content');
    if (content) {
        const text = content.textContent;
        const words = text.split(/\s+/).length;
        const readingTime = Math.ceil(words / 200); // Average reading speed
        
        const indicator = document.createElement('div');
        indicator.className = 'reading-time';
        indicator.innerHTML = `📖 ${readingTime} min read`;
        indicator.style.cssText = `
            position: fixed;
            bottom: 20px;
            right: 20px;
            background: var(--sr-primary, #2196f3);
            color: white;
            padding: 8px 12px;
            border-radius: 20px;
            font-size: 0.8rem;
            z-index: 1000;
            opacity: 0.8;
        `;
        
        document.body.appendChild(indicator);
    }
}

// Initialize reading time on page load
setTimeout(addReadingTime, 1000);

// Add theme toggle enhancement
function enhanceThemeToggle() {
    const themeToggle = document.querySelector('[data-md-component="palette"]');
    if (themeToggle) {
        themeToggle.addEventListener('change', function() {
            // Add smooth transition when theme changes
            document.body.style.transition = 'background-color 0.3s ease, color 0.3s ease';
            
            setTimeout(() => {
                document.body.style.transition = '';
            }, 300);
        });
    }
}

// Initialize theme toggle enhancement
setTimeout(enhanceThemeToggle, 500);

// Add scroll progress indicator
function addScrollProgress() {
    const progressBar = document.createElement('div');
    progressBar.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 0%;
        height: 3px;
        background: linear-gradient(45deg, var(--sr-primary, #2196f3), var(--sr-secondary, #ff9800));
        z-index: 9999;
        transition: width 0.1s ease;
    `;
    
    document.body.appendChild(progressBar);
    
    window.addEventListener('scroll', function() {
        const scrolled = (window.scrollY / (document.body.scrollHeight - window.innerHeight)) * 100;
        progressBar.style.width = `${Math.min(scrolled, 100)}%`;
    });
}

// Initialize scroll progress
addScrollProgress();