#!/usr/bin/env bash
#
# 🦥 Sloth Runner - Bootstrap Agent Script
# Installs and configures sloth-runner agent with systemd
#
# Usage:
#   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) [OPTIONS]
#
# Options:
#   --name NAME              Agent name (required)
#   --master HOST:PORT       Master server address (default: localhost:50053)
#   --port PORT              Agent port (default: 50051)
#   --bind-address IP        IP address for the agent to bind to
#   --user USER              User to run the agent as (default: current user)
#   --install-dir DIR        Installation directory (default: /usr/local/bin)
#   --no-systemd            Skip systemd service creation
#   --no-sudo               Install without sudo to ~/.local/bin
#   --version VERSION        Install specific version
#   --help                  Show this help message
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

# Default values
AGENT_NAME=""
MASTER_ADDRESS="localhost:50053"
AGENT_PORT="50051"
BIND_ADDRESS=""
SERVICE_USER="${USER}"
INSTALL_DIR=""
SKIP_SYSTEMD=false
USE_SUDO="auto"
VERSION=""
INSTALL_SCRIPT_URL="https://raw.githubusercontent.com/${OWNER}/${REPO}/master/install.sh"

# Helper functions
print_banner() {
    echo -e "${CYAN}"
    cat << "EOF"
    _____ __      __  __       ____                             
   / ___// /___  / /_/ /_     / __ \__  ______  ____  ___  _____
   \__ \/ / __ \/ __/ __ \   / /_/ / / / / __ \/ __ \/ _ \/ ___/
  ___/ / / /_/ / /_/ / / /  / _, _/ /_/ / / / / / / /  __/ /    
 /____/_/\____/\__/_/ /_/  /_/ |_|\__,_/_/ /_/_/ /_/\___/_/     
                                                                 
           Agent Bootstrap - Enterprise-Grade Automation            
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
    exit 1
}

show_help() {
    print_banner
    cat << EOF
${GREEN}Bootstrap Options:${NC}

  ${YELLOW}Required:${NC}
  --name NAME              Agent name (required)

  ${YELLOW}Optional:${NC}
  --master HOST:PORT       Master server address (default: localhost:50053)
  --port PORT              Agent port (default: 50051)
  --bind-address IP        IP address for agent to bind to (auto-detected if not set)
  --user USER              User to run the agent as (default: $USER)
  --install-dir DIR        Installation directory (default: /usr/local/bin)
  --version VERSION        Install specific version (default: latest)
  --no-systemd            Skip systemd service creation
  --no-sudo               Install without sudo to ~/.local/bin
  --help                  Show this help message

${GREEN}Examples:${NC}

  # Basic installation with agent name
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) --name myagent

  # Full configuration
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) \\
    --name production-agent-1 \\
    --master 192.168.1.10:50053 \\
    --port 50051 \\
    --bind-address 192.168.1.20

  # User installation (no systemd, no sudo)
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) \\
    --name myagent \\
    --no-sudo \\
    --no-systemd

${GREEN}Post-Installation:${NC}

  After installation, the agent will be:
  - Installed to $INSTALL_DIR/sloth-runner
  - Configured as systemd service (if --no-systemd not used)
  - Started and enabled on boot
  - Registered with master at $MASTER_ADDRESS

${GREEN}Check Agent Status:${NC}

  systemctl status sloth-runner-agent
  journalctl -u sloth-runner-agent -f

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --name)
            AGENT_NAME="$2"
            shift 2
            ;;
        --master)
            MASTER_ADDRESS="$2"
            shift 2
            ;;
        --port)
            AGENT_PORT="$2"
            shift 2
            ;;
        --bind-address)
            BIND_ADDRESS="$2"
            shift 2
            ;;
        --user)
            SERVICE_USER="$2"
            shift 2
            ;;
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --no-systemd)
            SKIP_SYSTEMD=true
            shift
            ;;
        --no-sudo)
            USE_SUDO="no"
            shift
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        *)
            error "Unknown option: $1\nUse --help for usage information"
            ;;
    esac
done

