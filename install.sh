#!/usr/bin/env bash
#
# 🦥 Sloth Runner - Universal Installer
# Install the latest version of Sloth Runner
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash
#   Or: bash install.sh [OPTIONS]
#
# Options:
#   --version VERSION    Install specific version (e.g., v3.23.1)
#   --install-dir DIR    Installation directory (default: $HOME/.local/bin or /usr/local/bin)
#   --no-sudo           Install to $HOME/.local/bin without sudo
#   --help              Show this help message
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# GitHub repository
OWNER="chalkan3-sloth"
REPO="sloth-runner"
BINARY_NAME="sloth-runner"

# Default values
VERSION=""
INSTALL_DIR=""
USE_SUDO="auto"
FORCE_INSTALL=false

# Helper functions
print_banner() {
    echo -e "${CYAN}"
    cat << "EOF"
    _____ __      __  __       ____                             
   / ___// /___  / /_/ /_     / __ \__  ______  ____  ___  _____
   \__ \/ / __ \/ __/ __ \   / /_/ / / / / __ \/ __ \/ _ \/ ___/
  ___/ / / /_/ / /_/ / / /  / _, _/ /_/ / / / / / / /  __/ /    
 /____/_/\____/\__/_/ /_/  /_/ |_|\__,_/_/ /_/_/ /_/\___/_/     
                                                                 
           Enterprise-Grade Task Automation Platform            
EOF
    echo -e "${NC}"
}

info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

success() {
    echo -e "${GREEN}✅${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠️${NC}  $1"
}

error() {
    echo -e "${RED}❌${NC} $1"
}

cleanup() {
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

show_help() {
    print_banner
    cat << EOF
${GREEN}Installation Options:${NC}

  --version VERSION       Install specific version (e.g., v3.23.1)
  --install-dir DIR       Installation directory
                         (default: \$HOME/.local/bin or /usr/local/bin)
  --no-sudo              Install to \$HOME/.local/bin without sudo
  --force                Overwrite existing installation
  --help                 Show this help message

${GREEN}Examples:${NC}

  # Install latest version to default location
  bash install.sh

  # Install specific version
  bash install.sh --version v3.23.1

  # Install to user directory (no sudo)
  bash install.sh --no-sudo

  # Install to custom directory
  bash install.sh --install-dir /opt/bin

${GREEN}Default Installation Locations:${NC}

  With sudo:    /usr/local/bin/$BINARY_NAME
  Without sudo: \$HOME/.local/bin/$BINARY_NAME

${GREEN}One-Line Install:${NC}

  curl -fsSL https://raw.githubusercontent.com/$OWNER/$REPO/master/install.sh | bash

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --no-sudo)
            USE_SUDO="no"
            shift
            ;;
        --force)
            FORCE_INSTALL=true
            shift
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        *)
            error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Detect OS and Architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64)
            ARCH="arm64"
            ;;
        arm64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            info "Supported: x86_64 (amd64), arm64, aarch64"
            exit 1
            ;;
    esac

    case "$OS" in
        linux|darwin)
            ;;
        *)
            error "Unsupported operating system: $OS"
            info "Supported: Linux, macOS (Darwin)"
            exit 1
            ;;
    esac
}

# Get latest release tag
get_latest_release() {
    # Try with gh CLI first
    if command -v gh &> /dev/null; then
        LATEST=$(gh release list --repo "$OWNER/$REPO" --limit 1 2>/dev/null | awk '{print $3}' | head -1)
        if [ -n "$LATEST" ]; then
            echo "$LATEST"
            return
        fi
    fi
    
    # Fallback to curl + jq
    if command -v jq &> /dev/null; then
        LATEST=$(curl -fsSL "https://api.github.com/repos/$OWNER/$REPO/releases/latest" | jq -r '.tag_name')
        if [ -n "$LATEST" ] && [ "$LATEST" != "null" ]; then
            echo "$LATEST"
            return
        fi
    fi
    
    # Fallback to curl + grep
    LATEST=$(curl -fsSL "https://api.github.com/repos/$OWNER/$REPO/releases/latest" | \
             grep '"tag_name":' | \
             sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST" ]; then
        error "Could not fetch latest release tag" >&2
        info "Please specify version with --version flag" >&2
        exit 1
    fi
    
    echo "$LATEST"
}

# Check if binary already exists
check_existing() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        EXISTING_PATH=$(command -v "$BINARY_NAME")
        EXISTING_VERSION=$("$BINARY_NAME" --version 2>/dev/null | head -1 || echo "unknown")
        
        warn "$BINARY_NAME is already installed at: $EXISTING_PATH"
        info "Installed version: $EXISTING_VERSION"
        
        if [ "$FORCE_INSTALL" = false ]; then
            read -p "Overwrite? [y/N] " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Installation cancelled"
                exit 0
            fi
        fi
    fi
}

