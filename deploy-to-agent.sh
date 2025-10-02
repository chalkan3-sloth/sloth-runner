#!/bin/bash

# Deploy sloth-runner to remote agent
if [ -z "$1" ]; then
    echo "Usage: $0 <agent-ip>"
    exit 1
fi

AGENT_IP="$1"
USER="chalkan3"

echo "Building sloth-runner for Linux..."
GOOS=linux GOARCH=amd64 go build -o sloth-runner-linux ./cmd/sloth-runner

echo "Stopping agent service on $AGENT_IP..."
ssh $USER@$AGENT_IP "sudo systemctl stop sloth-runner-agent"

echo "Copying binary to $AGENT_IP..."
scp sloth-runner-linux $USER@$AGENT_IP:/tmp/sloth-runner

echo "Installing binary on $AGENT_IP..."
ssh $USER@$AGENT_IP "sudo mv /tmp/sloth-runner /usr/local/bin/sloth-runner && sudo chmod +x /usr/local/bin/sloth-runner"

echo "Starting agent service on $AGENT_IP..."
ssh $USER@$AGENT_IP "sudo systemctl start sloth-runner-agent"

echo "Checking service status..."
ssh $USER@$AGENT_IP "sudo systemctl status sloth-runner-agent"

echo "Done!"
