#!/bin/bash
# Fix Incus container agent connectivity
# Usage: ./fix-incus-agent.sh <container-name> <host-ip> <agent-port>

set -e

CONTAINER_NAME=${1:-lady-arch}
HOST_IP=${2:-192.168.1.16}
AGENT_PORT=${3:-50052}
PROXY_DEVICE_NAME="sloth-agent-proxy"

echo "🔧 Fixing Incus container agent: $CONTAINER_NAME"
echo "   Host IP: $HOST_IP"
echo "   Agent Port: $AGENT_PORT"
echo ""

# Check if we're on the correct host
if [[ "$(hostname -I | grep -c "$HOST_IP")" -eq 0 ]]; then
    echo "⚠️  Warning: This script should run on the host ($HOST_IP)"
    echo "   Current IPs: $(hostname -I)"
    echo ""
    echo "Run this on the host or via SSH:"
    echo "   ssh user@$HOST_IP 'bash -s' < $0 $CONTAINER_NAME $HOST_IP $AGENT_PORT"
    exit 1
fi

# Check if incus is installed
if ! command -v incus &> /dev/null; then
    echo "❌ Error: incus not found on this host"
    exit 1
fi

# Check if container exists
if ! incus list | grep -q "$CONTAINER_NAME"; then
    echo "❌ Error: Container '$CONTAINER_NAME' not found"
    echo ""
    echo "Available containers:"
    incus list
    exit 1
fi

# Check if container is running
CONTAINER_STATUS=$(incus list "$CONTAINER_NAME" --format json | jq -r '.[0].status')
if [[ "$CONTAINER_STATUS" != "Running" ]]; then
    echo "⚠️  Container is not running (status: $CONTAINER_STATUS)"
    echo "   Starting container..."
    incus start "$CONTAINER_NAME"
    sleep 2
fi

# Check if proxy device already exists
if incus config show "$CONTAINER_NAME" | grep -q "$PROXY_DEVICE_NAME"; then
    echo "ℹ️  Proxy device '$PROXY_DEVICE_NAME' already exists"
    echo "   Removing old proxy..."
    incus config device remove "$CONTAINER_NAME" "$PROXY_DEVICE_NAME" || true
    sleep 1
fi

# Add proxy device
echo "📡 Adding proxy device to forward port $AGENT_PORT..."
incus config device add "$CONTAINER_NAME" "$PROXY_DEVICE_NAME" proxy \
    listen=tcp:0.0.0.0:$AGENT_PORT \
    connect=tcp:127.0.0.1:$AGENT_PORT

echo ""
echo "✅ Proxy device added successfully!"
echo ""

# Verify configuration
echo "📋 Container network configuration:"
incus config show "$CONTAINER_NAME" | grep -A 5 "devices:" || echo "No devices configured"
echo ""

# Get container IP for reference
CONTAINER_IP=$(incus list "$CONTAINER_NAME" --format json | jq -r '.[0].state.network.eth0.addresses[] | select(.family=="inet") | .address')
echo "📍 Container internal IP: $CONTAINER_IP"
echo "📍 Accessible via host: $HOST_IP:$AGENT_PORT"
echo ""

# Test if port is accessible
echo "🧪 Testing connectivity..."
if nc -zv -w 2 "$HOST_IP" "$AGENT_PORT" 2>&1 | grep -q "succeeded"; then
    echo "✅ Port $AGENT_PORT is accessible on $HOST_IP"
else
    echo "⚠️  Port $AGENT_PORT is not yet accessible"
    echo "   The agent might need to be restarted inside the container"
    echo ""
    echo "   Run inside container:"
    echo "   incus exec $CONTAINER_NAME -- systemctl restart sloth-runner-agent"
fi

echo ""
echo "🎉 Done! Test with:"
echo "   ./sloth-runner agent run $CONTAINER_NAME 'ls -la'"
echo ""

# Optional: Check if agent is running inside container
echo "🔍 Checking agent status inside container..."
if incus exec "$CONTAINER_NAME" -- systemctl is-active sloth-runner-agent &>/dev/null; then
    echo "✅ Agent service is running inside container"
else
    echo "⚠️  Agent service is not running inside container"
    echo ""
    echo "To start the agent:"
    echo "   incus exec $CONTAINER_NAME -- systemctl start sloth-runner-agent"
    echo ""
    echo "Or manually:"
    echo "   incus exec $CONTAINER_NAME -- /usr/local/bin/sloth-runner agent start \\"
    echo "       --name $CONTAINER_NAME \\"
    echo "       --master <MASTER_IP>:50053 \\"
    echo "       --port $AGENT_PORT \\"
    echo "       --bind-address 0.0.0.0 \\"
    echo "       --daemon"
fi

echo ""
echo "📚 For more details, see: AGENT_TROUBLESHOOTING.md"
