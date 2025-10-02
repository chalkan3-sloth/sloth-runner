#!/bin/bash

# Script para conectar agents remotos ao master
MASTER_IP="192.168.1.29"
MASTER_PORT="50053"

AGENTS=(
    "192.168.1.16"
    "192.168.1.17"
)

echo "🎯 CONECTANDO AGENTS AO MASTER"
echo "Master: ${MASTER_IP}:${MASTER_PORT}"
echo ""

for agent_ip in "${AGENTS[@]}"; do
    echo "🔗 Conectando agent: $agent_ip"
    
    # Teste de conectividade
    if ping -c 1 -W 3 "$agent_ip" > /dev/null 2>&1; then
        echo "  ✅ Host $agent_ip acessível"
        
        # Verificar se sloth-runner existe no host remoto
        if ssh "chalkan3@$agent_ip" "which sloth-runner" > /dev/null 2>&1; then
            echo "  ✅ sloth-runner encontrado em $agent_ip"
            
            # Iniciar agent remoto
            echo "  🚀 Iniciando agent em $agent_ip..."
            ssh "chalkan3@$agent_ip" "sloth-runner agent start --master ${MASTER_IP}:${MASTER_PORT} --daemon" &
            
            # Dar tempo para conectar
            sleep 3
            
            echo "  ✅ Agent $agent_ip iniciado"
        else
            echo "  ❌ sloth-runner não encontrado em $agent_ip"
            echo "  📋 Tentando copiar binário..."
            
            # Copiar binário se não existir
            if scp ./sloth-runner "chalkan3@$agent_ip:~/sloth-runner"; then
                ssh "chalkan3@$agent_ip" "sudo mv ~/sloth-runner /usr/local/bin/ && sudo chmod +x /usr/local/bin/sloth-runner"
                echo "  ✅ Binário copiado e instalado em $agent_ip"
                
                # Tentar iniciar novamente
                echo "  🚀 Iniciando agent em $agent_ip..."
                ssh "chalkan3@$agent_ip" "sloth-runner agent start --master ${MASTER_IP}:${MASTER_PORT} --daemon" &
                sleep 3
                echo "  ✅ Agent $agent_ip iniciado"
            else
                echo "  ❌ Falha ao copiar binário para $agent_ip"
            fi
        fi
    else
        echo "  ❌ Host $agent_ip não acessível"
    fi
    
    echo ""
done

echo "⏱️ Aguardando agents se registrarem..."
sleep 5

echo "📋 LISTANDO AGENTS REGISTRADOS:"
./sloth-runner agent list --master ${MASTER_IP}:${MASTER_PORT}