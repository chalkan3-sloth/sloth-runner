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
#   --report-address IP      IP address to report to master (useful for NAT/port forwarding)
#   --incus HOST_IP:PORT     Auto-configure for Incus container (sets bind-address to 0.0.0.0 and report-address)
#   --user USER              User to run the agent as (default: current user)
#   --install-dir DIR        Installation directory (default: /usr/local/bin)
#   --no-systemd            Skip systemd service creation
#   --no-sudo               Install without sudo to ~/.local/bin
#   --version VERSION        Install specific version
#   --update                Update existing installation and restart service
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
REPORT_ADDRESS=""
SERVICE_USER="${USER}"
INSTALL_DIR=""
SKIP_SYSTEMD=false
USE_SUDO="auto"
VERSION=""
UPDATE_MODE=false
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
  --report-address IP:PORT IP address to report to master (useful for NAT/port forwarding)
  --incus HOST_IP:PORT     Auto-configure for Incus container (sets bind to 0.0.0.0, report to HOST_IP:PORT)
  --user USER              User to run the agent as (default: $USER)
  --install-dir DIR        Installation directory (default: /usr/local/bin)
  --version VERSION        Install specific version (default: latest)
  --update                 Update existing installation and restart service
  --no-systemd            Skip systemd service creation
  --no-sudo               Install without sudo to ~/.local/bin
  --help                  Show this help message

${GREEN}Examples:${NC}

  # Basic installation with agent name
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) --name myagent

  # Incus container setup (auto-configures bind and report addresses)
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) \\
    --name main \\
    --master 192.168.1.29:50053 \\
    --incus 192.168.1.17:50052

  # Full configuration with port forwarding/NAT
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) \\
    --name production-agent-1 \\
    --master 192.168.1.10:50053 \\
    --port 50051 \\
    --bind-address 0.0.0.0 \\
    --report-address 192.168.1.20:50052

  # User installation (no systemd, no sudo)
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) \\
    --name myagent \\
    --no-sudo \\
    --no-systemd

  # Update existing installation
  bash <(curl -fsSL $INSTALL_SCRIPT_URL) \\
    --name myagent \\
    --update

${GREEN}Incus/LXC Container Setup:${NC}

  When deploying in Incus/LXC containers, use --incus flag:

  1. On the host, configure port forwarding:
     sudo incus config device add <container> sloth-proxy proxy \\
       listen=tcp:0.0.0.0:50052 connect=tcp:127.0.0.1:50051

  2. Inside the container, run bootstrap:
     bash <(curl -fsSL $INSTALL_SCRIPT_URL) \\
       --name main \\
       --master <master_ip>:50053 \\
       --incus <host_ip>:50052

  This automatically sets:
    --bind-address 0.0.0.0  (listen on all interfaces)
    --report-address <host_ip>:50052  (master connects via host)


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
        --report-address)
            REPORT_ADDRESS="$2"
            shift 2
            ;;
        --incus)
            # Auto-configure for Incus container
            BIND_ADDRESS="0.0.0.0"
            REPORT_ADDRESS="$2"
            info "Incus mode: bind-address=0.0.0.0, report-address=$2"
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
        --update)
            UPDATE_MODE=true
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
if [ -z "$AGENT_NAME" ] && [ "$UPDATE_MODE" = false ]; then
    error "Agent name is required. Use --name to specify it.\nExample: $0 --name myagent\nUse --help for more information."
fi

# Detect OS
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    case "$OS" in
        linux)
            # Check if systemd is available and working
            if ! command -v systemctl &> /dev/null; then
                if [ "$SKIP_SYSTEMD" = false ]; then
                    warn "systemd not found. Skipping service creation."
                    SKIP_SYSTEMD=true
                fi
            elif ! systemctl --version &> /dev/null 2>&1; then
                if [ "$SKIP_SYSTEMD" = false ]; then
                    warn "systemd not functioning properly. Skipping service creation."
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
    elif [ "$BIND_ADDRESS" = "0.0.0.0" ]; then
        info "Using bind address: 0.0.0.0 (all interfaces)"
    else
        info "Using bind address: $BIND_ADDRESS"
    fi
    
    # If report address is set, show it
    if [ -n "$REPORT_ADDRESS" ]; then
        info "Using report address: $REPORT_ADDRESS (for master connection)"
    fi
}

# Install sloth-runner
install_sloth_runner() {
    if [ "$UPDATE_MODE" = true ]; then
        info "Update mode: stopping existing service..."
        
        # Stop service if running
        if systemctl is-active --quiet sloth-runner-agent 2>/dev/null; then
            if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
                sudo systemctl stop sloth-runner-agent
            else
                systemctl stop sloth-runner-agent
            fi
            success "Service stopped"
        fi
        
        info "Updating sloth-runner binary..."
    else
        info "Downloading and running installer..."
    fi
    
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
    
    if [ "$UPDATE_MODE" = true ]; then
        success "sloth-runner updated successfully"
    else
        success "sloth-runner installed successfully"
    fi
}

