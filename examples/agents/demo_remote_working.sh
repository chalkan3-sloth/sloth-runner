#!/bin/bash

# Script demonstrando execução remota funcionando nos agents

echo "🚀 Demonstração de Execução Remota - Sloth Runner"
echo "==============================================="
echo ""

echo "📡 Listando agents registrados:"
sloth-runner agent list --master 192.168.1.29:50053
echo ""

echo "🔍 Executando LS em múltiplos agents..."
echo ""

echo "📍 Agent: ladyguica"
echo "-------------------"
sloth-runner agent run ladyguica "echo '🖥️ Executando em:' && hostname && echo '👤 Usuário:' && whoami && echo '📂 Arquivos:' && ls -la \$HOME | head -8" --master 192.168.1.29:50053

echo ""
echo "📍 Agent: keiteguica"  
echo "--------------------"
sloth-runner agent run keiteguica "echo '🖥️ Executando em:' && hostname && echo '👤 Usuário:' && whoami && echo '📂 Arquivos:' && ls -la \$HOME | head -8" --master 192.168.1.29:50053

echo ""
echo "✅ Execução remota concluída com sucesso!"
echo ""

echo "💡 Outros exemplos úteis:"
echo "========================"
echo ""

echo "📊 Informações do sistema:"
echo "sloth-runner agent run ladyguica \"uname -a && uptime && df -h /\" --master 192.168.1.29:50053"
echo ""

echo "🔧 Verificar processos:"
echo "sloth-runner agent run keiteguica \"ps aux | head -10\" --master 192.168.1.29:50053"
echo ""

echo "🌐 Verificar conectividade:"
echo "sloth-runner agent run ladyguica \"ping -c 3 google.com\" --master 192.168.1.29:50053"