#!/bin/bash

# 🚀 SLOTH-RUNNER REMOTE EXECUTION DEMO
# This script demonstrates the current working functionality

echo "🔥 SLOTH-RUNNER REMOTE EXECUTION DEMO"
echo "======================================="
echo

# Check master connection
echo "1️⃣ Testing Master Connection..."
sloth-runner agent list --master 192.168.1.29:50053
echo

# Test basic remote commands
echo "2️⃣ Testing Basic Remote Commands..."
echo "📍 Executing 'hostname' on ladyguica:"
sloth-runner agent run ladyguica "hostname" --master 192.168.1.29:50053
echo

echo "📍 Executing 'whoami' on keiteguica:"  
sloth-runner agent run keiteguica "whoami" --master 192.168.1.29:50053
echo

echo "📍 Executing 'ls -la $HOME | head -3' on ladyguica:"
sloth-runner agent run ladyguica "ls -la \$HOME | head -3" --master 192.168.1.29:50053
echo

# Test delegate_to (in development)
echo "3️⃣ Testing delegate_to (IN DEVELOPMENT)..."
echo "📋 Running delegate_to example:"
echo "Command: sloth-runner run -f examples/agents/legacy_syntax_delegate.sloth simple_remote_test"
echo "Status: Script transmission works, execution debugging in progress"
echo

echo "✅ DEMO COMPLETE"
echo "💡 Use 'sloth-runner agent run AGENT_NAME \"command\"' for reliable remote execution"
echo "🔧 delegate_to functionality is implemented and being debugged"