# Create systemd service
create_systemd_service() {
    if [ "$SKIP_SYSTEMD" = true ]; then
        return
    fi
    
    info "Creating systemd service..."
    
    # Determine working directory based on user
    local work_dir="/var/lib/sloth-runner"
    if [ "$SERVICE_USER" != "root" ]; then
        work_dir="/home/$SERVICE_USER/.sloth-runner"
    fi
    
    # Create working directory
    if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
        sudo mkdir -p "$work_dir"
        sudo chown "$SERVICE_USER:$SERVICE_USER" "$work_dir" 2>/dev/null || true
    else
        mkdir -p "$work_dir"
        chown "$SERVICE_USER:$SERVICE_USER" "$work_dir" 2>/dev/null || true
    fi
    
    # Build agent start command (without --daemon for simple service)
    local agent_cmd="$INSTALL_DIR/sloth-runner agent start"
    agent_cmd="$agent_cmd --name $AGENT_NAME"
    agent_cmd="$agent_cmd --master $MASTER_ADDRESS"
    agent_cmd="$agent_cmd --port $AGENT_PORT"
    agent_cmd="$agent_cmd --bind-address $BIND_ADDRESS"
    
    # Add report-address if specified
    if [ -n "$REPORT_ADDRESS" ]; then
        agent_cmd="$agent_cmd --report-address $REPORT_ADDRESS"
    fi
    
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
WorkingDirectory=$work_dir
Restart=always
RestartSec=5s
StartLimitInterval=60s
StartLimitBurst=5

# Agent Configuration
ExecStart=$agent_cmd

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=sloth-runner-agent

# Performance
LimitNOFILE=65536

# Security
NoNewPrivileges=true
PrivateTmp=true

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
    
    if [ "$UPDATE_MODE" = true ]; then
        info "Restarting sloth-runner-agent service..."
        if command -v sudo &> /dev/null && [ "$(id -u)" -ne 0 ]; then
            sudo systemctl restart sloth-runner-agent
        else
            systemctl restart sloth-runner-agent
        fi
        success "Service restarted"
    else
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
    fi
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
        if [ -n "$REPORT_ADDRESS" ]; then
            echo -e "    --report-address $REPORT_ADDRESS \\"
        fi
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
    if [ -n "$REPORT_ADDRESS" ]; then
        echo -e "  Report Address:${YELLOW}$REPORT_ADDRESS${NC}"
    fi
    echo -e "  User:          ${YELLOW}$SERVICE_USER${NC}"
    echo ""
    
    if [ "$SKIP_SYSTEMD" = false ]; then
        echo -e "${CYAN}🔧 Service Management:${NC}"
        echo -e "  sudo systemctl status sloth-runner-agent    # Check status"
        echo -e "  sudo systemctl restart sloth-runner-agent   # Restart agent"
        echo -e "  sudo systemctl stop sloth-runner-agent      # Stop agent"
        echo -e "  sudo journalctl -u sloth-runner-agent -f    # View logs"
        echo ""
    else
        echo -e "${CYAN}🔧 Agent Management (no systemd):${NC}"
        echo -e "  $INSTALL_DIR/sloth-runner agent start --name $AGENT_NAME --master $MASTER_ADDRESS --port $AGENT_PORT --bind-address $BIND_ADDRESS --daemon"
        echo -e "  pkill -f 'sloth-runner agent'    # Stop agent"
        echo -e "  ps aux | grep sloth-runner       # Check if running"
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

# Start agent directly if systemd not available
start_agent_directly() {
    if [ "$SKIP_SYSTEMD" = false ]; then
        return
    fi
    
    info "Starting agent directly (systemd not available)..."
    
    # Build command
    local cmd="$INSTALL_DIR/sloth-runner agent start"
    cmd="$cmd --name $AGENT_NAME"
    cmd="$cmd --master $MASTER_ADDRESS"
    cmd="$cmd --port $AGENT_PORT"
    cmd="$cmd --bind-address $BIND_ADDRESS"
    if [ -n "$REPORT_ADDRESS" ]; then
        cmd="$cmd --report-address $REPORT_ADDRESS"
    fi
    cmd="$cmd --daemon"
    
    # Start agent
    if $cmd; then
        success "Agent started successfully"
        sleep 2
        
        # Verify it's running
        if ps aux | grep -v grep | grep "sloth-runner agent" | grep "$AGENT_NAME" > /dev/null; then
            success "Agent is running!"
            ps aux | grep -v grep | grep "sloth-runner agent" | grep "$AGENT_NAME"
        else
            warn "Agent may not be running. Check logs with: cat agent.log"
        fi
    else
        error "Failed to start agent"
    fi
}

# Main bootstrap flow
main() {
    print_banner
    
    if [ "$UPDATE_MODE" = true ]; then
        info "🔄 Update mode: Updating sloth-runner agent..."
    else
        info "Starting Sloth Runner Agent bootstrap..."
    fi
    echo ""
    
    # Detect OS
    detect_os
    
    # Detect bind address (skip in update mode if service exists)
    if [ "$UPDATE_MODE" = false ]; then
        detect_bind_address
    fi
    
    # Install/Update sloth-runner
    install_sloth_runner
    echo ""
    
    # Create systemd service or start directly
    if [ "$SKIP_SYSTEMD" = false ]; then
        if [ "$UPDATE_MODE" = false ]; then
            create_systemd_service
            echo ""
        fi
        start_service
        echo ""
        verify_agent
    else
        if [ "$UPDATE_MODE" = false ]; then
            start_agent_directly
            echo ""
        else
            warn "Update complete. Please restart your agent manually."
        fi
    fi
    
    # Show post-install info
    if [ "$UPDATE_MODE" = true ]; then
        echo ""
        echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
        echo -e "${GREEN}    🎉 Sloth Runner Agent Update Complete!${NC}"
        echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
        echo ""
        if [ "$SKIP_SYSTEMD" = false ]; then
            echo -e "${CYAN}Service restarted and running with updated binary${NC}"
        fi
    else
        show_post_install
    fi
}

# Run main function
main "$@"
