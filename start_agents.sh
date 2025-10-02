#!/bin/bash

# Script para iniciar os agents em daemon nos hosts especificados

echo "🚀 Iniciando agents em daemon..."

# Agent ladyguica no 192.168.1.16
echo "📡 Iniciando agent ladyguica no 192.168.1.16:50051"
ssh chalkan3@192.168.1.16 "cd ~/.local/bin && ./sloth-runner agent start --name ladyguica --port 50051 --master 192.168.1.29:50053 --bind-address 192.168.1.16 --daemon"

# Agent keiteguica no 192.168.1.17
echo "📡 Iniciando agent keiteguica no 192.168.1.17:50051"
ssh chalkan3@192.168.1.17 "cd ~/.local/bin && ./sloth-runner agent start --name keiteguica --port 50051 --master 192.168.1.29:50053 --bind-address 192.168.1.17 --daemon"

echo "✅ Agents iniciados em daemon!"
echo ""
echo "Para verificar os agents registrados:"
echo "sloth-runner agent list --master 192.168.1.29:50053"