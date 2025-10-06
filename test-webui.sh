#!/bin/bash

# Start the web UI for testing
echo "Starting Sloth Runner Web UI..."
echo "Access it at: http://localhost:8080"
echo "Chart test page: http://localhost:8080/chart-test"
echo "Agent control: http://localhost:8080/agent-control"
echo ""
echo "Press Ctrl+C to stop"
echo ""

$HOME/.local/bin/sloth-runner ui
