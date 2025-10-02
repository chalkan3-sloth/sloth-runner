#!/bin/bash

# Example script showing how to run commands on specific agents
# This demonstrates direct agent execution for targeted commands

echo "🎯 SLOTH RUNNER - AGENT DIRECT EXECUTION EXAMPLES"
echo "=================================================="
echo ""

# Check if agents are available
echo "📋 Checking available agents..."
./sloth-runner agent list --master 192.168.1.29:50053
echo ""

# Execute ls command on ladyguica specifically
echo "📂 Running 'ls -la' on ladyguica agent..."
./sloth-runner agent run ladyguica "ls -la" --master 192.168.1.29:50053
echo ""

# Execute ls command on keiteguica specifically  
echo "📂 Running 'ls -la' on keiteguica agent..."
./sloth-runner agent run keiteguica "ls -la" --master 192.168.1.29:50053
echo ""

# Get system info from ladyguica
echo "📊 Getting system info from ladyguica..."
./sloth-runner agent run ladyguica "hostname && uptime && df -h /" --master 192.168.1.29:50053
echo ""

# Get system info from keiteguica
echo "📊 Getting system info from keiteguica..."
./sloth-runner agent run keiteguica "hostname && uptime && df -h /" --master 192.168.1.29:50053
echo ""

# Run different commands in parallel using background jobs
echo "🚀 Running parallel commands on both agents..."
./sloth-runner agent run ladyguica "ps aux | head -10" --master 192.168.1.29:50053 &
./sloth-runner agent run keiteguica "free -h" --master 192.168.1.29:50053 &

# Wait for background jobs to complete
wait

echo ""
echo "✅ Agent command execution examples completed!"
echo ""
echo "💡 Tips:"
echo "   - Use 'agent run' for specific agent targeting"
echo "   - Use workflow files for complex multi-step operations"
echo "   - Combine both approaches for maximum flexibility"