# Validate required parameters
if [ -z "$AGENT_NAME" ]; then
    error "Agent name is required. Use --name to specify it.\nExample: $0 --name myagent\nUse --help for more information."
fi

# Detect OS
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    case "$OS" in
        linux)
            # Check if systemd is available
            if ! command -v systemctl &> /dev/null; then
                if [ "$SKIP_SYSTEMD" = false ]; then
                    warn "systemd not found. Skipping service creation."
                    SKIP_SYSTEMD=true
                fi
            fi
            ;;
        darwin)
            if [ "$SKIP_SYSTEMD" = false ]; then
                warn "macOS detected. systemd not available. Skipping service creation."
                SKIP_SYSTEMD=true
            fi
            ;;
        *)
            error "Unsupported operating system: $OS"
            ;;
    esac
}

# Auto-detect bind address if not specified
detect_bind_address() {
    if [ -z "$BIND_ADDRESS" ]; then
        # Try to get the primary network interface IP
        if command -v ip &> /dev/null; then
            BIND_ADDRESS=$(ip route get 1 2>/dev/null | awk '{print $7}' | head -1)
        elif command -v hostname &> /dev/null; then
            BIND_ADDRESS=$(hostname -I 2>/dev/null | awk '{print $1}')
        fi
        
        # Fallback to 0.0.0.0 if detection fails
        if [ -z "$BIND_ADDRESS" ]; then
            BIND_ADDRESS="0.0.0.0"
            warn "Could not detect IP address, using 0.0.0.0"
        else
            info "Detected bind address: $BIND_ADDRESS"
        fi
    fi
}

# Install sloth-runner
install_sloth_runner() {
    info "Downloading and running installer..."
    
    # Build install command
    INSTALL_CMD="bash"
    INSTALL_ARGS=()
    
    if [ -n "$VERSION" ]; then
        INSTALL_ARGS+=("--version" "$VERSION")
    fi
    
    if [ -n "$INSTALL_DIR" ]; then
        INSTALL_ARGS+=("--install-dir" "$INSTALL_DIR")
    fi
    
    if [ "$USE_SUDO" = "no" ]; then
        INSTALL_ARGS+=("--no-sudo")
        # Set install dir if not specified
        if [ -z "$INSTALL_DIR" ]; then
            INSTALL_DIR="$HOME/.local/bin"
        fi
    else
        # Default install dir with sudo
        if [ -z "$INSTALL_DIR" ]; then
            INSTALL_DIR="/usr/local/bin"
        fi
    fi
    
    INSTALL_ARGS+=("--force")
    
    # Download and run installer
    if ! curl -fsSL "$INSTALL_SCRIPT_URL" | bash -s -- "${INSTALL_ARGS[@]}"; then
        error "Failed to install sloth-runner"
    fi
    
    success "sloth-runner installed successfully"
}

# Create systemd service
create_systemd_service() {
    if [ "$SKIP_SYSTEMD" = true ]; then
        return
    fi
    
    info "Creating systemd service..."
    
    # Build agent start command
    local agent_cmd="$INSTALL_DIR/sloth-runner agent start"
    agent_cmd="$agent_cmd --name $AGENT_NAME"
    agent_cmd="$agent_cmd --master $MASTER_ADDRESS"
    agent_cmd="$agent_cmd --port $AGENT_PORT"
    agent_cmd="$agent_cmd --bind-address $BIND_ADDRESS"
    agent_cmd="$agent_cmd --daemon"
    
    # Create systemd service file
    local service_file="/etc/systemd/system/sloth-runner-agent.service"
    local temp_service="/tmp/sloth-runner-agent.service.$$"
    
    cat > "$temp_service" << EOF
[Unit]
Description=Sloth Runner Agent - $AGENT_NAME
Documentation=https://chalkan3.github.io/sloth-runner/
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$SERVICE_USER
Restart=on-failure
RestartSec=5s
StartLimitInterval=60s
StartLimitBurst=3

# Agent Configuration
ExecStart=$agent_cmd

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=sloth-runner-agent

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=/var/log

# Performance
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
    
    # Install service file
    if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
        sudo mv "$temp_service" "$service_file"
        sudo chmod 644 "$service_file"
    else
        mv "$temp_service" "$service_file"
        chmod 644 "$service_file"
    fi
    
    success "Systemd service created at $service_file"
}

