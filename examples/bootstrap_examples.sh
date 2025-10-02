#!/usr/bin/env bash
#
# Example: Install Sloth Runner Agent via Bootstrap
#
# This script demonstrates various ways to use bootstrap.sh
#

set -e

echo "=== Sloth Runner Agent Bootstrap Examples ==="
echo ""

# Example 1: Basic installation
echo "1. Basic Installation:"
echo "   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \\"
echo "     --name myagent"
echo ""

# Example 2: Production setup
echo "2. Production Setup:"
echo "   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \\"
echo "     --name production-agent-01 \\"
echo "     --master 192.168.1.10:50053 \\"
echo "     --port 50051 \\"
echo "     --bind-address 192.168.1.20 \\"
echo "     --user slothrunner"
echo ""

# Example 3: Development setup (no systemd)
echo "3. Development Setup (No Systemd):"
echo "   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \\"
echo "     --name dev-agent \\"
echo "     --no-sudo \\"
echo "     --no-systemd"
echo ""

# Example 4: Specific version
echo "4. Install Specific Version:"
echo "   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \\"
echo "     --name myagent \\"
echo "     --version v3.23.1"
echo ""

# Example 5: Multiple agents on same host
echo "5. Multiple Agents on Same Host:"
echo "   # Agent 1"
echo "   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \\"
echo "     --name agent-01 \\"
echo "     --port 50051"
echo ""
echo "   # Agent 2"
echo "   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \\"
echo "     --name agent-02 \\"
echo "     --port 50052"
echo ""

# Example 6: Remote installation via SSH
echo "6. Remote Installation via SSH:"
echo "   ssh user@remote-host \"bash <(curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh) \\"
echo "     --name remote-agent \\"
echo "     --master 192.168.1.10:50053 \\"
echo "     --bind-address 192.168.1.20\""
echo ""

# Example 7: Vagrant installation
echo "7. Vagrant Installation:"
echo "   vagrant ssh -c \"bash <(curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh) \\"
echo "     --name vagrant-agent \\"
echo "     --master 192.168.1.10:50053\""
echo ""

echo "=== Post-Installation Commands ==="
echo ""
echo "Check agent status:"
echo "  sudo systemctl status sloth-runner-agent"
echo ""
echo "View logs:"
echo "  sudo journalctl -u sloth-runner-agent -f"
echo ""
echo "List agents on master:"
echo "  sloth-runner agent list"
echo ""
echo "Test agent:"
echo "  sloth-runner agent run myagent \"hostname\""
echo ""
echo "Delete agent:"
echo "  sloth-runner agent delete myagent --yes"
echo ""
