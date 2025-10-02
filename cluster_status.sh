#!/bin/bash

# Agent Status Dashboard
MASTER_IP="192.168.1.29"
MASTER_PORT="50053"

echo "🎛️ SLOTH RUNNER CLUSTER STATUS DASHBOARD"
echo "========================================"
echo "Master: ${MASTER_IP}:${MASTER_PORT}"
echo "Timestamp: $(date)"
echo ""

echo "📡 MASTER STATUS:"
echo "----------------"
if ps aux | grep "sloth-runner master" | grep -q -v grep; then
    echo "✅ Master is running"
    echo "   PID: $(ps aux | grep "sloth-runner master" | grep -v grep | awk '{print $2}')"
    echo "   Port: ${MASTER_PORT}"
else
    echo "❌ Master is not running"
fi
echo ""

echo "👥 REGISTERED AGENTS:"
echo "-------------------"
./sloth-runner agent list --master ${MASTER_IP}:${MASTER_PORT} 2>/dev/null || echo "❌ Failed to connect to master"
echo ""

echo "🔍 AGENT PROCESSES:"
echo "-----------------"
echo "Agent 192.168.1.16:"
if ssh -o ConnectTimeout=3 chalkan3@192.168.1.16 "ps aux | grep sloth-runner | grep -v grep" 2>/dev/null; then
    echo "✅ Agent process running"
else
    echo "❌ No agent process found"
fi

echo ""
echo "Agent 192.168.1.17:"
if ssh -o ConnectTimeout=3 chalkan3@192.168.1.17 "ps aux | grep sloth-runner | grep -v grep" 2>/dev/null; then
    echo "✅ Agent process running"
else
    echo "❌ No agent process found"
fi

echo ""
echo "🧪 CONNECTIVITY TEST:"
echo "--------------------"
for ip in "192.168.1.16" "192.168.1.17"; do
    echo -n "Testing $ip: "
    if ping -c 1 -W 1 "$ip" > /dev/null 2>&1; then
        echo "✅ Reachable"
    else
        echo "❌ Unreachable"
    fi
done

echo ""
echo "📊 CLUSTER SUMMARY:"
echo "-----------------"
total_agents=$(./sloth-runner agent list --master ${MASTER_IP}:${MASTER_PORT} 2>/dev/null | grep -c "Active" || echo "0")
echo "Active Agents: $total_agents"
echo "Expected Agents: 2"
if [ "$total_agents" = "2" ]; then
    echo "Status: ✅ All agents connected"
elif [ "$total_agents" = "1" ]; then
    echo "Status: ⚠️ Partial connectivity"
else
    echo "Status: ❌ No agents connected"
fi

echo ""
echo "🎯 READY FOR DISTRIBUTED EXECUTION!"
echo "Example command:"
echo "  ./sloth-runner run -f test_distributed.sloth distributed_test_workflow"