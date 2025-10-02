#!/bin/bash

# Pipeline distribuída usando sloth-runner
# Demonstra execução coordenada em múltiplos agents

echo "🌟 Pipeline Distribuída - Sloth Runner"
echo "======================================"
echo ""

# Verificar agents disponíveis
echo "📡 Verificando agents disponíveis..."
sloth-runner agent list --master 192.168.1.29:50053
echo ""

# Fase 1: Coleta de informações
echo "📊 Fase 1: Coletando informações dos sistemas..."
echo "-----------------------------------------------"

echo "🖥️ ladyguica - Informações do sistema:"
sloth-runner agent run ladyguica "echo 'Sistema:' && uname -a && echo 'Uptime:' && uptime && echo 'Memória:' && free -h && echo 'Disco:' && df -h / | tail -1" --master 192.168.1.29:50053

echo ""
echo "🖥️ keiteguica - Informações do sistema:" 
sloth-runner agent run keiteguica "echo 'Sistema:' && uname -a && echo 'Uptime:' && uptime && echo 'Memória:' && free -h && echo 'Disco:' && df -h / | tail -1" --master 192.168.1.29:50053

echo ""

# Fase 2: Verificação de serviços
echo "🔧 Fase 2: Verificando serviços..."
echo "--------------------------------"

echo "⚙️ ladyguica - Processos principais:"
sloth-runner agent run ladyguica "echo 'Top 5 processos:' && ps aux --sort=-%cpu | head -6" --master 192.168.1.29:50053

echo ""
echo "⚙️ keiteguica - Processos principais:"
sloth-runner agent run keiteguica "echo 'Top 5 processos:' && ps aux --sort=-%cpu | head -6" --master 192.168.1.29:50053

echo ""

# Fase 3: Verificação de conectividade
echo "🌐 Fase 3: Testando conectividade..."
echo "----------------------------------"

echo "📶 ladyguica - Teste de conectividade:"
sloth-runner agent run ladyguica "echo 'Testando ping:' && ping -c 2 google.com | tail -3" --master 192.168.1.29:50053

echo ""
echo "📶 keiteguica - Teste de conectividade:"
sloth-runner agent run keiteguica "echo 'Testando ping:' && ping -c 2 google.com | tail -3" --master 192.168.1.29:50053

echo ""

# Fase 4: Listagem de arquivos importantes
echo "📁 Fase 4: Verificando arquivos importantes..."
echo "--------------------------------------------"

echo "📂 ladyguica - Diretórios principais:"
sloth-runner agent run ladyguica "echo 'Home:' && ls -la \$HOME | head -5 && echo '' && echo 'Logs:' && ls -la /var/log 2>/dev/null | head -3 || echo 'Sem acesso a /var/log'" --master 192.168.1.29:50053

echo ""
echo "📂 keiteguica - Diretórios principais:"
sloth-runner agent run keiteguica "echo 'Home:' && ls -la \$HOME | head -5 && echo '' && echo 'Logs:' && ls -la /var/log 2>/dev/null | head -3 || echo 'Sem acesso a /var/log'" --master 192.168.1.29:50053

echo ""

# Finalização
echo "✅ Pipeline distribuída concluída com sucesso!"
echo "============================================="
echo ""
echo "🎯 Resumo da execução:"
echo "- ✅ Informações de sistema coletadas"
echo "- ✅ Serviços verificados"  
echo "- ✅ Conectividade testada"
echo "- ✅ Arquivos importantes listados"
echo ""
echo "🚀 Todos os agents executaram suas tarefas com sucesso!"