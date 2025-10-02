#!/bin/bash
# 🚀 Script para executar LS em múltiplos agents remotos
# Uso: ./ls_multiple_agents.sh

set -e  # Para em caso de erro

MASTER="192.168.1.29:50053"

echo "🔍 Executando LS em todos os agents remotos..."
echo "================================================"

echo ""
echo "📍 Agent: ladyguica (192.168.1.16)"
echo "--------------------------------"
sloth-runner agent run ladyguica "hostname && echo 'Listagem de arquivos:' && ls -la \$HOME | head -15" --master $MASTER

echo ""
echo "📍 Agent: keiteguica (192.168.1.17)"  
echo "--------------------------------"
sloth-runner agent run keiteguica "hostname && echo 'Listagem de arquivos:' && ls -la \$HOME | head -15" --master $MASTER

echo ""
echo "✅ LS executado com sucesso em todos os agents!"
echo "==============================================="