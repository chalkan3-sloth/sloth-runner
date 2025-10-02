#!/bin/bash

echo "🎯 DEMONSTRAÇÃO DE EXECUÇÃO REMOTA COM SLOTH-RUNNER"
echo "=================================================="
echo ""

# Verificar agents conectados
echo "📋 1. Verificando agents conectados..."
./sloth-runner agent list --master 192.168.1.29:50053
echo ""

# Executar comando direto em cada agent (FUNCIONA!)
echo "🚀 2. Executando comandos DIRETAMENTE nos agents (funciona!)..."
echo ""

echo "📂 2.1. Executando 'ls -la' no LADYGUICA (192.168.1.16):"
echo "---"
./sloth-runner agent run ladyguica "echo 'Executando no:' && hostname && echo 'Diretório atual:' && pwd && echo 'Conteúdo:' && ls -la | head -10" --master 192.168.1.29:50053 || echo "ERRO: Falha na execução"
echo ""

echo "📂 2.2. Executando 'ls -la' no KEITEGUICA (192.168.1.17):"
echo "---"
./sloth-runner agent run keiteguica "echo 'Executando no:' && hostname && echo 'Diretório atual:' && pwd && echo 'Conteúdo:' && ls -la | head -10" --master 192.168.1.29:50053 || echo "ERRO: Falha na execução"
echo ""

# Executar workflow com delegate_to (ainda executa local)
echo "⚠️  3. Executando workflow com :delegate_to() (ainda local)..."
echo "---"
./sloth-runner run -f examples/agents/simple_ls.sloth distributed_ls_workflow
echo ""

echo "📊 4. Verificando database SQLite..."
echo "---"
sqlite3 .sloth-cache/agents.db "SELECT 'Agent: ' || name || ', Address: ' || address || ', Status: ' || CASE WHEN (strftime('%s','now') - last_heartbeat) < 60 THEN 'ACTIVE' ELSE 'INACTIVE' END as status FROM agents;"
echo ""

echo "✅ DEMONSTRAÇÃO CONCLUÍDA!"
echo ""
echo "📝 RESULTADOS:"
echo "   ✅ Comando direto via 'agent run' → EXECUTA REMOTAMENTE"
echo "   ⚠️  Workflow com :delegate_to()    → Ainda executa localmente"
echo ""
echo "🔧 PRÓXIMO PASSO: Corrigir integração do AgentResolver no taskrunner"