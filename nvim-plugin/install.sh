#!/bin/bash

# Sloth Runner DSL Neovim Plugin Installer
# This script installs the Sloth Runner DSL syntax highlighting plugin for Neovim

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default Neovim config directory
NVIM_CONFIG_DIR="$HOME/.config/nvim"

# Plugin source directory (relative to this script)
PLUGIN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}🦥 Sloth Runner DSL - Neovim Plugin Installer${NC}"
echo "================================================="

# Check if Neovim is installed
if ! command -v nvim &> /dev/null; then
    echo -e "${RED}❌ Neovim is not installed. Please install Neovim first.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Neovim found: $(nvim --version | head -n1)${NC}"

# Check Neovim config directory
if [ ! -d "$NVIM_CONFIG_DIR" ]; then
    echo -e "${YELLOW}📁 Creating Neovim config directory: $NVIM_CONFIG_DIR${NC}"
    mkdir -p "$NVIM_CONFIG_DIR"
fi

# Function to install files
install_plugin_files() {
    echo -e "${BLUE}📦 Installing plugin files...${NC}"
    
    # Create required directories
    mkdir -p "$NVIM_CONFIG_DIR/syntax"
    mkdir -p "$NVIM_CONFIG_DIR/ftdetect" 
    mkdir -p "$NVIM_CONFIG_DIR/ftplugin"
    mkdir -p "$NVIM_CONFIG_DIR/lua"
    mkdir -p "$NVIM_CONFIG_DIR/queries/sloth"
    
    # Copy syntax files
    if [ -f "$PLUGIN_DIR/syntax/sloth.vim" ]; then
        cp "$PLUGIN_DIR/syntax/sloth.vim" "$NVIM_CONFIG_DIR/syntax/"
        echo "  ✅ Copied syntax/sloth.vim"
    fi
    
    # Copy filetype detection
    if [ -f "$PLUGIN_DIR/ftdetect/sloth.vim" ]; then
        cp "$PLUGIN_DIR/ftdetect/sloth.vim" "$NVIM_CONFIG_DIR/ftdetect/"
        echo "  ✅ Copied ftdetect/sloth.vim"
    fi
    
    # Copy filetype plugin
    if [ -f "$PLUGIN_DIR/ftplugin/sloth.vim" ]; then
        cp "$PLUGIN_DIR/ftplugin/sloth.vim" "$NVIM_CONFIG_DIR/ftplugin/"
        echo "  ✅ Copied ftplugin/sloth.vim"
    fi
    
    # Copy Lua module
    if [ -f "$PLUGIN_DIR/lua/sloth-runner.lua" ]; then
        cp "$PLUGIN_DIR/lua/sloth-runner.lua" "$NVIM_CONFIG_DIR/lua/"
        echo "  ✅ Copied lua/sloth-runner.lua"
    fi
    
    # Copy Tree-sitter queries if they exist
    if [ -f "$PLUGIN_DIR/queries/sloth/highlights.scm" ]; then
        cp "$PLUGIN_DIR/queries/sloth/highlights.scm" "$NVIM_CONFIG_DIR/queries/sloth/"
        echo "  ✅ Copied queries/sloth/highlights.scm"
    fi
    
    if [ -f "$PLUGIN_DIR/queries/sloth/indents.scm" ]; then
        cp "$PLUGIN_DIR/queries/sloth/indents.scm" "$NVIM_CONFIG_DIR/queries/sloth/"
        echo "  ✅ Copied queries/sloth/indents.scm"
    fi
}

# Function to add setup to init.lua
setup_init_lua() {
    local init_lua="$NVIM_CONFIG_DIR/init.lua"
    
    echo -e "${BLUE}⚙️  Setting up init.lua configuration...${NC}"
    
    # Check if init.lua exists
    if [ ! -f "$init_lua" ]; then
        echo -e "${YELLOW}📝 Creating init.lua file${NC}"
        touch "$init_lua"
    fi
    
    # Check if sloth-runner setup already exists
    if grep -q "require.*sloth-runner" "$init_lua"; then
        echo -e "${YELLOW}⚠️  Sloth Runner setup already exists in init.lua${NC}"
        return 0
    fi
    
    # Add setup configuration
    cat >> "$init_lua" << 'EOF'

-- Sloth Runner DSL Plugin Configuration
require("sloth-runner").setup({
  runner = {
    command = "sloth-runner", -- Path to your sloth-runner binary
    keymaps = {
      run_file = "<leader>sr",
      list_tasks = "<leader>sl",
      dry_run = "<leader>st", 
      debug = "<leader>sd",
    }
  },
  completion = {
    enable = true,
    snippets = true,
  },
  folding = {
    enable = true,
    auto_close = false,
  }
})
EOF
    
    echo "  ✅ Added Sloth Runner configuration to init.lua"
}

