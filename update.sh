#!/usr/bin/env bash
#
# 🦥 Sloth Runner - Universal Updater
# Update Sloth Runner to the latest version
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/update.sh | bash
#   Or: bash update.sh [OPTIONS]
#
# Options:
#   --version VERSION    Update to specific version (e.g., v3.23.1)
#   --check-only        Only check for updates, don't install
#   --force             Force update even if already up to date
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
CHECK_ONLY=false
FORCE_UPDATE=false
INSTALL_DIR=""
USE_SUDO="auto"

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
${GREEN}Update Options:${NC}

  --version VERSION       Update to specific version (e.g., v3.23.1)
  --check-only           Only check for updates, don't install
  --force                Force update even if already up to date
  --help                 Show this help message

${GREEN}Examples:${NC}

  # Update to latest version
  bash update.sh

  # Check for updates without installing
  bash update.sh --check-only

  # Update to specific version
  bash update.sh --version v3.23.1

  # Force update
  bash update.sh --force

${GREEN}One-Line Update:${NC}

  curl -fsSL https://raw.githubusercontent.com/$OWNER/$REPO/master/update.sh | bash

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --check-only)
            CHECK_ONLY=true
            shift
            ;;
        --force)
            FORCE_UPDATE=true
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
        exit 1
    fi
    
    echo "$LATEST"
}

# Check current installation
check_current_version() {
    if ! command -v "$BINARY_NAME" &> /dev/null; then
        error "$BINARY_NAME is not installed"
        info "Use install.sh to install Sloth Runner first"
        info "curl -fsSL https://raw.githubusercontent.com/$OWNER/$REPO/master/install.sh | bash"
        exit 1
    fi
    
    CURRENT_PATH=$(command -v "$BINARY_NAME")
    INSTALL_DIR=$(dirname "$CURRENT_PATH")
    
    # Get current version
    CURRENT_VERSION=$("$BINARY_NAME" --version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown")
    
    if [ "$CURRENT_VERSION" = "unknown" ]; then
        warn "Could not determine current version"
    else
        info "Current version: $CURRENT_VERSION"
        info "Installed at: $CURRENT_PATH"
    fi
}

# Compare versions
version_gt() {
    # Remove 'v' prefix for comparison
    local v1="${1#v}"
    local v2="${2#v}"
    
    # Split versions into arrays
    IFS='.' read -ra V1 <<< "$v1"
    IFS='.' read -ra V2 <<< "$v2"
    
    # Compare major, minor, patch
    for i in 0 1 2; do
        local num1=${V1[$i]:-0}
        local num2=${V2[$i]:-0}
        
        if [ "$num1" -gt "$num2" ]; then
            return 0
        elif [ "$num1" -lt "$num2" ]; then
            return 1
        fi
    done
    
    return 1
}

# Check if update is needed
check_update_needed() {
    local latest=$1
    
    if [ "$CURRENT_VERSION" = "unknown" ]; then
        warn "Cannot compare versions, will proceed with update"
        return 0
    fi
    
    if [ "$CURRENT_VERSION" = "$latest" ]; then
        success "Already up to date! ($CURRENT_VERSION)"
        
        if [ "$FORCE_UPDATE" = false ]; then
            return 1
        else
            warn "Force update enabled, reinstalling..."
            return 0
        fi
    fi
    
    if version_gt "$latest" "$CURRENT_VERSION"; then
        info "Update available: $CURRENT_VERSION → $latest"
        return 0
    else
        success "Current version ($CURRENT_VERSION) is newer than latest ($latest)"
        
        if [ "$FORCE_UPDATE" = false ]; then
            return 1
        else
            warn "Force update enabled, downgrading to $latest..."
            return 0
        fi
    fi
}

# Determine if sudo is needed
determine_sudo() {
    if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
        USE_SUDO="no"
    elif command -v sudo &> /dev/null; then
        USE_SUDO="yes"
    else
        error "Cannot write to $INSTALL_DIR and sudo is not available"
        exit 1
    fi
}

# Download and install binary
update_binary() {
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
    
    info "Updating $INSTALL_DIR/$BINARY_NAME..."
    
    # Backup current binary
    local backup_path="$INSTALL_DIR/${BINARY_NAME}.backup"
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        info "Creating backup: $backup_path"
        if [ "$USE_SUDO" = "yes" ]; then
            sudo cp "$INSTALL_DIR/$BINARY_NAME" "$backup_path"
        else
            cp "$INSTALL_DIR/$BINARY_NAME" "$backup_path"
        fi
    fi
    
    # Install new binary
    if [ "$USE_SUDO" = "yes" ]; then
        sudo install -m 755 "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    else
        install -m 755 "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    success "Binary updated successfully!"
    
    # Remove backup
    if [ -f "$backup_path" ]; then
        if [ "$USE_SUDO" = "yes" ]; then
            sudo rm -f "$backup_path"
        else
            rm -f "$backup_path"
        fi
    fi
}

# Verify update
verify_update() {
    local expected_version=$1
    
    info "Verifying update..."
    
    if command -v "$BINARY_NAME" &> /dev/null; then
        local new_version
        new_version=$("$BINARY_NAME" --version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown")
        
        if [ "$new_version" = "$expected_version" ]; then
            success "Update verified: $new_version"
            success "Binary location: $(command -v $BINARY_NAME)"
        else
            warn "Version mismatch: expected $expected_version, got $new_version"
        fi
    else
        error "Binary not found after update"
        exit 1
    fi
}

# Show post-update info
show_post_update() {
    local old_version=$1
    local new_version=$2
    
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}    🎉 Sloth Runner updated successfully!${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}📦 Update Details:${NC}"
    echo -e "  Old Version: ${old_version}"
    echo -e "  New Version: ${new_version}"
    echo -e "  Location:    $(command -v $BINARY_NAME 2>/dev/null || echo $INSTALL_DIR/$BINARY_NAME)"
    echo ""
    echo -e "${CYAN}📖 Changelog:${NC}"
    echo -e "  https://github.com/$OWNER/$REPO/releases/tag/${new_version}"
    echo ""
    echo -e "${CYAN}📚 Documentation:${NC}"
    echo -e "  https://chalkan3-sloth.github.io/sloth-runner/"
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
}

# Main update flow
main() {
    print_banner
    
    info "Starting Sloth Runner update..."
    
    # Detect platform
    detect_platform
    success "Platform detected: $OS/$ARCH"
    
    # Check current installation
    check_current_version
    
    # Get version to update to
    if [ -z "$VERSION" ]; then
        info "Fetching latest release..."
        VERSION=$(get_latest_release)
    fi
    success "Target version: $VERSION"
    
    # Check if update is needed
    if ! check_update_needed "$VERSION"; then
        if [ "$CHECK_ONLY" = true ]; then
            exit 0
        fi
        
        info "No update performed"
        exit 0
    fi
    
    # If check-only mode, exit here
    if [ "$CHECK_ONLY" = true ]; then
        info "Update available but not installing (--check-only mode)"
        exit 0
    fi
    
    # Determine sudo requirements
    determine_sudo
    info "Installation directory: $INSTALL_DIR"
    
    # Save old version for display
    OLD_VERSION="$CURRENT_VERSION"
    
    # Download and update
    update_binary "$VERSION"
    
    # Verify
    verify_update "$VERSION"
    
    # Show post-update info
    show_post_update "$OLD_VERSION" "$VERSION"
}

# Run main function
main "$@"
