#!/bin/bash

# DEMO FINAL: Execução Remota Funcionando
# Este script demonstra como executar comandos remotos nos agentes

echo "🚀 DEMO: Execução Remota Distribuída"
echo "=================================="
echo

# Configuração
MASTER="192.168.1.29:50053"

echo "📋 1. Listando agentes registrados..."
sloth-runner agent list --master $MASTER
echo

echo "🏃 2. Executando hostname em cada agente..."
echo
echo "🔍 ladyguica:"
sloth-runner agent run ladyguica "hostname && echo 'Executado no ladyguica!'" --master $MASTER
echo

echo "🔍 keiteguica:"
sloth-runner agent run keiteguica "hostname && echo 'Executado no keiteguica!'" --master $MASTER
echo

echo "📂 3. Listando arquivos home em cada agente..."
echo
echo "📁 ladyguica:"
sloth-runner agent run ladyguica "ls -la \$HOME | head -10" --master $MASTER
echo

echo "📁 keiteguica:" 
sloth-runner agent run keiteguica "ls -la \$HOME | head -10" --master $MASTER
echo

echo "💾 4. Informações do sistema..."
echo
echo "💻 ladyguica:"
sloth-runner agent run ladyguica "uname -a && date" --master $MASTER
echo

echo "💻 keiteguica:"
sloth-runner agent run keiteguica "uname -a && date" --master $MASTER
echo

echo "✅ DEMO CONCLUÍDO!"
echo "Use 'sloth-runner agent run AGENT_NAME \"comando\"' para execução remota"