# Function to create example file
create_example() {
    local example_file="$HOME/example.sloth.lua"
    
    echo -e "${BLUE}📄 Creating example file...${NC}"
    
    if [ -f "$PLUGIN_DIR/example.sloth.lua" ]; then
        cp "$PLUGIN_DIR/example.sloth.lua" "$example_file"
        echo "  ✅ Created example file: $example_file"
    else
        # Create a simple example if the file doesn't exist
        cat > "$example_file" << 'EOF'
-- Example Sloth Runner DSL file
-- This file demonstrates syntax highlighting

local build = task("build_app")
    :description("Build the application")
    :command(function(params, deps)
        local result = exec.run("go build -o app ./cmd/main.go")
        return result.success, result.stdout
    end)
    :timeout("5m")
    :build()

workflow.define("simple_build", {
    description = "Simple build workflow",
    version = "1.0.0",
    tasks = { build }
})
EOF
        echo "  ✅ Created simple example file: $example_file"
    fi
}

# Function to check sloth-runner binary
check_sloth_runner() {
    echo -e "${BLUE}🔍 Checking for sloth-runner binary...${NC}"
    
    if command -v sloth-runner &> /dev/null; then
        echo -e "${GREEN}✅ sloth-runner found: $(which sloth-runner)${NC}"
        echo "  Version: $(sloth-runner --version 2>/dev/null || echo 'unknown')"
    else
        echo -e "${YELLOW}⚠️  sloth-runner binary not found in PATH${NC}"
        echo "   You may need to install sloth-runner or update the binary path in init.lua"
    fi
}

# Main installation process
main() {
    echo -e "${BLUE}🚀 Starting installation...${NC}"
    echo ""
    
    # Install plugin files
    install_plugin_files
    echo ""
    
    # Setup init.lua
    setup_init_lua
    echo ""
    
    # Create example file
    create_example
    echo ""
    
    # Check for sloth-runner binary
    check_sloth_runner
    echo ""
    
    echo -e "${GREEN}🎉 Installation completed successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Restart Neovim or run :source ~/.config/nvim/init.lua"
    echo "2. Open the example file: nvim ~/example.sloth.lua"
    echo "3. Test syntax highlighting and features"
    echo ""
    echo "Key mappings:"
    echo "  <leader>sr - Run current file"
    echo "  <leader>sl - List tasks"
    echo "  <leader>st - Dry run"
    echo "  <leader>sd - Debug"
    echo ""
    echo "Snippets:"
    echo "  _task     - Insert task template"
    echo "  _workflow - Insert workflow template"
    echo ""
    echo "Text objects:"
    echo "  vit - Select task block"
    echo "  viw - Select workflow block"
    echo ""
    echo -e "${BLUE}Happy coding with Sloth Runner DSL! 🦥⚡${NC}"
}

# Run with --help flag
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Sloth Runner DSL Neovim Plugin Installer"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --help, -h     Show this help message"
    echo "  --uninstall    Remove the plugin"
    echo ""
    echo "This script installs syntax highlighting and IDE features"
    echo "for Sloth Runner DSL files in Neovim."
    exit 0
fi

# Uninstall option
if [ "$1" = "--uninstall" ]; then
    echo -e "${YELLOW}🗑️  Uninstalling Sloth Runner DSL plugin...${NC}"
    
    # Remove files
    rm -f "$NVIM_CONFIG_DIR/syntax/sloth.vim"
    rm -f "$NVIM_CONFIG_DIR/ftdetect/sloth.vim" 
    rm -f "$NVIM_CONFIG_DIR/ftplugin/sloth.vim"
    rm -f "$NVIM_CONFIG_DIR/lua/sloth-runner.lua"
    rm -rf "$NVIM_CONFIG_DIR/queries/sloth"
    
    echo -e "${GREEN}✅ Plugin files removed${NC}"
    echo -e "${YELLOW}⚠️  Please manually remove Sloth Runner configuration from init.lua${NC}"
    exit 0
fi

# Run main installation
main