# Determine installation directory
determine_install_dir() {
    if [ -n "$INSTALL_DIR" ]; then
        # User specified directory
        return
    fi
    
    if [ "$USE_SUDO" = "no" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        USE_SUDO="no"
    elif [ "$USE_SUDO" = "auto" ]; then
        # Try to detect if we can write to /usr/local/bin
        if [ -w "/usr/local/bin" ] || [ "$(id -u)" -eq 0 ]; then
            INSTALL_DIR="/usr/local/bin"
            USE_SUDO="no"
        elif command -v sudo &> /dev/null; then
            INSTALL_DIR="/usr/local/bin"
            USE_SUDO="yes"
        else
            INSTALL_DIR="$HOME/.local/bin"
            USE_SUDO="no"
        fi
    fi
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        info "Creating directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR" 2>/dev/null || {
            if [ "$USE_SUDO" = "yes" ]; then
                sudo mkdir -p "$INSTALL_DIR"
            else
                error "Failed to create directory: $INSTALL_DIR"
                exit 1
            fi
        }
    fi
}

# Download and install binary
install_binary() {
    local version=$1
    local artifact_name="${BINARY_NAME}_${version}_${OS}_${ARCH}.tar.gz"
    local download_url="https://github.com/$OWNER/$REPO/releases/download/${version}/${artifact_name}"
    
    TEMP_DIR=$(mktemp -d)
    local temp_archive="$TEMP_DIR/$artifact_name"
    
    info "Downloading $BINARY_NAME $version..."
    info "Platform: $OS/$ARCH"
    info "URL: $download_url"
    
    if ! curl -fsSL "$download_url" -o "$temp_archive"; then
        error "Failed to download from $download_url"
        info "Available releases:"
        curl -fsSL "https://api.github.com/repos/$OWNER/$REPO/releases/latest" | \
            grep "browser_download_url.*tar.gz" | cut -d'"' -f4 | sed 's/.*\//  - /'
        exit 1
    fi
    
    success "Downloaded successfully"
    
    info "Extracting archive..."
    if ! tar -xzf "$temp_archive" -C "$TEMP_DIR"; then
        error "Failed to extract archive"
        exit 1
    fi
    
    # Find the binary
    local binary_path
    binary_path=$(find "$TEMP_DIR" -name "$BINARY_NAME" -type f | head -1)
    
    if [ -z "$binary_path" ]; then
        error "Binary not found in archive"
        exit 1
    fi
    
    info "Installing to $INSTALL_DIR/$BINARY_NAME..."
    
    # Install binary
    if [ "$USE_SUDO" = "yes" ]; then
        sudo install -m 755 "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    else
        install -m 755 "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    success "Binary installed successfully!"
}

# Verify installation
verify_installation() {
    info "Verifying installation..."
    
    # Check if directory is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warn "Installation directory not in PATH: $INSTALL_DIR"
        echo ""
        info "Add to your PATH by adding this line to ~/.bashrc or ~/.zshrc:"
        echo -e "${YELLOW}  export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
        echo ""
    fi
    
    # Try to run the binary
    if command -v "$BINARY_NAME" &> /dev/null; then
        local version
        version=$("$BINARY_NAME" --version 2>/dev/null | head -1 || echo "installed")
        success "Installation verified: $version"
        success "Binary location: $(command -v $BINARY_NAME)"
    else
        warn "Binary installed but not found in PATH"
        info "Installed at: $INSTALL_DIR/$BINARY_NAME"
        info "Add $INSTALL_DIR to your PATH to use it"
    fi
}

# Show post-install instructions
show_post_install() {
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}    🎉 Sloth Runner installed successfully!${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}📚 Quick Start:${NC}"
    echo -e "  ${BINARY_NAME} --help          # Show help"
    echo -e "  ${BINARY_NAME} --version       # Show version"
    echo -e "  ${BINARY_NAME} run workflow.sloth  # Run a workflow"
    echo ""
    echo -e "${CYAN}📖 Documentation:${NC}"
    echo -e "  https://chalkan3.github.io/sloth-runner/"
    echo ""
    echo -e "${CYAN}🐛 Issues & Support:${NC}"
    echo -e "  https://github.com/$OWNER/$REPO/issues"
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
}

# Main installation flow
main() {
    print_banner
    
    info "Starting Sloth Runner installation..."
    
    # Detect platform
    detect_platform
    success "Platform detected: $OS/$ARCH"
    
    # Get version to install
    if [ -z "$VERSION" ]; then
        info "Fetching latest release..."
        VERSION=$(get_latest_release)
    fi
    success "Version to install: $VERSION"
    
    # Check for existing installation
    check_existing
    
    # Determine installation directory
    determine_install_dir
    info "Installation directory: $INSTALL_DIR"
    
    # Download and install
    install_binary "$VERSION"
    
    # Verify
    verify_installation
    
    # Show post-install info
    show_post_install
}

# Run main function
main "$@"