# Enable and start service
start_service() {
    if [ "$SKIP_SYSTEMD" = true ]; then
        return
    fi
    
    info "Reloading systemd daemon..."
    if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
        sudo systemctl daemon-reload
    else
        systemctl daemon-reload
    fi
    
    info "Enabling sloth-runner-agent service..."
    if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
        sudo systemctl enable sloth-runner-agent
    else
        systemctl enable sloth-runner-agent
    fi
    
    info "Starting sloth-runner-agent service..."
    if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
        sudo systemctl start sloth-runner-agent
    else
        systemctl start sloth-runner-agent
    fi
    
    success "Service started and enabled"
}

# Verify agent is running
verify_agent() {
    if [ "$SKIP_SYSTEMD" = true ]; then
        warn "Systemd not configured. Please start agent manually:"
        echo ""
        echo -e "${YELLOW}  $INSTALL_DIR/sloth-runner agent start \\"
        echo -e "    --name $AGENT_NAME \\"
        echo -e "    --master $MASTER_ADDRESS \\"
        echo -e "    --port $AGENT_PORT \\"
        echo -e "    --bind-address $BIND_ADDRESS \\"
        echo -e "    --daemon${NC}"
        echo ""
        return
    fi
    
    info "Verifying agent status..."
    sleep 2
    
    if systemctl is-active --quiet sloth-runner-agent; then
        success "Agent is running!"
        
        # Show status
        echo ""
        if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
            sudo systemctl status sloth-runner-agent --no-pager -l
        else
            systemctl status sloth-runner-agent --no-pager -l
        fi
    else
        error "Agent failed to start. Check logs with: journalctl -u sloth-runner-agent -n 50"
    fi
}

# Show post-install instructions
show_post_install() {
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}    🎉 Sloth Runner Agent Bootstrap Complete!${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}📋 Agent Configuration:${NC}"
    echo -e "  Name:          ${YELLOW}$AGENT_NAME${NC}"
    echo -e "  Master:        ${YELLOW}$MASTER_ADDRESS${NC}"
    echo -e "  Port:          ${YELLOW}$AGENT_PORT${NC}"
    echo -e "  Bind Address:  ${YELLOW}$BIND_ADDRESS${NC}"
    echo -e "  User:          ${YELLOW}$SERVICE_USER${NC}"
    echo ""
    
    if [ "$SKIP_SYSTEMD" = false ]; then
        echo -e "${CYAN}🔧 Service Management:${NC}"
        echo -e "  sudo systemctl status sloth-runner-agent    # Check status"
        echo -e "  sudo systemctl restart sloth-runner-agent   # Restart agent"
        echo -e "  sudo systemctl stop sloth-runner-agent      # Stop agent"
        echo -e "  sudo journalctl -u sloth-runner-agent -f    # View logs"
        echo ""
    fi
    
    echo -e "${CYAN}📊 Check Agent on Master:${NC}"
    echo -e "  sloth-runner agent list                     # List all agents"
    echo -e "  sloth-runner agent run $AGENT_NAME \"hostname\"  # Test agent"
    echo ""
    echo -e "${CYAN}📖 Documentation:${NC}"
    echo -e "  https://chalkan3.github.io/sloth-runner/"
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
}

# Main bootstrap flow
main() {
    print_banner
    
    info "Starting Sloth Runner Agent bootstrap..."
    echo ""
    
    # Detect OS
    detect_os
    
    # Detect bind address
    detect_bind_address
    
    # Install sloth-runner
    install_sloth_runner
    echo ""
    
    # Create systemd service
    create_systemd_service
    echo ""
    
    # Start service
    start_service
    echo ""
    
    # Verify
    verify_agent
    
    # Show post-install info
    show_post_install
}

# Run main function
main "$@"
